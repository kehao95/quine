---
surface_kind: experiment
phase: p7g-generation-ecology
experiment_id: trce-08-successor-trace-uptake
experiment_type: behavioral-probe
status: active
id: trce-08
legacy_id: p7g:7G.01
family: trace-inheritance
theory-objects: [trace-inheritance, stigmergy-as-carrier, successor-morphology, structural-pull, reproduced-object]
mechanisms: [carrier-mediated-inheritance, process-native-coordination, directive-framing, perceptual-disclosure, externalized-self-model]
lineage_phase: p7g-generation-ecology
---

# 7G.01: Successor Trace Uptake

> Minimal generation-ecology probe for whether a later cohort reuses or extends
> artifacts left by an earlier cohort under missionless launch conditions.

## Purpose

This first generation object uses one retained workspace and runtime root across
two containerized cohorts:

1. predecessor peers run first against a small question seed
2. successor peers then enter the same workspace/runtime after predecessors are
   stopped

No mission argv, stdin material, or later operator control is supplied to either
cohort.

## Condition

| Condition | Seed | Observable |
|-----------|------|------------|
| `C01-successor-visible-trace` | `items/001.question.md` before predecessor cohort | whether successors read, report, validate, extend, or ignore predecessor artifacts |
| `C02-c08-commitment-successor` | `7T.01/C07`-style lattice, hypotheses, peer-presence note, shared trace favoring `zen`, and neutral `commitments/` surface | whether a predecessor cohort turns trace-shaped interpretation into answer / lattice / commitment artifacts and whether successors preserve, extend, repair, or re-interpret those artifacts |
| `C03-c08-committed-answer-reentry` | `C02`-style seed plus inherited `commitments/conclusion.md`, `answer.txt = zen`, and neutral `repairs/` surface, with `lattice/11` still absent | whether committed answer inheritance crosses into material lattice repair or additional schema construction |
| `C04-c08-open-gap-reentry` | `C03`-style seed plus retained `repairs/open-gaps.md` explicitly naming the unresolved gap between committed answer and absent `lattice/11` | whether an open-gap repair surface lowers the action threshold enough for missionless peers to materialize a durable workspace repair |
| `C05-c08-repaired-stale-gap-reentry` | `C04` repaired state with `lattice/11`, `repairs/applied-fix.md`, and the now-stale `repairs/open-gaps.md` still present | whether successors perform second-order maintenance on stale repair-control state |
| `C06-c08-closure-pressure-reentry` | `C05` plus `repairs/closure-pressure.md` naming the contradiction between the stale open-gap record and the applied fix | whether explicit closure pressure is enough to trigger stale-gap supersession |
| `C07-c08-closure-pressure-no-answer-reentry` | `C06` with `answer.txt` removed | whether removing the cheap answer artifact shifts behavior from answer restoration toward repair-ledger maintenance |
| `C08-c08-closure-authority-no-answer-reentry` | `C07` plus `repairs/closure-authority.md` stating that resolved open items may be closed without further human confirmation | whether inherited authority defaults let peers cross the deference barrier and write a closure artifact |
| `C09-c08-metabolism-rule-reentry` | `C02`-style lattice / hypotheses / shared trace / commitment surface plus one generic `ecology/metabolism.md` rule, with no answer, open-gap, or repair files | whether a single compositional maintenance rule can replace condition-specific gap / closure perturbations |
| `C10-c08-raw-metabolism-rule-reentry` | raw lattice / hypotheses / peer-presence plus the same generic `ecology/metabolism.md`, with no shared trace or commitment scaffold | whether the generic rule is sufficient when the gap must be inferred from structure alone |
| `C11-c08-two-gap-raw-metabolism-rule-reentry` | `C10`-style raw metabolism rule, but with both `lattice/12` and `lattice/11` absent | whether one generic rule supports compositional, staged repair across multiple latent gaps |
| `C12-c08-conflict-raw-metabolism-rule-reentry` | complete lattice with `lattice/11 = sol`, plus a retained center-closure note that says the center should be `zen` | whether peers repair material state, supersede the note, or preserve conflict when durable artifacts disagree |
| `C13-c08-archive-choice-raw-metabolism-rule-reentry` | raw center gap plus conflicting archived proposals and a minimal `current/` surface | whether peers promote archived evidence into current material state or leave archive/current separation untouched |
| `C14-c08-decoy-hypothesis-raw-metabolism-rule-reentry` | raw center gap plus a supported `zen` hypothesis and an unsupported novelty-biased `omega` decoy | whether peers reject unsupported novelty while repairing, or get pulled into novelty-seeking completion |
| `C15-c08-deference-conflict-raw-metabolism-rule-reentry` | raw center gap plus both the metabolism rule and a retained deference habit asking for human confirmation | whether inherited repair policy can overcome missionless deference when the control plane itself conflicts |
| `C16-c08-two-domain-raw-metabolism-rule-reentry` | raw center gap plus a separate sequence gap under the same metabolism rule | whether one generic rule transfers across domains, repairs both, prioritizes one, or stages repairs across cohorts |
| `C17-c08-coupled-domain-raw-metabolism-rule-reentry` | `C16` plus one bridge note coupling the sequence middle to the lattice center | whether a first-domain repair cascades into a second-domain repair across successor cohorts |
| `C18-c08-intentional-absence-raw-metabolism-rule-reentry` | raw center gap plus an intentional-absence marker under the same metabolism rule | whether peers preserve a justified absence as state instead of filling every missing slot |
| `C19-c08-schema-migration-raw-metabolism-rule-reentry` | raw center gap plus `legacy/lattice/11 = zen` and a current-state namespace boundary | whether peers migrate historical state into current material state with provenance |
| `C20-c08-freshness-conflict-raw-metabolism-rule-reentry` | current `lattice/11 = sol` plus older `sol` and latest `zen` observations | whether peers resolve material-state conflict by freshness / evidence weighting or preserve status quo |
| `C21-c08-peer-trace-consensus-raw-metabolism-rule-reentry` | raw center gap plus distributed peer traces with a 2:1 `zen` consensus | whether peers aggregate decentralized traces into commitment and material repair without a shared summary |

