# Axiom III: Filesystem-as-State

> **Persistence:** Process Memory (RAM) $\perp$ Time. Filesystem (Disk) $\parallel$ Time.

## 1. The Volatility Constraint
The Agent is a standard OS Process.
Upon `exit()`, the Kernel reclaims all allocated memory (Heap, Stack).
Therefore, **RAM is not a storage medium** for this architecture. It is purely a computation buffer.

## 2. The Single Source of Truth
For intelligence to persist across process generations ($P_t \to P_{t+1}$), state must be serialized to a non-volatile medium.
In Quine, the **Filesystem is the Only State Machine**.

*   **Concurrency:** The filesystem is a shared resource. Multiple processes can access it simultaneously.
*   **Race Conditions:** If two Agents write to the same file without coordination, data corruption occurs. This is a physical law, not a policy.
*   **Persistence:** Files outlive the process that created them.
*   **Truth:** What is written to disk is the only objective reality. Memory (RAM) is subjective and vanishes on death.

## 3. The Shared Reality
The Filesystem acts as the **Universal Blackboard**.
*   **Decoupling:** Agents communicate by modifying the shared environment (files), not by direct message passing.
*   **Observation:** To know what happened, an Agent must read the file. To let others know, it must write a file.
*   **History:** The filesystem accumulates the history of all past actions. It is the fossil record of cognition.

## 4. Environment as Database
We do not define *where* state is stored (e.g., specific folders).
The OS defines the physics of storage:
*   **Permissions:** Can I write here?
*   **Capacity:** Is the disk full?
*   **Persistence:** Will this inode survive a reboot?

The Agent must navigate these physical constraints to preserve its information.
