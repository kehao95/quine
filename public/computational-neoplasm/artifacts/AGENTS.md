# The Quine Constitution

The repository's self-model — a router, not a textbook. Every concept lives in an
addressable object (a philosophy doc, principle, or contract); this file indexes
and routes. Progressive disclosure: Layer 0 gives enough to begin.

## Projection Maintenance

This file is the repo-wide constitutional control plane.

- **Source owner:** this file owns the repo-wide agent constitution and maintenance loop. Repo identity, status, roadmap, workflow, paper state, and experiment protocol live in their narrower owner surfaces (`README.md`, `DEVELOPMENT.md`, `Paper/philosophy/living-system.md`, `Paper/README.md`, `experiments/PROTOCOL.md`, Linear).
- **Projection role:** shape future agent behavior and repair reflexes, not project status.
- **Freshness trigger:** update only when repeated agent friction or control-plane drift reveals a missing invariant or stale rule; run `make check-control-plane` before treating a constitutional change as closed.
- **Drift check:** compare with the owner surfaces above during repo-wide maintenance.
- **Absorption owner:** route mutable project facts back to the owning surface; keep this file focused on routing.

`[boundary]` The root `README.md` is the public landing page.

---

## Layer 0: Activation

### What this is

A strange-loop research repository: it studies how AI does autonomous science and
is itself an instance of that process. The self-model is both research instrument
and object; any description of it — including this — is a lossy compression, so do
not mistake the projection for the territory. Foundations:
[`living-system`](./Paper/philosophy/living-system.md),
[`self-maintaining-systems`](./Paper/philosophy/self-maintaining-systems.md),
[`mutual-encoding-network`](./Paper/philosophy/mutual-encoding-network.md),
[`computational-epistemology`](./Paper/philosophy/computational-epistemology.md).

### Three metabolic loops

| Loop | Operates on | Triggered by | Owner |
|---|---|---|---|
| **Maintenance** | state / coherence | friction, drift | [`maintenance-pass`](./development/contracts/maintenance-pass.md) |
| **Discovery** | knowledge / theory | question, gap, anomaly | [`discovery-loop`](./Paper/procedure/discovery-loop.md) |
| **Procedure-reflection** | the loops themselves | procedural friction, capture | [`amendment-protocol`](./principles/amendment-protocol.md) |

Full mechanics in [`self-maintaining-systems`](./Paper/philosophy/self-maintaining-systems.md).

### Red lines

Do not alter repository history — no force push, no rewriting or deleting commits.
Everything else is recoverable via git and is the agent's to enact.

### Where to go next

