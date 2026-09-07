import { expect, test } from "bun:test"
import { testRender } from "@opentui/solid"
import type { TestRendererSetup } from "@opentui/core/testing"
import type { JSX } from "solid-js"
import { App } from "../src/App"
import { EntryView, StreamingView } from "../src/components/CellView"
import type { RenderCell } from "../src/store"
import { initialState, reduceQcli, renderEntries } from "../src/store"
import { fixtureStore, makeCell } from "./fixtures"

async function frameOf(node: () => JSX.Element, opts: { width?: number; height?: number; ready?: (f: string) => boolean } = {}): Promise<string> {
  const app: TestRendererSetup = await testRender(node, { width: opts.width ?? 100, height: opts.height ?? 30 })
  let out = ""
  const ready = opts.ready ?? ((f: string) => f.trim().length > 0)
  // Markdown parsing is async (cold parser worker): poll frames until the
  // content predicate holds instead of sleeping a fixed number of passes.
  for (let i = 0; i < 60; i++) {
    await app.renderOnce()
    out = app.captureCharFrame()
    if (ready(out)) break
    await new Promise((resolve) => setTimeout(resolve, 15))
  }
  app.renderer.destroy()
  return out
}

function entryOf(cell: RenderCell) {
  return renderEntries([cell])[0]!
}

// --- settled cell kinds -----------------------------------------------------

test("message renders markdown with code fence", async () => {
  const cell = makeCell({
    seq: 1,
    kind: "message",
    role: "assistant",
    text: "The **novelty gate** missed empty packages:\n\n```go\nif len(pkg.Decls) == 0 {\n\treturn false\n}\n```",
  })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />, {
    ready: (frame) => frame.includes("novelty gate") && frame.includes("len(pkg.Decls)") && !frame.includes("**"),
  })
  expect(f).toContain("novelty gate")
  expect(f).not.toContain("**") // markdown bold markers concealed, not raw
  expect(f).toContain("len(pkg.Decls)")
  expect(f).toMatchSnapshot()
})

test("human_post renders action tag, ladder and body", async () => {
  const cell = makeCell({
    seq: 1,
    kind: "human_post",
    author: "human",
    action: "inject",
    text: "look at the empty-package bug",
  }) as RenderCell
  cell.receipt = { queued: true, received: true }
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />)
  expect(f).toContain("◀ inject")
  expect(f).toContain("[▷✓··]")
  expect(f).toContain("empty-package")
  expect(f).toMatchSnapshot()
})

test("receipt ladder renders each stage distinctly", async () => {
  const stages: Array<[string, RenderCell["receipt"]]> = [
    ["queued-only", { queued: true }],
    ["received", { queued: true, received: true }],
    ["woke", { queued: true, received: true, woke: true }],
    ["delivered", { queued: true, received: true, woke: true, delivered: true }],
  ]
  const cells = stages.map(([label, receipt], i) => {
    const cell = makeCell({ seq: i + 1, kind: "human_post", author: "human", action: "post", text: label }) as RenderCell
    cell.receipt = receipt
    return cell
  })
  const f = await frameOf(() => (
    <box flexDirection="column">
      {cells.map((cell) => (
        <EntryView entry={entryOf(cell)} />
      ))}
    </box>
  ))
  expect(f).toContain("[▷···]")
  expect(f).toContain("[▷✓··]")
  expect(f).toContain("[▷✓⚡·]")
  expect(f).toContain("[▷✓⚡▣]")
  expect(f).toMatchSnapshot()
})

test("reply renders with distinct accent rail", async () => {
  const cell = makeCell({ seq: 1, kind: "reply", author: "sess_peer_3", text: "pasted diff to peer-3" })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />, { ready: (frame) => frame.includes("pasted diff") })
  expect(f).toContain("◀ reply")
  expect(f).toContain("pasted diff")
  expect(f).toMatchSnapshot()
})

test("reasoning collapses to summary with expand affordance", async () => {
  const cell = makeCell({ seq: 1, kind: "reasoning", text: "the novelty gate seems to skip empty pkgs\nsecond line\nthird line" })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />)
  expect(f).toContain("▸ reasoning · 3 lines")
  expect(f).toContain("(→ expand)")
  expect(f).not.toContain("second line")
  expect(f).toMatchSnapshot()
})

