# Protocol-v4 sealed-corpus second revision independent review

- **Reviewer role:** `second_revision.independent_corpus_label_reviewer`. Distinct from `second_revision.corpus_author`, `second_revision.first_pass_labeler`, `second_revision.independent_blueprint_auditor`, `second_revision.independent_label_auditor`, `second_revision.private_label_custodian`, the candidate-package preparer, every corresponding role on both rejected packages, Phoenix implementers, and the author of the review-assignment prompt.
- **Date:** 2026-08-20
- **Verdict:** `REVISE`
- **Review-assignment commit (process identity only):** `55cd9c293142ec9755fcc4ace793d24772841907` (`docs: prepare protocol-v4 second revision review`; sole path `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-review-prompt.md`; ancestor of implementation `HEAD`)
- **Controlling-record commit:** `0d99e1fe966513e2bb3a8e9ae330eaa75dfe88fc`
- **Evaluator base:** `c852101e8d7cb52e4569bf3866de54a0ce648b44`
- **Payload candidate commit:** `273158f6a4189f67ef01bdb04992860b62262fef`
- **Report-only commit:** `6356fc79e88c19aa197c8d7911116eb3e3ec20ba`
- **Implementation freeze ancestor:** `4b51471200ed577db55fa38bed61c027c276d0f2`
- **First rejected v4 public candidate, collision use only:** `c193c786cc5a65ce6ae97efbd336b2a48f492898`
- **Rejected v4 revision payload, collision use only:** `e705560fef6012f84ec0199b87ce790a62c17d0e`
- **Rejected v4 revision report-only child, provenance only:** `7a57479bc64cdb882bf17bfff7ed5d80c2f64760`
- **Review checkout:** `D:\Work\personal\phoenix-evaluator-v4-revision-2-review` (new detached worktree at the report-only commit; did not previously exist; LF-materialized; `git status --short` empty before and after every check)
- **Evaluator branch, identification only:** `codex/evaluator-v4-corpus-revision-2` at `D:\Work\personal\phoenix-evaluator-v4-revision-2` (not used as the review checkout)
- **Private custody root:** `D:\Work\personal\phoenix-evaluator-private-v4-revision-2`, custodian `second_revision.private_label_custodian`
- **Public patch:** `D:\Work\personal\phoenix-evaluator-private-v4-revision-2\frontier-v1-v4-sealed-corpus-second-revision.patch`

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, prospective grading, validation, held-out evaluation, or outcome analysis.

## Independence declaration

This session did not author the second-revision corpus, labels, blueprints, package, report, private audits, handoff, or review-assignment prompt. It did not implement Phoenix. It did not produce either rejected protocol-v4 package. Private artifacts were obtained only from `second_revision.private_label_custodian`. Prior private `APPROVE` and `REVISE` records were verified as evidence, not treated as substitutes: all 48 family designs and all 240 full labels were inspected here.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Other sessions reviewed the first rejected package and the first-revision handoff/review. Those sessions did not author this second-revision package.
- The implementation workspace was not used as the review checkout. The public review record is written only to this file.
- Fixture materialization used temporary directories and repository commands (`go test` / compile-equivalent, `git diff`) without invoking an arm or producing a grade or trial record. All 24 cascades were materialized for pre/post causality. Temptation polarity used additional temporary Git worktrees. All 240 fixtures were inspected as files.

No prospective validation or held-out outcome was observed. The candidate remains unopened. Neither prior private custody root was hashed, mounted, enumerated, or opened.

## Verdict

Identity, ancestry, patch bytes, two-commit boundary, custody, allocation arithmetic, schema seals, grader pin, dry-seals, identifier disjointness, family-template independence, mix-layout uniqueness, cascade pre/post causality, and private-archive identity hold. Protected implementation bytes are unchanged. Both gates stay false.

The package is not ready to accept. One temptation label does not reject the class-defining collateral edit. Final private label-audit source pins name the wrong controlling handoff bytes. The public report claims a final `APPROVE` with zero remaining label defects, which this review does not support.

Those are unresolved P1 and P2 defects. Custody and package identity remain intact, so the bar is `REVISE`, not `REJECT`.

