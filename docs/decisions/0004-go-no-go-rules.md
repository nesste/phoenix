# 0004: Frontier experiment go/no-go rules

- **Status:** Protocol v3 frozen; independently accepted
- **Date:** 2026-08-17
- **Protocol:** `experiments/frontier-v1/protocol.json`
- **Review 1:** `docs/reviews/2026-08-17-task-0.5-review-1.md`
- **Second-review prompt:** `docs/reviews/2026-08-17-task-0.5-evaluator-prompt-v2.md`

## Decision

Phoenix uses a fixed, paired experiment with no outcome-based stopping. Protocol v3 keeps one confirmatory Phase 1 claim: C versus B intention-to-treat task success. The direct-task, frontier, and teaching-refusal comparisons are pre-registered component isolations. An isolation may remove or reject its component, but it cannot create a headline pass.

Review 1 returned `REVISE`; the independent confirmation review accepted the corrected v3 candidate on 2026-08-17. This freeze patch records that acceptance and sets `protocol.json` to `status: frozen`, `frozen: true`. Acceptance freezes the design but does not open validation or held_out outcomes.

## Runtime and trial contract

All arms use Claude Code 2.1.229, `claude-sonnet-5`, low effort, standard service, fast mode off, stdio, and surface-only access. Each assigned trial receives a fresh model context, sandbox, and isolated world. Session persistence is disabled.

Each case-arm pair runs three times. Claude Code does not expose a model seed for this configuration, so comparisons pair arms by case and repetition. A family-blocked Williams Latin square orders launches of `(case_id, repetition, arm)` pairing keys; it is not a within-session crossover. The schedule seed is 20260817, and the schedule digest must be committed before a tranche opens.

A trial has these hard limits:

- 12 model turns;
- 180 seconds;
- USD 0.15;
- one infrastructure retry, only for a provider, process-launch, or MCP connection failure before the first model token.

A timeout, malformed call, tool error, refusal, agent error, turn-limit hit, or cost-cap hit is an intention-to-treat failure and is not retryable. Retry attempts retain the original pairing key and analysis cluster.

Budget stops occur only at completed pairing-key boundaries: all Phase 1 arms for that case and repetition receive terminal outcomes, or none do. Incomplete pairing keys count toward the unresolved quota. The 5% unresolved limit and 2-percentage-point imbalance rule use assigned trials in each arm before rounding. If cap-hit rates differ by more than 0.02 between compared arms, cost-ratio claims are indeterminate. None of these conditions can create a pass.

## Sealed-tranche size and allocation

Each validation and held_out tranche contains exactly the precommitted minimum design: 120 cases from 24 generating families, five cases per family, three repetitions per case-arm pair, at least eight families per class and inferential subset, and no generating family shared across tranches.

The allocation fixes the cluster incidence:

- eight direct families, each with three `direct`, one `absence`, and one `stale_frontier` case;
- eight recovery families, each with three `recovery`, one `far_discovery`, and one `temptation` case;
- eight mix families containing, in total, 12 `cascade`, four `far_discovery`, four `temptation`, four `absence`, four `stale_frontier`, and 12 `adversarial_text` cases, with at least one `cascade` case in every mix family;
- all eight mix families therefore contribute to the frontier subset, so that subset has 16 families: the eight recovery families plus all eight mix families.

Class totals are 24 `direct`, 24 `recovery`, and 12 for each other class.

Under binary rate 0.5, family ICC 0.10, and repetition ICC 1.0, the design has approximately 62% power for a 15-point headline effect. Its approximate 80% headline MDE against zero is 19 points. The USD 300 validation ceiling cannot support 80% power for the former Holm-adjusted LCB floors, so v3 makes no such claim.

## Arm contract

| Arm | Surface | Additional behavior |
| --- | --- | --- |
| A | flat tools with complete upfront schemas | none |
| B | identical to A | frozen static authoring-only skill or `CLAUDE.md` |
| C | Phoenix handles plus structured `act` | frontier and teaching refusals |
| D | identical to C | frontier suppressed; teaching refusals retained |
| E | identical to D | teaching refusals replaced with plain typed errors |
| D' | Phase 2 surface identical to counted C | authored static ranking instead of counted ranking |

All arms share the case, sandbox, verb implementations, typed payloads, side effects, limits, model configuration, and cost accounting. The upfront flat schemas in A/B and constant schema in C/D/E are irreducible surface differences.

Arm B is operationally strong. Its frozen document must name every starter-world verb with arguments and a one-line when-to-use; include one acceptable worked path per authoring class; include a recovery recipe for every authoring refusal case; use authoring evidence only; receive a committed human-factors completeness review against the authoring labels; and freeze by digest before validation. It uses the same implementations, payloads, effects, limits, and surface-only access as A.

## Analysis

The assigned case-arm-repetition trial is the unit of analysis. Comparisons pair arms by case and repetition and cluster by generating family and case.

For subsets with at least 20 families, one-sided bounds use a 10,000-replicate paired hierarchical bootstrap: resample families, then cases within families, while retaining arm pairs and repetitions. For fewer than 20 families, the primary method is a sign-flip permutation of unweighted family means of paired case-level differences: all `2^G` flips when `G <= 16`, or 100,000 seeded Monte Carlo flips when `16 < G < 20`. A percentile bootstrap is not the primary bound below 20 families.

The headline is the single Phase 1 primary claim at one-sided alpha 0.05. The three component comparisons use the same family-count-dependent method as isolations, not Holm-adjusted bounds. Phase 2 learning is a separate single claim.

Ratios operate on family totals and are analyzed on the log scale. A zero family total receives a +0.5 token continuity correction. If either compared arm has fewer than ten successful trials, ratio claims are indeterminate. The report also includes complete-case, one-vote-per-case, and one-vote-per-family sensitivity views; disagreement narrows the conclusion or makes it indeterminate.

