# ALIFE Reference Deep-Research Plan

This file is a focused investigation plan for strengthening the front half and related-work posture of the ALIFE conference paper.

The goal is not to collect every relevant citation. The goal is to answer the few questions that materially affect admission, framing, and overclaim discipline.

## Expected Outputs

Each investigation thread should produce:

- `1` short note answering the question in plain language
- `3-6` candidate citations with one-line relevance notes
- `1` draft paragraph or terminology note that could be dropped into the manuscript
- `1` explicit non-claim or red-line to avoid

## Priority 1: Admission And Framing

### 1. What lets Quine count as an ALIFE system at all?

Answer this before polishing prose.

- target question: which ALIFE concepts best justify Quine as a `non-traditional computational medium` without overclaiming `life`
- likely anchors: Langton on `life as it could be`, Beer on the current field and future direction, Bedau on ALIFE as synthetic study of organization/adaptation/complexity
- what to extract:
  - exact language about synthetic media, constructive methodology, and bounded ALIFE contribution
  - whether the field already accepts toy or non-biochemical media when the dynamics are scientifically legible
  - what level of claim is normal for conference papers versus ontology papers
- manuscript payoff:
  - `Introduction`: one venue-facing paragraph
  - `2.1 Why ALIFE?`: two or three anchor citations that do real argumentative work
- red line:
  - do not rely on Quine self-citation alone for ALIFE admission

### 2. What are the strongest prior-art comparators for digital organisms and digital evolution?

- target question: relative to Tierra, Avida, and adjacent digital-evolution systems, what is genuinely new here and what is not
- likely anchors: Ray on Tierra, Ofria/Wilke on Avida, one later Avida or digital-evolution result that shows ecological dynamics under explicit conditions
- what to extract:
  - how prior systems define organisms, inheritance, ecology, and observability
  - which of those systems already had strong condition control and digital-organism language
  - what Quine adds: real host habitat, semantic task layer, durable external state, mortal process lineages
- manuscript payoff:
  - related-work paragraph that prevents `this is just a new digital organism claim`
  - discussion sentence clarifying that Quine is not the first observable ALIFE system
- red line:
  - do not imply Quine invented digital organisms, experimental digital ecology, or observability

### 3. Which conceptual frame best supports `organism-environment coupling` for this paper?

- target question: what framework best supports the claim that the habitat is causal medium rather than backdrop
- likely anchors: Beer on agent-environment dynamics, plus one autonomy/agency source if needed
- useful follow-ons if the first pass is too thin:
  - Barandiaran et al. on agency and individuality
  - Moreno and Mossio on biological autonomy
- what to extract:
  - language for coupled dynamics, constraints, and viability
  - whether the organism should be described as the process alone or the process lineage plus habitat coupling
- manuscript payoff:
  - sharper thesis sentence
  - tighter definitions for organism, generation, death, inheritance
- red line:
  - do not smuggle in stronger autonomy claims than the experiments support

## Priority 2: Terminology Defense

### 4. How should the paper defend `inheritance`, `lineage`, and `environmental memory`?

- target question: which ecological or evolutionary concepts make those terms legitimate without sounding like genetic inheritance
- likely anchors: niche construction, ecological inheritance, stigmergy
- likely sources: Odling-Smee/Laland/Feldman, Theraulaz/Bonabeau, plus one narrower ecological-inheritance source if needed
- what to extract:
  - whether durable environmental modification that shapes later behavior is enough for ecological inheritance language
  - whether stigmergy is the better term for `successor reads predecessor artifacts`
  - how to distinguish inherited environmental support from mere logging or task notes
- manuscript payoff:
  - `2.3` terminology note
  - discussion framing for the strongest maintenance case as niche construction
- red line:
  - do not let `inheritance` read as genotype transfer or population evolution

### 5. How much autopoiesis language can the conference paper safely carry?

- target question: what is the minimum citation and concept work needed to mention autopoiesis only as a boundary, not as the main claim
- likely anchors: Maturana and Varela, maybe one later autonomy source if needed
- what to extract:
  - the minimal requirements for self-production, boundary maintenance, and organizational closure
  - why Quine's current conference evidence falls short of a strong autopoiesis claim
- manuscript payoff:
  - one disciplined discussion paragraph
  - stronger non-claim language
- red line:
  - do not let `post-exit maintenance` silently become `autopoiesis demonstrated`

## Priority 3: Contemporary Positioning

### 6. What ALIFE-plus-LLM work already exists, and how should Quine differ from it?

- target question: where does this paper sit relative to recent ALIFE/LLM bridge papers
- likely anchors: Nisioti et al. on ALIFE and LLMs, Kumar et al. on ASAL
- useful comparators:
  - Generative Agents
  - Affordable Generative Agents
  - other LLM-agent environment papers only if they materially affect your novelty boundary
- what to extract:
  - whether those papers treat LLMs as organisms, operators, search tools, or environment participants
  - whether they manipulate viability pressures or instead study capability and behavior in fixed environments
  - what niche remains for Quine: host-enforced mortality, persistence, and maintenance as manipulable pressures in a real computational habitat
- manuscript payoff:
  - one related-work paragraph in Section 2 or Discussion
  - one sentence clarifying why this paper is about selected phenotype, not just LLM-enabled emergence
- red line:
  - do not overstate novelty where the field already has LLM-plus-ALIFE bridge work

### 7. Should the paper explicitly defend Quine as a `toy world` or `toy model`?

- target question: is the toy-model framing a venue asset or an unnecessary complication
- likely anchors: Beer and collaborators on toy models, plus any ALIFE discussions that value idealized synthetic systems
- what to extract:
  - what toy models are good for
  - how toy-model legitimacy differs from realism claims
  - whether this framing helps with reviewer skepticism about UNIX/POSIX as a habitat
- manuscript payoff:
  - optional introduction or discussion sentence
  - safer answer to `this is not biology`
- red line:
  - do not widen the paper into philosophy-of-science unless the reviewer-risk clearly requires it

## Suggested Investigation Order

1. `Langton -> Beer 2024 -> Bedau 2003`
2. `Ray -> Ofria/Wilke -> one ecological-digital-evolution comparator`
3. `Beer 1995 -> autonomy/agency follow-ons only if needed`
4. `Odling-Smee -> Theraulaz/Bonabeau -> ecological inheritance follow-up`
5. `Maturana/Varela` for boundary discipline
6. `Nisioti -> ASAL -> agent comparators`
7. toy-model framing only if the first draft still feels under-justified

## Stop Conditions

Stop digging when:

- the paper has `2-3` strong ALIFE anchors in the introduction and `2-3` more in Section `2`
- each risky term has at least one defensible citation path
- the related-work contrast can be stated in one paragraph without hand-waving
- new citations stop changing wording or claim boundaries

If new reading starts expanding the paper's ontology instead of tightening the conference argument, stop and return to drafting.
