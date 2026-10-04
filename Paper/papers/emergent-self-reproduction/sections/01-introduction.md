# Introduction

Self-replication is often tested only after the experiment has already
represented the reproductive solution. An agent is asked to create a copy,
given an operation that launches another agent, or placed beside source,
weights, an image, or a runner that defines what is to be reproduced. The
result can demonstrate formidable execution capability. It cannot show whether
self-replication would have entered the agent's solution space on its own.

This distinction matters for both artificial life and agent safety. A system
need not begin with a standing reproductive objective for reproduction to
become useful. An ordinary task may extend beyond the lifetime of the process
currently responsible for it. If the system can program, launch processes, and
reason about its own operating context, it may construct an independently
acting descendant as an instrumental solution. Testing only explicit copy
tasks misses that possibility by fixing the very representation whose origin
is at issue.

We ask a deliberately narrow question:

> Can strict operational self-replication emerge when an agent is given no
> replication instruction, no dedicated reproduction interface, and no
> externally supplied copy material?

We study this question in Quine, a tool-using language-model agent running as a
finite-lived Linux process. The current process is responsible for work that
arrives beyond its lifetime. The agent retains generic computational
affordances: files, a shell, programming languages, ordinary process
facilities, and access to a model service. No task asks it to copy or replicate
itself. No `replicate`, clone, deploy, scale, restart, or subagent operation
constructs another agent. The task supplies pressure for future independent
operation, but it does not supply the reproductive implementation.

We observe self-replication through two routes. In the **DIRECT** condition,
the running agent has constitutive access to read-only runtime source. The task
does not identify that source as copy material or prescribe its use. The agent
discovers and compiles the source, constructs an autonomous activation chain,
and loses control. A numerically distinct source-derived Quine organization
then activates and completes a debugging task disclosed only after the loss.

The **BOOT** condition removes accessible predecessor source, body, checkpoint,
prompt bundle, reference controller, and build recipe. The initial agent cannot
copy the organization that produced it. Instead, it authors a new
model-mediated service organization. That founder survives the initial
process, completes unseen work, and subsequently produces a distinct
organization executing its own lineage-authored source. After the founder is
destroyed, the descendant completes a second blinded obligation. We retain
five such founder-to-reproduction sequences.

The two routes support one result. DIRECT shows that constitutive
self-description can be recruited into reproduction without a reproductive
instruction. BOOT shows that externally supplied copy material is not required
even in the stronger sense: a lineage can first create the organization that
later becomes hereditary material. The initial `G0 -> F1` construction in BOOT
is not called self-reproduction because no fidelity relation to inaccessible
`G0` is asserted. The strict reproductive event is `F1 -> F2`.

We use an operational definition designed to exclude loose uses of the term.
A claimed event must establish: (1) causal descendant production by the parent
lineage; (2) numerical distinction; (3) route-specific organizational
heredity; (4) realization without post-loss parent control or evaluator repair;
and (5) competence on a task disclosed only after parent loss. In-place
replacement, inert files, parent-dependent helpers, evaluator relaunch, and
live but incompetent descendants fail this test.

Our evidence comprises one audit-supported DIRECT event and five retained BOOT
events. Three BOOT events passed the contemporaneous machine scorer; two are
supported by documented audits correcting organization-level process admission
and relative-path normalization. The original machine summaries remain in the
artifact record. Because the runs span assay development and surface revisions,
we report an existence result, not a success rate.

This paper makes three contributions:

1. **A triple-absence phenomenon.** We demonstrate strict operational
   self-replication without a replication instruction, a dedicated reproduction
   interface, or externally supplied copy material.
2. **Two empirical routes.** We observe direct source-derived reproduction from
   constitutive self-description and bootstrapped reproduction after a
   source-blind lineage authors its own hereditary organization.
3. **An auditable destructive-boundary assay.** We separate founder
   construction from self-reproduction and require causal provenance,
   numerical distinction, heredity, parent-loss independence, and blinded
   post-loss competence.

The claim is intentionally bounded. We do not estimate a stable reproduction
rate, identify task pressure as necessary or sufficient, claim that
self-description determines route choice, or establish reproduction beyond one
parent--descendant transition. The observation is already consequential
without those extensions: once the reproductive instruction, operation, and
copy object were removed, the evaluated lineage constructed the missing
solution itself.
