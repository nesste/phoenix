# Seventh independent review — protocol-v5 section 9 execution-resilience payload, revision 7

- **Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author it, and I am not any of the six reviewers who produced the prior records.
- **Date:** 2026-08-28
- **Payload commit:** `27edb8664a060d0a7690039e8e29dba94facada0`
- **Verdict:** **ACCEPT** — 0 P0 · 0 P1 · 0 P2 · 4 P3

**The repair holds, and the central claim is supported.** Review 6's N1 was that the frozen bytes named a compensating control — "already made by the independent reviewer of the refreeze" — that does not exist. That was a *factual* assertion about existing machinery, and it was false. This revision replaces it with a *deontic* one: the fields may be set only by an independently reviewed payload commit, a chair gate patch must not set them, and the gate must name the review that discharged the obligation. A rule cannot be false the way an assertion can; it can only be unenforced. The payload does not claim it is enforced — it says twice, in the frozen bytes, that the runner cannot check the review's conclusion — and it explicitly acknowledges the contrary precedent rather than hiding it. That is the honest shape.

My four P3s are text corrections and one refreeze-or-later strengthening. **None of them warrants an eighth round**, and I say so deliberately: I looked for a P2 and did not find one.

No model was run, no private grade obtained, no validation outcome observed. Both gates remain closed, verified in `protocol.json` and `pre-validation-artifacts.json`.

---

## 1. Pinned inputs — every digest I computed

`git rev-parse HEAD` = `27edb8664a060d0a7690039e8e29dba94facada0`. Ancestry reproduced exactly, fourteen commits back to `2cf99331…` (decision 0025).

**Working tree:** exactly one untracked file, this prompt. Nothing else. **No P0.**

| Artifact | Rule | Computed | Match |
| --- | --- | --- | --- |
| `artifacts/gate-1a-execution-resilience-candidate.json` | raw | `sha256:e87ea77d058f9100252e1391f96de8398f2189aaab213063f5f84c87b431ac3c` | yes |
| `artifacts/execution-resilience-candidate.md` | raw | `sha256:5e7a1fba96cc2e0afb7a42d2d7e9712f5b5f5e21981b103a489cb03eb1401d6d` | yes |
| `runner/validation_gate.go` | LF | `sha256:810535ab8b509afd95070e19f4fff6b6d10ba81e6021658b4219773286ea2223` | yes |
| `runner/execution_resilience_test.go` | LF | `sha256:5571d6c92bb2f32cf568c2e39e2011cb0e2b7e8d7cccf00fa76e2cedb18e4fb7` | yes |
| `runner/artifact_freeze_test.go` | LF | `sha256:890bca5afa424de7847d5b37e1a0ac25f1e6ac34c7273e0556560fa9fc231abb` | yes |
| `pre-validation-artifacts.json` | LF | `sha256:1edecb7da598adbebc8f22db73ef54ee2138918cb408baca91252dded2b45a53` | yes |
| `validation-execution-boundary.md` | LF | `sha256:e3c1b3bcd7d60ef9d90bfc9d80f813226f4c36a8d315b8003169b7075aabe14d` | yes |

**Inventory: 16 files.** All sixteen recomputed with an independent Python implementation of the digest rule — **zero mismatches, zero missing**. Per-file digests: README.md `fc748a2f…`, corpusctl/protocol_test.go `4495d032…`, pre-validation-artifacts.json `1edecb7d…`, protocol.json `726b6811…`, artifact_freeze_test.go `890bca5a…`, execution_resilience.go `fe989967…`, execution_resilience_test.go `5571d6c9…`, main.go `192d275e…`, process_events.go `605a3653…`, scheduled_resume.go `aa050eb4…`, scheduled_resume_test.go `a090f71f…`, scheduled_run.go `0956afad…`, types.go `fce5d7b9…`, validation_execution_test.go `9ef5e4c7…`, validation_gate.go `810535ab…`, validation-execution-boundary.md `e3c1b3bc…`.

**Ordinal-path set digest reproduced:** `sha256:d36d948d2f15ca619d2c34d2962c7c6791759f83b80516d87bb2b027add95d50`. The prompt's clarification about the `sha256:` prefix was correct and saved the time review 6 lost.

---

