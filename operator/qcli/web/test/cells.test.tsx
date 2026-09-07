// Per-kind rendering + registers (WEB.md §6): every contract cell kind has a
// DOM rendering, registers select the visual layer, markdown fences get
// syntax-highlighted, unknown kinds fall back to a raw block.
import { describe, expect, test } from "vitest"
import { Divider, EntryView, StreamingView } from "../src/components/Cell"
import type { RenderCell, RenderEntry } from "../src/store"
import { entryKey } from "../src/store"
import { makeCell, mount } from "./fixtures"

function entryFor(cell: RenderCell): RenderEntry {
  return { key: entryKey(cell), kind: "cell", cell }
}

function renderCell(cell: RenderCell): { container: HTMLElement; dispose: () => void } {
  return mount(() => <EntryView entry={entryFor(cell)} selectedKey={null} onSelect={() => {}} />)
}

describe("cell kinds", () => {
  test("message renders markdown with highlighted fences", () => {
    const { container, dispose } = renderCell(
      makeCell({ seq: 1, kind: "message", role: "assistant", text: "Some **bold** and\n\n```go\nfunc main() {}\n```" }) as RenderCell,
    )
    expect(container.querySelector(".md strong")?.textContent).toBe("bold")
    const code = container.querySelector(".md pre code")
    expect(code?.className).toContain("hljs")
    expect(code?.querySelector(".hljs-keyword")?.textContent).toBe("func")
    dispose()
  })

  test("human_post renders action tag, ladder, and text", () => {
    const cell = makeCell({ seq: 2, kind: "human_post", action: "inject", author: "human", text: "look at this" }) as RenderCell
    cell.receipt = { queued: true, received: true }
    const { container, dispose } = renderCell(cell)
    expect(container.querySelector(".action-tag")?.textContent).toContain("inject")
    expect(container.querySelectorAll(".ladder .lit")).toHaveLength(2)
    expect(container.querySelectorAll(".ladder .unlit")).toHaveLength(2)
    expect(container.textContent).toContain("look at this")
    expect(container.querySelector(".author")?.textContent).toBe("you")
    dispose()
  })

  test("reasoning collapses to a head line and expands on click", () => {
    const { container, dispose } = renderCell(
      makeCell({ seq: 3, kind: "reasoning", text: "line one is quite something\nline two\nline three" }) as RenderCell,
    )
    expect(container.textContent).toContain("reasoning · 3 lines")
    expect(container.textContent).not.toContain("line two")
    ;(container.querySelector(".reasoning .head") as HTMLElement).click()
    expect(container.textContent).toContain("line two")
    dispose()
  })

  test("tool pair nests result under call; orphan call shows no-result marker", () => {
    const call = makeCell({ seq: 4, kind: "tool_call", tool_id: "t1", tool_name: "bash", text: `{"command":"ls"}` }) as RenderCell
    const result = makeCell({ seq: 5, kind: "tool_result", tool_id: "t1", tool_name: "bash", status: "ok", text: "file-a\nfile-b" }) as RenderCell
    const pair: RenderEntry = { key: entryKey(call), kind: "pair", call, result }
    const { container, dispose } = mount(() => <EntryView entry={pair} selectedKey={null} onSelect={() => {}} />)
    expect(container.querySelector(".tool-call .head")?.textContent).toContain("bash")
    expect(container.querySelector(".tool-result .status.ok")?.textContent).toContain("ok")
    expect(container.querySelector(".tool-result pre")?.textContent).toContain("file-a")
    dispose()

    const orphan: RenderEntry = { key: "k2", kind: "pair", call, result: null }
    const { container: c2, dispose: d2 } = mount(() => <EntryView entry={orphan} selectedKey={null} onSelect={() => {}} />)
    expect(c2.textContent).toContain("no result")
    d2()
  })

  test("outcome, reply, and unknown kinds render distinctly", () => {
    const { container: c1, dispose: d1 } = renderCell(
      makeCell({ seq: 6, kind: "outcome", termination_mode: "exit", exit_code: 0, turn_count: 3 }) as RenderCell,
    )
    expect(c1.querySelector(".outcome-cell")?.textContent).toContain("exit")
    d1()
    const { container: c2, dispose: d2 } = renderCell(makeCell({ seq: 7, kind: "reply", text: "pong", author: "sess_x" }) as RenderCell)
    expect(c2.querySelector(".reply-cell")).toBeTruthy()
    d2()
    const { container: c3, dispose: d3 } = renderCell(makeCell({ seq: 8, kind: "novel_kind", raw: `{"kind":"novel_kind"}` }) as RenderCell)
    expect(c3.querySelector(".unknown-cell .raw")?.textContent).toContain("novel_kind")
    d3()
  })

  test("registers: demoted dims, interrupted marks not-persisted, local marks sending", () => {
    const demoted = makeCell({ seq: 9, kind: "message", text: "old" }) as RenderCell
    const { container: c1, dispose: d1 } = mount(() => (
      <EntryView entry={entryFor(demoted)} selectedKey={null} onSelect={() => {}} demoted />
    ))
    expect(c1.querySelector(".entry")?.className).toContain("register-demoted")
    d1()

    const interrupted = makeCell({ seq: 10, kind: "message", text: "partial" }) as RenderCell
    interrupted.register = "interrupted"
    const { container: c2, dispose: d2 } = renderCell(interrupted)
    expect(c2.querySelector(".not-persisted")?.textContent).toContain("not persisted")
    d2()

    const local = makeCell({ seq: 11, kind: "human_post", action: "inject", text: "echo" }) as RenderCell
    local.local = true
    local.register = "local"
    const { container: c3, dispose: d3 } = renderCell(local)
    expect(c3.querySelector(".sending")?.textContent).toContain("sending")
    d3()
  })

  test("streaming view shows generation, thinking line, forming tool, cursor", () => {
    const { container, dispose } = mount(() => (
      <StreamingView
        live={{
          generation: 3,
          parts: [
            { type: "stream_delta", generation: 3, seq: 1, kind: "reasoning_delta", text: "checking", tool_id: null, tool_name: null, ts: 1 },
            { type: "stream_delta", generation: 3, seq: 2, kind: "text_delta", text: "The score", tool_id: null, tool_name: null, ts: 2 },
            { type: "stream_delta", generation: 3, seq: 3, kind: "tool_call_delta", text: `{"cmd`, tool_id: "t1", tool_name: "bash", ts: 3 },
          ],
        }}
      />
    ))
    expect(container.textContent).toContain("streaming · gen 3")
    expect(container.textContent).toContain("thinking: checking")
    expect(container.textContent).toContain("The score")
    expect(container.textContent).toContain("forming: bash")
    expect(container.querySelector(".cursor")).toBeTruthy()
    dispose()
  })

  test("divider renders", () => {
    const { container, dispose } = mount(() => <Divider>── context reset ──</Divider>)
    expect(container.querySelector(".divider")?.textContent).toContain("context reset")
    dispose()
  })
})
