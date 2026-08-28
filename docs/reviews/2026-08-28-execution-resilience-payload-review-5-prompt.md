# Fifth independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 5

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** any of the four reviewers who produced the previous review records. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a fifth review

Four independent reviews have run, and all four returned REVISE:

| Review | Target | P0 | P1 | P2 | P3 |
| --- | --- | --- | --- | --- | --- |
| 1 — `…-payload-review.md` | `7a1ea42…` | 0 | 2 | 6 | 9 |
| 2 — `…-payload-review-2.md` | `e92eaac…` | 0 | 0 | 2 | 10 |
| 3 — `…-payload-review-3.md` | `60c6aee…` | 0 | 1 | 0 | 4 |
| 4 — `2026-08-28-…-review-4.md` | `70ff787…` | 0 | 0 | 2 | 4 |

**Every round's findings have been defects introduced by the previous round's fixes.** Rounds 3, 4 and now 5 have all centred on one component: the closure-record verdict parser in `validation_gate.go`, which corroborates that the record the post-closure gate names really does state ACCEPT.

You should know two things about how this round was written.

**First, the implementer declined part of review 4's prescription, deliberately.** Review 4's fix for its N2 was "reject the token whenever any other verdict word appears anywhere on the decisive line". The implementer rejected that as over-reaching and adopted a narrower rule instead — a token followed by a comma or slash, with another verdict word later on the line, is an enumeration; a token followed by prose is a statement. The stated reason: review 3's own record has the verdict line `**Verdict: REVISE** — … it refuses 17 committed ACCEPT records`, which review 4's literal rule would refuse to read, and review 4 verified its prescription against *accepting* records only. **Judge that call.** Was declining correct? Is the replacement rule sound in both directions, or does it now admit something review 4's rule would have caught?

**Second, the chair has an open question about whether this component should exist in this form at all.** An alternative design was raised: pin the closure record by digest in the gate document and take the verdict from the frozen gate fields, reducing the runner's check to a digest comparison with no prose parsing. If your review concludes the parser is still generating defects, say so plainly and say whether you think the heuristic approach is the wrong shape for a governance gate — that judgment is wanted, not out of scope.

Your scope is the **whole payload**. Read all four prior records, but treat their verdicts as claims, not settled fact.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`a368433089a857dbec4d72ee74123a4ee8743f47`**. Ancestry: `fec6d554156eb749057fb8b6cd3db927c2fecf05` (fourth review record) → `70ff7875340071324a0006a5504bbf527c21a5f5` (superseded revision 4) → `6f54e8ca2903c4653bd319c116df57f0edf8317c` (third review record) → `60c6aee55b71716743b30d084a57daeae48b3373` → `9f53fa73fbede7e240aeb22e9b02d287c9c39a94` → `e92eaacf6ab54b547b331d14749717da64b58e0d` → `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` → `7a1ea42513acd3d55e276eabc2459da4a037acc0` → `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree. It is intentionally dirty with exactly one untracked file, this prompt. Anything else dirty is a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:2309d017d58e49e4dceb9c1dec37823e0902200ac31e98813844913a3da61bad`, **18 files**, ordinal-path set digest `sha256:44192b3c80911ee22489b0537f30dbddaea5003576b93fa8543e051f59995f67`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:b69c7ff802a10d124493ac1c14603a2a601bd11453cc6258e71d1919a9a42d0b`.
- LF-normalized SHA-256: `runner/validation_gate.go` `sha256:5f97d6e13b8fcad5f3c32155eec4359f72aba2254da5856e3a036ededf6061c4`; `runner/closure_verdict_corpus_test.go` `sha256:a400c96336cef8814dd9c8d1548aedc1733e2910b8979c44112363b451c18c7d`; `artifacts/closure-verdict-corpus-expectations.json` `sha256:b761d71d9279775ca659ef2e1785ef3148ef2329dfa4cf4bec0cffba7d992694`; `validation-execution-boundary.md` `sha256:778926ba22b888964f8d40509d08684eaa9225c02076dd0be7aeb3d47f643351`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## What you must check

**A. The verdict parser and its new pinned corpus — this round's highest risk.**

1. Drive the parser yourself over every `.md` under `docs/reviews/`. Report your counts. Then check them against `artifacts/closure-verdict-corpus-expectations.json`, which is this round's committed ground truth. **The expectations file was derived by the implementer, by a scan the implementer also wrote.** It is asserted to be independent of `validation_gate.go`; verify that claim by spot-reading a sample of records and confirming the label matches what the document actually states. Are any labels wrong? A wrong label is worse than a wrong parser, because it silently defines correctness.
2. Attack the enumeration rule specifically, in both directions. It fires only on a comma or slash directly after the token. Try: `Verdict: ACCEPT or REVISE` (no comma); `Verdict: ACCEPT / REVISE` (spaced slash); `Verdict: ACCEPT; REVISE; REJECT`; a prompt that writes the verdict list across two lines; and conversely a genuine ACCEPT record whose line reads `ACCEPT, with the residuals below` (comma, no other verdict word).
3. Attack the block-context tracking: nested fences, a fence opened inside a block quote, `~~~` mixed with backticks, a fence never closed, a fence marker inside an indented block.
4. Attack the separator requirement and the label furniture: a table row that happens to satisfy the furniture rule, a verdict inside an HTML comment, a heading whose token line is a list item.
5. Judge whether the parser is now correct, or merely correct against everything anyone has thought to try. Say which.

**B. Interaction defects.** Rounds 2, 3 and 4 each found one. Check the revision-5 changes — the parser rewrite, the removal of `>` from the furniture cutset, the addition of `\r` to it, the pinned expectations replacing the live sweep, and the boundary document's new paragraphs — against assumptions elsewhere in the payload.

**C. The whole payload, as if fresh.** Every § 9 condition located in code with **both** questions answered — can an invalid resume pass, can a valid one be refused. Partial summary outcome-free. Process-event log honest about what the host does not expose. Post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure. Authoring exemption. Records matching the implementation with no overclaim. Freeze hygiene: exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`), world-build and grader digests unmoved, the inventory's 18 files matching the working tree.

**D. Residuals.** Review 4 § 8 carries the operative list. For each, say whether it is correctly characterized and genuinely acceptable, or should be a finding.

**E. Quality gates.** Run independently: `go test ./...`; corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`. Race detector via `wsl -e bash -lc 'cd /mnt/d/Work/personal/phoenix && go test -race ./experiments/frontier-v1/runner/'`. State whether you ran it yourself.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. For every finding state explicitly whether it permits an invalid resume or refuses a valid one. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed.

Four rounds have each found something real, and the count is now 0 P0 / 0 P1 / 2 P2 for two rounds running. **A clean ACCEPT is a legitimate and expected outcome if the payload is sound** — do not manufacture findings to justify a sixth round. Equally, do not soften a real finding to end the loop. If your honest conclusion is that the payload is sound but the parser is the wrong mechanism for this job, ACCEPT with that recorded as a residual and say so — the chair is weighing exactly that question.
