// Per-kind cell rendering with registers (TUI.md §4 semantics, WEB.md §5).
// One renderer; kind selects the body, register selects the visual layer.
import { createSignal, For, Match, Show, Switch } from "solid-js"
import type { JSX } from "solid-js"
import type { LiveState, RenderCell, RenderEntry } from "../store"
import { cellText, liveView, receiptSegments, shortSession } from "../store"
import { renderMarkdown } from "../md"

function authorLabel(cell: RenderCell): { label: string; human: boolean } {
  if (cell.author === "human" || cell.kind === "human_post") return { label: "you", human: true }
  if (cell.kind === "reply") return { label: shortSession(cell.author) || "peer", human: false }
  return { label: cell.author ? shortSession(cell.author) : "agent", human: false }
}

export function Ladder(props: { cell: RenderCell }) {
  return (
    <span class="ladder" title="delivery receipts: queued · received · woke · delivered">
      [
      <For each={receiptSegments(props.cell.receipt)}>
        {(seg) => <span class={seg.lit ? "lit" : "unlit"}>{seg.lit ? seg.glyph : "·"}</span>}
      </For>
      ]
    </span>
  )
}

function Markdown(props: { text: string }) {
  return <div class="md" innerHTML={renderMarkdown(props.text)} />
}

function CellBody(props: { cell: RenderCell }) {
  const c = () => props.cell
  const [expanded, setExpanded] = createSignal(false)
  return (
    <Switch
      fallback={
        <div class="unknown-cell">
          <span class="kind-tag">{c().kind}</span>
          <pre class="block raw">{c().raw ?? cellText(c())}</pre>
        </div>
      }
    >
      <Match when={c().kind === "message"}>
        <Show when={c().register === "interrupted"}>
          <div class="not-persisted">▒ interrupted — not persisted</div>
        </Show>
        <Markdown text={cellText(c())} />
      </Match>
      <Match when={c().kind === "reply"}>
        <div class="reply-cell">
          <span class="kind-tag">reply</span>
          <Markdown text={cellText(c())} />
        </div>
      </Match>
      <Match when={c().kind === "human_post"}>
        <div class="human-post">
          <span class="action-tag">◀ {c().action ?? "post"}</span>
          <Ladder cell={c()} />
          <Show when={c().local}>
            <span class="sending"> (sending)</span>
          </Show>
          <Show when={c().error_code}>
            <div class="send-error">send error: {c().error_code}</div>
          </Show>
          <div style={{ "white-space": "pre-wrap" }}>{c().text ?? ""}</div>
        </div>
      </Match>
      <Match when={c().kind === "reasoning"}>
        <div class="reasoning">
          <div
            class="head"
            onClick={(e) => {
              e.stopPropagation()
              setExpanded(!expanded())
            }}
          >
            {expanded() ? "▾" : "▸"} reasoning · {cellText(c()).split("\n").length} lines
            <Show when={!expanded()}> — {firstLine(cellText(c()))}</Show>
          </div>
          <Show when={expanded()}>
            <div style={{ "white-space": "pre-wrap" }}>{cellText(c())}</div>
          </Show>
        </div>
      </Match>
      <Match when={c().kind === "outcome"}>
        <div class="outcome-cell">
          <span class="kind-tag">outcome</span>
          {c().termination_mode ?? "done"} · exit={c().exit_code ?? "n/a"} · turns={c().turn_count ?? "n/a"} · tokens=
          {c().tokens_in ?? "n/a"}/{c().tokens_out ?? "n/a"}
        </div>
      </Match>
      <Match when={c().kind === "tool_result"}>
        <ToolResultBody cell={c()} />
      </Match>
    </Switch>
  )
}

function firstLine(s: string): string {
  const line = s.split("\n", 1)[0] ?? ""
  return line.length > 120 ? line.slice(0, 119) + "…" : line
}

