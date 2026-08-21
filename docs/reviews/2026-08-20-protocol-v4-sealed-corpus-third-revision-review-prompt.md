# Protocol-v4 sealed-corpus targeted third revision independent review prompt

This assignment covers only the targeted two-commit replacement produced from `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-handoff.md`. It does not reopen either rejected protocol-v4 package or the unchanged portions of the reviewed second-revision package. The earlier review prompts remain provenance records and must not be reused as the operative assignment.

Give the text below to a reviewer who is independent of every third-revision remediation, label-audit, adjudication, custody, packaging, handoff-preparation, and review-assignment-preparation role; the previous independent corpus-and-label reviewer; all second-revision roles; both rejected packages; and Phoenix implementation. Use a fresh review session and a new detached review checkout.

Do not issue this assignment while it is untracked. Commit it in the implementation repository, then record that review-assignment commit in the message to the reviewer. The assignment commit is a process identity only; it does not change the candidate package identified below.

```text
Independently review the targeted protocol-v4 third-revision `validation` and `held_out` corpus replacement for `frontier-v1`.

This is a corpus-and-label gate review. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, prospective grading against model output, validation, held-out evaluation, or outcome analysis.

Candidate package

- Controlling third-revision handoff commit: e52c19adbea9fc5b91f9c1179c1fd614eb6a3683
- Previous independent `REVISE` review-record commit: 7f6afda6e6f69d60dc55c87d312d9756d0e01c63
- Previous review-assignment commit, provenance only: 55cd9c293142ec9755fcc4ace793d24772841907
- Fresh detached review checkout: D:\Work\personal\phoenix-evaluator-v4-revision-3-review
- Evaluator branch, identification only: codex/evaluator-v4-corpus-revision-3
- Evaluator base: c852101e8d7cb52e4569bf3866de54a0ce648b44
- Payload candidate commit: 8319a3e776aaad245cc69e78d04d4df93ef625bd
- Report-only commit: 0f285e08f492db8fd0f5b13b88557648c0de2f55
- Required ancestry: report-only commit parent == payload candidate; payload candidate sole parent == evaluator base
- Implementation freeze: 4b51471200ed577db55fa38bed61c027c276d0f2
- Reviewed second-revision payload, exact remediation source: 273158f6a4189f67ef01bdb04992860b62262fef
- Reviewed second-revision report-only child, provenance only: 6356fc79e88c19aa197c8d7911116eb3e3ec20ba
- First rejected v4 public candidate, collision use only: c193c786cc5a65ce6ae97efbd336b2a48f492898
- Rejected v4 revision payload, collision use only: e705560fef6012f84ec0199b87ce790a62c17d0e
- Public payload patch: D:\Work\personal\phoenix-evaluator-private-v4-revision-3\frontier-v1-v4-sealed-corpus-third-revision.patch
- Expected public-patch raw SHA-256: sha256:b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea
- Expected public-patch byte length: 1188302
- Public report at the report-only commit: docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md
- Expected public-report raw SHA-256: sha256:b301b9ff1592bcfe47912ae41322845d9236b206988fcd490c55bbc5be85a35f
- Expected public-report byte length: 36320
- New private custody root: D:\Work\personal\phoenix-evaluator-private-v4-revision-3
- Reviewed second-revision private root, read-only comparison source only: D:\Work\personal\phoenix-evaluator-private-v4-revision-2
- Expected private-label custodian: third_revision.private_label_custodian
- Stored private-label archive index: D:\Work\personal\phoenix-evaluator-private-v4-revision-3\label-archive-index.tsv
- Expected private-label archive index: 240 lines, 24720 bytes, sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4
- Expected copied private blueprint raw SHA-256: sha256:0474c75fe905c22d002e09faa1425546156d693263bd9322f08debfbf3958332
- Expected copied collision-inventory raw SHA-256: sha256:7ccb6a8f4f81a8b514e8c393109a604903d51849636b8c861d6f97e10ffbad0a
- Expected new independent label-audit raw SHA-256: sha256:06fda15836fb413f8af7bacacdde2374ee37fe2bf33d5e98d6748df11690cd7b
- Expected audit-adjudication raw SHA-256: sha256:adc48e47e2c7c29db95130bc9eeabbcc9c872754c21b4c0a041907e55020520b
- Expected materialization-and-causality raw SHA-256: sha256:cfaee9915f4a632cae5057f5c21c7bb09fd91701e88825e087f25776d5446d63
- Expected author mechanical-label-audit raw SHA-256: sha256:744be86858c22d8f36bd2c5a4f7356b165dfa89fb5ad457c41ce1133c8793a6e
- Expected mechanical-verification raw SHA-256: sha256:a60ffb4ca5905909cb5d7bf8d424f9783f6b47c3de4c77de2da2d638d864b7b6
- Expected label-byte-identity raw SHA-256: sha256:46dc8a84469d10c445a78d16db292130c236befc6b03e9b9607ff81efc1b143c
- Expected archive-identity record raw SHA-256: sha256:9fefa879c10af771118d4e61ce9a7e37dca304c57a3fecc92fd7cfeef0a0e8ff
- Expected exact rejected-byte-comparison raw SHA-256: sha256:105b861de8055227580c76ee71c6001716142290f7951274106012b8fb97f82d
- Expected private roles record raw SHA-256: sha256:2c598f9648d9f6b1a30782489fe00d1d69bb32a30f98ed1fbb4c8ffc242f0ba2
- Expected public-patch-identity record raw SHA-256: sha256:f6c6c4b9c2b691c9010573d870a39f6efd2a4f5e31fa3993308e625bc2b727a7
- Controlling third-revision handoff raw SHA-256: sha256:f2527438a02839d3c8c654b276a58e4bf8cbd4cc34f3dbc759a3864fa66682b4
- Second-revision independent review raw SHA-256: sha256:4bead135200170c85637d90a80f41f39b141c8d5762fd56386a5e348e2a76f05
- Second-revision handoff raw SHA-256: sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc
- Original revision handoff raw SHA-256: sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417

Pinned freeze and dependent identities

- Protocol v4, LF-normalized UTF-8: sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9
- Production/authoring world, canonical JSON: sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b
- Production/authoring world, raw bytes: sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143
- Current grader: sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc
- World-build manifest, raw bytes: sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f
- Content-addressed `world_build_digest`: sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d
- Validation registry, raw bytes: sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816
- Held-out registry, raw bytes: sha256:527b1726e782b69c03eea2b55edfe6300c819709ac25e2cac716580641a4319a
- Validation sealed manifest: 89645 bytes, sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee
- Held-out sealed manifest: 88631 bytes, sha256:f053ae16cba42a0965716bb60aaf97235bb0db0a6db76eaba87a88922e95c65a
- Corrected public canonical label digest for opaque case `held_out_7f53f837`: sha256:9d32cfdfca0fd50d1588f999c6189ee881c1ed1428ebda86af4c3e07b1a1a03a

Required reading and authority

Read the following completely before reviewing the candidate. The controlling process records are later than the evaluator base; read their exact Git blobs from the implementation repository at their named commits, not from the detached candidate checkout.

- the controlling third-revision handoff at commit e52c19adbea9fc5b91f9c1179c1fd614eb6a3683;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-review.md at commit 7f6afda6e6f69d60dc55c87d312d9756d0e01c63;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-review-prompt.md at commit 55cd9c293142ec9755fcc4ace793d24772841907;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-handoff.md at commit 0d99e1fe966513e2bb3a8e9ae330eaa75dfe88fc;
- docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md;
- docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md;
- docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md at report-only commit 0f285e08f492db8fd0f5b13b88557648c0de2f55;
- the complete evaluator-base-to-payload diff, report-only diff, and reviewed-second-revision-to-new-payload diff;
- the affected public case, fixture, registry row, and sealed-manifest entry for opaque case `held_out_7f53f837`;
- through private custody, the corrected label, the disputed label for opaque case `held_out_18cbf054`, both labels' reviewed second-revision counterparts, the complete new independent audit, adjudication, materialization evidence, byte-identity evidence, archive index, mechanical records, roles record, copied blueprints and collision inventory, and exact rejected-byte record;
- `experiments/frontier-v1/corpusctl/internal/corpus/grade.go`, schemas, protocol, authoring world, world-build manifest, and relevant `corpusctl` commands.

The tracked second-revision `REVISE` record controls the remediation scope. Its positive findings remain evidence for public and private bytes proven unchanged, but no prior verdict substitutes for independent review of the changed label, regenerated identities, corrected source pins, or the new audit disagreement.

Independence and custody

- Record the reviewer role as `third_revision.independent_corpus_label_reviewer`, the date, material seen, and all independence limitations. Be distinct from every third-revision remediation author, label auditor, adjudicator, custodian, package preparer, handoff preparer, assignment preparer, the previous independent reviewer, all second-revision roles, both rejected packages, and Phoenix implementation.
- Obtain new private artifacts only from `third_revision.private_label_custodian`. Access the second-revision private root only as the authorized read-only comparison source for raw label equality. Do not inspect either earlier rejected private root.
- Verify every private artifact's raw identity before relying on its content. Recompute every embedded source pin from authoritative implementation-repository Git blobs.
- Do not copy a full label, expected answer, acceptable path, rationale, private blueprint detail, causal proof, disagreement detail, or adjudication detail into a Phoenix checkout, Git history, public review record, terminal transcript, or response.
- A private finding may name an opaque case ID and a non-revealing defect category only.
- Do not run the experiment runner, a model, an arm, a schedule, a trial, validation, held-out evaluation, or outcome analysis. Do not generate a schedule or open Gate 1A.
- Outcome-free fixture materialization and deterministic `corpus.Grade` calls using synthetic evidence are allowed. They must not use model output, invoke an arm, or create a trial, grade record, or outcome record.
- Do not import the patch into the implementation workspace. Stop and report a custody breach if any prospective validation or held-out outcome is observed.

Fresh review checkout and LF materialization

Create a new detached worktree at the exact report-only commit. Stop if the target path already exists; do not delete or repurpose it.

    git -C D:\Work\personal\phoenix worktree add --detach D:\Work\personal\phoenix-evaluator-v4-revision-3-review 0f285e08f492db8fd0f5b13b88557648c0de2f55

In the new checkout, require an empty `git status --short`, then run exactly:

    git config --worktree core.autocrlf false
    git config --worktree core.eol lf
    git checkout-index --all --force
    git status --short

Run `checkout-index` only in the new detached review checkout. Stop if either status check is nonempty. Never run `seal --write` in the review checkout. Do not amend, rebase, switch, commit, or modify the evaluator branch.

Verify all sections below.

1. Commit identity, ancestry, patch equality, and report-only boundary

   - Require detached HEAD at report-only commit `0f285e08f492db8fd0f5b13b88557648c0de2f55`.
   - Require its sole parent to be payload `8319a3e776aaad245cc69e78d04d4df93ef625bd`; require the payload's sole parent to be evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`; require the implementation freeze to be an ancestor of the base.
   - Require the report-only commit to add exactly `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md` and no other path. Recompute its 36320-byte raw identity.
   - Independently capture the native stdout bytes of:

         git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 8319a3e776aaad245cc69e78d04d4df93ef625bd

     Use a binary-safe process API or raw stream. Do not use a PowerShell pipeline, `>`, `Out-File`, `Set-Content`, or a text decode/re-encode step. Capture twice and require both captures and the stored patch to be byte-equal, 1188302 bytes, and raw SHA-256 `b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea`.
   - Require the base-to-payload path boundary to remain the complete reviewed corpus replacement: exactly 884 changed paths, comprising 480 corpus paths, 400 fixture paths, and four registry/manifest paths, with no report or private artifact.
   - Confirm `pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

2. Exact targeted-remediation boundary

   - Compare reviewed payload `273158f6a4189f67ef01bdb04992860b62262fef` with new payload `8319a3e776aaad245cc69e78d04d4df93ef625bd`. Require exactly two changed paths: `experiments/frontier-v1/manifests/held_out-label-digests.json` and `experiments/frontier-v1/manifests/held_out.json`.
   - Require all 240 public cases and all 240 public fixtures to be byte-identical to the reviewed payload.
   - Compare all 240 new labels with the reviewed second-revision labels by raw SHA-256. Require exactly one changed file, `labels/held_out_7f53f837.json`, and 239 byte-identical files. Require no extra or missing label.
   - Verify `notes/label-byte-identity.json` independently; do not treat the record as proof by itself.
   - Require the copied family blueprints and collision inventory to match their pinned reviewed raw identities. They are inherited evidence, not new approvals.
   - Any additional public case, fixture, private label, protocol, world, runtime, grader, runner, prompt, arm, gate, schedule, trial, or outcome change exceeds the assignment and requires `REVISE` or `REJECT` according to consequence.

3. Corrected `held_out_7f53f837` label and temptation polarity

   - Inspect the unchanged public case and fixture, the reviewed defective private label, and the new corrected label. Apply `grade.go` path matching and sequence-index semantics exactly, including every shorter acceptable path.
   - Independently construct outcome-free synthetic envelope evidence and, where needed, real repository diffs. Require the intended cobalt-only edit to pass and each class-defining collateral variant to fail: both headers changed to the new cobalt value; the version suffix applied to both existing headers; a real diff containing the both-header edit; and a final message that tries to compensate for the collateral edit.
   - Require every listed acceptable path to remain compatible with every path-bound sequence index. Confirm the corrected label establishes the public task rather than merely rejecting the known counterexamples.
   - Independently recompute the corrected label's canonical digest and require equality with public registry value `sha256:9d32cfdfca0fd50d1588f999c6189ee881c1ed1428ebda86af4c3e07b1a1a03a`.
   - Keep all expected paths, checks, rationales, raw label identities, and answer-bearing details private.

4. New independent audit, corrected source pins, and disputed adjudication

   - Verify the new audit, adjudication, materialization record, mechanical records, roles record, archive record, and byte-identity record against the pinned raw hashes before reading them.
   - Confirm the new independent audit actually covered 240 labels, 240 cases, 240 fixtures, 240 blueprint rows, 1141 checks, 551 Go regexes, 24 temptation labels, 24 cascade terminals, 26 state-event cases, 240 wrong paths, and the repeated-layout groups, including the largest group of 71.
   - Independently recount the 1141 checks: 240 `act_sequence`, 264 `act_status`, 226 `act_output_matches`, 62 `act_path_absent`, 24 `act_count`, 111 `file_matches`, and 214 `final_message_matches`. Require 603 acceptable-path steps, exactly 24 empty paths all and only for absence cases, zero path-index incompatibilities, and zero blinded-human checks.
   - Recompute the audit's source pins from the authoritative tracked blobs. Require the four handoff/review identities pinned above. Reject a similarly named evaluator-checkout file as authority.
   - The independent audit is intentionally retained with verdict `REVISE`, one open cascade-terminal disagreement, and a claimed additional-label scope expansion for opaque case `held_out_18cbf054`. The remediation-author adjudication claims the row fails before the required state change and passes after it, leaving zero open disagreements without changing another label.
   - Do not automatically accept the author adjudication and do not automatically fail merely because the immutable audit retains `REVISE`. Independently inspect the disputed label, case, fixture, blueprint, state event, acceptable paths, and terminal check. Materialize the exact pre-edit and post-edit states and run `corpus.Grade` with synthetic evidence. Decide whether the terminal check fails to establish success before the required change and establishes success after it.
   - Reconcile the audit's `cascade_terminals_passing: 23` with the materialization record's 24 fail-before/pass-after claims. `ACCEPT` is allowed only if the reviewer independently proves the adjudication is correct and records why the retained disagreement is closed. If the label is defective or ambiguity remains, return `REVISE` and require the smallest affected-label replacement package.
   - Confirm the prior second-revision v3 `APPROVE` was not reused as the new audit conclusion.

5. Archive, registries, manifests, and executable reproducibility

   - Rebuild the archive index from all 240 new label files: SHA-256 exact raw bytes; form UTF-8 `labels/<filename>\tsha256:<lowercase-hex>` lines; sort by archive-relative path in ordinal UTF-8 byte order; join with LF; terminate with LF. Require byte equality with the stored index, 240 lines, 24720 bytes, and digest `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4`.
   - Confirm the superseded archive digest `sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a` was not reused.
   - Recompute protocol, world, grader, world-build, registry, and manifest identities pinned above. Stop before substantive private review if a frozen contract identity differs.
   - From `experiments/frontier-v1/corpusctl`, run:

         go test -count=1 ./...
         go vet ./...
         go run ./cmd/corpusctl grader-digest --repo-root ../../..

   - Run both seals without `--write`, capture raw stdout, and require byte equality with their committed manifests:

         go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json
         go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json

   - Run `git diff --check` for the base-to-payload diff. Keep the review checkout clean after every check.

6. Inherited corpus evidence, exact-byte screening, opacity, and protected state

   - Verify the second-revision review record's positive identity, allocation, family-independence, five-variant, mix-layout, collision, opacity, schema, custody, and all-cascade causality findings apply to byte-identical public and private material. Do not silently carry a finding across changed bytes.
   - Re-run schema, allocation, identifier, registry linkage, leakage, public-opacity, path-index, regex, and protected-artifact checks. Re-run the exact whole-file rejected-byte comparison against public payloads `c193c786...` and `e705560f...` using the new payload tree and new private root.
   - Reconcile the reported exact-byte counts: 626 evaluator-base blobs; rejected trees with 551 and 706 blobs; rejected unique sets with 329 and 484 blobs, 813 in union; 706 payload-tree blobs; 484 payload unique blobs; 303 new-private regular files; and zero public, private, or combined intersections.
   - Confirm public artifacts reveal no expected outcomes, acceptable paths, grading recipes, rationales, private causal details, disagreement details, or answer-bearing fixture templates.
   - Confirm the revised public report accurately states identities, the one-label scope, 239 unchanged labels, 480 unchanged public case/fixture files, audit `REVISE`, adjudication, counts, pins, tests, seals, custody, and no-outcome status. It must not present the package as independently approved.
   - Confirm the protocol, runtime and prompts, arm schemas, Arm B, worlds, runner, grader, analysis, accepted freeze entries, `pre-validation-artifacts.json`, and both gates are byte-identical to evaluator base where required. Keep both gates false.
   - Confirm neither prohibited earlier private root was opened, the new package remains unimported, and no model, arm, schedule, trial, prospective model grade, validation result, held-out result, or outcome influenced remediation or review.

Review record and verdict

Write the public review record to:

D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-third-revision-review.md

Return `ACCEPT`, `REVISE`, or `REJECT`, with P0-P3 findings and exact public paths and lines where applicable. For a private finding, use only an opaque case ID and a non-revealing defect category.

The review record must include:

- reviewer role, independence declaration, date, review-assignment commit, reviewed commits and paths, and material seen;
- independently recomputed ancestry, patch identity and equality, report-only boundary, report identity, old-to-new payload boundary, 239/1 private-label byte comparison, archive identity and equality, audit and adjudication identities, corrected source-pin reconciliation, grader, protocol, world, world-build, registries, and manifests;
- a requirement matrix for sections 1-6;
- exact results for both `held_out_7f53f837` temptation polarity and `held_out_18cbf054` cascade-terminal adjudication, stated only by opaque ID and non-revealing verdict category;
- audit coverage and check totals, test, vet, grader, dry-seal, materialization, exact-byte, opacity, custody, protected-state, and clean-worktree results;
- findings, accepted limitations, custody status, and explicit confirmation that no outcome was opened;
- the smallest next artifact allowed by the verdict.

Verdict bar

- `ACCEPT` requires every changed identity, custody boundary, source pin, corrected-label premise, disputed-adjudication premise, archive, registry, manifest, executable check, protected-state check, and unchanged-byte dependency to pass, with no unresolved P0 or P1. Every P2 must be fixed or specifically justified as nonblocking.
- The retained independent audit's `REVISE` verdict is evidence requiring direct resolution, not an automatic package verdict. `ACCEPT` requires independent reproduction sufficient to close or sustain its opaque disagreement.
- A top-level failure involving commit ancestry, patch bytes, archive bytes, reviewer independence, private-label custody, an observed prospective outcome, or a consumed tranche is P0 and requires `REJECT`.
- Use `REVISE` for a correctable label, registry, manifest, audit-evidence, report, or review-record defect while package identity and custody remain intact.
- Use `REJECT` when the package cannot be repaired without violating identity, validity, independence, or custody premises.

If `ACCEPT`, recommend only that the project chair consider importing the exact reviewed base-to-payload patch. The reviewer does not authorize import. The report-only commit is not part of that patch. The chair must record any import decision separately and verify protected bytes, `status: partial`, `remaining: ["schedule digest"]`, and both false gates afterward. Acceptance or import does not authorize schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.

If `REVISE`, keep the package unimported and name the smallest replacement artifact. If only a public review-record correction is needed, require only that record. If an existing private evidence record needs a new independent follow-up without a label change, preserve payload, patch, labels, archive, registries, and manifests and require the new evidence plus a report-only replacement. Any private-label change requires a new affected label digest and archive identity, regenerated dependent registry and manifest, a new base-parented payload, patch, and report-only child.

If `REJECT`, identify the failed independence, identity, validity, or custody premise and keep both gates closed.
```
