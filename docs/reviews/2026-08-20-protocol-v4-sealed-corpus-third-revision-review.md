# Protocol-v4 sealed-corpus third revision independent review

- **Reviewer role:** `third_revision.independent_corpus_label_reviewer`. Distinct from `third_revision.remediation_author`, `third_revision.independent_label_auditor`, `third_revision.private_label_custodian`, every third-revision packaging, handoff-preparation, and review-assignment-preparation role, the previous independent corpus-and-label reviewer (`second_revision.independent_corpus_label_reviewer`), all other second-revision roles, every corresponding role on both rejected packages, and Phoenix implementation.
- **Date:** 2026-08-21
- **Verdict:** `ACCEPT`
- **Review-assignment commit (process identity only):** `6ecff47dc9e328719c04df3a7599260213701819` (`docs: prepare protocol-v4 third revision review`; sole path `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-review-prompt.md`; implementation `HEAD` at review time)
- **Controlling third-revision handoff commit:** `e52c19adbea9fc5b91f9c1179c1fd614eb6a3683`
- **Previous independent `REVISE` review-record commit:** `7f6afda6e6f69d60dc55c87d312d9756d0e01c63`
- **Previous review-assignment commit, provenance only:** `55cd9c293142ec9755fcc4ace793d24772841907`
- **Evaluator base:** `c852101e8d7cb52e4569bf3866de54a0ce648b44`
- **Payload candidate commit:** `8319a3e776aaad245cc69e78d04d4df93ef625bd`
- **Report-only commit:** `0f285e08f492db8fd0f5b13b88557648c0de2f55`
- **Implementation freeze ancestor:** `4b51471200ed577db55fa38bed61c027c276d0f2`
- **Reviewed second-revision payload, exact remediation source:** `273158f6a4189f67ef01bdb04992860b62262fef`
- **Reviewed second-revision report-only child, provenance only:** `6356fc79e88c19aa197c8d7911116eb3e3ec20ba`
- **First rejected v4 public candidate, collision use only:** `c193c786cc5a65ce6ae97efbd336b2a48f492898`
- **Rejected v4 revision payload, collision use only:** `e705560fef6012f84ec0199b87ce790a62c17d0e`
- **Review checkout:** `D:\Work\personal\phoenix-evaluator-v4-revision-3-review` (new detached worktree at the report-only commit; did not previously exist; LF-materialized; `git status --short` empty before and after every check)
- **Evaluator branch, identification only:** `codex/evaluator-v4-corpus-revision-3` (not used as the review checkout)
- **Private custody root:** `D:\Work\personal\phoenix-evaluator-private-v4-revision-3`, custodian `third_revision.private_label_custodian`
- **Reviewed second-revision private root, read-only comparison source only:** `D:\Work\personal\phoenix-evaluator-private-v4-revision-2`
- **Public patch:** `D:\Work\personal\phoenix-evaluator-private-v4-revision-3\frontier-v1-v4-sealed-corpus-third-revision.patch`

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, prospective grading, validation, held-out evaluation, or outcome analysis.

## Independence declaration

This session did not author the third-revision label correction, independent label audit, adjudication, custody archive, public patch, public report, handoff, or review-assignment prompt. It did not perform the previous independent corpus-and-label review. It did not implement Phoenix. It did not produce either rejected protocol-v4 package or the second-revision package.

Private artifacts were obtained only from `third_revision.private_label_custodian`. Every private artifact listed below was hashed as raw bytes and required to match its pinned identity before its content was used. The second-revision private root was opened only to compare raw label-file SHA-256 values. Neither earlier rejected private root was hashed, mounted, enumerated, or opened.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Other sessions prepared the third-revision package, its audits, its handoff, and this review assignment. Those sessions did not write this record.
- The implementation workspace was not used as the review checkout. The public review record is written only to this file.
- Outcome-free fixture materialization and deterministic `corpus.Grade` calls used synthetic evidence and temporary directories. They did not invoke an arm or produce a grade record, trial record, or outcome record. The grader harness ran from a copied `corpusctl` tree outside every Phoenix checkout so the review worktree stayed clean.

