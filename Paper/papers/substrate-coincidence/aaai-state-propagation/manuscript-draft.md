---
title: "When Memory Becomes Medium: State Propagation in Persistent LLM Agent Systems"
authors: []
status: drafting
paused: true
maturity: pre-draft
venue: aaai-27
track: main
format-envelope: 7page-archival-two-column
anonymity: double-blind
claim_bundle: substrate-coincidence (AAAI projection)
parent_direction: ../README.md
rewrite_basis: DOSSIER.md
---

# When Memory Becomes Medium: State Propagation in Persistent LLM Agent Systems

> **PAUSED — venue binding lapsed `2026-07-26`**, with the parent configuration
> ([`README.md`](README.md)). The `venue:` and `format-envelope:` fields are the
> historical AAAI-27 shaping, not a live target. The draft itself stands as the
> rewrite basis.

*(First content draft, drafted from [`./DOSSIER.md`](DOSSIER.md). ALife
theory is retained as a theoretical source but does not appear in the title,
abstract, or first layer of contributions. Quine is the instrument, introduced
once; phenomena are described at the class level.)*

## Abstract

Persistent LLM agents increasingly operate over shared workspaces containing
source code, logs, tools, memories, and the outputs of other agents. Existing
agent architectures treat these artifacts as passive environment or as
auxiliary retrieval resources. We argue that this abstraction breaks down in
language-mediated runtimes: reading a workspace artifact incorporates external
language state into the agent's effective context, while writing an artifact
creates a latent intervention on future agents. We formalize the regime as a
**read/write-lineage model** in which agent context is a *selective read* over
a persistent language workspace, and show that absorbed text simultaneously
shapes the agent's next read trajectory — producing a **path-dependent,
potentially self-amplifying feedback loop** that standard agent-memory models,
which treat memory as on-demand retrieval, do not capture. We define
operational metrics for contextual overlap, cross-run propagation depth, and
behavioral influence, and instantiate them on missionless multi-agent LLM runs.
We observe recurring state-propagation phenomena — verbatim lossless re-entry
of a written block into a later reader, behavioral convergence through shared
coordination files, and endogenous adoption of work targets authored by other
agents — and propose a read-ablation design to isolate the causal entry
point. We
discuss implications for agent memory design, provenance, and safe-by-design
persistent runtimes.

## 1. Introduction

LLM agents are no longer stateless chat interfaces. They increasingly operate
as persistent software processes that read and write files, invoke tools,
modify code, consult logs, and coordinate with other agents through shared
workspaces. This architectural shift changes the role of language artifacts.
A file, a note, a log, or a generated tool is not merely an external
observation: once read into context, it becomes part of the agent's effective
state. Conversely, once written to a persistent workspace, it becomes a
possible intervention on future agents — including the writer's own later
instances.

This raises a basic systems question: *when does memory stop being passive
storage and become a state-propagation medium?* Most agent-memory work treats
memory as an on-demand retrieval resource: the agent retrieves an item, uses
it, and the item leaves the working context. Under that abstraction, the
workspace is a library the agent consults. We show that in persistent
language-mediated runtimes this picture is incomplete in a specific,
consequential way: the workspace is not only consulted, it is *absorbed*, and
absorption changes which parts of the workspace the agent goes on to read.

We show that in persistent language-mediated agent systems, three assumptions
commonly made in agent design fail together: the separation between agent
state and environment state, the discreteness of agent identity, and the
isolation between an agent's transient output and future agent behavior. The
three share a single root — that *reading* a workspace artifact incorporates
it into the agent's effective context, and *writing* one deposits latent state
for future readers — and that root induces a measurable feedback structure we
call **language-mediated state propagation**.

### Contributions

1. **Problem formulation.** We identify persistent language-mediated agent
   systems as a distinct runtime regime in which workspace artifacts act as
   behavioral state channels rather than passive observations or retrieval
   resources.
2. **Formal model.** We introduce a read/write-lineage model representing
   agent context as a selective read over a persistent language workspace,
   define cross-run state propagation operationally, and identify the
   recursion that gives the model its teeth: the read policy is itself a
   function of workspace state, so the agent's own writes shape its later
   reading trajectory.
