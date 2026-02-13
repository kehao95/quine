# Axiom IV: argv-as-Mission

> **Input:** `argv` $\equiv$ Code Segment (`.text`). `stdin` $\equiv$ Data Stream.

## 1. Definition

The Agent is not a chatbot; it is a **Harvard Architecture Machine**.

To separate **Instruction** (what to do) from **Information** (what to process), we map POSIX primitives to memory segments:

*   **System Prompt + `argv`:** The **Code Segment (`.text`)**. Immutable, Read-Only. Defined at process start.
*   **`stdin`:** The **Data Stream**. Read-Only. The raw material to be processed.

> For the volatile working memory (Context Window), see [Axiom VI: RAM-as-Context](./axiom-06-ram-as-context.md).

## 2. The Cognitive Mapping

| Component | Unix Concept | Memory Segment | Cognitive Function |
| :--- | :--- | :--- | :--- |
| **`argv`** | Program Arguments | **`.text` (Code)** | **The Mission.** Immutable goal. |
| **`stdin`** | File Descriptor 0 | **Data Bus** | **The Raw Material.** Input to process. |

## 3. Mission vs. Intent

We distinguish two related but distinct concepts:

| Term | Perspective | Semantics |
| :--- | :--- | :--- |
| **Mission** | Process-internal | The immutable goal for the lifetime of *this* process. Once started, `argv` cannot change. |
| **Intent** | Parent → Child | The directive passed when spawning a child. From the parent's view, it is a "next step" — one phase in a larger plan. |

When a Parent calls `fork(intent="subtask")`:
*   The **Intent** is what the Parent *sends*.
*   Once the Child starts, that Intent becomes the Child's **Mission**.

This distinction matters because:
1.  **Mission** emphasizes immutability and authority (the process cannot question its own `argv`).
2.  **Intent** emphasizes decomposition and planning (the parent breaks work into sub-goals).

## 4. Invariants

*   **Mission Immutability:** The Agent *cannot* modify its own `argv`. The Mission is fixed for the lifetime of the process.
*   **Code/Data Isolation:** Data from `stdin` cannot overwrite `argv` (Code).
*   **`exec` Semantics:** When `exec` is called:
    *   **Mission (`argv`) is preserved** (unless explicitly changed).
    *   **Context Window is wiped** (see [Axiom VI](./axiom-06-ram-as-context.md)).

## 5. Security: The Harvard Architecture

By physically separating the **Mission Bus** (`argv`) from the **Data Bus** (`stdin`), Quine implements a cognitive **Harvard Architecture**.

*   **Von Neumann Failure (Traditional Agents):** User data is concatenated with system prompts. "Ignore previous instructions" works because they live in the same address space.
*   **Harvard Solution (Quine):** The Agent knows that **Authority** comes *only* from `argv`. The `stdin` is just "what I am reading." It cannot countermand the Mission.

## 6. Composability

Because the Agent accepts standard input (`stdin`) and produces standard output (`stdout`), it satisfies the **Universal Interface** contract.

This allows the Agent to be placed in any pipeline position:
*   **As a Filter:** Transforming data stream A to stream B.
*   **As a Source:** Generating data from an internal model.
*   **As a Sink:** Consuming data to perform side effects.

The OS guarantees that `argv` (the logic) remains isolated from the data stream, preventing instruction injection.
