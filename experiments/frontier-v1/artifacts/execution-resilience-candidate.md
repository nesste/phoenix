# Protocol-v5 section 9 execution-resilience candidate

Status: **review candidate, not frozen.** Both outcome gates stay closed (`may_open_validation: false`, `may_open_held_out: false`). No model was run, no private grade obtained, and no validation outcome observed while authoring this payload.

Base commit: `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, section 5 absence-acceptance refreeze).
Inventory: [`gate-1a-execution-resilience-candidate.json`](gate-1a-execution-resilience-candidate.json), 16 files.

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
| (c) resume within 72 hours | `interrupted_at` → `resume_authorized_at`, timestamps required in order | refuses |
| (d) cause classified before inspection and before the resume decision, and on the frozen list | `cause_classified_before_inspection_and_resume_decision`, `cause` ∈ {`host_restart`, `power_loss`, `hardware_failure`, `process_kill`} | refuses |
| kill establishment | `kill_principal_established_as_non_participant` plus `process_event_log_sha256` re-verified against the log on disk | refuses |
| at most one resume | checkpoint `resumes` field, `resume_index` must equal `resumes + 1` | refuses |
| promptness | `promptness_statement` required; a `resume_authorized_at` later than `conditions_verified_at` additionally requires `lapse_explanation` | refuses |

**OOM is struck** — it is simply absent from `frozenResumeCauses`, so it falls to the generic "not on the frozen outcome-uncorrelated cause list" refusal along with any other unlisted cause, including the empty string. A refusal is not an abandonment: the run stops without launching, and the tranche closes indeterminate under the existing stop rules. Abandonment by choice has no representation anywhere in the runner.

The single-resume rule is carried in the durable checkpoint, which gains a `resumes` counter: a fresh run writes 0, the one permitted resume writes 1, and a checkpoint already at 1 is refused before any attestation is read (`already used its single resume`). `verifyCheckpointProgress` also bounds the counter, so a hand-edited checkpoint claiming a higher count is a mismatch, not a licence.

**Scope.** The attestation rule applies to `--tranche validation`. Authoring resumes are development loops and are exempt (`TestAuthoringResumeIsExemptFromTheAttestationRule` pins the exemption so it is a decision, not an oversight). Reviewer: this scoping is the one place the payload chose narrower than the proposal's unqualified wording, on the grounds that section 9's subject throughout is the validation tranche and its indeterminate closure.

### 4. Post-closure gate

`requirePostClosureAudit` in `validation_gate.go`: when the freeze document's `validation_execution_status` is anything other than empty or `complete`, `requireValidationGate` refuses unless `validation_execution_closure_review` names an existing record and `validation_execution_closure_review_verdict` is `ACCEPT`. The gate document gains those two fields plus `validation_execution_closure_gate` stating the rule, and `verifyPostClosureGateState` in the freeze test asserts the rule text and that both fields stay **empty**.

**Consequence the chair should note.** The 1,475-launch execution closed by decision 0023 ended without a completed schedule. The gate therefore applies to it retroactively, exactly as section 9 specifies ("after any validation execution that ends without a completed schedule"). A new blocker is registered in `protocol.json`: an independently reviewed and committed closure record for that execution is now required before any later validation execution, on top of the existing blockers. This payload does not produce that review — it is a separate artifact, an independent audit of decision 0023's cause classification, custodian logs, and attestations.

### 5. Protocol and boundary records

`protocol.json` gains `amendment.execution_resilience` carrying the whole rule, and its blocker list is rewritten: the section 9 clause is struck from the remaining-v5-amendment blocker (leaving path witnesses and the witness-derived cap), and the closure-record blocker is inserted, taking the list from 7 to 8 entries. `validation-execution-boundary.md` gains an "Execution resilience, mandatory resume, and the post-closure gate" section carrying the host requirements (automatic OS restarts disabled, no interactive login sessions, supervised custodian process), the three per-checkpoint records, the five numbered resume conditions, and the post-closure gate including its retroactive application to the decision-0023 closure. `README.md` describes the two new output records and the resume flag.

## Deliberate freeze red

Exactly one test is expected to fail at this commit, and must:

- `TestPreValidationFreezeMatchesAcceptedCandidates` — the frozen `scheduled_runner` set digest and per-file digests still pin the decision-0025 bytes.

`TestLocalArtifactCandidateMatchesImplementation` and `TestValidationBuildReproducesFrozenWorldBuildDigest` stay **green**: the Phoenix binary is untouched (nothing under `cmd/phoenix` or its import graph changed), so the world-build identity remains `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4`, and the grader digest is unchanged at `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b` because no corpusctl non-test source, `go.mod`/`go.sum`, or `schema/*.schema.json` file was touched. No authoring label, manifest, or grading_script pin moves in this payload.

Everything else is green: root suite, corpusctl suite, surface-spike, `go vet`, staticcheck, `gocyclo -over 15`, `dupl -t 100`, `validate-spec`, `validate-authoring`.

Note for the reviewer: `pre-validation-artifacts.json` is edited in this payload commit, which is normally a refreeze-only file. The edit is confined to the `gates` block — the three new post-closure fields — and touches no digest, no `frozen` flag, and no accepted-candidate pointer. The digests in that document stay stale until the refreeze commit, which is why the freeze test is red.

## Refreeze plan after acceptance

1. Re-pin the `scheduled_runner` block in `pre-validation-artifacts.json`: 24 files (the existing 21 plus `process_events.go`, `execution_resilience.go`, `execution_resilience_test.go`), new set digest, `candidate_commit` = this payload commit, review record and its raw digest, `replaces_candidate_commit` = `54e256e9f4347d34844a3ef9b2156600360dcd13`, and the accepted-findings list extended with any residual the review records.
2. Update the freeze-test constants in `artifact_freeze_test.go` (`verifyReplacementRunnerFreeze`: candidate artifact path and raw digest, review path/commit/digest, replaces-commit, findings, notes, file count 24).
3. Convert this payload's working-tree guard to the commit-pinned historical form via `verifyCandidateSetAtCommit`.
4. Decision document 0026 (`IMPORT_AND_FREEZE`).

The local-artifact blocks (`runtime_invocation_and_exact_per_arm_system_prompts`, `arm_a_schemas`, `world_definition_and_world_build_digest`, `grader_digest`) do **not** move: this payload changes none of the bytes they cover, so their shared `candidate_commit`, local-candidate digest, and review pointer stay at the decision-0025 values and `verifyAcceptedLocalArtifacts` keeps passing throughout.

## Accepted residuals

- The attestation rule is scoped to the validation tranche; authoring resumes are exempt by construction and by pinned test.
- The process-event log cannot attribute a kill to its sender on either supported host. The runner records the unavailability honestly and pushes establishment onto the custodian's written attestation, as section 9 requires; a custodian who attests falsely is outside what any runner check can reach.
- `requirePostClosureAudit` runs inside `requireValidationGate` after the opening-flag checks, so a closed gate still reports "gate is closed" first. Both are refusals; the ordering is only about which message surfaces.
- The runner verifies that the attestation *claims* digest identity and that its recorded schedule and world-build digests equal the live ones. Byte-level re-verification of the whole frozen set at resume remains the custodian's step, recorded in the attestation, exactly as section 9 words condition (b).
- Carried from earlier cycles: `classifyScheduledOutputEntries` was at gocyclo 15 with zero headroom and is decomposed here (`scheduledOutputEntryKind`, `recordScheduledAssignment`, `finishScheduledInventory`), which retires carried finding P2-1; `verifyScheduledCheckpoint` is likewise split into identity and progress halves. `applyScheduledResume` still returns `done=true` with a non-nil error on stop paths (carried P2-3), and the `validateExternalGrade` manual/fail precedence divergence is still unfixed and still unreachable under v5 (carried P3).
