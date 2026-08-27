# Protocol-v5 section 9 execution-resilience candidate

Status: **review candidate, not frozen.** Both outcome gates stay closed (`may_open_validation: false`, `may_open_held_out: false`). No model was run, no private grade obtained, and no validation outcome observed while authoring this payload.

Base commit: `fec6d554156eb749057fb8b6cd3db927c2fecf05` (the fourth review record). Frozen state being replaced: decision 0025, `2cf99331688f4705ed7d6cd3efed98941f59defc`.
Inventory: [`gate-1a-execution-resilience-candidate.json`](gate-1a-execution-resilience-candidate.json), 18 files.

**Revision 5.** Four independent reviews have run, by four different reviewers, and all four returned REVISE.

- Review 1 ([record](../../../docs/reviews/2026-08-27-execution-resilience-payload-review.md)) against payload `7a1ea42513acd3d55e276eabc2459da4a037acc0`: two P1, six P2, nine P3, no P0.
- Review 2 ([record](../../../docs/reviews/2026-08-27-execution-resilience-payload-review-2.md)) against revision `e92eaacf6ab54b547b331d14749717da64b58e0d`: **no P0 and no P1**; it confirmed both P1s genuinely fixed and all six P2s fixed or substantively fixed, and raised two new P2s and ten P3s.

Every P1 and P2 from all four reviews is fixed here, along with most of the P3 notes. The two second-review P2s were both *in the fix code*, and both refused a **legitimate** resume rather than permitting an invalid one — see "What the second review changed" below. Nothing in the reviewed design has been discarded across either round; every finding has been a defect in the enforcement, not in the contract.

## What this implements

Proposal section 9 (`docs/plans/2026-08-27-protocol-v5-proposal.md`) and the fourth review's freeze obligation 5: the process-event log created per section 9 with its digest committed alongside the custody record, and the post-closure audit-before-retry gate encoded in the frozen gate rules.

### 1. Process-event log — a new artifact

`experiments/frontier-v1/runner/process_events.go` writes `process-events.jsonl` into the scheduled output directory. It is append-only, each line written with `O_APPEND` and fsynced, so an abrupt termination cannot lose the entries already recorded.

