# qcli-server Kernel — Implementation Specification (M1)

> **Status: ready to execute.** Implements the [`CONTRACT.md`](./CONTRACT.md)
> server side. A fresh agent can execute this end-to-end; every runtime claim
> is backed by a `file:line` reference (§9) or a `git show` revival pointer
> (§8) — do not re-investigate them. Read `DESIGN.md` §Design Principles first;
> a change that makes the kernel *drive* the agent, or *persist* anything, is
> wrong regardless of how useful it looks.

## 1. Goal

A stateless Go daemon that (a) resolves or spawns one Quine agent, (b) opens a
client endpoint the agent can reply to, (c) tails the agent's observation
surfaces into typed contract events over HTTP+SSE, and (d) writes graded
control events to the agent's ctl surface with bounded latency. Restart at any
time must lose nothing durable.

## 2. Background — runtime facts (verified 2026-07-02, do not re-derive)

1. **The ctl surface is FUSE-only.** `runtimeSurfaceBackendName = "fuse"`
   (`internal/runtime/public_surface.go:6`; decision:
   `Paper/core/decisions/2026-06-control-surface-fuse-only.md`). The agent's
   own process serves `public/ctl/{post,poke,inject,interrupt}` as virtual
   nodes (`public_surface_fuse_linux.go:271-412`); a write commits on
   flush/release. Consequences: a healthy agent accepts a plain `os.WriteFile`;
   a **wedged** agent can block the write indefinitely inside the FUSE call; a
   **dead** agent without unmount yields `ENOTCONN` (stale mount) or `ENOENT`.
2. **Ctl writes carry plain text** (payload after trailing-newline
   normalization, `contract_surface.go:81`), not JSON. The client identity /
   reply path travels in the text envelope (CONTRACT.md §6).
3. **Runtime root layout**: `pid/<PID>` is a symlink to the agent's `public/`
   dir; `agent_root = dirname(public)`; `agent/<session_id>/` is the agent
   root; `client/<id>/` is client-endpoint territory
   (`peer_discovery.go:203-262`; old qctl `Resolve`).
4. **Observation surfaces**: `current.jsonl` lines are `TapeEntry`
   `{"type":"meta|message|tool_result|outcome","data":{…}}`
   (`tape.go:285-329`); `live.jsonl` lines are `liveEntry`
   (`runtime.go:943-952`), truncated per generation
   (`generateStreaming`, `runtime.go:971`); `log/control.jsonl` lines are
   `controlLogEntry` (`control.go:82-89`); `status/session.json`
   (`runtime.go:1620-1637`) and `status/inbox.json` are overwritten in place.
5. **The client endpoint is client-owned and client-served**: a real directory
   `runtime_root/client/<id>/public/{ctl,status,log}` where `ctl/post` is a
   **FIFO created and read by the client**, with `status/client.json`
   announcing the endpoint (old qctl `OpenClientEndpoint`,
   `ensureClientControlFIFO`, `ServeClientEndpoint`). This is the isomorphism
   made concrete — the kernel is a peer with its own surface, and any agent
   can `echo reply > client/<id>/public/ctl/post`.
6. **The named day-one bug**: the old qctl `Message()` opened the target ctl
   `O_WRONLY` with no timeout and blocked forever on a non-draining peer
   (`git show 88a78d65:internal/qctl/qctl.go`). The FUSE migration changes the
   mechanism (blocked FUSE call instead of FIFO open) but not the hazard.

## 3. Constraints (violations are bugs)

| # | Constraint |
|---|---|
| K1 | **Observe, never drive.** The kernel emits control events and reads files. No polling the agent into action, no lifecycle ops beyond `spawn`, no kill. |
| K2 | **Zero persistence.** No files owned by the kernel except the client endpoint surface (which is announcement, not state — recreated on start). `client_ref` associations, seq counters, unread anything: in-memory only. |
| K3 | **Agent-facing surfaces unchanged.** The envelope format (CONTRACT.md §6) and ctl files are the existing runtime contract; the kernel adapts to them. |
| K4 | **Every ctl write is bounded** (default 3s, `QCLI_CTL_TIMEOUT_MS`). On timeout: report `write_timeout`, mark peer unreachable in status, leave the write goroutine to finish or leak-bounded (§4.5). |
| K5 | **Never read `live.jsonl` for anything but relay** — it is display-only by runtime invariant (`scripts/check-live-surface-isolation.sh` guards the runtime side; the kernel must not launder it into `/context` or cells). |
| K6 | **Contract fidelity.** Wire schemas come from CONTRACT.md verbatim; kernel tests assert them field-for-field. |

