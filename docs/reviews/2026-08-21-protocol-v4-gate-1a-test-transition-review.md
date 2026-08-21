# Protocol-v4 Gate 1A test-transition independent review

- **Reviewer role:** `gate_1a.test_transition_independent_reviewer`. Distinct from `gate_1a.test_transition_author`, the earlier execution-boundary reviewers, the import/refreeze chair of decision 0010, and any later Gate 1A chair.
- **Date:** 2026-08-21
- **Verdict:** `ACCEPT`
- **Review assignment:** conversational, from the project chair via `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-review-prompt.md`. No `docs: assign …` commit exists for this round; `main` was at `94761cf43b938f643173b486c7f1fca5149e8a08` when the review opened and was not advanced before this record.
- **Frozen implementation base:** `94761cf43b938f643173b486c7f1fca5149e8a08`
- **Prior accepted payload:** `74d06da13f62cffa4fd635e048e6331a6e2d95a6`
- **Prior refreeze decision:** `docs/decisions/0010-protocol-v4-validation-execution-boundary-import-refreeze.md` at `94761cf43b938f643173b486c7f1fca5149e8a08`; independently verified 5,869 raw bytes; `sha256:5b7ff5965c5f8c3a26f62c3c0e3643539b16a5475888a4a090c7b9071a550c28`
- **Replacement payload:** `71f9648789decf4cd56ef8a24bc840b0dda7efd9` (`experiments: decouple gate tests from chair state`)
- **Replacement-report commit:** `7ffc4a1966b4127d9445ebc840723dc6a7a245d4` (`docs: report Gate 1A test transition candidate`)
- **Candidate report:** `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-candidate.md`; independently verified 4,928 raw bytes; `sha256:689ad8cee3a6f5af2daed05081e10d1263b469086555d63757a7250a368dbc3f`
- **Review prompt:** `docs/reviews/2026-08-21-protocol-v4-gate-1a-test-transition-review-prompt.md`; independently verified 5,894 raw bytes; `sha256:8c48f9e3bcb8d213324ebf1b2f8b943499c5dde271a44c6e91d71f850b8d467e`
- **Candidate inventory:** `experiments/frontier-v1/artifacts/gate-1a-test-transition-candidate.json` at the replacement payload; independently verified 4,118 raw bytes; `sha256:ca7eae05c72f970a592e4a82ab88d2167564323f08a788a0081a349f32c49d75`
- **Replacement 16-file identity:** `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`
- **Review checkout:** `D:\Work\personal\phoenix-gate-1a-test-transition-review` (new detached worktree at the report commit; did not previously exist; `core.autocrlf=false`, `core.eol=lf`; `git status --porcelain` empty of tracked changes before and after every check)
- **Disposable reconstruction and probes:** `D:\Work\personal\phoenix-gate-1a-test-transition-review-base-probe` (detached at the frozen base), `D:\Work\personal\phoenix-gate-1a-test-transition-review-open-gate-probe` (detached at the payload), and `D:\Work\personal\phoenix-gate-1a-test-transition-review-scratch` (outside every Phoenix checkout). Probe edits were restored or discarded and were not committed.

This review does not import the replacement, modify a frozen runner or grader, refreeze replacement bytes, open Gate 1A on `main`, authorize validation, or authorize held-out. Both gates on `main` remain false. No model, arm, validation trial, private grade, outcome analysis, or prospective evaluation was run.

## Independence declaration

This session did not author the candidate payload, the candidate report, the review prompt, Phoenix, decision 0010, or any freeze/import decision. It is distinct from the execution-boundary author and reviewer roles recorded through `6b158ea0…`.

Recorded limitations:

- The reviewer is a fresh session on the same machine as earlier Phoenix work. The same-machine / shared-workspace limitation is recorded and does not reuse a prohibited role.
- The implementation workspace `D:\Work\personal\phoenix` at `main` `94761cf43b938f643173b486c7f1fca5149e8a08` was inspected only to confirm it was not modified. Every technical audit, reconstruction, test, quality run, blocker reproduction, mutation, and disposable probe ran in the detached review worktree or the disposable probe trees.
- The candidate branch `codex/gate-1a-test-transition-candidate` and its author worktree `D:\Work\personal\phoenix-gate-1a-test-transition-candidate` were not modified. `git status --porcelain` there was empty before and after.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. `experiments/frontier-v1/labels/` in the review checkout contains only `authoring/`; no validation or held-out label directory exists and none was created or opened. Existing authoring tests and `make quality`'s authoring validator may read authoring labels as they already do. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. No real custodian executable was substituted. `seal --write` was not run.

