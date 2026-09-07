# M1 Same-Referent Detector — Tightening (Specificity Weighting)

Script: `/tmp/claude-1000/-home-kehao95-quine/5052fd24-5ffc-4613-9fde-f0a06a012c7d/scratchpad/m1_tightening.py`
Raw run log: `/tmp/claude-1000/-home-kehao95-quine/5052fd24-5ffc-4613-9fde-f0a06a012c7d/scratchpad/m1_tightening_output.txt`
Baseline (Phase-2, hard hub cutoff): `/tmp/claude-1000/-home-kehao95-quine/5052fd24-5ffc-4613-9fde-f0a06a012c7d/scratchpad/m1-evaluation.md` (archived: `/home/kehao95/quine/development/frontmatter-object-graph-feasibility/phase2-m1-evaluation-2026-07-02.md`)

Corpus unchanged from Phase-2: 14 files, 117 triples (grep-verified, assertion in
script). No `Paper/` files touched; nothing committed.

## Part A — specificity-weighted score

**Formula:** for a pair of files `(a, b)`, `score(a,b) = Σ 1/in-degree(t)` over
targets `t` shared by both files' edge sets (any relation type except
`home-family`, of which there are none). In-degree recomputed over the same
14-file / 117-triple graph. **No hard hub cutoff** — hubs are simply
down-weighted (`selection-signal-collapse`, in-degree 7, weighs 0.143;
`ai-self-replication`, in-degree 6, weighs 0.167), replacing Phase-2's
binary "exclude if in-degree > 3."

### Full ranked list (36 pairs with ≥1 shared target, out of C(14,2)=91)

| rank | score | pair | shared targets |
|---|---|---|---|
| 1 | **1.976** | 6m-corpus-as-ecology-seed ↔ research-brief | symmetry-trap-phase(.33), ai-self-replication(.17), body-ecology-coupling(.50), network-utility-selection(.50), selection-signal-collapse(.14), successor-morphology(.33) |
| 2 | **1.333** | encoding-recursion-productive-spiral ↔ idea-space-density-not-tier | exec-dependent-idea-lifecycle(.50), computational-epistemology(.33), mutual-encoding-network(.50) |
| 3 | 1.143 | 6m-corpus-as-ecology-seed ↔ exploit-niche-succession | symmetry-trap-phase(.33), unfakeable-selection(.50), ai-self-replication(.17), selection-signal-collapse(.14) |
| 4 | 0.833 | idea-space-density-not-tier ↔ negentropy-is-not-independence | collision-is-not-automated(.50), computational-epistemology(.33) |
| 5 | 0.810 | co-evolving-evaluator ↔ exploit-niche-succession | unfakeability-gradient(.50), ai-self-replication(.17), selection-signal-collapse(.14) |
| 5 | 0.810 | co-evolving-evaluator ↔ research-brief | drq-03-host-lineage-drq(.50), ai-self-replication(.17), selection-signal-collapse(.14) |
| 7 | 0.750 | idea-space-density-not-tier ↔ metadata-as-object-graph | never-paraphrase-fields(.50), compression-loss-experiment(.25) |
| 7 | 0.750 | idea-space-density-not-tier ↔ read-layer-prior | living-system(.50), compression-loss-experiment(.25) |
| 9 | 0.667 | carrier-pull-coupling ↔ drift-vs-selection-confound | aaai-c3-safety-non-inheritance(.50), ai-self-replication(.17) |
| ... | ≤0.643 | (25 more pairs, all `selection-signal-collapse`/`ai-self-replication`-only or single narrow target) | — full list in run log |

(Full 36-row list with every shared target and weight is in
`m1_tightening_output.txt`; only the top 9 are reproduced here since nothing
below rank 4 is relevant to the ground truth.)

### Ground-truth pairs — score + rank

| pair | verdict wanted | score | rank / 36 |
|---|---|---|---|
| 6m-corpus-as-ecology-seed ↔ research-brief | **TP** (must survive) | **1.976** | 1 |
| encoding-recursion-productive-spiral ↔ idea-space-density-not-tier | **TP** (must survive) | **1.333** | 2 |
| 6m-corpus-as-ecology-seed ↔ exploit-niche-succession | **FP** (must drop) | 1.143 | 3 |
| idea-space-density-not-tier ↔ negentropy-is-not-independence | **FP** (must drop) | 0.833 | 4 |

### Threshold separability

