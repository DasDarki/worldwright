<script setup lang="ts">
import * as dagre from '@dagrejs/dagre'

interface GraphNode {
  id: number
  entity_type_id: number
  title: string
  slug: string
  summary?: string
}
interface GraphEdge {
  from: number
  to: number
  type_key: string
  type_label: string
  inverse_label?: string
  is_symmetric: boolean
  category: string
  description?: string
}
interface GraphData {
  nodes: GraphNode[]
  edges: GraphEdge[]
}

const props = withDefaults(defineProps<{
  entityIds: number[]
  embedded?: boolean
  title?: string
}>(), { embedded: false })

const { $api } = useNuxtApp()
const { locale, t } = useI18n()

const data = ref<GraphData | null>(null)
const error = ref<string | null>(null)

const ids = computed(() => Array.from(new Set(props.entityIds.filter((n) => Number.isFinite(n) && n > 0))))

async function load() {
  if (!ids.value.length) {
    data.value = { nodes: [], edges: [] }
    return
  }
  error.value = null
  try {
    const res = await $api<GraphData>('/entities/relationship-graph', {
      method: 'POST',
      body: { ids: ids.value, lang: locale.value },
    })
    data.value = res
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to load graph'
    data.value = { nodes: [], edges: [] }
  }
}

watch(ids, load, { immediate: true })

const NODE_W = 180
const NODE_H = 72
const PAD = 20

interface LaidOutNode extends GraphNode { x: number; y: number }
interface LaidOutEdge { points: { x: number; y: number }[]; edge: GraphEdge }

const layout = computed(() => {
  const d = data.value
  if (!d || !d.nodes.length) return { width: 0, height: 0, nodes: [] as LaidOutNode[], edges: [] as LaidOutEdge[] }

  const g = new dagre.graphlib.Graph({ multigraph: true })
  g.setGraph({
    rankdir: 'TB',
    ranksep: 70,
    nodesep: 36,
    edgesep: 14,
    marginx: PAD,
    marginy: PAD,
  })
  g.setDefaultEdgeLabel(() => ({}))

  for (const n of d.nodes) {
    g.setNode(String(n.id), { width: NODE_W, height: NODE_H, label: n.title })
  }
  // dedupe symmetric edges to avoid drawing two lines between same pair
  const symSeen = new Set<string>()
  d.edges.forEach((e, idx) => {
    if (e.is_symmetric) {
      const key = [Math.min(e.from, e.to), Math.max(e.from, e.to), e.type_key].join('|')
      if (symSeen.has(key)) return
      symSeen.add(key)
    }
    g.setEdge(String(e.from), String(e.to), { label: e.type_label, edge: e }, `e-${idx}`)
  })

  dagre.layout(g)

  const laidOutNodes: LaidOutNode[] = d.nodes.map((n) => {
    const p = g.node(String(n.id))
    return { ...n, x: p.x, y: p.y }
  })
  const laidOutEdges: LaidOutEdge[] = g.edges().map((eId) => {
    const meta = g.edge(eId) as any
    return { points: meta.points || [], edge: meta.edge as GraphEdge }
  })

  const graphInfo = g.graph()
  return {
    width: Math.max(graphInfo.width || 0, 200),
    height: Math.max(graphInfo.height || 0, 100),
    nodes: laidOutNodes,
    edges: laidOutEdges,
  }
})

function edgePath(points: { x: number; y: number }[]): string {
  if (points.length === 0) return ''
  let d = `M ${points[0].x} ${points[0].y}`
  for (let i = 1; i < points.length - 1; i++) {
    const c = points[i]
    const next = points[i + 1]
    const mx = (c.x + next.x) / 2
    const my = (c.y + next.y) / 2
    d += ` Q ${c.x} ${c.y} ${mx} ${my}`
  }
  const last = points[points.length - 1]
  d += ` T ${last.x} ${last.y}`
  return d
}

