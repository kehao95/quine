import { render } from "@opentui/solid"
import { App } from "./App"
import { createQcliStore } from "./store"

const baseUrl = process.env.QCLI_SERVER ?? "http://127.0.0.1:7777"
const store = createQcliStore(baseUrl)

await render(() => <App store={store} />, {
  exitOnCtrlC: false,
  targetFps: 60,
  autoFocus: false,
  useMouse: true,
})
