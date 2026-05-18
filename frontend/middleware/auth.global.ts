// Routes that require authentication. Everything else is openly viewable —
// the backend filters content by visibility (public / player / secret), so an
// anonymous visitor lands on the same URL as a logged-in player and just sees
// the public subset. Edit / new pages additionally use the `admin-only`
// middleware which redirects non-admins to the corresponding view page (or to
// /login if they are not signed in at all).
const AUTH_REQUIRED_PREFIXES = ['/account', '/onboarding']

export default defineNuxtRouteMiddleware(async (to) => {
  if (to.path === '/login') return
  if (to.path.startsWith('/share/')) return

  const auth = useAuthStore()
  if (!auth.fetched) {
    await auth.fetchMe()
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

  const needsAuth = AUTH_REQUIRED_PREFIXES.some((p) => to.path === p || to.path.startsWith(`${p}/`))
  if (needsAuth && !auth.isAuthenticated) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
})
