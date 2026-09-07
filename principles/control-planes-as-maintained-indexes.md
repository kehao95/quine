---
id: control-planes-as-maintained-indexes
type: principle
gist: keep indexes aligned with their owners and make stale routes repairable
status: canonical
---

# control-planes-as-maintained-indexes

A current control-plane owner defines a contract, decision, or route. An index
or projection points to those owners and must remain aligned with them. State
which source changes can make the projection stale and where a repair belongs.
These are authoring questions, not a required metadata card.

A new surface should help a participant or tool make a relevant decision. Its
presence, reachability, use, and benefit are distinct. Test the relation that
failed rather than adding another declaration of intent.

## Failure mode its absence produces

Silent drift: an index stops matching reality and no one knows to refresh it,
because nothing states when it goes stale or who absorbs the update. The
[`development/contracts/control-plane-causal-path.md`](../development/contracts/control-plane-causal-path.md)
contract carries the boundary as authoring guidance; ground-truth validators
(`check-active-doc-links.sh`, generated-projection freshness checks) own the
machine-checkable half. A hook that enforced the boundary as a five-field card
was retired 2026-07-18 after collapsing into template replication (Goodhart);
"Index projection drift" in
[`evolution/patterns.md`](../evolution/patterns.md) records the recurrence.

## Related

- [`accessible-coherence.md`](./accessible-coherence.md) — reachability is the agent-facing face of the same invariant
