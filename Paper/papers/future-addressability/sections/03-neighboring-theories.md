# Neighboring Theories

Future addressability is not proposed in a vacuum. It composes pieces that
nearby theories already own. The contribution is the composition, so the
neighbors must be named precisely enough to show what is borrowed and what is
new. We group them into five families and then make explicit the two that do
the most load-bearing work: a metaphysics of continuity (closest-continuer) and
a theory of naming (name, address, route).

## Self-reference and construction

| Neighbor | What it gives | What remains different here |
|---|---|---|
| Self-reproducing automata | self-description, description copying, universal construction [@vonneumann1966selfreproducing] | delayed future-obligation testing after a boundary |
| Recursion theorem | the logical root of self-reference: a program with access to its own description [@kleene1952introduction] | self-description as a reconstructive resource, not as the continuity criterion |
| Computational reflection | a running system that inspects and rebuilds its own structure [@smith1984reflection; @maes1987concepts] | reflection across a real process-image boundary, judged by post-handoff carrier viability |
| Digital life platforms | substrate-defined organisms, genomes, offspring placement, selection [@ray1991approach; @ofria2004avida] | cases where the substrate does not predefine the reproduced self-object |

The recursion theorem and reflection matter because they are the formal home of
"a system that holds and acts on its own description." Quine's source projection
and runtime rebuild are reflective operations across a destructive `exec`. The
novelty is not the reflection; it is that reflection alone does not settle
continuity. A reflective rebuild can still fail the future carrier.

## Naming, addressing, and authority

| Neighbor | What it gives | What remains different here |
|---|---|---|
| Naming and binding | the canonical split of name, address, and route to a destination [@saltzer1982naming] | which bound organization must exist for a future obligation to be served, not just resolved |
| Identifier/locator separation | stable identity under a changing locator [@farinacci2013lisp] | identity indexed by served future obligations, not by a registered identifier |
| Capability systems | authority-bearing references; sturdy versus live references [@miller2006robust; @shapiro2003capability] | authority as one carrier in the topology; a sturdy reference is latent addressability |
| POSIX process semantics | concrete `fork`, `exec`, file descriptors, image replacement [@ritchie1974unix] | individuality criteria under future-task carriers |

Naming theory is what keeps "addressability" from collapsing into either
endpoint naming or live reachability. A name resolves to an address; an address
routes to an attachment point; an attachment point must be backed by a competent
organization. Future addressability is a claim about that last binding: the
obligation must reach an organization that can serve it, not merely a name that
still resolves. Capability theory supplies the cleanest temporal version of
this: a *sturdy* reference can be stored and later reactivated, while a *live*
reference is currently connected. Latent addressability is the sturdy-reference
case; constant live reachability is the over-strong reading the criterion
rejects.

## Concurrency and behavioral identity

| Neighbor | What it gives | What remains different here |
|---|---|---|
| Actor model | mailbox addressability, message passing, replacement behavior [@hewitt1973actor; @agha1986actors] | actor identity as one carrier-indexed case, not the whole theory |
| Virtual actors | location-transparent identity reactivated on demand on a new host [@bykov2011orleans] | the engineered industrial form of latent, carrier-indexed reactivation |
| Futures and promises | explicit forward-looking obligations [@friedman1978applicative; @liskov1988promises] | obligation plus carrier plus competent successor organization |
| Process calculi | channel mobility and name passing [@milner1992calculus] | a lightweight criterion tied to real system boundaries |
| Session types | structured communication traces [@honda1998language] | a way to describe `tau_E`, not a full proof burden |
| Temporal logic | liveness properties [@pnueli1977temporal] | checking that later obligations are eventually acted on |
| Coalgebra and bisimulation | observer-relative behavioral equivalence [@rutten2000universal] | equivalence indexed by carrier and ecology, not over a fixed interface |

Virtual actors deserve emphasis because they are almost exactly the engineered
version of this paper's claim: a grain has a stable identity, no guaranteed live
instance, and is reactivated on whatever host can serve the next message. That
industrial pattern shows the criterion is not exotic; it is already how
location-transparent systems decide "the same individual."

## Systems continuity and reconfiguration

