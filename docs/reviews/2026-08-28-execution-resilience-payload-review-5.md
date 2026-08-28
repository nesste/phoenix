# Fifth independent review — protocol-v5 § 9 execution-resilience payload

- **Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author the payload, and I am not any of the four reviewers who produced the prior records. Fifth review.
- **Date:** 2026-08-28
- **Payload commit:** `a368433089a857dbec4d72ee74123a4ee8743f47` (`experiments: fix the section 9 payload findings from the fourth review`)
- **Verdict:** `REVISE`

No model was run, no private grade obtained, no validation or held-out outcome observed. `may_open_validation` and `may_open_held_out` were `false` at the start and end of this review and I did not touch them. I made no change of any kind to the repository: my only writes were to the session scratchpad (a standalone extraction of the parser, a patched variant of it, a 43-case attack corpus) and to this record.

---

## 1. Pinned inputs — verified

| Input | Expected | Computed | |
| --- | --- | --- | --- |
| `git rev-parse HEAD` | `a368433089a857dbec4d72ee74123a4ee8743f47` | identical | match |
| Ancestry | `fec6d554…` → `70ff7875…` → `6f54e8ca…` → `60c6aee5…` → `9f53fa73…` → `e92eaacf…` → `46030d4f…` → `7a1ea425…` → `2cf99331…` | exact, in order | match |
| Working tree | exactly one untracked file, this prompt | that single line, nothing else | match |

**No P0 on tree cleanliness.**

### Digests I computed

| Artifact | Kind | Computed | |
| --- | --- | --- | --- |
| `artifacts/gate-1a-execution-resilience-candidate.json` | raw | `2309d017d58e49e4dceb9c1dec37823e0902200ac31e98813844913a3da61bad` | match |
| `artifacts/execution-resilience-candidate.md` | raw | `b69c7ff802a10d124493ac1c14603a2a601bd11453cc6258e71d1919a9a42d0b` | match |
| `runner/validation_gate.go` | LF | `5f97d6e13b8fcad5f3c32155eec4359f72aba2254da5856e3a036ededf6061c4` | match |
| `runner/closure_verdict_corpus_test.go` | LF | `a400c96336cef8814dd9c8d1548aedc1733e2910b8979c44112363b451c18c7d` | match |
| `artifacts/closure-verdict-corpus-expectations.json` | LF | `b761d71d9279775ca659ef2e1785ef3148ef2329dfa4cf4bec0cffba7d992694` | match |
| `validation-execution-boundary.md` | LF | `778926ba22b888964f8d40509d08684eaa9225c02076dd0be7aeb3d47f643351` | match |
| Inventory ordinal-path set digest | — | `44192b3c80911ee22489b0537f30dbddaea5003576b93fa8543e051f59995f67` | match |
| Inventory file count | — | **18** | match |

All **18** per-file digests recomputed independently: **zero mismatches**.

Digests I did **not** recompute from source: the world-build and grader digests. I verified them indirectly and sufficiently: `git diff --name-only fec6d554…HEAD` over `experiments/frontier-v1/corpusctl`, `go.mod`, `go.sum`, `experiments/frontier-v1/schema`, `cmd/phoenix` and `internal` returns **empty**, so neither digest's input set was touched; and both pinning tests are green.

---

## 2. Quality gates — all run by me

| Gate | Result |
| --- | --- |
| `go test ./...` | one failure, the expected one; every other package `ok` |
| corpusctl suite | **pass** |
| `go vet ./...` | **pass**, no output |
| `staticcheck v0.7.0 ./...` | **pass**, no findings |
| `gocyclo -over 15` | **pass**, no output |
| `dupl -plumbing -t 100` | **pass**, no output |
| `make validate-spec` | **pass** — 9 examples valid |
| `make validate-authoring` | **pass** |
| `go test -race` on WSL2 | **ran it myself**; no data race; only the expected freeze red |

The single red, on both hosts, is `TestPreValidationFreezeMatchesAcceptedCandidates` — the deliberate payload-commit red. **No other failure anywhere.**

---

## 3. Area A — the verdict parser and its new pinned corpus

I extracted `validation_gate.go:143-261` **verbatim** into a standalone module and drove it directly. Nothing below rests on the payload's own tests.

### 3.1 My corpus run

Over every `.md` in `docs/reviews/` (78 files: 77 committed plus this prompt):

