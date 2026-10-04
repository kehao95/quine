---
configuration_id: lba
claim_bundle: ../
venue: alife-2026-lba
format-envelope: 2page-lba
role: venue-probe
maturity-target: mechanism
status: draft
paused: true
state-mode: full-manuscript-home
---

# Substrate Coincidence LBA Configuration

> **PAUSED — venue binding lapsed `2026-07-26`.** The ALIFE 2026 LBA window
> closed `2026-07-20` and this configuration never submitted (see
> [`alife-2026-lba.md`](../../../../venues/objects/alife-2026-lba.md) § Cycle
> outcome). It is withdrawn and not re-bound.
>
> This is a `role: venue-probe` configuration, so with the venue gone it has no
> live job — but it owns no claim of its own to lose. **The claim bundle
> ([`../`](../README.md)) is untouched**: the three-separation-collapse thesis
> and the null-05 evidence live there, and the compression work in
> [`manuscript.md`](manuscript.md) stays as a usable 2-page projection for the
> next short-format window. Re-bind the probe when there is one.

This configuration binds the substrate-coincidence claim bundle to the
[`alife-2026-lba`](../../../../venues/objects/alife-2026-lba.md) venue and the
[`2page-lba`](../../../../venues/envelopes/2page-lba.md) format envelope.

## Configuration Role

The LBA is a venue probe: it compresses the substrate-coincidence structural
theory into one bounded phenomenon for the ALIFE late-breaking abstract route —
the simultaneous collapse of the three evolutionary separations (substrate-type,
individual boundary, Weismann barrier) from a single root mechanism
(read = dissolution) in a language-based agent medium, illustrated by the
null-05 missionless-agent runs.

## Source Surfaces

- claim bundle / manuscript source: [`../manuscript-draft.md`](../manuscript-draft.md)
- manuscript source (LBA): [`manuscript.md`](manuscript.md)
- build metadata: [`metadata.yaml`](metadata.yaml)
- venue object: [`../../../../../venues/objects/alife-2026-lba.md`](../../../../venues/objects/alife-2026-lba.md)
- experimental basis: [`experiments/active/null-ecology/null-05-reflexive-self-embedding/`](../../../../../experiments/active/null-ecology/null-05-reflexive-self-embedding)

## Build

```
make pdf     # output/Substrate_Coincidence_LBA.pdf
make tex     # intermediate LaTeX
make html    # standalone HTML
```

## Query Maintenance

This file owns only the configuration binding and local manuscript entrance.
The claim, full theory ladder, and cultural-evolution positioning stay in the
parent dossier (`../README.md`, `../manuscript-draft.md`).
