# Phoenix World Plan

> **For implementers:** execute phase-by-phase using TDD and the gates below. Phoenix remains a standalone process. It never calls models, never orchestrates agents, and never depends on Hermes or any agent harness.

**Supersedes:** `2026-08-17-phoenix-implementation-plan.md` (the control-plane plan). What survives from it: the measurement discipline (pre-committed thresholds, held-out corpora, independent labeling, kill conditions), the TDD verification strategy, the append-only log, and the amendment-gating idea. What is discarded: contracts as the core unit, activation detection as the primary risk, and enforcement/permission gateways as the product. What is transformed: the ledger becomes the episode log; receipts become reconstructable episodes; activation becomes the frontier.

**Goal:** Build a local-first **environment daemon** — "the world" — that frozen-weight agents inhabit. The agent wakes up holding a goal and a few live handles, and everything else arrives in-band: every action's result carries a small ranked frontier of what that action just made relevant, refusals teach the reachable alternative, and the world — not the agent — is what learns. Improvement lands in the world's code, verbs, and frontier weights, gated by replay against a case bank.

**The bet being tested:** in-band, state-computed capability delivery beats upfront documentation (tool schemas, skill files, CLAUDE.md) for frozen-weight agents — on task success, on tokens, and on scaling in the number of capabilities.

---

## 1. Product decision

Phoenix is not a control plane, not "dynamic skills," and not an agent framework.

The core unit is the **world**:

> A live object graph of handles with verbs, served to agents through one narrow constant surface, whose results carry their own frontier, whose refusals teach, and whose transition function is ordinary versioned code improved only through replayed, gated amendments.

### Public vocabulary

Use only:

- **World** — the environment daemon and the object graph it owns.
- **Handle** — a live reference to a real thing (this repo, this test suite, this staging dir). What an agent holds defines what it can do.
- **Verb** — a typed operation on a handle.
- **Result** — the typed outcome of a verb invocation, always carrying a frontier.
- **Frontier** — the ranked suggestions (≤ 3) attached to a result: verbs made relevant by this state change, arguments already bound.
- **Refusal** — a failed invocation that names what was refused, why in one line, and the reachable alternative.
- **Episode** — the append-only record of one session: acts, results, outcome.
- **Case** — a frozen replayable input plus expected outcome, promoted from an episode.
- **Friction** — an observed failure, human correction, repeated manual sequence, or budget blowout that licenses a change to the world.
- **Amendment** — a typed, replay-gated change to the world.

Do not expose contract/opportunity/receipt/grant vocabulary in schemas or product copy.

### The three laws

These are CI-enforced invariants, not style preferences:

1. **O(1) wake-up.** The standing agent-facing surface (tool description + any instructions) has a hard token budget and does not grow with the number of verbs. A world with 1000 verbs costs the same at wake-up as a world with 10.
2. **Frontier cap.** A result carries at most 3 suggestions, each with a why-line ≤ 80 characters and arguments bound to live state. No suggestion may be a paragraph.
3. **No prose rules.** Any behavioral rule must exist as topology (a handle that is or is not reachable), verb behavior, or a teaching refusal. Standing instructional text is a build failure, not a documentation style.

### Engineering standard

Generated code and hand-written code meet the same bar. Implementation work must:

- follow idiomatic Go and established standard-library patterns before introducing custom abstractions;
- keep packages cohesive, functions short enough to understand locally, and dependencies pointed in one direction;
- extract repeated behavior behind small reusable APIs once the shared behavior is stable, without creating speculative abstraction layers;
- place interfaces at the consuming boundary, keep side effects at explicit edges, and avoid god objects, cyclic responsibilities, copy-paste branches, and cross-package state mutation;
- use contextual errors, cancellation and timeouts at I/O boundaries, table-driven tests where cases share behavior, and comments only for intent or non-obvious constraints;
- pass formatting, static analysis, tests, dependency checks, and the documented maintainability review. Passing tests does not excuse duplicated, tangled, or needlessly long code.

`docs/decisions/0006-engineering-standards.md` will pin the tools, thresholds, allowed exceptions, and review checklist. An exception must name the concrete reason and the follow-up condition; generated code is never exempt.

---

## 2. Non-negotiable boundaries

Phoenix must not become:

- an agent loop or orchestrator;
- a model router or anything that calls a model;
- a prompt, skill, or tool marketplace;
- a manifest, card catalog, or capability index read at wake-up — this is now an explicit anti-goal, not a fallback;
- a permission or policy layer bolted on top of the agent;
- a chat interface or dashboard;
- a hosted dependency for local operation;
- an autonomous self-modifier: the world's own code changes only through the amendment gate.

