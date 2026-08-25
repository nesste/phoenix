# 0015: Protocol-v4 Gate 1A validation build-recipe import and refreeze

- **Status:** Validation build-recipe repair imported and frozen; outcome gates remain closed
- **Date:** 2026-08-25
- **Decision:** `IMPORT_AND_FREEZE`
- **Independent review:** `docs/reviews/2026-08-25-protocol-v4-gate-1a-validation-build-recipe-review.md` at `8dd34b5eb013dcbe46243ce7768ff5a60b93ed66`; verdict `ACCEPT`; 13,061 raw bytes; `sha256:8a66ba51266c5289535b7472610ff81a47058a19a3706ac72d0ab4968040eedb`
- **Accepted payload:** `ed3860708c931ddc848b1bc90e4d6435585ff0d6`
- **Payload commit:** `ed38607` (already on `main`)
- **Review-record import commit:** `8dd34b5`
- **Focused refreeze commit:** `c3b006ed501765e283ce59a6dc6c069b2f127bde`

## Decision

The project chair accepts the independent review, keeps payload `ed38607` as the implementation, and refreezes the scheduled-runner identity that the repair changed.

The replacement runner inventory is `experiments/frontier-v1/artifacts/gate-1a-validation-build-recipe-candidate.json`, with raw identity `sha256:069b3c202d92845194265c720574fe46cb798637b9edb12965ef77401b24bb3b` over 4,995 bytes. The 20-file runner identity is `sha256:c9e2f4ff1b6c37ea64b8d35b0124d432c1741be01788e51744cdfbd9b56f447f`. It supersedes accepted payload `610fa588488579cf5551a627795d4dc7f771f053` and its 19-file identity `sha256:25cce910f76aaac88cfde6d35c4e44080134547013b19dc4386b6e1a0bc603af`.

The frozen world definition, canonical world identity, linux/amd64 world-build digest `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`, validation schedule, manifest, label registry, grader digest, prompts, and analysis artifacts are unchanged by this decision.

`may_open_validation` remains `false`, and `may_open_held_out` remains `false`. This decision does not open Gate 1A, execute a model or arm, create a validation output, resume or splice the 720-trial diagnostic archive, contact a real custodian, obtain a private grade, inspect a validation or held-out outcome, or authorize a new disjoint 1,800-launch run. Decision 0012 remains spent.

This chair session authored payload `ed38607` and transcribed the independent review record. The independent review was performed by a separate reviewer session with a fresh context that did not author the payload; the same-machine / shared-workspace overlap recorded in the review record stands. The freeze relies on that `ACCEPT`, not on this session's self-tests.

## Accepted replacement

The runner freeze set grows from 19 files to 20 by adding `validation_build_recipe_test.go`. Live hashes also change for `main.go`, `README.md`, and `validation-execution-boundary.md`.

The replacement repairs the last decision-0014 residual: validation scheduled runs now build Phoenix with the frozen world-build recipe — `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.26.6`, `-trimpath -buildvcs=false -ldflags "-s -w -buildid= -X main.version=dev"`, at `bin/phoenix` — and record that target in the world-build manifest. Authoring runs and probes keep the host recipe (`version=authoring`, temporary directory). The closed validation gate still fails first, and the world-build pin from decision 0014 still compares the live digest to the freeze after the build and before any trial; the repair makes that comparison satisfiable instead of removing it. `TestValidationBuildReproducesFrozenWorldBuildDigest` reproduces the frozen world-build digest live from the validation recipe, including from a Windows cross-compile host.

## Known non-blocking residuals

The chair records the independent review's `ACCEPT` with no P0 or P1 findings. The following residuals do not block this freeze:

- P2-1: validation builds write and then remove `bin/phoenix`, overwriting any untracked user binary parked there; `bin/` is gitignored and the location is dictated by the frozen manifest, which hashes the repo-relative executable path;
- P2-2: the target-pin test asserts the build command and environment by substring containment; the live reproduction test closes the gap because any effective recipe deviation changes the world-build digest;
- a validation run still requires a linux/amd64 execution host to run the built binary; this repair fixes build identity only;
- `gocyclo -over 15` still fails on `runScheduledCases`, `loadScheduledResume`, `New` in `internal/activate`, and `verifySpentValidationExecutionGateState`. Those functions are in previously accepted payloads or pre-existing freeze tests. This freeze does not refactor them.

P2-1 through P2-3 from decisions 0010, 0011, and 0013/0014 remain in force for unchanged cap, timeout, closed-gate ordering, and resume behavior.

## Import and refreeze boundary

Payload `ed38607` already landed the implementation. Review import `8dd34b5` adds only the independent review record.

Focused refreeze commit `c3b006e` changes:

- `experiments/frontier-v1/pre-validation-artifacts.json` (scheduled_runner block: new 20-file inventory and review provenance);
- `experiments/frontier-v1/runner/artifact_freeze_test.go`;
- `experiments/frontier-v1/runner/validation_candidate_artifact_test.go`.

The previous 19-file authoring-repair candidate remains in the tree and is now verified against payload `610fa588`, not the working tree. No world, world-build, schedule, grader, prompt, or analysis byte is changed in this freeze. Both gates stay false.

## Verification

After import and refreeze:

- the full test suite passes with both gates closed (`go test -count=1 ./...` exit 0), including the replacement freeze guard, the new 20-file candidate guard, and the historical guards for both superseded candidates;
- the authoritative guard reproduces the 20-file identity `sha256:c9e2f4ff…447f` from the live files and verifies the raw candidate-inventory and independent-review identities;
- `TestValidationBuildReproducesFrozenWorldBuildDigest` reproduces world-build `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e` live from the validation recipe;
- `gofmt`, `go vet`, and staticcheck v0.7.0 are clean on the runner; `gocyclo -over 15` fails only on the four functions named above;
- `git diff --check` is clean;
- no model, arm, schedule, trial, real custodian, private grade, validation result, held-out result, or outcome was run or observed.

## Remaining execution blocker

The repaired boundary is frozen and still closed. A later project-chair decision may authorize a **new disjoint** 1,800-launch validation execution only after setting `may_open_validation` true, leaving `may_open_held_out` false, preserving every frozen identity, and refusing the 720-trial diagnostic archive. The build-recipe precondition from decision 0014 is now satisfied: live validation builds reproduce the frozen world-build digest. Execution additionally requires a linux/amd64 host. Decision 0012 remains spent.
