## Period XII: ALIFE Submission, Self-Source & Structural Reflection (Apr 12–25)

Two ALIFE papers were submitted, the self-source process surface was implemented (enabling the agent to inspect and rebuild from its own code), FUSE-based public control surfaces landed, the context-first cognition reframe triggered an alignment cascade, and the repository began reflecting on its own structural evolution.

### Apr 12 — Quine-lifelike manuscript rewrite
**Key commits:** `6eee708`, `869eab8`, `a5546d1`, `5437ff3`
**What changed:** Major rewrite: sections reorganized from 7 to 5 files. Added BibTeX references and three new reference dossiers. Tightened ontology boundaries (axiom-02, evolution-14, evolution-16). Final terminology tightening.

---

### Apr 12 — QUINE_PROMPT_IMPL_DETAILS gate
**Key commits:** `0cc0f98`
**What changed:** Added env var controlling whether detailed fd lifecycle descriptions appear in system prompt. Default off to avoid polluting non-autopoiesis experiments.

---

### Apr 12 — Experiment evidence retention
**Key commits:** `b550862`, `720c067`, `d4f0a86`, `e16fae4`
**What changed:** Retained large batch of experiment run artifacts from quine-lifelike redesign wave (+203K lines). Pruned invalid pure-zero-perception runs (−13K lines). Retained cognition-lifecycle and guardian-horizon runs.

---

### Apr 12 — Minimal-perceptual-prerequisites refresh
**Key commits:** `0ef2d13`
**What changed:** Updated degradation gradient figures, revised experiment data YAML, tightened all section prose and references.

---

### Apr 12 — ALIFE paper track split
**Key commits:** `b1537e3`
**What changed:** Created journal-track paper `quine-process-habitat` as separate lane from conference-length `quine-lifelike` with full dossier (manuscript-plan, metadata, plans, references, scripts, sections).

---

### Apr 12 — Submission metadata standardization
**Key commits:** `c02237a`
**What changed:** Standardized build output filenames and submission metadata across all three TRACK/alife papers.

---

### Apr 12 — Autopoiesis frontier paper bootstrap
**Key commits:** `ba09a35`
**What changed:** Bootstrapped `quine-cognitive-autopoiesis` dossier with 299-line manuscript draft, initial manuscript plan, and workspace. Updated experiment paper-feed annotations for 6.1, 6.2, 6.26, 6.3, 6.39, 6.53.

---

### Apr 12 — Paper control-plane root creation
**Key commits:** `90f54a1`
**What changed:** Created Paper/README.md as the research program's control-plane root (+98 lines). Recorded both ALIFE 2026 submissions. Added cross-paper editorial repair rules and paper-workspace-template.

---

### Apr 12 — Final quine-lifelike submission
**Key commits:** `12cec7e`
**What changed:** Submission-grade edits: kernel-enforced distinction, fork/exec/exit lifecycle vocabulary, cross-references, compressed auxiliary paragraphs, substrate-enforced vs tunable comparison. Submitted as fp190s1, 9 pages.

---

### Apr 13 — Post-submission planning
**Key commits:** `ade81ef`, `1ca5c80`
**What changed:** Created `ai-self-reproduction-maintenance` paper dossier in TRACK/system with continuity ladder brainstorm, prior work scan, and theoretical overview. Added Borges sidecar experiment runs and `paper-public.sh` script for sharing paper drafts.

---

### Apr 14 — Idle suspension primitive
**Key commits:** `72408f6`, `ebd7502`, `01a6361`
**What changed:** Implemented explicit idle/standby state in the runtime. Added model evaluations (quiet-standby-pressure prompt), runtime test coverage, updated coverage map. Documented client-signal design decisions and signal class hierarchy. Retained baseline evaluation runs.

---

### Apr 14 — qctl companion client
**Key commits:** `b2d5514`, `7a0616e`, `ce99209`
**What changed:** Created `cmd/qctl/main.go` and `internal/qctl/qctl.go` — thin CLI for interacting with runtime control surfaces (+978 lines). Major rewrite of model evaluation harness. Added callback protocol evaluations.

---

### Apr 14 — Shared ALIFE build infrastructure
**Key commits:** `e52e4cf`
**What changed:** Extracted common arXiv packaging logic into shared template Makefile, reducing duplication.

---

### Apr 15 — Paper program systematization
**Key commits:** `696627`
**What changed:** Created `split-merge-protocol.md`, dossier template, shape-checking script. Standardized multiple paper dossiers to new template. Expanded 6.54 README.

---

### Apr 15 — Autopoiesis evidence absorption
**Key commits:** `0cf950e`, `7f71fae`
**What changed:** Repaired detached fd wording. Scaffolded 6.58 rerun. Absorbed all retained autopoiesis evidence through experiment 6.57 (+257K lines). Major p6 README rewrite.

---

### Apr 15 — Process-surface callback pilot
**Key commits:** `03705e8`, `b9cfbe5`
**What changed:** Added cleanroom pilot evaluation for process-surface callbacks — testing model discovery of peer callback mechanism.