| Verdict read | Count |
| --- | --- |
| `ACCEPT` | 23 |
| `REVISE` | 15 |
| `REJECT` | 1 |
| none | 39 (38 committed + this prompt) |

Against the expectations file: **77 records pinned, 77 compared, zero mismatches, in both directions.** The pinned set is also *complete* — `git ls-tree -r HEAD docs/reviews` yields exactly 77 `.md` files and the pinned key set is identical, nothing tracked-but-unpinned, nothing pinned-but-untracked.

### 3.2 The expectations file's independence — round-specific question (2)

**The independence claim is overstated as written, but the labels are correct, and I verified that by reading rather than by scanning.**

The `derivation` field says the labels were "derived by a scan independent of `validation_gate.go`", then describes a scan that restates `validation_gate.go`'s algorithm clause for clause, including this round's *novel and contested* comma-or-slash rule. A second implementation of the same specification is not independent ground truth; a shared specification error would produce agreeing labels and silently define correctness wrongly. That is precisely the hazard the prompt names.

So I checked the labels the only way that settles it — against what the documents say:

- I extracted every line containing `verdict` for all 77 records and read the deciding line of each in context. **I saw the literal deciding text for every one of the 77.**
- I read in full the headers of the seven records taking the heading-continuation path, where the label depends on a second line: arm-b-human-factors (**ACCEPT** ✓), absence-acceptance-payload-review (**ACCEPT** ✓), protocol-v5-amendment-payload-review (**ACCEPT** ✓), protocol-v5-proposal-review (**REVISE** ✓), -review-2 (**REVISE** ✓), -review-3 (**REVISE** ✓), -review-4 (**ACCEPT** ✓).
- I opened the five `none`-labelled documents that are not evaluator prompts. None states a verdict of its own. The validation-schedule candidate is the interesting one: it *cites* another record's ACCEPT, and `none` is correct. ✓
- I checked every `-prompt` document's boilerplate individually. All 39 correctly `none`.

**Conclusion: I found no wrong label.** The ground truth is sound. But the `derivation` sentence should not claim independence from `validation_gate.go` while reciting that file's algorithm. Recorded as **N4 (P3)**.

### 3.3 The enumeration rule — round-specific question (1)

**Declining review 4's prescription was correct.** I verified the counterexample directly: review 3's decisive line is `**Verdict: REVISE** — … it refuses 17 committed ACCEPT records`, and review 4's literal rule refuses to read it. Review 4 stated it validated its prescription against accepting records only, and that omission is exactly what its rule would have broken. Adopting it would have re-created the round-3 false-refusal class one round after closing it. Declining was right, and the reasoning in the candidate note is accurate.

**The replacement rule is sound in the false-refusal direction and incomplete in the false-accept direction.** Measured:

Correctly read as statements: `ACCEPT, with the residuals below` → `ACCEPT`; review 3's line → `REVISE`; all 23 committed accepting records → `ACCEPT`.

**False accepts**, all measured:

| Input | Read | Should be |
| --- | --- | --- |
| `Verdict: ACCEPT or REVISE` | **`ACCEPT`** | none |
| `Verdict: ACCEPT / REVISE` (spaced slash) | **`ACCEPT`** | none |
| `Verdict: ACCEPT; REVISE; REJECT` | **`ACCEPT`** | none |
| `Verdict: ACCEPT \| REVISE` | **`ACCEPT`** | none |
| `Verdict: ACCEPT — REVISE — REJECT` | **`ACCEPT`** | none |
| `Verdict: ACCEPT,` / `REVISE, or REJECT.` (across two lines) | **`ACCEPT`** | none |

`Verdict: ACCEPT/REVISE` unspaced is refused only incidentally — `closureVerdictWord` trims non-letters from the ends only, so the token never matches the map — not because the enumeration rule fires. The rule fires on exactly two characters, `,` and `/`, and only attached to the token.

In the other direction: `Verdict: ACCEPT, superseding the earlier REVISE` — a natural accepting line — is **refused**. Fail-closed and minor, but the rule is not clean in either direction.

**A sound fix exists and I verified it before prescribing it.** I patched the extraction to distinguish enumeration from statement by *adjacency* rather than by separator character: walk the fields after the token; a verdict word reached having passed only punctuation and coordinating conjunctions (`or`/`and`/`nor`) is an enumeration; the first field of real prose ends the scan and makes it a statement; a token trailing a separator with only connectives after it is a continued list. Measured:

- **All 77 committed records: identical to the pinned labels. Zero regressions, both directions.**
- Closes `or`, spaced `/`, `;`, `|`, dash, and comma-at-end-of-line.
- Keeps `ACCEPT, with the residuals below`, and now also correctly reads `ACCEPT, superseding the earlier REVISE`.
- Keeps review 3's line → `REVISE`.

Roughly fifteen lines. I mention it so the finding is actionable, not to insist on this shape.

### 3.4 Block-context tracking

Fence tracking is new this round and closes most of review 4's N3. It leaks in these constructions, all measured:

| Construction | Read | Should be |
| --- | --- | --- |
| `~~~` fence containing a backtick fence containing `Verdict: ACCEPT` | **`ACCEPT`** | none |
| backtick fence containing a `~~~` fence | **`ACCEPT`** | none |
| a 4-space-indented fence marker inside an open fence (closes it early) | **`ACCEPT`** | none |
| `~~~` closing a backtick fence (marker types not distinguished) | **`REJECT`** from inside the block | none |
| multi-line HTML comment containing a verdict | **`ACCEPT`** | none |
| tab-indented code block | **`ACCEPT`** | none |
| `## Verdict` followed by a 4-space-indented `ACCEPT` | **`ACCEPT`** | none |
| `## Verdict` followed by `>ACCEPT` (no space) | **`ACCEPT`** | none |

Three root causes: `fenced` is a single boolean with no record of marker type or length; `strings.HasPrefix(line, "    ")` is the only indented-code test, so tabs pass; and `nextNonBlankLine` on the heading branch applies **none** of the skips the main loop applies. The payload's own test pinning `## Verdict` + `> ACCEPT` as refused passes only because the space after `>` makes it a separate field that reduces to the empty string. Remove the space and the case fails.

Correctly handled: a fence opened inside a block quote does not corrupt outer state; an unclosed fence swallows the rest of the document (fail-closed); single-line HTML comments, table rows, and `## Verdict` + a list item, a fenced token or a table row are all refused.

### 3.5 Separator requirement and label furniture

No defect found. Pipe-separated cells, `Verdict ACCEPT was withheld…`, `Verdict-setting is hard`, `Verdicts of prior rounds`, `## Verdict bar`, `Review record and verdict`, and `**Verdict: ACCEPTED**` are all refused. Review 4's N4 false refusals are fixed: em dash, en dash, footnote `[^1]`, NBSP and CRLF all read correctly, each confirmed. Removing `>` from `closureLabelFurniture` and adding `\r` are both correct; without `\r` in `closureTokenFurniture` a CRLF `## Verdict\r` would yield `\r` as the separator candidate and refuse the heading branch, so the addition is load-bearing and pinned.

### 3.6 Is the parser correct, or correct against everything anyone has thought to try?

**The latter, and demonstrably so.** It is exactly correct against the 77-record corpus and the 31 constructed cases the payload pins. I wrote 43 cases in about an hour and eleven are wrong. That is not a comment on this round's care — round 5's parser is markedly better than round 4's, which was markedly better than round 3's — it is a comment on the shape of the problem. See § 9.

---

## 4. Area B — interaction defects in the revision-5 changes

**No cross-component interaction defect found.**

- **`>` removed from the furniture cutset** — safe; block-quote lines are skipped before the label test, and neither guard is redundant given the other's failure modes.
- **`\r` added** — necessary, not merely harmless (§ 3.5).
- **Pinned expectations replacing the live sweep** — closes review 4's N5 cleanly. `loadClosureVerdictExpectations` guards `V == 1`, `len(Records) >= 60`, and two required `purpose` substrings, so the pin cannot be silently emptied or reworded to weaken the claim.
- **The corpus test's two directions** — review 4's N1 genuinely fixed. `expected` now comes from committed bytes rather than the parser, so the false-refusal direction has real force. The `gateErr` block is now redundant with `assertClosureVerdict` rather than unreachable; harmless.
- **The candidate guard test** rebased to base commit `fec6d554…` and file count 18, consistent with the inventory, and correctly not in the frozen 18.

---

## 5. Area C — the whole payload, read fresh

