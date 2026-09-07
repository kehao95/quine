// Char-frame capture harness: renders the full representative session for a
// human visual pass. Usage: bun dev/capture.tsx [W=100 H=44] [EXPAND=1]
import { testRender } from "@opentui/solid"
import { createStore, unwrap } from "solid-js/store"
import { App } from "../src/App"
import type { Command, CommandResponse } from "../src/contract"
import type { QcliState, QcliStore, RenderCell } from "../src/store"
import { initialState, reduceQcli } from "../src/store"

function cell(partial: Partial<RenderCell> & { seq: number; kind: string }): RenderCell {
  return {
    type: "cell",
    role: null,
    text: null,
    ts: null,
    tool_id: null,
    tool_name: null,
    status: null,
    is_error: false,
    action: null,
    author: null,
    raw: null,
    ...partial,
  }
}

const [state, setState] = createStore<QcliState>({
  ...initialState(),
  peer: { session: "sess_demo_aaaa1111", run_id: "run_demo", pid: 48213, agent_root: "/tmp/a", runtime_root: "/tmp/r" },
  statusSnapshot: {
    type: "status",
    live: true,
    working: true,
    pid: 48213,
    ppid: 1,
    session: "sess_demo_aaaa1111",
    run_id: "run_demo",
    incarnation: 2,
    parent_session: "",
    model: "gpt-5.5",
    depth: 2,
    pending: 0,
    tokens_in: null,
    tokens_out: null,
  },
  connection: "connected",
  backfilled: true,
  demoted: [
    {
      reason: "compaction",
      cells: [
        cell({ seq: 1, kind: "message", role: "assistant", text: "earlier cognition that the agent has since compacted away", register: "demoted" }),
      ],
    },
  ],
  cells: [
    {
      ...cell({ seq: 10, kind: "human_post", author: "human", action: "inject", text: "look at the empty-package bug in the novelty gate" }),
      receipt: { queued: true, received: true, woke: true, delivered: true },
    },
    cell({ seq: 11, kind: "reasoning", text: "the novelty gate seems to skip empty packages\nbecause the scorer counts decls\nand empty pkgs have none" }),
    cell({
      seq: 12,
      kind: "message",
      role: "assistant",
      text: "The **novelty gate** misses empty packages. Minimal check:\n\n```go\nif len(pkg.Decls) == 0 {\n\treturn false\n}\n```",
    }),
    cell({ seq: 13, kind: "tool_call", tool_id: "call_1", tool_name: "bash", text: `{"command":"awk -f scorer.awk gen-*/decls | sort | uniq -c"}` }),
    cell({ seq: 14, kind: "tool_result", tool_id: "call_1", tool_name: "bash", status: "ok", text: "3/5 packages flagged empty\ngen-2/pkg: 0 decls\ngen-4/pkg: 0 decls" }),
    cell({ seq: 15, kind: "reply", author: "sess_demo_aaaa1111", text: "pasted the diff to peer-3 for review" }),
  ],
  live: {
    generation: 3,
    parts: [
      { type: "stream_delta", generation: 3, seq: 1, kind: "reasoning_delta", text: "now checking the scorer output for gen-5", tool_id: null, tool_name: null, ts: 1 },
      { type: "stream_delta", generation: 3, seq: 2, kind: "text_delta", text: "Checking the score for the remaining packages no", tool_id: null, tool_name: null, ts: 2 },
    ],
  },
  roster: [
    { pid: 48213, session: "sess_demo_aaaa1111", run_id: "run_demo", parent_session: "", depth: 2, model: "gpt-5.5", live: true, pending: 0, agent_root: "/tmp/a", attached: true },
    { pid: 48300, session: "sess_demo_bbbb2222", run_id: "run_demo2", parent_session: "sess_demo_aaaa1111", depth: 3, model: "gpt-5.5", live: false, pending: 1, agent_root: "/tmp/b", attached: false },
  ],
})

const store: QcliStore = {
  state,
  connect: async () => {},
  command: async (cmd: Command): Promise<CommandResponse> => ({ type: "queued", action: cmd.op === "send" ? cmd.action : "inject", client_ref: "c-1" }),
  fetchRoster: async () => state.roster,
  fetchPeerContract: async () => null,
  dispatch: (event) => setState(reduceQcli(unwrap(state), event)),
  clearUnread: () => setState("unread", 0),
}

const width = Number(process.env.W ?? 100)
const height = Number(process.env.H ?? 44)
const app = await testRender(() => <App store={store} />, { width, height })

async function pump(cond: (f: string) => boolean, passes = 80) {
  for (let i = 0; i < passes; i++) {
    await app.renderOnce()
    if (cond(app.captureCharFrame())) return
    await new Promise((resolve) => setTimeout(resolve, 15))
  }
}

await pump((f) => f.includes("Minimal check:") && f.includes("pasted the diff") && f.includes("[▷✓⚡▣]") && !f.includes("**"), 200)

if (process.env.EXPAND === "1") {
  // Focus up to the reasoning cell and expand it, then expand the tool pair.
  app.mockInput.pressArrow("up", { ctrl: true }) // reply
  app.mockInput.pressArrow("up") // pair
  app.mockInput.pressArrow("right")
  app.mockInput.pressArrow("up") // message
  app.mockInput.pressArrow("up") // reasoning
  app.mockInput.pressArrow("right")
  await pump((f) => f.includes("▾ reasoning") && f.includes("because the scorer counts decls"))
}

process.stdout.write(app.captureCharFrame() + "\n")
app.renderer.destroy()
