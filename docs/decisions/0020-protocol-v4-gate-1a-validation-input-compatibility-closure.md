# 0020: Protocol-v4 Gate 1A closure for validation-input compatibility repair

- **Status:** Validation closed; held-out closed; execution stopped during public preflight
- **Date:** 2026-08-26
- **Decision:** `CLOSE_GATE_1A`
- **Closed opening:** `docs/decisions/0019-protocol-v4-gate-1a-validation-reopening-after-refactor.md`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.

## Decision

The project chair closes Gate 1A and withdraws decision 0019 before any retry. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: false`, retains `may_open_held_out: false`, and records this closure without changing the prior interrupted-execution evidence from decision 0013.

At 2026-08-26 09:05 Europe/Bucharest, the reviewed validation-2 launcher was invoked once. The runner stopped during public case preflight because validation case `validation_01afd934` pins world reference `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b`, while the production world frozen after the authoring repair has canonical identity `sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc`.

The stop occurred before output-directory creation, Phoenix build, runtime verification, model invocation, trial creation, grading, checkpoint creation, or cost. `scheduled-validation-2` remained absent. The launcher produced an empty stdout log and a 225-byte stderr log with raw SHA-256 `99ac550ff46ce3f0b671359d215a8c57be9f33bb43b7aa92c0f13d13fce0f1a1`. This was a failed public preflight, not a validation outcome. No validation data was opened or burned.

The mismatch is systematic: all 120 frozen validation case files and the validation manifest pin the former world identity. The authoring repair imported by decision 0014 changed and refroze the production world but intentionally left the validation schedule, manifest, and sealed corpus bytes unchanged. The incompatibility was not exercised by the closed-gate tests or caught by the later reviews.

## Required repair

Before Gate 1A may reopen:

1. Migrate the public validation case `world_ref` values to the already frozen production-world identity without changing case IDs, fixtures, labels, grading scripts, schedule order, arms, repetitions, or outcomes.
2. Regenerate the public validation manifest case-input identities and its canonical digest.
3. Update the validation boundary and external custodian description to the new public manifest identity without changing the frozen grader semantics, label registry, private archive, or schedule identity.
4. Add a no-model test that checks every scheduled validation case against the frozen production world while the gate is closed.
5. Obtain an independent review of the exact replacement, refreeze the changed identities, and reopen validation in a later chair decision.

## Gate patch boundary

This closure patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, closing validation and recording this decision;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the reclosed state;
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, requiring the reclosed state; and
- this decision record.

No validation case, manifest, schedule, world, runner execution path, custodian, private label, prompt, analysis, or archived result changes in this closure commit. The compatibility repair follows as a separate replacement candidate.

## Authorization boundary

Decision 0019 authorizes nothing further. A later decision may authorize one new disjoint 1,800-launch validation execution only after the compatibility repair is independently accepted and refrozen. Held-out remains closed.