**min(TP) = 1.333 > max(FP) = 1.143.** A single threshold in the open-closed
interval **(1.143, 1.333]** — e.g. **1.2** or **1.25** — keeps both true
positives and drops both false positives. Checked against the **full 36-pair
list**, not just the 4 GT pairs: only rank 1 and rank 2 clear 1.238 (or any
value in the separating interval). No other pair in the whole corpus sneaks
in above the line — the resulting candidate list is exactly
`{6m-corpus ↔ research-brief, encoding-recursion ↔ idea-space-density}`, i.e.
**both TPs, zero FPs, zero surprises.**

**Margin quality (is this a coincidence?).** The gap (0.19, ~14% of TP2's
score) is real structural separation, not curve-fitting luck: FP1's raw
shared-target count (4) is actually *higher* than TP2's (3), but two of FP1's
four targets are the corpus's two biggest hubs
(`selection-signal-collapse` w=0.143, `ai-self-replication` w=0.167) — exactly
the paradigm-neighborhood furniture the Phase-2 diagnosis identified as the
false-positive mechanism — while every one of TP2's three shared targets is
narrow (w≥0.333). This is precisely what specificity weighting is designed to
punish: raw shared-count no longer wins, weighted specificity does.

### Blind-test re-confirmation

Removed `dedup-candidate` + `follow-up-to` (2 triples), recomputed in-degree
on the reduced 115-triple graph, rescored `6m-corpus-as-ecology-seed` ↔
`research-brief`: **score unchanged, 1.976** (identical shared-target set).
This is expected — neither deleted edge was one of this pair's 6 shared
targets — but it confirms the blind-test pair still clears the threshold
by the same wide margin (1.976 vs. threshold ~1.2) even with the answer-key
edge gone, using purely structural signal.

**Part A verdict: specificity weighting ALONE achieves full separation.**
Clean list, both FPs gone, both TPs kept, blind-test pair recovered with an
unchanged score. Part B is not required by the task's own branching logic,
but was run anyway below as an independent robustness check (does an
unrelated signal — consequence overlap — agree with the score-based cut, or
was the 0.19 margin a fluke?).

## Part B — consequence-overlap (robustness check only, not required)

Reused existing extraction for the two discovered-sparks already classified in
`development/compression-loss-experiment-2026-07-02/artifacts/E-consequence-classes.md`:

- **6m-corpus-as-ecology-seed (S1):** licenses re-hydrating the 61-row 6M
  corpus as standing variation in a demand arena; predicts character
  displacement by body class; **prohibits** claiming this closes the
  heritability condition (2 of 3 met); requires reporting adapter/inference
  cost vs. a fresh-induction baseline.
- **exploit-niche-succession (C3):** predicts unanchored evaluator
  co-evolution collapses to collusion; predicts a "cured" exploit channel
  relocates to the next-cheapest niche; predicts a reclassified
  "exec-dependent" channel still collapses if any sub-signal stays
  narration-satisfiable; predicts optimal audit rate falls toward zero as
  cheat-rate falls.

For the four files not covered by that scheme (all seed-cluster / brief, not
discovered-sparks), I extracted consequences myself, 2-4 bullets from
thesis+body — **flagged: this is my own LLM Level-0 judgment**, not a
mechanical extraction:

- **idea-space-density-not-tier:** predicts collision cross-section keeps
  degrading without an installed density-physics mechanism; **prohibits**
  installing any density/dedup mechanism on a narrative-layer (text/embedding)
  distance metric; licenses deferring any physics fix until the metric
  problem is answered at the exec layer; predicts the first-deposit QD/tier
  frame is subsumed, not falsified.
- **negentropy-is-not-independence:** **prohibits** crediting pure
  generation+reflection loops as surprise-independent verification of the
  repo's own blind spots; predicts internally-consistent, locally-compressing
  "order" can still be externally false without an independent check;
  licenses crediting real compression only against mechanical/exec
  predicates; predicts the stakes of the missing independence channel rise,
  not fall, as models get more capable.
- **encoding-recursion-productive-spiral:** **prohibits** labeling a new turn
  of the vertical encoding recursion as pathological/treadmill without
  checking it; predicts a turn that adds structure with ~zero compression
  gain is bloat (a positive, checkable discriminator); licenses continuing to
  re-enter the recursion rather than shutting it down by policy; predicts
  vertical re-encoding and horizontal collision are one co-search, not
  opposed activities.
- **research-brief (11-diversification):** predicts one fixed base model
  under varied material/demand conditions yields 8+ qualitatively distinct
  successor morphologies; licenses using that operator deliberately as an
  agent-population-generation mechanism; **prohibits** conflating it with
  weight/prompt/architecture diversity or evolutionary population search;
  predicts body-ecology-coupling is stageable via network-utility-selection
  over those morphologies.

### Overlap per GT pair

