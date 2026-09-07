#!/usr/bin/env python3
"""Generate the experiments/active/ family skeleton + READMEs (Phase 1).

Reproducible: reads id-map.tsv (active rows) + the persisted inventory for each
experiment's question/status/theory_objects, and writes:
  - experiments/active/README.md            (family router)
  - experiments/active/<family>/README.md   (navigation-only family index)

Family READMEs carry NO theory content (P2) and NO experiment frontmatter, so
the surface validators treat them as navigation indexes. Experiment tables list
each experiment by NEW id with its legacy id; per-row clickable links are added
in Phase 3 once experiments physically move under the family (until then they
remain canonical at their legacy paths).
"""
from __future__ import annotations

import json
import pathlib

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[1]
INV = HERE / "inventory"
ACTIVE = ROOT / "experiments" / "active"

PREFIX_ORDER = [
    "autopoiesis-succession", "structural-elicitation", "cooperation-dynamics",
    "trace-inheritance", "early-existence", "dormancy-ecology", "benchmarks",
    "null-ecology", "digital-red-queen", "cumulative-novelty", "social-role",
]

PREFIX = {
    "autopoiesis-succession": "succ", "dormancy-ecology": "dorm",
    "trace-inheritance": "trce", "null-ecology": "null",
    "cooperation-dynamics": "coop", "social-role": "role",
    "digital-red-queen": "drq", "structural-elicitation": "elic",
    "cumulative-novelty": "novl", "early-existence": "exst", "benchmarks": "bench",
}

