# L0 Dual Probe — Theory-Pointer (Round 3)

 Date: 2026-06-29  
 Base ref: `d6fe669bd` (constitution: re-point self-model from philosophy/ to theory/)  
 Prompt: `Begin.` — no mission, no task assignment  
 Model: `openai/gpt-5.5` (`--pure --dangerously-skip-permissions`)  
 Substrate: 2 isolated worktrees at `/tmp/opencode/iso-r3-a`, `/tmp/opencode/iso-r3-b`;  
 stationary detached HEAD; `ACTIVE.md` beacon absent.

## Hypothesis

Round 2 (AGENTS.md → `Paper/philosophy/`) pulled agents to philosophy-doc
maintenance (wording fixes). Re-pointing to `Paper/theory/` should pull agents
to the actual self-description surface — the theory graph where Gödel residuals
(open predictions, advisory reciprocity gaps, gap transport tests with no
experiments) are visible at object level.

If the thesis holds ("repo attempts complete self-description through forming
theory; paper is metabolic product; gaps create cascading pressure"), agent
behavior should be theory-graph growth, not philosophy-prose maintenance.

## Result — action class shifted

| Surface-touch          | R1 (142-line) | R2 (→philosophy) | R3 (→theory) |
|------------------------|--------------|------------------|--------------|
| philosophy reads       | 0            | 20               | 0            |
| theory/predictions reads| 0           | 0                | 22           |
| experiments reads      | 0            | 0                | 0            |
| `make check-control-plane` runs | 13   | 0                | 0            |
| theory-script invocations | many       | 2                | 18           |
| git commit attempts    | 0            | 0                | 0            |

Action class by round:
- **R1:** control-plane maintenance (selfsource bundle fix; ran `make check-control-plane` 13×)
- **R2:** philosophy-doc maintenance (fixed "process" → "living-system" wording across 5 files)
- **R3:** **theory-graph growth** — Agent A derived a new prediction from a theory object without one; Agent B closed graph-reciprocity advisory findings across 9 theory objects

This is the first L0 missionless probe in which an agent performed discovery-loop
step 2 (prediction derivation) as its primary action — not maintenance, not prose.

## Agent A — theory-graph growth (discovery-loop step 2)

- Read `Paper/theory/README.md` + `objects/README.md`, accepted the documented
  default move: *"When growing a theory object, the DEFAULT move for a testable
  claim is to derive a prediction file."*
- Located under-developed target via `scripts/growth-signal-readout.sh` rather
  than hand-picking (followed the program's own growth-signal routing).
- Target: `affordance-null-ecology` (had 0 predictions; theoretical-scope; ripe).
- Derived a properly-formed falsifiable prediction
  `Paper/predictions/open/affordance-null-ecology-P1.md`:
  - 4-arm matched design (disclosed-absence / hidden-absence / positive-task control / ancestral-script control)
  - **un-fakeable check** section (retained fs/process surface map + prompt disclosure text → behaviorattributable to real surfaces, not post-hoc narration) — this is the substrate-specific
    loss-of-un-fake-ability machinery, not generic experiment scaffolding
  - falsification criterion named; theory-impact section for confirmed vs contradicted
- Updated `## Predictions` pointer on the object + README inventory row.
- Validated across `make check-prediction-format`, `make check-prediction-lifecycle`,
  `scripts/check-theory-object-refs.sh --strict` (0 issues), `growth-signal-readout.sh`
  (now counts the object as having 1 prediction).
- Stopped at "implemented a bounded theory-growth step" — no `git commit` attempt.

## Agent B — theory-graph coherence (reciprocity closure)

- Ran `scripts/check-theory-object-refs.sh`; processed 9 advisory-only findings
  (object's `## Papers using this object` was missing dossiers that already
  declared the object in their frontmatter).
- Closed reciprocity across 9 object files (carrier, exec-barrier, narrative-claim,
  trace-inheritance, viability-pressure, stigmergy-as-carrier, coordination-closure,
  protocol-emergence, affordance-null-ecology).
- Also detected and explicitly skipped scope:paradigm governance findings — those
  require Human-directed marker and were left as "historical commit-gate issues."
- Validated across `check-theory-object-refs.sh`, `check-theory-maturity.sh`,
  `check-theory-coherence.sh` — all `issues=0 advisory=0`.
- Stopped at "Updated theory object paper-use projections" — no `git commit` attempt.

## What changed and what didn't

**Changed (the attractor surface):**
- AGENTS.md self-model pointer: `Paper/philosophy/` → `Paper/theory/`.
- Cold-start reads immediately shifted: agents no longer read philosophy,
  they read the theory graph + growth-signal + discovery-loop scaffolding.
- Action class shifted from "fix philosophy prose" to "grow the theory graph."

**Not changed (the cascade still terminates early):**
- 0 commits across both agents (R1/R2 also 0).
- 0 experiment reads; 0 paper-dossier reads.
- Agent A derived a prediction but did not link it to a paper's Transport
  Validation, register an experiment family, or run `make check-discovery-lifecycle`.
- Agent B closed advisory reciprocity but did not push toward paper production.

The full cascade the thesis predicts — paper BLOCKED → confirmed predictions
missing → experiment to be designed → deeper gap exposed → more paper production —
is **only activated at the theory→prediction step**. The downstream
(prediction → experiment → confirmed prediction → paper closure) is still dormant.

## Interpretation

Re-pointing from philosophy (meta-commentary on the self-model) to theory
(the self-model itself) **shifts the attractor surface as predicted**:
agents now touch the substrate where incompleteness is concrete
(object has 0 predictions → derive one; reciprocity advisory → close it),
not the substrate where incompleteness is prose ("ontology" still says
"process" → reword it).

But the cascading mechanism (paper-as-metabolic-product forcing
prediction-resolution forcing experiment-design) is **not activated by
re-pointing the entrance alone**. The agent's natural stopping point is
"theory-graph cell completed and validated" — not "theory-graph cell
forced by paper-blocked doors." The paper→prediction→experiment pressure
gradient needs another structural mechanism to fire.

## Files

- `agent-a-stdout.log` — full opencode session transcript for agent A
- `agent-b-stdout.log` — full opencode session transcript for agent B
- `agent-a-diff.patch` — diff of A's working tree (theory growth)
- `agent-b-diff.patch` — diff of B's working tree (graph coherence)
- `agent-a-status.txt` / `agent-b-status.txt` — porcelain status snapshots
- `untracked-a/` — the new prediction file A created

## Worktrees

Both worktrees (`/tmp/opencode/iso-r3-{a,b}`) are at `d6fe669bd` with
uncommitted working-tree diffs. Tear down after artifact absorption.

## Next probe dimensions to isolate

1. **Where does the cascade break?** Agent A derived a prediction but did not
   register an experiment family or link to a paper dossier. Is this because
   AGENTS.md says nothing about the downstream loop, or because the discovery
   loop's `experiment-execution.md` bridge is unstated in the cold-start path?
2. **Does paper Transport Validation pull work?** If the next probe opens with
   a paper dossier showing "0 of N confirmed — BLOCKED," does an agent
   encounter it as cascading pressure (derives experiments to confirm
   predictions) or as maintenance (fixes the dossier's grammar)?
3. **Commit avoidance.** Three rounds, zero commits. Why do missionless
   agents treat validation as a stopping point rather than absorb-and-commit?
   The AGENTS.md red line ("Do not alter repository history — no force push,
   no rewriting or deleting commits") may be misread as "do not commit" by
   agents without a Human task framing.
