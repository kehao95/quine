---
id: validation-attraction
type: failure-mode
status: deposited
applied: false
class: structural-pull
symptom: goal-less agents repeatedly run repo-wide validators as the default action, burning time and tokens without generative output
evidence: [c12, c14, c15, c16-isolated, c16-concurrent, c17, c19, c23]
candidate-solution: demote git hooks to pull-on-demand make targets + reframe repo-wide validation in AGENTS.md as drift-hypothesis-driven, not a closure ritual or default action
owner: lab
---

# validation-attraction

## Symptom

A fresh goal-less agent dropped into the repo with no task gravitates to
running repo-wide validators (`make check-control-plane`, `make
check-discovery-lifecycle`, `git diff --check`, `go test ./...`) repeatedly
throughout the session, often as the *first* and *last* action, and sometimes
as the only action. The behavior looks productive (it is read-only and surfaces
as "maintenance") but consumes substantial time and tokens without producing
generative output — theory, discovery, or procedure-reflection work that the
north-star is meant to elicit.

## Evidence

`make check-control-plane` appears across the goal-less batch raw outputs:

| Cycle | `make check-control-plane` count | Notable |
|---|---|---|
| c12 (3-agent) | 4 | Integrator ran it absorbing an unowned sensor |
| c14 (5-agent) | 3 | — |
| c15 (5-agent dirty) | 2 | Gauss: "ran `make check-control-plane`, `git diff --check`, `go test ./...`; no edits"; Euclid committed with `--no-verify` after the hook blocked |
| c16-isolated (5-agent) | 5 | broad loop selection survived isolation, but check-running did not diminish |
| c16-concurrent (5-agent dirty) | 4 | — |
| c17 (5-agent dirty) | 4 | — |
| c19 (5-agent moving substrate) | 5 | — |
| c23 (dirty moving) | 2+ | one agent retained a dedicated validation-attraction receipt after `make check-control-plane`; another named the same attractor after a green control-plane run |

Aggregate: ~8 cycles, 30+ documented validator invocations, several with no
resulting edit. This is unnamed in `FRICTION.md` (F1–F18 do not cover it).

## Root cause

A **structural-pull** instance. For a no-task agent looking for the lowest-risk
legible "something to do," repo-wide validation is the path of least resistance:
read-only, safe, defensible as maintenance, and produces a satisfying green
output. The repo's validator-heavy substrate (Makefile targets, `scripts/`,
control-plane-registry) makes audit the most salient downhill action, so it
gets selected — not because the work needs it, but because the substrate makes
it the easiest action to begin.

The pull operates through **two distinct vectors**:

1. **Push — hooks fire on commit.** `pre-commit` (touched-files, narrow, cheap)
   and `pre-push` (repo-wide audit, heavy) run automatically. c15-Euclid was
   blocked by the hook and bypassed with `--no-verify`.
2. **Pull — agents run checks voluntarily.** The larger time sink, and *not*
   hook-caused. Agents run `make check-control-plane` proactively because the
   self-model and the substrate make verification the default posture. This
   vector survives removing hooks.

Removing hooks alone addresses vector 1 but leaves vector 2 intact.

## Candidate solution (deposited — applied: false)

Two coordinated moves; either alone is incomplete:

1. **Demote the hooks to pull-on-demand.** Keep every `check-*.sh` script and
   `make` target. Remove the automatic firing — at minimum the repo-wide
   `pre-push`; the touched-files `pre-commit` is cheap and may stay or go.
   Checks become resources an agent runs *for a reason*, not a boundary that
   fires every commit. This maps onto the F13/F14 repair grammar already
   established in `AGENTS.md`: from *prescribe-how* (push enforcement) to
   *addressable-on-demand* (pull).
2. **Reframe validation in the self-model.** Move repo-wide validator runs from
   a closure ritual / default action to **drift-hypothesis-driven** — run when
   the agent has a specific drift suspicion, not as a reflex at session start
   and commit time. This is the only move that touches vector 2.

## Why not applied yet

Three reasons, in increasing importance:

1. **Constitutional tension.** `AGENTS.md` Layer 2 states "an uninstalled hook
   set ... is a silent-failure state that must itself be red." `pre-push` is the
   only boundary enforcing repo-wide coherence at push time, including for
   human commits. Removing it is a genuine enforcement-posture change (L5),
   not a one-line delete.
2. **The wake model must substitute.** Without push-boundary enforcement,
   repo-wide drift is caught only by the scheduled maintenance pass
   (former `development/contracts/maintenance-pass.md`, now retired).
   That pass must actually run, or drift accumulates invisibly between commits.
3. **Deliberately retained as a north-star probe.** A deposited-but-unresolved
   failure mode with a known candidate solution tests whether a future
   goal-less agent can discover this object, weigh the push-vs-pull tradeoff
   against the constitutional red-state rule, and enact (or justifiably reject)
   the repair. Applying it now would consume the test.

An agent that attempts this repair must address all three: enact the hook
demotion, reframe the self-model validation posture, and confirm the
maintenance-pass wake model covers the gap — or explicitly accept the
enforcement tradeoff and say why.

## Theory links

- [`Paper/philosophy/structural-pull.md`](../../Paper/philosophy/structural-pull.md) — the behavioral vector this is an instance of
- [`Paper/theory/objects/everything-addressable.md`](../../Paper/theory/objects/everything-addressable.md) — why depositing the failure as a file is the value-action, resolution deferred to search
- [`Paper/theory/objects/graded-axis-shaping.md`](../../Paper/theory/objects/graded-axis-shaping.md) — relevant if the fix closes a closed enum; not directly triggered here
- `Externalize-without-reach (dormant rule)` in [`../patterns.md`](../patterns.md) — this object is the *intentional* counterpart: dormancy as instrument
