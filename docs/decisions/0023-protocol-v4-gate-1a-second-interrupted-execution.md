# 0023: Protocol-v4 Gate 1A second interrupted execution

- **Status:** Authorized validation execution spent; result `indeterminate`; held-out closed
- **Date:** 2026-08-27
- **Decision:** `CLOSE_INTERRUPTED_EXECUTION`
- **Opening decision:** `docs/decisions/0022-protocol-v4-gate-1a-validation-reopening-after-world-compatibility-repair.md`
- **Evidence:** `experiments/frontier-v1/results/scheduled-validation-2`
- **Custody index:** `experiments/frontier-v1/results/scheduled-validation-2-custody.json`
- **Custody aggregate:** `sha256:3dff34f5c5e6fbd1958497cdca961b5a04fda1f7c84fd8215bb5f49bfb710241`
- **Exploratory report:** `docs/reports/2026-08-27-protocol-v4-second-interrupted-validation-exploratory.md`

## Decision

The project chair records that the single Gate 1A execution authorized by decision 0022 is spent and `indeterminate`. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: false`, keeps `may_open_held_out: false`, and names this close-out and custody record.

The host restart terminated the external launcher, Go runner, Claude process, and custodian process. A process audit on 2026-08-27 found none of them running. No process required a forced stop.

The durable checkpoint covers 1,475 of 1,800 launches, indexes 0-1474, with 295 complete A-E pairing keys and 60.4877118 USD recorded spend. It contains 1,347 graded trials and 128 runtime-terminal ITT failures. The archive also contains completed launch 1475, but that launch follows the last checkpoint and belongs to incomplete pairing key 295. It is excluded from every reported metric.

The run produced no `scheduled-summary.json`. The full frozen schedule was not completed, so the result is indeterminate regardless of the partial effect estimates. Outcomes have been inspected for engineering diagnosis, which burns this validation tranche for future confirmatory use.

This decision does not resume, splice, extend, or complete the archive; open held-out; change a frozen runner, prompt, world, grader, schedule, or analysis byte; authorize another validation execution; or start Phase 2.

## Diagnostic finding

The world-compatibility repair is effective: retained trials match the frozen world build, and every graded Phoenix trial contains a matched orientation. The current protocol still performs poorly. In the checkpointed prefix, B succeeds on 41.4% of assignments, C on 19.0%, D on 10.8%, and E on 8.1%. The family-weighted C-B difference is -22.5 percentage points and remains -22.1 points across the first 19 complete family blocks. C incurs 73 turn-limit failures and costs about four times as much as B for equal assignments.

These are exploratory diagnostics, not Gate 1A claims. The committed report records the full arm, class, termination, cost, and sensitivity extraction.

## Why the execution is indeterminate

Decision 0022 authorized one complete 1,800-launch execution. A host restart stopped it after 295 of 360 pairing keys. Pairing-key checkpointing preserved a coherent prefix and prevented the loss seen in decision 0013, but it cannot turn an incomplete schedule into a gate result. The absent scheduled summary is a second mechanical indication that execution did not complete.

## What this archive may be used for

The directory is a diagnostic archive only. It may inform authoring-side repairs to the Phoenix protocol, turn economy, class reachability, and teaching/refusal design. It may not support a capability claim, an efficiency claim, held-out opening, Phase 2, or a splice into a later run.

## Authorized next action

Keep held-out and validation closed. Do not run more sealed-tranche tests yet.

Authoring-only work may investigate C's turn-limit behavior, its capability tax relative to B, the zero-success absence and cascade classes, and whether D/E should be redesigned or removed. Any frozen-boundary change requires independent review and refreeze. Any later Gate 1A attempt requires a new disjoint sealed validation tranche because outcomes from this tranche and the earlier 720-run tranche have been inspected.
