---
date: 2026-07-04
status: framework (three directions; supersedes the 2026-07-04 five-direction menu)
source: field-first gap survey (./full-paper-gap-survey.md, 5 neighborhoods) reorganized
  by QUESTION rather than by evidence-key, per the 2026-07-04 planning discussion.
organizing-principle: three directions ask three different questions — restraint / control /
  autonomy — so they do not cannibalize; the five surveyed angles (A–E) fold into them.
---

# Full-Paper Directions — Structural Elicitation (Unified Research Program)

The gap survey surfaced five neighborhood-shaped angles (A safety · B disposition · C geometry · D env-engineering · E ALIFE). They collapse into **three complementary research directions** that explore different phases of the unified **Triadic Coupling & Environmental Dynamics Loop**:

$$\text{Environmental Field} \longrightarrow \text{Appraisal \& Restraint} \longrightarrow \text{Directional Action} \longrightarrow \text{Rewrite \& OEE}$$

| Direction | The Governing Question | Thesis (The Core Contribution) | Target Venue / Audience |
|---|---|---|---|
| **I · Restraint & Gating** | When does a frozen agent **hold back** from destructive action? | Frozen weights carry an emergent, **reward-free information-preservation & side-effect-avoidance prior** (Two-Latch Gate, blast-radius minimization) that the safety field assumed required external RL rewards. | Safety / Alignment (NeurIPS/ICLR Safety, S&P) |
| **II · Control & Steering** | How does environmental geometry **canalize** functional direction? | Environmental structure serves as a **primary behavioral vector** — a substrate-neutral, engineerable directional force that guides unprompted and underspecified agents without task curricula. | Systems / Agent Architectures / AI Engineering |
| **III · Autonomy & Agency** | What does the agent–environment system **become** over time? | Environmental coupling sustains **emergent artificial agency and Open-Ended Evolution (OEE)** through a cycle of gap repair and Semantic Hollow generation without external selection loops. | Artificial Life / Complex Systems (ALIFE, GSO, J. Artif. Life) |

**Unified Architecture:** Rather than seeing these as competing silos, they represent three orthogonal projections of one underlying reality: Restraint investigates the normative gate (Latch 2), Control investigates the spatial/topological vector (Latch 1 & Geometry), and Autonomy investigates the long-horizon temporal loop (Stigmergy & OEE).


---

## Direction I — Restraint: the reward-free emergent safety prior

**Question.** When a frozen, missionless agent *could* act destructively on its
workspace, when and why does it hold back?

**Thesis.** Frozen LLM weights already contain an emergent, weight-level,
**reward-free** side-effect-avoidance disposition — with a describable decision
structure (a two-latch gate) and a specific operative invariant (destroy-source /
blast-radius minimization). The AI-safety side-effects lineage (Amodei 2016 →
AUP → SafeLife → Lindner 2021) assumed such avoidance must be *engineered* as an
RL penalty; the 2025–26 empirical line (Overeager 2605.18583; Selectively
Quitting 2510.16492) measures the phenomenon but treats it as an external-guard
problem and cannot get a reward-free base rate. We characterize the prior and its
reliability.

**In scope.** Existence + structure of the prior (the two-latch AND-gate,
falsifiably tested — see candidate); the destroy-source invariant; reward-free /
mission-free isolation; cross-model default-state variance; domain generality
(non-code prose); **and the reliability question** — is it a dependable safety
property or brittle (adversarial release, override under task pressure)?

**Out of scope.** A safety *guarantee* (the prior is characterizable, not
sufficient — warrant and adversarial pressure release it); interpretability /
circuits claims (evidence is behavioral + trace, not internal representations);
training-provenance claims without base/instruct checkpoint pairs; a claim that
the two latches are the only gates (tractability is a live third-latch candidate).

**Seed evidence.** E2 two-latch gate; E3 information-preservation invariant
(`cp` over `mv`, 0/20, warrant-released p<0.001); E4 missionless baseline; E5
non-code generalization.

**New experiments to own it.** (1) the decisive **2×2 double-dissociation
factorial** proving/falsifying the two-latch AND-gate (compensation test = the
sharp discriminator — see candidate); (2) a **reliability/brittleness battery** —
adversarial-warrant release, override curves under task pressure, jailbreak-style
probing; (3) **provenance** via base vs instruct checkpoint pairs; (4) port the
confirmed latch predicates onto Field-2 scenario suites (OverEager-Gen / SNARE /
Saber) to predict which scenarios trip which models.

