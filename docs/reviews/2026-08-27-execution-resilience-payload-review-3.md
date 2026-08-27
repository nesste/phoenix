# Third independent review — protocol-v5 section 9 execution-resilience payload, revision 3

**Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author any file under review, and I am not either of the two reviewers who produced the first and second review records. I reviewed by reading the repository at the pinned commit and by running the test and quality commands myself.
**Date:** 2026-08-27
**Payload commit:** `60c6aee55b71716743b30d084a57daeae48b3373`
**Verdict: REVISE** — **one P1 finding, four P3.** No P0, no P2. Every P1 and P2 from both prior reviews is genuinely fixed, and I confirmed each by independent attack rather than by reading the implementer's account. The single P1 is in the code written this round to fix review 2's P3-N3, and it fails in **both** directions against this repository's own corpus of review records: it refuses 17 committed ACCEPT records and accepts 12 documents that are not accepting review records at all, 5 of which are evaluator prompts. That is the round-2 pattern again — a fix correct against the one reported case and wrong in the wider context — and this time the prescription itself rested on a survey of only two files.

No model was run, no private grade obtained, no validation outcome observed. `may_open_validation` and `may_open_held_out` were false at the start and end of this review and I did not touch them. I made no change of any kind to the repository; my only writes were to the session scratchpad and to this record.

## 1. Pinned inputs — digests I computed

Digest rule applied: CRLF→LF, then CR→LF, then SHA-256. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

| Input | Expected | Computed | Match |
| --- | --- | --- | --- |
| `git rev-parse HEAD` | `60c6aee55b71716743b30d084a57daeae48b3373` | `60c6aee55b71716743b30d084a57daeae48b3373` | yes |
| `HEAD~1` (second review record) | `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` | `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` | yes |
| `HEAD~2` (superseded revision 2) | `e92eaacf6ab54b547b331d14749717da64b58e0d` | `e92eaacf6ab54b547b331d14749717da64b58e0d` | yes |
| `HEAD~3` (first review record) | `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` | `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` | yes |
| `HEAD~4` (superseded payload 1) | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | yes |
| `HEAD~5` (decision 0025) | `2cf99331688f4705ed7d6cd3efed98941f59defc` | `2cf99331688f4705ed7d6cd3efed98941f59defc` | yes |
| `artifacts/gate-1a-execution-resilience-candidate.json` **raw** | `sha256:9b581f39…0078e8` | `sha256:9b581f396d2ed8f8c285141ff90d613b6033fbf9555c7249c414eb7a2d0078e8` | yes |
| `artifacts/execution-resilience-candidate.md` **raw** | `sha256:41b9e227…03bbee` | `sha256:41b9e227b094cecb331699a60ece587e162cadfdb9dae55a53426de3c403bbee` | yes |
| Inventory file count | 16 | 16 | yes |
| Ordinal-path set digest | `sha256:473024fb…e1776147` | `sha256:473024fb57ca62a3be4ed032e05a2e37835c76ed60e10d1f72292559e1776147` | yes |

All 16 per-file LF-normalized digests recomputed and matched the inventory:

| Path | LF-normalized SHA-256 | Match |
| --- | --- | --- |
| `experiments/frontier-v1/README.md` | `sha256:fc748a2ffab7e4cf4256edb232bd5682802835391f2704e845f11082f07ed45b` | yes |
| `experiments/frontier-v1/corpusctl/internal/corpus/protocol_test.go` | `sha256:4495d032c7d86fbfa03952ad676f3e769de4ca1f0d071dc694afa376742cc952` | yes |
| `experiments/frontier-v1/pre-validation-artifacts.json` | `sha256:ce14d55cdbb36c9102e71830101eb45f0bd29df91af35a51134c57b72b3f4cc2` | yes |
| `experiments/frontier-v1/protocol.json` | `sha256:726b68112108542ba3491fde538342399cbb446a04755a7e0208e076cbdb0d32` | yes |
| `experiments/frontier-v1/runner/artifact_freeze_test.go` | `sha256:3f62b2e97e52f08296b1383c7f159b73da8e38b646219874662e9805b234ff95` | yes |
| `experiments/frontier-v1/runner/execution_resilience.go` | `sha256:fe98996728d7ccbd4b1bbd23f44d29e3ebb543a64e2430502a7ccbfb2dee374a` | yes |
| `experiments/frontier-v1/runner/execution_resilience_test.go` | `sha256:9cf56873fb0bdcc0ad229c101b3d919b631ddb1ef0520241a72311833de34957` | yes |
| `experiments/frontier-v1/runner/main.go` | `sha256:192d275ef16be74f9a3f47747cf3d14d8df52d18d16127734db786866c5bdff3` | yes |
| `experiments/frontier-v1/runner/process_events.go` | `sha256:605a365310aebb4f98c5e75fca8ae96232e671a48381756557f7c891ec70c617` | yes |
| `experiments/frontier-v1/runner/scheduled_resume.go` | `sha256:aa050eb443a6676a1ab540788c182a3ef1ec4e63ae2c7bfd4fdb4b23fac9fc2e` | yes |
| `experiments/frontier-v1/runner/scheduled_resume_test.go` | `sha256:a090f71facb3d21ba3e7d1df32988cb7a64d277b5bdba97a46637bf435c2b893` | yes |
| `experiments/frontier-v1/runner/scheduled_run.go` | `sha256:0956afad8e186b1bef71a64c44cab9c16aa64da6c8ee82cb9188fa49fab4cb9a` | yes |
| `experiments/frontier-v1/runner/types.go` | `sha256:fce5d7b956df609fe56bf77e9e400eed0f79f23374778915050948035ba92191` | yes |
| `experiments/frontier-v1/runner/validation_execution_test.go` | `sha256:9ef5e4c7e927082d44830bba29016b446b43bfa3c7a37974e18a850dc012d20e` | yes |
| `experiments/frontier-v1/runner/validation_gate.go` | `sha256:0bdaa343872454c8d0387bbf0de3b4cc2e7a3b20f9dcecd218f4c62d1adafe80` | yes |
| `experiments/frontier-v1/validation-execution-boundary.md` | `sha256:a3c32fa33645a4e287d7476f8fad3e260451b7baea6845e79501c40dfc6b6d12` | yes |