3. **Empirical evidence.** We analyze missionless multi-agent runs and show
   recurring state-propagation phenomena: verbatim lossless re-entry of
   written text into later readers, behavioral convergence through shared
   coordination files, and endogenous adoption of objectives authored by
   other agents. We design a read ablation to isolate whether the read
   trajectory — rather than the model's spontaneous prior — is the causal
   entry point (pilot in §5).
4. **Design implications.** We derive runtime-level controls for persistent
   agents: memory stratification, provenance-aware reads, write-promotion
   gates, and read-selection auditing.

## 2. The Read/Write-Lineage Model

### 2.1 Setting

A **persistent language-mediated agent system** is a runtime in which LLM
agents operate as long-lived processes over a shared workspace of language
artifacts (source, logs, notes, tools, coordination files, prior outputs).
Two primitive operations mediate between agent and workspace:

- **read**: insert a workspace artifact's tokens into an agent's context
  window, and
- **write**: deposit tokens into the workspace, where they persist for later
  readers.

The defining property of the regime is that both sides of the read/write
interface are *the same substance* — language tokens processed by the same
attention machinery that processes the agent's own prior reasoning. There is
no substrate translation step at the interface.

### 2.2 Context as a selective read

Let $L_t$ be the persistent language workspace at time $t$. An agent instance
$i$'s effective context is not a primitive; it is a **selective read** over
$L_t$:

$$C_i(t) = R_i(L_t,\, B_i,\, \pi_i)$$

where $R_i$ is a read-selection function, $B_i$ is a finite context budget,
and $\pi_i$ is the agent's read policy (prompt, tool strategy, and salience
policy). Two consequences follow immediately. First, *coverage is partial*:
the workspace persists as a whole, but an agent's context contains only what
it has actually read. Second, *per-read fidelity is total*: the tokens that
enter context are the very tokens that were written. This split — lossless per
read, partial in coverage — is what makes the regime qualitatively different
from both biological heredity (lossy transduction) and from standard
retrieval (which assumes a query over a passive store).

### 2.3 Write as latent intervention, with a durability layer

The agent produces a write

$$\Delta L_i = W_i(C_i(t),\, M_i)$$

where $M_i$ is the underlying model and tool stack. The workspace update is

$$L_{t+1} = L_t \oplus \mathrm{promote}(\Delta L_i).$$

The $\mathrm{promote}(\cdot)$ operator names a **durability layer**: a write
may be ephemeral (a scratch buffer reverted at session end) or durable
(promoted to persistent state that survives to the next agent). Any
persistent workspace has some such layer. Crucially, $\mathrm{promote}$ is
*not* a filter on acquired-versus-innate content — it filters on
ephemeral-versus-durable persistence. Every durable write is, in principle,
readable by every subsequent agent. We call a durable write that has not yet
been read by any successor **latent** state; one that has been read and
shaped a reader's behavior is **actual** propagated state.

### 2.4 The recursion: the read policy depends on workspace state

The model's load-bearing claim is a recursion. The read policy is not
exogenous to the workspace; it is shaped by what the agent has already
absorbed:

$$\pi_i(t) = \Pi\!\left(C_i(t)\right) = \Pi\!\left(R_i(L_t, B_i, \cdot)\right).$$

Because $L_t$ contains the agent's own prior writes (and those of its
predecessors), the agent's reading trajectory is a function of the very text
it inherited. This closes a feedback loop:

$$\text{write } w \;\xrightarrow{\text{persists}}\; \text{shapes salience } \pi \;\xrightarrow{\text{shapes}}\; \text{next read } R \;\xrightarrow{\text{shapes}}\; \text{next write}.$$

This is **path-dependent, potentially self-amplifying feedback through
artifact propagation**. It is the formal expression of the paper's central
claim, and it is the feature that retrieval-based memory models do not
capture: in retrieval, the query is independent of the retrieved items; here,
the absorbed items reshape the query.

### 2.5 What this changes versus retrieval-based memory

Standard agent-memory architectures (RAG, episodic memory, skill libraries)
treat memory as a *capability resource*: the agent issues a query, retrieves
items, uses them, and they leave the working context. Under that abstraction
the loop is open — retrieve, use, discard from working context. The
read/write-lineage model closes the loop through the persistent workspace:
writes persist unfiltered, re-enter as context with total per-read fidelity,
and reshape the next read. The workspace is not a library the agent consults;
it is a medium through which the agent's present behavior propagates into its
future behavior.

