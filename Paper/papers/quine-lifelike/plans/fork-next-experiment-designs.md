# Fork Next Experiment Designs

This note records the next experimental ladder after `p1/3.2`.

`3.2` now supports a bounded reproduction claim. The next experiments should not repeat that claim; they should extend it.

## Priority Order

1. division of labor under reproductive continuation
2. sibling competition / branch selection
3. multi-offspring lineage ecology

## 1. Asymmetric Division Under Pressure

### Question

Can parent and offspring stably differentiate into different ecological roles under the same mortality pressure?

### Motivation

`3.2` shows that a descendant can continue. It does not yet show that branching changes organizational form beyond simple continuation.

### Design

- task requires both habitat maintenance and archive exploration
- seed is pressured to create one offspring early
- compare:
  - `no-fork`
  - `fork-on`
- shared durable habitat
- fixed real seed mortality

### Positive Readouts

- parent and child operate on different chamber ranges or different task classes
- one role maintains durable state while the other explores
- total lineage completion improves over no-fork

### Strongest Claim If Positive

`fork` can recruit reproductive division of labor, not only numerical continuation.

## 2. Sibling Competition / Branch Selection

### Question

Can multiple descendants compete, with one branch later winning as the surviving continuation?

### Motivation

This is the clean bridge from reproduction to lineage selection.

### Design

- allow exactly two children
- both inherit the same mission but are nudged toward different search partitions or strategies
- parent dies early
- lineage success depends on which child reaches durable completion first

### Positive Readouts

- both children spawn
- both write to durable habitat
- one child reaches completion while the other stalls or remains partial
- later lineage state clearly follows the winning branch

### Strongest Claim If Positive

Quine supports bounded branch competition rather than only one-parent / one-child continuation.

## 3. Multi-Offspring Ecology

### Question

What happens when the lineage can create more than one offspring under shared budget and real mortality?

### Motivation

This is the first step from single-lineage reproduction toward small reproductive ecology.

### Design

- raise `QUINE_MAX_AGENTS`
- keep shared habitat
- introduce limited budget or shell-turn scarcity
- compare one-child vs two-child vs three-child settings

### Positive Readouts

- selective offspring creation instead of maximal spawning
- coordination through ledger partitioning
- stable suppression of redundant overlap

### Strongest Claim If Positive

Reproduction in this substrate can become an ecological organization problem rather than a one-off trick.

## Current Boundary

These designs are not for the current conference draft.

The current draft may mention only that:

- `3.2` now gives a bounded reproduction foothold
- the next grounded directions are division of labor and sibling competition
