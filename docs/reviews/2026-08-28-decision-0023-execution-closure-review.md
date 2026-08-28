# Third independent review — decision-0023 execution closure, revision 3

- **Reviewer role:** independent evaluation reviewer. I did not author the closure record. I performed the two prior independent reviews that returned `REVISE` and re-audited the complete revision rather than treating either prior prescription as predetermined acceptance criteria.
- **Date:** 2026-08-28
- **Pinned candidate commit:** `238fcf53248671634e407993896f00751d3d7cfb`
- **Candidate record:** `docs/reports/2026-08-28-decision-0023-execution-closure-record.md`
- **Independently computed candidate LF-normalized UTF-8 SHA-256:** `sha256:aea6154d48e2cfb1fd0c96fd6ac34b36149232a97692537bc2614c16b2c3eb4e`
- **Verdict:** **ACCEPT** — 0 P0 · 0 P1 · 0 P2 · 0 P3

## Judgment

The revision is an accurate and sufficiently complete protocol-v5 § 9 closure record for the interrupted execution closed by decision 0023. It classifies the cause conservatively, distinguishes committed evidence from chair attestations and inference, reconstructs the archive custody correctly, acknowledges the failed cause-ordering condition, permanently rules out resume, exposes the surviving conditioning risks, and authorizes nothing.

`host_restart` is the correct classification against the frozen cause list. Decision 0023 and the exploratory report are committed chair records of the restart. They support the classification, but no surviving host log establishes the restart's initiator or reason. The candidate now preserves that limit instead of converting the frozen category into proof that this particular interruption was unengineered.

The raw 40 stderr bytes are **neutral** toward the cause and initiator of the original interruption. The chair's separate provenance attestation makes the artifact corroborative only of the reported post-restart sequence: a relaunch was attempted after the restart and failed at the WSL2 mount step. The artifact is not original-termination evidence. No committed termination artifact was omitted: the available set is decision 0023, the exploratory report, the custody index, the archive, and the captured relaunch stderr. The chair attests that no additional host event log, supervisor record, shell history, launcher output, or contemporaneous note exists.

## Checks performed

I independently performed the following checks at the pinned commit:

- verified the candidate commit and its commit metadata;
- verified that the working tree contained exactly the one permitted untracked review prompt before this record was created;
- read the candidate, decisions 0022, 0023, and 0026, the exploratory report, protocol v5's `amendment.execution_resilience`, the validation execution boundary and post-closure gate, `pre-validation-artifacts.json`, the custody index, and the evidence archive;
- computed the candidate digest after CRLF→LF and CR→LF normalization and confirmed that raw and normalized digests are identical;
- independently enumerated every archive file and classified assignment, runtime, trial, and grade records;
- independently hashed the raw checkpoint and rebuilt the ordinal-name custody aggregate from every raw file;
- parsed every assignment, checked launch-index uniqueness and continuity, compared the full archive prefix with the frozen v4 schedule, and checked declared assignment/trial/grade/runtime paths;
- reconstructed pairing keys and reconciled the checkpoint with the one post-checkpoint launch;
- checked for missing, duplicate, mismatched, unclassified, and orphaned evidence;
- compared the execution-era accepted payload with the candidate commit for runner, world-build, grader, validation corpus, manifest, and schedule identities;
- verified that the validation corpus, manifest, and v4 schedule remain byte-identical but retired, while the runner set, world build, grader, protocol, and arm construction changed;
- verified the gate state: `may_open_validation: false`, `may_open_held_out: false`, execution status `indeterminate`, and all four post-closure fields empty; and
- verified that the only revision-3 change is the closure record's corrected treatment of the unestablished 72-hour expiry.

No model was run, no private grade was obtained, no arm effect or outcome statistic was recomputed, and no new validation or held-out outcome was produced. The already-committed exploratory outcome statements were read only as required source context.

## Reconstructed custody values

| Check | Independently reconstructed value |
| --- | --- |
| Archive files | 5,649 |
| Assignment records | 1,476 |
| Runtime streams | 1,476 |
| Trial records | 1,348 |
| Grade records | 1,348 |
| Assignment launch-index range | 0–1475 |
| Distinct assignment launch indices | 1,476 |
| Assignment gaps or duplicates | none |
| `scheduled-summary.json` | absent |
| Checkpoint raw SHA-256 | `sha256:f35ec2f891f04725ea787aa103205fbf49573796503294bcebeb443566a557f4` |
| Custody aggregate SHA-256 | `sha256:3dff34f5c5e6fbd1958497cdca961b5a04fda1f7c84fd8215bb5f49bfb710241` |
| Checkpoint next launch index | 1475 |
| Checkpoint complete pairing keys | 295 |
| Checkpointed graded trials | 1,347 |
| Checkpointed terminal failures without trial or grade | 128 |
| Post-checkpoint evidence | launch 1475, arm C, pairing key 295 |
| Schedule-prefix mismatches | none |
| Orphaned assignments, streams, trials, or grades | none |

