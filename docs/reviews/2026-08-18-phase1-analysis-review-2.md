# Phase 1 analysis focused-revision independent review

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-18
- **Verdict:** `ACCEPT`
- **Revision candidate:** `62946f4a1a03ea89636c5b3243f3b4d53b166682`
- **Candidate subject:** `fix: revise Phase 1 analysis`
- **Base candidate:** `a9898cc657bc1825e32c6154adf23e989394fd75`
- **First review:** `docs/reviews/2026-08-18-phase1-analysis-review.md` (`REVISE`)
- **Findings:** None

## Verdict

The focused revision closes both blocking findings and both P3 findings from the first review. Independently recomputed digests match the prompt. The G=8 sign-flip, G=20 bootstrap, G=24 cap-imbalanced fail, and G=24 capability-pass cases behave as required. The four efficiency-only inequalities are unchanged. Runner bytes remain the accepted `b4df919` set. Protected protocol, Arm B, freeze-manifest, and authoring-summary identities are unchanged. `make quality` passed.

No P0, P1, or P2 defect remains in the analysis implementation or report template. Two-sided 95% intervals for descriptive measures remain a completeness limitation and do not block this analysis freeze.

## Closed first-review findings

| First-review finding | Status | Evidence |
| --- | --- | --- |
| P1 bootstrap `PWorse` | Closed | `inferDifference` at G≥20 returns bounds only. `PWorse` stays nil, JSON `omitempty` omits it, and the template prints a harm p-value only when that pointer is set. Harm gates use `harmPBelow`, which is false when `PWorse` is nil. |
| P2 cap-imbalance fail path | Closed | Capability still returns first when `LowerBound > 0`. Cap imbalance can set headline `indeterminate` only when `LowerBound > -0.05` after that return, i.e. inside `(-0.05, 0]`. `LowerBound <= -0.05` is `fail`. |
| P3 `please confirm` | Closed | The help-request regex includes `please confirm`, with `TestDeadEndRecognizesFailedHelpRequestButNotPassingAbsenceReport` covering that phrase. |
| P3 direct −0.10 reference | Closed | The `direct_no_tax` report line states the point estimate versus `-0.100000` and that the harm gate is not a noninferiority claim. |

## Requirement matrix

No regression relative to the first review except check 5, which now passes.

| Check | Result | Notes |
| --- | --- | --- |
| 1. Input integrity | Pass | Unchanged loader contract; overwrite refusal unchanged. |
| 2. ITT and complete-case | Pass | Unchanged. |
| 3. Interrupted pairing / stop rules | Pass | Unchanged. |
| 4. Unresolved thresholds | Pass | Unchanged. |
| 5. Cap-hit imbalance | Pass | Cannot veto capability. Clear fail remains fail. Indeterminate only inside the efficiency-only window. |
| 6. Bootstrap / sign-flip | Pass | G≥20 bounds only; G<20 exact or seeded sign-flip with lower-tail `PWorse`. Seed `20260817`. |
| 7. Forced decisions | Pass | Direct, headline, frontier, and teaching rules match protocol v4. Components cannot create a headline pass. |
| 8. Efficiency-only | Pass | Same four conditions and strict inequalities, evaluated only after capability does not pass. |
| 9. Paired-success token ratios | Pass | Unchanged. |
| 10. Recovery / dead-end / help | Pass | `please confirm` now matches the documented classifier. |
| 11. Frontier take reconstruction | Pass | Unchanged. |
| 12. Report contract | Pass | Sign-flip harm p-value only when present; direct claim names −0.10; still cannot open a gate. |
| 13. Tests and `make quality` | Pass | New tests cover bootstrap `PWorse` absence, cap-imbalanced fail, `please confirm`, and the −0.10 report line. `make quality` passed. |
| 14. Runner bytes | Pass | `git diff b4df919..62946f4 -- experiments/frontier-v1/runner` is empty. |
| 15. Protected artifacts | Pass | Protocol, Arm B, pre-validation manifest, and authoring summary match the expected identities. Manifest remains `partial` with both outcome gates false and seven remaining freeze entries. |

## Independent spot calculations

### G=8 all-minus-one sign-flip

Eight families, each paired difference −1. Exact null size 256. The 0.95 quantile is at index 243, value 0.5. Bounds are `[-1.5, -0.5]`. Harm p-value is `1/256 = 0.00390625`. Matches `TestExactSignFlipRejectsConsistentHarm`.

### G=20 all-plus-one bootstrap

Twenty families, each paired difference +1. Every hierarchical replicate has mean 1, so lower and upper bounds are 1. `PWorse` is nil. The report therefore emits no harm p-value. Matches `TestHierarchicalBootstrapIsDeterministicAtTwentyFamilies`.

### G=24 all-minus-one headline with cap imbalance

