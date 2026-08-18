package world

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ValidateDefinition checks cross-references that JSON Schema cannot express.
func ValidateDefinition(definition *Definition) error {
	if definition == nil {
		return fmt.Errorf("world definition is required")
	}
	roots, err := validateRoots(definition)
	if err != nil {
		return err
	}
	if err := validateHandleTypes(definition, roots); err != nil {
		return err
	}
	return validateTransitions(definition, roots)
}

func validateRoots(definition *Definition) (map[string]Root, error) {
	roots := make(map[string]Root, len(definition.Roots))
	for _, root := range definition.Roots {
		if _, exists := roots[root.Name]; exists {
			return nil, fmt.Errorf("duplicate root name %q", root.Name)
		}
		if _, exists := definition.HandleTypes[root.Type]; !exists {
			return nil, fmt.Errorf("root %q references unknown handle type %q", root.Name, root.Type)
		}
		roots[root.Name] = root
	}
	return roots, nil
}

func validateHandleTypes(definition *Definition, roots map[string]Root) error {
	for handleType, descriptor := range definition.HandleTypes {
		for verbName, verb := range descriptor.Verbs {
			if err := validateGrantTypes(definition, handleType, verbName, verb); err != nil {
				return err
			}
			if err := validateRefusals(definition, roots, handleType, verbName, verb); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateGrantTypes(definition *Definition, handleType, verbName string, verb Verb) error {
	for _, grantedType := range verb.Grants {
		if _, exists := definition.HandleTypes[grantedType]; !exists {
			return fmt.Errorf("%s.%s grants unknown handle type %q", handleType, verbName, grantedType)
		}
	}
	return nil
}

func validateRefusals(definition *Definition, roots map[string]Root, handleType, verbName string, verb Verb) error {
	refusalIDs := make(map[string]struct{}, len(verb.Refusals))
	for _, refusal := range verb.Refusals {
		if _, exists := refusalIDs[refusal.ID]; exists {
			return fmt.Errorf("%s.%s has duplicate refusal id %q", handleType, verbName, refusal.ID)
		}
		refusalIDs[refusal.ID] = struct{}{}
		if refusal.Instead == nil {
			continue
		}
		if err := validateCall(definition, roots, handleType, verb, *refusal.Instead, false); err != nil {
			return fmt.Errorf("%s.%s refusal %q: %w", handleType, verbName, refusal.ID, err)
		}
	}
	return nil
}

func validateTransitions(definition *Definition, roots map[string]Root) error {
	transitionIDs := make(map[string]struct{}, len(definition.Transitions))
	for _, transition := range definition.Transitions {
		if _, exists := transitionIDs[transition.ID]; exists {
			return fmt.Errorf("duplicate transition id %q", transition.ID)
		}
		transitionIDs[transition.ID] = struct{}{}
		descriptor, exists := definition.HandleTypes[transition.Match.HandleType]
		if !exists {
			return fmt.Errorf("transition %q matches unknown handle type %q", transition.ID, transition.Match.HandleType)
		}
		verb, exists := descriptor.Verbs[transition.Match.Verb]
		if !exists {
			return fmt.Errorf("transition %q match verb %q is not attached to handle type %q", transition.ID, transition.Match.Verb, transition.Match.HandleType)
		}
		if transition.Match.Status == "absent" && len(transition.Suggestions) > 0 {
			return fmt.Errorf("transition %q cannot attach suggestions to absent results", transition.ID)
		}
		for index, suggestion := range transition.Suggestions {
			if err := validateCall(definition, roots, transition.Match.HandleType, verb, suggestion.Call, true); err != nil {
				return fmt.Errorf("transition %q suggestion %d: %w", transition.ID, index, err)
			}
		}
	}
	return nil
}

func validateCall(definition *Definition, roots map[string]Root, sourceType string, sourceVerb Verb, call CallTemplate, allowResult bool) error {
	targetTypes, err := selectedTypes(definition, roots, sourceType, sourceVerb, call.Handle, allowResult)
	if err != nil {
		return err
	}
	for _, targetType := range targetTypes {
		descriptor := definition.HandleTypes[targetType]
		targetVerb, exists := descriptor.Verbs[call.Verb]
		if !exists {
			return fmt.Errorf("verb %q is not attached to selected handle type %q", call.Verb, targetType)
		}
		if err := validateArguments(call.Args, targetVerb.ArgsSchema, sourceVerb.ResultSchema, allowResult); err != nil {
			return fmt.Errorf("call to %s.%s: %w", targetType, call.Verb, err)
		}
	}
	if call.State != nil {
		if err := validateBinding(*call.State, sourceVerb.ResultSchema, allowResult); err != nil {
			return fmt.Errorf("state binding: %w", err)
		}
	}
	return nil
}

func selectedTypes(definition *Definition, roots map[string]Root, sourceType string, sourceVerb Verb, selector HandleSelector, allowResult bool) ([]string, error) {
	switch selector.Source {
	case "self":
		return []string{sourceType}, nil
	case "root":
		root, exists := roots[selector.Name]
		if !exists {
			return nil, fmt.Errorf("selected root %q does not exist", selector.Name)
		}
		return []string{root.Type}, nil
	case "result":
		if !allowResult {
			return nil, fmt.Errorf("result handles are unavailable in a refusal alternative")
		}
		if !pointerResolves(sourceVerb.ResultSchema, selector.ResultPointer) {
			return nil, fmt.Errorf("result handle pointer %q does not resolve", selector.ResultPointer)
		}
		if len(sourceVerb.Grants) == 0 {
			return nil, fmt.Errorf("result handle selector is used by a verb that declares no grants")
		}
		return sourceVerb.Grants, nil
	default:
		return nil, fmt.Errorf("unsupported handle selector source %q", selector.Source)
	}
}

func validateArguments(arguments map[string]Binding, argumentSchema, resultSchema json.RawMessage, allowResult bool) error {
	var schema map[string]any
	if err := json.Unmarshal(argumentSchema, &schema); err != nil {
		return fmt.Errorf("decode argument schema: %w", err)
	}
	if required, ok := schema["required"].([]any); ok {
		for _, item := range required {
			name, _ := item.(string)
			if _, exists := arguments[name]; !exists {
				return fmt.Errorf("does not bind required argument %q", name)
			}
		}
	}
	if additional, exists := schema["additionalProperties"].(bool); exists && !additional {
		properties, _ := schema["properties"].(map[string]any)
		for name := range arguments {
			if _, exists := properties[name]; !exists {
				return fmt.Errorf("binds undeclared argument %q", name)
			}
		}
	}
	for name, binding := range arguments {
		if err := validateBinding(binding, resultSchema, allowResult); err != nil {
			return fmt.Errorf("argument %q: %w", name, err)
		}
	}
	return nil
}

func validateBinding(binding Binding, resultSchema json.RawMessage, allowResult bool) error {
	if binding.ResultPointer == nil {
		return nil
	}
	if !allowResult {
		return fmt.Errorf("result pointer %q is unavailable", *binding.ResultPointer)
	}
	if !pointerResolves(resultSchema, *binding.ResultPointer) {
		return fmt.Errorf("result pointer %q does not resolve", *binding.ResultPointer)
	}
	return nil
}

func pointerResolves(rawSchema json.RawMessage, pointer string) bool {
	var schema any
	if err := json.Unmarshal(rawSchema, &schema); err != nil {
		return false
	}
	if pointer == "" {
		return true
	}
	if !strings.HasPrefix(pointer, "/") {
		return false
	}
	current := schema
	for _, encoded := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(encoded, "~1", "/"), "~0", "~")
		member, ok := current.(map[string]any)
		if !ok {
			return false
		}
		if properties, ok := member["properties"].(map[string]any); ok {
			if next, exists := properties[token]; exists {
				current = next
				continue
			}
			if additional, exists := member["additionalProperties"].(bool); !exists || additional {
				current = map[string]any{}
				continue
			}
			return false
		}
		if items, exists := member["items"]; exists {
			if _, err := strconv.ParseUint(token, 10, 64); err != nil {
				return false
			}
			current = items
			continue
		}
		return false
	}
	return true
}