| § 9 condition | Where enforced | Invalid resume passes? | Valid one refused? |
| --- | --- | --- | --- |
| (a) no outcome inspection | `execution_resilience.go:284` | No | No |
| (b) frozen bytes digest-identical | `:284`, `:293-306` (compared to live values) | Not for schedule/world-build; whole-set re-verification is the custodian's — recorded residual | No |
| (c) resume within 72 hours | `:342-345`, anchored `:358-374` to the log and the runner clock | No — back-dating refused by the log anchor | No — skew allowance and staleness bound both pinned |
| (d) cause on the frozen list, classified before the decision | `:30-35`, `:280-291` | No — no `oom` entry; `process_kill` needs `KillPrincipal` | No |
| Promptness / lapse | `:346-349` | No | No — required only past a 1-hour lapse |
| At most one resume | `:195-210`, corroborated against `start` entries | No — a downward-edited counter is caught; missing log refuses | No |
| Post-closure gate | `validation_gate.go:93-103`, `:109-141` | **Weakened by N1/N2**; the `ClosureVerdict != "ACCEPT"` test itself is exact | No |
| Partial summary outcome-free | `:41-66`, guard at `execution_resilience_test.go:521-535` | No | No |
| Process-event log honesty | `process_events.go:31-40` | — | — |

**Post-closure gate encoded in frozen gate rules and applying to decision 0023.** Confirmed by diffing `pre-validation-artifacts.json` against decision 0025's state: three fields added, `validation_execution_status` is `indeterminate`, both closure fields empty, so `requirePostClosureAudit` refuses today. Both gate flags remain `false`. **The gate is armed, in frozen bytes.**

**Authoring exemption** cannot launder a validation execution: `validateScheduleForTranche` refuses a schedule whose declared tranche does not match the flag.

**Records match the implementation**, with the exceptions at N3 and N5. The boundary document is *honest* about the enumeration rule's actual shape — it states the comma-or-slash rule verbatim rather than claiming a general one. That is the right instinct; it also means the frozen prose codifies the incomplete rule.

**Freeze hygiene.** Exactly one deliberate red, the named one. World-build and grader digest inputs untouched. All 18 inventory files match. The payload diff is tight — seven files, `+319/-85`.

---

## 6. Disposition of every prior review's P1 and P2

**Review 1:** P1-1, P1-2, P2-1, P2-2, P2-3, P2-5, P2-6 — **all still fixed.** P2-4 still fixed in form and materially better this round; the residual of the same class survives as **N1/N2**, not a regression.

**Review 2:** P2-N1, P2-N2 and every P3 — **all still fixed, none regressed.** P3-N1 carried as a residual, still accurate.

**Review 3:** **P1-N1 still fixed**, and now against committed ground truth rather than a self-referential sweep — I re-measured 0 errors across all 77 records, both directions, where review 3 measured 17 false refusals and 12 false accepts. P3-N1(new-a) correctly carried; (new-b), (new-c), (new-d) still fixed.

**Review 4:**

| Finding | Disposition |
| --- | --- |
| **N1 (P2) corpus test's false-refusal direction unreachable** | **Genuinely fixed.** Ground truth is committed and parser-independent; both directions have force. One stale sentence survives in the unfrozen note (**N5**). |
| **N2 (P2) enumeration rejection fitted to prompt wording** | **Partially fixed.** All three of review 4's measured cases are closed — I re-ran each — and the rule is order-independent as claimed. But it is still *character*-fitted rather than *structure*-fitted, and five further spellings leak. Not regressed; the same class at one further remove. Recorded as **N1**. |
| N3 (P3) no block context | **Substantially fixed.** Five constructions still leak; recorded as **N2**. |
| N4 (P3) fitted cutset refuses genuine records | **Fixed**, by review 4's own second route: the cutsets are widened and the accepted forms are stated in the frozen boundary document. |
| N5 (P3) frozen test coupled to a live directory | **Genuinely fixed.** The pin is complete against `HEAD`. |
| N6 (P3) count inaccuracy | **Fixed** — 23/15/1/38 matches my measurement exactly. |

**No prior P1 or P2 has regressed, and none has been papered over.**

---

## 7. New findings

**0 P0 · 0 P1 · 1 P2 · 4 P3.**

### N1 — P2. The enumeration rule fires on two characters; five other spellings of a verdict list read as an accepting record

