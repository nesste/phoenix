# Protocol-v4 validation execution-boundary second revision independent review

- **Reviewer role:** `validation_execution_boundary.second_revision_independent_reviewer`. Distinct from `validation_execution_boundary.candidate_author`, `validation_execution_boundary.independent_reviewer`, `validation_execution_boundary.revision_author`, `validation_execution_boundary.revision_independent_reviewer`, `validation_execution_boundary.second_revision_author`, any later replacement-import/refreeze chair, and any later Gate 1A chair.
- **Date:** 2026-08-21
- **Verdict:** `ACCEPT`
- **Review assignment:** conversational, from the project chair. No `docs: assign …` commit exists for this round; `main` was at `41c75deab949091cae6d1757a214e835a4f7af18` when the review opened and was not advanced before this record.
- **Frozen implementation base:** `31bdabc0ea2a3395ab50a83c2fcbab5e3938f537`
- **Prior replacement payload (`REVISE`d):** `4bba4b4f988e2f9bbc4dbc986fbe7413078b9989`
- **Prior review:** `41c75deab949091cae6d1757a214e835a4f7af18`; independently verified 19,624 raw bytes; `sha256:a0edbae4a695aaeec0ce45ed967427c733ed94a5014154705ac3e39477dc9e6a`
- **Second replacement payload:** `74d06da13f62cffa4fd635e048e6331a6e2d95a6` (`experiments: reject padded validation cost caps`)
- **Second replacement-report commit:** `05ab169acdfabae4ce7b42b4ffb8d6ecb276e935` (`docs: report validation execution boundary second revision`)
- **Second replacement report:** `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-candidate.md`; independently verified 5,235 raw bytes; `sha256:c77207aaebb6876805354e720adc7e72301ab759e2affaa3df3bb7c88526c546`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/validation-execution-boundary-candidate.json` at the replacement payload; independently verified 4,244 raw bytes; `sha256:14b0761362709780ded9f6fc47e7ff8d49f688e09062d8dcc6d13724b9109d1d`
- **Boundary artifact-set identity:** `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1`
- **Review checkout:** `D:\Work\personal\phoenix-validation-execution-boundary-second-revision-review` (new detached worktree at the report commit; did not previously exist; `core.autocrlf=false`, `core.eol=lf`; `git status --porcelain` empty of tracked changes before and after every check)
- **Disposable reconstruction and probes:** `D:\Work\personal\phoenix-validation-execution-boundary-second-revision-review-scratch` (outside every Phoenix checkout; not committed)

This review does not import the replacement, modify a frozen runner or grader, refreeze replacement bytes, open Gate 1A, authorize validation, or authorize held-out. Both gates on `main` remain false. No model, arm, validation trial, private grade, outcome analysis, or prospective evaluation was run.

## Independence declaration

This session did not author the rejected candidate, the original independent `REJECT`, the first replacement payload, the prior `REVISE`, the second replacement payload, the second replacement report, either inventory, Phoenix, or any freeze/import decision.

Recorded limitations:

- The reviewer is a fresh session on the same machine as earlier Phoenix work. The same-machine / shared-workspace limitation is recorded and does not reuse a prohibited role.
- The implementation workspace `D:\Work\personal\phoenix` at `main` `41c75dea…` was used only to write this review record. Every technical audit, reconstruction, test, quality run, closed-gate CLI probe, and disposable probe ran in the detached review worktree or the scratch directory.
- The candidate branch `codex/validation-execution-candidate-revision-2` and its author worktree were not modified. `git status --porcelain` there was empty before and after.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. `experiments/frontier-v1/labels/` in the review checkout contains only `authoring/`; no validation or held-out label directory exists and none was created or opened. Existing authoring tests and `make quality`'s authoring validator may read authoring labels as they already do. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. No real custodian executable was substituted — every custodian in this review was the synthetic, outcome-free helper built from the candidate's own test source, running only in temporary directories. `seal --write` was not run.

## Verdict

The single unresolved P1 from `41c75dea…` is closed. `strings.TrimSpace` is gone from the validation cap preflight, both Go parse sites now read the same unmodified cap string, and a whitespace-padded `0.15` is rejected before schedule preparation, output-directory inspection or creation, Phoenix build, runtime verification, and custodian path resolution or handshake.

I confirmed this by differential probe rather than by reading the diff alone: the identical probe, on identically populated synthetic roots, contacts the custodian on `4bba4b4f…` and does not contact it on `74d06da1…`. A 31-case parser-agreement probe found no divergence in either direction between the preflight parser and the execution parser.

Everything the prior review passed — identity, ancestry, the 16-path replacement boundary, inventory and artifact-set digests, closed-gate refusal before mutation or custodian contact, public schedule reconstruction, held-out unreachability, authoring label derivation, historical freeze-guard retargeting, tests, quality, and the closed-gate CLI probe — was re-verified against this payload and still holds. The four changed numeric and duration values (`0.14`, `0.16`, `179s`, `181s`) remain rejected at the same point. The 300 USD run budget remains locked and is also checked before custodian contact.

Two P2 findings are recorded below. Neither is `ACCEPT`-blocking: one is a precision defect in the report's description of its own test, not in the guard; the other is a pre-existing residual of the same class as the closed P1, at an external parser outside this payload's scope. Both are flagged for the import/refreeze chair.

This `ACCEPT` authorizes only a later project-chair import/refreeze decision. It does not import the replacement, refreeze bytes, open Gate 1A, authorize validation execution, authorize use of a real custodian, or authorize any outcome access.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Prior review at `41c75dea…` | SHA-256 of git blob; 19,624 bytes; LF-only | `sha256:a0edbae4a695aaeec0ce45ed967427c733ed94a5014154705ac3e39477dc9e6a` |
| Prior inventory at `4bba4b4f…` | SHA-256 of git blob; 4,205 bytes; LF-only | `sha256:680fc24683f5d8dedf1c74634026d94da40b703591678a6fcc3b214b9a690d8d` |
| Second payload vs parent `31bdabc0…` | `git diff-tree --name-status -r 74d06da1…` | exactly the 16 claimed candidate paths |
| Second payload vs prior payload `4bba4b4f…` | `git diff-tree --name-status -r 4bba4b4f… 74d06da1…` | exactly 5 paths: cap guard, validation execution test, inventory, boundary contract, README |
| Report commit tree vs payload `74d06da1…` | `git diff-tree --name-status -r 05ab169a…` | exactly the second replacement report |
| Second replacement report at `05ab169a…` | SHA-256 of git blob; 5,235 bytes; LF-only | `sha256:c77207aaebb6876805354e720adc7e72301ab759e2affaa3df3bb7c88526c546` |
| Candidate inventory at `74d06da1…` | SHA-256 of git blob; 4,244 bytes; LF-only | `sha256:14b0761362709780ded9f6fc47e7ff8d49f688e09062d8dcc6d13724b9109d1d` |
| Inventory file digests and ordinal-path aggregate | LF-normalized SHA-256 of each of the 16 inventory paths, then SHA-256 of ordinal-sorted `path\tsha256:<hex>\n` lines, in an independent Python reimplementation | all 16 file digests matched; aggregate `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1` |
| Validation schedule canonical JSON | Independent RFC 8785 subset (UTF-16 key order) written for this review | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Validation manifest canonical JSON | Same independent encoder | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| Validation label-registry raw bytes | SHA-256 of worktree bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| World-build digest from `make quality` | `cmd/build-manifest` then `quality-check compare` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: detached `HEAD` `05ab169a…`; sole parent payload `74d06da1…`; that commit's sole parent frozen base `31bdabc0…`. The replacement was reviewed as a complete candidate from the frozen base, not as a patch over `4bba4b4f…`.

The aggregate digest was recomputed twice by different means: my own Python implementation of the documented algorithm, and the candidate's in-repo guard `TestValidationExecutionBoundaryCandidateMatchesWorkingTree`, which implements it independently in Go. Both produce the pinned value.

Inventory state at the payload commit: `status: review_candidate`, `frozen: false`, `execution_performed: false`, `gates.may_open_validation: false`, `gates.may_open_held_out: false`, 16 boundary files, `run_budget_usd: 300`, `max_cost_usd_per_trial: 0.15`, `timeout_seconds: 180`. `revision_of` names prior payload `4bba4b4f…`, review record `41c75dea…`, and the whitespace finding. The runner's frozen constants (`frozenGraderDigest`, `frozenValidationScheduleDigest`, `frozenValidationManifest`, `frozenValidationRegistryRaw`, `frozenPrivateArchive`) match the inventory contract values exactly.

Live `pre-validation-artifacts.json` in the review checkout remains `status: complete`, `remaining: []`, both gates false, and still describes the accepted authoring runner freeze rather than this replacement. `protocol.json` `gate.may_open_validation` and `gate.may_open_held_out` remain false.

The 16-path payload set and the 16-file inventory set differ, as in the prior round, and the split is again accurate. The payload additionally changes the inventory JSON and three freeze-guard/candidate-guard tests, which sit outside the artifact-set digest and were reviewed as provenance controls. The inventory additionally covers four test files whose blobs are byte-identical to parent `31bdabc0…` (`arm_b_test.go`, `runtime_test.go`, `schedule_test.go`, `scheduled_run_test.go`); I confirmed each equals the base. No private, label, result, runtime-stream, trial, grade, analysis-output, or gate-state path appears anywhere in the payload.

## The closed P1, verified by differential probe

The remediation is a one-token change at `experiments/frontier-v1/runner/validation_gate.go:74`:

```
-	costCap, err := strconv.ParseFloat(strings.TrimSpace(budget), 64)
+	costCap, err := strconv.ParseFloat(budget, 64)
```

`strings` remains imported and used by `normalizeTranche`. The two cap parse sites are now `validation_gate.go:74` and `scheduled_run.go:23` (`strconv.ParseFloat(config.budgetUSD, 64)`), both untrimmed. I traced the string between them: `prepareScheduledCLI` passes `budget` to the preflight, then to `prepareRunConfig`, which stores it as `budgetUSD: budget` with no normalization. `prepareRunConfig`'s `strings.TrimSpace(budget) == ""` is an emptiness test only and does not mutate the stored value.

### Reproduction against the prior payload

I materialized `4bba4b4f…` into the scratch directory with `git archive`, injected the same whitespace case into its five-case test, and ran it. The case failed: the preflight accepted `" 0.15 "` and execution continued past it.

To reach custody, I then ran a disposable public-only probe on that same tree with a fully populated synthetic root — an open synthetic gate, all 120 public validation cases and fixtures copied in, the frozen schedule, and the candidate's own synthetic custodian, which writes a marker file on any invocation:

```
prepareScheduledCLI error = open …\worlds\dev-repo\world.json: The system cannot find the path specified.
PROBE RESULT: custodian WAS contacted (marker present)
```

The prior review's P1 reproduces exactly: `" 0.15 "` passes the preflight, the schedule loads, the custodian `describe` handshake runs, and the run only fails later.

### The same probe against the second payload

Byte-identical probe, byte-identical procedure, `74d06da1…`:

```
prepareScheduledCLI error = validation requires --max-budget-usd exactly 0.15
PROBE RESULT: custodian NOT contacted (marker absent)
```

The cap is rejected at the preflight. Nothing downstream runs.

### Parser agreement across 31 cap spellings

A second disposable probe compared the preflight verdict against a mirror of the execution parser (`ParseFloat`, then `perTrialCap <= 0` rejected) for 31 cap strings, and asserted that any preflight-accepted string must also parse at execution *and* equal `0.15`. No divergence appeared in either direction.

Rejected at the preflight, as required: `" 0.15 "`, `" 0.15"`, `"0.15 "`, `"\t0.15\n"`, `"  0.15  "`, `"\v0.15"`, `"\f0.15"`, `"0.15\r"`, `"\u00a00.15"` (non-breaking space), `""`, `" "`, `"NaN"`, `"nan"`, `"Inf"`, `"+Inf"`, `"-Inf"`, `"-0.15"`, `"abc"`, `"0,15"`, `"0.14"`, `"0.16"`, `"0.1"`.

Accepted at the preflight, and each parses to exactly `0.15` at the execution parser: `"0.15"`, `"0.150"`, `".15"`, `"15e-2"`, `"1.5e-1"`, `"+0.15"`, `"0x1.3333333333333p-3"`, `"1_5e-2"`, `"0.1500000000000000000001"`. These are numeric-equivalent spellings, not changed caps, and match the limitation the prior review already accepted. See finding P2-2 for the one place that equivalence is not fully established.

### Ordering, re-verified

`prepareScheduledCLI` for tranche `validation`:

1. `requireValidationGate`
2. `requireFrozenValidationTrialLimits(budget, timeout)`
3. default output path string only
4. `prepareScheduledInputs` — case load, schedule decode/validate, digest, run-budget parse, then for validation: schedule-digest pin, 300 USD equality, custodian path resolution and `describe`
5. case preflight
6. `requireEmptyScheduledOutput`
7. `prepareRunConfig` — output `MkdirAll`, Phoenix build
8. `claudeDriver.Verify`

The trial-limit guard sits at step 2, ahead of every mutation and every custody contact. `runScheduleCLI` is the only caller of `runScheduledCases`, and it always routes through `prepareScheduledCLI`. Timeout equality is at the `time.Duration` level (`timeout != frozenValidationTimeout`), so `3m` is accepted and `179s`/`181s` are not. The run budget is parsed once at `main.go:235`, compared numerically to `300`, and carried as a `float64` — there is no second parse and therefore no analogous split; a malformed or padded run budget also fails before custodian contact.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Closed gate and tranche boundary | Pass | `normalizeTranche` admits only `authoring` and `validation`. `--tranche validation` requires a schedule; `--write-schedule` is rejected for validation and remains authoring-only. No `held_out` loader, fixture, schedule, grader, or CLI path exists in the runner — the only `held_out` references in non-test runner source are the two lines that require the gate to stay closed. `requireValidationGate` precedes the trial-limit check, custodian contact, output inspection, `MkdirAll`, Phoenix build, and runtime verification, and its error is `validation gate is closed`. Open-gate checks require status `complete`, `may_open_validation: true`, `may_open_held_out: false`, the frozen schedule path and digest, and the source-manifest path and digest, then re-hash the live schedule, manifest, and public registry. |
| 2. Public input and frozen execution contract | Pass | The committed schedule reproduces version 1, tranche `validation`, seed `20260817`, three repetitions, arms A–E, 24 family blocks, 360 pairing keys, 1,800 launches, 120 distinct cases, 360 launches per arm, and canonical digest `sha256:b38a0eaa…e075813` under my independent encoder. `TestFrozenValidationScheduleLoadsOnlyPublicInputs` additionally regenerates it deterministically from the 120 public cases and deep-compares. Manifest and registry pins matched. The cap, timeout, and 300 USD run budget are all locked and all enforced before custody. The whitespace class that failed in the prior round now fails at the preflight. |
| 3. Private-grading custody | Pass, with the recorded TOCTOU limitation | Unchanged from `4bba4b4f…` — `validation_grader.go` is byte-identical between the two payloads, so the prior review's Pass carries. Re-confirmed here: absolute path, `EvalSymlinks` outside the implementation repository, regular file, strict `describe` equality to version 1 / tranche `validation` / 120 cases / grader `sha256:36abfbec…dcfc` / schedule / manifest / registry `sha256:847c510c…c1c816` / private archive `sha256:5a320a87…922e4`; executable hashed, re-hashed after handshake, re-hashed before each grade; grade argv only `grade --case-id <public-id> --trial <retained-trial-path>`. `TestValidationGraderMustRemainExternal` and `TestExternalValidationGraderHandshakeAndGrade` pass. |
| 4. Evidence, stops, and authoring compatibility | Pass | `scheduledRunConfiguration` records the frozen `TimeoutSeconds`, `MaxCostUSDPerTrial` (the parsed `float64`), grader digest, `validation_grader_boundary`, and `validation_grader_adapter_sha256`. Authoring is untouched: `requireFrozenValidationTrialLimits` is called only under `tranche == "validation"`, and `runProbeCLI` never calls it, so authoring retains its configurable cap and timeout. Historical freeze guards still pin the accepted runner freeze at `b4df9190…` and local runtime files at `795ce71a…`. The candidate guard verifies the live 16-file inventory and does not mark it frozen. |
| 5. Reproducible public checks | Pass | See below. The closed-gate CLI probe printed `validation gate is closed`, produced no grader-path error — so the nonexistent custodian path was never resolved — and did not create `experiments/frontier-v1/results/scheduled-validation`. |

## Public checks

From detached `HEAD` `05ab169acdfabae4ce7b42b4ffb8d6ecb276e935`:

| Command | Result |
| --- | --- |
| `git status --porcelain` | empty of tracked changes before and after every check |
| `git diff-tree --no-commit-id --name-status -r 74d06da1…` | the 16 claimed candidate paths |
| `git diff-tree --no-commit-id --name-status -r 4bba4b4f… 74d06da1…` | the 5 claimed changed paths |
| `git diff-tree --no-commit-id --name-status -r 05ab169a…` | only the second replacement report |
| `go test ./experiments/frontier-v1/runner -run TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact -count=1 -v` | pass; all five subtests |
| `go test ./...` | pass (18 packages) |
| nested `experiments/frontier-v1/corpusctl` `go test ./...` | pass (additional; nested module is outside root `./...`) |
| `make quality` | pass; world-build `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| closed-gate CLI probe: `--tranche validation --case all --schedule …/validation.json --validation-grader D:\does-not-exist\grader.exe --max-budget-usd 0.15 --run-budget-usd 300 --timeout 180s` | exit 1; stderr `validation gate is closed`; no grader-path error; `scheduled-validation` absent |
| `git diff --check 31bdabc0… 74d06da1…` | empty |
| `git diff --check 74d06da1… 05ab169a…` | empty |

