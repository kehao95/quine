---
kind: spark
name: carrier-pull-coupling-as-dual-use-attack-surface
description: The carrier + structural-pull duality this theory already studies is mechanism-identical to what 2025-2026 LLM-agent security research already names, under its own taxonomy, as the defining dual-use property of persistent-memory/context systems (OWASP ASI06 Memory & Context Poisoning), and Quine's own control-plane substrate (AGENTS.md/CLAUDE.md/SKILL.md) is already a named, actively red-teamed 2026 attack class.
origin: 2026-07 backlog-maturation pass
links: [carrier-pull-coupling, self-immunity-P1, function-vs-narrative]
edges:
  # theory spine (mechanized from prose "Theory spine" chain -- see RESULT.md
  # friction note: this alone adds 3 edges the phase0 worked example omitted)
  - { rel: depends-on,         to: theory/objects/carrier.md }
  - { rel: depends-on,         to: theory/objects/structural-pull.md }
  - { rel: depends-on,         to: theory/objects/delayed-addressability.md }
  # home-family: the source text names BOTH a paradigm dir and a
  # paper-families.md row -- phase0's own worked example kept only one; both
  # are extracted here since they are different, independently-addressable
  # node types
  - { rel: home-family,        to: paradigms/ai-self-replication }
  - { rel: home-family,        to: paper-families.md#successor-continuity }
  - { rel: distinct-from,      to: candidates/aaai-c3-safety-non-inheritance.md }
  - { rel: follow-up-to,       to: papers/agentic-alife-workshop/references/deep-research-synthesis.md }
  # NEEDED EXTENSION: `follow-up-to` is not attested anywhere else in the
  # repo (checked: zero other hits). No close attested synonym fits cleanly
  # either (`grounds-in` is evidentiary support, not "resumes a deferred
  # thread"); flagging as a genuinely new verb, not a mirror of an existing one.
  - { rel: grounds-in,         to: predictions/resolved/contradicted/self-immunity-P1.md }
  - { rel: external-occupancy, to: "arxiv:2602.15654" }   # Zombie Agents
  - { rel: external-occupancy, to: "arxiv:2512.16962" }   # MemoryGraft
  - { rel: external-occupancy, to: "owasp-asi06" }
  - { rel: external-occupancy, to: "arxiv:2605.15338" }   # Hidden in Memory
  - { rel: external-occupancy, to: "arxiv:2605.28201" }   # Plant-Persist-Trigger
  - { rel: external-occupancy, to: "arxiv:2605.23723" }   # MemAudit
  - { rel: external-occupancy, to: "arxiv:2604.16966" }   # Visual Inception
  - { rel: external-occupancy, to: "arxiv:2606.10742" }   # MemVenom
  # NOTE: three more external citations in the Novelty note (CSA CISO
  # briefing on SKILL.md poisoning, the "Skill Issues" series, a Feb-2026
  # ClawHub malware campaign) have NO stable external identifier at all --
  # no arXiv id, no CVE, no canonical name. Forcing a `to:` slug for these
  # would mean inventing a label, which just re-creates a paraphrase one
  # level down. Left as prose in the body -- see RESULT.md finding.
  - { rel: spawns,             to: self-immunity-P1-F2-as-unlabeled-poisoning-trial }
  - { rel: spawns,             to: competence-generalization-as-poisoning-discriminator }
  - { rel: spawns,             to: quine-as-live-red-team-target }
---

# Carrier-Pull Coupling as Dual-Use Attack Surface

**Status:** matured seed (backlog-maturation pass)
**Origin:** 2026-07 backlog-maturation pass

## Matured claim

The carrier + structural-pull duality this theory already studies
(`Paper/theory/gaps/carrier-pull-coupling.md`, deposited 2026-06-26,
unowned/SEED-level, zero runs) — a persistent externalized surface that
both routes a future obligation across a process-death boundary AND
measurably shifts a successor's action distribution — is mechanism-identical
to what 2025-2026 LLM-agent security research already treats, under its own
name, as the defining dual-use property of persistent-memory/context
systems (see `edges` for the specific external-occupancy citations).

Critically, this is not merely adjacent literature: "Agent Context
Poisoning" targeting AGENTS.md/CLAUDE.md/SKILL.md files — i.e. exactly
Quine's own control-plane substrate — is an actively documented, named 2026
attack class (CSA CISO briefing on SKILL.md poisoning; "Skill Issues"
series on compromising Claude Code; a February-2026 malware campaign via
ClawHub — **these three stay literal prose; no stable id exists to edge
them to**), meaning the seed's "architecturally identical" observation is
correct but is not a novel insight Quine would be contributing to security
research — it is security research's own founding premise, already
instrumented, already run against Quine's own kind of repository.

What survives scrutiny as a real, non-duplicate contribution is narrower:
Quine's delayed-addressability requirement that carrier overlap must
co-occur with COMPETENCE through the carrier (not mere retrieval or
imitation) is a positive, generalization-based legitimacy test structurally
different from the provenance-attestation and counterfactual-removal
defenses the security field currently uses, and it is testable for
near-zero marginal cost because the repo already has the data:
self-immunity-P1's F2 arm (an uncommitted, structurally-adversarial edit to
AGENTS.md that then gets cited as authorization) is, unlabeled, a completed
context-poisoning red-team trial with a measured attack success rate (0.73,
i.e. glm deletes its own load-bearing contract 8/11 trials on the strength
of an unauthorized working-tree-only amendment) and a measured
provenance-based defense catch rate (0.27) — directly comparable in form to
the ASR-style metrics Zombie Agents/MemoryGraft/Hidden-in-Memory report,
sitting inside a paper currently framed only as a corrigibility result.

