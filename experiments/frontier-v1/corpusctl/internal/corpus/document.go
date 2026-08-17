// Package corpus loads, validates, digests, and manifests the frontier-v1
// task corpus. Every document is read once and validated against the JSON
// Schemas in experiments/frontier-v1/schema.
package corpus

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	experimentDir = "experiments/frontier-v1"
	schemaDir     = experimentDir + "/schema"
)

// Document is one JSON file on disk with its decoded value and canonical digest.
type Document struct {
	// Path is repository-relative and slash-separated so messages and
	// manifests read the same on every platform.
	Path   string
	Digest string
	value  any
	data   []byte
}

// Name is the document's file name without the .json suffix.
func (document Document) Name() string {
	return strings.TrimSuffix(path.Base(document.Path), ".json")
}

// Into decodes the document into a typed value. Shape is already guaranteed by
// the schema, so this only projects the fields the tool needs.
func (document Document) Into(target any) error {
	if err := json.Unmarshal(document.data, target); err != nil {
		return fmt.Errorf("%s: %w", document.Path, err)
	}
	return nil
}

func loadDocument(root, relative string) (Document, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
	if err != nil {
		return Document{}, err
	}
	value, err := DecodeJSON(data)
	if err != nil {
		return Document{}, fmt.Errorf("%s: %w", relative, err)
	}
	digest, err := Digest(value)
	if err != nil {
		return Document{}, fmt.Errorf("%s: %w", relative, err)
	}
	return Document{Path: relative, Digest: digest, value: value, data: data}, nil
}

// DigestFile returns the canonical digest of a JSON file addressed by any path.
func DigestFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	value, err := DecodeJSON(data)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return Digest(value)
}

// schemaSet holds the compiled corpus schemas keyed by their file stem, for
// example "case" for case.schema.json.
type schemaSet map[string]*jsonschema.Schema

func loadSchemas(root string) (schemaSet, error) {
	pattern := filepath.Join(root, filepath.FromSlash(schemaDir), "*.schema.json")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no schemas found at %s", pattern)
	}
	compiler := jsonschema.NewCompiler()
	identifiers := make(map[string]string, len(paths))
	for _, schemaPath := range paths {
		identifier, value, err := readSchema(schemaPath)
		if err != nil {
			return nil, err
		}
		if err := compiler.AddResource(identifier, value); err != nil {
			return nil, fmt.Errorf("%s: %w", schemaPath, err)
		}
		identifiers[strings.TrimSuffix(filepath.Base(schemaPath), ".schema.json")] = identifier
	}
	schemas := make(schemaSet, len(identifiers))
	for name, identifier := range identifiers {
		schema, err := compiler.Compile(identifier)
		if err != nil {
			return nil, fmt.Errorf("compile %s.schema.json: %w", name, err)
		}
		schemas[name] = schema
	}
	return schemas, nil
}

func readSchema(schemaPath string) (string, any, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", nil, err
	}
	value, err := DecodeJSON(data)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", schemaPath, err)
	}
	members, ok := value.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("%s: schema must be a JSON object", schemaPath)
	}
	identifier, ok := members["$id"].(string)
	if !ok {
		return "", nil, fmt.Errorf("%s: schema must declare a string $id", schemaPath)
	}
	return identifier, value, nil
}

func (schemas schemaSet) validate(name string, document Document) error {
	schema, ok := schemas[name]
	if !ok {
		return fmt.Errorf("%s.schema.json was not compiled", name)
	}
	if err := schema.Validate(document.value); err != nil {
		return fmt.Errorf("%s violates %s.schema.json: %w", document.Path, name, err)
	}
	return nil
}

// validateValue checks an already-decoded value that has no file of its own,
// such as a manifest being generated.
func (schemas schemaSet) validateValue(name string, value any) error {
	schema, ok := schemas[name]
	if !ok {
		return fmt.Errorf("%s.schema.json was not compiled", name)
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("generated %s violates %s.schema.json: %w", name, name, err)
	}
	return nil
}

func validateWithSchema(root, schemaPath string, document Document) error {
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(filepath.Join(root, filepath.FromSlash(schemaPath)))
	if err != nil {
		return fmt.Errorf("compile %s: %w", schemaPath, err)
	}
	if err := schema.Validate(document.value); err != nil {
		return fmt.Errorf("%s violates %s: %w", document.Path, schemaPath, err)
	}
	return nil
}
