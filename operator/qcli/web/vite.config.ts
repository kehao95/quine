/// <reference types="vitest/config" />
import { defineConfig } from "vite"
import solid from "vite-plugin-solid"

// Dev loop (WEB.md §3): `bunx vite dev` proxies the qcli/1 API routes to a
// running `qcli serve` (QCLI_PROXY overrides the target). Production build is
// served by the kernel itself at / — same origin, no proxy, no CORS.
const API_ROUTES = ["/events", "/context", "/status", "/roster", "/command", "/healthz", "/peer-contract", "/history"]
const proxyTarget = process.env.QCLI_PROXY ?? "http://127.0.0.1:7777"

export default defineConfig({
  plugins: [solid()],
  server: {
    fs: { allow: ["../.."] },
    proxy: Object.fromEntries(API_ROUTES.map((route) => [route, { target: proxyTarget, changeOrigin: false }])),
  },
  build: { target: "es2022" },
  test: {
    environment: "jsdom",
    include: ["test/**/*.test.{ts,tsx}"],
  },
  resolve: {
    conditions: ["browser", "development"],
  },
})
