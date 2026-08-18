package devrepo

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

func TestRepoReadEditAndTraversalRefusal(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "note.txt")
	if err := os.WriteFile(path, []byte("before\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := newTestExecutor(t, Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git"})
	resource := world.Resource{Kind: "path", Value: root}

	read := executor.Execute(context.Background(), verb.Request{
		HandleType: "repo", Verb: "read", Resource: resource, Args: map[string]any{"path": "note.txt"},
	})
	if read.Status != verb.StatusOK || read.Value.(map[string]any)["content"] != "before\n" {
		t.Fatalf("read result = %#v", read)
	}
	edit := executor.Execute(context.Background(), verb.Request{
		HandleType: "repo", Verb: "edit", Resource: resource,
		Args: map[string]any{"path": "note.txt", "old": "before", "new": "after"},
	})
	if edit.Status != verb.StatusOK || edit.Value.(map[string]any)["changed"] != true {
		t.Fatalf("edit result = %#v", edit)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(contents), "after\n"; got != want {
		t.Fatalf("edited file = %q, want %q", got, want)
	}

	outside := executor.Execute(context.Background(), verb.Request{
		HandleType: "repo", Verb: "read", Resource: resource, Args: map[string]any{"path": "../outside.txt"},
	})
	if outside.Status != verb.StatusFail || outside.Error.Code != "path_outside_repo" {
		t.Fatalf("traversal result = %#v", outside)
	}
}

func TestFindUsesOnlyReachableTopologyAndRepositoryFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "catalog.txt"), []byte("TestRelayHandshake\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := newTestExecutor(t, Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git", MaxSearchMatches: 10})
	finder := staticFinder{matches: []world.ReachableMatch{{Handle: "h_1234567890123456", Type: "tests", Label: "suite", Verb: "focus"}}}

	result := executor.Execute(context.Background(), verb.Request{
		HandleType: "repo", Verb: "find", Resource: world.Resource{Kind: "path", Value: root},
		Args: map[string]any{"query": "relay"}, Reachable: finder,
	})
	if result.Status != verb.StatusOK {
		t.Fatalf("find result = %#v", result)
	}
	matches := result.Value.(map[string]any)["matches"].([]any)
	joined := strings.ToLower(strings.Join(anyStrings(matches), "\n"))
	if !strings.Contains(joined, "catalog.txt") {
		t.Fatalf("find matches do not include repository result: %s", joined)
	}
	if finder.lastQuery != "" {
		t.Fatal("static finder must be immutable in this test")
	}
}

type staticFinder struct {
	matches   []world.ReachableMatch
	lastQuery string
}

func (finder staticFinder) FindReachable(string) []world.ReachableMatch {
	return append([]world.ReachableMatch(nil), finder.matches...)
}

func anyStrings(values []any) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.(string))
	}
	return result
}
