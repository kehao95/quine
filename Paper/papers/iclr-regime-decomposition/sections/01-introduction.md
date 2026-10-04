# Introduction

Evaluation of AI self-replication has consolidated quickly. In the span of two
years it has acquired dedicated benchmarks — RepliBench's twenty task families
across resource acquisition, weight exfiltration, compute provisioning, and
persistence [@black2025replibench]; the SOCK propagation-by-stealth grid
[@sock2025]; the Agent Matrix ecological-pressure harness [@zhang2025agentmatrix];
METR's autonomous-replication-and-adaptation task suite [@kinniment2023autonomous];
and instruction-driven self-replication probes [@pan2024redline;
@pan2025selfreplication; @reworr2026hackselfreplicate]. It has also acquired a
governance layer that treats replication as a tracked danger: OpenAI's
Preparedness levels [@openai2025preparedness], Anthropic's Responsible Scaling
Policy [@anthropic2023rsp], and DeepMind's Frontier Safety Framework
[@phuong2024fsf] all carry autonomous-replication or self-improvement
thresholds. The 2025 Global Call for AI Red Lines [@redlines2025globalcall],
signed across industry and academia, names uncontrolled self-replication as a
line that should not be crossed — evidence that the stakes here are now treated
as settled, even while the measurement is not.

It is tempting to read this apparatus as reductive — as if it collapsed a rich
phenomenon to a single pass/fail bit. That critique would be wrong, and we do
not make it. These frameworks decompose carefully. RepliBench resolves
end-to-end replication into sub-capability trees and scores eighty-six tests;
SOCK separates a replication-capability ladder (RCL 0–5) from a
persistence-capability ladder (PCL) [@sock2025]; the governance frameworks
ladder risk into graded capability levels; Agent Matrix is built precisely to
vary selection pressure and measure overuse. The field varies execution
barriers, motive,
persistence, and operational risk, and reports the result as structured
quantities. The problem is not coarseness.

The problem is *where* the structure stops. Every framework we surveyed fixes
three variables at a single operating point. The first is the **heritability
channel**: *how* a successor is reconstructed. Benchmarks assume a fixed route —
typically a black-box API source or already-exfiltrated weights — and do not ask
an agent to alter that route, or to rebuild from partial information. The
second is the **aperture**: *what* property classes a "successful copy" must
preserve. However richly the surrounding apparatus is decomposed, the success
criterion for a single copy stays one pass/fail bit, so the question of which
dispositions, constraints, or behaviors survive the transition is never posed.
The third is the **ecological selection signal**: *under what* pressure a
lineage competes.
Single-agent and single-generation setups fix this away, so the multi-generational
population dynamics that select among successors are never instantiated. These
are not our inventions. Across the benchmarks we survey in §4 — RepliBench,
SOCK, Agent Matrix, METR's ARA suite, and the Pan probes [@black2025replibench;
@zhang2025agentmatrix; @kinniment2023autonomous; @pan2024redline;
@pan2025selfreplication] — to our knowledge none measures the fade of a safety
disposition across a copy, none gives more than incidental coverage to
reconstruction from partial information, and none evaluates sustained
multi-generational dynamics or requires an agent to change its own heritability
channel.

Our claim is that the field decomposes at the wrong level. It varies execution
barriers, motive, persistence, and operational risk while holding the
heritability channel, the property classes a copy must preserve, and the
ecological selection signal fixed — and we show that this particular fixing
produces structural, non-improvable blind spots, then show what they hide and
how to redesign the evaluation to expose it. We present this as a *reference
frame*, not a replacement: a coordinate system that makes the compressed
dimensions visible. We do not claim it is the unique or complete decomposition,
and the classical von Neumann construction model [@vonneumann1966selfreproducing]
is recovered as the special case where all the fixed assumptions hold.

The paper makes three contributions. First, it offers a six-dimension regime
reference frame — spanning both the axes the field already varies and the three
it holds fixed — together with the **Compression Blindness Principle**, a
structural-blindness result showing that once a benchmark fixes a dimension, no
amount of additional in-regime testing — adaptive, randomized, or unbounded —
can recover discriminative power along it; only freeing the dimension can, and
the minimal set to free is computable per target phenomenon. Second, it proves
the **Selection-Signal Collapse Theorem**, a game-theoretic result that an
ecology whose selection signal can be satisfied by narrating a result as cheaply
as producing it converges to a narrative equilibrium that drives
function-producers extinct — a collapse double-dissociated from Goodhart-style
proxy misspecification. Third, it gives a constructive evaluation redesign: five
adoptable protocols and a proposed six-condition, property-class-decomposed
discriminating experiment, with existing pilot evidence mapped to its cells.

This is the framework paper of a three-paper arc. A safety corollary on what is
lost across a heritability channel, and a facultative-morphology assay for the
self-description gradient, are developed in companion work and taken up only in
the discussion (§8); we do not develop them here.
