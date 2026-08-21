# Protocol-v4 validation-schedule independent review

- **Reviewer role:** `validation_schedule.independent_reviewer`. Distinct from `validation_schedule.candidate_author`, the project-chair corpus import session, every third-revision corpus remediation, packaging, handoff, and review role, and every session that wrote the candidate schedule, generator, evidence record, or this assignment.
- **Date:** 2026-08-21
- **Verdict:** `ACCEPT`
- **Review-assignment commit (process identity only):** `397291503eef7fa6656049780d234898d588294c` (`docs: prepare validation schedule review`; sole path `docs/reviews/2026-08-21-protocol-v4-validation-schedule-review-prompt.md`; implementation `HEAD` at review time)
- **Review base and candidate parent:** `2822b1e8d2b3c5675d52eec9fe3435349618bcc0`
- **Candidate commit:** `44fccf51cb1684a1b71da9128cab30ad2b2fb6af`
- **Candidate-report commit:** `bb473e22b1e2a5b21221a0b7b57e2aec4836c52b`
- **Corpus import commit:** `ac47da46e9e864191e742f8717d54537a746a55a`
- **Reviewed evaluator payload:** `8319a3e776aaad245cc69e78d04d4df93ef625bd`
- **Project-chair import decision:** `docs/decisions/0008-protocol-v4-sealed-corpus-import.md`
- **Independent corpus review:** `docs/reviews/2026-08-20-protocol-v4-sealed-corpus-third-revision-review.md`; verdict `ACCEPT`
- **Review checkout:** `D:\Work\personal\phoenix-validation-schedule-review` (new detached worktree at the report commit; did not previously exist; materialized with `core.autocrlf=false` and `core.eol=lf`; `git status --short` empty of tracked changes before and after every check)
- **Disposable reconstruction:** `D:\Work\personal\phoenix-schedule-review-scratch` (outside every Phoenix checkout; not committed)

This review does not authorize a schedule freeze patch, schedule execution, model or arm run, prospective grade, Gate 1A, validation, held-out evaluation, outcome analysis, or access to a private label. It decides only whether the exact public validation schedule is eligible for a later focused digest-freeze patch.

## Independence declaration

This session did not author the candidate generator, schedule bytes, candidate report, review-assignment prompt, corpus import, or import decision. It did not perform the third-revision corpus review. It did not implement Phoenix.

Recorded limitations:

- The reviewer is a fresh Cursor session on the same machine as earlier Phoenix work. Other sessions prepared the candidate, the report, the import, and this assignment. Those sessions did not write this record.
- The implementation workspace was not used as the review checkout. The public review record is written only to this file.
- The independent reconstruction reimplemented the frozen authoring scheduler against public validation inputs. It did not import or copy `experiments/frontier-v1/scheduletool`.

No private evaluator root was opened, hashed, mounted, or enumerated. No validation or held-out label was opened. No schedule was executed. No model, arm, trial, prospective grade, validation result, held-out result, or outcome was produced or observed. `seal --write` was not run.

## Verdict

Identity, ancestry, six-path boundary, source-manifest linkage, generator safety, frozen-algorithm equivalence, independent reconstruction, Williams construction, raw-byte and canonical digests, tests, vet, quality, protected state, custody, and the authoring-only execution boundary all hold. There are no P0, P1, P2, or unresolved P3 findings.

`ACCEPT` makes this exact schedule digest eligible for a focused project-chair freeze decision. It does not freeze the schedule, open Gate 1A, or establish that the accepted runner can execute the schedule.

## Independently recomputed identities

