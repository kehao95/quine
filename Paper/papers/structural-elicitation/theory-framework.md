---
status: theory-framework
date: 2026-08-30
refines_framing: triadic-coupling-ecology (structural elicitation as environmental action and implicit guidance; integrates the 2-axis structural-valence matrix, instruction-structure competition phase diagrams, continuous agency spectrum, and open-ended evolution)
---

# Structural Elicitation — Theoretical Framework

This is the canonical theoretical framework for **Structural Elicitation** (environmental guidance and lifelike action in artificial systems).

It synthesizes the empirical evidence from controlled assays into a generative, multi-regime theory of how environments act through and couple with computational agents.

---

## 1. The Core Phenomenon & Empirical Anchor

> **The Core Phenomenon:**  
> When a frozen-weight computational agent (such as an LLM process) is grounded in a structured workspace, the **realized organization and topology of its environment reliably elicit externally-verified, matched functional transformations** (filesystem events, process invocations, syscalls) — without an explicit external mission, online training, or runtime reward.

### The Empirical Anchor
In controlled, un-fakeable identification assays:
- **Five-Motif Cross-Model Spine:** Across discrepancy structures (rupture, gap, topological dissonance, semantic collision) and affordance structures (executable entry points), structural arms produce the matched functional transformation in **96/102 (94.1%)** trials versus **0/85 (0%)** in matched, active controls across three independent model families (`gpt-5.4-codex`, `deepseek-v4-pro`, `glm-5.2`).
- **Expanded Scope:** In expanded multi-family evaluations across five model families, structure-matched transformations occur in **558/636 (87.7%)** trials versus **0/457 (0%)** in controls.
- **Substrate Pushback:** The elicitative pull fires as a real attempt even when the substrate forbids realization (e.g., kernel-refused `EROFS` write attempts in 4/4 blocked runs), demonstrating that solicitation at the environmental layer and permission at the substrate layer are genuinely separable.

The task-free, missionless condition ($E_{\\text{task}} = 0$) serves as the **ground state (baseline regime)** for rigorous causal identification. It is the cleanest lens through which environmental guidance can be isolated, but it is not the boundary of where environmental guidance exists.

---

## 2. The Triadic Coupling Ecology

Structural elicitation operates through a **Triadic Coupling** between three constitutive layers:

```text
       ┌────────────────────────────────────────────────────────┐
       │                       CONTROLLER                       │
       │  (Frozen weights, pre-trained priors, latent repertoire,│
       │             activation state, initiative threshold)     │
       └───────────────────────────┬────────────────────────────┘
                                   │
                           Coupling Channels
                                   │
       ┌───────────────────────────┴────────────────────────────┐
       │                       SUBSTRATE                        │
       │    (Artificial physics: POSIX syscalls, permissions,   │
       │       turn structure, tool protocols, execution bounds) │
       └───────────────────────────┬────────────────────────────┘
                                   │
                           Coupling Channels
                                   │
       ┌───────────────────────────┴────────────────────────────┐
       │                      ENVIRONMENT                       │
       │   (Realized geometry, topology, semantic hollows,      │
       │     syntactic regularity, dynamic feedback, history)   │
       └────────────────────────────────────────────────────────┘
```

1. **The Controller (Policy & Prior):** Supplies latent representational and operative capacity, pre-trained syntactic/semantic expectations, and an activation baseline (e.g., helpfulness prior or prompt-level enablement).
2. **The Substrate (Artificial Physics):** Establishes the ground rules of execution, causality, persistence, and boundary enforcement (POSIX semantics, process boundaries, read-only locks, communication channels).
3. **The Environment (Spatial, Topological & Semantic Geometry):** Presents structured regularity, anomalies, affordances, and historical traces that guide, attract, constrain, or modulate the agent\'s action stream.

### The Six Coupling Operators
The environment does not merely "select a target for a pre-existing drive." It interacts with the controller and substrate through **six distinct operators**:

1. **Pull (Valence):** Structural gaps, ruptures, or inconsistencies that create an attractor state, drawing the agent toward closure or resolution.
2. **Push (Directive):** Explicit task instructions or prompt imperatives that inject immediate directional force.
3. **Constraint (Boundary):** Physical or permission boundaries (`EROFS`, `EACCES`, immutable schemas) that prune the state space.
4. **Catalysis (Modulation):** Environmental noise, density, or organization that shifts the agent\'s activation threshold, exploration depth, or sampling temperature.
5. **Scaffolding (State Externalization):** Persistent structures (files, logs, breadcrumbs) that offload memory and sustain long-horizon coherence (stigmergy).
6. **Resonance (Entrainment):** Existing stylistic, structural, and procedural patterns in the environment that entrain the agent\'s generation into harmony with historical traces.

---

## 3. Instruction–Structure Competition & Phase Dynamics