---

### Apr 15 — Process surface refinement
**Key commits:** `b1af657`
**What changed:** Major runtime.go rewrite (+121 lines), qctl expansion (+34 lines), runtime_test expansion (+118 lines).

---

### Apr 16 — Deep theory: two-act manuscript structure
**Key commits:** `0363e9d`, `2538d41`, `9ae8635`, `2648342`, `694773`
**What changed:** Restructured autopoiesis paper around two ontological questions: Act I "What is reproduced?" (fresh synthesis, identity collapse, target-shaping) and Act II "Does reproduction actually continue?" (pressure-elicited renewal, handoff physics). Defined "bounded" along five dimensions. Added explicit spine hierarchy table, split Stream 3. Named identity collapse/climb as conceptual pair.

---

### Apr 16 — 57-experiment paper-facing analysis
**Key commits:** `e94c79e`
**What changed:** Created EXPERIMENT-LOG.md (1083 lines) with six-dimension analysis for all 57 p6 experiments. Each experiment received standardized "Paper-Facing Read" section. Major p6 README restructure.

---

### Apr 16 — Themes layer creation
**Key commits:** `3d967b9`, `5b576d7`, `f90c9f6`, `45f0548`, `7445296`
**What changed:** Created `Paper/themes/` as new structural layer. Moved autopoiesis-renewal from incubator to themes with plans and references. Formalized theme metadata contract with frontmatter checking. Absorbed theory notes from paper plans into shared theme layer. Trimmed duplicate routing.
**Friction:** Cross-track theory was accumulating inside individual paper dossiers, widening them beyond their manuscript boundary. Moved to incubator, but incubator was wrong — this was durable cross-track ownership, not incubation.

---

### Apr 16 — Paper validation refactoring
**Key commits:** `884e74f`, `c29a3bd`, `6852e37`
**What changed:** Refactored paper validation from path-based to label-based discovery. Created `paper_dossier_labels.py`, `list-paper-dossiers.py`, `audit-paper-feed-links.py`, `resolve-paper-dossier.py`. Hardened dossier schema with migration guardrails.

---

### Apr 16 — Paper research workflow surface
**Key commits:** `bb5eaf1`
**What changed:** Created `Paper/research.md` (129 lines) as theory routing surface. Added `autopoiesis-renewal-theory` incubator (307-line README).

---

### Apr 16 — Dead reservoir cleanup
**Key commits:** `72c80a6`
**What changed:** Stripped stale paper content (−4633 lines). Created `Paper/control-plane-topology.md`. Updated paper navigation surfaces.
**Friction:** `quine-process-habitat` still carried a full manuscript tree despite being classified as a delta reservoir — it looked like a live competing paper home.

---

### Apr 16 — Autopoiesis renewal research synthesis
**Key commits:** `bfc5626`, `cf22644`
**What changed:** Four new research synthesis documents: handoff-physics-continuity-carriers, identity-through-handoff, reproduced-interpretive-organization, strongest-counterexamples-missing-regimes. Recentered theme family on bounded organizational closure.

---

### Apr 17 — Preprint convergence
**Key commits:** `15d0163`, `cf33db9`, `4454c0a`
**What changed:** Created preprint-first publication surfaces with paper-b-preprint-cut plan (85 lines). Converged Paper-A manuscript framing: created preprint subdirectory with 388-line manuscript-plan. Absorbed ACM article template, Beer reference intake, LaTeX infrastructure.

---

### Apr 18 — Context-first cognition reframe
**Key commits:** `bb20257`
**What changed:** Major runtime refactor shifting from instruction-heavy to context-first cognition model. Tools and schemas reorganized, test coverage updated. 41 files, +1253/−565 lines.

---

### Apr 18 — Context-first alignment cascade
**Key commits:** `7df1364`, `a2cf0d0`, `9c228d9`, `3a8683f`, `a63f501`
**What changed:** Five alignment commits in 41 minutes propagated context-first reality through: root docs (ROADMAP, STATUS) → Paper/core architecture, primitives, registries → ROADMAP rewrite to single-active-pull structure → Paper program → control-plane-topology.md. Created `world-suspension.md` primitive doc. ROADMAP.md rewritten from topic map to focused pull (+182/−173 lines). Added "Roadmap Pull Contract" to DEVELOPMENT.md.
**Friction:** The deep architectural change immediately made every downstream control-plane document stale.

---

### Apr 18 — FUSE public control surface
**Key commits:** `a98b6ed`, `f47918d`, `3c11cf2`, `4bb3663`
**What changed:** Added FUSE-based public control surface: `control.go` (91 lines), `public_surface.go` (68 lines), `public_surface_fuse_linux.go` (221 lines). Major test expansion (+328 runtime tests, +305 integration tests). Extended LLM protocol layer with FUSE-aware logic. Merged feature branch.

---

