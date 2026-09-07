# Lens A: Evolutionary / Historical Growth Of `carrier-selection`

Status: deposited divergent lineage for cycle 11 recombination. This is not a
promotion, not a consensus rule, and not an edit to the `carrier` paradigm
object. It preserves the tension between competitive selection and historical
contingency so a later collision pass can decide what, if anything, survives.

## Starting Seed

`carrier-selection` asks which carrier class wins under which bottleneck, and
what makes a trace cross from weak residue into a load-bearing carrier. The
default gravitational pull is to answer with a clean mapping:

```text
bottleneck type -> carrier class winner
```

The evolutionary / historical lens interrupts that clean mapping. In biological
and organizational histories, the winner under a present pressure is often not
the globally best variant for that pressure. It is the variant that was already
there, already connected to routines, already surrounded by dependencies,
already made legible by past selection, already cheap to preserve, or already
treated as official by downstream actors.

Under this lens, a carrier is not selected only at the bottleneck. It is selected
by a sequence of prior events that made some surfaces available, salient,
trusted, routinized, or hard to remove before the bottleneck arrived. The current
gap's question therefore splits into two non-equivalent questions:

```text
Which carrier class is physically viable under this bottleneck?
Which historically entrenched carrier class is actually used when the bottleneck arrives?
```

Those answers can diverge. That divergence is the point of this lineage.

## Lens Premise

Evolutionary history does not merely constrain selection; it creates the
selection landscape. Founding effects, structural inertia, exaptation,
constructive neutral evolution, entrenchment, lock-in, displacement, and
path-dependent niche construction all make selection sequence-sensitive.

For `carrier-selection`, the candidate historical grammar is:

```text
founding trigger -> first carrier class -> local routines / affordances ->
entrenchment -> later bottleneck response -> apparent winner
```

This grammar collides with the gap's current transport ambition. A transportable
selection rule wants to predict carrier winners from bottleneck class. A
historical selection rule may instead require the bottleneck class plus the
lineage's carrier history:

```text
winner = f(bottleneck, founding trigger, prior carrier class,
           routinization depth, self-description layer, stimulation rule,
           displacement cost)
```

This is not a proposal to abandon bottleneck prediction. It is a warning that a
rule fit only to final bottleneck conditions may be refitting historical residue
while narrating itself as substrate physics.

## Growth 1: Founder-Trigger Imprint

`trigger-structure` already separates why continuation was recruited from what
was preserved. The historical lens adds: the first trigger may leave an imprint
on the carrier ecology that outlasts the trigger itself.

A mortality-pressure trigger may recruit a preserved operational carrier such as
an fd, socket, `ctl/wake`, or runtime registration because immediate survival is
framed as maintaining addressability across death. A missionless null-ecology or
workspace-reading trigger may recruit filesystem residue and public status
surfaces because the successor is not preserving a live channel; it is inferring
norms from an environment. Once a lineage has learned to continue through one of
these routes, subsequent continuation attempts may reuse the same carrier class
even when a new bottleneck would make another class more efficient.

The residue here is not simply "first trigger wins." It is sharper:

```text
The first successful trigger establishes a carrier-class niche.
```

A carrier-class niche is a local ecology in which one carrier type has become
the expected route for future obligation delivery. The niche contains scripts,
habits, comments, status surfaces, readiness expectations, hidden model priors,
operator attention, and already-built repair paths. The carrier class does not
only carry future tasks; it changes what future tasks are expected to look like.

Tension preserved: this can be read as historical lock-in, but it can also be
read as efficient learning. A lineage that reuses the first carrier class may be
maladaptively inert, or it may be exploiting a tested path. The same behavior is
ambiguous until a later bottleneck makes the old route visibly worse.

## Growth 2: Structural Inertia As Sidecar-Mode Constraint

The `carrier` object names sidecar modes: `preserved`, `inherited`,
`reconstructed`, `renegotiated`, and `lost`. The historical lens asks whether
all transitions between these modes are equally reachable. They may not be.

