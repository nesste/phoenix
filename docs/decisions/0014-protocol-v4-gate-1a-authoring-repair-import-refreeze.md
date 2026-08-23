# 0014: Protocol-v4 Gate 1A authoring-repair import and refreeze

- **Status:** Authoring repairs imported and frozen; outcome gates remain closed
- **Date:** 2026-08-23
- **Decision:** `IMPORT_AND_FREEZE`
- **Independent review:** `docs/reviews/2026-08-23-protocol-v4-gate-1a-authoring-repair-review.md` at `1204153d6228a0a06fbb15f70eae7d600da5f93a`; verdict `ACCEPT`; 6,289 raw bytes; `sha256:c770a58bf809cca0039a19054972d8a493bd4fa206255d3a754e52644be9e110`
- **Accepted payload:** `610fa588488579cf5551a627795d4dc7f771f053`
- **Payload commit:** `610fa588` (already on `main`)
- **Review-record import commit:** `1204153`
- **Focused refreeze commit:** `e6738b15cbf5c586f845c62d4198fbf721c90eae`

## Decision

The project chair accepts the independent review, keeps payload `610fa588` as the implementation, and refreezes the validation execution boundary together with the local world and world-build identities that the repairs changed.

The replacement runner inventory is `experiments/frontier-v1/artifacts/gate-1a-authoring-repair-candidate.json`, with raw identity `sha256:d9a85ffce0a06389046b5926001cbde7e3c8f5a8f2d99c075aea8baa8ad0747a` over 4,681 bytes. The 19-file runner identity is `sha256:25cce910f76aaac88cfde6d35c4e44080134547013b19dc4386b6e1a0bc603af`. It supersedes accepted payload `71f9648789decf4cd56ef8a24bc840b0dda7efd9` and its 16-file identity `sha256:9f8c48fb60ced7442924e91ca4e52436961efc8c94dc9c1118256a4b8c1de2d9`.

The frozen world definition raw identity is `sha256:efed86542a6586e4e347e67e7ea1c9fe12038170aaab877981137e22a3fbae96`. The canonical world identity is `sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc`. The frozen linux/amd64 world-build digest is `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`.

`may_open_validation` remains `false`, and `may_open_held_out` remains `false`. This decision does not open Gate 1A, execute a model or arm, create a validation output, resume or splice the 720-trial diagnostic archive, contact a real custodian, obtain a private grade, inspect a validation or held-out outcome, or authorize a new disjoint 1,800-launch run.

This chair session also authored payload `610fa588` and transcribed the independent review record. It did not perform the independent review. The freeze relies on that `ACCEPT`, not on this session’s self-tests. Same-machine / shared-workspace overlap is recorded.

## Accepted replacement

The runner freeze set grows from 16 files to 19 by adding `scheduled_resume.go`, `scheduled_resume_test.go`, and `world_build_pin_test.go`. Live hashes also change for `main.go`, `scheduled_run.go`, `types.go`, and `validation_gate.go`.

The replacement implements decision 0013’s three authoring repairs:

1. unmatched orientation binds inspect/status before the first executable act, and an authored miss does not fall through;
2. validation scheduled runs pin live world-build to the freeze after Phoenix is built and before any trial;
3. pairing-key checkpoint/resume writes `scheduled-checkpoint.json` before the first launch and refuses the 720-trial archive.

Authoring still skips the world-build pin. The closed validation gate still runs first.

## Known non-blocking residuals

The chair records the independent review’s `ACCEPT` with no P0 or P1 findings. The following residuals do not block this freeze:

- leftover `requireEmptyScheduledOutput` is unused on the production path;
- MCP fallback tests use `inspect` while production binds `status`;
- runtime-failure assignments without `*.trial.json` brick resume of that directory by design;
- `buildPhoenix` still uses host `GOOS`/`GOARCH` and `-X main.version=authoring`. The freeze recipe is linux/amd64 and `version=dev`. A later validation opening must use the freeze recipe or the pin will refuse;
- `gocyclo -over 15` fails on `runScheduledCases`, `loadScheduledResume`, `New` in `internal/activate`, and `verifySpentValidationExecutionGateState`. Those functions are in the accepted payload or pre-existing freeze tests. This freeze does not refactor them.

P2-1, P2-2, and P2-3 from decisions 0010 and 0011 remain in force for unchanged cap, timeout, and closed-gate ordering tests.

## Import and refreeze boundary

Payload `610fa588` already landed the implementation. Review import `1204153` adds only the independent review record.

Focused refreeze commit `e6738b1` changes:

- `experiments/frontier-v1/artifacts/gate-1a-authoring-repair-candidate.json` (new 19-file inventory);
- `experiments/frontier-v1/pre-validation-artifacts.json`;
- `experiments/frontier-v1/artifacts/pre-validation-local-candidate.json`;
- `experiments/frontier-v1/artifacts/world-build.linux-amd64.json`;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`;
- `experiments/frontier-v1/runner/validation_candidate_artifact_test.go`.

The previous 16-file test-transition candidate remains in the tree and is now verified against payload `71f9648`, not the working tree. No schedule, grader, prompt, or analysis byte is changed in this freeze. Both gates stay false.

## Verification

After import and refreeze:

- the runner suite passes with both gates closed, including the replacement freeze guard and the new 19-file candidate guard;
- the authoritative guard reproduces the 19-file identity `sha256:25cce910…03af` from the live files;
- the guard verifies the raw candidate inventory and independent-review identities;
- authoring corpus validation passes;
- spec examples and `worlds/dev-repo/world.json` validate;
- `make verify-local-freeze-candidate` reproduces world-build `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`;
- `git diff --check` is clean;
- `make quality` fails only `gocyclo -over 15` on the four functions named above; tests, vet, staticcheck, module verification, and vuln scan passed before that gate;
- no model, arm, schedule, trial, real custodian, private grade, validation result, held-out result, or outcome was run or observed.

## Remaining execution blocker

The repaired boundary is frozen and still closed. A later project-chair decision may authorize a **new disjoint** 1,800-launch validation execution only after setting `may_open_validation` true, leaving `may_open_held_out` false, preserving every frozen identity, using the freeze world-build recipe for validation builds, and refusing the 720-trial diagnostic archive. Decision 0012 remains spent.
