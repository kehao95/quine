# Signal-First Interaction Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md) § Navigation And Routing.

When a change touches interactive operation, CLI shape, or long-lived session behavior, treat the canonical problem as signal handling plus agent-chosen response rather than conversation simulation.

Ontology owner: [`operating-philosophy.md`](../../Paper/philosophy/operating-philosophy.md) § Control: Textual Payloads, Physical Signals (signals vs payload surfaces, signal-class separation, adapter lowering). Client/command owner: [`operator/qcli/`](../../operator/qcli/README.md) — `CONTRACT.md` for the wire surface and command set, `DESIGN.md` for the control-action semantics (`post` / `poke` / `inject` / `interrupt`) and client object model. This contract owns only the development-time review rule.

Named defect:

- the historical required single-mission `argv` contract made interactive work look like a question of "how to inject more user messages" or "how to hand the process another task" instead of first asking what external signals the process can receive and how it should respond
- conventional coding-agent UX implies a conversational request/response loop in which every submitted message expects a reply; Quine's client surface must not inherit that assumption, because a signal may be delivered with no expected response and any reply is agent-chosen

Review rules:

- when a design starts from "how do we append another user message?" or "how do we inject another task?", stop and restate the problem as "what signal arrived, when is it delivered, and what responses are available to the process?"
- do not treat Human as the default source of input; a signal may originate from a Human, another Quine, a scheduler, a file-mediated control plane, or some later runtime source, and the isomorphism between Human-originated interaction and peer-Quine coordination must survive at the Quine boundary
- if a provider protocol requires role-shaped messages such as `user`, treat that as a lowering detail at the model boundary rather than as Quine's primary self-model
- do not smuggle queue semantics into the runtime contract unless ordering, backlog, or fairness is actually required by the owning surface
- do not overclaim a control action before the runtime can actually perform it; bridge behavior that only interrupts active tool work is a precursor, not the final contract
- keep Human-facing client / TUI / UI concerns in the companion client process rather than adding interaction-specific command modes to the `quine` binary; keep any unavoidable lowering to role-shaped messages or task-shaped bookkeeping narrow, reversible, and outside the control-plane description of Quine itself

Reflection checkpoint: ask first what signals exist, how they are transported and delivered, their priority/preemption semantics, and who originates them; only then whether a derived internal work unit (task, commitment, plan) helps. If any answer is borrowed lazily from today's adapter shape, the design is overfit to the current UI.

Deferred lifecycle unification note (no owner doc yet):

- `init`, `continue`, `spawn`, `fork`, and `exec` should be reasoned about as combinations of session construction and session activation, not as unrelated lifecycle cases
- session construction decides session identity, incarnation allocation, mission projection, context seed, parent relation, and workspace/world seed
- session activation decides process image, environment lowering, stdio/capture policy, working directory, and whether the current process, a child process, or an exec replacement runs that constructed session
- do not collapse the public affordances prematurely; first keep the shared internal boundary clean enough that `fork` can remain context-preserving, `spawn` can remain fresh, `continue` can remain same-session/current-incarnation activation with a new physical run, and `exec` can remain same-session/image-replacement
