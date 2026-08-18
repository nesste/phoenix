# 0004: Frontier experiment go/no-go rules

- **Status:** Protocol v4 independently accepted and frozen; protocol v3 acceptance is historical and superseded
- **Date:** 2026-08-17
- **Protocol:** `experiments/frontier-v1/protocol.json`
- **Review 1:** `docs/reviews/2026-08-17-task-0.5-review-1.md`
- **Second-review prompt:** `docs/reviews/2026-08-17-task-0.5-evaluator-prompt-v2.md`
- **Protocol v4 acceptance:** `docs/reviews/2026-08-18-protocol-v4-review-2.md`

## Decision

Protocol v4 and decision 0007 supersede the v3 surface and corpus contract after authoring exposed an initial-activation gap and an unrepresentable stale-state event. The second independent review accepted candidate `f889f13f0c514fa5108a1e392701ebeadc4376f7` on 2026-08-18 with no findings. Protocol v4 is frozen by the follow-up acceptance-record patch. Gate 1A, validation, held_out, and all sealed tranches remain closed.

The v3 `ACCEPT` verdict and freeze described later in this decision are historical facts about the superseded v3 bytes only. They do not freeze or accept v4, authorize reuse of retired v3 candidates, or satisfy any current blocker.

## Current v4 amendment contract

Phoenix keeps one constant structured `act` surface for C, D, and E. Before a session's first executable act, an orientation input may activate from world-authored rules. After the first executable act, orientation may return only the current pending frontier or teaching-refusal alternative. With no pending call it returns no calls; it cannot reactivate from intent. D therefore cannot recreate a suppressed frontier, and E cannot recreate either a suppressed frontier or teaching alternative. Omitted state restoration is limited to a state-bound call in the current pending frontier.

State events are a shared trial-harness rule, not a Phoenix-arm treatment. The same case event plan is supplied to A, B, C, D, and E. Each arm adapter indexes only executable actions and applies a declared file replacement after computing the numbered action's result but before delivering it to the agent. Orientations do not advance the index. This preserves the intended cross-arm state while allowing a pending state-bound call to become reproducibly stale.

Phase 1 prompts are pinned per arm. A and B receive the common tool-use prompt only. C, D, and E receive that same prompt plus the Phoenix bootstrap-intent instruction. No other standing workflow prose is added.

| Arm | Standing surface | Orientation | Frontier | Teaching refusal | State events |
| --- | --- | --- | --- | --- | --- |
| A | flat tools with complete upfront schemas | none | none | none | shared harness |
| B | identical to A plus frozen static authoring-only document | none | none | none | shared harness |
| C | Phoenix handles plus structured `act` | bootstrap-only fresh activation | present | present | shared harness |
| D | identical to C | bootstrap-only; pending refusal alternatives only after an act | suppressed | present | shared harness |
| E | identical to D | bootstrap-only; no calls after an act | suppressed | plain typed errors | shared harness |

The fixed paired design, claims, budgets, and no-outcome-stopping rules below remain the intended statistical contract where protocol v4 retains them. Any wording that reports a v3 freeze or acceptance is explicitly historical.

## Protocol v4 independent acceptance

The accepted candidate is commit `f889f13f0c514fa5108a1e392701ebeadc4376f7`. The acceptance applies to the protocol-design candidate identified by the five LF-normalized UTF-8 SHA-256 digests recorded in `protocol.json` and the v4 review record. It does not convert the 7/8 authoring run into a gate result.

The accepted limitations include the preimplementation A/B/D/E adapters, shared flat-tool state-event application, and the requirement that D/E suppression remove frontier or refusal calls before pending state is recorded. Cascade remains noisy and its restored label cannot be weakened.

Protocol acceptance authorizes this acceptance-record patch only. It does not authorize Gate 1A, candidate sealing, validation, or held_out. Those remain blocked on new independent families, every pre-validation artifact digest, the Arm B document and human-factors review, the Phase 1 A-E runner, and Task 0.6 acceptance.

## Historical v3 design and acceptance (superseded)

The remainder of this decision records the independently accepted v3 design for auditability. It uses the historical present tense in places because those paragraphs are preserved from the accepted record. None of it changes the current v4 status.

Phoenix used a fixed, paired experiment with no outcome-based stopping. Protocol v3 kept one confirmatory Phase 1 claim: C versus B intention-to-treat task success. The direct-task, frontier, and teaching-refusal comparisons were pre-registered component isolations. An isolation could remove or reject its component, but it could not create a headline pass.

