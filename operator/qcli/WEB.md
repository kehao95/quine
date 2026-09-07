# qcli Web Client — Implementation Specification (M5)

> **Status: ready to execute.** Consumes [`CONTRACT.md`](./CONTRACT.md)
> `qcli/1` **verbatim — zero contract changes** (verified: kernel cells carry
> full untruncated text; truncation was always render-side). Interaction
> semantics inherit [`DESIGN.md`](./DESIGN.md) (five deviations) and
> [`UX.md`](./UX.md); component semantics inherit [`TUI.md`](./TUI.md) §4–§7
> unless overridden here. The web client is the **observation-rich surface**
> (density, detail, always-visible population); the TUI remains the
> terminal-native operating surface. Both are pure projections over the same
> kernel.

## 1. Role — what web adds that a terminal cannot

| Need | Web answer |
|---|---|
| see *everything* about a cell | inspector pane: full text, full tool output, args, receipt timeline with timestamps, raw event JSON |
| population at a glance | always-visible roster sidebar with lineage tree, liveness, pending/unread badges |
| long-session reading | virtualized transcript, client-side kind-filters + text search (volatile, in-memory) |
| rich rendering | real markdown + syntax highlighting, wide code blocks, proper typography |
| native input | browser composer: Enter sends, Shift+Enter newline — no terminal keyboard-protocol caveats |

Non-goals (v1): tape history (`/history` stays M3-reserved), multi-peer
simultaneous streams (kernel is single-attach; switching re-attaches), any
write surface beyond the four ctl actions, authentication (loopback trust
boundary, same as the contract).

## 2. Principles (inherited, restated for the browser)

- **Zero persistence — including the browser.** No localStorage, no cookies,
  no IndexedDB. Refresh = rebuild from `hello → backfill_complete`. Theme
  follows `prefers-color-scheme`. Losing volatile state (unread, filters,
  selection) on refresh is honest, not a defect.
- **Observation-only.** Same as the TUI: no permission dialogs, no reply
  correlation, receipts not spinners, 0..N async replies tolerated.
- **Same-origin serving.** The kernel serves the built app at `/`, so the SSE
  and command endpoints are same-origin — no CORS surface is added to the
  contract.

## 3. Architecture & code layout

```
operator/qcli/
  client-core/            NEW — framework-free shared TS data layer
    contract.ts           (moved from tui/src; the qcli/1 types + tolerance rules)
    reducer.ts            (moved/extracted from tui/src/store.ts: event→state
                           reducer table, renderEntries pairing, receiptSegments,
                           liveView — everything with no OpenTUI/DOM imports)
    fixtures.ts           (moved from tui/test; shared by both test suites)
  tui/                    imports ../client-core (bun resolves relative TS)
  web/                    NEW — SolidJS web app
    index.html  vite.config.ts  package.json (pins)  src/
      main.tsx  App.tsx  connect.ts (fetch-SSE, same lifecycle as TUI)
      components/{Topbar,Roster,Transcript,Cell,Inspector,Composer}.tsx
      theme.css            (semantic tokens mirroring tui theme.ts)
    dist/                  build output (gitignored; built by make)
```

- **client-core extraction is a refactor with a green-tests invariant**: the
  TUI suite (46 tests) must pass unchanged after the move — that suite *is*
  the acceptance test for the extraction. `contract.ts` was spec-required to
  be UI-free (TUI.md §3.1) precisely for this.
- **Kernel increment (only one):** `internal/qcli/server.go` serves static
  files at `/` from a directory resolved as `QCLI_WEB_DIST`, else
  `<qcli-source-dir>/web/dist` (same discovery pattern as the TUI launch);
  when absent, `/` returns a small plain-text index of contract endpoints
  (current behavior preserved as fallback). SPA fallback: unknown non-API
  paths serve `index.html`. API routes take precedence — `/events`,
  `/context`, `/status`, `/roster`, `/command`, `/healthz`, `/peer-contract`,
  `/history` are never shadowed.
- **Entry point:** `qcli web [selector…]` = attach/spawn + serve + print the
  URL (and try `xdg-open`/`open`, non-fatal if absent). `qcli serve` alone
  also serves the app if dist exists — one daemon, both surfaces.
- **Build:** `make build-qcli-web` → `bun install && bunx vite build` in
  `operator/qcli/web`. Vite + `vite-plugin-solid`, versions pinned exact.
  Dev loop: `bunx vite dev` with a proxy of the API routes to a running
  `qcli serve` (config included).

## 4. Layout