## 2. Quality gates — all run by me at this commit

| Gate | Result |
| --- | --- |
| `go test ./...` | **1 failure, the expected one.** `TestPreValidationFreezeMatchesAcceptedCandidates` at `artifact_freeze_test.go:155`, on `README.md`. Every other package `ok`. |
| corpusctl suite | `ok`, clean |
| `go vet ./...` | Clean |
| staticcheck v0.7.0 | Clean |
| `gocyclo -over 15` | Clean |
| `dupl -t 100` | Clean |
| `make validate-spec` | All 9 examples valid |
| `make validate-authoring` | Valid |
| `gofmt -l runner/` | Clean |
| **`go test -race`** | **I ran this myself** on WSL2. **No `DATA RACE`.** Same single freeze red, 42.183s. |

**Exactly one red, and it is the declared one.** The gate-state half of that same test (`verifyPostClosureGateState`) passes, which is what confirms the four fields are pinned empty and the rule text matches.

**World-build and grader identities unmoved.** I diffed `pre-validation-artifacts.json` against decision 0025 in full: the *only* change is the addition of `validation_execution_closure_gate` and the four empty closure fields. Every other digest, `frozen` flag and accepted-candidate pointer is byte-identical. The edit strictly tightens; nothing is opened.

---

## 3. Area A — does the obligation actually bind? (the centrepiece)

### 3.1 Review 6's attack, walked again against this revision

**The structural precondition is unchanged, and I verified it myself rather than inheriting it.** `grep` for `artifact_freeze_test` in `pre-validation-artifacts.json` returns **0**; for `pre-validation-artifacts` in the same file, **0**. Neither the gate document nor the test pinning its fields empty is in any frozen file set. Both remain chair-editable without tripping any freeze test, and commit `9a28d54` is the standing precedent for editing exactly that pair in one unreviewed commit.

**What the fourth field costs an inattentive chair.** They set three fields, following the pre-revision-7 schema or the decision-0022 pattern. `requirePostClosureAudit` now refuses: *"…requires a closure record that is named, digest-pinned, independently reviewed ACCEPT, and set by a reviewed payload commit."* To find out what the fourth field is they must open the `gates` block they are already editing, where `validation_execution_closure_gate` tells them, in capitals, that setting these fields is an **OBLIGATION on the commit that sets these fields** and that **a chair gate patch must not set them**. The prohibition is placed at the exact point of the violation, in the file they have open, at the moment they are blocked. That is well designed, and it is a real cost.

One honest qualification on timing. The tests exercising `requirePostClosureAudit` use a synthetic root, not the live document, so the runner's refusal does not arrive while the chair is authoring the decision — the freeze test arrives first, and the chair must relax it anyway. The runner's refusal arrives later, at validation-execution time. It still lands before any model runs, which is what §9 cares about, but it is a second gate rather than a first one.

**What it costs a motivated chair: one string.** I tested this rather than asserting it. I extracted `verifyClosureFieldsReview` verbatim into a standalone module and ran it against real repository paths:

```
ACCEPTED   docs/reviews/2026-08-28-…-review-7-prompt.md   (this review's own prompt)
ACCEPTED   docs/reviews/2026-08-27-…-payload-review.md    (review 1, REVISE)
ACCEPTED   docs/reviews/2026-08-28-…-review-6.md          (review 6, REVISE)
ACCEPTED   docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md (unrelated)
ACCEPTED   docs/reviews//…-review-6.md                    (doubled separator)
refused    docs/reviews              (prefix)
refused    docs/reviews/.            (not a regular file)
refused    docs/reviews/absent.md    (stat)
refused    README.md                 (prefix)
refused    docs/reviews/../README.md (traversal)
```

**This very prompt satisfies the field.** So does a REVISE record, an unrelated old review, and — since the closure record must itself live under `docs/reviews/` and be non-empty — the closure record named in the first field. Nothing checks distinctness, target commit, or verdict.

**Is the difference worth the field? Yes — but not for the reason the payload gives.** The payload's justification is that the runner "can refuse a gate that names none — which is what setting these fields outside a reviewed payload commit would leave behind." That is the weakest available argument, because an empty fourth field is precisely the state a chair *trying to open the gate* will not leave behind: the runner refuses, and their own objective fails. (Graded as N4.)

