// Composer behavior (WEB.md §6): Enter sends, Shift+Enter does not, action
// chips select the ctl action, and interrupt is a guarded two-click with a
// 2s auto-disarm.
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest"
import { Composer } from "../src/components/Composer"
import { fixtureStore, mount } from "./fixtures"

function keydown(el: HTMLElement, init: KeyboardEventInit) {
  el.dispatchEvent(new KeyboardEvent("keydown", { bubbles: true, cancelable: true, ...init }))
}

describe("composer", () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  test("Enter sends the composer text with the selected action", async () => {
    const { store, commands } = fixtureStore()
    const { container, dispose } = mount(() => <Composer store={store} />)
    const ta = container.querySelector("textarea")!
    ta.value = "hello quine"
    keydown(ta, { key: "Enter" })
    await vi.waitFor(() => expect(commands).toHaveLength(1))
    expect(commands[0]).toEqual({ op: "send", action: "inject", text: "hello quine" })
    await vi.waitFor(() => expect(ta.value).toBe(""))
    dispose()
  })

  test("Shift+Enter does not send", async () => {
    const { store, commands } = fixtureStore()
    const { container, dispose } = mount(() => <Composer store={store} />)
    const ta = container.querySelector("textarea")!
    ta.value = "line"
    keydown(ta, { key: "Enter", shiftKey: true })
    await Promise.resolve()
    expect(commands).toHaveLength(0)
    dispose()
  })

  test("action chips switch the ctl action; empty post does not send", async () => {
    const { store, commands } = fixtureStore()
    const { container, dispose } = mount(() => <Composer store={store} />)
    const chips = Array.from(container.querySelectorAll(".actions label"))
    const post = chips.find((c) => c.textContent?.includes("post"))!
    ;(post.querySelector("input") as HTMLInputElement).click()
    const ta = container.querySelector("textarea")!
    keydown(ta, { key: "Enter" })
    await Promise.resolve()
    expect(commands).toHaveLength(0) // empty post rejected client-side
    ta.value = "queued note"
    keydown(ta, { key: "Enter" })
    await vi.waitFor(() => expect(commands).toHaveLength(1))
    expect(commands[0]).toEqual({ op: "send", action: "post", text: "queued note" })
    dispose()
  })

  test("interrupt arms on first click, sends on second, auto-disarms after 2s", async () => {
    const { store, commands } = fixtureStore()
    const { container, dispose } = mount(() => <Composer store={store} />)
    const btn = container.querySelector(".interrupt") as HTMLButtonElement

    btn.click()
    expect(btn.className).toContain("armed")
    expect(commands).toHaveLength(0)
    await vi.advanceTimersByTimeAsync(2100)
    expect(btn.className).not.toContain("armed")
    expect(commands).toHaveLength(0) // auto-disarm fired, nothing sent

    btn.click()
    btn.click()
    await vi.waitFor(() => expect(commands).toHaveLength(1))
    expect(commands[0]).toEqual({ op: "send", action: "interrupt", text: "" })
    dispose()
  })

  test("disabled composer blocks input when the peer is dead", () => {
    const { store } = fixtureStore()
    const { container, dispose } = mount(() => <Composer store={store} disabled />)
    expect((container.querySelector("textarea") as HTMLTextAreaElement).disabled).toBe(true)
    expect((container.querySelector(".send") as HTMLButtonElement).disabled).toBe(true)
    dispose()
  })
})
