# Protocol-v4 validation execution-boundary revision independent review

- **Reviewer role:** `validation_execution_boundary.revision_independent_reviewer`. Distinct from `validation_execution_boundary.candidate_author`, `validation_execution_boundary.independent_reviewer`, `validation_execution_boundary.revision_author`, this assignment’s project-chair author, any later replacement-import/refreeze chair, and any later Gate 1A chair.
- **Date:** 2026-08-21
- **Verdict:** `REVISE`
- **Review-assignment commit (process identity only):** `8eb219322b91b9008d097f1d7089fc51dbda2e8e` (`docs: assign validation execution boundary revision review`; sole parent `6f1696542c4033fe3edd660c95b07b2e74f292f5`; sole path this file’s sibling prompt)
- **Assignment raw bytes:** 11,008; `sha256:2d03bc786068f58956e19800a368976922f8ce599f52c2d7b5177b65b7db427d`
- **Frozen implementation base:** `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`
- **Rejected payload:** `9de7861d473b9ec14778638938944b75aeff937c`
- **Reject review:** `6f1696542c4033fe3edd660c95b07b2e74f292f5`; independently verified 15,808 raw bytes; `sha256:22bf9ce094021c709d9005f90f4a486167a5bf5aa4a2dd65017ad300210bc7f0`
- **Replacement payload:** `4bba4b4f988e2f9bbc4dbc986fbe7413078b9989` (`experiments: revise validation execution boundary candidate`)
- **Replacement-report commit:** `7190e8cb4b304d437af05ecdddd9082537bed036` (`docs: report validation execution boundary revision`)
- **Replacement report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-revision-candidate.md`
- **Replacement report raw bytes:** 5,478; `sha256:19b073280e3f1770680b109ff389307c9ef3651ff4b27a71cd814e87c1dae0c2`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json` at the replacement payload
- **Candidate inventory raw bytes:** 4,205; `sha256:680fc24683f5d8dedf1c74634026d94da40b703591678a6fcc3b214b9a690d8d`
- **Boundary artifact-set identity:** `sha256:012784de3eb93f6c84f7b93008bd7f8e6834b2d53427c043661c0846e670ea1e`
- **Review checkout:** `D:\Work\personal\phoenix-validation-execution-boundary-revision-review` (new detached worktree at the report commit; did not previously exist; materialized with `core.autocrlf=false` and `core.eol=lf`; `git status --porcelain` empty of tracked changes before and after every check)
- **Disposable reconstruction and probes:** `D:\Work\personal\phoenix-validation-execution-boundary-revision-review-scratch` (outside every Phoenix checkout; not committed)

This review does not import the replacement, modify a frozen runner or grader, refreeze replacement bytes, open Gate 1A, authorize validation, or authorize held-out. Both gates on `main` remain false. No model, arm, validation trial, private grade, outcome analysis, or prospective evaluation was run.

## Independence declaration

This session did not author the rejected candidate, the original independent `REJECT`, the replacement payload, the replacement report, the replacement inventory, this assignment, Phoenix, or any freeze/import decision.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Same-machine / shared-workspace limitation is recorded and does not reuse a prohibited role.
- The implementation workspace `D:\Work\personal\phoenix` at `main` `8eb21932…` was used only to read the assignment and to write this review record. All technical audit, reconstruction, tests, quality, the closed-gate CLI probe, and disposable probes ran in the detached review worktree or the scratch directory.
- The replacement branch `codex/validation-execution-candidate-revision` and its existing author worktree were not modified.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. No validation or held-out label directory exists in the review worktree; none was opened. Authoring labels remain under `labels/authoring/` and were not inspected. Existing authoring tests and `make quality`’s authoring validator may read those files as they already do. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. No real custodian executable was substituted. Synthetic, outcome-free helpers ran only in temporary directories and the scratch path above. `seal --write` was not run.

## Verdict

Identity, ancestry, the 16-path replacement payload, inventory/artifact-set digests, closed-gate refusal before mutation or custodian contact, public schedule reconstruction, held-out unreachability, authoring label derivation, historical freeze-guard retargeting, tests, quality, and the closed-gate CLI probe all hold.

The original P1 is closed for the four changed numeric/duration values named in the replacement report: `0.14`, `0.16`, `179s`, and `181s` are rejected after the closed-gate check and before output mutation or custodian contact. The 300 USD run budget remains locked. Timeout equality is at the `time.Duration` level.

