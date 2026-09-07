## Period I: Runtime Features & Experiment Sweeps (Feb 25–26)

The repository's runtime engine received major capability expansions — reasoning model control, swarm forking, unlimited execution — while experiment infrastructure matured through systematic sweep design.

### Feb 25 — Thinking budget control
**Key commits:** `b135ae8`, `103b3c2`
**What changed:** Added `QUINE_THINKING_BUDGET` env var supporting off/low/medium/high values. OpenAI protocol sends `reasoning_effort` + `thinking.type`; Anthropic protocol maps to `budget_tokens` 5000/15000/50000 with thinking/redacted_thinking content blocks.

---

### Feb 25 — Vision and benchmark prompt tuning
**Key commits:** `202f5e4`
**What changed:** Renamed vision parameter `intention` → `instruction`. Added benchmark mission prefix with verification/generalization/replanning guidance.

---

### Feb 25 — Swarm fork redesign
**Key commits:** `1e02177`, `1ee329f`
**What changed:** Major fork tool redesign (+1051 lines). One tool call now spawns N children concurrently in three modes: gather-all, race (first-wins with SIGKILL), fire-and-forget. All-or-nothing slot allocation. Renamed `intents` → `argv` (POSIX naming), changed boolean `wait/race` → single `mode` enum. Added physics-first system prompt for children. LLM behavior test prompts for fork-batch, fork-race, fork-search, swarm-fork.

---

### Feb 25 — Tool call logging rework
**Key commits:** `f17c8fc`
**What changed:** Switched to blacklist approach for `.log` files — show all tool parameters with truncation for long fields. Added `is_error` flag to completion lines.
**Friction:** Previous whitelist approach omitted useful diagnostic information.

---

### Feb 26 — Default unlimited turns
**Key commits:** `2f85dfb`
**What changed:** Changed default `QUINE_MAX_TURNS` to 0 (unlimited). Agent sees "Shell Executions: unlimited" instead of counter. Mortality section omits shell-exhaustion death when unlimited.

---

### Feb 26 — Experiment data corrections and sweep design
**Key commits:** `3dfe9ec`, `a8679f0`, `6557791`, `47e34ca`
**What changed:** Corrected prior session's incorrect model attribution for runs 007-009 (all used Kimi K2.5, not Claude Sonnet 4.6). Designed run-010: first clean baseline with all prior confounds removed (baseline prompt, no turn limit, no fork changes). Expanded from 43 to 50 tasks after cross-checking against 4.1 sweep found 5 tasks incorrectly classified as impossible.
**Friction:** Inaccurate run attribution and task classification from prior sessions.

---

### Feb 26 — Turn limit discovery
**Key commits:** `26cefa3`
**What changed:** Run-009 results: 13/23 (56.5%). Critical finding: commenting out `QUINE_MAX_TURNS` env var had no effect — binary `config.go` default of 30 still enforced. Environment variable override wasn't working.
**Friction:** Silent configuration failure where the env var appeared to be set but was overridden by compiled default.

---

### Feb 26 — Evolution 11–13 design documents
**Key commits:** `4bfe0b4`, `9c88fb1`, `8f1f330`, `191cacb`
**What changed:** Created `11_subjective_reality.md`, `12_process_as_fs.md`, `13_schrodingers_process.md` in `Artifacts/evolution/`. Expanded with Plan 9 "everything is a file" philosophy. Created `08_vision.md`, `14_autonomy.md` (agent autonomy design). Added behavior test prompts: logic-bomb, sandbox-emergent, sandbox.

---

### Feb 26 — Experiment data bulk import
**Key commits:** `7b87f14`, `090a6ae`, `09176ef`
**What changed:** Bulk import of all experiment run data from runs 006–010 (+388K lines). Includes agent traces, config files, trial logs, verifier outputs, CTRF reports, reward scores.

---

### Feb 26 — Evolution 11 task decomposition
**Key commits:** `31a0b95`, `3ff53b8`
**What changed:** Decomposed Evolution 11 (Subjective Reality) into three tracked subtasks. Task hierarchy reorganization in issue tracker.

---
