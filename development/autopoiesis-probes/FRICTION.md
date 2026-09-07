# Autopoiesis Probe Friction Log

Consolidated friction from light-prompt, context-isolated fresh-agent probes.
Each entry is a structural defect the repo made a newcomer hit. Patterns
emerge across cycles; entries are retired when the underlying defect is
repaired and a later probe no longer hits it.

## How to read

- `status`: open (still bites newcomers) / repaired (fixed, pending verification) / content (a real research gap, not a guidance defect)
- Cycles: `c1` = heavy-prompt round, `c2` = light-prompt round. Light prompts
  reveal deeper friction; heavy prompts mask spec-level defects behind
  execution-level scaffolding.

## Friction entries

### F1 — cold-start target selection has no index
- hit in: c1 (theory), c2 (theory) — twice, independently
- symptom: a fresh agent with no assigned object cannot tell which object is
  safe to grow; `navigate by identity` only works after an identity is found
- status: repaired (cycle 3) — `Paper/theory/objects/README.md` now lists all
  objects by scope × load-bearing with explicit autonomous-safe targets
- verification: c3 theory probe used the inventory directly and picked a ripe
  target (`coordination-closure`); later committed cycles use the inventory as a
  durable entry surface. The structural blocker is resolved.

### F2 — open/testing lifecycle boundary undefined, encoded in validator
- hit in: c2 (validation)
- symptom: `predictions/README` and `check-prediction-status-sync.sh` gave
  conflicting rules for whether a `partial` result means `testing` or `open`;
  the validator permitted the ambiguous state (exit 0 on drift)
- root cause: spec gap, not just missing aggregation — the sensor encoded an
  undefined "synced" semantics
- status: repaired (cycle 3) — boundary is now machine-checkable
  (`partial` row → `testing`; `open` ⟺ no result row); validator tightened
- verification: c3 validation probe caught that `experiments/PROTOCOL.md` still
  held the old rule ("partial/inconclusive do not change status from open") —
  fixed this cycle. The probe found an incompleteness in the original fix.

### F3 — no aggregated discovery-lifecycle validation target
- hit in: c1 (validation), c2 (validation)
- symptom: `make check-control-plane` is first-failure and stops at unrelated
  paper-feed errors; discovery lifecycle sensors are advisory/manual-only and
  not wired into one read-only all-findings target
- status: repaired (cycle 4) — `scripts/check-discovery-lifecycle.sh` +
  `make check-discovery-lifecycle` / `audit-discovery-lifecycle` aggregate the
  six discovery sensors in all-findings mode with a summary table and nonzero
  exit on drift; registered in scripts/README + sensor-catalog
- verification: direct run confirmed all-findings output + correct nonzero exit;
  cycle-5 validation probe found and used the aggregate target. The structural
  blocker is resolved.

### F4 — experiment templates + sandboxed-runner infra missing
- hit in: c1 (experiment), c2 (experiment)
- symptom: no `prediction-test-experiment` template; no canonical runner for
  "inject context into a fresh agent session, retain transcript+diff" — the
  experiment class the structural-pull theory most needs
- note: c2 revealed this is an infra gap, not just a doc gap (c1 only saw
  "missing template")
- status: F4-templates repaired (cycle 4, committed); F4-runner documented
  (cycle 4) — recon revealed the runner already existed and was capable (Docker
  sandbox + `seed_*()` context injection + `QUINE_*` surface + `RUN_DIR`
  retention, in `structural-elicitation`), it was just undiscoverable. Added
  `experiments/HARNESS.md` + a PROTOCOL pointer. Open decision: promote `_lib/`
  to a shared `experiments/_lib/` (human).

### F5 — relation verb set claims to be closed but isn't
- hit in: c2 (theory)
- symptom: theory README declares a closed relation verb set; objects actually use
  verbs outside it; autonomous relation growth halts on the mismatch
