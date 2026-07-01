# The Runtime

The implementation is a unitary executable that serves as both runtime and agent template.

## Host-Guest Architecture: Separation of Control and Compute

Quine enforces a physical separation between deterministic control flow and probabilistic computation, as shown in Figure \ref{fig:architecture}.
This separation is not merely architectural preference; it reflects a fundamental asymmetry in where state lives and how failures manifest.

\begin{figure*}[t]
\centering
\includegraphics[width=0.5\textwidth]{figures/Architecture.pdf}
\caption{Host-Guest Architecture: The runtime separates deterministic control (Host) from probabilistic computation (Guest). The Host manages syscalls, file descriptors, and signals; the Guest provides decisions, tool calls, and reasoning.}
\label{fig:architecture}
\end{figure*}

**Host (Local OS Process).**
The Host is a conventional compiled program that maintains the agent's context, lifecycle, and filesystem state.
It is responsible for:

- Parsing `argv` to extract the mission
- Providing an annotated shell environment with file descriptors (including stdin) for the Guest
- Serializing context (conversation history, tool results) for transmission to the Guest
- Managing child processes spawned by the Guest's tool invocations
- Calling `exit` with the status code determined by the Guest

The Host maintains conversation history, control state, and I/O buffers for tool execution.
Because the Host is an ordinary OS process, agent instances inherit standard resource management: they can be scheduled, signaled, killed, or resource-limited by POSIX tools (`nice`, `kill`, `ulimit`, `cgroups`) without requiring framework-specific supervision.

**Guest (Remote LLM API).**
The Guest acts as a stateless cognitive oracle.
It receives the serialized POSIX state (mapped to prompts) and returns structured tool invocations.
The Guest has no persistent state between calls; all context must be provided in each request.
This statelessness is a design choice: it ensures that the Host maintains authoritative state and that the Guest can be replaced, rate-limited, or load-balanced without coordination.

The Host contains no task-specific logic; decision-making is delegated to the Guest.
The Host is purely reactive: it provides the environment the Guest operates in, captures the results, and iterates.

**Context Assembly.**
The Host preserves OS-level I/O boundaries when assembling context:

1. The mission (`argv`) is placed in the System Prompt.
2. Material (`stdin`) is announced in the User Message; access is provided via file descriptor remapping (Section 3.2).
3. Tool outputs are returned as Tool Result messages.

This layered structure aligns with instruction hierarchy training [@Wallace2024], where models prioritize System > User > Tool.
The Host does not rewrite or summarize streams; it annotates boundaries and delegates interpretation to the Guest.

**Tool Interface.**
The Host exposes four tools—`sh`, `fork`, `exec`, `exit`—named after their POSIX counterparts with semantics as defined in Section 2.4:

- **`sh`:** Execute shell command, return stdout/stderr/exit status.
- **`fork`:** Spawn child Quine process(es) with new `argv`; optionally block until completion.
- **`exec`:** Replace current process image; preserve pid, env vars, and file descriptors. By default Quine re-execs its own binary with the current mission `argv`, but explicit target/`argv` can hand control to a different executable.
- **`exit`:** Terminate with status code.

The `sh` tool is the primary interaction mechanism—file operations, compilation, and system commands all go through shell invocations.
The `fork` tool enables delegation without leaving the process model; each child receives its own mission and returns structured results.
The `exec` tool is how agents manage context growth: in the default self re-entry path, agents checkpoint progress in env vars or the filesystem, then `exec` into a fresh Quine instance with the same mission and a clean context window.
The same primitive can also hand off directly to a non-Quine executable when the task is better finished as an ordinary POSIX process.

## POSIX Conformance: File Descriptor Mapping

From the shell's perspective, Quine is a standard POSIX filter: stdin in, stdout out, stderr for diagnostics, exit status for outcome.
This means Quine composes anywhere a traditional Unix filter can: pipelines, shell scripts, subprocesses.

The challenge is maintaining clean context windows while honoring this contract.
The runtime solves this by exposing annotated file descriptors within the shell environment.
The runtime's stdin, stdout, and stderr are passed into each shell invocation as higher-numbered file descriptors:

- **fd 3:** Runtime's stdin (the material stream)
- **fd 4:** Runtime's stdout (the deliverable channel)
- **fd 5:** Runtime's stderr (the failure-signal channel)

This mapping leaves the shell's standard file descriptors (0, 1, 2) available for normal command I/O.
To read input material, the Guest reads from fd 3.
To emit a deliverable, the Guest writes to fd 4 (e.g., `echo "result" >&4`).
To emit streaming failure signals, the Guest writes to fd 5 (e.g., `echo "failed" >&5`).
The shell command's own stdout/stderr are captured and returned to the Guest as tool results for reasoning.

This separation allows the Guest to capture command output for reasoning while emitting deliverables to downstream processes through a separate channel.
Without it, the agent would face a dilemma: pollute the deliverable stream with intermediate outputs, or lose visibility into command results.

## Illustrative Shell Compositions

Quine's adherence to standard streams enables composition with classic Unix utilities and multi-agent pipelines.

**Example 1: Single-process replacement.**
A cognitive filter that understands intent:

```bash
cat server.log | ./quine "Extract lines indicating auth failures" | sort | uniq -c
```

Here Quine replaces `grep` in a traditional pipeline.
The cognitive filter receives log lines on stdin, applies semantic understanding to identify authentication failures, and emits matching lines to stdout.
Downstream tools (`sort`, `uniq -c`) process the output normally.

**Example 2: Implicit DAG.**
Three agents form a reasoning chain:

```bash
git diff HEAD~1 | \
  ./quine "Identify the changed components" | \
  ./quine "Assess risk level for each change" | \
  ./quine "Generate a review checklist"
```

Each agent receives the previous agent's output as stdin and contributes its analysis to stdout.
The shell provides the topology; no orchestrator or shared memory is required.
Each agent runs in a separate process with separate context.

**Example 3: Control flow without framework.**
Exit codes drive branching:

```bash
./quine "Apply the patch" < fix.patch && \
  ./quine "Run the test suite and report results"
```

The first agent attempts to apply a patch; it exits 0 on success, non-zero on failure (e.g., conflicts, malformed patch).
The shell's `&&` operator conditionally executes the second agent only if the patch applied cleanly.
This is cognitive branching using standard shell control flow.

Because Quine is packaged as a standard POSIX executable, the same composability extends reflexively: a Quine instance may invoke `./quine "mission"` through the `sh` tool. This should be understood not as a second lifecycle primitive, but as an external composition path enabled by the single-image design. The runtime's internal delegation mechanism remains `fork`; shell-mediated re-invocation is a consequence of POSIX conformance, not a separate runtime protocol.

## Implementation Scale

The Host is implemented in approximately 3,100 lines of Go (cognitive loop, tool execution, job management, and session recording). LLM provider abstractions add another 2,600 lines (protocol adapters, OAuth flows, and configuration). The complete implementation totals approximately 5,700 lines, compiled to a single ~9.8 MB binary that supports multiple LLM providers (Anthropic Claude, OpenAI GPT, Google Gemini) through a provider abstraction layer. Provider selection is controlled by environment variables, enabling the same binary to use different backends without recompilation.

The runtime contains no domain-specific task logic; decision-making is delegated to the Guest (LLM).
The properties that emerge from this design—containment, composition, and continuity—are examined in the next section.
