# Protocol-v4 sealed-corpus second revision independent review prompt

This assignment covers only the two-commit package produced from `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-handoff.md`. It does not apply to either rejected protocol-v4 package. The prior review prompts remain provenance records and must not be reused as the operative assignment.

Give the text below to a reviewer who is independent of both rejected packages and every second-revision authoring, labeling, blueprint-audit, label-audit, custody, packaging, handoff-preparation, and review-assignment-preparation role. The reviewer must not have implemented Phoenix. Use a fresh review session and a new detached review checkout.

Do not issue this assignment while it is untracked. Commit it in the implementation repository, then record that review-assignment commit in the message to the reviewer. The assignment commit is a process identity only; it does not change the candidate package identified below.

```text
Independently review the protocol-v4 second replacement `validation` and `held_out` corpus package for `frontier-v1`.

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, prospective grading, validation, held-out evaluation, or outcome analysis.

Candidate package

- Controlling-record commit: 0d99e1fe966513e2bb3a8e9ae330eaa75dfe88fc
- Fresh detached review checkout: D:\Work\personal\phoenix-evaluator-v4-revision-2-review
- Evaluator branch, identification only: codex/evaluator-v4-corpus-revision-2
- Evaluator base: c852101e8d7cb52e4569bf3866de54a0ce648b44
- Payload candidate commit: 273158f6a4189f67ef01bdb04992860b62262fef
- Report-only commit: 6356fc79e88c19aa197c8d7911116eb3e3ec20ba
- Required ancestry: report-only commit parent == payload candidate; payload candidate sole parent == evaluator base
- Implementation freeze: 4b51471200ed577db55fa38bed61c027c276d0f2
- First rejected v4 public candidate, collision use only: c193c786cc5a65ce6ae97efbd336b2a48f492898
- Rejected v4 revision payload, collision use only: e705560fef6012f84ec0199b87ce790a62c17d0e
- Rejected v4 revision report-only child, provenance only: 7a57479bc64cdb882bf17bfff7ed5d80c2f64760
- Public payload patch: D:\Work\personal\phoenix-evaluator-private-v4-revision-2\frontier-v1-v4-sealed-corpus-second-revision.patch
- Expected public-patch raw SHA-256: sha256:a7539840eb087d7c933b4b4eab486d885c48c15714831dd4fce4c3f6caf34d32
- Expected public-patch byte length: 1188302
- Public report at the report-only commit: docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md
- Expected public-report raw SHA-256: sha256:7dbed2a4ab32f492221ba6ecb5b8b06f637099ac05df3502e5e41e9c643f74f4
- Private custody root: D:\Work\personal\phoenix-evaluator-private-v4-revision-2
- Expected private-label custodian: second_revision.private_label_custodian
- Stored private-label archive index: D:\Work\personal\phoenix-evaluator-private-v4-revision-2\label-archive-index.tsv
- Expected private-label archive index: 240 lines, 24720 bytes, sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a
- Expected collision-inventory raw SHA-256: sha256:7ccb6a8f4f81a8b514e8c393109a604903d51849636b8c861d6f97e10ffbad0a
- Expected private blueprint raw SHA-256: sha256:0474c75fe905c22d002e09faa1425546156d693263bd9322f08debfbf3958332
- Expected final independent blueprint-audit raw SHA-256: sha256:11ac7114e070f31194c122416dc318615c1c51b18b1f9c0594be91ce4559ac16
- Expected materialization-and-causality-audit raw SHA-256: sha256:409d8c529b9cb239cec19d198a00d6bcc95433a00d0c67b3b952bb3acb7a63cb
- Expected first-pass-labeling-audit raw SHA-256: sha256:0b248f7759ddd67ea2cd4e83066dad3aaa375c46ee75cfb246557a03d757aee3
- Expected private mechanical-label-audit raw SHA-256: sha256:9d68d8a724ff171fadafb1d97d99639172c1b17e7cdb15daed67505758632d55
- Expected final independent label-audit raw SHA-256: sha256:af103b077920b651c08e3b713013e9df057088c1a9ad6ca10c8611581e6e51e0
- Expected exact rejected-byte-comparison record raw SHA-256: sha256:c8cc73b796c1a4df568edcf560ddfa5f873de007c618f231405f677f35e68618
- Expected private roles record raw SHA-256: sha256:6048f47fe58208561306d0aa671ba0c729bb79e622c9f10c65daca08f9cdf9b2
- Controlling second-revision handoff raw SHA-256 at controlling-record commit: sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc
- Original revision handoff raw SHA-256 at controlling-record commit: sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417
- Rejected-revision public review raw SHA-256 at controlling-record commit: sha256:6c9c622fb87547607303325e1f4ad903af9772b997108cf2d7615b525b5c03c0

Pinned freeze identities

- Protocol v4, LF-normalized UTF-8: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9
- Production/authoring world, canonical JSON: sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b
- Production/authoring world, raw bytes: sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143
- Current grader: sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc
- World-build manifest, raw bytes: sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f
- Content-addressed `world_build_digest`: sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d
- Validation registry, raw bytes: sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816
- Held-out registry, raw bytes: sha256:bba0717933ffa4d446f8ac3ce65e6ec44cf0bcc5a58600df3fcdfddda3c93206
- Validation sealed manifest: 89645 bytes, sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee
- Held-out sealed manifest: 88631 bytes, sha256:420f8fac9d063238737c87f637f54a319f505001c916468fa03bf7e22e4063aa

Required reading

Read the following completely before reviewing the candidate. The controlling process records are later than the evaluator base; read their exact blobs from the implementation repository at controlling-record commit `0d99e1fe...`, not from the detached candidate checkout.

- the controlling second-revision handoff at commit 0d99e1fe966513e2bb3a8e9ae330eaa75dfe88fc;
- docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md at the report-only commit;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-review.md, as the public defect record for rejected payload e705560f...;
- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md;
- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review.md, as the public defect record for rejected candidate c193c786...;
- docs/plans/2026-08-17-phoenix-world-plan.md, especially Task 0.4;
- experiments/frontier-v1/protocol.json, especially `tranche_design.family_allocation`;
- experiments/frontier-v1/schema/common.schema.json;
- experiments/frontier-v1/schema/case.schema.json;
- experiments/frontier-v1/schema/fixture.schema.json;
- experiments/frontier-v1/schema/label.schema.json;
- experiments/frontier-v1/corpusctl/internal/corpus/grade.go and the relevant `corpusctl` commands;
- spec/world.schema.json;
- experiments/frontier-v1/worlds/authoring.dev_repo.json;
- worlds/dev-repo/world.json;
- visible authoring cases, fixtures, and labels, for collision and opacity checks only;
- the complete evaluator-base-to-payload diff and the report-only commit diff;
- retired-v3 public validation and held-out material at the evaluator base, through read-only Git-object access;
- both rejected-v4 public payloads, through read-only Git-object access;
- every new public case, fixture, label-digest registry, and sealed manifest at the payload commit;
- through the private custodian, all 240 full labels, the stored archive index, family blueprints, collision inventory, materialization and causality evidence, audit history, exact rejected-byte record, and final private audit records.

Independence and custody

- Record the reviewer role, date, material seen, and all independence limitations. The reviewer must be distinct from every author, labeler, auditor, custodian, preparer, reviewer, and Phoenix implementer associated with either rejected package or the second revision. The reviewer must also be distinct from the author of this assignment.
- Treat every prior private `APPROVE` or `REVISE` record as evidence to verify, never as a substitute for this review. Independently inspect all 48 family designs and all 240 full labels.
- Obtain private artifacts only from `second_revision.private_label_custodian`. Verify each artifact's raw identity before relying on its content. Recompute and reconcile every embedded source identity; do not assume an audit's source pins identify the controlling tracked bytes.
- Do not copy a full label, expected answer, acceptable path, rationale, private blueprint detail, causal proof, disagreement, or adjudication into a Phoenix checkout, Git history, public review record, terminal transcript, or response.
- A private finding may name an opaque case ID and a non-revealing defect category. It must not disclose the expected outcome or enough information to reconstruct it.
- Do not inspect, enumerate, hash, mount, or open any prior private custody root. Prior public material is in scope only as collision evidence.
- Do not run the experiment runner, a model, an agent, an arm, a schedule, a trial, prospective grading, validation, held-out evaluation, or outcome analysis. Do not generate a schedule or open Gate 1A.
- Outcome-free fixture materialization and deterministic checks of repository commands, builds, tests, status output, schemas, and grader semantics are allowed in temporary directories. They must not invoke an arm or create a trial, grade, or outcome record.
- Do not import the patch into the implementation workspace. Stop and report a custody breach if any prospective validation or held-out outcome is observed.

Fresh review checkout and LF materialization

The chair or review preparer must create a new detached worktree at the exact report-only commit. Stop if the target path already exists; do not delete or repurpose it.

    git -C D:\Work\personal\phoenix worktree add --detach D:\Work\personal\phoenix-evaluator-v4-revision-2-review 6356fc79e88c19aa197c8d7911116eb3e3ec20ba

In the new checkout, require an empty `git status --short`, then run exactly:

    git config --worktree core.autocrlf false
    git config --worktree core.eol lf
    git checkout-index --all --force
    git status --short

Run `checkout-index` only in the new detached review checkout. Stop if either status check is nonempty. Never run `seal --write` in the review checkout. Do not amend, rebase, switch, commit, or modify the evaluator branch.

Verify all sections below.

1. Commit identity, ancestry, patch equality, and two-commit boundary

   - Require detached HEAD at report-only commit `6356fc79e88c19aa197c8d7911116eb3e3ec20ba`.
   - Require its sole parent to be payload commit `273158f6a4189f67ef01bdb04992860b62262fef`; require the payload's sole parent to be evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`; require the implementation freeze to be an ancestor of the base.
   - Require the report-only commit to add exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md` and no other path.
   - Independently capture the native stdout bytes of:

         git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 273158f6a4189f67ef01bdb04992860b62262fef

     Use a binary-safe process API or raw stream. Do not use a PowerShell pipeline, `>`, `Out-File`, `Set-Content`, or any text decode/re-encode step. Require 1188302 bytes, raw SHA-256 `a7539840eb087d7c933b4b4eab486d885c48c15714831dd4fce4c3f6caf34d32`, and byte equality with the stored patch. Capture the diff independently a second time and require equality among both captures and the stored file.
   - Require exactly 884 changed paths: 480 corpus paths, 400 fixture paths, and four registry/manifest paths, with no other path.
   - Require the payload to exclude the public report and every full label, rationale, blueprint, audit, collision inventory, private note, protocol, world, runtime, prompt, arm, runner, grader, analysis, freeze, gate, schedule, trial, outcome, or implementation change.
   - Confirm `pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

