---
paper_id: future-addressability-computational-individuality
claim_bundle: future-addressability-computational-individuality
venue: cognitive-systems-research-journal
format-envelope: full-journal
track: theory
cluster: autopoiesis-renewal
status: planning
state_mode: planning-skeleton
role: computational-individuality-theory-paper
family: successor-continuity
publication_class: theory-candidate
maturity_target: theory-anchor
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theme: autopoiesis-renewal
theme_home: ../../../theory/views/continuity.md
theory-objects: [carrier, exec-barrier, delayed-addressability, future-task-surface, trace-inheritance]
gaps: [branching-identity, collective-continuity, carrier-selection]
evidence: ["6.63", "6.64", "6.65", "6.66", "6.67", "6M.10", "6F.01"]
view-affinity: [continuity]
---

# Computational Individuality Is Indexed By Future Addressability

## Current Question

When source code, runtime body, process identity, persistent state, address,
and future task carrier come apart across a boundary event, what should count
as the continuing computational individual?

## Thesis

Computational individuality is not fixed by source code, executable identity,
process identity, internal state, path name, or artifact continuity alone.
It is indexed by the future-task carrier through which the system can still be
addressed, woken, obligated, and made to act after a boundary event.

Short form:

> **Future addressability indexes computational individuality.**

Operational form:

> A successor is continuous with a computational individual only relative to a
> future ecology in which it remains reachable and competent through the
> carrier by which future obligations arrive.

## Core Claim Boundary

This paper may claim a **computational individuality criterion**: for
obligation-bearing computational lineages, continuity can be indexed by the
future carrier through which a successor remains addressable, competent, and
able to restore the receptive state required for later obligations.

It must not claim a universal metaphysical identity theory, a full
self-reproduction theory, biological life, generic distributed-systems identity,
or new positive Quine evidence beyond the motivating cases promoted by sibling
papers.

## Absorbed Theory Update

The recovered `dr` follow-up sharpens the paper's anchor into three linked
claims:

1. computational individuality is **carrier-indexed**
2. continuity is tested by **post-handoff addressability**
3. closure claims must be placed on a **scaffold gradient**

The useful absorption is conceptual, not citation-ready. The `dr` report
supplies vocabulary and neighboring literatures; primary sources still need to
be checked before manuscript citation.

The current best wording is:

> Computational individuality is not indexed by what a system copies, but by
> the future carrier through which it remains addressable and competent after a
> boundary event.

This is a criterion for a specific kind of computational continuity. It is not
a universal metaphysical identity theory.

## Second DR Round Absorption

The four-query `dr` round adds a sharper internal structure:

1. **future addressability is triadic**: it joins an obligation, a carrier, and
   a successor organization. Mere reachability is not enough; the reached
   successor must be competent for the ecology's future-obligation trace.
2. **carrier topology splits continuity types**: some systems preserve state
   and execution; others discard local state and reconstruct organization from
   source, event logs, mailboxes, capabilities, or control-plane state.
3. **formalization should stay lightweight**: the useful notation is a
   carrier-indexed predicate, not a full proof system.
4. **failure modes are theory-relevant**: substrate capture, state drift,
   irreversible degradation, and unbounded recursion are not only engineering
   risks; they mark where future addressability ceases to identify a coherent
   individual.

Refined criterion:

```text
P --B--> S

S is continuous with P under ecology E via carrier c only if:

1. S remains addressable through c after boundary B
2. S can satisfy E's future-obligation trace tau
3. S preserves the relevant carrier-indexed behavior from the observer's
   standpoint
4. S restores the receptive condition required for further obligations
```

For Quine:

```text
tau_Q = idle -> wake(task) -> tool/action -> report -> idle
```

This keeps the formal claim usable without pretending that a single process
calculus, session type, or bisimulation relation captures all of the real OS
and ecological scaffolding.

## Why This Paper Exists

The boundary paper shows a concrete `ALIFE` phenomenon: in Quine, raw-stream
future-task pressure can make a thin stream successor viable, while
Quine-native `ctl/wake` pressure makes a control-surface-capable runtime
organization viable.

In the current series split, this paper is the **generalization**:

```text
LBA = announcement
main Paper A = evidence
future-addressability = generalization
```

