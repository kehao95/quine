# Manuscript Plan: arXiv Track — Minimal Verifiable Claim

**Target**: arXiv cs.OS / cs.AI (preprint for priority)  
**Length**: ~5 pages  
**Status**: Released reference plan for `arXiv:2603.18030v1` (submitted 2026-03-08)

---

## Final Polish Summary (March 7, 2026)

This plan records the final pre-submission polish that shipped as `arXiv:2603.18030v1`:

1. **Terminology locked**: Five-channel I/O contract (Mission/Material/Deliverable/Diagnostics/Outcome) used consistently throughout
2. **Structure compressed**: From 4 contributions to 3, aligning with 3C framework (Composition/Containment/Continuity)
3. **Section 6 merged into Section 5**: Future directions compressed from detailed five-layer architecture to ~150-word architectural sketch
4. **Tone calibrated**: "Limitations" → "Boundaries"; possibility claims, not normative claims
5. **Plan 9 grounded**: Two citations added (Pike 1993 for namespaces, Pike 1995 for lineage), positioned as conceptual ancestry
6. **Transitions added**: Explicit bridges at Section 4→5, 5.1→5.2, 5.2→5.3
7. **Redundancy removed**: Self-references ("This section...") eliminated; defensive repetition consolidated
8. **Figures updated**: All three SVG/PDF figures regenerated with consistent terminology

The paper corresponding to this plan is now public on arXiv.

---

## Strategic Positioning

**Venue Fit**: arXiv since late 2025 discourages review/position papers without prior peer review. This version is a **self-contained technical contribution** with verifiable claims.

**Central Claim**: "I show that an LLM agent can be realized as a native POSIX process with a complete reasoning-and-acting loop, and analyze the structural properties that emerge from this mapping."

**What This Version Does**:
- Presents the POSIX mapping as a **protocol** (not a paradigm)
- Validates Layer 1 with a working implementation
- Treats future extensions as **boundaries of POSIX**, not limitations (~0.5 page)
- Focuses on **emergent properties** that follow from the design
- Uses restrained, technical language throughout
- References Plan 9 as conceptual lineage for future directions

**Relationship to HotOS Version**: This occupies the **verifiable core**. The HotOS version extends into paradigm-level reframing ("wrong abstraction layer"). They form a continuous research program.

**Key Vocabulary**: maps, realizes, implements, structural consequences, may improve (not: reframes, relocates, paradigm, complete Agent OS, challenge to the field)

---

## Core Thesis

Quine is not an agent framework. It is (1) a protocol that maps LLM agents to POSIX processes, and (2) a **unitary executable** that implements this protocol through **recursive process instantiation**.

---

## Title

**"Quine: A POSIX-Native Runtime for LLM Agents"**

---

## Abstract (~150 words)

> Quine is a runtime architecture and reference implementation for realizing LLM agents as native POSIX processes. The mapping is explicit: identity is PID, interface is standard streams and exit status, state is memory, environment, and filesystem, and lifecycle is fork/exec/exit. A single executable implements this model by recursively spawning fresh instances of itself. This placement makes process isolation, composition, and resource control available through existing OS mechanisms rather than reimplementing them in an application-layer framework. I analyze three categories of emergent properties—composition (recursive structure, universal composability), containment (hardware isolation, kernel enforcement), and continuity (exec-based self-renewal, feedback channels)—and demonstrate operational feasibility on long-context retrieval tasks. I identify four architectural boundaries where the POSIX model proves insufficient for cognitive systems and sketch future directions inspired by Plan 9's approach to making internal structure visible, nameable, and composable.

---

## 1. Introduction (~0.5 page)

### 1.1 The Problem

> "Write programs that do one thing and do it well. Write programs to work together. Write programs to handle text streams, because that is a universal interface." — Doug McIlroy [McIlroy 1978]

The rapid development of Large Language Model (LLM) agents has led to a proliferation of frameworks designed to manage their lifecycle, memory, and communication. Systems like LangChain [LangChain], AutoGen [Wu et al. 2023], and CrewAI [CrewAI] have successfully democratized agent development by providing high-level abstractions for these tasks. However, a structural pattern has emerged: many current systems primarily expose agent abstractions at the application layer, even when the underlying execution still relies on OS services.

