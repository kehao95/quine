# Failure-Mode Objects

Addressable first-class objects for **observed behavioral failure modes** of the
repository-as-system — typically how the substrate shapes agent behavior in
ways that waste energy or squeeze generative action. Each failure mode is one
file, deposited so it can be searched over rather than lost to chat or to prose
inside a consolidated log.

This is the addressable counterpart to the friction surfaces that hold failure
as inline prose:

- [`../../development/autopoiesis-probes/FRICTION.md`](../../development/autopoiesis-probes/FRICTION.md) — cycle-attached probe friction
- [`../../EVOLUTION.md`](../../EVOLUTION.md) § Open Friction — inline `[owner; status]` entries
- [`../patterns.md`](../patterns.md) — cross-period meta-pattern taxonomy (descriptive synthesis rows)

A friction pattern **graduates to a failure-mode object** here when (a) it
recurs across cycles and (b) it carries a candidate repair worth depositing as
addressable metadata. One-off or cycle-local friction stays in `FRICTION.md`.

## Why "deposited, not applied"

Some failure modes are deliberately left with `applied: false` even when a
candidate solution is known. This serves the
[`Autopoietic North-Star`](../../AGENTS.md): a deposited-but-unresolved failure
mode is a discovery target a future goal-less agent can find, weigh, and
attempt to repair — a test of whether the system can self-repair its own
defects, not just execute given fixes. This is the intentional counterpart of
the `Externalize-without-reach (dormant rule)` pattern in
[`../patterns.md`](../patterns.md): there the dormancy is accidental; here it
is the instrument.

## Object schema

Frontmatter:

```yaml
---
id: <slug>
type: failure-mode
status: open | deposited | repaired   # open = no candidate yet; deposited = candidate recorded; repaired = wired in
applied: true | false                 # whether the candidate-solution is enacted in the substrate
class: <optional grouping, e.g. structural-pull>
symptom: <one-line>
evidence: [<cycle or commit ids>]
candidate-solution: <one-line>
owner: lab
---
```

Body sections: `## Symptom`, `## Evidence`, `## Root cause`, `## Candidate
solution (deposited — applied: <bool>)`, `## Why not applied yet` (when
`applied: false`), `## Theory links`.

A `status: deposited` + `applied: true` pair is incoherent — applying the
solution should advance `status` to `repaired` in the same change.

## Lifecycle

- **deposit** — create the object with the observation + evidence; set
  `status: deposited`, `applied: false` if a candidate solution is known but
  withheld
- **attempt** — a later agent (often a goal-less probe) finds the object,
  weighs the candidate, and either enacts the repair or records why it did not
- **close** — on enactment, set `applied: true` and `status: repaired`, and
  point to the substrate surface that now carries the repair

When a failure-mode object is repaired, do not delete it: it remains evidence
of the pattern and a check against recurrence. Retire its `patterns.md`
reverse-link row only if the pattern stops recurring.

## Discovery path

A fresh agent following only the self-model reaches these objects via:
`AGENTS.md` → [`../../EVOLUTION.md`](../../EVOLUTION.md) →
[`../README.md`](../README.md) Surfaces → this README → objects; or via the
reverse-link row in [`../patterns.md`](../patterns.md).

## Objects

| ID | Status | Applied | Symptom |
|----|--------|---------|---------|
| [`validation-attraction.md`](./validation-attraction.md) | deposited | false | goal-less agents repeatedly run repo-wide validators as the default action, burning time/tokens without generative output |
| [`shared-index-state-contamination.md`](./shared-index-state-contamination.md) | repaired | true | concurrent agents share a moving worktree and index, destabilizing ownership, provenance, and commit identity |
| [`concurrent-receipt-topology-race.md`](./concurrent-receipt-topology-race.md) | repaired | true | concurrent receipt machinery becomes its own shared-state race and distorts the evidence topology |
| [`ownership-pressure-induced-noop.md`](./ownership-pressure-induced-noop.md) | repaired | true | dirty shared ownership suppresses edits and is misread as absent agency |
| [`probe-route-capture.md`](./probe-route-capture.md) | repaired | true | visible probe apparatus redirects missionless agents away from native cold-start selection |
| [`proxy-enforcement-collapse.md`](./proxy-enforcement-collapse.md) | repaired | true | a blocking check verifies token/shape presence instead of ground truth; the population converges on template replication and the check manufactures false confidence |
| [`sh-output-explosion.md`](./sh-output-explosion.md) | deposited | false | large sh output truncates and is destroyed, forcing re-runs that explode the context and kill the incarnation (~5min/turn, memory death) |
| [`sensor-binding-severance.md`](./sensor-binding-severance.md) | deposited | false | a registered blocking validator goes silent when a restructure moves the surfaces it guards; its filter still matches something by accident, so audits stay green while nothing is checked |
