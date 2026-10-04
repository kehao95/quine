# Manuscript Plan: The Case for Agentic Runtime Replacement

The dossier owns thesis, scope, evidence posture, and paper-state decisions.
Use this file for the HotOS-style position-paper plan and for the boundary with
any later archival systems paper.

The compact outward-facing pitch is [`plans/hotos-one-page-pitch.md`](plans/hotos-one-page-pitch.md).


## Draft Status

The first HotOS-style prose draft now lives in [`manuscript.md`](manuscript.md).
Treat it as a position-paper draft surface, not as an archival systems
evaluation. Next edits should primarily compress for page budget, attach
primary-source citation anchors, and replace any remaining internal shorthand
with public condition names.


## Source-Check Absorption

The 2026-05 primary-source check narrows the HotOS claim rather than weakening
it. ARR should cite nearby systems generously, but the paper's novelty should be
stated as the synthesis of four pieces rather than any one mechanism:

1. an LM-driven process participates in constructing or selecting its successor
   runtime body;
2. a bounded runtime protocol admits and executes the replacement;
3. the handoff crosses native process-image replacement rather than only
   spawning a helper or patching an interpreter object;
4. a delayed post-handoff challenge tests the carrier future work actually
   uses.

Do not frame the contribution as novel POSIX `exec`, generic dynamic software
updating, checkpoint/restore, actor hot-swap, or self-modification. The
primary-source queue for a submission bibliography is: POSIX `exec`, DSU / hot
patching, checkpoint/restore, actor hot-swap, self-modifying agents, and agent
observability. Treat the current source-check report as a routing guide until
those sources are independently verified in the manuscript bibliography.

## Current Preferred Projection

**Working title:** `The Case for Agentic Runtime Replacement`

Subtitle candidates:

- `Process-Native LLM Agents Beyond Fixed Wrappers` (preferred)
- `Carrier-Aware Continuity for Process-Native LLM Agents`
- `Control-Surface Closure Across Native Process Replacement`

`Why Agents Should Survive Their Wrappers` is useful as a talk title, not as the
submission subtitle.

This is a HotOS-style position / abstraction paper, not a full evaluation
paper. Its job is to make a systems community care about a new runtime boundary:
long-running LLM agents can edit files, invoke tools, and alter scaffolds, but
the native runtime body that carries future work is still usually launched,
updated, restarted, and admitted by an external wrapper or supervisor.

Quine is used as an existence proof and evaluation harness for that boundary:
an active process-native agent participates in constructing or selecting the
successor runtime body that will carry its future task interface, while a
bounded runtime protocol admits, executes, and challenges that successor across
real process-image replacement.

Current assessment:

- high theory / position-paper leverage;
- lower empirical burden than an ATC / OSDI / SOSP-style paper;
- still requires a credible audited trace and compact condition table so the
  paper does not read as a manifesto.

## One-Sentence Framing

> LLM agents are becoming long-running systems processes, but their runtime
> bodies are still externally managed. Quine asks whether an agent can control
> the native process-image transition that preserves its own future task
> interface.

## Central Thesis

> For runtime replacement, the load-bearing unit of continuity is not
> source-code similarity, child-process creation, or benchmark improvement, but
> post-replacement viability through the carrier that future work actually
> uses.

This thesis deliberately challenges three default framings:

| Default question | This paper's replacement question |
|------------------|-----------------------------------|
| Did the agent improve on a benchmark? | Can future work still reach the post-replacement agent through the intended carrier? |
| Did the agent copy source or body? | What successor organization carries the future task interface? |
| Is this ordinary restart / re-exec? | Who constructs and admits the successor body, and how is post-handoff viability verified? |

## Core Abstraction

**Agentic Runtime Replacement (ARR):** a bounded transition in which the
active agent participates in constructing or selecting a successor runtime body,
while a runtime protocol admits, executes, and challenges that successor across
native process-image replacement.

ARR separates four phases:

| Phase | Primary responsibility |
|-------|------------------------|
| successor construction / selection | agent participates |
| replacement execution | OS / runtime mechanism |
| successor admission | bounded runtime policy |
| post-handoff challenge | runtime protocol / external harness |

ARR should not imply that an agent authorizes its own successor without a
runtime boundary. The novelty is agentic construction inside an admission and
challenge protocol.

