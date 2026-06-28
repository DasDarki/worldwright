// Maps the user-facing direct-download URLs (e.g. /export.pdf.zip) onto
// the backend endpoints that actually stream the file. Doing this in
// Nitro middleware means the browser only ever sees one round-trip:
// request /export.pdf.zip → 302 → /api/export/wiki.pdf.zip → file body.
//
// We intentionally use a server-side redirect (not a Nuxt page) so the
// URL keeps working for crawlers, curl, and direct paste-into-the-bar
// scenarios without first booting Vue.

const REDIRECTS: Record<string, string> = {
  '/export.md':      '/api/export/wiki.md',
  '/export.md.zip':  '/api/export/wiki.md.zip',
  '/export.pdf':     '/api/export/wiki.pdf',
  '/export.pdf.zip': '/api/export/wiki.pdf.zip',
}

export default defineEventHandler((event) => {
  // event.path includes any query string, so split it off
  const path = (event.path || '').split('?')[0]
  const target = REDIRECTS[path]
  if (!target) return
  return sendRedirect(event, target, 302)
})
