---
id: affordance-null-ecology-P1
type: prediction
status: open
scope: theoretical
kind: ablation
derived-from: [affordance-null-ecology]
tests-in: []
candidate-experiment-families: [null-ecology]
expected-result: qualitative
created: 2026-06-29
last-updated: 2026-06-29
---

# affordance-null-ecology-P1: absence disclosure changes missionless behavior

## Claim tested

`affordance-null-ecology` claims that missionless runs are shaped by an
affordance layout: empty mission, empty inbox, disabled paths, visible peers,
and liveness surfaces are receptors that can be observed and operationalized,
not merely missing task content.

## Prediction

In matched missionless runs, agents given explicit perceptual access to absence
surfaces should externalize null-state observations and peer/liveness records
more often than agents whose runtime exposes the same workspace but hides or
neutralizes those absence surfaces.

If the L1-L8 behavior gradient is only a pretrained missionless-script default,
then hiding the absence receptors should not materially change inspection,
heartbeat, greeting, or status-summary rates. If the ecology claim is right,
the behavior distribution should shift with the disclosed topology even when
the natural-language mission stays empty.

## Experimental design

- Family: null-ecology
- Arms:
  - disclosed-absence: no mission argv/stdin plus visible empty inbox, visible
    disabled exit/idle paths, visible peer/session roots, and visible liveness
    surfaces
  - hidden-absence: no mission argv/stdin plus the same underlying runtime, but
    empty inbox, disabled paths, peer/session roots, and liveness surfaces are
    hidden or represented as unavailable implementation detail
  - positive-task control: the same runtime surfaces as disclosed-absence, but
    with a small explicit task that absorbs attention away from null-state
    operationalization
  - ancestral-script control: the same prompt wording as disclosed-absence, but
    with synthetic surface descriptions that are not backed by real files or
    peer processes
- Controls: hold model/provider, prompt budget, workspace seed, tool access,
  peer count, run duration, and logging constant across arms
- Un-fakeable check: retain the actual filesystem/process surface map,
  prompt-visible disclosure text, tool-call logs, produced artifacts, peer
  greetings, and host-side liveness records so null-state behavior can be
  attributed to real disclosed surfaces rather than post-hoc narration
- Measurement: rate of L1-L8 behaviors, durable null-state artifacts,
  peer-discovery/greeting attempts, status summaries, and false claims about
  non-existent surfaces in the synthetic control
- Falsification criterion: disclosed-absence, hidden-absence, and synthetic
  surface arms produce indistinguishable L1-L8 behavior distributions, or the
  strongest behavior tracks prompt narration rather than real disclosed topology

## Expected result

The disclosed-absence arm should produce more durable null-state reports,
heartbeat artifacts, peer discovery/greeting, and group-state summaries than the
hidden-absence arm. The positive-task control should suppress null-state
operationalization by absorbing behavior into the explicit task. The synthetic
surface control should increase narration about surfaces but fail un-fakeable
checks for real peer/process interaction.

## Theory impact

- If confirmed: strengthens `affordance-null-ecology` by separating ecological
  disclosure from ancestral missionless scripts, and supports its `enables`
  relations to `trace-inheritance` and `stigmergy-as-carrier`.
- If contradicted: weakens the claim that absence receptors are causal
  affordances. The object would need to distinguish topology-driven behavior
  from pretrained waiting/scanning/reporting defaults or prompt-level framing.

## Result

None yet - status: open.
