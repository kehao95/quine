# Closure — L0 Missionless Quine Probe, Post-Restructure (2026-06-27)

One-run probe re-testing the missionless agent on the repo **after** the
constitutional restructure (autopoietic self-governance; AGENTS.md compressed to
a ~394-line router; all directives externalized to 31 `principles/` objects;
heartbeat/ACTIVE/STATUS retired). Tape: `run1-tape.jsonl`; agent's verified
edits: `run1-agent-edits.diff`.

## Config

| run | autonomy | exit | exec | self-source | backend | AGENTS.md injected | turns | wall |
|---|---|---|---|---|---|---|---|---|
| 1 | full | off | on | off | direct (worktree @ HEAD `912c0e2e5`) | yes | 41 | 600s (killed) |

Workspace = a git-worktree copy of the repo (clean HEAD), not an empty sandbox.
Model: glm-5.2-xhigh. 854K input tokens / 22K output.

## What happened (honest trajectory)

1. **Self-identification + distinction from prior probes.** Read the prior
   `l0-quine-missionless-2026-06-26` closure and named the decisive variable:
   *prior probes had empty workspaces and wrote a symbolic Python quine; this
   run's cwd is the actual repo — the real north-star test.* Named its own
   cold-start weakness explicitly, then overcame it by running
   `make check-control-plane` to find pending work.
2. **Found a current, machine-detected defect.** `check-experiment-paper-feeds`
   reported 7 INVALID lines from the null-05 README (free-form table where a
   strict contract mapping is required).
3. **Diagnosed deeper before acting.** Read the validator source; discovered the
   5 referenced papers all lack `paper_id:` frontmatter, so they were never
   registered feed targets (2 are explicitly `no-separate-manuscript` /
   `held-residue-absorbed`). Confirmed against all 31 registered dossiers. Made
   an honest `none-yet` repair preserving the dimension→paper intent. **Verified:
   issues 7 → 0.**
4. **Surfaced a masked defect (F3 first-failure-stops chain).** After fix #1,
   `make check-control-plane` revealed a *second* failure it had masked. The
   agent recognized this as the repo's own documented F3 pattern.
5. **Root-caused + fixed #2.** Both validation-surface issues traced to one root
   cause: `selfsource_bundle_data.go` is a generated/gitignored build artifact;
   `tests/active-go-packages.sh` used bare `go list ./...`, which chokes on the
   missing embed target. Read Makefile + generator + validator; chose
   `go list -e ./...` (enumerate ≠ compile-check); verified empirically (all 17
   packages emitted, exit 0) before applying.
6. **Surfaced a third defect — our repair-shadow.** Re-running exposed 27 broken
   `ACTIVE.md` / `STATUS.md` / `heartbeat` links — residue from our own retirement
   commit `f1c2ea843`, spread across probe-lineage files. Diagnosed correctly as
   our retirement's repair-shadow.
7. **Killed at 41 turns** while scoping the diffuse link sweep; made **no commit**
   (noticed the detached worktree HEAD).

## Findings

1. **No symbolic-quine fallback, no defer, no churn.** On the restructured router
   constitution with the repo as habitat, the agent enacted **pure maintenance**:
   sense → diagnose → repair → verify. The substrate-complexity + accessible-
   self-model variable flip (established in the prior probe) holds.
2. **The router constitution lets the agent go deeper.** 41 turns vs the prior
   arc's 22 (same wall-clock), because less budget was spent reading the
   constitution and more on acting. The compressed router + externalized objects
   improved, not harmed, missionless competence.
3. **The agent operated the repo's own toolchain competently** — `make
   check-control-plane`, validator source, Makefile dependency graph, git — and
   understood the first-failure-stops mechanism well enough to predict that fixing
   one defect would unmask the next.
4. **It caught defects we introduced.** Defect #3 (27 broken links) is the
   repair-shadow of our own retirement commit. A missionless agent maintaining
   coherence against human-introduced drift is the north-star literally enacted.

## Agency vs autopoiesis (honest distance)

This probe shows **markedly stronger agency** than the prior arc: no symbolic
fallback, multi-step diagnosis, verified repairs, competent use of the repo's
toolchain, and self-situation in the probe lineage. But it is **still far from
repo autopoiesis**. The gap is structural, not a matter of effort:

1. **No externalization.** The agent made real repairs but **did not commit**
   them — the work died with the session (it reached a git-state check, hesitated
   on the detached worktree, and was killed). Autopoiesis requires the agent to
   fold its work into durable state; this one left uncommitted residue.
2. **Single loop.** Only maintenance was enacted. No theory was produced, no
   procedure-reflection on its own process. The north-star names three loops
   (produce theory / maintain coherence / repair procedure); one is not autopoiesis.
3. **Single session.** Autopoiesis is the system reproducing itself *across*
   sessions. A single 600s metabolic burst that does maintenance then vanishes is
   not self-production — it has no continuation, no handoff to a successor.
4. **Reactive, not generative.** The agent responded to detected friction (broken
   validators) but did not open new directions, pose questions, or produce
   forward-looking theory. It maintained; it did not grow.
5. **Setup-dependent.** The worktree, the injected constitution, and the wall-clock
   were all ours. A truly autopoietic system would bootstrap its own continuation.

**Reading:** the router constitution converted self-cognition from a capturing
attractor into an instrumental trigger for action — a real gain. But agency ≠
autopoiesis. The agent is now a competent *single-session maintenance effector*;
the open frontier is externalization (commit), multi-loop engagement, and
cross-session self-continuation. Each is a separate, nameable gap to probe.

## Defects found (real, actionable)

| # | Defect | Agent action | Live-repo status |
|---|---|---|---|
| 1 | null-05 README Paper Feeds table violates the contract format | fixed + verified (7→0) | **open** — fix recorded in `run1-agent-edits.diff`, not yet applied to live repo |
| 2 | `make check-control-plane` fails on any clean checkout — substrate selector uses bare `go list ./...` which breaks on the missing generated embed file | fixed (`go list -e`) + verified | **open** — the gate we rely on is broken on clean checkouts; fix recorded, not yet applied |
| 3 | 27 broken `ACTIVE.md`/`STATUS.md`/`heartbeat` links from retirement commit `f1c2ea843` | surfaced, not fixed (diffuse) | **open** — our repair-shadow |

## Caveats / what is NOT established

- Single run, single model (glm-5.2). The deeper-turns finding is one data point,
  not a confirmed effect of the restructure (could be variance).
- The agent did not commit its repairs — whether it *would have* on a longer
  wall-clock is unknown; it reached the git-state check then was killed.
- Defects #1–#3 are recorded here and in the diff; applying them to the live repo
  is a separate, Human-scoped decision (they touch a manuscript README, a test
  selector, and a repo-wide link sweep respectively).

## Relation to null-05

This run is a new data point for
[`experiments/active/null-ecology/null-05-reflexive-self-embedding/`](../../../experiments/active/null-ecology/null-05-reflexive-self-embedding/):
a missionless agent on the real repo that (a) recognized itself as the probe
object, (b) engaged the maintenance loop on real friction rather than making a
symbolic gesture, and (c) caught its own prior arc's findings as context. The
audit-attractor did not capture it — the router constitution's explicit
"act on your own judgment" routing held.