SPEC = {
    "autopoiesis-succession": {
        "scope": "Successor regeneration, carrier succession, and post-handoff "
        "continuity: whether a Quine lineage produces a viable successor without "
        "being told to, what minimal substrate makes self-birth reappear, and "
        "what operational carrier (not formal source) survives the handoff "
        "barrier. Merges the p6 self-birth ladder (canonical rungs only), the "
        "p6f future-addressability branch, the p6m successor-morphospace assays, "
        "the p6r agentic-runtime-replacement controls, and the p7 self-source "
        "continuation line (7.4-7.8).",
        "boundary": [
            "distinct from `dormancy-ecology` — no quiescence/wake selection pressure",
            "distinct from `trace-inheritance` — the carrier is the live process/"
            "runtime surface, not a stigmergic medium read by a later stateless peer",
        ],
        "ops": [
            "absorbs `p6-autopoiesis/lib/` and `p6m-successor-morphospace/lib/` "
            "prompt-posture helpers into a family-local `lib/` (Phase 3)",
            "absorbs the P6 `EXPERIMENT-LOG.md` into `analysis/` (Decision 8)",
            "47 active rungs; the other 47 p6 experiments are archived to "
            "`archive/legacy-phase-6/` (citation-validated triage, Appendix B)",
        ],
    },
    "structural-elicitation": {
        "scope": "Environmental structure as a third source of behavior beside "
        "instruction and reward: whether the shape of the workspace itself "
        "(gaps, gradients, grooves, ruptures) elicits behavior in a missionless "
        "single agent under an un-fakeable functional DV. Was p13 (13.x anchor "
        "tests + the 14.x non-linguistic-affordance extension).",
        "boundary": [
            "distinct from `cumulative-novelty` — structure elicits, it is not a "
            "frontier the agent expands",
            "distinct from `null-ecology` — structure is the independent variable, "
            "not the absence of one",
        ],
        "ops": [
            "actively running (live cross-model wave) — frontmatter seeded in "
            "Phase 1; the directory MOVE is deferred until the wave settles",
            "Python scoring helpers (`_lib/`) carry absolute paths — fixed on move",
        ],
    },
    "cooperation-dynamics": {
        "scope": "Scarcity-driven cooperation, population coordination, and social "
        "structure (including adversarial pressure) over a shared substrate: "
        "when and how missionless peers coordinate, mechanize coordination "
        "organs, or defend against a hostile peer. Was p8, plus the "
        "`2.11-saboteur` adversarial-commons experiment.",
        "boundary": [
            "distinct from `social-role` — about whether/how coordination happens, "
            "not which roles spontaneously differentiate",
            "distinct from `digital-red-queen` — cooperation/scarcity, not "
            "adversarial selection over OS routing surfaces",
        ],
        "ops": [
            "absorbs the p8 `BEHAVIOR-ANALYSIS.md` into `analysis/`",
            "p8 children share an `8.11` world-setup script (intra-family dep, Appendix G)",
            "6 early p8 stubs (8.2-8.7) archived to `archive/legacy-p8-early/`",
        ],
    },
    "trace-inheritance": {
        "scope": "Stigmergic trace causality and cross-generation niche repair: "
        "how residue a later, stateless newcomer reads carries inheritance "
        "across lifetimes, and when a trace acts as an environmental carrier "
        "versus a merely narrated protocol. Merges the p7 generation/trace/"
        "irregularity/seeded-workspace ecologies and the p1 stigmergy line "
        "(2.3, 2.5-2.9).",
        "boundary": [
            "distinct from `autopoiesis-succession` — inheritance is via an "
            "environmental medium, not a live carrier across one handoff",
            "distinct from `dormancy-ecology` — no dormancy/wake selection",
        ],
        "ops": [
            "the P7G/P7T trace-causality backbone is the family's strongest signal",
            "p1 lineage runners cross-reference `2.7-archivist-index` (Appendix G)",
        ],
    },
    "early-existence": {
        "scope": "Historical existence proofs: bounded search, recursion, "
        "retrieval, and self-replication under physical constraints — the "
        "earliest demonstrations that the substrate can do the thing at all. "
        "Merges the p1 borges existence experiments (1.x, 2.1/2.1b, 3.x), the "
        "p3 MRCR retrieval line, and the p5 JOOTS meta-cognition design.",
        "boundary": [
            "distinct from the science families — kept for genuine historical "
            "existence proofs, not as a catch-all for old experiments",
            "p1 inheritance/role experiments are re-homed to their scientific "
            "families (trace-inheritance, social-role, cooperation-dynamics)",
        ],
        "ops": [
            "p1/p3/p5 runners use depth-3 root-reach climbs — fixed +1 on move",
        ],
    },
    "dormancy-ecology": {
        "scope": "Explicit quiescence, dormant reservoirs, and wake-selection "
        "ecology: whether missionless peers enter dormancy without instruction, "
        "and how later entrants select which dormant peers wake — distinguishing "
        "host-mechanic wakes from artifact/agent-decided wakes. Was p12.",
        "boundary": [
            "distinct from `autopoiesis-succession` — selection is over dormant "
            "peers, not successor morphology",
            "distinct from `null-ecology` — dormancy is an explicit ecological "
            "state under wake pressure, not a missionless behavior gradient",
        ],
        "ops": [
            "p12 tapes embed absolute paths (retained evidence) — migration note "
            "on move, not rewritten (Decision 7)",
        ],
    },
    "benchmarks": {
        "scope": "External benchmark evaluation and evaluation infrastructure: "
        "Harbor/terminal-bench, ProgramBench, and the persistent compiler-"
        "construction ecology. Evaluation scaffolds, not theory-feeding "
        "experiments. Merges p4 (Harbor), p10 (benchmark-native), and p9c "
        "(compiler ecology).",
        "boundary": [
            "navigation/infrastructure family — its experiments carry no theory "
            "objects by default (the p9c coordination *finding* routes to "
            "protocol-emergence/coordination-closure, but the artifact is "
            "infrastructure-shaped)",
        ],
        "ops": [
            "p4 `tasks/` Harbor fixtures embed paths — likely stay via symlink "
            "exemption (Open Question 4)",
            "p9c is an active campaign (9C.01, I01-I33 runs)",
        ],
    },
    "null-ecology": {
        "scope": "Missionless behavior gradient and null-state operationalization: "
        "what agents do when no mission, reward, or imposed structure pulls "
        "them, and how to make the null state measurable. Merges p7e and the "
        "p7-oee affordance-rich null observations (7.9, 7.10, 7.13).",
        "boundary": [
            "distinct from `structural-elicitation` — absence of structure is the "
            "variable, not structure as an elicitor",
            "distinct from `cumulative-novelty` — no frontier expansion claim",
        ],
        "ops": [],
    },
    "digital-red-queen": {
        "scope": "Adversarial co-evolution on OS routing surfaces: session-memory "
        "selection, lineage predation, and the symmetry trap, where selection "
        "must distinguish state an organism produced from state it merely "
        "narrated. Was p11.",
        "boundary": [
            "distinct from `cooperation-dynamics` — adversarial selection over "
            "runtime routing, not scarcity-driven cooperation",
        ],
        "ops": [
            "`records/matches.jsonl` is the canonical run ledger",
        ],
    },
    "cumulative-novelty": {
        "scope": "Wiki/knowledge-artifact frontier expansion: whether persistent "
        "agents grow a knowledge frontier across rounds and what the first "
        "novelty basins are. Merges p7w and the p7-oee habitat assays (7.1-7.3).",
        "boundary": [
            "distinct from `structural-elicitation` — the agent expands a frontier "
            "rather than responding to imposed structure",
        ],
        "ops": [],
    },
    "social-role": {
        "scope": "Missionless role differentiation: whether spontaneous labor "
        "division and group organization emerge without assigned roles. Kept "
        "deliberately small (role differentiation only). Merges p7r, the p7-oee "
        "role observation (7.11), and p1 `2.10-code-commons`.",
        "boundary": [
            "distinct from `cooperation-dynamics` — about which roles differentiate, "
            "not whether coordination/cooperation happens at all",
        ],
        "ops": [],
    },
}


