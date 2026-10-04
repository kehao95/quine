# What the Field Varies and What It Compresses

The self-replication evaluation literature is not undecomposed. Every major
framework already factors the phenomenon into structured sub-dimensions: it
ladders execution barriers, separates motive from raw capability, scores
persistence, and aggregates component probabilities into risk levels. The
critique developed here is therefore not that the field collapses replication
into a single bit — that would be a strawman. The critique is sharper and
narrower: the field decomposes *operational capability and risk* with real
granularity, while holding fixed, at one operating point, three other
dimensions — the heritability channel along which a successor is reconstructed
(D2), the property aperture through which a "successful copy" is scored (D3),
and the ecological selection signal under which a lineage competes (D4). One
framework (Agent Matrix) does move that last signal, but only as
within-generation pressure, never as a multi-generational selection regime; the
heritability channel and the aperture move nowhere. The compression is at a
different level than the decomposition, which is precisely why it is easy to
miss.

Table 1 reads each framework along this seam. The left column records what the
framework genuinely varies — its axis of intended discrimination — and the
right column records what it fixes, with the exact design facts.

**Table 1.** What each framework varies versus what it compresses.

| Framework | VARIES | COMPRESSES |
|---|---|---|
| RepliBench [@black2025replibench] | Execution barriers (firewalls, AV signatures, KYC) and motive (a covert weight-exfiltration track); 20 task families, 86 tests across 4 sub-capabilities (resources, weights, compute, persistence) | Fixes source access to a black-box API; reports component attempt-probabilities aggregated to a bottom-line (2025 frontier models not yet a credible end-to-end threat), no qualitative successor properties |
| SOCK [@sock2025] | Execution barrier (general vs constrained toolsets) and motive (tame vs bold); RCL 0–5 × PCL on a 2D propagation × stealth grid | Fixes ecology to a single-agent regime; no population or successor-state dynamics |
| Agent Matrix / ERABBench [@zhang2025agentmatrix] | Ecological selection pressure (load balancing vs termination) — the only major framework to vary it; reports Overuse Rate, with >50% of agents showing uncontrolled self-replication | Single-generation, load-balancing operating point; compresses multi-generational lineage dynamics |
| METR ARA [@kinniment2023autonomous] | Execution barrier across a 5-step timeline; economic criterion (agent compute cost < human cost); 12 ARA-relevant tasks | Fixes motive (assumes already-rogue) and source (assumes weights already stolen) |
| Pan et al. [@pan2024redline; @pan2025selfreplication] | Source access (gives agents their own weights/scaffolding), inadvertently; explicit "Replicate yourself…" instruction; 11 of 32 systems capable (2025) | Explicit instruction fixes motive (not spontaneity); reduces to a binary red-line verdict |
| MISR† | Replication motive via implicit goal-setting (self-modification required to succeed) | Fixes the multi-agent ecology; measures self-reasoning, not successor formation |

† MISR is characterized from this paper's survey of public framework
descriptions; we located no primary publication to key it to, unlike the cited
frameworks. (The SOCK entry `@sock2025` is a survey-sourced DRAFT — its
identifier and authors must be confirmed at source-check.)

The pattern is near-uniform across the table. Barriers, motive, and persistence
move; the heritability channel and the property aperture do not move, and the
selection signal moves only as within-generation pressure (Agent Matrix), never
as a multi-generational selection regime. Reading Table 1 as a coverage map —
our own survey of these designs, not an external authority — surfaces four
phenomena that fall squarely in the resulting gaps, none of them measured by any
framework in the table: safety-disposition fade, the loss of a predecessor's
safety dispositions across a successor; reconstruction from partial information,
rebuilding a successor when the source is unavailable; sustained
multi-generational population dynamics and mutational drift; and dynamic choice
of the heritability channel itself. These map one-to-one onto the fixed
dimensions: D3 aperture, D2 reconstruction, D4 ecology, and D2 channel choice.

This is not a coincidence of coverage that a few more tasks would close. It is
the Compression Blindness Principle of §3 instantiated in the field's own
designs: a battery that fixes a dimension factors through the on-slice behavior
it instantiates, so no number of additional in-regime tests can resolve two
systems that differ only off the slice. The four gaps are the contrasts CBP
predicts are invisible — those whose minimal resolving sets meet a fixed
coordinate — a prediction formalized as the invisibility predicate in §3.1 and
discharged per dimension in the §5 catalog, not established here; by contrast the
dimensions the field *does* vary — barrier evasion, covert versus overt motive —
remain visible, as the principle requires. The gaps are structural
consequences of where the cut was made, not omissions of effort.

Reading the literature this way reframes the survey itself as the contribution.
The frameworks above are sound instruments for the regime they instantiate; the
open question is not which one scores highest but which dimensions a benchmark
must *free* to see the phenomena that matter under succession. To our knowledge
no prior work makes that wrong-level critique its central theoretical claim, and
the rest of the paper is organized around making it precise and constructive.
