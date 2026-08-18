package world

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoadValidatesSchemaAndProducesStableDigest(t *testing.T) {
	first := writeWorld(t, baseWorldJSON)
	second := writeWorld(t, `{
  "transitions": [],
  "handle_types": {
    "detail": {"verbs": {"show": {"refusals": [], "may_revoke": false, "grants": [], "result_schema": {"type": "object"}, "args_schema": {"type": "object", "additionalProperties": false}}}},
    "repo": {"verbs": {"inspect": {"refusals": [], "may_revoke": true, "grants": ["detail"], "result_schema": {"type": "object", "additionalProperties": false}, "args_schema": {"type": "object", "additionalProperties": false}}}}
  },
  "roots": [{"resource": {"value": "repo", "kind": "logical"}, "label": "repository", "type": "repo", "name": "repo"}],
  "id": "test_world",
  "v": 1
}`)

	left, err := Load(schemaPath(t), first)
	if err != nil {
		t.Fatalf("Load(first): %v", err)
	}
	right, err := Load(schemaPath(t), second)
	if err != nil {
		t.Fatalf("Load(second): %v", err)
	}
	if left.Digest() != right.Digest() {
		t.Fatalf("equivalent world digests differ: %s != %s", left.Digest(), right.Digest())
	}
	if !strings.HasPrefix(left.Digest(), "sha256:") || len(left.Digest()) != len("sha256:")+64 {
		t.Fatalf("Digest() = %q, want a sha256 digest", left.Digest())
	}

	var unknown map[string]any
	if err := json.Unmarshal([]byte(baseWorldJSON), &unknown); err != nil {
		t.Fatal(err)
	}
	unknown["unexpected"] = true
	unknownPath := writeWorldValue(t, unknown)
	if _, err := Load(schemaPath(t), unknownPath); err == nil || !strings.Contains(err.Error(), "world schema") {
		t.Fatalf("Load(world with unknown field) error = %v, want schema rejection", err)
	}
}

func TestLoadRejectsInvalidSemanticReferences(t *testing.T) {
	tests := []struct {
		name string
		edit func(map[string]any)
		want string
	}{
		{
			name: "root handle type",
			edit: func(value map[string]any) {
				value["roots"].([]any)[0].(map[string]any)["type"] = "missing"
			},
			want: `root "repo" references unknown handle type "missing"`,
		},
		{
			name: "declared grant type",
			edit: func(value map[string]any) {
				verbs(value, "repo")["inspect"].(map[string]any)["grants"] = []any{"missing"}
			},
			want: `grants unknown handle type "missing"`,
		},
		{
			name: "transition verb",
			edit: func(value map[string]any) {
				value["transitions"] = []any{transition("missing", map[string]any{})}
			},
			want: `match verb "missing" is not attached to handle type "repo"`,
		},
		{
			name: "unbound required argument",
			edit: func(value map[string]any) {
				verbs(value, "repo")["use_value"] = map[string]any{
					"args_schema":   map[string]any{"type": "object", "required": []any{"value"}, "properties": map[string]any{"value": map[string]any{"type": "string"}}, "additionalProperties": false},
					"result_schema": map[string]any{"type": "object"}, "grants": []any{}, "may_revoke": false, "refusals": []any{},
				}
				value["transitions"] = []any{transition("inspect", map[string]any{"verb": "use_value"})}
			},
			want: `does not bind required argument "value"`,
		},
		{
			name: "result pointer",
			edit: func(value map[string]any) {
				verbs(value, "repo")["use_value"] = map[string]any{
					"args_schema":   map[string]any{"type": "object", "required": []any{"value"}, "properties": map[string]any{"value": map[string]any{"type": "string"}}, "additionalProperties": false},
					"result_schema": map[string]any{"type": "object"}, "grants": []any{}, "may_revoke": false, "refusals": []any{},
				}
				value["transitions"] = []any{transition("inspect", map[string]any{
					"verb": "use_value",
					"args": map[string]any{"value": map[string]any{"result_pointer": "/missing"}},
				})}
			},
			want: `result pointer "/missing" does not resolve`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := decodeWorld(t)
			test.edit(value)
			_, err := Load(schemaPath(t), writeWorldValue(t, value))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load() error = %v, want it to contain %q", err, test.want)
			}
		})
	}
}

func TestLoadAcceptedWorldExamples(t *testing.T) {
	root := repositoryRoot(t)
	paths := []string{
		filepath.Join(root, "spec", "examples", "world.dev_repo.json"),
		filepath.Join(root, "experiments", "frontier-v1", "worlds", "authoring.dev_repo.json"),
	}
	for _, path := range paths {
		if _, err := Load(schemaPath(t), path); err != nil {
			t.Fatalf("Load(%s): %v", filepath.ToSlash(path), err)
		}
	}
}

func TestEmptyResultPointerIsNotTreatedAsMissing(t *testing.T) {
	empty := ""
	err := validateBinding(Binding{ResultPointer: &empty}, json.RawMessage(`{"type":"object"}`), false)
	if err == nil || !strings.Contains(err.Error(), `result pointer "" is unavailable`) {
		t.Fatalf("validateBinding error = %v, want empty-pointer rejection", err)
	}
}

func transition(matchVerb string, overrides map[string]any) map[string]any {
	call := map[string]any{
		"handle": map[string]any{"source": "self"},
		"verb":   "show",
		"args":   map[string]any{},
	}
	for key, value := range overrides {
		call[key] = value
	}
	return map[string]any{
		"id":          "next",
		"match":       map[string]any{"handle_type": "repo", "verb": matchVerb, "status": "ok", "result_when": map[string]any{}},
		"suggestions": []any{map[string]any{"call": call, "why": "continue with reachable work", "score": 1}},
	}
}

func verbs(value map[string]any, handleType string) map[string]any {
	return value["handle_types"].(map[string]any)[handleType].(map[string]any)["verbs"].(map[string]any)
}

func decodeWorld(t *testing.T) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal([]byte(baseWorldJSON), &value); err != nil {
		t.Fatal(err)
	}
	return value
}

func writeWorldValue(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return writeWorld(t, string(encoded))
}

func writeWorld(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "world.json")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func schemaPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test source")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "spec", "world.schema.json")
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	return filepath.Clean(filepath.Join(filepath.Dir(schemaPath(t)), ".."))
}

const baseWorldJSON = `{
  "v": 1,
  "id": "test_world",
  "roots": [{"name": "repo", "type": "repo", "label": "repository", "resource": {"kind": "logical", "value": "repo"}}],
  "handle_types": {
    "repo": {"verbs": {"inspect": {
      "args_schema": {"type": "object", "additionalProperties": false},
      "result_schema": {"type": "object", "additionalProperties": false},
      "grants": ["detail"], "may_revoke": true, "refusals": []
    }}},
    "detail": {"verbs": {"show": {
      "args_schema": {"type": "object", "additionalProperties": false},
      "result_schema": {"type": "object"},
      "grants": [], "may_revoke": false, "refusals": []
    }}}
  },
  "transitions": []
}`