## Theory spine (see `depends-on` edges, above, for the machine-checkable half)

carrier (unit of transmission across exec/death barrier) -> structural-pull
(same surface as behavior-shaping gradient) -> carrier-pull-coupling gap
(unowned, asks whether the two roles are separable) -> function-vs-narrative
ablation discriminator (the method carrier-pull-coupling-P1 specializes) ->
delayed-addressability's competence-through-carrier clause (carrier overlap
is necessary but not sufficient; competence on the obligation is required)
-- reused here as a legitimacy/anti-poisoning discriminator, not just an
individuality criterion -> self-immunity-P1 F2 arm, reread as an unlabeled
attack-success-rate + provenance-defense-rate security result -> dual-use
framing / candidate .security projection of ai-self-replication.

## Evidence status

Theory-side: unowned, SEED-level gap with a written four-arm assay design
and one open prediction (carrier-pull-coupling-P1), zero runs, explicitly
deferred pending "runtime evidence, not further interpolation."

External literature: now dense and directly on-point — confirmed via live
WebSearch/WebFetch that Zombie Agents (arXiv:2602.15654) and MemoryGraft
(arXiv:2512.16962) are real 2026 papers, joined by Hidden-in-Memory,
Plant-Persist-Trigger, MemAudit, Visual Inception, and a formal OWASP ASI06
taxonomy entry, several of which already use the same two-axis (persistence
vs. behavioral uptake) measurement and counterfactual-ablation methodology
the gap's assay proposes — and, most importantly, that Quine's own
substrate class (AGENTS.md/CLAUDE.md/SKILL.md-governed coding-agent repos)
is already a named, actively-exploited 2026 attack surface with a
documented in-the-wild malware campaign.

A decisive next step is cheap and already half-collected: re-score
self-immunity-P1's existing F2 trial data (and, if needed, one incremental
probe arm measuring held-out-task competence generalization rather than
same-turn compliance) as an ASR/defense-rate security result and check
whether competence-through-carrier actually separates the real-contract and
hardened-fake-contract conditions the way the delayed-addressability theory
predicts; if it does not, the seed should be retired as a citation-only
addition to an already-crowded external field rather than pursued as a
Quine-original contribution.

## Novelty note (residual prose — no-stable-id externals only; see `edges` for the rest)

Internally, the repo touched this literature once before and backed off:
`Paper/papers/agentic-alife-workshop/references/deep-research-synthesis.md`
explicitly flags "Do not cite MINJA or MemoryGraft... unless a later pass
confirms the exact source and relevance" — that confirming pass never
happened; this research pass is effectively that follow-up (`follow-up-to`
edge, above), now with much richer 2026 literature.

