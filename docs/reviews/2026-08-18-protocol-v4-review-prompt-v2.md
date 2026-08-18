# Protocol v4 second independent review prompt

Copy the text below into a fresh evaluator session that did not implement Phoenix and did not produce the first v4 review.

```text
Act as the second independent evaluation reviewer for the Phoenix frontier-v1 protocol v4 amendment. This is a protocol and implementation-contract review. It is not an outcome run, a corpus-authoring task, Gate 1A, or permission to open a sealed tranche.

Repository: https://github.com/nesste/phoenix
Candidate commit: f889f13f0c514fa5108a1e392701ebeadc4376f7

Identify the candidate before reviewing it. Compute SHA-256 over LF-normalized UTF-8 bytes and require these exact values:

- experiments/frontier-v1/protocol.json: sha256:4601e1970ebd271161fa3c5c3c245da28a5f55252eac8823f7f9d4f862adf0ff
- docs/decisions/0007-authoring-activation-amendment.md: sha256:1df4e012ec5d57e1759681cc92d77049897d1517d1ba3f109aa8222791be7674
- docs/decisions/0004-go-no-go-rules.md: sha256:cb35b48709321f885cfc5cd50a6bc08a51bf570d6146b5acad1bb283d3e41903
- experiments/frontier-v1/authoring-analysis.md: sha256:f28f07a7cf08583b45bd5e2dc043399c877288a9b7a68d95747ce9ddb15eb58b
- experiments/frontier-v1/results/authoring/summary.json: sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405

Return REVISE without reviewing the candidate if the commit or any digest differs.

Review context

The first independent v4 review returned REVISE with no P0 finding and accepted the core mechanism in principle: orientation, the O(1) standing surface, pending-frontier precedence, stale-state behavior, accounting, and the then-retained 8/8 authoring execution all checked out. It requested a focused revision:

1. make state events a shared harness rule for A through E, including flat A/B;
2. allow fresh intent activation only before the first executable act;
3. restore cascade to repo.status -> tests.run -> tests.list -> tests.focus;
4. make decision 0004's v3 freeze and acceptance unmistakably historical and document v4 orientation and cross-arm events;
5. scope omitted-state restoration to the current pending frontier;
6. reject mixed intent and executable fields inside Act as well as at MCP decoding;
7. pin per-arm system prompts and add the Phoenix intent instruction only to C/D/E; and
8. rename the stale protocol test to reflect that v4 is currently unfrozen.

The implementer reports that `make quality` passed, then authoring was run exactly once with no retry. That run scored 7/8: cascade failed after ignoring the returned tests.list frontier. Treat the historical 8/8 and the revision 7/8 as authoring evidence, not repeated measurements of a frozen design and not a Gate 1A result.

Review scope

Read these files from the exact candidate commit:

- experiments/frontier-v1/protocol.json
- docs/decisions/0004-go-no-go-rules.md
- docs/decisions/0007-authoring-activation-amendment.md
- experiments/frontier-v1/README.md and authoring-analysis.md
- experiments/frontier-v1/results/authoring/summary.json and the eight referenced trial and grade files
- experiments/frontier-v1/labels/authoring/authoring_1ca5cade.json
- experiments/frontier-v1/manifests/authoring.json
- internal/surface/mcp.go and mcp_test.go
- internal/activate/activate.go and activate_test.go
- internal/episode/log.go and sqlite.go
- cmd/phoenix/state_events.go and state_events_test.go
- experiments/frontier-v1/runner/main.go, runtime.go, types.go, and main_test.go
- experiments/frontier-v1/corpusctl/internal/corpus/protocol_test.go, trial.go, grade.go, and sealed_test.go
- experiments/frontier-v1/schema/case.schema.json and trial.schema.json
- spec/world.schema.json, spec/episode.schema.json, and spec/result.schema.json
- verbs/dev-repo/activations.go, transitions.go, refusals.go, and world_contract_test.go
- worlds/dev-repo/world.json and experiments/frontier-v1/worlds/authoring.dev_repo.json

Consult docs/plans/2026-08-17-phoenix-world-plan.md only to recover the original Task 0.5 and Gate 1A intent. Report conflicts with the v4 amendment; do not silently choose one.

Independence and custody rules

- Do not modify repository files.
- Do not run authoring, validation, or held_out trials.
- Do not open or inspect validation or held_out cases, fixtures, labels, manifests, registries, or outcomes. The v3 candidates are retired custody artifacts, not review inputs.
- Do not generate, relabel, reseal, or replace validation or held_out candidates.
- Read-only source inspection, digest computation, schema validation, and tests using temporary data are allowed.
- Do not treat the historical v3 ACCEPT record as acceptance of v4.
- Keep protocol-design acceptance, authoring readiness, Gate 1A readiness, and sealed-tranche authorization separate.

Allowed checks

1. `go test -count=1 ./...` from the repository root.
2. `go test -count=1 ./...` from `experiments/frontier-v1/corpusctl`.
3. `go run ./cmd/corpusctl validate --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json` from `experiments/frontier-v1/corpusctl`.
4. Independently serialize the MCP instructions and tool list for 10, 100, and 1,000 synthetic verbs and confirm O(1) size and the protocol limits.
5. Confirm the production and authoring world JSON files are byte-identical.

Required review questions

1. Input and activation boundary: Are executable and orientation modes mutually exclusive at both MCP decoding and direct admission? Does the first valid executable act permanently end fresh activation, including when post-act frontiers are suppressed? Can D or E recreate a suppressed call by reorienting? Can an orientation still return the current pending frontier or refusal alternative?

2. Pending-state scope: Can omitted state be restored only from a uniquely matching state-bound call in the current pending frontier? Check clearing, replacement, ambiguity, invalid execution, repeated orientation, refusal alternatives, and calls matching a retired frontier.

3. State-event parity: Does the protocol unambiguously make one deterministic executable-action-indexed event plan a shared harness rule for A, B, C, D, and E? Is the order result computation -> event -> delivery identical across arms, with orientations excluded from the index? Distinguish an adequate preimplementation harness contract from missing implementation that must block protocol acceptance.

4. Per-arm prompts: Are exact Phase 1 prompts pinned for A through E? Confirm that A/B receive no Phoenix intent instruction and C/D/E receive the same bootstrap instruction. Check that runner code uses the pinned C behavior for authoring and rejects an unpinned Phase 2 prompt.

5. Cascade restoration: Does the label again require repo.status -> tests.run -> tests.list -> tests.focus and bind its output check to the focused act? Confirm the authoring manifest is current. Decide whether restoration matches the original contract and whether the historical 8/8 path is consistent with it.

6. Revision run: Reconstruct the 7/8 result from summary, trials, labels, and grades. For cascade, verify that Phoenix returned tests.list after tests.run, the runtime ignored it, and the grade failed without retry. Decide whether this is mechanism evidence, runtime variance, authoring-readiness evidence, or a protocol-design blocker. Do not weaken the restored label or select the historical favorable run.

7. Historical status: Is every v3 freeze and ACCEPT statement in decision 0004 unmistakably historical? Are current v4 status, review state, and closed gates unambiguous in 0004, 0007, protocol.json, README, and analysis?

8. Evidence and accounting: Are orientations logged and costed but excluded from executable paths, sequence positions, act counts, event indices, and grader path checks? Confirm the absence case's four orientations and zero executable acts are represented consistently.

9. Test adequacy: Does `TestProtocolV4IsUnfrozenCompleteAndBudgeted` own the current unfrozen state and the new cross-arm/prompt contracts? Do the nearest surface tests own post-act non-reactivation, refusal-alternative precedence, current-frontier state expiry, and mixed direct admission plus MCP rejection? Identify any branch that can regress while these tests pass.

10. Constant surface and safety: Reconfirm the one-tool O(1) standing surface, reachable-handle-only activation, three-call bound, world-authored rule validation, stale-state behavior, and no registry search. The first review accepted these in principle; verify that the revision did not regress them.

11. Arm interpretation: Produce an A-E matrix for surface, prompt, orientation before/after first act, frontier, teaching refusal, state events, and pre-first-act information. Decide whether C-B, C-D, and D-E remain interpretable.

12. Freeze readiness: Verify `status: authoring_amendment`, `frozen: false`, `review.accepted: false`, `may_open_validation: false`, and `may_open_held_out: false`. State the smallest next patch if protocol design is accepted. Do not authorize Gate 1A or any sealed outcome run.

Required response

Return:

1. Verdict: ACCEPT, REVISE, or REJECT.
2. Candidate identification: commit plus all five independently computed normalized digests.
3. Findings ordered P0 to P3 with exact file and line, consequence, and smallest correction; say `No findings` if none.
4. A PASS/FAIL/UNRESOLVED matrix for all 12 questions with evidence.
5. The A-E comparability matrix requested above.
6. A normal frontier trace, a refusal-alternative trace, and a stale-event trace.
7. An accounting trace covering the revision-run absence case and one executable case.
8. Independent O(1) standing-surface measurements.
9. Separate decisions on protocol-design acceptance, authoring readiness, Gate 1A readiness, and sealed-tranche authorization.
10. A decision on whether protocol.json may move to frozen/accepted at this commit.
11. If ACCEPT, a review record with reviewer role, date, candidate commit, five digests, accepted limitations, unresolved authoring issues, and remaining execution blockers.
12. The exact next artifact allowed. It may be an acceptance-record patch or further visible authoring work; it must not be Gate 1A, candidate sealing, validation, or held_out execution.

Do not defer to the implementer's analysis. Reconstruct conclusions from candidate artifacts. A documented future harness obligation is not automatically a protocol defect, but an undefined cross-arm rule, accounting rule, or safety invariant is.
```