test("reasoning expanded shows the full text", async () => {
  const cell = makeCell({ seq: 1, kind: "reasoning", text: "first line\nsecond line\nthird line" })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} expanded />)
  expect(f).toContain("▾ reasoning")
  expect(f).toContain("second line")
  expect(f).toContain("(← collapse)")
  expect(f).toMatchSnapshot()
})

test("outcome renders the termination line", async () => {
  const cell = makeCell({ seq: 1, kind: "outcome", termination_mode: "exit", exit_code: 0, turn_count: 12, tokens_in: 900, tokens_out: 210 })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />)
  expect(f).toContain("■ outcome · exit · exit=0 · turns=12 · tokens=900/210")
  expect(f).toMatchSnapshot()
})

test("unknown kind renders raw block", async () => {
  const cell = makeCell({ seq: 1, kind: "novel_kind", raw: `{"kind":"novel_kind","payload":42}` })
  const f = await frameOf(() => <EntryView entry={entryOf(cell)} />)
  expect(f).toContain("▸ novel_kind")
  expect(f).toContain(`"payload":42`)
  expect(f).toMatchSnapshot()
})

// --- tool pairing -----------------------------------------------------------

const callCell = makeCell({ seq: 2, kind: "tool_call", tool_id: "call_1", tool_name: "bash", text: `{"command":"awk -f scorer.awk gen-*/decls"}` })
const resultCell = makeCell({ seq: 3, kind: "tool_result", tool_id: "call_1", tool_name: "bash", status: "ok", text: "3/5 packages flagged empty\nline2\nline3" })

test("tool pair collapsed: call header + result preview", async () => {
  const f = await frameOf(() => <EntryView entry={renderEntries([callCell, resultCell])[0]!} />)
  expect(f).toContain("▸ bash")
  expect(f).toContain("(→ expand)")
  expect(f).toContain("✓ ok")
  expect(f).toContain("3/5 packages flagged empty")
  expect(f).toMatchSnapshot()
})

test("tool pair expanded: args as code + full output", async () => {
  const f = await frameOf(() => <EntryView entry={renderEntries([callCell, resultCell])[0]!} expanded />)
  expect(f).toContain("▾ bash")
  expect(f).toContain(`"command"`)
  expect(f).toContain("line3")
  expect(f).toMatchSnapshot()
})

test("orphaned tool_call renders without result and without spinner", async () => {
  const f = await frameOf(() => <EntryView entry={renderEntries([callCell])[0]!} />)
  expect(f).toContain("▸ bash")
  expect(f).toContain("└ · no result")
  expect(f).toMatchSnapshot()
})

test("pairing consumes the result cell (single entry)", () => {
  const entries = renderEntries([callCell, resultCell])
  expect(entries).toHaveLength(1)
  expect(entries[0]!.kind).toBe("pair")
})

// --- registers ---------------------------------------------------------------

test("streaming register renders rail, thinking line and cursor block", async () => {
  const live = {
    generation: 3,
    parts: [
      { type: "stream_delta" as const, generation: 3, seq: 1, kind: "reasoning_delta", text: "the gate seems wrong", tool_id: null, tool_name: null, ts: 1 },
      { type: "stream_delta" as const, generation: 3, seq: 2, kind: "text_delta", text: "checking the score no", tool_id: null, tool_name: null, ts: 2 },
      { type: "stream_delta" as const, generation: 3, seq: 3, kind: "tool_call_delta", text: `{"command":"awk`, tool_id: "call_9", tool_name: "bash", ts: 3 },
    ],
  }
  const f = await frameOf(() => <StreamingView live={live} />)
  expect(f).toContain("▒ streaming · gen 3")
  expect(f).toContain("▸ thinking:")
  expect(f).toContain("checking the score no▊")
  expect(f).toContain("▸ bash forming…")
  expect(f).toMatchSnapshot()
})

