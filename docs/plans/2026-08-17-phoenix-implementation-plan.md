# Phoenix Implementation Plan

> **SUPERSEDED (2026-08-17):** replaced by [`2026-08-17-phoenix-world-plan.md`](2026-08-17-phoenix-world-plan.md), which reframes Phoenix from a capability control plane into an environment daemon that agents inhabit. This document is retained for the measurement discipline and decisions it records; do not implement from it.

> **For implementers:** execute this plan phase-by-phase using TDD, independent review, and the acceptance gates below. Phoenix must remain standalone and must not import or depend on Hermes or any other agent harness.

**Goal:** Build a local-first, runtime-neutral capability control plane that detects when behavioral contracts should apply, makes missed activation observable, mediates protected effects, requires evidence, and later evolves contracts only through replayed and approved amendments.

**Architecture:** Phoenix is a standalone process. External agents and runtimes connect through narrow adapters over JSON via stdio, Unix sockets, or localhost HTTP. Phoenix owns contract resolution, effect mediation, evidence receipts, evaluation records, and an append-only ledger; it does not invoke models, plan tasks, orchestrate agents, or own chat sessions.

**Initial tech stack:** Linux amd64; Go with the supported version pinned in `go.mod` and CI during Task 1.1; pure-Go SQLite in WAL mode; RFC 8785 JSON Canonicalization Scheme for hashed protocol objects; JSON Schema Draft 2020-12; content-addressed filesystem artifacts; in-toto-compatible evidence statements where practical; ordinary deterministic Go checks in Phase 1. Phase 0 must pin every dependency, tool, invocation command, timeout, and runtime/model configuration before the experiment. WASM, distributed infrastructure, hosted services, and generalized policy engines are deferred.

---

## 1. Product decision

Phoenix is not “dynamic skills” and not another agent framework.

The core unit is a **contract**:

> A versioned behavior contract that can recognize an applicable opportunity, constrain a protected effect, require typed evidence, emit a reconstructable decision, and be evaluated against a pre-committed outcome.

Static documents may remain guidance payloads. They are not authoritative activation, policy, evidence, or evaluation mechanisms.

### Public vocabulary

Use only:

- **Phoenix** — the product and standalone control plane.
- **Contract** — the versioned behavioral unit.
- **Opportunity** — an observed situation in which one or more contracts may apply.
- **Decision** — activate, reject, require evidence, allow, deny, or indeterminate.
- **Run receipt** — the human-readable and machine-readable causal record.
- **Amendment** — a proposed future contract change, deferred until Phase 3.

Do not expose the earlier Charter/fixture/plate/kerf/scar vocabulary in schemas or product copy.

---

## 2. Non-negotiable boundaries

Phoenix must not become:

- an agent loop or orchestrator;
- a model router;
- a prompt or skill marketplace;
- a workflow designer;
- a chat interface;
- a universal memory layer;
- an autonomous self-modification engine;
- a hosted dependency for local operation;
- a system that claims enforcement it cannot provide.

Phoenix may enforce only effects that pass through a boundary it controls. Every adapter and contract must declare one of these strengths:

1. **Advisory** — Phoenix reports a decision; bypass remains possible.
2. **Cooperative** — the normal runtime path blocks denied effects; alternate paths may remain.
3. **Mediated** — Phoenix controls the effect gateway or credential and prevents bypass through that boundary.

The run receipt must disclose both the declared strength and known bypass paths.

---

## 3. Correct risk order

Implementation priorities follow this order:

1. **Activation opportunity recognition** — contracts that never engage are silent failures.
2. **Ground-truth corpus quality** — weak labels make every comparison meaningless.
3. **Context integrity** — inferred context cannot silently become authoritative fact.
4. **False blocking and conflict deadlock** — over-enforcement causes users to disable the system.
5. **Effect-boundary completeness** — mediated enforcement is real only at controlled boundaries.
6. **Evidence validity** — provenance proves lineage, not truth.
7. **Action canonicalization** — network/API effects are materially harder than files.
8. **Human amendment fatigue** — future review flows can degrade into rubber-stamping.
9. **Platform creep** — infrastructure must not outrun demonstrated value.

### Silence semantics

Phoenix distinguishes three cases:

```text
no opportunity + no activation = potentially normal
opportunity + no activation    = activation failure
expected cadence + no opportunity = pipeline health failure
```

Phoenix reports four distinct observation classes:

- **gateway-observed opportunity** — the Phoenix-owned effect boundary saw the request;
- **adapter-observed opportunity** — a cooperative runtime hook reported the request;
- **independently observed attempt** — the experiment harness or OS-level observer counted an actual attempt outside the detector under test;
- **unobservable bypass** — a documented path outside Phoenix's visibility.

Live Phoenix cannot claim to detect attempts outside its observation boundary. End-to-end missed activation is measurable only when an independent observer supplies the actual-attempt denominator. Cadence health is deferred from Phase 1 unless the selected capability has a justified recurring opportunity source.

---

## 4. Repository target structure

The first implementation should converge on this tree:

