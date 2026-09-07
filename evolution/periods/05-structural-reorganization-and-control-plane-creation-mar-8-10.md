## Period V: Structural Reorganization & Control Plane Creation (Mar 8–10)

A comprehensive restructuring wave touched nearly every surface — paper workspace, experiment organization, testing infrastructure, and project governance — establishing the repository's first formal control plane documents.

### Mar 8 — Repository positioning clarification
**Key commits:** `8353af0`
**What changed:** Rewrote AGENTS.md, Paper README, and main README to clarify repository positioning and branch strategy.

---

### Mar 8 — Paper workspace restructure
**Key commits:** `3092794`, `38502d0`
**What changed:** Renamed `track-arxiv` → `arxiv` directories. Deleted old files and `full_paper.md`. Created `Paper/evidence-map.md`. Moved old paper sections to archive. Net reduction of 1400+ lines.
**Friction:** Inconsistent naming between paper tracks; stale artifacts from submission sprint.

---

### Mar 8 — FD+3 alignment
**Key commits:** `3d60e40`
**What changed:** Updated all fd references across behavior test prompts and e2e tests to use fd+3 convention. Synchronized documentation with actual runtime behavior across 37 files.
**Friction:** Documentation and tests used inconsistent fd numbering.

---

### Mar 8 — Testing guide creation
**Key commits:** `2e60196`, `4ca4f63`
**What changed:** Created `TESTING.md` (314 lines) covering e2e tests, behavior tests, and experiment protocols. Then compressed AGENTS.md and TESTING.md — net reduction of 173 lines through standardization.

---

### Mar 8 — Academic paper build pipeline
**Key commits:** `470a3fd`
**What changed:** Added complete ACM article LaTeX class (`acmart.cls` — 3655 lines), bibliography styles, section templates.

---

### Mar 8 — Experiment index reorganization
**Key commits:** `2780812`, `2c77d74`, `32a145d`, `4a08875`, `43a18c2`
**What changed:** Moved sweep logs to dedicated directory. Moved failure analysis to `analysis/`. Created dedicated outlines for all ALife tracks: `outline-alife-conf.md`, `outline-alife-journal.md`, `outline-arxiv.md`. Added `guideline-alife-conf.md` (365 lines).

---

### Mar 8 — Project control plane creation
**Key commits:** `18e8d5a`
**What changed:** Created `ROADMAP.md` (144 lines) and `STATUS.md` (168 lines) — first formal project control plane documents. Added evolution `TEMPLATE.md` (51 lines).

---

### Mar 8 — Bead issue tracker removal
**Key commits:** `4fc5265`
**What changed:** Exported all 94 bead tasks to `Paper/archive/notes/bd-tasks-archive.md` (757 lines). Deleted entire `.beads/` directory. Removed bead references from AGENTS.md, STATUS.md, ROADMAP.md.
**Friction:** Bead-based issue tracking had become maintenance overhead without proportional value.

---

### Mar 8 — Runtime physics refactoring
**Key commits:** `a70bd06`, `81a7130`
**What changed:** Execution budget now driven by environment variables. `e2e.sh` gained 177 lines of new tests. Includes budget-hard-fail behavior test recordings. Updated acceptance criteria in AGENTS.md.

---

### Mar 8 — Job tools and detach behavior
**Key commits:** `b8c39a0`, `bc794d0`
**What changed:** Added job tool behavior tests. Recorded daemon, detach, and detach-emergent test runs. Net reduction of 856 lines from cleanup alongside new feature testing.

---

### Mar 8 — Habitat guardian experiment
**Key commits:** `77f3cfd`
**What changed:** Created `experiments/active/early-existence/exst-05-habitat-guardian/` with `init-habitat.sh`, `tamper.sh`. Added behavior prompt `detach-overlap-emergent.md`. Updated system prompt and tool schemas. Expanded behavior test runner.

---

### Mar 9 — Material stream philosophy
**Key commits:** `011806e`, `6ae09c2`
**What changed:** Created `18_material_streams.md` (228 lines) defining how data flows through agent processes as physical material. Updated `prompt.go` with material/no-material prompt differentiation and tests.

---

### Mar 9 — Stdin stream physics
**Key commits:** `6f31e14`, `3fb956d`
**What changed:** Created `stdin-physics.md` behavior prompt (40 lines). Recorded stdin-physics test runs. Added stdin-specific e2e tests.

---

### Mar 9 — Subjective reality and time control redesign
**Key commits:** `7a7eecf`
**What changed:** Major rewrite of evolution 11-13: `11_subjective_reality.md` (464 lines), `12_process_as_fs.md`, `13_schrodingers_process.md`. New time control concepts for agent subjective experience.

