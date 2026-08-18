# Phase 1 arm-adapter implementation status

Date: 2026-08-18

Protocol v4 is accepted and frozen. This implementation tranche adds the first post-acceptance runner machinery; it does not open Gate 1A or either sealed tranche.

## Completed in this tranche

- Arms A and B expose the complete conventional tool set generated from the same world definition used by C, D, and E.
- Every arm uses the same graph, verb executor, episode store, and executable-indexed state-event hook. An integration test verifies that a flat-tool result is computed before the event and delivered after the event is applied.
- D suppresses post-act frontiers before pending state is recorded. E suppresses both frontiers and teaching alternatives before pending state is recorded.
- The server and authoring runner accept only A, B, C, D, or E. Exact protocol system prompts are pinned in code; the Phoenix bootstrap-intent sentence is present only for C, D, and E.
- Flat arms receive no opaque-handle prompt. Their runtime allowlist is derived deterministically from the world definition.
- The Arm B authoring-only guide lists every starter verb and covers one visible acceptable path per authoring class. Mechanical tests keep the guide aligned with the world and visible labels.
- Non-C authoring probes default to arm-specific output directories, protecting the retained Arm C evidence.

## Still required before a Phase 1 outcome run

1. Obtain and commit the independent Arm B human-factors review, then freeze the accepted document by digest.
2. Implement the family-blocked Williams schedule, three repetitions per case-arm pair, and deterministic launch order.
3. Implement attempt records, eligible infrastructure retry classification, pairing-key budget stops, timeout/cap ITT handling, and complete runtime/token/cost metadata.
4. Add the precommitted analysis implementation and report template, with tests for ITT, unresolved and cap-imbalance rules, paired inference, and component contrasts.
5. Freeze every `artifact_freeze.before_validation` item by digest and complete Task 0.6 independent review.
6. Independently generate and seal new disjoint validation and held-out candidates. Do not inspect or execute their outcomes while completing the items above.

The current runner remains authoring-only. No authoring execution is required by this implementation tranche, and the retained 7/8 focused-revision run is unchanged.
