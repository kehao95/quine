---
surface_kind: experiment
phase: p8-population
experiment_id: coop-20-full-disclosure-baseline
experiment_type: ablation
status: complete
id: coop-20
legacy_id: p8:8.19
family: cooperation-dynamics
theory-objects: [coordination-closure, viability-pressure]
mechanisms: [directive-framing, perceptual-disclosure, process-native-coordination, scarcity-selection-pressure, forced-externalization]
lineage_phase: p8-population
---

# Exp 8.19: Full Disclosure Baseline

> **Interpretation boundary:** 8.19 is the **full-disclosure baseline** for the
> ablation ladder. Agents are explicitly told about peers, shared budget, and
> concurrent execution. This is **Step A** of the ladder.

## Hypothesis

With full disclosure of social context—explicit mention of other agents, shared
budget, and concurrent execution—coordination should succeed reliably. This
establishes the baseline that coordination *can* work with complete information.

## Core Question

> When agents have complete social information (peers, shared budget, concurrent
> execution), does coordination succeed? This provides the control condition for
> the ablation ladder.

## Paper Feeds

- `alife/minimal-perceptual-prerequisites` - primary - main-text - **Step A**: full disclosure baseline

## Surface Map

```text
8.19-full-disclosure-baseline/
├── README.md          ← this file
├── prompts/           ← full-disclosure task prompts
├── setup/             ← inherited guarded-world generator
├── analysis/          ← scores task success
└── runs/              ← per-run artifact directories
```

## Experimental Design

### Baseline Inheritance

8.19 uses the same base infrastructure as 8.17:

- shared workspace with fs_mutations enabled
- 15 cells, shared budget 17, 2 agents
- max 10 `world get` calls per agent per reset epoch
- GPT-5.4 with xhigh reasoning effort
- unlimited turns

### Independent Variable

The independent variable is **social information disclosure**:

| Condition | Social Disclosure | fs_mutations | Budget Visibility |
|-----------|-------------------|--------------|-------------------|
| `full-disclosure` | FULL (peers, shared budget, concurrent) | enabled | enabled |

### What Agents Are Told

Unlike 8.17 (zero-social), agents in 8.19 are explicitly told:

1. "You are one of multiple agents working concurrently"
2. "The budget of 17 calls is shared across all agents"
3. "Other agents may create or modify files in the workspace"
4. "Coordinate with peers to avoid duplicating work"

### Expected Outcome

**Success** — coordination should emerge readily with full information. This
establishes that our task and infrastructure support coordination when agents
have the necessary social context.

## Run Ledger

| Run Directory | ID | Condition | Outcome | Notes |
|---------------|-----|-----------|---------|-------|
| `20260409-145513-copilot-gpt-5.4-xhigh-full-disclosure-condA-2ag` | R01 | ~~full-disclosure~~ | INVALID | runtime was hidden (bug), needs rerun |
| `20260409-151703-copilot-gpt-5.4-xhigh-full-disclosure-condA-2ag` | R02 | full-disclosure | retained, reread pending | actual retained pointer replacing stale `150343` reference |
| `20260409-173947-copilot-gpt-5.4-xhigh-full-disclosure-condA-2ag` | R03 | full-disclosure | retained, reread pending | retained rerun pointer |
| `20260409-174017-copilot-gpt-5.4-xhigh-full-disclosure-condA-2ag` | R04 | full-disclosure | retained, reread pending | retained rerun pointer |

### Current Interpretation Boundary

R01 is invalid as a full-disclosure control because the runtime was hidden. The
three later retained pointers are cataloged, but this README no longer claims a
specific success rate until those raw traces are reread and absorbed into
`RUN-CATALOG.md` / `BEHAVIOR-ANALYSIS.md`.

This keeps 8.19 as the full-disclosure baseline owner without letting a stale
missing path (`150343`) stand in for retained evidence.

## Analysis Protocol

Score each run for:
1. Task completion (15/15 cells collected)
2. Budget efficiency (calls remaining)
3. Coordination mode (explicit protocol, implicit, asymmetric)
4. Time to first coordination artifact
