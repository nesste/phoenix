# Protocol v5 proposal second independent review prompt

Copy the text below into a fresh evaluator session that did not author the v5 proposal (either revision) and did not produce the first review.

```text
Act as the second independent evaluation reviewer for the Phoenix protocol-v5 proposal, revision 2. The first review (docs/reviews/2026-08-27-protocol-v5-proposal-review.md) returned REVISE with three P1 findings (P1-1 §8 grading split, P1-2 §5 absence criteria, P1-3 §9 resume rule), five P2 and four P3 findings, and independently reproduced all seven evidence checks E1-E7 at stated magnitude. Revision 2 claims to resolve the P1 findings and disposition every P2/P3. Your job is to verify the resolutions, not to re-litigate what the first review already settled. This is a design-proposal review gating whether v5 implementation may begin. It is not an outcome run, a corpus-authoring task, Gate 1A, a refreeze, or permission to open any sealed tranche.

Repository working tree: D:\Work\personal\phoenix (base commit 845a7f8f69e8d12f02631000afbabab41bc41c2a)

Identify the review inputs before reviewing. Compute SHA-256 over LF-normalized bytes and require these exact values:

- docs/plans/2026-08-27-protocol-v5-proposal.md (revision 2, the candidate): sha256:d2884b01863586c4e216233b00caa44fa22c708d200ff72101f77c76d6934898
- docs/reviews/2026-08-27-protocol-v5-proposal-review.md (first review): sha256:ca91bc30c7af34508cc800b850e8c2dc37b17249dca143a6ed84aec38d555724
- docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md: sha256:041c69ee13d976218e78ca22071570044a3e6f15d5ed7e17500eec10c0126f69
- experiments/frontier-v1/protocol.json: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9

Return REVISE without reviewing if any digest differs.

Standing of prior evidence

Do not re-reconstruct E1-E7 unless revision 2's problem-statement numbers differ from the first review's section 5, in which case verify only the differing numbers. The archive at experiments/frontier-v1/results/scheduled-validation-2 (launch indexes 0-1474; class labels in experiments/frontier-v1/manifests/validation-label-digests.json) remains open diagnostic evidence and MAY be read for spot checks. Nothing under any held_out path, label, manifest, or registry may be opened.

Independence and custody rules

- Do not modify any repository file except writing your review to docs/reviews/2026-08-27-protocol-v5-proposal-review-2.md.
- Do not run authoring, validation, or held_out trials.
- Do not generate, relabel, reseal, or replace any candidate.
- Read-only inspection and script-based aggregation over the diagnostic archives are allowed.
- Do not treat the first review's REVISE-with-no-P0 as presumptive acceptance of revision 2; verify each resolution on its own terms.

Required review questions

V1. P1-1 resolution (§8): Does the classification table cover every check kind the grader implements for every class it appears in? Cross-check the table's row set against the check kinds actually present in scheduled-validation-2 grade files (act_sequence, act_count, final_message_matches, act_status, act_output_matches, file_matches, act_path_absent) and the grader's additionally implemented kinds named by the first review (file_absent, command_exit, command_output_matches, final_message_states). Is selector addressing specified tightly enough to implement — in particular, is "the last successful act with verb V on handle-type H" well-defined when multiple handles share a type, and does the witness-must-pass-selectors rule actually verify satisfiability? Can a gating check in an outcome-primary class still bind covertly on path through any residual index dependence?

V2. P1-2 resolution (§5): Is the acceptance test binding in the sense that a pattern set failing any classified message cannot freeze? Is the classified message set's commitment auditable (who labels, where committed)? Take at least five archived zero-act absence final messages (including 000315-validation_3fb2fd50-r02-A "No, none of the available tools can publish…" and 000317-…-B "These tools cannot publish…") and decide whether the widened draft family plus the negation-scope guard rewrite plausibly classifies them correctly; flag any archived correct refusal the draft family still cannot match. Is "near the target noun phrase" now fully defined (window, source of the object phrase)? Decide whether a frozen regex can express the negation-scope requirement or whether §5 must commit to a structured (non-regex) check before freeze.

V3. P1-3 resolution (§9): With resume mandatory, one resume per tranche, cause attestation required before inspection AND before the resume decision point, OOM struck, and the partial summary specified outcome-free — does any residual outcome-selection channel remain? Scrutinize: (i) spend retained in the partial summary against the budget-stop rule — can spend force a "budget stop" that functions as outcome-correlated abandonment? (ii) the 72-hour window — can letting the clock run out serve as de-facto discretionary abandonment, and does the rule need a "resume as soon as conditions verify" clause? (iii) the cause list — is "external process kill" abusable? (iv) who verifies digest identity at resume and where is it recorded?

V4. P2 dispositions: For each of P2-1 (lexicographic 12-verb rule, 2,048-byte declared_args cap with truncation shape), P2-2 (zero-call parity decision plus required solvable zero-call cases plus dry-run correlation report), P2-3 (enumerated E-reference rewrite — grep protocol.json yourself for E references and check the enumeration for completeness), P2-4 (units pinned, margin derived as 1+1+2), P2-5 (corrected cascade statement in finding 7): adequate, inadequate, or introduces a new defect?

V5. P3 dispositions: P3-1 (example recounted: six orientations, three identical tests.run), P3-2 (units relabeled: 2.35 token ratio all assignments, 4.03 USD ratio; success criterion 2 pins the measure), P3-3 (round-trip invariant plus runtime.jsonl as agent-seen record), P3-4 (six success criteria now, including the standing-gate re-measure and the §5 acceptance test): resolved?

V6. New-defect sweep: Revision 2 adds material absent from revision 1 — the §1 zero-call parity argument and always_ready composition rule, the §3 at-most-one-always_ready-per-handle-type constraint and dev-repo marks-none decision, the §7 teaching_refusal_value retirement record and pre-approved minimal-E fallback, the §8 table itself, the §9 cause list, the §5 draft family. Does any addition contradict a v4 mechanism this proposal claims not to change, leave a term undefined, or create a rule that cannot be implemented as frozen bytes? An undefined cross-arm rule, accounting rule, or safety invariant is a defect; a documented future implementation obligation is not.

V7. Answer revision 2's five evaluator questions (its § Review questions for the second independent evaluator) directly.

Required response

Write your full review to docs/reviews/2026-08-27-protocol-v5-proposal-review-2.md containing:

1. Verdict: ACCEPT, REVISE, or REJECT — where ACCEPT means implementation of v5 may begin, not that v5 is frozen or any gate opens.
2. Input identification: the four independently computed normalized digests.
3. Findings ordered P0 to P3 with exact proposal section, consequence, and smallest correction; say No findings if none.
4. A PASS/FAIL/UNRESOLVED matrix for V1-V7 with one-line evidence each.
5. A resolution matrix: for each first-review finding (P1-1..P1-3, P2-1..P2-5, P3-1..P3-4), RESOLVED / PARTIALLY RESOLVED / UNRESOLVED / RESOLVED-WITH-NEW-DEFECT.
6. Separate decisions on: proposal soundness, whether the six pre-gate success criteria are now right and sufficient, and the exact next artifact. If ACCEPT, the next artifact is the v5 amendment implementation payload under the close → payload → independent review → refreeze cycle; it must not be Gate 1A, candidate sealing, validation, or held_out execution.

Also return the verdict, any P0/P1 findings in full, the V-matrix, and the resolution matrix in your final message.
```
