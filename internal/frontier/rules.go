// Package frontier compiles authored transition rules and binds their calls to
// the live topology after each act.
package frontier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/nesste/phoenix/internal/jsonptr"
	"github.com/nesste/phoenix/internal/world"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	MaxEntries       = 3
	MaxWhyCharacters = 80
)

var (
	identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
	digestPattern     = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type Observation struct {
	HandleType string
	Handle     string
	Verb       string
	Status     string
	Result     any
	State      *world.LiveState
}

type Call struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	State  *string        `json:"state,omitempty"`
}

type Entry struct {
	Call       Call    `json:"call"`
	Why        string  `json:"why"`
	Provenance string  `json:"provenance"`
	Score      float64 `json:"score"`
}

type Topology interface {
	RootHandle(string) (world.Handle, bool)
	ReachableHandle(string) (world.Handle, bool)
}

type Engine struct {
	rules []compiledRule
	verbs map[string]*jsonschema.Schema
}

type compiledRule struct {
	definition world.Transition
	matcher    *jsonschema.Schema
}

func New(definition *world.Definition) (*Engine, error) {
	if definition == nil {
		return nil, fmt.Errorf("world definition is required")
	}
	for _, transition := range definition.Transitions {
		if err := validateTransition(transition); err != nil {
			return nil, fmt.Errorf("transition %q: %w", transition.ID, err)
		}
	}
	if err := world.ValidateDefinition(definition); err != nil {
		return nil, fmt.Errorf("validate authored rules: %w", err)
	}
	engine := &Engine{verbs: make(map[string]*jsonschema.Schema)}
	if err := engine.compileVerbs(definition); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(definition.Transitions))
	for _, transition := range definition.Transitions {
		if _, exists := seen[transition.ID]; exists {
			return nil, fmt.Errorf("duplicate transition id %q", transition.ID)
		}
		seen[transition.ID] = struct{}{}
		matcher, err := compileSchema("transition "+transition.ID+" result matcher", transition.Match.ResultWhen)
		if err != nil {
			return nil, err
		}
		engine.rules = append(engine.rules, compiledRule{definition: transition, matcher: matcher})
	}
	return engine, nil
}

func (engine *Engine) compileVerbs(definition *world.Definition) error {
	for handleType, descriptor := range definition.HandleTypes {
		for name, verb := range descriptor.Verbs {
			compiled, err := compileSchema(handleType+"."+name+" argument", verb.ArgsSchema)
			if err != nil {
				return err
			}
			engine.verbs[verbKey(handleType, name)] = compiled
		}
	}
	return nil
}

func validateTransition(transition world.Transition) error {
	if !identifierPattern.MatchString(transition.ID) {
		return fmt.Errorf("invalid id")
	}
	if !validStatus(transition.Match.Status) {
		return fmt.Errorf("invalid match status %q", transition.Match.Status)
	}
	if len(transition.Suggestions) > MaxEntries {
		return fmt.Errorf("contains more than %d suggestions", MaxEntries)
	}
	for index, suggestion := range transition.Suggestions {
		if utf8.RuneCountInString(suggestion.Why) == 0 {
			return fmt.Errorf("suggestion %d has an empty why line", index)
		}
		if utf8.RuneCountInString(suggestion.Why) > MaxWhyCharacters {
			return fmt.Errorf("suggestion %d why line exceeds %d characters", index, MaxWhyCharacters)
		}
		if suggestion.Score < 0 || suggestion.Score > 1 {
			return fmt.Errorf("suggestion %d score is outside 0..1", index)
		}
		if err := validateTemplate(suggestion.Call); err != nil {
			return fmt.Errorf("suggestion %d: %w", index, err)
		}
	}
	return nil
}

func validateTemplate(call world.CallTemplate) error {
	if !identifierPattern.MatchString(call.Verb) {
		return fmt.Errorf("invalid verb %q; calls must be structured", call.Verb)
	}
	if call.Args == nil {
		return fmt.Errorf("call arguments must be a bound object")
	}
	for name, binding := range call.Args {
		if err := validateBinding(binding); err != nil {
			return fmt.Errorf("argument %q: %w", name, err)
		}
	}
	if call.State != nil {
		if err := validateBinding(*call.State); err != nil {
			return fmt.Errorf("state binding: %w", err)
		}
	}
	return nil
}

func validateBinding(binding world.Binding) error {
	count := 0
	if binding.Literal != nil {
		count++
		var value any
		if err := json.Unmarshal(binding.Literal, &value); err != nil {
			return fmt.Errorf("literal is invalid JSON")
		}
	}
	if binding.ResultPointer != nil {
		count++
		if !jsonptr.Valid(*binding.ResultPointer) {
			return fmt.Errorf("result pointer %q is invalid", *binding.ResultPointer)
		}
	}
	if binding.StatePointer != nil {
		count++
		if !jsonptr.Valid(*binding.StatePointer) {
			return fmt.Errorf("state pointer %q is invalid", *binding.StatePointer)
		}
	}
	if count != 1 {
		return fmt.Errorf("binding must select exactly one literal, result pointer, or state pointer")
	}
	return nil
}

func compileSchema(label string, raw json.RawMessage) (*jsonschema.Schema, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("%s schema is required", label)
	}
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode %s schema: %w", label, err)
	}
	compiler := jsonschema.NewCompiler()
	const resource = "urn:phoenix:frontier-schema"
	if err := compiler.AddResource(resource, value); err != nil {
		return nil, fmt.Errorf("add %s schema: %w", label, err)
	}
	compiled, err := compiler.Compile(resource)
	if err != nil {
		return nil, fmt.Errorf("compile %s schema: %w", label, err)
	}
	return compiled, nil
}

func validStatus(status string) bool {
	return status == "ok" || status == "fail" || status == "refused" || status == "absent"
}

func verbKey(handleType, verb string) string {
	return handleType + "\x00" + verb
}
