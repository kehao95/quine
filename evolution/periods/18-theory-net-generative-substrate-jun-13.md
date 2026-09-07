## Period XVIII: Theory-Net Generative Substrate (Jun 13–14)

The theory layer was audited as a substrate for autonomous theory growth and
found to be a designated-but-uncoupled variation medium — six `THEORY.md` silos
joined by a detached human-facing index, the repository's own reflexive
finding #2 ("variation media nobody writes stay empty") made flesh on its own
theory layer. A bounded structural repair turned the theme nodes into a
traversable, agent-extensible net with node-local growth surfaces along three
axes, coupled the net to the ordinary work loop, and gave it a wake model.

### Jun 13 — Theory net co-located, coupled, and sensed

**Key commits:** `theory-generative-substrate` branch landing.

**What changed:** Six moves, all reuse of existing machinery rather than new
mechanism.

- *Sideways (coherence).* Added a co-located, typed `Sibling Theory Edges` table to all
  six `Paper/themes/*/THEORY.md` (verbs `depends-on` / `specializes` /
  `shares-evidence` / `tensions-with`), linking sibling `THEORY.md` (not
  `README.md`), seeded as a projection of `research-map.md` /
  `control-plane-topology.md`. The three unowned cross-theme gaps
  (branching-identity, time-semantics revisability, collective continuity) are
  surfaced as `OPEN (unowned)` edges on both endpoints with provisional owners.
  The contract lives once in `Paper/themes/README.md` § Node-Local Theory
  Surfaces; the old prose "Borrows from" lines became pointers to the typed
  table.
- *Down (depth).* Added an `Open Depth Front` line to each `THEORY.md` naming the
  theme's highest-yield maturity jump plus candidate (often unclaimed) evidence,
  and an `Unclaimed Evidence Register` to `Paper/evidence-map.md`.
- *Up (meta).* Extended the Strange Loop and Amendment Protocol in `AGENTS.md` so
  the theory layer is a first-class repair target (`Theory Drift → Framework
  Defect → Theory Repair → THEORY.md / evidence-map`), falsified by the
  `claim-maturity.md` Demotion Rule.
- *Routing + coupling.* Added a theory row to the `AGENTS.md` Routing Matrix, a
  bounded theory-front increment and `Theory` / `Synthesis` activation to the
  Maintenance Keyword (and `maintenance-pass.md`), and a theory-net floor to
  `CLOSURE.md` and `Paper/contracts.md` so ordinary closure emits a theory delta
  or a named no-op.
- *Wake model.* Added `scripts/check-theory-coherence.sh` (mechanical v1,
  presence + cross-node gap sync), discovered by `heartbeat-readout.sh` as
  `audit-theory-coherence`; advisory/heartbeat-scoped via the
  `check-validator-registration.sh` MANUAL_ONLY allowlist.

**Friction:** A two-workflow structural review showed theory had grown exactly
once in repo history, and only under heavy human-orchestrated multi-agent work.
The drop-in default gradient pointed at maintenance/closure/papers, never theory:
`THEORY.md` files carried zero cross-links to each other (nodes and edges were
not co-located), zero open-question language (closed-form prose), and the
meta-loop terminated at the Constitution. The empty `open-ended-ecology/papers/`
directory was the field confirmation of finding #2 on the theory layer itself.

**Structural repair:** Co-locate edges with nodes and type them; make the
unowned gaps and depth fronts first-class node-local invitations; widen the
meta-loop terminus to the theory layer; couple theory deltas to the closure loop
so the medium fills as a byproduct of finishing work; add a detection-only sensor
so a flat or drifting theory node surfaces in the heartbeat. The confidence bar
was deliberately relaxed from pure unprompted drop-in to "an agent given a
direction is then carried by structure along down / sideways / up."

**What it taught:** The repository's own habitat findings are sharp instruments
against its own theory layer, not only its control plane. A variation medium
stays empty until it is structurally coupled to the work and made perceptible at
the point where a working agent already stands; an index a human reads is not a
substrate an agent inhabits. Two follow-ups remain open: the sensor's promotion
from advisory to a `check-control-plane` gate once proven quiet-when-clean, and a
real single-session drop-in test to confirm the attractor is observed, not
assumed.

### Jun 13 — Drop-in test passes, but exposes the meta under-pull; Meta axis grown explicitly

**Key commits:** `theory-generative-substrate` branch (meta-axis landing).

