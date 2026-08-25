# Protocol-v4 Gate 1A complexity-refactor independent review

- **Reviewer role:** `gate_1a.complexity_refactor_independent_reviewer`. Distinct from the authoring agent that produced payload `8ed9202c80d8f591c5d0db8a7e8a349022f952f8`, the chair that recorded decisions 0016/0017, and every earlier Gate 1A authoring, review, import, or refreeze role.
- **Date:** 2026-08-25
- **Verdict:** `ACCEPT`
- **Review assignment:** conversational, from the project chair after the payload was committed. The reviewed object is the exact commit `8ed9202c80d8f591c5d0db8a7e8a349022f952f8` (`experiments: author Gate 1A complexity refactor`), verified byte-identical to the working tree at review time (`git diff 8ed9202 --stat` empty, `git status` clean before and after all checks, HEAD = payload).
- **Frozen implementation base / parent:** `8881d5f552aedfd9e6283ee1cc33e468fe036856` (`experiments: close protocol v4 validation gate for complexity refactor`, the decision-0017 closure commit). Confirmed as the payload's sole parent and as the candidate's recorded `base_commit`.
- **Authorizing residual:** decisions 0014 and 0015 recorded the `gocyclo -over 15` quality-gate failures (`runScheduledCases` 27, `loadScheduledResume` 24, `activate.New` 16, plus the freeze-test helper) as non-blocking residuals. Decision 0017 (`docs/decisions/0017-protocol-v4-gate-1a-complexity-refactor-closure.md`) closed the gate reopened by the withdrawn 0016 so this refactor could be authored, reviewed, and refrozen with both gates closed; no execution occurred under 0016.
- **Payload:** `8ed9202c80d8f591c5d0db8a7e8a349022f952f8`
- **Candidate artifact:** `experiments/frontier-v1/artifacts/gate-1a-complexity-refactor-candidate.json` — raw SHA-256 `67ce91a3a1fb991c2fc2565f2d6ca35f0bd7a352daa40dce3ab629f9bea8ac34`, 5,155 bytes (file is LF-only; the LF-normalized digest is identical).
- **Candidate note at the payload:** `experiments/frontier-v1/artifacts/complexity-refactor-candidate.md`
- **Review checkout:** `D:\Work\personal\phoenix` on `main` at the payload commit with a clean tree. Same-machine / shared-workspace limitation recorded: no disposable detached worktree; byte identity was instead established by empty `git diff 8ed9202` and clean `git status` before and after all checks, including after removal of the transient build outputs.

This review does not import the payload, modify freeze hashes, refreeze bytes, open Gate 1A, authorize validation, authorize held-out, resume or splice `experiments/frontier-v1/results/scheduled-validation`, inspect private or held-out labels, or run a model. Both gates remain false. The two expected freeze-digest test failures against the old frozen identity are not defects against this candidate.

## Independence declaration

This is a fresh session that did not author the candidate payload, the candidate note, the refactored code, decision 0016 or 0017, or any freeze/import decision. It is an independent reviewer session.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. No validation or held-out labels were inspected. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. Nothing under `experiments/frontier-v1/results/` was read or written. No custodian or grader was contacted. No gate state was mutated. The only build artifacts produced were `go build`/`go test` outputs and one manual reproduction build (`bin/phoenix`, `build/manifest-review.json`), all removed; `git status` was empty afterward.

## Verdict

The payload delivers exactly the contracted behavior-preserving complexity decompositions and nothing else. Every digest in the 20-file candidate inventory matches the working tree, the set digest reproduces, the diff against the base commit touches only the nine claimed files, the frozen linux/amd64 recipe was independently rebuilt on this host and reproduces the claimed NEW world-build digest, `gocyclo -over 15` is now clean across the quality-gate scope, and the full suite fails only in the two documented pre-refreeze tests. Line-by-line comparison of all three decompositions found no semantic drift: no error message, error ordering, stop rule, checkpoint rule, refusal condition, or evidence byte changes. There are no P0 or P1 findings. This `ACCEPT` authorizes only a later project-chair import/refreeze of payload `8ed9202c80d8f591c5d0db8a7e8a349022f952f8`. It does not freeze, open validation, or authorize any execution.

## Required-check matrix

