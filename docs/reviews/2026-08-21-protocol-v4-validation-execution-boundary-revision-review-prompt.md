# Protocol-v4 validation execution-boundary revision: independent review assignment

- **Assignment owner:** project-chair workflow
- **Assigned reviewer role:** `validation_execution_boundary.revision_independent_reviewer`
- **Date:** 2026-08-21
- **Frozen implementation base:** `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`
- **Rejected payload:** `9de7861d473b9ec14778638938944b75aeff937c`
- **Reject review:** `6f1696542c4033fe3edd660c95b07b2e74f292f5`; 15,808 raw bytes; `sha256:22bf9ce094021c709d9005f90f4a486167a5bf5aa4a2dd65017ad300210bc7f0`
- **Replacement branch:** `codex/validation-execution-candidate-revision`
- **Replacement payload:** `4bba4b4f988e2f9bbc4dbc986fbe7413078b9989`
- **Replacement report commit:** `7190e8cb4b304d437af05ecdddd9082537bed036`
- **Replacement report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-revision-candidate.md`
- **Replacement report identity:** 5,478 raw bytes; `sha256:19b073280e3f1770680b109ff389307c9ef3651ff4b27a71cd814e87c1dae0c2`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json` at the replacement payload
- **Candidate inventory identity:** 4,205 raw bytes; `sha256:680fc24683f5d8dedf1c74634026d94da40b703591678a6fcc3b214b9a690d8d`
- **Boundary artifact-set identity:** `sha256:012784de3eb93f6c84f7b93008bd7f8e6834b2d53427c043661c0846e670ea1e`

Review the replacement candidate for `ACCEPT`, `REVISE`, or `REJECT`. This assignment authorizes independent public-only review. It does not authorize import, refreeze, Gate 1A, a model or arm run, validation execution, a real custodian handshake, a private grade, outcome access, held-out access, or analysis.

## Independence and checkout

The reviewer must be distinct from:

- the original candidate author and reviewer;
- the revision author and this assignment author;
- any later import/refreeze chair or Gate 1A chair.

Declare any same-machine or shared-workspace limitation. Create a clean disposable worktree or detached checkout at report commit `7190e8cb4b304d437af05ecdddd9082537bed036`. Do not modify the replacement branch. Record the checkout path and require `git status --porcelain` to be empty before and after the review.

Do not open, mount, enumerate, hash, copy, or inspect a private evaluator root or any validation/held-out label file. Do not substitute a real custodian executable. Synthetic, outcome-free helpers in temporary directories are permitted. Do not execute a schedule or invoke a model, runtime, Phoenix arm, grader on model output, or outcome analysis.

Read the reject review from commit `6f1696542c4033fe3edd660c95b07b2e74f292f5` and verify its raw identity before relying on it. Treat the replacement as a complete candidate from the frozen base, not as a patch to apply over the rejected payload.

## Identity and ancestry

Verify independently:

1. payload `4bba4b4f988e2f9bbc4dbc986fbe7413078b9989` has sole parent `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`;
2. report commit `7190e8cb4b304d437af05ecdddd9082537bed036` has sole parent `4bba4b4…` and adds only the named revision report;
3. the payload changes exactly the 16 paths named by the report and no private, result, runtime-stream, trial, grade, analysis-output, or gate-state path;
4. the inventory is raw-byte identical to the identity above, contains 16 boundary files, names the rejected candidate and finding, and records `run_budget_usd: 300`, `max_cost_usd_per_trial: 0.15`, and `timeout_seconds: 180`;
5. inventory state remains `review_candidate`, `frozen: false`, both gates false, and `execution_performed: false`;
6. every LF-normalized file digest and the ordinal-path aggregate digest reproduce exactly.

The report-only commit is context and claimed verification evidence. It is not part of the implementation payload.

## Remediation finding

The rejected candidate passed the prior identity, gate, schedule, custody, evidence, authoring, and held-out audits but failed to lock two frozen trial limits. Determine whether the replacement fully closes that P1:

- validation accepts a per-trial cost cap only when its numeric value is exactly `0.15` USD;
- validation accepts a timeout only when its duration is exactly `180s`;
- a changed cap or timeout is rejected before schedule preparation, output-directory inspection or creation, Phoenix build, runtime verification, custodian path resolution, or custodian `describe` contact;
- the existing closed-gate check still runs first;
- the exact-limit checks apply only to validation and do not change authoring configuration behavior;
- the validation run budget remains exactly 300 USD;
- the accepted values reach the existing validation preparation path without changing their representation in summary evidence or pairing-budget calculations.

Inspect the call order in `prepareScheduledCLI` and the implementation of `requireFrozenValidationTrialLimits`. Check parsing edge cases that could admit a value other than 0.15, including invalid decimals, `NaN`, infinities, negative values, and surrounding whitespace. Check timeout equality at the `time.Duration` level.

