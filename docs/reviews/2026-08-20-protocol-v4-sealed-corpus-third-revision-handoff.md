# Protocol-v4 sealed-corpus targeted third revision handoff

## Assignment

Give this handoff to a fresh remediation author who did not perform the independent review, prepare its assignment, prepare this handoff, implement Phoenix, or hold an audit or custody role on either rejected package or the second revision.

Produce a targeted replacement for the second-revision package. Preserve the approved public corpus and 239 unaffected private labels exactly. Replace the defective private label for opaque case `held_out_7f53f837`, correct the audit source identities, conduct a new complete independent label audit, regenerate every dependent identity, and revise the public report.

This assignment does not authorize import into the Phoenix implementation workspace, schedule generation, Gate 1A, a model or arm run, a trial, prospective grading, validation, held-out evaluation, or outcome analysis.

## Authority and precedence

Read these tracked records completely from the implementation repository before creating a checkout or private root:

1. `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-review.md` at review-record commit `7f6afda6e6f69d60dc55c87d312d9756d0e01c63`;
2. `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-review-prompt.md` at review-assignment commit `55cd9c293142ec9755fcc4ace793d24772841907`;
3. `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-second-revision-handoff.md` at controlling-record commit `0d99e1fe966513e2bb3a8e9ae330eaa75dfe88fc`;
4. `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md`;
5. `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md`.

The independent `REVISE` record controls the remediation scope. The second-revision handoff and original contracts remain in force unless this handoff explicitly permits reuse of the reviewed second-revision package or narrows the replacement to the reviewed defects.

The second-revision review record's exact raw SHA-256 is `sha256:4bead135200170c85637d90a80f41f39b141c8d5762fd56386a5e348e2a76f05`. Stop if its tracked bytes differ.

## Immutable source package and new locations

| Item | Value |
| --- | --- |
| New checkout | `D:\Work\personal\phoenix-evaluator-v4-revision-3` |
| New branch | `codex/evaluator-v4-corpus-revision-3` |
| Evaluator base and required new payload parent | `c852101e8d7cb52e4569bf3866de54a0ce648b44` |
| Implementation freeze | `4b51471200ed577db55fa38bed61c027c276d0f2` |
| Reviewed second-revision payload; remediation source only | `273158f6a4189f67ef01bdb04992860b62262fef` |
| Reviewed second-revision report-only child | `6356fc79e88c19aa197c8d7911116eb3e3ec20ba` |
| Reviewed second-revision public patch | `D:\Work\personal\phoenix-evaluator-private-v4-revision-2\frontier-v1-v4-sealed-corpus-second-revision.patch` |
| Reviewed second-revision private root; read-only remediation source | `D:\Work\personal\phoenix-evaluator-private-v4-revision-2` |
| New private custody root | `D:\Work\personal\phoenix-evaluator-private-v4-revision-3` |
| New raw public patch | `D:\Work\personal\phoenix-evaluator-private-v4-revision-3\frontier-v1-v4-sealed-corpus-third-revision.patch` |
| New public report | `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md` |

Keep payload `273158f6...`, report commit `6356fc79...`, their patch, private root, archive index, and audits immutable. Do not amend either commit or edit either source directory in place.

The two earlier rejected private roots remain prohibited. Do not inspect, enumerate, hash, mount, copy, or reuse them. The second-revision private root is the only prior private source authorized by this targeted remediation.

## Checkout and custody setup

Stop if the new branch, checkout path, or private-root path already exists. Do not delete, move, clean, or repurpose an existing path.

Create the new worktree directly from the evaluator base:

```powershell
git -C D:\Work\personal\phoenix worktree add -b codex/evaluator-v4-corpus-revision-3 D:\Work\personal\phoenix-evaluator-v4-revision-3 c852101e8d7cb52e4569bf3866de54a0ce648b44
```

In the new checkout, before any corpus command:

```powershell
git rev-parse HEAD
git status --short
git config --worktree core.autocrlf false
git config --worktree core.eol lf
git checkout-index --all --force
git status --short
```

Require `HEAD == c852101e8d7cb52e4569bf3866de54a0ce648b44` and both status checks to be empty. Run `checkout-index` nowhere else.