One unresolved P1 remains: surrounding whitespace on an otherwise-numeric `0.15` cap passes `requireFrozenValidationTrialLimits` and then reaches custodian `describe` contact. `ACCEPT` is therefore withheld.

This `REVISE` does not authorize a replacement import, refreeze, Gate 1A, validation execution, use of a real custodian, or any outcome access.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Assignment prompt at `8eb21932…` | SHA-256 of git blob; 11,008 bytes; LF-only | `sha256:2d03bc786068f58956e19800a368976922f8ce599f52c2d7b5177b65b7db427d` |
| Reject review at `6f169654…` | SHA-256 of git blob; 15,808 bytes; LF-only | `sha256:22bf9ce094021c709d9005f90f4a486167a5bf5aa4a2dd65017ad300210bc7f0` |
| Replacement payload vs parent `31bdabc0…` | `git diff-tree --name-status -r 4bba4b4f…` | exactly the 16 claimed candidate paths |
| Report commit tree vs payload `4bba4b4f…` | `git diff-tree --name-status -r 7190e8cb…` | exactly `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-revision-candidate.md` |
| Replacement report raw bytes at `7190e8cb…` | SHA-256 of git blob; 5,478 bytes; LF-only | `sha256:19b073280e3f1770680b109ff389307c9ef3651ff4b27a71cd814e87c1dae0c2` |
| Candidate inventory raw bytes at `4bba4b4f…` | SHA-256 of git blob; 4,205 bytes; LF-only | `sha256:680fc24683f5d8dedf1c74634026d94da40b703591678a6fcc3b214b9a690d8d` |
| Inventory file digests and ordinal-path aggregate | LF-normalized SHA-256 of each of the 16 inventory paths, then SHA-256 of sorted `path\tdigest\n` lines | every file digest matched; aggregate `sha256:012784de3eb93f6c84f7b93008bd7f8e6834b2d53427c043661c0846e670ea1e` |
| Validation manifest canonical JSON | Independent RFC 8785 integer/string/object/array subset (UTF-16 key order) | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| Validation label-registry raw bytes | SHA-256 of git/worktree bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Validation schedule canonical JSON | Independent reconstruction then same encoder; committed file deep-equal | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| World-build digest from `make quality` | `cmd/build-manifest` then `quality-check compare` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: detached HEAD `7190e8cb…`; sole parent payload `4bba4b4f…`; that commit’s sole parent frozen base `31bdabc0…`. Assignment `8eb21932…` is a child of the reject-review commit on `main`, not an ancestor of the reviewed payload. The replacement was reviewed as a complete candidate from the frozen base, not as a patch over `9de7861d…`. Versus that rejected payload, the replacement changes six paths: the trial-limit lock, its tests, the inventory, the contract document, and README.

Inventory state at the payload commit: `status: review_candidate`, `frozen: false`, `execution_performed: false`, `gates.may_open_validation: false`, `gates.may_open_held_out: false`, 16 boundary files, `run_budget_usd: 300`, `max_cost_usd_per_trial: 0.15`, `timeout_seconds: 180`. It names rejected payload `9de7861d…` and the P1 finding. Live `pre-validation-artifacts.json` in the review checkout remains `status: complete`, `remaining: []`, both gates false. Protocol `gate.may_open_validation` and `may_open_held_out` remain false.