`make quality` had to be run from PowerShell rather than Git Bash. The `build` target selects a `cmd.exe` idiom (`if not exist bin mkdir bin`) whenever `$(OS)` is `Windows_NT`, which fails when `make` dispatches to `/usr/bin/sh`. This is a property of the reviewer's shell, not a candidate defect, and is unchanged from the frozen base.

`make quality` created gitignored `bin/phoenix` and `build/manifest.json` in the disposable review worktree. They are not candidate evidence, are outside the payload and the artifact set, and were removed afterward. Tracked and ignored status were both empty at the end.

## Findings

Neither finding blocks `ACCEPT`. Both are recorded for the import/refreeze chair.

### [P2-1] The report overstates what its own regression test proves — `docs/reviews/2026-08-21-protocol-v4-validation-execution-boundary-second-revision-candidate.md`

The report states: "The synthetic custodian writes a marker on any invocation. The new case therefore tests the exact custody-contact regression, not only the returned error."

The first sentence is true. The conclusion does not follow. `TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact` calls `prepareScheduledCLI` with `ids = nil`, and `loadSelectedCasesForTranche` returns an empty slice for a nil id list. Schedule validation therefore fails at `generateScheduleForTranche` with `schedule requires cases, repetitions, and at least two arms` before the custodian handshake or `MkdirAll` can be reached — regardless of whether the cap guard fired. Both structural assertions in that test, the output-directory check and the custodian-marker check, are consequently vacuous for all five cases.

