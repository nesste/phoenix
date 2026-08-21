# Protocol-v4 validation execution-boundary independent review

- **Reviewer role:** `validation_execution_boundary.independent_reviewer`. Distinct from `validation_execution_boundary.candidate_author`, this assignment’s project-chair author, any later replacement-import/refreeze chair, and any later Gate 1A chair.
- **Date:** 2026-08-21
- **Verdict:** `REJECT`
- **Review-assignment commit (process identity only):** `5f7b138c40ed3b3efa9c8a80fb52c1fa9b53a994` (`docs: assign validation execution boundary review`; sole parent `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`; sole path this file’s parent commit)
- **Assignment raw bytes:** 9,265; `sha256:76e6dd3c76dc06f8dfecd344b4543767a33af39efd18c80c0f0490a9638fcb77`
- **Frozen implementation base:** `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`
- **Candidate commit:** `9de7861d473b9ec14778638938944b75aeff937c` (`experiments: add validation execution boundary candidate`)
- **Candidate-report commit:** `32aee102b78ce47e327b43a0887a8def47549de0` (`docs: report validation execution boundary candidate`)
- **Candidate report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-candidate.md`
- **Candidate report raw bytes:** 8,193; `sha256:52b8c6db6ca8ae2314951ec81d11e62384bccf977641521ea63cf90363180d89`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json` at candidate commit
- **Candidate inventory raw bytes:** 3,783; `sha256:70eabdd0edca570709f2053beed39fa56bf5acedd0fd850d5a985a784bc145a1`
- **Boundary artifact-set identity:** `sha256:9d7497742ca7fe0ccceaae6e9dd33df0347923dd382b06d73103a07f78299550`
- **Review checkout:** `D:\Work\personal\phoenix-validation-execution-boundary-review` (new detached worktree at the report commit; did not previously exist; materialized with `core.autocrlf=false` and `core.eol=lf`; `git status --porcelain` empty of tracked changes before and after every check)
- **Disposable reconstruction:** `D:\Work\personal\phoenix-validation-execution-boundary-review-scratch` (outside every Phoenix checkout; not committed)

This review does not import the candidate, modify a frozen runner or grader, refreeze replacement bytes, open Gate 1A, authorize validation, or authorize held-out. Both gates on `main` remain false. No model, arm, validation trial, private grade, outcome analysis, or prospective evaluation was run.

## Independence declaration

This session did not author the candidate, the candidate report, the candidate inventory, this assignment, Phoenix, or any freeze/import decision. A prior same-machine session prepared the candidate; that session did not write this record.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Same-machine / shared-workspace limitation is recorded and does not reuse a prohibited role.
- The implementation workspace `D:\Work\personal\phoenix` at `main` `5f7b138c…` was used only to read the assignment and to write this review record. All technical audit, reconstruction, tests, quality, and the closed-gate CLI probe ran in the detached review worktree.
- The candidate branch `codex/validation-execution-candidate` and its existing worktree were not modified.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. No validation or held-out label directory exists in the review worktree; none was opened. Authoring labels remain under `labels/authoring/` and were not inspected. Existing authoring tests and `make quality`’s authoring validator may read those files as they already do. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. No real custodian executable was substituted. `seal --write` was not run.

## Verdict

Identity, ancestry, the 16-path candidate commit, inventory/artifact-set digests, closed-gate refusal before mutation or custodian contact, public schedule reconstruction, held-out unreachability, authoring label derivation, historical freeze-guard retargeting, tests, quality, and the closed-gate CLI probe all hold.

One unresolved P1 remains: a later authorized validation run can still change the frozen per-trial cost cap and timeout through CLI flags. That can change the A–E trial contract and pairing-budget behavior. `ACCEPT` is therefore withheld.