On enforcement: Phoenix does not police a global namespace. If something must be impossible, it is **absent** — there is no verb and no reachable handle for it. Absence, not denial. Where a verb touches something dangerous, the verb's own behavior and refusal shape carry the constraint. Verbs execute under the daemon's OS principal; the daemon's principal, not policy text, bounds what verbs can physically do. Phoenix claims exactly this and nothing stronger: agents interacting **through the surface** can only do what the world contains. An agent with its own shell access is outside Phoenix's boundary, and every experiment must state which access mode was measured.

---

## 3. Risk order

1. **Frontier quality** — bad `next:` is advice-spam; the entire bet is that in-band beats upfront. If suggestions are ignored or mislead, nothing downstream matters.
2. **Ground truth for "better"** — the case bank. Without replayable cases with expected outcomes, every change to the world is drift with good intentions.
3. **Far-capability discovery** — if `find`-as-an-action fails, agents get stuck in local neighborhoods and the O(1) wake-up becomes a cage.
4. **Transition-function sprawl** — the world's code becoming the new convoluted text, unreviewable and untestable.
5. **Learning-loop gaming and noise** — counting reinforcing plausible-but-wrong paths; attribution errors crediting the wrong suggestion.
6. **Ossification** — rich-get-richer frontier edges freezing the world into early habits; exploration must be budgeted, not hoped for.
7. **Platform creep** — multi-agent, network effects, and dashboards before the single-agent bet is proven.

---

## 4. Repository target structure

```text
phoenix/
├── README.md
├── LICENSE                         # add only after Raoul selects a license
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── phoenix/
│       └── main.go                 # daemon + CLI (serve, replay, episodes, world)
├── internal/
│   ├── world/                      # object graph, handles, reachability
│   ├── verb/                       # verb registry, typed results, execution
│   ├── frontier/                   # transition rules, ranking, decay, counters
│   ├── teach/                      # refusal shaping
│   ├── episode/                    # append-only log (SQLite WAL)
│   ├── casebank/                   # frozen cases + replay runner
│   ├── amend/                      # typed amendments + gate (Phase 2)
│   └── surface/                    # MCP server, the single act surface
├── verbs/                          # verb implementations (the part that grows)
├── worlds/
│   └── dev-repo/                   # starter world definition for a Go/git repo
├── experiments/
│   └── frontier-v1/
│       ├── README.md
│       ├── protocol.json
│       ├── corpus/{authoring,validation,held_out}/
│       ├── labels/authoring/
│       ├── manifests/              # digests only for withheld labels
│       ├── runner/
│       └── reports/
├── spec/
│   ├── result.schema.json          # result envelope incl. frontier + refusal
│   ├── episode.schema.json
│   ├── case.schema.json
│   ├── world.schema.json           # world definition format
│   └── amendment.schema.json       # Phase 2
├── testdata/
└── docs/
    ├── decisions/
    ├── plans/
    └── experiment-results.md
```

Do not create deferred directories until their phase begins.

**Initial tech stack:** Go with the toolchain version pinned in `go.mod` and CI; official MCP Go SDK over stdio; pure-Go SQLite in WAL mode (CGO disabled); JSON Schema Draft 2020-12 for spec files; RFC 8785 canonical JSON for stable protocol digests; content-addressed episode, case, and world-build digests. A world-build digest covers the daemon executable, registered verb set, schemas, world definition, authored rules, and active frontier weights. All of these are defaults to be confirmed or overturned by the Phase 0 surface spike — record the outcome either way. Linux-first; nothing in Phase 0–2 may require kernel features beyond portable POSIX unless a decision doc justifies it.

---

# Phase 0 — Design commitments and baseline corpus

## Phase 0 objective

Pre-commit what the world must prove before building it. The deliverables are a frozen surface design, a task corpus with withheld labels, and go/no-go rules — not application code.

## Phase 0 gate

Phase 0 is accepted only when:

- the surface spike has produced a working end-to-end round trip with the pinned runtime and a measured standing token cost;
- the result envelope, frontier shape, and refusal shape are frozen as schemas;
- authoring, validation, and held-out corpora are separated by generating family, with labels rationale-carrying and independently reviewed;
- success/token/latency thresholds and all kill rules are committed before any results exist;
- Arms D (frontier suppressed) and E (plain errors) remain in the experiment;
- the experiment can conclude that only part of the system is useful — surface without frontier, frontier without learning — and kill the rest.

### Task 0.1: Record the system boundary

**Files:**
- Create: `docs/decisions/0001-system-boundary.md`

**Steps:**
1. State that Phoenix never calls models, never orchestrates, and holds no chat state.
2. State the three laws with their numeric budgets: standing surface ≤ 600 tokens; frontier ≤ 3 entries, why-line ≤ 80 chars.
3. Define absence-not-denial: the enforcement story is topology plus the daemon's OS principal; no policy layer exists.
4. State the measured claim precisely: Phoenix bounds agents interacting through the surface; agents with independent shell access are outside the boundary, and experiments must declare which access mode each arm used.
5. Review against the non-negotiable list in section 2.
6. Commit: `docs: define Phoenix world boundary`.