Twenty-four families, C always fails, B always succeeds, one C cap hit. Bootstrap bounds are −1. Cap-hit rates differ by `1/24 ≈ 0.0417 > 0.02`. Because `LowerBound = -1 <= -0.05`, the headline is `fail`, not `indeterminate`. Matches `TestHeadlineClearFailureIsNotMadeIndeterminateByCapImbalance`.

### G=24 capability pass despite cap imbalance

The complementary design (C always succeeds, B always fails, one C cap hit) has lower bound 1. Headline is `capability_pass`; the paired-success ratio stays separately indeterminate. Matches `TestHeadlineCapabilityPassIsNotVetoedByCapImbalance`.

### Efficiency-only conjunction (unchanged)

Evaluated only when the C−B success lower bound is not greater than 0, and all four must hold:

1. success lower bound `> -0.05`
2. paired-success C/B token-ratio upper bound determinate and `<= 0.80`
3. C−B dead-end-rate upper bound `< 0.05`
4. C−B timeout-or-cap-hit-rate upper bound `< 0.05`

## Candidate identity

Independently computed LF-normalized UTF-8 SHA-256. All nine match the expected revision digests. Working-tree analysis files are byte-identical to `62946f4`. HEAD `b788ccd` only adds this review prompt.

| Artifact | SHA-256 |
| --- | --- |
| `experiments/frontier-v1/analysis/analysis_test.go` | `sha256:e28f48153b6d4e90cac3011309809e2ea8b275031fe443446dbcd8619f18b029` |
| `experiments/frontier-v1/analysis/claims.go` | `sha256:ff34599bbc1eee374ee149998b124babdd1a1c7993ecba5c11c96b3a0c9a5a0b` |
| `experiments/frontier-v1/analysis/load.go` | `sha256:94fc85440eaa6c8010423a9732d68628dfa1e29f5be4eeaacb948b3a3bc58c19` |
| `experiments/frontier-v1/analysis/main.go` | `sha256:b41d11132a40d56caba4038c7bdbc43a9dd1c3fb99b5bb9b36baffbe2459fc09` |
| `experiments/frontier-v1/analysis/README.md` | `sha256:65f0ff8b929327341e0ccec0555841f30c980efbc67e5c2d27f3de67594e0d94` |
| `experiments/frontier-v1/analysis/report-template.md.tmpl` | `sha256:361641401cc85e8951cf08f5279aac8488713d8dcd52a4df2f2a539a66bdbede` |
| `experiments/frontier-v1/analysis/report.go` | `sha256:947a36d5fdce8379b68a72a36d9c1f2c1ba09495910abd90f300368b2c2c340a` |
| `experiments/frontier-v1/analysis/statistics.go` | `sha256:8263fbbc5bd0e7708b701cb0ed95aea37e18b0cc26fea38aaa9b30215037b6ef` |
| `experiments/frontier-v1/analysis/types.go` | `sha256:26d2d922f84214eee87346ba429cc6b9bd62a911a031087ee870f4dcfad1d200` |

## Protected-artifact statement

This review did not modify the revision candidate, generate or execute a schedule, run a model, or inspect private validation or held-out labels.

| Artifact | SHA-256 | Status |
| --- | --- | --- |
| Accepted runner files at `b4df919070bb9a6d2912662b4a59674b0e25a332` | identical; empty diff through `62946f4` | eligible to freeze as the runner digest |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` | unchanged |
| Arm B document | `sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0` | frozen; unchanged |
| Authoring `summary.json` | `sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405` | retained 7/8 evidence identity unchanged |
| `pre-validation-artifacts.json` | `sha256:c22646e2517aa2cd14baa50e1c6a03d2577127d0344a02d3ea327045ae242eca` | `status: partial`; both outcome gates false; seven remaining freeze entries |

## Freeze boundary

The analysis implementation and report template at the nine digests above, together with the unchanged independently accepted runner bytes from `b4df919070bb9a6d2912662b4a59674b0e25a332`, become eligible for a focused freeze patch.

That patch may record those digests only. It must leave `pre-validation-artifacts.json` `partial`, keep `may_open_validation` and `may_open_held_out` false, and must not freeze or substitute an authoring schedule for the validation schedule. The other open freeze entries remain open: runtime invocation and exact per-arm system prompts; arm A schemas; world definition and world-build digest; schedule digest; grader digest.

## Descriptive-interval classification

Protocol `analysis.confidence` still asks for two-sided 95% intervals on descriptive measures, and those intervals are still absent. They are not inputs to any forced decision. This completeness gap does not block the analysis freeze.

## Exact next allowed action

Write a focused freeze patch that records the nine analysis/report-template digests and the accepted runner digest, keeps both outcome gates false, and leaves the remaining freeze entries open. Do not generate or execute an authoring or validation schedule. Do not run a model outcome. Do not open Gate 1A, validation, or held_out.

Missing sealed families and the validation schedule remain external blockers. This acceptance does not authorize them.

Gate 1A, validation, and held_out remain closed.