| Artifact | Algorithm | Result |
| --- | --- | --- |
| Candidate commit tree vs parent `2822b1e8…` | `git diff --name-only` | exactly the six candidate paths |
| Report commit tree vs candidate `44fccf51…` | `git diff-tree --name-only -r bb473e22…` | exactly `docs/reviews/2026-08-21-protocol-v4-validation-schedule-candidate.md` |
| Candidate report raw bytes at `bb473e22…` | SHA-256 of git blob; 7,731 bytes; byte-equal to the LF worktree file | `sha256:e491d484e0c768feedcbc6571153530fe7a5f6ea91b4be241415c3451469968a` |
| `experiments/frontier-v1/README.md` | SHA-256 of git blob at `44fccf51…`; 17,689 bytes | `sha256:c424ef79e44594069a7d0b0dd8d4dc4311ccd41933bf411944b6309262567558` |
| `experiments/frontier-v1/schedules/validation.json` raw bytes | SHA-256 of git blob at `44fccf51…`; 391,950 bytes | `sha256:6b264a8daff60d9b507f11d760f1bbf19dec556f2aebd700dc5fbdc8a8c56aec` |
| `experiments/frontier-v1/scheduletool/main.go` | SHA-256 of git blob; 8,138 bytes | `sha256:cd959297da5cf796c6d0d386ba5cab67522b186216999fa81352a43aa5600de4` |
| `experiments/frontier-v1/scheduletool/main_test.go` | SHA-256 of git blob; 1,154 bytes | `sha256:1e6d5cad12c29b1db2301c0239fc80083d5cd489ce85ecebd22d509b14b21495` |
| `experiments/frontier-v1/scheduletool/schedule.go` | SHA-256 of git blob; 5,520 bytes | `sha256:37325fa40e645da066e69c9377f0f66f312799d1862861362c6f3ff36789a182` |
| `experiments/frontier-v1/scheduletool/schedule_test.go` | SHA-256 of git blob; 2,784 bytes | `sha256:b7b27ef528ce3bb51548bcad255a45aee26e0cb41f32ad83650e69bb84c7544d` |
| Validation manifest canonical JSON | Independent RFC 8785 (integer subset, UTF-16 key order) of parent blob `2822b1e8…`; frozen `corpusctl digest` agreed | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| Validation schedule canonical JSON | Same independent encoder on the candidate schedule; frozen `corpusctl digest` and candidate `--verify` agreed | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Frozen runner `schedule.go` | SHA-256 of worktree/git bytes; equals the accepted freeze pin in `pre-validation-artifacts.json` | `sha256:d04884d84036d3d8dd8cdd6f1fccdc4f889f3289831bb9ecec0e84292e894c95` |
| World-build digest from `make quality` | `cmd/build-manifest` then `quality-check compare` against `experiments/frontier-v1/artifacts/world-build.linux-amd64.json` | `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d` |

Ancestry: detached HEAD `bb473e22…`; sole parent `44fccf51…`; that commit's sole parent `2822b1e8…`. Import `ac47da46…` is an ancestor of the candidate parent.

`experiments/frontier-v1/pre-validation-artifacts.json` remains `status: partial`, `remaining: ["schedule digest"]`, `may_open_validation: false`, and `may_open_held_out: false`. No `validation_schedule` freeze entry exists.

## Requirement matrix

| Check | Result | Notes |
| --- | --- | --- |
| 1. Identity, ancestry, and six-path boundary | Pass | Required parents; candidate changes exactly the six listed paths and no other path, all additions except a 13-line README insertion; report commit changes only the candidate report; recomputed raw lengths and SHA-256 values match the candidate report and the assignment pins; none of those paths is a private label, result, runtime stream, trial, grade, analysis output, gate change, or frozen runner/grader change. Candidate files contain LF only. |
| 2. Exact public source and case-manifest linkage | Pass | `git diff --name-only 8319a3e7… 2822b1e8… -- experiments/frontier-v1/corpus experiments/frontier-v1/fixtures experiments/frontier-v1/manifests` is empty. Git blob maps for the eight sealed public groups (validation and held-out cases, fixtures, manifests, and label-digest registries) are identical across payload `8319a3e7…`, import `ac47da46…`, and parent `2822b1e8…`. Manifest is sealed version 1 for `validation`, 120 cases, 24 families. All 120 scheduled case IDs and family IDs equal the manifest set. Independently recomputed canonical digest of every public validation case equals its manifest `input_digest`. Each case's `case_id`, `family_id`, and `sandbox_fixture` equal its manifest row. Five cases per family. Source identity is blob equality with the reviewed payload, not filename matching. |
| 3. Generator safety and frozen-algorithm equivalence | Pass | Static inspection of all four scheduletool files: the command reads only `experiments/frontier-v1/manifests/validation.json` and the 120 public files under `experiments/frontier-v1/corpus/validation/` named by that manifest. It compares `sandbox_fixture` as a string and does not read fixture content, labels, private roots, results, runtime output, trials, grades, or outcomes. No `net`, `net/http`, `os/exec`, subprocess, model, runtime, runner, grader, or analysis import. `--write` refuses an existing path (`main.go` 98–105) and `ensureInside` keeps output inside the repository (`main.go` 257–262). Shared scheduler functions with frozen `experiments/frontier-v1/runner/schedule.go` are byte-identical after the allowed adaptation: `generateSchedule` takes a `tranche` argument instead of hard-coding `"authoring"`, and cases are `scheduleCase` rather than `runnableCase`. Seed `20260817`, three repetitions, arms `A,B,C,D,E`, sorted family/case inputs, SplitMix64, rejection sampling, Fisher-Yates, family blocking, `(case_id, repetition)` pairing keys, and the ten-row odd-treatment Williams design are unchanged. Candidate tests and vet passed; they were not treated as independent proof. |
| 4. Independent schedule reconstruction | Pass | Disposable Python reconstruction outside the Phoenix checkout, written from the frozen authoring algorithm and public validation inputs, not from scheduletool. Exact deep equality with `experiments/frontier-v1/schedules/validation.json`, including entry order. Version 1, tranche `validation`, seed `20260817`, repetitions 3, arms `A,B,C,D,E`; 120 cases, 24 families, five cases per family; 360 pairing keys; 1,800 launches; launch indexes `0..1799`; pairing indexes `0..359`; each pairing key has one case, one repetition, one family, and all five arms consecutively; each case appears once per arm in each repetition; each family occupies exactly one contiguous 75-entry block; every five-arm row is the Williams row selected by local pairing position and family block; all ten Williams rows are permutations and every directed predecessor pair appears exactly twice. Raw-byte and canonical digests match the pins. Candidate `--verify` was run only after that reconstruction succeeded and printed the same schedule and manifest digests. Frozen `corpusctl digest` independently agreed. |
| 5. Reproducible public checks | Pass | From the detached review worktree: `go test -count=1 ./experiments/frontier-v1/scheduletool` exit 0; `go vet ./experiments/frontier-v1/scheduletool` exit 0; `go run ./experiments/frontier-v1/scheduletool --repo-root . --verify experiments/frontier-v1/schedules/validation.json` printed `schedule: sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813 (1800 launches); manifest: sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490`; corpusctl `go test -count=1 ./...` and `go vet ./...` exit 0; `make quality` exit 0 and reproduced world-build `sha256:27c2f53775ac783b9085698068fa4660bead917352e11e1305a41dcfcce0188d`; `git diff --check 2822b1e8… 44fccf51…` empty. Gitignored `bin/` and `build/` from `make quality` were removed afterward. Tracked worktree status remained empty. |
| 6. Protected state and authorization boundary | Pass | Candidate does not change protocol, runtime/prompts, Arm A schemas, Arm B, worlds, accepted runner files, grader, analysis, world build, imported corpus, registries, manifests, or existing freeze entries. `pre-validation-artifacts.json` still partial with both gates false. No validation or held-out private-label directory exists in the review worktree or the implementation checkout. Authoring labels and authoring results remain; they were not opened. Candidate report (`docs/reviews/2026-08-21-protocol-v4-validation-schedule-candidate.md` 114–124) and README (`experiments/frontier-v1/README.md` 151, 188) state that the accepted runner is authoring-only and cannot execute this schedule. Frozen runner `loadCase` rejects non-`authoring_` IDs (`runner/materialize.go` 21–22) and `validateSchedule` requires tranche `authoring` (`runner/schedule.go` 83–84). That limitation is a separate pre-Gate blocker, not a defect in this candidate. |

