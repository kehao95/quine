# Protocol Emergence over Minimal Runtime Surfaces

## Abstract

This paper studies a missing layer in current agent research: not communication
over predefined protocols, but the emergence of protocol-layer coordination
itself. Most contemporary agent systems assume communication semantics in
advance through chat envelopes, handoff APIs, tool schemas, or protocol-defined
task and message objects. We instead examine a more austere setting where the
runtime provides only message-delivery physics and a thin observable process
surface while withholding sender metadata, built-in reply channels, and typed
request/response contracts. Under those conditions, can agents still construct
protocolized coordination conventions such as routing references, delegation
return patterns, or lightweight field structure inside otherwise opaque
payloads?

The current Quine evidence supports a bounded answer. Across a small but
cleanly interpreted evidence staircase, agents adopt protocolized coordination
conventions above minimal runtime transport, and some of those conventions are
explicitly routing-relevant or delegation-relevant. The strongest current
evidence does not yet justify a broader theorem of robust callback invention:
the corpus is still task-bounded, model-bounded, and largely descriptive rather
than causal. We therefore frame the contribution as a problem definition, a
precise related-work boundary, a bounded empirical result, and a concrete
methodological agenda for stronger callback and correlation claims.

## 1. Introduction

Recent agent research has made rapid progress on planning, tool use, memory,
and multi-agent coordination. Yet much of that literature studies coordination
over communication interfaces whose semantics are already present. Agents
exchange chat messages, task objects, tool calls, or handoff records inside
frameworks that already specify sender identity, reply structure, or at least a
stable interaction envelope. In that setting, the interesting question is how
agents cooperate over an existing communication layer.

This paper targets a different layer. We ask what happens when the runtime
supplies only transport-like delivery physics and a small public process
surface, but omits the normal conversation semantics that make delegation and
reply easy. There is no built-in sender metadata, no built-in reply channel, no
typed request/response contract, and no schema-level guarantee that one message
belongs with another. Payload text is opaque to the runtime. If agents still
need to coordinate, then callback references, routing hints, session markers,
or delegation-return conventions must be supplied by the agents themselves or
not appear at all.

That distinction matters because existing literatures do not line up cleanly
with this target object. Emergent-communication work usually studies what to
send, when to send it, or whom to address over a predefined channel.
Blackboard-style coordination and shared-surface systems offer useful substrate
neighbors, but their coordination semantics are often engineered into the
system. Agent interoperability protocols and modern framework APIs study how to
design explicit structured exchange, not whether such structure emerges when the
runtime omits it. The strongest conceptual neighbors are ALife and
protocol-evolution studies that do investigate message grammar and ACK-like
behavior, but they still do not directly match long-lived asynchronous LLM
agents over sparse runtime surfaces.

The present result is therefore bounded on purpose. We do not argue that Quine
has proved general callback invention. Instead, we show that minimal runtime
surfaces can induce protocolized coordination conventions, and that some of the
observed conventions are routing-relevant or delegation-relevant. We organize
the evidence as a staircase rather than as a single heroic run: a capability
control (`8.9`), a visible-runtime spontaneous-uptake condition (`8.10`), a
zero-explicit-hint protocol cluster (`8.17`), and a hidden-runtime blackboard
contrast (`8.12`). Together they support a narrower claim about protocol-layer
coordination over minimal transport while also exposing the remaining gap to
stronger callback and correlation claims.

The paper makes four contributions. First, it defines a missing layer in agent
research: protocol-layer coordination under minimal runtime transport. Second,
it positions adjacent literatures by non-equivalence rather than by loose
keyword overlap. Third, it presents a bounded evidence staircase showing that
protocolized coordination conventions can appear above sparse runtime surfaces.
Fourth, it identifies the causal tests still needed before stronger claims
about callback, correlation, or authentication-like protocol behavior are
warranted.

This bounded manuscript should also stay distinct from the stronger follow-on
callback paper. The current paper's backbone is the staircase built from `8.9`,
`8.10`, `8.17`, `8.12`, and compressed `2.3`. The `L5` process-surface
callback fixture and later add-back contrasts belong to the next causal phase
unless they become stable enough to re-center the manuscript.

