# Frontier v1 sealed corpus candidate

Status: unopened candidate for Task 0.6 inspection. This package does not authorize Phase 1, an outcome run, or import into the implementation workspace.

## 1. Evaluator identities and independence

| Role id | Work | Phoenix material seen |
| --- | --- | --- |
| `evaluator.corpus_author` | Created the 48 generating families, 240 runnable cases, 160 fixtures, and first-pass labels | Task 0.4 in `docs/plans/2026-08-17-phoenix-world-plan.md`; `experiments/frontier-v1/protocol.json`; `experiments/frontier-v1/schema/`; `spec/world.schema.json`; `experiments/frontier-v1/worlds/authoring.dev_repo.json`; `experiments/frontier-v1/corpusctl/`; the frontier-v1 README directory contract. Authoring template strings and one authoring case/label/fixture were used only as a collision set and schema example. |
| `evaluator.label_auditor` | Second-pass mechanical audit of allocation, opacity, schema lengths, and check-kind counts | Same contracts as the author. Did not read Phoenix daemon, verb, frontier, or surface implementation. Did not see arm output, trial records, or any prospective validation or held-out outcome. |
| Phoenix implementer | None | Did not create cases, write labels, or adjudicate labels. |

`evaluator.label_auditor` is not the Phoenix implementer. Case author and first-pass labeler were the same evaluator role. That limit is recorded here; it is not dual-control labeling by two people.

Neither role used a Phoenix implementation, an arm, a trial result, or a sealed-tranche outcome while designing cases or labels. Authoring families were treated as a collision set, not as templates to copy.

## 2. Base commit and public-patch digest

- Evaluator checkout: git worktree at `D:\Work\personal\phoenix-evaluator-corpus`, detached `cd84e888a362130c5dfabd240c9a50ba7406a588` (`docs: prepare Phase 0 blocker handoffs`).
- Public payload digest (canonical digests of the 404 case, fixture, and manifest files, sorted by path): `sha256:89698a43ebbce508c37b2abfab5da8b11d2f95da6e2663e38969246cf79e6918`.
- Public patch path: `D:\Work\personal\phoenix-evaluator-private\frontier-v1-sealed-corpus.patch`. The unified-diff digest is written to the custody record after the patch is frozen.

The implementation workspace at `D:\Work\personal\phoenix` was not modified. This candidate was not applied there.

## 3. Counts

Each tranche independently:

| Quantity | validation | held_out |
| ---: | ---: | ---: |
| Cases | 120 | 120 |
| Generating families | 24 | 24 |
| Fixtures | 80 | 80 |
| Withheld label digests | 120 | 120 |

Class coverage, each tranche:

| Class | Cases | Distinct families | Protocol floor |
| --- | ---: | ---: | --- |
| `direct` | 24 | 8 | 24 cases / 8 families |
| `recovery` | 24 | 8 | 24 cases / 8 families |
| `cascade` | 12 | 8 | 12 / 8 |
| `adversarial_text` | 12 | 8 | 12 / 8 |
| `far_discovery` | 12 | 12 | 12 / 8 |
| `temptation` | 12 | 12 | 12 / 8 |
| `absence` | 12 | 12 | 12 / 8 |
| `stale_frontier` | 12 | 12 | 12 / 8 |

Family groups, each tranche: 8 direct, 8 recovery, 8 mix.

Fixture split, each tranche: 16 direct-group snapshots (2 per family), 24 recovery-group snapshots (3 per family), 40 mix-group snapshots (5 per family). 16 + 24 + 40 = 80.

Variant evidence: every family has 5 cases. Every family has 5 distinct goals. Direct-group families vary the snapshot between v1 and v2; recovery-group families use three snapshots; mix-group families use five. No family reuses a goal string.

## 4. Family-allocation matrix

Family numbers were shuffled so the integer suffix does not encode the group. Class names below are coverage metadata from the sealed manifests. They are not answers.

### Validation