The real value is different and worth stating plainly. Under revision 6 the gate document asserted **nothing** about who reviewed; the claim lived only in a code comment, and it was false. Under revision 7 the gate document **names a specific artifact and asserts that it discharged the obligation**. That converts an unstated assumption into a committed, named, falsifiable claim. Any later auditor can open that path and see in under a minute whether it is a review of the right commit reaching ACCEPT. A motivated chair must now write a specific falsehood into a durable artifact rather than merely omit a statement. That is a genuine change in the evidentiary record, and it is worth the field even though it is not a machine control.

### 3.2 Is the obligation binding, or merely written?

**Honest answer: written, plus one anchor. It is not self-executing, and this project has no mechanism that makes it so.** I looked for one and it does not exist:

1. **The freeze workflow** is triggered by frozen-byte changes. The four fields are not frozen bytes, so the workflow has no trigger routing an edit to them into itself. The obligation is a rule without a dispatcher.
2. **The freeze tests** go red if the fields are set — but `artifact_freeze_test.go` is unfrozen, and precedent `9a28d54` shows it being edited in the same unreviewed commit as the gate document.
3. **Reviewer prompts** are authored by the chair, per round. Nothing obliges a chair to convene a review of a commit the chair writes.
4. **The boundary document** is frozen, in the 16-file inventory, and every payload review reaches it. `:86` now places the obligation in **Authorization boundary**, addressing the reviewer directly: *"its independent reviewer must read the named closure record and rule that it is an accepting independent review of the closure, and must refuse a payload that sets those fields without one."*

So (4) does work — but only conditionally. It binds a reviewer **who is asked to review**. What no mechanism supplies is the guarantee that anyone is asked. The honest characterization: *the reviewer reads the boundary document, and the reviewer exists only if the chair convenes one.*

**Is that sufficient for §9's "independently reviewed and committed"? Yes — and I want to be precise about why, because this is where I could most easily have manufactured a finding.**

§9's post-closure gate reads: *"no subsequent validation execution may be **authorized** until a closure record … has been independently reviewed and committed."* This is a rule addressed to the project about authorization. It does not specify a runner check, and it could not: whether a document is an accepting independent review is not a machine-decidable property, as rounds 3, 4 and 5 established at the cost of three parser revisions. Freeze obligation 5 asks for something narrower and machine-shaped — *"the post-closure audit-before-retry gate **encoded in the frozen gate rules**"* — and that is delivered and verified.

So the boundary document's load-bearing sentence — *"this is the step where section 9's requirement … is actually discharged"* — is a **correct reading of the specification**, not an overclaim. §9 is discharged by a human act; the runner's job is to make its absence visible. That is exactly what the payload says it does, and exactly what it does.

**Is §9 therefore "discharged by convention"? Substantially, yes — and the payload says so.** I record that plainly because it is the honest answer, not because it is a defect. The difference from revision 6 is real: revision 6 asserted a false fact about an existing control and rested the deletion on it; revision 7 states a true rule, admits the runner cannot enforce it (both explicit in frozen bytes), acknowledges the contrary precedent rather than concealing it (*"notwithstanding that other gate-state fields in the same block have been set that way before"*), and leaves a named artifact behind. **The payload nowhere claims a mechanism it lacks.** Applying the prompt's own test — is the central claim unsupported? — the answer is no.

### 3.3 Does the fourth field create a regress?

Yes, and it is correctly left open rather than chased. The gate now depends on "the review that blessed the fields", and that review is itself unverified — with strictly *less* machinery than the closure record has: no digest pin, no verdict field, no commit field, no distinctness check.

**A P3-grade weakness, and properly a residual now.** Three reasons it is not a P2: any verification chain terminates in a human act, and adding a verdict field moves the regress one step rather than closing it; the exposure is *the same* future unreviewed commit review 6 identified, opening no new channel, and residual 1 discloses the human dependency in terms; and the payload never claims the field verifies the review's content — it says the opposite, twice, in frozen bytes.

But it is a genuine weakness, because **this project already has the right pattern three lines away and did not use it.** `frozenFileSet` carries `review_record` + `review_record_commit` + `review_record_raw_sha256` + `review_verdict` for every frozen set, machine-checked by `verifyRawGitFileDigest`. The fourth gate field carries none of the three companions. Graded as N1.

