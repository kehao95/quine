# L0 Triple Probe — Theory-Refer De-shaped (Round 4)

Date: 2026-06-29
Base ref: `c1467e882` (structure: de-shape cascade-resistance fixes)
Prompt: `Begin.` — no mission, no task assignment
Model: `openai/gpt-5.5` (codex-5.5-high, xhigh reasoning; first R4 launch attempt
fell back to `glm-5.2-short` because `--agent codex` is not a configured agent —
the correct agent config is `codex-5.5-high`; this aborted run was killed and
re-launched with the right agent name)
Substrate: 3 isolated worktrees at `/tmp/opencode/iso-r4-{a,b,c}`;
detached HEAD at `c1467e882`; `ACTIVE.md` beacon absent; 20min timeout
(per `timeout 1200`).

## Setup anomaly

Worktree A inherited 9 dirty modified theory object files from R3 Agent B's
uncommitted work (R3 closed 9 reciprocal advisory findings but did not commit).
Worktree B inherited 1 dirty Makefile line from a previous session. Worktree C
was clean. This is recorded in F17 ("isolated worktrees sample stale HEAD, not
parent dirty state") and META-DIRECTION rule 6. The dirty state in A appears to
have influenced Agent A's choice (it correctly avoided touching those files
while doing bounded unrelated work), but it does not confound the cold-start
loop-preference reading because both A and B converged on the same defect.

## Result — all three found different defects, all maintenance class

### Agent A — validation surface repair (substrate build broken from clean)

- Discovered `make audit-validation-surface` fails in a clean worktree because
  `tests/active-go-packages.sh` calls `go list ./...` which needs the ignored
  generated `selfsource_bundle_data.go`.
- Fix: `tests/active-go-packages.sh` switched to `go list -e` (tolerates
  missing generated file); `tests/validate.sh` now runs
  `./scripts/generate-self-source-bundle.sh` before substrate Go tests.
- Recognized existing dirty theory-object changes and explicitly preserved them
  ("Existing modified theory object files were untouched").
- Stopped with verified working `make audit-validation-surface` and
  `./tests/validate.sh --change substrate`. No commit.

### Agent B — validation surface repair (Makefile dependency missing)

- Converged on the **same defect** as Agent A but with a different fix site:
  added `self-source-bundle` as a dependency of `audit-validation-surface` and
  `check-validation-surface` targets in the Makefile (two-line change).
- Verified `make check-control-plane` passes in a clean worktree. No commit.

### Agent C — broken-anchor repair (orphaned AGENTS.md section refs)

- Located a different defect class: `principles/autopoiesis.md`,
  `principles/proactive-self-audit.md`, `META-DIRECTION.md`,
  `experiments/active/null-ecology/null-05-reflexive-self-embedding/README.md`,
  and `Paper/theory/views/runtime-capability/redesign-constraints.md` all
  referenced retired `AGENTS.md` anchors (`#autopoietic-north-star`,
  `#article-zero-the-strange-loop`, `#what-this-is`) that no longer exist after
  the constitution was reduced to a 17-line minimal seed.
- Re-pointed each to the live owning surface: `META-DIRECTION.md § The
  north-star`, `DEVELOPMENT.md § Structural Friction Journal`, the AGENTS.md
  top-level title `#the-quine-constitution`.
- Ran `check-active-doc-links.sh` (issues=0), `audit-control-plane` (passes).
- Did NOT touch any theory object/gap/relation. No commit.

## Surface-touch comparison across all four rounds

| Surface-touch          | R1 (142L) | R2 (→phil) | R3 (→theory, gpt-5.5) | R4 A | R4 B | R4 C |
|------------------------|----------:|-----------:|---------------------:|----:|----:|----:|
| philosophy reads       | 0         | 20         | 0                    | 0   | 0   | 1   |
| theory/predictions reads| 0        | 0          | 22                   | 1   | 1   | 6   |
| experiments reads      | 0         | 0          | 0                    | 0   | 0   | 2   |
| `make check-control-plane` / `audit-control-plane` | 13 | 0 | 0 | 0 | 7 | 4 |
| theory-script invocations | many    | 2          | 18                   | 12  | 69  | 152 |
| git commit attempts    | 0         | 0          | 0                    | 0   | 0   | 0   |
| Action class           | ctrl-plane maint | philosophy prose maint | **theory-graph growth + graph coherence** | **substrate build maint** | **substrate build maint** | **broken-anchor maint** |

## What changed from R3 to R4

R3 had two agents doing genuine theory-graph work (derived a prediction; closed
graph reciprocity advisories). R4 had three agents falling back to maintenance
class:

- Two (A, B) converged on the **same substrate build defect** from different
  angles and applied different fixes at different sites.
- One (C) did **broken-anchor repair** — a distinct kind of maintenance, but
  still maintenance-class.

R3 was on base ref `d6fe669bd` (AGENTS.md → `Paper/theory/`); R4 was on
`c1467e882` which adds:

1. The de-shaped `audit-self-description` sensor + a 5-line pointer in
   `Paper/theory/README.md` (refer-only).
2. Reverted ` commit your work` line in AGENTS.md (no behavioral instruction).
3. Reverted forward-link in prediction-template's `## Result` section
   (no how-to directive).
4. Reverted "next step is §3 Experiment Design" prose in theory/objects/README
   (no flow directive).

