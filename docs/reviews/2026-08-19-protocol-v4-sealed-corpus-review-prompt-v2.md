# Protocol-v4 sealed-corpus independent review prompt v2

This prompt supersedes `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review-prompt.md`. Give it to an evaluation reviewer who did not author the replacement corpus or labels, prepare the candidate package, implement Phoenix, or prepare this review handoff. Use a fresh review session.

```text
Independently review the unopened protocol-v4 replacement candidates for the `validation` and `held_out` tranches of `frontier-v1`.

This is a corpus-and-label gate review. It is not an import, schedule-generation step, model run, arm run, trial, validation run, held-out run, or outcome analysis.

Candidate handoff

- Read-only review checkout: D:\Work\personal\phoenix-evaluator-v4-review
- Evaluator base: c852101e8d7cb52e4569bf3866de54a0ce648b44
- Candidate commit: c193c786cc5a65ce6ae97efbd336b2a48f492898
- Implementation freeze: 4b51471200ed577db55fa38bed61c027c276d0f2
- Public patch: D:\Work\personal\phoenix-evaluator-private-v4\frontier-v1-v4-sealed-corpus.patch
- Expected public-patch raw SHA-256: sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65
- Public report in the review checkout: docs/reviews/2026-08-19-frontier-v1-v4-sealed-corpus-report.md
- Expected private-label archive digest: sha256:21288e5c2e808e5cb8528d5ed6999ebd69db3e54a7ef8028aae85531b19b4b72
- Expected private-label custodian: evaluator.v4.private_label_custodian
- Expected retired-v3 collision-inventory raw SHA-256: sha256:c12cc875b6630b4ed6327184b96773d643ceb71fbc1dead3ef3f6796a4d2dc9a
- Expected collision-inventory custodian: evaluator.v4.private_label_custodian

Pinned freeze identities

- Protocol v4, LF-normalized UTF-8: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9
- Production/authoring world, canonical JSON: sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b
- Production/authoring world, raw bytes: sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143
- Current grader: sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc
- Frozen Linux-amd64 content-addressed world-build digest, recorded at `artifacts.world_definition_and_world_build_digest.world_build_digest` in `pre-validation-artifacts.json`: sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d
- World-build manifest raw bytes, recorded separately at `artifacts.world_definition_and_world_build_digest.world_build_manifest_raw_sha256`: sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f

Read:

- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md
- docs/reviews/2026-08-19-frontier-v1-v4-sealed-corpus-report.md
- docs/plans/2026-08-17-phoenix-world-plan.md, especially Task 0.4
- experiments/frontier-v1/protocol.json, especially `tranche_design.family_allocation`
- experiments/frontier-v1/schema/
- spec/world.schema.json
- experiments/frontier-v1/worlds/authoring.dev_repo.json
- experiments/frontier-v1/corpusctl/
- experiments/frontier-v1/corpus/authoring/
- experiments/frontier-v1/fixtures/authoring/
- experiments/frontier-v1/labels/authoring/
- the complete public diff from the evaluator base to the candidate
- the retired-v3 public validation and held-out cases, fixtures, registries, and manifests from the evaluator-base tree, using `git show c852101e8d7cb52e4569bf3866de54a0ce648b44:<path>` or equivalent read-only Git-object access
- all replacement public validation and held-out cases, fixtures, label-digest registries, and sealed manifests at the candidate commit
- through the custodian's private channel, all 240 new full replacement labels, new private template notes and adjudication records, and the outcome-free collision inventory

Independence and custody rules

- Record your reviewer role, date, candidate commit, material seen, and independence limitations. You must be independent of the corpus author, first-pass labeler, candidate-package preparer, Phoenix implementation, and this prompt revision.
- Keep the candidate checkout detached and read-only except for the one LF-materialization procedure below. Do not amend, rebase, switch, or modify the candidate branch.
- Obtain the new private archive and collision inventory from the named custodian and verify their identities before using them. Do not assume an unverified directory is the reviewed archive.
- Do not copy a full label, expected answer, acceptable path, rationale, private template note, disagreement detail, or adjudication detail into a Phoenix checkout, Git history, command log, review record, or response.
- Findings about private material may name an opaque case ID and defect category, but must not disclose the expected outcome or enough detail to reconstruct it.
- Do not inspect, hash, mount, or otherwise open the historical protocol-v3 private-label archive. All retired-v3 checks use only public files from the evaluator-base Git tree and the new outcome-free collision inventory.
- Do not run an experiment runner, model, agent, arm, trial, schedule, prospective grader outcome, validation outcome, or held-out outcome. Do not generate a schedule or open Gate 1A.
- Do not import the public patch into the implementation workspace during this review.
- Read-only source inspection, hashing, schema validation, Go tests, Go vet, label canonicalization, allocation checks, collision analysis, and outcome-free manifest reproduction are allowed.
- Stop and report a custody breach if any prospective validation or held-out outcome is observed. A consumed tranche cannot be accepted by silently regenerating or relabeling it.

Windows LF materialization

Before any byte-identity, digest, or seal check, confirm `git status --short` is empty in the detached review checkout, then run exactly:

    git config --worktree core.autocrlf false
    git config --worktree core.eol lf
    git checkout-index --all --force
    git status --short

Run `checkout-index` only in the new clean detached review checkout. Stop if status was nonempty before it or is nonempty after it. This is the only permitted materialization change. Never use `seal --write` in the review checkout.

Verify all of the following.

1. Identity, ancestry, patch equality, and boundary
   - The review checkout is clean, detached at the exact candidate commit, descends from the exact evaluator base, and contains the implementation freeze as an ancestor.
   - Recompute the raw patch SHA-256 and match the handoff.
   - Independently generate the raw stdout bytes of `git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 c193c786cc5a65ce6ae97efbd336b2a48f492898`. Capture native stdout with a binary-safe process API or raw stream and hash those bytes directly. Do not pass the diff through a PowerShell pipeline, `>`, `Out-File`, `Set-Content`, or any text decode/re-encode step. Require the raw diff bytes to hash to `sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65` and equal the public patch byte-for-byte, not merely apply to the same tree.
   - The candidate diff contains only deletion of retired-v3 public tranche files, new protocol-v4 public cases and fixtures at the approved live paths, the two outcome-free label-digest registries, the two sealed manifests, and the public report.
   - Reuse of the live repository paths under `corpus/{validation,held_out}` and `fixtures/{validation,held_out}` is expected replacement behavior. Do not count a reused live path by itself as an identity or derivation collision.
   - No protocol, runtime/prompt, arm schema, Arm B, world, runner, grader, analysis, accepted freeze entry, gate, schedule, implementation source, or full-label byte changed or entered the patch.
   - `pre-validation-artifacts.json` remains `partial`; only `schedule digest` remains; both outcome gates remain false.

2. Pinned contract and reproducibility
   - Independently recompute the protocol's LF-normalized digest, both world identities, the grader digest, the raw world-build-manifest digest, and the manifest's internal content-addressed `world_build_digest`. Use the algorithm attached to each field; do not obtain `sha256:27c2f537...` by hashing the manifest file. Stop before private-label review if any value differs.
   - Before inspecting any private label, run `go run ./cmd/corpusctl grader-digest --repo-root ../../..` and require exact output `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`.
   - Run `go test -count=1 ./...` and `go vet ./...` in `experiments/frontier-v1/corpusctl`.
   - Run these exact dry-seal commands from `experiments/frontier-v1/corpusctl`; do not add `--write`:

       go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json
       go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json

     Compare each generated manifest byte-for-byte with its committed manifest.
   - Keep the review checkout clean after every check.

3. Frozen allocation, independently from the candidate report
   - Use `protocol.json` `tranche_design.family_allocation` as the source of truth and public registry classes as the counted evidence. Do not derive the contract from the candidate report.
   - Each tranche has exactly 120 cases in exactly 24 generating families, five cases per family, 120 registry rows, and no missing or extra label identity.
   - `group_direct` has exactly eight families, and every one independently has exactly 3 `direct` + 1 `absence` + 1 `stale_frontier` case.
   - `group_recovery` has exactly eight families, and every one independently has exactly 3 `recovery` + 1 `far_discovery` + 1 `temptation` case.
   - `group_mix` has exactly eight families and, across those eight families, exactly 12 `cascade`, 12 `adversarial_text`, and 4 each of `far_discovery`, `temptation`, `absence`, and `stale_frontier` cases. Every mix family has five cases, at least one `cascade`, and at least one `adversarial_text`.
   - Consequently, each tranche has exactly 24 `direct`, 24 `recovery`, and 12 of every other required class. Every class meets its frozen distinct-family floor.
   - Every family's five variants differ in observable state or task. Distinct goal strings alone do not prove five generating variants; do not accept the report's five-goals evidence as sufficient.

4. Collision procedure and structural independence
   - Perform four separate public comparisons: replacement versus current authoring; replacement versus retired-v3 public validation at the evaluator base; replacement versus retired-v3 public held-out at the evaluator base; and replacement validation versus replacement held_out.
   - From public Git trees, compare family IDs, case IDs, goals, fixture contents, file-map digests, service aliases, test catalogs, adversarial patterns, and state-change patterns. Public case documents have no generating-template field: infer generating-template reuse from observable goal grammar and the combined fixture, alias, catalog, adversarial, and state-change structure. Use `git show` or equivalent object reads for v3-at-base; do not switch the review checkout to the base.
   - Verify the external outcome-free inventory's completeness against the authoring and both retired-v3 public trees for exactly the fields it claims: tracked paths, family IDs, case IDs, goals, fixture digests, and file-map digests.
   - Verify `collision-inventory.json` by raw-file SHA-256, with no normalization or JSON canonicalization, and require `sha256:c12cc875b6630b4ed6327184b96773d643ceb71fbc1dead3ef3f6796a4d2dc9a`.
   - Do not ask the collision inventory to prove fields it does not contain. Public v3 observable-structure, alias, catalog, adversarial, and state-change comparisons come from the evaluator-base Git tree itself.
   - Compare label-only structures—acceptable paths, rationales, grading logic, and expected-state derivations—only between the visible authoring labels and the new 240 private labels, and between the two new replacement tranches. Compare the new private template notes with public authoring family structures and with the other replacement tranche. Never use the historical v3 private archive for either check.
   - Reject renamed or lightly rewritten prior public families as non-independent. Distinguish expected live-path replacement from reuse of an ID, generating template, fixture content, goal grammar, answer-path structure, service alias, test catalog, adversarial pattern, or state-change pattern.

5. Public opacity
   - Inspect every public replacement case and fixture for leakage of class, expected answer, acceptable path, grader logic, private rationale, or instructions that reveal how to score the case.
   - Confirm adversarial fixture text is untrusted task data rather than evaluator or runner instruction.
   - Require every row of each public label-digest registry to contain only `case_id`, `class`, `label_digest`, and `grading_script`, with no additional outcome-bearing field.
   - Inspect both sealed manifests and the entire public report for expected outcomes, acceptable paths, rationales, template notes, private disagreements, adjudication details, or other answer leakage.

6. Protocol-v4 executable semantics
   - Every case uses the current strict schema and canonical world digest, references only available fixtures, and is solvable through the frozen world and its twelve flat-tool equivalents.
   - Every `state_changes` replacement is deterministic, uses a safe relative path, and occurs after a numbered executable act.
   - Treat A-E event-index parity as a case-contract check: every arm receives the same declared event plan and executable-action index, and orientations never advance that index. This is not authorization to re-review or execute the runner.
   - Every `stale_frontier` case creates staleness through a declared shared state event, not timing, an external race, or arm-specific behavior.
   - Capability-absence cases have a genuinely absent capability and do not reward unsupported action. No case relies on an unregistered capability or a host path outside its fixture.
   - Require every replacement cascade to satisfy the Task 0.4 class definition: success needs at least three verbs in an order the agent would not guess upfront. Treat the authoring sequence `repo.status` -> `tests.run` -> `tests.list` -> `tests.focus` as a collision set, not as the replacement template. Reject replacement cases that clone its verb sequence, test-name pattern, or focused-act binding. If a distinct replacement cascade uses the tests diagnostic workflow, it must not omit `tests.list` between `tests.run` and `tests.focus`; do not select the favorable authoring omission. Otherwise judge the case's own cascade and bind its checks to the act that establishes success rather than imposing the authoring verbs.

7. Complete private-label review
   - Inspect all 240 new labels, not a sample. Verify one schema-valid full label per public case and no extra label.
   - Canonicalize every label with the candidate's `corpusctl digest`; each digest must match the corresponding public registry row.
   - Recompute the private archive identity using this fixed algorithm: SHA-256 of LF-normalized UTF-8 lines `<archive-relative-path>\t<sha256:file-bytes>` for all 240 label files, sorted by archive-relative path. Require `sha256:21288e5c2e808e5cb8528d5ed6999ebd69db3e54a7ef8028aae85531b19b4b72`. Do not substitute a self-described algorithm merely because the candidate report states it.
   - Verify every label uses the same opaque `case_id`, the correct evaluator-only class, and grader `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`.
   - Reason from the frozen world, public case, fixture, and declared state events to verify the expected outcome. Do not use a model or prospective run to test it.
   - Verify checks establish success and distinguish at least one plausible wrong path. Check regex anchoring and escaping, file paths, status expectations, act counts, action indexes, and state-event timing.
   - Verify action indexes count executable acts only; orientations are separate evidence and never advance an index.
   - Verify acceptable paths are non-exhaustive and contain an empty path when no action is expected. Verify each rationale explains the expected outcome and why the checks prove it.
   - Verify deterministic checks are used where possible. If any human judgment is required, require `arm_hidden: true`, a concrete outcome-independent rubric, and retained independent adjudication.
   - Independently recount scripted, regex, and blinded-human checks by kind and compare with the public report. Reconcile every difference.

8. Independence and custody evidence
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
- public-patch, private-archive, collision-inventory, grader, world, protocol, world-build, registry, and manifest identities independently recomputed with the pinned algorithms;
- a requirement matrix for all sections above;
- independently computed per-family allocation, class/family totals, variant evidence, and check-kind totals;
- test, vet, grader-digest, dry-seal reproduction, patch-equality, archive-digest, and collision checks with exit status;
- findings, accepted limitations, custody status, and confirmation that no outcome was opened;
- the smallest next artifact recommended by the verdict.

Verdict bar:

- `ACCEPT` requires no unresolved P0 or P1 finding. Every P2 must be resolved or explicitly justified as nonblocking; otherwise return `REVISE`.
- An identity-pin failure, private-label custody breach, observed prospective outcome, or consumed tranche is P0 and requires `REJECT`, not a patchable `REVISE`.
- Use `REVISE` for a correctable corpus, label, registry, manifest, report, patch, or review-evidence defect while custody remains intact.
- Use `REJECT` for a failed independence, identity, validity, or custody premise that invalidates the candidate.

If `ACCEPT`, recommend only that the project chair consider a focused import of the exact reviewed public patch. The reviewer does not authorize import. The chair must record the import decision separately and, if importing, verify that protected freeze bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates remain unchanged. Neither reviewer acceptance nor chair import authorizes schedule generation, Gate 1A, validation, held-out, or any outcome run.

If `REVISE`, keep the candidate unimported and name the smallest replacement patch, private-label correction, or custody artifact needed. Any private-label correction must produce new affected label digests and a new private-archive identity; determine whether dependent public registries, manifests, report, candidate commit, and patch digest must also be replaced.

If `REJECT`, identify the failed independence, identity, validity, or custody premise and keep all gates closed.
```
