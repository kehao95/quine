# qcli — Client Design

> The current client design for operating Quine. It replaces the removed qctl
> design surface and keeps the surviving constraints here: protocol-equal
> projection, control events rather than request/response, zero persistence,
> optional 0..N agent replies, and a single-view first cut over Quine's live
> process surface.
>
> UX maturity layer: [`UX.md`](./UX.md) — opencode-parity inventory, delivery
> receipts, gap decisions beyond the first cut, and milestone gates.
> Executable specs: [`CONTRACT.md`](./CONTRACT.md) (frozen `qcli/1` wire
> schema), [`KERNEL.md`](./KERNEL.md) (Go kernel, M1), [`TUI.md`](./TUI.md)
> (terminal client, M2), [`WEB.md`](./WEB.md) (web client, M5). Doc map:
> [`README.md`](./README.md).

## Design Principles

Three rules govern every decision below. A design that breaks any is wrong.

### 1. Client adapts to Quine's shape, not the reverse.

The contract is a **projection of Quine's actual filesystem surface**, not a
spec the runtime must satisfy. When Quine evolves (new tool, new cell kind, new
ctl action), the client follows — Quine never grows an interface *for the
client*. Concretely: every contract field traces to a real ctl file,
`current.jsonl` entry, or status projection.

This runs in both directions — host-facing protocols lower *through* this
client, they do not reach past it (moved from the Signal-First Interaction
Contract, 2026-07-26; that contract now lives at
[`development/contracts/signal-first-interaction.md`](../../development/contracts/signal-first-interaction.md)):

- if ACP or a similar editor/host protocol becomes useful, expose it **from this
  client** as an adapter over Quine's existing process surface; do not implement
  it inside the `quine` binary
- editor- and host-facing protocol assumptions must not leak backward into
  Quine's native signal/control ontology
- keep the command surface non-transport-shaped: `attach` / `spawn` (see
  § Entry points) plus filesystem `ctl` for non-interactive operation, rather
  than transport verbs such as `send`

### 2. Protocol-equal projection bridge, not a runtime container.

qcli-server is a **communication peer** on Quine's ctl surface (it speaks the
same `public/ctl/{post,poke,inject,interrupt}` protocol any peer would), but
**ontologically it is Human's foreground projection** — not a second Quine. It
has no agent loop, no tools, no mission, no selfsource. The "equality" is
protocol-level only. This is the isomorphism invariant from
the [Signal-First Interaction Contract](../../development/contracts/signal-first-interaction.md),
held precisely: protocol peer, not master, not a peer-agent.

The decisive contrast with opencode: opencode's client-server is a master that
drives the agent via RPC drain. qcli-server **cannot drive Quine** — it only
emits control events (`reply_required: false` always) and observes the surface.
This is non-negotiable because Quine's research (dormancy-ecology, selection
substrate, agent-owns-lifecycle) is built on the agent being autonomous.

### 3. Zero persistence. Pure projection.

qcli-server owns **no durable state**. It tails Quine's surface, parses it,
serves it over the contract. Restart = rebuild from surface. This is what
makes it safe as a daemon without becoming a single point of failure: nothing
is lost when it dies, because Quine's `current.jsonl`, tape, and ctl files are
the only truth. (Contrast opencode's daemon, which owns durable session state
and must recover on crash.)

**Red line:** any proposal to add "client-owned session history", "unread
state", or "preference store" to qcli-server violates this principle. All truth
lives in Quine's surface.

## Architecture

```
                 control events (inbound)                 observation (outbound)
                 ─────────────────────                    ─────────────────────
Human ──post/poke/inject/interrupt──▶ qcli-server ──ctl──▶ Quine agent
                                     (Go, stateless)       (owns its lifecycle)
Human ◀──cognition cells / status──── qcli-server ◀──────── current.jsonl
       (rendered by TUI)             ◀──────────────────── client/<id>/ctl/post
                                                          public/status/session.json
```

