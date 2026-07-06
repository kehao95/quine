# Quine

> **A POSIX-native runtime for large language model agents**

Quine realizes an LLM agent as a native operating-system process. An agent has a
PID, standard streams, a filesystem, resource limits, and lifecycle operations
such as `fork`, `spawn`, `exec`, and `exit`. The runtime uses Unix process
semantics as the substrate for agent coordination instead of rebuilding those
semantics inside an application framework.

```sh
# Pipe material into a mission.
echo "What is 2+2?" | quine "answer the question"

# Run a mission directly.
quine "write a haiku about recursion"
```

> **Note.** This is a trimmed public build: the runtime's self-source /
> self-reproduction feature (embedding and re-emitting its own repository) is
> compiled out, so the public binary carries no embedded source. The full
> mechanism is described in the [systems paper](./public/quine-arxiv/).

## At A Glance

Quine is both a runtime and a small public portfolio.

- The **runtime** lives at the repository root — [`cmd/`](./cmd/),
  [`internal/`](./internal/), [`habitat/`](./habitat/).
- The **published papers** and their experiments live under
  [`public/`](./public/README.md).

## Published Work

The full portfolio — papers and essays, published/accepted work only, newest
first. Date is the earliest of submitted/accepted/published (whichever is
recorded). Full descriptions, PDFs, and experiment data live under
[`public/`](./public/README.md).

| Date | Work | Venue | Status |
|------|------|-------|--------|
| 2026-07-03 | [A Tumor in the Repository (essay)](./public/tumor-in-the-repository-essay/) | Blog essay (Substack) | published |
| 2026-07-02 | [Structure Grows](./public/structure-grows/) | Blog essay (Substack) | published |
| 2026-07-02 | [A Tumor in the Repository](./public/computational-neoplasm/) | ALIFE 2026 LBA | accepted |
| 2026-06-26 | [Structural Elicitation in a Frozen POSIX Agent](./public/structural-elicitation/) | ALIFE 2026 LBA | accepted |
| 2026-06-18 | [Facultative Self-Reproduction in Quine](./public/facultative-self-reproduction/) | ALIFE 2026 LBA | accepted |
| 2026-06-14 | [From Simulated Worlds to Process Habitats](./public/agentic-alife-workshop/) | ALIFE 2026 workshop | accepted |
| 2026-04-12 | [Coordination Under Existential Unawareness](./public/minimal-perceptual-prerequisites/) | ALIFE 2026 | accepted |
| 2026-03-27 | [The Autopoietic Repository](./public/the-autopoietic-repository/) | Blog essay (Substack) | published |
| 2026-03-17 | [Harness Engineering Is Architectural Amnesia](./public/harness-engineering-is-architectural-amnesia/) | Blog essay (Substack) | published |
| 2026-03-15 | [Why Terminal Bench 2 Is Broken, and Why I Still Love It](./public/terminal-bench-2-love-letter/) | Blog essay (Substack) | published |
| 2026-03-08 | [Quine: Realizing LLM Agents as Native POSIX Processes](./public/quine-arxiv/) | arXiv:2603.18030 | published |

## Documentation

