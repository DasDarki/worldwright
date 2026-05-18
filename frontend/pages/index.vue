<script setup lang="ts">
import type { EntitySummary, EntityType } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const { t } = useI18n()
useHead({ title: () => t('app.name') })

const auth = useAuthStore()
const { $api } = useNuxtApp()

const { data: entitiesData } = await useAsyncData('home-entities', () =>
  $api<{ entities: EntitySummary[] }>('/entities')
)
const { data: typesData } = await useAsyncData('home-types', () =>
  $api<{ entity_types: EntityType[] }>('/entity-types')
)
const { data: tagsData } = await useAsyncData('home-tags', () =>
  $api<{ tags: string[] }>('/tags')
)
const { data: recentData } = await useAsyncData('home-recent', () =>
  $api<{ entities: EntitySummary[] }>('/recent?limit=8'),
)

const greeting = computed(() => auth.user?.display_name || auth.user?.email?.split('@')[0] || '')

function typeNameFor(e: { entity_type_id: number }, types?: { id: number; name_en: string }[]): string {
  return types?.find(t => t.id === e.entity_type_id)?.name_en || ''
}

useReveal()
</script>

<template>
  <section class="py-16 md:py-24">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger">
        <div class="ww-eyebrow mb-8 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I — {{ t('app.name') }}
        </div>
        <h1 class="hero-title mb-6">
          {{ t('home.welcomePrefix') }} <em>{{ greeting }}</em>,<br />
          {{ t('home.welcomeSuffix') }}
        </h1>
        <p class="lede mb-16">{{ t('home.lede') }}</p>
      </div>

      <div class="grid md:grid-cols-3 gap-px stats reveal">
        <div class="stat">
          <div class="num">{{ entitiesData?.entities.length ?? 0 }}</div>
          <div class="ww-label">{{ t('home.stats.entries') }}</div>
        </div>
        <div class="stat">
          <div class="num">{{ typesData?.entity_types.length ?? 0 }}</div>
          <div class="ww-label">{{ t('home.stats.types') }}</div>
        </div>
        <div class="stat">
          <div class="num">{{ tagsData?.tags.length ?? 0 }}</div>
          <div class="ww-label">{{ t('home.stats.tags') }}</div>
        </div>
      </div>

      <div class="mt-20 reveal">
        <Ornament>{{ t('home.recent') }}</Ornament>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-px mt-10 entries">
          <NuxtLink
            v-for="e in (recentData?.entities || []).slice(0, 6)"
            :key="e.id"
            :to="`/entities/${e.slug}`"
            class="entry"
          >
            <div class="ww-label mb-2 entry-type">{{ typeNameFor(e, typesData?.entity_types) }}</div>
            <h3 class="entry-title">{{ e.title }}</h3>
            <p v-if="e.summary" class="entry-summary">{{ e.summary }}</p>
          </NuxtLink>
        </div>

        <div class="mt-10 flex justify-center">
          <NuxtLink to="/entities" class="ww-btn-primary">
            {{ t('nav.entities') }}
            <span class="arrow" aria-hidden="true">
              <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
                <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 6vw, 88px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.lede {
  font-size: 20px;
  color: var(--ww-ink-faint);
  max-width: 36em;
  font-style: italic;
}

.stats {
  background: var(--ww-ink-hairline);
  border-top: 1px solid var(--ww-ink-hairline);
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.stat {
  background: rgb(var(--ww-parchment));
  padding: 36px 30px;
  transition: background-color .5s ease;
}
.stat:hover { background: rgb(var(--ww-parchment-deep)); }
.stat .num {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: 64px;
  line-height: 1;
  color: rgb(var(--ww-gold-deep));
  letter-spacing: -0.03em;
  margin-bottom: 12px;
}

.entries {
  background: var(--ww-ink-hairline);
  border-top: 1px solid var(--ww-ink-hairline);
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.entry {
  background: rgb(var(--ww-parchment));
  padding: 30px 26px 28px;
  transition: background-color .5s ease, transform .35s cubic-bezier(.22,1,.36,1);
  position: relative;
}
.entry::after {
  content: '';
  position: absolute;
  left: 26px; right: 26px; bottom: 0;
  height: 2px;
  background: rgb(var(--ww-gold));
  transform: scaleX(0);
  transform-origin: left;
  transition: transform .6s cubic-bezier(.22,1,.36,1);
}
.entry:hover { background: rgb(var(--ww-parchment-deep)); }
.entry:hover::after { transform: scaleX(1); }

.entry-type { color: rgb(var(--ww-vermilion)); }

.entry-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
  font-size: 22px;
  line-height: 1.15;
  margin: 0 0 10px;
}
.entry-summary {
  font-size: 15px;
  color: var(--ww-ink-faint);
  line-height: 1.5;
}
</style>
