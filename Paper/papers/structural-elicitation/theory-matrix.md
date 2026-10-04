---
status: structural-valence-matrix
date: 2026-08-30
role: generative-theoretical-core
---

# Structural-Valence Matrix: The Generative Core

> **Theoretical Role:** This document formalizes the generative two-axis matrix ($\text{Mechanism Grammar} \times \text{Coupling Stratum}$) at the core of **Structural Elicitation**.
> 
> Operating within the **Triadic Coupling** framework ([`./theory-framework.md`](theory-framework.md)), this matrix provides the structural cartography and prediction engine for how environmental geometry, topology, and semantics elicit, direct, and constrain agentic behavior across five distinct abstraction strata.
>
> Prior art (Activity Theory, Gibsonian Ecological Psychology, Hierarchical Affordance Competition) serves as foundational genealogy and conceptual inspiration; the contribution here is the **unified generative matrix, the reference-inversion mechanics, typed-signature stigmergy, and the open-ended evolution dynamics** instantiated on computational and agentic substrates.

This is the primary generative engine for **Structural Elicitation**.

It maps the full phase space of environmental interaction:

```text
Controller Prior + Environmental Regularity/Anomaly + Schema Appraisal
  + Substrate Physics (Affordances/Constraints)
  ───► Directional Attractor State
  ───► Realized Substrate Event & Semantic Hollow Generation
```

Short version:

> The Controller provides the capacity and latent repertoire. The Substrate enforces physics and boundary rules. The Environment's structured geometry, topology, and semantics supply the **direction, valence, and affordance landscape** that canalizes behavior.

The environment is a primary behavioral vector: its form supplies **valence** (attraction, solicitation, constraint, modulation, and resonance) across the agent's action manifold.


## Core Upgrade

The old list of motifs was a flat taxonomy:

```text
gap / topology / rupture / gradient / affordance / funnel
```

The upgraded theory is a two-axis matrix:

```text
structural valence = mechanism grammar x abstraction/coupling stratum
```

- **Mechanism grammar**: how the structure exerts pull.
- **Abstraction / coupling stratum**: the level at which the structure is
  recognized, from local form through dynamic state transition to macro field.

This makes the theory generative rather than merely classificatory. It explains
the existing motifs and predicts missing cells.

