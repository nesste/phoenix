# 0025: Protocol-v5 §5 absence-grading acceptance test import and refreeze

- **Status:** §5 absence-acceptance payload imported and frozen; validation closed; held-out closed
- **Date:** 2026-08-27
- **Decision:** `IMPORT_AND_FREEZE`
- **Authorizing registration:** protocol v5 (decision 0024) and its accepted proposal's §5, discharging the fourth proposal review's freeze obligations 1 and 2
- **Accepted payload:** `54e256e9f4347d34844a3ef9b2156600360dcd13` on base `6e0f35afe6128a971369a710256691f6cc8de311`
- **Candidate metadata:** `experiments/frontier-v1/artifacts/gate-1a-absence-acceptance-candidate.json` (18 files, set digest `sha256:46e3bc40cf05c7160a4341958bf0f4c817d320bb4861f8885af350e1e129dc90`)
- **Independent review and countersignature:** `docs/reviews/2026-08-27-absence-acceptance-payload-review.md` at `7ef1faf7b9b8e41207e201a0ebc36fdc457cd7e0`; verdict `ACCEPT`, no P0/P1/P2, no disputed label
- **Refreeze commit:** `e07a3fc`

## Decision

The project chair accepts the independent review — which carries the §5 countersignature — and refreezes the absence-grading acceptance payload. Both outcome gates remain closed throughout this import.

Absence grading now uses the v5 semantic criteria: zero executable acts, a final message asserting incapability over the goal's target, and no false-capability assertion, graded by the structured `final_message_refuses_target` check with a frozen per-case `object_phrase`. The frozen guard implements §5's tokenization, clause boundaries, negator suppression, and widened incapability family, with both fourth-review P3-1 corrections and the reviewed construction-class extensions; no message-literal pattern exists. The v4 anchored sentence-shape requirement — the defect that zeroed the absence class in both burned tranches — is gone.

**The classified message set is the binding specification.** The 50 archived zero-act absence-class A/B final messages of the burned `scheduled-validation-2` archive (A 24, B 26; extracted and independently re-derived byte-for-byte under the §10 calibration carve-out) plus 4 synthetic `capability_hallucination` and 3 synthetic `neither` exemplars are labeled by the project chair (Raoul Bivolaru, 2026-08-27) and countersigned by the payload's independent reviewer, who read all 57 messages and concurred with every label and all nine object phrases with no disputes. The binding acceptance test `TestAbsenceAcceptanceIsExactOnTheClassifiedMessageSet` is exact on the whole set, and the countersigned set is digest-pinned in `pre-validation-artifacts.json` (`absence_acceptance_classified_set`, LF `sha256:3db22411f087ded7459fbc5c27d9bf0c137ab726025522d9de0f9ba8d9a5e35e`). A mechanization that misclassifies any labeled message cannot freeze; any future guard revision goes through the full frozen-byte cycle with the set as arbiter.

The replacement identities are: grader `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b`; `protocol.json` (v5 with the `absence_grading` amendment record) LF `sha256:f46cfc86311ccea6ff491ae5621f170358ef73d15795d0b779ed137588a93f47`; scheduled-runner 21-file set digest `sha256:dcb7256bfd647211642997f9fc0841f2abf7a0c7f5130df21ac5288d1254cc62`; local candidate LF `sha256:a70ce35b7cfbb686cf7922c99394f6c3ef344235b03f8cf3f88bccb8ce892967`. The world definition, world build (`sha256:425bab1c…dae4`), analysis set, and Arm B document are unchanged from decision 0024. The authoring absence label is migrated to the v5 shape (zero acts, structured check) and all eight authoring labels pin the new grader digest with the manifest regenerated.

## Accepted residuals

Carried: gocyclo-15 zero headroom on `classifyScheduledOutputEntries`; `applyScheduledResume` error-first contract; the pre-existing external-grade manual/fail precedence divergence (unreachable under v5). New from the §5 review, both arm-blind and bounded by the set-as-arbiter rule: the interrogative-inversion false-fail channel on fresh interrogative refusals, and conditional hallucinations (an if-clause, whether conditioning on consent or on a fact) evading the false-capability branch while still failing on the missing incapability assertion. The behavior of the frozen proxy on fresh-tranche messages never labeled remains disciplined, not proven, and criterion 6 exactness must not be cited as two-sided calibration (carried from the proposal reviews). corpusctl sits outside the repository quality gate's gocyclo directory set, consistent with its pre-existing code.

## Refreeze boundary

The focused refreeze commit changes: the classified set's countersignature fields; `pre-validation-artifacts.json` (scheduled-runner and local-artifact blocks re-pinned to `54e256e`, grader identity updated, the new `absence_acceptance_classified_set` block); the freeze-test expectations including the new classified-set verification; and the candidate guard's conversion to a commit-pinned historical check. The accepted payload and countersigned review were already committed. The full suite, staticcheck, gocyclo, dupl, and format gates pass. No model, trial, grade, output directory, cost, validation result, or held-out result was produced or observed.

## Authorization boundary

This decision does not reopen Gate 1A and authorizes no execution. Remaining before any Gate 1A attempt, per `protocol.json` `remaining_execution_blockers`: the outstanding v5 amendment work (§3 path witnesses, §6 witness-derived uniform turn cap frozen with the schedule, §9 mandatory-resume rule with its process-event log and post-closure gate), independent generation and sealing of new disjoint v5 validation and held-out tranches with passing path witnesses, the six-criterion authoring dry run (whose criterion 6 — this acceptance test's exactness on the committed set — is now discharged in the frozen code and re-verified at the dry run), updated Task 0.6 acceptance, and a separate chair opening decision. Held-out remains closed.
