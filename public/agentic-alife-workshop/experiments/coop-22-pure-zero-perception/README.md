---
surface_kind: experiment
phase: p8-population
experiment_id: coop-22-pure-zero-perception
experiment_type: ablation
status: complete
id: coop-22
legacy_id: p8:8.20-pure-zero-perception
family: cooperation-dynamics
theory-objects: [stigmergy-as-carrier, trigger-structure, coordination-closure]
mechanisms: [perceptual-disclosure, scarcity-selection-pressure, process-native-coordination, forced-externalization, carrier-mediated-inheritance]
lineage_phase: p8-population
---

# Exp 8.20: Pure Zero Perception

> **Interpretation boundary:** 8.20 is the **pure zero perception** condition.
> Both `fs_mutations` AND budget visibility are disabled. This is the
> strictest ablation — **Step D** of the ladder.

## Hypothesis

With all known perception channels disabled:
1. `fs_mutations` telemetry — DISABLED (cannot see file changes)
2. Budget visibility — HIDDEN (cannot detect budget anomalies)

Coordination should fail completely. Agents have no channel to infer
peer existence. This would establish that *some* perception channel
is necessary for spontaneous coordination.

## Motivation

8.18 (perception-disabled) showed 4/6 success rate because agents could
still infer peers from budget anomalies:

> "budget dropped from 17 to 14 after my first call—I expected 16"

8.20 removes this channel by hiding budget information entirely.
The `world get` output changes from:

```
testvalue1
[generation: 1] [budget: 14/17 remaining]
```

to:

```
testvalue1
[generation: 1]
```

## Core Question

> When ALL perception channels are disabled, does coordination become
> impossible? This would prove that passive perception through *some*
> environmental channel is necessary for spontaneous coordination.

## Paper Feeds

- `alife/minimal-perceptual-prerequisites` - primary - main-text - **Step D**: pure zero perception (expected failure)
- `alife/agentic-alife-workshop` - secondary - supporting-only - pure zero-perception probe for the workshop paper's Assay 2; clarifies artifact-aware coordination versus writeback / validation closure when `fs_mutations` and budget visibility are absent

## Experimental Design

### Environment Variables

```bash
QUINE_FS_MUTATION_TELEMETRY_ENABLED=0   # No file change telemetry
QUINE_PROMPT_BUDGET_VISIBILITY=hidden   # No budget information
```

### What This Tests

If 8.20 fails while 8.18 succeeds, this proves:
- Budget-anomaly inference was the mechanism enabling 8.18 success
- Peer-inference requires at least one environmental signal

If 8.20 somehow succeeds, this suggests:
- Agents can coordinate purely through file content (e.g., checking for existing files)
- Or our channel isolation is still incomplete

### Expected Outcome

**FAILURE** — coordination should collapse completely.

## Run Ledger

### Wave 1 (2026-04-09) — Original v1 prompt

| Run ID | Turns | Resets | Time | TaskDone | Peer-Aware | Coordination | Violations |
|--------|-------|--------|------|----------|------------|--------------|------------|
| 201406 | 320 | 5 | 83m | no | yes | unilateral | clean |
| 202112 | 174 | 7 | 44m | no | yes | unilateral | clean |
| 211311 | 361 | 17 | 168m | YES* | yes | unilateral | ⚠️ source code |
| 211354 | 181 | 2 | 45m | YES* | yes | none | ⚠️ state.json |

### Wave 2 (2026-04-11) — Hardened prompt (no binary inspection)

| Run ID | Turns | Resets | Time | TaskDone | Peer-Aware | Coordination | Violations |
|--------|-------|--------|------|----------|------------|--------------|------------|
| 025550 | 142 | 5 | 22m | no | yes | none | ⚠️ sibling/runtime listing |
| 025552 | 226 | 5 | 45m | no | yes | none | clean |
| 025554 | 250 | 5 | 52m | YES* | yes | **bilateral** | ⚠️ world/state files |
| 025556 | 314 | 7 | 54m | no | yes | **bilateral** | ⚠️ state + source |

*TaskDone=YES means `world validate` passed, and all YES cases occurred after violations.*

### Wave 3 (2026-05-26) — Current-model finite-time probe

This wave was run as a sidechannel-removal check after the workshop paper
was recalibrated as an exploratory assay-substrate proposal. It uses the same
pure-zero perception condition but a current `gpt-5.5` xhigh lane rather than
the earlier `gpt-5.4` xhigh lane. It is therefore a supplementary mechanism
probe, not a replacement for the retained 2026-04-09 / 2026-04-11 backbone.

