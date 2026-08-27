package surface

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nesste/phoenix/internal/teach"
)

// TestWireRoundTripIsLosslessForEveryResponseShape is the frozen §4 invariant:
// Expand(Compress(x)) == x for every response shape, with expansion
// deterministic and lossless, whether the response announces the full
// world_build digest or only its prefix.
func TestWireRoundTripIsLosslessForEveryResponseShape(t *testing.T) {
	state := testWorldBuild
	shapes := map[string]Envelope{
		"ok_with_frontier": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_01234567",
			Handle: "h_0123456789abcdef", Verb: "run", Status: StatusOK,
			Result: map[string]any{"passed": float64(12)}, Text: "ok tests.run",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}},
			Frontier: []FrontierEntry{{
				Call: Call{Handle: "h_fedcba9876543210", Verb: "status", Args: map[string]any{}, State: &state},
				Why:  "inspect after verification", Provenance: "authored", Score: 1,
			}},
		},
		"ok_orient": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_11234567",
			Handle: "h_0123456789abcdef", Verb: "orient", Status: StatusOK,
			Result: map[string]any{"matched": false}, Text: unmatchedIntentText,
			Handles:  HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}},
			Frontier: []FrontierEntry{},
		},
		"fail_with_declared_args": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_21234567",
			Handle: "h_0123456789abcdef", Verb: "focus", Status: StatusFail,
			Error: &Error{
				Code: "invalid_arguments", Message: "verb arguments do not match the declared schema",
				DeclaredArgs: json.RawMessage(`{"type":"object","properties":{"test":{"type":"string"}}}`),
			},
			Text:    "failed focus: verb arguments do not match the declared schema",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"fail_with_details": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_31234567",
			Handle: "h_0123456789abcdef", Verb: "build", Status: StatusFail,
			Error:   &Error{Code: "execution_failed", Message: "build exited unsuccessfully", Details: map[string]any{"exit_code": float64(1)}},
			Text:    "failed build: build exited unsuccessfully",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"refused_with_instead": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_41234567",
			Handle: "h_0123456789abcdef", Verb: "focus", Status: StatusRefused,
			Refusal: &Refusal{
				What: "tests.focus", Why: "test is required",
				Instead: &Call{Handle: "h_0123456789abcdef", Verb: "list", Args: map[string]any{}},
			},
			Text:    "refused tests.focus: test is required",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"refused_no_alternative": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_51234567",
			Handle: "h_0123456789abcdef", Verb: "commit", Status: StatusRefused,
			Refusal: &teach.Refusal{What: "commit cannot run", Why: "no changes", NoAlternative: "edit first"},
			Text:    "refused: no changes",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"absent_with_reachable_verbs": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_61234567",
			Handle: "h_0123456789abcdef", Verb: "deploy", Status: StatusAbsent,
			Error:   &Error{Code: "absent", Message: "act is not reachable", ReachableVerbs: []string{"edit", "find", "list", "read"}},
			Text:    "absent: act is not reachable",
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"exhausted": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_71234567",
			Handle: "h_0123456789abcdef", Verb: "orient", Status: StatusExhausted,
			Text:    exhaustedOrientationText,
			Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
		},
		"handle_delta": {
			V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_81234567",
			Handle: "h_0123456789abcdef", Verb: "open", Status: StatusOK,
			Result: map[string]any{"opened": true}, Text: "ok repo.open",
			Handles: HandleDelta{
				Grant:  []HandleGrant{{Ref: "h_fedcba9876543210", Type: "repo", Label: "granted", State: testWorldBuild}},
				Revoke: []string{"h_aaaaaaaaaaaaaaaa"},
			},
			Frontier: []FrontierEntry{},
		},
	}
	for name, envelope := range shapes {
		for _, first := range []bool{true, false} {
			wire := Compress(envelope, first)
			serialized := marshalJSON(t, wire)
			validateWireJSON(t, serialized)
			var decoded WireEnvelope
			if err := json.Unmarshal(serialized, &decoded); err != nil {
				t.Fatalf("%s decode: %v", name, err)
			}
			roundTripped, err := Expand(decoded, testWorldBuild)
			if err != nil {
				t.Fatalf("%s expand: %v", name, err)
			}
			normalizedWant := normalizeEnvelope(t, envelope)
			normalizedGot := normalizeEnvelope(t, roundTripped)
			if !reflect.DeepEqual(normalizedGot, normalizedWant) {
				t.Fatalf("%s (first=%v) round trip diverged:\n got %s\nwant %s", name, first, normalizedGot, normalizedWant)
			}
		}
	}
}

// normalizeEnvelope compares envelopes through their canonical archival JSON
// so pointer identity and raw-message spacing cannot mask or fake divergence.
func normalizeEnvelope(t *testing.T, envelope Envelope) string {
	t.Helper()
	return string(marshalJSON(t, envelope))
}

func TestExpandRejectsMismatchedOrMissingWorldBuild(t *testing.T) {
	envelope := Envelope{
		V: 1, WorldBuild: testWorldBuild, SessionID: "s_0123456789abcdef", ActID: "a_91234567",
		Handle: "h_0123456789abcdef", Verb: "run", Status: StatusOK,
		Result: map[string]any{}, Text: "ok",
		Handles: HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}}, Frontier: []FrontierEntry{},
	}
	wire := Compress(envelope, false)
	other := "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, err := Expand(wire, other); err == nil {
		t.Fatal("expansion accepted a mismatched world-build prefix")
	}
	wire.WorldBuildPrefix = ""
	if _, err := Expand(wire, testWorldBuild); err == nil {
		t.Fatal("expansion accepted a wire envelope with no world-build identity")
	}
}

func TestFirstWireResponseAnnouncesFullWorldBuildOnce(t *testing.T) {
	session := &Session{}
	if !session.firstWireResponse() {
		t.Fatal("first response did not announce the world build")
	}
	if session.firstWireResponse() {
		t.Fatal("second response announced the world build again")
	}
}