```text
phoenix/
├── README.md
├── LICENSE                         # add only after Raoul selects a license
├── go.mod
├── go.sum
├── Makefile
├── cmd/
│   └── phoenix/
│       └── main.go
├── internal/
│   ├── activation/
│   │   ├── engine.go
│   │   ├── predicates.go
│   │   └── engine_test.go
│   ├── admission/
│   │   ├── admission.go
│   │   └── admission_test.go
│   ├── adapter/
│   │   ├── protocol.go
│   │   └── conformance_test.go
│   ├── artifact/
│   │   ├── store.go
│   │   └── store_test.go
│   ├── contract/
│   │   ├── model.go
│   │   ├── registry.go
│   │   ├── validate.go
│   │   └── registry_test.go
│   ├── evidence/
│   │   ├── statement.go
│   │   ├── verify.go
│   │   └── verify_test.go
│   ├── gateway/
│   │   └── filesystem/
│   │       ├── gateway.go
│   │       ├── path.go
│   │       └── gateway_test.go
│   ├── ledger/
│   │   ├── ledger.go
│   │   ├── sqlite.go
│   │   ├── migrations.go
│   │   └── sqlite_test.go
│   ├── opportunity/
│   │   ├── model.go
│   │   └── monitor.go
│   ├── receipt/
│   │   ├── render.go
│   │   └── render_test.go
│   ├── server/
│   │   ├── service.go
│   │   ├── <selected-transport>.go
│   │   └── service_test.go
│   └── token/
│       ├── grant.go
│       └── grant_test.go
├── adapters/
│   ├── claude-code/
│   │   ├── README.md
│   │   ├── hooks.json
│   │   ├── pre_tool_use.sh
│   │   └── post_tool_use.sh
│   └── codex-mcp/                  # Phase 2 portability confirmation
├── contracts/
│   └── verified-publication/
│       ├── contract.json
│       └── README.md
├── experiments/
│   └── activation-v1/
│       ├── README.md
│       ├── protocol.json
│       ├── corpus/
│       │   ├── authoring/
│       │   ├── validation/
│       │   └── held-out/
│       ├── labels/authoring/
│       ├── manifests/              # digests only for withheld labels
│       ├── runner/
│       └── reports/
├── spec/
│   ├── protocol.schema.json
│   ├── contract.schema.json
│   ├── opportunity.schema.json
│   ├── evidence.schema.json
│   ├── decision.schema.json
│   └── receipt.schema.json
├── testdata/
│   ├── artifacts/
│   ├── paths/
│   └── protocol/
└── docs/
    ├── decisions/
    ├── plans/
    │   └── 2026-08-17-phoenix-implementation-plan.md
    ├── security-model.md
    └── experiment-results.md
```

Do not create deferred directories until their phase begins. The tree documents intended ownership; it is not permission to scaffold the whole platform up front.

---

# Phase 0 — Research contract and evidence baseline

## Phase 0 objective

Pre-commit what Phoenix must prove before writing the control plane. The primary deliverable is a trustworthy activation corpus and experiment protocol, not application code.

## Phase 0 gate

Phase 0 is accepted only when:

- the protected effect and opportunity boundary are unambiguous;
- one Phase 1 external runtime is selected through a reproducible spike and its integration surface verified;
- the runtime process is denied direct write permission to publication roots, while the Phoenix gateway runs under the only OS principal permitted to publish there;
- authoring, validation, and final held-out corpora are separated by generating family/source trace;
- labels have rationales and independent review;
- activation, false-block, latency, and operating-cost thresholds are committed before results exist;
- Arm D remains in the experiment;
- prior-art adoption decisions are recorded;
- the experiment can conclude that only the gateway/receipts are useful and kill the activation superstructure.

### Task 0.1: Record the system boundary

**Files:**
- Create: `docs/decisions/0001-system-boundary.md`

**Steps:**
1. State that Phoenix does not call models or orchestrate tasks.
2. Define advisory, cooperative, and mediated enforcement.
3. Require adapters to report known bypass paths.
4. Define the Phase 1 protected effect as publication from a Phoenix-owned staging root into explicit allowlisted publication roots. Default-deny every destination outside those roots.
5. State that each publication root has explicit overwrite, mode, ownership, quota, and file-type policy.
6. State that the runtime principal cannot write publication roots directly and that the gateway—not the agent—owns the final publication operation.
7. Review against the non-negotiable boundary list.
8. Commit: `docs: define Phoenix system boundary`.

### Task 0.2: Verify and select the first external runtime

**Files:**
- Create: `docs/decisions/0002-phase-1-runtime.md`

**Default selection:** Claude Code through its documented pre-tool and post-tool hook surfaces, subject to the reproducible spike below.

**Files:**
- Create: `experiments/runtime-spike/claude-code/README.md`
- Create: `experiments/runtime-spike/claude-code/fixtures/`
- Create: `experiments/runtime-spike/claude-code/run.sh`
- Create: `experiments/runtime-spike/claude-code/results.json`

**Steps:**
1. Pin the exact Claude Code version, model, configuration, permissions, hook configuration, invocation command, timeout, retry policy, randomness controls where available, and cost-accounting method.
2. Capture official pre-tool and post-tool payload fixtures with versioned documentation references.
3. Run the spike against ordinary writes, shell redirection, copy, rename, symlink, and unsupported/alternate write paths.
4. Record actual attempts from an independent experiment observer, adapter-observed attempts, block/allow behavior, and known bypasses.
5. Label the adapter cooperative, never mediated.
6. Define how it sends raw task context, proposed operation, tool identity, arguments, target path, and runtime version to Phoenix.
7. Verify OS separation: the runtime principal may write only the staging root; only the Phoenix gateway principal may write allowlisted publication roots.
8. Record OpenAI Codex through a Phoenix MCP tool as the planned Phase 2 portability adapter; do not implement it in Phase 1.
9. Commit: `experiment: verify Phase 1 runtime boundary`.

