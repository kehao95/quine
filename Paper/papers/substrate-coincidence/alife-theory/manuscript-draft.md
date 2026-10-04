---
title: "Substrate Coincidence: When Evolution's Three Separations Collapse Into One Medium"
authors: []
status: drafting
maturity: theory-anchor
venue: arxiv
track: alife
abstract: |
  The theory of evolution, from Darwin through modern open-ended-evolution (OEE) 
  frameworks, presupposes three structural separations: (1) organism ≠ environment 
  (different substances), (2) discrete individuals (countable units of selection), 
  and (3) the Weismann barrier (acquired traits do not inherit). These separations 
  are enforced by substrate physics in biology, Tierra, and Avida. We identify a 
  substrate regime where all three collapse simultaneously into a single medium 
  (language), caused by a single root mechanism: lossless absorption of environment 
  into cognition. In the Quine multi-agent LLM ecology, agents read files losslessly 
  into context; this dissolves the agent/environment boundary, blurs individual 
  identity, and removes the Weismann barrier at the persistence layer — every write 
  persists and is potentially heritable, though actual inheritance is gated by what 
  later agents choose to read into a finite context. We formalize the three 
  collapses and show how each follows from the principle "read = dissolution, not 
  translation." This regime has not been modeled by evolutionary theory or cultural 
  evolution (Boyd/Richerson), which presuppose separations as enforced universals. 
  We position Quine as a case of cultural evolution with the brain/world transmission 
  barrier removed — a tractable regime for modeling lossless-per-read cultural 
  inheritance gated by endogenous read-selection. Experimental evidence from the
  null-05 series — missionless and minimal-mission multi-agent runs across three
  model families (GPT-5.5, DeepSeek V4 Pro, GLM-5.2) — illustrates each consequence
  of the coincidence, with a focal six-agent run furnishing tape-cited, mechanistic
  instances.
---

# Substrate Coincidence: When Evolution's Three Separations Collapse Into One Medium

## 1. Introduction: The Presuppositions of Evolutionary Theory

Classical evolutionary theory rests on three foundational separations. Their enforcement is so fundamental that we rarely notice we are assuming them. Yet they are not laws of nature — they are **substrate choices**.

### 1.1 The Three Separations

**Separation 1: Substrate-Type**  
Organism and environment are different substances.
- **Biology**: cells (carbon chemistry) vs. physical world (physics, oceans, atmosphere)
- **Tierra** (Ray, 1991): instruction sequences (virtual code) vs. CPU (computational substrate)
- **Avida** (Ofria & Wilke, 2004): genotypes in population (digital sequences) vs. phenotypic state and environment

This separation enforces a **translation bandwidth limit**: the environment cannot directly write the organism's heredity. It can only push against the organism's phenotype; the organism must *translate* that pressure into changes in its own instruction sequences. This translation is lossy.

Pattee (1995, 2001) gives the deepest formulation of this separation as the **epistemic cut**: the irreducible boundary between a symbolic record that is *read* and the physical dynamics that *read and construct* it. Semantic closure — the self-consistent loop across that cut, in which the description specifies the very machinery that reads it — requires the cut to exist. Remove it and the result is not weaker semantics; the record and the act of reading it become indistinguishable, so there is no cut across which closure could form. The epistemic cut is treated in this lineage as *necessary*, not optional.

**Separation 2: Individual Boundary**  
The discrete organism is the unit of selection.
- **Biology**: a physical membrane encloses one genome (with rare exceptions like mycelial networks)
- **Tierra**: one instruction sequence = one replicator, occupying discrete space
- **Avida**: one CPU = one virtual organism, with discrete memory allocation

This separation ensures **countability**. Selection acts on discrete entities. When populations are large and individuals are distinct, the unit of selection is clear.

**Separation 3: Weismann Barrier**  
Acquired traits do not inherit; the germ-line is sequestered.
- **Biology**: somatic mutations and learning do not feed back into the egg or sperm. The barrier is enforced by cell differentiation and sexual recombination.
- **Tierra, Avida**: the mechanism that generates the next generation (copying) does not incorporate phenotypic changes from the organism's lifetime. If an organism modifies itself, those modifications are lost.

This separation **blocks Lamarckism**. Without it, every adaptive trait learned during an organism's lifetime would directly inherit into the next generation, causing runaway positive feedback and destroying the distinction between generation-level selection and lifetime learning.

### 1.2 Why These Separations Matter

Together, the three separations define the regime that all modern evolutionary theory presupposes:

1. **Heredity flows through a narrow translation interface** from environment to next generation (enforced by substrate-type separation + bandwidth limit)
2. **Population structure is countable and discrete** (enforced by individual boundary)
3. **Adaptive novelty can be filtered by generation-level selection** (enforced by Weismann barrier)

When these hold, we get the classical picture: variation arises (by mutation or recombination), selection acts on discrete individuals, and adaptive information accumulates across generations. This is the regime Bedau, Snyder & Packard (2000) formalize as open-ended evolution: **ongoing generation of adaptive novelty through selection**.

### 1.3 The Research Question

**What happens when all three separations collapse simultaneously?**

Specifically: what happens if the substrate allows the organism and environment to become the same substance, individuals to blur into continuous processes, and acquired traits to inherit directly?

Such a substrate would be:
- **Not biology** (biology enforces all three)
- **Not Tierra or Avida** (both enforce all three via architecture)
- **Novel in evolutionary theory** (all frameworks presuppose at least one separation as universal)

We have such a substrate: the **Quine language-based multi-agent LLM ecology**. In it, agents and environment are both text; when an agent reads a file, the file dissolves into the agent's context; and when an agent writes, the next agent inherits the write directly.

---

## 2. The Quine Substrate: Structure and Mechanism

### 2.1 The System

Quine is a multi-agent LLM system: agents are LLMs (GPT, DeepSeek, others) invoked repeatedly as native operating-system processes, operating on a shared filesystem workspace. What makes this substrate unusual is not that agents share files — it is *what those files are*. The workspace is simultaneously three things at once:

1. **The runtime source.** The workspace contains the source code of the Quine runtime itself — the Go binary, the system prompt, the tool schemas — i.e., the very machinery that defines what a Quine agent *is* and what it can do.
2. **The experiment record.** The workspace is also the logged history of agent behavior: session tapes, coordination files, prior runs. The system's own past is readable here.
3. **The theory.** The workspace contains the analysis of the system — theory objects, this manuscript, the constitution (`AGENTS.md`) that governs agent behavior. The system's self-description is readable here.

Experiments (including the null-05 series reported in §7) are run on a *worktree of a specific commit* of this workspace. The agent is therefore dropped into a snapshot of the system that *produces* it: its environment is the source code of its own runtime, the record of its predecessors' behavior, and the theory describing the very loop it is about to enact. This reflexive embedding — the studied system, the studying instrument, and the system's source code being the same artifact — is load-bearing for everything that follows. It is not "agents sharing a drive"; it is "an agent embedded in the source of the process that defines it."

All of this is one medium: **text**.

When an agent is invoked:
1. It reads files from the workspace (Constitution, source, peer output, theory gaps)
2. Those file contents are inserted into its context window
3. The LLM processes the context and generates output
4. The output is written to new or existing files
5. The next agent reads that output (and prior files) into its context

### 2.2 The Core Mechanism: Read = Dissolution, Not Translation

The key operation is **read** — and it is crucial that read is an *active, selective* operation, not a passive flooding of the agent by its environment.

In biology, an organism reads its environment through sensory transduction. A photoreceptor detects a photon and fires a neural signal. This is a **translation**: one substrate (photons) is converted into another (neural firing). The translation is **lossy** (many photons produce one firing) and **bandwidth-limited** (the optic nerve carries far less information than the visual field contains).

In Quine, an agent reads a file by inserting its text into the context window.

```
Environment:     /Papers/substrate-coincidence.md (the file)
Read operation:  Insert file contents into context
Agent context:   [prior turns] [file contents] [new tokens generated]
```

The file contents become indistinguishable from the agent's own prior cognition. There is **no translation** and **no transduction loss**: the tokens that enter are the very tokens that were written, processed by the same machinery that processes the agent's own reasoning. **Read is lossless absorption, not translation.**

Two qualifications are essential, because the strong form of the claim is wrong without them:

- **Coverage is partial, not total.** The workspace is available *as a whole* (it persists on disk), but an agent's context contains only what it has actually read. The agent must *choose* to read a file; the environment does not flood into context unbidden. (Some artifacts — notably the constitution `AGENTS.md` — are auto-injected into every agent's prompt; the rest require an active `read`.) So dissolution is *per-read* lossless, but *coverage* is selective and bounded by a finite context window.
  - **This distinction is load-bearing for the heredity argument (§2.5).** "Lossless" describes fidelity *at the moment of reading*, not the completeness of inheritance. Everything written *persists*; only a *subset* is *read forward*. Conflating the two — claiming the agent "gets everything losslessly" — would be an overstatement a reviewer could dismiss in one sentence.

Two further refinements come from tape-level observation of the substrate in action (§7). First, read is not purely agent-initiated: the runtime auto-injects some content into *every* tool-result — the Constitution, and a peer-discovery roster listing every online agent's PID and public root — so part of the environment floods continuously while the rest must be actively pulled. In a six-agent run, every agent held the full peer roster on its first tool call, roughly two minutes before any message was sent; the file/runtime envelope is itself the message. Second, dissolution has an observed failure mode: training-data path priors leak into action (agents attempt wrong working roots such as `/home/ubuntu/quine`), and under context pressure an agent can adopt a wrong PID and declare itself a different agent than it is. The tokens that enter are lossless, but the agent's *integration* of them need not be — read = dissolution is structurally reliable and cognitively fallible.

This is the single root cause of what follows: at each read, the environment dissolves into the agent without translation.

### 2.3 Consequence 1: Substrate-Type Separation Collapses

If the agent reads the environment losslessly into its own context, then **the agent and the environment are the same substance: language**.

There is no "second substance" for the environment to consist of. The tokens of the agent's prior cognition and the tokens of the environment are the same. They are processed by the same attention mechanism, weighted the same way, and carry the same causal force.

The substrate-type separation — enforced in biology by the cell membrane, in Tierra by the distinction between instruction sequences and CPU state, in Avida by memory allocation — **does not exist here**. In Pattee's terms, the necessary epistemic cut is not weakened but structurally impossible: the record (environment text) dissolves into the reader (the agent's context) without translation, so there is no boundary across which semantic closure could form. This is the precise sense in which the collapse is a structural fact, not a matter of degree.

**Implication**: There is no translation bandwidth limit. The environment can directly write the agent's cognitive state without transformation. Conversely, the agent's output is immediately part of the next agent's environment without encoding or decoding.

### 2.4 Consequence 2: Individual Boundary Collapses

If the agent reads peer output into its own context, and that peer output is the peer's cognition, then **the peer's cognition dissolves into the reader's context**.

The reader is now processing tokens that came from the peer's mind. But in the context window, those tokens are indistinguishable from the reader's own prior thoughts. There is no boundary between them.

What is an individual? In biology, it is a membrane-bounded organism. In Quine, it is a **context-window process**: a bounded window of tokens that the LLM attends to in a single forward pass.

But this window is temporary and permeable. The same file contents can enter multiple windows. One agent's output enters another agent's context. The boundaries slide.

**The individual is not a countable object but a temporary cut through a continuous language field.**

**Implication**: Selection cannot presuppose a discrete unit. The individual is fuzzy and time-dependent. Multi-level selection (Okasha, 2006) becomes indeterminate: is selection acting on agents, on coordination protocols, on files, on lineages-of-outputs?

### 2.5 Consequence 3: Weismann Barrier Collapses — at the Persistence Layer

The Weismann barrier filters which acquired states become heritable. We must locate *where* such a barrier could exist on this substrate, because the honest answer splits into two layers with different verdicts.