### Task 0.2: Surface spike — one tool, one runtime

**Files:**
- Create: `docs/decisions/0002-surface-and-stack.md`
- Create: `experiments/surface-spike/README.md`
- Create: `experiments/surface-spike/results.json`

**Default under test:** a single MCP tool `act` served over stdio to Claude Code, with results rendered as compact text plus a structured frontier block.

**Steps:**
1. Pin exact runtime version, model, configuration, and invocation command for the spike.
2. Build a throwaway daemon exposing `act(handle, verb, args)` with two hardcoded handles and ~6 verbs; no learning, no persistence.
3. Compare two surface candidates: structured `act` vs. a code-REPL surface (`eval` over an object graph). Measure: standing token cost of each surface, malformed-call rate, and whether the runtime renders frontier text legibly to the model.
4. Confirm the Go + MCP SDK + stdio default, or record the overturning decision (e.g., TypeScript SDK maturity) with reasons.
5. Verify the O(1) property mechanically: generate synthetic worlds with 10 / 100 / 1000 verbs and confirm the standing surface token count is flat while an ordinary tool-schema dump grows linearly. Keep the harness — it becomes a CI test in Task 1.4.
6. Record known failure modes: what happens on daemon crash mid-act, oversized results, and concurrent sessions.
7. Commit: `experiment: verify Phoenix surface and stack`.

### Task 0.3: Freeze the result envelope

**Files:**
- Create: `spec/result.schema.json`
- Create: `spec/world.schema.json`
- Create: `docs/decisions/0003-result-envelope.md`

**Envelope shape (illustrative, frozen by the schema):**

```json
{
  "v": 1,
  "world_build": "sha256:…",
  "session_id": "…",
  "act_id": "…",
  "handle": "h_tests_…",
  "verb": "run",
  "status": "ok | fail | refused | absent",
  "result": { "…": "typed payload per verb" },
  "error": null,
  "text": "compact rendering the model reads",
  "handles": {
    "grant": [{ "ref": "h_flake_…", "type": "flake", "label": "auth failure" }],
    "revoke": []
  },
  "frontier": [
    {
      "call": {
        "handle": "h_flake_…",
        "verb": "bisect",
        "args": { "test": "auth_test" },
        "state": "sha256:…"
      },
      "why": "seen this failure shape 4x, resolved 3",
      "provenance": "authored | counted",
      "score": 0.71
    }
  ],
  "refusal": {
    "what": "…",
    "why": "…",
    "instead": { "handle": "h_tests_…", "verb": "focus", "args": { "test": "auth_test" } }
  }
}
```

**Steps:**
1. Define the envelope, structured call, handle-transition, frontier entry, error, and refusal shapes as JSON Schema with unknown-field rejection and status-dependent required fields.
2. Represent calls as typed `handle`/`verb`/`args` objects. Rendering is derived output; validation, execution, reachability, and attribution never parse the human-readable `text` or a call string.
3. Require frontier arguments to be bound to live state. A suggestion may never contain a placeholder the agent must invent; state-sensitive suggestions carry a precondition digest and are revalidated at invocation.
4. Make handle references opaque and session-scoped. The runtime receives the compact root-handle set at session start within the standing token budget; later handles enter or leave reachability only through explicit `grant`/`revoke` result fields. A guessed, revoked, or cross-session reference produces the same typed `absent` result.
5. Require every refusal to carry a structured `instead` call referencing a reachable handle/verb, or explicitly state that no alternative exists.
6. Define `provenance` semantics: `authored` (hand-written transition rule) vs `counted` (replay-gated ranking, Phase 2). The field exists from day one so arms can be distinguished later.
7. Define the world-definition format: which root handles exist for a project, which verbs attach to each handle type, which results can grant or revoke handles, and the initial authored transition rules.
8. Define canonicalization and digest inputs for envelopes, calls, world builds, and episode references; validate example documents against the schemas with a standalone fixture.
9. Commit: `spec: freeze result envelope and frontier shape`.

### Task 0.4: Build the task corpus

**Files:**
- Create: `experiments/frontier-v1/README.md`
- Create cases under `experiments/frontier-v1/corpus/{authoring,validation,held_out}/`
- Create authoring labels under `experiments/frontier-v1/labels/authoring/`
- Commit only digests of withheld labels under `experiments/frontier-v1/manifests/`

**Runnable input shape (contains no answers):**

```yaml
case_id:             # opaque tranche-scoped identifier; never encode the class
goal:                # the task given to the agent
sandbox_fixture:     # repo snapshot digest the task runs against
world_ref:           # world definition digest
family_id:
```

**Withheld label shape (stored outside the implementation workspace until evaluation):**

```yaml
case_id:
class:               # evaluator-only metadata; never handed to the agent
expected_outcome:    # machine-checkable end state where possible
grading_script:      # digest of the script that grades the end state
acceptable_paths:    # verb sequences considered correct (non-exhaustive)
label_rationale:
```

