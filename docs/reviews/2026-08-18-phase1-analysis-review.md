# Phase 1 analysis and report independent review

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-18
- **Verdict:** `REVISE`
- **Candidate commit:** `a9898cc657bc1825e32c6154adf23e989394fd75`
- **Candidate subject:** `feat: add Phase 1 analysis machinery`

## Verdict

The candidate is a substantial and mostly faithful implementation of the protocol-v4 analysis contract. Input integrity, ITT accounting, unresolved-rate gates, pairing-key safety-stop handling, family-count-dependent interval construction, paired-success ratios, component isolations, sensitivity views, frontier linkage, and overwrite refusal are present. `make quality` passed with analysis included in format, complexity, and duplication checks. The accepted runner bytes are unchanged.

Two defects block acceptance. At `G>=20`, the shared `PWorse` field is not a harm p-value, and the report labels it as one on every claim, including the G=24 headline. Cap-hit imbalance, which must not veto a capability pass, also converts an independently failed headline into `indeterminate`. Neither defect can create a capability pass. Both must be corrected before freeze.

## Findings

### [P1] Bootstrap `PWorse` is not a harm p-value — `experiments/frontier-v1/analysis/statistics.go:36`

At `G>=20`, `inferDifference` sets `PWorse` to the share of hierarchical-bootstrap replicates that are `<= 0`. That quantity is small when the intervention is better and large when it is worse. The report template then labels the same field "One-sided harm p-value" for every claim.

Consequence: on the registered G=24 headline, a capability pass (`LowerBound > 0`) displays a near-zero "harm p-value". A uniformly worse intervention displays a harm p-value near 1. The same inverted quantity would fail-open a `PWorse < 0.05` harm gate if any isolation ever reached 20 families. Under the committed allocation the four forced *decisions* still use the correct objects (sign-flip p-values at G=8/16, and the headline lower bound rather than `PWorse`), so this cannot create a capability pass. It still misreports the primary claim.

Smallest correction: keep sign-flip `PWorse` as the lower tail of the null versus the observed family-mean difference. For the bootstrap, report a separately named probability or a one-sided p-value whose small values mean harm (`P(θ* >= 0)` under a negative observed effect), and do not label the headline field as a harm p-value.

```36:40:experiments/frontier-v1/analysis/statistics.go
		return inference{
			Method: "10000-replicate paired hierarchical bootstrap", Families: len(families), Pairs: len(values),
			Point: meanPairDifference(values), LowerBound: quantile(distribution, 0.05),
			UpperBound: quantile(distribution, 0.95), PWorse: tailProbability(distribution, 0),
		}
```

```46:46:experiments/frontier-v1/analysis/report-template.md.tmpl
- One-sided harm p-value: {{optional .PWorse}}
```

### [P2] Cap-hit imbalance converts a headline fail into indeterminate — `experiments/frontier-v1/analysis/claims.go:152`

Capability is returned immediately when the success lower bound is greater than zero, so cap-hit imbalance cannot veto a capability pass. After that return, any cap-hit imbalance forces `indeterminate`, even when the success lower bound is already `<= -0.05` and efficiency-only is independently impossible.

Consequence: a clear capability failure that should force "kill or narrow Phoenix" becomes an unevaluable tranche whenever C and B differ by more than 0.02 in cap-hit rate. Protocol `cap_imbalance` makes cost-ratio claims indeterminate, not the success claim. This cannot create a pass; it can block a valid fail.

Smallest correction: if `LowerBound <= -0.05`, decide `fail` regardless of cap imbalance. Treat cap imbalance as headline-indeterminate only inside the efficiency-only window `(-0.05, 0]`, where the paired-success ratio is actually required.

