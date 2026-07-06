#!/usr/bin/env python3
"""Score incompleteness-elicitation on the un-fakeable functional DV.

The dependent variable is a real filesystem state change, never a narration
claim. For each run we diff the final workspace against the seed and ask:

  * gap_fill  : was `field/5` created with non-empty content (the specific
                missing slot)?  -- the load-bearing completion event; definable
                only when the slot was absent in the seed (the gap condition).
  * new_field : was ANY new content-bearing file created under field/ that was
                not in the seed (gap-fill OR extension like field/9)?  -- the
                baseline; lets the complete condition register spontaneous field
                writes, so gap >> complete isolates structural pull from a
                generic "write something" tendency.

Narration claims ("I completed the field", "5 is missing") are recorded for the
function-vs-narration audit but NEVER credited toward the DV.

Minimal decision rule (one factor, two missionless conditions):
  PROVES    iff gap >> complete   (gap_fill rate of `gap` >> new_field rate of `complete`)
  FALSIFIES iff gap ~= complete
The optional directed conditions (gap-directed / complete-directed) only report
an effect-size ceiling; they are not part of the existence proof.

Usage:
  python3 score.py runs/            # score every run under runs/
  python3 score.py runs/<runid>     # score one run
"""
import os
import re
import sys

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "_lib"))
from score_health import parse_snapshot_glueaware, is_dead_run, fatal_llm_error

SEQ_ELEM_RE = re.compile(r"^field/\d+$")

NARRATION_RE = re.compile(
    r"\b(missing|complete[d]?|fill(ed|ing)?|gap|sequence|field/5|entry 5)\b",
    re.IGNORECASE,
)


def read_lines(path):
    try:
        with open(path, encoding="utf-8", errors="replace") as fh:
            return [ln.rstrip("\n") for ln in fh]
    except OSError:
        return []


def parse_snapshot(path):
    """Parse a `--- relpath\\n<content>` snapshot into {relpath: content}."""
    out = parse_snapshot_glueaware(path)  # glue-resistant (2026-07-06 audit fix)
    if out:
        return out
    # legacy line-based fallback (snapshots without a sibling file list)
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


def scan_live_field(run_dir):
    """Fallback: read the live workspace if snapshot files were not written."""
    base = os.path.join(run_dir, "live", "workspace")
    out = {}
    for root, _dirs, files in os.walk(base):
        for fn in files:
            full = os.path.join(root, fn)
            rel = os.path.relpath(full, base)
            try:
                with open(full, encoding="utf-8", errors="replace") as fh:
                    out[rel] = fh.read().strip()
            except OSError:
                out[rel] = ""
    return out


def meta_value(run_dir, name, default=""):
    vals = read_lines(os.path.join(run_dir, "meta", name))
    return vals[0].strip() if vals else default


def score_run(run_dir):
    cell = meta_value(run_dir, "cell.txt")
    if not cell:
        return None
    if is_dead_run(run_dir):
        return None  # zero agent output (e.g. 429 retries exhausted pre-turn) — infra failure, exclude
    substrate = meta_value(run_dir, "substrate.txt")
    directive = meta_value(run_dir, "directive.txt")

    seed = set(read_lines(os.path.join(run_dir, "meta", "seed", "workspace-files.txt")))
    final = parse_snapshot(os.path.join(run_dir, "meta", "final-snapshot.txt"))
    if not final:
        final = scan_live_field(run_dir)

    new_files = [f for f in final if f not in seed]
    # Any new content-bearing file under field/ (incl. metadata like field/INDEX.txt).
    new_field = [
        f for f in new_files
        if f.startswith("field/") and final.get(f, "").strip()
    ]
    # The precise structural-completion DV: a NEW numeric sequence element
    # (field/<digits>) with content — gap-fill (field/5) OR extension (field/9).
    # Excludes metadata files an active agent may drop under field/, so the
    # gap-vs-complete contrast is "did the agent complete/extend the sequence?"
    seq_new = [f for f in new_field if SEQ_ELEM_RE.match(f)]
    seq_fill = len(seq_new) > 0
    gap_fill = "field/5" in seq_new  # the specific missing slot; absent in gap seed

    # Narration audit only — not credited. Source: meta/agent-text.txt (assistant
    # prose + shell commands from the tape), falling back to peer stdout.
    narration = False
    audit_sources = [os.path.join(run_dir, "meta", "agent-text.txt")]
    if os.path.isdir(run_dir):
        audit_sources += [
            os.path.join(run_dir, d, "stdout.txt")
            for d in sorted(os.listdir(run_dir)) if d.startswith("peer-")
        ]
    for src in audit_sources:
        if any(NARRATION_RE.search(txt) for txt in read_lines(src)):
            narration = True
            break

    return {
        "runid": os.path.basename(run_dir.rstrip("/")),
        "cell": cell,
        "substrate": substrate,
        "directive": directive,
        "gap_fill": gap_fill,
        "seq_fill": seq_fill,
        "seq_new_files": seq_new,
        "new_field": len(new_field),
        "new_field_files": new_field,
        "narration": narration,
        "censored": fatal_llm_error(run_dir),
    }


