---
id: maintainability
type: principle
gist: prefer clear structure over clever shortcuts; converge on one canonical surface
status: canonical
---

# maintainability

- Prefer clear structure over clever shortcuts
- Update existing control-plane docs instead of spawning orphan files
- Keep naming, responsibility, and compression level aligned; when one surface starts carrying multiple jobs, split it
- Avoid near-duplicate configuration surfaces; keep one canonical variable unless there is a strict operational need
- Similar or overlapping runtime capabilities should converge on one canonical control surface; mirrors or projections need a sharply distinct role or they are drift
- Primitives and tools should have one clear meaning; do not overload one interface with multiple unrelated transformations when those responsibilities can be separated cleanly
- **Code is status**: implementation facts should not be duplicated in prose
