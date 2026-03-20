# Quine

> **A POSIX-Native Runtime for Large Language Model Agents**

Quine is a runtime that realizes LLM agents as native operating system processes. Each agent runs as a standard POSIX process with PID, standard streams (stdin/stdout/stderr), and lifecycle management through `fork`/`exec`/`exit`. This design inherits process isolation, composition, and resource control from the OS rather than reimplementing them at the application layer.

```sh
# Basic usage
$ echo "What is 2+2?" | ./quine "answer the question"
4

# Self-reproduction (computational autopoiesis)
$ ./quine "Output a binary that is an implementation of yourself." > q
$ chmod +x q
$ ./q "say hello to the world"
Hello, World!
```

## Follow

Follow ongoing Quine updates, writing, and releases:

- X: [@kehao95](https://x.com/kehao95)
- Substack: [kehao95.substack.com](https://kehao95.substack.com/)

## At A Glance

Quine is both:

- a runnable POSIX-native agent runtime
- a testable engineering artifact
- a systems project built around native POSIX process semantics

The project matters less for "how to wrap a model with tools" than for a stronger claim:

> **The right abstraction for an agent is not a chat session or orchestration graph, but a process living under the laws of an operating system.**

That shift changes what counts as memory, communication, safety, failure, and even intelligence itself.

## What You Can Do Here

- build and run the runtime from [`cmd/quine/`](./cmd/quine/)
- inspect the core implementation under [`internal/`](./internal/)
- run the validation stack under [`tests/`](./tests/)

## Highlights

- standard-stream interface: mission via `argv`, material via `stdin`,
  deliverables via `stdout`, diagnostics via `stderr`
- process-native composition: agents chain with pipes and delegate via child
  processes instead of bespoke orchestration graphs
- replayable runtime state: tape logs and filesystem-backed artifacts make runs
  inspectable after the fact
- layered validation: deterministic Go tests, live runtime checks, and
  model-behavior scenarios live in the same repo

## Quick Start

Install:

```bash
go install github.com/kehao95/quine/cmd/quine@latest
```

Configure:

```bash
cp .env.example .env
# edit .env
source .env
```

Run:

```bash
quine "Write a haiku about recursion"
echo "What is 2+2?" | quine "Answer the question"
```

For the full setup, see [QUICKSTART.md](./QUICKSTART.md).

## Citation

If you use Quine in research, please cite the systems paper:

- [Quine: Realizing LLM Agents as Native POSIX Processes](https://arxiv.org/abs/2603.18030)
- DOI: [10.48550/arXiv.2603.18030](https://doi.org/10.48550/arXiv.2603.18030)

## Repository Map

| Path | Role in the project |
|------|---------------------|
| `cmd/quine/` | CLI entrypoint for the runtime |
| `internal/` | Runtime substrate: config, LLM protocols, tools, tape, lifecycle |
| `tests/` | Unified 4-layer validation control plane |
| `tests/runtime/` | Live runtime contract tests against the real binary |
| `tests/model/` | Instructional and emergent model-layer acceptance tests |
| `experiments/` | Experiment designs, run records, and empirical exploration |

If you only want to use Quine rather than navigate the broader repository,
start with [QUICKSTART.md](./QUICKSTART.md) and [TESTING.md](./TESTING.md).

## Why Quine

Most agent stacks build a runtime in user space. Quine takes the opposite bet:
the operating system already provides the core primitives an agent runtime
needs.

| Agent concern | POSIX primitive |
|--------------|-----------------|
| Agent instance | Process |
| Communication | `stdin` / `stdout` / `stderr` / pipes |
| Durable state | Filesystem |
| Judgment | Exit code |
| Delegation | Child process |
| Scheduling | OS scheduler |

The result is a runtime that composes with normal Unix tools, keeps failure
surfaces legible, and stays close to substrate-level behavior instead of hiding
it behind framework abstractions.

## Start Here

- [QUICKSTART.md](./QUICKSTART.md) for installation, configuration, and first runs
- [TESTING.md](./TESTING.md) for the validation stack and runtime/model checks
- [tests/README.md](./tests/README.md) for the test layout and layer map
- [DEVELOPMENT.md](./DEVELOPMENT.md) for development workflow and control-plane rules
- [experiments/p3-mrcr/3.1-needle-retrieval/README.md](./experiments/p3-mrcr/3.1-needle-retrieval/README.md) for a representative experiment walkthrough

## License

<a href="https://github.com/torvalds/linux/blob/master/COPYING">
  <img src="https://img.shields.io/badge/License-GPLv2-blue.svg" alt="License: GPLv2">
</a>

**Quine is Free Software.**

It is released under the **GPLv2**, the same license as the Linux Kernel.
I chose this license to assert that the "physics" of AI agents—like the physics of the OS—must remain common infrastructure.

See [LICENSE](./LICENSE) for details.
