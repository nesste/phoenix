# 0022: Protocol-v4 Gate 1A validation reopening after world-compatibility repair

- **Status:** Validation open for one new disjoint execution; held-out closed
- **Date:** 2026-08-26
- **Decision:** `OPEN_GATE_1A`
- **Gate commit:** This decision record and the focused gate-state patch are committed atomically.
- **Prior opening:** `docs/decisions/0019-protocol-v4-gate-1a-validation-reopening-after-refactor.md`, withdrawn by decision 0020 before any model launch or validation outcome
- **Compatibility refreeze:** `docs/decisions/0021-protocol-v4-gate-1a-validation-world-compatibility-import-refreeze.md` at `370b52580b2bb038366d805264041243e3874da6`
- **Accepted execution payload:** `215b30ae89933c532468b452be6238a6028a740e`
- **Independent review:** `docs/reviews/2026-08-26-protocol-v4-validation-world-compatibility-review.md` at `a71e59b457cdfdb2014b434033946d319036c060`; verdict `ACCEPT`
- **Frozen validation schedule:** `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`
- **Frozen validation manifest:** `sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6`
- **Frozen execution boundary:** `sha256:57e071790f2e49e7a5ab5c45a2685578ccaf0f9bd4fbfea8712f7d172c9e6424` (21 files)
- **Frozen world-build digest:** `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4`

## Decision

The project chair reopens Gate 1A for one **new disjoint** execution of the frozen protocol-v4 validation schedule. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: true`, retains `may_open_held_out: false`, records this decision as the opening record with decision 0019 retained as the prior opening, and limits execution to `frozen_validation_schedule_only` with public validation inputs.

This decision authorizes one execution of the accepted 1,800-launch validation schedule under the precommitted A–E protocol, 300 USD run budget, numeric 0.15 USD per-trial cap, and 180-second timeout. It authorizes the runner to contact the designated external validation custodian through the reviewed interface. Private labels remain outside the Phoenix workspace and Git history.

The disjointness conditions remain strict:

- The 720-trial diagnostic archive at `experiments/frontier-v1/results/scheduled-validation` and its custody record remain closed evidence of decision 0013's `indeterminate` result. They must not be resumed, spliced, extended, or counted.
- The pre-execution attempt closed by decision 0020 reached the public-input compatibility guard before any model call, trial, grade, cost, output directory, or validation outcome. It spent no execution authorization.
- The new run must use `experiments/frontier-v1/results/scheduled-validation-2`, which is distinct from the archived path and absent at opening. Fresh evidence begins at launch index zero with a new `scheduled-checkpoint.json`.

This decision does not itself execute the schedule, create the output directory, authorize held-out access, permit a frozen-byte change, or authorize a retry outside the precommitted infrastructure-retry rule. It does not open Phase 2. Gate 1A results must be reported and independently reviewed before the project can continue past the stop gate.

## Basis

The chair verified the recorded preconditions for this reopening:

- decision 0020 closed Gate 1A immediately after the public-input compatibility guard rejected the former validation `world_ref`; no model or outcome was observed;
- the repair changes exactly 120 validation case `world_ref` values to the already-frozen production world and changes no case ID, goal, family, fixture, class, schedule entry, label digest, grader rule, prompt, runtime, arm behavior, world byte, world-build input, analysis byte, or held-out input;
- the replacement boundary was independently reviewed `ACCEPT` and frozen by decision 0021;
- the new gate-neutral compatibility guard preflights all 120 validation cases against the production world and rejects the former identity as a negative control;
- the complete repository quality gate passes, including all tests, formatting, vet, static analysis, module verification, vulnerability scan, complexity and duplication checks, corpus validation, and frozen world-build reproduction; and
- `pre-validation-artifacts.json` reports `status: complete` and `remaining: []`.

The matching external custodian adapter reports the replacement manifest identity and all unchanged custody identities. Its reviewed Linux/amd64 binary is `sha256:a4f0e8c0f3356ef2dff28e09319e0e68a5a6398d9eed45f9814527070eb5012f`; its launcher is `sha256:492ada05475aac47a809864f4d1379d48f7df17efd27aec9a3b944cef89e0237`. The adapter accepts only the new evidence directory and rejects the burned archive.

## Accepted residuals

The residuals recorded in earlier refreezes remain in force. In particular, the external custodian directory is not version-controlled, so its provenance depends on recorded hashes, preserved binaries, exact source reconstruction, and behavioral probes. Execution still requires a Linux/amd64 host. Any frozen-boundary change requires another independent replacement review and refreeze while the gate is closed.

## Gate patch boundary

This Gate 1A patch changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, opening validation, keeping held-out closed, recording this decision and the prior-opening pointer, and removing decision 0020's superseded closure fields;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, requiring the reopened-validation, closed-held-out state and this decision's provenance;
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the completed-freeze gate assertion; and
- this decision record.

The patch changes no accepted execution-boundary byte, schedule byte, manifest byte, protocol byte, prompt, schema, Arm B byte, world-build input, grader byte, analysis byte, corpus input, label registry, archived evidence byte, or private artifact.

## Verification

After the gate patch:

- the authoritative freeze guard reproduces every accepted artifact identity, including the 21-file boundary and world-build identity;
- the full test suite passes with validation open and held-out closed;
- `git diff --check` is clean; and
- no model, arm, schedule, trial, custodian call, private grade, validation result, held-out result, or outcome was run or observed while opening the gate.

## Authorized next action

External validation may run once on the Linux/amd64 execution host, using the frozen schedule, reviewed custodian adapter, and fresh output directory:

```text
D:\Work\personal\phoenix-validation-custodian\run-validation.sh
```

The underlying runner invocation remains the frozen 1,800-launch command with `--output-dir experiments/frontier-v1/results/scheduled-validation-2` and the 300 USD run budget. Any identity mismatch, budget mismatch, custodian failure, safety stop, or incomplete evidence must stop the run and produce an `indeterminate` result. Interrupted progress may resume only within this new run's own evidence directory under the frozen pairing-key checkpoint contract. Held-out remains closed.
