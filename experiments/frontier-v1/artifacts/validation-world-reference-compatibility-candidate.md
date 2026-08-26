# Gate 1A validation world-reference compatibility candidate

Status: review candidate. Both outcome gates are closed.

## Problem

The authoring repair accepted by decision 0014 changed the production and authoring world canonical identity from `sha256:f5f6b2f4ea695705e333b237296f0197fd8f22af77d66e9ecb08ef769fd1615b` to `sha256:d7f93051030c3f7a03442a99b12be77949a55f1ced1985bc34e16bd5d57289bc`. The 120 sealed validation case files and their manifest retained the former identity. The first launch attempt under decision 0019 stopped during public preflight before creating output or invoking a model.

## Candidate

The candidate makes the minimum compatibility repair:

- replace only `world_ref` in each of the 120 public validation case files;
- regenerate `manifests/validation.json` with the accepted `corpusctl seal --write` command, preserving case IDs, family IDs, fixtures, classes, label digests, counts, and coverage;
- update the runner's frozen validation-manifest identity and custody-boundary documentation;
- add a no-model runner test that loads all 120 validation cases and applies the production world preflight while both outcome gates are closed; and
- update the schedule-tool digest assertion without changing the schedule bytes.

The validation schedule remains `sha256:b38a0eaab063ba39dcbbc896c7b74ef085587177d3f58edcea0439d56e075813`. The label registry remains raw `sha256:847c510ca4449856db075fa85a75801dd6e62629fca028f54039e1f962c1c816`. Grader semantics remain `sha256:36abfbec8dd5365605d43ddbce796954ee24348acf1ea0a76b65365c2ee7dcfc`. The private archive remains `sha256:5a320a8742185e0471c6add885d1861950bbde8c7c312afbe907ae99762922e4`.

The replacement public validation manifest canonical digest is `sha256:39acbad5e45ad65302659cd0875bdfe589165ede9b60b6448ac4b09ccfb1e0c6`.

## Exclusions

The candidate does not change the production world, world build, prompts, arm behavior, runtime, case goals, case IDs, fixtures, families, labels, grading logic, schedule order, repetitions, arms, statistical analysis, held-out corpus, or any validation outcome. No private label is opened. No model, trial, grade, output directory, or cost is produced.

Independent review and a later refreeze are required before Gate 1A may reopen.
