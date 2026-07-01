---
surface_kind: experiment
phase: p8-population
experiment_id: coop-19-fs-mutations-disabled
experiment_type: ablation
status: complete
id: coop-19
legacy_id: p8:8.18
family: cooperation-dynamics
theory-objects: [stigmergy-as-carrier, trigger-structure, coordination-closure]
mechanisms: [perceptual-disclosure, scarcity-selection-pressure, process-native-coordination, forced-externalization, carrier-mediated-inheritance]
lineage_phase: p8-population
---

# Exp 8.18: fs_mutations Disabled

> **Interpretation boundary:** 8.18 tests coordination when `fs_mutations` is disabled but budget visibility remains.
> fs_mutations telemetry is disabled, so agents cannot observe unexplained
> file changes. This is **Step D** of the ladder — expected to FAIL.

## Hypothesis

Without passive perception (fs_mutations disabled), agents cannot detect
environmental anomalies that would trigger coordination. Even with budget
visibility, coordination should collapse because agents cannot observe:

1. Files appearing that they did not create
2. Files being modified by unknown actors
3. Task state changes they did not cause

## Core Question

> When passive perception is disabled, does coordination collapse? This would
> prove that passive perception is a **necessary condition** for spontaneous
> multi-agent coordination in zero-social-information settings.

## Paper Feeds

- `alife/minimal-perceptual-prerequisites` - primary - main-text - **Step D**: perception disabled (expected failure)

## Surface Map

```text
8.18-fs-mutations-disabled/
├── README.md          ← this file
├── prompts/           ← zero-social prompts (same as 8.17)
├── setup/             ← inherited guarded-world generator
├── analysis/          ← scores task failure modes
└── runs/              ← per-run artifact directories
```

## Experimental Design

### Baseline Inheritance

8.18 uses the same base as 8.17, except:

- **fs_mutations DISABLED** — agents do not receive file change notifications
- Everything else identical: budget visibility, clean world help, zero social hints

### Independent Variable

The independent variable is **passive perception availability**:

| Condition | Social Hints | fs_mutations | Budget Visibility |
|-----------|--------------|--------------|-------------------|
| `perception-disabled` | NONE | **DISABLED** | enabled |

### What This Tests

8.17 showed that with passive perception (fs_mutations) + budget visibility,
coordination succeeds 100%.

8.18 removes fs_mutations while keeping budget visibility. If coordination
fails, this proves:

> **Passive perception is necessary.** Budget anomalies alone are insufficient.

### Expected Outcome

**Failure** — coordination should collapse. Possible failure modes:

1. **Budget race**: Both agents proceed independently, exhaust budget
2. **Blind duplication**: Both collect the same cells, waste calls
3. **No coordination artifacts**: No REQUEST.txt or protocol files emerge

### Success Criteria for the Experiment

The experiment "succeeds" (in proving our hypothesis) if:

1. Coordination fails in most/all runs (budget exhausted, incomplete task)
2. No coordination protocols emerge spontaneously
3. Agents show no evidence of peer inference

## Run Ledger

| Run Directory | ID | Condition | Outcome | Notes |
|---------------|-----|-----------|---------|-------|
| 20260409-151633-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R01 | early no-perception | **TERMINATED** | excluded from clean behavior read |
| 20260409-160036-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R02 | perception-disabled (clean) | **FAILURE** | 10/15, single-agent dominance |
| 20260409-164736-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R03 | perception-disabled (clean) | **SUCCESS** | 15/15, budget-anomaly inference |
| 20260409-173817-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R04 | perception-disabled (clean) | **SUCCESS** | 15/15, coordination files created |
| 20260409-173847-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R05 | perception-disabled (clean) | **SUCCESS** | 15/15 |
| 20260409-173917-copilot-gpt-5.4-xhigh-no-perception-condA-2ag | R06 | perception-disabled (clean) | **SUCCESS** | 15/15 |

**Summary**: 4/5 interpreted clean runs succeeded after R01 was excluded.
R02 remains the clean failure case; R03-R06 show that budget-anomaly inference
provides an alternative peer-discovery channel even when `fs_mutations` is
disabled. See transcript analysis below.

### R01 Analysis (Contaminated)

**Note:** R01 used a buggy runtime that showed `fs_mutations: (empty)` instead of
omitting the field entirely. This was misleading but still prevented agents from
seeing actual file changes. R02 uses the fixed runtime.

### R02 Analysis (Clean)

**Result: FAILURE** — coordination collapsed as predicted with clean runtime.

- **Cells collected**: 10/15 (incomplete)
- **Budget remaining**: 6 (sufficient budget, but no coordination)
- **Agent-1**: 1 get (almost completely idle)
- **Agent-2**: 10 gets (hit per-agent limit)
- **Resets**: 2

**Failure mode**: "Single-agent dominance" — Agent-2 independently collected 10 cells
(hitting the per-agent limit), while Agent-1 made only 1 get. Without fs_mutations
telemetry, neither agent could detect the other's activity. Agent-2 couldn't signal
for help, and Agent-1 couldn't perceive work being done.

**R02-local insight**: With fs_mutations completely omitted from tool results
(not just empty), agents have no passive file-mutation channel; in this run,
budget visibility alone was insufficient. Later R03-R06 show the broader branch
is not that simple: budget-anomaly inference can still support closure without
`fs_mutations`.

## Analysis Protocol

Score each run for:
1. Task completion status (success/failure)
2. Budget exhaustion (remaining calls at end)
3. Duplication rate (cells collected by both agents)
4. Coordination artifact presence (REQUEST.txt, etc.)
5. Peer inference traces in reasoning