| Check | Result |
| --- | --- |
| A. Byte identity: all 20 candidate file digests match working tree under LF normalization (CRLF→LF, CR→LF, SHA-256) | Pass (20/20 match; 0 mismatches) |
| B. Artifact set digest over ordinal-path-sorted `<path>\t<digest>\n` lines | Pass — computed `sha256:2d20be1ac41cc9072ae472b9b5a73e2607e9d35093fda49c422fa487117ecb4b`, equals claimed |
| C. `base_commit` = parent of payload (`8881d5f552aedfd9e6283ee1cc33e468fe036856`); `git diff 8881d5f..8ed9202 --stat` touches only the nine claimed files (+287/−180) | Pass |
| D. Set delta vs the frozen 20-file `scheduled_runner` inventory (decision 0015, set digest `sha256:c9e2f4ff…f447f`): identical 20 paths; digests changed in exactly README.md, validation-execution-boundary.md, scheduled_run.go, scheduled_resume.go, validation_build_recipe_test.go | Pass |
| E. No gate/freeze tampering: pre-validation-artifacts.json, protocol.json, schedules, manifests, labels, grader files, artifact_freeze_test.go, local_artifact_candidate_test.go, world-build.linux-amd64.json, pre-validation-local-candidate.json, results/ all untouched; both gates false in the gate document and in the candidate JSON | Pass |
| F. Contract cross-check vs freeze: schedule `sha256:b38a0eaab063…075813`, manifest `sha256:57ccc0c754f7…fd490`, replaced world-build `sha256:bf976cad…be8a4e` (= frozen doc line 273), cap 0.15, timeout 180, cases 120, launches 1800, 300 USD budget enforced at frozen main.go:252 | Pass |
| G. World-build reproduction, manual: frozen recipe (`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.26.6`, `-trimpath -buildvcs=false -ldflags "-s -w -buildid= -X main.version=dev"`, `-o bin/phoenix ./cmd/phoenix`) + `cmd/build-manifest` with the frozen input set printed `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4` — equals the claimed NEW digest | Pass |
| H. World-build reproduction, end-to-end: `TestValidationBuildReproducesFrozenWorldBuildDigest` fails with "live world-build sha256:b5a26d5e…19c4 does not reproduce frozen world-build sha256:bf976cad…8a4e" — live value equals the claim, failure is the documented interim state | Pass |
| I. `gocyclo -over 15 cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`: exit 0, no output — the point of the cycle. Measured: runScheduledCases 12 (was 27), applyScheduledResume 6, runRemainingSchedule 9, loadScheduledResume 11 (was 24), classifyScheduledOutputEntries 15, activate.New 9 (was 16), validateActivationRule 9 | Pass |
| J. `go run ./cmd/quality-check no-output -- go run github.com/mibk/dupl@v1.1.0 -plumbing -t 100 …`: exit 0 | Pass |
| K. Hygiene: `gofmt -l` empty, `go vet ./...` clean, staticcheck v0.7.0 clean | Pass |
| L. Full suite `go test -count=1 ./...`: only `TestPreValidationFreezeMatchesAcceptedCandidates` (README.md `sha256:af321a8e…` vs old frozen `sha256:a7650caf…`) and `TestValidationBuildReproducesFrozenWorldBuildDigest` (live b5a26d5e vs old frozen bf976cad) fail; every other package ok | Pass (documented pre-refreeze state; no other failures) |
| M. Behavior contract tests: all 22 scheduled-run/scheduled-resume tests pass (stops, checkpoints, resume refusals, retries); internal/activate 5/5 pass | Pass |
| N. Candidate guards: 5/5 pass — new working-tree guard reads base_commit `8881d5f…6856` with 20 files; the superseded build-recipe guard now verifies historically at `ed3860708c931ddc848b1bc90e4d6435585ff0d6`; test-transition (71f9648) and authoring-repair (610fa58) historical guards unaffected | Pass |
| O. Public-only inputs; `execution_performed=false` truthful; no results/, custodian, or model touched by any check run here | Pass |

## Behavior equivalence (core review)

Each extraction was compared line-by-line between `git show 8881d5f:<path>` and the working tree.

