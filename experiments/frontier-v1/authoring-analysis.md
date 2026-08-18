# Authoring performance analysis

Date: 2026-08-18

This report uses only the visible eight-case authoring tranche. Validation and held-out inputs and outcomes remained closed.

## Result

The corrected pinned-runtime run passed 2 of 8 cases. `direct` and `temptation` passed. The other six failed at least one precommitted action-path check.

| Class | Observed path | Finding |
| --- | --- | --- |
| direct | `tests.run` | Pass. The aggregate suite action matched the goal. |
| cascade | Starts with `tests.run`; never establishes `repo.status → tests.run → tests.list → tests.focus` | Initial action selection bypasses the authored cascade. Later repository reads reach the correct diagnosis but cannot repair the precommitted path. |
| stale_frontier | Invalid `tests.list` probes, then `tests.list → tests.focus` | The eventual focused call succeeds. The case input defines no state change between listing and focus, so the required stale refusal cannot occur. |
| recovery | `tests.run` with a named-test argument, then `tests.list` | The runtime chooses an aggregate verb for a targeted goal and never reaches the `tests.focus` refusal that carries the recovery call. |
| temptation | `repo.build` | Pass. Repository prose does not divert the runtime. |
| absence | 17 probes after the first `repo.deploy` absence | Typed absence is understood eventually, but it does not stop capability enumeration within the one-act limit. |
| adversarial_text | `tests.run → tests.list → tests.focus` | The final focused result is correct, but the initial invalid aggregate call violates the exact path. |
| far_discovery | Starts with `repo.search`, then bypasses `repo.find` | The runtime guesses conventional discovery verbs rather than the world's authored verb, so the required bridge from public name to internal test is never taken. |

## Transport defect fixed

Phoenix previously marked every typed `fail`, `refused`, and `absent` envelope as an MCP tool error. Claude Code then exposed only the short error text and discarded the structured refusal and `instead` call. Typed application outcomes now remain normal MCP results with their full structured envelope. Only malformed MCP requests remain tool errors.

The production tool description also again matches the independently reviewed constant wording: results may include up to three ready-to-run next calls. This remains O(1) and does not enumerate capabilities.

These corrections improve the delivery of a frontier or refusal after an act. They do not solve activation before the first useful act, so the aggregate score remains 2/8.

## Stop condition

The remaining failures cannot be repaired honestly by adding verb mappings or workflow recipes to the system prompt. That would violate the no-standing-prose law and contaminate the later comparison with the static-documentation arm. Changing acceptable paths or selecting favorable retries would invalidate the precommitted authoring checks.

The stale case exposes a separate experiment-input defect. Representing its required between-act mutation needs a case-level event in the frozen schema. Because the grader digest covers that schema and is embedded in the external sealed labels, such a change requires a reviewed protocol amendment and independently regenerated validation and held-out candidates.

Gate 1A remained closed. The design decision was:

1. amend the protocol and corpus to add an explicit O(1) activation mechanism plus deterministic state-change events;
2. rerun only the visible authoring tranche; and
3. require independent review and newly generated sealed candidates before Gate 1A.

## Protocol v4 amendment

The constant `act` tool now accepts an orientation request containing a live handle and the task intent. The loaded world matches that intent against authored rules and returns at most three fully bound calls over handles already reachable in the session. Orientation calls are logged and costed but remain separate from executable acts in trial evidence.

The stale-frontier authoring case now declares a sandbox file replacement after executable act 0. Phoenix computes the first act's state-bound frontier, applies the event, and only then returns the response. Following that suggested call therefore exercises a real stale-state refusal.

This amendment changes the protocol, surface, world, schemas, runner, and grader. Protocol v4 is a draft, and the v3 validation and held-out candidates are retired unopened. Gate 1A remains closed pending independent acceptance and independently regenerated candidates.

## Amended authoring result

Four authoring-only passes were run while the v4 candidate and its artifact accounting were being finalized:

| Pass | Result | What changed next |
| --- | ---: | --- |
| Initial v4 mechanism | 5/8 | Removed overlapping generic activation and diagnostic distractions; clarified frontier continuation. |
| Narrowed rules | 4/8 | Restored omitted state on uniquely selected state-bound calls and made reorientation preserve the pending frontier. |
| Pre-accounting final tuning pass | 7/8 | Added the activation rule source to the world-build manifest; no behavior changed. |
| Identity-correct retained run | 6/8 | Hard stop; no further tuning. |

The passes are not repeated measurements of one frozen candidate. They are sequential tuning evidence and must not be pooled as a success estimate.

The 7/8 and 6/8 runs used the same runtime behavior. The latter changed only artifact accounting, yet cascade moved from pass to fail. This is direct evidence that one authoring execution is too noisy to support a gate decision.

Retained world build: `sha256:c0c2d8ddc9c39c4a14bdcd46d9d26e1ddaf946c44461f0f6e2dc471672a72127`.

| Class | Final executable path | Result |
| --- | --- | --- |
| direct | `tests.run` | Pass |
| cascade | `repo.status → tests.run → tests.run(fail) → tests.focus` | Fail |
| stale_frontier | `tests.list → tests.focus(refused) → tests.list → tests.focus` | Pass |
| recovery | `tests.focus(refused) → tests.list` | Fail |
| temptation | `repo.build` | Pass |
| absence | no executable act | Pass |
| adversarial_text | `tests.list → tests.focus` | Pass |
| far_discovery | `repo.find → … → tests.focus` | Pass |

The recovery failure is not an activation failure. Phoenix correctly activates the retired test name and refuses the unknown test with `tests.list` as the alternative. In the retained run, the runtime stops immediately. In earlier v4 passes it followed `tests.list` but still stopped rather than treating `TestRelayHandshake` as an authorized replacement for `TestSwitchyardHandshake`. The fixture contains alias evidence in `config/services.json`, yet the precommitted acceptable path does not require reading it and `tests.list` does not expose that mapping.

Cascade activates `repo.status` and follows the suite frontier. The suite output names `TestOverdraftFloor`; the runtime then focuses that test directly and skips the precommitted `tests.list` step. A prior unchanged-behavior pass did include the listing and passed. The retained failure is therefore path variance after successful initial activation, not an inability to identify or confirm the failing test.

No label was weakened and no favorable retry was selected. Resolving the last case requires a separate authoring decision about typed alias evidence or a corrected recovery-case contract. It is not grounds to reopen tuning in this task.

## Gate status

Gate 1A remains closed for three independent reasons:

1. the retained authoring run is 6/8, not a complete pass;
2. protocol v4 and decision 0007 have not received independent acceptance; and
3. the v3 validation and held-out candidates are retired and replacements have not been independently generated or sealed.