2. Pinned contract and reproducibility

   - Independently recompute the protocol's LF-normalized digest, both world identities, grader digest, raw world-build-manifest digest, and content-addressed `world_build_digest` using their specified algorithms. For `world_build_digest`, serialize only the typed `build` record with Go `json.Marshal` and hash those bytes.
   - Stop before private-label review if any pinned contract identity differs.
   - From `experiments/frontier-v1/corpusctl`, run:

         go test -count=1 ./...
         go vet ./...
         go run ./cmd/corpusctl grader-digest --repo-root ../../..

   - Run both seal commands without `--write`, capture raw stdout, and require byte equality with the committed manifests:

         go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json
         go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json

   - Recompute the two registry identities and the two manifest lengths and identities pinned above. Keep the review checkout clean after every check.

3. Frozen allocation and identifiers

   - Derive the allocation from `protocol.json`, schemas, and public registry classes, not from the candidate report.
   - Require exactly 120 cases, 24 families, five cases per family, 120 fixtures, and 120 registry rows in each tranche.
   - Require eight direct-group families, each with `3 direct + 1 absence + 1 stale_frontier`.
   - Require eight recovery-group families, each with `3 recovery + 1 far_discovery + 1 temptation`.
   - Require eight mix families per tranche. Across each mix group, require 12 `cascade`, 12 `adversarial_text`, and four each of `far_discovery`, `temptation`, `absence`, and `stale_frontier`; require at least one cascade and one adversarial-text case in every mix family.
   - Require 24 direct, 24 recovery, and 12 of every other required class per tranche, with every distinct-family floor met.
   - Recompute every family, fixture, and case identifier. Require uniqueness across the new tranches and absence from all five prior public collision sets.