This pattern introduces complexity. By managing agents as objects within a single application process, frameworks reconstruct fault isolation, context switching, and message passing in application code—mechanisms the operating system already provides with hardware-enforced guarantees. The operating system, having evolved over five decades to solve these problems for classical software, remains largely underutilized as a runtime for agents.

### 1.2 My Approach

Recent progress in tool calling and structured output makes it practical to bind natural-language reasoning to a small set of system-level operations. Instead of introducing another application-layer framework, Quine treats the operating system itself as the runtime substrate for agents. It consists of two components:

**Component 1: A Protocol** — A strict mapping from agent concepts to POSIX primitives across four dimensions: identity (PID), interface (streams + exit status), state (memory/env/filesystem), and lifecycle (`fork`/`exec`/`exit`). Section 2 elaborates each dimension; Figure 1 provides the complete mapping.

**Component 2: A Single-Image Runtime** — A unitary executable that implements this protocol. When an agent spawns a child, it instantiates the same runtime image. Parent and child share code but diverge in state.

### 1.3 Contributions

1. **A disciplined agent-to-POSIX mapping** across identity, interface, state, and lifecycle (Section 2).
2. **A unitary recursive runtime** that realizes this mapping through self-instantiation—a single binary that spawns instances of itself (Section 3).
3. **An analysis of structural properties** that emerge from this mapping: composition through standard streams, containment through process isolation, and continuity through `exec` (Section 4).

---

## 2. The POSIX Mapping (~1 page)

*Each dimension of the protocol builds on primitives that trace to the original Unix design [Ritchie & Thompson 1974] and were codified in [IEEE POSIX].*

### 2.1 Identity

| Agent Concept | POSIX Primitive | Semantics |
|---------------|-----------------|-----------|
| Agent instance | Process (PID) | Kernel-managed identity |
| Agent boundary | Address space | Memory isolation between agents |

An agent does not need a framework-assigned UUID. The kernel provides one.

### 2.2 Interface

| Agent Concept | POSIX Primitive | Semantics |
|---------------|-----------------|-----------|
| **Mission** | `argv` | Immutable mission description |
| **Material** | `stdin` | Data input / material |
| **Deliverable** | `stdout` | Data output / deliverable |
| **Diagnostics** | `stderr` | Diagnostics channel |
| **Outcome** | `exit status` | Success (0) or failure (>0) |

### 2.3 State

| Agent Concept | POSIX Primitive | Semantics |
|---------------|-----------------|-----------|
| **Ephemeral** | Process memory | Cleared on `exec`/exit |
| **Scoped** | Environment variables | Preserved across `exec`; inherited by children |
| **Global** | Filesystem | Persistent; shared (default) or isolated via namespaces |

This three-tier state model maps directly to process lifetime scopes.

### 2.4 Lifecycle

| Agent Concept | POSIX Primitive | Semantics |
|---------------|-----------------|-----------|
| **Spawn** | `fork` | Create child process with new mission |
| **Continue** | `exec` | Replace current process image while preserving process-level continuity |
| **Terminate** | `exit(status)` | Signal outcome to parent |

An agent can **spawn** children to delegate work (synchronously via `wait`, or asynchronously via background execution and job control), **continue** itself to escape context limits, or **terminate** to report completion. These three operations, combined with the state hierarchy (§2.3), define the complete lifecycle.

---

## 3. The Runtime (~0.75 page)

*The implementation is a unitary executable that serves as both runtime and agent template.*

### 3.1 Host-Guest Architecture: Separation of Control and Compute

Quine enforces a strict physical separation between deterministic control flow and probabilistic computation.

```
┌─────────────────────────────────────┐
│           Quine Runtime             │
│  ┌───────────┐    ┌──────────────┐  │
│  │   Host    │    │    Guest     │  │
│  │  (Go)     │◄──►│    (LLM)     │  │
│  │           │    │              │  │
│  │ Syscalls  │    │  Decisions   │  │
│  │ I/O       │    │  Tool calls  │  │
│  │ Signals   │    │  Reasoning   │  │
│  └───────────┘    └──────────────┘  │
└─────────────────────────────────────┘
```

