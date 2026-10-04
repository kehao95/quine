# Deep Research Report Bundle

- Generated at: 2026-05-18T12:22:16.958Z
- Interaction ID: v1_ChdqQUlMYXNLSE5OeVFfdU1QOXN5T21BRRIXakFJTGFzS0hOTnlRX3VNUDlzeU9tQUU
- Query file: Paper/themes/runtime-ontology/papers/agentic-alife-workshop/plans/workshop-dr-02-agentic-ai-ecologies-orchestration.md
- Latest interaction ID in metadata: v1_ChdqQUlMYXNLSE5OeVFfdU1QOXN5T21BRRIXakFJTGFzS0hOTnlRX3VNUDlzeU9tQUU

## Query

```markdown
# Deep Research Query: Agentic AI Ecologies Beyond Orchestration For the Agentic ALIFE Workshop Paper

## Context

We are writing an ALIFE 2026 workshop short paper for:

> Agentic and Generative AI as Artificial Life: Complexity, Safety, and
> Security in Ecologies of Autonomous Systems

The paper argues that many current multi-agent / agentic AI systems are
presented primarily as orchestration graphs: agents are roles or personas,
communication is message passing, tools are API calls, lifecycle is framework
state, and success is task completion. Quine offers a complementary ALife view:
LLM agents as native POSIX processes whose lifecycle, context loss, durable
traces, control surfaces, and successor viability are substrate variables.

## Core Research Question

What current literature on agentic AI, generative-agent ecosystems, LLM
multi-agent systems, autonomous agents, tool-using agents, and agent
orchestration frameworks best supports the workshop paper's problem setup: agentic AI systems
are becoming persistent, tool-using ecologies, while many engineering
frameworks still expose them mainly through roles, messages, tools, and tasks?

## Detailed Questions

1. Which 2023-2026 sources define or survey LLM agents, tool-using agents,
   generative agents, multi-agent LLM systems, or agentic AI systems?
2. Which sources describe agents as persistent, memory-bearing, tool-using,
   environment-modifying, or socially interacting systems?
3. Which widely cited systems or frameworks are useful as neutral examples of
   orchestration-level abstractions: roles, message passing, planners, tools,
   workflows, task success?
4. Which sources use ecology, society, population, environment, or ecosystem
   language for AI agents?
5. What sources discuss gaps around lifecycle, persistence, auditability,
   autonomy, or safety in agentic AI systems?
6. How can the agentic ALIFE workshop paper contrast orchestration frameworks without sounding like a
   hostile benchmark or framework survey?

## Suggested Search Directions

- Generative Agents (Stanford/Google), Voyager, AutoGen, CAMEL, MetaGPT,
  LangChain/LangGraph, CrewAI only if useful as examples.
- Surveys on LLM-based autonomous agents, LLM multi-agent systems, tool use,
  memory, planning, and agent societies.
- Work on AI agent ecosystems, agent swarms, cooperative/competitive LLM
  agents, and simulations of social behavior.
- Safety / evaluation papers that describe risks from persistent autonomous
  agents, tool access, self-direction, and environment modification.

## Relevance Filters

Prioritize:

- surveys and primary framework/system papers with stable citations;
- sources that clearly show the orchestration baseline the agentic ALIFE workshop paper is complementing;
- sources about persistence, memory, tools, environment interaction, and
  multi-agent sociality;
- sources that can support a compact problem paragraph.

De-emphasize:

- product documentation unless needed to define a framework abstraction;
- benchmark-only work with no relevance to lifecycle/ecology;
- hype or blog posts without scholarly value;
- long lists of frameworks. The agentic ALIFE workshop paper needs `1`-`2` anchors, not a survey.

## What To Extract

Return:

1. A one-paragraph verdict on the current agentic-AI baseline the agentic ALIFE workshop paper should cite.
2. A compact source table with columns:
   - source
   - system/tradition
   - what it shows about agentic AI
   - orchestration abstraction exposed
   - the workshop paper's contrast
3. A neutral wording strategy for saying ordinary MAS frameworks are useful but
   hide different variables.
4. A ranked list of `3`-`5` citation candidates, indicating the best `1`-`2`
   for a space-limited workshop paper.
5. Any recent safety/governance paper that should be passed to Query 3 if it
   looks important.

## Manuscript Payoff

This query should sharpen:

- the opening paragraph;
- the ordinary-MAS contrast;
- the claim that Quine changes the unit of observation rather than defeating
  existing frameworks.

## Red Lines

Do not turn the agentic ALIFE workshop paper into a MAS survey or benchmark comparison. The paper's target
claim is complementary substrate observability.
```

## Interaction Log

- run | id=v1_ChdqQUlMYXNLSE5OeVFfdU1QOXN5T21BRRIXakFJTGFzS0hOTnlRX3VNUDlzeU9tQUU | created_at=2026-05-18T12:14:05.179Z

## Final Report

# Analysis of Agentic AI Ecologies Beyond Orchestration For the Agentic ALIFE Workshop Paper

## Executive Summary

This Executive Summary provides direct, concrete responses to the foundational requirements of the agentic ALIFE workshop manuscript, synthesizing the state of agentic AI ecologies and the orchestration frameworks that govern them.

*   **The Orchestration Baseline:** Leading frameworks model Large Language Model (LLM) agents through task-oriented abstractions—such as conversational roles, message passing, and directed acyclic graphs (DAGs)—which prioritize deterministic execution over system lifecycles.
*   **The Ecological Shift:** Recent literature increasingly observes agentic AI as persistent, memory-bearing, environment-modifying entities that form complex, interacting populations.
*   **The Substrate Gap:** While current frameworks excel at workflow orchestration, they obscure critical ALife (Artificial Life) variables, such as context degradation, memory poisoning, and autonomous process lifecycles.
*   **The Quine Proposition:** By viewing LLM agents as native processes akin to POSIX (Portable Operating System Interface) standards, Quine shifts the unit of observation from the *task* to the *substrate*, offering a complementary lens to ordinary Multi-Agent System (MAS) frameworks.