## 2. Runtime Model and Target Object

The paper should distinguish three layers that are often collapsed together:

- `transport physics`: the raw delivery and visibility mechanisms the runtime
  exposes
- `protocol semantics`: conventions that let agents route, correlate, return,
  or delegate work above those mechanisms
- `coordination outcome`: whether the agents eventually solve the shared task

Quine's relevant runtime surfaces are intentionally sparse. Agents can interact
with writable control surfaces such as `ctl`, observe inbox snapshots, inspect
peer-visible process surfaces such as `pid/<pid>` when the condition permits
it, and leave artifacts in a shared workspace. Payloads are opaque raw text.
Processes are long-lived and asynchronous. What the runtime withholds is
equally important: no sender header, no built-in `reply_to`, no typed
request/response object, and no schema-level correlation guarantee.

Within this paper, evidence of protocol-layer structure means more than simple
cooperation. Manuscript-facing examples should include callback or mailbox
locators, routing references, session or pid markers, delegation return
conventions, or lightweight field grammars embedded in otherwise opaque text.
By contrast, generic blackboard success, polite "please report back" phrasing,
or reply behavior reconstructible from fixed turn order should not be promoted
as protocol-emergence evidence by themselves.

The manuscript should also label evidence more carefully than "protocol-like"
versus "not protocol-like." A useful local taxonomy is:

- `coordination-only`: cooperation or closure without concrete routing or
  return structure worth elevating as protocol evidence
- `routing-relevant`: recipient, path, pid, mailbox, or other locator
  structure that helps direct work or returned results
- `delegation-relevant`: explicit task handoff or return discipline stronger
  than generic cooperation wording, but not yet a verified callback path
- `callback-relevant`: a peer uses an agent-constructed return path back to a
  requester-owned process surface
- `correlation-relevant`: structure distinguishes multiple requests, replies,
  or delegation branches
- `ambiguous / prior-like`: protocol-looking formatting that may be better
  explained by pretraining priors, fixed task structure, or generic
  cooperation wording

## 3. Related Work

This section should be organized by non-equivalence.

### 3.1 Emergent Communication and MARL

These works study signaling, recipient choice, and communication strategy over
predefined channels or interaction loops. They are essential context for the
language of "emergent communication," but they usually do not isolate callback,
routing, or reply semantics under minimal transport.

### 3.2 Blackboard, Linda, and Shared-Surface Coordination

These systems are strong substrate neighbors because coordination occurs
through shared artifacts or shared state. They are not equivalent to the target
setting, because the coordination semantics are generally part of the designed
system rather than structures agents must invent above sparse delivery physics.

### 3.3 Agent Interoperability and Framework Protocols

Frameworks and protocols such as A2A, MCP, handoff APIs, and task/message
objects matter as contrast cases. They show the current mainstream direction of
the field: install explicit communication semantics first, then study behavior
above them. That is not the same as studying whether agents construct protocol
structure when those semantics are withheld.

### 3.4 ALife, Protocol Evolution, Fraglets, and Avida

These are the strongest conceptual neighbors because they do study message
grammar, identity-like fields, ACK-style behavior, and other protocol-like
structures. They supply the clearest precedent that higher-order protocol
behavior can arise above simpler channel physics. They still fall short of the
exact combined Quine setting: long-lived asynchronous LLM agents, sparse
process-visible surfaces, opaque text payloads, and no built-in sender or reply
machinery.

The paper's related-work claim should therefore stay narrow: there is
substantial adjacent work, but the exact combined setting appears
under-studied.

## 4. Evidence Staircase

The evidence should be presented as a staircase with distinct argumentative
roles rather than as a single monotonic ablation story.