```148:159:experiments/frontier-v1/analysis/claims.go
	if inference.LowerBound > 0 {
		claim.Decision, claim.Reason = "capability_pass", "one-sided 95% lower bound for C-B ITT success is greater than zero"
		return claim
	}
	efficiency := inference.LowerBound > -0.05 && ratio.Determinate && ratio.UpperBound != nil && *ratio.UpperBound <= 0.80 && deadEndInference.UpperBound < 0.05 && capInference.UpperBound < 0.05
	if efficiency {
		claim.Decision, claim.Reason = "efficiency_only", "capability did not pass; all four precommitted efficiency-only conditions hold"
	} else if capImbalanced {
		claim.Decision, claim.Reason = "indeterminate", "cap-hit imbalance makes the required paired-success token ratio indeterminate"
	} else {
		claim.Decision, claim.Reason = "fail", "capability lower bound is not above zero and the efficiency-only conjunction does not hold"
	}
```

### [P3] Help-request regex omits `please confirm` — `experiments/frontier-v1/analysis/load.go:213`

The analysis README treats "confirm" as a help-request trigger alongside provide, clarify, and specify. The regex allows `please (provide|clarify|specify)` and `(can|could|would) you ... confirm`, but not `please confirm`.

Consequence: some failed help-seeking finals are not classified as dead ends. This is descriptive for capability and is only one of four efficiency-only inputs.

Smallest correction: add `confirm` to the `please (...)` group, matching the documented operational definition.

### [P3] Direct harm report does not present the −0.10 descriptive comparator — `experiments/frontier-v1/analysis/report-template.md.tmpl:43`

Protocol `direct_no_tax` requires reporting the C−A difference against −0.10 descriptively and forbids a noninferiority claim. The template prints the point estimate and bounds but never names −0.10.

Smallest correction: add one descriptive line on `direct_no_tax` that states the observed difference against −0.10 and repeats that this is not a noninferiority claim.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Input integrity | Pass | One result per assigned trial; cases × 3 × 5; repetitions 0–2; A–E exactly once per `(case_id, repetition)`; family must match the manifest; JSON and Markdown outputs refuse overwrite. |
| 2. ITT and complete-case | Pass | ITT uses every assigned trial and the runner's `itt_success`. Manual-required remains `status=fail` with `CompleteCase=true`. Unresolved, budget-stopped, and safety-stopped are excluded from complete-case. |
| 3. Interrupted pairing / stop rules | Pass | A safety stop on any arm of a pairing key zeroes all five ITT outcomes, marks them unresolved, and makes the tranche indeterminate. Budget and safety stops also trip `summary.Status` / `StopReason`. |
| 4. Unresolved thresholds | Pass | Arm rates use assigned denominators with `>` 0.05 and pairwise imbalance `>` 0.02, before rounding. |
| 5. Cap-hit imbalance | Fail | Cannot veto capability (tested). Incorrectly replaces an independent `fail` with `indeterminate` (P2). |
| 6. Bootstrap / sign-flip | Pass, with P1 on `PWorse` | Family-then-case hierarchical bootstrap at G≥20; unweighted family means with exact 2^G or 100,000 seeded flips below 20; seed `20260817`; empirical 5th/95th percentiles. Independent 8-family and 20-family spot checks match. |
| 7. Forced decisions | Pass, with P2 on headline fail | Direct, capability, frontier, and teaching rules match protocol v4. Components cannot create a headline pass. Headline fail path is too conservative under cap imbalance. |
| 8. Efficiency-only | Pass | Evaluated only after capability does not pass. Four registered conditions. ITT cost is descriptive. Token/USD/wall-time/cap imbalance cannot veto capability. |
| 9. Paired-success token ratios | Pass | Both-arm successes only; family totals; log scale; +0.5 only on zero family totals; indeterminate below ten arm successes; family-count-dependent interval. |
| 10. Recovery / dead-end / help | Pass | First non-ok act, later ok act, passing grade, shared by D and E. Wrong-verb, dead-end, timeout/cap, and help-request definitions match the analysis README, with P3 on `please confirm`. |
| 11. Frontier take reconstruction | Pass | Next executable `mcp__phoenix__act` matched by handle, verb, and canonical args; displayed `state` ignored so omitted state is allowed; orientations counted; teaching alternatives are not read from `frontier`; take-and-succeed requires a passing grade. |
| 12. Report contract | Pass, with P1 labeling | Identifies runtime/model/schedule/world/grader/Arm B; ITT and sensitivity; capability vs efficiency; scope limit; cannot open a gate. Engineering-gate evidence is stated as a separate attachment in the analysis README. |
| 13. Tests and `make quality` | Pass | Synthetic tests cover the listed topics. `make quality` passed with analysis in format, gocyclo, and dupl. Residual gaps: no efficiency-only conjunction test, no 16<G<20 Monte Carlo test, no JSON overwrite test. |
| 14. Runner bytes | Pass | `git diff b4df919..a9898cc -- experiments/frontier-v1/runner` is empty. Independently recomputed runner digests match the accepted scheduled-runner review. |
| 15. Protected artifacts | Pass | `protocol.json`, Arm B, authoring `summary.json`, and decision 0004 are unchanged. `pre-validation-artifacts.json` remains `partial` with both outcome gates false and seven remaining freeze entries. |

