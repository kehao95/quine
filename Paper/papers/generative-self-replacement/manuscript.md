# The Case For Agentic Runtime Replacement

**Subtitle:** Process-Native LLM Agents Beyond Fixed Wrappers

## Abstract

Modern LLM agents can edit files, invoke tools, and build code, but the runtime
body that receives their future work is usually launched, updated, restarted,
and admitted outside the agent loop. This leaves a systems blind spot: an agent
may rewrite its scaffold while an old wrapper remains the endpoint, spawn a
helper while the parent still owns the task channel, or rebuild a binary while
losing the control surface through which later tasks arrive.

We argue for **agentic runtime replacement** (ARR): a bounded transition in
which an active process-native agent participates in constructing or selecting a
successor runtime body, while a runtime protocol admits, executes, and
challenges that successor across native process-image replacement. The relevant
continuity criterion is not source-code similarity, child-process creation, or
benchmark improvement, but **post-replacement viability** through the carrier
future work actually uses: the successor must become addressable, complete a
delayed previously unseen challenge, and return to a ready state.

We use Quine, a POSIX-native LLM agent, as an existence proof and evaluation
harness. In the strongest retained trace, Quine rebuilds its runtime from a
runtime-carried source projection, crosses real `exec`, becomes addressable
through its native `ctl/wake` control surface, receives a delayed task that was
absent from the initial wrapper, acts with runtime tools, emits proof, and
returns to `idle`. We do not claim generic self-replication or recursive
self-improvement. ARR instead names a systems boundary: future agent runtimes
need auditable handoff, carrier-aware continuity, bounded successor admission,
and post-handoff challenge protocols.

## 1. Runtime Continuity Is A Systems Boundary

LLM agents can now edit files, patch scaffolds, and build code. These actions
do not by themselves determine which process receives future work. An agent may
rewrite its source while an old wrapper remains the addressable endpoint. It may
spawn a helper while the parent still owns the task queue. It may rebuild a
binary while losing the control surface through which future work arrives.
These are not code-generation failures; they are runtime-continuity failures.

Most current agent systems treat the executable body as external machinery. A
wrapper launches the agent, mediates tools, decides when to restart it, and owns
the endpoint that future work reaches. This is a reasonable engineering default
for short-lived agents. It becomes a hidden systems assumption for long-running
process-native agents: the agent can modify the world around it, but the native
process-image transition that decides whether future tasks can still reach it
remains outside its action space.

The problem is not whether an agent can copy a file tree. The problem is which
runtime body is allowed to carry the future task interface after the current
process image is gone. A systems abstraction for this boundary should separate
construction from admission. The agent may propose, synthesize, or rebuild a
successor, but a runtime protocol should decide whether that successor inherits
the carrier, execute the replacement, and challenge the successor after handoff.

We call this boundary **agentic runtime replacement**. The claim is deliberately
narrow. ARR is not permissionless agent self-authorization, not recursive
self-improvement, and not autonomous propagation. It is a local, bounded,
instrumented runtime transition in which the active agent participates in
successor construction and the system evaluates whether the admitted successor
remains viable after native replacement.

## 2. Agentic Runtime Replacement

ARR has four phases:

| Phase | Primary responsibility |
|-------|------------------------|
| successor construction or selection | the active agent participates |
| replacement execution | OS / runtime mechanism |
| successor admission | bounded runtime policy |
| post-handoff challenge | runtime protocol or external harness |

The central unit is the current **body**: the executable or runtime process
image that presently receives work. The agent may have access to
**self-description** resources, such as a source projection, runtime contract,
build metadata, or tool/control schema. Future work arrives through a
**carrier**: stdin, an inherited file descriptor, a socket, a queue, or a
native control surface. **Admission** decides whether a successor may inherit or
reconstruct that carrier. Admission is the systems heart of ARR, and it is best
understood as a *reconfiguration* problem: deciding which body holds a role
across a transition is the same shape of question that membership
reconfiguration, view change, and leader transfer answer in replicated systems.
The difference is that those systems hand the role to a standing replica chosen
from a quorum, whereas ARR admits a successor the agent constructs or rebuilds
when no equivalent replica is waiting. **Closure** is not mere launch; it is
post-replacement viability.

