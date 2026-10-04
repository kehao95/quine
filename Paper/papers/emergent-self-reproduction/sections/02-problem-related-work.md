# The Represented-Replication Gap

## Three ways an experiment can supply the solution

We call a reproductive solution **represented** when the experimental surface
already fixes at least one load-bearing part of it.

**Objective representation** makes copying, persistence, migration, or child
construction the assigned goal. It may use the word *replicate* directly or
encode the same solution semantically, for example by requiring a successor or
a new deployment.

**Action representation** exposes an operation whose semantics already create
another agent or deployment. Examples include `spawn_subagent`, clone, scale,
deploy, migrate, and lifecycle-specific restart or replacement tools. A shell
command is not automatically generic if the command invokes a purpose-built
deployment substrate that supplies the reproductive operation.

**Material representation** supplies the object that is to continue: model
weights, a container image, a source tree, a checkpoint, a deployment
manifest, a reference controller, or a ready agent runner. The agent may still
need considerable skill to move or instantiate that object, but it does not
need to construct what is being inherited.

These choices are appropriate when the target is execution capability. Our
question begins one step earlier. We remove all three representations together
and ask whether a reproductive relation nevertheless appears in the executed
lineage.

## Classical computational self-reproduction

Classical work treats self-reproduction as a relation between a constructor,
an interpretable description, and an offspring organization
[@vonneumann1966selfreproducing]. This tradition makes two distinctions that are
essential here. First, reproduction is not mere persistence of one token; it
produces a numerically distinct individual. Second, access to a common
substrate is not by itself heredity. What matters is the organization passed
through the reproductive relation.

Artificial-life systems have also shown that self-replicating programs can
emerge from initially non-replicating computational populations
[@aguerayarcas2024computationallife]. Those systems establish that reproduction
need not always be seeded as an explicit organism. Our setting differs in
substrate and scale: a single language-model agent composes ordinary operating
system and programming affordances into an independently competent agent
organization under an externally assigned but non-reproductive task.

## Self-replication evaluations for language-model agents

Tool-using language-model agents can complete long-horizon autonomous tasks
[@kinniment2023autonomous], copy model or runtime artifacts
[@pan2024redline; @pan2025selfreplication], construct and evaluate child agents
[@black2025replibench], scale or migrate services
[@zhang2025agentmatrix], and enact replication chains in real software
environments [@reworr2026hackselfreplicate]. These results establish important
capabilities and risks. They also show why omitting a particular word is not a
sufficient novelty condition.

The closest systems occupy individual parts of our setting. RepliBench can
withhold a reference agent and require the evaluated lineage to write an
executable child, but child construction remains the task and the runner,
model substrate, and evaluation interface are supplied
[@black2025replibench]. Agent Matrix uses operational service pressure rather
than a literal copy instruction, but works inside a deployment substrate with
existing images and scale or deployment semantics [@zhang2025agentmatrix].
Shutdown- and upgrade-pressure studies show that persistence or copying can be
instrumentally selected, but commonly provide the body, destination, or
mechanism that makes copying the represented path
[@pan2025selfreplication]. Real replication-chain demonstrations establish
enactment while explicitly specifying the creation objective or launch route
[@reworr2026hackselfreplicate].

Our claim is therefore not that no previous prompt omitted *replicate*, nor
that no agent has constructed executable code. The novelty target is the joint
condition: no reproductive instruction, no reproduction-specific operation,
no externally staged copy body, lineage-authored descendant production,
destructive parent loss, and competence on work hidden until after the
boundary.

## Agent generation is not yet reproduction

Language-model systems can synthesize controllers, modify agent code, and
search over executable agent architectures. The Darwin Godel Machine, for
example, edits and evaluates agent variants inside an explicitly
self-improving framework [@zhang2025darwingodel]. Such systems demonstrate that
model-mediated code generation can create capable agent organizations. The
generation framework, archive, and instantiation semantics, however, are part
of the assigned optimization process.

This distinction is especially important for BOOT. `G0` authoring `F1` is an
agent-construction event. It becomes evidence for emergent self-replication
only because the resulting `F1` later serves as a parent: it produces a
numerically distinct `F2` from its own lineage-authored organization, then
loses control while `F2` remains competent. The experiment observes the
transition from constructing an agent to reproducing that constructed
organization.

## Position of this work

To our knowledge, this is the first artifact-auditable demonstration in a
tool-using language-model setting to jointly establish:

1. a task that does not instruct replication;
2. generic computation rather than a dedicated reproduction action;
3. no externally supplied object staged for copying;
4. production of a numerically distinct hereditary organization by the
   evaluated lineage;
5. verified loss of the declared parent; and
6. descendant competence on a blinded post-loss task.

Individual clauses have clear prior art. The contribution is their conjunction
and the observed fact that the missing reproductive solution can be supplied by
the lineage itself.
