# The POSIX Mapping

Each dimension of the protocol builds on primitives that trace to the original Unix design [@RitchieThompson1974] and were codified in [@IEEEPOSIX].
The mapping adopts a disciplined one-to-one correspondence: each agent concept maps to a primary POSIX primitive, and the semantics are inherited directly from the operating system.

## Identity

Identity concerns how an agent instance is distinguished and bounded at runtime. In Quine, both properties are inherited directly from the process model.

- **Agent instance -> Process (PID):** Kernel-managed identity.
- **Agent boundary -> Address space:** Memory isolation between agents.

The kernel provides a unique identifier—the PID—which is globally unique within the system for the process lifetime; in the local runtime, this removes the need for a separate framework-level identifier.
This identifier is used by the scheduler, the memory manager, and the signal subsystem; Quine inherits these associations rather than reimplementing them.

The address space boundary defines what memory an agent can access.
Two agents (two processes) cannot read or write each other's memory without explicit arrangement (shared memory segments, files).
This isolation is enforced by the MMU at hardware level, not by convention or access control lists in application code.

## Interface

Interface concerns how an agent receives instructions, accepts input, produces output, and signals completion. These interactions map directly to standard process I/O channels.

- **Mission (`argv`):** Immutable mission description.
- **Material (`stdin`):** Data input / material.
- **Deliverable (`stdout`):** Data output / deliverable.
- **Diagnostics (`stderr`):** Diagnostics channel.
- **Outcome (`exit status`):** Success (0) or failure (>0).

The `argv` is set at process creation and cannot be modified during execution.
This immutability makes it suitable for carrying the agent's mission—the high-level instruction that defines what the agent should accomplish.
Because `argv` is separate from `stdin`, the instruction channel is structurally distinct from the data channel.

Standard input (`stdin`) carries material—the content the agent operates on—produced by an upstream process or file.
Standard output (`stdout`) carries the deliverable—the agent's contribution—consumed by downstream processes or redirection targets.
Standard error (`stderr`) carries diagnostics—progress indicators, warnings, and error messages—optionally consumed by observers without polluting the deliverable stream.

The exit status is a small integer (0–255) that serves as a control signal indicating the outcome of execution.
By convention, 0 indicates success; non-zero values indicate various failure modes.
This signal is available to the parent process and enables conditional composition in shell scripts (`&&`, `||`).
Figure \ref{fig:interface} summarizes this five-channel I/O contract.

\begin{figure*}[t]
\centering
\includegraphics[width=0.7\textwidth]{figures/Interface.pdf}
\caption{The Five-Channel I/O Contract: Mission (argv), Material (stdin), Deliverable (stdout), Diagnostics (stderr), and Outcome (exit status). The parent process provides mission via argv and may listen for diagnostics and outcome. Upstream and downstream processes communicate via standard streams.}
\label{fig:interface}
\end{figure*}

## State

State concerns what an agent retains during execution, what survives across continuation, and what can be externalized for coordination. Under POSIX, these forms of state map to three tiers with distinct lifetime scopes.

- **Ephemeral -> Process memory:** Cleared on `exec`/exit.
- **Scoped -> Environment variables:** Preserved across `fork` and `exec`; inherited by children as independent copies.
- **Global -> Filesystem:** Persistent; shared (default) or isolated via namespaces.

**Ephemeral state** (process memory) exists only while the process runs.
When the process calls `exec` or terminates, this state is lost.
For LLM agents, this corresponds to the "working memory" accumulated during a single execution—the context window contents, intermediate computations, and any in-memory data structures.

**Scoped state** (environment variables) survives both `fork` and `exec`.
On `fork`, the child inherits a copy; on `exec`, the calling process retains its environment.
This makes environment variables suitable for passing compact metadata between generations: progress markers, configuration flags, or compressed summaries of prior execution.
The copy-on-fork semantics mean children can modify their environment without affecting the parent.

**Global state** (filesystem) persists beyond process lifetime and is visible to all processes (subject to permissions).
This is the only state tier that survives both `exec` and process termination.
For agents, the filesystem serves as long-term memory, shared artifacts, and coordination medium.

## Lifecycle

Lifecycle concerns how agents are created, delegate work, renew themselves, and terminate. These transitions are expressed directly through process control primitives.

- **Spawn (`fork`):** Create child process with new mission.
- **Continue (`exec`):** Replace current process image while preserving process-level continuity.
- **Terminate (`exit(status)`):** Signal outcome to parent.

An agent can **spawn** children to delegate work (synchronously via `wait`, or asynchronously via background execution and job control).
Each child receives its own `argv` (a distinct sub-mission) and inherits environment variables and context history from the parent.
The parent can block until the child completes, or continue executing while monitoring child status through signals and job control.

Beyond structural delegation, `fork` provides cognitive decomposition: each child operates with an independent context window, allowing complex problems to be partitioned into subproblems that individually fit within context limits.
The parent aggregates results without carrying the full reasoning burden of each subtask.

An agent can **continue** itself by calling `exec` with its own image.
This replaces the process image—clearing process memory—while preserving PID, parent relationship, environment variables, and (optionally) open file descriptors.
In the general POSIX case, `exec` may also install a different executable image and a different `argv`.
Quine's default self re-entry path reuses its current image and mission `argv`, so the agent can reset cognitive context while maintaining the same directive.
For LLM agents, this provides a mechanism to escape context limits: the agent can checkpoint progress in environment variables (or offload to the filesystem for larger state), then `exec` to start fresh with a clean context window while continuing the same task.

An agent **terminates** by calling `exit` with a status code.
This releases all resources (memory, file descriptors, child processes if the parent does not wait), notifies the parent, and makes the exit status available for inspection.
The three operations—spawn, continue, terminate—combined with the state hierarchy, define the lifecycle model illustrated in Figure \ref{fig:lifecycle}.

\begin{figure*}[t]
\centering
\includegraphics[width=\textwidth]{figures/Lifecycle.pdf}
\caption{Agent Lifecycle: The three lifecycle operations (spawn, continue, terminate) and their relationship to state tiers. Spawn creates child processes with distinct missions; continue replaces the process image while preserving process-level continuity, with mission continuity arising in Quine's default self re-entry path; terminate signals completion to the parent.}
\label{fig:lifecycle}
\end{figure*}

The next section describes how this mapping is realized in a concrete runtime.
