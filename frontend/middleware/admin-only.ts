import { useAuthStore } from '~/stores/auth'
import { useToastsStore } from '~/stores/toasts'

function fallbackFor(path: string): string {
  if (/^\/entities\/new\/?$/.test(path)) return '/entities'
  const editEntity = path.match(/^\/entities\/([^/]+)\/edit\/?$/)
  if (editEntity) return `/entities/${editEntity[1]}`

  if (/^\/events\/new\/?$/.test(path)) return '/events'
  if (/^\/events\/\d+\/edit\/?$/.test(path)) return '/events'

  if (/^\/maps\/new\/?$/.test(path)) return '/maps'
  const editMap = path.match(/^\/maps\/(\d+)\/edit\/?$/)
  if (editMap) return `/maps/${editMap[1]}`

  if (/^\/timelines\/new\/?$/.test(path)) return '/timelines'
  const editTimeline = path.match(/^\/timelines\/(\d+)\/edit\/?$/)
  if (editTimeline) return `/timelines/${editTimeline[1]}`

  return '/'
}

export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()
  if (!auth.fetched) await auth.fetchMe()
  if (auth.isAdmin) return

  if (!auth.isAuthenticated) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }

  if (import.meta.client) {
    try {
      const toasts = useToastsStore()
      const { $i18n } = useNuxtApp() as any
      const message = $i18n?.t?.('auth.needAdmin') || 'Editing requires admin role.'
      toasts.error(message)
    } catch {
      /* noop */
    }
  }
  return navigateTo(fallbackFor(to.path), { replace: true })
})
