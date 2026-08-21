# 0010: Protocol-v4 validation execution-boundary import and refreeze

- **Status:** Replacement boundary imported and frozen; outcome gates remain closed
- **Date:** 2026-08-21
- **Decision:** `IMPORT_AND_FREEZE`
- **Independent review:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-review.md` at `6b158ea0de163b43f8c0fdeb1e8cc3fd40609a52`; verdict `ACCEPT`; 26,902 raw bytes; `sha256:e1444fcca48e7d5b13e1b0e4cd385898d387d8aba4799fae0e3d49f27326ddaf`
- **Accepted payload:** `74d06da13f62cffa4fd635e048e6331a6e2d95a6`
- **Payload import commit:** `f0326ab`
- **Candidate report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-candidate.md` at `05ab169acdfabae4ce7b42b4ffb8d6ecb276e935`; 5,235 raw bytes; `sha256:c77207aaebb6876805354e720adc7e72301ab759e2affaa3df3bb7c88526c546`
- **Focused refreeze commit:** `e7ca24d`

## Decision

The project chair accepts the independent second-revision review, imports the exact accepted payload onto `main`, and refreezes the reviewed 16-file validation execution boundary in `experiments/frontier-v1/pre-validation-artifacts.json`.

The frozen boundary identity is `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1`. Its candidate inventory is `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json`, with raw identity `sha256:14b0761362709780ded9f6fc47e7ff8d49f688e09062d8dcc6d13724b9109d1d` over 4,244 bytes. The accepted payload replaces scheduled-runner candidate `b4df919070bb9a6d2912662b4a59674b0e25a332`.

The manifest source limit is now `public_validation_inputs_only`. `may_open_validation` remains `false`, and `may_open_held_out` remains `false`. This decision imports and freezes an outcome-free execution path; it does not open Gate 1A, execute a model or arm, contact the real custodian, obtain a private grade, create a validation outcome, open held-out, or authorize outcome analysis.

## Accepted behavior

For validation execution, the closed Gate 1A check runs first. Once that gate is separately opened, the runner rejects a cap that does not parse to numeric `0.15`, a timeout that is not exactly `180s` as a `time.Duration`, or a run budget that is not exactly 300 USD. These checks run before schedule preparation, output-directory inspection or creation, Phoenix build, runtime verification, and custodian path resolution or handshake. Authoring behavior remains configurable under its existing contract.

The boundary also fixes the frozen public validation schedule and world-build identity, permits only public validation cases, obtains private grading through the synthetic or external custodian interface without importing private labels, and keeps held-out inaccessible.

## Known non-blocking residuals

The independent reviewer recorded two P2 findings for the chair. Both are accepted as non-blocking:

- **P2-1:** The candidate report overstates the strength of the marker and output assertions in `TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact`. Because that test supplies no case IDs, schedule validation also prevents the later side effects. Its error assertion and guard-order mutation remain load-bearing, and the reviewer's independent populated probe established that the corrected guard prevents custodian contact.
- **P2-2:** Numeric-equivalent Go spellings of `0.15` are forwarded verbatim to the external runtime parser. The frozen operator command uses the ordinary `0.15` spelling, so the reviewed operating path is unaffected. Any external-parser rejection is loud and enters the existing safety-stop and `indeterminate` evidence paths. A later replacement may narrow validation to the literal string `0.15` or forward a canonical spelling, but doing so would require a new review and refreeze.

## Focused refreeze boundary

Import commit `f0326ab` applies the exact accepted payload. Focused refreeze commit `e7ca24d` then changes only:

- `experiments/frontier-v1/pre-validation-artifacts.json`, replacing the prior nine-file scheduled-runner entry with the accepted 16-file boundary, adding candidate/report/review provenance and the two accepted findings, and setting the source limit to public validation inputs; and
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, verifying the exact candidate, report, review, replacement, file-set, source-limit, and closed-gate identities.

The focused refreeze changes no file inside the accepted 16-file boundary. The candidate report remains a report-only commit rather than an imported implementation artifact; the freeze guard verifies its raw identity directly at its recorded commit.

## Verification

After the import and focused refreeze:

- the runner test suite passes, including the candidate-artifact and authoritative-freeze guards;
- the freeze guard reproduces the 16-file identity `sha256:3056e95a...83f1` from the live accepted files;
- the guard verifies the raw candidate, report, and independent-review identities;
- the artifact manifest reports `complete`, zero remaining artifact-freeze items, a public-validation-only source limit, and both gates false;
- `make quality` passes, including root tests, vet, staticcheck, module verification, vulnerability scan, complexity and duplication checks, schema validation, authoring corpus validation, and reproducible world-build comparison;
- `git diff --check` is clean; and
- no model, arm, trial, real custodian, private grade, validation result, held-out result, or outcome was produced or observed.

## Remaining execution blocker

The validation execution boundary is ready for a separate project-chair Gate 1A decision. That is the final decision step before external validation can begin. It must explicitly open validation while leaving held-out closed and must not alter any frozen boundary byte.
