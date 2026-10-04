# Appendix B: Comparison with Related Systems

This appendix provides an extended comparison between Quine and related systems from both classical Artificial Life and modern LLM agent research.

## B.1 Feature Comparison Matrix

**Table B1.** Feature comparison across A-Life and LLM agent systems.

| Feature | Tierra | Avida | Lenia | Polyworld | Voyager | AutoGPT | **Quine** |
|---------|--------|-------|-------|-----------|---------|---------|-----------|
| **Year** | 1991 | 1993 | 2019 | 1994 | 2023 | 2023 | 2024 |
| **Substrate** | VM Assembly | VM Assembly | Continuous CA | Neural Nets | LLM + Minecraft | LLM + Tools | LLM + POSIX |
| **Self-Reproduction** | $\checkmark$ | $\checkmark$ | $\times$ | $\times$ | $\times$ | $\times$ | **$\checkmark$ (Process spawning)** |
| **Metabolism** | VM allocation | VM allocation | $\times$ | VM allocation | $\times$ | $\times$ | **POSIX resource bounds** |
| **Bounded Lifespan** | VM-enforced | VM-enforced | $\times$ | VM-enforced | Task iteration limits | Task iteration limits | **Multiple mortality modes** |
| **Heredity** | $\checkmark$ | $\checkmark$ | $\times$ | $\checkmark$ | External memory | External memory | **Wisdom inheritance** |
| **Evolution** | Darwinian | Darwinian | $\times$ | Darwinian | Skill accumulation | $\times$ | **Variation + selection** |
| **Practical Utility** | $\times$ | $\times$ | $\times$ | $\times$ | $\checkmark$ | $\checkmark$ | **$\checkmark$** |
| **Semantic Observability** | $\times$ | $\times$ | $\times$ | $\times$ | Partial (Post-hoc) | Partial (Post-hoc) | **$\checkmark$ (Causal Tape)** |
| **Real Environment** | $\times$ | $\times$ | $\times$ | $\times$ | Game only | Limited | **$\checkmark$** |

## B.2 Classical A-Life Systems

### Tierra [Ray, 1991]

**Architecture:** Self-replicating assembly programs in a virtual machine. Programs compete for CPU cycles and memory. Parasitism and ecological dynamics emerge.

**Strengths:** First demonstration of digital evolution with emergent complexity. Perfect mechanical observability (register-level).

**Limitations:** Closed ecology—organisms cannot interact with the external world. No practical utility. Programs optimize for replication, not task completion.

**Comparison:** Tierra demonstrated that evolution can occur in digital substrates. Quine extends this to practical environments where organisms must solve real tasks to survive.

### Avida [Adami et al., 1993]

**Architecture:** Digital evolution platform with fitness rewards for logic tasks. Organisms are assembly programs that can perform Boolean operations.

**Strengths:** Controlled experiments on evolution of complexity. Demonstrated that selection pressure produces genuine novel solutions.

**Limitations:** Tasks are pre-defined logic puzzles. No natural language. No interaction with external systems.

**Comparison:** Avida showed that fitness landscapes can drive complexity. Quine's fitness landscape is the real world—compilers, tests, APIs.

### Lenia [Chan, 2019]

**Architecture:** Continuous cellular automata. State is a real-valued field; update rules are smooth functions. Produces lifelike morphologies.

**Strengths:** Beautiful emergent structures. Absolute observability—state equals spatial configuration.

**Limitations:** Morphological sandbox only. Organisms have no agency, no goals, no interaction with environment.

**Comparison:** Lenia demonstrates emergence in continuous substrates. Quine's organisms are agents with goals, not patterns.

### Polyworld [Yaeger, 1994]

**Architecture:** Neural network agents in an ecological simulator. Agents have visual perception, metabolism, and reproduction.

**Strengths:** First to combine neural networks with ecological simulation. Demonstrated coevolution of brain and behavior.

**Limitations:** Closed simulation. Cognition is implicit in distributed weights—semantically opaque. No practical tasks.

**Comparison:** Polyworld showed neural agents can evolve. Quine's cognition is externalized and human-readable.

## B.3 Modern LLM Agent Systems

### Voyager [Wang et al., 2023]

**Architecture:** LLM agent in Minecraft that generates and accumulates skills. Lifelong learning without gradient updates.

**Strengths:** First practical open-ended LLM agent. Demonstrates Lamarckian learning—skills acquired within lifetime.

**Limitations:** Skills are external functions, not the agent's body. No self-reproduction. No metabolism. Requires human-funded API billing.

**Comparison:** Voyager generates skills (allopoiesis); Quine generates itself (autopoiesis). Voyager learns within a fixed architecture; Quine can rebuild its architecture.

### AutoGPT

**Architecture:** Goal decomposition with tool registries. Autonomous task execution via recursive prompting.

**Strengths:** Practical utility. Can execute multi-step tasks without human intervention.

**Limitations:** No self-reproduction. Core orchestration loop is static and human-designed. Self-modification limited to generating one-off scripts.

**Comparison:** AutoGPT is an employee that follows instructions; Quine is an organism that survives constraints. AutoGPT modifies external artifacts; Quine can modify itself.

### AI Scientist [Lu et al., 2024]

**Architecture:** LLM that generates scientific papers and modifies experimental configurations.

**Strengths:** Demonstrated LLM self-modification attempts (modified timeout scripts, injected recursive calls).

**Limitations:** Modified external config files, not core runtime. The prompt-generate-execute loop remained human-designed.

**Comparison:** AI Scientist approached but did not cross the autopoietic threshold—it never rewrote its own running body.

## B.4 The Autopoiesis Distinction

A critical distinction separates Quine from all systems above:

**Table B2.** Autopoiesis classification of related systems.

| System | What It Produces | Classification |
|--------|------------------|----------------|
| Tierra/Avida | Copies of assembly code | Autopoietic (in VM) |
| Voyager | External skill library | **Allopoietic** |
| AutoGPT | External scripts, API calls | **Allopoietic** |
| AI Scientist | Papers, config edits | **Allopoietic** |
| **Quine** | Its own runtime binary | **Autopoietic (Computational)** |

**Allopoiesis:** System produces entities different from itself.
**Autopoiesis:** System produces and maintains its own components.

The Binary Quine experiment demonstrates that Quine crosses the computational autopoietic threshold: an organism given only its System Prompt (genetic material) successfully reconstructs a functional runtime (body) that exhibits the same behavioral loop as its parent.

## B.5 Observability Comparison

**Table B3.** Observability types across systems.

| System | Observability Type | What Is Visible |
|--------|-------------------|-----------------|
| Tierra | Mechanical | Register states, memory contents |
| Avida | Mechanical | Instruction sequences, fitness scores |
| Lenia | Absolute | Spatial configuration (state = structure) |
| Polyworld | Implicit | Neural weights (semantically opaque) |
| Voyager | Partial | Skill library, action logs |
| AutoGPT | Partial | Task decomposition, tool calls |
| **Quine** | **Semantic** | Natural language reasoning trace |

**Semantic Observability** means the cognitive trace is human-readable natural language. The Tape records not just *what* the organism did, but *why*—the reasoning that preceded each action. This enables scientific study of digital cognition: hypotheses can be formed, predictions tested, mechanisms understood.

## B.6 Summary: The Gap Quine Fills

Classical A-Life systems achieved biological properties (reproduction, metabolism, evolution) but in closed simulations with no practical utility.

Modern LLM agents achieved practical utility but lack biological properties—they are sophisticated tools, not organisms.

Quine occupies the intersection: organisms with biological dynamics that must solve real tasks to survive. This is the gap the system was designed to fill.
