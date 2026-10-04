# Context

We have a formal result called the "Compression Blindness Principle" about AI evaluation. It states: if a benchmark B operates at a fixed point r* in a multi-dimensional regime space R, fixing a subset of dimensions F, then for any property σ that depends only on the fixed dimensions, no in-regime measurement (no matter how many tests, how adaptive, how randomized) can resolve systems that agree on the varied dimensions but differ on the fixed ones. The key non-trivial parts are: (1) a closure floor — the ENTIRE sigma-algebra of conceivable in-regime evidence is capped by the on-slice observation, not just individual tests; (2) a visible/blind discrimination — some fixed-dimension-touching phenomena ARE visible (if they also perturb varied dimensions), making the principle falsifiable; (3) a constructive minimal repair — which dimension to free to see a specific phenomenon.

# Research questions

1. Does measurement theory or psychometrics have analogous results about what fixed-factor evaluation designs cannot detect? (e.g., generalizability theory, facet theory, differential item functioning)

2. Does the causal inference / identification literature (Pearl, Manski, Rubin) have results about what you cannot learn when you fix vs. condition on variables? Specifically results about the LIMITS of inference under fixed experimental design.

3. Does the experimental design literature (factorial design, response surface methodology) formally characterize what phenomena are invisible to designs that fix certain factors?

4. Does the benchmark/evaluation methodology literature in ML have any existing formal results about the structural limits of fixed-configuration benchmarks? Not just "benchmarks have gaps" but formal impossibility/information-theoretic bounds.

5. Does the "underspecification" literature (D'Amour et al. 2022) or the "predictive multiplicity" literature (Marx, Calmon & Ustun 2020, Breiman's Rashomon set) contain anything like a closure principle — that adding more tests within the same specification cannot resolve the multiplicity?

6. Is there an existing formalization of "you cannot see what you do not vary" that goes beyond the trivial partial-derivative-is-zero fact, specifically one that proves a closure/saturation result over ALL possible in-regime tests?

7. Has anyone applied Rice's theorem (non-trivial semantic properties are undecidable) or similar computability results to benchmark evaluation limits?

# Scope
- Measurement theory, psychometrics, causal inference, experimental design, ML evaluation methodology
- 2015-2025 primarily; foundational works any date

# Output
For each question: yes/no + specific papers + 1-sentence assessment of how close.
