# Third independent review prompt — decision-0023 execution closure, revision 3

You are an evaluation reviewer independent of the author of the closure record. Review the interrupted validation execution's closure against protocol v5 § 9. Do not run a model, obtain a private grade, modify an outcome artifact, open either outcome gate, or authorize a new execution. Treat every statement in the candidate as a claim to verify, not as settled fact.

## Why there is a third review

The first independent review of commit `3da0a920ac3d967f9a60a78cda0b753b29e6cb33` returned **REVISE** with one P1, two P2 findings, and two P3 findings:

1. The record treated the 40 stderr bytes as a termination artifact and claimed the restart was not engineered, although the raw bytes concern only a later relaunch and no evidence establishes who initiated the original restart or why.
2. It claimed all frozen bytes had moved, although the validation corpus, manifest, and five-arm schedule remain byte-identical but retired.
3. It inferred all authorized execution bytes from the matching world build, which establishes only world-build identity; the schedule identity has separate support.
4. Its exhaustive termination-evidence list omitted the exploratory report.
5. Its lapse explanation said no resume could have been attempted, although v4 allowed resume before outcome inspection.

Revision 2 resolved all five, but its second independent review returned **REVISE** with one P2: it claimed the 72-hour backstop was long past although only about 47.69 hours separated the archive's latest timestamp and the pinned candidate commit. Revision 3 removes elapsed time as a ground, states that expiry had not been established, and relies on the independently dispositive outcome-inspection, cause-ordering, identity, and arm-set grounds.

Recheck the whole record, not only this edit, and treat the prescriptions from both prior reviews as claims rather than predetermined acceptance criteria.

## Pinned candidate

- Candidate commit: `238fcf53248671634e407993896f00751d3d7cfb`
- Candidate record: `docs/reports/2026-08-28-decision-0023-execution-closure-record.md`
- Candidate record LF-normalized UTF-8 SHA-256: `sha256:aea6154d48e2cfb1fd0c96fd6ac34b36149232a97692537bc2614c16b2c3eb4e`
- The raw SHA-256 is identical because the committed file uses LF line endings.
- The working tree may contain exactly one untracked file: this prompt. Anything else dirty must be reported before review.

Verify the commit and digest before beginning. Digest normalization is CRLF→LF, then CR→LF, followed by SHA-256 over the UTF-8 bytes.

## Required sources

Read at least:

- `docs/reports/2026-08-28-decision-0023-execution-closure-record.md`
- `docs/decisions/0022-protocol-v4-gate-1a-validation-reopening-after-world-compatibility-repair.md`
- `docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md`
- `docs/reports/2026-08-27-protocol-v4-second-interrupted-validation-exploratory.md`
- `docs/decisions/0026-protocol-v5-execution-resilience-import-refreeze.md`
- `experiments/frontier-v1/protocol.json`, especially `amendment.execution_resilience`
- `experiments/frontier-v1/validation-execution-boundary.md`, especially the post-closure gate
- `experiments/frontier-v1/pre-validation-artifacts.json`
- `experiments/frontier-v1/results/scheduled-validation-2-custody.json`
- the evidence archive at `experiments/frontier-v1/results/scheduled-validation-2`

The chair's three confirmations are committed in the candidate and must be assessed as custodian attestations: the captured stderr came from the failed post-restart relaunch; no additional termination evidence exists; and cause attribution was contemporaneous with outcome inspection rather than preceding it.

## What to verify

### 1. Cause and termination evidence

- Decide whether `host_restart` is the correct frozen cause classification and whether the committed evidence supports it without overclaim.
- Assess the 40 captured stderr bytes and the chair's interpretation. Say whether they corroborate, qualify, contradict, or are neutral toward `host_restart`.
- Decide whether section 2 honestly and completely enumerates the available custodian evidence given that the execution predates `process-events.jsonl` and the durable `resumes` counter.
- Do not require nonexistent v5 artifacts to be fabricated. Do identify any actually committed termination artifact the record omitted.

### 2. Mechanical custody reconstruction

Independently reproduce, rather than transcribe, the archive file count; counts of assignments, streams, trials, and grades; launch-index range and distinctness; absence of `scheduled-summary.json`; checkpoint raw digest; and custody aggregate digest. Reconcile 295 complete five-arm pairing keys with checkpoint launch index 1475 and the one post-checkpoint launch. Report any mismatch, duplicate, gap, or orphan.

This is a closure audit, not a new diagnostic analysis. Do not recompute arm effects or introduce a new interpretation of outcomes.

### 3. Attestations, ordering, and resume eligibility

- Determine whether the record clearly distinguishes committed evidence, chair attestations, and inference.
- Confirm that it does not falsely claim the § 9 cause classification occurred before outcome inspection.
- Decide whether the failure of that ordering condition is handled conservatively and whether resume is permanently unavailable on the independent grounds listed in section 4.
- Decide whether the lapse-explanation clause is correctly treated as inapplicable rather than silently omitted.
- Check the claims that the frozen bytes and arm set have since changed.

### 4. Post-closure validity channel

- Decide whether the record adequately audits why the execution ended and exposes every plausible way a later attempt could be conditioned on this archive's side signals.
- Check that the next attempt is described only as a fresh, disjoint, independently sealed tranche, never a resume, splice, or continuation.
- Confirm that the record itself authorizes nothing and that validation and held-out remain closed.

### 5. Accuracy and completeness

Review the whole candidate for unsupported claims, internal contradictions, missing evidence limitations, misleading certainty, or requirements imported retrospectively in a non-conservative way. Separate closure-record defects from other project blockers such as § 3 witnesses, the § 6 cap, candidate sealing, the authoring dry run, and Task 0.6.

## Verdict and review record

Return exactly one verdict: **ACCEPT**, **REVISE**, or **REJECT**. List findings by P0/P1/P2/P3 with exact file and line, the failure each permits, and the smallest correction. State explicitly whether each finding could permit an outcome-conditioned retry or instead concerns record quality only.

If and only if the verdict is ACCEPT, produce an independent review record under `docs/reviews/` containing:

- reviewer role and date;
- the pinned candidate commit and independently computed LF-normalized digest;
- the checks performed and reconstructed custody values;
- the verdict;
- accepted evidence limitations and residual risks;
- confirmation that neither outcome gate is opened and no execution is authorized; and
- the exact next artifact: a separate payload commit that sets the four post-closure gate fields to the accepting review, followed by independent review of that payload and the required refreeze. It is not a chair gate patch and not a Gate 1A opening decision.

Do not set any post-closure gate field as part of this review. Under the frozen boundary, those fields belong to the later independently reviewed payload.
