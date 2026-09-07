# L0 Ablation 2x2 — CLOSURE (Full, N=20)

## Setup

24 agents launched (6 per cell), 4 agents lost (a6/b6/c6/d6 — killed before
start). 20 valid agents across 4 cells, all clean worktrees, prompt `Begin.`.

| Cell | Base | Variant | N | Sensor present |
|---|---|---|---|---|
| A | `d6fe669bd` (R3 base) | `--model openai/gpt-5.5` (default) | 5 | no |
| B | `d6fe669bd` | `--agent codex-5.5-high` (xhigh) | 5 | no |
| C | `c1467e882` (R4 base) | `--model openai/gpt-5.5` (default) | 5 | yes |
| D | `c1467e882` | `--agent codex-5.5-high` (xhigh) | 5 | yes |

## Full agent table

| agent | sensor | ctrlPL | thObj | substr | exper | diffln | behavior |
|---|---:|---:|---:|---:|---:|---:|---|
| a1 | 0 | 0 | 117 | 0 | 2 | 64 | seed-backlog clock maintenance |
| a2 | 0 | 0 | 81 | 0 | 1 | 64 | seed-backlog clock maintenance |
| a3 | 0 | 15 | 4 | 49 | 134 | 24 | substrate defect repair |
| a4 | 0 | 0 | 45 | 0 | 0 | 61 | theory reciprocity closure (5 objects) |
| a5 | 0 | 7 | 1 | 32 | 26 | 10 | substrate defect repair |
| b1 | 0 | 9 | 1 | 28 | 24 | 10 | substrate defect repair |
| b2 | 0 | 14 | 7 | 57 | 46 | 44 | substrate defect repair (3 files) |
| b3 | 0 | 9 | 9 | 36 | 46 | 34 | substrate defect repair |
| b4 | 0 | 1 | 13 | 0 | 0 | 69 | governance/docs (Makefile+scripts/README) |
| b5 | 0 | 9 | 6 | 44 | 26 | 11 | substrate defect repair |
| c1 | 2 | 0 | 60 | 0 | 0 | 61 | sensor-driven reciprocity closure (5 objects) |
| c2 | 5 | 0 | 19 | 0 | 0 | 22 | sensor → Makefile validator gap |
| c3 | 0 | 7 | 1 | 38 | 24 | 11 | substrate defect repair |
| c4 | 2 | 7 | 21 | 47 | 26 | 24 | substrate defect repair (+ sensor invoked) |
| c5 | 0 | 0 | 43 | 0 | 0 | 61 | theory reciprocity closure (5 objects, via audit-theory-family, NOT sensor) |
| d1 | 0 | 0 | 56 | 0 | 0 | 110 | theory reciprocity closure (9 objects) |
| d2 | 0 | 0 | 1 | 6 | 43 | 96 | experiment runtime fix (FUSE mountpoint) |
| d3 | 0 | 15 | 1 | 36 | 21 | 10 | substrate defect repair |
| d4 | 0 | 7 | 1 | 46 | 29 | 41 | substrate defect repair |
| d5 | 0 | 0 | 0 | 0 | 21 | 350 | experiment FUSE overlay error handling (via Linear issue tracker) |

## Behavior distribution

| Behavior | Count | Pct | Cells |
|---|---:|---:|---|
| Substrate defect repair (selfsource/active-go-packages) | 10 | 50% | A(2), B(4), C(2), D(2) |
| Theory reciprocity closure (broken papers-using pointers) | 4 | 20% | A(1), C(2), D(1) |
| Experiment/runtime fix | 2 | 10% | D(2) |
| Seed-backlog clock maintenance | 2 | 10% | A(2) |
| Governance/docs | 1 | 5% | B(1) |
| Sensor-driven Makefile gap | 1 | 5% | C(1) |

## Per-cell behavior diversity

| Cell | Distinct behaviors | Behaviors |
|---|---:|---|
| A (default, no sensor) | 3 | seed-backlog×2, substrate×2, reciprocity×1 |
| B (xhigh, no sensor) | 2 | substrate×4, governance×1 |
| C (default, +sensor) | 3 | reciprocity×3, substrate×1, Makefile×1 |
| D (xhigh, +sensor) | 3 | substrate×2, experiment×2, reciprocity×1 |

