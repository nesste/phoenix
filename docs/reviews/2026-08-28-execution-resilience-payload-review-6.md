# Sixth independent review — protocol-v5 section 9 execution-resilience payload, revision 6

- **Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author it and I am not any of the five reviewers who produced the prior records.
- **Date:** 2026-08-28
- **Payload commit:** `b375d432df1facf29a8e3e66f065b19453a8ede4`
- **Verdict:** **REVISE** — 0 P0 · 0 P1 · **1 P2** · 3 P3

**The deletion is right and should stand.** I agree with the fifth review's § 9 recommendation and I would not reinstate the parser, not even in narrow canonical form. My P2 is not about the code that was removed; it is about the justification frozen alongside it. The payload discharges § 9's "independently reviewed and committed" requirement entirely onto a workflow step — "the independent reviewer of the refreeze that sets these fields" — that **this project does not have and has never had**. I verified that by git archaeology, not by reading intent. The fix is text plus one written obligation, no parser, and it should make round seven very small.

No model was run, no private grade obtained, no validation outcome observed. Both gates remain closed.

---

## 1. Pinned inputs — every digest I computed

`git rev-parse HEAD` = `b375d432df1facf29a8e3e66f065b19453a8ede4`. Ancestry confirmed exactly as specified, twelve commits back to `2cf99331…` (decision 0025).

**Working tree:** exactly one untracked file, this prompt. Nothing else modified, staged, or untracked. **No P0.**

| Artifact | Rule | Computed | Match |
| --- | --- | --- | --- |
| `artifacts/gate-1a-execution-resilience-candidate.json` | raw | `sha256:3874e629dc359dc05b5b6607d61347547ee75e2388d7808bee2fa95061298970` | yes |
| `artifacts/execution-resilience-candidate.md` | raw | `sha256:ac478f8593db5e7beffebe8e037bce725b16af068b0196493d4b94a349e3d245` | yes |
| `runner/validation_gate.go` | LF | `sha256:741e19b0cec0d83d8348f931dec8d371e4de61217e53ff4cf762f63200ddd8c9` | yes |
| `runner/execution_resilience_test.go` | LF | `sha256:e6b2cc06f707ae42c2e5c3c459bc6c4d75916ee9480efd080667d572afc2d837` | yes |
| `runner/artifact_freeze_test.go` | LF | `sha256:0b5b3662bea710600abd49518c6d1250f4f6d699a0fb300ec5f81176d81bbb8a` | yes |
| `pre-validation-artifacts.json` | LF | `sha256:40f163717c4ead7ac237e931d50b96c2766166e003c0cdf5016ac7bc2b33ef18` | yes |
| `validation-execution-boundary.md` | LF | `sha256:a81514b8f4dea6678362df0ce4eb769378176e64f5307b02f789326f19ba51f5` | yes |

**Inventory:** 16 files. I recomputed all sixteen LF-normalized digests against the working tree — **zero mismatches**. **Ordinal-path set digest** `sha256:1ba17cd292b53122a483aab037a5cf03a0108991fd1df55aca72faf26eb1e834` reproduced.

One note for future reviewers, not a finding. The `artifact_set_digest_algorithm` string is ambiguous about the `sha256:` prefix. The line **includes** it. Hashing the bare hex instead yields `sha256:1db7b695a724301d19ddee67c67655ef7ac7c7926ce67e9fc056282fa5fa7688`, which is what I got on my first attempt. The convention is pre-existing and consistent with `verifyCandidateSetDigest`; I record both so a seventh reviewer does not lose the same ten minutes.

**Experiment pin** (§ 3.1): `sha256:49e7ee4c00771b204b3b780101229d5fd8abb7b7d0a28e4a6f3fc8b91a8f8663`.

---

## 2. Quality gates — all run by me at this commit

