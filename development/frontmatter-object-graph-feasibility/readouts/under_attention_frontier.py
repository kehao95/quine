#!/usr/bin/env python3
"""
GENERATED VIEW SUPPORT SCRIPT — regenerate readouts/under-attention-frontier-*.md
by running this file; not a canonical store, per read-layer-prior.

Tier-1 under-attention frontier: every Paper/theory/objects/*.md and
Paper/theory/gaps/*.md node with idea-layer in-degree 0 or 1, computed over
the STEP-1-patched full 101-file Paper/incubator/ corpus. In-degree here is
the RAW count of distinct files pointing at a target via any non-home-family
edge (external-occupancy vs internal is not relevant to this readout — it
measures whether the idea layer has looked at an existing theory node at
all, not same-referent risk).
"""
import glob
import os
from collections import defaultdict

from graph_lib import in_scope_incubator_files, build_corpus_triples, REPO


def main():
    files = in_scope_incubator_files()
    all_triples = build_corpus_triples(files)

    target_sources = defaultdict(set)
    for (src, rel, to) in all_triples:
        if rel == "home-family":
            continue
        target_sources[to].add(src)
    indegree = {t: len(s) for t, s in target_sources.items()}

    objects = sorted(glob.glob(os.path.join(REPO, "Paper/theory/objects/*.md")))
    gaps = sorted(glob.glob(os.path.join(REPO, "Paper/theory/gaps/*.md")))

    def name_key(p, kind):
        base = os.path.basename(p)[:-3]
        return f"theory/{kind}/{base}"

    zero, one, two_plus = [], [], []
    for p in objects:
        if os.path.basename(p) == "README.md":
            continue
        k = name_key(p, "objects")
        d = indegree.get(k, 0)
        (zero if d == 0 else one if d == 1 else two_plus).append((k, d))
    for p in gaps:
        if os.path.basename(p) == "README.md":
            continue
        k = name_key(p, "gaps")
        d = indegree.get(k, 0)
        (zero if d == 0 else one if d == 1 else two_plus).append((k, d))

    n_objects = len([p for p in objects if os.path.basename(p) != "README.md"])
    n_gaps = len([p for p in gaps if os.path.basename(p) != "README.md"])
    total_nodes = n_objects + n_gaps

    print(f"theory/objects: {n_objects}  theory/gaps: {n_gaps}  total: {total_nodes}")
    print(f"in-degree 0: {len(zero)}   in-degree 1: {len(one)}   in-degree 2+: {len(two_plus)}")
    print(f"Tier-1 frontier (0 or 1): {len(zero) + len(one)} / {total_nodes}")

    zero.sort()
    one.sort()

    print("\n-- ZERO --")
    for k, d in zero:
        print(f"  {k}")
    print("\n-- ONE --")
    for k, d in one:
        print(f"  {k}")

    return zero, one, two_plus, n_objects, n_gaps


if __name__ == "__main__":
    main()