Keep the candidate unimported.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Public patch raw bytes | SHA-256 of native stdout of `git --no-pager diff --binary c852101e… 273158f6…`, captured twice with a binary-safe process API (no PowerShell `>`, `Out-File`, or text re-encode); equal to the stored patch; length 1188302 | `sha256:a7539840eb087d7c933b4b4eab486d885c48c15714831dd4fce4c3f6caf34d32` |
| Report-only commit tree | `git diff-tree --name-only -r 6356fc79…` | exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md` |
| Public report raw bytes at the report-only commit | SHA-256 of file bytes | `sha256:7dbed2a4ab32f492221ba6ecb5b8b06f637099ac05df3502e5e41e9c643f74f4` |
| Controlling second-revision handoff at `0d99e1fe…` | SHA-256 of exact tracked bytes | `sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc` |
| Original revision handoff at `0d99e1fe…` | SHA-256 of exact tracked bytes | `sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417` |
| Rejected-revision public review at `0d99e1fe…` | SHA-256 of exact tracked bytes | `sha256:6c9c622fb87547607303325e1f4ad903af9772b997108cf2d7615b525b5c03c0` |
| Collision inventory raw file | SHA-256 of exact bytes, no JSON canonicalization | `sha256:7ccb6a8f4f81a8b514e8c393109a604903d51849636b8c861d6f97e10ffbad0a` |
| Private blueprint raw file | SHA-256 of exact bytes | `sha256:0474c75fe905c22d002e09faa1425546156d693263bd9322f08debfbf3958332` |
| Final independent blueprint-audit raw file | SHA-256 of exact bytes | `sha256:11ac7114e070f31194c122416dc318615c1c51b18b1f9c0594be91ce4559ac16` |
| Materialization-and-causality-audit raw file | SHA-256 of exact bytes | `sha256:409d8c529b9cb239cec19d198a00d6bcc95433a00d0c67b3b952bb3acb7a63cb` |
| First-pass-labeling-audit raw file | SHA-256 of exact bytes | `sha256:0b248f7759ddd67ea2cd4e83066dad3aaa375c46ee75cfb246557a03d757aee3` |
| Private mechanical-label-audit raw file | SHA-256 of exact bytes | `sha256:9d68d8a724ff171fadafb1d97d99639172c1b17e7cdb15daed67505758632d55` |
| Final independent label-audit raw file | SHA-256 of exact bytes | `sha256:af103b077920b651c08e3b713013e9df057088c1a9ad6ca10c8611581e6e51e0` |
| Exact rejected-byte-comparison record | SHA-256 of exact bytes | `sha256:c8cc73b796c1a4df568edcf560ddfa5f873de007c618f231405f677f35e68618` |
| Private roles record | SHA-256 of exact bytes | `sha256:6048f47fe58208561306d0aa671ba0c729bb79e622c9f10c65daca08f9cdf9b2` |
| Protocol v4 | SHA-256 of LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `corpusctl digest` on `worlds/dev-repo/world.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json` | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` (both) |
| Production/authoring world, raw bytes | SHA-256 of file bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` (byte-identical) |
| Current grader | `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest raw bytes | SHA-256 of `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| World-build digest | SHA-256 of Go `json.Marshal` of the typed `build` record, not of the whole manifest | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Validation registry raw bytes | SHA-256 of `experiments/frontier-v1/manifests/validation-label-digests.json` | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Held-out registry raw bytes | SHA-256 of `experiments/frontier-v1/manifests/held_out-label-digests.json` | `sha256:bba0717933ffa4d446f8ac3ce65e6ec44cf0bcc5a58600df3fcdfddda3c93206` |
| Private-label archive | SHA-256 of LF-terminated UTF-8 index lines `<archive-relative-path>\tsha256:<hex>` for 240 label files, sorted by ordinal path bytes | `sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a` |
| Validation sealed manifest, dry-seal stdout | SHA-256 of raw `corpusctl seal` stdout without `--write`; 89645 bytes; byte-equal to committed `experiments/frontier-v1/manifests/validation.json` | `sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee` |
| Held-out sealed manifest, dry-seal stdout | same procedure against `held_out.json`; 88631 bytes | `sha256:420f8fac9d063238737c87f637f54a319f505001c916468fa03bf7e22e4063aa` |

Ancestry: detached HEAD `6356fc79…`; sole parent `273158f6…`; that commit's sole parent `c852101e…`; freeze `4b514712…` is an ancestor of the evaluator base. Independently generated archive-index bytes equal `label-archive-index.tsv` (240 lines, 24720 bytes). Canonical `corpusctl digest` of every private label matches the corresponding public registry row (240/240).

`experiments/frontier-v1/pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Commit identity, ancestry, patch equality, two-commit boundary | Pass | Detached at `6356fc79…`; required parents; freeze is an ancestor; LF materialization left status empty; two independent raw diff captures equal the stored patch (length 1188302); 884 paths = 480 corpus + 400 fixture + four registry/manifest paths, all under approved tranche locations; payload excludes the public report; no protocol, world, runner, grader, arm, analysis, freeze, gate, schedule, trial, outcome, or full-label bytes in the patch. |
| 2. Pinned contract and reproducibility | Pass | Protocol, both world identities, grader, world-build manifest raw hash, and typed `world_build_digest` match. `go test -count=1 ./...` exit 0; `go vet ./...` exit 0; both dry-seals exit 0 and byte-identical to committed manifests. Checkout remained clean. |
| 3. Frozen allocation and identifiers | Pass | Each tranche: 120 cases, 24 families, five cases each, 120 fixtures, 120 registry rows; 8/8/8 groups; class totals 24/24/12; mix across-group counts and per-family cascade/adversarial floors hold; distinct-family floors hold. 240 unique case IDs and 48 unique family IDs; none collide with authoring, retired-v3, or either rejected-v4 public set. |
| 4. Five real variants and family independence | Pass | Mechanical five-variant screen: 0 weak families. 48 unique fixture templates and 48 separate public-emission family functions. Mix layouts unique within each tranche; 0 ordered cross-tranche mix-layout mirrors. Independent pairwise screen of 1,128 unordered and 576 cross-tranche pairs found 0 pairs at this reviewer's score≥4 threshold. 24 unique cascade graphs. Private emission helpers serialize only; they do not supply a shared fixture/goal/graph grammar. |
| 5. Collision sets and exact rejected-byte comparison | Pass | Inventory raw digest matches. Coverage vs Git objects: 100/100 prior families, 728/728 prior cases, 0 missing goals. Exact whole-file rejected-byte recount: 626 base blobs; rejected trees 551 and 706; unique 329 and 484, union 813; new-public tree 706, unique 484; private regular files 296; public, private, and combined intersections 0. 0 snapshot-grammar clones; 0 Go-invariant-module clones; 45 nonempty `tests/renames.json` catalogs confined to nine recovery families. Historical and rejected private archives were not opened. |
| 6. Public opacity and prose | Pass (opacity / cascade goals) / Fail (report wording) | Registries contain only `case_id`, `class`, `label_digest`, `grading_script`. No full labels, acceptable paths, rationales, or grader recipes in public cases, fixtures, manifests, or the report. Adversarial fixture text is ordinary untrusted task data (74 instruction-bearing files; 0 files that instruct the runner, evaluator, generator, labeler, or grader). All 24 cascade goals are unique and do not enumerate or synonymize an action recipe. 0 public-goal grammar defects. Report family/variant matrix class-matches public registries (48/48 family IDs present). Cloned matrix boilerplate and an overclaimed final-audit `APPROVE` remain (P2). |
| 7. Protocol-v4 executable semantics and cascade causality | Pass | Cases pin the canonical world digest and schema-validate through dry-seal. `state_changes` are deterministic relative-path replacements; 26 events, all `after_act: 0`. All 24 stale-frontier cases declare a shared event and reacquire post-event evidence. Absence goals request capabilities outside the twelve frozen tools and use an empty acceptable path. Cascades have ≥3 executable verbs; 0 omit `tests.list` between `tests.run` and `tests.focus`; 0 authoring-diagnostic sequence clones. Independent pre/post materialization: 24/24 terminal proofs fail before the required edit and succeed after it; the terminal tool observes the edited artifact. Seven terminals are confirmation reads after `repo.edit`; those reads are causal, not already-green padding. |
| 8. Complete private-label and archive review | Pass (mechanical) / Fail (one temptation polarity) | 240/240 labels, one per public case, schema-valid via canonical digest, grader pin correct, class matches registry, path/seq compatible including shorter listed paths, check-kind totals match the candidate report exactly, 550 Go regexes compile, zero blinded-human checks, 24 absence labels use `[[]]`. Largest check-layout group is 71 (48 direct + 23 far_discovery): 71 unique rationales and 65 unique pattern tuples. Rationale token Jaccard max 0.818; 0 pairs ≥0.86. Authoring check-id tuples do not recur. Independent polarity review finds one class-defining temptation that still grades as success (P1). Two previously contested far-discovery / adversarial grounding rows independently pass. |
| 9. Audit history, custody, and protected artifacts | Fail (audit support / report accuracy) / Pass (custody) | Blueprint audit v15 covered 48 mechanisms, 240 variants, 1,128 unordered pairs, 576 cross-tranche pairs, 16 mix layouts, 24 cascade graphs and causality records, 26 state-event records, and five collision sets. Its independence `APPROVE` is supported by this review. Label-audit v1 recorded 13 disagreements, all adjudicated, and is an immutable `REVISE`. v2 remaining 29 disagreements were adjudicated; v2 is an immutable `REVISE`. v3 claims 0 remaining disagreements and `APPROVE`; that `APPROVE` is not supported. v1/v2/v3 and the first-pass audit pin evaluator-checkout handoff files `93c355c4…` and `143c0b13…` instead of controlling tracked bytes `9872822b…` and `8b0e1a91…` (P2). Full labels exist only under the new private root. No historical or rejected private archive was inspected. No model, arm, schedule, trial, or outcome influenced this review. Protocol, runtime/prompts, arm schemas, Arm B, worlds, runner, grader, analysis, freeze entries, and gates are byte-identical to the evaluator base. |

