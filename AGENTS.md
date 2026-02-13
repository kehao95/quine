# Agent Instructions

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

---

## Experiment Protocol

All experiments MUST be **well-bounded** and **self-contained**. Each experiment produces a complete, reproducible record. Each run gets a unique **runid** directory.

### Directory Structure

```
experiments/
├── README.md                        # Master index with results summary
├── <phase>-<theme>/                 # Phase directory (e.g. p1-physics/)
│   ├── README.md                    # Phase overview & experiment index
│   ├── <id>-<experiment>/           # Individual experiment (e.g. 1.1-insurance/)
│   │   ├── README.md               # Experiment design & results
│   │   ├── prompt.md               # Prompt template (single source of truth)
│   │   ├── run.sh                  # Parameterized run script
│   │   ├── setup/                  # Environment setup scripts (optional)
│   │   │   └── *.sh               # Pre-experiment environment preparation
│   │   └── runs/                   # Each run isolated by runid
│   │       └── <YYYYMMDD-HHMMSS>-<model-short>/
│   │           ├── workspace/      # Agent working directory (cd here to run)
│   │           │   └── ...        # All agent-generated files/artifacts
│   │           ├── quine/          # Quine runtime data (QUINE_DATA_DIR)
│   │           │   ├── *.jsonl    # Session tape files (flat)
│   │           │   ├── *.log      # Runtime logs (flat)
│   │           │   └── locks/     # Semaphore locks
│   │           └── meta/           # Experiment metadata
│   │               ├── prompt-used.md  # Prompt snapshot for this run
│   │               ├── stdout.txt      # Standard output
│   │               └── stderr.txt      # Standard error
│   └── ...
└── ...
```

**Key conventions:**
- `experiments/<phase>-<theme>/<id>-<experiment>/` — experiment definition (checked in)
- `runs/<runid>/` — ephemeral run data (may be large, selectively committed)
- **runid format:** `YYYYMMDD-HHMMSS-<model-short>` (e.g. `20260209-143022-opus`)

### Run Execution

Experiments execute directly on the host. The `run.sh` script sets up the run directory and launches quine:

```bash
# run.sh handles all of this automatically:
RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="./runs/${RUNID}"
mkdir -p "${RUN_DIR}/workspace" "${RUN_DIR}/quine" "${RUN_DIR}/meta"

# Snapshot the prompt
cp prompt.md "${RUN_DIR}/meta/prompt-used.md"

# Execute quine from workspace directory
cd "${RUN_DIR}/workspace"
QUINE_DATA_DIR="../quine" \
QUINE_MODEL_ID="${MODEL}" \
  quine < "../meta/prompt-used.md" \
  > "../meta/stdout.txt" \
  2> "../meta/stderr.txt"
```

**Directory purposes:**
- `workspace/` — Agent's working directory; all generated files appear here
- `quine/` — Quine runtime data (tapes, logs, locks); set via `QUINE_DATA_DIR`
- `meta/` — Experiment scaffolding (prompt snapshot, stdout/stderr captures)

### Experiment Workflow

Experiments follow a **Design-First** approach: complete the design document before writing any execution code.

1. **Design**: Write `README.md` with:
   - Core hypothesis and predictions
   - Independent / dependent variables
   - Control group vs. experimental group
   - Expected results (even sketch a predicted curve)
   - Risks & mitigations
   - Open questions that need human decision
   → **Checkpoint: Get human review before proceeding.**
2. **Implement**: Build the experiment tooling:
   - `prompt.md` — the prompt template (single source of truth)
   - `run.sh` — parameterized runner (accepts `MODEL`, `EXTRA_ARGS`, etc.)
   - `setup/` — environment preparation scripts (if needed)
   - Scorer / analyzer scripts (if needed)
3. **Execute**: Run via `./run.sh` — all outputs auto-captured in `runs/<runid>/`:
   - Tape files → `runs/<runid>/quine/${SESSION_ID}.jsonl`
   - Runtime logs → `runs/<runid>/quine/${SESSION_ID}.log`
   - Prompt snapshot → `runs/<runid>/meta/prompt-used.md`
   - Generated files → `runs/<runid>/workspace/`
