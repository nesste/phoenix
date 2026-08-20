# Protocol-v4 sealed-corpus revision-handoff independent review 2

- **Reviewer role:** Evaluation reviewer. Independent of the rejected-candidate corpus author, that candidate’s independent review session, Phoenix implementation in this session, the original revision-handoff authoring, and the first handoff-review session. This session ran in the Phoenix implementation workspace at `D:\Work\personal\phoenix` and therefore had read access to implementation files.
- **Date:** 2026-08-20
- **Verdict:** `ACCEPT`
- **Reviewed artifact:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md` at `c49d4812f5f7d3b48a11a0498ccd02aa3c3b08cc` (`docs: close corpus revision handoff process gaps`)
- **Reviewed handoff, LF-normalized UTF-8:** `sha256:8b0e1a919e9afecdd054070f1c3361bc940deceea690c679abc896f10bd52417` (raw bytes identical; LF-only)
- **First handoff review:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff-review.md` (`REVISE`)
- **Original contract:** `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md`
- **Prior corpus review:** `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review.md` (`REVISE`, candidate `c193c786…` unimported)
- **Revision checkout at review:** `D:\Work\personal\phoenix-evaluator-v4-revision`, branch `codex/evaluator-v4-corpus-revision`, HEAD `c852101e8d7cb52e4569bf3866de54a0ce648b44`, worktree clean before and after LF materialization
- **New private custody root:** not created (`D:\Work\personal\phoenix-evaluator-private-v4-revision` does not exist)

This is a handoff-contract review of the patched authoring assignment. It authorizes issuing that assignment to a fresh corpus author. It does not authorize import, a validation schedule, Gate 1A, validation, held-out, or any outcome run.

## Verdict

The three first-review P1 process holes are closed in the current text: path-bound `seq` now names `act_status` and `act_output_matches`; blueprint audit is a named role that cannot be the corpus author or the later acceptance reviewer, with a hard stop if that person is missing; and the two-commit public-patch capture forbids PowerShell text re-encode and requires a second native-stdout byte-equality check. The first-review P2 and P3 items are also present.

Pinned freeze identities recompute in the clean revision checkout. Allocation, cascade-recipe, five-variant, family-ID, four-collision-set, ordinal-byte archive-index, and no-outcome rules remain in force. No P0, P1, or P2 defect remains in the authoring handoff.

Issue this text to a fresh corpus author. Keep Gate 1A closed. Do not import `c193c786…`.

## Independently recomputed identities

Recomputed in the clean revision checkout at evaluator base `c852101e…` after `core.autocrlf=false`, `core.eol=lf`, and `git checkout-index --all --force`. Both `git status --short` checks were empty. Freeze commit `4b514712…` is an ancestor. Rejected candidate `c193c786…` is present as a Git object and was not checked out.

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Protocol v4 | SHA-256 of LF-normalized UTF-8 of `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `corpusctl digest` of `worlds/dev-repo/world.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json` | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` (both) |
| Production/authoring world, raw bytes | SHA-256 of file bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` (byte-identical) |
| Current grader | `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest, raw bytes | SHA-256 of `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| Content-addressed `world_build_digest` | SHA-256 of Go `json.Marshal` of the typed `build` record; independently equal to the embedded digest; not the manifest-file hash | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Rejected public patch | From prior corpus review; not re-hashed in this session | `sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65` |

`pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, `may_open_held_out: false` in both the revision checkout and the implementation workspace.

No protocol, runner, world, grader, Arm B, analysis, freeze, gate, schedule, or full-label byte was modified by this review. `seal --write` was not run. No model, arm, trial, or prospective outcome was observed. The revision worktree remained clean.

## Closed first-review findings