**File:** `validation_gate.go:237-248`, specifically the `HasSuffix(trailing, ",") || HasSuffix(trailing, "/")` guard at `:238-239`.

**Measured false accepts** (§ 3.3): `ACCEPT or REVISE`, `ACCEPT / REVISE`, `ACCEPT; REVISE; REJECT`, `ACCEPT | REVISE`, `ACCEPT — REVISE — REJECT`, and a list continued onto a second line. Each returns `ACCEPT` and each opens the post-closure gate.

**Direction: permits an invalid resume.** A `validation_execution_closure_review` naming an evaluator *prompt* rather than the closure record — the exact one-token slip review 3 named as the class this check exists to catch — satisfies `requirePostClosureAudit` with no review performed, whenever that prompt lists the verdicts in any of those forms.

**Why this is not hypothetical.** All five payload-review prompts in this project, including the one I am reviewing under, carry:

```
## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, …
```

The parser refuses these for exactly one reason: the sentence begins with `Return`, which lands in field 0 and is not a verdict token. I verified the margin directly — delete that single word and the shipped parser returns **`ACCEPT`**. The decision-0023 closure-review prompt has not been written yet, and it will be drafted from this template. Whether the gate still guards its own primary failure mode at the moment it matters is currently decided by one word of boilerplate that nobody has been asked to preserve.

**Mitigations, recorded honestly.** `requirePostClosureAudit` is unreachable while `may_open_validation` is `false`. Both gate fields are frozen bytes, so exploitation needs a chair error *plus* a refreeze that is itself independently reviewed. Those are why this is P2 and not P1.

**Fix, verified before prescribing (§ 3.3):** the adjacency test. Measured over all 77 records, **zero regressions in either direction**; it also repairs a false refusal the current rule produces. **But see § 9 before spending a sixth round on this.**

### N2 — P3. Five block constructions still read a quoted verdict as the record's own

**File:** `validation_gate.go:160-188`. Root causes at `:163` (bare boolean, blind to marker type and length), `:169` (misses tab-indented code), `:182` (`nextNonBlankLine` applies none of the skips).

**Direction: permits an invalid resume.** A record whose first verdict-shaped line is another record's verdict, quoted in one of these ways, reads as its own. No committed record trips any of them; I checked all 77.

**Fix:** track the opening marker's rune and length and require a matching closing fence; treat a leading tab as indented code; route the heading-continuation line through the same skips as the main loop. The last is two lines and removes an inconsistency that the payload's own test passes only by accident of a space character.

### N3 — P3. The frozen boundary document claims block-context coverage the code does not have

**File:** `validation-execution-boundary.md:78` — "a verdict inside a fenced block, a block quote, or an indented block, which belongs to whatever record is being quoted". Stated without qualification, in bytes about to be frozen; falsified by every case in N2. Same class as review 2's P3-N10.

**Direction: neither.** A claimed control stated more strongly than it exists, which is how a future reader stops checking.

**Fix:** land N2, or qualify the sentence.

### N4 — P3. The expectations file's `derivation` claims independence while reciting the implementation's algorithm

**File:** `closure-verdict-corpus-expectations.json`, `derivation`. A second implementation of one specification is not independent ground truth. **The labels are in fact correct** — I verified them by reading (§ 3.2) — so this is a wording defect, not a correctness defect. But the file is the payload's answer to "how do we know the parser is right", it is frozen, and the sentence invites a future reader to trust an independence that was not obtained.

**Direction: neither.**

**Fix:** say what was done — "labels derived by a second scan written from the specification, and confirmed by reading each record's own verdict statement" — and, for real independence in frozen bytes, record who read the records.

### N5 — P3. A stale paragraph in the candidate note names a test that no longer exists

**File:** `execution-resilience-candidate.md:126`. Names `TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus` and describes it sweeping every committed record; the test is now `TestClosureVerdictParserMatchesTheCommittedCorpusExpectations` and deliberately does not sweep, which line 146 of the same document then says. The counts in the paragraph are correct.

**Direction: neither.** Not among the 18 frozen files, so lighter than N3.

**Fix:** past-tense the paragraph or point it at the current test name.

---

## 8. Residual list — **not operative at this verdict**

For the refreeze that follows the N1 fix, these are the entries I would accept. I re-derived each from the code at this commit. Three of review 4's entries are **retired** by this payload and must not be carried forward: the live-directory sweep entry, and the two block-context/cutset entries, which are **replaced** by the reworded ones below.

