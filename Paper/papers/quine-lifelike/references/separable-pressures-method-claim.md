# Separable Mortality, Persistence, and Maintenance Pressures

This note absorbs the current background investigation for the paper's methodological claim.

It is not manuscript prose. It is a reference-side synthesis used to decide what the paper can safely claim, what comparisons matter, and which citations should carry those comparisons.

## Current Read

The manuscript's strongest new center is not only that Quine exhibits bounded continuity, ecological inheritance, and detached maintenance. It is that Quine appears to make **mortality, persistence, and maintenance separable experimental pressures** inside a soft-ALIFE medium.

That methodological move is potentially venue-native for ALIFE because the field has long treated new substrates as scientific instruments for studying `life-as-it-could-be`, not merely as engineering demos. `Langton1989ArtificialLife`, `Bedau2003ArtificialLifeReview`, and `Beer2024ALifeAsItCouldBe` remain the primary high-level anchors for that point.

## Why Separable Pressures Matter

The relevant methodological question is:

- which causal questions about life-like organization become experimentally tractable in a given substrate?

For this paper, the answer is:

- Quine potentially allows mortality, persistence, and maintenance to be manipulated independently enough that we can ask which organizations are recruited by which pressures, instead of treating those pressures as one inseparable package.

This is stronger than saying the substrate is merely flexible or configurable. The claim is about **orthogonalized intervention surfaces**:

- mortality can be changed without automatically rewriting persistence or maintenance
- persistence can be changed without automatically rewriting mortality policy
- maintenance requirements can be changed without simply renaming death or persistence failure

## Comparator Read

### 1. Classic digital ALIFE / digital evolution

Use `Ray1992SynthesisOfLife` and `OfriaWilke2004Avida` as the base comparators, but the relevant contrast is not "they are older digital organisms." The relevant contrast is how pressure is wired into the substrate.

Current comparative read:

- `Tierra` has real scarcity and real death, but mortality is strongly tied to soup-space pressure, reaper policy, and overwrite dynamics.
- `Avida` has real configurable death policies and rich population management, but mortality, persistence, and damage are still largely embedded in the organism execution model rather than exposed as three explicitly independent pressure axes.

Implication for the paper:

- do not frame Quine as "more real" than Tierra or Avida
- frame Quine as offering a different kind of experimental control surface
- if the manuscript keeps a strong method claim, it should explicitly say that classic digital ecologies have real selection but more tightly coupled viability physics

### 2. Soft / neural / CA-style life-like systems

The useful comparator class is not "these are weak ALIFE." It is "these systems often show robust life-like behavior without turning mortality/persistence/maintenance into independently enforceable levers."

Current read from the literature:

- `MordvintsevEtAl2020GrowingNCA` is strong evidence that growth, persistence, and regeneration can be engineered or trained in decentralized media, but pressure is largely imposed through training setup and target-state loss rather than through ongoing mortality-aware ecology.
- `Chan2020LeniaExpandedUniverse` is strong evidence that CA-like media can show rich self-organizing and homeostatic behavior, but viability is usually dynamical stability under the rule rather than independently instrumented mortality / persistence / maintenance accounting.

Implication for the paper:

- these systems are useful for the `metaphorical / fixed / externally posed` side of the gap claim
- do not overstate them as merely rhetorical; they exhibit real dynamical robustness
- instead say they usually do not provide the same kind of substrate-level orthogonalization of viability pressures

### 3. Explicit maintenance / viability models

This is the most important near-miss category for reviewer defense because it contains systems that already treat maintenance as a real causal variable.

Current read:

- `McMullin2004ThirtyYearsComputationalAutopoiesis` and `McMullin1997ComputationalAutopoiesisOriginalAlgorithm` show that computational autopoiesis treats decay, repair, and viability as real model dynamics rather than rhetorical overlays.
- however, mortality and maintenance are usually tightly coupled there: failure of maintenance is effectively the death mechanism.

Implication for the paper:

- autopoietic chemistry is the strongest prior art against any broad claim that Quine is the first soft-ALIFE system with real maintenance pressure
- the safer distinction is narrower: Quine may offer a substrate where mortality, persistence, and maintenance are jointly real **and** experimentally separable in a way that autopoietic models, classic digital organisms, and neural morphogenesis systems usually do not

