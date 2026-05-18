<script setup lang="ts">
import type { MapAsset } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const { t } = useI18n()
useHead({ title: () => t('maps.title') })

const auth = useAuthStore()
const { $api } = useNuxtApp()
const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)

const { data } = await useAsyncData('maps-list', () =>
  $api<{ maps: MapAsset[] }>('/maps'),
)
const maps = computed(() => data.value?.maps || [])

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-2xl">
      <div class="stagger mb-12 flex items-end justify-between gap-6">
        <div>
          <div class="ww-eyebrow mb-6 flex items-center gap-3">
            <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
            Vol. I · The atlas proper
          </div>
          <h1 class="hero-title">
            <em>{{ t('maps.title') }}</em>
          </h1>
          <p class="lede">{{ t('maps.lede') }}</p>
        </div>
        <NuxtLink v-if="auth.isAdmin" to="/maps/new" class="ww-btn-primary new-btn">
          {{ t('maps.new') }}
          <span class="arrow" aria-hidden="true">
            <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
              <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
        </NuxtLink>
      </div>

      <p v-if="!maps.length" class="empty">{{ t('maps.empty') }}</p>

      <ul v-else class="grid">
        <li v-for="m in maps" :key="m.id" class="card reveal">
          <NuxtLink :to="`/maps/${m.id}`" class="link">
            <div class="thumb">
              <img :src="`${apiBase}/assets/${m.asset_id}`" :alt="m.name" />
              <span class="overlay" />
            </div>
            <div class="meta">
              <h3 class="name">{{ m.name }}</h3>
              <span class="count">{{ m.pins?.length || 0 }} {{ t('maps.pins') }}</span>
            </div>
          </NuxtLink>
        </li>
      </ul>
    </div>
  </section>
</template>

<style scoped lang="scss">
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(48px, 8vw, 120px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 18px; margin-top: 6px; }

.new-btn { white-space: nowrap; }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 80px 0;
}

.grid {
  list-style: none;
  padding: 0;
  margin: 40px 0 0;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 28px;
}
.card { transition: transform .35s cubic-bezier(.22,1,.36,1); }
.card:hover { transform: translateY(-3px); }

.link { display: block; }

.thumb {
  position: relative;
  aspect-ratio: 4 / 3;
  overflow: hidden;
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment-deep) / .5);
}
.thumb img {
  width: 100%; height: 100%;
  object-fit: cover;
  transition: transform .8s cubic-bezier(.22,1,.36,1), filter .8s ease;
  filter: sepia(.18);
}
.card:hover .thumb img { transform: scale(1.04); filter: sepia(.05); }

.overlay {
  position: absolute;
  inset: 0;
  pointer-events: none;
  box-shadow: inset 0 0 60px rgb(14 36 44 / .35);
}

.meta {
  margin-top: 12px;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}
.name {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
  font-size: 20px;
  margin: 0;
}
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
</style>
