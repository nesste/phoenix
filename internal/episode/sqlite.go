package episode

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gowebpki/jcs"
	"github.com/nesste/phoenix/internal/redact"
	_ "modernc.org/sqlite"
)

const schemaVersion = 1

var (
	idPattern     = regexp.MustCompile(`^[a-z]_[A-Za-z0-9_-]+$`)
	digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	ErrNotFound   = errors.New("episode not found")
)

type Options struct {
	Secrets     []string
	RecallLimit int
	Now         func() time.Time
}

type Store struct {
	db          *sql.DB
	redactor    redact.Redactor
	recallLimit int
	now         func() time.Time
}

type startedPayload struct {
	WorldBuild      string          `json:"world_build"`
	Request         any             `json:"request"`
	RequestDigest   string          `json:"request_digest"`
	SuggestionTaken *SuggestionLink `json:"suggestion_taken"`
}

type completedPayload struct {
	Result       any    `json:"result"`
	ResultDigest string `json:"result_digest"`
}

type terminalPayload struct {
	Reason  string `json:"reason,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

func Open(path string, options Options) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("episode database path is required")
	}
	if options.RecallLimit <= 0 {
		options.RecallLimit = 20
	}
	if options.Now == nil {
		options.Now = time.Now
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open episode database: %w", err)
	}
	db.SetMaxOpenConns(1)
	store := &Store{
		db: db, redactor: redact.New(options.Secrets),
		recallLimit: options.RecallLimit, now: options.Now,
	}
	if err := store.initialize(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (store *Store) initialize(ctx context.Context) error {
	for _, statement := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=FULL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
	} {
		if _, err := store.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("configure episode database: %w", err)
		}
	}
	if _, err := store.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS episodes (
			episode_id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL UNIQUE,
			world_build TEXT NOT NULL,
			started_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS episode_events (
			sequence INTEGER PRIMARY KEY AUTOINCREMENT,
			episode_id TEXT NOT NULL REFERENCES episodes(episode_id),
			act_id TEXT NOT NULL,
			kind TEXT NOT NULL CHECK(kind IN ('act_started','act_completed','act_interrupted','episode_completed','episode_interrupted')),
			payload BLOB NOT NULL,
			payload_digest TEXT NOT NULL,
			created_at TEXT NOT NULL,
			UNIQUE(episode_id, act_id, kind)
		);
		CREATE INDEX IF NOT EXISTS episode_events_lookup ON episode_events(episode_id, sequence);
	`); err != nil {
		return fmt.Errorf("create episode schema: %w", err)
	}
	if err := store.recoverInterrupted(ctx); err != nil {
		return fmt.Errorf("recover interrupted episodes: %w", err)
	}
	return nil
}

func (store *Store) StartEpisode(ctx context.Context, sessionID, worldBuild string) (string, error) {
	if !idPattern.MatchString(sessionID) {
		return "", fmt.Errorf("session id is invalid")
	}
	if !digestPattern.MatchString(worldBuild) {
		return "", fmt.Errorf("world build must be a sha256 digest")
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin episode: %w", err)
	}
	defer tx.Rollback()
	var existingID, existingBuild string
	err = tx.QueryRowContext(ctx, "SELECT episode_id, world_build FROM episodes WHERE session_id = ?", sessionID).Scan(&existingID, &existingBuild)
	if err == nil {
		if existingBuild != worldBuild {
			return "", fmt.Errorf("session already belongs to a different world build")
		}
		return existingID, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("find session episode: %w", err)
	}
	episodeID := "e_" + rand.Text()
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO episodes(episode_id, session_id, world_build, started_at) VALUES(?, ?, ?, ?)",
		episodeID, sessionID, worldBuild, formatTime(store.now()),
	); err != nil {
		return "", fmt.Errorf("insert episode: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit episode: %w", err)
	}
	return episodeID, nil
}

func (store *Store) BeginAct(ctx context.Context, act StartedAct) error {
	if err := validateStartedAct(act); err != nil {
		return err
	}
	var episodeBuild string
	if err := store.db.QueryRowContext(ctx, "SELECT world_build FROM episodes WHERE episode_id = ?", act.EpisodeID).Scan(&episodeBuild); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if err != nil {
		return fmt.Errorf("read episode world build: %w", err)
	}
	if episodeBuild != act.WorldBuild {
		return fmt.Errorf("act world build does not match its episode")
	}
	request, _, requestDigest, err := store.prepare(act.Request)
	if err != nil {
		return fmt.Errorf("prepare act request: %w", err)
	}
	payload := startedPayload{
		WorldBuild: act.WorldBuild, Request: request, RequestDigest: requestDigest,
		SuggestionTaken: act.SuggestionTaken,
	}
	return store.appendActEvent(ctx, act.EpisodeID, act.ActID, "act_started", payload, true)
}

