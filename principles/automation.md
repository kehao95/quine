---
id: automation
type: principle
gist: automate concrete recurring checks and verify their real trigger and consumer paths
status: canonical
---

# automation

Automate a recurring contract when a reliable observable and a useful consumer
exist. Reuse existing tools before adding a gate. A conceptual quality judgment
without such an observable remains review guidance.

[Sensor types](../scripts/README.md#sensor-typology) answer different questions:
presence, compliance, synchronization, and inventory. They are not an ascending
scale of scientific validity. A registered hook can still miss its intended
files; exercise actual routing and meaningful failures when changing that path.

Event-driven checks and scheduled reviews address different failure modes.
Add a schedule when drift between events creates a concrete need, with an owner
and execution mechanism. A documented cadence is not an installed scheduler.

Report findings with enough context to repair the cause. Counts can describe a
population; they should not become a universal target that rewards suppressing
alarms or producing mechanically countable research artifacts.