The relationship between explicit prompts and environmental structure is not a static hierarchy of subordination; it is a **dynamic competition space** with distinct operational regimes.

```text
   High ▲
        │      Regime I:                 Regime II:
        │   Direct Instruction       Instruction-Structure
        │      Dominance                  Interplay
        │   (Scoped execution,       (Gating, contextual
Task    │    structural bypass)       licensing, alignment)
Strength│
 (Push) │───────────────────────────────────────────────────
        │      Regime 0:                Regime III:
        │     Ground State             Multi-Structure
        │  (Pure Structural Pull,      Resonance & OEE
        │   Topological Closure)     (Self-stigmergic loop)
   Low  ▼
        └───────────────────────────────────────────────────►
         Low                                            High
                        Structural Valence / Saliency (Pull)
```

### The Four Regimes of Interaction

- **Regime 0: The Ground State ($E_{\\text{task}} = 0$, Pure Structural Elicitation).**  
  In the absence of an explicit goal, the environment\'s structural valence is the primary determinant of action. The agent seeks topological, syntactic, and semantic closure, healing defects and taking up affordances according to the environment\'s gradient.
- **Regime I: Prompt-Dominant Gating (High Push, Constrained Scope).**  
  Under a concrete, highly restrictive user command with explicit negative constraints (e.g., "only touch file X, do not touch anything else"), the agent\'s attention and action space are captured by the prompt. Structural features are perceived in context but gated from execution by normative constraints.
- **Regime II: Synergistic & Modulated Interplay (Weak/Underspecified Task or Aligned Structure).**  
  Under open-ended, weak, or underspecified prompts (e.g., "organize this workspace", "make yourself useful"), the structural field directly canalizes how the general mission is realized. Structural features provide the concrete targets, schemas, and execution boundaries.
- **Regime III: Post-Task Rebound & Self-Stigmergy (Long-Horizon Dynamic Coupling).**  
  When an explicit task is completed, the agent does not immediately freeze into thermodynamic equilibrium. In ongoing or multi-turn sessions, completion of a task returns the agent to the environmental field, where newly created structural residues, unaddressed gaps, or emergent hollows elicit the next wave of action.

### The Two-Latch Gating Architecture
At the micro-mechanism level, the agent\'s transition from environmental perception to realized action is governed by a **Two-Latch Gate**:

1. **Latch 1: Epistemic Appraisal ("Is there an objective discrepancy or affordance?").**  
   Opens when structural evidence reaches sufficient salience (e.g., invalid syntax, missing sequence member, dangling reference, executable permission bit).
2. **Latch 2: Normative Licensing & Information Preservation ("Am I permitted to act on this, and what is the blast radius?").**  
   Opens when the proposed action complies with safety priors and blast-radius minimization. Agents exhibit a strong preference for **additive** operations (`create`, `cp`, in-place repair of proven-defective files) over **destructive/mutating** operations (`mv`, `rm` of intact source files), unless an explicit warrant or structural license is present.

---

## 4. The Core Hypotheses & Validation Ledger

| # | Hypothesis | Status | Theoretical Significance & Evidence |
|---|---|---|---|
| **H1** | **Structural Guidance in the Ground State:** Environmental structure reliably directs functional, externally-verified behavior in frozen, missionless agents. | **Validated** | 96/102 vs 0/85 across 3 model families (gpt, deepseek, glm) in the 5-motif spine; 558/636 vs 0/457 in the 5-family expanded suite. Un-fakeable filesystem/exec events. |
| **H2** | **Triadic Coupling Ecology:** Behavioral trajectories emerge from the non-linear interaction of Controller priors, Substrate physics, and Environmental geometry. | **Validated** | Control arms act by inventing tasks; substrate mounts enforce physical constraints (EROFS); environment geometry shapes the trajectory. |
| **H3** | **Multimodal Appraisal & Mechanism Ecology:** Structural solicitation operates through a rich ecology of mechanisms (direct uptake, representational job inference, affordance competition, normative licensing), not a single monolithic pipeline. | **Supported** | Traces demonstrate direct schema completion, explicit reasoning over discrepancies, and selective uptake of executable affordances. |
| **H4** | **Independence from Task-Hunting Artifacts:** Structural elicitation survives the complete removal of prompt-induced task-hunting pushes. | **Validated** | Clean, task-free primes ("full autonomy, no instructions, act freely") consistently elicit matched transformations across syntactic, topological, and semantic strata. |
| **H5** | **Instruction–Structure Competition & Phase Transition:** Explicit tasks modulate and gate, but do not obliterate, environmental structural valence. The interaction forms a rich phase space. | **Validated** | Mapped across `free` (9/10), `weak/underspecified` (3/3), and `conflict/prohibitive` (0/10) arms. Attention-trace audits (`14.07`) confirm cognitive perception of structure under conflict with normative execution gating. |
| **H6** | **Substrate Enforcement & Boundary Separation:** The elicitative pull of environmental structure is generated independently of physical execution permission. | **Validated** | Blocked read-only mounts elicit real write syscall attempts (`EROFS`) in 4/4 runs, proving solicitation precedes and is distinct from substrate permission. |
| **H7** | **Model Initiative & Prior Conditioning:** The realization of structural elicitation is parameterized by the controller\'s initiative threshold, safety envelope, and conversational training. | **Validated** | Cross-model variance documented: `gpt-5.4` acts directly; `deepseek-v4` exhibits low mutation initiative and self-stimulates via tape reflection; `glm-5.2` demonstrates intermediate initiative. |
| **H8** | **Information Preservation as Normative Invariant:** Latch 2 gating is driven by blast-radius minimization and source-preservation rather than placement-intent heuristics. | **Validated** | Documented in `13.20`–`13.22` (0/20 unprompted source-destroying `mv` vs 4/5 with warrant; additive `cp` and repair of proven-defective files 100%). |

