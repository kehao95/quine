#!/usr/bin/env python3
"""Phase 3: repoint RELATIVE markdown links inside moved experiment docs.

repoint.py only rewrote `experiments/<phase>/<exp>` absolute-from-root paths.
Experiment READMEs also carry relative links — to sibling experiments
(`../6.59-...`), across phases (`../../p6m-.../6M.03`), and up to repo root
(`../../../Paper/...`). After the depth-3 -> depth-4 move these all shift.

General fix: for each markdown link in a moved doc, resolve the target against
the doc's OLD location, remap any path that fell under a moved directory
(id-map + supplemental-map), then re-express the target relative to the doc's
NEW location. This is exact for depth changes, sibling renames, and cross-phase
moves simultaneously.

Scope: docs under experiments/active/<family>/<id>/ whose <id> dir is a known
id-map experiment move (so the OLD location is unambiguous). Family-level
absorbed docs (active/<family>/{lib,analysis,notes,...}) are out of scope here.

DRY-RUN by default; --apply to write.
"""
from __future__ import annotations

import argparse
import os
import pathlib
import re

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]
LINK_RE = re.compile(r"\]\(([^)]+)\)")


def load_dir_map():
    old2new = {}
    new2old_exp = {}
    # merged family subsurfaces have an ambiguous single old location, so they
    # are not used to resolve a doc's OWN old location (their links are handled
    # separately if needed); they still participate in target remapping.
    merged_skip = {
        "experiments/active/autopoiesis-succession/analysis",
        "experiments/active/autopoiesis-succession/lib",
    }
    for fn in ("id-map.tsv", "supplemental-map.tsv"):
        for ln in (HERE / fn).read_text().splitlines():
            if ln.startswith("#") or not ln.strip():
                continue
            c = ln.split("\t")
            old, new = c[0], c[1]
            old2new[old] = new
            if fn == "id-map.tsv":
                new2old_exp[new] = old
            elif not new.endswith(".md") and new not in merged_skip:
                new2old_exp[new] = old  # 1:1 family-subsurface dir move
    pairs = sorted(old2new.items(), key=lambda kv: len(kv[0]), reverse=True)
    return pairs, new2old_exp


def remap(target: str, pairs) -> str:
    for o, n in pairs:
        if target == o or target.startswith(o + "/"):
            return n + target[len(o):]
    return target


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()
    pairs, new2old_exp = load_dir_map()

    files_changed = links_changed = 0
    skip_parts = {"runs", "tapes", "batches", "campaigns", "incubator"}
    for md in sorted((ROOT / "experiments").rglob("*.md")):
        rel = md.relative_to(ROOT).as_posix()
        parts = rel.split("/")
        if parts[1] in {"incubator", "protocol"}:
            continue
        if skip_parts & set(parts) or any(p.startswith("runs-") for p in parts):
            continue
        new_file_dir = "/".join(parts[:-1])
        # Find the doc's OLD location: longest new-dir prefix that is a known
        # move; otherwise the doc did not move (phase README, router) -> old=current.
        old_file_dir = new_file_dir
        for depth in range(min(len(parts) - 1, 6), 1, -1):
            cand = "/".join(parts[:depth])
            if cand in new2old_exp:
                sub = "/".join(parts[depth:-1])
                base = new2old_exp[cand]
                old_file_dir = base if not sub else f"{base}/{sub}"
                break

        text = md.read_text()
        n_local = [0]

        def repl(m):
            link = m.group(1).strip()
            if (not link or link[0] in "#/" or link.startswith(("http://", "https://", "mailto:"))
                    or "://" in link):
                return m.group(0)
            path_part, sep, anchor = link.partition("#")
            if not path_part:
                return m.group(0)
            old_target = os.path.normpath(f"{old_file_dir}/{path_part}")
            new_target = remap(old_target, pairs)
            new_rel = os.path.relpath(new_target, new_file_dir)
            newlink = new_rel + (sep + anchor if sep else "")
            if newlink != link:
                n_local[0] += 1
                return "](" + newlink + ")"
            return m.group(0)

        new_text = LINK_RE.sub(repl, text)
        if n_local[0]:
            files_changed += 1
            links_changed += n_local[0]
            if args.apply:
                md.write_text(new_text)

    mode = "APPLIED" if args.apply else "DRY-RUN"
    print(f"[{mode}] fix-readme-links: files_changed={files_changed} links_changed={links_changed}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