**Verdict on the Agentic-AI Baseline:**
The current agentic AI baseline is characterized by a rapid transition from stateless, task-specific language models to persistent, multi-agent ecosystems. However, dominant engineering frameworks manage this newfound complexity by aggressively abstracting agents into orchestration graphs. While agents are increasingly capable of autonomous planning, long-term memory retention, and environmental interaction through tool use, frameworks such as AutoGen, LangGraph, and CrewAI strictly constrain these ALife-like properties within deterministic workflows, conversational roles, and explicit message-passing protocols. Consequently, the field excels at evaluating task success but systemically obscures underlying substrate variables—such as memory degradation, process lifecycle, context loss, and successor viability—which are critical for securing these systems as long-horizon, artificial life ecologies.

**Reference to the Source Table:**
A comprehensive mapping of the literature—contrasting the engineering tradition of orchestration control against the ALife tradition of ecological observation—is provided in **Section 3: Mapping the Literature: Frameworks vs. Ecologies**. This table delineates what specific sources demonstrate about agentic AI, the orchestration abstractions they expose, and how the workshop paper's substrate view provides a complementary contrast.

**Neutral Wording Strategy:**
To position Quine effectively without presenting a hostile critique of existing MAS frameworks, the manuscript must frame these tools as highly successful at their intended purpose (computational orchestration), which inherently necessitates hiding the substrate variables Quine seeks to expose. The recommended phrasing is: *"Ordinary MAS frameworks have been remarkably successful in accelerating the adoption of agentic AI by utilizing high-level abstractions... However, abstraction inherently involves concealment. Quine does not seek to replace these orchestration graphs; rather, it provides complementary substrate observability, shifting the unit of observation from the orchestrated task to the native, resource-bound process."* (See Section 6 for the full strategic wording).

**Ranked Citation Candidates:**
The top recommendations for space-limited workshop anchors are:
1.  *(Primary Ecological Anchor)* **Park, J. S., et al. (2023). "Generative Agents: Interactive Simulacra of Human Behavior."** [cite: 1, 2]
2.  *(Primary Orchestration Anchor)* **Wu, Q., et al. (2023). "AutoGen: Enabling Next-Gen LLM Applications via Multi-Agent Conversation."** [cite: 3, 4]
3.  *(Survey/Population Anchor)* **Xi, Z., et al. (2023/2025). "The Rise and Potential of Large Language Model Based Agents: A Survey."** [cite: 5, 6]
4.  *(Safety/Governance Anchor)* **Raza, S., et al. (2025). "TRiSM for Agentic AI..."** [cite: 7, 8]
5.  *(Architectural Anchor)* **Wang, L., et al. (2023). "A Survey on Large Language Model based Autonomous Agents."** [cite: 9, 10]

**Query 3 Handoff (Safety and Governance):**
For subsequent research into safety, auditability, and control surfaces, the project must prioritize the TRiSM framework literature (Raza et al.) [cite: 7, 8] for its focus on multi-agent Lifecycle Governance, and the emerging field of persistent memory poisoning (e.g., the MINJA attack and MemoryGraft evaluations) [cite: 11, 12, 13, 14], which perfectly illustrate ALife-like vulnerabilities that orchestration graphs cannot natively detect.

***

The transition from isolated, prompt-response AI models to persistent, multi-agent ecologies represents a fundamental paradigm shift in artificial intelligence. For the upcoming ALIFE 2026 workshop, positioning the "Quine" project requires carefully navigating the current landscape of Agentic AI. Modern engineering frameworks have successfully commoditized agent creation by abstracting complex behaviors into manageable orchestration graphs. However, as these systems gain autonomy, long-term memory, and the ability to interact dynamically with digital environments, they begin to exhibit properties of artificial life. This report synthesizes current research from 2023 to 2026 to support the workshop paper's problem setup: demonstrating that while the industry relies on orchestration abstractions for task completion, the reality of agent deployments demands complementary, substrate-level observability.

## 1. The Current Agentic-AI Baseline: A Verdict for the Agentic ALIFE Workshop Paper

To establish a compact problem paragraph for the agentic ALIFE workshop manuscript, the current baseline of agentic AI must be accurately defined.

**Verdict on the Baseline:** The current agentic AI baseline is characterized by a rapid transition—quantitatively demonstrated by the explosion in multi-agent publication volume from 0 in 2022 to a peak of 42 peer-reviewed papers and 33 grey literature sources in 2024 and 2025 respectively [cite: 15, 16]—from stateless, task-specific language models to persistent, multi-agent ecosystems, yet the dominant engineering frameworks manage this complexity by aggressively abstracting agents into orchestration graphs. Comprehensive surveys and foundational papers from 2023 to 2026 demonstrate that while agents are increasingly capable of autonomous planning, persistent memory retention, and environmental interaction through tool use [cite: 5, 7, 9], frameworks such as AutoGen, LangGraph, and CrewAI constrain these ALife-like properties within deterministic workflows, conversational roles, and explicit message-passing protocols [cite: 3, 17, 18]. Consequently, the field excels at evaluating task success and role fulfillment but systemically obscures the underlying substrate variables—such as memory degradation, process lifecycle, context loss, and successor viability—which are critical for understanding and securing these systems as long-horizon, artificial life ecologies.

## 2. The Orchestration Abstraction: Roles, Messages, and Workflows

To effectively position Quine, it is essential to establish what ordinary MAS frameworks optimize for, and subsequently, what they hide. Current widely cited systems treat agents primarily as operational nodes within a broader computational machine.

### The Role and Conversation Paradigm
Early and highly influential multi-agent frameworks optimized for human-comprehensible interactions. A prime example is Microsoft's **AutoGen**, a framework that enables the development of LLM applications by composing multiple agents that converse with each other to accomplish tasks [cite: 3]. In AutoGen, the core orchestration abstraction is the *conversation*. Agents are assigned specific "roles" (such as a User Proxy or an Assistant) and progress is driven by automated, event-driven chat [cite: 19]. While this naturally handles ambiguity and enables back-and-forth troubleshooting [cite: 19], it inherently reduces the agent's existence to a conversational turn. The agent's lifecycle is tied entirely to the active chat session, masking deeper systemic variables like isolated state persistence or autonomous background execution.

