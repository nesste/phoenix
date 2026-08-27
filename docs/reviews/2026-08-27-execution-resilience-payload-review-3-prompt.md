# Third independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 3

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** either of the two reviewers who produced the previous review records. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a third review

Two independent reviews have run against earlier commits of this payload, and both returned REVISE:

| Review | Target | P0 | P1 | P2 | P3 |
| --- | --- | --- | --- | --- | --- |
| 1 — `docs/reviews/2026-08-27-execution-resilience-payload-review.md` | `7a1ea42…` | 0 | 2 | 6 | 9 |
| 2 — `docs/reviews/2026-08-27-execution-resilience-payload-review-2.md` | `e92eaac…` | 0 | 0 | 2 | 10 |

Review 2 confirmed both of review 1's P1s genuinely fixed and all six P2s fixed or substantively fixed. **Both of its own new P2s were in the fix code, and both refused a *legitimate* resume rather than permitting an invalid one — one of them caused by an interaction between two separate fixes.** That is the failure mode this round should watch hardest: a repair that trades a permissive defect for a restrictive one, or that is correct in isolation and wrong in combination.

Your scope is the **whole payload**, not the deltas. Read both prior records, but treat their verdicts as claims, not settled fact. If a prior finding was wrong, or a residual either reviewer accepted should not have been accepted, say so.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`60c6aee55b71716743b30d084a57daeae48b3373`**. Ancestry: parent `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` (second review record) → `e92eaacf6ab54b547b331d14749717da64b58e0d` (superseded revision 2) → `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` (first review record) → `7a1ea42513acd3d55e276eabc2459da4a037acc0` (superseded payload 1) → `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree. It is intentionally dirty with exactly one untracked file, this prompt. Anything else dirty is a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:9b581f396d2ed8f8c285141ff90d613b6033fbf9555c7249c414eb7a2d0078e8`, 16 files, ordinal-path set digest `sha256:473024fb57ca62a3be4ed032e05a2e37835c76ed60e10d1f72292559e1776147`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:41b9e227b094cecb331699a60ece587e162cadfdb9dae55a53426de3c403bbee`.
- LF-normalized SHA-256: `runner/execution_resilience.go` `sha256:fe98996728d7ccbd4b1bbd23f44d29e3ebb543a64e2430502a7ccbfb2dee374a`; `runner/process_events.go` `sha256:605a365310aebb4f98c5e75fca8ae96232e671a48381756557f7c891ec70c617`; `runner/validation_gate.go` `sha256:0bdaa343872454c8d0387bbf0de3b4cc2e7a3b20f9dcecd218f4c62d1adafe80`; `runner/execution_resilience_test.go` `sha256:9cf56873fb0bdcc0ad229c101b3d919b631ddb1ef0520241a72311833de34957`; `validation-execution-boundary.md` `sha256:a3c32fa33645a4e287d7476f8fad3e260451b7baea6845e79501c40dfc6b6d12`; `README.md` `sha256:fc748a2ffab7e4cf4256edb232bd5682802835391f2704e845f11082f07ed45b`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256, rendered `sha256:<hex>`. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## Governing sources

- `docs/plans/2026-08-27-protocol-v5-proposal.md` § 9 — the accepted specification.
- `docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md`, **freeze obligation 5**.
- The two prior review records named above.
- `experiments/frontier-v1/artifacts/execution-resilience-candidate.md` — the implementer's account, with "What the first review changed" and "What the second review changed" sections. Claims to verify.

## What you must check

**A. Both directions on every check, always.** For each § 9 condition, ask two questions, not one: can an *invalid* resume pass, and can a *valid* resume be refused? Give both answers explicitly. Concretely:

1. **The timing anchors** (`verifyResumeTimingAnchors`). The interruption anchor now allows `resumeClockSkew`. Is five minutes right — enough for an honest second-precision attestation of a signalled interruption, and not so much that it opens anything? Walk a realistic timeline for each of the four frozen causes (`host_restart`, `power_loss`, `hardware_failure`, `process_kill`) from interruption to resume, and say whether an honest custodian passes. Check the interaction between the skew allowance, the 72-hour measurement from `lastEvent`, the `resumeClockTolerance` freshness bound, and the promptness/lapse threshold — these four constants now interact and were added at different times.
2. **The closure verdict parser** (`verifyClosureReviewVerdictLine`, `closureReviewVerdict`). It must refuse this repository's own committed REVISE records and accept a genuine ACCEPT record. Test it against every real review record in `docs/reviews/`. Can a legitimate ACCEPT record be refused (verdict phrased differently, verdict in a table, verdict word split across lines)? Can a REVISE record pass (verdict line absent but ACCEPT stated elsewhere, "Verdict: ACCEPT" quoted inside a discussion of what would have been accepted)?
3. **The resume-counter corroboration** and **the outcome-free guard**: re-verify both hold, including against the newly fixed `artifact_digests` shape.

**B. Interaction defects.** Review 2's P2-N1 existed only because the P1-1 fix changed what the log contains. Look for the same pattern again among the revision-3 changes: the mutex in `append`, the injectable signal channel, the `omitempty` removal, the promptness threshold, the error-path spend. Does any of them change an assumption another check depends on?

**C. The whole payload, as if fresh.** Every § 9 condition located in code; the partial summary outcome-free in both directions; the process-event log honest about what the host does not expose; the post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure; the authoring-tranche exemption; records matching the implementation with no overclaim; freeze hygiene with exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`), and the world-build and grader digests unmoved.

**D. The frozen prose.** `validation-execution-boundary.md` and `README.md` were rewritten this round to state obligations they previously omitted. Verify the prose now matches the code exactly — including the claim that the one-hour freshness bound is an implementation control and not a § 9 requirement. A custodian following the boundary document literally must be able to construct an attestation the runner accepts; walk through that construction and say whether it succeeds.

**E. Residuals.** The note carries P3-N1 (coordinated log rewrite), the host-clock residual, the post-closure gate's chair-discipline and wipe-and-restart caveats, the authoring exemption, and P3-1/P3-2/P3-5. For each, say whether it is correctly characterized and genuinely acceptable, or should be a finding.

**F. Quality gates.** Run independently: `go test ./...`; corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`. For `-race`: the payload author reports running it on the WSL2 linux host with no DATA RACE and the single expected red. If you can reproduce that, do; if not, say so plainly rather than restating the claim.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed.

Two rounds have already produced substantive fixes and the finding count is falling. **A clean ACCEPT is a legitimate and expected outcome if the payload is sound** — do not manufacture findings to justify a third round. Equally, do not soften a real finding to end the loop. Report what you verified and what you could not.
