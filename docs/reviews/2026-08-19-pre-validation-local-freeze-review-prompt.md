# Pre-validation local-artifact freeze review prompt

Give this prompt to an evaluation reviewer who did not prepare the candidate package.

```text
Review the outcome-free four-artifact freeze candidate at the identified candidate commit.

Read:

- experiments/frontier-v1/protocol.json, especially runtime, system_prompts, trial, artifact_freeze, and gate
- docs/decisions/0004-go-no-go-rules.md
- docs/decisions/0007-authoring-activation-amendment.md
- experiments/frontier-v1/artifacts/pre-validation-local-candidate.json
- experiments/frontier-v1/artifacts/world-build.linux-amd64.json
- experiments/frontier-v1/pre-validation-artifacts.json
- experiments/frontier-v1/runner/runtime.go
- experiments/frontier-v1/runner/main.go
- experiments/frontier-v1/runner/scheduled_run.go
- internal/surface/flat.go
- worlds/dev-repo/world.json
- experiments/frontier-v1/worlds/authoring.dev_repo.json
- experiments/frontier-v1/corpusctl/internal/corpus/grade.go

Independence and safety rules:

- Do not modify the candidate, protocol, runner, world, schemas, grader, corpus, labels, or accepted freeze manifest.
- Do not generate or execute an authoring, validation, or held-out schedule.
- Do not run a model or inspect any private validation or held-out label.
- The historical v3 sealed files are retired custody evidence. Do not run, relabel, or reseal them.
- Read-only builds, tests, digest recomputation, schema inspection, and source comparison are allowed.

Verify all of the following:

1. Candidate boundary. The package contains exactly the four open local artifacts: runtime invocation and A-E prompts, Arm A schemas, world definition/build, and grader. The validation schedule remains outside the package.
2. Closed gates. The candidate is review-only and not frozen. `pre-validation-artifacts.json` remains `partial`; both outcome gates remain false; all five open entries remain listed.
3. Runtime identity. Runtime version 2.1.229, model `claude-sonnet-5`, low effort, standard service tier, 12 turns, 180-second timeout, USD 0.15 per-trial cap, stdio, surface-only access, no session persistence, strict MCP config, empty built-in tools, stream JSON, and one eligible pre-token infrastructure retry agree across protocol, candidate, runner, and scheduled summary construction.
4. Exact prompts. A and B receive only `Use only the configured tools. Follow the user request exactly.` C, D, and E receive that sentence plus the exact bootstrap-intent sentence. No flat arm receives live handles or the bootstrap instruction. Arm B alone appends the already-frozen static document.
5. Invocation. The Claude and Phoenix argument templates accurately describe `prepareRuntimeInvocation`, including per-attempt paths, arm allowlists, world-build identity, optional state events, and the Arm B document. Variable placeholders are clearly distinguished from fixed bytes.
6. Arm A schemas. The twelve tool names, descriptions, and argument-schema pointers are the deterministic sorted output of `internal/surface/flat.go` over `worlds/dev-repo/world.json`. Each schema is complete, rejects unknown arguments, and exposes no Phoenix frontier, refusal, or opaque-handle instruction.
7. World identity. Production and authoring world documents have the same canonical JSON digest. The committed Linux-amd64, CGO-disabled build manifest reproduces byte-for-byte under `make quality`, carries a self-consistent world-build digest, and covers the executable, registered verbs, schemas, world definition, authored activation/frontier/refusal rules, and null Phase 1 weights.
8. Grader identity. Independently recompute the grader digest. Confirm that its coverage includes LF-normalized non-test corpusctl sources, corpusctl module files, and every frontier-v1 schema, and that authoring labels pin the same digest.
9. Protected artifacts. Confirm that the accepted Arm B, scheduled-runner, and analysis/report bytes and their accepted freeze entries are unchanged. Confirm that no sealed outcome or schedule was opened.
10. Engineering gate. Run `make quality` and report the result. The final manifest comparison must pass.

Return `ACCEPT`, `REVISE`, or `REJECT` with P0-P3 findings and exact paths/lines. If ACCEPT:

- identify the candidate commit;
- record the LF-normalized UTF-8 SHA-256 digest of the candidate JSON;
- record the raw SHA-256 and internal world-build digest of the committed world-build manifest;
- record the grader digest;
- state accepted limitations;
- authorize only a focused acceptance patch that copies the four reviewed artifacts into `pre-validation-artifacts.json`, removes only those four strings from `remaining`, leaves `schedule digest`, keeps status `partial`, and keeps both gates false.

Do not authorize corpus import, validation schedule generation, Gate 1A, validation execution, or held-out execution merely because these four local artifacts pass review.
```
