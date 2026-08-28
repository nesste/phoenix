# Closure record — the interrupted 1,475-launch validation execution closed by decision 0023

- **Status:** closure record, awaiting independent review. Authorizes nothing.
- **Date:** 2026-08-28
- **Execution closed by:** `docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md`
- **Opening decision for that execution:** `docs/decisions/0022-protocol-v4-gate-1a-validation-reopening-after-world-compatibility-repair.md`
- **Evidence archive:** `experiments/frontier-v1/results/scheduled-validation-2` (diagnostic-only, burned)
- **Custody index:** `experiments/frontier-v1/results/scheduled-validation-2-custody.json`
- **Required by:** protocol v5 § 9 post-closure gate, installed by decision 0026

## Why this record exists

Protocol v5 § 9, imported by decision 0026, states that after any validation execution ending without a completed schedule — interruption, budget stop, or lapse — no subsequent validation execution may be authorized until a closure record carrying the cause classification, custodian logs, attestations, and the lapse explanation if any has been independently reviewed and committed. The gate is deliberately retroactive: the execution closed by decision 0023 ended without a completed schedule, so it falls inside the rule, and `protocol.json` carries the corresponding blocker.

The gate's purpose is narrow and worth restating, because it determines what this record must establish. Repeated close-and-retry would condition an eventually completed tranche on side signals: if a chair may keep closing executions and starting new ones, the one that finally completes is selected partly by what the closed ones showed. The gate bounds that channel by requiring each retry to survive an independent audit of **why the last execution died**. This record exists to answer that question and to expose anything that would make the next attempt conditional on this one's outcomes.

Decision 0023 is a sound closure decision and this record does not supersede it. But it was written on 2026-08-27, before § 9 existed, and it is not § 9-shaped: it does not classify the cause against a frozen cause list, does not enumerate the custodian evidence as such, and does not address resume eligibility, because no mandatory-resume rule was in force.

## 1. Cause classification

**Classified cause: `host_restart`.** On the § 9 frozen outcome-uncorrelated cause list.

The supporting statements are already committed. Decision 0023 records that *"The host restart terminated the external launcher, Go runner, Claude process, and custodian process. A process audit on 2026-08-27 found none of them running. No process required a forced stop."* The exploratory report records that the execution *"reached 1,475 of 1,800 checkpointed launches before the PC restart."*

**The classification and its evidence limit must be kept separate.** Section 9 places `host_restart` on its frozen outcome-uncorrelated cause list, unlike OOM, and its explicit originating-principal test applies to `process_kill`, not to this classification. But the category does not establish the provenance of this particular restart. No host log or other available evidence establishes whether it was automatic, manually initiated by a project participant, or initiated for a reason related to anything a participant knew. Decision 0023's statement that no process required a forced stop concerns the child processes, not who or what initiated the host restart. The cause classification is therefore `host_restart`; the initiator is unresolved, and this record does not claim the interruption itself was proven unengineered.

**One ordering condition of § 9 is not satisfied, and cannot be retrospectively.** § 9 requires the cause to be classified and attested *before any outcome inspection and before the resume decision point*. For this execution that ordering did not occur and could not have: the rule did not exist, and decision 0023 and the exploratory report were authored together, with the cause attribution and the diagnostic inspection contemporaneous. **The project chair confirmed on 2026-08-28 that the two were contemporaneous and that the cause was not classified before any outcome was inspected.** This record does not claim otherwise. The consequence is developed in section 4 — it is not a defect to be repaired but a fact that makes resume permanently unavailable, which is the conservative direction.

## 2. Custodian logs and the evidence that stands in their place

**The § 9 process-event log does not exist for this execution.** `process-events.jsonl` is an artifact the § 9 payload created; the execution predates it. So does the durable checkpoint's `resumes` counter. No custodian process-event log can be produced for this run, and none should be fabricated.

What does exist, and what I verified rather than transcribed. Every figure below was recomputed from the archive on 2026-08-28 and reproduces the committed custody index exactly:

| Fact | Recomputed | Custody index |
| --- | --- | --- |
| Archive file count | 5,649 | 5,649 |
| Assignment records | 1,476 | 1,476 |
| Runtime streams | 1,476 | 1,476 |
| Trial records | 1,348 | 1,348 |
| Grade records | 1,348 | 1,348 |
| Launch indices | 0–1475, 1,476 distinct | 0–1475 |
| `scheduled-summary.json` | absent | absent |
| Checkpoint raw SHA-256 | `sha256:f35ec2f891f04725ea787aa103205fbf49573796503294bcebeb443566a557f4` | identical |
| Aggregate index SHA-256 | `sha256:3dff34f5c5e6fbd1958497cdca961b5a04fda1f7c84fd8215bb5f49bfb710241` | identical |

