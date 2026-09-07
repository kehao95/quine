# qcli — UX Maturity Specification

> Extends [`DESIGN.md`](./DESIGN.md) from the single-view first cut toward an
> **opencode-grade client experience**. DESIGN.md owns the architecture, the
> contract, and the five interaction deviations; this doc owns the **UX layer**:
> the opencode-parity inventory, interaction detail beyond the first cut,
> failure-state UX, and the maturity milestones. Nothing here changes the three
> design principles — every feature below traces to an existing runtime
> surface, or is explicitly flagged as a **runtime RFC** (§7) rather than
> smuggled in as a client hack.

## 1. What "opencode-mature" decomposes into

"Feels like opencode" is not one property. It decomposes into three buckets,
and the bucket determines what we do with each feature.

### Bucket A — borrow directly (rendering & mechanical polish)

Zero conflict with Quine's shape. Port the pattern, keep the quality bar.

| Feature | Source pattern |
|---|---|
| markdown + syntax-highlighted code cells | OpenTUI `Markdown` + `Code` renderables |
| tool_call ↔ tool_result pairing, collapse/expand | opencode session transcript; revive deleted `App.tsx` pairing logic |
| in-place streaming cell (typewriter) | `live.jsonl` deltas → streaming register (DESIGN.md, landed runtime-side) |
| scrollback + sticky-bottom | OpenTUI `ScrollBox` |
| IME-safe multiline composer | OpenTUI `Textarea` (grapheme/width model; CJK verified by dogfood) |
| command palette + fuzzy pickers | opencode `command-palette` pattern |
| themes, truecolor→256→16 degradation | revive deleted `theme.ts` |
| mouse support (scroll, click-to-expand, selection) | OpenTUI events |
| char-frame snapshot tests | `testRender`/`captureCharFrame`; revive deleted `render.test.tsx` |

### Bucket B — re-derive under Quine's shape

The opencode feature answers a real UX need, but its mechanism assumes a
master-drain client. Keep the need, rebuild the mechanism on Quine's surfaces.

