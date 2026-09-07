import { expect, test } from "bun:test"
import { testRender } from "@opentui/solid"
import type { TestRendererSetup } from "@opentui/core/testing"
import { App } from "../src/App"
import type { FixtureHandle } from "./fixtures"
import { fixtureStatus, fixtureStore, makeCell } from "./fixtures"

async function mount(handle: FixtureHandle, width = 100, height = 30): Promise<TestRendererSetup> {
  const app = await testRender(() => <App store={handle.store} />, { width, height })
  await settle(app)
  return app
}

async function settle(app: TestRendererSetup, passes = 4) {
  for (let i = 0; i < passes; i++) {
    await app.renderOnce()
    await new Promise((resolve) => setTimeout(resolve, 15))
  }
}

// pump renders until the condition holds — key parsing (lone ESC) and async
// command dispatch need real time, not a fixed pass count.
async function pump(app: TestRendererSetup, cond: () => boolean, passes = 60): Promise<void> {
  for (let i = 0; i < passes; i++) {
    await app.renderOnce()
    if (cond()) return
    await new Promise((resolve) => setTimeout(resolve, 15))
  }
}

// The §6 driven-key regression gate: a plain Enter on the focused composer
// must issue POST /command {op:"send"} — this is the exact bug the first M2
// attempt shipped (Textarea defaults map return→newline, meta+return→submit).
test("plain Enter sends the composer text", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  await app.mockInput.typeText("hello quine")
  await settle(app, 2)
  app.mockInput.pressEnter()
  await pump(app, () => handle.commands.some((c) => c.op === "send"))
  app.renderer.destroy()

  const sends = handle.commands.filter((c) => c.op === "send")
  expect(sends).toHaveLength(1)
  expect(sends[0]).toMatchObject({ op: "send", action: "inject", text: "hello quine" })
})

test("alt+enter inserts a newline instead of sending", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  await app.mockInput.typeText("line-one")
  app.mockInput.pressEnter({ meta: true })
  await settle(app, 1)
  await app.mockInput.typeText("line-two")
  await settle(app, 2)
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(handle.commands.filter((c) => c.op === "send")).toHaveLength(0)
  expect(frame).toContain("line-one")
  expect(frame).toContain("line-two")
})

test("typing ? into a non-empty composer does not open help", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  await app.mockInput.typeText("what?")
  await settle(app, 2)
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(frame).not.toContain("qcli help")
  expect(frame).toContain("what?")
})

test("? on an empty composer opens the help overlay", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  await app.mockInput.typeText("?")
  await pump(app, () => app.captureCharFrame().includes("qcli help"))
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(frame).toContain("qcli help")
  expect(frame).toContain("receipt ladder")
  expect(frame).toMatchSnapshot()
})

test("tab cycles the delivery action chip", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  expect(app.captureCharFrame()).toContain("inject >")
  app.mockInput.pressTab()
  await pump(app, () => app.captureCharFrame().includes("interrupt >"))
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(frame).toContain("interrupt >")
})

test("esc arms the interrupt chip; second esc sends an empty interrupt", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  app.mockInput.pressEscape()
  await pump(app, () => app.captureCharFrame().includes("Esc again = interrupt"))
  expect(app.captureCharFrame()).toContain("Esc again = interrupt")
  expect(app.captureCharFrame()).toMatchSnapshot()
  app.mockInput.pressEscape()
  await pump(app, () => handle.commands.some((c) => c.op === "send"))
  app.renderer.destroy()

  const sends = handle.commands.filter((c) => c.op === "send")
  expect(sends).toHaveLength(1)
  expect(sends[0]).toMatchObject({ op: "send", action: "interrupt", text: "" })
})

test("up-arrow at the boundary recalls history", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  await app.mockInput.typeText("historic-payload-xyz")
  await settle(app, 1)
  app.mockInput.pressEnter()
  await pump(app, () => handle.commands.some((c) => c.op === "send"))
  app.mockInput.pressArrow("up")
  await settle(app, 3)
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(handle.commands.filter((c) => c.op === "send")).toHaveLength(1)
  expect(frame).toContain("historic-payload-xyz")
})

test("ctrl+up focuses a cell and right expands it", async () => {
  const handle = fixtureStore({
    cells: [
      makeCell({ seq: 1, kind: "reasoning", text: "first thought line\nsecond thought line\nthird thought line", author: null }),
    ],
  })
  const app = await mount(handle)
  app.mockInput.pressArrow("up", { ctrl: true })
  await pump(app, () => app.captureCharFrame().includes("cells"))
  expect(app.captureCharFrame()).toContain("[↑↓] cells")
  app.mockInput.pressArrow("right")
  await pump(app, () => app.captureCharFrame().includes("third thought line"))
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(frame).toContain("▾ reasoning")
  expect(frame).toContain("third thought line")
})

test("ctrl+p opens the roster overlay with lineage and liveness", async () => {
  const handle = fixtureStore()
  const app = await mount(handle)
  app.mockInput.pressKey("p", { ctrl: true })
  await pump(app, () => app.captureCharFrame().includes("peers"))
  const frame = app.captureCharFrame()
  app.renderer.destroy()

  expect(frame).toContain("peers")
  expect(frame).toContain("← you")
  expect(frame).toContain("◍")
  expect(frame).toMatchSnapshot()
})

test("dead peer disables the composer and sends nothing", async () => {
  const handle = fixtureStore({ statusSnapshot: fixtureStatus({ live: false }) })
  const app = await mount(handle)
  const frame = app.captureCharFrame()
  await app.mockInput.typeText("x")
  await settle(app, 1)
  app.mockInput.pressEnter()
  await settle(app, 2)
  app.renderer.destroy()

  expect(frame).toContain("peer is dead")
  expect(frame).toContain("composer disabled")
  expect(frame).toMatchSnapshot()
  expect(handle.commands.filter((c) => c.op === "send")).toHaveLength(0)
})