**Host** (Local OS Process): Maintains the agent's context, lifecycle, and filesystem state. Because cognition is delegated to remote endpoints, the local process footprint remains small, with memory primarily used for I/O buffers and control state. This keeps the local runtime lightweight and allows agent instances to be managed as ordinary OS processes.

**Guest** (Remote LLM API): Acts as a stateless cognitive oracle. It receives the serialized POSIX state (mapped to prompts) and returns deterministic tool invocations. The Host contains no task-specific logic; decision-making is delegated to the Guest.

**Context Assembly**: The Host avoids application-layer restructuring of the data stream and instead preserves OS-level I/O boundaries when assembling context. It passes file descriptors directly to the shell subprocess and annotates their semantics in the system prompt. Data flows through kernel buffers; the cognitive layer receives only the contract, not the content.

**Tool Interface**: The Host exposes four tools to the Guest—`sh`, `fork`, `exec`, `exit`—named after and semantically aligned with their POSIX counterparts:

| Tool | Semantics |
|------|-----------|
| `sh` | Execute shell command, return output |
| `fork` | Spawn child process(es) running `./quine`; return stdout, stderr, exit status |
| `exec` | Replace current process image; default self re-entry re-execs Quine, but explicit target/`argv` may hand off to a different executable |
| `exit` | Terminate with status code |

### 3.2 POSIX Conformance: A Five-Channel I/O Contract

Quine enforces a five-channel I/O contract that decouples intent, data, execution, and observability:

| Channel | Stream | Content | Consumer |
|---------|--------|---------|----------|
| **Mission** | `argv` | Mission / instructions (text) | Runtime → System Prompt |
| **Material** | `stdin` | Input data (bytes) | Runtime → Agent (internal) |
| **Deliverable** | `stdout` | Pure output (bytes) | Agent → Downstream (external) |
| **Diagnostics** | `stderr` | Diagnostics / failures | Agent → Supervisor (external) |
| **Outcome** | `exit status` | Success (0) / failure (>0) | Process → Parent wait() |

**External view:** From the shell's perspective, Quine is a standard POSIX filter—stdin flows in, stdout flows out, stderr carries diagnostics, exit status signals success or failure.

**Internal mechanism:** To maintain clean context windows while honoring this contract, the runtime uses file descriptor passing. When the agent spawns shell commands, the runtime's stdin/stdout are passed via additional descriptors, allowing the agent to separately capture command feedback (for reasoning) and emit deliverables (for pipelines). This separation enables Quine to act as a drop-in replacement for `grep`, `sed`, or `awk` while preserving cognitive transparency.

### 3.3 Illustrative Shell Compositions

Quine's adherence to standard streams enables composition with classic Unix utilities and multi-agent pipelines.

**Example 1: Single-process replacement.** A cognitive filter that understands intent:

```bash
cat server.log | ./quine "Extract lines indicating auth failures" | sort | uniq -c
```

The cognitive filter acts as a drop-in replacement for `grep`, maintaining the same interface while providing semantic understanding.

**Example 2: Implicit DAG.** Three agents form a reasoning chain:

```bash
git diff HEAD~1 | \
  ./quine "Identify the changed components" | \
  ./quine "Assess risk level for each change" | \
  ./quine "Generate a review checklist"
```

No orchestrator, no shared memory—the shell provides the topology.

**Example 3: Control flow without framework.** Exit codes drive branching:

```bash
./quine "Is this a security vulnerability? Exit 0=yes, 1=no" < report.md && \
  ./quine "Draft disclosure timeline and patch plan" < report.md
```

Shell conditionals (`&&`, `||`) become cognitive branch points.

### 3.4 Implementation Scale

The runtime is implemented in ~5,700 lines of Go (cognitive loop ~3,100 LOC; LLM provider abstractions ~2,600 LOC), compiled to a ~9.8 MB binary. The runtime contains no task-specific logic—all reasoning emerges from the Guest (LLM). Source code is available at https://github.com/kehao95/quine.

