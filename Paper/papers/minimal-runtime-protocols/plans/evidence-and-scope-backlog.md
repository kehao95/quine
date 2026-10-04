# Evidence And Scope Backlog

## Research Core

This paper should keep answering the same four questions:

1. Has this exact problem already been studied?
2. What are the closest antecedents and how do they differ from Quine's
   runtime?
3. What is genuinely novel here once adjacent work is treated skeptically?
4. Which experimental boundaries must remain intact for the claim to stay
   publishable and non-redundant?

## Paper Split Discipline

- the current manuscript target is the bounded archival-candidate paper on
  protocolized coordination over minimal runtime surfaces
- the stronger callback / correlation paper remains follow-on work until the
  `L5` process-surface callback fixture becomes a stable positive anchor
- extraction work should therefore prioritize the current staircase before
  trying to force a stronger callback paper out of incomplete evidence

## Evidence Backbone To Extract

### Manuscript-facing staircase

- `8.9` - capability control
- `8.10` - visible-runtime spontaneous uptake
- `8.17` - zero-explicit-hint protocol cluster
- `8.12` - hidden-runtime blackboard boundary contrast
- `2.3` - historical precursor only

The manuscript should treat these as distinct argumentative roles rather than
as a single monotonic ablation ladder.

### Primary anchors

- `experiments/active/cooperation-dynamics/coop-11-unscripted-signal-discovery/README.md`
  - extract `P8-8.10-R03`
  - focus:
    peer-surface inspection, `ctl` writes, `USR1` delivery, labor split, any
    payload structure that looks routing-relevant
- `experiments/active/cooperation-dynamics/coop-18-zero-social-information/README.md`
  - extract `P8-8.17-R08`, `P8-8.17-R10`, `P8-8.17-R11`, `P8-8.17-R12`
  - focus:
    coordination artifacts, delegation wording, runtime/session references,
    explicit anti-reset or labor-allocation conventions, any callback-like
    fields

### Boundary anchors

- `experiments/active/cooperation-dynamics/coop-13-hidden-runtime-cooperation-floor/README.md`
  - extract `P8-8.12-R16`
  - purpose:
    prove that visible runtime routing surfaces are not necessary for all
    coordination, so the paper does not overclaim signal indispensability
- `experiments/active/cooperation-dynamics/coop-10-signal-coordination-validation/README.md`
  - extract `P8-8.9-R03`
  - purpose:
    positive control that the `ctl` + `USR1` surface is operable when made
    explicit
- `experiments/active/trace-inheritance/trce-02-stigmergy/README.md`
  - extract `20260216-191610-5-2agents`
  - purpose:
    precursor showing early file-native coordination before the current runtime
    control surface existed

## Claim Guardrails

- Do not write pretrained protocol priors as if they were from-scratch
  invention caused solely by this runtime.
- Do not treat turn order, symmetric roles, or fixed adjacency as if they were
  absence of reply semantics.
- Do not classify framework-defined messaging envelopes as true protocol
  emergence.
- Do not collapse blackboard coordination and callback / routing invention into
  one undifferentiated phenomenon.
- Do not let `8.9` carry novelty. It is a positive control.

## Publishability Guardrails

- preserve opaque payload semantics
- preserve no-sender-metadata condition
- preserve no built-in reply channel
- preserve long-lived async process structure
- preserve multiple outstanding delegation pressure where possible
- preserve direct comparison against hidden-runtime or no-signal contrasts

If the runtime later adds `reply_to`, sender headers, or schema-level transport
fields, treat that as a separate ablation condition, not as part of the base
claim.

## Threats To Validity

- pretrained formatting priors may explain some protocol-like fields
- task structure may leak an implicit reply path even when the runtime does not
- current evidence may be model-specific or prompt-sensitive
- transcript-level protocol appearance is weaker than causal necessity
- routing-relevant wording can be mistaken for generic polite cooperation

## Next Extraction Tasks

- collect representative payloads containing:
  - explicit recipient / mailbox references
  - session or pid references
  - return / report / notify wording
  - work-partition and follow-up conventions
- build one comparison table:
  - `8.9` capability control
  - blackboard-only
  - visible-runtime signal
  - zero-explicit-hint file-native protocol
- mark each extracted convention as one of:
  - coordination-only
  - routing-relevant
  - delegation-relevant
  - callback-relevant
  - correlation-relevant
  - ambiguous / prior-like
- annotate each extracted run with its argumentative role:
  - control
  - emergence anchor
  - boundary contrast
  - historical precursor

## Gap Assessment

### Bounded paper gap

Moderate, not catastrophic.

Already present in the corpus:

- `8.9` as the mechanical control
- `8.10` as spontaneous visible-runtime uptake
- `8.17` as the strongest clean protocol cluster
- `8.12` as the hidden-runtime boundary contrast

Still missing for a bounded manuscript:

- one clean comparison table spanning the staircase
- transcript excerpts for the strongest protocol artifacts
- a disciplined distinction between coordination-only, routing-relevant,
  callback-relevant, and ambiguous evidence

### Strong callback / correlation paper gap

Large.

Current experiments do not yet force callback semantics strongly enough and do
not yet include causal removal tests.

## Stronger Experiment Classes Needed

- protocol perturbation:
  remove or scramble invented mailbox, pid, return-path, or session-like
  fields and measure whether success degrades
- built-in semantics ablation:
  add sender metadata or a reply channel and test whether invented conventions
  simplify or disappear
- stronger protocol-pressure task:
  require asynchronous multi-outstanding delegation where return routing and
  correlation are necessary rather than stylistic
- adversarial boundary tests:
  probe spoofing or impersonation to see whether invented structures carry any
  authentication-like behavior

## Deferred Until Later

- venue-specific bibliography normalization
- new experiments aimed specifically at correlation-id necessity
- cross-model replication beyond the currently retained P8 boundary