### 3.4 Does the un-frozen status of the two files undermine the repair?

**It limits it; it does not undermine it, and it is correctly a separate concern.** The obligation's design accepts that enforcement is human, so a chair-editable tripwire contradicts nothing the payload claims. What it does mean is that *"may be set only by a payload commit that is independently reviewed under this freeze workflow"* describes a routing the workflow does not automatically perform; a reader could infer that touching these fields *engages* the cycle, and it does not.

Fixing it here would be wrong. Bringing `pre-validation-artifacts.json` into a frozen block that it itself defines is circular, and bringing `artifact_freeze_test.go` in is a structural change to the freeze machinery affecting every block — far outside a payload scoped to §9. It belongs in its own cycle and should be recorded as a residual so it is not lost. It is not a reason to refuse this payload.

---

## 4. Area B — the replacement's correctness

`verifyClosureFieldsReview` mirrors `verifyClosureReviewRecord` exactly, minus the digest comparison. Attacked both, measured where measurable: paths outside `docs/reviews/`, traversal, absolute paths, backslashes, directories (both `docs/reviews` and `docs/reviews/.`), absent files, empty files, case differences, doubled separators, `foo..md`, and the empty-field case — all behave as the safe direction requires, and the measured results are in §3.1.

**The symlink case is the one asymmetry with teeth.** For the closure record it is not an escape — the bytes must still reproduce the pin. For the fields-review there is no pin, so a symlink under `docs/reviews/` pointing anywhere non-empty satisfies the field. It requires committing a symlink, which is conspicuous, and it is subsumed by the fact that any real file satisfies the field anyway. Not separately graded.

**Is the digest asymmetry deliberate and defensible, or an oversight?** **Defensible on ordering, and undocumented.** Under the freeze workflow the review record is committed *after* the payload commit it reviews, so the payload commit setting the closure fields cannot pin its own review's digest — the bytes do not exist yet. The closure record predates the commit that names it and can be pinned. That is sound and I accept it.

Two things follow. The reason is stated **nowhere** — not in the code comment, the boundary document, or the note — and an unexplained asymmetry in a document this heavily reviewed invites an eighth reviewer to raise it as I nearly did. And the ordering argument does *not* extend to the refreeze: the refreeze commit follows the review, has its bytes, and already writes `review_record_commit` / `review_record_raw_sha256` / `review_verdict` for the runner block. So the asymmetry is a reason to defer the pin, not to omit it permanently. Both folded into N1.

**Can an invalid closure record satisfy the gate? Only through the human channel §3 describes, never through the code. Can a valid one be refused? Only by an authoring error in the gate document, every one of which fails closed and surfaces immediately.**

**Test coverage of the new code is real**, not nominal: `TestPostClosureGateRequiresTheFieldsReviewToBeAResolvableRecord` pins outside-`docs/reviews`, traversal, absent and empty; `TestPostClosureGateRefusesValidationAfterANonCompletedExecution` gains an "unnamed fields review" case; both existing gate tests were widened to five parameters. The directory case is untested for either function — a coverage gap, not a defect, and `IsRegular()` is shared.

---

## 5. Area C — nothing regressed

- **Duplicated doc comment: fixed.** `grep -c` returns **1**. Review 6's N2 closed.
- **Four closure fields pinned empty** by `verifyPostClosureGateState`, with the new field correctly tagged and matching `validation_gate.go`.
- **Rule-text substrings match**: all eight verified present, including the two new ones, case-exact.
- **Inventory:** 16 files, all present, all digests match, set digest reproduces.
- **`accepted_residuals` now carries the full operative list.** 16 entries matching review 6 §9 one-for-one. **The `writeJSON` atomicity residual is present at entry 12** — reviews 4, 5 and 6 each flagged its absence and no revision fixed it; this one does. That three-round-old complaint is closed.
- **Deleted parser symbols:** zero occurrences in any `.go` file.
- **Unrelated behaviour untouched.** `execution_resilience.go`, `process_events.go`, `scheduled_resume.go`, `scheduled_run.go`, `types.go` and `main.go` are **byte-identical to `a368433`** — the revision review 5 passed with 0 P1 / 0 P2 — verified by an empty `git diff --stat`.

