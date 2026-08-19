# Pre-validation local-artifact freeze focused-revision independent review

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-19
- **Verdict:** `ACCEPT`
- **Revision candidate:** `795ce71acd51beb977190811c90cc538f3c6b928`
- **Candidate subject:** `fix: revise local pre-validation freeze candidate`
- **Base candidate:** `239fc2bf9e42dfdddbba067d3407005d1addd973`
- **First review:** `docs/reviews/2026-08-19-pre-validation-local-freeze-review.md` (`REVISE`)
- **Findings:** None

## Verdict

The focused revision closes the blocking P1. `--world-build` in the Phoenix template is now the placeholder `<world-build digest>`. The committed Linux-amd64 digest remains only on the world-build artifact. `verifyInvocationTemplates` compares both argument templates to `prepareRuntimeInvocation`. Frozen runner implementation bytes, Arm B, analysis/report, protocol, world, grader, and `pre-validation-artifacts.json` are unchanged. `make quality` passed, including byte-for-byte world-build manifest comparison.

No P0, P1, or P2 defect remains in the four local artifacts. This acceptance authorizes only the focused freeze patch described below. It does not authorize corpus import, validation schedule generation, Gate 1A, validation execution, or held-out execution.

## Closed first-review findings

| First-review finding | Status | Evidence |
| --- | --- | --- |
| P1 Phoenix `--world-build` treated a variable digest as fixed bytes | Closed | Candidate line 63 is `"--world-build", "<world-build digest>"`. `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` remains only under `world_definition_and_world_build_digest`. `verifyInvocationTemplates` fails if those identities are conflated. |

The first review also asked for a template-to-`prepareRuntimeInvocation` comparison. That test is present and passed under `make quality`.

## Requirement matrix

No regression relative to the first review except check 5, which now passes.

| Check | Result | Notes |
| --- | --- | --- |
| 1. Candidate boundary | Pass | Unchanged. Four local artifacts; `remaining_after_acceptance` is only `schedule digest`. |
| 2. Closed gates | Pass | Unchanged. Candidate is `review_candidate` and unfrozen. Manifest remains `partial` with both gates false and five open entries. |
| 3. Runtime identity | Pass | Unchanged. |
| 4. Exact prompts | Pass | Unchanged. |
| 5. Invocation | Pass | Claude and Phoenix templates match `prepareRuntimeInvocation` after substituting per-attempt paths. World-build, Arm B document, allowlists, and optional state events are placeholders or documented optional suffixes, distinct from fixed bytes. |
| 6. Arm A schemas | Pass | Unchanged twelve sorted tools. |
| 7. World identity | Pass | Unchanged canonical world digest and self-consistent Linux-amd64 manifest. Reproduced under this revision's `make quality`. |
| 8. Grader identity | Pass | Unchanged `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`. Authoring labels still pin it. |
| 9. Protected artifacts | Pass | Diff `239fc2b..795ce71` touches only the first review record, the candidate JSON, and `local_artifact_candidate_test.go`. Frozen runner implementation files, Arm B, analysis, protocol, world, grader, sealed corpus, and the accepted freeze manifest are untouched. |
| 10. Engineering gate | Pass | `make quality` exited 0. Regenerated world-build digest `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` matched the committed manifest byte-for-byte. |

## Candidate identity

| Artifact | SHA-256 |
| --- | --- |
| Revision candidate commit | `795ce71acd51beb977190811c90cc538f3c6b928` |
| Candidate JSON, LF-normalized UTF-8 | `sha256:19743d85232a2ec83b80bb71fc6b5dc3ee6cc669d319c446218aa248eb8f76b9` |
| World-build manifest, raw file | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| World-build internal digest | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Grader digest | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |

Working-tree candidate JSON is LF-only, so raw and LF-normalized hashes coincide.

## Accepted limitations

- Conclusions remain limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes.
- Expected headline power is about 62% at a 15-point effect; the approximate 80% MDE against zero is 0.19.
- The USD 300 validation ceiling cannot support 80% power for a 15-point Holm-adjusted LCB-floor design.
- Efficiency-only is a conservative conjunction.
- Cascade remains run-to-run noisy. Its restored four-step label must not be relaxed.
- Live `--world-build` argv is the digest of the host-local `buildPhoenix` executable (`-X main.version=authoring`, host `GOOS`/`GOARCH`). That value is distinct from the frozen Linux-amd64 world-build identity. The templates now record that distinction; they do not make the two digests equal.
- `service_tier: standard` and `fast_mode: false` are recorded in protocol, candidate, and scheduled summary. They are not Claude CLI flags in `prepareRuntimeInvocation`.
- `[Arm B only]` and `[when declared]` are documentation markers in the templates, not literal argv tokens. The implementation emits `--append-system-prompt-file` and `--state-events` without those prefixes.

## Exact next allowed action

Write a focused acceptance patch that:

- copies the four reviewed artifacts from `experiments/frontier-v1/artifacts/pre-validation-local-candidate.json` into `experiments/frontier-v1/pre-validation-artifacts.json`;
- records this revision commit and this review record on those four entries;
- removes only these four strings from `remaining`: `runtime invocation and exact per-arm system prompts`, `arm A schemas`, `world definition and world-build digest`, `grader digest`;
- leaves `schedule digest` in `remaining`;
- keeps `status` `partial`;
- keeps `may_open_validation` and `may_open_held_out` false.

Do not change protocol, runner implementation, world, schemas, grader, corpus, labels, Arm B, analysis, or the already-frozen freeze entries. Do not generate a validation schedule. Do not run a model. Do not open Gate 1A, validation, or held_out.

## Remaining execution blockers

- independent generation and sealing of new validation and held_out families;
- schedule digest;
- Gate 1A, validation execution, and held-out execution.

Passing these four local artifacts does not open those gates.
