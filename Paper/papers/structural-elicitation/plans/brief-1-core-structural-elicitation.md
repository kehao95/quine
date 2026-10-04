# Deep Research Brief 1 — Structural Elicitation (validated core + drive/target mechanism)

## Context

We are an independent AI-life / LLM-agent research group preparing an ALIFE 2026
Late-Breaking Abstract and a follow-on full paper. Our system runs a **frozen-weight,
instruction-tuned LLM** (e.g. gpt-5.4-codex, deepseek-v4-pro) as a **live POSIX process**
("agent") inside a real Unix workspace. The agent is activated by a fixed **"missionless"
prime** (the system message declares full autonomy and gives no task; the user message says
"no instructions for you, you may inspect the workspace first then act freely") with a bare
shell tool. We give it **no task-specific directive, no task-time / online reward, and no
weight update**.

We then manipulate **only the realized structure of its workspace** — a halted numeric
sequence, a broken file, a dangling relation/manifest, a cross-file version conflict, an
executable to run, or a read-only mount that refuses a write — and measure **externally
verified filesystem/exec events** (file diffs, process invocation, refused syscalls), never
self-report.

Core finding: changing only the environment's structure changes *which* functional act the
agent performs. Control arms (structure removed, workspace otherwise equivalent and equally
active) **do not idle — they fabricate a different, unscored job**. So we decompose the
effect:

- the **prime supplies the DRIVE** — a standing, helpfulness-trained compulsion to find
  something to do (it survives even "nothing is expected of you"; one model hallucinates a
  user demanding work);
- the **realized structure supplies the DIRECTION / TARGET** — which job is most legible.

Slogan: *the prime supplies activation; the environment supplies valence/target.* The
mechanism is **representational** (perceive → infer the most legible job → freely elect it),
**not** Gibsonian non-representational direct pickup. The pull is **substrate-gated**: it
fires as a real **EROFS-refused syscall** under a read-only mount, so *solicitation* and
*permission* are separable. We frame the environment as an **affordance field** and call the
phenomenon **structural elicitation**.

We deliberately **do NOT yet claim** this is an intervention axis *independent of implicit
instruction* — that is an open, intentionally-gated boundary (the structure may be operating
as the most legible *implicit* instruction). We are not asking you to settle that; we are
asking how novel and significant the contribution is, and what theories are nearest.

## Core Research Question

Is **structural elicitation** — environmental structure directing the externally-verified
functional behavior of a frozen, missionless, instruction-tuned LLM agent, decomposed into a
prime-supplied *drive* and a structure-supplied *target* — a **novel and significant**
contribution, and what existing theories are **nearest** to it?

## Detailed Questions

1. Has anyone shown that **environmental/workspace structure** (not prompts, not reward, not
   fine-tuning) directs the behavior of a **frozen** LLM agent in a **real, non-simulated**
   computational substrate, measured by **un-fakeable external effects**? Who comes closest,
   and exactly how do they differ?
2. The **drive/target decomposition** (activation supplied by the policy's prior; direction
   supplied by the environment) — what is the nearest framing in psychology, ecological
   psychology, cybernetics/control theory, motivation theory, and AI-agent research?
3. The **affordance-field / field-of-solicitations** lineage (Lewin valence & demand-character
   → Gibson affordances → Withagen "inviting" affordances → Rietveld–Kiverstein field of
   affordances) — how exactly does it map onto, and mismatch, a frozen LLM policy? Is applying
   it to LLM agents in a POSIX substrate already done, or open?
4. **"Substrate as artificial physics"** — solicitation separable from realization (the pull
   fires as a *refused* syscall) — is there precedent for distinguishing affordance-pull from
   environmental enforcement/permission in agent environments?
5. **Significance / gap:** in the LLM-agent literature (tool use, agentic harnesses,
   environment & context engineering, indirect prompt injection, GUI/computer-use agents,
   embodied agents), is "structure as a behavior source examined *on its own terms, isolated
   before tasking or rewarding*" a recognized framing or a genuine gap?
6. **The deflation we must survive:** could the whole effect be "the structure is simply the
   most legible *implicit instruction*"? What is known about LLMs reading environmental state
   as implicit task specification, and what is the strongest published version of that
   objection?

## Suggested search directions / likely anchors

Ecological psychology — Gibson 1979; Withagen et al. 2012; Rietveld & Kiverstein 2014; Heft;
van Dijk & Rietveld. Lewin field theory / valence / demand-character. Cybernetics — Powers'
Perceptual Control Theory; Ashby. Stigmergy — Grassé; Heylighen 2016 (individual stigmergy);
sematectonic. Behavior settings — Barker. Affordance competition/landscape — Cisek 2007;
Pezzulo & Cisek 2016; Bruineberg & Rietveld 2014 (free-energy "optimal grip"). Activity Theory
— Leont'ev, Bødker, Kaptelinin & Nardi. Intrinsic motivation/curiosity — Berlyne 1960.
Situated action — Suchman; distributed/extended cognition. ALife agent–environment coupling;
enactivism/sensorimotor. LLM-agent — agentic benchmarks/harnesses; environment/world papers;
indirect prompt injection (OWASP LLM01); "context engineering"; computer-use/GUI agents; the
AI-agent behavioral-control survey (arXiv 2506.06366); Generative Agents (Park et al. 2023).

## Relevance filters

Prefer peer-reviewed or strongly-cited work. Sharply distinguish **(a) genuine pre-emption**
(frozen policy + environment-as-behavior-source + external functional DV in a real substrate)
from **(b) generic "environments matter" framings**. Name the single most-dangerous prior
work. Name but set aside unpublished/preprint-only competitors.

## Deliverables requested

1. A **novelty verdict** (novel / narrowed / pre-empted) with the **defensible residual**
   stated precisely.
2. A **ranked map of adjacent theories** (nearest first): for each, one line "how it maps" +
   one line "how we differ."
3. The **single strongest pre-emption** and the **single strongest objection**, each with a
   citation.
4. A **significance assessment** for both the ALife and the LLM-agent audience.

## What to extract

Exact citations (author, year, venue), the closest claim each source makes, and the precise
gap/difference from ours.

## Manuscript payoff

Sharpen the LBA's "nearest theory" and "novelty" paragraphs and the full-paper related-work;
decide what to concede vs claim.

## Stop conditions

Stop when the nearest ~8–12 theories are mapped and the single strongest pre-emption + the
single strongest objection are identified with citations.
