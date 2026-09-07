// Transcript view filters + search (pure logic and DOM), demoted blocks with
// dividers, and the empty state (WEB.md §6).
import { describe, expect, test } from "vitest"
import { filterEntries, Transcript } from "../src/components/Transcript"
import type { RenderCell } from "../src/store"
import { renderEntries } from "../src/store"
import { fixtureStore, makeCell, mount } from "./fixtures"

const cells = [
  makeCell({ seq: 1, kind: "human_post", action: "inject", text: "find the bug" }),
  makeCell({ seq: 2, kind: "reasoning", text: "thinking about it" }),
  makeCell({ seq: 3, kind: "message", text: "I found the bug in scorer.awk" }),
  makeCell({ seq: 4, kind: "tool_call", tool_id: "t1", tool_name: "bash", text: `{"command":"awk"}` }),
  makeCell({ seq: 5, kind: "tool_result", tool_id: "t1", tool_name: "bash", status: "ok", text: "3 flagged" }),
  makeCell({ seq: 6, kind: "reply", text: "done" }),
  makeCell({ seq: 7, kind: "outcome", termination_mode: "exit" }),
] as RenderCell[]

const allOn = { messages: true, reasoning: true, tools: true, human: true, outcome: true }

describe("filterEntries", () => {
  test("kind groups filter entries", () => {
    const entries = renderEntries(cells)
    expect(entries).toHaveLength(6) // pair collapses two cells
    expect(filterEntries(entries, { ...allOn, tools: false }, "")).toHaveLength(5)
    expect(filterEntries(entries, { ...allOn, messages: false }, "")).toHaveLength(4) // drops message + reply
    expect(filterEntries(entries, { ...allOn, human: false, reasoning: false, outcome: false }, "")).toHaveLength(3)
  })

  test("substring search matches cell text and tool names", () => {
    const entries = renderEntries(cells)
    expect(filterEntries(entries, allOn, "scorer")).toHaveLength(1)
    expect(filterEntries(entries, allOn, "bash")).toHaveLength(1) // the pair, by tool name
    expect(filterEntries(entries, allOn, "zzz-nope")).toHaveLength(0)
  })
})

describe("transcript DOM", () => {
  test("renders entries, demoted block with dividers, and streaming cell", () => {
    const { store } = fixtureStore({
      cells,
      demoted: [{ reason: "compaction", cells: [makeCell({ seq: 0, kind: "message", text: "old cognition" }) as RenderCell] }],
      live: {
        generation: 2,
        parts: [{ type: "stream_delta", generation: 2, seq: 1, kind: "text_delta", text: "forming", tool_id: null, tool_name: null, ts: 1 }],
      },
    })
    const { container, dispose } = mount(() => (
      <Transcript store={store} selectedKey={null} onSelect={() => {}} onAtBottomChange={() => {}} />
    ))
    expect(container.querySelectorAll(".entry").length).toBeGreaterThanOrEqual(7)
    expect(container.textContent).toContain("context reset: compaction")
    expect(container.textContent).toContain("── live context ──")
    expect(container.textContent).toContain("old cognition")
    expect(container.textContent).toContain("streaming · gen 2")
    dispose()
  })

  test("filter chips hide kinds in the DOM", () => {
    const { store } = fixtureStore({ cells })
    const { container, dispose } = mount(() => (
      <Transcript store={store} selectedKey={null} onSelect={() => {}} onAtBottomChange={() => {}} />
    ))
    expect(container.querySelectorAll('[data-kind="reasoning"]')).toHaveLength(1)
    ;(container.querySelector('[data-testid="filter-reasoning"]') as HTMLElement).click()
    expect(container.querySelectorAll('[data-kind="reasoning"]')).toHaveLength(0)
    dispose()
  })

  test("search box narrows entries", () => {
    const { store } = fixtureStore({ cells })
    const { container, dispose } = mount(() => (
      <Transcript store={store} selectedKey={null} onSelect={() => {}} onAtBottomChange={() => {}} />
    ))
    const search = container.querySelector('[data-testid="search"]') as HTMLInputElement
    search.value = "scorer"
    search.dispatchEvent(new Event("input", { bubbles: true }))
    expect(container.querySelectorAll(".entry")).toHaveLength(1)
    dispose()
  })

  test("empty state names the next step", () => {
    const { store } = fixtureStore({ cells: [], peer: null })
    const { container, dispose } = mount(() => (
      <Transcript store={store} selectedKey={null} onSelect={() => {}} onAtBottomChange={() => {}} />
    ))
    expect(container.textContent).toContain("attach or spawn a peer")
    dispose()
  })
})