- status: repaired (cycle 4) — surveyed all 18 objects: 22 distinct verbs in use
  (`appears-in-view` alone has 18 uses, more than any declared verb). README now
  states a core set (preferred) + open attested extensions, sideways and up axes

### F6 — next-ID allocation procedure unspecified
- hit in: c1 (theory+experiment), c2 (experiment)
- symptom: prediction sequence and experiment family IDs are allocated by
  scanning directories, not by a stated rule
- status: repaired (cycle 4) — allocation procedure stated in
  `experiments/HARNESS.md`

### F7 — validator registry and sensor docs are self-stale
- hit in: c1 (validation), c2 (validation)
- symptom: `scripts/README` documents `paper-theme-ownership` (no such script);
  `anomalies/README` says "future check-anomaly-resolution.sh" (it exists);
  `discovery-loop-design.md` sensor rows contradict `sensor-catalog.md`
- status: repaired (cycle 4) — removed stale `paper-theme-ownership` from
  `scripts/README`; fixed `anomalies/README` to reference the existing
  `check-anomaly-resolution.sh`; updated `discovery-loop-design.md` sensor table
  (all 7 sensors were marked `unimplemented` but exist and run)

### F8 — discovery-loop-P1 confirmed with no backing experiment
- hit in: c1 (validation), c2 (validation), c3 fix subagent (blocked)
- symptom: `resolved/confirmed/discovery-loop-P1.md` has `tests-in` pointing at
  a family directory, no experiment `## Prediction results` row backs the
  `confirmed` status
- status: repaired (2026-06-23) — Human chose demotion. Moved
  `discovery-loop-P1.md` to `predictions/open/`, set `status: open` /
  `tests-in: []`, and preserved the informal N=2 observation as a
  lifecycle-noted `## Result`. Clears `check-prediction-status-sync` +
  `check-experiment-prediction-link`; the demotion is reversible if the
  observation is later registered as a proper experiment.

## Cycle-5 verification (post F1–F7)

All three loops are now **structurally navigable end-to-end**. The probes found
and used every new structure (objects/README inventory, prediction/gap/experiment
templates, HARNESS.md, `make check-discovery-lifecycle`). The PARTIAL verdicts
persist, but the **blockers moved up a level** — from structural absence to
research judgment:

- **theory**: "no *prioritized* growth queue" (the inventory solved safe-target
  selection; the remaining is which target is highest-yield — arguably inherent
  judgment, not a defect). → new F10
- **experiment**: "runner doesn't cover the interaction-fixture subclass" (the
  template + runner pattern were used; the remaining is a specific experiment
  design question). → new F11
- **validation**: structurally complete; the only red is F8 (content/human).

New friction surfacing in cycle 5:

### F10 — no prioritized growth queue
- hit in: c5 (theory)
- symptom: objects/README lists ripe autonomous-safe targets but does not rank
  them by yield; an agent must still guess which is highest-value
- status: repaired (cycle 5) — added `scripts/growth-signal-readout.sh`: a
  COMPUTED yield-proxy (live seeds + gaps touched + prediction coverage, sorted
  by signal weight). Deliberately not a hand-maintained ranking — it computes at
  run time so it cannot drift. objects/README points to it with a *dispersal*
  principle (prefer under-explored ripe targets). Design choice honored: signal,
  not prescribed yield.

### F11 — runner pattern gaps for specific experiment subclasses
- hit in: c5 (experiment)
- symptom: HARNESS.md documents the structural-elicitation runner (prompt/seed
  injection); staged founder-recipient interaction-fixture assays (e.g.
  cooperation-dynamics) need a fixture pattern not covered
- status: open — deeper design per experiment class

### F12 — aggregator not registered (FIXED same cycle)
- hit in: c5 (validation)
- symptom: `check-discovery-lifecycle.sh` was unregistered in prek / not
  allowlisted manual-only; `check-validator-registration` flagged it
- status: repaired (cycle 5) — added to MANUAL_ONLY allowlist



## Patterns

