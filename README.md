# Quine

> **POSIX as Physics for Emergent Artificial Life**

Quine is a runtime that treats the operating system as the fundamental physics for AI agents. By mapping LLM tool calls directly to POSIX system calls (`fork`, `exec`, `pipe`, `exit`), it enables agents to exhibit emergent behaviors through environmental constraints rather than explicit programming.

```sh
# Basic usage
$ ./quine "count the files in this directory" < /dev/null
42

# Computational autopoiesis (self-reproduction)
$ ./quine "Output a binary that is an implementation of yourself." > q
$ chmod +x q
$ ./q "say hello to the world"
Hello, World!
```

## Overview

This repository provides the **core physics engine** of Quine—the foundational implementation described in the paper ["Quine: POSIX as Physics for Emergent Artificial Life"](paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf) (ALife 2026).

**What's included:**
- Core POSIX runtime with `sh`, `fork`, `exec`, `exit` primitives
- Context lifecycle management (entropy accumulation & reset)
- Tape-based observability (JSONL event logs)
- Basic experiment data for three published emergent behaviors

**What this demonstrates:**
- **Metabolic Transcendence** — Agents handling datasets 6× their context window via generational succession
- **Stigmergic Coordination** — Multi-agent coordination through filesystem state without explicit communication
- **Computational Autopoiesis** — Self-reproduction: agents generating functional copies of their own runtime

## Architecture

Quine implements a **Host-Guest** model where:

- **Host (Go)** — Deterministic POSIX environment providing physical constraints
- **Guest (LLM)** — Probabilistic agent operating under resource limits
- **Tools** — Direct mapping from LLM function calls to OS primitives:
  - `sh` — Execute shell commands (basic perception/action)
  - `fork` — Spawn child processes (delegation, parallelism)
  - `exec` — Restart with fresh context (entropy reset)
  - `exit` — Declare fitness and terminate (selection pressure)

**Key invariants:**
- Context is volatile RAM that fills monotonically (tokens → entropy)
- All persistent state must be externalized to the filesystem
- Communication is pipelines (`stdin`/`stdout`), not conversation
- Death is a feature: resource limits create selection pressure

## Experimental Evidence

The `experiments/` directory contains data for the three core behaviors documented in the paper:

| Experiment | Path | Description |
|------------|------|-------------|
| **Needle Retrieval** | `p3-mrcr/3.1-needle-retrieval/` | 266K token dataset (6× context limit) handled via generational succession |
| **Stigmergy** | `p1-borges/2.3-stigmergy/` | 2–6 agent coordination via filesystem state without explicit messages |
| **Binary Quine** | `p3-mrcr/3.2-binary-quine/` | Self-reproduction: agent generates 5.4MB functional copy of itself |

Each experiment directory contains:
- `README.md` — Experimental design and results
- `prompt.md` — Mission specification
- `run.sh` — Reproduction script
- `runs/*/` — Complete execution traces (tape logs, generated artifacts)

For details on the emergent behaviors and their implications, see the [paper](paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf) (Section 4: Emergent Behaviors).

## Getting Started

### Installation

```bash
# Build from source
go build -o quine ./cmd/quine/

# Configure LLM provider (requires API key)
export ANTHROPIC_API_KEY="your-key"
# or export OPENAI_API_KEY="your-key"
```

### Quick Examples

```bash
# Simple task
echo "What is 2+2?" | ./quine "answer the question"

# With timeout and turn limits
QUINE_MAX_TURNS=10 QUINE_SH_TIMEOUT=30 ./quine "analyze this directory" < /dev/null

# Binary input/output
cat input.bin | ./quine -b "process this binary data" > output.bin
```

👉 **[Full Quick Start Guide](./QUICKSTART.md)**

## Reproducing Experiments

Each experiment can be reproduced via its `run.sh` script:

```bash
cd experiments/p3-mrcr/3.2-binary-quine
./run.sh claude-sonnet-4-20250514
```

Results appear in `runs/<timestamp>-<model>/` with complete tape logs and generated artifacts.

## Advanced Experimental Pipelines

The **advanced experimental suite** (including the 1000-file stigmergic test environment, MRCR longitudinal analysis tools, and autopoiesis bootstrap testbed) is maintained in a private lab branch.

**For academic collaborations, reproducing specific emergent behaviors, or access to the full experimental infrastructure:**
- Contact: [your-email@domain.com]
- Include: Institution/affiliation and research objectives

This ensures proper knowledge transfer and prevents fragmented replication attempts that might yield misleading results.

## Documentation

- **[Paper](paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf)** — Full theoretical framework and experimental results (ALife 2026)
- **[Quick Start](./QUICKSTART.md)** — Installation and basic usage
- **[Experiments](./experiments/)** — Reproduction instructions for published results

## Citation

If you use Quine in your research, please cite:

```bibtex
@inproceedings{quine2026,
  title={Quine: POSIX as Physics for Emergent Artificial Life},
  author={[Hao Ke]},
  booktitle={Proceedings of the 2026 Conference on Artificial Life},
  year={2026}
}
```

## License

<a href="https://github.com/torvalds/linux/blob/master/COPYING">
  <img src="https://img.shields.io/badge/License-GPLv2-blue.svg" alt="License: GPLv2">
</a>

Quine is released under the **GPLv2**, the same license as the Linux Kernel. This ensures that the "physics" of AI agents—like the physics of the OS—remains common infrastructure.

See [LICENSE](./LICENSE) for details.