| Neighbor | What it gives | What remains different here |
|---|---|---|
| Checkpoint/restart and migration | state preservation across host or time boundaries | distinguishing state restoration from obligation-bearing addressability |
| State-machine replication | a service identity maintained by a replicated role, not one body [@schneider1990implementing] | a single lineage reconstructing its own carrier, not a quorum maintaining a service |
| Membership reconfiguration / leader transfer | admitting a successor to carry a role across a view change [@ongaro2014raft] | successor admission as the moment future addressability transfers |

This family is the closest systems analogue for **successor admission**: the
question of which body inherits the future-work carrier after the current one is
gone is a reconfiguration question, not only an `exec` question. The difference
is that replication maintains a role through redundancy, while future
addressability asks what must be reconstructed when there is no standing quorum
to fall back on.

## Persistence, individuality, and organizational closure

| Neighbor | What it gives | What remains different here |
|---|---|---|
| Closest-continuer theory | a successor counts as continuant if it is the closest continuer along the relevant dimension [@nozick1981philosophical] | the relevant dimension is the future-task carrier, made operational and testable |
| Psychological continuity and fission | survival as continuity, and the branching problem for multiple successors [@parfit1984reasons] | obligation routing, not memory, as the continuity relation; carrier arbitrates branching |
| Interactor / replicator and lineages | lineage individuality; the entity that interacts with its environment [@hull1980individuality] | the interactor as carrier-indexed organization rather than a biological body |
| Scaffolded reproduction | reproduction supported by external machinery, on a gradient [@godfreysmith2009darwinian] | the scaffold gradient as the boundary between bounded closure and stronger claims |
| Biological individuality | physiological and immunological criteria for one persisting individual [@clarke2010problem; @pradeu2012limits] | continuity criteria for computational, not metabolic, individuals |
| Markov blanket / free energy | a statistical boundary that individuates a self-maintaining system [@kirchhoff2018markov] | a carrier-defined boundary of receptivity, not a stationary statistical blanket |
| Autopoiesis and agency | strong upper bounds: intrinsic normativity, adaptive individuality, self-production [@maturana1980autopoiesis; @barandiaran2009defining] | bounded computational continuity below full life or agency |
| Closure of constraints and semantic closure | organization and symbol/action relations [@montevil2015biological; @pattee1995evolving] | scaffolded, process-native recoupling rather than full closure |

## Two load-bearing neighbors

**Closest-continuer theory** is the metaphysical backbone the criterion has been
implicitly using. Nozick's schema says a later object is the continuant of an
earlier one when it is the closest continuer along whatever dimension is
relevant, provided it is close enough. Future addressability fixes that
dimension and makes it testable: the relevant continuer is the organization to
which the ecology's native future obligations are actually routed and which can
serve them. This is what lets the criterion reject promiscuous readings (not
every reachable endpoint) without retreating to a single fixed handle (not the
PID, not the source). It also exposes the **branching problem** directly. When a
lineage forks into several runnable successors, Parfit's fission case applies:
there may be no fact of "the one successor" prior to the ecology routing
obligations. The carrier is what breaks the tie. A copy that no future
obligation can reach is not a competing continuant, however identical.

**Naming theory** is what makes the central word, *addressability*, do real
work. Without the name/address/route distinction, "still addressable" slides
between "the name resolves," "the host is up," and "a competent successor is
bound and ready." The criterion needs the third. This is also why a stable
endpoint name is explicitly insufficient: a name can resolve to a freshly built
organization, a stale one, or nothing competent at all.

## The boundary

The nearest theories reach address, state, behavior, authority, description,
construction, reconfiguration, persistence, and organizational closure. The
novelty boundary is their composition around future obligations. When familiar
identity handles diverge, future addressability asks which organization can
still be reached and made to act by the ecology's native future-task carrier.

This is why the claim should stay narrow. It is not "addressability defines
life." It is not "every endpoint is an individual." It is:

> for obligation-bearing computational lineages, future addressability is an
> operational continuity criterion.

That wording leaves room for biological individuality to require stronger
conditions such as metabolic self-maintenance and intrinsic normativity. It
also leaves room for purely mechanistic accounts of computation. Obligations
in this paper are not private intentions. They are concrete protocol
expectations: queued messages, promises, wake events, probes, scheduler work,
or control-surface tasks.