## Verdict

Payload `71f9648` safely makes the accepted closed-gate tests independent of the live chair-controlled gate state. It does not change production behavior. Both safety assertions remain load-bearing. The exact four-path replacement may be imported and refrozen before Gate 1A opens.

I reproduced the stated blocker against the frozen base: with only the chair-controlled validation gate set true, held-out kept false, and the two non-frozen freeze-state assertions adjusted, `TestRepositoryValidationGateRemainsClosed` received no error from `requireValidationGate`, and `TestClosedValidationGatePrecedesOutputMutation` advanced to the deliberately invalid one-nanosecond timeout and returned `validation requires --timeout exactly 180s`.

The replacement moves those two assertions onto a synthetic temporary repository whose gate is explicitly false. The committed tests pass with the live gate closed. A disposable open-gate `make quality` probe also passes. Production files in the 16-file boundary are byte-identical to accepted payload `74d06da1…` except for `validation_execution_test.go`.

Three P2 findings are recorded below. Two are carried forward unchanged from decision 0010. The third is a fixture limitation of the new closed-gate helper: flipping only its boolean does not isolate the timeout error, because the helper omits frozen schedule identity fields. The committed closed-path assertion is nonetheless load-bearing, and a populated dual mutation on the existing open-gate helper exposes the timeout error. None of these findings blocks `ACCEPT`.

This `ACCEPT` authorizes only a later project-chair import/refreeze of payload `71f9648789decf4cd56ef8a24bc840b0dda7efd9` and boundary identity `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`. It does not import the replacement, refreeze bytes, open Gate 1A, authorize validation execution, authorize use of a real custodian, or authorize any outcome access. Gate 1A opening remains a separate chair decision and must keep held-out closed.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Payload vs parent `94761cf43b938f…` | `git diff-tree --name-status -r 71f9648789de…` | exactly four paths: candidate inventory (added), `artifact_freeze_test.go`, `validation_candidate_artifact_test.go`, `validation_execution_test.go` |
| Payload parents | `git cat-file -p 71f9648789de…` | sole parent `94761cf43b938f643173b486c7f1fca5149e8a08` |
| Report commit tree vs payload | `git diff-tree --name-status -r 7ffc4a1966b4…` | exactly the candidate report and this review prompt |
| Candidate report at `7ffc4a1966b4…` | SHA-256 of git blob; 4,928 bytes; LF-only | `sha256:689ad8cee3a6f5af2daed05081e10d1263b469086555d63757a7250a368dbc3f` |
| Review prompt at `7ffc4a1966b4…` | SHA-256 of git blob; 5,894 bytes; LF-only | `sha256:8c48f9e3bcb8d213324ebf1b2f8b943499c5dde271a44c6e91d71f850b8d467e` |
| Candidate inventory at `71f9648789de…` | SHA-256 of git blob; 4,118 bytes; LF-only | `sha256:ca7eae05c72f970a592e4a82ab88d2167564323f08a788a0081a349f32c49d75` |
| Inventory file digests and ordinal-path aggregate | LF-normalized SHA-256 of each of the 16 inventory paths, then SHA-256 of ordinal-sorted `path\tsha256:<hex>\n` lines, in an independent Python reimplementation | all 16 file digests matched the inventory; aggregate `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9` |
| 16-file set vs accepted payload `74d06da13f62…` | byte comparison of git blobs | only `experiments/frontier-v1/runner/validation_execution_test.go` differs; its LF-normalized digest changes from `sha256:212a0b78ddf29c6e1d7d8d4ab305e55707b86619fd1d96be6ee9ac9509d3a61d` to `sha256:460dc582292473ecd5b91e436ba726e2f409399cd254588ee8f0eded99a02b03` |
| Decision 0010 at `94761cf43b93…` | SHA-256 of git blob; 5,869 bytes; LF-only | `sha256:5b7ff5965c5f8c3a26f62c3c0e3643539b16a5475888a4a090c7b9071a550c28` |
| World-build digest from both `make quality` runs | `cmd/build-manifest` then `quality-check compare` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: detached `HEAD` `7ffc4a1966b4…`; sole parent payload `71f9648789de…`; that commit's sole parent frozen base `94761cf43b93…`. The replacement was reviewed as a complete candidate from the frozen base.

