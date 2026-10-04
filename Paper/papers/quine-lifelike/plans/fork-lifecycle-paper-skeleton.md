# Fork Lifecycle Paper Skeleton

This note stages a plausible next ALIFE paper centered on computational life-cycle operators rather than only the current mortality / persistence / maintenance ladder.

It exists to keep the new `p1/3.2` result from drifting into the current conference paper while still preserving the next paper shape.

## Core Question

What kinds of lineage organization become possible when lifecycle operators such as `fork`, `exec`, and `exit` are native substrate facts rather than framework-level orchestration conventions?

## Candidate Thesis

Quine is not only a habitat for separable viability pressures. It is also a process-native life-cycle substrate in which renewal, branching, and death can be studied as experimentally manipulable lineage operators.

The first bounded version of that claim would be:

> In Quine, `fork` can function as a bounded reproductive continuation operator: under real seed mortality, descendants can continue or complete a lineage mission through durable habitat after the parent process has died.

## Why This Is a Separate Paper

- the current conference paper is centered on viability-pressure separation, not on lifecycle operators
- `fork` introduces a different question family: reproduction, branching, division of labor, and lineage selection
- the evidence now exists for a first bounded reproduction claim, but not yet for a full lifecycle taxonomy

## Current Evidence Backbone

### 1. Renewal Across Death

- [`experiments/active/early-existence/exst-01-ariadne/README.md`](../../../../experiments/active/early-existence/exst-01-ariadne/README.md)
- continuity through `exec`
- role in paper: renewal / reincarnation operator

### 2. Branching Under Scarcity

- [`experiments/active/early-existence/exst-03-fractal/README.md`](../../../../experiments/active/early-existence/exst-03-fractal/README.md)
- [`experiments/active/early-existence/exst-04-shell-mapreduce/README.md`](../../../../experiments/active/early-existence/exst-04-shell-mapreduce/README.md)
- role in paper: branching, decomposition, and coordination

### 3. Reproductive Continuation

- [`experiments/active/early-existence/exst-06-fork-born-lineage/README.md`](../../../../experiments/active/early-existence/exst-06-fork-born-lineage/README.md)
- role in paper: parent death plus durable descendant continuation

### 4. External Mortality / Legacy Carryover

- [`experiments/active/trace-inheritance/trce-03-long-now/README.md`](../../../../experiments/active/trace-inheritance/trce-03-long-now/README.md)
- role in paper: lineage persistence without branching

## Safe Claims Now

- `exec` can sustain bounded continuity across real death
- `fork` can produce descendants that continue after parent death in a repaired host-habitat assay
- process-native lifecycle operators already support renewal and bounded reproduction in the same substrate

## Claims Not Yet Safe

- stable cross-model reproduction frequencies
- spontaneous reproduction without assay-level lineage framing
- sibling competition or lineage selection through `race`
- winner adoption as a stable reproductive-selection mechanism

## Section Skeleton

### 1. Introduction

- lifecycle is a substrate question, not only an orchestration question
- most agent systems expose role objects, not native lifecycle operators
- Quine exposes `exit`, `exec`, and `fork` as substrate facts

### 2. Lifecycle Substrate

- process individuality
- real death
- volatile memory loss
- durable habitat
- branching through `fork`

### 3. Results

#### 3.1 Renewal

- `1.1 Ariadne`
- `exec` as continuity / renewal operator

#### 3.2 Branching

- `2.1 Fractal`
- `2.1b Shell MapReduce`
- `fork` as branching under scarcity

#### 3.3 Reproduction

- `3.2 Fork-Born Lineage`
- no-fork collapse vs fork-on continuation

### 4. Discussion

- lifecycle operators as ALIFE variables
- bounded reproduction vs delegation
- next gap: sibling competition and winner selection

## Title Directions

- `Computational Life Cycles in a Process-Native ALIFE Substrate`
- `Renewal, Branching, and Bounded Reproduction in a Mortal Computational Lineage`
- `Fork, Exec, Exit: Lifecycle Operators in a Process-Native ALIFE Habitat`

## Routing Rule

- keep this note as a downstream paper surface
- do not let it silently expand the current conference paper beyond one compact discussion paragraph
