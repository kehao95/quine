# Quine Notes

> Short system-semantics notes that extend the core Quine paper without reopening the whole architecture each time.

Paper-level control-plane home:

- dossier: [`../README.md`](../README.md)

## Role in the System Track

The runtime-substrate family has three layers of writing:

1. **Core paper family** — see [`../../quine-arxiv/README.md`](../../quine-arxiv/README.md) and [`../../quine-hotos/README.md`](../../quine-hotos/README.md)
2. **Quine Notes** — compact papers on one primitive, one claim, one structural consequence
3. **Track B** — later synthesis once broader implementation and evaluation exist

These notes assume the core mapping:

> Quine models an agent as a native POSIX process: mission via `argv`, material via `stdin`, deliverable via `stdout`, diagnostics via `stderr`, outcome via `exit status`, with lifecycle through `fork` / `exec` / `exit`.

Each note should restate that mapping briefly, then isolate one consequence.

## Series Constraints

- **one note, one claim**
- **3–5 pages by default**
- **self-contained but not self-rebuilding**
- **implementation-aware, experiment-light unless the note truly needs data**
- **series label stays light**: use “Part of the Quine Notes series” in subtitle, footnote, or project page metadata

## Master Manuscript Plan

- [`manuscript-plan-series.md`](manuscript-plan-series.md)

## Planned Notes

### Wave 1 — Interface Semantics

1. [`manuscript-plan-01-cognitive-harvard-architecture.md`](manuscript-plan-01-cognitive-harvard-architecture.md)
2. [`manuscript-plan-02-stderr-semantic-gradient.md`](manuscript-plan-02-stderr-semantic-gradient.md)
3. [`manuscript-plan-03-exec-context-reset.md`](manuscript-plan-03-exec-context-reset.md)
4. [`manuscript-plan-04-exit-delegation-verdict.md`](manuscript-plan-04-exit-delegation-verdict.md)

### Wave 2 — Composition and Coordination

5. [`manuscript-plan-05-stdin-stdout-compositional-interface.md`](manuscript-plan-05-stdin-stdout-compositional-interface.md)
6. [`manuscript-plan-06-filesystem-information-medium.md`](manuscript-plan-06-filesystem-information-medium.md)
7. [`manuscript-plan-07-fork-wait-race.md`](manuscript-plan-07-fork-wait-race.md)

### Wave 3 — Data / Learning

8. [`manuscript-plan-08-tape-exit-synthetic-trajectory-generation.md`](manuscript-plan-08-tape-exit-synthetic-trajectory-generation.md)

## Suggested Writing Order

1. Cognitive Harvard Architecture
2. stderr as Semantic Gradient
3. exec as Context Reset
4. exit as Delegation Verdict
5. stdin/stdout as the Compositional Interface

The remaining notes should stay queued until the first wave is real.
