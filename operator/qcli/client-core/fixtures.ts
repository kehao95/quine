// Shared contract fixtures (pure data + builders). The reactive fixtureStore
// helpers live with each client — client-core stays framework-free, and a
// shared solid instance across packages would be a dual-instance trap.
import type { Cell, RosterEntry } from "./contract"
import type { QcliState } from "./reducer"

export function makeCell(partial: Partial<Cell> & { seq: number; kind: string }): Cell {
  return {
    type: "cell",
    role: null,
    text: null,
    ts: null,
    tool_id: null,
    tool_name: null,
    status: null,
    is_error: false,
    action: null,
    author: null,
    raw: null,
    ...partial,
  }
}

export const fixturePeer = {
  session: "sess_demo_aaaa1111",
  run_id: "run_demo",
  pid: 42,
  agent_root: "/tmp/a",
  runtime_root: "/tmp/r",
}

export function fixtureStatus(partial: Partial<QcliState["statusSnapshot"] & object> = {}) {
  return {
    type: "status" as const,
    live: true,
    working: true,
    pid: 42,
    ppid: 1,
    session: fixturePeer.session,
    run_id: fixturePeer.run_id,
    incarnation: 2,
    parent_session: "",
    model: "gpt-5.5",
    depth: 2,
    pending: 0,
    tokens_in: null,
    tokens_out: null,
    ...partial,
  }
}

export const fixtureRoster: RosterEntry[] = [
  {
    pid: 42,
    session: fixturePeer.session,
    run_id: fixturePeer.run_id,
    parent_session: "",
    depth: 2,
    model: "gpt-5.5",
    live: true,
    pending: 0,
    agent_root: "/tmp/a",
    attached: true,
  },
  {
    pid: 43,
    session: "sess_demo_bbbb2222",
    run_id: "run_demo2",
    parent_session: fixturePeer.session,
    depth: 3,
    model: "gpt-5.5",
    live: false,
    pending: 1,
    agent_root: "/tmp/b",
    attached: false,
  },
]