I established this directly rather than by inspection. Running the whitespace case against the vulnerable payload `4bba4b4f…`, the test failed on the error-message assertion at `validation_execution_test.go:108` and the marker was absent — precisely the reverse of what the report's sentence implies. Only the separately populated probe described above set the marker.

The test is nonetheless a sound regression test, which is why this is P2 and not P1. Its error-message assertion is load-bearing and does the work the report attributes to the marker: it fails on the vulnerable payload, and a mutation moving `requireFrozenValidationTrialLimits` to after `prepareScheduledInputs` — that is, after the custodian handshake — fails all five subtests. The guard's ordering is genuinely defended.

Recommendation for a later revision, not required for import: pass the 120 case IDs and copy the public cases and fixtures into the synthetic root, as the review probes do, so the marker and output assertions become load-bearing rather than incidentally satisfied. Alternatively, correct the report's sentence to claim only what the test establishes.

### [P2-2] The unmodified cap string is forwarded verbatim to a third, external parser — `experiments/frontier-v1/runner/runtime.go:120`

`prepareRuntimeInvocation` passes `request.BudgetUSD` — the same raw operator string, not the parsed `float64` — straight into the runtime argv as `--max-budget-usd`. The cap therefore has three consumers, not two: the Go preflight, the Go execution parser, and the external `claude` executable's own parser, whose numeric grammar is outside this repository and outside the frozen contract.

