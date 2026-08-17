# Sealed corpus evaluator handoff

## Assignment

Act as the corpus evaluator independent of Phoenix implementation. Create unopened `validation` and `held_out` candidates for `frontier-v1`. Work in a separate checkout and keep every full validation and held-out label outside that checkout.

This assignment creates and audits inputs. It does not run an agent, open outcomes, tune Phoenix, or authorize Phase 1.

## Independence boundary

- Do not use a Phoenix implementation, arm output, trial result, or prospective validation or held-out outcome while designing cases or labels.
- Do not copy full validation or held-out labels into the Phoenix checkout, a patch, a command log, or a review report. The repository-wide seal guard rejects them by design.
- Do not derive a family from an authoring fixture or template. Validation and held-out must also be disjoint from each other.
- Record who created the cases, who labelled them, what Phoenix material each person saw, and every disagreement or adjudication. The label reviewer must not be the Phoenix implementer.
- Treat a tranche as consumed if any outcome is observed. Stop and report the breach; do not replace a case silently.

The evaluator may use these contracts:

- `docs/plans/2026-08-17-phoenix-world-plan.md`, Task 0.4;
- `experiments/frontier-v1/protocol.json`;
- `experiments/frontier-v1/schema/`;
- `spec/world.schema.json` and `experiments/frontier-v1/worlds/authoring.dev_repo.json`;
- `experiments/frontier-v1/corpusctl/` for canonical digests, validation, grading, and sealing.

Avoid the authoring cases, fixtures, and labels while generating new families. The separate checkout may contain them so `corpusctl seal` can detect exact template or fixture reuse, but they are a collision set, not examples to imitate.

## Required design for each tranche

Create exactly 120 cases in 24 generating families, with five distinct case variants per family. Use opaque case IDs; an ID, goal, filename, fixture text, or family name must not reveal the class or expected answer.

Use this allocation independently for validation and held_out:

| Family group | Families | Cases in each family |
| --- | ---: | --- |
| Direct | 8 | 3 `direct`, 1 `absence`, 1 `stale_frontier` |
| Recovery | 8 | 3 `recovery`, 1 `far_discovery`, 1 `temptation` |
| Mix | 8 | 5 cases drawn from the totals below |

Across the eight mix families, allocate 12 `cascade`, 12 `adversarial_text`, and four each of `far_discovery`, `temptation`, `absence`, and `stale_frontier`. Put at least one `cascade` and one `adversarial_text` case in every mix family. The resulting tranche totals must be:

| Class | Cases | Minimum distinct families |
| --- | ---: | ---: |
| `direct` | 24 | 8 |
| `recovery` | 24 | 8 |
| every other required class | 12 each | 8 each |

Within a family, five cases may share a generating idea but must vary the observable state or task. Across families and tranches, change the generating template rather than renaming files or rewriting prose around the same answer path.

## Labels and grading

Create one schema-valid full label per case outside the checkout. Prefer deterministic checks over `final_message_states`. When judgment cannot be avoided, set `arm_hidden: true`, write a concrete rubric, and retain independent adjudication records.

Each label must:

- use the same opaque `case_id` as its input;
- state the evaluator-only class;
- pin the current `corpusctl grader-digest`;
- define checks that distinguish success from a plausible wrong path;
- list non-exhaustive acceptable verb paths, including an empty path where acting is not expected;
- explain the expected outcome and why the checks establish it.

Canonicalize each private label with `corpusctl digest`. Put only `case_id`, `class`, `label_digest`, and `grading_script` into the outcome-free registry.

## Candidate package

Prepare one public patch containing only:

```text
experiments/frontier-v1/corpus/validation/*.json
experiments/frontier-v1/corpus/held_out/*.json
experiments/frontier-v1/fixtures/validation/*.json
experiments/frontier-v1/fixtures/held_out/*.json
experiments/frontier-v1/manifests/validation-label-digests.json
experiments/frontier-v1/manifests/held_out-label-digests.json
experiments/frontier-v1/manifests/validation.json
experiments/frontier-v1/manifests/held_out.json
docs/reviews/<date>-frontier-v1-sealed-corpus-report.md
```

Keep outside the patch and repository:

- full labels and label rationales;
- class-to-template working notes that expose answers;
- disagreements and adjudication details that reveal expected outcomes;
- any trial or model output.

In the separate checkout, generate the manifests from `experiments/frontier-v1/corpusctl`:

```powershell
go test ./...
go vet ./...
go run ./cmd/corpusctl grader-digest --repo-root ../../..
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/validation-label-digests.json --write
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche held_out --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --label-digests experiments/frontier-v1/manifests/held_out-label-digests.json --write
```

`seal` enforces schemas, references, current grader digests, missing or extra label digests, and exact cross-tranche template and fixture-content collisions. It does not enforce every numeric allocation in `protocol.json`; verify the tables above independently from the generated manifests and include the arithmetic in the report.

## Report to the Phase 0 reviewers

Submit the public patch plus a report that contains no expected outcomes and records:

1. evaluator identities or stable role identifiers and independence declarations;
2. the base commit and public-patch digest;
3. counts by tranche, class, family, and fixture, with evidence that every family has five distinct case variants;
4. a family-allocation matrix proving the required distribution without describing answers;
5. confirmation that family templates and fixture contents do not overlap across authoring, validation, and held_out;
6. schema-validation and sealing commands with exit status;
7. the number of scripted, regex, and blinded-human checks by kind;
8. disagreement counts and whether all disagreements were adjudicated, without private label content;
9. confirmation that no agent or arm was run and no outcome was observed;
10. the external private-label archive digest and its custodian.

Do not apply the patch to the implementation workspace. The Task 0.6 reviewers inspect the candidate first. The project chair decides whether it may be imported.
