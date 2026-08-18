package devrepo

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/world"
)

func TestProductionWorldMatchesRegisteredVerbsAndAuthoredRules(t *testing.T) {
	definition, err := world.Load(
		filepath.Join("..", "..", "spec", "world.schema.json"),
		filepath.Join("..", "..", "worlds", "dev-repo", "world.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := rootNames(definition.Roots), []string{"episodes", "git", "repo", "tests"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("world roots = %#v, want %#v", got, want)
	}
	expected := authoredWorld(t)
	for handleType, descriptor := range expected.HandleTypes {
		actual := definition.HandleTypes[handleType]
		if len(actual.Verbs) != len(descriptor.Verbs) {
			t.Fatalf("%s verb count = %d, want %d", handleType, len(actual.Verbs), len(descriptor.Verbs))
		}
		for name, want := range descriptor.Verbs {
			got, exists := actual.Verbs[name]
			if !exists || !sameJSON(got.ArgsSchema, want.ArgsSchema) || !sameJSON(got.ResultSchema, want.ResultSchema) || !sameRules(got.Refusals, want.Refusals) {
				t.Fatalf("world verb %s.%s differs from its registered definition", handleType, name)
			}
		}
	}
	if !sameJSON(definition.Transitions, AuthoredTransitions()) {
		t.Fatal("world transitions differ from AuthoredTransitions")
	}
	if !sameJSON(definition.Activations, AuthoredActivations()) {
		t.Fatal("world activations differ from AuthoredActivations")
	}
}

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

func sameRules(left, right []world.RefusalRule) bool {
	if len(left) == 0 && len(right) == 0 {
		return true
	}
	return sameJSON(left, right)
}

func rootNames(roots []world.Root) []string {
	names := make([]string, 0, len(roots))
	for _, root := range roots {
		names = append(names, root.Name)
	}
	sort.Strings(names)
	return names
}

func sameJSON(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	var leftValue, rightValue any
	if json.Unmarshal(leftJSON, &leftValue) != nil || json.Unmarshal(rightJSON, &rightValue) != nil {
		return false
	}
	return reflect.DeepEqual(leftValue, rightValue)
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
