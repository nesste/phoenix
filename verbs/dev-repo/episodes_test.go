package devrepo

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/verb"
)

func TestRecallFailsClosedUntilEpisodeStoreExists(t *testing.T) {
	executor := newTestExecutor(t, Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git"})
	result := executor.Execute(context.Background(), verb.Request{
		HandleType: "episodes", Verb: "recall", SessionID: "s_0123456789abcdef", Args: map[string]any{"query": "failure"},
	})
	if result.Status != verb.StatusFail || result.Error.Code != "recall_unavailable" {
		t.Fatalf("recall without store = %#v", result)
	}
}

func TestRecallReturnsStoredEpisodePointers(t *testing.T) {
	ctx := context.Background()
	store, err := episode.Open(filepath.Join(t.TempDir(), "episodes.db"), episode.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	const sessionID = "s_recall_store"
	const worldBuild = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	episodeID, err := store.StartEpisode(ctx, sessionID, worldBuild)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.BeginAct(ctx, episode.StartedAct{
		EpisodeID: episodeID, ActID: "a_prior_failure", WorldBuild: worldBuild,
		Request: episode.ActRequest{Handle: "h_tests", Verb: "run", Args: map[string]any{}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := store.CompleteAct(ctx, episode.CompletedAct{
		EpisodeID: episodeID, ActID: "a_prior_failure",
		Result: map[string]any{"status": "fail", "message": "authentication failure"},
	}); err != nil {
		t.Fatal(err)
	}

	executor := newTestExecutor(t, Config{
		Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git", Episodes: store,
	})
	result := executor.Execute(ctx, verb.Request{
		HandleType: "episodes", Verb: "recall", SessionID: sessionID,
		Args: map[string]any{"query": "authentication"},
	})
	if result.Status != verb.StatusOK {
		t.Fatalf("recall result = %#v", result)
	}
	pointers := result.Value.(map[string]any)["episodes"].([]any)
	if len(pointers) != 1 {
		t.Fatalf("recall pointers = %#v, want one pointer", pointers)
	}
	pointer := pointers[0].(map[string]any)
	if pointer["episode_id"] != episodeID || pointer["act_id"] != "a_prior_failure" {
		t.Fatalf("recall pointer = %#v", pointer)
	}
}
