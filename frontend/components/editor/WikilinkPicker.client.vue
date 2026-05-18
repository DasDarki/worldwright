<script setup lang="ts">
import type { EntitySummary } from '~/types/api'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; pick: [{ slug: string; label: string }] }>()

const { $api } = useNuxtApp()
const query = ref('')
const results = ref<EntitySummary[]>([])
const pending = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)
const activeIndex = ref(0)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

async function search(q: string) {
  if (!q.trim()) {
    try {
      const { entities } = await $api<{ entities: EntitySummary[] }>('/entities')
      results.value = entities.slice(0, 12)
    } catch {
      results.value = []
    }
    return
  }
  try {
    const { hits } = await $api<{ hits: { id: number; slug: string; title: string; entity_type_id: number; summary?: string }[] }>(
      `/search?q=${encodeURIComponent(q)}&limit=12`,
    )
    results.value = hits.map((h) => ({
      id: h.id,
      slug: h.slug,
      title: h.title,
      summary: h.summary,
      entity_type_id: h.entity_type_id,
      visibility: 'player',
    }))
  } catch {
    results.value = []
  }
}

watch(query, (q) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  pending.value = true
  debounceTimer = setTimeout(async () => {
    await search(q)
    pending.value = false
    activeIndex.value = 0
  }, 160)
})

watch(
  () => props.open,
  async (isOpen) => {
    if (isOpen) {
      query.value = ''
      await search('')
      await nextTick()
      inputRef.value?.focus()
    }
  },
)

function pick(e: EntitySummary) {
  emit('pick', { slug: e.slug, label: e.title })
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    e.preventDefault()
    emit('close')
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, results.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
  } else if (e.key === 'Enter') {
    const hit = results.value[activeIndex.value]
    if (hit) {
      e.preventDefault()
      pick(hit)
    }
  }
}
</script>

<template>
  <Transition name="picker">
    <div v-if="open" class="picker-backdrop" @click.self="emit('close')">
      <div class="picker" role="dialog" aria-label="Insert wikilink" @keydown="onKeydown">
        <div class="ww-eyebrow picker-head">Insert wikilink</div>
        <input
          ref="inputRef"
          v-model="query"
          class="ww-input picker-input"
          placeholder="Search the codex…"
        />
        <ul class="results">
          <li
            v-for="(r, i) in results"
            :key="r.id"
            :class="['result', { active: i === activeIndex }]"
            @mouseenter="activeIndex = i"
            @click="pick(r)"
          >
            <div class="title">{{ r.title }}</div>
            <div v-if="r.summary" class="summary">{{ r.summary }}</div>
          </li>
          <li v-if="!results.length && !pending" class="empty">No matches.</li>
        </ul>
      </div>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.picker-backdrop {
  position: fixed; inset: 0;
  background: rgb(0 0 0 / .25);
  backdrop-filter: blur(2px);
  z-index: 200;
  display: flex; align-items: flex-start; justify-content: center;
  padding: 14vh 20px 20px;
}
.picker {
  width: min(560px, 100%);
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  padding: 22px 24px 18px;
  box-shadow: 0 60px 100px -40px rgb(0 0 0 / .45);
}
.picker-head { margin-bottom: 14px; }
.picker-input { margin-bottom: 14px; font-size: 18px; }
.results { list-style: none; margin: 0; padding: 0; max-height: 50vh; overflow-y: auto; }
.result {
  padding: 12px 8px;
  cursor: pointer;
  border-bottom: 1px solid var(--ww-ink-hairline);
  transition: background-color .25s ease, padding .25s ease;
  &:hover, &.active {
    background: rgb(var(--ww-gold) / .12);
    padding-left: 14px;
  }
}
.result .title {
  font-family: 'Fraunces', serif;
  font-size: 18px;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
}
.result .summary {
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin-top: 2px;
  font-style: italic;
}
.empty {
  padding: 30px 0;
  text-align: center;
  font-style: italic;
  color: var(--ww-ink-faint);
}

.picker-enter-active, .picker-leave-active { transition: opacity .25s ease; }
.picker-enter-active .picker, .picker-leave-active .picker {
  transition: transform .35s cubic-bezier(.22,1,.36,1), opacity .25s ease;
}
.picker-enter-from, .picker-leave-to { opacity: 0; }
.picker-enter-from .picker, .picker-leave-to .picker {
  transform: translateY(-8px) scale(.98);
  opacity: 0;
}
</style>
