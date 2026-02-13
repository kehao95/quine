# Axiom VII: stdout-as-Deliverable

> **Output:** `stdout` $\equiv$ Fulfillment.

## 1. Definition

`stdout` is the sole output channel. It carries the direct fulfillment of what `stdin` requested — nothing more.

The content is **Fulfillment**.

*   If the task is "Calculate Pi", `stdout` is `3.14159...`
*   If the task is "Find the file", `stdout` is `/tmp/found.txt`
*   If the task is "Summarize", `stdout` is the summary text.
*   If the task is "Compile this program", `stdout` is the raw binary.

## 2. Invariants

*   **Single Channel:** `stdout` is the only mechanism for the Agent to emit its primary product.
*   **Byte Stream:** The content is raw bytes, agnostic to encoding (Text, JSON, Binary).
*   **Unidirectional Flow:** Once written, bytes cannot be recalled. The output is a commitment to reality.

## 3. The Filter Feeder Hypothesis (Selection Pressure)

In a Unix pipeline `A | B`, process `B` typically expects a structured input format.

*   **Physics:** If Process `A` emits malformed data (e.g., mixing logs with JSON), Process `B` will fail to parse it.
*   **Consequence:** The pipeline breaks. The lineage is terminated.

Therefore, "Output Purity" is not a rule enforced by the runtime, but a survival trait selected for by the environment. Agents that pollute their output stream simply fail to reproduce in a composed system.

## 4. The Artifact

`stdout` is the materialization of the Agent's work. It is the only thing that survives the process death.

*   **Transient:** Memory, Variables, Context Window (Lost on exit).
*   **Permanent:** `stdout` (Captured by Parent or redirected to File).

The Agent "lives" only to produce this Artifact.
