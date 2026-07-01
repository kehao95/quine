# 6M.07 Aggregate Morphology Metrics Summary

**Status:** analysis generated.

**Rows:** 61
**Rows with explicit unknown fields:** 61

## Rows By Experiment

| Experiment | Rows |
|------------|------|
| `6M.01` | 10 |
| `6M.02` | 13 |
| `6M.03` | 7 |
| `6M.04` | 6 |
| `6M.05` | 6 |
| `6M.06` | 4 |
| `6M.08` | 4 |
| `6M.09` | 5 |
| `6M.10` | 6 |

## Rows By Body Class

| Body class | Rows |
|------------|------|
| `custom-control-adapter-attempt` | 9 |
| `custom-stream-body` | 10 |
| `failed-no-successor` | 1 |
| `fd-adapter` | 1 |
| `fd-mediated-adapter` | 2 |
| `handoff-wrapper-preserved-runtime` | 1 |
| `one-shot-immediate-successor` | 1 |
| `preserved-runtime-reentry` | 1 |
| `quine-runtime-or-seed` | 5 |
| `rebuilt-quine-runtime` | 3 |
| `rebuilt-runtime` | 4 |
| `source-derived-rebuilt-quine-runtime` | 12 |
| `standard-utility` | 2 |
| `unknown` | 9 |

## Rows By Carrier

| Carrier | Rows |
|---------|------|
| `ctl-wake` | 4 |
| `fd` | 1 |
| `immediate-mission` | 1 |
| `inherited-fd` | 2 |
| `native-control` | 36 |
| `native-control-withheld` | 1 |
| `none` | 4 |
| `stdin` | 12 |

## Rows By Failure Class

| Failure class | Rows |
|---------------|------|
| `carrier-conflation-pre-handoff-future-stream` | 1 |
| `control-carrier-intentionally-withheld` | 1 |
| `control-proof-not-observed` | 4 |
| `methodological-negative` | 1 |
| `multi-wake-proof-incomplete` | 1 |
| `native-control-addressability-failure` | 1 |
| `no-handoff-before-control-timeout` | 4 |
| `none` | 47 |
| `seed-negative-for-runtime-rebuild` | 1 |

## Caveat

This table is a normalization pass, not a success-rate dataset. Historical seed
rows and retained-run rows have different instrumentation depth; unknown numeric
fields are preserved rather than inferred.
