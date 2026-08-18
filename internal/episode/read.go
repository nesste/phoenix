package episode

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (store *Store) Episode(ctx context.Context, episodeID string) (Record, error) {
	var record Record
	var startedAt string
	err := store.db.QueryRowContext(ctx, `SELECT episode_id, session_id, world_build, started_at
		FROM episodes WHERE episode_id = ?`, episodeID).Scan(
		&record.EpisodeID, &record.SessionID, &record.WorldBuild, &startedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("read episode: %w", err)
	}
	record.V = schemaVersion
	record.Status = EpisodeActive
	record.StartedAt, err = parseTime(startedAt)
	if err != nil {
		return Record{}, err
	}
	record.Acts = []ActRecord{}
	rows, err := store.db.QueryContext(ctx, `SELECT sequence, act_id, kind, payload, created_at
		FROM episode_events WHERE episode_id = ? ORDER BY sequence`, episodeID)
	if err != nil {
		return Record{}, fmt.Errorf("read episode events: %w", err)
	}
	defer rows.Close()
	actIndexes := make(map[string]int)
	for rows.Next() {
		var sequence int64
		var actID, kind, createdAt string
		var payload []byte
		if err := rows.Scan(&sequence, &actID, &kind, &payload, &createdAt); err != nil {
			return Record{}, fmt.Errorf("scan episode event: %w", err)
		}
		if err := foldEvent(&record, actIndexes, sequence, actID, kind, payload, createdAt); err != nil {
			return Record{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return Record{}, fmt.Errorf("iterate episode events: %w", err)
	}
	return record, nil
}

// Latest returns the most recently started episode. Evaluation runners use a
// fresh database per trial, so this identifies the one isolated session
// without exposing episode contents through the agent-facing recall verb.
func (store *Store) Latest(ctx context.Context) (Record, error) {
	var episodeID string
	err := store.db.QueryRowContext(ctx, `SELECT episode_id FROM episodes ORDER BY started_at DESC, episode_id DESC LIMIT 1`).Scan(&episodeID)
	if errors.Is(err, sql.ErrNoRows) {
		return Record{}, ErrNotFound
	}
	if err != nil {
		return Record{}, fmt.Errorf("find latest episode: %w", err)
	}
	return store.Episode(ctx, episodeID)
}

func foldEvent(record *Record, actIndexes map[string]int, sequence int64, actID, kind string, payload []byte, createdAt string) error {
	eventTime, err := parseTime(createdAt)
	if err != nil {
		return err
	}
	switch kind {
	case "act_started":
		var started startedPayload
		if err := json.Unmarshal(payload, &started); err != nil {
			return fmt.Errorf("decode act start: %w", err)
		}
		request, err := decodeRequest(started.Request)
		if err != nil {
			return err
		}
		actIndexes[actID] = len(record.Acts)
		record.Acts = append(record.Acts, ActRecord{
			Sequence: sequence, ActID: actID, WorldBuild: started.WorldBuild,
			Status: ActStarted, Request: request, RequestDigest: started.RequestDigest,
			SuggestionTaken: started.SuggestionTaken, StartedAt: eventTime,
		})
	case "act_completed":
		var completed completedPayload
		if err := json.Unmarshal(payload, &completed); err != nil {
			return fmt.Errorf("decode act completion: %w", err)
		}
		index, exists := actIndexes[actID]
		if !exists {
			return fmt.Errorf("completion for unknown act %q", actID)
		}
		record.Acts[index].Status = ActCompleted
		record.Acts[index].Result = completed.Result
		record.Acts[index].ResultDigest = completed.ResultDigest
		record.Acts[index].CompletedAt = &eventTime
	case "act_interrupted":
		index, exists := actIndexes[actID]
		if !exists {
			return fmt.Errorf("interruption for unknown act %q", actID)
		}
		record.Acts[index].Status = ActInterrupted
		record.Acts[index].CompletedAt = &eventTime
	case "episode_completed", "episode_interrupted":
		var terminal terminalPayload
		if err := json.Unmarshal(payload, &terminal); err != nil {
			return fmt.Errorf("decode episode terminal event: %w", err)
		}
		if kind == "episode_completed" {
			record.Status = EpisodeCompleted
			record.Outcome = terminal.Outcome
		} else {
			record.Status = EpisodeInterrupted
		}
		record.EndedAt = &eventTime
	}
	return nil
}

func (store *Store) Recall(ctx context.Context, sessionID, query string) ([]Pointer, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []Pointer{}, nil
	}
	pattern := "%" + escapeLike(strings.ToLower(query)) + "%"
	rows, err := store.db.QueryContext(ctx, `
		SELECT starts.episode_id, starts.act_id
		FROM episode_events AS starts
		JOIN episode_events AS completed
		  ON completed.episode_id = starts.episode_id AND completed.act_id = starts.act_id
		JOIN episodes ON episodes.episode_id = starts.episode_id
		WHERE episodes.session_id = ?
		  AND starts.kind = 'act_started' AND completed.kind = 'act_completed'
		  AND (lower(CAST(starts.payload AS TEXT)) LIKE ? ESCAPE '\'
		       OR lower(CAST(completed.payload AS TEXT)) LIKE ? ESCAPE '\')
		ORDER BY completed.sequence DESC LIMIT ?`, sessionID, pattern, pattern, store.recallLimit)
	if err != nil {
		return nil, fmt.Errorf("recall episodes: %w", err)
	}
	defer rows.Close()
	pointers := []Pointer{}
	for rows.Next() {
		var pointer Pointer
		if err := rows.Scan(&pointer.EpisodeID, &pointer.ActID); err != nil {
			return nil, fmt.Errorf("scan recalled episode: %w", err)
		}
		pointers = append(pointers, pointer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recalled episodes: %w", err)
	}
	return pointers, nil
}

func decodeRequest(value any) (ActRequest, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return ActRequest{}, err
	}
	var request ActRequest
	if err := json.Unmarshal(encoded, &request); err != nil {
		return ActRequest{}, fmt.Errorf("decode act request: %w", err)
	}
	return request, nil
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "%", "\\%")
	return strings.ReplaceAll(value, "_", "\\_")
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse episode timestamp: %w", err)
	}
	return parsed, nil
}
