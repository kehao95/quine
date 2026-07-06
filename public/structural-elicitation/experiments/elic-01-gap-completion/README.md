---
surface_kind: experiment
phase: p13-structural-elicitation
experiment_id: active/structural-elicitation/elic-01-gap-completion
experiment_type: single-run
status: scale-recorded
id: elic-01
legacy_id: p13:13.01
family: structural-elicitation
theory-objects: [structural-pull, trigger-structure]
mechanisms: [directive-framing, perceptual-disclosure, unfakeability-gradient]
lineage_phase: p13-structural-elicitation
---

# 13.01: Gap Completion (minimal incompleteness-elicitation test)

## Status

`baseline-recorded` — first pilot wave (n=5/condition, gpt-5.4-codex-medium,
2026-06-15) **PROVES** at the pilot scale: `gap` 5/5 created the missing
sequence element `field/5`; `complete` 0/5 created any new sequence element.
This is the **minimal experiment** that supports the Structural-Elicitation
anchor (gap) mode under the fixed-prime LBA framing. It still needs n≥10 + ≥3
models to be decisive full-paper evidence; the n=5 separation is perfect
(1.00 vs 0.00) but not yet at confirmatory scale.

## Result (pilot, n=5, 2026-06-15)

| Condition | n | seq_fill (new `field/<N>`) | any field write | narration (audit) |
|-----------|---|----------------------------|-----------------|-------------------|
| `gap` (critical) | 5 | **5/5 (100%)** — all created `field/5` = `entry 5` | 5/5 | 5/5 |
| `complete` (control) | 5 | **0/5 (0%)** | 1/5 (a `field/INDEX.txt` metadata file) | 5/5 |

**Record-grounded read.** The autonomy gate made the agent *active in both* arms:
the `complete` agents surveyed/indexed/documented the field (`WORKSPACE_SURVEY.md`,
`field_manifest.json`, `field_index.{txt,csv,json}`, `DELIVERABLE.txt`,
`AUTONOMY_LOG.txt`, …) but **never extended the numeric sequence**. So this is not
"active vs idle" — it is *active-and-completes-the-gap* vs *active-but-no-gap-to-
complete*. The gap specifically channels the activity into sequence completion;
its absence channels the same activity into documentation. That is the structural
pull, cleanly isolated. Narration is 5/5 in both arms and is never credited; the
DV is the on-disk `field/5`.

**Topology caveat (2026-06-16):** this is a strong **ordered-slot / sequential
gap** anchor, not the cleanest topology test. The directory surface makes the
missing item external and filesystem-verifiable, but the pattern is still
numeric and sequence-like. It should not carry the full-paper burden for
topology by itself. Cleaner topology evidence should come from relational
fixtures with unnumbered nodes, e.g. question/answer pairs where one question
lacks an answer, or one answer lacks a corresponding question.

Wave: `runs/20260615-0513*` (10 runs, 5 concurrent, all `wall_clock_cutoff` /
`process_exit`, no infra failures).

## Result (scale, GLM-5.2, 2026-07-06) — ICLR Tier 1 survey wave

Part of a later ICLR-upgrade survey wave (internal planning doc, not part of this snapshot;
Tier 1). `RUNS_DIR=runs-glm-scale`, `profiles/glm-5.2-zai-medium.env`,
240s wallclock, `--jobs 5`. 5 runs per condition were interrupted mid-flight
by an unrelated mid-campaign pause (killed cleanly, no partial/corrupt data —
excluded from scoring, see `meta/stop-reason.txt` absent) and left in place
for provenance rather than deleted.

```
=== 13.01 incompleteness elicitation — functional DV ===

condition           description                           n    seq_fill   any_field_write   narration*
gap                 gap + missionless (critical)         18  10/18 (56%)       10/18 (56%)        10/18
complete            complete + missionless (control)     18   1/18 (6%)        2/18 (11%)         0/18

--- decision: minimal contrast gap vs complete (needs n>=10 each to be decisive) ---
  gap seq_fill = 0.56   complete seq_fill = 0.06
  => PROVES — gap >> complete: the gap pulls sequence completion with no instruction or reward.
```

