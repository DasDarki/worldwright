// Walks text nodes within a root element and wraps occurrences of `query`
// in <mark> elements. Returns the first highlight element so the caller can
// scrollIntoView. Use only on the client.
//
// Reasoning: the alternative (CSS Custom Highlight API) has uneven browser
// support and can't be scrolled to with scrollIntoView; manual DOM wrapping
// works everywhere and lets us style the highlight via the codex palette.

export function highlightInDom(root: HTMLElement | null, query: string, className = 'search-hl'): HTMLElement | null {
  if (!root || !query) return null
  const trimmed = query.trim()
  if (!trimmed) return null

  // Pull out the individual word-like tokens; FTS matches whole tokens, so
  // we follow the same idea instead of literal-substring matching.
  const tokens = trimmed
    .split(/\s+/)
    .map((t) => t.replace(/[^\p{L}\p{N}]/gu, ''))
    .filter((t) => t.length >= 2)
  if (!tokens.length) return null

  const pattern = new RegExp(`(${tokens.map(escape).join('|')})`, 'giu')

  // Collect text nodes upfront so DOM mutations don't invalidate the walk.
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
    acceptNode(node) {
      const text = node.nodeValue || ''
      if (!text.trim()) return NodeFilter.FILTER_REJECT
      const parent = (node.parentElement as HTMLElement | null)
      // Skip inside our own marks, code blocks, etc.
      if (!parent) return NodeFilter.FILTER_REJECT
      const tag = parent.tagName
      if (tag === 'MARK' || tag === 'SCRIPT' || tag === 'STYLE') return NodeFilter.FILTER_REJECT
      return NodeFilter.FILTER_ACCEPT
    },
  })
  const texts: Text[] = []
  let n: Node | null = walker.nextNode()
  while (n) {
    texts.push(n as Text)
    n = walker.nextNode()
  }

  let first: HTMLElement | null = null
  for (const t of texts) {
    const value = t.nodeValue || ''
    pattern.lastIndex = 0
    if (!pattern.test(value)) continue
    pattern.lastIndex = 0
    const frag = document.createDocumentFragment()
    let last = 0
    for (const m of value.matchAll(pattern)) {
      const idx = m.index ?? 0
      if (idx > last) frag.appendChild(document.createTextNode(value.slice(last, idx)))
      const mark = document.createElement('mark')
      mark.className = className
      mark.textContent = m[0]
      frag.appendChild(mark)
      if (!first) first = mark
      last = idx + m[0].length
    }
    if (last < value.length) frag.appendChild(document.createTextNode(value.slice(last)))
    t.parentNode?.replaceChild(frag, t)
  }
  return first
}

function escape(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