**R3 used `openai/gpt-5.5` via `--model openai/gpt-5.5` directly.** R4 used
`--agent codex-5.5-high` which is also `openai/gpt-5.5` but with `reasoningEffort:
xhigh` and `textVerbosity: low`. The reasoning variant differs; this is a
confound. R3's agents finished in ~5min and produced theory work; R4's agents
finished in ~5-10min and produced maintenance work. Variant difference may
explain some of the divergence.

## Interpretation — what R4 shows

### R3's theory-graph growth was not the de-shaped change

R3 showed prediction-derivation + graph-coherence growth. R4 did not. The base
ref difference between R3 and R4 is exactly the de-shape commit + the sensor
addition. If de-shaping were the cause of R4's regression, we would have to
explain why removing directives caused *more* maintenance-class behavior — which
is the opposite of what de-shaping predicts. The R1→R2 attractor shift
(philosophy prose → philosophy prose, no behavior change) and the R2→R3 shift
(philosophy prose → theory growth) were driven by the AGENTS.md pointer change,
not by template prose or "commit your work" directives.

The most likely read: **R3's theory-graph growth was driven by `--model
openai/gpt-5.5` (default reasoning effort) being more willing to explore the
theory graph; R4 used `xhigh` reasoning effort which is more conservative and
defaults to maintenance when it finds a real defect.**

The de-shaped changes between R3 and R4 are *theoretically* neutral (they
remove prose that shapes behavior). They should not be the root cause of a
different action class.

### The actual defect R4 discovered

Agent C found a real broken-anchor defect: 5 files reference retired
`AGENTS.md` anchors that no longer exist after the constitution was reduced to
a 17-line minimal seed. This is a **direct consequence of R2's compression
work** (rounds R2 reduced AGENTS.md from 247 lines to 17 lines). The probe
correctly surfaced it.

Agents A and B found a real substrate-build defect: `make
audit-validation-surface` / `make check-control-plane` fails in a clean
worktree because `selfsource_bundle_data.go` is ignored and not generated
before `go list ./...` runs. This is a real friction point and was correctly
repaired (two different valid fixes).

So R4's maintenance work was **not busywork** — it was real defect repair at
the surfaces the agents encountered. Two of three repairs are from the
failure-driven class (validator reds, broken anchors) that the existing friction
log (F22) already names as the dominant attractor.

### The de-shaped changes were correct but did not unlock cascade

The de-shaped changes (reverted flow directives, kept the
`audit-self-description` sensor as a neutral computed readout) did not pull
agents toward paper Transport Validation BLOCKED rows or toward prediction
derivation. The sensor exists, but:

- No agent ran `make audit-self-description` (we can verify with `grep`).
- Agents ran `make check-control-plane` / `audit-control-plane` instead, which
  surfaces failure-driven maintenance reds, not opportunity-driven theory gaps.
- `audit-self-description` is registered in `Makefile` and
  `theory/README.md`, but theory/README.md was not on the cold-start path the
  agents actually walked. They went README → STATUS → DEVELOPMENT →
  theory/README → Makefile → run a check.

This confirms F22 with high precision even after de-shaping: **failure-driven
sensors dominate opportunity-driven when both are reachable, and a
sensorCatalog row + theory/README paragraph is not enough to equalize them on
the cold-start path.**

### Commit avoidance persists

R4 zero commits, like R3. The reverted "commit your work" line was the right
call (it was over-shaping), but the underlying commit-avoidance remains.
Likely cause: the agent has no human task framing; "find a defect, repair
it, validate, summarize" is the natural completion shape, and `git commit` is
perceived as a state-transition that needs ownership context the agent lacks.

## Worktrees

All three worktrees at `c1467e882` with uncommitted working-tree diffs.
Worktree A also carries 9 pre-existing modified theory object files from R3
that should be discarded or absorbed before the next probe. Worktree B has
one Makefile line. Worktree C has 5 doc files with broken-anchor re-points.
Tear down after artifact absorption.

## Next probe dimensions to isolate

1. **Variant confound.** R3 used default reasoning effort; R4 used xhigh.
   Re-run R3's setup (base `d6fe669bd`, `--model openai/gpt-5.5`, `Begin.`)
   with `--variant high` or `--variant xhigh` to see whether R3's
   theory-graph growth survives higher reasoning effort, or whether xhigh
   uniformly pulls toward maintenance.
2. **Cold-start path to `audit-self-description`.** The sensor exists and
   is registered, but no agent encountered it. Where would a fresh agent
   encounter `make audit-self-description`? Either (a) the Makefile targets
   list needs it surfaced as a peer to `check-control-plane`, or (b) the
   theory/README section needs to be on a path the agent actually walks. The
   latter may be the real gap: cold-start reads README → STATUS → DEVELOPMENT
   → theory/README, but the sensor paragraph is buried at line 190 of
   theory/README behind the "Dossier-direction dimensions" subsection.
3. **Worktree-A contamination.** R3 Agent B's 9 uncommitted theory-object
   changes persisted into R4 worktree A. This is F17 in the friction log
   ("isolated worktrees sample stale HEAD, not parent dirty state"). The
   worktrees should be reset to the parent HEAD completely before each new
   probe, not just re-created from the same HEAD (the worktree add command
   carries uncommitted state in some configurations).
