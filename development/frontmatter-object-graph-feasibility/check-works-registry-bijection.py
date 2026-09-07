#!/usr/bin/env python3
"""Usage: python3 check-works-registry-bijection.py — verifies the field-baseline
bijection: every `external-occupancy` edge id in Paper/incubator/**/*.md appears
as a row in Paper/theory/references/works/registry.tsv, and every registry row
is cited by >=1 live edge. Reports mismatches; exits non-zero if any are found.
"""

import csv
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]
INCUBATOR_DIR = REPO_ROOT / "Paper" / "incubator"
REGISTRY_PATH = REPO_ROOT / "Paper" / "theory" / "references" / "works" / "registry.tsv"

EDGE_RE = re.compile(r'rel:\s*external-occupancy\s*,\s*to:\s*"([^"]+)"')


def find_edge_ids():
    ids = set()
    for path in INCUBATOR_DIR.rglob("*.md"):
        text = path.read_text(encoding="utf-8")
        ids.update(EDGE_RE.findall(text))
    return ids


def load_registry_ids():
    if not REGISTRY_PATH.exists():
        print(f"ERROR: registry not found at {REGISTRY_PATH}", file=sys.stderr)
        sys.exit(2)
    with REGISTRY_PATH.open(newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f, delimiter="\t")
        return {row["external_id"] for row in reader}


def main():
    edge_ids = find_edge_ids()
    registry_ids = load_registry_ids()

    unregistered = sorted(edge_ids - registry_ids)   # edges with no node
    uncited = sorted(registry_ids - edge_ids)         # nodes with no live edge

    print(f"edge ids (distinct): {len(edge_ids)}")
    print(f"registry rows: {len(registry_ids)}")

    if not unregistered and not uncited:
        print("BIJECTION HOLDS: every edge id has exactly one registry row, "
              "and every registry row is cited by >=1 edge.")
        return 0

    if unregistered:
        print(f"\nMISMATCH: {len(unregistered)} edge id(s) with no registry row "
              f"(run the regenerator):")
        for i in unregistered:
            print(f"  - {i}")
    if uncited:
        print(f"\nMISMATCH: {len(uncited)} registry row(s) with no live edge "
              f"(stale; re-derive with the regenerator):")
        for i in uncited:
            print(f"  - {i}")
    return 1


if __name__ == "__main__":
    sys.exit(main())