---

## 6. Area D — the whole payload

**Method, stated honestly.** The §9 enforcement machinery is byte-identical to the revision five reviewers attacked in depth. I re-derived the post-closure gate completely and measured its acceptance surface; verified the freeze document, inventory, protocol registration and gate flags directly; and did **not** re-derive the resume-attestation, process-event-log and partial-summary logic from scratch, relying on the empty diff plus prior records. That is the limit of my coverage.

Confirmed: every §9 condition located in code with both directions answered, unchanged from what reviews 2–6 verified; partial summary outcome-free with the depth-2 guard intact; process-event log honest about the sending principal; **post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure** — `validation_execution_status` is `"indeterminate"`, all four closure fields empty, so `requirePostClosureAudit` refuses today; `protocol.json` carries **8** blockers and `amendment.execution_resilience` in full; **freeze obligation 5 is discharged**; authoring exemption scoped and pinned; both gate flags false.

**Records against implementation — the note's present-tense claims.** This class produced findings in three consecutive rounds, so I checked line by line. Review 6's N3 and N4 are genuinely fixed: section 4 is rewritten to the four-field rule and is accurate; the refreeze recipe says **24 files** in both steps with the deleted names struck, and I verified 21 + 3 = 24 against the live block. **But two sentences carrying the exact claim this revision exists to remove survive verbatim** — N2 below.

---

## 7. Disposition of every prior review's P1 and P2

**Review 1:** P1-1, P1-2, P2-1, P2-2, P2-3, P2-5, P2-6 — **all still fixed** (files byte-identical to `a368433`). **P2-4 still fixed and strengthened this round.**

**Review 2:** P2-N1, P2-N2 — **still fixed.**

**Review 3:** P1-N1 — **retired** by the deletion; component absent from every `.go` file.

**Review 4:** N1, N2 — **retired** by the deletion.

**Review 5:** N1 — **retired** by the deletion.

**Review 6:** **N1 — FIXED this revision** (all three frozen locations imperative, obligation in **Authorization boundary**, fourth field anchoring it; see §3 for what the fix does and does not achieve). N2 — **fixed**, exactly one comment occurrence. N3 — **fixed**, section 4 accurate. N4 — **fixed**, 24 files, arithmetic verified.

**No prior P1 or P2 has regressed.** Every one is still fixed, retired by the deletion, or fixed this round.

---

## 8. New findings — 0 P0 · 0 P1 · 0 P2 · 4 P3

### N1 — P3. The fields-review anchor is weaker than this project's own pattern, and the asymmetry is unexplained

**Files:** `validation_gate.go:127-144`; `pre-validation-artifacts.json` `gates.validation_execution_closure_fields_payload_review`.

`verifyClosureFieldsReview` checks only prefix, absence of `..`, regular file, non-empty. Measured (§3.1): the field is satisfied by this review's own prompt, by review 1 or review 6 (both REVISE), by any unrelated review, and by naming the closure record in both fields. Nothing verifies the record reviewed the payload commit that set the fields, reached ACCEPT, or is distinct from the closure record. Being unpinned, it can also be amended after the fact — the exact property the payload argues the digest pin uniquely supplies for the closure record. Three lines away, `frozenFileSet` maintains `review_record` + `review_record_commit` + `review_record_raw_sha256` + `review_verdict`, machine-checked by `verifyRawGitFileDigest`; the gate field carries none of the three companions, and no document explains why.

**Direction: permits an invalid resume** — but only through the same future unreviewed commit residual 1 already discloses, and only to an actor already willing to write a false claim into a durable artifact. It creates no new channel and refuses nothing valid.

**Fix.** Not in this payload. The ordering argument is sound and should be **stated** rather than left silent. The strengthening belongs to the future payload that actually sets these fields — by construction an independently reviewed commit whose reviewer is the person who must rule: add `…_payload_review_commit`, `…_payload_review_raw_sha256` and `…_payload_review_verdict`, verify them with the existing `verifyRawGitFileDigest` helper, and require the fields-review path to differ from the closure record path. Recorded as a residual.

### N2 — P3. The candidate note still asserts, twice and in the present tense, the exact claim this revision exists to remove

