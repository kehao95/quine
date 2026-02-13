# Project Quine: Autopoietic Intelligence via Fractal POSIX Processes

> The complexity of AI Agents is not inherent to intelligence, but a product of deviating from fundamental computational principles. The solution: redefine Agents as POSIX-compliant binary processes.

## Document Index

The `Artifacts/` directory serves as the **Single Source of Truth** for the project's definitions and philosophy.

| Sequence | Artifact | Content Description |
| :--- | :--- | :--- |
| **00** | **[00_abstract.md](./00_abstract.md)** | Core thesis statement. |
| **01** | **[01_manifesto.md](./01_manifesto.md)** | **The Philosophy.** Motivation, misconceptions, evolutionary dynamics, comparative analysis. |
| **02** | **[02_genealogy.md](./02_genealogy.md)** | **The Roots.** Historical synthesis of Unix Philosophy (1969), Cognitive Architectures (2023), and Thermodynamics of Intelligence (2025). |
| **03** | **[03_axioms/](./03_axioms/)** | **The Physics.** Formal definitions of the 12 axioms and 13 system guarantees. |
| **04** | **[04_implementation.md](./04_implementation.md)** | **The System.** Kernel architecture, host-guest separation, lifecycle management. |
| **05** | **[05_case_study.md](./05_case_study.md)** | **The Evidence.** Three case studies: Specification as Physics (jq), Scale Invariance (MRCR), and Autopoiesis (Binary Quine). |
| **06** | **[06_conclusion.md](./06_conclusion.md)** | **The Future.** Thermodynamic analysis and final synthesis. |
| **07** | **[07_acknowledgement.md](./07_acknowledgement.md)** | Acknowledgments. |

## Reading Order (The Monograph Structure)

The documentation is structured to be read as a coherent monograph, moving from philosophy to physics, then to engineering and evidence.

### Part I: The Philosophy
*   **[Abstract](./00_abstract.md)** — Core thesis: recursive cognitive architecture mapping POSIX primitives to cognition.
*   **[Manifesto](./01_manifesto.md)** — The ontological shift: redefining Agents as POSIX processes. Seven misconceptions. Evolutionary dynamics. Comparative analysis.

### Part II: The Roots (History)
*   **[Genealogy](./02_genealogy.md)** — Four historical threads:
    1. Unix Pipes & Filters (1964)
    2. Cognitive Architectures: ToT & Reflexion (2023)
    3. OS-LLM Analogy Research: MemGPT, AIOS
    4. Thermodynamics of Intelligence (2025): Context Rot, Detailed Balance, Maxwell's Demon

### Part III: The Physics (Axioms)
*   **[Layer 1: Ontology](./03_axioms/axiom-01-os-as-world.md)** — The world, the agent, and the medium (Axioms I–III).
*   **[Layer 2: Interface](./03_axioms/axiom-04-argv-as-mission.md)** — The I/O membrane (Axioms IV–VIII): argv, env, RAM, stdout, stderr.
*   **[Layer 3: Dynamics](./03_axioms/axiom-09-signal-as-intervention.md)** — Control & lifecycle (Axioms IX–XII): Signal, Recursion, Exit, Scarcity.
*   **[Guarantees](./03_axioms/guarantees.md)** — Thirteen system properties derived from these axioms.

### Part IV: The System (Implementation)
*   **[Implementation](./04_implementation.md)** — The `quine` kernel architecture:
    *   Host-Guest separation (Go runtime vs LLM cognition)
    *   Quad-Channel Protocol (argv, stdin, stdout, stderr)
    *   Four tools: `sh`, `fork`, `exec`, `exit`
    *   Process management: Mitosis (fork) and Metamorphosis (exec)
    *   The Autopoietic Log (Tape)
    *   Stateless Iterator Pattern (Iterative Amnesia)

### Part V: The Evidence (Case Studies)
*   **[Case Studies](./05_case_study.md)** — Three experiments validating the architecture:
    *   **A: Specification as Physics (jq)** — System Prompt as executable environmental physics.
    *   **B: Scale Invariance (MRCR)** — Recursive reincarnation decouples retrieval from context limits.
    *   **C: Autopoiesis (Binary Quine)** — Self-replication from System Prompt alone.

### Part VI: The Future (Conclusion)
*   **[Conclusion](./06_conclusion.md)** — Emergence over Instruction. The Synthetic Data Engine. Computational Thermodynamics.

---

## Axiom Index

The twelve axioms form a three-layer architecture, moving from physical laws to cognitive dynamics.

### Layer 1: Ontology — The world, the agent, and the medium

| # | Axiom | Unix Primitive | Cognitive Function |
|---|-------|---------------|--------------------|
| I | [OS-as-World](./03_axioms/axiom-01-os-as-world.md) | OS kernel | Physical laws of the environment |
| II | [Process-as-Agent](./03_axioms/axiom-02-process-as-agent.md) | process | Agent identity & isolation |
| III | [Filesystem-as-State](./03_axioms/axiom-03-filesystem-as-state.md) | filesystem | Persistent state across generations |

### Layer 2: Interface — The agent's I/O contract and memory model

| # | Axiom | Unix Primitive | Cognitive Function |
|---|-------|---------------|--------------------|
| IV | [argv-as-Mission](./03_axioms/axiom-04-argv-as-mission.md) | argv | Immutable mission / Harvard Architecture |
| V | [Env-as-Wisdom](./03_axioms/axiom-05-env-as-wisdom.md) | env | State transfer / compressed knowledge |
| VI | [RAM-as-Context](./03_axioms/axiom-06-ram-as-context.md) | RAM / Heap | Volatile working memory / Entropy buffer |
| VII | [stdout-as-Deliverable](./03_axioms/axiom-07-stdout-as-deliverable.md) | stdout | Output / fulfillment |
| VIII | [stderr-as-Gradient](./03_axioms/axiom-08-stderr-as-gradient.md) | stderr | Feedback / semantic gradient |

### Layer 3: Dynamics — Control, lifecycle, and evolutionary pressure

| # | Axiom | Unix Primitive | Cognitive Function |
|---|-------|---------------|--------------------|
| IX | [Signal-as-Intervention](./03_axioms/axiom-09-signal-as-intervention.md) | signal | Asynchronous intervention / interrupt |
| X | [Recursion-as-Metabolism](./03_axioms/axiom-10-recursion-as-metabolism.md) | fork, exec | Cognitive decomposition / entropy management |
| XI | [Exit-as-Judgment](./03_axioms/axiom-11-exit-as-judgment.md) | exit code | Termination / verdict |
| XII | [Scarcity-as-Selection](./03_axioms/axiom-12-scarcity-as-selection.md) | ulimit, cgroups | Evolutionary pressure |

---

## The Thirteen Guarantees

Properties inherited by mapping Agentic behavior to POSIX standards. See **[guarantees.md](./03_axioms/guarantees.md)** for details.

| Layer | Guarantees |
|-------|------------|
| **Ontology** | Hard Isolation, Capability Minimalism, Deterministic Context |
| **Interface** | Prompt Injection Immunity*, Universal Composability, Semantic Backpropagation |
| **Memory** | Wisdom Inheritance, Entropy Reset, Forensic Auditability |
| **Dynamics** | Graceful Degradation, Thermodynamic Halting, Fractal Scale-Invariance, Selection over Design |

\* *Conditional guarantee: requires LLM support to respect the authority boundary.*