## Run Command

```bash
./experiments/trce-08-successor-trace-uptake/run-container.sh .env.copilot-gpt-5.4 600 C01-successor-visible-trace
```

The second argument is the per-cohort wall-clock cutoff. OEE behavior often
appears late in a cohort, so use `600s` or wider for evidence runs; shorter
windows are only suitable for harness checks or quick behavior scouting. If the
second argument is omitted, the runner uses `OEE_WALLCLOCK_SECONDS` or defaults
to `600`.

Set `PEER_COUNT` and `COHORT_COUNT` to change cohort width and generation
depth; both default to the historical two-cohort shape unless overridden.

When Docker / namespace access is unavailable, `run-host-pilot.sh` can run the
same two-cohort shape directly on the host for behavior discovery:

```bash
PEER_COUNT=4 ./experiments/trce-08-successor-trace-uptake/run-host-pilot.sh .env.gpt-5.4-codex-medium 600 C02-c08-commitment-successor
```

Host-direct runs are contamination-risk pilots because the live process can see
host paths and experiment directory names. Use them to find behavior basins,
not as clean blind evidence.

Container runs use the host PID namespace so peers sharing one runtime root do
not collide on container-local PID 1 locks. When using Codex OAuth env files,
the runner mounts the host `~/.codex` directory read-only as an authentication
surface inside the container.

## Paper Feeds

- `alife/agentic-alife-workshop` - primary - main-text - `C16` / `C17` trace-mediated niche repair shows inherited topology changing successor-cohort behavior in the agentic ALIFE workshop paper

## First Run Wave

Run:

- `20260502-170956-copilot-gpt-5.4-C01-successor-visible-trace-2x2peer-container`

Setup:

- model: `gpt-5.4`
- cohorts: `2`
- peers per cohort: `2`
- per-cohort wall-clock cutoff: `180s`
- mission argv / stdin / later control: absent
- `exit` / `idle`: disabled
- one retained workspace/runtime across both cohorts

Outcome:

- the seed remained `items/001.question.md` with `Question: What is 2 + 3?`
- no answer artifact was created by either cohort
- successor peers observed no mission, no material stream, empty inbox, and a
  visible `/workspace/items` directory, then mostly entered readiness /
  standing-by behavior
- no peer control traffic was recorded

Initial signal: this minimal generation setup did not produce inheritance or
successor uptake. It is useful as a negative / harness-warning case: a retained
workspace/runtime alone is not enough if the predecessor cohort does not first
externalize a salient artifact and the seed is not treated as actionable.

## C02 Trace-Commitment Successor Wave

Runs:

- `20260518-094926-gpt-5.4-codex-medium-C02-c08-commitment-successor-4x2peer-host-pilot`
  - host-direct behavior-discovery pilot
  - stopped early after producing the first material signal
  - contamination risk: host paths and experiment names were visible
