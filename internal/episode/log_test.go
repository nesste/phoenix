package episode

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

const testWorldBuild = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestAppendOnlyEpisodeReconstructsAndRecoversInterruptedAct(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "episodes.db")
	store, err := Open(path, Options{})
	mustSucceed(t, err)
	episodeID, err := store.StartEpisode(ctx, "s_session_one", testWorldBuild)
	mustSucceed(t, err)
	first := StartedAct{
		EpisodeID: episodeID, ActID: "a_first", WorldBuild: testWorldBuild,
		Request:         ActRequest{Handle: "h_repo", Verb: "edit", Args: map[string]any{"path": "main.go"}},
		SuggestionTaken: &SuggestionLink{ActID: "a_prior", Index: 1},
	}
	mustSucceed(t, store.BeginAct(ctx, first))
	result := map[string]any{
		"status":   "ok",
		"handles":  map[string]any{"grant": []any{}, "revoke": []any{"h_old"}},
		"frontier": []any{map[string]any{"call": map[string]any{"handle": "h_tests", "verb": "run", "args": map[string]any{}}, "why": "verify", "provenance": "authored", "score": 1}},
	}
	finished := CompletedAct{EpisodeID: episodeID, ActID: first.ActID, Result: result}
	mustSucceed(t, store.CompleteAct(ctx, finished))
	// Retrying the identical append must not create another event.
	mustSucceed(t, store.BeginAct(ctx, first))
	mustSucceed(t, store.CompleteAct(ctx, finished))
	mustSucceed(t, store.BeginAct(ctx, StartedAct{
		EpisodeID: episodeID, ActID: "a_half_written", WorldBuild: testWorldBuild,
		Request: ActRequest{Handle: "h_tests", Verb: "run", Args: map[string]any{}},
	}))
	mustSucceed(t, store.Close())

	reopened, err := Open(path, Options{})
	mustSucceed(t, err)
	defer reopened.Close()
	record, err := reopened.Episode(ctx, episodeID)
	mustSucceed(t, err)
	if record.SessionID != "s_session_one" || record.WorldBuild != testWorldBuild || record.Status != EpisodeInterrupted {
		t.Fatalf("episode identity/status = %#v", record)
	}
	if got, want := len(record.Acts), 2; got != want {
		t.Fatalf("act count = %d, want %d: %#v", got, want, record.Acts)
	}
	if record.Acts[0].Status != ActCompleted || record.Acts[1].Status != ActInterrupted {
		t.Fatalf("act statuses = %q, %q", record.Acts[0].Status, record.Acts[1].Status)
	}
	if record.Acts[0].Sequence >= record.Acts[1].Sequence {
		t.Fatalf("append order was not preserved: %#v", record.Acts)
	}
	if !reflect.DeepEqual(record.Acts[0].SuggestionTaken, first.SuggestionTaken) {
		t.Fatalf("suggestion link = %#v, want %#v", record.Acts[0].SuggestionTaken, first.SuggestionTaken)
	}
	if record.Acts[0].RequestDigest == "" || record.Acts[0].ResultDigest == "" {
		t.Fatalf("act digests are missing: %#v", record.Acts[0])
	}
	storedResult := record.Acts[0].Result.(map[string]any)
	if _, ok := storedResult["handles"]; !ok {
		t.Fatalf("handle delta is not reconstructable: %#v", storedResult)
	}
	if _, ok := storedResult["frontier"]; !ok {
		t.Fatalf("frontier is not reconstructable: %#v", storedResult)
	}
}

func TestIdempotentAppendRejectsActIdentityCollision(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, Options{})
	episodeID, err := store.StartEpisode(ctx, "s_collision", testWorldBuild)
	if err != nil {
		t.Fatal(err)
	}
	start := StartedAct{
		EpisodeID: episodeID, ActID: "a_retry", WorldBuild: testWorldBuild,
		Request: ActRequest{Handle: "h_repo", Verb: "status", Args: map[string]any{}},
	}
	if err := store.BeginAct(ctx, start); err != nil {
		t.Fatal(err)
	}
	start.Request.Verb = "build"
	if err := store.BeginAct(ctx, start); err == nil || !strings.Contains(err.Error(), "different payload") {
		t.Fatalf("colliding retry error = %v", err)
	}
}

func TestFinishEpisodeRecordsOutcomeAndRemainsIdempotent(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, Options{})
	episodeID, err := store.StartEpisode(ctx, "s_completed", testWorldBuild)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.FinishEpisode(ctx, episodeID, "success"); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishEpisode(ctx, episodeID, "success"); err != nil {
		t.Fatalf("idempotent finish: %v", err)
	}
	record, err := store.Episode(ctx, episodeID)
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != EpisodeCompleted || record.Outcome != "success" || record.EndedAt == nil {
		t.Fatalf("completed episode = %#v", record)
	}
}

