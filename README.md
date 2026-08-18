# Phoenix

Phoenix is a research project for giving AI agents capabilities as they need them, instead of loading every tool and instruction into the prompt at startup.

It runs as a local environment service called a **world**. An agent starts with a goal and a few references to real things, such as a repository or test suite. The agent performs typed actions through one small interface. Each result can suggest a few relevant next actions with their arguments already filled in.

Phoenix does not call a model or run an agent loop. It provides the environment an external agent acts inside.

> **Project status:** Phase 0 research contract accepted. The public validation and held-out corpora are sealed, their full labels remain external, and no sealed outcome has been opened. Phase 1 implementation is underway: Tasks 1.1–1.6 provide the command shell, strict world loader, isolated session graph, opaque handles, typed verb execution, the closed dev-repo starter verb set, the constant `act` MCP surface, authored frontier computation, teaching refusals, quality gate, and content-addressed build manifest.

## The problem Phoenix is testing

An agent with many conventional tools often receives every tool name, description, and input schema before it starts. That standing context grows with the number of capabilities, even though most of them are irrelevant to the current task.

Phoenix tests a different approach:

1. Give the agent one constant `act(handle, verb, args)` interface.
2. Expose only the handles currently reachable in the world.
3. Return a small **frontier** of useful next calls after each action.
4. Explain recoverable mistakes at the moment they happen through typed refusals.

The central research question is whether this in-band guidance improves task success and token use compared with complete upfront documentation.

## A small example

Suppose an agent starts with a handle for a repository and needs to check its health:

```text
act(handle: "h_repo_demo", verb: "status", args: {})

result: repo is clean on main
next:   h_tests_demo.run
why:    verify the suite before changing code
```

The agent did not need the complete test API at startup. The repository result made the test suite relevant, so Phoenix returned a ready-to-run suggestion. If an action is missing required input, Phoenix can return a structured refusal with the reason and a reachable alternative. If a handle is not reachable, Phoenix returns the same generic `absent` result whether the reference was guessed, revoked, cross-session, or nonexistent.

## Core concepts

| Term | Meaning |
| --- | --- |
| **World** | The local service and live object graph the agent interacts with. |
| **Handle** | An opaque, session-scoped reference to something real, such as a repository or test suite. |
| **Verb** | A typed operation attached to a handle. |
| **Result** | The typed outcome of a verb call: `ok`, `fail`, `refused`, or `absent`. |
| **Frontier** | Up to three ranked next calls made relevant by the latest result. |
| **Refusal** | A state-sensitive rejection that explains what failed, why, and the reachable alternative. |
| **Episode** | The append-only record of one agent session. |
| **Case** | A frozen, replayable task with an expected outcome. |
| **Amendment** | A proposed change to the world that must pass replay gates before it can become active. |

## The three design rules

Phoenix treats these as build and experiment constraints:

1. **Constant wake-up cost.** The complete standing interface is at most 600 tokens and does not grow with the number of verbs.
2. **Small frontiers.** A result contains at most three suggestions. Each explanation is at most 80 characters, and suggested arguments are already bound.
3. **No standing prose rules.** Behavioral constraints belong in reachability, verb behavior, or typed refusals instead of a large instruction document.

Reachability is the enforcement model. A session can call only verbs attached to handles it can reach. Verbs still run with the daemon process's operating-system permissions; Phoenix is not a replacement for filesystem, process, or network isolation.

## What Phoenix is not

Phoenix is not an agent framework, model router, chat application, prompt marketplace, permission engine, or autonomous self-modifier. It does not orchestrate agents and does not claim to constrain an agent that also has independent shell or filesystem access.

## What exists today