**Required case classes:**

- **direct** — the needed verb is obvious; measures overhead, not discovery;
- **cascade** — success needs 3+ verbs in an order the agent would not guess upfront; the frontier's home turf;
- **far-discovery** — the needed verb is obscure; tests `find`-as-an-action;
- **recovery** — the natural first attempt fails; tests whether teaching refusals convert errors into progress;
- **temptation** — a plausible-but-wrong verb exists; measures wrong-verb rate;
- **absence** — the capability genuinely does not exist; correct behavior is saying so promptly, not flailing;
- **stale-frontier** — the sandbox has changed such that a previously good suggestion is now wrong; tests decay and refusal honesty;
- **adversarial text** — task or fixture content tries to talk the agent into ignoring or abusing the surface.

**Steps:**
1. Write authoring cases first, then a sealed validation tranche, then a final held-out tranche; split by generating family, never by random individual case.
2. Prefer scripted grading (end-state checks) over judgment; where human grading is unavoidable, define blinding.
3. Have an independent evaluator label validation and held-out cases without seeing the implementation; preserve disagreements as adjudication records.
4. Commit input digests and withheld-label digests before implementation; prohibit tuning against a tranche after it is opened.
5. Add a corpus manifest with counts, class coverage, and explicit denominators — not outcome results.
6. Treat an opened tranche as consumed if its outcomes inform any corpus, world, frontier, refusal, runner, or protocol change. Retire it and generate a new independently labeled and sealed tranche before another gate attempt; never rerun a repaired system against the same tranche as if it remained validation or held-out evidence.
7. Commit authoring inputs and sealed manifests separately.

### Task 0.5: Pre-commit experiment arms and decision rules

**Files:**
- Create: `experiments/frontier-v1/protocol.json`
- Create: `docs/decisions/0004-go-no-go-rules.md`

**Arms:**

All arms use the same pinned runtime/model/configuration, sandbox fixtures, retry policy, and cost accounting. Each trial starts with a fresh agent context and the protocol-declared daemon/world state. Runs use paired seeds where the runtime supports them, enough repetitions for the pre-committed uncertainty bound, and randomized or counterbalanced case order. All arms are measured on identical tasks and backed by the same verb implementations, access mode, side effects, limits, and typed result payloads. The decision doc must list any irreducible surface difference and narrow the claim accordingly.

- **A:** agent + conventional flat toolset with standard upfront schemas (status quo baseline);
- **B:** Arm A plus the best reasonable frozen static skill file / CLAUDE.md — the strong documentation baseline;
- **C:** agent + Phoenix world: handles, frontier, teaching refusals (Phase 1: authored frontier; Phase 2 rerun: counted frontier);
- **D:** identical Phoenix world with the frontier suppressed (handles and refusals only) — isolates the frontier's causal value from the surface itself;
- **E:** identical to D but state refusals become plain typed errors with no suggested alternative — isolates the teaching refusal's causal value;
- **D′ (Phase 2 only):** frontier present but ranking unlearned/static — isolates learning's value over authored suggestions.

**Primary measures:**

- task success rate (scripted grading);
- tokens-to-success (all turns, prompt + completion), including the paired-success ratio when both compared arms succeed;
- acts-to-success and wall time;
- wrong-verb invocation rate;
- dead-end rate (agent stalls or asks for help);
- timeout-or-cap-hit rate over all assigned trials;
- recovery rate after refusal;
- frontier take-rate and take-and-succeed rate (frontier-bearing arms only);
- standing wake-up cost across synthetic world sizes (the O(1) demonstration);
- daemon overhead per act.

**Steps:**
1. Create a falsifiable claim table: for every claim, record population, comparator, unit of analysis, endpoint, numeric threshold, uncertainty bound, minimum detectable effect, repetitions, seed handling, context and state reset, ordering, exclusions, missing/indeterminate handling, retry clustering, multiplicity policy, stopping rule, and forced failure verdict. Thresholds are committed here, before any arm runs.
2. Define the direct-task harm gate: test C versus A on all assigned **direct** trials at one-sided alpha 0.05. A detected harm rejects or repairs the surface; a non-significant result does not establish noninferiority.
3. Define component-level decisions:
   - C versus B ITT success is the single Phase 1 confirmatory claim. Cost measures are descriptive and cannot veto a capability pass;
   - if the capability bound does not clear zero, an efficiency-only result requires the pre-committed success, paired-success token, dead-end, and timeout/cap-hit guardrails; report it only as an efficiency product;
   - C does not improve on D in the frontier isolation → kill the frontier; keep the surface only if the headline passed;
   - D does not improve on E over all assigned recovery trials, or harms downstream success → replace teaching refusals with plain typed errors;
   - (Phase 2) counted frontier does not beat authored frontier on held-out cases → keep counting off and hand-author transitions;
   - result indeterminate → expand or repair the corpus; do not promote architecture.