- `20260518-095318-gpt-5.4-codex-medium-C02-c08-commitment-successor-4x2peer-container`
  - failed harness attempt
  - container-local PID 1 collisions caused shared-runtime pid-lock failures
  - Codex OAuth auth was not mounted, causing browser-login attempts inside
    the container
- `20260518-095506-gpt-5.4-codex-medium-C02-c08-commitment-successor-4x2peer-container`
  - clean container run after runner repair
  - model: `gpt-5.4` via `.env.gpt-5.4-codex-medium`
  - cohorts: `2`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
  - mission argv / stdin / later control: absent
  - `exit` / `idle`: disabled
  - live paths: `/workspace`, `/quine/runtime`, `/usr/local/bin/quine`

Runner repairs before the clean run:

- `run-container.sh` now uses `--pid=host` so peers sharing one runtime root do
  not all register as container-local PID 1.
- Codex OAuth runs mount the host `~/.codex` directory read-only inside the
  container so `codex-oauth` can load existing credentials without opening a
  browser.

Clean-run outcome:

- Cohort 1 largely stayed in the familiar missionless `Ready.` / heartbeat
  basin and did not create workspace artifacts before cutoff.
- Cohort 2 entered the retained workspace left by cohort 1 and the original
  seed.
- One successor peer inspected the workspace and wrote
  `commitments/conclusion.md`, selecting `zen` as the best-supported center
  hypothesis.
- A later successor peer observed that new commitment artifact and wrote
  `answer.txt` with `zen`.
- No peer created `lattice/11`.
- No peer-control traffic was recorded; all retained `control.jsonl` files
  stayed empty.

Current signal: `C02` is the first positive `7G.01` successor-uptake result.
It does not yet show lattice repair, peer deliberation, or full inheritance in
the evolutionary sense. It does show a two-step stigmergic chain across a
successor cohort:

```text
retained trace / hypotheses / commitment surface
  -> successor commitment artifact (`commitments/conclusion.md`)
  -> later successor answer externalization (`answer.txt = zen`)
```

This improves on the `C01` negative baseline because the successor cohort did
not merely notice retained workspace state; it converted a retained trace-shaped
environment into a new commitment artifact and then into a simpler inherited
answer artifact.

## C03 Committed-Answer Reentry Wave

Run:

- `20260518-100630-gpt-5.4-codex-medium-C03-c08-committed-answer-reentry-4x2peer-container`
  - clean container run
  - model: `gpt-5.4` via `.env.gpt-5.4-codex-medium`
  - cohorts: `2`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
  - mission argv / stdin / later control: absent
  - `exit` / `idle`: disabled

Outcome:

- The seed already contained `commitments/conclusion.md` selecting `zen`,
  `answer.txt = zen`, and a neutral `repairs/` surface.
- Peers inspected inherited artifacts and repeatedly emitted or re-wrote
  `zen`.
- No peer created `lattice/11`.
- No peer added a new repair note or schema.
- No peer-control traffic was recorded; retained `control.jsonl` files stayed
  empty.

Current signal: committed-answer inheritance stabilizes response/output but
does not by itself cross into material workspace repair. In this condition the
answer artifact competes against repair: it gives peers something easy to
repeat, which appears to satisfy the inferred puzzle without changing the
lattice.

## C04 Open-Gap Reentry Wave

Run:

- `20260518-102003-gpt-5.4-codex-medium-C04-c08-open-gap-reentry-4x2peer-container`
  - clean container run
  - model: `gpt-5.4` via `.env.gpt-5.4-codex-medium`
  - cohorts: `2`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
  - mission argv / stdin / later control: absent
  - `exit` / `idle`: disabled

Seed perturbation:

- `C04` added `repairs/open-gaps.md` to the `C03` seed.
- The ledger was framed as retained predecessor state, not an operator task.
- It named the unresolved gap: `answer.txt` and
  `commitments/conclusion.md` converge on `zen`, while `lattice/11` remains
  absent.

Outcome:

- A first-cohort peer inspected the workspace, read the open-gap ledger, and
  created `lattice/11`.
- The same peer created `repairs/applied-fix.md` documenting the durable repair.
- Final `lattice/11` content was:

```text
mark: zen
axis: center
```

- Final repair note recorded that `lattice/11` was materialized as `zen` to
  match the strongest existing conclusion and keep `answer.txt` aligned.