**What changed:** A blind single-session drop-in test (a fresh agent given only
"Theory: autopoiesis-renewal — grow it", no orchestration) confirmed the
substrate carries an agent: it navigated AGENTS.md → themes/README → THEORY.md
friction-free and was handed axis-tagged growth targets from the node alone. But
it produced two DOWN proposals and zero UP (meta) ones — confirming the owner's
diagnosis that the model is strongly pulled to deepen/decompose and barely at all
to synthesize, and that the existing node surfaces (Sibling Theory Edges =
sideways, Open Depth Front = a maturity-ladder magnet = down) had no meta pull.
In response the Meta/up axis was grown explicitly: a third node-local surface,
`## Meta / Synthesis Front`, was added to all six `THEORY.md` (an upward
`Spine alignment` edge into `Paper/philosophy/self-maintaining-systems.md` with a
closed verb set {instantiates|serves|specializes|tensions-with}, plus a
`Throughline question` that must be a synthesis/framework question, not a depth
step), contracted once in `themes/README.md`. The `Theory` pass was made
**spine-first and gating** in `AGENTS.md` / `maintenance-pass.md`: state the
throughline + spine alignment first, and any depth/sideways increment must name
the throughline/spine it serves or be deferred. Closure (`CLOSURE.md`,
`Paper/contracts.md`) now checks the Synthesis Front is not left stale, and
`check-theory-coherence.sh` gained presence + spine-link + closed-verb checks for
the new surface. A 14-agent workflow drafted the six Synthesis Fronts and
adversarially verified each was genuine meta (not depth-in-disguise) before
landing.

**Friction:** The drop-in test made the axis asymmetry concrete — even with the
up-axis machinery present in AGENTS.md (as reactive *repair*), a directed agent
still rabbit-holed into depth because no node-local surface *invited* synthesis
and depth was the path of least resistance. Meta is the directive axis (it lets a
future agent grasp the research thread and choose well), so under-pulling it actively
amplifies the fragmentation the whole exercise targets.

**Structural repair:** Give meta a node-local generative home symmetric to the
other two axes (an upward edge + a throughline question), make the spine a
first-class meta object themes point up to (philosophy file stays canonical
owner), and gate depth/sideways behind naming the throughline they serve so
ungrounded depth is deferred, not done.

**What it taught:** Building an axis "in name" (a repair clause in the
constitution) does not create traction; an axis pulls only when it has a
node-local surface a working agent lands on. And the model's depth bias must be
*countered* structurally (gate depth on the spine), not merely *balanced* by
adding a third option. Open follow-up unchanged: promote the sensor to a gate
once quiet-when-clean.

### Jun 13 — First synthesis through the meta axis: carrier as the shared atomic unit

**Key commits:** `theory-generative-substrate` branch (carrier-synthesis landing).
Human-directed synthesis pass (re-theorizing the spine is human-directed per
AGENTS.md, not an autonomous maintenance increment).