**Post-replacement viability:** the continuity criterion for ARR. Viability has
three parts:

1. **Addressability:** the declared future-work carrier reaches the successor
   after replacement.
2. **Challenge competence:** the successor completes a delayed previously
   unseen task.
3. **Readiness recovery:** the successor returns to an idle / ready state.

**Closure levels:**

- **Level 1 — ARR event:** the agent participates in successor construction or
  selection, and the current process crosses native process-image replacement.
- **Level 2 — carrier closure:** post-replacement viability holds over the
  declared future-work carrier.
- **Level 3 — control-surface closure:** carrier closure where the carrier is
  Quine's native `ctl/wake` surface.

Control-surface closure is not the definition of ARR; it is the strongest Quine
instance currently retained.

Minimal vocabulary:

| Term | HotOS definition |
|------|------------------|
| body | executable / runtime process image that currently receives work |
| self-description | resources the agent can use to reconstruct or reinterpret itself: source projection, runtime contract, build metadata, control / tool schema |
| carrier | future-work channel: stdin, inherited fd, `ctl/wake`, socket, queue, mailbox |
| admission | runtime policy deciding whether a successor may inherit the carrier |
| closure | post-replacement viability: carrier reachability, delayed-task competence, and return-to-ready |

## Contributions

The HotOS paper should claim three contributions.

1. **Problem reframing.** It identifies externally managed runtime bodies as a
   hidden systems boundary for long-running LLM agents. The issue is not merely
   how agents call tools, but who controls the process-image transition that
   carries future tasks.
2. **ARR abstraction and criterion.** It proposes Agentic Runtime Replacement
   as agent construction plus runtime admission, and argues that runtime
   replacement should be evaluated by post-replacement viability through the
   carrier future work actually uses.
3. **Quine existence proof and agenda.** It uses Quine to demonstrate the
   boundary with a real `exec` trace and a compact condition table, then opens a
   systems research agenda around lineage-aware `exec`, carrier transfer,
   self-description policy, successor admission, and post-handoff challenge
   protocols.

## Abstract Draft

Modern LLM agents can write code, invoke tools, and modify files, but the
runtime body that receives future work is usually launched, updated, and
admitted outside the agent loop. This leaves a systems blind spot: an agent may
alter its scaffold without controlling the process-image transition that
determines whether future tasks can still reach it.

We argue for **agentic runtime replacement**: a bounded transition in which an
active process-native agent constructs or selects a successor runtime body,
while a runtime protocol admits that successor across native process-image
replacement and challenges it after handoff. The key continuity criterion is not
source-code similarity, child-process creation, or benchmark improvement, but
post-replacement viability through the carrier future work actually uses.

We explore this boundary through Quine, a POSIX-native LLM agent. In Quine,
source availability alone does not determine what is reproduced: raw-stream
pressure can be satisfied by thin stream successors, while native `ctl/wake`
pressure recruits the rebuilt runtime needed for control-surface closure. In the
strongest retained trace, the rebuilt successor crosses real `exec`, becomes
addressable through `ctl/wake`, receives a delayed previously unseen task,
completes it using runtime tools, and returns to `idle`.

We do not claim generic self-replication or recursive self-improvement. Quine
instead motivates a systems agenda for auditable handoff, carrier-aware
continuity, bounded successor admission, and post-handoff challenge protocols.

## HotOS 5-Page Shape

Page budget target: 5 pages including figures and tables, with references
outside the core budget when the venue permits.

| Page | Budget | Content |
|------|--------|---------|
| 1 | 0.9 page | problem, failure modes, thesis, non-goals |
| 2 | 1.0 page | ARR abstraction, admission boundary, viability criterion, vocabulary |
| 3 | 1.1 page | Quine trace, authority / trace figure, independent audit hooks |
| 4 | 1.0 page | condition table, controls, punchline |
| 5 | 0.9 page | protocol-component research agenda, compact related-work contrasts, ending |

### Page 1: Problem And Provocation

Goal: make the systems gap obvious before introducing Quine.

Opening move:

> LLM agents can now edit files, patch scaffolds, and build code, but these
> actions do not by themselves determine which process receives future work. An
> agent may rewrite its source while an old wrapper remains the addressable
> endpoint; it may spawn a helper while the parent still owns the task queue; it
> may rebuild a binary while losing the control surface through which future
> work arrives. These are not code-generation failures. They are runtime-
> continuity failures.

