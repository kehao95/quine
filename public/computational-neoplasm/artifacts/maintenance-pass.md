# Maintenance Pass Contract

> Extracted operational contract. Activation trigger, stance, and authorization
> boundary: [`AGENTS.md`](../../AGENTS.md) § Maintenance Keyword.
> Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).

## Causal Path

| Field | Value |
|-------|-------|
| reader / effector | agents running a pure self-maintenance pass |
| behavior effect | turn the maintenance stance into a checkable pass: diagnose deeply, repair the smallest coherent surface or prune loadless machinery, validate, close |
| perception / validation | maintenance sessions drifting into patch sprees, broad redesign, or undocumented closure |
| absorption owner | this file for checklist drift; `AGENTS.md` for stance and boundary drift |

---

Named friction: the Constitution used to declare itself non-procedural while
carrying the maintenance checklist. Detailed pass criteria are procedure, not
self-model, so they live with the extracted development contracts.

## A Valid Maintenance Pass

A valid maintenance pass:

- inspects current repository state before acting, including local diffs,
  recent commits, and relevant control-plane surfaces
- investigates one repository structure, control plane, or subsystem deeply
  enough to name its current coherence state
- names the largest bounded friction or entropy source revealed by that
  investigation
- watches for authority-creep: any rule, core principles included, being
  treated as immutable rather than functional-and-revisable. Amendment friction
  tracks stability-need and destructive scope, not rank. (Stance:
  [`AGENTS.md`](../../AGENTS.md) § Descriptive vs Normative;
  [`Paper/philosophy/mutual-encoding-network.md`](../../Paper/philosophy/mutual-encoding-network.md)
  Live Seeds.)
- chooses the owning layer before editing: implementation, tests, evidence,
  status, roadmap, paper, or Constitution
- leaves unrelated Human work untouched except for explicit classification or
  handoff
- treats repair as downstream of diagnosis: prefer one small durable repair
  over broad redesign once the owning layer and direction are clear
- discusses repair direction with Human before changing policy or structure
  when the diagnosis admits multiple plausible futures
- updates navigation when a repair changes where future agents should look
- routes venue, deadline, and format-envelope questions to
  `Paper/venues/objects/<id>.md` and `Paper/venues/README.md` before consulting
  `publication-strategy.md` or `TRACK/venues/`
- checks the semantic residual that
  `scripts/check-phase-interpretation-ownership.sh` cannot see: prose that
  restates a claim its projection card assigns elsewhere, plus phase-root
  sibling docs and child experiment READMEs outside the script's v1 scope. A
  parallel fact-home is named friction; this step complements the mechanical
  header/card/index checks.
- validates the repair with the narrowest meaningful check available
- commits the repair when it forms a coherent state transition; otherwise
  reports why commit is deferred
- ends with `Closure complete` or `Closure incomplete` according to
  [`CLOSURE.md`](../../CLOSURE.md)

## Good Autonomous Targets

Broken links, stale indexes, duplicate explanations, missing ownership,
unabsorbed standing Human preferences, temporary residue that has a clear
destination, and control-plane drift visible from current repository state.

One bounded **theory-front increment** is also a good autonomous target when the
session would otherwise only tidy control planes: absorb one unclaimed run into
a `Paper/theory/objects/<id>.md` relation/evidence delta, forge or close one
typed relation, or move one `Paper/theory/gaps/<id>.md` frontier one rung. The
graph-native contract lives in [`../../Paper/theory/README.md`](../../Paper/theory/README.md)
and [`../../Paper/contracts.md`](../../Paper/contracts.md); the completed re-cut
is recorded in
[`../../Paper/theory-restructure-design.md`](../../Paper/theory-restructure-design.md).

A fourth bounded theory-front target is growing the **Meta / Synthesis Front**:
sharpen a stale **Throughline question**, correct a stale **Spine alignment**
verb ({instantiates | serves | specializes | tensions-with}), or forge an honest
`tensions-with` edge into the spine
([`../../Paper/philosophy/self-maintaining-systems.md`](../../Paper/philosophy/self-maintaining-systems.md)).
A theory pass leads with this directive axis: state the throughline and spine
alignment first, and defer depth or sideways increments that cannot name what
they serve. The Throughline question is synthesis/framework work, not an
experiment step.

## Out Of Scope

