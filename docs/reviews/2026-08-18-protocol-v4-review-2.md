# Protocol v4 independent acceptance record

- **Reviewer role:** Evaluation reviewer, independent of implementation
- **Date:** 2026-08-18
- **Verdict:** `ACCEPT`
- **Candidate commit:** `f889f13f0c514fa5108a1e392701ebeadc4376f7`
- **Findings:** None

## Candidate identity

The reviewer independently computed LF-normalized UTF-8 SHA-256 and matched all five pinned artifacts:

| Artifact | SHA-256 |
| --- | --- |
| `experiments/frontier-v1/protocol.json` | `sha256:4601e1970ebd271161fa3c5c3c245da28a5f55252eac8823f7f9d4f862adf0ff` |
| `docs/decisions/0007-authoring-activation-amendment.md` | `sha256:1df4e012ec5d57e1759681cc92d77049897d1517d1ba3f109aa8222791be7674` |
| `docs/decisions/0004-go-no-go-rules.md` | `sha256:cb35b48709321f885cfc5cd50a6bc08a51bf570d6146b5acad1bb283d3e41903` |
| `experiments/frontier-v1/authoring-analysis.md` | `sha256:f28f07a7cf08583b45bd5e2dc043399c877288a9b7a68d95747ce9ddb15eb58b` |
| `experiments/frontier-v1/results/authoring/summary.json` | `sha256:114fab6f1a1d2b6f8c98c3b4a8ef544e9aa22cf5f7e19b53c674409568435405` |

The reviewer ran the root and corpus test suites, validated the authoring manifest, confirmed the production and authoring worlds are byte-identical, and independently measured the constant standing surface at 535 compact JSON bytes with zero spread across 10, 100, and 1,000 verbs.

## Decision

The eight first-review revisions are present. There is no P0, P1, or P2 protocol-design defect. The retained 7/8 authoring run is runtime variance on the restored cascade contract, not a reason to weaken the label, select the historical favorable run, or reopen protocol design.

Protocol-design acceptance is granted. Authoring readiness remains descriptive, Gate 1A remains closed, and no sealed tranche is authorized.

## Accepted limitations

- Conclusions remain limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes.
- Headline power and efficiency-only conservatism are unchanged from the historical v3 record.
- A/B/D/E adapters, including shared state-event application on flat tools, remain preimplementation harness work.
- D/E suppression must remove frontier or refusal calls before pending state is recorded.
- Cascade remains run-to-run noisy. Its restored four-step label must not be relaxed, and the historical 8/8 must not be selected as a success estimate.

## Remaining execution blockers

- independent generation and sealing of new validation and held_out families;
- every `artifact_freeze.before_validation` digest;
- the Arm B document and human-factors review;
- Phase 1 A-E runner adapters, including shared state-event application;
- Task 0.6 Phase 0 review.

The exact authorized next artifact was this acceptance-record patch. It may set protocol v4 to accepted and frozen while keeping `may_open_validation` and `may_open_held_out` false. It does not authorize Gate 1A, candidate sealing, validation execution, or held_out execution.