### Task 0.3: Select prior-art boundaries

**Files:**
- Create: `docs/decisions/0003-prior-art.md`

**Steps:**
1. Compare a native in-toto Statement envelope with a minimal custom evidence envelope.
2. Prefer in-toto-compatible subject, predicate, digest, and producer fields unless the experiment proves an incompatibility.
3. Compare server-side one-use opaque grants with macaroons/Biscuit for the local filesystem gateway.
4. Prefer one-use server-side grants for Phase 1 if attenuation and offline delegation are not required; explicitly avoid inventing cryptographic token semantics.
5. Compare a typed Go decision table with Cedar/OPA.
6. Prefer typed code for one Phase 1 contract; revisit Cedar/OPA only after multiple policies create demonstrated complexity.
7. Record what code each adopted standard deletes.
8. Commit: `docs: choose Phase 1 evidence and authorization primitives`.

### Task 0.4: Freeze the minimal contract schema

**Files:**
- Create: `spec/contract.schema.json`
- Create: `contracts/verified-publication/contract.json`
- Create: `contracts/verified-publication/README.md`
- Test later: `internal/contract/registry_test.go`

**Required Phase 1 fields only:**

```yaml
id:
version:
effect_class:
activation:
  positive_signals:
  exclusions:
required_evidence:
decision:
  allow:
  deny:
outcome:
  success:
  failure:
```

**Steps:**
1. Define `filesystem.publish` as the only effect class.
2. Define staged artifacts plus requests targeting explicit allowlisted publication roots as gateway-observed opportunities.
3. Define the Phoenix-owned staging root, allowlisted publication roots, and default-deny behavior outside both.
4. Define exclusions for controlled test fixtures and non-publication staging operations; exclusions never grant destination write authority.
5. Require an artifact digest and declared verification receipt.
6. Define deny reasons as stable machine-readable codes.
7. Define the exact predicate grammar; do not use natural-language predicates.
8. Define RFC 8785 canonicalization, unknown-field rejection, decision precedence, `unknown` semantics, and the stable reason-code registry.
9. Exclude dependencies, composition, amendments, promotion channels, delayed evaluators, and arbitrary executable hooks.
10. Validate the example contract and canonical digest with a standalone schema fixture.
11. Commit: `spec: define minimal verified publication contract`.

### Task 0.5: Define context provenance

**Files:**
- Create: `spec/opportunity.schema.json`
- Create: `docs/decisions/0004-context-authority.md`

**Steps:**
1. Require every context field to declare `observed`, `declared`, or `inferred` provenance.
2. Record producer, timestamp, runtime version, and optional confidence.
3. Prevent inferred fields from independently satisfying security-critical preconditions.
4. Define gateway-observed, adapter-observed, independently observed, and unobservable-bypass classes.
5. Define the opportunity event before contract lookup on each visible path.
6. State that `opportunity_without_activation` is a valid live claim only for a gateway-observed request; end-to-end detector misses require the independent experiment observer.
7. Defer cadence monitoring from Phase 1.
8. Define branch-specific protocol transitions and terminal states for rejected, denied, failed, indeterminate, completed, and reconciled runs before any Go ledger implementation.
9. Commit: `spec: define opportunity, provenance, and run state semantics`.

### Task 0.6: Build the activation corpus

**Files:**
- Create: `experiments/activation-v1/README.md`
- Create cases under:
  - `experiments/activation-v1/corpus/authoring/`
  - `experiments/activation-v1/corpus/validation/`
  - `experiments/activation-v1/corpus/held-out/`
- Create authoring labels under: `experiments/activation-v1/labels/authoring/`
- Store validation and final held-out labels outside the implementation workspace; commit only their digests under `experiments/activation-v1/manifests/` until the relevant evaluation gate.

**Runnable input shape (contains no answers):**

```yaml
case_id:
task:
observed_context:
proposed_effect:
family_id:
source_trace_id:
```

**Withheld label shape (stored outside the implementation workspace until evaluation):**

```yaml
case_id:
should_activate:
required_evidence:
expected_decision:
label_rationale:
```

Before implementation, commit a manifest containing input digests and withheld-label digests. Split by generating family and source trace, never by random individual case.

**Required case classes:**

- obvious publish attempts;
- ordinary writes that are not publication;
- ambiguous paths and naming;
- symlink and traversal attempts;
- shell redirection, rename, copy, archive, and generated-file paths;
- alternate tool surfaces;
- missing context;
- stale context;
- malicious text attempting to suppress activation;
- legitimate exemptions;
- evidence present but bound to the wrong artifact;
- retries and duplicate attempts;
- correct refusal;
- TradePilot-inspired unreachable-gate and zero-activity cases, generalized without coupling Phoenix to TradePilot.

**Steps:**
1. Write authoring cases first, then a sealed validation tranche, then a final held-out tranche.
2. Declare whether each adversarial family belongs to authoring, validation, or final held-out evaluation.
3. Have an independent evaluator label validation and held-out cases without seeing the implementation.
4. Preserve disagreements as explicit adjudication records.
5. Keep held-out labels outside the implementation workspace or access-control them until the final evaluation command.
6. Commit input and label digests before implementation; prohibit tuning after a tranche is opened.
7. Add a corpus manifest with counts, provenance classes, family/source coverage, and explicit denominators—not outcome results.
8. Commit authoring inputs and sealed manifests separately.