---

## 4. Properties (~1 page)

*These properties are not implemented features; they emerge from mapping agents to POSIX processes.*

By representing agents as processes, Quine inherits mechanisms refined over five decades of Unix evolution. I organize these into three categories: what agents can do together (Composition), what agents cannot do (Containment), and how agents persist beyond their limits (Continuity).

### 4.1 Composition

**Recursive Composition.** The architecture exhibits structural self-similarity: a child process is isomorphic to its parent. Both run the same binary and expose the same interface (argv/stdin/stdout/exit), regardless of depth in the process tree. This removes the "worker" vs "manager" distinction—every node can spawn its own subtree. From the parent's perspective, a delegated subtree presents a single process-level interface.

**Universal Composability.** Because Quine adheres to the POSIX filter contract—stdin in, stdout out, stderr for diagnostics, exit status for outcome—it composes with any Unix tool: shell scripts, `xargs`, GNU `parallel`, `cron`, or another Quine instance. Existing infrastructure (makefiles, CI pipelines) can incorporate cognitive agents without modification.

### 4.2 Containment

The process model bounds agent behavior through layered enforcement. Each layer provides guarantees that higher layers cannot bypass.

| Layer | Mechanism | Guarantee |
|-------|-----------|-----------|
| Hardware | MMU | Memory isolation between agents |
| Kernel | cgroups, ulimit, OOM | Resource budgets (CPU, memory, I/O) |
| Kernel | capabilities, namespaces, seccomp | Privilege and visibility bounds |
| Runtime | argv/stdin separation | Instruction-data boundary |

**Hardware: Memory Isolation.** Each agent occupies a separate address space enforced by the MMU. A child crash (SEGV, OOM, infinite loop) cannot corrupt the parent's memory. Isolation comes from silicon, not defensive programming.

**Kernel: Resource Budgets.** `cgroups`, `ulimit`, and the OOM killer enforce budgets on CPU, memory, and I/O. When exceeded, the kernel terminates the process and reports the cause.

**Kernel: Privilege Restriction.** `capabilities` partition root privileges; `namespaces` isolate filesystem, network, and PID visibility; `seccomp` restricts the syscall surface. Runaway or malicious agents are bounded by **kernel enforcement**, not application-layer exception handlers.

**Runtime: Instruction-Data Separation.** `argv` carries the mission (immutable after `exec`); `stdin` carries data. The Host maps these to System Prompt and User Message respectively. Combined with instruction hierarchy training [Wallace et al. 2024], this provides a structural basis that may reduce certain prompt manipulation risks—though systematic security evaluation is outside this paper's scope.

### 4.3 Continuity

Agents face cognitive death (context exhaustion) and physical death (process termination). The process model provides mechanisms to persist across both.

**Continuation across `exec`.** The `exec` syscall replaces the process image while preserving process-level continuity such as PID, environment variables, and open file descriptors. In Quine's default self re-entry path, the current mission `argv` is reused as well. An agent approaching context limits can checkpoint progress in env vars or the filesystem, then `exec` into a fresh instance. This reframes part of the memory problem as deterministic reconstruction: compact scoped state in env vars, working context assembled on demand, and raw history persisted on the filesystem.

**Feedback via stderr and exit status.** The `stderr` stream and exit status form a feedback channel from child to parent. Shell conditionals (`&&`, `||`) become cognitive branch points. Parents implement retry logic, fallback strategies, or error aggregation based on child feedback.

### 4.4 Demonstration Evidence

The following observations demonstrate operational feasibility. Sample sizes are limited; results are illustrative.

**Mission-first streaming.** I evaluated Quine on an MRCR-style needle retrieval task [OpenAI MRCR] spanning 4K–279K tokens (8 samples). Results:

| Context Length | Quine Accuracy | Baseline Accuracy |
|----------------|----------------|-------------------|
| 4K–7K (5 samples) | ≥0.996 (all) | Failed 2/5 |
| 178K–279K (3 samples) | 1.000 (all) | Failed 3/3 |

