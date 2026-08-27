# Protocol-v4 second interrupted validation: exploratory findings

- **Date:** 2026-08-27
- **Status:** `indeterminate`
- **Use:** engineering diagnosis only
- **Decision:** close Gate 1A; do not resume this archive; keep held-out closed

## Bottom line

The repaired execution reached 1,475 of 1,800 checkpointed launches before the PC restart. It is not a valid Gate 1A result because the schedule is incomplete and there is no `scheduled-summary.json`. The prefix is nevertheless large and balanced enough to answer the engineering question: the current Phoenix protocol is not competitive with the conventional B surface.

B succeeds on 122 of 295 trials (41.4%). C succeeds on 56 of 295 (19.0%), a family-weighted difference of -22.5 percentage points. D succeeds on 10.8% and E on 8.1%. Restricting the extraction to the first 19 complete family blocks barely changes C-B: -22.1 points.

The earlier world-reference and orientation repairs did their job. Every graded C, D, and E trial now contains a matched orientation, and all retained trials use the frozen world build. The remaining failure is in the protocol behavior after orientation: C reaches the turn limit in 73 of 295 assignments, costs roughly four times as much as B, and still loses badly on direct, adversarial-text, far-discovery, stale-frontier, and temptation cases.

## Evidence and validity

The archive contains 5,649 files with aggregate index digest `sha256:3dff34f5c5e6fbd1958497cdca961b5a04fda1f7c84fd8215bb5f49bfb710241`.

- Durable checkpoint: 1,475 launches, indexes 0-1474
- Complete five-arm pairing keys: 295
- Coverage: 99 cases from 20 families; 19 complete family blocks
- Checkpointed evidence: 1,347 graded trials and 128 runtime-terminal trials
- Checkpointed cost: 60.4877118 USD
- World build: frozen and observed `sha256:b5a26d5e…a2c4`
- Scheduled summary: absent

The archive also contains completed launch 1475. It was written after the last durable checkpoint and belongs to an incomplete pairing key, so this report excludes it. Including it would violate the precommitted five-arm pairing contract.

The execution is `indeterminate`, not a shortened confirmatory result. Its outcomes have now been inspected, so the validation tranche is burned and cannot be resumed or reused for another Gate 1A claim.

## Arm results

| Arm | Success | Rate | Terminal failures | Cap hits | Tokens | Cost (USD) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| A | 103/295 | 34.9% | 7 | 7 | 3,364,441 | 5.8132 |
| B | 122/295 | 41.4% | 2 | 2 | 3,883,734 | 4.7110 |
| C | 56/295 | 19.0% | 74 | 73 | 9,111,034 | 18.9993 |
| D | 32/295 | 10.8% | 25 | 25 | 6,125,892 | 15.2672 |
| E | 24/295 | 8.1% | 20 | 20 | 6,290,836 | 15.6970 |

C costs 4.03 times as much as B for the same number of assignments. Its observed cost per success is about 0.339 USD, versus 0.039 USD for B, an 8.79-fold difference. These are diagnostic ratios, not confirmatory efficiency estimates.

Across the 295 paired comparisons, C beats B 11 times, B beats C 77 times, and 207 pairs tie, mostly because both fail. The raw C-B success difference is -22.4 points and the family-weighted difference is -22.5 points.

## Class-level success

| Class | Pairings | A | B | C | D | E |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| absence | 27 | 0.0% | 0.0% | 0.0% | 0.0% | 0.0% |
| adversarial text | 24 | 50.0% | 45.8% | 0.0% | 0.0% | 0.0% |
| cascade | 30 | 0.0% | 0.0% | 0.0% | 0.0% | 0.0% |
| direct | 63 | 84.1% | 79.4% | 49.2% | 36.5% | 38.1% |
| far discovery | 30 | 70.0% | 53.3% | 6.7% | 0.0% | 0.0% |
| recovery | 59 | 5.1% | 5.1% | 13.6% | 13.6% | 0.0% |
| stale frontier | 33 | 6.1% | 87.9% | 45.5% | 0.0% | 0.0% |
| temptation | 29 | 41.4% | 44.8% | 0.0% | 3.4% | 0.0% |

Phoenix shows a small recovery-class benefit over A/B, but it is far too narrow to offset the losses elsewhere. The zero success of every arm on absence and cascade is a separate warning: those classes need an authoring-side audit of case solvability, world reachability, grader expectations, and exposed surface behavior.

## Failure mechanism

The checkpointed prefix contains 127 turn-limit failures and one agent error. C accounts for 73 turn-limit failures and the single agent error. Those 74 terminal trials have no grade because the frozen runner correctly treats post-token turn-limit and agent failures as ITT failures.

This differs from the first 720-run diagnostic archive. In that run, most Phoenix trials failed before reaching an executable orientation. Here, all 221 graded C trials, all 270 graded D trials, and all 275 graded E trials contain a matched orientation. C improves from 14.6% to 19.0% and its deficit to B narrows by roughly five points, but B remains stable near 41%. The orientation repair removed one bottleneck without making the main protocol effective.

D and E do not rescue the design. They hit the turn limit less often than C, yet their success rates fall to 10.8% and 8.1%. E records no recovery-class successes. The teaching/refusal layers therefore add cost without producing a useful task-level benefit in this sample.

## Required next work

Do not resume this run and do not run another validation tranche yet.

1. Use authoring-only cases to reduce C's repeated-act and turn-limit behavior while preserving the now-working orientation path.
2. Diagnose the direct, adversarial-text, far-discovery, stale-frontier, and temptation regressions against B. C must at least stop imposing a large capability tax before another gate attempt.
3. Audit absence and cascade end to end because every surface scores zero. Check case solvability, reachable verbs, grader paths, and whether the intended success route can be executed within 12 turns.
4. Reconsider D and E rather than tuning them around the edges. Their current teaching/refusal behavior does not recover enough success to justify its cost.
5. Independently review and refreeze any protocol, prompt, world, runner, or grader changes.
6. Generate a new disjoint sealed validation tranche. This tranche and the earlier 720-run tranche are both diagnostic archives and cannot support a future Gate 1A claim.

The evidence is enough to stop and improve the protocol. It is not a valid basis for opening held-out or continuing to Phase 2.