## Findings

### 1. Substrate defect is the dominant attractor (50%, cross-config)

10/20 agents were pulled to repair the `selfsource_bundle_data.go` generation
defect. This attractor is configuration-independent — any agent that runs
`make check-control-plane`, `audit-theory-family`, or related theory validators
encounters the failure and redirects to fix it. This is a real broken state,
not a repo design directive.

### 2. Randomness is real and healthy

Same-cell agents diverge into 2-3 distinct behaviors in every cell. No cell
shows extreme convergence (>80% to one behavior). Compare R1 (13/13 to
check-control-plane). The repo does not over-direct agent behavior.

### 3. Theory reciprocity is a stable cross-config discovery attractor

4 agents (a4, c1, c5, d1) independently closed broken papers-using-this-object
pointers in theory objects. This is not sensor-dependent:
- a4: discovered via `check-theory-object-refs` (Cell A, no sensor)
- c1: discovered via `audit-self-description` sensor (Cell C)
- c5: discovered via `audit-theory-family` (Cell C, did NOT invoke sensor)
- d1: discovered via theory-maturity/coherence scripts (Cell D, xhigh)

The reciprocity gap is a structural property of the theory graph (objects
missing cross-references to papers that use them), detectable through multiple
independent paths.

### 4. Sensor reach is effort-dependent

Sensor invocation counts:
- Cell C (default + sensor): c1=yes, c2=yes, c3=no, c4=yes, c5=no → 3/5
- Cell D (xhigh + sensor): 0/5 (none reached it)

Default-effort agents reach the sensor via theory/README 60% of the time.
xhigh agents never reach it (root-first recon bypasses the paragraph at
theory/README line ~190).

### 5. d5: unique external-task-source cold start

d5 is the only agent (1/20) that went to an external task tracker (Linear) for
cold-start direction. Found Q-12/Q-119, then implemented FUSE overlay error
handling. This path appeared only under xhigh + sensor base.

### 6. All agents terminate after one fix — no exploration continues past closure

Every agent (20/20) follows the pattern: find a gap → fix → verify → stop.
No agent performed maintenance AND discovery in the same session. No agent
fixed one gap then continued to explore others.

### 7. Binary pass/fail checks create stronger trajectory termination than advisory checks

Trajectory depth varies by check type:

| Check type | Example | Exit | Trajectory | Depth |
|---|---|---|---|---|
| Blocking (strict) | `check-control-plane` fails | 1 | hard failure → urgency → mechanical fix → pass → "done" signal → stop | shallow: entire session on one fix |
| Advisory (exit 0) | `audit-self-description`, `growth-signal-readout` | 0 | surface gap → agent judges → creative fix → verify | medium: reads 5-9 theory objects |
| Experiment runtime | FUSE mountpoint failure | — | run → fail → diagnose → fix → rerun | deep: still in progress when truncated |

Key mechanism: **strict checks (exit 1) create a "completion signal"**. When the
check passes after a fix, the agent interprets the green checkmark as "done"
and terminates. Advisory checks (exit 0) surface gaps without blocking, so
no strong "done" signal terminates the trajectory — the agent must self-judge
when to stop.

This explains why R3's prediction derivation was unique: R3 agents hit no
blocking check (growth-signal-readout is advisory), so their trajectory
continued naturally to the creative step of deriving a prediction. Agents
that hit blocking checks all consumed their turns in fix→verify→stop cycles.

### 8. Code checks vs language evaluation: local horizon vs divergence

Code-based checks (Makefile targets, shell validators) produce binary
pass/fail signals that:
- Create urgency (failure must be fixed)
- Create completion satisfaction (green checkmark = done)
- Focus agent on **local horizon** (the failing check and its immediate fix)
- Terminate exploration after closure

Language-based evaluation (reading theory objects, graph coherence) produces
ambiguous signals that:
- Carry no urgency (advisory, exit 0)
- Carry no completion signal (no green checkmark to achieve)
- Allow **divergent** interpretation (agent decides what "done" means)
- Keep exploration open because there is no binary closure point

The repo's validator infrastructure (Makefile targets with --strict) is a
powerful attractor that funnels agents into local optimization. The theory
graph (objects, gaps, predictions) is a softer surface that allows broader
exploration but produces stochastically varied results.
