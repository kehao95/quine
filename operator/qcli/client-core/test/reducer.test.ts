import { expect, test } from "bun:test"
import type { Cell, QcliEvent } from "../contract"
import {
  actionDescription,
  applyLocalQueued,
  initialState,
  liveView,
  parseSlashCommand,
  receiptLadder,
  reduceQcli,
  renderEntries,
} from "../reducer"

const baseCell: Cell = {
  type: "cell",
  seq: 1,
  kind: "message",
  role: "assistant",
  text: "hello",
  ts: 1,
  tool_id: null,
  tool_name: null,
  status: null,
  is_error: false,
  action: null,
  author: null,
  raw: null,
}

// --- reducer table (TUI.md §3.3) --------------------------------------------

test("hello resets volatile state", () => {
  let state = initialState()
  state = reduceQcli(state, { ...baseCell, kind: "reply", text: "old" })
  state = reduceQcli(state, {
    type: "hello",
    contract: "qcli/1",
    endpoint_id: "qcli-1",
    peer: { session: "sess", run_id: "run", pid: 1, agent_root: "/a", runtime_root: "/r" },
  })
  expect(state.cells).toHaveLength(0)
  expect(state.peer?.session).toBe("sess")
  expect(state.endpointId).toBe("qcli-1")
})

test("status patches the snapshot", () => {
  let state = initialState()
  state = reduceQcli(state, {
    type: "status",
    live: true,
    working: true,
    pid: 9,
    ppid: 1,
    session: "s",
    run_id: "r",
    incarnation: 1,
    parent_session: "",
    model: "m",
    depth: 0,
    pending: 3,
    tokens_in: null,
    tokens_out: null,
  })
  expect(state.statusSnapshot?.pending).toBe(3)
})

test("meta cells feed the header, not the stream", () => {
  let state = initialState()
  state = reduceQcli(state, { ...baseCell, kind: "meta" })
  expect(state.cells).toHaveLength(0)
})

test("backfill gates unread notifications", () => {
  let state = initialState()
  state = reduceQcli(state, { ...baseCell, kind: "reply", text: "replay" })
  expect(state.unread).toBe(0)
  state = reduceQcli(state, { type: "backfill_complete", count: 1 })
  state = reduceQcli(state, { ...baseCell, seq: 2, kind: "reply", text: "live" })
  state = reduceQcli(state, { ...baseCell, seq: 3, kind: "outcome" })
  state = reduceQcli(state, { ...baseCell, seq: 4, kind: "message" })
  expect(state.unread).toBe(2)
})

test("stream reducer opens and closes live register", () => {
  let state = initialState()
  state = reduceQcli(state, { type: "stream_delta", generation: 3, seq: 1, kind: "text_delta", text: "hel", tool_id: null, tool_name: null, ts: 1 })
  state = reduceQcli(state, { type: "stream_delta", generation: 3, seq: 2, kind: "text_delta", text: "lo", tool_id: null, tool_name: null, ts: 2 })
  expect(state.live?.parts).toHaveLength(2)
  state = reduceQcli(state, { type: "stream_delta", generation: 3, seq: 3, kind: "completed", text: null, tool_id: null, tool_name: null, ts: 3 })
  expect(state.live).toBeNull()
})

test("generation rollover without completed synthesizes an interrupted cell", () => {
  let state = initialState()
  state = reduceQcli(state, { type: "stream_delta", generation: 1, seq: 1, kind: "reasoning_delta", text: "think", tool_id: null, tool_name: null, ts: 1 })
  state = reduceQcli(state, { type: "stream_delta", generation: 1, seq: 2, kind: "text_delta", text: "partial", tool_id: null, tool_name: null, ts: 2 })
  state = reduceQcli(state, { type: "stream_delta", generation: 2, seq: 1, kind: "text_delta", text: "restart", tool_id: null, tool_name: null, ts: 3 })
  expect(state.cells).toHaveLength(1)
  expect(state.cells[0]!.register).toBe("interrupted")
  expect(state.cells[0]!.text).toContain("partial")
  expect(state.cells[0]!.text).toContain("think")
  expect(state.live?.generation).toBe(2)
})

test("completed for a stale generation is ignored", () => {
  let state = initialState()
  state = reduceQcli(state, { type: "stream_delta", generation: 2, seq: 1, kind: "text_delta", text: "x", tool_id: null, tool_name: null, ts: 1 })
  state = reduceQcli(state, { type: "stream_delta", generation: 1, seq: 9, kind: "completed", text: null, tool_id: null, tool_name: null, ts: 2 })
  expect(state.live?.generation).toBe(2)
})

test("context reset demotes existing cells and drops backfilled", () => {
  let state = initialState()
  state = reduceQcli(state, baseCell)
  state = reduceQcli(state, { type: "backfill_complete", count: 1 })
  state = reduceQcli(state, { type: "context_reset", reason: "compaction", ts: 2 })
  expect(state.cells).toHaveLength(0)
  expect(state.demoted).toHaveLength(1)
  expect(state.demoted[0]?.cells[0]?.register).toBe("demoted")
  expect(state.backfilled).toBe(false)
})