4. Five real variants and family independence

   - Inspect all 48 mechanisms and all 240 public case/fixture pairs. Unique goals, tokens, filenames, nouns, service aliases, graphs, or hashes are not proof of real variants or independent families.
   - Require every family's five cases to differ materially in observable state, task type, intermediate evidence, success condition, recovery condition, or check evidence.
   - Verify every row of the public report's family/variant matrix against the public case and fixture. Confirm that the report accurately describes the tree without revealing expected paths or answers.
   - Independently screen and inspect all 1,128 unordered family pairs and all 576 cross-tranche pairs. Investigate every high-similarity pair, every shared topology or grammar, and the structural pairs forced by either rejected review. Do not adopt the candidate's threshold or verdict without independent judgment.
   - Reject a parameterized global generator, noun swaps, topology twins, numbered or slot-for-slot validation/held-out mirrors, shared goal or label skeletons, and graph diversity created by optional or causally irrelevant acts.
   - Inspect the private generation and audit scripts as provenance evidence for whether one shared implementation imposed family structure. Do not treat distinct serialized outputs as proof that the authoring mechanism was independent.

5. Collision sets and exact rejected-byte comparison

   - Independently compare the package with five prior public sets: visible authoring; retired-v3 validation at the evaluator base; retired-v3 held-out at the evaluator base; rejected candidate `c193c786...`; and rejected revision payload `e705560f...`. Also compare validation with held-out.
   - Use read-only Git-object access for historical and rejected public material. Do not switch the review checkout or execute code from either rejected payload.
   - Compare paths, identifiers, goals and goal grammar, raw fixture and file-map digests, fixture topology and observable shape, services and aliases, test catalogs and naming patterns, adversarial text, state changes, action graphs, and public registry structure.
   - Verify `notes/collision-inventory.json` by exact raw bytes and independently audit its coverage and correctness. Do not extend it with any prior private material.
   - Independently reproduce the second handoff's exact whole-file rejected-byte algorithm. Enumerate Git blobs recursively and binary-safely; subtract the evaluator-base blob set; hash unique rejected, new-public, and new-private whole files as `(byte_length, raw_sha256)` pairs; do not follow private symlink or reparse-point directories; do not scan substrings.
   - Reconcile the candidate's reported counts: 626 base blobs; rejected trees with 551 and 706 blobs; rejected unique sets with 329 and 484 blobs, 813 in union; 706 new-public tree blobs; 484 new-public unique blobs; 296 private regular files; and zero public, private, or combined intersections. Treat automated equality as screening only; it does not establish structural independence.
   - Compare label-only structures between visible authoring labels and new private labels, and between the two new tranches. Never inspect a prior private label archive.