---

### Mar 9 — Transactional workspace physics
**Key commits:** `6cb06ab`
**What changed:** Implemented transactional workspace using shadow filesystem (+5710 lines across 92 files). `e2e.sh` gained 69 lines of workspace transaction tests.

---

### Mar 9 — Linux overlay workspace
**Key commits:** `9cf141b`
**What changed:** Implemented Linux overlay filesystem for workspace isolation using overlayfs for copy-on-write semantics (+20K lines across 517 files). `e2e.sh` gained 86 lines of overlay-specific tests.

---

### Mar 9 — Archivist and lineage experiments
**Key commits:** `57f4439`, `9f2d25f`
**What changed:** Created `p1-borges/2.6-legacy-protocol-ablation/README.md` (280 lines), `p1-borges/2.7-archivist-index/` with full infrastructure (README 487 lines, `analyze.sh` 140 lines, `run.sh` 132 lines, `init-archive.sh` 97 lines), and `p1-borges/2.8-endogenous-index-lineage/` (README 272 lines, `analyze.sh` 195 lines, `run.sh` 129 lines).

---

### Mar 9 — Massive experiment archive cleanup
**Key commits:** `6ecf970`
**What changed:** Deleted ~961K lines of bulky harbor experiment artifacts across 20,969 files.
**Friction:** Repository size had grown unsustainably from raw experiment data.

---

### Mar 10 — Control plane maturation
**Key commits:** `5b6889c`, `cd27b87`, `6cbcb7c`
**What changed:** Major compression of ROADMAP.md and STATUS.md. Expanded `experiments/PROTOCOL.md` (+164 lines, then +481 lines). Refined AGENTS.md. Updated Paper README.

---

### Mar 10 — Experiment standardization
**Key commits:** `1d2fba5`
**What changed:** Added standardized README headers to all experiments across phases p1–p5.

---

### Mar 10 — Research control plane scripts
**Key commits:** `0890956`, `2966426`, `3f97fd4`, `b2a5fd6`
**What changed:** Created `scripts/check-experiment-paper-feeds.sh` (140 lines), `scripts/check-paper-dossier-links.sh` (157 lines), `scripts/check-paper-evidence-backbone.sh` (156 lines), `scripts/focus-paper.sh` (110 lines), and `scripts/code2prompt-paper-context.sh` (225 lines). Added paper-feed and reverse links between experiments and paper dossiers.

---

### Mar 10 — ALife paper reframing
**Key commits:** `0c0ee15`
**What changed:** Restructured ALife conference outline around "constraint-selected dynamics" framing. Added `_meta/calls/alife-2026-cfp.md` (102 lines — call for papers). Updated all experiment paper-feed links.

---

### Mar 10 — Lineage environmental cognition series
**Key commits:** `88d96d1`
**What changed:** Created `LINEAGE_ENVIRONMENTAL_COGNITION_SERIES.md` (274 lines) — multi-experiment series design studying how agent cognition emerges from environmental constraints.

---

### Mar 10 — Guideline consolidation
**Key commits:** `8b59811`
**What changed:** Removed three standalone guideline files (686 lines total). Guidelines consolidated into paper dossier READMEs.
**Friction:** Standalone guidelines duplicated content already in paper dossiers.

---

### Mar 10 — Outline → manuscript-plan rename
**Key commits:** `83e8942`
**What changed:** Renamed all `outline-*.md` files to `manuscript-plan-*.md` across incubator papers (8 notes papers renamed). Updated template and scripts.

---

### Mar 10 — Paper directory consolidation
**Key commits:** `55135e1`
**What changed:** Collapsed paper directories: each paper now lives in a single directory with `README.md`, `WORKSPACE.md`, and `manuscript-plan.md`.
**Friction:** Papers scattered across nested subdirectories impeded navigation.

---

### Mar 10 — Control plane documentation expansion
**Key commits:** `a64c6bc`, `dad6f94`
**What changed:** Expanded AGENTS.md (+65 lines), Paper/README.md (+67 lines). Created `paper-workspace-template.md` and `paper-manuscript-plan-template.md`. Updated arxiv-acmart template with section stubs.

---

### Mar 10 — ALife conference paper scaffold and initial draft
**Key commits:** `907b7af`, `1780eb5`, `51d721e`
**What changed:** Scaffolded ALife conference paper: `.gitignore`, `Makefile`, `metadata.yaml`, `references.bib`, `figures/README.md`, `WORKSPACE.md` (64 lines). Created all section files with manuscript plan (+211 lines). First full prose draft with actual content across all sections.

---