test("local queued cell advances through receipt ladder without regression", () => {
  let state = applyLocalQueued(initialState(), "inject", "hello", "c-1")
  expect(state.cells[0]?.local).toBe(true)
  expect(receiptLadder(state.cells[0]?.receipt)).toBe("[▷···]")
  const events: QcliEvent[] = [
    { type: "receipt", stage: "received", action: "inject", delivery: null, message_id: "message-1", client_ref: "c-1", pending: null, ts: 1 },
    { type: "receipt", stage: "woke", action: null, delivery: "inject", message_id: null, client_ref: "c-1", pending: 1, ts: 2 },
    { type: "receipt", stage: "delivered", action: null, delivery: "inject", message_id: "message-1", client_ref: "c-1", pending: null, ts: 3 },
  ]
  for (const event of events) state = reduceQcli(state, event)
  expect(receiptLadder(state.cells[0]?.receipt)).toBe("[▷✓⚡▣]")
  expect(state.messageRefs["message-1"]).toBe("c-1")
})

test("woke receipt without client_ref lands in the lane only", () => {
  let state = initialState()
  state = reduceQcli(state, { type: "receipt", stage: "woke", action: null, delivery: "poke", message_id: null, client_ref: null, pending: 2, ts: 1 })
  expect(state.laneNotice).toContain("woke")
  expect(state.cells).toHaveLength(0)
})

test("authoritative human_post reconciles local echo by text and action", () => {
  let state = applyLocalQueued(initialState(), "inject", "hello", "c-1")
  state = reduceQcli(state, { ...baseCell, seq: 2, kind: "human_post", action: "inject", author: "human", text: "hello" })
  expect(state.cells).toHaveLength(1)
  expect(state.cells[0]?.local).toBeUndefined()
  expect(state.cells[0]?.client_ref).toBe("c-1")
})

test("send error attaches to the pending local echo", () => {
  let state = applyLocalQueued(initialState(), "inject", "doomed", "c-1")
  state = reduceQcli(state, { type: "error", scope: "send", code: "write_timeout", message: "ctl blocked", recoverable: true })
  expect(state.cells[0]?.error_code).toBe("write_timeout")
  expect(state.laneNotice).toContain("write_timeout")
})

test("non-send error lands in the lane", () => {
  let state = initialState()
  state = reduceQcli(state, { type: "error", scope: "tail", code: "peer_dead", message: "gone", recoverable: false })
  expect(state.lastError?.code).toBe("peer_dead")
  expect(state.cells).toHaveLength(0)
})

// --- derived models -----------------------------------------------------------

test("renderEntries pairs tool_call with its result by tool_id", () => {
  const call: Cell = { ...baseCell, seq: 2, kind: "tool_call", tool_id: "c1", tool_name: "bash", text: "{}" }
  const result: Cell = { ...baseCell, seq: 3, kind: "tool_result", tool_id: "c1", tool_name: "bash", status: "ok", text: "out" }
  const other: Cell = { ...baseCell, seq: 4 }
  const entries = renderEntries([call, result, other])
  expect(entries).toHaveLength(2)
  expect(entries[0]!.kind).toBe("pair")
  const pair = entries[0] as Extract<(typeof entries)[0], { kind: "pair" }>
  expect(pair.result?.text).toBe("out")
})

test("renderEntries leaves orphans standalone", () => {
  const call: Cell = { ...baseCell, seq: 2, kind: "tool_call", tool_id: "c1", tool_name: "bash", text: "{}" }
  const strayResult: Cell = { ...baseCell, seq: 3, kind: "tool_result", tool_id: "cX", tool_name: "bash", status: "ok", text: "stray" }
  const entries = renderEntries([call, strayResult])
  expect(entries).toHaveLength(2)
  expect(entries[0]!.kind).toBe("pair")
  expect((entries[0] as Extract<(typeof entries)[0], { kind: "pair" }>).result).toBeNull()
  expect(entries[1]!.kind).toBe("cell")
})

test("liveView groups deltas into forming sub-elements", () => {
  const view = liveView({
    generation: 1,
    parts: [
      { type: "stream_delta", generation: 1, seq: 1, kind: "reasoning_delta", text: "why ", tool_id: null, tool_name: null, ts: 1 },
      { type: "stream_delta", generation: 1, seq: 2, kind: "text_delta", text: "answer", tool_id: null, tool_name: null, ts: 2 },
      { type: "stream_delta", generation: 1, seq: 3, kind: "tool_call_delta", text: `{"a`, tool_id: "t1", tool_name: "bash", ts: 3 },
      { type: "stream_delta", generation: 1, seq: 4, kind: "tool_call_delta", text: `":1}`, tool_id: "t1", tool_name: null, ts: 4 },
    ],
  })
  expect(view.reasoning).toBe("why ")
  expect(view.text).toBe("answer")
  expect(view.tools).toHaveLength(1)
  expect(view.tools[0]!.args).toBe(`{"a":1}`)
  expect(view.tools[0]!.name).toBe("bash")
})

test("action description comes from raw peer contract when present", () => {
  const state = {
    ...initialState(),
    peerContract: {
      control_actions: {
        inject: { description: "queue-and-deliver from peer contract" },
      },
    },
  }
  expect(actionDescription(state, "inject")).toBe("queue-and-deliver from peer contract")
  expect(actionDescription(state, "poke")).toBe("poke")
})

test("slash commands parse", () => {
  expect(parseSlashCommand("/switch")).toEqual({ kind: "switch" })
  expect(parseSlashCommand("/action interrupt")).toEqual({ kind: "action", action: "interrupt" })
  expect(parseSlashCommand("/spawn build a nest")).toEqual({ kind: "spawn", mission: "build a nest" })
  expect(parseSlashCommand("plain text")).toEqual({ kind: "none" })
})