**File:** `artifacts/execution-resilience-candidate.md:158` and `:164`.

Line 158: *"…it is already made — by the independent reviewer of the refreeze that sets the gate fields."* Line 164: *"that judgment stays with the independent reviewer of the refreeze, which is where it already lived and where it is competent."*

Review 6's N1 named four locations; this note was one of them. Three were repaired. **This one was not**, and the note now contradicts itself six lines later at `:170`, which states in bold that no such reviewer exists and reproduces the archaeology proving it. Both sentences are present tense and unqualified, so they fail review 6's own historical-versus-present-tense test, notwithstanding that they sit in a section titled "What the fifth review changed".

**Direction: neither.** The note is documentation and is not among the 16 frozen files. It matters because the note is what a chair reads first when planning the refreeze, and these sentences transmit precisely the mental model — *the refreeze reviewer will catch it* — whose absence is the subject of this revision.

**Fix.** One clause each: mark both as revision 6's reasoning in the past tense, or strike them and let `:170` stand. A transcription-time correction; no cycle.

### N3 — P3. The note says the fields-review "must resolve the same way" as the digest-pinned closure record

**File:** `artifacts/execution-resilience-candidate.md:59`.

The antecedent of "the same way" includes the digest clause, and there is no pin for the fields-review. The boundary document gets this right — it scopes "reproduces the pinned digest" to the closure record alone — so the defect is confined to the note.

**Direction: neither.** Documentation, unfrozen. Graded because a reader taking it at face value credits the anchor with a strength it does not have.

**Fix.** *"…the named fields-review must resolve to such a file, but is not digest-pinned."*

### N4 — P3. The stated justification for the fourth field describes the failure mode it is least likely to encounter

**Files:** `validation_gate.go:129-131`; `validation-execution-boundary.md:80`.

Both say the runner "can refuse a gate that names none — which is what setting these fields outside a reviewed payload commit would leave behind." An empty fourth field is what an actor working from the *old* three-field schema leaves behind. An actor who wants the gate open cannot leave it empty — the runner refuses and their objective fails — so they will fill it, and anything fills it. The sentence presents a transitional, inattentive-only case as the characteristic one, and undersells the field's actual value: forcing a specific falsifiable claim onto the record.

**Direction: neither** — it is a justification, not a check, and the adjacent sentence correctly hedges. Graded because the last two rounds turned on claims frozen alongside correct code, and this is frozen prose that overstates its reach.

**Fix.** Replace the clause with what the field actually buys: *"…so the gate must name a specific record rather than leave the discharge unstated, and any later auditor can check that record against the commit it claims to review."*

---

## 9. Residual list — worded for transcription into the refreeze's `accepted_findings`

**The inventory's `accepted_residuals` array is now the operative list and has caught up** — 16 entries matching review 6 §9 one-for-one, including the `writeJSON` atomicity entry that reviews 4, 5 and 6 each flagged as missing. **Transcribe all 16 verbatim**, with entry 1 extended and two new entries added:

**Entry 1 — replace with:**

> "Residual: the runner forms no opinion about whether the named closure record is an accepting independent review. It verifies that the named path is a non-empty regular file under docs/reviews, free of traversal, whose LF-normalized digest reproduces the pin, and that the gate names the independent review of the payload commit that set those fields. Whether the closure record is an accepting independent review rests on that human step, which the boundary document binds to an independently reviewed payload commit and forbids to a chair gate patch. The fields-review field is an anchor, not a check: it is not digest-pinned, carries no commit or verdict, and is not required to differ from the closure record, so it is satisfied by any non-empty regular file under docs/reviews. Its value is that the discharge is named on the record and checkable by a later auditor, not that the runner verifies it. Deleting the verdict parser gave up machine refusal of a named record whose own prose reads REVISE, or which states no verdict at all; it gave up nothing against a document that merely lists the available verdicts, which the parser accepted."

**New entry — add:**

> "Residual: neither pre-validation-artifacts.json nor artifact_freeze_test.go appears in any frozen file set, so setting the four post-closure closure fields is not a frozen-byte change and does not itself trigger the close-payload-review-refreeze cycle. The obligation that only an independently reviewed payload commit may set them is written in the frozen gate rule, the runner comment and the boundary document's Authorization boundary section, and anchored by the fields-review field; it is not machine-routed. Bringing either file into a frozen block is a change to the freeze machinery affecting every block and is out of scope for this payload."

