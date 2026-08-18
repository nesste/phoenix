package devrepo

import (
	"encoding/json"

	"github.com/nesste/phoenix/internal/world"
)

// AuthoredTransitions is the initial static frontier for the dev-repo world.
// Task 1.8 attaches these rules when it assembles the concrete world roots.
func AuthoredTransitions() []world.Transition {
	return []world.Transition{
		{
			ID: "build_failure_next",
			Match: world.Match{
				HandleType: "repo", Verb: "build", Status: "ok",
				ResultWhen: raw(`{"properties":{"exit_code":{"minimum":1}},"required":["exit_code"]}`),
			},
			Suggestions: []world.Suggestion{
				suggest(rootCall("git", "diff"), "inspect changes related to the build failure", 1),
				suggest(rootCall("tests", "run"), "separate test failures from build failures", .8),
				suggest(callWithResult("episodes", "recall", "query", "/stderr"), "recall fixes for this failure output", .7),
			},
		},
		{
			ID: "test_failure_next",
			Match: world.Match{
				HandleType: "tests", Verb: "run", Status: "ok",
				ResultWhen: raw(`{"properties":{"failed":{"minimum":1}},"required":["failed"]}`),
			},
			Suggestions: []world.Suggestion{
				suggest(selfCall("list"), "list tests before focusing the failure", 1),
				suggest(callWithResult("episodes", "recall", "query", "/stderr"), "recall prior resolutions for this failure", .9),
				suggest(rootCall("git", "diff"), "inspect changes associated with the failure", .7),
			},
		},
		{
			ID: "edit_changed_next",
			Match: world.Match{
				HandleType: "repo", Verb: "edit", Status: "ok",
				ResultWhen: raw(`{"properties":{"changed":{"const":true}},"required":["changed"]}`),
			},
			Suggestions: []world.Suggestion{
				suggest(rootCall("tests", "run"), "run tests affected by the edit", 1),
				suggest(rootCall("git", "diff"), "inspect the resulting edit", .8),
			},
		},
		{
			ID: "test_list_focus_next",
			Match: world.Match{
				HandleType: "tests", Verb: "list", Status: "ok",
				ResultWhen: raw(`{"properties":{"tests":{"minItems":1}},"required":["tests"]}`),
			},
			Suggestions: []world.Suggestion{
				suggest(callWithResult("", "focus", "test", "/tests/0"), "focus the first bound test", 1),
			},
		},
	}
}

func rootCall(root, verb string) world.CallTemplate {
	return world.CallTemplate{
		Handle: world.HandleSelector{Source: "root", Name: root}, Verb: verb,
		Args: map[string]world.Binding{},
	}
}

func selfCall(verb string) world.CallTemplate {
	return world.CallTemplate{
		Handle: world.HandleSelector{Source: "self"}, Verb: verb,
		Args: map[string]world.Binding{},
	}
}

func callWithResult(root, verb, argument, pointer string) world.CallTemplate {
	call := selfCall(verb)
	if root != "" {
		call.Handle = world.HandleSelector{Source: "root", Name: root}
	}
	call.Args[argument] = resultPointer(pointer)
	return call
}

func resultPointer(pointer string) world.Binding {
	return world.Binding{ResultPointer: &pointer}
}

func suggest(call world.CallTemplate, why string, score float64) world.Suggestion {
	return world.Suggestion{Call: call, Why: why, Score: score}
}

func raw(value string) json.RawMessage {
	return json.RawMessage(value)
}
