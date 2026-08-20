# Protocol-v4 sealed-corpus revision independent review prompt

This prompt is for the two-commit replacement package produced from `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md`. It supersedes neither historical record, but it is the only review assignment that matches the revision package. Do not reuse `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review-prompt-v2.md`; that prompt identifies the rejected one-commit candidate.

Give the text below to a reviewer who did not author the replacement corpus, labels, blueprints, package, report, private audits, or this prompt; did not perform the blueprint or label audits; and did not implement Phoenix. Use a fresh review session and a new detached review checkout.

```text
Independently review the protocol-v4 replacement `validation` and `held_out` corpus package for `frontier-v1`.

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, validation, held-out evaluation, or outcome analysis.

Candidate package

- Fresh detached review checkout: D:\Work\personal\phoenix-evaluator-v4-revision-review
- Evaluator branch: codex/evaluator-v4-corpus-revision
- Evaluator base: c852101e8d7cb52e4569bf3866de54a0ce648b44
- Payload candidate commit: e705560fef6012f84ec0199b87ce790a62c17d0e
- Report-only commit: 7a57479bc64cdb882bf17bfff7ed5d80c2f64760
- Required ancestry: report-only commit parent == payload candidate; payload candidate parent == evaluator base
- Implementation freeze: 4b51471200ed577db55fa38bed61c027c276d0f2
- Rejected v4 public candidate, collision use only: c193c786cc5a65ce6ae97efbd336b2a48f492898
- Public payload patch: D:\Work\personal\phoenix-evaluator-private-v4-revision\frontier-v1-v4-sealed-corpus-revision.patch
- Expected public-patch raw SHA-256: sha256:96d25bd5b0340d71b71d46b8d030bd746323e29f224ffe6ba4c275f7765e4ded
- Expected public-patch byte length: 1366164
- Public report at the report-only commit: docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-revision-report.md
- Private custody root: D:\Work\personal\phoenix-evaluator-private-v4-revision
- Expected private-label archive digest: sha256:c66d73e6d3db16b52feac05b6ceacc5ffb441b6c6ab4fec8705a87ea19d076bc
- Expected private-label custodian: revision.private_label_custodian
- Expected collision-inventory raw SHA-256: sha256:af6c8d1a1630efb54c404aebf4410448ad5c324075b30413c6c12a0ce94e888f
- Expected private blueprint raw SHA-256: sha256:0c4f2ff6db43128a604e8c824e5ae07acb76f410c4ffdf093d0ed7ebf76d66f0
- Expected independent blueprint-audit raw SHA-256: sha256:ac776a68cb1ad6b15fe82fe0115f2213b383bc6d44fade651872d41928bb4cbc
- Expected independent label-audit raw SHA-256: sha256:ae3fccb58b4dc7b2e2504f7a6ae3995ce98b309cc8dd879a037467cacb76cf81
- Expected private mechanical-audit raw SHA-256: sha256:25a1ac529a49d29cbd8aa7048919e235570a4c06571cce7a366bef2308b188bd

Pinned freeze identities

- Protocol v4, LF-normalized UTF-8: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9
- Production/authoring world, canonical JSON: sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b
- Production/authoring world, raw bytes: sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143
- Current grader: sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc
- World-build manifest, raw bytes: sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f
- Content-addressed `world_build_digest`: sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d

Required reading

- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md
- docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-revision-report.md at the report-only commit
- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md
- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review.md, as the public defect record for the rejected candidate
- docs/plans/2026-08-17-phoenix-world-plan.md, especially Task 0.4
- experiments/frontier-v1/protocol.json, especially `tranche_design.family_allocation`
- experiments/frontier-v1/schema/case.schema.json
- experiments/frontier-v1/schema/label.schema.json
- experiments/frontier-v1/schema/common.schema.json
- experiments/frontier-v1/corpusctl/internal/corpus/grade.go and the relevant `corpusctl` commands
- spec/world.schema.json
- experiments/frontier-v1/worlds/authoring.dev_repo.json
- worlds/dev-repo/world.json
- experiments/frontier-v1/corpus/authoring/
- experiments/frontier-v1/fixtures/authoring/
- experiments/frontier-v1/labels/authoring/, for collision and opacity checks only
- the complete evaluator-base-to-payload diff
- the report-only commit diff
- retired-v3 public validation and held-out material at the evaluator base, using read-only Git-object access
- rejected-v4 public validation and held-out material at commit c193c786cc5a65ce6ae97efbd336b2a48f492898, using read-only Git-object access
- every replacement public case, fixture, label-digest registry, and sealed manifest at the payload commit
- through the private custodian, all 240 full replacement labels, the blueprint file, blueprint-audit record, label-audit record, mechanical-audit record, stored archive index, and collision inventory

Independence and custody

- Record the reviewer role, date, material seen, and all independence limitations. The reviewer must be distinct from `revision.corpus_author`, `revision.first_pass_labeler`, `revision.independent_blueprint_auditor`, `revision.independent_label_auditor`, `revision.private_label_custodian`, the candidate-package preparer, Phoenix implementers, and the author of this prompt.
- Treat the prior private audits as evidence to verify, not as substitutes for this review. Independently inspect all 48 family designs and all 240 full labels.
- Keep the review checkout detached and read-only except for the LF materialization procedure below. Do not amend, rebase, switch, commit, or modify the evaluator branch.
- Obtain private artifacts only from `revision.private_label_custodian`. Verify every pinned private identity before relying on its content.
- Do not copy a full label, expected answer, acceptable path, rationale, private blueprint detail, disagreement detail, or adjudication detail into a Phoenix checkout, Git history, public review record, terminal transcript, or response.
- A private finding may name an opaque case ID and non-revealing defect category. It must not disclose the expected outcome or enough information to reconstruct it.
- Do not inspect, hash, mount, or open the retired protocol-v3 or rejected-v4 private archives. Their public material is a collision set; their private material remains outside scope.
- Do not run the experiment runner, a model, an agent, an arm, a schedule, a trial, prospective grading, validation, held-out evaluation, or outcome analysis. Do not generate a schedule or open Gate 1A.
- Outcome-free fixture materialization and deterministic checks of repository commands, builds, tests, status output, schemas, and grader semantics are allowed in temporary directories. They must not invoke an arm or create a trial or outcome record.
- Do not import the patch into the implementation workspace during review.
- Stop and report a custody breach if any prospective validation or held-out outcome is observed. A consumed tranche cannot be accepted by regenerating or relabeling it.

Fresh review checkout and LF materialization

The chair or review preparer must create a new detached worktree at the exact report-only commit. Do not reuse `D:\Work\personal\phoenix-evaluator-v4-review` or any authoring checkout. Before review, require:

    git -C D:\Work\personal\phoenix worktree add --detach D:\Work\personal\phoenix-evaluator-v4-revision-review 7a57479bc64cdb882bf17bfff7ed5d80c2f64760

In the new checkout, confirm `git status --short` is empty, then run exactly:

    git config --worktree core.autocrlf false
    git config --worktree core.eol lf
    git checkout-index --all --force
    git status --short

Run `checkout-index` only in the new detached review checkout. Stop if either status check is nonempty. This is the only permitted checkout materialization change. Never run `seal --write` in the review checkout.

Verify all sections below.

1. Commit identity, ancestry, patch equality, and two-commit boundary
   - Require detached HEAD at report-only commit `7a57479bc64cdb882bf17bfff7ed5d80c2f64760`.
   - Require its sole parent to be payload commit `e705560fef6012f84ec0199b87ce790a62c17d0e`, and require the payload's sole parent to be evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`.
   - Require the implementation freeze to be an ancestor of the evaluator base.
   - Require the report-only commit to change exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-revision-report.md` and no other path.
   - Independently capture the native stdout bytes of:

       git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 e705560fef6012f84ec0199b87ce790a62c17d0e

     Use a binary-safe process API or raw stream. Do not use a PowerShell pipeline, `>`, `Out-File`, `Set-Content`, or any text decode/re-encode step. Hash the raw bytes, require SHA-256 `96d25bd5b0340d71b71d46b8d030bd746323e29f224ffe6ba4c275f7765e4ded`, require length 1366164, and require byte-for-byte equality with the stored patch. Capture the command's raw stdout independently a second time and require equality with the first capture and stored file.
   - Require the base-to-payload diff to contain exactly 884 paths: 480 corpus paths, 400 fixture paths, and four registry/manifest paths, all under the approved public tranche locations.
   - Require the payload patch to exclude the public report. Require it to contain no full label, rationale, private note, protocol, world, runtime, prompt, arm, runner, grader, analysis, freeze, gate, schedule, trial, outcome, or implementation change.
   - Confirm `experiments/frontier-v1/pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

