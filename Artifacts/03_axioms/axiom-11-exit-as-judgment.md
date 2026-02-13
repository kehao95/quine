# Axiom XI: Exit-as-Judgment

> **Termination:** Exit Code $\equiv$ Judgment.

## 1. Definition

Termination is an explicit decision computed by the Agent, not a side effect of resource exhaustion. The Agent declares its terminal state via a standard POSIX integer exit code.

## 2. Invariants

*   **Explicit Verdict:** The exit code determines the process health, aligning with strict POSIX standards:
    *   `exit 0` = **Success** (Task Completed). The process executed its logic successfully and believes it fulfilled the intent.
    *   `exit > 0` = **Failure** (Task Aborted).
        *   `exit 1`: **General Error** (Runtime Error, Cognitive Failure).
        *   `exit 2`: **Misuse** (Invalid Intent/Arguments).
        *   `exit 124`: **Timeout** (`SIGTERM`).
        *   `exit 137`: **OOM Kill** (`SIGKILL`).
*   **The Single Source of Truth:** Judgment is compressed into this 8-bit integer. There is no separate "Status String" or "Meta Object".
    *   **Agent's Opinion:** Encoded in `exit code`.
    *   **Agent's Reasoning:** Encoded in `Tape`.
    *   **Agent's Output:** Encoded in `stdout`.
*   **Self-Authority:** The Agent is the sole authority over its own cognitive status, but the OS is the authority over its health.

## 3. Thermodynamic Halting

The system guarantees convergence to rest state.
1.  **Energy Injection:** Spawning child processes.
2.  **Energy Decay:** Each process consumes tokens and approaches its context limit.
3.  **Collapse:** Processes execute `exit()`, returning control and freeing resources.
4.  **Zero State:** When Root exits, system energy = 0.

Infinite loops are structurally impossible within the DAG topology + $D_{max}$ constraint.

## 4. The Coordination: Volition vs. Physics

A robust system cannot rely solely on the Agent's internal volition to terminate. Reliability requires the coordination of two forces:

### 1. Internal Volition (The Exit)
*   **Source:** The Agent's logic.
*   **Meaning:** "I have finished my reasoning."
*   **Semantics:** This is a **Judgment**. It carries semantic weight (Success/Failure/Partial).

### 2. External Physics (The Kill)
*   **Source:** The OS Kernel / Parent Process.
*   **Meaning:** "You have exhausted your resources (Time/RAM/Budget)."
*   **Semantics:** This is a **Constraint**. It carries no semantic weight regarding the task, only regarding the system state.

### Resolution Logic
The system's integrity relies on the race condition between Volition and Physics:
*   **Ideal State:** Agent `exits` *before* the Kernel `kills`. (Judgment is preserved).
*   **Failure State:** Kernel `kills` *before* Agent `exits`. (Judgment is lost; resource bankruptcy declared).

Therefore, **Exit-as-Judgment** is valid *if and only if* it occurs within the bounds of **OS-as-World** (Axiom I).
