# 0004: Frontier experiment go/no-go rules

- **Status:** Draft; numeric thresholds and cost envelope not frozen
- **Date:** 2026-08-17
- **Protocol:** `experiments/frontier-v1/protocol.json`

## Decision structure

No validation or held-out run starts while this decision is draft. Freezing it requires an evaluation reviewer independent of implementation to accept the unit of analysis, effect thresholds, uncertainty method, repetitions, missing-data treatment, multiplicity policy, and total cost ceiling.

Every claim records:

- population and task classes;
- intervention and comparator;
- unit of analysis and clustering variables;
- primary endpoint and numeric threshold;
- uncertainty bound and minimum detectable effect;
- repetitions, seed handling, context reset, and case ordering;
- exclusions fixed before the run;
- retry, timeout, missing, and indeterminate handling;
- multiplicity policy and stopping rule;
- one forced verdict when the threshold is missed.

## Fixed qualitative verdicts

These verdicts are frozen now; only their numeric thresholds remain open:

- C must not regress against A on direct-task success. A regression kills the current surface.
- C versus B is the headline package comparison. Failure after token, latency, and daemon cost kills or narrows the headline claim.
- C versus D determines whether the frontier survives.
- D versus E determines whether teaching refusals survive.
- In Phase 2, counted C versus authored D′ determines whether counting remains enabled.
- An indeterminate result does not count as a pass. It requires a new, adequately powered sealed tranche.

## Proposed endpoints for independent review

The review must assign numbers to these endpoint forms without inspecting validation or held-out outcomes:

| Claim | Primary endpoint form | Forced failure verdict |
| --- | --- | --- |
| no direct-task tax | paired success-rate difference, C minus A | repair or reject the surface |
| headline value | C minus B success plus tokens-to-success ratio | kill or narrow Phoenix |
| frontier value | C minus D on cascade, far-discovery, recovery, and temptation | remove the frontier |
| refusal value | D minus E recovery after a state refusal | replace teaching refusals with typed errors |
| O(1) wake-up | maximum standing-token spread across 10/100/1,000 verbs | reject the constant-surface claim |
| daemon cost | per-act latency distribution and wall-time contribution | optimize or reject the surface |
| learning value | counted C minus authored D′ on the untouched held-out tranche | keep counting disabled |

## Analysis requirements

- Each trial starts with a fresh model context, sandbox, and isolated world state.
- Arms use identical underlying verb implementations, side effects, limits, and typed payloads.
- Seeds are paired where the runtime exposes them; otherwise repetitions remain paired by case and order block.
- Case order is randomized or counterbalanced within a pre-generated block.
- Confidence intervals cluster repeated observations by case family and seed.
- Tokens-to-success excludes failed trials from neither the denominator nor cost reporting; failures report their consumed tokens separately.
- Results publish all primary measures with explicit numerators, denominators, exclusions, retries, and missing trials.
- Any outcome-informed change burns the opened tranche before another gate attempt.

## Unresolved gate inputs

The following values remain `null` in `protocol.json` and block every outcome run:

- minimum cases per class and generating family;
- repetitions per case and arm;
- confidence level and interval method;
- minimum detectable effects and pass thresholds;
- per-trial timeout, retry limit, and cost ceiling;
- total experiment budget;
- handling of provider errors and runtime outages;
- correction for multiple primary comparisons.

These are product and research-budget choices, not implementation defaults. They must be accepted before the protocol's `frozen` field can become `true`.