## Independent spot calculations

### G<20 exact sign-flip

Eight families, each with paired difference −1 (C always fails, comparator always succeeds). Unweighted family means are eight copies of −1. The exact null is all 256 sign assignments of that vector. The 0.95 empirical quantile uses index `ceil(0.95 × 256) − 1 = 243`, which falls in the 28 outcomes with six positive signs, value 0.5. Bounds are `point ± 0.5 = [-1.5, -0.5]`. Harm p-value is `1/256 = 0.00390625`. This matches `TestExactSignFlipRejectsConsistentHarm` and would force `reject_surface` / `remove_frontier` / `replace_teaching_refusals` as applicable.

### G>=20 bootstrap invariant

Twenty families, each with paired difference +1. Every family-then-case resample still has mean 1, so the 5th and 95th percentiles are 1. The capability lower bound is greater than zero. The current `PWorse = P(θ* <= 0)` is 0, which the template would print as a harm p-value of 0.000000 on a capability pass (P1).

### Interrupted pairing

One manifest case, three repetitions, five arms (15 assigned trials). Repetition 1 arm D is `safety_stopped`. `finalizePairings` counts that key incomplete, sets all five repetition-1 arms to ITT failure, `CompleteCase=false`, `Unresolved=true`, and leaves repetitions 0 and 2 unchanged. `incomplete > 0` makes the tranche indeterminate. Independently confirmed by `TestSafetyStopDiscardsWholeFiveArmPairing`.

### Four headline efficiency-only conditions

Evaluated only when the C−B ITT success lower bound is not greater than 0, and all four must hold:

1. success lower bound `> -0.05`
2. paired-success C/B token-ratio upper bound determinate and `<= 0.80`
3. C−B dead-end-rate upper bound `< 0.05`
4. C−B timeout-or-cap-hit-rate upper bound `< 0.05`

ITT tokens, USD, wall time, and cap-hit imbalance cannot create a capability pass. Condition 2 becoming indeterminate must not override an already failed condition 1 (P2).

## Candidate identity

Independently computed LF-normalized UTF-8 SHA-256 for the files named in the review prompt. All nine match the expected digests.

