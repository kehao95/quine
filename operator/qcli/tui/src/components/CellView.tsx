import { For, Match, Show, Switch, createMemo } from "solid-js"
import { TextAttributes } from "@opentui/core"
import { useTerminalDimensions } from "@opentui/solid"
import type { LiveState, RenderCell, RenderEntry } from "../store"
import { liveView, receiptSegments, renderEntries } from "../store"
import { cellText, shortSession } from "../contract"
import { authorColor, syntaxStyle, theme, tint, truncate } from "../theme"

const AUTHOR_COL = 10
const RESULT_LINE_CAP = 40

export interface EntryViewProps {
  entry: RenderEntry
  focused?: boolean
  expanded?: boolean
  demoted?: boolean
}

// EntryView renders one transcript entry (a standalone cell or a tool pair).
export function EntryView(props: EntryViewProps) {
  return (
    <Switch>
      <Match when={props.entry.kind === "pair"}>
        <PairView
          call={(props.entry as Extract<RenderEntry, { kind: "pair" }>).call}
          result={(props.entry as Extract<RenderEntry, { kind: "pair" }>).result}
          focused={props.focused}
          expanded={props.expanded}
          demoted={props.demoted}
        />
      </Match>
      <Match when={props.entry.kind === "cell"}>
        <CellView cell={(props.entry as Extract<RenderEntry, { kind: "cell" }>).cell} focused={props.focused} expanded={props.expanded} demoted={props.demoted} />
      </Match>
    </Switch>
  )
}

function useBudget() {
  const dims = useTerminalDimensions()
  return createMemo(() => Math.max(24, dims().width - AUTHOR_COL - 24))
}

interface RowProps {
  author: string | null | undefined
  railColor: string
  focused?: boolean
  demoted?: boolean
  children?: unknown
}

// Row is the shared cell chrome: author column + left-rail bordered body.
function Row(props: RowProps) {
  const fg = () => (props.demoted ? theme.dim : authorColor(props.author))
  const rail = () => (props.focused ? theme.accent : props.demoted ? theme.borderSubtle : props.railColor)
  const who = () => {
    const base = !props.author || props.author === "human" ? "you" : shortSession(props.author)
    return props.focused ? `› ${base}` : base
  }
  return (
    <box flexDirection="row" gap={1} marginTop={1}>
      <text width={AUTHOR_COL} flexShrink={0} fg={fg()} attributes={props.focused ? TextAttributes.BOLD : undefined}>
        {props.author === undefined || props.author === null || props.author === "" ? "" : who()}
      </text>
      <box flexGrow={1} minWidth={0} border={["left"]} borderColor={rail()} paddingLeft={1} flexDirection="column">
        {props.children as never}
      </box>
    </box>
  )
}

function expandHint(expanded: boolean | undefined, collapsible: boolean): string {
  if (!collapsible) return ""
  return expanded ? "  (← collapse)" : "  (→ expand)"
}

export function CellView(props: { cell: RenderCell; focused?: boolean; expanded?: boolean; demoted?: boolean }) {
  const c = () => props.cell
  const budget = useBudget()
  const demoted = () => props.demoted || c().register === "demoted"
  const interrupted = () => c().register === "interrupted"
  const bodyFg = () => (demoted() ? theme.dim : theme.text)

  return (
    <Switch fallback={<UnknownCell cell={c()} focused={props.focused} demoted={demoted()} budget={budget()} expanded={props.expanded} />}>
      <Match when={c().kind === "message" && interrupted()}>
        <Row author={c().author ?? "agent"} railColor={tint(theme.error, theme.bg, 0.35)} focused={props.focused} demoted={demoted()}>
          <text fg={theme.error}>{`▒ interrupted — not persisted`}</text>
          <text fg={tint(theme.text, theme.bg, 0.35)} wrapMode="word">
            {cellText(c())}
          </text>
        </Row>
      </Match>
      <Match when={c().kind === "message"}>
        <Row author={c().author ?? "agent"} railColor={theme.border} focused={props.focused} demoted={demoted()}>
          <Show when={demoted()} fallback={<markdown content={cellText(c())} syntaxStyle={syntaxStyle()} fg={bodyFg()} />}>
            <text fg={theme.dim} wrapMode="word">
              {cellText(c())}
            </text>
          </Show>
        </Row>
      </Match>
      <Match when={c().kind === "reply"}>
        <Row author={c().author ?? "peer"} railColor={theme.secondary} focused={props.focused} demoted={demoted()}>
          <text fg={demoted() ? theme.dim : theme.secondary} attributes={TextAttributes.BOLD}>
            {`◀ reply`}
          </text>
          <Show when={demoted()} fallback={<markdown content={cellText(c())} syntaxStyle={syntaxStyle()} fg={bodyFg()} />}>
            <text fg={theme.dim} wrapMode="word">
              {cellText(c())}
            </text>
          </Show>
        </Row>
      </Match>
      <Match when={c().kind === "human_post"}>
        <HumanPost cell={c()} focused={props.focused} demoted={demoted()} />
      </Match>
      <Match when={c().kind === "reasoning"}>
        <Reasoning cell={c()} focused={props.focused} expanded={props.expanded} demoted={demoted()} budget={budget()} />
      </Match>
      <Match when={c().kind === "tool_call"}>
        <PairView call={c()} result={null} focused={props.focused} expanded={props.expanded} demoted={demoted()} />
      </Match>
      <Match when={c().kind === "tool_result"}>
        <OrphanResult cell={c()} focused={props.focused} expanded={props.expanded} demoted={demoted()} budget={budget()} />
      </Match>
      <Match when={c().kind === "outcome"}>
        <Row author="" railColor={theme.warning} focused={props.focused} demoted={demoted()}>
          <text fg={demoted() ? theme.dim : theme.warning}>
            {`■ outcome · ${c().termination_mode ?? "done"} · exit=${c().exit_code ?? "n/a"} · turns=${c().turn_count ?? "n/a"} · tokens=${c().tokens_in ?? "n/a"}/${c().tokens_out ?? "n/a"}`}
          </text>
        </Row>
      </Match>
    </Switch>
  )
}

