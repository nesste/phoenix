# Second independent review prompt — protocol-v5 section 9 execution-resilience payload, revision 2

You are an evaluation reviewer independent of this payload's implementation. You did not author it, and you are **not** the reviewer who produced the first review record. Review by reading the repository at the pinned commit; do not run a model, do not obtain a private grade, and do not observe any validation outcome. Both outcome gates are closed and must stay closed.

## Why there is a second review

The first payload (`7a1ea42513acd3d55e276eabc2459da4a037acc0`) was independently reviewed and returned **REVISE**: two P1 findings, six P2, nine P3, no P0. That record is committed at `docs/reviews/2026-08-27-execution-resilience-payload-review.md`. The revision under review claims to fix both P1s, all six P2s, and five of the nine P3s.

Your scope is the **whole payload**, not only the deltas. A fix that resolves the reported symptom while leaving the underlying weakness, or that introduces a new defect elsewhere, is exactly what a second review exists to catch. Read the first review record — but treat its verdicts as claims, not as settled fact. If you think a finding it raised was wrong, or a residual it accepted should not have been accepted, say so.

## Pinned inputs

Verify these before reviewing anything else, and report the values you computed.

- Payload commit: `git rev-parse HEAD` must be **`e92eaacf6ab54b547b331d14749717da64b58e0d`**. Its parent must be `46030d4fa668b0e1b00d80ea91ee29717dfaf67e` (the first review record), whose parent is the superseded payload `7a1ea42513acd3d55e276eabc2459da4a037acc0`, whose parent is `2cf99331688f4705ed7d6cd3efed98941f59defc` (decision 0025, the frozen state being replaced).
- Do not commit anything, and do not modify the working tree, while you review. The tree is intentionally dirty with exactly one untracked file, this prompt. If anything else is dirty, report it as a P0.
- Candidate inventory `experiments/frontier-v1/artifacts/gate-1a-execution-resilience-candidate.json`, raw SHA-256 `sha256:cb36e2087f1d2b4e0a7d61965d472a6c9ee03ae498cf913fc2b1b1e9813aea47`, 16 files, ordinal-path set digest `sha256:95f33f2af58442464557b77287647f93393cdb313397721cd0bac8525511808e`.
- Candidate note `experiments/frontier-v1/artifacts/execution-resilience-candidate.md`, raw SHA-256 `sha256:05799da56672700913b0462ff457541a402c07cad9bdfaec84ede1c31a3ea3ed`.
- LF-normalized SHA-256: `runner/execution_resilience.go` `sha256:28d9c174fb5c457271b382281329b8b12f6cd48cfb42989427fc1cf4c34e14d8`; `runner/process_events.go` `sha256:cdfa5f84e8dd92dfc41d1904b289ce704e1000a1a8fe36504b5b85fa9c7792d4`; `runner/validation_gate.go` `sha256:df56726da5dc61753f1b761c09079e61174fe6c4368d7ce47a4bf1a687d68bb7`; `runner/scheduled_run.go` `sha256:e58360c3e83656c8cc996a577c6eed89380956c8c912e3ba380a70b664e427ad`; `runner/execution_resilience_test.go` `sha256:f39d36ed26e8c707e10aae9de1fbb94bab8e59f1bce9a4a7d519edadd8964546`; `pre-validation-artifacts.json` `sha256:ce14d55cdbb36c9102e71830101eb45f0bd29df91af35a51134c57b72b3f4cc2`.

Digest rule: CRLF→LF, then CR→LF, then SHA-256, rendered `sha256:<hex>`. Set digest: SHA-256 over ordinal-path-sorted UTF-8 lines `<path>\t<lf-normalized-file-digest>\n`.

## Governing sources

- `docs/plans/2026-08-27-protocol-v5-proposal.md` § 9 — the accepted specification. The payload must implement it, not reinterpret it.
- `docs/reviews/2026-08-27-protocol-v5-proposal-review-4.md`, **freeze obligation 5**.
- `docs/reviews/2026-08-27-execution-resilience-payload-review.md` — the first review.
- `experiments/frontier-v1/artifacts/execution-resilience-candidate.md` — the implementer's account, including a "What the first review changed" section. Claims to verify, not evidence.

