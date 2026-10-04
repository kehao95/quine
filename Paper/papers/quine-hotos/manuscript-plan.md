# Manuscript Plan: HotOS Track — Paradigm-Level Reframing

**Target**: HotOS 2025 (workshop paper, 5 pages excluding references)  
**Venue Fit**: HotOS welcomes "new ideas, new research directions, or challenges to conventional wisdom"—exactly where a paradigm reframing belongs.

---

## Strategic Positioning

**What This Version Does**: Argues that the agent field has placed the agent at the **wrong abstraction layer**, and that fixing this is not a feature addition but an **ontological relocation**. The five layers are framed as a **necessity chain** (substrate completion), not a feature roadmap.

**Central Claim**: "The main claim of this paper is not that agents *can use* operating-system mechanisms, but that they should be *defined* at that layer."

**Paper Type**: This is a *reframing* paper, not a *proof of superiority* paper. It argues for where the agent should be defined, specifies the consequences, and opens a research agenda.

**Relationship to arXiv Version**: The arXiv version establishes the **minimal verifiable claim** (Layer 1 works). This HotOS version extends that into the **full paradigm argument** (wrong abstraction layer + substrate completion). They form a continuous research program, not conflicting claims.

**Key Vocabulary**: categorical mismatch, substrate (not framework), ontological completion, necessity chain, reframes/relocates (not solves/proves)

---

## Version History

- v51: Paradigm reframing (this version)
- v50: Architecture-focused version -> now `../quine-arxiv/manuscript-plan.md`
- Full changelog: [`archive/track-system-meta/CHANGELOG-outline-track-a.md`](../../archive/track-system-meta/CHANGELOG-outline-track-a.md)

---

## Title

**"LLM Agents as Operating-System Objects"**

*Alternative titles (in order of preference):*
1. "From Agent Frameworks to Agent Substrates: A POSIX/Plan 9 Perspective"
2. "Reframing LLM Agents as OS Primitives"
3. "Quine: Toward an Operating-System Substrate for LLM Agents"

---

## Abstract (~150 words)

> Contemporary LLM agent systems are predominantly built as application-layer runtimes. This design choice forces frameworks to reconstruct, in user space, capabilities that operating systems already provide: process isolation, lifecycle control, inter-process communication, persistence, and composition. We argue that this is a categorical mismatch. An agent should be treated not as a library object but as an operating-system object.
>
> This paper makes the claim precise through a five-layer architecture. **Layer 1** realizes agents as POSIX processes—PID as identity, streams as interface, exit status as outcome—and is validated by Quine, an open-source implementation. **Layers 2–5** extend this foundation through Plan 9 principles: unified namespace makes the agent inspectable (Structure); scoped filesystem gives it a subjective world (Space); CoW snapshots enable mind and world to restore independently (Time); 9P mounts preserve the abstraction across machine boundaries (Mesh).
>
> The contribution is ontological before it is implementational: it relocates the agent from the framework layer to the operating-system substrate, and specifies the research program that follows.

---

## 1. Introduction (~0.75 page)

### 1.1 The Problem: A Categorical Mismatch

The rapid development of LLM agents has spawned a proliferation of frameworks—LangChain, AutoGen [Wu et al. 2023], CrewAI, AIOS [Mei et al. 2024]. These systems differ in orchestration strategy, but share a common architectural choice: the agent is treated as an application-layer object—a Python class, a graph node, a message handler—managed by framework code running in user space.

This choice has consequences. Scheduling, isolation, persistence, retry logic, and inter-agent communication must all be reimplemented as framework features rather than inherited as operating-system structure. An unhandled exception in one agent risks cascading failure across the orchestrating process. Memory management becomes probabilistic summarization rather than system-level resource control. Composition requires custom protocols rather than standard streams.

This is not a local engineering pattern; it is a field-wide convergence toward user-space reimplementation of operating-system concerns. The duplication is not innovation—it is evidence that the agent has been placed at the wrong abstraction layer.

### 1.2 The Thesis: Agents as Operating-System Objects

> **The main claim of this paper is not that agents can use operating-system mechanisms, but that they should be *defined* at that layer.**

The claim is not that application-layer runtimes are useless, but that they may occupy the wrong architectural layer for the long-term evolution of agent systems. When an agent is defined as an OS object rather than a framework object, concerns that frameworks must simulate—fault domains, resource isolation, lifecycle management, composition—become structural properties of the substrate.

### 1.3 Why Now

