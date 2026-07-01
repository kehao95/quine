# Properties

The previous sections described the mapping and the runtime; this section examines what the design yields.
The properties discussed here—Containment, Composition, and Continuity—are not implemented features; they emerge from representing agents as POSIX processes.
By grounding agents in the process model, Quine inherits mechanisms refined over five decades of Unix evolution.

## Containment

Containment in Quine is inherited from the layered enforcement structure of the POSIX process model, not added by the runtime. Each layer constrains a different dimension of agent behavior, and higher layers cannot override lower-layer guarantees.

**Hardware and kernel enforcement.**
Each agent occupies a distinct address space enforced by the MMU; a child crash cannot corrupt the parent's memory.
The kernel further constrains resource consumption (cgroups, ulimit, OOM killer) and privilege (capabilities, namespaces, seccomp).
When limits are exceeded or unauthorized operations attempted, enforcement occurs at the kernel level—the process is terminated with an observable cause.
These bounds apply to agents that are buggy or adversarial: fork bombs, unauthorized network access, and container escapes are bounded by kernel enforcement, not application-layer exception handlers.

**Runtime: instruction-data separation.**
Above these layers, Quine preserves a boundary at the interface level: mission (argv) and material (stdin) travel through distinct OS channels.
The Host maintains this separation when constructing the LLM context: argv maps to the System Prompt, stdin content enters only the User Message.
Combined with instruction hierarchy training [@Wallace2024], this yields a structural basis that may reduce some prompt manipulation risks—control over input data does not by itself grant control over the instruction channel.
This is a claim about architectural separation, not a complete security guarantee.

The result is a containment model that is inherited rather than reimplemented: hardware isolates memory, the kernel bounds resources and privileges, and the runtime preserves a structural distinction between instruction and data.

## Composition

Composition arises from two sources: recursive delegation through `fork`, and external invocability as a standard POSIX command.
The former is a runtime mechanism; the latter is a structural consequence of the process model.

**Recursive composition via `fork`.**
A child process is structurally isomorphic to its parent: same executable, same interface (argv/stdin/stdout/stderr/exit status).
This removes the distinction between "manager" and "worker"—any instance can delegate by spawning children, and any child can coordinate its own subtree.
From the parent's perspective, a delegated subtree remains encapsulated behind the process boundary.

This recursive structure also distributes cognitive load.
Each child begins with an independent context window and reasons only about its subproblem; the parent carries only the child's returned result, not its intermediate reasoning.
Problems exceeding a single agent's reasoning budget can be partitioned into tractable subproblems, with the process tree serving as an implicit divide-and-conquer structure.

**Shell-level composability.**
Because Quine conforms to the standard command contract, it is directly invocable by the shell as an ordinary POSIX executable.
The shell can launch, sequence, redirect, and supervise Quine exactly as it would any other Unix command—composing with pipelines, `xargs`, GNU `parallel`, `cron`, or existing CI workflows without adaptation.
No separate workflow language, daemon, or application-layer protocol is required.

This composability extends reflexively: a Quine instance may invoke `./quine "sub-mission"` through the `sh` tool, treating another agent as an ordinary command.
This is not a second delegation primitive but a consequence of POSIX conformance—the runtime's internal mechanism remains `fork`; shell-mediated invocation is simply the external view of the same single-image design.

## Continuity

Agents face two forms of mortality: context exhaustion (cognitive death) and process termination (physical death).
Quine uses standard POSIX mechanisms to persist across both.

**Surviving cognitive death: continuation across `exec`.**
The `exec` syscall replaces the process image—clearing memory and conversational context—while preserving process ID, parent relationship, environment variables, and open file descriptors.
In Quine's default self re-entry path, the current image and mission `argv` are reused as well.
This gives the agent a way to renew itself without becoming a different computational entity in the process graph.

As an agent approaches context limits, it can `exec` into a fresh instance while preserving several distinct continuity surfaces at once: process identity in the process graph, live material and pipeline state in open file descriptors, and scoped or durable state in environment variables or the filesystem. In Quine's default self re-entry path, mission continuity is preserved as well through reuse of the current `argv`.
The key idea is not to preserve raw cognition, but to preserve enough structured state to make cognition reconstructible.
`wisdom` is one Quine-specific convenience for encoding compact scoped state in environment variables, not the definition of `exec` continuity itself.
Long-lived history can remain on the filesystem; active context is selectively reassembled after renewal.

**Surviving physical death: feedback through stderr and exit status.**
When a child process terminates—whether by completing its task or failing—the parent must decide how to continue.
A failing child emits diagnostics on stderr before terminating; the parent reads this and decides whether to retry, skip, or escalate.
Exit status (0–255) provides a compact outcome signal; shell conditionals (`&&`, `||`) become cognitive branch points.
Standard Unix supervision semantics thus function as adaptive agent coordination.

## Operational Validation

The following observations demonstrate that the architecture is operational—agents can exercise these properties in practice.
These are qualitative feasibility demonstrations, not performance benchmarks.
Detailed execution traces are provided in Appendix A.

**Composition: recursive delegation.**
In an exploratory search task exceeding single-agent budgets, agents used `fork` to spawn parallel workers, assigned disjoint sectors, and coordinated results through the filesystem.
One run produced a 3-level process tree (36 sessions) mirroring the target directory structure—demonstrating recursive delegation and inter-process coordination.
(Appendix A.1)

**Continuity: exec-based self-renewal.**
In MRCR-style needle retrieval tasks [@OpenAIMRCR], agents processed material ranging from 4K to 279K tokens via stdin.
Short contexts required no renewal; long contexts triggered adaptive `exec` cycles—one run required 9 cycles, externalizing progress to environment variables between each renewal.
A baseline comparison (loading full context without streaming) failed on 5 of 8 samples; the streaming architecture with `exec` renewal succeeded on all.
(Appendix A.2)

**Reproducibility.**
Implementation, prompts, and execution logs are available at https://github.com/kehao95/quine.

The same mapping that yields these properties also clarifies where Quine sits among existing systems and where the POSIX model stops being sufficient.
