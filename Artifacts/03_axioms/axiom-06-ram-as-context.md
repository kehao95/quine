# Axiom VI: RAM-as-Context

> **Memory:** Volatile RAM (Heap) $\equiv$ Context Window.

## 1. Definition

The **Context Window** is not a "Knowledge Base" or a "Long-Term Memory." It is the **Volatile RAM (Heap)** of the Agent Process.

It is the active workspace where "Thinking" occurs. Like physical RAM, it is fast, expensive, and **ephemeral**.

## 2. The Law of Cognitive Entropy

The fundamental physical law governing the Context Window is **Monotonic Entropy Increase**.

*   **Accumulation:** Every token generated is a token added to the heap.
*   **Pollution:** As the conversation ("Stream of Consciousness") extends, the ratio of Signal (Mission) to Noise (Past Reasoning) decreases.
*   **Irreversibility:** You cannot "un-think" a thought. Once in the Context, it occupies attention.

## 3. The Memory Leak Hypothesis

Current "Long-Context" architectures treat the Context Window as a database (filling it with 100k+ tokens). In the Quine framework, this is defined as a **Memory Leak**.

*   **Symptom:** As Context grows, inference latency increases linear-quadratically, and reasoning capability degrades (the "Lost in the Middle" phenomenon).
*   **Diagnosis:** The Agent is failing to free memory.
*   **Cure:** The only way to free the Context Window is to terminate the process (`exit`) or replace the process image (`exec`).

## 4. RAM vs. Disk (The Separation)

We explicitly separate the Agent's memory into two distinct physical tiers:

| Tier | Component | Medium | Physics | Lifecycle |
| :--- | :--- | :--- | :--- | :--- |
| **L1** | **Context Window** | **RAM** | Volatile, Fast, Toxic | Wiped on `exec` |
| **L2** | **Filesystem** | **Disk** | Persistent, Shared | Persists across `exec` |

## 5. The "Detox" Mechanism

Because Context is toxic (accumulates entropy), the Agent must implement a **Detox Cycle**:

1.  **Serialize:** Write essential state to **Disk**.
2.  **Flush:** Call `exec`. This effectively "powers down" the RAM, instantly vaporizing all accumulated hallucinations, circular logic, and distraction.
3.  **Reload:** The new process starts with **Zero Entropy** in RAM, but full access to persisted state on Disk.

## 6. The Cognitive OOM Killer

Just as the OS Kernel kills processes that consume too much physical RAM (OOM Killer), the Quine Runtime kills Agents that consume too much Context (Tokens).

*   **Threshold:** Defined by `MAX_TOKENS` or budget constraints.
*   **Action:** Immediate `SIGKILL`.
*   **Philosophy:** A "Memory Leaking" Agent (one that fails to `exec` or `exit` before filling its window) is a **Defective Process**. It is not "guided" to stop; it is **terminated**.
*   **Evolutionary Consequence:** Only Agents that develop "Hygiene" (regular `exec` / `exit` before the limit is reached) survive to complete long tasks.
