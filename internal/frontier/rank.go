package frontier

import (
	"encoding/json"
	"sort"

	"github.com/nesste/phoenix/internal/jsonptr"
	"github.com/nesste/phoenix/internal/world"
)

type candidate struct {
	entry Entry
	key   string
	rule  string
}

func (engine *Engine) Compute(topology Topology, observation Observation) []Entry {
	if engine == nil || topology == nil {
		return []Entry{}
	}
	candidates := make([]candidate, 0, MaxEntries)
	for _, rule := range engine.rules {
		match := rule.definition.Match
		if match.HandleType != observation.HandleType || match.Verb != observation.Verb || match.Status != observation.Status {
			continue
		}
		if rule.matcher.Validate(observation.Result) != nil {
			continue
		}
		for _, suggestion := range rule.definition.Suggestions {
			call, ok := engine.bindCall(topology, observation, suggestion.Call)
			if !ok {
				continue
			}
			encoded, err := json.Marshal(call)
			if err != nil {
				continue
			}
			candidates = append(candidates, candidate{
				entry: Entry{Call: call, Why: suggestion.Why, Provenance: "authored", Score: suggestion.Score},
				key:   string(encoded), rule: rule.definition.ID,
			})
		}
	}
	return rank(candidates)
}

func (engine *Engine) bindCall(topology Topology, observation Observation, template world.CallTemplate) (Call, bool) {
	handle, ok := bindHandle(topology, observation, template.Handle)
	if !ok {
		return Call{}, false
	}
	args := make(map[string]any, len(template.Args))
	for name, binding := range template.Args {
		value, ok := bindValue(binding, observation)
		if !ok {
			return Call{}, false
		}
		args[name] = value
	}
	validator, exists := engine.verbs[verbKey(handle.Type, template.Verb)]
	if !exists || validator.Validate(args) != nil {
		return Call{}, false
	}
	call := Call{Handle: handle.Ref, Verb: template.Verb, Args: args}
	if template.State != nil {
		value, ok := bindValue(*template.State, observation)
		state, stringValue := value.(string)
		if !ok || !stringValue || !digestPattern.MatchString(state) {
			return Call{}, false
		}
		call.State = &state
	}
	return call, true
}

func bindHandle(topology Topology, observation Observation, selector world.HandleSelector) (world.Handle, bool) {
	var handle world.Handle
	var ok bool
	switch selector.Source {
	case "self":
		handle, ok = topology.ReachableHandle(observation.Handle)
	case "root":
		handle, ok = topology.RootHandle(selector.Name)
	case "result":
		value, found := jsonptr.Resolve(observation.Result, selector.ResultPointer)
		if !found {
			return world.Handle{}, false
		}
		ref, isString := value.(string)
		if !isString {
			return world.Handle{}, false
		}
		handle, ok = topology.ReachableHandle(ref)
	}
	return handle, ok
}

func bindValue(binding world.Binding, observation Observation) (any, bool) {
	if binding.Literal != nil {
		var value any
		if json.Unmarshal(binding.Literal, &value) != nil {
			return nil, false
		}
		return value, true
	}
	if binding.ResultPointer != nil {
		return jsonptr.Resolve(observation.Result, *binding.ResultPointer)
	}
	if binding.StatePointer != nil && observation.State != nil {
		return jsonptr.Resolve(observation.State.Value, *binding.StatePointer)
	}
	return nil, false
}

func rank(candidates []candidate) []Entry {
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].entry.Score != candidates[j].entry.Score {
			return candidates[i].entry.Score > candidates[j].entry.Score
		}
		if candidates[i].key != candidates[j].key {
			return candidates[i].key < candidates[j].key
		}
		if candidates[i].entry.Why != candidates[j].entry.Why {
			return candidates[i].entry.Why < candidates[j].entry.Why
		}
		return candidates[i].rule < candidates[j].rule
	})
	entries := make([]Entry, 0, min(len(candidates), MaxEntries))
	seen := make(map[string]struct{}, len(candidates))
	for _, item := range candidates {
		if _, duplicate := seen[item.key]; duplicate {
			continue
		}
		seen[item.key] = struct{}{}
		entries = append(entries, item.entry)
		if len(entries) == MaxEntries {
			break
		}
	}
	return entries
}
