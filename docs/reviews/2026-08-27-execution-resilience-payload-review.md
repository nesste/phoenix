# Independent review — protocol-v5 section 9 execution-resilience payload

**Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author any file under review. I reviewed by reading the repository at the pinned commit and running the test and quality commands.
**Date:** 2026-08-27
**Payload commit:** `7a1ea42513acd3d55e276eabc2459da4a037acc0`
**Verdict: REVISE** — two P1 findings, both in newly added code, both cheap to fix. No P0. No model was run, no private grade obtained, no validation outcome observed; `may_open_validation` and `may_open_held_out` were false throughout and I did not touch them.

## 1. Pinned inputs — digests I computed

Digest rule applied: CRLF→LF, then CR→LF, then SHA-256.

| Input | Expected | Computed | Match |
| --- | --- | --- | --- |
| `git rev-parse HEAD` | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | yes |
| `git rev-parse HEAD^` | `2cf99331688f4705ed7d6cd3efed98941f59defc` | `2cf99331688f4705ed7d6cd3efed98941f59defc` | yes |
| `artifacts/gate-1a-execution-resilience-candidate.json` raw | `sha256:2e6c3c90…f99f93` | `sha256:2e6c3c90d73f15c1713bd026237fefec48adbdf3a6abec369350cc3b66f99f93` | yes |
| `artifacts/execution-resilience-candidate.md` raw | `sha256:37494387…43eb80` | `sha256:3749438786e13fa822fd9f0b99052b4ce8aca53bd6fb4eabf016248c4243eb80` | yes |
| `protocol.json` LF | `sha256:726b6811…cbdb0d32` | `sha256:726b68112108542ba3491fde538342399cbb446a04755a7e0208e076cbdb0d32` | yes |
| `validation-execution-boundary.md` LF | `sha256:247e5f8c…2f081e` | `sha256:247e5f8cc4ee355ea23e8414c3ca4389070231e5cc0570deb1853d008f2b081e` | yes |
| `pre-validation-artifacts.json` LF | `sha256:47f4bd4a…a1bed3` | `sha256:47f4bd4af3635510eb932dc8528372747eb665acd47b2beda653192a1a61bed3` | yes |
| `runner/execution_resilience.go` LF | `sha256:4f6cf801…b5a3696` | `sha256:4f6cf8015f66d6814711d27409bdbc9ac070405ee9cab10c042660a66b5a3696` | yes |
| `runner/process_events.go` LF | `sha256:0dabcd78…c080aa9e17` | `sha256:0dabcd7855ab2a457bb5cea8394fdbba1558add0136b560fddb723c080aa9e17` | yes |
| Inventory file count | 16 | 16 | yes |
| Ordinal-path set digest | `sha256:1f0b6f87…fcde7c609` | `sha256:1f0b6f872f58d17ddb987c0106c0c0df3dbd6385d060ac7991c97fafcde7c609` | yes |

I independently recomputed the LF digest of **all 16** inventory files and compared each against the inventory's `files` map: **zero mismatches**. The set digest above was rebuilt from those 16 recomputed values, not copied from the document.

Other identities I verified rather than trusted:

- World-build identity in `pre-validation-artifacts.json`: `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` — unchanged, and `TestValidationBuildReproducesFrozenWorldBuildDigest` passes.
- Grader digest `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b` — unchanged; the string does not appear in this commit's diff at all.
- Accepted local candidate `sha256:a70ce35b7cfbb686cf7922c99394f6c3ef344235b03f8cf3f88bccb8ce892967` — recomputed, matches.

**Working tree:** exactly one uncommitted path, `docs/reviews/2026-08-27-execution-resilience-payload-review-prompt.md` (untracked). Nothing else is dirty. No P0 on this ground.

**Files changed vs parent:** 19 (16 inventoried + the inventory, the candidate note, and the uninventoried guard test `runner/execution_resilience_candidate_artifact_test.go`). The scope matches the declaration.

## 2. Test and quality command results

| Command | Result |
| --- | --- |
| `go test ./...` (root module) | 1 FAIL, all else ok — **only** `TestPreValidationFreezeMatchesAcceptedCandidates` (`artifact_freeze_test.go:149`, README digest `f8d657c5…` vs frozen `b486fc55…`) |
| corpusctl suite (`cd experiments/frontier-v1/corpusctl && go test ./...`) | ok |
| `experiments/surface-spike` suite (extra, not required) | ok |
| `go vet ./...` | clean, exit 0 |
| `staticcheck@v0.7.0 ./...` | clean, no diagnostics |
| `gocyclo@v0.6.0 -over 15` over the gated dirs | clean, no output |
| `dupl@v1.1.0 -plumbing -t 100` under `quality-check no-output` | clean, no output |
| `make validate-spec` | all 9 examples valid |
| `make validate-authoring` | valid |
| `gofmt -l` over changed dirs | clean |