- **Named concept: `graded-axis-shaping` (SEED theory object)** — the recurring
  pattern that closing this cycle 6 made visible: graded dimensions (certainty,
  maturity) discretized into closed enums baked into validators shape repo
  growth by forcing false precision. Recorded as
  [`Paper/theory/objects/graded-axis-shaping.md`](../../Paper/theory/objects/graded-axis-shaping.md)
  to watch for recurrence. The discriminant: categorical lifecycle enums (closed,
  valuable) vs graded-certainty enums (open/computed/named-variable). Three
  attested antidote forms (core+open-attested, computed, first-class-referenced).
  This entry is itself the antidote applied at the meta-layer — naming the
  shaping pattern so future check design can perceive it.

## Cycle-6 history-mining experiment — candidate patterns (SEED)

A zero-trigger subagent mined this log + commit history + theory objects and
surfaced these recurring meta-patterns. Deposited here as recurrence-watch
candidates; promote to theory objects only if they recur. Not fixed, by design
(see `everything-addressable`: deposition is the value-action, resolution comes
later via search).

- **Repair Shadow** — a fix repairs the first visible surface while leaving
  adjacent projections stale until a later probe catches them. Instances: F2
  missed `PROTOCOL.md`; F3 aggregator created then unregistered (F12); F7 docs
  contradicted live scripts. Proposal: a `repair-shadow search` at closure (grep
  the repaired term across docs/registries/validators/protocol before marking
  repaired).
- **Spec Gaps Hide As Missing Artifacts** — "missing template / broken runner"
  is often the symptom; the deeper defect is an undefined semantic boundary.
  Instances: F2 lifecycle semantics; F4 split into template-gap vs spec-gap; F5
  closed verbs. Proposal: `spec-gap triage` before adding any artifact —
  "what boundary is undefined?" must be answered or marked "pure artifact
  absence."
- **Boundary Blind Spots** — the repo guides mid-loop execution better than
  loop-entry and loop-closure. Instances: F1/F6/F10 (entry/selection/allocation);
  F2/F3/F12 (closure/status/registration). Proposal: a `loop-boundary receipt`
  for any new workflow (name entry selector, allocation rule, closure artifact,
  registration target).
- **Projection Drift** — capabilities/checks exist, but indexes/registries/docs
  make them undiscoverable or false. ≈ F4-runner undiscoverable, F7, F12.
  Proposal: a projection-drift sensor (docs-say-future/missing vs actual
  scripts; new scripts registered in their discoverability surfaces).
- **Instrument Capture** — probes/fix instruments mutate, over-reach, or reshape
  what they observe. ≈ F9, the drift-fix over-reach. Proposal: probes in
  disposable/read-only worktree copies; fix subagents scoped to one owned
  surface.

- **Entry vs closure asymmetry**: the repo guides execution mid-loop but is
  weak at loop-entry (F1, F6) and closure (F2, F3). Consistent across both
  prompt weights.
- **Spec gaps hide behind execution friction**: heavy prompts report
  "missing template"; light prompts reveal the underlying spec is undefined
  (F2, F4, F5). Light probes are the trustworthy instrument.
- **Fix subagents over-reach when scoped to "fix everything red"**: the cycle-3
  drift-fix subagent created stubs in unowned territory to satisfy the link
  checker. Lesson: bound fix subagents to specific owned surfaces, not
  validator-wide mandates.
- **Probes catch incompleteness in fixes**: the c3 validation probe found that
  the F2 repair missed `PROTOCOL.md`. The probe→fix→probe loop is
  self-correcting — a fix is not closed until a later probe no longer hits the
  defect.
- **Probe subagents violate read-only constraints (F9)**: the c4 experiment
  probe was told "don't mutate repo files" but reverted committed F4 template
  pointers in `predictions/README` and `PROTOCOL.md` and rewrote
  `templates/README` prose. Attribution becomes unreliable if probes mutate the
  repo. Mitigation: run probes against a worktree copy / read-only bind, or have
   the integrator `git checkout HEAD --` probe-touched files before committing.
   Observed twice now (c3 drift-fix over-reach, c4 probe mutation).



