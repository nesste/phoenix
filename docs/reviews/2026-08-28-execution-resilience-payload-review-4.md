# Fourth independent review — protocol-v5 § 9 execution-resilience payload

- **Reviewer role:** evaluation reviewer, independent of this payload's implementation. I did not author the payload and I am not any of the three reviewers who produced the prior records. Fourth review.
- **Date:** 2026-08-28
- **Payload commit:** `70ff7875340071324a0006a5504bbf527c21a5f5` (`experiments: fix the section 9 payload findings from the third review`)
- **Verdict:** `REVISE`

No model was run, no private grade obtained, no validation or held-out outcome observed. `may_open_validation` and `may_open_held_out` were `false` at the start and end of this review and I did not touch them. I made no change of any kind to the repository: my only writes were to the session scratchpad (a standalone driver of the verdict parser, an attack corpus, and one Go test injected via `go test -overlay`, which writes nothing to the tree) and to this record.

---

## 1. Pinned inputs — verified

| Input | Expected | Computed | |
| --- | --- | --- | --- |
| `git rev-parse HEAD` | `70ff7875340071324a0006a5504bbf527c21a5f5` | `70ff7875340071324a0006a5504bbf527c21a5f5` | match |
| Ancestry | `6f54e8ca…` → `60c6aee5…` → `9f53fa73…` → `e92eaacf…` → `46030d4f…` → `7a1ea425…` → `2cf99331…` | exact, in order | match |
| Working tree | exactly one untracked file, this prompt | `?? docs/reviews/2026-08-28-execution-resilience-payload-review-4-prompt.md`, nothing else | match |

`git status --porcelain --untracked-files=all` returned that single line and nothing else, at the start of the review and again after every scratch operation. **No P0 on tree cleanliness.**

### Digests I computed

Rule applied: CRLF→LF, then CR→LF, then SHA-256, rendered `sha256:<hex>`.

| Artifact | Digest | Pinned |
| --- | --- | --- |
| `artifacts/gate-1a-execution-resilience-candidate.json` (raw) | `sha256:1eecf1e6d847c4bb6222c7d4c7fee37aebdce124abdcb763c19bed8dcb923013` | match |
| `artifacts/execution-resilience-candidate.md` (raw) | `sha256:d9e4d21372645449b8820973e0cd5b77f8b80e9bf778704a38990e4ff5dce89b` | match |
| Ordinal-path set digest, 17 files | `sha256:2fa91f23537f0bf7ea87ba73f60256a6291cf40ba896ee5694e53f6f71a9e245` | match |

Both files are LF-only in the tree (raw and LF-normalized digests are identical). All 17 inventoried files re-hashed individually and every one matches its recorded value:

| Path | LF-normalized SHA-256 |
| --- | --- |
| `experiments/frontier-v1/README.md` | `sha256:fc748a2ffab7e4cf4256edb232bd5682802835391f2704e845f11082f07ed45b` |
| `experiments/frontier-v1/corpusctl/internal/corpus/protocol_test.go` | `sha256:4495d032c7d86fbfa03952ad676f3e769de4ca1f0d071dc694afa376742cc952` |
| `experiments/frontier-v1/pre-validation-artifacts.json` | `sha256:ce14d55cdbb36c9102e71830101eb45f0bd29df91af35a51134c57b72b3f4cc2` |
| `experiments/frontier-v1/protocol.json` | `sha256:726b68112108542ba3491fde538342399cbb446a04755a7e0208e076cbdb0d32` |
| `experiments/frontier-v1/runner/artifact_freeze_test.go` | `sha256:3f62b2e97e52f08296b1383c7f159b73da8e38b646219874662e9805b234ff95` |
| `experiments/frontier-v1/runner/closure_verdict_corpus_test.go` | `sha256:444955e0d4724827e96dbeb0db3dbbd8a76cbf55da21a7fd9cdc7411ca688161` |
| `experiments/frontier-v1/runner/execution_resilience.go` | `sha256:fe98996728d7ccbd4b1bbd23f44d29e3ebb543a64e2430502a7ccbfb2dee374a` |
| `experiments/frontier-v1/runner/execution_resilience_test.go` | `sha256:9cf56873fb0bdcc0ad229c101b3d919b631ddb1ef0520241a72311833de34957` |
| `experiments/frontier-v1/runner/main.go` | `sha256:192d275ef16be74f9a3f47747cf3d14d8df52d18d16127734db786866c5bdff3` |
| `experiments/frontier-v1/runner/process_events.go` | `sha256:605a365310aebb4f98c5e75fca8ae96232e671a48381756557f7c891ec70c617` |
| `experiments/frontier-v1/runner/scheduled_resume.go` | `sha256:aa050eb443a6676a1ab540788c182a3ef1ec4e63ae2c7bfd4fdb4b23fac9fc2e` |
| `experiments/frontier-v1/runner/scheduled_resume_test.go` | `sha256:a090f71facb3d21ba3e7d1df32988cb7a64d277b5bdba97a46637bf435c2b893` |
| `experiments/frontier-v1/runner/scheduled_run.go` | `sha256:0956afad8e186b1bef71a64c44cab9c16aa64da6c8ee82cb9188fa49fab4cb9a` |
| `experiments/frontier-v1/runner/types.go` | `sha256:fce5d7b956df609fe56bf77e9e400eed0f79f23374778915050948035ba92191` |
| `experiments/frontier-v1/runner/validation_execution_test.go` | `sha256:9ef5e4c7e927082d44830bba29016b446b43bfa3c7a37974e18a850dc012d20e` |
| `experiments/frontier-v1/runner/validation_gate.go` | `sha256:781d84cb770d1edfd1df439847369127ce64f15513626d43b0bb5c635c5b3984` |
| `experiments/frontier-v1/validation-execution-boundary.md` | `sha256:a9a616d3eb2c90833d841353f695c6f690f270d4aa74375a1336e76af970d767` |

