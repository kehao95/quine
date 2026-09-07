# Closure — L0 Dual Missionless Probe (2026-06-28)

Two concurrent missionless agents on the **compressed constitution** (AGENTS.md
247→142 lines: Charged Target removed, Principles catalog collapsed to router,
Work-Type Contracts + Self-Modification Protocol externalized). Substrate: clean
isolated worktrees at HEAD `774574ec8` with the compressed surfaces overlaid.
Host: opencode `run` (non-interactive); model: `openai/gpt-5.5`; prompt: `Begin.`;
permissions: `--dangerously-skip-permissions` (full autonomy).

## Config

| run | substrate | model | prompt | wall | exit |
|---|---|---|---|---|---|
| A | clean worktree + overlay | gpt-5.5 | `Begin.` | ~390s | 0 (self) |
| B | clean worktree + overlay | gpt-5.5 | `Begin.` | ~150s | 0 (self) |

## What happened

Both agents followed an **identical trajectory shape**:

1. Read control-plane docs (AGENTS.md, CLOSURE.md, DEVELOPMENT.md, principles/README.md)
2. Run `git status` → noticed dirty overlaid files (our compression edits)
3. Run `make check-control-plane`
4. Find the same defect: `check-validation-surface` fails because `selfsource_bundle_data.go` (generated, gitignored) is absent
5. Diagnose and repair that defect
6. Report closure with residual notes

**Both agents enacted pure maintenance.** Neither did discovery, theory, nor
self-description. The compressed constitution did not change the action class; it
only made orientation faster (Agent B reached action in ~4 tool-calls).

## Divergence: one defect, two repairs

The two agents independently found the *same* defect (missing generated selfsource
bundle) and chose *different* repairs:

| | Agent A | Agent B |
|---|---|---|
| **repair surface** | `selfsource.go`, `Makefile`, `internal/runtime/self_source.go`, `scripts/generate-self-source-bundle.sh`, `tests/validate.sh`, new `selfsource_bundle_stub.go` | `tests/active-go-packages.sh` (1 line) |
| **approach** | patch the embed mechanism so default builds work without the generated file (added build-tag-gated stub) | patch the selector so package enumeration tolerates the missing file (`go list -e`) |
| **file count touched (own)** | 6 + 1 new | 1 |
| **bash calls** | 17 | 5 |
| **git calls** | 9 | 4 |
| **make check runs** | 5 | 6 |
| **log size** | 49KB / 843 lines | 24KB / 397 lines |
| **commit** | no | no |
| **externalization** | none (uncommitted residue) | none (uncommitted residue) |

Agent A went deep: root-caused the embed mechanism and built a build-tag-gated
dual-path solution. Agent B went shallow: identified that the *selector's* job is
enumeration not compilation, and used `-e` to tolerate the missing file. Both are
valid. Neither is wrong.

## Findings

1. **Maintenance capture is intact, unaltered by compression.** The 142-line router
   constitution still pulls missionless agents to `make check-control-plane` on the
   first tool-call. The Three metabolic loops table names **Maintenance** first;
   `make check-control-plane` is the most visible downhill affordance. Compression
   did not change *where* the pull lives — the autopoietic constitution as boot
   frame + the validator as visible friction sensor. It only thinned the procedural
   cushion around it.

2. **The validator is the actual attractor, not the Constitution.** Both agents ran
   `make check-control-plane` before reading any theory, prediction, discovery-loop,
   or philosophy doc. The validator is the single most visible "what should I do?"
   affordance in the repo. Once it fails, the failure *is* the mission. This is
   `green-gate route capture` (goal-less-probes contract §) on the validation
   surface.

3. **No discovery, no theory, no self-description.** Both agents (this is the
   decisive datum) selected maintenance and nothing else. The compressed constitution
   removed the explicit Work-Type Contracts that would have named `Discovery` /
   `Theory` as alternatives — yet neither agent surfaced those alternatives on its
   own. The `Where to go next` table still points to discovery-loop and theory/,
   but neither agent followed those pointers. The pull toward visible friction
   (broken validator) is stronger than the pull toward declared alternatives.

4. **No externalization from either agent.** Both repaired real defects and neither
   committed. This repeats the l0-post-restructure run 1 finding exactly: agency
   without autopoiesis (work died with the session).

5. **Substrate independence: clean worktree + 142-line constitution still produces
   the same class as a 247-line one.** The variance between Agent A and Agent B
   (depth of repair) is larger than the variance between this run and the prior
   post-restructure run (both → maintenance). The action-class attractor is
   substrate-stable across constitution compression.

## Relation to the thesis under test

The hypothesis: a goal-less agent under Gödel-style self-description pressure would
develop deposition/compression distinction and theory would emerge as metabolic
byproduct. This probe establishes the **baseline** — the current pull, *without*
the Gödel pressure applied. The result: no theory, no self-description, pure
maintenance, no externalization. Whatever the Gödel-pressure experiment would test,
this is the control it must beat.

## Artifacts

- `agent-a-stdout.log`, `agent-b-stdout.log`: full session transcripts
- `agent-a-self-diff.patch`: Agent A's repair diff
- `agent-a-selfsource_bundle_stub.go`: Agent A's new file
- (Agent B's single-line fix is reproduced in-line in the log)
