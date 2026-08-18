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

type activationResolver struct{}

func (activationResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{"revision": "test"}, nil
}