test("interrupted generation flips into an interrupted cell", async () => {
  let state = initialState()
  state = reduceQcli(state, { type: "stream_delta", generation: 1, seq: 1, kind: "text_delta", text: "half-finished thought", tool_id: null, tool_name: null, ts: 1 })
  state = reduceQcli(state, { type: "stream_delta", generation: 2, seq: 1, kind: "text_delta", text: "fresh start", tool_id: null, tool_name: null, ts: 2 })
  expect(state.cells).toHaveLength(1)
  expect(state.cells[0]!.register).toBe("interrupted")

  const f = await frameOf(() => <EntryView entry={renderEntries(state.cells)[0]!} />)
  expect(f).toContain("▒ interrupted — not persisted")
  expect(f).toContain("half-finished thought")
  expect(f).toMatchSnapshot()
})

test("compaction divider demotes prior cells", async () => {
  const { store } = fixtureStore({
    demoted: [
      {
        reason: "compaction",
        cells: [makeCell({ seq: 1, kind: "message", text: "old cognition", role: "assistant" }) as RenderCell],
      },
    ],
    cells: [makeCell({ seq: 10, kind: "message", text: "fresh cognition", role: "assistant" }) as RenderCell],
  })
  const f = await frameOf(() => <App store={store} />, { ready: (frame) => frame.includes("fresh cognition") && frame.includes("old cognition") })
  expect(f).toContain("── context compaction")
  expect(f).toContain("no longer the agent's working memory")
  expect(f).toContain("old cognition")
  expect(f).toContain("── live context ──")
  expect(f).toContain("fresh cognition")
  expect(f).toMatchSnapshot()
})

// --- app chrome ---------------------------------------------------------------

test("app renders header, lane, footer and composer", async () => {
  const { store } = fixtureStore({
    cells: [
      { ...(makeCell({ seq: 1, kind: "human_post", author: "human", action: "inject", text: "look at the empty-package bug" }) as RenderCell), receipt: { queued: true, received: true, woke: true, delivered: true } },
      makeCell({ seq: 2, kind: "message", role: "assistant", text: "I'll check the novelty gate." }) as RenderCell,
    ],
  })
  const f = await frameOf(() => <App store={store} />, { ready: (frame) => frame.includes("novelty gate") && frame.includes("[▷✓⚡▣]") })
  expect(f).toContain("qcli · ")
  expect(f).toContain("gpt-5.5")
  expect(f).toContain("tokens: n/a")
  expect(f).toContain("[enter] send")
  expect(f).toContain("inject >")
  expect(f).toContain("[▷✓⚡▣]")
  expect(f).toMatchSnapshot()
})

test("CJK content renders with correct width handling", async () => {
  const { store } = fixtureStore({
    cells: [
      { ...(makeCell({ seq: 1, kind: "human_post", author: "human", action: "inject", text: "看一下空包 bug，novelty gate 漏了空包" }) as RenderCell), receipt: { queued: true, received: true } },
      makeCell({ seq: 2, kind: "message", role: "assistant", text: "我来查 novelty gate 的空包问题——先跑 scorer 复核。" }) as RenderCell,
    ],
  })
  const f = await frameOf(() => <App store={store} />, { ready: (frame) => frame.includes("novelty gate 漏了空包") && frame.includes("先跑 scorer") })
  expect(f).toContain("看一下空包 bug")
  expect(f).toContain("我来查 novelty gate")
  expect(f).toMatchSnapshot()
})

test("narrow terminal (60 cols) still renders coherently", async () => {
  const { store } = fixtureStore({
    cells: [
      makeCell({ seq: 1, kind: "message", role: "assistant", text: "a moderately long line of cognition that must wrap on narrow terminals" }) as RenderCell,
      makeCell({ seq: 2, kind: "tool_call", tool_id: "c1", tool_name: "bash", text: `{"command":"echo narrow"}` }) as RenderCell,
    ],
  })
  const f = await frameOf(() => <App store={store} />, { width: 60, height: 24, ready: (frame) => frame.includes("cognition") && frame.includes("▸ bash") })
  expect(f).toContain("cognition")
  expect(f).toContain("▸ bash")
  expect(f).toMatchSnapshot()
})
