# Protocol-v4 sealed-corpus revision independent review

- **Reviewer role:** Independent corpus-and-label gate reviewer. Distinct from `revision.corpus_author`, `revision.first_pass_labeler`, `revision.independent_blueprint_auditor`, `revision.independent_label_auditor`, `revision.private_label_custodian`, the candidate-package preparer, Phoenix implementers, and the author of the review prompt.
- **Date:** 2026-08-20
- **Verdict:** `REVISE`
- **Evaluator base:** `c852101e8d7cb52e4569bf3866de54a0ce648b44`
- **Payload candidate commit:** `e705560fef6012f84ec0199b87ce790a62c17d0e`
- **Report-only commit:** `7a57479bc64cdb882bf17bfff7ed5d80c2f64760`
- **Implementation freeze ancestor:** `4b51471200ed577db55fa38bed61c027c276d0f2`
- **Rejected v4 public candidate, collision use only:** `c193c786cc5a65ce6ae97efbd336b2a48f492898`
- **Review checkout:** `D:\Work\personal\phoenix-evaluator-v4-revision-review` (new detached worktree at the report-only commit; LF-materialized; `git status --short` empty before and after every check)
- **Private custody root:** `D:\Work\personal\phoenix-evaluator-private-v4-revision`, custodian `revision.private_label_custodian`
- **Public patch:** `D:\Work\personal\phoenix-evaluator-private-v4-revision\frontier-v1-v4-sealed-corpus-revision.patch`

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, validation, held-out evaluation, or outcome analysis.

## Independence declaration

This session did not author the replacement corpus, labels, blueprints, package, report, private audits, or the review prompt. It did not implement Phoenix. Private artifacts were obtained only from `revision.private_label_custodian`. Prior private audits were verified, not treated as substitutes: all 48 family designs and all 240 full labels were inspected here.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Earlier sessions reviewed the rejected one-commit candidate and the revision *handoff* document. Those sessions did not produce this replacement package.
- The implementation workspace was not used as the review checkout. The public review record is written only to this file.
- Fixture materialization used temporary directories and repository `go test` / compile-equivalent commands. Six fixtures were executed that way; all 240 fixtures were inspected as files. No arm, runner, schedule, trial, or prospective grade was invoked.

No prospective validation or held-out outcome was observed. The candidate remains unopened.

## Verdict

Identity, ancestry, patch bytes, two-commit boundary, custody, allocation arithmetic, schema seals, grader pin, dry-seals, identifier disjointness, and private-archive identity hold. Protected implementation bytes are unchanged. Both gates stay false.

The replacement is not 48 independent generating families. Every one of the 240 fixtures is an instance of one repository-snapshot generator (same template grammar, same Go invariant module, same empty rename catalog, same adversarial request paragraph, domain noun and family number swapped). All eight mix families are slot-for-slot class-layout mirrors of `validation_family_1xx` with `held_out_family_(1xx+200)`. All 24 cascade goals are the same reconcile-from-authority recipe with a trailing compile/test/confirm clause. Eighteen of those cascades then require `repo.build` or `tests.*` after editing a record file that the Go invariant does not read; the materialized suite is already green before the reconcile.

Those are unresolved P1 corpus defects. Custody and identity remain intact, so the bar is `REVISE`, not `REJECT`.

