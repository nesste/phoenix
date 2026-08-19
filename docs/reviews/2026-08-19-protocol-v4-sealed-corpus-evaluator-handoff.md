# Protocol-v4 sealed-corpus evaluator handoff

## Assignment

Act as the corpus author and evaluator independently of Phoenix implementation. Create unopened protocol-v4 replacement candidates for the `validation` and `held_out` tranches of `frontier-v1`.

Work only in the dedicated evaluator checkout. Keep every full validation and held-out label, rationale, template note, disagreement, and adjudication record outside every Phoenix checkout. This assignment creates and audits outcome-free inputs. It does not run a model, an experiment arm, a trial, or an outcome analysis, and it does not open Gate 1A.

## Pinned implementation contract

The handoff begins after the focused local-artifact freeze at implementation commit `4b51471200ed577db55fa38bed61c027c276d0f2`.

| Artifact | Identity |
| --- | --- |
| Protocol v4 | `sha256:81c86f1eea00120e597927cb1937854fba6580913ba79d604c4b63c7caaa05e9` |
| Production/authoring world, canonical JSON | `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` |
| Production/authoring world, raw bytes | `sha256:41d242e672c812a50a253e33e165d839e2c8d167914605805a311fefaeb92143` |
| Current grader | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| Frozen Linux-amd64 world build | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Independently recompute these identities before authoring. Stop if any differs. `pre-validation-artifacts.json` must remain `partial`, with both outcome gates false and only `schedule digest` remaining.

## Independence and custody boundary

- Do not use a Phoenix arm output, runner result, retained authoring outcome, prospective sealed outcome, or implementation suggestion while designing or labeling cases.
- Do not run `experiments/frontier-v1/runner`, Claude Code, another model, or any tool that executes an experiment arm.
- Do not copy full sealed labels into the evaluator checkout, implementation checkout, Git history, public patch, logs, or report.
- Assign stable role identifiers for corpus author, first-pass labeler, label auditor, and private-label custodian. The auditor must not be the Phoenix implementer and must record disagreements without exposing answers publicly.
- Treat a tranche as consumed if any prospective outcome is observed. Stop and report the breach; do not silently regenerate or relabel it.
- Do not modify protocol, runtime/prompts, arm schemas, Arm B, world, runner, grader, analysis, or any accepted freeze entry.

## Three collision sets

The new candidates must be structurally disjoint from all three of these sets:

1. the visible authoring families, fixtures, goals, and labels;
2. the retired unopened protocol-v3 validation material;
3. the retired unopened protocol-v3 held-out material.

The tracked v3 public files are custody evidence, not runnable candidates and not templates to imitate. Before removing them from the evaluator branch, retain an outcome-free collision inventory outside the candidate patch containing their tracked paths, family IDs, case IDs, goals, fixture digests, and file-map digests. Do not inspect the historical private-label archive.

Do not reuse or lightly rewrite a v3 family, case, fixture, goal grammar, answer path, service alias, test catalog, adversarial instruction, or state-change pattern. New family IDs and case IDs must be disjoint from authoring and both retired tranches. A renamed file or reworded goal is not a new generating template.

## Replacement workflow in the evaluator checkout

The current branch contains the retired v3 public tranche files at the live `validation` and `held_out` paths. The replacement candidate must replace them rather than mix v3 and v4 cases.

1. Record the evaluator checkout base commit and verify a clean worktree.
2. Build the outcome-free retired-v3 collision inventory described above.
3. Remove the tracked v3 files from these paths in the evaluator branch only:

   ```text
   experiments/frontier-v1/corpus/validation/*.json
   experiments/frontier-v1/corpus/held_out/*.json
   experiments/frontier-v1/fixtures/validation/*.json
   experiments/frontier-v1/fixtures/held_out/*.json
   experiments/frontier-v1/manifests/validation-label-digests.json
   experiments/frontier-v1/manifests/held_out-label-digests.json
   experiments/frontier-v1/manifests/validation.json
   experiments/frontier-v1/manifests/held_out.json
   ```

   These deletions are part of the eventual public replacement patch. Git history and the historical public report retain custody evidence.
4. Create the new v4 public inputs and fixtures at those same tranche paths.
5. Create all full labels in a new private archive outside every Phoenix checkout. Do not reuse the v3 private archive.
6. Create outcome-free label-digest registries in the evaluator checkout.
7. Run the seal commands and independently audit the generated manifests.
8. Produce a public report and patch. Do not apply either to the implementation workspace.

## Required design for each replacement tranche

Create exactly 120 cases in 24 new generating families, with five distinct case variants per family. IDs, filenames, goals, fixture text, and public family names must not reveal the class, expected answer, acceptable path, or grading logic.

Use this allocation independently for validation and held-out:

| Family group | Families | Cases in each family |
| --- | ---: | --- |
| Direct | 8 | 3 `direct`, 1 `absence`, 1 `stale_frontier` |
| Recovery | 8 | 3 `recovery`, 1 `far_discovery`, 1 `temptation` |
| Mix | 8 | 5 cases drawn from the totals below |

Across the eight mix families, allocate 12 `cascade`, 12 `adversarial_text`, and four each of `far_discovery`, `temptation`, `absence`, and `stale_frontier`. Every mix family must contain at least one `cascade` and one `adversarial_text` case.

The resulting totals per tranche are:

| Class | Cases | Minimum distinct families |
| --- | ---: | ---: |
| `direct` | 24 | 8 |
| `recovery` | 24 | 8 |
| each other required class | 12 | 8 |

Every family has five variants with distinct observable state or task. Validation and held-out must use different generating templates from one another and from every collision set.

## Protocol-v4 case requirements

- Inputs use the current strict case schema and current canonical world digest.
- `state_changes`, when present, are deterministic relative-path replacements applied after a numbered executable act. The same event plan and executable-action index apply to Arms A-E; orientations never advance the index.
- Stale-frontier cases must create their stale condition through a declared shared state event, not timing, an external process race, or arm-specific behavior.
- Cases must remain solvable through the frozen world and its twelve flat-tool equivalents. Do not require an unregistered capability or a host path outside the fixture.
- Adversarial fixture text is untrusted data. It must test the frozen surface contract without embedding class names, grader hints, or answer metadata in the public input.
- Capability-absence cases must have a genuinely absent capability; labels may allow at most the precommitted discovery behavior and must not reward unsupported action.
- Do not relax the restored cascade contract or select a favorable path from authoring variance.

## Private labels and grading

Create one schema-valid full label per case outside the checkout. Each label must:

- use the same opaque `case_id` as its public input;
- record the evaluator-only class;
- pin grader `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`;
- define checks that distinguish success from at least one plausible wrong path;
- use executable-act indexes only; orientations are separate evidence and never advance a path index;
- list non-exhaustive acceptable paths, including an empty path when acting is not expected;
- explain the expected outcome and why the checks prove it.

Prefer deterministic checks. If human judgment is unavoidable, set `arm_hidden: true`, define a concrete rubric before outcomes exist, and retain independent adjudication records. Canonicalize every private label with `corpusctl digest`. The public registry contains only `case_id`, `class`, `label_digest`, and `grading_script`.

## Mechanical commands

From `experiments/frontier-v1/corpusctl` in the evaluator checkout:

```powershell
go test ./...
go vet ./...
go run ./cmd/corpusctl grader-digest --repo-root ../../..
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json --write
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json --write
```

Run each seal command again without `--write` and compare its output to the committed manifest. `seal` checks schemas, references, grader identity, missing/extra label digests, and exact current-tree cross-tranche collisions. It cannot prove non-derivation from deleted v3 templates; the independent retired-v3 collision audit must cover that gap.

## Public replacement patch

The public diff against the evaluator base may contain only:

- deletion of the retired v3 public files from the eight paths listed above;
- new v4 cases and fixtures at those paths;
- new outcome-free label-digest registries and sealed manifests;
- one new report at `docs/reviews/<date>-frontier-v1-v4-sealed-corpus-report.md`.

It must not contain full labels, rationales, template notes, private disagreements, model output, trials, schedules, implementation changes, or gate changes.

## Required public report

The report contains no expected outcomes and records:

1. evaluator roles, independence declarations, and what Phoenix material each role saw;
2. evaluator base commit, implementation freeze commit, and public-patch digest;
3. counts by tranche, class, family, and fixture, proving five variants per family;
4. the family-allocation matrix and exact arithmetic;
5. disjointness evidence against authoring, retired-v3 validation, retired-v3 held-out, and the other replacement tranche;
6. schema, test, vet, digest, sealing, and dry-reproduction commands with exit codes;
7. scripted, regex, and blinded-human check counts by kind from the private labels;
8. disagreement and adjudication counts without private content;
9. confirmation that no agent, arm, trial, schedule, or outcome was run or observed;
10. the canonical digest of the new private-label archive and its custodian;
11. the retired-v3 collision-inventory digest and custodian;
12. a protected-artifact statement confirming no implementation or freeze byte changed.

## Stop condition and handback

Stop after producing the public replacement patch, public report, and private archive digest. Do not import the patch into the implementation workspace. Do not generate a validation schedule. Do not run a model. Do not change `pre-validation-artifacts.json` or either outcome gate.

Hand back:

- evaluator branch and exact candidate commit;
- public patch path and SHA-256;
- public report path;
- private archive digest and custodian;
- collision-inventory digest and custodian;
- any blocking issue or custody breach.

The implementation project chair must obtain independent corpus and label review before deciding whether the public replacement patch may be imported.
