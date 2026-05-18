<script setup lang="ts">
import type { TagWithCount } from '~/types/api'

const { t } = useI18n()
useHead({ title: () => t('tags.title') })

const { $api } = useNuxtApp()
const { data } = await useAsyncData('tags-cloud', () =>
  $api<{ tags: TagWithCount[] }>('/tags?with_counts=true'),
)

const tags = computed(() => data.value?.tags || [])
const maxCount = computed(() => tags.value.reduce((m, t) => Math.max(m, t.count), 1))

function weight(count: number) {
  if (maxCount.value === 0) return 1
  return 0.55 + (count / maxCount.value) * 0.85
}

const query = ref('')
const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return tags.value
  return tags.value.filter((t) => t.name.includes(q))
})

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-12">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I · The marginalia
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('tags.title') }}</em>
        </h1>
        <p class="lede">{{ t('tags.lede') }}</p>
      </div>

      <div class="filter">
        <input v-model="query" type="search" class="ww-input" :placeholder="t('tags.filter')" />
      </div>

      <p v-if="!tags.length" class="empty">{{ t('tags.empty') }}</p>

      <TransitionGroup v-else name="cloud" tag="ul" class="cloud reveal">
        <li
          v-for="tag in filtered"
          :key="tag.name"
          class="tag-item"
          :style="{ fontSize: `${weight(tag.count)}rem` }"
        >
          <NuxtLink :to="{ path: '/entities', query: { tag: tag.name } }" class="tag-link">
            <span class="name">{{ tag.name }}</span>
            <span class="count">{{ tag.count }}</span>
          </NuxtLink>
        </li>
      </TransitionGroup>
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
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 18px; }

.filter { margin: 20px 0 36px; max-width: 28em; }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 80px 0;
}

.cloud {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 14px 18px;
  align-items: baseline;
}
.tag-item {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 80, "opsz" 24, "wght" 400;
}
.tag-link {
  display: inline-flex;
  align-items: baseline;
  gap: 6px;
  color: rgb(var(--ww-ink));
  border-bottom: 1px solid transparent;
  transition: color .25s ease, border-color .25s ease, transform .35s cubic-bezier(.22,1,.36,1);
}
.tag-link:hover {
  color: rgb(var(--ww-vermilion));
  border-color: rgb(var(--ww-vermilion));
  transform: translateY(-2px);
}
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 9px;
  letter-spacing: .22em;
  color: var(--ww-ink-faint);
  vertical-align: super;
}

.cloud-enter-active, .cloud-leave-active { transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1); }
.cloud-enter-from { opacity: 0; transform: translateY(-4px) scale(.94); }
.cloud-leave-to { opacity: 0; transform: scale(.92); }
</style>