The arithmetic reconciles: 295 complete pairing keys × five distinct A–E arms = 1,475 launches at indices 0–1474. Launch 1475 is the single recorded member of incomplete pairing key 295. Trial and grade index sets are identical; every runtime stream is declared by its assignment; every retained assignment matches the corresponding frozen schedule entry.

## Attestations, ordering, and resume eligibility

The record correctly distinguishes:

- committed statements in decisions 0022 and 0023 and the exploratory report;
- chair attestations concerning the relaunch stderr's provenance, absence of further evidence, and contemporaneous cause attribution/outcome inspection; and
- inference from the byte sequence, host layout, and mechanical archive structure.

It does not claim that the § 9 cause classification preceded outcome inspection. The chair confirms the opposite. That failed ordering condition is handled conservatively and cannot be repaired retrospectively.

Resume is permanently unavailable on independently dispositive grounds: outcomes were inspected; the cause classification did not precede outcome inspection; the runner, world-build, grader, protocol, and arm construction no longer satisfy the execution-era identity set; and v5 retired arm E, so a v4 five-arm prefix cannot form a tranche with a v5 A–D remainder. The candidate correctly declines to use the 72-hour backstop as a ground because its expiry had not been established when the record was authored.

The lapse-explanation clause is correctly treated as inapplicable. This execution closed after a host restart, not by clock lapse, and outcome inspection had already made resume unavailable.

## Post-closure validity channel

The closure record exposes the unresolved possibility that the restart was participant-initiated or otherwise engineered; the surviving evidence cannot decide that question. It also discloses the known conditioning channel: diagnostics from this burned archive informed the authoring-side v5 repairs. Protocol v5 permits that engineering use while forbidding either burned archive from seeding new v5 cases, labels, or witnesses, except for the separately frozen absence-message calibration carve-out.

No later attempt may resume, splice, extend, or continue this archive. Any later Gate 1A attempt must use a fresh, disjoint, independently generated and sealed v5 tranche. The retained v4 validation corpus, manifest, and five-arm schedule are retired and cannot serve as that tranche.

## Accepted evidence limitations and residual risks

- No `process-events.jsonl` or durable `resumes` counter exists because the execution predates the v5 § 9 payload. They must not be fabricated retrospectively.
- No surviving evidence establishes who initiated the host restart or why. `host_restart` is the frozen cause classification, not proof of non-participant provenance.
- The 40 stderr bytes establish no original-interruption fact by themselves; their post-restart provenance rests on the chair's custodian attestation.
- The archive independently establishes the recorded world-build identity and schedule prefix, but does not independently re-establish every execution-era runner, corpus, or grader byte.
- Custody reconstruction establishes the committed archive's internal coherence and pinned bytes; it cannot supply a missing contemporaneous host log or independently prove the restart was unengineered.
- The archive informed v5 engineering work. Future tranche disjointness, independent generation, path witnesses, sealing, and refreeze remain separate controls and blockers.
- Section 3 witnesses, the section 6 witness-derived cap, candidate generation and sealing, the authoring dry run, Task 0.6, and a later chair opening decision remain outside this closure review and unresolved.

These limitations are disclosed and conservative. None permits this tranche to resume or this review to authorize another execution.

## Findings and authorization boundary

There are no P0, P1, P2, or P3 findings. Accordingly, no finding permits an outcome-conditioned retry and no record-quality correction remains required for this candidate.

This review opens neither outcome gate and authorizes no execution. Validation and held-out remain closed. This review does not set any post-closure gate field and is not a Gate 1A opening decision.

## Exact next artifact

The exact next artifact is a **separate payload commit** that sets the four post-closure gate fields in `experiments/frontier-v1/pre-validation-artifacts.json` to this accepting review, its `ACCEPT` verdict, and its LF-normalized digest, while naming the independent review of that payload commit:

- `validation_execution_closure_review`
- `validation_execution_closure_review_verdict`
- `validation_execution_closure_review_lf_normalized_utf8_sha256`
- `validation_execution_closure_fields_payload_review`

That payload is not a chair gate patch and must not open Gate 1A. It must receive its own independent review, whose reviewer reads this record and rules that it is an accepting independent review of the closure, followed by the required replacement refreeze. Only after that separate reviewed payload and refreeze may a distinct chair decision consider opening validation while keeping held-out closed and satisfying every other blocker.
