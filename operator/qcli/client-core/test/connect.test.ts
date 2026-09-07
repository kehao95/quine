// Fixture fetch-SSE driver (WEB.md §6): connectSSE against a mocked fetch
// streaming recorded frames — no live kernel needed. Proves frame parsing,
// dispatch, and abort behavior for every client that shares connect.ts.
import { expect, test } from "bun:test"
import { connectSSE } from "../connect"
import type { QcliEvent } from "../contract"

function sseResponse(frames: string[]): Response {
  const enc = new TextEncoder()
  const stream = new ReadableStream<Uint8Array>({
    start(controller) {
      for (const frame of frames) controller.enqueue(enc.encode(frame))
      controller.close()
    },
  })
  return new Response(stream, { status: 200, headers: { "content-type": "text/event-stream" } })
}

test("connectSSE parses frames, dispatches typed events, and stops on abort", async () => {
  const frames = [
    `data: {"type":"hello","contract":"qcli/1","endpoint_id":"qcli-1","peer":{"session":"s","run_id":"r","pid":1,"agent_root":"/a","runtime_root":"/r"}}\n\n`,
    // split one frame across two chunks to prove partial-buffer handling
    `data: {"type":"cell","seq":1,"kind":"message","role":"assistant","text":"hi","ts":1,`,
    `"tool_id":null,"tool_name":null,"status":null,"is_error":false,"action":null,"author":null,"raw":null}\n\n`,
    `data: not-json\n\n`,
    `data: {"type":"backfill_complete","count":1}\n\n`,
  ]
  const events: QcliEvent[] = []
  const states: string[] = []
  const ac = new AbortController()
  const realFetch = globalThis.fetch
  globalThis.fetch = (async () => sseResponse(frames)) as unknown as typeof fetch
  try {
    const done = connectSSE(
      "http://mock",
      (ev) => {
        events.push(ev)
        if (ev.type === "backfill_complete") ac.abort()
      },
      (s) => states.push(s),
      ac.signal,
    )
    await done
  } finally {
    globalThis.fetch = realFetch
  }
  expect(events.map((e) => e.type)).toEqual(["hello", "cell", "backfill_complete"])
  expect(states).toContain("rebuilding")
  expect(states.at(-1)).toBe("disconnected")
})
