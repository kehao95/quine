# ICLR-1 Paper: Collision Implications

> Feed from `theory/references/facultative-succession-theory-collision-2026-06-30.md`.
> Insights for the evaluation blind-spots paper.

## Compression Blindness Principle: stat mech formalization

The collision gives a non-trivial angle on compression blindness:

**"Irrelevant operator" framing**: under low demand, source code is an
irrelevant operator in the RG sense — it exists in the microstate but flows
to zero in the macroscopic behavior. This is WHY the benchmark at the
field's operating point (low demand, high material) is blind: at that
operating point, material is irrelevant. The benchmark confirms "source
present → success" but the source is ceremonial (Meyer-Rowan decoupling).
The real action happens at high demand where source becomes a relevant
operator.

This means the compression blindness principle can be stated as:

> A benchmark operating at point r* has zero discriminative power along
> dimension D_i if D_i is an irrelevant operator at r*. Adding more tests
> at r* cannot make D_i relevant — only moving to a different point in
> regime space can.

## Selection-signal collapse: spin glass connection

The degradation-basin result connects selection-signal collapse to spin
glass physics:
- narration-satisfiable selection = zero-cost signal → pooling equilibrium
  (all forms look the same to the selector)
- exec-dependent selection = costly signal → separating equilibrium
- The spin glass framing predicts: under quenched disorder (degraded
  material), the free energy landscape is rugged with many local minima
  → compensation basins are metastable states

## Evaluation redesign: reaction norm protocol

The evo-devo lens suggests the evaluation redesign should be:

> Instead of measuring replication success at one point, measure the
> REACTION NORM: the full function from (material, demand) to successor
> form. A benchmark that reports one point on the reaction norm is like
> a geneticist who measures phenotype at one temperature and claims
> the genotype has no plasticity.

This is a concrete, publishable recommendation for benchmark designers.

## Figures for ICLR-1

1. **Compression blindness diagram**: field's operating point marked on
   phase diagram; arrow showing "blind directions"
2. **Reaction norm surface**: 3D plot of material × demand → form complexity
3. **Benchmark comparison table**: what each benchmark varies vs what it
   compresses (already in survey)
