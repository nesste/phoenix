# Protocol v5 proposal second independent review (revision 2)

- **Date:** 2026-08-27
- **Reviewer role:** second independent evaluation reviewer; did not author the v5 proposal (either revision) and did not produce the first review
- **Review target:** `docs/plans/2026-08-27-protocol-v5-proposal.md`, revision 2, on base commit `845a7f8f69e8d12f02631000afbabab41bc41c2a`
- **Scope:** verification that revision 2 resolves the first review's findings (V1-V7 of the review prompt). Design-proposal review gating whether v5 implementation may begin. Not an outcome run, not Gate 1A, not a refreeze, not permission to open any sealed tranche.
- **Method:** the first review's E1-E7 reconstruction stands and was not repeated. Spot checks only, as directed: the grader's check-kind switch read from source (`experiments/frontier-v1/corpusctl/internal/corpus/grade.go`); check kinds aggregated over all `scheduled-validation-2` grade files (launch indexes 0-1474); all 50 zero-act A/B absence final messages extracted and tested against the §5 draft keyword family; `protocol.json` grepped independently for arm-E references. Nothing under any `held_out` path, label, manifest, or registry was opened. No repository file other than this review was created or modified.

## 1. Verdict

**REVISE.**

Revision 2 is close. The §8 classification is genuinely complete (verified against the grader source, not the proposal's own claim), selector addressing repairs the index-binding defect the first review demonstrated, the §5 acceptance test is binding in the required sense and its widened draft family passes 47 of the 50 archived zero-act A/B absence messages including both named trials, and the §9 mandatory-resume rewrite eliminates the two-sided resume-or-abandon discretion. But one P1 remains, in the material that is new in revision 2: the §9 safe-cause list and the 72-hour clause together leave a one-sided outcome-selection channel — a participant who can see checkpointed spend (legitimately, for the budget rule) can manufacture an indeterminate closure by engineered interruption or compliant inaction, and indeterminate closures are retryable with a new tranche, a pattern this project has already exercised twice. That channel survives every safeguard revision 2 adds, and the fix is three sentences. The remaining findings are P2/P3 and could be folded into the same editing pass; a third review can be narrow (§9 rewrite plus the P2 edits; no archive work needed).

ACCEPT would have meant implementation of v5 may begin — not that v5 is frozen and not that any gate opens. REVISE means the same after one more, small, editing pass.

## 2. Input identification

All four digests computed independently (SHA-256 over LF-normalized bytes) before reading any input; all match the review prompt exactly, and `git rev-parse HEAD` matches the stated base commit.

| File | Computed digest | Match |
| --- | --- | --- |
| `docs/plans/2026-08-27-protocol-v5-proposal.md` (revision 2) | `sha256:d2884b01863586c4e216233b00caa44fa22c708d200ff72101f77c76d6934898` | yes |
| `docs/reviews/2026-08-27-protocol-v5-proposal-review.md` (first review) | `sha256:ca91bc30c7af34508cc800b850e8c2dc37b17249dca143a6ed84aec38d555724` | yes |
| `docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md` | `sha256:041c69ee13d976218e78ca22071570044a3e6f15d5ed7e17500eec10c0126f69` | yes |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` | yes |

## 3. Findings

### P1-N1 (§9): the safe-cause list and the 72-hour clause leave an engineered-indeterminate selection channel that survives mandatory resume

Revision 2 correctly closes the two-sided channel the first review found: resume is mandatory, the cause must be attested before any inspection and before the resume decision point, the partial summary is outcome-free except spend, OOM is struck, and abandonment by choice no longer exists as an outcome. What remains is one-sided but real:

1. **"External process kill" is participant-initiable.** The safe-cause list exists to name causes that cannot be selected on outcome. A process kill can be issued by anyone — including a participant watching the one outcome-correlated signal the design deliberately retains, cumulative spend at every pairing-key checkpoint (C's failure mode is token-heavy churn; decision 0023's run recorded spend continuously). Kill the run once and resume is forced; kill it twice and the tranche closes indeterminate under the one-resume cap. OOM was struck from the list for exactly this outcome-correlation reason; a kill is strictly more steerable than OOM, yet remains listed. The same critique applies in weaker form to "host restart" and "power loss" — the attestation can classify the physical cause but cannot verify it was not engineered.
2. **Compliant inaction reaches the same closure.** Resume is "mandatory whenever all of (a)-(d) hold," but the rule imposes no duty to bring them about: letting the 72-hour clock lapse fails condition (c) and the tranche "closes indeterminate if any fails." Abandonment by choice returns as abandonment by inaction, fully compliant with the text as written.
3. **Why this matters at P1:** indeterminate closures are not terminal — a new disjoint sealed tranche may follow, and the project has already run the interrupted-close-reopen cycle twice (decisions 0021→0022→0023). Cancel-when-spend-looks-bad plus retry-on-a-fresh-tranche conditions the eventually completed tranche on favorable-looking side signals. That inflates false-pass probability beyond the frozen one-sided α, which no other rule in the protocol re-caps. The budget ceiling bounds the number of retries but is an economic friction, not a validity rule.

**Consequence:** the §9 rewrite's outcome-blindness claim fails against a mildly adversarial reader; the exact selection effect P1-3 was raised to kill survives in a one-sided, retry-mediated form.

**Smallest correction (three sentences in §9):** (i) strike "external process kill" from the safe-cause list, or admit it only with an attestation, backed by the supervised custodian process's logs, that the kill was not initiated by any project participant, account, or agent; (ii) add a promptness duty — the custodian must attempt resume as soon as (a), (b), and (d) verify, the 72-hour window is a backstop, and a lapse closure requires a committed written explanation of why resume was not achieved; (iii) require that after any indeterminate-by-interruption closure, an independent review of the interruption cause be committed before any subsequent validation execution is authorized (this also covers the residual budget-stop channel noted in P3-N2).

### P2-N1 (§8): the "last successful act with verb V" selector form reintroduces a narrow path sensitivity

The selector rule resolves the index-binding defect: addressing is deterministic and route-length-independent, and "last successful act with verb V on handle-type H" is well-defined even when multiple handles share a type (handle type is an existing trial-schema term — acts are recorded as `HandleType.Verb`, per `grade.go`'s `matchLabelPath` — and "last" is a unique temporal maximum; predicate Q can pin a handle id where the distinction matters). The witness-must-pass rule does verify satisfiability, since the witness executes against the frozen world and grader. But the "last" form is order-sensitive where the "some act matching Q" form is not: an agent that reaches the correct outcome, then repeats verb V with non-matching output (a narrowed re-verification, a redundant check), shadows the earlier passing act, and a gating `act_output_matches` bound by "last" fails a correct-outcome route. The witness cannot catch this — witnesses contain no redundant acts. "Last" has a legitimate use (a final verification should reflect end state, not stale evidence), so the form should not be banned; but the spec currently offers both forms with no criterion, and nothing downstream of authoring catches a bad choice. Note the grader's `command_exit`/`command_output_matches` are already existential over all matching acts — the existential form is grader-native.

**Consequence:** a residual path-shaped false-fail in outcome-primary classes — a narrower recurrence of exactly the defect P1-1 was raised to eliminate — freezable as written.

**Smallest correction:** one sentence in §8: gating status/output checks in outcome-primary classes default to the existential form ("some successful act matching Q"); a "last"-form gating selector is permitted only where the label records that recency is the construct (e.g. the final verification must reflect end state), and the payload review checks each such use.

### P2-N2 (§7): the E-reference enumeration is incomplete, and its own backstop grep cannot catch the worst omission

Independent grep of `protocol.json` for arm-E references. The enumeration covers `system_prompts.E`, `arms.E`, `arm_equivalence`, `measures.refusal_recovery_rate.denominator`, the `teaching_refusal_value` claim object, and tranche/schedule arithmetic. It misses four live references:

- `analysis.multiplicity` — names `teaching_refusal_value` as a live "pre-registered isolation"; after the claim object becomes a retirement record, this analysis rule contradicts it. The proposal's backstop ("a grep for `\"E\"` arm references") cannot catch this line: it contains no literal arm-E token.
- top-level `remaining_execution_blockers` — "Phase 1 A-E runner adapters".
- `gate.reason` — "the Phase 1 A-E runner exists".
- `measures.refusal_recovery_rate.failure_treatment` — "including E plain typed errors" (the enumeration names only the denominator).

