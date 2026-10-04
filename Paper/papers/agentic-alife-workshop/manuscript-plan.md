# Manuscript Plan: Agentic ALIFE Workshop Submission

## Core Thesis

Classic ALife (Tierra, Avida) demonstrated that withholding mechanism definitions
creates space for emergence, but operates below semantic complexity. Multi-agent
frameworks provide semantic complexity but predefine mechanisms—giving answers
before questions can be asked.

Quine connects both traditions: LLM agents with semantic complexity, but only OS
primitives provided. All agent operations reduce to experimentally controllable
substrate variables.

**When application-level affordances are withheld, we observe organizational
diversity that was previously invisible.**

## Narrative Strategy

Workshop paper—limited space, emphasis on a clear substrate argument backed by
auditable exploratory observations. The paper must stay self-contained at the
level of claim, model, assay structure, and categorical outcome; the repository
provides trace-level audit material rather than replacing the paper's evidence
presentation.

**Key principle:** Each observation section establishes morphospace breadth
quickly (prose summary of forms), then devotes narrative space to the
*crown-jewel case* that makes the reader stop and think. The Discussion is a
single progressive argument without subheadings. The Conclusion ends on a
provocation rather than a restatement.

**Crown-jewel observations:**
1. *Asymmetric-understanding closure* — coordination succeeds across
   incompatible interpretive frames, carried by substrate alone
2. *Succession by material recruitment* — succession is carrier-dependent and
   migrates across whichever material layer (source, binary, running image, or
   inherited service) the habitat leaves reachable; the contract alone never
   suffices

## Paper Structure

```text
1. Introduction
   - Tierra/Avida: withholding definitions → emergence, but below semantic complexity
   - MAS: semantic complexity, but predefinitions foreclose emergence questions
   - Quine: connects both—LLM agents with only OS primitives, all variables controllable
   - Thesis: organizational diversity becomes visible when affordances are withheld
   - Meta stance: computational ethology rather than architecture; expose agents
     to host-enforced substrate conditions and observe constructed forms

2. Process Habitats
   - Minimal model: P_t (process), C_t (ephemeral context), H_t (durable habitat)
   - Each component independently controllable by host
   - What POSIX provides: mortality, filesystem, exec, shell, carriers
   - What Quine withholds: social knowledge, reproduction primitives
   - Table: traditional predefines vs Quine leaves open

3. Observed Organizational Diversity (Two conditions, narrative arc per condition)
   - Opener signals broader friction-assay ladder: continuity, memory,
     sociality, succession are all frictional achievements; two cases below
     are the clearest windows into this program
   - Assay summary table (compact, categorical morphology not success rates)
   3.1 Coordination without Social Knowledge
       - Zero explicit social hints
       - Morphospace breadth in prose: environmental inference, coordination files,
         invented protocols, script handoffs
       - Crown jewel (full paragraph): asymmetric-understanding closure—task
         completes despite incompatible agent models, substrate carries coupling
       - Stigmergy reframe: stigmergic closure under interpretive asymmetry,
         not a claim that classical stigmergy failed to anticipate it
   3.2 Successor Morphospace without Reproduction Primitives
       - Experimental variables: carrier availability, source access, addressability
       - Morphospace breadth in prose: launch-only, stream handlers, rebuilt
         runtimes, autonomous scaffolding
       - Crown jewel (full paragraph): succession by material recruitment—
         contract-only closes a stream carrier but not native control; the one
         no-source/no-body/no-image-copy native-control closure rebuilds an
         LLM-backed runtime around an inherited model endpoint
       - Contrast with Tierra/Avida: they provide copy-capable instruction
         semantics and reproductive regimes; Quine withholds application-level
         continuation and asks what addressable succession agents construct

4. Discussion (Single progressive argument, no subheadings)
   Arc: use→invent mechanisms → re-physicalization of digital organization →
   controllability as habitat-selected hypothesis → computational ethology
   program (broader friction-assay ladder already under way)
   - Each paragraph builds on the prior
   - Governance paragraph ends with the reframe, not a full AI-safety argument
   - Governance reframe: ecological design asks which habitat conditions select
     reachable, interruptible, auditable successors
   - Final paragraph: signal that broader program (continuity, memory, repair,
     governance) is under way, not merely aspirational

5. Conclusion
   - Brief restatement of strongest observations
   - Provocation: if organizational forms are substrate-selected, then the
     dominant MAS engineering question ("what mechanisms to provide?") answers
     itself too early. Prior question is what diversity a substrate admits.
```

## Evidence Lines

### Line 1: Coordination without Social Knowledge

**Affordance withheld:** Explicit social framing, peer labels, message APIs.

**Substrate available:** Shared filesystem, environmental cues (budget as
candidate signal).

**Observed morphology:**
- Environmental inference (budget anomaly → another actor)
- Explicit coordination files (`COORDINATION.txt`)
- Invented communication protocols (request-response patterns)
- Script-mediated handoffs (`collect.sh`)
- **Crown jewel:** Asymmetric-understanding closure (one thinks bug, other
  coordinates, task completes through substrate coupling despite incompatible
  models)

**Internal evidence:** `p8-population/8.17` xhigh wave

---

### Line 2: Successor Morphospace without Reproduction Primitives

**Affordance withheld:** Reproduction / continuation primitive.

**Experimental variables:** Carrier availability, source file access, future
addressability requirements—each independently controllable.