**No failure beyond the single expected freeze red.** I also ran the named greens explicitly: `TestLocalArtifactCandidateMatchesImplementation` PASS, `TestValidationBuildReproducesFrozenWorldBuildDigest` PASS, `TestExecutionResilienceCandidateMatchesWorkingTree` PASS, and all eleven new section 9 tests PASS.

**One correction to the candidate note's green list.** `verifyAcceptedLocalArtifacts` does **not** run at this commit: it is called at `artifact_freeze_test.go:152`, after the `t.Fatalf` inside `verifyReplacementRunnerFreeze` at line 149. The note's claim that it "keeps passing throughout" is therefore not demonstrated by the suite. I verified it independently by reimplementing its logic outside the repository against `pre-validation-artifacts.json` and `artifacts/pre-validation-local-candidate.json`: candidate digest matches, review record exists, and all four local-artifact blocks (`runtime_invocation_and_exact_per_arm_system_prompts`, `arm_a_schemas`, `world_definition_and_world_build_digest`, `grader_digest`) match their reviewed candidate bodies with the expected metadata. **It would pass.** Recorded as P3-8, not as a defect in the payload.

**Carried P2-1 is genuinely retired.** Measured: `classifyScheduledOutputEntries` gocyclo **9** (was 15 at threshold with zero headroom); the new helpers `scheduledOutputEntryKind` 7, `finishScheduledInventory` 8, `recordScheduledAssignment` 3; `verifyScheduledCheckpoint` **4**, with `verifyCheckpointIdentity` 8 and `verifyCheckpointProgress` 8. The decomposition is real and the headroom claim is accurate. See P2-5 for what replaced it.

## 3. Findings

### P1-1 — The signal handler swallows SIGINT and SIGTERM, and the log then lies

`experiments/frontier-v1/runner/process_events.go:150-165`, wired at `experiments/frontier-v1/runner/main.go:139-143`.

`signal.Notify(signals, os.Interrupt, syscall.SIGTERM)` disables Go's default action for those signals for the whole program. The watcher goroutine records one termination event and returns; nothing re-raises, exits, or cancels the run loop. Three consequences, all introduced by this payload:

1. **A scheduled run can no longer be stopped cleanly.** Ctrl+C on the recorded Windows host and SIGTERM on the WSL2/Linux host both become inert for the entire duration of the run. Before this payload both killed the process. Section 9's host requirements call for a *supervised* custodian process; a multi-hour, 300 USD run that ignores its supervisor's stop signal is the opposite.
2. **The custody log records a falsehood.** A `termination` event with detail `received interrupt` is appended, and then the runner keeps launching trials. An auditor reading the log sees a termination that did not happen, at a launch index the run went on past.
3. **The artifact fails in exactly the case §9 built it for.** The only signal that still stops the process is SIGKILL, which cannot be caught, so it produces *no* log entry. The process-event log therefore records catchable signals that do not kill and misses the uncatchable one that does — inverting its evidentiary purpose for the `process_kill` cause.

Only one signal is ever handled, too: the `select` receives once and the goroutine exits, so a second signal is neither recorded nor acted on.

**Fix:** after recording, restore the default disposition and re-raise (`signal.Stop(signals)` then re-send the signal to self), or cancel the run loop so `runScheduledCases` returns and the stop event is written. Either way the process must actually stop, and a `termination` entry must only be written when it does.

### P1-2 — The single-resume bound is laundered by editing one integer, and the corroborating evidence is on disk and unread

`experiments/frontier-v1/runner/execution_resilience.go:153-155` (the bound), `experiments/frontier-v1/runner/scheduled_resume.go:302-305` (the only validation of the counter), `experiments/frontier-v1/runner/scheduled_run.go:90` (the increment).

The ordering the prompt asks about is **correct**: the counter is checked at line 153, before `--resume-attestation` is even read at line 157. And a hand-edited checkpoint claiming a *higher* count is refused — `verifyCheckpointProgress` rejects `Resumes > 1`. But the attack is downward, and it is unbounded:

> After the permitted resume, `scheduled-checkpoint.json` holds `"resumes": 1`. A second interruption arrives. The custodian edits that one field to `0`. `verifyCheckpointProgress` accepts `0` (it is in `[0, 1]`, and no other field constrains it — `next_launch_index` and `spent_usd` reconcile against the assignment records, `resumes` reconciles against nothing). `requireResumeAuthorization` then sees `Resumes == 0 < 1`, computes `wantIndex = 1`, and accepts a fresh attestation declaring `resume_index: 1`. The tranche takes a second resume with every check green.

