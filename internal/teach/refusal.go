// Package teach compiles state-sensitive refusal rules and binds their
// alternatives to a session's live topology.
package teach

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/jsonptr"
	"github.com/nesste/phoenix/internal/world"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
	digestPattern     = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type Observation struct {
	HandleType string
	Handle     string
	Verb       string
	Args       map[string]any
	State      *world.LiveState
	Failure    map[string]any
}

type Refusal struct {
	What          string         `json:"what"`
	Why           string         `json:"why"`
	Instead       *frontier.Call `json:"instead"`
	NoAlternative string         `json:"no_alternative,omitempty"`
}

type Topology interface {
	RootHandle(string) (world.Handle, bool)
	ReachableHandle(string) (world.Handle, bool)
}

type Engine struct {
	rules map[string][]compiledRule
	verbs map[string]*jsonschema.Schema
}

type compiledRule struct {
	definition world.RefusalRule
	matcher    *jsonschema.Schema
}

func New(definition *world.Definition) (*Engine, error) {
	if definition == nil {
		return nil, fmt.Errorf("world definition is required")
	}
	if err := validateRefusalShapes(definition); err != nil {
		return nil, err
	}
	if err := world.ValidateDefinition(definition); err != nil {
		return nil, fmt.Errorf("validate authored refusals: %w", err)
	}
	engine := &Engine{
		rules: make(map[string][]compiledRule),
		verbs: make(map[string]*jsonschema.Schema),
	}
	for handleType, descriptor := range definition.HandleTypes {
		for name, verb := range descriptor.Verbs {
			key := verbKey(handleType, name)
			arguments, err := compileSchema(handleType+"."+name+" argument", verb.ArgsSchema)
			if err != nil {
				return nil, err
			}
			engine.verbs[key] = arguments
			for _, rule := range verb.Refusals {
				matcher, err := compileSchema(handleType+"."+name+" refusal "+rule.ID, rule.When)
				if err != nil {
					return nil, err
				}
				engine.rules[key] = append(engine.rules[key], compiledRule{definition: rule, matcher: matcher})
			}
		}
	}
	return engine, nil
}

func (engine *Engine) Evaluate(topology Topology, observation Observation) (Refusal, bool, error) {
	if engine == nil || topology == nil {
		return Refusal{}, false, nil
	}
	features := observationFeatures(observation)
	for _, rule := range engine.rules[verbKey(observation.HandleType, observation.Verb)] {
		if rule.matcher.Validate(features) != nil {
			continue
		}
		refusal := Refusal{
			What: rule.definition.What, Why: rule.definition.Why,
			NoAlternative: rule.definition.NoAlternative,
		}
		if rule.definition.Instead == nil {
			return refusal, true, nil
		}
		instead, ok := engine.bindAlternative(topology, observation, *rule.definition.Instead)
		if !ok {
			return Refusal{}, false, fmt.Errorf("refusal %q alternative is not reachable and bound", rule.definition.ID)
		}
		refusal.Instead = &instead
		return refusal, true, nil
	}
	return Refusal{}, false, nil
}

func (engine *Engine) bindAlternative(topology Topology, observation Observation, template world.CallTemplate) (frontier.Call, bool) {
	handle, ok := alternativeHandle(topology, observation, template.Handle)
	if !ok {
		return frontier.Call{}, false
	}
	args := make(map[string]any, len(template.Args))
	for name, binding := range template.Args {
		value, ok := bindValue(binding, observation)
		if !ok {
			return frontier.Call{}, false
		}
		args[name] = value
	}
	validator, exists := engine.verbs[verbKey(handle.Type, template.Verb)]
	if !exists || validator.Validate(args) != nil {
		return frontier.Call{}, false
	}
	call := frontier.Call{Handle: handle.Ref, Verb: template.Verb, Args: args}
	if template.State != nil {
		value, ok := bindValue(*template.State, observation)
		state, stringValue := value.(string)
		if !ok || !stringValue || !digestPattern.MatchString(state) {
			return frontier.Call{}, false
		}
		call.State = &state
	}
	return call, true
}

