<script setup lang="ts">
import type { MapAsset, MapPin } from '~/types/api'

definePageMeta({ layout: 'share' })

const { t } = useI18n()
const route = useRoute()
const { $api } = useNuxtApp()
const config = useRuntimeConfig()

const id = computed(() => Number(route.params.id))

const { data, error } = await useAsyncData(`share-map-${id.value}`, () =>
  $api<{ map: MapAsset }>(`/share/maps/${id.value}`),
)

if (error.value && (error.value as any).statusCode === 404) {
  throw createError({ statusCode: 404, statusMessage: 'Map not found' })
}

const map = computed(() => data.value?.map || null)
const siteName = computed(() => (config.public.siteName as string) || 'Worldwright')
const { absolute } = useSiteUrl()
const ogImageUrl = computed(() => map.value
  ? absolute(`/api/assets/${map.value.asset_id}`)
  : absolute('/banner.png'))
const pageUrl = computed(() => map.value ? absolute(`/share/map/${map.value.id}`) : absolute('/'))

useHead({ title: () => map.value?.name || t('maps.title') })
useSeoMeta({
  title: () => map.value?.name,
  ogTitle: () => map.value ? `${map.value.name} · ${siteName.value}` : siteName.value,
  ogDescription: () => map.value ? `A map from ${siteName.value}` : siteName.value,
  ogImage: () => ogImageUrl.value,
  ogImageAlt: () => map.value?.name || siteName.value,
  ogType: 'website',
  ogUrl: () => pageUrl.value,
  ogSiteName: () => siteName.value,
  twitterCard: 'summary_large_image',
  twitterTitle: () => map.value?.name,
  twitterImage: () => ogImageUrl.value,
})

const selectedPin = ref<MapPin | null>(null)
function onSelect(p: MapPin) {
  selectedPin.value = selectedPin.value?.id === p.id ? null : p
}
</script>

<template>
  <article v-if="map" class="share-map py-12 md:py-16">
    <div class="mx-auto max-w-screen-2xl px-6 md:px-12">
      <header class="head">
        <div>
          <div class="ww-label public-badge">{{ t('share.publicBadge') }}</div>
          <h1 class="title"><em>{{ map.name }}</em></h1>
        </div>
      </header>

      <MapCanvas
        :asset-id="map.asset_id"
        :pins="map.pins"
        :selected-pin-id="selectedPin?.id ?? null"
        @select-pin="onSelect"
      />

      <Transition name="fade">
        <aside v-if="selectedPin" class="pin-detail ww-panel">
          <h4 class="pd-title">{{ selectedPin.label || t('maps.unnamedPin') }}</h4>
          <a v-if="selectedPin.target_entity_slug"
            :href="`/share/entity/${selectedPin.target_entity_slug}`"
            class="ww-link"
          >{{ t('maps.linkedEntity') }}</a>
          <a v-if="selectedPin.target_map_id"
            :href="`/share/map/${selectedPin.target_map_id}`"
            class="ww-link"
          >{{ t('maps.linkedMap') }}</a>
        </aside>
      </Transition>
    </div>
  </article>
</template>

<style scoped lang="scss">
.head { margin-bottom: 20px; }
.public-badge {
  color: rgb(var(--ww-gold-deep));
  margin-bottom: 8px;
}
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(36px, 6vw, 80px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}

.pin-detail {
  position: fixed;
  right: 24px; bottom: 24px;
  width: min(360px, calc(100% - 48px));
  z-index: 60;
  display: grid;
  gap: 10px;
}
.pd-title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  margin: 0;
}

.fade-enter-active, .fade-leave-active { transition: opacity .3s ease, transform .3s cubic-bezier(.22,1,.36,1); }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(8px); }
</style>
