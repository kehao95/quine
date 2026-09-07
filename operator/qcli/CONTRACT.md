# qcli Contract — `qcli/1` (M0, frozen)

> The **normative boundary** between qcli-server (Go kernel) and every qcli
> client (terminal TUI now, web later). Both sides build against this document;
> neither side may depend on anything not written here. Design rationale lives
> in [`DESIGN.md`](./DESIGN.md); UX decisions in [`UX.md`](./UX.md); this doc is
> schema only.
>
> **Freeze rule:** additive changes (new event type, new optional field, new
> cell kind) bump nothing — clients must tolerate them (§8). Any change to an
> existing field's meaning or type bumps the contract string and is a design
> event, not an implementation detail.

## 1. Transport

- HTTP + SSE over `127.0.0.1` (default bind; configurable via
  `QCLI_HTTP_ADDR`). No auth in v1 — the trust boundary is the local user, same
  as the filesystem surfaces themselves. Binding non-loopback is refused
  without an explicit `QCLI_HTTP_UNSAFE=1`.
- SSE framing: each frame is `data: <one JSON object>\n\n`. No SSE `event:` /
  `id:` fields — the discriminator is the JSON `type` field. Reconnect is
  full-rebuild (§5), so `Last-Event-ID` is deliberately unused.
- One server observes **one attached peer** in v1. Multi-peer fan-in is
  reserved (§9); `author` fields already exist where noted.

## 2. Endpoints

| Endpoint | Method | Purpose |
|---|---|---|
| `/events` | GET | SSE: `hello` → `status` → backfill `cell`s → `backfill_complete` → live tail |
| `/context` | GET | parsed settled cells as a JSON array (non-SSE consumers, tests) |
| `/status` | GET | latest status object (JSON) |
| `/roster` | GET | population snapshot from `runtime-root/pid/*` scan |
| `/peer-contract` | GET | raw attached peer `public/status/contract.json` for UI labels/help |
| `/command` | POST | `attach` / `spawn` / `send` |
| `/healthz` | GET | `200 {"ok":true,"contract":"qcli/1","attached":bool}` |
| `/history` | GET | **reserved** (tape backfill, M3) — returns `501` in v1 |

`GET /events` before any successful `attach`/`spawn` returns `409` with an
`error` body (`code:"no_peer"`).

## 3. Event schemas (`/events`)

### 3.1 `hello` — first frame on every connection

```json
{"type":"hello","contract":"qcli/1","endpoint_id":"qcli-1719300000-48100",
 "peer":{"session":"sess_…","run_id":"run_…","pid":48213,"agent_root":"/…","runtime_root":"/…"}}
```

### 3.2 `status` — snapshot after `hello`, then on change

```json
{"type":"status","live":true,"working":false,"pid":48213,"ppid":48100,
 "session":"sess_…","run_id":"run_…","incarnation":2,"parent_session":"",
 "model":"gpt-5.5","depth":2,"pending":0,"tokens_in":null,"tokens_out":null}
```

- `live` — PID liveness (signal-0 check). `working` — derived: `current.jsonl`
  advance or an in-flight generation (`live.jsonl` activity) or unresolved
  tool_call pairing. `pending` — from `status/inbox.json` `pending_count`.
- `tokens_in/out` are `null` until runtime RFC R1 lands (UX.md §7) — the
  kernel **never** fabricates them.

### 3.3 `cell` — one settled cognition/track cell

```json
{"type":"cell","seq":42,"kind":"message","role":"assistant",
 "text":"I'll check the novelty gate","ts":1719300000123,
 "tool_id":null,"tool_name":null,"status":null,"is_error":false,
 "action":null,"author":null,"raw":null}
```

`seq` is a kernel-local monotonic counter stamped at the tailer boundary
(display ordering only, not a causal clock). `author` is populated for
authority-bearing cells such as `human_post`; otherwise it is `null`. Fields by
kind — absent/irrelevant fields are `null`:

| kind | populated fields | source |
|---|---|---|
| `meta` | `text` (model), `depth`, `session` | current.jsonl `meta` entry |
| `message` | `text`, `role`, `ts` | `message` entry, assistant/user text |
| `reasoning` | `text` (summary text, **never** encrypted id blobs) | `reasoning_content` / reasoning items of a `message` entry |
| `tool_call` | `tool_name`, `tool_id`, `text` (args as compact JSON) | `tool_calls[]` of an assistant `message` |
| `tool_result` | `tool_id`, `tool_name`, `status`, `text` (output preview), `is_error` | `tool_result` entry |
| `human_post` | `text` (payload body, envelope stripped), `action`, `author` | client-envelope user message observed in current.jsonl (what was actually **delivered**) |
| `reply` | `text`, `ts` | payload arriving on the client endpoint FIFO |
| `outcome` | `text` (packed summary), plus `termination_mode`, `exit_code`, `turn_count`, `tokens_in`, `tokens_out` | `outcome` entry |

