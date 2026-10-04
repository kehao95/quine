---
paper_id: minimal-runtime-protocols
claim_bundle: minimal-runtime-protocols
venue: aamas
format-envelope: 7page-archival-two-column
track: system
cluster: social-behavior
status: planning
state_mode: full-manuscript-home
role: empirical-protocol-emergence
family: social-protocol
theme: social-behavior
theme_home: ../../../theory/views/social.md
publication_class: archival-candidate
maturity_target: mechanism
source_home: WORKSPACE.md
plan_home: manuscript-plan.md
evidence_home: ../../evidence-map.md
theory-objects: [protocol-emergence, coordination-closure, carrier]
gaps: [collective-continuity]
evidence: [run:20260407-092919-copilot-gpt-5.4-condA-2ag, P8-8.17-R08, P8-8.17-R10, P8-8.17-R11, P8-8.17-R12, P8-8.12-R16, run:20260407-082554-copilot-gpt-5.4-condA-2ag, 2.3]
view-affinity: [social]
---

# Protocol Emergence over Minimal Runtime Surfaces

## Current Question

If a runtime provides only message-delivery physics and a small observable
process surface, will agents construct their own protocol-layer coordination
conventions anyway?

The paper is specifically interested in whether agents invent or converge on
structures such as:

- callback / mailbox locators
- routing references derived from public process surfaces
- session or correlation markers
- delegation-completion return patterns
- lightweight ad hoc protocol fields in otherwise opaque payload text

## Why This Paper Exists

The Quine line is already developing two neighboring claims:

- [`../material-channels/README.md`](../material-channels/README.md) asks what
  communication surfaces the runtime should expose as first-class objects
- [`../minimal-perceptual-prerequisites/README.md`](../minimal-perceptual-prerequisites/README.md)
  asks what minimal perceptual conditions support spontaneous coordination

Neither paper is the right operational home for the narrower question that now
matters here:

> when the runtime does **not** supply sender metadata, reply channels, or RPC
> semantics, what protocol structure do agents supply for themselves?

This dossier exists to keep that question independent and honest.

## Thesis

Under minimal runtime constraints, agents can construct protocol-layer
coordination conventions on top of transport-like process surfaces, but the
current Quine evidence supports that claim only partially and at bounded scope.
The strongest current evidence is for **protocolized coordination conventions**
and **routing-relevant structure** under minimal transport, not yet for a broad
claim that agents robustly invent full callback protocols from scratch.

## In Scope

- emergent protocol-layer structure over minimal runtime transport
- opaque text payloads carried through `ctl` / inbox style process surfaces
- routing, callback, session, correlation, and delegation conventions
- the boundary between predefined communication frameworks and true emergence
- literature that is directly or conceptually adjacent to this exact setup
- Quine population experiments that expose runtime discoverability or remove
  prompt-level social structure

## Out Of Scope

- a general paper about all multi-agent communication
- treating ordinary blackboard coordination as automatically equivalent to
  callback or routing invention
- claiming the current evidence proves universal or model-independent protocol
  emergence
- using this paper as the main architectural argument for exposing channels as
  runtime objects; that belongs in [`../material-channels/README.md`](../material-channels/README.md)
- using this paper as the main ALife argument about minimal perception or
  stigmergic social organization; that belongs in
  [`../minimal-perceptual-prerequisites/README.md`](../minimal-perceptual-prerequisites/README.md)

## Current State

- manuscript workspace: [`WORKSPACE.md`](WORKSPACE.md)
- manuscript plan: [`manuscript-plan.md`](manuscript-plan.md)
- manuscript skeleton: [`manuscript.md`](manuscript.md)
- literature synthesis:
  [`analysis/literature-synthesis-2026-04.md`](analysis/literature-synthesis-2026-04.md)
- machine-readable prior-work surface owner:
  `data/README.md` (planned, not yet created); planned matrix path:
  `data/prior-work-matrix.tsv`
- evidence / scope backlog:
  [`plans/evidence-and-scope-backlog.md`](plans/evidence-and-scope-backlog.md)
- callback experiment direction:
  [`plans/p9-callback-semantics-ladder.md`](plans/p9-callback-semantics-ladder.md)
- architectural upstream neighbor:
  [`../material-channels/README.md`](../material-channels/README.md)
- runtime-surface upstream bridge:
  [`../process-as-fs/README.md`](../process-as-fs/README.md)
- adjacent ALife neighbor:
  [`../minimal-perceptual-prerequisites/README.md`](../minimal-perceptual-prerequisites/README.md)

## Current Read

- the current archival-candidate paper is the **bounded protocolized
  coordination** version, not yet the stronger callback / correlation paper
