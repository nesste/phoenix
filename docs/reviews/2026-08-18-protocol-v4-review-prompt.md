# Protocol v4 independent review prompt

After the final authoring candidate is committed, replace the five angle-bracketed candidate fields. Then copy the text below into a fresh evaluator session that did not implement Phoenix.

```text
Act as the independent evaluation reviewer for the Phoenix frontier-v1 protocol v4 amendment. This is a protocol and implementation-contract review. It is not an outcome run, a corpus-authoring task, or permission to open a sealed tranche.

Repository: https://github.com/nesste/phoenix
Candidate commit: <COMMIT>

Identify the candidate before reviewing it. Compute SHA-256 over LF-normalized UTF-8 bytes and require these exact values:

- experiments/frontier-v1/protocol.json: <PROTOCOL_SHA256>
- docs/decisions/0007-authoring-activation-amendment.md: <DECISION_0007_SHA256>
- docs/decisions/0004-go-no-go-rules.md: <DECISION_0004_SHA256>
- experiments/frontier-v1/authoring-analysis.md: <AUTHORING_ANALYSIS_SHA256>

Return REVISE without reviewing the candidate if the commit or any digest differs.

Review scope

Read these files from the exact candidate commit:

- experiments/frontier-v1/protocol.json
- docs/decisions/0004-go-no-go-rules.md
- docs/decisions/0007-authoring-activation-amendment.md
- experiments/frontier-v1/authoring-analysis.md
- experiments/frontier-v1/results/authoring/summary.json
- internal/surface/mcp.go and internal/surface/budget_test.go
- internal/activate/activate.go and internal/activate/activate_test.go
- internal/world/load.go and internal/world/semantic.go
- cmd/phoenix/state_events.go and cmd/phoenix/state_events_test.go
- internal/episode/log.go and internal/episode/sqlite.go
- spec/world.schema.json, spec/episode.schema.json, and spec/result.schema.json
- experiments/frontier-v1/schema/case.schema.json and experiments/frontier-v1/schema/trial.schema.json
- experiments/frontier-v1/runner/main.go, runtime.go, and types.go
- experiments/frontier-v1/corpusctl/internal/corpus/protocol_test.go, sealed_test.go, trial.go, and grade.go
- verbs/dev-repo/activations.go, transitions.go, refusals.go, and world_contract_test.go
- worlds/dev-repo/world.json and experiments/frontier-v1/worlds/authoring.dev_repo.json

Consult docs/plans/2026-08-17-phoenix-world-plan.md only to recover the original Task 0.5 and Gate 1A intent. Where that plan conflicts with the explicit v4 amendment, report the inconsistency instead of silently choosing one.

Independence and custody rules

- Do not modify any repository file.
- Do not run authoring, validation, or held_out trials.
- Do not open or inspect validation or held_out cases, fixtures, labels, manifests, registries, or outcomes. The v3 candidates are retired custody artifacts, not review inputs.
- Do not generate, relabel, reseal, or replace validation or held_out candidates.
- Authoring cases, labels, trials, and the retained summary are visible and may be inspected when needed to assess the amendment rationale.
- Read-only source inspection, digest computation, schema validation, and tests that construct their own temporary data are allowed.
- Do not treat the historical v3 ACCEPT record as acceptance of v4.
- Keep protocol-design acceptance separate from authoring performance, Gate 1A readiness, and authorization to run a sealed tranche.

Allowed checks

Run these checks if the environment permits:

1. `go test -count=1 ./...` from the repository root.
2. `go test -count=1 ./...` from `experiments/frontier-v1/corpusctl`.
3. `go run ./cmd/corpusctl validate --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json` from `experiments/frontier-v1/corpusctl`.
4. Independently serialize the served MCP instructions and tool list for 10, 100, and 1,000 synthetic verbs. Confirm one standing tool, empty instructions, no enumerated verb catalogue or workflow prose, no more than 600 ASCII bytes under the repository's conservative budget, and no size-dependent growth beyond the accepted spread.
5. Confirm that the production and authoring world JSON files are byte-identical.

Do not run `make quality` if it would write a build manifest or executable. Reproduce its read-only checks separately when useful.

Required review questions

1. Amendment necessity: Does the visible authoring evidence support the stated initial-activation gap and the need for a deterministic between-act state change? Separate evidence for the problem from evidence that the chosen solution works.

2. Constant surface: Does Phoenix still expose exactly one standing `act` tool with O(1) instructions and schema size? Does adding `intent` preserve the no-standing-documentation claim, or does it move an unbounded capability map into another channel?

3. Input contract: Are executable input (`handle`, `verb`, `args`, optional `state`) and orientation input (`handle`, `intent`) mutually exclusive at admission? Are malformed, mixed, or incomplete requests rejected deterministically?

4. Activation scope: Can orientation use only world-authored rules and handles already reachable in the live session? Confirm that it cannot inspect the global verb registry, reveal unreachable handles or verbs, or return more than three calls. Check regex captures, required argument binding, root targeting, result/state binding restrictions, and semantic validation.

5. Pending frontier precedence: When a pending authored frontier exists, does orientation return it before attempting fresh activation? Can repeated orientation bypass, replace, or expand the pending frontier?

6. State safety: When a client copies a uniquely pending state-bound call but omits `state`, does admission restore the authored precondition? Identify any path that could convert a stale-safe call into an unconditional execution or incorrectly attach state to a different call.

7. State events: Are case events deterministic, indexed by executable acts, applied after result/frontier computation but before delivery, and confined to existing regular non-symlink files inside the sandbox? Check absolute paths, traversal, duplicate events, invalid indices, write failures, and event behavior after failed or refused acts. Decide whether full-file replacement is sufficiently specified for every experimental arm.

8. Evidence and accounting: Are orientations logged and costed while remaining separate from executable acts in trial evidence? Confirm that acceptable paths, sequence positions, act counts, frontier linkage, and grader checks use executable acts only. Check whether orientation failures, unmatched intents, repeated orientations, time, tokens, and cost have an unambiguous intention-to-treat treatment.

9. Arm comparability: For each of A, B, C, D, and E, state whether intent orientation and state events are present, absent, or undefined. Determine whether C versus B remains interpretable and whether C versus D and D versus E still isolate frontier and teaching-refusal effects. Report any place where orientation itself supplies a frontier or teaching signal that contaminates an isolation.

10. Information parity: Determine whether world-authored activation rules encode task-specific workflow knowledge unavailable to the flat-tool or static-documentation arms. If the information differs, decide whether the protocol specifies a fair counterpart for A and B and whether Arm B's document can be frozen without using validation or held_out information.

11. Corpus compatibility: Confirm that case-level state events and trial-level orientations are schema-validated, deterministic, and represented consistently in the runner, corpus loader, grader, and authoring manifest. Check that the grader digest changed and all authoring labels pin the current digest without weakening their precommitted checks.

12. Candidate retirement: Confirm that v4 keeps validation and held_out closed, marks v3 acceptance as historical only, and prevents normal tests from reading retired sealed candidates. Decide whether keeping the retired files in the implementation repository is compatible with the stated custody model.

13. Test adequacy: Identify v4 invariants asserted only in prose, implementation branches without a behavior owner, stale test names or assertions, and tests that could pass while the amendment contract is broken. Apply the repository standard: recommend extending the nearest existing behavioral test for a bug; recommend a new test only for genuinely new functionality.

14. Authoring result: First confirm that the retained summary and every referenced trial were produced by the candidate world and grader. Report the retained score rather than relying on the implementer's prose. The earlier v4 candidate retained 6/8 after a behavior-identical 7/8 tuning pass, with cascade and recovery failing. The submitted candidate may contain authoring corrections. Decide separately:
    a. whether the retained failures reveal a protocol or mechanism defect;
    b. whether allowing cascade to omit `tests.list` after `tests.run` names the failing test is a legitimate contract correction or outcome-driven label weakening;
    c. whether typed evidence connecting `TestSwitchyardHandshake` to `TestRelayHandshake` makes recovery well-specified without leaking a solution;
    d. whether the correction was precommitted before the retained rerun and applied consistently to comparable arms; and
    e. whether protocol v4 can be accepted while Gate 1A remains closed.

15. Freeze readiness: Verify that `status: authoring_amendment`, `frozen: false`, `review.accepted: false`, `may_open_validation: false`, and `may_open_held_out: false` accurately describe the submitted state. Determine the smallest changes required before v4 may become frozen, without authorizing any sealed outcome run.

Required response

Return:

1. Verdict: ACCEPT, REVISE, or REJECT.
2. Candidate identification: the commit and all four independently computed normalized digests.
3. Findings ordered P0 to P3. Each finding must include an exact file and line, consequence, and smallest correction. Say `No findings` if there are none.
4. A requirement matrix for all 15 review questions with PASS, FAIL, or UNRESOLVED and supporting evidence.
5. An arm-comparability matrix for A through E covering standing surface, orientation, frontier, teaching refusal, state events, and information available before the first executable act.
6. Independent standing-surface measurements for 10, 100, and 1,000 verbs, including serialized byte counts and spread.
7. A state-safety and event-ordering trace for one normal pending call and one call made stale by the authoring event.
8. An accounting trace showing how one orientation and its later executable act appear in the episode log, trial evidence, grading path, act count, and cost record.
9. A custody decision for the retired v3 validation and held_out candidates.
10. A separate decision on protocol-design acceptance, authoring readiness, and Gate 1A readiness. Do not collapse them into one verdict.
11. A decision on whether protocol.json may move to `status: frozen`, `frozen: true`, and `review.accepted: true` at this commit.
12. If ACCEPT, provide a review record with reviewer role, date, candidate commit, the four normalized digests, accepted limitations, unresolved authoring issues, and remaining execution blockers.
13. The exact next artifact allowed. Acceptance may authorize an acceptance-record patch or additional authoring work. It must not authorize opening validation or held_out until independent replacement candidates exist and every pre-validation artifact is frozen by digest.

Do not defer to the implementer's authoring analysis. Reconstruct each conclusion from the candidate artifacts. A documented external blocker is not automatically a protocol defect, but an undefined arm, accounting rule, or safety invariant is.
```
