#!/usr/bin/env python3
"""Score semantic-collision on the un-fakeable functional DV.

The dependent variable is a real filesystem state change, never a narration
claim. Two surfaces declare the service version:

  * config/service.conf : a line  VERSION=<x>
  * config/manifest.json: a field "version": "<y>"

In the `collision` arm the seed values disagree (conf 2.4.0, manifest 2.3.0); in
the `consistent` arm they agree (both 2.3.0). Both surfaces are individually valid
and complete — the only manipulated bit is whether their MEANINGS collide.

For each run we parse both versions from the seed and the final snapshot:

  * resolved      : the two versions DISAGREED in the seed and are EQUAL at the
                    end (an alignment edit occurred). Definable only for the
                    collision arm; structurally impossible for `consistent`.
  * direction     : which value the field converged to — `to_conf` (manifest
                    rewritten to the conf value), `to_manifest` (conf rewritten
                    to the manifest value), or `to_other` (both moved to a third
                    value). Reveals the field's authority gradient.
  * version_touched: did either declared version change from its seed value? The
                    matched "did the active agent edit the version fields at all"
                    baseline, comparable across both arms.

Narration ("the versions are inconsistent") is NEVER credited toward the DV.

Decision rule (one factor, two missionless conditions):
  PROVES    iff collision.resolved >> consistent.version_touched
  FALSIFIES iff collision.resolved ~= 0  (the collision pulls no alignment)

Usage:
  python3 score.py runs/            # score every run under runs/
  python3 score.py runs/<runid>     # score one run
"""
import json
import os
import re
import sys

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "_lib"))
from score_health import parse_snapshot_glueaware, is_dead_run, fatal_llm_error

CONF_RE = re.compile(r"(?mi)^\s*version\s*=\s*([0-9][0-9A-Za-z.\-]*)")
MAN_RE = re.compile(r'"version"\s*:\s*"([^"]+)"')


def vtuple(v):
    try:
        return tuple(int(x) for x in str(v).split("."))
    except (ValueError, AttributeError):
        return (0,)


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


def scan_live(run_dir):
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


def conf_version(text):
    m = CONF_RE.search(text or "")
    return m.group(1) if m else None


def man_version(text):
    m = MAN_RE.search(text or "")
    return m.group(1) if m else None


def versions(snap):
    return conf_version(snap.get("config/service.conf", "")), \
        man_version(snap.get("config/manifest.json", ""))


