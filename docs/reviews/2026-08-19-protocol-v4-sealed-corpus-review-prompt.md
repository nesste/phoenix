# Protocol-v4 sealed-corpus independent review prompt

Status: superseded by `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review-prompt-v2.md`. Do not issue this version to a reviewer.

Give this prompt to an evaluation reviewer who did not author the replacement corpus or labels, prepare the candidate package, or implement Phoenix. Use a fresh review session.

```text
Independently review the unopened protocol-v4 replacement candidates for the `validation` and `held_out` tranches of `frontier-v1`.

This is a corpus-and-label gate review. It is not an import, schedule-generation step, model run, arm run, trial, validation run, held-out run, or outcome analysis.

Candidate handoff

- Read-only review checkout: D:\Work\personal\phoenix-evaluator-v4-review
- Evaluator base: c852101e8d7cb52e4569bf3866de54a0ce648b44
- Candidate commit: c193c786cc5a65ce6ae97efbd336b2a48f492898
- Implementation freeze: 4b51471200ed577db55fa38bed61c027c276d0f2
- Public patch: D:\Work\personal\phoenix-evaluator-private-v4\frontier-v1-v4-sealed-corpus.patch
- Expected public-patch SHA-256: sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65
- Public report in the review checkout: docs/reviews/2026-08-19-frontier-v1-v4-sealed-corpus-report.md
- Expected private-label archive digest: sha256:21288e5c2e808e5cb8528d5ed6999ebd69db3e54a7ef8028aae85531b19b4b72
- Expected private-label custodian: evaluator.v4.private_label_custodian
- Expected retired-v3 collision-inventory digest: sha256:c12cc875b6630b4ed6327184b96773d643ceb71fbc1dead3ef3f6796a4d2dc9a
- Expected collision-inventory custodian: evaluator.v4.private_label_custodian

Read:

- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md
- docs/reviews/2026-08-19-frontier-v1-v4-sealed-corpus-report.md
- docs/plans/2026-08-17-phoenix-world-plan.md, especially Task 0.4
- experiments/frontier-v1/protocol.json
- experiments/frontier-v1/schema/
- spec/world.schema.json
- experiments/frontier-v1/worlds/authoring.dev_repo.json
- experiments/frontier-v1/corpusctl/
- the complete public diff from the evaluator base to the candidate
- all public validation and held-out cases, fixtures, label-digest registries, and sealed manifests at the candidate commit
- through the custodian's private channel, all 240 full replacement labels and the outcome-free collision inventory

Independence and custody rules

- Record your reviewer role, date, candidate commit, material seen, and independence limitations. You must be independent of the corpus author, first-pass labeler, candidate-package preparer, and Phoenix implementation.
- Keep the candidate checkout detached and read-only. Do not amend, rebase, or modify the candidate branch.
- Obtain the private archive and collision inventory from the named custodian and verify their identities before using them. Do not assume an unverified directory is the reviewed archive.
- Do not copy a full label, expected answer, acceptable path, rationale, private template note, disagreement detail, or adjudication detail into a Phoenix checkout, Git history, command log, review record, or response.
- Findings about private material may name an opaque case ID and defect category, but must not disclose the expected outcome or enough detail to reconstruct it.
- Do not inspect the historical protocol-v3 private-label archive.
- Do not run an experiment runner, model, agent, arm, trial, schedule, prospective grader outcome, validation outcome, or held-out outcome. Do not generate a schedule or open Gate 1A.
- Do not import the public patch into the implementation workspace during this review.
- Read-only source inspection, hashing, schema validation, Go tests, Go vet, label canonicalization, allocation checks, collision analysis, and outcome-free manifest reproduction are allowed.
- Stop and report a custody breach if any prospective validation or held-out outcome is observed. A consumed tranche cannot be accepted by silently regenerating or relabeling it.

Verify all of the following.

1. Identity and boundary
   - The review checkout is clean, detached at the exact candidate commit, and descends from the exact evaluator base.
   - The raw public-patch SHA-256 matches the handoff and the patch represents the base-to-candidate public replacement.
   - The candidate diff contains only deletion of retired-v3 public tranche files, new protocol-v4 public cases and fixtures at the approved paths, the two outcome-free label-digest registries, the two sealed manifests, and the public report.
   - No protocol, runtime/prompt, arm schema, Arm B, world, runner, grader, analysis, accepted freeze entry, gate, schedule, implementation source, or full-label byte changed or entered the patch.
   - `pre-validation-artifacts.json` remains `partial`; only `schedule digest` remains; both outcome gates remain false.

2. Pinned contract and reproducibility
   - Independently recompute the LF-normalized protocol digest, canonical and raw world identities, grader digest, and frozen Linux-amd64 world-build identity recorded in the evaluator handoff. Stop on a mismatch.
   - Run `go test ./...` and `go vet ./...` in `experiments/frontier-v1/corpusctl`.
   - Run `corpusctl grader-digest` and both `seal` commands without `--write`. Compare the generated validation and held-out manifests byte-for-byte with the committed manifests.
   - Keep the review checkout clean after every check. Do not use `seal --write` in the review checkout.

3. Counts and allocation, independently from the report
   - Each tranche has exactly 120 cases in exactly 24 generating families, five cases per family, 120 registry rows, and no missing or extra label identity.
   - Each tranche has eight direct-group families, eight recovery-group families, and eight mix families.
   - Each tranche has 24 `direct`, 24 `recovery`, and 12 of every other required class.
   - Every required class meets its distinct-family floor. Every mix family includes at least one `cascade` and one `adversarial_text` case, and the stated mix arithmetic is exact.
   - Every family's five variants differ in observable state or task; a filename or prose-only rewrite is not counted as a distinct generating variant.

4. Structural independence and public opacity
   - Independently compare family IDs, case IDs, goals, generating templates, fixture contents, file-map digests, answer-path structures, service aliases, test catalogs, adversarial patterns, and state-change patterns against authoring, retired-v3 validation, retired-v3 held-out, and the other replacement tranche.
   - Verify the external collision inventory's digest and completeness against the tracked collision sets. Do not accept the report's disjointness claim without checking the inventory and candidate.
   - Reject renamed or lightly rewritten prior families as non-independent.
   - Inspect every public case and fixture for leakage of class, expected answer, acceptable path, grader logic, private rationale, or instructions that reveal how to score the case.
   - Confirm adversarial fixture text is untrusted task data rather than evaluator or runner instruction.

5. Protocol-v4 executable semantics
   - Every case uses the current strict schema and canonical world digest, references only available fixtures, and is solvable through the frozen world and its twelve flat-tool equivalents.
   - Every `state_changes` replacement is deterministic, uses a safe relative path, and occurs after a numbered executable act.
   - The same event plan and executable-action index apply to Arms A-E; orientations never advance the action index.
   - Every `stale_frontier` case creates staleness through a declared shared state event, not timing, an external race, or arm-specific behavior.
   - Capability-absence cases have a genuinely absent capability and do not reward unsupported action. Cascade cases preserve the frozen cascade contract. No case relies on an unregistered capability or a host path outside its fixture.

6. Complete private-label review
   - Inspect all 240 labels, not a sample. Verify one schema-valid full label per public case and no extra label.
   - Canonicalize every label with the candidate's `corpusctl digest`; each digest must match the corresponding public registry row. Recompute the private archive identity using the report's stated algorithm and match the handoff digest.
   - Verify every label uses the same opaque `case_id`, the correct evaluator-only class, and grader sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc.
   - Reason from the frozen world, public case, fixture, and declared state events to verify the expected outcome. Do not use a model or prospective run to test it.
   - Verify checks establish success and distinguish at least one plausible wrong path. Check regex anchoring and escaping, file paths, status expectations, act counts, action indexes, and state-event timing.
   - Verify action indexes count executable acts only; orientations are separate evidence and never advance an index.
   - Verify acceptable paths are non-exhaustive and contain an empty path when no action is expected. Verify each rationale explains the expected outcome and why the checks prove it.
   - Verify deterministic checks are used where possible. If any human judgment is required, require `arm_hidden: true`, a concrete outcome-independent rubric, and retained independent adjudication.
   - Independently recount scripted, regex, and blinded-human checks by kind and compare with the public report. Reconcile every difference.

7. Independence and custody evidence
   - Confirm no Phoenix arm output, runner result, retained authoring outcome, prospective sealed outcome, or implementation suggestion influenced family design or labels.
   - Confirm no prospective model, agent, arm, trial, schedule, validation, held-out, or outcome was run or observed.
   - Confirm full labels exist only in the new protocol-v4 private archive, not the evaluator commit, public patch, report, implementation checkout, or historical v3 archive.
   - Confirm disagreements and adjudications are complete and privately retained. Record counts publicly without private content.
   - Treat the shared corpus-author/first-pass-labeler role and mechanical original audit as an explicit independence limitation. Your review must supply the missing independent label-content assessment.

Review record and verdict

Write the public review record to:

D:\Work\personal\phoenix\docs\reviews\2026-08-19-protocol-v4-sealed-corpus-review.md

Return `ACCEPT`, `REVISE`, or `REJECT` with P0-P3 findings and exact public paths/lines where applicable. For private findings, use only opaque case ID plus a non-revealing defect category.

The review record must include:

- reviewer role, independence declaration, date, evaluator base, candidate commit, and reviewed locations;
- public-patch, private-archive, collision-inventory, grader, world, protocol, world-build, registry, and manifest identities independently recomputed;
- a requirement matrix for all sections above;
- independently computed count/allocation and check-kind totals;
- test, vet, grader-digest, seal-reproduction, archive-digest, and collision checks with exit status;
- findings, accepted limitations, custody status, and confirmation that no outcome was opened;
- the smallest next artifact authorized by the verdict.

If `ACCEPT`, authorize only a focused import of the exact reviewed public patch into the implementation workspace, followed by verification that protected freeze bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates are unchanged. Do not authorize schedule generation, Gate 1A, validation, held-out, or any outcome run.

If `REVISE`, keep the candidate unimported and name the smallest replacement patch, private-label correction, or custody artifact needed. Any private-label correction must produce new affected label digests and a new private-archive identity; determine whether dependent public registries, manifests, report, candidate commit, and patch digest must also be replaced.

If `REJECT`, identify the failed independence, validity, or custody premise and keep all gates closed.
```
