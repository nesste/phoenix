# 0024: Protocol-v5 amendment payload import and refreeze

- **Status:** v5 amendment payload (proposal §§1, 2, 4, 7, 8) imported and frozen; validation closed; held-out closed
- **Date:** 2026-08-27
- **Decision:** `IMPORT_AND_FREEZE`
- **Authorizing proposal:** `docs/plans/2026-08-27-protocol-v5-proposal.md` (revision 4), fourth independent review `ACCEPT` at `docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md`
- **Closure in force:** `docs/decisions/0023-protocol-v4-gate-1a-second-interrupted-execution.md` — both outcome gates closed throughout
- **Accepted payload:** `0f8d9c72c59bea5da5abf792f493f1d75299b1f8` on base `24895c4a272167cf8ea1445d3b81b2133ae62c4c`
- **Candidate metadata:** `experiments/frontier-v1/artifacts/gate-1a-protocol-v5-amendment-candidate.json` (62 files, set digest `sha256:fcea7dcd02fcc90165a11004e92b9233bbe6da213fbd86a6671ba5785b65874e`)
- **Independent review:** `docs/reviews/2026-08-27-protocol-v5-amendment-payload-review.md` at `bd096c6d89f6a9a4cd70fb49fa4f16cdc574c227`; verdict `ACCEPT`, no P0/P1/P2 findings
- **Refreeze commit:** `1f80d7b`

## Decision

The project chair accepts the independent review and refreezes the protocol-v5 amendment payload. Both outcome gates remain closed throughout this import.

The payload implements the five mechanisms of the accepted v5 proposal that are implementable before a new sealed tranche exists: bootstrap-only orientation in fact (the standalone post-act `orient` returns `exhausted`; zero-call bootstrap with `always_ready`-gated fallback; conditional C/D prompts), reachability-scoped discovery on `absent` and `invalid_arguments` errors, wire-envelope compression with a frozen lossless round-trip invariant, the arm-E removal rewrite of `protocol.json` and the A–D runner, schedule, scheduletool, and analysis machinery, and selector-addressed outcome-primary grading with the frozen gating/descriptive classification and the `final_message_states` bar. Proposal §§3, 5, 6, and 9 are later v5 amendment work registered in `protocol.json` `remaining_execution_blockers`.

The replacement identities are:

| Artifact | Identity |
| --- | --- |
| World-build digest (linux/amd64, `version=dev`, `bin/phoenix`) | `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` |
| World-build manifest, LF-normalized | `sha256:04508e1c3b54df5e725e6753590ab039d2e6f2469b08216cea7128724e0cd070` |
| Grader semantics | `sha256:7b7438f84164d69157bbde87fd3ffb2f88e069bb01a97785a737910e29df0e05` |
| `protocol.json` (v5), LF-normalized | `sha256:39dfcd2f0554a97657e7140da041b7ea76bebdc9e3dc8f63cc88732b30534d69` |
| Scheduled-runner boundary, 21 files, set digest | `sha256:5fe8632a1a7eb37b412f280703e72f113fe65e34ed0a4624113df0a948d22a24` |
| Analysis implementation, 9 files, set digest | `sha256:824c152f6c3835b79471a15fa9c6463c2db68cade11ed399f470cc04bfc71223` |
| Local candidate (prompts, invocation, world build, grader), LF-normalized | `sha256:55b212d95f90c4bb39461277ed16db993d09868804dc13d60d8e5384ce3cb37b` |

The world definition is unchanged: raw `sha256:efed8654…3fbae96`, canonical `sha256:d7f93051…89bc`. The world-build digest was reproduced byte-identically — executable and manifest — on the WSL2 linux/amd64 execution host with the frozen recipe at the frozen `bin/phoenix` path before this import.

## Retirement

All v4-sealed unopened candidates are retired: the five-arm validation schedule (`sha256:b38a0eaa…5813`), its manifest, public label-digest registry, and private-archive identities are preserved as custody evidence, marked retired in `pre-validation-artifacts.json`, and refused by the deterministic A–D schedule construction. They must not be run, relabeled, or resealed. The retired identities remain pinned in `validation_grader.go` as custody evidence until a new disjoint sealed v5 tranche replaces them; validation execution is impossible in this state by construction. Both burned diagnostic archives remain diagnostic-only under decisions 0013 and 0023.

The `teaching_refusal_value` isolation is retired at engineering grade per the v5 proposal; no future artifact may cite teaching-refusal value as confirmed, and the archived D/E contrast is its permanent evidence record.

## Accepted residuals

- Carried P2-1: `classifyScheduledOutputEntries` measures gocyclo 15, exactly at the threshold with zero headroom.
- Carried P2-3: `applyScheduledResume` returns `done=true` with a non-nil error on stop paths; callers must check the error first.
- Review P3-1 (pre-existing): `validateExternalGrade` lets a gating `manual_required` dominate a gating `fail`, the reverse of the grader's `overallStatus` precedence; unreachable under v5 because no check kind produces `manual_required`. The optional alignment may ride a later payload.
- Review P3-2 is resolved in the refreeze commit: the README and execution-boundary world-build literals were updated to the refrozen digest, so those two set members' digests differ from the payload commit's bytes; the accepted-findings note in `pre-validation-artifacts.json` records this, and the payload-commit bytes remain verifiable through the commit-pinned candidate guard.
- The zero-call-bootstrap/absence correlation is unbroken until the new sealed tranche includes solvable zero-call cases; the authoring dry run must report it.

## Refreeze boundary

The focused refreeze commit changes: `pre-validation-artifacts.json` (re-pinned scheduled-runner, analysis, runtime-invocation, arm-A, world-build, and grader blocks; retired schedule marking; both gates retained closed), `pre-validation-local-candidate.json` (conditional C/D prompts, A–D server template, new grader and world-build identities), `world-build.linux-amd64.json` (regenerated), the two frozen documents' world-build literals, the freeze-test expectations, and the candidate guard's conversion to a commit-pinned historical check. The accepted payload and review records were already committed. The full suite, staticcheck, gocyclo, dupl, validate-spec, validate-authoring, and the local freeze-candidate comparison pass against this refreeze. No model, trial, grade, output directory, cost, validation result, or held-out result was produced or observed.

## Authorization boundary

This decision does not reopen Gate 1A and authorizes no execution. Before any Gate 1A attempt: the remaining v5 amendment work (§5 absence-grading acceptance test with its chair-labeled, evaluator-countersigned classified message set; §3 path witnesses; §6 witness-derived uniform turn cap frozen with the schedule; §9 mandatory-resume rule with its process-event log and post-closure gate), independent generation and sealing of new disjoint v5 validation and held-out tranches with passing path witnesses, the six-criterion authoring dry run, updated Task 0.6 acceptance, and a separate chair opening decision are all required. Held-out remains closed.
