// Roster sidebar (WEB.md §5): lineage tree from parent_session, liveness
// glyphs, 10s polling while the tab is visible, click → confirm chip →
// attach → full rebuild, plus the spawn form.
import { createSignal, For, onCleanup, onMount, Show } from "solid-js"
import type { QcliStore, RosterEntry } from "../store"
import { shortSession } from "../store"

interface TreeRow {
  entry: RosterEntry
  depth: number
}

export function rosterTree(peers: RosterEntry[]): TreeRow[] {
  const bySession = new Map(peers.map((p) => [p.session, p]))
  const children = new Map<string, RosterEntry[]>()
  const roots: RosterEntry[] = []
  for (const p of peers) {
    if (p.parent_session && bySession.has(p.parent_session)) {
      const list = children.get(p.parent_session) ?? []
      list.push(p)
      children.set(p.parent_session, list)
    } else {
      roots.push(p)
    }
  }
  const rows: TreeRow[] = []
  const visit = (entry: RosterEntry, depth: number) => {
    rows.push({ entry, depth })
    for (const child of children.get(entry.session) ?? []) visit(child, depth + 1)
  }
  for (const root of roots) visit(root, 0)
  return rows
}

export function Roster(props: { store: QcliStore }) {
  const [confirming, setConfirming] = createSignal<number | null>(null)
  const [spawnOpen, setSpawnOpen] = createSignal(false)
  let missionRef: HTMLInputElement | undefined

  onMount(() => {
    void props.store.fetchRoster()
    const tick = setInterval(() => {
      if (typeof document === "undefined" || document.visibilityState === "visible") void props.store.fetchRoster()
    }, 10_000)
    onCleanup(() => clearInterval(tick))
  })

  async function attach(peer: RosterEntry) {
    setConfirming(null)
    await props.store.command({
      op: "attach",
      selector: { pid: peer.pid, runtime_root: props.store.state.peer?.runtime_root },
    })
    void props.store.fetchRoster()
    void props.store.fetchPeerContract()
  }

  async function spawn() {
    const mission = missionRef?.value?.trim() ?? ""
    if (!mission) return
    await props.store.command({ op: "spawn", mission })
    setSpawnOpen(false)
    void props.store.fetchRoster()
  }

  return (
    <div class="roster">
      <h3>peers</h3>
      <Show when={props.store.state.roster.length > 0} fallback={<div class="meta">no peers — spawn one below, or start a quine</div>}>
        <For each={rosterTree(props.store.state.roster)}>
          {(row) => (
            <div
              class={`peer-row ${row.entry.attached ? "attached" : ""}`}
              style={{ "padding-left": `${6 + row.depth * 14}px` }}
              onClick={() => {
                if (row.entry.attached) return
                setConfirming(confirming() === row.entry.pid ? null : row.entry.pid)
              }}
            >
              <span class={`glyph ${row.entry.live ? "live" : "dead"}`}>{row.entry.live ? "●" : "◍"}</span>
              <span>{shortSession(row.entry.session)}</span>
              <span class="meta">
                {row.entry.model} · d{row.entry.depth}
              </span>
              <Show when={row.entry.pending > 0}>
                <span class="pending">✉{row.entry.pending}</span>
              </Show>
              <Show when={row.entry.attached}>
                <span class="meta">← you</span>
              </Show>
              <Show when={confirming() === row.entry.pid}>
                <button
                  class="confirm-chip"
                  onClick={(e) => {
                    e.stopPropagation()
                    void attach(row.entry)
                  }}
                >
                  switch? discards view
                </button>
              </Show>
            </div>
          )}
        </For>
      </Show>
      <div class="spawn">
        <Show
          when={spawnOpen()}
          fallback={
            <button onClick={() => setSpawnOpen(true)}>+ spawn</button>
          }
        >
          <input ref={missionRef} placeholder="mission…" onKeyDown={(e) => e.key === "Enter" && void spawn()} />
          <button onClick={() => void spawn()}>spawn</button>
          <button onClick={() => setSpawnOpen(false)}>cancel</button>
        </Show>
      </div>
    </div>
  )
}
