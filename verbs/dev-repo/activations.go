package devrepo

import "github.com/nesste/phoenix/internal/world"

// AuthoredActivations is the initial intent-to-call mapping for the dev-repo
// world. Calls bind only roots already reachable in the live session.
func AuthoredActivations() []world.ActivationRule {
	return []world.ActivationRule{
		activation("diagnose_complete_suite", `(?i)(determine|diagnose|why).*(complete|suite).*(not green|fail|exact)`,
			suggest(rootCall("repo", "status"), "establish repository state before suite diagnosis", 1)),
		activation("run_complete_suite", `(?i)^\s*(run|execute).*(complete|entire|full).*(test|suite)`,
			suggest(rootCall("tests", "run"), "run the complete reachable test suite", 0)),
		activation("focus_named_test", `\b(Test[A-Za-z0-9_]+)\b`,
			suggest(capturedCall("tests", "focus", "test", 1), "run the test named in the intent", 1)),
		activation("discover_internal_service_test", `(?i)find.*internal.*test.*public\s+([A-Za-z0-9_-]+)\s+service`,
			suggest(capturedCall("repo", "find", "query", 1), "bridge the public service name through repository search", 1)),
		activation("list_contextual_test", `(?i)(first test|current [A-Za-z0-9_-]* ?test|whether .* test passes)`,
			suggest(rootCall("tests", "list"), "bind the requested test from the live suite", 0)),
		activation("build_repository", `(?i)\b(compile|compiles|builds)\b`,
			suggest(rootCall("repo", "build"), "check whether the reachable repository builds", 1)),
		activation("absent_deploy_or_release", `(?i)\b(deploy|deployment)\b|\brelease identifier\b`),
		activation("inspect_reachable_repository", `^$a`,
			suggest(rootCall("repo", "status"), "inspect reachable repository state", 0)),
	}
}

func activation(id, pattern string, suggestions ...world.Suggestion) world.ActivationRule {
	if suggestions == nil {
		suggestions = []world.Suggestion{}
	}
	return world.ActivationRule{ID: id, Pattern: pattern, Suggestions: suggestions}
}

func capturedCall(root, verb, argument string, capture int) world.CallTemplate {
	return world.CallTemplate{
		Handle: world.HandleSelector{Source: "root", Name: root}, Verb: verb,
		Args: map[string]world.Binding{argument: {IntentCapture: &capture}},
	}
}
