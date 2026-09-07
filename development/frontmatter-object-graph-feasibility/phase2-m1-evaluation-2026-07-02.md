# M1 Same-Referent Detector — Decisive Behavior-Gate Evaluation

Script: `/tmp/claude-1000/-home-kehao95-quine/5052fd24-5ffc-4613-9fde-f0a06a012c7d/scratchpad/m1_eval.py`
Raw run log: `/tmp/claude-1000/-home-kehao95-quine/5052fd24-5ffc-4613-9fde-f0a06a012c7d/scratchpad/m1_output.txt`

Corpus located mechanically via `grep -rl "^edges:" Paper/incubator/`: exactly the
14 files named in the task (8 seed-cluster + 5 discovered-sparks + 1 brief). The
script asserts this file list matches the grep output — it does.

## Step 1 — Parse

**Total triples: 117.** Sanity check against the task's "~28 + ~89" estimate: the
8 seed-cluster files sum to exactly **28** triples; the 5 discovered-sparks + 1
brief sum to exactly **89** triples. Exact match.

Per-file counts:

| file | triples |
|---|---|
| incubator/collision-is-not-automated.md | 3 |
| incubator/consequence-identity.md | 3 |
| incubator/encoding-recursion-productive-spiral.md | 3 |
| incubator/idea-space-density-not-tier.md | 10 |
| incubator/metadata-as-object-graph.md | 2 |
| incubator/negentropy-is-not-independence.md | 3 |
| incubator/never-paraphrase-fields.md | 2 |
| incubator/read-layer-prior.md | 2 |
| **seed-cluster subtotal** | **28** |
| discovered-sparks/6m-corpus-as-ecology-seed.md | 11 |
| discovered-sparks/carrier-pull-coupling-as-dual-use-attack-surface.md | 21 |
| discovered-sparks/co-evolving-evaluator-as-the-general-antidote.md | 11 |
| discovered-sparks/drift-vs-selection-confound-in-single-descendant-lineage-chains.md | 13 |
| discovered-sparks/exploit-niche-succession.md | 17 |
| divergent-directions/11-diversification/research-brief.md | 16 |
| **sparks+brief subtotal** | **89** |

By relation type: `depends-on` 33, `external-occupancy` 19, `relates-to` 17,
`grounds-in` 15, `spawns` 13, `distinct-from` 7, `originates-in` 5,
`tensions-with` 3, `refines` 1, `superseded-by` 1, `dedup-candidate` 1,
`follow-up-to` 1, `specializes` 1. `home-family` edges: **0** — confirms the
design's claim that this edge type was dismantled, not merely discounted.

## Step 2 — In-degree / hubs

Hubs (in-degree > 3, excluded from shared-target counting):

| target | in-degree | sources |
|---|---|---|
| `theory/objects/selection-signal-collapse` | 7 | consequence-identity, 6m-corpus, co-evolving-evaluator, drift-vs-selection, exploit-niche-succession, research-brief, idea-space-density |
| `paradigms/ai-self-replication` | 6 | 6m-corpus, carrier-pull-coupling, co-evolving-evaluator, drift-vs-selection, exploit-niche-succession, research-brief |
| `theory/references/compression-loss-experiment-2026-07-02` | 4 | idea-space-density, metadata-as-object-graph, never-paraphrase-fields, read-layer-prior |

These are exactly the kind of paradigm/experiment hub the design predicted
(`selection-signal-collapse`, the paradigm object itself). Discounting them is
load-bearing: without it, every one of the 5 discovered-sparks files would
pairwise-collide on `selection-signal-collapse` + `paradigms/ai-self-replication`
alone (10 spurious pairs), which is exactly the false-positive mode the design
document's Phase-0 probe reported.

Non-hub targets with in-degree 2-3 (the candidate basis) — 16 targets, full list
in the run log; notably includes `theory/gaps/body-ecology-coupling` (in-degree
2: 6m-corpus, research-brief) and `candidates/symmetry-trap-phase` (in-degree 3).

## Step 3 — Same-referent candidates (full graph, real edges)

Raw pairs with ≥2 shared non-hub targets: **4**, out of C(14,2)=91 possible pairs.

Applying the guards:

- **distinct-from guard**: 2 `distinct-from` edges resolve to another file in
  this 14-file corpus (`collision-is-not-automated` ↔ `idea-space-density-not-tier`,
  `consequence-identity` ↔ `idea-space-density-not-tier`). **Neither** of these
  pairs appears in the raw candidate-pair list anyway (they don't share ≥2
  non-hub targets), so the guard fires 0 times on real data — it did not need to
  suppress anything here, but it is correctly implemented and would have fired
  had the overlap existed.
