<script setup lang="ts">
import type { EntitySummary, EntityType } from '~/types/api'

const { t, locale } = useI18n()
const route = useRoute()
const tagFilter = computed(() => (route.query.tag as string | undefined) || '')

useHead({ title: () => tagFilter.value ? `#${tagFilter.value} · ${t('entities.title')}` : t('entities.title') })
useSeoMeta({ description: () => t('entities.subtitle') })

const { $api } = useNuxtApp()

const { data: typesData } = await useAsyncData('entity-types', () =>
  $api<{ entity_types: EntityType[] }>('/entity-types'),
)
const { data: entitiesData } = await useAsyncData(
  () => `all-entities-${tagFilter.value}`,
  () => {
    const qs = tagFilter.value ? `?tag=${encodeURIComponent(tagFilter.value)}` : ''
    return $api<{ entities: EntitySummary[] }>(`/entities${qs}`)
  },
  { watch: [tagFilter] },
)

const grouped = computed(() => {
  const map = new Map<number, EntitySummary[]>()
  for (const e of entitiesData.value?.entities || []) {
    const arr = map.get(e.entity_type_id) || []
    arr.push(e)
    map.set(e.entity_type_id, arr)
  }
  return (typesData.value?.entity_types || [])
    .map(t => ({ type: t, entries: map.get(t.id) || [] }))
    .filter(g => g.entries.length > 0)
})

function typeLabel(t: EntityType) {
  return locale.value === 'de' ? t.name_de : t.name_en
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-2xl">
      <div class="stagger mb-16 md:mb-24">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I · The Codex
        </div>
        <h1 class="hero-title mb-6">
          <em>{{ t('entities.title') }}</em>
        </h1>
        <p class="lede">{{ t('entities.subtitle') }}</p>
        <Transition name="filter-bar">
          <div v-if="tagFilter" class="filter-bar">
            <span class="ww-label">{{ t('entities.filteredByTag') }}</span>
            <span class="tag-pill">{{ tagFilter }}</span>
            <NuxtLink to="/entities" class="clear">{{ t('entities.clearFilter') }} ×</NuxtLink>
          </div>
        </Transition>
      </div>

      <div class="space-y-20">
        <section v-for="g in grouped" :key="g.type.id" class="group reveal">
          <div class="group-head">
            <h2 class="group-title">{{ typeLabel(g.type) }}</h2>
            <span class="group-count">{{ g.entries.length }}</span>
          </div>
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-px entries">
            <NuxtLink
              v-for="e in g.entries"
              :key="e.id"
              :to="`/entities/${e.slug}`"
              class="entry"
            >
              <h3 class="entry-title">{{ e.title }}</h3>
              <p v-if="e.summary" class="entry-summary">{{ e.summary }}</p>
              <div class="entry-meta">
                <span class="ww-label visibility" :data-v="e.visibility">{{ t(`entities.visibility.${e.visibility}`) }}</span>
              </div>
            </NuxtLink>
          </div>
        </section>

        <p v-if="!grouped.length" class="empty">{{ t('entities.empty') }}</p>
      </div>
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
  font-size: 18px;
  color: var(--ww-ink-faint);
  max-width: 38em;
}

.group { position: relative; }
.group-head {
  display: flex;
  align-items: baseline;
  gap: 18px;
  margin-bottom: 26px;
}
.group-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(28px, 3.4vw, 44px);
  letter-spacing: -0.02em;
}
.group-count {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .32em;
  color: var(--ww-ink-faint);
}

.entries {
  background: var(--ww-ink-hairline);
  border-top: 1px solid var(--ww-ink-hairline);
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.entry {
  background: rgb(var(--ww-parchment));
  padding: 26px 24px 24px;
  position: relative;
  transition: background-color .5s ease, transform .35s cubic-bezier(.22,1,.36,1);
}
.entry::after {
  content: '';
  position: absolute;
  left: 24px; right: 24px; bottom: 0;
  height: 2px;
  background: rgb(var(--ww-gold));
  transform: scaleX(0);
  transform-origin: left;
  transition: transform .6s cubic-bezier(.22,1,.36,1);
}
.entry:hover { background: rgb(var(--ww-parchment-deep)); }
.entry:hover::after { transform: scaleX(1); }

.entry-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
  font-size: 22px;
  line-height: 1.15;
  margin: 0 0 10px;
}
.entry-summary {
  font-size: 15px;
  line-height: 1.5;
  color: var(--ww-ink-faint);
  margin: 0;
}
.entry-meta { margin-top: 14px; }
.visibility {
  font-size: 9px;
  color: var(--ww-ink-faint);
}
.visibility[data-v="secret"] { color: rgb(var(--ww-vermilion)); }
.visibility[data-v="player"] { color: rgb(var(--ww-gold-deep)); }
.visibility[data-v="public"] { color: rgb(var(--ww-ink) / .55); }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 80px 0;
}

.filter-bar {
  margin-top: 24px;
  display: inline-flex;
  align-items: center;
  gap: 14px;
  padding: 10px 16px;
  background: rgb(var(--ww-gold) / .12);
  border: 1px solid rgb(var(--ww-gold) / .35);
}
.tag-pill {
  font-family: 'Fraunces', serif;
  font-style: italic;
  font-size: 17px;
  color: rgb(var(--ww-vermilion));
}
.clear {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  border-bottom: 1px solid var(--ww-ink-hairline);
  padding-bottom: 1px;
  transition: color .25s ease, border-color .25s ease;
}
.clear:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); }

.filter-bar-enter-active, .filter-bar-leave-active {
  transition: opacity .3s ease, transform .35s cubic-bezier(.22,1,.36,1);
}
.filter-bar-enter-from, .filter-bar-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