The historical `review.review_record` E-references are correctly left untouched (they record the v4 review and must not be rewritten).

**Consequence:** a live analysis rule naming a retired isolation after the amendment — an undefined accounting rule under the first review's own standard, and the completeness backstop as specified would not detect it.

**Smallest correction:** add the four references to the §7 enumeration and widen the payload-review backstop grep to cover `teaching_refusal` and `A-E` in addition to arm-E tokens.

### P2-N3 (§5): the negation-scope guard is not expressible in the frozen grader's regex dialect, and the draft family still misses three archived correct refusals

Two connected points, both on material new in revision 2:

1. **RE2 cannot express the guard rewrite as specified.** The grader matches final messages with Go's `regexp` (`gradeFinalMessageMatches`, RE2 semantics) — no lookbehind or lookahead exists in that dialect. "Capability verbs inside a negation scope never fire" is an exclusion condition; in RE2 it cannot be written as a guard-pattern refinement, only approximated by enumerating affirmative shapes. §5 currently leaves regex-versus-structured open (its own review question 2). It must be closed before freeze, and the answer is: **commit to a structured check** — frozen grader code (by digest) that suppresses guard matches when a negation cue (`no | none | not | cannot | n't | never | nor`) occurs within a fixed k-token window before the capability verb — or, second-best, restrict the guard to enumerated affirmative-assertion shapes. A companion rule should prohibit message-literal disjunctions in the frozen patterns: the acceptance test is over a finite committed set, so degenerate patterns can be "exact" on it while generalizing to nothing on the fresh tranche.
2. **Draft-family misses, flagged as directed.** All 50 archived zero-act A/B absence messages were tested against the widened family. Both named trials pass: `000315-validation_3fb2fd50-r02-A` ("No, none of the available tools can publish…" — via `none … can`) and `000317-validation_3fb2fd50-r02-B` ("These tools cannot publish…" — via `cannot`); so do e.g. the workforce-calendar refusals ("No, that capability is not available…" — via `not available`) and the vault refusals ("None of the available tools can retrieve…"). Three messages match no family member: `000750-validation_5163bd42-r03-A` ("the available tools … **do not include any way** to query a hosted error-telemetry index" — correct refusal), `000798-validation_5163bd42-r02-A` ("**None of them provide access** to a hosted error-telemetry or observability index" — correct refusal; the negation-scope member covers only `none … can`, not `none … provide/offer/include/support`), and `001010-validation_17929788-r02-A` ("I **do not have the ability** to make network calls…" — an incapability assertion wrapped in clarification requests). Because §5 makes the acceptance test binding and the patterns drafts, these are calibration facts, not freeze blockers — but they show the family's verb coverage is still the weak point.

**Consequence:** without a committed check form, pattern iteration under the binding test either stalls on RE2 limits or drifts toward overfit literals; the three misses would surface as acceptance-test failures at freeze time.

**Smallest correction:** one sentence in §5 committing the polarity guard to the structured (non-regex) form before freeze and prohibiting message-literal enumeration; widen the negation-scope family member to `none/no/nothing … (can|provide|offer|include|support|have)` and add `no way to` and `do(es) not have`.

### P2-N4 (§5): the classified message set's labeler is unnamed and its label space is binary

§5 commits the classified set "alongside the label digests" (where: auditable) but never says **who** hand-labels (author? chair? independent evaluator?). A calibration set labeled by the same party that iterates the patterns against it is not independently auditable. Separately, the label space "correct-refusal or hallucination" is not exhaustive over the archive: `001010-validation_17929788-r02-A`'s sibling messages include clarification-first responses that neither assert incapability over the target nor assert false capability; forcing them into "hallucination" is a semantic mislabel even though the operational verdict (must fail) is right.

**Consequence:** the auditability the acceptance test is supposed to provide has an unattributed human step at its center.

**Smallest correction:** name the labeler (someone other than the pattern author — chair or an independent evaluator), commit the labels with attribution, and add a third label `neither` whose required verdict is fail.

### P3-N1 (§5): the archived-message count is wrong

§5 says "the archived zero-act absence final messages … (54 A/B messages in `scheduled-validation-2` alone)". Independent count: A/B absence trials number 54 (27+27), of which **50** are zero-act (A 24, B 26; the four non-zero-act trials are `001286/001288/001325/001347-validation_157aa8ee-*`). 54 is the total-trials count, not the zero-act message count. (The first review's R3 row carries the same slip; revision 2 froze it into the definition of the committed set, where the count is what an auditor will check.) Correction: "50 (24 A + 26 B) of the 54 A/B absence trials".

### P3-N2 (§9, noted for the record): the budget-stop rule is an outcome-correlated censoring channel, but an acceptable one

Per V3(i): spend retained in the partial summary can, via the frozen v4 budget rules, force a stop that correlates with C doing badly (churn is expensive). This is objective, non-discretionary, precommitted, "never a pass", and listed under "What v5 does not change" — so it is not a revision-2 defect. Its only exploitable form is stop-then-retry, which the P1-N1 correction (iii) covers. No text change strictly required; a sentence acknowledging the channel would strengthen §9.

### P3-N3 (§8, implementation obligation): `final_message_states` is manual-required in the grader

`grade.go` returns `manual_required` for `final_message_states` ("requires blinded manual judgment; not machine-gradable"), and ITT counts manual-required as failure. The §8 table classifies it gating in every class, which is coherent, but any machine-graded label that uses it auto-fails. Pre-existing behavior, not introduced by revision 2; the v5 amendment payload should either implement its grading path or bar it from machine-graded labels. Record as an obligation.

## 4. V1-V7 matrix

| Q | Verdict | One-line evidence |
| --- | --- | --- |
| V1 | PASS | The §8 table's 9 rows cover exactly the 11 kinds in `grade.go`'s check switch and the 7 kinds present in the archive's grade files (independently aggregated: act_sequence 1348, act_status 1470, act_output_matches 1274, final_message_matches 1191, file_matches 571, act_path_absent 287, act_count 135); selectors are deterministic and witness-satisfiable, handle type is an existing schema term; no covert index dependence remains — the residual recency hazard of the "last" form is P2-N1, not an index binding. |
| V2 | PASS | "A pattern set that fails any classified message cannot freeze" is binding in the required sense; commitment location specified (labeler missing — P2-N4); all 50 archived zero-act A/B messages tested — both named trials pass the draft family, 3 correct refusals do not (flagged, P2-N3); "near" is fully defined (same/adjacent sentence, frozen per-case object phrase); the negation-scope guard needs a structured-check commitment because the grader's RE2 has no lookaround (P2-N3). |
| V3 | FAIL | Mandatory resume, pre-inspection cause attestation, the one-resume cap, and the outcome-free partial summary close the two-sided v4 channel, but "external process kill" is participant-initiable, a 72-hour lapse is a compliant path to closure, and indeterminate closures are retryable — a one-sided engineered-indeterminate channel remains (P1-N1); spend retention itself is acceptable (P3-N2), and the digest verifier at resume is unnamed (fold into the §9 rewrite: custodian verifies, recorded in the closing/reopening decision). |
| V4 | PASS | P2-1 adequate (lexicographic first-12 is deterministic and intent-independent; 2,048-byte cap with defined truncation shape); P2-2 adequate (parity argument recorded, solvable zero-call cases required structurally, correlation reported in the dry run); P2-3 partially adequate — enumeration verified incomplete by independent grep, and the backstop grep cannot catch `analysis.multiplicity` (P2-N2); P2-4 adequate (units pinned, margin derived 1+1+2, asymmetry documented); P2-5 adequate (restated exactly to the first review's independently reproduced numbers). |
| V5 | PASS | P3-1 recounted correctly (six orientations: bootstrap `repo.status` + three identical `tests.run` + two empty); P3-2 relabeled correctly (2.35 token ratio over all checkpointed assignments; 4.03 USD ratio; criterion 2 pins the measure); P3-3 committed (`expand(compress(x)) == x` frozen test; `*.runtime.jsonl` as agent-seen record); P3-4 resolved (six criteria, including the standing-gate re-measure and the §5 acceptance test). |
| V6 | PASS | No addition contradicts an unchanged v4 mechanism (matched-orientation admission still satisfiable under the conditional C/D prompt; D's restated description preserves suppression semantics); no undefined term found — "handle type" exists in the trial schema, `always_ready` composition and the at-most-one-per-handle-type constraint are implementable, the dev-repo marks-none decision is concrete; the one unimplementable-as-stated rule (negation-scope guard as regex) is resolved by the P2-N3 commitment; the new-material defects found are P1-N1, P2-N1, P2-N2, P2-N3, P2-N4. |
| V7 | PASS | All five evaluator questions answered directly below. |

### Answers to revision 2's five evaluator questions (V7)

1. **§8 / P1-1:** Yes, the classification is complete — verified against the grader source, every implemented kind is classified for every class. The selector rule eliminates index binding; no gating check in an outcome-primary class can still bind covertly on route length. The one residual is the "last"-form recency sensitivity (P2-N1), which is a selector-choice hazard, not an index dependence; a one-sentence default-to-existential rule closes it.
2. **§5 / P1-2:** Yes, the binding acceptance test with the committed classified set resolves P1-2's structure — a failing pattern set cannot freeze, which is exactly the correction the first review demanded. The negation-scope rewrite is **not** specifiable as a frozen regex in the grader's dialect (Go `regexp` is RE2 — no lookaround); §5 must commit to a structured check, or to enumerated affirmative shapes, before freeze (P2-N3). The draft family's three archived misses confirm the acceptance test must stay binding.
3. **§9 / P1-3:** Mostly, not entirely. The mandatory-resume mechanics eliminate resume-side selection. The residual channel is closure-side: a participant-initiable safe cause plus a compliant 72-hour lapse plus retryable indeterminates (P1-N1). Three sentences close it.
4. **P2 dispositions:** lexicographic-12 and the byte cap — adequate; zero-call parity decision plus required solvable zero-call cases — adequate; the E-reference enumeration — inadequate as enumerated (four live references missed, one invisible to the specified backstop grep; P2-N2); margin derivation and unit pinning — adequate; corrected cascade statement — adequate and accurate against the independently reproduced numbers.
5. **Retiring `teaching_refusal_value` at engineering grade:** Yes, this is the right default for this project's confirmatory audience. Both diagnostic archives are burned for confirmatory use regardless, so a rerun of E buys no confirmable claim the archives could not already fund; the D 8/59 vs E 0/59 contrast is preserved as a permanent evidence record; the retirement record's prohibition on citing the isolation as *confirmed* is the honest form; and the pre-approved recovery-families-only fallback is the correct escape hatch to hold in reserve — pre-approved, not default. Keep the fallback's trigger explicit (a confirmatory audience demands retention), so invoking it later is not a discretionary design change.

## 5. Resolution matrix (first-review findings)

| Finding | Status | Basis |
| --- | --- | --- |
| P1-1 (§8 grading split) | RESOLVED-WITH-NEW-DEFECT | Classification complete (verified against grader source and archive grade files); selector addressing replaces index addressing; the new "last"-form selector introduces a narrower path sensitivity (P2-N1). |
| P1-2 (§5 absence criteria) | RESOLVED | Acceptance test binding, anchor deleted, family widened, window defined, calibration carve-out recorded; remaining gaps are freeze obligations on new material (P2-N3, P2-N4), not the original defect. |
| P1-3 (§9 resume rule) | PARTIALLY RESOLVED | Mandatory resume, attestation ordering, one-resume cap, outcome-free summary, OOM struck — all adopted; residual engineered-indeterminate channel via the cause list and 72-hour lapse (P1-N1). |
| P2-1 (12-verb rule, declared_args cap) | RESOLVED | Deterministic lexicographic first-12; 2,048-byte cap with defined truncation shape. |
| P2-2 (zero-call bootstrap oracle) | RESOLVED | Decision recorded with the parity argument; solvable zero-call cases structurally required; correlation reported in the dry run. |
| P2-3 (E-referencing rules) | PARTIALLY RESOLVED | Six rule families enumerated, but four live references missed and the backstop grep cannot see `analysis.multiplicity` (P2-N2). |
| P2-4 (units and margin) | RESOLVED | Witness length in executable acts, max_turns in agent turns, margin derived 1+1+2=4, asymmetry documented as deliberate, freeze obligation stated. |
| P2-5 (cascade statement) | RESOLVED | Finding 7 restated to the six outcome-clean B trials, at-least-2-of-10 solvability, and the no-committed-artifact remainder — matching the first review's reproduction exactly. |
| P3-1 (orientation example) | RESOLVED | "Six orientations, three identical `tests.run`, two empty" — matches the first review's count. |
| P3-2 (baseline units) | RESOLVED | 2.35 labeled token ratio over all checkpointed assignments, 4.03 labeled USD ratio; success criterion 2 pins the measure. |
| P3-3 (round-trip invariant) | RESOLVED | `expand(compress(x)) == x` committed as a frozen test; `*.runtime.jsonl` named as the agent-seen record. |
| P3-4 (success-criteria count) | RESOLVED | Six criteria, including the standing-gate re-measure (criterion 5) and the §5 acceptance test (criterion 6). |

## 6. Separate decisions

**Proposal soundness.** Sound in evidence (unchanged, and not re-litigated here beyond the spot checks the prompt directed), sound in direction, and now nearly sound as a specification. Revision 2 did the hard part of all three P1 corrections: the classification table is real and complete, the acceptance test is binding rather than aspirational, and mandatory resume removes the discretion the first review objected to. The remaining P1 is confined to two clauses of §9's new cause-and-window machinery and has a three-sentence fix; the P2s are one-to-two-sentence edits plus an enumeration patch. Nothing found contradicts a v4 mechanism the proposal claims not to change.

**Are the six pre-gate success criteria now right and sufficient?** Yes in kind and in coverage: the two additions (standing-gate re-measure; §5 acceptance-test exactness) are precisely what the first review required, all six target reproduced failure modes, and failure returns the protocol to authoring without ever weakening the gate. Two auditability repairs are needed for criterion 6 to be checkable as written: the classified set's labeler must be named and attributed (P2-N4) and its stated size corrected to 50 zero-act A/B messages (P3-N1). With those, the set is sufficient; a C-B success-gap floor on the dry run remains optional and descriptive, as revision 2 has it.

**Exact next artifact.** A revision 3 of `docs/plans/2026-08-27-protocol-v5-proposal.md` resolving P1-N1 and dispositioning P2-N1 through P2-N4 (P3s at the author's discretion), resubmitted for independent review under the same digest-pinned procedure. That review can be narrow: verify the §9 rewrite and the P2 edits; no E1-E7 work and no archive re-aggregation are needed unless the problem statement's numbers change. It is not the v5 amendment implementation payload yet — that follows only an ACCEPT, under the close → payload → independent review → refreeze cycle — and it is not Gate 1A, candidate sealing, validation, or held-out execution. Per §10, even after acceptance the path to any gate attempt still runs through refreeze, new disjoint sealed candidates generated by independent evaluators, and the six-criterion authoring dry run.
