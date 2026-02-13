# System Implementation

The `quine` binary is a **POSIX-compliant runtime** for Large Language Models. It serves as a hypervisor, translating probabilistic token streams into deterministic system calls.

## 1. The Host-Guest Architecture

The runtime enforces a strict separation between the **Host** (Physics) and the **Guest** (Cognition).

*   **The Host (Go):** Deterministic. Handles file I/O, process spawning, signal handling, and resource quotas.
*   **The Guest (LLM):** Probabilistic. Resides in the context window. Predicts the next tool execution based on the Tape.

**The binary contains zero business logic.** It does not know what a "task" is. It is simply a read-eval-print loop for an LLM.

### 1.1 The Execution Cycle

The runtime implements a strict turn loop. The execution trace (Tape) is built incrementally in memory for context construction; the JSONL file on disk is a write-only audit log.

1.  **Infer:** Serialize in-memory conversation history to Context → Call LLM API.
2.  **Execute:** Run the syscall (e.g., `sh`).
3.  **Persist:** Append result to history (RAM) and Audit Log (disk, `fsync`).

## 2. The I/O Contract

Quine enforces a **Quad-Channel Protocol** to decouple Intent, Data, Execution, and Observability.

| Channel | Stream | Content | Consumer |
| :--- | :--- | :--- | :--- |
| **Intent** | `argv` | Instructions / Prompt (Text) | `quine` Runtime (System Prompt) |
| **Material** | `stdin` | Data Stream (Bytes) | Agent (via `cat /dev/stdin` in `sh`) |
| **Deliverable** | `stdout` | Pure Output (Bytes, via `>&3`) | Downstream Process (`\|`) |
| **Signal** | `stderr` | Failure Gradient (Post-Mortem) | Parent Process / Supervisor |
| **Audit** | `Tape` | Full Event Log (JSONL) | Human / UI / Debugger (Forensics) |

### 2.1 The Four Tools

The Agent has four specialized tools, mapping directly to POSIX primitives:

| Tool | Purpose | POSIX Equivalent | Use Case |
| :--- | :--- | :--- | :--- |
| **`sh`** | Execute command | `system()` | Interacting with the OS (File I/O, Network, stdin/stdout) |
| **`fork`** | Horizontal Scaling | `fork()` + `exec()` | Exploration (Spawns Child with **Cloned Context** + New Intent) |
| **`exec`** | Vertical Scaling | `exec()` | Detox (Replaces Self with **Empty Context** + New Intent) |
| **`exit`** | Terminate | `exit()` | Judgment / Completion |

### 2.2 The fd 3 Stdout Mechanism

Each `sh` command spawns an isolated shell process. The child process has access to:

*   **fd 0 (stdin):** The process's real stdin, wired via `cmd.Stdin`. The agent reads piped input data with `cat /dev/stdin`.
*   **fd 1 (stdout):** Captured into a buffer. The content appears in the tool result for the agent's context window.
*   **fd 2 (stderr):** Captured into a buffer. Also appears in the tool result.
*   **fd 3 (real stdout):** The process's real stdout, passed via `cmd.ExtraFiles[0]`. The agent writes deliverables with `echo "result" >&3` or `cat file.txt >&3`.

This separation serves two purposes:

1.  **Context Visibility:** Regular command output (fd 1) is captured and returned to the agent, so it can observe what happened.
2.  **Output Purity:** Deliverables written to fd 3 flow directly to the parent process's stdout without polluting the context window. This enables binary output, large file streaming, and pipeline composition.

## 3. Process Management

The Agent manages its cognitive lifecycle using two POSIX primitives: `fork` (Scaling) and `exec` (Renewal).

### 3.1 Mitosis (Fork) — "I need help"
Used for **Horizontal Scaling** and **Exploration**.

*   **Action:** The Agent calls `fork` with `intent="subtask"` and `wait=true/false`.
*   **Mechanism:**
    1.  Host flushes current Tape to disk.
    2.  Host copies `tape.jsonl` → `child_tape.jsonl` (Full Context Clone).
    3.  Host spawns child process pointed at `child_tape`.
    4.  If `wait=true`, Parent blocks until Child exits. If `wait=false`, Parent continues immediately (Fire-and-Forget).
*   **State:** Child starts with **Parent's Memories** + **New Intent**.
*   **Use Case:** "Based on my previous findings (Context), please investigate X (New Intent) and tell me the result (Wait)."

### 3.2 Metamorphosis (Exec) — "I need a fresh brain"
Used for **Vertical Scaling** and **Context Detox**.

