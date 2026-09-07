## Period IX: Control-Plane Hardening & ALIFE Manuscript (Mar 22–23)

The repository's control plane became machine-enforceable through surface contracts, frontmatter metadata, and Obsidian integration. The ALIFE manuscript underwent a fundamental reorientation from theoretical positioning to evidence-centric structure.

### Mar 22 — Surface contract infrastructure
**Key commits:** `18d3e7b`
**What changed:** 8 new validation scripts in `scripts/` (1000+ lines total), 8 new experiment surface declarations, new `experiments/PROTOCOL.md` (+38 lines), `.githooks/` expansion with pre-push, Makefile restructured with control-plane checks. 82 files changed.

---

### Mar 22 — Frontmatter metadata pilot
**Key commits:** `5e9b53e`
**What changed:** Introduced frontmatter metadata as a queryable control-plane surface. Piloted on p6-autopoiesis experiments (6.1–6.4) with a 149-line validator script.

---

### Mar 22 — Obsidian observer layer
**Key commits:** `1134dfe`, `be217c4`
**What changed:** Made Obsidian a checked-in observation surface with queryable Bases for frontmatter gaps, experiment inventory, and paper dossiers. Then pulled back automation to keep the layer declarative — removed premature helper script.
**Friction:** Initial approach was too automated; Bases are useful as artifacts but scripted interaction was premature.

---

### Mar 22 — Experiment frontmatter rollout
**Key commits:** `b15ac87`, `ec51256`, `62b5e40`, `c06819d`
**What changed:** Frontmatter metadata added to all experiment families: p3-mrcr (3.2–3.4), p1 lineage (2.7–2.9), p4 harbor (4.1–4.3), p1 coordination (2.1, 2.1b, 2.3), and remaining p1 experiments (1.1, 1.2, 2.5, 2.10, 2.11, 3.1). Entire experiment inventory now structurally queryable.

---

### Mar 22 — Phase overview surfaces
**Key commits:** `2001e9d`
**What changed:** Phase-level README entries for p1–p6. Created `phase-overviews.base` for coherent Obsidian querying.

---

### Mar 22 — Massive validator expansion
**Key commits:** `c980ca6`, `df5218c`
**What changed:** 13 new validation scripts (2000+ lines) covering experiment incubators, tasks, Obsidian bases, paper templates, tracks, root docs, scripts, surfaces, tests. Git hooks surface made explicitly indexed and owner-resolved. Pre-commit expanded +151 lines.

---

### Mar 22 — 2.6 promoted as ALIFE §4.2 hinge
**Key commits:** `4c76e80`, `e97a265`
**What changed:** Promoted 2.6 redesign as structural backbone for ALIFE manuscript section 4.2 (environmental memory). Bounded causal-memory ablation evidence reshaped the paper's central argument.

---

### Mar 22 — Massive 2.6 evidence retention
**Key commits:** `d25db60`, `6ec956c`
**What changed:** Preserved complete canonical run artifacts for 2.6: library workspace files (14,325 files), no-lineage prompt control variant (2,045 files) as paired comparison evidence.

---

### Mar 22 — ALIFE manuscript evidence-centric refocus
**Key commits:** `b1ebc3b`
**What changed:** Rewrote all 6 manuscript sections to build narrative around concrete experimental evidence rather than theoretical positioning.

---

### Mar 22 — Runtime retention export surface
**Key commits:** `9d6b1ef`
**What changed:** Added runtime infrastructure for structured retention metadata export: `runtime.go` (+87 lines), config (+14 lines), `.env.example` (+3 lines), protocol (+17 lines).

---

### Mar 22 — Evidence ladder refresh with control experiments
**Key commits:** `1008d8c`, `a441b27`
**What changed:** Refreshed evidence ladder with control experiments for 1.1-ariadne (paired exec-enabled vs exec-disabled conditions, `run.sh` +292 lines) and 3.1-habitat-guardian. Realigned manuscript narrative.

---

### Mar 22 — Manuscript public-facing reset
**Key commits:** `abb9a27`
**What changed:** Manuscript-plan restructured (499 changes). All sections rewritten for public-facing clarity. Transition from internal documentation to conference-ready manuscript.

---

### Mar 22 — Explicit retention directory migration
**Key commits:** `3807dea`, `fdfade1`
**What changed:** Migrated all experiment runners (p1-borges, p3-mrcr, p4-harbor, and all archived experiments) from implicit `QUINE_DATA_DIR` to explicit retention directory specification.
**Friction:** Implicit data directory made experiment reproducibility fragile.

---

### Mar 22 — Archive and incubator surface normalization
**Key commits:** `22a4693`
**What changed:** Normalized archive (legacy) and incubator (early-stage) experiment surfaces into the control-plane schema with validators.

---

### Mar 23 — Evidence-centric manuscript reorientation
**Key commits:** `057eac8`
**What changed:** Replaced traditional manuscript sections (introduction/positioning/evidence/discussion/conclusion) with mechanism-focused structure grounded in experimental evidence (4.1-continuity, 4.2-environmental-memory, 4.3-detached-maintenance). Manuscript-plan +341 lines. 418 files changed.

---

### Mar 23 — AGENTS.md streamlining
**Key commits:** `e3cd308`
**What changed:** Streamlined AGENTS.md with refined agent role definitions.

---

### Mar 23 — CLOSURE.md extraction
**Key commits:** `60dadf9`
**What changed:** Extracted closure protocol from AGENTS.md into dedicated `CLOSURE.md` (+49 lines).
**Friction:** Closure protocol had grown into a sub-system with its own contracts but was embedded inside AGENTS.md. Repeated refinement commits signaled it was outgrowing its container.

---

### Mar 23 — Accessible coherence principle
**Key commits:** `1a0cf67`
**What changed:** Added "Accessible Coherence" principle to the Constitution: a rule only contributes to coherence if it is reachable on the path where decisions are made.

---

### Mar 23 — Interactive service design direction
**Key commits:** `388814c`
**What changed:** Established architectural priority: reliable signal/control protocol mechanisms take precedence over natural language chat injection for interactive services. Documented in DEVELOPMENT.md (+45 lines) and ROADMAP.md (+10 lines).

---