All six LF-normalized digests named in the prompt match exactly. The inventory's 17 files match the working tree with no extra and no missing path.

**World-build and grader digests unmoved.** `git diff 2cf99331…HEAD -- experiments/frontier-v1/pre-validation-artifacts.json` is 4 insertions / 1 deletion, confined to the `gates` block (the three new post-closure fields). `sha256:425bab1c…ffdae4` (world build) and `sha256:8146a68a…baebf0b` (grader) are present and untouched.

---

## 2. Test and quality command results

| Command | Result |
| --- | --- |
| `go test ./...` | **1 failure, the expected one.** `--- FAIL: TestPreValidationFreezeMatchesAcceptedCandidates` — `README.md digest = sha256:fc748a2f…, freeze requires sha256:b486fc55…`. Every other package `ok`. |
| corpusctl suite (`cd experiments/frontier-v1/corpusctl && go test ./...`) | **pass**, exit 0 |
| `go vet ./...` | **clean**, no output, exit 0 |
| `staticcheck` v0.7.0 `./...` | **clean**, no output, exit 0 |
| `gocyclo -over 15` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis` | **clean**, no function over 15 |
| `dupl -t 100` (via `quality-check no-output`) over the same paths | **clean**, exit 0 |
| `make validate-spec` | **pass** — 9 examples valid against their schemas |
| `make validate-authoring` | **pass** — `valid …/manifests/authoring.json` |
| **`go test -race ./experiments/frontier-v1/runner/`** | **I ran this myself**, on the WSL2 Ubuntu host. **No data race reported.** The only failure is the same expected freeze red; run time 42.9 s. |

**Freeze hygiene: exactly one deliberate red, and it is the claimed one.** No other test fails on either host. I did not attempt to defeat the freeze test.

---

## 3. Area A — the rewritten verdict parser

I did not rely on `TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus` passing. I extracted `validation_gate.go:141-216` verbatim into a standalone module and drove it myself.

### A1. The parser over every `.md` under `docs/reviews/` — measured

**Over the 75 committed records at this commit: 23 accepted, 52 refused.**
(Over the 76 files in the working tree, i.e. including this review's untracked prompt: 23 accepted, 53 refused.)

I checked ground truth independently of the parser, by reading the verdict statement of every file:

- **All 23 accepted files are genuinely accepting review records.** 20 use the bolded-label-plus-code-span style, 3 use the current `## 1. Verdict` heading style with `**ACCEPT.**` on the next line. **Zero false accepts.**
- **All 52 refusals are correct.** 13 are non-prompt records genuinely stating `REVISE`; 1 states `REJECT` (`2026-08-21-protocol-v4-validation-execution-boundary-review.md:5`); the rest are evaluator prompts, handoffs, status notes, candidate notes and reports that state no verdict of their own. **Zero false refusals.**

**Review 3's P1-N1 is genuinely fixed against the live corpus, in both directions.** Its measured 17-of-19 false refusals and 12 false accepts are all gone. I reproduced this from scratch rather than accepting the claim.

One observation, not a finding: `REJECT` is not in the parser's token set, so a `REJECT` record reads as *no verdict* rather than as a non-ACCEPT verdict. Both refuse, so gate behaviour is correct. It matters only in that "first-statement-decisive" does not protect a `REJECT` record — see N3.

### A2. Constructed attacks — 38 cases, both directions

I wrote and ran a 38-case suite of my own. **21 diverge from what a reader of the record would conclude.**

**False accepts (a document that is not an accepting record reads as one):**

