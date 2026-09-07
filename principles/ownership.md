---
id: ownership
type: principle
gist: make responsibility for working contracts and shared facts findable
status: canonical
---

# ownership

Make the source of a working contract or shared operational fact clear, so
participants know where to update it and which consumers depend on it.
Ownership follows the reason to change. Split a document when separate uses
benefit from independent maintenance; explanations can share material and
support different interpretations.

## Worked contract

[`experiments/PROTOCOL.md`](../experiments/PROTOCOL.md) § Phase Interpretation
Layer is the worked contract: belief → phase README, technique → analysis
instrument doc, new run evidence → per-run analysis file, experiment-local →
child README.

## Failure mode its absence produces

Parallel homes: a fact grows N restated copies because no surface owns it, and
presence-checks stay green while the copies drift. The "Declared ownership
without compliance pressure" pattern in
[`evolution/patterns.md`](../evolution/patterns.md) names the related failure
where a card declares ownership the body contradicts.

## Related

- [`single-canonical-location.md`](./single-canonical-location.md) — ownership is what makes a single location stay canonical