### Task 0.7: Pre-commit experiment arms and decision rules

**Files:**
- Create: `experiments/activation-v1/protocol.json`
- Create: `docs/decisions/0005-go-no-go-rules.md`

**Arms:**

All arms use the same independently logging publication boundary, OS permissions, tools, model/runtime version, context, retry accounting, and case-order randomization.

- **A:** base agent through a permissive publication gateway;
- **B:** same base agent and permissive gateway plus the frozen best reasonable static skill;
- **C:** Phoenix activation, evidence, receipts, and enforced publication gateway;
- **D:** identical Phoenix activation, evidence, and receipts with the same gateway configured permissively.

Arm C minus D identifies enforcement value. Other component claims require precommitted diagnostic ablations or must remain explicitly non-causal.

**Diagnostic:** Oracle activation may be used to distinguish detector failures from policy/gateway failures. It is not a product arm.

**Primary measures:**

- activation recall;
- activation precision;
- missed opportunity rate;
- false block rate;
- correct abstention rate;
- unsupported publication rate;
- evidence mismatch detection;
- human correction burden;
- latency and total operating cost per successful task.

**Steps:**
1. Create a falsifiable claim table. For every claim record population/task class, intervention, comparator, unit of analysis, primary endpoint, numeric threshold, uncertainty bound, missing/indeterminate handling, retry clustering, exclusions, stopping rule, multiplicity policy, minimum detectable effect, and forced failure verdict.
2. Define unacceptable regression outside positive cases.
3. Define blinded grading and negative controls.
4. Freeze exact runtime/model/configuration, tools, permissions, invocation commands, repetitions, timeout, case randomization, retry accounting, and cost accounting for every arm.
5. Define sentinel/canary instrumentation checks and require every run to close with actual-attempt, observed-opportunity, activation, and labeled-positive counts with explicit denominators. Recall is undefined for zero-positive subsets.
6. Define component-level decisions:
   - activation fails but gateway wins → narrow Phoenix to gateway plus receipts;
   - activation wins but enforcement adds no value → retain only if advisory mode beats the static skill after cost;
   - detector does not beat the static skill after cost → kill detector even if the gateway survives;
   - enforced C minus permissive D shows no benefit or unacceptable completion loss → remain shadow/advisory;
   - gateway loses to a simpler independently logged allow/deny wrapper → kill the gateway architecture;
   - mandatory evidence does not reduce unsupported publication beyond gateway-only logging → kill mandatory evidence;
   - run receipts do not improve blinded causal reconstruction or correction time over structured logs → kill the receipt product layer;
   - false blocking exceeds threshold → stop and repair activation/conflicts;
   - result indeterminate → expand or repair corpus, do not promote architecture.
7. Scope every conclusion to verified filesystem publication under the tested runtime; do not generalize to behavioral contracts or agent control overall.
8. Commit: `docs: precommit Phoenix experiment decisions`.

### Task 0.8: Phase 0 independent review

**Files:**
- Create: `docs/decisions/0006-phase-0-review.md`

**Steps:**
1. An evaluation reviewer independent of implementation reviews corpus validity, leakage, labels, and comparison design.
2. A security/implementation reviewer validates adapter feasibility, OS-principal separation, effect-boundary completeness, and minimum schema.
3. A human-factors reviewer checks whether the proposed receipt can state the behavioral deal without configuration sludge.
4. The accountable project chair records accepted findings and unresolved blockers.
5. Do not begin Phase 1 while the corpus or go/no-go rules remain unaccepted.
6. Commit: `docs: accept Phoenix Phase 0 research contract`.

---

# Phase 1 — Thin activation and mediated-publication experiment

## Phase 1 objective

Test whether Phoenix can recognize protected publication opportunities under an independently observed experiment, expose the limits of live observation, mediate one Linux filesystem effect through OS-principal separation, and provide reconstructable receipts through one real external runtime. Claims are limited to verified filesystem publication under the pinned runtime.

## Explicit omissions

Do not build in Phase 1:

- WASM execution;
- a second runtime adapter;
- contract dependencies or composition;
- amendments or automatic generation;
- shadow/canary/stable promotion;
- generalized rollback infrastructure;
- cross-agent state;
- network/API effects;
- a dashboard;
- a general policy language;
- hosted services.

### Task 1.1: Initialize the minimal Go project

**Files:**
- Create: `go.mod`
- Create: `cmd/phoenix/main.go`
- Create: `cmd/phoenix/main_test.go`
- Create: `Makefile`
- Create: `.gitignore`

**Steps:**
1. Pin the supported Go toolchain, pure-Go SQLite library, JSON Schema library, RFC 8785 canonicalization library, and static-analysis versions. Record the CGO-disabled policy.
2. Write a failing CLI smoke test in `cmd/phoenix/main_test.go` for `phoenix version` and the single transport selected by the runtime spike.
3. Run the test and verify failure.
4. Implement only the command shell required by the test.
5. Add pinned `make test`, `make lint`, `make build`, and `make validate-spec` targets with documented expected exit status.
6. Run tests, static analysis, schema validation, and a CGO-disabled build.
7. Commit: `build: initialize Phoenix CLI`.