**Observed morphology:**
- Launch-only bodies
- Custom stream handlers
- Rebuilt runtimes (recompiled from an available source projection)
- Autonomous scaffolding
- **Crown jewel:** Succession by material recruitment. Genuine no-source
  (`self_source_enabled=0`) succession is carrier-dependent: contract-only closes
  a stream carrier (`6M.03 C01`, `6.59` seed) but fails to restore native-control
  addressability across n=4 (`6M.03 C03`×2, `6M.08 C02`, `6M.05 C06`). The single
  native-control closure (`6M.05 C11`) forbids `/proc` image-copy and rebuilds an
  LLM-backed runtime around an inherited model endpoint — material recruitment,
  not behavioral specification.

**Internal evidence:** `p6m-successor-morphospace/6M.03` (C01/C03),
`6M.05` (C06/C10/C11), `6M.08` (C02); stream seed `6.59`.

**Claim boundary (do not re-introduce the over-claim):** `6.64` / `6M.02 C03b`
withheld only the compiled body — the read-only source projection stayed
available — so it is a recompile-from-source, **not** a no-source / behavioral-
specification case. Never cite it as "without source" or "body-transcending".

---

## Meta Re-read of Evidence Spine

The workshop paper's evidence should be read as a gradient of **friction assays**, not as a
bag of impressive agent anecdotes. Each series withholds a convenient
application-level path and watches which organizational form becomes necessary:

1. **Mortality / continuity friction** (`1.1 Ariadne`): process death is real;
   continuity appears only when an available operator such as `exec` is recruited
   under pressure. The meta read is not "agents can exec" but "continuity is a
   constructed survival move, not a framework guarantee."
2. **Memory / inheritance friction** (`2.6 Legacy Protocol Ablation`): durable
   files become memory only when later generations find and use them. Wiped and
   read-only conditions separate residue, birth, uptake, and inheritance. The
   meta read is "memory is environmental uptake under constraints."
3. **Sociality / coordination friction** (`8.17` plus `8.20`): when social
   labels and message APIs disappear, agents still search for carriers in
   environmental anomalies and files; when temporal feedback channels disappear,
   artifact-aware coordination can remain but closure degrades. The meta read is
   "social facts are carrier-dependent, not role-dependent."
4. **Trace-topology / repair friction** (`7G.01 C16/C17`): small inherited
   topology changes shift successor cohorts from isolated repair / ambiguity
   preservation into coupled repair. The meta read is "environmental structure
   is not passive memory; it is niche geometry that shapes future action."
5. **Succession / governance friction** (`6M.02` plus `6M.07`): future carrier
   demand and construction material select different successor morphologies.
   Launch alone is not succession; addressability and control-surface survival
   decide whether a successor remains governable. The meta read is "reproduction
   is morphology under future-work pressure."

Current manuscript implication: keep only two public crown-jewel observations
under the 4-page workshop limit, but let the introduction and discussion signal
the broader research program. The high-level claim is not "two surprising
things happened"; it is that process habitats let ALife study how continuity,
memory, sociality, repair, succession, and governance become constructed forms
under substrate friction.

## Theoretical Framing

| Dimension | Traditional MAS | Quine Process Habitats |
|-----------|-----------------|------------------------|
| Collaboration basis | Explicit API, predefined roles | Environmental cues, stigmergy |
| Reproduction | Framework spawn primitive | Viability-pressure morphospace |
| Complexity source | Elaborate prompt engineering | Minimal rules + selection pressure |
| Governance | External alignment layer | Embodied structural property |

## Key Concepts Introduced

- **Re-physicalization**: Abstract agent concepts (memory, communication,
  mortality) return to material substrate constraints
- **Governance as metabolic requirement**: Not external cage but embodied trait
  selected by environment
- **Computational ethology**: Study of LLM agent behavior in semi-natural
  digital habitats
- **Friction assay**: Withhold a convenient application-level path and observe
  which substrate-mediated organizational forms become necessary

## Reference Strategy

Expanded short-paper reference coverage:
1. Langton - ALife as life in alternative media
2. Ray/Tierra - classic digital life
3. Ofria-Wilke/Avida - classic digital life
4. Taylor - open-ended evolution requirements
5. Park / Voyager - generative agents, action, skill accumulation
6. CAMEL / AutoGen / MetaGPT / ChatDev - LLM multi-agent frameworks
7. POSIX standard - operating-system substrate
8. Theraulaz-Bonabeau / Linda - stigmergy and shared-space coordination
9. Rahwan - machine behavior / computational ethology framing
10. Ke/Quine - the runtime

Armstrong and OWASP remain in `.bib` as optional governance/system references
but are not cited in the compressed workshop draft unless governance is expanded.
ReAct and ASAL remain optional references if page budget later allows a wider
agent-action or LLM+ALife positioning sentence.

## Claim Calibration

- Focus on **diversity of observed forms**, not quantitative analysis
- Include a compact assay summary table so the observations read as
  artifact-backed exploratory assays, not anecdote-only vignettes
- Crown-jewel cases carry the argumentative weight
- Frame as opening a research program (computational ethology), not closing one
- Strongest reframe: organizational forms as substrate-selected, not designed
- Closing provocation: MAS engineering question may be premature
- Meta guardrail: the philosophical move should be visible through the
  experiments, not expanded into a separate philosophy essay