The candidate note answers the wrong direction: "a hand-edited checkpoint claiming a *higher* count is a mismatch, not a licence." Nothing addresses lowering it.

What makes this a P1 rather than an accepted limit of a trust-the-custodian design is that **this payload already writes the independent evidence and does not consult it.** `process-events.jsonl` is append-only, fsynced per append, and receives exactly one `start` entry per custodian process (`scheduled_run.go:99`). The count of `start` entries *is* the resume count. Further, each resume's attestation commits `process_event_log_sha256` into git, so a rewritten log is detectable across commits in a way an edited checkpoint integer is not. The runner already opens and hashes this file at resume (`verifyAttestedProcessEventLog`, `execution_resilience.go:209-221`) — it hashes it and never parses it.

**Fix:** in `requireResumeAuthorization`, parse the process-event log, count `processEventStart` entries, and refuse unless `checkpoint.Resumes >= startCount - 1`. Roughly five lines, using a file the payload already reads.

### P2-1 — The outcome-free field-set test is depth-1; a new field slips past inside `artifact_digests`

`experiments/frontier-v1/runner/execution_resilience_test.go:268-282`.

The both-directions guard is real *at the top level*: `assertPartialSummaryFieldsAreOutcomeFree` decodes into `map[string]json.RawMessage` and checks membership in both directions against `outcomeFreePartialFields`. But nested objects are `json.RawMessage` and are never inspected. Adding a field to `partialArtifactDigests` (`execution_resilience.go:35-40`) — say `PassCount int \`json:"pass_count"\`` or a per-arm grade digest — changes nothing at the top level, so the suite stays green and the field ships. `arm_assignment_counts` is likewise a free-form `map[string]int` whose keys are unchecked.

This is precisely the "can a new field slip past it" failure the prompt asks about, and it lands in the one substructure a future author is most likely to extend.

**Fix:** recurse into `artifact_digests` with its own admitted key set, and assert that `arm_assignment_counts` keys are a subset of `phase1Arms`.

### P2-2 — Nothing tests that the resume counter is ever incremented or persisted

`experiments/frontier-v1/runner/scheduled_run.go:90`.

Deleting `progress.resumes++` breaks no test. I checked every `Resumes` reference in the suite: `execution_resilience_test.go:31` and `:60` assert `== 0` on a *fresh* run (true either way); `:81` passes the value directly to `recordTermination`; `:109` and `:363` construct a `scheduledCheckpoint` literal without going through `writeScheduledCheckpoint`. The added helper `writeTestCheckpointWithResumes` (`scheduled_resume_test.go:349`) is called exactly once, from `writeTestCheckpoint`, with a hardcoded `0`.

So the round trip that the candidate note advertises as the mechanism — resume increments, checkpoint persists, next resume reads it back and refuses — is never exercised end to end, and `verifyCheckpointProgress` refusing `Resumes: 2` is never exercised at all. Combined with P1-2, the single-resume bound is the weakest of the five conditions both in code and in coverage.

**Fix:** one test that runs a schedule, truncates to an interruption, resumes with a valid attestation, and asserts the on-disk checkpoint reads `resumes: 1` and that a further resume is refused.

### P2-3 — The 72-hour window is entirely self-reported; two independent anchors exist on disk and neither is used

`experiments/frontier-v1/runner/execution_resilience.go:223-251`.

`verifyResumeAttestationTiming` parses four timestamps out of the attestation and compares them **only to each other**. Ordering is enforced correctly (interruption → classification → verification → authorization, `:237-242`), and `elapsed > resumeWindow` refuses (`:243-246`). But every value is a number the custodian typed. The function never reads `time.Now()`, and never reads the interruption time that the payload itself wrote to disk.

> A custodian resuming three weeks after the interruption writes `interrupted_at: <now-1h>`, `resume_authorized_at: <now>`, and passes. The 72-hour condition is satisfied by arithmetic on invented inputs.

Section 9 words condition (b) as custodian-verified and recorded — the payload is right to treat digest identity that way. But condition (c), "resume occurs within 72 hours of the interruption," is a fact about the world, and the runner has two independent handles on it that this payload created: the `recorded_at` of the last entry in `process-events.jsonl`, and `updated_at` in `scheduled-summary.partial.json`. It consults neither, nor its own clock.

**Fix:** require `interrupted_at` to be at or after the last process-event `recorded_at` (a resume cannot precede the last thing the interrupted process did), and require `resume_authorized_at` to be within a small tolerance of `time.Now()`. That converts (c) from a claim into a check without touching the conditions §9 deliberately leaves to the custodian.

### P2-4 — The post-closure gate's "record" check is existence only

`experiments/frontier-v1/runner/validation_gate.go:100-102`.

