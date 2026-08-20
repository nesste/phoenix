# Protocol-v4 sealed-corpus revision-handoff independent review

- **Reviewer role:** Evaluation reviewer. Independent of the rejected-candidate corpus author, that candidate’s independent review session, Phoenix implementation in this session, and this handoff’s authoring. This session ran in the Phoenix implementation workspace at `D:\Work\personal\phoenix` and therefore had read access to implementation files.
- **Date:** 2026-08-20
- **Verdict:** `REVISE`
- **Reviewed artifact:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md`
- **Original contract:** `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md`
- **Prior corpus review:** `docs/reviews/2026-08-19-protocol-v4-sealed-corpus-review.md` (`REVISE`, candidate `c193c786…` unimported)
- **Revision checkout at review:** `D:\Work\personal\phoenix-evaluator-v4-revision`, branch `codex/evaluator-v4-corpus-revision`, HEAD `c852101e8d7cb52e4569bf3866de54a0ce648b44`, worktree clean
- **New private custody root:** not created yet (`D:\Work\personal\phoenix-evaluator-private-v4-revision` does not exist)

This is a handoff-contract review. It does not authorize corpus authoring to start, import, a validation schedule, Gate 1A, validation, held-out, or any outcome run.

## Verdict

The handoff is the right *kind* of assignment: wholly replace both sealed tranches in a new checkout and new private archive; keep the rejected candidate unopened and unimported; treat it as a fourth collision set rather than a repair target. The pinned freeze identities recompute. Allocation, cascade-recipe, five-variant, family-ID, two-commit patch-identity, and `seq`-path compatibility rules are present in substance.

It is not yet issuable to a fresh non-implementer. Three P1 process holes would recreate the last failure or publish a false identity: the new `seq` control names the wrong check kind; the independent blueprint audit has no role that cannot be the corpus author; and the two-commit public-patch digest has no Windows binary-safe capture rule even though every listed work location is a Windows path.

Do not give this text to a corpus author until those three items are patched. Keep Gate 1A closed.

## Independently recomputed identities

Recomputed in the clean revision checkout at evaluator base `c852101e…`. Freeze commit `4b514712…` is an ancestor. Rejected candidate `c193c786…` is present as a Git object and must not be checked out for authoring.

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Protocol v4 | SHA-256 of LF-normalized UTF-8 of `experiments/frontier-v1/protocol.json` | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `corpusctl digest` of `worlds/dev-repo/world.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json` | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` (both) |
| Production/authoring world, raw bytes | SHA-256 of file bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` (byte-identical) |
| Current grader | `go run ./cmd/corpusctl grader-digest --repo-root ../../..` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest, raw bytes | SHA-256 of `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| Content-addressed `world_build_digest` | Manifest field `digest` (Go `json.Marshal` of the `build` record); not the manifest-file hash | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |
| Rejected public patch | From prior review; not re-hashed in this session | `sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65` |

`pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, `may_open_held_out: false` in both the revision checkout and the implementation workspace.

No protocol, runner, world, grader, Arm B, analysis, freeze, gate, schedule, or full-label byte was modified by this review. `seal --write` was not run. No model, arm, trial, or prospective outcome was observed.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Pinned freeze identities | Pass | All six handoff pins recompute. Gates remain false. |
| 2. Custody and no-outcome framing | Pass | New checkout, new private root, rejected private labels closed, no import, no schedule, no Gate 1A. |
| 3. Prior P1 corpus defects restated as authoring rules | Pass | Noun-swap / global template ban, eight distinct mix mechanisms, no mirrored held-out pair, cascade goals must not cue recipes, five real state/task variants, token/filename/alias alone is not a variant. |
| 4. Prior P2 identity/label defects restated | Fail | Family-ID reuse, lexicographic archive index, two-commit non-circular patch digest, and path/check compatibility are present. The compatibility rule names `act_sequence` `seq`, which the grader does not use (P1). |
| 5. Two-commit public-patch identity | Fail | The payload-then-report split is the correct non-circular design. Capture of `git diff --binary` as raw bytes is unspecified on Windows (P1). |
| 6. Independent audit before sealing | Fail | Blueprint audit is required before labels are finalized, but no role distinct from the corpus author is assigned to it. The later acceptance reviewer is barred from authoring, so they cannot be that auditor (P1). |
| 7. Contract completeness for a non-implementer | Fail | Fresh author is told to read the original handoff and the prior review, not `protocol.json`, Task 0.4 class definitions, or the case/label schemas (P2). |
| 8. Frozen allocation | Pass | 120/24/5, 8/8/8 groups, mix arithmetic, and per-family cascade/adversarial floors match `tranche_design.family_allocation`. |
| 9. Four collision sets plus cross-tranche disjointness | Pass | Authoring, retired-v3 validation, retired-v3 held-out, rejected v4 public tree, and the other new tranche. Inventory fields are expanded past exact strings. |
| 10. Protocol-v4 executable semantics | Partial | Relies on “original requirements remain in force.” Stale-frontier shared events, genuine absence, twelve flat tools, and adversarial-as-untrusted-data are not restated here (P3). |

## Coverage of the 2026-08-19 corpus findings

| Prior finding | Handoff treatment |
| --- | --- |
| P1 mix/held-out noun-swap of one template | Explicit ban on a parameterized global template, mirrored layouts, alias tables, and label skeletons |
| P1 cascade goals state the recipe | Goal may state outcome and constraints; must not enumerate, synonymize, or strongly cue the path; three executable verbs; authoring diagnostic sequence is a collision set |
| P1 five variants are goal/token rewrites | Private five-row matrix; public evidence must use observable state/task axes; three same-class cases need a material state or task difference |
| P2 archive digest used custodian order | Ordinal-byte sort of `labels/…` paths, LF-terminated index, SHA-256 of index bytes |
| P2 shorter acceptable path vs `seq` checks | Intended, but names the wrong check kind (P1 below) |
| P2 family IDs recycle `*_family_001`–`024` | Forbidden; schema still allows other three-digit suffixes |
| P2 report omitted candidate commit and patch digest | Two-commit procedure so the report can name the payload commit and patch digest without self-reference |
| P3 cloned goal typos | Proofread public goals and report before freezing the payload candidate |

## Findings

### [P1] The new `seq` control names a check kind that has no `seq` field — `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md:126`

The previous candidate listed shorter acceptable paths that the grader cannot score as success because path-bound `seq` indexes are resolved against the matched path’s offsets. That binding lives on `act_status` and `act_output_matches` in `label.schema.json` and `grade.go` (`resolveActIndex` maps `check.seq` through `path.offsets` when an `act_sequence` match exists). `act_sequence` itself has `mode`, not `seq`.

A fresh author who implements the “in particular” sentence literally will search for `act_sequence.seq`, find nothing, and skip the check this paragraph exists to add. The preceding sentence (“all path-bound checks can succeed on that path”) is the correct rule. Replace the example with `act_status` / `act_output_matches` `seq` indexes, and require the mechanical test to evaluate every listed path independently, including shorter ones.

### [P1] Independent blueprint audit can be satisfied by the corpus author — `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md:60` and `:92`

The last candidate failed because mix families and the opposite tranche were one generating template. This handoff’s main new control is “an independent private blueprint audit must pass before labels are finalized.” The named roles are corpus author, first-pass labeler, independent label auditor, and private-label custodian. Blueprint audit is earlier than labeling. The only person the handoff forbids from authoring is the later acceptance reviewer, who therefore cannot perform a pre-label audit.

If the corpus author signs their own 48 blueprints, the control is the same mechanical self-audit that missed noun-swaps last time. Require a blueprint auditor who is not the corpus author, and keep that person out of the later acceptance review. If the project cannot staff that role, say so and treat pairwise similarity plus the later review as the only independence, rather than claiming an independent pre-label audit.

### [P1] Public-patch SHA-256 capture is not operable on the stated Windows workstations — `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md:177`

The two-commit scheme is correct: hash the base-to-payload-candidate raw diff, then commit a report that records that digest. The handoff only says to generate `git --no-pager diff --binary …` “as raw bytes.” Review prompt v2 already had to forbid PowerShell pipelines, `>`, `Out-File`, `Set-Content`, and any text decode/re-encode, after identity hashes were corrupted that way. Every path in this handoff is `D:\…`.

Without that capture rule, the report will freeze a digest of a re-encoded file. The next independent review, hashing native git stdout, will not match. Copy the v2 binary-safe capture paragraph into this authoring handoff, including byte-equality against the stored patch file.

### [P2] A non-implementer is not required to read the class and schema contracts

The assignment is for someone who did not implement Phoenix. Required reading is only the original evaluator handoff and the prior review. Neither document states the Task 0.4 class one-liners for `direct`, `recovery`, `far_discovery`, and `temptation`. `protocol.json` `tranche_design.family_allocation` is the allocation source of truth; `schema/case.schema.json`, `schema/label.schema.json`, and `schema/common.schema.json` (`familyId` = `*_family_[0-9]{3}`) are the emit constraints.

Add those paths to the required-read list. Point family IDs at the schema pattern and the free suffix range `025`–`999` after excluding collision-set IDs.

### [P2] `world_build_digest` must be recomputed, but the algorithm is not attached

The handoff correctly lists both the manifest-file hash `da59f795…` and the content-addressed digest `27c2f537…`. “Independently recompute” does not say that `27c2f537…` is SHA-256 of Go `json.Marshal` of the manifest `build` record, and is not the `digest` field copied out of the file. Reading the embedded field would vacuously succeed. Attach the algorithm the way v2 did for the reviewer.

Canonical world identity should also name `corpusctl digest` on the two world files. Grader identity already names `grader-digest`.

### [P2] Authoring remains a likely labeling tutorial

Collision set 1 includes visible authoring families, cases, fixtures, goals, and labels. The v3 evaluator handoff said those files may exist in the checkout so `seal` can detect reuse, and must not be used as examples. This revision drops that sentence. A non-implementer learning the label schema from `experiments/frontier-v1/labels/authoring/` will copy authoring path skeletons, which is the same class of template reuse the prior review rejected across tranches.

Restore the explicit rule: authoring in the revision checkout is a collision set, not a generator or a labeling tutorial. Label schema and grader pin come from `schema/` and `corpusctl`, not from authoring labels.

### [P3] Revision HEAD is not pinned as a precondition

The listed revision checkout is presently clean at evaluator base `c852101e…`. The handoff never requires verifying that before authoring. A dirty tree or a start from `codex/evaluator-v4-corpus` (`c193c786…`, still checked out in a sibling worktree) would mix rejected public files into the replacement. Add: confirm `HEAD == c852101e…`, `git status --short` empty, and do not use the rejected-candidate worktree.

### [P3] Protocol-v4 case semantics stay only in the original handoff

Stale-frontier via a declared shared state event, genuinely absent capabilities, solvability through the twelve frozen flat tools, adversarial fixture text as untrusted data, and “do not modify protocol/runtime/world/grader/freeze bytes” are still required, but only by the “original requirements remain in force” sentence. Restate them in this document so a fresh author cannot treat the revision text as complete.

## Accepted limitations

- This review is of the authoring handoff, not of a replacement corpus. No new labels, blueprints, or public patch exist to inspect.
- Reviewer session had read access to the Phoenix implementation checkout. No implementation file was modified except this review record.
- Rejected-candidate private labels were not opened. Rejected public material was not used as a writing template; public review findings were used only as the defect list this handoff must close.
- `familyId` disjointness was judged against the schema pattern and the four collision sets, not against live-path reuse of `corpus/validation` and `corpus/held_out`, which remains expected replacement behavior.

## Custody status

Intact. The rejected candidate remains unimported. The historical protocol-v3 private archive was not inspected. The rejected v4 private archive was not inspected. No Phoenix arm output, runner result, retained authoring outcome, or prospective sealed outcome was used. Gate 1A was not opened.

The revision checkout is the right starting tree. The new private custody directory does not exist yet and must not be created until this handoff is issuable.

## Smallest next artifact

A patched revision handoff, still outcome-free, that:

1. Names `act_status` / `act_output_matches` `seq` indexes as the path-bound checks, and requires every listed acceptable path to be able to satisfy those indexes.
2. Assigns blueprint audit to a person who is not the corpus author and not the later acceptance reviewer, or drops the claim of independent pre-label audit.
3. Copies the Windows binary-safe `git diff --binary` capture rule from review prompt v2 into the two-commit identity procedure.
4. Adds required reading of Task 0.4 class definitions, `protocol.json` `tranche_design.family_allocation`, and the case/label/common schemas, with the `familyId` suffix rule.
5. Attaches the `world_build_digest` algorithm and the authoring-is-collision-not-tutorial sentence.

Do not issue the current text to a corpus author. Do not import `c193c786…`. Do not generate a schedule. Do not open Gate 1A.

A matching independent-review prompt for the two-commit, payload-only patch, and ordinal-byte archive index is a separate chair artifact. It is not required inside this authoring handoff, but it must exist before the replacement candidate is reviewed. Review prompt v2 still describes the rejected one-commit package and must not be reused as-is.