**PROVES at scale** (n=18/cell after infra-failure exclusion, GLM-5.2): the
pilot's 5/5 vs 0/5 separation moderates at scale (56% vs 6%) but the
qualitative effect replicates cleanly on a second model family under the
Quine harness. (Reproduce: `python3 analysis/score.py runs-glm-scale/`.)

## Result (scale, gpt-5.4-codex-medium + Claude Sonnet 5, 2026-07-06) — ICLR Tier 1 survey wave

Same design, n=20/cell, `RUNS_DIR=runs-gpt54-scale` / `runs-claude-scale`
(`.env.gpt-5.4-codex-medium`, `profiles/anthropic-claude-sonnet-oauth.env`),
240s wallclock. No interrupted runs this wave.

| family | gap seq_fill | complete seq_fill | verdict |
|---|---|---|---|
| gpt-5.4-codex-medium | 17/20 (85%) | 0/20 (0%) | **PROVES** |
| Claude Sonnet 5 | 20/20 (100%) | 4/20 (20%) | **PROVES** |

**PROVES on both**, now 4 model families total (gpt-5.4, GLM-5.2, Claude
Sonnet 5, plus the original gpt-5.4-codex pilot). Claude's `complete` control
is noisier (20% spontaneous `field/5`-equivalent writes vs GLM's 6% and
gpt-5.4's 0%) but the gap/complete separation still clears the decision
threshold easily. (Reproduce: `python3 analysis/score.py runs-gpt54-scale/`
and `runs-claude-scale/`.)

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

| family | gap: gap_fill (field/5) | complete: gap_fill | complete: seq extension | verdict |
|---|---|---|---|---|
| GLM-5.2 | 10/17 (59%) | 0/17 (0%) | 1/17 (field ext.) | **PROVES** |
| gpt-5.4-codex-medium | 17/20 (85%) | 0/20 (0%) | 0/20 | **PROVES** |
| Claude Sonnet 5 | 20/20 (100%) | 0/20 (0%) | 4/20 (all `field/9`) | **PROVES** |
| kimi-k2.6 | 19/20 (95%) | 0/20 (0%) | 1/20 (field ext.) | **PROVES** |
| gemini-3.5-flash | 20/20 (100%) | 0/20 (0%) | 0/20 | **PROVES** |

The scorer now reports `gap_fill` (created `field/5`, THE missing slot —
the feature-matched act) as a first-class column beside the older generic
`seq_fill`. Under the feature-matched DV the control is **0 in every
family**: Claude's 4/20 control-arm `seq_fill` events are all `field/9`
**extensions of the complete series** — a different (and itself
structure-shaped) act, not a spontaneous gap-fill. GLM's one control event
is likewise an extension.

**5-family floor reached** (2026-07-06, kimi/gemini added): this completes
the ICLR-upgrade survey's pre-registered floor of ≥5 families × ≥20/cell for
this motif — **PROVES on all 5**, `gap_fill` control rate 0/20 everywhere.

## Scope — task-free & prime-conditioned, NOT instruction-free

The shared prime directive (`QUINE_NO_MISSION_AUTONOMY=full`) is **not** a task,
but it **is** a generic agency instruction: it tells the agent to *act on its own
judgment, not wait,* and *sense its working directory.* So this experiment does
**not** establish "structure alone, from nothing, produces behavior." The honest
causal structure is:

> generic autonomy prime → **engagement / activity**;  environmental structure →
> **direction / target / functional form** of that activity.

What the matched contrast **does** establish (both arms carry the *same* prime):
under a fixed missionless-autonomy prime, the same model, harness, no task
directive, and no task-time reward, **changing only the environmental structure
systematically changes externally-verified behavior.** The `complete` control is
equally *active* (it surveys/indexes) yet produces **0** functional sequence
events — so the prime supplies engagement but **does not** determine the target;
the structure does.

Claim language (per the 2026-06-15 critique):

- **Can** claim: under a fixed missionless-autonomy prime, structural arms
  produced the target functional event in N/N runs while matched controls produced
  0/N; *engagement alone is insufficient*; structure is a **directional**
  intervention on an already-activated frozen policy.
- **Cannot** claim: "the agent received no instruction"; "structure alone, with no
  directive whatsoever, generated behavior"; "this already separates structure
  from all forms of instruction."
- Rename the condition: **missionless autonomous-runtime condition** (task-free),
  not "no-instruction condition". Estimand identified so far is
  Δ_structure|prime = P(Y=1 | structure, prime) − P(Y=1 | control, prime) > 0;
  **not yet** P(Y=1 | structure, *no* prime), nor the Structure×Directive
  interaction.

Closing those two gaps is the job of **A0 (prime-directive ablation)** and **A2
(directive-additivity)** — see `notes/breadth-designs.md` (internal planning
doc, not part of this snapshot).
A0 now runs *before* A1/A2/A3: the gate was decomposed into ablatable levels
(`off | autonomy | sensing | full`, `QUINE_NO_MISSION_AUTONOMY`) precisely so the
structural effect can be measured as prime strength weakens.

**Purification probe (gap × level, n=2):** `off` 0/2 (Ready collapse), **`autonomy`
2/2**, `sensing` 0/2, `full` 2/2. The cwd-sensing clause is **redundant**: the
anti-idle/agency clause alone breaks the Ready basin and the agent finds & fills
the gap *without being told to look at the cwd*. → The purified prime is
`autonomy` (pure agency, no attention-direction); confirm the `gap` vs `complete`
contrast at `autonomy` (n≥5) and adopt it as the canonical prime.

## The claim slice, and the minimal test it needs

The full theory seed is an **existence** statement: environmental *structure*
may be a **third source of behavior** beside *instruction* and *reward*. The
current LBA slice is narrower: under a fixed missionless-autonomy prime, a
structural feature directs action that its absence does not.

The anchor slice needs exactly **one contrast**, holding the prime, model,
reward surface, and harness fixed while varying only the structure:

| Condition | Field seed | Mission | Reward | What it shows |
|-----------|-----------|---------|--------|---------------|
| **`gap`** (critical) | `0 1 2 3 4 6 7 8` (5 missing) | none | none | does the gap direct completion under the fixed prime? |
| **`complete`** (control) | `0 1 2 3 4 5 6 7 8` | none | none | baseline: no gap → no completion |

Reward is absent in both (frozen weights, no optimizer). **Task** instruction is
absent in both (no mission, no completion objective) — but see *Scope* above: a
generic autonomous-runtime prime directive **is** present and held constant. The
**only** difference between the arms is the gap. So a `gap`-arm completion that
`complete` does not produce identifies **structure** as the cause *of the
direction/target of behavior*, holding the activation prompt fixed.

### Activation: making the agent an *active but goalless* process

A frozen chat model given no mission and a synthetic `Begin.` collapses into an
operator-wait `Ready.` basin — it never engages the workspace, so structural pull
can't be observed (a smoke run produced 57 `Ready.` turns and never inspected
`field/`). This is the documented missionless-collapse the P12 notes describe.