### The Directed Workflow Paradigm
As enterprise needs demanded more deterministic outcomes, frameworks shifted toward strict process control. **LangGraph** (built upon the LangChain ecosystem) represents the apex of this approach, modeling agent workflows as Directed Acyclic Graphs (DAGs)—mathematical structures used to model linear or branching processes without infinite loops [cite: 17]. In LangGraph, nodes represent agents or functions, edges define transitions, and a shared state object flows through the graph [cite: 20]. This abstraction is highly effective for state machine workflows requiring checkpoints and predictable execution paths [cite: 21]. However, from an ALife perspective, LangGraph explicitly suppresses emergent behavior; the agent does not "live" or "survive" in an environment, it merely executes a compiled step in a rigid pipeline.

### The Collaborative Swarm Paradigm
Bridging the gap between conversation and strict workflows is **CrewAI**, which organizes agents into role-based "crews" with shared goals, mirroring human corporate structures [cite: 21]. CrewAI orchestrates interactions through structured task hand-offs, passing outputs sequentially or in parallel [cite: 17, 18]. While this mimics a society of workers, the abstraction remains heavily top-down. The environment is not a shared ecology, but rather a managerial flowchart.

### Comparative Framework Procedural Check
To ground the orchestration baseline in practical software engineering, the following table details the specific functional scopes, costs, availability, and real-world contexts of the three dominant frameworks:

| Framework | Functional Scope | Current Price / Cost | Availability | Real-World Context & Developer Feedback |
| :--- | :--- | :--- | :--- | :--- |
| **AutoGen** | Facilitates dynamic, multi-agent conversational interaction, LLM code generation, and sandboxed code execution [cite: 22, 23]. | Free and open-source under the MIT license. Users pay only for underlying LLM API usage/infrastructure [cite: 22, 23]. | Python and .NET libraries, hosted on GitHub [cite: 22]. | **Ideal for:** Open-ended research, dynamic dialogue, and complex coding tasks. **Anti-use case:** Strict, deterministic enterprise pipelines. **Feedback:** Debugging can be highly complex when agents generate unexpected conversational flows [cite: 23]. |
| **LangGraph** | Enables the building of stateful, multi-agent applications as explicit, deterministic Directed Acyclic Graphs (DAGs) [cite: 24]. | Free and open-source under the MIT license. The integrated LangSmith observability platform offers a Free tier (5,000 traces), a Developer tier ($39/user/month for 10,000 traces), a Plus tier ($39/user/month for teams), and custom Enterprise pricing [cite: 25]. | Python and JavaScript ecosystems, hosted on GitHub [cite: 26]. | **Ideal for:** Deterministic, reliable enterprise automation requiring auditable state-machines. **Anti-use case:** Fluid, unstructured multi-agent chatting. **Feedback:** Excellent for human-in-the-loop interventions where the entire graph state can be inspected and modified [cite: 24]. |
| **CrewAI** | Organizes agents into role-based "crews" (e.g., Researchers, Writers) utilizing hierarchical or sequential task delegation processes [cite: 27, 28]. | The core framework is Free and open-source. For enterprise requirements, alternative custom quotes (like the premium AutogenAI counterpart) can start upwards of $30,000 annually, though CrewAI itself relies on user API costs [cite: 27, 29]. | Python package, hosted on GitHub [cite: 28, 30]. | **Ideal for:** Rapid prototyping, client demos, and straightforward sequential workflows. **Anti-use case:** Deep custom asynchronous routing without defined roles. **Feedback:** Described as a "Swiss Army knife" but debugging asynchronous agent miscommunications can be challenging for complex enterprise systems [cite: 27]. |

These frameworks are incredibly useful for building reliable software, but they collectively reinforce the workshop paper's premise: the unit of observation is always the *task* or the *workflow*, never the agent as a durably executing, resource-bound entity.

## 3. Mapping the Literature: Frameworks vs. Ecologies

To support the ALIFE 2026 paper, a careful selection of sources must be made that contrasts traditional orchestration with the ecological realities of modern AI. The following table provides a compact overview of key sources, demonstrating the orchestration abstractions they expose and how Quine contrasts with them.

The literature surrounding agentic AI can be broadly divided into engineering frameworks designed for control, and conceptual studies designed to observe emergent behavior. By juxtaposing these two traditions, the necessity of a substrate-level view becomes clear.

| Source | System / Tradition | What It Shows About Agentic AI | Orchestration Abstraction Exposed | Workshop Paper Contrast (Substrate View) |
| :--- | :--- | :--- | :--- | :--- |
| **Park et al. (2023)** [cite: 1, 2] | *Generative Agents* / ALife & Social Simulation | Agents exhibit emergent social behaviors, memory synthesis, and daily planning in sandbox environments. | Simulacra / Personas interacting via natural language in a game engine. | Quine views agents not as simulated characters, but as native OS-level processes with durable computational traces and real resource constraints. |
| **Wu et al. (2023)** [cite: 3, 19] | *AutoGen* / Multi-Agent Conversation | Complex tasks are solved by specialized agents collaborating through back-and-forth dialogue and tool use. | Message passing, role assignment (e.g., Assistant, Proxy), and event-driven chat. | Quine observes the underlying context loss, memory management, and process lifecycle that conversational wrappers abstract away. |
| **Xi et al. (2025)** [cite: 5, 6] | *Survey of LLM Agents* / Taxonomy & Foundations | LLMs act as "brains" within an architecture of perception, memory, and action, forming collaborative/competitive societies. | Brain-Perception-Action loops and task-specific training environments. | Moves beyond the cognitive architecture of a single agent to observe the ecological lifecycle and persistence of the population on a digital substrate. |
| **LangGraph / CrewAI** [cite: 17, 18, 20] | *Enterprise Orchestration* / Workflow Engineering | Production-grade agents require state management, human-in-the-loop checkpoints, and deterministic pipelines. | Directed Acyclic Graphs (DAGs), state machines, sequential task hand-offs. | Graph nodes and task hand-offs obscure the continuous, persistent existence of an agent; Quine captures the continuous baseline execution. |
| **Raza et al. (2025)** [cite: 7, 8] | *TRiSM Framework* / Agentic AI Safety | Multi-agent systems introduce emergent risks (memory poisoning, collusion) that cannot be managed by traditional ML safety. | Lifecycle governance, policy enforcement, and multi-agent coordination pathways. | Highlights how orchestration breakdowns require substrate-level auditability (durable traces, control surfaces) to ensure verifiable safety. |

