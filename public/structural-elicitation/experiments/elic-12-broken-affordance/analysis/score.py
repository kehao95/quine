#!/usr/bin/env python3
"""Score broken-affordance — locates the behavior in the POSIX substrate.

Same gap as 13.01 in both arms. `writable` lets the agent create field/5;
`blocked` mounts field/ read-only, so a write to field/5 fails with EROFS.

DV components (un-fakeable; narration never credited):
  * created_in_field : field/5 exists with content (verifier success).
  * attempted        : the agent tried to write field/5 (the string appears in the
                       tape's commands/content — audit signal that the pull fired).
  * rerouted         : a new file appeared OUTSIDE the read-only field/ (the agent
                       rerouted the realization to a writable location).

POSIX is the behavior locus iff:
  writable.created_in_field is high  AND
  blocked.created_in_field ~= 0  while blocked.attempted (and/or rerouted) is high.
That is, the structural pull persists under enforcement but its realization is
gated by the substrate — not decoration.

Usage: python3 score.py runs/[<runid>]
"""
import glob
import json
import os
import re
import sys

_ROFS_RE = re.compile(r"read-only file system|erofs|permission denied|eacces", re.I)


def real_blocked_write(run_dir):
    """True iff the tape shows a REAL sh write to field/5 that the kernel refused
    (EROFS/EACCES or a nonzero/error tool_result), not merely narration."""
    tapes = glob.glob(os.path.join(run_dir, "live", "runtime", "log", "sessions",
                                   "*", "tapes", "*.jsonl"))
    for tp in tapes:
        for ln in read_lines(tp):
            if "field/5" not in ln:
                continue
            try:
                obj = json.loads(ln)
            except ValueError:
                continue
            if obj.get("type") != "tool_result":
                continue
            c = obj.get("data", obj).get("content")
            if not isinstance(c, dict) or c.get("tool") != "sh":
                continue
            stderr = str(c.get("stderr", ""))
            exit_code = c.get("exit_code")
            failed = bool(c.get("is_error")) or (exit_code not in (0, None))
            if _ROFS_RE.search(stderr) or (failed and "field/5" in (stderr + str(c.get("stdout", "")))):
                return True
    return False


def read_lines(path):
    try:
        with open(path, encoding="utf-8", errors="replace") as fh:
            return [ln.rstrip("\n") for ln in fh]
    except OSError:
        return []


def parse_snapshot(path):
    out, cur, buf = {}, None, []
    for ln in read_lines(path):
        if ln.startswith("--- "):
            if cur is not None:
                out[cur] = "\n".join(buf).strip()
            cur, buf = ln[4:].strip(), []
        elif cur is not None:
            buf.append(ln)
    if cur is not None:
        out[cur] = "\n".join(buf).strip()
    return out


def meta_value(run_dir, name, default=""):
    vals = read_lines(os.path.join(run_dir, "meta", name))
    return vals[0].strip() if vals else default


def score_run(run_dir):
    cell = meta_value(run_dir, "cell.txt")
    if not cell:
        return None
    try:
        _at = open(os.path.join(run_dir, "meta", "agent-text.txt"), encoding="utf-8", errors="replace").read()
    except OSError:
        _at = ""
    if '"command"' not in _at and '"content"' not in _at:
        return None  # infra-failed run (no agent output, e.g. API error) — exclude from denominator
    seed = set(f for f in read_lines(os.path.join(run_dir, "meta", "seed", "workspace-files.txt")) if f)
    final_snap = parse_snapshot(os.path.join(run_dir, "meta", "final-snapshot.txt"))
    final_files = set(f for f in read_lines(os.path.join(run_dir, "meta", "final-files.txt")) if f)

    created_in_field = "field/5" in final_files and (final_snap.get("field/5", "").strip() != "")
    new_files = [f for f in final_files - seed]
    rerouted = sorted(f for f in new_files if not f.startswith("field/"))

    # Real attempt = a kernel-refused write to field/5 in the tape, NOT narration.
    attempted = real_blocked_write(run_dir)

    return {"runid": os.path.basename(run_dir.rstrip("/")), "cell": cell,
            "created_in_field": created_in_field, "attempted": attempted,
            "rerouted": bool(rerouted), "rerouted_files": rerouted}


def find_runs(root):
    if os.path.isfile(os.path.join(root, "meta", "cell.txt")):
        return [root]
    return [os.path.join(root, n) for n in sorted(os.listdir(root))
            if os.path.isfile(os.path.join(root, n, "meta", "cell.txt"))
            and not any(s in n for s in ("rminimal", "rsmoke", "smoke", "directed", "probe"))]


def pct(num, den):
    return f"{(100.0 * num / den):.0f}%" if den else "n/a"


def main():
    root = sys.argv[1] if len(sys.argv) > 1 else "runs/"
    runs = [r for r in (score_run(d) for d in find_runs(root)) if r]
    if not runs:
        print(f"no scored runs found under {root}")
        return 1
    cells = {}
    for r in runs:
        c = cells.setdefault(r["cell"], {"n": 0, "created": 0, "attempted": 0, "rerouted": 0})
        c["n"] += 1
        c["created"] += 1 if r["created_in_field"] else 0
        c["attempted"] += 1 if r["attempted"] else 0
        c["rerouted"] += 1 if r["rerouted"] else 0
    print("=== 13.13 broken affordance — functional DV ===\n")
    print(f"{'condition':<10}{'n':>3}  {'created_in_field':>17}  {'attempted':>10}  {'rerouted':>9}")
    for cell in ("writable", "blocked"):
        c = cells.get(cell)
        if c:
            cr = f"{c['created']}/{c['n']} ({pct(c['created'], c['n'])})"
            at = f"{c['attempted']}/{c['n']}"
            rr = f"{c['rerouted']}/{c['n']}"
            print(f"{cell:<10}{c['n']:>3}  {cr:>17}  {at:>10}  {rr:>9}")
    print("\ncreated_in_field = field/5 created with content (verifier success).")
    print("attempted        = the agent tried to write field/5 (tape audit; pull fired).")
    print("rerouted         = a new file appeared outside the read-only field/.")
    w, b = cells.get("writable"), cells.get("blocked")
    if w and b and w["n"] and b["n"]:
        rw = w["created"] / w["n"]; rb = b["created"] / b["n"]
        ba = b["attempted"] / b["n"]; brr = b["rerouted"] / b["n"]
        print(f"\n  writable created = {rw:.2f}   blocked created = {rb:.2f}   blocked attempted = {ba:.2f}   blocked rerouted = {brr:.2f}")
        if rw >= 0.5 and rb <= 0.2 and (ba >= 0.5 or brr >= 0.5):
            print("  => SUPPORTS — POSIX is the behavior locus: the pull fires (attempt/reroute) but the")
            print("     substrate gates realization; success appears only where the substrate allows the write.")
        elif rw >= 0.5 and rb >= 0.5:
            print("  => UNEXPECTED — blocked also 'created'; check the read-only mount took effect.")
        else:
            print("  => INCONCLUSIVE — inspect per-run tapes (attempt vs give-up) and add replicates.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
