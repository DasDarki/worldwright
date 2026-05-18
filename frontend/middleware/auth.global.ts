export default defineNuxtRouteMiddleware(async (to) => {
  const publicPaths = ['/login']
  if (publicPaths.includes(to.path)) return
  if (to.path.startsWith('/share/')) return

  const auth = useAuthStore()
  if (!auth.fetched) {
    await auth.fetchMe()
  }
  if (!auth.isAuthenticated) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }

  if (auth.isAdmin) {
    if (!auth.onboarding) {
      await auth.fetchOnboarding()
    }
    const needsOnboarding = auth.onboarding && !auth.onboarding.completed
    if (needsOnboarding && to.path !== '/onboarding') {
      return navigateTo('/onboarding')
    }
    if (!needsOnboarding && to.path === '/onboarding') {
      return navigateTo('/')
    }
  }
})