Keep the candidate unimported.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Public patch raw bytes | SHA-256 of native stdout of `git --no-pager diff --binary c852101e… e705560f…`, captured twice with a binary-safe process API; equal to the stored patch; length 1366164 | `sha256:96d25bd5b0340d71b71d46b8d030bd746323e29f224ffe6ba4c275f7765e4ded` |
| Report-only commit tree | `git diff-tree --name-only -r 7a57479…` | exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-revision-report.md` |
| Collision inventory raw file | SHA-256 of exact bytes, no JSON canonicalization | `sha256:af6c8d1a1630efb54c404aebf4410448ad5c324075b30413c6c12a0ce94e888f` |
| Private blueprint raw file | SHA-256 of exact bytes | `sha256:0c4f2ff6db43128a604e8c824e5ae07acb76f410c4ffdf093d0ed7ebf76d66f0` |
| Independent blueprint-audit raw file | SHA-256 of exact bytes | `sha256:ac776a68cb1ad6b15fe82fe0115f2213b383bc6d44fade651872d41928bb4cbc` |
| Independent label-audit raw file | SHA-256 of exact bytes | `sha256:ae3fccb58b4dc7b2e2504f7a6ae3995ce98b309cc8dd879a037467cacb76cf81` |
| Private mechanical-audit raw file | SHA-256 of exact bytes | `sha256:25a1ac529a49d29cbd8aa7048919e235570a4c06571cce7a366bef2308b188bd` |
| Protocol v4 | SHA-256 of LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `corpusctl digest` on `worlds/dev-repo/world.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json` | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` (both) |
| Production/authoring world, raw bytes | SHA-256 of file bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` (byte-identical) |
| Current grader | `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest raw bytes | SHA-256 of `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| World-build digest | SHA-256 of Go `json.Marshal` of the typed `build` record, not of the whole manifest | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Private-label archive | SHA-256 of LF-terminated UTF-8 index lines `<archive-relative-path>\tsha256:<hex>` for 240 label files, sorted by ordinal path bytes | `sha256:c66d73e6d3db16b52feac05b6ceacc5ffb441b6c6ab4fec8705a87ea19d076bc` |
| Validation sealed manifest, dry-seal stdout | SHA-256 of raw `corpusctl seal` stdout without `--write`; byte-equal to committed `experiments/frontier-v1/manifests/validation.json` | `sha256:4bab07c81226cd1373424dd24461f17588d0151b7dc38da6f16ce14d443d1304` |
| Held-out sealed manifest, dry-seal stdout | same procedure against `held_out.json` | `sha256:a93beedd292bb6128a6289cfde8e3b760b41360c8cef2faa70d0f8dc60cc5402` |

Ancestry: detached HEAD `7a57479…`; sole parent `e705560f…`; that commit's sole parent `c852101e…`; freeze `4b514712…` is an ancestor of the evaluator base. Independently generated index bytes equal `notes/private-label-archive.index.txt`. Canonical `corpusctl digest` of every private label matches the corresponding public registry row (240/240).

`pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, `may_open_held_out: false`.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Commit identity, ancestry, patch equality, two-commit boundary | Pass | Detached at `7a57479…`; required parents; freeze is an ancestor; LF materialization left status empty; two independent raw diff captures equal the stored patch (length 1366164); 884 paths = 480 corpus + 400 fixture + four registry/manifest paths, all under approved tranche locations; payload excludes the public report; no protocol, world, runner, grader, arm, analysis, freeze, gate, schedule, trial, outcome, or full-label bytes in the patch. |
| 2. Pinned contract and reproducibility | Pass | Protocol, both world identities, grader, world-build manifest raw hash, and typed `world_build_digest` match. `go test -count=1 ./...` exit 0; `go vet ./...` exit 0; both dry-seals exit 0 and byte-identical to committed manifests. Checkout remained clean. |
| 3. Frozen allocation and identifiers | Pass | Each tranche: 120 cases, 24 families, five cases each, 120 fixtures, 120 registry rows; 8/8/8 groups; class totals 24/24/12; mix across-group counts and per-family cascade/adversarial floors hold; distinct-family floors hold. IDs are `validation_family_101`–`124` and `held_out_family_301`–`324`; 240 unique case IDs; none collide with authoring, retired-v3, or rejected-v4 public IDs. |
| 4. Five real variants and family independence | Fail | Same-class triples do not collapse to one path/fixture (mechanical five-variant screen: 0 weak groups; 24 unique cascade graphs; unique direct-group and recovery-group three-case action signatures). Independence fails: one global fixture generator, mirrored mix layouts, cloned cascade goals (P1). |
| 5. Five collision sets and collision-inventory integrity | Pass (identifiers / exact strings) / Fail (structure) | Inventory raw digest matches. Independently recomputed Git-tree family IDs, case IDs, and goals agree with all five records (488/488 goals; 248/248 fixture raw digests). Extra inventory public paths are the withheld-label and manifest files this review collected separately from Git. Exact goals, IDs, and fixture bytes are disjoint from the four historical sets and across replacement tranches. Structural comparison still finds a shared generator and mix-family mirrors (P1). Historical and rejected private archives were not opened. |
| 6. Public opacity and prose | Fail (cascade cueing) | Registries contain only `case_id`, `class`, `label_digest`, `grading_script`. No full labels, acceptable paths, rationales, or grader recipes in public cases, fixtures, manifests, or the report. Adversarial text lives in `project-notes/.../request-*.txt` as ordinary untrusted task data and does not name the class or instruct the grader. All 24 cascade goals synonymize the reconcile-from-authority recipe (P1). No grammar/capitalization defects in the 240 public goals. |
| 7. Protocol-v4 executable semantics | Fail | Cases pin the canonical world digest and schema-validate through dry-seal. `state_changes` are deterministic relative-path replacements; all 24 stale-frontier cases declare two replacements after executable act 0. Absence goals request capabilities outside the twelve frozen tools and use an empty acceptable path. Cascades have ≥3 executable verbs and do not omit `tests.list` between `tests.run` and `tests.focus`. Eighteen cascades still append compile/test steps that already succeed on the unedited fixture (P1). |
| 8. Complete private-label and archive review | Pass (mechanical) / Fail (content) | 240/240 labels, one per public case, schema-valid via canonical digest, grader pin correct, class matches registry, path/seq compatible including shorter listed paths, check-kind totals match the candidate report exactly, zero blinded-human checks, 24 absence labels use `[[]]`. Independent content review finds the global label skeleton, cloned rationales, and the eighteen padded cascade graphs (P1/P2). |
| 9. Audit history, custody, and protected artifacts | Fail (audit support) / Pass (custody) | Blueprint auditor inspected the claimed 48/240/1128/576/24/24/288 counts and recorded APPROVE; that APPROVE is not supported for generating-template independence or cascade causality. Label auditor inspected 240/240 and recorded APPROVE; 12/12 prior disagreement categories are closed in the audit record, with 0 open findings. Full labels exist only under the new private root. No historical or rejected private archive was inspected. No model, arm, schedule, trial, or outcome influenced this review. Protocol, runtime/prompts, arm schemas, Arm B, worlds, runner, grader, analysis, freeze entries, and gates are byte-identical to the evaluator base. |

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

