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

**Two properties of this cause matter under § 9 and both hold.** A host restart is outcome-uncorrelated: it is unrelated to what any trial produced, unlike OOM, which § 9 strikes precisely because transcript-heavy failing trials make memory pressure correlate with outcomes. And it is not a process kill, so the § 9 requirement that no project participant, account, or agent initiated the termination does not arise as a question needing the establishment § 9 demands — decision 0023's committed statement that no process required a forced stop is consistent with that, and no evidence contradicts it.

**One ordering condition of § 9 is not satisfied, and cannot be retrospectively.** § 9 requires the cause to be classified and attested *before any outcome inspection and before the resume decision point*. For this execution that ordering did not occur and could not have: the rule did not exist, and decision 0023 and the exploratory report were authored together, with the cause attribution and the diagnostic inspection contemporaneous. This record does not claim otherwise. The consequence is developed in section 4 — it is not a defect to be repaired but a fact that makes resume permanently unavailable, which is the conservative direction.

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

The durable checkpoint reads: schedule digest `sha256:b38a0eaa…5813`, world build `sha256:b5a26d5e…a2c4`, run budget 300 USD, spend 60.4877118 USD, next launch index 1475, completed pairing keys 295. The observed world build equals the frozen world build of the day, so the execution ran the bytes it was authorized to run.

**The arithmetic reconciles, and it is worth showing because it is the mechanical evidence that the interruption was clean.** 295 complete pairing keys × 5 arms (v4 A–E) = 1,475 launches at indices 0–1474, exactly the checkpoint's `next_launch_index`. The archive holds one further assignment, launch 1475, written after the last checkpoint and belonging to incomplete pairing key 295; decision 0023 and the exploratory report both exclude it from every reported metric, correctly, because including it would break the precommitted pairing contract. Of the 1,475 checkpointed launches, 1,347 produced graded trials and 128 ended in runtime-terminal failure with no trial; the archive's 1,348 grades are those 1,347 plus launch 1475's. There is no gap, no duplicate, and no orphaned evidence.

**The absent scheduled summary is the second mechanical indication** that the schedule did not complete, independent of the checkpoint. Both agree.

### The one piece of evidence no committed document has explained

The custody index records 40 bytes of captured standard error:

```
mkdir /mnt/d: file exists
exit status 1
```

Neither decision 0023 nor the exploratory report mentions it. It is the only captured artifact of process termination, so this record must address it rather than leave it uninterpreted.

`/mnt/d` is the WSL2 mount point of the Windows `D:` volume, on which the repository lives. `mkdir /mnt/d` failing with `file exists` is what a mount-setup step produces when the mount is already present. This is consistent with a **post-restart re-launch attempt**: the host came back, WSL2 restarted with the drive already mounted, a wrapper step tried to create the mount point, failed, and exited 1. It is not consistent with a failure inside the Go runner, the Claude process, or the custodian grader, none of which creates mount points.

**On this reading the stderr corroborates the `host_restart` classification rather than qualifying it**, and it also indicates that the re-launch attempt failed immediately, before touching the output directory — which is consistent with the archive containing no evidence written after launch 1475. **This reading is mine, from the byte sequence and the host layout; it is not a committed custodian statement.** It is flagged for the independent reviewer and for chair confirmation in section 5.

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

Neither is a gap to be filled. Both are facts, and both point the same way: resume is unavailable. Section 5 lists what the chair must confirm before this record is reviewed.

## 4. Resume eligibility, and the lapse explanation

**§ 9's lapse-explanation clause does not apply, because no resume was attempted and none could have been.** § 9 requires a written explanation where a closure occurs by clock lapse — the 72-hour backstop expiring while resume remained possible. That is not this case.

Resume was unavailable from the moment the outcomes were inspected, and the mandatory-resume rule that would have compelled an attempt did not exist at the time. Under the v4 protocol in force, the chair's options on an incomplete schedule were to close the execution as indeterminate or to resume it; decision 0023 closed it, and the exploratory report states plainly that the tranche *"cannot be resumed or reused for another Gate 1A claim."* Every subsequent condition confirms this:

