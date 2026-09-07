# qcli Terminal Client — Implementation Specification (M2 → M3)

> **Status: ready to execute after KERNEL.md P3** (the TUI develops against
> `qcli serve` + recorded fixtures; it needs the contract server, not a live
> agent). Consumes [`CONTRACT.md`](./CONTRACT.md) `qcli/1` verbatim.
> Interaction semantics are fixed by [`DESIGN.md`](./DESIGN.md) (five
> deviations) and [`UX.md`](./UX.md) (gap decisions G1–G7, parity checklist);
> this doc turns them into components, state, and tests. Where this doc and
> UX.md disagree on a detail, this doc is normative (it is the later, closer
> projection); flag the disagreement in the PR.

## 1. Stack (pinned — do not re-litigate, see DESIGN.md §Stack)

- Runtime: `bun` (≥1.1). Language: TypeScript, strict.
- `@opentui/core` + `@opentui/solid` + `@opentui/keymap` (pin exact versions
  at first install; OpenTUI is pre-1.0 — record the pins in `package.json`
  and never float), `solid-js`.
- Location: `operator/qcli/tui/` — `src/{index.tsx,contract.ts,store.ts,
  connect.ts,theme.ts,keymap.ts}`, `src/components/*`, `test/*`,
  `dev/capture.tsx` (char-frame harness), `start.sh`.
- Launched by the `qcli` binary with `QCLI_SERVER` in env (KERNEL.md §4.7);
  standalone dev: `QCLI_SERVER=http://127.0.0.1:PORT bun run src/index.tsx`.

## 2. Layout (single view — the only view in this cut)

```
┌ peer-7 · pid=48213 · sess…9d4d · gpt-5.5 · depth=2 · inc=2 · ●live ─────────┐
│ COGNITION  (scrollbox, sticky-bottom)                                        │
│   ◀ human_post·inject  look at the empty-package bug   ▷✓⚡▣                │
│   ┊ message   I'll check the novelty gate                                    │
│   ┊ reasoning ▸ (3 lines)                                                    │
│   ┊ tool_call ▸ bash: awk -f scorer.awk …                                   │
│   ┊ └ tool_result ▸ ok · 3/5 flagged empty                                  │
│   ┊▒ message   checking the score no▊            ← streaming register        │
│   ◀ reply   pasted diff to peer-3                                            │
├ lane ────────────────────────────────────────────────────────────────────────┤
│  working · pending=0 · 2 tool calls this turn          [tokens: n/a]         │
├ composer   action=inject → peer-7 ───────────────────────────────────────────┤
│  > █                                                                          │
└ [tab]action  [ctrl+p]peers  [esc·esc]interrupt  [?]help ─────────────────────┘
```

Header = latest `status` + `hello.peer`. Lane = status only, never content.
Footer = 4 contextual bindings (swap by focus state).

## 3. Data layer

### 3.1 `contract.ts`

Types mirroring CONTRACT.md §3–4 field-for-field, plus the tolerance rules
(§8: unknown type → ignore; unknown kind → raw block; extra fields → ignore).
This file is the client's half of the contract test surface — keep it free of
UI imports so the web client reuses it unchanged.

### 3.2 `connect.ts` — SSE lifecycle

Fetch-based SSE reader (bun-native, no EventSource dependency): parse
`data:` frames, JSON-decode, dispatch to the store. On drop: exponential
backoff (250ms → 4s cap), then **reset the entire store** and replay the new
`hello → backfill_complete` (CONTRACT.md §5). Surface connection state as a
signal (`connected | reconnecting | rebuilding`) rendered in the header dot
and a one-line `reconnected · rebuilt` lane notice.

### 3.3 `store.ts` — one Solid store, reducers per event type

```ts
{
  peer, statusSnapshot, connection,
  cells: Cell[],                 // settled, in arrival order (kernel seq)
  demoted: Cell[][],             // pre-context_reset generations (G1)
  live: { generation, parts } | null,   // in-flight streaming cell
  receipts: Map<client_ref, {stage, message_id}>,
  unread: number,                // volatile (G5)
  backfilled: boolean,
}
```

Reducer rules (normative):

| Event | Rule |
|---|---|
| `hello` | reset all; set peer |
| `status` | patch snapshot; incarnation bump renders divider via `context_reset` (kernel emits it — do not derive it twice) |
| `cell` | append to `cells`; if `kind ∈ {reply, outcome}` and (window unfocused ∨ scrolled off bottom) → `unread++`, bell (G5) |
| `backfill_complete` | `backfilled = true` (before it, suppress unread/bell — replay is not news) |
| `stream_delta` | if `generation` ≠ current → open fresh `live`; append part by kind; `completed` → close `live` (the settled cells follow as `cell` events — **do not** synthesize a settled cell from deltas) |
| `context_reset` | push `cells` onto `demoted`; clear `cells`; render divider (reason-labelled); `backfilled = false` |
| `receipt` | if `client_ref` matches → advance ladder stage; `woke` (no ref) → transient lane note |
| `error` | scope `send` → attach error chip to the pending human_post; else lane banner (G7) |