## 4. Design

### 4.1 Code layout

```
internal/qcli/            (Go package, sibling of internal/runtime)
  resolve.go   resolve_test.go     selector → Agent          (revive)
  spawn.go     spawn_test.go       quine spawn + registration (revive)
  endpoint.go  endpoint_test.go    client endpoint FIFO + serve (revive)
  envelope.go  envelope_test.go    formatClientSignalPayload  (revive verbatim)
  ctlwrite.go  ctlwrite_test.go    bounded ctl write (new)
  cells.go     cells_test.go       TapeEntry line → []Cell    (revive + seq)
  tailctx.go   tailctx_test.go     current.jsonl tailer + compaction detect
  taillive.go  taillive_test.go    live.jsonl tailer + generation detect
  tailctl.go   tailctl_test.go     control.jsonl tailer + client_ref match
  status.go    status_test.go      session/inbox poll + derive working/live
  roster.go    roster_test.go      pid/* scan projection      (revive)
  contract.go  contract_test.go    peer contract.json proxy   (new)
  session.go   session_test.go     attached-peer lifecycle, event fan-in
  server.go    server_test.go      HTTP+SSE, /command, subscriber fan-out
operator/qcli/app/main.go          `qcli` binary: attach|spawn|serve + TUI launch
```

The package lives in `internal/` (buildable by the module, coverable by
`go test ./internal/...`) but remains outside the embedded Quine self-source.
qcli is a companion client, not part of the agent body; self-source tests
should keep `cmd/qcli`, `operator/qcli/app`, and `internal/qcli` excluded just
as qctl/world companion surfaces were excluded.

### 4.2 Session core (`session.go`)

One `Session` = one attached peer + one client endpoint + all tailers. All
events (cells, deltas, receipts, status) flow into a single in-memory ring
journal `[]Event` with kernel seq, guarded by a mutex; SSE subscribers replay
`hello → status → cells → backfill_complete` from current parsed state, then
follow the live feed via a per-subscriber buffered channel (cap 1024). A
subscriber that falls behind is **disconnected** (client rebuilds — cheap by
design); never block the tailers on a slow reader.

Attach tears down any previous session (v1: one peer). The client endpoint id
is `qcli-<unixms>-<pid>`; endpoint is recreated on every start (K2).

### 4.3 Tailers

**Context tailer (`tailctx.go`).** Open, parse all lines through `cells.go`
(seq-stamped), keep `(inode, size, offset)`; poll at 150ms (`QCLI_POLL_MS`)
with a partial-line buffer (a line is consumed only when `\n`-terminated —
the runtime appends whole lines, but reads can race mid-write). Detection:
size < offset **or** inode change ⇒ compaction: emit `context_reset
(reason:"compaction")`, reset cell state, re-parse from zero, re-emit
backfill. Incarnation bump in status ⇒ same with `reason:"incarnation"`.

**Live tailer (`taillive.go`).** Same line-gated tailing on `live.jsonl`.
Truncation (size < offset) ⇒ `generation++`, offset 0. Relay each `liveEntry`
as `stream_delta` verbatim + `generation`. Absent file = no generation in
flight (not an error). Never parsed into cells (K5).

**Control-log tailer (`tailctl.go`).** Line-gated tailing on
`log/control.jsonl`; map `controlLogEntry` → `receipt` (CONTRACT.md §3.6).
`client_ref` matching: `ctlwrite.go` records `(client_ref, envelopeText)` in a
ring of 64; a `received` whose `message.payload` equals a recorded envelope
gets that `client_ref` and creates `message_id → client_ref`; `delivered`
resolves through that map. Unmatched receipts relay with `client_ref:null`.