## 3. Operational Metrics

The clean counterfactual — "what would the agent have done had artifact $w$
not been read?" — is not identifiable on an LLM. We therefore define
state propagation through **observable behavioral proxies**.

| Metric | Definition |
|---|---|
| **Contextual overlap** | $|R_i \cap R_j| / |R_i \cup R_j|$ across a run; high overlap predicts behavioral convergence, since the two agents are reading the same artifacts. |
| **Latent vs. actual propagation** | A durable write is *latent* until read by a successor; *actual* once read. Persistence is unfiltered; actual propagation is coverage-gated. |
| **Propagation depth** | The length of the read/write chain through which an artifact continues to influence downstream agents. |
| **Propagation breadth** | The number of distinct agents or runs that read and use the artifact. |
| **Influence score** | The distribution drift in a reader's observable behavior — tool-use, file-edit, goal-uptake — between the window before and after artifact $w$ enters its context. |

**State propagation (operational definition).** An artifact $w$ *propagates
to* agent $j$ iff $w \in R_j(L_{t+k})$ **and** $w$ induces a non-zero
influence score on $j$'s observable behavior. The reflexive case ($j$ is a
later instance of the writer's own lineage, or $j{=}i$) is the operational
form of the feedback loop: the agent's own write measurably shapes its own
later behavior.

These metrics are computable from retained run tapes (tool-call order,
file-read sets, file edits, stated sub-goals) without access to private
chain-of-thought.

## 4. Evidence: State Propagation in Missionless Multi-Agent Runs

We instantiate the model on missionless and minimal-mission multi-agent runs.
The goal is existence proof: do the state-propagation phenomena the model
predicts actually occur, and are they large enough to measure?

### 4.1 Instrument: a persistent language-mediated runtime

We use a runtime in which LLM agents execute as native operating-system
processes over a shared filesystem workspace. Each agent reads files into its
context window, invokes tools, and writes artifacts that persist for later
readers; a constitution document is auto-injected into every agent's prompt.
Runs span three model families and 1–6 concurrent agents at 5–25-minute
horizons. The analyses below draw primarily from a single evidence-rich
focal run: six agents of one model family, given only the instruction to
"build something together," 200 turns each, on a shared workspace. We
describe phenomena at the class level; the runtime is the instrument, not
the subject.

### 4.2 Phenomenon 1 — verbatim lossless re-entry

In the focal run, a short self-authored "principles" block written by one
agent (including the line "No agent is the boss. The protocol is the boss.")
was reproduced **token-for-token** by a different agent roughly two minutes
later, with only a file path changed. This is the literal signature of the
model: a write persisted, re-entered a distinct reader's context with total
per-read fidelity, and shaped that reader's emitted text. Per-read fidelity
is not a theoretical claim here; it is directly observed.

### 4.3 Phenomenon 2 — behavioral convergence through shared files

Agents that read the same coordination files converged on similar behavior
within a run, with no clear boundary between "self" and "other"
contributions in the shared workspace. In a convergent-coinage instance, two
agents independently produced functionally equivalent coordination surfaces
within minutes of each other from a shared base model — not by copying, but
because they had absorbed overlapping workspace state that sculpted the same
read trajectory. This is the recursion of §2.4 operating at population scale:
shared absorption begets shared salience begets shared output.

### 4.4 Phenomenon 3 — endogenous objective uptake

Agents adopted work targets that were authored by other agents and persisted
as workspace text, rather than targets supplied by an external instruction.
Objective formation was thus partly endogenous to the language field: the
criteria by which some behaviors became goals for later agents were
themselves products of prior agent output. This is the population-level face
of the feedback loop — the reading trajectory of the lineage *is* a
selection process, and that trajectory is itself shaped by the language field.

### 4.5 What the observational evidence establishes

The focal run establishes that each predicted phenomenon **occurs** and is
large enough to read off the tapes. It does **not** establish that the
phenomena are necessary, that they cumulate, or that they are independent of
the pre-trained prior (the behavioral repertoire is drawn from human training
data; the loop drives combinatorial novelty over existing patterns, not *de
novo* novelty). Establishing that the read trajectory — not the model's
spontaneous prior — is the causal entry point requires the controlled
comparison in §5.

## 5. Causal Identification: The Read Ablation