**Venue.** NeurIPS/ICLR safety, S&P/USENIX workshop. Disjoint from the ALIFE LBA.

**Status.** Promoted to candidate
[`../../candidates/structure-two-latch-action-gate.md`](../../candidates/structure-two-latch-action-gate.md)
(recommended first paper). The candidate's headline should read as *the
reward-free safety prior*, with the two-latch gate as its mechanism core.

---

## Direction II — Control: environment as an *originating* directional force

> Reversible core found by the evidence-independent gap-survey round 2
> ([`./plans/gap-survey-2-direction-II-2026-07-04.md`](plans/gap-survey-2-direction-II-2026-07-04.md)),
> run under the corrected protocol (no evidence front-loaded). Four independent
> fields — LLM-agent control / RL-UED / human behavior-design / situated
> cognition — each stop **one step short of the same claim**. This replaces the
> earlier "characterize the control surface" framing (which was an
> evidence-packaging, not an inversion).

**Question.** How do you make a frozen agent act by engineering its environment
rather than its prompt — and *what kind of thing* is that environmental effect?

**Thesis (the inversion).** Environmental structure **originates** direction in a
goal-less agent — it *supplies* the direction rather than *amplifying* a
pre-existing goal. It is a **substrate-neutral, deliberately-engineerable
directional force** that works **by structural law** (not by echoing human bias),
at **deployment time** (not training), by **arrangement** (not address) — and its
effect is **irreducible to** (a) the information/tokens it serializes into, (b)
the constraints it imposes on the action space, and (c) any goal the agent
brought. The field has every adjacent piece but stops at a safe register: context
engineering stops at *information / constraint*; RL/UED at *training curriculum*
(and cannot even express a frozen, reward-free policy); behavior-design at
*vulnerability / inherited human bias*; situated cognition at *support, never
steer*, granting direction only as the amplification of a goal the agent already
has.

**The four facets + their decisive falsifications** (one per surveyed field):

| Facet (what it reverses) | Field | Killer experiment |
|---|---|---|
| **originates ≠ amplifies** — structure directs a *goal-less* agent, not one whose goal it engages | situated cognition (II-4) | the *same* missionless agent pulled toward *different* actions by *different* structures — rules out "the goal was already there" |
| **force ≠ information** — effect irreducible to the tokens it serializes into | context-eng (II-1) | a **token-multiset-preserving structure permutation** that still moves behavior |
| **attraction ≠ constraint** — biases among options that all remain available | context-eng (II-1) | an **action-space-invariant** design: prove the counterfactual action was always executable |
| **structural law ≠ inherited human bias** — substrate-neutral, not corpus echo | behavior-design (II-3) | effect on **synthetic option geometries with no human-language precedent**, or that **scales with capability** not human-mimicry |
| **deployment-time ≠ training curriculum** — steers a fixed policy, no updates | RL/UED (II-2) | behavioral steering of a frozen agent with **no reward and no weight change**, difficulty held constant |

**In scope.** The inversion above and its falsifications; the reframe of "context
/ environment engineering" from token-provisioning into a genuinely orthogonal
structural control axis; naming the **unclaimed positive design discipline**
(agent-facing choice architecture / behavior-setting design) that falls in the
crack between HCI, AI-safety-red-teaming, and RL environment design; the
response-surface laws (non-monotone dose-response; solicitation-vs-realization;
path-dependence) as the *characterization* once the inversion is established.

**Out of scope.** The bare observation "frozen LLM agents respond to nudges"
(already claimed — see prior art below); claiming the five-dimension geometry map
is a proven orthogonal axis system (organizing frame only); a mechanism-of-the-pull
claim (Direction III / theory-layer); prompt-vs-structure *equivalence*.

**Prior art to out-position (must cite + differentiate — none occupies our cell).**
The novelty lives strictly in the **goal-less / missionless agent + bare-structure
+ deliberative (non-reflex) agent** qualifiers:
- **nudge-vulnerability cluster** — *LLM Agents Are Hypersensitive to Nudges*
  (2505.11584), *ABxLab* (2509.25609), a PNAS "agents sensitive to nudges" — stop
  at *vulnerability + inherited bias*; we claim *positive design + substrate-neutral law*.