function HumanPost(props: { cell: RenderCell; focused?: boolean; demoted?: boolean }) {
  const c = () => props.cell
  const rail = () => (c().local ? tint(theme.accent, theme.bg, 0.5) : theme.accent)
  return (
    <Row author="human" railColor={rail()} focused={props.focused} demoted={props.demoted}>
      <box flexDirection="row" gap={1}>
        <text fg={props.demoted ? theme.dim : theme.accent} attributes={TextAttributes.BOLD}>
          {`◀ ${c().action ?? "post"}${c().local ? " (sending)" : ""}`}
        </text>
        <Ladder cell={c()} demoted={props.demoted} />
      </box>
      <Show when={c().error_code}>
        <text fg={theme.error}>{`⚠ send failed: ${c().error_code}`}</text>
      </Show>
      <text fg={props.demoted ? theme.dim : theme.text} wrapMode="word">
        {c().text ?? ""}
      </text>
    </Row>
  )
}

// Ladder renders the receipt stages: lit stages keep their glyph and color,
// unlit stages render as `·` (char-testable, TUI.md §5).
function Ladder(props: { cell: RenderCell; demoted?: boolean }) {
  const segments = () => receiptSegments(props.cell.receipt)
  return (
    <box flexDirection="row">
      <text fg={theme.dim}>[</text>
      <For each={segments()}>
        {(seg) => <text fg={props.demoted ? theme.dim : seg.lit ? (seg.label === "delivered" ? theme.success : theme.accent) : theme.dim}>{seg.lit ? seg.glyph : "·"}</text>}
      </For>
      <text fg={theme.dim}>]</text>
    </box>
  )
}

function Reasoning(props: { cell: RenderCell; focused?: boolean; expanded?: boolean; demoted?: boolean; budget: number }) {
  const text = () => cellText(props.cell)
  const lines = () => text().split("\n")
  const collapsible = () => lines().length > 1 || text().length > props.budget
  return (
    <Row author={props.cell.author ?? "agent"} railColor={theme.borderSubtle} focused={props.focused} demoted={props.demoted}>
      <Show
        when={props.expanded && collapsible()}
        fallback={
          <text fg={theme.dim} attributes={TextAttributes.ITALIC}>
            {`${collapsible() ? "▸" : "·"} reasoning · ${lines().length} line${lines().length === 1 ? "" : "s"} — ${truncate(lines()[0] ?? "", props.budget)}${expandHint(props.expanded, collapsible())}`}
          </text>
        }
      >
        <text fg={theme.dim} attributes={TextAttributes.ITALIC}>{`▾ reasoning${expandHint(true, true)}`}</text>
        <text fg={theme.textMuted} attributes={TextAttributes.ITALIC} wrapMode="word">
          {text()}
        </text>
      </Show>
    </Row>
  )
}

export function PairView(props: { call: RenderCell; result: RenderCell | null; focused?: boolean; expanded?: boolean; demoted?: boolean }) {
  const budget = useBudget()
  const call = () => props.call
  const name = () => call().tool_name ?? "tool"
  const args = () => call().text ?? ""
  return (
    <Row author={call().author ?? "agent"} railColor={theme.border} focused={props.focused} demoted={props.demoted}>
      <Show
        when={props.expanded}
        fallback={
          <text fg={props.demoted ? theme.dim : theme.textMuted}>
            {`▸ ${name()}  ${truncate(compactOneLine(args()), budget())}${expandHint(false, true)}`}
          </text>
        }
      >
        <text fg={props.demoted ? theme.dim : theme.textMuted}>{`▾ ${name()}${expandHint(true, true)}`}</text>
        <Show when={args()}>
          <code content={prettyArgs(args())} filetype="json" syntaxStyle={syntaxStyle()} />
        </Show>
      </Show>
      <Show when={props.result} fallback={<text fg={theme.dim}>{`└ · no result`}</text>}>
        <ResultBody result={props.result!} expanded={props.expanded} demoted={props.demoted} budget={budget()} />
      </Show>
    </Row>
  )
}