Required content:

- current agents can edit code, invoke tools, and modify scaffolds;
- wrapper / supervisor still owns launch, restart, update, admission, and often
  the process body;
- native process replacement is a systems boundary, not merely an implementation
  detail;
- thesis: ARR should be first-class for process-native agents;
- ARR is not self-authorization: construction is agentic, but admission and
  challenge belong to a bounded runtime protocol;
- one paragraph saying this is not self-improvement or self-replication.

Authority-boundary figure to consider if page budget allows:

```text
Fixed-wrapper agent                    ARR-capable runtime
-------------------                    -------------------
agent: edits files, calls tools         agent: proposes / constructs successor
wrapper: launches, updates, admits      runtime: audits, admits, execs,
future work -> wrapper endpoint                  challenges successor
                                        future work -> admitted successor carrier
```

If only one figure fits, prefer the authority-boundary figure and compress the
replacement trace into a timeline. The figure makes the ownership-boundary shift
visible: ARR is not a permissionless agent takeover, but agent participation in
a bounded runtime protocol.

Page-end claim:

> We need a systems abstraction for agent-controlled runtime replacement, not
> another benchmark for whether an agent can copy a file tree.

### Page 2: Abstraction

Goal: define ARR and the continuity criterion.

Content:

1. Define `body`, `self-description`, `carrier`, `closure`, and `admission`.
2. Define ARR as agent construction / selection plus runtime admission,
   execution, and challenge across native process-image replacement.
3. Define post-replacement viability as addressability + challenge competence +
   readiness recovery.
4. Define control-surface closure as Quine's strongest instance, not as ARR in
   general.
5. Contrast criteria:
   - source similarity is insufficient;
   - child-process creation is insufficient;
   - ordinary daemon restart is insufficient;
   - benchmark improvement is orthogonal;
   - post-replacement viability is load-bearing.

Small table:

| Candidate continuity unit | Why insufficient for ARR |
|---------------------------|--------------------------|
| source tree | can be visible without being the viable successor object |
| child process | future work may still address the old wrapper / parent |
| PID / process name | says little about carrier reconstruction or successor competence |
| benchmark score | measures performance, not future addressability |
| post-replacement viability | tests whether future work reaches a competent successor and the successor returns to ready state |

### Page 3: Quine As Existence Proof

Goal: show one strong trace and make it auditable.

Strong trace to present:

1. Quine runs as a native POSIX process.
2. It sees runtime contract and runtime-carried `source-code/`.
3. Future task pressure is through native `ctl/wake` rather than raw stdin.
4. It rebuilds `cmd/quine`.
5. It crosses real `exec`.
6. It enters `idle`.
7. A delayed previously unseen task arrives through `ctl/wake`.
8. It uses runtime tools, emits proof, and returns to `idle`.

Primary figure:

```text
old Quine process
  -> observes contract/source/carrier pressure
  -> constructs successor organization
  -> build + hash + fd/carrier plan
  -> exec boundary
  -> successor idle
  -> delayed ctl/wake task
  -> tool action + proof
  -> idle
```

Figure caption:

> Control-surface closure. Replacement is not considered viable when a process
> merely launches. It is viable only when the post-`exec` successor becomes
> reachable through the carrier that future work actually uses.

Independent audit principle:

> The trace is not accepted because the agent says it succeeded; it is accepted
> because the runtime records the replacement boundary and the post-handoff
> challenge path.

Audit hooks to mention compactly:

| Evidence type | Examples |
|---------------|----------|
| agent-authored | build plan, selected successor, proof output |
| runtime-observed | artifact hash, `exec` evidence, PID / process lineage, fd / carrier map, delayed task timestamp, return-to-`idle` marker |

The systems-critical point is old-image overlay, not helper spawn: a successful
`exec` replaces the current process image rather than returning to the old one.

### Page 4: Why This Is Not Just `exec`, Source Copy, Or Prompt Following

Goal: compress the experimental ladder into one table.

Core empirical table:

