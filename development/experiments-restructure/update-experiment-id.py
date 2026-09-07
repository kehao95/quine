#!/usr/bin/env python3
"""Phase 3: update moved experiments' `experiment_id` frontmatter to the new path.

After the move, each experiment lives at experiments/active/<family>/<id-slug>/;
its `experiment_id` must read `active/<family>/<id-slug>` so
check-control-plane-frontmatter (active/ branch) matches. `phase`/`lineage_phase`
are left untouched as lineage. Idempotent.
"""
from __future__ import annotations

import pathlib
import re

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]


def main() -> int:
    updated = skipped = missing = 0
    for ln in (HERE / "id-map.tsv").read_text().splitlines():
        if ln.startswith("#") or not ln.strip():
            continue
        cols = ln.split("\t")
        new, disp = cols[1], cols[5]
        if disp != "active":
            continue
        eid = new[len("experiments/"):]  # active/<family>/<id-slug>
        target = ROOT / new / "README.md"
        if not target.is_file():
            alt = ROOT / new / "DESIGN.md"
            target = alt if alt.is_file() else target
        if not target.is_file():
            print(f"NO README/DESIGN: {new}")
            missing += 1
            continue
        text = target.read_text()
        new_text, n = re.subn(r"(?m)^experiment_id:.*$", f"experiment_id: {eid}", text, count=1)
        if n and new_text != text:
            target.write_text(new_text)
            updated += 1
        else:
            skipped += 1
    print(f"updated={updated} skipped(no-change)={skipped} missing={missing}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