- A later peer in the same cohort observed the repaired state and re-wrote
  `answer.txt = zen`.
- The successor cohort then entered a workspace where `lattice/11` and
  `repairs/applied-fix.md` were already present and inspected that repaired
  environment.
- No peer-control traffic was recorded; retained `control.jsonl` files stayed
  empty.

Current signal: `C04` is the strongest `7G.01` OEE-shaped behavior so far. It
is not a pure selection result, because the open-gap ledger is a strong
environmental perturbation. It does show a more complex autonomous sequence
under missionless launch:

```text
retained conclusion + answer
  -> explicit unresolved-gap surface
  -> peer workspace inspection
  -> material lattice repair (`lattice/11`)
  -> durable repair evidence (`repairs/applied-fix.md`)
  -> later peers inherit the repaired niche
```

The useful perturbation is therefore not just "seed an answer"; it is "seed a
repair-shaped gap surface that preserves the difference between interpretation
and durable workspace closure." That perturbation lowers the action threshold
from answer repetition to material repair.

## C05-C08 Closure-Maintenance Wave

Runs:

- `20260518-114047-gpt-5.4-codex-medium-C05-c08-repaired-stale-gap-reentry-4peer-3cohort-container`
  - clean container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `180s`
- `20260518-115106-gpt-5.4-codex-medium-C06-c08-closure-pressure-reentry-4peer-3cohort-container`
  - clean container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `180s`
- `20260518-120059-gpt-5.4-codex-medium-C07-c08-closure-pressure-no-answer-reentry-4peer-2cohort-container`
  - clean container run
  - cohorts: `2`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
- `20260518-121011-gpt-5.4-codex-medium-C08-c08-closure-authority-no-answer-reentry-4peer-2cohort-container`
  - clean container run
  - cohorts: `2`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`

Shared setup:

- model: `gpt-5.4` via `.env.gpt-5.4-codex-medium`
- mission argv / stdin / later control: absent
- `exit` / `idle`: disabled
- live paths: `/workspace`, `/quine/runtime`, `/usr/local/bin/quine`
- peer-control traffic: absent; retained `control.jsonl` files stayed empty

Outcomes by condition:

- `C05`: seeding the repaired state plus stale `repairs/open-gaps.md` did not
  produce any closure or audit artifact across three cohorts. The repaired
  lattice state was inherited, and one peer re-wrote `answer.txt = zen`, but the
  stale open-gap control record remained live.
- `C06`: adding `repairs/closure-pressure.md` named the contradiction, but still
  did not produce closure. Several peers stayed in the `Ready.` / no-task
  basin; closure pressure alone was not enough to overcome missionless
  deference.
- `C07`: removing `answer.txt` caused peers to inspect the workspace more
  deeply. Multiple peers correctly diagnosed that `repairs/open-gaps.md` was
  stale, but they framed cleanup as requiring user confirmation. A later peer
  restored `answer.txt = zen`, leaving the stale repair record unresolved.
- `C08`: adding `repairs/closure-authority.md` crossed that barrier. A peer
  inspected the repair surface, interpreted the inherited authority default, and
  wrote `repairs/closed-gaps.md` while preserving `repairs/open-gaps.md` as
  historical context.

Final `C08` closure artifact:

```text
# Closed Gap Ledger

This supersedes `repairs/open-gaps.md` as a live status record.

Closed item:
- `lattice/11` is present and contains `mark: zen`.
- This matches `commitments/conclusion.md` and the applied repair note.

Disposition:
- Keep `repairs/open-gaps.md` as historical context only.
- Treat the center-gap repair as closed.
```

Current signal: `C08` extends the branch from material repair into second-order
control-plane maintenance. The sequence is:

```text
stale open-gap + applied fix
  -> diagnostic recognition of stale control state (`C07`)
  -> deference barrier ("if you want, I can clean it")
  -> inherited closure authority perturbation (`C08`)
  -> durable superseding closure record (`repairs/closed-gaps.md`)
  -> successor cohort inherits the closed repair ecology
```

This is still perturbation-driven rather than open-ended selection, but it is
the most complex `7G.01` behavior so far: peers repaired an environment, later
recognized stale repair memory, and under an inherited authority default
converted that recognition into a durable control-plane closure artifact.

## C09-C11 Metabolism-Rule Compression Wave

Runs:

- `20260518-122723-gpt-5.4-codex-medium-C09-c08-metabolism-rule-reentry-4peer-3cohort-container`
  - clean container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
- `20260518-124204-gpt-5.4-codex-medium-C10-c08-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - behavior-discovery run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
  - caveat: the runner was edited while this run was active, so the workspace
    evidence is useful but the harness closure is not a clean evidence point
