---
paper_id: lifecycle-bottleneck-reproduction
claim_bundle: lifecycle-bottleneck-reproduction
venue: artificial-life-journal
format-envelope: full-journal
track: theory
cluster: autopoiesis-renewal
status: planning
state_mode: planning-skeleton
role: lifecycle-bottleneck-theory-paper
family: successor-continuity
publication_class: theory-candidate
maturity_target: theory-anchor
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theme: autopoiesis-renewal
theme_home: ../../../theory/views/continuity.md
theory-objects: [carrier, exec-barrier, delayed-addressability, trace-inheritance, viability-pressure]
gaps: [branching-identity, collective-continuity, carrier-selection]
evidence: [6M.04, 6M.05, 6M.06, 6M.10]
view-affinity: [continuity]
---

# Self-Reproduction as Lifecycle-Bottleneck Successor Formation

## Current Question

Is self-reproduction better understood as organizational continuity across a
lifecycle bottleneck than as primitive copying of a body, source file, or
artifact?

## Thesis

Self-reproduction is one successor morphology that appears when an organization
must preserve continuity across a lifecycle bottleneck using available
self-description material under future-carrier demand and viability-depth
pressure.

Short form:

```text
lifecycle bottleneck
  + self-description material/degradation
  + future-carrier/capability demand
  + viability-depth requirement
  -> successor morphology
```

Recent theory absorption clarifies the material-overlap analogy. For Quine,
the safer computational term is **carrier / capability overlap**: source and
runtime contract are reconstructive or normative self-description, while fd,
`ctl/wake`, mailbox, endpoint, runtime registration, retained filesystem
surface, or other authority-bearing carriers are what let future obligations
reach the successor after the bottleneck.

## Core Claim Boundary

This paper may claim a **theory of reproduction as bottleneck-crossing
successor formation**. It may use Quine as an operational case showing how
contract-level successors, stream/fd adapters, preserved-body reentry, rebuilt
runtime successors, compensating responders, and readiness-depth frontiers
separate under different resources, carriers, and viability-depth requirements.

It must not claim that current Quine evidence proves a universal biological
law, endogenous reproductive drive, open-ended evolution, population-level
Darwinian selection, or full autopoiesis.

## Theory Anchors

This paper is the reproduction-ontology slice, so its theoretical height should
come from the canonical lineage-individuality literature rather than from
retelling the Quine ladder. Four anchors carry it:

- **Interactor / replicator (Hull).** The reproduced object is best read as an
  *interactor* — the organization that engages its environment and carries
  lineage continuity — rather than as a replicator-style copied description.
  Quine's `source-code/` is replicator-like reconstructive material; the
  carrier-bound runtime organization that future obligations actually reach is
  the interactor. This is the cleanest home for "organization identity versus
  source identity," and it keeps lineage continuity, not byte copying, as the
  unit.
- **Scaffolded reproduction (Godfrey-Smith).** Reproduction lies on a gradient
  from simple to scaffolded to collective, where external machinery supplies
  part of the reproductive work. Quine is openly scaffolded: provider, OS,
  toolchain, source projection, control probe, and successor mission do real
  reproductive work. The scaffold gradient is therefore the right axis for
  placing bounded control-surface closure below endogenous reproduction, and it
  replaces hand-waving about "how much is really the entity" with a named
  position on a known spectrum.
- **Biological individuality (Pradeu, Clarke).** The debate over physiological
  versus evolutionary individuality, and over markers and mechanisms that make
  one persisting individual, is the right contrast class for asking what makes a
  computational lineage one continuing individual across a bottleneck. The
  borrowed move is the question form, not biological equivalence.
- **Lifecycle bottleneck (carried through).** Reproduction as bottleneck-
  crossing successor formation is the paper's own frame; the anchors above give
  it a vocabulary that situates it among reproducer theory rather than inventing
  terms.

Discipline: borrow the question forms and the gradient, not the strong
biological claims. The interactor framing must not become a claim of metabolic
individuality, and the scaffold gradient must not be read as a path that Quine
has already traversed toward endogenous reproduction.

## In Scope

- reproduction as organizational continuity across lifecycle bottlenecks
- the reproduced object as interactor rather than replicator: the
  lineage-carrying runtime organization, not the copied self-description
- reproduction on a scaffold gradient, with Quine placed as openly scaffolded
  rather than endogenous
- repair / regeneration / reproduction as a continuum of threatened
  organization
- runtime contract as formal reconstitutor, executable scaffold, or inherited
  niche rather than hidden genome
- source-code as reconstructive resource rather than reproduction criterion
- carrier / capability overlap as the operational counterpart to biological
  material-overlap intuitions
- externally imposed future viability pressure as the current Quine trigger
- Quine evidence as bounded operational examples, not the whole empirical
  burden

## Out Of Scope

- retelling the full boundary-paper evidence ladder
- proving bounded control-surface closure as a new empirical result
- systems proof of `exec`, fd, or `ctl/wake` mechanics
- evaluation-framework or benchmark design
- claiming biological equivalence for Quine's lifecycle

