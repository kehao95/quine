# Manuscript Plan

## Status

This file is now the active writing plan for the essay, restructured to center on a suspense-and-conflict-driven narrative infused with personal meta-reflections.

The full earlier brainstorm, including raw chronology and uncompressed pain-point capture, has been preserved without deletion in:
- [`./analysis/ideation-capture.md`](analysis/ideation-capture.md)

Use this file for section flow, emphasis, and drafting order.

## Working Shape

Recommended article shape:
- A first-person engineering case study infused with strong philosophical meta-reflections.
- Organized around cognitive dissonance (the gap between "open source" norms and LLM research realities) and engineering deadlocks.
- Deeply personal in its aesthetic choices (the obsession with historical authenticity).
- Reusable for other research repositories without pretending Quine is generic.

Recommended non-goals:
- Do not turn this into a generic DVC tutorial.
- Do not turn this into a long chronology of cleanup work.
- Do not flatten the story into "I removed secrets and made the repo public".

## Working Thesis

In the LLM era, open-sourcing pure runtime code has diminishing value; what truly matters is sharing the research, the failures, and the raw evolutionary history. However, open-sourcing a living research workbench without destroying its historical authenticity or leaking sensitive host data requires fundamentally separating the repository's structural skeleton from its heavy payloads using DVC as a boundary control plane.

## Narrative Center Of Gravity

The strongest version of this essay centers on four pillars:
1. **The Code-Devaluation Reality**: Releasing an "AI wrapper" is trivial; releasing the messy, living research is what actually matters.
2. **The Security Deadlock**: The stochastic nature of LLMs means logs are irreplaceable memory, yet these same logs are the biggest source of host environment leaks.
3. **The Authenticity Obsession**: A stubborn refusal to take the easy way out (squashing/deleting history). The historical footprint *is* the work.
4. **The Boundary Innovation**: Using DVC not as a large-file garbage can, but as a precise publication-boundary protocol to resolve the deadlock.

## Recommended Organization

### 1. Opening & Meta-Premise 1: The Diminishing Value of "Clean Code"

**Job of the section**
- Hook the reader with a critique of modern AI open-source culture (just scrubbing API keys and pushing).
- Introduce the first meta-reflection: In the LLM era, code is cheap. The ideas, the experiments, and the raw traces are the real assets.
- Establish the goal: We didn't want to open-source a runtime; we wanted to open-source the `Living Research Workbench`.
- The twist: Doing this honestly is exponentially harder than a standard `git push`.

**Key points to cover**
- Quine combined runtime code, theory, experiments, logs, and drafts.
- Why a sanitized, isolated `public` branch feels like a museum exhibit, not a working lab.

### 2. The Conflict: Logs as Memory vs. Logs as Leakage

**Job of the section**
- Establish the core engineering deadlock of the repository.
- Make the reader understand why the standard advice ("just delete the logs") fails here.

**Key points to cover**
- **The Research Need**: LLM experiments have real stochasticity. You cannot deterministically rerun them. Logs are not residue; they are the primary read surface for future theory-building.
- **The Security Nightmare**: Because the project was private-first, the agent freely interacted with the host environment. Shell outputs, trace files, and absolute paths became contaminated with sensitive host data.
- **The Trap**: You can't delete them (destroys the science), and you can't publish them (leaks the secrets).

### 3. Meta-Premise 2: The Stubborn Obsession with Historical Truth

**Job of the section**
- Inject strong personal engineering aesthetics and emotional stakes.
- Reject the easy "industry standard" compromise.

**Key points to cover**
- The pragmatic engineering solution would be to squash the history, delete the 1000+ messy commits, and start a fresh public repo.
- The meta-reflection: A deep, almost obsessive "cleanliness" regarding historical truth. Deleting the history feels like faking the research process.
- The constraint is set: We must scrub the poison (secrets) from the repository, but we absolutely refuse to alter the 1000+ footprints (commits) that got us here.

### 4. The Pivot: DVC as a Boundary, Not a Storage Bin

**Job of the section**
- Present the main engineering breakthrough.
- Redefine how the reader thinks about DVC.

**Key points to cover**
- Since we can't delete history and can't publish payloads, we must decouple the *structure* from the *payload*.
- DVC was brought in not just to offload large files, but to act as a semantic boundary.
- Private Retention vs. Public Publication: The DVC server is self-hosted, turning DVC into a "manifest-scoped allowlist." We keep everything privately, but only "publish" what is audited.

### 5. The Surgery: Rewriting 1000+ Commits Without Faking Them

**Job of the section**
- Show, practically, how the obsession with authenticity was achieved.
- Focus on the precision of the history rewrite.

**Key points to cover**
- HEAD cleanup vs. History Rewrite.
- Using `git-filter-repo` and DVC backfilling to perform a surgical strike on the repository's past.
- We kept the commit metadata, the timestamps, and the directory structures, but swapped the radioactive log files with safe DVC pointers.

### 6. The Reveal: What Actually Became Public

**Job of the section**
- Close the story with concrete results and generalize the lessons.

**Key points to cover**
- The numbers: Pack size dropped from 469 MiB to 26 MiB, yet the commit history remains 100% intact.
- Define the new architectural reality of the public Quine:
  - Public Source/History (The authentic skeleton)
  - Private Retention (The unabridged memory)
  - Published Public Evidence (The audited, restorable data)
- Conclusion: This tri-boundary architecture is what it actually takes to practice Open Science in the era of LLM agents.

## Recommended Emphasis

Approximate emphasis split for a strong draft:
- `15%` Meta-Premise 1 (Code devaluation & Open Research vision)
- `25%` The Deadlock (Logs as memory vs. Logs as leak vectors)
- `15%` Meta-Premise 2 (The obsession with historical authenticity)
- `25%` The Pivot & Surgery (DVC boundaries and history replay)
- `20%` The Reveal (Metrics, the final public object, and takeaways)

## Quantitative Material

Use numbers to anchor the "Surgery" and "Reveal" sections. (Refer to [`./analysis/size-and-blob-audit.md`](analysis/size-and-blob-audit.md))
- Fresh clone comparison: `469.28 MiB` (lab) vs `26.74 MiB` (main).
- `lab -> main` reachable blob bytes reduction: `-87.5%`.
- These numbers prove that the "payloads" were successfully extracted while preserving the structure.

## Drafting Sequence

Recommended writing order:
1. Draft **Section 2 & 3** first: The core conflict and the personal obsession are the heart of the essay. Get the tone right here.
2. Draft **Section 1**: Now that you know the emotional core, write the hook. Frame the AI open-source landscape.
3. Draft **Section 4 & 5**: Explain the DVC boundary and the history rewrite mechanics.
4. Draft **Section 6**: Bring it home with the metrics and the overarching architectural lessons.
