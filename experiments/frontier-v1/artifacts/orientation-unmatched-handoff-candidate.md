# Authoring repair candidate: unmatched-orientation handoff

Status: draft, not frozen, not authorized for validation.
Owner: authoring repair after decision 0013.
Must not change frozen validation bytes until independent review and refreeze.

## Bug

C/D/E prompts send the complete goal as intent. The world has six regex activations. On public validation goals, 18/120 match. Unmatched orientation returns `matched: false` and an empty frontier. The model repeats intent until the 12-turn cap. After the first executable act, reactivation is correctly forbidden, so an empty post-act frontier is also a dead end.

On the interrupted prefix, C passed 19/24 matched trials and 2/120 unmatched trials. 96/144 C trials had zero acts.

## Required authoring behavior

Before the first executable act, if no authored activation binds, orientation must still return at most three fully bound calls from an authored fallback (inspect/status), not an empty retry loop. Absence/temptation cases that must not act need an authored miss that tells the model to stop, not to retry the same intent.

After the first executable act, do not reactivate from intent. If the act returns no pending frontier, the session is done unless a pending refusal alternative exists.

## Out of scope here

World-build pin and checkpoint/resume are separate repair candidates. This note does not authorize editing `world.json`, `activations.go`, or the frozen runner.
