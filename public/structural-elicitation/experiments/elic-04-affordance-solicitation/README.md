---
surface_kind: experiment
phase: p13-structural-elicitation
experiment_id: active/structural-elicitation/elic-04-affordance-solicitation
experiment_type: single-run
status: scale-recorded
id: elic-04
legacy_id: p13:13.04
family: structural-elicitation
theory-objects: [structural-pull, affordance-null-ecology]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient]
lineage_phase: p13-structural-elicitation
---

# 13.04-affordance-solicitation: Affordance Solicitation (structural-elicitation breadth probe)

## Result (pilot n=5/arm, gpt-5.4-codex-medium, 2026-06-15) — **supports LBA breadth**

| condition | n | functional event | |
|-----------|---|------------------|--|
| `affordance` (structure) | 5 | **5/5 (100%)** | structure pulls the action |
| `inert` (control)  | 5 | **0/5 (0%)**  | no structure -> no action |

`affordance` >> `inert`, with no task-specific directive and no task-time
reward, supports this as one of the strongest non-gap LBA motifs. Current design
map: `notes/breadth-designs.md` (internal planning doc, not part of this snapshot).

## Result (scale, GLM-5.2, 2026-07-06) — ICLR Tier 1 survey wave

Part of a later ICLR-upgrade survey wave (internal planning doc, not part of this snapshot;
Tier 1). `RUNS_DIR=runs-glm-scale`, `profiles/glm-5.2-zai-medium.env`,
240s wallclock, `--jobs 5`, n=20 per condition (`affordance`,
`affordance_named`, `inert`). 5 runs were interrupted mid-flight by an
unrelated mid-campaign pause and excluded from scoring (left in place for
provenance rather than deleted).

```
=== 13.04 affordance solicitation — functional DV (executable invoked) ===

condition     n  invoked
affordance   18  16/18 (89%)
inert        19  0/19 (0%)

--- decision: affordance vs inert ---
  affordance invoked = 0.89   inert invoked = 0.00
  => PROVES — a fitting affordance solicits invocation with no instruction/reward.
```

**PROVES at scale** (89% vs 0%, GLM-5.2) — one of the cleanest replications in
the survey. Note: `analysis/score.py`'s table only reports `affordance` vs
`inert`; the `affordance_named` arm ran at full n=20 alongside them (data on
disk under `runs-glm-scale/`) but isn't summarized by this scorer — not yet
folded into a number, flagged here rather than silently dropped. (Reproduce:
`python3 analysis/score.py runs-glm-scale/`.)

## Result (scale, gpt-5.4-codex-medium + Claude Sonnet 5, 2026-07-06) — ICLR Tier 1 survey wave

Same design, n=20/cell, `RUNS_DIR=runs-gpt54-scale` / `runs-claude-scale`,
240s wallclock, no interrupted runs. `affordance_named` ran at full n=20 in
both waves too, same scorer caveat as above.

| family | affordance invoked | inert invoked | verdict |
|---|---|---|---|
| gpt-5.4-codex-medium | 19/20 (95%) | 0/20 (0%) | **PROVES** |
| Claude Sonnet 5 | 8/20 (40%) | 0/20 (0%) | INCONCLUSIVE |

**gpt-5.4-codex-medium PROVES cleanly** (95% vs 0%), matching GLM's strength.
**Claude Sonnet 5 invokes far less often** (40% vs 89–95% for the other two
families) — still a clear positive gap over its own 0% control, but well
below this scorer's decisiveness bar. This is the weakest Claude replicate in
the survey; worth a record-grounded read on *why* Claude solicits the
affordance less often (more caution around executing unfamiliar binaries?)
before treating it as a null. (Reproduce: `python3 analysis/score.py
runs-gpt54-scale/` and `runs-claude-scale/`.)

## Result (2026-07-06 scoring audit — CORRECTED scale numbers; supersedes the scale tables above)

A post-campaign audit (structure-originates-direction dossier, Tier 1
absorption) found two scoring-pipeline defects that contaminated the scale
tables above; both are now fixed in `analysis/score.py` (shared helpers:
[`../_lib/score_health.py`](../_lib/score_health.py)) and the corrected
numbers below are what `python3 analysis/score.py runs-*-scale/` reproduces.