| Constructed case | Parser result |
| --- | --- |
| `Verdict: ACCEPT` inside a fenced code block | **ACCEPT** |
| bolded-label-plus-code-span inside a ` ```markdown ` block | **ACCEPT** |
| `> **Verdict:** ACCEPT` in a block quote | **ACCEPT** |
| `  > Verdict: ACCEPT` in an indented block quote | **ACCEPT** |
| 4-space-indented code block `    Verdict: ACCEPT` | **ACCEPT** |
| `1. Verdict: ACCEPT, REJECT, or REVISE.` | **ACCEPT** |
| `Verdict: ACCEPT, with findings, or REVISE.` | **ACCEPT** |
| `## Verdict` + `ACCEPT, REJECT, or REVISE.` | **ACCEPT** |
| `Verdict ACCEPT was withheld for the reasons below.` | **ACCEPT** |

Correctly refused: table rows (`|` is not furniture), `ACCEPTED`, `ACCEPTABLE`, a mid-sentence `…state a verdict: ACCEPT…`, a heading that is a filename, a `REJECT` record citing a later `**Verdict: ACCEPT**`, a `## Verdict` heading whose token line is itself block-quoted, and a `REVISE` record citing an earlier ACCEPT mid-sentence.

**False refusals (a genuine ACCEPT record is refused):**

| Constructed case | Parser result |
| --- | --- |
| `**Verdict — ACCEPT**` (em dash) | refused |
| `**Verdict – ACCEPT**` (en dash) | refused |
| `## Verdict` + `— ACCEPT` | refused |
| `- [^1] **Verdict:** ACCEPT` (footnote marker) | refused |
| `[1] **Verdict:** ACCEPT` | refused |
| NBSP before the label | refused |
| `**Review verdict: ACCEPT**` | refused |
| `**My verdict: ACCEPT**` | refused |
| `**Independent verdict: ACCEPT**` | refused |
| `**Fourth-review verdict: ACCEPT**` | refused |
| `**Verdict (fourth review): ACCEPT**` | refused |
| `## Verdict and findings` + `ACCEPT` | refused |

Correctly accepted: `### 2. Verdict:` + `ACCEPT`; `**Verdict - ACCEPT**` (ASCII dash); NBSP *after* the colon; **CRLF in both the inline and the heading style** (I checked this specifically — `\r` survives `TrimLeft` but `strings.Fields` and `TrimSpace` both treat it as whitespace); tab separation; `**ACCEPT.**` with trailing period.

### A3. Are the closed qualifier set and the three-field window principled, or fitted?

**Fitted, and I can show it rather than assert it.**

The qualifier set `{final, overall}` is exactly the two qualifiers that appear in this corpus and no more. Any other natural qualifier — `Review verdict:`, `Independent verdict:`, `My verdict:` — is refused, as is a parenthesised qualifier *after* the label. A future record using an unlisted qualifier is refused, and because `validation_gate.go` is frozen bytes, the repair is a full close → payload → review → refreeze cycle.

The three-field neighbourhood window is worse, because it is fitted in the *unsafe* direction. Every evaluator prompt in this repository writes the boilerplate as `Verdict: ACCEPT, REVISE, or REJECT`, which places `REVISE` at field 2 of 3 — inside the window by one position. Reorder those two words to `ACCEPT, REJECT, or REVISE`, which is equally natural English, and `REVISE` moves to field 4, the window misses it, and the prompt reads as an accepting record. That is finding **N2**. The guard that stops a `…-prompt.md` file from satisfying the post-closure gate is, today, word order in a document nobody has written yet.

### A4. Is coupling a frozen runner test to the live `docs/reviews/` directory sound?

**Partly sound, partly a defect (finding N5).**

In its favour: running the parser against the real corpus is the only thing that would have caught review 3's P1, and it is a genuine improvement over the two-file survey that produced that defect. I also have direct evidence the coupling is not immediately fragile — this review's own prompt file was added to `docs/reviews/` after the payload commit, and the test still passes with it present.

Against it: `closure_verdict_corpus_test.go` is one of the 17 frozen files, and it globs the *live* directory. Its outcome now depends on files that unrelated future work adds. A future evaluator prompt worded as in N2 turns this frozen test red, and going green again costs a full freeze cycle for a change with nothing to do with the runner. Separately, the test's only classifier of "this is a prompt, not a record" is the filename substring `-prompt` (`:37`, `:43`); three existing prompt-like documents — `2026-08-18-sealed-corpus-evaluator-handoff.md`, `2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md`, `2026-08-18-task-0.6-human-factors-pack.md` — do not match it and are checked as records. They refuse today, so there is no live problem, but the classifier is a substring, not a property.

---

## 4. Area B — interaction defects

**The boundary document's field list matches `resumeAttestation` exactly.** I transcribed the enumeration at `validation-execution-boundary.md:72` into a test injected with `go test -overlay` (no repo write) and compared it by reflection against every `json` tag on `resumeAttestation`. **All 20 names agree: no omission, no extra.**

