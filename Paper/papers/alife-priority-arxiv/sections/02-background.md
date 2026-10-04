# 2. Background

Classical A-Life systems (Tierra, Avida, Lenia) demonstrated that biological dynamics—reproduction, metabolism, evolution—can emerge in digital substrates, but they remain *closed simulations*: organisms evolve to self-replicate, not to accomplish practical work [Ray, 1991; Adami et al., 1993]. Modern LLM agents (Voyager, AutoGPT, AI Scientist) achieve practical utility but lack *genuine mortality*: they cannot truly die from resource exhaustion; time is reversible; failure has no irreversible consequences [Park et al., 2023; Wang et al., 2023]. While modern systems can programmatically instantiate "swarms" or sub-agents, these are merely allopoietic task delegations, not biological reproduction.

Table 1 summarizes the gap:

| Criterion | Traditional A-Life | Modern LLM Agents | **Quine** |
|-----------|-------------------|-------------------|-----------|
| Self-reproduction | $\checkmark$ (Code copying) | $\triangle$ (Logical API orchestration) | $\checkmark$ (OS-level process mitosis & sporulation) |
| Metabolism | $\checkmark$ (CPU cycles / energy units) | $\times$ (Agent-invisible) | $\checkmark$ (POSIX resource bounds) |
| Bounded lifespan | $\checkmark$ (Cycle limits) | $\triangle$ (Iteration limits) | $\checkmark$ (Resource mortality) |
| Heredity | $\checkmark$ (Genome copying) | $\triangle$ (External memory / RAG) | $\checkmark$ (Wisdom inheritance) |
| Practical Utility | $\times$ (Closed simulation) | $\checkmark$ (Task completion) | $\checkmark$ (Survival-driven tasks) |
| Semantic Observability | $\times$ (Registers) | $\checkmark$ (CoT traces) | $\checkmark$ (Causal Tape) |

*Table abbreviations: VM = Virtual Machine; RAG = Retrieval-Augmented Generation; CoT = Chain-of-Thought.*

Quine occupies the intersection: organisms with biological dynamics that must solve real tasks to survive, with semantically observable cognition. Appendix B provides detailed comparisons with each system.