| pair | consequence overlap? |
|---|---|
| 6m-corpus ↔ research-brief (TP) | **YES, strong** — both license/predict the same operationalized test: stage body-ecology-coupling via network-utility-selection over successor-morphology variation, both flag heritability (condition 3) as the open gap. (research-brief literally `spawns` 6m-corpus.) |
| encoding-recursion ↔ idea-space-density (TP) | **YES, graded/productive-disagreement** — both license a test of the *same* substrate's growth (the vertical encoding recursion / idea-space concentration) and propose competing discriminators for when it is healthy (marginal-compression-per-turn vs. global density-physics), i.e. they share a prediction *target* while disagreeing on the verdict — exactly the "genuine related-but-distinct" shape the task specifies. |
| 6m-corpus ↔ exploit-niche-succession (FP) | **NO** — embodiment/heritability-assay consequences vs. evaluator-collusion/exploit-relocation consequences share no predicted test, prohibition, or action. |
| idea-space-density ↔ negentropy-is-not-independence (FP) | **NO** — density/collision-cross-section physics vs. epistemic status (order-vs-truth) of self-reflective loops share no predicted test, prohibition, or action. |

**Part B agrees with Part A on all four pairs** — an independent signal
(licensed consequences, extracted from prose, no reference to the edge graph)
draws the identical TP/FP line the specificity-weighted score drew. This
corroborates that the Part-A threshold margin is a real structural fact about
this corpus, not a numerical coincidence.

## Verdict

1. **Does specificity weighting ALONE separate TP from FP?** **YES.** Score
   `Σ 1/in-degree(shared target)`, threshold in (1.143, 1.333] (recommend
   **1.2**), applied over the full 14-file/117-triple graph: clean list,
   exactly `{6m-corpus↔research-brief, encoding-recursion↔idea-space-density}`
   survive, both known FPs drop below threshold, blind-test pair recovers at
   unchanged score 1.976 with the answer-key edge deleted. This is the
   recommended detector — it is simpler than adding a second signal and it
   already achieves full separation on the available ground truth.

2. **Was Part B (consequence-overlap) needed?** **No** — Part A alone
   achieves full separation, so per the task's own branching Part B was not
   required. It was run anyway as an independent check and **confirms** the
   Part-A cut on all four pairs (agreement, not disagreement), which is
   evidence the 0.19-score margin is structural rather than incidental.

3. **False-positive rate vs. Phase-2 baseline.** Reframed like Phase-2's own
   Step 5 (excluding the pair already known via a declared `dedup-candidate`
   edge, i.e. treating TP1 as "already resolved" bookkeeping rather than a
   newly-surfaced candidate): the **undetected-candidate set at threshold 1.2
   is exactly 1 pair** — `encoding-recursion ↔ idea-space-density` — and it is
   a true positive (productive-disagreement), not a false positive.
   **New false-positive rate: 0/1 (0%), vs. Phase-2 baseline 2/3 (67%).**
   Equivalently, scored against all 4 GT pairs directly (not excluding the
   declared pair): 2/2 TPs kept, 0/2 FPs kept — **0% FP rate** either way you
   slice it.

4. **Recommended M1 detector spec:**
   - **Signal:** specificity-weighted shared-target score,
     `score(a,b) = Σ_{t ∈ shared(a,b)} 1/in-degree(t)`, in-degree computed
     over the same corpus's full edge graph (all relation types except
     `home-family`). Replaces the Phase-2 hard hub cutoff (in-degree > 3
     excluded) outright — no cutoff needed, weighting subsumes it.
   - **Threshold:** ≥ 1.2 (any value in (1.143, 1.333] works on this corpus;
     1.2 leaves comfortable margin on both sides).
   - **Consequence-overlap second signal:** **not needed.** Adding it would
     reintroduce an LLM Level-0 extraction step (as flagged above, the 4
     seed-cluster/brief files required manual consequence extraction — no
     existing mechanical source covers them) for zero marginal precision
     gain on this ground truth.
   - **Cost:** **pure-mechanical.** The whole detector is arithmetic over the
     already-parsed edge graph — same inputs Phase-2 used (frontmatter
     `edges:` blocks), no prose reading, no LLM judgment call, fully
     re-runnable as a CI-style check. This is strictly cheaper than Phase-2's
     hard-cutoff version (same parse, one extra division per shared target)
     and does not require the consequence-extraction step Part B would have
     needed had Part A failed.

**Caveat carried forward from Phase-2, unchanged by this tightening:** the
external-citation asymmetry gap (files that cite external work only in prose,
not via `external-occupancy` edges, contribute nothing to any edge-based
detector, weighted or not) still applies here and was not addressed by this
task.
