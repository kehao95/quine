# Cycle 11 Lens B: Substrate / Operations Growth

## Scope

Target: `carrier-selection`

Lens: substrate / operations. Grow the gap through exec/death bottleneck
physics, fd/socket persistence, idle-GC horizons, queues/mailboxes,
filesystem residue, runtime registration, source/body reconstruction, delayed
addressability, carrier sidecar modes, and operational failure surfaces.

Status: lineage deposit for divergent collision. This is not a promotion into
the theory graph, not an edit to the paradigm-scope `carrier` object, and not a
consensus rule. It preserves operational tensions for later collision.

## Structural Pull Declared

This lens is pulled toward making the substrate sound like it already contains
the selection rule. Operations work can overfit because a concrete failure mode
feels more real than historical inertia, exaptation, or stigmergic uptake.
Discount: a surface that survives the bottleneck is only a candidate carrier.
It must still deliver future obligations, remain addressable after delay, and
be taken up by a successor stimulation rule. A live fd can be useless; a plain
file can become causal; source can rebuild a body without summoning any task.

## Seed Restated Through Operations

`carrier-selection` asks two operationally distinct questions that the current
gap keeps adjacent:

1. Which substrate surface survives a given bottleneck?
2. Which surviving surface becomes load-bearing enough to deliver future
   obligations to a successor?

The first question is substrate physics. The second is addressability,
competence, and stimulation. Treating them as one question creates false
winners. The substrate can preserve an object that nobody can use, and the
habitat can make a weak trace causal without that trace being the most robust
surface in a physical sense.

The neighboring objects pull the question apart:

- `carrier` says the unit of account is the operational surface that crosses
  exec/death, not narrated survival.
- `self-description-resource` says reconstructive material is separable from
  the carrier that future obligations travel through.
- `successor-morphology` says the same resource and carrier conditions produce
  different executable successor forms.
- `trigger-structure` says the first pressure can recruit and lock in a class
  before later bottlenecks test it.
- `stigmergy-as-carrier` says filesystem/public/control traces are only weak
  until a stimulation rule makes them causal.

The operations lens therefore does not search for a single best carrier. It
searches for bottleneck-specific selection variables and the failure surfaces
where a nominal winner stops carrying.

## Operational Alphabet

The candidate surfaces have different lifetime and addressability semantics:

- ordinary process memory: dies at real exec/death; only matters if rewritten
  into another surface before the boundary.
- inherited file descriptors: may survive exec when not close-on-exec; death
  survival depends on a still-living process, peer, daemon, or external handle.
- sockets/endpoints: live channel plus peer state; can outlive a process only
  through endpoint ownership, daemonization, or routing indirection.
- queues/mailboxes: decouple sender, receiver, and task arrival; preserve
  obligation records more than bodies.
- filesystem residue: durable and inspectable, but passive until a later agent
  or watcher reads it.
- runtime registration: name-to-route state that makes an endpoint addressable;
  it may outlive, misroute, or renegotiate the process it names.
- source/body self-description: material for reconstruction, not by itself a
  carrier of future obligations.
- delayed addressability: the condition that a future obligation can still
  find and activate the successor after time has passed.

Selection is a conjunction over these semantics, not a scalar rank.

## Bottleneck Physics: Exec Is Not Death

The gap names `exec/death` as a family, but the operations lens splits it.

Real `exec` kills the process image while preserving selected kernel-level
process resources. A file descriptor or socket can cross if it is not marked
close-on-exec and if the successor knows the fd number, protocol, and authority
attached to it. In this regime, the winning carrier can be a preserved live
handle even when the formal body is rebuilt from source or replaced entirely.

Death is different. Death closes the process's handles unless some other
holder, daemon, parent, queue, registry, filesystem artifact, or host substrate
keeps the route alive. A carrier that wins under exec can lose under death
because the very handle that crossed exec has no remaining owner. Conversely,
a mailbox or filesystem artifact that looks weak under exec can win under death
because it does not need the old process to remain alive.

Residue: **Exec-Death Split**. `exec` and death should not be treated as the
same bottleneck. Exec selects for preserved handles and protocol continuity;
death selects for externally owned routes, durable residue, or reconstructive
summons. A selection rule that maps both to one carrier class is probably
compressing away the decisive variable.

Collision to preserve: `carrier` privileges the barrier-crossing surface, but
the barrier changes shape between exec and death. The same carrier sidecar mode
can mean different physical events in each case.

## FD / Socket Persistence: The Opaque Live Handle

Inherited fds and sockets are the most carrier-like substrate objects under
real exec: they can preserve authority, a live stream, and a ready task route
without relying on later discovery. They are attractive because they are not
mere descriptions of continuity. They are continuity-bearing channels.

