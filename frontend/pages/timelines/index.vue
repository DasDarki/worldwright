<script setup lang="ts">
import type { Timeline } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const { t } = useI18n()
useHead({ title: () => t('timelines.title') })

const auth = useAuthStore()
const { $api } = useNuxtApp()

const { data } = await useAsyncData('timelines-list', () =>
  $api<{ timelines: Timeline[] }>('/timelines'),
)
const timelines = computed(() => data.value?.timelines || [])

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-12 flex items-end justify-between gap-6">
        <div>
          <div class="ww-eyebrow mb-6 flex items-center gap-3">
            <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
            Vol. I · The threads of time
          </div>
          <h1 class="hero-title">
            <em>{{ t('timelines.title') }}</em>
          </h1>
          <p class="lede">{{ t('timelines.lede') }}</p>
        </div>
        <NuxtLink v-if="auth.isAdmin" to="/timelines/new" class="ww-btn-primary new-btn">
          {{ t('timelines.new') }}
          <span class="arrow" aria-hidden="true">
            <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
              <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
        </NuxtLink>
      </div>

      <p v-if="!timelines.length" class="empty">{{ t('timelines.empty') }}</p>

      <ul v-else class="tl-list">
        <li v-for="tl in timelines" :key="tl.id" class="tl-item reveal">
          <NuxtLink :to="`/timelines/${tl.id}`" class="link">
            <div class="left">
              <h2 class="name">{{ tl.name }}</h2>
              <p v-if="tl.description" class="desc">{{ tl.description }}</p>
            </div>
            <div class="right">
              <span class="count">{{ tl.event_ids.length }} {{ t('timelines.events') }}</span>
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

.tl-list {
  list-style: none; padding: 0; margin: 40px 0 0;
  display: grid; gap: 1px;
  background: var(--ww-ink-hairline);
  border-top: 1px solid var(--ww-ink-hairline);
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.tl-item { background: rgb(var(--ww-parchment)); transition: background-color .4s ease; }
.tl-item:hover { background: rgb(var(--ww-parchment-deep)); }
.link {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 24px;
  padding: 24px 28px;
}
.left { display: flex; flex-direction: column; gap: 6px; min-width: 0; }
.name {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(24px, 3vw, 36px);
  letter-spacing: -0.02em;
  margin: 0;
}
.desc { font-style: italic; color: var(--ww-ink-faint); margin: 0; font-size: 15px; }
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
</style>