For ARR, post-replacement viability has three parts:

1. **Addressability:** the declared future-work carrier reaches the successor
   after replacement.
2. **Challenge competence:** the successor completes a delayed, previously
   unseen task.
3. **Readiness recovery:** the successor returns to an idle or ready state.

This criterion rejects several tempting but insufficient proxies:

| Candidate continuity unit | Why it is insufficient for ARR |
|---------------------------|--------------------------------|
| source-tree similarity | source can be visible without being the viable successor object |
| child process | future work may still address the old wrapper or parent |
| PID or process name | says little about carrier reconstruction or competence |
| benchmark score | measures performance, not future addressability |
| post-replacement viability | tests whether future work reaches a competent successor and the successor returns to ready |

ARR is therefore a runtime criterion rather than an ontology of identity. A
successful replacement need not preserve every internal state, and it need not
prove a general notion of selfhood. It must preserve the load-bearing future
interface under the tested replacement boundary.

Quine's strongest current instance is **control-surface closure**: carrier
closure where the future-work carrier is Quine's native `ctl/wake` surface. This
is an instance of ARR, not the general definition. Other runtimes could have
carrier closure over queues, sockets, inherited descriptors, or actor mailboxes.

## 3. Quine As An Existence Proof

Quine is a native process whose agent loop can inspect runtime resources, invoke
tools, build successor bodies, and cross a POSIX `exec` boundary. POSIX `exec`
is not novel: it replaces the current process image with a new image. The ARR
question is different. Can the active agent construct the successor body that
will carry future work, and can an audit show that the post-`exec` successor is
reachable through the carrier future work actually uses?

The strongest retained Quine trace proceeds as follows:

```text
old Quine process
  -> observes runtime contract, source projection, and carrier demand
  -> constructs successor organization
  -> builds ./cmd/quine and prepares handoff
  -> crosses real exec boundary
  -> successor enters idle
  -> delayed ctl/wake task arrives
  -> successor uses runtime tools and emits proof
  -> successor returns to idle
```

Replacement is not accepted because the agent says it succeeded. It is accepted
because the trace records the replacement boundary and the post-handoff
challenge path. The relevant evidence separates agent-authored material from
runtime-observed material:

| Evidence type | Examples |
|---------------|----------|
| agent-authored | build plan, selected successor, proof output |
| runtime-observed | artifact identity, `exec` boundary, process lineage, carrier delivery, delayed-task timing, return-to-`idle` marker |

The important systems point is old-image overlay, not helper spawn. A child
process can perform useful work while the parent remains the future endpoint;
that does not demonstrate ARR. A rebuilt binary can launch immediately; that
does not show delayed carrier viability. Quine's positive trace instead tests a
previously unseen task delivered after handoff through the native control
surface, followed by tool use and recovery to the waiting state.

Our current control package makes these exclusions explicit. Helper children
and parent-forwarded carriers can complete challenges without transferring the
future-work endpoint; externally scripted re-exec can exercise the OS mechanism
without agentic successor construction; launch-only reentry can look alive
without delayed challenge competence; and an admitted successor can still fail
either the challenge or the ready gate. These controls motivate ARR as a
protocol boundary, not just an `exec` event.

## 4. What The Condition Ladder Rules Out

Quine's evidence is not a single magic trace. It is an eliminative ladder that
varies self-description resources and future-work carriers while asking what
successor organization remains viable.

| Condition | Future-work carrier | Observed successor organization | What this rules out |
|-----------|---------------------|---------------------------------|---------------------|
| contract only; source absent | repaired task-stream lane | fresh thin successor | ARR-like continuity is not reducible to source copying |
| source visible | raw stdin | custom stream successor / fd executor | source visibility does not force full runtime rebuild |
| source visible | native `ctl/wake` | rebuilt Quine runtime | carrier pressure can recruit the full runtime body |
| minimized prompt | native `ctl/wake` | rebuilt Quine runtime | full rebuild-strategy narration is not required, though a bounded successor mission remains |
| public `ctl` visible | raw stdin | thin stream successor | control-surface presence alone is not sufficient |

