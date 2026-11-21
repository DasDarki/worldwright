<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, useId } from "vue"

type CrystalState = 'idle' | 'loading' | 'success'

const props = withDefaults(defineProps<{
  size?: string
  interactive?: boolean
  state?: CrystalState
  glow?: boolean
}>(), {
  size: '48px',
  interactive: false,
  state: 'idle',
  glow: true,
})

interface Facet {
  name: string
  points: [number, number][]
  normal: [number, number]
  base: number
  hue?: number
}

const top: [number, number] = [513, 105]
const upperRight: [number, number] = [860, 495]
const frontTopRight: [number, number] = [748, 505]
const lowerRight: [number, number] = [860, 1080]
const frontBottomRight: [number, number] = [748, 1012]
const bottom: [number, number] = [510, 1452]
const lowerLeft: [number, number] = [165, 1075]
const upperLeft: [number, number] = [163, 522]
const ridge: [number, number] = [445, 390]
const frontTopLeft: [number, number] = [367, 700]
const frontBottomLeft: [number, number] = [380, 1210]
const bevelTop: [number, number] = [245, 592]
const bevelBottom: [number, number] = [250, 1012]
const splitTop: [number, number] = [748, 600]
const splitBottom: [number, number] = [447, 1335]

const facets: Facet[] = [
  { name: 'crown-left', points: [top, upperLeft, ridge], normal: [-0.55, -0.83], base: 34 },
  { name: 'crown-bright', points: [upperLeft, ridge, frontTopLeft], normal: [-0.85, -0.5], base: 50 },
  { name: 'crown-front', points: [top, frontTopLeft, frontTopRight], normal: [0.15, -0.95], base: 13 },
  { name: 'crown-right', points: [top, frontTopRight, upperRight], normal: [0.8, -0.6], base: 38, hue: 216 },
  { name: 'side-right', points: [frontTopRight, upperRight, lowerRight, frontBottomRight], normal: [1, -0.05], base: 27 },
  { name: 'side-left-outer', points: [upperLeft, bevelTop, bevelBottom, lowerLeft], normal: [-1, 0], base: 24 },
  { name: 'side-left-inner', points: [bevelTop, frontTopLeft, frontBottomLeft, bevelBottom], normal: [-0.7, 0.15], base: 33 },
  { name: 'front-left', points: [frontTopLeft, frontTopRight, splitTop, splitBottom, frontBottomLeft], normal: [-0.1, 0.1], base: 7 },
  { name: 'front-right', points: [splitTop, frontBottomRight, bottom, splitBottom], normal: [0.25, 0.2], base: 3 },
  { name: 'base-left', points: [lowerLeft, frontBottomLeft, bottom], normal: [-0.55, 0.83], base: 9 },
  { name: 'base-right', points: [frontBottomRight, lowerRight, bottom], normal: [0.6, 0.8], base: 21 },
]

const outline: [number, number][] = [top, upperRight, lowerRight, bottom, lowerLeft, upperLeft]

const innerEdges: [number, number][][] = [
  [top, frontTopLeft],
  [ridge, upperLeft],
  [top, frontTopRight],
  [frontTopLeft, frontTopRight, upperRight],
  [frontTopLeft, frontBottomLeft, bottom],
  [frontTopRight, frontBottomRight, lowerRight],
  [frontBottomRight, bottom],
  [frontBottomLeft, lowerLeft],
  [bevelTop, bevelBottom],
  [splitTop, splitBottom],
]

const vertices: [number, number][] = [top, upperRight, lowerRight, bottom, lowerLeft, upperLeft, frontTopRight]
const center: [number, number] = [512, 780]

function toPoints(points: [number, number][]) {
  return points.map(([x, y]) => `${x},${y}`).join(' ')
}

const outlinePath = `M${toPoints(outline).split(' ').join(' L')} Z`

const light = reactive({ x: 0.62, y: -0.78 })
const target = { x: 0.62, y: -0.78 }
let lastPointer = 0
let frame = 0
let reducedMotion = false
const root = ref<SVGSVGElement | null>(null)
const gradientId = useId()
const glowId = useId()

function facetFill(facet: Facet) {
  const dot = facet.normal[0] * light.x + facet.normal[1] * light.y
  const lightness = Math.min(70, Math.max(2, facet.base + dot * 13))
  const saturation = 14 + Math.max(0, dot) * 18
  return `hsl(${facet.hue ?? 221} ${saturation.toFixed(1)}% ${lightness.toFixed(1)}%)`
}

const glintVertex = computed(() => {
  let best = vertices[0]!
  let bestScore = -Infinity
  for (const vertex of vertices) {
    const deltaX = vertex[0] - center[0]
    const deltaY = vertex[1] - center[1]
    const length = Math.hypot(deltaX, deltaY)
    const score = (deltaX / length) * light.x + (deltaY / length) * light.y
    if (score > bestScore) {
      bestScore = score
      best = vertex
    }
  }
  return best
})

