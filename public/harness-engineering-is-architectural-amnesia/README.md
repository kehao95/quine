# Harness Engineering Is Architectural Amnesia

> Stop rebuilding the operating system in user space. A blog essay — published on Substack.

This entry is the self-contained public record of the essay: the manuscript
source and its figure, alongside the canonical published version.

**Read it:** <https://kehao95.substack.com/p/harness-engineering-is-architectural>

## What it's about

The series opener on runtime surface. The AI industry has noticed that runtime
matters — and then set about rebuilding, in application code, the isolation,
composition, and lifecycle machinery an operating system already provides. The
essay argues that agent harnesses keep re-deriving the OS from scratch, and asks
what changes when you inherit those primitives from the kernel instead.

## Contents

| Path | What it is |
|---|---|
| [`manuscript.md`](manuscript.md) | The essay source |
| [`figures/`](figures/) | The process-IO diagram the essay embeds |

The runtime this essay is about is the code at the repository root
([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Cite

```
Hao Ke. "Harness Engineering Is Architectural Amnesia."
kehao95.substack.com, 2026.
https://kehao95.substack.com/p/harness-engineering-is-architectural
```
