# Protocol-v4 Gate 1A validation build-recipe independent review

- **Reviewer role:** `gate_1a.validation_build_recipe_independent_reviewer`. Distinct from the authoring agent that produced payload `ed3860708c931ddc848b1bc90e4d6435585ff0d6`, the chair that recorded decision 0014, and any earlier Gate 1A authoring, review, import, or refreeze role.
- **Date:** 2026-08-25
- **Verdict:** `ACCEPT`
- **Review assignment:** conversational, from the project chair after the payload was committed. The reviewed object is the exact commit `ed3860708c931ddc848b1bc90e4d6435585ff0d6`, verified byte-identical to the working tree at review time (`git diff ed38607 --stat` empty, `git status` clean, HEAD = payload).
- **Frozen implementation base / parent:** `0ca98a1ae3f1ca1dedc7db86e5a676b0f1d8f7cc` (`docs: record Gate 1A authoring-repair import and refreeze`). Confirmed as the payload's sole parent and as the candidate's recorded `base_commit`.
- **Authorizing residual:** decision 0014 (`docs/decisions/0014-protocol-v4-gate-1a-authoring-repair-import-refreeze.md`) recorded that `buildPhoenix` used host `GOOS`/`GOARCH` and `-X main.version=authoring`, so a validation opening would always fail closed at the world-build pin. This payload repairs exactly that residual.
- **Payload:** `ed3860708c931ddc848b1bc90e4d6435585ff0d6` (`experiments: author Gate 1A validation freeze-recipe build`)
- **Candidate artifact:** `experiments/frontier-v1/artifacts/gate-1a-validation-build-recipe-candidate.json` — raw SHA-256 `069b3c202d92845194265c720574fe46cb798637b9edb12965ef77401b24bb3b`, 4,995 bytes (file is LF-only; the LF-normalized digest is identical).
- **Candidate note at the payload:** `experiments/frontier-v1/artifacts/validation-build-recipe-candidate.md`
- **Review checkout:** `D:\Work\personal\phoenix` on `main` at the payload commit with a clean tree. Same-machine / shared-workspace limitation recorded: no disposable detached worktree was used this round; byte identity between the reviewed tree and the payload commit was instead established by empty `git diff ed38607` and clean `git status` before and after all checks.

This review does not import the payload, modify freeze hashes, refreeze bytes, open Gate 1A, authorize validation, authorize held-out, resume or splice `experiments/frontier-v1/results/scheduled-validation`, inspect private or held-out labels, or run a model. Both gates remain false. Expected freeze-digest test failures against the old 19-file identity are not defects against this candidate.

## Independence declaration

This session did not author the candidate payload, the candidate note, Phoenix, decision 0014, or any freeze/import decision. It is an independent reviewer session.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. No validation or held-out labels were inspected. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. Nothing under `experiments/frontier-v1/results/` was read or written. No custodian or grader was contacted. No gate state was mutated. The only build artifacts produced were `go build`/`go test` outputs removed by their own cleanups.

## Verdict

The decision-0014 residual repair is present, correct, and fail-closed. Every digest in the 20-file candidate inventory matches the working tree, the set digest reproduces, the frozen linux/amd64 recipe reproduces the frozen world-build digest live on this host, and the full suite fails only in the two documented pre-refreeze freeze tests. There are no P0 or P1 findings. This `ACCEPT` authorizes only a later project-chair import/refreeze of payload `ed3860708c931ddc848b1bc90e4d6435585ff0d6`. It does not freeze, open validation, or authorize any execution.

## Required-check matrix