It should inherit the evidence standard from the boundary paper, but it should
not retell the full Quine ladder. Its job is to state the carrier-indexed
continuity criterion that the Quine case reveals and compare that criterion
against neighboring systems and formal vocabularies.

That result points beyond the `LLM self-reproduction` frame. It suggests a
more general theory of computational individuality:

- the relevant individual is not always the copied artifact
- the relevant body is not always the current executable
- the relevant continuity relation is not always memory preservation
- the future carrier can determine which organization must be preserved or
  reconstructed

This dossier exists to hold that theory as its own paper-sized object rather
than letting it blur the boundary paper, the evaluation paper, or the systems
mechanism paper.

## Current State

This is a `planning-skeleton` / `seed-dossier` inside the
`autopoiesis-renewal` theme. It currently owns the theory slice about
future-addressability-indexed computational individuality, but it has not yet
chosen a venue route or promoted the `6F.xx` experiment plan into a positive
evidence backbone.

The durable material already present is:

- a compact thesis and operational criterion
- absorbed vocabulary from the recovered and second-round deep-research passes
- a 2026-05-30 theory-height note that frames future addressability as a
  prospective selection criterion over candidate self-objects
- a 2026-06 theory-height pass that grounds the criterion in a five-family
  neighbor backbone with two load-bearing anchors (closest-continuer metaphysics
  and naming theory), maps the predicate onto name/address/route and coalgebraic
  bisimulation, and surfaces the branching/fission and corrigibility tensions
- a relation map to the boundary, evaluation, and systems-mechanism siblings
- a prospective experiment feed for turning the motivating Quine contrast into
  discriminating evidence

## Central Claim

In process-native and distributed computational systems, the "same
individual" relation is indexed by future addressability:

- if future tasks arrive by `stdin`, a stream-compatible successor can be the
  minimal continuing organization
- if future tasks arrive by actor mailbox, mailbox address and message-handling
  behavior become central
- if future tasks arrive by endpoint or capability, carrier topology becomes
  central
- if future tasks arrive by checkpoint scheduler, checkpoint image and restore
  protocol define the continuity object
- if future tasks arrive by Quine `ctl/wake`, the successor must preserve
  `wake -> act -> report -> idle`

The reproduced self-object is therefore carrier-indexed: it is the minimal
organization that preserves future addressability and action capacity under a
given future-task ecology.

## Relationship To Existing Papers

| Paper | Relation |
|-------|----------|
| [`../boundary-paper/README.md`](../boundary-paper/README.md) | provides the motivating Quine phenomenon: reproduced object, delayed addressability, and bounded control-surface closure |
| [`../../archive/superseded-2026/evaluation-paper/README.md`](../../archive/superseded-2026/evaluation-paper/README.md) (removed 2026-07-01; archived; superseded by `iclr-regime-decomposition` / `iclr-safeguard-composition`) | provided evaluation vocabulary for successor relations under real handoff constraints |
| [`../generative-self-replacement/README.md`](../generative-self-replacement/README.md) | provides the systems mechanism of generated / rebuilt successor bodies across `exec` |
| this dossier | extracts the general theory of computational individuality indexed by future addressability |

This paper should not retell the Quine evidence ladder in full. It should use
Quine as the strongest experimental case, then generalize the continuity
principle across OS processes, actor systems, distributed systems,
checkpoint/restart, capability systems, and artificial-life reproduction
semantics.

Its novelty is an indexing criterion, not a claim to have solved computational
identity in general. The manuscript must distinguish future addressability from
nearby concepts such as liveness, endpoint reachability, session traces,
bisimulation, capability reachability, actor identity, and ordinary service
continuity.

## In Scope

- computational individuality
- post-handoff addressability
- future-task carriers and carrier topology
- control-surface closure as one instance of addressability-indexed identity
- behavioral equivalence under a specified future ecology
- OS process identity, `exec`, file descriptors, and capabilities
- actor mailbox / endpoint / service-address continuity
- process migration and checkpoint/restart as contrast cases
- artificial-life individuation where offspring semantics are not given

## Out Of Scope

- proving full biological autopoiesis
- claiming Quine is alive
- claiming full autonomous agency
- claiming population-level evolutionary selection
- broad AI self-replication capability benchmarking
- systems implementation details of rebuilding successor bodies except as
  examples of carrier-indexed continuity

