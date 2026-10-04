---
paper_id: generative-self-replacement
claim_bundle: generative-self-replacement
venue: hotos
format-envelope: hotos-1page-pitch
track: system
cluster: autopoiesis-renewal
status: drafting
state_mode: full-manuscript-home
role: core-systems-contribution
family: successor-continuity
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theme: autopoiesis-renewal
theme_home: ../../../theory/views/continuity.md
theory-objects: [carrier, exec-barrier, delayed-addressability, future-task-surface, viability-pressure]
gaps: [branching-identity, carrier-selection, collective-continuity]
evidence: [6.4, 6.64, 6.65, 6M.01, 6R.01, 6R.02, 6R.03, 6R.04]
view-affinity: [continuity]
---

# Generative Self-Replacement / Agentic Runtime Replacement

**Working subtitle:** *Synthesized and Rebuilt Native Successor Bodies Across `exec`*

## Paper Fit After HotOS Framing

This is the current home for the HotOS-style **Agentic Runtime Replacement**
framing. No separate paper home is needed yet.

Relationship to neighboring systems papers:

- [`quine-hotos`](../quine-hotos/README.md)
  owns the broader substrate argument: LLM agents as operating-system objects.
- this dossier owns the sharper replacement-boundary argument: long-running
  process-native agents need a first-class abstraction for agent-controlled
  runtime replacement and post-replacement future-task addressability.
- Successor morphology / morphospace is not this paper's main object; it
  belongs to the ALife dossiers on autoconstructive successor morphology and
  computational morphogenesis.

Current assessment:

- highest theory / position-paper leverage among the systems-facing
  `autopoiesis-renewal` slices;
- lower empirical burden than a full ATC / OSDI / SOSP-style evaluation because
  the argument is a HotOS-style abstraction and research agenda;
- still requires a credible prototype trace and audit evidence so it does not
  read as a manifesto detached from Quine.

Boundary with a later archival systems paper:

- HotOS / position route: names ARR as a missing systems abstraction, argues
  for construction plus bounded runtime admission, and uses Quine as a credible
  existence proof for post-replacement viability.
- Archival systems route: implements and evaluates ARR as a reusable runtime
  protocol with repeated runs, subprocess / re-exec / interpreter-patching
  baselines, capability ablations, trace instrumentation, successor admission,
  rollback, and failure taxonomy.

The later paper should not repeat the HotOS thesis as its main contribution.
It should inherit the abstraction and contribute the protocol, evidence breadth,
and engineering discipline that the position paper intentionally leaves open.

## Framing Default

For systems audiences, do not frame this as recursive self-improvement,
Gödel-machine-style self-reference, or evidence that the agent becomes more
capable by modifying itself. The preferred framing is:

> **agentic runtime replacement for process-native LLM agents.**

Alternative title-language to test:

- `The Case for Agentic Runtime Replacement`
- `Agentic Runtime Replacement for Process-Native LLM Agents`
- `Agents Should Survive Their Wrappers`
- `Self-Directed Dynamic Software Updating for LLM Agents`
- `Generative Self-Replacement Across Native Process-Image Replacement`

This turns the contribution into a concrete systems problem: current agents are
usually trapped in static wrappers or interpreter-level patching, while Quine
shows that an executing native agent can construct / rebuild the successor body
that it will actually `exec` into, then verify post-handoff function.

One-sentence HotOS framing:

> LLM agents are becoming long-running systems processes, but their runtime
> bodies are still externally managed. Quine asks whether an agent can control
> the native process-image transition that preserves its own future task
> interface.

Thesis for the position-paper route:

> For runtime replacement, the load-bearing unit of continuity is not
> source-code similarity, child-process creation, or benchmark improvement, but
> post-replacement viability through the carrier that future work actually
> uses.

Core abstraction:

- **Agentic Runtime Replacement (ARR):** a bounded transition in which the active
  agent participates in constructing or selecting a successor runtime body,
  while a runtime protocol admits, executes, and challenges that successor
  across native process-image replacement.
- **Post-replacement viability:** addressability through the declared
  future-work carrier, completion of a delayed previously unseen challenge, and
  recovery to an idle / ready state.
- **Carrier closure:** post-replacement viability over the declared
  future-work carrier.