## What you must check

**A. Each claimed fix, adversarially.** For every P1 and P2 the revision claims to fix, do not accept that the reported symptom is gone — determine whether the underlying weakness is gone.

1. **P1-1 (termination handler).** Does a signalled run now actually stop? Is the exit reachable in production, not only under the injected test double? Does the handler still write the entry *before* stopping, and is that write durable? What happens on a second signal, or on a signal arriving during the exit? Does moving the watcher into `runScheduledCases` leave any window at the start or end of a run where a signal is recorded but not acted on, or acted on but not recorded? Is anything now leaked or left uncleaned that the previous default behaviour also would have left?
2. **P1-2 (resume-counter corroboration).** The counter is now cross-checked against `start` entries in the process-event log. Can that be defeated? Consider: truncating or rewriting the log, deleting individual lines, a log with zero starts, a log whose entries are valid JSON but semantically wrong, and the interaction with `verifyAttestedProcessEventLog`'s digest check (which of the two runs first, and does either accept a state the other would reject?). Is the failure mode closed in every case?
3. **P2-3 (timing anchors).** The window is now anchored to the last process-event timestamp and the runner's clock. Try to defeat it: a custodian who controls the host clock, an attestation authored at resume time for a weeks-old interruption, a log whose last entry is recent because a later process appended to it. Are `resumeClockTolerance` and `resumeClockSkew` defensible values, and does the anchor change break any legitimate resume the specification intends to permit? A fix that refuses a *valid* resume is as much a defect as one that permits an invalid one.
4. **P2-4 (closure record).** Is requiring the string `ACCEPT` anywhere in the file a meaningful check, or is it defeatable and/or over-strict? Consider a REVISE record that quotes the word, and a legitimate ACCEPT record that does not contain that exact token.
5. **P2-1, P2-2, P2-6, P2-5** and the fixed P3s: verify each is real and complete.

**B. New defects introduced by the revision.** In particular: the mutex-guarded `liveProgress` handle is read from the signal goroutine while the run loop writes it — is the synchronization correct and complete, and is there any remaining data race (consider running `go test -race ./experiments/frontier-v1/runner/`)? Does the reordering inside `TestPreValidationFreezeMatchesAcceptedCandidates` change what that test guarantees? Does any new refusal path fire on a legitimate first run rather than only on a resume?

**C. The whole payload, as if fresh.** Everything the first prompt asked, restated: each § 9 condition located in code and tested for whether a resume can succeed with it unmet; the partial summary genuinely outcome-free; the process-event log honest about what the host does not expose; the post-closure gate encoded in the frozen gate rules and applying to the decision-0023 closure; the authoring-tranche exemption judged; the protocol, boundary, and README records matching the implementation with no overclaim; freeze hygiene with exactly one deliberate red (`TestPreValidationFreezeMatchesAcceptedCandidates`) and the world-build and grader digests unmoved.

**D. Residuals.** The note offers P3-1, P3-2, P3-5, P3-6, the authoring exemption, and the chair-discipline caveat on `validation_execution_status` as accepted residuals. For each, say whether it is correctly characterized and genuinely acceptable, or whether it should be a finding.

**E. Quality gates.** Run independently: `go test ./...`; the corpusctl suite; `go vet ./...`; staticcheck v0.7.0; `gocyclo -over 15` and `dupl -t 100` over `cmd internal verbs experiments/frontier-v1/runner experiments/frontier-v1/analysis`; `make validate-spec`; `make validate-authoring`; and `go test -race` on the runner package. Report anything beyond the single expected red.

## Verdict

Return **ACCEPT** or **REVISE**, with findings graded P0/P1/P2/P3, each with file, line, the failure it permits, and the fix. If you accept, give the explicit residual list worded for direct transcription into the refreeze's `accepted_findings`. Report every digest you computed. Do not soften findings to reach ACCEPT, and do not manufacture findings to look rigorous — a clean second review after a substantive revision is a legitimate outcome. Report what you verified and what you could not.
