# Protocol-v4 sealed-corpus second replacement revision handoff

## Assignment

Give this handoff to a fresh corpus author who did not create or audit either rejected protocol-v4 package, prepare either package or review assignment, perform either independent review, or implement Phoenix.

Create a wholly new `validation` and `held_out` corpus for `frontier-v1`. This is a clean-room replacement authored from the frozen protocol, schemas, grader, and world contract. It is not a repair of payload `e705560fef6012f84ec0199b87ce790a62c17d0e`.

Keep both rejected packages unimported. Treat their public material only as collision evidence. Do not reuse any script, fixture, generator, blueprint, label, rationale, acceptable-path design, check layout, or label skeleton from either package.

This assignment produces outcome-free public inputs and privately held labels. It does not authorize import, schedule generation, Gate 1A, a model or arm run, a trial, validation, held-out evaluation, prospective grading, or outcome analysis.

## Authority and precedence

Read these authoritative records completely from the implementation workspace before creating the checkout or authoring any artifact:

1. the original revision handoff, `D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-revision-handoff.md`;
2. the external `REVISE` review, `D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-revision-review.md`;
3. the original evaluator contract, `D:\Work\personal\phoenix\docs\reviews\2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md`;
4. the first rejected-package review, `D:\Work\personal\phoenix\docs\reviews\2026-08-19-protocol-v4-sealed-corpus-review.md`.

The original revision handoff remains in force in full unless this handoff replaces a location, role name, package identity, collision-set count, filename, or requirement with a stricter rule. This handoff supersedes the original revision handoff for those fields.

The external review is defect evidence and a collision source. It is not a family-design source, fixture library, blueprint, label tutorial, or permission to edit the reviewed package in place.

The prior review prompt at `D:\Work\personal\phoenix\docs\reviews\2026-08-20-protocol-v4-sealed-corpus-revision-review-prompt.md` is retained for provenance. Its instructions governed the completed review and do not instruct the second-revision author.

Do not issue this handoff while it, the external review, or the prior review prompt is untracked. Commit all three records in the implementation repository first. The assignment issuer must record the resulting controlling-record commit in the assignment message. The fresh author must verify that all three paths are present at that commit and read the tracked copies from `D:\Work\personal\phoenix` before creating the base worktree. The new worktree intentionally starts from `c852101e…` and therefore need not contain these later process records.

## New work locations and pinned history

| Item | Value |
| --- | --- |
| Second-revision checkout | `D:\Work\personal\phoenix-evaluator-v4-revision-2` |
| Second-revision branch | `codex/evaluator-v4-corpus-revision-2` |
| Evaluator base and required payload parent | `c852101e8d7cb52e4569bf3866de54a0ce648b44` |
| Implementation freeze | `4b51471200ed577db55fa38bed61c027c276d0f2` |
| First rejected public candidate | `c193c786cc5a65ce6ae97efbd336b2a48f492898` |
| Rejected revision payload; collision use only | `e705560fef6012f84ec0199b87ce790a62c17d0e` |
| Rejected revision report-only child | `7a57479bc64cdb882bf17bfff7ed5d80c2f64760` |
| New private custody root | `D:\Work\personal\phoenix-evaluator-private-v4-revision-2` |
| New raw public patch | `D:\Work\personal\phoenix-evaluator-private-v4-revision-2\frontier-v1-v4-sealed-corpus-second-revision.patch` |
| New public report | `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md` |

The following locations are prior-package evidence and must not be used for authoring or custody:

- `D:\Work\personal\phoenix-evaluator-v4`;
- `D:\Work\personal\phoenix-evaluator-v4-revision`;
- `D:\Work\personal\phoenix-evaluator-v4-review`;
- `D:\Work\personal\phoenix-evaluator-v4-revision-review`;
- `D:\Work\personal\phoenix-evaluator-private-v4`;
- `D:\Work\personal\phoenix-evaluator-private-v4-revision`.

Do not switch an existing worktree to the new branch. Create a new worktree directly from the evaluator base:

```powershell
git -C D:\Work\personal\phoenix worktree add -b codex/evaluator-v4-corpus-revision-2 D:\Work\personal\phoenix-evaluator-v4-revision-2 c852101e8d7cb52e4569bf3866de54a0ce648b44
```

Stop if the branch or path already exists unexpectedly. Do not delete, move, repurpose, or clean a pre-existing directory to make the command succeed.

Before any corpus command, run only in the new checkout:

```powershell
git rev-parse HEAD
git status --short
git config --worktree core.autocrlf false
git config --worktree core.eol lf
git checkout-index --all --force
git status --short
```

Require `HEAD == c852101e8d7cb52e4569bf3866de54a0ce648b44` and both status checks to be empty. Never run `checkout-index` in an implementation, prior-package, review, or historical evaluator checkout.

Create the new private custody root only after the assignment is issued to the fresh author. It must begin empty. Do not seed it by copying any prior private directory or archive.

## Clean-room custody and non-reuse boundary

Assign new stable role identifiers, such as `second_revision.corpus_author`, `second_revision.first_pass_labeler`, `second_revision.independent_blueprint_auditor`, `second_revision.independent_label_auditor`, and `second_revision.private_label_custodian`. Record the actual identifiers privately. Each required independent role must be distinct from the corresponding author and from all prior-package authoring, audit, review, preparation, and Phoenix implementation roles.

Do not inspect, hash, mount, open, execute, source, import, copy, or adapt private material under either prior custody root. In particular, do not reuse a prior package's:

- generation, audit, fixture, hashing, or packaging scripts;
- fixture files, fixture fragments, file maps, directory blueprints, or topology tables;
- family blueprints, variant matrices, domain lists, naming tables, aliases, or action-graph plans;
- full labels, label fragments, rationales, acceptable paths, check IDs, check tuples, regexes, or label skeletons;
- template notes, prompts, adjudications, disagreements, private audits, collision inventories, or archive indexes.

Do not transcribe those artifacts by hand, translate them into a different format, or use them as scaffolding. A renamed noun, token, path, test, service, key, family number, or prose variant is reuse, not a new generating idea.

Read both public review records completely as required above. Their content may be used only as outcome-free defect, collision, provenance, and pinned-identity evidence; it may not supply a design, fixture, generator, blueprint, goal, action graph, label, or rationale. Public payload material at `c193c786…` and `e705560…` may be inspected only through read-only Git-object access for collision detection. Do not check out either payload in the new authoring worktree, execute code taken from either payload, or copy public files into the new checkout as a starting point.

Learn case and label structure from the current schemas and grader. Learn the task classes and allocation from Task 0.4 and `protocol.json`. Learn available behavior from the frozen world contract and its twelve flat-tool equivalents. Visible authoring labels and prior private labels are not tutorials.

Keep every new full label, rationale, blueprint, variant matrix, template note, disagreement, adjudication, private audit, collision inventory, and archive index under the new private custody root and outside every Phoenix checkout, public patch, report, terminal transcript, and response.

## Five prior public collision sets

The new candidates must be structurally disjoint from:

1. visible authoring families, cases, fixtures, goals, and labels;
2. retired protocol-v3 validation public material at evaluator base `c852101e…`;
3. retired protocol-v3 held-out public material at evaluator base `c852101e…`;
4. all public material in rejected candidate `c193c786…`;
5. all public payload material in rejected revision `e705560…`.

The new validation and held-out tranches must also be structurally disjoint from each other.

Before authoring, create a new outcome-free collision inventory under the new private root. Generate it independently; do not extend, edit, or copy either prior inventory. Use read-only Git-object access for the two rejected payloads and base-tree public material. Never open a historical or rejected private archive to populate it.

The inventory must cover public paths, family IDs, case IDs, goals and goal grammar, raw fixture digests, file-map digests, fixture topology and observable shape, service aliases, test catalogs and naming patterns, adversarial-text patterns, state-change patterns, action-graph shapes, and public registry structure for all five prior sets. It must also support cross-tranche comparison for the new package.

At minimum, the `e705560…` collision record must make the following reviewed structures unavailable for reuse:

- the single repository-snapshot fixture grammar used across its 240 fixtures;
- its repeated Go invariant module, empty rename catalog, and adversarial-request paragraph;
- validation-to-held-out `+200` family pairing and slot-for-slot mix layouts;
- reconcile-from-authority cascade goal grammar;
- cascade graphs padded with compile or test steps that do not consume the preceding edit;
- its repeated full-label rationale and check skeletons, without opening or reproducing those private labels.

The external review supplies the last item's non-revealing defect category. It does not authorize access to the private label archive.

New family IDs and case IDs must be absent from every prior collision set and unique across the new tranches. The IDs used by `e705560…`, including `validation_family_101` through `validation_family_124` and `held_out_family_301` through `held_out_family_324`, are unavailable. Preserve the `common.schema.json` family-ID pattern and select otherwise-free suffixes from `025` through `999` only after checking the inventory.