| Check | Result |
| --- | --- |
| A. Byte identity: all 20 candidate file digests match working tree under LF normalization | Pass (20/20 match; 0 mismatches) |
| B. Artifact set digest over ordinal-path-sorted `<path>\t<digest>\n` lines | Pass — computed `sha256:c9e2f4ff1b6c37ea64b8d35b0124d432c1741be01788e51744cdfbd9b56f447f`, equals claimed |
| C. `base_commit` = parent of payload; diff `0ca98a1..ed38607` touches only the six claimed files | Pass (README.md, main.go, validation_build_recipe_test.go, validation-execution-boundary.md, two new artifacts; +291/−16) |
| D. Set delta vs frozen 19-file `scheduled_runner`: adds exactly `validation_build_recipe_test.go`; digest changes only in README.md, main.go, validation-execution-boundary.md | Pass |
| E. No gate/freeze tampering: pre-validation-artifacts.json, protocol/schedule/manifest/label/grader files, artifact_freeze_test.go, validation_candidate_artifact_test.go, results/ all untouched; both gates false in freeze doc and candidate JSON | Pass |
| F. Validation recipe = linux/amd64/`version=dev`/`bin/phoenix`; all other tranches = host GOOS/GOARCH, `version=authoring`, temp dir under `build/` | Pass |
| G. Build command/env implement the frozen `target` block (`goos=linux`, `goarch=amd64`, `cgo_enabled=false`, `version=dev`, `trimpath`, `buildvcs=false`, `empty_buildid`) and the frozen manifest's executable path `bin/phoenix` | Pass |
| H. Manifest GOOS/GOARCH inputs come from the recipe, not `runtime.GOOS/GOARCH` | Pass |
| I. Tranche threading: `runProbeCLI` passes `"authoring"`; `prepareScheduledCLI` passes the real tranche; all consumers go through `effectiveTranche`, so probe behavior is unchanged | Pass |
| J. Ordering: `requireValidationGate` and `requireFrozenValidationTrialLimits` run before output-path resolution, output creation, and Phoenix build; `requireFrozenWorldBuild` runs after build and before any trial; no production path builds a validation binary while the gate is closed | Pass |
| K. Focused tests: 4/4 pass, including live reproduction of frozen digest `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e` from a Windows host cross-compile at `bin/phoenix` | Pass |
| L. Full suite `go test -count=1 ./...`: only `TestPreValidationFreezeMatchesAcceptedCandidates` and `TestGate1AAuthoringRepairCandidateMatchesWorkingTree` fail, both at the README.md digest vs the OLD frozen identity | Pass (documented pre-refreeze state; no other failures) |
| M. Public-only tests; `execution_performed=false` truthful | Pass |
| N. Hostile-environment probe: parent env `GOOS=windows GOARCH=386` still reproduces the frozen digest (appended overrides win via documented `os/exec` last-entry dedupe) | Pass |
| O. Hygiene: `gofmt -l` empty, `go vet` clean, staticcheck v0.7.0 clean, gocyclo -over 15 shows only the three pre-existing violations (runScheduledCases 27, loadScheduledResume 24, verifySpentValidationExecutionGateState 16) — no regressions | Pass |

## Findings

No P0 findings. No P1 findings. P2 findings below; none blocks.

- **P2-1 (cleanup can destroy a pre-existing user file).** `phoenixExecutablePath` for the validation recipe builds directly to `<root>/bin/phoenix`, overwriting any pre-existing file there, and its cleanup `os.Remove(executable)` then deletes it. A user binary parked at `bin/phoenix` would be silently clobbered by a validation build or by `TestValidationBuildReproducesFrozenWorldBuildDigest`. Mitigations observed: `bin/` is gitignored (`.gitignore` line 1, `/bin/`), the path is dictated by the frozen manifest (the manifest hashes the repo-relative executable path, so the location is load-bearing and cannot move), and no tracked file can be affected. Non-blocking.
- **P2-2 (substring assertions in the target-pin test).** `TestValidationBuildRecipeMatchesFrozenWorldBuildTarget` asserts the build command and environment by substring containment on the joined args/env, so a hypothetical superset command carrying additional contradictory flags would still pass that test. In practice `TestValidationBuildReproducesFrozenWorldBuildDigest` closes the gap: any effective deviation from the frozen recipe changes the executable digest and fails the reproduction test. Non-blocking.

## Behavior confirmed