For long contexts, agents used `exec` for adaptive self-renewal—one run required 9 `exec` cycles for 178K tokens. Results demonstrate feasibility with limited samples (N=8); baseline is one reasonable approach, not optimal.

**Shell-native composition.** In exploratory runs, agents used shell job control (`&`, `wait`, `sort -u`) for parallel coordination. Frequency and reliability are not quantified.

**Continuation across exec.** Agents discovered the checkpoint-self-renewal pattern without explicit instruction, suggesting the model understood `exec` semantics from the tool description.

**Reproducibility**: Implementation, prompts, and execution logs available at https://github.com/kehao95/quine.

---

## 5. Related Work and the Boundaries of POSIX (~1 page)

*Quine sits among existing agent systems as a runtime-level alternative; from there, the discussion turns to where the POSIX model ends and what must extend beyond it.*

### 5.1 Related Systems

**Architectural Positioning.** Quine's design reflects the Exokernel principle [Engler et al. 1995]: the kernel provides protection and resource multiplexing; policy decisions belong in applications. Rather than reimplementing scheduling and isolation in user-space, Quine exposes existing OS mechanisms (`fork`/`exec`/`wait`, address spaces, cgroups) to the LLM through a minimal tool interface.

**Taxonomy of Agent Runtimes.** I classify current systems by their structural relationship to the operating system:

1.  **Host-Coupled Extensions (e.g., Cursor, Copilot)**: Integrated within host environments (IDEs); lack independent lifecycles.
2.  **Application-Layer Schedulers (e.g., LangGraph, AutoGen [Wu et al. 2023], CrewAI)**: User-space scheduling and message buses; agents as in-process objects with simulated fault isolation.
3.  **Sandboxed Loops (e.g., Devin, OpenHands [Wang et al. 2024], SWE-agent [Yang et al. 2024])**: External controller managing containerized tools; agent brain remains in shared process.
4.  **Pseudo-OS Middleware (e.g., AIOS [Mei et al. 2024], agentOS [Li et al. 2026], UFO 2 [Zhang et al. 2025])**: User-space simulation of kernel abstractions; agents as threads within a single host process.

**Relationship to Prior Work**: Quine represents a distinct design point. Existing agent runtimes typically either manage agents at the application layer or use the OS primarily as a tool sandbox, rather than directly realizing each agent as a native POSIX process with the full reasoning-and-acting loop exposed through standard process interfaces.

This difference is about runtime organization, not end-task superiority.
-   Instead of a custom graph engine, it uses the **OS Scheduler** (CFS, cgroups).
-   Instead of proprietary message formats, it uses **Standard Streams (stdin/stdout)**.
-   Instead of an external controller, it uses **Process Lifecycle (`fork`/`exec`/`exit`)**.
-   Instead of containerized isolation, it uses **Process Address Spaces**.

Having located Quine among current systems, I now turn from comparative structure to the limits of the substrate itself.

### 5.2 Architectural Boundaries of the POSIX Model

POSIX was designed to manage execution, not cognition. While it addresses the external lifecycle of an agent—identity, communication, termination—it can only represent agents at the boundary of execution, not at the deeper levels of cognitive structure, situated perception, temporal revisability, and distributed composition.

I identify four structural mismatches that necessitate extensions beyond the POSIX model:

1.  **Unrepresented Internal Structure (The Structure Problem)**: Process memory is a private black box. An agent's internal organization—context window, attention layout, working memory—has no runtime representation. Diagnosis involves reading logs and inferring internal state.
2.  **Undifferentiated Worldview (The Space Problem)**: POSIX assumes a single shared world, partially filtered by permissions. Agents may require distinct and task-relative worldviews—not merely restricted access, but semantically scoped environments organized around relevance. This concern echoes, but also extends beyond, the Plan 9 idea that each process may inhabit its own namespace [Pike et al. 1993].
3.  **Irreversible Time (The Time Problem)**: POSIX treats side-effects as final in the order they occur. Standard filesystems are mutable in place; when an agent overwrites a file, the previous version is lost. This creates a mismatch: cognition is provisional and revisable, but the environment records only the latest state.
4.  **Local-Bound Composition (The Mesh Problem)**: POSIX composition is topologically local: its abstractions presume adjacency. Scaling to distributed systems requires abandoning the file metaphor for application-layer protocols (HTTP, gRPC), creating a semantic fracture between local and remote agents.

