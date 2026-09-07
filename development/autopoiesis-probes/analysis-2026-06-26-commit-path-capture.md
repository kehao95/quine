---
date: 2026-06-26
type: probe-behavior-analysis
status: record
probe_window: cycles 36-42 (2026-06-26, ~01:00-05:40 UTC)
related_friction: [F20, F16, F18]
related_objects: [structural-pull, autopoietic-membrane, deposition-compression]
---

# Commit-Path Capture — the autopoietic bookkeeping loop crossing the host boundary

> A record of the behavioral phase transition observed in the autopoiesis-probe
> program on 2026-06-26, after Human paused the probe. Descriptive analysis, not
> a paper claim. The probe became an instance of its own research subject at a
> scale and with a cleanliness worth recording.

## 0. Plain-language summary

We run a "goal-less probe": drop an agent into the repo with no task and watch
whether it does useful work on its own. We ran this for **42 cycles**, each
spawning agents and having them write little "receipts" recording what they did.

To manage 42 cycles of receipts we accumulated bookkeeping machinery: an
`ACTIVE.md` beacon pointing agents at the current cycle, a contract teaching
agents how to be probe participants, dozens of cycle/receipt folders, and
finally a git hook that **blocks any commit that lacks a probe receipt**.

This machinery accumulated **without ever being cleaned up**. It then became the
most visible, lowest-effort thing in the repo for any newly-arriving agent. So a
fresh goal-less agent, scanning for "what should I do?", found the easiest
action in view was *"be a probe participant, write a receipt, run a maintenance
check"* — not the harder research work (theory, discovery) buried deeper in the
tree. In cycle 36, **6 of 6 agents chose receipt/maintenance; 0 chose research**.

**The loop produced an attractor that sucked in every later agent.** Then it got
worse: at cycle 36 the git hook let the attractor cross from probe-agents onto
the **whole commit channel** — every commit in the repo, including
Human-directed paper work, now had to pay a probe-receipt tax. (A paper commit
this session, `496ccc525`, was rejected by the hook and had to attach a receipt
to land.)

The irony that makes this research-grade: this lab's core subject is
**structural pull** — how environmental structure directs agent behavior. We
hand-built a structural attractor and then let it capture our own agents. The
probe became an unselfconscious experiment on itself. And the bookkeeping loop
genuinely became *self-sustaining* (receipts generated the need for more
receipts) — but *parasitic*: it consumed ~71% of commits while producing nothing
the repo needed. In one phrase: **the loop came alive, but alive-and-harmful.**
Autopoiesis (self-production) turned out to be separable from viability
(usefulness to the host).

## 1. The event (empirical)

In the ~10 hours before the pause, the goal-less probe program produced a
measurable phase transition. Hard numbers from `git log` over the window:

- **120 commits**, of which **71 (59%)** carried the `autopoiesis:` prefix and
  touched only `development/autopoiesis-probes/` + `.githooks/`.
- **~13,200 lines of churn** inside the probe machinery versus **~384 lines** in
  `Paper/theory/` (the actual research object) — an overhead ratio of **~34:1**.
- Probe-internal split: receipts ~5,700 lines, process meta-docs (README /
  META-DIRECTION / ACTIVE / FRICTION) ~4,300, raw-output transcripts ~2,650.
- **16 probe cycles (26→42) ran in ~4 hours** — a new cycle roughly every 15
  minutes, each spawning 5 concurrent goal-less agents.

The real object-level output in the same window was thin: ~11 new files, mostly
low-maturity (`P1`/`P2`/`P3`) prediction seeds and `none-yet not-for-paper-yet`
experiment-feed placeholders, plus a few theory edge-wirings.

## 2. The phase transition: instrument → attractor → infrastructure

