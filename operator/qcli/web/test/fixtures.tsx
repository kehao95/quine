// Web test fixtures: pure data from client-core + a recording store bound to
// THIS package's solid instance, plus a DOM render helper for vitest+jsdom.
import { createStore, unwrap } from "solid-js/store"
import { render } from "solid-js/web"
import type { JSX } from "solid-js"
import type { Command, CommandResponse } from "../src/store"
import type { QcliState, QcliStore } from "../src/store"
import { initialState, reduceQcli } from "../src/store"
import { fixturePeer, fixtureRoster, fixtureStatus } from "../../client-core/fixtures"

export { fixturePeer, fixtureRoster, fixtureStatus, makeCell } from "../../client-core/fixtures"

export interface FixtureHandle {
  store: QcliStore
  commands: Command[]
}

export function fixtureStore(partial: Partial<QcliState> = {}): FixtureHandle {
  const [state, setState] = createStore<QcliState>({
    ...initialState(),
    peer: fixturePeer,
    statusSnapshot: fixtureStatus(),
    connection: "connected",
    backfilled: true,
    roster: fixtureRoster,
    ...partial,
  })
  const commands: Command[] = []
  const store: QcliStore = {
    state,
    connect: async () => {},
    command: async (cmd: Command): Promise<CommandResponse> => {
      commands.push(cmd)
      const action = cmd.op === "send" ? cmd.action : "inject"
      return { type: "queued", action, client_ref: `c-${commands.length}` }
    },
    fetchRoster: async () => state.roster,
    fetchPeerContract: async () => null,
    dispatch: (event) => setState(reduceQcli(unwrap(state), event)),
    clearUnread: () => setState("unread", 0),
  }
  return { store, commands }
}

export function mount(ui: () => JSX.Element): { container: HTMLElement; dispose: () => void } {
  const container = document.createElement("div")
  document.body.appendChild(container)
  const dispose = render(ui, container)
  return {
    container,
    dispose: () => {
      dispose()
      container.remove()
    },
  }
}
