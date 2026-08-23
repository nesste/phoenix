# Authoring repair: durable pairing-key checkpoint/resume

Status: implemented in authoring sources; not frozen; not authorized for validation.
Owner: authoring repair after decision 0013.
Frozen validation bytes are unchanged. Freeze tests for the scheduled runner are expected to fail until independent review and refreeze.

## Bug

Decision 0013 closed Gate 1A as indeterminate after the host process was lost at launch 720 of 1,800. The frozen runner required an empty output directory and kept spend only in memory. Re-invoking the same command could not reconstruct pairing-key progress, so the interrupted prefix could not continue and a later authorized run would fail the same way.

This repair does not resume, splice, or reopen that 720-trial archive. The closed validation gate still fails first. Existing trial world-builds in that archive do not match the freeze, so a later reopened validation run cannot attach to it.

## Implemented authoring behavior

Scheduled runs now pin schedule identity before the first model call:

- Before the first launch the runner writes `scheduled-checkpoint.json` with schedule digest, live world-build, run budget, spend 0, and next launch index 0. After every complete A–E group it rewrites that checkpoint with reconstructed spend and the next launch index.
- A nonempty output directory with assignments or evidence is resumable only when that checkpoint is present. Assignment records must be a contiguous prefix of the pinned schedule, evidence files must belong to those assignments, and launched assignments must carry trial world-build evidence matching the live digest.
- Spend is reconstructed from assignment `total_cost_usd` before the next pairing-key budget reserve. An in-progress pairing key continues its remaining arms without repeating the reserve.
- A finished `scheduled-summary.json`, missing checkpoint, unrecognized files, launch gaps, schedule-digest or identity mismatch, world-build mismatch, missing trial world-build evidence, or evidence without an assignment fail closed.
- Empty output still starts a fresh run.

Authoring and, after a later chair decision and refreeze, validation use the same resume path. The 720-trial diagnostic directory remains an archive only.

## Out of scope here

This note does not update freeze hashes, reopen validation, resume the interrupted prefix, or claim that a documented host can finish 1,800 launches inside the provider window. Independent review and a new refreeze are required before any later chair decision can authorize a disjoint validation execution.
