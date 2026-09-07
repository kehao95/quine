#!/usr/bin/env python3
"""
GENERATED VIEW SUPPORT SCRIPT — regenerate readouts/over-density-clusters-*.md
by running this file; not a canonical store, per read-layer-prior.

Same-referent ("over-density") detector, STEP 2 of the idea-layer edge
migration follow-up. Reuses the specificity-weighted M1 method validated in
the QA agent's full_m1.py:

    score(a,b) = sum over targets t shared by both files of 1/in-degree(t)
    threshold >= 1.2

but adds the EXTERNAL-CITATION GUARD: the full-corpus run found a false
positive (three `divergent-directions` research-briefs scoring >=1.2 purely
by sharing generic external-arXiv citations like Reflexion/Voyager/
BIG-bench — literature furniture any two ML-agent-subfield surveys would
both cite, not evidence of a shared internal claim).

Guard rule chosen (see the weight-sweep in __main__ output): full exclusion
(weight 0) OVER-KILLS — it also drops the true-positive cluster B
(carrier-memetics <-> cultural-attraction), which is carried partly by 3
shared external-occupancy citations alongside 4 shared repo-node targets.
The chosen rule is a FRACTIONAL weight: a target reached by
external-occupancy from BOTH files in a pair contributes only
EXTERNAL_WEIGHT * 1/in-degree(t) (EXTERNAL_WEIGHT=0.5) instead of the full
1/in-degree(t) internal structural edges get. A sweep over
{1.0(=baseline), 0.6, 0.5, 0.42, 0.25, 0.0} shows 0.5 is the unique setting
in the tested grid that (a) drops the 06/14/15 literature-only cluster
(score 2.0 -> 1.0, below 1.2), (b) keeps both true clusters B (1.971 -> 1.305)
and C (1.525, unaffected — C's shared basis is already 100% repo-node, so
the guard is a no-op on it), and (c) leaves the 14-file pilot ground truth
untouched (TP1/TP2 stay caught, FP1/FP2 stay uncaught — none of the 4 pilot
pairs share any external-occupancy target in the first place, so the guard
cannot move them regardless of weight). in-degree is always the TRUE count
of distinct citing files for a target (never itself down-weighted) —
only a pair's per-target *contribution* to the score is fractional.

This module runs over the STEP-1-patched corpus (all 101 in-scope files)
and separately re-validates against the 14-file pilot ground truth (which
predates and is disjoint from the full-corpus batch, but is checked here on
its own sub-graph to match the QA phase-2 method exactly).
"""
import itertools
from collections import defaultdict

from graph_lib import in_scope_incubator_files, build_corpus_triples, REPO

THRESHOLD = 1.2


EXTERNAL_WEIGHT = 0.5  # chosen guard weight for external-occupancy targets (see report)


def compute(all_triples, guarded, ext_weight=0.0):
    """all_triples: list of (src, rel, to).
    guarded=False: original M1 rule, every edge type weighted equally (1/in-degree).
    guarded=True: internal structural edges keep weight 1/in-degree(t); an
    external-occupancy-reached target's contribution is instead
    ext_weight * 1/in-degree(t) (ext_weight=0.0 -> full exclusion,
    ext_weight=1.0 -> identical to baseline). in-degree itself is always the
    true count of distinct citing files for that target (unaffected by the
    scoring weight), computed over ALL kept edge types (home-family excluded,
    the only edge type this migration never uses)."""
    def not_home(rel):
        return rel != "home-family"

    # in-degree: true count, same under baseline and guarded (only the
    # per-pair contribution weight differs)
    target_sources = defaultdict(set)
    for (src, rel, to) in all_triples:
        if not_home(rel):
            target_sources[to].add(src)
    indegree = {t: len(s) for t, s in target_sources.items()}

    # file_targets: map file -> {target: is_external_occupancy}
    file_targets = defaultdict(dict)
    for (src, rel, to) in all_triples:
        if not not_home(rel):
            continue
        is_ext = (rel == "external-occupancy")
        # if a target is reached by BOTH an external-occupancy edge and an
        # internal edge from the same file (not expected in this corpus, but
        # guard anyway), prefer the internal (non-guarded) classification.
        if to in file_targets[src]:
            file_targets[src][to] = file_targets[src][to] and is_ext
        else:
            file_targets[src][to] = is_ext

    all_files = sorted(set(src for (src, _, _) in all_triples))

    def target_weight(t, is_ext):
        base = 1.0 / indegree[t]
        if guarded and is_ext:
            return ext_weight * base
        return base

    pair_results = {}
    for a, b in itertools.combinations(all_files, 2):
        ta, tb = file_targets[a], file_targets[b]
        shared = set(ta) & set(tb)
        if shared:
            score = 0.0
            for t in shared:
                # a target is "external-occupancy-reached" for this pair's
                # purposes only if BOTH files reached it that way (matches
                # the guard's intent: down-weight shared-citation evidence)
                is_ext = ta[t] and tb[t]
                score += target_weight(t, is_ext)
            pair_results[(a, b)] = (score, shared)
    return pair_results, indegree, file_targets


