// Fetch-based SSE lifecycle shared by every client. Framework-free: only
// fetch/ReadableStream — usable from bun (TUI) and the browser (web) alike.
import type { ConnectionState } from "./reducer"
import type { QcliEvent } from "./contract"
import { isQcliEvent } from "./contract"

export async function connectSSE(
  baseUrl: string,
  dispatch: (event: QcliEvent) => void,
  setConnection: (state: ConnectionState) => void,
  signal: AbortSignal,
): Promise<void> {
  let delay = 250
  while (!signal.aborted) {
    try {
      setConnection(delay === 250 ? "connecting" : "reconnecting")
      const res = await fetch(baseUrl + "/events", { signal, headers: { accept: "text/event-stream" } })
      if (!res.ok || !res.body) throw new Error("events HTTP " + res.status)
      setConnection("rebuilding")
      delay = 250
      const reader = res.body.getReader()
      const dec = new TextDecoder()
      let buf = ""
      for (;;) {
        const { value, done } = await reader.read()
        if (done) break
        buf += dec.decode(value, { stream: true })
        let idx: number
        while ((idx = buf.indexOf("\n\n")) >= 0) {
          const frame = buf.slice(0, idx)
          buf = buf.slice(idx + 2)
          for (const line of frame.split("\n")) {
            if (!line.startsWith("data:")) continue
            const payload = line.slice(5).trim()
            if (!payload) continue
            try {
              const parsed = JSON.parse(payload)
              if (isQcliEvent(parsed)) dispatch(parsed)
            } catch {
              // Ignore malformed frames; reconnect rebuilds state if the stream breaks.
            }
          }
        }
      }
    } catch {
      if (signal.aborted) break
    }
    setConnection("reconnecting")
    await sleep(delay, signal)
    delay = Math.min(delay * 2, 4000)
  }
  setConnection("disconnected")
}

function sleep(ms: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve) => {
    const t = setTimeout(resolve, ms)
    signal.addEventListener(
      "abort",
      () => {
        clearTimeout(t)
        resolve()
      },
      { once: true },
    )
  })
}
