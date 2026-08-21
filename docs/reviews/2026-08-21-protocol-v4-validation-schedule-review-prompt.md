# Independent review assignment: protocol-v4 validation schedule

Review the outcome-free validation-schedule candidate as `validation_schedule.independent_reviewer`.

You must be distinct from `validation_schedule.candidate_author`, the project-chair corpus import session, the third-revision corpus remediation and review roles, and every session that wrote the candidate schedule, generator, evidence record, or this assignment. You may be a fresh agent session on the same machine. Record that limitation.

This review does not authorize a schedule execution, model or arm run, prospective grade, Gate 1A, validation, held-out evaluation, outcome analysis, or access to a private label. It decides only whether the exact public validation schedule is eligible for a focused digest-freeze patch.

## Pinned public identities

- Review base and candidate parent: `2822b1e8d2b3c5675d52eec9fe3435349618bcc0`
- Candidate commit: `44fccf51cb1684a1b71da9128cab30ad2b2fb6af`
- Candidate-report commit: `bb473e22b1e2a5b21221a0b7b57e2aec4836c52b`
- Required ancestry: report commit's sole parent is the candidate; candidate's sole parent is the review base
- Candidate report: `docs/reviews/2026-08-21-protocol-v4-validation-schedule-candidate.md`
- Candidate report raw bytes: 7,731 bytes; `sha256:e491d484e0c768feedcbc6571153530fe7a5f6ea91b4be241415c3451469968a`
- Corpus import commit: `ac47da46e9e864191e742f8717d54537a746a55a`
- Project-chair import decision: `docs/decisions/0008-protocol-v4-sealed-corpus-import.md`
- Reviewed evaluator payload: `8319a3e776aaad245cc69e78d04d4df93ef625bd`
- Independent corpus review: `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-review.md`; verdict `ACCEPT`
- Public validation manifest canonical JSON digest: `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490`
- Proposed validation schedule canonical JSON digest: `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`
- Proposed validation schedule raw bytes: 391,950 bytes; `sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec`

The review-assignment commit is the commit containing this file. Record its full identity before beginning substantive review.

## Allowed material and custody boundary

Use a new detached worktree at candidate-report commit `bb473e22...`. The worktree must be clean before and after every check. Materialize with LF line endings and do not modify the candidate or report.

You may inspect:

- the candidate and report commits;
- the six candidate paths;
- the imported public validation cases, fixtures, registries, and manifests;
- the frozen protocol, runner, grader, analysis, worlds, Arm B, and pre-validation artifact manifest;
- the public corpus handoff, independent review, and import decision; and
- historical public authoring-scheduler code needed to compare the deterministic construction.

Do not open, enumerate, hash, mount, or inspect any private evaluator root or any validation or held-out label. Do not inspect results directories for a prospective sealed-tranche run. Do not invoke a model, runtime, Phoenix arm, experiment runner, grader on model output, schedule execution, validation analysis, or held-out analysis.

The schedule tool itself is outcome-free and may be run after static inspection. Independently reconstruct the expected schedule without relying on the candidate tool or its tests as the sole source of truth.

## Required checks

### 1. Identity, ancestry, and six-path boundary

- Require the report commit's sole parent to be candidate `44fccf51...`, and the candidate's sole parent to be `2822b1e8...`.
- Require candidate commit `44fccf51...` to change exactly the six paths listed in the candidate report and no other path.
- Recompute every raw artifact length and SHA-256 in the candidate report.
- Require the report commit to change only the candidate report and require its raw identity to match the pin above.
- Confirm no candidate or report path is a private label, result, runtime stream, trial, grade, analysis output, gate change, or frozen runner/grader change.

### 2. Exact public source and case-manifest linkage

- Independently recompute the validation manifest's canonical JSON digest and require the pinned value.
- Require the manifest to be sealed version 1 for `validation`, with 120 cases in 24 families and five cases per family.
- Require the 120 scheduled case IDs and family IDs to equal the manifest set exactly: no missing, extra, or duplicate row.
- Independently recompute canonical JSON digest for every public validation case and require equality with its manifest `input_digest`.
- Require each case's `case_id`, `family_id`, and `sandbox_fixture` to equal its manifest row.
- Confirm the source bytes at candidate parent `2822b1e8...` are the exact imported public payload bytes accepted by the independent corpus review. Do not infer source identity only from matching filenames.

### 3. Generator safety and frozen-algorithm equivalence

- Statically inspect every schedule-tool source and test file.
- Require the command to read only the public validation manifest and public validation case files. It must not read a label, fixture content, private root, result, runtime output, trial, grade, or outcome.
- Require no network access, subprocess, model, runtime, runner, grader, or analysis invocation.
- Require output paths to stay inside the repository and existing outputs not to be overwritten.
- Compare the deterministic algorithm with the independently accepted authoring scheduler: seed `20260817`, three repetitions, A–E arm order, sorted inputs, SplitMix64 stream and rejection sampling, Fisher-Yates shuffles, family blocking, `(case_id, repetition)` pairing keys, and the ten-row odd-treatment Williams design. The only semantic changes allowed are public validation inputs and `tranche: validation`.
- Run the candidate tests and vet, but do not treat self-tests as independent proof.

