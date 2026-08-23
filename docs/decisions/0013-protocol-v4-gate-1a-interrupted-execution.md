# 0013: Protocol-v4 Gate 1A interrupted execution

- **Status:** Authorized validation execution spent; result `indeterminate`; held-out closed
- **Date:** 2026-08-23
- **Decision:** `CLOSE_INTERRUPTED_EXECUTION`
- **Opening decision:** `docs/decisions/0012-protocol-v4-gate-1a-validation-opening.md`
- **Evidence:** `experiments/frontier-v1/results/scheduled-validation`
- **Custody index:** `experiments/frontier-v1/results/scheduled-validation-custody.json`
- **Custody aggregate:** `sha256:4ff9b0d3996a14606c90f1612ab1823296aafe1278ecf18da0ae2dbf7572c114`

## Decision

The project chair records that the single Gate 1A validation execution authorized by decision 0012 is finished and `indeterminate`. `experiments/frontier-v1/pre-validation-artifacts.json` now sets `may_open_validation: false`, keeps `may_open_held_out: false`, and records this close-out.

The retained prefix is 720 of 1,800 launches (indexes 0–719), with 144 complete A–E pairs, 674 graded trials, 720 runtime streams, no `scheduled-summary.json`, and empty stderr. Every retained trial names world-build `sha256:bbcf6e83…` instead of the frozen `sha256:27c2f537…`. Outcomes from this prefix have been inspected. The current validation tranche cannot be reused as an unopened confirmatory sample.

This decision does not resume the prefix, merge it with a later run, open held-out, change a frozen runner, prompt, world, grader, schedule, or analysis byte, authorize a replacement 1,800-launch execution, or start Phase 2.

## Why the execution is indeterminate

Decision 0012 authorized one complete frozen-schedule execution. This execution stopped after 720 launches when the host process was lost. The frozen runner cannot resume a nonempty output directory or reconstruct in-memory spend. Independently, the live world-build identity does not match the freeze. Either defect is enough.

## What this archive may be used for

The 720-trial directory is a diagnostic archive only. It may inform authoring repairs: unmatched orientation dead-ends, world-build enforcement, and durable pause/resume. It may not support a capability claim, an efficiency claim, a held-out opening, or a splice into a future confirmatory run.

## Authorized next action

Keep held-out closed. Do not launch validation again under decision 0012.

Authoring-only repair may proceed for:

1. unmatched-orientation handoff, so a no-match does not become a repeated-intent dead end;
2. a load-bearing check that the live world-build digest equals the frozen digest before the first model call;
3. pairing-key checkpoint/resume, or a host that can finish 1,800 launches inside the provider window.

Those repairs change frozen bytes. They require independent review and a new refreeze before any later chair decision can authorize a **new** disjoint validation execution.