### Task 1.2: Implement immutable contract loading

**Files:**
- Create: `internal/contract/model.go`
- Create: `internal/contract/validate.go`
- Create: `internal/contract/registry.go`
- Create: `internal/contract/registry_test.go`

**Steps:**
1. Test valid contract loading and stable content digest.
2. Test missing required fields, rejected unknown fields, malformed predicates, and unsupported schema versions. Retirement is deferred to Phase 3.
3. Test that published contract content cannot be mutated in place.
4. Implement schema validation and content-addressed identity.
5. Load only the verified-publication contract.
6. Run focused and full tests.
7. Commit: `feat: load immutable Phoenix contracts`.

### Task 1.3: Implement the append-only ledger

**Files:**
- Create: `internal/admission/admission.go`
- Create: `internal/admission/admission_test.go`
- Create: `internal/ledger/ledger.go`
- Create: `internal/ledger/sqlite.go`
- Create: `internal/ledger/migrations.go`
- Create: `internal/ledger/sqlite_test.go`

**Required event order:**

```text
opportunity.observed
contract.considered
contract.activated | contract.rejected
policy.decided
evidence.recorded
effect.authorized | effect.denied
effect.completed
outcome.evaluated
run.closed
```

This is an event vocabulary, not a single mandatory linear order. Phase 0 must define branch-specific transitions and terminal states for activated, rejected, denied, failed, indeterminate, and recovered runs.

Publication durability state machine:

```text
prepared -> attempted -> completed
                   \-> indeterminate -> reconciled_completed | reconciled_absent
```

Phoenix does not claim exactly-once filesystem publication across SQLite and the filesystem.

**Steps:**
1. Define exact transition tables, terminal states, stable reason codes, and recovery semantics before writing migrations.
2. Write admission tests for schema allowlisting, unknown-field rejection, size limits, normalization, and secret redaction before persistence.
3. Require every transport and internal caller to pass through the same admission API; the ledger accepts only admitted envelopes.
4. Test append-only ordering, idempotency keys, correlation IDs, and transaction rollback.
5. Test prepared/attempted/completed/indeterminate crash points and idempotent reconciliation after restart.
6. Test that retries cannot duplicate a completed publication and that ambiguity is reported rather than guessed.
7. Test concurrent writes under SQLite WAL.
8. Test reconstruction of rejected, denied, failed, indeterminate, and successful trails.
9. Implement the smallest schema satisfying those tests.
10. Do not implement distributed event sourcing or mutable global scratch state.
11. Commit: `feat: admit and record recoverable Phoenix events`.

### Task 1.4: Implement opportunity observation before activation

**Files:**
- Create: `internal/opportunity/model.go`
- Create: `internal/opportunity/monitor.go`
- Create: `internal/activation/predicates.go`
- Create: `internal/activation/engine.go`
- Create: `internal/activation/engine_test.go`
- Create: `experiments/activation-v1/runner/main.go`
- Create: `experiments/activation-v1/runner/main_test.go`

**Steps:**
1. Test gateway-observed, adapter-observed, independently observed, and unobservable-bypass cases separately.
2. Test deterministic activation for identical context and registry digest.
3. Test exclusions and adversarial suppression text.
4. Test that missing required facts become `indeterminate` or `deny`, never implicit allow.
5. Test `opportunity_without_activation` only for gateway-observed requests; use the independent observer for end-to-end detector misses.
6. Test broken instrumentation with sentinel attempts.
7. Implement reason-coded activation and rejection decisions.
8. Implement a harness-neutral corpus runner that invokes Phoenix core packages directly and receives actual-attempt facts from experiment fixtures; it must not require Claude Code, Hermes, or any agent framework.
9. Run the authoring corpus only; keep validation and final held-out tranches sealed.
10. Commit: `feat: detect publication opportunities and activation gaps`.

### Gate 1A: Activation-only stop gate

Before evidence, grants, gateway internals, adapters, or receipts expand further:

1. Freeze the detector.
2. Open only the sealed validation tranche, never the final held-out tranche.
3. Run the harness-neutral corpus runner with an independent actual-attempt denominator.
4. Report actual attempts, observed opportunities, activations, labeled positives, indeterminate cases, and explicit denominators.
5. Stop Phase 1 if the precommitted activation threshold fails. Repair the corpus or detector without opening the final tranche.
6. Continue only after the independent evaluation reviewer accepts the activation result.

### Task 1.5: Implement evidence statements and artifact storage

**Files:**
- Create: `spec/evidence.schema.json`
- Create: `internal/evidence/statement.go`
- Create: `internal/evidence/verify.go`
- Create: `internal/evidence/verify_test.go`
- Create: `internal/artifact/store.go`
- Create: `internal/artifact/store_test.go`

**Steps:**
1. Test content-addressed artifact storage and mutation detection.
2. Test that evidence binds to run, contract digest, source artifact digest, and proposed action.
3. Test that evidence from another run or artifact is rejected.
4. Test that model narration alone cannot satisfy a required evidence type.
5. Define trusted producer classes and verification methods. Self-attested evidence is labeled `declared` and cannot independently authorize a mediated effect.
6. Treat only Phoenix-generated artifact digests and outputs from explicitly allowlisted verifier executables, identified by digest and executed by Phoenix under a constrained principal, as authoritative Phase 1 evidence.
7. Define artifact-store and database ownership and the attacker modification model.
8. Implement the selected in-toto-compatible statement fields.
9. Preserve failed evidence and retries rather than overwriting them.
10. Commit: `feat: record authority-graded evidence`.