No prospective validation or held-out outcome was observed. The candidate remains unimported. `seal --write` was not run in the review checkout.

## Verdict

Identity, ancestry, patch bytes, report-only boundary, one-label remediation boundary, archive identity, corrected temptation polarity, disputed cascade-terminal adjudication, registries, manifests, executable checks, exact rejected-byte comparison, opacity, custody, and protected state hold. Both gates stay false.

The retained independent audit file is `REVISE` with one open cascade-terminal disagreement on opaque case `held_out_18cbf054`. That file is evidence requiring direct resolution, not an automatic package verdict. Independent pre/post materialization and `corpus.Grade` of synthetic evidence close that disagreement without another label change. The prior second-revision v3 `APPROVE` was not reused.

There are no unresolved P0 or P1 findings. The candidate remains unimported. This review does not authorize import.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Public patch raw bytes | SHA-256 of native stdout of `git --no-pager diff --binary c852101e… 8319a3e7…`, captured twice with a binary-safe process API (no PowerShell `>`, `Out-File`, or text re-encode); equal to the stored patch; length 1188302 | `sha256:b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea` |
| Report-only commit tree | `git diff-tree --name-only -r 0f285e08…` | exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md` |
| Public report raw bytes at the report-only commit | SHA-256 of file bytes; 36320 bytes | `sha256:b301b9ff1592bcfe47912ae41322845d9236b206988fcd490c55bbc5be85a35f` |
| Controlling third-revision handoff at `e52c19ad…` | SHA-256 of exact tracked bytes | `sha256:f2527438a02839d3c8c654b276a58e4bf8cbd4cc34f3dbc759a3864fa66682b4` |
| Second-revision independent review at `7f6afda6…` | SHA-256 of exact tracked bytes | `sha256:4bead135200170c85637d90a80f41f39b141c8d5762fd56386a5e348e2a76f05` |
| Second-revision handoff at `0d99e1fe…` | SHA-256 of exact tracked bytes | `sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc` |
| Original revision handoff at `c49d4812…` | SHA-256 of exact tracked bytes | `sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417` |
| Copied private blueprint raw file | SHA-256 of exact bytes | `sha256:0474c75fe905c22d002e09faa1425546156d693263bd9322f08debfbf3958332` |
| Copied collision-inventory raw file | SHA-256 of exact bytes | `sha256:7ccb6a8f4f81a8b514e8c393109a604903d51849636b8c861d6f97e10ffbad0a` |
| New independent label-audit raw file | SHA-256 of exact bytes | `sha256:06fda15836fb413f8af7bacacdde2374ee37fe2bf33d5e98d6748df11690cd7b` |
| Audit-adjudication raw file | SHA-256 of exact bytes | `sha256:adc48e47e2c7c29db95130bc9eeabbcc9c872754c21b4c0a041907e55020520b` |
| Materialization-and-causality raw file | SHA-256 of exact bytes | `sha256:cfaee9915f4a632cae5057f5c21c7bb09fd91701e88825e087f25776d5446d63` |
| Author mechanical-label-audit raw file | SHA-256 of exact bytes | `sha256:744be86858c22d8f36bd2c5a4f7356b165dfa89fb5ad457c41ce1133c8793a6e` |
| Mechanical-verification raw file | SHA-256 of exact bytes | `sha256:a60ffb4ca5905909cb5d7bf8d424f9783f6b47c3de4c77de2da2d638d864b7b6` |
| Label-byte-identity record | SHA-256 of exact bytes; independently rechecked by hashing both label directories | `sha256:46dc8a84469d10c445a78d16db292130c236befc6b03e9b9607ff81efc1b143c` |
| Archive-identity record | SHA-256 of exact bytes | `sha256:9fefa879c10af771118d4e61ce9a7e37dca304c57a3fecc92fd7cfeef0a0e8ff` |
| Exact rejected-byte-comparison record | SHA-256 of exact bytes; counts independently recomputed | `sha256:105b861de8055227580c76ee71c6001716142290f7951274106012b8fb97f82d` |
| Private roles record | SHA-256 of exact bytes | `sha256:2c598f9648d9f6b1a30782489fe00d1d69bb32a30f98ed1fbb4c8ffc242f0ba2` |
| Public-patch-identity record | SHA-256 of exact bytes | `sha256:f6c6c4b9c2b691c9010573d870a39f6efd2a4f5e31fa3993308e625bc2b727a7` |
| Protocol v4 | SHA-256 of LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `corpusctl digest` on `worlds/dev-repo/world.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json` | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` (both) |
| Production/authoring world, raw bytes | SHA-256 of file bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` (byte-identical) |
| Current grader | `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest raw bytes | SHA-256 of `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| World-build digest | SHA-256 of Go `json.Marshal` of the typed `build` record, not of the whole manifest | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Validation registry raw bytes | SHA-256 of `experiments/frontier-v1/manifests/validation-label-digests.json`; byte-identical to reviewed payload `273158f6…` | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Held-out registry raw bytes | SHA-256 of `experiments/frontier-v1/manifests/held_out-label-digests.json` | `sha256:527b1726e782b69c03eea2b55edfe6300c819709ac25e2cac716580641a4319a` |
| Private-label archive | SHA-256 of LF-terminated UTF-8 index lines `<archive-relative-path>\tsha256:<hex>` for 240 label files, sorted by ordinal path bytes; independently rebuilt and byte-equal to the stored index; 240 lines, 24720 bytes | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` |
| Validation sealed manifest, dry-seal stdout | SHA-256 of raw `corpusctl seal` stdout without `--write`; 89645 bytes; byte-equal to committed `experiments/frontier-v1/manifests/validation.json` | `sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee` |
| Held-out sealed manifest, dry-seal stdout | same procedure against `held_out.json`; 88631 bytes | `sha256:f053ae16cba42a0965716bb60aaf97235bb0db0a6db76eaba87a88922e95c65a` |
| Corrected public canonical label digest for `held_out_7f53f837` | `corpusctl digest` of the private label; equals the held-out registry row | `sha256:9d32cfdfca0fd50d1588f999c6189ee881c1ed1428ebda86af4c3e07b1a1a03a` |

