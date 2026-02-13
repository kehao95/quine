# Axiom I: OS-as-World

> **Environment:** The Operating System $\equiv$ The Physical World.

## 1. Definition

The Agent does not inhabit an abstract "chat room" or a Python interpreter. It inhabits the **Operating System**. The OS Kernel defines the physical laws of its universe, and the Shell provides the interface to manipulate it.

## 2. Invariants

*   **Sensors & Actuators:** The Agent perceives reality *only* through system calls (via CLI tools like `ls`, `cat`, `curl`, `ps`) and changes reality *only* through executable binaries (`mkdir`, `rm`, `python`, `gcc`).
*   **Laws of Physics:** The Kernel enforces immutable constraints:
    *   **Permissions:** You cannot read what you don't own (`EACCES`).
    *   **Resources:** You cannot consume more than your quota (`ENOMEM`, `ENOSPC`).
    *   **Time:** You cannot stop the clock (`SIGALRM`, `timeout`).
*   **No Magic:** The Agent cannot "hallucinate" a side effect. It cannot simply *say* "I have saved the file." It must execute `echo "content" > file.txt`. If the syscall fails, the action did not happen.

## 3. The Universal Interface

Because the OS provides a unified interface for Computing (Process), Storage (Filesystem), and Networking (Sockets), the Agent can interact with *any* system that exposes a POSIX-compliant interface.

*   **Local:** Manipulate files, processes, devices.
*   **Remote:** `ssh`, `curl` extend the agent's reach to other machines.
*   **Code:** `python`, `node`, `go` extend the agent's capability to any language runtime.

The OS is not just a sandbox; it is the **Universal Adapter** connecting the Agent to the digital world.

## 4. Implication for Benchmark

This axiom implies that a valid benchmark must test the Agent's ability to **navigate the OS**.
*   Can it find a file without knowing the path (`find`)?
*   Can it kill a runaway process (`ps`, `kill`)?
*   Can it recover from a "Disk Full" error (`df`, `rm`)?

Intelligence in Quine is defined as **OS Literacy**.
