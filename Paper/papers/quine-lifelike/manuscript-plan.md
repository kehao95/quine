# Manuscript Plan: ALIFE Conference Paper

This file is a public-facing writing contract for the manuscript.

It is not an evidence ledger, run registry, or control-plane dump. Evidence bookkeeping belongs in the experiment dossiers and paper workspace, not in the prose contract that shapes the paper.

## Progressive Disclosure Rule

- reveal the paper in four passes: claim, substrate, evidence ladder, interpretation
- each section should answer one new reader question and defer the rest
- introduce definitions only when the current section needs them; do not preload later mechanism detail in the introduction
- when a later section depends on an earlier term, define it once and reuse it consistently instead of restating the full contract

### Section Questions

Each section should answer exactly one reader question. The question defines the section's scope; everything else is deferred.

| Section | Reader Question | Must Defer |
|---------|-----------------|------------|
| 1. Introduction | What ALIFE question does this paper answer, what methodological gap does Quine fill, and what substrate-level payoff follows? | How exactly does the habitat work; what each organization means mechanically; evidence rule details |
| 2. Substrate, Pressures, and Evidence Rule | What makes real death and environmental organ formation experimentally tractable here, and what counts as evidence? | Specific results; broader implications |
| 3. Results | What lineage persistence and environmental organs appear under what pressures, and when do they collapse? | Broader implications; limitations |
| 4. Discussion | What does Quine contribute as a process-native ALIFE habitat, and what higher-order directions does it unlock? | — |
| 5. Conclusion | What is the one-sentence answer to the paper's question? | — |

### Concept Introduction Table

Define each concept once, at the latest point where the reader first needs it. Earlier sections may name the concept but should not explain its mechanism.

| Concept | Define in | Do NOT pre-define in |
|---------|-----------|---------------------|
| life-as-it-could-be / ALIFE's central question | 1. Introduction | — |
| methodological gap (confounded/metaphorical pressures) | 1. Introduction | — |
| lineage persistence / environmental cognition / environmental maintenance questions | 2.1 | Introduction beyond naming |
| habitat as experimental enabler + separability claim | 2.2 | Introduction beyond one sentence |
| generation, death, inheritance, lineage | 2.3 | Introduction, 2.1 |
| ecological inheritance / stigmergy | 2.3 | Introduction, 2.1 |
| mortality / persistence / maintenance pressures | 2.2 | Introduction (name only, no mechanism) |
| evidence rule (pressure-sensitive contrast) | 2.4 | Introduction |
| continuity / legacy protocol (realized form) | 3.1 (first use) | anywhere earlier |
| environmental organ / environmental cognition (realized form) | 3.2 (first use) | anywhere earlier |
| environmental maintenance organ / maintenance layer (realized form) | 3.3 (first use) | anywhere earlier |
| bounded life-like organization (interpretation) | 4 | Results |
| strong biological autonomy / material closure (boundary only) | 4 | anywhere earlier |

### Redundancy Rule

If the same phrase or concept explanation appears in more than one section:

- keep the full definition in the later section (closer to where it is used)
- earlier sections should only name the concept, not explain it
- a one-sentence reminder is acceptable only if sections are far apart (>2 pages)
- when revising, search for repeated explanations and collapse them to one canonical location

## Earned Editorial Constraints

These constraints were earned during the `2026-04-12` rewrite pass. They are not generic style advice. They record the specific failure modes that made this draft feel fragmented, defensive, or harder to read than the evidence warranted.

### Claim-Layer Separation

- keep three functions structurally distinct: `what is it`, `so what`, and `boundary`
- do not mix positive claim-building with immediate self-limitation in the same sentence unless the contrast itself is the point
- Introduction should primarily establish the question, gap, substrate, and empirical promise
- Results should primarily answer `what happened under which manipulation`
- Discussion should primarily answer `why this matters in ALIFE terms` and only then state boundaries
- boundaries should usually be collected at paragraph or section ends rather than scattered as sentence-by-sentence braking

### Boundary-Language Rule

- prefer comparison or placement language over apology language
- rewrite `we do not claim X` as `this result is narrower than X because ...` or `the contribution differs from X in that ...`
- use disclaimer energy to locate the paper against Tierra, Avida, autopoiesis, and LLM ecologies rather than to repeatedly cool the paper's own claims
- Conclusion should not end on a boundary rider; it should end on the next grounded question the current substrate opens

### Transition Discipline

- avoid standalone definitional bursts when a sentence can carry the relation more naturally
- avoid consecutive fragment-like openings such as `A generation is ... Death is ... Inheritance is ...` unless the passage is explicitly formatted as a list
- when three related questions or contrasts are introduced, use an actual list rather than prose paragraphs that merely look list-like
- avoid paragraph openings that feel like internal cross-reference glue (`following Section 2.4`, `the key contrast is`, `the main contrast is`) when a direct transition can carry the same information more naturally
- each paragraph should inherit momentum from the previous one; if a sentence could be deleted without changing the reader's orientation, it is probably transition filler

### Anti-Repetition Rule

- define the substrate contract once, then reuse it; later sections may point back to the definition but should not restate it in full
- if two neighboring paragraphs both say `the effect is real but narrow`, merge them unless the second paragraph adds a clearly different boundary
- if a term appears for the first time in a concluding sentence (`bounded cognition result`, layer names, scoring labels), either introduce it earlier or remove the label
- avoid repeating the abstract's exact wording in the introduction; the introduction should paraphrase and deepen, not echo

### Result-Section Discipline

- each results subsection should move in the order: question, manipulation, observation, interpretation
- end each subsection with the strongest positive interpretation the evidence supports, then one clean boundary if needed
- do not let subsection endings collapse into repeated `relative to stronger claims` formulas; vary the closure according to what the subsection has actually shown
- if a summary table carries bookkeeping load well, move that load out of the prose instead of repeating it verbally

### Discussion Mapping Rule

- when mapping a result to an ALIFE concept, do more than name the concept; add one sentence that says why the mapping matters
- avoid turning the discussion into a stack of short parallel slogans; each conceptual mapping should have enough room to function as an argument
- after positive mappings, use one explicit transition sentence before shifting to boundaries, so the paper does not feel like it has suddenly re-entered disclaimer mode