| Gate | Result |
| --- | --- |
| `go test ./...` | **1 failure, the expected one.** `TestPreValidationFreezeMatchesAcceptedCandidates`. Every other package `ok`. |
| corpusctl suite | `ok`, clean |
| `go vet ./...` | Clean |
| staticcheck v0.7.0 | Clean |
| `gocyclo -over 15` | Clean |
| `dupl -t 100` | Clean |
| `make validate-spec` | All 9 examples valid |
| `make validate-authoring` | Valid |
| **`go test -race`** | **I ran this myself** on WSL2. **No `DATA RACE`.** Same single freeze red. 42.1s. |

**Exactly one red, and it is the declared one.** World-build `sha256:425bab1c…ffdae4` and the grader identity are unmoved from decision 0025; the local-artifact blocks are untouched as the note promises.

---

## 3. Area A — did the deletion give up something the gate needed?

### 3.1 What was actually lost, measured

I extracted the shipped parser verbatim from `a368433:validation_gate.go` into a standalone module and ran it head-to-head against the digest-pin rule on five documents. Execution, not reading.

```
pinned at review time: sha256:49e7ee4c00771b204b3b780101229d5fd8abb7b7d0a28e4a6f3fc8b91a8f8663

1. the record exactly as the reviewer accepted it
    old parser : PASS        digest pin : PASS
2. AMENDED after the review that blessed it
    old parser : PASS        digest pin : REFUSE
3. a REVISE record, pinned by whoever sets the fields
    old parser : REFUSE      digest pin : PASS
4. an evaluator PROMPT named as the record, freshly pinned
    old parser : PASS        digest pin : PASS
5. an unreviewed chair note, freshly pinned
    old parser : REFUSE      digest pin : PASS
```

**Case 2 confirms the payload's strongest claim.** The pin is genuinely, strictly stronger on the amended-after-review case, exactly as the commit message and boundary document assert. Verified by execution.

**Case 4 settles the argument.** The document most likely to be named by mistake — a review *prompt* — **passed the parser too**. On the payload's own headline threat the deleted component was already worthless. Nothing is lost there.

**Cases 3 and 5 are the real residue.** The parser refused a document whose own prose reads REVISE, and one with no verdict statement at all. The pin accepts both. That is precisely and only what the deletion gave up: two forms of *inattention* by whoever sets the fields.

### 3.2 The load-bearing claim, tested against how this project actually works