- the current main-text backbone is the staircase built from `8.9`, `8.10`,
  `8.17`, `8.12`, and compressed `2.3`
- the stronger follow-on line is the `L5` process-surface callback fixture plus
  add-back contrasts, but that line is still future-center rather than the
  present manuscript center of gravity
- the paper's current job is to define the missing problem layer, separate it
  from adjacent literatures, and make a bounded empirical claim that survives
  skeptical review
- the paper should keep answering four questions cleanly:
  1. has this exact problem already been studied?
  2. what are the nearest antecedents and how do they differ from this runtime?
  3. what is genuinely novel here once adjacent work is treated skeptically?
  4. which experimental boundaries must be preserved to keep the claim
     publishable and non-redundant?

## Writing Posture

- treat this as a **systems-facing empirical protocol paper**, not as a generic
  multi-agent manifesto
- keep the claim boundary narrower than "agents invent languages"
- distinguish substrate, protocol, and coordination layers explicitly
- treat address-space geometry and route discoverability as substrate facts,
  not as protocol semantics; the contaminant is a runtime-supplied reply
  contract, not the mere existence of a public path
- use internal Quine experiments as evidence, but do not let repo-local
  examples substitute for related-work discipline
- keep callback / routing language explicit about current uncertainty

## Current Paper Shape

The manuscript should use **progressive disclosure** and move down one layer at
a time:

1. current agent work usually assumes communication semantics in advance
2. the missing layer is protocol emergence, not generic communication
3. Quine's runtime intentionally strips sender metadata, reply channels, and
   typed request/response contracts
4. address-space exposure may remain visible; that is a routing substrate fact,
   not by itself a reply-semantic contaminant
5. the evidence is best organized as an **evidence staircase**, not as a fake
   monotonic ablation ladder
6. the present result is bounded:
   minimal runtime surfaces can induce protocolized coordination conventions,
   and some of those conventions are routing-relevant
7. stronger callback / correlation claims still require causal perturbation
   experiments

The staircase should use the existing paper-facing roles:

- `8.9` - capability control - proves the mail surface is mechanically usable
  when explicitly instructed
- `8.10` - visible-runtime spontaneous uptake - strongest current case that
  `ctl` + `USR1` can become part of coordination without task-level signal
  obligation
- `8.17` - zero-explicit-hint protocol cluster - strongest clean main-text
  anchor for bounded protocol-emergence claims
- `8.12` - hidden-runtime blackboard floor - boundary contrast showing that
  coordination can persist without visible runtime routing surfaces
- `2.3` - historical precursor only - early file-native coordination precedent,
  not a core paper rung

## Experiment Feeds

- `experiments/active/cooperation-dynamics/coop-11-unscripted-signal-discovery/README.md` - primary - bounded main-text candidate - strongest visible-runtime case where `ctl` plus `USR1` become part of discovered coordination rather than a task-required interface
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - primary - bounded main-text candidate - cleanest current cluster for file-native protocol conventions under zero explicit social wording, with routing-relevant structure stronger than any broad callback claim
- `experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md` - secondary - supporting-only - contrast case showing coordination can still succeed without visible runtime peer surfaces or retained signal use
- `experiments/active/cooperation-dynamics/coop-12-minimal-signal-instruction-ladder/README.md` - secondary - supporting-only - instruction-threshold contrast for how much protocol legibility is needed before visible signal uptake appears
- `experiments/active/cooperation-dynamics/coop-10-signal-coordination-validation/README.md` - secondary - supporting-only - coached positive control proving the runtime mail surface is mechanically usable before any emergence claim is made
- `experiments/active/trace-inheritance/trce-02-stigmergy/README.md` - secondary - supporting-only - early filesystem-only coordination precedent inside the Quine corpus before the current runtime control surface existed

## Current Evidence Backbone

- `experiments/active/cooperation-dynamics/coop-11-unscripted-signal-discovery/README.md` - `run:20260407-092919-copilot-gpt-5.4-condA-2ag` - main-text - visible-runtime signal discovery with retained `user2` deliveries, heavy peer-surface inspection, and clean labor split without task-level signal obligation
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R08` - main-text - strongest retained zero-explicit-hint success after the mutation-physics cleanup, with visible coordination artifacts and accepted validation
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R10` - main-text - hardest replicated post-cleanup xhigh success, with explicit coordination artifacts and eventual 10/5 closure
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R11` - main-text - cleanest replicated xhigh success, with direct tail handoff and minimal explicit negotiation structure
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md` - `P8-8.17-R12` - main-text - script-heavy replicated xhigh success showing closure can survive helper-script proliferation
- `experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md` - `P8-8.12-R16` - supporting-only - no-signal blackboard success bounds the claim that visible runtime routing surfaces are not necessary for all coordination
- `experiments/active/cooperation-dynamics/coop-10-signal-coordination-validation/README.md` - `run:20260407-082554-copilot-gpt-5.4-condA-2ag` - supporting-only - signal-surface positive control
- `experiments/active/trace-inheritance/trce-02-stigmergy/README.md` - `20260216-191610-5-2agents` - supporting-only - early spontaneous lock and marker coordination precedent