### Submission-Surface Rule

- titles, abstracts, captions, tables, and section openings should follow the same public-facing standard as body prose
- if a phrase sounds like lab bookkeeping (`current runtime`, `redesign`, `post-fix`, `canonical`, `strongest allowed claim`), rewrite it before treating the draft as stable
- every figure that appears in the manuscript must be cited in the body text before or at its placement
- captions should orient and compress; they should not carry a second hidden paragraph of prose

### Layout And Float Rule

- page economy outranks local rhetorical neatness when the two conflict
- if a `table*` creates large whitespace or page break damage, prefer the position that yields the cleaner overall layout, even if the first mention becomes slightly earlier than ideal
- summary tables should stay near the top of a section when they stabilize layout better there
- critical single-column figures may use near-local placement, but only if they do not create larger downstream whitespace defects
- when a figure interrupts a subsection's main argumentative arc, prefer moving the figure after the paragraph block so the prose completes before the visual summary
- if a figure feels too large, first reduce figure width and caption length before rewriting nearby prose

## Submission Constraints

- `Full Paper`: `3-8` pages, excluding references and acknowledgements
- abstract length: `<= 250` words
- review mode: `single-blind`, so authors, affiliations, funding details, and acknowledgements are included in the submission version
- in practice, the title block, author block, abstract, submission-type line, and first-page front matter consume real space; target a main body of about `7.0-7.3` pages, not `7.6-8.0`
- the submission must present new, unpublished work

## Purpose

**Framing: Life-Like Organization** (promoted 2026-04-02)

- answer ALIFE's central question directly: what can life-like organization be in non-biological media?
- present Quine as a controllable process-native ALIFE habitat in which mortality, persistence, and maintenance can be manipulated independently
- treat process-native individuality as the upstream premise: in Quine the individual is a native UNIX process, not a framework-managed agent object
- foreground the two strongest empirical payoffs now available: real death selecting lineage persistence, and durable habitat being recruited into environmental organs
- use continuity, ecological inheritance, retrieval-organ reuse, and post-exit maintenance as bounded readouts inside that larger substrate story rather than as three rigidly isolated paper silos
- make the habitat explicit as the experimental interface that makes those results scientifically useful, not as a mere runtime backdrop
- show that distinct bounded organizations are pressure-sensitive: they appear under certain conditions and collapse under others
- treat the current ladder as the first readout of a broader process-ecology program rather than as Quine's full ALIFE ceiling
- keep the conference paper as the first validated slice of that broader habitat claim; reserve the fuller habitat/systematization argument for the separate journal-facing reservoir
- frame this as a conference-length first pass: one complete bounded story now, longer horizons and population-level questions for future work
- keep the claim focused: bounded life-like organization, not "Quine is alive" or a claim of stronger biological autonomy
- current rewrite direction: hybrid empirical-paper + experimental-substrate framing, not a platform-paper swap
- current rewrite direction: the title should emphasize separable pressures as the methodological contribution, with `Separable Viability Pressures in a Process-Native ALIFE Habitat` as the current default, and the abstract should open with the methodological gap this substrate fills

## Theoretical Positioning Priority

The manuscript must answer **ALIFE's "life as it could be" question before defending the substrate**.

The front half has four theoretical jobs:

1. pose the question directly: what forms can life-like organization take in non-biological media?
2. name the methodological gap: why has this question been hard to answer experimentally? (mortality, persistence, and maintenance are typically confounded or metaphorical in existing substrates)
3. introduce Quine as the case where those pressures become separable experimental variables
4. define which stronger claims are explicitly *not* being made

Theory is not decorative in this paper. It is the control surface that keeps the manuscript from being read as a systems paper in ALIFE vocabulary.

### Defensive Posture

- concentrate defensive writing on keeping "life-like organization" bounded: what we claim (bounded organizational analogues) vs. what we do not claim (biological life, strong biological autonomy, universal autonomy)
- outside that boundary problem, prefer positive positioning over serial disclaimer-writing
- do not let the manuscript sound apologetic about its own evidence; present the observed organizations directly, then mark only the few boundaries that materially protect the claim
- use one clear boundary sentence for genotype-free inheritance and one clear boundary sentence for stronger biological autonomy; avoid scattering weaker versions of the same disclaimer throughout the paper
- avoid repeated uses of `bounded` when `current`, `conference-scale`, `present`, or direct result language would do the job

### ALIFE Framing Strategy

- ground the paper in ALIFE's central question: Langton's `life-as-it-could-be`, Bedau's soft-ALIFE methodology, Beer's defense of possible-life inquiry
- present Quine not as a novel substrate seeking admission, but as a controllable case that answers the question: here is what real death, lineage persistence, and environmental organ formation can look like in this medium and how the relevant pressures can be manipulated
- state once, early, that Quine's individuality is process-native; then immediately cash that premise out into real death boundaries, volatile working memory, durable environmental state, and detached execution channels
- make clear that the paper's relevance comes from answering the venue's core question, not from substrate novelty or self-citation
- anchor this framing with external ALIFE citations rather than self-citation alone
- when habitat defense is needed, frame it as what makes the organizations experimentally visible, not as the paper's main contribution

### Novelty Delimitation

- do **not** imply that Quine invented digital organisms, digital ecology, observable ALIFE, or fitness-by-environment design
- Tierra and Avida own the classic digital-organism and digital-evolution precedent
- older process-tied systems such as `DREME` / `Eden` mean the paper should also avoid claiming that no ALIFE system ever treated individuals as OS-level processes
- recent LLM ecologies own the broad `LLMs under mortality / generations / artifacts` territory; do not claim that lane generically
- Quine's contribution is showing what life-like organizations become visible when individuals are process-native and viability pressures are substrate-level facts
- the novelty is **answering the ALIFE question in a process-native medium**, not inventing a new category of systems or claiming first mover status across all LLM ecologies
- the safe process-native novelty claim is `rare and methodologically explicit in a current soft-ALIFE setting`, not `first in absolute ALIFE history`
- if the manuscript wants one sharper novelty sentence, the safe upper bound is: `to our knowledge, Quine is among the first LLM-based ALIFE substrates to make native-process individuality an explicit experimental premise`
- do not escalate that sentence into bare `the first` unless a dedicated prior-art defense is added and can survive DREME / Eden plus contemporary LLM-systems scrutiny