func (store *Store) CompleteAct(ctx context.Context, act CompletedAct) error {
	if !idPattern.MatchString(act.EpisodeID) || !idPattern.MatchString(act.ActID) {
		return fmt.Errorf("episode or act id is invalid")
	}
	result, _, resultDigest, err := store.prepare(act.Result)
	if err != nil {
		return fmt.Errorf("prepare act result: %w", err)
	}
	payload := completedPayload{Result: result, ResultDigest: resultDigest}
	return store.appendActEvent(ctx, act.EpisodeID, act.ActID, "act_completed", payload, false)
}

func (store *Store) FinishEpisode(ctx context.Context, episodeID, outcome string) error {
	if !idPattern.MatchString(episodeID) || strings.TrimSpace(outcome) == "" {
		return fmt.Errorf("episode id and outcome are required")
	}
	return store.appendActEvent(ctx, episodeID, "", "episode_completed", terminalPayload{Outcome: outcome}, false)
}

func (store *Store) appendActEvent(ctx context.Context, episodeID, actID, kind string, payload any, allowNewAct bool) error {
	_, encoded, payloadDigest, err := store.prepare(payload)
	if err != nil {
		return fmt.Errorf("prepare %s event: %w", kind, err)
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin %s event: %w", kind, err)
	}
	defer tx.Rollback()
	exists, same, err := existingEvent(ctx, tx, episodeID, actID, kind, encoded)
	if err != nil {
		return err
	}
	if exists {
		if !same {
			return fmt.Errorf("%s event already exists with a different payload", kind)
		}
		return tx.Commit()
	}
	if err := validateEventTarget(ctx, tx, episodeID, actID, kind, allowNewAct); err != nil {
		return err
	}
	if err := appendEvent(ctx, tx, episodeID, actID, kind, encoded, payloadDigest, store.now()); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit %s event: %w", kind, err)
	}
	return nil
}

func validateEventTarget(ctx context.Context, tx *sql.Tx, episodeID, actID, kind string, allowNewAct bool) error {
	var exists int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM episodes WHERE episode_id = ?", episodeID).Scan(&exists); err != nil {
		return fmt.Errorf("find episode: %w", err)
	}
	if exists == 0 {
		return ErrNotFound
	}
	if kind == "episode_completed" {
		if err := ensureNoActiveActs(ctx, tx, episodeID); err != nil {
			return err
		}
	}
	if kind != "episode_completed" {
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM episode_events
			WHERE episode_id = ? AND act_id = ? AND kind = 'act_started'`, episodeID, actID).Scan(&exists); err != nil {
			return fmt.Errorf("find act start: %w", err)
		}
		if allowNewAct && exists == 0 {
			return ensureEpisodeActive(ctx, tx, episodeID)
		}
		if exists == 0 {
			return fmt.Errorf("act %q has no start event", actID)
		}
		if kind == "act_completed" {
			var interrupted int
			if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM episode_events
				WHERE episode_id = ? AND act_id = ? AND kind = 'act_interrupted'`, episodeID, actID).Scan(&interrupted); err != nil {
				return fmt.Errorf("find interrupted act: %w", err)
			}
			if interrupted != 0 {
				return fmt.Errorf("act %q is already interrupted", actID)
			}
		}
	}
	return ensureEpisodeActive(ctx, tx, episodeID)
}