4. **Analyze**: Update `README.md` with actual results:
   - **Model used** (e.g., `claude-sonnet-4-20250514`, `gpt-4o`)
   - **runid** reference (e.g., `20260209-143022-opus`)
   - Execution summary (tokens, turns, duration)
   - Key observations and findings
   - Pass/Fail determination with reasoning
   - Comparison tables across runs
5. **Commit**: Experiment definition always committed; run data selectively committed

### Principles

- **Reproducibility**: `run.sh` + `prompt.md` = fully reproducible
- **Isolation**: Each runid gets its own workspace directory
- **Traceability**: Every artifact lives under a unique runid, linked to its tape files
- **Repeatability**: Same experiment, different models/params → new runid, same script

### Run Script Conventions

Every `run.sh` MUST follow these patterns:

```bash
#!/bin/bash
# Usage: ./run.sh [MODEL] [EXTRA_ARGS...]
# Example: ./run.sh claude-opus-4-6
#          ./run.sh gpt-4o --max-depth 3

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MODEL="${1:-claude-sonnet-4-20250514}"
MODEL_SHORT="${MODEL##*-}"  # extract last segment: "opus", "4o", etc.

RUNID="$(date +%Y%m%d-%H%M%S)-${MODEL_SHORT}"
RUN_DIR="${SCRIPT_DIR}/runs/${RUNID}"
mkdir -p "${RUN_DIR}/workspace" "${RUN_DIR}/quine" "${RUN_DIR}/meta"

# Snapshot the prompt used for this run
cp "${SCRIPT_DIR}/prompt.md" "${RUN_DIR}/meta/prompt-used.md"

echo "Run ID: ${RUNID}"
echo "Run dir: ${RUN_DIR}"
echo "Model: ${MODEL}"

# Execute from workspace directory
cd "${RUN_DIR}/workspace"
QUINE_DATA_DIR="../quine" \
QUINE_MODEL_ID="${MODEL}" \
  quine < "../meta/prompt-used.md" \
  > "../meta/stdout.txt" \
  2> "../meta/stderr.txt"
```

### Hard-Won Lessons (from 3.1-js-quine v1–v6)

1. **Always `timeout` agent-generated code** — Use `timeout 5 node "$F"` etc. to prevent runaway processes
2. **Log workspace path early** — `RUN_DIR` is logged at the start so you can find artifacts if the run fails
3. **Tape file attribution** — `QUINE_DATA_DIR` ensures all tapes land in `runs/<runid>/quine/tapes/`

---

## Dual Repository Strategy (Public + Private Lab)

This project uses a **Shadow Branch + Dual Remotes** strategy to keep open-source code clean while preserving private experiment data.

### Repository Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Local Repository                                           │
│                                                             │
│  main ──────────────────────► public/main (kehao95/quine)   │
│    │       (pure code)              (public repo)           │
│    │                                                        │
│    └──── merge ────► lab ───► origin/lab (kehao95/quine-lab)│
│                      (code + experiments)  (private repo)   │
└─────────────────────────────────────────────────────────────┘
```

### Remotes

| Remote | URL | Purpose |
|--------|-----|---------|
| `public` | `git@github.com:kehao95/quine.git` | Open-source release |
| `origin` | `git@github.com:kehao95/quine-lab.git` | Private lab with experiments |

### Branches

| Branch | Contains | Push To |
|--------|----------|---------|
| `main` | Source code, docs, Artifacts/ | `public/main` |
| `lab` | Everything in main + experiments/, Paper/, .beads/, AGENTS.md | `origin/lab` |

### Daily Workflow

**Scenario A: Modify code and test**
```bash
# Stay on main branch
git checkout main
# Edit code, run tests
./quine < prompt.md
# Commit and publish
git add . && git commit -m "Fix logic"
git push public main
```

**Scenario B: Save experiment results worth keeping**
```bash
# Switch to lab and sync code
git checkout lab
git merge main  # Keep lab up-to-date with main!

# Commit experiment data
git add -f experiments/p3-mrcr/runs/20260213-*/
git commit -m "Record MRCR experiment results"
git push origin lab

# Return to main
git checkout main
```

### Key Principles

1. **Unidirectional flow**: Code flows `main → lab`, never the reverse
2. **Always merge before committing data**: `git merge main` prevents divergence
3. **Public repo never sees lab**: Only `main` is pushed to `public`
4. **`.gitignore` on main excludes**: `.beads/`, `experiments/`, `Paper/`, `AGENTS.md`

---

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

