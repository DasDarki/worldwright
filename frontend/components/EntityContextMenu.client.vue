<script setup lang="ts">
const props = defineProps<{
  open: boolean
  x: number
  y: number
  nodeId: number | null
  nodeTitle: string
  canEdit?: boolean
}>()

const emit = defineEmits<{ close: [] }>()

const { t } = useI18n()
const router = useRouter()

const menuEl = ref<HTMLElement | null>(null)

// Clamp position so the menu doesn't escape the viewport.
const adjusted = computed(() => {
  if (typeof window === 'undefined') return { left: props.x, top: props.y }
  const w = 220
  const h = 180
  let left = props.x
  let top = props.y
  if (left + w > window.innerWidth - 8) left = Math.max(8, window.innerWidth - w - 8)
  if (top + h > window.innerHeight - 8) top = Math.max(8, window.innerHeight - h - 8)
  return { left, top }
})

function close() { emit('close') }

function newSubpage() {
  if (!props.nodeId) return close()
  router.push({ path: '/entities/new', query: { parent: String(props.nodeId) } })
  close()
}

// Close on Escape and on any document click outside the menu.
if (import.meta.client) {
  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && props.open) close()
  }
  function onDocClick(e: MouseEvent) {
    if (!props.open) return
    if (menuEl.value && menuEl.value.contains(e.target as Node)) return
    close()
  }
  onMounted(() => {
    document.addEventListener('keydown', onKey)
    document.addEventListener('click', onDocClick)
  })
  onBeforeUnmount(() => {
    document.removeEventListener('keydown', onKey)
    document.removeEventListener('click', onDocClick)
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="ctx-fade">
      <div
        v-if="open"
        ref="menuEl"
        class="ctx-menu"
        :style="{ left: adjusted.left + 'px', top: adjusted.top + 'px' }"
        @contextmenu.prevent
      >
        <header class="ctx-head">{{ nodeTitle }}</header>
        <button
          v-if="canEdit"
          type="button"
          class="ctx-item"
          @click="newSubpage"
        >
          <svg viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
            <path d="M8 2 V14 M2 8 H14"/>
          </svg>
          {{ t('sidebar.menu.newSubpage') }}
        </button>
        <slot />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.ctx-menu {
  position: fixed;
  z-index: 9800;
  min-width: 220px;
  background: var(--ww-card-bg);
  border: 1px solid var(--ww-ink-hairline);
  box-shadow:
    0 30px 60px -25px rgb(0 0 0 / .45),
    0 10px 20px -10px rgb(0 0 0 / .25);
  padding: 6px 0;
}
.ctx-head {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  padding: 6px 16px 10px;
  border-bottom: 1px solid var(--ww-ink-hairline);
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.ctx-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 8px 16px;
  background: transparent;
  border: 0;
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  color: rgb(var(--ww-ink));
  text-align: left;
  cursor: pointer;
  transition: background-color .2s ease, padding .2s ease;
}
.ctx-item:hover {
  background: rgb(var(--ww-gold) / .12);
  color: rgb(var(--ww-gold-deep));
  padding-left: 20px;
}
.ctx-item svg { color: var(--ww-ink-faint); flex-shrink: 0; }

.ctx-fade-enter-active, .ctx-fade-leave-active {
  transition: opacity .15s ease, transform .15s cubic-bezier(.22,1,.36,1);
}
.ctx-fade-enter-from, .ctx-fade-leave-to {
  opacity: 0;
  transform: scale(.96);
}
</style>