func ensureNoActiveActs(ctx context.Context, tx *sql.Tx, episodeID string) error {
	var active int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*)
		FROM episode_events AS started
		WHERE started.episode_id = ? AND started.kind = 'act_started'
		  AND NOT EXISTS (
		    SELECT 1 FROM episode_events AS terminal
		    WHERE terminal.episode_id = started.episode_id AND terminal.act_id = started.act_id
		      AND terminal.kind IN ('act_completed','act_interrupted')
		  )`, episodeID).Scan(&active); err != nil {
		return fmt.Errorf("find active acts: %w", err)
	}
	if active != 0 {
		return fmt.Errorf("episode %q has active acts", episodeID)
	}
	return nil
}

func ensureEpisodeActive(ctx context.Context, tx *sql.Tx, episodeID string) error {
	var terminal int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM episode_events
		WHERE episode_id = ? AND kind IN ('episode_completed','episode_interrupted')`, episodeID).Scan(&terminal); err != nil {
		return fmt.Errorf("find episode terminal event: %w", err)
	}
	if terminal != 0 {
		return fmt.Errorf("episode %q is already terminal", episodeID)
	}
	return nil
}

func appendEvent(ctx context.Context, tx *sql.Tx, episodeID, actID, kind string, payload []byte, payloadDigest string, now time.Time) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO episode_events
		(episode_id, act_id, kind, payload, payload_digest, created_at) VALUES(?, ?, ?, ?, ?, ?)`,
		episodeID, actID, kind, payload, payloadDigest, formatTime(now),
	); err != nil {
		return fmt.Errorf("append %s event: %w", kind, err)
	}
	return nil
}

func existingEvent(ctx context.Context, tx *sql.Tx, episodeID, actID, kind string, payload []byte) (bool, bool, error) {
	var existing []byte
	err := tx.QueryRowContext(ctx, `SELECT payload FROM episode_events
		WHERE episode_id = ? AND act_id = ? AND kind = ?`, episodeID, actID, kind).Scan(&existing)
	if errors.Is(err, sql.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, fmt.Errorf("find existing %s event: %w", kind, err)
	}
	return true, bytes.Equal(existing, payload), nil
}

func (store *Store) recoverInterrupted(ctx context.Context) error {
	interrupted := terminalPayload{Reason: "daemon restarted before completion"}
	_, payload, payloadDigest, err := store.prepare(interrupted)
	if err != nil {
		return err
	}
	now := formatTime(store.now())
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO episode_events(episode_id, act_id, kind, payload, payload_digest, created_at)
		SELECT started.episode_id, started.act_id, 'act_interrupted', ?, ?, ?
		FROM episode_events AS started
		WHERE started.kind = 'act_started'
		  AND NOT EXISTS (
		    SELECT 1 FROM episode_events AS terminal
		    WHERE terminal.episode_id = started.episode_id AND terminal.act_id = started.act_id
		      AND terminal.kind IN ('act_completed','act_interrupted')
		  )`, payload, payloadDigest, now); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO episode_events(episode_id, act_id, kind, payload, payload_digest, created_at)
		SELECT episodes.episode_id, '', 'episode_interrupted', ?, ?, ?
		FROM episodes
		WHERE NOT EXISTS (
		  SELECT 1 FROM episode_events AS terminal
		  WHERE terminal.episode_id = episodes.episode_id
		    AND terminal.kind IN ('episode_completed','episode_interrupted')
		)`, payload, payloadDigest, now); err != nil {
		return err
	}
	return tx.Commit()
}

func (store *Store) prepare(value any) (any, []byte, string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, nil, "", err
	}
	var normalized any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		return nil, nil, "", err
	}
	normalized = store.redactor.Value(normalized)
	encoded, err = json.Marshal(normalized)
	if err != nil {
		return nil, nil, "", err
	}
	canonical, err := jcs.Transform(encoded)
	if err != nil {
		return nil, nil, "", err
	}
	return normalized, canonical, digestBytes(canonical), nil
}

func validateStartedAct(act StartedAct) error {
	if !idPattern.MatchString(act.EpisodeID) || !idPattern.MatchString(act.ActID) {
		return fmt.Errorf("episode or act id is invalid")
	}
	if !digestPattern.MatchString(act.WorldBuild) {
		return fmt.Errorf("world build must be a sha256 digest")
	}
	if act.Request.Handle == "" || act.Request.Verb == "" || act.Request.Args == nil {
		return fmt.Errorf("structured act request is incomplete")
	}
	if act.SuggestionTaken != nil && (!idPattern.MatchString(act.SuggestionTaken.ActID) || act.SuggestionTaken.Index < 0 || act.SuggestionTaken.Index > 2) {
		return fmt.Errorf("suggestion linkage is invalid")
	}
	return nil
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func digestBytes(contents []byte) string {
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (store *Store) Close() error {
	if store == nil || store.db == nil {
		return nil
	}
	return store.db.Close()
}