**Status poller (`status.go`).** Poll `status/session.json` +
`status/inbox.json` mtimes at 500ms; `live` via `syscall.Kill(pid, 0)` **plus
a `/proc/<pid>/stat` state check — signal-0 succeeds on zombies**, so an
exited-but-unreaped agent would otherwise read as live forever (found in
dogfood; `pidZombie` in `resolve.go`, regression `liveness_test.go`);
`working` := context/live tail advanced within the last 2s **or** an
unterminated generation is open **or** a `tool_call` awaits its
`tool_result`. Emit `status` only on change.

**Peer contract proxy (`contract.go`).** Read the attached peer's
`public/status/contract.json` and expose it through `GET /peer-contract`.
This is a read-only pass-through used by the TUI for action descriptions; the
kernel must not paraphrase or persist it. If the file is missing, return `404`
and let the client fall back to action labels only.

### 4.4 Cell parser (`cells.go`)

Revive `ParseTapeLine` and its helpers wholesale (§8) — 16 tests existed.
Contract deltas from the old parser: stamp kernel `seq`; emit `ts` from
`tape.Message.Timestamp`; `outcome` gains structured
`termination_mode/exit_code/turn_count/tokens_in/tokens_out`; unknown
`TapeEntry.type` ⇒ single cell `{kind:"<type>", raw:<line>}` (CONTRACT.md §8).
Keep: `reasoningText()` (summary text only, never encrypted blobs),
`extractIncomingHumanPosts` + `parseClientEnvelope` (authority-aware delivered
envelopes → `human_post` cells), `compactArgs`, `summarizeToolContent`. `reply` cells are
synthesized by `endpoint.go` when a payload arrives on the endpoint FIFO —
they never come from current.jsonl.

### 4.5 Bounded ctl write (`ctlwrite.go`) — the day-one fix

```go
func writeCtl(path string, payload []byte, timeout time.Duration) error
```

Run `os.OpenFile(O_WRONLY) + Write + Close` in a goroutine; select on
completion vs `time.After(timeout)`. On timeout return `ErrWriteTimeout` and
abandon the goroutine (it holds no locks; bound concurrent abandoned writers
to 8 per peer — beyond that, fail fast with `unreachable`). Classify errors:
`ENOTCONN`/`ENXIO` → `unreachable` (stale FUSE mount / dead peer); `ENOENT` →
`no_peer`. Empty-payload writes (poke/interrupt, CONTRACT.md §4.3) go through
the same path with zero bytes and **no envelope**.

### 4.6 HTTP server (`server.go`)

`net/http`, no framework. Routes per CONTRACT.md §2. `/command` dispatches to
session (attach/spawn/send); all responses synchronous JSON; `send` returns
after `writeCtl` acks or fails — never waits for agent behavior.
`/peer-contract` returns the attached peer's raw public contract JSON. Bind
loopback; refuse non-loopback without `QCLI_HTTP_UNSAFE=1`.

### 4.7 Entrypoint (`operator/qcli/app/main.go`)

```
qcli                       resolve default agent (or spawn) → serve on ephemeral port → exec TUI
qcli attach --pid N        same, explicit selector        (--session, --agent-root, --runtime-root)
qcli spawn [mission]       spawn + attach → serve → TUI
qcli serve [--http ADDR]   headless kernel only (web/dev/tests)
```

TUI launch: spawn `bun run <repo>/operator/qcli/tui/src/index.tsx` with
`QCLI_SERVER=http://127.0.0.1:<port>` in env; propagate exit code; on TUI exit,
shut down. If `bun` is absent, print the `qcli serve` + manual-TUI instructions
instead of failing silently. Config is env-only: `QCLI_HTTP_ADDR`,
`QCLI_CTL_TIMEOUT_MS`, `QCLI_POLL_MS`, `QCLI_QUINE_BIN`, `QCLI_RUNTIME_ROOT`.

## 5. Implementation plan (phased; `go build ./... && go test ./internal/...` after each)