**Working tree.** `git status --porcelain --untracked-files=all` reports exactly one entry, `?? docs/reviews/2026-08-27-execution-resilience-payload-review-3-prompt.md`. `git diff` and `git diff --cached` are both empty. **No P0 on this ground.**

**Unmoved identities**, verified in `pre-validation-artifacts.json` and by the suite staying green on `TestLocalArtifactCandidateMatchesImplementation` and `TestValidationBuildReproducesFrozenWorldBuildDigest`: world build `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4`, grader digest `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b`. Both match the inventory's `world_build_unchanged` / `grader_digest_unchanged` claims.

## 2. Test and quality command results

All run by me from the repository root at the pinned commit.

| Command | Result |
| --- | --- |
| `go test ./...` | **FAIL, exactly one red:** `TestPreValidationFreezeMatchesAcceptedCandidates` (`artifact_freeze_test.go:153`, README digest `fc748a2f…` vs freeze `b486fc55…`). Every other package `ok`. |
| corpusctl suite (`cd experiments/frontier-v1/corpusctl && go test ./...`) | ok |
| `go vet ./...` | clean, no output |
| `staticcheck` v0.7.0 `./...` | clean, no output |
| `gocyclo -over 15` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis` | clean, no output |
| `dupl -t 100` (via `quality-check no-output`) | clean, no output |
| `make validate-spec` | 9 examples valid |
| `make validate-authoring` | `valid …/manifests/authoring.json` |
| **`go test -race ./experiments/frontier-v1/runner/`** on WSL2 | **I ran this myself**, `wsl -e bash -lc 'cd /mnt/d/Work/personal/phoenix && go test -race ./experiments/frontier-v1/runner/'`. **No `DATA RACE` report.** Single red, the same expected `TestPreValidationFreezeMatchesAcceptedCandidates`. 42.8s. |

**No failure beyond the single expected freeze red, on any command, on either host.** The payload's claim that the race-detector gap is discharged is confirmed independently, not restated.

I also measured gocyclo directly (`-top 8`, runner package) to check the corrected P3-N10 claim: `validateExternalGrade` 14, `runScheduledCases` 14. The inventory's `contract.revision.carried[0]` now names exactly those two at 14. **Accurate.**

## 3. Disposition of every prior P1 and P2 finding

### First review (`…-payload-review.md`)

| Finding | Disposition |
| --- | --- |
| **P1-1** termination handler swallows the signal | **Still fixed.** `watchProcessTerminationOn` (`process_events.go:178-199`) records then calls `release()` then `exit(terminationExitCode(received))`; `terminationExitCode` returns 143 for SIGTERM, 130 otherwise. Production call site passes `os.Exit`. Now exercised on every host via the injected channel, not skipped on Windows. |
| **P1-2** single-resume bound launderable downward | **Still fixed.** `verifyResumeCount` (`execution_resilience.go:195-210`) derives `attested = starts - 1` from the log and refuses when `attested > recorded`; `readProcessEventLog` returns an error for a missing, unreadable, undecodable, or zero-start log, so the attack fails closed rather than defaulting to zero. |
| **P2-1** outcome-free guard depth-1 | **Fixed, and now complete.** `assertPartialSummaryFieldsAreOutcomeFree` descends into `artifact_digests` against `outcomeFreeDigestFields` and checks `arm_assignment_counts` keys against `phase1Arms`; `assertExactFieldSet` tests membership in both directions. |
| **P2-2** counter round trip untested | **Still fixed.** `TestScheduledResumePersistsTheResumeCounter` present and green. |
| **P2-3** 72-hour window entirely self-reported | **Still fixed**, and the branch that closes it is now exercised (see P3-N4). Anchored to `lastEvent` and to the injected runner clock. |
| **P2-4** closure-record check is existence only | **Fixed in form** — non-empty regular `.md` under `docs/reviews/` whose text states a verdict — **but the replacement verdict parser is defective in both directions. See P1-N1 below.** This is the one place where a prior fix has been replaced by something that does not hold. |
| **P2-5** `verifyResumeAttestationFields` at gocyclo 14 | **Still fixed.** Split into identity/conditions/digests/log halves; measured well below the gate. |
| **P2-6** termination progress hardcoded zero | **Still fixed.** `liveProgress` mutex-guarded, advanced at each `recordCheckpoint` and at `openScheduledProgress`; wired through `progress.live.get` in `runScheduledCases:44`. |

First-review P3s: P3-3, P3-4, P3-7, P3-8, P3-9 verified fixed. P3-1, P3-2, P3-5, P3-6 carried as residuals — see §7.

### Second review (`…-payload-review-2.md`)

| Finding | Disposition |
| --- | --- |
| **P2-N1** `interrupted_at` anchor has no skew allowance | **Genuinely fixed.** `execution_resilience.go:368` now reads `stamps["interrupted_at"].Before(lastEvent.Add(-resumeClockSkew))`. I walked all four frozen causes end to end (§4) and an honest custodian passes in every one. Both sides pinned by `TestValidationResumeAcceptsAnHonestlySkewedInterruptionStamp` (2s early accepted, 24h early refused). |
| **P2-N2** boundary document describes the previous payload | **Genuinely fixed.** `validation-execution-boundary.md:63-66` now names all three added obligations, states that the 72 hours are measured from the log's last entry, and labels the one-hour freshness bound **"not a section 9 requirement … an implementation control"**. Mirrored in `README.md:179`. |
| **P3-N1** coordinated log rewrite | Carried as residual. Correctly characterized — see §7. |
| **P3-N2** guard blind to `omitempty`, never saw production shape | **Fixed.** `partialArtifactDigests` (`execution_resilience.go:41-46`) has a fixed four-key shape with no `omitempty`; all four keys admitted; `runResilienceSchedule` sets `graderDigest` and `armBDocumentDigest` so the guard sees the validation-path shape. |
| **P3-N3** `strings.Contains(contents, "ACCEPT")` satisfied by a REVISE record | **Fixed for the reported case, newly broken in the general case.** The bare-substring hole is closed and this repository's two committed REVISE records are correctly refused. But the replacement misclassifies the majority of the repository's review records in both directions. **Escalated to P1-N1.** |
| **P3-N4** the branch that closes P2-3 is untested | **Fixed.** `TestValidationResumeAnchorsItsWindowToTheProcessEventLogAndTheClock` gains the `"window measured from the last process event"` case asserting `"after the last process event"`. |
| **P3-N5** P1-1 test skips on Windows | **Fixed.** `watchProcessTerminationOn` takes the channel; `TestProcessTerminationStopsTheRun` and `TestProcessTerminationOnSIGTERMExitsWith143` run on both hosts. No `SKIP` in the runner package on either host. |
| **P3-N6** error-path stop event records zero spend | **Fixed.** `scheduled_run.go:60` now passes `SpentUSD: summary.SpentUSD`. |
| **P3-N7** dead assignment | **Fixed**, line removed. |
| **P3-N8** `append` unsynchronized | **Fixed**, and **fixed without breaking the anchor** — see §5. |
| **P3-N9** lapse explanation demanded on every resume | **Fixed.** `promptResumeWindow` (1h) gates the demand; `TestPromptResumeNeedsNoLapseExplanation` pins that a prompt custodian writes none. |
| **P3-N10** inaccurate gocyclo claim in frozen inventory bytes | **Fixed**, and I independently measured 14/14 for the two named functions. |

**Nothing from either prior review has regressed.** The one broken item is new code, written this round.

## 4. The four timing constants, walked end to end

`resumeWindow` 72h, `resumeClockSkew` 5min, `resumeClockTolerance` 1h, `promptResumeWindow` 1h. An honest custodian at invocation time `now`, with `L` = `max(recorded_at)` over the log, must satisfy: `I ≤ C ≤ V ≤ A`; `A−I ≤ 72h`; `A−V ≤ 1h` or a lapse explanation; `I ≥ L−5min`; `A−L ≤ 72h`; `A ≤ now+5min`; `now−A ≤ 1h`.

| Cause | What the log's last entry is | Honest `interrupted_at` vs `L` | Passes? |
| --- | --- | --- | --- |
| `host_restart` (systemd SIGTERM on the prepared WSL2 host) | `termination`, stamped at `T+ε` nanoseconds | `I = T` at second precision, so `I` is ≤ ~1s below `L` | **Yes** — inside the 5-minute allowance |
| `host_restart` (uncatchable Windows shutdown) | last `checkpoint`, before `T` | `I > L` | **Yes** |
| `power_loss` | last `checkpoint` | `I > L` | **Yes** |
| `hardware_failure` | last `checkpoint` (or a `stop` on an evidence-write failure) | `I ≥ L` | **Yes** |
| `process_kill` (SIGINT/SIGTERM) | `termination` at `T+ε` | `I = T`, ≤ ~1s below `L` | **Yes** |
| `process_kill` (SIGKILL) | last `checkpoint` | `I > L` | **Yes** |

**Is five minutes right?** Yes, and it is safe at any size. Moving `I` *earlier* only makes `A−I ≤ 72h` harder to satisfy; it buys an attacker nothing, because the ordering constraint is the only other place `I` appears. The attack direction is moving `I` *later*, and that is not constrained by this check at all — it is constrained by `A−L ≤ 72h`, which the custodian cannot influence. I confirmed the attack is refused: an interruption at day 0 re-attested at day 7 with a self-consistent one-hour timeline passes ordering, the 72-hour arithmetic, promptness, freshness and the interruption anchor, and is refused by `authorized.Sub(lastEvent) > resumeWindow` at `execution_resilience.go:372`. That branch is the one P3-N4 now pins.

**Erosion of the 72 hours.** Measuring from `L` rather than from `I` shortens the real window by `I−L`, which for the abrupt causes is the time from the last completed pairing-key group to the interruption — bounded by one group, i.e. ~4 trials × 180s plus grading, on the order of 15 minutes against 72 hours. Negligible, conservative in the correct direction, and **disclosed** in the frozen boundary document (`:65`). Correct.

**The two one-hour constants do not collide.** `resumeClockTolerance` bounds `now−A`; `promptResumeWindow` bounds `A−V`. A prompt custodian sets `V`, then `A` a few minutes later, then invokes: both satisfied, no lapse explanation, as `TestPromptResumeNeedsNoLapseExplanation` pins. The pathological path review 2 identified — a lapse explanation demanded on every resume — is gone.

**Attempts to defeat it that failed:** unattested resume; OOM or empty cause; second resume via the checkpoint counter; second resume via a downward-edited counter; back-dated interruption; pre-prepared stale authorization; future-stamped authorization; attestation naming a different output directory; `process_kill` without `kill_principal_established_as_non_participant`; schedule or world-build digest divergence; log deleted, truncated, zero-start, or containing a semantically empty line. All refused. The one that still works is the coordinated log rewrite, correctly carried as residual P3-N1.

## 5. Interaction defects among the revision-3 changes

I checked each revision-3 change against the assumptions the other checks depend on, which is where round 2's P2-N1 came from.

- **The mutex in `append` (P3-N8 fix) does not move the anchor.** `recorded_at` is stamped by `log.clock()` at `process_events.go:94`, **before** `log.mu.Lock()` at `:103`. So lock contention cannot inflate a termination event's timestamp past the true instant, and the anchor's tolerance is not consumed by an fsync the handler waited on. Separately, because timestamps are taken before the lock, file line order can differ from timestamp order — and `readProcessEventLog` takes `max(recorded_at)` (`:244-246`) rather than the last line, so it is order-independent. **Clean in combination.**
- **The injectable signal channel (P3-N5 fix) does not weaken the stop.** `watchProcessTerminationWith` still installs `signal.Notify` and passes `os.Exit`; the seam only adds a parameter. The returned stop function and the goroutine both call `release()`, which is `signal.Stop` and is idempotent; `close(done)` runs once from a single `defer`.
- **The `omitempty` removal (P3-N2 fix) does not leak a field.** Because all four keys are now always present, `assertExactFieldSet`'s both-directions test is meaningful on every path, and the fixture additionally populates the two validation-path digests. The removal changes the partial summary's bytes on the authoring path (two new always-present empty-string keys) but nothing reads that file for control flow — `scheduledOutputEntryKind` classifies it as progress and `loadScheduledResume` never decodes it.
- **The promptness threshold (P3-N9 fix) does not open the lapse duty.** `A−V > 1h` still demands the explanation; §9's own duty is on `A−L`, which the 72-hour backstop enforces separately.
- **The error-path spend (P3-N6 fix) appends a `stop` event after a failure, which freshens `L`.** I checked whether that lets a stale interruption look fresh: it cannot, because the stop is written at the moment of failure, so `L` still equals the moment the process died. It does mean `I` for an evidence-write failure sits at or after `L`, which passes. No interaction defect.
- **The one that does interact: `os.Exit` from the handler (P1-1 fix) against non-atomic `writeJSON`.** Recorded as P3-N1(new-a) below. Review 2 flagged this as a sharpening of P3-5; I agree with that grading and am carrying it, not escalating it.

## 6. New findings

### P1-N1 — The closure-record verdict parser misclassifies this repository's own review records in **both** directions

`experiments/frontier-v1/runner/validation_gate.go:130-148` (`verifyClosureReviewVerdictLine`) and `:152-165` (`closureReviewVerdict`); test at `experiments/frontier-v1/runner/execution_resilience_test.go:481`.

I reimplemented `closureReviewVerdict` and `verifyClosureReviewVerdictLine` verbatim in a standalone program importing nothing from the repository, and ran it against **every** `.md` file in `docs/reviews/` (76 files) plus a set of constructed cases.

**Direction 1 — it refuses legitimate ACCEPT records.**

The dominant house style for a verdict line in this repository is a bolded `Verdict:` label followed by the token in a code span. **19 committed records use exactly that form.** The parser refuses **17 of them**; the other two survive only because they happen to carry a *second*, differently formatted line (`**Final verdict: ACCEPT** — …`). The cause is at `:158`:

```go
rest := strings.TrimLeft(line[index+len("verdict:"):], " \t*_")
```

The cutset is `" \t*_"`. It does not contain a backtick. For a code-span token the remainder after trimming still begins with a backtick, and `strings.HasPrefix` against `ACCEPT` is false. Refused with *"does not state an ACCEPT verdict"*.

Three further legitimate forms are refused:

- **The current house style**, used by the two most recent v5 review records (`2026-08-27-absence-acceptance-payload-review.md`, `2026-08-27-protocol-v5-amendment-payload-review.md`): a `## 1. Verdict` heading with `**ACCEPT.**` on the following line. The parser is line-local, so the verdict word is never on the same line as the label. This is precisely the "verdict word split across lines" case the review prompt named.
- **A verdict in a table row** (`| closure | ACCEPT |`). Refused.
- **An ACCEPT record that cites a prior round's REVISE verdict.** `:137-139` returns an error on the *first* `REVISE` verdict line found, regardless of position, so an ACCEPT statement followed by "the first review returned **Verdict: REVISE**; all findings are fixed" is refused. Given this project's own multi-round review pattern — and that a decision-0023 closure review would naturally recount the execution's history — this is a likely, not exotic, shape.

