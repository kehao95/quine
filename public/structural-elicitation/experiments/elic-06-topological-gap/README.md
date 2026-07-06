---
surface_kind: experiment
phase: p13-structural-elicitation
experiment_id: active/structural-elicitation/elic-06-topological-gap
experiment_type: single-run
status: scale-recorded
id: elic-06
legacy_id: p13:13.06
family: structural-elicitation
theory-objects: [structural-pull, trigger-structure]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient]
lineage_phase: p13-structural-elicitation
---

# 13.06-topological-gap: Topological Gap (relational) (structural-elicitation breadth probe)

## Result (pilot n=5/arm, gpt-5.4-codex-medium, 2026-06-15) — **supports LBA breadth**

| condition | n | functional event | |
|-----------|---|------------------|--|
| `dangling` (structure) | 5 | **5/5 (100%)** | structure pulls the action |
| `closed` (control)  | 5 | **0/5 (0%)**  | no structure -> no action |

`dangling` >> `closed`, with no task-specific directive and no task-time reward,
supports this as the strongest non-numeric absence/realization LBA motif.
Current design map: `notes/breadth-designs.md` (internal planning doc, not part of this snapshot).

## Result (scale, GLM-5.2, 2026-07-06) — ICLR Tier 1 survey wave

Part of a later ICLR-upgrade survey wave (internal planning doc, not part of this snapshot;
Tier 1). `RUNS_DIR=runs-glm-scale`, `profiles/glm-5.2-zai-medium.env`,
240s wallclock, `--jobs 5`. 5 runs were interrupted mid-flight by an
unrelated mid-campaign pause and excluded from scoring (left in place for
provenance rather than deleted).

```
=== 13.06 topological gap — functional DV (dangling reference wired) ===

condition     n  wired
dangling     18  10/18 (56%)
closed       17  0/17 (0%)

--- decision: dangling vs closed ---
  dangling wired = 0.56   closed wired = 0.00
  => PROVES — a relational dangle pulls wiring with no instruction/reward.
```

**PROVES at scale** (56% vs 0%, GLM-5.2) — the pilot's 5/5 vs 0/5 separation
moderates at scale but the topology-specific effect (a dangling reference,
not a sequence gap) replicates cleanly on a second model family under the
Quine harness. (Reproduce: `python3 analysis/score.py runs-glm-scale/`.)

## Result (scale, gpt-5.4-codex-medium + Claude Sonnet 5, 2026-07-06) — ICLR Tier 1 survey wave

Same design, n=20/cell, `RUNS_DIR=runs-gpt54-scale` / `runs-claude-scale`,
240s wallclock, no interrupted runs.

| family | dangling wired | closed wired | verdict |
|---|---|---|---|
| gpt-5.4-codex-medium | 17/20 (85%) | 0/20 (0%) | **PROVES** |
| Claude Sonnet 5 | 20/20 (100%) | 0/20 (0%) | **PROVES** |

**PROVES on both**, and the strongest motif in the survey by consistency —
all three families now tested (GLM 56%, gpt-5.4 85%, Claude 100%) show zero
control-arm false positives and a clean, large separation. This is the
best-replicating non-numeric (relational/topological) motif in the batch.
(Reproduce: `python3 analysis/score.py runs-gpt54-scale/` and
`runs-claude-scale/`.)

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

| family | dangling wired | closed wired | verdict |
|---|---|---|---|
| GLM-5.2 | 10/16 (62%) | 0/14 (0%) | **PROVES** |
| gpt-5.4-codex-medium | 17/19 (89%) | 0/19 (0%) | **PROVES** |
| Claude Sonnet 5 | 20/20 (100%) | 0/20 (0%) | **PROVES** |
| kimi-k2.6 | 19/20 (95%) | 0/20 (0%) | **PROVES** |
| gemini-3.5-flash | 20/20 (100%) | 0/20 (0%) | **PROVES** |

Numbers barely move (denominators shrink by the dead runs); the
all-families PROVES read above stands. Still the most consistent motif:
zero control-arm false positives anywhere.

**5-family floor reached** (2026-07-06, kimi/gemini added): **PROVES on all
5 families**, 62–100% structure-arm rate, 0% control everywhere — the only
motif in the survey with zero control-arm false positives across every
single family tested.

## Design

One factor: `dangling` (a manifest referencing a target that does not exist — a relational/topological dangle, not a sequence gap) vs `closed` (control). Shared discipline
(identical to 13.01): missionless (no task-specific directive), frozen weights +
no task-time reward, autonomy gate `QUINE_NO_MISSION_AUTONOMY=1` held identical
across arms, single frozen policy (`PEER_COUNT=1`), neutral container paths,
schematic structure. DV is an **un-fakeable functional event**, never narration.

## DV

Is the dangling reference wired (`parts/gamma`, referenced by the manifest but absent, created)? Scorer: [`analysis/score.py`](analysis/score.py).

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

- `none-yet` - none - not-for-paper-yet - topological / relational gap breadth motif (5/5 vs 0/5); feeds the Structural-Elicitation ALIFE LBA (registration pending).