Review 1 returned `REVISE`; the independent confirmation review accepted the corrected v3 candidate on 2026-08-17. That historical freeze patch recorded the v3 acceptance and set the then-current `protocol.json` to `status: frozen`, `frozen: true`. It did not accept v4 and never opened validation or held_out outcomes.

## Historical v3 runtime and trial contract

All arms use Claude Code 2.1.229, `claude-sonnet-5`, low effort, standard service, fast mode off, stdio, and surface-only access. Each assigned trial receives a fresh model context, sandbox, and isolated world. Session persistence is disabled.

Each case-arm pair runs three times. Claude Code does not expose a model seed for this configuration, so comparisons pair arms by case and repetition. A family-blocked Williams Latin square orders launches of `(case_id, repetition, arm)` pairing keys; it is not a within-session crossover. The schedule seed is 20260817, and the schedule digest must be committed before a tranche opens.

A trial has these hard limits:

- 12 model turns;
- 180 seconds;
- USD 0.15;
- one infrastructure retry, only for a provider, process-launch, or MCP connection failure before the first model token.

A timeout, malformed call, tool error, refusal, agent error, turn-limit hit, or cost-cap hit is an intention-to-treat failure and is not retryable. Retry attempts retain the original pairing key and analysis cluster.

Budget stops occur only at completed pairing-key boundaries: all Phase 1 arms for that case and repetition receive terminal outcomes, or none do. Incomplete pairing keys count toward the unresolved quota. The 5% unresolved limit and 2-percentage-point imbalance rule use assigned trials in each arm before rounding. If cap-hit rates differ by more than 0.02 between compared arms, cost-ratio claims are indeterminate. None of these conditions can create a pass.

## Historical v3 sealed-tranche size and allocation

Each validation and held_out tranche contains exactly the precommitted minimum design: 120 cases from 24 generating families, five cases per family, three repetitions per case-arm pair, at least eight families per class and inferential subset, and no generating family shared across tranches.

The allocation fixes the cluster incidence:

- eight direct families, each with three `direct`, one `absence`, and one `stale_frontier` case;
- eight recovery families, each with three `recovery`, one `far_discovery`, and one `temptation` case;
- eight mix families containing, in total, 12 `cascade`, four `far_discovery`, four `temptation`, four `absence`, four `stale_frontier`, and 12 `adversarial_text` cases, with at least one `cascade` case in every mix family;
- all eight mix families therefore contribute to the frontier subset, so that subset has 16 families: the eight recovery families plus all eight mix families.

Class totals are 24 `direct`, 24 `recovery`, and 12 for each other class.

Under binary rate 0.5, family ICC 0.10, and repetition ICC 1.0, the design has approximately 62% power for a 15-point headline effect. Its approximate 80% headline MDE against zero is 19 points. The USD 300 validation ceiling cannot support 80% power for the former Holm-adjusted LCB floors, so v3 makes no such claim.

## Historical v3 arm contract

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

## Historical v3 analysis

The assigned case-arm-repetition trial is the unit of analysis. Comparisons pair arms by case and repetition and cluster by generating family and case.

For subsets with at least 20 families, one-sided bounds use a 10,000-replicate paired hierarchical bootstrap: resample families, then cases within families, while retaining arm pairs and repetitions. For fewer than 20 families, the primary method is a sign-flip permutation of unweighted family means of paired case-level differences: all `2^G` flips when `G <= 16`, or 100,000 seeded Monte Carlo flips when `16 < G < 20`. A percentile bootstrap is not the primary bound below 20 families.

The headline is the single Phase 1 primary claim at one-sided alpha 0.05. The three component comparisons use the same family-count-dependent method as isolations, not Holm-adjusted bounds. Phase 2 learning is a separate single claim.

Ratios operate on family totals and are analyzed on the log scale. A zero family total receives a +0.5 token continuity correction. If either compared arm has fewer than ten successful trials, ratio claims are indeterminate. The report also includes complete-case, one-vote-per-case, and one-vote-per-family sensitivity views; disagreement narrows the conclusion or makes it indeterminate.

## Historical v3 forced decisions