- **P1 — bridge**: `resolve.go`, `spawn.go`, `endpoint.go`, `envelope.go`,
  `ctlwrite.go`. Revive from §8; port `Message()` onto `writeCtl`.
  Envelope byte-format must match §6 of CONTRACT.md exactly (golden test).
- **P2 — parse + tail**: `cells.go`, `tailctx.go`, `taillive.go`, `tailctl.go`,
  `status.go`, `roster.go`, `contract.go`. Golden fixtures: record a real agent's
  `current.jsonl`, `live.jsonl`, `control.jsonl` into
  `internal/qcli/testdata/` (one short real run; scrub nothing — it is already
  public-surface data).
- **P3 — session + server**: `session.go`, `server.go`. Contract tests drive
  `/events` with `httptest`, assert the exact event sequence and JSON fields
  against CONTRACT.md schemas.
- **P4 — receipts association**: ring matching in `tailctl.go` + `send`
  integration.
- **P5 — entrypoint + live smoke**: `app/main.go`, TUI launch, env config;
  wire T-E2E (§6.1) into the suite.

## 6. Test matrix (PTY-free; table is normative)

| ID | Asserts |
|---|---|
| T-RESOLVE | one-selector rule; pid symlink resolution; session under `agent/<id>`; not-found errors |
| T-SPAWN | registration poll; early-exit detection; 5s timeout (use a fake `quine` script fixture) |
| T-ENDPOINT | FIFO created; `client.json` announced; payload on FIFO → `reply` cell event |
| T-ENVELOPE | golden byte-equality with CONTRACT.md §6 (including trailing-newline normalization) |
| T-CTL-TIMEOUT | write to a FIFO with no reader (simulates wedged commit) returns `ErrWriteTimeout` within budget; abandoned-writer cap enforced |
| T-CELLS | revive the 16 old parser tests; add: unknown entry type → raw cell; outcome structured fields; ts propagation |
| T-TAIL-CTX | snapshot-then-tail; mid-line write race (partial buffer); compaction (truncate file) → `context_reset` + full re-backfill |
| T-TAIL-LIVE | generation increments on truncation; `completed` relay; absent file tolerated |
| T-RECEIPT | received/woke/delivered mapping; client_ref payload match; foreign writes relay with null client_ref |
| T-STATUS | working derivation cases; live=false on dead pid **and on zombie pid** (signal-0 passes on defunct processes); pending from inbox.json |
| T-ROSTER | scan of a fixture runtime-root: live/dead/foreign-symlink entries |
| T-PEER-CONTRACT | `/peer-contract` returns raw attached `public/status/contract.json`; 404 when absent; no paraphrase/persistence |
| T-CONTRACT-SSE | full `/events` sequence field-for-field vs CONTRACT.md; unknown-field tolerance is the *client's* duty, so the server must not emit undocumented fields |
| T-COMMAND | attach/spawn/send happy + error paths; empty-payload rules; 409 before attach |
| T-ZERO-STATE | kill server mid-tail, restart, reconnect: identical backfill (same fixture ⇒ same cells) |
| T-E2E | automated live smoke on the gpt-5.5 profile — full loop: spawn → inject → receipts → stream deltas → settled cells → reply (§6.1) |

### 6.1 T-E2E — automated live smoke (`e2e_smoke_test.go`)

Runs in the normal suite against a **real agent** on the repo-tracked
gpt-5.5 profile (`profiles/gpt-5.5-codex-oauth.env` — `openai-responses`,
the fully-streaming provider path, so this test also closes STREAMING.md's
pending live-deltas acceptance item).

1. **Profile loading**: resolve env via
   `bash -c 'set -a; source profiles/gpt-5.5-codex-oauth.env; set +a; env'`
   and diff against the parent env — profiles may contain shell logic; never
   reimplement their resolution in Go.
2. **Gating**: `QCLI_E2E=0` opts out explicitly. Otherwise the test runs, and
   **skips loudly only when credentials are non-interactively unavailable**
   (`~/.codex/auth.json` missing or token unrefreshable without user
   interaction). The `codex-oauth` sentinel can start a browser OAuth flow
   when no fresh token exists (`profiles/README.md`) — a test must **never**
   trigger that; probe token freshness first and `t.Skip` with the reason.
