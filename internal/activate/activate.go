// Package activate binds a natural-language intent to a bounded set of
// world-authored calls over handles that are already reachable in a session.
package activate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"unicode/utf8"

	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/jsonschemautil"
	"github.com/nesste/phoenix/internal/world"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	maxEntries        = 3
	inspectFallbackID = "inspect_reachable_repository"
)

type Topology interface {
	RootHandle(string) (world.Handle, bool)
	ReachableHandle(string) (world.Handle, bool)
}

type Engine struct {
	rules []compiledRule
	verbs map[string]*jsonschema.Schema
}

type compiledRule struct {
	definition world.ActivationRule
	pattern    *regexp.Regexp
}

type candidate struct {
	entry frontier.Entry
	key   string
	rule  string
}

func New(definition *world.Definition) (*Engine, error) {
	if definition == nil {
		return nil, fmt.Errorf("world definition is required")
	}
	if err := world.ValidateDefinition(definition); err != nil {
		return nil, fmt.Errorf("validate authored activations: %w", err)
	}
	engine := &Engine{verbs: make(map[string]*jsonschema.Schema)}
	for handleType, descriptor := range definition.HandleTypes {
		for name, verb := range descriptor.Verbs {
			compiled, err := jsonschemautil.Compile(handleType+"."+name+" argument", "urn:phoenix:activation-schema", verb.ArgsSchema)
			if err != nil {
				return nil, err
			}
			engine.verbs[verbKey(handleType, name)] = compiled
		}
	}
	for _, rule := range definition.Activations {
		if len(rule.Suggestions) > maxEntries {
			return nil, fmt.Errorf("activation %q contains more than %d suggestions", rule.ID, maxEntries)
		}
		if rule.ID == inspectFallbackID && len(rule.Suggestions) == 0 {
			return nil, fmt.Errorf("activation %q must include at least one suggestion", inspectFallbackID)
		}
		for index, suggestion := range rule.Suggestions {
			if utf8.RuneCountInString(suggestion.Why) == 0 || utf8.RuneCountInString(suggestion.Why) > 80 {
				return nil, fmt.Errorf("activation %q suggestion %d has an invalid why line", rule.ID, index)
			}
			if suggestion.Score < 0 || suggestion.Score > 1 {
				return nil, fmt.Errorf("activation %q suggestion %d score is outside 0..1", rule.ID, index)
			}
		}
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("compile activation %q: %w", rule.ID, err)
		}
		engine.rules = append(engine.rules, compiledRule{definition: rule, pattern: pattern})
	}
	return engine, nil
}

func (engine *Engine) Compute(topology Topology, anchor, intent string) []frontier.Entry {
	if engine == nil || topology == nil {
		return []frontier.Entry{}
	}
	anchorHandle, reachable := topology.ReachableHandle(anchor)
	if !reachable {
		return []frontier.Entry{}
	}
	candidates := []candidate{}
	matchedSpecific := false
	for _, rule := range engine.rules {
		if rule.definition.ID == inspectFallbackID {
			continue
		}
		captures := rule.pattern.FindStringSubmatch(intent)
		if captures == nil {
			continue
		}
		matchedSpecific = true
		for _, suggestion := range rule.definition.Suggestions {
			call, ok := engine.bind(topology, anchorHandle, captures, suggestion.Call)
			if !ok {
				continue
			}
			encoded, err := json.Marshal(call)
			if err != nil {
				continue
			}
			candidates = append(candidates, candidate{
				entry: frontier.Entry{Call: call, Why: suggestion.Why, Provenance: "activation", Score: suggestion.Score},
				key:   string(encoded), rule: rule.definition.ID,
			})
		}
	}
	ranked := rank(candidates)
	if len(ranked) > 0 || matchedSpecific {
		return ranked
	}
	return engine.bindFallback(topology, anchorHandle, intent)
}

func (engine *Engine) bindFallback(topology Topology, anchor world.Handle, intent string) []frontier.Entry {
	for _, rule := range engine.rules {
		if rule.definition.ID != inspectFallbackID {
			continue
		}
		candidates := []candidate{}
		captures := []string{intent}
		for _, suggestion := range rule.definition.Suggestions {
			call, ok := engine.bind(topology, anchor, captures, suggestion.Call)
			if !ok {
				continue
			}
			encoded, err := json.Marshal(call)
			if err != nil {
				continue
			}
			candidates = append(candidates, candidate{
				entry: frontier.Entry{Call: call, Why: suggestion.Why, Provenance: "activation", Score: suggestion.Score},
				key:   string(encoded), rule: rule.definition.ID,
			})
		}
		return rank(candidates)
	}
	return []frontier.Entry{}
}

func (engine *Engine) bind(topology Topology, anchor world.Handle, captures []string, template world.CallTemplate) (frontier.Call, bool) {
	var handle world.Handle
	var ok bool
	switch template.Handle.Source {
	case "self":
		handle, ok = topology.ReachableHandle(anchor.Ref)
	case "root":
		handle, ok = topology.RootHandle(template.Handle.Name)
	default:
		return frontier.Call{}, false
	}
	if !ok {
		return frontier.Call{}, false
	}
	args := make(map[string]any, len(template.Args))
	for name, binding := range template.Args {
		value, bound := bindValue(binding, captures)
		if !bound {
			return frontier.Call{}, false
		}
		args[name] = value
	}
	validator, exists := engine.verbs[verbKey(handle.Type, template.Verb)]
	if !exists || validator.Validate(args) != nil {
		return frontier.Call{}, false
	}
	return frontier.Call{Handle: handle.Ref, Verb: template.Verb, Args: args}, true
}

func bindValue(binding world.Binding, captures []string) (any, bool) {
	if binding.Literal != nil {
		var value any
		if json.Unmarshal(binding.Literal, &value) != nil {
			return nil, false
		}
		return value, true
	}
	if binding.IntentCapture == nil || *binding.IntentCapture >= len(captures) {
		return nil, false
	}
	value := captures[*binding.IntentCapture]
	return value, value != ""
}

func rank(candidates []candidate) []frontier.Entry {
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].entry.Score != candidates[j].entry.Score {
			return candidates[i].entry.Score > candidates[j].entry.Score
		}
		if candidates[i].key != candidates[j].key {
			return candidates[i].key < candidates[j].key
		}
		return candidates[i].rule < candidates[j].rule
	})
	entries := make([]frontier.Entry, 0, min(len(candidates), maxEntries))
	seen := make(map[string]struct{}, len(candidates))
	for _, item := range candidates {
		if _, duplicate := seen[item.key]; duplicate {
			continue
		}
		seen[item.key] = struct{}{}
		entries = append(entries, item.entry)
		if len(entries) == maxEntries {
			break
		}
	}
	return entries
}

func verbKey(handleType, verb string) string {
	return handleType + "\x00" + verb
}
