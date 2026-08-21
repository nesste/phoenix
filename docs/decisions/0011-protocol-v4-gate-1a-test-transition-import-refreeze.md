# 0011: Protocol-v4 Gate 1A test-transition import and refreeze

- **Status:** Test transition imported and frozen; outcome gates remain closed
- **Date:** 2026-08-21
- **Decision:** `IMPORT_AND_FREEZE`
- **Independent review:** `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-review.md` at `2feec72b66e262d2f43479d3037b2e3ebc8b0644`; verdict `ACCEPT`; 23,202 raw bytes; `sha256:06789567c2aad2351d719285642e3d4eea258d14fe8b4e8f0f9668101766360e`
- **Accepted payload:** `71f9648789decf4cd56ef8a24bc840b0dda7efd9`
- **Payload import commit:** `37ff470`
- **Review-record import commit:** `1c8c61d`
- **Candidate report:** `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-candidate.md` at `7ffc4a1966b4127d9445ebc840723dc6a7a245d4`; 4,928 raw bytes; `sha256:689ad8cee3a6f5af2daed05081e10d1263b469086555d63757a7250a368dbc3f`
- **Focused refreeze commit:** `0c7453c`

## Decision

The project chair accepts the independent review, imports payload `71f9648`, and refreezes the 16-file validation execution boundary at `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`.

The candidate inventory is `experiments/frontier-v1/artifacts/gate-1a-test-transition-candidate.json`, with raw identity `sha256:ca7eae05c72f970a592e4a82ab88d2167564323f08a788a0081a349f32c49d75` over 4,118 bytes. The replacement supersedes accepted payload `74d06da13f62cffa4fd635e048e6331a6e2d95a6` and its 16-file identity `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1`.

`may_open_validation` remains `false`, and `may_open_held_out` remains `false`. This decision does not open Gate 1A, execute a model or arm, create a validation output, contact a real custodian, obtain a private grade, inspect a validation or held-out outcome, or authorize outcome analysis.

## Accepted replacement

The only changed file inside the 16-file execution boundary is `experiments/frontier-v1/runner/validation_execution_test.go`. Its LF-normalized identity is now `sha256:460dc582292473ecd5b91e436ba726e2f409399cd254588ee8f0eded99a02b03`.

The replacement moves two closed-gate assertions from the live implementation workspace to a synthetic temporary repository whose validation and held-out gates are explicitly false. The tests remain load-bearing when the chair later changes the live validation gate. Production runner code is byte-identical to the prior accepted boundary.

## Known non-blocking residuals

The chair accepts all three P2 findings recorded by the independent reviewer:

- **P2-1:** The trial-limit test's marker and output assertions remain incidentally satisfied because the test supplies no case IDs. Its error assertion and ordering mutation remain load-bearing, and an independent populated probe established no custodian contact.
- **P2-2:** Numeric-equivalent Go spellings of `0.15` are forwarded verbatim to the external runtime parser. The frozen operator spelling `0.15` is unaffected, and an external-parser rejection remains loud and auditable.
- **P2-3:** Flipping only the synthetic closed-gate boolean reaches the missing-schedule-identity error before the timeout. The committed ordering assertion still fails if the gate check moves behind the timeout check, and the independent populated open-gate probe exposed the timeout error.

These findings do not change the frozen operating path. Addressing any of them by changing a boundary byte requires another replacement review and refreeze.

## Import and refreeze boundary

Import commit `37ff470` applies the exact four-path payload:

- the replacement candidate inventory;
- the historical accepted-freeze guard;
- the live replacement-candidate guard; and
- `validation_execution_test.go`.

Review import `1c8c61d` adds only the independent review record. The candidate report and review prompt remain report-only files at `7ffc4a1`; the authoritative guard verifies the candidate report directly at that commit.

Focused refreeze commit `0c7453c` changes only `experiments/frontier-v1/pre-validation-artifacts.json` and `experiments/frontier-v1/runner/artifact_freeze_test.go`. It records the new candidate, report, review, findings, file identity, and aggregate identity, then switches the authoritative guard from historical verification of `74d06da…` to live verification of `71f9648…`. No accepted 16-file boundary byte changes in the focused refreeze.

## Verification

After import and refreeze:

- the runner suite passes with both gates closed;
- the authoritative guard reproduces the 16-file identity `sha256:9f8c48fb...de2d9` from the live files;
- the guard verifies the raw candidate, candidate-report, and independent-review identities;
- `make quality` passes, including root tests, vet, staticcheck, module verification, vulnerability scan, complexity and duplication checks, schema validation, authoring corpus validation, and reproducible world-build comparison;
- the world-build identity remains `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d`;
- `git diff --check` is clean; and
- no model, arm, schedule, trial, real custodian, private grade, validation result, held-out result, or outcome was run or observed.

## Remaining execution blocker

The refrozen boundary is ready for a separate project-chair Gate 1A decision. That final pre-execution decision must set `may_open_validation` true, leave `may_open_held_out` false, preserve every frozen identity, and authorize only the frozen validation schedule.
