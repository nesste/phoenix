# Frontier-v1 custody-safe validation execution boundary

Status: protocol-v5 amendment candidate. This document specifies an execution path; it does not authorize execution. Both outcome gates remain closed, and `held_out` has no runner path.

Protocol v5 retires every v4-sealed unopened candidate. The committed five-arm validation schedule, manifest, label-digest registry, and private archive identities below are the retired v4 pins; the runner's deterministic A–D schedule construction refuses the retired schedule, so no validation execution is possible until independent evaluators seal a new disjoint v5 tranche and its identities replace these rows at refreeze.

## Fixed inputs

The validation path accepts the public `corpus/validation` cases and their content-addressed public fixtures only. It executes only a committed version-1 validation schedule with seed `20260817`, three repetitions, arms A–D, and complete contiguous pairing keys. The runner currently records these identities:

| Input | Recorded identity | Status |
| --- | --- | --- |
| Validation schedule, canonical JSON | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` | retired v4; refused by the A–D construction |
| Validation manifest, canonical JSON | `sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6` | retired v4; replacement pending v5 seal |
| Validation label-digest registry, raw bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` | retired v4; replacement pending v5 seal |
| Private label archive index | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` | retired v4; replacement pending v5 seal |
| Deterministic grader | `sha256:8146a68a11a7593d8bfdeed102175143267a80c018f9222e678b20512baebf0b` | v5 §5 absence-acceptance candidate |

The A–D runtime, conditional C/D prompts, surfaces, trial limits, retry policy, Arm B document, world, and evidence schema are those of the v5 amendment payload. Validation requires exactly a 300 USD run-budget boundary, a 0.15 USD per-trial cap, and a 180-second per-trial timeout.

## Gate and custody protocol

The runner performs these checks before output creation, Phoenix build, runtime verification, or model invocation:

1. `pre-validation-artifacts.json` is complete, `may_open_validation` is true, and `may_open_held_out` is false.
2. The gate document names the frozen schedule and source manifest identities above.
3. The committed schedule, public manifest, and public label-digest registry still match their pins.
4. The requested untrimmed per-trial cap parses to exactly 0.15 USD and the timeout is exactly 180 seconds. Whitespace-padded cap values are invalid.
5. The selected schedule is the exact deterministic A–D validation schedule and its canonical digest matches the pin; the retired v4 five-arm schedule fails this check by construction.
6. The custodian grader path is absolute, resolves outside the implementation repository, and names a regular file.
7. The custodian's strict `describe` JSON matches all five identities, tranche `validation`, version 1, and 120 cases.

The trial-limit check occurs before schedule preparation, output-directory inspection or creation, and custodian path resolution or handshake. Any different, whitespace-padded, or malformed cap, and any different timeout, is rejected without output mutation or custodian contact.

Validation runs build Phoenix with the frozen world-build recipe: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOTOOLCHAIN=go1.26.6`, `-trimpath -buildvcs=false -ldflags "-s -w -buildid= -X main.version=dev"`, at `bin/phoenix`. After that build and before any trial or model call, the live world-build digest must equal the frozen `sha256:425bab1cdf8528a1eb962cd06945268e519a1cea56d3169d6c0465e8e2ffdae4` or the run stops and the build is cleaned up. Authoring runs keep the host build recipe and skip the pin.

The runner hashes the external executable after resolving the boundary and records that adapter hash plus the `describe` document in retained summary evidence. Before each grade it re-hashes the executable and stops if its identity changed.

For grading, the runner invokes:

```text
<external-grader> grade --case-id <validation_case_id> --trial <absolute-trial-evidence-path>
```

No label path, private root, registry contents, expected outcome, or archive location is passed. Standard error is discarded, error messages are generic, command time is capped at 30 seconds, and standard output is capped at 4 MiB. The returned JSON must contain exactly `case_id`, `checks`, and `status`; check objects must contain exactly `id`, `kind`, `verdict`, `reason`, and the v5 `gating` flag. Case identity, unique nonempty checks, allowed verdicts, and aggregate status consistency — computed over gating checks only, so descriptive route checks are reported without failing the trial — are verified before normalized grade evidence is retained.

## Evidence and stop rules

The existing scheduled runner retains sanitized runtime streams, trial records, normalized grade records, assignment records, and one scheduled summary. The summary's tranche is `validation` and includes the frozen schedule digest, grader digest, custodian description, and executable hash. Existing pairing-budget, retry, ITT, indeterminate, and safety-stop rules remain in force.

Any gate, identity, schedule, custody handshake, grade-shape, or executable-change failure stops safely. It does not fall back to authoring labels or an in-repository grader. A runtime or grading safety stop leaves the scheduled run indeterminate under the existing contract.

## Execution resilience, mandatory resume, and the post-closure gate

Host requirements for the execution window: automatic OS restarts disabled (Windows Update deferred, or the prepared WSL2/Linux host used), no interactive login sessions, and the custodian process supervised.

