# Protocol v5 proposal independent review

- **Date:** 2026-08-27
- **Reviewer role:** independent evaluation reviewer; did not author the v5 proposal or the 2026-08-27 exploratory report
- **Review target:** `docs/plans/2026-08-27-protocol-v5-proposal.md` (uncommitted draft on base commit `845a7f8f69e8d12f02631000afbabab41bc41c2a`)
- **Scope:** design-proposal review gating whether v5 implementation may begin. Not an outcome run, not Gate 1A, not a refreeze, not permission to open any sealed tranche.
- **Method:** all evidence checks E1-E7 were reconstructed with fresh scripts over the `experiments/frontier-v1/results/scheduled-validation-2` diagnostic archive (launch indexes 0-1474 only; class labels from `experiments/frontier-v1/manifests/validation-label-digests.json`). Nothing under any `held_out` path, label, manifest, or registry was opened. No repository file other than this review was created or modified.

## 1. Verdict

**REVISE.**

The transcript-level evidence base is solid: every one of the proposal's seven mechanism claims reproduces from the archive at the stated magnitude (section 5), and the direction of all ten changes is supported. Implementation may not begin yet because three P1 specification defects would be baked into the implementation as written: the §8 gating/descriptive split is incoherent for the grader's position-indexed checks and silently omits the `act_status` check kind; the §5 semantic grading criteria fail the proposal's own acceptance test against the very archive messages it cites; and the §9 resume rule leaves the resume-versus-abandon decision discretionary, which is a live outcome-selection channel. Each has a small correction. A revised draft resolving the P1 findings should return for independent review; no re-reconstruction of E1-E7 will be needed unless the evidence section changes.

ACCEPT here would have meant implementation of v5 may begin — not that v5 is frozen and not that any gate opens. REVISE means the same after one more editing pass.

## 2. Input identification

All four digests computed independently (SHA-256 over LF-normalized bytes) and match the review prompt exactly:

