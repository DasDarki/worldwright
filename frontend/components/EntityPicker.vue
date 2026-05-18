<script setup lang="ts">
import type { EntitySummary } from '~/types/api'

interface Hit {
  id: number
  slug: string
  title: string
  summary?: string
  entity_type_id: number
}

const props = withDefaults(defineProps<{
  modelValue: EntitySummary | null
  exclude?: number[]
  placeholder?: string
}>(), {
  exclude: () => [],
  placeholder: 'Search the codex…',
})

const emit = defineEmits<{ 'update:modelValue': [v: EntitySummary | null] }>()

const { $api } = useNuxtApp()

const query = ref('')
const hits = ref<Hit[]>([])
const open = ref(false)
const active = ref(0)
const pending = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
const wrapRef = ref<HTMLElement | null>(null)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

async function refresh(q: string) {
  pending.value = true
  try {
    if (!q.trim()) {
      const { entities } = await $api<{ entities: EntitySummary[] }>('/entities')
      hits.value = entities.slice(0, 12)
    } else {
      const { hits: results } = await $api<{ hits: Hit[] }>(`/search?q=${encodeURIComponent(q)}&limit=12`)
      hits.value = results
    }
  } catch {
    hits.value = []
  } finally {
    pending.value = false
    active.value = 0
  }
}

watch(query, (q) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  open.value = true
  debounceTimer = setTimeout(() => refresh(q), 160)
})

function pick(h: Hit | EntitySummary) {
  emit('update:modelValue', {
    id: h.id,
    slug: h.slug,
    title: h.title,
    entity_type_id: h.entity_type_id,
    summary: h.summary,
    visibility: 'player',
  })
  open.value = false
  query.value = ''
}

function clear() {
  emit('update:modelValue', null)
  query.value = ''
}

async function onFocus() {
  if (!hits.value.length) await refresh('')
  open.value = true
}

function onKeydown(e: KeyboardEvent) {
  const filtered = hits.value.filter((h) => !props.exclude.includes(h.id))
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    active.value = Math.min(active.value + 1, filtered.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    active.value = Math.max(active.value - 1, 0)
  } else if (e.key === 'Enter') {
    const h = filtered[active.value]
    if (h) { e.preventDefault(); pick(h) }
  } else if (e.key === 'Escape') {
    open.value = false
  }
}

const filteredHits = computed(() => hits.value.filter((h) => !props.exclude.includes(h.id)))

if (import.meta.client) {
  onMounted(() => {
    function onDocClick(e: MouseEvent) {
      if (!wrapRef.value) return
      if (!wrapRef.value.contains(e.target as Node)) open.value = false
    }
    document.addEventListener('click', onDocClick)
    onBeforeUnmount(() => document.removeEventListener('click', onDocClick))
  })
}
</script>

<template>
  <div ref="wrapRef" class="picker">
    <div v-if="modelValue" class="selected">
      <span class="selected-title">{{ modelValue.title }}</span>
      <button type="button" class="x" @click="clear" aria-label="Clear">×</button>
    </div>
    <div v-else class="input-wrap">
      <input
        ref="inputRef"
        v-model="query"
        class="ww-input"
        :placeholder="placeholder"
        @focus="onFocus"
        @keydown="onKeydown"
      />
      <Transition name="dropdown">
        <ul v-if="open && (filteredHits.length || pending)" class="dropdown">
          <li
            v-for="(h, i) in filteredHits"
            :key="h.id"
            :class="['hit', { active: i === active }]"
            @mouseenter="active = i"
            @click="pick(h)"
          >
            <div class="title">{{ h.title }}</div>
            <div v-if="h.summary" class="summary">{{ h.summary }}</div>
          </li>
          <li v-if="!filteredHits.length && !pending" class="empty">No matches.</li>
        </ul>
      </Transition>
    </div>
  </div>
</template>

<style scoped lang="scss">
.picker {
  position: relative;
}
.selected {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px 6px 14px;
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-gold) / .08);
  font-family: 'EB Garamond', serif;
}
.selected-title {
  font-style: italic;
  font-size: 16px;
}
.x {
  background: none; border: 0;
  font-size: 18px; line-height: 1;
  color: var(--ww-ink-faint);
  cursor: pointer;
  padding: 0 4px;
  transition: color .25s ease;
}
.x:hover { color: rgb(var(--ww-vermilion)); }

.input-wrap { position: relative; }
.dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0; right: 0;
  list-style: none; margin: 0; padding: 0;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 25px 50px -20px rgb(0 0 0 / .35);
  max-height: 40vh;
  overflow-y: auto;
  z-index: 50;
}
.hit {
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--ww-ink-hairline);
  transition: background-color .25s ease, padding .25s ease;
}
.hit:last-child { border-bottom: 0; }
.hit.active, .hit:hover { background: rgb(var(--ww-gold) / .12); padding-left: 18px; }
.hit .title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
  font-size: 15px;
}
.hit .summary {
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin-top: 2px;
  font-style: italic;
}
.empty {
  padding: 14px;
  text-align: center;
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 13px;
}

.dropdown-enter-active, .dropdown-leave-active {
  transition: opacity .2s ease, transform .25s cubic-bezier(.22,1,.36,1);
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
