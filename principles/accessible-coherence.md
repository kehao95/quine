---
id: accessible-coherence
type: principle
gist: a rule only counts if reachable on the decision path; intent must route to the concrete protocol
status: canonical
---

# accessible-coherence

A rule only contributes to coherence if it is reachable on the path where
decisions are made. High-level intent must route to the concrete protocol; if an
agent has to remember the rule instead of following the index, the system is
only nominally coherent and has leaked entropy at the control-plane boundary.

A container that two or more sessions will edit should carry its local contract
at the entrance; a rule that lives only in a distant central doc predictably
fails at write time. The experiments analysis job-map is the worked example.

## Failure mode its absence produces

Nominal coherence: a rule exists in the self-model but is not reachable where
the decision is made, so it is violated in practice. Presence-checks stay green
because they verify the rule's existence, not its reach. This is the
"Declared ownership without compliance pressure" pattern in
[`evolution/patterns.md`](../evolution/patterns.md).

## Related

- [`control-planes-as-maintained-indexes.md`](./control-planes-as-maintained-indexes.md) — the stale-state boundary is how a control plane stays reachable-as-truth
- [`experiments/PROTOCOL.md`](../experiments/PROTOCOL.md) § Phase Interpretation Layer — the worked contract