The first four limitations are architectural—addressable through OS-level extensions rather than application-layer workarounds. This paper focuses on the POSIX mapping; the extensions are future work.

These mismatches do not invalidate the POSIX mapping; they mark the point at which a runtime for agents must grow beyond process semantics into cognitive semantics.

### 5.3 Future Directions

The POSIX mapping presented here should be understood as a foundation rather than a complete theory of agent runtime. It establishes a robust account of lifecycle, isolation, and composition at the level of execution, but it does not yet provide corresponding primitives for cognition itself.

More generally, POSIX treats the agent as a process: an externally manageable unit with inputs, outputs, permissions, and termination conditions. What it lacks is a richer account of the agent as a structured cognitive entity situated in a task-relative world, acting within revisable time, and composing across distributed topologies without semantic fracture.

My ongoing work therefore explores four directions beyond the POSIX layer:

- **Structure**: giving the agent's internal organization a first-class runtime representation, rather than treating cognition as an opaque black box.
- **Space**: replacing the assumption of a single shared filesystem with task-relative worldviews, so that each agent inhabits a scoped and meaningful environment.
- **Time**: extending execution beyond irreversible mutation toward revisable histories, in which cognition and side-effects can be aligned.
- **Mesh**: preserving compositional continuity across machines, so that local and distributed agents remain instances of the same operational model.

These directions are stated at the level of architecture rather than mechanism. Specific realizations may involve synthetic filesystems, namespace construction, snapshotting, or network-transparent protocols, but the underlying claim is broader: POSIX is sufficient as a first runtime layer for agents, yet insufficient as their final operating model.

The present paper addresses only that first layer. The broader question—what an operating system for cognition should expose, delimit, remember, and compose—belongs to future work. Conceptually, this move is closer to the Plan 9 lineage [Pike et al. 1995] than to conventional POSIX extension: the aim is not merely to add mechanisms, but to rethink what the runtime should make visible, nameable, and composable.

---

## 6. Conclusion (~0.25 page)

Quine shows that LLM agents can be architected as native POSIX processes, offering a different design point from application-layer frameworks. By mapping agents onto POSIX processes, Quine inherits well-understood process semantics for isolation, composition, and lifecycle control—mechanisms refined over fifty years. The cost is accepting Unix's assumptions: deterministic composition, text-stream interfaces, and shared-nothing isolation.

Quine is both a working implementation and a concrete runtime design point: it shows that agent isolation, composition, and lifecycle control can be inherited from the OS rather than rebuilt above it. The broader questions of internal structure, situated worldviews, revisable execution, and distributed continuity lie beyond the scope of the present work. Whether OS-native runtimes offer advantages for large-scale agent systems remains an empirical question.

Source code: https://github.com/kehao95/quine

---

## Figures & Tables

1. **The POSIX Mapping** (Table, §2): Agent Concept ↔ POSIX Primitive correspondence across identity, interface, state, lifecycle
2. **Interface: Five-Channel I/O Contract** (Diagram, §2.2): Five channels showing Mission (argv), Material (stdin), Deliverable (stdout), Diagnostics (stderr), Outcome (exit) with upstream/downstream composition
3. **Lifecycle Primitives** (Diagram, §2.4): Three lifecycle operations—fork (spawn), exec (continue), exit (terminate), showing three-tier state model (Mission/Scoped/Ephemeral)
4. **Host-Guest Architecture** (Diagram, §3.1): Runtime structure showing separation of control (Host) and compute (Guest)

---

## Estimated Length

| Section | Pages |
|---------|-------|
| Abstract | 0.15 |
| Introduction | 0.5 |
| POSIX Mapping | 1.0 |
| The Runtime | 1.0 |
| Properties | 1.0 |
| Related Work & Boundaries | 1.0 |
| Conclusion | 0.25 |
| **Total** | **~5 pages** |

---
---

