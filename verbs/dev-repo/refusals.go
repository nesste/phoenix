package devrepo

import (
	"encoding/json"

	"github.com/nesste/phoenix/internal/world"
)

type RefusalDefinition struct {
	HandleType string
	Verb       string
	Rules      []world.RefusalRule
}

// AuthoredRefusals converts recoverable starter-world constraints into
// in-band teaching results. Task 1.8 attaches these rules to the world verbs.
func AuthoredRefusals() []RefusalDefinition {
	return []RefusalDefinition{
		{
			HandleType: "tests", Verb: "focus",
			Rules: []world.RefusalRule{
				focusAlternative("missing_test", json.RawMessage(`{"properties":{"args":{"not":{"required":["test"]}}},"required":["args"]}`), "tests.focus requires a test name", "the requested test was not bound"),
				focusAlternative("unknown_test", failureCode("unknown_test"), "tests.focus declined the test name", "the requested test is not in the live suite"),
				focusAlternative("stale_suite", failureCode("stale_state"), "tests.focus declined stale suite state", "the test suite changed after the call was suggested"),
			},
		},
		{
			HandleType: "repo", Verb: "read",
			Rules: []world.RefusalRule{
				noAlternative("path_outside_repo", "repo.read declined the path", "the path is outside the reachable repository", "use a relative path inside the repository"),
			},
		},
		{
			HandleType: "repo", Verb: "edit",
			Rules: []world.RefusalRule{
				noAlternative("path_outside_repo", "repo.edit declined the path", "the path is outside the reachable repository", "use a relative path inside the repository"),
				noAlternative("edit_conflict", "repo.edit declined the replacement", "the expected text did not match exactly once", "read the current file and retry with an exact match"),
			},
		},
		{
			HandleType: "git", Verb: "commit",
			Rules: []world.RefusalRule{
				noAlternative("nothing_to_commit", "git.commit cannot run", "no changes are available to commit", "make a reachable edit before committing"),
			},
		},
		{
			HandleType: "episodes", Verb: "recall",
			Rules: []world.RefusalRule{
				noAlternative("recall_unavailable", "episodes.recall cannot run", "the episode store has not started", "start the episode store before recalling history"),
			},
		},
	}
}

func focusAlternative(id string, when json.RawMessage, what, why string) world.RefusalRule {
	return world.RefusalRule{
		ID: id, When: when, What: what, Why: why,
		Instead: &world.CallTemplate{
			Handle: world.HandleSelector{Source: "self"}, Verb: "list",
			Args: map[string]world.Binding{},
		},
	}
}

func noAlternative(code, what, why, reason string) world.RefusalRule {
	return world.RefusalRule{
		ID: code, When: failureCode(code), What: what, Why: why,
		Instead: nil, NoAlternative: reason,
	}
}

func failureCode(code string) json.RawMessage {
	encoded, _ := json.Marshal(map[string]any{
		"properties": map[string]any{
			"failure": map[string]any{
				"properties": map[string]any{"code": map[string]any{"const": code}},
				"required":   []string{"code"},
			},
		},
		"required": []string{"failure"},
	})
	return encoded
}
