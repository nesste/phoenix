# 0008: Protocol-v4 sealed-corpus import

- **Status:** Exact reviewed payload imported; outcome gates remain closed
- **Date:** 2026-08-21
- **Decision:** `IMPORT`
- **Independent review:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-review.md` at `20e8c4a0d8feb7fef732c2efd593d1675a717dc7`; verdict `ACCEPT`
- **Evaluator base:** `c852101e8d7cb52e4569bf3866de54a0ce648b44`
- **Reviewed payload:** `8319a3e776aaad245cc69e78d04d4df93ef625bd`
- **Report-only commit, not imported:** `0f285e08f492db8fd0f5b13b88557648c0de2f55`
- **Implementation import commit:** `ac47da46e9e864191e742f8717d54537a746a55a`

## Decision

The project chair accepts the independent third-revision review and imports only the exact reviewed evaluator-base-to-payload patch.

The imported raw diff is 1,188,302 bytes with identity `sha256:b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea`. The import commit changes exactly the reviewed 884 public paths: 480 corpus paths, 400 fixture paths, and four label-digest registry or sealed-manifest paths. The evaluator's report-only commit is not part of the import.

This decision does not import private labels, open an outcome, generate a schedule, open Gate 1A, authorize a model or arm run, or authorize validation, held-out evaluation, grading, or outcome analysis.

## Chair basis and independence

The chair relied on the public independent review record and its `ACCEPT` verdict. That review reports no P0, P1, P2, or unresolved P3 findings and independently verified identity, ancestry, patch bytes, the one-label remediation boundary, custody, registry and manifest linkage, executable reproducibility, opacity, protected state, and the disputed cascade-terminal adjudication.

This chair session did not author the corpus revision, corrected label, audit, adjudication, private archive, evaluator report, review assignment, or independent review. It did not inspect validation or held-out labels or observe a prospective result. It inspected the public patch bytes, public review record, public workspace state, and reproducible public checks needed for this import decision.

The chair accepts the review's recorded limitations: unchanged public case and fixture bytes were allowed to inherit the second-revision pairwise screen; 23 unchanged cascade proofs were not re-executed by the third-revision reviewer; and deep polarity probes concentrated on the corrected temptation label and disputed cascade row.

## Exact import verification

Before import:

- the stored patch was independently hashed as 1,188,302 bytes and `sha256:b607b25bfeee850fd4d6827b943e076af5d2382836e9393ac84b4d169b9f78ea`;
- all 884 approved paths were unchanged in the implementation workspace relative to evaluator base `c852101e...`;
- `git apply --check --index` passed; and
- validation and held-out private-label directories were absent.

Import commit `ac47da46...` has sole parent `20e8c4a0...`. A binary-safe capture of native stdout from `git --no-pager diff --binary 20e8c4a0... ac47da46...` is byte-identical to the stored reviewed patch: 1,188,302 bytes with the same SHA-256. The commit changes 884 paths and no decision, report, private, implementation, protocol, world, runtime, prompt, arm, runner, grader, analysis, freeze, gate, schedule, trial, or outcome path.

After import:

- the eight approved path groups are byte-identical to reviewed payload `8319a3e...`;
- protocol, runtime and prompt inputs, Arm A schema derivation, Arm B, worlds, runner, grader, analysis, world-build and accepted pre-validation freeze artifacts remain byte-identical to evaluator base `c852101e...` where required;
- no validation or held-out private label is present in the implementation workspace;
- `experiments/frontier-v1/pre-validation-artifacts.json` remains `status: partial` with exactly `remaining: ["schedule digest"]`;
- `may_open_validation` remains `false`;
- `may_open_held_out` remains `false`; and
- `git diff --check` is clean.

## Reproducible checks

All checks ran in the implementation workspace after applying the exact patch and before committing it:

| Check | Result |
| --- | --- |
| `go test ./...` at repository root | Pass |
| `go test -count=1 ./...` in `experiments/frontier-v1/corpusctl` | Pass |
| `go vet ./...` in `experiments/frontier-v1/corpusctl` | Pass |
| `corpusctl grader-digest` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| Validation dry seal vs committed manifest | Byte-equal; 89,645 bytes; `sha256:0c33fcd516c7b87bf9730b23200e109812a364f5e6d1cf5e7a15ec33583e35ee` |
| Held-out dry seal vs committed manifest | Byte-equal; 88,631 bytes; `sha256:f053ae16cba42a0965716bb60aaf97235bb0db0a6db76eaba87a88922e95c65a` |

No seal used `--write`. No schedule, trial, model call, prospective grade, validation result, held-out result, or outcome was produced or observed.

## Authorization boundary

The public corpus import is complete. The only declared pre-validation remainder is the schedule digest, but this decision does not authorize generating it. Both outcome gates remain closed until a separate, explicit artifact and gate decision satisfies the frozen protocol.
