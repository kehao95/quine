// Browser composer (WEB.md §5): Enter sends, Shift+Enter newline (native —
// no terminal keyboard-protocol caveats), action radio chips carrying the
// peer contract's own descriptions, and a separate guarded interrupt button
// (arm → 2s auto-disarm → confirm click).
import { createSignal, For, onCleanup, Show } from "solid-js"
import type { ControlAction } from "../store"
import type { QcliStore } from "../store"
import { actionDescription } from "../store"

const ACTIONS: ControlAction[] = ["post", "poke", "inject", "interrupt"]

export function Composer(props: { store: QcliStore; disabled?: boolean }) {
  const [action, setAction] = createSignal<ControlAction>("inject")
  const [armed, setArmed] = createSignal(false)
  let textRef: HTMLTextAreaElement | undefined
  let disarmTimer: ReturnType<typeof setTimeout> | undefined
  onCleanup(() => clearTimeout(disarmTimer))

  async function send() {
    const text = textRef?.value ?? ""
    const a = action()
    if (!text.trim() && a !== "poke" && a !== "interrupt") return
    const res = await props.store.command({ op: "send", action: a, text })
    if (res.type !== "error" && textRef) textRef.value = ""
  }

  async function interrupt() {
    if (!armed()) {
      setArmed(true)
      clearTimeout(disarmTimer)
      disarmTimer = setTimeout(() => setArmed(false), 2000)
      return
    }
    setArmed(false)
    clearTimeout(disarmTimer)
    const text = textRef?.value ?? ""
    const res = await props.store.command({ op: "send", action: "interrupt", text })
    if (res.type !== "error" && text && textRef) textRef.value = ""
  }

  function onKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey && !e.isComposing) {
      e.preventDefault()
      void send()
    }
  }

  return (
    <div class="composer">
      <div class="actions">
        <For each={ACTIONS}>
          {(a) => (
            <label class={action() === a ? "on" : ""} title={actionDescription(props.store.state, a)}>
              <input type="radio" name="action" value={a} checked={action() === a} onChange={() => setAction(a)} />
              {a}
            </label>
          )}
        </For>
        <span class="hint">{actionDescription(props.store.state, action())}</span>
      </div>
      <div class="row">
        <textarea
          ref={textRef}
          placeholder={props.disabled ? "peer exited — composer disabled" : "type a message — Enter sends, Shift+Enter newline"}
          disabled={props.disabled}
          onKeyDown={onKeyDown}
        />
        <button class="send" onClick={() => void send()} disabled={props.disabled}>
          send
        </button>
        <button class={`interrupt ${armed() ? "armed" : ""}`} onClick={() => void interrupt()} disabled={props.disabled}>
          <Show when={armed()} fallback={<>⏸ interrupt</>}>click again to interrupt</Show>
        </button>
      </div>
      <div class="keys">Enter send · Shift+Enter newline · replies are agent-chosen (0..N, maybe never)</div>
    </div>
  )
}