At every pairing-key boundary the runner writes three records beside the evidence: the durable `scheduled-checkpoint.json`, the outcome-free `scheduled-summary.partial.json` (launch index, completed pairing-key list, per-arm assignment counts, artifact digests, timestamps, and cumulative spend only — no success counts, no grade tallies, no per-arm outcome field), and an appended entry in `process-events.jsonl`. That process-event log is append-only and carries start, checkpoint, stop, and termination events with timestamps and, where the host exposes it, the originating principal; a signal-driven termination records that the host does not expose the sending principal. The log's digest is committed with the custody record, and the completed `scheduled-summary.json` remains the completion criterion.

If an execution is interrupted by infrastructure failure, resume from the last durable checkpoint is **mandatory**, not discretionary; abandonment by choice does not exist as an outcome, and the tranche closes indeterminate whenever a condition fails. The runner refuses the resume unless `--resume-attestation` names a custodian record establishing all of:

1. no human or agent inspected any outcome artifact between interruption and resume;
2. the frozen runner, world, corpus, schedule, and grader bytes are digest-identical at resume, with the live schedule and world-build digests recorded in the attestation;
3. resume is authorized within 72 hours of the interruption, with the timestamps ordered interruption, cause classification, condition verification, authorization. Three further requirements make that window checkable rather than self-reported, and a custodian must satisfy them explicitly:
   - `interrupted_at` must be at or after the last `recorded_at` in `process-events.jsonl`, allowing five minutes of clock skew. A signalled interruption writes its termination entry a moment *after* the signal arrives, so attesting the true interruption instant is correct and the skew allowance covers the difference; an interruption attested materially earlier than the log's last entry is refused, because the process demonstrably ran past it.
   - the 72 hours are measured from that same last log entry, not only from the attested interruption, so back-dating the interruption cannot manufacture a fresh window.
   - `resume_authorized_at` must be stamped within one hour of the moment the runner is invoked, and no more than five minutes ahead of it. **This one-hour bound is not a section 9 requirement**; it is an implementation control giving the promptness duty effect, by forcing the attestation to be authored at the moment of resume rather than prepared in advance. The section 9 obligations remain the 72-hour backstop and the duty to attempt resume as soon as conditions 1, 2, and 4 verify.
4. the cause was classified and attested before any outcome inspection and before the resume decision point, and is on the frozen outcome-uncorrelated list: host restart, power loss, hardware failure. OOM is struck because transcript-heavy failing trials make memory pressure outcome-correlated. A process kill qualifies only when the attestation, citing the process-event log by digest, establishes that no project participant, account, or agent initiated it;
5. this is the tranche's first resume. A second interruption closes the tranche indeterminate.

The 72 hours are a backstop, not an allowance: resume must be attempted as soon as conditions 1, 2, and 4 verify, an authorization more than an hour after that moment requires a written lapse explanation in the attestation, and a closure by clock lapse requires a committed written explanation of why resume could not be attempted earlier. Authoring resumes are development loops and are exempt from the attestation rule.

The attestation is a JSON document and is decoded strictly, so an unknown or misspelled field is refused. It carries exactly: `v` (1), `tranche` (`validation`), `output_dir` (the repository-relative output directory), `custodian`, `cause`, `resume_index`, the four timestamps `interrupted_at` / `cause_classified_at` / `conditions_verified_at` / `resume_authorized_at` (RFC3339), the four booleans `no_outcome_inspection_between_interruption_and_resume` / `cause_classified_before_inspection_and_resume_decision` / `frozen_bytes_digest_identical_at_resume` / `kill_principal_established_as_non_participant`, `process_event_log` (the literal `process-events.jsonl`) and `process_event_log_sha256`, `frozen_digests_verified` (an object carrying `schedule_digest` and `world_build`), `promptness_statement`, and the optional `lapse_explanation` and `process_event_log_review`.

One state has no resume path: a custodian process that wrote its `start` event but died before its first pairing-key checkpoint leaves progress records with no checkpoint, which the runner refuses wholesale. The window is milliseconds, and the rule is what stops a phantom `start` entry from permanently refusing a later genuine resume; such a tranche closes indeterminate.

**Post-closure gate.** After any validation execution that ends without a completed schedule — interruption, budget stop, or lapse — no subsequent validation execution is authorized until a closure record (cause classification, custodian logs, attestations, and the lapse explanation if any) has been independently reviewed and committed. The gate document records that review path and its verdict; `requireValidationGate` refuses an otherwise-open gate whose last execution status is not `complete` until the named record exists and its verdict is ACCEPT. The record must be a non-empty markdown file under `docs/reviews/` whose own first verdict statement reads ACCEPT: the runner parses the verdict label rather than searching for the token, because every review record in this project contains the word ACCEPT and every evaluator prompt lists the available verdicts. A record whose first verdict statement is REVISE is refused even when the gate document claims otherwise. This bounds the engineered-indeterminate retry channel: repeated close-and-retry would condition the eventually completed tranche on side signals, so each retry must survive an independent audit of why the last execution died. The 1,475-launch execution closed by decision 0023 has no such review yet, so the gate refuses a further validation execution on that ground alone.

## Authorization boundary

Independent review must verify the exact replacement bytes, public-only tests, closed-gate behavior, and custody contract. Acceptance must be followed by a replacement refreeze. Only then may a distinct project-chair decision set `may_open_validation` true. That decision must leave `may_open_held_out` false. Neither this candidate nor its review may run a model, obtain a private grade, or observe a validation outcome.