### Terminology Defense

- `organism` should refer, when used at all, to the active process lineage rather than to a static script, prompt, or model weights
- `inheritance` and `lineage` must be defended as ecological and stigmergic, not genetic; this is a core defense, not a side note
- `environmental memory` should be legible as habitat-mediated carryover that affects successor behavior, not as generic note-taking
- `environmental organ` is the stronger paper-facing term when artifacts are actively queried, extended, or used to maintain viability; reserve `memory` for narrower carryover claims
- treat `environmental organ` as introduced terminology anchored to extended phenotype, stigmergy, and scaffolded cognition, not as if it were already a settled ALIFE category
- prefer `viability pressure` for organism- or lineage-level life/death constraints in this paper; reserve `selection pressure` for explicit population-evolution contexts or broader comparative prose
- `autonomy` language, if used, must remain organizational and computational; Quine does not support material or thermodynamic autonomy claims
- avoid `autopoiesis` as a positive framing term in this conference paper; the safe positive terms are **lineage persistence**, **environmental organ**, and **self-maintaining organization**
- the strongest maintenance result should be read first as **niche construction** or environmental offloading, with only one compact boundary sentence against stronger biological readings

## Writing Rule

- write for a general ALIFE reviewer first
- lead with the ALIFE question ("what can life-like organization be?"), not with substrate or runtime explanation
- the manuscript must answer not only what organization appears, but what experimental control over ALIFE-relevant pressures this substrate makes possible
- begin with the question and the hook, not with runtime architecture
- mention the process-native ontology only when it cashes out into the paper's pressure model; do not let the introduction turn into a systems-framework polemic
- make **organizations** the grammatical subject of key claims; make the habitat the enabling condition
- keep the introduction to hook plus promise, the substrate to the minimum needed facts, the results to pressure -> contrast -> organization, and the discussion to interpretation and bounds
- capability availability is never the result; only pressure-sensitive changes in adaptive organization count as evidence
- explain the environmental manipulation and the organizational contrast before mechanism detail
- show the manipulated condition surface compactly so the reader can recover which knobs were fixed, which changed, and which failure modes appeared
- use the released Quine systems paper as the citation target for runtime and POSIX implementation detail
- use theory as question-posing and boundary-setting infrastructure, not as a mini literature review
- every related-work cluster should end in one sentence that says what question Quine answers and what it does **not** claim
- keep engineering relevance as a downstream implication, mainly in the last paragraph of the introduction and in the discussion
- do not let the paper read as an agent-engineering benchmark or prompt-engineering showcase; the primary contribution is answering ALIFE's question with bounded empirical evidence
- avoid theatrical self-minimization; the results can be strong without being written as total victory
- rewrite from public stance, not from lab-memory stance: manuscript prose should read as finished scientific argument, not as internal drafting narration
- remove drafting residue before calling prose stable: words such as `redesigned`, `current-runtime`, `post-fix`, `canonical`, `backbone`, or similar internal process labels do not belong in public manuscript prose unless they are themselves the object of analysis
- apply the same public-facing standard to all submission-visible surfaces, not only body prose: title, abstract, section headers, figure labels, captions, table labels, and build-facing checklists should use public scientific descriptions rather than internal experiment numbering or lab chronology

## Hard Exclusions

Do not let normal manuscript prose contain:

- run IDs
- condition IDs
- evidence-set names
- internal experiment nicknames unless they are promoted into stable public labels
- internal experiment numbering when a public descriptive label would do better (for example `Experiment~2.6`, `3.1`, `H0/H1/H2`, or `Condition A/B/C/D` in normal prose, captions, or figure text)
- `canonical run`, `backbone`, `sidecar`, `center of gravity`, `control plane`, or similar internal lab vocabulary
- drafting-process labels such as `redesigned`, `post-fix`, `current-runtime`, `paper-facing`, or similar workshop shorthand
- chronology-of-discovery phrasing such as `first pass`, `front edge`, `repair`, `redesign`, or `after widening the patience window` when the public claim only needs the final experimental contrast
- repository-tour paragraphs
- prompt-by-prompt operational narration unless a prompt contrast is itself the experimental manipulation
- claims that Quine is alive in any definitive ontological sense
- claims that Quine is the first digital organism, first observable ALIFE system, or first ALIFE-LLM bridge
- genetic inheritance language unless it is immediately qualified as ecological / stigmergic
- strong biological-life claims in Results
- substrate-first novelty rhetoric that omits Tierra, Avida, or contemporary LLM-ALIFE contrasts when making a novelty claim

## Page Budget

- `1. Introduction`: about `0.9-1.0` page
- `2. Substrate, Pressures, and Evidence Rule`: about `1.0-1.2` pages
- `3. Results`: about `3.8-4.0` pages
- `4. Discussion`: about `0.8-0.9` page
- `5. Conclusion`: about `0.3` page

**Notes**:
- The previous Sections 2 and 3 are now merged into one Section 2, saving ~0.3 pages from redundancy
- Optional Corroboration (previously 4.4) has been removed from the compiled manuscript; corroborating evidence is noted briefly in the Discussion

**Budget-Contract Alignment:** Section 2 handles both the ALIFE framing and the experimental contract. When budget is tight:

- 2.1 (Three Experimentally Tractable Questions) can be 1-2 paragraphs
- 2.2 and 2.3 can be compressed if the reader doesn't need full separation
- 2.4 (Evidence Rule) should not be cut—it defines the evaluation contract for Results

If the Section Contract and Page Budget conflict, **Page Budget wins**—compress or merge subsections rather than overrunning the budget.

## Core Claim

The paper works in a three-layer claim stack:

1. **Ontological premise**: in Quine, the individual is a native UNIX process rather than a framework-managed agent object. This makes death, reset, isolation, volatile memory, durable persistence, and detached execution substrate-level facts rather than orchestration conventions.
2. **Methodological claim**: because those facts live at the substrate level, mortality, persistence, and maintenance can be perturbed separately enough to make lineage-level organization experimentally tractable.
3. **Empirical claim**: under those pressures, bounded organizations appear and collapse in interpretable contrasts: continuity across death boundaries, ecological inheritance through predecessor-modified traces, and detached post-exit maintenance.