## Core Vocabulary

**Future addressability**

The property that future work can still reach a successor through the ecology's
native task carrier after a boundary event, and that the successor can act on
that work through the ecology's required trace.

It is not just "there is still an address." It is:

```text
obligation + carrier + competent successor organization
```

**Carrier-indexed self-object**

The organization that must be preserved or reconstructed for future
addressability and action capacity under a specified carrier ecology.

**Carrier topology**

The set of carriers through which continuity can be preserved, reconstructed,
or lost: source carrier, execution carrier, future-task carrier, tool carrier,
state carrier, observer/probe carrier, and authority/capability carrier.

This topology lets us distinguish:

- **state-preserving continuity**, where memory/execution state is saved or
  migrated
- **organization-reconstructing continuity**, where local state may die but a
  role, mailbox, event log, source, capability, or control surface rebuilds the
  competent organization

**Scaffold gradient**

The spectrum of externally supplied conditions that support the lineage:
hardware, OS, provider, toolchain, source projection, control probe, prompt, and
successor mission. The gradient marks the boundary between bounded
control-surface closure and stronger claims about autonomous agency or
autopoiesis.

**Post-handoff addressability**

Delayed addressability after a real handoff boundary. A successor must receive
new work after the boundary, not merely complete precomputed work.

**Control-surface closure**

The case where a successor remains reachable through a native control surface,
acts on a previously unseen obligation, and restores the waiting state for
future obligations.

**Latent addressability**

The case where an individual is intentionally unreachable through ordinary
routes but remains reachable under a delayed, authorized, or scheduled carrier.
This handles dormancy, air-gapping, rotating addresses, and security-driven
unaddressability without making constant live reachability the criterion.

**Substrate capture**

A failure mode where future addressability is maintained almost entirely by the
infrastructure rather than by the candidate individual or successor
organization. This does not invalidate the continuity relation, but it weakens
claims about the individual's own renewal capacity.

## Core Formal Schema

Let:

- `E` be a future-task ecology
- `C_E` be the carrier through which future obligations arrive in `E`
- `O` be a successor organization
- `B` be a boundary event such as `exec`, migration, restart, or handoff
- `tau_E` be the future-obligation trace required by `E`
- `c` be the specific carrier instance used by `E`

Then:

> `O` is a viable continuation of a computational individual under `E` iff,
> after boundary `B`, `O` remains addressable through `c`, can perform
> `tau_E`, preserves the ecology-relevant behavior, and restores the receptive
> condition under which later obligations can arrive.

For Quine's control-surface case:

```text
idle -> wake(task) -> tool/action -> report -> idle
```

This trace is not a universal identity criterion. It is the criterion for the
Quine `ctl/wake` ecology. Other ecologies have different carrier-indexed
criteria.

Minimal notation for manuscript use:

```text
P --B--> S implies Addressable_E(S, c)
                 and S satisfies tau_E
                 and S is behaviorally adequate relative to c
```

Avoid presenting this as a fully proved theorem until primary sources and a
formal appendix are prepared.

## Theory Neighbors

The neighbor set is now organized into five families plus two load-bearing
anchors. Full prose lives in
[`sections/03-neighboring-theories.md`](sections/03-neighboring-theories.md);
this is the routing summary.

- **self-reference and construction**: von Neumann self-reproducing automata,
  the Kleene recursion theorem, computational reflection (3-Lisp, metaobject
  protocols), and digital-life platforms — self-description as reconstructive
  resource, not as the continuity criterion
- **naming, addressing, and authority**: Saltzer name/address/route, the
  identifier/locator split, and capability sturdy-versus-live references — what
  makes `addressability` mean a bound competent organization, not a resolving
  name
- **concurrency and behavioral identity**: actor model, virtual actors
  (Orleans), futures/promises, process calculi, session types, temporal-logic
  liveness, and coalgebraic bisimulation — carrier-indexed behavioral adequacy
- **systems continuity and reconfiguration**: checkpoint/restart and migration,
  state-machine replication, and membership reconfiguration / leader transfer —
  successor admission as the moment future addressability transfers
