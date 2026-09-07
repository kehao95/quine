import { For, Show, createMemo } from "solid-js"
import { TextAttributes } from "@opentui/core"
import type { RosterEntry } from "../contract"
import { shortSession } from "../contract"
import { HELP_ROWS } from "../keymap"
import { theme } from "../theme"

export function HelpOverlay(props: { width: number; height: number }) {
  return (
    <box position="absolute" left={0} top={0} width={props.width} height={props.height} alignItems="center" justifyContent="center" backgroundColor={theme.bg} zIndex={100}>
      <box backgroundColor={theme.menuBg} border borderStyle="rounded" borderColor={theme.borderActive} paddingLeft={2} paddingRight={2} paddingTop={1} paddingBottom={1}>
        <text fg={theme.accent} attributes={TextAttributes.BOLD}>
          qcli help
        </text>
        <box height={1} />
        <For each={HELP_ROWS}>{(row) => <text fg={theme.text}>{`  ${row[0].padEnd(14)} ${row[1]}`}</text>}</For>
        <box height={1} />
        <text fg={theme.dim}>{`receipt ladder: ▷ queued  ✓ received  ⚡ woke  ▣ delivered (unlit = ·)`}</text>
        <text fg={theme.dim}>press any key to close</text>
      </box>
    </box>
  )
}

function livenessGlyph(peer: RosterEntry): { glyph: string; color: string } {
  if (!peer.live) return { glyph: "◍", color: theme.dim }
  return { glyph: "●", color: theme.success }
}

// lineageIndent indents children under ancestors present in the same roster
// (parent_session chains), so fork lineage reads as a tree without a tree view.
function lineageDepth(peer: RosterEntry, bySession: Map<string, RosterEntry>): number {
  let depth = 0
  let parent = peer.parent_session
  const seen = new Set<string>()
  while (parent && bySession.has(parent) && !seen.has(parent)) {
    seen.add(parent)
    depth++
    parent = bySession.get(parent)!.parent_session
  }
  return Math.min(depth, 6)
}

export function RosterOverlay(props: {
  width: number
  height: number
  peers: RosterEntry[]
  selected: number
  filter: string
}) {
  const bySession = createMemo(() => new Map(props.peers.map((p) => [p.session, p])))
  return (
    <box position="absolute" left={0} top={0} width={props.width} height={props.height} alignItems="center" paddingTop={Math.floor(props.height / 6)} backgroundColor={theme.bg} zIndex={200}>
      <box width={Math.min(76, props.width - 4)} backgroundColor={theme.menuBg} border borderStyle="rounded" borderColor={theme.borderActive} paddingLeft={1} paddingRight={1} paddingTop={1} paddingBottom={1}>
        <box flexDirection="row" justifyContent="space-between">
          <text fg={theme.accent} attributes={TextAttributes.BOLD}>
            peers
          </text>
          <text fg={theme.dim}>{props.filter ? `filter: ${props.filter}` : "type to filter"}</text>
        </box>
        <box height={1} />
        <Show
          when={props.peers.length > 0}
          fallback={<text fg={theme.dim}>{`no peers found — /spawn <mission> from the composer, or start one: quine "mission"`}</text>}
        >
          <For each={props.peers}>
            {(peer, i) => {
              const liveness = livenessGlyph(peer)
              const indent = "  ".repeat(lineageDepth(peer, bySession()))
              return (
                <box backgroundColor={i() === props.selected ? theme.selectionBg : theme.menuBg} paddingLeft={1} paddingRight={1} flexDirection="row" gap={1}>
                  <text fg={liveness.color}>{liveness.glyph}</text>
                  <text fg={i() === props.selected ? theme.text : theme.textMuted}>{`${indent}${shortSession(peer.session)}`}</text>
                  <Show when={peer.attached}>
                    <text fg={theme.accent}>{`← you`}</text>
                  </Show>
                  <text fg={theme.dim}>{`${peer.model} · depth=${peer.depth} · pending=${peer.pending} · pid=${peer.pid}`}</text>
                </box>
              )
            }}
          </For>
        </Show>
        <box height={1} />
        <text fg={theme.dim}>↑/↓ select · enter attach · esc close · /spawn from composer</text>
      </box>
    </box>
  )
}