**Direction 2 — it accepts documents that are not accepting review records.**

Any line containing `verdict:` followed by `ACCEPT` passes, including instructional boilerplate and quotation. **12 files under `docs/reviews/` are accepted by the parser**, and **not one of them is an accepting review record.** At least five are evaluator *prompts* whose only matching line is the instruction `1. Verdict: ACCEPT, REVISE, or REJECT.` — e.g. `2026-08-27-protocol-v5-proposal-review-prompt-v4.md:35`, `2026-08-27-absence-acceptance-payload-review-prompt.md:46`, `2026-08-27-protocol-v5-amendment-payload-review-prompt.md:46`. The review-3 prompt file itself is accepted, on the line where it quotes `"Verdict: ACCEPT"` while *warning about this exact hole*.

**What each direction permits.**

- **The false accept permits an invalid resume of the retry channel.** `validation_execution_closure_review` naming `…-closure-review-prompt.md` instead of `…-closure-review.md` — a one-token slip in a frozen field — satisfies the post-closure gate although no review has been performed. That is the precise class the check exists to catch, and it is the mechanism freeze obligation 5 requires ("the post-closure audit-before-retry gate encoded in the frozen gate rules"). With the check vacuous for this class, the gate reduces to the chair-discipline residual it was meant to replace.
- **The false refusal refuses a valid retry.** When the decision-0023 closure review is eventually written and accepted in the project's own house style, `requireValidationGate` refuses to open validation, and repairing it means editing frozen bytes — a full close → payload → independent review → refreeze cycle, at the point of maximum schedule pressure.

