# Phoenix

Phoenix is a proposed standalone, local-first capability control plane for AI agents.

It is designed to make behavioral contracts:

- activate from observable opportunities rather than optional prompt retrieval;
- disclose whether enforcement is advisory, cooperative, or mediated;
- require evidence bound to exact actions and artifacts;
- make missed activation visible instead of silently doing nothing;
- evolve only through replayed, reviewed, and reversible amendments.

Phoenix is **not** an agent harness, model router, workflow engine, chat interface, or autonomous self-modification system.

## Status

Planning only. No implementation has been accepted or started.

The implementation sequence deliberately treats activation recall and the activation corpus as the primary research risks. The initial experiment is limited to one contract, one external runtime, one mediated filesystem publication gateway, and one evidence ledger.

## Plan

See [`docs/plans/2026-08-17-phoenix-implementation-plan.md`](docs/plans/2026-08-17-phoenix-implementation-plan.md).

The plan includes:

1. Phase 0 research protocol and held-out activation corpus;
2. thin Phase 1 activation and mediated-publication experiment;
3. Phase 2 second-runtime portability confirmation;
4. Phase 3 evidence-gated amendments;
5. Phase 4 multi-agent capability composition;
6. explicit component-level kill conditions.

Phoenix follows the currently prioritized research-kernel milestone unless Raoul explicitly changes that sequencing.
