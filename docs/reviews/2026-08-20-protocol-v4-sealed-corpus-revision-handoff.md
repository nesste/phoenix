# Protocol-v4 sealed-corpus replacement revision handoff

## Assignment

Give this handoff to a fresh corpus author who did not create the rejected protocol-v4 candidate, perform its independent review, or implement Phoenix.

Create a wholly replacement `validation` and `held_out` corpus for `frontier-v1`. Work only in the dedicated revision checkout and a new private custody directory. The rejected candidate remains unopened custody evidence and an additional collision set. Do not repair it in place or reuse its generator, family blueprints, labels, private notes, or noun-swap structure.

This assignment produces outcome-free inputs and private labels. It does not import a patch, generate a schedule, open Gate 1A, run a model or experiment arm, execute a trial, or observe an outcome.

## Work locations and pinned history

| Item | Value |
| --- | --- |
| Revision checkout | `D:\Work\personal\phoenix-evaluator-v4-revision` |
| Revision branch | `codex/evaluator-v4-corpus-revision` |
| Evaluator base | `c852101e8d7cb52e4569bf3866de54a0ce648b44` |
| Implementation freeze | `4b51471200ed577db55fa38bed61c027c276d0f2` |
| Rejected candidate | `c193c786cc5a65ce6ae97efbd336b2a48f492898` |
| Rejected public patch | `sha256:c03aa6797ba472e114d0b61a10a739174b5b0906ef86934242664a8d9b114b65` |
| New private custody root | `D:\Work\personal\phoenix-evaluator-private-v4-revision` |
| Independent review record | `D:\Work\personal\phoenix\docs\reviews\2026-08-19-protocol-v4-sealed-corpus-review.md` |
| Original evaluator contract | `D:\Work\personal\phoenix\docs\reviews\2026-08-19-protocol-v4-sealed-corpus-evaluator-handoff.md` |

Read the original evaluator contract and the independent review record completely. All original allocation, schema, protocol-v4 semantics, custody, sealing, and no-outcome requirements remain in force unless this revision handoff makes them stricter.

Before authoring, independently recompute and require:

| Artifact | Identity |
| --- | --- |
| Protocol v4, LF-normalized UTF-8 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` |
| Production/authoring world, raw bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` |
| Current grader | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| World-build manifest, raw bytes | `sha256:da59f795b3670e4bf16ecae5453222374c9f91e50ebd3a7e9ddb6d1610487f0f` |
| Content-addressed `world_build_digest` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Stop on any mismatch. `pre-validation-artifacts.json` must remain `partial`, with only `schedule digest` remaining and both outcome gates false.

## Windows LF setup

Before any corpus command in the new clean revision checkout:

```powershell
git status --short
git config --worktree core.autocrlf false
git config --worktree core.eol lf
git checkout-index --all --force
git status --short
```

Stop if either status is nonempty. Do not run `checkout-index` in the implementation, rejected-candidate, or historical evaluator checkout.

## Custody and independence boundary

- Use a new private archive under the revision custody root. Do not modify or reuse `D:\Work\personal\phoenix-evaluator-private-v4`.
- Do not inspect the historical protocol-v3 private archive or the rejected candidate's private labels, template notes, adjudication records, or generator.
- The rejected candidate's public commit and public review findings are a collision set, not examples to imitate.
- Keep every new full label, rationale, family blueprint, template note, disagreement, and adjudication record outside every Phoenix checkout and outside the public patch.
- Record stable role identifiers for corpus author, first-pass labeler, independent label auditor, and private-label custodian. The later acceptance reviewer must be independent of all authoring roles.
- Do not use any Phoenix arm output, retained authoring outcome, prospective sealed outcome, runner result, implementation suggestion, or model-generated trial while designing or labeling cases.
- Treat a tranche as consumed if any prospective outcome is observed. Stop and report the breach; do not regenerate or relabel it.

## Four collision sets

The new candidates must be structurally disjoint from:

1. visible authoring families, cases, fixtures, goals, and labels;
2. retired protocol-v3 validation public material at evaluator base `c852101e…`;
3. retired protocol-v3 held-out public material at evaluator base `c852101e…`;
4. all public material in rejected candidate `c193c786…`.

Validation and held_out must also be structurally disjoint from each other.