| If you want to... | Start with |
|-------------------|------------|
| build and run Quine | [Install](#install) |
| understand the runtime | [Core Idea](#core-idea), [Architecture Snapshot](#architecture-snapshot) |
| read the papers | [`public/README.md`](./public/README.md) |
| reproduce an experiment | the paper's `REPRODUCE.md` / `experiments/` under [`public/`](./public/README.md) |
| restore experiment tapes | [Artifacts](#artifacts) |

## Install

Install the command:

```bash
go install github.com/kehao95/quine/cmd/quine@latest
```

Or build from a checkout:

```bash
git clone https://github.com/kehao95/quine.git
cd quine
go build -o ./quine ./cmd/quine/
```

Configure a model endpoint, then run:

```bash
cp .env.example .env
# Edit .env with your provider credentials, then:
source .env
echo "What is 2+2?" | ./quine "answer the question"
```

## Core Idea

Quine's first-order object is the POSIX-native LLM agent runtime.

| Claim | Meaning | Why it matters |
|-------|---------|----------------|
| **Agent = Process** | The agent is a native OS process, not a simulated persona inside an app framework | Moves the substrate from conversation management to a mature computational ontology |
| **POSIX as Physics** | Permissions, signals, pipes, files, and resource limits are the environment's laws | Safety and coordination become environmental constraints, not prompt-only intentions |
| **Context as Entropy** | Tokens are volatile working memory, not durable knowledge | Makes externalization, memory structure, and `fork` principled cognitive operations |
| **Environment as Cultural Memory** | Durable traces can become local cultural priors, not just stored facts | Opens an ALife route from environmental residue to behavior replication without explicit goals or rewards |
| **Selection over Instruction** | Robust behavior emerges because bad strategies fail under scarcity | Reframes agent design from prescribing behavior to shaping a survivable environment |
| **Behavior is the Acceptance Test** | A feature is incomplete until a model can discover and use it from prompt plus schema | Treats legibility and emergence as first-class systems properties |

Quine matters less as another tool wrapper and more as a runtime cut: an agent
is a process living under operating-system physics. That shift changes what
counts as memory, communication, safety, failure, and intelligence.

## Architecture Snapshot

| Agent primitive | POSIX mechanism | Purpose |
|-----------------|-----------------|---------|
| Identity | `session_id`, `run_id`, PID | Durable lineage, physical run, live route |
| Interface | stdin/stdout/stderr | Material input, deliverable output, diagnostics |
| State | filesystem surfaces, retained traces, context, bounded env handoff | Durable state, volatile cognition, explicit inheritance |
| Lifecycle | `fork` / `spawn` / `exec` / `exit` | Context-preserving creation, fresh creation, image replacement, judgment |
| Composition | pipes, job control, process groups | Multi-agent coordination |

The host side is deterministic Go code that exposes OS-like constraints. The
guest side is a probabilistic model operating through that process surface.
Tools map model calls to filesystem, process, and lifecycle operations.

## Repository Map

| Path | Role |
|------|------|
| [`cmd/quine/`](./cmd/quine/) | CLI entrypoint for the runtime |
| [`cmd/world/`](./cmd/world/) | Habitat runtime used by environment-control studies |
| [`internal/`](./internal/) | Runtime substrate: config, LLM protocols, tools, tape, lifecycle, and workspace behavior |
| [`habitat/`](./habitat/) | Habitat world package used by `cmd/world` |
| [`public/`](./public/) | Published papers and their experiments — the portfolio |

## Artifacts

Some experiment outputs are represented by DVC pointer files rather than stored
directly in Git. Public artifact restore is manifest-scoped: the exact published
paths are listed in [`.dvc/public-manifest.txt`](./.dvc/public-manifest.txt) and
restored with:

```bash
./scripts/pull-public-dvc-manifest.sh
```

A plain `dvc pull` is not the public restore contract. Each experimental paper
under `public/` reads completely without the tapes; the tapes are the raw
evidence behind its result tables.

## Follow

- X: [@kehao95](https://x.com/kehao95)
- Substack: [kehao95.substack.com](https://kehao95.substack.com/)

## Citation

Quine has two citation anchors. For the runtime, implementation, or
agent-as-process design point, cite the systems paper. If you discuss Quine as a
process habitat for generative artificial life, also cite the ALife manuscript.

Runtime / systems:

- Hao Ke. [Quine: Realizing LLM Agents as Native POSIX Processes](https://arxiv.org/abs/2603.18030).
  arXiv:2603.18030, 2026. DOI: [10.48550/arXiv.2603.18030](https://doi.org/10.48550/arXiv.2603.18030)

Generative ALife / coordination:

- Hao Ke and Jingyun Wu. Coordination Under Existential Unawareness: Information
  Ablation and Closure Thresholds in LLM Multi-Agent Systems. ALIFE 2026.

See [`CITATION.cff`](./CITATION.cff) for machine-readable citation metadata.

## License

Quine is Free Software released under the GPLv2, the same license family as the
Linux kernel. See [`LICENSE`](./LICENSE) for details.