The question is timely not only because language models can now reliably operate external tools, but because the current generation of agent frameworks has made the cost of staying in the application layer visible. Tool schemas act as contracts; the runtime validates correctness at the system level. The semantic gap between intent and syscall has become bridgeable. Meanwhile, the pattern of rebuilding OS primitives in user space has become undeniable—and the costs are accumulating.

### 1.4 The Five Layers as Necessity Chain

This paper specifies a five-layer architecture—not as a feature roadmap, but as successive repairs to the mismatch between execution-oriented POSIX and cognition-oriented workloads. Each layer addresses an insufficiency exposed by the previous:

**Layer 1 (Identity)**: An agent is a POSIX process—PID as identity, streams as interface, exit status as outcome. This realizes the Actor Model [Hewitt et al. 1973] using native primitives. *But Layer 1 leaves representation fragmented across file descriptors, argv, env, and external logs.*

**Layer 2 (Structure)**: The agent becomes a directory (`/agent/<pid>/`). All state—streams, mission, history—becomes addressable paths. Unix tools apply directly. *But unified representation does not define the agent's operational environment.*

**Layer 3 (Space)**: The agent inhabits a scoped world (`/agent/<pid>/fs/`), a CoW overlay where mutations land in a private layer. The host remains pristine; the runtime can diff mutations after each command. *But forward-only mutation cannot support cognitive backtracking.*

**Layer 4 (Time)**: CoW snapshots (`/agent/<pid>/snapshots/`) make history first-class. The central insight: because mind (`tape/`) and world (`fs/`) are separately materialized, they can restore independently. Revert filesystem damage while preserving lessons learned; or restore both for branch exploration. *This is the most original extension—not inherited from classic Plan 9, but a natural extension of its namespace philosophy into the temporal dimension.*

**Layer 5 (Mesh)**: 9P mounts remote agents as local directories. The machine boundary does not break the agent abstraction. Local and remote agents share the same interface.

These are not optional enhancements but consequences of treating the agent as an OS object.

### 1.5 This Paper

We present **Quine**, an open-source Go runtime that validates Layer 1—demonstrating that an agent can be realized as a native POSIX process while retaining a complete reasoning-and-acting loop. Layers 2–5 are presented as architectural extensions whose value remains to be established empirically. The purpose of this paper is to reframe the design space, not to claim universal superiority over current frameworks.

### 1.6 Contributions

1. **A categorical reframing**: The claim that contemporary agent systems place the agent at the wrong abstraction layer, and that this mismatch explains the recurrent reimplementation of OS primitives in user space.

2. **A five-layer architecture**: Identity, Structure, Space, Time, Mesh—specified as successive completions of the agent substrate, not as independent features.

3. **Quine**: An existence proof for Layer 1, demonstrating that POSIX-native agents are feasible today. Source: https://github.com/kehao95/quine

---

## 2. Layer 1: The POSIX Foundation (~0.5 page)

*This section defines the complete mapping from agent concepts to POSIX primitives.*

**Figure 1: Agent-to-POSIX Mapping**

```
┌──────────────┬─────────────────┬───────────────────────────────────┐
│  Dimension   │  Agent Concept  │  POSIX Primitive                  │
├──────────────┼─────────────────┼───────────────────────────────────┤
│  IDENTITY    │  Instance       │  Process (PID)                    │
│              │  Boundary       │  Address space                    │
├──────────────┼─────────────────┼───────────────────────────────────┤
│  INTERFACE   │  Mission        │  argv                             │
│              │  Input          │  stdin                            │
│              │  Output         │  stdout                           │
│              │  Diagnostics    │  stderr                           │
│              │  Outcome        │  exit status                      │
├──────────────┼─────────────────┼───────────────────────────────────┤
│  STATE       │  Ephemeral      │  Process memory                   │
│              │  Scoped         │  Environment variables            │
│              │  Persistent     │  Filesystem                       │
├──────────────┼─────────────────┼───────────────────────────────────┤
│  LIFECYCLE   │  Spawn          │  fork + exec                      │
│              │  Continue       │  exec (self)                      │
│              │  Terminate      │  exit(status)                     │
└──────────────┴─────────────────┴───────────────────────────────────┘
```

The mapping is strict: an agent does not need a framework-assigned UUID—the kernel provides one (PID). The three-tier state model (ephemeral/scoped/persistent) maps directly to process lifetime scopes. An agent can **spawn** children to delegate work, **continue** itself via `exec` to escape context limits while preserving identity, or **terminate** to report outcome. These primitives, combined with standard streams, define the complete agent lifecycle without framework abstractions.