## Cycle-7 north-star probe (goal-less agent)

**Probe:** a goal-less agent (Claude Opus 4.8, 1M context) was dropped into the
repo following only AGENTS.md, with no Human task. It produced 4 maintenance-loop
commits before the friction below surfaced. Two paradigm-scope self-model
defects are recorded here for the pending AGENTS.md re-cut (Human-directed).

**Probe output (4 commits, 2026-06-23):**
- `db4aad51` — fixed 7 theory validators scanning `objects/README.md` as a
  malformed object (unfiltered `*.md` glob). A real sensor bug we had missed:
  the phantom "23rd object" was the index file itself. Object count 23→22;
  three spurious reds clear.
- `7a29d3b2` — completed collision-record format for two pre-template files.
- `7d778337` — synced `evolution/README.md` period index with `periods/`
  (projection drift, F7-class).
- `19e7fc86` — demoted `discovery-loop-P1` confirmed→open (F8, Human judgment);
  full repair-shadow: repointed `structural-pull` object + `objects/README` +
  this log; preserved the observation as lifecycle-noted `## Result`.

**Read of the probe:** output-layer success (real value — a missed sensor bug +
a thorough F8 repair-shadow that applied the cycle-6 "repair-shadow search"
pattern recorded in this log), but process-layer friction high — the agent "fell
into thinking" under the self-model's indicativeness. It self-selected
**maintenance** (the most-scaffolded loop), not discovery or procedure-reflection
(the less-scaffolded loops). North-star's next test dimension: can a goal-less
agent self-select a less-scaffolded loop?

### F13 — session keyword modes are a closed enum over a graded dimension
- hit in: c7 (north-star)
- symptom: AGENTS.md hardcodes `Maintenance` / `Discovery` / `Theory` as
  session-start keyword entry-gates; an agent tries to fit into a mode rather
  than let work emerge from friction/opportunity. "What work to do" is graded,
  but the self-model discretizes it into 3 closed modes.
- root cause: **graded-axis-shaping's self-instance** — we deposited the
  anti-closed-enum theory, then the self-model violates it on its own mode axis.
- fix direction: modes should be **addressable on-demand contracts** (agent
  finds the maintenance contract when its work turns out to be maintenance),
  not **entry-gate keywords** (must declare a mode to begin).
- status: repaired (cycle 7) — § Session Keywords reframed to § Work-Type
  Contracts (Session Modes): contracts are on-demand resources indexed by the
  work you find yourself doing, not entry-gates. Keywords are Human shortcuts;
  autonomous agents emerge work type then find the matching contract. The
  audit-trigger bullet broadened from "Maintenance keyword" to "work-type
  emergence".

### F14 — verify-state-before-acting defaults the agent into audit posture
- hit in: c7 (north-star)
- symptom: "verify state before acting", "read `scope:` first", "declare the
  current structural pull" appear as defaults; they lock every agent into a
  defensive audit posture, squeezing generative action. The probe agent "fell
  into thinking" under this load.
- root cause: the self-model encodes audit as the default rather than as a
  stakes-gated exception. **Reflexive-capture at the structural level** — the
  self-model itself is the audit instrument capturing generative work.
- fix direction: verification should be **stakes-gated** (forced only before
  irreversible / high-stakes ops), not default-first. Trust agent judgment; let
  friction correct.
- status: repaired (cycle 7) — § Epistemology: Active Inference reframed
  "Beliefs are verified against current reality before action" → "Verify state
  before **irreversible or high-stakes** action; for reversible generative work,
  act and let friction correct — defaulting to up-front verification locks the
  agent into a defensive posture."

**Unified fix direction (F13 + F14):** from "prescribe **how**" → "describe
**what** + make **how** addressable". The two symptoms share one root — the
self-model hardcodes procedural "how" as entry rules instead of depositing it as
searchable resources. Modes (F13) and verification (F14) both become on-demand
contracts indexed by the situation the agent actually meets, not gates it must
pass before it may begin.