2. Pinned contract and reproducibility
   - Independently recompute the protocol's LF-normalized digest, both world identities, grader digest, raw world-build-manifest digest, and content-addressed `world_build_digest` using each field's specified algorithm. For `world_build_digest`, decode the manifest, serialize only its typed `build` record with Go `json.Marshal`, and SHA-256 those bytes. Do not copy the embedded digest or hash the entire manifest.
   - Stop before private-label review if any pinned identity differs.
   - From `experiments/frontier-v1/corpusctl`, run:

       go test -count=1 ./...
       go vet ./...
       go run ./cmd/corpusctl grader-digest --repo-root ../../..

     Require grader output `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`.
   - Run both seal commands without `--write` and capture stdout as raw bytes:

       go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json
       go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json

     Require byte equality with the committed manifests. The expected raw manifest digests are `sha256:4bab07c81226cd1373424dd24461f17588d0151b7dc38da6f16ce14d443d1304` for validation and `sha256:a93beedd292bb6128a6289cfde8e3b760b41360c8cef2faa70d0f8dc60cc5402` for held-out.
   - Keep the review checkout clean after every check.

3. Frozen allocation and identifiers
   - Derive the contract from `protocol.json`, schemas, and public registry classes, not the candidate report.
   - Require exactly 120 cases, 24 families, five cases per family, 120 fixtures, and 120 registry rows in each tranche.
   - Require eight direct-group families, each with `3 direct + 1 absence + 1 stale_frontier`.
   - Require eight recovery-group families, each with `3 recovery + 1 far_discovery + 1 temptation`.
   - Require eight mix families. Across them, require 12 `cascade`, 12 `adversarial_text`, and four each of `far_discovery`, `temptation`, `absence`, and `stale_frontier`; require at least one cascade and one adversarial-text case in every mix family.
   - Consequently require 24 direct, 24 recovery, and 12 of every other required class per tranche, with every distinct-family floor met.
   - Require family IDs `validation_family_101` through `validation_family_124` and `held_out_family_301` through `held_out_family_324`, all unique and absent from every collision set. Require 240 unique replacement case IDs, all absent from every collision set.

