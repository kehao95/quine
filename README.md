# Quine

> **A POSIX-native runtime for large language model agents**

Quine treats an agent as a real operating-system process rather than a chat session or orchestration node. Mission comes in through `argv`, material comes in through `stdin`, deliverables leave on `stdout`, failures stay on `stderr`, and delegation becomes ordinary process creation.

```sh
$ echo "What is 2+2?" | quine "Answer the question"
4

$ cat article.md | quine "Summarize the argument" | quine "Turn it into three bullets"
```

The claim behind the project is simple: agent systems do need a runtime, but rebuilding that runtime in user space throws away the substrate that already solves execution, communication, state, supervision, and scheduling.

## Why an OS Runtime

Quine stays close to POSIX primitives on purpose. It prefers inheriting the operating system's contracts over inventing another application-layer harness.

| Agent concern | POSIX primitive |
|--------------|-----------------|
| Agent instance | Process |
| Working context | Process memory |
| Durable state | Filesystem |
| Communication | `stdin` / `stdout` / `stderr` / pipes |
| Judgment | Exit code |
| Delegation | Child process |
| Scheduling | OS scheduler |

This is not UNIX nostalgia. It is an engineering leverage argument: use the layer that already knows how to isolate, compose, observe, and clean up independent computational units.

## What Exists Today

- A Go runtime and CLI in [`cmd/quine/`](./cmd/quine/)
- Native stream-oriented I/O and shell composition
- Filesystem-backed state plus tape logs for replay and inspection
- Recursive delegation through child agents with bounded runtime limits
- Deterministic Go tests plus live API-backed runtime and behavior checks

## Start Here

| Goal | Read |
|------|------|
| Install and run Quine | [QUICKSTART.md](./QUICKSTART.md) |
| Understand the public test surfaces | [TESTING.md](./TESTING.md) |
| Read the systems paper | [paper/Quine_Realizing_LLM_Agents_as_Native_POSIX_Processes.pdf](./paper/Quine_Realizing_LLM_Agents_as_Native_POSIX_Processes.pdf) |
| Read the artificial-life paper | [paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf](./paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf) |
| Inspect runtime internals | [`internal/`](./internal/) |

## Public Papers
ArXiv preprints(on hold for now)
- [Quine: Realizing LLM Agents as Native POSIX Processes](./paper/Quine_Realizing_LLM_Agents_as_Native_POSIX_Processes.pdf)
- [Quine: POSIX as Physics for Emergent Artificial Life](./paper/Quine_POSIX_as_Physics_for_Emergent_Artificial_Life.pdf)

## Branch Surface

`main` is the curated public branch. It keeps runnable code, public manuscripts, and enough documentation to evaluate the runtime without mirroring the full internal research workspace.

## License

<a href="https://github.com/torvalds/linux/blob/master/COPYING">
  <img src="https://img.shields.io/badge/License-GPLv2-blue.svg" alt="License: GPLv2">
</a>

**Quine is Free Software.**

It is released under the **GPLv2**, the same license as the Linux Kernel.
See [LICENSE](./LICENSE) for details.