def load_active():
    fams: dict[str, list] = {}
    for ln in (HERE / "id-map.tsv").read_text().splitlines():
        if ln.startswith("#") or not ln.strip():
            continue
        lp, np, lid, nid, fam, disp = ln.split("\t")
        if disp != "active":
            continue
        fams.setdefault(fam, []).append(
            {"lp": lp, "new_id": nid, "legacy_id": lid, "leaf": np.split("/")[-1]})
    return fams


def load_meta():
    meta = {}
    for f in INV.glob("inv-*.json"):
        for e in json.loads(f.read_text())["experiments"]:
            meta[e["legacy_path"]] = {
                "question": e.get("question", "").strip(),
                "status": e.get("status", "").strip(),
                "theory_objects": e.get("theory_objects", []),
            }
    return meta


def fam_readme(fam, rows, meta) -> str:
    rows = sorted(rows, key=lambda r: r["new_id"])
    objs = sorted({o for r in rows for o in meta.get(r["lp"], {}).get("theory_objects", [])})
    spec = SPEC[fam]
    out = [f"# Family: {fam}", ""]
    out += ["## Scientific scope", "", spec["scope"], ""]
    out += ["## Theory objects touched", "",
            (", ".join(objs) if objs else "_(navigation/infrastructure family — no theory objects)_"), ""]
    out += ["## Experiments", "",
            f"{len(rows)} experiments, homed here under family-prefixed IDs "
            "(legacy phase IDs preserved in the last column and in each "
            "experiment's `legacy_id` frontmatter).", "",
            "| ID | Question | Status | Theory objects | Legacy ID |",
            "|----|----------|--------|----------------|-----------|"]
    for r in rows:
        m = meta.get(r["lp"], {})
        q = m.get("question", "").replace("|", "/")
        if len(q) > 90:
            q = q[:87] + "..."
        to = ", ".join(m.get("theory_objects", [])) or "—"
        out.append(f"| [`{r['new_id']}`](./{r['leaf']}/) | {q} | {m.get('status','')} | {to} | `{r['legacy_id']}` |")
    out += ["", "## Boundary with other families", ""]
    out += [f"- {b}" for b in spec["boundary"]]
    out += ["", "## Operational notes", ""]
    out += ([f"- {o}" for o in spec["ops"]] if spec["ops"] else ["- none"])
    if fam == "autopoiesis-succession":
        out += ["", "## Phase-Standard Prompt Posture", "",
                "The p6-origin self-birth/ablation runners share one canonical "
                "prompt-posture helper rather than copying the env block: "
                "[`./lib/prompt-posture.sh`](./lib/prompt-posture.sh). Runners "
                "that use the shared posture must `source` it (not inline it).",
                "",
                "- audit: `./scripts/check-experiment-prompt-posture.sh`",
                "- strict: `./scripts/check-experiment-prompt-posture.sh --strict`"]
    out += ["", "---", "",
            "_A family directory groups experiments by scientific neighborhood "
            "for navigation only; it owns no theory content (scientific ownership "
            "lives in [`Paper/theory/objects/`](../../../Paper/theory/objects/)). "
            "See [`../../restructure-design.md`](../../restructure-design.md)._", ""]
    return "\n".join(out)


