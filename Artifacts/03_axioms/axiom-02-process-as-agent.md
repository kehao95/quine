# Axiom II: Process-as-Agent

> **Identity:** POSIX Process $\equiv$ Agent.

## 1. Definition

The Agent is a standalone binary executable operating within the host OS process table. It is not an object, a class instance, or a service. It is a discrete unit of execution governed by standard Unix process semantics.

## 2. Invariants

*   **State Isolation:** Memory is private. Inter-agent state pollution is impossible by OS design.
*   **Standard Interface:** The runtime interacts solely via standard streams (`stdin`, `stdout`, `stderr`) and signals (`SIGINT`, `SIGKILL`).
*   **Resource Bounded:** CPU and memory usage are constrained by OS primitives (`ulimit`, `cgroups`), not application logic.
*   **Ephemerality:** Processes are disposable.
*   **Side Effects:** Agents are processes, not pure functions. Like `rm`, `mv`, and `touch`, they are defined by their ability to mutate system state (filesystem, network, environment).

## 3. Cognitive Mapping

| Unix Primitive | Cognitive Equivalent |
| :--- | :--- |
| **Binary** | The `quine` runtime (inference engine) |
| **Process** | An active instance of thought |
| **PID** | Unique cognitive session identifier |
| **Scheduler** | The OS kernel (time-slicing execution) |

## 4. Operational Semantics

The `quine` binary treats the Large Language Model (LLM) as a **probabilistic text transducer**. It reads text, computes a semantic transformation, and writes text. Isomorphic to `sed` or `awk`, but operating on meaning rather than pattern.

## 5. Process Morphologies (The Ecological Niche)

Since the Agent is strictly a POSIX process, it inherits all standard Unix operational modes. It is not limited to a "REPL" loop; it adapts to the niche defined by its invocation method.

### Type A: The Periodic Observer (Cron Job)
*   **Unix Primitive:** `cron` / `systemd timers`
*   **Role:** **Time-Driven Immunity.** The Agent wakes at specific intervals (e.g., `@daily`), performs a bounded check (log scan, news digest), and `exits`.
*   **Dynamics:** Zero residual entropy. The Agent has no "memory burden" between cycles; it is a fresh, stateless antibody.

### Type B: The Reflex (Event-Driven)
*   **Unix Primitive:** `inotify` / `udev` / Signals
*   **Role:** **Event-Driven Intervention.** The Agent remains dormant (consuming 0 tokens) until a physical event occurs (file save, device plug-in).
*   **Dynamics:** $O(1)$ latency. Replaces "Polling" (Active Waiting) with "Interrupts" (Passive Waking).

### Type C: The Supervisor (PID 1)
*   **Unix Primitive:** `init` / `supervisord`
*   **Role:** **Meta-Cognition / Router.** An Agent whose `argv` is not to process data, but to manage the lifecycle of other Agents.
*   **Dynamics:** **Headless Autonomy.** It consumes tokens only to analyze `exit codes` and `stderr` streams of its children, deciding whether to restart or escalate.

### Type D: The Service (IPC Node)
*   **Unix Primitive:** Unix Domain Sockets (`.sock`)
*   **Role:** **Stateful Utility.** The Agent binds to a socket, offering specific cognitive services (e.g., "RAG Search", "Planner") to multiple local clients.
*   **Dynamics:** Allows for Non-DAG topologies (e.g., Star topology) where state needs to be shared or queried interactively.

### Type E: The Sandboxed Entity
*   **Unix Primitive:** `chroot` / Namespaces / cgroups
*   **Role:** **Restricted Reality.** The Agent operates in a manufactured reality where the filesystem is empty or read-only.
*   **Dynamics:** **Physical Safety > Moral Safety.** "Hallucinated" destruction (`rm -rf /`) is contained within the namespace, rendering semantic alignment checks redundant.

### Type F: The Specialist (LD_PRELOAD)
*   **Unix Primitive:** `LD_PRELOAD` / Shared Objects (`.so`)
*   **Role:** **Skill Injection.** Injecting specific cognitive libraries (prompts/tools) into a generic Agent at runtime.
*   **Dynamics:** Decouples "Personality" from "Runtime". Skills become hot-swappable libraries.