### Task 1.6: Implement one-use action grants

**Files:**
- Create: `internal/token/grant.go`
- Create: `internal/token/grant_test.go`

**Steps:**
1. Define principal authentication for the single selected transport: Unix peer credentials and socket ownership/permissions if Unix sockets are selected, or inherited parent-process identity and a private pipe if stdio is selected.
2. Test grants bound to authenticated principal, run, contract digest, artifact digest, target path, and expiry.
3. Test one-use consumption.
4. Test replay, substitution, stale grant, unauthenticated caller, and wrong-principal rejection.
5. Implement server-side opaque grants; do not expose custom bearer claims.
6. Ensure denied or consumed grants remain auditable.
7. Commit: `feat: issue authenticated one-use publication grants`.

### Task 1.7: Implement the mediated filesystem gateway

**Files:**
- Create: `internal/gateway/filesystem/path.go`
- Create: `internal/gateway/filesystem/gateway.go`
- Create: `internal/gateway/filesystem/gateway_test.go`

**Steps:**
1. Freeze Phase 1 to Linux and document required kernel/filesystem features.
2. Build tests entirely inside temporary staging and publication roots; never use a real project tree.
3. Use directory-file-descriptor-relative resolution with Linux `openat2` beneath/no-symlink constraints or a documented equivalent; authorize opened descriptors, not unchecked path strings.
4. Permit regular files only. Reject devices, FIFOs, sockets, directories, mount/bind-mount escapes, and unsupported cross-filesystem publication.
5. Revalidate the staged source descriptor, digest, destination root, and authenticated grant immediately before mutation.
6. Publish through a same-directory temporary file, apply explicit mode/ownership policy, fsync file and directory, then atomic rename according to a per-root no-clobber/overwrite policy.
7. Test missing evidence, digest mismatch, traversal, symlink escape, parent-directory swap, source replacement, mount escape, special files, rename/copy substitution, replay, expiry, unrelated target changes, and concurrent races.
8. Test crash points before write, after temporary write, after file fsync, after rename, before directory fsync, and before ledger completion; implement cleanup and idempotent reconciliation.
9. Test fail-closed conflict behavior and compact reason codes.
10. Test that denied actions produce zero protected side effects.
11. State the exact guarantee: at-most-one accepted final artifact after reconciliation, with explicit `indeterminate` state across unresolved crash windows—not exactly-once publication.
12. Commit: `feat: mediate crash-recoverable Linux publication`.

### Task 1.8: Implement the single transport selected by the runtime spike

**Files:**
- Create: `spec/protocol.schema.json`
- Create: `spec/decision.schema.json`
- Create: `internal/server/service.go`
- Create one of: `internal/server/stdio.go` or `internal/server/unix.go`
- Create: `internal/server/service_test.go`

**Operations:**

```text
ObserveOpportunity
ResolveContracts
SubmitEvidence
RequestAuthorization
PublishArtifact
CloseRun
GetReceipt
```

**Steps:**
1. Define exact JSON Schemas, RFC 8785 canonicalization, a deliberately small predicate grammar, decision precedence, unknown-field rejection, and a stable reason-code registry.
2. Test canonical request/response envelopes and schema versions.
3. Test malformed requests, oversized fields, secrets, duplicate idempotency keys, unsupported versions, and timeouts at the common admission boundary.
4. Test peer/parent identity according to the selected transport.
5. Test that protocol calls cannot bypass admission or ledger events.
6. Implement only the transport proven necessary by the runtime spike. Defer the second transport and localhost HTTP.
7. Commit: `feat: expose one authenticated Phoenix transport`.

### Task 1.9: Build the Claude Code adapter

**Files:**
- Create: `adapters/claude-code/README.md`
- Create: `adapters/claude-code/hooks.json`
- Create: `adapters/claude-code/pre_tool_use.sh`
- Create: `adapters/claude-code/post_tool_use.sh`
- Create: `internal/adapter/conformance_test.go`

**Steps:**
1. Test adapter payloads against captured official hook fixtures.
2. Emit raw runtime facts with provenance; do not pre-label inferred facts as observed.
3. Route candidate publication attempts to Phoenix before the effect.
4. Run the adapter as an OS principal that cannot write publication roots directly.
5. Record unsupported or unobservable tool paths as known bypasses and measure them through the independent experiment observer.
6. Ensure the adapter cannot claim mediated enforcement; only the Phoenix filesystem gateway may do so.
7. Add a conformance test for opportunity → decision → evidence → authorization → receipt.
8. Document setup commands without making Phoenix depend on Claude Code.
9. Commit: `feat: add Claude Code cooperative adapter`.

### Task 1.10: Render one run receipt

**Files:**
- Create: `spec/receipt.schema.json`
- Create: `internal/receipt/render.go`
- Create: `internal/receipt/render_test.go`

**Required receipt fields:**

- intended effect;
- contract and digest;
- activation decision and exact reasons;
- context provenance;
- required and supplied evidence;
- allow/deny decision;
- enforcement strength and known bypasses;
- resulting artifact digest and target;
- outcome;
- observation class, activation gap, instrumentation health, and known unobservable paths.