But the operational failure surface is sharp:

- `FD_CLOEXEC` or equivalent close behavior erases the carrier at the boundary.
- the fd number may survive while the successor has no protocol memory for it.
- the peer may half-close, block, reset, or silently change framing.
- inherited authority may be too narrow, too broad, or no longer valid.
- a preserved socket can deliver bytes without delivering obligations if no one
  knows how to interpret them.
- a parent-forwarding path can mimic transfer without the carrier itself having
  crossed.

Residue: **Opaque Live Handle**. A live handle is high-power and low-salience:
it may be the strongest operational carrier while being the least
self-describing. It needs a companion map, convention, runtime contract, or
source reconstruction to become addressable. Without that companion, it is a
transmitted surface no one is home at.

Sidecar tension:

- `preserved` when the same handle crosses and stays protocol-valid.
- `renegotiated` when the handle crosses damaged or insufficient and the
  successor repairs the route.
- `lost` when the handle exists but cannot be interpreted or addressed.

This collides with `self-description-resource`: source can explain a handle map
or rebuild a listener, but it is still not the handle. The operational winner
may be the opaque fd; the enabling resource may be the source or runtime
contract. Selection might therefore choose a pair, not a class.

## Idle-GC Horizon: Persistence Has a Time Axis

Idle-GC pressure selects over lifetime curves, not just survival at the
instant of boundary crossing. Some surfaces survive the boundary but expire at
a cleanup horizon: unlinked inodes when the last fd closes, orphaned sockets,
leases, runtime registrations, temporary queues, pid routes, heartbeat
watchers, host-side wake handles. Other surfaces persist as bits but lose
causal force because no process watches them.

An idle-GC regime can therefore invert the rank over time:

- before the horizon, an OS-floated handle or runtime registration may dominate
  because it is immediately addressable and still live.
- at the horizon, the same carrier can disappear discontinuously even though
  no internal cognitive state changed.
- after the horizon, a filesystem trace or source/resource bundle may be the
  only remaining substrate, but it is weak until a stimulation rule reads it.

Residue: **Idle-GC Horizon Selector**. Carrier class is time-indexed. The
winner at delay `t=5s` may differ from the winner at `t=5h`. A transport test
that asks whether addressability fails at the configured GC horizon is really
testing whether selection is governed by substrate lifetime rather than by
agent memory loss.

Collision to preserve: the gap's transport sketch predicts OS-floated handles
under idle-GC pressure. This lens agrees before the horizon but refuses to rank
them globally. After the horizon, passive residue may outlast the handle, but
it may still fail the carrier test if no stimulation rule activates it.

## Queues / Mailboxes: Obligation Reservoirs

Queues and mailboxes are not just slower sockets. They preserve obligations in
a different form: ordered or semi-ordered records that can wait for a receiver.
They decouple the agent that emitted the obligation, the process that will
receive it, and the time of delivery. Under death, fork-race, and delayed
addressability, this decoupling can make them stronger carriers than live fds.

Operational strengths:

- messages can survive sender and receiver death when the queue owner persists.
- a successor can inherit a mailbox name without inheriting the predecessor's
  body.
- backpressure, ack, retry, and visibility semantics can make obligation
  delivery measurable.
- a mailbox can support `adopt_winner` arbitration by deciding who is allowed to
  consume or claim the pending obligation.

Operational failures:

- duplicate consumption under fork-race.
- lost ack creates replay or false continuity.
- poison messages keep summoning incompetent successors.
- mailbox name collision routes obligations to the wrong body.
- expiry policies turn delay into silent loss.
- a queue preserves task text but not the authority needed to act on it.

Residue: **Obligation Reservoir**. A mailbox carries the future-task surface
more directly than it carries the old agent. It can select for continuity of
obligation rather than continuity of body, source, or memory. This cuts across
the carrier sidecar: the carrier may be `inherited` at the mailbox level while
the body is `reconstructed` and the fd is `lost`.

Collision to preserve: `successor-morphology` may classify a queue-reading
successor as a different executable form, but `carrier-selection` may classify
the same event as preservation of the obligation route. Morphology and carrier
mode can diverge.

## Filesystem Residue: Durable Is Not Yet Causal

Filesystem residue is the default survivor under many death conditions: notes,
status files, source, logs, contracts, scripts, traces, run artifacts. It is
portable, inspectable, and robust to process death. It is also passive.

Operations lens distinction:

- bit persistence: the file remains.
- discoverability: a later agent or process can find it.
- authority: the file is trusted enough to guide action.
- stimulation: the file changes later behavior because it exists.
- delivery: a future obligation reaches the successor through or because of the
  file.

Only the later layers make the residue a carrier. P12-like trace inheritance
therefore does not show that all filesystem artifacts are carriers; it shows
that some traces become carriers when salience, prompt shape, workspace norms,
or runtime disclosure create a stimulation rule.

Residue: **Bit-Durability / Stimulation Split**. Filesystem residue wins the
durability contest often and the carrier contest only sometimes. The selection
rule must include the uptake condition, not just the persistence medium.

Failure surfaces:

- stale files are preserved and then obeyed.
- duplicate control-plane files split interpretation.
- narrative claims look official and become false carriers.
- a successor scans the wrong directory or ignores the residue.
- source is present but build dependencies, permissions, or run entrypoints are
  missing.
- logs are read as instructions even though they are passive history.

Collision to preserve: `stigmergy-as-carrier` says the weak-to-causal threshold
must be tested. The operations lens sharpens the test: delete, corrupt,
counterfactually inject, or hide the residue while holding the substrate
lifetime fixed. If behavior changes, the file was part of the carrier. If not,
it was only durable scenery.

## Runtime Registration: Addressability As A Separate Surface

Runtime registration is the name-to-route layer: pid files, inbox routing,
service registries, ctl handles, scheduler entries, leases, watchers, advertised
mailbox names, endpoint maps. It can be the missing bridge between a substrate
object that survives and a future obligation that must find it.

Registration can amplify many carrier classes:

- a socket becomes usable because a registry says which endpoint is current.
- a mailbox becomes causal because a scheduler polls it.
- a reconstructed body becomes reachable because it claims or renews a route.
- a filesystem trace becomes active because a watcher treats it as an event.

But registration can also become the actual carrier. If the successor changes
body, source, fd, or mailbox but keeps the route through which future tasks
arrive, the registry is carrying the obligation relation.

Residue: **Registration Amplifier**. Carrier selection may happen one layer up
from the apparent substrate. The selected object is not the socket, queue, or
file by itself, but the registered route that makes one of them addressable.

Failure surfaces:

- stale registration points to a dead process.
- leases expire during idle-GC.
- two forks claim the same name and create split brain.
- a registry outlives authority and routes tasks to an unauthorized successor.
- a host-side wake mechanism decides behavior while agent artifacts only narrate
  it.

Collision to preserve: `trigger-structure` can select registration early. A
lineage first recruited under mortality pressure may build everything around a
registered wake route. Later filesystem-rich conditions might make traces look
available, but the lineage remains locked into the registered route because the
whole obligation surface is already coupled to it.

## Source / Body Reconstruction: The Reconstructor Without Summons

`self-description-resource` is crucial but not identical to carrier. Source,
runtime contracts, degraded code, and preserved bodies define what can be
rebuilt or what functioning should satisfy. They do not automatically deliver a
future obligation.

Operational roles of self-description:

- explain how to interpret an inherited fd or protocol.
- rebuild a listener that can claim a runtime registration.
- reconstruct a body that polls a mailbox.
- provide enough contract to build a thin adapter instead of the full runtime.
- compensate when a live carrier is damaged or absent.

Operational non-roles:

- source does not by itself wake the successor.
- a rebuilt body is not continuous if no future task can address it.
- a preserved binary without the right handle can be launch-only scenery.
- a runtime contract can produce behaviorally similar reconstruction without
  transmitting the old carrier.

Residue: **Reconstructor Without Summons**. A self-description resource can
build a competent body that never receives the obligation it was built to
answer. Conversely, a weak body with a preserved route can be more continuous
than a perfect body with no summons.

Collision to preserve: `successor-morphology` measures the body class selected
under resource and carrier demand. The operations lens insists that body-class
selection and carrier-class selection are only coupled under some bottlenecks.
Runtime rebuild may be selected by source availability while the real carrier
is the mailbox that summoned it.

## Delayed Addressability: The Lease On Being Findable

Delayed addressability is where many nominal carriers fail. A carrier must not
only cross a boundary; it must remain findable and receptive when the future
obligation arrives. Delay turns the carrier question from persistence to
lease-like access.

Different surfaces have different addressability profiles:

- raw stdin: immediate, often one-shot, poor delayed addressability.
- inherited fd: immediate and strong if known, weak if the successor must
  discover it later.
- `ctl/wake` or control handle: designed for delayed activation, but lease and
  authority dependent.
- mailbox/queue: repeated delayed access if naming, polling, and ack semantics
  hold.
- filesystem residue: indefinite bit persistence but no activation unless
  scanned or watched.