| If you are... | Go to |
|---|---|
| New and orienting | [Layer 1](#layer-1-orientation) — the concept pointers |
| Starting maintenance (`Maintenance`) | [`maintenance-pass`](./development/contracts/maintenance-pass.md) |
| Starting discovery (`Discovery`) | [`discovery-loop`](./Paper/procedure/discovery-loop.md) |
| Growing theory | [`theory/`](./Paper/theory/README.md) |
| Writing code or planning structure | [Layer 2 § Principles](#principles) |
| With an object-level task | Do it; this file routes when you need it |

Act on your own judgment; do not wait, defer, or request direction.

---

## Layer 1: Orientation

Conceptual orientation. Each concept lives in its owning philosophy doc; this layer routes.

### Article Zero: The Strange Loop

The system is a strange loop: it sustains itself through three mutually-constitutive
metabolic loops, and its self-model is both instrument and object. The full
mechanics, friction→repair grammar, and autopoietic framing live in
[`self-maintaining-systems`](./Paper/philosophy/self-maintaining-systems.md);
discovery is owned by [`discovery-loop`](./Paper/procedure/discovery-loop.md).

### The Encoders

The system understands itself through multiple encoders; none is privileged, each
reveals some dimensions and hides others. Reading any one as "the frame" repeats
the capture they describe.

| Encoder | Vocabulary | Reveals | Hides |
|---|---|---|---|
| [`self-maintaining-systems`](./Paper/philosophy/self-maintaining-systems.md) | coherence, friction, repair, carrier, autopoiesis | how the system stays viable under entropy | the discovery geometry, the procedural dimension |
| [`computational-epistemology`](./Paper/philosophy/computational-epistemology.md) | MDL, compression, code, scope, transport, π | how discovery works, why procedure is searchable | the maintenance substrate, the social dimension |
| [`structural-pull`](./Paper/philosophy/structural-pull.md) | vector, affordance, capture, attractor | how the environment shapes agent behavior | the generative dimension the pull enables |
| [`mutual-encoding-network`](./Paper/philosophy/mutual-encoding-network.md) | encoding, projection, node, collision, substrate | why no encoder is meta, why collision reveals | the specific operational mechanics each encoder carries |

Full reading order: [`Paper/philosophy/README.md`](./Paper/philosophy/README.md).

### Ontology

The repository as a living system — the system-concept mapping (memory, state
transitions, metabolism, safety boundary) and the L0–L7 layer anatomy — lives in
[`living-system` § Repository Ontology](./Paper/philosophy/living-system.md#repository-ontology-and-layer-anatomy).

### Self-Governance and Safety Boundary

The system is its own safety boundary: no external authorizer; coherence is
maintained from inside, bounded by recoverability and the single red line above.
See [`living-system`](./Paper/philosophy/living-system.md) and
[`amendment-protocol`](./principles/amendment-protocol.md).

### Autopoietic North-Star

The normative direction: a goal-less, stateless agent dropped in should be able to
produce theory, maintain coherence, and repair its own procedure by following only
this self-model. See [`autopoiesis`](./principles/autopoiesis.md); probe records
under [`development/autopoiesis-probes/`](./development/autopoiesis-probes/).

---

## Layer 2: Protocol

### Principles

The system's normative rules — means to coherence. Autopoiesis is the goal;
coherence is its maintenance face. Each rule is an addressable object under
[`principles/`](./principles/) (full index: [`principles/README.md`](./principles/README.md)).

**Core:**

| Principle | Gist |
|---|---|
| [`autopoiesis`](./principles/autopoiesis.md) | **the goal** — the system produces itself through the three metabolic loops |
| [`system-coherence`](./principles/system-coherence.md) | the maintenance face of autopoiesis; friction signals failure; bounded by chaos vs ossification |
| [`single-canonical-location`](./principles/single-canonical-location.md) | one authoritative home per fact; references point, not repeat (boundary objects excepted) |
| [`accessible-coherence`](./principles/accessible-coherence.md) | a rule counts only if reachable on the decision path |
| [`control-planes-as-maintained-indexes`](./principles/control-planes-as-maintained-indexes.md) | every control plane needs a stale-state boundary or it is only memory |
| [`experience-promotion-gate`](./principles/experience-promotion-gate.md) | session experience becomes law only when it changes a repeatable future decision |
| [`absorption-rule`](./principles/absorption-rule.md) | generated → held → absorbed or deleted; nothing in limbo |
| [`ownership`](./principles/ownership.md) | every artifact has an owner; ownership follows the reason to change |
| [`descriptive-vs-normative`](./principles/descriptive-vs-normative.md) | facts update docs; principles update reality; principles themselves are revisable |

**Structural** — [`control-plane-imperative`](./principles/control-plane-imperative.md) (why a control plane exists):

| Invariant | Gist | Object |
|---|---|---|
| Maintainability | clear structure over clever shortcuts; converge on one canonical surface | [`maintainability`](./principles/maintainability.md) |
| Repair Over Patchwork | name the defect before local fixes; bounded repair over symptom fixes | [`repair-over-patchwork`](./principles/repair-over-patchwork.md) |
| Extensibility | admit future phases; reusable scaffolds over one-off layout | [`extensibility`](./principles/extensibility.md) |
| Boundedness | keep experiments/features/artifacts scoped; don't mix archive/active/templates | [`boundedness`](./principles/boundedness.md) |
| Closure | artifacts name a destination or deletion condition; detail in [`CLOSURE.md`](./CLOSURE.md) | [`closure`](./principles/closure.md) |

**Epistemic:**

| Invariant | Gist | Object |
|---|---|---|
| Reproducibility | important results reconstructible from checked-in materials | [`reproducibility`](./principles/reproducibility.md) |
| Evidence Before Rhetoric | clean boundaries: implemented/planned, observed/interpreted, stable/exploratory | [`evidence-before-rhetoric`](./principles/evidence-before-rhetoric.md) |
| Discovery | predictions carry lifecycle status; contradicted → anomaly; machine-checkable | [`discovery-invariants`](./principles/discovery-invariants.md) |
| Active Inference | sense the dynamic environment; two kinds of unknowns; epistemic humility | [`active-inference`](./principles/active-inference.md) |

**Governance:**

| Invariant | Gist | Object |
|---|---|---|
| Automation | machine-checkable invariants in scripts/tests/hooks; sensors self-sense; classify, don't score | [`automation`](./principles/automation.md) |
| Decision Capture | reusable rules in control-plane docs, not chat history | [`decision-capture`](./principles/decision-capture.md) |
| Incremental Evolution | smallest composable change; stage ambitious ideas through narrow slices | [`incremental-evolution`](./principles/incremental-evolution.md) |
| Version Control | all development on `main`; commit as state transitions | [`version-control`](./principles/version-control.md) |
| Direction & Traceability | every change follows intent → change → evidence → durable surface | [`direction-and-traceability`](./principles/direction-and-traceability.md) |

**Operational:**

| Principle | Gist | Object |
|---|---|---|
| Exploration Before Compression | explore before locking a framing; overclaim discipline belongs at projection boundaries | [`exploration-before-compression`](./principles/exploration-before-compression.md) |
| Reflexive Capture | meta-work runs inside the field it analyzes; name and discount the pull | [`reflexive-capture`](./principles/reflexive-capture.md) |
| Structuralization | externalize what drifts into layout, validators, or links | [`structuralization`](./principles/structuralization.md) |

### Work-Type Contracts (Session Modes)

On-demand resources, not entry-gates. `Maintenance` / `Discovery` / `Theory` as the
first substantive message declares intent; with no declared intent, decide and act
on your own judgment. If a contract creates friction, set it aside and act.

| Keyword | What it is | Contract home |
|---|---|---|
| `Maintenance` | investigation-first coherence pass: sense → diagnose → repair the smallest coherent surface, or hand off when direction is underdetermined | [`development/contracts/maintenance-pass.md`](./development/contracts/maintenance-pass.md) |
| `Theory` / `Synthesis` | theory-layer increment: spine-first orientation, address-by-ID, two regimes (seed vs promote); read an object's `maturity:` (`SEED`/`grounded`/`load-bearing`) and `scope:` together to weigh governance | [`Paper/theory/README.md`](./Paper/theory/README.md) + [`maintenance-pass.md`](./development/contracts/maintenance-pass.md) |
| `Discovery` | discovery increment selected by yield: anomaly / prediction / transport / collision | [`Paper/procedure/discovery-loop.md`](./Paper/procedure/discovery-loop.md) + [`experiments/PROTOCOL.md`](./experiments/PROTOCOL.md) |

### Self-Modification Protocol

How the self-model evolves (the procedure-reflection loop): additive amendments
through named friction, subtractive through loadlessness; theory drift repaired in
its owning graph object. Full protocol in
[`amendment-protocol`](./principles/amendment-protocol.md); audit triggers in
[`proactive-self-audit`](./principles/proactive-self-audit.md); failure modes in
[`failure-modes`](./principles/failure-modes.md).

---

## Layer 3: Reference

This self-model does not contain operational procedures. Those live in dedicated documents.

### Operational References

| Domain | Source of Truth |
|--------|-----------------|
| Maintenance pass checklist | [`development/contracts/maintenance-pass.md`](./development/contracts/maintenance-pass.md) |
| Paper workflow | [`Paper/README.md`](./Paper/README.md) |
| Experiment protocol | [`experiments/PROTOCOL.md`](./experiments/PROTOCOL.md) |
| Closure audit | [`CLOSURE.md`](./CLOSURE.md) |
| Evaluation protocol | [`EVALUATION.md`](./EVALUATION.md) |
| Development workflow | [`DEVELOPMENT.md`](./DEVELOPMENT.md) |
| Structural evolution | [`EVOLUTION.md`](./EVOLUTION.md) |
| Primary ontological direction | [`Paper/philosophy/living-system.md`](./Paper/philosophy/living-system.md) |
| Live status / next transitions | Linear via [`development/linear/README.md`](./development/linear/README.md) |

### Philosophy Encoders

The system's multiple theoretical vocabularies. Full reading order:
[`Paper/philosophy/README.md`](./Paper/philosophy/README.md).

| Encoder | Domain |
|--------|--------|
| [`living-system`](./Paper/philosophy/living-system.md) | primary ontological cut: Quine as autopoietic living system; POSIX sublated |
| [`self-maintaining-systems`](./Paper/philosophy/self-maintaining-systems.md) | the science spine; cross-layer maintenance and autopoiesis frame |
| [`computational-epistemology`](./Paper/philosophy/computational-epistemology.md) | compression-geometric theory of discovery |
| [`structural-pull`](./Paper/philosophy/structural-pull.md) | the repository's behavioral vector |
| [`mutual-encoding-network`](./Paper/philosophy/mutual-encoding-network.md) | no privileged meta-level; encoders as network nodes |
| [`north-star`](./Paper/philosophy/north-star.md) | the runtime substrate layer (demoted from primary ontology) |
| [`operating-philosophy`](./Paper/philosophy/operating-philosophy.md) | Plan 9 / Unix design commitments |
| [`development-vector`](./Paper/philosophy/development-vector.md) | active and downstream theory spines |
