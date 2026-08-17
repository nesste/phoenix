# Phoenix

Phoenix is a proposed local-first **environment daemon** — a "world" — that frozen-weight AI agents inhabit.

Instead of describing capabilities to an agent upfront (tool schemas, skill files, CLAUDE.md walls), the world delivers capability knowledge in-band:

- the agent wakes up holding a goal and a few live **handles**; standing context is O(1) in the number of tools;
- every action's **result carries a frontier** — at most three state-computed suggestions of what that action just made relevant, arguments already bound;
- **refusals teach**: a declined action names what, why, and the reachable alternative, at the only moment that knowledge was ever going to be read;
- constraints are **topology, not policy**: what must be impossible is absent, not denied;
- the **world learns, not the agent**: frontier rankings come from counted outcomes, friction becomes candidate verbs, and no change lands without replaying a case bank.

Phoenix is **not** an agent harness, model router, prompt marketplace, permission layer, chat interface, or autonomous self-modifier. It never calls a model.

## Status

Planning only. No implementation has been accepted or started.

The core bet — that in-band frontier delivery beats upfront documentation for frozen-weight agents — is treated as a falsifiable research question with pre-committed thresholds, held-out task corpora, and component-level kill conditions.

## Plan

See [`docs/plans/2026-08-17-phoenix-world-plan.md`](docs/plans/2026-08-17-phoenix-world-plan.md):

1. Phase 0 — surface spike, frozen result envelope, task corpus, pre-committed experiment;
2. Phase 1 — the static world: handles, verbs, authored frontiers, teaching refusals, episode log;
3. Phase 2 — the world that learns: counters, case bank, replay-gated typed amendments;
4. Phase 3 — the shared world: one world improving across many sessions.

The earlier control-plane plan is retained as [`docs/plans/2026-08-17-phoenix-implementation-plan.md`](docs/plans/2026-08-17-phoenix-implementation-plan.md) (superseded).

Phoenix follows the currently prioritized research-kernel milestone unless Raoul explicitly changes that sequencing.
