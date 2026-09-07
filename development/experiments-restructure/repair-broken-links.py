#!/usr/bin/env python3
"""Idempotent repair of broken relative markdown links under experiments/.

fix-readme-links.py was not idempotent; running it more than once over-shifted
some relative links (each extra pass re-resolved an already-correct link against
its old location). Every such over-shift produces a link that does NOT resolve
(verified), so this pass is safe: it only touches links that fail to resolve
from their CURRENT location, and reconstructs the intended target from the
id-map. Correct links resolve and are left untouched -> idempotent.

Target reconstruction for a broken link (after stripping leading `../`):
  1. repo-root file/dir exists (e.g. DEVELOPMENT.md, internal/, Paper/...)
  2. experiments/<tail> exists (legacy phase dirs/READMEs still present)
  3. tail is a legacy experiment path in the move map -> its new path
  4. a path component is a NEW id/dir leaf (half-remapped link) -> its new path
"""
from __future__ import annotations

import argparse
import os
import pathlib
import re

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]
LINK_RE = re.compile(r"\]\(([^)]+)\)")
# mirror check-active-doc-links excluded_segments: links into / docs under these
# are not link-checked (historical/generated/absent payloads), so leave them.
EXCLUDED = {"runs", "archive", "artifacts", "jobs", "workspaces", "sections",
            "analysis", "external-sources", "templates", "tapes", "evidence",
            "batches", "campaigns"}


def load_maps():
    old2new = {}
    leaf2path = {}
    for fn in ("id-map.tsv", "supplemental-map.tsv", "phase-dir-map.tsv"):
        for ln in (HERE / fn).read_text().splitlines():
            if ln.startswith("#") or not ln.strip():
                continue
            c = ln.split("\t")
            old, new = c[0], c[1]
            old2new[old] = new
            if fn == "id-map.tsv":
                leaf2path[new.split("/")[-1]] = new
    return old2new, leaf2path


def strip_dotdot(p: str):
    n = 0
    while p.startswith("../"):
        p = p[3:]
        n += 1
    return p, n


def reconstruct(tail: str, old2new, leaf2path):
    if (ROOT / tail).exists():
        return tail
    if (ROOT / "experiments" / tail).exists():
        return f"experiments/{tail}"
    # strip a trailing /README.md (or /) to test the dir, remember the suffix
    base, suffix = tail, ""
    for suf in ("/README.md", "/DESIGN.md", "/"):
        if tail.endswith(suf):
            base, suffix = tail[: -len(suf)], suf
            break
    if base in old2new:
        return old2new[base] + suffix
    # a relative link climbing to experiments/ omits the prefix in its tail
    if f"experiments/{base}" in old2new:
        return old2new[f"experiments/{base}"] + suffix
    comps = base.split("/")
    for i, comp in enumerate(comps):
        if comp in leaf2path:
            rest = "/".join(comps[i + 1:])
            target = leaf2path[comp] + ("/" + rest if rest else "")
            return target + suffix
    return None


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()
    old2new, leaf2path = load_maps()

    files_changed = fixed = unfixable = 0
    unfixable_list = []
    for md in sorted((ROOT / "experiments").rglob("*.md")):
        rel = md.relative_to(ROOT).as_posix()
        parts = rel.split("/")
        if parts[1] in {"incubator", "protocol"} or EXCLUDED & set(parts):
            continue  # only repair docs that are themselves link-checked
        new_dir = "/".join(parts[:-1])
        text = md.read_text()
        n_local = [0]

        def repl(m):
            link = m.group(1).strip()
            if (not link or link[0] in "#/" or "://" in link
                    or link.startswith(("http", "mailto:"))):
                return m.group(0)
            path_part, sep, anchor = link.partition("#")
            if not path_part or not path_part.startswith("../"):
                return m.group(0)
            if (ROOT / os.path.normpath(f"{new_dir}/{path_part}")).exists():
                return m.group(0)  # already correct
            tail, _ = strip_dotdot(path_part)
            target = reconstruct(tail, old2new, leaf2path)
            if not target or not (ROOT / os.path.normpath(target.rstrip("/"))).exists():
                unfixable_list.append(f"{rel} :: {link}")
                return m.group(0)
            new_rel = os.path.relpath(target.rstrip("/"), new_dir)
            if target.endswith("/"):
                new_rel += "/"
            n_local[0] += 1
            return "](" + new_rel + (sep + anchor if sep else "") + ")"

        new_text = LINK_RE.sub(repl, text)
        if n_local[0]:
            files_changed += 1
            fixed += n_local[0]
            if args.apply:
                md.write_text(new_text)
    unfixable = len(unfixable_list)
    mode = "APPLIED" if args.apply else "DRY-RUN"
    print(f"[{mode}] repair-broken-links: files={files_changed} fixed={fixed} unfixable={unfixable}")
    for u in unfixable_list[:40]:
        print(f"    UNFIXABLE {u}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
