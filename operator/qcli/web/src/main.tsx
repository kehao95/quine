import { render } from "solid-js/web"
import { App } from "./App"
import { createQcliStore } from "./store"
import "./theme.css"

// Same-origin: the kernel serves this app at / next to the qcli/1 endpoints
// (WEB.md §2). Dev mode proxies instead (vite.config.ts).
const store = createQcliStore("")

render(() => <App store={store} />, document.getElementById("root")!)