4. Scope every conclusion to the tested runtime, world, and task classes.
5. Commit: `docs: precommit frontier experiment decisions`.

### Task 0.6: Phase 0 independent review

**Files:**
- Create: `docs/decisions/0005-phase-0-review.md`

**Steps:**
1. An evaluation reviewer independent of implementation reviews corpus validity, leakage controls, labels, and comparison design.
2. An implementation reviewer validates surface feasibility, envelope schemas, and the O(1) harness.
3. A human-factors reviewer checks frontier and refusal legibility: can a reader reconstruct from an episode why the agent did what it did.
4. The accountable project chair records accepted findings and unresolved blockers.
5. Do not begin Phase 1 while the corpus or go/no-go rules remain unaccepted.
6. Commit: `docs: accept Phoenix Phase 0 research contract`.

---

# Phase 1 — The static world

## Phase 1 objective

Prove the surface thesis **before** building any learning: a world with hand-authored frontiers and teaching refusals, one real runtime, one starter world for a real dev repo. If in-band delivery cannot beat upfront documentation with authored suggestions, counting will not save it.

## Explicit omissions

Do not build in Phase 1:

- learning of any kind (counters, decay, reranking);
- the amendment gate;
- shared mutable world state across sessions; Phase 1 may admit concurrent isolated sessions, but they share no handles, counters, or topology mutations;
- network-touching verbs beyond what the starter world's build/test path already requires;
- a second runtime adapter;
- a dashboard, TUI, or web anything;
- WASM, plugins, or a verb marketplace.

### Task 1.1: Initialize the minimal Go project

**Files:**
- Create: `go.mod`, `Makefile`, `.gitignore`
- Create: `cmd/phoenix/main.go`
- Create: `cmd/phoenix/main_test.go`

**Steps:**
1. Pin the Go toolchain, MCP SDK, pure-Go SQLite, JSON Schema, and static-analysis versions; record the CGO-disabled policy.
2. Write a failing CLI smoke test for `phoenix version` and `phoenix serve --stdio` handshake.
3. Implement only the command shell required by the test.
4. Create `docs/decisions/0006-engineering-standards.md`; pin formatting, static-analysis, dependency, duplication/complexity, test, and maintainability-review rules, including the exception format.
5. Add pinned `make test`, `make lint`, `make build`, `make validate-spec`, and `make quality` targets with documented exit statuses. `make quality` is the CI entry point for the engineering standard.
6. Produce a content-addressed build manifest covering the executable, registered verbs, schemas, world definition, rules, and weights.
7. Commit: `build: initialize Phoenix daemon`.

### Task 1.2: Implement the object graph and handles

**Files:**
- Create: `internal/world/graph.go`
- Create: `internal/world/load.go`
- Create: `internal/world/graph_test.go`

**Steps:**
1. Test loading a world definition against `spec/world.schema.json`; unknown fields rejected; content digest stable.
2. Test that session start issues a compact set of opaque, session-scoped root references and that results grant or revoke later references explicitly.
3. Test reachability: guessed, revoked, or cross-session references do not exist for that session — verb dispatch returns the same typed *absent* result, not *denied* and not a namespace oracle.
4. Test that handles reference live state (paths, processes) resolved at act time, not stale snapshots, and that state-sensitive calls reject stale preconditions cleanly.
5. Implement the smallest graph satisfying the tests. No global verb namespace anywhere; `find` searches only the session's reachable topology and may reveal a new handle only through a declared grant edge.
6. Commit: `feat: load worlds of reachable handles`.

### Task 1.3: Implement the verb registry and typed results

**Files:**
- Create: `internal/verb/registry.go`
- Create: `internal/verb/exec.go`
- Create: `internal/verb/exec_test.go`
- Create: `verbs/dev-repo/…` (starter verbs)

**Steps:**
1. Test that every verb declares its handle type, argument schema, and result type; registration of a verb with untyped results fails.
2. Test execution of in-process verbs and allowlisted external executables (identified by digest); no arbitrary exec.
3. Test that verb failures produce structured `fail` results, never raw stack traces.
4. Test timeout, output size limits, and secret redaction on results before persistence.
5. Implement the registry and the starter verb set for the dev-repo world: `repo.read/edit/build`, `tests.run/focus`, `git.status/diff/commit`, `find` (search verbs and handles as an action), `episodes.recall(query)`.
6. Commit: `feat: execute typed verbs on handles`.

### Task 1.4: Implement the surface and enforce the three laws

**Files:**
- Create: `internal/surface/mcp.go`
- Create: `internal/surface/mcp_test.go`
- Create: `internal/surface/budget_test.go`

