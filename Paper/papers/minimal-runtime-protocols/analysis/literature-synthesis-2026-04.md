# Literature Synthesis: Minimal Runtime Protocols

> Paper-local synthesis for the question:
> when agents are given only transport-like runtime physics, do they invent
> callback, routing, delegation, or correlation structure for themselves?
>
> Exact bibliographic normalization is deferred. This file keeps the current
> paper-local read precise enough for scope control and later bibliography
> conversion.

## Framing

The relevant comparison class is narrower than generic "emergent
communication." The closest target object is:

- minimal transport primitives
- opaque payloads
- no built-in sender metadata
- no built-in reply channel
- no predefined request/response schema
- long-lived asynchronous processes

The most important discriminator is not whether symbols emerge, but whether
agents must construct **protocol-layer structure** such as addresses, callback
locators, correlation markers, or delegation-return conventions.

## Directly Relevant Prior Work

| Citation | Short summary | Why relevant | True protocol emergence? | Callback / routing / delegation? | Closeness |
|----------|---------------|--------------|--------------------------|----------------------------------|-----------|
| Knoester et al. 2007, *Directed Evolution of Communication and Cooperation in Digital Organisms* | Avida organisms evolve message use and message grammar in a neighbor-inbox substrate. Identity-like payload fields matter for cooperation. | Strong precedent that message conventions can evolve inside the payload rather than being fully supplied by the substrate. | Partial but real. Channel physics are predefined, but message grammar is evolved. | Partial yes. Identity-like fields and message recognition are central. | High conceptually; lower on LLM / OS-runtime specifics. |
| Yamamoto and Tschudin 2005, *Experiments on the Automatic Evolution of Protocols using Genetic Programming* | Evolves protocol implementations under changing network conditions, including ACK-style return traffic and retransmission logic. | One of the clearest precedents for literal protocol emergence rather than only semantic signaling. | Yes within a constrained execution model. | Yes. ACK / return-path behavior is part of the object. | High on protocol logic, lower on agent cognition. |
| Tschudin 2003, *Fraglets* | Introduces a tiny symbol-string execution model for communication protocols, including send targets and confirmed-delivery examples. | Very close to the "minimal physics, higher-order protocol above it" framing. | Mixed. Primitives are fixed, higher-order protocol behavior can be synthesized. | Yes. Address-like targeting and ACK behavior are explicit. | High conceptually. |
| Jim and Giles 2000, *Talking Helps* | Evolved predator coordination through a shared message-board substrate. | Relevant as an early shared-surface coordination system where conventions are not fully specified in advance. | Partial. Shared board semantics are given, symbol use emerges. | Weak. Shared board reduces pressure for explicit reply routing. | Moderate. |
| Saunders and Pollack 1996, *The Evolution of Communication Schemes Over Continuous Channels* | Studies evolved signaling over continuous or aggregated channels. | Useful precedent for sender ambiguity created by channel physics rather than by missing APIs. | Partial. Channel is predefined; signaling scheme emerges. | No strong callback focus. | Moderate. |
| Ackley and Littman 1994, *Altruism in the Evolution of Communication* | Agents communicate over aggregated channels where direct sender labeling is blurred. | Another precedent for anonymity-by-physics rather than identity-by-envelope. | Partial. Channel and world are given; signaling evolves. | No explicit callback or delegation structure. | Moderate. |

## Closely Related Neighboring Work