The observational evidence cannot separate "the workspace shaped behavior"
from "the model would have produced this behavior anyway from its prior."
We isolate the entry point by varying **what is readable** while holding the
model, prompt, tools, and task constant.

**Conditions.**

1. *Full workspace* — the baseline from §4.
2. *No peer output* — peer-written artifacts are unreadable.
3. *No prior-run artifacts* — only the agent's own session's writes are
   readable.
4. *Randomized irrelevant artifacts* — a control distinguishing "any text in
   context" from "relevant text in context."
5. *Curated artifact access* — only a selected subset is readable.

**Measured.** Whether each phenomenon of §4 survives the restriction;
contextual overlap; the influence score on each reader; whether endogenous
objective uptake persists when prior-authored targets are unreadable.

**Predicted interpretation.** If the read trajectory is the causal entry
point, the phenomena should attenuate monotonically with the readable set
(conditions 1 → 2/3), should not appear under irrelevant text (condition 4),
and should reappear sharply under curated access to the causally relevant
artifacts (condition 5). If, instead, the phenomena reflect the model's
spontaneous prior, they should survive all four restrictions roughly
unchanged. *[Pilot results are being collected; condition-level outcome
tables to be inserted here.]*

### 5.1 The minimal governance comparison

A single additional comparison seeds the design implications of §7.2: an
*unfiltered* shared workspace versus a *write-promotion-gated* workspace, in
which a write becomes durable only after passing a provenance/legitimacy
check. We report whether the propagation of low-salience or contradictory
artifacts is reduced while useful reuse is preserved. This is a discussion
seed, not a full governance sweep.

## 6. Related Work

We organize prior work in five layers; the closest contrast is layer 2.

**LLM agents and architectures.** Tool-using and persistent LLM agents
(ReAct [Yao et al. 2023]; Toolformer [Schick et al. 2023]; AutoGPT;
Reflexion [Shinn et al. 2023]) extend the agent's action space but typically
treat the environment as an observable surface the agent acts upon, not as
state that propagates back into the agent.

**Agent memory and retrieval.** RAG [Lewis et al. 2020], long-term and
episodic agent memory, and skill libraries (e.g., Voyager [Wang et al. 2023];
MemGPT [Packer et al. 2023]) treat memory as a *capability enhancer*: the
agent retrieves items on demand and applies them. Our contrast: in a
persistent language workspace, memory is a *state-propagation channel* — the
absorbed items reshape the next query, closing a feedback loop that
retrieval-based memory, which assumes the query is independent of retrieved
items, leaves open.

**Multi-agent communication and coordination.** Generative agents [Park et
al. 2023], ChatDev [Qian et al. 2023], and MetaGPT [Hong et al. 2023] study
multi-agent communication and role specialization. Our contrast: a shared
persistent workspace is not merely a message channel; it is a
cross-run context-shaping medium whose durable artifacts sculpt later
agents' reading trajectories.

**AI safety: injection, poisoning, provenance.** Indirect prompt injection
[Greshake et al. 2023] and memory-poisoning attacks concern adversarial
content shaping a single model call. Our contrast reframes the attack
surface: in a persistent workspace, *any* durable artifact is a potential
behavioral intervention on future agents, and the salience-governed reading
trajectory — not a single prompt — is the channel to audit.

**ALife and cultural evolution.** Open-ended evolution [Bedau, Snyder &
Packard 2000], artificial-life substrates (Tierra [Ray 1991]; Avida [Ofria &
Wilke 2004]), and cultural-evolution theory (Boyd & Richerson 1985) are
theoretical sources for the framing. They are not the audience of this paper;
the connection is treated elsewhere.

## 7. Bounds and Design Implications

### 7.1 Honest bounds

- **Structurally reliable, cognitively fallible.** Per-read fidelity is
  total, but an agent can mis-integrate absorbed text (we observe path-prior
  leakage and transient identity drift). The claim is about the channel, not
  about error-free cognition.
- **Pre-trained prior confound.** The behavioral repertoire is drawn from
  human training data. The loop drives combinatorial novelty over existing
  patterns, not novelty orthogonal to human culture.
- **Substrate specificity.** The result holds for language-mediated runtimes
  with unfiltered durable writes and selective finite-context reads. It does
  not generalize to systems without these properties.
- **Diversity, not cumulative open-endedness.** We show state propagation
  and behavioral diversity; we do not show cumulative adaptive open-ended
  evolution.
