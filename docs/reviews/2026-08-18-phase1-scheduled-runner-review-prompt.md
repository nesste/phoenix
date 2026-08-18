# Phase 1 scheduled-runner review prompt

Give this prompt to an evaluation reviewer who did not implement the scheduled runner.

```text
Review the frontier-v1 authoring-only scheduled runner at the identified candidate commit.

Read:

- experiments/frontier-v1/protocol.json, especially trial, arm_equivalence, measures, and artifact_freeze
- docs/decisions/0004-go-no-go-rules.md
- experiments/frontier-v1/runner/schedule.go and schedule_test.go
- experiments/frontier-v1/runner/scheduled_run.go and scheduled_run_test.go
- experiments/frontier-v1/runner/runtime.go and runtime_test.go
- experiments/frontier-v1/runner/main.go, main_test.go, and types.go
- experiments/frontier-v1/README.md

Do not run authoring, validation, or held_out model outcomes. Do not inspect private labels. Read-only tests and outcome-free schedule generation in a temporary directory are allowed.

Required checks

1. Seed 20260817 produces a deterministic family-blocked schedule with three repetitions and Arms A-E.
2. The five-arm Williams construction is valid for an odd number of treatments and balances directed first-order predecessors across its ten rows.
3. Every (case_id, repetition) pairing key is contiguous. A family cannot reappear after its block closes.
4. The loader rejects any change to the seed, repetitions, arms, family order, pairing order, or launch order.
5. The budget check occurs before a pairing key and reserves the five trial caps plus the pooled 10% infrastructure capacity. A budget stop launches none of that key and cannot create a pass.
6. Only provider 429/5xx, process launch failure, or MCP connection failure before the first model token receives one retry. Timeout, turn limit, cost cap, malformed output, agent error, tool error, refusal, and post-token failure never retry.
7. A retry receives a fresh context, sandbox, world state, handles, episode store, and state-event plan while retaining its original pairing key.
8. Assignment records preserve every attempt and correctly encode ITT failure, manual-required failure, cap hit, exhausted infrastructure, budget stop, and safety stop.
9. Token accounting captures input, cache-creation input, cache-read input, and output buckets without double counting; USD, turns, API time, and wall time remain attributable to each attempt.
10. The summary records the schedule and world digests, actual grader digest, Arm B digest, exact system prompts and tool allowlists, runtime/model settings, caps, retry limit, and run budget.
11. The runner remains authoring-only, refuses a nonempty scheduled-output directory, and cannot open Gate 1A, validation, or held_out.
12. Existing retained authoring evidence, protocol.json, pre-validation-artifacts.json, and the accepted Arm B document are unchanged.

Return ACCEPT, REVISE, or REJECT. List findings by priority with exact files and lines. If ACCEPT, identify the candidate commit, record the reviewed file digests, state whether the runner is ready to freeze after the remaining analysis/report work, and preserve all closed-gate limitations. This review does not authorize a model outcome run or freeze a validation schedule.
```
