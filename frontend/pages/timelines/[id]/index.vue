<script setup lang="ts">
import type { Calendar, Timeline, WorldEvent } from '~/types/api'
import { formatDate } from '~/composables/useCalendar'
import { useAuthStore } from '~/stores/auth'

const { t, locale } = useI18n()
const route = useRoute()
const auth = useAuthStore()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data } = await useAsyncData(`timeline-${id.value}`, () =>
  $api<{ timeline: Timeline }>(`/timelines/${id.value}`),
)
const { data: calsData } = await useAsyncData('tl-cals', () =>
  $api<{ calendars: Calendar[] }>('/calendars'),
)

const timeline = computed(() => data.value?.timeline || null)
useHead({ title: () => timeline.value?.name || t('timelines.title') })
useSeoMeta({
  title: () => timeline.value?.name,
  description: () => timeline.value?.description,
  ogTitle: () => timeline.value?.name,
  ogDescription: () => timeline.value?.description,
})

const calsById = computed(() => {
  const m = new Map<number, Calendar>()
  for (const c of calsData.value?.calendars || []) m.set(c.id, c)
  return m
})

function dateFor(ev: WorldEvent): string {
  return formatDate(calsById.value.get(ev.calendar_id), ev, (locale.value as 'en' | 'de'))
}

useReveal()
</script>

<template>
  <article v-if="timeline" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="topbar">
        <NuxtLink to="/timelines" class="ww-btn-ghost back">
          <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
            <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('timelines.backList') }}
        </NuxtLink>
        <NuxtLink v-if="auth.isAdmin" :to="`/timelines/${timeline.id}/edit`" class="ww-btn-ghost edit">
          {{ t('timelines.edit') }}
        </NuxtLink>
      </div>

      <header class="stagger mt-8 mb-12">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I · {{ t('timelines.titleSingular') }}
        </div>
        <h1 class="title">
          <em>{{ timeline.name }}</em>
        </h1>
        <p v-if="timeline.description" class="desc">{{ timeline.description }}</p>
      </header>

      <Ornament>{{ t('timelines.chronicle') }}</Ornament>

      <ol class="timeline mt-10">
        <li v-for="ev in timeline.events || []" :key="ev.id" class="event reveal">
          <div class="date">{{ dateFor(ev) }}</div>
          <div class="rule"><span class="dot" aria-hidden="true" /></div>
          <div class="body">
            <h3 class="evt-title">{{ ev.title }}</h3>
            <div v-if="ev.importance" class="importance">
              <span v-for="n in ev.importance" :key="n" class="star">✦</span>
            </div>
          </div>
        </li>
        <li v-if="!(timeline.events || []).length" class="empty">{{ t('timelines.noEvents') }}</li>
      </ol>
    </div>
  </article>
</template>

<style scoped lang="scss">
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(44px, 7vw, 96px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.desc {
  font-style: italic;
  font-size: 20px;
  color: rgb(var(--ww-ink-shade));
  margin: 18px 0 0;
  max-width: 36em;
}

.timeline {
  list-style: none;
  padding: 0;
  margin: 0;
}
.event {
  display: grid;
  grid-template-columns: 200px 32px 1fr;
  gap: 24px;
  align-items: start;
  padding: 22px 0;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
@media (max-width: 640px) {
  .event { grid-template-columns: 1fr; }
  .rule { display: none; }
}
.date {
  font-family: 'Cormorant SC', serif;
  font-size: 13px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
  padding-top: 2px;
}
.rule { position: relative; min-height: 20px; }
.rule::before {
  content: '';
  position: absolute;
  left: 50%; top: 6px; bottom: -22px;
  width: 1px;
  background: var(--ww-ink-hairline);
}
.event:last-of-type .rule::before { bottom: 0; height: 12px; }
.dot {
  position: absolute;
  left: 50%; top: 6px;
  transform: translate(-50%, 0);
  width: 9px; height: 9px;
  border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 0 3px rgb(var(--ww-parchment)), 0 0 0 4px rgb(var(--ww-vermilion) / .35);
}
.evt-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(20px, 2.4vw, 28px);
  margin: 0;
}
.importance .star {
  color: rgb(var(--ww-gold));
  font-size: 12px;
  margin-right: 2px;
}
.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 80px 0;
  border: 0;
}
</style>
