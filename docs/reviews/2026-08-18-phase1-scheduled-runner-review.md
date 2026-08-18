# Phase 1 scheduled-runner independent review

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-18
- **Verdict:** `ACCEPT`
- **Candidate commit:** `b4df919070bb9a6d2912662b4a59674b0e25a332`
- **Findings:** None

## Candidate identity

The reviewed candidate is `feat: add scheduled Phase 1 runner` at `b4df919070bb9a6d2912662b4a59674b0e25a332`. Independently computed LF-normalized UTF-8 SHA-256 for the files named in the review prompt:

| Artifact | SHA-256 |
| --- | --- |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| `docs/decisions/0004-go-no-go-rules.md` | `sha256:883678b0799c76509ac306e7c38f322c35dab12d8dfe19b71228045b22652e9e` |
| `experiments/frontier-v1/runner/schedule.go` | `sha256:d04884d84036d3d8dd8cdd6f1fccdc4f889f3289831bb9ecec0e84292e894c95` |
| `experiments/frontier-v1/runner/schedule_test.go` | `sha256:73c9bb9238ab26b1c64cd49e5989ee2e16d92f44b3203a465bced09f962729ec` |
| `experiments/frontier-v1/runner/scheduled_run.go` | `sha256:4a86b8d84789fa60c9b1991dcfaf0c9b82e9cee7d37459af79df84d69d0ba1cf` |
| `experiments/frontier-v1/runner/scheduled_run_test.go` | `sha256:080ce37de7acd2c7d741a10cb3b1658bd77a39809b546d606f43c78af960be83` |
| `experiments/frontier-v1/runner/runtime.go` | `sha256:fe7d608fb2b64defb2ae5097235b93191d7e386e75693a9b8cfd6103940aeb56` |
| `experiments/frontier-v1/runner/runtime_test.go` | `sha256:899f1ac588082420f652b1d6fdd7eb0ca20065dc2b11eabd50839eeac93183b2` |
| `experiments/frontier-v1/runner/main.go` | `sha256:e0949c0abf26fa6fe9f3b2ccd8c8c8c937d69098110a8454c08fad8b96c06c35` |
| `experiments/frontier-v1/runner/main_test.go` | `sha256:a8ad317b9e10f541d5ace06ea0ac9422ccf8b04077bb573b9e6bcfa4e8cf9847` |
| `experiments/frontier-v1/runner/types.go` | `sha256:54da663664141753fab65d89ed8db98101bb01613237c8a5dd75df78f908fc94` |
| `experiments/frontier-v1/README.md` | `sha256:0a9770860cadc0ad12d1e2463c3b621cf54ca5efc8f8e70a6abac27449d54b72` |

The candidate commit does not modify `protocol.json`, `pre-validation-artifacts.json`, `experiments/frontier-v1/arms/arm-b.md`, authoring labels, or retained authoring evidence. Rechecked freeze identities:

| Artifact | SHA-256 |
| --- | --- |
| Arm B document | `sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0` |
| Authoring `summary.json` | `sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405` |

Those match the accepted Arm B freeze and the protocol-v4 authoring-summary pin. `pre-validation-artifacts.json` remains `partial` with both outcome gates false.

The runner package tests passed. Outcome-free schedule generation with seed `20260817` against the eight authoring cases produced 120 launches and digest `sha256:165449201a6146c7ba90bb4c2ed3fbe39528974cf16ee81283e9110a27188684`. That digest is a reviewer check, not a committed freeze.

## Required checks

