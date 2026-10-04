# Manuscript Plan: Computational Individuality Is Indexed By Future Addressability

This file owns the initial manuscript plan for the future-addressability theory
paper.

## Working Title

Current default:

- `Computational Individuality Is Indexed by Future Addressability`

Alternate titles:

- `Future Addressability and Computational Individuality`
- `When the Future Selects the Self: Carrier-Indexed Individuality in Computational Systems`
- `Post-Handoff Addressability as a Criterion of Computational Continuity`

## One-Sentence Thesis

In computational systems that continue through handoff, replacement, or
reconstruction, the relevant individual is not fixed by source, state, name, or
artifact identity; it is indexed by the future carrier through which obligations
can still reach and be acted on by the successor.

## Target Contribution

The paper should contribute a theory of computational individuality that is:

- general enough to speak beyond LLM agents
- precise enough to distinguish carriers and continuity criteria
- grounded by Quine's `6.63 -> 6.67` control-surface contrast and the
  `6M.10` carrier/capability-overlap ladder
- bounded enough not to claim full life, full agency, or full autopoiesis

In the series stack, this is the **generalization** paper. The LBA announces
the result, the main boundary paper carries the Quine evidence ladder, and this
paper abstracts the continuity criterion to obligation-bearing computational
lineages.

The claim is not that future addressability explains all identity. The claim is
that, for computational lineages judged by future obligations, future
addressability is the indexing relation that determines which organization must
count as the continuing self-object.

## 2026-05-30 Theory-Height Push

The next framing pass should raise the paper from a Quine continuity criterion
to a general computational-individuality theory candidate.

Core move:

```text
identity handles diverge -> future carrier selects the criterion ->
successor organization satisfies the trace -> receptive state is restored
```

This should make the paper prospective rather than retrospective. Existing
neighbor theories often ask what description, state, address, behavior,
authority, or organization persists from the past. This paper asks what future
obligation can still bind the lineage after a boundary event, through which
carrier, and what organization must exist for that obligation to be acted on.

That gives the paper its theory height: future addressability is not only a
reachability property, but a selection criterion over candidate self-objects.
The future-task carrier determines which organization must be preserved or
reconstructed for the lineage to remain obligation-bearing.

Immediate manuscript consequence:

- introduce future addressability as prospective selection over identity
  handles, not only as delayed reachability
- make `obligation + carrier + competent successor organization +
  receptive-state recovery` the main criterion
- move the cross-media carrier comparison into the central argument
- bridge to successor morphology with the compressed relation:

```text
future carriers select individuality criteria;
resource and viability pressures shape the morphology that can satisfy them.
```

Paper-local planning note:
[`plans/theory-height-push-2026-05-30.md`](plans/theory-height-push-2026-05-30.md).

## Novelty Boundary

The paper's novelty should be stated as an indexing criterion, not as a claim
to have discovered computational identity in general.

Already known:

- computational entities can preserve process, state, endpoint, address, or
  behavior across change
- quines and self-reproducing automata separate description from construction
- actors, process calculi, services, and capability systems already treat
  addressability and observable interaction as central
- autopoiesis, semantic closure, and closure of constraints already theorize
  organizational self-production at stronger biological or semiotic levels

New claim:

> when these identity handles diverge, the continuing computational individual
> is indexed by the future-task carrier through which obligations can still
> reach and be acted on by a successor.

This makes future addressability a continuity criterion:

- source identity can be present without continuity
- source rebuild can be present without delayed addressability
- endpoint or PID continuity can be present without competent action
- competent action can be carrier-indexed rather than body-indexed
- the relevant self-object changes when the future carrier changes

The claim is strongest where the system has a real boundary event and a delayed
obligation:

```text
boundary crossed -> future carrier delivers new obligation -> successor acts
through the required interface -> successor restores receptive condition
```

For Quine, the receptive condition is `idle`; for other systems it may be a
mailbox, endpoint, scheduler slot, file descriptor, capability, or restore
protocol.

## Closest Theory Frontier

The manuscript should use a "near but not identical" structure for related
work:

