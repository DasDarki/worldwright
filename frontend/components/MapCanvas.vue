<script setup lang="ts">
import type { MapPin } from '~/types/api'

const props = withDefaults(defineProps<{
  assetId: number
  pins: MapPin[]
  editable?: boolean
  selectedPinId?: number | null
}>(), {
  editable: false,
  selectedPinId: null,
})

const emit = defineEmits<{
  placePin:  [pos: { x: number; y: number }]
  selectPin: [pin: MapPin]
  movePin:   [pin: MapPin, pos: { x: number; y: number }]
  loaded:    [dims: { width: number; height: number }]
}>()

const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)

const wrapRef = ref<HTMLElement | null>(null)
const imgRef = ref<HTMLImageElement | null>(null)
const draggingId = ref<number | null>(null)

function clientToFraction(clientX: number, clientY: number) {
  if (!wrapRef.value) return { x: 0, y: 0 }
  const r = wrapRef.value.getBoundingClientRect()
  return {
    x: Math.min(1, Math.max(0, (clientX - r.left) / r.width)),
    y: Math.min(1, Math.max(0, (clientY - r.top) / r.height)),
  }
}

function onCanvasClick(e: MouseEvent) {
  if (!props.editable) return
  if (draggingId.value !== null) { draggingId.value = null; return }
  if ((e.target as HTMLElement).closest('.pin')) return
  emit('placePin', clientToFraction(e.clientX, e.clientY))
}

function onPinClick(pin: MapPin, e: MouseEvent) {
  e.stopPropagation()
  if (draggingId.value !== null) return
  emit('selectPin', pin)
}

function onPointerDown(pin: MapPin, e: PointerEvent) {
  if (!props.editable) return
  e.stopPropagation()
  ;(e.currentTarget as HTMLElement).setPointerCapture(e.pointerId)
  draggingId.value = pin.id
}
function onPointerMove(pin: MapPin, e: PointerEvent) {
  if (draggingId.value !== pin.id) return
  const pos = clientToFraction(e.clientX, e.clientY)
  emit('movePin', pin, pos)
}
function onPointerUp(pin: MapPin, e: PointerEvent) {
  if (draggingId.value !== pin.id) return
  ;(e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId)
  setTimeout(() => { draggingId.value = null }, 50)
}

function onImgLoad() {
  if (!imgRef.value) return
  emit('loaded', { width: imgRef.value.naturalWidth, height: imgRef.value.naturalHeight })
}
</script>

<template>
  <div
    ref="wrapRef"
    :class="['map-canvas', { editable }]"
    @click="onCanvasClick"
  >
    <img
      ref="imgRef"
      :src="`${apiBase}/assets/${assetId}`"
      alt="Map"
      class="map-img"
      draggable="false"
      @load="onImgLoad"
    />
    <TransitionGroup name="pin" tag="div" class="pins">
      <button
        v-for="pin in pins"
        :key="pin.id"
        type="button"
        class="pin-wrap"
        :class="{
          selected: pin.id === selectedPinId,
          secret: pin.visibility === 'secret',
          player: pin.visibility === 'player',
          public: pin.visibility === 'public',
          dragging: pin.id === draggingId,
        }"
        :style="{ left: `${pin.x * 100}%`, top: `${pin.y * 100}%` }"
        @click="(e) => onPinClick(pin, e)"
        @pointerdown="(e) => onPointerDown(pin, e)"
        @pointermove="(e) => onPointerMove(pin, e)"
        @pointerup="(e) => onPointerUp(pin, e)"
      >
        <span class="halo" aria-hidden="true" />
        <span class="dot" aria-hidden="true" />
        <span v-if="pin.label" class="label">{{ pin.label }}</span>
      </button>
    </TransitionGroup>
  </div>
</template>

<style scoped lang="scss">
.map-canvas {
  position: relative;
  display: block;
  width: 100%;
  background: rgb(var(--ww-parchment-deep) / .4);
  border: 1px solid var(--ww-ink-hairline);
  user-select: none;
  overflow: hidden;
}
.map-canvas.editable { cursor: crosshair; }
.map-img {
  display: block;
  width: 100%;
  height: auto;
  pointer-events: none;
  user-select: none;
}
.pins { position: absolute; inset: 0; }

.pin-wrap {
  position: absolute;
  transform: translate(-50%, -100%);
  background: transparent;
  border: 0;
  padding: 0;
  cursor: grab;
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  outline: none;
  animation: pinDrop .55s cubic-bezier(.22,1,.36,1) backwards;
}
.pin-wrap.dragging { cursor: grabbing; }
.halo {
  position: absolute;
  bottom: -4px;
  width: 24px; height: 24px;
  border-radius: 50%;
  background: rgb(var(--ww-vermilion) / .25);
  transform: scale(.7);
  transition: transform .35s cubic-bezier(.22,1,.36,1), background-color .35s ease;
  pointer-events: none;
}
.pin-wrap:hover .halo { transform: scale(1.2); background: rgb(var(--ww-vermilion) / .4); }
.pin-wrap.selected .halo {
  background: rgb(var(--ww-gold) / .45);
  animation: pulse 1.6s ease-in-out infinite;
}
.dot {
  width: 14px; height: 14px;
  border-radius: 50% 50% 50% 0;
  transform: rotate(-45deg);
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 4px 8px rgb(0 0 0 / .35), inset 0 0 0 2px rgb(var(--ww-parchment));
  position: relative;
  z-index: 1;
  transition: transform .35s cubic-bezier(.22,1,.36,1), background-color .35s ease;
}
.pin-wrap.player .dot { background: rgb(var(--ww-gold)); }
.pin-wrap.public .dot { background: rgb(var(--ww-ink)); }
.pin-wrap.selected .dot {
  background: rgb(var(--ww-gold-bright));
  transform: rotate(-45deg) scale(1.2);
}
.pin-wrap:hover .dot { transform: rotate(-45deg) scale(1.1) translateY(-2px); }
.label {
  position: absolute;
  bottom: calc(100% + 6px);
  background: rgb(var(--ww-ink));
  color: rgb(var(--ww-parchment));
  padding: 4px 10px;
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  white-space: nowrap;
  opacity: 0;
  transform: translateY(4px);
  transition: opacity .25s ease, transform .25s cubic-bezier(.22,1,.36,1);
  pointer-events: none;
}
.pin-wrap:hover .label, .pin-wrap.selected .label { opacity: 1; transform: translateY(0); }

@keyframes pinDrop {
  0%   { transform: translate(-50%, -160%) scale(.6); opacity: 0; }
  60%  { transform: translate(-50%, -90%)  scale(1.1); opacity: 1; }
  100% { transform: translate(-50%, -100%) scale(1); opacity: 1; }
}
@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.4); }
}

.pin-enter-active, .pin-leave-active { transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1); }
.pin-enter-from { opacity: 0; transform: translate(-50%, -130%) scale(.6); }
.pin-leave-to { opacity: 0; transform: translate(-50%, -80%) scale(.6); }
</style>