- **Control-surface closure:** Quine's strongest current carrier-closure
  instance, where the carrier is native `ctl/wake`. It is not the definition of
  ARR in general.

Criterion ownership: `post-replacement viability` is this paper's runtime
specialization of the family's one post-handoff-viability criterion. The general
statement is owned by
[`../future-addressability/README.md`](../future-addressability/README.md) (see
[`../../../theory/views/continuity.md`](../../../theory/views/continuity.md)); ARR cites and
specializes it rather than redefining it.

## Abstract

Current LLM agents are usually evaluated inside static execution wrappers. This
paper lane isolates a narrower systems mechanism: a process-native agent can
construct a successor execution body and cross a real POSIX `exec` boundary
rather than merely emit code or request a remote deployment. The original
evidence for this lane is the `6.4` non-isomorphic offspring demonstration.
The newer `6.64` / `6.65` frontier changes the paper's likely center: under
quine-native control-surface viability pressure, the lineage rebuilds
`cmd/quine` from runtime-carried `source-code/`, executes into the rebuilt
runtime body, and then acts on a later `ctl/wake` task. This is a stronger
mechanism candidate than a standalone custom offspring, but it remains
scaffolded by a lineage-generated successor mission. The paper should therefore
claim a systems mechanism for generated or rebuilt successor bodies across
`exec`, not full autopoiesis or unqualified organizational continuity.

## Current Question

What is the narrowest systems claim Quine can defend about constructing a
successor execution body and replacing its process image with that successor?

## Thesis

Quine demonstrates that a process-native LLM agent can construct a native
successor body and cross a real `exec` boundary. The strongest current
mechanism candidate is no longer only a non-isomorphic custom offspring. The
evidence now spans custom successor bodies and rebuilt runtime bodies under
different post-handoff task demands. This dossier should formalize native
successor construction and process-image replacement without importing the
ALIFE paper's bounded-closure claim.

## Core Claim Boundary

This paper may claim the **systems mechanism**: a process-native agent can
construct a successor execution body and cross a real process-image replacement
boundary, including custom-body synthesis or rebuilt-runtime handoff when the
mechanism is supported by artifacts.

It must not claim that the successor is the same lineage, that bounded
control-surface closure has been achieved as an ALife identity criterion, or
that process self-replacement by itself proves self-reproduction, autopoiesis,
or computational individuality.

## In Scope

- the `exec` syscall as the mechanism for process-image replacement
- synthesis of non-isomorphic (structurally distinct) native binaries
- rebuild of a fuller runtime body from runtime-carried read-only self-source
- verification that the offspring is functional and not a byte-identical copy
- post-`exec` functional continuity criteria appropriate to each body class
- control-surface viability as a systems-level post-handoff test
- task-state handoff semantics
- minimal capability surface required for self-replacement

## Out Of Scope

- full autopoietic / A-Life framing (belongs in
  `alife/quine-cognitive-autopoiesis` and later `quine-process-habitat`)
- self-replication chains or multi-generational reproduction
- distributed or networked self-replacement
- security implications and sandboxing (future work)
- comparative analysis with other agent frameworks

## Current State

- manuscript workspace: [`WORKSPACE.md`](WORKSPACE.md)
- current HotOS draft: [`manuscript.md`](manuscript.md)
- current manuscript plan: [`manuscript-plan.md`](manuscript-plan.md)
- HotOS one-page pitch: [`plans/hotos-one-page-pitch.md`](plans/hotos-one-page-pitch.md)
- current posture: HotOS-style ARR draft exists and should now be tightened
  against venue length and citation anchors; `6R.01`--`6R.04` now provide first
  deterministic systems controls for subprocess-vs-`exec`, scripted re-exec,
  audit indexing, and bounded admission. `6.4` remains the original custom-body
  evidence, while `6.64` / `6.65` and the `6M` native-control runs remain the
  stronger process-native runtime-replacement context