| Claim | Rule | Forced verdict |
| --- | --- | --- |
| direct-task harm gate, C versus A | On all assigned direct trials, fail if the one-sided family-mean permutation test that C is worse than A has `p < 0.05`. Otherwise do not block the headline, do not claim noninferiority, and report the difference against -0.10 descriptively. | repair or reject the Phoenix surface |
| headline capability, C versus B | On all assigned trials, pass capability only if the one-sided alpha-0.05 lower bound for C minus B ITT success is greater than zero. ITT tokens, wall time, USD, and paired-success tokens are descriptive and cannot veto this pass. | if this and efficiency-only fail, kill or narrow Phoenix |
| headline efficiency-only | Test only when the primary bound is at most zero. Require success lower bound greater than -0.05, paired-success token-ratio upper bound at most 0.80, C-minus-B dead-end-rate upper bound below 0.05, and C-minus-B timeout-or-cap-hit upper bound below 0.05. | label only as an efficiency product; never claim capability |
| frontier isolation, C versus D | On all assigned cascade, far-discovery, recovery, and temptation trials in the 16-family subset fixed by the allocation, remove the frontier if the family-mean point estimate is at most zero or C is worse at one-sided `p < 0.05`. Tokens are descriptive. | remove the frontier; retain the surface only if the headline passed |
| teaching-refusal isolation, D versus E | On all assigned recovery trials, define recovery from the first failed or declined state-sensitive act through a later successful alternative and successful grade. Replace teaching if D minus E recovery is at most zero, D is worse at one-sided `p < 0.05`, or the downstream ITT-success lower bound is below -0.05. | use plain typed errors |
| learning, counted C versus D' | On all held_out cases, keep counting disabled if the one-sided alpha-0.05 success lower bound is at most zero. Efficiency is descriptive. | keep authored ranking and counting disabled |

The approximate 80% MDEs against zero are 0.19 for the headline and learning claims, 0.26 for the frontier isolation at its allocated 16 families and 60 cases, and 0.39 for both the direct harm gate and the teaching isolation at their effective sample sizes of 20. Direct-gate power is about 0.16 for a 10-point tax and 0.35 for a 20-point tax; the evaluator-specified 0.30 effect size is not an 80% MDE. These values describe sensitivity, not pass floors.

## Historical v3 Phase 2

Phase 2 disables online updates and freezes the daemon, verbs, authored rules, candidate weights, and training-episode digest. Its arms are counted C, D', and B. Counted C versus D' is the learning claim. Counted C versus B repeats the Phase 1 headline success threshold without consuming Phase 1 alpha. If Phase 1 passed but the held_out replication fails, held_out is indeterminate and cannot pass.

## Historical v3 engineering gates

The inferential claims do not override these hard checks:

- incremental standing input cost versus empty MCP is at most 600 tokens at 10, 100, and 1,000 verbs;
- standing-token spread across those world sizes is at most five tokens;
- on Linux, daemon-only act latency has median at most 2 ms and p95 at most 10 ms;
- daemon time is at most 5% of median successful trial wall time;
- malformed-call rate is at most 5% across at least 20 fresh one-call sessions, with the Wilson 95% interval reported.

Failure blocks the surface claim. Repairing a failed gate after seeing validation outcomes burns the tranche.

## Historical v3 budget

The hard budget remains USD 565:

| Work | Arithmetic | Ceiling |
| --- | --- | ---: |
| authoring | stated allocation | USD 75 |
| validation | 120 cases x 3 repetitions x 5 arms x USD 0.15 x 1.10 = USD 297 | USD 300 |
| held_out | 120 cases x 3 repetitions x 3 arms x USD 0.15 x 1.10 = USD 178.20 | USD 180 |
| surface and protocol | stated allocation | USD 10 |

The 10% allowance is pooled infrastructure capacity, not authorization for every trial to retry. A budget or safety stop is indeterminate, never a pass.

## Historical v3 freeze boundary

The v3 contract required separate digests for the runtime invocation, system prompt, A schemas, B document, world build, runner, schedule, grader, analysis implementation, and report template before validation could open. Those historical digests cannot freeze v4 artifacts.

Review 1 identified the v2 protocol by LF-normalized UTF-8 SHA-256 and returned `REVISE`; it did not issue an acceptance record. The independent confirmation review accepted and froze only the corrected v3 candidate. An independent reviewer accepted Task 0.2 on 2026-08-18. Those historical acceptances do not apply to v4.

## Historical independent v3 acceptance record

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

This record applied only to the superseded v3 candidate. It is retained for custody and must not be read as a current freeze, current acceptance, permission to create sealed families, or authorization to run validation or held_out outcomes.
