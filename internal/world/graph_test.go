package world

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestSessionHandlesAreOpaqueScopedAndExplicitlyMutated(t *testing.T) {
	definition := loadBaseWorld(t)
	resolver := &sequenceResolver{values: []any{map[string]any{"revision": 1}}}
	graph := NewGraph(definition, resolver)
	first, firstRoots, err := graph.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	second, secondRoots, err := graph.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(firstRoots), 1; got != want {
		t.Fatalf("root count = %d, want %d", got, want)
	}
	firstRoot := firstRoots[0]
	if !strings.HasPrefix(firstRoot.Ref, "h_") || strings.Contains(firstRoot.Ref, "repo") {
		t.Fatalf("root handle %q is not opaque", firstRoot.Ref)
	}
	if firstRoot.Ref == secondRoots[0].Ref {
		t.Fatal("two sessions received the same root handle")
	}

	guessed := first.Prepare(context.Background(), "h_0123456789abcdef", "inspect", "")
	crossSession := second.Prepare(context.Background(), firstRoot.Ref, "inspect", "")
	unknownVerb := first.Prepare(context.Background(), firstRoot.Ref, "not_attached", "")
	for name, result := range map[string]Access{"guessed": guessed, "cross-session": crossSession, "unknown verb": unknownVerb} {
		if result.Status != AccessAbsent {
			t.Fatalf("%s status = %q, want %q", name, result.Status, AccessAbsent)
		}
		if !reflect.DeepEqual(result.Problem, guessed.Problem) {
			t.Fatalf("%s problem = %#v, want generic %#v", name, result.Problem, guessed.Problem)
		}
	}
}

func TestSessionReachabilityChangesOnlyThroughExplicitDelta(t *testing.T) {
	definition := loadBaseWorld(t)
	resolver := &sequenceResolver{values: []any{map[string]any{"revision": 1}}}
	session, roots, err := NewGraph(definition, resolver).StartSession()
	if err != nil {
		t.Fatal(err)
	}
	root := roots[0]
	delta, err := session.ApplyDelta(root.Ref, "inspect", Mutation{Grants: []GrantRequest{{
		Type: "detail", Label: "inspection detail", Resource: Resource{Kind: "logical", Value: "detail"},
	}}})
	if err != nil {
		t.Fatalf("ApplyDelta(grant): %v", err)
	}
	if got, want := len(delta.Grant), 1; got != want {
		t.Fatalf("grant count = %d, want %d", got, want)
	}
	granted := delta.Grant[0].Ref
	if result := session.Prepare(context.Background(), granted, "show", ""); result.Status != AccessReady {
		t.Fatalf("granted handle status = %q, want %q; problem: %#v", result.Status, AccessReady, result.Problem)
	}

	revoked, err := session.ApplyDelta(root.Ref, "inspect", Mutation{Revokes: []string{granted}})
	if err != nil {
		t.Fatalf("ApplyDelta(revoke): %v", err)
	}
	if !reflect.DeepEqual(revoked.Revoke, []string{granted}) {
		t.Fatalf("revoked refs = %#v, want %#v", revoked.Revoke, []string{granted})
	}
	guessed := session.Prepare(context.Background(), "h_0123456789abcdef", "show", "")
	if result := session.Prepare(context.Background(), granted, "show", ""); result.Status != AccessAbsent || !reflect.DeepEqual(result.Problem, guessed.Problem) {
		t.Fatalf("revoked handle leaked its prior existence: %#v", result)
	}
}

func TestPrepareResolvesLiveStateAndRejectsStalePrecondition(t *testing.T) {
	definition := loadBaseWorld(t)
	resolver := &sequenceResolver{values: []any{
		map[string]any{"revision": 1},
		map[string]any{"revision": 2},
		map[string]any{"revision": 2},
	}}
	graph := NewGraph(definition, resolver)
	session, roots, err := graph.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	if resolver.callCount() != 0 {
		t.Fatal("session start resolved a resource snapshot")
	}

	first := session.Prepare(context.Background(), roots[0].Ref, "inspect", "")
	if first.Status != AccessReady {
		t.Fatalf("first status = %q, want ready; problem: %#v", first.Status, first.Problem)
	}
	second := session.Prepare(context.Background(), roots[0].Ref, "inspect", "")
	if second.Status != AccessReady {
		t.Fatalf("second status = %q, want ready; problem: %#v", second.Status, second.Problem)
	}
	if first.State.Digest == second.State.Digest {
		t.Fatalf("live state digest did not change: %s", first.State.Digest)
	}
	if got := second.State.Value.(map[string]any)["revision"]; got != 2 {
		t.Fatalf("second live revision = %v, want 2", got)
	}

	stale := session.Prepare(context.Background(), roots[0].Ref, "inspect", first.State.Digest)
	if stale.Status != AccessStale {
		t.Fatalf("stale status = %q, want %q", stale.Status, AccessStale)
	}
	if stale.Problem == nil || stale.Problem.Code != "stale_state" {
		t.Fatalf("stale problem = %#v, want stale_state", stale.Problem)
	}
	if stale.State != nil {
		t.Fatal("stale result exposed live state")
	}
	if got, want := resolver.callCount(), 3; got != want {
		t.Fatalf("resolver calls = %d, want %d", got, want)
	}
}

func TestApplyDeltaRejectsUndeclaredGrantAtomically(t *testing.T) {
	definition := loadBaseWorld(t)
	graph := NewGraph(definition, &sequenceResolver{values: []any{map[string]any{}}})
	session, roots, err := graph.StartSession()
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.ApplyDelta(roots[0].Ref, "inspect", Mutation{Grants: []GrantRequest{{
		Type: "repo", Label: "not declared", Resource: Resource{Kind: "logical", Value: "other"},
	}}})
	if err == nil || !strings.Contains(err.Error(), "does not declare grants of handle type") {
		t.Fatalf("ApplyDelta error = %v, want undeclared grant rejection", err)
	}
	if got, want := session.ReachableCount(), 1; got != want {
		t.Fatalf("reachable count = %d, want %d after rejected delta", got, want)
	}
}

func loadBaseWorld(t *testing.T) *Definition {
	t.Helper()
	definition, err := Load(schemaPath(t), writeWorld(t, baseWorldJSON))
	if err != nil {
		t.Fatal(err)
	}
	return definition
}

type sequenceResolver struct {
	mu     sync.Mutex
	values []any
	calls  int
}

func (resolver *sequenceResolver) Resolve(_ context.Context, resource Resource) (any, error) {
	resolver.mu.Lock()
	defer resolver.mu.Unlock()
	if resource.Value == "" {
		return nil, fmt.Errorf("empty resource")
	}
	if resolver.calls >= len(resolver.values) {
		return resolver.values[len(resolver.values)-1], nil
	}
	value := resolver.values[resolver.calls]
	resolver.calls++
	return value, nil
}

func (resolver *sequenceResolver) callCount() int {
	resolver.mu.Lock()
	defer resolver.mu.Unlock()
	return resolver.calls
}
