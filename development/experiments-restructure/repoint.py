#!/usr/bin/env python3
"""Deterministic experiment-path repointer for the experiments/ restructure.

Reads the legacy->new path map (id-map.tsv) and rewrites every occurrence of a
legacy experiment path to its new family-home path across control-plane text
files (markdown, shell, python). Operates by *literal substring replacement* of
the experiment-directory path, which correctly handles every reference form:

    experiments/p6-autopoiesis/6.1-minimal-self-birth
    experiments/p6-autopoiesis/6.1-minimal-self-birth/
    experiments/p6-autopoiesis/6.1-minimal-self-birth/README.md
    ../experiments/p6-autopoiesis/6.1-minimal-self-birth/README.md
    ../../../experiments/p6-autopoiesis/6.1-minimal-self-birth/

Relative `../` link prefixes are preserved automatically: only the
experiments-rooted portion deepens, and that portion is expressed in full by the
link, so the resolved target stays correct. (Runner scripts that climb OUT of
their experiment with `${SCRIPT_DIR}/../../../` are a SEPARATE concern handled by
fix-runner-depth.py — this tool does not touch SCRIPT_DIR climbs.)

What is intentionally NOT rewritten:
  - run/tape/batch/campaign payload trees (retained evidence; Decision 7 forbids
    rewriting paths embedded in immutable run metadata)
  - generated heartbeat readout snapshots (development/heartbeat/readouts/)
  - .dvc pointer files (leaf-only `path:` values; unaffected by directory moves)

Default mode is DRY-RUN (report only). Pass --apply to write changes.

Usage:
  python3 repoint.py --map id-map.tsv                 # dry-run, whole repo
  python3 repoint.py --map id-map.tsv --apply         # apply
  python3 repoint.py --map id-map.tsv --only Paper     # limit to a subtree
"""
from __future__ import annotations

import argparse
import pathlib
import sys

# Directory path-components that mark a tree we must never rewrite.
EXCLUDE_DIR_PARTS = {
    "runs", "batches", "campaigns", "tapes", "workspace",
    ".git", "node_modules", "__pycache__", ".cache-bin",
}
# Specific subtree prefixes (repo-relative) to exclude wholesale.
EXCLUDE_PREFIXES = (
    "development/heartbeat/readouts/",
    "development/experiments-restructure/",  # the map + this tool live here
)
# File extensions we treat as repointable control-plane text. Intentionally
# excludes data/generated formats (.csv/.json/.yaml manifests, analysis records)
# whose embedded paths are historical values, not control-plane links — those
# are repointed by regeneration or left as retained evidence (cf. Decision 7).
TEXT_EXTS = {".md", ".sh", ".py"}
# Root-level files (no dir) that are control-plane and should be repointed.
ROOT_FILES = {"AGENTS.md", "CLOSURE.md", "README.md", "ROADMAP.md", "STATUS.md",
              "EVALUATION.md", "DEVELOPMENT.md", "QUICKSTART.md"}


def load_map(map_path: pathlib.Path) -> list[tuple[str, str]]:
    pairs: list[tuple[str, str]] = []
    for ln, raw in enumerate(map_path.read_text().splitlines(), 1):
        line = raw.rstrip("\n")
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        cols = line.split("\t")
        if len(cols) < 2:
            continue
        legacy, new = cols[0].strip(), cols[1].strip()
        if not legacy or not new or new in {"-", "(archived)", "(none)"}:
            continue
        if legacy == new:
            continue
        pairs.append((legacy, new))
    # Longest legacy path first so no shorter path is a prefix of a longer one
    # during sequential replacement.
    pairs.sort(key=lambda p: len(p[0]), reverse=True)
    return pairs


def is_excluded(rel: str) -> bool:
    for pre in EXCLUDE_PREFIXES:
        if rel.startswith(pre):
            return True
    parts = set(pathlib.PurePosixPath(rel).parts)
    if parts & EXCLUDE_DIR_PARTS:
        return True
    for part in pathlib.PurePosixPath(rel).parts:
        if part.startswith("runs-"):
            return True
    return False


def candidate_files(root: pathlib.Path, only: str | None):
    base = root / only if only else root
    for p in base.rglob("*"):
        if not p.is_file():
            continue
        rel = p.relative_to(root).as_posix()
        if is_excluded(rel):
            continue
        if p.suffix in TEXT_EXTS or (p.parent == root and p.name in ROOT_FILES):
            yield p, rel


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--map", required=True)
    ap.add_argument("--root", default=None, help="repo root (default: 3 levels up)")
    ap.add_argument("--only", default=None, help="restrict to a repo subtree, e.g. Paper")
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()

    here = pathlib.Path(__file__).resolve()
    root = pathlib.Path(args.root).resolve() if args.root else here.parents[2]
    map_path = pathlib.Path(args.map)
    if not map_path.is_absolute():
        map_path = (here.parent / map_path).resolve()

    pairs = load_map(map_path)
    if not pairs:
        print(f"no replacement pairs loaded from {map_path}", file=sys.stderr)
        return 2

    total_repls = 0
    files_changed = 0
    per_file: list[tuple[int, str]] = []

    for path, rel in candidate_files(root, args.only):
        try:
            text = path.read_text()
        except (UnicodeDecodeError, OSError):
            continue
        n = 0
        new_text = text
        for legacy, new in pairs:
            if legacy in new_text:
                c = new_text.count(legacy)
                if c:
                    new_text = new_text.replace(legacy, new)
                    n += c
        if n:
            per_file.append((n, rel))
            total_repls += n
            files_changed += 1
            if args.apply:
                path.write_text(new_text)

    per_file.sort(reverse=True)
    mode = "APPLIED" if args.apply else "DRY-RUN"
    print(f"[{mode}] repoint over {root}{'/' + args.only if args.only else ''}")
    print(f"  replacement pairs: {len(pairs)}")
    print(f"  files changed:     {files_changed}")
    print(f"  total occurrences: {total_repls}")
    for n, rel in per_file[:40]:
        print(f"    {n:6d}  {rel}")
    if len(per_file) > 40:
        print(f"    ... and {len(per_file) - 40} more files")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