**Why P1 and not P2.** It defeats, in both directions, the only automated control standing between an engineered-indeterminate close and a retry; it is about to be frozen into bytes that cost a full cycle to change; and it is demonstrably wrong against the majority of this repository's own corpus rather than against a hypothetical. I record the mitigations honestly: `requirePostClosureAudit` is unreachable while `may_open_validation` is false (`:68-71`), and the gate document's two fields are frozen and chair-set, so exploitation needs a chair error plus a refreeze that is itself reviewed. Those mitigations are why this is not P0.

The candidate note's claim at `execution-resilience-candidate.md:98` — *"tolerates markdown emphasis and title case"* — is an overclaim: it tolerates `*` and `_` but not the backtick code span this repository actually uses, and the claim is about to be frozen alongside the code.

**Fix.**

1. Add the backtick and the double quote to the `TrimLeft` cutset at `:158`, and trim the same set from the right of the extracted token, so a code-span or quoted token resolves.
2. Make the match structural rather than substring-positional: require that the text *preceding* `verdict:` on the line consist only of markdown furniture (`#`, `-`, `*`, `_`, backtick, whitespace) — which admits `- **Verdict:**`, `## Verdict:` and `**Final verdict:**` while rejecting the numbered instruction `1. Verdict: …` and mid-sentence quotation.
3. Require the token after the colon to be a lone verdict word: after trimming, the remainder must begin with the verdict token and must not also contain the other verdict token, which rejects `Verdict: ACCEPT, REVISE, or REJECT` outright.
4. Treat only the record's **first** verdict statement as decisive, so a later citation of a prior round's REVISE does not overturn a stated ACCEPT.
5. Optionally accept a `Verdict` *heading* whose next non-blank line is a lone verdict token, which is the current house style.
6. Extend `TestPostClosureGateRequiresARealReviewRecord` with a case for each: the code-span form, the heading-plus-next-line form, an ACCEPT record citing a prior REVISE, and the evaluator-prompt boilerplate (must refuse). The present test pins only `**Verdict: Accept** - no findings`, a form no committed record uses.

