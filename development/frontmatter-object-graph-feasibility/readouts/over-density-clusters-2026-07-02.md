# Over-Density (Same-Referent) Clusters — 2026-07-02

**GENERATED VIEW — regenerate by running `python3 same_referent_detector.py`
(in this directory); not a canonical store, per read-layer-prior.**

Scope: full STEP-1-patched `Paper/incubator/**/*.md` corpus (101 in-scope
files, 1300 edge triples). Method: specificity-weighted M1 same-referent
score, `score(a,b) = sum over targets shared by both files of weight(t)`,
threshold `>= 1.2` (validated on the 14-file pilot to separate true from
false positives).

## The guard

The unguarded score (`weight(t) = 1/in-degree(t)` for every shared target,
regardless of edge type) produces a demonstrated false positive: three
`divergent-directions/*/research-brief.md` files (06-phase-diagram-capability,
14-bureaucratization, 15-reaction-norm-benchmark) score >= 1.2 pairwise
purely by sharing generic external-arXiv citations (Reflexion, Generative
Agents, Voyager, MemGPT / BIG-bench, HELM, ChatbotArena, AgentBench) — well
known background citations any two ML-agent-subfield literature surveys
would both cite, not evidence of a shared internal claim.

**Chosen rule:** a shared target reached by `external-occupancy` from BOTH
files in a pair contributes only `0.5 * 1/in-degree(t)` to the score
(instead of the full `1/in-degree(t)` that internal structural edges —
`depends-on`, `grounds-in`, `relates-to`, `spawns`, `originates-in`,
`dedup-candidate`, etc. — always get). In-degree itself is always the true,
un-discounted count of citing files; only the *contribution* of an
external-occupancy-mediated shared target is fractional.

**Why 0.5, not full exclusion (weight 0):** full exclusion also kills a
real true-positive cluster (carrier-memetics-mechanism-correction <->
cultural-attraction-...-weismann-reading), whose score is carried partly by
3 shared external citations and partly by 4 shared repo-node targets — at
weight 0 its score drops to 0.638, below threshold. A weight sweep
(`{1.0, 0.6, 0.5, 0.42, 0.25, 0.0}`, see script `__main__` output) shows 0.5
is the value that simultaneously (a) drops the literature-only cluster
(2.0 -> 1.0), (b) keeps both real clusters above threshold, and (c) leaves
the 14-file pilot ground truth untouched (none of the 4 pilot-truth pairs
share any external-occupancy target at all, so the guard cannot move them
regardless of weight chosen).

## Before -> after

| Pair | Before (unguarded) | After (guarded, w=0.5) | Verdict |
|---|---|---|---|
| `06-phase-diagram-capability` <-> `14-bureaucratization` | 2.000 | 1.000 | **DROPPED** (false positive) |
| `06-phase-diagram-capability` <-> `15-reaction-norm-benchmark` | 2.000 | 1.000 | **DROPPED** (false positive) |
| `carrier-memetics-mechanism-correction` <-> `cultural-attraction-rival-...-weismann-reading` | 1.971 | 1.305 | **KEPT** (true cluster) |
| `sematectonic-vs-marker-trace-taxonomy` <-> `web-scale-stigmergic-eval-contamination-bridge` | 1.525 | 1.525 (unaffected — 100% repo-node basis) | **KEPT** (true cluster) |

Pilot ground truth (checked on its own 14-file subgraph, guarded rule):

| Pair | Score | Threshold check |
|---|---|---|
| TP1: `6m-corpus-as-ecology-seed` <-> `11-diversification` research-brief (blind-test pair) | 1.976 | >= 1.2, caught (correct) |
| TP2: `encoding-recursion-productive-spiral` <-> `idea-space-density-not-tier` | 1.333 | >= 1.2, caught (correct) |
| FP1: `6m-corpus-as-ecology-seed` <-> `exploit-niche-succession` | 1.143 | < 1.2, uncaught (correct) |
| FP2: `idea-space-density-not-tier` <-> `negentropy-is-not-independence` | 0.833 | < 1.2, uncaught (correct) |

All three constraints hold simultaneously: the guard drops the false
positive, keeps both plausible true clusters, and leaves the pilot ground
truth (both true and false pairs) exactly where it was.

## Clusters after the guard (2 total, 4 files, ~4%)

### Cluster 1 — `sematectonic-vs-marker-trace-taxonomy` <-> `web-scale-stigmergic-eval-contamination-bridge`

Score 1.525. Shared targets: `candidates/stigmergic-trace-causality`,
`experiments/active/benchmarks/bench-06-persistent-cleanroom-c-compiler/workspace`,
`theory/gaps/carrier-selection`, `theory/gaps/function-vs-narrative`,
`theory/objects/carrier`, `theory/objects/stigmergy-as-carrier`.

**Judgment: true positive.** Both files explicitly extend the same
stigmergy/marker-trace theoretical toolkit (`stigmergy-as-carrier`,
`trace-inheritance`) toward adjacent applications — same-referent family,
correctly surfaced. Entirely carried by repo-node (not literature) overlap,
so the guard is a no-op here — this cluster was never at risk of being a
citation artifact.

### Cluster 2 — `carrier-memetics-mechanism-correction` <-> `cultural-attraction-rival-to-the-discrete-replicator-weismann-reading`

Score 1.305 (down from 1.971 pre-guard). Shared targets:
`arxiv:2403.08882`, `arxiv:2412.10270`, `arxiv:2605.23054` (down-weighted
0.5x), `papers/cultural-evolution-past-brain-world-barrier`,
`theory/objects/carrier`, `theory/objects/exec-barrier`,
`theory/views/substrate`.

**Judgment: true positive, worth a human collision check.** Both files are
from the same 2026-07 backlog-maturation pass and both propose mechanisms
for carrier-mediated cultural/memetic inheritance — genuinely adjacent, not
obviously a duplicate, but same neighborhood on real theory-node substance,
not just shared literature.

## What dropped and why it's correct to drop it

`06-phase-diagram-capability` / `14-bureaucratization` / `15-reaction-norm-benchmark`
are three independent literature reviews over the same ML-agent subfield on
genuinely different research questions (capability phase diagrams vs.
organizational bureaucratization vs. reaction-norm benchmarking). Their
only shared targets are generic background citations; after the guard,
their combined connected component disappears entirely (both pairwise
scores fall to 1.000, below 1.2).

## Count

**2 clusters, 4 files (~4% of the 101-file corpus), both true positives**
under the guarded rule (post-guard, pre-guard was 3 clusters / 7 files with
1 of the 3 a demonstrated false positive).

## Tier-2 caveat

This detector is structurally blind to genuinely empty regions with no
node at all — same-referent detection over a relation graph can only
compare files that already have edges into shared targets; it cannot
surface duplication in ideas that never got any edge at all, nor can it
detect a would-be duplicate pair that happens to share zero repo-node
targets (e.g. two files independently citing entirely disjoint literature
about the same idea, worded differently, would score 0 regardless of any
weighting scheme).
