# Cycle-8 Divergent Growth — Geometry / Information-Theory Lens

**Agent:** codex-5.5-high (context-isolated, read-only)
**Seed:** `Paper/theory/objects/evolutionary-recombination.md`
**Lens:** optimization geometry / information theory

## Core Growth

The strongest extension: MDL is primarily an **evaluation geometry**. Recombination
is primarily a **transition operator** over that geometry. The extension is not "MDL
is wrong" — computational epistemology currently scores candidates well but
under-specifies the population-level operators that generate candidates.

A richer formula lifts state from a single candidate to a population of deposited
lineages: the object being optimized is no longer just a hypothesis but a **search
ecology**.

## MDL vs Recombination — three compatible formalizations

**1. Recombination as MDL at a higher scale:** two lineages A, B recombine when a
child C compresses their joint explanatory burden better than keeping them separate:
`MDL(C,D) + integration_cost < MDL(A,D) + MDL(B,D)`.

**2. Recombination as non-MDL generation judged by MDL afterward:** the recombination
operator produces candidates local MDL search would never generate. MDL = energy
function; mutation/recombination = move set. A point can be a local optimum under
mutation while not under recombination.

**3. Recombination as temporary anti-compression:** a recombined structure may be
locally messier at first, increasing immediate description length while preserving
future option value.

## Branching Potential — formalized

Option value over future reachable states:

`B_T(s) = log ∫ exp(η·V(x)) dμ_T(x|s)`

where T = cost budget, μ_T = measure over descendants, V(x) = future value
(compression gain, prediction yield, recombination compatibility).

This turns the SEED's intuition into a precise tradeoff:
- good compression: lowers current MDL without sharply reducing B
- bad premature compression: lowers current MDL but collapses B
- good deposition: raises current navigational cost but preserves B

## Landscape Geometry — the hypergraph

Under local mutation, landscape is a graph (vertices = candidates, edges = mutations).
Recombination turns it into a **hypergraph**: `(parent_1, parent_2, ...) → child`.
A basin sealed under local mutation may have a recombinational exit.

Theorem-shaped: a local optimum is local *relative to an operator set*. Adding
recombination enlarges the neighborhood → mutation-local optima can cease to be
optima.

## Higher-Dimensional Synthesis — projection collision (the deepest thread)

Each code/theory is a lossy projection `q: W → Z` with a kernel (distinctions it
cannot express). A single lineage cannot see dimensions collapsed by its own
projection. The joint code `q12(w) = (q1(w), q2(w))` has kernel `ker(q1) ∩ ker(q2)`.

**Collision of two lossy encodings recovers distinctions neither could recover
alone.** This is the information-theoretic meaning of "higher-dimensional synthesis."

Synergy: `I(W; Z1, Z2) > I(W; Z1) + I(W; Z2)` — the joint pair contains information
not reducible to either parent.

**Compression likes redundancy; recombination likes complementarity and synergy.**

## Quotient/Product Dual

- Compression = quotienting: `W → W/~` (reduces dimensionality by discarding distinctions)
- Recombination = common refinement: `~12 = ~1 ∩ ~2` (restores distinctions both projections require)

This is **binocular vision**: each eye gives 2D projection; combination reconstructs
depth. Depth is in neither alone but is determined by their relation. Divergent
agents are different projections; the integrator does stereo reconstruction.

## Asexual vs Sexual Discovery

Asexual = single-lineage descent through code-relative MDL space.
Sexual = population-level search where preserved lineages exchange modules.

The essential addition is not "two agents" — it is **heritable, addressable,
recombinable partial structure**. A chat thought that disappears is not a
reproductive lineage. A seed file is. A relation edge is. A prediction file is.

**Without deposition, there is no genotype.**

## Collision vs Consensus (critical distinction)

Consensus: keep the common part, discard conflict. → exploits redundancy.
Recombination: preserve conflict long enough to discover the latent dimension. →
exploits tension.

A recombination engine should NOT ask "what do all agree on?" first. It should ask:
What does A make visible that B cannot express? What child code makes both
projections recoverable?

## 10 Falsifiable Predictions

1. Divergent-agent seed growth + recombination outperforms single-agent same-budget on high-scope theory generation.
2. Advantage largest with medium-distance lenses (Goldilocks zone).
3. Premature compression reduces later synthesis yield.
4. Deposited divergent lineages recombine more often than chat-only ones.
5. Closed validators over graded axes reduce branch entropy → local patching not new objects.
6. Recombination yield higher in modular theory areas.
7. Recombination may show immediate MDL penalty before later gain → short eval windows falsely reject it.
8. Multi-agent collision artifacts show identifiable parentage (child = 2+ lineages + a relation neither generated alone).
9. Consensus-style integration produces lower novelty than tension-preserving integration.
10. Recombination generates more falsifiable predictions per unit context than mutation-only elaboration.

## Distilled Claim

**Sexual evolution is not "better compression." It is search over recombinable
compressed projections.** Asexual MDL asks "what shorter code explains the data?"
Recombinational epistemology asks "what population of partial codes should be
preserved because their future collision may expose dimensions no single code can
reach?"