```
┌─────────────────────────────────────────────────────────────┐
│  Quine runtime  (UNCHANGED, owns its lifecycle)              │
│  surface: public/ctl/{4} · current.jsonl · client/ctl · status│
└────────────┬───────────────────────────────────┬───────────┘
             │ filesystem ctl surface            │
             ▼                                   ▼
┌─────────────────────────────────────────────────────────────┐
│  qcli-server  (Go kernel — stateless projection bridge)      │
│  ① ctl bridge      — resolve/spawn agent, write ctl events   │
│  ② surface tailer  — current.jsonl + client/ctl + status     │
│  ③ cell parser     — tape line → typed Cell (schema-driven)  │
│  ④ fan-in seam     — reserved for later N-peer merge (§6)     │
│  ⑤ contract server — HTTP+SSE                                │
└────────────────────────────┬────────────────────────────────┘
                             │ contract (Cell schema + events)
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  qcli  (TS + OpenTUI + SolidJS — pure projection renderer)   │
│  · individual view (the only view)                           │
│  · control-event composer (4 actions, no addressing)         │
│  · zero own state (rebuilds from SSE on reconnect)           │
└─────────────────────────────────────────────────────────────┘
```

### Roles, precisely

| Layer | Is | Is not |
|---|---|---|
| **qcli-server** | ctl surface peer (protocol-equal), surface projector, contract server | Quine's master, a second agent, a session owner |
| **qcli (TUI)** | projection renderer, control-event editor | ctl writer (goes through qcli-server), state owner |
| **Quine agent** | the autonomous runtime, owns lifecycle | driven by qcli-server (it receives events, decides itself) |

### Stack (validated — do not re-litigate)

| Layer | Choice | Evidence |
|---|---|---|
| Terminal client | **TypeScript + `@opentui/core` + `@opentui/solid` + `@opentui/keymap` + `solid-js`, run via `bun`** | OpenTUI is OpenCode's own production TUI core ("powers OpenCode in production today"); MIT; actively maintained (`@opentui/core@0.4.2`, 281 versions, 440k weekly downloads). A v1 client on this exact stack was built and rendering-verified in this repo (deleted `operator/qctl/tui-client/` — removed only because of the qctl→qcli architecture rewrite, not stack failure). OpenTUI is a native Zig core + TS bindings; SolidJS is the reactive layer, so the same Solid components can re-render to a future web client. |
| Server kernel | **Go** | same language as the Quine runtime; stateless projection bridge. |

**Accepted tradeoffs (informed, not blocking):** two toolchains (Go kernel + bun client) is the deliberate price of OpenCode-level rendering polish + web-layer sharing; OpenTUI ships a per-platform native core (bundle accordingly); OpenTUI is pre-1.0 but production-hardened via OpenCode; the deleted v1 targeted the old group-chat model, so **revive its patterns** (markdown, tool_call/result pairing, theme, SSE→Solid store, char-frame snapshot tests) but **rebuild for the single-view interaction model**.

## Quine Surface (the shape to adapt to)

The contract is a direct projection of these surfaces. Nothing in the contract
exists without a traceable source here.

### Control surface (inbound — Human → agent)

```
agent/<session>/public/ctl/{post,poke,interrupt,inject}
```

Four graded delivery actions, always `reply_required: false`:

| Action | Semantics |
|---|---|
| `post` | inbox only, lowest priority |
| `poke` | queue + resume, no auto-delivery |
| `inject` | **default** — deliver as context at next safe point |
| `interrupt` | urgent, preempts active tool work |

No request-id, no correlation, no callback. The agent may reply 0..N times,
arbitrarily delayed, or never. qcli-server's `send` writes to exactly one of
these four files. It never blocks waiting for a reply.

### Observation surface (outbound — agent → Human)