The aggregate digest was recomputed twice by different means: my own Python implementation of the documented algorithm, and the candidate's in-repo guard `TestGate1ATestTransitionCandidateMatchesWorkingTree`, which implements it independently in Go. Both produce the pinned value. Worktree bytes of the inventory and of all 16 files matched the git blobs.

Inventory state at the payload commit: `status: review_candidate`, `frozen: false`, `execution_performed: false`, `gates.may_open_validation: false`, `gates.may_open_held_out: false`, 16 boundary files, `run_budget_usd: 300`, `max_cost_usd_per_trial: 0.15`, `timeout_seconds: 180`. `replaces.accepted_candidate_commit` names prior payload `74d06da13f62…`.

Live `pre-validation-artifacts.json` in the review checkout remains `status: complete`, `remaining: []`, both gates false, and still describes the accepted execution-boundary freeze at `74d06da13f62…` / identity `sha256:3056e95a696bb7fbe6a8e0aec9df960256de085490ebbd853c57353a5d9383f1`. `protocol.json` `gate.may_open_validation` and `gate.may_open_held_out` remain false.

The 16-path payload set and the 16-file inventory set differ, and the split is accurate. The payload additionally changes the inventory JSON, the live-candidate guard, and the historical freeze guard, which sit outside the artifact-set digest and were reviewed as provenance controls. The historical freeze guard now verifies the accepted 16-file set at commit `74d06da13f62…` rather than from the live working tree, which is required once `validation_execution_test.go` is allowed to change. No private, label, result, runtime-stream, trial, grade, analysis-output, or gate-opening path appears in the payload.

## Required-check matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Ancestry and four-path payload boundary | Pass | Sole parent is the frozen base. Changed paths are exactly the inventory, `artifact_freeze_test.go`, `validation_candidate_artifact_test.go`, and `validation_execution_test.go`. |
| 2. Inventory, 16-file, and aggregate identities | Pass | 4,118 raw bytes; inventory `sha256:ca7eae05…c49d75`; all 16 LF-normalized file digests match; aggregate `sha256:9f8c48fb…de2d9`. |
| 3. Only `validation_execution_test.go` changes inside the 16-file boundary | Pass | The other 15 inventory files are byte-identical to `74d06da13f62…`. No non-test runner source, protocol, schedule, prompt, schema, Arm B, world, grader, analysis, corpus, registry, or manifest byte changes. |
| 4. Reproduce the blocker against the base | Pass | See below. The two live-repository closed-gate tests fail for the stated reasons. |
| 5. Replacement uses a synthetic closed-gate repository | Pass | Both tests call `closedValidationGateTestRoot`, which writes `may_open_validation: false` and `may_open_held_out: false` into a temporary root. |
| 6. Closed-gate error assertion remains load-bearing | Pass | Mutating the synthetic gate to true makes `TestClosedValidationGateIsRejected` fail. |
| 7. Ordering assertion remains load-bearing | Pass, with P2-3 | Closed path returns `validation gate is closed` and leaves the output path absent. A one-bit helper mutation does not isolate the timeout; a populated dual does. See findings. |
| 8. Closed-gate runner suite and `make quality` | Pass | Review checkout, both committed gates false. |
| 9. Open-gate `make quality` probe | Pass | Disposable payload checkout; only the chair-controlled validation gate set true; held-out false; two non-frozen freeze-state assertions updated; probe restored afterward. |
| 10. Production behavior byte-identical to the accepted boundary | Pass | Unchanged production files retain the numeric 0.15 USD cap check, exact 180-second timeout, 300 USD run budget, pre-output and pre-custodian ordering, frozen schedule and public identities, A–E launch contract, external grading boundary, authoring compatibility, and held-out unreachability. |
| 11. Public, outcome-free inputs only | Pass | No real custodian, private label, grade, validation result, held-out result, or outcome was accessed. |
| 12. Carry forward decision 0010 P2 residuals | Pass | P2-1 and P2-2 remain non-blocking and are not upgraded or discarded. P2-3 is new evidence about the replacement helper, not a change to those residuals. |

## Blocker reproduction against the base

Disposable checkout `D:\Work\personal\phoenix-gate-1a-test-transition-review-base-probe` at `94761cf43b93…`. I set only `experiments/frontier-v1/pre-validation-artifacts.json` `gates.may_open_validation` to true and kept `may_open_held_out` false. `protocol.json` was not modified. To reach the runner suite I adjusted only the two non-frozen live freeze-state assertions:

- `TestPreValidationFreezeMatchesAcceptedCandidates` in `artifact_freeze_test.go`
- `TestCompletedArtifactFreezeKeepsOutcomeGatesClosed` in `local_artifact_candidate_test.go`

The two live-repository closed-gate tests then failed exactly as the candidate report states:

```
=== RUN   TestRepositoryValidationGateRemainsClosed
    validation_execution_test.go:48: gate error = <nil>
--- FAIL: TestRepositoryValidationGateRemainsClosed (0.03s)
=== RUN   TestClosedValidationGatePrecedesOutputMutation
    validation_execution_test.go:59: gate error = validation requires --timeout exactly 180s
--- FAIL: TestClosedValidationGatePrecedesOutputMutation (0.03s)
```

Those probe edits were restored and not committed.

## Replacement and mutation results

`closedValidationGateTestRoot` at `validation_execution_test.go:226-237` writes a synthetic `pre-validation-artifacts.json` with `status: complete`, `may_open_validation: false`, and `may_open_held_out: false`. `TestClosedValidationGateIsRejected` (`:45-50`) and `TestClosedValidationGatePrecedesOutputMutation` (`:52-63`) both use that root. The ordering test still passes timeout `1` (one nanosecond) into `prepareScheduledCLI`, whose validation path checks the gate before `requireFrozenValidationTrialLimits`.

Committed closed-gate tests, review checkout, gates false:

```
=== RUN   TestClosedValidationGateIsRejected
--- PASS: TestClosedValidationGateIsRejected (0.01s)
=== RUN   TestClosedValidationGatePrecedesOutputMutation
--- PASS: TestClosedValidationGatePrecedesOutputMutation (0.00s)
```

The ordering test's `os.Stat` assertion confirmed the synthetic output path `results/must-not-exist-gate-test` was absent.

One-bit mutation of `closedValidationGateTestRoot` to `may_open_validation: true`, held-out still false:

```
=== RUN   TestClosedValidationGateIsRejected
    validation_execution_test.go:48: gate error = validation gate does not name the frozen validation schedule
--- FAIL: TestClosedValidationGateIsRejected (0.01s)
=== RUN   TestClosedValidationGatePrecedesOutputMutation
    validation_execution_test.go:58: gate error = validation gate does not name the frozen validation schedule
--- FAIL: TestClosedValidationGatePrecedesOutputMutation (0.00s)
```

Check 6 is satisfied: the closed-gate error assertion is load-bearing. Check 7's dual is not isolated by that one-bit flip, because the closed helper omits the frozen schedule identity fields that `requireValidationGate` demands once the gate is open. See P2-3.

Populated dual: the same ordering test pointed at `openValidationGateTestRoot`, which already copies the public schedule, manifest, and registry and sets the synthetic validation gate true while keeping held-out false:

```
=== RUN   TestClosedValidationGatePrecedesOutputMutation
    validation_execution_test.go:58: gate error = validation requires --timeout exactly 180s
--- FAIL: TestClosedValidationGatePrecedesOutputMutation (0.05s)
```

That is the later timeout error required by check 7. `prepareScheduledCLI` still returns before `requireEmptyScheduledOutput` and `prepareRunConfig`, so the output directory is not created. Both mutation edits were restored and not committed.

## Public checks

From detached `HEAD` `7ffc4a1966b4127d9445ebc840723dc6a7a245d4`, committed gates false:

| Command | Result |
| --- | --- |
| `git status --porcelain` | empty of tracked changes before and after every check |
| `git diff-tree --no-commit-id --name-status -r 71f9648789de…` | the four claimed candidate paths |
| `git diff-tree --no-commit-id --name-status -r 7ffc4a1966b4…` | only the candidate report and review prompt |
| `go test ./experiments/frontier-v1/runner -run TestClosedValidationGateIsRejected\|TestClosedValidationGatePrecedesOutputMutation -count=1 -v` | pass |
| `go test ./experiments/frontier-v1/runner -count=1` | pass, as part of `make quality` (5.651s) |
| `go test ./...` | pass (18 root packages) |
| nested `experiments/frontier-v1/corpusctl` `go test ./...` | pass |
| `make quality` | pass; world-build `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| disposable open-gate `make quality` | pass; same world-build digest; runner tests 7.121s |
| `git diff --check 94761cf43b93… 71f9648789de…` | empty |

`make quality` had to be run from PowerShell rather than Git Bash. The `build` target selects a `cmd.exe` idiom (`if not exist bin mkdir bin`) whenever `$(OS)` is `Windows_NT`, which fails when `make` dispatches to `/usr/bin/sh`. This is a property of the reviewer's shell, not a candidate defect, and is unchanged from the frozen base.

`make quality` created gitignored `bin/phoenix` and `build/manifest.json` in the disposable review and probe worktrees. They are not candidate evidence, are outside the payload and the artifact set, and were removed afterward. Tracked and ignored status were both empty at the end of the review checkout.

Live committed gate state after every check: `pre-validation-artifacts.json` `may_open_validation: false`, `may_open_held_out: false`; candidate inventory both false; `protocol.json` both false. Held-out was never opened.

## Findings

None of these findings blocks `ACCEPT`. P2-1 and P2-2 are carried forward from decision 0010. P2-3 is new evidence from this replacement.

### [P2-1] The prior candidate report overstates what its own regression test proves — `experiments/frontier-v1/runner/validation_execution_test.go:80-117`

Carried forward unchanged from decision 0010. `TestValidationRejectsChangedTrialLimitsBeforeOutputOrCustodianContact` still supplies `ids = nil`. Schedule validation therefore still prevents later side effects, so the marker and output assertions remain incidentally satisfied. The error assertion and guard-order mutation remain load-bearing. This payload does not claim to fix that residual, and the trial-limit test body is unchanged.

### [P2-2] The unmodified cap string is forwarded verbatim to a third, external parser — `experiments/frontier-v1/runner/runtime.go:120`

Carried forward unchanged from decision 0010. `runtime.go` is byte-identical to accepted payload `74d06da13f62…`. Numeric-equivalent Go spellings of `0.15` are still forwarded to the external runtime parser. The frozen operator spelling `0.15` is unaffected.

### [P2-3] A one-bit mutation of the new closed-gate helper does not isolate the timeout error — `experiments/frontier-v1/runner/validation_execution_test.go:226-237`

`closedValidationGateTestRoot` writes only `status` and the two gate booleans. That is enough for the closed path, because `requireValidationGate` returns `validation gate is closed` before it inspects schedule identity. Once the boolean is flipped to true, the next failure is `validation gate does not name the frozen validation schedule`, not `--timeout exactly 180s`.

The ordering assertion is nevertheless load-bearing on the committed closed path: if `requireValidationGate` were removed or moved after `requireFrozenValidationTrialLimits`, `TestClosedValidationGatePrecedesOutputMutation` would fail today with the timeout error while the synthetic gate was still false. Independently, pointing that same test at `openValidationGateTestRoot` — a synthetic root that can pass the open-gate identity checks — exposes `validation requires --timeout exactly 180s`.

This is the same class as P2-1: a side mutation is weaker than the report's dual suggests, while the error assertion still does the required work. It is P2, not P1, because production ordering is unchanged, the closed-path test fails if ordering is inverted, and the populated dual confirms the timeout is the next check.

Recommendation for a later replacement, not required for import: include the frozen schedule identity fields in the closed helper, or have the ordering test share the open helper's public files, so that flipping only `may_open_validation` isolates the timeout error.

## Accepted limitations

- Same-machine fresh session, recorded above.
- Hash-then-exec TOCTOU on the custodian path, carried forward from the execution-boundary review. This payload does not touch `validation_grader.go`.
- `--runtime` remains an operator path, as in authoring.
- Numeric-equivalent spellings of `0.15` remain accepted by both Go parsers. See P2-2.
- Independent reconstruction used public cases, fixtures, manifest, and registry only. No label file or private archive was opened. The private-archive digest is checked only as a custodian `describe` string in already-accepted tests.
- Root `go test ./...` does not enter the nested `corpusctl` or `surface-spike` modules. Nested `corpusctl` tests were run separately and pass; `make quality` still ran `validate-authoring` and spec validation.
- `make quality` was run from PowerShell for the shell-dialect reason recorded above.

## Smallest next artifact

A project-chair import/refreeze decision that imports `71f9648789decf4cd56ef8a24bc840b0dda7efd9` onto `main`, refreezes the 16-file boundary at `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`, and updates the historical freeze guard to describe this replacement. That decision should record P2-1, P2-2, and P2-3 as known, non-blocking residuals. This review does not specify that change and does not authorize it.

Gate 1A opening remains a separate decision and must leave held-out closed.
