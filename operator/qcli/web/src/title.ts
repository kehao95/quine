// Unread attention surface (UX.md G5, WEB.md §5): the badge lives in
// document.title and a favicon dot — volatile, no Notification API in v1.
import { shortSession } from "../../client-core/contract"

export function computeTitle(unread: number, peerSession: string | null | undefined): string {
  const base = peerSession ? `qcli — ${shortSession(peerSession)}` : "qcli"
  return unread > 0 ? `(${unread}) ${base}` : base
}

const BASE_FAVICON = (() => {
  if (typeof document === "undefined") return ""
  return (document.getElementById("favicon") as HTMLLinkElement | null)?.href ?? ""
})()

export function setFaviconDot(active: boolean): void {
  if (typeof document === "undefined") return
  const link = document.getElementById("favicon") as HTMLLinkElement | null
  if (!link) return
  if (!active) {
    if (BASE_FAVICON) link.href = BASE_FAVICON
    return
  }
  const canvas = document.createElement("canvas")
  canvas.width = 32
  canvas.height = 32
  const ctx = canvas.getContext("2d")
  if (!ctx) return
  const img = new Image()
  img.onload = () => {
    ctx.drawImage(img, 0, 0, 32, 32)
    ctx.fillStyle = "#f85149"
    ctx.beginPath()
    ctx.arc(24, 8, 7, 0, Math.PI * 2)
    ctx.fill()
    link.href = canvas.toDataURL("image/png")
  }
  img.src = BASE_FAVICON
}