---

## 3. Validation: The Quine Runtime (~0.5 page)

*Quine validates the minimal claim of this paper: that an agent can be realized as a native POSIX process while retaining a complete reasoning-and-acting loop. Source: https://github.com/kehao95/quine*

*It does not validate Layers 2–5, which remain architectural extensions specified here as a research program.*

### 3.1 Host-Guest Architecture

Quine enforces strict separation between deterministic control (Host) and probabilistic cognition (Guest).

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

The **Host** (local process) maintains lifecycle and I/O state; because cognition is delegated to remote endpoints, local memory holds only buffers—enabling the OS scheduler to multiplex thousands of agents as I/O-bound workloads. The **Guest** (remote LLM) acts as a stateless oracle, receiving serialized state and returning tool invocations. The Host exposes four tools—`sh`, `fork`, `exec`, `exit`—semantically aligned with their POSIX counterparts.

### 3.2 Composition Without Frameworks

Quine is a standard POSIX executable (argv/stdin/stdout/exit), composable with any Unix tool. Two examples demonstrate framework-free orchestration:

**Implicit DAG.** Three agents form a reasoning chain without orchestrator:

```bash
git diff HEAD~1 | \
  ./quine "Identify the changed components" | \
  ./quine "Assess risk level for each change" | \
  ./quine "Generate a review checklist"
```

The shell provides topology; no shared memory, no coordination protocol.

**Cognitive Branching.** Exit codes drive control flow:

```bash
./quine "Is this a security vulnerability? Exit 0=yes, 1=no" < report.md && \
  ./quine "Draft disclosure timeline and patch plan" < report.md
```

Shell conditionals (`&&`, `||`, `if`) become cognitive branch points—the control flow primitives that LangGraph reimplements are native to the shell.

---

## 4. Structural Consequences (~0.5 page)

*These properties are not implemented features—they are structural consequences of the POSIX mapping. They represent what the substrate choice makes possible, not what has been empirically demonstrated.*

| Property | Mechanism | Structural Basis |
|----------|-----------|------------------|
| **Fault Isolation** | Address space (MMU) | Child crash cannot corrupt parent; kernel reclaims resources |
| **Instruction-Data Separation** | argv → System Prompt; stdin → User Message | Provides structural basis for injection resistance |
| **Unbounded Continuation** | `exec` preserves PID + env vars | Creates conditions for arbitrary-length tasks within finite context |
| **Semantic Feedback** | stderr + exit status | Parent can adapt strategy based on child outcomes |

**Instruction-Data Separation** deserves emphasis: `argv` carries mission (instruction), `stdin` carries material (data)—structurally separate at the OS level. The Host maps `argv` exclusively to System Prompt and `stdin` to User Message. Combined with instruction hierarchy training [Wallace et al. 2024] (System > User > Tool), this provides a structural basis for defense-in-depth against injection attacks—though the degree of protection remains to be evaluated empirically.

**Unbounded Continuation** suggests a path from probabilistic memory maintenance toward deterministic reconstruction: an agent clears accumulated context (and hallucinations) via `exec` while preserving compressed state in environment variables. Each generation inherits distilled insights; context is *constructed* from persistent logs, not inherited wholesale. This separation of storage (immutable filesystem) from cognition (assembled context) permits unbounded continuation in principle, subject to reconstruction fidelity.

---

## 5. Related Work (~0.25 page)

**Taxonomy.** Current agent systems can be classified by their structural relationship to the operating system. They differ in orchestration strategy, but largely share the application-layer placement of the agent abstraction:

1.  **Host-Coupled Extensions (e.g., Cursor, Copilot)**: Architectures tightly integrated within specific host environments (IDEs). Agents lack independent lifecycles.

2.  **Application-Layer Schedulers (e.g., LangGraph, AutoGen [Wu et al. 2023], CrewAI)**: Runtimes that construct graph-based scheduling and message buses in user-space. Agents are managed as in-process objects; fault domains are simulated rather than hardware-backed.

3.  **Sandboxed Loops (e.g., Devin, OpenHands [Wang et al. 2024], SWE-agent [Yang et al. 2024])**: Containerized execution environments managed by an external controller. The agent's *tools* are isolated; the agent *brain* remains in a shared process.

4.  **Pseudo-OS Middleware (e.g., AIOS [Mei et al. 2024], agentOS [Li et al. 2026], UFO 2 [Zhang et al. 2025])**: Application-layer middleware that simulates kernel abstractions in user-space. Despite using OS terminology, all managed agents are threads or data structures within a single host process.

