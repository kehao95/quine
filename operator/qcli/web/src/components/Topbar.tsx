// Topbar: peer identity + status chips + connection state (WEB.md §4).
import { Show } from "solid-js"
import type { QcliStore } from "../store"
import { shortSession } from "../store"

export function Topbar(props: { store: QcliStore }) {
  const s = () => props.store.state
  const status = () => s().statusSnapshot
  return (
    <div class="topbar">
      <span>
        {status()?.live === false ? "◍" : "●"} {s().peer ? shortSession(s().peer!.session) : "no peer"}
      </span>
      <div class="chips">
        <Show when={s().peer}>
          <span>pid {status()?.pid ?? s().peer!.pid}</span>
        </Show>
        <Show when={status()}>
          <span>{status()!.model}</span>
          <span>depth {status()!.depth}</span>
          <span>inc {status()!.incarnation}</span>
          <span>pending {status()!.pending}</span>
          <span>{status()!.working ? "working" : "idle"}</span>
          <span>
            tokens {status()!.tokens_in == null ? "n/a" : `${status()!.tokens_in}/${status()!.tokens_out}`}
          </span>
        </Show>
      </div>
      <div class="spacer" />
      <Show when={s().unread > 0}>
        <span class="unread-badge" data-testid="unread-badge">
          {s().unread}
        </span>
      </Show>
      <span class="conn">{s().connection}</span>
    </div>
  )
}