*   **Action:** The Agent calls `exec` with `context="wisdom_summary"`.
*   **Mechanism:**
    1.  Host replaces the current process image (Same PID).
    2.  **Mission is Preserved** (Original `argv` is retained).
    3.  **Tape is Reset** (Empty file).
    4.  **Context is Reset** (Zero Entropy).
    5.  `context` arg is injected to bootstrap the new instance.
*   **Continuity:** Inherits file descriptors (stdin remains available to the new process image).
*   **Use Case:** "I have finished processing this chunk. Reincarnate me with a clean brain (but keep the wisdom) to process the next chunk."

### 3.3 Session Lineage

Each quine process has a unique `SESSION_ID`. Child processes do **not** inherit the parent's `SESSION_ID` — they generate their own. This ensures:

*   Each process writes to its own Tape file.
*   Multiple children spawned from a single `sh` call (e.g., via `&`) have distinct identities.
*   Lineage is tracked via `PARENT_SESSION` environment variable.

### 3.4 Signal Handling (The Nervous System)

> *Design in progress — details to be finalized.*

The runtime maps OS signals to cognitive interventions. The goal is to ensure **pipeline purity**: a process must either deliver a complete result or signal failure cleanly.

## 4. The Autopoietic Log (The Tape)

While **State** is Filesystem-Native (Axiom III), the runtime maintains a separate **Audit Log** (the "Tape") for forensic analysis. This is an append-only JSONL file located at `${QUINE_DATA_DIR}/${SESSION_ID}.jsonl`.

*   **Content:** It records the exact sequence of Inputs (`stdin`), Thoughts (Reasoning Traces), Actions (`syscalls`), and Outcomes (`stdout/stderr`).
*   **Purpose:** It is not the primary mechanism for agent state persistence (which is the Filesystem), but the source of truth for **Process Supervision**.
*   **Artifact:** This log serves as the "Experience Replay" buffer for training future models (Synthetic Data) and debugging current ones (Forensics).

## 5. The Stateless Iterator Pattern (Iterative Amnesia)

To process infinite data streams with a finite context window, Quine implements an **Iterative Amnesia** pattern.

### 5.1 The Tri-Zone Context (Harvard Architecture)

The runtime virtually segments the Context Window into three zones based on memory segment analogies:

| Zone | Content | Memory Segment | Persistence |
| :--- | :--- | :--- | :--- |
| **Zone A: The Law** | System Prompt + Mission (`argv`) | **.text (Code)** | **Immutable** (Preserved across `exec`) |
| **Zone B: The Wisdom** | Learned Insights (`env`) | **Registers / Stack** | **Transferred** (Passed via `exec`) |
| **Zone C: The Heap** | Conversation History | **.data (Heap)** | **Disposable** (Wiped on `exec`) |

### 5.2 The "Reincarnation Protocol" Modes

The Agent chooses between two modes of operation based on task interdependence:

#### Mode A: Independent (The Factory Worker)
> **"Each item is a new world."**
*   **Scenario:** Translation, Image Processing, Data Normalization.
*   **Action:** `exec ./quine "$@"` (No env changes).
*   **Result:** Zero context transfer. The next instance starts fresh, knowing only the System Prompt (`argv`).
*   **Entropy:** ΔE = 0 (Perfect reset).

#### Mode B: Sequential (The Historian)
> **"I need to remember what I just read."**
*   **Scenario:** Reading a novel, parsing multi-line logs, stateful protocol handling.
*   **Action:**
    1. `export QUINE_WISDOM_SUMMARY="Summary of Chapter 1"`
    2. `exec ./quine "$@"`
*   **Result:** Wisdom transfer. The next instance starts with the System Prompt + The `QUINE_WISDOM_*` variables.
*   **Entropy:** ΔE = len(QUINE_WISDOM_*) (Controlled growth).

### 5.3 The Wisdom Namespace Protocol

The runtime enforces a strict naming convention for state transfer to prevent environment pollution and ensure structured memory.

1.  **Write:** The Agent uses `export QUINE_WISDOM_KEY="Value"` to save a memory fragment.
2.  **Read:** The Runtime injects these variables into the System Prompt of the next incarnation.
3.  **Persistence:** Standard POSIX `exec` semantics ensure these variables survive the process replacement.

This turns the OS Environment into a **Key-Value Store** for the Agent's cognitive state.

### 5.4 Benefits

1.  **O(1) Latency:** Processing line 1,000,000 is as fast as line 1.
2.  **Security:** Malicious data in Chunk N is destroyed before reading Chunk N+1. It cannot "poison" the agent's long-term memory because Zone C is never persisted.
3.  **No Drift:** Hallucinations cannot compound. Each chunk is an independent statistical event.