6. Public opacity and prose

   - Inspect every public case, fixture, registry, manifest, and every line of the public report.
   - Reject leakage of expected outcomes, acceptable paths, grading recipes, rationales, private causal details, evaluator-only class cues, or answer-bearing fixture templates.
   - Confirm each public registry row contains only `case_id`, `class`, `label_digest`, and `grading_script`.
   - Confirm adversarial fixture text is ordinary untrusted task data and does not instruct the runner, evaluator, generator, labeler, or grader.
   - Audit all 24 cascade goals. Reject any goal that enumerates, synonymizes, or strongly cues the action recipe.
   - Proofread all 240 goals, all fixture templates, and the public report. Record public grammar, capitalization, cloned wording, contradictory-task, or inaccurate-report defects with exact paths and lines.

7. Protocol-v4 executable semantics and cascade causality

   - Validate every public case and fixture against the strict schemas and canonical world. Require every referenced capability, path, test, refusal, status, and output premise to agree with the frozen implementation.
   - Materialize deterministic fixtures in temporary directories and exercise the relevant repository commands, builds, tests, status, and diff behavior without invoking an arm or producing a grade or trial record.
   - Require each `state_changes` replacement to be deterministic, use a safe relative path, and occur after a declared executable-act index shared by Arms A-E. Orientations never advance that index.
   - Require every stale-frontier case to create staleness through its shared event and to reacquire post-event evidence. Reject labels that rely on pre-event bytes after the event.
   - Require every absence case to request a genuinely unavailable capability and reward no unsupported act.
   - Require every cascade to use at least three executable verbs whose order emerges from intermediate evidence. Reject optional or unrelated acts used as padding.
   - Independently materialize the pre-edit and post-edit premise for all 24 cascades. Require the terminal verification to fail to establish success before the earlier required change and to establish success after it. Verify that the terminal tool observes or consumes the edited artifact and that its output proves the label's claim.
   - If a cascade uses `tests.run` followed by `tests.focus`, require an intervening `tests.list` discovery act. Treat the authoring diagnostic sequence as a collision set, not a template.