function onPointerMove(event: PointerEvent) {
  const element = root.value
  if (!element) {
    return
  }
  const rect = element.getBoundingClientRect()
  const deltaX = event.clientX - (rect.left + rect.width / 2)
  const deltaY = event.clientY - (rect.top + rect.height / 2)
  const length = Math.hypot(deltaX, deltaY) || 1
  target.x = deltaX / length
  target.y = deltaY / length
  lastPointer = performance.now()
}

function tick(time: number) {
  if (!reducedMotion && time - lastPointer > 4000) {
    const angle = time / 5200
    target.x = Math.cos(angle) * 0.8
    target.y = Math.sin(angle) * 0.55 - 0.45
  }
  light.x += (target.x - light.x) * 0.06
  light.y += (target.y - light.y) * 0.06
  frame = requestAnimationFrame(tick)
}

onMounted(() => {
  reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (props.interactive) {
    lastPointer = -10000
    window.addEventListener('pointermove', onPointerMove, { passive: true })
    frame = requestAnimationFrame(tick)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('pointermove', onPointerMove)
  cancelAnimationFrame(frame)
})
</script>

<template>
  <svg
    ref="root"
    class="crystal"
    :class="[`is-${state}`, { 'has-glow': glow }]"
    :style="{ width: size }"
    viewBox="110 60 804 1440"
    role="img"
    aria-label="Onyx"
  >
    <defs>
      <radialGradient :id="gradientId">
        <stop offset="0" stop-color="#ffffff" stop-opacity="0.95" />
        <stop offset="0.25" stop-color="#bcd6ff" stop-opacity="0.55" />
        <stop offset="1" stop-color="#5b9dff" stop-opacity="0" />
      </radialGradient>
      <filter :id="glowId" x="-30%" y="-30%" width="160%" height="160%">
        <feGaussianBlur stdDeviation="14" />
      </filter>
    </defs>
    <path v-if="glow" class="halo" :d="outlinePath" :filter="`url(#${glowId})`" />
    <g class="facets">
      <polygon
        v-for="facet in facets"
        :key="facet.name"
        :points="toPoints(facet.points)"
        :style="{ fill: facetFill(facet) }"
      />
    </g>
    <g class="edges">
      <polyline v-for="(edge, index) in innerEdges" :key="index" :points="toPoints(edge)" />
    </g>
    <path class="outline" :d="outlinePath" />
    <path class="runner" :d="outlinePath" pathLength="100" />
    <circle
      class="glint"
      :cx="glintVertex[0]"
      :cy="glintVertex[1]"
      r="70"
      :fill="`url(#${gradientId})`"
    />
  </svg>
</template>

<style scoped>
.crystal {
  display: block;
  height: auto;
  overflow: visible;
}

.facets polygon {
  transition: fill 500ms cubic-bezier(0.2, 0.8, 0.2, 1);
}

.edges polyline {
  fill: none;
  stroke: rgba(188, 214, 255, 0.28);
  stroke-width: 3;
  stroke-linejoin: round;
}

.outline {
  fill: none;
  stroke: rgba(140, 186, 255, 0.75);
  stroke-width: 5;
  stroke-linejoin: round;
}

.halo {
  fill: none;
  stroke: #5b9dff;
  stroke-width: 18;
  opacity: 0.32;
  transition: opacity 600ms cubic-bezier(0.2, 0.8, 0.2, 1);
}

.runner {
  fill: none;
  stroke: #bcd6ff;
  stroke-width: 7;
  stroke-linecap: round;
  stroke-dasharray: 9 91;
  opacity: 0;
}

.glint {
  mix-blend-mode: screen;
  transition: cx 600ms cubic-bezier(0.2, 0.8, 0.2, 1), cy 600ms cubic-bezier(0.2, 0.8, 0.2, 1);
}

.is-idle .runner {
  opacity: 0.85;
  animation: run 9s linear infinite;
}

.is-loading .facets {
  animation: breathe 1.8s ease-in-out infinite;
}

.is-loading .halo {
  animation: halo 1.8s ease-in-out infinite;
}

.is-loading .runner {
  opacity: 1;
  stroke-dasharray: 18 82;
  animation: run 1.4s linear infinite;
}

.is-success .runner {
  opacity: 1;
  stroke-dasharray: 100 0;
  animation: flash 900ms cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
}

.is-success .halo {
  opacity: 0.6;
}

@keyframes run {
  to {
    stroke-dashoffset: -100;
  }
}

@keyframes breathe {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.55;
  }
}

@keyframes halo {
  0%,
  100% {
    opacity: 0.18;
  }
  50% {
    opacity: 0.5;
  }
}

@keyframes flash {
  from {
    opacity: 1;
    stroke-width: 12;
  }
  to {
    opacity: 0.2;
    stroke-width: 5;
  }
}

@media (prefers-reduced-motion: reduce) {
  .is-idle .runner,
  .is-loading .runner {
    animation: none;
    opacity: 0;
  }
}
</style>
