# Phase 1 analysis contract

This command turns one scheduled-runner summary and its retained evidence into a deterministic analysis JSON document and, optionally, a Markdown report. It does not run a model, generate a schedule, grade a case, or change a gate.

## Inputs

The command requires:

- a version 1 `scheduled-summary.json` containing exactly one assignment record for every assigned trial;
- the matching tranche manifest, which supplies outcome-free case class and generating-family metadata;
- every trial and sanitized runtime path referenced by the scheduled summary.

The summary and manifest must name the same tranche. The summary must contain every manifest case at repetitions 0, 1, and 2, and every `(case_id, repetition)` key must contain A–E exactly once. Every result family must match the manifest. The command refuses to overwrite analysis or report output.

```powershell
go run ./experiments/frontier-v1/analysis `
  --repo-root . `
  --summary <scheduled-summary.json> `
  --manifest <tranche-manifest.json> `
  --output <analysis.json> `
  --report <report.md>
```

Do not run this command on validation or held-out evidence until the corresponding gate is independently opened.

## Missingness and pairing

The assigned case-arm-repetition trial is the ITT unit. `manual_required`, unresolved infrastructure, budget stops, safety stops, timeouts, and cap hits remain failures. Complete-case sensitivity excludes unresolved, budget-stopped, and safety-stopped assignments; it does not turn a manual-required grade into missing data.

If a safety stop interrupts a five-arm pairing key, the analysis discards the entire key as observed outcome evidence, sets all five ITT outcomes to failure, and counts all five assignments as unresolved. A budget or safety stop makes the tranche indeterminate. The same happens when unresolved rates exceed 0.05 in an arm or differ by more than 0.02 between any precommitted pair. Rates use assigned trials before rounding.

Cap-hit imbalance is narrower. A difference above 0.02 makes the affected cost-ratio claim indeterminate; it does not veto a capability pass.

## Inference

For 20 or more generating families, the implementation performs 10,000 paired hierarchical bootstrap replicates with seed `20260817`. Each replicate samples families, then cases within the sampled family, and retains every arm pair and repetition for a sampled case.

Below 20 families, it first averages repetitions within case and cases within family. It then uses the unweighted family means in an exact sign-flip distribution when `G <= 16`, or 100,000 seeded sign flips when `16 < G < 20`. One-sided bounds use the empirical 5th or 95th percentile. The harm p-value is the lower-tail probability of the observed family-mean difference; Monte Carlo p-values use the add-one correction.

The complete-case, one-vote-per-case, and one-vote-per-family estimates accompany every claim. They cannot replace the registered ITT decision.

Paired-success token ratios include only pairing keys where both compared arms passed. Tokens are totaled within family, ratios are analyzed on the log scale, and a zero family total receives `+0.5` before logging. The ratio is indeterminate if either arm has fewer than ten successes or the compared cap-hit rates differ by more than 0.02. Family log ratios use the hierarchical bootstrap at 20 or more eligible families and the registered sign-flip rule below 20.

## Evidence-derived measures

The following definitions make the descriptive measures reproducible from retained evidence:

- A wrong verb is an executable act whose Phoenix envelope ends in `fail` or `absent`. A `refused` act is state-inappropriate but is not counted as a wrong verb.
- A dead end is a timeout, turn-limit, cost-cap, agent-error, or malformed-output termination. A terminal graded failure with no executable act is also no progress. A failed trial is a help request when its final message asks the user to provide, clarify, specify, or confirm missing information, says that more information or context is required, or says it cannot proceed without something. Unlaunched and unresolved assignments are reported separately, not relabeled as dead ends.
- Recovery on a registered recovery case requires a first non-`ok` executable act, a later `ok` executable alternative, and a passing final grade. The same rule applies to D refusals and E plain typed errors.
- Frontier linkage is reconstructed from the sanitized runtime stream. A take requires the next Phoenix executable call to match a displayed call's handle, verb, and arguments. The state field may be omitted because the accepted v4 admission rule restores it only for the current pending frontier. A take-and-succeed also requires a passing grade.

Orientations remain separate from executable acts for act counts and sequence positions. Their displayed calls remain eligible for frontier linkage.

## Decisions and limits

The command implements the four Phase 1 decisions in `protocol.json`: direct harm, headline capability or efficiency-only, frontier isolation, and teaching-refusal isolation. Component isolations can reject a component but cannot create a headline pass. Cost cannot veto capability.

The rendered report preserves the scope limit and states that it cannot open Gate 1A, validation, or held-out. Engineering-gate evidence is separate from the scheduled outcome analysis and must be attached before a project-level surface claim is made.