The 16-path payload patch and the 16-file inventory set are not the same set. That split is accurate: the inventory covers the live replacement contract/README/production runner sources/behavior tests/new validation tests, including four unchanged test files (`arm_b_test.go`, `runtime_test.go`, `schedule_test.go`, `scheduled_run_test.go`) whose blobs equal parent `31bdabc0…`. The commit additionally changes the inventory JSON and three freeze-guard/candidate-guard tests, which are outside the artifact-set digest and were reviewed as provenance controls. No private, result, runtime-stream, trial, grade, analysis-output, or gate-state path is in the payload.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Closed gate and tranche boundary | Pass | `--tranche validation` requires `--case all` and a schedule. `--write-schedule` and single-arm probes remain authoring-only. `normalizeTranche` admits only `authoring` and `validation`. No `held_out` loader, fixture, schedule, grader, or CLI path exists in the runner. `prepareScheduledCLI` calls `requireValidationGate` before `requireFrozenValidationTrialLimits`, `prepareScheduledInputs` (custodian), `requireEmptyScheduledOutput`, `prepareRunConfig` (output mkdir / Phoenix build), and `claudeDriver.Verify`. Closed-gate error is `validation gate is closed`. Open-gate checks require status `complete`, `may_open_validation: true`, `may_open_held_out: false`, frozen schedule path/digest, and source-manifest path/digest, then re-hash the live schedule, manifest, and public registry. |
| 2. Public input and frozen execution contract | Fail (P1, whitespace only) | Independent reconstruction from all 120 public validation cases and fixtures reproduced the committed schedule: version 1, tranche `validation`, seed `20260817`, three repetitions, arms A–E, 24 five-case family blocks, 360 pairing keys, 1,800 launches, canonical digest `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`. Manifest and registry pins matched. Frozen A–E prompts, surfaces, runtime invocation template, Arm B document, world, 12-turn limit, retry/ITT constants, pairing-budget formula, and evidence schema are preserved except for the added custodian metadata and the new trial-limit preflight. Run budget 300 USD remains required. Numeric caps `0.14`/`0.16` and timeouts `179s`/`181s` are rejected before output or custodian contact. Surrounding whitespace on `0.15` is not. |
| 3. Private-grading custody | Pass, with recorded TOCTOU limitation | Absolute path, `EvalSymlinks` outside the implementation repository, regular file, strict `describe` JSON equality to version 1 / tranche `validation` / 120 cases / grader `sha256:36abfbec…dcfc` / schedule / manifest / registry `sha256:847c510c…c1c816` / private archive `sha256:5a320a87…922e4`. Executable is hashed, re-hashed after handshake, and re-hashed before each grade. Grade argv is only `grade --case-id <public-id> --trial <retained-trial-path>`. Timeout 30s, stdout cap 4 MiB, stderr discarded, generic command failures. Returned JSON is strict, case-bound, nonempty, uniquely checked, verdict-limited, aggregate-consistent, re-marshaled, and retained. `corpusGrader` rejects any non-authoring tranche/case and is not used on the validation branch; a failed external handshake returns without falling back. |
| 4. Evidence, stops, and authoring compatibility | Pass | Scheduled summary records tranche, schedule digest, grader digest, `validation_grader_boundary`, `validation_grader_adapter_sha256`, frozen `timeout_seconds`, and `max_cost_usd_per_trial`. Runtime/trial/grade/assignment evidence remains retained. Grade or runtime safety failures call `recordSafetyStop` and leave the scheduled run `indeterminate`. Authoring still derives only `labels/authoring/<case_id>.json` inside `corpusGrader.Grade`. Authoring probes and authoring scheduled runs do not call `requireFrozenValidationTrialLimits`. Historical runner freeze still verifies accepted bytes at `b4df9190…`; local runtime-file hashes are verified at `795ce71a…`; the candidate guard verifies the live 16-file inventory and does not mark it frozen. Live `pre-validation-artifacts.json` still describes the accepted authoring runner freeze, not this replacement. |
| 5. Reproducible public checks | Pass | See commands below. Closed-gate CLI probe printed `validation gate is closed` and did not create `experiments/frontier-v1/results/scheduled-validation`. The nonexistent grader path was not inspected. |

## Call order and trial-limit lock

`prepareScheduledCLI` for tranche `validation`:

1. `requireValidationGate`
2. `requireFrozenValidationTrialLimits(budget, timeout)`
3. default output path string only
4. `prepareScheduledInputs` (schedule load/validation, 300 USD check, then custodian path resolution and `describe`)
5. case preflight
6. `requireEmptyScheduledOutput`
7. `prepareRunConfig` (output mkdir, Phoenix build)

`requireFrozenValidationTrialLimits` parses the cap with `strconv.ParseFloat(strings.TrimSpace(budget), 64)` and requires equality to the `0.15` constant; timeout must equal `180 * time.Second`. `NaN`, infinities, negatives, empty strings, and invalid decimals fail that check. `ParseFloat("0.15")` has the same bits as the Go constant. Timeout `3m` equals `180s`; `179s` and `181s` do not.

The accepted string `"0.15"` and duration `180s` then reach `prepareRunConfig` unchanged. Summary evidence records `MaxCostUSDPerTrial` as that same `float64` and `TimeoutSeconds` as `int(timeout.Seconds())`. Pairing-budget capacity remains `perTrialCap * 5 * 1.10`.

## Public checks

From detached `HEAD` `7190e8cb4b304d437af05ecdddd9082537bed036`:

| Command | Result |
| --- | --- |
| `git status --porcelain` | empty of tracked changes before and after every check |
| `git diff-tree --no-commit-id --name-status -r 4bba4b4f…` | the 16 claimed candidate paths |
| `git diff-tree --no-commit-id --name-status -r 7190e8cb…` | only the replacement report |
| `go test ./experiments/frontier-v1/runner -run TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact -count=1` | pass; synthetic custodian writes a marker on any invocation; all four wrong-limit cases produced the named error, no output directory, and no marker |
| `go test ./...` | pass |
| nested `experiments/frontier-v1/corpusctl` `go test ./...` | pass (additional; nested module is outside root `./...`) |
| `make quality` | pass; world-build `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| closed-gate CLI probe with `D:\does-not-exist\grader.exe`, `--max-budget-usd 0.15`, `--timeout 180s`, and `--run-budget-usd 300` | exit 1; stderr `validation gate is closed`; no grader path error; `scheduled-validation` absent |
| `git diff --check 31bdabc0… 4bba4b4f…` | empty |

`make quality` created gitignored `bin/phoenix` and `build/manifest.json` in the disposable review worktree. They are not candidate evidence and were removed afterward. Tracked worktree status remained empty.

## Custody of the external executable (TOCTOU)

The runner hashes path contents, then execs that same path string. There is a time-of-check/time-of-use window between `digestRawFile` / `os.ReadFile` and `exec.Command`. The repeated identity checks (after `describe`, and before every `grade`) detect substitution across handshake and between grades. They cannot bind a single invocation to the hashed bytes on a hostile co-resident writer.

For this threat model — operator-supplied custodian, prevent in-repo label loading, detect adapter replacement across the run, no private material in the implementation workspace — those repeated checks are sufficient. A process that can replace the custodian between hash and exec already has write access to that executable. This is an accepted limitation, not an `ACCEPT`-blocking finding.

## Findings

### [P1] Reject whitespace-padded validation caps before custodian contact — `experiments/frontier-v1/runner/validation_gate.go:74`

`requireFrozenValidationTrialLimits` trims the cap string before `ParseFloat`. `runScheduledCases` later parses `config.budgetUSD` without trimming. A value such as `" 0.15 "` therefore passes the frozen-limit preflight as numeric `0.15`, then is invalid for the execution parser.

A disposable public-only probe opened a synthetic gate, copied the 120 public validation cases, and supplied a synthetic custodian that writes a marker on any invocation. `prepareScheduledCLI(..., budget: " 0.15 ", timeout: 180s, runBudget: "300")` did not return the frozen-limit error and set the custodian marker. Independent `strconv` probes also showed `" 0.15"`, `"0.15 "`, `"\t0.15\n"`, and `"  0.15  "` as preflight-pass / later-fail. `NaN`, infinities, negatives, `0.14`, `0.16`, and invalid decimals remain rejected at the preflight.

On a complete repository tree, that same successful preflight continues into `prepareScheduledInputs` (custodian `describe`), output-directory inspection, Phoenix build, and runtime verification before `runScheduledCases` can fail. The replacement’s claimed “before schedule preparation, output-directory inspection or creation, … or custodian describe contact” guarantee therefore does not hold for surrounding whitespace. This does not admit a completed trial at a changed numeric cap, leak labels, or open held-out. It does contact custody code after an invalid preflight, so `ACCEPT` is withheld.

## Accepted limitations

- Same-machine fresh session, recorded above.
- Hash-then-exec TOCTOU on the custodian path, judged sufficient for this threat model.
- `--runtime` remains an operator path, as in authoring; the frozen invocation template is still built in `prepareRuntimeInvocation` and recorded in summary configuration.
- Numeric-equivalent spellings of `0.15` that `ParseFloat` accepts without trimming (`0.150`, `.15`, `15e-2`, hex-float forms that round to the same bits) pass both parsers and do not change pairing-budget arithmetic or summary `float64` evidence. They are not a changed cap under the stated numeric rule.
- Independent reconstruction used public cases, fixtures, manifest, and registry only. It did not open label files or a private archive. The private-archive digest is checked only as a custodian `describe` string.
- Root `go test ./...` does not enter nested `corpusctl` or `surface-spike` modules; `make quality` still ran `validate-authoring` and spec validation. Nested corpusctl tests were run separately and passed.

## Smallest next artifact

A revised candidate that rejects a validation `--max-budget-usd` value unless the untrimmed string parses to exactly `0.15` (remove the `TrimSpace`, or otherwise make the preflight parser identical to `runScheduledCases`), plus a test that `" 0.15 "` with `180s` produces the cap error, no output, and no custodian marker. Keep the existing four changed-value cases. This review does not specify that patch and does not authorize it.

This `REVISE` does not authorize import, refreeze, Gate 1A, validation execution, use of a real custodian, or any outcome access.