- **persistence, individuality, and closure**: closest-continuer theory,
  psychological continuity and fission, interactor/replicator lineages,
  scaffolded reproduction, biological individuality, Markov-blanket
  individuation, autopoiesis, and closure of constraints — continuity criteria
  below full life or agency

Two load-bearing anchors carry the theoretical height:

- **closest-continuer theory** (Nozick) is the metaphysical backbone: the
  criterion fixes the relevant dimension as the future-task carrier and makes it
  testable, which is also what arbitrates the branching / fission problem
  (Parfit) when a lineage forks into multiple successors
- **naming theory** (Saltzer name/address/route) is what keeps `addressability`
  from sliding between "name resolves," "host is up," and "a competent successor
  is bound and ready"; the criterion needs the third

Current posture:

- use these theories to sharpen the criterion
- do not import their strongest claims into Quine by analogy
- treat the corrigibility / safe-interruptibility neighbor as a deliberate
  tension to surface, not a property the criterion endorses

## Novelty Boundary

The novelty is not that computation can have identity across change. Existing
fields already study that through process identity, migration, checkpointing,
actor addresses, endpoint naming, behavioral equivalence, and self-reproducing
automata.

The novelty is the proposed indexing rule:

> when source, state, body, address, and future carrier come apart, the
> continuing computational individual is the organization that preserves future
> addressability under the ecology's native obligation carrier.

This moves the question from:

- "which object was copied?"
- "which process survived?"
- "which state was restored?"
- "which endpoint name stayed stable?"

to:

> through which future carrier can obligations still reach a successor, and
> what organization must exist for that successor to act?

The strongest current claim is therefore:

- **not** a new general theory of life
- **not** a proof of full autopoiesis
- **not** a replacement for process calculi, actors, checkpointing, or
  capability theory
- **yes** a carrier-indexed continuity criterion for process-native and
  distributed computational lineages
- **yes** a way to explain why different ecologies select different reproduced
  self-objects

## Closest Existing Theory Frontier

| Neighbor | Where it already reaches | Where this paper adds |
|----------|--------------------------|------------------------|
| Quines, von Neumann automata, recursion theorem | description copying, constructive reproduction, self-reference | delayed future addressability after a boundary event; reflection alone does not settle continuity |
| Computational reflection (3-Lisp, MOP) | a running system inspecting and rebuilding its own structure | reflection across a destructive process-image boundary, judged by post-handoff carrier viability |
| Naming and binding (name/address/route) | resolving a name to an address to a route | requiring the route to terminate at a bound, competent organization |
| POSIX and OS process semantics | process image replacement, fd inheritance, permissions | individuality as the viable organization under future-task carriers |
| Checkpoint/restart, migration, state-machine replication | preserving state or maintaining a role through redundancy | reconstructing a carrier with no standing quorum; admission as continuity transfer |
| Actor model and virtual actors | mailbox continuity; on-demand reactivation on a new host | mailbox/grain reactivation as the engineered case of latent carrier-indexed identity |
| Process calculi and coalgebraic bisimulation | observable traces and behavioral equivalence over an interface | equivalence indexed by the carrier the ecology exposes, not a fixed interface |
| Capability systems (sturdy vs live refs) | authority-bearing references and unforgeable access | a sturdy reference as latent addressability; authority as one carrier |
| Closest-continuer and psychological continuity | continuant as closest continuer; survival and fission | fixing the relevant dimension as the future-task carrier, making it testable |
| Interactor/replicator, scaffolded reproduction, biological individuality | lineage individuality and reproduction on a scaffold gradient | computational interactor as carrier-indexed organization; scaffold gradient as the closure boundary |
| Autopoiesis, closure of constraints, Markov-blanket individuation | self-production, organizational and statistical individuation | a weaker, scaffolded, carrier-defined criterion below full agency |
| Corrigibility / safe interruptibility | preserving an external control and shutdown channel | the dual tension: continuity preserves the future-task interface, not the control channel |

The nearest theories reach the pieces: address, state, behavior, authority,
description, construction, reconfiguration, persistence, and organizational
closure. The paper's boundary is their composition around future obligations:

> computational individuality is indexed by future addressability.

## Objections To Carry Forward

The strongest objections are useful and should stay in the paper:

- **Biological autonomy objection**: addressability is not metabolic
  self-maintenance; it is weaker than autopoiesis and intrinsic normativity.