### 4. Independent schedule reconstruction

Build an independent, disposable reconstruction from the public manifest and cases. Do not copy or import the candidate generator. The disposable reconstruction must remain outside the Phoenix checkout and must not be committed.

Require exact deep equality with `experiments/frontier-v1/schedules/validation.json`, including entry order. Independently verify:

- version 1 and tranche `validation`;
- seed `20260817`, repetitions 3, and arms exactly `A,B,C,D,E`;
- 120 cases, 24 families, five cases per family;
- 360 pairing keys and 1,800 launches;
- launch indexes exactly `0..1799`;
- pairing indexes exactly `0..359`;
- each pairing key has one case, one repetition, one family, and all five arms consecutively;
- each case appears once per arm in each repetition;
- each family occupies exactly one contiguous 75-entry block;
- every five-arm row belongs to the correct Williams row selected by local pairing position and family block; and
- all ten Williams rows are permutations and every directed predecessor pair appears exactly twice across the ten-row design.

Recompute both the raw-byte identity and canonical JSON digest of the candidate schedule and require the pinned values. Run the candidate verifier only after the independent reconstruction succeeds, and require it to report the same schedule and manifest digests.

### 5. Reproducible public checks

From the detached review worktree, run:

```powershell
go test -count=1 ./experiments/frontier-v1/scheduletool
go vet ./experiments/frontier-v1/scheduletool
go run ./experiments/frontier-v1/scheduletool --repo-root . --verify experiments/frontier-v1/schedules/validation.json
Set-Location experiments/frontier-v1/corpusctl
go test -count=1 ./...
go vet ./...
Set-Location ../../..
make quality
git diff --check 2822b1e8d2b3c5675d52eec9fe3435349618bcc0 44fccf51cb1684a1b71da9128cab30ad2b2fb6af
```

All must pass. `make quality` must reproduce the frozen world-build identity. Do not run either sealed-tranche seal with `--write`.

### 6. Protected state and authorization boundary

- Require protocol, runtime and prompts, Arm A schemas, Arm B, worlds, accepted runner files, grader, analysis, world build, imported corpus, registries, manifests, and existing accepted freeze entries to be unchanged by the candidate.
- Require `pre-validation-artifacts.json` to remain `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`.
- Confirm no validation or held-out private label exists in the implementation or review checkout.
- Confirm no schedule was executed and no model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed.
- Confirm the candidate report and README clearly state that the accepted runner is authoring-only and cannot execute the validation schedule.
- Treat that execution limitation as a separate pre-Gate blocker, not as permission to modify the candidate or silently assume an execution path. `ACCEPT` is allowed only if the schedule candidate itself is correct and the limitation remains explicit.

## Review record

Write the public review record only to:

`D:\Work\personal\phoenix\docs\reviews\2026-08-21-protocol-v4-validation-schedule-review.md`

Return `ACCEPT`, `REVISE`, or `REJECT`, with P0-P3 findings and exact public paths and lines. The record must include:

- reviewer role, independence declaration, date, review-assignment commit, candidate and report identities, worktree, and material seen;
- ancestry, six-path boundary, raw identities, source-manifest identity, 120/24/5 linkage, and exact imported-source result;
- generator safety and frozen-algorithm comparison;
- the independent reconstruction method and exact equality result;
- schedule structure, Williams-design, raw-byte, and canonical-digest results;
- test, vet, quality, protected-state, custody, no-outcome, and clean-worktree results;
- findings and accepted limitations;
- the authoring-only execution blocker; and
- the smallest next artifact allowed by the verdict.

## Verdict bar

- `ACCEPT` requires every identity, source, custody, construction, ordering, structural, reproducibility, and protected-state check to pass, with no unresolved P0 or P1. Every P2 must be fixed or specifically justified as nonblocking.
- Any observed sealed outcome, opened private label, source-corpus mismatch, candidate/report ancestry mismatch, schedule-byte mismatch, non-independent reviewer, or generated schedule informed by outcome data is P0 and requires `REJECT`.
- Use `REVISE` for a correctable generator, schedule, evidence-record, documentation, or review-record defect when source identity, custody, and independence remain intact.
- If `ACCEPT`, authorize only a focused project-chair freeze decision for this exact schedule digest. That later patch may add an accepted `validation_schedule` entry to `pre-validation-artifacts.json`, set `status: complete`, set `remaining: []`, and keep both gates false. It must not modify any frozen implementation byte or authorize execution.
- Even after schedule freeze, Gate 1A must remain closed until the authoring-only execution boundary is resolved by a separate independently reviewed artifact and explicit project-chair decision. Held-out remains closed.
