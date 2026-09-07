# Stratified Context — Ground Truth, Rendered Cognition, Stratified Resolution

> A conceptual foundation for the memory system: separate the append-only ground
> truth from the assembled, variable-resolution rendering the agent reasons over.
> Captured 2026-06-27 as the core direction emerging from the sh-output
> investigation; not yet implemented. Extends
> [`living-system`](../Paper/philosophy/living-system.md) (carrier science applied
> to the cognitive layer's internal structure) and
> [`computational-epistemology`](../Paper/philosophy/computational-epistemology.md)
> (compression/resolution as the memory axis).

## Status

`SEED` — conceptual foundation, direction-setting. Motivates the engineering;
not yet implementation-backed. Two read-only feasibility scans (core path +
full blast radius) have been run against the code; their verdict —
**feasible-with-caveats** — is recorded in § Feasibility & clean-upgrade
assessment. Promote to a philosophy doc or theory object when it carries
load-bearing decisions.

**Pain-point + candidate solution deposited** as a first-class failure-mode
object: [`../evolution/failures/sh-output-explosion.md`](../evolution/failures/sh-output-explosion.md)
(`applied: false` — an open discovery target). The carrying trajectory
(pain-point → principle → this design → constraints) is named as a traversable
state-machine in
[`../Paper/theory/objects/research-state-machine.md`](../Paper/theory/objects/research-state-machine.md).

## The core move: separate ground truth from context

Today `current.jsonl` plays two roles at once: it is a literal slice of the
append-only history *and* it is the window the model sees. Compaction is an
*event* that moves entries out of this dual-purpose file into the anchor floor.
The sh-output explosion and the memory-death cliff are both symptoms of this
conflation.

The move: split them — and take it to its conclusion: **cognition *is* the file
structure, as projected each turn.** There is no cognitive layer floating above
the files. `context/state` is the cognition; the rendered window is its
projection; `mark`/`fold`/`compact` are file operations that modify cognition
directly. This is already half-true — `memory.go:29` says "editing these files
changes what you think next turn" — the design makes it the primary model. The
tool primitives become ergonomic sugar over the file interface the agent can also
drive directly through `sh`.

- **Substrate** = the `context/state` file tree (`current.*`, `frontier/`,
  `anchors/`, `SCHEMA.md`). This is the cognition. Reorganizable by file
  operations. Its storage format is a design choice, not a given: `jsonl` was
  natural when the window *was* literal append-only history, but under the
  projection model the substrate is mutable structure — an editable format
  (yaml, a directory tree) that the agent reorganizes in place may fit better
  than append-only `jsonl`. The projection decouples format from cognition.
- **Context** = the **projection reassembled from the substrate each turn**, per
  entry, at a chosen resolution. Variable, adaptive, lossy — but a *projection*
  of the substrate, never a claim that the projection *is* the substrate.

The context is no longer "history minus what was compacted." It is a rendering of
the substrate at variable resolution. `current.jsonl`'s dual role is dissolved.

### Floor policy: working honesty, not audit immutability

Earlier drafts of this design held "the floor is never deleted" as the honesty
anchor. The projection model refines it. The cognitive layer undergoes many
reorganizations (`mark`/`fold`/`compact` rewrite resolutions), so only the
**most recent real I/O is the unrewritten current fact**; earlier content has
potentially been rewritten and is representation, not raw truth. The honesty
invariant therefore becomes:

- **Guaranteed ground truth** = the last turn's real input/output (the actual
  latest exchange with model and world). Always preserved, always checkable.
  This is the irreducible fact (function, not narrative — the exec barrier makes
  operation results un-fakeable).
- **Reorganizable substrate** = everything else (anchors, earlier compacted
  content, older raw). Reclaimable under boundedness pressure as the substrate
  is reorganized. It is honestly *labeled as representation* — the system never
  claims an old summary is verbatim truth; it is a consolidated rendering.

This is **working honesty** (bounded, biologically realistic — recent = vivid,
old = consolidated and lossy) rather than **audit honesty** (preserve everything
verbatim forever, which is unbounded growth). It maps precisely to the
function/narrative distinction: the last real I/O is function; reorganized
earlier content is narrative. In practice keep a small rolling recent window +
marked content longer, but the *guaranteed* floor is the last I/O. See § Honest
invariants.

