# Frontier-v1 custody-safe validation execution boundary

Status: replacement candidate. This document specifies an execution path; it does not authorize execution. Both outcome gates remain closed, and `held_out` has no runner path.

## Fixed inputs

The validation path accepts the 120 public `corpus/validation` cases and their content-addressed public fixtures only. It executes only the committed version-1 validation schedule with seed `20260817`, three repetitions, arms A–E, 360 complete pairing keys, and 1,800 launches. The runner requires these accepted identities:

| Input | Required identity |
| --- | --- |
| Validation schedule, canonical JSON | `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813` |
| Validation manifest, canonical JSON | `sha256:57ccc0c754f7c2beb74063dacc4bad26ab8b598ac2da9b6a8d637afd391fd490` |
| Validation label-digest registry, raw bytes | `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816` |
| Private label archive index | `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4` |
| Deterministic grader | `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc` |

The frozen A–E runtime, prompts, surfaces, trial limits, retry policy, Arm B document, world, and evidence schema are unchanged. Validation requires exactly a 300 USD run-budget boundary, a 0.15 USD per-trial cap, and a 180-second per-trial timeout.

## Gate and custody protocol

The runner performs these checks before output creation, Phoenix build, runtime verification, or model invocation:

1. `pre-validation-artifacts.json` is complete, `may_open_validation` is true, and `may_open_held_out` is false.
2. The gate document names the frozen schedule and source manifest identities above.
3. The committed schedule, public manifest, and public label-digest registry still match their pins.
4. The requested untrimmed per-trial cap parses to exactly 0.15 USD and the timeout is exactly 180 seconds. Whitespace-padded cap values are invalid.
5. The selected schedule is the exact deterministic validation schedule and its canonical digest matches the pin.
6. The custodian grader path is absolute, resolves outside the implementation repository, and names a regular file.
7. The custodian's strict `describe` JSON matches all five identities, tranche `validation`, version 1, and 120 cases.

The trial-limit check occurs before schedule preparation, output-directory inspection or creation, and custodian path resolution or handshake. Any different, whitespace-padded, or malformed cap, and any different timeout, is rejected without output mutation or custodian contact.

The runner hashes the external executable after resolving the boundary and records that adapter hash plus the `describe` document in retained summary evidence. Before each grade it re-hashes the executable and stops if its identity changed.

For grading, the runner invokes:

```text
<external-grader> grade --case-id <validation_case_id> --trial <absolute-trial-evidence-path>
```

No label path, private root, registry contents, expected outcome, or archive location is passed. Standard error is discarded, error messages are generic, command time is capped at 30 seconds, and standard output is capped at 4 MiB. The returned JSON must contain exactly `case_id`, `checks`, and `status`; check objects must contain exactly `id`, `kind`, `verdict`, and `reason`. Case identity, unique nonempty checks, allowed verdicts, and aggregate status consistency are verified before normalized grade evidence is retained.

## Evidence and stop rules

The existing scheduled runner retains sanitized runtime streams, trial records, normalized grade records, assignment records, and one scheduled summary. The summary's tranche is `validation` and includes the frozen schedule digest, grader digest, custodian description, and executable hash. Existing pairing-budget, retry, ITT, indeterminate, and safety-stop rules remain in force.

Any gate, identity, schedule, custody handshake, grade-shape, or executable-change failure stops safely. It does not fall back to authoring labels or an in-repository grader. A runtime or grading safety stop leaves the scheduled run indeterminate under the existing contract.

## Authorization boundary

Independent review must verify the exact replacement bytes, public-only tests, closed-gate behavior, and custody contract. Acceptance must be followed by a replacement refreeze. Only then may a distinct project-chair decision set `may_open_validation` true. That decision must leave `may_open_held_out` false. Neither this candidate nor its review may run a model, obtain a private grade, or observe a validation outcome.
