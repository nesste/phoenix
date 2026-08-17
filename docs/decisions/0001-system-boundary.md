# 0001: Phoenix system boundary

- **Status:** Accepted
- **Date:** 2026-08-17
- **Scope:** Phase 0 and every later Phoenix phase

## Decision

Phoenix is a local environment daemon. It owns a live graph of handles, dispatches typed verbs, computes result frontiers, shapes refusals, and records episodes. Phoenix does not call a model, run an agent loop, orchestrate agents, own chat state, or depend on an agent harness.

The agent-facing surface obeys three numeric laws:

1. The complete standing surface, including the tool description and required instructions, is at most 600 tokens and does not grow with the number of verbs.
2. Each result contains at most three frontier entries. Each `why` line is at most 80 characters, and every suggested call has a reachable handle and fully bound arguments.
3. Behavioral rules live in topology, verb behavior, or teaching refusals. Standing prose rules fail the build.

## Enforcement boundary

Phoenix uses absence rather than a global permission policy. A session can invoke only verbs attached to handles reachable from that session. Guessed, revoked, cross-session, and nonexistent handle references produce the same typed `absent` result. Phoenix does not expose a global verb namespace that could reveal unreachable capabilities.

Verbs run under the daemon's operating-system principal. That principal and the host filesystem/process permissions bound what a verb can physically do. Phoenix adds no separate policy engine and makes no claim that prose instructions constrain effects.

The measured claim is limited to agents interacting through the Phoenix surface. An agent with independent shell, filesystem, network, or process access can bypass Phoenix and is outside this boundary. Every experiment arm must record its access mode and use the same access mode as its comparator.

## Non-goals

Phoenix is not:

- an agent framework, runtime adapter marketplace, or multi-agent coordinator;
- a model router, prompt manager, skill registry, or capability catalog loaded at wake-up;
- a general permission or policy layer;
- a chat interface, dashboard, hosted dependency, or autonomous self-modifier.

World code, topology, authored transitions, refusals, and active weights change only through the amendment process defined in later phases. Phase 0 and Phase 1 do not implement autonomous amendments.

## Consequences

- Runtime integrations remain outside the daemon and may be replaced without changing world semantics.
- Security claims must name the daemon principal and the agent's independent access explicitly.
- If an operation must be impossible through Phoenix, the world omits the verb or keeps its handle unreachable.
- A refusal constrains only a state-sensitive verb that exists; it is not evidence of a hidden global capability.
- The constant surface and frontier budgets are CI failures, not documentation preferences.

## Review against the plan

This decision preserves every non-negotiable boundary in section 2 of the world plan. Production deployment, hosted operation, a second runtime, plugins, dashboards, and shared worlds remain outside the accepted Phase 0 scope.