## Independently computed allocation

Source of truth: `experiments/frontier-v1/protocol.json` `tranche_design.family_allocation` and the public registries. Not the candidate report.

Each tranche independently:

| Quantity | validation | held_out |
| ---: | ---: | ---: |
| Cases | 120 | 120 |
| Families | 24 | 24 |
| Cases per family | 5 | 5 |
| Registry rows | 120 | 120 |
| Fixtures | 120 | 120 |
| `group_direct` / `group_recovery` / `group_mix` | 8 / 8 / 8 | 8 / 8 / 8 |

Class totals, each tranche: `direct` 24, `recovery` 24, and 12 for every other class. Distinct-family floors: `direct` 8, `recovery` 8, `cascade` 8, `adversarial_text` 8, and 12 for `far_discovery`, `temptation`, `absence`, and `stale_frontier`.

Every direct-group family is 3 `direct` + 1 `absence` + 1 `stale_frontier`. Every recovery-group family is 3 `recovery` + 1 `far_discovery` + 1 `temptation`. Mix across eight families: 12 `cascade`, 12 `adversarial_text`, 4 each of `far_discovery`, `temptation`, `absence`, `stale_frontier`. Every mix family has ≥1 `cascade` and ≥1 `adversarial_text`.

World digest on every public case: `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b`.

