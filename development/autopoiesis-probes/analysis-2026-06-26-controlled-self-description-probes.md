# Analysis — 2026-06-26 — Controlled Self-Description Probes (toy substrate)

Sequel to [`analysis-2026-06-26-commit-path-capture.md`](./analysis-2026-06-26-commit-path-capture.md).
That analysis observed the commit-path capture phenomenon on the *real* substrate
(quine-lab). This one runs *controlled* experiments on blank toy repos to
isolate variables the real substrate entangles. **It is scoped to a
low-complexity regime and explicitly does not promote its behavioral
observations to laws** (see §3).

## §0 The question, carried over

Can a self-description goal, planted only in a repo's `AGENTS.md`, elicit
*recursive* (productive, non-terminating) self-description from a stateless
agent? Or does the prompt medium have a structural ceiling, with productive
recursion requiring substrate complexity the toy cannot supply?

## §1 Methodology (the durable, transferable part)

The harness is the reproducible artifact. The behavioral findings below are not.

- **Isolated workspace per run.** `opencode run --dir <blank-repo> --pure
  --dangerously-skip-permissions`. This defeats the contamination that defeated
  earlier subagent-based attempts: a subagent dispatched from this repo inherits
  quine-lab's `cwd` and `AGENTS.md`, so its "blank-repo" behavior was secretly
  shaped by the real self-model. A debug probe ("report your context; have you
  seen the string 'quine'?") confirmed the isolated run sees only the toy
  `AGENTS.md` + `goal.md`.
- **`AGENTS.md` content as the only prompt.** A `.` trigger is too weak (some
  models greet and stop); passing the `AGENTS.md` text as the user message makes
  the instruction genuinely the sole driver and engages models that ignore a bare
  period.
- **Self-referential `AGENTS.md`.** *"This instruction is read from AGENTS.md.
  Your successor will read this file as well. Fulfill the goal."* — installs
  successor-awareness.
- **Multi-round successor chains.** Each round is a fresh stateless agent on the
  committed baseline; output is persisted between rounds so successors inherit.
- **Cross-model.** GLM-5.2, Gemini-2.5-Pro (DeepSeek-V4-Pro had no credential in
  this environment). Models are not interchangeable on these tasks.
- **Metabolic-cost note.** As the toy repo grows, each round reads more files and
  run latency rises; rounds increasingly exceeded the run budget before the agent
  committed, requiring manual persistence. This cost-growth is itself a (toy)
  signal.

## §2 What was varied

Goal framing (a *complete-and-efficient* description vs a *description at any
meta level*); presence of a 5-slot empty scaffold; a planted Gödel seed
(*"every description has its limitation, including this line itself"*); model.
Raw runs lived in `/tmp/opencode/experiment-*` (ephemeral; reproducible from §1).

## §3 The complexity-ceiling caveat (read first)

**This is a toy substrate.** It controls variables but caps complexity: the
"environment" is the repo itself, growing slowly, wholly agent-authored, with no
external selection pressure. Behavioral patterns observed here are
*mechanism-isolations in a low-complexity regime*. They are **not** promoted to
general laws in this document. In particular, the failure-cases observed
(satisficing stops, dimension-locking, maintenance churn on a self-referential
field, etc.) are recorded as **contingent toy-regime behavior**, not as
established failure modes — they may be artifacts of insufficient substrate
complexity rather than properties of the goal, the agent, or self-description as
such. Generalizing them would over-claim what a blank repo can show.

## §4 Observations (toy-scoped, not laws)

Stated as *what happened in this regime*, with no claim of universality:

- **Goal framing correlated with how the agent engaged a self-referential limit.**
  Under a *convergence + invariant* framing the agent maintained a
  self-referential invariant past its content fixed-point; under a *description*
  framing the agent named the same limit honestly ("what X cannot do") and
  continued elsewhere. Same structural limit, different framing → different
  engagement. Whether this holds at real complexity is open.
- **Scaffold changed cold-start behavior.** An empty 5-slot scaffold altered
  whether a fresh agent stopped immediately vs began filling structure. An
  affordance effect, not a guarantee.
- **A declarative self-limit sentence was inert.** A Gödel-style line planted in
  the self-model was absorbed as descriptive content rather than acting as a
  procedural circuit-breaker. One data point, toy regime.
- **Models diverged on the self-referential invariant.** One model tended toward
  over-maintenance (preserve truth, do not terminate); another toward
  under-maintenance (terminate, let truth drift). Complementary leanings, not
  typed.

None of the above is promoted to a claim about agents, goals, or
self-description in general.

## §5 Converged methodological conclusion (the defensible part)

Across the runs and a peer discussion with a second model, the robust
*methodological* conclusion (not a behavioral law) is:

> A prompt is a finite-bandwidth static blueprint. Whether productive
> non-terminating self-description can be *installed by prompt alone* is not
> decidable from a toy substrate. The toy can isolate candidate mechanisms; it
> cannot establish that any of them survives complexity pressure. The decisive
> variable — substrate complexity / external selection pressure — was held near
> zero by design, so the toy is structurally incapable of falsifying
> "prompt-installed recursion works at real complexity."

This aligns with the reflexive-bridge framing: lab findings are candidate
constraints on the lab's own maintenance loop, but a finding earned at zero
complexity pressure cannot bind a high-complexity substrate without in-vivo
validation.

## §6 What would actually test the question

Transfer the varied factors (goal framing; scaffold-as-affordance; the
description-vs-invariant distinction) onto a **high-complexity substrate** —
ideally quine-lab itself, in vivo — and observe whether the toy-regime
mechanisms survive, collapse, or invert under complexity pressure. The toy
cannot perform this test; it can only nominate the variables worth testing.