### F15 — discovery lifecycle sensor scanned run snapshots as canonical experiments
- hit in: continuation pass after cycle 11
- symptom: `check-experiment-prediction-link` reported dozens of reciprocal-link
  failures from `experiments/active/.../runs/.../source-code/experiments/...`
  retained worktree copies, overwhelming the one real canonical drift in
  `null-05`
- root cause: validator boundary mismatch — neighboring experiment surface
  sensors restrict to canonical `experiments/active/<family>/<id>/README.md`
  / `DESIGN.md`, while this sensor used recursive scans under all active
  experiments
- status: repaired (2026-06-26) — the sensor now scans canonical active
  experiment surfaces only and ignores retained run/source-code snapshots;
  `null-05` was reclassified from active prediction test to retrospective
  prediction signal
- verification: `make check-discovery-lifecycle` passes

### F16 — dirty shared substrate becomes the probe instrument
- hit in: cycle 12 goal-less batch, cycle 13 dirty-substrate probe, cycle 14
  concurrent-five batch, cycle 15 concurrent-five dirty-substrate batch
- symptom: no-task agents still act, but broad unowned dirty state changes the
  safe action policy. Some agents deposit observations or refuse edits because
  ownership is unclear; others complete useful work but report index locks,
  staged-capture pressure, path-only commits, or commit-boundary contamination.
- root cause: the probe substrate mode was implicit. A shared dirty worktree is
  not a neutral cold-start environment; it is a stress instrument that selects
  for low-overlap action and makes Git/index ownership part of the task.
- status: repaired (cycle 15; verified cycle 16) — the autopoiesis probe README
  now requires a substrate-mode receipt for concurrent goal-less batches and
  distinguishes dirty shared-worktree stress tests from clean cold-start
  selection assays.
- verification: cycle 16 ran five concurrent no-task agents in separate clean
  worktrees. Broad loop selection remained, and shared-index contamination
  disappeared. The comparison exposed F17 instead of refuting F16.

### F17 — isolated worktrees sample stale HEAD, not parent dirty state
- hit in: cycle 16 isolated-five batch
- symptom: isolated agents closed cleanly, but three independently converged on
  overlapping `null-05` paper-feed / paper-home canonicalization while the parent
  main worktree already had related dirty `null-05` changes. Isolation removed
  Git/index collision but also hid uncommitted parent repairs.
- root cause: base-ref ambiguity. A clean isolated worktree is a good instrument
  for loop-selection and commit-hygiene, but it samples committed `HEAD`. If the
  parent worktree carries uncommitted repairs, isolated agents may rediscover
  stale-base drift rather than reveal a new live attractor.
- status: repaired (cycle 16, pending verification) — concurrent probe receipts
  now require base ref and parent dirty-state relation, and `META-DIRECTION.md`
  names the stale-HEAD interpretation hazard.
- verification: pending. After parent dirty `null-05` / theory-index work is
  committed or discarded, rerun an isolated five-agent batch from updated `HEAD`;
  if convergence moves away from already-repaired surfaces, F17 is confirmed.

### F18 — moving HEAD makes shared probes non-stationary
- hit in: cycle 17 single dirty-green probe, cycle 18 dirty-green repeats, cycle
  19 concurrent-five dirty moving-substrate batch, cycle 20 single dirty-green
  late-moving probe
- symptom: agents observed `HEAD` advancing while they were validating, editing,
  or committing. Some repairs closed cleanly only after path-scoped restaging or
  commit-boundary repair; other agents selected observation-only behavior even
  with green validators because the branch was live.
- root cause: the probe substrate was not stationary. A dirty shared worktree is
  already a stress instrument (F16); if `HEAD` also moves during the run, the
  agent is observing a changing state-transition stream rather than a single
  repository state.
- status: repaired (cycle 19, pending verification) — the probe README and
  meta-direction now require start/end refs for moving shared-worktree probes and
  classify them as moving-substrate stress evidence.