The verdict test is exact and correct (`ClosureVerdict != "ACCEPT"` refuses, verified by `TestPostClosureGateRefusesValidationAfterANonCompletedExecution` with `REVISE`). The record test is `os.Stat` on a path from the same document. So:

> `validation_execution_closure_review: "README.md"` with `validation_execution_closure_review_verdict: "ACCEPT"` satisfies the gate. So does any directory — `os.Stat` does not check for a regular file. The named record's own verdict line is never read; the verdict is taken from the gate document that the same edit sets.

Both fields are frozen bytes under the freeze workflow, which is the real control, so this is not a P1. But the prompt asks directly whether the gate "can be satisfied by a record that does not exist, or by a non-ACCEPT verdict": non-ACCEPT is genuinely refused; a record that does not exist *as a review* is not.

**Fix:** require the path to be a regular file under `docs/reviews/`, non-empty, and containing the verdict string it claims.

### P2-5 — P2-1 is retired and immediately re-created on a new function

`experiments/frontier-v1/runner/execution_resilience.go:169`.

`verifyResumeAttestationFields` measures **gocyclo 14**, one below the `-over 15` gate — the same near-threshold condition that carried finding P2-1 recorded for `classifyScheduledOutputEntries`. It is the joint-highest non-test function in the repository alongside `validateExternalGrade` (14) and `runScheduledCases` (14). The gate passes; the structural point that made P2-1 worth recording now applies to the payload's central check function.

**Fix or accept explicitly.** Splitting the identity checks (`:170-187`) from the digest and log checks (`:188-191`) would restore headroom cheaply.

### P2-6 — The termination event's progress is a hardcoded zero in production

`experiments/frontier-v1/runner/main.go:141`.

The progress callback passed to `watchProcessTermination` is `func() processEventProgress { return processEventProgress{} }`. Every termination entry a real run writes therefore records `next_launch_index: 0`, `completed_pairing_keys: 0`, `spent_usd: 0`, `resumes: 0`, regardless of actual progress. The unit test at `execution_resilience_test.go:74-88` passes progress explicitly and so does not cover the production wiring.

This makes the candidate note's claim inaccurate: "Each entry carries `v`, the event kind, an RFC3339Nano UTC timestamp, tranche, schedule digest, world build, the launch index / completed pairing keys / cumulative spend at that moment, the resume count, and the originating principal." That is true of `start`, `checkpoint`, and `stop`, and false of `termination` — the one entry a resume attestation is meant to cite. Worse, an entry reading `resumes: 0, next_launch_index: 0` on the second interruption of a resumed tranche actively misleads an auditor and would defeat the P1-2 cross-check if that were built on the termination entry rather than on `start` counts.

**Fix:** thread the live progress into the closure (the same value `scheduledProgress.events` already computes), or move the watcher inside `runScheduledCases` where that state exists.

### P3 notes

- **P3-1** `experiments/frontier-v1/runner/execution_resilience.go:324-338`. `recordCheckpoint` writes the partial summary *after* appending the checkpoint event, so the summary's `process_event_log_sha256` is correct at each boundary — but `recordStop` (`:287-297`) appends the terminal entry and nothing rewrites the summary. The runner-committed log digest therefore never covers the terminal entry. Harmless given that the custody record is hashed by hand over the whole file at closure, but the boundary document's "the log's digest is committed with the custody record" is discharged by a human step, not by code.
- **P3-2** `experiments/frontier-v1/runner/execution_resilience.go:47-60`. The partial summary is **broader** than §9's literal list (which admits launch index, completed pairing keys, per-arm counts, artifact digests, timestamps, cumulative spend). It adds `v`, `tranche`, `status`, `schedule_digest`, `resumes`. I judge this **defensible**: none carries an outcome, and identity fields are needed to bind the record to a tranche. Naming it because the payload must be judged against the accepted text.
- **P3-3** `experiments/frontier-v1/runner/execution_resilience.go:291-294` with `scheduled_run.go:322`. `recordStop`'s `detail` is `summary.Status + ": " + summary.StopReason`, and a safety-stop `StopReason` concatenates an arbitrary `cause.Error()` — which for grader-shape failures can embed a per-trial status string (e.g. `custody-safe grade status "fail" is inconsistent with its checks`). `process_events.go:57` describes the log as carrying "outcome-free lifecycle events". Record a fixed stop-kind token instead of the full error text.
- **P3-4** `experiments/frontier-v1/pre-validation-artifacts.json:17-19`. The three added `gates` lines are indented with 2 spaces where the surrounding block uses 4. The digest is computed over LF-normalized raw bytes, so this inconsistency becomes part of the frozen identity. Fix before the refreeze re-pins it.
- **P3-5** `experiments/frontier-v1/runner/main.go` `writeJSON`. The checkpoint and partial summary are written with `os.WriteFile` — no fsync, no temp-and-rename — unlike the process-event log, which is genuinely `O_APPEND` + `Sync` per entry (`process_events.go:102-114`). A kill during a boundary write can truncate the checkpoint that §9 calls "durable". The failure mode is **closed**: `decodeStrict` errors and the resume is refused, closing the tranche indeterminate. Pre-existing for the checkpoint, newly load-bearing now that it carries the resume counter.
- **P3-6** `experiments/frontier-v1/protocol.json:350`. Blocker 5 loses two clauses, not one: the §9 clause this payload discharges, and the §5 evaluator-countersignature clause. I verified the latter is genuinely discharged (`docs/decisions/0025-*.md` records the countersignature and the ACCEPT verdict), so the edit is correct — but it is outside this payload's stated scope and should be named in the decision document rather than absorbed silently.
- **P3-7** `experiments/frontier-v1/runner/main.go:219` and `execution_resilience.go:194-207`. The `--resume-attestation` path gets `resolveRepositoryPath` but no `ensureInside`, and `verifyResumeFrozenDigests` skips comparison when the expected digest is empty. Neither is reachable on the validation path (`requireFrozenWorldBuild` and the frozen schedule digest guarantee both are non-empty), but the skip is a silent-pass shape in a check function.
- **P3-8** `experiments/frontier-v1/runner/artifact_freeze_test.go:149-152`. `verifyAcceptedLocalArtifacts` is unreachable at a payload commit because the deliberate runner-freeze `t.Fatalf` precedes it. The candidate note asserts it "keeps passing throughout" without that being observable. I verified it independently (see §2). Consider hoisting it above the runner-freeze check so payload commits keep exercising it.
- **P3-9** `experiments/frontier-v1/runner/scheduled_run.go:96-98`. When `runRemainingSchedule` returns an error, `recordStop` is skipped, so the log ends on a `checkpoint` with no terminal entry — indistinguishable from an abrupt kill. An `error` stop kind would keep the log unambiguous.

