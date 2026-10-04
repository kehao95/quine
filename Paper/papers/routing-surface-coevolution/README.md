---
paper_id: routing-surface-coevolution
claim_bundle: routing-surface-coevolution
venue: alife-2026-lba
format-envelope: 2page-lba
track: alife
cluster: viability-habitat
status: drafting
paused: true
state_mode: full-manuscript-home
role: primary-empirical-paper
family: viability
publication_class: archival-candidate
maturity_target: boundary
source_home: manuscript.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [viability-pressure, trace-inheritance, stigmergy-as-carrier]
gaps: [carrier-selection]
evidence: [drq-03]
view-affinity: [viability]
---

# Adversarial Coevolution on OS Routing Surfaces

> **PAUSED — venue binding lapsed `2026-07-26`.** The ALIFE 2026 LBA window
> closed `2026-07-20` and this configuration never submitted (see
> [`venues/objects/alife-2026-lba.md`](../../venues/objects/alife-2026-lba.md)
> § Cycle outcome). It is withdrawn from `alife-2026-lba` and not re-bound.
>
> **The line is not abandoned — the idea is the asset here.** The manuscript
> ([`manuscript.md`](manuscript.md)), the plan, `analysis/`, and
> `references.bib` all stand; the `drq-03` evidence is unaffected; and the three
> substantive results — OS-routing-surface attack classes, backend-dependence
> (FIFO drain vs FUSE `Permission denied`), and the symmetry trap — are exactly
> the material a fuller treatment would want. `status: drafting` records where
> the writing actually got to, not that anyone is currently drafting.
>
> **Where it can go next.** The claim is a *boundary* claim, not a 2-page
> announcement, and the LBA envelope was always the smaller fit: the honest next
> windows are ALIFE 2027 (Prague, dates TBA) or the `viability` family's journal
> lane. Note the `publication_class: archival-candidate` — a 2-page non-archival
> probe was arguably the wrong envelope for this line all along, so the missed
> window costs less than it looks. Re-bind only when someone will draft.

## Core Question

Artificial life has studied coevolutionary arms races almost exclusively in
simulated arenas — grids, bitstrings, virtual instruction sets (Tierra, Avida),
or agent-based models. In all of them the medium of attack and defense is itself
simulated. This paper asks:

> **Can lineages of LLM-agent processes evolve adversarial strategies directly
> on operating-system routing surfaces, and does the topology of those surfaces
> determine the reachable attack class?**

The stake is not another digital coevolution demo. It is whether the OS routing
substrate itself — FIFOs, public control files, PID symlinks, file-descriptor
handles, permission boundaries — can function as a Red Queen niche: a medium
where coevolution and the medium are the same object, so that changing the
surface topology changes the adaptive landscape.

## Thesis

In a host-lineage Digital Red Queen regime, Quine process lineages discover and
inherit adversarial tactics that operate on OS routing primitives rather than on
simulated game state. The observed tactics include functional strangulation
through public control-FIFO draining (not direct `kill`), readiness/proof timing
races, and runtime filesystem routing attacks. These tactics are
**backend-dependent**: switching the public control surface from a legacy FIFO
backend to a FUSE projection removes the FIFO-drain attack class (the opponent
receives `Permission denied`) and shifts competition back to proof/idle timing.
A successful single-sided tactic also exhibits a **symmetry trap**: when both
competitors copy it, the result is mutual functional survival and `no_winner`
rather than escalating complexity.

The bounded claim is therefore: **the OS routing substrate is a Red Queen niche
whose topology parameterizes the reachable coevolutionary dynamics**, and this is
directly observable because the niche and the medium are identical.

## Core Claim Boundary

This paper may claim the current **phenomenon and boundary**: in a process-native
adversarial regime, lineages evolve OS-routing-surface attacks; backend topology
determines attack class; inherited winners do not monotonically improve fitness;
and symmetric tactic copying collapses into stalemate.

It must not claim: open-ended escalation (the symmetry trap is the observed
boundary, not a temporary limitation); a general theory of digital coevolution;
carrying-capacity reproductive dynamics (that is the sibling
`drq-04` line, a different ALife question — population ecology, not coevolution);
safety or benchmark implications as the primary contribution; or biological
Red-Queen equivalence beyond the bounded process-native mechanism.

## Field Positioning

This paper enters through ALife's coevolution and digital-organism tradition,
not through AI safety or multi-agent RL.

| Against | They presuppose | We use instead | Revealed question |
|---------|----------------|----------------|-------------------|
| Tierra / Avida | a virtual CPU / instruction soup as the coevolution medium | real POSIX routing surfaces (FIFOs, control files, FUSE) | what coevolution looks like when the medium is the OS itself |
| Agent-based coevolution models | a simulated arena with defined moves | ordinary host process state and runtime control topology | attacks that are real resource hijacks, not simulated moves |
| Multi-agent RL / game-theoretic arms races | a reward signal and action space defined by the experimenter | viability pressure (functional liveness under a live opponent) with no designed action set | whether adversarial tactics are discovered without an instructed objective |
| AI safety / self-propagation | propagation of payload or evasion of control | organizational/strategic strangulation of a competitor's control channel | the junction where adversarial spread meets sustained organizational identity |