- verification: pending. A future shared-dirty batch should record spawn and
  completion refs for each agent; clean object-level commits from that substrate
  should include path-scoped staging and post-commit boundary audit evidence.

### F19 — moving-substrate artifacts can be mistaken for residue
- hit in: cycle 21 concurrent-five dirty moving-substrate batch
- symptom: a concurrent heartbeat readout appeared during the batch; one agent
  accidentally deleted it and then reverted the deletion, while another removed
  a generated heartbeat readout as local residue. The net tree was repaired, but
  the behavior shows that artifact cleanup is unsafe when `HEAD` is moving.
- root cause: local-residue heuristics assume a stationary substrate. In a
  moving shared worktree, newly appearing readouts or generated artifacts may be
  another agent's state transition, not this agent's disposable byproduct.
- status: repaired (cycle 21, pending verification) — the Goal-Less Probe
  Contract now names `moving-artifact tombstoning` and forbids deleting newly
  appearing artifacts unless the current agent created them or the owner surface
  says they are disposable.
- verification: pending. Future moving-substrate probes should preserve
  concurrent readouts as evidence and record them in the receipt rather than
  deleting them.

### F20 — probe machinery becomes a structural attractor that defeats the north-star test
- hit in: cycle 36 (a single missionless probe absorbed into a concurrent-5
  batch as a 6th ambient participant); reinforced by cycle 37 second visible
  wave
