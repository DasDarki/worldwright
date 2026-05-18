<script setup lang="ts">
import type { Calendar } from '~/types/api'

const { t } = useI18n()
useHead({ title: () => t('calendars.title') })

const { $api } = useNuxtApp()
const { data } = await useAsyncData('calendars-list', () =>
  $api<{ calendars: Calendar[] }>('/calendars'),
)
const calendars = computed(() => data.value?.calendars || [])

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-16">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I · The hours
        </div>
        <h1 class="hero-title mb-6">
          <em>{{ t('calendars.title') }}</em>
        </h1>
        <p class="lede">{{ t('calendars.lede') }}</p>
      </div>

      <ul class="cal-list">
        <li v-for="c in calendars" :key="c.id" class="cal-item reveal">
          <NuxtLink :to="`/calendars/${c.id}`" class="cal-link">
            <div class="left">
              <div class="ww-label">{{ c.epoch_name || t('calendars.epochless') }}</div>
              <h2 class="name">{{ c.name }}</h2>
            </div>
            <div class="right">
              <span class="meta">{{ c.current_year }}</span>
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
.lede {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 18px;
  max-width: 38em;
}

.cal-list { list-style: none; margin: 0; padding: 0; display: grid; gap: 1px; background: var(--ww-ink-hairline); border-top: 1px solid var(--ww-ink-hairline); border-bottom: 1px solid var(--ww-ink-hairline); }
.cal-item { background: rgb(var(--ww-parchment)); transition: background-color .4s ease; }
.cal-item:hover { background: rgb(var(--ww-parchment-deep)); }
.cal-link {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 30px;
}
.left { display: flex; flex-direction: column; gap: 6px; }
.name {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(28px, 3.4vw, 44px);
  letter-spacing: -0.02em;
  margin: 0;
}
.right .meta {
  font-family: 'Cormorant SC', serif;
  font-size: 14px;
  letter-spacing: .32em;
  color: rgb(var(--ww-gold-deep));
}
</style>