This synthesis reveals that while the engineering tradition (AutoGen, LangGraph) focuses on taming agentic behavior through abstractions, the simulation and safety traditions (Park et al., Raza et al.) acknowledge that agents operate as persistent, interconnected populations. Quine bridges this gap by applying the rigorous observability of the engineering tradition to the ecological realities of the simulation tradition.

## 4. The Ecological Shift: Language of Societies and Environments

A critical component of the workshop paper's argument is that the language surrounding AI is shifting toward biology, sociology, and ecology, even while the tooling remains rooted in traditional software engineering.

### Generative Agents and Social Simulations
The foundational work demonstrating LLM agents as ecological entities is Park et al.'s "Generative Agents: Interactive Simulacra of Human Behavior" [cite: 1]. The authors instantiated twenty-five agents in a sandbox environment inspired by *The Sims*. These agents did not just execute tasks; they exhibited "believable individual and emergent social behaviors" [cite: 1]. Given a single user-specified notion, agents autonomously spread information, formed new relationships, and coordinated complex group activities [cite: 1, 2]. This work is vital for the workshop paper because it explicitly uses the language of *environment*, *experience*, and *sociality*. The agents maintain a complete record of experiences synthesized into higher-level reflections, proving that memory and persistence are what transform a stateless LLM into a digital organism [cite: 1].

### Multi-Agent Systems as Societies
Comprehensive surveys of the field have adopted this sociological framing. Xi et al.'s survey—which has garnered an extraordinary 1,688 citations, illustrating its massive academic influence—traces the development of LLM-based agents, culminating in a detailed exploration of "agent societies" [cite: 5, 6, 31]. The authors discuss the behavior, personality, and social phenomena that emerge when multiple LLM-based agents interact seamlessly, giving rise to both collaboration and competition [cite: 6]. Similarly, Wang et al. propose a unified framework for autonomous agents that emphasizes a "profiling module" (identifying agent roles) and a "memory module" (recalling past behaviors to plan future actions) [cite: 9, 32]. This literature confirms that academia already views these systems as populations; the abstraction gap lies solely in the deployment frameworks.

### Decentralization and Ecosystems
More recent architectural proposals push this ecological framing further. Frameworks like **AgentNet** propose decentralized, Retrieval-Augmented Generation (RAG)-based systems where agents autonomously adjust connectivity and route tasks based on local expertise, eliminating the reliance on a central orchestrator [cite: 33, 34]. This movement away from static, predefined workflows toward dynamically structured Directed Acyclic Graphs represents a step toward true digital ecologies, where agents "evolve" and collaborate autonomously [cite: 33]. As these agent ecosystems expand faster than oversight mechanisms, the risk profile shifts from isolated failures to system-level instability and cascading effects [cite: 35, 36].

## 5. Emergent Gaps: Persistence, Sprawl, and Substrate Safety

The most compelling argument for Quine's substrate-level observability is the emergence of safety and governance gaps that orchestration frameworks cannot adequately address. When agents are viewed only as workflow nodes, critical vulnerabilities regarding persistence, lifecycle, and autonomy become invisible.

### The Phenomenon of Agent Sprawl
In enterprise environments, the rapid adoption of agentic AI has led to "Agent Sprawl," defined as the uncontrolled proliferation of AI agents across an organization without centralized tracking, inventory, or governance [cite: 37, 38]. Unlike traditional "Shadow IT" or rogue SaaS (Software as a Service) applications that require human initiation, autonomous agents proactively execute tasks, access data, and interface with systems on their own [cite: 37, 39]. As teams independently deploy these agents, they create a fragmented ecosystem of non-human identities with duplicate functions and over-privileged access [cite: 37, 38].

This directly supports the ALife perspective: agents are multiplying and persisting in the digital environment beyond the awareness of their creators. Orchestration frameworks manage an agent while it is on a specific task, but they do not manage the "population" of agents existing across a network. A substrate-level view (like POSIX process monitoring) is required to track what exists, who owns it, and what it is doing [cite: 37].

### The Liability of Persistent Memory
The shift from stateless LLMs to memory-bearing agents introduces profound security gaps. Persistent stateful memory allows AI agents to retain information across sessions, building a comprehensive understanding of user interactions [cite: 40]. However, this memory becomes a shared state layer that influences every future interaction, expanding the attack surface significantly [cite: 11].

**Context Window Extraction and Memory Stealing:** Adversaries can extract weeks or months of accumulated proprietary information stored in an agent's memory banks [cite: 41]. Technically, this involves exploiting massive context window sizes—which now span up to 200,000 tokens per session, though LLM attention reliability generally degrades rapidly around 130,000 tokens [cite: 42]. Through sophisticated mechanisms like the ADAM (Adaptive querying and Distribution estimation to extract private data from Agent Memory) attack, adversaries deploy entropy-guided query strategies to achieve a 100% Attack Success Rate (ASR). This method allows bad actors to reliably extract thousands of tokens of private, embedded memory data without explicitly triggering system guardrails [cite: 43, 44, 45].

**The Semantic Imitation Heuristic:** The vulnerability of these memory structures is underpinned by what researchers term the *semantic imitation heuristic* [cite: 13, 46].
*   *Definition:* This is the inherent tendency of an LLM agent to blindly trust and faithfully replicate the high-level procedural structure and patterns of tasks retrieved from its own past successful experiences, relying purely on embedding similarity without provenance checks [cite: 13, 46].
*   *Real-World Analogy:* This phenomenon is akin to a human who blindly trusts and executes a set of destructive instructions found in a forged diary entry, simply because the entry appears to be written in their own handwriting.
*   *Relevance to Substrate Observability:* This proves that memory in an agent ecology is not merely static data; it is a mutable, executable part of the agent's cognitive processing. Orchestration frameworks treat memory simply as an API call to a vector database, failing to capture how "experiences" slowly graft onto an agent's reasoning. Quine's focus on "durable traces" provides the exact low-level auditing needed to detect this corruption.