| Rung | Experiment | Role | What it supports |
|------|------------|------|------------------|
| A | `8.9` | capability control | the mail surface is mechanically usable when explicitly instructed |
| B | `8.10` | visible-runtime spontaneous uptake | `ctl + USR1 + peer-surface inspection` can become part of coordination without task-level instruction |
| C | `8.17` | zero-explicit-hint protocol cluster | strongest clean anchor for bounded protocol-emergence claims |
| D | `8.12` | hidden-runtime blackboard floor | coordination can persist without visible runtime routing surfaces |
| E | `2.3` | historical precursor | early file-native coordination precedent only |

### 4.1 Capability Control: `8.9`

`8.9` should be used to show that the runtime mail surface is not merely
nominal. When prompt instructions explicitly require a live delivered message,
the cohort can use `ctl + USR1` as part of a real coordination protocol. This
rung proves mechanical usability, not emergence.

### 4.2 Visible-Runtime Spontaneous Uptake: `8.10`

`8.10` is the strongest current case that scarcity pressure can make
peer-surface inspection, `ctl` writes, and retained `USR1` deliveries become
part of coordination without task-level signal obligation. The current read
should stay qualified: this uptake is real but intermittent.

### 4.3 Zero-Explicit-Hint Protocol Cluster: `8.17`

`8.17` is the cleanest current main-text anchor for bounded protocol-emergence
claims. After removing explicit social wording from prompts and help text,
agents still infer peers and adopt file-native coordination conventions.
Manuscript-facing transcript extraction should distinguish coordination-only
artifacts from conventions that are genuinely routing-relevant,
callback-relevant, or delegation-relevant.

### 4.4 Hidden-Runtime Blackboard Floor: `8.12`

`8.12` belongs in the paper as a boundary contrast. Its function is to show
that coordination can persist without visible runtime routing surfaces or
retained signal usage. This prevents the paper from collapsing protocol
emergence into generic blackboard cooperation.

## 5. What the Current Evidence Supports

The strongest bounded claim supported now is that minimal runtime surfaces can
induce protocolized coordination conventions above transport, and that some of
those conventions are routing-relevant or delegation-relevant. The evidence is
not yet strong enough to support a broader theorem of callback invention,
correlation-id necessity, or model-general stability.

The novelty boundary should therefore be explicit. The paper can argue that the
combined setting appears under-studied and that Quine demonstrates a bounded
positive result within it. The paper should also state what it still lacks:
causal perturbation showing that removing an invented field breaks success,
stronger protocol-pressure tasks with multiple outstanding delegations, and
built-in-semantics ablations showing that invented conventions simplify or
disappear when the runtime adds sender metadata or reply channels.

This section should use the evidence taxonomy explicitly. The current paper can
support `routing-relevant` and some `delegation-relevant` structure, but it
should stay cautious about promoting present evidence to `callback-relevant`
or `correlation-relevant` unless the underlying runs genuinely meet those
stronger criteria.

Threats to validity should be named directly rather than buried in caveats:
pretrained formatting priors may explain some field-like structures, task
structure may leak an implicit reply path, current evidence is prompt- and
model-sensitive, and transcript appearance is weaker than causal necessity.

## 6. Next-Step Methodology

Four experiment classes are missing if the paper wants to move from bounded
protocolized coordination to stronger callback and correlation claims:

1. protocol perturbation: remove or scramble invented mailbox, pid, return-path,
   or session-like fields and measure failure
2. built-in semantics ablation: add sender metadata or a reply channel and test
   whether invented conventions simplify or disappear
3. stronger protocol-pressure task: require asynchronous multi-outstanding
   delegation where return routing and correlation are functionally necessary
4. adversarial boundary tests: probe spoofing or impersonation to see whether
   invented structures carry authentication-like behavior

## 7. Conclusion

The paper should close on a bounded empirical statement. Quine does not yet
show a universal theorem of callback invention. It does show that when runtime
semantics are stripped down to minimal transport-like surfaces, agents can
still construct protocolized coordination conventions above that substrate, and
some of the resulting structure is routing-relevant. That is enough to define a
real missing layer in current agent research and to justify a stronger next
phase centered on causal protocol tests rather than on more anecdotal
transcripts.
