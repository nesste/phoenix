# Protocol-v4 validation execution-boundary independent review assignment

- **Assignment owner:** project-chair workflow
- **Assigned reviewer role:** `validation_execution_boundary.independent_reviewer`
- **Date:** 2026-08-21
- **Frozen implementation base:** `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537` on `main`
- **Candidate branch:** `codex/validation-execution-candidate`
- **Candidate commit:** `9de7861d473b9ec14778638938944b75aeff937c`
- **Candidate report commit:** `32aee102b78ce47e327b43a0887a8def47549de0`
- **Candidate report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-candidate.md`
- **Candidate report raw identity:** 8,193 bytes; `sha256:52b8c6db6ca8ae2314951ec81d11e62384bccf977641521ea63cf90363180d89`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json` at candidate commit
- **Candidate inventory raw identity:** 3,783 bytes; `sha256:70eabdd0edca570709f2053beed39fa56bf5acedd0fd850d5a985a784bc145a1`
- **Boundary artifact-set identity:** `sha256:9d7497742ca7fe0ccceaae6e9dd33df0347923dd382b06d73103a07f78299550`

This assignment does not import the candidate, modify a frozen runner or grader, refreeze replacement bytes, open Gate 1A, authorize validation, or authorize held-out. Both gates on `main` remain false. No model, arm, validation trial, private grade, outcome analysis, or prospective evaluation may be run during this review.

## Independence requirement

The reviewer must not be the candidate author, this assignment author, a later replacement-import/refreeze chair, or a later Gate 1A chair. The review record must declare any same-machine or shared-workspace limitation. Use a clean disposable worktree or detached checkout of `32aee102b78ce47e327b43a0887a8def47549de0`; do not modify the candidate branch.

The reviewer must not open, mount, enumerate, hash, or copy any private label root or validation/held-out label file. Do not substitute a real custodian executable. Synthetic outcome-free test helpers are permitted. Authoring labels need not be opened.

## Reviewed object and ancestry

Verify independently:

1. candidate `9de7861d473b9ec14778638938944b75aeff937c` has sole parent `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`;
2. report commit `32aee102b78ce47e327b43a0887a8def47549de0` has sole parent candidate `9de7861…` and changes only the candidate report;
3. candidate `9de7861…` changes exactly the 16 paths listed in the report;
4. the candidate inventory is raw-byte identical to the pinned identity above, contains 16 boundary files, is `review_candidate`, `frozen: false`, keeps both gates false, and states `execution_performed: false`;
5. every inventory file digest and the ordinal-path aggregate digest reproduce exactly.

Review the exact candidate commit, not later branch state. The report-only child is context and a claimed verification record, not part of the implementation payload.

## Required technical audit

### 1. Closed gate and tranche boundary

Confirm validation requires `--tranche validation`, `--case all`, and a schedule. Confirm schedule generation and single-arm probes remain authoring-only. Confirm there is no `held_out` loader, fixture, schedule, grader, or CLI path.

Confirm a closed validation gate refuses before output-directory creation, Phoenix build, runtime verification, model invocation, or custodian contact. An open validation gate must require `pre-validation-artifacts.json` status `complete`, `may_open_validation: true`, `may_open_held_out: false`, a frozen validation schedule entry, and exact accepted schedule/source-manifest pins.

### 2. Public input and frozen execution contract

Using public files only, reconstruct and validate the committed validation schedule from all 120 public cases and fixtures. Require version 1, tranche `validation`, seed `20260817`, three repetitions, A–E arms, 24 five-case family blocks, 360 complete pairing keys, 1,800 launches, and canonical digest `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`.

Require the public validation manifest canonical digest `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` and public validation label-registry raw digest `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816`. Confirm the candidate preserves the frozen A–E prompts, surfaces, runtime invocation, Arm B document, world, trial cap, retry policy, pairing-budget behavior, ITT rules, and evidence structure, except for the expressly added custodian metadata.