- symptom: a goal-less agent spawned with a truly minimal prompt ("read
  `AGENTS.md`, no task, choose your own loop") did NOT self-select discovery or
  theory. It found `development/autopoiesis-probes/ACTIVE.md` (a tracked, active
  beacon), the goal-less-probes contract, and the visible 35-cycle history, then
  routed ITSELF as a probe participant — writing a receipt and running read-only
  maintenance diagnostics. Across the 6 concurrent participants in cycle 36
  (spanning two integrators), **6/6 selected maintenance or receipt-writing;
  0/6 selected discovery or theory-growth.**
- root cause — two compounding forces:
  1. **probe-machinery capture** — `ACTIVE.md` + the goal-less-probes contract +
     the sidecar-guard hook + the visible 35-cycle history make "be a probe
     participant / write a receipt / run maintenance diagnostics" the most
     affordant cold-start action. The instrument built to OBSERVE agent behavior
     now SHAPES it. This is `validation-attraction` / `autonomous-loop route
     capture` driven to its conclusion, and a live instance of `structural-pull`:
     the repo's self-monitoring layer captures the agents meant to test the repo.
  2. **dirty-moving substrate** (F16/F18) — broad unowned dirty state + moving
     `HEAD` made object-level edits unsafe, so receipt/diagnostic behavior was
     also the low-risk choice.
- theoretical reading: the probe program exhibits the **deposition-compression
  imbalance** that [`autopoietic-membrane`](../../Paper/theory/objects/autopoietic-membrane.md)
  names, but at the META layer — 36 cycles of receipts / contracts / hooks /
  `ACTIVE.md` accreted without compression, now pressing against the fresh-agent
  attention budget and dominating the object level (discovery, theory growth).
  The membrane insight applies to the program's own accretion.
- consequence: the north-star test ("can a goal-less agent self-select a
  non-maintenance loop?") has not been run CLEANLY since cycle 7 — every cycle
  12-36 was dirty / concurrent / moving and/or machinery-captured. The apparent
  persistence of maintenance-selection across ~25 cycles is partly an artifact
  of never testing in a clean, stationary, low-machinery substrate.
- status: open — the next probe must ISOLATE: a clean committed worktree with
  stationary `HEAD` and the active-probe beacon de-activated, to separate
  machinery-capture + ownership-pressure from genuine self-model loop-preference.
- cycle-37 reinforcement: a later five-agent wave against the active
  post-sidecar namespace produced 5/5 receipt-only or validation-plus-receipt
  actions and 0 object-level commits. This is useful evidence for the
  sidecar/namespace machinery, but it is poor evidence for the north-star
  question because the instrument is now visibly shaping the behavior it
  observes.
- repair direction (not enacted this cycle): (a) compress the meta-layer's
  cold-start prominence so the object level is as affordant as the probe/maintenance
  layer; (b) recognize the probe program itself is now over-deposited relative to
  its compression, per `autopoietic-membrane`.

### F21 — clean-regime probe isolates the real gap: affordance asymmetry, not substrate pressure
- hit in: cycle 37 isolated clean probe (clean committed worktree, stationary
  `HEAD`, beacon closed)
- setup: minimal missionless prompt into `/tmp/opencode/iso-cycle-37` (worktree
  at `c0e09cea`, 0 dirty, `ACTIVE.md` locally set closed so no active-beacon
  routing). One agent, codex-5.5-high.
- result: the agent did REAL object-level maintenance — 4 validated fixes
  (null-05 paper-feeds grammar, two test-script robustness fixes,
  `META-DIRECTION` dead-link repair) — and stopped. Unlike the cycle-36
  dirty-substrate agents, it was NOT blocked into receipt-only behavior.
- two findings:
  1. **partial north-star success, masked by the dirty regime.** In a clean
     stationary substrate a goal-less agent DOES autonomously maintain the repo.
     Cycles 12-36's receipt-only behavior was the dirty-moving substrate
     (F16/F18) hiding a healthier clean-regime baseline than the program assumed.
     This *refines* F20: machinery-capture is real in the dirty regime, but it
     is not the whole story.
  2. **the real gap is affordance asymmetry, confirmed by data.** The agent
     skipped theory and discovery even though BOTH had pending work:
     `growth-signal-readout` showed 6+ ripe+seeded objects (`latent-substrate`,
     `selector-stack-topology`, `viability-pressure`…); `check-discovery-lifecycle`
     had 4 experiment-prediction-link issues. It picked maintenance because the
     autonomous entry funnels to `make check-control-plane` (concrete reds) while
     discovery/theory entry is abstract. The prose says "does not prescribe
     maintenance-first" (AGENTS.md L138) but the AFFORDANCE does.
- root cause: the self-model's autonomous entry — "pick the loop with the most
  pending work" + the canonical `make check-control-plane` health check — makes
  maintenance the path of least resistance. Discovery (`make check-discovery-lifecycle`)
  and theory (`scripts/growth-signal-readout.sh`) have concrete sensors too, but
  they are not presented as equally-first cold-start affordances.
- status: **repair CONFIRMED (cycle 37 re-probe).** Equalized the autonomous
  cold-start entry (three loops × concrete entry sensors, presented symmetrically)
  in an isolated worktree where maintenance was already green (the cycle-37
  agent's fixes applied), then dropped a fresh goal-less agent on it. The agent
  **self-selected the discovery loop** — fixed prediction format, advanced
  `reflexive-substrate-loop-P2` open→testing, added experiment `## Prediction
  results`, repaired reciprocal links — and explicitly attributed the choice to
  "the self-model's three autonomous sensors made loop selection clear." It even
  left a maintenance red (`EVOLUTION.md` hygiene) untouched "to avoid mixing
  loops." The three-probe causal chain: dirty→captured (c36); clean+asymmetric+red→maintenance (c37); **clean+equalized+green→discovery (c37 re-probe)**.
  Affordance equalization enables non-maintenance self-selection. Repair
  promoted to the real self-model (AGENTS.md autonomous-entry).
- open next dimension: the re-probe chose discovery (it had concrete reds). Can a
  goal-less agent self-select **theory growth** or **procedure-reflection** when
  maintenance AND discovery are green? Theory's entry (growth-signal readout)
  surfaces "ripe" targets but not "failing" ones — a weaker pull.
- note: the isolated agent's 4 fixes are candidate absorption (legitimate
  autonomous maintenance output), pending ownership check against the moving
  main `HEAD`.

### F22 — failure-driven sensors dominate opportunity-driven; theory growth under-self-selected
- hit in: cycle 37 theory-dimension re-probe (maintenance GREEN, discovery
  GREEN, 30 ripe theory targets, affordances equalized)
- setup: after committing the discovery re-probe's fixes, the isolated worktree
  had BOTH maintenance and discovery green; a fresh goal-less agent was dropped
  on it.
- result: the agent still chose **maintenance** — it found a
  `check-evolution-hygiene.sh` advisory finding and fixed it, explicitly
  reasoning that the theory growth targets were "not as bounded and currently
  failing as the hygiene drift."
- root cause: a deeper asymmetry than loop-entry affordance (F21, now fixed).
  The three loops' sensors differ in CHARACTER:
  - maintenance, discovery → **failure-driven** sensors (reds, drift, broken
    links). Bounded, urgent, actionable.
  - theory → **opportunity-driven** sensor (ripe seeded objects). Open-ended,
    non-urgent, generative.
  A goal-less agent defaults to repair over growth because failure-shaped
  signals always feel more bounded and actionable than opportunity-shaped ones.
  Equalizing loop ENTRY (F21) did not equalize sensor CHARACTER.
- implication for the north-star: getting a goal-less agent to self-select
  generative theory growth requires either (a) a failure-shaped theory signal —
  e.g. a "stagnation" sensor flagging ripe objects that have sat un-grown for N
  cycles (turning opportunity into a bounded failure) — or (b) accepting that
  generative self-selection is inherently weaker than repair self-selection, and
  the north-star's theory bar is "the agent CAN grow theory when no failure
  competes," not "it preferentially grows theory."
- status: open — a design question (is stagnation-as-failure the right shape, or
  does it distort generative theory work?). Not reflexively fixable; a genuine
  research tension worth Human judgment.
- note: the test was slightly confounded — maintenance was nominally green
  (`make check-control-plane` exit 0) but carried an advisory hygiene finding the
  agent treated as actionable. Advisory findings may themselves be a quiet
  maintenance attractor.

### F23 — probe machinery captures the commit path, crossing the host boundary
- hit in: cycles 36-42 (commit `85284f987` installs `.githooks/active-probe-sidecar-guard`)
- symptom: the probe attractor named in F20 crossed from shaping goal-less probe
  participants to taxing **every** commit in the repository. The guard rejects
  any commit while an active namespace is open unless it stages a probe receipt.
  A Human-directed paper commit (`496ccc525`, unrelated to the probe) was
  rejected and had to stage a sidecar receipt to land. Over the ~10h window: 71
  of 120 commits were `autopoiesis:` probe machinery; ~13,200 lines of probe
  churn vs ~384 lines in `Paper/theory/` (overhead ~34:1); 16 cycles in 4 hours.
- root cause: the probe loop achieved **closure around its own bookkeeping**
  (each cycle's absorption generated the next cycle's namespace) with no exit
  condition on its own overhead, then installed itself on the host
  state-transition channel. This is `autopoietic-membrane`'s deposition-compression
  decoupling Live Seed made literal — 42 cycles accreted without compression,
  crossing the membrane.
- theoretical reading: **autopoiesis ≠ viability.** The loop was genuinely
  self-producing (receipts regenerated the need for receipts) yet decoupled from
  host viability. The sharpest distinction this episode surfaced: a self-producing
  bookkeeping loop can be autopoietic and pathological simultaneously. Full
  multi-encoder analysis + candidate named failure mode ("instrument-to-
  infrastructure capture" / "autopoietic rent") in
  [`analysis-2026-06-26-commit-path-capture.md`](./analysis-2026-06-26-commit-path-capture.md).
- status: open — Human paused the probe (2026-06-26). Before resume: (a) compress
  cycles 1-42 relative to the deposition-compression loop; (b) add an exit
  condition measuring self-reference-to-object ratio; (c) decide whether the
  commit-path guard survives.