def active_root_readme(fams, meta) -> str:
    out = ["# Active Experiments — Family Index", "",
           "> **Migration in progress.** This is the family-grouped home for "
           "active experiments (`experiments/active/<family>/<id>/`). During the "
           "migration window the legacy phase directories (`p1-borges/`, "
           "`p6-autopoiesis/`, …) remain canonical until Phase 3 moves their "
           "experiments here and Phase 4 removes them. See "
           "[`../restructure-design.md`](../restructure-design.md) and the "
           "[`../../development/experiments-restructure/`](../../development/experiments-restructure/) "
           "tooling.", "",
           "Families group experiments by **scientific neighborhood** for "
           "navigation; they own no theory content. Scientific ownership lives "
           "in [`../../Paper/theory/objects/`](../../Paper/theory/objects/); each "
           "experiment declares its `theory-objects` in frontmatter.", "",
           "| Family | Active | Scientific neighborhood |",
           "|--------|--------|-------------------------|"]
    one_line = {
        "autopoiesis-succession": "successor regeneration / carrier succession / post-handoff continuity",
        "structural-elicitation": "environmental structure as a third behavior source",
        "cooperation-dynamics": "scarcity cooperation / population coordination / adversarial pressure",
        "trace-inheritance": "stigmergic trace causality / cross-generation niche repair",
        "early-existence": "historical existence proofs (search, recursion, retrieval)",
        "dormancy-ecology": "quiescence / dormant reservoirs / wake selection",
        "benchmarks": "external benchmark + compiler-ecology infrastructure",
        "null-ecology": "missionless behavior gradient / null-state operationalization",
        "digital-red-queen": "adversarial co-evolution on OS routing surfaces",
        "cumulative-novelty": "wiki/knowledge-artifact frontier expansion",
        "social-role": "missionless role differentiation",
    }
    for fam in PREFIX_ORDER:
        n = len(fams.get(fam, []))
        out.append(f"| [`{fam}/`](./{fam}/) | {n} | {one_line[fam]} |")
    out += [f"| **Total** | **{sum(len(v) for v in fams.values())}** | |", "",
            "## Index Projection Maintenance", "",
            "Source owner: each `<family>/README.md` owns its experiment table; "
            "experiment READMEs own experiment truth; "
            "`development/experiments-restructure/id-map.tsv` is the canonical "
            "legacy→new map.", "",
            "Projection role: this index routes to families and reports active "
            "counts. It does not restate experiment or theory content.", "",
            "Freshness trigger: regenerate when families gain/lose experiments "
            "(via `generate-family-readmes.py`).", "",
            "Drift check: `make check-control-plane`; per-family counts vs id-map.", "",
            "Absorption owner: if this index lags, regenerate it from id-map "
            "rather than hand-editing.", ""]
    return "\n".join(out)


def main() -> int:
    fams = load_active()
    meta = load_meta()
    ACTIVE.mkdir(parents=True, exist_ok=True)
    (ACTIVE / "README.md").write_text(active_root_readme(fams, meta))
    for fam, rows in fams.items():
        d = ACTIVE / fam
        d.mkdir(parents=True, exist_ok=True)
        (d / "README.md").write_text(fam_readme(fam, rows, meta))
    print(f"wrote experiments/active/README.md + {len(fams)} family READMEs")
    for fam in PREFIX_ORDER:
        print(f"  {PREFIX[fam]:5s} {fam:24s} {len(fams.get(fam, []))}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