- "Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian's written attestation, as section 9 requires."
- "Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b)."
- "Residual: the process-event log is unsigned custodian-controlled local state, so a coordinated rewrite that deletes a start entry and re-attests the new digest defeats both the single-resume corroboration and the 72-hour anchor; the runner retains no previously accepted digest and cannot detect it. Detection is the committed attestation digest compared across commits, under independent review."
- "Residual: the 72-hour anchor is measured against the runner's own clock, which the custodian's host controls; a clock set to shortly after the last process event passes every timing check, while a clock set before it is refused by the interruption anchor. No trusted time source exists on the execution host."
- "Residual: measuring the 72 hours from the log's last entry rather than from the attested interruption shortens the real window by the time from the last completed pairing-key group to the interruption, bounded by one group; the erosion is disclosed in the boundary document and is conservative in the correct direction."
- "Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports 'gate is closed' first; both are refusals and only the surfaced message differs."
- "Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag. One consequence: the end-to-end resume-counter round-trip test runs on the authoring path and so does not traverse requireResumeAuthorization."
- "Residual: nothing in the runner writes validation_execution_status. A custodian who wipes the output directory presents as a fresh run and bypasses the resume rule entirely, so the post-closure gate — armed by the chair setting that field at closure, under the freeze workflow — is the sole control on the retry channel after a future non-completed execution. It is armed by frozen bytes for the current one."
- "Residual (P3): the closure verdict parser skips fenced blocks, block quotes and four-space-indented blocks, but its fence state is a single boolean blind to marker type and length, its indented-code test does not recognize a leading tab, and the heading-continuation branch applies none of the three skips; nested or mixed-marker fences, a tab-indented block, a multi-line HTML comment, and a quoted or indented token on the line after a Verdict heading therefore read as the record's own verdict. No committed record trips any of them, and the gate document's two fields are frozen and chair-set."
- "Residual (P3): the parser's qualifier set is closed to `final` and `overall`; a future closure record written outside the documented forms is refused. The failure is closed, and the accepted verdict-line forms are stated in the boundary document."
- "Residual (P3): the pinned corpus is complete against the 77 committed review records at this commit and is not extended automatically; a closure record added later is not covered until the pin is deliberately extended, which is the intended trade."
- "Residual (P3): the partial summary's committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure."
- "Residual (P3): the partial summary carries five non-outcome identity fields beyond section 9's literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome."
- "Residual (P3): the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log. Because the termination handler deliberately calls os.Exit, a supervised stop as well as an abrupt kill can truncate a boundary write; the failure is closed — decodeStrict refuses and the tranche closes indeterminate — but a routine Ctrl+C at an unlucky microsecond can burn the tranche. A temp-and-rename in writeJSON would remove the class."
- "Residual (P3): a validation output directory whose first custodian process wrote its start event but died before the initial checkpoint is refused wholesale by the progress-implies-checkpoint rule and cannot be resumed; the window is milliseconds and the rule is what prevents a phantom start entry from permanently refusing a later genuine resume."
- "Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first."
- "Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required."
- "Carried P2-1 and P2-5 remain retired by this payload; runScheduledCases and validateExternalGrade remain at gocyclo 14, both pre-existing and unchanged by this payload."

### My judgment on each carried residual

All are **correctly characterized and genuinely acceptable**; none should be promoted. Two notes:

- **The `writeJSON` atomicity residual** is still **absent from the inventory's `accepted_residuals` array**, exactly as review 4 flagged. Transcribe from this § 8, not from the array.
- **Review 4's conditional on the post-closure-gate residual applies, and I am invoking it.** Review 4 wrote that if the chair carries its N2 rather than fixing it, the residual must gain a sentence stating that the runner's corroboration of the named closure record remains heuristic and can accept a document that merely lists ACCEPT. My N1 is that carried remainder. If the chair accepts N1 rather than fixing it, **that sentence must be added** — and I would not accept that trade either, which is why N1 is a finding.

---

## 9. The mechanism-shape question — my answer, which is the most useful thing in this record

**The chair's alternative is right, and I would take it over another parser round.**