Run `TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact`. Confirm its synthetic custodian writes a marker on any invocation and that all four wrong-limit cases prove both no output and no contact. Candidate tests are evidence, not a substitute for source inspection. Add disposable public-only probes if necessary; do not edit the tracked gate or candidate.

## Regression audit

Confirm the replacement preserves every previously passing boundary.

### Gate, tranche, and held-out boundary

- Validation requires `--tranche validation`, `--case all`, and a schedule.
- Validation schedule generation and single-arm validation probes remain disabled.
- A closed gate refuses before output mutation, build, runtime verification, model invocation, or custodian contact.
- An open synthetic gate requires status `complete`, `may_open_validation: true`, `may_open_held_out: false`, and exact frozen public identities.
- No `held_out` loader, fixture, schedule, grader, or CLI path exists.

### Public schedule and trial contract

Using public files only, independently reconstruct or otherwise verify the exact validation schedule:

- version 1, tranche `validation`, seed `20260817`, three repetitions, arms A–E;
- 120 cases in 24 five-case family blocks;
- 360 complete pairing keys and 1,800 launches;
- canonical schedule digest `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`;
- validation manifest canonical digest `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490`;
- public label-registry raw digest `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816`.

Confirm the frozen A–E prompts, surfaces, runtime invocation, Arm B document, world, 12-turn limit, retry policy, pairing-budget rule, ITT treatment, safety stops, and evidence structure remain unchanged except for the already-reviewed custodian metadata and the corrected trial-limit enforcement.

### External custodian boundary

Re-audit the external grader without private material:

- absolute regular-file path outside the implementation repository after symlink resolution;
- strict `describe` binding to version 1, validation, 120 cases, grader digest `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`, schedule, manifest, public registry, and private archive `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4`;
- executable hash recorded, rechecked after handshake, and rechecked before each grade;
- grade arguments limited to the public case ID and completed retained trial path;
- no label path, private root, expected outcome, archive path, or registry content passed or leaked;
- 30-second command bound, 4 MiB stdout cap, discarded stderr, generic failures, strict normalized grade JSON, and no fallback to the authoring grader.

Record and judge the existing hash-then-exec TOCTOU limitation under the stated threat model.

### Evidence and authoring compatibility

Confirm validation summaries retain tranche, schedule digest, grader digest, custodian description, executable hash, frozen cap, and frozen timeout. Runtime, trial, grade, assignment, retry, unresolved, budget-stop, and safety-stop evidence must preserve the accepted semantics. Grade or runtime safety failure must leave the scheduled run indeterminate.

Confirm authoring still derives only `labels/authoring/<case_id>.json`, existing authoring tests pass, historical freeze guards verify accepted historical bytes, and the candidate guard does not claim the replacement is already frozen.

## Required public checks

Run at least these commands from the detached report checkout:

```powershell
git status --porcelain
git diff-tree --no-commit-id --name-status -r 4bba4b4f988e2f9bbc4dbc986fbe7413078b9989
git diff-tree --no-commit-id --name-status -r 7190e8cb4b304d437af05ecdddd9082537bed036
go test ./experiments/frontier-v1/runner -run TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact -count=1
go test ./...
make quality
go run ./experiments/frontier-v1/runner --repo-root . --tranche validation --case all --schedule experiments/frontier-v1/schedules/validation.json --arm-b-document experiments/frontier-v1/arms/arm-b.md --validation-grader D:\does-not-exist\grader.exe --max-budget-usd 0.15 --timeout 180s --run-budget-usd 300
git diff --check 31bdabc0ea2a3395ab50a83c2fcbab5e3938f537 4bba4b4f988e2f9bbc4dbc986fbe7413078b9989
```

The direct CLI probe must stop with `validation gate is closed`, must not inspect the nonexistent grader, and must not create `experiments/frontier-v1/results/scheduled-validation`. Do not edit the tracked gate to reach later CLI stages. Temporary synthetic-gate unit tests are allowed.

If `make quality` creates ignored `bin/` or `build/` products, record them. Remove only those verified generated paths if cleanup is permitted. They are not candidate evidence.

## Review record and verdict

Write a new review record that pins:

- this assignment's commit and raw identity;
- replacement payload and report commits;
- inventory, report, file-set, schedule, manifest, registry, and world-build identities;
- exact checkout, commands, outputs, worktree status, custody statement, and independence declaration;
- every finding classified P0–P3.

`ACCEPT` requires no unresolved finding that could alter the frozen validation trial contract, execute or mutate evidence while Gate 1A is closed, contact custody code after an invalid preflight, expose or misbind private outcomes, admit held-out access, bypass the frozen schedule, change retry/ITT semantics, regress authoring, or misstate freeze status.

An `ACCEPT` verdict authorizes only a later exact-payload import/refreeze decision. It does not authorize that import, Gate 1A, validation execution, use of a real custodian, private grade access, held-out, or outcome analysis.