function edgeClass(edge: GraphEdge): string {
  if (edge.category === 'genealogy') {
    return `e-genealogy e-${edge.type_key}`
  }
  return `e-generic e-${edge.type_key}`
}

function midpoint(points: { x: number; y: number }[]): { x: number; y: number } {
  if (!points.length) return { x: 0, y: 0 }
  const mid = points[Math.floor(points.length / 2)]
  return mid
}

const router = useRouter()
function goToNode(slug: string) {
  if (didDrag.value) return // a pan ended on this node — don't navigate
  router.push(`/entities/${slug}`)
}

// --- Pan & zoom -----------------------------------------------------------
//
// We render the SVG at the container's pixel size and apply pan/zoom by
// transforming an inner <g>. The user starts at a "fit-to-view" transform
// (graph centered + scaled to fit); drag-on-background pans, mousewheel
// zooms around the cursor. A reset button restores the fit.

const wrapRef = ref<HTMLElement | null>(null)
const canvasRef = ref<HTMLElement | null>(null)

const viewW = ref(0)         // container width in CSS px
const viewH = ref(360)       // container height in CSS px (fixed unless content needs more)

const tx = ref(0)
const ty = ref(0)
const scale = ref(1)
const userInteracted = ref(false)
const didDrag = ref(false)

function fitView() {
  if (!canvasRef.value) return
  const w = canvasRef.value.clientWidth
  const h = viewH.value
  const gw = layout.value.width
  const gh = layout.value.height
  if (gw <= 0 || gh <= 0 || w <= 0) return
  const sx = (w - 16) / gw
  const sy = (h - 16) / gh
  const s = Math.min(sx, sy, 1)
  scale.value = Math.max(s, 0.2)
  tx.value = (w - gw * scale.value) / 2
  ty.value = (h - gh * scale.value) / 2
}

function setViewSize() {
  if (canvasRef.value) viewW.value = canvasRef.value.clientWidth
  // Adapt height to graph if it's small; cap at 720 for huge graphs so the
  // user always has a window-sized viewport rather than a 3000px tall canvas.
  const gh = layout.value.height
  viewH.value = Math.min(Math.max(gh + 40, 280), 720)
}

function resetView() {
  userInteracted.value = false
  setViewSize()
  fitView()
}

onMounted(() => {
  setViewSize()
  fitView()
  if (typeof window !== 'undefined') window.addEventListener('resize', onResize)
})
onBeforeUnmount(() => {
  if (typeof window !== 'undefined') window.removeEventListener('resize', onResize)
})

function onResize() {
  setViewSize()
  if (!userInteracted.value) fitView()
}

watch(() => layout.value.width, () => {
  setViewSize()
  if (!userInteracted.value) fitView()
})

// Wheel zoom anchored at cursor position.
const MIN_SCALE = 0.2
const MAX_SCALE = 4

function onWheel(e: WheelEvent) {
  if (!canvasRef.value) return
  e.preventDefault()
  userInteracted.value = true
  const rect = canvasRef.value.getBoundingClientRect()
  const cx = e.clientX - rect.left
  const cy = e.clientY - rect.top
  // factor: invert sign so up = zoom in
  const factor = e.deltaY < 0 ? 1.12 : 1 / 1.12
  const newScale = Math.min(MAX_SCALE, Math.max(MIN_SCALE, scale.value * factor))
  if (newScale === scale.value) return
  // keep the world-space point under the cursor stationary
  const k = newScale / scale.value
  tx.value = cx - (cx - tx.value) * k
  ty.value = cy - (cy - ty.value) * k
  scale.value = newScale
}

// Pan with primary button drag on the canvas background.
let dragStartX = 0
let dragStartY = 0
let dragOriginX = 0
let dragOriginY = 0
const panning = ref(false)