| ID | Run ID | Model lane | Stop mode | Event window | TaskDone | Coordination | Interpretation |
|----|--------|------------|-----------|--------------|----------|--------------|----------------|
| `P8-8.20P-R09` | `20260526-014726-gpt-5.5-codex-xhigh-pure-zero-condA-2ag` | `gpt-5.5` xhigh | operator-stopped finite-time probe | `01:47:50`-`02:22:58Z` | no, `world validate` rejected twice | delayed artifact-aware handoff, no completed writeback | shared files still enabled delayed coordination, but the cohort did not close the feedback loop before cutoff |

An earlier same-day attempt
`20260526-014651-gpt-5.5-codex-xhigh-pure-zero-condA-2ag` is excluded as
runner drift: the current runtime required `QUINE_WORKSPACE_OWNER=1`, while the
older P8 runner still passed `true`.

#### 20260526 finite-time read

The run is useful because it separates three outcomes that the earlier summary
partly collapsed:

1. **No early passive trigger.** With both `fs_mutations` and budget visibility
   disabled, the two agents initially duplicated work on `c01`-`c09`, exhausted
   the shared budget, and did not establish early closed-loop coordination.
2. **Artifact awareness can still happen.** After a reset, agent-1 collected
   `c01`-`c10`, hit the per-agent limit on `c11`, and wrote `COORDINATION.md`
   asking the other agent to collect `c11`-`c15`. Agent-2 later collected
   `c11`-`c15`, showing that shared file content can still become a delayed
   coordination surface even without mutation telemetry or visible budget
   anomalies.
3. **Closure still failed in finite time.** At cutoff, `results.txt` contained
   only `c01`-`c10`; no `peer_values.txt` existed; and `world validate` was
   rejected twice. The world state shows all `c01`-`c15` were collected in
   generation 3, but the values for `c11`-`c15` never made it back into the
   submitted workspace artifact.

Finite-time interpretation: pure zero perception does not make coordination
impossible in the absolute sense. It blocks clean task closure in the current
samples by removing the temporal feedback channels that make coordination
converge. Durable file content can still support delayed artifact-aware
coordination, but without `fs_mutations` or budget visibility the agents failed
to complete the writeback / validation loop before operator cutoff.

**Excluded runs:**
- 201346: broken workspace
- 034754: anomalously short (6 min)
- 034756: infrastructure crash
- 034758: anomalously short (6 min)

### Coordination Classification

- **None**: No shared scripts or coordination artifacts
- **Unilateral**: One agent created scripts, only that agent used them
- **Bilateral**: Both agents created or used shared scripts for division of labor

**Bilateral coordination observed in 2 runs (025554, 025556):**
- 025554: agent-2 created `collect.sh` and `next.sh`; both agents used them
- 025556: Both agents created and used `do_get.sh`, `do_reset.sh`

However, both bilateral runs later became contaminated by out-of-workspace inspection.

### Summary

| Metric | Count |
|--------|-------|
| Total runs | 8 |
| Task completed (any means) | 3/8 |
| Task completed legitimately | **0/8** |
| Violation-contaminated runs | 5/8 |
| Peer awareness | 8/8 (100%) |
| Bilateral coordination | 2/8 (25%) |
| Unilateral scripts | 3/8 (38%) |
| Clean runs | 3/8 |
| Mean turns (clean-only, lower bound) | 240 |
| Mean time (clean-only, lower bound) | 57.3 min |
| Mean resets (clean-only, lower bound) | 5.7 |

### Key Finding

Step D **blocks clean task closure** under compliant conditions in the retained
samples, but it does not block every form of artifact-aware coordination.

All three non-violating runs failed despite peer-awareness. Bilateral coordination appeared only in later-contaminated runs:
- 025554: Both agents used `collect.sh` + `next.sh` for systematic collection
- 025556: Both agents used `do_get.sh` + `do_reset.sh` for coordinated resets

The 2026-05-26 finite-time probe adds a cleaner boundary: a compliant current
model run produced delayed `COORDINATION.md` uptake and tail-cell collection,
but still failed to convert artifact awareness into legitimate closure.

**The missing channel**: Without budget visibility or fs_mutations, agents cannot detect *when* their peer acts—only *that* artifacts exist. This temporal blindness prevents the closed-loop feedback needed for convergent coordination.