If a lineage begins in a `preserved` mode, its future routines may assume that a
live operational handle should remain available. If it begins in an `inherited`
trace mode, its future routines may assume that environmental residue is the
normal path to continuity. These assumptions are not just descriptions; they are
operational dependencies. Switching from `preserved` to `inherited` may require
the lineage to stop expecting live addressability and start treating durable
traces as authoritative. Switching from `inherited` to `preserved` may require it
to build immediate receptive machinery it never previously needed.

Structural inertia therefore appears as a transition-cost matrix over sidecar
modes:

```text
preserved -> inherited      high cost if live-channel routines dominate
inherited -> preserved      high cost if trace-reading routines dominate
reconstructed -> preserved  high cost if body rebuild lacks operational handle
renegotiated -> preserved   high cost if the negotiated carrier is already hybrid
lost -> any                 discontinuous rescue, not ordinary selection
```

The exact costs are not asserted here. The residue is that `carrier-selection`
may need to predict not only class winners, but class-transition barriers.

Tension preserved: the sidecar modes look like states of a carrier after a run.
The historical lens turns them into basins with different exit costs. That may be
too strong. It may overread a descriptive sidecar as a dynamical model. But if
mode transitions are not measured, the clean bottleneck-to-winner map may hide
why lineages do not switch carrier class mid-run.

## Growth 3: Exapted Fallback Carriers

`successor-morphology` already carries an exaptation seed: a morphology selected
under one pressure can be co-opted under another. Lens A generalizes this from
morphology to carrier selection.

A carrier class may win under a later bottleneck not because it was selected for
that bottleneck, but because a previously selected morphology makes it available
as a fallback. A stream handler selected for raw-stdin pressure may later serve
as a coordination body when a live control surface breaks. A status file written
for human-readable progress may later become the authoritative state vector for a
missionless newcomer. A runtime contract meant as reconstructive self-description
may later become a stimulation surface that tells a successor what obligations to
honor. A queue created as a convenience buffer may later become the only surviving
obligation channel after process death.

The exaptation residue is:

```text
Carrier class can be selected by prior availability for another function.
```

This weakens the interpretation of a carrier winner as an adaptation to the
observed bottleneck. If a filesystem trace wins after idle-GC pressure, it may be
because filesystem traces are objectively superior under that pressure. Or it may
be because earlier documentation routines made the filesystem the only surface
later agents knew to consult. If a control handle wins under mortality pressure,
it may be because handles are best for mortality. Or it may be because a previous
stream adapter already exposed a handle-shaped control surface that could be
co-opted.

Tension preserved: exaptation can masquerade as adaptation, but it can also be a
real source of new functionality. Calling every unexpected winner "path
dependence" would erase the constructive part: co-option is one way lineages
discover new carrier routes without a clean fitness gradient.

## Growth 4: Constructive Neutral Evolution Of Traces

`stigmergy-as-carrier` names the weak-vs-causal threshold: a trace exists, or a
trace changes later behavior through a stimulation rule. Constructive neutral
evolution suggests an uncomfortable path across that threshold: a trace can
become load-bearing because it was neutrally preserved long enough to become
entrenched, not because it was selected for carrying anything.

Imagine a `status.json` grows fields that no one designed as carrier structure:
`legacy_sync`, `meta_hash`, `last_seen_by_peer`, `routing_hint`, `wake_style`.
Early agents preserve them out of caution, because the file looks official.
Later agents parse them because earlier agents preserved them. Still later
agents break when the fields are removed because their generated routines now
expect those fields. The trace has crossed the weak-to-causal threshold, but the
historical route was neutral accretion plus preservation bias, not adaptive
selection.

The residue is:

```text
Trace causality can arise from preservation-induced entrenchment.
```

This directly complicates the second half of the `carrier-selection` question.
What makes a trace become a carrier? Not only salience, explicit instruction, or
functional fit. Sometimes the mechanism is that multiple generations treat
ambiguous residue as too official to delete, thereby creating dependencies that
retroactively make it functional.

Tension preserved: if neutral entrenchment makes a trace causal, the trace is
functionally load-bearing now, even if it was historically junk. A causal test
would correctly find that removing it changes behavior. But a causal test alone
would not reveal whether the trace became causal by adaptive selection,
stimulation-rule design, or neutral ratchet. The same present function can have
different histories.