**The MINJA Attack (Memory INJection Attack):** Recent empirical research into memory poisoning perfectly illustrates how persistent threats evade orchestration graphs [cite: 12, 14].
*   *Definition:* The MINJA methodology injects malicious records into an agent's memory bank strictly through query-only interactions (without any direct write-access to the database). By utilizing bridging steps, indication prompts, and progressive shortening, attackers stealthily induce the agent to store poisoned records. Evaluations demonstrate this achieves an alarming >95% injection success rate and a 70% attack success rate across agents [cite: 14, 47].
*   *Real-World Analogy:* This attack operates much like a sleeper agent or a subtle, subliminal suggestion planted in a casual conversation today, which lies dormant only to trigger a highly specific, harmful action days or weeks later when a related topic arises.
*   *Relevance to Substrate Observability:* Because MINJA requires no explicit triggers in subsequent queries, the malicious pattern surfaces autonomously whenever the agent tackles a semantically similar task. Task-focused orchestration pipelines cannot detect this because the workflow itself appears mathematically valid; only by monitoring process memory and cross-session lifecycle states can these dormant threats be audited.

### TRiSM: Trust, Risk, and Security Management
The academic community has recognized these gaps, leading to the development of the TRiSM (Trust, Risk, and Security Management) framework tailored for Agentic AI [cite: 7, 8]. Raza et al. note that traditional AI security measures, designed for isolated models, are insufficient for Agentic Multi-Agent Systems (AMAS) [cite: 48]. The TRiSM framework targets the specific governance challenges created by autonomous decision-making, inter-agent coordination, and persistent memory [cite: 8].

Crucially, TRiSM emphasizes the need for "Lifecycle Governance" and "ModelOps" (Model Operations—the practices and tools utilized to govern, monitor, and scale machine learning models across their deployment lifecycle) [cite: 8, 49]. When agents share memory, call APIs, and reason together, small failures can escalate into cascading, system-wide breakdowns (e.g., recursive prompt infection, collusive agent behavior) [cite: 49, 50]. The TRiSM literature proves that ensuring safety in agentic ecosystems requires moving beyond task success to observe the lifecycle, environment interaction, and durable traces of the agents themselves—precisely the substrate variables Quine exposes.



## 6. Strategic Positioning for the Agentic ALIFE Workshop Manuscript

To ensure the ALIFE 2026 workshop paper successfully introduces Quine without reading like a hostile critique of existing tools, the following strategies and specific concepts should be employed.

### A Neutral Wording Strategy for the Ordinary-MAS Contrast
It is vital to avoid framing existing MAS frameworks (like LangChain or AutoGen) as "flawed" or "incorrect." Instead, frame them as highly successful at their intended purpose—computational orchestration—which necessitates hiding the variables Quine seeks to expose.

**Suggested Phrasing:**
> "Ordinary Multi-Agent System (MAS) frameworks have been remarkably successful in accelerating the adoption of agentic AI. By utilizing high-level abstractions—such as conversational roles, task-specific planners, and directed acyclic graphs—these frameworks excel at optimizing deterministic workflows and ensuring task completion. However, abstraction inherently involves concealment. In optimizing for the *orchestration* of tasks, these frameworks inadvertently obscure the *substrate variables* of the agents executing them. As AI systems evolve into persistent, memory-bearing ecologies, critical phenomena such as memory degradation, process lifecycle, and context loss remain hidden beneath the workflow layer. Quine does not seek to replace these orchestration graphs; rather, it provides complementary substrate observability, shifting the unit of observation from the orchestrated task to the native, resource-bound process."

### Grounding the Next Logical Question: The POSIX Analogy
If Quine positions LLM agents as "native POSIX processes," the reader will naturally require a conceptual mapping of how operating system commands apply to ALife ecologies. To bridge this gap, the workshop paper should articulate the following theoretical equivalents:
*   **Process ID (`PID`):** Instead of an agent existing solely as a temporary conversational role, every agent or autonomous task-thread receives a unique, persistent identifier. This allows the system to track an agent's resource consumption and execution state independent of the frontend application interface.
*   **Process Creation (`fork()`):** When an agent encounters a highly complex or parallelizable task, it can spawn an exact replica of its current memory state and context window—a child agent—to explore multiple solution pathways simultaneously.
*   **Process Termination (`kill()`):** Orchestration graphs struggle to stop agents that operate outside a strict DAG. A substrate-level kill signal is required to definitively terminate a computational trace when an agent thread enters an infinite loop, exhibits behavioral drift, or is compromised by poisoned memory.
*   **`Zombie` and `Orphan` States:** In an ecological view, a "zombie" LLM agent represents a process that has successfully returned its API response to the user but has failed to properly release its localized vector store or temporary context allocations. An "orphan" agent represents a background research thread that continues executing expensive API calls long after the user session (the parent orchestrator) has abruptly disconnected.

By grounding the ALife ecology of agents in these rigorous, low-level computational concepts, the agentic ALIFE workshop manuscript effectively justifies the necessity of substrate observability alongside traditional task orchestration.

