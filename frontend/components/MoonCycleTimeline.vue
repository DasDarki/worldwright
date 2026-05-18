<script setup lang="ts">
import type { CalendarMoon, CalendarMoonPhase } from '~/types/api'
import { resolveMoonPhase } from '~/composables/useCalendar'

const props = defineProps<{
  moon: CalendarMoon
  currentCycle?: number | null
}>()

const activePhase = computed(() => {
  if (props.currentCycle == null) return null
  return resolveMoonPhase(props.moon, props.currentCycle)
})

interface WindowSegment {
  phase: CalendarMoonPhase
  left: number
  width: number
}

const windows = computed<WindowSegment[]>(() => {
  const out: WindowSegment[] = []
  for (const p of props.moon.phases || []) {
    const r = Math.max(0, p.randomness)
    if (r === 0) continue
    const start = p.cycle_position - r
    const end = p.cycle_position + r
    if (start < 0) {
      out.push({ phase: p, left: (start + 1) * 100, width: -start * 100 })
      out.push({ phase: p, left: 0, width: end * 100 })
    } else if (end > 1) {
      out.push({ phase: p, left: start * 100, width: (1 - start) * 100 })
      out.push({ phase: p, left: 0, width: (end - 1) * 100 })
    } else {
      out.push({ phase: p, left: start * 100, width: (end - start) * 100 })
    }
  }
  return out
})

const phases = computed(() => props.moon.phases || [])

const hoveredId = ref<number | null>(null)
function setHover(p: CalendarMoonPhase | null) {
  hoveredId.value = p?.id ?? p?.sort_order ?? null
}
function isHovered(p: CalendarMoonPhase) {
  return hoveredId.value === (p.id ?? p.sort_order ?? null)
}

const isActive = (p: CalendarMoonPhase) => activePhase.value && activePhase.value === p
</script>

<template>
  <div class="cycle">
    <div class="track">
      <div class="bar" aria-hidden="true" />

      <div
        v-for="(w, i) in windows"
        :key="`win-${i}`"
        :class="['window', { active: isActive(w.phase) || isHovered(w.phase) }]"
        :style="{ left: `${w.left}%`, width: `${w.width}%` }"
        :aria-label="`${w.phase.name} window`"
      />

      <div
        v-for="(p, i) in phases"
        :key="`m-${p.id ?? `s-${i}`}`"
        :class="['marker', { active: isActive(p), hover: isHovered(p) }]"
        :style="{ left: `${p.cycle_position * 100}%` }"
        @mouseenter="setHover(p)"
        @mouseleave="setHover(null)"
      >
        <div class="marker-stem" aria-hidden="true" />
        <div class="marker-glyph">
          <MoonGlyph :cycle="p.cycle_position" :size="18" />
        </div>
        <div class="marker-label">
          <span v-if="p.icon" class="icon">{{ p.icon }}</span>
          <span class="name">{{ p.name }}</span>
        </div>
      </div>

      <div
        v-if="currentCycle != null"
        class="now"
        :style="{ left: `${currentCycle * 100}%` }"
        :aria-label="`Now at cycle ${Math.round(currentCycle * 100)}%`"
      >
        <span class="now-line" aria-hidden="true" />
        <span class="now-flag">
          <MoonGlyph :cycle="currentCycle" :size="22" :glow="true" />
        </span>
      </div>
    </div>

    <div class="ticks" aria-hidden="true">
      <span class="tick" style="left: 0%">0</span>
      <span class="tick" style="left: 25%">¼</span>
      <span class="tick" style="left: 50%">½</span>
      <span class="tick" style="left: 75%">¾</span>
      <span class="tick" style="left: 100%">1</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.cycle {
  padding: 18px 8px 26px;
}

.track {
  position: relative;
  height: 90px;
  margin-top: 8px;
}

.bar {
  position: absolute;
  inset: 56px 0 auto 0;
  height: 3px;
  background: linear-gradient(to right,
    rgb(var(--ww-ink-shade) / .8) 0%,
    rgb(var(--ww-ink-shade) / .35) 12%,
    rgb(var(--ww-gold)) 50%,
    rgb(var(--ww-ink-shade) / .35) 88%,
    rgb(var(--ww-ink-shade) / .8) 100%);
  border-radius: 1px;
  box-shadow: 0 0 0 1px rgb(var(--ww-gold) / .2);
}

.window {
  position: absolute;
  top: 50px;
  height: 15px;
  background: rgb(var(--ww-vermilion) / .12);
  border-top: 1px solid rgb(var(--ww-vermilion) / .25);
  border-bottom: 1px solid rgb(var(--ww-vermilion) / .25);
  transition: background-color .35s ease;
  pointer-events: none;
}
.window.active { background: rgb(var(--ww-gold) / .22); }

.marker {
  position: absolute;
  top: 0;
  transform: translateX(-50%);
  width: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: transform .35s cubic-bezier(.22,1,.36,1);
}
.marker:hover, .marker.hover { transform: translateX(-50%) translateY(-2px); }

.marker-glyph {
  margin-bottom: 2px;
  transition: transform .35s cubic-bezier(.22,1,.36,1);
}
.marker.active .marker-glyph { transform: scale(1.25); }

.marker-stem {
  width: 1px;
  height: 18px;
  background: var(--ww-ink-hairline);
  transition: background-color .25s ease;
}
.marker.active .marker-stem, .marker.hover .marker-stem { background: rgb(var(--ww-gold)); }

.marker-label {
  position: absolute;
  top: 100%;
  margin-top: 6px;
  display: flex;
  align-items: baseline;
  gap: 4px;
  font-family: 'Cormorant SC', serif;
  font-size: 9px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  white-space: nowrap;
  opacity: 0;
  transform: translateY(-2px);
  transition: opacity .25s ease, transform .25s ease, color .25s ease;
  pointer-events: none;
}
.marker:hover .marker-label, .marker.active .marker-label {
  opacity: 1;
  transform: translateY(0);
}
.marker.active .marker-label { color: rgb(var(--ww-gold-deep)); }
.icon { font-size: 11px; }

.now {
  position: absolute;
  inset: 0 auto 0 0;
  width: 0;
  transform: translateX(-50%);
  pointer-events: none;
}
.now-line {
  position: absolute;
  top: 32px;
  bottom: 0;
  left: 50%;
  width: 1px;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 8px rgb(var(--ww-vermilion) / .5);
  transform: translateX(-50%);
  animation: now-pulse 2.5s ease-in-out infinite;
}
.now-flag {
  position: absolute;
  top: -2px;
  left: 50%;
  transform: translateX(-50%);
}

.ticks {
  position: relative;
  height: 14px;
  margin-top: 4px;
}
.tick {
  position: absolute;
  top: 0;
  transform: translateX(-50%);
  font-family: 'Cormorant SC', serif;
  font-size: 9px;
  letter-spacing: .22em;
  color: var(--ww-ink-faint);
}

@keyframes now-pulse {
  0%, 100% { opacity: .8; }
  50% { opacity: 1; box-shadow: 0 0 12px rgb(var(--ww-vermilion) / .8); }
}
</style>
