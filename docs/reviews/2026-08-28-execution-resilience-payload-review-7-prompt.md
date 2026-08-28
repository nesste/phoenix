# Seventh independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 7

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** any of the six reviewers who produced the previous review records. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a seventh review

Six independent reviews have run, and all six returned REVISE:

| Review | Target | P0 | P1 | P2 | P3 |
| --- | --- | --- | --- | --- | --- |
| 1 | `7a1ea42…` | 0 | 2 | 6 | 9 |
| 2 | `e92eaac…` | 0 | 0 | 2 | 10 |
| 3 | `60c6aee…` | 0 | 1 | 0 | 4 |
| 4 | `70ff787…` | 0 | 0 | 2 | 4 |
| 5 | `a368433…` | 0 | 0 | 1 | 4 |
| 6 | `b375d43…` | 0 | 0 | 1 | 3 |

Rounds 3, 4 and 5 each found defects in a closure-record verdict *parser* that read a review record's prose. On the fifth review's recommendation and the chair's decision, revision 6 **deleted** it and replaced it with a digest pin. The sixth review endorsed that deletion — and found that the *justification* frozen alongside it named a compensating control this project does not have: it asserted the judgment was "already made by the independent reviewer of the refreeze that sets these fields", and no refreeze in this project has ever been independently reviewed.

**This revision replaces that assumption with a written obligation and a frozen anchor.** That is what you are reviewing, and the question is whether it actually discharges § 9's requirement that the closure record "has been independently reviewed and committed", or whether it is still assumption-shaped one level up.

