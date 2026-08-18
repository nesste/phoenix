# 0005: Phase 0 independent review

- **Status:** Phase 0 research contract accepted
- **Date:** 2026-08-18
- **Candidate commit:** `cd84e888a362130c5dfabd240c9a50ba7406a588` (`docs: prepare Phase 0 blocker handoffs`)
- **Evaluation review:** `docs/reviews/2026-08-18-task-0.6-evaluation.md`
- **Implementation review:** `docs/reviews/2026-08-18-task-0.6-implementation.md`
- **Human-factors review:** `docs/reviews/2026-08-18-task-0.6-human-factors.md`

## Decision

**`ACCEPT`**

Phoenix Phase 0 is an accepted research contract. The surface spike, frozen result and world schemas, frozen protocol v3, authoring tranche, and independently created unopened validation and held_out candidates together satisfy the Phase 0 gate.

This acceptance authorizes importing the reviewed public sealed-corpus patch into the implementation workspace and beginning Phase 1 implementation. It does not open validation or held_out. Before either tranche opens, freeze the runtime invocation, system prompt, Arm A schemas, Arm B document, world build, runner, schedule, grader, analysis implementation, and report template by digest as required by `protocol.json`.

No sealed outcome was opened.

## Candidate prerequisites

At `cd84e888a362130c5dfabd240c9a50ba7406a588`:

- the implementation workspace contains no validation or held_out trial records, outcomes, cases, or full labels;
- `protocol.json` is `status: frozen`, `frozen: true`, and independently accepted (Task 0.5, 2026-08-17);
- Task 0.2 evidence is independently accepted (2026-08-18);
- `may_open_validation` and `may_open_held_out` are false.

The public sealed-corpus package exists outside this workspace and was reviewed before import.

## Reviewer roles and independence

| Role | Date | Verdict | Independence |
| --- | --- | --- | --- |
| Evaluation reviewer | 2026-08-18 | `ACCEPT` | Independent of Phoenix implementation and of the corpus author. Inspected private labels through the private channel. Did not copy expected outcomes into this repository. |
| Implementation reviewer | 2026-08-18 | `ACCEPT` | Independent of corpus labeling. Did not inspect private labels. Did not apply the sealed-corpus patch. |
| Human-factors reviewer | 2026-08-18 | `ACCEPT` | Independent of implementation authorship. Did not inspect a sealed-tranche run. |
| Project chair | 2026-08-18 | `ACCEPT` | Did not implement Phoenix, author the corpus, or inspect private labels. Combined the three reviews. |

Recorded limits:

- All four roles are Cursor agents operated by the project chair, who is also custodian of the private-label archive. This is not dual-control by four people.
- Corpus author and first-pass labeler were the same evaluator role. The label auditor was a mechanical second pass, not an independent re-label. The Phoenix implementer did not create, label, or adjudicate cases.
- One session may not approve its own implementation or corpus work. The implementation and corpus were produced in earlier sessions; this review used three separate reviewer sessions plus this chair record.

## Candidate identification

SHA-256 over LF-normalized UTF-8 bytes at the candidate commit, except sealed artifacts which were hashed in the evaluator checkout or private archive and were not present in the implementation workspace.

### Decisions, protocol, and spike results

| Artifact | Digest |
| --- | --- |
| `docs/decisions/0001-system-boundary.md` | `sha256:ba19dc35c956197446bba1f1d05043d6c7df7eb676c179d1e16f680a1b2dfbef` |
| `docs/decisions/0002-surface-and-stack.md` | `sha256:fce1fa773df3383a6d3c8f792c0a2842cc20e842184207310f7ce9d9bc0d194b` |
| `docs/decisions/0003-result-envelope.md` | `sha256:ab6ee936362a857fc4d3c6c7032ff563d70e8844e334e78cd0fdaebaf1c3ec52` |
| `docs/decisions/0004-go-no-go-rules.md` | `sha256:44410c1b56857c498daf3d267af376fe46339df6e670b610a684fa234bdae573` |
| `experiments/frontier-v1/protocol.json` | `sha256:95b284da8fd16a0bca1a42358e4158e9cb926be58c41b86d34d7ec74ea873a7c` |
| `experiments/surface-spike/results.json` | `sha256:a11c8d271d1c5eb7eaab098c4aa086e30b7728a34835e58b7c50e1c6b4cd96d0` |
| `experiments/frontier-v1/manifests/authoring.json` | `sha256:4a9a93a510d3a4c7379d8c962fd84c7d1867138cee0ecca55cbf39e417249f87` |