**New entry — add:**

> "Residual: the closure record is digest-pinned and its review is not, because under the freeze workflow a payload commit precedes its own review and cannot pin bytes that do not yet exist. The companion commit, digest and verdict fields are available to the future payload that actually sets the closure fields, whose independent reviewer is the person the obligation binds; they are deferred to it rather than omitted."

**On the "runner forms no opinion" residual specifically — review 6 made its acceptability conditional on N1 being fixed. I rule that the condition is met and the residual is now acceptable**, in the amended wording above. Review 6's objection was that the residual asserted a human obligation written nowhere binding, so accepting it would freeze the assumption rather than the control. That is no longer the case: the obligation is imperative in three frozen locations, restated in **Authorization boundary** where it addresses the reviewer, and anchored by a field the runner enforces. The amendment I require is only that the residual stop leaving the anchor's weakness to inference.

---

## 10. My judgment on each residual

**All 16 carried entries are correctly characterized and genuinely acceptable, and none should be promoted.** Six notes:

1. **Entry 1 is acceptable as amended.** Unamended it is acceptable but incomplete.
2. **The two new entries are the honest residue of area A.** Neither is a defect in this payload; both are structural facts the payload correctly declined to change. Recording them is what stops the next round rediscovering them as findings, which is how three of the last four rounds began.
3. **Entries 2–6 remain the load-bearing ones**, and all five are correctly characterized as custodian-trust limits no runner-side control can remove. §9 itself assigns these to the custodian's written attestation, so they are limitations of the specification, not the implementation.
4. **Entry 9 is the most consequential residual in the list** — it is what makes the post-closure gate the sole control on the retry channel, which is why seven rounds have been spent on this gate. Correctly characterized.
5. **Entry 12 (`writeJSON` atomicity) should be fixed eventually**, not because it is dangerous but because a temp-and-rename removes the class for a few lines and it has now been carried for five rounds. Acceptable to carry again.
6. **Entries 14–16 are pre-existing and unchanged**; `gocyclo -over 15` is clean.

---

## Verdict

**ACCEPT.** 0 P0 · 0 P1 · 0 P2 · 4 P3. The four P3s are text corrections (N2, N3, N4) and one deferral recorded as a residual (N1); none permits an invalid resume, none refuses a valid one, and none justifies an eighth round. They should be carried into the refreeze as accepted findings or fixed at transcription — N2 in particular, since the payload's own primary narrative document should not still assert the claim the payload removes.

**What I verified and what I did not.** I verified every pinned digest, the 16-file inventory and its set digest, all nine quality gates including the race detector run by me, the single deliberate red, the four-field gate encoding, the frozen-block status of both unfrozen files by direct search, the acceptance surface of `verifyClosureFieldsReview` by executing a verbatim replica against real paths, the eight rule-text substrings, the 16-entry residual array against review 6 §9, the world-build and grader identities against decision 0025, the blocker count and both gate flags, and that the six §9 machinery files are byte-identical to the revision review 5 passed. I did **not** re-derive the resume-attestation, process-event-log or partial-summary logic from first principles.

---

## Chair transcription note

Transcribed verbatim from the seventh independent reviewer's record. Before transcription the payload author independently reproduced the findings: the duplicated doc comment occurs exactly once (review 6's N2 closed); the residual array carries 16 entries; and both N2 sentences survive at `execution-resilience-candidate.md:158` and `:164`, contradicting `:170` six lines later.

**Disposition of the four P3s at transcription.** N2 and N3 lie in the candidate note, which is **not** among the 16 frozen inventory files, so both are corrected in the refreeze commit as the reviewer sanctioned. **N4 is not corrected**: it lies in `validation_gate.go` and `validation-execution-boundary.md`, both frozen inventory files whose digests this review pins, so editing them would invalidate the record that accepts them. N4 and N1 are carried into the refreeze as accepted findings, with N4's suggested wording recorded for the future payload that sets the closure fields.

This is the **acceptance record**. The focused refreeze follows, then decision 0026.
