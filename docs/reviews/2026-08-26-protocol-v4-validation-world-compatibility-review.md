# Protocol-v4 validation world-reference compatibility review

- **Date:** 2026-08-26
- **Verdict:** `ACCEPT`
- **Base closure:** `acc9730dfa553e56dc20af992176717a9c4f6044`
- **Reviewed payload:** `215b30ae89933c532468b452be6238a6028a740e`
- **Candidate metadata:** `a9f4e9a`
- **Review mode:** independent read-only sub-agent review, followed by a focused amendment review
- **Outcome access:** none

## Scope

The review covered the Gate 1A compatibility repair required by decision 0020:

- 120 public validation case files;
- regenerated `experiments/frontier-v1/manifests/validation.json`;
- the 21-file validation execution boundary;
- schedule and manifest identity handling;
- the positive and negative production-world compatibility tests; and
- the external validation-2 compatibility custodian source, launcher, and Windows/Linux binaries under `D:\Work\personal\phoenix-validation-custodian`.

Both outcome gates remained closed. The reviewer did not inspect private label contents, validation outcomes, or the archived 720-trial results. No model, trial, grade, output directory, cost, or outcome was produced.

## Verdict

`ACCEPT`. No P0 or P1 finding blocks import or refreeze.

The focused amendment review also returned `ACCEPT` and confirmed that the overall verdict remains valid. Amendment `215b30a` changes only `validation_world_compatibility_test.go`: it keeps the 120-case positive integration test valid after a legitimate gate reopening and adds an explicit negative control using the former world identity. It changes no production code.

## Public validation inputs

Exactly 120 validation case files changed. Each file differs from the closure base only in `world_ref`:

```text
sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b
sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc
```

Case IDs, goals, family IDs, fixtures, classes, and all other case fields are unchanged. The raw validation-case index is `sha256:16ef12548affa352dd8919bbb498cf6e88173e6f3c793968df1807bd48e6e8d4` under the candidate's ordinal-path index algorithm.

The accepted `corpusctl seal --write` command reproduced the replacement manifest byte-for-byte:

| Identity | SHA-256 |
| --- | --- |
| Manifest canonical JSON | `39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6` |
| Manifest raw bytes | `6f18ce9e8f7954d585647a6f8949556cd58c053e5ebb6ff098b9744242fce7d9` |

Only the top-level world reference and each case's world reference and derived input digest changed. Counts, class coverage, family coverage, grading metadata, fixtures, and label digests are unchanged.

## Preserved identities

The reviewer independently reproduced or compared these unchanged identities:

| Artifact | Identity |
| --- | --- |
| Validation schedule, 1,800 launches | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Public label registry, raw | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Grader semantics | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |
| Private archive index | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` |
| Production world, canonical JSON | `sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc` |

The production world, world build, prompts, runtime, analysis, schedule construction, held-out corpus, and archived evidence have no payload diff.

## Execution-boundary inventory

All 21 candidate file hashes reproduce at payload `215b30a`. Their aggregate identity is:

```text
sha256:57e071790f2e49e7a5ab5c45a2685578ccaf0f9bd4fbfea8712f7d172c9e6424
```

The authoritative historical-payload test reproduces this identity from Git rather than trusting the working tree.

The compatibility tests establish both directions:

- all 120 validation cases pass production-world preflight; and
- a case carrying the former `f5f6…1615b` reference fails with the expected production-world mismatch.

These tests are independent of gate-open state. Runtime gate enforcement remains unchanged and still executes before validation input preparation.

## External custodian

The reviewed compatibility adapter identities are:

| Artifact | SHA-256 |
| --- | --- |
| `main.go` | `4e63984830b1621a909ec8479401942542c4bfb90576e43bace92d266a533ef6` |
| `main_test.go` | `7021e39e62cc5cc3dedc2617ac315b8c7059d37e38a0f7d238a712e356324c6d` |
| `run-validation.sh` | `492ada05475aac47a809864f4d1379d48f7df17efd27aec9a3b944cef89e0237` |
| Linux/amd64 binary | `a4f0e8c0f3356ef2dff28e09319e0e68a5a6398d9eed45f9814527070eb5012f` |
| Windows/amd64 binary | `0b8b228bd1e92518ff29bdab3411b3385a5706a2aa2c8dadf538b6ba286a0598` |

Replacing the one new manifest constant in-memory reproduces the previously reviewed source hash. Removing only the `-compat` binary suffix from the launcher reproduces its previous reviewed hash. Both binaries return the replacement manifest identity and every unchanged frozen identity through `describe`, and both reject a trial from the burned `scheduled-validation` archive.

## Verification

The independent reviewer reproduced:

- `corpusctl` tests, vet, and byte-stable sealing;
- schedule-tool tests, vet, and deterministic schedule verification;
- runner focused tests, vet, candidate inventory, compatibility checks, and world-build reproduction;
- custodian tests, vet, shell syntax, hashes, handshakes, and old-archive rejection; and
- the expected pre-refreeze runner state, where only the authoritative accepted-freeze test remains stale.

`scheduled-validation-2` remained absent throughout review.

## Accepted residual

P2: the external custodian directory is not version-controlled. Prior-source provenance therefore relies on its recorded hashes. The exact one-constant source reconstruction, preserved earlier binaries, new artifact hashes, and behavioral probes materially reduce this auditability risk but do not remove it.

## Authorization boundary

This `ACCEPT` authorizes only a later import/refreeze decision. It does not reopen Gate 1A, authorize held-out access, start validation, invoke a model, or produce a grade.