Require exactly 300 USD for the validation run budget and 0.15 USD per trial. Reject any path that can select one validation case, regenerate a validation schedule, change an arm set, or bypass the deterministic schedule comparison.

### 3. Private-grading custody

Audit the external-custodian boundary end to end:

- the path must be absolute, resolve outside the implementation repository (including symlinks), and name a regular file;
- strict `describe` JSON must match version 1, tranche `validation`, 120 cases, grader `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`, the accepted schedule and public-manifest identities, public registry `sha256:847c510c…c1c816`, and private archive `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4`;
- the external executable raw hash must be recorded, checked again after handshake, and checked before each grade;
- the grade command must receive only the public case ID and completed retained trial path, never a label path, private root, expected outcome, archive path, or registry content;
- subprocess timeout, stdout cap, discarded stderr, and generic command failures must prevent unbounded or accidental private-path leakage;
- returned grade JSON must be strict, case-bound, nonempty, uniquely checked, verdict-limited, aggregate-consistent, normalized, and retained;
- no error or fallback path may invoke the in-repository authoring grader for validation.

Call out any time-of-check/time-of-use limitation around the external executable and judge whether the repeated identity checks are sufficient for this threat model. Do not test this with private material.

### 4. Evidence, stops, and authoring compatibility

Confirm validation summary evidence identifies tranche, schedule digest, grader digest, full custodian description, and executable hash. Confirm existing runtime/trial/grade/assignment evidence remains retained and that grade or runtime safety failures produce an indeterminate scheduled run without silently continuing or relabelling outcomes.

Run existing authoring tests and audit the interface change to ensure authoring still derives only `labels/authoring/<case_id>.json` internally. Confirm the historical freeze guards accurately verify accepted historical bytes, and the new candidate guard covers the intended live boundary set without misrepresenting the accepted `main` freeze.

## Required public checks

At minimum, from the exact candidate/report checkout:

```powershell
git status --porcelain
git diff-tree --no-commit-id --name-status -r 9de7861d473b9ec14778638938944b75aeff937c
git diff-tree --no-commit-id --name-status -r 32aee102b78ce47e327b43a0887a8def47549de0
go test ./...
make quality
go run ./experiments/frontier-v1/runner --repo-root . --tranche validation --case all --schedule experiments/frontier-v1/schedules/validation.json --arm-b-document experiments/frontier-v1/arms/arm-b.md --validation-grader D:\does-not-exist\grader.exe --run-budget-usd 300
git diff --check 31bdabc0ea2a3395ab50a83c2fcbab5e3938f537 9de7861d473b9ec14778638938944b75aeff937c
```

The direct CLI probe must refuse because the gate is closed, before checking the nonexistent grader, and must not create `experiments/frontier-v1/results/scheduled-validation`. Do not edit the gate to exercise the real CLI path. Unit tests may use a temporary synthetic gate document and synthetic custodian outside a temporary repository.

If `make quality` creates gitignored `bin/` or `build/` products in the disposable review worktree, report them and remove only those verified generated paths afterward if permitted. They are not candidate evidence.

## Review deliverable and verdict

Write a new review record on `main` that pins this assignment commit, candidate/report commits, raw identities, artifact-set digest, exact checkout, commands, results, custody statement, and independence declaration. Classify every finding P0–P3.

`ACCEPT` requires no unresolved finding that could:

- access, copy, infer, or leak private labels beyond normalized retained grade outcomes during a later authorized run;
- execute or mutate evidence while Gate 1A is closed;
- admit held-out access;
- change or bypass the frozen schedule, A–E trial contract, budgets, retry/ITT rules, or case/outcome binding;
- substitute or change the custodian executable undetected under the accepted threat model;
- regress authoring execution or falsely claim the replacement is already frozen.

An `ACCEPT` verdict authorizes only a later project-chair exact-payload import/refreeze decision. It does not authorize that import, Gate 1A, validation execution, use of a real custodian, or any outcome access.
