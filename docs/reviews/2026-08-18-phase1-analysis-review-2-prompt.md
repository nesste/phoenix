# Phase 1 analysis focused revision review prompt

Copy the text below into a fresh evaluation-review session that did not implement the revision. The reviewer may read the first review but must check the revised bytes independently.

```text
Act as the independent evaluation reviewer for the focused revision of the Phoenix frontier-v1 Phase 1 analysis and report template. This is a code, statistics, and reporting review. Do not run a model or inspect any validation or held-out outcome.

Repository: https://github.com/nesste/phoenix
Revision candidate: 62946f4a1a03ea89636c5b3243f3b4d53b166682
Candidate subject: fix: revise Phase 1 analysis
Base candidate: a9898cc657bc1825e32c6154adf23e989394fd75
First review: docs/reviews/2026-08-18-phase1-analysis-review.md (`REVISE`)

Expected LF-normalized UTF-8 SHA-256 digests:

- experiments/frontier-v1/analysis/analysis_test.go: sha256:e28f48153b6d4e90cac3011309809e2ea8b275031fe443446dbcd8619f18b029
- experiments/frontier-v1/analysis/claims.go: sha256:ff34599bbc1eee374ee149998b124babdd1a1c7993ecba5c11c96b3a0c9a5a0b
- experiments/frontier-v1/analysis/load.go: sha256:94fc85440eaa6c8010423a9732d68628dfa1e29f5be4eeaacb948b3a3bc58c19
- experiments/frontier-v1/analysis/main.go: sha256:b41d11132a40d56caba4038c7bdbc43a9dd1c3fb99b5bb9b36baffbe2459fc09
- experiments/frontier-v1/analysis/README.md: sha256:65f0ff8b929327341e0ccec0555841f30c980efbc67e5c2d27f3de67594e0d94
- experiments/frontier-v1/analysis/report-template.md.tmpl: sha256:361641401cc85e8951cf08f5279aac8488713d8dcd52a4df2f2a539a66bdbede
- experiments/frontier-v1/analysis/report.go: sha256:947a36d5fdce8379b68a72a36d9c1f2c1ba09495910abd90f300368b2c2c340a
- experiments/frontier-v1/analysis/statistics.go: sha256:8263fbbc5bd0e7708b701cb0ed95aea37e18b0cc26fea38aaa9b30215037b6ef
- experiments/frontier-v1/analysis/types.go: sha256:26d2d922f84214eee87346ba429cc6b9bd62a911a031087ee870f4dcfad1d200

Read:

- docs/reviews/2026-08-18-phase1-analysis-review.md
- experiments/frontier-v1/protocol.json
- docs/decisions/0004-go-no-go-rules.md
- every file under experiments/frontier-v1/analysis
- docs/reviews/2026-08-18-phase1-scheduled-runner-review.md

Independence and safety rules:

- Do not modify the revision candidate.
- Do not generate or execute an authoring or validation schedule.
- Do not run any authoring, validation, or held-out model outcome.
- Do not inspect private validation or held-out labels.
- Unit tests, `make quality`, synthetic calculations, and digest checks are allowed.
- Gate 1A, validation, held_out, and all seven open freeze entries remain closed during this review.

Close or reopen every first-review finding:

1. P1 bootstrap `PWorse`: At G>=20, the bootstrap must contribute the registered one-sided bounds without emitting or displaying a `p_worse` value. At G<20, `PWorse` must remain the lower tail of the exact or seeded sign-flip null and the template may label only that value as the one-sided sign-flip harm p-value. Confirm JSON `omitempty` and Markdown rendering do not present a bootstrap tail proportion as a p-value.
2. P2 cap imbalance: Capability must still return first when the success lower bound is greater than zero. When the lower bound is at or below -0.05, the headline must be `fail` even if C and B cap-hit rates differ by more than 0.02. Cap imbalance may make the headline `indeterminate` only when the success lower bound lies inside the efficiency-only window `(-0.05, 0]`, where the paired-success ratio is required.
3. P3 help request: `please confirm` must match the documented help-request classifier, with a regression test.
4. P3 direct reference: the `direct_no_tax` report must show the observed C-A point estimate against -0.10 and state that the harm gate is not a noninferiority claim.

Then rerun the original requirement matrix from checks 1-15 and look for regressions. In particular, independently verify:

- the G=8 all-minus-one sign-flip still has p=1/256 and bounds [-1.5, -0.5];
- the G=20 all-plus-one bootstrap has lower and upper bounds equal to 1 and no harm p-value;
- a G=24 all-minus-one headline with cap imbalance remains `fail`;
- a capability lower bound above zero remains `capability_pass` despite cap imbalance;
- the four efficiency-only conditions and their strict inequalities are unchanged;
- runner bytes remain identical to accepted candidate b4df919070bb9a6d2912662b4a59674b0e25a332;
- protocol.json, Arm B, pre-validation-artifacts.json, and the retained authoring summary remain unchanged;
- `make quality` passes.

Expected protected identities:

- protocol.json: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9
- Arm B: sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0
- pre-validation-artifacts.json: sha256:c22646e2517aa2cd14baa50e1c6a03d2577127d0344a02d3ea327045ae242eca
- retained authoring summary.json: sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405

Return:

1. Verdict: ACCEPT, REVISE, or REJECT.
2. Findings ordered P0 to P3 with exact file and line, consequence, and smallest correction.
3. A closed/open matrix for the two blocking findings and two P3 findings from the first review.
4. The original 15-check requirement matrix, noting any regression.
5. Independent spot calculations for the four cases listed above.
6. Candidate commit and LF-normalized UTF-8 SHA-256 digests for all nine analysis artifacts.
7. A protected-artifact statement covering the runner, protocol, Arm B, partial freeze manifest, and retained authoring evidence.
8. If ACCEPT, state the exact freeze boundary: these analysis implementation and report-template digests, plus the unchanged independently accepted runner bytes, become eligible for a focused freeze patch. Do not freeze or substitute an authoring schedule for the validation schedule.
9. The exact next allowed action. Do not authorize a schedule or model outcome run, Gate 1A, validation, or held_out merely because the analysis revision is accepted.

The first review noted that two-sided intervals for descriptive measures remain a completeness limitation rather than a forced-decision defect. Reassess that classification independently and state whether it blocks this analysis freeze.
```
