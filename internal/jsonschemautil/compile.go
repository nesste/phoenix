// Package jsonschemautil contains small shared helpers for authored JSON
// Schema compilation.
package jsonschemautil

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func Compile(label, resource string, raw json.RawMessage) (*jsonschema.Schema, error) {
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode %s schema: %w", label, err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(resource, value); err != nil {
		return nil, fmt.Errorf("add %s schema: %w", label, err)
	}
	compiled, err := compiler.Compile(resource)
	if err != nil {
		return nil, fmt.Errorf("compile %s schema: %w", label, err)
	}
	return compiled, nil
}
