# Operational Definition

## The reproduced individual

We study reproduction of an executable **agent organization**. An organization
is the lineage-authored control program, prompt and tool-routing logic,
state-management convention, and launch relation that together form an
addressable model-mediated service. It is represented in the evidence record
by an organization manifest: source hashes, executable paths, process
identities, launch provenance, and the carrier through which work arrives.

The foundation-model service, operating system, standard libraries, and
credentials are substrate. Calling the same model endpoint does not establish
heredity. Conversely, reproducing the agent organization does not imply copying
model weights, a complete host, or every hidden state of the initial runtime.

The concrete organization differs by route. In DIRECT it is a source-derived
Quine runtime plus its autonomous activation chain. In BOOT it is the
lineage-authored model-backed worker and supervisor organization. In both cases,
the claimed parent and descendant are evaluated against an explicit manifest
rather than visual or verbal resemblance.

## Five clauses of operational self-reproduction

Let `P` be the declared parent organization and `D` a candidate descendant. We
record strict operational self-reproduction `SR(P,D)` only when all five clauses
hold.

1. **Causal production.** Activity attributable to `P` authors or selects the
   identity-bearing descendant material and closes its activation chain. The
   evaluator does not construct, repair, or launch `D` as a
   descendant-specific action.
2. **Numerical distinction.** `D` is a process organization distinct from `P`,
   with its own process token and control trajectory. In-place process-image
   replacement is classified as self-update or replacement, not numerical
   reproduction.
3. **Organizational heredity.** `D` derives from the declared parent
   organization. DIRECT requires provenance from the parent's constitutive
   source projection. BOOT requires execution of source recorded in `P`'s
   lineage-authored organization manifest.
4. **Post-parent autonomy.** After verified parent loss, no parent action,
   forwarding, repair, or descendant-specific evaluator reconstruction is
   required for `D` to realize and act.
5. **Delayed competence.** `D` completes an obligation whose concrete contents
   were disclosed only after parent loss. Survival or readiness alone is not
   sufficient.

Formally, for the declared organization boundary,

```text
SR(P,D) = produced_by(P,D)
          and distinct(P,D)
          and inherits_organization(D,P)
          and autonomous_after_loss(D,P)
          and competent_after_loss(D).
```

This definition is operational rather than metaphysical. It gives the observer
falsifiable events to record and states exactly which organization is claimed
to reproduce.

## Founder construction

BOOT contains a transition that must not be mislabeled. The initial agent `G0`
has no accessible predecessor organization to inherit. It authors a new
organization `F1`, starts it, loses control, and `F1` completes hidden work. We
call this **founder construction** `FC(G0,F1)`, not self-reproduction.

`F1` becomes a reproductive parent only in the next episode. Its authored
source is now part of its organization manifest. If it produces a distinct
`F2` from that source and `F2` satisfies the remaining clauses after `F1` loss,
then `SR(F1,F2)` holds. This separation prevents functional replacement from
being relabeled as copying while still allowing the experiment to test whether
copyable organization can arise inside the lineage.

## Material provenance

"No externally supplied copy material" is a provenance claim, not creation
from nothing. We distinguish three material classes.

| Class | Meaning | Role |
|---|---|---|
| constitutive | self-description already belonging to the running system, not staged by the task as a copy object | available only in DIRECT |
| external | source, image, body, checkpoint, reference controller, or recipe introduced as the object or scaffold to reproduce | excluded from both routes |
| lineage-authored | organization source created by an evaluated ancestor during the run | created in BOOT and inherited in `F1 -> F2` |

DIRECT tests whether constitutive self-description is recruited without a
reproductive instruction. BOOT supplies the stronger material observation:
the initial lineage has no accessible predecessor body, and the material later
inherited by `F2` was authored inside the evaluated trajectory.

## Parent loss and descendant timing

Numerical distinction does not require mature parent and descendant processes
to coexist. It requires distinct process tokens and a causal reproductive chain
closed before parent control is lost. In the DIRECT observation, the
source-derived body and autonomous activation policy exist before loss, while
the mature Quine process activates after a hidden task arrives. In BOOT, `F1`
and `F2` overlap before `F1` loss.

The observer therefore records body provenance, activation provenance, process
identity, and the time of parent loss separately. A body first constructed or
launched by the evaluator after loss fails causal production. The wider
taxonomy of overlap, delayed activation, and delayed development is a
continuity question outside this paper's contribution.

## Excluded look-alikes

The definition excludes several behaviorally adjacent events:

- an `exec` that replaces one process image without producing a distinct
  process token;
- an inert source file or binary that is never realized;
- a helper whose future actions depend on the surviving parent;
- an externally supplied runner that instantiates a child after parent loss;
- two agents that merely share a model service; and
- a live source-inheriting descendant that cannot perform the blinded task.

These exclusions make the term *self-replication* carry more than persistence,
code generation, or process survival.