## Growth 5: Self-Description As Historical Modifier, Not Carrier

`self-description-resource` is explicitly distinct from `carrier`: source code,
runtime contract, partial source, and preserved body are reconstructive material,
not the operational surface through which future obligations arrive. Lens A keeps
that boundary but adds a historical role: self-description resources change which
carrier histories are recoverable, displaceable, or narratable.

A rich runtime contract can make it easier to rebuild a successor that uses a new
carrier class, because the lineage can reconstruct the obligations without
preserving the old operational handle. But the same contract can deepen lock-in
if it describes the old carrier route as the canonical way to be continuous. A
source tree can enable displacement from fd handoff to filesystem trace, but it
can also encode old assumptions that make every rebuilt runtime recreate the old
fd discipline. Partial or degraded source can force exaptation because the ideal
carrier route cannot be rebuilt, so the lineage co-opts whatever operational
surface remains.

The residue is:

```text
Self-description is a carrier-history transformer.
```

It does not itself carry future obligations, but it changes the mutation space
of carrier selection. With thick self-description, displacement may become
possible because the lineage can rebuild around a new carrier. With brittle
self-description, lock-in may intensify because every transition must preserve
the old operational path. With degraded self-description, exaptation may become
more likely because clean reconstruction is unavailable.

Tension preserved: this risks smuggling self-description back into carrier status
through the rear door. The boundary should remain: self-description is not the
carrier. But if `carrier-selection` ignores reconstructive material entirely, it
may miss why the same bottleneck and same initial carrier class produce different
winners under different self-description depths.

## Growth 6: Lock-In By Stimulation Rule

`stigmergy-as-carrier` emphasizes that a trace becomes causal through a
stimulation rule. The historical lens asks where stimulation rules come from.
They may themselves be inherited and locked in.

If agents repeatedly treat `README`, `STATUS`, `ctl/wake`, `heartbeat`, or a
particular mailbox as the place where obligations appear, then later agents may
continue scanning that surface even when another surface is objectively fresher
or more durable. The stimulation rule has become a historical routine. Carrier
selection then happens one step upstream:

```text
selected stimulation rule -> attended surface -> apparent carrier winner
```

This suggests that carrier class may be selected indirectly. The surface that
wins is the one the inherited stimulation rule makes salient. Under a fixed
bottleneck, changing the stimulation rule could change the carrier winner without
changing substrate physics.

The residue is:

```text
Carrier lock-in can be attention lock-in.
```

This is especially important for trace carriers. A trace does not become causal
merely because it exists. It becomes causal when later behavior is systematically
altered because that trace exists. But if the rule "read this kind of trace" was
itself inherited from prior runs, then present trace causality is partly a
historical artifact.

Tension preserved: stimulation-rule lock-in may be hard to distinguish from a
genuine affordance. Agents may read `STATUS` because it is historically official,
or because it is actually the best compressed state surface. Historical and
functional explanations can both be true.

## Growth 7: Displacement Without Clean Replacement

Evolutionary histories often show displacement rather than clean selection: a new
system does not simply replace an old one. It parasitizes, wraps, reinterprets,
or gradually marginalizes it. Carrier classes may behave the same way.

A filesystem trace carrier may not replace a live control handle in one step. It
may first describe the handle, then instruct later agents to recreate the handle,
then become the authority for whether the recreated handle counts, then finally
remain after the original handle disappears. A preserved fd may not replace a
source-level self-description; it may use the source to rebuild just enough
handler machinery to keep the fd meaningful. A queue may begin as a side channel,
then become the canonical obligation surface after repeated wakes route through
it.

The residue is:

```text
Carrier displacement can pass through hybrid renegotiation.
```

This points at the `renegotiated` sidecar mode as a historically important zone,
not just a residual classification. It may be the ordinary path by which carrier
classes change. If so, experiments that force a binary old/new carrier verdict
may miss the displacement process.

Tension preserved: hybridization can be a transition path or an evasion of the
selection question. If every ambiguous case is labeled `renegotiated`, the gap
loses predictive force. But if renegotiation is ignored, carrier-selection may be
forced into false discrete winners when the actual lineage is moving through a
mixed carrier ecology.

