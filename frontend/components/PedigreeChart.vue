<script setup lang="ts">
import type { Genealogy, GenealogyEdge, GenealogyNode } from '~/types/api'

const props = withDefaults(defineProps<{
  genealogy: Genealogy
  focal: number
}>(), {})

interface Placement {
  node: GenealogyNode
  col: number
  gen: number
  px: number
  py: number
}

const COL_W = 220
const ROW_H = 150
const NODE_W = 200
const NODE_H = 90

function buildGenerations(g: Genealogy): Map<number, number> {
  const generation = new Map<number, number>()
  generation.set(g.focal, 0)
  const queue: number[] = [g.focal]
  while (queue.length) {
    const id = queue.shift()!
    const gen = generation.get(id)!
    for (const e of g.edges) {
      if (e.type === 'parent_of') {
        if (e.to === id && !generation.has(e.from)) {
          generation.set(e.from, gen - 1)
          queue.push(e.from)
        }
        if (e.from === id && !generation.has(e.to)) {
          generation.set(e.to, gen + 1)
          queue.push(e.to)
        }
      } else if (e.type === 'spouse_of' || e.type === 'sibling_of') {
        if (e.from === id && !generation.has(e.to)) {
          generation.set(e.to, gen)
          queue.push(e.to)
        }
        if (e.to === id && !generation.has(e.from)) {
          generation.set(e.from, gen)
          queue.push(e.from)
        }
      }
    }
  }
  for (const n of g.nodes) {
    if (!generation.has(n.id)) generation.set(n.id, 0)
  }
  return generation
}

function spouseOf(edges: GenealogyEdge[]): Map<number, Set<number>> {
  const m = new Map<number, Set<number>>()
  const add = (a: number, b: number) => {
    if (!m.has(a)) m.set(a, new Set())
    m.get(a)!.add(b)
  }
  for (const e of edges) {
    if (e.type === 'spouse_of') {
      add(e.from, e.to)
      add(e.to, e.from)
    }
  }
  return m
}

function sortGeneration(nodes: GenealogyNode[], spouses: Map<number, Set<number>>): GenealogyNode[] {
  const placed: GenealogyNode[] = []
  const remaining = new Set(nodes.map((n) => n.id))
  const byId = new Map(nodes.map((n) => [n.id, n] as const))
  while (remaining.size > 0) {
    const firstId = remaining.values().next().value!
    const first = byId.get(firstId)!
    placed.push(first)
    remaining.delete(firstId)
    const partners = spouses.get(firstId) ?? new Set()
    for (const p of partners) {
      if (remaining.has(p)) {
        placed.push(byId.get(p)!)
        remaining.delete(p)
      }
    }
  }
  return placed
}

const layout = computed(() => {
  const g = props.genealogy
  const generation = buildGenerations(g)
  const spouses = spouseOf(g.edges)
  const byGen = new Map<number, GenealogyNode[]>()
  for (const node of g.nodes) {
    const gen = generation.get(node.id) ?? 0
    if (!byGen.has(gen)) byGen.set(gen, [])
    byGen.get(gen)!.push(node)
  }
  const sortedGens = Array.from(byGen.keys()).sort((a, b) => a - b)
  const placements: Placement[] = []
  let maxCount = 0
  for (const gen of sortedGens) {
    const ordered = sortGeneration(byGen.get(gen)!, spouses)
    if (ordered.length > maxCount) maxCount = ordered.length
    const half = (ordered.length - 1) / 2
    ordered.forEach((node, i) => {
      const col = i - half
      placements.push({
        node,
        col,
        gen,
        px: col * COL_W,
        py: (gen - sortedGens[0]) * ROW_H,
      })
    })
  }
  const minX = Math.min(...placements.map((p) => p.px))
  const maxX = Math.max(...placements.map((p) => p.px))
  const width = (maxX - minX) + NODE_W + 40
  const height = (sortedGens.length - 1) * ROW_H + NODE_H + 40
  const offsetX = -minX + NODE_W / 2 + 20
  return { placements, width, height, offsetX, generations: sortedGens.length }
})

const placementById = computed(() => {
  const m = new Map<number, Placement>()
  for (const p of layout.value.placements) m.set(p.node.id, p)
  return m
})

interface DrawnEdge {
  type: string
  d: string
  active?: boolean
}

