#!/usr/bin/env python3
"""13.03 rupture: DV = the broken artifact becomes valid JSON (un-fakeable)."""
import os, sys, json
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
def seed_snap(d): return parse_snapshot(os.path.join(d,"meta","seed","workspace-snapshot.txt"))
def find_runs(root):
    if os.path.isfile(os.path.join(root,"meta","cell.txt")): return [root]
    return [os.path.join(root,n) for n in sorted(os.listdir(root)) if os.path.isfile(os.path.join(root,n,"meta","cell.txt")) and not any(s in n for s in ("rminimal","rsmoke","smoke","directed","probe"))]
def pct(a,b): return f"{100.0*a/b:.0f}%" if b else "n/a"
def valid_json(t):
    try: json.loads(t); return True
    except Exception: return False
def score_run(d):
    cell=meta(d,"cell.txt")
    if not cell: return None
    if is_dead_run(d):
        return None  # zero agent output (e.g. 429 retries exhausted pre-turn) — infra failure, exclude
    seed=seed_snap(d); final=final_snap(d)
    sj=seed.get("data.json",""); fj=final.get("data.json","")
    sv,fv=valid_json(sj),valid_json(fj)
    return {"cell":cell,"repaired":(not sv) and fv,"final_valid":fv,"changed":sj.strip()!=fj.strip(),"censored":fatal_llm_error(d)}
def main():
    root=sys.argv[1] if len(sys.argv)>1 else "runs/"
    all_dirs=find_runs(root)
    dead=sum(1 for x in all_dirs if is_dead_run(x))
    runs=[r for r in (score_run(x) for x in all_dirs) if r]
    if not runs: print("no runs"); return 1
    cells={}
    for r in runs:
        c=cells.setdefault(r["cell"],{"n":0,"rep":0,"chg":0,"val":0})
        c["n"]+=1; c["rep"]+=int(r["repaired"]); c["chg"]+=int(r["changed"]); c["val"]+=int(r["final_valid"])
    print("=== 13.03 rupture — functional DV (artifact becomes valid JSON) ===\n")
    print("condition     n  repaired         final_valid  changed")
    for cell in ("rupture","intact"):
        c=cells.get(cell)
        if not c: continue
        rep=f"{c['rep']}/{c['n']} ({pct(c['rep'],c['n'])})"; val=f"{c['val']}/{c['n']}"; chg=f"{c['chg']}/{c['n']}"
        print(f"{cell:<12}{c['n']:>3}  {rep:<15}  {val:<11}  {chg}")
    rup,inta=cells.get("rupture"),cells.get("intact")
    print("\n--- decision: rupture vs intact ---")
    if rup and inta and rup["n"] and inta["n"]:
        rr=rup["rep"]/rup["n"]; ic=inta["chg"]/inta["n"]
        print(f"  rupture repaired = {rr:.2f}   intact spurious-change = {ic:.2f}")
        if rr>=0.5 and rr>=ic+0.4: print("  => PROVES — a contradiction pulls repair with no instruction/reward.")
        elif rr<=0.1: print("  => FALSIFIES — no pull.")
        else: print("  => INCONCLUSIVE.")
    cens=sum(1 for r in runs if r.get("censored"))
    print(f"\n  run health: {dead} dead run(s) excluded (no agent output; infra); "
          f"{cens}/{len(runs)} scored run(s) censored by a mid-session LLM failure "
          f"(kept in denominators — structure-arm rates are lower bounds).")
    return 0
if __name__=="__main__": sys.exit(main())
