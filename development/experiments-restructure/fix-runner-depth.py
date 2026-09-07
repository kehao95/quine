#!/usr/bin/env python3
"""Fix `${SCRIPT_DIR}` relative-depth climbs in runner scripts after the move.

The migration inserts exactly one new level (`active/<family>` in place of
`<phase>`, plus the new `active/` parent) between `experiments/` and the
experiment dir:

    legacy:  experiments/<phase>/<exp>/...                 (exp at repo-depth 3)
    new:     experiments/active/<family>/<exp>/...         (exp at repo-depth 4)

So a `${SCRIPT_DIR}` climb only needs `+1 ../` when it reaches `experiments/`
or repo root — i.e. when it climbs ABOVE the experiment's grouping dir. A climb
that lands AT the grouping dir (phase→family) keeps its depth, but if it then
names a sibling experiment that sibling has been renamed/re-homed, so those are
FLAGGED for manual Appendix-G reconciliation rather than rewritten.

Depth-aware rule, for a script whose dir sits `d` levels below its experiment
dir (d=0 for `<exp>/run.sh`, d=1 for `<exp>/setup/x.sh`):

    climb count N  >= d+2   ->  root/experiments-crossing  -> rewrite to N+1
    climb count N  == d+1   ->  grouping-level (sibling)    -> FLAG
    climb count N  <= d     ->  inside the experiment       -> leave

Default mode is DRY-RUN. Pass --apply to write changes.

Usage (pre-move scope scan over legacy phases):
  python3 fix-runner-depth.py --only experiments
Usage (post-move fix over the new tree):
  python3 fix-runner-depth.py --only experiments/active --apply
"""
from __future__ import annotations

import argparse
import pathlib
import re

CLIMB = re.compile(r"\$\{?SCRIPT_DIR\}?((?:/\.\.)+)")


def exp_depth(rel: str) -> int | None:
    """Levels the script's dir sits below its experiment dir, or None if the
    path is not under a recognizable experiment layout."""
    parts = rel.split("/")
    if len(parts) < 3 or parts[0] != "experiments":
        return None
    if parts[1] == "active":
        exp_idx = 3  # experiments/active/<family>/<exp>/...
    elif parts[1] == "archive":
        return None  # archive layout handled separately
    else:
        exp_idx = 2  # experiments/<phase>/<exp>/...
    if len(parts) <= exp_idx + 1:
        return None
    # dirs strictly between the exp dir and the script file
    return len(parts[exp_idx + 1:-1])


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--root", default=None)
    ap.add_argument("--only", default="experiments")
    ap.add_argument("--apply", action="store_true")
    args = ap.parse_args()

    here = pathlib.Path(__file__).resolve()
    root = pathlib.Path(args.root).resolve() if args.root else here.parents[2]
    base = root / args.only

    files_rewritten = 0
    climbs_rewritten = 0
    flagged: list[str] = []
    skipped_layout = 0

    for p in sorted(base.rglob("*.sh")):
        rel = p.relative_to(root).as_posix()
        if "/runs/" in rel or "/runs-" in rel or "/tapes/" in rel:
            continue
        d = exp_depth(rel)
        if d is None:
            skipped_layout += 1
            continue
        try:
            text = p.read_text()
        except (UnicodeDecodeError, OSError):
            continue

        file_hits = 0

        def repl(m: re.Match) -> str:
            nonlocal file_hits
            run = m.group(1)
            n = run.count("/..")
            if n >= d + 2:               # crosses experiments/ or root -> +1
                file_hits += 1
                return m.group(0).replace(run, run + "/..")
            return m.group(0)            # grouping/inside -> leave (maybe flag)

        new_text = CLIMB.sub(repl, text)

        # Flag grouping-level sibling climbs (N == d+1) for manual remap.
        for m in CLIMB.finditer(text):
            n = m.group(1).count("/..")
            if n == d + 1:
                tail = text[m.end():m.end() + 40].split()[0] if text[m.end():m.end() + 1] not in ("", " ", '"', "'") else ""
                flagged.append(f"{rel}: climb={n} (grouping-level sibling) ...{m.group(0)}{tail}")

        if file_hits:
            files_rewritten += 1
            climbs_rewritten += file_hits
            if args.apply:
                p.write_text(new_text)

    mode = "APPLIED" if args.apply else "DRY-RUN"
    print(f"[{mode}] fix-runner-depth over {base}")
    print(f"  scripts with +1 climbs: {files_rewritten}")
    print(f"  climbs rewritten (+1 ../): {climbs_rewritten}")
    print(f"  grouping-level sibling climbs flagged for Appendix-G remap: {len(flagged)}")
    print(f"  scripts skipped (non-experiment layout / archive): {skipped_layout}")
    for f in flagged[:50]:
        print(f"    {f}")
    if len(flagged) > 50:
        print(f"    ... and {len(flagged) - 50} more")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
