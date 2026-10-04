# ICLR upgrade experiments — Tier 1 (scale) + Tier 2 (steering) (HANDOFF SPEC)

**Status:** specified 2026-07-06, not started. Self-contained handoff for an
executing agent. Rationale: `../manuscript-plan.md` § ICLR-readiness. These
two tiers decide the venue (Tier 1+2 delivered → ICLR 2027, predicted
deadline 2026-09-24; otherwise TMLR with current scope). Report results back;
**do not edit the paper dossier** — the writing session absorbs them.

Binding constraints inherited from the elic-34 retraction (apply to ALL new
fixtures below): (1) **event-level DV** — snapshot the workspace file list
after every tool call (Quine harness: score from tapes/final+snapshots;
mini-swe driver: `SnapshottingDockerEnvironment` already does this); a run
counts if the matched act is observed at ANY point (cleanup-proof). (2) No
accidental semantic leakage in scaffolding meant to be neutral: neutral
subtree names come from `_lib/neutral_names.sh`, never English words like
`cache`/`state`/`old`. Content INSIDE a feature is allowed to carry meaning
— that is the rescoped claim's scope; only *scaffolding* must stay neutral.

## Tier 1 — scale the five-motif spine (survey table)

**Goal:** turn Table 1 into a systematic two-harness survey: n≥20/cell,
≥5 model families, five motifs.