The durable checkpoint reads: schedule digest `sha256:b38a0eaa…5813`, world build `sha256:b5a26d5e…a2c4`, run budget 300 USD, spend 60.4877118 USD, next launch index 1475, completed pairing keys 295. The observed world build equals the frozen world build of the day, establishing the authorized world-build identity; the checkpoint's schedule digest and exact prefix comparison establish the schedule identity separately. The archive does not independently re-establish every runner, corpus, and grader byte, and this record does not claim that it does.

**The arithmetic reconciles, and it is worth showing because it is the mechanical evidence that the interruption was clean.** 295 complete pairing keys × 5 arms (v4 A–E) = 1,475 launches at indices 0–1474, exactly the checkpoint's `next_launch_index`. The archive holds one further assignment, launch 1475, written after the last checkpoint and belonging to incomplete pairing key 295; decision 0023 and the exploratory report both exclude it from every reported metric, correctly, because including it would break the precommitted pairing contract. Of the 1,475 checkpointed launches, 1,347 produced graded trials and 128 ended in runtime-terminal failure with no trial; the archive's 1,348 grades are those 1,347 plus launch 1475's. There is no gap, no duplicate, and no orphaned evidence.

**The absent scheduled summary is the second mechanical indication** that the schedule did not complete, independent of the checkpoint. Both agree.

**The project chair confirmed on 2026-08-28 that no other termination evidence exists** beyond decision 0023, the exploratory report, the custody index, the evidence archive, and the captured standard error addressed below. In particular, there is no additional host event log, supervisor record, shell history, launcher output, or contemporaneous note to enumerate. This is an evidence limitation, not an omitted artifact; the record therefore identifies the complete available custodian evidence for this execution.

### The captured relaunch artifact no committed document has explained

The custody index records 40 bytes of captured standard error:

```
mkdir /mnt/d: file exists
exit status 1
```

Neither decision 0023 nor the exploratory report mentions it. It is captured standard error from the later relaunch attempt, not an artifact of the original termination; this record must address it without treating its raw content as evidence of how the interruption occurred.

`/mnt/d` is the WSL2 mount point of the Windows `D:` volume, on which the repository lives. `mkdir /mnt/d` failing with `file exists` is what a mount-setup step produces when the mount is already present. This is consistent with a **post-restart re-launch attempt**: the host came back, WSL2 restarted with the drive already mounted, a wrapper step tried to create the mount point, failed, and exited 1. It is not consistent with a failure inside the Go runner, the Claude process, or the custodian grader, none of which creates mount points.

**The raw stderr is neutral toward the cause and initiator of the original interruption.** Its link to the restart sequence comes from the chair's provenance attestation, not from the 40 bytes alone. **The project chair confirmed on 2026-08-28** that the bytes came from a post-restart relaunch attempt failing on an already-present WSL2 mount, not from the runner, the Claude process, or the custodian grader. With that attestation, the artifact corroborates that a relaunch was attempted after the reported restart and failed immediately, before touching the output directory — consistent with the archive containing no evidence written after launch 1475. It does not establish who initiated the restart or why.

## 3. Attestations

§ 9 requires attestations in a closure record. The following are **already committed** in decision 0023 and are cited here, not re-authored:

- The host restart terminated the external launcher, Go runner, Claude process, and custodian process.
- A process audit on 2026-08-27 found none of them running.
- No process required a forced stop.
- The execution produced no `scheduled-summary.json`; the full frozen schedule was not completed; the result is indeterminate regardless of the partial effect estimates.
- Outcomes were inspected for engineering diagnosis, which burns this validation tranche for future confirmatory use.

**Attestations this record does not supply, and must not.** An attestation is a statement by a person about facts within their knowledge. Two § 9 attestations cannot be given for this execution by anyone, at any time:

1. That the cause was classified before any outcome inspection and before the resume decision point. It was not (section 1).
2. That no outcome artifact was inspected between interruption and a resume. Outcomes were inspected, deliberately and on the record.

Neither is a gap to be filled. Both are facts, and both point the same way: resume is unavailable. Section 5 records the relevant chair confirmations.

## 4. Resume eligibility, and the lapse explanation

**§ 9's lapse-explanation clause does not apply because this execution closed after a `host_restart`, not by clock lapse.** Section 9 requires a written explanation where the 72-hour backstop expires while resume remains possible. Under v4, resume was an available option before outcome inspection; it became permanently unavailable when the outcomes were inspected for diagnosis. No clock-lapse closure occurred.