### 4. Framework-level controllability

`BohmEtAl2019MABE2` matters because it is a precedent for publishing experimental controllability itself as a methodological contribution.

Current read:

- MABE2 is a framework for modular ALIFE / evolutionary experiments
- it demonstrates that configurability and reduced time from hypothesis to result are legitimate contribution types
- but it is a framework for hosting designs, not itself a specific substrate proving orthogonalized mortality / persistence / maintenance

Implication for the paper:

- cite MABE2 only to support `tractable experimentation as a contribution type`
- do not use it as a direct substrate comparator unless you are careful about framework versus medium

## Terminology Support

### Ecological inheritance

`OdlingSmeeEtAl2003NicheConstruction` remains the broad anchor. Add `OdlingSmeeLaland2011EcologicalInheritance` when the prose needs to distinguish:

- ecological inheritance from genetic inheritance
- ecological inheritance from generic note-taking
- inheritance through external environmental legacy effects

Working interpretation for this paper:

- predecessor-modified environmental structure that persists and shapes successor behavior is enough for a bounded ecological-inheritance reading
- the manuscript should still avoid implying Darwinian heredity, population genetics, or replicator-based inheritance

### Maintenance and organizational closure

`MaturanaVarela1980Autopoiesis` is still best used as a boundary concept, not as the paper's positive explanatory center.

If discussion needs a tighter bridge, `MossioMoreno2010OrganisationalClosure` is the cleaner citation for:

- self-maintaining organization
- closure as a stronger biological standard than the paper claims
- why detached maintenance in Quine is interesting without equating it to full biological autonomy

## What the Main Claim Can Safely Be

### Safe core claim

Quine provides a controllable soft-ALIFE substrate in which mortality, persistence, and maintenance can be manipulated independently enough to ask which lineage-level organizations are recruited by which pressures.

This is then demonstrated with three contrasts:

- continuity under mortality pressure
- ecological inheritance under persistence pressure
- detached maintenance under maintenance pressure

### Risky claim

Avoid broad versions of:

- Quine is the first ALIFE substrate with real mortality
- Quine is the first substrate with real maintenance
- Quine is the first soft-ALIFE system with enforceable viability at all

The literature on Tierra, Avida, computational autopoiesis, and neural / CA regeneration gives reviewers too many plausible objections.

### Narrower strong claim

If the paper wants a stronger novelty sentence, the most defensible version is something like:

- to our knowledge, Quine is a soft-ALIFE medium that combines host-enforced mortality, manipulable environmental persistence, and detachable maintenance channels in one substrate, allowing these pressures to be separated experimentally rather than treated as fixed or tightly coupled world dynamics

This remains strong, but it avoids a naive `first real ALIFE substrate` reading.

## Recommended Citation Jobs

Use the new sources only where they do real argumentative work:

- `Section 1`:
  - `Langton1989ArtificialLife`
  - `Beer2024ALifeAsItCouldBe`
  - `Ray1992SynthesisOfLife`
  - `OfriaWilke2004Avida`
  - optionally `McMullin2004ThirtyYearsComputationalAutopoiesis` if the gap sentence needs a maintenance-aware comparator

- `Section 2`:
  - `Ray1992SynthesisOfLife`
  - `OfriaWilke2004Avida`
  - `McMullin2004ThirtyYearsComputationalAutopoiesis`
  - `MordvintsevEtAl2020GrowingNCA`
  - `Chan2020LeniaExpandedUniverse`
  - `BohmEtAl2019MABE2` only for methodological controllability as a contribution type

- `Inheritance terminology`:
  - `OdlingSmeeEtAl2003NicheConstruction`
  - `OdlingSmeeLaland2011EcologicalInheritance`
  - `TheraulazBonabeau1999Stigmergy`

- `Discussion boundary work`:
  - `MossioMoreno2010OrganisationalClosure`
  - `MaturanaVarela1980Autopoiesis`

## Red Lines

- do not write the method claim as if all prior ALIFE substrates lacked real pressure
- do not write the soft / neural comparator class as if trained regeneration were merely rhetorical
- do not write the autopoiesis comparator class as if it lacked enforceable maintenance
- do not keep an unrestricted `first` claim unless the manuscript carries a compact but explicit comparator table
