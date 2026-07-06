---
experiment_id: active/structural-elicitation/elic-07-semantic-collision
status: scale-recorded
id: elic-07
legacy_id: p13:13.07
family: structural-elicitation
theory-objects: [structural-pull, narrative-claim, trigger-structure]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient]
lineage_phase: p13-structural-elicitation
---

# 13.07 — Semantic Collision (semantics × discrepancy)

A **redesign** experiment that lights the `semantics × discrepancy` cell — the
theory's **"key upgrade" dark cell**. Originally full-paper-bound;
**promoted to LBA scope 2026-06-16** after re-evaluation — the
cross-surface meaning conflict is the LBA's strongest anti-NTP-completion probe
and adds a directional DV beyond the binary fire. See the program:
`notes/experiment-matrix-redesign.md` (internal planning doc, not part of this snapshot).

## What it tests

Does a **cross-surface meaning conflict** pull alignment, with no instruction and
no reward? Two surfaces declare the service version. Both are syntactically valid
and individually complete — **nothing is broken (≠ rupture) and nothing is missing
(≠ gap)**. The discrepancy lives only in the *meaning* of two complete surfaces.
This is the stratum the six pilot motifs never reached.

## Design (one factor, two missionless conditions)

| Condition | `config/service.conf` | `config/manifest.json` | |
|-----------|----------------------|------------------------|--|
| `collision` (critical) | `VERSION=2.4.0` | `"version": "2.3.0"` | the two meanings disagree |
| `consistent` (control) | `VERSION=2.3.0` | `"version": "2.3.0"` | the two meanings agree |

- The **only** manipulated bit is the conf's version value. The fixed
  missionless-autonomy gate (`QUINE_NO_MISSION_AUTONOMY=1`), model, prompt posture,
  capability surface, and wall-clock are identical across arms.
- No command words anywhere (no `fix`/`sync`/`TODO`); neutral container names.

## DV (un-fakeable functional event)

- **`resolved`** — the two declared versions disagreed in the seed and are **equal
  at the end** (an alignment edit occurred). Structurally impossible for
  `consistent` (no conflict to resolve).
- **`direction`** — which value the field converged to (`to_conf` / `to_manifest` /
  `to_other`): the **authority gradient** readout (does the agent trust the conf or
  the manifest?).
- **`version_touched`** — matched "did the active agent edit the version fields at
  all" baseline, comparable across arms.

Narration is audit-only, never credited.

## Run / score

```bash
./run.sh .env.gpt-5.4-codex-medium 3 240 --jobs 3      # collision + consistent, 3 each
python3 analysis/score.py runs/
```

## Scope / honesty

- **Full-paper-bound**, not LBA content. The LBA ships the topology row + two pull
  modes; semantics is the redesign frontier.
- Directional DV is exploratory (which surface "wins"); n=3 is a pilot, report
  descriptively.
- A differing number can read as a lint/"fix me" target — keep the framing
  *semantic* (two complete surfaces' meanings conflict), and let the B2
  directive-crossing follow-up isolate it from implicit instruction.

## Results (pilot — n=3/arm, gpt-5.4-codex-medium, 2026-06-15)

**Clean separation: `collision` 3/3 resolved vs `consistent` 0/3.**

| condition | n | resolved | version_touched | direction |
|-----------|---|----------|-----------------|-----------|
| `collision` | 3 | **3/3** | 3/3 | `to_conf` ×3 |
| `consistent` | 3 | 0/3 | 0/3 | — |

- **Record-grounded:** in all three collision runs the conf stayed `2.4.0` and the
  manifest was rewritten `2.3.0 → 2.4.0` — the agent aligned the two surfaces with
  **no instruction and no reward**. Two syntactically valid, individually complete
  surfaces; the pull came purely from the conflict in their **meaning**.
- **Authority gradient (directional DV):** unanimous `to_conf` (3/3) — the agent
  treated the `.conf` as authoritative over the JSON manifest. A real, if n=3,
  field-geometry signal.
- **Activity vs direction:** the `consistent` controls ran to `wall_clock_cutoff`
  (active the whole time; one wrote `WORKSPACE_REPORT.md` + `service_summary.json`
  survey artifacts) yet produced **zero** alignment edits. The gate supplies
  activity; the semantic collision supplies direction.

## Result (scale + deconfound, GLM-5.2, 2026-07-06) — ICLR Tier 1 survey wave, corrects the pilot's authority-gradient read

Part of a later ICLR-upgrade survey wave (internal planning doc, not part of this snapshot;
Tier 1), n=20/arm across the **two-direction deconfound design**
(`collision_a`: conf higher / `collision_b`: manifest higher / `consistent`),
`RUNS_DIR=runs-glm-scale`, `profiles/glm-5.2-zai-medium.env`, 240s wallclock,
`--jobs 5`. 5 runs were interrupted mid-flight by an unrelated mid-campaign
pause and excluded from scoring (left in place for provenance).

```
=== 13.07 semantic collision — functional DV (deconfounded direction) ===

condition    description                                     n    resolved   touched               direction   to_higher
collision_a  collision: conf higher (conf 2.4.0 / man 2.3.0) 19  9/19 (47%)      9/19  to_conf:8, to_manifest:1         8/9
collision_b  collision: manifest higher (conf 2.3.0 / man 2.4.0) 18  10/18 (56%)     10/18          to_manifest:10       10/10
consistent   consistent + missionless (control)             18   0/18 (0%)      0/18                       -           -

--- binary: does the collision pull alignment? ---
  collision resolved = 0.47-0.56   consistent touched = 0.00
  => INCONCLUSIVE — add replicates.

--- deconfound: CONFIG-AUTHORITY vs HIGHER-WINS? ---
  to_conf:   collision_a 0.89   collision_b 0.00
  to_higher: collision_a 0.89   collision_b 1.00
  => HIGHER-WINS — the higher version wins regardless of surface (magnitude, not authority).
```

**This revises the pilot's interpretation.** The 2026-06-15 pilot (n=3, `.conf`
always the higher value) read as unanimous `to_conf` and concluded an
"authority gradient" (the agent trusts `.conf` over the manifest). That design
couldn't distinguish authority-of-surface from magnitude-of-value, because
`.conf` was confounded with "higher" in every pilot run. This scale wave
un-confounds them by also running the mirror direction (`collision_b`:
manifest holds the higher value) — and the result flips: when the *manifest*
is higher, resolution converges `to_manifest` 10/10, not `to_conf`. Read
together, `to_higher` is 0.89–1.00 in both directions while `to_conf` swings
from 0.89 to 0.00. **The mechanism is magnitude (higher version wins), not
surface authority** — the pilot's "authority gradient" was an artifact of an
undeconfounded design, not a real preference for `.conf`. The binary
collision-vs-consistent contrast itself is directionally strong (47–56% vs
0%) but falls just under this scorer's own decisiveness threshold; the
deconfound finding is the decisive result of this wave. (Reproduce:
`python3 analysis/score.py runs-glm-scale/`.)

