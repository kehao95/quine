# Introduction

Computational systems routinely continue across boundaries. A process may call
`exec` and become a new program while inheriting open file descriptors. An
actor may crash and be restarted behind a mailbox. A service may fail over
behind a stable endpoint. A checkpointed computation may return after a long
suspension. In each case, there is a tempting handle that looks like identity:
the source file, the process identifier, the memory image, the endpoint name,
the mailbox, the checkpoint, or the observable behavior.

None of these handles is globally sufficient. Source can be copied without a
living successor. A process identifier can survive a complete program-image
replacement. A checkpoint can restore memory without restoring the future
interface by which obligations arrive. An endpoint can route to a newly
constructed organization rather than a preserved body. A behaviorally similar
clone can answer the same input while not continuing the same obligation chain.

This paper proposes **future addressability** as an operational criterion for a
narrower and more tractable class of cases: obligation-bearing computational
lineages. A successor is continuous with such a lineage only relative to a
future ecology in which later obligations can still reach it through the
relevant carrier and be acted on through the ecology's required trace. The
claim is not a universal metaphysics of digital identity. It is a practical
criterion for cases where a lineage is judged by future work.

The central idea is deliberately simple:

> computational individuality is indexed by future addressability.

The work is in making that sentence nontrivial. Future addressability is not
mere endpoint naming. It is a triadic relation:

```text
obligation + carrier + competent successor organization
```

The future obligation is what must still be done. The carrier is the channel,
mailbox, file descriptor, capability, endpoint, control surface, scheduler, or
restore protocol through which that obligation can arrive. The successor
organization is the minimal body, runtime, protocol handler, state image, or
control-capable organization that can receive the obligation and act on it.

Quine provides the motivating case. Quine realizes language-model lineages as
native POSIX processes [@ke2026quine]. POSIX gives execution and inheritance
primitives, but not offspring semantics [@ritchie1974unix]. In the Quine
experiments motivating this paper, source, executable body, process boundary,
future-task carrier, and continuity can come apart. With runtime-carried
self-source available, a raw-stream future can make a thin stream successor
viable. Under a Quine-native `ctl/wake` future, the viable successor must
reconstruct a runtime organization capable of receiving a wake event, acting
with tools, reporting, and returning to `idle`.

This contrast motivates the paper's thesis:

> When source, state, body, endpoint, and future-task carrier diverge, the
> continuing computational individual is the carrier-indexed organization that
> preserves post-boundary addressability, competence, and receptive-state
> restoration under the specified ecology.

This is the generalization layer of a three-part Quine paper stack. The
late-breaking abstract announces the bounded result, and the main boundary
paper carries the Quine evidence ladder. This paper uses that evidence as the
motivating case for a broader continuity criterion. It does not claim that
future addressability is all of computational identity, nor that Quine proves
generic self-replication.

The contribution is threefold. First, the paper separates identity handles
from continuity criteria. Second, it connects artificial-life reproduction
questions to concrete system carriers such as file descriptors, mailboxes,
capabilities, checkpoints, and control surfaces. Third, it turns the Quine
result into a general experimental program for distinguishing source rebuild,
ordinary handoff, state preservation, organizational reconstruction,
carrier/capability overlap, control-surface closure, substrate capture, and
stronger autonomous renewal.