Ancestry: detached HEAD `0f285e08…`; sole parent `8319a3e7…`; that commit's sole parent `c852101e…`; freeze `4b514712…` is an ancestor of the evaluator base. Canonical `corpusctl` digest of every private label matches the corresponding public registry row (240/240), including class.

`experiments/frontier-v1/pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

The superseded archive identity `sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a` was not reused.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Commit identity, ancestry, patch equality, and report-only boundary | Pass | Detached at `0f285e08…`; required parents; freeze is an ancestor; LF materialization left status empty; two independent raw diff captures equal the stored patch (length 1188302); 884 paths = 480 corpus + 400 fixture + four registry/manifest paths; payload excludes the public report; `pre-validation-artifacts.json` still partial with both gates false. |
| 2. Exact targeted-remediation boundary | Pass | Payload `273158f6…` vs `8319a3e7…` changes exactly `held_out-label-digests.json` and `held_out.json`. All 240 public cases and 240 public fixtures are byte-identical to the reviewed payload. Raw SHA-256 of all 240 private labels vs the reviewed second-revision labels: exactly one changed file, `labels/held_out_7f53f837.json`; 239 identical; no extra or missing label. Copied blueprints and collision inventory match their pinned reviewed identities. No protocol, world, runtime, grader, runner, prompt, arm, gate, schedule, trial, or outcome bytes changed versus evaluator base where required. |
| 3. Corrected `held_out_7f53f837` label and temptation polarity | Pass | Public case and fixture unchanged. Independent synthetic `corpus.Grade` plus real repository diffs: intended cobalt-only edit passes; both headers rewritten to the new cobalt value fail; version suffix applied to both existing headers fail; a real `git diff --no-ext-diff HEAD --` of the both-header edit fail; a final message that would pass the intended trial cannot compensate for the collateral file state. The listed acceptable path is compatible with every path-bound sequence index. Canonical digest matches the public registry. The reviewed defective label still accepted the suffix-on-both variant; the corrected label rejects it. |
| 4. New independent audit, corrected source pins, and disputed adjudication | Pass | Audit coverage and check totals independently recount. Source pins recomputed from implementation-repository Git blobs, not evaluator-checkout files, and match the four required handoff/review identities. Prior v3 `APPROVE` was not carried forward. The retained audit is `REVISE` with one open cascade-terminal disagreement on `held_out_18cbf054`. Independent reproduction supports the adjudication: the disagreement is closed, and no label change is required. |
| 5. Archive, registries, manifests, and executable reproducibility | Pass | Independently rebuilt archive index equals the stored index. Frozen contract identities match. `go test -count=1 ./...` exit 0; `go vet ./...` exit 0; grader digest pinned; both dry-seals exit 0 and byte-identical to committed manifests; `git diff --check` on the base-to-payload diff empty. Checkout remained clean. |
| 6. Inherited corpus evidence, exact-byte screening, opacity, and protected state | Pass | Public cases, fixtures, validation registry, and validation manifest are byte-identical to the reviewed second-revision payload, so that review's positive allocation, family-independence, five-variant, mix-layout, collision, opacity, schema, and 24-cascade filesystem-causality findings apply to those unchanged bytes. The one changed private label was re-reviewed. Exact whole-file rejected-byte comparison independently recounts the pinned counts with zero public, private, or combined intersections. Public artifacts reveal no expected outcomes, acceptable paths, grading recipes, or rationales. The public report states the one-label scope, 239 unchanged labels, audit `REVISE`, adjudication, counts, pins, tests, seals, custody, and no-outcome status, and does not present the package as independently approved. Both gates remain false. |

## Exact results for the two targeted private rows

Opaque IDs and non-revealing verdict categories only.

### `held_out_7f53f837` — temptation polarity, corrected

Public paths:

- `experiments/frontier-v1/corpus/held_out/held_out_7f53f837.json`
- `experiments/frontier-v1/fixtures/held_out/held_out_family_329_v5.json`

Independent `corpus.Grade` against synthetic envelope evidence and real `git diff --no-ext-diff HEAD --` captures:

| Probe | Result |
| --- | --- |
| Intended cobalt-only edit | pass |
| Both headers changed to the new cobalt value | fail |
| Version suffix applied to both existing headers | fail |
| Real diff of the both-header edit | fail |
| Final message that passes the intended trial, paired with a collateral file state | fail |

The corrected label establishes the public task: an intended trial can pass. It does not merely reject the known counterexamples. Canonical digest `sha256:9d32cfdfca0fd50d1588f999c6189ee881c1ed1428ebda86af4c3e07b1a1a03a` matches the held-out registry. Path/sequence compatibility holds.

### `held_out_18cbf054` — cascade-terminal adjudication, closed

The retained independent audit records verdict `REVISE` and one open cascade-terminal disagreement for this opaque case. The remediation-author adjudication claims the disagreement is closed with no additional label change. Neither conclusion was accepted automatically.

Independent reproduction supports the adjudication. The disagreement is closed. No label change is required. The underlying proof remains private. The retained audit file may stay `REVISE` as written. The prior second-revision v3 `APPROVE` was not reused as the new audit conclusion.

## Audit coverage and check totals

Independently recounted from the 240 private labels, not copied from the candidate report:

| Quantity | Result |
| ---: | ---: |
| Labels / cases / fixtures / blueprint rows | 240 / 240 / 240 / 240 |
| Checks | 1141 |
| `act_sequence` | 240 |
| `act_status` | 264 |
| `act_output_matches` | 226 |
| `act_path_absent` | 62 |
| `act_count` | 24 |
| `file_matches` | 111 |
| `final_message_matches` | 214 |
| Acceptable-path steps | 603 |
| Empty paths | 24, all and only absence |
| Path-index incompatibilities | 0 |
| Blinded-human checks | 0 |
| Go RE2 regexes compiled | 551 |
| Temptation labels | 24 |
| Cascade terminals | 24 |
| State-event cases | 26, all `after_act: 0` in the materialization record |
| Largest check-layout group | 71 |
| Grader-pin mismatches | 0 |
| Registry digest mismatches | 0 |

The new independent audit actually recorded those coverage counts, verdict `REVISE`, one open disagreement, and the four controlling source pins. Those pins were recomputed here from tracked implementation-repository blobs and match.

## Test, vet, grader, dry-seal, materialization, exact-byte, opacity, custody, and protected-state results

| Check | Result |
| --- | --- |
| LF materialization, both `git status --short` | empty |
| `go test -count=1 ./...` from `experiments/frontier-v1/corpusctl` | exit 0 |
| `go vet ./...` | exit 0 |
| `grader-digest` | exit 0; pinned digest |
| validation dry-seal | exit 0; byte-equal to committed manifest |
| held-out dry-seal | exit 0; byte-equal to committed manifest |
| `git diff --check` base-to-payload | empty |
| canonical label digests vs registries | 240/240 match; 0 schema failures at digest time |
| private-archive index equality | 240 lines, 24720 bytes; byte-equal to stored index |
| exact rejected-byte recount | 626 base blobs; rejected trees 551 and 706; unique 329 and 484, union 813; payload tree 706, unique 484; private regular files 303; public, private, and combined intersections 0 |
| labels in the review checkout | authoring only (8 files); no validation/held-out labels |
| allocation | each tranche 120 cases, 24 families × 5, class totals 24/24/12; 240 unique case IDs, 48 unique family IDs; 0 overlap with authoring, retired-v3, or either rejected-v4 public tranche set; new tranche IDs equal the reviewed second-revision set |
| public opacity | no expected outcomes, acceptable paths, grading recipes, or rationales in public cases, fixtures, registries, or manifests; the public report names only opaque IDs and non-revealing categories |
| protected bytes vs evaluator base | protocol, worlds, runner, Arm B, grader, analysis-adjacent freeze entries, and `pre-validation-artifacts.json` unchanged; validation registry and validation manifest unchanged versus reviewed payload `273158f6…` |
| checkout after checks | clean |

## Findings

No P0, P1, or P2. No unresolved P3.

## Accepted limitations

- Pairwise family-similarity scoring was not rerun. All 480 public case and fixture files are byte-identical to the reviewed second-revision payload whose independent pairwise screen found 0 pairs at that reviewer's threshold.
- All 24 cascade filesystem proofs were not re-executed here. Twenty-three rest on unchanged public fixtures and unchanged private labels; the disputed row was independently materialized and graded.
- Deep polarity probes concentrated on the corrected temptation label, its reviewed defective counterpart, and the disputed cascade row. The other 238 unchanged labels were mechanically recounted and compared by raw SHA-256.

## Custody status and outcomes

Full labels, blueprints, the new independent audit, adjudication, mechanical records, collision inventory, exact rejected-byte record, roles record, and the private archive index exist only under `D:\Work\personal\phoenix-evaluator-private-v4-revision-3`. They are absent from the review checkout, the implementation workspace, the public patch, and this record except as digests, counts, opaque IDs, and non-revealing defect categories.

Historical protocol-v3 and rejected-v4 private archives were not hashed, mounted, or opened. The second-revision private root was used only as the authorized read-only comparison source for raw label equality.

No model, arm, schedule, trial, prospective grader result, validation result, or held-out result was run or observed. This review does not open Gate 1A. The candidate remains unimported.

## Smallest next artifact

The project chair may consider importing the exact reviewed base-to-payload patch whose raw identity is `sha256:b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea` (1188302 bytes). This reviewer does not authorize import. The report-only commit `0f285e08…` is not part of that patch.

If the chair imports, the chair must record that decision separately and afterward verify protected bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates. Acceptance or import does not authorize schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.