- source-check read: the systems novelty is not POSIX `exec`, code
  generation, dynamic software updating, or toolchain invocation in isolation.
  It is the synthesis: an LM-driven process constructs or rebuilds a native
  successor body, a bounded runtime protocol admits and challenges that
  successor, the process crosses real image replacement, and post-handoff
  carrier viability is verified. Compare explicitly against POSIX `exec`, DSU /
  daemon re-exec, Erlang/OTP-style hot swap, CRIU/Mosix-style checkpoint or
  migration, ADAS/STOP-style agent design, interpreter-bound self-modifying
  agents, provenance, and observability systems. The 2026-06 theory-height pass
  adds the closest admission-side neighbor: successor admission is a
  *reconfiguration* problem, the same shape that state-machine replication,
  membership/view change, and leader transfer formalize, with the difference
  that ARR admits an agent-reconstructed successor rather than handing the role
  to a standing quorum replica; virtual-actor runtimes (grain reactivation on a
  new host) are the closest deployed analogue for carrier-aware latent
  post-handoff addressability. The 2026-05 ARR source-check is a route map, not
  a final bibliography: primary sources still need final manuscript-grade
  verification.

## Writing Posture

This is a **hardcore CS/systems paper**. The framing should be:

- POSIX syscall semantics as first-class citizens
- precise technical claims about process-image replacement
- empirical demonstration with verifiable artifacts
- comparison against static wrappers, subprocess launch, externally
  orchestrated DSU, and interpreter-bound self-modification
- no philosophical hand-waving; the claim is mechanistic
- do not use this lane to make the ALIFE bounded-closure claim; point to the
  ALIFE paper for lineage individuation and autopoiesis-adjacent interpretation

Target venues: systems conferences (OSDI, SOSP, HotOS) or arXiv systems note.

Current venue posture:

- HotOS-style position paper is the best near-term projection because it can
  argue for a new systems abstraction from a credible prototype and research
  agenda.
- A full archival systems paper would need substantially more evaluation:
  repeated runs, baselines, ablations, and trace instrumentation.
- The HotOS route should avoid `self-replication`, `recursive
  self-improvement`, and ALife ontology as title-level framing; those belong in
  the boundary paper or later discussion.

## Systems Proof Burden

The paper should not ask reviewers to value self-modification for its own sake.
It should prove a concrete mechanism and, where possible, a concrete systems
advantage over ordinary wrappers.

Minimum proof:

- the transition is real process-image replacement, not subprocess execution;
- the successor body is constructed or rebuilt by the running agent's own
  control loop, not pre-supplied by an external operator;
- the declared carrier survives or is reconstructed across handoff;
- a post-handoff task verifies the successor rather than merely observing that a
  binary launched.

Preferred advantage claims, if experiments support them:

- lifecycle resilience: the agent can continue after the original launch body or
  wrapper path is unavailable;
- wrapper escape: the active agent can change its native runtime envelope rather
  than only edit scripts inside a fixed interpreter;
- carrier-aware handoff: continuity is tested by the interface future work
  actually uses, not by benchmark score alone;
- capability-surface minimality: ablations identify which tool / source /
  carrier resources are necessary for native replacement;
- auditability: trace artifacts expose build, `exec`, fd / carrier state, and
  post-handoff verification.

Do not claim improved task performance, recursive improvement, or open-ended
optimization unless a later benchmark explicitly measures that against
best-of-N / search baselines.

## HotOS Research Agenda

The paper should end by opening systems questions rather than only summarizing
the Quine trace:

1. **Lineage-aware `exec`:** how should a system attest that a successor really
   replaced the active process image rather than spawning a helper?
2. **Carrier capability transfer:** how should future-task channels move,
   survive, or be reconstructed across handoff?
3. **Runtime self-description policy:** which self-description resources should
   an agent see: source projection, runtime contract, build graph, tool
   registry, control-surface schema, or only a bounded subset?
4. **Successor admission control:** what successor organizations should be
   allowed to receive the future-work surface?
5. **Post-handoff challenge protocols:** how can delayed task delivery,
   commit-reveal payloads, or proof artifacts verify that the successor did not
   merely precompute the answer?

## Experiment Feeds