Read all six prior records, but treat their verdicts as claims, not settled fact — including the sixth review's N1 prescription, which this revision implements and which you may judge insufficient or misdirected.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`27edb8664a060d0a7690039e8e29dba94facada0`**. Ancestry: `3fc001e247332a5dfdb9d99fc54fc0a9681e8ea1` (sixth review record) → `b375d432df1facf29a8e3e66f065b19453a8ede4` (superseded revision 6) → `950df9345c34ddab89cb93cfc8944fe37c861763` → `a368433089a857dbec4d72ee74123a4ee8743f47` → `fec6d554156eb749057fb8b6cd3db927c2fecf05` → `70ff7875340071324a0006a5504bbf527c21a5f5` → `6f54e8ca2903c4653bd319c116df57f0edf8317c` → `60c6aee55b71716743b30d084a57daeae48b3373` → `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` → `e92eaacf6ab54b547b331d14749717da64b58e0d` → `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` → `7a1ea42513acd3d55e276eabc2459da4a037acc0` → `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree. It is intentionally dirty with exactly one untracked file, this prompt. Anything else dirty is a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:e87ea77d058f9100252e1391f96de8398f2189aaab213063f5f84c87b431ac3c`, **16 files**, ordinal-path set digest `sha256:d36d948d2f15ca619d2c34d2962c7c6791759f83b80516d87bb2b027add95d50`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:5e7a1fba96cc2e0afb7a42d2d7e9712f5b5f5e21981b103a489cb03eb1401d6d`.
- LF-normalized SHA-256: `runner/validation_gate.go` `sha256:810535ab8b509afd95070e19f4fff6b6d10ba81e6021658b4219773286ea2223`; `runner/execution_resilience_test.go` `sha256:5571d6c92bb2f32cf568c2e39e2011cb0e2b7e8d7cccf00fa76e2cedb18e4fb7`; `runner/artifact_freeze_test.go` `sha256:890bca5afa424de7847d5b37e1a0ac25f1e6ac34c7273e0556560fa9fc231abb`; `pre-validation-artifacts.json` `sha256:1edecb7da598adbebc8f22db73ef54ee2138918cb408baca91252dded2b45a53`; `validation-execution-boundary.md` `sha256:e3c1b3bcd7d60ef9d90bfc9d80f813226f4c36a8d315b8003169b7075aabe14d`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`, and **the digest field in each line carries the `sha256:` prefix** — the sixth reviewer lost time to that ambiguity, so it is stated here.

## What you must check

**A. Does the obligation actually bind? — this round's central question.**

The repair has three parts: imperative text in three frozen locations; the obligation restated in the boundary document's **Authorization boundary** section; and a fourth gate field, `validation_execution_closure_fields_payload_review`, naming the independent review of the payload commit that set the closure fields, enforced by `verifyClosureFieldsReview`.

1. **Walk the sixth review's attack again against this revision.** A future `OPEN_GATE_1A` decision commits atomically: `may_open_validation: true`, the four closure fields, and any relaxation of `verifyPostClosureGateState` it needs. Does anything now stop it, or make it visible? Be concrete about what the fourth field costs an inattentive chair and what it costs a motivated one, and say plainly whether the difference is worth the field.
2. **Is the obligation binding, or merely written?** It lives in prose in the boundary document and the gate rule. Nothing executes it. Determine whether this project has any mechanism that makes a written obligation on a *future reviewer* stick — and if the honest answer is "the reviewer reads the boundary document", say whether that is sufficient for § 9's "independently reviewed and committed" or whether § 9 is now discharged by convention.
3. **Does the fourth field create a regress?** It names a review record. Nothing verifies that record reviewed *this* payload, reached ACCEPT, or is not the closure record itself. Is that a defect, a residual, or correctly out of scope? Consider whether the field could be satisfied by naming this very prompt, or the closure record, or an unrelated old review.
4. `artifact_freeze_test.go` and `pre-validation-artifacts.json` are still in no frozen block — the sixth review established this and this revision does not change it. Judge whether that undermines the repair, and whether it should be fixed here or is properly a separate concern.

**B. The replacement's correctness.** `verifyClosureFieldsReview` mirrors `verifyClosureReviewRecord` but without a digest pin. Attack both: paths, traversal, symlinks, directories, empty files, case sensitivity. Is the asymmetry — the closure record is digest-pinned, its review is not — deliberate and defensible, or an oversight?

**C. Nothing regressed.** The duplicated doc comment is claimed removed; confirm exactly one occurrence. Confirm the four closure fields are pinned empty by `verifyPostClosureGateState` and its rule-text substrings match the gate document. Confirm the inventory's 16 files match the working tree, and that the `accepted_residuals` array now carries the full operative list from review 6 § 9 — reviews 4, 5 and 6 each flagged it as lagging.

**D. The whole payload, as if fresh.** Every § 9 condition located in code with **both** questions answered — can an invalid resume pass, can a valid one be refused. Partial summary outcome-free. Process-event log honest about what the host does not expose. Post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure. Authoring exemption. Records matching the implementation with no overclaim — the candidate note now carries a seven-revision narrative; check its present-tense claims against the current code, since stale claims in it have been a finding in three consecutive rounds. Freeze hygiene: exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`), world-build and grader digests unmoved.

**E. Residuals.** Review 6 § 9 is the operative list, now mirrored in the inventory. For each entry say whether it is correctly characterized and genuinely acceptable. Say explicitly whether the "runner forms no opinion" residual is now acceptable — review 6 made it conditional on N1 being fixed, and this revision claims to have fixed it.

**F. Quality gates.** Run independently: `go test ./...`; corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`. Race detector via `wsl -e bash -lc 'cd /mnt/d/Work/personal/phoenix && go test -race ./experiments/frontier-v1/runner/'`. State whether you ran it yourself.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. For every finding state explicitly whether it permits an invalid resume or refuses a valid one. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed.

Six rounds have each found something real, and the trend is 2 P1 → 0 → 1 → 0 → 0 → 0, with P2 counts falling to one. This revision is small: three prose changes, one field, one deleted duplicate comment. **A clean ACCEPT is a genuinely likely outcome** — do not manufacture findings to justify an eighth round.

Equally: the last two rounds each found that a claim frozen alongside correct code was false. If this revision's central claim — that a written obligation on a future reviewer discharges § 9 — is itself unsupported, that is the finding, and it does not matter that the code is fine. Say so plainly.
