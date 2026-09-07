#!/usr/bin/env python3
"""Usage: python3 regenerate-works-registry.py — rebuilds Paper/theory/references/works/registry.tsv from the live `external-occupancy` edges in Paper/incubator/**/*.md, merging (not overwriting) any existing non-default coverage_state/title.

Field-baseline layer, Phase 1 (see Paper/_design/proposals/field-baseline-and-novelty.md).
This is the additive node-set regenerator: it re-derives the id-set and
citing_count from the current edges every run, but preserves any coverage_state
or title already recorded in the registry for an id that survives (the false-
frontier guard depends on human/Level-4 coverage marks surviving regeneration).
"""

import csv
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]
INCUBATOR_DIR = REPO_ROOT / "Paper" / "incubator"
REGISTRY_PATH = REPO_ROOT / "Paper" / "theory" / "references" / "works" / "registry.tsv"

EDGE_RE = re.compile(r'rel:\s*external-occupancy\s*,\s*to:\s*"([^"]+)"')

COLUMNS = ["external_id", "citing_count", "coverage_state", "title"]
DEFAULT_COVERAGE_STATE = "unsearched"


def find_edges():
    """Return {external_id: set(citing file paths)} across the current incubator corpus."""
    occurrences = {}
    for path in sorted(INCUBATOR_DIR.rglob("*.md")):
        text = path.read_text(encoding="utf-8")
        ids_in_file = set(EDGE_RE.findall(text))
        for ext_id in ids_in_file:
            occurrences.setdefault(ext_id, set()).add(str(path.relative_to(REPO_ROOT)))
    return occurrences


def load_existing_registry():
    """Return {external_id: {"coverage_state": ..., "title": ...}} from the current
    registry.tsv, if present. Used to preserve human/Level-4 overlay on regeneration."""
    existing = {}
    if not REGISTRY_PATH.exists():
        return existing
    with REGISTRY_PATH.open(newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f, delimiter="\t")
        for row in reader:
            existing[row["external_id"]] = {
                "coverage_state": row.get("coverage_state", DEFAULT_COVERAGE_STATE),
                "title": row.get("title", ""),
            }
    return existing


def build_rows(occurrences, existing):
    rows = []
    for ext_id, files in occurrences.items():
        prior = existing.get(ext_id, {})
        coverage_state = prior.get("coverage_state") or DEFAULT_COVERAGE_STATE
        title = prior.get("title") or ""
        rows.append(
            {
                "external_id": ext_id,
                "citing_count": len(files),
                "coverage_state": coverage_state,
                "title": title,
            }
        )
    # Sort by citing_count desc, then id asc.
    rows.sort(key=lambda r: (-r["citing_count"], r["external_id"]))
    return rows


def write_registry(rows):
    REGISTRY_PATH.parent.mkdir(parents=True, exist_ok=True)
    with REGISTRY_PATH.open("w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, delimiter="\t", fieldnames=COLUMNS)
        writer.writeheader()
        for row in rows:
            writer.writerow(row)


def main():
    occurrences = find_edges()
    existing = load_existing_registry()
    rows = build_rows(occurrences, existing)
    write_registry(rows)

    dropped = sorted(set(existing) - set(occurrences))
    print(f"regenerated {REGISTRY_PATH.relative_to(REPO_ROOT)}: {len(rows)} rows "
          f"from {sum(len(f) for f in occurrences.values())} edges across "
          f"{len(set().union(*occurrences.values())) if occurrences else 0} citing files.")
    if dropped:
        print(f"WARNING: {len(dropped)} previously-registered id(s) no longer have a "
              f"live edge (kept out of the rebuilt registry, non-additive by design "
              f"if you re-run this): {', '.join(dropped)}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
