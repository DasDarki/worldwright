<script setup lang="ts">
import type { EntitySummary, EntityType } from '~/types/api'

interface DragState {
  id: number
  parentId: number | null
}

const props = defineProps<{
  node: EntitySummary
  childrenOf: Map<number, EntitySummary[]>
  types: EntityType[]
  depth: number
  activeSlug: string
  dragState: DragState | null
  canDrag: boolean
}>()

const emit = defineEmits<{
  'drag-start': [payload: DragState]
  'drag-end': []
  'reorder': [payload: { draggedId: number; targetId: number; position: 'before' | 'after' }]
}>()

const expanded = ref(props.depth < 1)
const children = computed(() => props.childrenOf.get(props.node.id) || [])
const hasChildren = computed(() => children.value.length > 0)

const type = computed(() => props.types.find((t) => t.id === props.node.entity_type_id))
const isActive = computed(() => props.node.slug === props.activeSlug)

const dropPosition = ref<'before' | 'after' | null>(null)
const isBeingDragged = computed(() => props.dragState?.id === props.node.id)

function onDragStart(e: DragEvent) {
  if (!props.canDrag) return
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', String(props.node.id))
  }
  emit('drag-start', { id: props.node.id, parentId: props.node.parent_id ?? null })
}

function onDragEnd() {
  dropPosition.value = null
  emit('drag-end')
}

function onDragOver(e: DragEvent) {
  if (!props.dragState) return
  if (props.dragState.id === props.node.id) return
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  const target = e.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const y = e.clientY - rect.top
  dropPosition.value = y < rect.height / 2 ? 'before' : 'after'
}

function onDragLeave(e: DragEvent) {
  const related = e.relatedTarget as Node | null
  const current = e.currentTarget as HTMLElement
  if (!related || !current.contains(related)) {
    dropPosition.value = null
  }
}

function onDrop(e: DragEvent) {
  if (!props.dragState) return
  const draggedId = props.dragState.id
  const pos = dropPosition.value
  dropPosition.value = null
  if (!pos || draggedId === props.node.id) return
  e.preventDefault()
  e.stopPropagation()
  emit('reorder', { draggedId, targetId: props.node.id, position: pos })
}
</script>

<template>
  <li
    class="tree-node"
    :class="{
      active: isActive,
      dragging: isBeingDragged,
      'drop-before': dropPosition === 'before',
      'drop-after':  dropPosition === 'after',
    }"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <div
      class="row"
      :style="{ paddingLeft: `${depth * 14 + 6}px` }"
      :draggable="canDrag"
      @dragstart="onDragStart"
      @dragend="onDragEnd"
    >
      <button
        v-if="hasChildren"
        type="button"
        class="chev"
        :class="{ open: expanded }"
        :aria-label="expanded ? 'Collapse' : 'Expand'"
        @click="expanded = !expanded"
      >
        <svg viewBox="0 0 10 10" width="9" height="9" aria-hidden="true">
          <path d="M3 1 L7 5 L3 9" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <span v-else class="chev-spacer" aria-hidden="true" />

      <span
        v-if="canDrag"
        class="drag-handle"
        aria-hidden="true"
        :title="'Drag to reorder'"
      >
        <svg viewBox="0 0 10 14" width="8" height="12">
          <circle cx="3" cy="3" r="1.1" fill="currentColor"/>
          <circle cx="7" cy="3" r="1.1" fill="currentColor"/>
          <circle cx="3" cy="7" r="1.1" fill="currentColor"/>
          <circle cx="7" cy="7" r="1.1" fill="currentColor"/>
          <circle cx="3" cy="11" r="1.1" fill="currentColor"/>
          <circle cx="7" cy="11" r="1.1" fill="currentColor"/>
        </svg>
      </span>

      <NuxtLink :to="`/entities/${node.slug}`" class="title">
        <span v-if="type" class="type-dot" :style="{ background: type.color || 'currentColor' }" />
        <span class="text">{{ node.title }}</span>
      </NuxtLink>
    </div>

    <Transition name="branch">
      <ul v-if="expanded && hasChildren" class="children">
        <EntityTreeNode
          v-for="c in children"
          :key="c.id"
          :node="c"
          :children-of="childrenOf"
          :types="types"
          :depth="depth + 1"
          :active-slug="activeSlug"
          :drag-state="dragState"
          :can-drag="canDrag"
          @drag-start="(p) => emit('drag-start', p)"
          @drag-end="emit('drag-end')"
          @reorder="(p) => emit('reorder', p)"
        />
      </ul>
    </Transition>
  </li>
</template>

<style scoped lang="scss">
.tree-node { display: block; position: relative; }
.row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 6px 4px 0;
  transition: background-color .25s ease, opacity .2s ease;
}
.row:hover { background: rgb(var(--ww-gold) / .08); }
.tree-node.active > .row { background: rgb(var(--ww-vermilion) / .08); }
.tree-node.active > .row .text { color: rgb(var(--ww-vermilion)); }
.tree-node.dragging > .row { opacity: .4; }

.tree-node.drop-before::before,
.tree-node.drop-after::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 0 2px rgb(var(--ww-vermilion) / .15);
  pointer-events: none;
  z-index: 2;
}
.tree-node.drop-before::before { top: -1px; }
.tree-node.drop-after::after { bottom: -1px; }

.chev {
  width: 18px; height: 18px;
  display: inline-flex; align-items: center; justify-content: center;
  color: var(--ww-ink-faint);
  background: transparent;
  transition: transform .35s cubic-bezier(.22,1,.36,1), color .25s ease;
  flex-shrink: 0;
}
.chev:hover { color: rgb(var(--ww-gold-deep)); }
.chev.open { transform: rotate(90deg); color: rgb(var(--ww-gold-deep)); }
.chev-spacer { display: inline-block; width: 18px; height: 18px; flex-shrink: 0; }

.drag-handle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 18px;
  color: rgb(var(--ww-ink-faint) / .6);
  cursor: grab;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity .2s ease, color .2s ease;
}
.row:hover .drag-handle { opacity: 1; }
.drag-handle:hover { color: rgb(var(--ww-gold-deep)); }
.row:active .drag-handle { cursor: grabbing; }

.title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-family: 'EB Garamond', serif;
  font-size: 15px;
  line-height: 1.3;
  color: rgb(var(--ww-ink));
  transition: color .25s ease;
  flex: 1;
  min-width: 0;
}
.title:hover .text { color: rgb(var(--ww-gold-deep)); }
.text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.type-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
  opacity: .75;
}

.children {
  list-style: none;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.branch-enter-active, .branch-leave-active {
  transition: grid-template-rows .35s cubic-bezier(.22,1,.36,1), opacity .25s ease;
  display: grid;
  grid-template-rows: 1fr;
}
.branch-enter-from, .branch-leave-to {
  grid-template-rows: 0fr;
  opacity: 0;
}
.branch-enter-active > *, .branch-leave-active > * {
  overflow: hidden;
}
</style>
