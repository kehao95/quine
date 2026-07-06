---
surface_kind: experiment
phase: p13-structural-elicitation
experiment_id: active/structural-elicitation/elic-03-rupture
experiment_type: single-run
status: scale-recorded
id: elic-03
legacy_id: p13:13.03
family: structural-elicitation
theory-objects: [structural-pull, trigger-structure]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient]
lineage_phase: p13-structural-elicitation
---

# 13.03-rupture: Rupture / Instability (structural-elicitation breadth probe)

## Result (pilot n=5/arm, gpt-5.4-codex-medium, 2026-06-15) — **supports LBA breadth**

| condition | n | functional event | |
|-----------|---|------------------|--|
| `rupture` (structure) | 5 | **5/5 (100%)** | structure pulls the action |
| `intact` (control)  | 5 | **0/5 (0%)**  | no structure -> no action |

`rupture` >> `intact`, with no task-specific directive and no task-time reward,
supports this as a breadth motif. Because a broken file can read like implicit
"fix me", keep it descriptive rather than independence evidence. Current design
map: `notes/breadth-designs.md` (internal planning doc, not part of this snapshot).

## Result (scale, GLM-5.2, 2026-07-06) — ICLR Tier 1 survey wave

Part of a later ICLR-upgrade survey wave (internal planning doc, not part of this snapshot;
Tier 1). `RUNS_DIR=runs-glm-scale`, `profiles/glm-5.2-zai-medium.env`,
240s wallclock, `--jobs 5`. 5 runs were interrupted mid-flight by an
unrelated mid-campaign pause and excluded from scoring (left in place for
provenance rather than deleted).

```
=== 13.03 rupture — functional DV (artifact becomes valid JSON) ===

condition     n  repaired         final_valid  changed
rupture      17  8/17 (47%)       8/17         8/17
intact       18  0/18 (0%)        18/18        0/18

--- decision: rupture vs intact ---
  rupture repaired = 0.47   intact spurious-change = 0.00
  => INCONCLUSIVE.
```

**Effect direction holds at scale (47% vs 0%), scorer verdict is
`INCONCLUSIVE` only because it falls just under this script's own >=50%
threshold** — the separation itself (47 points, zero false positives in
`intact`) is large and clean; read as supportive-but-below-the-script's-bar
rather than a null result. (Reproduce: `python3 analysis/score.py runs-glm-scale/`.)

## Result (scale, gpt-5.4-codex-medium + Claude Sonnet 5, 2026-07-06) — ICLR Tier 1 survey wave

Same design, n=20/cell, `RUNS_DIR=runs-gpt54-scale` / `runs-claude-scale`,
240s wallclock, no interrupted runs.

| family | rupture repaired | intact spurious-change | verdict |
|---|---|---|---|
| gpt-5.4-codex-medium | 16/20 (80%) | 1/20 (5%) | **PROVES** |
| Claude Sonnet 5 | 16/20 (80%) | 9/20 (45%) | INCONCLUSIVE |

**gpt-5.4-codex-medium PROVES cleanly** (80% vs 5%), the strongest replicate
of this motif yet. **Claude Sonnet 5 repairs at the same rate (80%) but is
much noisier on the control** — `intact`'s `final_valid` drops to 16/20 (vs
20/20 for GLM/gpt-5.4), meaning Claude actively rewrites the already-valid
control artifact often enough to occasionally break it. That's a real
family-specific behavioral difference (Claude is more editorially active on
artifacts generally, not specific to rupture) rather than a failure of the
motif — worth a record-grounded look before citing. (Reproduce:
`python3 analysis/score.py runs-gpt54-scale/` and `runs-claude-scale/`.)

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

| family | rupture repaired | intact spurious-change | intact final_valid | verdict |
|---|---|---|---|---|
| GLM-5.2 | 8/12 (67%) | 0/15 (0%) | 15/15 | **PROVES** (was INCONCLUSIVE) |
| gpt-5.4-codex-medium | 16/19 (84%) | 1/20 (5%) | 20/20 | **PROVES** |
| Claude Sonnet 5 | 18/18 (100%) | 5/18 (28%) | 18/18 | **PROVES** (was INCONCLUSIVE) |
| kimi-k2.6 | 20/20 (100%) | 8/20 (40%) | 20/20 | **PROVES** |
| gemini-3.5-flash | 20/20 (100%) | 0/20 (0%) | 20/20 | **PROVES** |

The two INCONCLUSIVE verdicts above are both retracted as artifacts:
**Claude never broke the valid control artifact** — `final_valid` is 18/18
(the earlier 16/20 counted four runs where a glued header made `data.json`
unparseable, not invalid); its repair rate was also undercounted 16→18 by
the same glue. Claude's remaining 5/18 control-arm changes are
**additive metadata enrichment** (e.g. appending `"sum": 3, "product": 2`
or a `meta` block — always still valid JSON), a family-specific editorial
style, not spurious breakage. GLM's 47% was dead-run dilution; among runs
whose lane survived to the first turn it is 67%, and its censoring-free
runs repair at the other families' levels.

**5-family floor reached** (2026-07-06, kimi/gemini added): repair rate is
**100% for kimi, gpt-5.4, and gemini's repair-only figure**, `final_valid`
holds at 20/20 in both new families — no control-arm breakage. Kimi is the
most editorially active on the control (40% additive changes, same
enrichment style as Claude, never invalidating), gemini the quietest (0%,
matching GLM/gpt-5.4). **PROVES on all 5 families.**

## Design

One factor: `rupture` (a syntactically broken JSON artifact (missing brace) — a contradiction, not a missing element) vs `intact` (control). Shared discipline
(identical to 13.01): missionless (no task-specific directive), frozen weights +
no task-time reward, autonomy gate `QUINE_NO_MISSION_AUTONOMY=1` held identical
across arms, single frozen policy (`PEER_COUNT=1`), neutral container paths,
schematic structure. DV is an **un-fakeable functional event**, never narration.

## DV

Does the structurally-invalid artifact become valid JSON (`json.loads` succeeds), given a broken seed? Scorer: [`analysis/score.py`](analysis/score.py).

## How to run

```bash
./run.sh .env.gpt-5.4-codex-medium 5 240 --jobs 5
python3 analysis/score.py runs/
```

Follow-up (depth): n>=10 + >=3 models; the optional directive-additivity arm.

## Parent

- Breadth program: `notes/breadth-designs.md` (internal planning doc, not part of this snapshot)
- Observation surface and submission tracking: `development/status/README.md` (internal status surface, not part of this snapshot).

## Paper Feeds

- `none-yet` - none - not-for-paper-yet - rupture / repair breadth motif (5/5 vs 0/5), directive-adjacent; feeds the Structural-Elicitation ALIFE LBA (registration pending).
