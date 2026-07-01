# Introduction

> "Write programs that do one thing and do it well. Write programs to work together. Write programs to handle text streams, because that is a universal interface." — Doug McIlroy [@McIlroy1978]

## The Problem

The question of how to structure LLM agents is often asked at the framework level.
This paper asks it one layer lower: what if the runtime substrate were not an application framework, but the operating system itself?

The rapid development of LLM agents has led to a proliferation of frameworks designed to manage their lifecycle, memory, and communication.
Systems such as LangChain [@LangChain], AutoGen [@Wu2023], and CrewAI [@CrewAI] have successfully democratized agent development by providing high-level abstractions for these tasks.
However, a structural pattern has emerged: these systems primarily expose agent abstractions at the application layer, even when the underlying execution still relies on OS services.

This pattern adds complexity that grows with system scale.
By managing agents as objects within a single application process, frameworks frequently reconstruct mechanisms the operating system already provides:

- **Fault isolation** is simulated through exception handlers and try-catch blocks, rather than enforced by hardware-backed address space separation.
- **Context switching** between agents requires application-level scheduling logic, rather than delegating to the kernel's mature scheduler.
- **Message passing** is mediated by in-process queues or database-backed channels, rather than kernel-managed pipes with backpressure.

The operating system, having evolved over five decades to solve these problems for classical software, is less commonly treated as a first-class runtime substrate for agents.
Treating the OS as substrate also sharpens a deeper question: if the process is the right first abstraction for agency, where do process semantics stop being expressive enough for cognition?

## My Approach

Recent progress in tool calling and structured output makes it practical to bind natural-language reasoning to a small set of system-level operations.
Instead of introducing another application-layer framework, I present a runtime architecture where the operating system itself serves as the execution substrate for agents.
The design consists of two components:

**Component 1: A Protocol.**
A disciplined mapping from agent concepts to POSIX primitives across four dimensions:

- **Identity** — The agent's unique identifier is its process ID (PID), assigned by the kernel.
- **Interface** — The agent communicates through standard streams (`stdin`/`stdout`/`stderr`) and reports outcomes via exit status.
- **State** — The agent's memory is process memory (cleared on exit), environment variables (inherited by children), and filesystem (persistent).
- **Lifecycle** — The agent spawns children via `fork`, continues itself via `exec`, and terminates via `exit`.

Section 2 elaborates each dimension; Figure 1 illustrates the interface contract.

**Component 2: A Single-Image Runtime.**
A unitary executable that implements this protocol.
When an agent spawns a child, it instantiates the same runtime image with different arguments.
Parent and child share code but diverge in state.
This recursive structure means the runtime and the agent template are the same artifact—there is no separate "agent definition" language or configuration format.
It also means that delegation is not orchestration from outside: a parent agent constructs a child's operational world using the same runtime it itself inhabits.

## Contributions

This paper makes three contributions:

1. **A systems design perspective for LLM agents.** I argue that the operating system can serve as the runtime substrate for agents, and I make this claim concrete through a disciplined mapping from agent abstractions to POSIX primitives across identity, interface, state, and lifecycle (Section 2).

2. **A reference runtime that instantiates this perspective.** I present Quine, a single-image runtime in which agents are realized as native POSIX processes and recursive delegation is implemented by self-instantiation rather than by an external orchestrator (Section 3).

3. **An analysis of both the reach and the limits of the process model for agents.** I show how POSIX process semantics directly yield isolation, composition, and self-renewal through exec, and I argue that this mapping, precisely because it works, also exposes where execution semantics cease to be sufficient for cognition—most immediately in the constitution of an agent's world and in the revisability of its actions over time (Sections 4–5).

## Scope and Non-Claims

This paper presents the architecture and a reference implementation of Quine, not a wholesale replacement for existing frameworks.
I make no claims about superior end-task performance; the contribution is a concrete runtime design point that inherits isolation, composition, and lifecycle control from POSIX rather than rebuilding them.
Its limitations and boundaries are discussed explicitly in Section 5, which identifies where process semantics remain effective and where they begin to leave important aspects of agency unrepresented.