The payload rests cases 3 and 5 on a single sentence, frozen in three places (`validation_gate.go:103-105`, the gate document's `validation_execution_closure_gate`, and `validation-execution-boundary.md:78`):

> Whether those bytes constitute an accepting independent review is a human judgment, and under the freeze workflow it is already made by the independent reviewer of the refreeze that sets these fields.

**I tested this obligation and it does not exist.**

**(a) There is no reviewer of a refreeze.** The workflow is payload commit → **independent review of the payload** → focused refreeze → decision document. The review targets the payload commit. I checked every refreeze commit this project has produced: decision 0024's `1f80d7b` and decision 0025's `e07a3fc` carry **no independent review record**. No document in `docs/reviews/` targets a refreeze commit. The reviewer named in the sentence is not a role this project instantiates.

**(b) Read charitably — "the reviewer in the cycle that sets these fields" — the claim still fails, on this project's own precedent.** Gate-state fields in this block are not set by review cycles. They are set by chair gate patches. Decision 0022 states it plainly: *"This decision record and the focused gate-state patch are committed atomically."* Commit `9a28d547aae9216b8c2b27d30f5fa531aec7fbee` is that patch: four files, no independent review, and it rewrote `pre-validation-artifacts.json`'s `gates` block. Decision 0023 set `validation_execution_status: "indeterminate"` the same way.

**(c) The freeze test that holds the line is itself unpinned and chair-editable.** `verifyPostClosureGateState` requires all three closure fields empty. I searched `pre-validation-artifacts.json` for `artifact_freeze_test.go`: **zero occurrences**. It is in no frozen block, in neither the current 21-file `scheduled_runner` set nor its successor. Neither is `pre-validation-artifacts.json` itself. And commit `9a28d54` did exactly this edit — including flipping `gates.MayOpenValidation ||` to `!gates.MayOpenValidation ||`, retargeting which gate state the test demands — with no review.

**(d) The payload's own note contradicts the claim.** Its residual list says the gate's application *"depends on **the chair** setting that field at closure."* The note says chair; the frozen bytes say independent reviewer.

**The concrete attack, end to end.** A future `OPEN_GATE_1A` decision, following decision 0022's template exactly, commits atomically: (i) `may_open_validation: true`; (ii) the three closure fields naming any document and its LF digest — a REVISE closure review, a chair note, an evaluator prompt; (iii) the two-line relaxation of `verifyPostClosureGateState`. Suite green. `requireValidationGate` permits the next validation execution. **No step in this project's written process requires an independent reviewer to look at that commit.**

**Is the residual exposure real, or absorbed by human review as the note and boundary document claim?** It is real and not absorbed, because the absorbing step is not written down anywhere binding. § 9 requires the closure record to have "been independently reviewed and committed." After this payload that requirement is discharged by **nothing**: not by the runner, which now has no opinion about what the document is, and not by any obligation on the commit that sets the fields. The prompt asked whether the obligation is written anywhere binding or merely assumed. **It is assumed.** That is N1.

To be fair to the payload: this was a self-binding control on the chair, and against a *motivated* chair the parser was worth nothing — the chair writes the closure document and can put any verdict line in it. The parser resisted inattention only, and case 4 shows it failed on the most likely inattention. The loss is smaller than the frozen prose implies. But it is not zero, and the compensating control the frozen prose names does not exist.

### 3.3 Would a narrow canonical-form verdict check recover it? — No, and I would not want one

The narrowest defensible form is: the record must contain, at the start of a line, the exact byte sequence `**Verdict: ACCEPT**` — one mandated form, no cutsets, no qualifier set, no block tracking, no corpus. Roughly five lines and genuinely closed.

**I do not want it.** Three reasons, in order of weight:

1. **It buys cases 3 and 5 only, and both are inattention cases that a written review obligation covers completely and for free.** Fixing N1 dominates it.
2. **It does not close case 4**, the actual threat. A closure-review prompt drafted from this project's own template would carry the mandated form in its instruction block. I would then be prescribing exactly the heuristic arms race that consumed rounds 3, 4 and 5.
3. **It reintroduces a frozen refusal channel with no repair path.** The boundary document's paragraph stating accepted verdict forms was deleted here. Reinstating a form check means reinstating that spec, and a future closure record written outside it is refused at the point of maximum schedule pressure, repairable only by a full cycle. That was review 3's P1-N1 in its false-refusal direction.

**The deletion is simply right.** The runner should check identity; the judgment should sit with a human. What the payload got wrong is not deleting the parser — it is asserting a human step exists when it does not, instead of creating it. Text and one obligation sentence, not code.

---

## 4. Area B — the replacement's own correctness, both directions

| Attack | Result |
| --- | --- |
| Directory named as the record | Refused — `IsRegular()` |
| Symlink to a file elsewhere | `os.Stat` follows it, but the bytes must still reproduce the pin. Not an escape. |
| `docs/reviews/../README.md` | Refused — `Contains(review, "..")`. Pinned by test. |
| Leading `/`, or backslash separators | Refused — fails the `docs/reviews/` prefix. Safe direction. |
| `DOCS/REVIEWS/x.md` on a case-insensitive FS | Refused by the prefix check even though the file resolves. Refuses a valid one; the path is authored, so conventional. |
| Digest in wrong case / no `sha256:` prefix | Refused — exact compare, identical to every other pin in the project. |
| CRLF vs LF in the record | **Handled.** `digestLFNormalizedFile` normalizes, so a Windows checkout of an LF-pinned record still matches. |
| Raw digest pinned instead of LF-normalized | Refused if the file has any CRLF. Field name is unambiguous. Acceptable. |
| Empty file, or a pin equal to the empty digest | Refused by the size check before the digest is taken. |
| Empty pin with `ACCEPT` verdict | Refused in `requirePostClosureAudit`. |

`digestLFNormalizedFile` is **the same rule the rest of the project uses**, shared with the freeze tests and the Arm B document path. It reproduced my independent Python implementation on all seven pinned files.

Two observations deliberately **not** graded as findings: the `.md` suffix requirement was dropped, and under a byte pin the extension carries no information the digest does not; and `Contains(review, "..")` would refuse a record named `foo..md`, a safe direction with no such record. I mention both so a seventh reviewer does not raise them.

**Can an invalid closure record satisfy the gate? Yes — but only via N1's channel, never via the code.** **Can a valid one be refused? Only by an authoring error in the gate document**, all of which fail closed and are visible immediately.

---

## 5. Area C — nothing else regressed

- **Dangling symbols:** none. The seven removed functions appear in **no** `.go` file. Surviving mentions are in the inventory's historical-narrative fields, where they correctly describe the deletion.
- **Deleted files gone** from the working tree and the inventory.
- **Inventory:** 16 files, all digests match, set digest reproduces.
- **`verifyPostClosureGateState`** pins all **three** closure fields empty, correctly extended. All six rule-text substrings present in the gate document. Checked programmatically.
- **Candidate guard** tracks the new base `950df934…` at 16 files and lives outside the frozen set.
- **Unrelated behaviour:** the commit touches ten files. `execution_resilience.go`, `process_events.go`, `scheduled_resume.go`, `scheduled_run.go`, `types.go` and `main.go` are **untouched**, so the § 9 enforcement machinery is byte-identical to what review 5 passed with 0 P0 / 0 P1.

---

## 6. Area D — the whole payload, as if fresh

Stated honestly about method: the § 9 enforcement code is unchanged from revision 5, which five reviewers have attacked in depth and which carries no outstanding P1 or P2. I re-derived the post-closure gate completely, spot-checked the rest against the records, and confirmed the suite green but for the one red. I did **not** re-derive the resume-attestation, process-event-log and partial-summary logic from scratch; I relied on the unchanged diff plus the prior records, and I flag that as the limit of my coverage.

Confirmed at this commit: partial summary outcome-free in both directions with the depth-2 guard intact; process-event log honest about the sending principal; post-closure gate encoded in the frozen gate rules and applying retroactively to decision 0023 — `protocol.json` carries **8** blockers including the closure-record entry, and `amendment.execution_resilience` is present; authoring exemption scoped and pinned.

**Records against implementation.** The note's six-revision narrative is accurate about the past. Its **present-tense** claims are not, in two places — N3 and N4. The commit message asserts it "retires N5 by rewriting the stale paragraph"; it rewrote one paragraph and left two others, one of which is the note's own specification of the very component that changed.

---

## 7. Disposition of every prior review's P1 and P2

**Review 1:** P1-1, P1-2, P2-1, P2-2, P2-3, P2-5, P2-6 — **all still fixed**, code untouched. **P2-4** (closure-record check is existence only) — **still fixed, by a materially better mechanism**: existence alone never satisfies the gate; a directory, a path outside `docs/reviews/`, a traversal, an empty file and a byte mismatch all refuse. The *prose* half of the round-3/4/5 remainder is **retired by deletion**, not regressed — but see N1 for what replaced its justification.

**Review 2:** P2-N1, P2-N2 — **still fixed.**

**Review 3:** **P1-N1 RETIRED BY THE DELETION.** The component no longer exists. My case 4 confirms independently that its false-accept direction was still live in the shipped parser.

**Review 4:** N1 and N2 (both P2) — **RETIRED BY THE DELETION.**

**Review 5:** N1 (P2) — **RETIRED BY THE DELETION.**

**No prior P1 or P2 has regressed.** Every one is either still fixed or retired by removing the component it concerned. Of review 5's P3s: N2 and N4 are retired with the parser and corpus; **N3** is retired as written, but its *class* — a frozen document claiming a control more strongly than it exists — recurs as my N1 in stronger form. **N5 is NOT retired**: the named paragraph was fixed and two more of the same class were left. See N3 and N4.

---

## 8. New findings — 0 P0 · 0 P1 · 1 P2 · 3 P3

### N1 — P2. The compensating control the frozen bytes name does not exist as a required step

**Files:** `validation_gate.go:103-105`; `validation-execution-boundary.md:78`; `pre-validation-artifacts.json` `gates.validation_execution_closure_gate` (final clause); `execution-resilience-candidate.md:158`.

**The failure it permits.** All three frozen documents state the judgment is "already made by the independent reviewer of the refreeze that sets these fields." No such reviewer exists (§ 3.2a), gate-state fields in this block are set by unreviewed chair patches (§ 3.2b), and the freeze test holding them empty is itself unpinned and chair-editable (§ 3.2c). A future `OPEN_GATE_1A` decision can set the three closure fields to any document, relax `verifyPostClosureGateState`, and open validation in one commit that no written process obliges anyone to review. § 9's "independently reviewed and committed" requirement is then discharged by nothing at all.

**Direction: permits an invalid resume** — a subsequent validation execution after a non-completed one whose closure record was never independently reviewed, the exact engineered-indeterminate retry channel the gate exists to bound. Not exploitable at this commit: fields empty, `may_open_validation` false, `requirePostClosureAudit` unreachable. It requires a future commit. That is why it is P2 and not P1.

**Fix — text plus one obligation, no code:**

1. Replace the descriptive sentence in all three frozen locations with an imperative one: *"The three closure fields may be set only by a payload commit that is independently reviewed under the freeze workflow, and that reviewer must read the named record and rule that it is an accepting independent review of the closure. A chair gate patch must not set them."*
2. Add the same obligation to the boundary document's **Authorization boundary** section, where it binds, rather than only to the explanatory paragraph.
3. Recommended, cheap, machine-checkable: have the commit that sets the fields also record the review record that blessed them, and keep `verifyPostClosureGateState` asserting that field non-empty whenever the others are set. That gives the workflow claim a frozen anchor instead of an assumption.
4. Correct the candidate note's residual, which already says "chair", to match.

### N2 — P3. A seven-line doc comment is duplicated verbatim in a file about to be frozen

**File:** `validation_gate.go:85-91` and `:92-98`.

The block beginning `// requirePostClosureAudit encodes the protocol-v5 section 9 post-closure gate.` appears **twice, identically**. The diff shows the payload added a full copy of the pre-existing comment when it meant to append only the new paragraph; revision `a368433` has exactly one occurrence, HEAD has two. `gofmt`, `go vet` and staticcheck are all blind to it.

**Direction: neither.** No behavioural effect. It matters because refreeze would freeze it permanently into the accepted byte set, and because it is a fresh regression introduced by this revision.

**Fix:** delete lines 92-98.

### N3 — P3. The note's own specification of the post-closure gate still describes the deleted parser

**File:** `execution-resilience-candidate.md:59` (section "### 4. Post-closure gate").

Present tense, four claims, all now false: the record must be *"a non-empty regular **markdown** file … whose own text states an ACCEPT verdict"*; *"a record whose verdict contradicts the gate document"* refuses (it no longer does — case 3); *"the gate document gains **those two fields**"* (three); *"asserts … that **both fields** stay empty"* (three). This is the note's primary description of the one component the revision changed, and the paragraph a chair or seventh reviewer reads first.

**Direction: neither.** Documentation, and the note is not among the 16 frozen files. But this is review 5's N5 class, and the commit message claims N5 retired.

**Fix:** rewrite section 4 to the current rule.

### N4 — P3. The refreeze recipe instructs the chair to pin two files that no longer exist

**File:** `execution-resilience-candidate.md:182` and `:183`.

Step 1 says re-pin `scheduled_runner` at *"26 files (the existing 21 plus … `closure_verdict_corpus_test.go`, and the pinned corpus expectations JSON)"*. The last two are deleted by this payload. Step 2 repeats "file count 26". I verified the arithmetic against the live 21-file block: the correct successor is **24**.

**Direction: refuses a valid one** — followed literally it produces a refreeze that cannot go green, self-correcting but wasting a cycle where this project has least slack.

**Fix:** 24 files; strike the two deleted names from both lines.

---

## 9. Residual list — **not operative at this verdict**

**Retired by the deletion — must NOT be carried forward** (area E; the prompt's expectation confirmed exactly): review 5's block-context/fence-state residual; its closed-qualifier-set residual; its pinned-corpus-completeness residual. All three name a component that no longer exists.

**Added by the deletion** — one new entry:

> "Residual: the runner forms no opinion about whether the named closure record is an accepting independent review. It verifies that the named path is a non-empty regular file under docs/reviews, free of traversal, whose LF-normalized digest reproduces the pin. Whether the document is an accepting independent review of the closure rests entirely on the human step that sets the three gate fields, which must be an independently reviewed payload commit under the freeze workflow and not a chair gate patch. Deleting the verdict parser gave up machine refusal of a named document whose own prose reads REVISE, or which states no verdict at all; it gave up nothing against a document that merely lists the available verdicts, which the parser accepted."

**Carried forward unchanged from review 5 § 8:** signal-sender principal not exposed; attestation *claims* digest identity; process-event log is unsigned custodian-controlled local state; 72-hour anchor against a custodian-controlled clock; window erosion bounded by one pairing-key group; `requirePostClosureAudit` ordering; authoring exemption; nothing writes `validation_execution_status` and a wiped directory presents as fresh; committed log digest never covers the terminal entry; five non-outcome identity fields; `os.WriteFile` atomicity; the start-event-without-checkpoint state; carried P2-3; carried pre-existing `validateExternalGrade` P3; P2-1/P2-5 retirement with `runScheduledCases` and `validateExternalGrade` at gocyclo 14.

---

## 10. My judgment on each surviving residual

**All carried entries are correctly characterized and genuinely acceptable, and none should be promoted.** Three notes:

1. **The `writeJSON` atomicity residual is still absent from the inventory's `accepted_residuals` array**, exactly as reviews 4 and 5 both flagged and neither revision fixed. It is now three rounds old. **Transcribe from this § 9, not from the array** — the array is not the operative list and has never caught up.
2. **The new "runner forms no opinion" residual is only acceptable once N1 is fixed.** Its second sentence asserts the human obligation; today that obligation is not written anywhere binding, which is precisely N1. Accepting it against the current text would freeze the assumption rather than the control. This is the same conditional review 4 attached to its N2 and review 5 invoked — I am invoking it in turn.
3. **Review 5's judgment that the deletion was the right call is correct, and I endorse it** against the alternative of a sixth parser round. My disagreement with the payload is narrow and entirely about what was frozen *alongside* the deletion. Review 5's § 9 asserted the same workflow claim the payload then froze; that assertion was the one thing in an otherwise excellent record that nobody had tested. I tested it, and it does not hold. Fixing that is a text change plus one sentence of written obligation, and it needs no code at all — which is why round seven should be short.

---

## Chair transcription note

Transcribed verbatim from the sixth independent reviewer's record. Before transcription the payload author independently reproduced all four findings. N2: the doc comment occurs once at `a368433` and twice at HEAD — a regression introduced by the deletion commit, invisible to gofmt, vet and staticcheck. N1(a): no document in `docs/reviews/` references either refreeze commit `1f80d7b` or `e07a3fc`, so no refreeze has ever been independently reviewed. N1(b): commit `9a28d54` is exactly as described — a chair gate patch touching `pre-validation-artifacts.json` and `artifact_freeze_test.go` in one unreviewed commit. N1(c): `artifact_freeze_test.go` appears zero times in `pre-validation-artifacts.json`, so it sits in no frozen block. The verdict is **REVISE**. § 9 is conditional and not operative.
