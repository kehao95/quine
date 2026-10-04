# Cross-Model Replication & Threshold-Dependence Plan

Paper: `minimal-perceptual-prerequisites` (ALIFE 2026 accepted poster)
Created: `2026-06-06`
Status: planned but paused. The June 2026 supplement wave was invalidated and
deleted; any restart is gated on a paper-runtime compatibility audit under
current HEAD.
Canonical execution surface:
[`experiments/active/cooperation-dynamics/coop-16-cross-model-replication`](../../../../experiments/active/cooperation-dynamics/coop-16-cross-model-replication/README.md)

## Why this plan exists

All four reviewers flagged that the closure-threshold result rests on a single
model family (GPT-5.4). A cross-model study is the highest-value supplement: it
is the empirical answer to "why only GPT" (Reviewer 1) and to the generalization
concern (Reviewers 3, 4), and it removes the need for any written model-choice
rationale (camera-ready checklist item 6, deferred).

## What we already know (do not re-discover)

- **A prior cross-model probe already exists and leans NEGATIVE.**
  `8.15-cross-model-replication` tested GPT-5.4 vs Opus 4.5 and found the GPT
  cooperation regime did **not** cleanly replicate to Opus 4.5 (partial
  blackboard behavior, no retained signal discovery, no accepted solve).
  **Caveat:** those Opus runs used the older `discoverable` / `shared-goal`
  prompt rungs from 8.11, **not** this paper's A-D ablation ladder, and N was
  tiny. So it is not yet a clean test of *this paper's* threshold — but it is a
  real warning that the threshold may be model-dependent.

This negative signal is the single most important input to the strategy: do not
assume "cross-model" means "confirm generality." It may not generalize.

## Prerequisite gate (blocks all batch runs)

Do not reopen batch runs until a paper-runtime compatibility audit has been
done and recorded against the original April paper contract. The failure mode
from the deleted June wave was not just a bad sample; it was starting reruns
before re-checking whether the newer runtime still matched the paper-defined
agent-visible surfaces.

At minimum the audit must compare:

- task prompt
- runtime prompt / context disclosure
- tool inventory and tool `--help` surfaces
- default workspace artifacts and directory naming
- agent-visible env and capability metadata
- telemetry, counters, timing, and other sidechannels
- any authority-bearing carrier that downstream tools still trust

Canonical rule:
[`../../../../experiments/PROTOCOL.md`](../../../../experiments/PROTOCOL.md).

## Reframed question (because of the 8.15 negative)

Do not frame this as "confirm the threshold generalizes." Frame it as the open
question that is robust to either outcome:

> **Is the minimal-information closure threshold model-dependent — does agent
> capability / family move where coordination breaks (the C→D boundary)?**

- Reproduces across families → robustness; broaden the claim.
- Moves (a stronger model closes even in pure-zero D, or a weaker one fails
  earlier) → the threshold is capability-gated, itself a finding about what
  coordination-under-unawareness requires.

Either way the study is publishable and honest. This framing also turns the
existing Opus negative from an embarrassment into a data point.

## Model set

| Model | Role | Note |
|-------|------|------|
| GPT-5.4 | reference | full A-D ladder already in hand |
| Opus 4.5 | re-run under A-D | existing 8.15 negative used the wrong rungs; needs a clean paper-aligned test before any "doesn't generalize" claim |
| DeepSeek (reasoning) | new family | **open weights** → also answers reproducibility |
| Kimi K2 (Moonshot) | new family | strong agentic / long-context |

Pilot 1-2 Condition-A runs per new model first (capability + harness sanity)
before committing N. Optional 4th family only if resources allow.

## Conditions & run counts (target the threshold, not a dense matrix)

The novel claim is the C→D boundary, so the minimal cross-model unit is the pair
{C succeeds, D fails}, plus A as a capability baseline.

| Model | A | B | C (budget-only) | D (pure zero) |
|-------|---|---|-----------------|---------------|
| GPT-5.4 | have | have | have | have |
| Opus 4.5 (re-run, full) | ~4 | ~4 | ~6 | ~8 |
| DeepSeek (full) | ~4 | ~4 | ~6 | ~8 |
| Kimi (threshold + baseline) | ~3 | skip | ~6 | ~8 |

- Priority under tight resources: **D > C > A > B**. D gets the most runs
  (negatives need more to convince).
- Lean fallback: A/C/D only on every new model.

## Comparability (hold everything but the model fixed)

- Budget 17, 15 cells, 2 agents, per-agent limit 10, identical prompts,
  `direct` workspace, `hidden` runtime.
- Reasoning effort is not 1:1 across vendors — use each model's high/max
  reasoning mode and document the per-model effort mapping as a stated caveat.
- Keep the `world --help` wording socially clean and identical across models.
- Treat prompt/tool/runtime parity as a first-class comparability condition, not
  as an informal setup detail.

## How it absorbs into the paper

1. **Results**: extend Table 2 into a model × condition matrix; add a
   "Cross-Model Replication / Threshold Dependence" subsection.
2. **Claim language**: if it reproduces, reverse the item-5 scope-narrowing
   (abstract / intro / related work / conclusion) toward "across N families"; if
   it moves, report the capability-gated threshold honestly.
3. **Closed-source → strength**: with open-weight DeepSeek in the set, state that
   at least one open model supports reproducibility — resolves Reviewer 1 without
   a rationale sentence.
4. **Limitations**: shift from "single model" to remaining limits (N per cell,
   one task family, informal coding).
5. **Page budget**: the conference limit is 8 body pages with little headroom —
   a multi-model table/subsection means trimming (compress R10 prose / move
   detail to appendix). An extended version removes this constraint.

## Bookkeeping

- Update `8.15-cross-model-replication` to the A-D ladder design (its current
  README predates this paper's ablation framing and carries the older
  GPT-vs-Opus negative). New runs land there.
- Single-model larger-N can reuse the 8.17-8.20 dirs.
- Update `experiment-data.yaml` and the dossier Evidence Anchor / Backbone when
  runs complete.

## Timeline decision

Recommendation: **route the full cross-model study to an extended version**, not
the `2026-06-21` camera-ready, because (a) ~2 weeks is too short for
runs + analysis + rewrite, (b) the compatibility audit still has to pass before
new evidence is countable, (c) the page
budget is tight, and (d) the existing negative Opus signal means the outcome is
genuinely uncertain and should not be rushed into the archival proceedings
version. See the dossier "Current Read" for the strategy rationale.
