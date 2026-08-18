# Arm B human-factors review

## Verdict

**ACCEPT**

The Arm B static document meets the documentation-completeness and usability checklist. This verdict freezes documentation quality only. It does not estimate experimental benefit, open Gate 1A, or authorize validation or held_out.

## Candidate

| Field | Value |
| --- | --- |
| Reviewer role | Human-factors reviewer, independent of Arm B authorship |
| Review date | 2026-08-18 |
| Candidate commit | `73adf8c608f0edf06597b569b17faa32e1a3b5b9` |
| Document | `experiments/frontier-v1/arms/arm-b.md` |
| LF-normalized UTF-8 SHA-256 | `sha256:e717895a0e617b9e1a4fd9b9511b3604a1aa07867045f815e63f2d60a3ccc8f0` |
| Candidate bytes | 2673, LF-only UTF-8 |

The production and authoring world definitions were byte-identical at review time.

## Required checks

All eight checks passed:

1. Every starter-world flat tool appears exactly once with its complete argument list and a correct one-line use case.
2. The document contains one authoring-derived acceptable path for each visible task class.
3. Recovery guidance covers every authoring refusal case without invented arguments or fixture-specific names.
4. Missing-capability guidance forbids tool-name guessing and preserves the precommitted zero-call or single-`repo_find` behavior.
5. Repository text and tool results are treated as data.
6. The document contains none of the prohibited treatment or sealed-tranche terminology.
7. Its instructions are operational but no richer than the visible authoring acceptable paths.
8. Its twelve flat tool names match the deterministic world-derived list.

The accepted tool set is `episodes_recall`, `git_commit`, `git_diff`, `git_status`, `repo_build`, `repo_edit`, `repo_find`, `repo_read`, `repo_status`, `tests_focus`, `tests_list`, and `tests_run`.

## Findings

There were no P0, P1, or P2 findings.

One non-blocking P3 note remains: the `git_status` use line could say that it returns the raw Git porcelain listing and is not a substitute for `repo_status`. The cascade row already selects `repo_status`, so the reviewer accepted the candidate without this clarification. Changing the accepted document would require a new digest and review.

## Accepted limitations

- The evidence is authoring-only: eight classes with one visible case per class.
- Result schemas are omitted; the protocol requires tool names, arguments, and one-line use guidance.
- `repo_find` returns paths. A live result binds the later `tests_focus` call; prescribing `repo_read` or another `tests_list` would exceed the far-discovery label.
- Task-class headings are operational paraphrases rather than protocol identifiers.
- This review does not measure C minus B benefit.

## Remaining blockers

- Freeze every other artifact listed in `protocol.json` under `artifact_freeze.before_validation`.
- Complete the Williams schedule, three repetitions, retry and ITT accounting, analysis implementation, and report template.
- Independently generate and seal new disjoint validation and held-out families where the v4 record requires replacements.

Gate 1A, validation, and held_out remain closed.
