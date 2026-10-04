# Main Paper Theory Expansion Discussion: April 2026

This note records the April 2026 theory-expansion discussion for the main
boundary paper. It is not a required revision checklist. It is a map of theory
directions that can strengthen the paper without inflating the current claim.

## Plain Central Thesis

The strongest theoretical claim is:

> In process-native artificial life, the reproduced self is not identical with
> source, body, or artifact. It is the minimal organization that preserves
> future addressability and action capacity across a mortality boundary under a
> given carrier ecology.

Short form:

> The future carrier selects what has to count as the self.

Current evidence supports this most clearly through the carrier contrast:

- raw-stream pressure can make a thin stream handler sufficient
- Quine-native `ctl/wake` pressure makes a control-surface-capable runtime
  organization necessary
- runtime-carried self-source is therefore an affordance, not a sufficient
  determinant of the reproduced self-object

## What The Paper Already Has

The current paper already has enough to claim:

- self-object is condition-dependent
- source visibility is not sufficient
- generic `exec` is not closure
- source rebuild is not closure
- delayed addressability after real `exec` is a stronger continuity test
- return to `idle` matters because it restores future addressability
- bounded control-surface closure is below full autopoiesis and above ordinary
  handoff or artifact copying

## Concepts To Preserve

Use these terms as stable paper vocabulary:

- **carrier-indexed self-object**: the reproduced self-object is defined
  relative to the future-task carrier
- **post-handoff addressability**: the successor remains reachable after real
  handoff
- **control-surface closure**: the successor can receive, interpret, act,
  report, and return to waiting state through the native control surface
- **viability-induced self selection**: constraints determine which
  organization is sufficient for continuation
- **scaffold gradient**: the path from explicit continuation instructions to
  more endogenous maintenance
- **process-native reproduction**: reproduction organized around process
  replacement, carrier continuity, and control interface, not file copying

Avoid over-strong uses:

- do not use `selection` as population-level evolutionary selection
- do not claim full autopoiesis
- do not claim full autonomous agency
- do not treat source visibility as selfhood

## Five-Layer Closure Frame

This is the most useful expansion structure.

| Layer | Question | Current status |
|-------|----------|----------------|
| Description closure | Can the system access and use a self-description? | `source-code/` can be inspected, but source alone is insufficient |
| Constructive closure | Can it build a viable successor body? | direct rebuild and fresh synthesis are available in bounded lanes |
| Carrier closure | Does the successor preserve or regenerate the future-task carrier? | central current result; `ctl/wake` creates a stronger requirement than raw stream |
| Control-surface closure | Can it wake, act with tools, report, and return to `idle` after real `exec`? | current strongest positive in `6.64` / `6.65` |
| Normative / autonomous closure | Does it maintain and repair future addressability with less scaffold? | future boundary, not current claim |

Use this frame to explain why the paper is strong without overstating it.

## Nearby Theory Lanes

These lanes can enrich the full paper or future work, but should stay
subordinate to the control-surface result.

### Semantic Closure / Epistemic Cut

Useful idea:

- `source-code/` is symbolic self-description
- build/runtime/toolchain interpret that description
- the rebuilt runtime re-enters a control ecology where it can interpret new
  tasks

Safe wording:

> Quine shows partial semantic recoupling, not full semantic closure.

### Closure Of Constraints

Useful idea:

- the system does not close every production condition
- it locally closes a control interaction lane
- provider, toolchain, source projection, and control-probe ecology remain
  supplied by the experimental environment

Safe wording:

> The result is local control-surface constraint closure, not closure of all
> constraints required for autonomous organization.

### Agency / Normativity

Useful idea:

- the next theoretical threshold is whether addressability becomes something
  the lineage protects and repairs, not merely something the prompt tells it
  to preserve

Possible term:

- **addressability normativity**

Safe wording:

> Current experiments show addressability as an imposed viability condition;
> future experiments should test whether the lineage internalizes it as a
> maintained norm.

### Process Theory / Behavioral Equivalence

Useful idea:

- the successor need not be byte-identical
- it must preserve the observable interaction trace required by the future
  control ecology

Candidate formal statement:

```text
P0 ≈_C P1 iff, under control ecology C,
P1 after exec can realize:
idle -> wake(task) -> tool/action -> report -> idle
```

### OS / Capability Ontology

Useful idea:

- reproduction is partly preserving or regenerating the capability topology
  needed for future control

Carrier topology categories:

- source carrier
- execution carrier
- future-task carrier
- tool carrier
- state carrier
- observer/probe carrier

### Philosophical Individuation

Useful idea:

- the individual is not a file, PID, or source tree
- it is a process individual formed through source projection, runtime rebuild,
  `exec` boundary, control surface, and idle ecology

Use this only as interpretation, not as the operational test.

### Minimal-Life Analogies

Chemoton and autocatalytic-set analogies can help explain why the result is
interesting, but they are dangerous if made central.

Safe use:

- structural analogy only
- not a claim that Quine is living
- not a claim of material metabolic closure

### Open-Ended Evolution

Current work is lineage continuity, not open-ended evolution. The path toward
open-endedness would be:

```text
single closure event
-> repeated closure cycles
-> closure under carrier perturbation
-> self-modified successor mission
-> population of carrier-adapted self-objects
-> new control-surface strategies
```

## Experiments Suggested By The Discussion

Prioritize:

1. **Successor mission ablation ladder**: remove explicit instructions for
   staying alive, calling `idle`, treating control payloads as tasks, using
   tools, and returning to `idle`.
2. **Multi-wake closure**: test several post-handoff `ctl/wake` tasks, not one.
3. **Carrier perturbation**: vary wake payload format, tool paths, state
   exposure, working directory, and source-code layout.
4. **Minimal runtime competitor**: allow a non-Quine `ctl/wake` adapter; if it
   succeeds, the self-object is control-surface-compatible organization rather
   than full Quine runtime.
5. **Source degradation**: remove or corrupt source projection pieces and test
   diagnosis / repair.
6. **Multi-generation recursion**: require P1 to generate P2 and P2 to pass a
   delayed wake test.

## Manuscript Use

For the current main paper:

- absorb the five-layer closure frame only if space allows
- use nearby theories to sharpen boundaries, not to broaden the title claim
- keep the main contribution centered on carrier-indexed self-object and
  post-handoff addressability
- keep the conclusion at bounded control-surface closure

For future work:

- successor-mission ablation and multi-wake closure are the cleanest next
  experiments
- addressability normativity is the most promising bridge from reproduction
  theory to agency theory
- carrier topology is the most promising bridge from ALife to OS/process
  theory
