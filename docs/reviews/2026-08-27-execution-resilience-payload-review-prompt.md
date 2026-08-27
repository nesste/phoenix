# Independent review prompt — protocol-v5 section 9 execution-resilience payload

You are an evaluation reviewer independent of this payload's implementation. You did not author it. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`7a1ea42513acd3d55e276eabc2459da4a037acc0`** (the commit titled "experiments: protocol-v5 section 9 execution-resilience payload"). Its parent must be `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025).
- Do not commit anything, and do not amend the working tree, while you review. If the working tree is dirty with a review-prompt file, ignore it; it is committed with your record afterwards.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:2e6c3c90d73f15c1713bd026237fefec48adbdf3a6abec369350cc3b66f99f93`, 16 files, ordinal-path set digest `sha256:1f0b6f872f58d17ddb987c0106c0c0df3dbd6385d060ac7991c97fafcde7c609`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:3749438786e13fa822fd9f0b99052b4ce8aca53bd6fb4eabf016248c4243eb80`.
- LF-normalized SHA-256 of the key changed bytes: `protocol.json` `sha256:726b68112108542ba3491fde538342399cbb446a04755a7e0208e076cbdb0d32`; `validation-execution-boundary.md` `sha256:247e5f8cc4ee355ea23e8414c3ca4389070231e5cc0570deb1853d008f2b081e`; `pre-validation-artifacts.json` `sha256:47f4bd4af3635510eb932dc8528372747eb665acd47b2beda653192a1a61bed3`; `runner/execution_resilience.go` `sha256:4f6cf8015f66d6814711d27409bdbc9ac070405ee9cab10c042660a66b5a3696`; `runner/process_events.go` `sha256:0dabcd7855ab2a457bb5cea8394fdbba1558add0136b560fddb723c080aa9e17`.

Digest rule everywhere: CRLF→LF, then CR→LF, then SHA-256, rendered `sha256:<hex>`. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## Governing sources

- `docs/plans/2026-08-27-protocol-v5-proposal.md` § 9 "Execution resilience and precommitted mandatory resume" — the accepted specification. This payload must implement it, not reinterpret it.
- `docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md` — the acceptance record; **freeze obligation 5** is the one this payload discharges.
- `experiments/frontier-v1/artifacts/execution-resilience-candidate.md` — the implementer's account. Treat it as a claim to verify, not as evidence.

## Scope

In scope: everything in the 16-file inventory, plus the two uninventoried payload-support files `runner/execution_resilience_candidate_artifact_test.go` and the candidate note. Out of scope: sections 1–8 and 10, already settled by decisions 0024 and 0025; the absence-grading guard; the schedule and corpus, which the v5 seal replaces.

## What you must check

**1. Section 9 conditions, one by one.** For each of (a) no-inspection, (b) digest identity, (c) the 72-hour window, (d) pre-decision cause classification on the frozen list — locate the enforcing code, and state whether a resume can succeed with that condition unmet. Specifically:

- Is OOM refused? Is an unlisted or empty cause refused? Is the refusal message accurate?
- Can a `process_kill` resume without `kill_principal_established_as_non_participant`? Is the process-event log digest re-verified against the file on disk, so a stale or forged citation fails?
- Is the timestamp ordering enforced (interruption → classification → verification → authorization), and is a resume authorized more than 72 hours after the interruption refused?
- Is the promptness duty encoded, and does a late authorization require a written lapse explanation?
- Is "at most one resume per tranche" enforced *before* the attestation is consulted, and is the counter carried durably in the checkpoint so a second interruption cannot be laundered by supplying a fresh attestation? Can a hand-edited checkpoint raise the count?
- Does any code path allow abandonment-by-choice as an outcome, i.e. a non-completed run that is neither resumed nor indeterminate?

**2. The partial summary is outcome-free.** Read `scheduledPartialSummary` and everything written into it. Does any field carry a success count, grade tally, pass/fail status, or per-arm outcome? Is the both-directions field-set test real, or can a new field slip past it? Is spend the only outcome-adjacent quantity, and is its retention justified by the budget-stop rule as § 9 says?

**3. The process-event log.** Does it record start, stop, checkpoint, and termination with timestamps? Is it genuinely append-only and durable against an abrupt kill? Is the "originating principal where the host exposes it" claim honest — in particular, is the payload right that a Go signal handler cannot identify the sender on either supported host, and does the termination entry say so rather than misattributing the kill to the receiving process's owner? Is the log's digest committed where § 9 requires (partial summary now, custody record at closure)?

**4. The post-closure gate.** Is it encoded in the *frozen gate rules*, not merely in prose? Does it fire for every non-completed ending — interruption, budget stop, and lapse — and not only for interruption? Does it correctly apply to the decision-0023 closure, which had no independent review, and is that registered as a blocker? Can it be satisfied by a record that does not exist, or by a non-ACCEPT verdict?

**5. Scope decision.** The payload exempts the authoring tranche from the attestation rule. Section 9 does not say "validation" in every sentence. Judge whether the exemption is defensible or whether it is a gap; say so explicitly either way.

**6. Records match the implementation.** Read `protocol.json` `amendment.execution_resilience`, the boundary document's new section, and the README paragraphs against the code. Flag any statement that overclaims. Verify the blocker list went 7→8, that the struck section 9 clause is genuinely discharged, that the added closure blocker is accurate, and that `corpusctl` `protocol_test.go` asserts the new contract rather than merely tolerating it.

**7. Freeze hygiene.** Confirm exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`) and that `TestLocalArtifactCandidateMatchesImplementation`, `TestValidationBuildReproducesFrozenWorldBuildDigest`, and `verifyAcceptedLocalArtifacts` stay green. Confirm the Phoenix binary is untouched (nothing under `cmd/phoenix`'s import graph changed) so the world-build identity `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` holds, and that the grader digest is unchanged at `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b` because no corpusctl non-test source, `go.mod`/`go.sum`, or `schema/*.schema.json` changed. Judge whether editing the `gates` block of `pre-validation-artifacts.json` in a payload commit is acceptable given that no digest, `frozen` flag, or accepted-candidate pointer moved.

**8. Quality gates.** Independently run: `go test ./...`, the corpusctl suite, `go vet ./...`, staticcheck, `gocyclo -over 15` and `dupl -t 100` over the gated directories, `validate-spec`, `validate-authoring`. Report anything beyond the single expected red. Note that carried finding P2-1 is claimed retired by decomposition — verify the measured complexity.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0 (blocks acceptance), P1 (must fix before refreeze), P2 (fix or accept explicitly as a recorded residual), P3 (note only). For every finding give file and line, the failure it permits, and what would fix it. If you accept, list the residuals you accept so they can be transcribed into the refreeze's `accepted_findings`. Report the digests you computed for every pinned input, so the transcription can be verified against your record rather than trusted.
