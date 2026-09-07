#!/usr/bin/env python3
"""Phase 3 move: git mv every mapped experiment dir to its new home.

`git mv` on a DIRECTORY renames the whole subtree on the filesystem (carrying
gitignored local run residue too) and records the rename for tracked files, so
history is preserved, the 1958 adjacent .dvc pointers travel, and per-experiment
`runs/.gitignore` rules keep ignored runs ignored at the new path.

Idempotent: skips a row whose old path is already gone or whose new path exists.
Run from anywhere; operates on the repo root. Use --apply to execute.
"""
from __future__ import annotations

import argparse
import pathlib
import subprocess
import sys

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()

    moved = skipped = errors = 0
    for ln in (HERE / "id-map.tsv").read_text().splitlines():
        if ln.startswith("#") or not ln.strip():
            continue
        old, new = ln.split("\t")[0], ln.split("\t")[1]
        oldp, newp = ROOT / old, ROOT / new
        if not oldp.exists() and newp.exists():
            skipped += 1
            continue
        if not oldp.exists():
            print(f"MISSING old (and no new): {old}")
            errors += 1
            continue
        if newp.exists():
            print(f"NEW already exists, old still present: {new}")
            errors += 1
            continue
        if not args.apply:
            moved += 1
            continue
        newp.parent.mkdir(parents=True, exist_ok=True)
        r = subprocess.run(["git", "-C", str(ROOT), "mv", old, new],
                           capture_output=True, text=True)
        if r.returncode != 0:
            print(f"git mv FAILED: {old} -> {new}\n  {r.stderr.strip()}")
            errors += 1
        else:
            moved += 1
    verb = "moved" if args.apply else "would move"
    print(f"{verb}={moved} skipped(already)={skipped} errors={errors}")
    return 1 if errors else 0


if __name__ == "__main__":
    raise SystemExit(main())