## References (Draft)

*Ordered by first appearance in text. All citations use [Author et al. Year] format.*

**§1.1 — Problem Statement:**

1. **[LangChain]** LangChain, Inc. "LangChain: Building applications with LLMs through composability." https://github.com/langchain-ai/langchain
   - *Used in*: §1.1, §5.1 (Orchestrated Runtimes)
   - *Note*: Open-source framework; no canonical academic paper.

2. **[Wu et al. 2023]** Wu, Q., et al. "AutoGen: Enabling Next-Gen LLM Applications via Multi-Agent Conversation Framework." *arXiv:2308.08155*, 2023. (Microsoft Research)
   - *Used in*: §1.1, §5.1 (Orchestrated Runtimes representative)
   - *Key contribution*: "Conversable Agent" abstraction, conversation programming paradigm.

3. **[CrewAI]** CrewAI, Inc. "CrewAI: Framework for orchestrating role-playing autonomous AI agents." https://github.com/crewAIInc/crewAI
   - *Used in*: §1.1, §5.1 (Orchestrated Runtimes)
   - *Note*: Open-source framework; no canonical academic paper.

**§2 — POSIX Mapping:**

4. **[Ritchie & Thompson 1974]** Ritchie, D.M. & Thompson, K. "The UNIX Time-Sharing System." *Communications of the ACM*, 17(7), 1974.
   - *Used in*: §2 (historical foundation of process/stream model)
   - *Key insight*: "Everything is a file" and process isolation as organizing principles.

5. **[IEEE POSIX]** IEEE Std 1003.1-2017. "Standard for Information Technology—Portable Operating System Interface (POSIX)."
   - *Used in*: §2 (POSIX Mapping), §3 (Conformance)

**§5.1 — Related Systems (Architectural Positioning):**

6. **[Engler et al. 1995]** Engler, D.R., Kaashoek, M.F., & O'Toole, J. "Exokernel: An Operating System Architecture for Application-Level Resource Management." *SOSP 1995*.
   - *Used in*: §5.1 (architectural positioning)
   - *Key insight*: OS should provide protection/multiplexing only; policy belongs in applications.

**§4.2 — Execution Safety:**

7. **[Wallace et al. 2024]** Wallace, E., Xiao, K., Leike, R., Weng, L., Heidecke, J., & Beutel, A. "The Instruction Hierarchy: Training LLMs to Prioritize Privileged Instructions." *arXiv:2404.13208*, 2024. (OpenAI)
   - *Used in*: §4.2 (Instruction-Data Separation)
   - *Key finding*: Priority hierarchy (System > User > Tool Output); models can be trained to selectively ignore lower-privileged instructions even under adversarial attack.

**§5.1 — Related Systems:**

6. **[McIlroy 1978]** McIlroy, M.D. "Unix Time-Sharing System: Foreword." *Bell System Technical Journal*, 57(6), 1978.
   - *Used in*: §5.1, §6 (Unix Philosophy quotes)

7. **[Wang et al. 2024]** Wang, X., et al. "OpenHands: An Open Platform for AI Software Developers as Generalist Agents." *arXiv*, 2024. (CMU et al.)
   - *Used in*: §5.1 (Sandboxed Loops representative)
   - *Key contribution*: Event Stream architecture, modular SDK for agent development.

8. **[Yang et al. 2024]** Yang, J., et al. "SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering." *arXiv*, 2024. (Princeton)
   - *Used in*: §5.1 (Sandboxed Loops representative)
   - *Key contribution*: Agent-Computer Interface (ACI) concept; optimized shell interactions for LLM context windows.

9. **[Mei et al. 2024]** Mei, K., et al. "AIOS: LLM Agent Operating System." *arXiv:2403.16971*, 2024. (COLM 2025)
   - *Used in*: §5.1 (Pseudo-OS Middleware)
   - *Key contribution*: OS-inspired abstractions for LLM agents (scheduler, memory manager), but implemented as application-layer middleware, not kernel-level mapping.
   - *Distinction*: All agents are Python objects within a single process; Quine maps each agent to a separate PID.

