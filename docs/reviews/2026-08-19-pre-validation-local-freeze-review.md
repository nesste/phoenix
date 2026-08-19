# Pre-validation local-artifact freeze independent review

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-19
- **Verdict:** `REVISE`
- **Candidate commit:** `239fc2bf9e42dfdddbba067d3407005d1addd973`
- **Candidate subject:** `docs: prepare local pre-validation freeze candidate`

## Verdict

Nine of the ten required checks hold. The candidate is review-only, the four local artifacts are the right freeze boundary, prompts and runtime identity agree, Arm A schemas match `flat.go`, production and authoring worlds share one canonical digest, the committed Linux-amd64 world-build reproduced under `make quality`, the grader digest recomputes and is pinned by every authoring label, and protected freeze bytes are unchanged.

One freeze-identity defect blocks acceptance. The Phoenix invocation template records the committed Linux-amd64 world-build digest as a fixed `--world-build` argument, but `prepareRuntimeInvocation` passes the digest of the host-local `buildPhoenix` executable (`-X main.version=authoring`, host `GOOS`/`GOARCH`). Accepting the candidate as written would freeze an invocation claim that contradicts the already-frozen runner.

This review does not authorize copying these artifacts into `pre-validation-artifacts.json`, corpus import, validation schedule generation, Gate 1A, validation execution, or held-out execution.

## Findings

### [P1] Phoenix template treats a variable world-build digest as fixed bytes — `experiments/frontier-v1/artifacts/pre-validation-local-candidate.json:63`

The review prompt requires the Claude and Phoenix argument templates to describe `prepareRuntimeInvocation`, with variable placeholders distinguished from fixed bytes. The Phoenix template lists `--world-build` followed by the committed Linux-amd64 digest `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d`. That digest is the correct identity of the separate world-build artifact. It is not the value `prepareRuntimeInvocation` writes.

`prepareRuntimeInvocation` copies `request.WorldBuild` into the MCP server args. That field is the digest returned by `buildPhoenix`, which rebuilds Phoenix with `-X main.version=authoring` and `runtime.GOOS`/`runtime.GOARCH`. The frozen `make quality` identity uses `version=dev`, `GOOS=linux`, and `GOARCH=amd64`. Those two digests are not the same object. The template also names a `<freshly built Phoenix executable>` in the same argv, which makes the hardcoded digest internally inconsistent: a freshly built host binary cannot carry the committed Linux-amd64 world-build identity.

Consequence: freezing this candidate would record that every trial passes `--world-build sha256:27c2f537…`, while the already-frozen scheduled runner would pass a different digest. The Linux-amd64 world-build freeze can stand on its own; it must not be restated as a fixed invocation byte.

Smallest correction: in `phoenix_server_invocation_template` only, replace the hardcoded digest with a placeholder such as `<world-build digest>`. Keep `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` solely under `world_definition_and_world_build_digest`. Do not change the frozen runner.

```56:66:experiments/frontier-v1/artifacts/pre-validation-local-candidate.json
      "phoenix_server_invocation_template": [
        "<freshly built Phoenix executable>",
        "serve", "--stdio",
        "--world", "worlds/dev-repo/world.json",
        "--world-schema", "spec/world.schema.json",
        "--episode-db", "<per-attempt episodes.db>",
        "--root-refs", "<per-attempt roots.json>",
        "--world-build", "sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d",
        "--arm", "<A|B|C|D|E>",
        "[when declared] --state-events", "<per-attempt state-events.json>"
      ],
```

```100:104:experiments/frontier-v1/runner/runtime.go
			"args": []string{
				"serve", "--stdio", "--world", request.WorldPath, "--world-schema", request.SchemaPath,
				"--episode-db", request.EpisodePath, "--root-refs", rootPath, "--world-build", request.WorldBuild,
				"--arm", request.Arm,
			},
```