function ResultBody(props: { result: RenderCell; expanded?: boolean; demoted?: boolean; budget: number }) {
  const r = () => props.result
  const glyph = () => (r().is_error ? "✗" : "✓")
  const statusFg = () => (props.demoted ? theme.dim : r().is_error ? theme.error : theme.success)
  const out = () => r().text ?? ""
  const lines = () => out().split("\n")
  const statusLabel = () => `${glyph()} ${r().status ?? (r().is_error ? "error" : "ok")}`
  return (
    <Show
      when={props.expanded}
      fallback={
        <box flexDirection="row">
          <text fg={theme.dim} flexShrink={0}>{`└ `}</text>
          <text fg={statusFg()} flexShrink={0}>
            {statusLabel()}
          </text>
          <text fg={theme.dim}>{` · ${truncate(compactOneLine(out()), props.budget)}`}</text>
        </box>
      }
    >
      <box flexDirection="row">
        <text fg={theme.dim} flexShrink={0}>{`└ `}</text>
        <text fg={statusFg()}>{statusLabel()}</text>
      </box>
      <box paddingLeft={2} flexDirection="column">
        <text fg={r().is_error ? tint(theme.error, theme.text, 0.4) : props.demoted ? theme.dim : theme.textMuted} wrapMode="word">
          {lines().slice(0, RESULT_LINE_CAP).join("\n")}
        </text>
        <Show when={lines().length > RESULT_LINE_CAP}>
          <text fg={theme.dim}>{`… (+${lines().length - RESULT_LINE_CAP} more lines)`}</text>
        </Show>
      </box>
    </Show>
  )
}

function OrphanResult(props: { cell: RenderCell; focused?: boolean; expanded?: boolean; demoted?: boolean; budget: number }) {
  return (
    <Row author={props.cell.author ?? "agent"} railColor={theme.border} focused={props.focused} demoted={props.demoted}>
      <ResultBody result={props.cell} expanded={props.expanded} demoted={props.demoted} budget={props.budget} />
    </Row>
  )
}

function UnknownCell(props: { cell: RenderCell; focused?: boolean; demoted?: boolean; budget: number; expanded?: boolean }) {
  const raw = () => props.cell.raw ?? props.cell.text ?? ""
  return (
    <Row author={props.cell.author ?? ""} railColor={theme.warning} focused={props.focused} demoted={props.demoted}>
      <text fg={theme.warning}>{`▸ ${props.cell.kind}`}</text>
      <text fg={theme.dim} wrapMode="word">
        {props.expanded ? raw() : truncate(raw(), props.budget * 3)}
      </text>
    </Row>
  )
}

export function Divider(props: { label: string }) {
  return (
    <box marginTop={1}>
      <text fg={theme.dim}>{`── ${props.label} ──`}</text>
    </box>
  )
}

export function DemotedBlockView(props: { reason: string; cells: RenderCell[] }) {
  const entries = createMemo(() => renderEntries(props.cells))
  return (
    <box flexDirection="column" marginTop={1}>
      <Divider label={`context ${props.reason} · ${props.cells.length} cell${props.cells.length === 1 ? "" : "s"} above ${props.cells.length === 1 ? "is" : "are"} no longer the agent's working memory`} />
      <For each={entries()}>{(entry) => <EntryView entry={entry} demoted />}</For>
    </box>
  )
}

// StreamingView renders the in-flight generation as a forming cell in the
// streaming register: dim, ▒ rail, trailing block cursor (DESIGN.md dev. #2).
export function StreamingView(props: { live: LiveState }) {
  const view = createMemo(() => liveView(props.live))
  const budget = useBudget()
  return (
    <box flexDirection="row" gap={1} marginTop={1}>
      <text width={AUTHOR_COL} flexShrink={0} fg={theme.dim}>
        {"▒"}
      </text>
      <box flexGrow={1} minWidth={0} border={["left"]} borderColor={tint(theme.accent, theme.bg, 0.5)} paddingLeft={1} flexDirection="column">
        <text fg={theme.dim}>{`▒ streaming · gen ${props.live.generation}`}</text>
        <Show when={view().reasoning}>
          <text fg={theme.dim} attributes={TextAttributes.ITALIC}>
            {`▸ thinking: ${tail(view().reasoning, budget())}`}
          </text>
        </Show>
        <Show when={view().text}>
          <text fg={tint(theme.text, theme.bg, 0.25)} wrapMode="word">
            {view().text + "▊"}
          </text>
        </Show>
        <For each={view().tools}>
          {(tool) => (
            <text fg={theme.dim}>{`▸ ${tool.name || tool.id} forming… ${tail(compactOneLine(tool.args), Math.max(12, budget() - 20))}`}</text>
          )}
        </For>
      </box>
    </box>
  )
}

function tail(s: string, n: number): string {
  const one = compactOneLine(s)
  return one.length > n ? "…" + one.slice(-(n - 1)) : one
}

function compactOneLine(s: string): string {
  return s.replace(/\s*\n\s*/g, " ⏎ ")
}

function prettyArgs(raw: string): string {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}