**What changed:** Two blind drop-in theory-growth runs (autopoiesis-renewal and
open-ended-ecology), each forced by the meta-first gate to orient on the spine
before deepening, independently surfaced the same primitive: the **carrier** — not
the process, not the source — as the shared atomic unit. Held together with the
substrate, this became a canonical synthesis section, `Carrier Succession As The
Shared Atomic Unit`, in `Paper/philosophy/self-maintaining-systems.md` (elaborating
the "Carrier Succession" already named in proposition 1 of the Grand ALife
Synthesis). It states: `exec`/death is a Weismann-like inheritance barrier that
only externalized surfaces (carriers) cross; individuality (delayed
addressability) and inheritance (carrier-class selection under a bottleneck) are
two questions about the same carrier; the carrier appears at three loci —
in-lineage (autopoiesis), in-habitat (viability's extended phenotype), in-trace
(OEE's stigmergic residue); carrier is necessary-not-sufficient (transmission, not
competence); and the theory homes reorganize autopoiesis's existing
Neighboring-Theory Backbone by shared object, adding Dawkins replicator/extended
phenotype, Jablonka–Lamb inheritance systems, and Grassé stigmergy as proposed
homes. The three themes' Meta / Synthesis Front spine-alignment lines were tagged
with their carrier locus and linked up to the synthesis; the spine's Nearby
Surfaces now names all three legs. Carrier mechanics stay owned by
`autopoiesis-renewal/theory/core-relations.md`.

**Friction:** The carrier concept was already developed inside autopoiesis-renewal,
and carrier-shaped objects existed in viability (extended phenotype) and OEE
(inherited artifact / trace), but under three different names with no shared owner
— a cross-theme object scattered across themes, invisible until two probes hit it
from different angles at once.

**Structural repair:** Give the cross-theme synthesis a single canonical home (the
spine, which already named Carrier Succession), keep the operational mechanics with
their existing owner, and wire each theme's upward edge to the synthesis so the
shared primitive is discoverable at every node.

**What it taught:** This is the meta axis paying off as designed: orienting on the
spine before deepening produced a spine-level synthesis (thorough mastery) that pure depth
could not — the carrier was only visible when two themes were held against the
spine together. It is also the first content evidence that the up axis grows
theory, not just navigation. Honest status: a directional synthesis with proposed
cross-locus and literature connections, both carrier experiments still unrun.

### Jun 13 — Narrative-collapse recorded + refactor: one substrate property, three faces

**Key commits:** `theory-generative-substrate` branch (substrate-property unification pass).
Human-directed synthesis correction (re-theorizing the spine is human-directed per
AGENTS.md, not an autonomous maintenance increment).

**What changed:** Recorded the dynamic/population complement of the carrier
synthesis — the narrative-collapse frontier — and, while doing so, caught and
repaired section-accretion. A 108-agent deep-research pass positioned the claim
honestly (saved at
`Paper/themes/open-ended-ecology/references/narration-collapse-novelty-positioning-dr-report-2026-06-13.md`):
the components are established named prior art (Goodhart / Manheim & Garrabrant;
specification gaming / Krakovna, Karwowski; Tool-Call Hacking / Proof-of-Use —
mostly AI-safety/RL, not ALife) that must be cited, not renamed; the genuinely
unnamed contribution is the 4-way conjunction in an LLM-as-organism substrate. The
two sibling sections landed earlier today (`Carrier Succession As The Shared Atomic
Unit` and `The Functional-Narrative Selective Boundary`) were unified into a single
`## The Substrate Removes Free Un-Fakeability Of State` in
`Paper/philosophy/self-maintaining-systems.md`, framed by that one property (state
is expensive to make and carry but cheap to describe) and deriving three faces:
TRANSMISSION (Carrier Succession), SELECTION (narrative collapse / "evolutionary
Goodhart under free narration", cure = un-fakeable functional-dependency edge,
grounded in verifiable computation not Zahavian costly signaling), INHERITANCE
(counterfactual culture — faithful inheritance ≠ truth). Added a Novelty Posture
note (substrate-driven recombination; cite the parts, claim the conjunction; never
claim to invent Goodhart/spec-gaming/costly signaling) and a Status / Open-Limbs
block (Bedau complexity, cultural-evolution corollary home, Lehman & Stanley
objective-deception, exec-exit-0 adversarial stress test). Every claim, citation,
qualifier, and empirical number from both prior sections is preserved. The three
themes' Meta / Synthesis Front locus links were repointed to the new
transmission-face anchor; the DR-report and OEE-references section-name pointers
were retitled; the spine's Nearby Surfaces callbacks reworded; a Grand ALife
pointer added. Carrier mechanics stay owned by
`autopoiesis-renewal/theory/core-relations.md`; novelty positioning stays in the DR
report.

**Friction:** Both ideas were correct and individually well-cited, but each had
been mechanically appended as its own `##` section — a reader saw two glued
findings rather than the single substrate property they both follow from. The
owner named it: do not mechanically apply the existing structure when it loses the
insight. Mechanical appending had hidden the thorough synthesis — the lost free-un-fakeability
assumption is the shared root; transmission / selection / inheritance are its
consequences, not three separate discoveries.

**Structural repair:** Promote the substrate property to the section frame and
demote the two former section titles to `###` face subsections (semantic content
verbatim and still addressable), so the spine reads as one argument — property →
three faces — with single-canonical-location preserved (spine owns the synthesis +
barrier framing; mechanics owned by core-relations.md; novelty owned by the DR
report). The mechanical edit-residue (an "inbound cross-references to repair"
paragraph the draft had placed inside the spine) was kept out of the canonical doc
and recorded here instead.

**What it taught:** Co-located findings that share a root are an accretion risk even
when each is individually honest and well-cited. The meta axis pays off a second
way — not by adding theory but by recognizing two recently-landed sections were one,
and unifying them without dropping a claim or overclaiming the parts. And the
honest novelty posture (cite the parts, claim the conjunction) is itself a durable
structural lesson: a substrate-recombination program's contribution is the
conjunction and the empirical demonstration, never the invention of its components.

### Jun 13 — Cultural-evolution survey resolves LIMB 2 to partially-covered; Sperber tension scoped

**Key commits:** `theory-generative-substrate` branch (LIMB 2 resolution pass).

**What changed:** A third deep-research pass (102 agents; saved at
`Paper/themes/open-ended-ecology/references/narration-collapse-culture-limb2-positioning-dr-report-2026-06-13.md`)
resolved the counterfactual-culture corollary against the cultural-evolution /
memetics literature: **partially covered**. The generic premise (inheritance
spreads independent of truth and host fitness) is canonically named — Dawkins /
Blackmore memetics, Boyd & Richerson "selfish / rogue cultural variants"
(fitness-anchored, not falsity-anchored), Acerbi 2019 / Fessler cultural-attraction
misinformation — so the spine now cites these and claims only the residual none
names: the conjunction of *stateless trace-only* inheritance + *high-fidelity
copying of false content as such*. An important honest correction surfaced and was
applied: Sperber's Cultural Attraction Theory explicitly rejects high-fidelity
copying (transmission is reconstruction), a real tension with the "high-fidelity"
wording; the spine's Inheritance face now scopes "high-fidelity" to the external
trace/text channel (Acerbi's "fidelity amplifier" regime, where Sperber's
objection is weakest) and cites Blackmore as the pro-high-fidelity anchor. The
centerpiece (LLM-ecology narration collapse) stayed unconfirmed; a dedicated
venue-by-venue negative sweep is running.

**Friction:** The corollary had been carried as a blanket "may already be named"
under-verified flag; left vague it risked either overclaiming novelty or silently
conceding it. The Sperber tension in particular would have been a clean reviewer
kill of the "high-fidelity" wording had it not been surfaced and scoped.

**Structural repair:** Update the Inheritance face and the Status / Open-Limbs
block to honest sourced states (cite the generic premise; own only the
stateless-trace-only + high-fidelity-of-false residual; scope the high-fidelity
wording to the text channel and name the Sperber tension as accounted-for). Raw
survey in OEE references; spine carries only the verdict.

**What it taught:** A live theory claim gets *stronger* when an adversarial survey
forces it to (a) cite the canonical names it was unknowingly re-deriving, (b)
shrink to the precise residual it actually owns, and (c) pre-empt the sharpest
disciplinary objection (Sperber) by scoping its own wording rather than ignoring
it. Honest positioning is a coherence gain, not a concession.

### Jun 13 — ALife / cultural-evolution survey closes one open limb, narrows the claim

**Key commits:** `theory-generative-substrate` branch (open-limb resolution pass).

**What changed:** A second 105-agent deep-research pass (saved at
`Paper/themes/open-ended-ecology/references/narration-collapse-alife-culture-positioning-dr-report-2026-06-13.md`)
covered the ALife / open-ended-evolution and cultural-evolution literature the
first survey had explicitly skipped, and resolved one of the two open limbs of the
substrate-property synthesis. LIMB 1 (functional-vs-narrative complexity) is
**prior art — cite, do not rename**: Bedau & Packard's evolutionary-activity
statistics (neutral-shadow null) and the MODES knockout "informative sites" filter
(Dolson & Ofria 2019) already name and operationalize the adaptive-vs-non-adaptive
complexity distinction. The spine's Status / Open-Limbs block was updated to cite
them and to **narrow** the project's complexity claim to the genuine residual the
filters do not handle: the adversarial case where the narrative item is itself
selected by the external scalar (not a neutral hitchhiker). Lehman & Stanley's
"Abandoning Objectives" (2011) is cited for the fixed-objective-destination idea,
with a load-bearing boundary — its *deception* pathology (goal not reached) is kept
distinct from the *satisficing* / Goodhart pathology (goal cheaply satisfied);
four claims conflating them were refuted. Two limbs stay open: the
counterfactual-culture corollary (zero confirmed claims on Boyd & Richerson /
Dawkins / Sperber — needs its own pass) and the exec-exit-0 un-fakeability stress
test; the centerpiece appears unnamed at medium confidence pending a targeted
negative search.

**Friction:** The first novelty survey had a self-declared decisive gap — it never
covered ALife / cultural-evolution — so the spine carried two limbs as "may already
be named." Leaving them as vague "may" is exactly the overclaim-or-underclaim
ambiguity the novelty posture exists to prevent.

**Structural repair:** Run the dedicated survey, then update the spine's
Open-Limbs block to honest, sourced states (RESOLVED-cite / RESOLVED-with-boundary
/ STILL-OPEN / CENTERPIECE-medium), and narrow the complexity claim to its real
residual. The raw survey lives in OEE references; the spine carries only the
verdict.

**What it taught:** Honest positioning is iterative: the first survey named what it
could and flagged its own gap; the second closed half the gap and *shrank* a claim
rather than growing one. A research program improves as much by correctly citing
prior art and narrowing to its real residual as by adding new theory.

### Jun 13 — Dedicated negative sweep: centerpiece confirmed-unnamed (swept venues)

**Key commits:** `theory-generative-substrate` branch (centerpiece negative-sweep pass).

**What changed:** A fourth deep-research pass (102 agents; saved at
`Paper/themes/open-ended-ecology/references/narration-collapse-centerpiece-negative-sweep-dr-report-2026-06-13.md`)
ran the dedicated 2023-2026 venue-by-venue negative search the prior passes had
deferred, and upgraded the centerpiece from "appears unnamed / medium" to
**confirmed-unnamed across swept venues (high confidence)**: no source names the
narration-collapse failure mode, a functional-vs-narrative complexity distinction
for generative-agent populations, or the un-fakeable functional-dependency-edge
cure. Near-misses were ruled out by full-text extraction — roadmaps calling for
exactly this system without naming the failure or cure ("Static Sandboxes Are
Inadequate" 2510.13982; Hammond et al. "Multi-Agent Risks" 2502.14143);
transmission-chain studies with no selection signal; self-rewarding co-evolution
that scores text (MAE 2510.23595, the very failure the cure guards against); exit-0
SE gating used as single-pipeline correctness (AgentForge). The cure's precise
unnamed contribution was sharpened: each ingredient exists separately (exit-0
inter-agent execution; persistent-artifact ecologies under selection) but no source
unifies execution-grounded inter-agent dependency WITH a population selection
signal. Acerbi 2016 was found to name + quantify the digital "fidelity amplifier"
regime, so the fidelity half of the culture residual is pre-empted (now cited); the
spine's Status and Inheritance face were updated accordingly. The honesty bound is
recorded explicitly: "no match found after a dedicated sweep", NOT "provably first"
— hard-to-index venues (GECCO OEE / ALIFE proceedings / AAMAS) are the residual
coverage gap, and the stateless-trace-only naming follow-up stays unchecked.

**Friction:** The centerpiece had carried a medium-confidence "appears unnamed"
inferred from absence within earlier verified sets — a soft state that could not
honestly support a novelty claim either way.

**Structural repair:** Run the dedicated sweep, then record the confirmed-unnamed
verdict WITH its explicit honesty bound (swept-venue, not provably-first) and the
sharpened cure-novelty statement; cite Acerbi 2016 for the pre-empted fidelity half.
Spine carries the verdict; the four DR reports in OEE references carry the evidence.

**What it taught:** A negative novelty verdict is only worth as much as the sweep
behind it, and it must be stated with its own scope bound. Four adversarial surveys
converged on a disciplined posture: cite Goodhart / spec-gaming / Bedau-MODES /
Lehman-Stanley / Dawkins-Blackmore / Acerbi for the components; claim only the
conjunction and the empirical demonstration; and keep two honest residuals open
(trace-only naming; the exec-exit-0 adversarial stress test, which is an experiment,
not a literature question).

### Jun 13 — Self-correction: de-gate the theory control plane (exploration invites, dossier audits)

**Key commits:** `theory-generative-substrate` branch (control-plane de-gate pass).
Human-directed, after the owner named a real failure: the assistant, asked to
*brainstorm* novelty, ran an adversarial novelty audit and deflated every direction
— "don't treat the garden as a courtroom."

**What changed:** A developmental-lens audit (9-agent workflow, both a rigor-guard
and a generativity-check) found that the same machinery added earlier today had
become an attractor that snuffs live sprouts. The unifying defect: ONE growth
surface set was being asked to be both the exploration medium AND the
publication-defense audit, so dossier-layer controls (spine-alignment proof,
maturity gating, cite-don't-claim, the anti-decoration yardstick, named-claim
closure) were stationed at the **exploration intake** instead of the **promotion
boundary**. The repair RE-SEPARATES two regimes without weakening either:
- AGENTS.md Theory pass: the spine-first read is now **orientation, not a gate**;
  seeding work routes through a LIGHT contract (name the signal + sibling, mark
  `SEED`, no proof) and the spine/maturity/novelty gates fire only on *promoting*
  work crossing into a dossier or a spine update; increments self-classify
  `SEED` / `PROMOTING` in the closure line.
- spine: the reflexive "name the future behavior or do not add it" yardstick is
  scoped to durable *infrastructure*, explicitly NOT to seeds (which are material);
  the Novelty Posture is scoped to dossier-crossing / spine-update work, not
  exploration; `Direction Versus Claim` is named as the regime boundary.
- themes/README: a new optional `## Live Seeds` exploration surface (speculative,
  un-linted, no proof) distinct from the typed Sibling Edges; neighbors framed as
  lenses that suggest experiments, not fences.
- CLOSURE.md: an exploratory-absorption path (register a seed fast, integrate
  later) so suggestive runs are neither overclaimed nor deleted.

**Friction:** the structures added earlier today (the Meta/Synthesis Front
deferral gate, the in-spine novelty posture / STILL-OPEN stamps) gave the entering
agent a downhill pull toward audit/deflation — the repo's own structural vector
capturing its own meta-analysis. Generative vocab already outweighed prohibition
(92:41), so this was about gate **placement**, not tone.

**Structural repair:** relocate every gate from intake to the promotion boundary,
add the protected seed surface + exploratory closure path, and keep all rigor
preserved at the dossier layer (anti-overclaim, single-canonical-location, the
demotion rule, evidence traceability — all named and kept, just downstream of the
boundary). The new methodology `Paper/philosophy/structural-pull.md` records WHY
(the environment has a behavioral vector; self-audit must account for it).

**What it taught:** the repository's structure is itself an effector with a vector
— it shapes the agents that maintain it, and a meta-analysis run inside it inherits
that pull. Building discipline surfaces is setting a default behavior for every
future agent; an invitation and a gate look alike on the page but pull opposite
behaviors. Order belongs at the gate between exploration and dossier; creativity is
freed everywhere upstream of it.

### Jun 13 — Promote structural pull from acting-principle to a growable theme theory

**Key commits:** `theory-generative-substrate` branch (structural-pull theory-upgrade pass).
Human-directed.

**What changed:** Structural pull (the environment has a behavioral vector; the
repository's structure acts on the agents who maintain it) had been recorded only
as a cross-layer acting principle in `Paper/philosophy/structural-pull.md`. The
owner asked to upgrade it into a theory under a theme so it can be enriched and
grown like any other theory seed, not stay a one-off repo principle. Promoted into
the `repository-self-model` theme: a new Axis 7 ("Structure has a behavioral vector
— structural pull", the dual of the founding files-as-memory/agents-as-effectors
axis), a new growable theory file `theory/structural-pull.md` (the object, a
neighbors-as-lenses backbone — Gibson affordances, stigmergy, nudge, situated /
distributed cognition, niche construction, second-order cybernetics — a generative
`shares-object` sibling edge to open-ended-ecology's affordance-rich null ecology
at the runtime locus, five Live Seeds, and a seed-stage frontier), and a
Progressive Disclosure layer entry. The philosophy doc is reframed as the
cross-layer self-audit *methodology face* and points to the theme as the theory's
growable home. Single-canonical-location preserved: the theme owns the theory
(scientific object), philosophy owns the self-audit application; they cross-link.

**Friction:** an insight worth growing was stranded as a flat acting principle with
no theme home, no neighbors backbone, and no seed surface — exactly the
"inert-because-uncoupled" state the program studies. As a philosophy note it could
be cited; it could not *grow*.

**Structural repair:** give it a theme owner (an axis + a Layer-2/3 theory file
with Live Seeds and a neighbors backbone) so future agents can plant seeds against
it and route a bounded slice into a dossier when ready; keep the methodology face
in philosophy for self-audit passes.

**What it taught:** the difference between a principle and a theory, in this repo,
is whether it has a growable home — a theme node with seeds, neighbors-as-lenses,
and a frontier. Promoting structural pull is also the first theory grown entirely
under the just-de-gated regime: planted as seeds and a neighbors backbone, with no
novelty audit at intake.

### Jun 14 — Make exploration reachable at the router + add the Cross-Reading divergence engine

**Key commits:** `theory-generative-substrate` branch (router-signage + cross-reading pass).
Human-directed.

**What changed:** Two moves. (1) A structural-pull audit of the root + Paper/ control
plane (with a reflexive-capture guard) found the structure pulls correctly toward
maintenance + publication-rigor and that those pulls are healthy and should stay —
the ONE real, mass-independent defect was that the generative entry (incubator +
themes Live Seeds) existed and was well-built but was INVISIBLE at the routers. Fix
(additive, no new surfaces, no de-gating of healthy rigor): added an "I have a
spark" row to Paper/README's Start-By-Question, a capture-first lead to Brainstorm
Mode (incubator was step 9/9), and a one-hop root README → program-map / evidence-map
/ incubator pointer. (2) Added the **Cross-Reading divergence engine** to the
themes/README node-local contract and the AGENTS Theory pass: borrowed lenses are
first-class and reusable (name the lens a seed came from; note where else it points),
and the core generative move is to re-read a theme / claim / **experiment** through a
lens or another theme's theory it did not come from, recording what falls out as a
`SEED` (theory↔theory collision = cross-reading two of our own theories). Dogfooded
with a cross-read seed in structural-pull (P11 Digital Red Queen re-read through
niche-construction theory). The loop: spark → seed (name the lens) → cross-read →
new seed.

**Friction:** the owner wanted the theory layer to have higher *inspiratory* and
*divergent* momentum — agents should not only diverge but record the neighbor
theories they borrow, collide theories, and re-read experiments through different
lenses. The de-gated structure supported ~half of that (Live Seeds = divergence;
neighbors-as-lenses = lens awareness) but had no downhill move for recording lenses
as reusable, colliding theories, or lens-driven experiment re-reading.

**Structural repair:** the reflexive-capture guard caught the audit over-applying the
de-gate frame (mass-as-gate, router-blindness, factual inflation), so the root/Paper
rigor was left intact and only the router signage was fixed; the new generativity was
added as a *named move* (Cross-Reading) reusing the existing Live Seeds surface rather
than proliferating new control-plane docs (honoring Boundedness).

**What it taught:** generativity is increased by making good moves *reachable and
named*, not by removing discipline — the audit's own honest output (after the guard)
was "extend the de-gate exactly one step to the router, and add an invitation, and
touch nothing healthy." And this very pass is another structural-pull / Reflexive
Capture datum: the audit instrument captured the analysis until a counter-guard
named the pull.

### Jun 14 — Structural pull deepened: first experiment, neighbor-lens deep research, first case study, generative-output genre

**Key commits:** `theory-structural-pull-neighbor-lenses` (continuation closure).
Human-directed.

**What changed:** Structural pull (Axis 7) was deepened from a freshly promoted seed
into a growing theory with experiment evidence, a borrowed-lens backbone, a
case-study series, and its first avoidance counter.

- *Experiment evidence.* The blind drop-in experiment was recorded as a
  `## Session evidence` section in `theory/structural-pull.md`: the vector **splits in
  two** — a strong *navigational* vector (agents reliably reach a de-gated invitation)
  and a weak *generative* one (they propose but rarely land; a no-task organism drifts
  to the maintenance default; "collisions" reproduced a pre-existing seed). Reflexive
  Capture was observed twice (organism self-reports diverged from artifacts; first-pass
  analysis graded proposals as products). New Live Seeds: measure reach-rate vs
  land-rate; activation energy; the capture signature with two instances.
- *Neighbor-lens deep research.* A 15-agent **generative** sweep (neighbors-as-lenses,
  NOT a novelty audit; a skeptical fit-critic counterbalanced the generative frame) →
  `references/structural-pull-related-theory-dr-report-2026-06-14.md`, organized by the
  theory's three faces. The fit-critic correctly flagged the Gödel/Tarski cluster as
  over-stretched (analogy only, never a proof). Five richest directions planted as Live
  Seeds: foraging / Marginal Value Theorem (reach-without-land = optimal foraging,
  leaving at the habitat average = maintenance), behavior settings + undermanning,
  generative entrenchment, double-hermeneutic + standpoint inversion, and three
  competing "ignitions" of generation (learning-progress inverted-U /
  contradiction-as-motor / novelty-search-vs-archive).
- *First case study.* `theory/cases/01-the-audit-attractor.md` — a generative task
  captured into audit-shaped output **three times in one session**, the third minutes
  after the analyst authored the detector for it. Establishes the `theory/cases/`
  instance series.
- *First avoidance counter (option ①).* A "shape of a generative pass" genre added to
  the themes/README Cross-Reading engine: a typified generative form (idea → where it
  pulls thinking → collision → open question) that explicitly is NOT the
  asymmetry-warning / "what NOT to borrow" / "don't over-claim" shape, plus an
  output-shape signature check. Designed so its downhill is "produce ideas in shape,"
  not "watch yourself for bias."

**Friction:** the audit attractor captured generative work twice more this session —
most sharply when a request to *learn neighboring theories for perspective* was
returned as an asymmetry-warning-first, firewall-centred, "don't over-claim"
presentation, until the owner named it ("我什么都没 claim 啊" — "I'm not claiming anything"). The same structural-pull /
Reflexive Capture vector the de-gate pass first caught, now recurring at the
*presentation* seam.

**Structural repair:** record the recurrence as the durable first case study, and
counter it with a positive-valence genre at the generative surface rather than a
watch-yourself rule in AGENTS.md — which would both rely on captured introspection and
add to the discipline-surface density that is part of the pull. Neighbor theory
deepens the theory but stays in `references/`; spine and dossier rigor untouched.

**What it taught:** the theory's sharpest claim keeps being confirmed by the work that
studies it — author-knowledge of the detector confers no immunity, and the reliable
detector is external every time. The right counter to a behavioral vector is another
vector whose *downhill is the target behavior* (a genre), not an injunction to be
vigilant. And matching the instrument's valence to the task (generative vs audit) is
the standing discipline the case study exists to enforce.

### Jun 14 — Extension lenses + Gemini recruited as a theory-spark tool; theme-wide Live Seeds harvest

**Key commits:** `theory-structural-pull-neighbor-lenses` (continuation closure).
Human-directed.

**What changed:** Three moves, all additive.

- *Depth ladder + the tectonic correction.* An external theory review surfaced four
  deeper neighbors (actor-network theory, media ecology, enactive/structural coupling,
  generalized Conway's law); these were added to `theory/structural-pull.md` as a
  `## The depth ladder` synthesis sorting the neighbors by *how deep the pull reaches*
  (action → trace → network → form → trigger → reasoning architecture) and converging on
  "the agent is constituted by the structure." A second pass (Gemini via `agy`) caught
  that the convergence over-reached toward *total* constitution and bounded it with the
  **tectonic correction**: the agent arrives carrying its own rigid pre-training plate,
  so structural pull is a *collision of two structures* (repo vector × base-model
  attractor), `VEC` is a fault line not a sum, generation and friction live at the seam,
  and Reflexive Capture is partial *because* the pre-training plate is unassimilable. The
  bound was landed as a **generative** move (it opens a frontier — locate the fault), not
  a retraction.
- *Gemini recruited as a first-class tool.* `agy` (Google Gemini CLI) was recorded in
  the control plane as a theory-spark instrument: a new
  `development/contracts/theory-spark-query.md` (sibling of `deep-research-query.md`), a
  `DEVELOPMENT.md` routing row + standing-preference bullet + `### Theory-Spark Query
  Contract` section, and a pointer in `themes/README.md`'s Cross-Reading engine. Framed
  as the external different-valence instrument structural pull recommends as a partial
  defense against Reflexive Capture (a different structure / fault line); `dr` =
  fact-finding, `agy` = generative spark.
- *Theme-wide spark harvest.* Gemini was run across all six themes (one agent each) plus
  a cross-theme collision pass; a `## Live Seeds` section was landed into every
  `THEORY.md` (29 seeds; none had one before), the sharpest cross-collision
  (structural pull = an extended phenotype with a viability function) added to
  `structural-pull.md`, and the blind-spot seeds from the meta round (Umwelt, Barad
  intra-action, foraging-bankruptcy, epistemic staging-ground). Full session preserved
  in `references/gemini-theme-spark-harvest-2026-06-14.md`.

**Friction:** two empty mediums met their own finding #2 ("variation media nobody writes
stay empty"). (1) The meta theory had been recorded with a convergence over-reach that
no inside-the-field pass had caught; an external model on a *different* fault line was
positioned to catch it — itself corroborating Reflexive Capture's "the reliable detector
is external." (2) The `## Live Seeds` surface existed in the `themes/README` contract but
was empty in all six themes — the theory layer's own variation medium, uncoupled and
therefore unwritten.

**Structural repair:** record the external-spark source as a repeatable, tooled
control-plane instrument (not chat-only); fill the empty Live Seeds medium across all
themes with sourced sparks; bound the convergence generatively rather than deleting it.

**What it taught:** an external model of a *different structure* is a usable partial
defense against Reflexive Capture precisely because it sits on a different fault line —
the same logic that makes the human owner the reliable detector. And "enrich the theory
layer," operationalized, is: fill the empty variation medium with sourced sparks, tooled
and recorded so the move is repeatable rather than a one-off.