### Attempted attacks that failed — the payload holds

Recording these so the negative results are on the record rather than assumed:

- **Launder a validation run through the authoring exemption.** Not possible. `validateScheduleForTranche` (`schedule.go`) refuses unless `schedule.Tranche == tranche`, so the frozen validation schedule cannot be run under `--tranche authoring`. The exemption cannot be reached by a validation execution.
- **Delete the checkpoint to look like a fresh run.** Refused. `finishScheduledInventory` (`scheduled_resume.go:180-182`) requires a checkpoint whenever assignments, evidence, *or* the new progress records are present — and the partial summary and process-event log are both classified as progress (`:140`), so they now hold the requirement open even after the assignments are removed.
- **Cite a stale or forged process-event log.** Refused. `verifyAttestedProcessEventLog` (`:209-221`) recomputes the digest from the file on disk and compares; covered by the `stale process-event log digest` case.
- **Resume with an unlisted or empty cause.** Refused, with an accurate message. `TestValidationResumeAcceptsOnlyFrozenOutcomeUncorrelatedCauses` covers `oom`, `out_of_memory`, `operator_choice`, and `""`, all hitting "is not on the frozen outcome-uncorrelated cause list".
- **Sneak an extra field into the attestation.** Refused. `decodeStrict` sets `DisallowUnknownFields` and rejects trailing JSON.
- **Point the attestation at another tranche's directory.** Refused (`output_dir` compared against the repo-relative output path).
- **Reach a launch before the resume refusal.** Not possible. `openScheduledProgress` runs before `applyScheduledResume` and before `runRemainingSchedule` (`scheduled_run.go:40-52`); no model is invoked first. The validation gate, including the post-closure audit, is the very first thing `prepareScheduledCLI` does (`main.go:188-194`), ahead of the world build, grader preparation, and output-directory work.

## 4. Answers to the eight check areas

**1. Section 9 conditions, one by one.**