One current.jsonl line may fan out to several cells (e.g. one assistant entry
→ `reasoning` + `message` + N × `tool_call`), in stable order.

### 3.4 `backfill_complete`

```json
{"type":"backfill_complete","count":57}
```

Everything before it is replayed state; everything after is live tail.

### 3.5 `stream_delta` — in-flight generation (relay of `live.jsonl`)

```json
{"type":"stream_delta","generation":3,"seq":7,"kind":"text_delta",
 "text":"checking the score no","tool_id":null,"tool_name":null,"ts":1719300000456}
```

- `kind` ∈ `text_delta | reasoning_delta | tool_call_delta | completed` —
  verbatim relay of the runtime's `liveEntry` (`internal/runtime/runtime.go:945`).
- `generation` is kernel-local, incremented each time the tailer observes the
  per-generation truncation. `completed` ends a generation; the authoritative
  settled cell(s) follow as normal `cell` events.
- Display-only by inheritance: this stream mirrors a surface that is never
  provider input (STREAMING.md C2). A generation that ends **without**
  `completed` (interrupt/error) simply stops; clients render the
  interrupted register (UX.md G7).

### 3.6 `receipt` — delivery lifecycle (relay of `log/control.jsonl`)

```json
{"type":"receipt","stage":"received","action":"inject","delivery":null,
 "message_id":"message-3","client_ref":"c-17","pending":null,"ts":1719300000789}
```

- `stage` ∈ `received | woke | delivered` — mapped 1:1 from control-log
  `kind` (`internal/runtime/control.go:82`). `received`/`delivered` carry
  `message_id` (`control.go:365,481`); `woke` carries `delivery` + `pending`
  and no `message_id`.
- `client_ref` is set when the kernel matches the entry's payload to one of
  its own recent sends (§4.3) — volatile, in-memory, best-effort. Receipts
  for writes made by *other* peers relay with `client_ref:null` (they are
  legible too; the surface is shared).
- A fourth stage `queued` exists **only** as the synchronous `/command` ack
  (§4.3), never as an SSE event.

### 3.7 `context_reset`

```json
{"type":"context_reset","reason":"compaction","ts":1719300001000}
```

Emitted when the context tailer observes truncation/inode change of
`current.jsonl` (`reason:"compaction"`) or an `incarnation` bump in status
(`reason:"incarnation"`). Followed immediately by a fresh backfill
(`cell`* → `backfill_complete`). Clients decide how to render what they
previously held (UX.md G1: demote, divider).

### 3.8 `error`

```json
{"type":"error","scope":"tail","code":"peer_dead","message":"pid 48213 gone","recoverable":false}
```

`scope` ∈ `send | attach | spawn | tail | server`. Errors never masquerade as
cells.

## 4. Commands (`POST /command`)

One JSON object per request; synchronous JSON response; never blocks on agent
behavior.

### 4.1 `attach`

```json
{"op":"attach","selector":{"agent_root":"","pid":0,"session":"","runtime_root":""}}
```

Exactly one of `agent_root | pid | session` (plus optional `runtime_root`
scope) — the old qctl selector semantics. Response:
`{"type":"attached","peer":{…as hello.peer…}}` or `error`
(`code:"target_required"|"not_found"`).

### 4.2 `spawn`

```json
{"op":"spawn","mission":"…","runtime_root":""}
```

Spawns the `quine` binary (resolution: `QCLI_QUINE_BIN`, else sibling of the
qcli binary, else `$PATH`), waits ≤5s for `pid/<PID>` registration, then
attaches. Response as `attach`, or `error` (`code:"spawn_failed"|"register_timeout"`).

### 4.3 `send`

```json
{"op":"send","action":"inject","text":"look at the empty-package bug"}
```