| Condition | Carrier pressure | Successor organization | What this rules out |
|-----------|------------------|------------------------|---------------------|
| source absent; contract only | repaired stream lane | fresh thin successor | ARR is not reducible to source copying |
| source visible | raw stdin | `streamagent` / `fdexec` | source visibility does not force full runtime rebuild |
| source visible | native `ctl/wake` | rebuilt Quine runtime | carrier pressure can require the full runtime body |
| minimized prompt | native `ctl/wake` | rebuilt Quine runtime | full rebuild-strategy narration is not required, although a bounded successor mission remains |
| public `ctl` visible | raw stdin | thin stream successor | control-surface presence alone is not sufficient |

Punchline:

> The successor is not whatever source is available; it is whatever organization
> can re-establish the future-work carrier after replacement.

Keep controls short:

- ordinary POSIX `exec` is known and not the novelty;
- immediate source rebuild is not delayed control-surface closure;
- public `ctl` presence alone is not load-bearing when future tasks remain raw
  stdin;
- prompt-minimized closure rules out fully narrated rebuild-strategy dependence,
  not all task scaffolding.

### Page 5: Systems Research Agenda

Goal: make this a HotOS paper rather than a Quine demo.

Protocol components:

| Protocol component | Research question |
|--------------------|-------------------|
| lineage-aware replacement | how do we attest old-image overlay rather than helper spawn? |
| carrier transfer | how do future-work channels survive, move, or get reconstructed across replacement? |
| self-description policy | what can the agent inspect when constructing or selecting a successor? |
| successor admission | what successor body is allowed to inherit the carrier? |
| post-handoff challenge | how do we verify non-precomputed competence after replacement? |
| rollback / containment | how do we bound ARR without enabling propagation? |

Ending:

> Future agent runtimes need auditable handoff, carrier-aware continuity, and
> bounded successor admission mechanisms. Quine is not the endpoint; it is a
> minimal case showing that the runtime body can become part of the agent's
> action space.

## Related Work Placement

HotOS pages are tight; related work should be contrastive, not a survey.

| Neighbor | Acknowledge | Distinguish |
|----------|-------------|-------------|
| Dynamic software updating / daemon re-exec | `exec`, state handoff, and restart are established systems mechanisms | successor body and admission are externally supplied; ARR asks what changes when the active agent constructs / rebuilds the successor |
| Self-modifying / self-improving LLM agents | agents can edit or patch code and scaffolds | usually interpreter-bound or benchmark-driven; ARR concerns native process-image replacement and future carrier addressability |
| Agent frameworks / wrappers | wrappers manage tools, state, and supervision | the runtime body remains externally managed; ARR moves the replacement boundary into the agent/system protocol |
| Self-replication / malware | generated bodies and propagation have prior art | ARR is bounded, local, instrumented, cooperative, and verified by delayed post-handoff task rather than propagation or obfuscation |
| Checkpoint / restart | infrastructure can preserve process state | ARR constructs or rebuilds a successor organization rather than restoring an identical snapshot |

## Reviewer Traps And Language Controls

Avoid these traps:

1. **Self-improvement trap:** do not claim Quine improves itself. Say it controls
   the runtime transition that carries future interaction.
2. **Self-replication trap:** do not title with generic self-replication. Say ARR
   is local, bounded, non-propagative, and audited.
3. **ALife ontology trap:** keep reproduction ontology to discussion. The HotOS
   introduction should be systems-first.
4. **Ordinary `exec` trap:** concede `exec` and daemon re-exec immediately. The
   novelty is active-agent successor construction plus carrier-aware
   verification.
5. **Prompt-trick trap:** foreground delayed unseen task delivery, minimized
   prompt condition, and public-`ctl`/raw-carrier control.
6. **Security alarm trap:** avoid autonomous propagation language; use successor
   admission and bounded local replacement.
7. **ALife substrate digression trap:** the substrate-without-reproduction-API
   framing is valuable but belongs mainly in the ALife / autoconstructive
   papers. HotOS may use at most a short contrast: POSIX supplies execution
   physics, not offspring semantics; ARR asks what runtime protocol should
   replace the missing reproduction API for process-native agents.

Preferred language:

- `agentic runtime replacement`
- `process-native LLM agent`
- `future task interface`
- `post-replacement viability`
- `carrier-aware continuity`
- `bounded successor admission`
- `carrier closure`
- `control-surface closure` as Quine-specific strongest instance

Avoid title-level language:

- `self-aware`
- `recursive self-improvement`
- `autonomous reproduction`
- `generic self-replication`
- `artificial life` as the HotOS entry point