| opencode feature | Why it can't transfer | Quine-native replacement |
|---|---|---|
| session list / resume | sessions are live processes + durable tapes, not client-owned records | peer roster (live/dormant/dead) + lineage; *attach*, not resume (§3 G4) |
| turn spinner ("working…" until reply) | no turns; replies are 0..N, agent-chosen | **delivery receipt ladder** from `log/control.jsonl` (§3 G2) + liveness lane |
| permission dialog | agent owns its lifecycle (deviation #1) | graded delivery actions + guarded interrupt (§3 G6) |
| token/cost HUD | runtime holds usage in memory only; no live surface | runtime RFC: usage projection (§7 R1) |
| completion notification | no completion event exists | volatile unread + bell on reply/outcome while backgrounded (§3 G5) |
| `/undo` `/redo` | context is the agent's working memory, not client state | rejected; honest intervention is `interrupt` + a new message |
| model picker | model belongs to the agent's spawn config | display-only in header; choose at `/spawn` time |

### Bucket C — reject (violates the shape)

Prompt-driven drain · per-tool-call permission · client-owned durable
history/preferences · reply correlation & waiting spinners · inline subagent
expansion. These are the five deviations of DESIGN.md restated as UX
non-goals. Any parity checklist item that requires one of these is dropped.

## 2. The flagship: delivery legibility (what opencode cannot do)

Maturity is not only parity — Quine's surfaces afford one UX element opencode
has no equivalent for. The graded control plane already logs a **delivery
lifecycle** per message (`internal/runtime/control.go:82` `controlLogEntry`;
contract `control_log_events`, `contract_surface.go:112`):

```
received  — the write was accepted, carries message.id + payload
woke      — an idle process was resumed (poke/inject/interrupt)
delivered — the payload was surfaced into the process context, carries message.id
```

`received` and `delivered` carry the stable inbox `message.id`
(`control.go:365,481`), so the client can associate them to a specific
`human_post` cell exactly (match own write by payload+time on `received`,
then track by id). The composer's send therefore renders an honest receipt
ladder on the cell:

```
◀ human_post·inject  look at the empty-package bug   ▷queued ✓received ⚡woke ▣delivered
```

- `post` may honestly stop at `✓received` (queue-only) until the agent
  retrieves its inbox — that is the semantics made visible, not a bug.
- No stage implies a reply is coming. The ladder ends at `delivered`; replies
  remain independent cells (deviation #5).
- This *replaces* the turn spinner: the anxiety opencode solves with a spinner
  ("did it hear me?") is solved with receipts, without faking turns.

Kernel increment: a third tailer on `log/control.jsonl` emitting a
`receipt` SSE event `{message_id, stage, action, ts}`. Pure observation,
zero runtime change.

## 3. Gap decisions (what the first cut leaves open)

### G1 — Transcript continuity under compaction

`current.jsonl` is mutable working memory; the agent may compact it, and the
first cut's transcript would silently shrink. Decisions:

1. **Backfill on attach/reconnect** via `/context` (already in the contract).
2. **Compaction boundary marker**: when the tailer observes truncation/rewrite
   (size or inode regression), render a divider cell —
   `── context compacted · N cells above are no longer the agent's working memory ──`
   — and keep the pre-compaction cells rendered but visually demoted (dim).
   The client shows history *it* witnessed without pretending the agent still
   has it. Volatile: on client restart only the surface truth remains (red
   line: zero persistence — a restarted client honestly shows less).
3. **Tape history register (M3)**: kernel gains `GET /history?peer=&before=` reading
   the durable tape (read-only, lazy, upward scroll). Rendered in a `history`
   register (dimmer, marked `⌁ tape`), sharply distinct from live cells.
   Tape stays out of the live event stream — it is backfill, not tail.

### G2 — Delivery receipts

§2. Contract adds the `receipt` event kind; cells keep `message_id` linkage
client-side only (volatile).

### G3 — Status richness

| Element | Source | Status |
|---|---|---|
| session/pid/ppid/depth/model/incarnation | `status/session.json` (`runtime.go:1620`) | available now |
| pending count | `status/inbox.json` `pending_count` (contract `inbox_schema`) | available now — kernel tails it |
| working/idle | derived: `current.jsonl` advance + tool activity (DESIGN.md status event) | available now |
| world state | `world/status.json` (contract `world_status`) | available; header badge only |
| tokens in/out, context % | **nowhere live** — tape totals are in-memory until outcome (`tape.go:269`) | runtime RFC R1; status bar shows `tokens: n/a` until then, never fakes |

### G4 — Peer navigation (the "session list" answer)

- **Roster overlay** (`/switch` or `ctrl+p`): fuzzy picker over the kernel's
  roster projection (scan `runtime-root/pid/*` → `status/session.json` +
  liveness verify — revive `roster.go` pattern). Each row: handle, liveness
  glyph (`●` live `○` dormant `◍` dead), model, depth, pending, parent handle.
- **Lineage rendering**: indent children under parents (`parent_session`,
  `depth` are already in session.json). This is a *list with lineage hints*,
  not a tree view (first cut stays single-view; fork points marked in-stream
  per DESIGN.md deviation #4).
- **Spawn flow**: `/spawn <mission>` opens a one-field composer; on success
  the client switches to the new peer, transcript starts at its `meta` cell.
- **Empty states**: no runtime root found → guided hint (`quine` not running,
  how to start); dead peer selected → banner + transcript preserved read-only.

### G5 — Attention routing (0..N async replies)

- Unread badge per peer in the roster + a status-bar aggregate
  (`2 unread from peer-9`), incremented on `reply`/`outcome` cells arriving for
  non-focused peers or while the terminal is unfocused.
- Terminal bell / OSC 9 notification, opt-out via env.
- **All volatile, client-side.** qcli-server keeps zero unread state (red
  line). Restart ⇒ all read; acceptable by design.

### G6 — Composer maturity

- Action selector as a cycling chip (`tab`): `post|poke|inject|interrupt`,
  showing the one-line semantic from the peer contract's `control_actions`
  descriptions, fetched through `GET /peer-contract` — the contract is the help
  text (never paraphrase; render the contract's own words).
- Default `inject` (DESIGN.md). `interrupt` is **guarded**: bare `Esc` never
  fires it; `Esc` focuses an armed confirm chip (`Esc again = interrupt`),
  inherited from REDESIGN's broadcast-interrupt confirm rationale.
- Multiline (`shift+enter`), IME-safe, paste-safe; input history (volatile,
  in-memory); `ctrl+e` opens `$EDITOR` for long payloads.
- **No `@file` / `!shell`** — the agent has its own filesystem and shell;
  material goes in the payload or via the agent's world. (REDESIGN deferred
  these; this doc closes them as rejected for the terminal client.)

### G7 — Failure & edge UX (every state renders; nothing hangs silently)

| Condition | Detection | Render |
|---|---|---|
| ctl FIFO not drained | kernel write timeout (day-one fix, DESIGN.md step 1) | send chip → `⚠ unreachable`; cell stays, receipt ladder stalls at `▷queued` |
| peer death | roster liveness verify fails | `◍ dead` banner; transcript read-only; no spinner anywhere |
| SSE drop | client-side | auto-reconnect + full rebuild from `/context` (cheap because zero state); one-line `reconnected · rebuilt` notice |
| incarnation change | `incarnation_id` bump in status | divider cell `── incarnation N started ──` |
| generation dies mid-stream | `live.jsonl` ends without `completed` | streaming cell → `interrupted` register, marked "not persisted" (DESIGN.md) |
| orphaned tool_call | result never arrives | rendered orphaned, no spinner (deviation #5 discipline) |

## 4. Keybindings, themes, config

- **Config stays environment-only** (`QCLI_THEME`, `QCLI_NOTIFY=0`, server
  addr) per QUICKSTART's env-only rule. No config file, no keymap file in v1 —
  a preference store is exactly what the red line forbids growing first.
- **Fixed keymap + discoverability**: `?` opens a cheat-sheet overlay;
  footer shows the 4 most contextual bindings (opencode's pattern).
  Core map: `ctrl+p` roster · `tab` action chip · `enter` send ·
  `shift+enter` newline · `Esc Esc` interrupt · `g/G` top/bottom ·
  `→/←` expand/collapse cell · `ctrl+e` $EDITOR · `?` help.
- **Themes**: revive `theme.ts` (git `88a78d65`); peer identity = stable
  color + handle, colorblind-safe, degrades to 256/16 (REDESIGN #13).

## 5. Parity acceptance checklist (the testable meaning of "mature")

Rendering/input: CJK + IME composition correct under snapshot tests · paste
of 10KB multiline payload intact · 10k-cell transcript scrolls without jank ·
markdown tables/code render inside width budget · mouse scroll + click-expand.

Liveness: delta→paint p95 < 100ms on local runtime · streaming cell snaps to
settled with zero visual jump for `message` cells · liveness lane never shows
content, only status.

Honesty: every failure state in G7 has a rendered state (grep the codebase
for silent `catch`) · no spinner exists anywhere in the binary · receipts
never imply replies · `tokens: n/a` until R1 lands.

Resilience: kill -9 qcli-server mid-stream → client reconnects and rebuilds
identically (snapshot-compare) · agent compaction mid-attach → marker cell,
no crash · dead FIFO → send fails in ≤ write-timeout with rendered error.

## 6. Milestones (revises DESIGN.md build sequence with the UX layer placed)

| M | Contents | Gate |
|---|---|---|
| **M0 contract freeze** | **spec: [`CONTRACT.md`](./CONTRACT.md)** — DESIGN.md schema **+** `receipt` event kind, `stream_delta` event kind, `/history` reserved, roster snapshot schema | both sides sign off; schema checked into repo |
| **M1 kernel** | **spec: [`KERNEL.md`](./KERNEL.md)** — ctl bridge (write timeout day one) · tailers: `current.jsonl`, `live.jsonl`, `status`, **`log/control.jsonl`** · roster projection · `/events` `/context` `/status` `/command` `/healthz` | PTY-free Go contract tests green |
| **M2 TUI core** | **spec: [`TUI.md`](./TUI.md)** — cell pipeline (all kinds + registers: streaming/settled/interrupted/history) · composer + action chip + guarded interrupt · receipt ladder · roster overlay + `/switch` `/spawn` · reconnect rebuild | dogfood: operate a real Quine day-to-day |
| **M3 maturity layer** | tape `/history` + history register · unread/notify · command palette · themes · cheat sheet · mouse · G7 edge matrix complete | §5 checklist green |
| **M4 runtime RFCs** | R1 usage projection → status bar tokens/context% · R2 Anthropic SSE (STREAMING.md Phase 3) | owner decision + validators |
| **M5 web client** | **spec: [`WEB.md`](./WEB.md)** — same contract, shared client-core data layer, kernel serves the app at `/` | §6 of WEB.md green + browser dogfood |

Revival index (all from `git show 88a78d65:<path>`): `internal/qctl/cell.go`
(tape-line → typed cell, 16 tests) · `roster.go` (pid scan projection) ·
`httpserver.go` (SSE server shape) · `qctl.go` (ctl bridge, spawn/attach,
selector resolution — and the FIFO-block bug to fix) ·
`operator/qctl/tui-client/src/{App.tsx,store.tsx,theme.ts,contract.ts}`
(SSE→Solid store, pairing, theme) · `test/render.test.tsx` (snapshot harness).
Revive patterns, rebuild for single-view (DESIGN.md stack note).

## 7. Runtime RFCs (owner decisions — not client work, never client hacks)

- **R1 usage projection.** A per-generation `status/usage.json`
  (`tokens_in/out` cumulative, last-generation delta, context-window estimate).
  Justification is *not* the client: it is peer-legible self-description of
  resource state (any peer can read it; cf. the self-description resource
  ladder line of research). Until it lands the client shows `n/a` — principle
  #1 forbids growing runtime surface *for the client*, so this ships only if
  it stands as agent-facing physics on its own.
- **R2 Anthropic streaming** (STREAMING.md Phase 3, deferred): without it,
  Anthropic-backed peers get settled cells but no streaming register. The
  client renders both honestly; R2 is additive.