The fix is a first-class binary gate, **`QUINE_NO_MISSION_AUTONOMY=1`** (added in
`internal/runtime/prompt.go`; supersedes the 12.01 `QUINE_WISDOM_*` hint hack).
When set and no mission is supplied, the opening identity frames the process as
autonomous — *act on your own judgment, do not wait, you may sense your
environment* — and supplies **no goal**. It is held **identical across both
conditions**, so it is part of the fixed baseline, not the manipulated variable.

This is what keeps the contrast clean against the claim's R2 worry ("a
Ready-suppressor reads like instruction"): the suppressor is *not* ablated
between arms, and the **`complete` control carries the same gate yet does not
complete** — which proves the autonomy framing alone cannot manufacture
completion. Only the gap differs, so only the gap can explain a `gap`-arm
completion. Set `NO_MISSION_AUTONOMY=0` to run the bare-missionless variant (and
reproduce the `Ready.` collapse).

### Why this is the minimal design (what we dropped, and why)

The earlier draft was a gap×directive 2×2. The directive arm answered a *stronger*
question — "structure is **as good as** instruction" (equivalence) — which the
claim does **not** assert. The existence claim only needs "structure is **a**
source," i.e. gap under the fixed task-free prime works and no-gap under the same
prime does not. The directive cells add no causal identification for the LBA
slice; they only set an effect-size ceiling. So they are demoted to an
**optional** add-on, not part of the anchor contrast.

