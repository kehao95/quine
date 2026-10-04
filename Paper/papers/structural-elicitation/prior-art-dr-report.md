# Deep Research: Prior Art & Novelty of "Structural Elicitation" (2026-06-15)

> Raw `dr` intake (100-agent fan-out; 18 sources fetched, 84 claims extracted,
> 25 verified, 21 confirmed, 4 refuted, 10 after synthesis). Research intake, not
> citation-ready authority. Source run: session task `wdb7gg7xp`. Durable
> consequence promoted to the locked claim note
> [`./claim.md`](claim.md)
> § Novelty (confirmed 2026-06-15).

## Claim Positioned

"Structural Elicitation": for a frozen-weight LLM agent, environmental structure
is a candidate third source of behavior beside instruction (goal) and reward
(optimization signal); structure exerts a pull/solicitation the fixed policy
resolves into action. The strong theory version asks whether this can be
isolated without instruction or reward. Anchor mode: the structural gap pulls
completion. Substrate: Quine = LLM agents as real POSIX processes.

**Current scope note (post-p13 pilot):** the LBA no longer tries to prove the
strong instruction-free version. It reports a fixed-prime breadth phenomenon:
under a missionless-autonomy runtime prime, structural motifs direct
externally-verified behavior with no task-specific directive and no task-time
reward. The prior-art objections below still own the full-paper independence
bar.

## Verdict

- **Viability:** can hold, but conditional — two objections survived adversarial
  verification and the load-bearing experiment is unrun in all surveyed prior
  art.
- **Novelty:** the *framework* is owned (not novel); the defensible residual is
  the *substrate transfer + third-source taxonomy*, upgraded by an isolated
  demonstration.

## What is owned (NOT novel) — concede

- "Structure solicits action" is established ecological psychology: affordances
  have "demand character", "soliciting power", an "inviting character" that
  "draws the individual agent into action" (Kiverstein, van Dijk & Rietveld 2019,
  *Synthese*, 10.1007/s11229-019-02123-x).
- The Lewin (demand-character / valence) -> Gibson -> Rietveld-Kiverstein
  field-of-affordances genealogy is *already integrated* into one framework (the
  "novel synthesis" claim was refuted 0-3). Do not claim the framework or the
  theory-combination as novel.

## The defensible residual (novel) — claim this

1. **Substrate transfer:** the affordance/solicitation framework (the field-of-
   solicitations sense, not the robotics manipulation sense) applied to
   frozen-weight LLM agents in a real, non-simulated POSIX substrate. The
   ecological-psychology literature never addresses computation/AI; the 2024-2026
   "affordance + LLM" landscape (3D-AffordanceLLM, Affordance-R1, A4-Agent, ...)
   uses only the manipulation sense; none uses the solicitation sense or a POSIX
   substrate.
2. **Third-source taxonomy:** structure counterposed as a third source beside
   instruction and reward. The dominant AI-agent behavioral survey (Chen et al.,
   *AI Agent Behavioral Science*, arXiv 2506.06366 / Nature HSSC) uses a different
   triad (intrinsic attributes / environmental constraints / behavioral feedback)
   and never counterposes structure to prompt vs reward (grep of the full PDF:
   affordance=0, stigmergy=0, Lewin=0, Gibson=0, behavior-setting=0).
3. **(Gated) the empirical payload:** the first isolated demonstration that
   structure elicits functional behavior in a frozen LLM with no instruction and
   no reward (the gap anchor; `experiments/active/structural-elicitation`).

## Viability — the two objections (both survived)

- **(a) "gap = implicit instruction" (from `cat` not `stdin`):** unrefuted. A
  perceived gap enters context and is plausibly an implicit directive; may be
  irrefutable in principle for a frozen LLM. Decide between empirical defeat
  (directive-additivity) and concede-and-reframe.
- **(b) "completion is the base NTP behavior, not a third source":** has teeth —
  next-token-prediction induces a structural completion/reasoning bias (Lin,
  Zhang & Xu, arXiv 2502.02007; 2-1/3-0; caveat: small synthetic depth task).
- **Load-bearing test (unrun in all prior art):** a design that shows the gap
  drives functional completion with no directive AND that structure contributes
  variance an explicit directive does not. The p13 `gap` vs `complete` one-factor
  run establishes gap-necessity but does not by itself isolate structure from
  implicit instruction (objection a); a Phase-2 discriminator (a non-default
  structure for (b), or a directive-additivity arm for (a)) is needed to move
  from "viable" to "novel".

## Nearest neighbors

- **TerraLingua** (Cognizant, arXiv 2603.16910, Mar 2026): frozen-weight agents,
  environment-as-locus, Gibsonian affordance gating + stigmergy; overlaps 3 of 4
  pillars. **Set aside as competition (PI call, 2026-06-15):** unpublished / not
  peer-reviewed, and its thesis (open-endedness / cumulative culture, with agent
  goals + a hard-coded variation+selection loop) is a different object than
  structure-as-third-source in a no-reward regime — the `dr` independently found
  it does NOT occupy the third-source / no-goal-no-reward territory. Acknowledge
  and position lightly; do not treat as pre-empting prior art.
- **Generative Agents** (Park et al., arXiv 2304.03442): behavior from an
  internal cognitive loop + goal seed, simulated Sims sandbox — contrast on
  source-locus AND substrate.
- **Vallinder & Hughes** (arXiv 2412.10270): frozen-weight but explicit payoff +
  top-50% selection — reward/instruction-driven, contrast on the reward axis.

## Refuted claims — do NOT lean on

- Fogg-model "leaves room for a third source" (refuted 1-2).
- Lewin -> Gibson -> Rietveld "novel synthesis" (refuted 0-3).
- TerraLingua / survey "environment is causally consequential / non-programmed
  emergence" as direct precedent (refuted 1-2 / 0-3).

## Open Questions

- Can objection (a) be empirically defeated (directive-additivity) or only
  conceded-and-reframed?
- Can structural pull be dissociated from the amortized NTP completion prior
  (objection b) via structures the base predictor would not default to?

Sources: *Synthese* 10.1007/s11229-019-02123-x; arXiv 2603.16910; arXiv
2506.06366; arXiv 2304.03442; arXiv 2412.10270; arXiv 2502.02007.
