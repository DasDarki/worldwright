<script setup lang="ts">
import type { Calendar, InWorldDate } from '~/types/api'
import { moonStatesForDate } from '~/composables/useCalendar'

const props = withDefaults(defineProps<{
  calendar: Calendar | null
  date?: InWorldDate | null
  size?: number
  variant?: 'inline' | 'tile' | 'minimal'
  glow?: boolean
}>(), {
  date: null,
  size: 28,
  variant: 'inline',
  glow: false,
})

const effectiveDate = computed<InWorldDate | null>(() => {
  if (props.date) return props.date
  if (!props.calendar) return null
  return {
    calendar_id: props.calendar.id,
    year: props.calendar.current_year,
    month_index: props.calendar.current_month_index,
    day: props.calendar.current_day,
  }
})

const states = computed(() => moonStatesForDate(props.calendar, effectiveDate.value))
</script>

<template>
  <div :class="['moons', variant]">
    <Transition v-for="state in states" :key="state.moon.id ?? state.moon.name" name="moon" appear>
      <div class="moon" :title="`${state.moon.name} · ${state.phase?.name || ''}`">
        <MoonGlyph :cycle="state.cycle" :size="size" :glow="glow" />
        <div v-if="variant !== 'minimal'" class="meta">
          <div class="name">{{ state.moon.name }}</div>
          <div v-if="state.phase" class="phase">
            <span v-if="state.phase.icon" class="icon">{{ state.phase.icon }}</span>
            <span>{{ state.phase.name }}</span>
          </div>
          <div v-else class="phase faint">—</div>
        </div>
      </div>
    </Transition>
    <div v-if="!states.length" class="empty">—</div>
  </div>
</template>

<style scoped lang="scss">
.moons {
  display: flex;
  align-items: center;
}
.moons.inline { gap: 18px; flex-wrap: wrap; }
.moons.tile { gap: 12px; flex-wrap: wrap; }
.moons.minimal { gap: 6px; }

.moon { display: flex; align-items: center; gap: 10px; }
.moons.minimal .moon { gap: 0; }

.meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}
.name {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
.phase {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  font-size: 13px;
  color: rgb(var(--ww-ink));
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.phase.faint { color: var(--ww-ink-faint); }
.icon { font-size: 12px; }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 13px;
}

.moon-enter-active { transition: opacity .6s ease, transform .6s cubic-bezier(.22,1,.36,1); }
.moon-enter-from { opacity: 0; transform: scale(.7) rotate(-10deg); }
</style>
