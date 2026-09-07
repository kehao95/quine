// Inspector detail (WEB.md §6): untruncated rendering for a >10KB
// tool_result, field table, and the receipt timeline with timestamps.
import { describe, expect, test } from "vitest"
import { Inspector } from "../src/components/Inspector"
import type { RenderCell } from "../src/store"
import { makeCell, mount } from "./fixtures"

describe("inspector", () => {
  test("renders a >10KB tool_result untruncated", () => {
    const big = Array.from({ length: 600 }, (_, i) => `line-${i}: ${"x".repeat(20)}`).join("\n")
    expect(big.length).toBeGreaterThan(10_000)
    const cell = makeCell({ seq: 1, kind: "tool_result", tool_id: "t1", tool_name: "bash", status: "ok", text: big }) as RenderCell
    const { container, dispose } = mount(() => <Inspector cell={cell} onClose={() => {}} />)
    const full = container.querySelector('[data-testid="inspector-text"]')
    expect(full?.textContent?.length).toBeGreaterThanOrEqual(big.length)
    expect(full?.textContent).toContain("line-599")
    dispose()
  })

  test("shows contract fields and raw event json", () => {
    const cell = makeCell({ seq: 2, kind: "tool_call", tool_id: "call_1", tool_name: "bash", text: "{}" }) as RenderCell
    const { container, dispose } = mount(() => <Inspector cell={cell} onClose={() => {}} />)
    const fields = container.querySelector("table.fields")?.textContent ?? ""
    expect(fields).toContain("tool_name")
    expect(fields).toContain("bash")
    expect(fields).toContain("call_1")
    expect(container.textContent).toContain("raw event json")
    expect(container.textContent).toContain(`"kind": "tool_call"`)
    dispose()
  })

  test("human_post shows the delivery timeline with stage timestamps", () => {
    const cell = makeCell({ seq: 3, kind: "human_post", action: "inject", text: "hi" }) as RenderCell
    cell.receipt = { queued: true, queued_ts: 1783044370000, received: true, received_ts: 1783044370234, message_id: "message-1" }
    const { container, dispose } = mount(() => <Inspector cell={cell} onClose={() => {}} />)
    const timeline = container.querySelector(".timeline")?.textContent ?? ""
    expect(timeline).toContain("queued")
    expect(timeline).toContain("received")
    expect(timeline).toContain("not reached") // woke + delivered unlit
    expect(timeline).toContain("message-1")
    expect(container.querySelectorAll(".timeline .stage.unlit").length).toBe(2)
    dispose()
  })

  test("close button fires onClose", () => {
    const cell = makeCell({ seq: 4, kind: "message", text: "x" }) as RenderCell
    let closed = false
    const { container, dispose } = mount(() => <Inspector cell={cell} onClose={() => (closed = true)} />)
    ;(container.querySelector(".close") as HTMLElement).click()
    expect(closed).toBe(true)
    dispose()
  })
})