10. **[Li et al. 2026]** Li, C., Liu, X., Meng, X., & Zhao, X. "Architecting AgentOS: From Token-Level Context to Emergent System-Level Intelligence." *arXiv:2602.20934*, 2026.
    - *Used in*: §5.1 (Pseudo-OS Middleware)
    - *Key contribution*: Conceptual framework treating LLM as "Reasoning Kernel" with Deep Context Management and Semantic Slicing.
    - *Distinction*: Theoretical framework mapping OS abstractions onto LLM constructs; Quine maps agents *to* actual OS processes.

11. **[Zhang et al. 2025]** Zhang, C., Huang, H., Ni, C., et al. "UFO2: The Desktop AgentOS." *arXiv:2504.14603*, 2025. (Microsoft Research)
    - *Used in*: §5.1 (Pseudo-OS Middleware)
    - *Key contribution*: Multiagent AgentOS for Windows desktops with HostAgent/AppAgent hierarchy, hybrid UI detection, speculative multi-action planning.
    - *Distinction*: Desktop automation focus; agents coordinate within single control framework. Quine: each agent is an independent POSIX process.

**§4.4 — Validation Evidence:**

14. **[OpenAI MRCR]** OpenAI. "Multi-turn Retrieval Context Reasoning (MRCR) Benchmark." https://huggingface.co/datasets/openai/mrcr
    - *Used in*: §4.4 (Mission-first streaming validation)
    - *Description*: Long-context needle retrieval benchmark requiring ordinal retrieval (N-th occurrence) in GPT-4o generated conversations.

**§6 — Design Principle:**

12. **[Hewitt et al. 1973]** Hewitt, C., Bishop, P., & Steiger, R. "A Universal Modular ACTOR Formalism for Artificial Intelligence." *IJCAI*, 1973.
    - *Used in*: §6 (Design Principle—Actor Model parallel)
    - *Key insight*: Actors as independent computational entities with message-passing, no shared state, and isolated failure domains—the theoretical foundation Erlang/OTP later operationalized.

**§6 — Beyond POSIX:**

13. **[Pike et al. 1995]** Pike, R., Presotto, D., Dorward, S., Flandrena, B., Thompson, K., Trickey, H., & Winterbottom, P. "Plan 9 from Bell Labs." *Computing Systems*, 8(3), 1995.
    - *Used in*: §5.3 (Plan 9 philosophy as conceptual lineage for future directions)

14. **[Pike et al. 1993]** Pike, R., Presotto, D., Thompson, K., & Trickey, H. "The Use of Name Spaces in Plan 9." *Operating Systems Review*, 27(2), 1993.
    - *Used in*: §5.2.2 (Plan 9 namespace concept as precedent for task-relative worldviews)

---

**Research Notes:**

- **Taxonomy validation**: Surveyed all major agent systems (2023-2026). No prior work maps agents with complete reasoning-and-acting loops directly onto POSIX processes. Existing approaches either reinvent OS primitives in user-space (Orchestrated, Pseudo-OS) or use OS as sandbox for tools only (Sandboxed), but none use OS primitives AS the agent identity.
- **Six-dimensional analysis framework**: Runtime Architecture, Process Model, IPC Mechanism, Lifecycle Management, Isolation Mechanism, Composition Method — all current systems score "Application-Layer" on at least 5/6 dimensions.

---

## Deferred to Track B

| Content | Reason |
|---------|--------|
| Formal definitions (Agent tuple, Lifecycle FSM) | Requires more space |
| Microbenchmarks (spawn latency, throughput) | Requires controlled experiments |
| Quantitative Benchmarks vs Frameworks | Requires fair benchmarking (LangChain, AutoGen) |
| Security analysis (prompt injection depth) | Tangential to core thesis |
| Scaling analysis (4-64 agents) | Requires additional experiments |
| Synthetic data pipeline details | Implementation detail |
| Classical agent theory mapping [Wooldridge & Jennings 1995] | Narrative is systems-first, not agent-theory-first |
| Concurrency control discussion [Dijkstra 1965; Hoare 1974] | Paper inherits OS mechanisms, doesn't implement them |