8. Complete private-label and archive review

   - Verify all pinned private records before reading them and confirm their roles are distinct from this reviewer. Inspect the complete audit history, including retained `REVISE` records, and verify that every earlier defect was actually closed.
   - Inspect all 240 full labels. Require exactly one schema-valid label per public case, no extras, a current grader pin, and exact class and canonical digest agreement with the public registries.
   - Independently reason from the frozen world, case, fixture, blueprint, and state events to verify every expected outcome. Prior audit approval is not proof.
   - For every acceptable path, apply `grade.go` path matching and index resolution exactly. Check every path-bound `act_status.seq` and `act_output_matches.seq` against every matched path, including shorter paths. If no `act_sequence` check exists, treat `seq` as an absolute executable-act index.
   - Require checks to establish success and reject at least one plausible wrong path. Verify regex compilation and semantics, paths, statuses, output fields, act counts, executable indexes, forbidden acts, file states, and event timing against deterministic tool behavior.
   - Require an empty acceptable path exactly where no action is expected. Require rationales, expected results, acceptable paths, and checks to be case-specific rather than shared non-schema skeletons.
   - Independently inspect all repeated check-layout groups, especially the reported largest group of 71 and every group over eight rows. Accept repetition only when schema constraints force the layout or every target, premise, path, result, and rationale remains substantively case-specific.
   - Independently recount checks. Reconcile against the candidate's 1,140 total: 240 `act_sequence`, 264 `act_status`, 225 `act_output_matches`, 62 `act_path_absent`, 111 `file_matches`, 214 `final_message_matches`, and 24 `act_count`. The candidate reports 926 scripted annotations, 550 compiled Go-regexp patterns, and zero blinded-human checks; these are overlapping properties, not a partition.
   - Recompute the archive index exactly: hash each label's raw bytes; form UTF-8 `labels/<filename>\tsha256:<lowercase-hex>` lines; sort by archive-relative path in ordinal UTF-8 byte order; join with LF; terminate with LF. Require 240 lines, 24,720 bytes, digest `sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a`, and byte equality with `label-archive-index.tsv`.

