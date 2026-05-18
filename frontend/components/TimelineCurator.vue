<script setup lang="ts">
import type { Calendar, WorldEvent } from '~/types/api'
import { formatDate } from '~/composables/useCalendar'

const props = defineProps<{
  modelValue: number[]
  events: WorldEvent[]
  calendars: Calendar[]
}>()

const emit = defineEmits<{ 'update:modelValue': [v: number[]] }>()

const { t, locale } = useI18n()

const selected = computed(() => new Set(props.modelValue))
const calsById = computed(() => {
  const m = new Map<number, Calendar>()
  for (const c of props.calendars) m.set(c.id, c)
  return m
})

const query = ref('')
const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return props.events
  return props.events.filter((e) => e.title.toLowerCase().includes(q))
})

function toggle(ev: WorldEvent) {
  const set = new Set(props.modelValue)
  if (set.has(ev.id)) set.delete(ev.id)
  else set.add(ev.id)
  emit('update:modelValue', Array.from(set))
}

function dateFor(ev: WorldEvent): string {
  return formatDate(calsById.value.get(ev.calendar_id), ev, (locale.value as 'en' | 'de'))
}

const selectedCount = computed(() => selected.value.size)
</script>

<template>
  <div class="curator">
    <div class="head">
      <div class="ww-label">{{ t('timelines.curator.title') }}</div>
      <div class="count">
        <span class="num">{{ selectedCount }}</span> · {{ t('timelines.curator.selected') }}
      </div>
    </div>

    <div class="search">
      <input
        v-model="query"
        type="search"
        class="ww-input"
        :placeholder="t('timelines.curator.search')"
      />
    </div>

    <TransitionGroup name="row" tag="ul" class="events">
      <li
        v-for="ev in filtered"
        :key="ev.id"
        :class="['event', { on: selected.has(ev.id) }]"
        @click="toggle(ev)"
      >
        <span class="check" aria-hidden="true">
          <svg viewBox="0 0 14 14" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <path d="M2 7 L6 11 L12 3" />
          </svg>
        </span>
        <span class="date">{{ dateFor(ev) }}</span>
        <span class="title">{{ ev.title }}</span>
      </li>
      <li v-if="!filtered.length" key="empty" class="empty">{{ t('timelines.curator.none') }}</li>
    </TransitionGroup>
  </div>
</template>

<style scoped lang="scss">
.curator { display: grid; gap: 12px; }
.head {
  display: flex; align-items: baseline; justify-content: space-between;
}
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.count .num {
  font-family: 'Fraunces', serif;
  font-style: italic;
  font-size: 18px;
  color: rgb(var(--ww-vermilion));
  letter-spacing: -0.02em;
  margin-right: 4px;
}

.search { display: grid; }

.events {
  list-style: none; margin: 0; padding: 0;
  display: grid; gap: 4px;
  max-height: 60vh;
  overflow-y: auto;
}
.event {
  display: grid;
  grid-template-columns: 24px 160px 1fr;
  gap: 12px;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  border: 1px solid transparent;
  background: rgb(var(--ww-parchment-deep) / .4);
  transition: background-color .25s ease, border-color .25s ease, padding .25s ease;
}
.event:hover { background: rgb(var(--ww-gold) / .12); padding-left: 16px; }
.event.on {
  background: rgb(var(--ww-gold) / .15);
  border-color: rgb(var(--ww-gold) / .5);
}
.event.on .check { background: rgb(var(--ww-gold)); color: rgb(var(--ww-ink-shade)); border-color: rgb(var(--ww-gold)); }
.check {
  width: 20px; height: 20px;
  border: 1px solid var(--ww-ink-hairline);
  display: inline-flex; align-items: center; justify-content: center;
  color: transparent;
  transition: color .25s ease, background-color .25s ease, border-color .25s ease;
}
.event.on .check svg { stroke: currentColor; }
.date {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
.title {
  font-family: 'EB Garamond', serif;
  font-size: 16px;
}
.empty {
  text-align: center;
  font-style: italic;
  color: var(--ww-ink-faint);
  padding: 24px 0;
  border: 0;
  background: transparent;
  cursor: default;
}

.row-enter-active, .row-leave-active { transition: opacity .3s ease, transform .3s cubic-bezier(.22,1,.36,1); }
.row-enter-from { opacity: 0; transform: translateY(-4px); }
.row-leave-to { opacity: 0; transform: translateX(6px); }
</style>