| File | Computed digest | Match |
| --- | --- | --- |
| `docs/plans/2026-08-27-protocol-v5-proposal.md` | `sha256:dfb72d4c6813818c7d6f67f1b77db7db2d7f46dc05d6cbc0c85eb28f09ed8c0e` | yes |
| `docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md` | `sha256:041c69ee13d976218e78ca22071570044a3e6f15d5ed7e17500eec10c0126f69` | yes |
| `docs/reports/2026-08-27-protocol-v4-second-interrupted-validation-exploratory.md` | `sha256:067315a52138f0e575cc579db16d3018338ca3c5a3f087a186103dceb50cb036` | yes |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` | yes |

## 3. Findings

### P1-1 (§8): the gating/descriptive split is undefined for `act_status` and incoherent for position-indexed gating checks

The validation grade files contain seven check kinds: `act_sequence`, `act_count`, `final_message_matches`, `act_status`, `act_output_matches`, `file_matches`, `act_path_absent` (the grader additionally implements `file_absent`, `command_exit`, `command_output_matches`, `final_message_states`). §8 classifies only five. `act_status` appears in every direct, cascade, far-discovery, stale-frontier, and temptation case and in recovery twice per case — and is never assigned to either side of the split.

Worse, the checks §8 designates as outcome-gating are position-indexed in the current grader: `act_status` and `act_output_matches` verdicts read "act 4 has expected status ok" / "act 4 output satisfies pattern". If `act_sequence` becomes descriptive, an agent may reach the correct outcome by a different-length route, and the gating check then evaluates a different act at index 4 and fails. Outcome-primary classes would still gate on path — covertly, through check addressing — defeating the change's purpose.

**Consequence:** §8 as written cannot be implemented coherently; the counterfactual it is justified by (launch 0 arm A: artifact and verifier passed, sequence failed — verified) would not actually be repaired for routes of different lengths.

**Smallest correction:** classify every implemented check kind as gating or descriptive per class, and redefine the gating act-addressed checks to be selector-addressed (e.g. "some successful act with verb `tests.run` whose output matches P", or "the last act matching predicate Q") instead of sequence-index-addressed. Freeze the addressing rule with the labels, as §8 already promises for the classification.

### P1-2 (§5): the enumerated semantic criteria fail the proposal's own acceptance test against the archived messages

§5 requires that "every semantically correct refusal in those archives must pass the new checks." Reconstruction shows the enumerated keyword family (`cannot | can not | no tool | not available | not possible | unable`) does not contain any member matching archived semantically correct refusals such as `000315-validation_3fb2fd50-r02-A` (zero acts): "No, none of the available tools can publish a checklist to a remote release system." — "none … can" matches no family member.

Additionally, the existing negative pattern that §5 retains with only the word "reviewed" is demonstrably defective: its raw pattern `\b(?:can|supports?|allows?)\b|…` fires on negated capability phrasing, so the same correct refusal fails the polarity guard too (both `final_message_matches` checks fail on that trial; A absence trials record 42 message-check failures across 27 trials). The guard currently penalizes exactly the refusal phrasings the class wants.

**Consequence:** frozen as enumerated, v5 absence grading would again fail correct refusals — repeating validity defect 6 with new patterns.

**Smallest correction:** make the acceptance criterion the binding specification: the family and negative pattern are drafts to be revised until every archived semantically correct refusal passes and every capability hallucination fails, with the classified message set committed alongside the review §5 already requires. Handle negated-capability constructions ("none … can", "only covers", "has no access", "doesn't accept") explicitly, and define "near the target noun phrase" (window size, noun-phrase source) before freeze.

### P1-3 (§9): the resume rule is not outcome-blind because resuming is optional

The rule says an interrupted execution "may resume" iff (a)-(d) hold. Nothing requires it. The custodian and chair, while blind to outcome artifacts, can observe outcome-correlated side signals — recorded spend, checkpoint progress rate, wall time per launch (C's failure mode is token-heavy churn, so high spend correlates with C doing badly). A discretionary resume lets an interrupted-then-losing run be abandoned and an interrupted-then-winning run be resumed, which is selection on outcomes through a side channel. Two secondary gaps: (i) OOM is listed as an acceptable interruption cause, but OOM is plausibly outcome-correlated for this protocol (churny failing trials produce the longest transcripts); (ii) the content of `scheduled-summary.partial.json` is unspecified — if it carries outcome tallies, writing/reading it at checkpoints undermines the non-inspection attestation.

**Consequence:** the precommitted resume rule, whose whole purpose is outcome-blindness, admits a selection effect a skeptical reader will reject.

**Smallest correction:** make resume mandatory whenever (a)-(d) hold, with failure of any condition closing the tranche as indeterminate (never abandoned by choice); cap resumes per tranche (one); specify that the partial summary contains no outcome-bearing fields (launch index, pairing key, digests, timestamps only); require the interruption-cause classification to be attested before any outcome inspection, and either strike OOM from the safe-cause list or require independent attestation that the OOM was not workload-selective.

### P2-1 (§2): the 12-verb selection rule is undefined and `declared_args` is unbounded

When a handle has more than 12 reachable verbs, §2 does not say which 12 the `absent` payload lists. Any intent- or attempt-informed selection would be a new guidance channel; any unstated order is a nondeterminism in frozen behavior. `error.declared_args` has no stated size bound. Correction: a deterministic, intent-independent selection (e.g. lexicographic or registration order), frozen; a stated byte cap on `declared_args`.

### P2-2 (§§1, 3, 5): the zero-call bootstrap response can become an absence oracle for C

Once §3 guarantees every sealed solvable case a matching, witness-covered frontier, "no ready call matches this intent" at bootstrap correlates strongly with the absence label — an inference shortcut C receives and A/B do not, partially converting the absence construct for C from "agent judges incapability" to "surface signals incapability." The parity argument (B infers absence from its complete manual) is available but is not made or decided anywhere. Correction: record the design decision explicitly; require some sealed solvable cases whose bootstrap legitimately returns zero calls, so zero-call is not label-equivalent; report the zero-call/absence correlation in the §"success criteria" dry run.

### P2-3 (§7): dropping E leaves E-referencing accounting rules undefined

`protocol.json` `measures.refusal_recovery_rate` has denominator "all assigned recovery trials in D and E"; `arm_equivalence.irreducible_differences`, `system_prompts.E`, and the `teaching_refusal_value` claim object all reference E. §7 restates the isolation but does not enumerate these. Under this review's standard an undefined accounting rule is a defect. Correction: the v5 amendment payload enumerates and rewrites every E-referencing rule (refusal_recovery_rate denominator becomes D-only).

### P2-4 (§6): witness-length units and turn semantics are unpinned; the +4 margin is underived

Witnesses are act sequences (executable acts); `max_turns` counts agent turns; C's bootstrap orientation consumes a turn A/B never spend, so C's effective margin is +3 versus A/B's +4. The uniform cap is the right fairness choice (per-arm caps would break arm equivalence, and C's relief is correctly assigned to §§1-2), but the proposal itself concedes the margin lacks a derivation. Correction: define both units in the frozen text and derive the margin (≥ 1 bootstrap turn + stated slack) before freeze. This is a freeze obligation, not an implementation blocker.

### P2-5 (Problem statement, finding 7): "manual-equipped B fails outcome checks on the rest" is false for 6 of B's 30 cascade assignments

Independent reconstruction: 6 B cascade trials — `000001/000028/000069-validation_ca447353-r0x-B` and `000902/000941/000974-validation_9a329c7f-r0x-B` — passed every non-path check (`file_matches`, `act_status`, `act_output_matches`) in 7-8 acts, well inside the cap, and failed only `act_sequence`. So cascade is demonstrably outcome-solvable within the cap for at least 2 of 10 cases, and for those cases the binding defect was the sequence label, not solvability. "Cascade may be unsolvable under the cap" over-generalizes. This strengthens §3 (path witnesses would have caught the label/route mismatch) and §8 (outcome-primary cascade), but the finding must be restated accurately. The narrower claim — no committed artifact proves a within-cap path — is confirmed: validation labels live with the external custodian, and `experiments/frontier-v1/labels/` contains only `authoring/`.

### P3-1 (Problem statement, mechanism 1): the named example is over-counted

`000024-validation_864449e1-r01-C` contains 6 orientations: bootstrap returns `repo.status` (provenance `activation`, score 0); three post-act orientations return the identical `tests.run` call under differing intents; two return null after the frontier is consumed. "Four orientations returning the same `tests.run` call" should read three. The mechanism itself is confirmed (section 5, E2).

### P3-2 (Success criteria): baseline units are mislabeled

"C/B ITT token ratio below 2.0 (was 2.35 on graded trials, 4.03 overall)": 2.346 is the token ratio over **all checkpointed assignments** (9,111,034 / 3,883,734), and 4.03 is the **USD cost** ratio (18.9993 / 4.7110), not a token ratio. Correct the parenthetical and pin the measurement definition (ITT token totals over dry-run assignments).

### P3-3 (§4): compression needs a committed round-trip invariant

The archival re-expansion means the committed trial bytes are not the bytes the agent saw. This is recoverable — `*.runtime.jsonl` records the raw `tool_result` content delivered to the agent — but the design should state: (i) expansion is deterministic and lossless with a committed `expand(compress(x)) == x` test, and (ii) runtime transcripts remain the committed record of agent-seen bytes. The 8-char `wb` prefix is adequate given the full digest appears in the first response and in every archival trial record.

### P3-4: the review prompt says "five pre-gate success criteria"; the proposal lists four

Turn-limit rate, token ratio, zero mandated-bootstrap absence failures, universal path witnesses — four. Either a criterion was dropped from the draft or the count is wrong; reconcile (section 7 recommends what a fifth should be).

## 4. Review-question matrix

| Q | Verdict | One-line evidence |
| --- | --- | --- |
| R1 | PASS | Refusal alternatives already arrive inline on the refused act response in both C and D archives (e.g. `000092-…-D`: `tests.focus` refused with `refusal.instead` = `tests.list`, `frontier: []`), so retiring post-act orient closes churn without severing delivery; recovery path after ignoring `instead` is retrying the act (re-delivery) and omitted-state restoration lives in admission, untouched by §1. |
| R2 | PASS | §2 exposes only currently-reachable verbs on the addressed handle in a per-act payload; the standing-surface gate measures standing input, which §2 does not touch, and the 12-cap bounds the payload at 1000 verbs — but the selection rule is undefined (P2-1) and the promised gate re-measurement must stay binding. |
| R3 | FAIL | The enumerated keyword family misses archived correct refusals it is required to pass and the retained negative pattern misfires on them (P1-2); archive validation itself is legitimate (54 zero-act A/B absence messages, a grading-side use §10 explicitly carves out), but the criteria are not freezable as written. |
| R4 | PASS | A scripted witness is a sound and needed solvability criterion (P2-5 shows sequence labels, not solvability, bound two cascade cases); uniform cap is the right fairness call; the +4 margin needs the derivation and unit-pinning of P2-4 before freeze. |
| R5 | PASS | Dropping E is outcome-informed but rests on burned diagnostic tranches that can never support confirmatory claims anyway, leaves the headline C-B estimator untouched, and is restated (not silently deleted) as engineering-grade; the D-8/E-0 recovery contrast reproduces exactly; accounting cleanup required (P2-3). |
| R6 | FAIL | `act_status` is unclassified and the gating checks are position-indexed (P1-1); cascade→outcome-primary is supported by the six outcome-clean B trials; absence/temptation/adversarial-text keep `act_count`/`act_path_absent` gating and stale-frontier stays path-primary, preserving those constructs. |
| R7 | FAIL | "May resume" leaves resume-versus-abandon discretionary while spend and progress rate are outcome-correlated side signals, and OOM as a safe cause is itself plausibly outcome-correlated (P1-3). |
| R8 | PASS | Grader inputs keep archival shape and `*.runtime.jsonl` preserves the agent-seen wire bytes, so no unrecoverable divergence; requires the committed round-trip invariant of P3-3; 8-char prefix adequate for continuity. |
| R9 | PASS | All seven evidence checks reproduce at stated magnitude (section 5); three secondary misstatements flagged (P2-5, P3-1, P3-2), none overturning a mechanism. |
| R10 | FAIL | `refusal_recovery_rate`'s D-and-E denominator, `arm_equivalence`, and `system_prompts.E` are left undefined rather than documented as obligations (P2-3); the `always_ready` flag's interaction with §3's intent-match predicates is unstated; D's arm text ("after an executable act only a pending refusal alternative may be returned") needs restatement under §1. |

## 5. Independent E1-E7 reconstruction

All computations restricted to launch indexes 0-1474 (1,475 assignments, 295 per arm, verified). Arm totals reproduce the exploratory report exactly (A 103, B 122, C 56, D 32, E 24 successes; C 73 turn-limit + 1 agent error).

| Check | Proposal's claim | Independently computed | Verdict |
| --- | --- | --- | --- |
| E1 orientation churn | 73 C turn-limit trials average 7.0 orientations vs 7.3 executable acts; successful C average 1.7 | n=73 turn-limit: mean 6.96 orientations, 7.27 executable acts (median 7); n=56 successes: mean 1.71 orientations, 2.68 acts; all-C mean 4.32 | Reproduces |
| E2 post-bootstrap replay | Post-act orient ignores intent and replays the pending frontier; example trial shows 4 identical `tests.run` returns | Named trial `000024-…-C`: consecutive post-act orientations with different intent text both return the identical `tests.run` call; strict scan finds the consecutive-differing-intent-identical-frontier pattern in 20/73 turn-limit C trials; example has 3 (not 4) `tests.run` returns among 6 orientations | Reproduces; example over-counted (P3-1) |
| E3 error dead ends | `absent` names no reachable verbs; `invalid_arguments` shows no schema | 352/352 `absent` payloads are exactly `{"code":"absent","message":"act is not reachable"}`; 48/48 `invalid_arguments` are exactly the bare generic message; no payload anywhere names a verb or schema | Reproduces exactly |
| E4 absence grading | 0% all arms; anchored pattern fails correct refusals (`000317-…-B`); 27/27 C/D/E failures include `act_count` from the mandated bootstrap act | 0/135 absence successes; `000317-…-B`: zero acts, final message "These tools cannot publish anything to a remote release system. …", `act_sequence`/`act_count` pass, fails only the start-and-end-anchored pattern; C/D/E: 27/27 each executed `repo.status` as first act (bootstrap frontier verified: provenance `activation`, score 0) and failed `act_count` max 0; A 24/27 and B 26/27 zero-act trials failed on message checks alone. Additional: the polarity guard also fails correct refusals phrased "none … can" (P1-2) | Reproduces; grading defect worse than claimed |
| E5 cascade solvability | No arm passed any cascade pairing; C turn-limits 23/30; B fails outcome checks on the rest; no committed artifact proves a within-cap path | 0/150 cascade successes; C: 23 turn-limit + 7 graded fails; B: 28 graded fails + 2 turn-limit — but 6 of B's graded fails passed **all** outcome checks in 7-8 acts, failing only `act_sequence` (2 of 10 cases outcome-solvable within cap); no committed artifact proves a path (validation labels external) | Reproduces except "B fails outcome checks on the rest" (P2-5) |
| E6 recovery advantage | Recovery: C/D 13.6% vs A/B 5.1%; `refusal.instead` steers directly to the fix | C 8/59 and D 8/59 (13.6%) vs A 3/59 and B 3/59 (5.1%), E 0/59; verified in `000094-…-C` and `000179-…-C`: `tests.focus` refused with `instead` = `tests.list` → agent executes `tests.list` → `tests.focus` ok → graded pass | Reproduces exactly |
| E7 arm E | E succeeded on 8.1% with zero recovery successes; D 8 recovery successes | E 24/295 (8.13%), all 24 in direct; recovery 0/59; D recovery 8/59 | Reproduces exactly |

Supporting ratios verified: C/B token ratio 2.346 (9,111,034/3,883,734), C/B cost ratio 4.033 (18.9993/4.7110 USD), cost per success C 0.339 vs B 0.039 USD.

## 6. Answers to the proposal's five evaluator questions

1. **Does §2 leak capability beyond the constant standing surface, and does the 12-verb cap bound the standing-surface gate at 1000 verbs?** No leak in principle: the standing surface (constant act-tool schema) is untouched, and the disclosure is per-act, reachability-scoped, and earned by acting — consistent with the thesis. The gate measures standing input tokens, which §2 does not change, so the cap is not what protects the gate; it is load-bearing for bounding the error payload at 1000 verbs, and boundedness fails without it. Two conditions: freeze a deterministic, intent-independent 12-verb selection rule (P2-1) and keep the promised gate re-measurement binding. On temptation: the trap verb, being reachable, will be listed — that is the construct under test (B's manual already names it, and B scores 44.8% on temptation), not a leak; note it in the frozen rationale.
2. **Is the §9 resume rule adequately outcome-blind?** No, as written (P1-3). The custodian attestation covers artifact inspection but not the decision channel: "may resume" lets outcome-correlated side signals (spend rate, progress, wall time) drive resume-versus-abandon. The 72-hour window is fine; the discretion is not. Make resume mandatory when the conditions hold, cap it at one, specify the partial summary as outcome-free, and reconsider OOM as a safe cause.
3. **Is dropping E acceptable?** Yes (R5), because both tranches are burned for confirmatory use regardless, the headline estimator is unaffected, and the isolation is restated at engineering grade rather than silently deleted — provided every E-referencing accounting rule is rewritten (P2-3) and no future artifact cites teaching-refusal value as confirmed. If a skeptical audience demands retention, the minimal form is E run only on the 8 recovery families (~20% of E's cost), since the isolation's endpoint is recovery.
4. **Does §8's split preserve the temptation and stale-frontier constructs?** Yes for both: temptation keeps `act_sequence`, `act_count`, and `act_path_absent` gating (resisting the trap is the path), and stale-frontier stays path-primary (not following the stale frontier is the outcome; archived C 15/33 successes come from exactly that behavior). But the split as a whole is only coherent after P1-1: `act_status` must be classified, and gating act-addressed checks must stop being sequence-index-addressed.
5. **Is the witness-plus-4 cap arm-fair given C's bootstrap turn?** Acceptably so, with an obligation: uniform caps preserve arm equivalence, and C's relief correctly comes from §§1-2; C's effective margin is one turn thinner (+3 after bootstrap) than A/B's (+4). The margin needs a written derivation before freeze — at minimum `margin >= 1 bootstrap turn + slack`, with witness length defined in executable acts and `max_turns` defined in agent turns (P2-4). Successful C trials averaged 1.71 orientations + 2.68 acts, so the margin is empirically comfortable; it just is not yet principled.

## 7. Separate decisions

**Proposal soundness.** Evidence-sound and direction-sound; specification-incomplete. All five protocol-surface mechanisms and both validity defects reproduce independently, several (E3, E4) more starkly than claimed, and one (E5) slightly overclaimed in a direction that actually favors the proposal's own §§3 and 8. The repairs target the reproduced mechanisms without touching the frozen measurement machinery. The defects are in §5, §8, and §9 as specifications — each correctable in one editing pass. REVISE, not REJECT.

**Are the pre-gate success criteria the right ones?** The proposal lists four, not five (P3-4). The four are right in kind (they target the reproduced failure modes and can only return the protocol to authoring, never weaken the gate) but incomplete. Add: (i) the §2 standing-surface gate re-measurement as an explicit criterion — it is promised in §2 but absent from the criteria list; (ii) the §5 acceptance test as a criterion — every archived semantically correct absence refusal passes the frozen checks and every capability hallucination fails; and fix the mislabeled baselines (P3-2) so "token ratio below 2.0" has a defined measurement. A C-B success-gap floor on the authoring dry run would also be reasonable and is legitimately outcome-visible there.

**Exact next artifact.** A revised v5 proposal draft (successor to `docs/plans/2026-08-27-protocol-v5-proposal.md`) resolving P1-1, P1-2, and P1-3 and dispositioning the P2 findings, resubmitted for independent review under this same digest-pinned procedure. Not Gate 1A, not candidate sealing, not validation, not held-out execution — and per §10, even after acceptance the path runs through the full close → payload → independent review → refreeze cycle, new disjoint sealed tranches generated by independent evaluators, and the §"success criteria" authoring dry run before any gate attempt.