4. Five real variants and family independence
   - Inspect all 48 family mechanisms and all 240 public case/fixture pairs. Do not accept unique goal wording, tokens, filenames, nouns, service aliases, or fixture hashes alone as proof of five variants.
   - Require each family's five cases to differ materially in observable state, task type, intermediate evidence, success condition, recovery condition, or check evidence.
   - Independently inspect the public report's 48 family summaries and 240-row variant matrices. Verify every statement against the public case and fixture. Confirm the report does not reveal expected paths or evaluator-only answers.
   - Recompute pairwise structural screens across all 1,128 unordered family pairs and all 576 cross-tranche pairs. Investigate every high-similarity pair rather than accepting the candidate threshold or conclusion.
   - Reject parameterized global templates, noun swaps, mirrored validation/held-out families, repeated label skeletons, and graph diversity achieved with causally irrelevant padding.

5. Five collision sets and collision-inventory integrity
   - Compare the replacement independently with: current visible authoring; retired-v3 validation public material at the evaluator base; retired-v3 held-out public material at the evaluator base; rejected-v4 public material at commit `c193c786cc5a65ce6ae97efbd336b2a48f492898`; and the other replacement tranche.
   - Use read-only Git-object access for base and rejected-candidate material. Do not switch the review checkout.
   - Compare family IDs, case IDs, goals and grammar, fixture contents and templates, fixture/file-map digests, observable fixture shapes, service aliases, test catalogs, adversarial patterns, state-change patterns, action graphs, and public label-registry shape.
   - Reject renamed or lightly rewritten prior families. Live-path replacement under the approved tranche directories is expected and is not itself a collision.
   - Verify `notes/collision-inventory.json` at the private root by SHA-256 of exact raw bytes, with no normalization or JSON canonicalization. Require `sha256:af6c8d1a1630efb54c404aebf4410448ad5c324075b30413c6c12a0ce94e888f`.
   - Independently verify the inventory's completeness for all five records and every field it claims. Do not ask it to prove fields it does not contain; obtain additional public comparisons directly from Git trees.
   - Compare label-only structures between visible authoring labels and the new private labels, and between the two replacement tranches. Never use historical or rejected private labels.