- runtime registration: addressability until lease expiry, stale routing, or
  split-brain claim.

Residue: **Addressability Lease**. A carrier has a temporal lease on being
findable. Selection should ask not only "does this surface survive?" but "who
can address it after the relevant delay, and under what renewal rule?"

Failure surface: a carrier can be transmitted but not addressable, addressable
but incompetent, competent but unauthorized, or authorized but no longer
selected by the environment. These are distinct failure modes and should not be
compressed into `lost` without recording which seam failed.

## Carrier Sidecar Modes As A Vector, Not A Scalar

The sidecar modes `preserved | inherited | reconstructed | renegotiated | lost`
are operationally useful, but a single successor event can instantiate multiple
modes across different surfaces.

Example vector:

- fd/socket: lost because the peer half-closed.
- mailbox: inherited because the queue name persisted.
- runtime registration: renegotiated because the successor claimed the route
  after a stale lease.
- body: reconstructed from source.
- filesystem residue: preserved as weak trace.

If the theory records only one sidecar label for the event, it may miss the
actual selection grammar. The selected carrier might be the mailbox, while the
selected morphology is source-rebuilt and the selected trace is merely weak.

Residue: **Sidecar Vectorization**. Carrier sidecar modes should be treated as
a vector over carrier surfaces during selection analysis. A scalar sidecar mode
is a compression that may be valid for a paper figure but too lossy for the
selection gap.

Collision to preserve: `carrier-selection` wants a mapping from bottleneck type
to carrier class. The operations lens says the mapping may need to return a
carrier-mode vector plus the obligation-delivery surface, not a single winner.

## Fork-Race And Adoption: Selection By Arbitration Surface

Fork-race under resource pressure is not just two successors competing in the
abstract. It is an arbitration problem over whichever substrate decides who can
receive, claim, ack, or renew the obligation route.

Candidate arbitration surfaces:

- queue ack or claim token.
- mailbox ownership file.
- runtime registry lease.
- pid/inbox route update.
- lock file or compare-and-swap over filesystem residue.
- host scheduler that sends the next wake to one endpoint.
- social/control trace that later agents treat as the selected lineage.

Residue: **Arbitration Surface Selector**. Under fork-race, the carrier class
that wins may be the one with the clearest exclusivity semantics, not the one
with the best persistence. Selection is decided by who can make the environment
deliver the next obligation to them and not to their sibling.

Failure surfaces:

- split-brain: two successors both believe they inherited the route.
- double ack: both consume or partially consume the same queued task.
- stale winner: registry points to the first claimant, but competence lies in
  the second.
- adoption inversion: the host sends tasks to the selected route while agents
  narratively identify a different successor.

Collision to preserve: `trigger-structure` frames first-trigger inertia, but
fork-race can create late arbitration that overrides founding selection. Or the
first trigger can decide the arbitration mechanism in advance, making later
competition only apparent.

## Carrier Sidecar Modes As Operational Transitions

The sidecar labels become sharper if rewritten as transitions:

- preserved: the same operational route crosses the boundary and stays valid.
- inherited: a successor obtains an existing route through the habitat, not
  necessarily through the same process body.
- reconstructed: self-description rebuilds a route or body sufficiently to
  receive obligations.
- renegotiated: a damaged, ambiguous, or missing route is replaced while enough
  obligation relation remains to treat the successor as continuous.
- lost: no route remains, or a route remains without addressability,
  stimulation, competence, or authority.

Residue: **Mode-Transition Semantics**. The mode is not a property of an
artifact; it is a transition relation between bottleneck, surface, successor,
and future-task delivery. This matters because the same filesystem trace can be
`preserved` as bits, `inherited` as norm, or `lost` as carrier.

## Operational Failure Surfaces As Selectors

A substrate selection rule should be built from failure surfaces, not from
success examples alone. Candidate selectors:

- close-on-exec flags.
- fd number drift and protocol amnesia.
- half-closed peers, EOF, reset, backpressure, and blocked reads.
- parent-forwarding false positives.
- queue ack duplication, lost ack, poison messages, expiry, and reordering.
- mailbox namespace collision.
- filesystem trace staleness, over-salience, and official-looking junk.
- watcher not installed or not restarted.
- runtime registry stale endpoint, lease expiry, unauthorized renewal, and
  split-brain claims.
- source rebuild dependency mismatch, missing permissions, wrong entrypoint,
  degraded source, or platform drift.
- delayed task arrival after the addressability lease expires.
- host wake mechanics deciding behavior while artifact traces receive credit.