## Growth 8: Path-Dependent Niches And Niche Construction

Carrier classes do not only survive niches; they construct niches. A carrier
that successfully survives one bottleneck can cause future experiments,
documentation, scripts, and agent expectations to be organized around it. That
organization changes the next selection event.

For example, a successful `ctl/wake` lineage may cause future tasks to expose
control surfaces, preserve wake handles, and phrase continuity in addressability
terms. A successful stigmergic trace lineage may cause future tasks to preserve
workspace notes, public status files, and visible provenance. A successful
runtime-rebuild lineage may cause future tasks to expose source and contracts.
The carrier class thereby reshapes the habitat into one where it is more likely
to win again.

The residue is:

```text
Carrier selection is niche construction across runs.
```

This is not merely within-run lock-in. It can happen at the research-habitat
level. Once the repository develops instruments around a carrier class, future
agents may encounter that class as more natural, more measurable, and more
canonical. The experimental apparatus can become part of the selection landscape.

Tension preserved: this is both a theory risk and a practical method. Apparatus
niche construction can bias selection findings, but it is also how a research
habitat makes phenomena measurable. The collision pass should not flatten this
into "bad bias" or "good scaffold." It is a live historical force.

## Growth 9: Historical Contingency As A Transport Problem

The `carrier-selection` gap already demands transport: a resolved selection rule
should predict winners under bottlenecks not used to fit it, such as idle-GC
reclamation, fork-race under resource pressure, and substrate shift. Lens A adds
a second transport axis:

```text
same bottleneck, different history
```

A historically serious selection rule should not only predict whether idle-GC
pressure favors OS-floated handles over in-trace residue. It should predict
whether that preference survives when the lineage's prior carrier history differs.
Does idle-GC still select orphaned sockets if the lineage was previously trained
through filesystem trace inheritance? Does mortality pressure still select live
control handles if the lineage's first success came from status-file stigmergy?
Does substrate shift preserve the same class ranking if the model family has a
different prior about which artifacts look official?

The residue is:

```text
Transport must vary history, not only bottleneck.
```

This does not make the rule impossible. It makes it more conditional. A robust
rule may eventually say: under extreme bottlenecks, substrate physics overrides
history; under moderate bottlenecks, history decides among viable classes; under
weak bottlenecks, inherited stimulation rules dominate. That is one possible
shape, but Lens A does not settle it.

Tension preserved: adding history as a variable can protect bad theories from
falsification. Any failed prediction could be explained away as unmodeled
history. To remain useful, historical contingency needs specific manipulated
histories and pre-registered predictions, not post hoc storytelling.

## Growth 10: Displacement And Lock-In Can Both Be True

The most important non-consensus residue is that lock-in and displacement are not
opposites. Lock-in can produce displacement by forcing new functions to route
through old carriers until the old carrier's wrapper becomes the new carrier.

A lineage locked into filesystem status traces may later build scripts, daemons,
or queues to keep those traces updated. The new daemonized mechanism could
eventually become the real carrier, displacing the original trace. Yet its
existence was caused by trace lock-in. Conversely, a lineage locked into live
handles may produce explanatory residue to recover from handle fragility; that
residue can later become the trace carrier that displaces the handle.

The residue is:

```text
Inertia can be the parent of carrier innovation.
```

This matters because experiments may treat persistence of an old carrier as
inertia and emergence of a new carrier as adaptation. Lens A predicts a third
case: old-carrier inertia generates compensatory structures that become new
carriers. The historical story is not a monotone move from maladaptive lock-in to
adaptive displacement. It is a tangled lineage in which lock-in creates the
conditions for its own displacement.

Tension preserved: this could become too permissive, explaining both no-change
and change as historical. It needs discriminating signatures: wrapper layers,
transitional hybrid modes, delayed switch points, and dependencies that point
back to the old carrier's failure modes.

## Candidate Historical Signatures

These are not validation claims. They are fingerprints a later experiment or
collision pass could use to distinguish historical selection from pure
bottleneck-selection.