| Source | What it is | Mutability |
|---|---|---|
| `context/state/current.jsonl` | agent's live cognition (message/reasoning/tool_call/tool_result/meta/outcome) | **mutable working memory** — may be compacted/truncated by the agent |
| `client/<id>/public/ctl/post` | agent-chosen replies written back to qcli-server's client endpoint | append-only |
| `public/status/session.json` | status projection (pid/session/model/depth/pending) | overwritten in place |

**Critical:** `current.jsonl` is working memory, not archive. It can be
compacted. The cell stream is a live tail, not a durable transcript. (Tape is
the durable record, but tape is out of scope for the first cut — offline
analysis tools consume it, not the live UI.)

## Individual View (the only view)

The first cut ships **one view**: the individual peer's cognition stream.
Group view, broadcast addressing, and parallel multi-peer observation are
**out of scope** for this cut (kernel reserves space for them, see §6).

### Why one view

Most operation is 1 Human ↔ 1 peer. Group/broadcast are research-era features
for the peer-ecology experiments; they add UI complexity without serving the
primary use case. A single excellent view beats two mediocre ones. Peer
switching is a command (`/switch <peer>`), not a view.

### What "approach opencode" means — and doesn't

opencode's session view is the quality target for **rendering and interaction
polish**. But its **interaction semantics are incompatible with Quine** and
must be rebuilt. Two layers:

| Layer | Borrow from opencode? |
|---|---|
| Rendering quality (transcript scroll, reasoning fold, tool call/result pairing, markdown/code, composer IME, visual density) | **Yes — directly** |
| Interaction semantics (prompt-driven drain, permission dialog, turn-based transcript, subagent inline) | **No — violates Quine's shape** |
| Token typewriter | **Partial** — transient stream on `live.jsonl` only; durable cells in `current.jsonl`/tape stay complete (see deviation #2) |

### Five interaction deviations from opencode

These five are where "approach opencode" stops. Each is a direct consequence of
Quine's shape, not a preference.

**1. No permission dialog.** opencode authorizes each tool call because the
client is master. Quine's agent owns its lifecycle; the client is a protocol
peer with no authority to permit or deny tool calls. Tool calls render in the
transcript directly, no confirmation. Human intervention is `interrupt` (a
peer-level control event that preempts active tool work), not per-call
authorization.

**2. Complete cells in the durable layer; live preview on a distinct transient
surface.** Two surfaces, sharply separated by role:

- **`tape` and `context/state/current.jsonl`** are **complete-cell surfaces**.
  `current.jsonl` is the provider-input and recovery surface — it is read back
  as provider messages (`runtime.readContextMessages`), and a lost entry breaks
  tool_use/tool_result pairing — so it carries one whole cell per append, never
  partial deltas. Tape is the durable record. Both crystallize a cell only when
  generation completes.
- **`context/state/live.jsonl`** is a **transient generation stream**: the
  in-flight token/reasoning deltas as they arrive from the provider. It is
  display-only, append-only, truncated at the start of each generation, and
  **never read as provider input, never recovered, never replayed.** Any peer
  may `tail -f` it for a live preview; the client renders it in-place as a
  streaming cell (see [§ Streaming register](#streaming-register-in-flight-cells)).
  When generation completes, one authoritative cell crystallizes into
  `current.jsonl` + tape and the transient stream is rotated.

The earlier "no token streaming engine and never will be" was grounded in the
buffering artifact (`provider.go` `io.ReadAll` collapsing the SSE stream
before decode), not in a Quine invariant. The real invariants are the two
above: durable surfaces stay complete-cell; transient progress lives on a
distinct never-provider-input surface. Making in-flight cognition legible
through a stream is directional progress (`north-star`: "hidden cognition
becomes legible through streams"), provided the role boundary is named in the
contract and enforced by a validator. Full surface contract and the
reproducibility boundary live in [§ Live generation stream](#live-generation-stream-livejsonl).

**3. Transcript is non-turn-based — two interleaved tracks.** opencode's
transcript is `prompt → drain → reply` turns. Quine's transcript interleaves
two semantically different tracks:

```
◀ human_post·inject  look at the empty-package bug     ← control event (Human track)
┊ message   I'll check the novelty gate                ← cognition (agent track)
┊ reasoning ▸ (collapsed)                              ← agent thinking
┊ tool_call ▸ bash: awk -f scorer.awk ...             ← agent action
┊ tool_result ▸ status=ok · 3/5 flagged empty         ← agent result
┊ message → room   novelty gate missed empty pkgs     ← agent spoke (would enter group layer)
 · peer-7 working · pending=0 ···                      ← liveness lane (not a cell)
◀ reply   pasted diff to peer-3                        ← async reply, not paired to any human_post
```

The Human's `human_post` and the agent's cognition are **different tracks**.
Never assume a `human_post` is followed by a reply. Never pair them.

**4. depth/fork lineage — lightweight.** opencode's `task` subagent is a
one-shot, expanded inline. Quine's `fork` creates a **lineage ecology** — child
peers are independent long-lived peers, not disposable tasks. First cut:
**lightweight** — `depth` shown in header, fork points marked in-stream. Child
peers are navigated via `/switch`, not expanded inline. (Inline fork-tree
browsing would conflict with the clean transcript goal.)

**5. 0..N async replies — transcript tolerates absence.** A control event may
get 0, 1, or many replies, at arbitrary delay. The transcript must not show a
"waiting for reply" spinner after a `human_post`. Replies are independent cells,
not correlated to any input.

### Layout

```
┌ peer-7 · pid=48213 · session=…9d4d · model=gpt-5.5 · depth=2 · ●live ──────┐
│ COGNITION (current.jsonl settled cells + live.jsonl in-flight)              │
│                                                                              │
│  ◀ human_post·inject  look at the empty-package bug                         │
│  ┊ message   I'll check the novelty gate                                    │
│  ┊ reasoning ▸ (3 lines, → to expand)                                      │
│  ┊ tool_call ▸ bash: awk -f scorer.awk gen-*/decls                         │
│  ┊ tool_result ▸ status=ok · 3/5 flagged empty                             │
│  ┊▒ message   checking the score no_            ← streaming (dim, cursor)  │
│  · peer-7 working · pending=0                    ← lane: status only       │
│  ◀ reply   pasted diff to peer-3           (async, unpaired)                │
├ status ────────────────────────────────────────────────────────────────────┤
│  peer-7 working · 2 tool calls this turn · pending=0                        │
├ input   action=inject  →  peer-7 ───────────────────────────────────────────┤
│  > ...                                                                       │
└ [esc]interrupt   /switch <peer>   /action post|poke|inject|interrupt ───────┘
```

### Cell rendering pipeline

Every cell flows through the same renderer; `kind` selects the component:

| Cell kind | Render | Default state |
|---|---|---|
| `meta` | header line (model, depth, session) | inline in header, not in stream |
| `human_post` | Human track (left-aligned, human color, action tag) | visible |
| `message` | text bubble (agent color) | visible |
| `reasoning` | collapsible block (first-line summary + expand) | **collapsed** |
| `tool_call` | paired with result by `tool_id`; collapsible | collapsed header (`tool_name`) |
| `tool_result` | paired under its `tool_call`; status + output preview | collapsed, expandable |
| `reply` | agent reply cell (distinct from message — it crossed the wire) | visible |
| `outcome` | termination line (exit_code, turn_count, mode) | visible, end-of-stream |

**tool_call ↔ tool_result pairing** by `tool_id`. Visually: call is parent,
result is indented child, connected. If a result never arrives (agent
interrupted), the call renders orphaned with no spinner.

**reasoning** never shows the encrypted reasoning id blob — only the summary
text (already handled by the cell parser: `reasoningText()` prefers
`reasoning_content`, falls back to reasoning summary items).

### Streaming register (in-flight cells)

During generation, `live.jsonl` deltas render **in-place in the cognition pane**
as a cell in a `streaming` state — not as a status-bar caption. This is the
"alive" seam, and it is Quine-justified rather than opencode-imitation:
Quine's agent is autonomous and the only intervention is `interrupt`, so seeing
in-flight content (especially a `tool_call` forming) carries more decision
value than in a per-call-gated client like opencode.

**Why in-place, not a status line.** One generation emits a structured sequence
`reasoning_delta → text_delta → tool_call_delta` — i.e. the sub-elements of a
single assistant cell forming in turn. The honest render is therefore a cell
that happens to be incomplete, reusing the cell pipeline. A status-bar caption
(one truncated line) cannot express multi-sub-element formation and loses the
interrupt-decision value.

**Streaming register — visual demotion of "forming" vs "settled".**

| Register | Visual | Content |
|---|---|---|
| `streaming` | dim / reduced opacity, left `▒` rail, trailing cursor block, optional pulse | in-flight deltas from `live.jsonl` |
| `settled` | full opacity, normal styling | crystallized cell from `current.jsonl` |
| `interrupted`/`errored` | strikethrough / red tint, marked "not persisted" | mid-stream failure/interrupt (see below) |

On the provider's `completed` event the streaming cell `snap`s to `settled`:
the authoritative crystallized cell replaces it at full styling.

**Liveness lane is status-only.** With in-flight content living in the stream,
the lane stops carrying partial text and becomes a pure heartbeat
(`working · pending=N · 2 tool calls this turn`). Content lives in the stream
(settled + in-flight); status lives in the lane — the opencode-proven split.

**Reasoning during streaming.** Settled `reasoning` is collapsed by default;
during streaming the collapsed header shows a live one-line summary
(`▸ thinking: the novelty gate seems…`), expandable on demand — the engagement
of "watching it think" without flooding the pane.

**Crystallization jitter.** `text_delta` concatenation ≡ the `completed` text
(protocol-guaranteed), so `message` snap is ~zero-jitter. `reasoning` /
`tool_call` may differ slightly; the streaming cell renders stable fields
(text) at full confidence and treats reasoning/tool_call as conservative
previews. The crystallized cell is authoritative and overwrites on `completed`
— the "sacrifice a bit of repeatability" the project already endorsed.

**Interrupted / errored mid-stream.** If generation fails or is interrupted
before `completed`, no crystallized cell is produced — the streaming cell never
reached `current.jsonl`. It must not masquerade as a persisted cell: it stays
in place but transitions to the `interrupted`/`errored` register (clearly
marked "not persisted"), and the lane shows the error/stop status. The user
sees "it produced this much, then died" — honest about non-durable state.

### Live generation stream (`live.jsonl`)

In-flight generation becomes observable without disturbing the complete-cell
durable surfaces. This is the seam that lets the GUI approach opencode's
liveness without importing its turn-based drain.

**Surface contract**

| Surface | Role | Mutability | Provider input? |
|---|---|---|---|
| `context/state/current.jsonl` | complete cognition cells | append-only, mutable working memory | **yes** (read back; recovery-critical) |
| `context/state/live.jsonl` | transient generation deltas | append-only, **truncated per generation** | **never** |
| `tape` | durable record | append-only | n/a (replay) |

`live.jsonl` entries are streaming deltas (`text_delta`, `reasoning_delta`,
`tool_call_delta`) stamped with a monotonic `seq`. They are produced by
consuming the provider's SSE response incrementally — the delta-decode logic
already exists in the protocol layer (`extractSSECompletedResponse` and its
`output_text.delta` / `reasoning_summary_text.delta` handlers); the buffering
`io.ReadAll` is removed so deltas emit as they arrive instead of after the full
stream is buffered. The authoritative final cell still comes from the
provider's terminal `completed` event — **shared final-assembly** with the
buffered path, so the crystallized cell is identical whether or not anything
tails the live stream.

**Reproducibility boundary (named explicitly)**

- **Transient layer is non-reproducible by design.** Delta chunking, timing,
  and partial snapshots are ephemeral and were never a reproducibility surface.
  This is the extent of "sacrifice a bit of repeatability for the experience" —
  it costs nothing durable.
- **Crystallized layer stays authoritative.** Tape and `current.jsonl` carry the
  same reproducible cell the buffered path would produce. The streaming decode
  is a *second* path; the concession licenses shipping it without first proving
  bit-identical equivalence to the buffered path in every edge case, **provided
  the final cell is still derived from the authoritative `completed` event.**
- Enabling `stream:true` on providers that today do not stream (Anthropic) is
  permitted under this boundary; their final cell likewise comes from the
  terminal event.

**Invariant (must be enforced)**

`live.jsonl` is never read by the provider-input path. `readContextMessages`
reads `current.jsonl` only; a validator confirms no code path feeds
`live.jsonl` into provider input or recovery. This is what keeps the second
live surface from being Plan-9 drift (a mirror of `current.jsonl`): its role is
sharply distinct — transient display, never input.

> **Implementation spec:** [`STREAMING.md`](./STREAMING.md) — phased build plan,
> test/verification matrix, acceptance criteria, and a full `file:line`
> reference index. Executable end-to-end by a fresh agent.

**Contract increment**

`contract.json` gains a sharply-distinct live surface entry:

```json
"surfaces": {
  "live_context":    "../context/state/current.jsonl",   // complete cells, provider input
  "live_generation": "../context/state/live.jsonl"        // transient deltas, display-only, never input
}
```

## qcli-server Contract

The public HTTP+SSE schema. Both the terminal client and any future web client
consume this. The kernel serves it; it renders nothing.

### Endpoints

```
GET  /events?peer=<id>    SSE: status snapshot, then cognition cells (live tail)
GET  /status?peer=<id>    current status snapshot (JSON)
GET  /context?peer=<id>   full parsed current.jsonl (for reconnect/backfill)
POST /command             a control command (send | interrupt)
GET  /healthz             liveness
```

### Cell event (the core unit)

One cell = one parsed tape entry, stamped with `seq` at the tailer boundary.
Schema-driven: known kinds have typed fields; unknown kinds carry `raw` and
render as a monospace block (forward-compatibility — Quine adding a new cell
kind does not break old clients).

```json
{
  "type": "cell",
  "seq": 42,
  "kind": "message",
  "text": "I'll check the novelty gate",
  "role": "assistant",
  "depth": 2,
  "tool_name": null,
  "tool_id": null,
  "status": null,
  "action": null,
  "author": null,
  "raw": null
}
```

Fields by kind:

| kind | uses |
|---|---|
| `meta` | `depth`, `model` (text) |
| `message` | `text`, `role` |
| `reasoning` | `text` |
| `tool_call` | `tool_name`, `tool_id`, `text` (args as compact JSON) |
| `tool_result` | `tool_id`, `tool_name`, `status`, `text` (output preview), `is_error` |
| `human_post` | `text`, `action`, `author` |
| `reply` | `text` |
| `outcome` | `exit_code`, `turn_count`, `termination_mode` (packed in `text`) |

### Status event

```json
{
  "type": "status",
  "live": true,
  "pid": 48213,
  "session": "9d4d...",
  "model": "gpt-5.5",
  "depth": 2,
  "pending": 0,
  "working": true
}
```

Emitted on `/events` initially (snapshot), then on change. `working` derives
from `current.jsonl` mtime/size advance + pending count + tool activity — it is
the liveness lane signal, not a token stream.

### Command

```json
{
  "op": "send",
  "action": "inject",
  "text": "look at the empty-package bug"
}
```

`action` ∈ `post | poke | inject | interrupt`. No `target` — always the
attached peer (no broadcast in this cut). qcli-server writes to the
corresponding `public/ctl/<action>` file with the standard envelope
(`authority: human`, `reply_required: false`). The POST returns immediately with an ack
(`{"type":"queued","action":"inject"}`) — **never waits for agent reply**.

### Attach / spawn

```
GET /events with no peer  →  needs attach first
POST /command {"op":"attach", "selector": {...}}  →  {"type":"attached", ...}
POST /command {"op":"spawn", "mission": "..."}    →  {"type":"attached", ...}
```

Selectors (unchanged from qctl): `--agent-root`, `--pid`, `--session`,
`--runtime-root`.

## Reserved: Multi-peer Fan-in (not in first cut)

M1 is **single-peer implementation only**. It keeps package boundaries and
reserved contract fields that make later N-peer fan-in straightforward, but it
does not implement N-peer watchers, group scope, broadcast, or author-stamped
merged streams.

When group view returns, the extension path is explicit: add N peer sessions
behind the existing tailer/session abstractions, populate reserved
`author`/`ingest_seq` fields, and add `/events?scope=group`. That is future
kernel + UI work, not hidden M1 behavior.

The rule for the first cut: **reserve seams, do not ship latent group
semantics**. A single-peer client must be excellent before the peer-ecology
surface returns.

## Build Sequence

Contract-first, dependency-ordered. Each step is independently verifiable.

**0. Contract** — freeze the Cell/Status/Command schema above, plus the
`live_generation` surface entry. This is the boundary both sides build against.

**1. Kernel: single-peer contract server (Go)**
- ctl bridge: resolve agent by selector, open client reply endpoint
- surface tailer: `current.jsonl` → cell parser → `/events` SSE
- live tailer: `live.jsonl` → `/events` streaming-delta events (transient,
  display-only; never feeds provider input)
- status tailer: `session.json` → `/status`
- command handler: `/command` → ctl file write (with the **write timeout** fix
  — the old qctl blocked forever on an unread FIFO; this is a real bug to fix
  on day one)
- PTY-free Go tests drive the contract, assert events

**2. Terminal client: individual view (TS + OpenTUI + SolidJS)**
- SSE consumer → SolidJS signal → cell renderer
- composer (IME-safe, action selector, `/switch` `/action` commands)
- cell components per kind (reasoning fold, tool pairing, human track)
- liveness lane — renders the in-progress cell from `live.jsonl` deltas; clears
  on crystallization into a complete cell
- attach on startup (`qcli` = spawn qcli-server + connect TUI)

**3. Live generation stream (runtime + invariant)**
- provider seam: turn `extractSSECompletedResponse` inside-out — emit deltas as
  the live `resp.Body` is scanned, remove the buffering `io.ReadAll`; final
  cell still derived from the authoritative `completed` event (shared
  final-assembly with the buffered path)
- runtime writes deltas to `context/state/live.jsonl` (truncate per
  generation); crystallizes one complete cell to `current.jsonl` + tape on
  completion
- enable `stream:true` on providers that today do not stream (Anthropic)
- **invariant validator**: confirm `live.jsonl` is never read by
  `readContextMessages` / recovery (this is what keeps the second live surface
  from being Plan-9 drift)
- self-source: `runtime.go` is a self-source sentinel; the change must keep the
  quine self-building

**4. Polish toward opencode parity**
- markdown/code rendering (OpenTUI `Markdown` + `Code`)
- scrollback with sticky-bottom
- command palette / fuzzy peer picker (for `/switch`)
- char-frame snapshot tests (`testRender`)

**5. (Future) Group view** — add group-scope kernel fan-in and UI over the
reserved seams; not part of M1/M2.

## Entry points

```
qcli                          attach to default / spawn, launch terminal UI
qcli attach --session <id>    attach to a live agent, launch UI
qcli spawn [mission]          spawn a new agent, attach, launch UI
qcli serve [--http ADDR]      headless: run qcli-server only (for webui / external)
```

`qcli` (no args) prefers a bundled `quine` binary as sibling.
Configuration stays **environment-variables-only** (QUICKSTART.md:151) — no
flags-as-config, no config files.
