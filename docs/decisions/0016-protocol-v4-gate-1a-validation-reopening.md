# 0016: Protocol-v4 Gate 1A validation reopening

- **Status:** Validation open for one new disjoint execution; held-out closed
- **Date:** 2026-08-25
- **Decision:** `OPEN_GATE_1A`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.
- **Prior opening:** `docs/decisions/0012-protocol-v4-gate-1a-validation-opening.md` — spent by the interrupted execution closed in `docs/decisions/0013-protocol-v4-gate-1a-interrupted-execution.md`; it authorizes nothing further.
- **Final boundary refreeze:** `docs/decisions/0015-protocol-v4-gate-1a-validation-build-recipe-import-refreeze.md` at `44fc5991`
- **Accepted execution payload:** `ed3860708c931ddc848b1bc90e4d6435585ff0d6`
- **Independent build-recipe review:** `docs/reviews/2026-08-25-protocol-v4-gate-1a-validation-build-recipe-review.md` at `8dd34b5eb013dcbe46243ce7768ff5a60b93ed66`; verdict `ACCEPT`
- **Frozen validation schedule:** `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`
- **Frozen execution boundary:** `sha256:c9e2f4ff1b6c37ea64b8d35b0124d432c1741be01788e51744cdfbd9b56f447f` (20 files)
- **Frozen world-build digest:** `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`

## Decision

The project chair reopens Gate 1A for one **new disjoint** execution of the frozen protocol-v4 validation schedule. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: true`, retains `may_open_held_out: false`, records this decision as the opening record with 0012 retained as the prior opening, and limits execution to `frozen_validation_schedule_only` with public validation inputs.

This decision authorizes one execution of the accepted 1,800-launch validation schedule under the precommitted A–E protocol, 300 USD run budget, numeric 0.15 USD per-trial cap, and 180-second timeout. It authorizes the runner to contact the designated validation custodian for private grading through the reviewed interface. Private labels must remain outside the Phoenix workspace and Git history.

The authorized run is disjoint from the spent 2026-08-23 execution:

- The 720-trial diagnostic archive at `experiments/frontier-v1/results/scheduled-validation` and its custody record remain closed evidence of decision 0013's `indeterminate` result. They must not be resumed, spliced, extended, or counted.
- The runner refuses that archive mechanically: it predates the frozen checkpoint contract, so `loadScheduledResume` rejects the directory ("scheduled resume is missing a checkpoint"), and any recomposed content fails the checkpoint and world-build verification added in decision 0014.
- The new run must therefore use a fresh output directory, explicitly passed via `--output-dir`, distinct from the archived path. Fresh evidence starts at launch index zero with a new `scheduled-checkpoint.json`.

This decision does not execute the schedule, select a custodian, create an output directory, authorize held-out access, permit a frozen-byte change, or authorize a retry outside the precommitted infrastructure-retry rule. It does not open Phase 2. Gate 1A results must be reported and independently reviewed before the project can continue past the stop gate.

## Basis

The chair verified the recorded preconditions for a reopening:

- decision 0013 closed the interrupted execution as `indeterminate` and authorized authoring repairs only;
- the three 0013 repairs (orientation handoff, world-build pin, pairing-key checkpoint/resume) are independently reviewed `ACCEPT` and frozen in decision 0014;
- the remaining 0014 residual — validation builds used the host recipe, so the world-build pin could never match — is repaired by payload `ed38607`, independently reviewed `ACCEPT`, and frozen in decision 0015 as the 20-file boundary `sha256:c9e2f4ff…447f`;
- live world-build matches the freeze: `TestValidationBuildReproducesFrozenWorldBuildDigest` rebuilds Phoenix with the frozen linux/amd64 `version=dev` recipe at `bin/phoenix` and reproduces `sha256:bf976cad…8a4e` from current sources;
- the frozen schedule, manifest, label-digest registry, grader, prompts, Arm B document, world definition, and analysis identities are unchanged since their acceptance; and
- `pre-validation-artifacts.json` reports `status: complete` and `remaining: []`.

The runner verifies the live validation schedule, public manifest, and label-digest registry before execution, refuses held-out access, rejects a changed trial cap or timeout before schedule preparation, output inspection or creation, Phoenix build, runtime verification, or custodian contact, and — after building with the freeze recipe and before any trial — stops unless the live world-build digest equals the frozen pin.

`protocol.json` is unchanged because it belongs to a previously accepted frozen identity. Its gate fields record the state at protocol-design acceptance. The runtime gate is the later chair-controlled `pre-validation-artifacts.json`, which the reviewed validation path reads directly.

## Accepted residuals

The residuals recorded in decisions 0010, 0011, 0014, and 0015 remain in force, including: the P2 cap, timeout, and closed-gate ordering findings; the P2 `bin/phoenix` overwrite-and-remove behavior and substring-based target-pin assertions from the 0015 review; the four pre-existing `gocyclo -over 15` functions; and the operational requirement that the built linux/amd64 binary can only execute on a linux/amd64 host. Any change that addresses these residuals by modifying a frozen boundary byte requires another independent replacement review and refreeze.

## Gate patch boundary

The Gate 1A patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, opening validation, keeping held-out closed, recording the date, scope, this decision path, and the prior-opening pointer, and leaving the spent-execution record from decision 0013 in place;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the reopened-validation, closed-held-out state with this decision's provenance and the unchanged spent-execution history; and
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the completed-freeze gate assertion; and
- this decision record.

The patch changes no accepted 20-file execution-boundary byte, schedule byte, protocol byte, prompt, schema, Arm B byte, world-build input, grader byte, analysis byte, corpus input, label registry, archived evidence byte, or private artifact.

## Verification

After the gate patch:

- the authoritative freeze guard reproduces every accepted artifact identity, including the 20-file boundary `sha256:c9e2f4ff…447f`;
- the full test suite passes with validation open and held-out closed, including the live world-build reproduction of `sha256:bf976cad…8a4e`;
- `git diff --check` is clean; and
- no model, arm, schedule, trial, custodian, private grade, validation result, held-out result, or outcome was run or observed while opening the gate.

## Authorized next action

External validation may run once, on a linux/amd64 host, with the frozen schedule and invocation and a fresh output directory disjoint from the 720-trial archive:

```text
go run ./experiments/frontier-v1/runner --repo-root . --tranche validation --case all --schedule experiments/frontier-v1/schedules/validation.json --output-dir experiments/frontier-v1/results/scheduled-validation-2 --arm-b-document experiments/frontier-v1/arms/arm-b.md --validation-grader <absolute custodian grader path> --run-budget-usd 300
```

Any identity mismatch, budget mismatch, custodian failure, safety stop, or incomplete evidence must stop the run and produce an `indeterminate` result under the frozen protocol. Interrupted progress may resume only within the new run's own evidence directory under the frozen pairing-key checkpoint contract. Held-out remains closed.