`phoenixBuildRecipeForTranche("validation")` returns `{goos: linux, goarch: amd64, version: dev, output: bin/phoenix}`; every other tranche returns `{runtime.GOOS, runtime.GOARCH, authoring, ""}` and builds into a fresh `build/authoring-runner-*` temp dir removed by `os.RemoveAll`. `phoenixBuildCommand` runs `go build -trimpath -buildvcs=false -ldflags "-s -w -buildid= -X main.version=<recipe>" -o <executable> ./cmd/phoenix` with `CGO_ENABLED=0 GOTOOLCHAIN=go1.26.6 GOOS=<recipe> GOARCH=<recipe>` appended after `os.Environ()`; the appended entries win under `os/exec` duplicate-key semantics, verified empirically with conflicting parent-environment values. The `.exe` suffix is keyed on `recipe.goos` and applies only to the temp-dir (authoring) branch, so a validation build on Windows correctly emits the suffix-free linux binary at `bin/phoenix`. The world-build manifest records `recipe.goos`/`recipe.goarch`, fixing the latent old behavior where the manifest recorded `runtime.GOOS/GOARCH` regardless of the environment actually used. All `buildPhoenix` error paths (build failure, `filepath.Rel` failure, manifest failure) invoke cleanup before returning.

In `prepareScheduledCLI`, the closed-gate check (which also requires `may_open_held_out` to remain false and re-verifies the frozen schedule, manifest, and label-registry identities) and the frozen trial-limit check (exactly 0.15 USD, exactly 180s) run before any output-directory default, schedule preparation, output creation, or Phoenix build. `requireFrozenWorldBuild` compares the live manifest digest to the frozen `world_build_digest` after the build and before any trial, cleaning up on mismatch. The only code that builds a validation-recipe binary while the gate is closed is the new reproduction test, which performs a build and digest comparison only — no trial, no output directory, no custodian contact, no model.

The candidate JSON's contract values were cross-checked against the freeze: schedule canonical digest `sha256:b38a0eaab063…075813`, manifest canonical digest `sha256:57ccc0c754f7…fd490`, world-build digest `sha256:bf976cadb46b…be8a4e`, budget 300 USD, cap 0.15, timeout 180, 120 cases, 1,800 launches — all equal to the frozen constants and freeze-document values. The new tests read only public inputs: `pre-validation-artifacts.json`, `artifacts/world-build.linux-amd64.json`, and the candidate JSON itself. `git status` was clean after all test runs; `execution_performed=false` is truthful.

## Residual notes (non-blocking)

These do not block `ACCEPT`:

- Freeze-digest tests (`TestPreValidationFreezeMatchesAcceptedCandidates`, `TestGate1AAuthoringRepairCandidateMatchesWorkingTree`) remain expected-fail until an authorized refreeze imports the 20-file identity; both currently fail at the first divergent file (README.md) against the old frozen digests.
- Authoring builds now pin `GOOS`/`GOARCH` to the host instead of inheriting operator-set environment values. This deviates from strict byte-for-byte equivalence only when an operator's environment already set `GOOS`/`GOARCH`, and the new behavior is the correct one — the old path would have cross-compiled per environment while misrecording `runtime.GOOS/GOARCH` in the manifest.
- Validation-recipe cleanup removes only the `bin/phoenix` file; an empty `bin/` directory created by `MkdirAll` remains afterward (gitignored, zero-byte residue).
- A real validation run still requires a linux/amd64 execution host to run the built binary; this repair fixes build identity only, as the candidate note itself states.
- The reproduction test is skipped under `-short`; the non-short full suite exercised it here.
- Same-machine / shared-workspace limitation: this review ran in the primary checkout at the payload commit with a clean tree rather than a disposable detached worktree.

## Public checks

Byte-identity script (scratchpad, LF-normalization rule): 20/20 file digests match; set digest reproduces `sha256:c9e2f4ff1b6c37ea64b8d35b0124d432c1741be01788e51744cdfbd9b56f447f`. Focused tests `go test ./experiments/frontier-v1/runner/ -run 'BuildRecipe|ReproducesFrozenWorldBuild|ValidationBuildRecipeCandidate' -v`: 4/4 pass, including live reproduction of `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`. Full suite `go test -count=1 ./...`: only the two documented freeze failures. Hostile-environment rebuild (`GOOS=windows GOARCH=386` in parent env): frozen digest still reproduced. `gofmt`, `go vet`, staticcheck v0.7.0: clean. gocyclo -over 15: pre-existing three violations only.

**Final verdict: ACCEPT** — the payload repairs the decision-0014 residual exactly as contracted, with verified byte identity, no gate or freeze tampering, preserved fail-closed ordering, live reproduction of the frozen world-build digest, and only the documented pre-refreeze freeze-test failures; the two P2 findings and all residuals are non-blocking.