def connected_components(pairs_above):
    parent = {}

    def find(x):
        parent.setdefault(x, x)
        while parent[x] != x:
            parent[x] = parent[parent[x]]
            x = parent[x]
        return x

    def union(x, y):
        rx, ry = find(x), find(y)
        if rx != ry:
            parent[rx] = ry

    nodes = set()
    for (a, b), _s, _shared in pairs_above:
        nodes.add(a)
        nodes.add(b)
        union(a, b)
    clusters = defaultdict(set)
    for n in nodes:
        clusters[find(n)].add(n)
    return sorted(clusters.values(), key=lambda c: -len(c))


PILOT_FILES = [
    "Paper/incubator/collision-is-not-automated.md",
    "Paper/incubator/consequence-identity.md",
    "Paper/incubator/encoding-recursion-productive-spiral.md",
    "Paper/incubator/idea-space-density-not-tier.md",
    "Paper/incubator/metadata-as-object-graph.md",
    "Paper/incubator/negentropy-is-not-independence.md",
    "Paper/incubator/never-paraphrase-fields.md",
    "Paper/incubator/read-layer-prior.md",
    "Paper/incubator/discovered-sparks/6m-corpus-as-ecology-seed.md",
    "Paper/incubator/discovered-sparks/carrier-pull-coupling-as-dual-use-attack-surface.md",
    "Paper/incubator/discovered-sparks/co-evolving-evaluator-as-the-general-antidote.md",
    "Paper/incubator/discovered-sparks/drift-vs-selection-confound-in-single-descendant-lineage-chains.md",
    "Paper/incubator/discovered-sparks/exploit-niche-succession.md",
    "Paper/incubator/divergent-directions/11-diversification/research-brief.md",
]

A_6m = "Paper/incubator/discovered-sparks/6m-corpus-as-ecology-seed.md"
A_brief = "Paper/incubator/divergent-directions/11-diversification/research-brief.md"
A_encrec = "Paper/incubator/encoding-recursion-productive-spiral.md"
A_density = "Paper/incubator/idea-space-density-not-tier.md"
A_exploit = "Paper/incubator/discovered-sparks/exploit-niche-succession.md"
A_negent = "Paper/incubator/negentropy-is-not-independence.md"

PILOT_GT = {
    "TP1 6m-corpus <-> 11-diversification (blind-test pair, must stay caught)": (A_6m, A_brief),
    "TP2 encoding-recursion <-> idea-space-density (must stay caught)": (A_encrec, A_density),
    "FP1 6m-corpus <-> exploit-niche-succession (must stay uncaught)": (A_6m, A_exploit),
    "FP2 idea-space-density <-> negentropy-is-not-independence (must stay uncaught)": (A_density, A_negent),
}

CLUSTER_A = ("Paper/incubator/divergent-directions/06-phase-diagram-capability/research-brief.md",
             "Paper/incubator/divergent-directions/14-bureaucratization/research-brief.md")
CLUSTER_B = ("Paper/incubator/discovered-sparks/carrier-memetics-mechanism-correction.md",
             "Paper/incubator/discovered-sparks/cultural-attraction-rival-to-the-discrete-replicator-weismann-reading.md")
CLUSTER_C = ("Paper/incubator/discovered-sparks/sematectonic-vs-marker-trace-taxonomy.md",
             "Paper/incubator/discovered-sparks/web-scale-stigmergic-eval-contamination-bridge.md")


