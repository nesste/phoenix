# Protocol v5 proposal fourth independent review prompt (minimal)

Copy the text below into a fresh evaluator session that did not author any revision of the v5 proposal and did not produce any of the three prior reviews.

```text
Act as the fourth independent evaluation reviewer for the Phoenix protocol-v5 proposal, revision 4. This is a MINIMAL review, verifiable by reading alone: three prior reviews (REVISE → REVISE → REVISE, the third with no P0/P1 findings) settled the evidence base and every section except the revision-4 edits. Revision 4 applies the third review's P2-1..P2-3 corrections and its P3 observations: frozen tokenization/clause-boundary/negator/scope definitions for the §5 polarity guard plus fail-side and pre-labeling clauses; §8 `selector_mode`/`recency_rationale` schema fields and a restored payload-review obligation; §9 process-event log declared as a new payload-created artifact and the post-closure gate widened to all non-completed executions; §7 historical-record exemption. Your job: verify these edits resolve the third review's findings and introduce no new defect. Everything else is settled and OUT OF SCOPE. This gates whether v5 implementation may begin. It is not an outcome run, Gate 1A, a refreeze, or permission to open any sealed tranche.

Repository working tree: D:\Work\personal\phoenix (base commit 845a7f8f69e8d12f02631000afbabab41bc41c2a)

Identify the review inputs before reviewing. Compute SHA-256 over LF-normalized bytes and require these exact values:

- docs/plans/2026-08-27-protocol-v5-proposal.md (revision 4, the candidate): sha256:1f977edfec3cce7ed1c53571cb73712edcaba53eab47486a4122acd6ee5518e5
- docs/reviews/2026-08-27-protocol-v5-proposal-review-3.md (third review): sha256:6d7759d773c67c2e3a504e7984023db6fe1210d5399bdee42d6195a41b1d9e6d

Return REVISE without reviewing if either digest differs.

Custody and independence rules

- Do not modify any repository file except writing your review to docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md.
- Do not run authoring, validation, or held_out trials; never open anything under any held_out path, label, manifest, or registry.
- Reading alone should suffice. Permitted spot checks if needed: the label/case schema files under experiments/frontier-v1/schema/ (to confirm the §8 field addition is coherent with the closed schema), and the trial file 001024-validation_*-r01-A (final_message field) in experiments/frontier-v1/results/scheduled-validation-2 for question M1.

Required review questions

M1. §5: Are tokenization ("splits on Unicode whitespace and punctuation, lowercases, keeps contraction suffixes"), the clause-boundary list, and the negator list (now including n't and nor) implementable as frozen grader code without further interpretation? Does the false-capability scope rule — unsuppressed capability hit in the same/adjacent sentence as an object-phrase token AND asserting capability over the goal's target, with genuinely-available-action offers out of scope and the classified set as arbiter — resolve the third review's 001024-r01-A false-fail concern without smuggling in a new undefined judgment ("asserts capability over the goal's target" — is the classified-set arbiter clause sufficient to make this freezable)?
M2. §5: Do the synthetic fail-side exemplar clause and the illustrative-only pre-labeling clause resolve the third review's P3 observations on vacuous fail coverage and prose pre-labeling?
M3. §8: Do the selector_mode ("some" default | "last") and required recency_rationale fields, stated as an explicit change to the closed label schema, plus "the payload's independent review checks every last-mode use against its rationale", close P2-2?
M4. §9: Does declaring the process-event log a new payload-created artifact — with stated content (start/stop/checkpoint/termination events, timestamps, originating principal where the host exposes it) and digest commitment to the custody record — close P2-3? Does widening the post-closure gate to "interruption, budget stop, or lapse" close the budget-stop-then-retry gap without contradicting the frozen v4 budget-exhaustion rule (which it leaves as the closure mechanism, adding only the audit-before-retry gate)?
M5. New-defect sweep over the revision-4 edits only: any new undefined term, nonexistent artifact presented as existing, or contradiction with settled sections (including the §7 historical-record exemption's consistency with the retirement rules)?

Required response

Write your full review to docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md containing:

1. Verdict: ACCEPT, REVISE, or REJECT — where ACCEPT means implementation of v5 may begin, not that v5 is frozen or any gate opens.
2. Input identification: the two independently computed normalized digests.
3. Findings ordered P0 to P3 with exact proposal section, consequence, and smallest correction; say No findings if none.
4. A PASS/FAIL/UNRESOLVED matrix for M1-M5 with one-line evidence each.
5. A resolution matrix for the third review's findings (P2-1, P2-2, P2-3, and its P3 observations): RESOLVED / PARTIALLY RESOLVED / UNRESOLVED / RESOLVED-WITH-NEW-DEFECT.
6. If ACCEPT: state the exact next artifact — the v5 amendment implementation payload under the close → payload → independent review → refreeze cycle, followed by independent generation of new disjoint sealed candidates and the authoring dry run against the six pre-gate success criteria; it must not be Gate 1A, candidate sealing, validation, or held_out execution. If ACCEPT, also produce the acceptance record: reviewer role, date, the two digests, accepted limitations, and remaining freeze obligations carried into the payload review.

Also return the verdict, any findings, the M-matrix, and the resolution matrix in your final message.
```