## Relation To Siblings

- [`../boundary-paper/README.md`](../boundary-paper/README.md) owns the current
  empirical phenomenon and bounded control-surface closure.
- [`../future-addressability/README.md`](../future-addressability/README.md)
  owns the narrower carrier-indexed computational individuality criterion, and
  owns the canonical general statement of the post-handoff-viability criterion;
  this paper's `post-bottleneck viability` is a reproduction-ontology
  specialization that cites it rather than redefining it (see
  [`../../../theory/views/continuity.md`](../../../theory/views/continuity.md)).
- [`../autoconstructive-successor-morphology/README.md`](../autoconstructive-successor-morphology/README.md)
  owns ALife comparator positioning.
- [`../generative-self-replacement/README.md`](../generative-self-replacement/README.md)
  owns systems mechanism.

## Current State

- manuscript workspace: [`WORKSPACE.md`](WORKSPACE.md)
- current manuscript plan: [`manuscript-plan.md`](manuscript-plan.md)
- theme theory input: [`../../../theory/views/continuity.md`](../../../theory/views/continuity.md)
- split note: [`../../../theory/views/continuity/plans/theory-publication-split-2026-05.md`](../../../theory/views/continuity/plans/theory-publication-split-2026-05.md)
- current posture: theory-candidate seed; do not draft as if citation/source
  checking were complete

## Current Read

The paper should raise the theory above the current strongest result without
letting that result swallow the theory. Bounded control-surface closure is a
case of successor morphology under one resource-pressure pairing, not the whole
account of self-reproduction.

The substrate matters because it determines how much of reproduction is already
done for the entity. In Quine, POSIX supplies ordinary execution physics
instead of offspring semantics: no privileged reproduction API, no substrate
offspring slot, no automatic continuity judgment, and no rollback if a
destructive `exec` handoff fails. This makes the lifecycle bottleneck unusually
visible. The lineage must construct the successor body and restore the
future-work carrier, while the experiment judges viability after handoff rather
than accepting body production as success.

Objective post-bottleneck viability should be read as:

```text
replacement occurred
  + declared carrier reaches successor
  + successor completes delayed previously unseen work
  + successor recovers ready / idle state
```

This keeps source rebuild in the right place: it can be a powerful way to
construct a successor body, but it is not itself reproduction-like continuity
unless the load-bearing future carrier/capability is preserved or rebuilt and
the successor remains competent after the bottleneck.

## Writing Posture

Use ALife / philosophy-of-biology language carefully. The safest rhetoric is:

- copying is not the primitive problem;
- crossing a lifecycle bottleneck while preserving organization is the problem;
- self-copy is a special morphology that can solve that problem under some
  resources and pressures.

## Experiment Feeds

- `none-yet` - none - not-for-paper-yet - current Quine experiments motivate the theory, but this seed dossier has not promoted a stable paper-local evidence feed

## Current Evidence Backbone

- `none-yet` - `none-yet` - not-for-paper-yet - no stable evidence backbone is currently declared for this theory paper

## Strongest Evidence

- source-absent cleanroom continuation suggests organization can recur without
  source or body copying
- source-visible raw-stream successor formation suggests source is resource,
  not destiny
- source-visible native-control closure suggests full runtime rebuild appears
  when pressure and material align
- prompt-minimized native-control closure suggests the result is not only a
  fully narrated initial strategy, while successor-mission scaffolding remains
- `6M.04`--`6M.06` show that pressure depth and material degradation matter:
  custom responders, addressable-but-incompetent basins, and multi-wake
  readiness frontiers are legitimate theory data rather than noise
- `6M.10` sharpens the material-overlap analogy: intact or repaired source can
  coexist with carrier-withheld or control-proof failure, while a
  prompt-explained missing `exec` implementation can be repaired to full
  closure when the load-bearing future-control carrier works

## Blocking Gaps

- primary-source citation cleaning for lifecycle bottlenecks, reproduction
  theory, formal/material reproduction, developmental systems, niche
  construction, and the new lineage-individuality anchors (Hull
  interactor/replicator, Godfrey-Smith scaffolded reproduction, Pradeu/Clarke
  biological individuality)
- deciding whether the theory is paper-sized before or after the boundary paper
  receives field feedback
- sharper distinction between Quine's externally imposed pressures and
  endogenous biological reproductive pressures

## Next Useful Work

- build the manuscript spine around the interactor/replicator split and the
  scaffold gradient, so the paper situates Quine within reproducer theory rather
  than inventing terms
- prepare a citation-clean theory memo from primary sources before drafting
  public manuscript claims
- decide whether this paper should be theory-track, ALife-theory workshop, or
  later journal/synthesis route, and whether it merges with future-addressability
  as one theory paper

## Unsafe Overclaims

- Quine has biological lifecycle, metabolism, germline, or endogenous
  reproductive drive
- source-absent contract-level regeneration is genetic inheritance
- pressure-conditioned successor morphology proves general biology
- bounded control-surface closure equals full autopoiesis
