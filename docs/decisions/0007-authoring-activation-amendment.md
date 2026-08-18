# 0007: Add scoped intent activation and deterministic state events

- **Status:** Independently accepted and frozen; authoring-only execution remains open
- **Date:** 2026-08-18
- **Supersedes:** The v3 surface and corpus contract in decision 0004

## Decision

Phoenix keeps one constant `act` tool and adds a second input mode. An executable input supplies `handle`, `verb`, `args`, and optional `state`. An orientation input supplies `handle` and `intent`.

Before the first executable act in a session, orientation returns the session's current pending authored frontier or refusal alternative when one exists; otherwise it may evaluate authored activation rules in the loaded world. After the first executable act, orientation may return only the current pending frontier or refusal alternative. If none is pending, it returns no calls and cannot reactivate from intent. This bootstrap boundary prevents frontier-suppressed Arms D and E from recreating a suppressed post-act frontier through fresh activation.

A live handle anchors orientation to its session. Returned calls may target any handle already reachable in that session, contain fully bound arguments, and use the same three-entry cap as a post-act frontier. Orientation cannot search the global verb registry or reveal an unreachable handle or verb. Executable fields and `intent` are mutually exclusive both at MCP decoding and at the admission boundary.

When a client selects a state-bound call from the current pending frontier but drops the optional `state` field, admission restores the authored precondition before execution. Retired frontiers cannot supply omitted state. A client therefore cannot accidentally turn a stale-safe current call into an unconditional act or revive stale protocol state from an earlier frontier.

The episode log records orientations, their match result, and their returned-call count. Trial evidence keeps orientations separate from executable acts. Graded paths, act counts, and act positions continue to refer only to executable acts.

Cases may also declare deterministic file replacements after a numbered executable act. State events are a shared trial-harness rule for every arm, including flat-tool Arms A and B. Each arm adapter receives the same plan and applies an event after computing the numbered action's result and before returning it to the runtime; orientations do not advance the index. This gives a state-bound frontier call a reproducible opportunity to become stale without relying on timing or an external process race while preserving cross-arm comparability.

The exact Phase 1 system prompt is pinned separately for each arm. A and B receive only the shared tool-use instruction. C, D, and E additionally receive the Phoenix bootstrap-intent instruction; no flat-tool arm receives it.

## Why

The first authoring run passed 2 of 8 cases. Six failures began before a frontier or teaching refusal could help because the runtime had to guess the first verb from an opaque handle. The stale-frontier case had a separate defect: its label required a state change that its input could not express.

Adding verb maps or worked recipes to the system prompt would turn the Phoenix arm into another documentation arm. Weakening the labels would erase the behaviors the corpus was designed to measure. The amended protocol puts initial selection and state change in typed, logged mechanisms instead.

## Consequences

Protocol v4 candidate `f889f13f0c514fa5108a1e392701ebeadc4376f7` was independently accepted on 2026-08-18 and is frozen by the follow-up acceptance-record patch. The v3 validation and held-out candidates pin the old world, schemas, and grader, so they remain retired unopened. They stay in the repository only as custody evidence and must not be opened, run, modified, or resealed in place.

Gate 1A remains closed. Before validation can open, independent evaluators must generate and seal new disjoint validation and held-out candidates, every pre-validation artifact must be frozen by digest, the Phase 1 A-E runner and Arm B document must be complete, and Task 0.6 must be accepted.
