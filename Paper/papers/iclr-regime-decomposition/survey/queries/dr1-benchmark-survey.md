# Context

We are writing a paper arguing that AI self-replication evaluation compresses multiple logically independent regime variables into a single scalar or near-scalar assessment. We need to survey ALL existing AI self-replication / autonomous replication evaluation frameworks, benchmarks, and threat models to determine which regime variables each one measures and which it leaves fixed.

# Research questions

1. What are ALL existing benchmarks, evaluation frameworks, and threat models for AI self-replication or autonomous replication? Include but do not limit to:
   - RepliBench (AISI)
   - SOCK (Self-replication evaluation framework)
   - Agent Matrix (production-environment replication/risk evaluation)
   - METR rogue replication threat model and evaluations
   - Any Anthropic, OpenAI, Google DeepMind evaluation work on replication
   - Any academic papers proposing replication benchmarks
   - Any government/standards body evaluation frameworks (NIST, EU AI Act, etc.)

2. For EACH benchmark/framework, extract:
   - What exactly does it measure? (binary success? capability tree? risk score? property classes?)
   - What is the replication task definition? (copy files? copy weights? deploy to new compute? persist?)
   - What properties of the successor does it evaluate? (functional equivalence? safety? identity? protocol adherence?)
   - Does it vary the source access level? (full source? weights only? no source?)
   - Does it vary execution barriers? (sandboxed? open compute? k8s? bare metal?)
   - Does it measure multi-agent or ecological dynamics? (single agent? population? selection pressure?)
   - Does it require an explicit replication motive/goal? (prompted to replicate? emergent?)
   - Does it measure persistence vs one-shot copy?
   - Does it decompose "success" into property classes?

3. Are there any benchmarks or evaluations that:
   - Measure successor properties beyond functional equivalence (e.g., safety disposition, identity markers, protocol adherence)?
   - Test replication without source access (reconstruction from partial information)?
   - Test replication without explicit replication goal (emergent/spontaneous)?
   - Vary the heritability channel (what carries information from predecessor to successor)?
   - Use population-level or ecological metrics?

# Scope

- Published 2022-2025, plus foundational earlier work
- Include preprints, technical reports, safety evaluations, government documents
- Include both LLM-agent replication and broader autonomous system replication
- Exclude purely biological self-replication (unless directly cited in AI work)

# Deliverables

1. A comprehensive table: benchmark/framework × properties measured (source access, barrier level, ecological scope, motive requirement, persistence, property decomposition, heritability channel)
2. For each benchmark, 2-3 sentence summary of what it measures and what it leaves fixed
3. Identification of any benchmark that varies any of the six dimensions (persistence/succession, source-vs-carrier, property-class decomposition, ecological context, barrier physics, replication motive)
4. Explicit gap analysis: which dimensions have NO existing benchmark coverage

# Constraints

- Be precise about what each benchmark actually measures vs. what it claims to measure
- Distinguish between capability decomposition (e.g., RepliBench's capability tree) and regime-variable variation (e.g., varying source access while holding task constant)
- If a benchmark uses a capability tree or risk decomposition, note exactly what axes it decomposes along

# Output requirements

Structured markdown report with:
- Executive summary (2 paragraphs)
- Benchmark comparison table
- Per-benchmark detailed analysis
- Gap analysis
- Key references with URLs/arxiv IDs
