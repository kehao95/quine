# Future Addressability

The paper uses five terms.

**Future-task ecology** is the environment that supplies later obligations and
specifies what counts as acting on them. The ecology may be a shell stream, an
actor mailbox, a service endpoint, a checkpoint scheduler, a capability-bearing
object reference, or a Quine control surface.

**Carrier** is the channel through which future obligations reach a candidate
successor. A carrier can be a file descriptor, a mailbox address, a network
endpoint, an object capability, a queue, a trace context, a scheduler slot, or
a control event. In capability systems, addressability and authority can be
closely linked: possessing a reference is also possessing the authority to use
it [@miller2006robust; @shapiro2003capability].

**Future-obligation trace** is the required interaction pattern. It is not just
an output behavior. It includes the receptive condition, the arrival of a later
obligation, action, reporting, and restoration of the state in which further
obligations can arrive.

**Carrier-indexed self-object** is the organization that must be preserved or
reconstructed to satisfy that trace under the specified carrier. It is not
necessarily the source, executable, memory state, endpoint name, or PID.

**Post-handoff addressability** is the delayed test. A successor must be
reachable after the boundary, not merely prepared before the boundary.

The criterion is:

```text
successor S is continuous under ecology E via carrier c iff,
after boundary B,
S remains addressable through c,
S satisfies the future-obligation trace tau_E,
S is behaviorally adequate relative to c,
and S restores the receptive condition for later obligations.
```

This criterion deliberately excludes several false positives. A copied source
tree is not enough. A successfully built binary is not enough. A PID that
survives `exec` is not enough. A stable endpoint name is not enough. A
functionally similar clone is not enough unless it is the organization to which
the relevant future obligation is actually routed.

The criterion also avoids an overly strong live-reachability requirement. Some
systems are intentionally dormant, air-gapped, address-rotated, or reachable
only through scheduled authority. These cases require **latent
addressability**: the successor is not continuously reachable through ordinary
routes, but remains reachable under the future ecology's delayed, authorized,
or scheduled carrier. Capability theory already draws this distinction as
*sturdy* versus *live* references: a sturdy reference can be stored and later
reactivated, while a live reference is currently connected [@miller2006robust].
Latent addressability is the sturdy-reference case, and a virtual-actor grain
that is reactivated on demand on a new host is its engineered form
[@bykov2011orleans].

The key contrast is therefore not "reachable versus unreachable" in the
abstract. It is: reachable how, by whom, for what obligation, after what
boundary, and with what competence?

## State Preservation And Organizational Reconstruction

The criterion clarifies a familiar systems split. Some continuity mechanisms
preserve state. Checkpoint/restart and process migration capture memory,
execution context, sockets, and related state so that an entity resumes with
minimal visible disruption. Other mechanisms reconstruct organization. Actor
supervision, service failover, and many cloud control planes allow local
execution state to die while preserving a role, mailbox, event log, source
image, or declarative control-plane state from which a competent organization
is rebuilt [@agha1986actors].

Future addressability cuts across both regimes. In state-preserving regimes,
the self-object may be the saved image plus restore protocol. In
organization-reconstructing regimes, the self-object may be the mailbox-bound
behavior, the service role behind an endpoint, the capability-bearing object,
or the wake/tool/report/idle runtime organization. The point is not that one
regime is more real than the other. The point is that the future carrier tells
us which organization matters.