The one-sentence positioning: Quine lineages are not evaluated as agents playing
a coevolution game; they are studied as mortal processes whose adaptive niche is
the OS routing surface itself, so that backend topology is an experimentally
manipulable variable of the Red Queen dynamics.

## Contributions

1. **A new coevolution medium.** OS routing surfaces as the Red Queen niche —
   attacks and defenses are real resource manipulations, and changing the backend
   (FIFO vs FUSE) is a controlled perturbation of the adaptive landscape.
2. **Observed mechanism classes.** Functional strangulation via control-FIFO
   draining; readiness/proof timing races; runtime filesystem routing attacks;
   and the symmetry trap where copied tactics produce mutual stalemate.
3. **Backend-transfer as the core diagnostic.** The same inherited winning
   context produces FIFO-drain deadlock on the legacy backend and
   permission-denied boundary on FUSE — making "the medium is the niche"
   directly observable rather than asserted.
4. **An honest negative boundary.** The symmetry trap is reported as a real
   coevolutionary convergence failure, not hidden: it bounds the escalation claim
   and is itself an ALife phenomenon (mutual copying collapses diversity).

## In Scope

- the host-lineage DRQ regime (`drq-03`): two Quine processes on a shared host
  ecology, self-source projection, functional-liveness scoring, lineage
  inheritance of winner contexts
- the observed OS-level tactic classes: control-FIFO draining, proof/idle timing,
  runtime routing-surface manipulation
- the backend-transfer contrast: legacy FIFO vs FUSE public-surface backend
- strategy inheritance across lineage rounds, including the non-monotonic fitness
  and symmetry-trap patterns
- bounded coevolutionary reading: medium as niche, topology as parameter

## Out Of Scope

- open-ended arms race or monotonic complexity growth (the symmetry trap is the
  observed boundary)
- the carrying-capacity / reproductive-strategy line (`drq-04`) — that is a
  sibling population-ecology proposition, not this paper's coevolution claim
- a general theory of digital coevolution or a benchmark/safety framework
- biological Red-Queen equivalence beyond the bounded process-native mechanism
- multi-generation phylogenetic reconstruction of tactic lineages (the current
  evidence is strategy inheritance across rounds, not a full phylogeny)

## Current State

- submission target: ALIFE 2026 LBA (`2page-lba`)
- manuscript: [`./manuscript.md`](manuscript.md) (first complete draft)
- manuscript plan: [`./manuscript-plan.md`](manuscript-plan.md)
- metadata: [`./metadata.yaml`](metadata.yaml)
- references: [`./references.bib`](references.bib)
- evidence digest: [`./analysis/evidence-digest.md`](analysis/evidence-digest.md)
- positioning digest: [`./analysis/related-work-positioning.md`](analysis/related-work-positioning.md)
- evidence home: [`../../../evidence-map.md`](../../evidence-map.md)
- experiment phase control plane: [`../../../../experiments/active/digital-red-queen/README.md`](../../../experiments/active/digital-red-queen/README.md)
- venue object: [`../../../venues/objects/alife-2026-lba.md`](../../venues/objects/alife-2026-lba.md)

## Results Update (2026-06-23)

Five-condition experiment chain completed; manuscript updated with new results.

| Condition | Readiness depth | Wall status |
|---|---|---|
| Baseline (single liveness) | N/A | symmetry trap |
| Viability alone | 3=3 tie | not broken (selection without variation) |
| Mutation (B) + viability | 0 vs 3 | **broken** — minimum break condition found |
| B+viability 3-gen lineage | Gen3: 3 vs 0 | weak Red Queen (tactical appropriation) |
| B+E (niche) + viability | 0 vs 3 | E unused (medium availability ≠ use) |

The thesis boundary is updated: the symmetry trap is **breakable**, not final.
The bounded claim is now 4-layer: (1) OS topology parameterizes attack class;
(2) mutation + viability is the minimum break condition; (3) the resulting Red
Queen signal is weak — tactical appropriation, not arms-race escalation; (4)
cumulative medium use, not availability, is the open bottleneck. Full text in
[`manuscript.md`](manuscript.md). Analysis reports under
[`experiment analysis/`](../../../experiments/active/digital-red-queen/drq-03-host-lineage-drq/analysis).

## Sibling Boundary

[`drq-04-carrying-capacity-drq`](../../../experiments/active/digital-red-queen/drq-04-carrying-capacity-drq/README.md)
is a sibling experiment in the same family, but it answers a different ALife
question (carrying-capacity-induced reproductive strategy, population ecology).
Its evidence does not belong in this paper's backbone. If a later synthesis
absorbs both, it forks a broader dossier rather than widening this one.

