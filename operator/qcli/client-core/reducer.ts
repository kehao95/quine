// Framework-free qcli/1 data layer shared by every client (terminal TUI, web).
// Nothing in this file may import OpenTUI, solid-js, or DOM APIs — the reducer
// table, pairing, and ladder models are pure functions over contract events
// (WEB.md §3: the client-core extraction). The TUI suite is the acceptance
// test of this move.
import type {
  Cell,
  Command,
  CommandResponse,
  ControlAction,
  ErrorEvent,
  PeerContractResponse,
  PeerInfo,
  QcliEvent,
  ReceiptEvent,
  ReceiptStage,
  RosterEntry,
  StatusEvent,
  StreamDeltaEvent,
} from "./contract"

export type ConnectionState = "connecting" | "connected" | "reconnecting" | "rebuilding" | "disconnected"

export interface ReceiptState {
  queued?: boolean
  received?: boolean
  woke?: boolean
  delivered?: boolean
  // Stage timestamps (ms) — receipt events carry ts; queued is stamped by the
  // client ack. The web inspector renders these as the delivery timeline.
  queued_ts?: number
  received_ts?: number
  woke_ts?: number
  delivered_ts?: number
  message_id?: string
  error?: string
}

export interface RenderCell extends Cell {
  local?: boolean
  client_ref?: string
  receipt?: ReceiptState
  register?: "settled" | "local" | "demoted" | "interrupted"
  error_code?: string
}

// RenderEntry is the derived transcript unit: tool_call/tool_result cells are
// paired by tool_id into one entry; everything else renders standalone.
export type RenderEntry =
  | { key: string; kind: "cell"; cell: RenderCell }
  | { key: string; kind: "pair"; call: RenderCell; result: RenderCell | null }

export function entryKey(cell: RenderCell): string {
  return `${cell.seq}:${cell.kind}:${cell.tool_id ?? ""}`
}

export function renderEntries(cells: RenderCell[]): RenderEntry[] {
  const used = new Set<number>()
  const out: RenderEntry[] = []
  cells.forEach((cell, i) => {
    if (used.has(i)) return
    if (cell.kind === "tool_call") {
      let result: RenderCell | null = null
      if (cell.tool_id) {
        for (let j = i + 1; j < cells.length; j++) {
          const candidate = cells[j]!
          if (!used.has(j) && candidate.kind === "tool_result" && candidate.tool_id === cell.tool_id) {
            result = candidate
            used.add(j)
            break
          }
        }
      }
      out.push({ key: entryKey(cell), kind: "pair", call: cell, result })
      return
    }
    out.push({ key: entryKey(cell), kind: "cell", cell })
  })
  return out
}

export interface LiveState {
  generation: number
  parts: StreamDeltaEvent[]
  interrupted?: boolean
}

export interface DemotedBlock {
  reason: string
  cells: RenderCell[]
}

export interface QcliState {
  endpointId: string
  peer: PeerInfo | null
  statusSnapshot: StatusEvent | null
  connection: ConnectionState
  cells: RenderCell[]
  demoted: DemotedBlock[]
  live: LiveState | null
  receipts: Record<string, ReceiptState>
  messageRefs: Record<string, string>
  unread: number
  backfilled: boolean
  lastError: ErrorEvent | null
  laneNotice: string
  roster: RosterEntry[]
  peerContract: PeerContractResponse["contract"] | null
}

export interface QcliStore {
  state: QcliState
  connect: (signal: AbortSignal) => Promise<void>
  command: (cmd: Command) => Promise<CommandResponse>
  fetchRoster: () => Promise<RosterEntry[]>
  fetchPeerContract: () => Promise<PeerContractResponse["contract"] | null>
  dispatch: (event: QcliEvent) => void
  clearUnread: () => void
}

export function initialState(): QcliState {
  return {
    endpointId: "",
    peer: null,
    statusSnapshot: null,
    connection: "connecting",
    cells: [],
    demoted: [],
    live: null,
    receipts: {},
    messageRefs: {},
    unread: 0,
    backfilled: false,
    lastError: null,
    laneNotice: "",
    roster: [],
    peerContract: null,
  }
}

export function reduceQcli(state: QcliState, event: QcliEvent): QcliState {
  switch (event.type) {
    case "hello":
      return {
        ...initialState(),
        endpointId: event.endpoint_id,
        peer: event.peer,
        connection: "connected",
      }
    case "status":
      return { ...state, statusSnapshot: event, connection: state.connection === "rebuilding" ? "connected" : state.connection }
    case "cell":
      return reduceCell(state, event as Cell)
    case "backfill_complete":
      return { ...state, backfilled: true, laneNotice: `backfill complete (${event.count})` }
    case "stream_delta":
      return reduceStream(state, event)
    case "receipt":
      return reduceReceipt(state, event)
    case "context_reset":
      return {
        ...state,
        demoted: state.cells.length > 0 ? [...state.demoted, { reason: event.reason, cells: state.cells.map((c) => ({ ...c, register: "demoted" as const })) }] : state.demoted,
        cells: [],
        live: null,
        backfilled: false,
        laneNotice: `context reset: ${event.reason}`,
      }
    case "error":
      return reduceError(state, event)
    default:
      return state
  }
}