Three consecutive independent reviews have now found defects in this one component, and each round's defects were introduced or left by the previous round's fix. That pattern is usually read as insufficient care. I do not think that is what is happening. Round 5's parser is a genuine improvement on round 4's, which was a genuine improvement on round 3's, and this round's implementer caught a real error in the prescription they were given. The work is good. The problem is that the input space is *English prose written by a future author under schedule pressure*, and a line-furniture heuristic over that space has no closure condition. Every round closes the cases someone thought to try and leaves the ones nobody did. I wrote 43 cases in an hour and broke it eleven times; the next reviewer will write 60 and break it again.

Look at what the rest of `validation_gate.go` does. Schedule identity: digest comparison. Manifest identity: digest comparison. Label registry: digest comparison. World-build: digest comparison. Grader: digest comparison. Process-event log: digest comparison. Every identity check in this payload is a digest comparison except one, and that one is the only thing that has produced a finding in three rounds.

**The proposal: pin the closure record by SHA-256 in the gate document; take the verdict from the frozen field; reduce the runner's check to `digest(named file) == pinned digest`.** Delete `closureReviewVerdict`, `closureVerdictLabel`, `closureVerdictToken`, `closureVerdictEnumerates`, `closureVerdictWord`, `nextNonBlankLine`, all three cutset constants, the corpus test and the expectations file. A net deletion of roughly 200 lines of frozen code and data, replaced by about ten.

What is lost is the corroboration that the named document *says* ACCEPT. That loss is smaller than it looks. The parser never established that a human independently reviewed anything — it established that a string appears in a shape. The judgment "this document is an accepting independent review of the closure" is a human judgment, and under this project's own freeze workflow it is already made by a human: the independent reviewer of the refreeze that sets those gate fields. A reviewer reading a closure record cannot be fooled by `ACCEPT / REVISE` under a `## Verdict` heading, by a verdict quoted inside a nested fence, or by a prompt named in place of a record. The frozen parser can be fooled by all three, today. The digest pin moves the check from a machine guessing at prose to a machine verifying an identity, and leaves the prose judgment where it already lives and where it is competent.

The digest pin is also strictly stronger against the failure the parser cannot touch at all: a named record edited *after* the review that blessed it. The parser re-reads whatever bytes are on disk and would happily re-bless a rewritten record. A digest pin refuses it.

If the chair prefers to keep the parser — and there is a reasonable argument that a defence-in-depth check, even a leaky one, is worth ten lines — then the adjacency fix in § 3.3 closes N1 with zero regressions across all 77 records, and the three changes in N2 close the block-context leaks. I verified the first myself. But I would rather the chair spent this round deleting the component than hardening it, and my verdict is REVISE either way.

---

## 10. What I did not verify

- I did not run a model, obtain a private grade, or observe any validation outcome. Both gates were `false` throughout.
- I did not construct a live end-to-end validation resume; that requires the validation gate open.
- I did not exercise the signal path against a live multi-hour scheduled run.
- I did not recompute the world-build or grader digests from source; I confirmed no file in either input set is touched by the payload diff, and both pinning tests are green.
- I did not verify that the process-event log's digest is committed with the custody record in practice; no v5 custody record exists yet.
- I did not inject any test into the repository. I drove the parser by extracting `validation_gate.go:143-261` verbatim into a scratchpad module, and drove the proposed fix by patching a copy. Both were run against the live `docs/reviews/` directory and a 43-case constructed suite of my own design. **Every count and parser result in this record is measured, not estimated.**
- The 77 labels I verified by reading the deciding line of all 77 records and reading twelve in full. I did not read all 77 end to end.
- I made no change of any kind to the repository.

---

## Chair transcription note

Transcribed verbatim from the fifth independent reviewer's record. Before transcription the payload author independently reproduced N1 by driving the shipped `closureReviewVerdict` against every listed spelling: `ACCEPT or REVISE`, `ACCEPT / REVISE`, `ACCEPT; REVISE; REJECT`, `ACCEPT | REVISE`, `ACCEPT — REVISE — REJECT` and the two-line list all return `ACCEPT`; `ACCEPT, superseding the earlier REVISE` is falsely refused; and the shipped prompt boilerplate is refused **only** because the sentence opens with the word `Return` — deleting that one word yields `ACCEPT`. The verdict is **REVISE**. § 8 is conditional and is not operative. The § 9 mechanism-shape recommendation is put to the chair as a decision, not actioned unilaterally.