The two Go parsers now agree exactly, which is what the prior P1 demanded. But the spellings both Go parsers accept include forms no non-Go parser is likely to accept: `0x1.3333333333333p-3`, `1_5e-2`, and to a lesser degree `.15`, `15e-2`, and `+0.15`. For those, the runner would certify the cap as frozen at `0.15` and then hand the runtime a string it may reject or interpret differently. This is the same class as the closed P1 — a normalization mismatch across a parse boundary — displaced one hop outward.

It is P2, not P1, for three reasons. The exposure requires an operator to type an exotic spelling deliberately; the frozen operator instruction is `--max-budget-usd 0.15`, so no realistic input reaches it. The guard defends operator discipline over a frozen contract, not an adversary who already controls the CLI. And the failure mode is loud rather than silent: a runtime that rejects the value fails the trial into the existing safety-stop and `indeterminate` paths, while summary evidence independently records `MaxCostUSDPerTrial` as the parsed `float64`, so any divergence is auditable after the fact.

The condition also pre-exists this payload. `runtime.go` is byte-identical between `4bba4b4f…` and `74d06da1…` and unchanged from the frozen base in this respect, so it is outside what the prior `REVISE` asked this revision to fix.

Recommendation for the import/refreeze chair, not required for import: either require the exact string `"0.15"` for validation, or forward `strconv.FormatFloat(costCap, 'f', -1, 64)` to the runtime instead of the operator's raw string, so all three consumers see one canonical spelling.

