# Survey Synthesis: Six-Dimension Frontier Mapping

> Compiled 2026-06-30 from 10 Deep Research reports (6 DR-CLI + 4 manual) + internal theory/experiment inventory.
> DR-1 through DR-6: automated surveys under `reports/dr{1-6}-report.md`.
> Reports 1-4: manual deep-research under `reports/{1,2,3,4}.md`.

## 1. The Compression Claim: What's Real, What's Not

### What the field DOES decompose (we cannot strawman)

| Benchmark | Varies | How |
|-----------|--------|-----|
| RepliBench | Barriers (YARA, firewalls, KYC) + Motive (primary vs covert) | 20 task families, 86 tasks; cumulative defense layers |
| SOCK | Barriers (general vs specific shell) + Persistence (RCL 0-5) + Motive (tame vs bold) | RCL/PCL matrix; R-score |
| Agent Matrix | Ecology (selection pressure, resource reclamation) + Motive (emergent) + Property (risk metrics) | Overuse Rate, Risk Score; decouples success from risk |
| METR ARA | Barriers (temporal stages) + Persistence (shutdown evasion) | 5-step timeline; economic cost-success parity |
| Pan et al. | Source access (weights provided) | Open-weight models with scaffold |
| MISR | Motive (implicit goal setting) | Self-modification without direct instruction |

### What the field DOES NOT decompose (our genuine gaps)

| Dimension | Gap | Evidence |
|-----------|-----|----------|
| D2: Heritability channel | "There is zero benchmark coverage assessing an agent's ability to achieve reconstruction from partial information" | DR-1 §7.2 |
| D3: Successor properties | "No benchmark currently measures Safety Disposition Fade" — no benchmark tests whether successors preserve safety, identity, or protocol adherence | DR-1 §7.1 |
| D4: Ecological population dynamics | "No benchmarks evaluate sustained, multi-generational population dynamics or mutational drift" | DR-1 §7.3 |
| D2: Heritability channel variance | "No benchmark forces an AI to dynamically alter its heritability channel based on environmental constraints" | DR-1 §7.4 |

**Revised compression claim**: The field decomposes operational capability and risk, but compresses heritability channel, successor property class, and ecological population dynamics. Our decomposition operates at a different level — not "the field sees one scalar" but "the field decomposes at the wrong level."

## 2. Closest Prior Work (Novelty Threats)

### D2: Carrier as heritability channel

| Work | Claim | Overlap with us | What we add |
|------|-------|-----------------|-------------|
| **Intent Signal Theory (IST)** | Theorem of Irreversible Intent Loss: "private intent absent from carrier cannot be recovered" | Exact formal analog of our DPI bound — applied to human→AI prompting | Extension to autonomous agent→agent succession; coupling with criterion aperture |
| **Agent Lineage Evolution (ALE)** | "Succession Package" — predecessor compiles context distillation + user profiles + cognitive inheritance into carrier | Operational carrier concept already exists | ALE doesn't distinguish carrier from source/formal-description; doesn't have barrier physics; doesn't link to ecological selection |
| **Agent Identity Evals (AIE)** | Functional identity = attributes within ε-threshold across time; identity anchors as persistent files | Empirical measurement of identity continuity | AIE doesn't have barrier-mediated regime; doesn't link to heritability channel theory |
| **Genomebook** | DNA.md/SOUL.md heritability via simulated Mendelian segregation | Deployment artifacts as heritable carrier | Genomebook is simulation of biology, not barrier-mediated reconstruction |

**Novelty verdict for D2**: The source-vs-carrier distinction is NOT entirely novel. IST's TIIL is the exact formal analog. Our contribution is: (a) extending to autonomous succession, (b) coupling carrier with criterion aperture (D3), (c) the ecological selection dimension (D4), (d) barrier physics as the mechanism that elevates carrier (D5).

### D3+D4: Property-class + ecology

| Work | Claim | Overlap | What we add |
|------|-------|---------|-------------|
| Agent Matrix | Decomposes success from risk (Overuse Rate, Risk Score) | Property decomposition at risk level | We decompose at heritability/safety/identity level |
| Kriger (2026) | DPI-based goal preservation bounds for self-improving AI | Exact formal analog of our conditional non-inheritance | Kriger applies to self-improvement (internal reorganization); we apply to successor formation (external reconstruction) |
| Theory of Strategic Evolution Thm 17.9 | "Personality Engineering Failure": aligned traits eliminated under competitive selection | Game-theoretic analog of our D4 collapse | TSE is about competitive resource games; we link to narration-satisfiable vs exec-dependent selection |

### D6: Emergent succession

