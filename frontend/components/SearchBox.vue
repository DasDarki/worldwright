<script setup lang="ts">
interface Hit {
  id: number
  slug: string
  title: string
  summary?: string
  entity_type_id: number
  snippet?: string
}

const { t } = useI18n()
const { $api } = useNuxtApp()
const router = useRouter()

const query = ref('')
const hits = ref<Hit[]>([])
const open = ref(false)
const pending = ref(false)
const active = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)
const wrapRef = ref<HTMLElement | null>(null)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(query, (q) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (!q.trim()) {
    hits.value = []
    open.value = false
    pending.value = false
    return
  }
  pending.value = true
  open.value = true
  debounceTimer = setTimeout(async () => {
    try {
      const res = await $api<{ hits: Hit[] }>(`/search?q=${encodeURIComponent(q)}&limit=8`)
      hits.value = res.hits
      active.value = 0
    } catch {
      hits.value = []
    } finally {
      pending.value = false
    }
  }, 180)
})

function go(h: Hit) {
  open.value = false
  query.value = ''
  router.push(`/entities/${h.slug}`)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    active.value = Math.min(active.value + 1, hits.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    active.value = Math.max(active.value - 1, 0)
  } else if (e.key === 'Enter') {
    const h = hits.value[active.value]
    if (h) { e.preventDefault(); go(h) }
  } else if (e.key === 'Escape') {
    open.value = false
    inputRef.value?.blur()
  }
}

function onFocus() {
  if (query.value && hits.value.length) open.value = true
}

if (import.meta.client) {
  onMounted(() => {
    function onDocClick(e: MouseEvent) {
      if (!wrapRef.value) return
      if (!wrapRef.value.contains(e.target as Node)) open.value = false
    }
    function onShortcut(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        inputRef.value?.focus()
        inputRef.value?.select()
      }
    }
    document.addEventListener('click', onDocClick)
    document.addEventListener('keydown', onShortcut)
    onBeforeUnmount(() => {
      document.removeEventListener('click', onDocClick)
      document.removeEventListener('keydown', onShortcut)
    })
  })
}
</script>

<template>
  <div ref="wrapRef" class="search-wrap">
    <div class="search-input">
      <svg viewBox="0 0 20 20" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" aria-hidden="true">
        <circle cx="9" cy="9" r="6"/>
        <path d="M13.5 13.5 L18 18"/>
      </svg>
      <input
        ref="inputRef"
        v-model="query"
        type="search"
        :placeholder="t('search.placeholder')"
        @keydown="onKeydown"
        @focus="onFocus"
      />
      <kbd class="kbd">⌘K</kbd>
    </div>

    <Transition name="dropdown">
      <div v-if="open && (hits.length || pending)" class="dropdown">
        <ul v-if="hits.length" class="hits">
          <li
            v-for="(h, i) in hits"
            :key="h.id"
            :class="['hit', { active: i === active }]"
            @mouseenter="active = i"
            @click="go(h)"
          >
            <div class="title">{{ h.title }}</div>
            <div v-if="h.snippet" class="snippet" v-html="h.snippet" />
            <div v-else-if="h.summary" class="snippet">{{ h.summary }}</div>
          </li>
        </ul>
        <div v-else-if="pending" class="empty">{{ t('common.loading') }}</div>
        <div v-else class="empty">{{ t('search.empty') }}</div>
      </div>
    </Transition>
  </div>
</template>

<style scoped lang="scss">
.search-wrap {
  position: relative;
  width: 280px;
  @media (max-width: 880px) { display: none; }
}
.search-input {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1px solid var(--ww-ink-hairline);
  padding: 0 12px;
  height: 38px;
  transition: border-color .3s ease;
  color: rgb(var(--ww-ink));
}
.search-input:focus-within { border-color: rgb(var(--ww-gold)); }
.search-input input {
  flex: 1;
  border: 0;
  background: transparent;
  outline: none;
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  color: rgb(var(--ww-ink));
}
.search-input input::-webkit-search-cancel-button { display: none; }
.kbd {
  font-family: 'Cormorant SC', serif;
  letter-spacing: .1em;
  font-size: 10px;
  color: var(--ww-ink-faint);
  border: 1px solid var(--ww-ink-hairline);
  padding: 1px 5px;
}

.dropdown {
  position: absolute;
  top: calc(100% + 6px);
  left: 0; right: 0;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 30px 60px -25px rgb(0 0 0 / .4);
  max-height: 60vh;
  overflow-y: auto;
  z-index: 200;
}
.hits { list-style: none; margin: 0; padding: 0; }
.hit {
  padding: 12px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--ww-ink-hairline);
  transition: background-color .25s ease, padding .25s ease;
}
.hit:last-child { border-bottom: 0; }
.hit.active, .hit:hover {
  background: rgb(var(--ww-gold) / .12);
  padding-left: 20px;
}
.hit .title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 24, "wght" 500;
  font-size: 16px;
}
.hit .snippet {
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin-top: 4px;
  font-style: italic;
  line-height: 1.4;
}
.hit .snippet :deep(mark) {
  background: rgb(var(--ww-gold) / .35);
  color: inherit;
  padding: 0 2px;
}

.empty {
  padding: 18px 14px;
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  font-size: 14px;
}

.dropdown-enter-active, .dropdown-leave-active {
  transition: opacity .2s ease, transform .25s cubic-bezier(.22,1,.36,1);
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