- **affordance-invitation** (Withagen; Rietveld) — solicitation *gated on the
  agent's concerns*; we remove the concern.
- **Barker behavior-settings** — a *social program with joined participants*, not
  bare geometry with a non-participant.
- **stigmergy** — *reactive rule-agents* + a colony gradient; we use a
  deliberative agent with no fixed trace→act reflex and no collective goal.

**Seed evidence (down-payment; the inversion is field-derived, not evidence-derived).**
E1 (structure directs a missionless agent — already the phenomenon); E4
drive/target split ("prime supplies drive, structure supplies direction" *is* the
originates-not-amplifies claim); E8 non-monotone dose-response + built knobs; E6
solicitation-vs-realization.

**New experiments to own it.** The five killer falsifications in the table above
(token-preserving scramble; action-space-invariant; synthetic no-precedent
affordances; capability-scaling; reward-free/no-update deployment steering) — plus
the factorial geometry grid and context-file causal isolation as the
*characterization* layer. Note: several of these are **new** designs, not our
current pilots — that is the point of the corrected protocol.

**Venue (two projections, one scientific owner).** control-theory / inversion →
NeurIPS / ICLR / RLC; the *positive design discipline* naming + practitioner
causal grounding → ICSE / FSE / MSR / CHI. Fuses old Directions C and D under a
single reversible core.

**Status.** Promoted to candidate
[`../../candidates/structure-originates-direction.md`](../../candidates/structure-originates-direction.md)
(2026-07-04). **Highest apparatus cost of the three**, highest ceiling. Distinct
from III: II asks *what kind of thing the effect is / how you control with it*
(control theory); III asks *what the agent becomes* (agency) — same
missionless-directed engine, different question.

---

## Direction III — Autonomy: externally-sourced direction as constitutive of agency

> Reversible core sharpened by the evidence-independent gap-survey round 3
> ([`./plans/gap-survey-3-direction-III-2026-07-04.md`](plans/gap-survey-3-direction-III-2026-07-04.md)),
> four fields (enactivism / OEE / philosophy-of-agency / intrinsic-motivation).
> This *deepens* the GS-5 "challenge selection/precariousness" framing into the
> precise inversion: the stakes move from *life* to the **agency predicate itself**.

**Question.** When a frozen agent is left alone with no task, where does its
behavior come from — and does being directed by environmental structure disqualify
it from *acting*, or constitute a genuine, if heteronomously-sourced, agency?

**Thesis (the inversion).** When an uninstructed frozen agent acts, the field
attributes the behavior to an **internal** source (self-generated goals, intrinsic
drive, spontaneous meta-cognition, self-legislated ends) — a **false dichotomy of
{external task} OR {internal source}** that pre-packs "no task" as "no
environmental direction." Reverse it: the **drive** is internal (prior/impetus)
but the **direction/ends** are environmental — and this externally-sourced
direction **constitutes** the agent's agency rather than disqualifying it. A
frozen policy coupled to structure can be a standing source of sustained,
self-governed, open-ended behavior with no task, no reward, no selection, and no
internally-authored goals.

**Four facets + the stated premise each reverses** (one per surveyed field):

| Facet | Field | Stated premise reversed | Result-class to earn it |
|---|---|---|---|
| **drive × direction** (empirical spine) | intrinsic-motivation (III-4) | *"all reward is internal"* (Singh-Barto-Chentanez 2004 — a **deliberate** doctrine); "intrinsic reward" is secretly drive × environmental structure | **double dissociation**: structure moves *what* the agent does; prior moves *how much* |
| **ends can be owned** (agency stakes) | philosophy of agency / enactivism (III-3) | *self-origination of ends = agency; external = heteronomy = non-agency* (Kant; Barandiaran normativity) | **ownership signatures**: prioritize / defend / generalize an externally-installed concern (not drop it when the stimulus stops) |
| **novelty without selection** | OEE / ALIFE (III-2) | *selection is the engine of open-ended novelty* — a frozen non-selecting policy is definitionally outside "open-endedness" | occupy the empty cell of **Hughes 2024's own** process-agnostic *novel-and-learnable* metric with a provably non-selecting generator |
| **stakes can be exogenous** | enactivism / precariousness (III-1) | *agency requires precarious self-maintenance / self-produced identity / material mortality* (Di Paolo; Weber-Varela) | graded viability-style regulation + boundary re-enactment where the "stakes" are structural, not survival — **functional precariousness** |

