export function useSiteUrl() {
  const config = useRuntimeConfig()
  const event = useRequestEvent()

  function origin(): string {
    const configured = (config.public.siteUrl as string) || ''
    if (configured) return configured.replace(/\/+$/, '')
    if (import.meta.server && event) {
      const host = (event.node.req.headers['x-forwarded-host'] as string)
        || (event.node.req.headers['host'] as string)
        || ''
      const proto = (event.node.req.headers['x-forwarded-proto'] as string)
        || ((event.node.req.socket as any)?.encrypted ? 'https' : 'http')
      if (host) return `${proto}://${host}`
    }
    if (import.meta.client && typeof window !== 'undefined') {
      return window.location.origin
    }
    return ''
  }

  function absolute(path: string): string {
    if (!path) return path
    if (/^https?:\/\//i.test(path)) return path
    const base = origin()
    if (!base) return path
    if (path.startsWith('/')) return base + path
    return `${base}/${path}`
  }

  return { origin, absolute }
}
