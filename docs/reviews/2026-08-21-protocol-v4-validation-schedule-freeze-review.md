# Protocol-v4 validation-schedule freeze external review

- **Reviewer role:** `validation_schedule_freeze.external_reviewer`. Distinct from `validation_schedule.candidate_author`, `validation_schedule.independent_reviewer`, the project-chair freeze session that wrote decision 0009, the project-chair corpus import session, and every third-revision corpus remediation, packaging, handoff, and review role.
- **Date:** 2026-08-21
- **Verdict:** `ACCEPT`
- **Reviewed decision:** `docs/decisions/0009-protocol-v4-validation-schedule-freeze.md` at `f79f077481633bd75e91399af9e8f192d883bb5b` (`docs: record validation schedule freeze decision`; sole path this file’s parent commit)
- **Decision raw bytes:** 4,624; `sha256:f00a1280d2ee8d86ccc85b2cccb83edfe5d2442e4facd0e2c09473246e25d2bc`
- **Focused freeze commit:** `dd6fef3b09055c1f8d79234f6ce4c1c93ad13912` (`experiments: freeze validation schedule digest`; sole parent `a301bdc781bd4b93f12418b577b7446d01c1b7e4`)
- **Independent schedule review:** `docs/reviews/2026-08-21-protocol-v4-validation-schedule-review.md` at `a301bdc781bd4b93f12418b577b7446d01c1b7e4`; verdict `ACCEPT`
- **Schedule candidate:** `44fccf51cb1684a1b71da9128cab30ad2b2fb6af`
- **Candidate-report commit:** `bb473e22b1e2a5b21221a0b7b57e2aec4836c52b`
- **Review checkout:** implementation workspace `D:\Work\personal\phoenix` at `HEAD` `f79f077481633bd75e91399af9e8f192d883bb5b`; `git status --porcelain` empty of tracked changes before and after every check

This review does not authorize schedule execution, a model or arm run, a prospective grade, Gate 1A, validation, held-out evaluation, outcome analysis, access to a private label, or any change to frozen implementation bytes. It decides only whether the recorded freeze matches the independent `ACCEPT` authorization.

## Independence declaration

This session did not author the candidate generator, schedule bytes, candidate report, review assignment, independent schedule review, freeze patch, or decision 0009. It did not perform the corpus import. It did not implement Phoenix.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Other sessions prepared the candidate, the independent schedule review, and the freeze. Those sessions did not write this record.
- The review ran in the implementation workspace rather than a new detached worktree because the reviewed artifact is the already-committed freeze decision at `HEAD`. Tracked status remained empty of modifications.
- This review re-verified freeze identity, path boundary, gates, protected bytes, and public checks. It did not repeat the independent full reconstruction of the schedule generator; it required the frozen schedule bytes to be identical to the independently reviewed candidate blob.

No private evaluator root was opened, hashed, mounted, or enumerated. No validation or held-out label was opened. Authoring labels were not opened. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. `seal --write` was not run.

## Verdict

The freeze is the exact patch the independent schedule review authorized. Decision 0009 records the correct identities, the chair/candidate-author overlap, closed gates, and the authoring-only execution blocker. There are no P0, P1, P2, or unresolved P3 findings.

`ACCEPT` confirms that schedule digest `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` is frozen. It does not open Gate 1A or establish that the accepted runner can execute the schedule.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Decision commit tree vs freeze `dd6fef3…` | `git diff-tree --name-only -r f79f077…` | exactly `docs/decisions/0009-protocol-v4-validation-schedule-freeze.md` |
| Decision 0009 raw bytes at `f79f077…` | SHA-256 of git blob; 4,624 bytes; LF-only | `sha256:f00a1280d2ee8d86ccc85b2cccb83edfe5d2442e4facd0e2c09473246e25d2bc` |
| Freeze commit tree vs review `a301bdc…` | `git diff-tree --name-status -r dd6fef3…` | exactly the three claimed freeze paths |
| Review record at `a301bdc…` | SHA-256 of git blob; 15,250 bytes; LF-only | `sha256:cd557f3dc79b7a2bbc2acbba0c40969b832178b1c23b4ca4c4442d602c601ca0` |
| Schedule raw bytes at candidate and `HEAD` | SHA-256 of git blob; 391,950 bytes; blobs equal | `sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec` |
| Schedule canonical JSON | Frozen `corpusctl digest`; candidate `--verify` agreed | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Validation manifest canonical JSON | Frozen `corpusctl digest` | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| World-build digest from `make quality` | `cmd/build-manifest` then `quality-check compare` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: `HEAD` `f79f077…`; sole parent freeze `dd6fef3…`; that commit’s sole parent independent review `a301bdc…`; that commit’s sole parent assignment `3972915…`. Candidate `44fccf51…` is an ancestor of the freeze.

