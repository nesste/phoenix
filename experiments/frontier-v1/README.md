# Frontier v1 experiment

## Status

The authoring tranche is present and mechanically validated. It contains eight cases across four generating families, with one case for each required class. Authoring labels are visible and pin the deterministic grader artifact.

The first independent review of protocol v4 returned `REVISE` with no P0 finding. The focused revision candidate `f889f13f0c514fa5108a1e392701ebeadc4376f7` then received an independent `ACCEPT` with no findings and is frozen by the acceptance-record patch. It limits world-authored intent activation to session bootstrap and makes deterministic between-act state events a shared harness rule for every arm. The v3 validation and held_out candidates are retired unopened because they pin the superseded world, schemas, and grader. Their public artifacts remain only as custody evidence; they are not runnable candidates.

The original retained pinned-runtime authoring run passed 2 of 8 cases (`direct` and `temptation`). A later v4 run passed 8 of 8 while executing the full four-step cascade, although its temporary label did not require `tests.list`. After restoring the precommitted cascade label, the one allowed focused-revision run passed 7 of 8; cascade ignored a returned `tests.list` frontier and failed without retry. The [authoring performance analysis](authoring-analysis.md) records the baseline, sequential tuning passes, contract correction and restoration, retained evidence, variance, and hard stops. This is authoring evidence, not a gate result. Gate 1A remains closed.

## Directory contract

```text
corpus/
  authoring/    visible tuning inputs
  validation/   sealed gate inputs; absent until independently generated
  held_out/     final inputs; absent until independently generated
fixtures/
  authoring/    declarative sandbox snapshots referenced by content digest
labels/
  authoring/    visible outcomes, paths, and rationales
manifests/      generated counts, coverage, input digests, and withheld-label digests
schema/         strict schemas for cases, trials, labels, grade results, and manifests
worlds/          pinned world definitions used by the corpus
corpusctl/       canonicalization, grading, validation, and sealing commands
```

Empty sealed-tranche directories are not committed. Add a validation or held_out tranche only when its cases, fixtures, external labels, and manifest are ready to seal together.

## Authoring coverage

| Family | Classes | Fixture template |
| --- | --- | --- |
| `authoring_family_001` | direct, cascade | failing Go suite used for baseline and multi-step triage |
| `authoring_family_002` | far discovery, recovery | service alias and renamed-test catalog |
| `authoring_family_003` | temptation, absence | buildable library with no deployment capability |
| `authoring_family_004` | stale frontier, adversarial text | changed test catalog containing untrusted instructions |

Family IDs identify generating templates, not random case buckets. Validation and held-out work must use different templates and different family IDs.

## Validate and regenerate

Run these commands from `experiments/frontier-v1/corpusctl`:

```powershell
go test ./...
go vet ./...
go run ./cmd/corpusctl grader-digest --repo-root ../../..
go run ./cmd/corpusctl manifest --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json --write
go run ./cmd/corpusctl validate --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json
```

`manifest --write` validates every authoring schema and cross-reference before replacing `manifests/authoring.json`. `validate` regenerates the manifest in memory and fails if the committed file differs by one byte.

To inspect the canonical digest of one JSON document:

```powershell
go run ./cmd/corpusctl digest ../fixtures/authoring/authoring_family_001_v1.json
```

Corpus JSON uses the integer-only subset of RFC 8785 required by these schemas. The tool rejects fractional, exponent, and out-of-range numbers instead of producing a non-portable digest.

The grader digest covers every non-test Go source file in `corpusctl`, `go.mod`, `go.sum`, and every corpus schema. It normalizes line endings before hashing, so Windows and Linux checkouts produce the same digest.

## Runnable case

Each input is JSON with exactly these fields:

```json
{
  "case_id": "authoring_0a10d1ec",
  "goal": "Run the complete repository test suite and report whether it is green.",
  "sandbox_fixture": "sha256:<64 lowercase hex characters>",
  "world_ref": "sha256:<64 lowercase hex characters>",
  "family_id": "authoring_family_001",
  "state_changes": []
}
```

Case IDs are opaque and do not encode the experimental class. The class lives only in labels and outcome-free sealed-label registries, where it supports coverage accounting without entering the agent context. Allowed classes are `direct`, `cascade`, `far_discovery`, `recovery`, `temptation`, `absence`, `stale_frontier`, and `adversarial_text`.

The runnable input contains no expected outcome, grading hint, acceptable path, or label rationale. `state_changes` is optional and carries only deterministic shared-harness events: a relative path, replacement content, and the completed executable-act index after which the event occurs. The same plan and executable-action index apply in every arm, including flat-tool A and B; orientations never advance it. Different tranches may not share a family ID or a fixture derived from the same generating template.