- **Outcome inspection burns the tranche.** § 9 preserves the v4 rule that inspecting any outcome before resume burns the tranche. Inspection happened, deliberately, for engineering diagnosis. That alone is dispositive and permanent.
- **The frozen bytes are no longer identical.** § 9 condition (b) requires digest identity of the runner, world, corpus, schedule, and grader at resume. All have since moved: decisions 0024, 0025 and 0026 replaced the protocol, the grader (`sha256:8146a68a…ebf0b`), the runner set, and the world build (`sha256:425bab1c…dae4` against this execution's `sha256:b5a26d5e…a2c4`). The archive's schedule digest `sha256:b38a0eaa…5813` is the retired v4 five-arm schedule, which the A–D construction now refuses by design.
- **The arms no longer exist.** This execution ran five arms; v5 retired arm E. A resumed prefix and a v5 remainder could not form one tranche.
- **The 72-hour backstop is long past**, and the cause classification could not have preceded the resume decision point in any case.

**Conclusion: this execution can never be resumed, and no future validation execution will be a continuation of it.** Any Gate 1A attempt requires a new disjoint sealed tranche, which decision 0023 already required and which `protocol.json` carries as a separate blocker.

## 5. What the chair must confirm before this record is independently reviewed

This record is authored from committed documents and from evidence recomputed on 2026-08-28. Three points are outside what the archive can establish and require the chair's confirmation, because they are statements about what happened and what was known:

1. **The stderr reading in section 2** — that the 40 captured bytes are a post-restart re-launch attempt failing on an already-present WSL2 mount, and not a failure inside the runner, the Claude process, or the custodian grader. If the chair's recollection or any uncommitted host record differs, section 2 must be corrected before review.
2. **That no other termination evidence exists** — no host event log, supervisor record, shell history, or note beyond decision 0023, the custody index, and the archive itself. If more exists it should be committed with this record, since § 9 asks for the custodian logs and this record's answer is that they are limited to what section 2 enumerates.
3. **That the cause classification was not made before outcome inspection.** Section 1 states this as a limitation. The chair is the only person who can confirm the ordering. If the cause was in fact determined before any outcome was inspected, section 1 should say so — it would not change the resume conclusion, but the record should be accurate.

## 6. What this record establishes for the post-closure gate

The gate asks whether the next validation execution would be conditioned on side signals from this one. On the evidence:

- **The execution died of a cause unrelated to its outcomes** — a host restart, corroborated by the only captured termination artifact, on the frozen outcome-uncorrelated list, with a clean checkpoint boundary and reconciling arithmetic showing no partial or manipulated evidence.
- **It was not engineered.** The execution ran to 82% of its schedule and stopped on an external event. Nothing in the archive suggests a stop chosen for what the outcomes showed, and the spend of 60.49 USD against a 300 USD ceiling shows the budget-stop rule was nowhere near engaged.
- **The retry is not a retry of this tranche.** It is a fresh, disjoint, independently sealed tranche, under a protocol whose arms, grader, world build, schedule and runner have all since been replaced through three reviewed and refrozen cycles.
- **The one genuine conditioning channel is disclosed and is not this record's to close.** The v5 protocol changes were themselves informed by this archive's diagnostics — that is exactly what decision 0023 authorized the archive for, and what protocol v5 § 10 permits as engineering repair while forbidding the archive from seeding v5 cases, labels, or witnesses. The next tranche's *content* must be disjoint and independently generated; that is a separate blocker with its own controls.

**This record does not authorize any validation execution.** It is one of several outstanding blockers, and both outcome gates remain closed.

## Verification note

Every archive figure in section 2 was recomputed on 2026-08-28 from `experiments/frontier-v1/results/scheduled-validation-2` and compared against the committed custody index: file and record counts, launch-index range and distinctness, absence of the scheduled summary, the checkpoint's raw digest, and the ordinal-name aggregate index digest. All reproduce exactly. No outcome statistic was recomputed for this record and none is reported here; the arm, class, termination and cost figures remain in the exploratory report, where decision 0023 placed them.