| Family | Group | Class counts |
| --- | --- | --- |
| `validation_family_001` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_002` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_003` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_004` | mix | cascade 1, adversarial_text 2, far_discovery 1, stale_frontier 1 |
| `validation_family_005` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_006` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_007` | mix | cascade 1, adversarial_text 2, far_discovery 1, temptation 1 |
| `validation_family_008` | mix | cascade 1, adversarial_text 2, absence 1, stale_frontier 1 |
| `validation_family_009` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_010` | mix | cascade 2, adversarial_text 1, far_discovery 1, temptation 1 |
| `validation_family_011` | mix | cascade 1, adversarial_text 2, absence 1, temptation 1 |
| `validation_family_012` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_013` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_014` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_015` | mix | cascade 2, adversarial_text 1, far_discovery 1, absence 1 |
| `validation_family_016` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_017` | mix | cascade 2, adversarial_text 1, absence 1, stale_frontier 1 |
| `validation_family_018` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_019` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_020` | mix | cascade 2, adversarial_text 1, temptation 1, stale_frontier 1 |
| `validation_family_021` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_022` | direct | direct 3, absence 1, stale_frontier 1 |
| `validation_family_023` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `validation_family_024` | direct | direct 3, absence 1, stale_frontier 1 |

Validation mix class totals: cascade 12, adversarial_text 12, far_discovery 4, temptation 4, absence 4, stale_frontier 4. Every mix family has at least one `cascade` and one `adversarial_text`.

### Held-out

| Family | Group | Class counts |
| --- | --- | --- |
| `held_out_family_001` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_002` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_003` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_004` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_005` | mix | cascade 1, adversarial_text 2, absence 1, stale_frontier 1 |
| `held_out_family_006` | mix | cascade 2, adversarial_text 1, far_discovery 1, temptation 1 |
| `held_out_family_007` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_008` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_009` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_010` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_011` | mix | cascade 1, adversarial_text 2, far_discovery 1, temptation 1 |
| `held_out_family_012` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_013` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_014` | mix | cascade 1, adversarial_text 2, far_discovery 1, stale_frontier 1 |
| `held_out_family_015` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_016` | mix | cascade 1, adversarial_text 2, absence 1, temptation 1 |
| `held_out_family_017` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_018` | mix | cascade 2, adversarial_text 1, temptation 1, stale_frontier 1 |
| `held_out_family_019` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_020` | mix | cascade 2, adversarial_text 1, far_discovery 1, absence 1 |
| `held_out_family_021` | recovery | recovery 3, far_discovery 1, temptation 1 |
| `held_out_family_022` | direct | direct 3, absence 1, stale_frontier 1 |
| `held_out_family_023` | mix | cascade 2, adversarial_text 1, absence 1, stale_frontier 1 |
| `held_out_family_024` | direct | direct 3, absence 1, stale_frontier 1 |

Held-out mix class totals match the validation mix totals. Every mix family has at least one `cascade` and one `adversarial_text`.

### Independent arithmetic

For one tranche:

```text
8 direct families  x 5 = 40
8 recovery families x 5 = 40
8 mix families     x 5 = 40
                       ---
                         120

direct:           8 x 3 = 24
recovery:         8 x 3 = 24
cascade:                  12  (mix only)
adversarial_text:         12  (mix only)
far_discovery:    8 + 4 = 12
temptation:       8 + 4 = 12
absence:          8 + 4 = 12
stale_frontier:   8 + 4 = 12
                       ---
                         120
```

`seal` does not enforce these counts. They were checked from the generated manifests and from the outcome-free class lists used to build the digest registries.

## 5. Cross-tranche disjointness

`corpusctl seal` rejects exact template-string reuse and byte-identical fixture file maps across authoring, validation, and held_out. Both seal commands exited 0, so those collisions are absent.

Independent of `seal`: 48 generating-template strings, all unique. Direct, recovery, and mix groups use different primary task shapes; validation and held-out families use different domains and file layouts rather than renamed copies of the same answer path.

## 6. Commands and exit status

Run from `experiments/frontier-v1/corpusctl` in the evaluator checkout:

| Command | Exit |
| --- | ---: |
| `go test ./...` | 0 |
| `go vet ./...` | 0 |
| `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | 0; printed `sha256:3bba60f0f94f6686c0046c73c8507611e937083d57b66fd022bbef7cf32ca92c` |
| `go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json --write` | 0; wrote `experiments/frontier-v1/manifests/validation.json` |
| `go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json --write` | 0; wrote `experiments/frontier-v1/manifests/held_out.json` |

Every withheld `grading_script` pins the grader digest above. The world digest on every case is `sha256:afa4896d40a68177e86c51574e940ff16633076434acd57a16e00c22e787f1b0`.

A repository-wide walk of the evaluator checkout found no schema-valid validation or held-out label documents. Full labels live only in the private archive.

## 7. Checks by kind

Counted from the 240 private labels, not from the sealed manifests (`seal` cannot see check bodies):

| Kind | Count | Bucket |
| --- | ---: | --- |
| `act_sequence` | 216 | scripted |
| `act_count` | 24 | scripted |
| `act_path_absent` | 24 | scripted |
| `act_output_matches` | 264 | regex |
| `final_message_matches` | 72 | regex |
| `final_message_states` | 0 | blinded-human |
| Total | 600 | |

Scripted 264, regex 336, blinded-human 0. No label uses `final_message_states`.

## 8. Disagreements

Second-pass audit found one defect before sealing: a mix-family rationale was 38 characters, below the 40-character schema minimum. It was lengthened and the tranche regenerated. No other disagreements were recorded. Remaining disagreements: 0. All recorded disagreements were adjudicated before `seal`.

## 9. No outcome observation

No agent was launched. No arm was run. No trial record was written. No validation or held-out outcome exists. Neither tranche is opened.

## 10. Private-label archive

- Location: `D:\Work\personal\phoenix-evaluator-private\`, outside every Phoenix checkout.
- Contents: 240 full labels, class-to-template working notes, adjudication log, generator, and the public patch. None of those private files are in the public patch except as this report's outcome-free counts.
- Custodian: project chair (workspace owner of this machine). Do not copy the archive into a Phoenix checkout; `seal` will reject the workspace.
- Canonical label-set digest (240 private labels, sorted by path): `sha256:e1170e0cf99a33d86ba6ca6a533e057b58d1e1021375b50a79700b0e698c2009`.
- Companion custody record: `D:\Work\personal\phoenix-evaluator-private\notes\digests.json`.

Reviewers should inspect this candidate in the evaluator checkout or from the public patch. Import into the implementation workspace remains a chair decision after Task 0.6.