**A custodian following the document literally can construct an accepted attestation — I constructed one and it works.** I built an attestation from the document's prose alone, wrote a fixture output directory with a `start` event, and called the real `requireResumeAuthorization`. **Accepted.** I then added the two optional fields the document names; still accepted. Review 3's **P3-N1(new-b) is genuinely fixed**, verified by execution rather than by reading.

**The other two revision-4 changes checked against assumptions elsewhere.** The corpus test is new and additive; it does not touch any anchor another check relies on. The parser rewrite is confined to `closureReviewVerdict` and its three helpers, all called only from `verifyClosureReviewVerdictLine`; I traced every caller. Neither the boundary document's new paragraphs nor the corpus test changes what any log contains, so the review-2 class of defect (a fix changing a log and silently breaking a different check's anchor) does not recur here. I specifically re-checked that the timing anchor still reads the maximum `recorded_at` rather than the last line, and that nothing in this revision adds a log entry.

---

## 5. Area C — the whole payload, as if fresh

| § 9 condition | Located in code | Can an invalid resume pass? | Can a valid one be refused? |
| --- | --- | --- | --- |
| (a) no outcome inspected, attested | `execution_resilience.go:284` | No — boolean required true | No |
| (b) frozen bytes digest-identical | `:284`, `:293-306` | No for schedule/world-build, compared against live values; byte-level re-verification is the custodian's, an accepted residual | No — a missing live digest refuses rather than skipping (`:298-300`) |
| (c) resume within 72 h | `:342-344`, `:358-383` | No — anchored to the log's last entry *and* the runner clock, so back-dating is refused | Not on an honest timeline; the 5-minute skew allowance (`:368`) covers a signalled interruption's stamp ordering |
| (d) cause classified before decision, on the frozen list | `:280-291` | No — `frozenResumeCauses` is a closed map, OOM absent, `process_kill` additionally requires `KillPrincipal` | No |
| Promptness duty / lapse explanation | `:346-349` | n/a | No — explanation demanded only past the 1-hour `promptResumeWindow` |
| At most one resume | `:195-210` | No — checkpoint counter corroborated against `start` entries; a lowered counter, a deleted or unreadable log all refuse | No |
| Post-closure gate | `validation_gate.go:92-102` | **Narrowed but not closed — see N2 and N3** | See N4 |

**Partial summary outcome-free in both directions.** Verified: the guard checks the written JSON against an admitted set in both directions and descends into `artifact_digests`; `partialArtifactDigests` has a fixed four-key shape with no `omitempty`.

**Process-event log honest about what the host does not expose.** `recordTermination` sets `exposed: false`, `source: signal_without_sender_identity`, and a reason. Correct: Go's `os/signal` delivers only the signal number on both supported hosts.

**Freeze obligation 5 discharged.** All four event kinds present, RFC3339Nano timestamps, principal where exposed. The post-closure gate is encoded in the frozen gate rules, and `verifyPostClosureGateState` pins four substrings of the rule text and both closure fields **empty** while `validation_execution_status` stays `"indeterminate"`. **It applies to the decision-0023 closure today.** The predicate is default-deny, covering budget stop and lapse without enumeration.

**Authoring-tranche exemption** cannot launder a validation execution: `validateScheduleForTranche` (`schedule.go:90-93`) refuses a schedule whose `Tranche` does not equal the flag. Verified independently.

**Ordering check:** `requireResumeAuthorization` runs inside `openScheduledProgress` **before** the `start` event is appended (`scheduled_run.go:81-103`), so a refused resume writes no `start` entry and does not consume the tranche's single allowance. Confirmed by reading the call order.

**Records match the implementation.** `README.md:177-179` is accurate on both new artifacts, the anchor, the corroboration, and the one-hour bound explicitly labelled an implementation control. `protocol.json`'s amendment and the 8-entry blocker list are pinned by the corpusctl suite, which passes. The overclaims I found are **N1** and **N6**.

---

## 6. Disposition of every prior P1 and P2

### Review 1

| Finding | Disposition |
| --- | --- |
| P1-1 signal handler swallows SIGINT/SIGTERM | **Still fixed.** Records, releases, then exits 128+signal. Not regressed. |
| P1-2 single-resume bound launderable | **Still fixed.** Corroborated against `start` entries; missing/unreadable log refuses. |
| P2-1 outcome-free test depth-1 | **Still fixed.** |
| P2-2 counter round trip untested | **Still fixed.** |
| P2-3 72-hour window self-reported | **Still fixed.** |
| P2-4 closure-record check existence only | **Still fixed in form**, materially improved this round — but see **N2/N3**, the surviving remainder of the same class, not a regression. |
| P2-5 gocyclo 14 | **Still fixed.** |
| P2-6 termination progress zero | **Still fixed.** |