```465:488:experiments/frontier-v1/runner/main.go
	command := exec.Command("go", "build", "-trimpath", "-buildvcs=false", "-ldflags=-s -w -buildid= -X main.version=authoring", "-o", executable, "./cmd/phoenix")
	command.Dir = root
	command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOTOOLCHAIN=go1.26.6")
	// ...
	manifest, err := buildmanifest.Build(root, buildmanifest.Input{
		// ...
		GOOS:            runtime.GOOS, GOARCH: runtime.GOARCH, CGOEnabled: false,
	})
```

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Candidate boundary | Pass | The package is exactly runtime invocation and A–E prompts, Arm A schemas, world definition/build, and grader. `remaining_after_acceptance` is only `schedule digest`. No validation schedule is included. |
| 2. Closed gates | Pass | Candidate `status` is `review_candidate`, `frozen` is false, both outcome gates are false, and all four nested artifacts are unfrozen. `pre-validation-artifacts.json` remains `partial` with both gates false and all five open entries listed. |
| 3. Runtime identity | Pass | Protocol, candidate, runner constants, `prepareRuntimeInvocation`, and `scheduledRunConfiguration` agree on Claude Code 2.1.229, `claude-sonnet-5`, low effort, standard service tier, 12 turns, 180 s, USD 0.15, stdio, surface-only, no session persistence, `--strict-mcp-config`, empty `--tools`, stream-json, and one pre-token infrastructure retry. |
| 4. Exact prompts | Pass | A and B are exactly `Use only the configured tools. Follow the user request exactly.` C, D, and E append the bootstrap-intent sentence. Flat arms do not receive live handles or that sentence. Arm B alone appends the already-frozen static document. |
| 5. Invocation | Fail | Claude argv, per-attempt paths, allowlists, optional state events, and the Arm B file match `prepareRuntimeInvocation`. `--world-build` is recorded as the Linux-amd64 digest instead of the runner's `request.WorldBuild` variable (P1). |
| 6. Arm A schemas | Pass | Independently derived twelve sorted tools from `worlds/dev-repo/world.json` via `root_name + '_' + verb` and `Invoke <type>.<verb> on <label>.` match the candidate. Every `args_schema` is a complete object with `additionalProperties: false`. Descriptions contain no frontier, refusal, or opaque-handle instruction. |
| 7. World identity | Pass | Production and authoring worlds are byte-identical (`sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143`) and share canonical digest `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b`. Committed Linux-amd64 CGO-disabled manifest is self-consistent, covers executable, verbs, schemas, world definition, authored activation/frontier/refusal rules, and null Phase 1 weights, and reproduced byte-for-byte under `make quality`. |
| 8. Grader identity | Pass | Independently recomputed `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`. Coverage is LF-normalized non-test corpusctl sources, `go.mod`/`go.sum`, and all eight frontier-v1 schemas (17 files). All eight authoring labels pin that digest. |
| 9. Protected artifacts | Pass | Arm B, scheduled-runner file set, and analysis/report file set match their accepted freeze hashes. `pre-validation-artifacts.json` is untouched by the candidate commit. No sealed validation or held-out path was modified; historical v3 sealed files were not opened. |
| 10. Engineering gate | Pass | `make quality` exited 0. `verify-local-freeze-candidate` compared `build/manifest.json` to `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` with no difference. Regenerated world-build digest was `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d`. |

## Independent digest record

LF-normalized UTF-8 SHA-256 unless noted. Working-tree files at `239fc2b` are LF-only, so raw and LF-normalized hashes coincide.

| Artifact | SHA-256 |
| --- | --- |
| `experiments/frontier-v1/artifacts/pre-validation-local-candidate.json` | `sha256:dd3dc41134890836726887d28b39aa766eb5ed789d332941de140f54c7f5ad42` |
| `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` (raw file) | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| World-build internal digest (`json.Marshal` of `build`) | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Production/authoring world canonical JSON | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` |
| Grader digest | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| `internal/surface/flat.go` | `sha256:99ceffc43aa848b4f241a34cc1d70c5e278d5578f376c502e9787f5785b2e8ff` |
| `experiments/frontier-v1/corpusctl/internal/corpus/grade.go` | `sha256:89304bb0ac6b45a8236fd4562b9c2b2212f40ebc80c6bbbaac907f8f233e719d` |
| Arm B document | `sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0` |

Accepted scheduled-runner and analysis file hashes were recomputed and match `pre-validation-artifacts.json` exactly.

## Decision

No P0 defect. One P1 freeze-identity defect in the candidate JSON. The other nine checks, including `make quality` and protected-artifact custody, pass.

Revise only the Phoenix `--world-build` template entry so it is a placeholder, then return the same four-artifact package for confirmation. Do not unfreeze the scheduled runner to make the Linux-amd64 digest appear in live argv. Do not treat this review as authorization to patch `pre-validation-artifacts.json`.

`local_artifact_candidate_test.go` checks prompts, Arm A tool derivation, world canonical equality, and grader digest, but it does not compare the invocation templates to `prepareRuntimeInvocation`. That gap is how the P1 reached review; the focused revision should add that comparison.

## Remaining execution blockers

Unchanged by this review, and still not authorized:

- independent generation and sealing of new validation and held_out families;
- schedule digest;
- Task 0.6 Phase 0 review, if not already accepted on its own record;
- Gate 1A, validation execution, and held-out execution.

Passing the other local artifacts does not open those gates.