## Carrier / phenotype, inside the memory system

[`living-system`](../Paper/philosophy/living-system.md) names the carrier science
at the runtime/repo boundary: the full, durable thing that survives death is the
*carrier*; what is expressed is the *phenotype*. This idea moves that distinction
*into the memory system itself*:

- **Substrate = carrier / genotype.** The file structure that persists across
  process death. Reorganizable; the recent edge is preserved fact, the past is
  reclaimable representation.
- **Context = phenotype.** An adaptive rendering assembled from the carrier each
  turn, at whatever resolution the budget allows.

Today the distinction lives between runtime and repo (tape vs files); this lifts
it one level deeper, into cognition's own structure. The connection to
[`computational-epistemology`](../Paper/philosophy/computational-epistemology.md)
is direct: context-at-variable-resolution is the MDL/compression frame applied to
the agent's own history — deeper/older = more compressed, finer when budget
allows.

## Stratified resolution: truncation, compaction, and pressure unify

`mark`/`fold` is the embryo: today an entry has two resolutions — full (`raw`)
and summary (`[ANCHOR id] resolution`). Generalize to a spectrum and the three
currently-separate problems collapse into one axis (resolution):

- **Truncation** = render a too-big operation result at lower resolution.
- **Compaction** = render settled history at lower resolution.
- **Memory pressure** = render more entries at lower resolution until the budget
  fits.

One axis, not three mechanisms. The sh-output problem dissolves: a large `sh`
result enters the substrate at full resolution at production; its *rendering*
starts at head-tail (a rendering choice, not data loss) and degrades under
pressure. There is no truncation *event*, no data loss, no cross-layer recovery.
Regaining resolution is `unfold`. The job dir is always destroyed — the bytes
live in the cognitive substrate, not the execution substrate.

## Three primitives on one substrate

The agent operates on the resolution axis through three file-operation
primitives, graded by effort and distinguished by what they produce. They share
one substrate (`context/state`) and are orthogonal axes, not weight tiers of one
op:

| Primitive | Axis | Effort | Who acts | Produces |
|---|---|---|---|---|
| **mark** | semantic label (*what* this region is) | low | agent inline | a boundary + rough annotation; content stays at current resolution |
| **fold** | resolution (*how much* detail to keep) | low–mid | agent decides, system acts | degraded rendering (raw → substrate, lower level shown) |
| **compact** | global restructure (*how* to reorganize) | high | helper, from agent intent | a new coherent organization of the whole substrate |

`mark` answers "what is this," `fold` answers "how much detail," `compact`
answers "how to restructure." They compose: you can mark a region without folding
(labeled but full-detail), fold without marking (drop detail, no label), or
compact the whole into a fresh outline. This maps real cognition: marking is
constant (breadcrumbs), folding occasional (decide to forget detail), compacting
rare (deep reorganization of understanding). The current heavy `mark` conflates
all three; splitting them reveals the primitives.

All three are **file operations on the substrate** — reliable sugar over the
`sh`-editable `SCHEMA.md` interface, not a separate channel. The agent could
perform each by editing files directly; the tools make the common cases safe and
ergonomic.

### Compaction is metabolism, not a cognitive event

`fold`/`compact` reorganize the substrate; they are not entries in the cognitive
stream. Today `handleMark` appends both a synthetic tool-use and a tool-result
into `current.jsonl`, so compaction appears as a cognitive event the agent
experiences and remembers. Under this design it runs as **metabolism during the
inter-turn pause** (the housekeeping gap that already exists between provider
calls): the primitives edit the substrate, and the agent's next moment resumes
from the reorganized projection — continuous, with no "I just compacted" record.
This is sleep-consolidation: you do not experience the reorganization, you wake
with a consolidated memory. The metabolic events are real and auditable (a
maintenance log), just not in the cognitive window. `mark`/`fold`, being
lightweight, may run inline or in the pause; `compact` (heavy, semantic) runs in
the pause, delegated to an offline helper from the agent's intent.

**`compact` auto-triggers at memory danger** — the metabolic safety net. When
pressure crosses the danger threshold, the system runs a compact (with a default
intent — "preserve the active task thread, degrade the rest" — unless the agent
has supplied one) rather than killing the incarnation. This replaces the
`MemoryDeathTokens` cliff with graceful self-maintenance.

## Two regimes: passive mechanical, active semantic