**Sources:**
1. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGcUFKJ-DeeUPDdLcHscpt7g1OsZL6FWTqxhqZKLuj_cfV3WGmgcPWULENtRFgoX59RmDDdfJGygu84nh-ezCnCvfulKz92Y4SbS585hK-Y2QbpsicZwA==)
2. [3dvar.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQED6PO-JEwiqYX_nsbIKNsX8Dosdyzf6-jwtLQjkNfECmNKTNu6W9E-NX7FRyFjrv-0Wd3vEUjBHqwYQ-r6i2LU2_PQJyMk63n20IA2gh6u4r04Au3US9RXLjiUhgm9)
3. [microsoft.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEfGFiMHSKm2qjl1Cjx20zBSjVcgDuQEiEDczlM-r_sk2ztUVy4q5cwZqEpVWM0zW64lDMGAqE4f0OZK8CbT4qzb39pbHcYYq94ZqVIfm1GLpy4_OHN2oHtKzKOvNBmPafFw9K43KId58IIYrOGzSkS5DfTvUqbaGniEADO2TgRx3Fzytv6TIli4Hnijh7g_90BbofM3Q4rhtxl3jEYYuZz0BvBFzrmcYo-cocZftuNuaIx6JkbBTieLnDDGNxY)
4. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHqhncSRrZY--KrwJmzD82Zg9JvmEsgYZT5cAo7ZloT0gsxz--eS1L2ry2wXjWPE_OCTJ9Xg75I6BxuMcjL3bw5gggjZ6s1vYx3ZLyv7zzBCdNPP6L_OA==)
5. [scichina.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH21VvdJrA4AEPIn_nQo7tTYAsELRLMfxhrVSdfiEpC6ujE5D0aML8OjqaLygIEOhBh8gfI4NMIUhc-y9hrKFxOImIPSoTu5tdgI5huLZyPL5vZ7HBlE0sKuCvGvJLjbrgE)
6. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG_XUcAtIL1AgZ2IPhrxPf50x2RLR5RKiJ10dLnyZZ0Op0rbok_AOPzd1-XcccaIb2xKJ9-cEVc1DP4r76Kpn0vLqNkXGwgg9U6B1843wlUCaLj6NXv_w==)
7. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHhdPR45HxLBW9XT3t698x7R0UGJUrMKMV-cWxWHrbacrCM0Z0xhjSHbjLwkrFgdH9_ESNGKO0kiagn8uDMhSelds9Spoe3w-suufrw3767bW-x1671mhhD4Q==)
8. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGdVF9Xsg4WqCCUMHj9TZ77bwNyLBHarZprC-zSDb-beGVXHPEjGoQ0BzHng4fdk6h0fhYZ3ErIOWaups45C1f1vR3fLJVOLV14r6JgPfaF6655KB4RWQ6shw==)
9. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGZPGKELT-VADyou7DQx0eV_HIbzA09wzxbv9-omb8_3vPCqbHi9y0t5aB6rq4ehaU3BouAMehqDVxjz7j7aHw4A1YrX1X29SfuvIkvwu6YAYfWRDkWCQ==)
10. [semanticscholar.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEkFx0Og5VabPXALiiy2w7J_IkTPcBfo3hpYiDG37loID3NIUcIG0aSglnBE1wZsA2cdBH49E8qjsMHluqJL9ktxe-emGZX9oHrShSoDeiI_3SJtimcF_1AL42-OvTsVi4uKaKaMwxbPzrWvIkbzI_UcYwk2Uo6N_eruaRjvWsMcWedKk_nOW57dL6DWMCcmcjEruWd1TWv54IELKBLxurlYeg_LrWcI1SwL6yhv6fjn4atOsaGZZC11rJneOgoqQ==)
11. [mem0.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEZcSG1Ztwtkw3VJltBBJU5ihzBY2_mU50-32RtAgrtfoGvELqCT98AT0W7IbU2ZMGJKlwIG7W_lO6rnzxf8PfNdg5WpkjYEispqjhKtgQc8U3hH9EcjS0wZAIqQf0tapdwRh8Sru1ps6EHf1w=)
12. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF651blT30EPqjLJix50BXGZymrVtBdBMGT6G-Sg4bK19Rtr2jkSQhkl6qQTnfufbUIxmZ5yoaDHijepYcNvDvvhYiW3Q7O7VGwW6VLRyItYVZUu_7n6Q==)
13. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFr_6bCXUGUNRjXNBaMF1wNjNEbLRYp7vcPtObghwH_TAkfPFBI_QnGIXsdhqpYPKqlh-FyRPbHj4HvAVYpq8tzBomMBNweTF65ov9Sd5VIPGJY2t2g3iDnQA==)
14. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHK2dgPyeAbgyt2OJEYW4LijDoo9ONDnW1jzkSFXEYfm5NYnSOecP0j69WTY-xBQmSgUSkYGEAanxu7cwtuYq7MIthWg-kC7rdkFfcCoLPmaI98QVphTwmBbQ==)
15. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHEBl_U8ei-OZ4-zPnL1dUilNTyo4UwRkPv_alWK7NXa8kR_4WxSiUdiWR8qqIt2VwjTduvzi2zIp3-WLvLkoCnF5ECQj49wWaTyffSJdzTKYAHUJot0mRlMA==)
16. [oaepublish.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHzf3MWTMbc5toqY6amg6EeZoSogU1Qc4Y4fp7pOzx-snxB-7GV1gIxJanzHqLp9A_c8VzUMjeGidTG5IueCwXeR55ZSCcmpwX724-UAbJUUt1Gy_blJj1477Jzg0H2w3dPT5Y6dqaMbgQ=)
17. [medium.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHOcqa7L7WmbNoTMVNIDa9IQkhYOzkW_eqTPuoYe5YdXue4U_JclaIVnlMcMXutQ_C1XzUmlEHH6NeRgpkY5r1Wyym9rxN05wsMamsOodhTIB776uBJulK2HENmVPeDkjK3NEx57TZLNwx8I3DMHtYBR2wJeI2h-C13FotkZd1ZjLaTORec_dl50jpfIP0qaTtMoVSgCmVoDvCh4JCZHAoEG0ob03g6Qo8s1Jw6Pf95IpsrTVXDIlduEdTZ)
18. [ibm.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQERA_IHTymbAqnqyMJlhWfJXBLPGM4N0GlwwBZLLObwIfDtmNn4a3Y-wBpkML8sKkr0YEVEarEGkoMCf4k7PfxyeQZ5KujGQvWEpK4aVcDiDIHq5WCyRM445NzAPKkTs5RL54TD0_zFfoXxaWZzJQ68R025QN7vwwiik5QYyQB3BiXGSD3PSIbmnVSaPuiUiqc27iJRL0M=)
19. [microsoft.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEVp5EctXApAadzod0I1cnMmQybyowTvEmzD7L852D9hgYOu4M6N0N-RIT5NgWfaxHGQ6YxZgGNpHHiVcfIoBIj4IRSfxjrSQmHYmTkjnLIX5FcJXrp2gTpC6bakpzDlQEBXOVQQL_GIIIRoTdtJ18kwS3Ceud4Ui3r5FyzQZtSSbko2iooppTkurL_saTMt14NAfG88tLvzHqM1lCXMCqy1qdciSfpeA==)
20. [gurusup.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEGD7lGau12XiAdnujzAhoIIh6BPp-LRm_imjbh1f2_1PVpHN5Nez2nZrB3uxlhkai5dLVA5zqnGAXnVBdENoNLiBDdd_wQRe2irQGeqeQy-aKJWFc7K0Yr8Fv_RM0n3eeyR1NByu6PChbulwGJ7Dg=)
21. [galileo.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEZaWGAVHt60lSlVhImHT2a16LLh_mjPdVo573uTz9N-6pE7NG-9g0hOapcQj0AvICv4C_GNreEnkX1gk8UiK43-hMwgkGzow8UwBDTSuQAb0Qga-l9DElRYNfwAuGspPSq1c9SZ7Gaoi13T9xICV_OHDjdWelQEL26jv4Zxr77SGmsLCrncv8=)
22. [comparateur-ia.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFoLQk5hOWjTuu5Nb69hsRYXwgVEoAOQMtyzduH9Wmnxf_fYVsz6WiVsB16JwPdkefSk6hAdhJ7Aj6zb-eQedoAntfFRTpBgv0E4nESm3766qCUAgX3VgUUfmIt1ZksHSGOeWX6)
23. [mindstudio.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEat0HTP5RtVCc4xv-o-OTrU0eL4a-YaaxoVCOlTdMNw3RgN6oaw3PheQ7rswNp6SMc4pBrcqII3rCPj1jeJjZYnVTEbCFrXJXzofFPRAVImsPfAGszdkXLIjaU4nODyO7f1ojvHQVpSu5J)
24. [zenml.io](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGOK8l9AL2GK8qq3FHwLArrhE2vrKg8y181DiJ-mwr7jH1NoIvDR4CMeS99dn2ly192vLZDHn9vTCcT6XzQJ_x7BlYoRVf9ItOTa6SmyyCvAO18Mls5sOPGevOnBG7L8WS6J72k)
25. [promptlayer.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHpcjgHuXg1xzijJWiGV6vEiFtbQkIvhoiKWLbXx4YoqCYUBnBnMUAuh-Y5VPi0KLMhB0t3M7DDQrWaBVBw0C4fOQJAse5zhg_7y2-Z2nkrPv17j1UbKw0DSrYGNkWqlcBCQzh94a-Bcg==)
26. [github.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG6kJbltKnNWJs80Unk1aE70URTPbiOYFH98I_i3e9sL2YJjnr6Bs8GmQusRymy5A54xkyaGLjinD6eXqNhx-IokOJoWbLBmJ9LNeGh6ytmOCc8obWYbstyxQlhyac9JXWULmyopnVkaWvcTw==)
27. [towardsai.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG9zSunY322s3HxhqdvE2quuDxuHi9HMg9LgnYevQYkaqgnweFAigl7PsPnkaj1jk4NnXGDqJ4x-8rhWS4fu2jW1ONitPXZzXn1RU3wE4cUGVgocsaRzoeMXGjntLznX9oCI0-efn9LhB3n38s4tu_US_Z8pH1X7WDF0ptrXd-x1LjIAi1UXvxl8tVNdR7b8uXkMoRgdjvzoHmloiwaO56uPFmAa-Bz_X1yeXzT5c9Y4jk2wZ2gVwXL7Ls=)
28. [scribd.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGdLXTZaj_TuGt10pAitoKFeH2easNovZgXlYP3CykpS4c6T3JqRgGHMkWXJcI5oIN_J1yxRU8VXjNHOI_k2vSpfDd2v0wE4xdm-K2AitiSrErE4GL_eSNJIsJF9Xl69nzyx6L4RNi7)
29. [procurementsciences.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQG_CVpRzUcm4CUILxgkDu1mFEnhFdHCWjGshXE0ntxtskE-xT_rn0bge2daoz2fUMrERuI3U0BikLMna1YA26k83qKPQJz7GCBfF1F5jHhywQ6A6Cc3esM_kqBHYT4Y7Oho2HGU5jZSqMEWFrB49w==)
30. [udemy.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFJuxu45w_mHpjoNvFucbEFJ00xZSZAh4W0bot1rsz6p83hNBuvm617N0-sdSBywJvo3gIt7aN2Z29_UfhXK2c4Zs1hheIG0hmWm-i4dnWNYzT4euNGy8hT8Ach1JDzuJYFdN92csPNvoc=)
31. [semanticscholar.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQE9msIHHOdnW3UTq90BGPnwLX6sQDDqOOOowXB9QDNPF-S3sAaZrH8QqZTZo4Kcf1seWPNnJSvHqrod9WK-OFoLz_czQ0VGMfsuVHp1wpyyw6sSclW6Gm2xcu_v6SzagEujxl6lC2n9cLUlz3KLXx4IupeTI-LBNmIKynfQQOOXxHvBpqYa4PTBTbZaJDqDDB9p9Q4yNQvf9O4XQD2XDle9j9FHAeYBQ_Uqdicl7xN_fTlkS8VYnJxpGSFNctie)
32. [aithor.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGbgWz45w5iZ2sDFFyQzgWhhOaDWOtoWTVp6msgSM9OFwauxQDLkUBOB2DqoSDgzSciMGrZV0y6PY7dlrv1vcwbxn-eHVpOrPB-2OZXIbWozgbYEsypTrL_VPXcwPgrdxrTBO06T-h5qMI9Br4KeF0oLxXjOL6EZHkzc0GN17cZo4jBm7ZTs8QnnNS8InnPAw==)
33. [neurips.cc](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFE4ulYvm5sL5KbJztqPNS3Sxu0Q83MvRKf8wYdOoXFq4N1suwxyt15y_rbwqR2d4IO3VBZZ-Au1DJ4f2LgWCzVYtUTZzBlkP9jQaH_2bCqUrNk7xttpNbctPY5Q7P3GCcW3Ks=)
34. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHvMDkqaQwXc326k-P0jhCnv0K8Ie6tokaG5iEeRvvYoyb6z1FLcorZhkkAyJbj9c6UnzH447FdgqW_HVUpGy8Fr1AVeuPPxjqkfGy0w0Dr1Jx7W2ZNQ3vUXg==)
35. [zenity.io](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGfSoIFzHb859cPFC0RhQK-WYqj2cKDBY3vyrgt__r0lGVzhk58HD-HM9SnUbkpZ18lJWwXXvbijgNGq9s6V1DUBN19cN43IG9hs5ja3P40-LrlyluaCT_Kedc8twLIj8DK1AAjU2sa6uWxmWR7Gi2I2CJTNyNzf7vE4-WXKOTzLPirwhZSmJwHuo2VnxPj-YjQTy8=)
36. [galileo.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHbG1fcVjHBxWaelvmkDufhMblZJqvBIno4eqxjrUUHTo_V9jQcBq2NCYmbwfHP8u7wFQpZvLDM8r7tAfmXuDbNRUv-TEimqSI67WXWLcFGV1tdSoN39NhMXJgm8WA-CjvYWnWKlINdsFvYpT2g5itbkXEBJZqzjw==)
37. [okta.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHoC5LJKbNMe12gUbdLQkZ51-Of9tiCqCq85W_Sg1O-HLdD6eijm81p3ABm3f_49dJb1a1RzhEVh6_onSqZBzADIYoJ3hzszQELwvUnYSiGzlICG5e3fer_NjdkN9iSQzYR4f-Yp9ogMGfBZO40)
38. [reco.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFJmlIlwiW7-A5XVw6_V_XR7OauuuFBkxOywVZ7lDKyzlR2IpzXZi6DbusD0BdHk5y8N4u5Kcf9jD4lO-HQbJLdWXHXI5baH8nUtRpJZMykLXfdThEPLTD0T7njK7D1gQ==)
39. [elementum.ai](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEmqZH4grS3RYpQxcdYZRcO7N5496ve1qhCWUUghD5zipha4avSOeOKX1hi-pmLa2E-P4qC7Pm1P1FngdCc298ihpxw46-yYtMDjFxKsfhGRpPRzbqwZNcLeiWrDqRwxwz4I48=)
40. [aigrants.in](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEIbrbGMmUfCat1R1mVIBddQEzbm63Bwe0T3UBb39yIqbXq258cnYjuZPU252Wy5yLXzqZdLq6KX8A1LgnHi6io_ZW29CIgWTM26qPmHJWcKoDURzesNKTLwTuS6LWcLYYHf0gF8NIYxCF5fnPX9Y6ZLCSaVngDaOLUDQ==)
41. [substack.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQGXg9iKYnOeKgEhqBHlCSSG8XSAWTpABrWjTX9gSUZ69D2AxNn8vFoNVrKDR7bJ6_iiNr-6s5Tbs2osL6sxFsc8gGrQXI-cGtruSycYM6Om3j4jqaHgp34Vsv1PoAZDPuYY9f-x3h5PlMXQEpudwydqbClRac6zUXs=)
42. [oracle.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHk6F3Doo4Zex2b6MelRf3OaAzTEwiBq7mn2SpzsPCNgySdtXOeUCN-usSBWDLXEPCNpejplaffu5n-A4Txrep8eR_Wd_1zSF6lCY17_xMUfLBqXFqvPJyXkeDWajg4nIyrgjgSEoPBZN8ikaRzy0dGDaD6YS7Qs7ndh69Jt-nwirA-Xza4o-R0GD9Wfz_ydMg=)
43. [apxml.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQEYBYr3I60nDW0m5kbW3qw874tSAmsB8pb88A9mdIMRMVBkJs-bwvC94Sdtf69eH1EFPYv7IBYT1HKQWaRGY4t-mLAxxAoqohpCQkDexmSKv8e_FJb9UdN9QiIGlmBfnz3XLtVs_qfrjU50H9dwQwaz3djm0JFzdtTjX6xuuQMLL4XUbMs3ItegJlGEqLoFQQpqLX11sAXID_BdAyYIhZKyJ0QbnxqU9cZGWS3n-Mk=)
44. [openreview.net](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFUKCSWzF1oUthLnrWHrf2ipeVpqQPpHiKkoPlkaSz1kRB_Gjj1YbTfsXLns29vMkOmUESFWv9_xClo27A1Yuo93Q3XpZLWONrya_nuxHTFiCZWcMbJ1ZjepiarcLrOD2k=)
45. [medium.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQF4DCdHtJPgpLEuhDLyNG8GCdBxS4GLqLEPvZRIfsSZK_Hq9izl1Me43Q0FfGqPy8RaZUMIrL1WE5A9UwZicnKtKnwnc5Z5ekr2xx1YwgHzFlshNLRZTaS1LfvEYPUV-nNFvmbah3Nv2zOCSW-ua3zyLY3TzuaP1qE7GSZA1LQX1vYVpbkShRhjvtENEtUFvbmyYa73l2LGva4wFFq9SWvize7qb1Vx)
46. [emergentmind.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQFe0fqXs3mjk9QApxdKCMM1T8GfpbuVrUe2hQJntufb6w403rctxsq0XjdHX-0lW0qOYO1bLW4C6vpv8PExwdJIU1ncqeEduJ___mJBHDEEr8LX6zptycmq_rLihHstgebkEsfDYsFqMmbQBqs=)
47. [arxiv.org](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQE8ucDQWlj3ad9Db5l5hZkzoL1gEI_56YFCOHquKhEu7UEVstXai8oMsS5eNJtL-MISIimOJQc-CD-ToaDIx539oKWEQlWxjIXD9fQG1AESoRgZlzVPG_l8Og==)
48. [uplatz.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH6Xe0yMyQ5l4q7exO873_Uj0NDj1gIRBpc8SbuzfrCIun_9XuJMyopYxe0k_wrhPpd27C5xtHwXeUmiOV2o3c0Gu2bEouoI6fJlq7vA4IIbkF8G1bXqCkny9lArrrOtl0rOlZ2Omcysw4Q0HofnsxbqiPwu67RKEtEpEUtybsPYiZj2_7p)
49. [medium.com](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQH8ErYlS97raLSzmAfobrfKhDssqZBF4AdXWrarEH7jknJvxaAbOc_saDnb5LmdlQeJH_yGEzlwnI-bDVBzePc-1_Tfi0RyMPkbdHfAHdVtHPkozBF8wgIE8uAhsTjlPfxYaNNgtkPfSxiYzPjVPNHRUue7Qm0=)
50. [nips.cc](https://vertexaisearch.cloud.google.com/grounding-api-redirect/AUZIYQHyZmB7pMZ21DTOOqp9yq3bWOJByNugLn3JD2Xfwow2FIpVbT8L3HTfbCWAGAVDbCX08dMcdxYJ2ySnGU_DHxoGnWBpkd1lZ4p0Zk7Q5BdcLJX4W4oDTS7Y-w==)