## Independent reconstruction method

The reconstruction lived only under `D:\Work\personal\phoenix-schedule-review-scratch`. It loaded the public sealed validation manifest and the 120 public validation case files from the review worktree. It grouped cases by family, sorted family IDs and case IDs, then applied SplitMix64 state `20260817` with the unsigned `(0-bound)%bound` rejection threshold and Fisher-Yates shuffles of families, per-family case order, and `(case_id, repetition)` pairing keys. Each family became one contiguous block. Each pairing key emitted one Williams row, selected by `(local_pairing_index + family_block) mod 10`, with all five arms consecutive.

Canonical digests used an independent RFC 8785 encoder for the integer/string/object/array subset frozen `corpusctl` accepts (UTF-16 code-unit key order, integer numbers only). That encoder, frozen `corpusctl digest`, and the candidate verifier agreed on the manifest and the schedule. The reconstruction did not call scheduletool except for the later `--verify` required by the assignment.

## Findings

No P0, P1, P2, or P3 findings.

## Accepted limitations

- The reviewer is a same-machine fresh session. That is recorded above and does not reuse a prohibited role.
- Candidate self-tests are not independent proof. The reconstruction and blob/source checks are.
- `make quality` format-check, gocyclo, and dupl still omit `experiments/frontier-v1/scheduletool`, as they omit `corpusctl`. Expanding those paths would have been a seventh-path Makefile change. `go test ./...`, `go vet ./...`, and staticcheck include the package. `gofmt -l` on the four files was empty. This does not affect schedule identity.
- The generator canonicalizes with `github.com/gowebpki/jcs`. Frozen `corpusctl` Canonical independently agreed on this payload.
- The accepted scheduled runner remains authoring-only. Independent acceptance of this schedule does not create an execution path.

## Authoring-only execution blocker

The frozen runner still rejects validation cases and a schedule whose tranche is `validation`. This candidate correctly leaves that boundary in place. `ACCEPT` here does not authorize modifying runner or grader bytes, executing the schedule, or opening Gate 1A.

Gate 1A must remain closed until a separate independently reviewed artifact and an explicit project-chair decision resolve that execution boundary. Held-out remains closed.

## Smallest next artifact

A focused project-chair freeze decision for this exact schedule digest:

`sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`

That later patch may add an accepted `validation_schedule` entry to `experiments/frontier-v1/pre-validation-artifacts.json`, set `status: complete`, set `remaining: []`, and keep `may_open_validation: false` and `may_open_held_out: false`. It must not modify any frozen implementation byte or authorize execution.