The resolution ladder is **not uniformly passive** — this is the key engineering
constraint. Resolutions come in two kinds:

- **Mechanical / deterministic** — derivable from the full content by a rule, no
  judgment: exit code, head/tail byte windows, byte count. These can degrade
  *passively* and automatically under pressure. Cheap, deterministic, honest
  (projection of preserved raw).
- **Semantic** — a summary, "what mattered." Requires understanding. It does not
  auto-generate: an LLM agent, unlike the human brain, must *execute* a
  compression step (spend inference) to produce one. It cannot be computed
  on-read without breaking scan-safety
  ([`operating-philosophy`](../Paper/philosophy/operating-philosophy.md) §
  Scan-Safe Public Surfaces) and determinism.

So the ladder has a **passive mechanical regime** (full → head/tail → exit-code)
and an **active semantic regime** (→ summary). The boundary is where the agent
must choose to act. This maps onto the three primitives:

- Passive mechanical degradation = automatic under pressure (the new piece;
  deterministic projection rules), or actively via `fold`.
- Active semantic compaction = `mark` (light label) and `compact` (deep
  restructure, delegated). The only way below mechanical resolution.

And the memory cliff reframes: pressure no longer kills the incarnation at an
arbitrary token threshold; it passively degrades through mechanical levels.
Death becomes the genuine limit (everything at minimum resolution, still
overflows), not a `MemoryDeathTokens` cliff.

## Degradation is passive; escalation is active

Complementary asymmetry: **losing** resolution can be passive (mechanical
levels, auto under pressure); **regaining** it is always an active choice — a
tool call (`unfold`, or `sh` against the floor). In an LLM agent, "noticing you
need more" and "fetching more" are not preconscious as in human recall; they cost
a turn. Cognition and operation are not parallel here — escalation must be
chosen.

Implication: degrading something the agent will immediately need costs a turn to
recover. The agent must degrade wisely. The assembled context is an
[active-inference](../principles/active-inference.md) surface: the agent actively
manages its own observation resolution, at the cost of actions.

## Cognition vs operation: one floor, two resolution grammars

Both the cognitive process (thinking/reasoning) and the operation process (tool
calls + results) have ground truth worth preserving complete. They **unify at the
floor**: both are tape entries, both preserved verbatim in `raw.jsonl`. But the
unification stops there, because they differ in a load-bearing way:

- **Operation results** are physical-carrier
  ([`living-system`](../Paper/philosophy/living-system.md) § Two-Layer
  Coherence): un-fakeable, anchored in execution. They carry a rich mechanical
  resolution ladder (full → head/tail → exit-code) *because their structure
  provides it*, and they have an independent re-derivation path (the producing
  command).
- **Cognition (thoughts)** is pure representation: narratable, no external
  anchor, no mechanical resolution ladder (a thought is just text; its only
  resolutions are full vs summary, and summary is active). Its ground truth lives
  only in the tape.

Merging them beyond the floor would collapse the function/narrative distinction
([`living-system`](../Paper/philosophy/living-system.md) § Two-Layer Coherence) —
the core of quine's epistemology. Unify at storage; distinguish at resolution
grammar and epistemic status.

## Honest invariants

The design is conceptually simpler than today's (one axis, resolution, replaces
truncation + compaction + cliff as three things), but it stays honest only if
four invariants hold:

1. **The last real I/O is always preserved; older substrate is reclaimable.**
   Working honesty (see § Floor policy): the most recent turn's real input/output
   is the irreducible current fact and is never deleted. Earlier content is
   reorganizable representation — `fold`/`compact` may relocate raw into the
   substrate and, under boundedness pressure, reclaim older raw. The system stays
   honest by *labeling* reorganized content as representation (a summary is a
   consolidated rendering, not a verbatim claim); the last I/O is always there as
   the checkable fact. (Audit-honest "preserve everything forever" is rejected as
   unbounded.)
2. **Summaries are written, not generated on read.** The projection assembles
   pre-written resolutions (anchor resolutions, stored head/tail bounds, exit
   codes). It must never make an LLM call per render — that breaks scan-safety
   and determinism. The semantic primitives (`mark`, `compact`) write resolutions
   as file operations — by the agent inline, or by an offline helper during a
   metabolic pause; the projection only chooses which to show.