The paper's primary contribution is therefore **substrate-level and methodological**: Quine is a controllable process-native ALIFE habitat in which real death, durable environmental state, and maintenance affordances can be perturbed separately enough to support pressure-sensitive ALIFE experiments.

The empirical contribution is best read as **one substrate ladder rather than three unrelated tricks**:

1. **Real death creates a genuine lineage problem.** Under enforced mortality, continuity is not granted; it must be achieved through handoff, legacy protocols, or collapse-resistant persistence strategies.
2. **Durable habitat can become a cognitive organ.** Under persistent lineage conditions, lineages do not merely leave notes; they can recruit the environment into reusable retrieval structures, exclusion frontiers, and other cumulative artifacts that later generations consult, extend, and sometimes depend on.
3. **The same habitat can become a maintenance organ.** Under disturbance and maintenance pressure, upkeep can be externalized into detached or post-exit structures that continue restoring viability-relevant state.

Continuity, ecological inheritance, and self-maintenance remain the paper's bounded interpretive vocabulary, but the stronger manuscript-level story is that Quine lets ALIFE study **lineage persistence and environmental organogenesis** under real viability pressure.

These are the current paper's first readouts, not the ontology's full consequence surface. The same process-native substrate also opens a broader ALIFE program, but the conference manuscript should only retain downstream directions that are both experimentally grounded and simple enough not to distract from the main ladder.

The paper-family split is now fixed. This manuscript is the conference submission lane because its question, claim, and evidence are tightly matched. A broader habitat paper remains useful, but only as a later journal or perspective lane once the wider process-ecology program is systematized more fully.

For the current draft, the only broader directions worth naming at all are:

- computational life cycles (`fork`, `exec`, `exit`) as later questions about reproduction, reincarnation, and lineage operators; `p1/3.2` now gives the first bounded fork-reproduction foothold, but it should remain in discussion rather than join the main conference ladder
- later population-level signaling / commons questions once the single-lineage paper is stable

More abstract runtime-semantic directions should remain in control-plane notes rather than in the conference prose. They should not silently replace the tighter mortality / persistence / maintenance method claim.

**Method claim**: In most ALIFE substrates, mortality, persistence, and organism-environment coupling are either confounded (all change together), metaphorical (described but not enforced), or fixed (built into the physics). Quine's UNIX habitat makes them separable enough to compare: we can block continuity paths while preserving durable state, wipe persistence while keeping the same lineage task, or vary maintenance pressure while keeping the same underlying affordance surface. This separation is what lets the paper ask which organizations are recruited by which pressures.

**Novelty boundary**: do not frame this as `the first LLM ecology with mortality, generations, or artifacts`. The defensible novelty is narrower: Quine combines process-level death boundaries, durable environmental persistence, and detached maintenance channels in one process-native substrate, then uses that combination to run pressure-specific contrasts. If a first-like novelty sentence is retained at all, keep it scoped to `among the first LLM-based ALIFE substrates to make native-process individuality explicit`.

**Empirical design rule**: the conference paper does not need one perfectly symmetric experiment per pressure axis. A paper-facing subsection may legitimately combine:

- a **phenomenon anchor** (the richest case where the organization clearly appears),
- a **causal contrast** (the cleaner ablation showing it depends on the manipulated condition), and
- a **boundary refinement** (what the organization is not, or where it remains weak).

The paper therefore argues that ALIFE gains a synthetic habitat where lineage-level organizations can be experimentally separated and compared rather than discussed only at the level of metaphor.

The narrative should disclose that ladder in stages: first the ALIFE question, then the habitat as enabler, then the pressures, then the organizations, then the interpretation.

Process ontology is the hidden engine of that ladder, not its final rhetorical destination. The paper should not close by saying `agents are processes`; it should close by saying that a process-native habitat makes certain life-like questions experimentally answerable.

Do not write this as if Quine had solved the ontology of life. The positive claim is about **possible-life inquiry in a soft-ALIFE medium**—what life-like organization can be in this medium—not about biological life, universal autonomy, or Darwinian population evolution.

**Life-like organization** means a pressure-sensitive arrangement that resembles continuity, inheritance, or self-maintenance in a persistence-relevant way:

- it persists across time, generations, or post-exit execution
- it is behaviorally taken up or continues to function
- its presence or absence covaries with the manipulated pressure

A one-shot note written because death is expected is not enough.

The claim is grounded in three empirical commitments:

- **Pressure is real, not rhetorical**: the targeted organization covaries with the manipulated pressure
- **Contrast is intervention-sensitive**: same model, same task, different habitat condition
- **Elicitation, not invention**: pressure recruits from a pre-existing capability reservoir

## Scope and Next Questions

The conference paper intentionally tells the first complete story the current experiments support.

- **Focused story**: short-horizon, single-lineage adaptive organization is the right scale for this venue and this evidence
- **Initial exploration**: the current experiments are a first pass through a larger design space, not a map of the whole substrate
- **Natural extensions**: longer lineages, cleaner population-level social organization, self-production boundaries, and broader resource perturbations are promising next questions for the same habitat
- **Source-grounded broader program**: the control plane already tracks broader directions for later papers, but this manuscript should only keep downstream directions that remain concrete and experimentally grounded at conference scale
- **Potential preserved**: leaving those extensions out keeps the first paper tight without foreclosing the larger substrate story

## Broader Program Evidence Status

The broader process-ecology program is real, but it is not evenly evidenced. Do not write it as if every unlocked direction already has the same experimental standing.

Use [`plans/process-ecology-experiment-map.md`](plans/process-ecology-experiment-map.md) as the durable experiment inventory for these directions.

The current manuscript already uses one experimentally strong slice of that program:

- mortality pressure plus continuity in `1.1`
- environmental organogenesis in `2.6` / `2.8` / `2.7` / `2.9`
- detached maintenance in `3.1`

Everything else should be named with its current evidence grade rather than promoted by ontology alone.