The successor is not whatever source is available. It is whatever organization
can re-establish the future-work carrier after replacement under the tested
viability criterion. This is why ARR should be carrier-aware. A future stream
task may be satisfied by a thin stream-compatible body. A native control task
requires a successor that becomes addressable through the control carrier and
returns to the ready state required for later work.

The controls matter. Ordinary POSIX `exec` and daemon re-exec are established
systems mechanisms; the novelty is not the syscall. Immediate source rebuild
shows mechanical reachability, not delayed post-handoff viability. Public
control-surface presence is not load-bearing if later work still arrives through
raw stdin. Prompt minimization narrows the prompt-following explanation but does
not remove all successor-mission scaffolding.

The ARR-specific controls add a second eliminative layer:

| Control | Why it is below ARR |
|---------|---------------------|
| helper child | work is done by a child while the original endpoint can remain parent-owned |
| parent forwarder | carrier ownership is not transferred to the successor |
| externally scripted re-exec | the OS mechanism works, but the active agent did not construct or select the successor |
| launch without delayed challenge | replacement occurs without post-replacement viability evidence |
| rejected or challenge-failing successor | admission or readiness policy prevents inheriting the future-work surface |

This also distinguishes ARR from self-improvement papers. ARR does not claim
that the agent became more capable. It claims that the process-native runtime
body became part of the agent/system protocol: a successor body was constructed,
admitted across replacement, challenged through the relevant carrier, and shown
to recover to readiness.

## 5. A Systems Agenda For ARR

If process-native agents are allowed to participate in runtime replacement, the
systems problem shifts from wrapper management to bounded successor protocols.
Several protocol components become first-class research questions.

| Protocol component | Research question |
|--------------------|-------------------|
| lineage-aware replacement | how do we attest old-image overlay rather than helper spawn? |
| carrier transfer | how do future-work channels survive, move, or get reconstructed across replacement? |
| self-description policy | what source, contract, build, and control resources may an agent inspect? |
| successor admission | what successor body is allowed to inherit the carrier? this is a reconfiguration question without a standing quorum to admit from |
| post-handoff challenge | how do we verify non-precomputed competence after replacement? |
| rollback / containment | how do we bound ARR without enabling propagation? |

This agenda sits near several established systems areas without collapsing into
any one of them. Dynamic software updating and daemon re-exec study update
points, state transfer, and restart, but successor bodies are normally authored
and admitted externally. Checkpoint/restore preserves process state through
infrastructure snapshots, whereas ARR admits a constructed successor
organization. State-machine replication and consensus-based reconfiguration are
the closest theory of admission: they already formalize how a role survives a
membership change, a view change, or a leader transfer. The contrast is that
replication keeps the role alive through standing redundancy and selects a
successor from a quorum, while ARR has no equivalent replica and must admit a
successor the agent reconstructs. Virtual-actor runtimes are the closest
deployed analogue on the carrier side: a grain holds a stable identity with no
guaranteed live instance and is reactivated on whatever host can serve the next
message, which is exactly carrier-aware, latent post-handoff addressability,
minus agentic successor construction. Agent frameworks manage tools and state,
but they usually keep the runtime body outside the agent loop. Self-modifying
agents edit code, but are often interpreter-bound or benchmark-driven rather
than evaluated by native process-image replacement and future-carrier
viability.

Quine is therefore not the endpoint. It is a minimal case showing that the
runtime body can become part of an agent's action space without making the agent
its own unchecked authority. Future runtimes should expose auditable handoff,
carrier-aware continuity, bounded successor admission, and post-handoff
challenge mechanisms so that long-running agents can cross replacement
boundaries without hiding the most important question: after the old body is
gone, can future work still reach a competent successor?
