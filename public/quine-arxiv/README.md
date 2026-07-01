# Quine: Realizing LLM Agents as Native POSIX Processes

> Systems / runtime paper. Published as arXiv:2603.18030.

This entry is the self-contained public record of the Quine systems paper: the
manuscript source, figures, and the released PDF.

## What the paper says

Current LLM agent frameworks reimplement isolation, scheduling, and
communication at the application layer, even though mature operating systems
already provide these mechanisms. Quine instead realizes an LLM agent as a
**native POSIX process**: identity is the PID, interface is the standard streams
and exit status, state is memory / environment / filesystem, and lifecycle is
`fork` / `exec` / `exit`. A single executable implements the model by recursively
spawning fresh instances of itself. Grounding the agent abstraction in the OS
process model inherits isolation, composition, and resource control directly from
the kernel, while supporting recursive delegation, self-renewal through `exec`,
and shell-native composition. The paper also marks where the process model stops
and points to two extensions beyond process semantics: task-relative worlds and
revisable time.

## Canonical version

- **arXiv:** [2603.18030](https://arxiv.org/abs/2603.18030)
- **DOI:** [10.48550/arXiv.2603.18030](https://doi.org/10.48550/arXiv.2603.18030)
- **Released PDF (v1, 2026-03-08):** [`data/2603.18030v1.pdf`](data/2603.18030v1.pdf)

## Contents

| Path | What it is |
|---|---|
| [`sections/`](sections/) | Manuscript source (abstract, introduction, POSIX mapping, runtime, properties, related work, conclusion, appendix) |
| [`figures/`](figures/) | Architecture, Interface, and Lifecycle diagrams |
| [`references.bib`](references.bib) | Bibliography |
| [`data/2603.18030v1.pdf`](data/2603.18030v1.pdf) | The released arXiv PDF |
| [`metadata.yaml`](metadata.yaml) | Title / author / abstract metadata |

The runtime described by this paper is the code at the repository root
([`cmd/quine/`](../../cmd/quine/), [`internal/`](../../internal/)).

## Cite

```bibtex
@misc{ke2026quine,
  title        = {Quine: Realizing LLM Agents as Native POSIX Processes},
  author       = {Hao Ke},
  year         = {2026},
  eprint       = {2603.18030},
  archivePrefix= {arXiv},
  doi          = {10.48550/arXiv.2603.18030}
}
```