func TestFinishEpisodeRejectsActiveAct(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, Options{})
	episodeID, err := store.StartEpisode(ctx, "s_active_act", testWorldBuild)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.BeginAct(ctx, StartedAct{
		EpisodeID: episodeID, ActID: "a_active", WorldBuild: testWorldBuild,
		Request: ActRequest{Handle: "h_repo", Verb: "status", Args: map[string]any{}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishEpisode(ctx, episodeID, "success"); err == nil || !strings.Contains(err.Error(), "active acts") {
		t.Fatalf("finish with active act error = %v", err)
	}
}

func TestLatestReturnsNewestEpisode(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, time.August, 18, 12, 0, 0, 0, time.UTC)
	store := openTestStore(t, Options{Now: func() time.Time {
		now = now.Add(time.Second)
		return now
	}})
	first, err := store.StartEpisode(ctx, "s_first", testWorldBuild)
	mustSucceed(t, err)
	second, err := store.StartEpisode(ctx, "s_second", testWorldBuild)
	mustSucceed(t, err)
	latest, err := store.Latest(ctx)
	mustSucceed(t, err)
	if latest.EpisodeID != second || latest.EpisodeID == first {
		t.Fatalf("latest episode = %q, want %q", latest.EpisodeID, second)
	}
}

func TestSecretsAreRedactedAndWALAcceptsConcurrentWriters(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, Options{Secrets: []string{"sk-live-secret"}})
	var mode string
	if err := store.db.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&mode); err != nil {
		t.Fatal(err)
	}
	if strings.ToLower(mode) != "wal" {
		t.Fatalf("journal mode = %q, want wal", mode)
	}

	const writers = 16
	var wait sync.WaitGroup
	errors := make(chan error, writers)
	for index := 0; index < writers; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			sessionID := fmt.Sprintf("s_concurrent_%02d", index)
			episodeID, err := store.StartEpisode(ctx, sessionID, testWorldBuild)
			if err != nil {
				errors <- err
				return
			}
			actID := fmt.Sprintf("a_concurrent_%02d", index)
			if err := store.BeginAct(ctx, StartedAct{
				EpisodeID: episodeID, ActID: actID, WorldBuild: testWorldBuild,
				Request: ActRequest{Handle: "h_repo", Verb: "read", Args: map[string]any{
					"password": "visible-password", "token": "prefix sk-live-secret suffix",
				}},
			}); err != nil {
				errors <- err
				return
			}
			if err := store.CompleteAct(ctx, CompletedAct{EpisodeID: episodeID, ActID: actID, Result: map[string]any{"status": "ok"}}); err != nil {
				errors <- err
			}
		}(index)
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		t.Fatal(err)
	}

	record, err := store.Episode(ctx, episodeIDForSession(t, store, "s_concurrent_00"))
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if strings.Contains(text, "visible-password") || strings.Contains(text, "sk-live-secret") {
		t.Fatalf("stored episode contains a secret: %s", text)
	}
	if !strings.Contains(text, "[REDACTED]") {
		t.Fatalf("stored episode does not show redaction: %s", text)
	}
}

func TestRecallReturnsPointersWithoutSummaryText(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, Options{})
	episodeID, err := store.StartEpisode(ctx, "s_recall", testWorldBuild)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.BeginAct(ctx, StartedAct{
		EpisodeID: episodeID, ActID: "a_auth_failure", WorldBuild: testWorldBuild,
		Request: ActRequest{Handle: "h_tests", Verb: "run", Args: map[string]any{}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CompleteAct(ctx, CompletedAct{
		EpisodeID: episodeID, ActID: "a_auth_failure",
		Result: map[string]any{"status": "fail", "error": map[string]any{"message": "authentication failure"}},
	}); err != nil {
		t.Fatal(err)
	}
	pointers, err := store.Recall(ctx, "s_recall", "authentication")
	if err != nil {
		t.Fatal(err)
	}
	want := []Pointer{{EpisodeID: episodeID, ActID: "a_auth_failure"}}
	if !reflect.DeepEqual(pointers, want) {
		t.Fatalf("Recall() = %#v, want pointers %#v", pointers, want)
	}
}

func openTestStore(t *testing.T, options Options) *Store {
	t.Helper()
	store, err := Open(filepath.Join(t.TempDir(), "episodes.db"), options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func mustSucceed(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func episodeIDForSession(t *testing.T, store *Store, sessionID string) string {
	t.Helper()
	var episodeID string
	if err := store.db.QueryRow("SELECT episode_id FROM episodes WHERE session_id = ?", sessionID).Scan(&episodeID); err != nil {
		t.Fatal(err)
	}
	return episodeID
}
