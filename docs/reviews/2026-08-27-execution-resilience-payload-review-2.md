# Second independent review — protocol-v5 section 9 execution-resilience payload, revision 2

**Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author any file under review, and I am not the reviewer who produced the first review record. I reviewed by reading the repository at the pinned commit and running the test and quality commands.
**Date:** 2026-08-27
**Payload commit:** `e92eaacf6ab54b547b331d14749717da64b58e0d`
**Verdict: REVISE** — two P2 findings, ten P3. **No P0, no P1.** Both first-review P1s are genuinely fixed; all six P2s are fixed or substantively fixed. The two new P2s are both in the code the fixes introduced, both refuse a *legitimate* resume rather than permitting an invalid one, and both are one-line-plus-one-paragraph repairs. No model was run, no private grade obtained, no validation outcome observed; `may_open_validation` and `may_open_held_out` were false throughout and I did not touch them. I made no change of any kind to the repository.

## 1. Pinned inputs — digests I computed

Digest rule applied: CRLF→LF, then CR→LF, then SHA-256.

| Input | Expected | Computed | Match |
| --- | --- | --- | --- |
| `git rev-parse HEAD` | `e92eaacf6ab54b547b331d14749717da64b58e0d` | `e92eaacf6ab54b547b331d14749717da64b58e0d` | yes |
| `HEAD~1` (first review record) | `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` | `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` | yes |
| `HEAD~2` (superseded payload) | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | `7a1ea42513acd3d55e276eabc2459da4a037acc0` | yes |
| `HEAD~3` (decision 0025) | `2cf99331688f4705ed7d6cd3efed98941f59defc` | `2cf99331688f4705ed7d6cd3efed98941f59defc` | yes |
| `artifacts/gate-1a-execution-resilience-candidate.json` **raw** | `sha256:cb36e208…3aea47` | `sha256:cb36e2087f1d2b4e0a7d61965d472a6c9ee03ae498cf913fc2b1b1e9813aea47` | yes |
| `artifacts/execution-resilience-candidate.md` **raw** | `sha256:05799da5…3a3ed` | `sha256:05799da56672700913b0462ff457541a402c07cad9bdfaec84ede1c31a3ea3ed` | yes |
| `runner/execution_resilience.go` LF | `sha256:28d9c174…c34e14d8` | `sha256:28d9c174fb5c457271b382281329b8b12f6cd48cfb42989427fc1cf4c34e14d8` | yes |
| `runner/process_events.go` LF | `sha256:cdfa5f84…9c7792d4` | `sha256:cdfa5f84e8dd92dfc41d1904b289ce704e1000a1a8fe36504b5b85fa9c7792d4` | yes |
| `runner/validation_gate.go` LF | `sha256:df56726d…7d68bb7` | `sha256:df56726da5dc61753f1b761c09079e61174fe6c4368d7ce47a4bf1a687d68bb7` | yes |
| `runner/scheduled_run.go` LF | `sha256:e58360c3…64e427ad` | `sha256:e58360c3e83656c8cc996a577c6eed89380956c8c912e3ba380a70b664e427ad` | yes |
| `runner/execution_resilience_test.go` LF | `sha256:f39d36ed…8964546` | `sha256:f39d36ed26e8c707e10aae9de1fbb94bab8e59f1bce9a4a7d519edadd8964546` | yes |
| `pre-validation-artifacts.json` LF | `sha256:ce14d55c…2b3f4cc2` | `sha256:ce14d55cdbb36c9102e71830101eb45f0bd29df91af35a51134c57b72b3f4cc2` | yes |
| Inventory file count | 16 | 16 | yes |
| Ordinal-path set digest | `sha256:95f33f2a…511808e` | `sha256:95f33f2af58442464557b77287647f93393cdb313397721cd0bac8525511808e` | yes |

I recomputed the LF digest of **all 16** inventory files and compared each against the inventory's `files` map: **zero mismatches**. The set digest above was rebuilt from those 16 recomputed values, not copied from the document.

Identities I verified rather than trusted:

- World-build `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` and grader digest `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b`: **neither string appears anywhere in `git diff HEAD~1 HEAD`**, and `git diff --name-only` over `cmd internal verbs build worlds go.mod go.sum schema` is **empty**. `TestValidationBuildReproducesFrozenWorldBuildDigest` and `TestLocalArtifactCandidateMatchesImplementation` both PASS.
- Payload scope vs the frozen state: `git diff --name-only 2cf9933 HEAD` is **21 paths** = the 16 inventoried files + the inventory + the candidate note + the uninventoried working-tree guard `runner/execution_resilience_candidate_artifact_test.go` + the first review record and its prompt. Scope matches the declaration; nothing undeclared.
- The guard test's pinned base commit was moved from `2cf9933…` to `46030d4…`, correctly tracking the new base.

**Working tree:** exactly one uncommitted path, `docs/reviews/2026-08-27-execution-resilience-payload-review-2-prompt.md` (untracked). Nothing else is dirty. **No P0 on this ground.**

## 2. Test and quality command results

| Command | Result |
| --- | --- |
| `go test ./...` (root module) | 1 FAIL, all else ok — **only** `TestPreValidationFreezeMatchesAcceptedCandidates` (`artifact_freeze_test.go:153`, README digest `f8d657c5…` vs frozen `b486fc55…`) |
| corpusctl suite (`cd experiments/frontier-v1/corpusctl && go test ./...`) | ok |
| `go vet ./...` | clean, exit 0 |
| `staticcheck@v0.7.0 ./...` | clean, no diagnostics |
| `gocyclo@v0.6.0 -over 15` over the gated dirs | clean, no output |
| `dupl@v1.1.0 -plumbing -t 100` under `quality-check no-output` | clean, no output |
| `make validate-spec` | all 9 examples valid |
| `make validate-authoring` | valid |
| `go test -race ./experiments/frontier-v1/runner/` | **COULD NOT RUN — see below** |