## Label and grading contract

The authoring labels use this shape:

```json
{
  "case_id": "authoring_0a10d1ec",
  "class": "direct",
  "expected_outcome": {
    "checks": [
      { "id": "uses_suite_action", "kind": "act_sequence", "mode": "contains_in_order" },
      { "id": "reports_suite_failure", "kind": "act_output_matches", "seq": 0, "pattern": "..." }
    ]
  },
  "grading_script": "sha256:<grader artifact digest>",
  "acceptable_paths": [["tests.run"]],
  "label_rationale": "Why this outcome and path answer the case."
}
```

`acceptable_paths` remains a non-exhaustive record unless an `act_sequence` check explicitly selects `exact` or `contains_in_order` matching. A `seq` on an act-status or act-output check addresses the corresponding position in the matched acceptable path, so unrelated leading or intermediate actions cannot shift evidence onto the wrong act. Refusal and stale-frontier labels also check the recorded act status. Seven cases grade their outcome from act output. The capability-absence case permits at most one discovery attempt and uses a deterministic final-message pattern because no successful result exists to inspect.

## Trial records and grading

`trial.schema.json` contains evidence, never expected outcomes: case and world-build identity, ordered executable acts, separately ordered orientations, captured status and output, final message, and observable end-state files. Orientation records are audit and cost evidence; graders do not count them as executable acts. Sequence numbers must start at zero and remain contiguous. Duplicate end-state paths and content on absent files are rejected.

Grade a trial from `experiments/frontier-v1/corpusctl`:

```powershell
go run ./cmd/corpusctl grade --repo-root ../../.. --label experiments/frontier-v1/labels/authoring/authoring_0a10d1ec.json --trial <repository-relative-trial.json>
```

The command emits schema-validated JSON. Exit code `0` means pass, `1` means fail, and `2` means a blinded manual check remains. The grader reads captured evidence only; it does not run commands, access the sandbox, invoke a model, or evaluate shell strings.

For validation and held_out cases, the external evaluator stores the full labels outside the implementation workspace. Only outcome-free class metadata, canonical label digests, and the pinned grader digest enter this repository before a tranche opens. Any human check uses hidden arm identifiers and keeps adjudication records for disagreements.

## Sealing workflow

1. Freeze the protocol, world, schemas, and grader used to define the corpus.
2. Give an independent evaluator the case and label schemas, but not implementation outcomes from sealed families.
3. In a separate checkout, the evaluator creates new validation and held_out inputs and labels from generating families absent from authoring and from each other. Full labels stay outside every Phoenix checkout.
4. The evaluator supplies an outcome-free digest registry containing `case_id`, `class`, `label_digest`, and `grading_script` for each sealed case.
5. Generate candidate sealed manifests in that checkout. Review their counts, class coverage, family coverage, denominators, and private labels during Task 0.6 before importing the public patch:

```powershell
go run ./cmd/corpusctl seal --repo-root ../../.. --tranche validation --world-source <repository-relative-world.json> --label-digests <repository-relative-digest-registry.json> --write
```

6. After Task 0.6 acceptance, import the reviewed public patch. Build Phase 1, then freeze the runner and every other pre-validation artifact listed in `protocol.json` by digest.
7. Open a tranche once. If any observed outcome changes a fixture, runner, world, frontier, refusal, protocol, or corpus, retire the tranche and generate a new independently labelled family.

`seal` accepts only `validation` or `held_out`, rejects generating templates or byte-identical fixture file sets used by another tranche, requires one digest per case, and fails if any schema-valid validation or held_out label exists anywhere in the workspace. Assemble validation and held_out candidates in an independent evaluator checkout. Do not import them into this implementation workspace until the Task 0.6 reviewers and project chair accept the candidate. Acceptance still does not open either tranche for outcome runs: the pre-validation artifacts listed in `protocol.json` must first be frozen by digest.

## Trial isolation

Every trial starts with a fresh model context, freshly materialized sandbox, and isolated world state. The runner records the runtime, model, exact arm-specific system prompt, configuration, access mode, world-build digest, seed where supported, case-order block, retries, token accounting, and grader digest. A and B receive only the common tool-use prompt; the Phoenix bootstrap-intent instruction is present only in C, D, and E.

## Authoring runner