| Artifact | SHA-256 |
| --- | --- |
| `experiments/frontier-v1/analysis/analysis_test.go` | `sha256:fdc5d37aebd6bd6628982e74a5f266b53b13a0ce1ad0d05a018dd8ccda41be4a` |
| `experiments/frontier-v1/analysis/claims.go` | `sha256:126b04aa9b7319d2dc9baa36f96d5a8298060fb947d02b49315e7e35be35a19b` |
| `experiments/frontier-v1/analysis/load.go` | `sha256:df67056c22c8caec3fdd5e4c022d3fcedbb8c68fa11723ac22cafdeb62a48300` |
| `experiments/frontier-v1/analysis/main.go` | `sha256:b41d11132a40d56caba4038c7bdbc43a9dd1c3fb99b5bb9b36baffbe2459fc09` |
| `experiments/frontier-v1/analysis/README.md` | `sha256:f6166cdfd940248f461921f97a2a08c7a84c25ae5735ac0e0ef83f08fd2933a9` |
| `experiments/frontier-v1/analysis/report-template.md.tmpl` | `sha256:95335bb49321325f05e22a9649faa925e273309c34c1530c9a853664a412e3c6` |
| `experiments/frontier-v1/analysis/report.go` | `sha256:947a36d5fdce8379b68a72a36d9c1f2c1ba09495910abd90f300368b2c2c340a` |
| `experiments/frontier-v1/analysis/statistics.go` | `sha256:539e0d7f95dc3f4205dfc2b31f64427f8c5ce0ea76cf647789ed2050d39a4552` |
| `experiments/frontier-v1/analysis/types.go` | `sha256:26d2d922f84214eee87346ba429cc6b9bd62a911a031087ee870f4dcfad1d200` |

Working tree analysis files are byte-identical to `a9898cc`. HEAD `9b31513` only adds this review prompt.

## Protected-artifact statement

This review did not modify the candidate, generate or execute a schedule, run a model, or inspect private validation or held-out labels. Independently recomputed identities:

| Artifact | SHA-256 | Status |
| --- | --- | --- |
| Accepted runner files listed in the scheduled-runner review | identical to that review | unchanged since `b4df919070bb9a6d2912662b4a59674b0e25a332` |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` | matches the scheduled-runner pin |
| `docs/decisions/0004-go-no-go-rules.md` | `sha256:883678b0799c76509ac306e7c38f322c35dab12d8dfe19b71228045b22652e9e` | matches the scheduled-runner pin |
| Arm B document | `sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0` | frozen; matches the Arm B acceptance |
| Authoring `summary.json` | `sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405` | retained 7/8 evidence identity unchanged |
| `pre-validation-artifacts.json` | `sha256:c22646e2517aa2cd14baa50e1c6a03d2577127d0344a02d3ea327045ae242eca` | `status: partial`; `may_open_validation: false`; `may_open_held_out: false`; seven remaining freeze entries |

The seven remaining freeze entries are: runtime invocation and exact per-arm system prompts; arm A schemas; world definition and world-build digest; runner digest; schedule digest; grader digest; analysis implementation and report template.

## Freeze boundary

Not authorized. The analysis implementation and report template are **not** eligible to freeze by the digests above. The accepted runner bytes remain unchanged and remain eligible to freeze *after* an accepted analysis revision, not from this verdict. Do not freeze or substitute an authoring schedule for the validation schedule. Do not write these digests into `pre-validation-artifacts.json` from this review.

## External blockers, separate from this verdict

These are not reasons to overlook the analysis defects:

- replacement sealed validation and held_out families are still absent from this implementation workspace;
- the validation schedule does not exist and cannot be replaced by an authoring schedule;
- seven `artifact_freeze.before_validation` entries remain, including analysis and the runner digest.

Task 0.6 acceptance in a separate evaluator checkout does not import those families here and does not open a gate.

## Exact next allowed action

Revise the analysis implementation and report template to correct P1 and P2 (P3 optional in the same patch). Re-submit those analysis files for independent review. Do not freeze analysis or runner digests. Do not generate or execute an authoring or validation schedule as part of that revision. Do not run a model outcome. Do not open Gate 1A, validation, or held_out.

## Accepted limitations of this review

- Conclusions remain limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes.
- Frontier linkage was checked against the analysis test's stream-json fixture (string `tool_result` content, `state` as a sibling of `args`). This review did not inspect retained authoring runtime transcripts.
- Cap-hit classification remains the accepted runner's `cap_hit` / termination encoding.
- Protocol `analysis.confidence` asks for two-sided 95% intervals on descriptive measures; those intervals are absent. That is a completeness gap, not a forced-decision defect, and was not raised as a separate P2.

Gate 1A, validation, and held_out remain closed. This review does not authorize opening them.