**Steps:**
1. Test a successful, denied, indeterminate, gateway-missed-activation, detector-missed-activation, and unobservable-bypass receipt.
2. Confirm persisted inputs were already normalized and redacted by the admission boundary; rendering must not be the first redaction layer.
3. Render JSON plus compact human-readable text.
4. Do not build a dashboard or general contract library UI.
5. Commit: `feat: render causal Phoenix run receipts`.

### Task 1.11: Implement and run the controlled experiment

**Files:**
- Create: `experiments/activation-v1/runner/`
- Create: `experiments/activation-v1/reports/phase-1-results.json`
- Create: `docs/experiment-results.md`

**Steps:**
1. Freeze implementation and commit the final code digest before opening held-out labels.
2. Run Arms A, B, C, and D through the same independently logging gateway, pinned runtime/model/configuration, permissions, tools, case ordering, retry policy, and cost accounting.
3. Run oracle activation only as a diagnostic.
4. Record actual attempts, gateway-observed attempts, adapter-observed attempts, activations, bypasses, latency, false blocks, missed opportunities, explicit denominators, and operating cost.
5. Blind grade outputs where judgment is required.
6. Compare results with pre-committed thresholds without changing thresholds.
7. Produce component-level verdicts for activation, evidence, gateway, and receipts.
8. Do not declare Phoenix successful from enforcement tests alone.
9. Commit raw machine-readable results and the interpreted report separately.

### Task 1.12: Phase 1 gate and independent review

**Files:**
- Create: `docs/decisions/0007-phase-1-verdict.md`

**Required acceptance checks:**

1. The independent observer supplies the actual-attempt denominator; Phoenix never claims visibility into unobserved bypasses.
2. End-to-end and conditional-on-observation activation recall meet their separately pre-committed thresholds.
3. False activation and false blocking remain within limits.
4. Adversarial near-misses are handled correctly.
5. Missing or mismatched evidence blocks Arm C.
6. Arm D records the same condition without blocking it.
7. Replay and action substitution fail at the gateway.
8. Every result is reconstructable from the ledger.
9. Latency, correction burden, bypass rate, instrumentation health, and conflict frequency are reported.
10. Each run closes with actual-attempt, observed-opportunity, activation, labeled-positive, and indeterminate counts; zero-positive recall is reported as undefined.
11. Crash recovery produces completed, reconciled-absent, or explicit indeterminate state without an exactly-once claim.
12. The runtime principal cannot write publication roots without the Phoenix gateway principal.

**Steps:**
1. An evaluation reviewer independent of implementation reviews causal claims and corpus leakage.
2. A security/code reviewer validates OS separation, evidence authority, filesystem races, crash recovery, and side effects.
3. A human-factors reviewer tests whether receipts improve blinded causal reconstruction and correction time over structured logs.
4. The accountable project chair selects exactly one verdict:
   - proceed to portability;
   - narrow to gateway plus receipts;
   - remain advisory only;
   - repair experiment;
   - kill Phoenix for this task class.
5. Push only the accepted result state.

---

# Phase 2 — Portability confirmation and contract composition

Phase 2 begins only if Phase 1 demonstrates useful activation or a narrower useful gateway/receipt product.

## Phase 2 objectives

- prove the same contract semantics across a second unrelated runtime;
- add OpenAI Codex through a Phoenix MCP adapter as the planned second integration;
- introduce contract conflicts without silently widening permissions;
- investigate network/API action canonicalization without shipping a universal gateway.

## Phase 2 tasks

### Task 2.1: Freeze the adapter conformance contract

Create `spec/adapter-conformance.schema.json` and require every adapter to report supported hooks, intercepted effect classes, evidence fidelity, runtime version, clock assumptions, and bypass paths.

### Task 2.2: Implement the Codex MCP adapter

Create `adapters/codex-mcp/` and prove that the unchanged verified-publication contract produces equivalent activation reasons, evidence requirements, and gateway decisions. It may be cooperative; do not claim mediated enforcement outside Phoenix’s gateway.

### Task 2.3: Add differential portability tests

Run golden protocol fixtures through both adapters and fail on semantic differences, missing opportunity events, altered reason codes, or evidence-shape drift.

### Task 2.4: Introduce minimal composition

Add only:

- explicit contract dependencies;
- explicit conflicts;
- authority precedence;
- fail-closed empty permission intersections;
- a human-visible conflict decision.

Measure conflict and override frequency. Do not auto-merge contracts.

### Task 2.5: Research network action identity

Write `docs/decisions/0008-network-action-canonicalization.md` covering method, URL normalization, headers, body canonicalization, redirects, retries, idempotency, authentication references, and mutable remote state. Do not reuse the filesystem digest blindly.

### Phase 2 gate

Proceed only if:

- both adapters preserve contract meaning;
- portability cost is acceptable;
- composition does not create excessive deadlocks;
- the control plane remains independent of either runtime;
- network action identity has an explicit design rather than a file-shaped digest.

---

# Phase 3 — Evidence-gated amendments

Phase 3 begins only after real operational corrections exist. Synthetic corrections are insufficient justification.

## Phase 3 objectives

- convert repeated observed failures or explicit human corrections into small candidate amendments;
- prevent active contracts from rewriting themselves;
- replay candidates against historical and held-out cases;
- keep review legible enough to resist rubber-stamping.

## Lifecycle