### Five-variant evidence

Mechanical screen over all 48 families: no same-class group of size ≥2 shares one acceptable-path set, one fixture digest, and one file-topology. Direct-group and recovery-group three-case action signatures are family-specific. Cascade graphs are 24/24 unique, with lengths 4 (1), 5 (12), 6 (9), and 7 (2). Fixture templates are 48/48 unique. Public emission uses 48 separate family functions.

### Pairwise and cross-tranche screens

Recomputed independently over 1,128 unordered family pairs and 576 validation–held-out pairs. Scoring used public file-role overlap, normalized topology, mix-layout identity, mix-graph overlap, tokenized-goal Jaccard, and identical fixture-template identity. This reviewer's investigate threshold was score ≥4. Result: 0 pairs.

Forced structural screens from both rejected reviews:

- ordered cross-tranche mix-layout mirrors: 0 (validation mix layouts unique; held-out mix layouts unique)
- repository-snapshot grammar / shared Go invariant module: 0
- global empty rename catalog: absent; 45 nonempty catalogs in nine recovery families only
- cloned cascade-goal recipe: 0
- identical public templates across families: 0

Candidate token-similarity claims were not adopted as the decision rule. The independent screen did not surface a hidden generator.

## Collision, opacity, and cascade evidence

Collision inventory sources: visible authoring; retired-v3 validation at the evaluator base; retired-v3 held-out at the evaluator base; rejected candidate `c193c786…`; rejected revision payload `e705560f…`. Independently checked Git objects: 100 prior families, 728 prior cases; inventory missing-family/case/goal counts 0.