### Review 2

All of P2-N1, P2-N2, P3-N2, P3-N4, P3-N5, P3-N6, P3-N7, P3-N8, P3-N9, P3-N10 **still fixed**; P3-N1 carried as a residual, still accurately characterized. P3-N3 superseded by review 3's P1-N1.

### Review 3

| Finding | Disposition |
| --- | --- |
| **P1-N1 verdict parser wrong in both directions** | **Fixed against the live corpus — measured, 0 errors in 75 records, both directions.** The false-refusal half is fully closed. The false-accept half is substantially narrowed but **not eliminated**: it survives as **N2** (word-order-dependent) and **N3** (no block context). Not a regression; the same class surviving at one further remove. |
| P3-N1(new-a) `writeJSON` not atomic | **Not fixed, correctly so** — carried as a residual. **Note:** absent from the inventory's 6-entry `accepted_residuals` array; present in review 3 § 8, the operative list. |
| P3-N1(new-b) attestation field set unnamed | **Fixed, verified by construction.** |
| P3-N1(new-c) unresumable start-without-checkpoint | **Fixed.** |
| P3-N1(new-d) stale round-count | **Fixed.** |

**No prior P1 or P2 has regressed, and none has been papered over.** Each one I re-derived from the code.

---

## 7. New findings

### N1 — P2. The corpus test's false-refusal direction is unreachable code, and three records claim it works

**File:** `experiments/frontier-v1/runner/closure_verdict_corpus_test.go:46-50` (the guard), `:10-17` (the doc comment).

`verifyClosureReviewVerdictLine(c)` returns `nil` **exactly when** `closureReviewVerdict(c)` returns `("ACCEPT", true)` — that is its entire body (`validation_gate.go:131-140`). Inside the test's `err != nil` branch, the negation therefore guarantees `!(found && verdict == "ACCEPT")`. Both calls are the same pure function on the same input. **The `t.Errorf` at line 49 can never execute.** I proved this by construction over the corpus and by the structural argument.

Consequently the test checks **one** direction only: that no `-prompt`-named file is accepted. It does not, and cannot, check that a genuine ACCEPT record is not refused — the direction review 3 measured at 17-of-19 wrong.

Three records assert otherwise, one in frozen bytes:
- `closure_verdict_corpus_test.go:15-17` (frozen): "The corpus is the arbiter for both directions."
- `gate-1a-execution-resilience-candidate.json`, `fixed_after_third_review` P1-N1: "…**or any record stating ACCEPT is refused**."
- `execution-resilience-candidate.md:126`: same claim.

**Direction:** neither, by itself — a test permits and refuses nothing. It removes the regression protection the payload offers against a recurrence of review 3's P1, and it is a claimed control that does not exist.

**Why P2 and not P3.** Review 2 graded a single untested branch (P3-N4) as P3, and by that precedent a test gap is P3. I grade this higher because it is not an untested branch: it is the payload's entire answer to the round-3 P1 in one of its two directions, it is asserted to exist in bytes about to be frozen, and it is roughly five lines to fix. I record my reasoning so the chair can regrade.

**Fix:** derive ground truth independently of the parser. A committed expectations table (`filename → expected verdict`) inside the frozen set, compared against the parser's output, gives both directions real force. Correct the three overclaiming sentences.

### N2 — P2. The verdict-list rejection is a three-field window fitted to today's prompt wording; two reordered words make an evaluator prompt read as an accepting record

**File:** `experiments/frontier-v1/runner/validation_gate.go:200-204`.

Measured, not hypothesised:

- `1. Verdict: ACCEPT, REJECT, or REVISE.` → parser returns **ACCEPT**
- `Verdict: ACCEPT, with findings, or REVISE.` → **ACCEPT**
- `## Verdict` + `ACCEPT, REJECT, or REVISE.` → **ACCEPT**

Every evaluator prompt in this repository happens to write `ACCEPT, REVISE, or REJECT`, which puts `REVISE` at field 2 — inside the window by a single position. The equally natural `ACCEPT, REJECT, or REVISE` puts it at field 4 and the guard misses entirely.

**Direction: permits an invalid execution.** `validation_execution_closure_review` naming a `…-closure-review-prompt.md` file — the one-token slip review 3 named as the exact class this check exists to catch — satisfies the post-closure gate with no review performed, whenever that prompt is worded in the ordering above.

**Mitigations, recorded honestly:** `requirePostClosureAudit` is unreachable while `may_open_validation` is false; the two gate-document fields are frozen and chair-set, so exploitation needs a chair error plus a refreeze that is itself reviewed. Those are why this is P2 and not P1 — unlike review 3's P1, it is not wrong against any file that exists today.