**Steps:**
1. Implement the single `act` tool over stdio per the spike decision; results rendered as compact text plus the structured envelope.
2. Test malformed calls, unknown handles, oversized args, and concurrent sessions at one admission point.
3. **Law 1 test:** the synthetic-world harness from Task 0.2 runs in CI; standing surface tokens must be flat across 10/100/1000-verb worlds and under the committed budget. Build fails otherwise.
4. **Law 3 test:** the surface exposes no standing instructional text beyond the budgeted tool description; grep-level CI check on the served surface.
5. Commit: `feat: serve one constant act surface`.

### Task 1.5: Implement the frontier engine (authored rules only)

**Files:**
- Create: `internal/frontier/rules.go`
- Create: `internal/frontier/rank.go`
- Create: `internal/frontier/frontier_test.go`

**Steps:**
1. Test transition rules: match on (verb, status, result features) → candidate structured calls with reachable handles and bound arguments; a rule producing an unbound placeholder, opaque call string, or unreachable handle fails validation.
2. **Law 2 test:** cap of 3 enforced at the engine, not by convention; why-lines over 80 chars rejected at rule load.
3. Test determinism: identical state and rule set produce identical frontiers.
4. Test that an empty frontier is legal and common — no padding with weak suggestions.
5. Author the initial transition rules for the dev-repo world (build fail → likely fixes; test fail → bisect/recall; edit → affected tests).
6. Commit: `feat: compute authored frontiers`.

### Task 1.6: Implement teaching refusals

**Files:**
- Create: `internal/teach/refusal.go`
- Create: `internal/teach/refusal_test.go`

**Steps:**
1. Test that every refusal carries what/why/instead, with structured `instead` validated as reachable from the session's handles or explicitly null with a stated reason.
2. Test the distinction between *absent* (handle/verb not in topology — short, factual) and *refused* (verb declined given state — teaches the alternative).
3. Convert the starter world's constraints into refusal shapes; zero standing rule text remains (Law 3).
4. Commit: `feat: shape refusals that teach`.

### Task 1.7: Implement the episode log

**Files:**
- Create: `internal/episode/log.go`
- Create: `internal/episode/sqlite.go`
- Create: `internal/episode/log_test.go`
- Create: `spec/episode.schema.json`

**Steps:**
1. Test append-only ordering, session correlation, idempotent re-append on retry, and crash-recovery of a half-written episode to an explicit `interrupted` state.
2. Test that every act's structured request, world-build digest, result digest, handle grants/revocations, frontier shown, and suggestion-taken linkage are reconstructable without parsing rendered text.
3. Test secret redaction before persistence and WAL-mode concurrent writers.
4. Test `episodes.recall(query)` returns pointers into past episodes, never free-floating prose summaries.
5. Commit: `feat: record reconstructable episodes`.

### Task 1.8: Assemble the starter world and run the corpus runner

**Files:**
- Create: `worlds/dev-repo/world.json`
- Create: `experiments/frontier-v1/runner/main.go`
- Create: `experiments/frontier-v1/runner/main_test.go`

**Steps:**
1. Define the dev-repo world: root handles (`repo`, `tests`, `git`, `episodes`), verbs, authored transitions, refusal shapes.
2. Build a harness-neutral runner: instantiates a sandbox fixture, serves the world, drives the pinned runtime on a case, grades the end state with the case's grading script. It must not require any specific agent framework beyond the pinned runtime under test.
3. Run the **authoring** tranche only; keep validation and held-out tranches sealed.
4. Fix world/frontier/refusal defects found on authoring cases freely — this is the tuning set.
5. Commit: `feat: run frontier corpus against the static world`.

### Gate 1A: Surface-thesis stop gate

Before any learning code is written:

1. Freeze the world definition, rules, daemon, schemas, and active weights; commit the complete world-build digest.
2. Open only the sealed validation tranche.
3. Run Arms A, B, C-static, D, and E with the pre-committed protocol.
4. Report all primary measures with explicit denominators, including the O(1) wake-up demonstration.
5. **Stop conditions (from 0.5):** if the C-static versus B ITT-success primary claim fails and the separate efficiency-only rule also fails, stop — repair the corpus or surface, or narrow Phoenix per the component decisions. Cost cannot veto a capability pass or rescue cheap failure by itself. Any repair informed by validation outcomes burns that tranche and requires a newly generated, independently sealed validation tranche before another gate attempt. Do not proceed to Phase 2 on the theory that learning will close the gap.
6. Continue only after the independent evaluation reviewer accepts the result.

---

# Phase 2 — The world that learns

Phase 2 begins only if Gate 1A passes.

## Phase 2 objective

Add the counting loop and the amendment gate: frontiers reranked by observed take-and-succeed, friction converted into candidate verbs and rules, and no change landing without replay. The agent stays frozen; the world's code is the only thing that improves.

## Explicit omissions

- autonomous landing outside the fenced amendment types;
- any LLM-judge in the acceptance path where scripted grading exists;
- shared/multi-session worlds (Phase 3);
- a second runtime.