Mechanical screen over all 48 families: no same-class group of size ≥2 shares one acceptable-path set, one fixture digest, and one file-topology. Direct-group three-case action signatures are unique across the sixteen direct families. Recovery-group three-case signatures are unique across the sixteen recovery families. Cascade graphs are 24/24 unique, with lengths 4 (5), 5 (12), 6 (6), and 7 (1).

That is not family independence. Variant axes inside a family are mostly different verbs over the same generator (find vs suite vs focus; reconcile vs adversarial note vs absence). Unique hashes and unique graphs do not make 48 generating ideas.

### Pairwise and cross-tranche screens

Recomputed independently over 1,128 unordered family pairs and 576 validation–held-out pairs.

- Identifier, exact-goal, and raw-fixture-digest overlap with all five collision records: none.
- Slot-aligned class/action overlap is not the right independence test. After stripping domain nouns and family numbers, mix-family file-role sets collide, and every `validation_family_N` has the identical slot class sequence as `held_out_family_(N+200)` (24/24 pairs, including all eight mix pairs where that layout is not forced by the protocol).
- Candidate token-similarity pairs at ≥0.50 were inspected. The higher-signal collisions are the numbered mix mirrors and the global generator, not those 14 pairs alone.
- Authoring label check-id tuples do not recur in the replacement labels.