```
┌ topbar: ● o_aaaa1111 · pid 48213 · gpt-5.5 · depth 2 · inc 2 │ pending 0 · working │ ⚡connected ┐
├───────────┬──────────────────────────────────────────────┬───────────────────────┤
│ ROSTER    │ TRANSCRIPT (virtualized, sticky-bottom)      │ INSPECTOR             │
│ ● you→7   │  [filters: msg ✓ reasoning ✓ tools ✓] [search]│ (click any cell)      │
│ ● peer-7  │  ◀ you · inject          ▷✓⚡▣  02:06:10     │ kind: tool_result     │
│  └○ p-9   │  ┊ reasoning ▸ 3 lines                       │ tool: bash (call_YB…) │
│ ◍ p-2     │  ┊ message  (markdown, full width)           │ status: ok            │
│           │  ┊ tool_call ▸ bash …                        │ ── full output ──     │
│ [spawn]   │  ┊ └ ✓ ok · 3/5 flagged                      │ (untruncated, mono)   │
│           │  ┊▒ streaming · gen 3 … ▊                    │ ── receipts ──        │
│           ├──────────────────────────────────────────────┤ received 02:06:10.234 │
│           │ [post|poke|inject|interrupt]  [⏸ interrupt]  │ delivered 02:06:10.235│
│           │ ┌ composer (Enter send · Shift+Enter newline)┐ ── raw event ──       │
│           │ └──────────────────────────────────────────┘ │ {…json…}              │
└───────────┴──────────────────────────────────────────────┴───────────────────────┘
```

Inspector and roster collapse on narrow viewports (≤900px: single column,
inspector becomes a modal). The transcript column is the only scroll owner.

## 5. Component behaviors (deltas from TUI.md §4–§7; unlisted = same semantics)

- **Transcript**: virtualized (10k cells smooth); "jump to live ↓" pill when
  scrolled off bottom; kind-filter chips and a substring search box — both
  pure view filters over the in-memory store, cleared on refresh. Registers
  render as in TUI.md §4 (streaming dim + cursor, demoted under divider,
  local echo until reconciled, interrupted "not persisted").
- **Cell → Inspector**: click selects (highlight rail); inspector shows every
  field the contract carries for that cell, untruncated, plus the receipt
  ladder as a timestamped timeline (for `human_post`/`local`), plus the raw
  event JSON (pretty-printed, copy button). Esc deselects.
- **Roster**: rendered from `GET /roster`, refreshed on open/focus and every
  10s while visible (polling is v1 contract behavior); lineage as an
  indented tree via `parent_session`; rows show liveness glyph, model,
  depth, pending; attached peer marked. Click → confirm chip ("switch
  discards this view — zero persistence") → `attach` → full rebuild.
  `[spawn]` opens a one-field mission form → `spawn` op.
- **Composer**: Enter sends, Shift+Enter newline (native browser semantics —
  the TUI's keyboard-protocol caveat does not exist here). Action selector as
  4 radio chips with the contract's own `control_actions` descriptions
  (via `/peer-contract`) as tooltips/hint line. Interrupt is a separate
  guarded button: first click arms (turns red, "click again to interrupt",
  2s auto-disarm), second click sends; with composer text and interrupt
  armed, the text rides along as the urgent payload.
- **Unread/attention** (G5): when `document.hidden` or scrolled off bottom,
  `reply`/`outcome` cells increment a badge in `document.title`
  (`(2) qcli — peer-7`) and the favicon dot; cleared on focus+bottom. No
  Notification API in v1 (permission prompt noise > value on localhost).
- **Connection**: SSE drop → full-width banner (`reconnecting…` →
  `rebuilding…`), then store reset + replay, one-line "reconnected" toast.
  Same backoff as TUI connect.ts (250ms → 4s cap).
- **Dead peer / vanished root**: `status.live=false` → persistent banner
  `◍ peer exited — transcript read-only (last observed state)`, composer
  disabled. (The kernel's zombie-aware liveness landed in `028a8dff`.)
- **Theme**: `prefers-color-scheme` auto; tokens mirror `tui/src/theme.ts`
  semantic names so the palettes stay in lockstep by review.

## 6. Tests (blocking, same discipline as TUI.md §10–§11)

- **client-core**: the moved reducer/fixture tests run in `client-core`
  (they are the extraction's proof); TUI suite passes unchanged.
- **web components** (bun test + happy-dom or vitest+jsdom — pick one, pin
  it): reducer-to-DOM rendering per cell kind + registers; inspector
  full-detail rendering (assert untruncated output for a >10KB tool_result
  fixture); composer Enter/Shift+Enter behavior; guarded interrupt arming;
  filter/search view logic; unread title badge logic.
- **fixture SSE driver**: reuse the shared fixtures through a tiny
  fetch-SSE mock (no live kernel needed in CI).
- **kernel serving test (Go)**: with a temp dist dir, `/` serves index.html,
  unknown path serves index.html (SPA fallback), API routes unshadowed;
  without dist, `/` falls back to the plain-text endpoint index.
- **acceptance**: all suites green (`web`, `client-core`, `tui`, `go`);
  `make build-qcli-web` produces dist; `qcli serve` + `curl -s :PORT/ | grep
  <title>` sanity; dogfood in a real browser against a real agent (send all
  four actions, watch streaming, inspect a cell, switch peers) — noted in
  the landing PR. These boxes are blocking, not aspirational.

## 7. Execution order (single implementing agent, phased)

P1 extract client-core (TUI green) → P2 kernel static serving + Go test →
P3 web app shell (connect, store wiring, transcript+composer) → P4 roster +
inspector + polish → P5 tests to §6, `make build-qcli-web`, dogfood.