| Neighbor | What it gives us | Why it is not the same claim |
|----------|------------------|------------------------------|
| Quines | self-description and textual fixed points | no delayed future-obligation test |
| von Neumann reproduction | description plus construction | offspring semantics are usually specified by the model |
| POSIX process semantics | concrete `exec`, fd, permission, and process-image boundaries | gives process physics, not individuality criteria |
| Process migration | state continuity across hosts | often preserves state rather than reconstructing viable future addressability |
| Checkpoint/restart | saved state plus restore protocol | can restore state without testing future-obligation competence |
| Actor model | mailbox address and message handling | one carrier case, not a general theory of carrier-indexed individuality |
| Process calculi | traces, bisimulation, observable behavior | formalizes interface equivalence but does not ask which self-object future ecology selects |
| Capability systems | authority-bearing carriers | authority is necessary in some ecologies but not sufficient for individuality |
| Semantic closure | symbol/action relation | Quine has bounded recoupling, not full semantic closure |
| Autopoiesis and agency | upper bound for self-produced individuality and normativity | current Quine evidence is scaffolded and below intrinsic agency |

The phrase to protect:

> The nearest theories each own one component of the problem; this paper
> composes them around future obligations as the indexing relation for
> computational individuality.

The paper must make the distinction from neighboring vocabularies explicit:
future addressability is not merely liveness, endpoint reachability, session
trace conformance, bisimulation, capability reachability, actor identity, or
service continuity. It is the triadic relation among a future obligation, the
carrier that delivers it after a boundary event, and the competent successor
organization that can act and restore a receptive condition.

## Core Argument

1. Computational systems often inherit misleading identity handles: file, PID,
   source tree, checkpoint image, service name, endpoint, mailbox, persona, or
   artifact.
2. These handles do not by themselves settle individuality across handoff.
3. Handoff creates a future-oriented selection test: which future obligation
   can still bind the lineage, through which carrier, and what organization
   must exist for that obligation to be acted on?
4. Different future-task carriers select different self-objects, and receptive
   state restoration distinguishes one-shot response from continuing
   addressability.
5. Quine demonstrates the separation experimentally: source availability stays
   constant while raw-stream versus `ctl/wake` carrier pressure changes the
   viable successor organization.
6. This motivates a carrier-indexed theory of computational individuality that
   couples naturally to successor morphology without collapsing into it.
7. The theory clarifies existing systems rather than replacing them:
   checkpoint/restart, process migration, actors, capabilities, service
   endpoints, and ALife reproduction all instantiate different carrier-indexed
   criteria.

## New Findings From DR Round

The completed four-query `dr` round should reshape the manuscript in four
ways.

First, the paper should frame future addressability as a **triadic relation**:

```text
future addressability = obligation + carrier + competent successor organization
```

This prevents the claim from collapsing into endpoint naming. An address alone
does not make an individual continuous; the successor must be able to satisfy
the future-obligation trace that the ecology imposes.

Second, the systems comparison should foreground the split between:

- **state preservation**: checkpoint/restart, process migration, transparent
  memory preservation
- **organizational reconstruction**: actors, service failover, Kubernetes-like
  control planes, Quine rebuilds

This split is now a central paper axis. Quine belongs on the organizational
reconstruction side, but with a process-native `exec` boundary and a native
control surface.

Third, the formal section should use a small predicate rather than a heavy
formal system:

```text
P --B--> S

Continuity_E(P, S, c) iff:
  Addressable_E(S, c)
  and S satisfies tau_E
  and S is behaviorally adequate relative to c
  and S restores the receptive condition for later obligations
```

The surrounding text can name the nearest formal traditions: process calculi
for channel mobility, capabilities for carrier authority, session types for
obligation traces, temporal logic for liveness, and bisimulation/coalgebra for
observer-relative behavioral equivalence. The paper should not pretend to
complete a full formal proof.

Fourth, the objections are no longer side notes. They define the boundary of
the theory:

- biological autonomy: future addressability is not metabolic autopoiesis
- mechanistic rigor: obligations must be grounded in concrete protocols
- triviality: every reachable endpoint is not automatically an individual
- dormancy: constant live reachability is too strong; use latent
  addressability where appropriate
- substrate capture: if infrastructure does all the preservation, the claim is
  ecology-maintained continuity rather than individual renewal

## Section Skeleton

### 1. Introduction

Open with the problem:

- computational continuity is often inferred from the wrong handle
- the future's route into the system often matters more than the copied object

State the anchor:

> computational individuality is indexed by future addressability.

Use Quine as the motivating anomaly, not as a full evidence retelling.

### 2. The Problem Of Computational Individuality

Distinguish:

- artifact identity
- source identity
- memory / checkpoint identity
- process identity
- endpoint / mailbox identity
- behavioral identity
- obligation-bearing continuity

Show why none is globally sufficient.

### 3. Future Addressability

Define:

- future-task ecology
- future-task carrier
- post-handoff addressability
- carrier-indexed self-object
- carrier topology

Core criterion:

```text
successor S is continuous under ecology E iff,
after boundary B,
S remains addressable through carrier C_E
and can perform E's required future-obligation trace.
```

Refined criterion:

```text
successor S is continuous under ecology E via carrier c iff,
after boundary B,
S remains addressable through c,
S satisfies the future-obligation trace tau_E,
S is behaviorally adequate relative to c,
and S restores the receptive condition for later obligations.
```

### 4. Quine As Motivating Case

Use only the compact contrast:

- `6.63`: source available, raw-stream pressure, thin stream body viable
- `6.64`: source available, `ctl/wake` pressure, rebuilt runtime organization
  viable
- `6.65`: minimized initial prompt preserves the same control-surface basin
- `6.66`: generic `exec` is not closure
- `6.67`: source rebuild is not delayed addressability
- `6M.10`: source/body material, carrier availability, prompt-explained
  implementation repair, and post-handoff proof separate cleanly

Interpretation:

- source visibility is not sufficient
- body rebuild is not sufficient
- formal reconstruction is not operational continuity
- future carrier/capability selects the minimal viable successor organization

### 5. Comparison Across Computational Media

Build one table and make state-preservation versus organizational
reconstruction explicit:

| Medium | Boundary | Future carrier | Candidate self-object | False positive |
|--------|----------|----------------|-----------------------|----------------|
| POSIX stream utility | `exec` | stdin/stdout | stream handler | source/runtime rebuild |
| Actor system | actor restart / code swap | mailbox/address | mailbox-bound behavior | implementation identity |
| Service endpoint | redeploy / failover | network endpoint / capability | endpoint-capable service organization | same container image |
| Checkpoint/restart | restore | checkpoint scheduler | checkpoint image + restore protocol | behavioral rewrite |
| Process migration | migration boundary | preserved channels / OS state | migrated process state | reconstructed lookalike |
| Quine control surface | `exec` | `ctl/wake` | wake/tool/idle runtime organization | source rebuild alone |

Add a short paragraph:

> No system preserves all carriers. Each continuity regime chooses which
> carrier matters: state, endpoint, authority, mailbox, source, or future-task
> control. Future addressability is the criterion that makes this choice
> explicit.

### 6. Theoretical Neighbors And Boundaries

Use nearby theories with explicit limits:

- autopoiesis: upper bound, not current claim
- semantic closure: description / interpretation distinction, not full closure
- closure of constraints: local control-surface closure, not all production
  conditions
- process calculi / behavioral equivalence: trace criterion, not life claim
- actor model: mailbox/address continuity, not biological individuation
- capability theory: authority-bearing carrier, not full identity by authority
- quines / von Neumann: description and construction, not delayed
  addressability

This section should make the novelty boundary explicit:

> We are not first to treat address, state, behavior, description, or
> construction as continuity-relevant. The contribution is to show that, when
> these handles diverge, future addressability provides the operational index
> for which organization counts as the continuing computational individual.

### 7. Implications

Show what changes if the theory is accepted:

- self-replication benchmarks should specify reproduced object and future
  carrier
- systems should distinguish state preservation from future addressability
- ALife systems should ask whether offspring semantics are substrate-given or
  carrier-induced
- agent handoff should report which carrier makes obligations survive

### 8. Limits And Future Work

Non-claims:

- not full life
- not full agency
- not population selection
- not universal identity theory

Next experiments:

- successor-mission ablation
- multi-wake closure
- carrier perturbation
- non-Quine `ctl/wake` adapter
- source degradation / repair
- multi-generation recursion

## Evidence Plan

The paper should use Quine as a theory-generating case, not as the sole proof.

Quine evidence:

- cite through boundary-paper dossier and experiment feeds
- use only representative contrasts
- avoid run-log density unless the venue demands an empirical appendix

Non-Quine evidence:

- examples from process migration, checkpoint/restart, actor systems, and
  service failover should be used as conceptual comparators, not empirical
  claims about Quine.

## Evaluation Plan

The empirical program should prioritize tests that separate addressability
from mere construction:

| Priority | Test | Why it matters |
|----------|------|----------------|
| 1 | minimal non-Quine `ctl/wake` adapter | tests whether the selected object is Quine runtime specifically or any control-surface-compatible organization |
| 2 | successor-mission ablation | tests how much closure depends on the generated continuation contract |
| 3 | multi-wake closure | tests whether receptive condition restoration is stable beyond one delayed task |
| 4 | carrier perturbation | tests whether the successor can diagnose or repair future-task carrier damage |
| 5 | checkpoint/restart contrast | separates state-preserving continuity from organization-reconstructing continuity |
| 6 | source degradation / repair | tests whether self-description is a passive artifact or enters active reconstruction |
| 7 | carrier/capability overlap ladder | separates formal reconstruction, authority/addressability carrier restoration, payload competence, and receptive-state recovery |

Pass/fail criteria should be adapted to our bounded claim:

- pass: post-boundary addressability through the declared carrier
- pass: previously unseen obligation completed through the required trace
- pass: receptive condition restored
- pass: scaffold dependencies declared
- fail: precomputed task completion without post-boundary obligation
- fail: address survives but competent action is lost
- fail: substrate capture is hidden as individual renewal
- fail: recursion or queue growth becomes unbounded

## Citation Plan

Required source families:

- POSIX `exec` semantics
- actor model / mailbox identity
- process calculi / behavioral equivalence
- capability-based systems
- process migration
- checkpoint/restart
- autopoiesis
- semantic closure / epistemic cut
- closure of constraints
- quines and von Neumann self-reproduction
- AI self-replication benchmarks only as contrast

Do not cite the `dr` report directly. Use it as a source-finding and framing
input.

## Language Guardrails

Preferred:

- `carrier-indexed self-object`
- `future addressability`
- `post-handoff addressability`
- `carrier topology`
- `future-obligation trace`
- `control-surface closure`
- `scaffold gradient`

Avoid:

- `the system is alive`
- `full autopoiesis`
- `the LLM copied itself`
- `selection` without `operational` or `viability-induced`
- `same self` without saying under which carrier ecology

Use carefully:

- `identity`: prefer `individuality` or `continuity criterion` unless the
  paragraph is explicitly about identity handles
- `selection`: use `viability-induced` or `operational` selection, not
  population-level evolutionary selection
- `closure`: specify `control-surface closure`, `carrier closure`, or
  `bounded closure`
- `obligation`: ground it in concrete protocol expectations, queued messages,
  probes, tasks, or control events
- `addressability`: specify live, delayed, or latent addressability

## Current Readiness

This is a seed dossier, not a draft-ready paper.

Ready next:

- comparison matrix
- source-checked bibliography
- short theory note
- first formal-definition pass

Not ready yet:

- final venue selection
- full manuscript
- citation-ready related work

## Post-Split Claim Boundary

This paper owns the future-addressability criterion, not the whole theory of
self-reproduction. It should read as a general computational-continuity theory
paper motivated by Quine, actors, services, checkpoint/restart, process
migration, and capability-bearing systems.

## Future Experiment Design For This Paper

The paper already has a detailed roadmap in
[`plans/experimental-roadmap-2026-04.md`](plans/experimental-roadmap-2026-04.md).
The manuscript plan should prioritize the roadmap this way:

| Priority | Experiment idea | Question it answers | Expected manuscript role |
|----------|-----------------|---------------------|--------------------------|
| 1 | minimal non-Quine `ctl/wake` adapter | Is the continuity object Quine runtime specifically, or any control-compatible organization? | discriminates carrier-indexed self-object |
| 2 | successor-mission ablation ladder | How much scaffold is required for post-handoff addressability? | scaffold gradient |
| 3 | multi-wake closure | Does the successor restore a receptive state robustly? | strengthens `idle -> wake -> act -> idle` trace |
| 4 | checkpoint/restart contrast | How does state preservation differ from organizational reconstruction? | theory comparison case |
| 5 | actor/mailbox analog | Does the criterion generalize beyond Quine's process-native control surface? | generality test |
| 6 | carrier perturbation | Can the successor maintain or repair the future carrier? | stronger autonomy boundary |

Design rule:

- future addressability is a triad: obligation + carrier + competent successor
  organization;
- endpoint survival alone is a false positive;
- competence without restored receptive condition is also incomplete.
