export const QCLI_CONTRACT = "qcli/1" as const

export type ControlAction = "post" | "poke" | "inject" | "interrupt"

export interface PeerInfo {
  session: string
  run_id: string
  pid: number
  agent_root: string
  runtime_root: string
}

export interface HelloEvent {
  type: "hello"
  contract: typeof QCLI_CONTRACT
  endpoint_id: string
  peer: PeerInfo
}

export interface StatusEvent {
  type: "status"
  live: boolean
  working: boolean
  pid: number
  ppid: number
  session: string
  run_id: string
  incarnation: number
  parent_session: string
  model: string
  depth: number
  pending: number
  tokens_in: number | null
  tokens_out: number | null
}

export type CellKind =
  | "meta"
  | "message"
  | "reasoning"
  | "tool_call"
  | "tool_result"
  | "human_post"
  | "reply"
  | "outcome"
  | string

export interface Cell {
  type?: "cell"
  seq: number
  kind: CellKind
  role: string | null
  text: string | null
  ts: number | null
  tool_id: string | null
  tool_name: string | null
  status: string | null
  is_error: boolean
  action: ControlAction | string | null
  author: string | null
  raw: string | null
  depth?: number | null
  session?: string | null
  termination_mode?: string | null
  exit_code?: number | null
  turn_count?: number | null
  tokens_in?: number | null
  tokens_out?: number | null
}

export interface BackfillCompleteEvent {
  type: "backfill_complete"
  count: number
}

export type StreamKind = "text_delta" | "reasoning_delta" | "tool_call_delta" | "completed" | string

export interface StreamDeltaEvent {
  type: "stream_delta"
  generation: number
  seq: number
  kind: StreamKind
  text: string | null
  tool_id: string | null
  tool_name: string | null
  ts: number
}

export type ReceiptStage = "queued" | "received" | "woke" | "delivered"

export interface ReceiptEvent {
  type: "receipt"
  stage: Exclude<ReceiptStage, "queued">
  action: ControlAction | string | null
  delivery: string | null
  message_id: string | null
  client_ref: string | null
  pending: number | null
  ts: number
}

export interface ContextResetEvent {
  type: "context_reset"
  reason: "compaction" | "incarnation" | string
  ts: number
}

export interface ErrorEvent {
  type: "error"
  scope: "send" | "attach" | "spawn" | "tail" | "server" | string
  code: string
  message: string
  recoverable: boolean
}

export type QcliEvent =
  | HelloEvent
  | StatusEvent
  | Cell
  | BackfillCompleteEvent
  | StreamDeltaEvent
  | ReceiptEvent
  | ContextResetEvent
  | ErrorEvent

export interface RosterEntry {
  pid: number
  session: string
  run_id: string
  parent_session: string
  depth: number
  model: string
  live: boolean
  pending: number
  agent_root: string
  attached: boolean
}

export interface RosterResponse {
  type: "roster"
  peers: RosterEntry[]
}

export interface QueuedResponse {
  type: "queued"
  action: ControlAction
  client_ref: string
}

export interface AttachedResponse {
  type: "attached"
  peer: PeerInfo
}

export interface PeerContractResponse {
  type: "peer_contract"
  contract: {
    control_actions?: Record<string, { description?: string }>
    [key: string]: unknown
  }
}

export type CommandResponse = QueuedResponse | AttachedResponse | ErrorEvent

export interface AttachCommand {
  op: "attach"
  selector: {
    agent_root?: string
    pid?: number
    session?: string
    runtime_root?: string
  }
}

export interface SpawnCommand {
  op: "spawn"
  mission: string
  runtime_root?: string
}

export interface SendCommand {
  op: "send"
  action: ControlAction
  text: string
}

export type Command = AttachCommand | SpawnCommand | SendCommand

export function isQcliEvent(value: unknown): value is QcliEvent {
  return !!value && typeof value === "object" && typeof (value as { type?: unknown }).type === "string"
}

export function shortSession(session: string | undefined | null): string {
  if (!session) return "peer"
  return session.length > 10 ? session.slice(-10) : session
}

export function cellText(cell: Pick<Cell, "text" | "raw">): string {
  return cell.text ?? cell.raw ?? ""
}
