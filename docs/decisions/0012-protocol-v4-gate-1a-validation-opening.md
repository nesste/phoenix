# 0012: Protocol-v4 Gate 1A validation opening

- **Status:** Validation open; held-out closed
- **Date:** 2026-08-21
- **Decision:** `OPEN_GATE_1A`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.
- **Final boundary refreeze:** `docs/decisions/0011-protocol-v4-gate-1a-test-transition-import-refreeze.md` at `22ee54c`
- **Accepted execution payload:** `71f9648789decf4cd56ef8a24bc840b0dda7efd9`
- **Independent execution-boundary review:** `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-review.md` at `2feec72b66e262d2f43479d3037b2e3ebc8b0644`; verdict `ACCEPT`
- **Frozen validation schedule:** `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`
- **Frozen execution boundary:** `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`

## Decision

The project chair opens Gate 1A for the frozen protocol-v4 validation schedule. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: true`, retains `may_open_held_out: false`, and limits execution to `frozen_validation_schedule_only` with public validation inputs.

This decision authorizes one execution of the accepted 1,800-launch validation schedule under the precommitted A–E protocol, 300 USD run budget, numeric 0.15 USD per-trial cap, and 180-second timeout. It authorizes the runner to contact the designated validation custodian for private grading through the reviewed interface. Private labels must remain outside the Phoenix workspace and Git history.

This decision does not execute the schedule, select a custodian, create an output directory, authorize held-out access, permit a frozen-byte change, or authorize a retry outside the precommitted infrastructure-retry rule. It does not open Phase 2. Gate 1A results must be reported and independently reviewed before the project can continue past the stop gate.

## Basis

The chair verified the recorded pre-validation requirements:

- protocol v4 is independently accepted and frozen;
- the replacement public validation and held-out corpus package is accepted and imported without private labels;
- runtime invocation, prompts, Arm A schemas, Arm B, world definition and build, grader, analysis, report template, and validation schedule are frozen by digest;
- the validation schedule contains 120 public cases, five arms, three repetitions, and 1,800 launches;
- the custody-safe validation execution boundary and its gate-transition tests are independently accepted, imported, and refrozen; and
- `pre-validation-artifacts.json` reports `status: complete` and `remaining: []`.

The runner verifies the live validation schedule, public manifest, and label-digest registry before execution. It refuses held-out access. It rejects a changed trial cap or timeout before schedule preparation, output inspection or creation, Phoenix build, runtime verification, or custodian contact.

`protocol.json` is unchanged because it belongs to a previously accepted frozen identity. Its gate fields record the state at protocol-design acceptance. The runtime gate is the later chair-controlled `pre-validation-artifacts.json`, which the reviewed validation path reads directly.

## Accepted residuals

Decisions 0010 and 0011 record three non-blocking P2 findings. They remain in force:

- the trial-limit test's no-output and no-contact assertions are not independently load-bearing when it supplies no case IDs, although its error and ordering assertions are load-bearing and an independent populated probe established no contact;
- validation accepts numeric-equivalent Go spellings of `0.15` and forwards the supplied spelling to the external runtime, while the precommitted invocation uses the literal `0.15`; and
- flipping only the synthetic closed-gate boolean reaches the missing-schedule-identity error before the timeout, although the committed ordering assertion remains load-bearing and the populated open-gate probe reaches the timeout.

Any change that addresses these residuals by modifying a frozen boundary byte requires another independent replacement review and refreeze.

## Gate patch boundary

The Gate 1A patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, opening validation, keeping held-out closed, and recording the date, scope, and this decision path;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the open-validation, closed-held-out state and decision provenance; and
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the completed-freeze gate assertion; and
- this decision record.

The patch changes no accepted 16-file execution-boundary byte, schedule byte, protocol byte, prompt, schema, Arm B byte, world-build input, grader byte, analysis byte, corpus input, label registry, or private artifact.

## Verification

After the gate patch:

- the authoritative freeze guard reproduces every accepted artifact identity;
- the runner tests pass with validation open and held-out closed;
- `make quality` passes, including the reproducible world-build comparison at `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d`;
- `git diff --check` is clean; and
- no model, arm, schedule, trial, custodian, private grade, validation result, held-out result, or outcome was run or observed while opening the gate.

## Authorized next action

External validation may run once with the frozen schedule and invocation. Any identity mismatch, budget mismatch, custodian failure, safety stop, or incomplete evidence must stop the run and produce an `indeterminate` result under the frozen protocol. Held-out remains closed.