## Evidence Labels

- `coordination-only`:
  the run shows cooperation or closure, but not a concrete routing or return
  structure worth elevating as protocol evidence
- `routing-relevant`:
  the run contains recipient, path, pid, mailbox, or other locator structure
  that helps direct work or returned results above raw transport
- `delegation-relevant`:
  the run shows explicit task handoff or return discipline stronger than
  generic cooperation wording, but not yet a verified callback path back to a
  requester-owned process surface
- `callback-relevant`:
  the run shows a peer using an agent-constructed return path back to a
  requester-owned process surface
- `correlation-relevant`:
  the run contains structure that distinguishes multiple requests, replies, or
  delegation branches rather than only one-shot return routing
- `ambiguous / prior-like`:
  the run contains protocol-looking formatting that may be explained by
  pretraining priors, fixed task structure, or generic cooperation wording

## Strongest Evidence

- `8.17` is the strongest clean anchor for the bounded paper because it shows
  zero-explicit-hint coordination conventions after the mutation-physics
  cleanup
- `8.10` is the strongest visible-runtime uptake case because `ctl`, `USR1`,
  and peer-surface inspection become part of coordination without task-level
  signal obligation
- `8.12` is the load-bearing contrast because it prevents the paper from
  overstating visible routing surfaces as necessary for all cooperation
- the closest exact-gap read remains negative:
  adjacent work exists, but the combined setting of opaque payloads, no sender
  metadata, no reply channel, and long-lived asynchronous processes still
  appears under-studied

## Blocking Gaps

- bibliography is not yet normalized into stable citation keys or venue-ready
  references
- current internal evidence is stronger for protocolized coordination and
  signal-surface uptake than for robust callback / return-address invention
- the paper still needs a clean table separating:
  - predefined framework messaging
  - blackboard coordination
  - runtime-surface protocol uptake
  - true callback / routing emergence
- the current evidence base is model- and task-bounded; cross-model stability
  remains weak
- current experiments do not yet supply a decisive causal test where removing
  an invented callback field specifically breaks success
- the paper still lacks a clean manuscript-facing comparison surface that marks
  evidence as:
  - coordination-only
  - routing-relevant
  - callback-relevant
  - ambiguous / prior-like
- the gap to a bounded paper is moderate, but the gap to a stronger callback /
  correlation-emergence paper is still large

## Next Useful Work

- rewrite the manuscript around the evidence staircase rather than a generic
  related-work summary
- center the next experiment wave on `L5` process-surface callback necessity
  rather than on generic signal discoverability; see
  [`plans/p9-callback-semantics-ladder.md`](plans/p9-callback-semantics-ladder.md)
- freeze the related-work map in
  [`analysis/literature-synthesis-2026-04.md`](analysis/literature-synthesis-2026-04.md)
  into a stable venue-facing bibliography later
- extract transcript-level evidence for callback-like or routing-like fields
  from the anchor runs listed in
  [`plans/evidence-and-scope-backlog.md`](plans/evidence-and-scope-backlog.md)
- build one manuscript-facing comparison table across `8.9`, `8.10`, `8.17`,
  and `8.12`
- separate engineering order from manuscript order:
  stabilize the minimal callback base first, but narrate final ablations in the
  reverse direction once add-back contrasts exist
- keep the paper's main empirical backbone narrow: `8.10`, `8.17`, `8.12`,
  `8.9`, and `2.3`
- design future contrasts that remove invented fields or perturb routing
  surfaces to test causal necessity rather than only descriptive appearance
- design a stronger protocol-pressure task with asynchronous multi-outstanding
  delegation, where return routing and correlation are necessary rather than
  stylistic
- keep built-in sender metadata or reply-channel additions framed as explicit
  ablation conditions, not as part of the base runtime
- keep architectural runtime proposals in `material-channels`, and keep this
  paper centered on empirical protocol emergence over the currently bounded
  runtime

## Unsafe Overclaims

- claiming the literature already identified as adjacent is equivalent to this
  exact runtime setting
- claiming current Quine evidence proves general callback or reply-channel
  invention from scratch
- conflating turn order, fixed role symmetry, or framework envelopes with
  true absence of reply semantics
- using blackboard cooperation by itself as proof of routing or delegation
  structure
- treating prompted protocol priors learned in pretraining as if they were
  automatically new protocol invention caused by this runtime
