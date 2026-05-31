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
  router.push(`/entities/${slug}`)
}

// Optional zoom: scale to container width
const wrapRef = ref<HTMLElement | null>(null)
const scale = ref(1)
function recompute() {
  if (!wrapRef.value) return
  const w = wrapRef.value.clientWidth - 4
  const gw = layout.value.width
  if (gw <= 0) return
  scale.value = gw > w ? w / gw : 1
}
onMounted(() => {
  recompute()
  if (typeof window !== 'undefined') window.addEventListener('resize', recompute)
})
onBeforeUnmount(() => {
  if (typeof window !== 'undefined') window.removeEventListener('resize', recompute)
})
watch(() => layout.value.width, recompute)
</script>

<template>
  <figure ref="wrapRef" class="ww-graph" :class="{ embedded }">
    <figcaption v-if="title" class="caption">{{ title }}</figcaption>
    <div v-if="error" class="empty">{{ error }}</div>
    <div v-else-if="!data || !data.nodes.length" class="empty">{{ t('graph.empty') }}</div>
    <div v-else class="canvas" :style="{ height: (layout.height * scale + 8) + 'px' }">
      <svg
        :viewBox="`0 0 ${layout.width} ${layout.height}`"
        :width="layout.width * scale"
        :height="layout.height * scale"
        preserveAspectRatio="xMidYMid meet"
      >
        <defs>
          <marker id="arrowhead" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
            <path d="M0,0 L10,5 L0,10 Z" fill="currentColor" />
          </marker>
        </defs>

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
  display: flex;
  justify-content: center;
}

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
