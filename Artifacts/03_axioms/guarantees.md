# The Thirteen System Guarantees

> Properties inherited by mapping Agentic behavior to POSIX standards.

These guarantees are not implemented features; they are emergent properties of the twelve axioms, structured across the same four layers: **Ontology**, **Interface**, **Memory**, and **Dynamics**.

## Layer 1: Ontological Guarantees
*Properties derived from the immutable laws of the OS environment.*

| # | Guarantee | Principle | Derived From |
|---|-----------|-----------|--------------|
| 1 | **Hard Isolation** | Trust the Kernel. A child process crash is strictly isolated to its own memory space. | Axiom I (OS-as-World), II (Process-as-Agent) |
| 2 | **Capability Minimalism** | Least Privilege. Safety is a physical constraint (permissions, namespaces), not a polite request. | Axiom I (OS-as-World) |
| 3 | **Deterministic Context** | Structure is Memory. Information has a deterministic address (filepath), not a probabilistic embedding. | Axiom III (Filesystem-as-State) |

## Layer 2: Interface Guarantees
*Properties derived from the agent's I/O contract.*

| # | Guarantee | Principle | Derived From |
|---|-----------|-----------|--------------|
| 4 | **Prompt Injection Immunity**\* | Harvard Architecture. Data (`stdin`) cannot overwrite Instructions (`argv`). Authority comes *only* from the Code Segment. | Axiom IV (argv-as-Mission) |

\* *Conditional guarantee: Quine provides the architectural separation; actual immunity requires LLM support to respect the authority boundary.*
| 5 | **Universal Composability** | Streams are Universal. An agent applies logic (`argv`) to any data stream (`stdin`). Pipelines work because downstream treats upstream output as raw material, not instruction. | Axiom IV (argv-as-Mission), VII (stdout-as-Deliverable) |
| 6 | **Semantic Backpropagation** | Stderr is the Gradient. The error stream enables optimization without inspecting internal state. | Axiom VIII (stderr-as-Gradient) |

## Layer 3: Memory Guarantees
*Properties derived from the agent's three-tier storage model.*

| # | Guarantee | Principle | Derived From |
|---|-----------|-----------|--------------|
| 7 | **Wisdom Inheritance** | Soul Transfer. Environment variables survive `exec`, enabling compressed knowledge to outlive any single process body. | Axiom V (Env-as-Wisdom) |
| 8 | **Entropy Reset** | Cognitive Garbage Collection. `exec` guarantees complete context purification — accumulated hallucinations are physically destroyed. | Axiom VI (RAM-as-Context), X (Recursion-as-Metabolism) |
| 9 | **Forensic Auditability** | No Hidden States. Reasoning is serialized to persistent storage. We debug the thought process, not the neural weights. | Axiom III (Filesystem-as-State) |

## Layer 4: Dynamic Guarantees
*Properties derived from control, lifecycle, and evolutionary pressures.*

| # | Guarantee | Principle | Derived From |
|---|-----------|-----------|--------------|
| 10 | **Graceful Degradation** | Anytime Algorithm. Under time pressure (`SIGALRM`), the agent yields best-effort output rather than failing silently. | Axiom IX (Signal-as-Intervention) |
| 11 | **Thermodynamic Halting** | Energy State → Zero. A process tree with bounded depth and finite budget is guaranteed to halt. | Axiom XI (Exit-as-Judgment), XII (Scarcity-as-Selection) |
| 12 | **Fractal Scale-Invariance** | The Part Contains the Whole. A single agent fixing a typo uses the same mechanism as a swarm designing an OS. | Axiom X (Recursion-as-Metabolism) |
| 13 | **Selection over Design** | Correctness by Survival. Effective behaviors are not programmed but selected by resource scarcity. The environment is the spec. | Axiom XII (Scarcity-as-Selection) |
