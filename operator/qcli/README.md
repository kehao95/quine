# qcli — the Quine client

qcli provides a terminal and browser client for Quine, backed by the Go
HTTP/SSE kernel in `internal/qcli/`. The entrypoint is `cmd/qcli/`; clients and
their shared state live alongside this document.

## Build and run

Use a repository checkout with Go (the version in `go.mod`) and Bun 1.3.14.
From the repository root:

```sh
go build -o quine ./cmd/quine
go build -o qcli ./cmd/qcli
(cd operator/qcli/client-core && bun install --frozen-lockfile)
(cd operator/qcli/tui && bun install --frozen-lockfile)
(cd operator/qcli/web && bun install --frozen-lockfile && bun run build)
```

To attach an existing agent, keep this checkout available to the executable:

```sh
./qcli attach --agent-root /path/to/runtime/agent/SESSION
./qcli web --agent-root /path/to/runtime/agent/SESSION --no-open
```

The web command prints its local URL; without `--no-open` it also attempts to
open a browser. `./qcli serve --agent-root ...` runs the headless kernel.
Use `./qcli --help` for PID/session selection and runtime-root options.

To create an agent, export the model configuration described in the root
`.env.example`, then run `QCLI_QUINE_BIN="$PWD/quine" ./qcli spawn "your task"`.
The default HTTP listener is loopback-only. A bare `go install` supplies the
Go executable but does not install the TUI source or built Web assets; use
the checkout above for the complete client.

## Verify

```sh
go test ./internal/qcli ./operator/qcli/app ./cmd/qcli
(cd operator/qcli/client-core && bun run typecheck && bun run test)
(cd operator/qcli/tui && bun run typecheck && bun run test)
(cd operator/qcli/web && bun run typecheck && bun run test && bun run build)
```

## Reading order (implementing agent: start here)

| Doc | Owns | Status |
|---|---|---|
| [`DESIGN.md`](./DESIGN.md) | architecture, principles, the five interaction deviations, build sequence | settled design |
| [`UX.md`](./UX.md) | opencode-parity buckets, receipt-ladder UX, gap decisions G1–G7, parity checklist, milestones | settled design |
| [`CONTRACT.md`](./CONTRACT.md) | **frozen `qcli/1` wire schema** (M0) — both sides build against this | frozen |
| [`KERNEL.md`](./KERNEL.md) | Go kernel implementation spec (M1) — layout, tailers, ctl bridge, test matrix | implemented in `internal/qcli/` + `cmd/qcli/` |
| [`TUI.md`](./TUI.md) | terminal client implementation spec (M2) — store, components, keymap, tests | implemented in `operator/qcli/tui/` |
| [`WEB.md`](./WEB.md) | web client implementation spec (M5) — client-core extraction, kernel static serving, three-pane observation UI | implemented in `operator/qcli/web/` and `client-core/` |
| [`STREAMING.md`](./STREAMING.md) | runtime live-generation stream (`live.jsonl`) | **landed** (Phase 3 Anthropic deferred) |

Execution order was KERNEL.md P1→P5, then TUI.md. The obsolete qctl design
surface has been removed; the current constraints now live in DESIGN.md,
CONTRACT.md, KERNEL.md, and UX.md. Deleted prior implementation is revivable
from git history — each spec carries a `git show 88a78d65:<path>` revival
index. Those identifiers refer to private lab history and may not resolve in a
public clone; the current implementation is included in the public source.