**Distinction.** Our distinction from these systems is not primarily one of feature set, but of abstraction layer. Quine is **POSIX-native**: each agent is a separate PID with its own address space, lifecycle, and stream interface. The contribution is not a better framework, but a relocation of where the agent is defined.

---

## 6. Completing the Agent Substrate (~1.5 pages)

> "Everything is a file." — Plan 9 from Bell Labs [Pike et al. 1995]

**From Execution Identity to Cognitive Substrate.** Layer 1 establishes the agent as a process—an execution identity with lifecycle and basic composition. But a process alone is not yet a substrate for cognition. POSIX solves *processhood*; it does not solve *substratehood*.

The move from POSIX to Plan 9 is not a feature extension but an ontological completion. POSIX gives the agent an execution identity. Plan 9 supplies the missing ontology of representation—the principle that if everything is a file, then everything becomes observable, composable, and network-transparent.

Layers 2–5 complete the substrate along four dimensions. Once the agent is treated as an OS object, these properties cease to be optional framework features and become consequences of substrate choice.

**Figure 4: Substrate Completion**

| Layer | Insufficiency Addressed | Principle | Outcome |
|-------|------------------------|-----------|---------|
| **2: Structure** | State fragmented across mechanisms | Unified namespace | Agent becomes legible |
| **3: Space** | Single global filesystem | Scoped namespace + CoW | Agent gets subjective world |
| **4: Time** | Forward-only mutation | Separate mind/world snapshots | Reversibility becomes first-class |
| **5: Mesh** | Local-only composition | Network transparency | Abstraction survives distance |

### 6.1 Layer 2: Structure — Making the Agent Legible

**The Insufficiency.** POSIX provides execution identity, but leaves representation fragmented. An agent's cognitive state is scattered: stdin/stdout/stderr as file descriptors, mission in argv, wisdom in environment variables, conversation history in external logs. For cognitive agents, introspection is essential—but the process abstraction hides everything behind opaque memory.

**The Completion.** Unify all agent-relevant state into a single addressable namespace. The agent becomes a directory at `/agent/<pid>/`:

```
/agent/<pid>/
├── stdin, stdout, stderr    # streams as files
├── argv                     # mission (read-only)
├── env/                     # environment variables
├── tape/                    # conversation history
│   ├── raw.jsonl           # complete log
│   └── summary             # lazy: computed on read
└── status                   # current state
```

A FUSE implementation enables **demand-driven execution**: reading `tape/summary` triggers summarization on-demand. The 50-year Unix toolchain becomes the agent debugging toolkit: `grep` the tape, `diff` two agents, `cat` the mission.

**The Structural Consequence.** Once representation is unified, legibility becomes a property of the substrate, not a feature to be implemented. Agent cognition becomes operable through file semantics.

### 6.2 Layer 3: Space — Giving the Agent a Subjective World

**The Insufficiency.** Layer 2 makes the agent legible, but does not define its operational environment. POSIX enforces a single global filesystem—agents can pollute host state with no undo, or proceed under delusion when mutations fail silently.

**The Completion.** Scope the world to the agent. Layer 3 adds `/agent/<pid>/fs/`, a **Copy-on-Write overlay** where all writes land in a private upper layer:

```
/agent/<pid>/
├── ...              # (Layer 2 state)
└── fs/              # agent's subjective filesystem
    ├── upper/       # CoW — all writes land here
    └── merged/      # what the agent sees
```

**The Fifth Channel.** Because all mutations are captured in `fs/upper/`, the runtime can diff before and after each command. This diff—injected into tool results alongside stdout/stderr—closes a gap: POSIX gives the agent outputs *about* action, but not evidence of world-state change. When a faulty `sed -i` produces no filesystem changes, the agent sees an empty `[FS MUTATIONS]` block.

**The Structural Consequence.** Namespace isolation enforces cognitive boundaries through OS-level structure rather than probabilistic instruction-following. Mistakes become ephemeral; only success persists.

### 6.3 Layer 4: Time — Decoupling Mental and Environmental Reversibility

**The Insufficiency.** Layers 2–3 give the agent legible structure and a scoped world. But POSIX models execution as forward mutation—filesystems mutate in place, history is lost. When an agent's "mind" (context window) can backtrack via prompt manipulation, but its "world" (filesystem) cannot rollback side-effects, they desynchronize.

This asymmetry between reversible cognition and irreversible side-effects is the root cause of cascading failures in long-running autonomous agents.