6. Public opacity and prose
   - Inspect every public replacement case, fixture, registry, manifest, and every line of the public report.
   - Reject leakage of expected outcomes, acceptable paths, grading recipes, rationales, private causal details, evaluator-only class cues, or answer-bearing fixture templates.
   - Confirm public registries contain only `case_id`, `class`, `label_digest`, and `grading_script` per row.
   - Confirm adversarial fixture text is ordinary untrusted task data and does not identify its class or instruct the runner, evaluator, generator, labeler, or grader.
   - Audit all 24 cascade goals. Reject goals that enumerate, synonymize, or strongly cue the action recipe.
   - Proofread all 240 public goals and 240 fixture templates. Record grammar, capitalization, cloned wording, or contradictory-task defects with exact public paths.

7. Protocol-v4 executable semantics
   - Validate every public case and fixture against the current strict schemas and canonical world. Require every referenced capability, test, path, refusal, status, and output premise to agree with the frozen implementation.
   - For deterministic execution premises, materialize fixtures in temporary directories and exercise the relevant repository command, build, test, list, focus, status, or diff behavior without invoking an arm or creating a trial.
   - Require every `state_changes` replacement to be deterministic, use a safe relative path, and occur after a declared executable-act index shared by Arms A-E. Orientations never advance that index.
   - Require every stale-frontier case to create staleness through its declared shared event and to reacquire post-event evidence. Reject causal edges that rely on pre-event bytes after the event or claim evidence a tool cannot produce.
   - Require every absence case to request a genuinely unavailable capability and reward no unsupported action.
   - Require every cascade to use at least three executable verbs, with each next step justified by prior evidence. Reject optional or unrelated steps inserted only to make a graph unique.
   - If a cascade uses `tests.run` followed by `tests.focus`, require the intervening `tests.list` discovery step. Treat the authoring sequence as a collision set, not a template.
   - Verify final build/test/status/diff checks establish exactly what their tool output can prove.

8. Complete private-label and archive review
   - Verify the private blueprint, blueprint audit, label audit, and mechanical audit against their pinned raw SHA-256 identities before reading them. Confirm their declared roles are distinct from this reviewer.
   - Inspect all 240 full labels, not a sample. Require one schema-valid label per public case and no extra label.
   - Canonicalize every label with the candidate's `corpusctl` logic. Require every digest and class to match the public registry and every label to pin the current grader.
   - Independently reason from the frozen world, public case, fixture, and declared state events to verify each expected outcome. Prior audit approval is not proof.
   - For every listed acceptable path, apply `grade.go` path-matching and index-resolution rules. Require every path-bound `act_status.seq` and `act_output_matches.seq` position to exist on every matched path, including shorter paths. If no `act_sequence` check exists, treat `seq` as an absolute executable-act index.
   - Require checks to establish success and reject at least one plausible wrong path. Verify regex escaping, file paths, output fields, statuses, act counts, executable indexes, and state-event timing against actual deterministic tool behavior.
   - Require an empty acceptable path when no action is expected. Require rationales to explain the case-specific evidence and plausible wrong path without boilerplate or a false edit target.
   - Require deterministic grading where possible. If human judgment exists, require `arm_hidden: true`, a concrete outcome-independent rubric, and retained independent adjudication.
   - Independently recount check kinds. The candidate reports 216 `act_sequence`, 214 `act_output_matches`, 72 `act_path_absent`, 95 `file_matches`, 37 `act_status`, 24 `final_message_matches`, 24 `act_count`, and zero blinded-human checks. Reconcile every difference.
   - Recompute the private archive identity exactly:
     1. SHA-256 each label file's raw bytes. Do not normalize label contents.
     2. Form UTF-8 lines `<archive-relative-path>\tsha256:<lowercase-hex>`. Paths are relative to the private root, begin `labels/`, and use `/` separators.
     3. Sort all 240 lines by ordinal byte order of the archive-relative path.
     4. Join with LF and terminate the index with LF.
     5. SHA-256 the exact index bytes.
   - Require digest `sha256:c66d73e6d3db16b52feac05b6ceacc5ffb441b6c6ab4fec8705a87ea19d076bc`. Require the independently generated index to equal `notes/private-label-archive.index.txt` byte-for-byte.