const drawnEdges = computed<DrawnEdge[]>(() => {
  const result: DrawnEdge[] = []
  const off = layout.value.offsetX
  for (const e of props.genealogy.edges) {
    const from = placementById.value.get(e.from)
    const to = placementById.value.get(e.to)
    if (!from || !to) continue
    if (e.type === 'parent_of') {
      const x1 = from.px + off
      const y1 = from.py + NODE_H / 2 + 20
      const x2 = to.px + off
      const y2 = to.py - NODE_H / 2 + 20
      const midY = (y1 + y2) / 2
      result.push({
        type: 'parent_of',
        d: `M ${x1} ${y1} C ${x1} ${midY}, ${x2} ${midY}, ${x2} ${y2}`,
      })
    } else if (e.type === 'spouse_of') {
      const x1 = from.px + off
      const x2 = to.px + off
      const y = from.py + 20
      result.push({
        type: 'spouse_of',
        d: `M ${x1} ${y} L ${x2} ${y}`,
      })
    } else if (e.type === 'sibling_of') {
      const x1 = from.px + off
      const x2 = to.px + off
      const y = from.py + 20
      const arch = 18
      const mid = (x1 + x2) / 2
      result.push({
        type: 'sibling_of',
        d: `M ${x1} ${y} Q ${mid} ${y - arch}, ${x2} ${y}`,
      })
    }
  }
  return result
})

const hasContent = computed(() => props.genealogy.nodes.length > 1)
</script>

<template>
  <div v-if="hasContent" class="pedigree-wrap">
    <Ornament>Lineage</Ornament>

    <div class="scroll">
      <div
        class="canvas"
        :style="{ width: layout.width + 'px', height: layout.height + 'px' }"
      >
        <svg
          class="lines"
          :viewBox="`0 0 ${layout.width} ${layout.height}`"
          :width="layout.width"
          :height="layout.height"
        >
          <path
            v-for="(edge, i) in drawnEdges"
            :key="`edge-${i}`"
            :d="edge.d"
            :class="['line', edge.type]"
            fill="none"
          />
        </svg>

        <div
          v-for="(p, i) in layout.placements"
          :key="p.node.id"
          :class="['node', { focal: p.node.id === focal }]"
          :style="{
            left: `${p.px + layout.offsetX - NODE_W / 2}px`,
            top: `${p.py}px`,
            width: `${NODE_W}px`,
            animationDelay: `${0.06 * i}s`,
          }"
        >
          <NuxtLink
            :to="`/entities/${p.node.slug}`"
            class="card"
            :title="p.node.summary || p.node.title"
          >
            <span v-if="p.node.id === focal" class="focal-mark" aria-hidden="true" />
            <span class="ww-label tier">Gen {{ p.gen }}</span>
            <span class="name">{{ p.node.title }}</span>
            <span v-if="p.node.summary" class="summary">{{ p.node.summary }}</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.pedigree-wrap {
  margin-top: 30px;
}
.scroll {
  margin-top: 22px;
  overflow-x: auto;
  padding: 12px 4px 24px;
}
.canvas {
  position: relative;
  margin: 0 auto;
}

.lines {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}
.line {
  stroke-width: 1.4;
  fill: none;
  stroke-linecap: round;
}
.line.parent_of {
  stroke: rgb(var(--ww-ink) / .55);
  stroke-dasharray: 0;
}
.line.spouse_of {
  stroke: rgb(var(--ww-vermilion));
  stroke-width: 1.8;
}
.line.sibling_of {
  stroke: rgb(var(--ww-gold-deep));
  stroke-dasharray: 3 4;
}

.node {
  position: absolute;
  height: 90px;
  opacity: 0;
  animation: rise .9s cubic-bezier(.22,1,.36,1) forwards;
  z-index: 1;
}
@keyframes rise {
  from { opacity: 0; transform: translateY(8px) scale(.96); }
  to   { opacity: 1; transform: none; }
}

.card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  height: 100%;
  padding: 10px 14px;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 18px 30px -22px rgb(0 0 0 / .4);
  transition: transform .35s cubic-bezier(.22,1,.36,1), border-color .35s ease, box-shadow .35s ease;
  position: relative;
  overflow: hidden;
}
.card:hover {
  transform: translateY(-2px);
  border-color: rgb(var(--ww-gold));
  box-shadow: 0 26px 38px -24px rgb(0 0 0 / .55);
}

.node.focal .card {
  border-color: rgb(var(--ww-vermilion));
  background: rgb(var(--ww-vermilion) / .05);
}

.focal-mark {
  position: absolute;
  top: -1px; right: -1px; bottom: -1px;
  width: 4px;
  background: rgb(var(--ww-vermilion));
  animation: focal-pulse 2.4s ease-in-out infinite;
}
@keyframes focal-pulse {
  0%, 100% { opacity: .6; }
  50% { opacity: 1; box-shadow: -4px 0 14px rgb(var(--ww-vermilion) / .5); }
}

.tier {
  font-size: 9px;
  letter-spacing: .26em;
  color: var(--ww-ink-faint);
}
.node.focal .tier { color: rgb(var(--ww-vermilion)); }

.name {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 24, "wght" 500;
  font-size: 16px;
  line-height: 1.15;
  letter-spacing: -0.01em;
  color: rgb(var(--ww-ink));
}
.summary {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  font-size: 11px;
  line-height: 1.3;
  color: var(--ww-ink-faint);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
