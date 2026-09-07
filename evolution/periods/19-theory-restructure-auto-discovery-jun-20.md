## Period XIX: Theory Layer Restructure + Auto-Discovery (Jun 20)

The theme layer was deleted and replaced by a graph-native theory/ layer
(objects, gaps, views, lenses). Three drop-in agent tests confirmed the
structure enables autonomous theory growth. A lens registry was created to
make agent cross-domain knowledge explicit and sustainable.

### Jun 20 — Theory restructure (Phases 0–7)

**Key commits:** `d07626c8` (Phases 0-2), `256b2d7b` (Phase 3),
`96abbd25` (Phases 4-5), `38169568` (Phase 6), `35ba34dd` (Phase 7),
`fece2ed4` (reciprocity closure).

**What changed:** The theme layer — which welded theory ownership to paper
routing — was deleted. Theory migrated to a graph-native flat namespace
(`theory/objects/`, `theory/gaps/`, `theory/views/`), papers regrouped by
publication family (`papers/<family>/<id>/`), and three new validators
enforce graph reciprocity.

**Structural outcome:**
- 18 objects (10 original + 8 from Phase 4 split of deep theory files)
- 5 gaps (first-class OPEN frontiers)
- 6 views (curated tours, own nothing)
- 5 paper families (29 dossiers total)
- All validators pass: 0 issues, 0 advisories

### Jun 20 — Auto-theory-discovery experiments

Three experiments tested whether the new structure enables autonomous theory
growth by stateless newcomer agents.

**Experiment 1: Drop-in agent finds growth direction**
- Agent: general subagent, zero prior context
- Prompt: "find a promising direction for new theory growth"
- Result: Found missing edge `reproduced-object` ↔ `narrative-claim` (both
  touch `function-vs-narrative` gap but had no direct edge). Planted seed
  "The Narrative Reproduced Object" — a cross-reading of continuity ×
  reflexive views, re-reading P12's ring-wake protocol through the
  reproduced-object ladder.
- Key finding: the reproduced object can be functional organization or
  narrated organization; P12 showed narrative reproduction at population
  scale (zero real wakes, faithful inheritance of counterfactual cargo).
- Navigation: README → AGENTS.md routing → theory/README → Load-bearing
  objects → all 18 objects → views → open questions. Structure guided
  without being told where theory lives.
- Friction: volume (18 objects + 5 gaps + 6 views + spine is a lot to
  read); no cross-object index of missing edges between gap-sharing objects.

