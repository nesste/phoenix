# Authoring repair: world-build pin before first model call

Status: implemented in authoring sources; not frozen; not authorized for validation.
Owner: authoring repair after decision 0013.
Frozen validation bytes are unchanged. Freeze tests for the scheduled runner are expected to fail until independent review and refreeze.

## Bug

The interrupted Gate 1A prefix recorded live world-build `sha256:bbcf6e83…` on every trial. The freeze names `sha256:27c2f537…`. `buildPhoenix` rebuilds with `-X main.version=authoring` and host `GOOS`/`GOARCH`. The freeze recipe is linux/amd64, `version=dev`, `bin/phoenix`. The runner never compared those identities before the first model call.

## Implemented authoring behavior

After Phoenix is built for a **validation** scheduled run, and before any trial or model call, the runner compares the live world-build digest to `world_definition_and_world_build_digest.world_build_digest` in `experiments/frontier-v1/pre-validation-artifacts.json`. A mismatch, missing freeze, or empty live digest fails the run and cleans up the build.

Authoring runs do not use this check. They may still build a host-local authoring binary.

The closed validation gate still fails first. This pin is the next load-bearing identity once a later chair decision reopens validation.

## Out of scope here

Checkpoint/resume, or a host that can finish 1,800 launches inside the provider window, remains a separate repair. This note does not update freeze hashes, reopen validation, or change the authoring build recipe.