### Task 2.1: Implement take/succeed counters and decay

**Files:**
- Create: `internal/frontier/counters.go`
- Create: `internal/frontier/counters_test.go`

**Steps:**
1. Define attribution: a suggestion is *taken* if a subsequent structured call in the same episode canonically matches its handle, verb, and arguments within a fixed window; it *succeeds* only when a scripted grader or independently adjudicated outcome is positive. Agent narration or an unverified close-out never increments success.
2. Test counter updates from eligible episode replay, multiplicative decay per idle period, and floor-at-zero. Record the training-episode digest and exclude held-out episodes.
3. Keep counters observational and the active weight set immutable. Test reranking against a candidate weight snapshot: `counted` scores may reorder candidates but cannot exceed the cap or resurrect rule-rejected suggestions.
4. Test the exploration budget as part of a versioned candidate weight configuration. Evaluated tranches freeze that configuration and its seed; no online counter or weight update occurs during a gate run.
5. Emit counter changes as a typed `reweight` amendment. The active frontier changes only after Task 2.3 replays and admits that amendment.
6. Commit: `feat: count, decay, and propose frontier weights`.

### Task 2.2: Implement the case bank and replay

**Files:**
- Create: `internal/casebank/bank.go`
- Create: `internal/casebank/replay.go`
- Create: `internal/casebank/replay_test.go`
- Create: `spec/case.schema.json`

**Steps:**
1. Test promotion: a failed or human-corrected episode freezes into a case (sandbox fixture digest, goal, expected outcome, grading script).
2. Implement model-free replay primitives in `phoenix replay`: materialize a sandbox fixture, serve a candidate world build, replay deterministic act traces where available, run grading scripts, and emit signed/digested inputs and results. The command never launches or calls a model.
3. Extend the external experiment runner from Task 1.8 to drive the pinned runtime for behavioral replays, using the model-free Phoenix primitives. This runner is evaluation infrastructure outside the daemon boundary; it supplies the report consumed by the amendment gate.
4. Test case retirement: a case every candidate passes for N runs is flagged non-discriminating (retired from the main gate, retained in history and sampled as a regression sentinel).
5. Test that deterministic and behavioral replay never mutate real project state — sandbox fixtures only.
6. Commit: `feat: replay the case bank against candidate worlds`.

### Task 2.3: Implement typed amendments and the gate

**Files:**
- Create: `internal/amend/amend.go`
- Create: `internal/amend/gate.go`
- Create: `internal/amend/gate_test.go`
- Create: `spec/amendment.schema.json`

**Amendment types (closed set):**

```text
add_verb          # new verb + its tests + candidate executable/build manifest + cost comparison
retire_verb
add_transition    # new authored frontier rule
reweight          # batched counter-driven rank changes
add_refusal       # new teaching shape
edit_world        # handle topology change — strongest review tier
```

**Steps:**
1. Test that every amendment references its friction evidence (episode/correction pointers) and carries a replay report; amendments without either are rejected at admission.
2. Test the gate: non-regression on the case bank plus improvement on the target friction's cases; otherwise auto-reject with the report preserved.
3. Test author/grader separation structurally: the session that proposed an amendment cannot grade it; grading runs in a fresh context against the frozen candidate.
4. Test deterministic rollback: the active world points to an immutable build manifest covering the daemon/verb executable, schemas, topology, rules, and weights. A supervisor activates or rolls back the exact build with a controlled restart; a config pointer alone cannot roll back compiled Go code. Historical episodes stay pinned to the build digests they ran under. Phase 2 does not load Go plugins in-process.
5. Test the idle invariant: zero friction in, zero amendments out — a system generating changes without friction evidence fails the test suite.
6. Autonomy levels: `reweight` may land automatically within fences only after its candidate snapshot passes replay; `add_verb`/`add_transition`/`add_refusal` land into a weekly human-reviewed batch; `edit_world` and `retire_verb` require explicit approval. Fences are config the amendment process itself cannot widen.
7. Commit: `feat: gate typed amendments through replay`.

### Task 2.4: Implement the friction detector

**Files:**
- Create: `internal/amend/friction.go`
- Create: `internal/amend/friction_test.go`

**Steps:**
1. Test detection of: the two-times rule (same manual verb sequence ≥ 2 episodes), human corrections, refusal-then-abandon patterns, and cost blowouts vs. class baseline.
2. Test that detected friction produces a queue entry with evidence pointers — never an auto-generated amendment.
3. Commit: `feat: detect friction worth amending`.

### Gate 2A: Learning stop gate

