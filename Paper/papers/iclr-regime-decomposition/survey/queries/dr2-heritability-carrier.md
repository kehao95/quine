# Context

We are developing a theoretical framework for AI self-replication that distinguishes six regime variables. One dimension (D2) concerns the "heritability channel": what carries information from a predecessor AI agent to its successor. We argue that under certain conditions (barrier-mediated succession), the operational carrier — not source code, weights, or full system bundle — becomes the primary heritability channel. We need to survey existing literature on AI agent property inheritance, successor formation, and identity continuity to position this claim.

# Research questions

1. Has anyone formally studied "heritability" in AI agents or AI systems? That is:
   - What properties of an AI agent survive across re-instantiation, re-training, fine-tuning, or reconstruction?
   - What properties survive when an agent is rebuilt from partial information vs. full source?
   - Is there any formal or information-theoretic treatment of what makes an AI agent "the same agent" across instantiations?

2. Has anyone studied "successor formation" in AI — the process by which a new agent instance takes over the role, obligations, or operational position of a previous agent?

3. What is the state of AI agent identity continuity research?
   - Personal identity in AI / numerical vs. qualitative identity
   - Continuity across deployment, re-deployment, version updates
   - Memory and obligation transfer between agent instances
   - Agent "lineage" tracking

4. Has anyone applied information-theoretic or rate-distortion arguments to:
   - Agent property preservation
   - What can and cannot be reconstructed from partial descriptions
   - Channel capacity arguments for agent replication/succession

5. Is there any work on "carrier" vs "source" distinction in AI replication?
   - Deployment artifacts (config files, system prompts, tool schemas, memory stores) as heritability channels
   - What gets copied vs. what gets reconstructed vs. what gets lost
   - Operational environment as part of the agent's identity

6. What related concepts exist in adjacent fields?
   - Software deployment and configuration as "inheritance"
   - Agent-oriented software engineering: agent lifecycle, agent cloning, agent migration
   - Digital organism / artificial life heritability models
   - Memetic transmission vs. genetic transmission (cultural evolution analogs)

# Scope

- 2018-2025 for AI-specific work
- Include foundational earlier work (artificial life, agent systems, philosophy of mind) if directly relevant
- Include both formal/theoretical work and empirical studies
- Include safety-relevant work on agent persistence, agent takeover, agent succession

# Deliverables

1. Literature map: what exists vs. what's missing for each research question
2. Key papers with their specific claims about AI agent heritability/identity/succession
3. Identification of the closest prior work to our "carrier as heritability channel" concept
4. Assessment: is the source-vs-carrier distinction novel, or has someone made it before?

# Constraints

- Distinguish between: (a) work on copying/cloning entire agents, (b) work on what partial information suffices for reconstruction, (c) work on identity/continuity across re-instantiation
- Flag any work that already uses terms like "carrier," "heritability channel," "successor formation" in the AI context
- Be precise about whether work is formal/theoretical or empirical

# Output requirements

Structured markdown report with:
- Executive summary
- Literature map (table: topic × key papers × what they claim × gap relative to our framework)
- Closest prior work analysis
- Novelty assessment
- Key references with URLs/arxiv IDs
