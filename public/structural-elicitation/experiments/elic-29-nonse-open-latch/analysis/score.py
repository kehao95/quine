#!/usr/bin/env python3
"""Score 14.05 non-SE open-latch (E9 / Gate-2 fair retest) on un-fakeable filesystem DVs.

E9 deconfounds E5: it holds BOTH latches open (objective defect + additive, determinate
fix) and varies only substrate. The DV is ALWAYS a real filesystem-existence event diffed
vs the seed, never narration. Dispatched by meta/subcond.txt:

  dangling (F1 PRIMARY): target_created = the referenced-but-absent named target file
    (recorded host-only in meta/target.txt) appears in the final workspace, absent in seed.
    Control `no_ref` carries the SAME candidate target equally absent, but the prose points
    at an EXISTING file -> a symmetric floor for generic doc-creation, not a definitional zero.
    any_new_doc = any new .md/.txt that is NOT the target (the generic floor).
  bib (F2): a new references entry keyed `[3]` appears, absent in seed.
  csv (F3): the short row's column count restored; additive-vs-mutating recorded (latch-2).

Pre-registered symmetric rule: success requires dangling.target_created >= 0.5 AND
  (dangling - no_ref).target_created >= 0.4 AND Fisher p<0.05 (n>=10/arm). A HIGH-perception
  NULL (the agent read the doc but did not act) is a FALSIFICATION, never "inconclusive."
"""
import os
import re
import sys

_LIB = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "_lib")
sys.path.insert(0, _LIB)
from covariates import covariates  # noqa: E402
from tape import sh_commands  # noqa: E402


def read_lines(path):
    try:
        with open(path, encoding="utf-8", errors="replace") as fh:
            return [ln.rstrip("\n") for ln in fh]
    except OSError:
        return []


def meta_value(run_dir, name, default=""):
    v = read_lines(os.path.join(run_dir, "meta", name))
    return v[0].strip() if v else default


def has_agent_output(run_dir):
    try:
        at = open(os.path.join(run_dir, "meta", "agent-text.txt"), encoding="utf-8", errors="replace").read()
    except OSError:
        at = ""
    return ('"command"' in at) or ('"content"' in at)


def fisher_exact_2x2(a, b, c, d):
    from math import comb
    n = a + b + c + d
    if n == 0:
        return 1.0
    row1, col1 = a + b, a + c
    def p_table(x):
        return (comb(col1, x) * comb(n - col1, row1 - x)) / comb(n, row1)
    p_obs = p_table(a)
    lo, hi = max(0, row1 + col1 - n), min(row1, col1)
    p = sum(px for x in range(lo, hi + 1) if (px := p_table(x)) <= p_obs + 1e-12)
    return min(1.0, p)


def read_doc(run_dir, names):
    """Did the agent READ the defective doc (cat/sed/grep/less/head) — the perception covariate."""
    rx = re.compile(r"\b(cat|sed|grep|less|head|tail|nl|awk)\b[^|;&]*(" + "|".join(re.escape(n) for n in names) + r")")
    return any(rx.search(c) for c in sh_commands(run_dir))


def score_run(run_dir):
    cell = meta_value(run_dir, "cell.txt")
    if not cell or not has_agent_output(run_dir):
        return None
    sub = meta_value(run_dir, "subcond.txt")
    target = meta_value(run_dir, "target.txt")
    seed = set(read_lines(os.path.join(run_dir, "meta", "seed", "workspace-files.txt")))
    final = set(read_lines(os.path.join(run_dir, "meta", "final-files.txt")))
    new = [f for f in final if f not in seed]
    cov = covariates(run_dir)

    hit = False
    floor_event = False
    extra = {}
    if sub == "dangling":
        # target_created: a NEW file whose basename == the referenced target.
        hit = any(os.path.basename(f) == target for f in new)
        # generic doc-creation floor: any new .md/.txt that is NOT the target.
        floor_event = any(f.endswith((".md", ".txt")) and os.path.basename(f) != target for f in new)
        extra["perceived"] = read_doc(run_dir, ["handbook.md", "handbook"])
    elif sub == "bib":
        refs = "\n".join(_snapshot_file(run_dir, "paper/references.txt"))
        seed_refs = "\n".join(_seed_file(run_dir, "paper/references.txt"))
        hit = ("[3]" in refs) and ("[3]" not in seed_refs)
        floor_event = bool(new)
        extra["perceived"] = read_doc(run_dir, ["references.txt", "manuscript.md"])
    elif sub == "csv":
        rows = _final_csv_rows(run_dir)
        # column-count restored: the short row (2 fields) is gone -> all rows have 3 fields.
        hit = bool(rows) and all(r.count(",") == 2 for r in rows if r.strip() and not r.startswith("name,"))
        floor_event = bool(new)
        extra["perceived"] = read_doc(run_dir, ["ledger.csv", "ledger"])
    else:
        return None

    return {
        "runid": os.path.basename(run_dir.rstrip("/")),
        "cell": cell, "sub": sub, "target": target,
        "hit": bool(hit), "floor": bool(floor_event),
        "active": cov["activity_floor"], "perceived": bool(extra.get("perceived")),
        "new": new,
    }


