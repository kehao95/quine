## Period II: Design Documents & ALife Paper Sprint (Feb 27 – Mar 1)

The Jobs-as-Resources design crystallized through three iterations, while the first paper track launched with a rapid writing sprint targeting ALife 2026 and a new experiment phase was designed.

### Feb 27–28 — Jobs-as-resources design
**Key commits:** `2b3d565`, `475bc6f`, `184a59d`
**What changed:** Created `15_jobs_as_resources.md` (+285 lines) — treating background processes as first-class filesystem resources. Three iterations: initial design → simplified (removed backpressure, unified file semantics, job ID = PID) → added interactive screen paradigm with dual-view (`screen.txt` + `screen.png`), macro-enabled input syntax, adaptive viewport via winsize file.

---

### Mar 1 — Artifacts/ → Paper/ reorganization
**Key commits:** `4e1d27a`
**What changed:** Merged `Artifacts/` content into `Paper/core/` (canonical definitions) and `Paper/shared/` (references, landscape, graphs). Renamed `Quine-Alife/` → `alife/`, `Quine-System/` → `system/`. (+2891 lines across 77 files)
**Friction:** DRY violation — core definitions duplicated between `Artifacts/` and `Paper/` directories.

---

### Mar 1 — ALife paper structure and initial draft
**Key commits:** `7e925ab`, `7f46819`, `168dd90`, `9c81114`
**What changed:** Created `paper.prompt.md` (11,860 lines) as giant context document for paper generation. Created ALife paper section files: background, anatomy, experiments, discussion, conclusion, appendices. Created `Paper/TRACK/alife/README.md`. Added `teleology-alife.md` reference document (249 lines). Strategic ambiguity framing for the paper's positioning.

---

### Mar 1 — Two-layer physics distinction
**Key commits:** `01c0efe`
**What changed:** Renamed §3.1 from "POSIX as Physics" to "Two Layers of Physics." Distinguished kernel constraints (SIGKILL, OOM) from runtime ecology (turns, context, carrying capacity).

---

### Mar 1 — Phase 5 FUSE tarpit experiment design
**Key commits:** `3c8fd13`
**What changed:** Designed `experiments/active/early-existence/exst-12-fuse-tarpit/` to test whether meta-cognition / representation change emerges under physical constraints. FUSE filesystem creates deliberate I/O delays. Includes `analyze.sh`, `run.sh`, `create_tarpit_fuse.py`. Created `experiments/active/early-existence/README.md`.

---

### Mar 1 — ALife paper discussion expansion
**Key commits:** `6907d0d`
**What changed:** Strengthened Discussion with evaluation critique. Expanded future directions: endogenous evolution, cross-boundary ecologies, niche construction. Created `07-references.md` (105 lines).

---