- **already-resolved (declared relation exists)**: **1** pair —
  `6m-corpus-as-ecology-seed` ↔ `11-diversification/research-brief`, linked by a
  direct `dedup-candidate` edge. Shared non-hub targets:
  `candidates/symmetry-trap-phase`, `theory/gaps/body-ecology-coupling`,
  `theory/objects/network-utility-selection`, `theory/objects/successor-morphology`
  (4 shared targets — the richest overlap in the corpus). Correctly bucketed as
  "already resolved," not surfaced as an undetected candidate.

**Surviving candidate pairs (undetected, guards applied): 3**

1. `6m-corpus-as-ecology-seed` ↔ `exploit-niche-succession`
   shared: `candidates/symmetry-trap-phase`, `papers/unfakeable-selection`
2. `encoding-recursion-productive-spiral` ↔ `idea-space-density-not-tier`
   shared: `_design/proposals/exec-dependent-idea-lifecycle`,
   `philosophy/computational-epistemology`, `philosophy/mutual-encoding-network`
3. `idea-space-density-not-tier` ↔ `negentropy-is-not-independence`
   shared: `incubator/collision-is-not-automated`, `philosophy/computational-epistemology`

## Step 4 — THE BLIND TEST

Removed all `dedup-candidate` and `follow-up-to` edges (2 triples: the one
`dedup-candidate` edge above, and one `follow-up-to` edge in
`carrier-pull-coupling-as-dual-use-attack-surface.md` → an external paper). 115
triples remain.

**Question:** does `6m-corpus-as-ecology-seed` ↔
`divergent-directions/11-diversification/research-brief` still surface from
structural edges alone (`grounds-in`, `external-occupancy`, `distinct-from`,
`depends-on`, `specializes`, `relates-to`, `originates-in`)?

**Result: YES — it surfaces**, with 4 shared non-hub targets carrying it over
threshold:

- `theory/gaps/body-ecology-coupling` (via `grounds-in` on both sides — in-degree 2)
- `theory/objects/successor-morphology` (via `depends-on` on both — in-degree 3, non-hub)
- `theory/objects/network-utility-selection` (via `depends-on` on both — in-degree 2)
- `candidates/symmetry-trap-phase` (via `relates-to` on both — in-degree 3, non-hub)

No `distinct-from` edge exists between this pair, so the guard does not suppress
it. This is exactly the pair the phase-0 probe predicted, recovered from
structural edges with the answer-key (`dedup-candidate`) edge deleted.

**One caveat that qualifies the win.** The phase-0 probe's prediction was that
the pair shares "the narrow gap `theory/gaps/body-ecology-coupling` **and
non-generic external cites**." The gap half of that prediction materialized as
an edge and did the real work. The external-cite half did **not**: `6m-corpus-
as-ecology-seed` carries **zero** `external-occupancy` edges (it discusses
related literature, if at all, only in prose), while `research-brief` carries
two (`arxiv:2603.16910`, `arxiv:2409.00853`). Those never became a shared edge
because only one side encoded them. The pair still surfaced — carried by
`grounds-in`/`depends-on`/`relates-to` to shared theory objects instead — but
this is a concrete, real gap the corpus exposed: **external citations that
exist in only one of two same-referent files' prose (not promoted to
`external-occupancy` edges on both) contribute nothing to same-referent
detection**, even when they would have been corroborating evidence. For a
purely literature-adjacent pair (no shared internal theory-object edges at
all), this failure mode would be fatal — the pair would not surface at all. It
happened not to be fatal here only because this pair also shares internal
theory-object edges.

## Step 5 — Judging the 3 surviving candidate pairs (false-positive rate)

Read each pair's thesis/matured-claim prose directly (not the edges) and
classified:

1. **`6m-corpus-as-ecology-seed` ↔ `exploit-niche-succession` → FALSE POSITIVE.**
   6m-corpus is about re-hydrating a body-morphology corpus as a common-garden
   assay for the *body-ecology-coupling* gap (embodiment/diversity domain).
   exploit-niche-succession is about whether the cheapest viable exploit
   relocates through an ordered sequence of niches as evaluator grounding rises
   (a Goodhart-ladder/security domain). They share only paradigm-neighborhood
   citations (`symmetry-trap-phase`, `unfakeable-selection` — one file
   `relates-to` these, the other `distinct-from`s them) and the "ecological
   succession" metaphor. Different problems, connected by shared paradigm
   furniture, not by a shared referent.

