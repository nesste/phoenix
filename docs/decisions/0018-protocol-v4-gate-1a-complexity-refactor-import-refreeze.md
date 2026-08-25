# 0018: Protocol-v4 Gate 1A complexity-refactor import and refreeze

- **Status:** Complexity refactor imported and frozen; outcome gates remain closed
- **Date:** 2026-08-25
- **Decision:** `IMPORT_AND_FREEZE`
- **Independent review:** `docs/reviews/2026-08-25-protocol-v4-gate-1a-complexity-refactor-review.md` at `8889539448dbb6eab88c0c5997f4f81ea17080e0`; verdict `ACCEPT`; 17,715 raw bytes; `sha256:e3055d077da82e856a4968e17a81013173ee0d6fb5854b75d7aad357a7c1a60d`
- **Accepted payload:** `8ed9202c80d8f591c5d0db8a7e8a349022f952f8`
- **Payload commit:** `8ed9202` (already on `main`)
- **Review-record import commit:** `8889539`
- **Focused refreeze commit:** `54712bdb40ddff937bcbe2459de744058d79877c`

## Decision

The project chair accepts the independent review, keeps payload `8ed9202` as the implementation, and refreezes the scheduled-runner and world-build identities that the refactor changed. This completes the cycle decision 0017 opened by closing the gate.

The replacement runner inventory is `experiments/frontier-v1/artifacts/gate-1a-complexity-refactor-candidate.json`, with raw identity `sha256:67ce91a3a1fb991c2fc2565f2d6ca35f0bd7a352daa40dce3ab629f9bea8ac34` over 5,155 bytes. The 20-file runner identity is `sha256:2d20be1ac41cc9072ae472b9b5a73e2607e9d35093fda49c422fa487117ecb4b`. It supersedes accepted payload `ed3860708c931ddc848b1bc90e4d6435585ff0d6` and its 20-file identity `sha256:c9e2f4ff1b6c37ea64b8d35b0124d432c1741be01788e51744cdfbd9b56f447f`.

Because `internal/activate` is compiled into `cmd/phoenix`, the frozen linux/amd64 world-build digest moves from `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e` to `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4`, with the regenerated manifest committed at `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` (LF identity `sha256:9210f42822f918b4fe1f7f73e7fafc760504f81cd012be09eb5eb2c472370cba`). The build recipe is unchanged. The world definition, canonical world identity, validation schedule, manifest, label registry, grader digest, prompts, and analysis artifacts are unchanged.

`may_open_validation` remains `false`, and `may_open_held_out` remains `false`. This decision does not open Gate 1A, execute a model or arm, create a validation output, resume or splice the 720-trial diagnostic archive, contact a real custodian, obtain a private grade, inspect a validation or held-out outcome, or authorize a new disjoint 1,800-launch run.

This chair session authored payload `8ed9202` and transcribed the independent review record. The independent review was performed by a separate reviewer session with a fresh context that did not author the payload; the same-machine / shared-workspace overlap recorded in the review record stands. The freeze relies on that `ACCEPT`, not on this session's self-tests.

## Accepted replacement

The runner freeze set stays at 20 files. Live hashes change for `scheduled_run.go`, `scheduled_resume.go`, `validation_build_recipe_test.go`, `README.md`, and `validation-execution-boundary.md`.

The replacement is a behavior-preserving complexity decomposition that brings `gocyclo -over 15` back under the quality gate: `runScheduledCases` (27 → 12) extracts `applyScheduledResume` and `runRemainingSchedule`; `loadScheduledResume` (24 → 11) extracts `classifyScheduledOutputEntries`; `activate.New` (16 → 9) extracts `validateActivationRule`; the gate-state freeze helper was restructured under decision 0017. The independent review verified line-by-line that no error message, error ordering, stop rule, checkpoint rule, refusal condition, or evidence byte changed. Duplicate candidate-guard tests were consolidated into shared helpers; the superseded build-recipe candidate guard is now a historical check at `ed38607`, and new working-tree guards live outside the frozen set.

## Known non-blocking residuals

The chair records the independent review's `ACCEPT` with no P0 or P1 findings. The following residuals do not block this freeze:

- P2-1: `classifyScheduledOutputEntries` measures gocyclo 15, exactly at the threshold with zero headroom;
- P2-2: the candidate note's phrasing implies three expected interim failures where exactly two tests fail;
- P2-3: `applyScheduledResume` returns `done=true` with a non-nil error on stop paths, so the flag is meaningful only when the error is nil; the sole caller checks the error first;
- a validation run still requires a linux/amd64 execution host.

Earlier accepted residuals from decisions 0010, 0011, 0014, and 0015 remain in force where their subjects are unchanged. The `gocyclo -over 15` residual recorded in 0014 and 0015 is discharged: the quality gate's complexity check now passes across its full scope.

## Import and refreeze boundary

Payload `8ed9202` already landed the implementation. Review import `8889539` adds only the independent review record.

Focused refreeze commit `54712bd` changes:

- `experiments/frontier-v1/pre-validation-artifacts.json` (scheduled_runner block: new 20-file inventory and review provenance; world-build block: new manifest and digest identities; local artifact blocks re-attributed to the accepted payload and review);
- `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` (regenerated from the accepted sources under the unchanged frozen recipe);
- `experiments/frontier-v1/artifacts/pre-validation-local-candidate.json` (world-build manifest and digest identities);
- `experiments/frontier-v1/runner/artifact_freeze_test.go`.

The superseded build-recipe candidate remains in the tree and is verified against payload `ed38607`, not the working tree. No schedule, grader, prompt, world-definition, corpus, or archive byte is changed in this freeze. Both gates stay false.

## Verification

After import and refreeze:

- the full test suite passes with both gates closed (`go test -count=1 ./...` exit 0), including the replacement freeze guard, the new candidate guard, all historical guards, and the live world-build reproduction of `sha256:b5a26d5e…a2c4`;
- `gocyclo -over 15` and the duplication gate pass across the full quality-gate scope — the purpose of this cycle;
- `gofmt`, `go vet`, and staticcheck v0.7.0 are clean;
- `git diff --check` is clean;
- no model, arm, schedule, trial, custodian, private grade, validation result, held-out result, or outcome was run or observed.

## Remaining execution blocker

The refactored boundary is frozen and closed. A later project-chair decision may reopen validation for one new disjoint 1,800-launch run only under the conditions decisions 0016 and 0017 preserved: `may_open_held_out` stays false, every frozen identity is preserved as refrozen here, the 720-trial diagnostic archive is refused, and execution occurs on a linux/amd64 host. Decisions 0012 and 0016 remain spent and withdrawn respectively.
