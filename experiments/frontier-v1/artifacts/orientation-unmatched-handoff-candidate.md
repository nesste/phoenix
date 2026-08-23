# Authoring repair: unmatched-orientation handoff

Status: implemented in authoring sources; not frozen; not authorized for validation.
Owner: authoring repair after decision 0013.
Frozen validation bytes are unchanged. Freeze tests for the world and world-build are expected to fail until independent review and refreeze.

## Bug

C/D/E prompts send the complete goal as intent. The world had six regex activations. On public validation goals, 18/120 match. Unmatched orientation returned `matched: false` and an empty frontier. The model repeated intent until the 12-turn cap. After the first executable act, reactivation is correctly forbidden, so an empty post-act frontier is also a dead end.

On the interrupted prefix, C passed 19/24 matched trials and 2/120 unmatched trials. 96/144 C trials had zero acts.

## Implemented authoring behavior

Before the first executable act:

- Specific authored activations still bind first.
- If a specific rule matches and binds no calls, that is an authored miss. Orientation returns no calls and tells the model not to repeat the intent. Absence uses `absent_deploy_or_release` so deploy/release goals do not fall through to inspect.
- If no specific rule matches, orientation binds `inspect_reachable_repository` (`repo.status`), not an empty retry loop.

After the first executable act, orientation still does not reactivate from intent. If the act returns no pending frontier, the session is done unless a pending refusal alternative exists.

## Out of scope here

World-build pin and checkpoint/resume remain separate repair candidates. The frozen runner, prompts, validation corpus, and freeze hashes are not updated. Those require independent review and a new refreeze before any later chair decision can authorize a disjoint validation execution.
