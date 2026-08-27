# Fourth independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 4

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** any of the three reviewers who produced the previous review records. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a fourth review

Three independent reviews have run against earlier commits of this payload, and all three returned REVISE:

| Review | Target | P0 | P1 | P2 | P3 |
| --- | --- | --- | --- | --- | --- |
| 1 — `docs/reviews/2026-08-27-execution-resilience-payload-review.md` | `7a1ea42…` | 0 | 2 | 6 | 9 |
| 2 — `docs/reviews/2026-08-27-execution-resilience-payload-review-2.md` | `e92eaac…` | 0 | 0 | 2 | 10 |
| 3 — `docs/reviews/2026-08-27-execution-resilience-payload-review-3.md` | `60c6aee…` | 0 | 1 | 0 | 4 |

**Every round's findings have been defects introduced by the previous round's fixes**, each at one remove further:

- Review 2's two P2s were both in review 1's fix code. One existed *only* because the P1-1 fix made a signalled interruption write a log entry, which silently changed what a different check's timestamp anchor compared against.
- Review 3's single P1 was in review 2's fix code. Review 2 prescribed "match a verdict line rather than the bare token" on the basis of a survey of **two** files; the parser written to that prescription, and tested against those same two files, then misclassified most of this repository's 76 committed review records in both directions.

That is the pattern to watch hardest. A fix that is correct against the case that prompted it, and wrong in the wider context, has now happened three times running.

Your scope is the **whole payload**, not the deltas. Read all three prior records, but treat their verdicts as claims, not settled fact. If a prior finding was wrong, or a residual any reviewer accepted should not have been accepted, say so.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`70ff7875340071324a0006a5504bbf527c21a5f5`**. Ancestry: parent `6f54e8ca2903c4653bd319c116df57f0edf8317c` (third review record) → `60c6aee55b71716743b30d084a57daeae48b3373` (superseded revision 3) → `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` (second review record) → `e92eaacf6ab54b547b331d14749717da64b58e0d` (superseded revision 2) → `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` (first review record) → `7a1ea42513acd3d55e276eabc2459da4a037acc0` (superseded payload 1) → `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree. It is intentionally dirty with exactly one untracked file, this prompt. Anything else dirty is a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:1eecf1e6d847c4bb6222c7d4c7fee37aebdce124abdcb763c19bed8dcb923013`, **17 files**, ordinal-path set digest `sha256:2fa91f23537f0bf7ea87ba73f60256a6291cf40ba896ee5694e53f6f71a9e245`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:d9e4d21372645449b8820973e0cd5b77f8b80e9bf778704a38990e4ff5dce89b`.
- LF-normalized SHA-256: `runner/validation_gate.go` `sha256:781d84cb770d1edfd1df439847369127ce64f15513626d43b0bb5c635c5b3984`; `runner/closure_verdict_corpus_test.go` `sha256:444955e0d4724827e96dbeb0db3dbbd8a76cbf55da21a7fd9cdc7411ca688161`; `runner/execution_resilience.go` `sha256:fe98996728d7ccbd4b1bbd23f44d29e3ebb543a64e2430502a7ccbfb2dee374a`; `runner/process_events.go` `sha256:605a365310aebb4f98c5e75fca8ae96232e671a48381756557f7c891ec70c617`; `validation-execution-boundary.md` `sha256:a9a616d3eb2c90833d841353f695c6f690f270d4aa74375a1336e76af970d767`; `README.md` `sha256:fc748a2ffab7e4cf4256edb232bd5682802835391f2704e845f11082f07ed45b`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256, rendered `sha256:<hex>`. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## Governing sources

- `docs/plans/2026-08-27-protocol-v5-proposal.md` § 9 — the accepted specification.
- `docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md`, **freeze obligation 5**.
- The three prior review records named above.
- `experiments/frontier-v1/artifacts/execution-resilience-candidate.md` — the implementer's account, with a "What the Nth review changed" section per round. Claims to verify.

## What you must check

**A. The rewritten verdict parser — the round-3 P1, and this round's highest risk.** It has been rewritten structurally: a furniture-only prefix with a closed qualifier set (`final`, `overall`), lone-token extraction with non-letters trimmed, heading-then-next-non-blank-line support, verdict-list rejection via the other verdict word's absence from the token's neighbourhood, and first-statement-decisive. A new test makes the committed review corpus the arbiter.

Do not take the corpus test's passing as sufficient — it was written by the same author as the parser, and a test and an implementation sharing an author share blind spots. Independently:

1. Reimplement the parser standalone, or drive it directly, and run it over **every** `.md` under `docs/reviews/`. Report the accept/refuse counts and check each accepted file really is an accepting review record and each refused non-prompt really does state a non-ACCEPT verdict or none.
2. Attack both directions with constructed cases. False accepts to try: a verdict word inside a fenced code block or an indented block quote; a table row; `Verdict: ACCEPTED` or `ACCEPTABLE`; a document whose *first* verdict statement is an ACCEPT quoted from a different record; a filename or heading containing the furniture characters. False refusals to try: a genuine ACCEPT record whose verdict line carries a leading footnote marker, a numbered `### 2. Verdict:` heading, `Verdict — ACCEPT` with a dash instead of a colon, non-breaking or unusual whitespace, CRLF line endings.
3. Judge the closed qualifier set (`final`, `overall`) and the three-field neighbourhood window: are they principled, or fitted to the corpus that happens to exist? What happens to a future record using a qualifier not in the set?
4. Judge whether making the live `docs/reviews/` corpus a test input is sound. It couples the frozen runner's test outcome to files added later by unrelated work. Could a future review record legitimately break this test, and is that acceptable or a defect?

**B. Interaction defects, again.** Check each revision-4 change against assumptions elsewhere: the parser rewrite, the boundary document's new attestation field list, and the new corpus test. In particular, does the boundary document's enumerated field set exactly match what `resumeAttestation` declares and `decodeStrict` will accept — every name, no omission, no extra? A custodian following it literally must be able to construct an accepted attestation; construct one and say whether it works.

**C. The whole payload, as if fresh.** Every § 9 condition located in code, with **both** questions answered for each — can an invalid resume pass, can a valid one be refused. The partial summary outcome-free in both directions. The process-event log honest about what the host does not expose. The post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure. The authoring-tranche exemption. Records matching the implementation with no overclaim. Freeze hygiene: exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`), the world-build and grader digests unmoved, and the inventory's 17 files matching the working tree.

**D. Residuals.** The note and review 3's §8 carry a residual list. For each, say whether it is correctly characterized and genuinely acceptable, or should be a finding. Note especially that review 3 said the chair-discipline residual would need extra wording *only if* the P1 went unfixed — check whether the fix means review 2's original wording now stands.

**E. Quality gates.** Run independently: `go test ./...`; corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`. For the race detector, Windows cannot run it (no GNU toolchain); use `wsl -e bash -lc 'cd /mnt/d/Work/personal/phoenix && go test -race ./experiments/frontier-v1/runner/'`. State whether you ran it yourself.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. For every finding state explicitly whether it permits an invalid resume or refuses a valid one. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed.

Three rounds have produced substantive fixes and the count is falling (2 P1 → 0 P1 → 1 P1, and 6 P2 → 2 P2 → 0 P2). **A clean ACCEPT is a legitimate and expected outcome if the payload is sound** — do not manufacture findings to justify a fifth round. Equally, do not soften a real finding to end the loop: three rounds have each caught a defect introduced by the previous round's repair, so the prior on "the fix introduced something" is not low. Report what you verified and what you could not.