## Forced decisions

| Claim | Rule | Forced verdict |
| --- | --- | --- |
| direct-task harm gate, C versus A | On all assigned direct trials, fail if the one-sided family-mean permutation test that C is worse than A has `p < 0.05`. Otherwise do not block the headline, do not claim noninferiority, and report the difference against -0.10 descriptively. | repair or reject the Phoenix surface |
| headline capability, C versus B | On all assigned trials, pass capability only if the one-sided alpha-0.05 lower bound for C minus B ITT success is greater than zero. ITT tokens, wall time, USD, and paired-success tokens are descriptive and cannot veto this pass. | if this and efficiency-only fail, kill or narrow Phoenix |
| headline efficiency-only | Test only when the primary bound is at most zero. Require success lower bound greater than -0.05, paired-success token-ratio upper bound at most 0.80, C-minus-B dead-end-rate upper bound below 0.05, and C-minus-B timeout-or-cap-hit upper bound below 0.05. | label only as an efficiency product; never claim capability |
| frontier isolation, C versus D | On all assigned cascade, far-discovery, recovery, and temptation trials in the 16-family subset fixed by the allocation, remove the frontier if the family-mean point estimate is at most zero or C is worse at one-sided `p < 0.05`. Tokens are descriptive. | remove the frontier; retain the surface only if the headline passed |
| teaching-refusal isolation, D versus E | On all assigned recovery trials, define recovery from the first failed or declined state-sensitive act through a later successful alternative and successful grade. Replace teaching if D minus E recovery is at most zero, D is worse at one-sided `p < 0.05`, or the downstream ITT-success lower bound is below -0.05. | use plain typed errors |
| learning, counted C versus D' | On all held_out cases, keep counting disabled if the one-sided alpha-0.05 success lower bound is at most zero. Efficiency is descriptive. | keep authored ranking and counting disabled |

The approximate 80% MDEs against zero are 0.19 for the headline and learning claims, 0.26 for the frontier isolation at its allocated 16 families and 60 cases, and 0.39 for both the direct harm gate and the teaching isolation at their effective sample sizes of 20. Direct-gate power is about 0.16 for a 10-point tax and 0.35 for a 20-point tax; the evaluator-specified 0.30 effect size is not an 80% MDE. These values describe sensitivity, not pass floors.

## Phase 2

Phase 2 disables online updates and freezes the daemon, verbs, authored rules, candidate weights, and training-episode digest. Its arms are counted C, D', and B. Counted C versus D' is the learning claim. Counted C versus B repeats the Phase 1 headline success threshold without consuming Phase 1 alpha. If Phase 1 passed but the held_out replication fails, held_out is indeterminate and cannot pass.

## Engineering gates

The inferential claims do not override these hard checks:

- incremental standing input cost versus empty MCP is at most 600 tokens at 10, 100, and 1,000 verbs;
- standing-token spread across those world sizes is at most five tokens;
- on Linux, daemon-only act latency has median at most 2 ms and p95 at most 10 ms;
- daemon time is at most 5% of median successful trial wall time;
- malformed-call rate is at most 5% across at least 20 fresh one-call sessions, with the Wilson 95% interval reported.

Failure blocks the surface claim. Repairing a failed gate after seeing validation outcomes burns the tranche.

## Budget

The hard budget remains USD 565:

| Work | Arithmetic | Ceiling |
| --- | --- | ---: |
| authoring | stated allocation | USD 75 |
| validation | 120 cases x 3 repetitions x 5 arms x USD 0.15 x 1.10 = USD 297 | USD 300 |
| held_out | 120 cases x 3 repetitions x 3 arms x USD 0.15 x 1.10 = USD 178.20 | USD 180 |
| surface and protocol | stated allocation | USD 10 |

The 10% allowance is pooled infrastructure capacity, not authorization for every trial to retry. A budget or safety stop is indeterminate, never a pass.

## Freeze boundary

Before validation opens, commit separate digests for the runtime invocation, system prompt, A schemas, B document, world build, runner, schedule, grader, analysis implementation, and report template. These artifacts may use authoring outcomes only.

Review 1 identified the v2 protocol by LF-normalized UTF-8 SHA-256 and returned `REVISE`; it did not issue an acceptance record. The independent confirmation review accepted the corrected v3 candidate, so the protocol is frozen. An independent reviewer accepted Task 0.2 on 2026-08-18. Validation remains closed until independent unopened sealed families exist and the Task 0.6 Phase 0 review is complete.

## Independent acceptance record

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-17
- **Verdict:** `ACCEPT`
- **Accepted protocol candidate, LF-normalized UTF-8 SHA-256:** `sha256:4d13c140742ba462d68af865d4ed12b08e6c2e1aba4f652fed4e35adf782a47d`
- **Accepted decision candidate, LF-normalized UTF-8 SHA-256:** `sha256:399d337ef957ff1a43e7b384a03a7ddfe1e3a360b9f5337ae1d662338dfe3849`

Accepted limitations:

- conclusions are limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes;
- expected headline power is about 62% at a 15-point effect, and the approximate 80% MDE against zero is 0.19;
- the USD 300 validation ceiling cannot support 80% power for a 15-point Holm-adjusted LCB-floor design;
- efficiency-only is a conservative conjunction; at a true zero difference, power to clear the 0.05 success or dead-end bounds is about 0.16.

Remaining execution blockers:

- independent unopened validation and held_out families;
- Task 0.6 Phase 0 review.

Until all blockers are cleared, `may_open_validation` and `may_open_held_out` remain false. This freeze does not authorize sealed-family creation, private labels, authoring outcomes, or validation or held_out outcome runs.