2. **`encoding-recursion-productive-spiral` ↔ `idea-space-density-not-tier` →
   RELATED-BUT-DISTINCT (genuine productive-disagreement candidate).** Both are
   about the idea-layer's own meta-dynamics and share real grounding
   (`philosophy/mutual-encoding-network`, `philosophy/computational-epistemology`,
   `_design/proposals/exec-dependent-idea-lifecycle`). encoding-recursion argues
   the vertical-recursion loop is a *productive spiral* (each turn compresses,
   nothing pathological). idea-space-density argues the idea ecology has a
   *pathology* (uncontrolled concentration, degraded collision cross-section)
   needing a physics fix. These are two angles on the same substrate — one
   optimistic about the recursion, one diagnosing a failure mode in the same
   space — a legitimate "collide these" invitation, not a same-claim dedup and
   not noise.

3. **`idea-space-density-not-tier` ↔ `negentropy-is-not-independence` → FALSE
   POSITIVE.** idea-space-density is about density physics / collision
   cross-section / the metric problem. negentropy-is-not-independence is about
   the epistemic distinction between negentropy (order) and surprise-
   independence (truth-tracking capacity) in self-reflective loops. Shared
   targets are a common philosophy grounding doc
   (`philosophy/computational-epistemology`) and a third idea file
   (`incubator/collision-is-not-automated`) that one references via
   `distinct-from` and the other via `relates-to` — i.e., the "shared referent"
   here is two different stances toward a third object, not a shared claim
   between the pair itself. Different problems under a common philosophical
   umbrella.

**False-positive rate: 2/3 (67%)** on the surviving (undetected) candidate set.
0/3 are true dedup pairs; 1/3 is a genuine related-but-distinct collision
candidate. (n=3 — small sample, stated plainly rather than dressed up.)

## Verdict

**(1) Clean candidate list, low false-positive, on real edges?** **Partially.**
The guards work exactly as designed and are load-bearing: hub-discounting kills
what would otherwise be 10+ spurious pairs among the 5 discovered-sparks files
alone (all sharing the paradigm + `selection-signal-collapse` hubs), and the
declared-relation bucket correctly separates the one already-known dedup pair
from "undetected" candidates. But on the 3 pairs the detector does surface as
*undetected* candidates, the false-positive rate is 2/3 — worse than the
"low false-positive" bar the design's Phase-1 gate asks for, though on a very
small n. The one true positive in that set is not a dedup but a productive-
disagreement collision, which is itself a legitimate M2 outcome per the design
(not every surfaced pair needs to be a dedup to be useful) — but the pure "same
problem" precision is weak on this sample.

**(2) Does the blind test recover the known `6m-corpus ↔ 11-diversification`
pair from structural edges alone?** **YES.** With the answer-key
`dedup-candidate` edge deleted, the pair still surfaces via 4 shared non-hub
structural targets (`grounds-in` to the same narrow gap, `depends-on` to the
same two theory objects, `relates-to` to the same candidate). This is the
decisive positive result: the un-fakeable edge signal — not the human's
already-written conclusion — reproduces the human prose judgment for the one
pair we have ground truth on. The design's central bet (shared edges, not
wording, carry the same-referent signal) holds on this instance.

**(3) Concrete gap the real data exposed.** External citations are asymmetric
across same-referent pairs: `6m-corpus-as-ecology-seed` has zero
`external-occupancy` edges even though it is plausibly literature-adjacent,
while its true dedup partner has two. The phase-0 probe's prediction that
"non-generic external cites" would help carry this pair was not realized as
edges — the pair surfaced despite that channel being empty, not because of it.
For a hypothetical pair whose *only* overlap is literature (no shared internal
theory-object edges), this same gap would mean the detector misses the pair
entirely. This argues for a follow-up check (not built here): flag idea files
that cite external work only in prose, as candidates for retroactive
`external-occupancy` promotion, before trusting M1's coverage on
literature-adjacent sparks.

**Bottom line:** the mechanical detector, run over real edges with the design's
two guards, recovers the one known ground-truth same-referent pair from
structural signal alone (a real, positive behavior-gate result for the
data-structure migration) — but its false-positive rate on the small set of
*newly surfaced* candidates (2/3) shows the in-degree discount alone is not
sufficient to reach "clean list" on this corpus; a human collision judgment
call (M2) is still doing real precision work downstream of M1, exactly as the
design's own phasing (M1 visibility → M2 human-judged collision, never
automatic merge) assumes.
