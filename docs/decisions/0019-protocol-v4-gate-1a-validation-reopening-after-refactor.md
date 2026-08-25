# 0019: Protocol-v4 Gate 1A validation reopening after the complexity refactor

- **Status:** Validation open for one new disjoint execution; held-out closed
- **Date:** 2026-08-25
- **Decision:** `OPEN_GATE_1A`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.
- **Prior openings:** `docs/decisions/0012-protocol-v4-gate-1a-validation-opening.md` (spent by the interrupted execution closed in decision 0013) and `docs/decisions/0016-protocol-v4-gate-1a-validation-reopening.md` (withdrawn unspent by decision 0017 before any execution). Neither authorizes anything further.
- **Final boundary refreeze:** `docs/decisions/0018-protocol-v4-gate-1a-complexity-refactor-import-refreeze.md` at `91ab8f5`
- **Accepted execution payload:** `8ed9202c80d8f591c5d0db8a7e8a349022f952f8`
- **Independent complexity-refactor review:** `docs/reviews/2026-08-25-protocol-v4-gate-1a-complexity-refactor-review.md` at `8889539448dbb6eab88c0c5997f4f81ea17080e0`; verdict `ACCEPT`
- **Frozen validation schedule:** `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`
- **Frozen execution boundary:** `sha256:2d20be1ac41cc9072ae472b9b5a73e2607e9d35093fda49c422fa487117ecb4b` (20 files)
- **Frozen world-build digest:** `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4`

## Decision

The project chair reopens Gate 1A for one **new disjoint** execution of the frozen protocol-v4 validation schedule. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: true`, retains `may_open_held_out: false`, records this decision as the opening record with 0016 retained as the prior opening, and limits execution to `frozen_validation_schedule_only` with public validation inputs.

This decision authorizes one execution of the accepted 1,800-launch validation schedule under the precommitted A–E protocol, 300 USD run budget, numeric 0.15 USD per-trial cap, and 180-second timeout. It authorizes the runner to contact the designated validation custodian for private grading through the reviewed interface. Private labels must remain outside the Phoenix workspace and Git history.

The disjointness conditions of decision 0016 carry over unchanged:

- The 720-trial diagnostic archive at `experiments/frontier-v1/results/scheduled-validation` and its custody record remain closed evidence of decision 0013's `indeterminate` result. They must not be resumed, spliced, extended, or counted. The runner refuses that archive mechanically (missing checkpoint; checkpoint and world-build verification reject any recomposed content).
- The new run must use a fresh output directory, explicitly passed via `--output-dir`, distinct from the archived path. Fresh evidence starts at launch index zero with a new `scheduled-checkpoint.json`.

This decision does not execute the schedule, select a custodian, create an output directory, authorize held-out access, permit a frozen-byte change, or authorize a retry outside the precommitted infrastructure-retry rule. It does not open Phase 2. Gate 1A results must be reported and independently reviewed before the project can continue past the stop gate.

## Basis

The chair verified the recorded preconditions for this reopening:

- decision 0017 closed the gate before any execution under 0016, so no validation outcome exists and disjointness is preserved;
- the complexity refactor is independently reviewed `ACCEPT` with a line-by-line behavior-equivalence finding of no semantic drift, and is frozen in decision 0018 as the 20-file boundary `sha256:2d20be1a…cb4b`;
- live world-build matches the freeze: `TestValidationBuildReproducesFrozenWorldBuildDigest` rebuilds Phoenix with the unchanged frozen linux/amd64 `version=dev` recipe at `bin/phoenix` and reproduces `sha256:b5a26d5e…a2c4` from current sources, and the independent reviewer reproduced the same digest manually;
- the quality gate is fully green for the first time since decision 0014: tests, vet, staticcheck, `gocyclo -over 15`, and the duplication gate all pass;
- the frozen schedule, manifest, label-digest registry, grader, prompts, Arm B document, world definition, and analysis identities are unchanged since their acceptance; and
- `pre-validation-artifacts.json` reports `status: complete` and `remaining: []`.

The runner verifies the live validation schedule, public manifest, and label-digest registry before execution, refuses held-out access, rejects a changed trial cap or timeout before schedule preparation, output inspection or creation, Phoenix build, runtime verification, or custodian contact, and — after building with the freeze recipe and before any trial — stops unless the live world-build digest equals the frozen pin.

`protocol.json` is unchanged because it belongs to a previously accepted frozen identity. The runtime gate is the chair-controlled `pre-validation-artifacts.json`, which the reviewed validation path reads directly.

## Accepted residuals

The residuals recorded in decisions 0010, 0011, 0014, 0015, and 0018 remain in force, including the three P2 findings from the complexity-refactor review and the operational requirement that the built linux/amd64 binary can only execute on a linux/amd64 host. Any change that addresses these residuals by modifying a frozen boundary byte requires another independent replacement review and refreeze, with the gate closed for the duration per the 0017 precedent.

## Gate patch boundary

The Gate 1A patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, opening validation, keeping held-out closed, recording the date, scope, this decision path, and the prior-opening pointer, and leaving the spent-execution record from decision 0013 in place;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the reopened-validation, closed-held-out state with this decision's provenance;
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the completed-freeze gate assertion; and
- this decision record.

The patch changes no accepted 20-file execution-boundary byte, schedule byte, protocol byte, prompt, schema, Arm B byte, world-build input, grader byte, analysis byte, corpus input, label registry, archived evidence byte, or private artifact.

## Verification

After the gate patch:

- the authoritative freeze guard reproduces every accepted artifact identity, including the 20-file boundary `sha256:2d20be1a…cb4b` and world-build `sha256:b5a26d5e…a2c4`;
- the full test suite passes with validation open and held-out closed;
- `git diff --check` is clean; and
- no model, arm, schedule, trial, custodian, private grade, validation result, held-out result, or outcome was run or observed while opening the gate.

## Authorized next action

External validation may run once, on a linux/amd64 host, with the frozen schedule and invocation and a fresh output directory disjoint from the 720-trial archive:

```text
go run ./experiments/frontier-v1/runner --repo-root . --tranche validation --case all --schedule experiments/frontier-v1/schedules/validation.json --output-dir experiments/frontier-v1/results/scheduled-validation-2 --arm-b-document experiments/frontier-v1/arms/arm-b.md --validation-grader <absolute custodian grader path> --run-budget-usd 300
```

Any identity mismatch, budget mismatch, custodian failure, safety stop, or incomplete evidence must stop the run and produce an `indeterminate` result under the frozen protocol. Interrupted progress may resume only within the new run's own evidence directory under the frozen pairing-key checkpoint contract. Held-out remains closed.
