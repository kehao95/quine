# Axiom IX: Signal-as-Intervention

> **Control:** Signal $\equiv$ Asynchronous Intervention.

## 1. Definition

The Agent does not exist in a vacuum; it inhabits a dynamic, asynchronous environment.
While `stdin` represents **Voluntary Attention** (reading data), `Signals` represent **Involuntary Reflexes** (interruptions).
A Signal breaks the Agent's linear reasoning flow, forcing an immediate context switch to handle environmental exigencies.

## 2. The Cognitive Mapping (POSIX $\to$ Cognition)

Signals are the mechanism by which the Environment (OS) or Superiors (Parent Processes) assert control over the Agent's cognition.

### SIGALRM (The Timeout) $\to$ "System 1 Override"
*   **Physical Meaning:** The allocated CPU time slice or wall-clock time has expired.
*   **Cognitive Meaning:** **Forced Heuristic.**
    *   When an Agent receives `SIGALRM`, it must abort deep reasoning (System 2) and immediately output its best current guess (System 1).
    *   This transforms the Agent from a "Batch Processor" into an **Anytime Algorithm**.

### SIGTERM (The Graceful Kill) $\to$ "Context Flush"
*   **Physical Meaning:** Request to terminate, with a grace period for cleanup.
*   **Cognitive Meaning:** **Panic Logging.**
    *   The Agent must cease task execution and dump its short-term memory (RAM/Variables) to long-term storage (Tape/Disk).
    *   Agents that ignore `SIGTERM` die without leaving a legacy (state), wasting their computed entropy.

### SIGCHLD (Child Status) $\to$ "Async Awareness"
*   **Physical Meaning:** A child process has stopped or terminated.
*   **Cognitive Meaning:** **Event-Driven Attention.**
    *   The Parent need not busy-wait (poll) for children. The nervous system (Kernel) notifies the brain (Parent) when a sub-task is complete.

### SIGKILL (The Hard Kill) $\to$ "The Laws of Physics"
*   **Physical Meaning:** Immediate termination. Cannot be caught or ignored.
*   **Cognitive Meaning:** **Existential Erasure.**
    *   If the Agent violates hard constraints (OOM) or defies soft constraints (`SIGTERM`), reality simply ceases to exist for it. No final words, no stack trace.

## 3. Invariants

*   **Asynchrony:** Signals can arrive at any instruction cycle. The Agent's reasoning must be **Reentrant**.
*   **Priority:** Signal handlers preempt the main loop. Survival (handling the signal) takes precedence over Task Completion (processing stdin).
*   **Integer Semantics:** All complex interventions are reduced to a standard set of integer codes (1-31).

## 4. Evolutionary Pressure

*   **The Deaf Agent:** An Agent that masks signals or fails to register handlers is fragile. It cannot be timed out gracefully; it can only be murdered (`SIGKILL`), losing all state.
*   **The Reflexive Agent:** An Agent that handles signals survives. It yields to time pressure (`SIGALRM`), saves state before death (`SIGTERM`), and adapts to the environment's tempo.