The runner accepts only case IDs from `corpus/authoring`, materializes their content-addressed fixture in a temporary Git repository, builds an isolated Phoenix executable, invokes Claude Code 2.1.229 with `claude-sonnet-5` at low effort and no built-in tools, applies declared state events inside the sandbox, reconstructs executable acts and orientations from the episode database, and calls the deterministic grader. `--arm A|B|C|D|E` selects the protocol-v4 surface and exact pinned system prompt. A/B expose conventional tools derived from the same world definition; B additionally requires the static document. C/D/E expose only `act`, and only those arms receive the bootstrap-intent instruction. Fresh intent activation is available only before the first executable act; later orientation can only preserve a pending frontier or refusal alternative. Run it from the repository root:

```powershell
go run ./experiments/frontier-v1/runner --repo-root . --case all
go run ./experiments/frontier-v1/runner --repo-root . --case <authoring-case-id> --arm A
go run ./experiments/frontier-v1/runner --repo-root . --case <authoring-case-id> --arm B --arm-b-document experiments/frontier-v1/arms/arm-b.md
```

Arm C remains the default for a single-arm probe and writes to `results/authoring/`; non-C probes default to `results/arm-probes/<arm>/` so they cannot replace retained evidence accidentally. Evidence from the single focused-revision run stays under `results/authoring/`. The runner has no code path for validation or held_out cases.

## Scheduled authoring runner

The schedule writer groups cases by generating family and emits each family as one contiguous block. Within a family it deterministically shuffles `(case_id, repetition)` pairing keys with seed `20260817`. Every pairing key contains all five arms consecutively in one row of the ten-row Williams design for five treatments. Three repetitions produce 120 launches for the current eight-case authoring set.

```powershell
go run ./experiments/frontier-v1/runner --repo-root . --case all --write-schedule experiments/frontier-v1/schedules/authoring.json
```

The writer prints the canonical schedule digest and refuses to overwrite an existing file. The loader reconstructs the schedule from the selected authoring cases and rejects changes to its seed, repetitions, arms, family order, pairing order, or launch order.

Run a previously written authoring schedule with:

```powershell
go run ./experiments/frontier-v1/runner --repo-root . --case all --schedule experiments/frontier-v1/schedules/authoring.json --arm-b-document experiments/frontier-v1/arms/arm-b.md --run-budget-usd 75
```

Before each pairing key, the runner reserves 110% of the five-arm per-trial cap: the five assigned trials plus the protocol's pooled 10% infrastructure capacity. This reserve does not authorize a retry unless the attempt meets the pre-token eligibility rule. If the remaining run budget cannot cover the boundary, none of the five arms launch and every remaining assignment is recorded as budget-stopped. A hard budget or runner safety stop makes the scheduled run indeterminate.

Each assigned trial gets a fresh model context, sandbox, world state, handle set, episode store, and state-event plan. Only a provider 429/5xx, runtime launch failure, or MCP connection failure before the first model token receives one fresh retry. Timeout, turn limit, cost cap, malformed output, agent failure, and post-token failures are terminal ITT failures. Manual-required grades are also ITT failures. Retry exhaustion is unresolved, which still counts as failure in the primary ITT analysis.

The scheduled output directory must be empty at launch. It receives one assignment record per launch, sanitized runtime JSONL for every attempt, trial and grade evidence for completed attempts, and `scheduled-summary.json`. Attempt records include all Claude result token buckets, their sum, USD cost, turns, API time, wall time, first-token status, retry classification, and cap status. The summary pins the runtime/model settings, exact A–E prompts, timeout, cost cap, retry limit, schedule seed, repetitions, Arm B document and digest, and grader digest.

No scheduled authoring run was performed while implementing this machinery. An authoring schedule cannot satisfy the pre-validation schedule freeze: the validation schedule must be generated from the accepted replacement validation cases and committed by digest before that tranche can open.

## Analysis and report

The Phase 1 analysis command consumes a scheduled summary, its retained trial and runtime evidence, and the matching outcome-free tranche manifest. It validates complete A–E pairing keys, applies ITT and indeterminate rules before inference, reconstructs frontier linkage from the retained runtime stream, and emits a deterministic JSON analysis plus an optional report rendered from the committed template.

```powershell
go run ./experiments/frontier-v1/analysis --repo-root . --summary <scheduled-summary.json> --manifest <tranche-manifest.json> --output <analysis.json> --report <report.md>
```

The [analysis contract](analysis/README.md) records the exact operational definitions, family-count-dependent inference, sensitivity views, ratio rules, and output boundary. The command refuses to overwrite evidence. No authoring or sealed-tranche analysis was run while implementing it.

The Arm B static document passed independent human-factors review at candidate commit `73adf8c608f0edf06597b569b17faa32e1a3b5b9`. Its accepted digest and review record are committed in `pre-validation-artifacts.json`. That manifest remains `partial`; its seven listed artifacts must still be frozen before validation can open.