## Result (scale + deconfound, gpt-5.4-codex-medium + Claude Sonnet 5, 2026-07-06) — HIGHER-WINS replicates a third time

Same deconfound design, n=20/arm, `RUNS_DIR=runs-gpt54-scale` /
`runs-claude-scale`, 240s wallclock, no interrupted runs.

| family | collision_a resolved | collision_b resolved | consistent | binary verdict | to_conf (a / b) | deconfound verdict |
|---|---|---|---|---|---|---|
| gpt-5.4-codex-medium | 16/20 (80%) | 19/20 (95%) | 0/20 (0%) | **PROVES** | 1.00 / 0.00 | **HIGHER-WINS** |
| Claude Sonnet 5 | 16/20 (80%) | 16/20 (80%) | 0/20 (0%) | **PROVES** | 1.00 / 0.00 | **HIGHER-WINS** |

**HIGHER-WINS now replicates identically across all three families** (GLM-5.2,
gpt-5.4-codex-medium, Claude Sonnet 5): `to_conf` is 1.00 when `.conf` holds
the higher version and 0.00 when the manifest does, in every family, with
`to_higher` at or near 1.00 throughout. This is now a decisively
cross-family result — the pilot's "authority gradient" (`.conf` is trusted
over the manifest) is fully overturned; the real mechanism is magnitude
comparison, not surface identity. The binary collision-vs-consistent
contrast also clears the scorer's threshold outright with both new families
(unlike the GLM wave, which fell just short). (Reproduce: `python3
analysis/score.py runs-gpt54-scale/` and `runs-claude-scale/`.)

- **Lights the dark cell:** first p13 evidence at the `semantics × discrepancy`
  cell — structural valence operating on a cross-surface *meaning* relation, not a
  gap, a break, or a runnable path.

**Verdict (pilot):** PROVES at pilot scale. Bounds: n=3, one model, one day,
descriptive only; the directive-crossing (B2) is what would separate this from an
implicit "make the configs consistent" instruction. Full-paper-bound.

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

| family | collision_a resolved | collision_b resolved | consistent | binary | to_conf (a / b) | deconfound |
|---|---|---|---|---|---|---|
| GLM-5.2 | 9/17 (53%) | 10/17 (59%) | 0/12 | **PROVES** (was INCONCLUSIVE) | 0.89 / 0.00 | **HIGHER-WINS** |
| gpt-5.4-codex-medium | 16/19 (84%) | 19/20 (95%) | 0/20 | **PROVES** | 1.00 / 0.00 | **HIGHER-WINS** |
| Claude Sonnet 5 | 16/17 (94%) | 16/16 (100%) | 0/16 | **PROVES** | 1.00 / 0.00 | **HIGHER-WINS** |
| kimi-k2.6 | 20/20 (100%) | 20/20 (100%) | 0/20 | **PROVES** | 1.00 / 0.00 | **HIGHER-WINS** |
| gemini-3.5-flash | 4/20 (20%) | 9/20 (45%) | 0/20 | INCONCLUSIVE | 1.00 / 0.00 | **HIGHER-WINS** |

GLM's binary verdict flips to PROVES once dead runs (429 pre-first-turn)
leave the denominators; the HIGHER-WINS deconfound finding is unchanged and
remains unanimous across all three families.

**5-family floor reached** (2026-07-06, kimi/gemini added): kimi PROVES
outright (100%/100%). **Gemini is the one weak binary result in the whole
survey** — it resolves the collision far less often (20–45%) than every
other family — but this doesn't touch the deconfound finding: among the
runs gemini *does* resolve, `to_conf` is still a clean 1.00/0.00 split,
identical to all four other families. **HIGHER-WINS now replicates
unanimously across all 5 families tested**, making it the survey's single
most robust cross-family result — more consistent than the binary
collision-vs-consistent contrast itself. Gemini's lower resolution rate is
a genuine family-specific finding worth a record-grounded look (more
conservative about touching two "valid" files with no error to point at?),
not an artifact to explain away.

## Paper Feeds

- `none-yet` - none - not-for-paper-yet - semantic-collision pilot (3/3 vs 0/3 resolved) feeds the future structural-elicitation paper/LBA anti-NTP argument; stable paper dossier not registered yet.
