# Phase 1 arm-adapter implementation status

Date: 2026-08-18

Protocol v4 is accepted and frozen. This status covers the arm adapters and scheduled authoring machinery. It does not open Gate 1A or either sealed tranche.

## Implemented

- Arms A and B expose the complete conventional tool set generated from the same world definition used by C, D, and E.
- Every arm uses the same graph, verb executor, episode store, and executable-indexed state-event hook. An integration test verifies that a flat-tool result is computed before the event and delivered after the event is applied.
- D suppresses post-act frontiers before pending state is recorded. E suppresses both frontiers and teaching alternatives before pending state is recorded.
- The server and authoring runner accept only A, B, C, D, or E. Exact protocol system prompts are pinned in code; the Phoenix bootstrap-intent sentence is present only for C, D, and E.
- Flat arms receive no opaque-handle prompt. Their runtime allowlist is derived deterministically from the world definition.
- The Arm B authoring-only guide lists every starter verb and covers one visible acceptable path per authoring class. Mechanical tests keep the guide aligned with the world and visible labels.
- An independent human-factors reviewer accepted the Arm B guide at commit `73adf8c608f0edf06597b569b17faa32e1a3b5b9`. `pre-validation-artifacts.json` freezes its accepted digest and keeps both outcome gates false.
- Non-C authoring probes default to arm-specific output directories, protecting the retained Arm C evidence.
- The schedule writer implements the protocol seed, three repetitions, family blocks, and the ten-row odd-treatment Williams design. Its loader reconstructs the expected schedule and rejects order or design drift.
- Scheduled execution keeps every five-arm pairing key contiguous and reserves its five trial caps plus the pooled 10% infrastructure capacity before launch.
- Every retry gets a fresh context, sandbox, world, handles, episode store, and state-event plan. Only one eligible pre-token infrastructure retry is permitted.
- Assignment records distinguish graded pass/fail, manual-required failure, terminal cap failure, exhausted infrastructure, budget stop, and safety stop. They retain every attempt's token buckets, USD, turns, API time, wall time, and first-token status.
- The scheduled runner at `b4df919070bb9a6d2912662b4a59674b0e25a332` received an independent `ACCEPT` with no findings. Its nine reviewed files are frozen as one file-digest set in `pre-validation-artifacts.json`.
- The Phase 1 analysis implements ITT accounting, whole-pair safety-stop handling, unresolved and cap-imbalance rules, paired hierarchical bootstrap and family sign flips, paired-success token ratios, the three component contrasts, sensitivity views, and a deterministic Markdown report template. It reconstructs frontier linkage from retained runtime evidence.
- The focused analysis revision at `62946f4a1a03ea89636c5b3243f3b4d53b166682` received an independent `ACCEPT` with no findings. The nine analysis and report artifacts are frozen by digest. The accepted non-blocking limitation is that descriptive measures do not yet carry two-sided 95% intervals; those intervals enter no forced decision.
- The outcome-free candidate for the other four locally available artifacts received an independent focused-revision `ACCEPT` at `795ce71acd51beb977190811c90cc538f3c6b928`. The runtime contract and exact prompts, twelve Arm A tool schemas, world definition/build, and grader are copied into `pre-validation-artifacts.json` with their reviewed identities. Mechanical tests bind every accepted object back to the candidate. Both outcome gates remain closed.

## Still required before a Phase 1 outcome run

1. Independently generate and seal new disjoint validation and held-out candidates. Do not inspect or execute their outcomes while completing the remaining artifacts.
2. Generate the validation schedule only from the accepted replacement validation cases and freeze its digest. The authoring schedule is not a substitute.
3. Keep the manifest `partial` and both outcome gates false until the schedule digest and every external blocker are closed.

The current runner remains authoring-only. No authoring execution is required by this implementation tranche, and the retained 7/8 focused-revision run is unchanged.
