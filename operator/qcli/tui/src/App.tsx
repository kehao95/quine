import { createEffect, createMemo, createSignal, For, on, onCleanup, onMount, Show } from "solid-js"
import { useKeyboard, useTerminalDimensions } from "@opentui/solid"
import type { KeyEvent } from "@opentui/core"
import { TextAttributes } from "@opentui/core"
import type { ControlAction } from "./contract"
import { shortSession } from "./contract"
import type { QcliStore, RenderEntry } from "./store"
import { actionDescription, parseSlashCommand, renderEntries } from "./store"
import { COMPOSER_KEY_BINDINGS, nextAction } from "./keymap"
import { theme } from "./theme"
import { DemotedBlockView, Divider, EntryView, StreamingView } from "./components/CellView"
import { HelpOverlay, RosterOverlay } from "./components/Overlays"

export function App(props: { store: QcliStore }) {
  const store = props.store
  const dims = useTerminalDimensions()
  const [action, setAction] = createSignal<ControlAction>("inject")
  const [helpOpen, setHelpOpen] = createSignal(false)
  const [rosterOpen, setRosterOpen] = createSignal(false)
  const [rosterIdx, setRosterIdx] = createSignal(0)
  const [rosterFilter, setRosterFilter] = createSignal("")
  const [armedInterrupt, setArmedInterrupt] = createSignal(false)
  const [history, setHistory] = createSignal<string[]>([])
  const [histIdx, setHistIdx] = createSignal(-1)
  const [quitArmedAt, setQuitArmedAt] = createSignal(0)
  const [focusKey, setFocusKey] = createSignal<string | null>(null)
  const [expandedKeys, setExpandedKeys] = createSignal<ReadonlySet<string>>(new Set<string>())
  let draft = ""
  let inputRef: any
  let scrollRef: any

  onMount(() => {
    const ac = new AbortController()
    void store.connect(ac.signal)
    void store.fetchRoster()
    void store.fetchPeerContract()
    onCleanup(() => ac.abort())
  })

  createEffect(() => {
    if (store.state.peer) void store.fetchPeerContract()
  })

  // Terminal bell on new unread replies/outcomes (G5); opt out via QCLI_NOTIFY=0.
  createEffect(
    on(
      () => store.state.unread,
      (unread, prev) => {
        if (prev !== undefined && unread > prev && process.env.QCLI_NOTIFY !== "0") {
          process.stdout.write("\x07")
        }
      },
      { defer: true },
    ),
  )

  const entries = createMemo<RenderEntry[]>(() => renderEntries(store.state.cells))
  const peerDead = createMemo(() => {
    const s = store.state.statusSnapshot
    return !!s && !s.live
  })

  const title = createMemo(() => {
    const peer = store.state.peer
    const status = store.state.statusSnapshot
    if (!peer) return "qcli · no peer"
    const liveDot = peerDead() ? "◍" : "●"
    if (dims().width < 80) return `qcli · ${shortSession(peer.session)} · ${status?.model ?? "?"} ${liveDot}`
    return `qcli · ${shortSession(peer.session)} · pid=${peer.pid} · ${status?.model ?? "?"} · depth=${status?.depth ?? 0} · inc=${status?.incarnation ?? 0} ${liveDot}`
  })

  const connectionLabel = createMemo(() => {
    if (dims().width < 80) return store.state.connection === "connected" ? "●" : "…"
    return store.state.connection
  })

  const lane = createMemo(() => {
    if (peerDead()) return "◍ peer is dead — transcript read-only"
    const s = store.state.statusSnapshot
    if (!s) return store.state.laneNotice || store.state.connection
    const tokens = s.tokens_in == null || s.tokens_out == null ? "tokens: n/a" : `tokens: ${s.tokens_in}/${s.tokens_out}`
    const unread = store.state.unread > 0 ? ` · ${store.state.unread} unread` : ""
    const notice = store.state.laneNotice ? ` · ${store.state.laneNotice}` : ""
    return `${s.working ? "working" : "idle"} · pending=${s.pending} · ${tokens}${unread}${notice}`
  })

  const footer = createMemo(() => {
    if (armedInterrupt()) return "Esc again = interrupt · any key cancels"
    if (focusKey() !== null) return "[↑↓] cells · [→/←] expand/collapse · [esc] back · [?] help"
    return `[enter] send · [alt+enter] newline · [tab] action=${action()} · [ctrl+p] peers · [?] help`
  })

  // Placeholder shows the peer contract's own words for the action when known
  // (never paraphrase), else a neutral typing hint.
  const composerPlaceholder = createMemo(() => {
    const desc = actionDescription(store.state, action())
    return desc === action() ? "type a message — enter sends" : desc
  })

  const filteredRoster = createMemo(() => {
    const q = rosterFilter().toLowerCase()
    if (!q) return store.state.roster
    return store.state.roster.filter((peer) =>
      [peer.session, String(peer.pid), peer.model, peer.agent_root].some((value) => fuzzyMatch(q, String(value).toLowerCase())),
    )
  })

  function focusableIndex(): number {
    const key = focusKey()
    if (key === null) return -1
    return entries().findIndex((e) => e.key === key)
  }

  function moveFocus(dir: number) {
    const list = entries()
    if (list.length === 0) return
    let idx = focusableIndex()
    if (idx === -1) idx = dir < 0 ? list.length - 1 : 0
    else idx = Math.max(0, Math.min(list.length - 1, idx + dir))
    setFocusKey(list[idx]!.key)
  }

  function toggleExpand(key: string, expand: boolean) {
    setExpandedKeys((prev) => {
      const next = new Set(prev)
      if (expand) next.add(key)
      else next.delete(key)
      return next
    })
  }

  async function sendComposer() {
    if (peerDead()) return
    const text = inputRef?.plainText ?? ""
    const slash = parseSlashCommand(text)
    if (slash.kind === "action") {
      setAction(slash.action)
      clearInput()
      return
    }
    if (slash.kind === "help") {
      setHelpOpen(true)
      clearInput()
      return
    }
    if (slash.kind === "switch") {
      clearInput()
      openRoster()
      return
    }
    if (slash.kind === "quit") {
      process.exit(0)
    }
    if (slash.kind === "spawn") {
      await store.command({ op: "spawn", mission: slash.mission })
      clearInput()
      void store.fetchRoster()
      return
    }

    const a = action()
    if (!text.trim() && a !== "poke" && a !== "interrupt") return
    if (text.trim()) setHistory((h) => [...h.filter((x) => x !== text), text].slice(-50))
    setHistIdx(-1)
    setFocusKey(null)
    const res = await store.command({ op: "send", action: a, text })
    if (res.type !== "error") clearInput()
  }

  async function sendEmptyInterrupt() {
    const text = action() === "interrupt" ? (inputRef?.plainText ?? "") : ""
    const res = await store.command({ op: "send", action: "interrupt", text })
    if (res.type !== "error" && text) clearInput()
  }

  function recall(dir: number) {
    const h = history()
    if (h.length === 0) return
    let idx = histIdx()
    if (idx === -1) {
      if (dir > 0) return
      draft = inputRef?.plainText ?? ""
      idx = h.length - 1
    } else {
      idx += dir
    }
    if (idx >= h.length) {
      setHistIdx(-1)
      setInputText(draft)
      return
    }
    if (idx < 0) idx = 0
    setHistIdx(idx)
    setInputText(h[idx]!)
  }

  function setInputText(text: string) {
    inputRef?.setText?.(text)
    inputRef?.gotoBufferEnd?.()
    inputRef?.focus?.()
  }

  function clearInput() {
    inputRef?.setText?.("")
    inputRef?.clear?.()
  }

  function openRoster() {
    void store.fetchRoster()
    setRosterOpen(true)
    setRosterIdx(0)
    setRosterFilter("")
  }

  async function attachSelected() {
    const peer = filteredRoster()[rosterIdx()]
    if (!peer) return
    await store.command({ op: "attach", selector: { pid: peer.pid, runtime_root: store.state.peer?.runtime_root } })
    setRosterOpen(false)
    setFocusKey(null)
    setExpandedKeys(new Set<string>())
  }

  async function openEditor() {
    const editor = process.env.EDITOR
    if (!editor) return
    const tmp = `${process.env.TMPDIR ?? "/tmp"}/qcli-compose-${Date.now()}.txt`
    await Bun.write(tmp, inputRef?.plainText ?? "")
    const proc = Bun.spawnSync([editor, tmp], { stdin: "inherit", stdout: "inherit", stderr: "inherit" })
    if (proc.exitCode === 0) {
      setInputText(await Bun.file(tmp).text())
    }
  }

  function scrollToTop() {
    scrollRef?.scrollTo?.(0)
    scrollRef?.scrollBy?.(-1e9)
  }

  function scrollToBottom() {
    scrollRef?.scrollBy?.(1e9)
    store.clearUnread()
  }

  // Global keys run BEFORE the focused textarea (subscription order) and gate
  // it via preventDefault. Rule (TUI.md §6): claim only chords the composer
  // never needs; printables and mid-buffer arrows belong to the composer.
  useKeyboard((evt: KeyEvent) => {
    if (evt.eventType && evt.eventType !== "press") return
    if (helpOpen()) {
      evt.preventDefault?.()
      setHelpOpen(false)
      return
    }
    if (rosterOpen()) {
      evt.preventDefault?.()
      if (evt.name === "escape" || (evt.ctrl && evt.name === "p")) {
        setRosterOpen(false)
      } else if (evt.name === "return") {
        void attachSelected()
      } else if (evt.name === "up") {
        setRosterIdx((i) => Math.max(0, i - 1))
      } else if (evt.name === "down") {
        setRosterIdx((i) => Math.min(filteredRoster().length - 1, i + 1))
      } else if (evt.name === "backspace") {
        setRosterFilter((f) => f.slice(0, -1))
        setRosterIdx(0)
      } else if (!evt.ctrl && !evt.meta && evt.sequence && evt.sequence.length === 1) {
        setRosterFilter((f) => f + evt.sequence)
        setRosterIdx(0)
      }
      return
    }
    if (evt.ctrl && evt.name === "c") {
      const now = Date.now()
      if (now - quitArmedAt() < 2000) process.exit(0)
      setQuitArmedAt(now)
      return
    }
    if (evt.ctrl && evt.name === "p") {
      evt.preventDefault?.()
      openRoster()
      return
    }
    if (evt.ctrl && evt.name === "e") {
      evt.preventDefault?.()
      void openEditor()
      return
    }
    if (evt.name === "tab") {
      evt.preventDefault?.()
      setAction((a) => nextAction(a))
      return
    }
    if (evt.name === "escape") {
      evt.preventDefault?.()
      if (armedInterrupt()) {
        setArmedInterrupt(false)
        void sendEmptyInterrupt()
      } else if (focusKey() !== null) {
        setFocusKey(null)
      } else {
        setArmedInterrupt(true)
        setTimeout(() => setArmedInterrupt(false), 2000)
      }
      return
    }
    if (armedInterrupt()) setArmedInterrupt(false)

    if (evt.ctrl && (evt.name === "up" || evt.name === "down")) {
      evt.preventDefault?.()
      moveFocus(evt.name === "up" ? -1 : 1)
      return
    }

    if (focusKey() !== null) {
      if (evt.name === "up" || evt.name === "down") {
        evt.preventDefault?.()
        moveFocus(evt.name === "up" ? -1 : 1)
        return
      }
      if (evt.name === "right" || evt.name === "left") {
        evt.preventDefault?.()
        toggleExpand(focusKey()!, evt.name === "right")
        return
      }
      // Printable input returns control to the composer.
      if (!evt.ctrl && !evt.meta && evt.sequence && evt.sequence.length === 1) {
        setFocusKey(null)
      }
    } else {
      // Boundary-gated history: only claim ↑/↓ when the composer cursor sits
      // on the first/last row; otherwise the textarea moves its own cursor.
      if (evt.name === "up" || evt.name === "down") {
        const row = inputRef?.logicalCursor?.row ?? 0
        const lineCount = inputRef?.lineCount ?? 1
        if (evt.name === "up" && row === 0) {
          evt.preventDefault?.()
          recall(-1)
        } else if (evt.name === "down" && row >= lineCount - 1) {
          evt.preventDefault?.()
          recall(1)
        }
        return
      }
    }

    // Empty-composer printables: ? help, g/G top/bottom.
    const composerEmpty = (inputRef?.plainText ?? "") === ""
    if (composerEmpty && (evt.name === "?" || evt.sequence === "?")) {
      evt.preventDefault?.()
      setHelpOpen(true)
      return
    }
    if (composerEmpty && (evt.name === "g" || evt.sequence === "g" || evt.sequence === "G")) {
      evt.preventDefault?.()
      if (evt.shift || evt.sequence === "G") scrollToBottom()
      else scrollToTop()
      return
    }
    if (evt.name === "pageup") {
      scrollRef?.scrollBy?.(-Math.max(3, Math.floor(dims().height / 2)))
      return
    }
    if (evt.name === "pagedown") {
      scrollRef?.scrollBy?.(Math.max(3, Math.floor(dims().height / 2)))
      return
    }
  })

  return (
    <box flexDirection="column" width="100%" height="100%" backgroundColor={theme.bg}>
      <box flexShrink={0} backgroundColor={theme.headerBg} paddingLeft={1} paddingRight={1} flexDirection="row" justifyContent="space-between">
        <text fg={theme.headerFg} attributes={TextAttributes.BOLD}>
          {title()}
        </text>
        <text fg={theme.headerFg}>{connectionLabel()}</text>
      </box>

      <scrollbox ref={(r: any) => (scrollRef = r)} flexGrow={1} minHeight={0} stickyScroll stickyStart="bottom" paddingLeft={1} paddingRight={1} paddingTop={1}>
        <Show when={store.state.cells.length === 0 && store.state.demoted.length === 0 && !store.state.live}>
          <text fg={theme.dim}>
            {store.state.peer
              ? "waiting for context…"
              : store.state.roster.length === 0
                ? "no peers found — /spawn <mission> to start one, or attach with ctrl+p once a quine is running"
                : "attach a peer: ctrl+p"}
          </text>
        </Show>
        <For each={store.state.demoted}>{(block) => <DemotedBlockView reason={block.reason} cells={block.cells} />}</For>
        <Show when={store.state.demoted.length > 0 && store.state.cells.length > 0}>
          <Divider label="live context" />
        </Show>
        <For each={entries()}>
          {(entry) => <EntryView entry={entry} focused={focusKey() === entry.key} expanded={expandedKeys().has(entry.key)} />}
        </For>
        <Show when={store.state.live}>
          <StreamingView live={store.state.live!} />
        </Show>
      </scrollbox>

      <box flexShrink={0} backgroundColor={theme.panelBg} border={["top"]} borderColor={theme.border} paddingLeft={1} paddingRight={1} flexDirection="row" justifyContent="space-between">
        <text fg={store.state.lastError ? theme.error : peerDead() ? theme.warning : theme.dim}>{lane()}</text>
        <text fg={armedInterrupt() ? theme.warning : theme.dim}>{footer()}</text>
      </box>

      <Show
        when={!peerDead()}
        fallback={
          <box flexShrink={0} paddingLeft={1} paddingTop={1} paddingBottom={1} backgroundColor={theme.panelBg}>
            <text fg={theme.dim}>composer disabled — attach a live peer with ctrl+p</text>
          </box>
        }
      >
        <box flexShrink={0} flexDirection="row" gap={1} paddingLeft={1} paddingRight={1} paddingTop={1} paddingBottom={1} backgroundColor={theme.panelBg}>
          <text fg={theme.accent} flexShrink={0}>{`${action()} >`}</text>
          <textarea
            ref={(r: any) => {
              inputRef = r
              setTimeout(() => r?.focus?.(), 1)
            }}
            flexGrow={1}
            height={3}
            keyBindings={COMPOSER_KEY_BINDINGS as never}
            placeholder={composerPlaceholder()}
            placeholderColor={theme.dim}
            cursorColor={theme.accent}
            onSubmit={() => void sendComposer()}
          />
        </box>
      </Show>

      <Show when={helpOpen()}>
        <HelpOverlay width={dims().width} height={dims().height} />
      </Show>
      <Show when={rosterOpen()}>
        <RosterOverlay width={dims().width} height={dims().height} peers={filteredRoster()} selected={rosterIdx()} filter={rosterFilter()} />
      </Show>
    </box>
  )
}

function fuzzyMatch(query: string, text: string): boolean {
  if (!query) return true
  let i = 0
  for (const ch of text) {
    if (ch === query[i]) i++
    if (i === query.length) return true
  }
  return false
}