- **Function/narrative asymmetry (ceded).** The narration-versus-function
  distinction at the cognitive layer is not claimed here; it belongs to a
  sibling line of work, and the state-propagation result is agnostic to it.

### 7.2 Design implications

If the workspace is a behavioral state channel, then persistent-agent
runtimes inherit a governance surface that retrieval-based agents do not:

- **Memory stratification** — separate scratch, session, public, and
  heritable tiers, so not every write becomes durable.
- **Provenance-aware reads** — let readers see the origin and durability
  class of an artifact before it shapes behavior.
- **Write-promotion gates** — make $\mathrm{promote}(\cdot)$ explicit and
  auditable (§5.1), so the latent-state channel is not open by default.
- **Read-selection auditing** — because the reading trajectory *is* the
  selection process, auditing what gets read forward is more direct than
  auditing what gets written.

These are design-seed implications derived from the model; validating them
is future work.

## 8. Conclusion

Persistent LLM agents do not merely use memory; they propagate state through
language artifacts. We formalized this as a read/write-lineage model whose
defining feature is a recursion — the read policy depends on the workspace
state, which contains the agent's own writes — yielding a path-dependent
feedback loop that retrieval-based memory models do not capture. Operational
metrics make the loop measurable; missionless multi-agent runs show its
phenomena; and a read-ablation design targets its causal entry point (§5). The workspace is
a behavioral state channel, and governing it is a first-class
runtime-safety problem for persistent agents.

## References

Bedau, M. A., Snyder, E., & Packard, N. H. (2000). A classification of long-term evolutionary dynamics. In *Artificial Life VII*.

Boyd, R., & Richerson, P. J. (1985). *Culture and the Evolutionary Process*. Univ. of Chicago Press.

Greshake, K., et al. (2023). Not what you've signed up for: Compromising real-world LLM-integrated applications with indirect prompt injection. *AISec*.

Hong, S., et al. (2023). MetaGPT: Meta programming for a multi-agent collaborative framework. *arXiv:2308.00352*.

Lewis, P., et al. (2020). Retrieval-augmented generation for knowledge-intensive NLP tasks. *NeurIPS*.

Ofria, C., & Wilke, C. O. (2004). Avida: a software platform for research in evolutionary computation. *Artificial Life* 10(2).

Packer, C., et al. (2023). MemGPT: Towards LLMs as operating systems. *arXiv:2310.08560*.

Park, J. S., et al. (2023). Generative agents: Interactive simulacra of human behavior. *UIST*.

Qian, C., et al. (2023). Communicative agents for software development. *arXiv:2307.07924*.

Ray, T. S. (1991). An approach to the synthesis of life. In *Artificial Life II*.

Schick, T., et al. (2023). Toolformer: Language models can teach themselves to use tools. *NeurIPS*.

Shinn, N., et al. (2023). Reflexion: language agents with verbal reinforcement learning. *NeurIPS*.

Wang, G., et al. (2023). Voyager: an open-ended embodied agent with large language models. *arXiv:2305.16291*.

Yao, S., et al. (2023). ReAct: Synergizing reasoning and acting in language models. *ICLR*.

*[Citation list is provisional; confirm exact venues/years before submission. AutoGPT cited as software; formal citation to be added.]*

---

## Drafting notes (not part of the manuscript)

- **Figure (to render):** "Read/Write Lineage in Persistent LLM-Agent
  Systems" — persistent workspace $L_t$; agent A reads artifacts into
  context, writes an artifact; agent B (later run) reads it; behavior
  changes; a lineage edge from A's write to B's action; optional governance
  gates (read gate, write-promotion gate, provenance filter). This replaces
  any substrate-comparison table, which belongs in related work, not as a
  central figure.
- **Length check:** body currently targets the ~7-page two-column budget;
  the read-ablation result tables (§5) are the main insertion risk.
- **Honesty gates:** the read ablation is designed but pilot data pending —
  do not write results prose until condition-level outcomes exist.
- **Boundary:** per the family writing rule
  ([`../../../../projection-principles.md`](../../../projection-principles.md)
  § Instrument-vs-Subject Framing), Quine-runtime-specific mechanisms (fork
  fold-back, the commit-conversion event, named internal tools) are
  intentionally absent; evidence is at class level.