- `action` ∈ `post | poke | inject | interrupt`.
- Non-empty `text`: the kernel wraps it in the client envelope (§6) and writes
  it to the attached peer's `public/ctl/<action>` with a bounded write
  (§7 write timeout). Empty `text` is valid **only** for `poke` and
  `interrupt` (the runtime's empty-write semantics) and writes zero bytes,
  no envelope; empty `post`/`inject` are rejected (`code:"empty_payload"`).
- Response (immediate): `{"type":"queued","action":"inject","client_ref":"c-17"}`
  — this is the ladder's `queued` stage. Failure:
  `{"type":"error","scope":"send","code":"no_peer|unreachable|write_timeout|empty_payload","message":"…"}`.
- `reply_required` is `false`, always, structurally: there is no field to set.

## 5. Reconnect / restart semantics

The client owns **zero** durable state; the server owns **zero** durable
state. On SSE drop the client reconnects and rebuilds from the new
`hello → … → backfill_complete` sequence, discarding everything it held. On
server restart the same holds (endpoint id changes; `client_ref` associations
and unread markers are honestly lost). No resume tokens exist by design.

## 6. Client envelope (kernel → agent ctl, informative here, owned by KERNEL.md)

The payload written to `ctl/<action>` for non-empty sends — qctl-derived text
preamble with an explicit client authority field:

```
[qcli-client]
authority: human
ctl_action: inject
reply_ctl: <runtime_root>/client/<endpoint_id>/public/ctl/post
reply_required: false

message:
<payload>
```

The `human_post` cell parser strips this envelope when it reappears in
current.jsonl; `authority` becomes the cell `author` (legacy envelopes without
it default to `human`), and `reply_ctl` is how agent-chosen replies reach the
client endpoint. The envelope is **not** part of this contract's wire schema —
it is the agent-facing surface, and changing it is a runtime-facing design
event.

## 7. `GET /roster`

```json
{"type":"roster","peers":[
  {"pid":48213,"session":"sess_…","run_id":"run_…","parent_session":"",
   "depth":2,"model":"gpt-5.5","live":true,"pending":0,
   "agent_root":"/…","attached":true}
]}
```

Projection of `runtime-root/pid/*` symlinks (→ each agent's `public/` root,
`internal/runtime/peer_discovery.go:203-262`) joined with each
`status/session.json` + liveness verify + `status/inbox.json`. Computed on
request; never cached beyond the request.

## 7.1 `GET /peer-contract`

Returns the attached peer's raw public contract document:

```json
{"type":"peer_contract","contract":{ "...": "raw public/status/contract.json" }}
```

The kernel reads `agent_root/public/status/contract.json` on request and wraps
the decoded object without paraphrase. Before attach, return `409 no_peer`;
when the peer lacks the file, return `404 not_found`. This endpoint exists so
clients can render control-action descriptions from the peer's own contract
surface without reading the filesystem directly.

## 8. Forward compatibility (client obligations)

- Unknown event `type`: ignore silently.
- Unknown cell `kind`: render `raw` (the original JSON line) as a monospace
  block — Quine growing a new cell kind must not break clients.
- Unknown/extra fields anywhere: ignore.
- Absent optional fields: treat as `null`.

## 9. Reserved (do not implement, do not repurpose)

- `author`, `ingest_seq` on `cell` — multi-peer group scope (`?scope=group`).
- `GET /history?peer=&before_seq=&limit=` — tape-backed backfill (M3), cells
  carry `"register":"history"`.
- `roster` push events on `/events` — v1 clients poll `/roster` on demand.

## 10. Surface truth (where every field comes from)

| Contract element | Runtime source | Evidence |
|---|---|---|
| ctl actions + semantics | FUSE ctl nodes; graded delivery | `internal/runtime/control.go:29-33`, contract manifest `contract_surface.go:79-105` |
| status identity fields | `status/session.json` | `internal/runtime/runtime.go:1620-1637` |
| `pending` | `status/inbox.json` `pending_count` | `contract_surface.go:106-111` |
| cell source lines | `context/state/current.jsonl` (TapeEntry framing `{"type","data"}`) | `internal/tape/tape.go:285-329`, `runtime.go:1376` |
| message payload shape | `tape.Message` | `internal/tape/tape.go:71-82` |
| stream_delta lines | `context/state/live.jsonl` (`liveEntry`) | `internal/runtime/runtime.go:943-952` |
| receipt lines | `log/control.jsonl` (`controlLogEntry`) | `internal/runtime/control.go:82-89,365,384,481` |
| roster scan | `runtime-root/pid/<PID>` → `public/` symlink | `internal/runtime/peer_discovery.go:203-262` |
| peer contract self-description | `public/status/contract.json` | `internal/runtime/contract_surface.go:56-125` |
