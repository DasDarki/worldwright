import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook('app:mounted', async () => {
    const auth = useAuthStore()
    if (!auth.fetched) await auth.fetchMe()
    await useOnyx().completePendingClaim()
  })
})
