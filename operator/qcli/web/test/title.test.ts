// Unread title badge logic (WEB.md §6).
import { describe, expect, test } from "vitest"
import { computeTitle } from "../src/title"

describe("computeTitle", () => {
  test("no peer, no unread", () => {
    expect(computeTitle(0, null)).toBe("qcli")
  })
  test("peer without unread", () => {
    expect(computeTitle(0, "sess_demo_aaaa1111")).toBe("qcli — o_aaaa1111")
  })
  test("unread prefixes a count badge", () => {
    expect(computeTitle(2, "sess_demo_aaaa1111")).toBe("(2) qcli — o_aaaa1111")
  })
})
