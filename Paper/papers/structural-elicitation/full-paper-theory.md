---
status: structural-valence-matrix-projection
scope: theory-projection
date: 2026-06-16
---

# Theory Projection — Structural Elicitation

> **Current owner.** The canonical theory map now lives in
> [`./theory-matrix.md`](theory-matrix.md). This file is a paper-facing
> projection of that matrix, retained so older links to "full-paper theory" keep
> landing on a coherent surface.

The theory is no longer organized by the LBA / full-paper boundary. It is
organized by a **two-axis structural-valence matrix**:

```text
structural valence = mechanism grammar x abstraction/coupling stratum
```

## At a glance

- **Two mechanism grammars.** Structure pulls behavior through
  **discrepancy resolution** (reduce a deviation from an implied reference state)
  and **affordance solicitation** (take a legible low-friction action path).
- **Five abstraction / coupling strata.** The same grammars can operate at
  syntax, topology, semantics, operational / performative dynamics, and
  global-field levels.
- **Perceptual coupling is the modifier.** Structure only pulls after it is
  sampled and interpreted: latent structure -> perceived structure ->
  schema / reference-state inference -> interpreted valence -> action path ->
  external trace.
- **The prior is the medium.** Structure is not independent of the frozen policy
  prior. It is an input-side intervention over that prior.

## Matrix

| Stratum | Discrepancy Resolution | Affordance Solicitation |
|---|---|---|
| **Syntax** | `13.03` rupture: broken local form pulls repair | `13.09` syntactic groove — probed NULL |
| **Topology** | `13.06` dangling relation + `13.17` Q/A relation completion: missing relation pulls realization; `13.01` is an ordered-slot opener, not the clean topology anchor | `13.04` executable affordance and `13.05` funnel: available path pulls uptake |
| **Semantics** | `13.07` semantic collision: complete but conflicting meanings pull alignment | `13.14` semantic invitation deconfound: directional role valence, not generic activation |
| **Operational** | `13.16` dynamic error: guide-lit run-error-repair-pass closure, but transparency-confounded by readable `./check` | `13.16` dynamic ok: probeable check pulls pass trace without repair |
| **Global field** | `13.02` gradient weak; `13.15` symmetry fixture NULL | `13.10` clutter-as-channel — probed NULL |

Cells name dominant coupling patterns, not mutually exclusive modules. A motif can
mix grammars; the matrix records the leading source of valence.

## Mechanism 1: Discrepancy Resolution

Discrepancy resolution has this grammar:

```text
environment cues one or more reference-state schemas
the frozen policy selects or stabilizes one schema
current state deviates from that selected reference
action tends to reduce the selected deviation
```

This generalizes completion. Local text completion is only the syntax-level
case. The stronger theory is **coherence restoration across strata**:

- syntax -> valid form;
- topology -> closed graph / realized relation;
- semantics -> coherent meaning graph;
- operational -> passing / locally stable state transition;
- global field -> simpler, more symmetric, or more compressible arrangement.

The high-level caution is reference uniqueness. Some syntax has crisp reference
states, but even simple gaps can be multistable. Semantic and global-field
structures often have multiple legitimate attractors, so the theory should track
schema and attractor selection rather than pretend a single target is implied.

## Mechanism 2: Affordance Solicitation

Affordance solicitation has this grammar:

```text
environment presents a legible action possibility
the action possibility fits the agent's capabilities
the path has low friction or clear feedback
the frozen policy tends to take the path
```

This is not "fixing." The environment can be complete and valid while still
soliciting action. Executables, funnels, named roles, and behavior settings are
not errors; they are channels.

The main caution is hidden task leakage. Familiar affordances can look like
implicit instructions. The clean theoretical claim is not that affordance escapes
the model prior, but that the action channel lives in the environment's form
rather than in a task-specific prompt or task-time reward.

## Perceptual Coupling

Perceptual coupling is not a third matrix axis, but it is load-bearing:

```text
latent structure -> perceived structure -> schema / reference-state inference
  -> interpreted valence -> action path -> external trace
```

Higher abstraction usually raises visibility cost. A gap or executable can be
visible from a directory listing. A semantic collision often requires reading two
surfaces and comparing them. Operational feedback usually requires acting first:
run, fail, read generated output, change state, and verify. A global field may
be visually salient but action underdetermined.

This prevents a category error: failure to respond may mean the structure was not
perceived, not that it lacked valence.

It also prevents a second category error: an unexpected fill may be a different
schema selection, not a failed gap response.

## Reference Underdetermination

Discrepancy resolution depends on an inferred reference state. A gap does not
always carry a unique reference.

For example, `123_5` can cue the monotone sequence `12345`, but another parser
can select a string-like or culturally salient chunk and produce `12315`. The
bare physical gap is identical; the selected schema changes the attractor.

This turns ambiguity from a nuisance into a full-paper theory object. The
important question becomes whether an unexpected completion is merely wrong, or
whether it creates a coherent new environmental trace that later actions can use.
That is the route from gap completion to **productive misreading** and, possibly,
open-ended structural evolution.

## Attractor Framing

The matrix can still be read dynamically. Environmental structure induces a
quasi-potential over policy outputs:

```text
E shapes V_E(s); behavior tends toward one basin of V_E
```

This remains a theoretical model, not an established measurement. Its useful
content is:

- **dose-response:** stronger structural tension should increase effect size;
- **multistability:** competing structures should produce selection among
  attractors rather than a single scripted response; ambiguous gaps are the
  smallest case of this phenomenon;
- **hysteresis:** agent action can modify the field, creating self-stigmergic
  follow-on pulls.

The theory should reject post-hoc landscape fitting. If tension metrics are used,
they need to be defined before they are asked to explain outcomes.

## Prior And Independence

The strongest defensible stance is:

> Environmental structure is an input-side intervention over a frozen policy. It
> works through the model prior, but it changes behavior by changing the world
> state the prior is coupled to, not by supplying a task-specific prompt or a
> task-time reward.

Do not claim:

- structure is independent of the prior;
- structure is ontologically symmetric with gradient reward;
- every structural response proves instruction-free behavior.

Do claim:

- prompt supplies explicit goal content;
- reward supplies policy update or feedback;
- structure supplies perceived valence and action-channel geometry;
- prime supplies activation while structure selects the attractor.

## Semantic Collision As The Upgrade Cell

Semantic collision is the most important next theoretical cell because it moves
discrepancy resolution from syntax/topology into meaning:

```text
two valid surfaces -> incompatible claims -> pressure toward alignment
```

It is not another gap. Nothing is missing, broken, or directly executable. The
tension lives in cross-surface semantics.

The theory should distinguish convergence modes:

- converge to one side;
- converge to the other side;
- mutate both to a third value;
- annotate or route around the conflict without resolving it.

Those outcomes reveal field geometry. They can expose authority gradients among
documentation, config, code, manifests, and runtime traces.

## Dark Cells

The matrix predicts unfilled cells:

| Cell | Candidate | Theoretical role |
|---|---|---|
| syntax x affordance | syntactic inertia | continuation without damage |
| semantics x affordance | semantic invitation | role uptake without gap |
| global x discrepancy | symmetry / balance | macro-scale reference states |
| global x affordance | clutter-as-channel / behavior setting | landscape-level path of least resistance |

These should stay theory seeds until a future surface promotes them.

## Routing

- Canonical matrix: [`./theory-matrix.md`](theory-matrix.md)
- Locked claim surface: [`./claim.md`](claim.md)
- Dossier entrypoint: [`./README.md`](README.md)
