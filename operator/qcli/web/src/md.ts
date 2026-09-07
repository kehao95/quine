// Markdown rendering for message/reply bodies: marked + highlight.js fences,
// sanitized with DOMPurify before it ever touches innerHTML (cells carry
// arbitrary model/tool text; loopback trust is no excuse for XSS-by-design).
import { Marked } from "marked"
import { markedHighlight } from "marked-highlight"
import hljs from "highlight.js/lib/common"
import DOMPurify from "dompurify"

const marked = new Marked(
  markedHighlight({
    langPrefix: "hljs language-",
    highlight(code: string, lang: string): string {
      try {
        if (lang && hljs.getLanguage(lang)) return hljs.highlight(code, { language: lang }).value
        return hljs.highlightAuto(code).value
      } catch {
        return code
      }
    },
  }),
)
marked.setOptions({ gfm: true, breaks: false })

export function renderMarkdown(text: string): string {
  const html = marked.parse(text, { async: false }) as string
  return DOMPurify.sanitize(html)
}
