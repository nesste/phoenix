# Protocol v5 proposal independent review prompt

Copy the text below into a fresh evaluator session that did not author the v5 proposal and did not produce the 2026-08-27 exploratory report.

```text
Act as the independent evaluation reviewer for the Phoenix protocol-v5 proposal. This is a design-proposal review, one stage earlier than a candidate-commit review: the proposal is an uncommitted draft on top of commit 845a7f8f69e8d12f02631000afbabab41bc41c2a and no v5 implementation exists. Your verdict gates whether implementation of v5 may begin. It is not an outcome run, a corpus-authoring task, Gate 1A, a refreeze, or permission to open any sealed tranche.

Repository working tree: D:\Work\personal\phoenix (base commit 845a7f8f69e8d12f02631000afbabab41bc41c2a)

Identify the review inputs before reviewing. Compute SHA-256 over LF-normalized bytes and require these exact values:

- docs/plans/2026-08-27-protocol-v5-proposal.md: sha256:dfb72d4c6813818c7d6f67f1b77db7db2d7f46dc05d6cbc0c85eb28f09ed8c0e
- docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md: sha256:041c69ee13d976218e78ca22071570044a3e6f15d5ed7e17500eec10c0126f69
- docs/reports/2026-08-27-protocol-v4-second-interrupted-validation-exploratory.md: sha256:067315a52138f0e575cc579db16d3018338ca3c5a3f087a186103dceb50cb036
- experiments/frontier-v1/protocol.json: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9

Return REVISE without reviewing if any digest differs.

Review context

Decision 0023 closed the second interrupted Gate 1A execution as indeterminate. Both sealed validation tranches (the 720-run archive at experiments/frontier-v1/results/scheduled-validation and the 1,475-launch archive at experiments/frontier-v1/results/scheduled-validation-2) are inspected diagnostic archives; they cannot support any future gate claim. The v5 proposal claims transcript-level root causes for the C-B deficit and proposes ten changes. Do not defer to the proposer's analysis: reconstruct the evidentiary claims yourself from the scheduled-validation-2 archive before judging the changes that rest on them.

Review scope

Read from the working tree:

- docs/plans/2026-08-27-protocol-v5-proposal.md (the candidate)
- experiments/frontier-v1/protocol.json (the v4 baseline it amends)
- docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md and 0013, 0021, 0022 for the interruption and repair history
- docs/reports/2026-08-27-protocol-v4-second-interrupted-validation-exploratory.md and its .json companion
- experiments/frontier-v1/manifests/validation-label-digests.json (class mapping)
- The scheduled-validation-2 archive: assignment, trial, grade, and runtime files as needed for the evidence checks below

Independence and custody rules

- Do not modify any repository file except writing your review to docs/reviews/2026-08-27-protocol-v5-proposal-review.md.
- Do not run authoring, validation, or held_out trials.
- The two burned validation archives are open diagnostic evidence under decision 0023 and MAY be read. Do not open or inspect anything under corpus/held_out, fixtures/held_out, or held_out labels, manifests, or registries.
- Do not generate, relabel, reseal, or replace any candidate.
- Read-only inspection and script-based aggregation over the diagnostic archives are allowed and expected.

Evidence checks (reconstruct independently; do not trust the proposal's numbers)

E1. Orientation churn: over the C-arm turn_limit assignments in launch indexes 0-1474, compute the distribution of orientation calls (tool inputs containing "intent") versus executable acts per trial, and the same for successful C trials. The proposal claims means near 7.0 vs 1.7 orientations respectively.
E2. Post-bootstrap orientation behavior: in at least three C runtime transcripts, verify that consecutive post-act orientation responses return the identical frontier call regardless of changed intent text.
E3. Error dead ends: verify that "absent" error payloads name no reachable verbs and "invalid_arguments" payloads include no schema, in actual archived responses.
E4. Absence grading: extract the final_message_matches patterns from absence-class grade files and at least two archived final messages; decide independently whether the proposal's "correct refusal fails on phrasing" reading is accurate, and whether the class's 0% across all five arms is fully explained by grading plus the bootstrap execute-mandate.
E5. Cascade solvability: verify that no arm passed any cascade pairing in the prefix, and assess from grades and acts whether a within-12-turn path plausibly exists (the proposal claims no committed artifact proves one).
E6. Recovery advantage: verify C/D versus A/B recovery success rates and that at least one successful C recovery trial shows refusal.instead steering.
E7. Arm E: verify E's overall and recovery-class success counts used to justify dropping the arm.

Required review questions

R1. Orientation retirement (proposal §1): Does making post-bootstrap orient return an idempotent exhausted response, combined with the conditional system prompt, close the churn mechanism without destroying anything the v4 amendment relies on (pending-frontier precedence, omitted-state restoration from the current pending frontier, refusal-alternative delivery)? Can an agent still recover a pending refusal alternative after the change, and by what path?
R2. Discovery leakage (proposal §2): Does listing up to 12 reachable verbs on absent errors and declared argument schemas on invalid_arguments errors leak capability beyond the constant standing surface? Consider the standing_surface engineering gate at 10/100/1000 verbs, adversarial-text cases where the world plants misleading text, and temptation cases where naming a reachable-but-wrong verb may steer the agent into the trap. Decide whether the 12-verb cap is load-bearing and whether error payload size is bounded.
R3. Zero-call bootstrap and the absence class (proposal §§1, 5): With bootstrap allowed to return zero calls and the execute-mandate removed, is the absence-class construct still measured (correctly refusing), or has the task been made trivially passable? Are the proposed semantic grading criteria (keyword family near target noun phrase, negative pattern) well-defined enough to freeze, and is validating them against burned-archive messages a legitimate use of inspected data or does it overfit grading to two spent samples?
R4. Path witnesses and the cap rule (proposal §§3, 6): Is a scripted non-model witness a sound solvability criterion? Does max_turns = max witness length + 4, uniform across arms, remain arm-fair given C spends a turn on bootstrap orientation that A/B do not? State whether the +4 margin needs a principled derivation before freeze.
R5. Dropping arm E (proposal §7): v4 pre-registered teaching_refusal_value as an isolation. Is retiring it on diagnostic evidence from burned tranches methodologically acceptable, or does it constitute outcome-informed protocol modification that a skeptical reader would flag? If retained, what is the minimal form?
R6. Grading altitude split (proposal §8): For each of the eight classes, decide whether the gating/descriptive assignment of act_sequence and act_count preserves the class construct. Pay attention to cascade (moved to outcome-primary) and stale_frontier (kept path-primary).
R7. Resume rule (proposal §9): Is the precommitted resume rule (durable pairing-key checkpoint, custodian non-inspection attestation, digest-identical bytes, 72-hour window, outcome-uncorrelated cause) adequately outcome-blind? Identify any channel through which an interrupted-then-resumed run could be selected for favorable outcomes, including the decision to resume versus abandon.
R8. Envelope compression (proposal §4): Does wire-form compression with archival re-expansion create any divergence between what the agent saw and what the grader or auditor reconstructs? Is the 8-char world_build prefix adequate for custody continuity?
R9. Evidence integrity: Do your independent reconstructions E1-E7 match the proposal's claims? Flag any claim the archive does not support at the stated magnitude.
R10. Completeness: Does the proposal leave any v4 mechanism in an undefined state (activation rules marked always_ready, frontier score semantics, D's post-act refusal-alternative delivery with frontiers suppressed, measures that reference E)? An undefined cross-arm rule, accounting rule, or safety invariant is a defect; a documented future implementation obligation is not.

Required response

Write your full review to docs/reviews/2026-08-27-protocol-v5-proposal-review.md containing:

1. Verdict: ACCEPT, REVISE, or REJECT — where ACCEPT means implementation of v5 may begin, not that v5 is frozen or any gate opens.
2. Input identification: the four independently computed normalized digests.
3. Findings ordered P0 to P3 with exact proposal section, consequence, and smallest correction; say No findings if none.
4. A PASS/FAIL/UNRESOLVED matrix for R1-R10 with one-line evidence each.
5. Your independently computed E1-E7 results beside the proposal's claimed values.
6. Answers to the proposal's own five evaluator questions.
7. Separate decisions on: proposal soundness, whether the five pre-gate success criteria are the right ones, and what the exact next artifact is. The next artifact must not be Gate 1A, candidate sealing, validation, or held_out execution.

Also return the verdict, the P0/P1 findings, and the R-matrix in your final message.
```
