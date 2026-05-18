<script setup lang="ts">
import type { Calendar, WorldEvent } from '~/types/api'
import { formatDate } from '~/composables/useCalendar'
import { useAuthStore } from '~/stores/auth'

const { t, locale } = useI18n()
useHead({ title: () => t('events.title') })

const auth = useAuthStore()
const { $api } = useNuxtApp()

const { data: eventsData } = await useAsyncData('events-list', () =>
  $api<{ events: WorldEvent[] }>('/events'),
)
const { data: calendarsData } = await useAsyncData('events-list-cals', async () => {
  const list = await $api<{ calendars: Calendar[] }>('/calendars')
  const full = await Promise.all(
    list.calendars.map((c) => $api<{ calendar: Calendar }>(`/calendars/${c.id}`).then((r) => r.calendar)),
  )
  return { calendars: full }
})

const events = computed(() => eventsData.value?.events || [])
const calendarsById = computed(() => {
  const m = new Map<number, Calendar>()
  for (const c of calendarsData.value?.calendars || []) m.set(c.id, c)
  return m
})

function dateFor(ev: WorldEvent) {
  return formatDate(calendarsById.value.get(ev.calendar_id), ev, (locale.value as 'en' | 'de'))
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-12 flex items-end justify-between gap-6">
        <div>
          <div class="ww-eyebrow mb-6 flex items-center gap-3">
            <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
            Vol. I · The chronicle
          </div>
          <h1 class="hero-title">
            <em>{{ t('events.title') }}</em>
          </h1>
          <p class="lede">{{ t('events.lede') }}</p>
        </div>
        <NuxtLink v-if="auth.isAdmin" to="/events/new" class="ww-btn-primary new-btn">
          {{ t('events.new') }}
          <span class="arrow" aria-hidden="true">
            <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
              <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
        </NuxtLink>
      </div>

      <p v-if="!events.length" class="empty">{{ t('events.empty') }}</p>

      <ol v-else class="timeline">
        <li v-for="ev in events" :key="ev.id" class="event reveal">
          <div class="date-cell">
            <div class="date">{{ dateFor(ev) }}</div>
            <MoonDisplay
              v-if="calendarsById.get(ev.calendar_id)"
              :calendar="calendarsById.get(ev.calendar_id)!"
              :date="{ calendar_id: ev.calendar_id, era_id: ev.era_id, year: ev.year, month_index: ev.month_index, day: ev.day }"
              :size="18"
              variant="minimal"
              class="row-moons"
            />
          </div>
          <div class="rule"><span class="dot" aria-hidden="true" /></div>
          <NuxtLink v-if="auth.isAdmin" :to="`/events/${ev.id}/edit`" class="body-link">
            <h3 class="title">{{ ev.title }}</h3>
            <div class="meta">
              <span class="ww-label visibility" :data-v="ev.visibility">{{ ev.visibility }}</span>
              <span v-if="ev.importance" class="importance">
                <span v-for="n in ev.importance" :key="n" class="star">✦</span>
              </span>
            </div>
          </NuxtLink>
          <div v-else class="body-link static">
            <h3 class="title">{{ ev.title }}</h3>
          </div>
        </li>
      </ol>
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

.timeline {
  list-style: none;
  margin: 40px 0 0;
  padding: 0;
}
.event {
  display: grid;
  grid-template-columns: 220px 32px 1fr;
  gap: 24px;
  align-items: start;
  padding: 22px 0;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.date-cell {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.row-moons { padding-top: 2px; }
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
.rule { position: relative; height: 100%; min-height: 20px; }
.rule::before {
  content: '';
  position: absolute;
  left: 50%; top: 6px; bottom: -22px;
  width: 1px;
  background: var(--ww-ink-hairline);
}
.event:last-child .rule::before { bottom: 0; height: 12px; }
.dot {
  position: absolute;
  left: 50%; top: 6px;
  transform: translate(-50%, 0);
  width: 9px; height: 9px;
  border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 0 3px rgb(var(--ww-parchment)), 0 0 0 4px rgb(var(--ww-vermilion) / .35);
}
.body-link {
  display: block;
  transition: transform .35s cubic-bezier(.22,1,.36,1);
}
.body-link:hover { transform: translateX(4px); }
.body-link.static:hover { transform: none; }
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(22px, 2.6vw, 32px);
  letter-spacing: -0.02em;
  margin: 0 0 6px;
}
.body-link:hover .title { color: rgb(var(--ww-gold-deep)); }
.meta {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 4px;
}
.visibility[data-v="secret"] { color: rgb(var(--ww-vermilion)); }
.visibility[data-v="player"] { color: rgb(var(--ww-gold-deep)); }
.visibility[data-v="public"] { color: var(--ww-ink-faint); }
.importance .star {
  color: rgb(var(--ww-gold));
  font-size: 12px;
}
</style>