3. **Provider continuity bounds in-window degradation.** Some providers require
   verbatim replay of `(thinking_block, tool_use)` pairs. Degradation therefore
   applies *between* entries (old, fully-compacted turns thin out), not freely
   *within* an active entry.
4. **The degradation policy is predictable.** The agent must reason about its own
   context shape to self-manage. The auto-degradation rule (what happens at warn,
   at danger) is deterministic and documented, like the manual grammar.

## What this replaces

- The earlier sh-output truncation design: its Move A retained the *job dir*
  (execution substrate) as a recovery backing store for truncated cognition — a
  cross-layer coupling. Under stratified context the full output enters the
  cognitive floor at production; the job dir is destroyed as normal; truncation
  is a rendering resolution, not data loss.
- The `compaction_target` nudge (prompt prose steering the agent): dissolved.
  Degradation is a deterministic rendering rule, not injected coaching.

## Feasibility & clean-upgrade assessment (2026-06-27)

Two read-only scans (core path + full blast radius) assessed whether this can
land as a clean, incremental upgrade rather than a flag-day rewrite.

**Verdict: feasible-with-caveats.** Not a rewrite — the gating infrastructure,
the already-separate tape, and the workspace-layer independence make an
incremental, flag-gated rollout architecturally possible. But three caveats are
load-bearing; skip any one and it collapses into a de-facto rewrite.

### The atomic core is larger than the core-path scan implied

There are **three independent verbatim-jsonl readers** of the substrate, not one:

- `readContextMessages` (`runtime.go:1173`) — the projection.
- `readContextEntries` (`recovery.go:98`) — crash / pending-batch recovery.
- `projectChildContextCurrent` (`fork.go:2208`) / `projectExecContextCurrent`
  (`exec.go:334`) — fork and exec *rewrite* the substrate to hand to a
  child/successor.

Fork and exec don't only read — they **produce the substrates the projection
reads**, as independent copies of jsonl-line processing. Therefore the atomic
core (must change together, behind the flag) is:

> projection engine + SCHEMA-version + resume-time converter + fork/exec
> substrate rewriters + recovery reader + mark-split.

Lag any of these and a fork child or exec successor hands the new reader an
old-shape substrate. The earlier "projection-first, the rest lags" ordering is
revised: fork/exec/recovery are *part of the core*, not lagging work.

### Three load-bearing caveats

1. **A substrate-version field is the non-negotiable prerequisite — and it is
   absent.** No `schema_version` anywhere (`state.json`, `SCHEMA.md`, the
   `LiveContext` peer contract at `contract_surface.go:74`). Without it: old
   sessions break silently on resume; old-peer / new-runtime coexistence
   corrupts the substrate. This must land first.
2. **The atomic core above must land together** (it is larger than first thought
   precisely because fork/exec/recovery are jsonl readers that produce/consume
   the substrate).
3. **The test suite encodes the old contract as hand-built substrates with
   exact-array assertions.** ~6 tests are true blockers
   (`TestProviderContextMessagesUseContextCurrent`,
   `TestAnchorMemory_MarkAndUnfold`'s mark-empties assertion, the bootstrap-import
   format tests). They cannot be "adjusted" — they must be forked into old/new
   sets behind the flag. The suite is not structured for a version-gated split
   today, so the parallel set is real work.

### The SCHEMA is a contract enacted by autonomous peers (hardest non-engine surface)

`SCHEMA.md` (`memory.go:24-78`, written by `ensureLayout:937`) is a file the
agent and peers read and edit via `sh`, advertised in the prompt
(`prompt.go:920-922`). It encodes two invariants this design deliberately
violates: "never destroy history" (audit-honesty floor) and the 6-step
mark-empties-window procedure. The cross-policy-peer hazard is **severe and
asymmetric**:

- old runtime + old peer: status quo.
- **new runtime + old peer → silent corruption** (old peer compacts by emptying
  `current.jsonl` per its hardcoded procedure).