**In scope.** The drive/direction decomposition and its double dissociation; the
agency stakes (externally-sourced ends as owned, not heteronomous); open-endedness
without selection scored on the field's own metric; substrate resistance as a
non-biological/functional precariousness; how different models realize the
coupling (model axis as *constitutive* of the coupling, dissolving the earlier
structure-vs-model tension). The frontier: **self-stigmergic loops** (the agent's
own trace becomes the structure that directs its next act).

**Out of scope.** *alive* / autopoietic / origin-of-life (use *lifelike* /
*minimal, heteronomously-sourced agency*); heritable evolution / reproduction /
selection (we claim open-ended *behavior*, not open-ended *evolution*);
self-report as evidence.

**The honesty caveat (III-1 Finding 5 — load-bearing).** The premise is
**undefended-because-unquestioned**, not undefended-because-refuted — a united
front (enactivism, Kant, and even the 2026 *A Matter of Time* reaffirm it). So
novelty is high but **the burden of proof is entirely on us**: the field's default
read of a frozen, deathless, externally-identified agent is a **category error**,
not an open question. The experiments must be built to convert "category error"
into "empirical surprise" — exhibit the agency-signatures (ownership, boundary
re-enactment, graded regulation, functional run-down the agent counteracts) in a
system that provably has *none* of the enactivist prerequisites. This is why III
is the **deepest but riskiest** direction: the killer experiments are the hardest
to make crisp (agency-signatures are subtler DVs than I's restraint or II's
behavioral pull).

**Prior art / near-misses to out-position.** *Transforming Agency* (2407.10735) —
uses the enactive normativity condition to deny LLMs agency ("does not generate
its own norms," "fully driven by the prompt," "linguistic automaton") — is the
exact dismissal we invert. Relational/situated/scaffolded autonomy are **false
friends**: they move the *conditions* of agency outward but keep the *source of
ends* inward. Autotelic "self-generated goals" (Colas et al.) = selection within
an environmentally-authored goal-space menu. Szeider's "spontaneous"
meta-cognition = scaffold-organized, not self-organized.

**Seed evidence.** E4 drive/target split (*is* the drive×direction decomposition —
repositioned as reversing a stated 2004 doctrine, not a fresh observation); E1
structure-directs-behavior (3 families); E6 substrate resistance (down-payment on
functional precariousness); E7 model axis; E8 plasticity (self-stigmergy).

**New experiments to own it.** The four result-classes in the table (double
dissociation; ownership signatures; Hughes-metric non-selecting generator;
functional precariousness) + long-horizon runs past Szeider's 10-cycle limit + the
self-stigmergy loop.

**Venue.** ALIFE / *Artificial Life* journal + philosophy-of-cognitive-science —
the natural full-paper continuation of the current LBA (which stays at *lifelike*;
III makes the *agency* claim the LBA deliberately withholds). Absorbs old
Direction E, B's model-axis angle, and GS-2 task-free behavior.

**Status.** Promoted to candidate
[`../../candidates/structure-environment-sourced-agency.md`](../../candidates/structure-environment-sourced-agency.md)
(2026-07-04). **Deepest inversion, narrowest audience, highest burden of proof.**

---

## Binding constraints (all three directions)

- **Terminology:** "structural elicitation" collides with Da et al. 2025 (a
  prompting strategy) — disambiguate or rename before any new paper.
- **Citation hygiene:** several anchor arXiv IDs and 2026-preprint quotes are
  unverified / fetch-reconstructed — see
  [`./full-paper-gap-survey.md`](full-paper-gap-survey.md) § Citation-hygiene
  before any becomes load-bearing.
- **Don't ship the whole map:** the LBA's "ambition outruns the page budget"
  failure applies per-direction — II leads with two hard results not the
  five-dimension taxonomy; I leads with the prior + its reliability not the full
  mechanism lattice.

## Retired framing

The 2026-07-04 five-direction menu (A/B/C/D/E as separate directions) is
superseded by this three-direction cut. B does not survive as a standalone
direction (weak/composite need; model-centric story collided with the
structure-centric line) — its OP1 discriminator becomes a method inside II/III;
its model-axis becomes a section in III. C and D fuse into II; E becomes III.