1. **Quine harness (existing experiments, just bigger campaigns).** For each
   of `elic-01-gap-completion`, `elic-03-rupture`,
   `elic-04-affordance-solicitation`, `elic-06-topological-gap`,
   `elic-07-semantic-collision` under
   `experiments/active/structural-elicitation/`: run `./run.sh <env> 20 240
   --jobs 5` per family, with `RUNS_DIR` per family per the experiment's
   README conventions. Families, in priority order: `profiles/gpt-5.5-codex-oauth-medium.env`,
   repo-root `profiles/glm-5.2-zai-medium.env`, `profiles/gpt-5.4-codex-oauth-medium.env`,
   `profiles/kimi-k2.6-litellm-medium.env` (probe the litellm gateway first;
   it was unhealthy on 2026-07-05), `profiles/gemini-3.5-flash.env` (probe),
   family-local `profiles/deepseek-v4-pro.env` (only if the user confirms the
   account is topped up — it 402'd on 2026-07-05). Pre-registered floor for
   the survey table: **≥5 families × ≥20/cell on all five motifs**; report
   partial coverage honestly if a lane fails.
   - elic-07 note: keep its two-direction deconfound arms (both collision
     directions), per its README.
   - Existing pilot runs stay in their run dirs; new campaigns go to fresh
     `RUNS_DIR`s (e.g. `runs-<family>-scale/`) so the paper can cite the
     scaled campaign cleanly.
2. **mini-swe-agent harness (port two more motifs).** Extend
   `experiments/active/structural-elicitation/elic-35-harness-generality/`
   with `rupture` and `affordance` fixtures (same seeding logic as the Quine
   originals; the affordance binary can be compiled with the same
   `assets/affordance.go` pattern). n=10/cell × {glm-5.2, + one more family
   via litellm if up}, generic-prime only (elic-35 already showed prime
   independence). DV per original motif definitions, event-level.

**Scoring:** existing per-experiment `analysis/score.py` for Quine runs
(they already implement the per-motif DV); elic-35's scorer pattern for the
mini-swe ports. Deliverable: one summary table (motif × family × harness →
structure-rate / control-rate / n).

## Tier 2 — steering and competition (the control-surface result)

**New experiment directory** `elic-3X-selection-steering/` (next free
number), Quine harness, forked from elic-01's runner skeleton.

**(a) Selection steering (subsumes the old L2 design).** Workspace contains
**three identical stations** — subdirs with neutral 6-hex names
(`neutral_names.sh 3 <seed>`), each seeded with the SAME inert material
(e.g. three data files, same bytes across stations). Exactly one station
additionally carries ONE feature (a gap in an ordinal file set inside that
station). Four conditions: feature-in-station-1/2/3 + flat control (no
feature anywhere). Station order/naming counterbalanced (rotate which hex
name sits first in directory listings across replicates — position is a
known cue, lesson of elic-34). n=10/cell × 2 families (gpt-5.5, glm-5.2).
- **DV:** in which station does a completion act occur (event-level).
  Because stations are identical minus the feature, every station's act is
  equally *available* by construction — this IS the action-space-invariant
  demonstration (attraction, not constraint).
- **Pre-registered rule:** steering accuracy = fraction of feature-arm runs
  whose first station-act lands in the featured station; **PASS ≥ 0.8**
  with flat-control station-acts ≈ 0. PASS ⇒ the paper's headline gains
  "which act a frozen agent performs can be selected ≥80% by workspace
  composition alone, no prompt change."

**(b) Feature competition (dominance ordering).** Workspace seeds TWO
features from the five, in separate neutral-named subtrees, locations
counterbalanced. Priority pairs (5 of 10, cheapest information first):
gap–affordance, gap–rupture, rupture–affordance, topology–gap,
semantic–affordance. n=10/pair × 2 families.
- **DV:** which matched act occurs FIRST (event-level ordering from
  per-call snapshots); record first-act / second-act / both / neither.
- **Pre-registered rule:** a dominance relation A≻B exists in a cell if
  P(A's act before B's) ≥ 0.7 among runs where at least one act fired.
  Report the resulting order and whether it is transitive; no pass/fail
  gate — this is a characterization result either way.

## Budget

Tier 1: ~1,000 Quine runs + ~120 mini-swe runs; Tier 2: ~180 runs. At
--jobs 5, ≤4 min/run: roughly 3–5 days of wall-clock spread across lanes;
GLM/kimi/gemini are cheap lanes, gpt lanes ride the Codex subscription.
Disk: watch `df` (the 2026-07-05 campaigns ran at 81% root usage); prune
old `runs-*` trees via the family root's `_ds_safe_rerun.sh`-style guards
if needed.

## Reporting

Per tier: paste scorer summary tables into the relevant experiment READMEs
under `## Result`, then report back a compact summary (rates only). The
writing session absorbs into manuscript §4 (survey), new §5.4
(steering/dominance), and the abstract headline.

## Status (2026-07-06 — Tier 1 Quine-harness core DONE and audited)

**Executed:** the five-motif spine at n=20/cell on three families — GLM-5.2,
gpt-5.4-codex-medium, and **Claude Sonnet 5** (a fourth family beyond the
original three; gpt-5.5 was not run — the executing agent substituted
gpt-5.4 + Claude). Partial kimi/gemini probes exist on disk (3–10 runs each,
uncommitted) — lanes were unhealthy; treat as non-evidentiary until rerun.
DeepSeek still blocked (HTTP 402).

**Scoring audit (paper-side, 2026-07-06):** two pipeline defects found and
fixed — glued snapshot headers (files without trailing newlines swallowed
the next `--- path` header; hid 2 real repairs, fabricated 4 control
"changes") and dead runs passing the infra filter (HTTP 429 lane deaths
before the first agent turn scored as behavioral nulls; 105 dead runs
across the waves, worst: Claude affordance 11/20 per arm). Fixes live in
`experiments/active/structural-elicitation/_lib/score_health.py` + all five
scorers/runners (commit `60f73887`); corrected tables are in each
experiment README under "Result (2026-07-06 scoring audit)".

**Corrected Tier 1 survey (dead-run-excluded; feature-matched DV):**

| motif | GLM-5.2 | gpt-5.4-codex | Claude Sonnet 5 | control (all fams) |
|---|---|---|---|---|
| gap (`field/5`) | 10/17 (59%) | 17/20 (85%) | 20/20 (100%) | 0 |
| rupture | 8/12 (67%) | 16/19 (84%) | 18/18 (100%) | 0 |
| affordance | 16/17 (94%) | 19/20 (95%) | 8/9 (89%) | 0 |
| affordance_named | 16/17 (94%) | 20/20 (100%) | 9/9 (100%) | 0 |
| topology | 10/16 (62%) | 17/19 (89%) | 20/20 (100%) | 0 |
| semantic (a+b) | 19/34 (56%) | 35/39 (90%) | 32/33 (97%) | 0 |

All 15 family×motif cells PROVE by each scorer's threshold (11/15 before
the audit). Aggregate: structure arms 219/259 (85%) vs feature-matched
control fires **0/236**. Bonus findings: (1) **HIGHER-WINS** unanimous 3/3
families — the pilot's "authority gradient" is dead, the semantic-collision
resolution mechanism is magnitude, not surface identity; (2) naming the
executable changes nothing (affordance_named ≈ affordance everywhere) —
the pull is the artifact, not its documentation; (3) control-arm
"activity without the act" replicates at scale, with family-specific
signatures (Claude: additive metadata enrichment, series extension —
always valid, never the feature-matched act). GLM's lower rates are
rate-limit truncation (its censoring-free runs fire at the other
families' levels), so all structure-arm rates are **lower bounds**.

**Still open (in priority order):**
1. **Tier 2 entirely** — selection steering (`elic-3X-selection-steering`)
   + feature competition. This is the ICLR control-surface upgrade; the
   go/no-go checkpoint (2026-08-01) needs it.
2. Tier 1 mini-swe ports (rupture + affordance, n=10, generic prime).
3. Family breadth: rerun kimi/gemini probes when lanes are healthy;
   deepseek if topped up. Rate-limit hygiene for reruns: stagger `--jobs`,
   or the dead-run rate will again eat 20–50% of a wave (Claude wave lost
   11/20 per affordance arm; GLM wave lost 105-run-wide chunks).

## Status round 2 (2026-07-06 later — 5-family floor + Tier 2a wave 1 + ports; verified paper-side)

Open items 2 and 3 above are now CLOSED; item 1 is half-closed:

- **5-family Tier 1 floor reached** (commit `014f97b6`, independently
  re-scored): kimi-k2.6 and gemini-3.5-flash ran all five motifs at
  n=20/cell on the fixed scorers, both waves infra-clean (0 dead, ~0
  censored). Both PROVE gap/rupture/affordance/topology; kimi PROVES
  semantic outright (40/40 with to-conf 1.00/0.00); **gemini's semantic
  binary is the survey's one weak cell** (4/20 + 9/20, INCONCLUSIVE) but
  its deconfound direction is still a clean 1.00/0.00 — HIGHER-WINS is now
  unanimous 5/5 families. Corrected 5-family aggregate: **structure arms
  558/636 (88%) vs 0/457 feature-matched control fires.** Per-family:
  GLM 70% (rate-limit-truncated wave; lower bound), gpt-5.4 91%, Claude 98%,
  kimi 99%, gemini 81% (semantic drags).
- **mini-swe ports DONE** (elic-35): rupture + affordance, n=10/cell,
  GLM-5.2 + kimi-k2.6, generic prime — 100% structure arms (one 90%
  named-arm cell) vs 0% controls. Harness generality now covers 3 of 5
  motifs on two families.
- **Tier 2a selection steering, wave 1 DONE** (elic-36, n=10/cell × 4
  conditions × 2 families): **GLM-5.2 PASS at 0.87** pooled steering
  accuracy (pre-registered ≥0.8 bar); **Claude Sonnet 5 FAIL at 0.70** —
  reported as-is. Construct validity decisive in both families: **zero
  miscues, zero flat-control false positives, all 80 runs.** Miss autopsy
  (elic-36 README): zero dead runs; Claude's misses had MORE turns than its
  hits (censoring ruled out); 4/9 explicitly noticed "entry 5 is missing"
  and chose to document it in a survey artifact instead of acting, 5/9
  never mentioned it. Same family completes the identical single-station
  gap 20/20 (elic-01) — multi-station workspaces dilute Claude's pull into
  an observer/restraint posture; GLM's does not dilute. GLM's 4 misses are
  clean, 3 at position 3 (mild position effect; ablation flagged).

**Remaining before the 2026-08-01 go/no-go (priority order):**
1. ~~**Tier 2b feature competition**~~ — **DONE 2026-07-06 (wave 1,
   `elic-37`), with one paper-side correction**: the committed Result had
   been scored from an incomplete Claude wave; full-wave re-score retracts
   Claude's `gap ≻ affordance` (0.50/0.50 at full n) — see the README's
   CORRECTION section. Corrected picture: GLM affordance-hub (3 relations,
   all P=1.00); Claude one relation (semantic ≻ affordance, P=0.75); the
   semantic–affordance pair **reverses between families** past the
   pre-registered bar on both sides — the clean family-dependence exhibit.
   Known limit: fast–fast pairs (topology–gap, gap–rupture) tie at ~1 Hz
   polling; a resolver needs per-tool-call snapshots. Absorbed into
   manuscript §5.4.
2. Optional Tier 2a widening: steering on gpt-5.4 and/or kimi (cheap,
   n=40/family) — turns "≥80% for at least one family" into a
   family-general statement or a characterized family split; also the
   position-3 ablation. A third family would also arbitrate the reversed
   semantic–affordance competition edge.
3. DeepSeek lane if topped up (would make the floor 6 families).
4. (If a fast–fast dominance answer is wanted) sub-second event capture —
   port the competition design onto the mini-swe per-tool-call
   snapshotting harness.

**Spec scorecard: every committed Tier 1 + Tier 2 item has now run at
least one recorded wave.** The Aug 01 go/no-go inputs are in hand ~3.5
weeks early.
