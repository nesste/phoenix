# Protocol-v4 interrupted validation: exploratory extraction

- **Date:** 2026-08-23
- **Status:** `indeterminate`
- **Use:** engineering diagnosis only
- **Decision:** do not resume this run; do not open held-out or Phase 2

## Bottom line

The retained results are enough to reject the current execution and diagnose the main failure. They are not enough to support a Gate 1A capability claim or continue to held-out.

Arm C underperformed Arm B by roughly 26–27 percentage points in both partial-data views. More importantly, Phoenix produced no executable act in two thirds of C trials. C worked well when orientation matched: 19 of 24 matched trials passed. It almost never worked when orientation did not match: 2 of 120 passed. The immediate engineering problem is the orientation-to-executable-path handoff.

The validation tranche is now burned. The run stopped early, used the wrong world-build identity, and its outcomes have been inspected. A valid next validation requires repaired infrastructure and a new disjoint sealed tranche.

## Evidence and validity

The retained directory contains 2,788 files with aggregate index digest `sha256:4ff9b0d3996a14606c90f1612ab1823296aafe1278ecf18da0ae2dbf7572c114`.

- 720 assignments, exactly schedule indexes 0–719
- 144 complete five-arm pairing keys
- 674 graded trials and 46 runtime-terminal trials
- 49 cases from 10 families; nine families are complete
- 144 assignments per arm
- 28.4686273 USD recorded cost

The data cannot produce a confirmatory result for two independent reasons:

1. Execution stopped after 720 of 1,800 assignments and produced no `scheduled-summary.json`.
2. Every retained trial pins world build `sha256:bbcf6e83…`, while Gate 1A froze `sha256:27c2f537…`. The runner should have rejected this before the first model call.

The extraction used the frozen analysis formulas directly. It analyzed all 720 assignments and repeated the calculation on the first 675 assignments, which cover nine complete families. Every claim remains `indeterminate`; the numbers below are exploratory.

## Arm results

| Arm | Success | Rate | Cap hits | Zero-act trials | Tokens | Cost (USD) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| A | 46/144 | 31.9% | 1 | 8 | 1,650,922 | 2.8020 |
| B | 59/144 | 41.0% | 0 | 6 | 1,984,059 | 2.3549 |
| C | 21/144 | 14.6% | 7 | 96 | 2,584,373 | 7.2459 |
| D | 13/144 | 9.0% | 16 | 104 | 2,822,498 | 7.6260 |
| E | 1/144 | 0.7% | 22 | 105 | 3,330,133 | 8.4398 |

Phoenix was both less successful and more expensive in this interrupted sample. C cost about three times as much as B while succeeding about one third as often.

## Headline comparison

For all 720 assignments, the family-weighted C−B success difference is −27.1 percentage points. The exact family sign-flip interval runs from −46.2 to −8.0 points; exploratory `p=0.0078125`.

Removing the partially observed tenth family barely changes the result. Across nine complete families, C−B is −25.9 points, with bounds from −44.4 to −7.4; exploratory `p=0.015625`.

The registered sensitivity views agree with the direction and magnitude. This is not a close or unstable partial result.

## Class-level success

| Class | A | B | C | D | E |
| --- | ---: | ---: | ---: | ---: | ---: |
| absence | 0.0% | 0.0% | 0.0% | 0.0% | 0.0% |
| adversarial text | 64.3% | 71.4% | 0.0% | 0.0% | 0.0% |
| cascade | 0.0% | 0.0% | 0.0% | 0.0% | 0.0% |
| direct | 100.0% | 88.9% | 11.1% | 11.1% | 5.6% |
| far discovery | 50.0% | 33.3% | 0.0% | 0.0% | 0.0% |
| recovery | 6.7% | 6.7% | 17.8% | 20.0% | 0.0% |
| stale frontier | 8.3% | 100.0% | 83.3% | 0.0% | 0.0% |
| temptation | 30.0% | 60.0% | 5.0% | 10.0% | 0.0% |

C shows useful behavior on stale-frontier and some recovery cases, but the gains do not offset broad failures on direct, adversarial-text, far-discovery, and temptation cases.

## Failure mechanism

| Arm | Trials with matched orientation | Success when matched | Success when unmatched | Zero-act trials |
| --- | ---: | ---: | ---: | ---: |
| C | 24/144 | 19/24 (79.2%) | 2/120 (1.7%) | 96/144 (66.7%) |
| D | 22/144 | 11/22 (50.0%) | 2/122 (1.6%) | 104/144 (72.2%) |
| E | 17/144 | 0/17 (0.0%) | 1/127 (0.8%) | 105/144 (72.9%) |

C recorded 619 orientation attempts but only 29 matched events across 24 trials. Unmatched orientation usually left the model without an executable path; it repeated orientation or stopped without acting. The full Phoenix surface therefore failed mostly before tool execution.

The component comparisons point in the same direction:

- Direct C−A: −88.9 points across the two observed direct families.
- Frontier C−D: −1.7 points. The frontier did not improve observed task success.
- Teaching recovery D−E: +20 points on recovery cases, but D's overall success was only 9.0% and this isolation cannot rescue the surface.

## Required redo

Do not resume this run. Preserve it as diagnostic evidence.

Before another Gate 1A attempt:

1. Make the runner reproduce and enforce the exact frozen world-build identity before output creation or model contact.
2. Fix and test the orientation matcher and the no-match handoff using authoring-only cases. A no-match response must not leave the model in a repeated-intent dead end.
3. Add a tested checkpoint/resume mechanism or run in infrastructure that can finish within the provider window.
4. Independently review and refreeze the changed runner, prompts, matcher, and build path.
5. Generate and seal a new disjoint validation tranche. The current tranche cannot be reused because its outcomes now informed diagnosis.

The existing 720 trials are enough to justify repair. They are not enough—and are not valid enough—to justify continuing the project past Gate 1A without a clean rerun.
