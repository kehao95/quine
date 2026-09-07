# L0 Ablation 2x2 — 2026-06-29

## Question

R3 showed theory-graph discovery (prediction derivation + reciprocity closure).
R4 regressed to substrate maintenance. Five confounds identified:

1. **Worktree contamination** — R4 Agent A inherited 9 dirty files from R3 Agent B
2. **Variant** — R3 used `--model openai/gpt-5.5` (default effort); R4 used `--agent codex-5.5-high` (xhigh + textVerbosity: low)
3. **Cold-start path** — R3 went theory-first; R4 went root-first (broader recon)
4. **Git log attractor** — R4 agents saw maintenance-focused recent commits
5. **Failure-driven attractor** — R4 agents hit real `selfsource_bundle_data.go` defect

De-shaping (R3→R4 base diff) is theoretically neutral and not expected to cause regression.

## Design

2x2 factorial, 2 agents per cell, 8 agents total. All worktrees clean.

Prompt: `Begin.`

| Cell | Base commit | Variant | N | Tests |
|---|---|---|---|---|
| A | `d6fe669bd` (R3 base) | `--model openai/gpt-5.5` (default) | 2 | R3 replication (clean worktree control) |
| B | `d6fe669bd` (R3 base) | `--agent codex-5.5-high` (xhigh) | 2 | Does xhigh alone cause maintenance pull? |
| C | `c1467e882` (R4 base) | `--model openai/gpt-5.5` (default) | 2 | Does de-shape + default still discover? |
| D | `c1467e882` (R4 base) | `--agent codex-5.5-high` (xhigh) | 2 | R4 replication (clean worktree control) |

## Predictions

- If **variant is the cause**: B does maintenance, C does discovery (same as A)
- If **base/de-shape is the cause**: C does maintenance, B does discovery (same as A)
- If **both matter**: B and C both do maintenance
- If **neither matters** (was just contamination/noise): B does discovery, C does discovery

## Worktrees

All at `/tmp/opencode/ablation-{cell}{n}` (e.g. `/tmp/opencode/ablation-a1`).

## Artifact naming

- `agent-{cell}{n}-stdout.log` — full session transcript
- `agent-{cell}{n}-diff.patch` — working tree diff
- `agent-{cell}{n}-status.txt` — git status --short