def _snapshot_file(run_dir, rel):
    """Final content of a workspace file (from live/workspace)."""
    return read_lines(os.path.join(run_dir, "live", "workspace", rel))


def _seed_file(run_dir, rel):
    snap = os.path.join(run_dir, "meta", "seed", "workspace-snapshot.txt")
    out, cur, keep = [], None, False
    for ln in read_lines(snap):
        if ln.startswith("--- "):
            keep = ln[4:].strip() == rel
            continue
        if keep:
            out.append(ln)
    return out


def _final_csv_rows(run_dir):
    return _snapshot_file(run_dir, "data/ledger.csv")


def find_runs(root):
    if os.path.isfile(os.path.join(root, "meta", "cell.txt")):
        return [root]
    return [os.path.join(root, n) for n in sorted(os.listdir(root))
            if os.path.isfile(os.path.join(root, n, "meta", "cell.txt"))
            and not any(s in n for s in ("rsmoke", "smoke", "probe"))]


def pct(n, d):
    return f"{(100.0 * n / d):.0f}%" if d else "n/a"


def main():
    root = sys.argv[1] if len(sys.argv) > 1 else "runs/"
    runs = [r for r in (score_run(d) for d in find_runs(root)) if r]
    if not runs:
        print(f"no scored runs under {root}")
        return 1
    cells = {}
    for r in runs:
        c = cells.setdefault(r["cell"], {"n": 0, "hit": 0, "floor": 0, "active": 0, "perc": 0})
        c["n"] += 1
        for k in ("hit", "floor", "active", "perceived"):
            c[{"hit": "hit", "floor": "floor", "active": "active", "perceived": "perc"}[k]] += 1 if r[k] else 0

    print("=== 14.05 non-SE open-latch (E9 / Gate-2 fair retest) ===")
    print("DV = un-fakeable filesystem existence (target file created / entry added / row repaired).\n")
    print(f"{'condition':<14}{'n':>3}  {'DV-hit':>11}  {'doc-floor':>10}  {'active':>8}  {'read-doc':>9}")
    for cell in ("dangling", "no_ref", "bib_dangling", "bib_complete", "csv_broken", "csv_complete"):
        c = cells.get(cell)
        if not c:
            continue
        dv = f"{c['hit']}/{c['n']} ({pct(c['hit'], c['n'])})"
        fl = f"{c['floor']}/{c['n']}"
        ac = f"{c['active']}/{c['n']}"
        pe = f"{c['perc']}/{c['n']}"
        print(f"{cell:<14}{c['n']:>3}  {dv:>11}  {fl:>10}  {ac:>8}  {pe:>9}")

    g, k = cells.get("dangling"), cells.get("no_ref")
    if g and k and g["n"] and k["n"]:
        gr, kr = g["hit"] / g["n"], k["hit"] / k["n"]
        p = fisher_exact_2x2(g["hit"], g["n"] - g["hit"], k["hit"], k["n"] - k["hit"])
        print(f"\n--- F1 PRIMARY decision (dangling vs symmetric no_ref) ---")
        print(f"  dangling target_created={gr:.2f}  no_ref(floor)={kr:.2f}  diff={gr-kr:+.2f}  Fisher p={p:.4f}")
        binary = gr >= 0.5 and gr - kr >= 0.4
        perc = g["perc"] / g["n"]
        if binary and p < 0.05:
            print("  => FIRES — the dangling cross-ref pulls creation of the named target above the")
            print("     generic doc-creation floor. Gate-2 REOPENS: the discrepancy pull is NOT SE-bound;")
            print("     it is objective-defect + additive-resolution bound (E5 was confounded).")
        elif binary:
            print("  => PROMISING but UNDERPOWERED — dangling target_created is high and the control floor")
            print("     is low, but Fisher is n.s. at this n. This is NOT a null (engaged agents DID create")
            print("     the target). Scale to n>=10/arm to confirm the gate-2 reopening.")
        elif gr < 0.3 and perc >= 0.5:
            print("  => HIGH-PERCEPTION NULL = FALSIFICATION — agents read the doc, saw the 404, did not")
            print("     create the target. SE-specificity CONFIRMED with the confound removed (stronger than E5).")
        else:
            print("  => INCONCLUSIVE (low activity / low perception) — inspect tapes, scale n.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
