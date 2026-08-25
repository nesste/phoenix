# Authoring repair: complexity refactor under the quality gate

Status: implemented in authoring sources; not frozen; not authorized for validation.
Owner: authoring repair after decision 0017, which closed the reopened gate for this cycle.
Frozen validation bytes are unchanged in `pre-validation-artifacts.json`. Freeze tests for the scheduled runner and the world-build reproduction test are expected to fail until independent review and refreeze.

## Residual

Decisions 0014 and 0015 recorded that `gocyclo -over 15` (the `make complexity` quality gate) fails on four functions: `runScheduledCases` (27), `loadScheduledResume` (24), the gate-state freeze helper in `artifact_freeze_test.go` (16+), and `New` in `internal/activate` (16). The first two are frozen runner boundary bytes; `activate.New` is compiled into the Phoenix executable, so any change to it changes the world-build digest.

## Implemented authoring behavior

All refactors are behavior-preserving decompositions; no error message, ordering, stop rule, or evidence byte changes:

- `runScheduledCases` (27 → 12) extracts `applyScheduledResume` (replaying retained records and terminal-stop handling) and `runRemainingSchedule` (the pairing-key group loop with budget/safety stops and checkpoints).
- `loadScheduledResume` (24 → 11) extracts `classifyScheduledOutputEntries` (directory-entry classification, orphan-evidence and missing-checkpoint refusals).
- `activate.New` (16 → 9) extracts `validateActivationRule` (suggestion-count, fallback, why-line, and score validation).
- The gate-state freeze helper was restructured table-driven under decision 0017's closure patch.

Because `internal/activate` is linked into `cmd/phoenix`, the frozen linux/amd64 world-build digest changes from `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e` to `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4`. The live README and boundary documents now name the new digest; the freeze documents still name the old one until refreeze, which is why the world-build reproduction test fails in the interim. The build recipe itself is unchanged.

The working-tree guard for the superseded build-recipe candidate is converted to a historical check at its accepted payload `ed3860708c931ddc848b1bc90e4d6435585ff0d6`; the new candidate's working-tree guard lives in `validation_candidate_artifact_test.go`, outside the frozen set, so future cycles do not need to edit frozen test bytes to retire a guard.

Public-only verification: the full runner suite passes except the two freeze guards and the world-build reproduction test named above; `gocyclo -over 15` is clean across the quality-gate scope; no model, arm, trial, custodian, gate state, or grade is touched.

## Out of scope here

This note does not update freeze hashes, reopen validation, change any stop rule or evidence format, resume or splice the 720-trial diagnostic archive, or change the build recipe. Reopening requires the refreeze and a distinct chair decision.
