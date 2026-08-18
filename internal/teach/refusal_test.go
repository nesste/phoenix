package teach

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nesste/phoenix/internal/world"
)

func TestEngineShapesRefusalWithReachableBoundAlternative(t *testing.T) {
	definition := teachingDefinition()
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, teachingResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := teachingRootRefs(roots)

	refusal, matched, err := engine.Evaluate(session, Observation{
		HandleType: "tests", Handle: refs["tests"], Verb: "focus", Args: map[string]any{},
		State: &world.LiveState{Digest: teachingDigest, Value: map[string]any{"ready": true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("missing test did not match its authored refusal")
	}
	if refusal.What == "" || refusal.Why == "" || refusal.Instead == nil {
		t.Fatalf("refusal is incomplete: %#v", refusal)
	}
	if refusal.Instead.Handle != refs["tests"] || refusal.Instead.Verb != "list" || refusal.Instead.Args == nil {
		t.Fatalf("alternative = %#v, want reachable tests.list", refusal.Instead)
	}
}

func TestEngineShapesExplicitNoAlternativeFromFailure(t *testing.T) {
	definition := teachingDefinition()
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, teachingResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := teachingRootRefs(roots)

	refusal, matched, err := engine.Evaluate(session, Observation{
		HandleType: "git", Handle: refs["git"], Verb: "commit", Args: map[string]any{"message": "checkpoint"},
		State:   &world.LiveState{Digest: teachingDigest, Value: map[string]any{"clean": true}},
		Failure: map[string]any{"code": "nothing_to_commit", "message": "no changes are available to commit"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !matched || refusal.Instead != nil || refusal.NoAlternative == "" {
		t.Fatalf("refusal = %#v, matched = %v; want explicit no-alternative refusal", refusal, matched)
	}
}

func TestEngineDoesNotRefuseWithoutMatchingState(t *testing.T) {
	definition := teachingDefinition()
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, teachingResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := teachingRootRefs(roots)

	refusal, matched, err := engine.Evaluate(session, Observation{
		HandleType: "tests", Handle: refs["tests"], Verb: "focus", Args: map[string]any{"test": "TestEngine"},
		State: &world.LiveState{Digest: teachingDigest, Value: map[string]any{"ready": true}},
	})
	if err != nil || matched || refusal != (Refusal{}) {
		t.Fatalf("Evaluate() = (%#v, %v, %v), want no refusal", refusal, matched, err)
	}
}

func TestNewRejectsIncompleteOrUnreachableRefusals(t *testing.T) {
	tests := []struct {
		name string
		edit func(*world.Definition)
		want string
	}{
		{
			name: "missing what",
			edit: func(definition *world.Definition) {
				definition.HandleTypes["tests"].Verbs["focus"].Refusals[0].What = ""
			},
			want: "what is required",
		},
		{
			name: "missing why",
			edit: func(definition *world.Definition) {
				definition.HandleTypes["tests"].Verbs["focus"].Refusals[0].Why = ""
			},
			want: "why is required",
		},
		{
			name: "missing no alternative reason",
			edit: func(definition *world.Definition) {
				definition.HandleTypes["git"].Verbs["commit"].Refusals[0].NoAlternative = ""
			},
			want: "no_alternative is required",
		},
		{
			name: "unreachable alternative root",
			edit: func(definition *world.Definition) {
				definition.HandleTypes["tests"].Verbs["focus"].Refusals[0].Instead.Handle = world.HandleSelector{Source: "root", Name: "hidden"}
			},
			want: `selected root "hidden" does not exist`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			definition := teachingDefinition()
			test.edit(definition)
			_, err := New(definition)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("New() error = %v, want it to contain %q", err, test.want)
			}
		})
	}
}

const teachingDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func teachingDefinition() *world.Definition {
	emptyArgs := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	focusArgs := json.RawMessage(`{"type":"object","additionalProperties":false,"required":["test"],"properties":{"test":{"type":"string"}}}`)
	commitArgs := json.RawMessage(`{"type":"object","additionalProperties":false,"required":["message"],"properties":{"message":{"type":"string"}}}`)
	result := json.RawMessage(`{"type":"object"}`)
	list := world.CallTemplate{Handle: world.HandleSelector{Source: "self"}, Verb: "list", Args: map[string]world.Binding{}}
	return &world.Definition{
		V: 1, ID: "teaching_test",
		Roots: []world.Root{
			{Name: "tests", Type: "tests", Label: "tests", Resource: world.Resource{Kind: "logical", Value: "tests"}},
			{Name: "git", Type: "git", Label: "git", Resource: world.Resource{Kind: "logical", Value: "git"}},
		},
		HandleTypes: map[string]world.HandleType{
			"tests": {Verbs: map[string]world.Verb{
				"focus": {ArgsSchema: focusArgs, ResultSchema: result, Refusals: []world.RefusalRule{{
					ID: "missing_test", When: json.RawMessage(`{"properties":{"args":{"not":{"required":["test"]}}},"required":["args"]}`),
					What: "tests.focus requires a test name", Why: "the requested test was not bound", Instead: &list,
				}}},
				"list": {ArgsSchema: emptyArgs, ResultSchema: result},
			}},
			"git": {Verbs: map[string]world.Verb{
				"commit": {ArgsSchema: commitArgs, ResultSchema: result, Refusals: []world.RefusalRule{{
					ID: "nothing_to_commit", When: json.RawMessage(`{"properties":{"failure":{"properties":{"code":{"const":"nothing_to_commit"}},"required":["code"]}},"required":["failure"]}`),
					What: "git.commit cannot run", Why: "no changes are available to commit", Instead: nil,
					NoAlternative: "make a reachable edit before committing",
				}}},
			}},
		},
		Transitions: []world.Transition{},
	}
}

func teachingRootRefs(roots []world.RootHandle) map[string]string {
	refs := make(map[string]string, len(roots))
	for _, root := range roots {
		refs[root.Name] = root.Ref
	}
	return refs
}

type teachingResolver struct{}

func (teachingResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{"ready": true}, nil
}