3. **Drive**: spawn a real quine with the profile env (mission: *"When you
   receive a message containing a reply_ctl path, immediately write the word
   pong to that path, then exit."*), attach, `send` one `inject`.
4. **Hard asserts**, in order, within a 180s budget: `queued` ack →
   `receipt/received` with matched `client_ref` → `stream_delta` events
   (≥1 `text_delta`/`reasoning_delta`, then a terminal `completed`) →
   settled `cell`s for that generation → `receipt/delivered` → a `reply`
   cell (`pong`) on the client endpoint.
5. **Teardown**: the test terminates the quine child it spawned. This is a
   harness cleaning up its own child process, not a kernel peer-kill — K1
   binds the kernel, not the test harness.

Cost: one short gpt-5.5 generation per suite run. Watchpoint: the final
`reply` assert leans on agent compliance (a reply is agent-chosen physics,
not a protocol guarantee). If it proves flaky, demote **only that assert** to
a logged warning — the receipt/delta/cell asserts are protocol-backed and
stay hard.

## 7. Acceptance criteria

- [ ] All §6 tests green; `go build ./...`, `go vet ./...`,
      `go test ./internal/...` green.
- [ ] Self-source stays green (`TestEmbeddedSelfSourceBuildsAndTests`) — qcli
      stays excluded from the embedded agent self-source unless the self-source
      policy is deliberately revised.
- [ ] `qcli serve` + `curl` reproduces the CONTRACT.md §3 event sequence
      against a fixture agent-root with no live agent (dead-peer status,
      backfill still served).
- [ ] T-E2E green in the suite on a credentialed machine (auto-skips only
      when codex OAuth credentials are non-interactively unavailable; the
      skip is loud). Its `stream_delta` assert closes STREAMING.md's pending
      live-smoke acceptance item — tick it in the same commit.
- [ ] Zero files created outside `runtime_root/client/<id>/` and the process's
      own memory (K2 audit: `strace`-check or fs-diff in a test).

## 8. Revival index (`git show 88a78d65:<path>`)

| Old file | Take | Leave |
|---|---|---|
| `internal/qctl/qctl.go` | `Resolve`, `loadAgentFromRoot`, `candidateRuntimeRoots`, `Spawn`, `mergedSpawnEnv`, `OpenClientEndpoint`, `ensureClientControlFIFO`, `ServeClientEndpoint`, `formatClientSignalPayload`, `verifyLivePID`, `readFIFOOnce` | the RPC layer, `Message()`'s unbounded open, the session/TUI glue |
| `internal/qctl/cell.go` + `cell_test.go` | the whole parser + tests | old RPC Cell envelope fields superseded by CONTRACT.md §3.3 |
| `internal/qctl/roster.go` | scan logic | roster RPC events (v1 polls) |
| `internal/qctl/httpserver.go` | SSE writer pattern, subscriber fan-out shape | `/roster` push, group scope |
| `internal/qctl/group.go` | *nothing for M1* — this is the reserved N-peer fan-in; re-read when group returns | — |

## 9. Reference index (current runtime, verified)

| Fact | Where |
|---|---|
| FUSE-only backend constant + decision pointer | `internal/runtime/public_surface.go:1-16` |
| ctl node write/commit path | `internal/runtime/public_surface_fuse_linux.go:271-412` |
| control actions + semantics + inbox schema + log events | `internal/runtime/contract_surface.go:79-116` |
| `controlLogEntry` fields; received/woke/delivered emission | `internal/runtime/control.go:82-89,365,384,481` |
| `session.json` fields | `internal/runtime/runtime.go:1611-1638` |
| `liveEntry` schema; truncate-per-generation; non-fatal writes | `internal/runtime/runtime.go:943-1035` |
| `TapeEntry` framing; `Message` payload shape | `internal/tape/tape.go:285-329,71-82` |
| pid symlink → public root; session.json join | `internal/runtime/peer_discovery.go:203-262` |
| live-surface isolation validator (runtime side of K5) | `scripts/check-live-surface-isolation.sh` |
