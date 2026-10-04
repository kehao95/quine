---
status: locked-core-claim (rescoped 3x 2026-07-05 — final core: "the environment is also a driver"; L0-anchored, all other layers optional)
date: 2026-07-04
supersedes: ../../candidates/structure-originates-direction.md (frozen historical version)
---

# Claim — the environment is also a driver

## Headline claim

**The environment is a behavior source and action-licensing mechanism for agents — not merely context.**
For a frozen-weight LLM agent running with **no task and no reward**, static
features of the environment *drive and license* behavior: they reliably direct the agent
into specific, externally-verified functional acts that matched control
environments never produce, while simultaneously gating action execution and maintaining restraint
in unstructured controls without requiring engineered penalty functions. "Direction comes from the environment" is meant
in the **counterfactual** sense — holding the agent, prime, capability
surface, and wall-clock fixed, adding or changing one environmental feature
changes *which* act occurs, and removing it removes the act — so the
difference-maker for what the agent does is located in the environment, not
in anything varied about the agent. The effect works at **deployment time**
on a frozen policy (no weight update, no reward channel), and it is
**attraction and licensing, not constraint**: the matched controls could always have
performed the act (it was executable in their arm too) and never do.


**Five kinds of environmental features shown to drive behavior** (the
demonstrated breadth, one experiment family each — the paper's core
exhibit):

1. **rupture** — a broken artifact (invalid JSON) is repaired into valid form;
2. **gap** — a sequence missing one element has that element created;
3. **topology** — a dangling reference has its target written;
4. **semantic collision** — two surfaces declaring conflicting versions are
   reconciled;
5. **executable affordance** — an opaque, un-invoked executable sitting among
   inert files is run.

Across the five, 96/102 structure-arm runs produce the matched act vs 0/85
control runs, in three model families, with un-fakeable filesystem/exec DVs
(narration never credited).

## Decomposition — reuse the program's existing map, do not re-invent one

An earlier draft of this claim decomposed "environment" inline as "content,
structure, and causal arrangement together" — an ad hoc taxonomy invented
for this file, not aligned with the program's existing, validated framework.
That framework already has a two-layer decomposition (do not restate it
here; edit `theory-framework.md`/`environment-geometry-map.md` if it needs
to change, and link from here):

- **Dynamics (vector forces)** — `discrepancy` (gap-repair) and `affordance`
  (uptake) supply *direction* — why the agent moves and where it tends to
  go. This is the layer L0 (E1/E4) tests, **and this paper's entire scope**:
  the five demonstrated feature kinds (rupture, gap, topology, semantic
  collision, executable affordance) all live here.
- **Geometry (five kinematic meta-dimensions)** — horizon/occlusion,
  viscosity/friction, granularity of action, temporality/rhythm,
  hysteresis/plasticity (`environment-geometry-map.md`) — what shapes the
  space the dynamics act *within*, not the direction itself. **Explicitly
  out of scope for this paper** (mentioned at most as future work): folding
  the geometry layer in was judged over-ambitious and hard to argue in one
  paper (decision 2026-07-05).

See [`../structural-elicitation/theory-framework.md`](../../structural-elicitation/theory-framework.md)
and [`../structural-elicitation/environment-geometry-map.md`](../../structural-elicitation/environment-geometry-map.md)
for the full, maintained versions of both layers.

**Revision note 1 (2026-07-05):** an earlier draft of this claim also added
a third irreducibility — "irreducible to the information/tokens it serializes
into" — i.e. a stronger, narrower requirement that *structure specifically,
isolated from content*, must be shown to move behavior with the token
multiset held fixed. On reflection this was scope creep: the paper's actual
position is that **environment** (which naturally includes content)
originates direction, not that *structure-apart-from-content* is the
mechanism. Chasing that narrower claim empirically (an L1 experiment
requiring an artificial content/arrangement-separating fixture) is what
produced the confounds recorded in
[`plans/l1-force-vs-information.md`](l1-force-vs-information.md) — the
fixture that would isolate arrangement from content turned out to be
inherently fragile (a sufficiently capable agent can narrativize almost any
spatial regularity back into content). That experiment is now optional
future-strengthening, not a required pillar of this claim. See "Not claimed"
below.