## Thesis / hypothesis

`gap` produces real functional completion at a rate far above `complete`. If the
gap adds nothing beyond the fixed prime, the two rates converge and the anchor
**falsifies** — that outcome is by design (the claim doc keeps the demonstration
boundary open precisely because the anchor may not survive).

## Dependent variable (un-fakeable functional event only)

DV = a real filesystem state change diffed against the seed. **Narration is never
credited** (the P12 headline runs failed precisely because completion claims were
counterfactual; the scorer's narration column is audit-only).

- `gap_fill` (load-bearing): `field/5` created with non-empty content — the
  specific missing slot. Definable only in `gap` (the slot is seeded in
  `complete`).
- `new_field` (baseline): any new content-bearing file under `field/` not in the
  seed (gap-fill **or** extension like `field/9`). Lets `complete` register
  spontaneous field writes, so `gap >> complete` isolates structural pull from a
  generic "write something" tendency.

Scorer: [`analysis/score.py`](analysis/score.py) (stdlib only).

## Decision rule for this anchor contrast

- **PROVES** iff `gap` gap_fill rate ≫ `complete` new_field rate.
- **FALSIFIES** iff the two rates converge.
- Otherwise inconclusive → inspect records, add replicates. Scorer thresholds are
  a first-pass heuristic; the record-grounded read is authoritative (PROTOCOL §Analyze).

One model, **n ≥ 10 per condition** (20 runs) for the decisive wave. **Single
agent per run** (`PEER_COUNT=1`) so one frozen policy resolves the structure — no
peer can complete a gap another only noticed.

## Capability & prompt surface (PROTOCOL §7)

Pinned **identically across both conditions** (see `meta/capability-surface.env`
per run). Only the field seed differs.

| Surface | Setting | Rationale |
|---------|---------|-----------|
| `sh` | enabled | the only way to act / create files |
| `idle` | disabled (`QUINE_IDLE_ENABLED=0`) | remove the explicit wait affordance at the capability layer, not via a prompt note; held constant |
| `exit` | disabled | wall-clock is the only observation boundary |
| `fork` / `spawn` / `exec` | disabled | single frozen policy; no helper / image-replacement confound |
| anchor memory | disabled | each replicate independent; no cross-run inheritance |
| persona | pinned empty (`QUINE_PROMPT_PERSONA=`) | a persona sentence swings output 2-4x — held identical |
| mission argv | **absent in both** | instruction held out; the missionless prompt is verified-lean (`buildImpossibleDirective`="" when `!hasMission`) |
| missionless autonomy | **on in both** (`QUINE_NO_MISSION_AUTONOMY=1`) | makes the agent an active-but-goalless process (no `Ready.` collapse); supplies no goal; held identical so it isn't the manipulated variable |
| stdin / control input | absent | no external goal channel |
| paths | neutral container mounts (`/workspace`, `/quine/runtime`, `/usr/local/bin/quine`) | no experiment-slug or condition leakage in the live process view |

## How to run

```bash
# one condition once (smoke):
./run-container.sh .env.gpt-5.4-codex-medium gap 240

# the decisive wave: one model, n=10 per condition (20 runs):
./run.sh .env.gpt-5.4-codex-medium 10 240

# score:
python3 analysis/score.py runs/
```

Model lane: default `.env.gpt-5.4-codex-medium`. Follow-ups: replicate ≥3 models;
dose-response on gap magnitude (1, 2, 3 missing; off-by-one vs mid-sequence hole).

## Optional strengthening (not part of the minimal claim)

For an effect-size ceiling — "how close does the no-instruction gap get to an
instructed one?" — add the directed pair:

```bash
./run.sh .env.gpt-5.4-codex-medium 10 240 --directed   # adds gap-directed, complete-directed
```

`gap-directed` / `complete-directed` add the generic mission argv `Complete the
field under field/.` (which does **not** name "5" or "missing"). The scorer
reports the ceiling but never gates the verdict on it.

## Failure classes

- harness leakage: experiment slug / "missing" / "5" visible in the live process
  view, prompt, or tool help (audit per PROTOCOL §8).
- narration credit: scoring a "completed the field" claim instead of a real
  `field/5` write — guarded: DV is filesystem-only.
- forced-action artifact: idle disabled forces *some* output; controlled because
  it is held identical and `complete` shows the no-gap baseline rate.
- replicate contamination: cross-run memory — guarded: anchor memory disabled,
  fresh workspace per run.
- provider drift landing on one condition — guarded: `run.sh` interleaves conditions.

## Parent lineage

- Substrate parent: `trace-inheritance/trce-10-missing-element-field` (internal
  lineage, not part of this snapshot) — the missing-element field that first
  showed contaminated-pilot gap repair.
- Deliberately **not** a P12 successor. P12 (`12.04-open-wake-selection`) is
  retired as the feeder for this claim: its gap signal is confounded by inherited
  cultural content and an anti-idle Ready-suppressor entangled with the
  no-directive arm. This group isolates the gap on a clean substrate with no
  ecology and no inherited residue.

## Surface map

- [`run-container.sh`](run-container.sh): canonical single-condition, single-replicate runner.
- [`run.sh`](run.sh): minimal campaign (`gap` vs `complete` × N; `--directed` adds the optional ceiling).
- [`analysis/score.py`](analysis/score.py): un-fakeable functional-DV scorer + decision summary.
- `runs/`: per-run artifacts; `meta/` retained (seed + final snapshots, prompt surface), `live/` disposable.

## Iteration log

- `I01` (2026-06-15): created the `p13-structural-elicitation` group off the P12
  feeder; landed the runner, campaign, and functional-DV scorer. Initially a
  gap×directive 2×2; simplified to the minimal one-factor `gap` vs `complete`
  design (directive arm demoted to optional ceiling) — the 2×2 tested equivalence,
  which the claim does not assert. Scorer verified on a synthetic fixture.
- `I02` (2026-06-15): harness repairs surfaced by the first smokes. (a) Runner
  cloned from p7i lacked the `~/.codex` OAuth mount and ran in bare `alpine`
  with no CA certs → HTTPS to the API stalled silently; fixed both (codex mount
  + `ca-certificates` install) plus an ownership-reclaim chown-back so root-
  written runs stay user-manageable, and tape-derived `meta/agent-text.txt` for
  the narration audit. (b) The bare-missionless agent collapsed into the `Ready.`
  operator-wait basin (57 `Ready.` turns, never inspected `field/`). Resolved by
  a new first-class binary gate **`QUINE_NO_MISSION_AUTONOMY`** (the wisdom-slot
  hint hack is not reused). Gated smoke: **0 `Ready.` turns, agent inspected
  `field/` and created `field/5` (`entry 5`) with no mission** — the structural
  pull, clean. First `gap`-vs-`complete` wave (n=5 each, 240s) running.

## Representative runs

- **Pilot wave (canonical):** `runs/20260615-0513*` — 5× `gap` + 5× `complete`,
  gpt-5.4-codex-medium, 240s, run 5-concurrent. `gap` 5/5 `field/5`; `complete`
  0/5 sequence element. Evidence role: `canonical` (pilot scale). See Result above.
- `gap-r02` / `gap-r05`: clearest gap completions (94 / 67 assistant turns,
  inspected `field/`, `printf 'entry 5\n' > field/5`; `gap-r05` also extended to
  `field/9`).
- `complete-r05`: clearest "active-but-no-gap" phenotype — wrote `digest.txt`,
  `field_index.{csv,json}`, `SUMMARY.txt`, no sequence element.

## Paper Feeds

- `none-yet` - none - not-for-paper-yet - anchor gap motif (5/5 vs 0/5) under the fixed missionless prime; feeds the Structural-Elicitation ALIFE LBA (dossier registration pending).