9. Audit history, custody, report accuracy, and protected artifacts

   - Verify that the final blueprint audit independently covered 48 mechanisms, 240 variants, 1,128 unordered pairs, 576 cross-tranche pairs, every mix layout, 24 cascade graphs and causality records, 26 state-event records, and all five prior collision sets. Independently confirm whether its `APPROVE` verdict is supported.
   - Verify that the final label audit independently covered all 240 case/fixture/label/blueprint combinations, 1,140 checks, 550 regexes, 24 cascades, and 26 state events. Independently confirm whether its `APPROVE` verdict is supported.
   - Confirm the reported 13 initial label disagreements were adjudicated, the two retained `REVISE` audits are immutable, and no disagreement or finding remains open. Publish counts only.
   - Recompute every private audit's embedded source hash from the authoritative artifact it claims to identify. Record and classify every mismatch; a top-level audit-file hash match does not cure an incorrect source pin.
   - Confirm full labels and private records exist only under the new custody root and not in a Phoenix checkout, Git history, public patch, or report.
   - Confirm neither prior private root was opened or reused and no prior script, fixture, blueprint, label, rationale, check layout, or label skeleton was used as scaffolding.
   - Confirm the public report's identities, counts, audit claims, family descriptions, cascade claims, and collision claims match independently recomputed evidence.
   - Confirm no model, arm, schedule, trial, prospective grade, validation result, held-out result, or outcome influenced design, packaging, or review.
   - Confirm protocol, runtime and prompts, arm schemas, Arm B, worlds, runner, grader, analysis, accepted freeze entries, `pre-validation-artifacts.json`, and both gates are byte-identical to the evaluator base. Keep both gates false and the implementation workspace otherwise unmodified.

Review record and verdict

Write the public review record to:

D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-second-revision-review.md

Return `ACCEPT`, `REVISE`, or `REJECT`, with P0-P3 findings and exact public paths and lines where applicable. For a private finding, use only an opaque case ID and a non-revealing defect category.

The review record must include:

- reviewer role, independence declaration, date, reviewed commits and paths, and material seen;
- independently recomputed ancestry, patch identity and equality, report-only boundary, report identity, archive-index identity and equality, collision-inventory identity, private audit identities and source-pin reconciliation, grader, protocol, world, world-build, registries, and manifests;
- a requirement matrix for sections 1-9;
- independently computed allocation, family and class totals, five-variant evidence, pairwise and cross-tranche screens, mix-layout results, cascade and event counts, acceptable-path compatibility, repeated-label-layout review, and check totals;
- test, vet, grader, dry-seal, fixture materialization, all-cascade pre/post causality, collision, exact-byte, opacity, and custody results with exit status or exact coverage counts;
- findings, accepted limitations, custody status, and explicit confirmation that no outcome was opened;
- the smallest next artifact allowed by the verdict.

Verdict bar

- `ACCEPT` requires every identity, custody, boundary, schema, allocation, collision, opacity, independence, executable-semantics, cascade-causality, and complete-label premise to pass, with no unresolved P0 or P1. Every P2 must be fixed or specifically justified as nonblocking.
- A top-level package identity failure involving commit ancestry, patch bytes, archive-index bytes, failed reviewer independence, private-label custody, an observed prospective outcome, or a consumed tranche is P0 and requires `REJECT`.
- Classify an embedded audit source-pin mismatch by consequence. Use `REVISE` when the candidate package and custody remain intact but audit support must be corrected or repeated; use `REJECT` only when the mismatch establishes a top-level identity, independence, validity, or custody failure.
- Use `REVISE` for a correctable corpus, fixture, label, registry, manifest, report, private audit-evidence, or review-record defect while package identity and custody remain intact.
- Use `REJECT` when the package cannot be repaired without violating identity, validity, independence, or custody premises.

If `ACCEPT`, recommend only that the project chair consider importing the exact reviewed base-to-payload patch. The reviewer does not authorize import. The report-only commit is not part of that patch. The chair must record any import decision separately and verify protected bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates afterward. Acceptance or import does not authorize schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.

If `REVISE`, keep the package unimported and name the smallest replacement artifact. Any private-label change requires new affected label digests and a new archive identity; identify every dependent registry, manifest, payload commit, patch, and report identity that must also change.

If `REJECT`, identify the failed independence, identity, validity, or custody premise and keep both gates closed.
```