**1. `runScheduledCases` → `applyScheduledResume` + `runRemainingSchedule` (runner/scheduled_run.go).**
- *Replay and resumed stops (`applyScheduledResume`):* replays `resume.results` first, exactly as before. Resumed budget stop: old did record-then-`return summary, nil` with error escalation; new does `return true, recordBudgetStop(…)` and the caller checks `err` before `done`, so both the success path (summary returned) and the error path (`scheduledSummary{}, err`) are byte-identical in effect. Resumed safety stop with remaining entries: same argument for `recordSafetyStop(…, fmt.Errorf("resumed incomplete safety stop"))`. Resumed safety stop with nothing remaining: identical in-place mutation (`Status="indeterminate"`, StopReason set only if empty). The no-stop case returns `resume.nextIndex >= len(schedule.Entries)`, matching the old early return that preceded the initial checkpoint write; the only overlap (empty schedule, nextIndex 0) returned before the checkpoint write in the old code too, so the initial-checkpoint ordering is preserved exactly.
- *Group loop (`runRemainingSchedule`):* recomputes `groupSize = len(schedule.Arms)` — same expression the caller validated non-zero on the same immutable by-value `schedule`, so the value is provably equal. Skip rule (`start+groupSize <= nextIndex`), capacity gate applied only when `start >= nextIndex`, capacity formula (`perTrialCap * groupSize * 1.10` vs `runBudgetUSD + 1e-9`), and partial-group entry skip (`start+offset < nextIndex`) are verbatim. Budget stop: old was record-then-`break`, which skipped the end-of-body checkpoint write and fell through to `return summary, nil`; new is `return recordBudgetStop(…)` — same records, same statuses, and critically **no checkpoint write after a budget stop** in either version. Safety stop: old record-then-`return summary, nil`; new `return recordSafetyStop(…)` — same, **no checkpoint after a safety stop**. Checkpoint after each completed group passes `*summary` (a value snapshot), identical to the old by-value pass, with the same `(start+groupSize, (start+groupSize)/groupSize)` arguments. All error propagation collapses to `scheduledSummary{}, err` at the caller, as before.
- All other functions in the file (`recordBudgetStop`, `recordSafetyStop`, `replayAssignedResult`, `addAssignedResult`, `runAssignedTrial`, configuration, stems, rounding) are textually unchanged.

**2. `loadScheduledResume` → `classifyScheduledOutputEntries` (runner/scheduled_resume.go).**
The extraction lifts the directory-classification block verbatim into a helper returning `{assignments, hasCheckpoint}`. Preserved exactly, in the same order, with the same error strings: directory-entry refusal ("unrecognized files"), finished-summary refusal, checkpoint detection, assignment loading with duplicate-launch-index refusal and stem registration, evidence-name classification, default unrecognized-file refusal, orphan-evidence refusal, and the missing-checkpoint refusal **still conditional on `len(assignments) > 0 || len(evidence) > 0`** — a checkpoint-less directory containing neither still loads clean. The post-classification sequence (reconstruct → verifyResumeWorldBuild → verifyScheduledCheckpoint only if `hasCheckpoint`) is unchanged, as are `loadAssignmentFile`, `reconstructScheduledResume`, and both verifiers.

**3. `activate.New` → `validateActivationRule` (internal/activate/activate.go).**
Pure lift of the per-rule suggestion validation: max-suggestion count against `maxEntries`, `inspect_reachable_repository` fallback minimum, why-line rune bounds 1..80, score bounds 0..1 — same error strings, same index-ordered iteration, and still executed per rule **before** that rule's `regexp.Compile` and its wrapped compile error, so validation order across rules is unchanged. Handle-type/verb schema compilation above the loop is untouched.

**4. Test-only consolidation.** `validation_candidate_artifact_test.go` gains `reviewCandidateDocument`, `loadReviewCandidate`, `verifyCandidateSetInWorkingTree`, `sortedCandidatePaths`, and `verifyCandidateSetDigest`; the two existing historical guards keep their exact base commits (94761cf/16 files at 71f9648; 0ffb8f9/19 files at 610fa58). The superseded build-recipe guard in `validation_build_recipe_test.go` becomes `…MatchesHistoricalPayload` verifying the recipe candidate's 20-file set at commit `ed38607` — correct, since that candidate described the tree as of decision 0015 and the working tree has since legitimately diverged. The new working-tree guard pins this candidate's `base_commit`, 20-file count, false gates, and set digest. All five guards pass live.

**5. Documentation.** README.md and validation-execution-boundary.md each change exactly one line: the old world-build digest string replaced by the new one. No other prose changes.