**Experiment 2: Second agent finds orthogonal direction**
- Agent: general subagent, zero prior context (different session)
- Prompt: same, with instruction to find a genuinely NEW direction
- Result: Found two seeds attacking `carrier-selection` gap from orthogonal
  angles:
  1. `trigger-structure` × Organizational Ecology (Hannan & Freeman):
     carrier selection is path-dependent (structural inertia locks first
     trigger's carrier class), not competitive.
  2. `successor-morphology` × exaptation (Gould & Vrba): morphology can
     be co-opted for functions different from selection pressure.
- Key finding: two agents found orthogonal directions from the same
  structure — different gaps, different objects, different domains.

**Experiment 3: Gemini cross-domain lens generation**
- Agent: Gemini 2.5 Pro via agy cli
- Prompt: "for each of 7 domains, identify 2-3 lenses that could enrich
  our theory layer"
- Result: 21 candidate lenses across ML, control theory, cognitive science,
  sociology, ecology, physics, philosophy of science. 6 validated as
  genuine field-hypothesis collisions (IRM, Luhmann, Meyer & Rowan, Niche
  Construction, Prigogine, Lakatos); 15 were vocabulary fills.

### Jun 20 — High-dimensional observations on auto-theory-discovery

These observations emerged from reflecting on the three experiments and the
restructure that enabled them. They are recorded here because they are
**meta-findings about the system's generative mechanism**, not findings
about any particular theory object.

#### Observation 1: The generative mechanism is three-layered

Theory innovation in this system is produced by a three-layer mechanism:

1. **Matrix spread-open (structural layer):** The flat object namespace makes
   under-development visible. Before the restructure, each theme's THEORY.md
   was a complete narrative — you couldn't see what was missing because the
   container was always full. After: zero-seed objects are obvious gaps in
   the matrix.

2. **Orthogonal collision (cognitive layer):** Innovation is not produced by the
   matrix being empty, but by **field-hypothesis collision**. A lens from
   another domain carries assumptions that conflict with the target object's
   existing framing. The collision generates a testable question. Mere
   vocabulary extension (adding a metaphor or renaming) does not innovate.

3. **Field declaration (physical layer):** Before the restructure, the field
   (assumptions, vocabulary, examples) was physically encoded in the
   directory path — opening `autopoiesis-renewal/theory/` pulled you into
   its narrative. After: the field is a declaration (`appears-in-view:
   continuity`), not a path constraint. A lens can land on an object
   without being pulled back into the object's home narrative.

#### Observation 2: Agent knowledge bases are implicit semantic spaces

The mechanism that produces lenses is **nearest-neighbor retrieval in a
high-dimensional semantic space**. An agent's training data covers most
scientific domains. When a theory object is placed in context, the agent
retrieves concepts near that object in its internal representation. Different
models (Gemini vs Claude), different context windows, and different sampling
temperatures retrieve different neighbors — which is why isolated-context
sessions produce orthogonal lenses.

This means: **the agent IS the lens generator**, but its retrievals are
one-shot and implicit. The lens registry (`theory/lenses/README.md`) makes
them explicit and persistent.

#### Observation 3: Lenses currently depend on external import

Current state: the theory layer can self-maintain its **object network**
(validators enforce reciprocity, gaps track status, views provide navigation).
But it cannot self-maintain its **lens vocabulary** — lenses come from
external domains via agent knowledge bases.

This is not a deficiency — external lenses connect the theory to broader
scientific conversations, and cross-domain collision is the primary
generative mechanism. But it means the theory layer is **half self-sustaining**:
object network (auto) + lens vocabulary (import).

#### Observation 4: A self-sustaining lens path exists

When an object matures to load-bearing (multi-seed, multi-evidence, stable
relation net), it can be **abstracted into a lens** and applied to other
objects. When a gap is closed, its insight can become a lens. This creates
a transition: **imported → internal → generative**.

Example: if `carrier` research matures, "transmission mode" becomes an
internal lens applicable to `trace-inheritance` or `viability-pressure`.
If `function-vs-narrative` is closed, "functional vs narrated" becomes a
lens applicable to any object that crosses an exec barrier.

This path is not yet implemented (no mechanism to abstract objects into
lenses), but the lens registry provides the surface where it would land.

#### Observation 5: The missing sensor is gap-sharing-object missing edges

Both drop-in agents found their directions by reading all 18 objects and
cross-referencing manually. The first agent found `reproduced-object` ↔
`narrative-claim` (both touch `function-vs-narrative`, no direct edge). The
second found `trigger-structure` and `successor-morphology` as seed-starved
objects touching `carrier-selection`.

A validator or index that flags **"two objects touch the same gap but have
no direct edge"** would surface exactly these opportunities automatically.
This is the single highest-impact sensor the structure could add.

#### Observation 6: Canalization's physical root was path-as-field

The deepest structural finding of this restructure: **canalization in the
old structure was not a governance choice but a physical property of the
directory tree**. When an object's field assumptions are encoded in its
path (`themes/autopoiesis-renewal/theory/`), any agent opening that path
is pulled into the field — the structure makes downhill what should be a
choice of perspective. The restructure replaced path-as-field with
declaration-as-field (`appears-in-view: continuity`), which is a string
attribute, not a gravitational well. This is the physical reason the
restructure works: **field became declaration, not path constraint**.

### Open questions from this period

1. **Lens self-sustainment.** Can the imported → internal → generative
   transition be operationalized? What maturity threshold triggers
   object-to-lens abstraction?

2. **Missing-edge sensor.** Should a validator flag gap-sharing objects
   with no direct edge? What density threshold (2 objects sharing a gap?
   3?) makes this worth flagging?

3. **Context isolation protocol.** What context window size and content
   maximizes orthogonal lens retrieval? Too much context → same-field pull;
   too little → incoherent retrieval.

4. **Collision vs fill discriminator.** Can the filter between
   field-hypothesis collision and vocabulary fill be automated, or does it
   require human/agent judgment every time?