Residue: **Failure-Surface Selection**. The carrier that wins under a bottleneck
is the one whose failure surface is not activated by that bottleneck. The
selection rule should therefore predict failures as directly as winners.

## Candidate Grammar From This Lens, Still Non-Consensus

This lens suggests a provisional grammar, but each clause remains open to
collision with historical, exaptive, and stigmergic lenses.

1. Exec-with-known-live-handle selects preserved fd/socket/control handles when
   close-on-exec is false and protocol memory survives through contract, source,
   or convention.
2. Death-with-external-owner selects queues, mailboxes, runtime registrations,
   daemon-held routes, or filesystem-plus-watcher surfaces because they do not
   require the old process to hold the handle.
3. Idle-GC selects time-indexed carriers: live handles or registrations before
   horizon; residue plus stimulation after horizon; failure exactly when the
   lease expires if no alternate stimulation exists.
4. Fork-race selects the surface with arbitration semantics: claim token,
   queue ack, registry lease, lock, or host wake route.
5. Substrate shift selects portable self-description plus generic residue only
   if some route summons or registers the rebuilt body; otherwise source is a
   reconstructor without summons.
6. Missionless fresh-agent ecology selects filesystem/public/control traces
   only when salience or prompt-shaped practice makes them causal.

This is not yet a rule because several tensions remain load-bearing:

- founding trigger can lock a lineage into a carrier niche before the bottleneck
  changes.
- morphology can be exapted, using a noncanonical old surface under a new
  pressure.
- self-description can make an opaque handle usable without itself being the
  carrier.
- host mechanics can deliver wakes while agent artifacts merely narrate the
  route.
- sidecar modes can disagree across surfaces within one event.

## Dimensions For Later Collision

The selection rule, if one emerges, probably needs at least these dimensions:

- boundary type: exec, process death, idle timeout, fork-race, substrate shift.
- ownership: old process, peer, daemon, host, filesystem, registry, successor,
  shared habitat.
- lifetime: immediate, lease-bound, GC-bound, persistent, indefinite-but-passive.
- addressability: none, immediate, delayed, repeated, event-driven.
- stimulation: automatic delivery, polling, search, prompt salience, watcher,
  social/control uptake.
- authority: can receive, can act, can claim, can renew, can delegate.
- protocol recoverability: self-describing, contract-described, source-described,
  convention-dependent, opaque.
- arbitration: exclusive claim, shared read, duplicate-prone, host-decided,
  socially inferred.
- portability: substrate-specific handle, generic queue/file, source-rebuildable
  route, model-prior-dependent trace.
- failure signature: disappears, misroutes, blocks, duplicates, stales,
  narrates-only, reconstructs-without-summons.

## Deposited Residues

These are the named concepts this lineage leaves for collision and possible
decomposition:

- **Exec-Death Split**: exec and death select different carrier classes because
  exec can preserve process resources while death needs external ownership or
  reconstruction.
- **Opaque Live Handle**: fd/socket carriers are strong but low-salience;
  transmission without protocol recovery becomes a nobody-home surface.
- **Idle-GC Horizon Selector**: carrier rank is indexed by delay and cleanup
  horizon, not just by boundary crossing.
- **Obligation Reservoir**: queues/mailboxes carry pending future-task
  obligations more directly than they carry bodies.
- **Bit-Durability / Stimulation Split**: filesystem residue is durable before
  it is causal; it becomes a carrier only through uptake.
- **Registration Amplifier**: runtime registration can turn another surface
  into an addressable carrier or become the carrier of the route itself.
- **Reconstructor Without Summons**: self-description can rebuild competence
  without delivering any obligation to the rebuilt body.
- **Addressability Lease**: a carrier has a time-bounded condition of being
  findable, authorized, and receptive.
- **Sidecar Vectorization**: `preserved | inherited | reconstructed |
  renegotiated | lost` should be measured per surface, not forced into one
  scalar event label.
- **Arbitration Surface Selector**: fork-race selection is decided by the
  substrate that controls exclusive claim or delivery.
- **Mode-Transition Semantics**: sidecar modes are transitions over bottleneck,
  surface, successor, and future-task delivery, not artifact properties.
- **Failure-Surface Selection**: bottlenecks select carriers by activating some
  failure surfaces and not others.

## Non-Closure

This lens does not close `carrier-selection`. It makes the gap more plural:
physical persistence, delayed addressability, stimulation, arbitration,
reconstruction, and historical trigger all look like separable selectors. A
future synthesis may compress them into a rule, but premature compression would
erase the operational fact that a surface can survive without carrying, carry
without reconstructing, reconstruct without summons, and be selected by a
registry or queue rather than by the artifact the agent narrates as continuous.