- **(a) No inspection.** `execution_resilience.go:182-184`. A pure boolean claim, refused when false. Unverifiable by construction — §9 words it as "attested in writing by the custodian", so this is correct scope, not a gap.
- **(b) Digest identity.** `:182-184` plus `verifyResumeFrozenDigests` (`:194-207`). Stronger than a bare claim: the attestation's recorded `schedule_digest` and `world_build` must equal the live values, so a stale attestation from a different world build fails. Byte-level re-verification of the whole frozen set stays with the custodian, exactly as §9 words it.
- **(c) 72-hour window.** `:243-246`. Enforced, but on self-reported timestamps only — **P2-3**. A resume can succeed with (c) unmet by back-dating.
- **(d) Pre-decision cause classification on the frozen list.** `:179-181` and `:182-184`. Enforced. **OOM is refused** — it is absent from `frozenResumeCauses` (`:28-33`) and falls to the generic refusal; so does the empty string and any unlisted value. **The refusal message is accurate**: "interruption cause %q is not on the frozen outcome-uncorrelated cause list" states exactly §9's reason for striking OOM.
- **`process_kill` without establishment.** Refused at `:185-187` ("a process kill resumes only when the process-event log establishes a non-participant principal"). The log digest is re-verified against the file on disk at `:209-221`, so a stale or forged citation fails. One presentational divergence: §9 lists three safe causes and treats `process_kill` as separately admissible; the code puts `process_kill` inside `frozenResumeCauses` and guards it with a second check. Behaviourally identical; `protocol.json` states it §9's way.
- **Timestamp ordering.** Enforced, all four in sequence, `:237-242`.
- **Promptness and lapse.** `:247-249` requires a `lapse_explanation` whenever `resume_authorized_at` is after `conditions_verified_at`, and `promptness_statement` is mandatory at `:173-175`. This is *broader* than §9, which attaches the written explanation to "a closure by clock lapse" — a case in which no resume happens and the runner is never invoked at all. §9's actual lapse duty is documentary and is recorded in the boundary document; it cannot be code-enforced. I judge the payload's stricter reading defensible and conservative, but it is a reinterpretation and should be named as such.
- **One resume per tranche.** Checked at `:153-155`, **before** the attestation is read — the ordering the prompt asks about is right. The counter is carried in the durable checkpoint. But it is launderable downward — **P1-2** — and untested — **P2-2**.
- **Abandonment by choice.** No code path represents it. A refusal returns an error and the run stops without launching; the tranche closes indeterminate under the existing stop rules. `frozenResumeCauses` contains no discretionary cause, and `operator_choice` is covered as a refused case in the tests. Correct.

**2. The partial summary is outcome-free.** Field set: `v`, `tranche`, `status` (always `in_progress`), `schedule_digest`, `artifact_digests{world_build, grader_digest, arm_b_document_digest, process_event_log_sha256}`, `started_at`, `updated_at`, `next_launch_index`, `completed_pairing_keys`, `arm_assignment_counts`, `spent_usd`, `resumes`. **No success count, grade tally, pass/fail status, or per-arm outcome field.** `arm_assignment_counts` counts *scheduled* assignments below `nextLaunchIndex` (`:88-96`) and is fully determined by the launch index given the frozen 4-arm pairing-complete schedule — it carries no outcome. The both-directions test is **real at the top level** and **not real one level down** — **P2-1**. **Spend is the only outcome-adjacent quantity**, and its retention is justified exactly as §9 says: the budget-stop rule requires it, and the mandatory-resume rule removes the channel by which it could steer a resume choice. Field set is broader than §9's literal list by five non-outcome identity fields — **P3-2**, judged defensible.

**3. The process-event log.** Records all four kinds — `start` once per custodian process after resume authorization and before any launch, `checkpoint` at every pairing-key boundary, `stop` at the terminal state, `termination` on signal — each with an RFC3339Nano UTC timestamp. Verified by `TestScheduledRunAppendsProcessEventsForEveryLifecyclePoint`, which asserts the exact kind sequence. **Genuinely append-only and durable**: `O_APPEND|O_CREATE|O_WRONLY` with `file.Sync()` before close on every entry (`process_events.go:102-114`), so an abrupt kill cannot lose already-recorded entries.