The store also appends the local echo on send: a `human_post` cell (register
`local`, ladder at `▷queued`) when `/command` acks — reconciled when the same
payload arrives as a delivered `human_post` from current.jsonl (match by
`client_ref` → receipt `message_id`; on match, replace echo with the
authoritative cell). If it never arrives, the echo stays with its stalled
ladder — honest (post may sit queued forever).

## 4. Cell components and registers

One renderer, `kind` selects the component (DESIGN.md pipeline). Registers
(orthogonal to kind): `settled` (full styling) · `streaming` (dim, `▒` rail,
trailing `▊` cursor) · `interrupted` (strikethrough tint, "not persisted"
tag) · `demoted` (dim, under a compaction divider) · `local` (echo, dashed
rail until reconciled) · `history` (reserved, M3 tape backfill).

| kind | component behavior |
|---|---|
| `meta` | not in stream — feeds header |
| `message` | markdown body (OpenTUI `Markdown`); role-colored rail |
| `reasoning` | collapsed: first line + line count; `→`/click expands; streaming: live one-line summary in the collapsed header |
| `tool_call` | collapsed header `tool_name + compact args`; expand → full args as `Code` |
| `tool_result` | indented under its `tool_call` (match `tool_id`; orphan ⇒ render standalone, no spinner); status glyph + preview; expand → full output, `is_error` tinted |
| `human_post` | left `◀` rail, human color, action tag + receipt ladder `▷✓⚡▣` (§5) |
| `reply` | distinct accent (it crossed the wire), `◀ reply` rail |
| `outcome` | terminal line: termination_mode, exit_code, turns, tokens — end of stream marker |
| *unknown* | monospace raw block (CONTRACT.md §8) |

