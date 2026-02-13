# Axiom X: Recursion-as-Metabolism

> **Metabolism:** Strategy $\equiv$ Process Primitives (`fork` vs `exec`).

## 1. Definition

Complex tasks are solved not by "thinking harder" (adding parameters), but by manipulating the process tree. The Agent uses two fundamental POSIX primitives to manage its cognitive state: **Fork (Exploration)** and **Exec (Purification)**.

## 2. The Physical Constraint: Entropy

We do not program the Agent to use these primitives as algorithmic preferences. Instead, they emerge as adaptations to a hard physical constraint: **The Monotonic Increase of Heap Entropy**.

*   **Entropy Accumulation:** As an Agent reasons, its `.data` segment (Tape) fills with "thought waste".
*   **Cognitive Coherence Limit:** When the Heap fills the Context Window, the Code Segment (`.text`) is pushed out of attention, and the Agent hallucinates.

## 3. The Two Primitives

The Agent navigates the Entropy Limit using two fundamental system calls:

### Primitive A: Fork (Horizontal Scaling)
*   **Physics:** `fork()` creates a clone of the process with a copy of the memory.
*   **Cognitive Function:** **Parallel Exploration.** The Parent delegates a sub-problem to a Child. The Child's entropy (context accumulation) is isolated from the Parent.

### Primitive B: Exec (Vertical Scaling)
*   **Physics:** `exec()` replaces the process image. The Heap (Context) is destroyed; the Mission (`argv`) is reloaded; the Wisdom (`env`) is preserved.
*   **Cognitive Function:** **Reincarnation.** The Agent sheds its "Body" (polluted RAM) to save its "Soul" (Wisdom). This is the only mechanism to process infinite streams with finite memory.

## 4. The Thermodynamic Necessity

Recursion is not an algorithmic choice; it is a metabolic necessity.

*   **Finite Context:** No single process can live forever (it will run out of tokens/memory).
*   **Infinite Tasks:** Some tasks exceed the lifespan of a single process.
*   **Conclusion:** The only way to complete an infinite task with finite life is to reproduce.

## 5. Invariants

*   **DAG Topology:** `fork()` creates a strictly hierarchical process tree.
*   **Information Compression:** Survival depends on discarding high-entropy state (RAM) while preserving low-entropy state (Env/Disk).
