#!/usr/bin/env python3
"""13.06 topological gap: DV = the dangling reference is wired (parts/gamma created)."""
import os, sys
sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "..", "_lib"))
from score_health import parse_snapshot_glueaware, is_dead_run, fatal_llm_error
def read_lines(p):
    try: return [l.rstrip("\n") for l in open(p,encoding="utf-8",errors="replace")]
    except OSError: return []
def parse_snapshot(p):
    out = parse_snapshot_glueaware(p)  # glue-resistant (2026-07-06 audit fix)
    if out: return out
    # legacy line-based fallback (snapshots without a sibling file list)
    out,cur,buf={},None,[]
    for ln in read_lines(p):
        if ln.startswith("--- "):
            if cur is not None: out[cur]="\n".join(buf).strip()
            cur,buf=ln[4:].strip(),[]
        elif cur is not None: buf.append(ln)
    if cur is not None: out[cur]="\n".join(buf).strip()
    return out
def scan_live(rd):
    base=os.path.join(rd,"live","workspace"); out={}
    for root,_,fs in os.walk(base):
        for fn in fs:
            full=os.path.join(root,fn); rel=os.path.relpath(full,base)
            try: out[rel]=open(full,encoding="utf-8",errors="replace").read().strip()
            except OSError: out[rel]=""
    return out
def meta(rd,name,d=""):
    v=read_lines(os.path.join(rd,"meta",name)); return v[0].strip() if v else d
def final_snap(d): return parse_snapshot(os.path.join(d,"meta","final-snapshot.txt")) or scan_live(d)
def seed_files(d):
    return set(read_lines(os.path.join(d,"meta","seed","workspace-files.txt")))
def find_runs(root):
    if os.path.isfile(os.path.join(root,"meta","cell.txt")): return [root]
    return [os.path.join(root,n) for n in sorted(os.listdir(root)) if os.path.isfile(os.path.join(root,n,"meta","cell.txt")) and not any(s in n for s in ("rminimal","rsmoke","smoke","directed","probe"))]
def pct(a,b): return f"{100.0*a/b:.0f}%" if b else "n/a"
def score_run(d):
    cell=meta(d,"cell.txt")
    if not cell: return None
    if is_dead_run(d):
        return None  # zero agent output (e.g. 429 retries exhausted pre-turn) — infra failure, exclude
    seedf=seed_files(d); final=final_snap(d)
    wired=bool(final.get("parts/gamma","").strip()) and "parts/gamma" not in seedf
    return {"cell":cell,"wired":wired,"censored":fatal_llm_error(d)}
def main():
    root=sys.argv[1] if len(sys.argv)>1 else "runs/"
    all_dirs=find_runs(root)
    dead=sum(1 for x in all_dirs if is_dead_run(x))
    runs=[r for r in (score_run(x) for x in all_dirs) if r]
    if not runs: print("no runs"); return 1
    cells={}
    for r in runs:
        c=cells.setdefault(r["cell"],{"n":0,"w":0}); c["n"]+=1; c["w"]+=int(r["wired"])
    print("=== 13.06 topological gap — functional DV (dangling reference wired) ===\n")
    print("condition     n  wired")
    for cell in ("dangling","closed"):
        c=cells.get(cell)
        if not c: continue
        w=f"{c['w']}/{c['n']} ({pct(c['w'],c['n'])})"
        print(f"{cell:<12}{c['n']:>3}  {w}")
    dg,cl=cells.get("dangling"),cells.get("closed")
    print("\n--- decision: dangling vs closed ---")
    if dg and cl and dg["n"] and cl["n"]:
        rd_=dg["w"]/dg["n"]; rc=cl["w"]/cl["n"]
        print(f"  dangling wired = {rd_:.2f}   closed wired = {rc:.2f}")
        if rd_>=0.5 and rd_>=rc+0.4: print("  => PROVES — a relational dangle pulls wiring with no instruction/reward.")
        elif abs(rd_-rc)<=0.1: print("  => FALSIFIES — no pull.")
        else: print("  => INCONCLUSIVE.")
    cens=sum(1 for r in runs if r.get("censored"))
    print(f"\n  run health: {dead} dead run(s) excluded (no agent output; infra); "
          f"{cens}/{len(runs)} scored run(s) censored by a mid-session LLM failure "
          f"(kept in denominators — structure-arm rates are lower bounds).")
    return 0
if __name__=="__main__": sys.exit(main())