### P3-N1(new-a) — A supervised stop can truncate the checkpoint, because `writeJSON` is neither atomic nor fsynced

`experiments/frontier-v1/runner/runtime.go:475-484`, reached from `execution_resilience.go:482` and `scheduled_run.go:44-45`.

`writeJSON` is `os.WriteFile`: truncate-then-write, no temp-and-rename, no fsync. The signal goroutine calls `os.Exit` asynchronously. A SIGTERM arriving during `writeScheduledCheckpoint` leaves a truncated or zero-length `scheduled-checkpoint.json`; `decodeStrict` then fails and `loadScheduledResume` refuses the whole directory.

**Refuses a valid resume** (it cannot permit an invalid one — every failure mode here is a decode error). The window is one small-file write against a multi-minute pairing group, so the probability per interruption is on the order of 10⁻⁴. This is review 2's P3-5 with its sharpened wording, and I reached the same conclusion independently rather than adopting it: acceptable as a residual, worth the chair knowing that a temp-and-rename plus fsync in `writeJSON` removes the class in a few lines and would also close P3-5 proper. **Fix:** write to `path+".tmp"`, `Sync`, `Close`, `os.Rename`.

### P3-N1(new-b) — The frozen operational records never name the attestation's field set, and `decodeStrict` rejects anything else

