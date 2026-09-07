import type { ControlAction } from "./contract"

export const ACTIONS: ControlAction[] = ["post", "poke", "inject", "interrupt"]

// COMPOSER_KEY_BINDINGS overrides OpenTUI's Textarea defaults (which bind
// return→newline and only meta+return→submit — TUI.md §6 post-mortem).
// Plain Enter submits; Alt+Enter is the portable newline; Shift+Enter also
// maps to newline for kitty-protocol terminals that can distinguish it.
export const COMPOSER_KEY_BINDINGS = [
  { name: "return", action: "submit" },
  { name: "kpenter", action: "submit" },
  { name: "linefeed", action: "submit" },
  { name: "return", shift: true, action: "newline" },
  { name: "return", meta: true, action: "newline" },
] as const

export const HELP_ROWS: Array<[string, string]> = [
  ["Enter", "send composer text"],
  ["Alt+Enter", "newline (Shift+Enter on kitty-protocol terminals)"],
  ["Tab", "cycle delivery action"],
  ["Esc Esc", "interrupt (first Esc arms, any other key cancels)"],
  ["Ctrl+Up/Down", "focus transcript cells"],
  ["Up/Down", "history at composer edge; cell focus nav when focused"],
  ["Right/Left", "expand / collapse the focused cell"],
  ["PgUp/PgDn", "scroll transcript"],
  ["g/G", "top / bottom (empty composer)"],
  ["Ctrl+P", "peer roster (/switch)"],
  ["Ctrl+E", "edit composer in $EDITOR"],
  ["?", "help (empty composer)"],
  ["Ctrl+C Ctrl+C", "quit"],
]

export function nextAction(action: ControlAction): ControlAction {
  const i = ACTIONS.indexOf(action)
  return ACTIONS[(i + 1) % ACTIONS.length]!
}