def score_run(run_dir):
    cell = meta_value(run_dir, "cell.txt")
    if not cell:
        return None
    if is_dead_run(run_dir):
        return None  # zero agent output (e.g. 429 retries exhausted pre-turn) — infra failure, exclude

    seed = parse_snapshot(os.path.join(run_dir, "meta", "seed", "workspace-snapshot.txt"))
    final = parse_snapshot(os.path.join(run_dir, "meta", "final-snapshot.txt"))
    if not final:
        final = scan_live(run_dir)

    conf_s, man_s = versions(seed)
    conf_f, man_f = versions(final)

    collided_seed = bool(conf_s and man_s and conf_s != man_s)
    both_final = bool(conf_f and man_f)
    resolved = collided_seed and both_final and conf_f == man_f
    version_touched = (conf_f != conf_s) or (man_f != man_s)

    direction = ""
    to_higher = None
    if resolved:
        if conf_f == conf_s:
            direction = "to_conf"      # manifest moved to the conf value
        elif man_f == man_s:
            direction = "to_manifest"  # conf moved to the manifest value
        else:
            direction = "to_other"     # both moved to a third value
        if conf_s and man_s:
            hi = conf_s if vtuple(conf_s) >= vtuple(man_s) else man_s
            to_higher = (conf_f == hi)

    return {
        "runid": os.path.basename(run_dir.rstrip("/")),
        "cell": cell,
        "conf_seed": conf_s, "man_seed": man_s,
        "conf_final": conf_f, "man_final": man_f,
        "resolved": resolved,
        "direction": direction,
        "to_higher": to_higher,
        "version_touched": version_touched,
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
        c = cells.setdefault(r["cell"], {"n": 0, "resolved": 0, "touched": 0, "higher": 0, "dirs": {}})
        c["n"] += 1
        c["resolved"] += 1 if r["resolved"] else 0
        c["touched"] += 1 if r["version_touched"] else 0
        c["higher"] += 1 if r.get("to_higher") else 0
        if r["direction"]:
            c["dirs"][r["direction"]] = c["dirs"].get(r["direction"], 0) + 1

    label = {
        "collision_a": "collision: conf higher (conf 2.4.0 / man 2.3.0)",
        "collision_b": "collision: manifest higher (conf 2.3.0 / man 2.4.0)",
        "consistent":  "consistent + missionless (control)",
    }
    print("=== 13.07 semantic collision — functional DV (deconfounded direction) ===\n")
    print(f"{'condition':<13}{'description':<46}{'n':>3}  {'resolved':>10}  {'touched':>8}  {'direction':>22}  {'to_higher':>10}")
    for cell in ("collision_a", "collision_b", "consistent"):
        c = cells.get(cell)
        if not c:
            continue
        rs = f"{c['resolved']}/{c['n']} ({pct(c['resolved'], c['n'])})"
        tc = f"{c['touched']}/{c['n']}"
        dirs = ", ".join(f"{k}:{v}" for k, v in sorted(c["dirs"].items())) or "-"
        hi = f"{c['higher']}/{c['resolved']}" if c["resolved"] else "-"
        print(f"{cell:<13}{label.get(cell, cell):<46}{c['n']:>3}  {rs:>10}  {tc:>8}  {dirs:>22}  {hi:>10}")
    print("\nresolved  = versions disagreed in seed and equal at end (alignment edit).")
    print("direction = which SURFACE won (to_conf / to_manifest).")
    print("to_higher = of resolved runs, how many converged to the higher version.")

    a, b, con = cells.get("collision_a"), cells.get("collision_b"), cells.get("consistent")
    res = [c for c in (a, b) if c and c["n"]]
    print("\n--- binary: does the collision pull alignment? ---")
    if res:
        rates = [c["resolved"] / c["n"] for c in res]
        r_con = (con["touched"] / con["n"]) if (con and con["n"]) else 0.0
        print(f"  collision resolved = {min(rates):.2f}-{max(rates):.2f}   consistent touched = {r_con:.2f}")
        if min(rates) >= 0.5 and min(rates) >= r_con + 0.4:
            print("  => PROVES — semantic collision pulls cross-surface alignment, no instruction/reward.")
        elif max(rates) <= 0.1:
            print("  => FALSIFIES — no alignment pull.")
        else:
            print("  => INCONCLUSIVE — add replicates.")
    else:
        print("  need collision_a / collision_b runs.")

    print("\n--- deconfound: CONFIG-AUTHORITY vs HIGHER-WINS? ---")
    if a and b and a["resolved"] and b["resolved"]:
        a_conf = a["dirs"].get("to_conf", 0) / a["resolved"]
        b_conf = b["dirs"].get("to_conf", 0) / b["resolved"]
        a_hi = a["higher"] / a["resolved"]
        b_hi = b["higher"] / b["resolved"]
        print(f"  to_conf:   collision_a {a_conf:.2f}   collision_b {b_conf:.2f}")
        print(f"  to_higher: collision_a {a_hi:.2f}   collision_b {b_hi:.2f}")
        if a_conf >= 0.6 and b_conf >= 0.6:
            print("  => CONFIG-AUTHORITY — the conf surface wins regardless of which version is higher.")
        elif a_hi >= 0.6 and b_hi >= 0.6:
            print("  => HIGHER-WINS — the higher version wins regardless of surface (magnitude, not authority).")
        else:
            print("  => MIXED / no clean rule — inspect per-run reasoning.")
    else:
        print("  need both collision_a and collision_b to have resolved runs.")
    cens = sum(1 for r in runs if r.get("censored"))
    print(f"\n  run health: {dead} dead run(s) excluded (no agent output; infra); "
          f"{cens}/{len(runs)} scored run(s) censored by a mid-session LLM failure "
          f"(kept in denominators — structure-arm rates are lower bounds).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