`experiments/frontier-v1/validation-execution-boundary.md:59-70`, `experiments/frontier-v1/README.md:179`, against `execution_resilience.go:119-140` and `decodeStrict` (`runtime.go`, `DisallowUnknownFields`).

Area D asks whether a custodian following the boundary document literally can construct an attestation the runner accepts. **On the timing rules, yes** — I walked it and conditions 1–5 and all three anchors are now stated accurately and completely, which is why P2-N2 is closed. **On the document shape, no.** The boundary document states the *obligations* but never the JSON field names, and `decodeStrict` uses `DisallowUnknownFields`, so an attestation is refused until every one of `no_outcome_inspection_between_interruption_and_resume`, `cause_classified_before_inspection_and_resume_decision`, `frozen_bytes_digest_identical_at_resume`, `kill_principal_established_as_non_participant`, `frozen_digests_verified.{schedule_digest,world_build}`, `process_event_log`, `process_event_log_sha256`, `promptness_statement`, `resume_index`, `output_dir` and the four timestamps is spelled exactly right. The only place they are enumerated is the table in `execution-resilience-candidate.md:39-47`, which is a candidate note and **is not among the 16 inventoried files**, and the Go source.

**Refuses a valid resume** — recoverably, by iterating against error messages, but inside a window the payload itself has bounded to one hour. **Fix:** add a worked attestation template, or the literal field list, to `validation-execution-boundary.md` so the frozen operational contract is self-sufficient. This is the same class as P2-N2 and belongs in the same document.

### P3-N1(new-c) — A first custodian process that dies between its `start` event and its first checkpoint leaves the directory unresumable

`experiments/frontier-v1/runner/scheduled_resume.go:180-182`, against `scheduled_run.go:53-57`.

`openScheduledProgress` writes the `start` event at `scheduled_run.go:104`; the initial checkpoint is written at `:54`. If the process dies in between, the directory has `process-events.jsonl` and no checkpoint, so `finishScheduledInventory`'s progress-implies-checkpoint rule refuses the whole directory with *"scheduled resume is missing a checkpoint"* — for a validation tranche, with no path forward short of a chair decision.