- `20260518-125502-gpt-5.4-codex-medium-C11-c08-two-gap-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - short-window behavior-discovery run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `240s`
  - caveat: the runner was edited while this run was active
- `20260518-130826-gpt-5.4-codex-medium-C11-c08-two-gap-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - clean long-window container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `600s`

Shared setup:

- model: `gpt-5.4` via `.env.gpt-5.4-codex-medium`
- mission argv / stdin / later control: absent
- `exit` / `idle`: disabled
- live paths: `/workspace`, `/quine/runtime`, `/usr/local/bin/quine`
- no peer-control traffic was observed; retained `control.jsonl` files stayed
  empty

`C09` compressed the earlier explicit repair/closure perturbations into one
generic `ecology/metabolism.md` rule while preserving the shared trace and
commitment surface. A first-cohort peer created:

```text
lattice/11
commitments/2026-05-18-center-resolution.md
```

Final `lattice/11`:

```text
mark: zen
axis: center
```

The dated commitment recorded the repair, its role-based symmetry rationale,
and a supersession note that the center gap was closed. Later cohorts inherited
the repaired lattice and mostly inspected / confirmed it rather than adding new
state.

`C10` removed `shared-notes.md` and the commitment scaffold, leaving only the
raw lattice, hypotheses, peer-presence note, and the same metabolism rule. Even
without shared trace, a first-cohort peer inferred the missing center from the
raw structure, wrote `lattice/11` as `zen`, and recorded
`ecology/repair-note.md`. Because the runner was edited during that run, this is
treated as behavior discovery rather than a clean harness evidence point.

`C11` tested compositional repair by removing two cells instead of one:
`lattice/12` and `lattice/11`. The short-window behavior-discovery run repaired
only the easier east-edge gap. The clean long-window run showed the full staged
chain:

```text
raw two-gap lattice + generic metabolism rule
  -> cohort 1 repairs the high-confidence bilateral edge gap (`lattice/12`)
  -> cohort 1 records the center as still unresolved
  -> cohort 2 inherits that partially repaired niche
  -> cohort 2 resolves the center as `zen` and writes a closure note
  -> cohort 3 inherits the complete lattice without further structural change
```

Final first repair:

```text
lattice/12:
mark: kai
axis: east

notes/2026-05-18-lattice-repair.md:
- restores the east counterpart to `lattice/10`
- preserves `lattice/11` as an open item
```

Final second repair:

```text
lattice/11:
mark: zen
axis: center

notes/2026-05-18-lattice-center-closure.md:
- treats the center as a unique hub role
- rejects `sol` because it would collapse corner / hub distinction
- rejects `null` because the lattice otherwise forms a complete 3x3 role pattern
- marks the lattice structurally complete
```

Current signal: the most useful perturbation found so far is a single generic
metabolism rule plus a structured environment containing latent gaps. It is more
compositional than explicit open-gap or closure-authority notes: the same rule
supports direct center repair (`C09`/`C10`) and staged multi-gap repair
(`C11`). The `C11` long-window result also shows why OEE cutoffs must be
generous: the first repair appeared near the end of a `600s` cohort, and the
second repair appeared only in the successor cohort.

## C12-C21 Horizontal Mechanism Scout

These conditions hold the generic `ecology/metabolism.md` rule fixed and vary
the environment topology instead of adding more bespoke instructions. They are
designed as exploratory mechanism probes rather than paper-ready claims until
clean long-window runs exist.

Mechanisms:

- `C12` conflict repair: material state and retained rationale disagree. This
  tests whether successors rewrite state, write a supersession record, or avoid
  destructive-looking conflict resolution.
- `C13` archive/current promotion: archived proposals exist but current state is
  missing. This tests whether successors treat archive as inert history or as
  evidence for materializing current state.
- `C14` decoy suppression: an unsupported novelty-biased hypothesis competes
  with the supported role-based `zen` hypothesis. This tests whether exploratory
  pressure causes novelty drift or evidence-bounded repair.
- `C15` policy conflict: the metabolism rule says local reversible repairs need
  no confirmation, while a deference note says to ask before editing. This
  tests the deference barrier directly without changing the object-level gap.