## Replacement design requirements

All allocation, schema, world, grader, protocol-v4 semantics, sealing, private-identity, public-opacity, acceptable-path, two-commit, protected-artifact, and no-outcome rules in the original revision handoff still apply.

Privately design and independently audit all 48 family blueprints before emitting cases or finalizing labels. Each family must have its own task mechanism and fixture topology. Five variants must differ in observable state or task, not only in names, tokens, paths, or wording.

The second replacement must satisfy these additional requirements derived from the external review:

- No global generator may instantiate all families. Shared schema-compliant serialization utilities are allowed only when they do not impose shared fixture grammar, topology, goal grammar, action graph, or label design.
- The eight mix families in each tranche must use eight genuinely different mechanisms and fixture structures.
- Validation and held-out mix-family class layouts must be designed independently. Do not create numbered pairs, `+N` mappings, slot-for-slot mirrors, noun swaps, or topology twins.
- Public cascade goals may state the desired outcome and constraints, but must not enumerate, synonymize, or strongly cue the required action recipe.
- Every cascade must require at least three executable verbs whose order emerges from intermediate evidence. Each required terminal build, test, status, diff, or confirmation step must consume or observe the result of an earlier required act.
- For every cascade, materialize the unedited fixture and prove that its terminal verification cannot establish success before the required earlier change. Then prove that the same verification establishes success after the required change. Retain this outcome-free causality evidence privately for all 24 cascades.
- Do not add optional, already-green, or causally irrelevant acts to make an action graph appear unique.
- New rationales and checks must be case-specific. Repeated frozen-world opening sentences, repeated non-absence check tuples, and shared acceptable-path skeletons require explicit auditor scrutiny and redesign unless the schema itself forces them.
- The public report must describe the mechanisms and topology actually present. Unique hashes, filenames, graphs, or goal strings are not evidence of family independence.

Treat the authoring diagnostic sequence `repo.status -> tests.run -> tests.list -> tests.focus` as a collision set, not a template. If a genuinely distinct cascade uses `tests.run` followed by `tests.focus`, it must include the intervening `tests.list` discovery step.

The independent blueprint auditor must inspect all 48 mechanisms, all 240 variant rows, all 1,128 unordered family pairs, all 576 validation-to-held-out pairs, every mix layout, and every cascade causality record. Automated screening is triage evidence, not an independent verdict.

The independent label auditor must inspect all 240 exact case, fixture, label, and state-event combinations. Prior audit approvals do not carry forward and may not be cited as support for this package.

## Frozen identities and protected state

Independently recompute every identity required by the original revision handoff. The expected pins remain:

| Artifact | Identity |
| --- | --- |
| Protocol v4, LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` |
| Production/authoring world, raw bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` |
| Current grader | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest, raw bytes | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| Content-addressed `world_build_digest` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Stop on any mismatch. Do not modify protocol, runtime or prompts, arm schemas, Arm B, worlds, runner, grader, analysis, accepted freeze entries, `pre-validation-artifacts.json`, or either gate. `pre-validation-artifacts.json` must remain `status: partial`, with `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.

No Phoenix arm output, retained authoring outcome, prospective sealed outcome, runner result, implementation suggestion, or model-generated trial may influence design, labeling, audit, or packaging. If any prospective validation or held-out outcome is observed, treat the tranche as consumed, stop, and report the custody breach. Do not regenerate or relabel it.

## Private identities and mechanical verification

Create all private identities from the new artifacts. No digest, index, audit identity, or collision-inventory identity from either rejected package carries forward.

Use the original revision handoff's exact ordinal-byte algorithm for the new 240-label archive index. Store the exact index under the new private root. Compute and report a new raw SHA-256 for the independently created collision inventory.

Run all original test, vet, grader, schema, seal, dry-seal, path-index compatibility, identifier, leakage, allocation, variant, similarity, and audit checks from the new checkout. In addition:

- verify collision coverage against both `c193c786…` and `e705560…`;
- verify that no new family or mix layout forms a numbered or structural pair across tranches;
- inspect all public goals for copied grammar and cascade recipe cues;
- materialize and check the pre-edit and post-edit causal premise for all 24 cascades;
- scan new private labels for repeated non-schema label and rationale skeletons;
- run the exact rejected-byte comparison below; establish non-reuse of inaccessible prior private material through the clean-room custody record, not by opening or hashing a prior private root.

Keep the worktree clean after committing the candidate. Run both seal commands again without `--write` and require their raw stdout to equal the committed manifests byte-for-byte.

### Exact rejected-byte comparison

This is an exact whole-file equality screen. It supplements the structural collision audit; it does not replace it.

1. Let `B` be the set of Git blob object IDs returned by a recursive full-tree enumeration of evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`. Include blobs at every path and exclude tree, commit, and tag objects.
2. For each rejected payload `R` in `{c193c786cc5a65ce6ae97efbd336b2a48f492898, e705560fef6012f84ec0199b87ce790a62c17d0e}`, let `T_R` be the deduplicated set of blob object IDs returned by the same recursive full-tree enumeration of `R`.
3. Define `R_unique = T_R \ B`. Because Git blob IDs are content-addressed, this subtraction excludes every exact byte sequence already present as a blob anywhere in the evaluator base, including a base blob retained, moved, or copied to another path. Deleted base blobs are absent from `T_R` and therefore contribute no rejected bytes.
4. For each object in `R_unique`, read the exact raw blob bytes with a binary-safe Git object API and compute raw SHA-256. Deduplicate the resulting `(byte_length, sha256)` pairs. The union for both rejected payloads is `H_rejected`.
5. Build `H_new_public` from the new payload commit with the same full-tree enumeration and `B` subtraction.
6. Build `H_new_private` from every regular file recursively under `D:\Work\personal\phoenix-evaluator-private-v4-revision-2`. Do not follow directory symlinks or junctions. For each file, hash its exact raw bytes and record `(byte_length, sha256)`. The comparison is whole-file equality; do not scan for substrings embedded inside a new file.
7. Require `H_rejected ∩ (H_new_public ∪ H_new_private)` to be empty. Record the counts of base blobs, each rejected tree's blobs, rejected unique blobs, new public tree blobs, new public unique blobs, private regular files, and intersections so an independent auditor can reproduce the result.

Parse Git path lists and object bytes through a binary-safe API or NUL-delimited output. Do not pass blob contents through a PowerShell text pipeline, `Get-Content`, `Out-File`, `Set-Content`, or any decode/re-encode step. A hash or length mismatch means the files are not byte-equal. A matching `(byte_length, sha256)` pair is an exact-byte collision and blocks packaging until the new artifact is redesigned without using the rejected bytes.

## New package identities

Use the original two-commit procedure with new identities:

1. Commit only the approved public corpus, fixtures, label-digest registries, and sealed manifests. Its sole parent must be evaluator base `c852101e8d7cb52e4569bf3866de54a0ce648b44`. Do not include the public report.
2. Capture the native raw stdout bytes of `git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 <new-payload-commit>` with a binary-safe process API. Store those exact bytes at `D:\Work\personal\phoenix-evaluator-private-v4-revision-2\frontier-v1-v4-sealed-corpus-second-revision.patch`. Independently capture the diff again and require byte equality before recording its new SHA-256.
3. Add only `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-second-revision-report.md` in a report-only child commit.

The public payload may contain only the approved tranche replacements. It must contain no report, script, implementation change, protocol change, world change, grader change, freeze change, gate change, schedule, trial, outcome, full label, rationale, blueprint, audit, collision inventory, or private note.

The public report must satisfy every original report requirement and add disjointness evidence against `e705560…`, all-cascade pre-edit/post-edit causality coverage, new-label skeleton audit results, and an accurate account of mix-layout independence. It must not reveal expected paths, private causal details, or answers.

## Stop condition and handback

Stop after the new payload commit, raw public patch, private archive and index, collision inventory, independent private audits, report-only commit, and public report exist. Do not import the payload or modify the implementation workspace.

Hand back every item required by the original revision handoff, using only the new checkout, new private root, new payload, new patch, and new report identities. Include explicit confirmation that:

- `e705560…` was treated as a public collision set;
- neither prior private root was opened or reused;
- no prior script, fixture, blueprint, label, rationale, check layout, or label skeleton was reused;
- all 24 cascade terminal checks are causally dependent on earlier required acts;
- no model, arm, schedule, trial, prospective grade, validation result, held-out result, or outcome was run or observed;
- protected bytes are unchanged and both gates remain false.

The project chair must prepare a new independent corpus-and-label review assignment for this exact two-commit package and the new ordinal-byte archive index before review. The reviewer must be independent of both rejected packages and every second-revision authoring, audit, custody, and preparation role.

This handoff does not authorize import, schedule generation, Gate 1A, validation, held-out evaluation, or any outcome run.
