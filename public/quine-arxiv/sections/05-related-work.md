# Related Work and the Boundaries of POSIX

Quine sits among existing agent systems as a runtime-level alternative; from there, the discussion turns to where the POSIX model ends and what must extend beyond it.

## Related Systems

I classify current agent systems by their structural relationship to the operating system.
The taxonomy is organized by where agent identity, lifecycle, interface, and isolation are realized—not by end-task capability or developer-facing features.

### Host-Coupled Extensions

Systems such as Cursor and GitHub Copilot are integrated within host environments such as IDEs and do not provide agents with independent process lifecycles.
The agent component exists within the host application rather than as a separately managed runtime entity: it is not exposed as a first-class POSIX process, does not present an independent stdin/stdout interface, and terminates with the host application.
This tight coupling enables deep integration—such as inline suggestions and direct access to editor state—but correspondingly limits independent spawning, termination, and shell-level composition with external tools.

### Application-Layer Schedulers

Frameworks such as LangGraph, AutoGen [@Wu2023], and CrewAI [@CrewAI] provide user-space scheduling and message buses.
Agents are realized as in-process objects (Python classes, coroutines) with fault isolation implemented through exception handling rather than hardware-enforced address space separation.
These systems have successfully demonstrated multi-agent coordination and have large ecosystems of tools and integrations; their design optimizes for developer productivity and rapid prototyping.
Structurally, agents within these frameworks typically share a single OS process: a crash in one agent (unhandled exception, memory corruption) can propagate to others, and scheduling is managed by the framework dispatcher rather than the kernel.

### Sandboxed Loops

Systems such as Devin, OpenHands [@Wang2024], and SWE-agent [@Yang2024] use an external controller to manage containerized tools.
The agent's "body" (shells, browsers, file access) runs in an isolated sandbox; the "brain" (LLM) runs in the controller process.
This architecture provides strong tool isolation—a runaway shell command cannot escape the container—while centralizing cognitive coordination, and has proven effective for complex software engineering tasks.
Structurally, however, the sandbox isolates tools rather than agents themselves: the cognitive loop (prompt -> LLM -> action -> observation) runs in a single controller process, and multiple "agents" are typically coroutines or threads within this controller rather than separate OS processes.

### Pseudo-OS Middleware

Systems such as AIOS [@Mei2024], agentOS [@Li2026], and UFO 2 [@Zhang2025] implement OS-like abstractions in user space, providing schedulers, memory managers, and IPC mechanisms for agents using terminology borrowed from operating systems.
These systems recognize that agent management resembles process management and attempt to provide similar abstractions.
However, despite OS-inspired naming, agents in these systems are typically objects or threads within a single host process.
"Process isolation" is simulated through software boundaries rather than hardware-enforced address spaces; the "scheduler" is a user-space dispatcher rather than the kernel's CFS; resources are accounted at the framework level rather than by cgroups or ulimit.

### Quine's Position

Quine represents a distinct design point.
Existing agent runtimes typically either manage agents at the application layer or use the OS primarily as a tool sandbox, rather than directly realizing each agent as a native POSIX process with the full reasoning-and-acting loop exposed through standard process interfaces.

This difference is about runtime organization, not end-task superiority.
Application-layer frameworks offer flexibility, rich ecosystems, and rapid development cycles.
Quine trades some of this flexibility for structural properties inherited from the OS: kernel scheduling rather than framework dispatch, standard streams rather than framework-specific message formats, process lifecycle primitives rather than object instantiation, and hardware-enforced isolation rather than exception handling.

Having located Quine among current systems, I now turn from comparative structure to the limits of the substrate itself.
The POSIX mapping provides a valid first abstraction for agents, but process semantics do not exhaust the runtime needs of cognition.
Among the directions that boundary exposes, two are especially immediate: one concerns space—whether an agent's world can be scoped by relevance rather than only permission; the other concerns time—whether cognition and side-effects can be revised on different terms.
Both arise directly from taking the process model seriously: once the agent is granted an execution boundary, the next question is what world that boundary encloses, and what kind of temporality its actions inhabit.

## From Namespace to World

