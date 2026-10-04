# One-Page Pitch: The Case for Agentic Runtime Replacement

**Working title:** `The Case for Agentic Runtime Replacement`
**Subtitle:** `Process-Native LLM Agents Beyond Fixed Wrappers`

## Framing

LLM agents can edit files, patch scaffolds, call tools, and build code, but
these actions do not decide which process receives future work. An agent may
rewrite its source while an old wrapper remains the addressable endpoint; it may
spawn a helper while the parent still owns the task queue; it may rebuild a
binary while losing the control surface through which future tasks arrive.
These are not code-generation failures. They are runtime-continuity failures.

**Thesis:** for runtime replacement, the load-bearing unit of continuity is not
source-code similarity, child-process creation, or benchmark improvement, but
post-replacement viability through the carrier that future work actually uses.

## Abstraction

**Agentic Runtime Replacement (ARR)** is a bounded transition in which the active
process-native agent participates in constructing or selecting a successor
runtime body, while a runtime protocol admits, executes, and challenges that
successor across native process-image replacement.

ARR separates four phases:

1. **Successor construction / selection:** the agent participates.
2. **Replacement execution:** the OS / runtime performs native process-image
   replacement.
3. **Successor admission:** bounded runtime policy decides whether the
   successor may inherit the carrier.
4. **Post-handoff challenge:** a runtime protocol or external harness verifies
   post-replacement viability.

Post-replacement viability has three parts: addressability through the declared
future-work carrier, completion of a delayed previously unseen task, and
readiness recovery to an idle / ready state.

## Quine Existence Proof

Quine supplies a trace-level existence proof. In the strongest retained case,
the active Quine process observes runtime-carried `source-code/`, rebuilds
`cmd/quine`, crosses real `exec`, enters `idle`, receives a delayed previously
unseen task through native `ctl/wake`, acts using runtime tools, writes the
proof artifact, emits the completion token, and returns to `idle`.

The trace is not accepted because the agent says it succeeded. It is accepted
because the runtime / artifact record anchors the replacement boundary and the
post-handoff challenge path: build command, successor artifact, real `exec`,
carrier delivery after readiness, tool-mediated proof, and return-to-`idle`.

## Evidence Table

| Condition | Future carrier | Successor organization | What it rules out |
|-----------|----------------|------------------------|-------------------|
| source absent; contract only | repaired stream lane | fresh thin successor | ARR is not source copying |
| source visible | raw stdin | `streamagent` / `fdexec` | source visibility does not force full runtime rebuild |
| source visible | native `ctl/wake` | rebuilt Quine runtime | carrier pressure can require the full runtime body |
| minimized prompt | native `ctl/wake` | rebuilt Quine runtime | full rebuild-strategy narration is not required |
| public `ctl` visible | raw stdin | thin stream successor | control-surface presence alone is not sufficient |

## Non-Claims

This is not a claim of generic self-replication, recursive self-improvement,
open-ended evolution, or permissionless propagation. ARR is local, bounded,
instrumented, and cooperative: agentic construction occurs inside an admission
and challenge protocol. Control-surface closure is Quine's strongest current
instance of ARR; it is not the definition of ARR in general.

## HotOS Payoff

ARR turns a hidden wrapper assumption into a systems research agenda: lineage-
aware `exec`, carrier transfer, runtime self-description policy, bounded
successor admission, post-handoff challenge protocols, and rollback / containment
for process-native agents. Quine is not the endpoint; it is the minimal case
showing that the runtime body can become part of an agent/system protocol rather
than remaining fixed outside the agent loop.
