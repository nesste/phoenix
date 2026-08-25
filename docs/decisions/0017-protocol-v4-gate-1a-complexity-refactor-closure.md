# 0017: Protocol-v4 Gate 1A closure for the complexity refactor

- **Status:** Validation closed; held-out closed; no execution occurred under decision 0016
- **Date:** 2026-08-25
- **Decision:** `CLOSE_GATE_1A`
- **Closed opening:** `docs/decisions/0016-protocol-v4-gate-1a-validation-reopening.md`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.

## Decision

The project chair closes Gate 1A before any execution under decision 0016. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: false`, retains `may_open_held_out: false`, and records this closure alongside the unchanged 0016 opening provenance and the spent 0013 execution record.

The reason for closure is an intended frozen-byte change: bringing `gocyclo -over 15` back under the quality gate requires refactoring `runScheduledCases` and `loadScheduledResume` (frozen runner boundary bytes) and `New` in `internal/activate` (compiled into the Phoenix executable, so the frozen world-build digest changes). The frozen validation execution boundary requires that replacement bytes be independently reviewed and refrozen before any opening; the chair therefore closes the gate for the duration of the authoring, review, and refreeze cycle rather than refreezing under an open gate.

No launch, trial, model call, custodian contact, output directory, or grade was produced under decision 0016 between its opening and this closure. Decision 0016's authorization is withdrawn unspent; it authorizes nothing further. A later chair decision may reopen validation only after the replacement bytes are independently accepted and refrozen, under the same conditions 0016 imposed: one new disjoint 1,800-launch run, held-out closed, frozen identities preserved, and the 720-trial diagnostic archive refused.

## Gate patch boundary

The closure patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, closing validation and recording the closure date and this decision path;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the reclosed state (this update also restructures the gate-state helper below the complexity threshold);
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the completed-freeze gate assertion; and
- this decision record.

No frozen boundary byte, schedule byte, world byte, world-build input, grader byte, prompt, or archive byte changes in this patch. The frozen-byte refactor itself follows as a separate replacement candidate.

## Verification

After the closure patch, the full test suite passes with both gates closed, and no model, arm, trial, custodian, private grade, or outcome was run or observed.