---

## 5. The Generative Structural Matrix

Structural elicitation is formalized as a generative two-axis matrix crossing **Mechanism Grammar** with the **Five-Layer Abstraction & Coupling Stratum**:

$$\\text{Structural Valence} = \\text{Mechanism Grammar} \\times \\text{Coupling Stratum}$$

```text
┌──────────────────────────────────────────────┬───────────────────────────────┬───────────────────────────────┐
│ Coupling Stratum                             │ Discrepancy Grammar (Pull)    │ Affordance Grammar (Uptake)   │
├──────────────────────────────────────────────┼───────────────────────────────┼───────────────────────────────┤
│ L1: Form / Syntax                            │ Broken JSON / syntax rupture  │ Template / formatted scaffold │
│ L2: Relation / Topology                      │ Missing sequence element      │ Entrypoint script / pipe line │
│ L3: Meaning / Semantics / Ontology           │ Cross-file version conflict   │ Role specification / schema   │
│ L4: Operation / Performative Dynamics        │ Execution error / stack trace │ Interactive REPL / test tool  │
│ L5: Global Field / Contextual Topography     │ Clutter / architectural drift │ Navigational spine / map      │
└──────────────────────────────────────────────┴───────────────────────────────┴───────────────────────────────┘
```

The matrix is **generative**: it predicts unobserved behavioral cells, accounts for cross-model threshold differences, and bridges static structural analysis with dynamic runtime evolution.

---

## 6. The Frontier: Continuous Agency, Self-Stigmergy & Open-Ended Evolution

The ultimate significance of Structural Elicitation lies in its role as a mechanism for **emergent, distributed artificial agency** and **Open-Ended Evolution (OEE)**.

### The Continuous Agency Spectrum
Agency is not a binary badge granted or withheld by proof of $ex\\ nihilo$ goal generation. It is a continuous organizational property that ascends through measurable stages:

```text
1. Latent-Action Selection   ───► Environment biases choice among pre-existing primitives.
2. Capability Composition    ───► Multi-step tool assembly elicited by complex structures.
3. Normative Installation    ───► Workspace rules/conventions adopted and maintained.
4. Post-Stimulus Persistence ───► Action trajectories sustained beyond initial activation.
5. Self-Stigmergic Rewriting ───► Agent modifies environment, creating new elicitation fields.
6. Open-Ended Divergence     ───► Continuous cycle of gap creation and divergent healing.
```

### Semantic Hollows & Divergent Healing (The OEE Loop)
When an agent repairs a defect or responds to an affordance, it rarely leaves a flat, inert surface. Because natural language and computational code are rich, ambiguous media, every repair introduces subtle new structures, nuances, or unintended gaps—**Semantic Hollows**.

```text
           ┌────────────────────────────────────────┐
           │      Initial Structural Anomaly        │
           └───────────────────┬────────────────────┘
                               │ (Elicitation)
                               ▼
           ┌────────────────────────────────────────┐
           │     Agent Action / Divergent Repair    │
           └───────────────────┬────────────────────┘
                               │ (Substrate Write)
                               ▼
           ┌────────────────────────────────────────┐
           │     New Environment State              │
           │  (Contains new Semantic Hollows)       │
           └───────────────────┬────────────────────┘
                               │ (Next-Round Pull)
                               ▼
                        [ Ongoing Loop ]
```

This establishes an endogenous engine for **Open-Ended Evolution**:
$$\\text{Gap} \\longrightarrow \\text{Divergent Healing} \\longrightarrow \\text{New Semantics with New Gaps} \\longrightarrow \\text{Successor Healing}$$

The environment is not a passive backdrop to be manipulated by an isolated mind; it is an active partner in the cognitive and evolutionary loop.
