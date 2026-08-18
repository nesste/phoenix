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

Gate 1A must remain closed. The next design decision is either:

1. amend the protocol and corpus to add an explicit O(1) activation mechanism plus deterministic state-change events, then independently reseal the unopened tranches; or
2. narrow or stop the single-`act` surface thesis before spending the validation budget.
