<script setup lang="ts">
import { moonGlyphPath } from '~/composables/useCalendar'

const props = withDefaults(defineProps<{
  cycle: number
  size?: number
  glow?: boolean
}>(), {
  size: 32,
  glow: false,
})

const R = 48
const path = computed(() => moonGlyphPath(props.cycle, R))
const isEmpty = computed(() => props.cycle < 0.005 || props.cycle > 0.995)
const fullness = computed(() => Math.abs(Math.sin(props.cycle * Math.PI)))
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="-52 -52 104 104"
    class="moon-glyph"
    :class="{ glow }"
    aria-hidden="true"
  >
    <defs>
      <radialGradient :id="`moon-lit-${size}`" cx="32%" cy="28%" r="65%">
        <stop offset="0%"  stop-color="rgb(var(--ww-gold-bright))" />
        <stop offset="100%" stop-color="rgb(var(--ww-gold))" />
      </radialGradient>
      <radialGradient :id="`moon-dark-${size}`" cx="68%" cy="72%" r="80%">
        <stop offset="0%"  stop-color="rgb(var(--ww-ink-shade) / .55)" />
        <stop offset="100%" stop-color="rgb(var(--ww-ink-shade) / .9)" />
      </radialGradient>
    </defs>

    <circle r="48" :fill="`url(#moon-dark-${size})`" />
    <path v-if="!isEmpty" :d="path" :fill="`url(#moon-lit-${size})`" />
    <circle r="48" fill="none" stroke="rgb(var(--ww-gold) / .25)" stroke-width="1" />
    <circle
      v-if="glow"
      r="48"
      fill="none"
      stroke="rgb(var(--ww-gold))"
      stroke-width="0.6"
      :opacity="0.15 + fullness * 0.4"
      class="halo-pulse"
    />
  </svg>
</template>

<style scoped lang="scss">
.moon-glyph {
  flex-shrink: 0;
  filter: drop-shadow(0 2px 4px rgb(0 0 0 / .25));
}
.moon-glyph.glow {
  filter: drop-shadow(0 0 8px rgb(var(--ww-gold) / .35)) drop-shadow(0 2px 4px rgb(0 0 0 / .25));
}
.halo-pulse {
  transform-origin: center;
  animation: halo 6s ease-in-out infinite;
}
@keyframes halo {
  0%, 100% { transform: scale(1); opacity: 0.2; }
  50%      { transform: scale(1.06); opacity: 0.5; }
}
</style>