This `REJECT` does not authorize a replacement import, refreeze, Gate 1A, validation execution, use of a real custodian, or any outcome access.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Assignment prompt at `5f7b138c…` | SHA-256 of git blob; 9,265 bytes; LF-only | `sha256:76e6dd3c76dc06f8dfecd344b4543767a33af39efd18c80c0f0490a9638fcb77` |
| Candidate commit tree vs parent `31bdabc0…` | `git diff-tree --name-status -r 9de7861d…` | exactly the 16 claimed candidate paths |
| Report commit tree vs candidate `9de7861d…` | `git diff-tree --name-status -r 32aee102…` | exactly `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-candidate.md` |
| Candidate report raw bytes at `32aee102…` | SHA-256 of git blob; 8,193 bytes; LF-only | `sha256:52b8c6db6ca8ae2314951ec81d11e62384bccf977641521ea63cf90363180d89` |
| Candidate inventory raw bytes at `9de7861d…` | SHA-256 of git blob; 3,783 bytes; LF-only | `sha256:70eabdd0edca570709f2053beed39fa56bf5acedd0fd850d5a985a784bc145a1` |
| Inventory file digests and ordinal-path aggregate | LF-normalized SHA-256 of each of the 16 inventory paths, then SHA-256 of sorted `path\tdigest\n` lines | every file digest matched; aggregate `sha256:9d7497742ca7fe0ccceaae6e9dd33df0347923dd382b06d73103a07f78299550` |
| Validation manifest canonical JSON | Independent RFC 8785 integer/string/object/array subset (UTF-16 key order) | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| Validation label-registry raw bytes | SHA-256 of git/worktree bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Validation schedule canonical JSON | Independent reconstruction then same encoder; committed file deep-equal | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| World-build digest from `make quality` | `cmd/build-manifest` then `quality-check compare` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: detached HEAD `32aee102…`; sole parent candidate `9de7861d…`; that commit’s sole parent frozen base `31bdabc0…`. Assignment `5f7b138c…` is a sibling of the candidate on `main`, not an ancestor of the reviewed payload.

Inventory state at the candidate commit: `status: review_candidate`, `frozen: false`, `execution_performed: false`, `gates.may_open_validation: false`, `gates.may_open_held_out: false`, 16 boundary files. Live `pre-validation-artifacts.json` in the review checkout remains `status: complete`, `remaining: []`, both gates false. Protocol `gate.may_open_validation` and `may_open_held_out` remain false.

The 16-path candidate patch and the 16-file inventory set are not the same set. That is the candidate report’s claimed provenance split, and it is accurate: the inventory covers the live replacement contract/README/production runner sources/behavior tests/new validation tests, including four unchanged test files (`arm_b_test.go`, `runtime_test.go`, `schedule_test.go`, `scheduled_run_test.go`) whose blobs equal parent `31bdabc0…`. The commit additionally changes the inventory JSON and three freeze-guard/candidate-guard tests, which are outside the artifact-set digest and were reviewed as provenance controls.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Closed gate and tranche boundary | Pass | `--tranche validation` requires `--case all` and `--schedule`. `--write-schedule` and single-arm probes remain authoring-only. `normalizeTranche` admits only `authoring` and `validation`. No `held_out` loader, fixture, schedule, grader, or CLI path exists in the runner. `prepareScheduledCLI` calls `requireValidationGate` before `prepareScheduledInputs` (custodian), `requireEmptyScheduledOutput`, `prepareRunConfig` (output mkdir / Phoenix build), and `claudeDriver.Verify`. Closed-gate error is `validation gate is closed`. Open-gate checks require status `complete`, `may_open_validation: true`, `may_open_held_out: false`, frozen schedule path/digest, and source-manifest path/digest, then re-hash the live schedule, manifest, and public registry. |
| 2. Public input and frozen execution contract | Fail (P1) | Independent reconstruction from all 120 public validation cases and fixtures reproduced the committed schedule: version 1, tranche `validation`, seed `20260817`, three repetitions, arms A–E, 24 five-case family blocks, 360 pairing keys, 1,800 launches, canonical digest `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`. Manifest and registry pins matched. Schedule comparison is `validateScheduleForTranche` DeepEqual against a regenerated Williams construction plus an exact canonical-digest pin; a one-case selection, regenerated write, or arm-set change is refused. Frozen A–E prompts, surfaces, runtime invocation template, Arm B document, world, retry/ITT constants, pairing-budget formula, and evidence schema are preserved except for the added custodian metadata. **Run budget 300 USD is required; per-trial 0.15 USD and timeout 180s are defaults only and can be changed on a later authorized run.** |
| 3. Private-grading custody | Pass, with recorded TOCTOU limitation | Absolute path, `EvalSymlinks` outside the implementation repository, regular file, strict `describe` JSON equality to version 1 / tranche `validation` / 120 cases / grader `sha256:36abfbec…dcfc` / schedule / manifest / registry `sha256:847c510c…c1c816` / private archive `sha256:5a320a87…922e4`. Executable is hashed, re-hashed after handshake, and re-hashed before each grade. Grade argv is only `grade --case-id <public-id> --trial <retained-trial-path>`. Timeout 30s, stdout cap 4 MiB, stderr discarded, generic command failures. Returned JSON is strict, case-bound, nonempty, uniquely checked, verdict-limited, aggregate-consistent, re-marshaled, and retained. `corpusGrader` rejects any non-authoring tranche/case and is not used on the validation branch; a failed external handshake returns without falling back. |
| 4. Evidence, stops, and authoring compatibility | Pass | Scheduled summary records tranche, schedule digest, grader digest, `validation_grader_boundary`, and `validation_grader_adapter_sha256`. Runtime/trial/grade/assignment evidence remains retained. Grade or runtime safety failures call `recordSafetyStop` and leave the scheduled run `indeterminate`. Authoring still derives only `labels/authoring/<case_id>.json` inside `corpusGrader.Grade`. Historical runner freeze now verifies accepted bytes at `b4df9190…` rather than the live replacement tree; local runtime-file hashes are verified at `795ce71a…`; the new candidate guard verifies the live 16-file inventory and does not mark it frozen. Live `pre-validation-artifacts.json` still describes the accepted authoring runner freeze, not this replacement. |
| 5. Reproducible public checks | Pass | See commands below. Closed-gate CLI probe printed `validation gate is closed` and did not create `experiments/frontier-v1/results/scheduled-validation`. |