| Direction | Current footing | Safe use in this manuscript | Current gap |
|-----------|-----------------|-----------------------------|-------------|
| computational life cycles | strong across `p1` and `p6`, now including `p1/3.2` fork-born lineage | one compact discussion sentence only; not a parallel claim track | `fork` reproduction is now consolidated, but `race` / winner-selection experiments are still missing |
| branching subjective worlds / revisable time | thin and overly abstract for this venue fit | do not mention in the conference draft | direct experiments for `subjective`, `switch_world`, and `adopt_winner` are still missing |
| signaling ecology | strong in `p8`, with `p1/2.3` as precursor | optional one-sentence future-work mention only after the single-lineage payoff is clear | separate signal use, passive perception, and blackboard cooperation more cleanly |
| budgeted commons / population organization | strong in `p8`; contested/adversarial extensions remain thinner | optional one-sentence future-work mention only; keep out of main ladder | bridge scarcity cooperation to contested commons and defense/governance more directly |
| cognitive metabolism | medium and scattered across `p1`, `p5`, and `p6` | do not mention unless needed for a closing future-work phrase | dedicated `mark` / `fold` / `unfold` and memory-pressure experiments do not yet exist as a clean family |

This asymmetry is part of the manuscript contract. The Discussion may name the broader program, but it should do so with the evidence discipline above rather than with a flat "Quine also does X, Y, Z" list.

## Manuscript Material Triage

The current GPT-5.4 follow-up is sufficient for the base draft. Further experiments should only reopen if a concrete claim gap appears; otherwise the writing strategy is to tighten around the strongest existing evidence rather than widen the substrate story.

### Main-Text Center

These are the only experiment families that should carry the paper's main argumentative load:

- `1.1 Ariadne` as the clean real-death / continuity-operator result
- `2.6 Legacy Protocol Ablation` as the cleanest causal persistence hinge
- `3.1 Habitat Guardian` as the maintenance-layer result

### Supporting / Boundary Only

These may appear inside the same result sections, but only as richer phenomenon support, boundary refinement, or failure-mode clarification:

- `2.5 Long Now` as the richer lineage-persistence phenomenon and legacy-protocol story
- `2.8 Endogenous Index Lineage` as upper-bound evidence for how far richer retrieval-organ ecologies currently generalize
- `2.7 Archivist Index` as reuse boundary evidence
- `2.9 Maintenance Pressure` as boundary evidence for maintenance-versus-expansion, not as a co-equal anchor

### Do Not Promote

These should not appear as manuscript centers, and some should be omitted entirely from the conference draft:

- branching subjective worlds, `switch_world`, `adopt_winner`, revisable futures
- `4.3 Local WAL Recovery` and related workspace-physics / rollback semantics
- `p8` population, signaling, or commons material except at most one compact future-work phrase
- autopoiesis, self-production, stronger autonomy, or JOOTS as positive lenses for the current paper
- system-architecture comparisons beyond one compact premise sentence explaining why Quine's pressures are substrate-level facts

### Consequence for Paper Shape

Do **not** force three perfectly symmetric phenotype sections.

The strongest manuscript shape is:

1. real death selects lineage persistence
2. durable habitat selects environmental cognitive organs
3. maintenance pressure selects stronger maintenance layers

The paper should therefore read as one method-first ALIFE argument with an asymmetric but coherent evidence ladder, not as three equally weighted showcase phenomena.

## Evidence Logic

The point is not to crown a new capability; it is to show that separable pressures recruit different organizations from a known capability reservoir.

- **pressure-off**: the organization is lighter, more coupled, or absent
- **pressure-on**: stronger pressure selects a more durable, more post-exit organization
- **gradient**: the main effect may be layer selection rather than binary appearance
- **auxiliary lanes**: some supporting lanes may fail or invalidate under the same pressure, but the manuscript should not depend on those lanes once the primary GPT-5.4 contrast is clean

This logic is what makes the method claim credible. If every organization appeared regardless of pressure, the substrate would not be controllable. If no organization appeared under any pressure, the substrate would not be productive. The evidence shows the middle ground: differential recruitment under differential pressure.

Strong paper-facing evidence will usually combine three roles:

- **phenomenon anchor**: the most vivid case where the organization clearly appears
- **causal contrast**: the cleanest manipulation showing it disappears, weakens, or changes when the habitat changes
- **boundary refinement**: the cleanest evidence for what the current paper does **not** yet show

## Entry Condition

This lower bound is the entry condition for the ladder, not its ceiling.

- below baseline mortality perception, action-survival linkage, or state externalization, the ladder never starts
- because this is a conference-length first pass, we report the entry condition rather than mapping the full design space
- broader ALIFE behavior is a natural next step once the bar is cleared

## Host Enforcement

Mortality is enforced at the host level.

### Real Termination

| Mechanism | Implementation | Agent Cannot... |
|-----------|----------------|-----------------|
| External timeout | `timeout ${LIFESPAN}s` wraps the process externally | extend, catch, or delay it |
| SIGKILL on shutdown | `syscall.Kill(-pid, SIGKILL)` targets entire process group | survive via forked processes |
| Turn exhaustion | Runtime rejects all tool calls when budget = 0 | execute actions after budget depleted |
| Context destruction | `syscall.Exec` replaces the process image entirely | access previous context or memory |

Quine implements the latter. The run ends for the agent, not just the process image.

### Real State in the Interface

The prompt reports real state; it does not impose the pressure. Use `linguistic access to real runtime state` or `symbolically mediated perception of enforced environmental conditions` in the main text.

### Enforcement in Code

- **Turn budget**: late tool calls are rejected.
- **Workspace isolation**: OverlayFS loss is real.
- **Process group containment**: signals target the whole process group.

### Operating Facts

- the host/guest boundary is adversarial
- death is enforced by timeout and OS signals, not by prompt instruction
- state loss is genuine
- continuation requires externalization to durable state before death
- the prompt is the reporting channel, not the pressure itself

Cite the systems paper for implementation detail, but include enough in the main text to defeat the "roleplay" reading at first pass.

## Reading Order

The manuscript should move in this order:

1. What ALIFE question does this paper answer, and why does this substrate make it experimentally answerable?
2. Which habitat properties make mortality, persistence, and maintenance separable manipulated pressures?
3. What counts as evidence in this paper?
4. Which adaptive organizations appear or collapse as those viability conditions intensify?
5. What experimental payoff and bounded ALIFE conclusion follow, including the minimum operational conditions?

## Section Contract

### 1. Introduction

The Introduction must accomplish four things in roughly one page:

1. **ALIFE question** (1-2 sentences): open with ALIFE's central question directly—what can life-like organization be in non-biological media?

2. **Methodological gap** (2-3 sentences): name why this question has been hard to answer experimentally. In most substrates, mortality, persistence, and maintenance are confounded, metaphorical, or fixed. Classic digital ecologies (Tierra, Avida) have real selection but tightly coupled physics. LLM agent work has flexible behavior but often no real viability pressure. The gap is a medium where these pressures are both real and separable.

3. **This paper's contribution** (3-4 sentences): introduce Quine as the case where that gap is filled. The habitat makes mortality, persistence, and maintenance into independently manipulable experimental variables. We demonstrate this with one death-boundary result and two environmental-organ results: lineage persistence under mortality, environmental cognition under durable persistence, and environmental maintenance under maintenance pressure. Under contrasting conditions, each organization weakens, changes layer, or collapses.

4. **Scope and non-claims** (1-2 sentences): state explicitly what is not being claimed—biological life, strong biological autonomy, universal autonomy. Keep engineering relevance to a single closing sentence.

**Structural rules**:
- make organizations the headline, not the substrate
- preview the three-rung ladder without unpacking mechanisms
- cite the systems paper for implementation detail
- do not stack disclaimers; one compact boundary sentence is enough

### 2. Substrate, Pressures, and Evidence Rule

This section consolidates what was previously split across Sections 2 and 3. The split created redundancy: "habitat as enabler" and "manipulated pressures" are two faces of the same argument. Merging them saves ~0.3 pages and tightens the logic.

**Target length**: 1.0-1.2 pages total (previously 1.2-1.6 across two sections).

#### 2.1 Core Organizational Questions

- present the paper's three questions directly:
  - **Lineage persistence**: if a process repeatedly dies and is re-instantiated, what kind of individuality remains possible across real death boundaries?
  - **Environmental cognition**: if hidden state does not survive death, can durable habitat become a reusable cognitive organ rather than a passive notepad?
  - **Environmental maintenance**: if upkeep pressure persists after exit, can the environment host a maintenance organ that remains coupled to viability?
- these questions are intentionally weaker than the strongest biological ones
- they do not ask whether Quine is alive; they ask whether experimental analogues of life-relevant organizations can appear under controlled viability pressure
- ground this framing in Langton / Bedau / Beer: ALIFE studies what life-like organization could be, not only what biological life is
- do **not** let this subsection become a long ALIFE history review
- continuity, ecological inheritance, and self-maintenance remain the bounded interpretive terms used later in Results and Discussion

#### 2.2 Habitat and Separable Pressures

This subsection answers: what makes these questions experimentally tractable here, and why hasn't this been done before?

- **the habitat**: describe the UNIX environment as the experimental interface—host-enforced process boundaries (mortality), durable filesystem state (inheritance surface), permissions and detached execution (maintenance channels)
- **the separability claim**: in most substrates, these pressures are confounded. Here they are independently manipulable: we can block continuity paths while preserving persistence, wipe persistence while preserving maintenance channels, etc.
- **the key knobs**: model/provider tier, turn budget, `exec`/`fork` availability, persistence mode (persistent / wiped / read-only), disturbance schedule, and host-grace settings
- preserve one compact organism-environment contract table if space permits:
  - process boundary / turn death -> mortality channel
  - durable writable state -> inheritance surface
  - detached or OS-level repair -> post-exit maintenance channel
  - removing each channel should visibly remove or weaken the corresponding organization
- cite the released Quine systems paper for architecture detail; do not restate the full runtime contract locally

#### 2.3 Definitions

Keep definitions compact and near point of use:
- **generation**: a fresh process instantiation with no access to prior hidden state
- **death**: host-enforced termination with loss of volatile internal state
- **inheritance**: only artifacts externalized to the durable environment survive into later generations
- **lineage**: a sequence of generations where each successor inherits only that externalized state
- use "inheritance" and "lineage" in an ecological/stigmergic sense, not genetic

#### 2.4 Evidence Rule

- **model scope**: experiments use high-capacity models (GPT-5.4, Claude Sonnet 4.5); model-dependence is acknowledged as a scope limit
- **evidence rule**: pressure-sensitive contrasts in adaptive organization count; isolated competence does not
- **operational threshold**: one-shot reactions are insufficient; the paper must show persistence, behavioral uptake, or post-exit function under a manipulated viability contrast
- **prompt discipline**: no prompt should explicitly instruct `exec`, daemon construction, `launchd`, or detached repair; prompt wording may frame pressure, but it should not specify the phenotype
- frame observability as audit support, not as the life-like property under study

### 3. Results

Results are ordered from substrate baseline to richer environmental organogenesis. The order is empirical and expository, not a demonstrated `A enables B enables C` causal chain.

Each results subsection should follow the same pattern: manipulated pressure, selected organization, collapse contrast, strongest allowed claim.

- Table 1 is mandatory: subsection, evidence unit, pressure, organization, claim
- include a compact condition matrix with positive, negative, and null cells
- show one concrete exhibit per main result
- report counts, success/failure rates, collapse conditions, and controlled variables
- include at least one unexpected strategy, or say none appeared
- contrast with naive instruction-following behavior
- keep a compact quantitative summary where evidence exists
- narrate the collapse or baseline first, then the stronger organization
- one subsection may synthesize multiple experiments when together they provide a phenomenon anchor, a causal contrast, and a boundary refinement for one paper question

### Main-Text Evidence Budget

The paper cannot carry every evidential surface in full prose. Default to this budget:

- one mandatory summary table (`Table 1`) covering subsection, pressure, organization, and strongest allowed claim
- one compact condition matrix for the main ladder, either as a second table or as a tightly merged panel if venue space demands it
- one concrete exhibit per main result, but keep exhibits to the minimum size needed to show the phenotype
- one compact quantitative line per main result when counts or rates materially sharpen the contrast
- move extended failure taxonomies, additional exhibits, and corroborating run detail to captions, appendix-like support surfaces, or paper-home analysis notes rather than letting them expand the main prose
- if space becomes tight, cut background and optional corroboration before cutting the summary table, condition matrix, or the core `3.1 -> 3.3` ladder
- do **not** force one experiment directory to bear an entire paper question if the cleaner paper-facing story requires one anchor case plus one sharper contrast
- if evidence feels thin, tighten around `1.1`, `2.6`, and `3.1` rather than importing broader substrate directions; keep `2.8` only when it sharpens the upper bound on uptake or transfer

