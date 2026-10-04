# Manuscript Plan: Protocol Emergence over Minimal Runtime Surfaces

**Working title:** `Protocol Emergence over Minimal Runtime Surfaces`  
**Role:** systems-facing empirical paper on protocol-layer emergence  
**Status:** planning

## Core Thesis

When a runtime provides only message-delivery physics and a thin observable
process surface, agents can still construct protocolized coordination
conventions above that substrate. The current paper should make a bounded
claim: Quine has evidence for protocolized coordination and some
routing-relevant structure, but not yet for a broad theorem of robust callback
invention.

## Claim Boundary

- **Positive claim:** minimal runtime transport can induce protocolized
  coordination conventions above the runtime layer
- **Stronger but partial claim:** some observed conventions are
  routing-relevant or delegation-relevant
- **Negative claim:** the exact combined setting of opaque payloads, no sender
  metadata, no reply channel, no predefined schema, and long-lived async
  processes appears under-studied
- **Non-claim:** Quine has already proved general callback or correlation-id
  emergence

## Paper Split Discipline

- the current archival-candidate manuscript is the bounded
  `protocolized-coordination over minimal runtime surfaces` paper
- the stronger callback / correlation paper is a follow-on that depends on the
  `L5` process-surface callback fixture becoming a stable anchor
- `8.9`, `8.10`, `8.17`, `8.12`, and compressed `2.3` should therefore carry
  the current manuscript backbone
- `L5` and later add-back contrasts should be framed as the next-step causal
  program unless they become stable enough to re-center the paper

## Evidence Labels

The manuscript should use one canonical evidence vocabulary:

- `coordination-only`: cooperation or closure without concrete routing or
  return structure worth elevating as protocol evidence
- `routing-relevant`: recipient, path, pid, mailbox, or other locator
  structure that helps direct work or returned results
- `delegation-relevant`: explicit task handoff or return discipline stronger
  than generic cooperation wording, but not yet a verified callback path
- `callback-relevant`: a peer uses an agent-constructed return path back to a
  requester-owned process surface
- `correlation-relevant`: structure distinguishes multiple requests, replies,
  or delegation branches
- `ambiguous / prior-like`: protocol-looking formatting that may be better
  explained by pretraining priors, fixed task structure, or generic
  cooperation wording

## Reader Journey

The manuscript should disclose one layer at a time:

1. agent research usually assumes communication semantics in advance
2. the missing layer is protocol emergence, not generic communication
3. Quine strips sender metadata, reply channels, and typed request/response
   contracts while preserving only minimal transport-like surfaces
4. address-space exposure can remain visible; that is substrate geometry, not
   protocol semantics
5. the empirical case is best presented as an **evidence staircase**
6. the bounded result is protocolized coordination with some routing-relevant
   structure
7. stronger callback claims still require causal perturbation experiments

## Section Plan

### 1. Introduction

Goals:

- open with the broad gap: most agent work studies coordination over predefined
  communication interfaces
- narrow immediately to the missing layer: protocol-layer coordination itself
  under minimal runtime transport
- state the exact runtime pressure:
  opaque payloads, no sender metadata, no built-in reply channel, no predefined
  schema, long-lived asynchronous agents
- preview the bounded thesis and the manuscript contribution order

Contribution preview order:

1. define the missing problem layer
2. position adjacent literatures without collapsing them together
3. present a bounded evidence staircase from control to emergence to boundary
   contrast
4. identify the causal tests still needed for stronger callback claims

### 2. Runtime Model and Target Object

Goals:

- distinguish `transport physics`, `protocol semantics`, and `coordination
  outcome`
- specify what the runtime gives:
  delivery surfaces, process-visible artifacts, opaque text payloads,
  asynchronous persistence
- specify what it withholds:
  sender identity, reply channel, typed schema, built-in request/response
  semantics
- define what counts as evidence of protocol-layer structure:
  callback locator, routing reference, correlation marker, delegation return
  convention, lightweight field grammar
- define what does **not** count by itself:
  generic blackboard coordination, polite cooperation wording, or turn-order
  reconstruction

### 3. Related Work

Structure by non-equivalence rather than chronology.

Bucket 1: emergent communication / MARL

- signaling, targeting, and communication strategy over predefined channels
- typically not about reply semantics or protocol-layer structure under minimal
  transport

