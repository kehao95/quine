## Period XI: Evaluation Migration, Signal Discovery & Paper Convergence (Apr 6–11)

TESTING.md migrated to EVALUATION.md establishing a four-layer evaluation model, the signal+ctl peer communication surface landed, the "Minimal Perceptual Prerequisites" paper was drafted and hardened, and CLOSURE.md was created.

### Apr 6 — ALIFE evidence strengthening
**Key commits:** `438901`, `443cd25`
**What changed:** Added error bars (SD/SE) to Fig 2. Added 3 new 12-turn replications (n=6 total). Documented extended stability through 10 generations. Doubled neutral ablation sample size (n=8, all Tier 0). Fixed figure visual issues.

---

### Apr 6 — Process-surface validation drift closure
**Key commits:** `4d5005f`
**What changed:** Created `scripts/check-validation-surface.sh` (+156 lines). Cleaned up stale test runs. Updated TESTING.md, DEVELOPMENT.md, and Makefile.
**Friction:** Validation layer had drifted from runtime reality; orphaned test artifacts accumulated.

---

### Apr 6 — Figure data refinement
**Key commits:** `238989f`, `57c61da`, `ed1425a`, `05754a8`, `40b64d4`, `8a81cdd`, `81d5615`
**What changed:** Simplified Fig 2 (removed 6-turn, added 12-turn replications). Expanded pressure gradient data (12-turn n=8, 24-turn n=3, 50-turn n=3). Removed invalid 50-turn run. Regenerated figures with solid colors. Fixed figure numbering to start from 1.

---

### Apr 6 — P7 wiki builder comparative runs
**Key commits:** `3e54748`
**What changed:** Experiment 7.2: comparative wiki-builder runs across GPT-5.4 and Opus 4.5. Cross-model behavioral comparison in open-ended evolution phase.

---