export function applyLocalQueued(state: QcliState, action: ControlAction, text: string, clientRef: string): QcliState {
  if (!text.trim() && (action === "poke" || action === "interrupt")) {
    return {
      ...state,
      receipts: { ...state.receipts, [clientRef]: { queued: true, queued_ts: Date.now() } },
      laneNotice: `${action} queued`,
    }
  }
  const receipt: ReceiptState = { queued: true, queued_ts: Date.now() }
  const local: RenderCell = {
    type: "cell",
    seq: nextLocalSeq(state),
    kind: "human_post",
    role: null,
    text,
    ts: Date.now(),
    tool_id: null,
    tool_name: null,
    status: null,
    is_error: false,
    action,
    author: "human",
    raw: null,
    local: true,
    client_ref: clientRef,
    receipt,
    register: "local",
  }
  return {
    ...state,
    cells: [...state.cells, local],
    receipts: { ...state.receipts, [clientRef]: receipt },
    laneNotice: `${action} queued`,
  }
}

function reduceCell(state: QcliState, cell: Cell): QcliState {
  const render = cell as RenderCell
  if (render.kind === "meta") {
    return state
  }
  const reconciled = reconcileLocalEcho(state.cells, render)
  const cells = reconciled.replaced ? reconciled.cells : [...state.cells, render]
  const unread = state.backfilled && (render.kind === "reply" || render.kind === "outcome") ? state.unread + 1 : state.unread
  return { ...state, cells, unread }
}

function reconcileLocalEcho(cells: RenderCell[], authoritative: RenderCell): { cells: RenderCell[]; replaced: boolean } {
  if (authoritative.kind !== "human_post") return { cells, replaced: false }
  const idx = cells.findIndex(
    (c) => c.local && c.kind === "human_post" && c.text === authoritative.text && c.action === authoritative.action,
  )
  if (idx < 0) return { cells, replaced: false }
  const next = cells.slice()
  next[idx] = { ...authoritative, receipt: next[idx]!.receipt, client_ref: next[idx]!.client_ref }
  return { cells: next, replaced: true }
}

function reduceStream(state: QcliState, event: StreamDeltaEvent): QcliState {
  if (event.kind === "completed") {
    if (!state.live || state.live.generation !== event.generation) return state
    return { ...state, live: null, laneNotice: "generation completed" }
  }
  let live = state.live
  let cells = state.cells
  if (!live || live.generation !== event.generation) {
    // A new generation opening while the previous never saw `completed` means
    // the previous one was interrupted/errored: flip it into an interrupted
    // cell in place (TUI.md §4) — clearly not-persisted, never synthesized as
    // a settled cell.
    if (live && live.parts.length > 0) {
      cells = [...cells, interruptedCellFromLive(live, nextLocalSeq(state))]
    }
    live = { generation: event.generation, parts: [] }
  }
  return { ...state, cells, live: { ...live, parts: [...live.parts, event] } }
}

function interruptedCellFromLive(live: LiveState, seq: number): RenderCell {
  const view = liveView(live)
  const toolNote = view.tools.length > 0 ? view.tools.map((t) => `${t.name || t.id} ${t.args}`).join("\n") : ""
  const text = [view.reasoning && `thinking: ${view.reasoning}`, view.text, toolNote].filter(Boolean).join("\n")
  return {
    type: "cell",
    seq,
    kind: "message",
    role: "assistant",
    text,
    ts: null,
    tool_id: null,
    tool_name: null,
    status: null,
    is_error: false,
    action: null,
    author: null,
    raw: null,
    register: "interrupted",
  }
}

function reduceReceipt(state: QcliState, event: ReceiptEvent): QcliState {
  if (!event.client_ref) {
    if (event.stage === "woke") {
      return { ...state, laneNotice: `peer woke${event.pending != null ? `; pending=${event.pending}` : ""}` }
    }
    return state
  }
  const prev = state.receipts[event.client_ref] ?? {}
  const nextReceipt = advanceReceipt(prev, event.stage, event.ts)
  if (event.message_id) nextReceipt.message_id = event.message_id
  const receipts = { ...state.receipts, [event.client_ref]: nextReceipt }
  const messageRefs = event.message_id ? { ...state.messageRefs, [event.message_id]: event.client_ref } : state.messageRefs
  const cells = state.cells.map((c) => (c.client_ref === event.client_ref ? { ...c, receipt: nextReceipt } : c))
  return { ...state, receipts, messageRefs, cells, laneNotice: `${event.stage}: ${event.client_ref}` }
}

