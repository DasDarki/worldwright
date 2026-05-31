<script setup lang="ts">
import type { EntitySummary } from '~/types/api'

const props = defineProps<{ open: boolean; initial?: number[] }>()
const emit = defineEmits<{
  close: []
  pick: [payload: { entityIds: number[] }]
}>()

const { t } = useI18n()
const { $api } = useNuxtApp()

const entities = ref<EntitySummary[]>([])
const loading = ref(false)
const query = ref('')
const selected = ref<Set<number>>(new Set())

watch(
  () => props.open,
  async (open) => {
    if (!open) return
    selected.value = new Set(props.initial || [])
    if (!entities.value.length) {
      loading.value = true
      try {
        const res = await $api<{ entities: EntitySummary[] }>('/entities')
        entities.value = res.entities
      } finally {
        loading.value = false
      }
    }
  },
)

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  const list = entities.value
  if (!q) return list
  return list.filter((e) =>
    e.title.toLowerCase().includes(q) ||
    (e.summary || '').toLowerCase().includes(q),
  )
})

function toggle(id: number) {
  if (selected.value.has(id)) selected.value.delete(id)
  else selected.value.add(id)
  // force reactivity on Set
  selected.value = new Set(selected.value)
}

function confirm() {
  emit('pick', { entityIds: Array.from(selected.value) })
  emit('close')
}

function close() { emit('close') }
</script>

<template>
  <Teleport to="body">
    <Transition name="picker-fade">
      <div v-if="open" class="overlay" @click.self="close">
        <div class="dialog">
          <header class="head">
            <h3 class="title">{{ t('graph.pickerTitle') }}</h3>
            <button type="button" class="close" @click="close" aria-label="Close">&times;</button>
          </header>
          <p class="hint">{{ t('graph.pickerHint') }}</p>

          <input
            v-model="query"
            type="search"
            class="ww-input search"
            :placeholder="t('search.placeholder')"
          />

          <div class="list-wrap">
            <p v-if="loading" class="state">{{ t('common.loading') }}</p>
            <ul v-else class="list">
              <li v-for="e in filtered" :key="e.id">
                <label class="row" :class="{ on: selected.has(e.id) }">
                  <input
                    type="checkbox"
                    :checked="selected.has(e.id)"
                    @change="toggle(e.id)"
                  />
                  <span class="t">{{ e.title }}</span>
                  <span v-if="e.summary" class="s">{{ e.summary }}</span>
                </label>
              </li>
            </ul>
          </div>

          <footer class="foot">
            <span class="count">{{ t('graph.selected') }}: {{ selected.size }}</span>
            <div class="actions">
              <button type="button" class="cancel" @click="close">{{ t('graph.cancel') }}</button>
              <button
                type="button"
                class="ww-btn-primary"
                :disabled="selected.size < 2"
                @click="confirm"
              >{{ t('graph.insert') }}</button>
            </div>
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.overlay {
  position: fixed; inset: 0; z-index: 9700;
  background: rgb(0 0 0 / .35);
  display: flex; align-items: center; justify-content: center;
  padding: 24px;
  backdrop-filter: blur(2px);
}
.dialog {
  background: var(--ww-card-bg);
  border: 1px solid var(--ww-ink-hairline);
  width: 560px;
  max-width: 100%;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 50px 80px -40px rgb(0 0 0 / .5);
}
.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 22px 8px;
}
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 460;
  font-size: 22px;
  margin: 0;
}
.close {
  background: transparent; border: 0;
  font-size: 22px; color: var(--ww-ink-faint);
  cursor: pointer;
}
.hint { padding: 0 22px 12px; font-style: italic; color: var(--ww-ink-faint); margin: 0; }
.search { margin: 0 22px 6px; width: calc(100% - 44px); }
.list-wrap { overflow-y: auto; flex: 1; padding: 0 22px 14px; }
.list { list-style: none; margin: 0; padding: 0; }
.state { padding: 24px 0; text-align: center; font-style: italic; color: var(--ww-ink-faint); }
.row {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 4px 12px;
  align-items: baseline;
  padding: 10px 6px;
  border-bottom: 1px dashed rgb(var(--ww-ink-hairline) / .6);
  cursor: pointer;
  transition: background-color .2s ease;
  input { grid-row: 1 / span 2; }
}
.row:hover { background: rgb(var(--ww-gold) / .08); }
.row.on { background: rgb(var(--ww-vermilion) / .06); }
.t {
  font-family: 'EB Garamond', serif;
  font-size: 15px;
  color: rgb(var(--ww-ink));
}
.s {
  grid-column: 2;
  font-size: 12px;
  font-style: italic;
  color: var(--ww-ink-faint);
  line-height: 1.3;
}
.foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 22px 18px;
  border-top: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment-stain) / .25);
}
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.actions { display: flex; gap: 12px; align-items: center; }
.cancel {
  background: transparent;
  border: 0;
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  cursor: pointer;
}

.picker-fade-enter-active, .picker-fade-leave-active {
  transition: opacity .2s ease;
  .dialog { transition: transform .25s cubic-bezier(.22,1,.36,1), opacity .25s ease; }
}
.picker-fade-enter-from, .picker-fade-leave-to {
  opacity: 0;
  .dialog { opacity: 0; transform: translateY(8px); }
}
</style>