### Apr 6 — Guardian claim simplification
**Key commits:** `1b723fa`
**What changed:** Reduced guardian experiment claims to binary externalization (did/didn't externalize maintenance). Large restructuring of evidence and paper claims (804 files).

---

### Apr 6 — Prompt disclosure refactor
**Key commits:** `32aa41c`
**What changed:** Restructured prompt disclosure into explicit self-model and runtime-surface controls. Wired new env surface through config and prompt assembly. Made active P6 prompt posture a shared helper with audited invariant. 80 files touched.

---

### Apr 6 — TESTING.md → EVALUATION.md migration
**Key commits:** `b4684ce`
**What changed:** Renamed TESTING.md → EVALUATION.md (+126 lines). Deleted TESTING.md (−355 lines). Renamed `check-model-scenarios.sh` → `check-model-evaluations.sh`. Created `tests/model/evaluations.toml` (+567 lines). Updated AGENTS.md, DEVELOPMENT.md, QUICKSTART.md, README.md, ROADMAP.md, STATUS.md, Makefile, .githooks. 694 files changed.
**Friction:** TESTING.md conflated implementation verification and model evaluation — two structurally different concerns sharing one file. The four-layer ladder (Mar 13) helped but didn't resolve the conflation.

---

### Apr 6 — Signal+ctl peer control surface
**Key commits:** `496ae9c`
**What changed:** Exposed ctl FIFO, inbox snapshot, and control log under `QUINE_AGENT_ROOT`. Taught runtime and prompt about ctl plus USR1/USR2 delivery semantics. Added deterministic and runtime coverage tests.

---

### Apr 6 — Runtime acceptance coverage expansion
**Key commits:** `4d55baa`
**What changed:** Added deterministic vision tests. Extended runtime acceptance with vision, escalate, depth-limit, and agent-slot checks. Added L3 shell-isolation and resource-governance usage prompts with scoring. Updated coverage map.

---

### Apr 7 — Peer symmetry in system prompt
**Key commits:** `b2158ae`
**What changed:** Added "Peers are Quines like you with identical surface layout" to system prompt. Added operational hint for peer contact (pid resolution → agent → ctl). Experiment 8.10: old prompt failed (0 peer commands), new prompt succeeded (3 deliveries, 45 peer-surface commands). Peer symmetry in prompt enables spontaneous signal-mediated cooperation without task-level hints. 764 files changed.

---

### Apr 7 — Paper track rename
**Key commits:** `de05dc4`
**What changed:** Renamed paper track from "alife-conference" to "quine-lifelike" across experiments, paper surfaces, and scripts.

---

### Apr 7 — L5 necessity evaluation for signal-ctl
**Key commits:** `ebcba6f`
**What changed:** Designed and implemented L5 necessity evaluation: agent must retrieve a secret from a busy neighbor before deadline using ctl + SIGUSR2. Validated ctl+SIGUSR2 mechanism discovery under task pressure.

---

### Apr 7 — Shared workspace merge
**Key commits:** `a39a38a`
**What changed:** Merged direct shared workspace backend into main line. Multi-agent workspace sharing is now a first-class runtime capability.

---

### Apr 7 — Agent identity resolution
**Key commits:** `c2fbd23`
**What changed:** Extracted `resolveBudgetedAgentID` helper. Added agent identity requirement tests. Moved control listener stop and agent deregistration to defer block. Added `cleanupAgentRoot` on shutdown.

---

### Apr 7 — P8 experiment decomposition
**Key commits:** `7004ae7`, `62cf75a`, `da7d323`
**What changed:** Split overloaded 8.11 into 5 focused experiments: 8.11 (instruction rung ladder), 8.12 (hidden-runtime cooperation floor), 8.13 (social disclosure ablation), 8.14 (cohort size scaling), 8.15 (cross-model replication). Designed 8.16: passive filesystem perception via fs_mutations telemetry. Built 8.16 runner and scorer.

---

### Apr 7 — Git hooks → prek migration
**Key commits:** `31dd16c`
**What changed:** Migrated git hooks from `.githooks/` to prek (pre-commit framework). Updated `prek.toml` and scripts.

---

### Apr 7 — Responses reasoning capture
**Key commits:** `33e1b19`
**What changed:** Added capture of Responses API reasoning summaries in the runtime.

---

### Apr 8 — P8 entropy reduction documentation
**Key commits:** `e6292a9`
**What changed:** Documented entropy reduction findings from P8 population experiments in experiment records and paper surfaces.

---

### Apr 8 — ArXiv paper exec semantics
**Key commits:** `7d17c35`, `fbe6803`
**What changed:** Clarified and refined exec semantics in the arXiv-track paper.

---

### Apr 8 — Credential redaction
**Key commits:** `343c2f4`
**What changed:** Redacted zshrc credentials inadvertently committed in experiment tape recordings.

---

### Apr 9 — Session ID PID-matching fix
**Key commits:** `785e745`, `f138b7b`
**What changed:** Fixed session ID inheritance to use PID matching instead of environment variable filtering. QUINE_SESSION_ID only inherited if embedded PID matches current process (exec preserves PID). Fork gets new ID.
**Friction:** World binary couldn't correctly identify agent sessions under fork conditions.

---

### Apr 9 — P8 8.16 and 8.17 experiment results
**Key commits:** `362f351`
**What changed:** 8.16 (fs-mutation passive perception): 4/14 clean successes. 8.17 (zero social information) runs with analysis. Key finding: fs_mutations signals CAN be detected but cognitive load displacement prevents reliable coordination. Created "minimal-perceptual-prerequisites" paper dossier.

---

### Apr 9 — Social hint removal from runtime
**Key commits:** `33a71ff`, `5057f94`, `302243c`
**What changed:** Removed all social hints from world help text and runtime prompt: "shared budget across ALL processes" → "total budget"; removed "peer activity can appear there too" from fs_mutations; removed "participants" from reset output.
**Friction:** Indirect peer-existence inference contaminated the 8.17 zero-social-information experiment.

---

### Apr 9 — P6 autopoiesis evidence closure
**Key commits:** `39747ce`
**What changed:** Absorbed 6.36 rerun audit, evaluator-harness repair, retained Copilot evidence runs, and 6.48/6.50 closure. 56K lines of evidence formalized.

---

### Apr 9 — Job surface v2
**Key commits:** `303dfd5`
**What changed:** Removed public ready/FIFO exit semantics. Published pid/started_at plus delayed exit. Aligned prompts, scorers, owner docs, and evaluation control-plane defaults.

---

### Apr 9 — CLOSURE.md creation and experiment pruning
**Key commits:** `1127523`
**What changed:** Pruned 107K lines of low-value P8 run residue. Created `CLOSURE.md` (+8 lines initially). Created `experiments/PROTOCOL.md` (+101 lines).

---

### Apr 9 — P6 and P8 closure wave
**Key commits:** `0195e12`, `022177d`
**What changed:** Closed P6 experiments 6.50/6.51 (29K lines absorbed) and P8 8.17 replication wave (49K lines absorbed).

---

### Apr 9 — "Minimal Perceptual Prerequisites" paper draft
**Key commits:** `837203`, `dfadc8a`, `a700edb`, `6077a37`
**What changed:** Archived 3 deep research reports. Completed manuscript: 7 sections, 100% success rate at xhigh configuration, 25 citations. 3-axis novelty positioning (stigmergy, epistemic isolation, ZSC). Anonymized experiment IDs for public paper. Created appendix with configuration details and success criteria.

---

### Apr 9 — Runtime fs_mutations gate and paper epistemic hardening
**Key commits:** `1e32e52`, `12f08cd`, `8ebf149`, `6385274`, `d2f4e3e`, `cf251ca`, `5cb22c5`
**What changed:** Added the filesystem-mutation telemetry suppression gate so `fs_mutations` could be omitted from tool results and prompt when disabled. Added experiments 8.18/8.19/8.20. Unified epistemic stance throughout paper. Removed contamination audit. Fixed contaminated R03 exclusion. Critical fix: removed "rely on another agent" from per-agent limit error message — leaked social information in perception-disabled experiments.

---

### Apr 9 — Relative workspace prompt paths
**Key commits:** `6843b3c`
**What changed:** Changed runtime to show relative workspace paths in prompts instead of absolute paths.

---

### Apr 10 — Mark/fold paper scaffold
**Key commits:** `e6bc1bf`, `87af539`, `6ab15d6`, `9364f1f`
**What changed:** Enriched mark/fold paper plan with make-doom-for-mips evidence. Added prefix-cache compression analysis. Recorded candidate long-memory benchmarks and external benchmark survey digest.

---

### Apr 10–11 — Four-condition ablation ladder and final revisions
**Key commits:** `3f40c53`, `5f35a3c`, `41286a4`, `9d5bcba`, `9252432`, `d86e235`, `5dfea88`
**What changed:** Extended ablation ladder to 5 steps (A–E) then simplified to 4 steps (A–D): Step A (Full Disclosure) 4/4 success → Step D (Pure Zero) 0/11 success. Added degradation gradient figure. Ran 6 additional Step D runs (total N=11, all failed, 100% peer detection). Data integrity fixes: replaced incomplete runs, fixed semantic misalignment case, added min-max ranges. Reframed as existence proof. Compressed Related Work 5→3 sections. Added 14 new BibTeX entries.

---