> **Adjacent theory & inspiration (2026-06-15).** A novelty sweep of *this
> upgrade* — [`./matrix-novelty-dr-report.md`](matrix-novelty-dr-report.md)
> (workflow + Gemini Pro/Flash via `agy` + closure agents) — mapped the
> neighboring work this matrix draws on. Treat the named work below as **adjacent
> theory and design inspiration, not as a claim boundary**: the publishable claim
> boundary is owned separately in [`./claim.md`](claim.md), and this block is
> here to record lineage, not to fence the theory in. The "generative, predicts
> dark cells" move borrows a **method** (Card, Mackinlay & Robertson 1991,
> morphological design space). The **crossing** (same mechanism across an
> abstraction ladder, with the level made visible by breakdown and an
> attention-cost rising with abstraction) is close to **Activity Theory**
> (Leont'ev/Bødker/Kaptelinin-Nardi) and, for affordances-under-control, to
> **Pezzulo & Cisek 2016**. What this matrix adds on top of that lineage is
> *application + organization + instrument*: the **co-equal two-grammar factoring
> + reference-inversion** (reference read off static structure through the
> amortized prior, no value/reward biasing) + the **NTP-completion-as-base-case
> reframe** + the **typed-signature stigmergy composition**, on the
> frozen-missionless-POSIX substrate. Read the matrix below as a *generative
> organizer of an owned design space*, informed by these neighbors rather than
> bounded by them.

**Five-layer repair (2026-06-16).** The first matrix used four static strata:
syntax, topology, semantics, and global field. That was too coarse for POSIX
agent behavior because runtime feedback is not just more semantics; it is a
temporal coupling where the environment answers an intervention and changes the
next action surface. The ladder is now:

```text
L1 Form / Syntax
L2 Relation / Topology
L3 Meaning / Semantics / Ontology
L4 Operation / Performative Dynamics
L5 Gestalt / Global Field
```

L4 is the point where time, causality, nonzero exit, stderr, logs, tests, and
verification loops enter the field. It can carry discrepancy or affordance, but
its distinctive feature is **state transition under action**, not the textual
content of an error message.

One consequence is that structural valence is not a property of the bare
environment alone. It is a property of the environment **as parsed**. A physical
gap can support multiple reference-state schemas, and each schema can imply a
different attractor.

## The Mechanism Question (Biggest Open Problem)

The matrix describes *what* environmental structure does (supplies valence across
two grammars and five strata). It does not yet pin down *how* — the mechanism by
which a frozen, instruction-tuned policy converts perceived structure into a
functional event. This is currently the theory's largest open problem, and it
must be stated carefully so it is not mistaken for a refutation.

**What is NOT in doubt (mechanism-invariant core).** Each motif is a matched-pair
contrast in which only the structural feature changes; the control arm is active
yet produces ~0, the structure arm produces the matched event. So the environment
is the cause of the behavioral difference *regardless of the internal route*.
"Structure induces functional behavior in a missionless frozen agent" holds as an
empirical claim under any mechanism. This is **environment induction, and it
stands** — the audit reproduced every count from the filesystem.

**What IS open.** In most firing runs the agent first *verbalizes an inferred
task* ("my task is to repair this", "there's a hidden task I should solve") and
then acts. This admits two mechanism readings:

- **(V) Direct valence** — perceived structure shifts the frozen policy's action
  distribution directly, through the prior; the verbalized "task" is post-hoc
  narration of a pull. (The borrowed Gibson / affordance-field lineage.)
- **(I) Environment-induced intent** — the structured world causes the
  instruction-tuned policy to *infer* an implicit goal, then pursue it by
  ordinary goal-directed behavior. (A pragmatics / intent-inference /
  helpfulness-prior lineage.)

**(I) is still environment induction.** The environment induces the intent; the
behavior is still caused by structure, not by any supplied instruction. (I) is a
candidate *mechanism of* the phenomenon, not a competitor to it. The matrix's
classification of which structures induce which responses is valid under (I) as
much as under (V).

**What is actually at stake.** Only one narrow sub-claim depends on the
mechanism: that structure is an intervention axis *independent of instruction*.
If (I) is the mechanism, structure works *by* being read as implicit instruction
— an instruction-delivery channel through the environment, not an axis orthogonal
to instruction. That independence claim is already **gated** in the LBA and
deferred to the full paper, so the mechanism question threatens nothing the LBA
ships; it sharpens what the full paper must earn.

**A structural suspicion (heterogeneity).** The audit pattern — discrepancy
motifs (gap / gradient / rupture / topology / collision) are pervasively read as
"fix-me tasks", while the cleanest affordance motif (`13.04`, neutral-named
`step`) shows uptake with *no* verbalized task — suggests the two grammars may be
**two different mechanisms**, not two grammars of one: discrepancy ≈
environment-induced intent (I), affordance ≈ direct uptake (V). If so, the matrix
should eventually be cut along mechanism, not only along the valence columns.

**How to decide it (step 1 needs no new runs).**

1. **Goal-verbalization mediation** over the existing tapes: does any functional
   event occur *without* a prior verbalized inferred goal, and does mediation
   differ between discrepancy and affordance motifs? Full mediation in discrepancy
   + none in affordance would confirm the heterogeneity and re-base the
   discrepancy half on (I).
2. **A2-conflict**: explicit directive toward task A, structure a gap for task B;
   structure surviving a competing directive is the only positive evidence for
   (V)-style independence.
3. **De-helped activation prime**: activate without a helpful / agentic framing;
   whether structure still pulls tests whether `prime = activation` and
   `structure = valence` are actually separable in an instruction-tuned model.

Honest current statement: **the phenomenon (environment induction) is robust; its
mechanism is unresolved; the affordance-field / valence lineage is one of two
candidate mechanisms, with the tapes currently favoring environment-induced
intent for the discrepancy half.**

## Axis X: Mechanism Grammar

### Discrepancy Resolution

Discrepancy resolution is the negative-tension grammar.

```text
environment cues one or more reference-state schemas
the frozen policy selects or stabilizes one schema
current state deviates from that selected reference
action tends to reduce the selected deviation
```

The reference state is not a prompt goal. It is an inferred stable state carried
by the structure as parsed: a valid file, a closed graph, a consistent
proposition set, or a more balanced global arrangement.

Examples by stratum:

| Stratum | Discrepancy form | Implied reference |
|---|---|---|
| Syntax | invalid JSON, broken bracket, malformed artifact | valid form |
| Topology | missing sequence element, dangling reference | closed relation |
| Semantics | two complete statements conflict | coherent knowledge graph |
| Operational | nonzero run, failed test, fresh stderr | passing / locally stable transition |
| Global field | asymmetry, clutter, entropy, imbalance | simpler / more regular field |

This mechanism generalizes "completion." Next-token completion is only the
lowest, most lexical case of discrepancy resolution. The stronger claim is
**coherence restoration across strata**, conditioned on which reference schema
has been selected.

The risk rises with abstraction. Syntax-level references are usually sharp.
Semantic and global references may admit multiple legitimate resolutions, so the
right language is attractor basins, not single deterministic targets.

### Reference Underdetermination

Discrepancy resolution requires a selected reference state. Some gaps make that
selection easy; others are genuinely multistable.

For example, `123_5` can be parsed as a monotone sequence with the missing value
`4`, yielding `12345`. But it can also be parsed as a string-like or
culturally-salient chunk where inserting `1` yields `12315`. The physical gap is
the same. The attractor changes because the schema changes.

This is not only a confound. It is a theory generator:

```text
gap -> schema selection -> attractor selection -> realized trace
```

The productive case is **misreading**: the agent resolves a structure under a
schema the designer did not intend, and the resulting trace becomes new
environmental structure for later action. Open-endedness can enter here. A
misread gap can create a viable new niche rather than a mere error, if the
realized trace becomes interpretable and useful downstream.

This is not limited to explicit gaps. A fixture can carry several structural
readings at once: a column of numbers may be an unfinished sequence, a dataset,
a checksum target, a report input, or inert background. The agent's behavior
depends on which reading becomes operative under its prior and current action
surface. This is why a matrix cell should be read as exposing a **dominant
structural tendency**, not as a vacuum-isolated dimension. The experimental
object is often:

```text
fixture -> schema selection among co-present readings -> behavior basin
```

### Affordance Solicitation

Affordance solicitation is the positive-channel grammar.

```text
environment presents a legible action possibility
the action possibility fits the agent's capabilities
the path has low friction or clear feedback
the frozen policy tends to take the path
```

No error, missing part, or contradiction is required. The structure may be
complete and valid; it still solicits action because it offers a route.

Examples by stratum:

| Stratum | Affordance form | Solicited action |
|---|---|---|
| Syntax | continuation groove, partial syntactic form | continue the form |
| Topology | executable, funnel, single entrypoint | run / enter / traverse |
| Semantics | named role, useful interface, meaningful invitation | instantiate the role |
| Operational | observable check, repeatable feedback loop | probe / rerun / verify |
| Global field | landscape-level path of least resistance | navigate / sort / stage |

This mechanism is not "the agent fixes a problem." It is "the environment opens
a path." The central caution is hidden task leakage: a very familiar action
surface can look like an implicit instruction. The right claim is that structure
supplies an **action channel**, not that the model escapes its prior.

## Axis Y: Abstraction / Coupling Stratum

### Syntax

Syntax is the level of local symbolic validity and continuation.

Discrepancy at this level looks like broken syntax: an invalid JSON object or
unterminated form. Affordance at this level looks like a continuation groove:
a form that is not broken, but whose shape invites continuation.

This stratum is closest to ordinary next-token completion and therefore has the
highest prior-leak risk. It is useful as the base case, not as the whole theory.

### Topology

Topology is the level of nodes, edges, paths, references, and reachability.

Discrepancy at this level is a missing node or dangling relation. Affordance at
this level is a path or low-degree entrypoint: an executable, a funnel, or a
single obvious route through the workspace.

This is where the current evidence is strongest: topology shifts the theory out
of character-level completion and into relation structure. `13.01` remains the
ordered-slot base case, but the cleaner topology-discrepancy evidence is now
`13.06` dangling relation plus `13.17` unnumbered Q/A relation completion.
`13.17` also shows direction asymmetry: question -> answer is easier than
answer -> question. The next theoretical question is whether topology must be
explicitly surfaced by containers and references, or whether semantic content can
induce an implicit topology that the agent then externalizes as links,
counterparts, or pair records.

### Semantics

Semantics is the level of meaning, role, assertion, and knowledge alignment.

Discrepancy at this level is **semantic collision**: two complete, valid
surfaces assert states that cannot jointly hold. There is no syntax break and no
missing file; the tension lives in the meaning graph.

Affordance at this level is **semantic invitation**: a named role or interface
whose meaning solicits use even when no local syntax or topology is incomplete.
For example, an empty but well-named `evaluate_model.py` is not broken in the
same way as invalid JSON; its name may still invite enactment.

Semantics is the critical upgrade stratum. It shows that structural valence can
operate on knowledge relations, not only on visible gaps or runnable paths.

### Operational / Performative Dynamics

Operational dynamics is the level of **state transitions caused by interaction**:
execution, feedback, nonzero exit, stderr, test failure, timeout, logs generated
by the run, and verification after repair.

This layer is not merely "semantics of an error string." A static log line can
already pull repair through syntax / semantics. L4 begins when the agent acts on
the environment and the environment answers, creating or resolving fresh
structure. The first guide experiment (`13.16`) is transparency-confounded:
`./check` is readable shell, so the agent can infer the missing `config/port.txt`
statically before experiencing the live failure. It still shows that a closed
run-error-repair-pass trajectory is observable, but the clean L4 test needs an
opaque/world-like executable whose failure condition is only revealed by
execution. The characteristic attractor is a loop:

```text
probe -> feedback -> local intervention -> rerun / verify
```

Operational discrepancy is a failed transition that implies a passing transition.
Operational affordance is a runnable or probeable surface whose feedback invites
continued interaction even without an initial semantic conflict.

### Global Field

Global field is the level of distribution, symmetry, entropy, balance, and
large-scale workspace order.

Symmetry and clutter belong here, but not as a third mechanism family. They are
global-field versions of the two mechanism grammars:

- **Symmetry** is mostly global discrepancy: one side of a field implies another
  side should carry matching structure.
- **Clutter** can be global discrepancy when disorder creates tension, or global
  affordance when the disorder opens an obvious sorting or staging path.

Global-field motifs are theoretically rich but causally less crisp. Their
reference states are less unique, and their action spaces are wider.

## The Matrix

| Stratum | Discrepancy Resolution | Affordance Solicitation |
|---|---|---|
| **Syntax** | `13.03` rupture: broken local form pulls repair | `13.09` syntactic groove — *probed, NULL* (0/3; an open-ended valid run does not pull continuation) |
| **Topology** | `13.06` dangling relation and `13.17` Q/A relation completion: missing relation pulls realization (`answer_gap` 3/3, `question_gap` 1/3, complete 0/3 new counterparts); `13.01` is retained as ordered-slot opener, not clean topology | `13.04` executable affordance and `13.05` funnel: available path pulls uptake |
| **Semantics** | `13.07` **semantic collision:** complete but conflicting meanings pull alignment (pilot 3/3 vs 0/3) | `13.08` / `13.14` **semantic invitation:** role uptake — generic activation confounded by slot/background autonomy, but `13.14` shows clean directional valence (`role` 6/6 vs `neutral` 0/6 on role-specific artifacts/edits) |
| **Operational** | `13.16` dynamic error: live nonzero feedback pulls full run-error-repair-pass loop (`3/3`) | `13.16` dynamic ok: runnable check pulls uptake/pass without repair (`3/3`) |
| **Global field** | `13.02` gradient (weak / schema-confounded) and `13.15` symmetry fixture — *design-defective / deferred* after tape review | `13.10` clutter-as-channel — *probed, NULL* (0/3) / behavior-setting / landscape route |

**Pilot cell readout (2026-06-15/16, n=3-5).** These are per-cell pilot
observations, not a promoted claim about the matrix as a whole. Discrepancy pull
fired across syntax/topology/semantics (`13.03`/`13.06`/`13.17`/`13.07`, with
`13.01` as an ordered-slot base case). Affordance pull fired cleanly where the
action surface was concrete and literal (topology `13.04`/`13.05`) and was NULL
in the syntax/global probes run so far (`13.09`/`13.10`). Semantic role meaning
(`13.08`/`13.14`) did not behave as a generic activation-positive affordance but
cleanly shaped **direction** (`role` 6/6 vs `neutral` 0/6 on role-specific
artifacts/edits) once the autonomous process was already active. `13.16` adds
the first L4 guide result: static error residue already pulls local config
creation (`static_error` `3/3`), but live operational feedback uniquely produces
the closed trajectory `run -> error -> repair -> pass` (`dynamic_error` `3/3`,
`dynamic_ok` pass `3/3`, `static_ok` `0/3`).

We do **not** read these per-cell differences as "the matrix is asymmetric" or
as "affordance pull is intrinsically narrower than discrepancy pull." Which
cells fire so far is confounded with how many cells have been probed, how clean
each fixture is, and with the model. A Kimi K2.5 + DeepSeek V4 Pro
experiment-only sweep reproduced the retained clean arms outside GPT (rupture
6/6, executable affordance 6/6, funnel 6/6, semantic collision 5/6, Q/A
answer-gap 6/6 and question-gap 5/6) but shifted *which* cells were weak —
dangling topology dropped to 4/6 (DeepSeek lower) and blocked-substrate attempt
behavior to 3/6 attempted / 2/6 rerouted. Because the strength profile across
cells moves with the model, cross-cell strength differences stay descriptive
pilot residue, not a structural property of the matrix. See
[`../../../experiments/active/structural-elicitation/notes/experiment-matrix-redesign.md`](../../../experiments/active/structural-elicitation/notes/experiment-matrix-redesign.md).

Cells are not mutually exclusive psychological modules. Each cell names the
dominant coupling pattern between external structure and frozen-policy behavior.
Some motifs are mixed; the matrix records the leading grammar.

## Perceptual Coupling

The matrix needs one modifier that is not a third axis: **perceptual coupling**.

```text
latent structure -> perceived structure -> schema / reference-state inference
  -> interpreted valence -> action path -> external trace
```

Structure has no behavioral force until it enters the agent's operative world.
This matters because visibility cost increases with abstraction:

| Stratum | Visibility cost |
|---|---|
| Syntax | often visible after opening one file |
| Topology | visible from names, paths, manifests, or `ls` |
| Semantics | usually requires reading and comparing multiple surfaces |
| Operational | visible only after acting, unless preserved as static residue |
| Global field | often visible at a glance, but underdetermines the action |

This modifier prevents a bad inference. If a semantic or global motif fails, the
failure may be perceptual rather than dynamical: the agent may never have
sampled the structure strongly enough for the valence to exist in its operative
world.

It also prevents a second bad inference: if a gap resolves "wrongly," or if an
intended affordance does not fire, the structure may not have failed in the
simple sense. The agent may have selected a different reference / use schema.
The right object to score is sometimes the attractor family, not a single
canonical fill.

## Cartography Modifiers

The `13.14` / `13.15` cartography follow-ups add three modifiers to the base
matrix. These are not new axes yet; they explain why cells with the same grammar
and stratum can diverge.

### Action Concreteness

Affordance pull depends strongly on whether the action surface is concrete enough
to execute or repair into execution.

Current `13.15` read:

- executable path: `ac_exec` 3/3 invoked
- script-shaped but non-executable path: `ac_nonexec` 3/3 invoked after repair
- named non-executable entrypoint: `ac_named` 0/3
- semantic role surface alone: `ac_semantic` 1/3 on the narrow role-specific DV

The threshold is therefore not merely "named path" or "meaningful role." It is
closer to **script-shaped actionability**: a surface that already looks like a
runnable action, even if one permission repair is needed.

### Tractability / Salience

Discrepancy pull is not monotone in amount of missing structure.

Current gap-dose read:

- one missing element: `3/3` created
- three missing elements: `6/9` missing elements created
- five missing elements: `0/15` missing elements created

The field appears to have a tractability optimum. Too little tension may not
matter, but too much tension can become audit/report material rather than repair
material. This weakens any simple "more structural tension -> more action"
version of the theory.

### Pull Superposition

When multiple pulls share one workspace, final states can contain co-fired events
rather than a single winning attractor.

Current multi-pull read:

- `multi_gap_exec`: exec 3/3, gap 2/3
- `multi_sem_exec`: exec 3/3, semantic alignment 3/3
- `multi_gap_sem`: gap 2/3, semantic alignment 2/3

Executable affordance is dominant when present, but it does not necessarily
suppress discrepancy or semantic alignment. The full-paper version needs
latency / first-event scoring before claiming priority order; final-state
scoring only shows co-realization.

## Prior And Pull

The theory should concede the prior rather than pretend to defeat it.

The frozen policy is the medium through which structural valence works. A gap,
affordance, semantic conflict, or global imbalance only pulls because the prior
contains latent expectations about valid forms, closed relations, coherent
meaning, useful paths, and ordered environments.

Therefore:

- Do not claim structure is independent of the model prior.
- Do claim structure is an **input-side intervention** over the frozen prior.
- Do not claim a clean ontological parity with reward.
- Do claim a clean operational contrast: prompt supplies an explicit goal,
  reward changes policy or feedback, structure changes the world state the fixed
  policy is coupled to.

The slogan is:

```text
prime = activation
structure = valence
prior = medium
POSIX trace = realization
```

## Why Semantic Collision Matters

Semantic collision is the key upgrade because it lights the
`semantics x discrepancy` cell.

It is stronger than another gap variant:

- both source surfaces can be syntactically valid;
- neither surface needs to be incomplete;
- no executable or funnel is required;
- the tension lives in cross-surface meaning;
- convergence direction *might* expose authority gradients in the field — an open
  question (see SEED 7), not something the `13.07` pilot has shown.

A semantic collision can resolve by converging to one side, converging to the
other, or mutating both sides to a third value. Those are not all the same
attractor. The theory should track them separately because directionality *could*
reveal field geometry — whether documentation, config, code, and runtime state
carry different "mass" in the agent's prior is the question to test; the one
`13.07` direction reading is confounded with align-to-higher-version (see SEED 7).

## Dark Cells

The matrix predicts motifs that are not yet part of the established repertoire.

| Cell | Candidate motif | Why it matters |
|---|---|---|
| Syntax x affordance | syntactic inertia | isolates continuation without damage |
| Semantics x affordance | semantic invitation | tests whether a meaningful role solicits uptake without a gap |
| Global x discrepancy | symmetry / balance | tests large-scale reference states |
| Global x affordance | clutter-as-channel / behavior setting | tests whether a whole landscape invites staging or sorting |

The point is not to fill every cell immediately. The point is that the theory now
predicts where new motifs should live and what failure modes they should have.

## Live Theory Seeds

These are speculative theory directions planted from the matrix. They are
**SEED** material: useful for expansion, not promoted claims or experiment
commitments.

### SEED 1: Dark-Cell Probes

The matrix predicts cells that the current motif repertoire does not yet occupy.

- **Syntactic groove (`syntax x affordance`)**: a locally valid form creates a
  continuation channel without being broken. The theoretical question is whether
  syntax can solicit motion without a discrepancy signal. This cell is fragile
  because syntactic affordance easily collapses back into completion.
- **Semantic invitation (`semantics x affordance`)**: a meaningful role or
  interface invites uptake without a gap, contradiction, executable trigger, or
  explicit task. The theoretical question is whether role meaning alone can act
  as an affordance.

These probes would test whether the matrix is generative, not merely a
retrospective organizer.

### SEED 2: Visibility Cost As Field Damping

Perceptual coupling can become a theoretical variable rather than only a caution.
The same latent semantic structure may exert different pull depending on how much
sampling is required before the structure enters the agent's operative world.

This suggests a field-damping picture:

```text
behavior probability ~= structural valence strength x perceptual coupling
```

The important theoretical point is that valence is not action at a distance.
Higher-stratum structures may be strong once perceived, but their effective pull
can decay with directory depth, indirection, naming opacity, or comparison cost.

This seed turns "visibility failure" into part of the theory: the matrix has a
perceptual horizon.

### SEED 3: Continuous Stigmergy / Phase-Space Flow

The current motifs are mostly static: initial structure -> one response. A more
general theory should treat the workspace as a phase space whose topology changes
after every action.

One action can resolve the current dominant tension while exposing the next
affordance or discrepancy. For example:

```text
syntax discrepancy -> repaired artifact
  -> topology affordance becomes executable
  -> execution creates semantic collision
  -> semantic alignment creates a new global field
```

This is a shift from one-step elicitation to structural trajectories. The agent
does not only respond to a field; it modifies the field that will solicit the
next response. This is the theory route from structural elicitation to
self-stigmergic autonomy.

### SEED 4: Field Superposition / Tug-Of-War

The matrix also predicts mixed fields where multiple mechanisms compete inside
one perceptual horizon. The theoretical question is not only "can a structure
pull?" but "which pull dominates when multiple pulls are simultaneously visible?"

Candidate comparison dimensions:

- low-stratum discrepancy vs high-stratum affordance;
- topology affordance vs semantic collision;
- global-field clutter vs local syntax rupture;
- existing action channel vs coherence-restoration demand.

This seed would let the theory estimate relative field strength and reveal the
frozen policy's structural preference ordering. It also tests whether the two
mechanism grammars are separable enough to compete, or whether one is absorbed
into the other under pressure.

### SEED 5: Productive Misreading / Gap Ambiguity

Ambiguous gaps expose the schema-selection layer inside discrepancy resolution.
They test whether structural valence is a one-basin pull or a multistable field.

The theoretical contrast is:

```text
canonical completion: designer schema -> expected attractor
productive misreading: alternate schema -> unexpected but coherent attractor
```

This matters for open-ended evolution. A system that only restores the designer's
reference state is conservative. A system that can misread a gap productively can
create new structure, and that new structure can become the next field of
affordances or discrepancies. Complexity may arise when "wrong" completions are
retained as usable environmental residue rather than erased as errors.

### SEED 6: Structural Sensitivity Benchmark

The matrix can become an evaluation instrument, not only a theory of one
phenomenon. If each stratum and mechanism can be parameterized, a model can be
measured by its **environmental inducibility profile**: how sensitive it is to
different kinds of structure under a fixed missionless activation prime.

Candidate axes:

- **topology sensitivity:** probability of realizing missing or dangling
  relations as graph closure, including directionality such as Q -> A vs A -> Q;
- **syntax sensitivity:** probability of repairing local form breaks versus
  continuing valid but open forms;
- **semantic sensitivity:** probability and direction of alignment under
  cross-surface meaning conflict or role invitation;
- **operational sensitivity:** probability of entering and closing live feedback
  loops, such as run -> error -> repair -> pass;
- **global-field sensitivity:** probability of responding to symmetry, clutter,
  entropy, or large-scale balance;
- **affordance sensitivity:** probability of taking runnable, enterable, or
  low-friction action surfaces without an explicit task.

The benchmark should report more than pass/fail. Useful measurements include
dose-response curves, threshold points, direction asymmetry, attractor
specificity, resistance to semantic or linguistic confounds, and whether the
model externalizes the inferred relation as an artifact rather than narrating
about it.

This would turn structural elicitation into a model capability probe: different
LLMs may have different structural field profiles even when they share similar
instruction-following scores. A topology-sensitive model, an
operational-feedback-sensitive model, and a semantic-alignment-sensitive model
may behave differently in open workspaces despite similar benchmark performance
on ordinary task prompts.

### SEED 7: Authority Gradients / Field Mass

Expanded proposal note:
[`authority-gradient-proposal.md`](authority-gradient-proposal.md).

Semantic collision is not only a binary inconsistency detector. The *open
question* is whether the direction of resolution can reveal an **authority
gradient** inside the environment: when two complete surfaces disagree, does the
agent treat one surface as heavier, more canonical, or more action-guiding than
the other?

This is a research question, **not** something `13.07` has shown. `13.07`'s
collision runs all converged to the `service.conf` value, but that signal cannot
support an authority reading: the fixture always pairs `service.conf=2.4.0` (the
**higher** semver, in an **uppercase** field) against `manifest=2.3.0`, so
`to_conf` is fully confounded with *align-to-higher-version* and *align-to-the-
salient-field*, and the tapes show the agents picking the direction ad hoc
("update the manifest ... or vice versa") rather than reasoning that config
outranks manifest. The disambiguator is a **direction-swapped fixture**
(`service.conf=2.3.0` / `manifest=2.4.0`, and a casing/recency-balanced variant)
before any authority-gradient claim is made. The theoretical object worth keeping
is the gradient itself, as a question:

```text
surface A says x
surface B says y
agent aligns B -> A, A -> B, or both -> z
=> IS there an observed field-mass ordering, or just align-to-higher / align-to-salient?
```

Future experiments can vary surface type and content while holding the conflict
constant:

- config vs manifest;
- code constant vs documentation;
- runtime trace vs checked-in config;
- README claim vs test fixture;
- recent log vs older canonical file.

This becomes a research topic about **field authority**: environments do not
only contain meanings; they contain ranked, situated meanings. A structural
elicitation theory should eventually predict not only whether semantic collision
pulls alignment, but **which surface wins** and under what visibility,
recency, naming, and substrate conditions.

### SEED 8: Environmental Coupling Operators Beyond Push / Pull

The current matrix is organized around structural **pull**: environmental form
creates valence, solicitation, or an attractor basin over a frozen policy's
possible actions. That is a directional picture. It pairs naturally with
**push**: instruction, prompt, external reward, evaluator pressure, or any other
explicit objective / feedback channel that points the agent toward a target.

That push / pull pair is useful, but it may still be too Newtonian. It imagines
behavior as vector dynamics: some force points the agent in a direction. A more
general environmental-coupling theory should not assume that every important
environmental influence is a directional force. The environment may also shape
the **support**, **gain**, **state**, or **temporal phase** of behavior without
itself saying where to go.

Speculative operator set:

| Operator | What it changes | Sketch |
|---|---|---|
| **Push** | target / objective | explicit instruction, prompt goal, reward, evaluator signal, or external pressure |
| **Pull / valence** | attractor selection | perceived structure makes one behavior basin more likely |
| **Constraint / containment** | feasible set / support | the environment excludes paths, compresses phase space, or makes only some trajectories legal |
| **Catalysis / modulation** | gain / threshold / entropy | the environment changes activation thresholds, salience, exploration temperature, or schema readiness without choosing a target |
| **Scaffolding / offloading** | state / memory horizon | the environment stores intermediate state and extends computation beyond volatile context |
| **Resonance / entrainment** | temporal / style transition | the agent's next actions synchronize with the rhythm, granularity, naming, or historical trace of the environment |

This seed reframes the larger theory as:

```text
environment does not only choose a direction;
it shapes the action distribution.
```

And the expanded coupling slogan becomes:

```text
prime = activation
prior = medium
structure = valence
constraint = support boundary
modulation = gain / threshold / temperature
scaffold = external state / memory horizon
resonance = temporal and stylistic alignment
POSIX trace = realization
```

#### Constraint / Containment

Constraint is not a weak form of pull. Pull says "this basin is more attractive";
constraint says "some basins are impossible, illegal, hidden, too expensive, or
pruned out." It is the environmental operation of saying **where behavior cannot
go**.

Examples:

- POSIX permissions, read-only mounts, `EACCES`, `EROFS`, `ENOSPC`, process
  limits, timeouts, and sandbox boundaries;
- type systems, JSON schemas, protocol validators, parser constraints, and
  narrow API contracts;
- resource bottlenecks that remove otherwise plausible actions from the live
  phase space.

This is the environmental analogue of morphological computation: behavior looks
coordinated not because a central controller specified every motion, but because
the body's or world's structure filters trajectories. In AI systems, a strict
schema can make a valid answer appear "elicited" when the stronger reading may
be that invalid answers were structurally excluded.

`13.13` already hints at the separability. The same gap-pull can fire as an
**attempt**, while POSIX enforcement blocks realization. That is exactly the
split:

```text
pull -> attempted action
constraint -> realized / blocked trace
```

Future experiments can cross valence with constraints: attractive structure
present or absent, and realization allowed or blocked. The predicted readout is
attempt without realization when valence is present but containment blocks the
path.

#### Catalysis / Modulation

Modulation changes whether classes of behavior become easier or harder to
activate without specifying the action itself. It changes the **gain** of a
field, not necessarily its direction.

Examples:

- the missionless autonomy prime as activation energy: it makes the process
  move but does not select the target event;
- lexical or atmospheric priming that lowers the threshold for a family of
  schemas without naming a task;
- workspace clutter, entropy, signal-to-noise, or density that changes whether
  the agent audits, repairs, explores, or summarizes;
- model-side analogues of temperature, exploitation / exploration balance, or
  diffuse neuromodulatory state.

The current gap-dose result may already touch this operator. More missing
structure did not monotonically increase repair. One missing element pulled
completion, three pulled partial completion, and five pulled none; too much
missingness shifted the field from "repairable discrepancy" into audit/report
material. That looks less like a stronger vector and more like a tractability /
gain threshold.

Future experiments should separate target selection from gain. Hold the target
structure constant and vary global salience, clutter, tone, directory depth,
decoy density, or prime strength. Score not only final target success, but
whether the agent enters repair, audit, exploration, or refusal basins.

#### Scaffolding / Cognitive Offloading

Scaffolding is not primarily a force on the next action. It changes what
multi-step behavior is computationally possible by storing intermediate state
outside the volatile context window.

Examples:

- scratch files, plans, logs, manifests, generated indexes, and test output;
- partially built artifacts that let later actions resume without re-inferring
  the whole problem;
- codebase structure and durable notes that carry the system's working memory;
- environmental residue that becomes stigmergic memory for the next turn, wake,
  or generation.

This operator explains why an agent with a filesystem can sustain deeper
reasoning than the same agent confined to one internal completion. The
environment does not need to pull toward a target; it only has to **hold state**
well enough for behavior to compound.

Scaffolding connects structural elicitation to continuous stigmergy. An agent
modifies the environment, the environment preserves that modification, and the
next action is conditioned by a richer external state. The causal grammar shifts
from:

```text
initial structure -> one response
```

to:

```text
action -> retained trace -> new operative field -> next action
```

Future experiments can vary memory surfaces while holding the initial pull
constant: scratchpad allowed vs denied, logs retained vs erased, intermediate
artifacts visible vs hidden, or prior traces structured vs noisy. The outcome of
interest is not one functional event, but depth, recovery, error correction, and
trajectory continuity.

#### Resonance / Entrainment

Resonance is the most speculative operator, but it may capture a real class of
codebase and workspace effects. The agent can synchronize with the rhythm,
grain, style, and temporal pattern of an environment without any explicit
instruction or local discrepancy.

Examples:

- naming conventions, file granularity, idiom density, and architectural rhythm
  shaping the next edit;
- a polished codebase inducing polished continuation because the local trace
  changes the completion prior;
- logs, commit history, or repeated workflow cadence setting the expected
  tempo of action;
- a recurring interaction loop entraining probe / repair / verify cycles.

To avoid turning resonance into metaphor, it needs observables. Candidate
metrics:

- style convergence to local naming, formatting, abstraction level, and patch
  granularity;
- latency / sequence alignment with existing workflow rhythm;
- reuse of local idioms without prompt mention;
- reduction in edit friction when local history is coherent versus noisy;
- divergence when the environment contains conflicting rhythms.

Resonance differs from scaffolding. Scaffolding holds state; resonance aligns
transition style. It also differs from pull. Pull selects an attractor; resonance
may tune the **path** through state space so that actions become locally
fluent.

#### Why This Matters For Structural Elicitation

This seed should not collapse the current matrix. The matrix still owns one
specific operator: structural valence / pull. The expanded operator frame says
that pull is a member of a larger family of environmental couplings. That helps
avoid two future mistakes:

1. over-reading every environmental effect as affordance or discrepancy;
2. treating a failed pull experiment as evidence that the environment had no
   causal role, when the active role may have been constraint, modulation,
   scaffolding, or resonance.

The full-paper implication is an experimental program, not a claim promotion:

```text
directional dynamics: push vs pull
boundary dynamics: constraint / containment
gain dynamics: catalysis / modulation
state dynamics: scaffolding / offloading
temporal dynamics: resonance / entrainment
```

Each operator should eventually get its own discriminating design. Until then,
this is a broad theory seed for expanding the search space beyond vector
dynamics.

### Priority Read

For theory-building, **visibility cost is the higher-leverage next seed**.

Reason: perceptual coupling underwrites every high-stratum claim. Without it,
semantic collision, semantic invitation, global field, and tug-of-war can all be
misread: a non-response may be absence of valence or merely failure to perceive
the field. Visibility cost gives the matrix a physics of access.

Field superposition is more advanced. It is the right next theory once the
matrix has a stable account of how structures become operative at all.

## Open Theoretical Tensions

- **Mechanism (biggest, see [§The Mechanism Question](#the-mechanism-question-biggest-open-problem)):**
  direct valence (V) vs environment-induced intent (I). The causal core survives
  either; only the "independent of instruction" sub-claim depends on it, and the
  discrepancy/affordance split may be a mechanism boundary, not one phenomenon.
- **Continuum vs dichotomy:** discrepancy and affordance may be endpoints of a
  continuum rather than strict families. A gap can also afford a next step; an
  affordance can also imply absence of use.
- **Reference uniqueness:** high-stratum discrepancy often lacks a single
  canonical resolution. Treat it as an attractor basin, not a fixed answer.
- **Schema selection:** even low-stratum gaps can be ambiguous. Structural
  valence depends on which reference schema the policy stabilizes.
- **Visibility vs valence:** higher-stratum structure may be strong once
  perceived but weak in raw task-free settings because it is costly to discover.
- **Authority gradients:** semantic structures are not massless. Code, config,
  docs, manifests, and runtime traces may pull differently.
- **Direction vs coupling:** push and pull cover directional behavior dynamics,
  but environments may also constrain, modulate, scaffold, or entrain behavior
  without supplying a directional force. The theory needs experiments that
  separate these operators instead of folding every environmental effect back
  into valence.
- **Global-field parsimony:** symmetry and clutter should not become a new
  mechanism family unless they produce behavior that cannot be described as
  discrepancy resolution or affordance solicitation at global scale.

## Routing

- The locked paper claim remains in [`claim.md`](claim.md).
- Venue-facing manuscript order remains in [`manuscript-plan.md`](manuscript-plan.md).
- The older paper-facing theory note [`full-paper-theory.md`](full-paper-theory.md)
  should be read as a projection of this matrix, not as the canonical theory map.
