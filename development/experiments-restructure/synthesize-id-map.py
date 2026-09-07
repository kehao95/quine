#!/usr/bin/env python3
"""Generate id-map.tsv from the Phase 0 inventory + p6 triage.

Deterministic and reproducible: reads ./inventory/inv-*.json (per-phase
classification) and ./inventory/audit-p6triage.json (citation-validated p6
triage), assigns family-prefixed sequential IDs, and emits the canonical
legacy->new path map consumed by repoint.py.

Routing:
  - active experiments  -> experiments/active/<family>/<prefix>-NN-<slug>
  - p6 superseded/design -> experiments/archive/legacy-phase-6/<dotted-id>-<slug>
  - p8 within-phase archive children -> experiments/archive/legacy-p8-early/<id>-<slug>
  - family == REVIEW    -> flagged, no new id (manual resolution)

Existing archive suites (ch4-7, p1-physics, p2-sociology, p3-metaphysics,
p6-bounded-continuity-branch) are intentionally NOT renamed in this pass: they
are stable and correctly cited, so cosmetic renames would add citation-break
risk for no behavioral gain (see Phase 0 decision in restructure-design.md).
"""
from __future__ import annotations

import json
import pathlib
import re

HERE = pathlib.Path(__file__).resolve().parent
INV = HERE / "inventory"

PREFIX = {
    "autopoiesis-succession": "succ",
    "dormancy-ecology": "dorm",
    "trace-inheritance": "trce",
    "null-ecology": "null",
    "cooperation-dynamics": "coop",
    "social-role": "role",
    "digital-red-queen": "drq",
    "structural-elicitation": "elic",
    "cumulative-novelty": "novl",
    "early-existence": "exst",
    "benchmarks": "bench",
}

# Stable phase ordering for NN assignment within a family.
PHASE_ORDER = [
    "p6-autopoiesis", "p6f-future-addressability", "p6m-successor-morphospace",
    "p6r-agentic-runtime-replacement", "p7-oee", "p1-borges", "p3-mrcr",
    "p5-joots", "p4-harbor", "p10-benchmarks", "p9c-compiler-ecology",
    "p7g-generation-ecology", "p7t-trace-ecology", "p7i-irregularity-ecology",
    "p7s-seeded-workspace-ecology", "p7e-open-ecology",
    "p7w-wiki-frontier-ecology", "p7r-role-ecology", "p8-population",
    "p11-digital-redqueen", "p12-dormancy-ecology", "p13-structural-elicitation",
]
PHASE_RANK = {p: i for i, p in enumerate(PHASE_ORDER)}

# Reverse map so an agent that wrote a prefix (e.g. "bench") instead of the
# family name ("benchmarks") still normalizes correctly.
PREFIX_TO_FAMILY = {v: k for k, v in PREFIX.items()}

# Manual resolutions for experiments the per-phase split rules left as REVIEW.
# p1:2.11-saboteur — multi-agent adversarial "saboteur in the commons"; routed
# to cooperation-dynamics by scientific question (adversarial social
# coordination / defensive adaptation), not by p1 emergence history (P2).
REVIEW_OVERRIDES = {
    "p1:2.11": "cooperation-dynamics",
}


def natkey(s: str):
    return [int(t) if t.isdigit() else t for t in re.split(r"(\d+)", s)]


def phase_of(legacy_path: str) -> str:
    # experiments/<phase>/...
    return legacy_path.split("/")[1]


def idpart(legacy_id: str) -> str:
    return legacy_id.split(":", 1)[1] if ":" in legacy_id else legacy_id


def main() -> int:
    # p6 triage dispositions keyed by dotted id (6.7 -> superseded ...)
    triage = {}
    tdata = json.loads((INV / "audit-p6triage.json").read_text())
    for e in tdata["experiments"]:
        triage[idpart(e["legacy_id"])] = e["disposition"]

    active = []        # (family, phase, legacy_id, legacy_path, slug, disposition)
    archived = []      # (legacy_path, new_path, legacy_id, disposition)
    review = []        # (legacy_path, legacy_id, reason)

    for f in sorted(INV.glob("inv-*.json")):
        if f.name == "inv-archive.json":
            continue
        data = json.loads(f.read_text())
        for e in data["experiments"]:
            lp = e["legacy_path"]
            lid = e["legacy_id"]
            slug = e["slug"]
            fam = e["family"]
            ph = phase_of(lp)

            # normalize prefix-as-family, then apply manual REVIEW resolutions
            fam = PREFIX_TO_FAMILY.get(fam, fam)
            if lid in REVIEW_OVERRIDES:
                fam = REVIEW_OVERRIDES[lid]

            leaf = lp.split("/")[-1]  # original on-disk dir name

            # within-phase archive children (p8-population/archive/...) -> legacy-p8-early
            if "/archive/" in lp:
                new = f"experiments/archive/legacy-p8-early/{leaf}"
                archived.append((lp, new, lid, "p8-early-archive"))
                continue

            # p6 triage override
            if ph == "p6-autopoiesis":
                disp = triage.get(idpart(lid), "active")
                if disp in ("superseded", "design-only"):
                    new = f"experiments/archive/legacy-phase-6/{leaf}"
                    archived.append((lp, new, lid, disp))
                    continue
                # canonical -> active

            if fam == "REVIEW" or fam not in PREFIX:
                review.append((lp, lid, e.get("disposition_reason", fam)))
                continue

            active.append((fam, ph, lid, lp, slug))

    # Assign NN per family
    rows = []  # (legacy_path, new_path, legacy_id, new_id, family, disposition)
    by_family: dict[str, list] = {}
    for rec in active:
        by_family.setdefault(rec[0], []).append(rec)

    family_counts = {}
    for fam, recs in by_family.items():
        recs.sort(key=lambda r: (PHASE_RANK.get(r[1], 99), natkey(idpart(r[2]))))
        pfx = PREFIX[fam]
        for i, (_, ph, lid, lp, slug) in enumerate(recs, 1):
            nid = f"{pfx}-{i:02d}"
            np = f"experiments/active/{fam}/{nid}-{slug}"
            rows.append((lp, np, lid, nid, fam, "active"))
        family_counts[fam] = len(recs)

    for lp, np, lid, disp in archived:
        rows.append((lp, np, lid, "(archived)", "-", disp))

    # Emit TSV (sorted by new_path for readability)
    rows.sort(key=lambda r: (r[4], r[1]))
    out = HERE / "id-map.tsv"
    with out.open("w") as fh:
        fh.write("# legacy_path\tnew_path\tlegacy_id\tnew_id\tfamily\tdisposition\n")
        for r in rows:
            fh.write("\t".join(r) + "\n")

    # Report
    print(f"id-map.tsv written: {len(rows)} rows")
    print(f"  active (moved to families): {sum(family_counts.values())}")
    for fam in PREFIX:
        if fam in family_counts:
            print(f"    {PREFIX[fam]:5s} {fam:24s} {family_counts[fam]}")
    print(f"  archived (p6 superseded/design + p8-early): {len(archived)}")
    print(f"  REVIEW (manual): {len(review)}")
    for lp, lid, reason in review:
        print(f"    {lid:28s} {lp}  -- {reason}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
