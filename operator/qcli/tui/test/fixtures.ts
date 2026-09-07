// Pure fixtures come from the shared client-core; the reactive fixtureStore
// stays here because it must use the TUI's own solid-js instance.
import { createStore, unwrap } from "solid-js/store"
import type { Command, CommandResponse } from "../src/contract"
import type { QcliState, QcliStore } from "../src/store"
import { initialState, reduceQcli } from "../src/store"
import { fixturePeer, fixtureRoster, fixtureStatus } from "../../client-core/fixtures"

export { fixturePeer, fixtureRoster, fixtureStatus, makeCell } from "../../client-core/fixtures"

export interface FixtureHandle {
  store: QcliStore
  commands: Command[]
}

// fixtureStore builds a reactive QcliStore whose command() records instead of
// fetching — the driven-key tests assert against `commands`.
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
