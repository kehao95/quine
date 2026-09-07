#!/usr/bin/env python3
"""Seed family/reciprocity frontmatter onto active experiment READMEs (Phase 1).

In-place and additive: for each of the 155 active experiments (id-map.tsv
disposition=active), adds the new control-plane keys to its README frontmatter
WITHOUT moving the directory and WITHOUT touching the existing
surface_kind/phase/experiment_id/experiment_type/status keys (so
check-control-plane-frontmatter stays green — experiment_id/phase still match
the legacy path until Phase 3).

Added keys (all top-level — the validator's parser is column-0 line-based, so
nesting is avoided to prevent a nested `phase:` from shadowing the required
one):
  id, legacy_id, family, theory-objects, lineage_phase

Idempotent: skips a README that already carries `id:` in its frontmatter.
"""
from __future__ import annotations

import json
import pathlib
import re

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]
INV = HERE / "inventory"
FM_RE = re.compile(r"\A(---\n)(.*?)(\n---\n)", re.S)


def load_meta():
    meta = {}
    for f in INV.glob("inv-*.json"):
        for e in json.loads(f.read_text())["experiments"]:
            meta[e["legacy_path"]] = {
                "theory_objects": e.get("theory_objects", []),
                "experiment_type": (e.get("experiment_type") or "single-run").strip(),
                "status": (e.get("status") or "unknown").strip(),
            }
    return meta


def main() -> int:
    meta = load_meta()
    seeded_insert = seeded_create = skipped = missing = 0
    for ln in (HERE / "id-map.tsv").read_text().splitlines():
        if ln.startswith("#") or not ln.strip():
            continue
        lp, np, lid, nid, fam, disp = ln.split("\t")
        if disp != "active":
            continue
        # prefer README.md; fall back to DESIGN.md for design-only experiments
        target = ROOT / lp / "README.md"
        if not target.is_file():
            alt = ROOT / lp / "DESIGN.md"
            target = alt if alt.is_file() else target
        if not target.is_file():
            print(f"NO README/DESIGN: {lp}")
            missing += 1
            continue

        m_info = meta.get(lp, {})
        objs = m_info.get("theory_objects", [])
        phase = lp.split("/")[1]
        leaf = lp.split("/")[-1]
        recip = [
            f"id: {nid}",
            f"legacy_id: {lid}",
            f"family: {fam}",
            f"theory-objects: [{', '.join(objs)}]",
            f"lineage_phase: {phase}",
        ]

        text = target.read_text()
        m = FM_RE.match(text)
        if m:
            body = m.group(2)
            if re.search(r"^id:\s", body, re.M):
                skipped += 1
                continue
            new_fm = m.group(1) + body + "\n" + "\n".join(recip) + m.group(3)
            target.write_text(new_fm + text[m.end():])
            seeded_insert += 1
        else:
            # no frontmatter block — create a full one (base + reciprocity keys)
            base = [
                "surface_kind: experiment",
                f"phase: {phase}",
                f"experiment_id: {phase}/{leaf}",
                f"experiment_type: {m_info.get('experiment_type', 'single-run')}",
                f"status: {m_info.get('status', 'unknown')}",
            ]
            block = "---\n" + "\n".join(base + recip) + "\n---\n\n"
            target.write_text(block + text)
            seeded_create += 1
    print(f"inserted={seeded_insert} created-block={seeded_create} "
          f"skipped(already)={skipped} missing={missing}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