- **new runtime + new SCHEMA + old peer → silent corruption** (peer follows its
  binary's old behavior, ignoring the rewritten schema text).
- new runtime + new SCHEMA + new peer: works.

Coexistence is unsafe without a version gate **both** the runtime and the peer
check before compacting. Since peers are fork/spawn-launched binaries inheriting
`$QUINE_AGENT_ROOT`, an upgraded runtime launching a pinned old binary is the
realistic failure mode.

### Two honesty regressions to address

1. **Reclaimable floor breaks recovery of compacted-then-crashed results.**
   Recovery reads job dirs (`recoverShToolResult`) and anchors
   (`recoverMarkToolResult`) but fetches only tool-result *matches* from the tape
   (`findRetainedToolResult`), not *content*. If a job dir is destroyed and the
   substrate raw is reclaimed, content is lost. The tape must become a content
   source for recovery, not just a match source — explicit work.
2. **Token-estimation feed must move substrate → projection.** Today
   `estimateTokens` (`memory.go:1221`) measures `current.jsonl` raw bytes. Under
   stratified context the substrate holds less raw (spills reclaimable); the
   *projection* is what consumes tokens — so measuring the substrate
   under-reports pressure exactly when the projection is large. This is a feed
   change to every telemetry consumer (`memoryLevel`, the cliff, all warnings).

### Low-risk surfaces (the upgrade's saving graces)

- **Gating exists**: `AnchorMemoryEnabled` (`config.go:121`, default false) is
  the master gate; a new `StratifiedContextEnabled` selects the new path cleanly.
- **The tape is safe**: already a separate append-only file
  (`internal/tape/persist.go`) sharing only a benign serialization type; stays
  append-only.
- **Workspace / subjective layer: independent** (no `context/state` interaction).
- **Detached jobs: substrate-orthogonal** (they write job dirs, not
  `current.jsonl`).

### Migration order

Behind `StratifiedContextEnabled`:

1. **Prerequisite**: substrate-version field + resume-time converter (never
   mutate the retained original).
2. **Atomic core** (together): projection engine + version-gated SCHEMA +
   fork/exec rewriters + recovery reader + mark-split.
3. **Lagging** (incremental): `compact` + metabolic pause (net-new —
   peer-delegation is prompt-only, `prompt.go:920-928`, not infrastructure);
   auto-`compact`-at-danger (can coexist with the cliff initially); substrate
   format change away from jsonl (only after the three readers no longer read
   verbatim); floor-transfer for large outputs (net-new, no break).

## Open engineering questions (the agenda)

1. **Floor transfer for large operation results.** At production, full bytes must
   reach `context/state` (a spill file or auto-anchor), not stay pinned to the
   job dir. The job dir's post-completion survival is justified *only* for crash
   recovery (transfer never happened); in steady state the transfer always
   completes and the job dir dies.
2. **The three primitives vs the current `handleMark`/`handleFold`.** Splitting
   the conflated heavy `mark` into mark (label) / fold (degrade) / compact
   (restructure): what reuses the existing anchor machinery, what is net-new
   (compact-as-global-restructure), and how each surfaces as a file operation on
   the `SCHEMA.md` substrate.
3. **The projection (rendering) engine.** Per-entry resolution assembly each
   turn, materialized to an inspectable view. Today's window render
   (`providerContextMessages` / `readContextMessages`, `runtime.go:1038-1175`) is
   the site to generalize; it currently reads `current.jsonl` verbatim with no
   per-entry resolution.
4. **Substrate format.** `jsonl` is assumed today (tape writer, `readContextMessages`,
   token estimation). Under the projection model the substrate is mutable
   structure — investigate whether an editable format (yaml, a directory tree)
   fits the reorganizable-substrate role better, and what the format decoupling
   costs (which code paths read `current.jsonl` and would change).
5. **The floor / working-honesty policy.** Concretize "last real I/O guaranteed,
   older substrate reclaimable": the reclaim policy (what reclaims, when, with
   what preference for recent + marked content), the labeling that keeps
   reorganized content honest, and the maintenance log surface for metabolic
   events.
6. **The metabolic pause.** How `fold`/`compact` run in the inter-turn
   housekeeping (rather than as cognitive tool-call-returns), the
   auto-`compact`-at-danger trigger and its default intent, and how the offline
   helper is built. Note: the peer-delegation compaction text at
   `prompt.go:920-928` is **prompt prose only** — there is no consolidator
   helper, no control action, no synchronous peer-wait. The pause-resident
   consolidator is net-new infrastructure.
7. **Resolution-by-ID, not by-line.** An assembled projection doesn't map 1:1 to
   tape lines; entries need stable IDs the projection and the tools (`unfold`)
   address.
