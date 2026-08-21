# Protocol-v4 validation-schedule candidate

- **Author role:** `validation_schedule.candidate_author`
- **Date:** 2026-08-21
- **Candidate commit:** `44fccf51cb1684a1b71da9128cab30ad2b2fb6af` (`experiments: prepare validation schedule candidate`)
- **Candidate parent:** `2822b1e8d2b3c5675d52eec9fe3435349618bcc0` (project-chair corpus import decision)
- **Imported corpus commit:** `ac47da46e9e864191e742f8717d54537a746a55a`
- **Reviewed evaluator payload:** `8319a3e776aaad245cc69e78d04d4df93ef625bd`
- **Independent corpus review:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-review.md`; verdict `ACCEPT`

This is an outcome-free schedule candidate. It does not freeze the schedule, complete the pre-validation manifest, open Gate 1A, authorize validation, run a model or arm, grade a prospective trial, inspect a private label, or observe an outcome.

## Candidate boundary

The candidate commit changes exactly six paths:

- `experiments/frontier-v1/README.md`
- `experiments/frontier-v1/schedules/validation.json`
- `experiments/frontier-v1/scheduletool/main.go`
- `experiments/frontier-v1/scheduletool/main_test.go`
- `experiments/frontier-v1/scheduletool/schedule.go`
- `experiments/frontier-v1/scheduletool/schedule_test.go`

The accepted scheduled runner remains byte-identical to its frozen artifact set. The generator is a separate command: it cannot execute a trial and is not on the model-run path. Protocol, prompts, Arm A schemas, Arm B, worlds, runner, grader, analysis, world build, accepted freeze entries, corpus bytes, manifests, and both gates are unchanged.

## Source identity and custody boundary

The generator reads only:

- `experiments/frontier-v1/manifests/validation.json`; and
- the 120 public files under `experiments/frontier-v1/corpus/validation/` named by that manifest.

It does not read `labels`, fixtures, results, private roots, runtime output, trial evidence, or grade evidence. It opens no network connection and starts no subprocess. It verifies the sealed version-1 validation manifest, 120-case and 24-family counts, exactly five cases per family, case IDs and family IDs, every public case's canonical input digest, and its fixture reference against the manifest row before scheduling.

The imported validation manifest canonical JSON digest is `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490`.

No validation or held-out private label is present in the implementation workspace. No private root was opened to generate or verify this candidate.

## Deterministic construction

The generator preserves the accepted authoring scheduler's construction except that the output tranche is `validation` and its inputs are the public validation cases:

- seed `20260817`;
- three repetitions;
- arms `A`, `B`, `C`, `D`, `E`;
- cases sorted within families and families sorted before deterministic SplitMix64 shuffling;
- each family emitted as one contiguous block;
- deterministic shuffling of `(case_id, repetition)` pairing keys within each family;
- the accepted ten-row Williams design for five treatments; and
- all five arms emitted consecutively for each pairing key.

Candidate counts:

| Quantity | Result |
| --- | ---: |
| Public cases | 120 |
| Families | 24 |
| Cases per family | 5 |
| Repetitions | 3 |
| Arms | 5 |
| Pairing keys | 360 |
| Launches | 1,800 |
| Entries per family block | 75 |
| Entries per case | 15 |

An independent structural pass, separate from the generator's verifier, found zero errors in launch-index continuity, family-block uniqueness and size, pairing-key contiguity, case/repetition coverage, and five-arm completeness.

## Artifact identities

| Artifact | Raw bytes | Raw SHA-256 |
| --- | ---: | --- |
| `experiments/frontier-v1/README.md` | 17,689 | `sha256:c424ef79e44594069a7d0b0dd8d4dc4311ccd41933bf411944b6309262567558` |
| `experiments/frontier-v1/schedules/validation.json` | 391,950 | `sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec` |
| `experiments/frontier-v1/scheduletool/main.go` | 8,138 | `sha256:cd959297da5cf796c6d0d386ba5cab67522b186216999fa81352a43aa5600de4` |
| `experiments/frontier-v1/scheduletool/main_test.go` | 1,154 | `sha256:1e6d5cad12c29b1db2301c0239fc80083d5cd489ce85ecebd22d509b14b21495` |
| `experiments/frontier-v1/scheduletool/schedule.go` | 5,520 | `sha256:37325fa40e645da066e69c9377f0f66f312799d1862861362c6f3ff36789a182` |
| `experiments/frontier-v1/scheduletool/schedule_test.go` | 2,784 | `sha256:b7b27ef528ce3bb51548bcad255a45aee26e0cb41f32ad83650e69bb84c7544d` |

The schedule's canonical JSON digest, which is the proposed freeze identity, is:

`sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`

Generation and verification both printed that digest and the pinned validation-manifest digest.

## Commands and results

All commands ran from the implementation workspace without a model or private label:

```powershell
go run ./experiments/frontier-v1/scheduletool --repo-root . --write experiments/frontier-v1/schedules/validation.json
go run ./experiments/frontier-v1/scheduletool --repo-root . --verify experiments/frontier-v1/schedules/validation.json
go test -count=1 ./experiments/frontier-v1/scheduletool
go vet ./experiments/frontier-v1/scheduletool
Set-Location experiments/frontier-v1/corpusctl
go test -count=1 ./...
Set-Location ../../..
make quality
git diff --check
```

Generation and verification passed with the same schedule and manifest digests. Schedule-tool tests and vet passed. Corpusctl tests passed. `make quality` passed, including root tests, vet, staticcheck, module verification, vulnerability scan, cyclomatic-complexity and duplication checks, schema validation, authoring corpus validation, reproducible Linux build, and byte comparison with the frozen world-build manifest. `git diff --check` was empty.

## Protected state

`experiments/frontier-v1/pre-validation-artifacts.json` remains:

- `status: partial`;
- `remaining: ["schedule digest"]`;
- `may_open_validation: false`; and
- `may_open_held_out: false`.

No schedule freeze entry has been written. No seal used `--write`. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed.

## Known execution boundary

The accepted scheduled runner remains authoring-only and rejects validation cases and a schedule whose tranche is `validation`. This candidate deliberately does not bypass or modify that boundary. Independent acceptance of the candidate may make only the validation schedule digest eligible for a focused freeze patch; it does not establish that the frozen runner can execute the schedule.

Before Gate 1A can open, a separate reviewed execution boundary must explain how the accepted runtime will consume public validation cases, obtain custody-safe private grading, preserve the frozen A–E trial machinery, and retain outcome evidence without modifying the schedule after outcomes exist. Any change to frozen runner or grader bytes requires a replacement review and refreeze of the affected identity.

## Smallest next artifact

Obtain an independent outcome-free review of candidate commit `44fccf51...`, this evidence record, the imported public corpus source, the standalone generator, and the exact schedule bytes. If that review returns `ACCEPT`, the project chair may consider a focused freeze patch that records the accepted schedule identity, sets the artifact manifest to complete with an empty `remaining` list, and keeps both gates false.

That freeze still must not open Gate 1A. The authoring-only execution boundary remains a separate blocker and must be resolved and reviewed before any validation run.