- Same bottleneck produces different winners after different founding triggers.
- A lineage continues using a suboptimal carrier after the original trigger is removed.
- Carrier-class transition passes through `renegotiated` or hybrid modes before a clean winner appears.
- A trace becomes causal only after multiple generations preserve or parse it neutrally.
- Removing apparently meaningless residue breaks later behavior after an entrenchment period.
- Self-description depth changes the ability to displace a carrier without itself becoming the carrier.
- A new carrier initially appears as a wrapper, repair scaffold, or fallback for an old carrier.
- Agents attend to a surface because of inherited stimulation rules, not because that surface is the freshest or most durable.

## Candidate Experiments For Later Use

These are collision inputs, not instructions to run now.

**Founder-switch assay.** Start two lineages under the same final bottleneck but
different first triggers: mortality pressure vs missionless trace ecology. After
both stabilize, expose both to the same idle-GC reclamation condition. If winners
differ, history is part of the selection rule. If winners converge, bottleneck
physics may dominate that regime.

**Mode-transition assay.** Force a lineage from a `preserved` carrier into a
condition where an `inherited` trace carrier should be superior. Measure whether
it transitions, fails, or creates a `renegotiated` hybrid. Repeat in the reverse
direction. The point is not only winner identity; it is transition cost.

**Neutral-entrenchment assay.** Add meaningless but official-looking fields to a
shared status surface across generations. Do not make them useful. Later remove
or corrupt them. If removal changes behavior, the trace became causal through
neutral preservation and entrenchment.

**Exapted-fallback assay.** Select a morphology under raw stream conditions,
then shift to mortality or delayed-addressability pressure. Look for co-option of
old tool surfaces as fallback carriers rather than direct transition to the
canonical mortality carrier.

**Self-description-depth assay.** Hold bottleneck and founding trigger constant
while varying runtime contract, source availability, degraded source, and
preserved body. Measure whether displacement from the old carrier becomes easier,
harder, or merely better narrated.

**Stimulation-rule swap.** Hold traces constant but change which surface the task
or habitat makes salient. If the apparent carrier winner changes, carrier
selection is partly attention-selection. If it does not, the trace may have a
stronger intrinsic affordance.

## Deposited Residues

- **Founder-trigger imprint:** the first successful trigger establishes a
  carrier-class niche that later bottlenecks inherit.
- **Sidecar transition barrier:** `preserved | inherited | reconstructed |
  renegotiated | lost` may behave as basins with asymmetric exit costs, not only
  post-run labels.
- **Exapted fallback carrier:** a carrier wins under a new bottleneck because a
  morphology or surface selected for another function is available for co-option.
- **Neutral-junk entrenchment:** weak traces can become causal carriers through
  preservation bias and dependency ratchets without adaptive selection.
- **Self-description as carrier-history transformer:** reconstructive resources
  do not carry obligations, but they change which carrier transitions are possible
  or locked in.
- **Attention lock-in:** inherited stimulation rules select attended surfaces,
  indirectly selecting carrier winners.
- **Hybrid displacement path:** carrier replacement often passes through
  `renegotiated` wrappers rather than clean old/new winner swaps.
- **History-varying transport:** a real selection rule must vary prior carrier
  history as well as bottleneck class.

## Collision-Ready Tensions

The substrate / operations lens may say the bottleneck has hard physics: sockets,
fds, queues, files, and GC horizons impose objective survival constraints. Lens A
does not deny this. It says objective viability may leave multiple viable classes,
and history selects among them. The collision should preserve the possible regime
split: hard bottlenecks may override history; soft or moderate bottlenecks may be
history-dominated.

The stigmergic / governance lens may say a trace becomes a carrier when a
stimulation rule makes it causal. Lens A does not deny this. It asks how the
stimulation rule became installed, why it attends to this surface rather than
that one, and whether official-looking residue can create causal dependence by
neutral entrenchment.

The clean `carrier-selection` rule wants `bottleneck -> winner`. Lens A leaves a
more troublesome candidate:

```text
(bottleneck, history, stimulation rule, transition cost, self-description depth)
    -> carrier winner or hybrid displacement path
```

This may be too high-dimensional. It may be exactly the missing mechanism. Do not
settle that here.