## Public checks

From detached `HEAD` `32aee102b78ce47e327b43a0887a8def47549de0`:

| Command | Result |
| --- | --- |
| `git status --porcelain` | empty of tracked changes before and after every check |
| `git diff-tree --no-commit-id --name-status -r 9de7861d…` | the 16 claimed candidate paths |
| `git diff-tree --no-commit-id --name-status -r 32aee102…` | only the candidate report |
| `go test ./...` | pass |
| nested `experiments/frontier-v1/corpusctl` `go test ./...` | pass (additional; nested module is outside root `./...`) |
| `make quality` | pass; world-build `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| closed-gate CLI probe with `D:\does-not-exist\grader.exe` and `--run-budget-usd 300` | exit 1; stderr `validation gate is closed`; no grader path error; `scheduled-validation` absent |
| `git diff --check 31bdabc0… 9de7861d…` | empty |

`make quality` created gitignored `bin/phoenix` and `build/manifest.json` in the disposable review worktree. They are not candidate evidence and were removed afterward. Tracked worktree status remained empty.

## Custody of the external executable (TOCTOU)

The runner hashes path contents, then execs that same path string. There is a time-of-check/time-of-use window between `digestRawFile` / `os.ReadFile` and `exec.Command`. The repeated identity checks (after `describe`, and before every `grade`) detect substitution across handshake and between grades. They cannot bind a single invocation to the hashed bytes on a hostile co-resident writer.

For this threat model — operator-supplied custodian, prevent in-repo label loading, detect adapter replacement across the run, no private material in the implementation workspace — those repeated checks are sufficient. A process that can replace the custodian between hash and exec already has write access to that executable. This is an accepted limitation, not an `ACCEPT`-blocking finding.

## Findings

### [P1] Pin validation per-trial cost and timeout to the frozen contract — `experiments/frontier-v1/runner/main.go:46`

A later authorized validation run can change frozen trial limits without changing the schedule file. `--run-budget-usd` must be exactly `300`, but `--max-budget-usd` and `--timeout` remain ordinary CLI flags. Their defaults are `0.15` and `180s`; `prepareScheduledInputs` never sees the per-trial budget, and `prepareRunConfig` only requires a positive timeout and a nonempty budget string.

That matters once Gate 1A is open. Pairing-budget capacity is `perTrialCap * 5 * 1.10`. Raising the per-trial cap changes when the 300 USD run budget stops the schedule and changes cost-cap ITT behavior. Lowering it, or shortening the timeout, likewise changes which trials hit caps. The assignment requires exactly 0.15 USD per trial; protocol `trial.max_cost_usd_per_trial` and `trial.timeout_seconds` are part of the frozen A–E contract.

This does not leak labels, execute while the gate is closed, or admit held-out. It does allow a later authorized run to change budgets and trial-contract limits, so `ACCEPT` is withheld.

## Accepted limitations

- Same-machine fresh session, recorded above.
- Hash-then-exec TOCTOU on the custodian path, judged sufficient for this threat model.
- `--runtime` remains an operator path, as in authoring; the frozen invocation template is still built in `prepareRuntimeInvocation` and recorded in summary configuration.
- Independent reconstruction used public cases, fixtures, manifest, and registry only. It did not open label files or a private archive. The private-archive digest is checked only as a custodian `describe` string.
- Root `go test ./...` does not enter nested `corpusctl` or `surface-spike` modules; `make quality` still ran `validate-authoring` and spec validation. Nested corpusctl tests were run separately and passed.

## Smallest next artifact

A revised candidate that additionally refuses validation unless `--max-budget-usd` is exactly `0.15` and `--timeout` is exactly `180s` (or equivalent locked constants), with tests that a different value is rejected before output creation or custodian contact. This review does not specify that patch and does not authorize it.

This `REJECT` does not authorize import, refreeze, Gate 1A, validation execution, use of a real custodian, or any outcome access.
