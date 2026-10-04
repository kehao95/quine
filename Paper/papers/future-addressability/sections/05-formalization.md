# A Lightweight Formalization

The paper does not need a full formal proof system in its first version. A
small predicate is enough to make the criterion auditable.

Let:

- `E` be a future-task ecology.
- `c` be the carrier through which obligations arrive in `E`.
- `tau_E` be the future-obligation trace required by `E`.
- `B` be a boundary event such as `exec`, restart, migration, checkpoint
  restore, or handoff.
- `P` be the predecessor organization.
- `S` be the candidate successor organization.

Then:

```text
P --B--> S

Continuity_E(P, S, c) iff:
  Addressable_E(S, c)
  and S satisfies tau_E
  and Adequate_E(S, P, c)
  and Receptive_E(S)
```

`Addressable_E(S, c)` says that the successor can be reached through the
carrier after the boundary. It is deliberately stronger than name resolution.
In the naming-theory decomposition of name, address, and route to a destination
[@saltzer1982naming], `Addressable_E` requires the route to terminate at a
*bound, competent* organization, not merely at a name that still resolves or a
host that is still up. This is also why a capability distinction is useful: a
stored *sturdy* reference that can be reactivated satisfies addressability
latently, while a currently connected *live* reference satisfies it directly
[@miller2006robust]. `S satisfies tau_E` says that the successor can perform the
required future-obligation trace. `Adequate_E(S, P, c)` says that the successor
is behaviorally adequate relative to the carrier and ecology, not necessarily
byte-identical or memory-identical; the natural formal reading is a
carrier-restricted coalgebraic bisimulation, equivalence over the observations
the carrier exposes rather than over a fixed global interface
[@rutten2000universal]. `Receptive_E(S)` says that the successor restores the
waiting condition needed for later obligations.

When more than one candidate satisfies the predicate, the criterion is
not silent. It inherits the closest-continuer tie-break [@nozick1981philosophical]:
the continuant is the candidate to which `E` actually routes the obligation
through `c`. A second identical build that no obligation reaches is not a
competing continuant. This keeps the predicate well-defined under branching
without smuggling in a stronger identity metaphysics.

For Quine:

```text
tau_Q = idle -> wake(task) -> tool/action -> report -> idle
```

The nearest formal tools can be named without making the whole manuscript a
formal-methods paper. Process calculi model channel mobility and handoff
[@milner1992calculus]. Capability semantics model authority-bearing carriers
[@miller2006robust]. Session types model communication traces
[@honda1998language]. Temporal logic models liveness [@pnueli1977temporal].
Coalgebraic equivalence and bisimulation model observer-relative behavioral
adequacy [@rutten2000universal].

The notation should stay honest about its limits. Real ecologies are
asynchronous. Capabilities are only as secure as their implementation. The
future-task ecology can change. Trace equivalence can miss branching behavior.
Full bisimulation can be too strong for practical systems. The predicate is
therefore a guide to experimental reporting, not a finished theorem.