function ToolCallBody(props: { call: RenderCell; result: RenderCell | null }) {
  const [expanded, setExpanded] = createSignal(false)
  const args = () => {
    const t = props.call.text ?? ""
    if (!expanded()) return t.length > 160 ? t.slice(0, 159) + "…" : t
    try {
      return JSON.stringify(JSON.parse(t), null, 2)
    } catch {
      return t
    }
  }
  return (
    <div class="tool-call">
      <div
        class="head"
        onClick={(e) => {
          e.stopPropagation()
          setExpanded(!expanded())
        }}
      >
        {expanded() ? "▾" : "▸"} <span class="kind-tag">tool</span>
        {props.call.tool_name ?? "tool"}
      </div>
      <div class="args">{args()}</div>
      <Show when={props.result} fallback={<div class="tool-result"><span class="status">└ · no result</span></div>}>
        <ToolResultBody cell={props.result!} nested />
      </Show>
    </div>
  )
}

function ToolResultBody(props: { cell: RenderCell; nested?: boolean }) {
  const c = () => props.cell
  const [expanded, setExpanded] = createSignal(false)
  const text = () => c().text ?? ""
  const preview = () => {
    const lines = text().split("\n")
    if (expanded() || lines.length <= 6) return text()
    return lines.slice(0, 6).join("\n") + `\n… (+${lines.length - 6} more lines — click to expand)`
  }
  return (
    <div
      class="tool-result"
      onClick={(e) => {
        e.stopPropagation()
        setExpanded(!expanded())
      }}
    >
      <span class={`status ${c().is_error ? "err" : "ok"}`}>
        └ {c().is_error ? "✗" : "✓"} {c().status ?? (c().is_error ? "error" : "ok")}
        <Show when={!props.nested}> · {c().tool_name ?? c().tool_id ?? "tool"} (orphan result)</Show>
      </span>
      <pre class="block">{preview()}</pre>
    </div>
  )
}

export function EntryView(props: {
  entry: RenderEntry
  selectedKey: string | null
  onSelect: (key: string, cell: RenderCell) => void
  demoted?: boolean
}) {
  const cell = () => (props.entry.kind === "pair" ? props.entry.call : props.entry.cell)
  const register = () => (props.demoted ? "demoted" : (cell().register ?? "settled"))
  const author = () => authorLabel(cell())
  return (
    <div
      class={`entry register-${register()} ${props.selectedKey === props.entry.key ? "selected" : ""}`}
      data-key={props.entry.key}
      data-kind={cell().kind}
    >
      <div class={`author ${author().human ? "human" : ""}`}>{author().label}</div>
      <div class="body" onClick={() => props.onSelect(props.entry.key, cell())}>
        <Show when={props.entry.kind === "pair"} fallback={<CellBody cell={cell()} />}>
          <ToolCallBody call={(props.entry as Extract<RenderEntry, { kind: "pair" }>).call} result={(props.entry as Extract<RenderEntry, { kind: "pair" }>).result} />
        </Show>
      </div>
    </div>
  )
}

export function StreamingView(props: { live: LiveState }) {
  const view = () => liveView(props.live)
  return (
    <div class="entry streaming" data-kind="streaming">
      <div class="author">▒</div>
      <div class="body">
        <div class="kind-tag">streaming · gen {props.live.generation}</div>
        <Show when={view().reasoning}>
          <div class="reasoning">▸ thinking: {view().reasoning}</div>
        </Show>
        <Show when={view().text}>
          <span style={{ "white-space": "pre-wrap" }}>{view().text}</span>
        </Show>
        <For each={view().tools}>
          {(tool) => (
            <div class="tool-call">
              <div class="head">▸ forming: {tool.name || tool.id}</div>
              <div class="args">{tool.args}</div>
            </div>
          )}
        </For>
        <span class="cursor">&nbsp;</span>
      </div>
    </div>
  )
}

export function Divider(props: { children: JSX.Element }) {
  return <div class="divider">{props.children}</div>
}
