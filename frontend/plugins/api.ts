export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  const event = useRequestEvent()

  const baseURL = import.meta.server
    ? (config.apiBaseInternal as string)
    : (config.public.apiBase as string)

  const api = $fetch.create({
    baseURL,
    credentials: 'include',
    onRequest({ options }) {
      if (import.meta.server && event) {
        const incoming = event.node.req.headers.cookie
        if (incoming) {
          const headers = new Headers(options.headers as HeadersInit | undefined)
          headers.set('cookie', incoming)
          options.headers = headers
        }
      }
    },
    onResponse({ response }) {
      if (import.meta.server && event) {
        const sc = response.headers.get('set-cookie')
        if (sc) {
          event.node.res.appendHeader('set-cookie', sc)
        }
      }
    },
  })

  return { provide: { api } }
})