9. Audit history, custody, and protected artifacts
   - Verify the blueprint auditor inspected all 48 families, 240 variants, 1,128 unordered pairs, 576 cross-tranche pairs, 24 cascades, 24 stale-frontier rows, and 288 collision records. Independently confirm its final APPROVE verdict is supported.
   - Verify the label auditor inspected all 240 exact label/case/fixture identities and deterministic premises. Independently confirm its final APPROVE verdict is supported.
   - Confirm all prior private disagreement categories are adjudicated and no finding remains open. Publish counts only, without private substance.
   - Confirm full labels and private records exist only under the new private custody root and not in any Phoenix checkout, Git history, public patch, or report.
   - Confirm no historical or rejected private archive was inspected or reused.
   - Confirm no model, arm, schedule, trial, prospective grader result, validation result, held-out result, or other outcome influenced design or review.
   - Confirm protocol, runtime and prompts, arm schemas, Arm B, worlds, runner, grader, analysis, accepted freeze entries, and gates are byte-identical to the evaluator base.
   - Keep both gates false and the implementation workspace unmodified.

Review record and verdict

Write the public review record to:

D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-revision-review.md

Return `ACCEPT`, `REVISE`, or `REJECT`, with P0-P3 findings and exact public paths and lines where applicable. For a private finding, use only an opaque case ID and non-revealing defect category.

The review record must include:

- reviewer role, independence declaration, date, reviewed commits, locations, and material seen;
- independently recomputed commit ancestry, patch bytes/digest/length/equality, report-only boundary, private-archive digest/index equality, collision-inventory digest, audit-record digests, grader, protocol, world, world-build, registry, and manifest identities;
- a requirement matrix for sections 1-9;
- independently computed allocation, family/class totals, five-variant evidence, pairwise/cross-tranche screens, cascade and stale-frontier counts, acceptable-path compatibility, and check-kind totals;
- test, vet, grader, dry-seal, fixture/tool premise, collision, opacity, and custody results with exit status or exact coverage counts;
- findings, accepted limitations, custody status, and confirmation that no outcome was opened;
- the smallest next artifact allowed by the verdict.

Verdict bar

- `ACCEPT` requires every identity, custody, boundary, schema, allocation, collision, opacity, executable-semantics, and complete-label premise to pass, with no unresolved P0 or P1. Every P2 must be fixed or explicitly justified as nonblocking.
- An identity-pin failure, patch mismatch, failed reviewer independence, private-label custody breach, observed prospective outcome, or consumed tranche is P0 and requires `REJECT`.
- Use `REVISE` for a correctable corpus, label, fixture, registry, manifest, report, audit-evidence, or review-record defect while identity and custody remain intact.
- Use `REJECT` when the candidate cannot be repaired without violating identity, validity, independence, or custody premises.

If `ACCEPT`, recommend only that the project chair consider importing the exact reviewed base-to-payload patch. The reviewer does not authorize import. The report-only commit is a separate review artifact and is not part of that patch. The chair must record any import decision separately and verify protected bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates afterward. Acceptance or import does not authorize schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.

If `REVISE`, keep the package unimported and name the smallest replacement artifact. Any private-label change requires new affected label digests and a new private-archive identity; determine all dependent registry, manifest, payload commit, patch, and report identities that must also change.

If `REJECT`, identify the failed independence, identity, validity, or custody premise and keep all gates closed.
```