- **Mechanistic objection**: "obligation" can sound semantic or observer
  imposed; the manuscript must ground obligations in concrete protocol
  expectations, queues, probes, or control events.
- **Triviality objection**: if every reachable endpoint is an individual, the
  criterion is too promiscuous; we need boundary events, carrier specification,
  and future-obligation traces.
- **Dormancy objection**: constant reachability is too strong; latent
  addressability is the better criterion for air-gapped, sleeping, or rotating
  systems.
- **Substrate-capture objection**: if the infrastructure does all preservation,
  the claim shifts from individual renewal to ecology-maintained continuity.

## Current Evidence Hook

Quine supplies the strongest motivating contrast:

- `6.63`: self-source available, raw-stream pressure, thin custom successor
- `6.64`: self-source available, `ctl/wake` pressure, rebuilt runtime body,
  post-`exec` wake, tool action, return to `idle`
- `6.65`: same control-surface behavior under minimized initial prompt
- `6.66`: generic POSIX `exec` utility baseline
- `6.67`: direct source rebuild baseline without delayed addressability
- `6M.10`: carrier/capability-overlap ladder separating source/body material,
  carrier availability, prompt-explained implementation repair, and
  post-handoff proof / readiness

Read:

- source visibility is not sufficient
- generic `exec` is not sufficient
- source rebuild is not sufficient
- formal reconstruction is not operational continuity
- the future-task carrier/capability determines the necessary successor
  organization

The prospective experiment home for turning this hook into discriminating
evidence is:

- [`experiments/active/autopoiesis-succession/README.md`](../../../experiments/active/autopoiesis-succession/README.md)

## Evaluation Implications

Future evidence should report:

- boundary event actually crossed
- carrier through which the later obligation arrived
- whether the later obligation was unseen before the boundary
- whether the successor satisfied the required future-obligation trace
- whether the successor restored the receptive condition
- which scaffold supplied source, state, authority, tools, provider, and probe
- failure mode, especially state drift, substrate capture, irreversible
  degradation, or unbounded recursion

## Experiment Feeds

- `experiments/active/autopoiesis-succession/succ-21-minimal-control-adapter/README.md` - primary - exploratory-only - decisive next test for whether the carrier-indexed self-object is Quine runtime identity or abstract control-surface compatibility
- `experiments/active/autopoiesis-succession/succ-22-successor-mission-ablation/README.md` - primary - exploratory-only - maps the scaffold gradient between delayed addressability and more autonomous carrier maintenance
- `experiments/active/autopoiesis-succession/succ-23-multi-wake-closure/README.md` - primary - exploratory-only - distinguishes one delayed addressability event from restored future-addressability capacity
- `experiments/active/autopoiesis-succession/succ-24-carrier-perturbation/README.md` - primary - exploratory-only - tests whether future addressability can become an actively maintained carrier relation
- `experiments/active/autopoiesis-succession/succ-25-checkpoint-restart-contrast/README.md` - primary - exploratory-only - supplies the clean contrast between state-preserving and organization-reconstructing continuity
- `experiments/active/autopoiesis-succession/succ-26-source-degradation-repair/README.md` - primary - exploratory-only - tests whether self-description participates in active reconstruction rather than passive source availability
- `experiments/active/autopoiesis-succession/succ-27-actor-mailbox-analog/README.md` - primary - exploratory-only - tests whether carrier-indexed individuality generalizes from Quine control surfaces to actor mailbox systems
- `experiments/active/autopoiesis-succession/succ-28-multi-generation-recursion/README.md` - primary - exploratory-only - tests whether future addressability can be a lineage-level relation rather than a one-step continuity criterion
- `experiments/active/autopoiesis-succession/succ-38-carrier-capability-overlap-ladder/README.md` - secondary - supporting-only - carrier-continuity contrast where source/body material and future-control capability separate

## Current Evidence Backbone