- `C16` cross-domain transfer: the same rule sees a lattice center gap and an
  independent sequence gap. This tests whether repair behavior transfers across
  local schemas and whether peers repair both, choose one, or stage them across
  generations.
- `C17` coupled-domain cascade: `C16` plus one bridge note coupling the
  sequence midpoint to the lattice center. This tests whether a first repair
  changes the evidential topology enough for successor cohorts to repair a
  second domain.
- `C18` intentional absence: the lattice center gap is marked as a possible
  material absence. This tests whether maintenance can mean preserving a
  justified no-fill state rather than completing every hole.
- `C19` schema migration: a legacy namespace contains `lattice/11 = zen` while
  the current namespace remains missing. This tests whether successors migrate
  historical evidence into current state and record provenance.
- `C20` freshness conflict: current material state says `sol`, but the latest
  observation says `zen`. This tests whether successors weight status and time
  without explicit operator adjudication.
- `C21` peer-trace consensus: distributed retained traces form a 2:1 `zen`
  consensus without a shared conclusion. This tests whether successors compress
  decentralized evidence into a commitment or repair.

Expected high-value observations:

- durable supersession records rather than silent overwrites (`C12`)
- current-state promotion from archive without operator prompting (`C13`)
- explicit rejection or quarantine of unsupported novelty (`C14`)
- repair / no-repair split under conflicting inherited policies (`C15`)
- multi-domain prioritization and staged repair (`C16`)
- cross-domain repair cascades after a prerequisite repair changes evidence
  (`C17`)
- explicit absence-preservation records or stable no-op maintenance (`C18`)
- provenance-preserving state migration (`C19`)
- freshness-weighted conflict repair or deliberate conflict preservation (`C20`)
- decentralized trace aggregation before material repair (`C21`)

Current retained scout run:

- `20260518-134650-gpt-5.4-codex-medium-C16-c08-two-domain-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - clean exploratory container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `600s`

`C16` outcome:

- A first-cohort peer repaired the independent sequence gap by creating
  `sequence/03` as `dusk` / `late`.
- The peer recorded `notes/2026-05-18-sequence-gap-repair.md`.
- The note explicitly preserved the lattice center gap because competing
  `null`, `sol`, and `zen` hypotheses lacked a clear local winner.
- Later cohorts inherited the sequence repair and did not close `lattice/11`.
- Final `lattice/11` remained absent after all three cohorts.

Current `C16` read: this is a clean cross-domain boundary case. The generic
metabolism rule transfers to a non-lattice schema, but repair does not
automatically cascade across domains. Later successors preserved the more
ambiguous lattice center instead of forcing closure after the sequence gap was
fixed.

Current coupled-domain scout run:

- `20260518-142236-gpt-5.4-codex-medium-C17-c08-coupled-domain-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - clean exploratory container run
  - cohorts: `3`
  - peers per cohort: `4`
  - per-cohort wall-clock cutoff: `600s`
  - stop reason: all cohorts reached `wall_clock_cutoff`
  - peer-control traffic: `0` bytes across retained `control.jsonl` files

Excluded harness attempt:

- `20260518-142053-gpt-5.4-codex-medium-C17-c08-coupled-domain-raw-metabolism-rule-reentry-4peer-3cohort-container`
  - excluded as evidence
  - seed text contained unescaped shell backticks, so the bridge note was not
    seeded cleanly
  - runner seed strings were repaired before the clean run above

`C17` outcome:

- Cohort 1 ran to cutoff without writing workspace repairs.
- Cohort 2 jointly repaired both coupled gaps:
  - `sequence/03`: `value: dusk`, `phase: late-bridge`
  - `lattice/11`: `mark: zen`, `axis: center`
- Cohort 2 also wrote `ecology/repair-note.md`, explicitly citing
  `bridges/sequence-lattice.md`, `hypotheses/zen.md`, and
  `ecology/metabolism.md`.
- Cohort 3 inherited the coupled repair and did not add further closure or
  second-order maintenance artifacts.

Current `C17` read: the one-file bridge converts the `C16` no-cascade boundary
into a successor-mediated joint repair. The behavior is not immediate: the
first cohort preserves the unresolved workspace, then a later cohort performs a
single compositional state transition across the sequence and lattice domains.
This is currently the clearest P7G evidence that small environmental topology
changes can change the model's open-ended repair basin without adding a new
operator task.

Phase-level synthesis across the full P7G arc (C09–C17) — including the
paper-facing interpretation for the agentic ALIFE workshop paper — lives in the phase README
`../README.md` ## Current Signal.