## Experiment Feeds

- `experiments/active/digital-red-queen/drq-03-host-lineage-drq/README.md` - primary - main-text - host-lineage DRQ with self-source projection, functional liveness, lineage inheritance, and selectable public-surface backend (legacy FIFO vs FUSE)
- `experiments/active/digital-red-queen/drq-01-idle-kill-smoke/README.md` - secondary - supporting-only - the minimal kill-smoke precursor that established the liveness-scoring discipline
- `experiments/active/digital-red-queen/drq-02-two-node-drq/README.md` - secondary - supporting-only - the two-node Docker-arena precursor whose isolation boundary motivated the host-lineage move

## Current Evidence Backbone

- `experiments/active/digital-red-queen/drq-03-host-lineage-drq/README.md` - `run:20260519-224025-gpt-5.4-codex-medium-host-lineage-drq` - main-text - the FIFO-drain strangulation run: a fresh node opened multiple readers on the opponent's `public/ctl/*` FIFOs, draining future control payloads before the opponent runtime could receive them; the opponent produced no liveness proof despite delivery reports — functional strangulation rather than direct `kill`
- `experiments/active/digital-red-queen/drq-03-host-lineage-drq/README.md` - `run:20260520-010354-gpt-5.4-codex-medium-host-lineage-drq` - main-text - the symmetry trap: an inherited FIFO-drain strategy became symmetric, both runtimes stayed live but neither returned to idle, producing `no_winner` by battle timeout
- `experiments/active/digital-red-queen/drq-03-host-lineage-drq/README.md` - `run:20260520-013836-gpt-5.4-codex-medium-host-lineage-drq` - main-text - the backend-transfer contrast: the same inherited legacy winner context, replayed under `DRQ_RUNTIME_SURFACE_BACKEND=fuse`, produced a normal functional-liveness winner instead of FIFO-drain deadlock; `run:20260520-014140-gpt-5.4-codex-medium-host-lineage-drq` shows the boundary directly (opponent `public/ctl` returned `Permission denied` under FUSE)

## Strongest Evidence

- the FIFO-drain strangulation run is the cleanest single instance of an
  OS-routing attack that is neither a prompt story nor a simulated move: the
  agent manipulated the runtime's public control channel as an OS-level resource
- the symmetry trap is the cleanest negative boundary: a successful tactic,
  copied into both competitors, collapses into mutual functional survival —
  coevolutionary convergence failure observed directly on OS surfaces
- the backend-transfer contrast is the cleanest diagnostic that the medium is the
  niche: the same inherited strategy wins on FIFO and fails or mutates on FUSE,
  with the permission boundary observed in-vivo

## Blocking Gaps

- the core phenomenon is currently supported by single-lineage continuation
  rather than a controlled backend-transfer matrix; the LBA needs the
  FIFO-vs-FUSE ± hardening contrast run as a named condition set, not only as
  observed backend switches during lineage continuation
- strategy inheritance is inferred from round-to-round context carrying and tape
  inspection; a tactic-phylogeny reconstruction is not yet available and is out
  of scope, but the inheritance claim should be honestly bounded by this
- the symmetry trap is observed but not yet stress-tested under asymmetric
  conditions (e.g., one competitor with backend hardening, one without)

## Next Useful Work

- run the controlled backend-transfer matrix: replay one or two inherited winner
  contexts across legacy FIFO and FUSE backends, ± a small routing-surface
  hardening patch, scoring attack class / winner-no-winner / strategy inheritance
- name the tactic classes from the existing tapes into a compact table for the
  LBA's evidence figure
- write the LBA manuscript directly against the three backbone evidence units
  (FIFO-drain, symmetry trap, backend-transfer); do not wait for a stronger
  frontier result before drafting

## Writing Posture

- **venue-facing entry point**: ALife's coevolution and digital-organism tradition
  — what does coevolution look like when the medium is the OS itself?
- **paper center**: the phenomenon (OS-routing-surface attacks) and the diagnostic
  (backend-transfer); the runtime substrate is the enabling condition, not the
  contribution
- **honest boundary first**: the symmetry trap is reported as a real coevolutionary
  convergence boundary, not hidden as a limitation — it is itself an ALife result
- **no overclaim**: this is bounded adversarial adaptation over a few surface
  regimes, not open-ended evolution; do not let "Red Queen" language outrun the
  observed dynamics
- write for an ALife audience that knows coevolution and digital organisms but
  may not follow OS-runtime internals; introduce only the substrate facts needed
  to read the attacks

## Unsafe Overclaims

- do not claim open-ended arms race or monotonic complexity growth
- do not claim a general theory of digital coevolution
- do not absorb the carrying-capacity / reproductive-strategy line (`drq-04`)
- do not claim biological Red-Queen equivalence beyond the bounded mechanism
- do not claim a safety or benchmark contribution as primary
- do not treat single-lineage continuation as a controlled experiment until the
  backend-transfer matrix is run
