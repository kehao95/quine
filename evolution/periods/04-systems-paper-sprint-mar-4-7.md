## Period IV: Systems Paper Sprint (Mar 4–7)

An intense four-day cycle produced the systems paper through ~35 outline iterations, elevating Plan 9 from a future vision to the paper's organizing principle, culminating in peer review and submission.

### Mar 4 — Systems paper infrastructure
**Key commits:** `136a842`
**What changed:** Created `a-priority/` directory with 7 section files, `Architecture.svg` diagram. Three outline tracks: `outline-track-a.md` (240 lines), `outline-track-b.md` (385 lines), `outline.md` (180 lines). Compared MemGPT, AIOS, LangChain, AutoGPT.

---

### Mar 5 — Outline iteration sprint (v6–v38)
**Key commits:** `398d1d2`, `7aad053`, `38b9208`, `8737d69`, `0efed73`, `2999f9e`, `207d680`, `73faedf`, `1946154`, `229d99d`, `d325b11`, `3536d57`, `3acd3a9`
**What changed:** ~30 outline iterations in one day. Key moves: consolidated section files into detailed outlines; added 202 lines of references; created `16_mesh_collaboration.md` (259 lines) and `17_tape_systems.md` (159 lines); switched to first-person singular voice; reframed agent as I/O-bound process; added "Why Now?" section on tool-calling maturity; promoted Plan 9 from subsection to full §6; achieved 4 limitations : 4 Plan 9 layers symmetry; merged limitations into §6 with problem→solution structure.

---

### Mar 5 — Changelog extraction
**Key commits:** `017e336`
**What changed:** Extracted complete version history (v6–v38) into `CHANGELOG-outline-track-a.md`. Outline reduced from 791 to 407 lines.
**Friction:** Version history was bloating the outline document itself.

---

### Mar 5 — Plan 9 as core organizing principle
**Key commits:** `51ec5ce`
**What changed:** Massive outline rewrite (+217/−179 lines) promoting Plan 9 from vision/future-work to the core organizing principle of the systems paper.

---

### Mar 6 — Unified agent namespace and outline split
**Key commits:** `fe276f6`, `ade84a9`
**What changed:** Agent as directory (`/agent/<pid>/`), all Layer 1 concepts mapped. Layers 3-5 use consistent paths. Split systems paper into two tracks: arxiv (738 lines) for full paper and hotos (466 lines) for workshop submission. Created `DEFERRED-IDEAS.md` (179 lines). Deleted unified `outline-track-a.md`.

---

### Mar 6 — FD mapping correction
**Key commits:** `8955cfa`, `74ebd40`
**What changed:** Documented and fixed a file descriptor mapping issue. Created `FD3-CORRECTION.md` (111 lines). Updated `core/implementation.md`.
**Friction:** Incorrect fd mapping documented in paper didn't match implementation.

---

### Mar 6 — Full systems paper written
**Key commits:** `2bd4d7d`, `a438ef9`, `4b46949`, `aad1b29`
**What changed:** All section files created: abstract, introduction, posix-mapping, runtime (153 lines), properties (181 lines), related-work (131 lines), beyond-posix, conclusion. Generated `full_paper.md` (727 lines) and PDF (250KB). Title changed from "Architecture" to "Runtime." Added availability statement and feasibility disclaimer. README restructured for dual-mission (Systems + ALife papers).

---

### Mar 7 — Peer review and submission
**Key commits:** `7cb2432`, `4c4a8fb`, `53bd087`, `92ecebe`, `5f07bc8`, `5fa539d`
**What changed:** Merged beyond-posix into conclusion. Added acmart LaTeX template. Major related work compression (−184 lines). Incorporated peer review feedback across all sections. Created `appendix-a.md` (103 lines). Final submission: PDF renamed to `Realizing_LLM_Agents_as_Native_POSIX_Processes.pdf` (331KB).

---