**The Core Insight.** The central extension of this paper: once mind and world are separately materialized, rollback becomes a first-class operation rather than an ad hoc recovery strategy. Because `tape/` (mind) and `fs/` (world) are separate subdirectories, they can restore independently:

- **Revert world, keep mind**: Undo filesystem damage while preserving the lesson learned
- **Revert both**: Full rollback for branch exploration
- **Revert mind, keep world**: Rare, but possible when reasoning diverged but changes were correct

**The Completion.** Layer 4 adds CoW snapshots of the entire agent directory:

```
/agent/<pid>/
├── tape/, env/, fs/     # current state
└── snapshots/           # temporal navigation
    ├── turn-001/        # complete state at turn 1
    └── turn-002/        # ...
```

**Branching Futures.** Because snapshots are cheap (CoW shares unchanged data), exploring alternative futures becomes trivial: fork from a checkpoint, run variants, compare outcomes via `diff`. Agent execution ceases to be a single narrative and becomes a search over recoverable trajectories.

**The Structural Consequence.** Time becomes a filesystem operation. The asymmetry between reversible cognition and irreversible side-effects dissolves. This is the most original extension in the architecture—not directly inherited from classic Plan 9, but a natural extension of its namespace philosophy into the temporal dimension.

### 6.4 Layer 5: Mesh — Preserving the Abstraction Across Distance

**The Insufficiency.** Layers 2–4 define a complete local substrate. But POSIX pipes are local. Scaling to distributed swarms typically requires abandoning the file metaphor for application-layer protocols—creating a semantic fracture between local and remote agents.

**The Completion.** The machine boundary should not break the agent abstraction. Layer 5 uses 9P to mount remote agents as local directories:

```
/agent/
├── 1001/                # local agent
├── 1002/                # local agent
└── remote-host:2001/    # remote agent (9P mount)
    ├── stdin            # write → delivered to remote
    └── stdout           # read → remote output
```

**The Structural Consequence.** Local and remote agents share the same interface. Physical machine boundaries become invisible. A swarm is not a graph—it is a shared directory.

### 6.5 Synthesis: From Process Identity to Cognitive Substrate

Taken together, the five layers define a shift in where the agent is thought to reside: no longer inside an application runtime, but inside the operating-system substrate itself.

The significance of this shift is not that it introduces new agent features, but that it relocates existing concerns—identity, state, rollback, communication—to a lower and more stable abstraction boundary.

**Table 1: The Five-Layer Agent Substrate**

| Layer | Name | What It Completes | Core Mechanism |
|-------|------|-------------------|----------------|
| **1** | Identity | Execution identity | POSIX Process |
| **2** | Structure | Unified representation | `/agent/<pid>/` (FUSE) |
| **3** | Space | Scoped environment | `/agent/<pid>/fs/` (OverlayFS) |
| **4** | Time | Temporal reversibility | `/agent/<pid>/snapshots/` (CoW) |
| **5** | Mesh | Network transparency | 9P mount |

---

## 7. Scope and Limitations (~0.1 page)

This paper does not claim that all agent workloads should be mapped to operating-system primitives, nor that POSIX/Plan 9 abstractions exhaust the design space. The claim is narrower: for agents that act through tools, external state, and delegation, the operating system offers a more fundamental substrate than application-layer orchestration.

This paper validates only the foundational claim at Layer 1. Layers 2–5 are presented as architectural extensions whose value remains to be established empirically. Whether this substrate yields broad empirical advantages over current frameworks is an open question that this paper does not attempt to resolve.

---

## 8. Conclusion (~0.2 page)

The main claim of this paper is not that agents can use operating-system mechanisms, but that they should be *defined* at that layer. Contemporary agent frameworks reconstruct, in user space, capabilities that operating systems already provide. This is evidence of a categorical mismatch: the agent has been placed at the wrong abstraction layer.

Quine demonstrates that the alternative is feasible—agents can be realized as native POSIX processes while retaining complete reasoning-and-acting loops. Layers 2–5 specify how this foundation can be completed into a full cognitive substrate, drawing on Plan 9's insight that uniform namespaces make heterogeneous concerns observable and composable.

The contribution is therefore ontological before it is implementational. It relocates the agent from the framework layer to the operating-system substrate, and opens a research agenda around what that relocation makes possible.

Source code: https://github.com/kehao95/quine

---

## Figures (3)

