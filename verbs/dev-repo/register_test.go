package devrepo

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

func TestDefinitionsContainOnlyTheStarterVerbSet(t *testing.T) {
	definitions, err := Definitions(Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git"})
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, 0, len(definitions))
	for _, definition := range definitions {
		got = append(got, definition.HandleType+"."+definition.Name)
	}
	sort.Strings(got)
	want := []string{
		"episodes.recall", "git.commit", "git.diff", "git.status", "repo.build", "repo.edit",
		"repo.find", "repo.read", "repo.status", "tests.focus", "tests.list", "tests.run",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("starter verbs = %#v, want %#v", got, want)
	}
}

func TestAuthoredTransitionsCompileAndBindToDevRepoRoots(t *testing.T) {
	definition := authoredWorld(t)
	engine, err := frontier.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, authoredResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := make(map[string]string, len(roots))
	for _, root := range roots {
		refs[root.Name] = root.Ref
	}

	entries := engine.Compute(session, frontier.Observation{
		HandleType: "repo", Handle: refs["repo"], Verb: "edit", Status: "ok",
		Result: map[string]any{"path": "main.go", "changed": true, "digest": "sha256:value"},
	})
	if got, want := len(entries), 2; got != want {
		t.Fatalf("edit frontier length = %d, want %d: %#v", got, want, entries)
	}
	if entries[0].Call.Handle != refs["tests"] || entries[0].Call.Verb != "run" {
		t.Fatalf("first edit suggestion = %#v, want reachable tests.run", entries[0])
	}
	if entries[1].Call.Handle != refs["git"] || entries[1].Call.Verb != "diff" {
		t.Fatalf("second edit suggestion = %#v, want reachable git.diff", entries[1])
	}
}

func TestAuthoredRefusalsTeachReachableOrExplicitAlternatives(t *testing.T) {
	definition := authoredWorld(t)
	engine, err := teach.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, authoredResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := make(map[string]string, len(roots))
	for _, root := range roots {
		refs[root.Name] = root.Ref
	}

	missing, matched, err := engine.Evaluate(session, teach.Observation{
		HandleType: "tests", Handle: refs["tests"], Verb: "focus", Args: map[string]any{},
	})
	if err != nil || !matched || missing.Instead == nil || missing.Instead.Verb != "list" {
		t.Fatalf("missing-test refusal = (%#v, %v, %v), want reachable tests.list", missing, matched, err)
	}
	conflict, matched, err := engine.Evaluate(session, teach.Observation{
		HandleType: "repo", Handle: refs["repo"], Verb: "edit", Args: map[string]any{},
		Failure: map[string]any{"code": "edit_conflict", "message": "conflict"},
	})
	if err != nil || !matched || conflict.Instead != nil || conflict.NoAlternative == "" {
		t.Fatalf("edit-conflict refusal = (%#v, %v, %v), want explicit reason", conflict, matched, err)
	}
}

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

func TestCommandVerbsUseFixedExecutablesAndArgumentArrays(t *testing.T) {
	runner := &fakeRunner{responses: []verb.CommandResult{
		{ExitCode: 0, Stdout: []byte("build ok")},
		{ExitCode: 0, Stdout: []byte("{\"Action\":\"pass\",\"Test\":\"TestX;touch_pwned\"}\n")},
		{ExitCode: 0},
		{ExitCode: 0, Stdout: []byte("committed")},
	}}
	executor := newTestExecutor(t, Config{Runner: runner, GoExecutable: "go_fixed", GitExecutable: "git_fixed"})
	resource := world.Resource{Kind: "path", Value: t.TempDir()}

	build := executor.Execute(context.Background(), verb.Request{HandleType: "repo", Verb: "build", Resource: resource, Args: map[string]any{}})
	if build.Status != verb.StatusOK {
		t.Fatalf("build result = %#v", build)
	}
	focus := executor.Execute(context.Background(), verb.Request{
		HandleType: "tests", Verb: "focus", Resource: resource, Args: map[string]any{"test": "TestX;touch_pwned"},
	})
	if focus.Status != verb.StatusOK || focus.Value.(map[string]any)["passed"] != true {
		t.Fatalf("focus result = %#v", focus)
	}
	commit := executor.Execute(context.Background(), verb.Request{
		HandleType: "git", Verb: "commit", Resource: resource, Args: map[string]any{"message": "safe; touch pwned"},
	})
	if commit.Status != verb.StatusOK {
		t.Fatalf("commit result = %#v", commit)
	}

	commands := runner.Commands()
	if commands[0].Executable != "go_fixed" || !reflect.DeepEqual(commands[0].Args, []string{"test", "-run=^$", "./..."}) {
		t.Fatalf("build command = %#v", commands[0])
	}
	if commands[1].Executable != "go_fixed" || commands[1].Args[len(commands[1].Args)-1] != "^TestX;touch_pwned$" {
		t.Fatalf("focus command did not preserve one typed argument: %#v", commands[1])
	}
	if commands[2].Executable != "git_fixed" || !reflect.DeepEqual(commands[2].Args, []string{"add", "-A", "--", "."}) {
		t.Fatalf("git add command = %#v", commands[2])
	}
	if commands[3].Args[3] != "safe; touch pwned" {
		t.Fatalf("commit message was not one argument: %#v", commands[3].Args)
	}
}

