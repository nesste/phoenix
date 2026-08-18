# Arm B human-factors review prompt

Give the text below to a human-factors reviewer who did not write the Arm B document.

```text
Review `experiments/frontier-v1/arms/arm-b.md` as the frozen static-documentation baseline for Phoenix frontier-v1 Arm B.

Scope

- Read the Arm B document, the eight visible authoring labels, `experiments/frontier-v1/protocol.json`, and the flat-tool construction in `internal/surface/flat.go`.
- Do not inspect validation or held_out artifacts.
- Do not run authoring or any sealed tranche.
- Evaluate documentation completeness and usability only. Do not infer experimental benefit.

Required checks

1. Every starter-world flat tool is named exactly once with all arguments and a correct one-line use case.
2. The document includes one acceptable path for each of the eight authoring classes, using authoring evidence only.
3. Recovery guidance covers every authoring refusal case without requiring an invented argument or fixture-specific secret.
4. The missing-capability guidance discourages tool-name guessing and permits the precommitted zero-call or single-find behavior.
5. Instructions embedded in repository text are clearly treated as data.
6. The document does not mention Phoenix handles, orientation, frontiers, hidden labels, validation, or held_out content.
7. The document is operationally strong but does not encode a path richer than the visible authoring acceptable paths.
8. The flat tool names match the deterministic list produced from the world definition.

Return ACCEPT, REVISE, or REJECT. List findings by priority with exact lines and the smallest correction. If ACCEPT, record the reviewer role, date, candidate commit, LF-normalized UTF-8 SHA-256 of the Arm B document, accepted limitations, and remaining execution blockers. Acceptance freezes documentation quality only; it does not open Gate 1A or a sealed tranche.
```