| Citation | Short summary | Why relevant | True protocol emergence? | Callback / routing / delegation? | Closeness |
|----------|---------------|--------------|--------------------------|----------------------------------|-----------|
| Foerster et al. 2016, *Learning to Communicate with Deep Multi-Agent Reinforcement Learning* | Modern MARL emergent-communication benchmark family with explicit communication actions. | Important baseline for "emergent communication" language, but usually assumes a known channel and interaction loop. | Mostly no for protocol structure; yes for message semantics. | Usually no. Reply path is implicit in the environment. | Low to moderate. |
| Sukhbaatar et al. 2016, *Learning Multiagent Communication with Backpropagation* | CommNet uses aggregated continuous communication channels. | Useful for anonymity-by-aggregation, but still a fixed communication architecture. | No for transport-layer emergence. | No. | Low to moderate. |
| Das et al. 2019, *TarMAC* | Learns what to send and whom to address in a predefined multi-agent graph. | Closest MARL precedent for emergent targeting. | Partial. Targeting emerges, but the architecture already defines distinct recipients. | Addressing yes, callback conventions no. | Moderate. |
| Lazaridou and Baroni 2020, *Emergent Multi-Agent Communication in the Deep Learning Era* | Survey of the modern emergent-communication literature. | Useful primarily as negative evidence: most benchmark families assume more interaction structure than this paper's runtime. | Mostly no for protocol-layer emergence. | Rare. | Contextual only. |
| Lowe et al. 2019, *On the Pitfalls of Measuring Emergent Communication* | Shows that naive metrics overstate communication. | Methodologically important: descriptive-looking protocols need causal tests. | Evaluation paper, not an emergence substrate. | Not directly. | High methodological relevance. |
| Carriero and Gelernter 1989, *Linda in Context* | Tuple-space / generative communication via shared substrate. | Canonical shared-surface communication precedent. | No. Coordination semantics are designer-specified. | Patterns can be programmed, but not agent-emergent. | Moderate substrate similarity. |
| Corkill 1991, *Blackboard Systems* | Shared blackboard with independent knowledge sources. | Key neighbor for shared-surface coordination language. | No. Blackboard representations and triggers are predefined. | Not in the target sense. | Moderate substrate similarity. |
| Pike et al. 1993 / 1995, *Use of Name Spaces in Plan 9* / *Plan 9 from Bell Labs* | Process surfaces and namespaces as file-like OS resources. | Strong substrate analogue for `ctl`, `pid/`, and process-surface visibility. | No. OS-interface design, not protocol-emergence study. | Enables routing, but does not study its emergence. | High substrate similarity only. |
| Wu et al. 2023, *AutoGen* | LLM multi-agent framework with predefined conversation interfaces. | Important contrast case for what this paper is not. | No. Messaging envelope is part of the framework. | Delegation yes, invented routing no. | Low. |
| Li et al. 2023, *CAMEL* | Role-playing chat agents with structured prompt scaffolding. | Another contrast case for framework-provided social structure. | No. Conversation semantics are scaffolded. | No in the target sense. | Low. |
| PheroPath 2024/2025 | Filesystem-based stigmergy protocol for LLM agents. | Relevant near-threat for filesystem-mediated coordination novelty. | No for this exact question. The coordination protocol is engineered, not an emergence study. | Some coordination structure yes, but framework-defined. | Moderate. |

## Important Negative Result

The current literature map suggests a real gap, but it is a narrow one.

What appears well covered:

- emergent signaling or symbol use under predefined channels
- shared-surface coordination through blackboards or tuple spaces
- OS-level process surfaces as interface design
- evolved protocol behavior in digital evolution or network-protocol settings
- LLM multi-agent systems with explicit messaging frameworks

What appears weakly covered or largely absent as a first-class target:

- long-lived asynchronous agent processes
- only minimal process surfaces such as `ctl`, inbox snapshots, `pid/`, and
  session roots
- opaque payloads
- no sender metadata
- no built-in reply channel
- no predefined request/response or tool-schema contract
- explicit interest in whether agents invent callback, routing, or correlation
  conventions above that substrate

That negative result is the core paper opportunity, but it should be stated
carefully:

- this does **not** mean there is no adjacent work
- it does mean the exact combined setting seems under-studied
- the strongest existing conceptual neighbors are protocol evolution and ALife,
  not current LLM agent frameworks

## Design Implications For This Paper

- Preserve the exact runtime asymmetry. If sender metadata or a reply channel is
  added by the runtime, the paper collapses toward ordinary messaging systems.
- Keep tasks asynchronous and allow multiple outstanding delegations. Otherwise
  reply structure can be reconstructed from turn order.
- Do not treat any protocol-looking field as evidence by itself. Use causal
  ablations or perturbations where possible.
- Separate blackboard coordination from callback / routing structure. They are
  adjacent, not identical.
- Keep routing claims partial until transcript extraction shows actual
  return-address or correlation behavior rather than generic "please continue"
  wording.
- Treat pretrained protocol priors as a threat to validity, not as something to
  hand-wave away.