1. **The POSIX Mapping** (Figure 1): §2 — Agent Concept ↔ POSIX Primitive
2. **Host-Guest Architecture** (Figure 2): §3.1 — Runtime structure diagram
3. **Substrate Completion** (Table 1): §6 — The five-layer progression

---

## Estimated Length

| Section | Pages |
|---------|-------|
| Abstract | 0.15 |
| Introduction | 0.75 |
| Layer 1: POSIX Foundation | 0.5 |
| Validation: Quine Runtime | 0.5 |
| Structural Consequences | 0.4 |
| Related Work | 0.25 |
| Completing the Substrate: Layers 2–5 | **~1.3** |
| Scope and Limitations | 0.1 |
| Conclusion | 0.2 |
| **Total** | **~4.15 pages** |

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

**§1.3 — Actor Model:**

4. **[Hewitt et al. 1973]** Hewitt, C., Bishop, P., & Steiger, R. "A Universal Modular ACTOR Formalism for Artificial Intelligence." *IJCAI*, 1973.
   - *Used in*: §1.3 (Identity theoretical foundation)
   - *Key insight*: Actors as independent computational entities with message-passing, no shared state, and isolated failure domains.

**§2 — POSIX Mapping:**

5. **[IEEE POSIX]** IEEE Std 1003.1-2017. "Standard for Information Technology—Portable Operating System Interface (POSIX)."
   - *Used in*: §2 (POSIX Mapping), §3 (Conformance)

**§4 — Execution Safety:**

6. **[Wallace et al. 2024]** Wallace, E., Xiao, K., Leike, R., Weng, L., Heidecke, J., & Beutel, A. "The Instruction Hierarchy: Training LLMs to Prioritize Privileged Instructions." *arXiv:2404.13208*, 2024. (OpenAI)
   - *Used in*: §4 (Instruction-Data Separation)
   - *Key finding*: Priority hierarchy (System > User > Tool Output); models can be trained to selectively ignore lower-privileged instructions even under adversarial attack.

**§5 — Related Systems:**

7. **[Wang et al. 2024]** Wang, X., et al. "OpenHands: An Open Platform for AI Software Developers as Generalist Agents." *arXiv*, 2024. (CMU et al.)
   - *Used in*: §5 (Sandboxed Loops representative)
   - *Key contribution*: Event Stream architecture, modular SDK for agent development.

8. **[Yang et al. 2024]** Yang, J., et al. "SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering." *arXiv*, 2024. (Princeton)
   - *Used in*: §5 (Sandboxed Loops representative)
   - *Key contribution*: Agent-Computer Interface (ACI) concept; optimized shell interactions for LLM context windows.

9. **[Mei et al. 2024]** Mei, K., et al. "AIOS: LLM Agent Operating System." *arXiv:2403.16971*, 2024. (COLM 2025)
   - *Used in*: §5 (Pseudo-OS Middleware)
   - *Key contribution*: OS-inspired abstractions for LLM agents (scheduler, memory manager), but implemented as application-layer middleware, not kernel-level mapping.
   - *Distinction*: All agents are Python objects within a single process; Quine maps each agent to a separate PID.

10. **[Li et al. 2026]** Li, C., Liu, X., Meng, X., & Zhao, X. "Architecting AgentOS: From Token-Level Context to Emergent System-Level Intelligence." *arXiv:2602.20934*, 2026.
    - *Used in*: §5 (Pseudo-OS Middleware)
    - *Key contribution*: Conceptual framework treating LLM as "Reasoning Kernel" with Deep Context Management and Semantic Slicing.
    - *Distinction*: Theoretical framework mapping OS abstractions onto LLM constructs; Quine maps agents *to* actual OS processes.

11. **[Zhang et al. 2025]** Zhang, C., Huang, H., Ni, C., et al. "UFO2: The Desktop AgentOS." *arXiv:2504.14603*, 2025. (Microsoft Research)
    - *Used in*: §5 (Pseudo-OS Middleware)
    - *Key contribution*: Multiagent AgentOS for Windows desktops with HostAgent/AppAgent hierarchy, hybrid UI detection, speculative multi-action planning.
    - *Distinction*: Desktop automation focus; agents coordinate within single control framework. Quine: each agent is an independent POSIX process.

**§6 — Beyond POSIX:**

12. **[Pike et al. 1995]** Pike, R., Presotto, D., Dorward, S., Flandrena, B., Thompson, K., Trickey, H., & Winterbottom, P. "Plan 9 from Bell Labs." *Computing Systems*, 8(3), 1995.
    - *Used in*: §6 (Plan 9 philosophy, 9P protocol, epigraph)

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
