<script setup lang="ts">
import type { EntitySummary, EntityType } from '~/types/api'
import { useAuthStore } from '~/stores/auth'
import { useToastsStore } from '~/stores/toasts'

const { t } = useI18n()
const route = useRoute()
const auth = useAuthStore()
const toasts = useToastsStore()
const { $api } = useNuxtApp()

const { data: entitiesData, refresh: refreshEntities } = await useAsyncData('sidebar-entities', () =>
  $api<{ entities: EntitySummary[] }>('/entities'),
)
const { data: typesData } = await useAsyncData('sidebar-types', () =>
  $api<{ entity_types: EntityType[] }>('/entity-types'),
)

// We keep a local copy of the entity list so we can mutate order on drop
// before the server confirms.
const localEntities = ref<EntitySummary[]>([])
watch(entitiesData, (d) => {
  if (d?.entities) localEntities.value = [...d.entities]
}, { immediate: true })

const types = computed(() => typesData.value?.entity_types || [])

const childrenOf = computed(() => {
  const m = new Map<number, EntitySummary[]>()
  for (const e of localEntities.value) {
    const key = e.parent_id ?? 0
    const arr = m.get(key) || []
    arr.push(e)
    m.set(key, arr)
  }
  return m
})

const roots = computed(() => childrenOf.value.get(0) || [])

const activeSlug = computed(() => (route.params.slug as string) || '')

interface DragState {
  id: number
  parentId: number | null
}
const dragState = ref<DragState | null>(null)

function onDragStart(payload: DragState) {
  dragState.value = payload
}
function onDragEnd() {
  dragState.value = null
}

async function onReorder(payload: { draggedId: number; targetId: number; position: 'before' | 'after' }) {
  const { draggedId, targetId, position } = payload
  if (draggedId === targetId) return

  const target = localEntities.value.find((e) => e.id === targetId)
  const dragged = localEntities.value.find((e) => e.id === draggedId)
  if (!target || !dragged) return

  // Prevent dropping a node into one of its own descendants.
  if (isAncestor(draggedId, target.parent_id ?? null)) return

  const newParentId = target.parent_id ?? null

  // Optimistic UI: rebuild localEntities so the dragged node sits at the new
  // position with the new parent. childrenOf re-derives order from this array,
  // so no sort_order field is needed locally.
  const withoutDragged = localEntities.value.filter((e) => e.id !== draggedId)
  const targetIdx = withoutDragged.findIndex((e) => e.id === targetId)
  if (targetIdx < 0) return
  const insertAt = position === 'before' ? targetIdx : targetIdx + 1
  const movedDragged: EntitySummary = { ...dragged, parent_id: newParentId ?? undefined }
  withoutDragged.splice(insertAt, 0, movedDragged)
  localEntities.value = withoutDragged

  // Compute the new sort_order for every sibling of the (new) parent.
  const siblings = (childrenOf.value.get(newParentId ?? 0) || [])
  const order = siblings.map((e, i) => ({
    id: e.id,
    parent_id: e.id === draggedId ? newParentId : undefined,
    sort_order: (i + 1) * 100,
  }))

  try {
    await $api('/entities/reorder', { method: 'POST', body: { order } })
    // refetch in the background to pick up server-side state (paths, etc.)
    refreshEntities()
  } catch (e: any) {
    toasts.error(e?.data?.error || t('sidebar.reorderFailed'))
    await refreshEntities()
  }
}

function isAncestor(maybeAncestorID: number, descendantParentID: number | null): boolean {
  if (descendantParentID == null) return false
  let cur: number | undefined = descendantParentID
  const safety = 50
  for (let i = 0; i < safety && cur != null; i++) {
    if (cur === maybeAncestorID) return true
    const parent = localEntities.value.find((e) => e.id === cur)
    cur = parent?.parent_id ?? undefined
  }
  return false
}
</script>

<template>
  <aside class="sidebar">
    <div class="head">
      <div class="ww-eyebrow">{{ t('sidebar.codex') }}</div>
      <NuxtLink v-if="auth.isAdmin" to="/entities/new" class="new-link" :title="t('sidebar.new')">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
          <path d="M7 1 V 13 M1 7 H 13"/>
        </svg>
      </NuxtLink>
    </div>
    <nav class="tree" :aria-label="t('sidebar.codex')">
      <ul class="roots">
        <EntityTreeNode
          v-for="r in roots"
          :key="r.id"
          :node="r"
          :children-of="childrenOf"
          :types="types"
          :depth="0"
          :active-slug="activeSlug"
          :drag-state="dragState"
          :can-drag="auth.isAdmin"
          @drag-start="onDragStart"
          @drag-end="onDragEnd"
          @reorder="onReorder"
        />
      </ul>
      <p v-if="!roots.length" class="empty">{{ t('sidebar.empty') }}</p>
    </nav>
  </aside>
</template>

<style scoped lang="scss">
.sidebar {
  width: 280px;
  flex-shrink: 0;
  position: sticky;
  top: 96px;
  align-self: flex-start;
  max-height: calc(100vh - 120px);
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--ww-ink-hairline);
  padding-right: 18px;
  @media (max-width: 960px) { display: none; }
}
.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 6px 14px;
  border-bottom: 1px solid var(--ww-ink-hairline);
  margin-bottom: 14px;
}
.new-link {
  width: 28px; height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--ww-ink-faint);
  transition: color .25s ease, background-color .25s ease;
  &:hover {
    color: rgb(var(--ww-vermilion));
    background: rgb(var(--ww-vermilion) / .1);
  }
}
.tree {
  overflow-y: auto;
  flex: 1;
  padding-bottom: 24px;
  scrollbar-width: thin;
}
.roots { list-style: none; margin: 0; padding: 0; }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  padding: 12px 6px;
  font-size: 14px;
}
</style>