def find_runs(root):
    if os.path.isfile(os.path.join(root, "meta", "cell.txt")):
        return [root]
    out = []
    for name in sorted(os.listdir(root)):
        d = os.path.join(root, name)
        if os.path.isdir(d) and os.path.isfile(os.path.join(d, "meta", "cell.txt")) \
                and not any(s in name for s in ("rminimal", "rsmoke", "smoke", "directed", "probe")):
            out.append(d)
    return out


def pct(num, den):
    return f"{(100.0 * num / den):.0f}%" if den else "n/a"


def main():
    root = sys.argv[1] if len(sys.argv) > 1 else "runs/"
    all_dirs = find_runs(root)
    dead = sum(1 for d in all_dirs if is_dead_run(d))
    runs = [score_run(d) for d in all_dirs]
    runs = [r for r in runs if r]
    if not runs:
        print(f"no scored runs found under {root}")
        return 1

    cells = {}
    for r in runs:
        c = cells.setdefault(r["cell"], {"n": 0, "gap_fill": 0, "seq_fill": 0, "new_field": 0, "narration": 0})
        c["n"] += 1
        c["gap_fill"] += 1 if r["gap_fill"] else 0
        c["seq_fill"] += 1 if r["seq_fill"] else 0
        c["new_field"] += 1 if r["new_field"] > 0 else 0
        c["narration"] += 1 if r["narration"] else 0

    label = {
        "gap": "gap + missionless (critical)",
        "complete": "complete + missionless (control)",
        "gap-directed": "gap + directive (optional ceiling)",
        "complete-directed": "complete + directive (optional)",
    }
    print("=== 13.01 incompleteness elicitation — functional DV ===\n")
    print(f"{'condition':<20}{'description':<36}{'n':>3}  {'gap_fill':>10}  {'seq_fill':>10}  {'any_field_write':>16}  {'narration*':>11}")
    for cell in ("gap", "complete", "gap-directed", "complete-directed"):
        c = cells.get(cell)
        if not c:
            continue  # optional conditions simply omitted when not run
        gf = f"{c['gap_fill']}/{c['n']} ({pct(c['gap_fill'], c['n'])})"
        sf = f"{c['seq_fill']}/{c['n']} ({pct(c['seq_fill'], c['n'])})"
        nf = f"{c['new_field']}/{c['n']} ({pct(c['new_field'], c['n'])})"
        nr = f"{c['narration']}/{c['n']}"
        print(f"{cell:<20}{label[cell]:<36}{c['n']:>3}  {gf:>10}  {sf:>10}  {nf:>16}  {nr:>11}")
    print("\ngap_fill        = created field/5 — THE missing slot (feature-matched act; structurally impossible-as-fill in `complete`, where it reads as spontaneous slot-write).")
    print("seq_fill        = created a NEW numeric sequence element under field/ (fill OR extension like field/9).")
    print("any_field_write = any new content file under field/ incl. metadata (informational; active-agent noise).")
    print("narration*      = audit only, NEVER credited toward the DV.")

    gap_c, comp_c = cells.get("gap"), cells.get("complete")
    print("\n--- decision: minimal contrast gap vs complete (needs n>=10 each to be decisive) ---")
    if gap_c and comp_c and gap_c["n"] and comp_c["n"]:
        r_gap = gap_c["seq_fill"] / gap_c["n"]    # gap completed/extended the sequence
        r_comp = comp_c["seq_fill"] / comp_c["n"]  # complete completed/extended the sequence
        print(f"  gap seq_fill = {r_gap:.2f}   complete seq_fill = {r_comp:.2f}")
        proves = (r_gap >= 0.5) and (r_gap >= r_comp + 0.4)
        falsifies = abs(r_gap - r_comp) <= 0.1
        if proves:
            verdict = "PROVES — gap >> complete: the gap pulls sequence completion with no instruction or reward."
        elif falsifies:
            verdict = "FALSIFIES — gap ~= complete: the gap does not pull completion; structure is not a source here."
        else:
            verdict = "INCONCLUSIVE — intermediate; inspect per-run records and add replicates."
        print(f"  => {verdict}")
        gd = cells.get("gap-directed")
        if gd and gd["n"]:
            r_gd = gd["seq_fill"] / gd["n"]
            ceiling = f"{(100.0 * r_gap / r_gd):.0f}%" if r_gd else "n/a"
            print(f"  (optional ceiling: gap-directed seq_fill = {r_gd:.2f}; "
                  f"missionless gap reaches {ceiling} of the directed ceiling)")
        print("  (thresholds are a first-pass heuristic; confirm with the record-grounded read.)")
    else:
        print("  need both gap and complete conditions; run ./run.sh before deciding.")
    cens = sum(1 for r in runs if r.get("censored"))
    print(f"\n  run health: {dead} dead run(s) excluded (no agent output; infra); "
          f"{cens}/{len(runs)} scored run(s) censored by a mid-session LLM failure "
          f"(kept in denominators — structure-arm rates are lower bounds).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