**The "where the host exposes it" claim is honest and the payload is right on the technical point.** A Go signal handler receives only the signal number; neither Windows nor Linux delivers the sending principal through `os/signal`. (On Linux the kernel *can* carry `si_pid` in a `SA_SIGINFO` handler, but Go's runtime does not surface it through `signal.Notify`, so the claim holds for this program on both supported hosts.) `recordTermination` (`:126-132`) sets `exposed: false`, `source: signal_without_sender_identity`, and a reason, rather than misattributing the kill to the receiving process's owner — which is the honest and correct choice, and is exactly why §9 makes `process_kill` unsafe absent a written establishment. Verified by `TestProcessTerminationRecordsAnUnexposedSendingPrincipal`.

Two defects sit on top of this correct design: the termination entry's progress is a constant zero in production (**P2-6**), and the handler prevents the termination it records (**P1-1**). The digest is committed in the partial summary as required now; commitment with the custody record at closure is a human step, and the runner-side digest never covers the terminal entry (**P3-1**).

**4. The post-closure gate.** **Encoded in the frozen gate rules, not merely prose.** Three respects:

- Code: `requirePostClosureAudit` (`validation_gate.go:91-104`), called from `requireValidationGate:71`, which is the first check `prepareScheduledCLI` performs for a validation tranche.
- Frozen document: `validation_execution_closure_gate`, `validation_execution_closure_review`, `validation_execution_closure_review_verdict` in `pre-validation-artifacts.json`.
- Frozen test: `verifyPostClosureGateState` (`artifact_freeze_test.go:210-226`) asserts the rule text and pins both closure fields **empty**, while `verifyValidationGateState:188` pins `validation_execution_status` to `"indeterminate"`.

**It fires for every non-completed ending, not only interruption.** The predicate is `status != "" && status != "complete"` — a default-deny, so `budget_stopped` and `lapsed` are covered without enumeration, and `TestPostClosureGateRefusesValidationAfterANonCompletedExecution` exercises all three status strings explicitly.

**It correctly applies to the decision-0023 closure, and that is registered as a blocker.** `validation_execution_status` is already `"indeterminate"` and both closure fields are empty, so the gate refuses on that ground today; the freeze test pins that state, so it cannot drift. `protocol.json` `remaining_execution_blockers` gains "an independently reviewed and committed closure record for the interrupted 1,475-launch validation execution closed by decision 0023". This is real, not prose-only.

**Can it be satisfied without a real reviewed record?** A non-ACCEPT verdict is genuinely refused. A record that does not exist is refused. But *any existing path* — including a directory or an unrelated file — plus an ACCEPT in the gate document satisfies it, because the record's own content is never read: **P2-4**. One structural caveat worth naming: nothing in the runner *writes* `validation_execution_status`, so the gate's teeth after a future interruption depend on the chair updating that field at closure. That is consistent with the pre-existing design (the field predates this payload) and is the freeze workflow's job, but it means the gate is armed by discipline, not by the runner.

**5. Scope decision — the authoring-tranche exemption. My judgment: defensible, and it is a decision rather than an oversight, but it is genuinely narrower than the accepted text and must be recorded as an accepted residual rather than passed silently.**

Stating the divergence plainly first: §9's opening sentence is "If **an execution** is interrupted by infrastructure failure" — unqualified. The payload reads it as validation-only (`execution_resilience.go:150`).

Reasons it is defensible:

1. §9's machinery is inoperable for authoring. Every consequence it names — "the tranche closes indeterminate", "burns the tranche exactly as under v4", the custody record, the 300 USD ceiling, the frozen schedule digest an attestation is checked against — is a validation construct. There is no indeterminate closure for an authoring loop to fall into, so a refusal there would have no defined meaning.
2. §9's own post-closure paragraph says "after any **validation** execution", showing the section's subject is the validation tranche even where individual sentences drop the qualifier.
3. Authoring outputs are not trusted on the custodian's word anyway: every authoring artifact reaches a frozen state only through payload → independent review → refreeze, which is a strictly stronger control than an attestation.
4. Most decisively for the security question: **the exemption cannot be used to launder a validation execution.** `validateScheduleForTranche` refuses unless `schedule.Tranche == tranche`, so the frozen validation schedule cannot be run under `--tranche authoring`; `--tranche validation` is the only route to it, and that route is not exempt.

The one substantive cost: an authoring resume after inspecting authoring outcomes could in principle condition authoring artifacts on side signals, and nothing mechanical prevents it. Given (3), I do not think that warrants the ceremony, but the chair should accept it knowingly. It is properly declared — in the inventory's `scope_decision`, in the note, in the boundary document, in the README, and pinned by `TestAuthoringResumeIsExemptFromTheAttestationRule`.

**6. Records match the implementation.** I read `protocol.json` `amendment.execution_resilience`, the boundary document's new section, and the two README paragraphs against the code. The amendment text tracks §9 closely and I found **no overclaim** in it. The boundary document's five numbered conditions match the five enforced checks, including the ordering requirement and the digest-recording nuance in condition 2, and it correctly carries both senses of the lapse duty. The README's two paragraphs are accurate.

**The one overclaim is in the candidate note, not in the frozen records**: "Each entry carries … the launch index / completed pairing keys / cumulative spend at that moment, the resume count" is false for `termination` entries in production (**P2-6**). A second, smaller one: "`verifyAcceptedLocalArtifacts` keeps passing throughout" is not observable at this commit (**P3-8**); I verified it independently and it holds.

Verified: blocker list **7 → 8** (counted at both commits). The struck §9 clause is genuinely discharged by this payload; the payload additionally strikes the §5 countersignature clause, which decision 0025 does discharge (**P3-6**). The added closure blocker is accurate. `corpusctl` `protocol_test.go` **asserts** the new contract rather than tolerating it: eleven required substrings against `amendment.execution_resilience` (including "OOM is struck", "At most one resume per tranche", "no per-arm outcome", "Post-closure gate"), an exact length check `len(protocol.Remaining) != 8` with `Remaining[7]` pinned to "Task 0.6", and the blocker-text requirements updated to "closure record" and "post-closure gate". A regression in any of these fails the corpusctl suite.

**7. Freeze hygiene.** Exactly one deliberate red, `TestPreValidationFreezeMatchesAcceptedCandidates`, and it is the one named. `TestLocalArtifactCandidateMatchesImplementation` and `TestValidationBuildReproducesFrozenWorldBuildDigest` are green (run explicitly). `verifyAcceptedLocalArtifacts` is unreachable at this commit but verified green by independent replication (**P3-8**).

The Phoenix binary is untouched: the changed-path list contains nothing under `cmd/`, `internal/`, `verbs/`, `build/`, `worlds/`, or `go.mod`/`go.sum`, so the world-build identity `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` holds, confirmed by the build-reproduction test. The grader digest is unchanged at `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b`: no corpusctl non-test source, no `go.mod`/`go.sum`, no `schema/*.schema.json` changed — the only corpusctl file in the diff is `internal/corpus/protocol_test.go`.

**Is editing the `gates` block in a payload commit acceptable?** **Yes, on these facts.** The edit is confined to three added keys; no digest, no `frozen` flag, and no accepted-candidate pointer moves; `may_open_validation` and `may_open_held_out` stay false; and both new closure fields are added **empty**, so the edit strictly *tightens* the gate rather than opening anything. It also could not sensibly wait for the refreeze: the freeze test's `verifyPostClosureGateState` and the runner's `requirePostClosureAudit` both need the fields present, so deferring them would have made the payload's own tests unrunnable. The stale digests are exactly what the deliberate red reports. One caveat: the added lines' indentation is inconsistent with the block and will be frozen as-is (**P3-4**).

**8. Quality gates.** All run independently; results in §2. Nothing beyond the single expected red. **Carried P2-1 is verified retired by measurement** (`classifyScheduledOutputEntries` 15 → 9; `verifyScheduledCheckpoint` → 4 with 8/8 halves), and the claim of headroom is accurate. The payload simultaneously introduces `verifyResumeAttestationFields` at gocyclo 14, one below the gate — **P2-5**.

## 5. Residuals I would accept, conditional on the P1s being fixed

The verdict is REVISE, so no residual list is operative yet. For the refreeze that follows the fix, these are the entries I would accept into `accepted_findings`, worded for direct transcription:

- "Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian's written attestation, as section 9 requires."
- "Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b)."
- "Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports 'gate is closed' first; both are refusals and only the surfaced message differs."
- "Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag."
- "Residual: nothing in the runner writes validation_execution_status; the post-closure gate's application after a future non-completed execution depends on the chair setting that field at closure, under the freeze workflow."
- "Residual (P3): the partial summary's committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure."
- "Residual (P3): the partial summary carries five non-outcome identity fields beyond section 9's literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome."
- "Residual (P3): the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log; a kill during a boundary write makes the tranche unresumable, which fails closed to indeterminate."
- "Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first."
- "Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required."
- "Carried P2-1 is retired by this payload: classifyScheduledOutputEntries is decomposed into scheduledOutputEntryKind, recordScheduledAssignment, and finishScheduledInventory (measured gocyclo 9), and verifyScheduledCheckpoint is split into identity and progress halves (measured 4/8/8), so both are back under the gocyclo threshold with headroom."

P2-1, P2-2, P2-3, P2-4, and P2-6 above are **not** in that list: each has a concrete cheap fix and each weakens a section 9 condition, so I would rather see them fixed alongside the P1s than carried. If the chair elects to carry P2-5 (`verifyResumeAttestationFields` at gocyclo 14) instead of splitting it, it should be recorded verbatim as the successor to carried P2-1 rather than dropped.

## 6. What I did not verify

- I did not run a model, obtain a private grade, or observe any validation outcome. Both gates were false at the start and end of this review.
- I did not execute the signal path (P1-1) against a live scheduled run; the finding rests on reading `signal.Notify` semantics against `process_events.go:150-165` and `main.go:139-143`, and on the absence of any re-raise, `os.Exit`, or run-loop cancellation anywhere in the runner (grep-verified: `watchProcessTermination` has exactly one call site).
- I did not empirically compile a mutated `partialArtifactDigests` to demonstrate P2-1; the finding rests on `assertPartialSummaryFieldsAreOutcomeFree` decoding into `map[string]json.RawMessage` and never descending into nested objects, which is unambiguous from the source.
- I made no change of any kind to the repository. My only writes were to the session scratchpad and to this record.

---

## Chair transcription note

Transcribed verbatim from the independent reviewer's record. Before transcription the payload author independently re-verified every digest in §1 against the working tree and `git`: all eleven match. The two P1 findings and P2-6 and P3-4 were independently reproduced by reading the cited lines. The verdict is **REVISE**, so this payload does **not** proceed to refreeze; a revision payload follows and is submitted for a second independent review.