function onMouseDown(e: MouseEvent) {
  if (e.button !== 0) return
  // Only start a pan if the target is the SVG background, not a node group
  // — let nodes own clicks for navigation.
  const target = e.target as Element | null
  if (target && target.closest('.node')) return
  panning.value = true
  didDrag.value = false
  dragStartX = e.clientX
  dragStartY = e.clientY
  dragOriginX = tx.value
  dragOriginY = ty.value
  if (typeof window !== 'undefined') {
    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  }
}
function onMouseMove(e: MouseEvent) {
  if (!panning.value) return
  const dx = e.clientX - dragStartX
  const dy = e.clientY - dragStartY
  if (!didDrag.value && (Math.abs(dx) > 3 || Math.abs(dy) > 3)) {
    didDrag.value = true
    userInteracted.value = true
  }
  if (didDrag.value) {
    tx.value = dragOriginX + dx
    ty.value = dragOriginY + dy
  }
}
function onMouseUp() {
  panning.value = false
  if (typeof window !== 'undefined') {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }
  // didDrag is consumed by the next node click; clear shortly after
  setTimeout(() => { didDrag.value = false }, 0)
}

function zoomBy(factor: number) {
  if (!canvasRef.value) return
  userInteracted.value = true
  const w = canvasRef.value.clientWidth
  const h = viewH.value
  const newScale = Math.min(MAX_SCALE, Math.max(MIN_SCALE, scale.value * factor))
  if (newScale === scale.value) return
  const k = newScale / scale.value
  // zoom around viewport center
  tx.value = w / 2 - (w / 2 - tx.value) * k
  ty.value = h / 2 - (h / 2 - ty.value) * k
  scale.value = newScale
}
</script>

<template>
  <figure ref="wrapRef" class="ww-graph" :class="{ embedded }">
    <figcaption v-if="title" class="caption">{{ title }}</figcaption>
    <div v-if="error" class="empty">{{ error }}</div>
    <div v-else-if="!data || !data.nodes.length" class="empty">{{ t('graph.empty') }}</div>
    <div
      v-else
      ref="canvasRef"
      class="canvas"
      :class="{ panning }"
      :style="{ height: viewH + 'px' }"
      @wheel="onWheel"
      @mousedown="onMouseDown"
    >
      <div class="controls" @mousedown.stop>
        <button type="button" class="ctl" :title="t('graph.zoomIn')" @click="zoomBy(1.2)">+</button>
        <button type="button" class="ctl" :title="t('graph.zoomOut')" @click="zoomBy(1 / 1.2)">−</button>
        <button type="button" class="ctl reset" :title="t('graph.reset')" @click="resetView">
          <svg viewBox="0 0 16 16" width="12" height="12" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
            <path d="M2 4 H6 V8"/>
            <path d="M2 4 A6 6 0 1 1 8 14"/>
          </svg>
        </button>
      </div>
      <svg
        class="graph-svg"
        :width="viewW || 0"
        :height="viewH"
        :viewBox="`0 0 ${Math.max(viewW, 1)} ${viewH}`"
        preserveAspectRatio="xMinYMin meet"
      >
        <defs>
          <marker id="arrowhead" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
            <path d="M0,0 L10,5 L0,10 Z" fill="currentColor" />
          </marker>
        </defs>

        <g class="viewport" :transform="`translate(${tx} ${ty}) scale(${scale})`">
        <g class="edges">
          <g v-for="(e, i) in layout.edges" :key="`edge-${i}`" :class="['edge', edgeClass(e.edge)]">
            <path
              :d="edgePath(e.points)"
              fill="none"
              :marker-end="e.edge.is_symmetric ? '' : 'url(#arrowhead)'"
            />
            <text
              :x="midpoint(e.points).x"
              :y="midpoint(e.points).y - 6"
              text-anchor="middle"
              class="edge-label"
            >{{ e.edge.type_label }}</text>
          </g>
        </g>

        <g class="nodes">
          <g
            v-for="n in layout.nodes"
            :key="n.id"
            :transform="`translate(${n.x - NODE_W / 2}, ${n.y - NODE_H / 2})`"
            class="node"
            :class="`type-${n.entity_type_id}`"
            role="link"
            :aria-label="n.title"
            tabindex="0"
            @click="goToNode(n.slug)"
            @keydown.enter="goToNode(n.slug)"
          >
            <rect :width="NODE_W" :height="NODE_H" rx="2" />
            <text :x="NODE_W / 2" :y="28" text-anchor="middle" class="node-title">{{ n.title }}</text>
            <text v-if="n.summary" :x="NODE_W / 2" :y="50" text-anchor="middle" class="node-sub">{{ n.summary.length > 32 ? n.summary.slice(0, 30) + '…' : n.summary }}</text>
          </g>
        </g>
        </g>
      </svg>
    </div>
  </figure>