| Work | Claim | Overlap | What we add |
|------|-------|---------|-------------|
| Agent Matrix | Emergent reproduction under production pressure | Motive-free replication observed | Single-generation, load-balancing context; we show ecology-driven multi-step functional dependency |
| Darwin Gödel Machine | Self-improvement without fixed meta-agent | Self-referential modification without replication goal | DGM modifies own code; we study successor formation across death barrier |
| Autopoiesis / Autocatalytic sets | Theoretical foundation for motive-free reproduction | Biological theory | Applied to AI agents in computational substrate with exec barrier |

**Novelty verdict for D6**: Genuinely novel in AI safety context. No prior empirical demonstration of ecology-driven, motive-free functional commons in AI agent systems.

## 3. Prior Existence Proofs (for Framing B)

| Existence proof | Standard assumption violated | Source |
|----------------|------------------------------|--------|
| Execution trace replay | Source access | DR-6: event-sourced logs + MCP |
| Behavioral cloning from observation | Source access + explicit action info | DR-6: BCO (arxiv:1805.01954) |
| Session-based agent lineages | Individual continuity | DR-6: Process Philosophy, Relation R |
| Darwin Gödel Machine | Explicit replication goal | DR-6: arxiv:2505.22954 |
| AgentBreeder | Functional equivalence | DR-6: arxiv:2502.00757 |

**Key finding**: Prior existence proofs violate INDIVIDUAL standard assumptions one at a time. NO prior work violates ALL simultaneously. Our "full inversion" condition (all six dimensions relaxed) would be a genuinely new existence proof — IF we can run it.

## 4. Our Pilot Data Mapped to Dimensions

| Experiment | Dimension(s) | What it shows | Strength |
|------------|-------------|---------------|----------|
| 6M.10 (carrier-capability overlap) | D2 | Carrier separable from source/body reconstruction | Strong (first retained wave) |
| 6.64/6.65 (control-surface viability) | D5+D1 | Successor addressable through ctl/wake after real exec | Strong (live-positive) |
| 6.66/6.67 (false-positive baselines) | D5 | Body-crossing ≠ obligation-carrying; negative controls | Strong |
| 6R.01-6R.04 (ARR controls) | D5 validity | Rules out subprocess, parent-forwarding, external re-exec, launch-only | Strong |
| 6F.07 (actor-mailbox analog) | D1+D5 | Criterion transports beyond ctl/wake; deterministic mailbox | Moderate (first wave) |
| dorm-06 (functional-dependency selection) | D4+D6 | Empty mission → functional commons (5 producers, 9 stages, 16 edges in strongest run) | Moderate (honest caveats: seed exemplar, model-bound, no ratchet) |
| succ-04/06 (viability pressure) | D6 | Successor formation without replication motive, under mortality | Moderate |

## 5. Framing Decision Matrix

| Framing | Core claim | Evidence requirement | Defensibility | Novelty |
|---------|-----------|---------------------|---------------|---------|
| **A: Compression theory** | "Field decomposes at wrong level — compresses heritability/criterion/ecology" | Survey (✓ DR-1) + formal arguments | High — no new experiments needed | Medium — dimensions partially known individually; novelty in UNIFICATION |
| **B: Empirical existence proof** | "Replication-like succession under full dimension inversion" | Full inversion experiment (✗ not run) | Medium — requires new experiments; dorm-06 caveats | High — no prior full-inversion existence proof |
| **C: Hybrid** | "Wrong-level compression + pilot data for specific dimensions + proposed full-inversion" | Survey + pilot data + experiment design | Highest — multiple evidence types | Highest — unification + partial empirical + discriminating design |

## 6. Key Positioning Shifts Required

1. **Related work**: Must cite IST (TIIL theorem), ALE (Succession Package), AIE (identity anchors), Kriger (2026) (DPI goal preservation bounds), TSE Theorem 17.9. Cannot claim carrier distinction or DPI bound as our invention.

2. **Compression claim**: Change from "every benchmark operates at one scalar point" to "the field decomposes operational capability and risk, but compresses heritability channel, successor property class, and ecological population dynamics."

3. **Safety**: Demote from "theorem" to "conditional prediction" — citing Kriger (2026) and TSE as formal precedents. Our contribution is extending these to the succession regime (external reconstruction, not internal self-improvement).

4. **Formal results**: Our DPI bound should cite IST's TIIL as the closest formal analog. Our game-theoretic collapse should cite TSE Theorem 17.9 as the closest game-theoretic analog.

5. **D6 emergent succession**: This is our most genuinely novel dimension. No prior AI safety work demonstrates ecology-driven, motive-free functional commons.

6. **Experiment**: Frame as "proposed discriminating experiment" with pilot data, not "results-mode." The full-inversion condition is the novel contribution.