func alternativeHandle(topology Topology, observation Observation, selector world.HandleSelector) (world.Handle, bool) {
	switch selector.Source {
	case "self":
		return topology.ReachableHandle(observation.Handle)
	case "root":
		return topology.RootHandle(selector.Name)
	default:
		return world.Handle{}, false
	}
}

func bindValue(binding world.Binding, observation Observation) (any, bool) {
	if binding.Literal != nil {
		var value any
		if json.Unmarshal(binding.Literal, &value) != nil {
			return nil, false
		}
		return value, true
	}
	if binding.StatePointer != nil && observation.State != nil {
		return jsonptr.Resolve(observation.State.Value, *binding.StatePointer)
	}
	return nil, false
}

func observationFeatures(observation Observation) map[string]any {
	args := observation.Args
	if args == nil {
		args = map[string]any{}
	}
	features := map[string]any{"args": args}
	if observation.State != nil {
		features["state"] = observation.State.Value
		features["state_digest"] = observation.State.Digest
	}
	if observation.Failure != nil {
		features["failure"] = observation.Failure
	}
	return features
}

func validateRefusalShapes(definition *world.Definition) error {
	for handleType, descriptor := range definition.HandleTypes {
		for verb, definition := range descriptor.Verbs {
			seen := make(map[string]struct{}, len(definition.Refusals))
			for _, rule := range definition.Refusals {
				if _, exists := seen[rule.ID]; exists {
					return fmt.Errorf("%s.%s has duplicate refusal id %q", handleType, verb, rule.ID)
				}
				seen[rule.ID] = struct{}{}
				if err := validateRefusalShape(rule); err != nil {
					return fmt.Errorf("%s.%s refusal %q: %w", handleType, verb, rule.ID, err)
				}
			}
		}
	}
	return nil
}

func validateRefusalShape(rule world.RefusalRule) error {
	if !identifierPattern.MatchString(rule.ID) {
		return fmt.Errorf("invalid id")
	}
	if strings.TrimSpace(rule.What) == "" {
		return fmt.Errorf("what is required")
	}
	if strings.TrimSpace(rule.Why) == "" {
		return fmt.Errorf("why is required")
	}
	if len(rule.When) == 0 {
		return fmt.Errorf("when schema is required")
	}
	if rule.Instead == nil {
		if strings.TrimSpace(rule.NoAlternative) == "" {
			return fmt.Errorf("no_alternative is required when instead is null")
		}
		return nil
	}
	if rule.NoAlternative != "" {
		return fmt.Errorf("no_alternative must be empty when instead is present")
	}
	if rule.Instead.Args == nil {
		return fmt.Errorf("instead arguments must be a bound object")
	}
	for name, binding := range rule.Instead.Args {
		if err := validateBinding(binding); err != nil {
			return fmt.Errorf("instead argument %q: %w", name, err)
		}
	}
	if rule.Instead.State != nil {
		if err := validateBinding(*rule.Instead.State); err != nil {
			return fmt.Errorf("instead state: %w", err)
		}
	}
	return nil
}

func validateBinding(binding world.Binding) error {
	count := 0
	if binding.Literal != nil {
		count++
		if !json.Valid(binding.Literal) {
			return fmt.Errorf("literal is invalid JSON")
		}
	}
	if binding.ResultPointer != nil {
		count++
		return fmt.Errorf("result pointers are unavailable in refusal alternatives")
	}
	if binding.StatePointer != nil {
		count++
		if !jsonptr.Valid(*binding.StatePointer) {
			return fmt.Errorf("state pointer %q is invalid", *binding.StatePointer)
		}
	}
	if count != 1 {
		return fmt.Errorf("binding must select exactly one literal or state pointer")
	}
	return nil
}

func compileSchema(label string, raw json.RawMessage) (*jsonschema.Schema, error) {
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode %s schema: %w", label, err)
	}
	compiler := jsonschema.NewCompiler()
	const resource = "urn:phoenix:teaching-schema"
	if err := compiler.AddResource(resource, value); err != nil {
		return nil, fmt.Errorf("add %s schema: %w", label, err)
	}
	compiled, err := compiler.Compile(resource)
	if err != nil {
		return nil, fmt.Errorf("compile %s schema: %w", label, err)
	}
	return compiled, nil
}

func verbKey(handleType, verb string) string {
	return handleType + "\x00" + verb
}