## Accepted limitations

- Same-machine fresh session, recorded above.
- Hash-then-exec TOCTOU on the custodian path. The runner hashes path contents and then execs that same path string; repeated identity checks after `describe` and before every grade detect substitution across the run but cannot bind one invocation to the hashed bytes against a hostile co-resident writer. For this threat model — operator-supplied custodian, prevent in-repo label loading, detect adapter replacement across the run, no private material in the implementation workspace — those checks are sufficient. Carried forward unchanged from the prior review.
- `--runtime` remains an operator path, as in authoring. The frozen invocation template is still built in `prepareRuntimeInvocation` and recorded in summary configuration.
- Numeric-equivalent spellings of `0.15` pass both Go parsers and change neither pairing-budget arithmetic nor summary `float64` evidence. They are not changed caps under the stated numeric rule. The set is now slightly wider than the prior review enumerated, since removing `TrimSpace` did not narrow the accepted grammar: `+0.15` and `1_5e-2` also parse. See P2-2 for the one consequence that is not fully closed.
- Independent reconstruction used public cases, fixtures, manifest, and registry only. No label file or private archive was opened. The private-archive digest is checked only as a custodian `describe` string.
- Root `go test ./...` does not enter the nested `corpusctl` or `surface-spike` modules. Nested `corpusctl` tests were run separately and pass; `make quality` still ran `validate-authoring` and spec validation.
- `make quality` was run from PowerShell for the shell-dialect reason recorded above.

## Smallest next artifact

A project-chair import/refreeze decision that imports `74d06da13f62cffa4fd635e048e6331a6e2d95a6` onto `main`, refreezes the 16-file boundary at `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1`, and updates `pre-validation-artifacts.json` to describe the replacement runner freeze. That decision should record P2-1 and P2-2 as known, non-blocking residuals. This review does not specify that change and does not authorize it.

Gate 1A opening remains a separate decision and must leave held-out closed.
