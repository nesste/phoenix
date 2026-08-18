package frontier

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/nesste/phoenix/internal/world"
)

func TestEngineBindsReachableCallsAndRanksDeterministically(t *testing.T) {
	definition := testDefinition()
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, frontierResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	root := rootRefs(roots)
	observation := Observation{
		HandleType: "repo", Handle: root["repo"], Verb: "inspect", Status: "ok",
		Result: map[string]any{"path": "internal/frontier/rules.go", "relevant": true},
		State:  &world.LiveState{Digest: testStateDigest, Value: map[string]any{"revision": testStateDigest}},
	}

	first := engine.Compute(session, observation)
	second := engine.Compute(session, observation)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("identical observations produced different frontiers:\n%#v\n%#v", first, second)
	}
	if got, want := len(first), 3; got != want {
		t.Fatalf("frontier length = %d, want engine cap %d: %#v", got, want, first)
	}
	if first[0].Call.Handle != root["tests"] || first[0].Call.Verb != "focus" {
		t.Fatalf("first call = %#v, want bound tests.focus", first[0].Call)
	}
	if got := first[0].Call.Args["test"]; got != "internal/frontier/rules.go" {
		t.Fatalf("bound test = %#v, want result path", got)
	}
	if first[0].Call.State == nil || *first[0].Call.State != testStateDigest {
		t.Fatalf("state = %#v, want live digest", first[0].Call.State)
	}
	for _, entry := range first {
		if entry.Provenance != "authored" {
			t.Fatalf("provenance = %q, want authored", entry.Provenance)
		}
	}
}

func TestEngineReturnsEmptyWithoutMatchingUsefulRules(t *testing.T) {
	definition := testDefinition()
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, frontierResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	entries := engine.Compute(session, Observation{
		HandleType: "repo", Handle: roots[0].Ref, Verb: "inspect", Status: "ok",
		Result: map[string]any{"path": "README.md", "relevant": false},
		State:  &world.LiveState{Digest: testStateDigest, Value: map[string]any{"revision": testStateDigest}},
	})
	if entries == nil || len(entries) != 0 {
		t.Fatalf("frontier = %#v, want explicit empty slice", entries)
	}
}

func TestEngineOmitsCandidatesWhoseHandlesAreNoLongerReachable(t *testing.T) {
	definition := testDefinition()
	repoVerb := definition.HandleTypes["repo"].Verbs["inspect"]
	repoVerb.MayRevoke = true
	descriptor := definition.HandleTypes["repo"]
	descriptor.Verbs["inspect"] = repoVerb
	definition.HandleTypes["repo"] = descriptor
	engine, err := New(definition)
	if err != nil {
		t.Fatal(err)
	}
	session, roots, err := world.NewGraph(definition, frontierResolver{}).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	refs := rootRefs(roots)
	if _, err := session.ApplyDelta(refs["repo"], "inspect", world.Mutation{Revokes: []string{refs["tests"]}}); err != nil {
		t.Fatal(err)
	}

	entries := engine.Compute(session, Observation{
		HandleType: "repo", Handle: refs["repo"], Verb: "inspect", Status: "ok",
		Result: map[string]any{"path": "README.md", "relevant": true},
		State:  &world.LiveState{Digest: testStateDigest, Value: map[string]any{"revision": testStateDigest}},
	})
	for _, entry := range entries {
		if entry.Call.Handle == refs["tests"] {
			t.Fatalf("frontier retained revoked handle: %#v", entries)
		}
	}
	if got, want := len(entries), 2; got != want {
		t.Fatalf("frontier length after revoke = %d, want %d useful calls: %#v", got, want, entries)
	}
}

func TestNewRejectsInvalidAuthoredRules(t *testing.T) {
	tests := []struct {
		name string
		edit func(*world.Definition)
		want string
	}{
		{
			name: "unbound required argument",
			edit: func(definition *world.Definition) {
				definition.Transitions[0].Suggestions[0].Call.Args = map[string]world.Binding{}
			},
			want: `does not bind required argument "test"`,
		},
		{
			name: "unreachable root",
			edit: func(definition *world.Definition) {
				definition.Transitions[0].Suggestions[0].Call.Handle.Name = "hidden"
			},
			want: `root "hidden" does not exist`,
		},
		{
			name: "why line",
			edit: func(definition *world.Definition) {
				definition.Transitions[0].Suggestions[0].Why = strings.Repeat("x", 81)
			},
			want: "why line exceeds 80 characters",
		},
		{
			name: "frontier cap",
			edit: func(definition *world.Definition) {
				suggestion := definition.Transitions[0].Suggestions[0]
				definition.Transitions[0].Suggestions = []world.Suggestion{suggestion, suggestion, suggestion, suggestion}
			},
			want: "more than 3 suggestions",
		},
		{
			name: "opaque call string cannot replace a structured call",
			edit: func(definition *world.Definition) {
				definition.Transitions[0].Suggestions[0].Call.Verb = "tests.focus()"
			},
			want: "invalid verb",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			definition := testDefinition()
			test.edit(definition)
			_, err := New(definition)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("New() error = %v, want it to contain %q", err, test.want)
			}
		})
	}
}

const testStateDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func testDefinition() *world.Definition {
	emptyArgs := json.RawMessage(`{"type":"object","additionalProperties":false}`)
	focusArgs := json.RawMessage(`{
		"type":"object","additionalProperties":false,"required":["test"],
		"properties":{"test":{"type":"string","minLength":1}}
	}`)
	result := json.RawMessage(`{
		"type":"object","additionalProperties":false,"required":["path","relevant"],
		"properties":{"path":{"type":"string"},"relevant":{"type":"boolean"}}
	}`)
	emptyResult := json.RawMessage(`{"type":"object"}`)
	path := "/path"
	revision := "/revision"
	suggestions := []world.Suggestion{
		{
			Call: world.CallTemplate{
				Handle: world.HandleSelector{Source: "root", Name: "tests"}, Verb: "focus",
				Args:  map[string]world.Binding{"test": {ResultPointer: &path}},
				State: &world.Binding{StatePointer: &revision},
			},
			Why: "run the affected test", Score: 1,
		},
		{Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "git"}, Verb: "diff", Args: map[string]world.Binding{}}, Why: "inspect the resulting diff", Score: .8},
		{Call: world.CallTemplate{Handle: world.HandleSelector{Source: "root", Name: "tests"}, Verb: "run", Args: map[string]world.Binding{}}, Why: "verify the complete suite", Score: .7},
	}
	return &world.Definition{
		V: 1, ID: "frontier_test",
		Roots: []world.Root{
			{Name: "repo", Type: "repo", Label: "repository", Resource: world.Resource{Kind: "logical", Value: "repo"}},
			{Name: "tests", Type: "tests", Label: "tests", Resource: world.Resource{Kind: "logical", Value: "tests"}},
			{Name: "git", Type: "git", Label: "git", Resource: world.Resource{Kind: "logical", Value: "git"}},
		},
		HandleTypes: map[string]world.HandleType{
			"repo": {Verbs: map[string]world.Verb{
				"inspect": {ArgsSchema: emptyArgs, ResultSchema: result},
				"status":  {ArgsSchema: emptyArgs, ResultSchema: emptyResult},
			}},
			"tests": {Verbs: map[string]world.Verb{
				"focus": {ArgsSchema: focusArgs, ResultSchema: emptyResult},
				"run":   {ArgsSchema: emptyArgs, ResultSchema: emptyResult},
			}},
			"git": {Verbs: map[string]world.Verb{
				"diff": {ArgsSchema: emptyArgs, ResultSchema: emptyResult},
			}},
		},
		Transitions: []world.Transition{
			{
				ID:          "relevant_inspection",
				Match:       world.Match{HandleType: "repo", Verb: "inspect", Status: "ok", ResultWhen: json.RawMessage(`{"properties":{"relevant":{"const":true}},"required":["relevant"]}`)},
				Suggestions: suggestions,
			},
			{
				ID:          "relevant_status",
				Match:       world.Match{HandleType: "repo", Verb: "inspect", Status: "ok", ResultWhen: json.RawMessage(`{"properties":{"relevant":{"const":true}},"required":["relevant"]}`)},
				Suggestions: []world.Suggestion{{Call: world.CallTemplate{Handle: world.HandleSelector{Source: "self"}, Verb: "status", Args: map[string]world.Binding{}}, Why: "refresh repository state", Score: .1}},
			},
		},
	}
}

func rootRefs(roots []world.RootHandle) map[string]string {
	refs := make(map[string]string, len(roots))
	for _, root := range roots {
		refs[root.Name] = root.Ref
	}
	return refs
}

type frontierResolver struct{}

func (frontierResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{"revision": testStateDigest}, nil
}