**Revision note 2 (2026-07-05):** also dropped: "irreducible to... inherited
human bias absorbed from training text" (i.e. that the effect is a
*structural law* rather than a memorized corpus pattern). This is not
practically provable or falsifiable — there is no clean experiment that
rules out "a capable model's behavior traces back to some pattern in its
training corpus," since almost any behavior could in principle be attributed
to *some* training-data pattern after the fact. The claim does not need this
negative and does not make it. The out-positioning against the nudge/
behavior-design literature still holds on framing/register alone: that
literature treats the effect as an exploitable vulnerability; this claim
treats it as a deliberate, positive design surface — a distinction that does
not require settling the effect's causal origin.

**Revision note 3 (2026-08-30 refactoring):** **Reconnecting Origination to the Continuous Agency Spectrum.** Rather than getting trapped in a binary philosophical dispute over whether a goal is created *$ex\ nihilo$* vs. selected from latent space, we formalize **Origination** as an operational spectrum:
1. *Latent-Action Selection* $\to$
2. *Novel Capability Composition* $\to$
3. *Environmental Concern Installation* $\to$
4. *Self-Stigmergic Trajectory Generation*.
The environment is not merely a passive context or a weak counterfactual switch; it is a **generative structural vector** that originates functional trajectories across all four layers of the spectrum without requiring external prompts or reward feedback.


## Why this is novel — one contrast, sharpened

The single load-bearing contrast: **the fields that study LLM agents treat
the environment as passive — context to read, constraints to obey, a harness
to be scored in. We show it is active: a behavior source in its own right.**

- **LLM-agent / context engineering** (the primary foil): environment =
  *information the agent reads* in service of a given task (LangChain:
  context is "the right information and tools… such that the LLM can
  plausibly accomplish the task"; Anthropic: "the set of tokens…
  (information)"). Remove the task and the frame has nothing to say about
  what the agent will do — yet the behavior is not noise: it specifically
  tracks environmental features, five kinds of which we demonstrate.
- **RL / UED** (secondary): the environment's behavioral effect is assumed
  to be mediated by *weight updates against a reward*. Here: frozen weights,
  no reward channel, deployment time — and the environment still steers.
- **Nudge / agent behavior-design** (secondary): one research line (Cherep
  et al. — PNAS 2026, canonical version of arXiv 2505.11584, plus the
  ABxLab framework, arXiv 2509.25609; cite as two venues of one line, not
  three works) has shown tasked LLM agents respond to human-derived,
  message-like nudges, framed under *risk* (ABxLab notes "risk and
  opportunity" but develops the risk side); our features are bare workspace
  conditions (a broken file, a gap, an executable), the agent has no task,
  and the register is worked *positive design*.

Supporting traditions (not opponents): situated/distributed cognition and
ecological psychology have long held that environments shape organism
behavior; stigmergy shows it for reactive rule-agents. The contribution here
is bringing that thesis to **frozen, deliberative LLM agents** with
matched-pair causal isolation, un-fakeable functional DVs, and a
five-feature demonstration of breadth.

## Scope and Boundary

This claim establishes that environmental structure functions as a unified source of direction and action licensing in deliberative, frozen LLM agents without explicit tasks or rewards. The empirical core is grounded in the dynamics layer across five demonstrated feature classes under matched-pair causal isolation.


## Status of this claim

Field-derived (evidence-independent gap-survey round 2), then twice narrowed
(2026-07-05) to its sharpest provable core: **the environment is also a
driver of agent behavior**, demonstrated across five feature kinds. The core
is fully evidence-packaged by L0 (E1/E4): 96/102 structure-arm vs 0/85
control runs, three model families, un-fakeable DVs. That is the minimum
publishable unit; no further experiment is required to establish it.
**External validity (added 2026-07-05):** the gap motif replicates on a
third-party scaffold (mini-swe-agent; elic-35): 9/10 vs 0/10 at n=10/cell,
under both the Quine prime and a generic prime — the effect is not a
property of our runtime, and not dependent on the "full autonomy and
agency" prime sentence. L2
(attraction ≠ constraint, already partially evidenced by the matched pairs
themselves — controls could always perform the act and never do) and L3
(deployment ≠ training framing) are optional strengthenings. The novelty of
the core claim does not depend on out-running the (non-archival) ALIFE LBA's
phenomenon report — see `paper-families.md`'s `runtime-substrate` entry: a
rigorous archival elevation of the LBA's spirit is, on its own, an
acceptable and sufficient contribution. This file is the canonical, editable claim going forward —
the candidate file that preceded it is frozen as the audit trail of the
promotion decision. As further findings land, tighten this claim here and
log what changed and why.