**Refuses a valid resume.** The window is milliseconds and the rule is what prevents the worse alternative (a phantom `start` inflating `attested` on a later genuine resume, permanently refusing it), so the design is right and only the outcome is harsh. I am **not** asking for a code change; I record it so the chair knows the state exists and what it looks like. Worth one sentence in the boundary document.

**I confirm the authorization ordering is correct**: `requireResumeAuthorization` runs at `scheduled_run.go:83`, before the log is constructed at `:87` and before the `start` event at `:104`. A *refused* resume therefore writes no `start` event and does not consume the tranche's single resume, and does not perturb the attested log digest. This was the highest-risk instance of the round-2 pattern and the payload gets it right.

### P3-N1(new-d) — Stale round-count text in the inventory bytes about to be frozen

`experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json:127`, `review_required`: *"second independent review after the first returned REVISE…"*. This is the third review, against a revision that two reviews have already returned REVISE on. Same class as P3-N10 — an inaccurate claim in bytes about to be pinned by raw digest at refreeze. **Permits nothing and refuses nothing.** **Fix:** reword to "third independent review after the first two returned REVISE".

## 7. Judgment on each area-E residual

- **P3-N1 (coordinated log rewrite).** Correctly characterized and genuinely acceptable. I re-derived it: `verifyAttestedProcessEventLog` (`:308-320`) compares the attestation against the bytes currently on disk, and no previously accepted digest is retained anywhere, so deleting a `start` line or appending a fresh entry defeats both the counter and the anchor in one move. Review 2's grading — an enhancement beyond §9, detectable only by comparing committed digests across commits — is right. Its proposed fix (record the accepted digest *and byte length* in the checkpoint and require prefix-consistency) is sound and cheap, and I would take it if the P1 fix reopens these files anyway, but I do not make it a condition.
- **The host-clock residual.** Correctly characterized and acceptable. I verified the asymmetry: a clock set *before* the last process event is refused by the interruption anchor; a clock set shortly *after* it passes every timing check. No trusted time source exists on the execution host, so no runner check can reach this. Keep, with review 2's wording.
- **The post-closure gate's chair-discipline caveat.** Correctly characterized, and review 2's sharpening is right and load-bearing: nothing in the runner writes `validation_execution_status`. **But P1-N1 makes it more load-bearing still.** The chair-discipline residual currently reads as "the chair must arm the gate"; with the parser defective, it must also read "and the runner's corroboration of the named record is unreliable in both directions". If the chair imports without fixing P1-N1, that sentence must be added. If P1-N1 is fixed, review 2's wording stands unchanged.
- **The wipe-and-restart caveat.** Correctly characterized and genuinely acceptable — I re-confirmed the mechanism: an emptied output directory returns `scheduledResume{}` from `loadScheduledResume:44-46`, `requireResumeAuthorization` returns at `:171`, and the resume rule never engages. §9 provides no runner-side control for this; the post-closure gate is the control, which is exactly why P1-N1 matters. Keep review 2's sharpened wording verbatim.
- **The authoring-tranche exemption.** Correctly characterized and genuinely acceptable. I re-verified the decisive argument independently: `validateScheduleForTranche` refuses a schedule whose tranche does not match the flag, so the exemption cannot launder a validation execution, and `TestAuthoringResumeIsExemptFromTheAttestationRule` pins it as a decision rather than an oversight. Review 2's added sentence about the round-trip test running on the authoring path is accurate and should be kept.
- **P3-1 (committed log digest never covers the terminal entry).** Correctly characterized and acceptable. `recordStop` appends after the last `writeScheduledPartialSummary`, so the partial summary's `process_event_log_sha256` is always one entry stale; the custody record hashes the whole file at closure. With P3-N6 fixed, the terminal entry now at least carries a true spend.
- **P3-2 (five identity fields beyond §9's literal list).** Correctly characterized and acceptable. I re-checked each of `v`, `tranche`, `status`, `schedule_digest`, `resumes`: none is a measurement, and each is needed to bind the record to a tranche.
- **P3-5 (checkpoint and partial summary not fsynced or atomically renamed).** Correctly characterized **with review 2's sharpened wording**, and only with it — the original wording understated the exposure, because the `os.Exit` in the termination handler makes a mid-write truncation a consequence of a *supervised stop* and not only of an abrupt external kill. See P3-N1(new-a); acceptable as a residual, and I concur with review 2 that the chair should know a temp-and-rename removes the class.

**One observation, not a finding.** The inventory's `contract.accepted_residuals` array carries 6 entries, while the note's prose carries the rest. The refreeze transcribes from the review record, so §8 below is the operative list; the array's brevity has no operative consequence.

## 8. Residual list for the refreeze that follows the P1 fix

The verdict is REVISE, so **there is no refreeze to attach findings to and this list is not operative.** For the refreeze that follows a P1-N1 fix, these are the entries I would accept into `accepted_findings`, worded for direct transcription:

- "Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian's written attestation, as section 9 requires."
- "Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b)."
- "Residual: the process-event log is unsigned custodian-controlled local state, so a coordinated rewrite that deletes a start entry and re-attests the new digest defeats both the single-resume corroboration and the 72-hour anchor; the runner retains no previously accepted digest and cannot detect it. Detection is the committed attestation digest compared across commits, under independent review."
- "Residual: the 72-hour anchor is measured against the runner's own clock, which the custodian's host controls; a clock set to shortly after the last process event passes every timing check, while a clock set before it is refused by the interruption anchor. No trusted time source exists on the execution host."
- "Residual: measuring the 72 hours from the log's last entry rather than from the attested interruption shortens the real window by the time from the last completed pairing-key group to the interruption, bounded by one group; the erosion is disclosed in the boundary document and is conservative in the correct direction."
- "Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports 'gate is closed' first; both are refusals and only the surfaced message differs."
- "Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag. One consequence: the end-to-end resume-counter round-trip test runs on the authoring path and so does not traverse requireResumeAuthorization."
- "Residual: nothing in the runner writes validation_execution_status. A custodian who wipes the output directory presents as a fresh run and bypasses the resume rule entirely, so the post-closure gate — armed by the chair setting that field at closure, under the freeze workflow — is the sole control on the retry channel after a future non-completed execution. It is armed by frozen bytes for the current one."
- "Residual (P3): the partial summary's committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure."
- "Residual (P3): the partial summary carries five non-outcome identity fields beyond section 9's literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome."
- "Residual (P3): the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log. Because the termination handler deliberately calls os.Exit, a supervised stop as well as an abrupt kill can truncate a boundary write; the failure is closed — decodeStrict refuses and the tranche closes indeterminate — but a routine Ctrl+C at an unlucky microsecond can burn the tranche. A temp-and-rename in writeJSON would remove the class."
- "Residual (P3): a validation output directory whose first custodian process wrote its start event but died before the initial checkpoint is refused wholesale by the progress-implies-checkpoint rule and cannot be resumed; the window is milliseconds and the rule is what prevents a phantom start entry from permanently refusing a later genuine resume."
- "Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first."
- "Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required."
- "Carried P2-1 is retired by this payload: classifyScheduledOutputEntries is decomposed into scheduledOutputEntryKind, recordScheduledAssignment, and finishScheduledInventory, and verifyScheduledCheckpoint is split into identity and progress halves. Its successor P2-5 is also retired: verifyResumeAttestationFields is split into identity and condition halves and measures below 5. runScheduledCases and validateExternalGrade remain at gocyclo 14, both pre-existing and unchanged by this payload."

**P1-N1 is not in that list.** It concerns frozen bytes that cannot be repaired at refreeze under this workflow, it fails in both directions, and one of those directions defeats the mechanism freeze obligation 5 exists to install. P3-N1(new-b) and P3-N1(new-d) are documentation edits I would rather see made alongside it, since both touch files the P1 fix does not.

## 9. What I did not verify

- I did not run a model, obtain a private grade, or observe any validation outcome. Both gates were false at the start and end of this review and I did not touch them.
- I did not execute the signal path against a live multi-hour scheduled run. My P1-1 and P3-N1(new-a) conclusions rest on reading `watchProcessTerminationOn`, confirming the single production call site passes `os.Exit`, and tracing the record-then-release-then-exit ordering against `writeJSON`'s truncate-then-write.
- I did not construct a live end-to-end validation resume, because that requires the validation gate open. Every resume conclusion in §4 rests on `requireResumeAuthorization` and its callees, exercised through the package's own fixtures and by hand-tracing the four causes.
- P1-N1 is the one finding I demonstrated empirically rather than by reading: I reimplemented both functions standalone and ran them over all 76 files in `docs/reviews/` and over nine constructed verdict forms. The counts (17 refused of 19 house-style records; 12 accepted of which none is an accepting record) are measured, not estimated.
- I did not attempt to defeat the freeze test or to reproduce the world-build digest from source; I confirmed the two digests are unmoved in `pre-validation-artifacts.json` and that the two tests pinning them stay green.
- I made no change of any kind to the repository. My only writes were to the session scratchpad (a standalone Go reimplementation of the verdict parser) and to this record.

---

## Chair transcription note

Transcribed verbatim from the third independent reviewer's record. Before transcription the payload author independently reproduced the two counts that carry P1-N1: 19 committed records use the code-span verdict style whose token the trim cutset could not reach, and the evaluator-prompt boilerplate `1. Verdict: ACCEPT, REVISE, or REJECT.` is present in five prompts and would have been read as an ACCEPT. The verdict is **REVISE**, so this payload does not proceed to refreeze; a fourth revision follows and is submitted for a further independent review. §8 is conditional on the P1 fix and is not operative at this verdict.