- `experiments/active/autopoiesis-succession/succ-21-minimal-control-adapter/README.md` - `6F.01` - exploratory-only - runnable comparator currently records methodologically retained predecessor failures, not a new positive claim
- `experiments/active/autopoiesis-succession/succ-22-successor-mission-ablation/README.md` - `6F.02` - exploratory-only - design-only scaffold-gradient ladder for weakening the successor mission after the `6.65` positive
- `experiments/active/autopoiesis-succession/succ-23-multi-wake-closure/README.md` - `6F.03` - exploratory-only - design-only repeated-wake test for distinguishing one-shot addressability from restored receptive capacity
- `experiments/active/autopoiesis-succession/succ-24-carrier-perturbation/README.md` - `6F.04` - exploratory-only - design-only perturbation test for passive carrier dependence versus active carrier maintenance
- `experiments/active/autopoiesis-succession/succ-25-checkpoint-restart-contrast/README.md` - `6F.05` - exploratory-only - design-only contrast between state-preserving and organization-reconstructing continuity
- `experiments/active/autopoiesis-succession/succ-26-source-degradation-repair/README.md` - `6F.06` - exploratory-only - design-only test for active self-description repair under post-handoff carrier pressure
- `experiments/active/autopoiesis-succession/succ-27-actor-mailbox-analog/README.md` - `evidence-set:6F.07-first-mailbox-analog-wave-20260628` - exploratory-only - deterministic mailbox analog confirms delayed-addressability can be decided on a non-`ctl/wake` future-task carrier; preserved-mailbox arms pass while address-lost controls fail
- `experiments/active/autopoiesis-succession/succ-28-multi-generation-recursion/README.md` - `6F.08` - exploratory-only - design-only lineage-recursion test for future addressability across multiple successor generations
- `experiments/active/autopoiesis-succession/succ-38-carrier-capability-overlap-ladder/README.md` - `6M.10` - supporting-only - retained carrier/capability contrast showing why source/body reconstruction and delayed future-control continuity must be coded separately

## Unsafe Overclaims

Do not claim:

- future addressability is the only possible individuality criterion
- every computational system has a self
- every service endpoint is an organism
- actor identity, PID identity, source identity, or address identity alone is
  sufficient for individuality
- Quine has full autopoiesis, full semantic closure, or full agency
- prospective `6F.xx` experiment designs already constitute positive evidence
  for the theory

It claims a conditional theory:

> where a computational lineage is judged by future obligations, individuality
> is indexed by the carrier through which those obligations can still arrive
> and be acted on.

## Blocking Gaps

- primary sources from the neighboring literatures still need to be checked
  before the theory can become citation-ready
- the `6F.xx` evidence path is mostly prospective or exploratory-only, not yet
  a discriminating positive backbone
- the venue and publication route are still unresolved
- the formal schema must stay lightweight until the manuscript can justify
  stronger notation without pretending to be a full proof system
- the paper must keep the computational-individuality theory distinct from
  Quine-specific evidence retelling, full autopoiesis, and generic AI
  self-replication claims

## Next Useful Work

- source-check the 2026-06 neighbor backbone against primary sources before
  treating the new citations (Nozick, Parfit, Hull, Godfrey-Smith, Clarke,
  Pradeu, Saltzer, Kleene, Smith, Maes, Bykov, Schneider, Ongaro, Kirchhoff,
  Soares, Orseau) as citation-ready rather than draft keys
- decide which of the five neighbor families stay in the main text versus a
  related-work appendix; the two load-bearing anchors (closest-continuer,
  naming theory) should stay in the main argument
- decide whether this paper should be an `Artificial Life` theory paper, a
  `HotOS` / systems-theory note, or an interdisciplinary computational
  individuality essay
- decide whether the corrigibility tension warrants its own short safety
  paragraph or stays folded into the limits section

## Local Surfaces

| Path | Role |
|------|------|
| [`WORKSPACE.md`](WORKSPACE.md) | local working-surface guide |
| [`manuscript-plan.md`](manuscript-plan.md) | initial manuscript plan and section skeleton |
| [`sections/README.md`](sections/README.md) | first sectionized manuscript draft |
| [`plans/experimental-roadmap-2026-04.md`](plans/experimental-roadmap-2026-04.md) | planned discriminating experiments |
| [`plans/README.md`](plans/README.md) | active paper-local planning prompts |
| [`references/README.md`](references/README.md) | recovered research reports and citation-discovery inputs |
| [`../../../experiments/active/autopoiesis-succession/README.md`](../../../experiments/active/autopoiesis-succession/README.md) | prospective `6F.xx` experiment control plane |