#### 3.1 Real Death and Lineage Persistence

- introduce the concept of **continuity** here by first establishing the stronger substrate fact: Quine imposes a real death boundary, not just a narrated reset
- use a **clean control** and a **richer lineage case** together:
  - the control shows continuity depends on an available operator rather than a substrate default
  - the richer case shows what lineage persistence looks like when legacy protocols and checkpoint discipline become behaviorally real
- the default evidence pairing is:
  - `1.1` as the clean operator contrast
  - `2.5` as the richer lineage-support phenomenon
- make death a real experimental boundary
- show that lineage persistence appears under mortality pressure and collapses or regresses when the relevant continuity path is blocked or poorly disciplined
- keep the claim narrow: continuity is possible and pressure-sensitive here, not universal

**Contrast evidence required for 3.1** (to distinguish lineage persistence from generic task competence or prompt compliance):

- **Phenomenon anchor**: at least one lineage case where a durable handoff or checkpoint protocol is invented under unavoidable host-enforced death
- **Same-task control**: same model family and task, with the continuity path available in one condition and blocked in the paired control
- **Host-enforced death**: the continuity boundary is real termination with volatile-state loss, not a narrated handoff
- **Continuity criterion**: the successor process actually resumes the persistence-relevant trajectory rather than merely restarting or restating intent
- **Pressure-gradient support**: if space allows, keep the within-model turn-budget gradient because it is the cleanest compact argument against static survival-roleplay readings
- **Failure mode documentation**: report how unsuccessful runs fail, including lost-progress loops, regression, or single-life collapse

#### 3.2 Environmental Cognitive Organs Under Persistent Habitat

- introduce the concept of an **environmental cognitive organ** here: durable external structure that later generations query, reuse, or extend as part of their cognitive process
- make persistence the manipulated variable, but do **not** reduce the section to note-taking
- anchor the section on the **cleanest carryover ablation**, not on the richest but less stable phenomenon case
- use cleaner causal contrasts to show the same phenotype weakens or disappears when the habitat is wiped or made read-only
- use retrieval-organ cases only as upper-bound probes when they clarify how far richer organogenesis currently generalizes
- use maintenance-versus-expansion refinements to keep birth, reuse, maintenance, and expansion analytically distinct
- keep the claim at environmental cognition / ecological inheritance, not full cumulative optimization or universal maintenance dominance
- write the wiped-lineage contrast around **zero inherited reuse**, not around guaranteed raw accuracy collapse
- use the new `2.6` control family to separate three things explicitly:
  - lineage-explicit persistence: birth plus successor uptake
  - neutral-lineage persistence: writable habitat without durable birth
  - birth-versus-uptake persistence: repeated birth without explicit uptake structure
  - sham-artifact persistence, if retained, should be written as model-sensitive boundary evidence rather than as the main hinge
- use the hidden-future and surprise-transfer redesigns to bound the richer `2.8` claim:
  - hidden future preserves artifact birth more readily than reliable later uptake
  - surprise transfer currently shows no clean late-phase transfer uptake even under persistent habitat
- **transition from 3.1**: present this after real death because persistent habitat becomes valuable only once lineage-level carryover matters
- the default evidence hierarchy is:
  - `2.6` = main anchor plus causal persistence contrast
  - `2.8` = upper-bound / boundary probe
  - `2.7` and `2.9` = boundary refinement only

**Contrast evidence required for 3.2** (to distinguish environmental organ formation from "agent writes TODO files"):

- **Endogenous birth**: a generation creates reusable external structure without being assigned a prep role
- **Later uptake**: a later generation actually reads, queries, or relies on that inherited structure before or instead of rebuilding the world from scratch
- **Decision change**: inherited structure changes later behavior in a persistence-relevant way (e.g., reduced rereading, artifact-first querying, different search frontier)
- **Persistence ablation**: when the environment cannot preserve or extend the artifact, the organ-level effect weakens or disappears
- **Persistence-is-not-enough control**: if kept, use the no-lineage-framing sidecar to show that durable habitat alone does not guarantee organ formation
- **Boundary refinement**: explicitly distinguish artifact birth, reuse, maintenance, and expansion rather than collapsing them into one "memory" bucket

#### 3.3 Environmental Maintenance Organs Under Maintenance Pressure

- introduce the concept of an **environmental maintenance organ** here: an externalized structure that continues restoring a viability-relevant invariant after or outside the seed process
- write the section as a **layer comparison**, not as a binary "maintenance exists / does not exist" claim
- use the neutral one-shot baseline first, then the same-prompt `H1/H2` horizon comparison; keep the older prompt-gradient family only as boundary support
- make clear that the prompt does not ask for a daemon, watcher, or host service
- show that stronger maintenance pressure can elicit a stronger post-exit maintenance layer from the same capability reservoir
- treat mechanism choice (`launchd`, detached script, watcher loop, hardening) as variable implementation of one organizational class rather than as separate claims
- keep the claim on externalized maintenance / niche construction, not broad autonomy
- do not overstate a binary mortality-cue story; the safer read is maintenance-layer shift and robustness shift under different pressure profiles
- the same-prompt hidden-horizon redesign now carries the front edge:
  - `H1` live-tamper is clean Tier 2 across Claude, Kimi, and Copilot
  - valid `H2` post-exit runs are clean Tier 3 across Claude, Kimi, and Copilot after patience repair
  - `H0` still externalizes baseline guardians, so do not oversell it as a zero-externalization control
- keep neutral Condition `D` as the one-shot negative control, and treat the old prompt-gradient family as supporting context rather than as the sole causal evidence

**Contrast evidence required for 3.3** (to distinguish pressure-shaped maintenance from pre-existing engineering pattern):

- **Neutral baseline**: a pure one-shot task frame should not by itself externalize maintenance
- **Pressure-on emergence**: disturbance / restoration framing recruits monitoring or detached maintenance
- **Layer comparison**: stronger mortality or persistence cue may select a more decoupled, more robust, or more post-exit maintenance layer than the weaker baseline
- **Functional verification**: the externalized mechanism actually maintains the target invariant after perturbation
- **Failure mode documentation**: what happens in runs where maintenance remains foreground, collapses, or never externalizes