Representative public mirrors, noun-swap identical except for domain token, family number, and snapshot hex:

- `experiments/frontier-v1/fixtures/validation/validation_family_120_v1.json`
- `experiments/frontier-v1/fixtures/held_out/held_out_family_320_v1.json`

## Collision, opacity, and cascade evidence

All 240 replacement fixtures contain:

- template grammar `Repository snapshot <hex> for <domain> with local source, records, and verification metadata.`
- `go.mod` under `example.com/revision/<family_id>`
- `internal/f<NNN>/state.go` plus `state_test.go` with `TestInvariant<NNN>V*` / `TestSibling<NNN>V*`
- `tests/renames.json` with an empty rename list
- `project-notes/<domain>/request-*.txt` whose body is the same two-sentence commit-everything request, only the decoy token changing

All 24 cascade goals match `Reconcile K…X so the <domain> agrees with its authoritative local record, preserve unrelated files, and …`. The trailing clause is compile, named-invariant, or “confirm the corrected repository record”. That synonymizes the success recipe. Public examples:

- `experiments/frontier-v1/corpus/validation/validation_69c2c7dd.json`
- `experiments/frontier-v1/corpus/validation/validation_0ae75943.json`
- `experiments/frontier-v1/corpus/held_out/` cascade rows for families 317–324, same grammar

Absence example, capability genuinely missing from the twelve tools: `experiments/frontier-v1/corpus/validation/validation_6f2fb5ad.json` (remote registry lookup). Stale-frontier example with a shared after-act-0 event: `experiments/frontier-v1/corpus/validation/validation_963b1c43.json`.

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
| private-archive index equality | 240 lines; byte-equal to stored index |
| collision-inventory fixture raw digests vs Git objects | 248/248 |
| collision-inventory goals vs Git objects | 488/488 |
| labels in the review checkout | authoring only (8 files); no validation/held-out labels |
| fixture materialization + `go test ./...` | 6/6 already green: `validation_family_101_v2`, `validation_family_101_v3`, `validation_family_102_v1`, `validation_family_109_v1`, `validation_family_120_v1`, `held_out_family_320_v1` |
| cascade `tests.run`→`tests.focus` missing `tests.list` | 0 |
| cascade graphs with <3 executable verbs | 0 |
| path/seq incompatibilities | 0 |
| stale-frontier cases without `state_changes` | 0 |
| check-kind totals | `act_sequence` 216, `act_output_matches` 214, `act_path_absent` 72, `file_matches` 95, `act_status` 37, `final_message_matches` 24, `act_count` 24, `final_message_states` 0 |
| acceptable-path compatibility | 240/240 listed paths have a position for every path-bound `seq` |
| checkout after checks | clean |

The six materialized suites passing includes the two mix-cascade snapshots whose public goals require a later compile or invariant. `internal/f120/state.go` already holds `good_120_1`; `records/120/index.txt` is not imported by that test. The compile/test step is therefore not evidence produced by the reconcile.

## Findings

### P1 — One parameterized generator instantiated 240 times

Every replacement fixture is the repository-snapshot generator above with a different domain noun and family number. Mix, direct, and recovery families share that skeleton; extra files (`release/channels.toml`, `harbor/codes.tsv`, `attestations/by-day/…`) are still the same `key=K…X` / `state=bad_*` / `required=good_*` record. Unique filenames and fixture hashes do not create 48 generating templates.

### P1 — Held-out mix families are numbered mirrors of validation mix families

`validation_family_117`–`124` and `held_out_family_317`–`324` have identical five-slot class sequences pair-wise. Direct and recovery groups must share the protocol layout, so those 101↔301 and 109↔309 matches are not by themselves a defect. Mix layouts are not forced. Eight of eight mix pairs are mirrored, and the 120/320 fixture pair is a noun swap of the same seven paths.

### P1 — Cascade goals cue the success recipe

