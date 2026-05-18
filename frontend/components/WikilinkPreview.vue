<script setup lang="ts">
import type { Entity } from '~/types/api'

const props = defineProps<{
  slug: string
  label?: string
}>()

const { fetchPreview } = useEntityPreview()
const { locale, t } = useI18n()

const linkEl = ref<any>(null)
const hovering = ref(false)
const entity = ref<Entity | null>(null)
const loadFailed = ref(false)

function anchorEl(): HTMLElement | null {
  const v = linkEl.value
  if (!v) return null
  return (v as any).$el || (v as HTMLElement)
}

let showTimer: number | null = null
let hideTimer: number | null = null
const SHOW_DELAY = 220
const HIDE_DELAY = 140

interface TooltipPos {
  left: number
  top: number
  arrow: 'top' | 'bottom'
}
const pos = ref<TooltipPos>({ left: 0, top: 0, arrow: 'top' })

function computePosition() {
  const el = anchorEl()
  if (!el || typeof window === 'undefined') return
  const rect = el.getBoundingClientRect()
  const margin = 8
  const tooltipW = 320
  const tooltipH = 180 // estimate; tooltip will size itself
  const viewportW = window.innerWidth
  const viewportH = window.innerHeight

  // prefer below the link; flip above if not enough room
  let top = rect.bottom + margin
  let arrow: 'top' | 'bottom' = 'top'
  if (top + tooltipH > viewportH && rect.top > tooltipH + margin) {
    top = rect.top - tooltipH - margin
    arrow = 'bottom'
  }
  // horizontal: center on the link, clamp to viewport
  let left = rect.left + rect.width / 2 - tooltipW / 2
  if (left < 8) left = 8
  if (left + tooltipW > viewportW - 8) left = viewportW - tooltipW - 8
  pos.value = { left: left + window.scrollX, top: top + window.scrollY, arrow }
}

async function onEnter() {
  if (hideTimer) { window.clearTimeout(hideTimer); hideTimer = null }
  if (showTimer) window.clearTimeout(showTimer)
  showTimer = window.setTimeout(async () => {
    computePosition()
    hovering.value = true
    if (!entity.value && !loadFailed.value) {
      const e = await fetchPreview(props.slug)
      if (e) {
        entity.value = e
        // recompute now that we may render more content
        nextTick(computePosition)
      } else {
        loadFailed.value = true
      }
    }
  }, SHOW_DELAY)
}

function onLeave() {
  if (showTimer) { window.clearTimeout(showTimer); showTimer = null }
  if (hideTimer) window.clearTimeout(hideTimer)
  hideTimer = window.setTimeout(() => {
    hovering.value = false
  }, HIDE_DELAY)
}

function onTooltipEnter() {
  if (hideTimer) { window.clearTimeout(hideTimer); hideTimer = null }
}

onBeforeUnmount(() => {
  if (showTimer) window.clearTimeout(showTimer)
  if (hideTimer) window.clearTimeout(hideTimer)
})

const typeName = computed(() => {
  const et = entity.value?.entity_type
  if (!et) return ''
  return locale.value === 'de' ? et.name_de : et.name_en
})

const display = computed(() => props.label || entity.value?.title || props.slug)
</script>

<template>
  <NuxtLink
    ref="linkEl"
    :to="`/entities/${slug}`"
    class="ww-link wikilink"
    @mouseenter="onEnter"
    @mouseleave="onLeave"
    @focus="onEnter"
    @blur="onLeave"
  >{{ display }}</NuxtLink>

  <Teleport to="body">
    <Transition name="wl-fade">
      <div
        v-if="hovering && (entity || loadFailed)"
        class="wl-preview"
        :class="`arrow-${pos.arrow}`"
        :style="{ left: pos.left + 'px', top: pos.top + 'px' }"
        @mouseenter="onTooltipEnter"
        @mouseleave="onLeave"
      >
        <template v-if="entity">
          <div class="head">
            <span v-if="typeName" class="type-chip">{{ typeName }}</span>
          </div>
          <h4 class="title">{{ entity.title }}</h4>
          <p v-if="entity.summary" class="summary">{{ entity.summary }}</p>
          <div v-if="entity.tags?.length" class="tags">
            <span v-for="tag in entity.tags.slice(0, 6)" :key="tag" class="tag">{{ tag }}</span>
          </div>
        </template>
        <template v-else>
          <p class="missing">{{ t('entity.previewMissing') }}</p>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.wikilink { cursor: pointer; }

.wl-preview {
  position: absolute;
  z-index: 9500;
  width: 320px;
  max-width: calc(100vw - 24px);
  padding: 14px 16px 16px;
  background: var(--ww-card-bg);
  border: 1px solid var(--ww-ink-hairline);
  box-shadow:
    0 30px 60px -25px rgb(0 0 0 / .35),
    0 12px 24px -12px rgb(0 0 0 / .18);
  pointer-events: auto;
  font-family: 'EB Garamond', serif;
}
.wl-preview::after {
  content: '';
  position: absolute;
  left: 50%;
  width: 10px;
  height: 10px;
  background: var(--ww-card-bg);
  border-right: 1px solid var(--ww-ink-hairline);
  border-bottom: 1px solid var(--ww-ink-hairline);
  transform: translateX(-50%) rotate(45deg);
}
.wl-preview.arrow-top::after {
  top: -6px;
  border: 1px solid var(--ww-ink-hairline);
  border-right: 0;
  border-bottom: 0;
  transform: translateX(-50%) rotate(45deg);
}
.wl-preview.arrow-bottom::after {
  bottom: -6px;
}

.head { margin-bottom: 6px; }
.type-chip {
  display: inline-block;
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  background: rgb(var(--ww-vermilion) / .08);
  padding: 3px 9px;
}
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 460;
  font-size: 22px;
  line-height: 1.15;
  margin: 0 0 6px;
  letter-spacing: -0.01em;
  color: rgb(var(--ww-ink));
}
.summary {
  font-style: italic;
  font-size: 14px;
  line-height: 1.45;
  color: rgb(var(--ww-ink-shade));
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.tag {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  border: 1px solid var(--ww-ink-hairline);
  padding: 2px 8px;
}
.missing {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 14px;
  margin: 0;
}

.wl-fade-enter-active, .wl-fade-leave-active {
  transition: opacity .18s ease, transform .18s cubic-bezier(.22,1,.36,1);
}
.wl-fade-enter-from, .wl-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