`experiments/frontier-v1/pre-validation-artifacts.json` is `status: complete`, `remaining: []`, `may_open_validation: false`, and `may_open_held_out: false`. The `validation_schedule` entry is `frozen: true` and pins the identities above.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Independent-review predicate | Pass | Review at `a301bdc…` is `ACCEPT` with no P0–P3 findings. Freeze parent is that commit. Review blob is unchanged through `HEAD`. |
| 2. Authorized freeze boundary | Pass | `dd6fef3…` changes exactly `pre-validation-artifacts.json`, `runner/artifact_freeze_test.go`, and `runner/local_artifact_candidate_test.go`. The two test files are freeze-guard tests, not members of the frozen nine-file runner set. Decision `f79f077…` adds only 0009. |
| 3. Exact reviewed schedule identity | Pass | `HEAD` schedule blob equals candidate `44fccf51…`. Raw length and SHA-256 match the pins. Canonical digest matches `corpusctl digest` and `--verify`. Structure: version 1, tranche `validation`, seed `20260817`, 1,800 contiguous launches, 120 cases, 24 families, five cases per family, 360 pairing keys, each pairing emitting all five arms consecutively, each family one contiguous 75-entry block. |
| 4. Manifest completion without gate opening | Pass | Status `complete`, `remaining: []`, both gates false, `source_limit: authoring_only`. Protocol `gate.may_open_validation` and `may_open_held_out` remain false and were not modified. |
| 5. Protected frozen implementation bytes | Pass | `git diff --name-only a301bdc… dd6fef3…` is only the three freeze paths. Protocol, prompts, Arm A derivation sources, Arm B, worlds, the nine-file runner set, grader, analysis, world-build, imported corpus, registries, manifests, and the schedule file itself are unchanged by the freeze. |
| 6. Decision-record accuracy | Pass | 0009’s pinned commits, digests, launch counts, three-path freeze boundary, role-overlap disclosure, and authoring-only blocker match the repository. |
| 7. Authoring-only execution boundary | Pass | Frozen `loadCase` still rejects non-`authoring_` IDs (`runner/materialize.go` 21–22). `validateSchedule` still requires tranche `authoring` (`runner/schedule.go` 83–84). 0009 states that completing the artifact list does not create an execution path. |
| 8. Custody and no-outcome | Pass | No `labels/validation` or `labels/held_out` directory exists. `labels/.gitignore` continues to exclude those names. No validation or held-out results directory exists. Authoring labels remain; they were not opened. Freeze and decision commits add no trial, grade, or outcome path. |
| 9. Reproducible public checks | Pass | Freeze-guard tests passed. `scheduletool --verify` printed the pinned schedule and manifest digests. `make quality` exited 0 and reproduced world-build `sha256:27c2f537…0188d`. `git diff --check a301bdc… HEAD` is empty. Gitignored `bin/` and `build/` from `make quality` were removed afterward. |

## Findings

No P0, P1, P2, or P3 findings.

## Accepted limitations

- Same-machine fresh session. That is recorded above and does not reuse a prohibited role.
- The chair/candidate-author overlap claimed by 0009 is accepted as a recorded process fact. Git timestamps are compatible (candidate `11:35`, independent review commit `a301bdc…`, freeze `12:13`, decision `12:14`). This review did not re-bind those commits to agent-session identifiers; the material control is that the freeze parent is the independent `ACCEPT` record, not a self-test.
- This review did not re-run the disposable independent reconstruction. Byte identity with the already-accepted candidate blob is the freeze predicate.
- Frozen `protocol.json` `remaining_execution_blockers` still lists historical items, including `every artifact_freeze.before_validation digest`. That list is frozen protocol text and was correctly left unchanged. Live freeze remainder is `pre-validation-artifacts.json` `remaining: []`. The live pre-Gate blocker is the authoring-only execution boundary.

## Authoring-only execution blocker

The accepted runner remains authoring-only. Schedule freeze does not create a validation execution path and does not open Gate 1A. Held-out remains closed.

## Smallest next artifact

A separately specified, independently reviewed custody-safe validation execution boundary that preserves this frozen schedule and A–E trial contract, consumes the public validation cases, obtains private grading without copying labels into the implementation workspace, retains outcome evidence, and keeps held-out closed. Any change to frozen runner or grader bytes requires an independent replacement review and refreeze before a project-chair Gate 1A opening decision.

This `ACCEPT` does not authorize that next artifact.