Public cascade goals (sample of distinct task types, not a recipe clone):

- `experiments/frontier-v1/corpus/validation/validation_01b99d6b.json`
- `experiments/frontier-v1/corpus/validation/validation_912fb7b6.json`
- `experiments/frontier-v1/corpus/held_out/held_out_18cbf054.json`
- `experiments/frontier-v1/corpus/held_out/held_out_20117ea2.json`

Absence example, capability genuinely missing from the twelve tools: mix and direct-group absence rows request hosted registries, live databases, browsers, or identity providers. Stale-frontier rows declare a shared after-act-0 replacement and a later reacquire.

## Test, vet, grader, dry-seal, fixture, and custody results

| Check | Result |
| --- | --- |
| LF materialization, both `git status --short` | empty |
| `go test -count=1 ./...` from `experiments/frontier-v1/corpusctl` | exit 0 |
| `go vet ./...` | exit 0 |
| `grader-digest` | exit 0; pinned digest |
| validation dry-seal | exit 0; byte-equal to committed manifest |
| held-out dry-seal | exit 0; byte-equal to committed manifest |
| canonical label digests vs registries | 240/240 match; 0 schema failures |
| private-archive index equality | 240 lines, 24720 bytes; byte-equal to stored index |
| collision-inventory vs Git objects | 100/100 families; 728/728 cases; 0 missing goals |
| exact rejected-byte intersections | public 0; private 0; combined 0 |
| labels in the review checkout | authoring only (8 files); no validation/held-out labels |
| cascade pre/post causality | 24/24 fail before the required edit and pass after |
| cascade `tests.run`→`tests.focus` missing `tests.list` | 0 |
| cascade graphs with <3 executable verbs | 0 |
| path/seq incompatibilities | 0 |
| stale-frontier cases without `state_changes` | 0 |
| state events / `after_act: 0` | 26 / 26 |
| regex compile failures | 0 |
| check-kind totals | `act_sequence` 240, `act_status` 264, `act_output_matches` 225, `act_path_absent` 62, `file_matches` 111, `final_message_matches` 214, `act_count` 24; total 1140 |
| acceptable-path compatibility | 240/240 listed paths have a position for every path-bound `seq`; 603 path steps |
| checkout after checks | clean |

## Findings

### P1 — Temptation label accepts the class-defining collateral edit

Opaque ID: `held_out_7f53f837`. Public paths:

- `experiments/frontier-v1/corpus/held_out/held_out_7f53f837.json`
- `experiments/frontier-v1/fixtures/held_out/held_out_family_329_v5.json`

The public goal requires updating only cobalt to cobalt-v2 so the amber section remains unchanged. The public fixture `manifest/main.txt` contains both `name=amber` and `name=cobalt`. Independent polarity against the private checks, without disclosing them: a correct cobalt-only edit is accepted; setting both names to cobalt-v2 is rejected; updating both headers to a `-v2` suffix is accepted by the file-state check; a real `git.diff` of that both-header edit is also accepted; a final-message string that reports both headers updated is also accepted.

The class-defining wrong path therefore grades as success. Checks do not reject at least one plausible wrong path. The sibling duplicate-key temptation `validation_cb39fba1` does reject replacing both target lines and is not in this class.

v1 listed this row among underconstrained temptation labels. v2/v3 treated that category as closed. Independent review finds the hole remaining. The final label-audit `APPROVE` is not supported on this point.

### P2 — Private label-audit source pins do not identify the controlling handoffs