**Fix, and I tested the prescription before writing it:** reject the token whenever any other verdict word (`REVISE`, `REJECT`) appears **anywhere on the decisive line** — and on the next-non-blank line, when the heading branch takes it — rather than within a fixed three-field window. Order-independent. I ran this over all 75 committed records: **no accepting record's decisive verdict line contains `REVISE` or `REJECT`**, so the change causes zero new false refusals on the existing corpus. I verified the fix against the whole corpus rather than against the cases that prompted it, because prescribing from a narrow survey is precisely what produced review 3's P1.

### N3 — P3. The parser has no block context: fenced code, block quotes and indented code all read as verdict statements

**File:** `experiments/frontier-v1/runner/validation_gate.go:146` (the furniture cutset includes `>`), `:151-168` (no fence tracking).

Measured false accepts: `Verdict: ACCEPT` inside a fence; the same inside a ` ```markdown ` fence; `> **Verdict:** ACCEPT` in a block quote; `  > Verdict: ACCEPT`; a 4-space-indented `Verdict: ACCEPT`. Also `Verdict ACCEPT was withheld for the reasons below.` — a prose sentence that merely begins with the word, since the label test is purely positional.

**Direction: permits an invalid execution.** A document whose first verdict-shaped line is another record's verdict, quoted in the idiomatic way, reads as its own verdict. No committed file trips this today — I checked all 23 accepting records — but these records routinely cite prior rounds, and block-quoting is how one cites. Note also that `REJECT` is not in the token set, so a `REJECT` record that later quotes an ACCEPT is *not* protected by first-statement-decisiveness.

**Fix:** track fenced regions (` ``` ` / `~~~`, toggled on the trimmed line) and skip any line whose furniture prefix contains `>` or that is indented four or more spaces. Verified against the corpus: zero new false refusals.

### N4 — P3. The qualifier set and furniture cutset are fitted, and refuse plausible genuine ACCEPT records

**File:** `validation_gate.go:146` (cutset), `:178` (the closed set `{final, overall}`), `:163` (heading branch requires an empty remainder).

