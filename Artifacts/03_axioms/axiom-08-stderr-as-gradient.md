# Axiom VIII: stderr-as-Gradient

> **Feedback:** `stderr` $\equiv$ Semantic Gradient.

## 1. Definition

`stderr` is the back-propagation channel for **Optimization Gradients** ($\nabla L$). It is **not** a general-purpose log stream. It is the channel for **Negative Feedback**.

## 2. Invariants

*   **Silence is Success:** If the Agent perfectly executes its mission (`L=0`), `stderr` MUST remain silent. Any output on `stderr` represents a **Non-Zero Gradient**—an indication of struggle, ambiguity, or failure.
*   **Channel Separation:** `stderr` is physically distinct from `stdout`.
*   **Signal Purity:** General "thoughts", "reasoning traces", and "info logs" belong to the **Audit Log (Tape)**. `stderr` is reserved strictly for **Error Signals** and **Correction Requests** required for the parent to optimize the prompt.
*   **Upstream Propagation:** While `stdout` flows downstream (to the next consumer), `stderr` flows upstream (to the Supervisor/Parent).

## 3. The Evolutionary Pressure (Selection)

The content of `stderr` is not enforced by the kernel. However, its utility is determined by the Parent Process.

*   **Case A (Noise):** If a Child fills `stderr` with useless information, the Parent cannot distinguish signal from noise. The Child is likely to be terminated for inefficiency.
*   **Case B (Signal):** If a Child emits precise error traces, the Parent can use this information to adjust the parameters for the next run (Self-Correction).

Therefore, "Gradient Fidelity" (useful error messages) is a survival trait. Agents that communicate their failure modes effectively are more likely to be "fixed" (re-run with better prompts) rather than discarded.

## 4. Emergent Protocol

Because of this selection pressure, semantic protocols naturally emerge. Agents evolve to output high-fidelity reasoning traces on `stderr` because it increases their systemic utility, maximizing their persistence in the process tree.

## 5. The Gradient Vector ($\nabla L$)

In this evolutionary context, the gradient is defined *post-hoc*:

$$ \nabla L \equiv \text{Information that reduces Parent Entropy} $$

If the content of `stderr` allows the Parent to reduce the search space for the solution, it is a valid gradient. If not, it is waste heat.

## 6. Signal Type

Unlike traditional software which emits structured logs, the Quine Agent emits **reasoning traces**. The parent LLM consumes this trace to "understand" the failure mode, enabling dynamic self-correction without hard-coded error handling logic.