Event kinds: `start` (once per custodian process, after resume authorization and before any launch), `checkpoint` (every pairing-key boundary), `stop` (the scheduled run's terminal state), and `termination` (SIGINT/SIGTERM, from a handler installed for the duration of the scheduled run).

Each entry carries `v`, the event kind, an RFC3339Nano UTC timestamp, tranche, schedule digest, world build, the launch index / completed pairing keys / cumulative spend at that moment, the resume count, and the originating principal.

**Principal, "where the host exposes it."** For events the custodian process itself emits, the host exposes the process owner: `exposed: true`, `source: host_process_owner`, with user, hostname, and PID. For a signal-driven termination neither Windows nor Linux exposes the *sending* principal to a Go signal handler, so the entry records `exposed: false`, `source: signal_without_sender_identity`, and a reason saying so, rather than presenting the receiving process's own owner as the killer. This is precisely why section 9 makes `process_kill` unsafe absent a written establishment: the log can prove *when* and *that* a kill arrived, never *who* sent it, so the custodian must establish the principal by other means.

### 2. Outcome-free partial summary

`scheduled-summary.partial.json` is rewritten at every pairing-key boundary (`execution_resilience.go`). Its field set is exactly the one section 9 admits — launch index, completed pairing-key list, per-arm assignment counts, artifact digests (world build, grader, Arm B document, process-event log), timestamps (`started_at` preserved across rewrites, `updated_at` stamped), and cumulative spend — plus the identity fields naming the tranche, schema version, in-progress status, schedule digest, and resume count, none of which carries an outcome. `TestScheduledRunWritesAnOutcomeFreePartialSummary` asserts the JSON object's key set in **both** directions against `outcomeFreePartialFields`, so a later field addition that leaks a success count or grade tally fails the suite rather than shipping.

Spend is retained because the budget-stop rule requires it; the mandatory-resume rule removes the channel through which spend could otherwise steer a resume choice. `scheduled-summary.json` remains the completion criterion, unchanged.

### 3. Mandatory resume, precommitted

`requireResumeAuthorization` runs in `openScheduledProgress`, before the first launch of a resumed run. It fires when a validation output directory already carries a checkpoint written by an earlier custodian process, and refuses unless `--resume-attestation` names a custodian record establishing all of:

| Condition | Field(s) | Refusal |
| --- | --- | --- |
| (a) no outcome inspected between interruption and resume | `no_outcome_inspection_between_interruption_and_resume` | refuses |
| (b) frozen bytes digest-identical, recorded in the attestation | `frozen_bytes_digest_identical_at_resume`, `frozen_digests_verified.{schedule_digest,world_build}` compared against the live values | refuses |
| (c) resume within 72 hours | timestamps required in order; `interrupted_at` at or after the log's last entry within five minutes' skew; 72 hours measured from that entry; authorization stamped within an hour of invocation | refuses |
| (d) cause classified before inspection and before the resume decision, and on the frozen list | `cause_classified_before_inspection_and_resume_decision`, `cause` ∈ {`host_restart`, `power_loss`, `hardware_failure`, `process_kill`} | refuses |
| kill establishment | `kill_principal_established_as_non_participant` plus `process_event_log_sha256` re-verified against the log on disk | refuses |
| at most one resume | checkpoint `resumes` field corroborated against the `start` entries in the process-event log; `resume_index` must equal `resumes + 1` | refuses |
| promptness | `promptness_statement` required; a `resume_authorized_at` later than `conditions_verified_at` additionally requires `lapse_explanation` | refuses |

**OOM is struck** — it is simply absent from `frozenResumeCauses`, so it falls to the generic "not on the frozen outcome-uncorrelated cause list" refusal along with any other unlisted cause, including the empty string. A refusal is not an abandonment: the run stops without launching, and the tranche closes indeterminate under the existing stop rules. Abandonment by choice has no representation anywhere in the runner.

The single-resume rule is carried in the durable checkpoint, which gains a `resumes` counter: a fresh run writes 0, the one permitted resume writes 1, and a checkpoint already at 1 is refused before any attestation is read (`already used its single resume`). `verifyCheckpointProgress` bounds the counter, so a hand-edited checkpoint claiming a *higher* count is a mismatch, not a licence. A counter edited *downward* is caught separately: `verifyResumeCount` counts the `start` entries in the append-only process-event log — one per custodian process that has opened the directory — and refuses when the log attests more resumes than the checkpoint records. A missing or unreadable log refuses rather than defaulting to zero, so the attack fails closed. Each resume's attestation also commits the log's digest, which makes a rewritten log detectable across commits in a way an edited integer is not.

The 72-hour window is likewise anchored to something the custodian did not author. `interrupted_at` may not precede the last timestamp the interrupted process wrote to its own log; `resume_authorized_at` must fall within the window measured from *that* timestamp, not only from the attested interruption; and the authorization must be neither stamped in the future nor stale against the runner's clock. Back-dating an interruption to manufacture a fresh window is therefore refused.

**Scope.** The attestation rule applies to `--tranche validation`. Authoring resumes are development loops and are exempt (`TestAuthoringResumeIsExemptFromTheAttestationRule` pins the exemption so it is a decision, not an oversight). Reviewer: this scoping is the one place the payload chose narrower than the proposal's unqualified wording, on the grounds that section 9's subject throughout is the validation tranche and its indeterminate closure.

### 4. Post-closure gate

`requirePostClosureAudit` in `validation_gate.go`: when the freeze document's `validation_execution_status` is anything other than empty or `complete`, `requireValidationGate` refuses unless `validation_execution_closure_review_verdict` is `ACCEPT` and `validation_execution_closure_review` names a non-empty regular markdown file under `docs/reviews/` whose own text states an ACCEPT verdict. Existence alone is not enough: a directory, a file elsewhere in the repository, or a record whose verdict contradicts the gate document all refuse. The gate document gains those two fields plus `validation_execution_closure_gate` stating the rule, and `verifyPostClosureGateState` in the freeze test asserts the rule text and that both fields stay **empty**.

**Consequence the chair should note.** The 1,475-launch execution closed by decision 0023 ended without a completed schedule. The gate therefore applies to it retroactively, exactly as section 9 specifies ("after any validation execution that ends without a completed schedule"). A new blocker is registered in `protocol.json`: an independently reviewed and committed closure record for that execution is now required before any later validation execution, on top of the existing blockers. This payload does not produce that review — it is a separate artifact, an independent audit of decision 0023's cause classification, custodian logs, and attestations.

### 5. Protocol and boundary records

`protocol.json` gains `amendment.execution_resilience` carrying the whole rule, and its blocker list is rewritten: the section 9 clause is struck from the remaining-v5-amendment blocker (leaving path witnesses and the witness-derived cap), and the closure-record blocker is inserted, taking the list from 7 to 8 entries. `validation-execution-boundary.md` gains an "Execution resilience, mandatory resume, and the post-closure gate" section carrying the host requirements (automatic OS restarts disabled, no interactive login sessions, supervised custodian process), the three per-checkpoint records, the five numbered resume conditions, and the post-closure gate including its retroactive application to the decision-0023 closure. `README.md` describes the two new output records and the resume flag.

## What the first review changed

The reviewer's two P1 findings were both real, both mine, and both in code this payload added. I reproduced each before fixing it.

**P1-1 — the termination handler swallowed the termination.** `signal.Notify` disables Go's default action, and the original watcher recorded an event and returned without exiting or cancelling. A scheduled run would have become un-interruptible for its whole multi-hour duration while the log asserted a termination that never happened; SIGKILL, the only signal that still stopped it, left no entry at all — inverting the artifact's evidentiary purpose for the `process_kill` cause. The handler now records the entry durably and then exits with 128 plus the signal number. `watchProcessTerminationWith` takes an injectable exit so `TestProcessTerminationStopsTheRun` asserts the process actually stops, rather than only that a line was written.

**P1-2 — the single-resume bound was launderable downward.** The counter check ordering was right and an inflated counter was refused, but nothing refused a *lowered* one: editing `"resumes": 1` to `0` in the checkpoint bought a second resume with every other check green. The reviewer's decisive point was that this payload already writes the corroborating evidence and never reads it — `process-events.jsonl` carries exactly one `start` entry per custodian process, and the runner hashed that file without parsing it. `verifyResumeCount` now parses it and refuses when the log attests more resumes than the checkpoint records. A deleted or unreadable log refuses too, so the attack fails closed rather than defaulting to zero.

The P2 findings are fixed rather than carried, because each weakened a section 9 condition:

- **P2-1** the outcome-free field-set guard was depth-1, so a field added inside `artifact_digests` would have shipped past it. It now descends into that object against its own admitted key set and requires `arm_assignment_counts` keys to be Phase 1 arms.
- **P2-2** nothing exercised the counter round trip; deleting the increment broke no test. `TestScheduledResumePersistsTheResumeCounter` now runs an interrupted schedule to completion and asserts the on-disk checkpoint and partial summary both read `resumes: 1`.
- **P2-3** the 72-hour window was arithmetic on four numbers the custodian typed. It is now anchored to evidence the custodian did not author: `interrupted_at` may not precede the last process-event timestamp, `resume_authorized_at` must be within the window *of that timestamp*, and the authorization must be neither future-stamped nor stale against the runner's clock. The clock is injected through the control struct so the anchors are testable.
- **P2-4** the closure record was an `os.Stat`, which any directory or unrelated file satisfied. It must now be a non-empty regular markdown file under `docs/reviews/` that states an ACCEPT verdict.
- **P2-5** `verifyResumeAttestationFields` had landed at gocyclo 14, one below the gate — re-creating the near-threshold condition that carried finding P2-1 recorded. Split into identity and condition halves; no non-test function in this payload now sits within one of the gate.
- **P2-6** every production termination entry recorded a hardcoded zero progress, which would have misled an auditor precisely on the second interruption of a resumed tranche. Live progress now reaches the handler through a mutex-guarded `liveProgress` handle updated at each checkpoint, and the watcher moved from `main.go` into `runScheduledCases` where that state exists.

Also fixed: **P3-3** (the stop event recorded the raw failure text, which for a grade-shape failure could quote a per-trial status; it now records a fixed stop-kind token), **P3-4** (the added gate fields were indented inconsistently and would have been frozen that way), **P3-7** (a missing live digest skipped its comparison instead of refusing), **P3-8** (`verifyAcceptedLocalArtifacts` was unreachable at a payload commit behind the deliberate red; hoisted above it), and **P3-9** (an evidence-write failure left no terminal entry, indistinguishable from an abrupt kill; it now records a stop).

Not fixed, and offered as residuals: P3-1 (the committed log digest never covers the terminal entry; the custody record hashes the whole file at closure), P3-2 (five non-outcome identity fields beyond §9's literal list), P3-5 (the checkpoint and partial summary are not fsynced, failing closed to indeterminate), P3-6 (the blocker edit also strikes the discharged §5 countersignature clause — named here and to be named in the decision document). The reviewer's judgment on the authoring-tranche exemption — defensible, but to be recorded knowingly rather than passed silently — is accepted, and it is carried as a residual with the reviewer's own reasoning attached.

## What the second review changed

The second reviewer found no P0 and no P1, and confirmed by independent attack that the first round's fixes hold: the resume-counter corroboration survives log deletion, truncation, zero-start logs, semantically empty lines, and injected starts; the termination handler's exit is reachable in production and its record precedes the stop; the post-closure gate fires on every non-completed ending and cannot drift. Its two new P2s were both defects the first payload could not have had, because they live in the code written to fix it.

**P2-N1 — the interruption anchor had no skew allowance, and the P1-1 fix is what made that bite.** Now that a catchable signal writes a `termination` entry, the log's last `recorded_at` is stamped at `T + ε` in nanoseconds, while a custodian attesting the interruption instant writes `T` at second precision. The anchor compared them with `Before` and no tolerance, so for the two frozen causes that arrive as signals — `host_restart` via systemd SIGTERM on the prepared Linux host, and `process_kill` — an *honest* attestation was refused, with a message that did not say what to write instead. Under a mandatory-resume rule the runner is the enforcement point, so that refusal burns clock against the 72-hour backstop and can close a ~300 USD tranche indeterminate. The anchor now allows the same five minutes of skew its counterpart already carried, and two tests pin both sides: a stamp two seconds early is accepted, one a day early is refused.

**P2-N2 — the frozen operational record described the previous payload.** Condition 3 of the boundary document still read as the first payload's timing rule, so it omitted both obligations the revision added. Worse, one of them — the one-hour authorization freshness bound — has **no basis in section 9 at all**; it is an implementation control. Both are now stated, and the one-hour bound is explicitly labelled as a control giving the promptness duty effect rather than as a section 9 requirement, so a custodian reading the frozen contract can satisfy the runner and can tell which obligations come from the protocol.

The P3s the reviewer preferred fixed are fixed:

- **P3-N3** the closure check matched the bare token `ACCEPT` anywhere in the file — which **this repository's own committed REVISE records satisfy**, since every reviewer is asked to return "ACCEPT or REVISE". It now reads the record's verdict *statement*. The third review showed the first attempt at this was wrong in both directions against the project's own corpus, and it is rewritten here — see "What the third review changed".
- **P3-N2** `partialArtifactDigests` carried two `omitempty` fields that are empty in every fixture and non-empty on the validation path, so the nested guard had never seen the shape it guards and would have *failed* against a real validation partial summary. The struct now has a fixed four-key shape, all four keys are admitted, and the fixture sets both digests.
- **P3-N4** the branch that actually closes P2-3 — an attestation authored at resume time for a weeks-old interruption — was unexercised; it now has a case.
- **P3-N5** the test pinning the P1-1 fix silently skipped on Windows, so the evidence for the fix ran on no host. `watchProcessTerminationOn` takes the signal channel directly, so the stop, the exit code, and the live-progress wiring are exercised everywhere; SIGTERM's 143 is pinned too.
- **P3-N6** the error-path stop event reintroduced a hardcoded zero spend; it now carries the real figure.
- **P3-N8** `append` is mutex-serialized, so the run loop and the signal goroutine cannot interleave a line.
- **P3-N9** the lapse explanation was demanded on every resume, which drained it of signal; it is now required only past a one-hour promptness window.
- **P3-N7** and **P3-N10** (a dead assignment, and an inaccurate gocyclo claim about to be frozen) are corrected.

Not fixed, and carried as residuals with the reviewer's own wording: **P3-N1** (a coordinated rewrite of the process-event log defeats both the counter and the anchor, since the runner retains no previously accepted digest — the reviewer graded this P3 as an enhancement beyond section 9, and detection remains the committed attestation digest compared across commits), the host-clock residual, and P3-1/P3-2/P3-5 from the first round with the second reviewer's sharpened wording on P3-5.

## What the third review changed

The third reviewer confirmed that **every P1 and P2 from both earlier rounds is genuinely fixed and none has regressed**, walked an honest custodian's timeline through all four frozen interruption causes and found each passes, verified that the mutex added for P3-N8 does not disturb the timing anchor (timestamps are stamped before the lock is taken, and the reader takes the maximum rather than the last line), and confirmed that a *refused* resume writes no `start` event and so does not consume the tranche's single allowance. It also ran the race detector on the linux host itself rather than restating this note's claim.

It found one P1, and it is the sharpest finding of the three rounds because it is the round-2 pattern repeating at one further remove.

**P1-N1 — the closure verdict parser was wrong in both directions.** Review 2 found that matching the bare token `ACCEPT` anywhere in a file accepts a REVISE record, and prescribed matching a verdict *line*. That prescription rested on a survey of two files. The parser written to it, tested against the two files, and about to be frozen, turned out to misclassify most of this repository's corpus. The third reviewer reimplemented it standalone and ran it over all 76 committed review records:

- It **refused 17 of the 19** records written in this project's dominant style, a bolded `Verdict:` label followed by the token in a code span. The trim cutset omitted the backtick, so the token never matched. It also refused the current house style, where a `## 1. Verdict` heading carries `**ACCEPT.**` on the next line, and refused any accepting record that cites an earlier round's REVISE.
- It **accepted 12 documents that are not accepting review records at all**, five of them evaluator prompts matching on the instruction `1. Verdict: ACCEPT, REVISE, or REJECT.` — and the review-3 prompt itself, on the line where it warns about this hole.

Either direction defeats the mechanism freeze obligation 5 exists to install: the false accept lets `validation_execution_closure_review` naming a `…-prompt.md` file satisfy the post-closure gate with no review performed, and the false refusal blocks the eventual decision-0023 closure review at the point of maximum schedule pressure, repairable only by another full freeze cycle.

I reproduced both counts before fixing: 19 records use the backtick style, and the prompt boilerplate is present in five prompts.

The parser is rewritten to be structural rather than substring-positional. A verdict statement is now a line whose text before the label is markdown furniture only (emphasis, code spans, quotes, list and heading markers, section numbering) optionally preceded by a closed set of qualifiers (`final`, `overall`); the token after the label is read as a lone word with non-letters trimmed, so `` `ACCEPT` ``, `"ACCEPT"`, `**ACCEPT**`, `ACCEPT.` and `Accept` all resolve; a label with no token takes the next non-blank line, which is the heading style; a verdict *list* is rejected by requiring that the other verdict word not appear in the token's immediate neighbourhood; and only the record's **first** verdict statement is decisive, so citing an earlier round does not overturn it.

**The corpus is now the arbiter, in the code.** `TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus` runs the parser over every committed review record and fails if any evaluator prompt reads as an accepting record or any record stating ACCEPT is refused — the same check the reviewer performed by hand, now permanent. `TestClosureVerdictParserReadsEveryCommittedVerdictStyle` pins all nine committed verdict forms plus five that must be refused. The pinned corpus records 23 ACCEPT, 15 REVISE, 1 REJECT and 38 documents stating no verdict of their own.

Also fixed from this round: the boundary document now enumerates the attestation's exact JSON field set, since `decodeStrict` refuses unknown or misspelled fields and the field names appeared only in this note and the Go source (P3-N1 new-b); it records the one unresumable state, a process that died between its `start` event and its first checkpoint (P3-N1 new-c); and the inventory's stale round-count text is corrected (P3-N1 new-d). The `os.Exit`-versus-`writeJSON` truncation window (P3-N1 new-a) is carried as the sharpened P3-5 residual, as both the second and third reviewers graded it.

## What the fourth review changed

The fourth reviewer confirmed the round-3 P1 genuinely fixed and proved it independently — extracting the parser into a standalone module and measuring **zero misclassifications across all 75 committed records in both directions**, where round 3 measured 17 false refusals and 12 false accepts. It confirmed every P1 and P2 from all three earlier rounds still fixed with none regressed, ran the race detector itself, and verified **by execution** that an attestation built from the boundary document's prose alone is accepted. Its two P2s were, once again, in the previous round's fix code.

**N1 — the corpus test guarded one direction while claiming two.** `verifyClosureReviewVerdictLine` returns nil exactly when `closureReviewVerdict` yields ACCEPT, so inside the test's error branch the false-refusal guard was unreachable code. The test could only ever catch a prompt being accepted — never a genuine ACCEPT record being refused, which is the direction round 3 measured at 17-of-19 wrong. Three records asserted it arbitrated both, one of them frozen.

Ground truth is now a **committed expectations table**, `artifacts/closure-verdict-corpus-expectations.json`, labelling every review record with the verdict it states, derived by a scan deliberately unlike the runner's. The test compares the parser against those labels in both directions and pins the per-verdict counts, so a regression in either direction fails.

**N2 — the verdict-list rejection was fitted to word order.** The three-field window caught `ACCEPT, REVISE, or REJECT` only because every prompt in this repository happens to put `REVISE` second. Reordered to `ACCEPT, REJECT, or REVISE` — equally natural English — the window misses and an evaluator prompt reads as an accepting record: exactly the false accept the check exists to prevent.

The replacement is order-independent: a verdict token followed by a comma or a slash, with another verdict word **anywhere** later on the line, is an enumeration. **I did not adopt the reviewer's literal prescription**, which was to reject on any other verdict word anywhere on the line. Round 3's own record disproves it: its verdict line reads `**Verdict: REVISE** — … it refuses 17 committed ACCEPT records`, and the literal rule would have refused to read it. The reviewer verified that prescription against accepting records only; taking a rule verified on a narrow slice is precisely what produced the round-3 P1, so the comma-or-slash refinement is verified against the whole corpus in both directions instead.

**N3 — no block context.** A verdict inside a fenced block, a block quote or an indented block read as the record's own, and the block-quote marker was itself in the furniture cutset. The parser now tracks fences and skips quoted and indented lines, and requires a separator between the label and the token, so `Verdict ACCEPT was withheld…` is prose again rather than a statement.

**N4 — a fitted cutset refusing plausible records.** Em and en dashes, footnote brackets, carriage returns and non-breaking spaces are now furniture. The qualifier set stays closed at `final` and `overall`, and instead of widening it speculatively the boundary document now **states the accepted verdict-statement forms**, so the closure record's author writes one the frozen parser reads.

**N5 — a frozen test coupled to a live directory.** The corpus is pinned, not swept, so a review record added by unrelated later work can no longer turn a frozen test red. **N6** — the note's corpus count is corrected.

## Deliberate freeze red

Exactly one test is expected to fail at this commit, and must:

- `TestPreValidationFreezeMatchesAcceptedCandidates` — the frozen `scheduled_runner` set digest and per-file digests still pin the decision-0025 bytes.

`TestLocalArtifactCandidateMatchesImplementation` and `TestValidationBuildReproducesFrozenWorldBuildDigest` stay **green**: the Phoenix binary is untouched (nothing under `cmd/phoenix` or its import graph changed), so the world-build identity remains `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4`, and the grader digest is unchanged at `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b` because no corpusctl non-test source, `go.mod`/`go.sum`, or `schema/*.schema.json` file was touched. No authoring label, manifest, or grading_script pin moves in this payload.

Everything else is green: root suite, corpusctl suite, surface-spike, `go vet`, staticcheck, `gocyclo -over 15`, `dupl -t 100`, `validate-spec`, `validate-authoring`.

**The race-detector gap is discharged.** The second reviewer could not run `go test -race` on the Windows host (no GNU-style C toolchain; the detector needs cgo and external linking) and reported that honestly rather than claiming it green. It has since been run on the project's linux/amd64 execution host — the same WSL2 Ubuntu with Go 1.26.6 that reproduces the frozen world-build digest — against `/mnt/d/Work/personal/phoenix`: **no `DATA RACE` report**, with the single expected freeze red and nothing else. That covers the `liveProgress` handle shared between the run loop and the signal goroutine, which is the concurrency the reviewer's manual audit predicted clean. The file-level append hazard the reviewer noted separately (P3-N8) is not visible to the detector and is closed by construction, the log's `append` now holding a mutex.

Note for the reviewer: `pre-validation-artifacts.json` is edited in this payload commit, which is normally a refreeze-only file. The edit is confined to the `gates` block — the three new post-closure fields — and touches no digest, no `frozen` flag, and no accepted-candidate pointer. The digests in that document stay stale until the refreeze commit, which is why the freeze test is red.

## Refreeze plan after acceptance

1. Re-pin the `scheduled_runner` block in `pre-validation-artifacts.json`: 26 files (the existing 21 plus `process_events.go`, `execution_resilience.go`, `execution_resilience_test.go`, `closure_verdict_corpus_test.go`, and the pinned corpus expectations JSON), new set digest, `candidate_commit` = this payload commit, review record and its raw digest, `replaces_candidate_commit` = `54e256e9f4347d34844a3ef9b2156600360dcd13`, and the accepted-findings list extended with any residual the review records.
2. Update the freeze-test constants in `artifact_freeze_test.go` (`verifyReplacementRunnerFreeze`: candidate artifact path and raw digest, review path/commit/digest, replaces-commit, findings, notes, file count 26).
3. Convert this payload's working-tree guard to the commit-pinned historical form via `verifyCandidateSetAtCommit`.
4. Decision document 0026 (`IMPORT_AND_FREEZE`).

The local-artifact blocks (`runtime_invocation_and_exact_per_arm_system_prompts`, `arm_a_schemas`, `world_definition_and_world_build_digest`, `grader_digest`) do **not** move: this payload changes none of the bytes they cover, so their shared `candidate_commit`, local-candidate digest, and review pointer stay at the decision-0025 values and `verifyAcceptedLocalArtifacts` keeps passing throughout.

## Accepted residuals

- The attestation rule is scoped to the validation tranche; authoring resumes are exempt by construction and by pinned test.
- The process-event log cannot attribute a kill to its sender on either supported host. The runner records the unavailability honestly and pushes establishment onto the custodian's written attestation, as section 9 requires; a custodian who attests falsely is outside what any runner check can reach.
- `requirePostClosureAudit` runs inside `requireValidationGate` after the opening-flag checks, so a closed gate still reports "gate is closed" first. Both are refusals; the ordering is only about which message surfaces.
- The runner verifies that the attestation *claims* digest identity and that its recorded schedule and world-build digests equal the live ones. Byte-level re-verification of the whole frozen set at resume remains the custodian's step, recorded in the attestation, exactly as section 9 words condition (b).
- The section 9 attestation rule is scoped to the validation tranche; authoring resumes are exempt by construction and by pinned test. The first reviewer judged this defensible and verified that the exemption cannot launder a validation execution, because `validateScheduleForTranche` refuses a schedule whose tranche does not match the flag. It is recorded here as a knowing acceptance, not an oversight.
- Nothing in the runner writes `validation_execution_status`, so the post-closure gate's application after a *future* non-completed execution depends on the chair setting that field at closure, under the freeze workflow. The gate is armed by discipline for the next closure, and by frozen bytes for the current one.
- Carried from earlier cycles: `classifyScheduledOutputEntries` was at gocyclo 15 with zero headroom and is decomposed here (`scheduledOutputEntryKind`, `recordScheduledAssignment`, `finishScheduledInventory`), which retires carried finding P2-1; `verifyScheduledCheckpoint` is likewise split into identity and progress halves. `applyScheduledResume` still returns `done=true` with a non-nil error on stop paths (carried P2-3), and the `validateExternalGrade` manual/fail precedence divergence is still unfixed and still unreachable under v5 (carried P3).
