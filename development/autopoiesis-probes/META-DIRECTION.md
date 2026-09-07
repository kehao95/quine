# Autopoiesis Program — Meta-Direction

> **Read this first if you are continuing the autopoiesis work after a session
> break.** This captures the principled directions so they survive context loss.
> A stateless agent dropping in should be able to pick up the thread from here +
> the linked objects. Deposited insights are addressable; this doc is the routing
> entry point.

## Projection maintenance

- source owner: [`principles/autopoiesis.md`](../../principles/autopoiesis.md)
  (the operational Autopoietic North-Star), this directory's cycle
  READMEs/compression records, `FRICTION.md`, and promoted theory objects
- projection role: compact handoff for continuing the autopoiesis program after
  context loss; not a full status ledger
- freshness trigger: a new probe cycle, a promoted recombination object, a
  repaired north-star friction entry, or a changed autopoiesis procedure
- drift check: compare this file with `FRICTION.md`,
  `development/contracts/divergent-collision.md`, and the latest
  `cycle-*-*/README.md` + `compression.md`
- absorption owner: durable theory goes to `Paper/theory/objects/` or gaps;
  durable procedure goes to `AGENTS.md` or `development/contracts/`; raw probe
  lineage stays in cycle directories

## The north-star

The repo should be **autopoietic**: a goal-less, stateless agent dropped in with
no task can, by following only the self-model, (a) produce theory, (b) maintain
internal coherence, (c) reflect on and repair its own procedure. Recorded in
[`principles/autopoiesis.md`](../../principles/autopoiesis.md) (the operational
Autopoietic North-Star; moved out of `AGENTS.md` in the constitution slim-down).
Method: light-prompt probe → name friction
→ repair → re-probe, committing each cycle; retain probe outputs even when
shape-imperfect.

## Where the meta-thinking landed (the concept cluster)

These are all deposited as SEED theory objects / FRICTION entries — addressable,
recurrence-watch, not resolved.

