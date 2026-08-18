# Phase 1 analysis and report independent review prompt

Copy the text below into a fresh reviewer session that did not implement the analysis candidate.

```text
Act as the independent evaluation reviewer for the Phoenix frontier-v1 Phase 1 analysis implementation and report template. This is a code, statistics, and reporting review. Do not run a model or inspect any validation or held-out outcome.

Repository: https://github.com/nesste/phoenix
Candidate commit: a9898cc657bc1825e32c6154adf23e989394fd75
Candidate subject: feat: add Phase 1 analysis machinery

Expected LF-normalized UTF-8 SHA-256 digests:

- experiments/frontier-v1/analysis/analysis_test.go: sha256:fdc5d37aebd6bd6628982e74a5f266b53b13a0ce1ad0d05a018dd8ccda41be4a
- experiments/frontier-v1/analysis/claims.go: sha256:126b04aa9b7319d2dc9baa36f96d5a8298060fb947d02b49315e7e35be35a19b
- experiments/frontier-v1/analysis/load.go: sha256:df67056c22c8caec3fdd5e4c022d3fcedbb8c68fa11723ac22cafdeb62a48300
- experiments/frontier-v1/analysis/main.go: sha256:b41d11132a40d56caba4038c7bdbc43a9dd1c3fb99b5bb9b36baffbe2459fc09
- experiments/frontier-v1/analysis/README.md: sha256:f6166cdfd940248f461921f97a2a08c7a84c25ae5735ac0e0ef83f08fd2933a9
- experiments/frontier-v1/analysis/report-template.md.tmpl: sha256:95335bb49321325f05e22a9649faa925e273309c34c1530c9a853664a412e3c6
- experiments/frontier-v1/analysis/report.go: sha256:947a36d5fdce8379b68a72a36d9c1f2c1ba09495910abd90f300368b2c2c340a
- experiments/frontier-v1/analysis/statistics.go: sha256:539e0d7f95dc3f4205dfc2b31f64427f8c5ce0ea76cf647789ed2050d39a4552
- experiments/frontier-v1/analysis/types.go: sha256:26d2d922f84214eee87346ba429cc6b9bd62a911a031087ee870f4dcfad1d200

Read at least:

- experiments/frontier-v1/protocol.json
- docs/decisions/0004-go-no-go-rules.md
- docs/reviews/2026-08-18-phase1-scheduled-runner-review.md
- docs/reviews/2026-08-18-phase1-arm-adapter-status.md
- every file under experiments/frontier-v1/analysis
- experiments/frontier-v1/runner/types.go
- experiments/frontier-v1/runner/scheduled_run.go
- experiments/frontier-v1/runner/runtime.go
- Makefile

Independence and safety rules:

- Do not modify the candidate.
- Do not generate or execute an authoring schedule.
- Do not run any authoring, validation, or held-out model outcome.
- Do not inspect private validation or held-out labels.
- Unit tests, make quality, static analysis, synthetic fixtures, independent calculations, and digest checks are allowed.
- Treat protocol.json, pre-validation-artifacts.json, Arm B, retained authoring evidence, and the reviewed runner bytes as protected.
- Gate 1A, validation, and held_out must remain closed regardless of this review's verdict.

Verify independently:

1. Input integrity requires one result per assigned trial, every manifest case, repetitions 0-2, and A-E exactly once per `(case_id, repetition)`. Result family metadata must match the manifest. Output files cannot be overwritten.
2. ITT success uses all assigned trials. Manual-required, timeout, cap hit, unresolved infrastructure, budget stop, and safety stop are failures. Complete-case sensitivity excludes unresolved/budget/safety assignments but retains manual-required failures.
3. If a safety stop interrupts a pairing key, all five arms in that key are discarded as observed outcomes, set to ITT failure, and counted unresolved. Budget or safety stops make the tranche indeterminate and can never produce a pass.
4. The unresolved thresholds use assigned-arm denominators before rounding: greater than 0.05 in any arm or more than 0.02 imbalance in a registered comparison makes the tranche indeterminate.
5. Cap-hit imbalance above 0.02 makes the affected cost ratio indeterminate but cannot veto a capability pass.
6. At G>=20, the 10,000-replicate paired hierarchical bootstrap resamples families and then cases while retaining paired arms and repetitions. At G<20, repetitions are averaged within case, cases within family, and unweighted family means use all exact sign flips for G<=16 or 100,000 seeded flips for 16<G<20. Seed 20260817 and empirical one-sided bounds are deterministic.
7. The direct C-A harm gate, headline C-B capability/efficiency decision, frontier C-D isolation, and teaching D-E isolation implement the exact protocol-v4 forced decisions. Component results cannot create a headline pass.
8. Efficiency-only is evaluated only after capability does not pass and requires all four registered conditions. ITT cost is descriptive. Capability is not vetoed by token, USD, wall-time, or cap imbalance.
9. Paired-success token ratios use only keys where both arms pass, total tokens within family, analyze family totals on the log scale, add 0.5 only to zero family totals, become indeterminate below ten successes in either arm, and use the family-count-dependent bootstrap/sign-flip rule.
10. Recovery uses the same first non-ok act, later ok alternative, and passing-grade rule for D and E. Wrong-verb, dead-end, timeout/cap, and help-request classifications match the documented operational definitions.
11. Frontier take reconstruction links the displayed call to the next Phoenix executable call by handle, verb, and canonical arguments. It allows the selected current pending call to omit state, counts orientation displays, excludes teaching alternatives, and requires a passing grade for take-and-succeed.
12. The report identifies runtime/model/schedule/world/grader/Arm B, presents ITT and sensitivity results, separates capability from efficiency, states the scope limit, and cannot claim that it opens a gate. Engineering-gate evidence remains a separate required attachment.
13. Synthetic tests materially cover ITT, unresolved and cap rules, interrupted pairings, exact and bootstrap inference, ratios, component decisions, sensitivity, frontier linkage, help requests, and report rendering/overwrite refusal. `make quality` passes with analysis included in formatting, complexity, and duplication checks.
14. `git diff b4df919070bb9a6d2912662b4a59674b0e25a332..a9898cc657bc1825e32c6154adf23e989394fd75 -- experiments/frontier-v1/runner` is empty. Independently compare the runner file digests with the accepted scheduled-runner review.
15. protocol.json, pre-validation-artifacts.json, Arm B, and the retained authoring summary are unchanged. The partial manifest still has both outcome gates false and seven remaining freeze entries.

Return:

1. Verdict: ACCEPT, REVISE, or REJECT.
2. Findings ordered P0 to P3 with exact file and line, consequence, and smallest correction.
3. A requirement matrix for checks 1-15.
4. Independent spot calculations for one G<20 sign-flip case, one G>=20 bootstrap invariant, one interrupted pairing, and the four headline efficiency conditions.
5. Candidate commit and LF-normalized UTF-8 digests for every analysis file and the report template.
6. A protected-artifact statement covering the accepted runner bytes, protocol, Arm B, partial freeze manifest, and retained authoring evidence.
7. If ACCEPT, an explicit freeze boundary: the analysis implementation and report template become eligible to freeze by these digests, and the unchanged accepted runner bytes become eligible to freeze as the runner digest. Do not freeze or substitute an authoring schedule for the validation schedule.
8. The exact next allowed action. Do not authorize a model outcome run, validation schedule freeze without the accepted replacement validation cases, Gate 1A, validation, or held_out.

Separate implementation defects from remaining external blockers. Missing replacement sealed families and the validation schedule are blockers, not reasons to overlook an analysis defect.
```
