<script setup lang="ts">
import type { TipTapNode } from '~/types/api'

defineProps<{ node: TipTapNode }>()
</script>

<template>
  <template v-if="node.type === 'doc'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </template>

  <p v-else-if="node.type === 'paragraph'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </p>

  <h2 v-else-if="node.type === 'heading' && (node.attrs?.level === 1 || node.attrs?.level === 2)">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </h2>
  <h3 v-else-if="node.type === 'heading'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </h3>

  <blockquote v-else-if="node.type === 'blockquote'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </blockquote>

  <ul v-else-if="node.type === 'bulletList'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </ul>
  <ol v-else-if="node.type === 'orderedList'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </ol>
  <li v-else-if="node.type === 'listItem'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </li>

  <SpoilerBlock v-else-if="node.type === 'spoiler'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </SpoilerBlock>

  <aside
    v-else-if="node.type === 'callout'"
    class="ww-callout"
    :data-variant="(node.attrs?.variant as string) || 'info'"
  >
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </aside>

  <pre v-else-if="node.type === 'codeBlock'" class="ww-code"><code><BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" /></code></pre>

  <div v-else-if="node.type === 'table'" class="ww-table-wrap">
    <table class="ww-table">
      <tbody>
        <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
      </tbody>
    </table>
  </div>
  <tr v-else-if="node.type === 'tableRow'">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </tr>
  <th v-else-if="node.type === 'tableHeader'" :colspan="Number(node.attrs?.colspan) || 1" :rowspan="Number(node.attrs?.rowspan) || 1">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </th>
  <td v-else-if="node.type === 'tableCell'" :colspan="Number(node.attrs?.colspan) || 1" :rowspan="Number(node.attrs?.rowspan) || 1">
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </td>

  <hr v-else-if="node.type === 'horizontalRule'" class="ww-hr" />

  <img
    v-else-if="node.type === 'wwimage' && node.attrs?.src"
    :src="(node.attrs.src as string)"
    :alt="(node.attrs.alt as string) || ''"
    class="ww-img-embed"
    loading="lazy"
  />

  <RelationshipGraph
    v-else-if="node.type === 'relationshipGraph'"
    :entity-ids="(node.attrs?.entityIds as number[]) || []"
    embedded
  />

  <WikilinkPreview
    v-else-if="node.type === 'wikilink' && node.attrs?.slug"
    :slug="(node.attrs.slug as string)"
    :label="(node.attrs.label as string) || ''"
  />

  <template v-else-if="node.type === 'text'">{{ node.text }}</template>

  <br v-else-if="node.type === 'hardBreak'" />

  <template v-else>
    <BodyNode v-for="(c, i) in node.content || []" :key="i" :node="c" />
  </template>
</template>

<style>
.ww-body .ww-img-embed {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 1.4em auto;
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 30px 60px -30px rgb(0 0 0 / .35);
}
</style>