Resume was unavailable from the moment the outcomes were inspected, and the mandatory-resume rule that would have compelled an attempt did not exist at the time. Under the v4 protocol in force, the chair's options on an incomplete schedule were to close the execution as indeterminate or to resume it; decision 0023 closed it, and the exploratory report states plainly that the tranche *"cannot be resumed or reused for another Gate 1A claim."* Every subsequent condition confirms this:

- **Outcome inspection burns the tranche.** § 9 preserves the v4 rule that inspecting any outcome before resume burns the tranche. Inspection happened, deliberately, for engineering diagnosis. That alone is dispositive and permanent.
- **The identity conditions required for resume no longer all hold.** § 9 condition (b) requires digest identity of the runner, world, corpus, schedule, and grader at resume. Decisions 0024, 0025 and 0026 changed the protocol, grader (`sha256:8146a68a…ebf0b`), runner set, and world build (`sha256:425bab1c…dae4` against this execution's `sha256:b5a26d5e…a2c4`). The validation corpus, manifest, and v4 five-arm schedule remain byte-identical, but they are retired: the A–D construction refuses that schedule by design, and the replacement v5 schedule has not yet been generated or sealed. Unchanged retired inputs do not restore whole-set identity or resume eligibility.
- **The arms no longer exist.** This execution ran five arms; v5 retired arm E. A resumed prefix and a v5 remainder could not form one tranche.
- **The 72-hour backstop is long past**, and the cause classification could not have preceded the resume decision point in any case.

**Conclusion: this execution can never be resumed, and no future validation execution will be a continuation of it.** Any Gate 1A attempt requires a new disjoint sealed tranche, which decision 0023 already required and which `protocol.json` carries as a separate blocker.

## 5. Chair confirmations

This record is authored from committed documents and from evidence recomputed on 2026-08-28. Three points lie outside what the archive can establish, because they are statements about what happened and what was known. The project chair was asked all three on 2026-08-28.

1. **The stderr reading in section 2 — confirmed.** The 40 captured bytes are a post-restart re-launch attempt failing on an already-present WSL2 mount, not a failure inside the runner, the Claude process, or the custodian grader. Section 2 records the confirmation.
2. **No further termination evidence — confirmed.** The chair confirms that no other termination evidence exists beyond decision 0023, the exploratory report, the custody index, the evidence archive, and the captured standard error addressed in section 2. Section 2 therefore enumerates the complete available custodian evidence for this execution.
3. **The cause-classification ordering — confirmed as stated.** The cause attribution and the diagnostic inspection were contemporaneous; the cause was not classified before any outcome was inspected. Section 1 stands, and the resume conclusion is unaffected: resume was already unavailable on the four other independent grounds in section 4.

## 6. What this record establishes for the post-closure gate

The gate asks whether the next validation execution would be conditioned on side signals from this one. On the evidence:

- **The classified cause is `host_restart`, on the frozen outcome-uncorrelated list**, with a clean checkpoint boundary and reconciling arithmetic showing no partial or manipulated archive evidence. The relaunch artifact, through the chair's provenance attestation, corroborates the post-restart sequence but not the restart's initiator or reason.
- **Whether the interruption was engineered cannot be established from the surviving evidence.** Nothing in the archive indicates a stop chosen for what the outcomes showed, the execution reached 82% of its schedule, and spend of 60.49 USD against a 300 USD ceiling rules out the budget-stop mechanism. But no host log or other record establishes that the restart was automatic or non-participant-initiated. That unresolved provenance is an accepted evidence limitation, not silently converted into proof. Its conservative consequence is that this tranche remains burned and permanently non-resumable.
- **The next attempt is not a retry of this tranche.** It must use a fresh, disjoint, independently sealed tranche under the changed A–D protocol, runner, grader, and world build. The old validation corpus, manifest, and five-arm schedule remain byte-identical but retired; a replacement v5 schedule is still unsealed and must not reuse them as a continuation.
- **The one genuine conditioning channel is disclosed and is not this record's to close.** The v5 protocol changes were themselves informed by this archive's diagnostics — that is exactly what decision 0023 authorized the archive for, and what protocol v5 § 10 permits as engineering repair while forbidding the archive from seeding v5 cases, labels, or witnesses. The next tranche's *content* must be disjoint and independently generated; that is a separate blocker with its own controls.

**This record does not authorize any validation execution.** It is one of several outstanding blockers, and both outcome gates remain closed.

## Verification note

Every archive figure in section 2 was recomputed on 2026-08-28 from `experiments/frontier-v1/results/scheduled-validation-2` and compared against the committed custody index: file and record counts, launch-index range and distinctness, absence of the scheduled summary, the checkpoint's raw digest, and the ordinal-name aggregate index digest. All reproduce exactly. No outcome statistic was recomputed for this record and none is reported here; the arm, class, termination and cost figures remain in the exploratory report, where decision 0023 placed them.