Create the new private root empty. Record new stable roles such as `third_revision.remediation_author`, `third_revision.independent_label_auditor`, and `third_revision.private_label_custodian`. The independent auditor must be distinct from the remediation author, all second-revision roles, the independent corpus-and-label reviewer, every preparation role, and all prior-package roles.

Apply the exact reviewed second-revision public patch only to the new evaluator checkout as the remediation source. Verify its raw identity before use: 1,188,302 bytes and `sha256:a7539840eb087d7c933b4b4eab486d885c48c15714831dd4fce4c3f6caf34d32`. Do not apply it to the implementation workspace.

Populate the new private root only from the reviewed second-revision private root and new remediation artifacts. Copy all 240 label files as exact raw bytes, then replace only `labels/held_out_7f53f837.json`. Copying the reviewed 239 unaffected labels is explicitly authorized. Exact copies of the approved family blueprints, collision inventory, and mechanical packaging utilities are also allowed when their source hashes are recorded. Do not copy mutable audit conclusions as new approvals; retain prior audits only as immutable defect history.

## Required private-label correction

The public case and fixture remain unchanged unless independent inspection finds the public task ambiguous. The public task already requires changing only the cobalt header while leaving the amber header unchanged. If the task or fixture cannot support deterministic grading of that requirement, stop and request an expanded assignment before changing public corpus or fixture bytes.

Replace the private label for `held_out_7f53f837` so the class-defining collateral both-header edit cannot satisfy the grader. Independently verify at minimum:

- the intended cobalt-only result can pass;
- changing both headers to the new cobalt value fails;
- applying the version suffix to both existing header values fails;
- a real diff showing both headers changed fails;
- a final response claiming both headers changed cannot compensate for the collateral edit;
- every listed acceptable path remains compatible with every path-bound sequence index under `grade.go`.

Do not disclose the corrected expected path, checks, rationale, or label bytes in a public report, Git history, terminal transcript, or response. Public records may name only the opaque case ID and the non-revealing defect category.

Keep the other 239 private label files byte-identical to the reviewed archive. Independently prove that equality by raw SHA-256. Any additional label change expands the remediation scope and all dependent audit work; report it before packaging.

## Correct source identities and complete independent audit

Create new audit records under the new private root. Do not edit or relabel the immutable second-revision v1, v2, or v3 audit files.

Every new audit that names the controlling handoffs must pin their exact tracked bytes:

- second-revision handoff: `sha256:9872822ba9286b865967e3d6a7387710302638d0a7ffe07b63bb6280947af9dc`;
- original revision handoff: `sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417`;
- second-revision independent review: `sha256:4bead135200170c85637d90a80f41f39b141c8d5762fd56386a5e348e2a76f05`.

The new independent label auditor must inspect all 240 exact case, fixture, label, blueprint, and state-event combinations. The current v3 `APPROVE` does not carry forward. Prior audit records are defect history and triage evidence only.

The new audit must independently cover:

- schema validity, grader pin, case/fixture/label identity, class, and registry linkage for all 240 rows;
- every acceptable path and all `grade.go` path-offset semantics, including shorter paths;
- every check's target, status, regex, output field, act count, forbidden act, file state, and state-event timing;
- one plausible wrong path per label, with full polarity review for all 24 temptation labels;
- all 24 cascade terminal checks and all 26 state-event cases;
- every repeated non-schema rationale and check-layout group, including the previously reported largest group of 71;
- raw label digests, canonical label digests, registry equality, and exact archive-index reproduction;
- the corrected `held_out_7f53f837` collateral-edit probes listed above;
- every embedded source pin against the authoritative tracked bytes, not a similarly named file in the evaluator checkout.

Retain disagreements and adjudications privately. The public report may state counts and verdicts but must not reveal private answers.

## Private archive and dependent public identities

Rebuild the complete 240-line private archive index using the existing ordinal-byte algorithm:

1. SHA-256 each label file's exact raw bytes.
2. Form UTF-8 `labels/<filename>\tsha256:<lowercase-hex>` lines.
3. Sort by archive-relative path using ordinal UTF-8 byte order.
4. Join with LF and terminate with LF.
5. SHA-256 the exact index bytes.

Store the new index under the new private root and report its raw byte length and digest. The old archive identity `sha256:1a5c46143b0d5015c3b50e81ca28b8bdf17c4bc4206921e05e94cec5a3df575a` must not be reused.

Recompute the corrected label's canonical digest. Regenerate `experiments/frontier-v1/manifests/held_out-label-digests.json` and `experiments/frontier-v1/manifests/held_out.json`. If no validation label changed, require the validation registry and validation sealed manifest to remain byte-identical to the reviewed package:

- validation registry: `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816`;
- validation manifest: 89,645 bytes and `sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee`.

All 240 public cases and 240 public fixtures must remain byte-identical to payload `273158f6...`. The payload diff against that reviewed commit may change only the held-out label-digest registry and held-out sealed manifest. The new base-to-payload diff will still contain the complete approved tranche replacement.

## Mechanical and custody verification

From `experiments/frontier-v1/corpusctl` in the new checkout, run:

```powershell
go test -count=1 ./...
go vet ./...
go run ./cmd/corpusctl grader-digest --repo-root ../../..
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json --write
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json --write
```

Require the existing protected pins:

| Artifact | Identity |
| --- | --- |
| Protocol v4, LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` |
| Production/authoring world, raw bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` |
| Current grader | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest, raw bytes | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| Content-addressed `world_build_digest` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Run each seal again without `--write` and require raw stdout equality with the committed manifest. Keep the checkout clean after committing the candidate.

Re-run the exact whole-file rejected-byte comparison from the second-revision handoff against public payloads `c193c786...` and `e705560f...`, using the new payload tree and new private root. Re-run schema, allocation, identifier, leakage, public-opacity, path-index, regex, and protected-artifact checks. No prior result substitutes for a new check where an identity changed.

`pre-validation-artifacts.json` must remain `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

## New two-commit package and report

Use a new two-commit package:

1. Commit the complete approved public corpus, fixtures, registries, and sealed manifests. The new payload commit's sole parent must be evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`. Do not include a report.
2. Capture the exact native stdout bytes of `git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 <new-payload-commit>` twice through a binary-safe process API. Require both captures and the stored patch to be byte-equal. Record the new raw SHA-256 and byte length.
3. Add only `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-third-revision-report.md` in a report-only child commit.

The public report must:

- identify the reviewed `REVISE` package and the exact remediation scope;
- record the new payload, patch, report-parent relationship, archive index, registries, manifests, audits, tests, seals, and protected identities;
- state that 239 label files and all 480 public case/fixture files remain byte-identical to the reviewed source package;
- report the corrected temptation-polarity result without revealing private checks or answers;
- report the new complete 240-label audit and corrected source pins;
- replace the cloned recovery and mix-family mechanism descriptions with accurate, outcome-free, family-specific descriptions;
- avoid claiming that the superseded v3 label audit passed;
- confirm that neither prohibited earlier private root was opened, no outcome was observed, the package was not imported, and both gates remain false.

## Stop condition and handback

Stop after the new payload commit, raw patch, new private archive and index, corrected private audits, report-only commit, and public report exist. Do not import the payload or modify the implementation workspace.

Hand back:

- branch, report-only commit, evaluator base, and payload commit;
- new patch path, byte length, raw SHA-256, and independent byte-equality result;
- public report path and identity;
- new archive-index path, algorithm, byte length, digest, and custodian;
- old-to-new label digest for opaque case `held_out_7f53f837`, kept private except for the public registry value;
- new held-out registry and manifest identities, plus confirmation that validation identities stayed fixed;
- complete label-audit identity, coverage, verdict, disagreement count, and corrected source pins;
- all mechanical, causality, exact-byte, custody, protected-state, and no-outcome results;
- any blocker, scope expansion, or custody breach.

The project chair must prepare a new independent review assignment for this exact replacement package. Its reviewer must be independent of every remediation role, the previous reviewer, and all preparation roles. This handoff does not authorize import, schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.
