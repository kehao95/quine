// qcli web client — three-pane observation surface over the qcli/1 contract
// (WEB.md §4). Pure projection: zero browser persistence; refresh rebuilds
// from hello → backfill_complete.
import { createEffect, createMemo, createSignal, onCleanup, onMount, Show } from "solid-js"
import type { QcliStore, RenderCell } from "./store"
import { renderEntries } from "./store"
import { computeTitle, setFaviconDot } from "./title"
import { Topbar } from "./components/Topbar"
import { Roster } from "./components/Roster"
import { Transcript } from "./components/Transcript"
import { Inspector } from "./components/Inspector"
import { Composer } from "./components/Composer"

export function App(props: { store: QcliStore }) {
  const store = props.store
  const [selectedKey, setSelectedKey] = createSignal<string | null>(null)
  const [atBottom, setAtBottom] = createSignal(true)
  const [visible, setVisible] = createSignal(typeof document === "undefined" || document.visibilityState === "visible")

  onMount(() => {
    const ac = new AbortController()
    void store.connect(ac.signal)
    void store.fetchPeerContract()
    const onVisibility = () => setVisible(document.visibilityState === "visible")
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setSelectedKey(null)
    }
    document.addEventListener("visibilitychange", onVisibility)
    window.addEventListener("focus", onVisibility)
    window.addEventListener("keydown", onKey)
    onCleanup(() => {
      ac.abort()
      document.removeEventListener("visibilitychange", onVisibility)
      window.removeEventListener("focus", onVisibility)
      window.removeEventListener("keydown", onKey)
    })
  })

  createEffect(() => {
    if (store.state.peer) void store.fetchPeerContract()
  })

  // Unread attention (G5): badge accumulates while hidden or scrolled away;
  // clears when the reader is back at the live tail.
  createEffect(() => {
    if (visible() && atBottom() && store.state.unread > 0) store.clearUnread()
  })
  createEffect(() => {
    if (typeof document === "undefined") return
    document.title = computeTitle(store.state.unread, store.state.peer?.session)
    setFaviconDot(store.state.unread > 0)
  })

  // The inspector shows the LIVE cell (receipts keep advancing after
  // selection), so selection is by entry key, resolved against the store.
  const selectedCell = createMemo<RenderCell | null>(() => {
    const key = selectedKey()
    if (!key) return null
    for (const entry of renderEntries(store.state.cells)) {
      if (entry.key !== key) continue
      return entry.kind === "pair" ? entry.call : entry.cell
    }
    for (const block of store.state.demoted) {
      for (const entry of renderEntries(block.cells)) {
        if (entry.key !== key) continue
        return entry.kind === "pair" ? entry.call : entry.cell
      }
    }
    return null
  })

  const peerDead = () => store.state.statusSnapshot?.live === false

  return (
    <div class="app">
      <Topbar store={store} />
      <Show when={store.state.connection === "reconnecting" || store.state.connection === "rebuilding"}>
        <div class="banner reconnect">{store.state.connection}… the view rebuilds from the kernel (zero client state)</div>
      </Show>
      <Show when={peerDead()}>
        <div class="banner dead">◍ peer exited — transcript read-only (last observed state)</div>
      </Show>
      <div class={`main ${selectedCell() ? "" : "no-inspector"}`}>
        <Roster store={store} />
        <Transcript store={store} selectedKey={selectedKey()} onSelect={(key) => setSelectedKey(key)} onAtBottomChange={setAtBottom} />
        <Show when={selectedCell()}>
          <Inspector cell={selectedCell()!} onClose={() => setSelectedKey(null)} />
        </Show>
      </div>
      <div class="lane">{store.state.laneNotice || " "}</div>
      <Composer store={store} disabled={peerDead()} />
    </div>
  )
}