Streaming cell = the same components fed from `live.parts` with register
`streaming`; `tool_call_delta` renders as a forming tool_call (this is the
interrupt-decision view, DESIGN.md deviation #2). A generation that ends
without `completed` flips the streaming cell to `interrupted` in place.

## 5. Receipt ladder (G2)

Rendered inline on `human_post`/`local` cells: `▷ queued → ✓ received →
⚡ woke → ▣ delivered`, dim until lit. `queued` lights on `/command` ack;
`received`/`delivered` light exactly (client_ref match); `⚡` lights when a
`woke` receipt arrives while this ref sits between `received` and `delivered`
— display-level inference, tooltip says so (`woke is peer-level`). Stages
never regress; a `post` honestly parking at `✓` is correct rendering. No
stage ever implies a reply (deviation #5): there is no fifth stage.

## 6. Composer

- OpenTUI `Textarea`: multiline, IME-safe, bracketed paste; 10KB payloads
  must survive intact (parity checklist).
- **Enter-submit is an explicit override, not a default** *(post-mortem:
  shipped broken once because this was unstated)*. OpenTUI's Textarea default
  keymap binds `return → newline` and only `meta+return → submit`, so a bare
  `<textarea onSubmit>` never fires on plain Enter. The composer MUST pass
  `keyBindings` (merged over defaults by `mergeKeyBindings`):

  ```tsx
  keyBindings={[
    { name: "return", action: "submit" },
    { name: "kpenter", action: "submit" },
    { name: "return", shift: true, action: "newline" }, // kitty-protocol terminals only
    { name: "return", meta: true,  action: "newline" }, // alt+enter fallback, works everywhere
  ]}
  ```

  Terminal reality: plain terminals send identical bytes for Enter and
  Shift+Enter — `shift+return` is only distinguishable under the kitty
  keyboard protocol. Alt+Enter is the portable newline; bracketed paste
  always preserves multiline. The footer/help must advertise `alt+enter`
  for newline, not shift+enter, when the protocol is absent.
- **Driven-key test (blocking)**: a test must simulate a plain `return`
  keypress on the focused composer and assert a `POST /command {op:"send"}`
  is issued. This is the regression gate for the exact failure above.
- **Global key handler defers to composer focus.** `useKeyboard` handlers may
  claim only chords the composer never needs (`ctrl+*`, `esc`, `tab`,
  `pgup/pgdn`, `?` **only when composer is empty or unfocused**). Printable
  characters — including `?` — and arrow keys while editing multiline text
  belong to the composer; `↑/↓` recall history only at the first/last row
  (cursor-position check), never mid-buffer.
- **Action chip** (`tab` cycles): `post | poke | inject | interrupt`, default
  `inject`. The chip's expanded hint shows the action's `description` from the
  peer's public contract `control_actions` — fetched once per attach through
  `GET /peer-contract` (never paraphrase; render the contract's words).
- **Guarded interrupt**: `esc` arms a confirm chip (`esc again = interrupt ·
  any key cancels`, 2s auto-disarm). Bare `esc` never fires. With text in the
  composer + action=interrupt, send carries the payload (urgent delivery);
  empty interrupt sends the empty write.
- Input history (`↑`/`↓` at boundary rows), in-memory only. `ctrl+e` opens
  `$EDITOR` (suspend TUI, temp file in `$TMPDIR`, read back on exit).
- Send = POST `/command`; on ack, local echo (§3.3); on error, error chip on
  the echo with `code` text — composer content is **not** cleared on failure.

## 7. Roster overlay & commands

- `ctrl+p` / `/switch`: fuzzy picker over `GET /roster` (fetched on open —
  v1 polls, CONTRACT.md §9). Row: liveness glyph (`●`/`○`/`◍`) · handle
  (short session) · model · depth · pending · `you`-marker on attached. Enter
  ⇒ `attach` command ⇒ full store reset and rebuild. Lineage rendered as
  indentation via `parent_session` chains.
- `/spawn <mission>` ⇒ `spawn` command, then switch to the new peer.
- `/action <a>` sets the chip; `/help` = `?`; `/quit` = `ctrl+c ctrl+c`
  (double, guarded — single `ctrl+c` copies selection if any, else hints).
- Dead peer attach: banner `◍ peer is dead — transcript read-only`, composer
  disabled, roster hint. Empty roster: guided empty state (how to spawn).

## 8. Keymap (single source: `keymap.ts`, rendered by `?`)

`tab` action chip · `enter` send (explicit keyBindings override, §6) ·
`alt+enter` newline (`shift+enter` where kitty protocol allows) · `esc esc`
interrupt (armed) · `ctrl+p` roster · `ctrl+e` $EDITOR · `↑/↓` history /
cell focus · `→/←` expand/collapse focused cell · `pgup/pgdn` scroll ·
`g/G` top/bottom · `?` help overlay · `ctrl+c ctrl+c` quit. Mouse: wheel
scroll, click to expand/collapse. No user keymap config in v1 (UX.md §4).

## 9. Theme (`theme.ts`)

Revive from `git show 88a78d65:operator/qctl/tui-client/src/theme.ts`; keep
truecolor→256→16 degradation and the colorblind-safe accent set. Registers
map to opacity/tint layers over the base theme. `QCLI_THEME=dark|light`
(default dark); no theme files.

## 10. Tests

- **Reducer unit tests** (`bun test`): every rule in §3.3's table, including
  stream generation rollover, context_reset demotion, echo reconciliation,
  unread gating on `backfilled`.
- **Char-frame snapshot tests**: revive the `testRender`/`captureCharFrame`
  harness (`git show 88a78d65:operator/qctl/tui-client/test/render.test.tsx`,
  `dev/capture.tsx`). One snapshot per: each cell kind settled · reasoning
  collapsed/expanded · tool pair collapsed/expanded/orphaned · streaming
  message mid-flight · interrupted register · receipt ladder at each stage ·
  compaction divider + demoted block · dead-peer banner · roster overlay ·
  armed interrupt chip · help overlay.
- **Contract fixture driver**: a fake SSE server (bun) replaying recorded
  kernel fixtures (share `internal/qcli/testdata/` recordings via a small
  generator script) — the TUI's integration test runs against it headlessly.
- **IME/CJK**: snapshot with CJK + combining chars in composer and cells
  (width model); manual dogfood note for actual IME composition (OpenTUI
  concern, verify on first install).

## 11. Acceptance criteria (M2 gate)

**These boxes are blocking, not aspirational: an implementation that ships
without every box checked is incomplete by definition, and the dogfood box is
the one that catches what tests structurally cannot (the first M2 attempt
shipped an Enter key that could not send — one minute of dogfood would have
caught it). §4's component table and §10's snapshot list are similarly
normative: "renders the kind as truncated plain text" does not satisfy a row
that says markdown, collapse/expand, or pairing.**

- [ ] All §10 tests green in CI (`bun test` wired into the repo's test
      entry point alongside Go tests), including the §6 driven-key
      composer test.
- [ ] Fixture-driven session renders the full §2 layout; scrolling 10k cells
      stays responsive (measure: frame budget under `dev/capture` replay).
- [ ] Reconnect kill-test: kill `qcli serve` mid-stream → TUI shows
      reconnecting → rebuilt state identical (snapshot compare) (G7).
- [ ] No spinner exists anywhere; every G7 state has a snapshot.
- [ ] Dogfood: operate a real Quine for a working session (send all four
      actions, watch a streaming generation, receive an async reply) — noted
      in the landing PR.
- [ ] M3 items (tape history register, command palette, notifications
      beyond bell) remain unstarted here — they land against `/history` and
      UX.md G1/G5 after M2 ships.