Plan 9 showed that a process need not inhabit a single global filesystem: namespaces can be per-process, and interfaces can be constructed rather than merely inherited [@Pike1993].
For agents, this insight is necessary but insufficient.
An agent requires not merely a different namespace, but a world organized by task relevance—a situated perspective in which some entities are present, others absent, and still others foregrounded or merely nameable.

A subjective world is not a false world; it is a selectively constituted one.
The distinction matters because it separates runtime responsibility from epistemic illusion.
When a debugging agent sees logs, stack traces, and failing tests while a planning agent sees deadlines, owners, and design intents, the difference is not merely what files each may access.
It is what objects are present as first-class entities at all.
Security asks what an agent may access; subjectivity asks what its world is made of.

This reframing has architectural consequences.
Traditional sandboxing restricts resources; a cognitive runtime must do more—it must scope reality.
The agent does not merely execute within constraints; it is situated within a world whose boundaries are drawn by relevance, role, and task.
Two agents on the same machine, with access to the same underlying storage, may nonetheless inhabit different operational worlds if their runtimes foreground different entities and relations.
In such a runtime, the difference may appear not only in which paths are mounted, but in which objects are surfaced as logs, goals, owners, hypotheses, or pending obligations.

In Quine, this constitutive role does not belong to an external control plane.
Because the runtime is recursively self-instantiating, the same executable that inhabits a world can also construct a different world for a child agent.
A parent does not merely launch a subprocess; it defines the visible environment into which that subprocess is born.
The agent is therefore neither a passive resident of a pre-given environment nor the object of a separate orchestrator's worldview—it is a constitutive participant in the production of local worlds, for itself through renewal and for others through delegation.
The distinctive point is not only that worlds can be scoped, but that world-construction is endogenous to the runtime itself.

POSIX can isolate processes; Plan 9 can differentiate namespaces.
But a cognitive runtime may need to expose worlds whose constitution is task-relative rather than permission-derived.
This does not yet define a mechanism; it marks a shift in what the runtime must be responsible for making visible.

## From Execution to Revision

POSIX time is operational and forward-moving: actions happen, effects accumulate, processes terminate.
Rollback, where available, is external, partial, and typically resource-centric—restoring a file, replaying a log, restarting a container.
But cognition is not merely sequential; it is provisional.
An agent may reconsider, backtrack, explore alternatives, and retain lessons from failed attempts.

The core difficulty is that conventional rollback conflates two kinds of state that agents need to treat differently.
Mental state comprises beliefs, plans, branches explored, rejected hypotheses, and learned constraints.
Environmental state comprises files changed, commands executed, resources allocated, and messages sent.
When these cannot be revised independently, the choice is stark: either both are lost, or neither is reversible.
What agents need is rollback without amnesia—the ability to undo effects while preserving experience.
A failed branch may need to retract file edits and spawned processes while preserving the constraints it discovered and the options it ruled out.

This is not simply an engineering problem of better checkpointing.
Agent work often involves speculative branches: multiple candidate futures explored from a common past, some committed, others abandoned.
A runtime that only knows committed execution cannot treat branching as first-class; it reduces counterfactual exploration to ad hoc application logic, with each framework inventing its own replay mechanism.

POSIX manages process lifetime but does not express the provisionality of cognition.
If revisability is central to how agents think, then time itself becomes part of the runtime contract—not a library feature bolted on afterward, but a structural concern that the operating layer must acknowledge.
The question is not how to implement snapshots, but whether alternative futures can be explored without reducing them to workarounds above the OS.

## Note on Scope

The two directions above—world and time—do not exhaust the boundaries of the POSIX model.
Internal cognitive structure (making the agent's reasoning addressable rather than opaque) and distributed composition (preserving file-based abstractions across machine boundaries) mark additional frontiers.
This paper focuses on the POSIX mapping itself; these extensions belong to future work.

The broader question—what an operating system for cognition should expose, delimit, remember, and compose—is closer to the Plan 9 lineage [@Pike1995] than to conventional POSIX extension.
The aim is not merely to add mechanisms, but to rethink what the runtime should make visible, nameable, and composable.