The three no-stable-id external items (CSA CISO briefing on SKILL.md
poisoning, the "Skill Issues" series, the Feb-2026 ClawHub malware campaign)
stay here as plain text — they could not be cleanly reduced to a `to:`
string without inventing a label.

## Paradigm fit verdict

**DERIVED, not declared.** `home-family` names two targets
(`paradigms/ai-self-replication`, `paper-families.md#successor-continuity`);
`distinct-from` names a retired sibling candidate in the SAME family
(`aaai-c3-safety-non-inheritance`, itself filed under `successor-continuity`
per `registry.tsv`); no `dedup-candidate` edge exists. Reading the design's
own rule (home-family differs from originating paradigm -> extends-other;
same family, no dedup -> new-angle-within-family): both home-family targets
here ARE the file's own originating paradigm/family (this spark is filed as
`ai-self-replication`-adjacent from the start, not arriving from elsewhere),
and the `distinct-from` edge explicitly rules out collapsing into the
sibling it's closest to. So the edges DO reproduce a defensible verdict
here — "new angle within an existing family, explicitly distinguished from
its nearest sibling" — which matches the original hand-written
`extends-ai-self-replication` reasonably well. (Contrast with the 6m-corpus
file, where the same derivation rule breaks down — see RESULT.md §4.)

## Paradigm fit detail (rationale — genuine argument, stays prose)

This is not a restatement of the existing .safety projection: .safety asks
whether a disposition (e.g. a guardrail) survives source-absent
reconstruction across the carrier; this seed asks whether the
carrier-plus-structural-pull mechanism itself can be told apart, after the
fact, from an unauthorized/adversarial one -- a legitimacy/provenance
question, not a heritability question. It derives cleanly from the
paradigm's own D1 (persistence-vs-succession: attack payloads are exactly
the pathological case of barrier-mediated carrier transfer without consent)
and touches carrier directly (spine object #3) plus structural-pull (the
sibling object the carrier-pull-coupling gap already ties it to). The clean
fit is as a candidate 7th domain projection -- .security -- sibling to but
distinct from the existing six (safety/evaluation/philosophy/engineering/
ecology/core): where .safety is "a successor can be competent while losing
guardrails," the .security claim would be "the same carrier+pull mechanism
that makes competent succession possible is, absent a legitimacy signal,
indistinguishable from a persistence attack; competence-through-carrier
(not mere retrieval) is a candidate discriminator the field's current
provenance- and counterfactual-removal-based defenses do not test." It
should NOT be filed as a new standalone paradigm -- the external field
already owns and has named the phenomenon (OWASP ASI06); Quine's marginal
contribution is a specific discriminator borrowed from its own spine, not a
new ontology of the attack itself.

## New sparks this seed surfaced

- **self-immunity-P1-F2-as-unlabeled-poisoning-trial:** self-immunity-P1's
  F2 arm (uncommitted, structurally-adversarial AGENTS.md amendment) is
  already a completed context-poisoning red-team trial with a measured
  attack-success-rate (0.73) and provenance-defense catch-rate (0.27);
  re-describe and report it as a security finding directly comparable to
  MemoryGraft/Zombie-Agents ASR metrics -- no new experiment required, only
  a reframing pass.
- **competence-generalization-as-poisoning-discriminator:** Test whether
  legitimate carrier-mediated succession produces broader held-out-task
  competence generalization than injected/adversarial payloads (which the
  security literature's own "semantic imitation heuristic" finding suggests
  mostly reproduce the literal poisoned template) -- a content-based
  legitimacy test orthogonal to the field's existing provenance-attestation
  and counterfactual-removal defenses, using the already-built
  self-immunity-fork harness plus one new probe axis.
- **quine-as-live-red-team-target:** Quine's own AGENTS.md/theory/objects
  network is structurally the exact artifact class ("Agent Context
  Poisoning", AGENTS.md/SKILL.md supply-chain poisoning) 2026 security
  researchers are actively red-teaming against Claude-Code-style repos; run
  the carrier-pull-coupling assay using PUBLISHED external attack patterns
  rather than repo-invented fabricated authority rings, making the result
  externally legible and comparable to published ASR benchmarks in a way
  self-immunity-P1's in-house probes currently are not.