The authorization boundary is constitutional and lives in
[`AGENTS.md`](../../AGENTS.md) § Maintenance Keyword:
no destructive operations, sweeping philosophical rewrites, runtime semantic
changes, evidence deletion, or modification of unowned local diffs. If the
strongest repair requires those, stop at a named handoff.

The bounded theory-front increment above carves out **one** evidence-anchored
claim/edge/maturity delta. Re-theorizing a family, rewriting a core formula, or
promoting a claim more than one maturity rung remains Human-directed theory work.

## Prune Mode

Named friction: visible-friction maintenance can leave loadless machinery green
under presence-only sensors, while requiring prior named friction makes some
deletion repairs fail the gate they are repairing.

A prune mode is a Human-directed maintenance pass whose target is deletion. For
an existing rule, field, gate, typed vocabulary, or validator requirement, name
the future behavior that depends on it. If none can be named, the mechanism is a
deletion candidate even without prior named friction.

A valid prune pass:

- classifies each probed mechanism as load-bearing, mixed, or narrative
- names one concrete deletion and the validation that would prove it did not
  break the system
- enacts only bounded deletions inside Human scope; otherwise hands off the
  deletion proposal
- does not replace deleted machinery with a larger rule unless a future behavior
  depends on the replacement

## Relationship To Other Surfaces

- [`AGENTS.md`](../../AGENTS.md) owns the activation keyword, the
  investigation-first stance, and the red-line boundary
- [`CLOSURE.md`](../../CLOSURE.md) owns the closure audit the pass must end with
- [`structural-friction-journal.md`](./structural-friction-journal.md) owns
  when a repair earns an evolution journal entry
- the scheduled wake-model requirement and loop-closure criterion below are owned here;
  [`AGENTS.md`](../../AGENTS.md) § Automation carries the principle-grade sentences
  that point to them
- the experiment family-index is the entry point; reach experiments by
  neighborhood there, or by question through a theory object's `## Evidence`. The
  migration plan and legacy→new id-map live under
  [`../experiments-restructure/`](../experiments-restructure/).
  *(experiment-tree path redacted in the public copy)*

## Scheduled Pass (Autonomous)

Named friction: keyword/session-triggered passes miss drift in untouched or
long-quiet files; event-driven sensing has no wake model.

Current state: the scheduled-detection surface (`scripts/heartbeat-readout.sh`,
`make heartbeat`, and `development/heartbeat/readouts/`) has been **removed**.
Scheduled sensing is now a named requirement and open capability gap, not an
executable repo target; the system relies on event-driven sensing plus on-demand
maintenance passes until a replacement wake model is introduced.

Any replacement scheduled maintenance pass has narrower authority than an
interactive pass, by design:

- it **detects**: run the available sensors, including the meta-sensors that
  catch an inert or buggy validator
- it **classifies**: label each finding (sensor-bug, content-rot,
  judgment-pending, in-flight) with age and owner, per the descriptive-readout
  rule in [`../../AGENTS.md`](../../AGENTS.md) § Automation
- it **reports**: leave the classified readout where the next session will
  find it
- it **never edits**: a scheduled pass produces no repairs; an interactive
  session remains the only effector for the repairs it names

This lets the system wake to look without silently rewriting itself. The readout
is a handoff surface, not a work log; a scheduled pass closes by producing it.

In an interactive session, use `make check-control-plane` or focused validators
for validation. Do not treat a generated maintenance readout as a freshness gate;
if a replacement wake model later restores readouts, retain at most one readout
per work boundary unless the commit records a before/after sensor change.

## Loop-Closure Criterion

The maintenance loop is closed when three detection-only conditions hold
together:

1. **Wake model exists** — the repository can observe itself without an
   interactive session, through a scheduled detection-only pass. This criterion
   is currently open: the prior heartbeat target has been removed without a
   replacement.
2. **Interpretation layers are compliance-sensed** — every interpretation layer
   is watched by a compliance sensor against its own declared ownership, not by
   presence alone (see the sensor typology in
   [`../../scripts/README.md`](../../scripts/README.md)).
3. **Learning layer is closed** — every recurring pattern carries an enforcing
   check or a named reason it cannot have one, and every open-friction entry
   carries an owner and a status.

All three are detection-only: closure means the system can see its own gaps, not
that the gaps are gone.
