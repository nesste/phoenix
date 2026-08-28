# Sixth independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 6

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** any of the five reviewers who produced the previous review records. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a sixth review, and what is different about it

Five independent reviews have run, and all five returned REVISE:

| Review | Target | P0 | P1 | P2 | P3 |
| --- | --- | --- | --- | --- | --- |
| 1 | `7a1ea42…` | 0 | 2 | 6 | 9 |
| 2 | `e92eaac…` | 0 | 0 | 2 | 10 |
| 3 | `60c6aee…` | 0 | 1 | 0 | 4 |
| 4 | `70ff787…` | 0 | 0 | 2 | 4 |
| 5 | `a368433…` | 0 | 0 | 1 | 4 |

Rounds 3, 4 and 5 each found defects in **one** component: the closure-record verdict parser, which read a review record's prose to corroborate that it stated ACCEPT. The fifth review recommended deleting that component rather than hardening it a fourth time, and **the chair accepted and directed the deletion.** This revision carries it out.

**Your review is therefore differently shaped from the five before it.** The obvious risk is no longer that a heuristic is subtly wrong. It is that **the deletion gave up something the gate needed.** Interrogate the removal at least as hard as the replacement. A finding that the payload has weakened a control is exactly what this round exists to catch, and it is a legitimate outcome.

Read all five prior records, but treat their verdicts as claims, not settled fact — including the fifth review's § 9 recommendation, which is the basis for this change and which you may disagree with.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`b375d432df1facf29a8e3e66f065b19453a8ede4`**. Ancestry: `950df9345c34ddab89cb93cfc8944fe37c861763` (fifth review record) → `a368433089a857dbec4d72ee74123a4ee8743f47` (superseded revision 5) → `fec6d554156eb749057fb8b6cd3db927c2fecf05` → `70ff7875340071324a0006a5504bbf527c21a5f5` → `6f54e8ca2903c4653bd319c116df57f0edf8317c` → `60c6aee55b71716743b30d084a57daeae48b3373` → `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` → `e92eaacf6ab54b547b331d14749717da64b58e0d` → `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` → `7a1ea42513acd3d55e276eabc2459da4a037acc0` → `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree. It is intentionally dirty with exactly one untracked file, this prompt. Anything else dirty is a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:3874e629dc359dc05b5b6607d61347547ee75e2388d7808bee2fa95061298970`, **16 files**, ordinal-path set digest `sha256:1ba17cd292b53122a483aab037a5cf03a0108991fd1df55aca72faf26eb1e834`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:ac478f8593db5e7beffebe8e037bce725b16af068b0196493d4b94a349e3d245`.
- LF-normalized SHA-256: `runner/validation_gate.go` `sha256:741e19b0cec0d83d8348f931dec8d371e4de61217e53ff4cf762f63200ddd8c9`; `runner/execution_resilience_test.go` `sha256:e6b2cc06f707ae42c2e5c3c459bc6c4d75916ee9480efd080667d572afc2d837`; `runner/artifact_freeze_test.go` `sha256:0b5b3662bea710600abd49518c6d1250f4f6d699a0fb300ec5f81176d81bbb8a`; `pre-validation-artifacts.json` `sha256:40f163717c4ead7ac237e931d50b96c2766166e003c0cdf5016ac7bc2b33ef18`; `validation-execution-boundary.md` `sha256:a81514b8f4dea6678362df0ce4eb769378176e64f5307b02f789326f19ba51f5`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## What you must check

**A. Did the deletion give up something the gate needed? — this round's central question.**

The removed component read the named closure record's text and refused if it did not state ACCEPT. It is replaced by: the gate document names the record, states its verdict, and pins the record's LF-normalized digest; `verifyClosureReviewRecord` refuses a path outside `docs/reviews/`, a traversal, a non-regular or empty file, and a digest mismatch.

1. Construct the concrete attack the parser used to stop and the digest pin does not. Walk it end to end: who does what, what the freeze workflow requires of each step, and where it is caught or is not. Is the residual exposure real, or is it fully absorbed by the human review of the refreeze that sets these fields — as the candidate note and the boundary document both claim?
2. The claim that the judgment is "already made by the independent reviewer of the refreeze" is load-bearing. **Test it against how this project actually works.** Does a refreeze reviewer, under the prompts this project issues, in fact read the closure record and rule on whether it is an accepting independent review? Is that obligation written anywhere binding, or is it assumed? If it is assumed, that is a finding.
3. Is the digest pin genuinely stronger on the amended-after-review case, as claimed? Verify by execution, not by reading.
4. Would a narrow, canonical-form verdict check — one exact form mandated by the boundary document, no heuristics — recover the lost property at acceptable cost? Say whether you would want it, and be concrete about the form. If you think the deletion is simply right, say that instead.

**B. The replacement's own correctness, in both directions.** Can an invalid closure record satisfy the gate? Can a valid one be refused? Attack: symlinks and directory entries; a path with `./` or backslashes or a leading slash; case differences on a case-insensitive filesystem; a digest in the wrong case or without the `sha256:` prefix; CRLF versus LF in the record; a record that is valid but pinned with the raw rather than the LF-normalized digest. Check that `digestLFNormalizedFile` is the same rule the rest of the project uses.

**C. Nothing else regressed in the deletion.** Confirm no dangling reference to any removed symbol; that the two deleted files are gone from the inventory and that the inventory's 16 files match the working tree; that `verifyPostClosureGateState` pins all three closure fields empty and its rule-text substrings still match the gate document; and that no unrelated behaviour changed. The candidate guard's base commit and file count must track the new base.

**D. The whole payload, as if fresh.** Every § 9 condition located in code with **both** questions answered — can an invalid resume pass, can a valid one be refused. Partial summary outcome-free. Process-event log honest about what the host does not expose. Post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure. Authoring exemption. Records matching the implementation with no overclaim — the candidate note now contains a long historical narrative across six revisions; check its present-tense claims against the current code. Freeze hygiene: exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`), world-build and grader digests unmoved.

**E. Residuals.** Review 5 § 8 is the operative list. Three of its entries concern the deleted parser (block context, qualifier set, pinned corpus) and should now be retired — confirm that, and say whether the deletion introduces any residual that should be added. For each surviving entry, say whether it is correctly characterized and genuinely acceptable.

**F. Quality gates.** Run independently: `go test ./...`; corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`. Race detector via `wsl -e bash -lc 'cd /mnt/d/Work/personal/phoenix && go test -race ./experiments/frontier-v1/runner/'`. State whether you ran it yourself.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. For every finding state explicitly whether it permits an invalid resume or refuses a valid one. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed.

Five rounds have each found something real. This round's payload is a **net deletion** — roughly 200 lines of frozen code and data removed for about ten added — so the usual "the fix introduced something" prior is weaker, and a clean ACCEPT is a genuinely likely outcome. Do not manufacture findings to justify a seventh round. Equally, do not accept a weakened control because deletion feels safer than heuristics: if the gate now permits something it should not, say so plainly.
