// The framework-free data layer (reducer table, pairing, ladder, live view)
// lives in the shared client-core (WEB.md §3). This module re-exports it and
// keeps only the Solid-reactive store factory — each client wraps the shared
// reducer with its OWN solid-js instance (a shared reactive core across
// package boundaries would mean two solid runtimes fighting).
import { createStore, unwrap } from "solid-js/store"
import type { Command, CommandResponse, PeerContractResponse, QcliEvent, RosterEntry, RosterResponse } from "../../client-core/contract"
import type { QcliState, QcliStore } from "../../client-core/reducer"
import { applyLocalQueued, initialState, reduceQcli } from "../../client-core/reducer"
import { connectSSE } from "../../client-core/connect"

export * from "../../client-core/reducer"
export { connectSSE } from "../../client-core/connect"

export function createQcliStore(baseUrl: string): QcliStore {
  const [state, setState] = createStore<QcliState>(initialState())
  const dispatch = (event: QcliEvent) => setState(reduceQcli(unwrap(state), event))

  async function connect(signal: AbortSignal) {
    await connectSSE(baseUrl, dispatch, (connection) => setState("connection", connection), signal)
  }

  async function command(cmd: Command): Promise<CommandResponse> {
    const res = await fetch(baseUrl + "/command", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(cmd),
    })
    const body = (await res.json()) as CommandResponse
    if (body.type === "error") {
      dispatch(body)
      return body
    }
    if (body.type === "queued" && cmd.op === "send") {
      setState(applyLocalQueued(unwrap(state), body.action, cmd.text, body.client_ref))
    }
    return body
  }

  async function fetchRoster(): Promise<RosterEntry[]> {
    try {
      const res = await fetch(baseUrl + "/roster")
      if (!res.ok) return []
      const body = (await res.json()) as RosterResponse
      const peers = Array.isArray(body.peers) ? body.peers : []
      setState("roster", peers)
      return peers
    } catch {
      return []
    }
  }

  async function fetchPeerContract(): Promise<PeerContractResponse["contract"] | null> {
    try {
      const res = await fetch(baseUrl + "/peer-contract")
      if (!res.ok) return null
      const body = (await res.json()) as PeerContractResponse
      setState("peerContract", body.contract)
      return body.contract
    } catch {
      return null
    }
  }

  const clearUnread = () => setState("unread", 0)

  return { state, connect, command, fetchRoster, fetchPeerContract, dispatch, clearUnread }
}
