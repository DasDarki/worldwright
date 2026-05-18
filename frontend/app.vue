<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useSiteUrl } from '~/composables/useSiteUrl'

const auth = useAuthStore()
if (!auth.providers.oauth.length) {
  await auth.fetchProviders()
}

const config = useRuntimeConfig()
const { absolute } = useSiteUrl()
const siteName = computed(() => (config.public.siteName as string) || 'Worldwright')
const siteDescription = computed(() => (config.public.siteDescription as string) || '')
const bannerUrl = computed(() => absolute('/banner.png'))

useHead({
  htmlAttrs: { lang: 'en' },
  titleTemplate: (chunk?: string) => chunk ? `${chunk} · ${siteName.value}` : siteName.value,
})

useSeoMeta({
  description: () => siteDescription.value,
  ogSiteName: () => siteName.value,
  ogTitle: () => siteName.value,
  ogDescription: () => siteDescription.value,
  ogType: 'website',
  ogImage: () => bannerUrl.value,
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: () => `${siteName.value} — an atlas of every world you will ever invent.`,
  twitterCard: 'summary_large_image',
  twitterTitle: () => siteName.value,
  twitterDescription: () => siteDescription.value,
  twitterImage: () => bannerUrl.value,
  themeColor: '#173842',
})
</script>

<template>
  <PaperOverlays />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <ToastContainer />
  <ShortcutHelp />
  <PWAStatus />
</template>