The independently accepted protocol-design freeze remains `protocol.json` `sha256:4d13c140742ba462d68af865d4ed12b08e6c2e1aba4f652fed4e35adf782a47d` and `0004` `sha256:399d337ef957ff1a43e7b384a03a7ddfe1e3a360b9f5337ae1d662338dfe3849`. The candidate bytes above differ only by the later Task 0.2 acceptance metadata (gate reason and remaining-blocker list). Claim tables, arms, budgets, and kill rules were not redesigned.

The independently accepted Task 0.2 evidence bytes remain `results.json` `sha256:376a96534eac11f08065ffc300bfb8b73fcec0c8c18a37336edf0c3731dbaf95` and `0002` `sha256:897c84d5891705997efeabb8f135fb127f75f48ed31e7fa63198b4d5449d973e`. The candidate `0002` and `results.json` bytes include the subsequent acceptance record.

### Sealed manifests and public evaluator report

Hashed in `D:\Work\personal\phoenix-evaluator-corpus` and `D:\Work\personal\phoenix-evaluator-private`. Not imported.

| Artifact | Digest |
| --- | --- |
| `experiments/frontier-v1/manifests/validation.json` | `sha256:bbc462f2241ae3c31fc95e6422e477a92cae4145a989712b511b541d0f077e3b` |
| `experiments/frontier-v1/manifests/held_out.json` | `sha256:0eb47d194bec2c5edb4a7e72289287777f35532a1006e233b8c8d01cf60fb2b1` |
| `experiments/frontier-v1/manifests/validation-label-digests.json` | `sha256:2e2febaf263db06d61a3109e997e96e0ebb2a3b387ee604213cb229f6c4da6dc` |
| `experiments/frontier-v1/manifests/held_out-label-digests.json` | `sha256:5742619559ba45777f0fa505b57359fcbb92733cd746a01f94782fa3415797be` |
| `docs/reviews/2026-08-18-frontier-v1-sealed-corpus-report.md` | `sha256:4c3d0451a88d31005adf30b0ef58882992a86bee2d944f5f713a5b6ac28ee69f` |
| Public payload (404 canonical file digests) | `sha256:89698a43ebbce508c37b2abfab5da8b11d2f95da6e2663e38969246cf79e6918` |
| Public patch `frontier-v1-sealed-corpus.patch` | `sha256:82f246b93c5e38454e4d7803f608b94e469a3198d6afb222423c5af5270df6ba` |
| Private label-set (240 labels, external) | `sha256:e1170e0cf99a33d86ba6ca6a533e057b58d1e1021375b50a79700b0e698c2009` |
| Grader | `sha256:3bba60f0f94f6686c0046c73c8507611e937083d57b66fd022bbef7cf32ca92c` |

Full labels remain with the external custodian. They are not in this repository.

### Task 0.6 review records

| Artifact | Digest |
| --- | --- |
| `docs/reviews/2026-08-18-task-0.6-evaluation.md` | `sha256:6958605e4d444f9d1d86a238b1c013e99402b2453df4d0f8830f22c9bbe6bfec` |
| `docs/reviews/2026-08-18-task-0.6-implementation.md` | `sha256:f9be6b9f7875e40da24c623cac74177df77b92d3f3426c925d10a44954f7597b` |
| `docs/reviews/2026-08-18-task-0.6-human-factors.md` | `sha256:cd5b351a9cc0f8ddfdb8e35f04d97ed4903e74995afd7faefade1dd26738c89d` |

## Corpus counts

Each sealed tranche, independently recomputed by the evaluation reviewer:

| Quantity | validation | held_out |
| --- | ---: | ---: |
| Cases | 120 | 120 |
| Generating families | 24 | 24 |
| Cases per family | 5 | 5 |
| Class totals | 24 `direct`, 24 `recovery`, 12 each other class | same |
| Frontier subset families | 16 | 16 |
| Withheld label digests | 120 | 120 |

Authoring remains 8 cases, 4 families, one of each class. Private labels were not copied into this workspace.

## Combined findings

No unresolved P0, P1, or P2 findings. Accepted P3 items are not Phase 0 blockers:

| ID | Owner | Finding | Artifact to clear |
| --- | --- | --- | --- |
| E-P3-1 | Phase 1 implementer, optional | Mix families ship 12 unused snapshot files per sealed tranche | Omit unreferenced mix versions, or record that fixture count is not used-snapshot count |
| E-P3-2 | corpusctl maintainer, optional | Authoring `validate` raw-byte compare fails on a CRLF working tree | LF-normalize in `validate`, or document the compare |
| I-P3-1 | Task 1.2 | Absent envelopes can still carry `error.details` or handle grants | Tighten the `absent` schema branch |
| I-P3-2 | Task 1.2 | No checked-in `fail` or no-alternative refusal fixtures | Add examples and `validate-spec` coverage |
| I-P3-3 | Task 1.4 | `measure-schema` ignores `--sizes` | Parse the flag or reject unknown args |
| H-P3-1 | Authoring runner, later | No retained full spike transcript | Authoring-only Episode 1 and 2 transcripts if claiming live why-line use |

## Implementation commands

Re-run by the implementation reviewer and confirmed by the chair at the candidate commit. All exited 0.

```powershell
Set-Location experiments/surface-spike
go test ./...
go vet ./...
go run ./cmd/validate-spec --repo-root ../..
go run ./cmd/measure-schema --sizes 10,100,1000

Set-Location ../frontier-v1/corpusctl
go test ./...
go vet ./...
go run ./cmd/corpusctl validate --repo-root ../../.. --tranche authoring --world-source experiments/frontier-v1/worlds/authoring.dev_repo.json
```

O(1) harness: `act` 1621 bytes at 10, 100, and 1,000 verbs (spread 0); `eval` 1490; flat tools 3861 / 38601 / 386001.

## Accepted limitations

- Conclusions remain limited to Claude Code 2.1.229, `claude-sonnet-5`, the frozen dev-repo world, and the eight task classes.
- Expected headline power is about 62% at a 15-point effect; the approximate 80% MDE against zero is 0.19.
- Linux evidence is Ubuntu 24.04.4 under WSL2 with `CGO_ENABLED=0`, not bare metal.
- Standing-token flatness at 10 / 100 / 1,000 verbs is inferred from identical serialized tool lists.
- With `n=20` and zero malformed calls, the Wilson 95% upper bound is 16.11% and is reported, not a pass floor.
- Authoring is one case per class; Arm B cannot be richer than that set.
- Held_out is a new-domain draw from the same puzzle grammar, not a structurally alien set.
- Shape legibility is accepted from constructed episodes plus aggregate spike evidence, not from a retained full model transcript.

## Unresolved blockers

No remaining Phase 0 contract blockers.

Pre-validation gates that this verdict keeps closed:

| Blocker | Owner | Artifact |
| --- | --- | --- |
| Sealed public files are not yet in the implementation workspace | implementer | Import `frontier-v1-sealed-corpus.patch` (digest above) without copying private labels |
| `may_open_validation` / `may_open_held_out` | protocol | Remain false until the `protocol.json` `artifact_freeze.before_validation` list is committed by digest |
| Outcome run | none | Not authorized |

## Exact next action

1. Import the reviewed public sealed-corpus patch into this workspace. Do not copy private labels. Keep every outcome gate closed.
2. Update `protocol.json` `gate.reason` to record that Task 0.6 accepted and sealed families exist, leaving `may_open_validation` and `may_open_held_out` false.
3. Begin Phase 1 implementation of the production daemon, session handles, verb registry, and authoring runner.
4. Before opening validation, freeze the artifacts listed in `protocol.json` `artifact_freeze.before_validation` by digest. Arm artifacts may use authoring outcomes only.
