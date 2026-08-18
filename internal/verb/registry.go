package verb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

type Definition struct {
	HandleType   string
	Name         string
	ArgsSchema   json.RawMessage
	ResultSchema json.RawMessage
	Handler      Handler
}

type registered struct {
	definition Definition
	args       *jsonschema.Schema
	result     *jsonschema.Schema
}

type Registry struct {
	mu      sync.RWMutex
	entries map[string]registered
}

func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]registered)}
}

func (registry *Registry) Register(definition Definition) error {
	if registry == nil {
		return fmt.Errorf("registry is required")
	}
	if !identifierPattern.MatchString(definition.HandleType) {
		return fmt.Errorf("handle type %q is invalid", definition.HandleType)
	}
	if !identifierPattern.MatchString(definition.Name) {
		return fmt.Errorf("verb name %q is invalid", definition.Name)
	}
	if definition.Handler == nil {
		return fmt.Errorf("%s.%s handler is required", definition.HandleType, definition.Name)
	}
	args, err := compileTypedSchema(definition.HandleType+"."+definition.Name+" argument", definition.ArgsSchema, "object")
	if err != nil {
		return err
	}
	result, err := compileTypedSchema(definition.HandleType+"."+definition.Name+" result", definition.ResultSchema, "object")
	if err != nil {
		return err
	}

	key := registryKey(definition.HandleType, definition.Name)
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, exists := registry.entries[key]; exists {
		return fmt.Errorf("verb %s.%s is already registered", definition.HandleType, definition.Name)
	}
	registry.entries[key] = registered{definition: definition, args: args, result: result}
	return nil
}

func (registry *Registry) lookup(handleType, name string) (registered, bool) {
	if registry == nil {
		return registered{}, false
	}
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	entry, exists := registry.entries[registryKey(handleType, name)]
	return entry, exists
}

func compileTypedSchema(label string, raw json.RawMessage, requiredType string) (*jsonschema.Schema, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("%s schema is required", label)
	}
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode %s schema: %w", label, err)
	}
	members, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s schema must be an object", label)
	}
	declaredType, ok := members["type"].(string)
	if !ok || declaredType == "" {
		return nil, fmt.Errorf("%s schema must declare a type", label)
	}
	if requiredType != "" && declaredType != requiredType {
		return nil, fmt.Errorf("%s schema type is %q, want %q", label, declaredType, requiredType)
	}

	compiler := jsonschema.NewCompiler()
	const resource = "urn:phoenix:verb-schema"
	if err := compiler.AddResource(resource, value); err != nil {
		return nil, fmt.Errorf("add %s schema: %w", label, err)
	}
	compiled, err := compiler.Compile(resource)
	if err != nil {
		return nil, fmt.Errorf("compile %s schema: %w", label, err)
	}
	return compiled, nil
}

func registryKey(handleType, name string) string {
	return handleType + "\x00" + name
}