| Area | State |
| --- | --- |
| System boundary | Accepted. Phoenix is a local world service and never calls models. |
| Result and world formats | Version 1 JSON Schemas and examples are checked in under [`spec/`](spec/). |
| Surface spike | Working Go/MCP prototype with structured `act` and restricted `eval` variants. Task 0.2 evidence is independently accepted. |
| Experiment protocol | Frontier experiment protocol v3 is independently accepted and frozen. |
| Corpus tooling | Authoring cases, schemas, deterministic grading, canonical digests, manifest generation, and sealing checks are implemented. |
| Sealed evaluation | 120 validation and 120 held-out cases are committed with outcome-free label digests. Full labels remain external and both outcome gates are closed. |
| Production daemon | Tasks 1.1–1.8 implemented: the assembled dev-repo world serves four opaque roots through one instruction-free `act` tool, with typed verbs, authored frontiers, teaching refusals, live state, and reconstructable WAL-backed episodes. The pinned authoring run passed 2/8 cases; the retained failures expose first-action discovery weakness before a frontier is reached, so Gate 1A remains closed. |

## Build the Phase 1 daemon shell

The production module pins Go 1.26.6 and builds a Linux-amd64 binary with CGO disabled. Run the complete engineering gate from the repository root:

```powershell
make quality
```

The gate tests and analyzes the production module, validates the accepted schemas, builds `bin/phoenix`, and writes the ignored content-addressed manifest at `build/manifest.json`. For a native development smoke check:

```powershell
go run ./cmd/phoenix version
go run ./cmd/phoenix serve --stdio
```

Run the visible tuning tranche with the pinned runtime (validation and held-out inputs are rejected by this runner):

```powershell
go run ./experiments/frontier-v1/runner --repo-root . --case all
```

The second command waits for an MCP client on standard input and output.

The accepted surface evidence recorded zero malformed calls in 20 fresh sessions for each candidate. On Ubuntu 24.04.4 under WSL2, the structured `act` spike measured 636.5 microseconds median daemon-only latency and 837 microseconds p95 across 100 calls.

## Run the Phase 0 spike

The spike uses Go 1.26.5 and the official MCP Go SDK over stdio.

From `experiments/surface-spike`:

```powershell
$env:GOTOOLCHAIN = 'go1.26.5'
go test ./...
go run ./cmd/measure-schema
go run ./cmd/validate-spec --repo-root ../..
```

To build and exercise the structured `act` surface on Windows:

```powershell
go build -o ./bin/phoenix-spike.exe ./cmd/phoenix-spike
go run ./cmd/probe-stdio --exe ./bin/phoenix-spike.exe --surface act
```

See the [surface-spike README](experiments/surface-spike/README.md) for the pinned Claude Code invocation and the corresponding `eval` probe.

## Repository guide

```text
docs/decisions/                 accepted and provisional design decisions
docs/plans/                     phased implementation plans
docs/reviews/                   independent-review records and prompts
experiments/surface-spike/      Phase 0 MCP surface prototype and measurements
experiments/frontier-v1/        frozen experiment protocol, authoring corpus, and tools
spec/                           result and world JSON Schemas with examples
```

Useful starting points:

- [System boundary](docs/decisions/0001-system-boundary.md)
- [Surface and stack decision](docs/decisions/0002-surface-and-stack.md)
- [Result envelope and world format](docs/decisions/0003-result-envelope.md)
- [Frozen experiment rules](docs/decisions/0004-go-no-go-rules.md)
- [Accepted Phase 0 research contract](docs/decisions/0005-phase-0-review.md)
- [Phoenix world plan](docs/plans/2026-08-17-phoenix-world-plan.md)
- [Frozen frontier v1 protocol](experiments/frontier-v1/protocol.json)

## Planned phases

- **Phase 0:** freeze the boundary, interface, schemas, corpus, and experiment before measuring outcomes.
- **Phase 1:** build the static world with handles, verbs, authored frontiers, teaching refusals, and episode records.
- **Phase 2:** add replay-gated learning for frontier weights and typed amendments.
- **Phase 3:** explore a shared world that improves across isolated sessions.

Phase 0 is accepted. Validation and held_out remain closed until the runtime, arm artifacts, world build, runner, schedule, grader, analysis implementation, and report template are frozen by digest.