1. **Seed, repetitions, arms.** `phase1ScheduleSeed` is `20260817`, repetitions are 3, and arms are A–E. Regeneration is deterministic. The authoring set yields 8 cases × 3 repetitions × 5 arms = 120 launches in four family blocks.
2. **Williams construction.** For five treatments the writer emits ten rows: five cyclic Williams sequences plus their reversals. Across those rows every directed first-order predecessor `X→Y` (`X≠Y`) appears exactly twice. Every generated pairing uses one of those rows.
3. **Contiguous pairing keys and family blocks.** Each `(case_id, repetition)` occupies five consecutive launches covering A–E once. A family does not reappear after its block closes. The eight cases each appear in repetitions 0, 1, and 2.
4. **Loader rejection.** `validateSchedule` requires version 1, tranche `authoring`, the frozen seed/repetitions/arms, and byte-level equality with regeneration from the selected authoring cases. Seed, arm-order, family-order, pairing-order, and launch-order drift all fail that comparison.
5. **Budget boundary.** Before each pairing key the runner reserves `5 × max_cost_usd_per_trial × 1.10`. If that reserve does not fit, none of the five arms launch, remaining assignments are `budget_stopped`, and the run status is `indeterminate`. A budget stop cannot create a run-level pass.
6. **Retry eligibility.** One retry is granted only for pre-token provider 429/5xx (and equivalent provider transient text), OS-level process launch failure, or MCP connection failure. Timeout, turn limit, cost cap, malformed output, agent error, and post-token failures are terminal.
7. **Retry isolation.** A retry calls `executeCaseAttempt` again: new sandbox, state directory, episode store, handle set, and state-event plan, same pairing key. Tests assert distinct sandboxes and roots, with attempt workspaces removed.
8. **Assignment encoding.** Records keep every attempt. Graded pass/fail, `manual_required`, cap hits, `infrastructure_unresolved`, `budget_stop`, and `safety_stop` are distinct terminations. Unresolved and manual-required outcomes have `itt_success` false.
9. **Token accounting.** Attempt metrics take Claude result `usage` buckets `input_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens`, and `output_tokens`, and total them without adding `modelUsage` on top when `usage` is present. USD, turns, API time, and wall time are stored per attempt. A retained 2.1.229 result event uses the same disjoint buckets (`6+3217+2530+389=6142`) and matching `modelUsage`.
10. **Summary pins.** `scheduled-summary.json` records schedule digest, world-build digest, live grader digest, Arm B path and digest, exact A–E system prompts, per-arm tool allowlists, runtime/model settings, caps, retry limit, and run budget.
11. **Authoring-only.** The loader accepts only `authoring_` cases from `corpus/authoring`. The writer hard-codes tranche `authoring`. Scheduled output must be empty. There is no path that sets `may_open_validation` or `may_open_held_out`, and no path that loads sealed cases.
12. **Unchanged freeze set.** Protocol, Arm B document, pre-validation manifest, and retained authoring evidence are untouched by this commit.

## Decision

No P0, P1, or P2 defect in the scheduled authoring runner. The twelve required checks hold.

This acceptance is limited to the authoring-only scheduled runner at the candidate commit. It does not authorize a model outcome run, freeze a validation schedule, open Gate 1A, or convert the retained 7/8 authoring probe into a gate result.

The runner implementation is ready to freeze as the runner digest after the remaining analysis implementation and report template are committed, and only if the runner bytes are unchanged at that freeze. An authoring schedule digest is not a substitute for the validation schedule freeze. `pre-validation-artifacts.json` must stay `partial` until every other `artifact_freeze.before_validation` item is committed by digest.

## Accepted limitations

- Conclusions remain limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes.
- Headline power and efficiency-only conservatism are unchanged from the protocol-v4 record.
- Cascade remains run-to-run noisy. Its restored four-step label must not be relaxed, and the historical 8/8 must not be selected as a success estimate.
- Cap-hit classification reads result `subtype` (`max_turn`, `budget`/`cost`) plus the harness timeout path. This review did not observe a live cap-hit stream.
- First-model-token detection is the stream-json `assistant` event. A later analysis must still discard any incomplete pairing key left by a safety stop.
- This writer cannot emit a validation or held_out schedule.

## Remaining execution blockers

- independent generation and sealing of new validation and held_out families;
- analysis implementation and report template;
- every remaining `artifact_freeze.before_validation` digest, including runner digest, validation schedule digest, and grader digest freeze;
- Task 0.6 Phase 0 review.

Gate 1A, validation, and held_out remain closed. This review does not authorize opening them.
