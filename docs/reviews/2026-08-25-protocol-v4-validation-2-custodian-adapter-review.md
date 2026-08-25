# Protocol-v4 validation-2 custodian adapter review

- **Date:** 2026-08-25
- **Verdict:** `ACCEPT`
- **Scope:** external validation custodian adapter for `scheduled-validation-2`
- **Review mode:** independent read-only sub-agent review
- **Execution state:** no model, trial, private grade, output directory, or validation execution was produced

## Reviewed change

The external custodian adapter now confines trial evidence to:

```text
D:\Work\personal\phoenix\experiments\frontier-v1\results\scheduled-validation-2
```

The preserved adapter for the interrupted execution confined evidence to the sibling `scheduled-validation` archive. The new launcher passes the required `--output-dir` and selects the replacement Linux/amd64 grader. The grading path, private-label verification, frozen identity checks, `corpusctl` invocation, strict JSON decoding, and grade consistency checks remain unchanged in the reviewed source.

The custodian source and binaries remain outside the Phoenix implementation repository. Private label contents were not inspected during this review.

## Reviewed artifacts

| Artifact | SHA-256 |
| --- | --- |
| `D:\Work\personal\phoenix-validation-custodian\main.go` | `31f6efdafc3d1fa446f5524265e13b26cf9579cec74da34b5bd7ae5bc49f560b` |
| `D:\Work\personal\phoenix-validation-custodian\main_test.go` | `7021e39e62cc5cc3dedc2617ac315b8c7059d37e38a0f7d238a712e356324c6d` |
| `D:\Work\personal\phoenix-validation-custodian\run-validation.sh` | `ba559912c5a2d4043d0bfc595172dc304bd7d4c707dafc3f687dbaf2ec67b6d4` |
| Replacement Linux/amd64 grader | `ebcdb440fe6f0dad900ddd343da625eef0fe066c602955b34d543aa04895c6d0` |
| Replacement Windows/amd64 grader | `02dd1a245665be7bfb4baad935f816936c87ef48305f0d32b727912073fa029b` |
| Preserved original Linux/amd64 grader | `f3220d84e098ec39df4cc16c731ea54f581e451fb509aba16a0d5f24a16efa24` |
| Preserved original Windows/amd64 grader | `9eba89f08dae4b19544ba5017b7fdf6a2c6b25740333b1c13c1e4bbe5a9191c1` |

Replacement artifact paths:

```text
D:\Work\personal\phoenix-validation-custodian\frontier-v1-grader-validation-2-linux-amd64
D:\Work\personal\phoenix-validation-custodian\frontier-v1-grader-validation-2.exe
```

## Frozen custodian description

Both replacement binaries returned the same accepted `describe` document:

| Field | Value |
| --- | --- |
| `v` | `1` |
| `tranche` | `validation` |
| `grader_digest` | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| `schedule_digest` | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| `manifest_digest` | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| `label_registry_raw_sha256` | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| `private_archive_digest` | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` |
| `cases` | `120` |

The read-only handshake succeeded on Windows and Linux. The private archive index matched the frozen digest; no label content was opened.

## Verification

The independent review confirmed:

- `go test` and `go vet` pass for the external custodian source;
- the boundary unit test accepts a regular file under `scheduled-validation-2` and rejects its `scheduled-validation` sibling;
- both replacement binaries reject a retained trial from the burned 720-trial archive with exit code 1 and the generic failure message;
- the Linux artifact is an executable linux/amd64 build and is selected by the launcher;
- `run-validation.sh` passes shell syntax validation and supplies the new output directory, replacement grader, 300 USD run budget, 0.15 USD per-trial cap, and 180-second timeout;
- the original grader binaries remain present with their recorded hashes; and
- `scheduled-validation-2` remained absent throughout review.

## Findings

No P0 or P1 finding blocks use of the replacement adapter.

One P2 auditability limitation remains: the external custodian directory is not version-controlled and no preserved copy of the original source exists. An exact historical source diff therefore cannot be reconstructed independently. Current-source inspection, preserved binary hashes, unchanged frozen `corpusctl` identity, exact `describe` results, unit tests, and direct old-archive rejection probes support the conclusion that the evidence-directory boundary is the only grading-code change.

## Authorization boundary

This record documents the adapter review. It changes no frozen execution-boundary file and does not itself launch or authorize an additional validation execution. Gate state and the single execution authorized by decision 0019 remain unchanged.
