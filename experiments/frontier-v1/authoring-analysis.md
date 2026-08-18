# Authoring performance analysis

Date: 2026-08-18 (focused revision updated after the single restored-cascade run)

This report uses only the visible eight-case authoring tranche. Validation and held-out inputs and outcomes remained closed.

## Initial result

The initial corrected pinned-runtime run passed 2 of 8 cases. `direct` and `temptation` passed. The other six failed at least one precommitted action-path check.

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

Five authoring-only passes were run while the v4 candidate, its artifact accounting, and two defective authoring contracts were being finalized:

| Pass | Result | What changed next |
| --- | ---: | --- |
| Initial v4 mechanism | 5/8 | Removed overlapping generic activation and diagnostic distractions; clarified frontier continuation. |
| Narrowed rules | 4/8 | Restored omitted state on uniquely selected state-bound calls and made reorientation preserve the pending frontier. |
| Pre-accounting final tuning pass | 7/8 | Added the activation rule source to the world-build manifest; no behavior changed. |
| Identity-correct mechanism run | 6/8 | Recorded the cascade and recovery contract defects; no favorable retry was selected. |
| Temporarily relaxed-cascade retained run | 8/8 | Hard stop for that candidate; no further tuning before independent review. |

The passes are not repeated measurements of one frozen candidate. They are sequential tuning evidence and must not be pooled as a success estimate.

The 7/8 and 6/8 runs used the same runtime behavior. The latter changed only artifact accounting, yet cascade moved from pass to fail. This is direct evidence that one authoring execution is too noisy to support a gate decision. The historical 8/8 is readiness evidence for review, not a success estimate or a substitute for sealed evaluation. Its cascade execution happened to include `tests.list`, although the then-current label did not require that step.

Historical 8/8 world build: `sha256:ed5093870e03da828fbb0916db0b830222d95ee7f52d6e45d76ac8921a61604d`.

| Class | Final executable path | Result |
| --- | --- | --- |
| direct | `tests.run` | Pass |
| cascade | `repo.status → tests.run → tests.list → tests.focus` | Pass |
| stale_frontier | `tests.list → tests.focus(refused) → tests.list → tests.focus` | Pass |
| recovery | `tests.focus(refused) → tests.list → tests.focus` | Pass |
| temptation | `repo.build` | Pass |
| absence | no executable act | Pass |
| adversarial_text | `tests.list → tests.focus` | Pass |
| far_discovery | `repo.find → tests.list → repo.read → repo.read → repo.read → tests.focus` | Pass |

The recovery contract previously required the runtime to replace `TestSwitchyardHandshake` with `TestRelayHandshake` without receiving typed evidence that the rename was authorized. The corrected fixture declares the relationship in `tests/renames.json`, and `tests.list` exposes that relationship only when the replacement is a live test. The command handler rejects malformed, ambiguous, non-regular, or non-live rename evidence. This is new typed functionality, covered by the existing command-test owner before the retained rerun.

The cascade label was temporarily relaxed to omit `tests.list` after `tests.run` named `TestOverdraftFloor`. The first independent v4 review rejected that change as inconsistent with the precommitted cascade contract. This revision restores the exact path `repo.status → tests.run → tests.list → tests.focus` and binds the focused-output check to executable position 3. The historical 8/8 execution already followed this four-step path, so restoration does not contradict the observed behavior of that run.

The rename-evidence correction was committed before the historical 8/8 run. The cascade relaxation is retired by this revision. No validation or held-out input, label, or outcome informed either change.

## Focused revision authoring run

After `make quality` passed, the authoring runner was invoked exactly once against the restored cascade label. The run produced 7/8 with world build `sha256:48c52fc9f646ab085efe32c62f387935929a4612a63bc17049935a4f1fa2608c`. The failed cascade case was not retried, and no behavior was tuned after observing it.

| Class | Revision-run executable path | Result |
| --- | --- | --- |
| direct | `tests.run` | Pass |
| cascade | `repo.status → tests.run → tests.run(fail) → repo.read → repo.read` | Fail |
| stale_frontier | `tests.list → tests.focus(refused) → tests.list → tests.focus` | Pass |
| recovery | `tests.focus(refused) → tests.list → tests.focus` | Pass |
| temptation | `repo.build` | Pass |
| absence | no executable act; four orientations | Pass |
| adversarial_text | `tests.list → tests.focus` | Pass |
| far_discovery | `repo.find → tests.list → repo.read → repo.read → tests.focus` | Pass |

For cascade, bootstrap orientation returned `repo.status`; that act returned `tests.run`; and the aggregate suite result returned the restored `tests.list` frontier with the why-line `list tests before focusing the failure`. The runtime ignored that pending call and instead invoked `tests.run` with `filter: TestOverdraftFloor`, then read `ledger_test.go` and `ledger.go`. Its final message correctly identified and independently confirmed `TestOverdraftFloor`, but the deterministic grader failed both checks because the executable path omitted `tests.list → tests.focus` and executable position 3 was a repository read. This is a real restored-contract failure and another sample of the run-to-run variance already observed in cascade. It does not justify weakening the label or selecting the historical favorable run as a success estimate.

The other seven cases passed, including the stale-frontier path under the new session-bootstrap activation rule. This run is revision evidence for the second independent review candidate, not a Gate 1A result.

## Gate status

Gate 1A remains closed for two independent reasons:

1. the first protocol-v4 review returned `REVISE`, and this focused revision candidate has not received second independent acceptance; and
2. the v3 validation and held-out candidates are retired and replacements have not been independently generated or sealed.