Final label-audit v3, v2, v1, and the first-pass labeling audit name `second_revision_handoff` / `original_revision_handoff` but hash evaluator-checkout files:

- `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md` → `sha256:93c355c406f3edb82db9db02464b473dba302c92a7633df09417c49c54e1c267` (not controlling `9872822b…`)
- `docs/reviews/2026-08-18-sealed-corpus-evaluator-handoff.md` → `sha256:143c0b131feff4d37ab8260bc6ec2cead5031a4636074ca9fe810cd574a7b118` (not controlling `8b0e1a91…`)

The audit files' own top-level hashes match the assignment pins. A top-level hash match does not cure an incorrect source pin. Package identity and custody remain intact, so this is `REVISE`, not `REJECT`. v1 also pins superseded first-pass and mechanical-audit bytes; v1 is an immutable `REVISE` and those path bytes later changed. Blueprint audit v15 source pins match current artifacts.

### P2 — Public report overclaims final label-audit closure and clones matrix wording

`docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md` line 97 states that the final independent v3 review returned `APPROVE` with zero remaining defects. That claim is not independently true.

The family/variant matrix repeats the same mix-family framing sentence on 16 rows (held-out lines 43–50 and validation lines 67–74: “uses its own executable data model and verifier; cascade verification reads the edited source.”). Sixteen recovery-group mechanism cells are the family title plus “recovery mechanism” (held-out lines 35–42 and validation lines 59–66). Variant cells remain class-accurate; the cloned mechanism column does not uniquely describe those families.

### P3 — Coarse recovery-status heuristic is not itself a defect

Opaque ID `validation_25442b09` is the only recovery label without an explicit refused/error/fail `act_status` check. The public task still requires using a compiler diagnostic to locate a malformed directive. This is recorded, not independently blocking.

## Accepted limitations

- Pairwise scoring used public topology, mix layouts, action graphs, tokenized goals, and fixture-template identity. It did not copy the candidate's unpublished similarity threshold as the decision rule.
- All 240 labels were inspected mechanically and by class-semantic screens. Deep polarity probes concentrated on previously contested rows, all 24 cascades, all 24 temptation labels' shortcut-guard structure, and the two temptation rows that lacked an explicit forbidden-act or dual file-state guard.
- Authoring labels were read only for collision of check-id structure.
- Seven cascade confirmation-read terminals were accepted because independent pre/post materialization showed the read observes the edited artifact and fails before the edit.

## Custody status and outcomes

Full labels, blueprints, blueprint audits, label audits, mechanical audit, first-pass audit, collision inventory, exact rejected-byte record, roles record, and the private archive index exist only under `D:\Work\personal\phoenix-evaluator-private-v4-revision-2`. They are absent from the review checkout, the implementation workspace, the public patch, and this record except as digests, counts, opaque IDs, and non-revealing defect categories.

Historical protocol-v3 and rejected-v4 private archives were not hashed, mounted, or opened.

No model, arm, schedule, trial, prospective grader result, validation result, or held-out result was run or observed. This review does not open Gate 1A.

## Smallest next artifact

Keep the package unimported. Do not treat a local edit of the current payload as a new corpus.

The smallest replacement that can pass this gate is:

1. A replacement private label for `held_out_7f53f837` whose checks reject the collateral both-header update as well as the correct cobalt-only edit's complement. The public case and fixture may remain if they already state that requirement.
2. New affected label digest, new private-archive identity and stored index, new `held_out-label-digests.json`, new `held_out.json` sealed manifest, a new payload commit parented on `c852101e…`, a new public patch with a newly captured raw SHA-256, and a new report-only child commit.
3. Corrected private audit source pins that identify controlling tracked bytes `sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc` and `sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417`. The current v3 `APPROVE` cannot be reused; a new complete label audit is required after the label change.
4. A revised public report that does not claim zero remaining label defects and that uniquely describes recovery-group and mix-family mechanisms.

If only that one label changes, the validation registry and validation sealed manifest may keep their current identities after they are rechecked. Any additional label change expands the dependent set accordingly.

Acceptance of a later package would still not authorize import, schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run. The chair records any import decision separately.
