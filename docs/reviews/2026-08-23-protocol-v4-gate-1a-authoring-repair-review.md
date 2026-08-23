# Protocol-v4 Gate 1A authoring-repair independent review

- **Reviewer role:** `gate_1a.authoring_repair_independent_reviewer`. Distinct from the authoring agent that produced payload `610fa588488579cf5551a627795d4dc7f771f053`, the chair that recorded decision 0013, and any later import/refreeze or Gate 1A chair.
- **Date:** 2026-08-23
- **Verdict:** `ACCEPT`
- **Review assignment:** conversational, from the project chair after that payload was committed. No `docs: assign …` commit exists for this round. The reviewed object is the exact commit `610fa588488579cf5551a627795d4dc7f771f053`, not later working-tree state.
- **Frozen implementation base / parent:** `0ffb8f9b44edca25582ec35541d25257d6a36dd0` (`Record the interrupted Gate 1A run as indeterminate and keep the 720-trial prefix as a diagnostic archive.`)
- **Authorizing decision:** `docs/decisions/0013-protocol-v4-gate-1a-interrupted-execution.md`
- **Payload:** `610fa588488579cf5551a627795d4dc7f771f053` (`experiments: author Gate 1A orientation, world-build pin, and pairing-key resume`)
- **Candidate notes at the payload:**
  - `experiments/frontier-v1/artifacts/orientation-unmatched-handoff-candidate.md`
  - `experiments/frontier-v1/artifacts/world-build-pin-candidate.md`
  - `experiments/frontier-v1/artifacts/pairing-key-checkpoint-resume-candidate.md`
- **Review checkout:** `D:\Work\personal\phoenix-review-610fa588` (disposable detached worktree at the payload; same-machine / shared-workspace limitation recorded)
- **Prior independent review of the uncommitted tree:** `REVISE`, two P1 resume defects. This review re-verifies those repairs on the committed payload.

This review does not import the payload, modify freeze hashes, refreeze bytes, open Gate 1A, authorize validation, authorize held-out, resume or splice `experiments/frontier-v1/results/scheduled-validation`, inspect private or held-out labels, or run a model. Both gates remain false. Expected freeze-digest test failures are not defects against this authoring candidate.

## Independence declaration

This session did not author the candidate payload, the three candidate notes, Phoenix, decision 0013, or any freeze/import decision. Same-machine / shared-workspace limitation: the review ran on the authoring host from a disposable detached worktree, not by editing `D:\Work\personal\phoenix` on `main`.

No private evaluator root was opened, hashed, mounted, enumerated, or copied. No validation or held-out labels were inspected. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. `seal --write` was not run.

## Verdict

Decision 0013’s three authoring repairs are present and fail-closed. There are no P0 or P1 findings. Both prior resume P1s are fixed. This `ACCEPT` authorizes only a later project-chair import/refreeze of payload `610fa588488579cf5551a627795d4dc7f771f053`. It does not freeze, open validation, or authorize a new disjoint 1,800-launch run.

## Required-check matrix

| Check | Result |
| --- | --- |
| A. Orientation: unmatched inspect/status fallback before first executable act; authored miss does not fall through; intent cannot reactivate after first executable act | Pass |
| B. World-build pin: validation scheduled runs compare live digest to freeze after Phoenix is built and before any trial/model call; mismatch/missing freeze/empty live digest fails and cleans up; authoring skips; closed gate still fails first | Pass |
| C. Resume: load/verify/write checkpoint; resume loop; nonempty output with assignments/evidence requires a checkpoint | Pass |
| D. Closed gate still precedes output mutation / model calls | Pass |
| E. Candidate does not silently update freeze hashes or reopen gates | Pass |
| F. Focused tests (no validation opening); freeze-digest failures not treated as defects | Pass |
| Prior P1: schedule identity checked even when no checkpoint exists yet | Pass (fixed) |
| Prior P1: launched assignment without `*.trial.json` must not skip world-build check | Pass (fixed) |

## Findings

No findings.

## Behavior confirmed

Unmatched orientation binds `inspect_reachable_repository` (`repo.status`) before the first executable act. `absent_deploy_or_release` is an authored miss and does not fall through to inspect. After the first executable act, intent is not recomputed.

Validation scheduled runs pin live world-build to `world_definition_and_world_build_digest.world_build_digest` in `experiments/frontier-v1/pre-validation-artifacts.json` after Phoenix is built and before trials. Authoring skips that pin. `may_open_validation` and `may_open_held_out` stay false. The closed gate still runs first.

Resume writes `scheduled-checkpoint.json` at next index 0 before the first launch, reconstructs spend from assignment `total_cost_usd`, continues incomplete pairing keys without re-reserving, and reserves `5 × per-trial cap × 1.10` only for new keys. It refuses finished summaries, unrecognized files, launch gaps, missing checkpoints when assignments or evidence exist, schedule digest or identity mismatch, world-build mismatch, launched assignments without trial world-build evidence, and evidence without an assignment. The 720-trial archive still has no checkpoint, so it cannot splice.

## Residual notes (non-blocking)

These do not block `ACCEPT`:

- leftover `requireEmptyScheduledOutput` is unused on the production path (`TestScheduledOutputMustBeEmpty` still tests the old helper)
- MCP fallback tests use `inspect`; production worlds bind `status` (covered by activation tests)
- runtime-failure assignments without `*.trial.json` brick resume of that directory by design
- freeze-digest tests remain expected-fail until refreeze
- `buildPhoenix` still uses host `GOOS`/`GOARCH` and `-X main.version=authoring`; the frozen world-build recipe is linux/amd64 and `version=dev`. The pin is fail-closed. A later validation opening must use the freeze recipe or validation will refuse at the pin.

## Public checks

Focused runner tests `ScheduledResume|RequireScheduledOutput|ScheduledOutputMust|ScheduledBudget|WorldBuildPin` and unmatched-orientation activation/MCP tests passed. Full runner freeze-digest tests were expected to fail and were not treated as defects.
