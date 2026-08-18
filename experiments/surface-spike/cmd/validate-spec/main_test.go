package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestResultSchemaRejectsUnknownFields(t *testing.T) {
	schema, instance := loadResultFixture(t, "result.ok.json")
	instance["unexpected"] = true
	if err := schema.Validate(instance); err == nil {
		t.Fatal("result with an unknown field passed validation")
	}
}

func TestResultSchemaEnforcesFrontierCap(t *testing.T) {
	schema, instance := loadResultFixture(t, "result.ok.json")
	frontier := instance["frontier"].([]any)
	instance["frontier"] = []any{frontier[0], frontier[0], frontier[0], frontier[0]}
	if err := schema.Validate(instance); err == nil {
		t.Fatal("result with four frontier entries passed validation")
	}
}

func TestAbsentResultCannotCarryGuidance(t *testing.T) {
	schema, instance := loadResultFixture(t, "result.absent.json")
	_, okResult := loadResultFixture(t, "result.ok.json")
	instance["frontier"] = okResult["frontier"]
	if err := schema.Validate(instance); err == nil {
		t.Fatal("absent result with a frontier passed validation")
	}
}

func TestAbsentResultCannotCarryDetailsOrHandleMutations(t *testing.T) {
	tests := []struct {
		name string
		edit func(map[string]any)
	}{
		{
			name: "error details",
			edit: func(instance map[string]any) {
				instance["error"].(map[string]any)["details"] = map[string]any{"leak": true}
			},
		},
		{
			name: "handle grant",
			edit: func(instance map[string]any) {
				instance["handles"].(map[string]any)["grant"] = []any{map[string]any{
					"ref": "h_1111111111111111", "type": "repo", "label": "leaked handle",
				}}
			},
		},
		{
			name: "handle revoke",
			edit: func(instance map[string]any) {
				instance["handles"].(map[string]any)["revoke"] = []any{"h_1111111111111111"}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schema, instance := loadResultFixture(t, "result.absent.json")
			test.edit(instance)
			if err := schema.Validate(instance); err == nil {
				t.Fatal("absent result carrying hidden state passed validation")
			}
		})
	}
}

func loadResultFixture(t *testing.T, name string) (*jsonschema.Schema, map[string]any) {
	t.Helper()
	root := repositoryRoot(t)
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(filepath.Join(root, "spec", "result.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "spec", "examples", name))
	if err != nil {
		t.Fatal(err)
	}
	var instance map[string]any
	if err := json.Unmarshal(data, &instance); err != nil {
		t.Fatal(err)
	}
	return schema, instance
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve test source path")
	}
	root, err := filepath.Abs(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}