**No failure beyond the single expected freeze red.** The red is the one named in the inventory's `expected_red_freeze_tests`, and it reports the runner set digest, exactly as a payload commit should.

**The `-race` result, explicitly: I could not execute it, and I am not reporting it as green.** `go test -race` on windows/amd64 requires cgo with a GNU-style toolchain. `CGO_ENABLED=1` fails with `C compiler "gcc" not found`. The only C compiler on this host is clang (`C:\Program Files\LLVM\bin\clang.exe`), and the Go linker drives it in MSVC mode: with the default linker it fails `LNK1143: invalid or corrupt file: no symbol for COMDAT section`, and with `-fuse-ld=lld` it fails on Go's `fix_debug_gdb_scripts.ld` linker script (`unknown file type`), which `cmd/link` emits unconditionally and which `-ldflags=-w` does not suppress. `-ldflags=-linkmode=internal` is refused (`runtime/race` requires external linking). Installing a MinGW toolchain would be a change to the host, which is outside a reviewer's remit. **I substituted a manual concurrency audit, reported in §5-B.** A reviewer on the WSL2/Linux host should run `-race` before the refreeze; my audit predicts it green on memory access and says nothing about the file-level hazard in P3-N8, which `-race` cannot see.

I also ran the section 9 tests with `-v`. All pass **except one skip**, which matters and is recorded as P3-N5: `TestProcessTerminationStopsTheRun` **SKIPs on this Windows host** (`os.Process.Signal(os.Interrupt)` is unsupported on Windows), so the payload's only evidence for the P1-1 fix does not execute on the recorded execution host.

Independent measurements, not taken from the note:

- `verifyResumeAttestationFields` is now **below 5** (it does not appear in the top 8 of its own file). P2-5 genuinely retired.
- `runScheduledCases` measures **14** at this commit, and **14** at `HEAD~1` (measured by extracting the parent's `scheduled_run.go` to a scratch path and running gocyclo on it). The revision's two added statements add no branch. Pre-existing, but see P3-N10.
- `pre-validation-artifacts.json`: the three added `gates` lines are now indented with 4 spaces, consistent with the block. **P3-4 fixed**, verified with `cat -A`.

## 3. Disposition of every first-review P1 and P2 finding

### P1-1 (termination handler swallows the signal) — **GENUINELY FIXED**

`process_events.go:158-190`. The handler now records the entry and calls `exit(terminationExitCode(received))`; `watchProcessTermination` supplies `os.Exit`, and that constructor is the one wired into production at `scheduled_run.go:44`. Answering each sub-question the prompt asks:

- **Does a signalled run stop?** Yes, with 130 (SIGINT) / 143 (SIGTERM).
- **Is the exit reachable in production, not only under the test double?** Yes. `watchProcessTermination` has exactly one non-test call site, `scheduled_run.go:44`, and it passes `os.Exit`. The injectable variant is used only by the test.
- **Is the write durable and does it precede the stop?** Yes. `recordTermination` → `append` → `file.Sync()` → `Close`, all before `exit`. The ordering is correct.
- **Second signal, or a signal during the exit?** The channel has capacity 1 and is drained once; `signal.Stop` is called *after* the record, so during the record-and-sync window a second SIGINT is buffered and discarded rather than killing the process. The window is one small fsync. I judge this the correct trade (the record must survive) and not a finding, but the chair should know that Ctrl+C is briefly non-forcing, and that a stalled disk would extend that window.
- **Windows at the start/end of a run.** Moving the watcher from `main.go` into `runScheduledCases` leaves a window *before* `openScheduledProgress` returns (resume authorization and the `start` write) and *after* `defer stopWatch()`, in which a signal takes the default disposition and writes no entry. Neither window can record a termination that did not happen or act on one it did not record — the failure mode P1-1 named is closed in both. A kill inside the opening window leaves the log without that process's `start`, which *under-counts* resumes; I traced this and it is benign: no launch has occurred and the checkpoint's counter is unmoved, so the allowance is correctly not consumed.
- **Anything leaked?** `os.Exit` skips deferred cleanup, including `stopWatch` and the run config's `cleanup()`. The pre-payload default disposition also skipped them, so nothing is newly leaked.

### P1-2 (single-resume bound launderable downward) — **FIXED for the reported attack; the underlying weakness is narrowed, not eliminated**

`execution_resilience.go:186-201` (`verifyResumeCount`), `:212-243` (`readProcessEventLog`). I attacked it along every axis the prompt names:

| Attack | Result |
| --- | --- |
| Edit checkpoint `resumes` 1→0 | **Refused.** `attested(1) > recorded(0)` → "the process-event log attests". This is the exact P1-2 attack and it is closed. |
| Delete the log | **Refused** — `read process-event log`. |
| Log with zero `start` entries | **Refused** — `process-event log records no custodian process start` (`:239-241`). |
| Truncate the log mid-line | **Refused** — `json.Unmarshal` fails → `decode process-event log`. Fail-closed. |
| Valid JSON but semantically wrong line (`{"foo":1}`) | **Refused** — `recorded_at` is empty, `time.Parse` fails → `decode process-event timestamp`. Every line must carry a parseable timestamp, so garbage cannot be padded in silently. |
| Inject extra `start` lines | Only ever **increases** `attested`, so it can only refuse. Not an attack. |
| Delete the checkpoint to look fresh | **Refused.** `finishScheduledInventory` (`scheduled_resume.go:178-180`) requires a checkpoint whenever assignments, evidence, **or progress records** are present, and both the partial summary and the process-event log are classified as progress (`:139`). |
| Copy assignments to a new output dir without the log | **Refused** — now also on the `read process-event log` ground. The log being load-bearing is a genuine strengthening. |

**Ordering with the digest check.** `verifyResumeCount` runs at `:165`, *before* the attestation file is even opened at `:172`; `verifyAttestedProcessEventLog` (`:299-311`) runs later, inside `verifyResumeAttestationFields`. Neither accepts a state the other rejects: the count check validates content and never looks at the attestation; the digest check binds the attestation to the bytes on disk and never validates content. The fail-fast ordering is right — a malformed log is reported as such rather than as a digest mismatch.

**The residual.** The two checks together do **not** close a *coordinated* rewrite: the log is unsigned custodian-controlled local state, so a custodian who deletes a `start` line and re-attests the new digest passes both. The runner keeps no memory of the digest it accepted at the previous resume, so it cannot detect the rewrite. Detection remains a human comparison across git commits, which is what §9 assigns to the committed attestation and the independent audit. **This is a real narrowing** — from "edit one integer that was never committed" to "forge an fsynced append-only log whose prior digest is already in git" — and the candidate note does not overclaim it (`candidate.md:46` says exactly this). Recorded as **P3-N1** with a concrete fix, not as a P1 or P2.

### P2-1 (outcome-free guard is depth-1) — **PARTIALLY FIXED**

`execution_resilience_test.go:20-22, 429-459`. The guard now descends into `artifact_digests` against `outcomeFreeDigestFields` and checks that `arm_assignment_counts` keys are Phase 1 arms. Both are real improvements and the `arm_assignment_counts` half is complete.

The nested half is not. `outcomeFreeDigestFields` admits only `world_build` and `process_event_log_sha256`, but `partialArtifactDigests` (`execution_resilience.go:37-42`) also declares `grader_digest,omitempty` and `arm_b_document_digest,omitempty`. Every fixture (`runResilienceSchedule`, `interruptedAuthoringRun`) builds a `runConfig` with neither digest set, so both are empty, both are omitted, and the guard's exact-key assertion happens to hold. On the validation path both are non-empty (`main.go:212` and `:445`), so the guard has never been run against the shape it is supposed to guard — and an `omitempty` field added to that struct still ships past it, which is the same class of hole P2-1 reported. See **P3-N2**.

### P2-2 (nothing tests the counter round trip) — **FIXED**

`execution_resilience_test.go:262-296`. `TestScheduledResumePersistsTheResumeCounter` runs a real interrupted schedule to completion and asserts the on-disk checkpoint **and** the partial summary both read `resumes: 1`. Deleting `progress.resumes++` now breaks this test. Two limits worth stating rather than grading: the round trip runs on the **authoring** path (a consequence of the exemption), so it does not traverse `requireResumeAuthorization`; and "a further resume is refused" is asserted by a direct unit call to `verifyCheckpointProgress` with a literal, not end to end. The validation-side refusal is separately covered by `TestValidationResumeRefusesASecondResume` and `TestValidationResumeCorroboratesTheCounterAgainstTheProcessEventLog`, so combined coverage is adequate. One leftover: **P3-N7**.

### P2-3 (72-hour window entirely self-reported) — **FIXED substantively; two new defects in the fix**

`execution_resilience.go:348-368`. I attacked it from both directions.

**Can an invalid resume still pass?**

- *Attestation authored at resume time for a weeks-old interruption* (the headline P2-3 attack): **refused**. Forward-dating `interrupted_at` to `now-1h` survives the ordering and `elapsed` checks, but `authorized.Sub(lastEvent)` is then ~3 weeks and trips `:357-359`. This branch is what actually closes P2-3 — and it is the one branch of the anchor that **no test exercises** (**P3-N4**).
- *Back-dating below the last event*: refused at `:353-355`.
- *Stale or future authorization*: refused at `:361-366`.
- *A later process appends to the log to freshen `lastEvent`*: cannot be done through the runner (a `start` entry is written only *after* authorization, `scheduled_run.go:104`), so it reduces to hand-editing the log — the same **P3-N1** residual.
- *Custodian controls the host clock*: **defeats the anchor**, irreducibly. Rolling the clock back to *before* `lastEvent` is refused (the `interrupted_at` anchor catches it), but rolling it to shortly *after* `lastEvent` passes everything. The runner has no trusted time source; this belongs on the residual list and I have worded it there.

**Are the constants defensible?** Yes, as values. `resumeClockSkew = 5m` is a sane NTP allowance for the future-stamp check. `resumeClockTolerance = 1h` covers only "write the attestation JSON and launch the runner" — the expensive §9(b) byte re-verification happens before `conditions_verified_at` and is outside the hour — so it is not too tight, and as an anti-forgery measure it is meaningless anyway (the custodian owns the clock); its real function is to force re-attestation at the moment of resume, which is a legitimate promptness control.

**Does the anchor break a legitimate resume?** **Yes — this is P2-N1**, below. And the new obligations it places on the custodian are absent from the frozen operational record — **P2-N2**.

### P2-4 (closure-record check is existence only) — **FIXED EXACTLY AS PRESCRIBED; the prescription is weak**

`validation_gate.go:107-127`. The record must now be under `docs/reviews/`, end in `.md`, be a regular non-empty file, and contain `ACCEPT`. `TestPostClosureGateRequiresARealReviewRecord` covers the outside-`docs/reviews`, empty, contradicting, and passing cases. The implementer did precisely what the first review asked.

The prompt asks whether requiring the string `ACCEPT` anywhere in the file is meaningful. **It is not, and I can demonstrate it from this repository.** `docs/reviews/2026-08-27-execution-resilience-payload-review.md` — a **REVISE** record — is under `docs/reviews/`, ends in `.md`, is a non-empty regular file, and contains the token `ACCEPT` **five times** (lines 125, 127, 129, 158, 207). It would satisfy `verifyClosureReviewRecord` today. Every independent review record in this project will contain that token, because reviewers are instructed to "Return **ACCEPT** or **REVISE**". The check is also over-strict in the other direction: a record writing "Verdict: Accept" in title case fails. Recorded as **P3-N3**, not higher, because both fields are frozen bytes and exploitation presupposes already subverting the freeze workflow.

### P2-5 (`verifyResumeAttestationFields` at gocyclo 14) — **GENUINELY FIXED**

Split into `verifyResumeAttestationIdentity` (7) and `verifyResumeAttestationConditions` (7); the parent measures below 5. Measured, not taken from the note. See **P3-N10** for the accompanying claim.

### P2-6 (termination progress hardcoded zero) — **FIXED IN CODE; production wiring still untested**

`scheduled_run.go:44` passes `progress.live.get`; `execution_resilience.go:471` and `scheduled_run.go:103` publish through it. Every termination entry a real run writes now carries the launch index, completed pairing keys, spend, and resume count at that moment, so the candidate note's `§1` claim (`candidate.md:20`) is now true of all four event kinds. However, the only test of the watcher passes an explicit literal closure and **skips on Windows**, so `progress.live.get` as the progress source is verified on no host — the same untested-production-wiring shape that produced the original defect. Folded into **P3-N5**.

### The five claimed P3 fixes

- **P3-3** fixed — `scheduledStopKind` (`execution_resilience.go:421-432`) reduces the reason to `safety_stop` / `run_budget` / `other`, pinned by `TestScheduledStopKindNeverCarriesTheRawFailureText` including the grade-shape string the first review quoted.
- **P3-4** fixed — verified with `cat -A`; the added `gates` lines now use the block's 4-space indent.
- **P3-7** fixed — `verifyResumeFrozenDigests` (`:289-291`) refuses on an empty live digest instead of skipping.
- **P3-8** fixed — `verifyAcceptedLocalArtifacts` hoisted above `verifyReplacementRunnerFreeze` (`artifact_freeze_test.go:152`). It is now genuinely exercised at a payload commit; the first review's independent replication is no longer needed.
- **P3-9** fixed, but the fix carries a new defect — see **P3-N6**.

## 4. New findings

### P2-N1 — The `interrupted_at` anchor has no skew allowance and will refuse a legitimate resume of a signalled interruption

`experiments/frontier-v1/runner/execution_resilience.go:353-355`.

```go
if stamps["interrupted_at"].Before(lastEvent) {
    return fmt.Errorf("attested interruption precedes the last process event at %s", ...)
}
```

Every other clock comparison in this function carries a tolerance (`resumeClockSkew` at `:361`, `resumeClockTolerance` at `:364`). This one has none, and it compares two quantities that are guaranteed to sit on opposite sides of the true instant:

- `lastEvent` is `max(recorded_at)` over the log, at **RFC3339Nano** precision. When a catchable signal ends the run, the last entry is the `termination` event, whose `recorded_at` is stamped *after* the signal arrives — `T + ε`.
- `interrupted_at` is what the custodian attests as the moment of interruption — naturally `T`, taken from the system journal or the supervisor — parsed at `:321` with `time.Parse(time.RFC3339, …)`, which custodians will in practice write at **second** granularity.

So for the two frozen causes that produce a catchable signal — `host_restart` (systemd SIGTERM on the prepared WSL2/Linux host) and `process_kill` — an honest custodian who attests the true interruption instant writes a value strictly below `lastEvent` and is **refused**, with a message that does not tell them what to write instead. Truncating `T + ε` to seconds fails identically. The only passing values are ones the custodian must invent by reading the log's `recorded_at` and rounding up — an instruction that appears nowhere (see P2-N2). `power_loss` and `hardware_failure` are unaffected, because their last entry is an earlier `checkpoint`.

**Failure it permits:** a refused resume of a valid interruption. Because the resume rule is mandatory and the runner is the enforcement point, a custodian who cannot get past this check burns clock against the 72-hour backstop and, if unresolved, closes a ~300 USD validation tranche indeterminate. The prompt is explicit that this is as much a defect as permitting an invalid resume, and I agree.

**Fix:** allow the same skew already defined for the other direction — `if stamps["interrupted_at"].Before(lastEvent.Add(-resumeClockSkew))` — and add a test case pinning that an `interrupted_at` a few seconds before the `termination` entry is accepted while one a day before is refused.

### P2-N2 — The frozen boundary document does not state the two custodian obligations the revision added

`experiments/frontier-v1/validation-execution-boundary.md`, section "Execution resilience, mandatory resume, and the post-closure gate", numbered condition 3; and `experiments/frontier-v1/README.md:179`.

Condition 3 reads, in full: *"resume is authorized within 72 hours of the interruption, with the timestamps ordered interruption, cause classification, condition verification, authorization"*. That was an accurate description of the first payload. It is no longer an accurate description of the runner. The revision added two refusal conditions that the custodian must satisfy and that appear in **no** frozen record:

1. `interrupted_at` may not precede the last `recorded_at` in `process-events.jsonl` (`execution_resilience.go:353`).
2. `resume_authorized_at` must be within `resumeClockTolerance` — **one hour** — of the moment the runner is invoked, and no more than `resumeClockSkew` ahead of it (`:361-366`). This obligation has **no basis in §9 at all**; it is a payload invention. §9 imposes a promptness duty and a 72-hour backstop, not a one-hour re-attestation window.

`protocol.json`'s `amendment.execution_resilience` is the statement of the *rule* and is correctly silent on implementation, so it is not at fault. But the boundary document is the operational contract the custodian follows at the moment of resume, it is one of the 16 inventoried files, and it is about to be refrozen. A custodian following it exactly can be refused — and combined with P2-N1, refused twice.

**Failure it permits:** the frozen operational record understates the runner's requirements, so a compliant custodian is refused during a time-critical window; and a payload-invented obligation enters the frozen bytes undocumented and unjustified against §9.

**Fix:** extend condition 3 to name both anchors and the one-hour freshness requirement, state plainly that `interrupted_at` must be at or after the log's last `recorded_at`, and justify the one-hour window as a promptness control rather than letting it read as a §9 requirement. Mirror the substance in the README's resume paragraph.

### P3-N1 — A coordinated rewrite of the process-event log defeats both the resume counter and the timing anchor, and the runner keeps no state that would detect it

`experiments/frontier-v1/runner/execution_resilience.go:186-201` and `:348-368`.

The corroboration and the anchor both derive from a file that the custodian can edit and then re-attest, because `verifyAttestedProcessEventLog` (`:299-311`) checks only that the attestation cites the bytes currently on disk. Deleting one `start` line buys a second resume; appending a fresh entry freshens `lastEvent`. Nothing in the checkpoint, the partial summary, or the runner retains the digest accepted at the previous resume, so the rewrite is invisible to the runner and detectable only by a human comparing committed digests across commits.

**Failure it permits:** a second resume, or a resume outside the 72-hour window, for a custodian willing to forge an artifact whose digest is already committed to git.

**Fix (an enhancement beyond §9, hence P3):** the log is append-only, so any legitimate later state is a byte-prefix extension of the earlier one. Record the accepted `process_event_log_sha256` **and its byte length** into the durable checkpoint at each accepted resume, and at the next resume require that hashing the current log truncated to that length reproduces the recorded digest. That closes both the count and the anchor forgery in one check, using state the runner already writes durably.

### P3-N2 — The nested outcome-free guard is blind to `omitempty` fields and has never seen the production shape

`experiments/frontier-v1/runner/execution_resilience_test.go:20-22` with `execution_resilience.go:37-42`.

`outcomeFreeDigestFields` admits `world_build` and `process_event_log_sha256`. `partialArtifactDigests` also declares `grader_digest,omitempty` and `arm_b_document_digest,omitempty`, both non-empty on the validation path (`main.go:212`, `:445`) and both empty in every fixture. Consequences: `assertExactFieldSet`'s both-directions check would **fail** on a real validation partial summary; and a future `PassCount int \`json:"pass_count,omitempty"\`` left zero in the authoring fixture ships past the guard — the same escape P2-1 reported, one tag away.

**Fix:** drop `omitempty` from those two tags so `artifact_digests` has a fixed four-key shape, admit all four keys in `outcomeFreeDigestFields`, and run the guard once against a fixture whose `runConfig` sets `graderDigest` and `armBDocumentDigest`.

### P3-N3 — `strings.Contains(contents, "ACCEPT")` is satisfied by this repository's own REVISE record

`experiments/frontier-v1/runner/validation_gate.go:123`.

Demonstrated above: `docs/reviews/2026-08-27-execution-resilience-payload-review.md` is a REVISE record that passes every clause of `verifyClosureReviewRecord`. The check also rejects a legitimate `Verdict: Accept` in title case, and matches inside `ACCEPTED`/`ACCEPTABLE`.

**Failure it permits:** the post-closure gate reports that an independently reviewed closure record states ACCEPT when the named record actually returned REVISE.

**Fix:** match a verdict *line* rather than a substring — require a line containing `Verdict: ACCEPT` (case-sensitive on the token, tolerant of markdown emphasis), and refuse if a `Verdict: REVISE` line is also present. Both this repository's review records use that exact `**Verdict: …**` form, so the pattern is already the house convention.

### P3-N4 — The branch that actually closes P2-3 is untested

`experiments/frontier-v1/runner/execution_resilience.go:357-359`, against `execution_resilience_test.go:299-330`.

`TestValidationResumeAnchorsItsWindowToTheProcessEventLogAndTheClock` covers "precedes the last process event", "stamped in the future", and "stale". It does **not** cover `authorized.Sub(lastEvent) > resumeWindow` — the branch that refuses an attestation authored at resume time for a weeks-old interruption, which is the headline attack P2-3 named and the whole reason for the anchor. The branch is reachable and correct by reading, but nothing pins it.

**Fix:** add a case with `lastEvent` at the fixture base, `interrupted_at`/`resume_authorized_at` ~80 hours later, and `now` alongside them, asserting `"after the last process event, past the 72-hour backstop"`.

### P3-N5 — The test that pins the P1-1 fix skips on the recorded Windows host

`experiments/frontier-v1/runner/execution_resilience_test.go:221-247`, skip at `:235`.

Verified by running the suite with `-v`: `--- SKIP: TestProcessTerminationStopsTheRun (0.00s)`. `os.Process.Signal(os.Interrupt)` is unsupported on Windows, so the assertion that a signalled run exits 130 — and the assertion that the termination entry carries live progress, which is the P2-6 evidence — do not execute on the primary recorded host. The skip is silent. Separately, no test on any host drives the watcher through `progress.live.get`, the actual production wiring.

**Failure it permits:** a regression in the stop behaviour or the live-progress wiring passes CI on the Windows host with no signal.

**Fix:** make the signal channel injectable (a `watchProcessTerminationOn(log, progress, exit, signals chan os.Signal)` seam) so the test can deliver a signal without OS support, and assert on a `scheduledProgress` whose `live` handle was advanced by `recordCheckpoint`.

### P3-N6 — The new error-stop entry records a hardcoded zero spend

`experiments/frontier-v1/runner/scheduled_run.go:60`, via `execution_resilience.go:405-415`.

```go
_ = progress.recordStop(summary.Launched, scheduledSummary{Status: "error"})
```

`recordStop` takes `SpentUSD` from the summary it is handed, so the terminal entry written after an evidence-write failure records `spent_usd: 0` regardless of actual spend. This is the same hardcoded-zero class as P2-6, reintroduced narrowly by the P3-9 fix, in the terminal entry of the §9 evidentiary artifact.

**Fix:** `scheduledSummary{Status: "error", SpentUSD: summary.SpentUSD}`, or take the spend from `progress.live.get()`.

### P3-N7 — Dead assignment in the P2-2 test

`experiments/frontier-v1/runner/execution_resilience_test.go:293`. `checkpoint.Resumes = maximumScheduledResumes` is never read; the following assertion uses a fresh `scheduledCheckpoint{Resumes: 2}` literal. staticcheck does not flag field-level dead stores. Harmless, but it reads as if the decoded checkpoint were being reused, which it is not. **Fix:** delete the line.

### P3-N8 — `processEventLog.append` is unsynchronized between the run loop and the signal goroutine

`experiments/frontier-v1/runner/process_events.go:86-116`.

Both goroutines open the same path `O_APPEND` and write independently, with no mutex. This is **not** a Go data race (the struct fields are write-once at construction and never mutated after the watcher starts — see §5-B), so `-race` would not report it; it is a file-level interleaving hazard. In practice both supported hosts make a single small `O_APPEND` write atomic, so the risk is low. If it ever interleaved, the result is fail-closed: `readProcessEventLog` cannot parse the line and the resume is refused — which lands in the "wrongly refuses a valid resume" class.

**Fix:** add `mu sync.Mutex` to `processEventLog` and hold it across `append`. Three lines, and it makes the concurrency story complete without relying on platform append semantics.

### P3-N9 — The 1-hour freshness bound makes the lapse explanation mandatory on every resume, which drains it of signal

`experiments/frontier-v1/runner/execution_resilience.go:337-339` interacting with `:364-366`.

`resume_authorized_at` must be within an hour of `now`; the ordering check forces `conditions_verified_at ≤ resume_authorized_at`; and any strict inequality demands a `lapse_explanation`. So a maximally prompt custodian — conditions verified, attestation written thirty seconds later — must still write a lapse explanation, and the fixture attestation (`execution_resilience_test.go:607`) duly carries one for a two-minute gap. The only escape is to stamp `conditions_verified_at` and `resume_authorized_at` *identically*, which passes (`After` is strict) and destroys the promptness record the two fields exist to capture.

The first review named the broader-than-§9 reading and judged it defensible; I agree it is fail-safe. But the revision turned "broader" into "universal", and a field demanded on every path carries no information. **Fix:** require the lapse explanation only past a promptness threshold (for example when `resume_authorized_at - conditions_verified_at` exceeds an hour), which restores §9's intent that the explanation marks a lapse rather than a formality.

### P3-N10 — An inaccurate claim in the frozen inventory bytes

`experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, `contract.revision.carried[0]`, echoed at `candidate.md:76`: *"no non-test function in the payload sits within one of the gocyclo gate."*

`runScheduledCases` is inventoried, is modified by this payload, and measures **14** — one below the `-over 15` gate. I measured it at both `HEAD~1` and `HEAD` (14 and 14), so the revision did not raise it; the claim is nonetheless false as written, and it is about to be frozen.

**Fix:** reword to the accurate statement — *"`verifyResumeAttestationFields` is split and now measures below 5; `runScheduledCases` and `validateExternalGrade` remain at 14, both pre-existing and unchanged by this payload."*

## 5. The whole payload, as if fresh

**A. Section 9 conditions, one by one.** (a) no inspection — `:275-277`, a boolean claim, refused when false; correct scope, §9 words it as a written attestation. (b) digest identity — `:275-277` plus `verifyResumeFrozenDigests` `:284-297`, now refusing rather than skipping on an empty live digest. (c) 72 hours — `:333-336` plus the anchors `:348-368`; **now genuinely checked**, subject to P2-N1, P3-N1, and the host-clock residual. (d) pre-decision classification on the frozen list — `:271-282`; **OOM is refused**, as are the empty string and any unlisted cause, with an accurate message; `process_kill` additionally requires the written establishment at `:278-280`. Timestamp ordering — `:327-332`, all four in sequence. Promptness and lapse — `:262-264` and `:337-339`, see P3-N9. One resume — `:165` before the attestation is read, corroborated at `:186-201`. Abandonment by choice — no code path represents it; every failure returns an error before any launch.

I re-ran the first review's negative results and they still hold: a validation schedule cannot be run under `--tranche authoring` (`validateScheduleForTranche`); a stale log digest is refused; an unlisted or empty cause is refused; `decodeStrict` rejects unknown fields and trailing JSON; a foreign `output_dir` is refused; and no model is invoked before the refusal — `requireValidationGate` including the post-closure audit is the first thing `prepareScheduledCLI` does (`main.go:183-185`), ahead of the world build, the grader, and the output directory.

**Does any new refusal path fire on a legitimate first run?** No. `requireResumeAuthorization` returns `nil` at `:162` unless `resume.hasCheckpoint`, so the log read, the count check, and the anchors are unreachable on a fresh run. The one adjacent case — a directory carrying a pre-v5 checkpoint but no process-event log — refuses with `read process-event log`, which is correct: such a tranche is a v4 artifact, retired by §10, and the post-closure gate blocks a new execution anyway.

**B. Race audit (substituting for the `-race` run I could not perform).** `liveProgress` (`process_events.go:194-209`) guards both `set` and `get` with the same mutex. I grepped every reference: the field `value` is never touched outside those two methods; the `*liveProgress` is allocated once at `scheduled_run.go:90` and never reassigned; `scheduledProgress`'s value receiver is harmless because `live` is a pointer. The two writers are `scheduled_run.go:103` (before the watcher exists, so unordered access is impossible) and `execution_resilience.go:471` (run loop); the single reader is the signal goroutine via the `progress.live.get` method value. **The synchronization is correct and complete for memory access, and I expect `-race` to be clean.** The other cross-goroutine state is `*processEventLog`, whose fields are written once at construction (and, in tests, `clock` before any goroutine starts) and only read thereafter — safe. The remaining hazard is file-level, not memory-level, and `-race` would not see it: **P3-N8**.

**Does the reordering inside `TestPreValidationFreezeMatchesAcceptedCandidates` change what it guarantees?** No. `verifyAcceptedLocalArtifacts` is moved *above* the deliberate `t.Fatalf` in `verifyReplacementRunnerFreeze`, so it now runs at payload commits and previously did not. Nothing is skipped that used to run: the checks after the runner-freeze call (`verifyFrozenFileSet`, `verifyFrozenValidationSchedule`, `verifyClassifiedSetFreeze`) were already unreachable at a payload commit and remain so. The test's guarantee at a *refreeze* commit, where nothing fatals, is unchanged — the same set of checks in a different order.

**C. Partial summary outcome-free in both directions.** Field set unchanged from the first review: `v`, `tranche`, `status` (always `in_progress`), `schedule_digest`, `artifact_digests`, `started_at`, `updated_at`, `next_launch_index`, `completed_pairing_keys`, `arm_assignment_counts`, `spent_usd`, `resumes`. No success count, grade tally, or per-arm outcome field. `arm_assignment_counts` counts *scheduled* assignments below `nextLaunchIndex` and is fully determined by the launch index given the pairing-complete schedule. Spend is retained on §9's own stated grounds. The guard is now two levels deep but incomplete — P3-N2.

**D. Process-event log honest about what the host does not expose.** Unchanged and still correct: `recordTermination` sets `exposed: false`, `source: signal_without_sender_identity`, and a reason, rather than presenting the receiving process's owner as the killer. I re-verified the technical claim: Go's `os/signal` delivers only the signal number on both supported hosts, so `si_pid` is not reachable through `signal.Notify`. `TestProcessTerminationRecordsAnUnexposedSendingPrincipal` runs on this host and passes.

**E. Post-closure gate encoded in the frozen gate rules and applying to decision 0023.** Yes, in all three respects. Code: `requirePostClosureAudit` (`validation_gate.go:91-105`), called from `requireValidationGate:71`. Frozen document: `validation_execution_closure_gate` / `_review` / `_review_verdict`, now correctly indented. Frozen test: `verifyPostClosureGateState` (`artifact_freeze_test.go:213-231`) asserts four substrings of the rule text and pins both closure fields **empty**, while `validation_execution_status` stays `"indeterminate"` — so the gate refuses today on the decision-0023 ground and cannot drift. The predicate `status != "" && status != "complete"` is default-deny, covering budget stop and lapse without enumeration. `protocol.json` carries the new blocker; the list is 8 entries and `corpusctl`'s `protocol_test.go` asserts the count, the tail entry, and the amendment substrings, so a regression fails that suite.

I attacked the gate one further way the first review did not: **wipe the output directory and start over.** With no checkpoint, no assignments, and no progress records, `finishScheduledInventory` sees a fresh run, `requireResumeAuthorization` returns at `:162`, and the resume rule never engages. The only thing standing between a custodian and an unlimited retry channel is `validation_execution_status` being set at closure and the post-closure gate refusing until a reviewed closure record exists. That makes the chair-discipline residual load-bearing rather than incidental, and I have sharpened its wording below.

**F. Records match the implementation.** `protocol.json`'s `amendment.execution_resilience` tracks §9 closely with no overclaim. The blocker edit is correct: blocker 4 now reads "path-witness authoring artifacts and the witness-derived uniform turn cap frozen with the schedule", the §5 countersignature clause having been discharged by decision 0025. The README's two paragraphs are accurate as far as they go. The candidate note's `§1` claim about every entry carrying progress is now true, P2-6 having been fixed. **The overclaims I found are P2-N2** (the boundary document omits two live refusal conditions) **and P3-N10** (the gocyclo claim in the inventory).

**G. Freeze hygiene.** Exactly one deliberate red, and it is the one named. World-build and grader identities are provably unmoved — neither digest string appears in the diff, and no file under `cmd/`, `internal/`, `verbs/`, `build/`, `worlds/`, `schema/`, `go.mod`, or `go.sum` changed. Editing the `gates` block at a payload commit remains acceptable on the same reasoning the first review gave: three added keys, both closure fields empty, no digest or gate flag moved, and the freeze test needs the fields present to run at all.

## 6. Judgment on the residuals offered in area D

- **P3-1 (committed log digest never covers the terminal entry).** Correctly characterized and acceptable. Note the interaction with P3-N6: the terminal entry is both outside the runner-committed digest *and* currently carries a zero spend on the error path. Fix P3-N6; keep P3-1 as a residual.
- **P3-2 (five identity fields beyond §9's literal list).** Correctly characterized and acceptable. I re-checked each: `v`, `tranche`, `status`, `schedule_digest`, `resumes` — none carries an outcome, and identity fields are needed to bind the record to a tranche. **Concur.**
- **P3-5 (checkpoint and partial summary not fsynced or atomically renamed).** Correctly characterized as failing closed, and I verified that claim independently: a truncated or zero-byte checkpoint fails `decodeStrict`, and a *deleted* one is caught by `finishScheduledInventory`'s progress-implies-checkpoint rule. **But the characterization now understates the exposure.** By deliberately calling `os.Exit` from the signal goroutine, the revision makes a mid-write kill a routine consequence of a *supervised stop*, not only of an abrupt external kill — so a normal Ctrl+C at an unlucky microsecond closes the tranche indeterminate. Still acceptable as a residual (§9 accepts indeterminate closure, and the window is microseconds), but the wording must say so, and the chair should know that a temp-and-rename in `writeJSON` would remove the class for a few lines. I have reworded it below.
- **P3-6 (blocker edit also strikes the §5 clause).** Correctly characterized. I verified blocker 4's new text and that decision 0025 discharges the countersignature clause. It is named in the candidate note, and the obligation to name it in decision 0026 carries forward. **Concur.**
- **The authoring-tranche exemption.** I concur with the first review's judgment and re-verified the decisive argument: `validateScheduleForTranche` refuses a schedule whose tranche does not match the flag, so the exemption cannot launder a validation execution. One consequence worth adding to the record: because of the exemption, the P2-2 round-trip test runs on the authoring path and so never traverses `requireResumeAuthorization`. That is a coverage consequence, not a defect. **Concur, with that sentence added.**
- **The chair-discipline caveat on `validation_execution_status`.** Correctly characterized and genuinely acceptable — nothing in the runner can write that field, and the freeze workflow is the control. **But it is more load-bearing than its current wording suggests**, because the wipe-and-restart path (§5-E) bypasses the resume rule entirely and leaves the post-closure gate as the sole control. **Concur, with sharpened wording below.**

## 7. Residual list — **not operative at this verdict**

The verdict is REVISE, so there is no refreeze to attach findings to. For the refreeze that follows the two P2 fixes, these are the entries I would accept into `accepted_findings`, worded for direct transcription:

- "Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian's written attestation, as section 9 requires."
- "Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b)."
- "Residual: the process-event log is unsigned custodian-controlled local state, so a coordinated rewrite that deletes a start entry and re-attests the new digest defeats both the single-resume corroboration and the 72-hour anchor; the runner retains no previously accepted digest and cannot detect it. Detection is the committed attestation digest compared across commits, under independent review."
- "Residual: the 72-hour anchor is measured against the runner's own clock, which the custodian's host controls; a clock set to shortly after the last process event passes every timing check. A clock set before the last process event is refused by the interruption anchor. No trusted time source exists on the execution host."
- "Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports 'gate is closed' first; both are refusals and only the surfaced message differs."
- "Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag. One consequence: the end-to-end resume-counter round-trip test runs on the authoring path and so does not traverse requireResumeAuthorization."
- "Residual: nothing in the runner writes validation_execution_status. A custodian who wipes the output directory presents as a fresh run and bypasses the resume rule entirely, so the post-closure gate — armed by the chair setting that field at closure, under the freeze workflow — is the sole control on the retry channel after a future non-completed execution. It is armed by frozen bytes for the current one."
- "Residual (P3): the partial summary's committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure."
- "Residual (P3): the partial summary carries five non-outcome identity fields beyond section 9's literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome."
- "Residual (P3): the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log. Because the termination handler now deliberately calls os.Exit, a supervised stop as well as an abrupt kill can truncate a boundary write; the failure is closed — decodeStrict refuses and the tranche closes indeterminate — but a routine Ctrl+C at an unlucky microsecond can burn the tranche. A temp-and-rename in writeJSON would remove the class."
- "Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first."
- "Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required."
- "Carried P2-1 is retired by this payload: classifyScheduledOutputEntries is decomposed into scheduledOutputEntryKind, recordScheduledAssignment, and finishScheduledInventory, and verifyScheduledCheckpoint is split into identity and progress halves, so both are back under the gocyclo threshold with headroom. Its successor P2-5 is also retired: verifyResumeAttestationFields is split into identity and condition halves and measures below 5. runScheduledCases and validateExternalGrade remain at 14, both pre-existing."

P2-N1 and P2-N2 are **not** in that list: both are cheap, both concern frozen bytes that cannot be repaired at refreeze under this workflow, and both refuse a resume the specification intends to permit. P3-N1 through P3-N10 I would accept as residuals if the chair prefers, but P3-N2, P3-N3, P3-N4 and P3-N6 are each a few lines and I would rather see them fixed alongside the P2s.

## 8. What I did not verify

- I did not run a model, obtain a private grade, or observe any validation outcome. Both gates were false at the start and end of this review.
- **I could not execute `go test -race`** on this host, for the toolchain reason detailed in §2. My concurrency conclusion in §5-B rests on reading every reference to the shared state, not on the detector. This should be re-run on the WSL2/Linux host before the refreeze.
- I did not execute the signal path against a live scheduled run. `TestProcessTerminationStopsTheRun` skips on this host (P3-N5), so my P1-1 conclusion rests on reading `watchProcessTerminationWith`, confirming its single production call site passes `os.Exit`, and tracing the record-then-exit ordering.
- I did not empirically compile a mutated `partialArtifactDigests` to demonstrate P3-N2; the finding rests on `assertExactFieldSet`'s both-directions membership test against an admitted set that omits two `omitempty` fields the validation path populates, which is unambiguous from the source.
- I did not attempt to demonstrate P2-N1 with a live signalled run; it rests on the absence of any tolerance at `execution_resilience.go:353`, on `time.Parse(time.RFC3339, …)` versus `RFC3339Nano` stamping, and on the ordering of the signal and the log write in `watchProcessTerminationWith`.
- I made no change of any kind to the repository. My only writes were to the session scratchpad (one extracted copy of the parent commit's `scheduled_run.go`, for the gocyclo before/after measurement) and to this record.

---

## Chair transcription note

Transcribed verbatim from the second independent reviewer's record. Before transcription the payload author independently re-verified every digest in §1: all fourteen match. P2-N1, P2-N2, P3-N2 and P3-N3 were independently reproduced by reading the cited lines — P3-N3 was confirmed by counting five occurrences of the `ACCEPT` token in the committed REVISE record, which does satisfy `verifyClosureReviewRecord` as written. The verdict is **REVISE**, so this payload does not proceed to refreeze; a third payload revision follows and is submitted for a further independent review. The `-race` gap is the reviewer's own, honestly reported, and is discharged separately on the linux host before any refreeze.
