<script setup lang="ts">
import type { MapAsset, MapPin } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const { t } = useI18n()
const route = useRoute()
const auth = useAuthStore()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data } = await useAsyncData(`map-${id.value}`, () =>
  $api<{ map: MapAsset }>(`/maps/${id.value}`),
)

const map = computed(() => data.value?.map || null)
useHead({ title: () => map.value?.name || t('maps.title') })
useSeoMeta({
  title: () => map.value?.name,
  ogTitle: () => map.value?.name,
  ogImage: () => map.value ? `/api/assets/${map.value.asset_id}` : '/banner.png',
  ogType: 'website',
})

const selectedPin = ref<MapPin | null>(null)
function onSelect(p: MapPin) {
  selectedPin.value = selectedPin.value?.id === p.id ? null : p
}
</script>

<template>
  <article v-if="map" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-2xl">
      <NuxtLink to="/maps" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('maps.backList') }}
      </NuxtLink>

      <header class="head mt-6 mb-8">
        <div>
          <div class="ww-eyebrow mb-4">Vol. I · {{ t('maps.title') }}</div>
          <h1 class="title"><em>{{ map.name }}</em></h1>
        </div>
        <NuxtLink v-if="auth.isAdmin" :to="`/maps/${map.id}/edit`" class="ww-btn-ghost edit">
          {{ t('maps.edit') }}
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M1 11 L3 11 L11 3 L9 1 L1 9 Z" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round"/>
          </svg>
        </NuxtLink>
      </header>

      <MapCanvas
        :asset-id="map.asset_id"
        :pins="map.pins"
        :selected-pin-id="selectedPin?.id ?? null"
        @select-pin="onSelect"
      />

      <Transition name="fade">
        <aside v-if="selectedPin" class="pin-detail ww-panel">
          <div class="pd-head">
            <h4 class="pd-title">{{ selectedPin.label || t('maps.unnamedPin') }}</h4>
            <button type="button" class="pd-close" @click="selectedPin = null" aria-label="Close">×</button>
          </div>
          <div class="pd-meta">
            <span class="ww-label visibility" :data-v="selectedPin.visibility">{{ selectedPin.visibility }}</span>
          </div>
          <div class="pd-actions">
            <NuxtLink v-if="selectedPin.target_map_id" :to="`/maps/${selectedPin.target_map_id}`" class="ww-link">
              {{ t('maps.linkedMap') }}
            </NuxtLink>
            <NuxtLink v-if="selectedPin.target_entity_slug" :to="`/entities/${selectedPin.target_entity_slug}`" class="ww-link">
              {{ t('maps.linkedEntity') }}
            </NuxtLink>
          </div>
        </aside>
      </Transition>
    </div>
  </article>
</template>

<style scoped lang="scss">
.back { margin-bottom: 14px; }
.head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 6vw, 84px);
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
  box-shadow: 0 30px 60px -25px rgb(0 0 0 / .4);
}
.pd-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}
.pd-title {
  font-family: 'Fraunces', serif;
  font-size: 20px;
  margin: 0;
}
.pd-close {
  background: transparent; border: 0;
  font-size: 20px; line-height: 1;
  color: var(--ww-ink-faint);
  cursor: pointer;
}
.pd-meta { margin-bottom: 12px; }
.visibility[data-v="secret"] { color: rgb(var(--ww-vermilion)); }
.visibility[data-v="player"] { color: rgb(var(--ww-gold-deep)); }
.visibility[data-v="public"] { color: var(--ww-ink-faint); }
.pd-actions { display: flex; gap: 14px; }

.fade-enter-active, .fade-leave-active { transition: opacity .3s ease, transform .3s cubic-bezier(.22,1,.36,1); }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(8px); }
</style>
