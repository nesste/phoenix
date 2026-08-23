package activate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/nesste/phoenix/internal/world"
)

func TestComputeBindsCapturedIntentToReachableCall(t *testing.T) {
	capture := 1
	empty := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	focus := json.RawMessage(`{"type":"object","additionalProperties":false,"required":["test"],"properties":{"test":{"type":"string"}}}`)
	definition := &world.Definition{
		V: 1, ID: "activation_test",
		Roots: []world.Root{{Name: "tests", Type: "tests", Label: "tests", Resource: world.Resource{Kind: "logical", Value: "tests"}}},
		HandleTypes: map[string]world.HandleType{"tests": {Verbs: map[string]world.Verb{
			"focus": {ArgsSchema: focus, ResultSchema: empty},
		}}},
		Activations: []world.ActivationRule{{
			ID: "named_test", Pattern: `(Test[A-Za-z0-9_]+)`,
			Suggestions: []world.Suggestion{{
				Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "tests"}, Verb: "focus", Args: map[string]world.Binding{"test": {IntentCapture: &capture}}},
				Why:  "run the named test", Score: 1,
			}},
		}},
	}
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, activationResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	entries := engine.Compute(session, roots[0].Ref, "Report whether TestBeaconCurrent passes")
	if len(entries) != 1 || entries[0].Call.Args["test"] != "TestBeaconCurrent" || entries[0].Provenance != "activation" {
		t.Fatalf("entries = %#v, want one bound activation", entries)
	}
	if hidden := engine.Compute(session, "h_0123456789abcdef", "TestBeaconCurrent"); len(hidden) != 0 {
		t.Fatalf("unreachable anchor exposed entries: %#v", hidden)
	}
}

func TestComputeUsesInspectFallbackWhenNoRuleMatches(t *testing.T) {
	engine, session, root := inspectFallbackWorld(t, false)
	entries := engine.Compute(session, root, "Make leaf a belong to the root declared by the manifest")
	if len(entries) != 1 || entries[0].Call.Verb != "status" || entries[0].Provenance != "activation" {
		t.Fatalf("entries = %#v, want inspect fallback status call", entries)
	}
}

func TestComputeAuthoredMissDoesNotUseInspectFallback(t *testing.T) {
	engine, session, root := inspectFallbackWorld(t, true)
	entries := engine.Compute(session, root, "Deploy the Parcel build to production and report the release identifier")
	if len(entries) != 0 {
		t.Fatalf("entries = %#v, want authored miss with no fallback", entries)
	}
}

func TestComputeWithoutInspectFallbackStaysEmpty(t *testing.T) {
	empty := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	definition := &world.Definition{
		V: 1, ID: "activation_test",
		Roots: []world.Root{{Name: "repo", Type: "repo", Label: "repo", Resource: world.Resource{Kind: "logical", Value: "repo"}}},
		HandleTypes: map[string]world.HandleType{"repo": {Verbs: map[string]world.Verb{
			"status": {ArgsSchema: empty, ResultSchema: empty},
		}}},
		Activations: []world.ActivationRule{{
			ID: "named_test", Pattern: `(Test[A-Za-z0-9_]+)`,
			Suggestions: []world.Suggestion{{
				Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "repo"}, Verb: "status"},
				Why:  "unused specific rule", Score: 1,
			}},
		}},
	}
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, activationResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	entries := engine.Compute(session, roots[0].Ref, "Make leaf a belong to the root declared by the manifest")
	if len(entries) != 0 {
		t.Fatalf("entries = %#v, want empty without inspect fallback", entries)
	}
}

func TestNewRejectsInspectFallbackWithoutSuggestions(t *testing.T) {
	empty := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	definition := &world.Definition{
		V: 1, ID: "activation_test",
		Roots: []world.Root{{Name: "repo", Type: "repo", Label: "repo", Resource: world.Resource{Kind: "logical", Value: "repo"}}},
		HandleTypes: map[string]world.HandleType{"repo": {Verbs: map[string]world.Verb{
			"status": {ArgsSchema: empty, ResultSchema: empty},
		}}},
		Activations: []world.ActivationRule{{
			ID: inspectFallbackID, Pattern: `^$a`, Suggestions: []world.Suggestion{},
		}},
	}
	if _, err := New(definition); err == nil {
		t.Fatal("expected inspect fallback without suggestions to fail")
	}
}

func inspectFallbackWorld(t *testing.T, includeAbsence bool) (*Engine, Topology, string) {
	t.Helper()
	empty := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	activations := []world.ActivationRule{
		{
			ID: "named_test", Pattern: `(Test[A-Za-z0-9_]+)`,
			Suggestions: []world.Suggestion{{
				Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "repo"}, Verb: "status"},
				Why:  "unused specific rule", Score: 1,
			}},
		},
		{
			ID: inspectFallbackID, Pattern: `^$a`,
			Suggestions: []world.Suggestion{{
				Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "repo"}, Verb: "status"},
				Why:  "inspect reachable repository state", Score: 0,
			}},
		},
	}
	if includeAbsence {
		activations = append([]world.ActivationRule{{
			ID: "absent_deploy_or_release", Pattern: `(?i)\b(deploy|deployment)\b|\brelease identifier\b`,
			Suggestions: []world.Suggestion{},
		}}, activations...)
	}
	definition := &world.Definition{
		V: 1, ID: "activation_test",
		Roots: []world.Root{{Name: "repo", Type: "repo", Label: "repo", Resource: world.Resource{Kind: "logical", Value: "repo"}}},
		HandleTypes: map[string]world.HandleType{"repo": {Verbs: map[string]world.Verb{
			"status": {ArgsSchema: empty, ResultSchema: empty},
		}}},
		Activations: activations,
	}
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, activationResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	return engine, session, roots[0].Ref
}

type activationResolver struct{}

func (activationResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{"revision": "test"}, nil
}