Bucket 2: blackboard / Linda / shared-surface coordination

- important substrate neighbors
- coordination semantics are usually designed into the system rather than
  invented by agents

Bucket 3: agent interoperability / framework protocols

- A2A, MCP, handoff systems, task/message frameworks
- explicit protocol design, not protocol construction under absent semantics

Bucket 4: ALife / protocol evolution / fraglets / Avida

- strongest conceptual neighbors because they study message grammar and
  protocol-like structure
- still not direct matches for long-lived asynchronous LLM agents over sparse
  runtime surfaces

Close with the exact negative result:
the combined Quine setting appears under-studied.

### 4. Evidence Staircase

Use argumentative roles, not a fake monotonic ablation story.

| Rung | Experiment | Role | Main read |
|------|------------|------|-----------|
| A | `8.9` | capability control | mail surface is mechanically usable when explicitly instructed |
| B | `8.10` | visible-runtime spontaneous uptake | `ctl + USR1 + peer-surface inspection` can become part of coordination without task-level signal obligation; strong but intermittent |
| C | `8.17` | zero-explicit-hint protocol cluster | strongest clean main-text anchor for bounded protocol-emergence claims |
| D | `8.12` | hidden-runtime blackboard floor | boundary contrast showing coordination can persist without visible runtime routing surfaces |
| E | `2.3` | historical precursor only | early file-native coordination precedent, not a core rung |

Section goals:

- separate emergence evidence from capability controls
- separate protocolized coordination from generic blackboard success
- keep `8.17` as the strongest clean anchor for the bounded claim
- use `8.12` to prevent overclaiming signal or visible routing necessity
- label manuscript-facing artifacts with the canonical evidence vocabulary
  rather than treating all protocol-looking structure as equally strong

### 5. What the Current Evidence Supports

Supported now:

- minimal runtime surfaces can induce protocolized coordination conventions
- some observed conventions are routing-relevant or delegation-relevant

Not yet supported:

- general callback emergence
- correlation-id necessity
- model-general stability
- strong causal claims that removing an invented field breaks success

This section should also hold the novelty boundary and threats-to-validity
material.

### 6. Next-Step Methodology

Specify the missing experiment classes for a stronger paper:

1. **Protocol perturbation**
   Remove or scramble invented mailbox, pid, return-path, or session-like
   fields and measure failure.
2. **Built-in semantics ablation**
   Add sender metadata or a reply channel and test whether invented
   conventions simplify or disappear.
3. **Stronger protocol-pressure task**
   Use asynchronous multi-outstanding delegation where return routing and
   correlation are necessary, not stylistic.
4. **Adversarial boundary tests**
   Probe spoofing or impersonation to see whether invented protocol structure
   carries authentication-like behavior.

### 7. Conclusion

- restate the bounded result:
  protocol-layer conventions can emerge over minimal runtime surfaces
- keep the close conservative:
  current evidence is partial, routing-relevant, and experimentally extensible
- point directly to the stronger causal program rather than overselling the
  present corpus

## Gap Assessment

For a bounded paper on `protocolized coordination over minimal runtime
surfaces`, the gap is moderate rather than catastrophic. The current corpus
already contains:

- `8.9` as control
- `8.10` as spontaneous visible-runtime uptake
- `8.17` as the strongest zero-hint protocol cluster
- `8.12` as the boundary contrast

Still missing even for the bounded version:

- a clean manuscript-facing comparison table across the staircase
- transcript extraction of the strongest protocol artifacts
- disciplined labeling of evidence as coordination-only, routing-relevant,
  delegation-relevant, callback-relevant, correlation-relevant, or ambiguous

For a stronger paper on callback / correlation emergence, the gap remains
large: the current experiments do not force callback semantics strongly enough
and do not yet include causal removal tests.

## Evidence Placement

- `8.9` belongs in the staircase as a control, not as novelty
- `8.10` and `8.17` are the main empirical anchors
- `8.12` belongs in the boundary section to show protocol emergence is not the
  same thing as generic cooperation
- `2.3` should stay historical and compressed
- the literature matrix belongs in related work plus a data pointer
- callback / routing claims should cite transcript-level excerpts only after
  extraction from the retained anchor runs