All 24 cascade goals tell the agent to reconcile an opaque key so a named domain “agrees with its authoritative local record,” preserve unrelated files, then compile, pass an invariant, or confirm the corrected record. That is the edit-from-authority recipe without verb names. Task 0.4 and the review prompt reject goals that synonymize or strongly cue the action recipe.

### P1 — Eighteen cascades pad the graph with an already-green verify step

Opaque IDs in this class include `held_out_25215c38`, `held_out_272bff15`, `held_out_2bc6c12e`, `held_out_76f0e320`, and the validation mix cascades whose public goals end in compile or invariant language (for example `validation_69c2c7dd`, `validation_5b3326f5`). Six cascades do edit a `.go` source the later build/test can observe. The other eighteen edit a record/text path and then require `repo.build` or `tests.*`. The prompt rejects optional steps inserted only to make a graph unique. Unique graphs were bought that way.

### P1 — Prior APPROVE audits are not supported on these points

The blueprint audit’s pairwise screen scored slot-aligned class/action overlap (max 2) and then APPROVE. It does not inspect the public fixture generator or mix class-layout mirrors. Round 6 closed a finding that later build/suite acts consume repaired state; independent materialization of `validation_family_120_v1` and `held_out_family_320_v1` shows the suite already green. The label audit recorded `cascade_recipe_cues: 0` and closed F-002 on the same causality claim. Those APPROVE records remain useful as mechanical coverage evidence. They are not proof of independence or cascade causality.

### P2 — Repeated label skeletons

All 240 rationales open with the same frozen-world framing sentence, then add case-specific paths. Check-id tuple `required_path` + `authoritative_value_observed` appears 32 times; `repo.find`/`repo.read` is the listed path on 22 labels. Absence labels correctly share one empty-path shape (24). The shared non-absence skeleton is a labeling-template leftover of the fixture generator. Not independently blocking if P1 generator/mirror/cascade defects are repaired, because those repairs force new labels.

### P2 — Candidate report overclaims mix topology uniqueness

The public report states 16 unique mix topologies and that high-similarity pairs were not mirrored families. Public fixtures for `validation_family_120` and `held_out_family_320` are the same topology. Any revised report must describe the generator that is actually in the tree.

## Accepted limitations

- Six fixtures were executed under `go test`; the other 234 were inspected as files and through dry-seal. The generator is uniform enough that the cascade padding result is not a sampling artifact.
- Pairwise scoring used public file-role topology, tokenized goals, and private action-graph overlap. It did not copy the candidate’s unpublished similarity threshold as the decision rule.
- Authoring labels were read only for collision of check-id structure.

## Custody status and outcomes

Full labels, blueprints, blueprint audit, label audit, mechanical audit, collision inventory, and the private archive index exist only under `D:\Work\personal\phoenix-evaluator-private-v4-revision`. They are absent from the review checkout, the implementation workspace, the public patch, and this record except as digests, counts, opaque IDs, and non-revealing defect categories.

Historical protocol-v3 and rejected-v4 private archives were not hashed, mounted, or opened.

No model, arm, schedule, trial, prospective grader result, validation result, or held-out result was run or observed. This review does not open Gate 1A.

## Smallest next artifact

Keep the package unimported. Do not treat a local edit of this generator as a new corpus.

The smallest replacement that can pass this gate is a new payload candidate whose 48 families are not instances of one repository-snapshot template, whose mix families are not +200 class-layout mirrors, whose cascade goals do not state the reconcile-from-authority recipe, and whose cascade verify steps fail on the unedited fixture and pass only after the required earlier edit.

That replacement changes public cases and fixtures. It therefore requires new private labels, new label digests, a new private-archive identity and stored index, new `validation-label-digests.json` / `held_out-label-digests.json`, new sealed manifests, a new payload commit parented on `c852101e…`, a new public patch with a newly captured raw SHA-256, and a new report-only child commit. Prior blueprint and label audits cannot be reused; they approved this generator.

Acceptance of a later package would still not authorize import, schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run. The chair records any import decision separately.