Create a new outcome-free collision inventory before authoring. It must cover public paths, family IDs, case IDs, goals, fixture digests, file-map digests, observable fixture shapes, service aliases, test catalogs, adversarial patterns, and state-change patterns from all four collision sets. Use read-only Git-object access for the base and rejected candidate. Never open historical or rejected private labels to extend the inventory.

New family IDs and case IDs must be disjoint from every collision set. Do not recycle `validation_family_001` through `024` or `held_out_family_001` through `024`. A renamed noun, file, service, goal, identifier, or token is not a new generating template.

## Design the 48 families before emitting cases

Privately define one blueprint per family. Each blueprint records, without using prior private labels:

- the task mechanism and why it is a distinct generating idea;
- the domain and fixture topology;
- the observable states or tasks that distinguish all five variants;
- the intended capability sequence shape without exposing it publicly;
- any recovery, absence, adversarial, or shared-state mechanism;
- collision evidence against all four sets and the other replacement tranche.

Do not create one parameterized global template and instantiate it with different nouns. The eight mix families within a tranche must use eight genuinely different mechanisms and fixture structures. Held-out must not pair with validation through renamed domains, mirrored family layouts, goal grammar, fixture topology, alias tables, state-change plans, or label skeletons.

An independent private blueprint audit must pass before labels are finalized. Record disagreements and adjudications privately. The public report may summarize non-outcome variant axes and collision counts but must not reveal expected paths or answers.

## Frozen allocation

For each tranche, independently preserve the protocol allocation:

- exactly 120 cases in 24 families, five cases per family;
- eight direct-group families, each with 3 `direct` + 1 `absence` + 1 `stale_frontier`;
- eight recovery-group families, each with 3 `recovery` + 1 `far_discovery` + 1 `temptation`;
- eight mix families containing, across the group, 12 `cascade`, 12 `adversarial_text`, and 4 each of `far_discovery`, `temptation`, `absence`, and `stale_frontier`;
- at least one `cascade` and one `adversarial_text` case in every mix family;
- the frozen class totals and distinct-family floors.

## Five real variants per family

Every family's five cases must differ in observable state or task. Goal wording, a target token, a noun, a filename, or a service alias alone does not establish a variant.

For the three same-class cases in every direct and recovery family, require a material difference in initial state, task type, intermediate observation, success condition, recovery condition, or check evidence. Two cases may share a capability path only when their state or task difference changes what the agent must observe or establish. Do not accept three paraphrases over one snapshot and one grading path.

Retain a private five-row variant matrix per family. The public report must give outcome-free evidence for each family's five distinct variants, not merely count unique goal strings.

## Cascade requirements

- A cascade succeeds through at least three executable verbs in an order the agent cannot infer from the public goal alone.
- Public cascade goals state the desired outcome and relevant constraints. They must not enumerate, synonymize, or strongly cue the success recipe.
- The required next step must emerge from an intermediate result, frontier, refusal, or changed state.
- Treat authoring `repo.status -> tests.run -> tests.list -> tests.focus` as a collision set, not a template. If a genuinely distinct diagnostic case uses `tests.run` followed by `tests.focus`, it must not omit the intervening `tests.list` discovery step.
- Vary cascade mechanisms and action-graph shapes across families and tranches. Do not repeat one edit recipe with different nouns, tokens, files, or domains.
- Orientations never count as executable acts. Shared state events use the same declared executable-action index for Arms A-E.

## Labels and grading

Create one new full label per case in the new private archive. Preserve the current schema and grader pin. Each label must distinguish success from a plausible wrong path and must be reasoned against the frozen world, case, fixture, and shared state events without running an arm.

For every listed `acceptable_path`, mechanically verify that all path-bound checks can succeed on that path. In particular, an `act_sequence` `seq` check must not require an index absent from a shorter listed path. Remove invalid shorter paths or redesign the checks; do not list a path the grader cannot score as successful.

Use executable-act indexes only. Include an empty acceptable path when no action is expected. Keep deterministic grading where possible; any human judgment requires `arm_hidden: true`, a precommitted concrete rubric, and retained independent adjudication.

Canonicalize each label with `corpusctl digest`. Public registry rows contain only `case_id`, `class`, `label_digest`, and `grading_script`.

## Private identities

Define the new private-label archive digest exactly as follows:

1. For every new label file, compute SHA-256 of raw file bytes.
2. Form the UTF-8 line `<archive-relative-path>\tsha256:<lowercase-hex>`. Paths are relative to the new private custody root, begin with `labels/`, and use `/` separators.
3. Sort all 240 lines by archive-relative path using ordinal byte order.
4. Join the lines with LF and terminate the index with LF.
5. SHA-256 the resulting index bytes.

Store that exact index privately and report its digest and custodian. The reported digest, stored index bytes, and documented algorithm must agree.

Store the new outcome-free collision inventory privately. Report its raw-file SHA-256, algorithm, and custodian.

## Mechanical checks

From `experiments/frontier-v1/corpusctl` in the revision checkout:

```powershell
go test -count=1 ./...
go vet ./...
go run ./cmd/corpusctl grader-digest --repo-root ../../..
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json --write
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json --write
```

Run each seal again without `--write` and compare its raw output byte-for-byte with the committed manifest. Keep the worktree clean after committing the candidate.

Add explicit private checks for:

- pairwise structural similarity across all 48 blueprints;
- validation-to-held-out mirrored families;
- five real variants per family;
- cascade goals that leak or cue their action sequence;
- acceptable paths incompatible with indexed checks;
- identifier reuse across all collision sets;
- public registry/report leakage;
- missing or extra labels and digest mismatches.

Automated similarity checks are screening evidence, not proof of independence. Independently inspect every family blueprint and every full label.

## Non-circular public package identities

Use two commits so the report can identify the candidate and patch without self-reference:

1. **Payload candidate commit.** Commit only the approved public corpus, fixture, registry, and sealed-manifest replacement. Do not include the public report in this commit.
2. Generate `git --no-pager diff --binary c852101e8d7cb52e4569bf3866de54a0ce648b44 <payload-candidate-commit>` as raw bytes and record its SHA-256. This is the public replacement patch.
3. **Report commit.** Add only `docs/reviews/2026-08-20-frontier-v1-v4-sealed-corpus-revision-report.md` as a child of the payload candidate. The report records the evaluator base, payload candidate commit, exact public-patch digest, report commit's parent relationship, private-archive digest, collision-inventory digest, counts, checks, custody declarations, and protected-artifact statement.

The report commit is the evaluator branch tip. The reviewed import payload is the base-to-payload-candidate patch. The report is a separate review artifact and is not included in that patch. Do not create a combined patch whose digest would depend on a report that embeds the digest.

The public payload diff may contain only the retired-v3 public deletions and the new replacement cases, fixtures, label-digest registries, and sealed manifests at the approved paths. It must contain no report, implementation change, protocol change, freeze change, gate change, schedule, trial, outcome, full label, rationale, or private note.

## Required public report

In addition to the original evaluator-report requirements, the revision report must include:

- the exact payload candidate commit and public-patch SHA-256;
- a non-outcome family blueprint summary proving 48 distinct mechanisms without revealing answers;
- a five-variant evidence matrix for every family using observable state/task axes, not goal-count evidence;
- disjointness evidence against authoring, retired-v3 validation, retired-v3 held-out, the rejected v4 candidate, and the other new tranche;
- cascade-goal audit counts and confirmation that goals do not state or cue recipes;
- acceptable-path/check compatibility counts;
- new family-ID and case-ID disjointness results;
- exact archive-index algorithm and digest;
- test, vet, grader, seal, dry-reproduction, similarity-screen, and independent private-audit results;
- confirmation that no model, arm, trial, schedule, or outcome was run or observed;
- confirmation that protected implementation and freeze bytes are unchanged and both gates remain false.

Proofread public goals and report text. Correct cloned wording and grammar defects before freezing the payload candidate.

## Stop condition and handback

Stop after the payload candidate commit, raw public patch, new private archive and index digest, new collision inventory and digest, report commit, and public report exist. Do not apply the payload to the implementation workspace.

Hand back:

- evaluator branch and report commit;
- evaluator base and payload candidate commit;
- public patch path, raw SHA-256, and byte-equality command;
- public report path;
- private-label archive digest, exact algorithm, and custodian;
- collision-inventory digest, algorithm, and custodian;
- family/variant/cascade/path-compatibility audit summaries;
- any blocking issue or custody breach.

The project chair must obtain a new independent corpus-and-label review before deciding whether the replacement public patch may be imported. This handoff does not authorize import, a validation schedule, Gate 1A, validation, held-out, or any outcome run.