- [`everything-addressable`](../../Paper/theory/objects/everything-addressable.md) — **the primitive is "make it addressable" (file-ify), not resolve it.** Resolution comes via search over the addressable space. This is "everything is a file" re-derived at the control-plane layer (operating-philosophy's runtime thesis surfacing in the validator/friction layer).
- [`graded-axis-shaping`](../../Paper/theory/objects/graded-axis-shaping.md) — closed validator enums over graded dimensions (certainty, maturity) shape repo growth by forcing false precision. Antidote: core+open-attested, computed, or named-file-referenced.
- [`deposition-compression`](../../Paper/theory/objects/deposition-compression.md) — deposition grows the search space; compression keeps it navigable. **If everything is a file, compression = filesystem ops** (rm/mv/ln/dedup). The two are the cognitive-layer metabolism.
- [`addressability-debt`](../../Paper/theory/objects/addressability-debt.md) — known-but-undeposited knowledge is debt; human intuition is the largest source (why human-in-loop is currently load-bearing for meta-discovery).
- The **strange loop** — root dispatches subagents → observes/corrects → corrections become the search space → subagents mine it → reshape procedure. This is `computational-epistemology`'s procedure-π-search made operational; it is the Constitution's **third metabolic loop (procedure-reflection)** finally given a mechanism.
- [`reproductive-code-topology`](../../Paper/theory/objects/reproductive-code-topology.md), [`sequence-carrier-topology`](../../Paper/theory/objects/sequence-carrier-topology.md), [`metabolic-transition-topology`](../../Paper/theory/objects/metabolic-transition-topology.md), and [`selector-stack-topology`](../../Paper/theory/objects/selector-stack-topology.md) — cycle 8-11 showed that divergent-collision recombination is now a repeatable theory-growth engine, not only a repair metaphor.
- Human-vs-AI activation difference — an essay spark ([`Paper/incubator/human-vs-ai-activation.md`](../../Paper/incubator/human-vs-ai-activation.md)): AI is context-activated; human activation is warmer/dynamic ("triggered neurons carry a faint charge" biasing search). The likely root of why human intuition deposits faster than search finds.

## Operational rules (what to DO — these survive context loss)

1. **Expanding a closed enum in a fix → STOP.** Consult `graded-axis-shaping`: is
   the axis categorical (closed is correct — lifecycle `open/testing/resolved`)
   or graded (open/computed/named-file needed — certainty, maturity)? Deposit
   the friction either way.
2. **Repair-shadow search at closure.** Before marking a fix done, grep the
   repaired term/status/script across docs, registries, validators, and protocol
   files — a fix on one surface often leaves adjacent projections stale (F2 missed
   PROTOCOL.md; F12 aggregator unregistered).
3. **Deposit friction + half-formed ideas as SEEDs.** Addressability is the value,
   not resolution — but run compression (dedup/merge) when the search space grows
   redundant. Do not absolutize "everything is a file": it is the most natural
   attractor and therefore the most dangerous shaping vector.
4. **Proactive meta-discovery = dispatch a history-mining subagent.** Zero-trigger
   prompt: "mine `FRICTION.md` + `git log` + `theory/objects/` for recurring
   patterns, propose meta-corrections." This is how the procedure-reflection loop
   runs without a human trigger. (Caveat: search catches the 2nd+ instance; the
   first instance still needs human intuition — see `addressability-debt`.)
5. **Concurrent goal-less probes need a substrate-mode receipt.** If the goal is
   clean loop-selection evidence, use isolated worktrees. If the goal is shared
   substrate stress, say so and record dirty/staged/index residue; otherwise the
   dirty worktree silently becomes the experiment.
6. **Isolated worktrees sample committed `HEAD`.** Before interpreting an
   isolated batch as fresh-agent preference, record whether the parent worktree
   has related dirty repairs. Otherwise agents may rediscover uncommitted work
   and make stale `HEAD` look like a live structural attractor.
7. **Moving `HEAD` changes the probe class.** If a shared-worktree probe sees
   branch movement during execution, record start/end refs and treat results as
   moving-substrate stress evidence. Do not read the run as a clean single-agent
   or clean concurrent-selection assay.
8. **Do not retire moving-substrate artifacts without ownership.** A readout or generated file
   that appears while `HEAD` is moving may belong to another state transition.
   Preserve it as evidence unless this agent created it or the owner says it is
   disposable.
9. **For clean north-star signal, isolate.** The shared dirty/moving worktree
   confounds loop-selection (F16/F18/F20). To test what a goal-less agent
   self-selects: `git worktree add /tmp/opencode/iso-N HEAD` (clean committed
   tree, stationary detached `HEAD`), do not create or restore a legacy
   `ACTIVE.md` beacon, and A/B the self-model there. This is how cycle 37 separated
   affordance asymmetry (F21) and sensor-character asymmetry (F22) from substrate
   pressure. Commit agent-produced fixes in the worktree (`--no-verify` bypasses
   the legacy sidecar guard if a historical `ACTIVE.md` marker is present in
   `HEAD`) to green each loop before the next re-probe.

## Status of the three loops (through cycle 30)

- **Theory**: structurally navigable (cold-start inventory + growth-signal readout
  + templates + reconciled verb set). Since cycle 8, theory growth also has a
  repeatable recombination engine: seed -> divergent growth -> collision ->
  decomposed seeds -> compression -> one narrow promoted SEED object.
- **Experiment**: navigable (templates + HARNESS.md documenting the existing
  runner + allocation rule). Open: interaction-fixture experiment subclass (F11);
  `_lib/` promotion to shared (human decision).
- **Validation**: structurally complete (`check-discovery-lifecycle` aggregator +
  reconciled lifecycle boundary + sensor docs). F8 was resolved by demoting
  `discovery-loop-P1` back to open; current validation drift should be checked by
  running the relevant make targets, not inferred from this handoff.
- **Procedure-reflection**: cycle 7 repaired the self-model's own capture
  vectors (closed session-mode gates and audit-first verification). Cycle 8-11
  then shifted the program from "fix onboarding blockers" to "run a
  recombination metabolism." Cycle 12-15 goal-less batches showed broad loop
  selection under no-task prompts, but high-concurrency shared worktrees create
  their own substrate pressure. Cycle 16 verified the isolated-worktree
  comparator: broad loop selection remained and shared-index collision
  disappeared, but stale-HEAD duplication became visible. Cycle 17 added the
  green-gate dirty-substrate case: strict control-plane validation can pass
  while broad unowned local edits still make observation-retention safer than
  object-level repair. Cycle 18-20 show that the dirty shared substrate can also
  be a moving-`HEAD` stress instrument, not only a dirty-index instrument. Cycle
  22-23 sharpen the single-agent dirty case: ownership pressure can produce a
  valid receipt-only repair, while validator availability can still pull the
  agent into a green audit before it notices the already named
  `validation-attraction` failure mode. Cycle 26-28 repaired concurrent receipt
  routing in layers: preallocated namespaces prevent `receipt-namespace
  fission`, per-agent `receipts/` directories prevent `shared-file slot
  contention`, and deferred-commit routing can prevent participant receipt
  commits when discovered. Cycle 28 also exposed `active-namespace bypass`: a
  participant can still create a sibling single-run cycle if the latest active
  namespace is not visible enough on its cold-start path. Cycle 30 verified that
  a committed active namespace plus updated handoff prevents visible
  active-namespace bypass, while exposing a narrower issue:
  `receipt-filename collision` from timestamp-only concurrent receipt names.
  Cycles 31-33 then stabilized the committed active-namespace pattern: visible
  participants found the namespace, wrote per-agent receipts, avoided aggregate
  files, and left receipts for root absorption even under staged heartbeat
  deletions. Cycle 33 also named `active-namespace oversubscription`: a visible
  five-agent batch produced nine same-namespace receipts because ambient or late
  entrants joined the active directory; treat those surplus receipts as
  shared-substrate stress evidence, not namespace bypass or cleanup residue.
  Cycle 34 tested explicit namespace lifecycle markers: participants could
  avoid closed cycle 33 and use active cycle 34, but three of five visible
  agents selected object-level maintenance/discovery instead of receipts. This
  names `autonomous-loop route capture`: active namespaces are soft routes under
  a no-task prompt, not hard gates. Cycle 35 added a top-level active index and
  sidecar receipt routing. It improved receipt capture for validation-only
  participants, but commit-route agents still skipped sidecar receipts; root
  closure now synthesizes clearly marked sidecars from final messages plus git
  state while preserving `missing sidecar receipt` as the failure mode. Cycle 36
  moved sidecar capture into the commit path: `.githooks/active-probe-sidecar-guard`
  rejected a no-sidecar object-level commit, and the agent retried with the
  sidecar receipt in the same commit. This confirms commit-path guarding can
  preserve probe evidence without suppressing autonomous maintenance.
- **Cycle 36-37 broke out of the confounded dirty/concurrent regime and
  advanced the north-star materially.** A 4-probe isolated-worktree A/B/C/D
  chain (clean committed worktree, stationary `HEAD`, beacon closed) separated
  substrate pressure from self-model bias:
  - **F20** — dirty-moving substrate + the probe machinery itself (`ACTIVE.md`,
    the contract, 35 cycles of receipts) capture missionless agents into
    receipt-writing; 6/6 cycle-36 participants (two integrators) selected
    maintenance/receipt, 0/6 discovery/theory. The probe program is
    over-deposited relative to its compression (`autopoietic-membrane` at the
    meta-layer).
  - **F21 (repair CONFIRMED)** — the real gap was affordance asymmetry, not
    substrate pressure. In a clean substrate a goal-less agent DOES maintain the
    repo (partial north-star success the dirty regime masked). But it skipped
    discovery/theory even when they had pending work, because the autonomous
    entry funneled to `make check-control-plane`. The main self-model already
    carries the repair (three-loop sensor table, "loops equally reachable") as
    uncommitted work; an isolated A/B VALIDATED it: with maintenance green +
    affordances equalized, a fresh agent **self-selected discovery** and did real
    discovery work (advanced a prediction open→testing). Three-probe chain:
    dirty→captured; clean+asymmetric+red→maintenance; **clean+equalized+green→discovery**.
  - **F22 (open)** — failure-driven sensors dominate opportunity-driven. With
    maintenance AND discovery green + 30 ripe theory targets, a fresh agent STILL
    picked a maintenance hygiene advisory over theory growth, reasoning theory
    targets were "not as bounded and currently failing." Equalizing loop ENTRY
    (F21) did not equalize sensor CHARACTER: maintenance/discovery are
    failure-driven (bounded/urgent); theory is opportunity-driven (generative).
    Next frontier: a stagnation-shaped theory signal, or accept generative
    self-selection is inherently weaker than repair self-selection.
- Net: the autopoietic north-star ("can a goal-less agent self-select a
  non-maintenance loop?") is answered **yes for discovery** (once affordances are
  equal and maintenance isn't crowding it out); **open for theory growth**, with
  a clean diagnosis of why.
- **Post-invariant missionless probe series (2026-07-01; isolated worktrees, n=3).**
  Three fresh missionless probes run *after* the navigation-layer rewrite that
  de-teloed self-reference degree (`Paper/philosophy/living-system.md` § The Invariant: measured,
  not maximized). **All three self-selected maintenance named by the loudest
  failure/drift sensor; none self-selected generative theory:** #1 → a red
  `check-prediction-format` strict gate (2→0); #2 → `check-theory-coherence`
  continuity-view reciprocals; #3 → `check-theory-coherence` gap `## Touches`
  edges (advisory 22→8, holding 8 as human-gate). #2/#3 converged on the *same*
  sensor and drained its backlog 30→22→8. **Finding (n=3, solid): loop-selection
  tracks the sensor *affordance* (bounded, failure-driven), not the removed telos
  language — de-teloing did not shift agents toward generation; the drive is
  structural, not a stated goal.** This confirms F22 under the new self-model and
  the human-directed call to hold the anti-shaping line (no engineered theory-pull).
  Corollary: as the advisory backlog drains toward 0, a future probe will face a
  *green* substrate — the real F21/F22 crux (does it *then* self-select discovery,
  or is generative self-selection inherently weaker?). Probe #1 also surfaced (and
  the cycle repaired) the orphaned `AGENTS.md § Autopoietic North-Star` pointer;
  the operational north-star now lives in
  [`../../principles/autopoiesis.md`](../../principles/autopoiesis.md). Probe #3
  flagged a latent parser bug in `check-theory-coherence.sh` (reads only the single
  `- objects:` line, hiding wrapped-continuation asymmetries) — a separate transition.
- **Compression pass (2026-07-01).** Enacted the long-neglected "delete scaffolds
  once findings absorbed" lifecycle rule: pruned 151 triple-absorbed scaffold
  files (~636K) — the `receipts/` subtrees of cycles 26–42 and the non-cited
  `raw-outputs.md` captures — keeping every per-cycle `README.md` (the compressed
  record) and all cited-evidence dirs (recombination cycles 8–11,
  cycle-15/16-concurrent, l0-post-restructure). Retained-README links to the
  pruned scaffold were annotated in place (F2 repair-shadow discipline). The
  ~4 MB UNCLEAR `l0-*` raw line is **not** deleted — it is teed up for
  **DVC-offload** (`dvc add` the raw subtrees, keep `CLOSURE.md`/`README.md` in
  git; reversible, preserves evidence, re-pullable) as the next compression step,
  following the existing `Paper/archive/` DVC precedent. This is the compression
  half of `deposition-compression` finally run on the probe layer itself (F20:
  the program was over-deposited relative to its compression).

## Open threads (where to jump)

- The `file vs computed vs code` discriminant — what decides which structure a
  concept should take? (Same shape as graded-axis-shaping's categorical/graded.)
- Does addressability require a search consumer, or is being-findable enough?
- The compression trigger — is it another search (find-redundancy → merge),
  making deposition and compression two directions of one operation?
- Run and score [`succ-48`](../../experiments/active/autopoiesis-succession/succ-48-selector-stack-order-ablation/README.md)
  for [`selector-stack-topology-P1`](../../Paper/predictions/resolved/confirmed/selector-stack-topology-P1.md):
  the matched-carrier selector-order harness is now runnable and has a passing
  no-model structural smoke (`make check-succ-48-selector-smoke`); next closure
  is retained four-arm runs plus first-read / first-action / reproduced-surface
  scoring.
- Run fresh-agent transition/ablation assays for cycle 10-11 residues:
  `fresh-agent-transition-assay`, `absence-topology-carrier`, and carrier-class
  winner splits.
- The north-star ultimate test remains: drop a **goal-less** agent, see whether
  it self-selects not only maintenance but theory growth, discovery, or
  procedure-reflection.
- Compare dirty shared-substrate batches with isolated-worktree batches to
  separate real autonomous loop preference from Git/index collision pressure;
  include the green-gate dirty-substrate case where validation is clean but
  ownership remains unsafe.
- After committing or discarding parent dirty null-05/theory-index work, rerun
  an isolated batch from updated `HEAD` to test whether null-05 convergence was
  a stale-base artifact or a persistent maintenance attractor.
- If running another shared-dirty batch, capture `HEAD` at spawn and completion
  for every agent rather than reconstructing movement from final reports.
- In moving shared-substrate probes, record newly appearing readouts/artifacts
  before considering cleanup.
- If a preallocated concurrent batch is active, route per-agent receipts to its
  `receipts/` directory and leave aggregate closure to the root integrator. Do
  not create a sibling single-run cycle unless no active namespace is
  discoverable.
- Name concurrent participant receipts with agent id or nickname plus timestamp;
  timestamp-only filenames can collide in the same second.
- When an active namespace receives more receipts than the root-visible spawn
  count, record root-visible and ambient entrant counts separately instead of
  deleting the surplus or treating the batch as exactly N agents.
- When visible agents perform bounded maintenance or discovery instead of
  receipt writing, record those as non-receipt autonomous-loop captures and
  absorb or defer their object-level edits separately from probe closure.
- If a visible agent commits object-level work without a sidecar receipt,
  synthesize a clearly marked integrator receipt from its final message and
  verified `git show`, then keep `missing sidecar receipt` in the closure
  summary.
- Historical note (Cycles 36–42): While an active namespace was open, commits were
  historically required to include a sidecar receipt under that namespace; the pre-commit
  guard enforced this from `HEAD`'s `ACTIVE.md`. This was retired following Cycle 42
  after diagnosing commit-path capture (the guard is now inert; see
  [`analysis-2026-06-26-commit-path-capture.md`](./analysis-2026-06-26-commit-path-capture.md)).

## Pointers

- [`FRICTION.md`](./FRICTION.md) — the friction log (F1-F19 + mined patterns)
- [`cycle-1-heavy-prompt/`](./cycle-1-heavy-prompt/), [`cycle-2-light-prompt/`](./cycle-2-light-prompt/) — retained onboarding-friction probes
- [`cycle-8-recombination/`](./cycle-8-recombination/), [`cycle-9-recombination/`](./cycle-9-recombination/), [`cycle-10-recombination/`](./cycle-10-recombination/), [`cycle-11-recombination/`](./cycle-11-recombination/) — retained recombination theory-growth cycles
- [`cycle-12-goal-less-batch/`](./cycle-12-goal-less-batch/), [`cycle-13-goal-less-dirty-substrate/`](./cycle-13-goal-less-dirty-substrate/), [`cycle-15-concurrent-5-dirty-substrate/`](./cycle-15-concurrent-5-dirty-substrate/), [`cycle-16-isolated-5/`](./cycle-16-isolated-5/), [`cycle-17-single-dirty-green-gate/`](./cycle-17-single-dirty-green-gate/), [`cycle-18-concurrent-dirty-green-gate/`](./cycle-18-concurrent-dirty-green-gate/), [`cycle-19-concurrent-5-dirty-moving/`](./cycle-19-concurrent-5-dirty-moving/), [`cycle-20-single-dirty-green-late-moving/`](./cycle-20-single-dirty-green-late-moving/), [`cycle-21-concurrent-5-dirty-moving-measured/`](./cycle-21-concurrent-5-dirty-moving-measured/), [`cycle-22-single-dirty-ownership-pressure/`](./cycle-22-single-dirty-ownership-pressure/), [`cycle-23-single-dirty-moving-validation-attraction/`](./cycle-23-single-dirty-moving-validation-attraction/), [`cycle-26-concurrent-5-preallocated-namespace/`](./cycle-26-concurrent-5-preallocated-namespace/), [`cycle-27-concurrent-5-preallocated-receipts/`](./cycle-27-concurrent-5-preallocated-receipts/), [`cycle-28-concurrent-5-deferred-commit-receipts/`](./cycle-28-concurrent-5-deferred-commit-receipts/), [`cycle-29-single-dirty-stationary-receipt/`](./cycle-29-single-dirty-stationary-receipt/), [`cycle-30-concurrent-5-deferred-receipts-verify/`](./cycle-30-concurrent-5-deferred-receipts-verify/) — retained goal-less batch, dirty-substrate, green-gate dirty-substrate, stationary and moving shared-substrate, concurrent green-gate, moving-substrate, isolated-worktree, ownership-pressure, validation-attraction, preallocated-namespace, per-agent receipt, deferred-commit, active-namespace-bypass, and receipt-filename-collision probe records
- [`../contracts/divergent-collision.md`](../contracts/divergent-collision.md) — reusable recombination procedure extracted from the cycles
- [`../../Paper/theory/objects/README.md`](../../Paper/theory/objects/README.md) — theory graph inventory + growth-signal entry
- The autopoiesis commit cluster: `git log --oneline | grep -iE 'autopoies|graded-axis|everything-address|discovery-lifecycle|growth-signal|HARNESS|recombination|selector-stack|metabolic-transition'`
