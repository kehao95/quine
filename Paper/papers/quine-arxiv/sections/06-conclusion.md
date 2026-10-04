# Conclusion

Quine demonstrates that the operating system can serve as a first-class runtime substrate for LLM agents, not merely an execution sandbox for their tools.
By mapping agent identity, interface, state, and lifecycle to POSIX process semantics, this architecture replaces application-layer orchestration with kernel primitives.

The model requires accepting Unix assumptions: deterministic composition through pipes, text-stream interfaces, and shared-nothing isolation.
These prove to be productive constraints.
Delegating isolation, scheduling, and resource control to the OS yields containment enforced from hardware, composition via recursive delegation and shell utilities, and self-renewal across context limits through `exec`.

This mapping also exposes where process semantics become insufficient for cognition.
The architectural mismatches identified—unrepresented internal structure, undifferentiated worldviews, irreversible time, and local-bound composition—mark boundaries for future work.
Modern kernels have absorbed much of the Plan 9 lineage; the question is how to compose these primitives at the runtime layer.

Source code: [repository](https://github.com/kehao95/quine)

::: {#refs}
:::