### Apr 18 — FUSE-first runtime surface design
**Key commits:** `f290732`
**What changed:** Major expansion of process-as-file-system evolution doc (+133 lines) with FUSE design rationale and boundaries.

---

### Apr 18 — Fork child projection repair
**Key commits:** `2a7c3c6`, `46e6d6a`
**What changed:** Fixed how fork children project runtime state (+264 lines in fork.go, +159 test lines). Slimmed fork bootstrap transaction payload.

---

### Apr 19 — Experiment 6.59 cleanroom closure
**Key commits:** `89d6ddd`, `e0f1f8b`, `07044df`, `7e8a495`
**What changed:** Retained first positive result from 6.59 (cleanroom surface isolation) with detailed analysis. Absorbed into all paper control planes. Promoted from "positive read" to "autonomous cleanroom closure."

---

### Apr 19 — Paper compression 9→7→6
**Key commits:** `cd970ff`, `f88df75`
**What changed:** Merged 9 preprint sections to 7 (§4+§5 combined, §7+§8 combined), then compressed to 6 with framework-first progressive disclosure. Drafting index rewritten. Build files aligned.
**Friction:** 9-section structure reflected editorial plan, not evidence arc. Experiment 6.59's closure proved §7 and §8 were one arc.

---

### Apr 19 — Environment control canonicalization
**Key commits:** `093abcf`
**What changed:** Created `.env.example` (+44 lines), `Paper/core/registries/env-controls.md` (+125 lines), expanded QUICKSTART.md (+62 lines). Added runtime doc sync checking.

---

### Apr 19–20 — Self-source implementation
**Key commits:** `4035367`, `5aa9b16`, `feac92f`, `8fdf158`, `0cb6e98`, `d934d94`, `addb188`, `fe14d40`
**What changed:** Implemented self-source process surface: `self_source.go` (124 lines) in runtime, root-level `selfsource.go` (33 lines) for source embedding. Added evaluation coverage for self-source inspection and self-rebuild exec handoff. Added incarnation-scoped runtime surfaces. Reorganized `cmd/` into scoped directories, then rolled back to root-canonical layout when 6.60 self-rebuild failed. Introduced `QUINE_SELF_REENTRY_TARGET` as explicit reentry mechanism. Created `selfsource_test.go` (73 lines).
**Friction:** Scoped `cmd/` layout (agent/quine/, operator/qctl/, habitat/world/) broke self-rebuild path — `go build ./cmd/quine` doesn't work from the projected surface. Rollback discovered same-day via cleanroom experiment.

---

### Apr 20 — Experiment 6.60 self-source cleanroom
**Key commits:** `d1e4abd`, `0246129`
**What changed:** Created experiment 6.60 (cleanroom self-source): successor inspects and rebuilds from its own source. Retained first sample. Major p6 README rewrite (+155 lines).

---

### Apr 21 — Theme growth and safety interfaces
**Key commits:** `caaf3cc`
**What changed:** Documented theme growth rules and safety interfaces: creation criteria, promotion from incubator, boundaries with paper dossiers.

---

### Apr 21 — P7-OEE evidence retention
**Key commits:** `f2d9f7c`
**What changed:** Retained all p7-oee experiment runs: 7.6 exit-restored-continuation, 7.7, 7.8 (+171K lines across 511 files).

---

### Apr 21–22 — Paper-B freeze and decomposition
**Key commits:** `4c2a94a`, `f9e7302`
**What changed:** Froze tri-axial evaluation schema (Continuity Target × Generative Enclosure × Handoff Regime) as explicit reference layer (+8433 lines of crystallized vocabulary). Then split into mini-route (publishable quickly) and full continuity framework. Added strongest non-quine rivals analysis and artificial succession theory note.
**Friction:** Paper-B's conceptual frame oscillated between scalar "replication" labels and actual physics for weeks without convergence.

---

### Apr 24 — macOS reentry hardening
**Key commits:** `fd47f2d`, `bd03be9`
**What changed:** Hardened macOS reentry environment gates: expanded config tests (+137 lines), control.go refinements (+52 lines). Gated shell interactive prompt surfaces with validation. Moved `active-go-packages.sh` from scripts/ to tests/.

---

### Apr 24 — External project skills evaluation
**Key commits:** `add069d`
**What changed:** Added evaluation coverage for external project skills catalog with new prompts and retained test runs.

---

### Apr 25 — EVOLUTION.md creation
**Key commits:** `470998c`
**What changed:** Created EVOLUTION.md (198 lines) as the structural evolution journal with 10 retrospective arcs. Added structural friction journal contract to DEVELOPMENT.md (+59 lines) and perception guidance to AGENTS.md. Bootstrapped "structure-grows" blog paper dossier.

---

### Apr 25 — Fragments.d prompt surfaces
**Key commits:** `1b30e9a`
**What changed:** Added fragments.d prompt surfaces and AGENTS evaluation coverage. Major test harness expansion (run.sh +494 lines, runtime run.sh +162 lines). Retained evaluation run batch.

---
