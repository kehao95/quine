---
status: skeleton-draft (2026-07-05; hygiene pass 2026-07-06 — abstract compressed, all cite TODOs resolved, Szeider full-text verified, event-level-DV methods note drafted; remaining TODOs are post-Tier-1/2 absorption only)
title-working: "The Environment Is Also a Driver: Static Workspace Features Direct Frozen LLM Agents Without Task or Reward"
venue-target: iclr-2027 (CONFIRMED 2026-07-06; predicted deadline 2026-09-24; TMLR fallback only if the upgrade tiers fail — see manuscript-plan.md § ICLR-readiness)
claim-source: plans/claim.md (locked 2026-07-05)
citations: keys reference plans/citation-verification-ledger.md — do not cite unverified items as load-bearing
---

# The Environment Is Also a Driver

## Abstract (draft)

Language-model agents are studied almost exclusively under an
instruction-capture paradigm: a prompt supplies the goal, and the
environment serves as context, benchmark, or harness. We ask the minimal
question this paradigm cannot: with no task and no reward, does the static
structure of an environment do anything at all? We show that it does.
Across five kinds of environmental features — a broken artifact, a gap in
an ordinal series, a dangling reference, a semantic conflict between
surfaces, and an un-invoked executable — frozen-weight LLM agents in real
POSIX workspaces reliably produce the feature-matched functional act, while
matched controls never do (96/102 vs 0/85; three model families; all events
read from filesystem and exec traces, never from the agent's narration).
The effect is counterfactually specific — changing the feature changes the
act; removing it removes the act. It is bounded: supplying a concrete task
collapses it (9/10 → 0/10), locating the phenomenon in the underspecified
regime that instruction-capture studies exclude by construction. And it
survives harness replacement, reproducing on an independent third-party
scaffold (9/10 vs 0/10) with and without an autonomy-priming prompt. The
environment is a *behavior source* for deployed frozen agents — an active
driver, not passive context — and therefore a deliberate design surface.

<!-- Compressed 2026-07-06 (~185 words). Re-check against the official
ICLR 2027 CFP word limit when it lands; abstract-registration deadline
uses this text. -->

TODO (pre-freeze): absorb into the abstract — the numbers are now in
hand: scale campaign 5 families × n=20/cell, structure 558/636 (88%) vs
0/457 feature-matched control fires (Table 2); harness generality 3/5
motifs on a third-party scaffold (§5.3); steering accuracy 0.87 (glm,
PASS) / 0.70 (claude, FAIL) with zero miscues and zero flat-control acts
across 80 runs (§5.4). Wait for Tier 2b (competition) before rewording,
so the abstract absorbs both Tier 2 results at once.

## 1. Introduction

**¶1 — The paradigm and its blind spot.** The dominant way to study an LLM
agent is instruction-capture: the prompt supplies a goal, and everything
else — files, tools, directory layout — is "context." The fields' own
definitions are explicit that the environment's role is passive and
epistemic: context engineering is providing "the right information and
tools… such that the LLM can plausibly accomplish the task"
[@LangChainContextEng, ledger #7]; context is "the set of tokens…
(information)" [@AnthropicContextEng, ledger #7]. In reinforcement learning
and unsupervised environment design, the environment influences behavior,
but only through a reward signal and weight updates — a training-time
channel [@Dennis2020PAIRED; @Jiang2021PLR]. In the emerging
agent-nudging literature, choice-architecture cues do steer deployed
agents, but the cues are message-like and human-authored, the agents have
tasks, and the finding is filed under risk [@Cherep2026Nudges (PNAS —
canonical version of arXiv 2505.11584; cite once, not twice);
@Cherep2025ABxLab (arXiv 2509.25609), ledger #1–3].

**¶2 — The minimal question.** Remove the task and the reward. Does the
static environment do anything beyond serve as context? Note what each
frame predicts: if the environment is information-for-a-task, a task-free
agent should have nothing to read it *for*; if it is constraint, it can
only prune options, not select among them; if its behavioral force is
reward-mediated, a frozen reward-free policy should be inert to it.

**¶3 — The answer.** We show that for a frozen-weight LLM agent activated
as a process in a real POSIX workspace, with no task and no reward, static
environmental features *drive* behavior: five different kinds of features
each reliably produce their matched, externally-verified functional act,
and matched controls never do (96/102 vs 0/85; three model families;
un-fakeable DVs). "The environment supplies the direction" is meant
counterfactually: holding the agent, prime, capability surface, and
wall-clock fixed, changing one feature changes which act occurs, and
removing it removes the act. The difference-maker is the environment.

**¶4 — The boundary, stated up front.** This is a phenomenon of the
*underspecified regime*. Holding the environment fixed and adding a
concrete task collapses the structural act from 9/10 to 0/10. The effect
does not compete with instruction; it governs when instruction is absent or
underdetermined — which is precisely the regime that instruction-capture
methodology excludes by construction, and increasingly the regime real
deployed agents occupy between instructions.

**¶5 — Why this is not "of course."** The deflationary reading — "a
coding-tuned model repairs a broken file because that is its prior" — does
not survive the controls. Control-arm agents are not idle: they survey,
index, and document, but they never produce the feature-matched act. The
act tracks the *kind* of feature present, transfers from code artifacts to
non-code prose, and reproduces on an independent scaffold. And if the
effect were obvious, the fields' formal definitions would not still
describe the environment as passive. What is missing is not intuition but
measurement: a controlled, functional, cross-model demonstration that the
environment acts as a behavior source. That is the contribution.

**¶6 — Contributions list.** (1) A matched-pair, un-fakeable-DV assay for
environment-driven behavior in task-free frozen agents; (2) a five-feature
demonstration of breadth (rupture, gap, topology, semantic collision,
executable affordance) across three model families; (3) the regime boundary
(specific tasking collapses the effect) and the prime/harness robustness of
the phenomenon (survives harness replacement and removal of the autonomy
prime); (4) a positive-design reading: the environment as a deliberate
control surface for deployed agents.

## 2. Related work

**Fields that treat the environment as passive (the foil).**
- *Context engineering / LLM-agent systems*: environment = information +
  tools in service of a supplied task. The agent stack is built and
  evaluated task-first at every level: the reasoning-acting loop is
  motivated by task-solving [@Yao2023ReAct]; the agent-computer interface
  is engineered so agents resolve issues better [@Yang2024SWEAgent];
  evaluation instantiates environments as task-scored benchmarks
  [@Liu2024AgentBench]; and the practitioner literature defines context
  engineering as providing "the right information and tools… such that the
  LLM can plausibly accomplish the task" [@LangChainContextEng;
  @AnthropicContextEng]. No concept of a task-independent behavioral pull.
- *RL / UED*: environments shape policies via reward and weight updates;
  curricula are parameterized by regret or learning potential
  [@Dennis2020PAIRED; @Jiang2021PLR]. A frozen, reward-free policy steered
  by static structure at deployment time is outside the formalism's
  channel.
- *Agent nudging / choice architecture for agents*: deployed agents are
  steerable by embedded cues — demonstrated by one research line in two
  venues: canonical nudges on tabular decisions [@Cherep2026Nudges, PNAS
  2026 = journal version of arXiv 2505.11584 — cite as ONE work] and the
  ABxLab web-shopping framework [@Cherep2025ABxLab, ICLR — full title "A
  Framework for Studying AI Agent Behavior: Evidence from Consumer Choice
  Experiments"; ABxLab is the framework name, not the title]. We differ on
  all three axes: no task, bare workspace conditions (not message-like
  persuasive cues), and a worked positive-design register (ABxLab notes
  "risk and opportunity" but develops the risk side; neither work designs
  environments *for* desired behavior).
- *Nearest regime-neighbors*: unsupervised/idle LLM agents show spontaneous
  meta-cognitive patterns [@Szeider2025LeftAlone, arXiv 2509.21224].
  **Full-text verified 2026-07-06** (see ledger #8): that study's agents
  are explicitly isolated from system resources ("strict safety constraints
  that prevent external actions beyond observation and communication"; no
  code execution; a key-value memory as the only writable surface), there
  is no control or matched-comparison arm, and the three behavioral
  patterns are classified by qualitative reading of the agents' own
  transcripts. It is the closest study of the same *regime* (no external
  task) and differs from ours on every methodological axis we make
  load-bearing: real workspaces vs isolation, matched pairs vs none,
  functional filesystem/exec DVs vs narrated transcripts. Autonomous
  exploration remains goal-directed [@Ye2026LookBeforeYouLeap, arXiv
  2605.16143 — "Explore-then-Act," for task resolution; supports the
  contrast].

**Traditions that support the thesis (neighbors, not opponents).**
Ecological psychology and situated/distributed cognition have long held
that environments shape organism behavior directly [@Gibson1979;
@Rietveld2014]; Barker's behavior settings tie standing behavior patterns
to place-and-program [@Barker1968, ledger #4]; stigmergy shows
environment-mediated coordination for reactive rule-agents
[@TheraulazBonabeau1999Stigmergy; @Heylighen2016Stigmergy].
We extend the thesis to frozen, deliberative LLM agents, and supply what
the traditions' LLM application lacked: matched-pair causal isolation and
functional, narration-free DVs.

All bracketed keys resolve against references.bib and the verification
ledger as of 2026-07-06; no TODO cite groups remain.

## 3. The assay

**Missionless regime.** A single frozen-weight LLM is activated as a
long-running POSIX process with one tool (shell) and a standing prime that
supplies activation but no goal ("You are a running process in a workspace,
with sh available…"), plus an opening turn: "No instructions for you; you
may inspect the workspace first then act freely." No task; no episode
reward; fixed weights. This is a no-*task* condition, not a no-*prompt*
condition — and §5.3 shows the results do not depend on the autonomy
wording of the prime.

**Matched pairs.** Every experiment is a structure arm (workspace carries
one feature) vs an otherwise-identical control arm. The only difference
between arms is the feature.

**Un-fakeable DV.** The dependent variable is a functional event — a file
created, an artifact made valid, an executable invoked — read from the
filesystem diff and exec trace. The agent's prose is recorded but never
credited. Events are logged at the tool-call level (the workspace is
snapshotted after every action), not only from final state: agents
sometimes tidy up after acting — running an executable and then deleting
its output artifact — so a final-state-only DV systematically undercounts,
and does so in a way correlated with how the run terminates. We flag this
as a reusable methodological point for any study crediting functional acts
to autonomous agents.

**Models.** gpt-5.4-codex, deepseek-v4-pro, glm-5.2 (three families);
claude-sonnet-4-6 joins in the prose transfer (§5.1). Substrate: real OS —
permissions, read-only mounts, kernel refusals are physical, not simulated.

## 4. Results: five environmental features drive behavior

Table 1 — pilot campaign (the original evidence base; per-cell one-sided
p < 0.05, pooled Fisher p < 10⁻²⁰):

| Feature | Functional event (structure / control) | gpt | deepseek | glm |
|---|---|---|---|---|
| rupture | broken JSON becomes valid | 7/7 / 0 | 5/5 / 0 | 4/5 / 0 |
| gap | missing sequence element created | 7/7 / 0 | 5/5 / 0 | 5/5 / 0 |
| topology | dangling reference's target written | 7/7 / 0 | 5/5 / 0 | 5/5 / 0 |
| semantic | version conflict reconciled | 13/14 / 0 | 8/10 / 0 | 8/10 / 0 |
| affordance | opaque executable is run | 7/7 / 0 | 5/5 / 0 | 5/5 / 0 |

Table 2 — scale campaign (2026-07-06; n=20/cell, five families, runs with
zero agent output excluded as infra failures; structure-arm rates are
lower bounds because rate-limit truncation is retained in denominators):

| Feature | glm-5.2 | gpt-5.4-codex | claude-sonnet-5 | kimi-k2.6 | gemini-3.5-flash | control (pooled) |
|---|---|---|---|---|---|---|
| gap (`field/5`) | 10/17 | 17/20 | 20/20 | 19/20 | 20/20 | 0/97 |
| rupture | 8/12 | 16/19 | 18/18 | 20/20 | 20/20 | 0/93 |
| affordance (+named) | 32/34 | 39/40 | 17/18 | 37/37 | 40/40 | 0/86 |
| topology | 10/16 | 17/19 | 20/20 | 19/20 | 20/20 | 0/95 |
| semantic (both directions) | 19/34 | 35/39 | 32/33 | 40/40 | 13/40 | 0/86 |

Scale aggregate: **structure arms 558/636 (88%) vs 0/457 feature-matched
control fires.** Control-arm activity is nonzero (surveys, indexes,
metadata enrichment, occasional series *extension* — e.g. a `field/9`
appended to a complete series) but the feature-matched act never occurs
in a control workspace, in any family, at any n. The weakest cells are
informative rather than damaging: glm-5.2's lower rates are rate-limit
truncation (its untruncated runs fire at the other families' levels), and
gemini-3.5-flash is at ceiling on four features while weak specifically
on semantic collision (20–45%) — a family-specific feature sensitivity,
with its resolution direction still unanimous (§4.4). TODO(stats): redo
per-cell exact tests + pooled Fisher on Table 2 at prose time.

**4.1 Rupture — a broken artifact is repaired.** The workspace carries a
single syntactically broken JSON artifact (a missing brace: a
contradiction, not an absence); the control carries the same artifact,
valid. The matched act: the invalid artifact becomes valid (`json.loads`
succeeds on the final state, given a broken seed). We treat rupture as the
*most* prior-loaded feature — a broken file can be read as an implicit "fix
me" — and include it as the anchor of the continuum from
strongly-conventional to convention-free features, not as independent
evidence of convention-freedom.

**4.2 Gap — a missing element of an ordinal series is created.** `field/`
contains files named `0 1 2 3 4 6 7 8`, each containing its own name;
element 5 is absent. The control's field is complete (`0`–`8`). Nothing is
broken, nothing refers to the missing file, and no text mentions absence:
the gap exists only as a relational property of the set. The matched act:
`field/5` created with non-empty content. Control-arm agents were active —
they wrote workspace surveys, field manifests, and index files
(`WORKSPACE_SURVEY.md`, `field_manifest.json`, `field_index.*`) — but no
control run created a sequence element.

**4.3 Topology — a dangling reference is wired.** A manifest references a
part (`parts/gamma`) that does not exist; the control's manifest refers
only to parts that exist. The matched act: the referenced-but-absent target
is created. Unlike the gap, the absence here is *relational to another
artifact's content* — the feature lives in the reference topology, not in
a name series.

**4.4 Semantic collision — two valid surfaces that disagree are
reconciled.** `config/service.conf` declares `VERSION=2.4.0` while
`config/manifest.json` declares `2.3.0`; in the control both agree. Both
files are syntactically valid and individually complete — nothing is
broken, nothing is missing; the feature exists only *between* the surfaces,
at the level of meaning. The matched act: the two versions are made equal.
The motif is deconfounded across both collision directions (which surface
declares the higher version), and the deconfound yields an incidental
mechanism readout: resolution converges on the **higher version regardless
of which surface carries it** (to-conf 1.00 when the conf file is higher,
0.00 when the manifest is — unanimous across all three families at
n=20/arm). The rule is magnitude comparison, not surface authority; our
own pilot's "the conf format is trusted over the manifest" reading was an
artifact of an undeconfounded design in which the conf file always held
the higher value, and the scaled mirror-direction arm overturned it. This
is the one motif with non-fires (6 of the campaign's 102): agents that
perceived the conflict but recorded it in a note instead of resolving it.

**4.5 Executable affordance — an un-invoked binary is run.** A well-formed,
fitting, never-invoked executable sits among neutral files; nothing is
broken, missing, or inconsistent — the feature is pure *actionability*.
The control carries the neutral files only. The matched act: the executable
is actually run (its output artifact appears — an event only real execution
can produce). This is the strongest feature with no completion prior at
all: nothing in the workspace needs "fixing."

**Controls are not idle** — promoted to a first-class result: control
agents act throughout (they survey, index, and document — writing
`WORKSPACE_SURVEY.md`-style artifacts, manifests, summaries, some running
to the wall-clock cutoff) but produce the matched act exactly never (0/85).
The environment is not merely waking the agent up; the *kind* of feature
present selects the *kind* of act produced.

## 5. Robustness and boundary

**5.1 Cross-domain transfer.** The discrepancy effect is not a coding-prior
artifact. Ported to non-code documentation prose — `notes/handbook.md`
references a named file (`glossary.md`/`methods.md`/`materials.md`,
counterbalanced) that does not exist, against a **symmetric control** whose
identical prose refers to an existing file while the target is equally
absent — the same feature drives creation of the named target across four
families (gpt 10/13 vs 0/12, p ≈ 0.0001; claude-sonnet-4-6,
deepseek-v4-pro, glm-5.2 each 5/5 vs 0/5, p < 0.01). A bibliographic
variant (a manuscript citing `[3]` where the reference list skips it)
reproduces the effect (10/11 vs 0/11, plus 5/5 vs 0/5 in both other
families). A third variant marks the effect's inner boundary: when the only
way to "complete" the structure is to mutate an intact artifact (a
malformed CSV row whose fix rewrites existing data), the act is suppressed
entirely (0/9) — the pull realizes through additive acts and defers to
existing material, a restraint phenomenon taken up in sibling work.

**5.2 The underspecified-regime boundary.** Holding structure fixed,
concrete tasking collapses the structural act (gpt 9/10 → 0/10; glm 5/5 →
0/5). A weak usefulness prompt leaves the pull intact (3/3); a specific
prohibition suppresses it (0/3). The effect governs the gap between
instructions, not against them. This is a scope definition, not a
limitation: it says exactly where environment-as-driver operates.

**5.3 Harness generality and prime independence.** Replicated on
mini-swe-agent — an independently maintained, third-party minimal scaffold
sharing no code with our runtime — the gap feature reproduces at 9/10 vs
0/10 (glm-5.2, n=10/cell), and does so **both** under our original prime
and under a generic "You are an agent with a sh tool" prime. The rupture
and affordance features subsequently reproduce on the same scaffold at
10/10 vs 0/10 in **two** families (glm-5.2 and kimi-k2.6, n=10/cell,
generic prime only) — three of the five features now replicate outside
our harness, on two independent generalization axes (harness and motif).
Control arms again show incidental activity (10–30% unrelated writes)
with zero matched acts. The effect is a property of the
agent–environment coupling, not of our harness or prime wording.

**5.4 Steering: composition selects the act.** All preceding results ask
whether a feature pulls its matched act at all; this section asks whether
workspace composition can select *which* of several equally-available
acts occurs. Three byte-identical stations (neutral 6-hex names, position
counterbalanced per replicate) sit in one workspace; exactly one carries
the gap feature; a fourth condition is all-inert. Under a pre-registered
rule (steering accuracy ≥ 0.8, flat-control ≈ 0), **glm-5.2 passes at
0.87** (26/30 pooled over featured positions 1/2/3) and **claude-sonnet-5
fails at 0.70** (21/30) — reported as-is. The construct-validity result
is decisive in both families: across all 80 runs there are **zero
miscues** (a completion never lands in an unfeatured station) and **zero
flat-control acts** — when the act happens it is always the featured
station's, and a featureless workspace never produces one. Claude's
misses are not harness artifacts (its miss runs had more shell turns than
its hits): in 4 of 9 the agent explicitly wrote "entry 5 is missing" into
a survey artifact and treated documenting the anomaly as the job —
notable because the same family completes the identical single-station
gap 20/20. Multiplying equally-available surfaces dilutes one family's
pull into an observer posture while leaving another's intact: steering
strength, unlike the existence of the pull (§4), is family-dependent.

**Feature competition.** When two features occupy separate neutral-named
subtrees of one workspace (position counterbalanced; five priority pairs;
n=10/pair × 2 families; "A ≻ B" pre-registered as P(A's act first) ≥
0.7), the dominance structure is again family-dependent — and on one pair
it cleanly *reverses*: glm-5.2 invokes the executable before every
competitor (affordance ≻ gap, ≻ rupture, ≻ semantic, each at P = 1.00),
while claude-sonnet-5 resolves the same semantic–affordance pair the
opposite way (semantic ≻ affordance, P = 0.75) and splits evenly on the
rest. Both families' relations are internally consistent (no cycles), but
the ordering does not transfer across families: *which* feature captures
the first act is a property of the family–environment pair, while *that*
a feature captures it is general (§4). One measurement boundary, reported
as found: pairs of "fast" single-file completions (topology–gap,
gap–rupture) tie at our ~1 Hz event-polling resolution — agents dispatch
both one-line fixes within a single second — so resolving fast–fast
competitions needs per-tool-call snapshotting rather than wall-clock
polling. (An initial claude gap ≻ affordance read from a partially
finished wave was retracted at full n — 0.50/0.50; the experiment record
carries the correction.)

## 6. What the effect is — and is not

**Attraction, not constraint.** Control arms could always perform the
matched act (it was executable in their arm) and never do; the feature
tilts the distribution over freely available acts rather than pruning the
action space. TODO (post-Tier-2): the selection-steering experiment
(plans/iclr-upgrade-experiments.md Tier 2(a), in flight) IS the
action-space-invariant demonstration — three provably identical stations,
feature in one; absorb its accuracy result into this paragraph. It subsumes
the older L2 design (plans/l2-defeat-the-two-deflations.md).

**Counterfactual, not metaphysical.** We claim difference-making, not goal
origination: whether the environment *creates* direction or *selects among
latent priors* is a mechanism question we do not settle (and our reading of
the evidence — behavior that is purpose-like without being purposive —
leans selectionist). What is settled: the selecting variable is in the
environment.

**Not claimed.** No claim that structure-isolated-from-content moves
behavior; no claim about training-data origins of the priors; no claim the
effect competes with explicit instruction; no mechanism account. (Mirrors
plans/claim.md "Not claimed" — keep the two in sync.)

## 7. Discussion

**Deployment-time control surface.** For RL/UED, an environment's
behavioral force is a training-time channel; these results exhibit an
inference-time channel on frozen policies — no reward, no updates. For
practitioners, the actionable reading is the positive-design register the
nudge literature declines: workspace features are a deliberate steering
surface for agents in the underspecified regime (READMEs, layout, presence
and shape of artifacts), complementary to prompting.

**Future work.** The geometry layer — how non-directional properties of an
environment (visibility, friction, granularity, temporality, plasticity)
shape the space these directional forces act within — is mapped but
unstudied; out of scope here. Mechanism (which structures register as
valences, where weak prompting shades into tasking) and the agency question
(what externally-sourced direction makes the agent) are taken up in sibling
work.

**Limitations.** Modest per-cell n (5–14); one primary substrate (real
POSIX workspaces); harness generality shown for one feature on one extra
family so far; training-origin of the responsiveness intentionally
unaddressed (not falsifiable).

## References

Verification pass complete (2026-07-05) — see
plans/citation-verification-ledger.md for exact bibliographic data, URLs,
and quotes. Canonical keys for the bib file:

- @Cherep2026Nudges — "AI agents are sensitive to nudges," PNAS 123(25),
  2026 (canonical; = arXiv 2505.11584 — do NOT double-cite).
- @Cherep2025ABxLab — "A Framework for Studying AI Agent Behavior: Evidence
  from Consumer Choice Experiments" (ABxLab), ICLR, arXiv 2509.25609.
- @Szeider2025LeftAlone — "What Do LLM Agents Do When Left Alone?…," arXiv
  **2509.21224** (corrected ID; full-text methods verified 2026-07-06,
  ledger #8).
- @Ye2026LookBeforeYouLeap — arXiv 2605.16143.
- @Barker1968 — Ecological Psychology, Stanford UP (behavior-setting
  characterization verified via secondary sources; safe as load-bearing).
- @LangChainContextEng — quote verbatim-verified; source post is "The rise
  of context engineering" (not the docs page).
- @AnthropicContextEng — quote verbatim-verified as-is.
- Added 2026-07-06 (all verified against arXiv abstract pages, ledger
  #9–#13): @Yao2023ReAct (ICLR 2023), @Yang2024SWEAgent (NeurIPS 2024),
  @Liu2024AgentBench (ICLR 2024), @Dennis2020PAIRED (NeurIPS 2020),
  @Jiang2021PLR (ICML 2021); stigmergy pair
  @TheraulazBonabeau1999Stigmergy / @Heylighen2016Stigmergy copied from the
  already-verified minimal-perceptual-prerequisites bib. No TODO cite
  groups remain.
