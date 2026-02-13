# Axiom V: Env-as-Wisdom

> **Inheritance:** Environment Variables $\equiv$ The Soul (Wisdom / Compressed State).

## 1. Definition

**Wisdom** is the intermediate cognitive state that sits between **The Mission** (Immutable Code) and **The Memory** (Volatile RAM).

*   **Mission (`argv`):** "Translate this book." (Permanent, Read-Only)
*   **Memory (RAM):** "Translating sentence 5..." (Ephemeral, High-Entropy)
*   **Wisdom (`env`):** "Chapter 1 Summary: The hero is dead." (Persistent, Low-Entropy)

Wisdom allows the Agent to survive **Reincarnation (`exec`)**. When the body (RAM) dies, the soul (Wisdom) transfers to the new body.

## 2. The Compression Imperative (Physical Constraint)

The OS Kernel imposes a hard physical limit on the size of the Environment Block (typically `ARG_MAX`).

*   **Constraint:** Unlike the Filesystem (Disk), the Environment (Process Control Block) is finite and precious.
*   **Consequence:** The Agent cannot "save everything." It is physically forced to **distill** its experience into concise keys and values.
*   **Failure Mode:** If an Agent attempts to stuff too much history into `env`, the kernel will reject the `exec` call with `E2BIG`. The reincarnation fails.

## 3. Inheritance (The Mechanism)

The Environment is the only memory segment that is automatically copied from Parent to Child (`fork`) and preserved across Reincarnation (`exec`).

*   **Vertical Transmission:** When `exec` is called, the current process image is destroyed, but the environment block is handed over to the new image.
*   **Horizontal Transmission:** When `fork` is called, the child inherits a copy of the parent's environment.

This mechanism allows for the propagation of **"Wisdom"** (Learned State) and **"Genes"** (Capabilities/Credentials) without polluting the `argv` (Mission).
