# Protocol v5 proposal third independent review prompt (narrow)

Copy the text below into a fresh evaluator session that did not author any revision of the v5 proposal and did not produce the first or second review.

```text
Act as the third independent evaluation reviewer for the Phoenix protocol-v5 proposal, revision 3. This is a NARROW review: revisions 1 and 2 were reviewed in full (first review: REVISE, three P1s, all resolved or dispositioned; second review: REVISE, one P1 and four P2s and three P3s against the new material). Revision 3 edits only §5, §7, §8, §9, the revision header, the evaluator-questions section, and success criterion 6. Sections not edited in revision 3 are settled by the prior reviews and are OUT OF SCOPE except where a revision-3 edit interacts with them. Your job: verify that revision 3 resolves P1-N1 and adequately dispositions P2-N1..P2-N4 and P3-N1..P3-N3 from the second review, and that the edits introduce no new defect. This gates whether v5 implementation may begin. It is not an outcome run, Gate 1A, a refreeze, or permission to open any sealed tranche.

Repository working tree: D:\Work\personal\phoenix (base commit 845a7f8f69e8d12f02631000afbabab41bc41c2a)

Identify the review inputs before reviewing. Compute SHA-256 over LF-normalized bytes and require these exact values:

- docs/plans/2026-08-27-protocol-v5-proposal.md (revision 3, the candidate): sha256:d396f5ffe37409f049bde5a00e96f1d1b10febeeb0c61317dfb6c290a6c51754
- docs/reviews/2026-08-27-protocol-v5-proposal-review.md (first review): sha256:ca91bc30c7af34508cc800b850e8c2dc37b17249dca143a6ed84aec38d555724
- docs/reviews/2026-08-27-protocol-v5-proposal-review-2.md (second review): sha256:0105d6eb62811135698cd86fb92c456ea9b74441663c2e6fdf3fe8fc54fe700a
- experiments/frontier-v1/protocol.json: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9

Return REVISE without reviewing if any digest differs.

Custody and independence rules

- Do not modify any repository file except writing your review to docs/reviews/2026-08-27-protocol-v5-proposal-review-3.md.
- Do not run authoring, validation, or held_out trials; never open anything under any held_out path, label, manifest, or registry.
- Do not re-reconstruct E1-E7 or re-review out-of-scope sections. Spot checks against experiments/frontier-v1/results/scheduled-validation-2 (launch indexes 0-1474) and protocol.json are allowed where the questions below direct.

Required review questions (these mirror the candidate's own five third-evaluator questions; answer both sets as one)

N1. §9 / P1-N1: Is the one-sided selection channel closed? Evaluate: (i) process kill now requires custodian-log establishment that no participant, account, or agent initiated it — is "establishes" verifiable in practice, and does the rule fail safe (no establishment → indeterminate)? (ii) the promptness duty ("resume as soon as (a), (b), (d) verify; 72h backstop; lapse closure requires committed written explanation") — does abandonment by inaction survive in any form? (iii) the post-closure gate (independent review of the interruption-cause record before any subsequent validation execution) — does it bound the engineered-indeterminate retry cycle, and is it enforceable given the project's decision-gated workflow? (iv) any NEW channel the rewrite itself opens.
N2. §8 / P2-N1: Does the existential default selector eliminate the shadowing hazard? Is the recency carve-out ("only where the label records recency as the construct") tight enough to freeze, or is "records recency as the construct" itself undefined? Does the P3-N3 disposition (final_message_states: implement or bar) leave grading coherent?
N3. §7 / P2-N2: Grep protocol.json yourself. Is the enumeration now complete — including analysis.multiplicity, remaining_execution_blockers, gate.reason, and refusal_recovery_rate.failure_treatment — and does the stated verification rule (read every arm-letter occurrence AND every retired claim id) close the letter-grep blindspot?
N4. §5 / P2-N3, P2-N4, P3-N1: Is the structured polarity guard (tokenize; suppress capability-verb hits in the scope of a preceding same-clause negator) specified tightly enough to implement as frozen grader code — in particular, is "same clause" defined or definable at freeze? Do the three-label scheme, chair-labels-evaluator-countersigns attribution, and corrected count (50: A 24, B 26) satisfy the audit requirements? Spot-check the widened family against the three previously missed refusals (000750-…-A "do not include any way", 000798-…-A "None of them provide access", 001010-…-A "do not have the ability") and at least two clarification-style messages that should be labeled neither.
N5. New-defect sweep over the revision-3 diffs only: does any edit contradict an unchanged v4 mechanism, leave a term undefined, or create a rule that cannot be implemented as frozen bytes? Check specifically: the §9 "custodian's process log" (does such a log exist as a committed artifact in the current custodian design, or is it a new obligation the payload must create?); the §5 negator list and prohibition on message-literal patterns; the §7 "A-D runner" rewrites; success criterion 6's three-label failure semantics.

Required response

Write your full review to docs/reviews/2026-08-27-protocol-v5-proposal-review-3.md containing:

1. Verdict: ACCEPT, REVISE, or REJECT — where ACCEPT means implementation of v5 may begin, not that v5 is frozen or any gate opens.
2. Input identification: the four independently computed normalized digests.
3. Findings ordered P0 to P3 with exact proposal section, consequence, and smallest correction; say No findings if none.
4. A PASS/FAIL/UNRESOLVED matrix for N1-N5 with one-line evidence each.
5. A resolution matrix for the second review's findings: P1-N1, P2-N1..P2-N4, P3-N1..P3-N3 — RESOLVED / PARTIALLY RESOLVED / UNRESOLVED / RESOLVED-WITH-NEW-DEFECT.
6. Separate decisions on: proposal soundness across all three reviews, whether the six pre-gate success criteria stand, and the exact next artifact. If ACCEPT, the next artifact is the v5 amendment implementation payload under the close → payload → independent review → refreeze cycle; it must not be Gate 1A, candidate sealing, validation, or held_out execution.

Also return the verdict, any P0/P1 findings in full, the N-matrix, and the resolution matrix in your final message.
```
