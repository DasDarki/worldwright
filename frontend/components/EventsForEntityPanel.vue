<script setup lang="ts">
import type { Calendar, WorldEvent } from '~/types/api'
import { formatDate } from '~/composables/useCalendar'

const props = defineProps<{ entityId: number }>()

const { t, locale } = useI18n()
const { $api } = useNuxtApp()

const { data } = await useAsyncData(
  () => `entity-events-${props.entityId}`,
  () => $api<{ events: WorldEvent[] }>(`/entities/${props.entityId}/events`),
  { watch: [() => props.entityId] },
)
const { data: calData } = await useAsyncData('panel-cals', () =>
  $api<{ calendars: Calendar[] }>('/calendars'),
)

const events = computed(() => data.value?.events || [])
const calsById = computed(() => {
  const m = new Map<number, Calendar>()
  for (const c of calData.value?.calendars || []) m.set(c.id, c)
  return m
})
</script>

<template>
  <div v-if="events.length" class="ww-panel">
    <h4 class="ww-label mb-4">{{ t('events.panelTitle') }}</h4>
    <ul class="events">
      <li v-for="ev in events" :key="ev.id" class="event">
        <span class="date">{{ formatDate(calsById.get(ev.calendar_id), ev, (locale as 'en' | 'de')) }}</span>
        <span class="title">{{ ev.title }}</span>
      </li>
    </ul>
  </div>
</template>

<style scoped lang="scss">
.events { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.event { display: grid; grid-template-columns: max-content 1fr; gap: 10px 14px; align-items: baseline; }
.date {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
.title {
  font-family: 'EB Garamond', serif;
  font-size: 15px;
}
</style>