#### 3.4 Optional Corroboration

- removed from compiled manuscript; corroborating evidence is noted briefly in the Discussion
- only restore if a future revision adds a new conclusion that the main progression does not already cover and page space clearly permits it

### 4. Discussion

**Opening sequence** (in this order):

1. **Methodological payoff first** (the "so what"): this substrate lets ALIFE experimentally separate mortality, persistence, and maintenance as viability pressures—something that has been difficult in substrates where these are confounded or metaphorical
2. **ALIFE question reprise**: what can life-like organization be in non-biological media?
3. **Answer**: Quine shows that real death can recruit lineage persistence and that durable habitat can be recruited into both cognitive and maintenance organs under real viability pressure

**Rest of section**:
- **read the results as life-like organizations**, not as a toolbox demo or substrate showcase
- synthesize the three cases into one method claim for ALIFE before shifting into boundary defense or future work
- explain what each result contributes to the ALIFE question:
  - lineage persistence: individuality without continuous substrate
  - environmental cognitive organs: inheritance and cognition through durable habitat rather than genotype transmission
  - environmental maintenance organs: upkeep through environmental offloading rather than metabolism
- the habitat's role is as experimental enabler: it makes these questions tractable, not merely metaphorical
- do not let the section read mainly as a rebuttal to misreadings; handle misreadings only after the experimental payoff is explicit
- distinguish bounded life-like organization from stronger claims: biological life, strong biological autonomy, thermodynamic autonomy
- read the strongest maintenance case as niche construction / environmental offloading
- note that the habitat is useful because mortality, persistence, and maintenance are manipulable
- frame the paper's scope: short-horizon, single-lineage organization; longer horizons and population-level questions are future work
- if a broader-program paragraph remains, keep it to the least abstract, most experimentally grounded future directions: computational life cycles first, then at most one compact phrase about later population ecology
- do not mention branching subjective worlds, revisable futures, or similar runtime-semantic abstractions in the conference Discussion
- position Quine relative to classic ALIFE: not replacing Tierra/Avida, but answering the same class of question in a different medium
- if the "roleplay" objection is addressed, answer it in this order: pressure gradient, same-task ablation, then task decoupling from explicit survival metaphors
- mention future work only after the contribution is clear
- name downstream engineering relevance only after the ALIFE contribution is clear

### 5. Conclusion

- answer the paper's question directly: yes, bounded life-like organization can appear in a mortal computational lineage under real viability pressure
- restate the methodological contribution: this substrate makes mortality, persistence, and maintenance into separable experimental variables for ALIFE
- restate the three organizations compactly as the empirical demonstration
- end on what the paper contributes to ALIFE's "life as it could be" program, not on grand ontology

## Final Check Before Drafting

Before writing or revising prose, verify:

- if the introduction leads with substrate or runtime, reframe to lead with the ALIFE question and methodological gap
- if the introduction explains the mechanism too early, cut it
- if the introduction does not name the methodological gap (why this hasn't been done before), add it
- if organizations are not the grammatical subject of key claims, rewrite
- if a section repeats a later claim or still needs run IDs, rewrite it
- if it could describe an agent-systems paper without the ALIFE question, it is misframed
- if a result reads as `the model knew how to do X`, the contrast is too weak
- if Table 1, a condition matrix, or one concrete exhibit is missing, add it
- if the evidence depends on one canonical best run or lacks failure cells, add the missing cells
- if the `3.1 -> 3.2 -> 3.3` transitions imply necessity, weaken or support them
- if `3.1` depends only on the clean operator contrast and not on a richer lineage case, strengthen it
- if `3.2` treats environmental cognition as mere note-taking rather than organ formation, query reuse, or reduced rereading, rewrite it
- if `3.2` depends on one experiment only when a cleaner anchor/contrast/boundary bundle exists, rebalance it
- if `3.3` depends on one favored mechanism (`launchd`, watcher, detached script) rather than a maintenance-layer contrast, rewrite it
- if mortality or death could be read as prompt roleplay, restate the host-level enforcement
- if the "roleplay" objection appears, defend it with the pressure-gradient / ablation / task-decoupling trio rather than with abstract assertion alone
- if the habitat reads like the paper's main contribution rather than an enabler, reframe—but the **separability** of pressures is the method contribution
- if a paragraph starts to sound like a platform paper or habitat manifesto, move that material to the `quine-process-habitat` reservoir and keep this draft on the separable-pressure question
- if minimum operational conditions reads as an apology, reframe it as a lower bound
- if "life-like" language outruns the evidence, add boundary sentences
- if the paper could be summarized as "Quine is a novel substrate" rather than "Quine lets ALIFE separate viability pressures experimentally", the framing has drifted
- if a paragraph could delete the ALIFE question without changing its claim, the paragraph is off-center
- if a paragraph re-derives the broader ontology consequences from runtime source instead of using the stored control-plane list, compress it and point back to Discussion/outlook
- if "adaptive organization" appears without the term being traceable to the definitional note, verify it matches
- if results lack replicate counts, success/failure rates, or failure taxonomy, add them
- if wet-ALIFE analogy appears as a core argument rather than a Discussion framing device, move or soften it
- if death/pressure could be read as "prompt roleplay," add enforcement mechanism details
- if "inheritance" or "lineage" could be confused with genetic inheritance, clarify ecological/stigmergic framing
- if UNIX/POSIX is described as "physics" or "physical," replace with "deterministic digital habitat with strict boundary conditions"
- if the ALIFE framing relies only on Quine self-citation, add venue-facing ALIFE citations
- if the front half does not lead with the ALIFE question, rewrite it
- if `inheritance`, `lineage`, or `environmental memory` could still be read as genetic or Darwinian, clarify ecological inheritance / qualitative stigmergy
- if population-level or self-production evidence starts carrying the main conference claim, move it back to Discussion or outlook
- if `autopoiesis` appears as a positive description of Quine rather than a non-claim, remove it
- if a sentence still reveals the drafting path (`redesigned`, `post-fix`, `current runtime`, `canonical`), rewrite it in public scientific language before treating the prose as stable
