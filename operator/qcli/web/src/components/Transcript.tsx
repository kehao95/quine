// Transcript column (WEB.md §5): kind-filter chips + substring search (pure
// view filters, volatile), demoted blocks under labelled dividers, the live
// streaming cell, sticky-bottom scrolling with a "jump to live" pill.
// Long-list smoothness comes from CSS content-visibility on .entry rather
// than JS windowing (see WEB.md deviation note in the landing report).
import { createEffect, createMemo, createSignal, For, on, Show } from "solid-js"
import type { QcliStore, RenderCell, RenderEntry } from "../store"
import { cellText, renderEntries } from "../store"
import { Divider, EntryView, StreamingView } from "./Cell"

export type FilterGroup = "messages" | "reasoning" | "tools" | "human" | "outcome"

export const FILTER_GROUPS: FilterGroup[] = ["messages", "reasoning", "tools", "human", "outcome"]

export function groupOf(entry: RenderEntry): FilterGroup {
  if (entry.kind === "pair") return "tools"
  switch (entry.cell.kind) {
    case "message":
    case "reply":
      return "messages"
    case "reasoning":
      return "reasoning"
    case "tool_call":
    case "tool_result":
      return "tools"
    case "human_post":
      return "human"
    case "outcome":
      return "outcome"
    default:
      return "messages"
  }
}

export function filterEntries(entries: RenderEntry[], groups: Record<FilterGroup, boolean>, search: string): RenderEntry[] {
  const q = search.trim().toLowerCase()
  return entries.filter((entry) => {
    if (!groups[groupOf(entry)]) return false
    if (!q) return true
    const cells = entry.kind === "pair" ? [entry.call, entry.result].filter(Boolean) as RenderCell[] : [entry.cell]
    return cells.some(
      (c) => cellText(c).toLowerCase().includes(q) || (c.tool_name ?? "").toLowerCase().includes(q),
    )
  })
}

export function Transcript(props: {
  store: QcliStore
  selectedKey: string | null
  onSelect: (key: string, cell: RenderCell) => void
  onAtBottomChange: (atBottom: boolean) => void
}) {
  const [groups, setGroups] = createSignal<Record<FilterGroup, boolean>>({
    messages: true,
    reasoning: true,
    tools: true,
    human: true,
    outcome: true,
  })
  const [search, setSearch] = createSignal("")
  const [atBottom, setAtBottom] = createSignal(true)
  let scrollRef: HTMLDivElement | undefined

  const entries = createMemo(() => filterEntries(renderEntries(props.store.state.cells), groups(), search()))

  const scrollToBottom = () => {
    if (scrollRef) scrollRef.scrollTop = scrollRef.scrollHeight
  }

  // Sticky bottom: follow the live tail unless the reader scrolled away.
  createEffect(
    on(
      () => [props.store.state.cells.length, props.store.state.live?.parts.length ?? 0],
      () => {
        if (atBottom()) queueMicrotask(scrollToBottom)
      },
    ),
  )

  function onScroll() {
    if (!scrollRef) return
    const nearBottom = scrollRef.scrollTop + scrollRef.clientHeight >= scrollRef.scrollHeight - 48
    setAtBottom(nearBottom)
    props.onAtBottomChange(nearBottom)
  }

  return (
    <div class="transcript-col">
      <div class="filterbar">
        <For each={FILTER_GROUPS}>
          {(g) => (
            <button
              class={`chip ${groups()[g] ? "on" : ""}`}
              data-testid={`filter-${g}`}
              onClick={() => setGroups({ ...groups(), [g]: !groups()[g] })}
            >
              {g}
            </button>
          )}
        </For>
        <input placeholder="search…" value={search()} onInput={(e) => setSearch(e.currentTarget.value)} data-testid="search" />
      </div>
      <div class="transcript" ref={scrollRef} onScroll={onScroll} data-testid="transcript">
        <Show when={props.store.state.cells.length === 0 && props.store.state.demoted.length === 0}>
          <div class="divider">{props.store.state.peer ? "waiting for context…" : "attach or spawn a peer"}</div>
        </Show>
        <For each={props.store.state.demoted}>
          {(block) => (
            <>
              <Divider>── context reset: {block.reason} · {block.cells.length} cell{block.cells.length === 1 ? "" : "s"} above are no longer the agent's working memory ──</Divider>
              <For each={filterEntries(renderEntries(block.cells), groups(), search())}>
                {(entry) => <EntryView entry={entry} selectedKey={props.selectedKey} onSelect={props.onSelect} demoted />}
              </For>
              <Divider>── live context ──</Divider>
            </>
          )}
        </For>
        <For each={entries()}>
          {(entry) => <EntryView entry={entry} selectedKey={props.selectedKey} onSelect={props.onSelect} />}
        </For>
        <Show when={props.store.state.live}>
          <StreamingView live={props.store.state.live!} />
        </Show>
        <Show when={!atBottom()}>
          <button
            class="jump-live"
            onClick={() => {
              scrollToBottom()
              setAtBottom(true)
              props.onAtBottomChange(true)
            }}
          >
            jump to live ↓
          </button>
        </Show>
      </div>
    </div>
  )
}