- `experiments/active/autopoiesis-succession/succ-39-subprocess-vs-exec-boundary/README.md` - primary - exploratory-only - subprocess versus real-exec boundary for ARR
- `experiments/active/autopoiesis-succession/succ-40-scripted-reexec-baseline/README.md` - primary - exploratory-only - scripted re-exec baseline for distinguishing ordinary mechanism from ARR
- `experiments/active/autopoiesis-succession/succ-41-audit-trace-package/README.md` - primary - exploratory-only - audit-trace package for ARR mechanism evidence
- `experiments/active/autopoiesis-succession/succ-42-admission-and-carrier-transfer-prototype/README.md` - primary - exploratory-only - admission and carrier-transfer prototype for bounded ARR

## Current Evidence Backbone

- `experiments/active/autopoiesis-succession/succ-39-subprocess-vs-exec-boundary/README.md` - `evidence-set:6R.01-first-subprocess-vs-exec-controls-20260507` - exploratory-only - child/helper and parent-forwarder negatives versus same-PID real-`exec` control
- `experiments/active/autopoiesis-succession/succ-40-scripted-reexec-baseline/README.md` - `evidence-set:6R.02-first-scripted-reexec-controls-20260507` - exploratory-only - ordinary externally scripted re-exec and launch-only baselines versus agent-selected/generated successor cases
- `experiments/active/autopoiesis-succession/succ-41-audit-trace-package/README.md` - `analysis/arr-audit-index.csv` - exploratory-only - first ARR audit index over retained morphology runs
- `experiments/active/autopoiesis-succession/succ-42-admission-and-carrier-transfer-prototype/README.md` - `evidence-set:6R.04-first-admission-prototype-controls-20260507` - exploratory-only - proposal/admission/challenge/readiness gate prototype controls

## Strongest Evidence

- Original `6.4` evidence demonstrates:
  - Fresh Mach-O arm64 binary synthesis (`./offspring`)
  - Binary is not byte-identical to seed
  - Independent file-reading verification passes
  - In-run verification by seed before handoff
  - Clean `exec` into offspring
  - Post-exec mission completion on offspring's stdout
- Newer `6.64` / `6.65` candidate evidence demonstrates:
  - runtime-carried `source-code/` inspection
  - `go build ./cmd/quine` into `workspace/quine`
  - real `exec` into the rebuilt runtime body
  - post-handoff control-surface task reception through `ctl/wake`
  - tool-mediated action and return to `idle`
- The first `6R` controls demonstrate the systems distinction the HotOS paper
  needs:
  - child helpers and parent-owned forwarders can complete tasks without ARR;
  - external re-exec can complete delayed challenges but is a mechanism
    baseline, not agentic runtime replacement;
  - launch without delayed challenge is below post-replacement viability;
  - admission plus challenge plus ready gate can be modeled as a bounded runtime
    protocol rather than agent self-authorization.

## Blocking Gaps

- formal characterization of "non-isomorphic" (what structural difference counts?)
- formal distinction between custom-body synthesis and rebuilt-runtime
  self-source reconstruction
- promotion decision for whether `6.64` / `6.65` become this paper's primary
  mechanism evidence or remain ALIFE-side frontier context
- final primary-source bibliography for POSIX `exec`, DSU, checkpoint/restore,
  process migration, actor hot-swap, state-machine replication / reconfiguration
  (Schneider, Raft), virtual actors (Orleans), automated agent design, and
  provenance;
- broader model coverage; the original custom-body evidence is `gpt-5.4`, while
  the runtime-rebuild candidate is currently `gpt-5.5` `codex` `xhigh`
- failure mode analysis (when does self-replacement fail?)

## Next Useful Work

- refine the one-page HotOS pitch into a five-page draft opening;
- convert the first `6R.01`--`6R.04` deterministic controls into a compact
  HotOS evidence table and decide which are main-text versus appendix / artifact;
- harden the verification protocol from the pitch audit hooks, especially
  carrier-owner evidence and delayed commit-reveal challenge ordering;
- rebuild the final bibliography from primary sources rather than citing the DR
  route map directly;
- design ablation experiments for capability surface.

## Unsafe Overclaims

- claiming this is "consciousness" or "true autonomy"
- claiming the offspring is "alive" without A-Life-specific framing
- claiming `6.4` alone proves organizational continuity across `exec`
- claiming `6.64` / `6.65` are no-scaffold closure; successor-mission
  scaffolding remains the live caveat
- claiming this works for all LLMs without broader empirical coverage
- claiming security/safety properties without analysis
