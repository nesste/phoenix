# 0009: Protocol-v4 validation-schedule freeze

- **Status:** Validation schedule frozen; outcome gates remain closed
- **Date:** 2026-08-21
- **Decision:** `FREEZE`
- **Independent review:** `docs/reviews/2026-08-21-protocol-v4-validation-schedule-review.md` at `a301bdc781bd4b93f12418b577b7446d01c1b7e4`; verdict `ACCEPT`
- **Schedule candidate:** `44fccf51cb1684a1b71da9128cab30ad2b2fb6af`
- **Candidate report:** `bb473e22b1e2a5b21221a0b7b57e2aec4836c52b`
- **Focused freeze commit:** `dd6fef3b09055c1f8d79234f6ce4c1c93ad13912`

## Decision

The project chair accepts the independent validation-schedule review and freezes the exact reviewed schedule in `experiments/frontier-v1/pre-validation-artifacts.json`.

The accepted schedule is `experiments/frontier-v1/schedules/validation.json`, with canonical JSON identity `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` and raw identity `sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec` over 391,950 bytes. It contains 1,800 launches: 120 validation cases in 24 five-case families, three repetitions, and five consecutive A–E assignments for each of 360 pairing keys.

The pre-validation artifact manifest is now `status: complete` with `remaining: []`. `may_open_validation` remains `false`, and `may_open_held_out` remains `false`.

This decision freezes an outcome-free schedule identity only. It does not execute the schedule, open Gate 1A, authorize a model or arm run, grade a prospective trial, open validation or held-out, or authorize outcome analysis.

## Chair basis and recorded limitation

The independent reviewer returned `ACCEPT` with no P0, P1, P2, or unresolved P3 findings. The reviewer independently verified ancestry, the six-path candidate boundary, exact imported-source bytes, case-manifest linkage, generator safety, equivalence to the frozen authoring scheduling algorithm, an independent full reconstruction, Williams ordering, raw and canonical identities, tests, protected state, custody, and absence of outcomes.

This project-chair session also authored the outcome-free schedule candidate. It did not perform the independent review. The role overlap is recorded: the freeze decision relies on the fresh independent review at `a301bdc...`, not on the candidate author's self-tests. No private label or prospective outcome was inspected by this chair session.

## Focused freeze boundary

Freeze commit `dd6fef3...` changes exactly:

- `experiments/frontier-v1/pre-validation-artifacts.json`, adding the accepted `validation_schedule` entry, setting `status: complete`, and emptying `remaining`;
- `experiments/frontier-v1/runner/artifact_freeze_test.go`, extending the non-frozen guard to verify the exact schedule, source-manifest, review, and completed-state identities; and
- `experiments/frontier-v1/runner/local_artifact_candidate_test.go`, updating the non-frozen gate-state guard for the completed artifact set.

No frozen implementation byte changed. Protocol, exact prompts, Arm A schemas, Arm B, worlds, the accepted nine-file runner set, grader, analysis, world build, public corpus, registries, manifests, and both gate booleans remain unchanged.

## Verification

After the focused freeze:

- the schedule verifier reproduced canonical schedule digest `sha256:b38a0eaa...5813` and source-manifest digest `sha256:57ccc0c7...d490`;
- the artifact manifest reports `complete`, zero remaining items, and both gates false;
- the updated freeze guards pass;
- `make quality` passes, including root tests, vet, staticcheck, module verification, vulnerability scan, complexity and duplication checks, schema validation, authoring corpus validation, and reproducible world-build comparison;
- `git diff --check` is clean; and
- no schedule, model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed.

## Remaining execution blocker

The accepted runner remains authoring-only. It rejects validation case IDs and a schedule whose tranche is `validation`. Completing the artifact-freeze list does not create an execution path and does not open Gate 1A.

The next artifact must separately define and implement a custody-safe validation execution boundary. It must preserve the frozen schedule and A–E trial contract, consume the public validation cases, obtain private grading without copying labels into the implementation workspace, retain outcome evidence, and keep held-out closed. Any change to frozen runner or grader bytes requires an independent replacement review and refreeze before a project-chair Gate 1A opening decision.