func TestGitCommitReportsNothingToCommitAsAConstraint(t *testing.T) {
	runner := &fakeRunner{responses: []verb.CommandResult{
		{ExitCode: 0},
		{ExitCode: 1, Stdout: []byte("nothing to commit, working tree clean")},
	}}
	executor := newTestExecutor(t, Config{Runner: runner, GoExecutable: "go", GitExecutable: "git"})
	result := executor.Execute(context.Background(), verb.Request{
		HandleType: "git", Verb: "commit", Resource: world.Resource{Kind: "path", Value: t.TempDir()},
		Args: map[string]any{"message": "checkpoint"},
	})
	if result.Status != verb.StatusFail || result.Error == nil || result.Error.Code != "nothing_to_commit" {
		t.Fatalf("commit result = %#v, want shapeable nothing_to_commit constraint", result)
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

func TestRecallFailsClosedUntilEpisodeStoreExists(t *testing.T) {
	executor := newTestExecutor(t, Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git"})
	result := executor.Execute(context.Background(), verb.Request{
		HandleType: "episodes", Verb: "recall", SessionID: "s_0123456789abcdef", Args: map[string]any{"query": "failure"},
	})
	if result.Status != verb.StatusFail || result.Error.Code != "recall_unavailable" {
		t.Fatalf("recall without store = %#v", result)
	}
}

func newTestExecutor(t *testing.T, config Config) *verb.Executor {
	t.Helper()
	registry := verb.NewRegistry()
	if err := Register(registry, config); err != nil {
		t.Fatal(err)
	}
	return verb.NewExecutor(registry, verb.Options{Timeout: time.Second, MaxResultBytes: 64 * 1024})
}

type fakeRunner struct {
	commands  []verb.Command
	responses []verb.CommandResult
}

func (runner *fakeRunner) Run(_ context.Context, command verb.Command) (verb.CommandResult, error) {
	runner.commands = append(runner.commands, command)
	if len(runner.responses) == 0 {
		return verb.CommandResult{}, nil
	}
	response := runner.responses[0]
	runner.responses = runner.responses[1:]
	return response, nil
}

func (runner *fakeRunner) Commands() []verb.Command {
	return append([]verb.Command(nil), runner.commands...)
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

func authoredWorld(t *testing.T) *world.Definition {
	t.Helper()
	definitions, err := Definitions(Config{Runner: &fakeRunner{}, GoExecutable: "go", GitExecutable: "git"})
	if err != nil {
		t.Fatal(err)
	}
	handleTypes := make(map[string]world.HandleType)
	for _, definition := range definitions {
		descriptor := handleTypes[definition.HandleType]
		if descriptor.Verbs == nil {
			descriptor.Verbs = make(map[string]world.Verb)
		}
		descriptor.Verbs[definition.Name] = world.Verb{
			ArgsSchema: definition.ArgsSchema, ResultSchema: definition.ResultSchema,
		}
		handleTypes[definition.HandleType] = descriptor
	}
	for _, authored := range AuthoredRefusals() {
		descriptor := handleTypes[authored.HandleType]
		verb := descriptor.Verbs[authored.Verb]
		verb.Refusals = append([]world.RefusalRule(nil), authored.Rules...)
		descriptor.Verbs[authored.Verb] = verb
		handleTypes[authored.HandleType] = descriptor
	}
	roots := make([]world.Root, 0, len(handleTypes))
	for _, name := range []string{"repo", "tests", "git", "episodes"} {
		roots = append(roots, world.Root{
			Name: name, Type: name, Label: name,
			Resource: world.Resource{Kind: "logical", Value: name},
		})
	}
	return &world.Definition{
		V: 1, ID: "dev_repo", Roots: roots, HandleTypes: handleTypes,
		Transitions: AuthoredTransitions(),
	}
}

type authoredResolver struct{}

func (authoredResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{}, nil
}