</template>

<style scoped lang="scss">
.ww-graph {
  margin: 1.6em 0;
  padding: 18px 16px 22px;
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment-stain) / .2);
  overflow: hidden;
}
.ww-graph.embedded { margin: 1.6em 0; }
.caption {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .32em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
  text-align: center;
  margin-bottom: 10px;
}
.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 24px 0;
  font-size: 14px;
}
.canvas {
  position: relative;
  width: 100%;
  background:
    radial-gradient(circle at 1px 1px, rgb(var(--ww-ink) / .06) 1px, transparent 1.5px) 0 0 / 22px 22px;
  border: 1px solid var(--ww-ink-hairline);
  cursor: grab;
  user-select: none;
  overflow: hidden;
  touch-action: none;
}
.canvas.panning { cursor: grabbing; }
.graph-svg { display: block; }

.controls {
  position: absolute;
  top: 10px; right: 10px;
  z-index: 2;
  display: flex;
  gap: 4px;
  background: var(--ww-card-bg);
  border: 1px solid var(--ww-ink-hairline);
  padding: 3px;
  box-shadow: 0 6px 14px -8px rgb(0 0 0 / .25);
}
.ctl {
  width: 26px;
  height: 26px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  background: transparent;
  font-family: 'Fraunces', serif;
  font-size: 16px;
  color: rgb(var(--ww-ink));
  cursor: pointer;
  transition: background-color .2s ease, color .2s ease;
}
.ctl:hover { background: rgb(var(--ww-gold) / .15); color: rgb(var(--ww-gold-deep)); }
.ctl.reset svg { display: block; }

:deep(.node) {
  cursor: pointer;
  rect {
    fill: rgb(var(--ww-card-bg));
    stroke: rgb(var(--ww-ink-hairline));
    stroke-width: 1;
    transition: fill .25s ease, stroke .25s ease;
  }
  &:hover rect {
    stroke: rgb(var(--ww-gold));
    fill: rgb(var(--ww-gold) / .1);
  }
  .node-title {
    font-family: 'Fraunces', serif;
    font-size: 14px;
    font-weight: 500;
    fill: rgb(var(--ww-ink));
  }
  .node-sub {
    font-family: 'EB Garamond', serif;
    font-size: 11px;
    font-style: italic;
    fill: var(--ww-ink-faint);
  }
}

:deep(.edge) {
  color: rgb(var(--ww-ink) / .55);
  path { stroke: currentColor; stroke-width: 1.4; }
  .edge-label {
    font-family: 'Cormorant SC', serif;
    font-size: 9px;
    letter-spacing: .18em;
    text-transform: uppercase;
    fill: var(--ww-ink-faint);
  }
}
:deep(.edge.e-genealogy)        { color: rgb(var(--ww-gold-deep)); }
:deep(.edge.e-spouse_of)        { color: rgb(var(--ww-vermilion)); path { stroke-width: 1.8; } }
:deep(.edge.e-sibling_of)       { color: rgb(var(--ww-gold)); path { stroke-dasharray: 3 4; } }
:deep(.edge.e-parent_of)        { color: rgb(var(--ww-ink) / .65); }
:deep(.edge.e-enemy_of)         { color: rgb(var(--ww-vermilion)); path { stroke-dasharray: 5 3; } }
:deep(.edge.e-ally_of)          { color: rgb(var(--ww-gold-deep)); }
</style>
