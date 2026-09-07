// Semantic theme tokens, adopted from OpenCode's role grouping at qcli scale.
// Components read role names, never raw hex. OpenTUI accepts hex strings for
// color props and quantizes to 256/16 colors itself on lesser terminals — the
// palette's job is semantic consistency, degradation is the renderer's.
// QCLI_THEME=dark|light selects the palette (env-only config, UX.md §4).
import { SyntaxStyle } from "@opentui/core"

export interface Palette {
  text: string
  textMuted: string
  dim: string
  bg: string
  panelBg: string
  elementBg: string
  menuBg: string
  primary: string
  accent: string
  secondary: string
  headerBg: string
  headerFg: string
  success: string
  warning: string
  error: string
  info: string
  border: string
  borderActive: string
  borderSubtle: string
  selectionBg: string
  syntaxComment: string
  syntaxKeyword: string
  syntaxString: string
  syntaxNumber: string
  syntaxFunction: string
  syntaxType: string
  syntaxVariable: string
  syntaxOperator: string
  syntaxPunctuation: string
  markdownHeading: string
  markdownLink: string
}

const dark: Palette = {
  text: "#d7dde5",
  textMuted: "#8b97a7",
  dim: "#6b7787",
  bg: "#0e1116",
  panelBg: "#11161d",
  elementBg: "#161d27",
  menuBg: "#1b2430",
  primary: "#1f6feb",
  accent: "#58a6ff",
  secondary: "#bc8cff",
  headerBg: "#1f6feb",
  headerFg: "#ffffff",
  success: "#3fb950",
  warning: "#d29922",
  error: "#f85149",
  info: "#58a6ff",
  border: "#2b3a4f",
  borderActive: "#3b5170",
  borderSubtle: "#1c2733",
  selectionBg: "#1b2533",
  syntaxComment: "#6b7787",
  syntaxKeyword: "#ff7b72",
  syntaxString: "#a5d6ff",
  syntaxNumber: "#79c0ff",
  syntaxFunction: "#d2a8ff",
  syntaxType: "#7ee787",
  syntaxVariable: "#ffa657",
  syntaxOperator: "#ff7b72",
  syntaxPunctuation: "#8b97a7",
  markdownHeading: "#79c0ff",
  markdownLink: "#58a6ff",
}

const light: Palette = {
  text: "#1f2328",
  textMuted: "#57606a",
  dim: "#77808a",
  bg: "#f6f8fa",
  panelBg: "#eaeef2",
  elementBg: "#e1e7ed",
  menuBg: "#ffffff",
  primary: "#0969da",
  accent: "#0969da",
  secondary: "#8250df",
  headerBg: "#0969da",
  headerFg: "#ffffff",
  success: "#1a7f37",
  warning: "#9a6700",
  error: "#cf222e",
  info: "#0969da",
  border: "#d0d7de",
  borderActive: "#8c959f",
  borderSubtle: "#e4e9ee",
  selectionBg: "#ddf4ff",
  syntaxComment: "#6e7781",
  syntaxKeyword: "#cf222e",
  syntaxString: "#0a3069",
  syntaxNumber: "#0550ae",
  syntaxFunction: "#8250df",
  syntaxType: "#1a7f37",
  syntaxVariable: "#953800",
  syntaxOperator: "#cf222e",
  syntaxPunctuation: "#57606a",
  markdownHeading: "#0550ae",
  markdownLink: "#0969da",
}

export const theme: Palette = process.env.QCLI_THEME === "light" ? light : dark

// getSyntaxRules maps tree-sitter / markdown scopes onto the semantic tokens
// (trimmed from OpenCode's getSyntaxRules). One palette drives markdown
// structure (heading/bold/italic/link/list) + fenced code.
export function getSyntaxRules(): Array<{
  scope: string[]
  style: { foreground?: string; bold?: boolean; italic?: boolean; underline?: boolean }
}> {
  return [
    { scope: ["comment"], style: { foreground: theme.syntaxComment, italic: true } },
    { scope: ["string", "markup.raw", "markup.inline.raw"], style: { foreground: theme.syntaxString } },
    { scope: ["number", "constant", "constant.numeric"], style: { foreground: theme.syntaxNumber } },
    { scope: ["keyword", "storage"], style: { foreground: theme.syntaxKeyword } },
    { scope: ["function", "entity.name.function", "support.function"], style: { foreground: theme.syntaxFunction } },
    { scope: ["type", "entity.name.type", "support.type"], style: { foreground: theme.syntaxType } },
    { scope: ["variable", "variable.parameter"], style: { foreground: theme.syntaxVariable } },
    { scope: ["operator", "keyword.operator"], style: { foreground: theme.syntaxOperator } },
    { scope: ["punctuation"], style: { foreground: theme.syntaxPunctuation } },
    { scope: ["markup.heading"], style: { foreground: theme.markdownHeading, bold: true } },
    { scope: ["markup.bold"], style: { foreground: theme.text, bold: true } },
    { scope: ["markup.italic"], style: { foreground: theme.text, italic: true } },
    { scope: ["markup.link", "markup.underline.link"], style: { foreground: theme.markdownLink, underline: true } },
    { scope: ["markup.list"], style: { foreground: theme.accent } },
    { scope: ["markup.quote"], style: { foreground: theme.textMuted, italic: true } },
  ]
}

let syntaxSingleton: SyntaxStyle | null = null

// syntaxStyle is the shared SyntaxStyle instance driving <markdown> and <code>.
export function syntaxStyle(): SyntaxStyle {
  if (!syntaxSingleton) syntaxSingleton = SyntaxStyle.fromTheme(getSyntaxRules())
  return syntaxSingleton
}

// tint linearly blends two hex colors — derive faded/hover shades from the
// base palette instead of minting new hex.
export function tint(a: string, b: string, alpha: number): string {
  const pa = hexToRgb(a)
  const pb = hexToRgb(b)
  const t = Math.max(0, Math.min(1, alpha))
  const mix = (x: number, y: number) => Math.round(x + (y - x) * t)
  return rgbToHex(mix(pa.r, pb.r), mix(pa.g, pb.g), mix(pa.b, pb.b))
}

// authorColor maps a peer session id to a stable hue. The Human is the accent.
export function authorColor(author?: string | null): string {
  if (!author || author === "human") return theme.accent
  let h = 0
  for (let i = 0; i < author.length; i++) h = (h * 31 + author.charCodeAt(i)) % 360
  return hslToHex(h, 55, process.env.QCLI_THEME === "light" ? 40 : 68)
}

export function truncate(s: string, n: number): string {
  return s.length > n ? s.slice(0, Math.max(1, n - 1)) + "…" : s
}

function hexToRgb(h: string): { r: number; g: number; b: number } {
  const n = parseInt(h.replace("#", ""), 16)
  return { r: (n >> 16) & 255, g: (n >> 8) & 255, b: n & 255 }
}

function rgbToHex(r: number, g: number, b: number): string {
  return "#" + [r, g, b].map((x) => x.toString(16).padStart(2, "0")).join("")
}

function hslToHex(h: number, s: number, l: number): string {
  s /= 100
  l /= 100
  const k = (n: number) => (n + h / 30) % 12
  const a = s * Math.min(l, 1 - l)
  const f = (n: number) => {
    const color = l - a * Math.max(-1, Math.min(k(n) - 3, Math.min(9 - k(n), 1)))
    return Math.round(255 * color)
      .toString(16)
      .padStart(2, "0")
  }
  return `#${f(0)}${f(8)}${f(4)}`
}
