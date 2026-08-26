# 0021: Protocol-v4 Gate 1A validation world-reference compatibility import and refreeze

- **Status:** Compatibility repair imported and frozen; validation closed; held-out closed
- **Date:** 2026-08-26
- **Decision:** `IMPORT_AND_FREEZE`
- **Closure:** `docs/decisions/0020-protocol-v4-gate-1a-validation-input-compatibility-closure.md` at `acc9730`
- **Accepted payload:** `215b30ae89933c532468b452be6238a6028a740e`
- **Candidate metadata:** `experiments/frontier-v1/artifacts/gate-1a-validation-world-compatibility-candidate.json` at `a9f4e9a`
- **Independent review:** `docs/reviews/2026-08-26-protocol-v4-validation-world-compatibility-review.md` at `a71e59b457cdfdb2014b434033946d319036c060`; verdict `ACCEPT`

## Decision

The project chair accepts the independent review and refreezes the validation public-input and execution boundary identities repaired after decision 0020. Both outcome gates remain closed throughout this import.

The repair migrates exactly 120 public validation case `world_ref` values from the former world identity to the production world already frozen by decisions 0014 and 0018. The repository's accepted seal command regenerates the case input digests and manifest. It changes no case ID, goal, family, fixture, class, label digest, grader rule, schedule entry, prompt, runtime, arm behavior, world byte, world-build input, analysis byte, held-out input, or outcome.

The replacement public identities are:

| Artifact | Identity |
| --- | --- |
| Production world, canonical JSON | `sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc` |
| Validation case raw index | `sha256:16ef12548affa352dd8919bbb498cf6e88173e6f3c793968df1807bd48e6e8d4` |
| Validation manifest, canonical JSON | `sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6` |
| Validation manifest, raw bytes | `sha256:6f18ce9e8f7954d585647a6f8949556cd58c053e5ebb6ff098b9744242fce7d9` |
| Validation schedule, canonical JSON | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Public label registry, raw bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Grader semantics | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| Private archive index | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` |

The replacement validation execution boundary contains 21 files with aggregate identity `sha256:57e071790f2e49e7a5ab5c45a2685578ccaf0f9bd4fbfea8712f7d172c9e6424`. The added compatibility test preflights all 120 validation cases against the production world and includes a negative control for the former identity. Its gate-neutral form remains valid after a later legitimate reopening.

## External custodian

The matching external adapter remains outside the Phoenix repository and private labels remain outside the workspace. Its reviewed identities are:

- source `sha256:4e63984830b1621a909ec8479401942542c4bfb90576e43bace92d266a533ef6`;
- launcher `sha256:492ada05475aac47a809864f4d1379d48f7df17efd27aec9a3b944cef89e0237`;
- Linux/amd64 binary `sha256:a4f0e8c0f3356ef2dff28e09319e0e68a5a6398d9eed45f9814527070eb5012f`; and
- Windows/amd64 binary `sha256:0b8b228bd1e92518ff29bdab3411b3385a5706a2aa2c8dadf538b6ba286a0598`.

Both adapters return the replacement manifest identity and every unchanged frozen custody identity. They accept only `scheduled-validation-2` evidence and reject the burned archive.

## Accepted residual

The external custodian directory is not version-controlled. Prior-source provenance depends on recorded hashes, preserved binaries, exact one-constant reconstruction, and behavioral checks. The carried residuals from decision 0018 remain in force.

## Refreeze boundary

This focused refreeze patch changes:

- `experiments/frontier-v1/pre-validation-artifacts.json`, importing the accepted 21-file boundary and replacement validation-manifest identity while retaining both closed gates;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, enforcing the replacement provenance, identities, file count, and closed gate; and
- this decision record.

The accepted payload and review records are already committed. No model, trial, grade, output directory, cost, validation result, or held-out result is produced or observed by this decision.

## Authorization boundary

This decision does not reopen Gate 1A. A later chair decision may authorize one new disjoint 1,800-launch execution only after the full suite and quality gate pass against this refreeze. Held-out remains closed.