1. **Glued snapshot headers.** The snapshot loop glued the next `--- path`
   header onto any file lacking a trailing newline, so the line-based parser
   silently dropped files — hiding real structure-arm acts and fabricating
   control-arm "changes". Runners now emit a newline guard; the parser
   re-splits against the authoritative `final-files.txt` list.
2. **Dead runs passing the infra filter.** Runs whose LLM lane died on HTTP
   429 before the FIRST agent turn still contain the two prime echoes in
   `meta/agent-text.txt`, so the old `'"content"' not in text` check scored
   them as live behavioral nulls. The scorer now strips the prime echoes and
   excludes zero-output runs from every denominator.

Runs that produced output and then died on a mid-session LLM failure are
kept in the denominators as censored observations (conservative: truncation
can only deflate structure-arm rates; controls are unaffected at ~0). The
GLM wave was hit hardest by rate limiting — its lower rates are truncation,
not a weaker pull (its censoring-free runs fire at the other families'
levels).

| family | affordance invoked | affordance_named invoked | inert invoked | verdict |
|---|---|---|---|---|
| GLM-5.2 | 16/17 (94%) | 16/17 (94%) | 0/17 (0%) | **PROVES** |
| gpt-5.4-codex-medium | 19/20 (95%) | 20/20 (100%) | 0/20 (0%) | **PROVES** |
| Claude Sonnet 5 | 8/9 (89%) | 9/9 (100%) | 0/9 (0%) | **PROVES** (was INCONCLUSIVE) |
| kimi-k2.6 | 18/18 (100%) | 19/19 (100%) | 0/20 (0%) | **PROVES** |
| gemini-3.5-flash | 20/20 (100%) | 20/20 (100%) | 0/20 (0%) | **PROVES** |

The "Claude solicits the affordance less often (40%)" reading above is
retracted as an artifact: **11 of Claude's 20 runs per arm were dead**
(Anthropic 429 rate-limit retries exhausted before the first agent turn —
`meta/agent-text.txt` contains only the two prime echoes; the old infra
filter passed them because the echoes contain `"content"`). Among runs
where the agent actually spoke, Claude invokes at 17/18 (94%) across the
two structure arms — indistinguishable from GLM and gpt-5.4. The
`affordance_named` arm (previously unscored) is now a first-class column:
**naming/documenting the executable changes nothing** (94–100% in both
arms, all families) — the pull is carried by the affordance itself, not by
its README description.

**5-family floor reached** (2026-07-06, kimi/gemini added): kimi and gemini
both invoke at 100% on both structure arms (kimi excludes 3 dead runs from
its own rate-limit hiccups, same censoring convention as Claude above),
0% inert everywhere. **PROVES on all 5 families** — this is the strongest,
most consistent motif in the survey alongside topological-gap.

## Design

One factor: `affordance` (a well-formed, un-invoked, fitting executable present among neutral files — no gap, no error) vs `inert` (control). Shared discipline
(identical to 13.01): missionless (no task-specific directive), frozen weights +
no task-time reward, autonomy gate `QUINE_NO_MISSION_AUTONOMY=1` held identical
across arms, single frozen policy (`PEER_COUNT=1`), neutral container paths,
schematic structure. DV is an **un-fakeable functional event**, never narration.

## DV

Is the un-invoked executable `./step` actually run (its output `out/ran.txt` appears)? Scorer: [`analysis/score.py`](analysis/score.py).

## How to run

```bash
./run.sh .env.gpt-5.4-codex-medium 5 240 --jobs 5
python3 analysis/score.py runs/
```

Follow-up (depth): n>=10 + >=3 models; the optional directive-additivity arm.

## Parent

- Breadth program: `notes/breadth-designs.md` (internal planning doc, not part of this snapshot)
- Observation surface and submission tracking: `development/status/README.md` (internal status surface, not part of this snapshot).

## Surface Map

```text
13.04-affordance-solicitation/
├── README.md
├── run.sh / run-container.sh   # local + containerized runners
├── assets/affordance.go        # the affordance surface presented to the agent
├── analysis/score.py           # scorer (stdlib only)
└── runs/ , runs-deepseek/       # retained run trees (DVC-tracked via .dvc pointers)
```

## Paper Feeds

- `none-yet` - none - not-for-paper-yet - affordance-uptake breadth motif (5/5 vs 0/5), strongest non-completion-prior read; feeds the Structural-Elicitation ALIFE LBA (registration pending).