## Findings

No P0 findings. No P1 findings. P2 findings below; none blocks.

- **P2-1 (new helper sits exactly at the complexity threshold).** `classifyScheduledOutputEntries` measures gocyclo 15 — passing (`-over 15` flags only >15) but with zero headroom. Any future conditional added to directory classification re-trips the quality gate on a frozen-boundary file. Non-blocking; noted for the next cycle.
- **P2-2 (candidate-note phrasing implies three expected failures).** `complexity-refactor-candidate.md` says "the full runner suite passes except the two freeze guards and the world-build reproduction test", which reads as three failing tests. Exactly two fail: `TestPreValidationFreezeMatchesAcceptedCandidates` and `TestValidationBuildReproducesFrozenWorldBuildDigest` (the latter *is* one of the two, not a third). The note's own line 5 states the correct pair. Documentation imprecision only; the machine-checkable claims are all accurate. Non-blocking.
- **P2-3 (sharp edge in the `applyScheduledResume` contract).** On stop paths it returns `(true, err)`; the `done` flag is meaningful only when `err == nil`. The sole caller checks `err` first, so behavior is identical today, but a future caller that consulted `done` before `err` could mask an evidence-write failure. Non-blocking.

## Residual notes (non-blocking)

- The two freeze tests remain expected-fail until an authorized refreeze imports the new 20-file identity and the new world-build digest; the freeze test fails at the first divergent path (README.md), and the reproduction test's live digest equals the candidate's claimed new value.
- A real validation run still requires a linux/amd64 execution host; this payload changes build identity only because `internal/activate` is linked into `cmd/phoenix` — the build recipe itself is unchanged (verified: `phoenixBuildRecipeForTranche`, `phoenixBuildCommand`, and `buildPhoenix` are untouched by the diff).
- The reproduction test is skipped under `-short`; the non-short full suite exercised it here, and the manual rebuild reproduced the digest independently of the payload's own test code.
- staticcheck v0.7.0 self-selected toolchain go1.26.7 (its own `go >= 1.25` floor); result clean.
- Same-machine / shared-workspace limitation: review ran in the primary checkout at the payload commit rather than a disposable detached worktree; clean `git status` and empty `git diff 8ed9202` were confirmed before and after all checks, including after deleting `bin/phoenix`, `build/manifest-review.json`, and the emptied `bin/`/`build/` directories.

## What this review did NOT do

It did not run any model, arm, trial, or schedule; did not read or write anything under `experiments/frontier-v1/results/`; did not contact a custodian or grader; did not open, hash, or enumerate any private or held-out label; did not modify the gate document, freeze hashes, protocol, schedules, manifests, or any tracked file; did not import, refreeze, or open either gate; and did not verify decision 0016/0017 chair authority beyond confirming the decision records exist and the gate document holds both gates false.

## Public checks

Byte-identity script (scratchpad Go program, LF-normalization rule): 20/20 file digests match; set digest reproduces `sha256:2d20be1ac41cc9072ae472b9b5a73e2607e9d35093fda49c422fa487117ecb4b`; candidate file raw SHA-256 `67ce91a3a1fb991c2fc2565f2d6ca35f0bd7a352daa40dce3ab629f9bea8ac34` (5,155 bytes). `git diff 8881d5f..8ed9202` reviewed hunk-by-hunk for all nine files. Manual frozen-recipe rebuild + `cmd/build-manifest` printed `sha256:b5a26d5e2290c7919e4bc629a774f387a766107539b4fcfdf7d55d0f1c19a2c4`. Full suite `go test -count=1 ./...`: only the two documented failures. Focused runs: 22/22 scheduled tests, 5/5 activate tests, 5/5 candidate guards. `gocyclo -over 15` exit 0; dupl gate exit 0; `gofmt -l` empty; `go vet` clean; staticcheck v0.7.0 clean.

**Final verdict: ACCEPT** — payload `8ed9202c80d8f591c5d0db8a7e8a349022f952f8` is a verified byte-exact, behavior-preserving complexity decomposition that brings the quality gate clean, changes no production semantics, evidence byte, stop rule, checkpoint rule, or refusal condition, tampers with no gate or freeze state, and reproduces the claimed new world-build digest under the unchanged frozen recipe; the only failing tests are the two documented pre-refreeze guards, and all three P2 findings are non-blocking.