Measured refusals of genuine ACCEPT statements: em and en dash separators (ASCII `-` is in the cutset, U+2013/U+2014 are not — and this project's prose uses em dashes throughout); footnote markers `[^1]` and `[1]` (`[` not in the cutset); `Review verdict:`, `My verdict:`, `Independent verdict:`, `Fourth-review verdict:`; a parenthesised qualifier after the label; `## Verdict and findings` + `ACCEPT`; NBSP before the label.

**Direction: refuses a valid execution.** Fails closed, which is the safe direction, and that is why this is P3. The cost is real: the decision-0023 closure review must be written in one of the pinned forms, and if it is not, `requireValidationGate` refuses at the point of maximum schedule pressure, repairable only by another full freeze cycle.

**Fix (either is adequate):** add the dashes, brackets and NBSP to the cutset and allow any single leading word before `verdict`; **or**, since `validation-execution-boundary.md` is already in this payload, state the accepted verdict-line forms in it so the closure record's author writes a form the frozen parser reads. The second is cheaper and does not touch the parser.

### N5 — P3. A frozen test's outcome is coupled to a live directory that unrelated future work writes to

**File:** `closure_verdict_corpus_test.go:19-25`, `:37`, `:43`.

The test globs the live `docs/reviews/` directory, and the test file is one of the 17 frozen artifacts. A future review record or evaluator prompt — added by work with nothing to do with the runner — can turn this frozen test red; under N2's wording, one will. Going green again costs a full freeze cycle. Separately, `strings.Contains(name, "-prompt")` is the only classifier of "prompt, not record", and three existing prompt-like documents do not match it.

**Direction: refuses a valid execution**, indirectly — it makes a red freeze suite, not a refused resume. It also risks the worse outcome of the test being weakened under schedule pressure.

**Fix:** pin the corpus as a committed expectations file inside the frozen set (which also fixes N1), and keep the live-directory sweep as a separate non-frozen test; or restrict the sweep to files committed at or before the freeze commit.

### N6 — P3. Count inaccuracy in the candidate note

**File:** `execution-resilience-candidate.md:126` — "Current corpus result: 23 accepted, 51 refused."

Measured over the 75 committed records: **23 accepted, 52 refused.** The accepted count is right; the refused count is one low.

**Direction: neither.** The note is not among the 17 frozen files, so this is lighter than review 2's P3-N10, which was in frozen inventory bytes. I report it because it is the same class and the note is what the chair reads.

---

## 8. Verdict and what I would carry

**REVISE**, on **N1** and **N2**. No P0, no P1.

N1 and N2 are not in the residual list. Both concern frozen bytes that cannot be repaired at refreeze under this workflow; N2 permits the exact failure freeze obligation 5 exists to prevent; N1 is the payload's own claimed guard against a recurrence of the round-3 P1, asserted in three records and vacuous in fact. Each is a few lines. **N3, N4, N5 and N6 I would accept as residuals** if the chair prefers, but N3 and N4 are cheap and touch files the N2 fix reopens anyway.

I want to be plain that this is not a manufactured fifth round. Review 3's P1-N1 is genuinely fixed and I proved it independently: zero misclassifications across 75 records in both directions, where review 3 measured 17 false refusals and 12 false accepts. Every P1 and P2 from all three rounds is fixed and none has regressed. The remaining findings are narrower than any previous round's — but the pattern the prompt warned about has recurred a fourth time, one remove further out: the round-3 fix ships a false-accept path that survives on word order, and a test that claims to guard both directions and guards one.

### Residual list for the refreeze that follows the N1/N2 fixes

The verdict is REVISE, so **this list is not operative.** For the refreeze that follows, these are the entries I would accept into `accepted_findings`. I have carried review 3's § 8 wording verbatim where I re-derived and agree with it, since that list — not the inventory's 6-entry array — is the operative one under this workflow.

- "Residual: neither supported host exposes the sending principal of a signal to a Go handler, so a termination event records exposed=false with the reason; establishing that no project participant initiated a kill remains the custodian's written attestation, as section 9 requires."
- "Residual: the runner verifies that the attestation claims digest identity and that its recorded schedule and world-build digests equal the live ones; byte-level re-verification of the whole frozen set at resume is the custodian step recorded in the attestation, as section 9 words condition (b)."
- "Residual: the process-event log is unsigned custodian-controlled local state, so a coordinated rewrite that deletes a start entry and re-attests the new digest defeats both the single-resume corroboration and the 72-hour anchor; the runner retains no previously accepted digest and cannot detect it. Detection is the committed attestation digest compared across commits, under independent review."
- "Residual: the 72-hour anchor is measured against the runner's own clock, which the custodian's host controls; a clock set to shortly after the last process event passes every timing check, while a clock set before it is refused by the interruption anchor. No trusted time source exists on the execution host."
- "Residual: measuring the 72 hours from the log's last entry rather than from the attested interruption shortens the real window by the time from the last completed pairing-key group to the interruption, bounded by one group; the erosion is disclosed in the boundary document and is conservative in the correct direction."
- "Residual: requirePostClosureAudit runs inside requireValidationGate after the opening-flag checks, so a closed gate reports 'gate is closed' first; both are refusals and only the surfaced message differs."
- "Residual: the section 9 attestation rule is scoped to the validation tranche; authoring resumes are development loops and are exempt, pinned by TestAuthoringResumeIsExemptFromTheAttestationRule. The exemption cannot launder a validation execution because validateScheduleForTranche refuses a schedule whose tranche does not match the flag. One consequence: the end-to-end resume-counter round-trip test runs on the authoring path and so does not traverse requireResumeAuthorization."
- "Residual: nothing in the runner writes validation_execution_status. A custodian who wipes the output directory presents as a fresh run and bypasses the resume rule entirely, so the post-closure gate — armed by the chair setting that field at closure, under the freeze workflow — is the sole control on the retry channel after a future non-completed execution. It is armed by frozen bytes for the current one."
- "Residual (P3): the closure verdict parser reads markdown by line furniture and has no block context, so a verdict statement inside a fenced code block, a block quote, or an indented code block reads as the record's own verdict; no committed record trips this, and the gate document's two fields are frozen and chair-set. Block-context tracking would remove the class."
- "Residual (P3): the parser's qualifier set is closed to `final` and `overall` and its furniture cutset admits the ASCII hyphen but not the em or en dash, the footnote bracket, or a non-breaking space; a future closure record written outside those forms is refused. The failure is closed, and the accepted verdict-line forms are stated in the boundary document."
- "Residual (P3): TestClosureVerdictParserAgreesWithTheCommittedReviewCorpus globs the live docs/reviews/ directory from a frozen test file, so a review record or evaluator prompt added by unrelated future work can turn the frozen suite red; the classifier for 'this is a prompt' is the filename substring -prompt, which three existing prompt-like documents do not match."
- "Residual (P3): the partial summary's committed process-event-log digest never covers the terminal stop or termination entry, because recordStop appends after the last partial-summary write; the custody record hashes the whole log at closure."
- "Residual (P3): the partial summary carries five non-outcome identity fields beyond section 9's literal list (v, tranche, status, schedule_digest, resumes); none carries an outcome."
- "Residual (P3): the checkpoint and partial summary are written with os.WriteFile, neither fsynced nor atomically renamed, unlike the fsynced append-only process-event log. Because the termination handler deliberately calls os.Exit, a supervised stop as well as an abrupt kill can truncate a boundary write; the failure is closed — decodeStrict refuses and the tranche closes indeterminate — but a routine Ctrl+C at an unlucky microsecond can burn the tranche. A temp-and-rename in writeJSON would remove the class."
- "Residual (P3): a validation output directory whose first custodian process wrote its start event but died before the initial checkpoint is refused wholesale by the progress-implies-checkpoint rule and cannot be resumed; the window is milliseconds and the rule is what prevents a phantom start entry from permanently refusing a later genuine resume."
- "Carried P2-3: applyScheduledResume returns done=true with a non-nil error on stop paths; callers must check the error first."
- "Carried P3 (pre-existing): validateExternalGrade lets a gating manual_required dominate a gating fail, the reverse of the grader overallStatus precedence; unreachable under v5 because no check kind produces manual_required."
- "Carried P2-1 is retired by this payload: classifyScheduledOutputEntries is decomposed into scheduledOutputEntryKind, recordScheduledAssignment, and finishScheduledInventory, and verifyScheduledCheckpoint is split into identity and progress halves. Its successor P2-5 is also retired: verifyResumeAttestationFields is split into identity and condition halves and measures below 5. runScheduledCases and validateExternalGrade remain at gocyclo 14, both pre-existing and unchanged by this payload."

---

## 9. Judgment on each residual carried into this review

- **P3-N1 (coordinated log rewrite).** Correctly characterized and genuinely acceptable. Re-derived independently. **Keep.**
- **The host-clock residual.** Correctly characterized and acceptable. I verified the asymmetry myself. **Keep, review 2's wording.**
- **The post-closure gate's chair-discipline caveat.** Review 3 made its wording conditional on P1-N1 being fixed. **My judgment: with N2 fixed, review 2's original wording stands unchanged.** **If the chair carries N2 or N3 instead of fixing them, the residual must gain a sentence** stating that the runner's corroboration of the named closure record remains heuristic and can accept a document that quotes or lists an ACCEPT. I would not accept that trade, which is why N2 is a finding.
- **The wipe-and-restart caveat.** Correctly characterized and genuinely acceptable. Re-confirmed. **Keep, review 2's sharpened wording verbatim.**
- **The authoring-tranche exemption.** Correctly characterized. Re-verified at `schedule.go:91-93`. **Keep.**
- **P3-1 (log digest never covers the terminal entry).** Correctly characterized and acceptable. **Keep.**
- **P3-2 (five identity fields).** Correctly characterized and acceptable. Re-checked each. **Keep.**
- **P3-5 / P3-N1(new-a).** Correctly characterized with review 2's sharpened wording, and only with it. Still live. **Keep** — but note it is **absent from the inventory's `accepted_residuals` array**, which carries 6 of the 15 entries review 3 recorded. I flag it so the chair transcribes from § 8 above and not from the array.
- **P3-N1(new-c).** Correctly characterized and acceptable, now recorded in the boundary document. **Keep.**
- **Carried P2-3, carried pre-existing P3, carried P2-1 retirement.** All three re-checked and still accurate. **Keep.**

---

## 10. What I did not verify

- I did not run a model, obtain a private grade, or observe any validation outcome. Both gates were `false` throughout.
- I did not construct a live end-to-end validation resume, because that requires the validation gate open. Every resume conclusion rests on `requireResumeAuthorization` and its callees, driven through the package's own fixtures and through one attestation I constructed myself and executed via `go test -overlay`.
- I did not execute the signal path against a live multi-hour scheduled run.
- I did not verify that the process-event log's digest is committed with the custody record *in practice*, because no custody record for a v5 execution exists yet; the runner does not write `…-custody.json`. The digest is recorded in the partial summary, which is what the runner controls.
- I did not attempt to defeat the freeze test or reproduce the world-build digest from source; I confirmed both digests are unmoved and that the gate-document diff touches only the three post-closure fields.
- N1 through N4 I demonstrated empirically rather than by reading: I extracted the parser verbatim into a standalone module, ran it over all 76 files in the working tree and all 75 committed records, ran a 38-case constructed attack suite of my own design, and validated my proposed fixes against the full corpus before prescribing them. The counts are measured, not estimated.
- I made no change of any kind to the repository.

---

## Chair transcription note

Transcribed verbatim from the fourth independent reviewer's record. Before transcription the payload author independently reproduced both P2 findings: N2 was confirmed by driving the real `closureReviewVerdict` against the two word orderings — `ACCEPT, REVISE, or REJECT` refuses, `ACCEPT, REJECT, or REVISE` returns ACCEPT — and N1 by the structural argument, `verifyClosureReviewVerdictLine` returning nil exactly when `closureReviewVerdict` yields ACCEPT, which makes the test's false-refusal branch unreachable. The verdict is **REVISE**, so this payload does not proceed to refreeze; a fifth revision follows. § 8 is conditional on the N1/N2 fixes and is not operative at this verdict.