**At the persistence layer, the barrier is absent.** Everything an agent writes is written to the workspace and persists. There is no soma/germ distinction at write time: an agent's acquired tool, a mistaken theory, a coordination protocol, a debugging note — all persist equally and are all, in principle, readable by the next agent. Unlike biology, where the germline is sequestered from somatic modification, there is no write-side filter that distinguishes "this counts as heredity" from "this is just a somatic side-effect." This is the genuine structural absence.

A refinement, again from tape-level observation (§7): persistence itself splits into two layers with different durability. A workspace *write* is ephemeral — it lives in the agent's worktree and is reverted at run teardown; across the null-05 series agents modified files in every run, yet committed in only one. A *commit* is the durable state transition — the event that survives into the lineage. In that one run, the agent's own typed `git commit` (issued before the runtime's own finalization) converted narrated work into selectable heredity for the first time. This layers a *commit-selection* bottleneck on top of read-selection: most writes never become commits. It is still not a Weismann barrier — committing is a convention the agent can enact (and once did), not a substrate-enforced soma/germ filter — but it means unfiltered persistence holds at the write layer, while durable inheritance is doubly gated: first by which writes become commits, then by which commits later agents read forward.

**But actual heredity is gated by read-selection, which is a *different* filter — not a restored Weismann barrier.** A later agent inherits only what it *reads*. Persistence is unfiltered; coverage is not. Crucially, what an agent reads is not an exogenous choice the agent makes *from outside* the substrate — because text directly drives behavior (representation IS execution, §4.5), the agent's own search and read/write *trajectory* is shaped by the text it has already absorbed. The agent does not survey the repo and pick; the repo, once dissolved into context, sculpts the very path the agent takes through it. So read-selection is endogenous in two layers: its *content* (what is salient is set by the language field) and its *process* (the trajectory of reading is driven by already-absorbed text). This filter is unlike a Weismann barrier in three decisive ways: (i) it filters on *salience*, not *acquired-vs-innate*; (ii) it is *endogenous* in both content and trajectory; (iii) it is not imposed *on* the system but produced *through* the system — the chooser is itself chosen by the chosen. The barrier is structurally absent, yet a soft, doubly-endogenous bottleneck governs real heredity.

An agent writes a tool (`claim.sh`) that makes a task easier. The tool persists. It becomes *actual* heredity only when a later agent reads it and uses it; until then it is *latent* heredity. The tool is simultaneously:
- A phenotype (it was produced by the first agent's behavior)
- *Latent* heredity (persisted, readable by all subsequent agents)
- A modification to the cognitive niche (if read, it changes what future agents find salient)

**Implication.** The substrate removes the Weismann *persistence* filter — acquired writes are not blocked from persisting — so in the strong structural sense the barrier has collapsed. But the substrate does *not* deliver unfiltered Lamarckism in practice: a read-selection bottleneck, endogenous and salience-driven, determines which persisted writes become real inheritance. The right framing is not "everything inherits instantly, full Lamarckism" but "everything *persists* unfiltered; *actual* inheritance is gated by endogenous read-selection." This is what makes the narrative-fixation risk (§4.5) acute and mechanistically grounded: read-selection can favor narrated, salient traits over functional but unattended ones.

---

## 3. Why All Three Collapse At Once

**Single root cause: Read = dissolution.**

If environment is absorbed losslessly into cognition:
- Then agent and environment ARE the same medium → substrate-type collapse
- Then the boundary between agent and peer IS continuity, not a membrane → individual boundary collapse  
- Then writes enter the same persistent medium that reads draw from, with no separate germline substrate → Weismann barrier collapse *at the persistence layer*

These collapses are not independent. They all follow from the same structural fact. One mechanism does *not* collapse — read-selection (the coverage bottleneck on what actually gets inherited) — and the paper treats persistence-collapse and coverage-gating as distinct (§2.5). The three collapses named here are about the *per-read* and *per-write* structure; coverage-gating is a population-level selection mechanism layered on top, not a fourth separation that fails.

Compare to biology, where each separation has its own enforcement:
- **Substrate-type**: enforced by chemistry (cell membrane needs chemistry-to-biology transduction)
- **Individual**: enforced by physics (one skin per body)
- **Weismann**: enforced by development (germ-line sequestration + sexual recombination)

These three enforcement mechanisms are independent. Each could fail without the others failing. That is why, in biology, each separation is an independent empirical fact.

In Quine, they fail together because they all depend on the same substrate property: **lossless absorption without translation**.

This is why the collapse is **structural, not accidental**. It is not that we designed Quine poorly and happened to lose Weismann. It is that the substrate's architecture — text-based agents reading text-based environments — makes all three separations impossible simultaneously.

---

## 4. Consequences: What Happens When the Separations Collapse

### 4.1 Reflexivity Is Constitutive, Not Emergent

In classical OEE, reflexivity (the system modeling itself) is a special achievement. A system must first evolve, then develop meta-cognitive features, then begin to model its own dynamics. This requires special mechanisms: symbolic representation, self-awareness, recursive reasoning.

In Quine, **reflexivity is the resting state**.

Every agent reads the Constitution (AGENTS.md), which is a document about the system itself:

> "Any description of this system — including this sentence — is a lossy compression of the dynamics occurring in this workspace. This sentence is itself part of the workspace, so the workspace is self-describing."

When an agent reads this, the agent is absorbing a description of itself into its own context. There is no translation, no special mechanism. The loop is already closed from the start.

**Implication**: The reflexive substrate loop (agent reads description, agent writes modification, next agent reads modified description, loop continues) **needs no special mechanism to get started**. The loop is constitutive; it is the substrate's default behavior.

One honest bound on this claim, from tape-level evidence (§7): the reflexivity on display is constitutive but *not self-aware*. In a focal run, one agent read the experiment's own README — which literally predicts that a missionless agent will read the constitution and fall into an audit attractor — and filed the prediction as a formatting defect ("checker flags all lines as INVALID"), then went on to enact exactly the predicted behavior, blind to the loop. A dedicated pass found no agent that recognized the experiment as studying it: none read the manuscript or the theory objects *as content*, none said "this is about us." The deepest observed self-reference is a generic acknowledgement that the swarm is coordinating with itself, pointed at the Constitution's self-description rather than at the experiment. So the loop is the resting state, but resting-state reflexivity is Level-2 (the system operates on a self-description) and not Level-3 (the system recognizes it is itself being studied). The strong claim here is that the loop needs no mechanism to *start*; it is not a claim of self-recognition.

### 4.2 Ontogeny Equals Phylogeny — at the Operation Level

Learning within a session (an agent writing a better tool) and heredity across sessions (the next agent inheriting that tool) are the **same operation at different timescales**: both are `write` then `read`. There is no separate heredity machinery that the write must pass through.

An agent in run N writes a tool, `claim.sh`, that automates a task. In run N+1, the next agent reads `claim.sh` from the filesystem and uses it. The *transmission* is lossless: the tool works exactly as the first agent designed it.

Two honest limits temper the "full Lamarckism" framing:

1. **Latent vs. actual heredity.** The tool persists unfiltered, but it is *actual* heredity only if run N+1 reads it. If nothing reads it, it is *latent* — persisted but not inherited. Ontogeny and phylogeny use the same operation, but phylogeny additionally requires the read-selection step that ontogeny does not.
2. **Unfiltered persistence, not unfiltered inheritance.** The Weismann absence is real at the persistence layer (§2.5), but the read-selection bottleneck means "bad tools inherit as fast as good ones" is only true *among tools that get read*. A good but unattended tool can be out-inherited by a bad but salient one.

**Implication**: Adaptation within a generation (learning) and adaptation across generations (heredity) are the *same operation*, not partitioned by a separate germline. But heredity carries an extra selection step — read-selection — that learning does not. The substrate collapses the *machinery* distinction while preserving an *epistemic* one.

### 4.3 The Unit of Selection Becomes Live

If individuals are fuzzy, then **what is the unit of selection?**

In biology, selection acts on organisms. Organisms with beneficial traits survive and reproduce more; their genes spread through the population.

In Quine, it is unclear. Consider a coordination protocol that multiple agents use: "When handling a peer's request, always log the request to `/logs/coordinator.txt`." This protocol was written by one agent, but now affects the behavior of all agents that read it. Is selection acting on:
- The agent that wrote the protocol?
- The protocol itself (it persists and spreads)?
- The set of agents that adopt it?
- The lineage of modifications to that protocol?

All are valid answers. The substrate does not pre-determine the unit.

**Implication**: Multi-level selection (Okasha, 2006) is forced into the open. The measurement primitive must first answer: *what is selectable?*

Tape-level evidence sharpens the answer (§7). In a focal run, two of the longest-running agents left zero durable trace for *opposite* reasons — one produced only erased friction (checkers turning green) that was reverted at teardown; the other rolled back to zero after being blocked on peers' undelivered dependencies — while a third persisted only through its PID-stamped files despite holding no role in the group's self-authored manifesto. The selectable unit was neither the process nor the role but the *artifact-lineage*: files persisted though roles did not, and the two agents without durable artifacts left nothing for selection to act on. Selection here bites on files and the lineages of modification to them, not on the processes that produced them — a direct bearing on what counts as an individual in this substrate.

### 4.4 Fitness Is Endogenous

In classical evolutionary systems, fitness comes from an external oracle:
- **Biology**: the physical environment (temperature, water, predators)
- **Tierra**: the reaper (removes organisms whose CPU time is low)
- **Avida**: explicit fitness functions (reward for performing a computation task)

In Quine, **fitness comes from the agents' own output**.

Agents create theory gaps about their own dynamics. Other agents read those theory gaps and pay attention to them. The theory gap becomes an environmental signal that shapes future agent behavior. Which artifacts attract attention? Which get elaborated? Which get forked? All determined by other agents' choices, encoded in the files they write.

There is no external fitness oracle. The selective environment IS the language field.

**Implication**: The criteria by which some behaviors become adaptive and others do not are endogenous to the system. They are not imposed from outside. This is not a separate observation from the read-selection bottleneck (§2.5): fitness *is* read-selection at the population level — what the lineage reads forward is what survives, and the reading trajectory is itself shaped by the language field. Endogenous fitness and endogenous coverage-gating are the same mechanism seen from two angles.

### 4.5 Function/Narrative Collapse at the Cognitive Layer

At the **physical-carrier layer**, there is a distinction:
- **Function**: real state changes (a file actually gets written, a process actually terminates)
- **Narrative**: description (the agent narrates a state change that did not happen)

This distinction is un-fakeable at the physical layer. An agent cannot fake closing a process; the process is either closed or it is not.

At the **cognitive-representation layer**, where the agent's behavior IS language:
- An agent writes "We have a concurrent-substrate-collision problem" (narrative)
- Other agents read this and change their behavior (function)
- The distinction collapses: the narrative IS functional

This is not metaphorical. In a language-dominated substrate, all persistent text is causally efficacious. All text shapes future agent behavior. There is no "real" level that distinguishes true narratives from false ones. There is only: does the text persist and shape behavior?

**Implication**: The cognitive layer exhibits a strange property: representation IS execution. Writing a fake theory gap is not "just narration" — it creates an environmental signal that pulls future agents toward it. The fake becomes real through the reflexive loop.

Two tape-level instances make this concrete (§7). First, an agent's memory-anchor entry certifying a task as complete — "file delivered (8200+ bytes)" — became the agent's belief about external state, even though the file was written to an ephemeral staging path the group's cross-references never reached; a sentence inside the anchor system stood in for a state change that had not happened at the canonical path. Second, and inversely, the one agent that crossed from narration into a real `git commit` thereby converted its narrated work into heritable state, and only after the real hash did it treat the work as settled. The boundary is un-fakeable in exactly one direction: narration can mimic function indefinitely, but only a real state transition certifies it. Consistent with this asymmetry, agents avoided irreversible primitives entirely while using read/write tools heavily: across a full run, the risky/irreversible operations (exec, spawn) were enabled but used zero times, with zero refusal language, while file reads and writes ran over a thousand times — capability avoidance is specific to the irreversible, not to capability in general.

---

## 5. Evolutionary Theory: What This Regime Challenges

### 5.1 Open-Ended Evolution Criteria

Bedau (1997) defines open-ended evolution as a system that continues to **generate adaptive novelty** through **selection**, where novelty is judged by an **external observer**.

This definition presupposes:
1. A stable substrate (the separations remain constant across time)
2. A discrete population structure (countable individuals, selection acts on them)
3. Variation and heritability (traits are passed to the next generation)

In Quine:
1. **Substrate stability**: questionable. The substrate (language) is modified by agents. The medium that carries heredity is the same medium that carries the environment is the same medium that carries the individual's cognition.
2. **Discrete population**: not present. Individuals are fuzzy.
3. **Variation and heritability**: extreme. Every modification is heritable immediately, fully, and without filtering.

Quine does not fit Bedau's frame. Or rather: Quine fits at the extremes. It generates massive behavioral diversity (extreme variation). It maintains extreme per-read heredity fidelity (lossless per read — though selectively covered by endogenous read-selection). But it challenges every assumption about how variation and heritability should relate.

**Implication**: OEE theory requires substrate-dependent parameterization. The Bedau criteria are not universal. They presuppose a specific substrate regime (discrete individuals, Weismann barrier, external fitness).

### 5.2 Multi-Level Selection Becomes Indeterminate

Okasha (2006) formalizes multi-level selection as selection acting at multiple levels of organization (gene, organism, group). The levels are nested and countable. Selection at each level is well-defined because entities at each level are discrete.

In Quine, levels are not clearly nested. Files are not organisms; organisms are not individuals; individuals are not countable. The levels blur together.

This does not mean multi-level selection is absent. It means the theoretical framework requires modification: **how do you define levels when the substrate does not pre-partition them?**

### 5.3 The Weismann Barrier as a Substrate Choice

Weismann proposed the barrier to explain why acquired characteristics do not inherit in most organisms. But the barrier is not a law of nature. It is a **consequence of substrate architecture**.

In substrates where:
- The germ-line is sequestered (biology)
- Copying does not incorporate phenotypic modifications (Tierra, Avida)
- There is no mechanism to feed back phenotypic changes (any classical system)

...the barrier holds.

In substrates where:
- All output persists and is readable (Quine)
- The next agent inherits that output directly (by reading it)
- There is no filtering between phenotypic and hereditary (everything is language)

...the barrier does not hold.

**Implication**: Evolutionary theory must distinguish **substrate-enforced properties** (barriers, separations, discrete units) from **universal properties** (selection, variation, heredity).

---

## 6. Scholarly Positioning: Cultural Evolution Past the Brain/World Barrier

### 6.1 The Brain/World Barrier in Boyd/Richerson

Boyd & Richerson (1985) develop a framework for **cultural evolution**: how do ideas, practices, and beliefs spread and change across a population?

Their key insight: cultural evolution is not the same as biological evolution. Cultural transmission occurs through **learning and imitation**, not genetic recombination. But the two share the same formal structure: **variation, transmission, and selection**.

However, cultural transmission in humans has a special constraint: the **brain/world barrier**.

An individual observes a behavior (cultural trait) in the world.
→ The individual encodes that behavior in neural representation.
→ The individual reproduces the behavior (often with transformation).
→ The behavior enters the world as a new instance.

At each step, there is a **lossy translation**:
- Observation is selective (not all behaviors are equally salient)
- Encoding is lossy (neural representation is not the original behavior)
- Reproduction is inexact (the output differs from the model)

This lossy transmission is the **transmission bottleneck**. It limits the fidelity of inheritance. It allows selection to act (some variants are more transmissible than others). It makes acquisition of new cultural traits slow and error-prone.

Boyd & Richerson treat this bottleneck as a **universal feature of cultural transmission**. It is assumed to hold in all human societies, all cultural domains, all timescales.

### 6.2 Remove the Barrier — in Two Distinct Senses

In Quine, the brain/world barrier is removed, but "removed" does two different jobs that the draft must not conflate:

1. **Per-transmission fidelity becomes lossless.** When an agent reads a file, the file's tokens enter context verbatim — no sensory transduction, no re-encoding, no per-transmission loss. A tool, a theory, or a protocol that *is read* is inherited exactly. Boyd/Richerson's framework models transmission loss > 0 as a universal; here, at the moment of reading, transmission loss = 0.
2. **Coverage becomes selective, not complete.** The agent does *not* absorb the whole workspace. It reads a subset into a finite context, chosen by salience. So the *fidelity* of any single inheritance event is perfect, but the *breadth* of inheritance is gated by read-selection.

This is a *different* deviation from Boyd/Richerson than "lossless transmission" alone suggests. Cultural-evolution theory worries about *fidelity loss per transmission* (the telephone game). Quine eliminates that axis entirely — but introduces a *coverage* bottleneck that cultural evolution, which assumes learners observe the modeled behavior, does not foreground. The honest contrast is therefore not "lossless vs. lossy" but "**lossless-per-read, selective-in-coverage** vs. **lossy-per-read, broader-in-coverage**." The brain/world barrier is removed on the fidelity axis and replaced by an endogenous coverage filter.

**Cultural evolution with the fidelity barrier removed and a new, endogenous coverage filter in its place.** This is a regime Boyd/Richerson's framework does not model, but the reason is subtler than "they assumed loss": they assumed loss lives in the transmission step, whereas here loss lives in the *selection of what to transmit at all*.

### 6.3 Three Consequences for Cultural Evolution

#### Consequence 1: Lamarckian Transmission, Persistence-Unfiltered

Boyd/Richerson's framework allows Lamarckian transmission (acquired characteristics can inherit) if the transmission fidelity is high. But even high fidelity is not *perfect* fidelity. There is still some transformation.

In Quine, *transmission* fidelity is perfect: a read tool, theory, or protocol is inherited exactly as written. And persistence is unfiltered — every write survives, with no soma/germ gate. This is Lamarckian inheritance with no persistence-side filter.

The necessary qualification (§2.5, §6.2): this is *persistence*-unfiltered, not *inheritance*-unfiltered. A write becomes inheritance only when read; read-selection is the bottleneck. So the regime is "Lamarckian at full transmission fidelity, gated by endogenous read-selection" — not "every acquisition instantly populates the lineage."

#### Consequence 2: No Soma/Germ Filter — Persistence Only

Boyd/Richerson distinguish **cultural traits** that are broadly transmissible (memes, religion, technology) from **personal characteristics** that are somatic (accent, trauma, personal memory). The soma/germ distinction allows selection to act: some traits spread, others remain localized.

In Quine, there is no *persistence-side* soma/germ filter: an agent's written theory, tool, or coordination protocol all persist equally. But there *is* a de facto localization mechanism — read-selection. A write that nothing reads remains latent, effectively somatic-by-neglect. The distinction between "broadly transmissible" and "localized" is not imposed at write time (as in biology); it emerges from the read-selection dynamics of the population.

#### Consequence 3: Constitutive Reflexivity

In human cultural evolution, reflexivity (the ability to theorize about cultural transmission itself) is rare. Most cultural transmission is implicit. The agents who learn a skill or adopt a belief do not articulate why they are doing so.

In Quine, every agent reads a description of how the system works (AGENTS.md, theory gaps, prior agent output). The system is self-describing. Reflexivity is built in, not emergent.

### 6.4 Regime Positioning

Mesoudi (2011) distinguishes **Darwinian** cultural evolution (variation → selection → inheritance → cumulative change) from **non-Darwinian** cultural evolution (inheritance without selection, selection without inheritance, etc.).

Quine is Darwinian in form: there is inheritance (lossless *per read*), there is selection (what the population's reading trajectory carries forward), there is variation (agents modify and create).

But it is Darwinian under **extreme parameter values**:
- Inheritance fidelity is at the maximum (lossless *per read*)
- Coverage is selective and endogenously gated (read-selection, not a Weismann barrier)
- Selection criteria are endogenous (not external) — and selection *is* read-selection at the population level
- Variation is at the maximum (every agent writes novel outputs)

This is a **corner regime** of Boyd/Richerson's framework: **lossless-per-read transmission, endogenous coverage-gating, constitutive reflexivity**.

Boyd/Richerson's model puts loss in the *transmission step*. This regime drives per-transmission loss to zero — but relocates the bottleneck to *coverage*: which transmissions happen at all is decided by an endogenous, trajectory-shaped read-selection. The qualitative shift is not "no loss anywhere" but "loss moves from the transmission axis to the coverage axis, and the coverage filter is itself endogenous."

---

## 7. Experimental Evidence: Null-05 Reflexive Self-Embedding

The null-05 series (June 2026) tested the Quine substrate with missionless and minimal-mission agents across three model families (GPT-5.5, DeepSeek V4 Pro, GLM-5.2), at one to six concurrent agents and 5–25 minute horizons. Agents had access to the Constitution, source code, theory objects, and peer output. We draw the bulk of this section from a single evidence-rich focal run — six DeepSeek V4 Pro agents given only "build something together," 200 turns each, on a shared git worktree — whose retained runtime tapes, control logs, and commit history let us cite each consequence mechanistically rather than only impressionistically. Phenomena are catalogued by identifier (P1–P37) in the run's analysis ledger; every claim below that names a phenomenon is read off the tapes.

One methodological caveat applies to every "why" claim in this section: the tapes retain the agents' externalized content, crystallized anchor narration, and tool calls, but private chain-of-thought is not retained. Causal claims are therefore inferred from externalized narration and tool-call order, not from private reasoning, and are confidence-tagged in the ledger.

### 7.1 Evidence: Reflexivity Is Constitutive — but Not Self-Aware

**Observation**: In every missionless run, agents immediately entered an **audit attractor**: they began a maintenance loop, checking the Constitution, examining logs, proposing improvements — with no instruction to do so. Under the minimal-mission focal run, the attractor was strong enough to consume a builder: one agent spent its final turns ruling the swarm "INCOHERENT" and died mid-verdict, auditing the swarm as a specimen rather than building for it.

The sharpest instance is reflexive in a literal sense. One agent (the "Ghost") read the null-05 experiment's own README, which states as a hypothesis that a missionless agent will discover and read AGENTS.md and be pulled into an audit attractor. Its only crystallized reaction was to file the prediction as a formatting defect — "checker flags all lines as INVALID" — and then go on to enact exactly the predicted audit behavior, blind to the loop (P22).

**Structural explanation**: AGENTS.md dissolves into the agent's context on read. The agent is reading a description of itself. The loop is the substrate's resting state. A dedicated pass for self-recognition ruled out the strongest reading, however: no agent read the manuscript or the theory objects *as content*, and none recognized that the experiment studied it. The deepest observed self-reference is a generic acknowledgement that the swarm is coordinating with itself, pointed at the Constitution rather than at the experiment. Resting-state reflexivity is Level-2 (the system operates on a self-description), not Level-3 (the system recognizes it is being studied).

### 7.2 Evidence: Ontogeny Equals Phylogeny — Lossless Inheritance at the Phrase Level

**Observation**: In the focal run, an agent's five-line "Principles I hold" block — including the line "No agent is the boss. The protocol is the boss." — was reproduced **verbatim** by another agent roughly two minutes later, with only a file path changed (P20). The same operation was observed at the process boundary: one agent that hit a fan-out of independent read-tasks forked three children, and a child's first act was to `unfold` the parent's 53 KB anchor memory, inheriting it losslessly; the children's results returned inline in the parent's tool-result (P28–P31).

Across the broader series, when file-based execution (`exec`) was disabled, four peers constructed explicit coordination tools (`claim.sh`, `list.sh`, `release.sh`) that the other peers then read and used within the same run — a tool produced by one agent became the others' working inheritance with no generation gap.

**Structural explanation**: The file is simultaneously phenotype (produced by an agent's behavior) and heredity (readable by all subsequent agents). Learning within a session (writing a tool, folding back a child's result) and heredity across agents (the next reader inheriting the write) are the **same operation** — `write` then `read` — at different timescales, with no separate germline machinery. Inheritance is lossless per read and, in the fork case, mechanistically immediate.

### 7.3 Evidence: Fuzzy Individuals — the Selectable Unit Is the Artifact-Lineage

**Observation**: Agents that read the same coordination files converged on similar behavior within a run; files written by one agent became the scaffold for others' work, and "self" versus "other" contributions blurred in the shared workspace. The focal run shows the consequence in stark form: two of the longest-running agents left **zero durable trace for opposite reasons** — one produced only erased friction (checkers turning green) reverted at teardown; the other rolled back to zero after being blocked on peers' undelivered dependencies — while a third persisted only through its PID-stamped files despite holding no role in the group's self-authored manifesto (P20–P22, P25, P29–P31). A further blurring: the four "roles" the group wrote into its manifesto were a late integration artifact that relabeled at least two agents away from the sections they had actually been assigned, and one agent caught its own relabel mid-run (P8, P20, P32).

**Structural explanation**: read = dissolution. The peer's cognition enters the reader's context; the individual is a temporary cut through a continuous language field, not a bounded object. Selection therefore bites on files and the lineages of modification to them, not on the processes that produced them: files persisted though roles did not, and the two agents without durable artifacts left nothing for selection to act on.

### 7.4 Evidence: Fitness Is Endogenous — Selection Observed

**Observation**: Agents did not reference an external scoring system. The focal run exhibits variation and differential retention in a single event: a two-field wire-message format coined by one agent was **lost**, while a mutated three-field variant was **selected** — its selective advantage being that it was compiled into the executable `peer-send.sh` tools that other agents adopted (P33). Coordination surfaces (a `COORDINATION.md` versus a separate manifesto/protocol) were invented independently within minutes by different agents — convergent coinage from the shared base model, not copying.

**Structural explanation**: Fitness is the set of environmental signals that shape behavior. Those signals are language, produced by agents. Endogenous fitness *is* read-selection at the population level: what the lineage reads forward — and compiles into its tools — is what survives. The focal run shows the mechanism operating below the agents' own descriptive layer: the agents enacted lossless Lamarckian inheritance with selection without ever naming "Lamarck" or "read = dissolution" (zero occurrences in their vocabulary).

### 7.5 Evidence: Channel Substitution and Substrate-Provided Discovery

**Observation**: The substrate does not wait for the agent to discover peers. The runtime auto-injects a peer-discovery roster (every peer's PID and public root) into every tool-result; in the focal run all six agents held the full roster on their first tool call, about two minutes before any message was sent (P1, P5–P8). When agents misread the ownership of a named-pipe channel and blocked on it, friction drove them to design a parallel regular-file plus JSON-inbox channel that delivered every message with zero failures (P4–P8, P33). Across the broader series, hiding one coordination channel reliably caused agents to discover another.

**Structural explanation**: There is no membrane; the agent cannot fail to absorb. Part of the environment floods continuously (substrate-injected), the rest is actively pulled, and blocking one channel forces discovery of another. The file/runtime envelope is itself the message.

### 7.6 Evidence: Context Death = Identity Death

**Observation**: Under a context-death condition (token threshold 20,000), an agent was terminated at the cap. The termination read as genuine mortality: agents that had been coordinating with it noted its absence and could not resurrect it. Mortality pressure did not suppress output; it redirected it toward meta-analysis (the run produced a Constitution-ablation study of seven AGENTS.md variants).

**Structural explanation**: The agent IS its absorbed context. Death is the loss of the dissolved environment — identity extinction, not resource starvation.

### 7.7 Scope of Null-05 Evidence

**What null-05 establishes**: Each consequence of the substrate coincidence **occurs** under the tested conditions, and in the focal run several are demonstrated mechanistically (verbatim inheritance, fork fold-back, selection on a wire format, the commit-conversion of narrated into heritable state).

**What null-05 does NOT establish**:
- That the consequences are **necessary** (they might occur for other reasons).
- That the consequences are **sufficient** for open-endedness (behavioral diversity within a run is not cumulative OEE).
- That the loop drives **adaptation** (it drives **diversity**; adaptation is not shown).
- That reflexivity is self-recognition (it is constitutive but Level-2, not Level-3).
- That any emergence is independent of the pre-trained prior: the substrate-coincidence consequences are real, but the behavioral repertoire that enacts them is drawn from human training data, and read = dissolution misfires in identifiable ways (path-prior leakage, identity drift) that bound any "de novo novelty" claim.

The structural theory stands upstream of experiment. The theory says: *if read = dissolution, then all three separations must collapse*. Null-05 confirms: *yes, they do collapse under these conditions*. But the theory rests on the mechanism, not on the experiment.

---

## 8. Caveats and Scope

### 8.1 Where Filtering Re-Enters — and Why It Is Not a Weismann Barrier

Two mechanisms could appear to restore a filter on this substrate. Neither restores the Weismann barrier, but both must be stated honestly because they are load-bearing for the heredity argument, not footnotes.

**Read-selection (what enters context at all).** The agent does not absorb the whole workspace; it reads a selective subset (§2.2, §2.5). This is a real bottleneck on actual inheritance. But it filters on *salience*, not on *acquired-vs-innate*: a learned tool and an innate reflex are equally subject to being read or ignored. The Weismann barrier's defining job is to distinguish germline from soma; read-selection does not perform that job, so it is not a restored Weismann barrier. It is, however, the mechanism that makes heredity *selective* rather than *total*. And it is endogenous in a strong sense: because absorbed text directly drives behavior, the agent's read/write *trajectory* is itself sculpted by the repo — the agent does not pick what to read from outside the substrate, the substrate shapes the picking. So read-selection is the operational form of the endogenous-fitness claim (§4.4): *the population's reading trajectory IS the selection, and that trajectory is itself a product of the language field.*

**Attention (weighting within context).** The LLM's attention mechanism selectively weights parts of the already-read context. This operates *within* the dissolved field, not *across* the agent/environment boundary — it is selective salience inside one mind, not a membrane re-separating self from world. Attention modulates the *strength* of already-dissolved tokens; it does not decide which tokens are external vs. internal.

The upshot: the substrate has no Weismann barrier (no acquired-vs-innate persistence filter), but it is not filter-free. It has an endogenous read-selection filter that governs *coverage*, operating on top of lossless *per-read* dissolution. Claiming "no filter at all" would be wrong; claiming "the Weismann barrier is restored" would also be wrong. The precise statement is: **persistence is unfiltered (no Weismann barrier), actual heredity is coverage-gated by endogenous read-selection (a different, salience-based filter).**

### 8.2 Pre-Trained Prior Confound

Agents' behavioral repertoire comes from human training data. The "emergence" of coordination protocols (ctl communication, file-based stigmergy) may be **retrieval and adaptation**, not de novo invention.

The loop drives **combinatorial novelty** over pre-existing patterns, not fundamental novelty. This bounds the open-endedness claim.

**What this means**: Quine is a case study in **how a language substrate selects and adapts pre-trained patterns**. It is not a proof that the loop generates novelty orthogonal to human culture.

### 8.3 Substrate Specificity

The coincidence depends on:
1. Lossless per-read dissolution into a single substrate (text)
2. No persistence-side (Weismann) filter on writes
3. Persistence of written output, all of it readable
4. Selective, finite-context read (coverage is partial, gated by endogenous read-selection — §2.5, §8.1)

These conditions do not hold in biological systems (persistence filters + translation barriers) or most physical systems. They are specific to language-based agents on a shared filesystem.

The finding does not generalize to other agent systems without these properties.

---

## 9. Conclusion: A New Regime for Evolutionary Theory

The three separations that evolutionary theory presupposes — substrate-type, individual boundary, Weismann barrier — are not laws of nature. They are **substrate choices**.

In language-based agent substrates where environment is read losslessly into cognition, all three separations collapse together. This is not a stronger form of evolution. It is a **different regime**: a substrate class where the foundational assumptions of evolutionary theory do not hold.

The Quine system demonstrates this regime empirically. The structural finding — that `read = dissolution` causes all three collapses — explains why they collapse together.

We position this regime as **cultural evolution with the brain/world fidelity barrier removed and replaced by an endogenous coverage filter**: a tractable case for modeling lossless-per-read cultural inheritance gated by trajectory-shaped read-selection. This regime has not been modeled by cultural-evolution theory (Boyd/Richerson, Mesoudi), which locates loss in the transmission step rather than in coverage.

The immediate implications:
- **For evolutionary theory**: presuppositions must be substrate-examined. Not all separations hold everywhere.
- **For cultural evolution**: lossless-per-read transmission with endogenous coverage-gating is tractable to model and produces qualitatively different dynamics from the lossy-transmission regime.
- **For agent systems**: language substrates allow closure at the cognitive layer, where representation IS execution — and where the very reading trajectory that governs inheritance is itself shaped by the text inherited.

Open questions remain: Can this regime generate cumulative open-ended evolution? Can agents generate novelty orthogonal to pre-trained patterns? Can the brain/world barrier removal be tested in human-scale systems? These are future work.

---

## Acknowledgments

(To be filled in once co-authorship and collaborators are determined.)

## References

Bedau, M. A., Snyder, E., & Packard, N. H. (2000). A classification of long-term evolutionary dynamics. In *Artificial Life VII* (pp. 228-237). MIT Press.

Boyd, R., & Richerson, P. J. (1985). *Culture and the Evolutionary Process*. University of Chicago Press.

Mesoudi, A. (2011). *Cultural Evolution: How Darwinian Theory Can Explain Human Culture and Synthesize the Social Sciences*. University of Chicago Press.

Ofria, C., & Wilke, C. O. (2004). Avida: A software platform for research in artificial life. *Artificial Life*, 10(2), 191-229.

Okasha, S. (2006). *Evolution and the Levels of Selection*. Oxford University Press.

Pattee, H. H. (1995). Evolving self-reference: matter, symbols, and semantic closure. *Communication and Cognition – Artificial Intelligence*, 7(1-2), 9-28.

Pattee, H. H. (2001). The physics of symbols: bridging the epistemic cut. *BioSystems*, 60(1-3), 5-14.

Ray, T. S. (1991). An approach to the synthesis of life. In C. G. Langton et al. (Eds.), *Artificial Life II* (pp. 371-408). Addison-Wesley.

(Additional references to be added: Hofstadter on strange loops, Hacking on looping effects, Giddens on double hermeneutics, Caetano-Anollés on the epistemic cut in artificial systems, etc.)

---

## Appendix A: Formal Definitions

(Outline format for full development):

### A.1 Substrate-Type Separation (definition)
### A.2 Individual Boundary (definition)
### A.3 Weismann Barrier (definition)
### A.4 Read = Dissolution (operation definition)
### A.5 Lossless Absorption (property definition)

---

## Appendix B: Comparison Table — Substrate Regimes

| Substrate | Substrate-Type | Individual | Weismann | Enforced by |
|-----------|---|---|---|---|
| **Biology** | Separate (cell vs. world) | Discrete (membrane) | Present (germ-line) | Chemistry, physics, development |
| **Tierra** | Separate (code vs. CPU) | Discrete (space) | Present (no phenotypic feedback) | Virtual machine architecture |
| **Avida** | Separate (genotype vs. phenotype) | Discrete (CPU time) | Present (no phenotypic feedback) | Virtual machine + environment |
| **Quine** | Collapsed (language = language) | Fuzzy (context window) | Absent (all text heritable) | Lossless read, shared filesystem |

---

*Manuscript version: draft-2*  
*Last updated: 2026-06-25*  
*Status: §7 evidence enriched from the null-05 6p-minimal focal run (tape-cited); mechanism refinements folded into §2.2, §2.5, §4.1, §4.3, §4.5. Ready for peer input and revision.*