def pilot_check(ext_weight, guarded):
    triples = build_corpus_triples(PILOT_FILES)
    pair_results, _, _ = compute(triples, guarded, ext_weight)
    out = {}
    for label, (a, b) in PILOT_GT.items():
        k = tuple(sorted([a, b]))
        s = pair_results[k][0] if k in pair_results else 0.0
        out[label] = s
    return out


def full_check(all_triples, ext_weight, guarded):
    pair_results, indegree, file_targets = compute(all_triples, guarded, ext_weight)
    out = {}
    for name, (a, b) in [("Cluster A (06<->14, FP)", CLUSTER_A),
                         ("Cluster B (carrier-memetics<->cultural-attraction, TP)", CLUSTER_B),
                         ("Cluster C (sematectonic<->web-scale-stigmergic, TP)", CLUSTER_C)]:
        k = tuple(sorted([a, b]))
        s = pair_results[k][0] if k in pair_results else 0.0
        out[name] = s
    return out, pair_results, indegree


def main():
    files = in_scope_incubator_files()
    all_triples = build_corpus_triples(files)
    print(f"Corpus: {len(files)} files, {len(all_triples)} edge triples (post STEP-1 patch)")

    print("\n" + "=" * 80)
    print("WEIGHT SWEEP — full-corpus cluster scores + pilot GT scores, by external_weight")
    print("(guarded=True; external_weight=1.0 is identical to the unguarded baseline)")
    print("=" * 80)
    for ext_weight in [1.0, 0.6, 0.5, 0.42, 0.25, 0.0]:
        full_scores, _, _ = full_check(all_triples, ext_weight, guarded=True)
        pilot_scores = pilot_check(ext_weight, guarded=True)
        print(f"\n-- external_weight={ext_weight} --")
        for k, v in full_scores.items():
            tag = "PASS" if (("FP" in k and v < THRESHOLD) or ("TP" in k and v >= THRESHOLD)) else "FAIL"
            print(f"    [{tag}] {k}: score={v:.3f}")
        for k, v in pilot_scores.items():
            want_catch = k.startswith("TP")
            tag = "PASS" if (v >= THRESHOLD) == want_catch else "FAIL"
            print(f"    [{tag}] {k}: score={v:.3f}")

    print("\n" + "=" * 80)
    print(f"CHOSEN RULE: guarded, external_weight={EXTERNAL_WEIGHT}")
    print("=" * 80)
    pair_results, indegree, file_targets = compute(all_triples, guarded=True, ext_weight=EXTERNAL_WEIGHT)
    above = [(p, s, shared) for (p, (s, shared)) in pair_results.items() if s >= THRESHOLD]
    above.sort(key=lambda x: -x[1])
    print(f"Pairs with >=1 shared target: {len(pair_results)}   Pairs >= {THRESHOLD}: {len(above)}")
    for (a, b), s, shared in above:
        print(f"  score={s:.3f}  {a}  <->  {b}")
        for t in sorted(shared):
            print(f"      {t} (in-degree={indegree[t]})")
    clusters = connected_components(above)
    print(f"\nConnected components: {len(clusters)}")
    for i, c in enumerate(clusters, 1):
        print(f"  Cluster {i} ({len(c)}): {sorted(c)}")

    print("\n--- BEFORE (unguarded baseline, for comparison) ---")
    pair_results_b, indegree_b, _ = compute(all_triples, guarded=False)
    above_b = [(p, s, shared) for (p, (s, shared)) in pair_results_b.items() if s >= THRESHOLD]
    above_b.sort(key=lambda x: -x[1])
    print(f"Pairs >= {THRESHOLD}: {len(above_b)}")
    for (a, b), s, shared in above_b:
        print(f"  score={s:.3f}  {a}  <->  {b}")
    clusters_b = connected_components(above_b)
    print(f"Connected components: {len(clusters_b)}")
    for i, c in enumerate(clusters_b, 1):
        print(f"  Cluster {i} ({len(c)}): {sorted(c)}")

    print("\n--- Pilot ground-truth re-check (pilot 14-file subgraph, chosen rule) ---")
    pilot_scores = pilot_check(EXTERNAL_WEIGHT, guarded=True)
    for k, v in pilot_scores.items():
        print(f"  {k}: score={v:.3f}  {'>=' if v >= THRESHOLD else '<'} {THRESHOLD}")


if __name__ == "__main__":
    main()
