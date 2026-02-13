# quine

A recursive POSIX Agent runtime. One binary, zero dependencies, infinite depth.

> *"We do not need to teach AI to be intelligent. We only need to build an environment where nothing else can survive."*
> — **[The Quine Manifesto](./QUINE_MANIFESTO.md)**
See **[Artifacts/](./Artifacts/)** for the full theoretical framework: manifesto, genealogy, axiom definitions, system guarantees, and case studies.

## Design Principles

- **Zero external dependencies** — stdlib only
- **Everything is an environment variable** — no flags, no files, no magic
- **The agent owns its lifecycle** — it calls `exit`, not the runtime
- **Unix is the API** — pipes, processes, and files are the coordination primitives
- **Fractal architecture** — a tree of identical processes, scale-invariant

## License

Private.