| First-review finding | Status | Evidence |
| --- | --- | --- |
| P1 `seq` control named `act_sequence.seq` | Closed | Handoff labels section: when an `act_sequence` check matches a path, resolve `act_status.seq` and `act_output_matches.seq` through that path’s offsets as `grade.go` `resolveActIndex` does; every listed path, including shorter ones, must have those positions. `label.schema.json` puts `seq` on `act_status` and `act_output_matches`; `act_sequence` has `mode`. |
| P1 blueprint audit could be the corpus author | Closed | Named role `blueprint auditor` must not be the corpus author and must not perform later acceptance review. Author may not approve their own 48 blueprints. Missing auditor is a stop, not a self-audit or similarity-screen substitute. |
| P1 Windows public-patch capture unspecified | Closed | Native `git --no-pager diff --binary` stdout, binary-safe capture to the revision patch path, forbidden PowerShell `>` / `Out-File` / `Set-Content` / pipeline re-encode, SHA-256 of stored bytes, second native capture, byte equality. |
| P2 non-implementer not required to read class/schema contracts | Closed | Required reading now includes Task 0.4 class definitions, `protocol.json` `tranche_design.family_allocation`, case/label/common schemas, `grade.go`, and both world files. Family IDs: schema pattern plus free suffixes `025`–`999` after excluding collision-set IDs. |
| P2 `world_build_digest` algorithm missing | Closed | Recompute with `corpusctl digest` on both worlds, `grader-digest` for the grader, and Go `json.Marshal` of the manifest `build` record only. Copying the embedded digest is forbidden. |
| P2 authoring usable as a labeling tutorial | Closed | Visible authoring material is a collision set, not a generator or tutorial. Inspect authoring only for collision or leakage. Learn labels from `schema/` and `corpusctl`, not `labels/authoring/`. |
| P3 revision HEAD not pinned | Closed | Before authoring: `HEAD == c852101e…`, empty status, do not use `D:\Work\personal\phoenix-evaluator-v4`. |
| P3 protocol-v4 case semantics only by incorporation | Closed | Restated: strict case schema, canonical world digest, twelve flat-tool solvability, shared executable-act `state_changes`, declared-event `stale_frontier`, genuine `absence`, adversarial text as untrusted data, no protocol/runtime/world/grader/freeze/gate edits. |

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Pinned freeze identities | Pass | All six handoff pins recompute. Gates remain false. |
| 2. Custody and no-outcome framing | Pass | New checkout, new private root, rejected private labels closed, no import, no schedule, no Gate 1A. |
| 3. Prior P1 corpus defects restated as authoring rules | Pass | Noun-swap / global template ban, eight distinct mix mechanisms, no mirrored held-out pair, cascade goals must not cue recipes, five real state/task variants. |
| 4. Prior P2 identity/label defects restated | Pass | Family-ID reuse ban with schema suffix rule, ordinal-byte archive index, two-commit non-circular patch digest, path-bound `act_status` / `act_output_matches` `seq` compatibility. |
| 5. Two-commit public-patch identity | Pass | Payload candidate then report child; reviewed import payload is the base-to-payload diff; report excluded from that patch; Windows binary-safe capture and byte-equality required. |
| 6. Independent audit before sealing | Pass | Named blueprint auditor, distinct from author and later reviewer, with staffing stop. Later acceptance reviewer isolated from every authoring and audit role. |
| 7. Contract completeness for a non-implementer | Pass | Required reading covers original contract, prior corpus review, Task 0.4, protocol allocation, schemas, grader, and worlds. |
| 8. Frozen allocation | Pass | 120/24/5, 8/8/8 groups, mix arithmetic, and per-family cascade/adversarial floors match `tranche_design.family_allocation`. Direct-group and recovery-group per-family layouts match. Class totals 24/24/12. |
| 9. Four collision sets plus cross-tranche disjointness | Pass | Authoring, retired-v3 validation, retired-v3 held-out, rejected v4 public tree, and the other new tranche. Inventory fields include shapes and patterns, not only exact strings. |
| 10. Protocol-v4 executable semantics | Pass | Restated in this handoff; twelve frozen Arm A tools remain the flat-tool surface. |

## Coverage of the 2026-08-19 corpus findings

Unchanged from the first handoff review except the former `seq` misname, which now matches `grade.go`. The patched handoff still forbids the mix/held-out noun-swap, recipe-stating cascade goals, goal-only variants, custodian-order archive digest, recycled `*_family_001`–`024` IDs, and a report that cannot name the payload commit and patch digest.

## Findings

None.

## Accepted limitations

- This review is of the patched authoring handoff, not of a replacement corpus. No new labels, blueprints, or public patch exist to inspect.
- Reviewer session had read access to the Phoenix implementation checkout and to the first handoff-review record. No implementation file was modified except this review record. The first review was used only as the defect list the patch must close.
- Rejected-candidate private labels were not opened. Rejected public material was not used as a writing template.
- `familyId` disjointness was judged against the schema pattern and the four collision sets, not against live-path reuse of `corpus/validation` and `corpus/held_out`, which remains expected replacement behavior.
- The assignment cannot be completed by one person: the blueprint auditor must be someone other than the corpus author, and the later acceptance reviewer must be someone other than every authoring and audit role. If those people are not available, the handoff requires a stop rather than self-audit.
- `fixture.schema.json` is not on the required-read list. Seal will still reject schema-invalid fixtures. Distinctness remains a generating-idea rule, not a unique `template` string.
- A matching independent-review prompt for the two-commit, payload-only patch and ordinal-byte archive index still does not exist. The handoff correctly places that on the chair before the replacement candidate is reviewed, not before this assignment is issued. Review prompt v2 must not be reused as-is.

## Custody status

Intact. The rejected candidate remains unimported. The historical protocol-v3 private archive was not inspected. The rejected v4 private archive was not inspected. No Phoenix arm output, runner result, retained authoring outcome, or prospective sealed outcome was used. Gate 1A was not opened.

The revision checkout is the right starting tree. The new private custody directory does not exist yet and may be created by the corpus author after this assignment is issued.

## Exact next allowed action

Issue `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md` at `c49d481…` to a fresh corpus author who did not create the rejected protocol-v4 candidate, perform its independent review, or implement Phoenix.

That author may create `D:\Work\personal\phoenix-evaluator-private-v4-revision` and work only in `D:\Work\personal\phoenix-evaluator-v4-revision` at evaluator base `c852101e…`. They may not use `D:\Work\personal\phoenix-evaluator-v4` or `D:\Work\personal\phoenix-evaluator-private-v4`.

The chair must prepare a matching review assignment for the payload-only patch, separate report commit, new collision set, and ordinal-byte archive index before the replacement candidate is reviewed.

Do not import `c193c786…`. Do not generate a schedule. Do not open Gate 1A.