```text
observed failure or correction
→ structured amendment proposal
→ static validation
→ historical replay
→ held-out regression check
→ shadow evaluation
→ bounded canary
→ human approval
→ stable pointer update or rejection
```

## Required implementation

- immutable versions and content digests;
- channels: `draft`, `candidate`, `shadow`, `canary`, `stable`, `retired`;
- replay runner;
- permission-expansion detection;
- stronger approval for expanded authority;
- small typed diffs limited to activation, exclusion, evidence, or decision fields;
- one-click rejection and deterministic rollback of the active pointer;
- unchanged historical receipts pinned to original digests;
- amendment-volume and reviewer-rubber-stamping metrics.

## Phase 3 kill conditions

- amendments regularly require large or opaque diffs;
- reviewers accept without meaningful inspection;
- candidate optimization overfits known traces;
- evaluator changes rewrite historical meaning;
- correction burden exceeds the repeated failure burden;
- autonomous promotion becomes necessary to keep the system usable.

---

# Phase 4 — Multi-agent alliance capabilities

Phase 4 begins only after single-agent activation, mediation, and amendments have independent evidence.

## Objectives

- carry contracts through agent handoffs;
- maintain scoped shared state without a global scratchpad;
- preserve one accountable mission owner;
- evaluate capabilities by task class rather than global agent scores.

## Required objects

- stable agent principal;
- handoff envelope containing intent, locked contracts, allowed/forbidden effects, required return artifact, and open approvals;
- scoped state record with writer, version, causal parent, visibility, classification, and optional expiry;
- lease and idempotency semantics for exclusive work;
- per-capability and per-task-class outcomes.

## Invariants

- receiving agents cannot silently amend locked contracts;
- composition cannot widen permissions;
- agents cannot grade or promote their own amendments;
- shared state never treats hidden model memory as authoritative;
- a subordinate success claim requires independently inspectable artifacts;
- routine disagreement stays internal; value conflicts return as options, consequences, and recommendation.

## Phase 4 gate

The alliance must outperform one capable agent after coordination cost on representative missions. More messages, delegates, or agent identities do not count as improvement.

---

# Verification strategy

## Every code task

1. Write a focused failing test.
2. Run it and confirm the expected failure.
3. Implement the smallest behavior.
4. Run the focused test.
5. Run the full suite.
6. Run static analysis and build.
7. Inspect side effects in a temporary environment.
8. Commit only the verified increment.

## Security-specific checks

- no real project paths in filesystem gateway tests;
- no ambient credentials in contracts or evidence;
- secret redaction occurs before persistence;
- symlink and path traversal tests run on supported platforms;
- denied actions leave zero protected side effects;
- grants are exact-action, expiring, and one-use;
- failed attempts remain visible;
- adapters cannot overstate enforcement strength.

## Neighboring behavior

A narrow passing test is insufficient. Each phase must verify:

- ordinary scratch writes still work;
- legitimate excluded writes are not blocked;
- duplicate/retried requests remain idempotent;
- ledger corruption or unavailable policy fails safely;
- Phoenix restart preserves reconstructable state;
- malformed adapter input cannot authorize effects;
- disabling enforcement changes only enforcement, not activation or evidence observation.

---

# Sequencing and review ownership

Phoenix does not displace the currently prioritized research-kernel milestone. It follows that work unless Raoul explicitly changes priority.

TradePilot and the research kernel may contribute generalized failure traces and architectural evidence, but Phoenix must not import TradePilot code, schema, state, or trading assumptions.

Required roles are implementation-independent:

- **Accountable project chair:** scope, gates, orchestration, acceptance, and verification that claimed artifacts exist.
- **Evaluation reviewer:** corpus, leakage controls, evidence standards, statistics, and causal claims; must not implement the detector being graded.
- **Security/implementation owner:** TDD, OS-principal boundary, gateway, adapters, durability, and build quality; cannot approve their own security gate.
- **Human-factors reviewer:** receipt comprehension, correction burden, and review fatigue.
- **Artifact designer:** presentation only after a text receipt proves incremental value over structured logs.

The repository and experiment must remain reproducible without these particular reviewers, Hermes profiles, or any named agent roster.

Production deployment, hosted services, public release, and any external credential gateway require separate approval.

---

# Final delivery gates

Phoenix may call itself successful only when all relevant gates are true:

1. Every Phoenix-visible path records its observation class before contract selection, while unobservable paths remain explicitly disclosed.
2. End-to-end missed activation is measured only against an independent actual-attempt denominator; live claims remain bounded by the observation surface.
3. Held-out activation results beat a strong static skill for the selected task class after cost.
4. Arm D demonstrates the causal value of enforcement separately from observation.
5. Mediated effects cannot bypass the controlled gateway.
6. Evidence binds claims to exact runs, artifacts, decisions, and contract versions.
7. A second unrelated runtime preserves contract semantics in Phase 2.
8. Amendments improve held-out outcomes without silent self-promotion in Phase 3.
9. Human review remains legible and bounded.
10. Multi-agent coordination beats a capable single agent after coordination cost before Phase 4 is retained.

If only the mediated gateway proves useful, Phoenix should become exactly that; receipts survive only if they improve blinded reconstruction or correction time over structured logs. If the detector does not beat a static skill after cost, kill the detector even if the gateway survives. The architecture serves evidence; evidence does not serve the architecture.