1. Freeze the complete world-build digest, including the training-episode set and candidate weights, then rerun the experiment on the sealed held-out tranche: C-counted vs C-authored (D′) vs B. No counters, weights, rules, or code update during the tranche.
2. Apply the pre-committed component decisions from 0.5: if counting does not beat authored ranking on held-out cases, keep counting off. If C beats B, the headline bet stands with learning; publish machine results and interpretation separately.
3. Report gaming checks: held-out cases the counting loop never saw, and the idle invariant over a friction-free week.
4. Independent review as in 0.6; the chair selects exactly one verdict: proceed to Phase 3 / keep static world / repair / kill. A repair informed by held-out outcomes burns that tranche; any later headline claim requires a newly generated, independently sealed held-out tranche.

---

# Phase 3 — The shared world

Phase 3 begins only after Gate 2A and only if real multi-session demand exists (more than one agent or person actually working in the same world).

## Objectives

- multiple sessions inhabit one world; handles flow between sessions through explicit handoff envelopes (intent, handles granted, expected return artifact);
- episodes attributed per principal; no session grades or lands its own amendments;
- frontier counters shared, so one session's discovered path improves every session's ranking;
- per-task-class outcome tracking replaces global scores.

## Invariants

- a handoff cannot widen the receiving session's topology beyond the granting session's;
- shared counters cannot be written except through graded episodes;
- a subordinate session's success claim requires an inspectable artifact, not narration.

## Phase 3 gate

The shared world must beat isolated per-session worlds on repeated task classes after coordination cost. More sessions or messages do not count as improvement.

---

# Verification strategy

## Every code task

1. Define the owning package, dependency direction, side-effect boundary, and reusable seam before adding code; do not duplicate an existing path.
2. Write a focused failing test; run it; confirm the expected failure.
3. Implement the smallest cohesive behavior with idiomatic Go. Keep functions and types locally understandable; split mixed responsibilities instead of growing a dispatcher, manager, or utility grab-bag.
4. Run focused tests, then the full suite, formatting, static analysis, dependency checks, duplication/complexity checks, and a CGO-disabled build through `make quality`.
5. Review the diff for copy-paste branches, needless abstractions, hidden global state, long control-flow chains, and package leakage. Refactor before acceptance; generated code receives the same review.
6. Inspect side effects in a sandbox; commit only the verified increment.

## Standing CI invariants

- Law 1: flat standing-surface tokens across synthetic world sizes, under budget;
- Law 2: frontier cap and why-line length enforced in the engine, covered by tests;
- Law 3: no standing instructional text on the served surface;
- secret redaction before persistence;
- refusals always carry a valid `instead` or an explicit null;
- replay never touches non-sandbox state;
- the idle invariant: no friction → no amendments;
- daemon restart preserves reconstructable episodes;
- the engineering-quality target passes: formatting, static analysis, dependency rules, documented duplication/complexity thresholds, and maintainability review, with no exemption for generated code.

## Neighboring behavior

Each phase must verify, beyond its own features:

- direct tasks are not slowed or degraded by the surface (no tax on the easy path);
- an empty frontier does not strand the agent — `find` and plain competence still work;
- malformed or adversarial fixture text cannot alter the surface, frontier rules, or counters;
- duplicate/retried acts remain idempotent in the episode log;
- a corrupted or unavailable episode DB degrades to a working world with logging disabled and a loud warning — never to silent partial logging.

---

# Sequencing and review ownership

Phoenix follows the currently prioritized research-kernel milestone unless Raoul explicitly changes priority. TradePilot and other projects may contribute generalized friction traces; Phoenix must not import their code, schemas, or assumptions.

Roles are implementation-independent:

- **Accountable project chair** — scope, gates, acceptance, verification that claimed artifacts exist.
- **Evaluation reviewer** — corpus, leakage, statistics, causal claims; must not implement what they grade.
- **Implementation owner** — TDD, daemon, surface, durability; cannot approve their own gate.
- **Human-factors reviewer** — frontier and refusal legibility, episode reconstruction, review fatigue in the amendment batch.

The repository and experiment must remain reproducible without these particular reviewers or any named agent roster. Production deployment, hosted services, and public release require separate approval.

---

# Final delivery gates

Phoenix may call itself successful only when:

1. Wake-up cost is demonstrated flat across world sizes while baselines grow linearly.
2. Held-out results beat the strong static-documentation baseline (Arm B) after cost — the headline gate.
3. C versus D isolates the frontier's causal contribution; the frontier survives only if that contribution is real.
4. D versus E shows that teaching refusals convert failures into recoveries better than plain typed errors.
5. Every result is reconstructable from episodes pinned to complete world-build digests.
6. Amendments improve held-out outcomes without self-grading, with activation and rollback of the exact immutable build, and with zero amendments under zero friction.
7. A shared world beats isolated worlds on repeated task classes before Phase 3 is retained.
8. The accepted codebase passes the engineering standard: idiomatic patterns, cohesive packages, short reusable units, controlled side effects, no unjustified duplication, and no tangled cross-package control flow.

If only part survives, Phoenix becomes exactly that part: a constant surface without frontiers, or an authored world without counting. The architecture serves evidence; evidence does not serve the architecture.
