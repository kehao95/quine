// Inspector pane (WEB.md §1/§5 — the "see everything" surface): every
// contract field untruncated, the receipt timeline with timestamps, and the
// raw event JSON with a copy button. Esc deselects (handled by App).
import { For, Show } from "solid-js"
import type { RenderCell } from "../store"
import { receiptSegments } from "../store"

const FIELD_ORDER: Array<keyof RenderCell> = [
  "kind",
  "seq",
  "role",
  "action",
  "author",
  "tool_name",
  "tool_id",
  "status",
  "is_error",
  "ts",
  "termination_mode",
  "exit_code",
  "turn_count",
  "tokens_in",
  "tokens_out",
  "register",
  "client_ref",
  "error_code",
]

function fmtTs(ts: number | null | undefined): string {
  if (ts == null) return ""
  const d = new Date(ts)
  return `${d.toLocaleTimeString(undefined, { hour12: false })}.${String(d.getMilliseconds()).padStart(3, "0")}`
}

export function Inspector(props: { cell: RenderCell; onClose: () => void }) {
  const c = () => props.cell
  const raw = () => JSON.stringify(c(), null, 2)
  async function copy() {
    try {
      await navigator.clipboard?.writeText(raw())
    } catch {
      // Clipboard access is best-effort (may be unavailable in some contexts).
    }
  }
  return (
    <div class="inspector" data-testid="inspector">
      <button class="close" onClick={() => props.onClose()} title="close (Esc)">
        ✕
      </button>
      <h3>cell detail</h3>
      <table class="fields">
        <tbody>
          <For each={FIELD_ORDER.filter((k) => c()[k] != null && c()[k] !== "")}>
            {(k) => (
              <tr>
                <td class="k">{k}</td>
                <td class="mono">{k === "ts" ? `${c().ts} (${fmtTs(c().ts)})` : String(c()[k])}</td>
              </tr>
            )}
          </For>
        </tbody>
      </table>

      <Show when={c().text != null && c().text !== ""}>
        <div class="section">full text</div>
        <pre class="block full" data-testid="inspector-text">{c().text}</pre>
      </Show>
      <Show when={c().raw}>
        <div class="section">raw line</div>
        <pre class="block full">{c().raw}</pre>
      </Show>

      <Show when={c().kind === "human_post"}>
        <div class="section">delivery timeline</div>
        <div class="timeline">
          <For each={receiptSegments(c().receipt)}>
            {(seg) => (
              <div class={`stage ${seg.lit ? "" : "unlit"}`}>
                <span>{seg.glyph}</span>
                <span>{seg.label}</span>
                <span class="ts">{seg.lit ? (seg.ts != null ? fmtTs(seg.ts) : "—") : "not reached"}</span>
              </div>
            )}
          </For>
          <Show when={c().receipt?.message_id}>
            <div class="stage">
              <span class="ts">message_id: {c().receipt?.message_id}</span>
            </div>
          </Show>
        </div>
      </Show>

      <div class="section">raw event json</div>
      <button class="copy" onClick={() => void copy()}>
        copy json
      </button>
      <pre class="block full">{raw()}</pre>
    </div>
  )
}