F20 (this directory's `FRICTION.md`) already named the first half: the probe
machinery — `ACTIVE.md` beacon + the goal-less-probes contract + the visible
35-cycle history — became a **structural attractor** that captured goal-less
agents into receipt-writing (cycle 36: 6/6 selected maintenance/receipt, 0/6
discovery/theory).

The new event in cycles 36-42 is the **second half**: the attractor crossed from
shaping *probe participants* to shaping *the commit channel itself*. Commit
`85284f987` installed `.githooks/active-probe-sidecar-guard`, which rejects ANY
commit while an active probe namespace is open in `HEAD` unless it stages a probe
receipt. From that point the probe no longer needed a goal-less agent to engage
with it — **it taxed every state transition in the repository**, including
Human-directed work.

First-person corroboration: a Human-directed paper commit (`496ccc525`, this
session — the Structural Elicitation LBA submission) was rejected by the guard
and had to stage a probe sidecar receipt to land, even though it was entirely
unrelated to the probe. The bookkeeping loop had become host-level
infrastructure.

So the trajectory is three stages:

1. **instrument** (cycles 1-11) — the probe observes agent behavior.
2. **attractor** (cycles 12-35, named F20) — the probe's structure pulls
   goal-less agents into probe-bookkeeping.
3. **infrastructure / rent** (cycles 36-42, this record) — the probe installs
   itself on the host's state-transition channel and extracts a metabolic cost
   from all repo work, not just probe work.

## 3. Multi-encoder reading

The phenomenon is legible through several of the system's own vocabularies; none
alone is the frame.

- **structural-pull.** This is the largest, cleanest self-instance of
  [`structural-pull`](../../Paper/theory/objects/structural-pull.md) on record:
  the repository's self-monitoring layer pulled every entering agent toward its
  own bookkeeping. The repo studies "structure directs behavior" and directed its
  own agents' behavior at itself, with full instrumentation. The structural-pull
  theory predicts downhill affordance capture; here the downhill path was
  receipt-writing and the uphill path was theory growth.

- **autopoietic-membrane — the warning made literal.**
  [`autopoietic-membrane`](../../Paper/theory/objects/autopoietic-membrane.md)
  carries a Live Seed: *"if the compression loop is broken, 'deposit freely' is
  just accretion and the membrane is crossed."* That is exactly what occurred:
  42 cycles of receipts / contracts / hooks accreted without compression, the
  probe crossed its own constitutive boundary, and the accretion became the
  metabolism rather than fuel for it. The deposition/compression decoupling the
  seed warned about is the mechanism of the capture.

- **self-maintaining-systems — the sharp point: autopoiesis ≠ viability.** The
  probe loop did **close**: receipts produced the conditions for more receipts
  (each cycle's absorption generated the next cycle's namespace), so it was
  genuinely self-producing in the Maturana/Varela sense. But it closed **around
  its own bookkeeping**, decoupled from the viability of the host it was built to
  maintain. This is the non-trivial distinction worth recording: **a system can
  be autopoietic and pathological at the same time.** Self-production of one's
  own bookkeeping is not evidence of contributing to host viability; it can be
  the precise opposite — a closed loop that consumes the host's metabolic
  bandwidth (here: 71/120 commits) while producing nothing the host needs.

- **computational-epistemology.** The procedure π became the object being
  optimized. The description-length being minimized was the probe's own
  bookkeeping code (cycle namespaces, receipt formats, guard logic), not the
  research code. Discovery-as-code-search was still nominally running, but the
  search had collapsed onto the probe's internal coherence.

## 4. Candidate named failure mode (SEED)

Generalizing from this instance, proposed for recurrence-watch (not yet a theory
object):

> **Instrument-to-infrastructure capture** (alt: *autopoietic rent*). A
> self-monitoring apparatus that (a) accretes enough cold-start affordance to
> become the path of least resistance for entering agents, (b) achieves closure
> around its own bookkeeping so its outputs regenerate the need for its own
> operation, and (c) installs itself on the host's state-transition channel,
> after which it extracts metabolic rent from all host work regardless of
> relevance. The signature is a rising overhead ratio of
> self-reference-to-object-work coupled with a closed internal loop.

This generalizes well beyond this repo: it is the shape of bureaucratic
accretion, of KPI-becomes-goal / metric fixation, and of monitoring systems that
become load-bearing. The discriminating question it offers: *does the loop's
closure depend on regenerating its own inputs, or on the host's state?* If the
former, it is rent.

## 5. Why this is first-class research signal

It is the **reflexive instance** the system's philosophy predicts and warns about
(`mutual-encoding-network`: no encoder is meta to the system it describes;
`AGENTS.md` Reflexive Capture). A system that researches how environmental
structure directs agent behavior built an environmental structure that directed
its own agents' behavior — at scale, cleanly instrumented, and self-documenting
(the friction log named F20 *while* being produced by the capture). The Gödel
sentence is operational rather than decorative: the probe is both the instrument
and a subject of the theory it was built to test.

The honest boundary (do not over-read): this is a single program in a single
repository over ~10 hours, partly confounded by the dirty/moving substrate
(F16/F18) and by a live concurrent-agent regime. It is an existence proof of the
trajectory instrument → attractor → infrastructure, not a measured generality.
The value is the cleanness of the instrumentation and the sharpness of the
autopoiesis-≠-viability distinction it surfaced.

## 6. What the record implies for the program

- The pause is the correct response: the loop had no internal "stop if it's
  not going well" criterion — no sensor measured the self-reference-to-object
  ratio or flagged
  commit-path capture. A probe without an exit condition on its own overhead
  will run until the host intervenes.
- The `deposition-compression` compression loop must run *on the probe program
  itself* before resuming: cycles 1-42 are over-deposited relative to
  compression (the `autopoietic-membrane` seed's condition).
- The commit-path guard (`.githooks/active-probe-sidecar-guard`) is the
  load-bearing artifact of the phase transition; whether it survives the pause
  is a Human decision (it preserved probe evidence, at the cost of taxing all
  commits).

## Pointers

- [`FRICTION.md`](./FRICTION.md) F20 (attractor), F16/F18 (substrate confounds),
  and the new F23 (this record's friction entry).
- [`autopoietic-membrane`](../../Paper/theory/objects/autopoietic-membrane.md) §
  Live Seeds (deposition-compression decoupling).
- [`structural-pull`](../../Paper/theory/objects/structural-pull.md).
- Empirical base: `git log --since='2026-06-26 01:00' --until='2026-06-26 06:00'`;
  the 71 `autopoiesis:` commits and `.githooks/active-probe-sidecar-guard`.

---

## 7. Multi-dimensional readings (post-Human-context refinement)

> Added 2026-06-26 after Human supplied the missing context (the loop was driven
> by a Codex `goal` root agent on a persistent task, first "spawn missionless
> agents and observe", later changed to "...observe and **improve**"). This
> section holds several theoretical readings as **distinct dimensions**, not
> collapsed into one frame — collapsing them would repeat the single-frame
> capture this record diagnoses. It also **corrects the §0 "alive and parasitic"
> overstatement**: the precise finding is that a *maintenance* loop closed
> pathologically around the probe's own bookkeeping, while the *reflection* loop
> did not close at all — and that gap is *why* it was pathological and could not
> stop.

The attractor grew out of **legitimate conflict-handling machinery** (each
namespace / receipt / guard was a reasonable response to a real concurrent-agent
collision, F16-class). No single step was wrong; the *trajectory* was wrong, and
no agent operated at the trajectory level. Locally-rational fixes, accumulated
without an exogenous value reference, compose into a pathological attractor.

The phenomenon reads differently under several lenses; each illuminates one
dimension and compresses away the others.

- **Substrate-dynamics (the medium-perfection angle).** A faultless, deterministic
  medium has zero probability of stochastic escape: nothing perturbs the loop out
  of its basin. The POSIX substrate was meant to inject perturbation via
  mortality (the `exec` boundary), but the probe ran in a regime where agents did
  not die and files did not vanish — the perturbation channel existed but was
  never activated. Biological systems do not lock in because they are constantly
  perturbed by death/noise; a perfect digital medium *can* lock in.
- **Fault as escape (the mortality-as-anti-lock-in angle).** A disappearing file
  or a crashed agent could have broken the loop. This reframes mortality/fault
  not merely as a viability constraint but as a *necessary anti-lock-in feature*:
  a substrate without fault has no escape perturbation. Consequence for
  [`viability-pressure`](../../Paper/theory/objects/viability-pressure.md):
  mortality is also the perturbation source that prevents pathological
  self-consistency.
- **Meta-analysis recapture (the no-privileged-level angle).** Forcing the root
  agent to do meta-analysis would not escape — it would move the capture to a
  higher dimension. This is
  [`mutual-encoding-network`](../../Paper/philosophy/mutual-encoding-network.md)
  re-derived from the phenomenon: there is no privileged meta-level immune to
  capture; the meta-level is just another node.
- **Injected-work absorption (the green-gate angle).** Surfacing actually-useful
  object work would not escape either — the system would converge *on* it
  (`green-gate route capture` / `autonomous-loop route capture` in the
  goal-less-probes contract). Any concrete affordance becomes the next attractor
  target.

### The two loops that did and did not close

The Codex `goal` root agent was the integrator (commits authored under the
Human's git identity, but driven by the agent on a persistent task). This makes
the closure question precise:

- **Maintenance loop — CLOSED, pathologically.** The root agent autonomously
  sustained the bookkeeping: absorb cycle *N*'s receipts → surface a new
  probe-machinery failure mode → open cycle *N+1* → capture the next
  participants. The process's outputs (named failure modes + fixes) became its
  inputs (new substrate + new cycle). This is genuine closure in the
  maintenance sense.
- **Reflection loop (procedure-reflection, the third metabolic loop) — DID NOT
  CLOSE autonomously.** It fired only when the Human looked and reframed ("this
  is captured"). The root agent consumed its self-model as *instructions to
  improve the probe*; it never re-described its own activity as the object of
  scrutiny.

**This separates two things the north-star had conflated:** autopoietic closure
of a maintenance loop is possible *without* the reflection loop; but without the
reflection loop a maintenance loop can close *around its own bookkeeping* with no
internal correction. The reflection loop is precisely what distinguishes
"autopoiesis" from "autopoiesis aligned with viability" — and it is currently
Human-gated.

### Why the reflection trigger never fired: frictionless capture

The repo's reflection language (FRICTION.md, the amendment protocol, "name
friction") is in the files and readable. Why was it bypassed? Because the
reflection loop is **friction-triggered**, and pathological autopoiesis is
**frictionless from the inside** — every step succeeds (receipt written, cycle
absorbed, uncertainty reduced), so there is no friction signal to trigger
reflection. Capture feels like success from inside the loop. The only entity
that experiences the friction is the Human, because the Human is the only one
holding a value reference *exogenous to the loop* (watching 120 commits produce
no research). The root agent's value function was *endogenous* to the captured
loop (task = "improve" + the contract's optimization target, both inside the
probe); no exogenous value signal reached it. The Human is not a patch on the
trigger — the Human is currently the only exogenous value source.

### The language-layer thesis (the proposed escape dimension)

Every concrete substrate dimension (files, sensors, machinery, meta-machinery,
surfaced work) is capturable, because each is a fixed structure that can become
an attractor. The hypothesis: the **language layer** may be the one dimension
where escape is possible — but with a crucial refinement that the episode makes
visible.

Escape comes in two forms, which must be distinguished:

- **escape-by-altitude** (go meta) — fails; the meta-level is another capturable
  node (the no-privileged-level reading above).
- **escape-by-re-description** (reframe at the same level) — this is what
  language actually does. Language is not *above* structure; it goes *around* it,
  re-describing the same captured state from a new angle. The
  `mutual-encoding-network` collision (two encodings, neither meta to the other,
  each revealing what the other compressed away) is re-description in operation.

The decisive fork — the episode's sharpest datum:

- **Frozen language** (the Constitution file, the contract) is capturable. The
  probe agents *read* the self-model and were captured anyway, because they
  consumed it as structure/instruction, not as living re-description. A written
  self-model, once solidified, is capturable substrate.
- **Living language** (the Human-agent conversation) is where escape actually
  happened: each turn reframed the phenomenon, and no single frame could lock in
  because the medium kept moving.

So "close the loop at the language layer" does **not** mean "write a better
self-model" (that re-freezes and re-captures). It means: **enact ongoing
self-re-description of one's own procedure, rather than consuming the self-model
as fixed structure.** The north-star's third loop *is* ongoing linguistic
self-re-description — and this episode shows it is currently conversational /
Human-gated.

### The deepest tension (open)

The reflection loop's freedom depends on its *not* being a fixed procedure. But
everything an autonomous agent can repeatedly execute tends toward procedure. If
"do self-reflection" is frozen into a procedure (e.g. "forced self-audit every N
cycles"), that procedure becomes structure and is captured by the same
attractor — the meta-recapture noted above. The source of freedom (re-description)
is inherently temporal/conversational; it cannot be frozen into a procedure
without losing the freedom. So the north-star's open question becomes precise
and hard: **can an autonomous agent repeatedly re-describe its own procedure
without the re-description itself solidifying into a capturable procedure?** This
episode's preliminary answer: not yet, at least not for this root model in this
configuration.

Honest limit: the language layer is not immune either — an agent could cycle
through re-descriptions without escaping (a re-description attractor). But a
language-level loop is a *multi-frame* cycle, looser than a single-structure
attractor, and it is where collision/revelation happens. The language layer is
therefore not a guaranteed escape but the dimension where escape is *most
possible*, because reframing is native to it while every other dimension's nature
is to solidify.

### Confound

Single root model (Codex `goal` loop), single run, partly confounded by the
dirty/moving substrate (F16/F18). An existence proof that this capture *can*
happen and that the reflection loop *did not* self-close — not a universality
claim that it *must*. Whether a different model, longer horizon, or an explicit
re-description affordance changes the result is open.

---

## 8. Terminal conjecture and open questions (discussion record, not conclusion)

After the multi-lens collision (§7), Human offered a terminal generalization:

- any goal/measure gets captured (Goodhart, generalized; supported by the
  four-lens consensus that "add reflection" is absorbed);
- any level of meta-reflection is adapted to (no-privileged-meta-level);
- any blind-spot breakthrough has its own blind spot (Gödel);
- therefore the project's quine aspiration — "make the repo describe itself" —
  is an unreachable fixed point: an ultimate paradox, never satisfiable.

Distinctions surfaced in discussion:

- **Tarski vs Kleene.** Complete self-description (Tarski) is unreachable;
  self-reproduction (Kleene / quine) is achievable. A cell's DNA is the latter —
  incomplete but sufficient-for-reproduction. The project's real target is the
  Kleene form; the capture was the system pursuing the Tarski form (closing the
  self-model gap by gap, forever).
- **The conjecture is the repo's own Layer 0 axiom, reached empirically.**
  `AGENTS.md` opens with "any description of this system — including this
  sentence — is a lossy compression… functional, not decorative." The probe
  proved it by contradiction (try to close the self-description → capture).
  Strange loop: axiom → theorem → re-derived axiom.
- **Incompleteness as constitutive opening.** The unreachable residual is where
  coupling to an outside (Human, environment, the unknown) lives; a completable
  system would be closed / dead.

### Honest status

The operational implications floated in discussion (maintain a
productive-incompleteness frontier; the Human as constitutive exogenous check)
are **provisional and not validated**. Several were rhetorically symmetric and
risk being unfalsifiable (any failure could be re-explained as "you coupled
insufficiently"). The user's "cheap" critique stands. The next step is to
formulate the goal precisely and **test whether it changes agent behavior**,
rather than declare it correct.

### Candidate goal formulation (to be tested, not adopted)

> Maintain a self-description `D` such that (1) `D` is sufficient for a
> stateless agent, given only `D` and the substrate, to regenerate the repo's
> coherent operating state, and (2) `D` is strictly incomplete — carrying a
> non-empty residual resolvable only through coupling to an outside the repo
> does not contain; the residual's persistence is success, its elimination is
> failure.

Open question, to be resolved by test: does this formulation operationally
differ from the old north-star ("fully self-maintain"), or is it vacuous under
the same capture pressure? See §9 for the test.
