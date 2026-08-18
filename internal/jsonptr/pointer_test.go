package jsonptr

import "testing"

func TestResolveObjectsArraysAndEscapedTokens(t *testing.T) {
	value := map[string]any{"items": []any{map[string]any{"a/b": "found"}}}
	resolved, ok := Resolve(value, "/items/0/a~1b")
	if !ok || resolved != "found" {
		t.Fatalf("Resolve() = (%#v, %v), want found", resolved, ok)
	}
	for _, invalid := range []string{"items/0", "/items/2", "/items/~2"} {
		if _, ok := Resolve(value, invalid); ok {
			t.Fatalf("Resolve(%q) succeeded", invalid)
		}
	}
}
