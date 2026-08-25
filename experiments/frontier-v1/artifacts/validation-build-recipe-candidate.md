# Authoring repair: validation builds use the freeze world-build recipe

Status: implemented in authoring sources; not frozen; not authorized for validation.
Owner: authoring repair after decision 0014.
Frozen validation bytes are unchanged in `pre-validation-artifacts.json`. Freeze tests for the scheduled runner are expected to fail until independent review and refreeze.

## Residual

Decision 0014 recorded this non-blocking residual: `buildPhoenix` used host `GOOS`/`GOARCH` and `-X main.version=authoring` for every tranche. The freeze recipe is linux/amd64, `version=dev`, `bin/phoenix`. With the host recipe, a later validation opening would always fail closed at the world-build pin, because the live digest could never equal the frozen `sha256:bf976cadb46b39169c130a2effee1016a5faa9990f63b39c33391f7a97be8a4e`.

## Implemented authoring behavior

The runner now selects a build recipe by tranche before Phoenix is built:

- **Validation scheduled runs** build `./cmd/phoenix` with `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.26.6` and `-trimpath -buildvcs=false -ldflags "-s -w -buildid= -X main.version=dev"`, writing to `bin/phoenix`. The world-build manifest records linux/amd64 and the `bin/phoenix` artifact path, matching the frozen manifest `experiments/frontier-v1/artifacts/world-build.linux-amd64.json`. Cleanup removes the built executable.
- **Authoring runs and probes** keep the previous host recipe: host `GOOS`/`GOARCH`, `version=authoring`, a fresh temporary directory under `build/`.

The closed validation gate still fails first, before any output creation or Phoenix build. After a validation build, the world-build pin from decision 0014 still compares the live digest to the freeze before any trial or model call; this repair makes that comparison satisfiable instead of removing it.

Public-only verification: `TestValidationBuildReproducesFrozenWorldBuildDigest` rebuilds Phoenix with the validation recipe and requires the live manifest digest to equal the frozen world-build digest. `TestValidationBuildRecipeMatchesFrozenWorldBuildTarget` pins the recipe to the frozen target block and the frozen manifest's executable path. No model, arm, trial, custodian, or gate state is touched.

## Out of scope here

This note does not update freeze hashes, reopen validation, change the frozen world-build identity, resume or splice the 720-trial diagnostic archive, or change authoring build behavior. A validation run additionally requires a linux/amd64 execution host; this repair fixes build identity only.