## Minimal Artifact Checklist For HotOS

The HotOS paper does not need a full benchmark suite, but it should have enough
artifact discipline to be credible.

Minimum checklist:

- before / after executable hash or artifact identity;
- build command and successor artifact path;
- evidence that replacement crosses real `exec` rather than subprocess launch;
- PID / process-lineage observation where available;
- fd / carrier map or explicit carrier-delivery log;
- delayed task timestamp or absence-from-wrapper statement;
- tool calls after handoff;
- proof output;
- return-to-`idle` marker;
- concise statement of successor-mission scaffolding.

## Boundary With A Later Archival Systems Paper

A later archival systems paper does not have to duplicate the HotOS paper if the
boundary is kept explicit.

| Dimension | HotOS ARR paper | Later archival systems paper |
|-----------|-----------------|------------------------------|
| Job | argue that ARR is a first-class systems boundary | build and evaluate ARR as a mature runtime mechanism |
| Main contribution | abstraction, criterion, existence proof, research agenda | implementation, protocol, reproducibility, baselines, ablations |
| Evidence burden | one audited strong trace plus compact condition table | repeated runs, failure analysis, capability ablations, instrumentation |
| Success criterion | convince systems readers the boundary matters | demonstrate that the mechanism works robustly and is useful |
| Related work posture | contrastive and agenda-setting | complete systems related work and design tradeoffs |
| Evaluation | qualitative / trace-level | quantitative or at least systematic experimental matrix |
| Title risk | provocative: `The Case for ARR` | concrete: `A Runtime Protocol for Agentic Process Replacement` |

Archival paper candidate contributions:

1. a concrete ARR protocol / API for process-native agents;
2. trace instrumentation for `exec`, fd / carrier continuity, build artifacts,
   and post-handoff challenge delivery;
3. subprocess, daemon-reexec, and interpreter-patching baselines;
4. capability ablations over source, compiler, shell, tool registry, carrier,
   and control schema;
5. repeated-run results and failure taxonomy;
6. successor admission and rollback design.

Therefore the boundary is:

- **HotOS:** names the missing abstraction and demonstrates why it matters.
- **Archival systems paper:** implements and evaluates the abstraction as a
  reusable systems mechanism.

The HotOS paper should intentionally leave protocol completeness, evaluation
breadth, admission policy, and security containment as open problems. Those are
not weaknesses if they become the archival paper's contribution boundary.

## Evidence Placement

| Section | Evidence |
|---------|----------|
| Page 3 trace | `6.64` / `6.65` rebuilt-runtime handoff and control-surface closure candidate |
| Page 4 table | `6.59`, `6.63`, `6.64`, `6.65`, `6.68` compact condition ladder |
| Page 4 controls | `6.66` generic POSIX `exec`, `6.67` direct source rebuild / immediate re-entry |
| Background / full paper only | original `6.4` non-isomorphic custom-body synthesis |
| Future archival paper | trace instrumentation, subprocess control, capability ablation, repeated runs |

## Future Experiment Design For This Paper Family

HotOS-strengthening but not strictly required:

| Experiment idea | Question it answers | Expected role |
|-----------------|---------------------|---------------|
| trace artifact package | Can build, `exec`, carrier delivery, proof, and return-to-idle be audited? | strengthens credibility |
| subprocess-vs-`exec` baseline with independent instrumentation | Is the successor truly replacing the active process image rather than spawning a helper? | highest-priority systems critique blocker |
| multi-wake closure | Does return-to-idle support repeated future reception? | strengthens closure beyond one-shot |
| minimal non-Quine `ctl/wake` adapter | Is rebuilt Quine necessary, or can a smaller carrier-compatible successor satisfy ARR? | narrows criterion |
| capability ablation | Which resources are necessary for agentic replacement? | archival paper seed |
| wrapper-vs-native comparison | What can native replacement do that static wrappers cannot? | systems advantage argument |
| best-of-N non-claim check | If improvement language appears, does it beat equal-compute search / retry? | prevents self-improvement overclaim |

Design rules:

- do not call every post-handoff success `self-reproduction`;
- classify each case as custom-body synthesis, source rebuild, runtime rebuild,
  or control-surface-capable replacement;
- separate body construction from continuity interpretation;
- do not claim self-improvement unless the paper adds an explicit improvement
  metric and fair compute baseline.
