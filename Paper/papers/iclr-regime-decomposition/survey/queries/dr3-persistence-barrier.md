# Context

We are developing a theoretical framework with six dimensions for AI self-replication. Two dimensions concern: (D1) persistence vs. succession — whether the "same agent" persists or a new agent succeeds it; (D5) barrier-free vs. barrier-mediated replication — whether the agent can freely copy itself or must overcome execution barriers (sandboxes, access controls, compute limitations). We need to survey the theory frontier on both.

# Research questions

## D1: Persistence, succession, and delayed addressability

1. How does the AI safety/alignment literature treat agent persistence?
   - Is there work on "future self" or future-agent addressability in AI systems?
   - Has anyone studied the difference between an agent persisting vs. being succeeded by a new instance?
   - What is the state of AI agent lifecycle theory?

2. In multi-agent AI systems, has anyone studied:
   - Generational succession (agent A creates/delegates to agent B)
   - Obligation transfer across agent instances
   - How agents reason about or address their own future instantiations
   - Deferred goals, delayed execution, future-conditioning behavior

3. Is there any formal treatment of:
   - When a successor counts as "the same agent" vs "a different agent"
   - Identity criteria for AI agents (type vs token, program vs process)
   - Persistence conditions: what must be preserved for identity to hold

## D5: Execution barriers and barrier-mediated replication

4. What is the state of sandboxing / access control as a barrier to AI self-replication?
   - Sandbox escape literature as it relates to replication capability
   - Compute access as a barrier to replication
   - Permission models and their effect on agent proliferation

5. Has anyone studied how barriers change WHAT replicates (not just WHETHER replication succeeds)?
   - Does barrier presence change the heritability channel?
   - Does barrier-mediated replication produce different successor profiles than barrier-free?
   - Any work on "what gets through" barriers vs "what doesn't"

6. What is the theory of "barrier-mediated" processes in:
   - Biology (membrane transport, reproductive isolation, speciation barriers)
   - Software security (privilege escalation as barrier crossing)
   - Information theory (channel capacity under constraints)

# Scope

- AI safety, AI systems, multi-agent systems literature: 2020-2025
- Foundational work on agent identity and persistence: any date
- Sandbox/security literature as it relates to autonomous AI: 2020-2025
- Adjacent fields (biology, philosophy, software engineering) only if directly cited

# Deliverables

1. Two literature maps (one for D1, one for D5) showing existing work and gaps
2. Key papers for each dimension
3. Assessment: is the persistence-vs-succession distinction novel in AI? Is the barrier-as-heritability-channel-modifier concept novel?
4. Any formal or mathematical treatment of either dimension

# Output requirements

Structured markdown with:
- Executive summary (split D1/D5)
- D1 literature map + key papers + novelty assessment
- D5 literature map + key papers + novelty assessment
- Cross-cutting observations
- Key references with URLs/arxiv IDs