function reduceError(state: QcliState, event: ErrorEvent): QcliState {
  if (event.scope !== "send") {
    return { ...state, lastError: event, laneNotice: `${event.scope}: ${event.code}` }
  }
  const cells = state.cells.slice()
  for (let i = cells.length - 1; i >= 0; i--) {
    const c = cells[i]
    if (c?.local && !c.error_code) {
      cells[i] = { ...c, error_code: event.code, receipt: { ...(c.receipt ?? {}), error: event.code } }
      break
    }
  }
  return { ...state, cells, lastError: event, laneNotice: `send error: ${event.code}` }
}

function advanceReceipt(prev: ReceiptState, stage: ReceiptStage, ts?: number): ReceiptState {
  const next = { ...prev }
  if (stage === "queued") {
    next.queued = true
    next.queued_ts = next.queued_ts ?? ts
  }
  if (stage === "received") {
    next.received = true
    next.received_ts = next.received_ts ?? ts
  }
  if (stage === "woke") {
    next.woke = true
    next.woke_ts = next.woke_ts ?? ts
  }
  if (stage === "delivered") {
    next.delivered = true
    next.delivered_ts = next.delivered_ts ?? ts
  }
  return next
}

function nextLocalSeq(state: QcliState): number {
  const last = state.cells.reduce((max, c) => Math.max(max, c.seq ?? 0), 0)
  return last + 0.001
}

export function actionDescription(state: QcliState, action: ControlAction): string {
  const desc = state.peerContract?.control_actions?.[action]?.description
  return typeof desc === "string" && desc.trim() ? desc : action
}

// liveView groups in-flight deltas into the sub-elements of the forming cell:
// reasoning summary, message text, and forming tool calls (by tool_id).
export interface LiveView {
  reasoning: string
  text: string
  tools: Array<{ id: string; name: string; args: string }>
}

export function liveView(live: LiveState | null): LiveView {
  const view: LiveView = { reasoning: "", text: "", tools: [] }
  if (!live) return view
  const toolIndex = new Map<string, number>()
  for (const part of live.parts) {
    if (part.kind === "text_delta") view.text += part.text ?? ""
    else if (part.kind === "reasoning_delta") view.reasoning += part.text ?? ""
    else if (part.kind === "tool_call_delta") {
      const id = part.tool_id ?? `tool-${toolIndex.size}`
      let idx = toolIndex.get(id)
      if (idx === undefined) {
        idx = view.tools.length
        toolIndex.set(id, idx)
        view.tools.push({ id, name: part.tool_name ?? "", args: "" })
      }
      if (part.tool_name) view.tools[idx]!.name = part.tool_name
      view.tools[idx]!.args += part.text ?? ""
    }
  }
  return view
}

export function liveText(live: LiveState | null): string {
  return liveView(live).text
}

// receiptSegments is the ladder model: ▷ queued → ✓ received → ⚡ woke →
// ▣ delivered, lit exactly per receipt state (woke is display-level, see
// TUI.md §5). Unlit stages render as `·` so char-frame tests can assert them.
export interface ReceiptSegment {
  glyph: string
  label: string
  lit: boolean
  ts?: number
}

export function receiptSegments(receipt: ReceiptState | undefined): ReceiptSegment[] {
  return [
    { glyph: "▷", label: "queued", lit: !!receipt?.queued, ts: receipt?.queued_ts },
    { glyph: "✓", label: "received", lit: !!receipt?.received, ts: receipt?.received_ts },
    { glyph: "⚡", label: "woke", lit: !!receipt?.woke, ts: receipt?.woke_ts },
    { glyph: "▣", label: "delivered", lit: !!receipt?.delivered, ts: receipt?.delivered_ts },
  ]
}

export function receiptLadder(receipt: ReceiptState | undefined): string {
  return "[" + receiptSegments(receipt).map((s) => (s.lit ? s.glyph : "·")).join("") + "]"
}

export function isControlAction(v: string): v is ControlAction {
  return v === "post" || v === "poke" || v === "inject" || v === "interrupt"
}

export function parseSlashCommand(
  text: string,
): { kind: "none" } | { kind: "action"; action: ControlAction } | { kind: "spawn"; mission: string } | { kind: "help" } | { kind: "quit" } | { kind: "switch" } {
  const trimmed = text.trim()
  if (!trimmed.startsWith("/")) return { kind: "none" }
  const [cmd, ...rest] = trimmed.slice(1).split(/\s+/)
  switch (cmd) {
    case "action": {
      const action = rest[0] ?? ""
      return isControlAction(action) ? { kind: "action", action } : { kind: "none" }
    }
    case "spawn":
      return { kind: "spawn", mission: rest.join(" ") }
    case "switch":
      return { kind: "switch" }
    case "help":
    case "?":
      return { kind: "help" }
    case "quit":
      return { kind: "quit" }
    default:
      return { kind: "none" }
  }
}
