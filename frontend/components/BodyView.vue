<script setup lang="ts">
import type { TipTapNode } from '~/types/api'

const props = defineProps<{ body: unknown }>()

const doc = computed<TipTapNode | null>(() => {
  try {
    if (!props.body) return null
    if (typeof props.body === 'string') return JSON.parse(props.body) as TipTapNode
    return props.body as TipTapNode
  } catch {
    return null
  }
})
</script>

<template>
  <div v-if="doc" class="ww-body">
    <BodyNode :node="doc" />
  </div>
</template>

<style scoped lang="scss">
.ww-body {
  font-family: 'EB Garamond', serif;
  font-size: 18px;
  line-height: 1.65;
  color: rgb(var(--ww-ink-shade));
  overflow-wrap: anywhere;
  min-width: 0;

  :deep(p) { margin: 0 0 1em; }
  :deep(img), :deep(svg), :deep(.ww-img-embed) {
    max-width: 100%;
    height: auto;
  }
  :deep(h2), :deep(h3) {
    font-family: 'Fraunces', serif;
    font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
    margin: 1.4em 0 .5em;
    letter-spacing: -0.02em;
  }
  :deep(blockquote) {
    border-left: 3px solid rgb(var(--ww-gold));
    margin: 1.2em 0;
    padding: .2em 0 .2em 1em;
    font-style: italic;
    color: rgb(var(--ww-ink-shade));
  }
  :deep(ul), :deep(ol) {
    padding-left: 1.6em;
    margin: 0 0 1em;
  }
  :deep(ul) { list-style: disc outside; }
  :deep(ul ul) { list-style: circle outside; }
  :deep(ul ul ul) { list-style: square outside; }
  :deep(ol) { list-style: decimal outside; }
  :deep(ol ol) { list-style: lower-alpha outside; }
  :deep(ol ol ol) { list-style: lower-roman outside; }
  :deep(li) {
    margin: .3em 0;
    line-height: 1.55;
  }
  :deep(li::marker) {
    color: rgb(var(--ww-gold-deep));
    font-family: 'Fraunces', serif;
  }

  /* Horizontal rule */
  :deep(.ww-hr) {
    border: 0;
    height: 1px;
    background: linear-gradient(to right, transparent, var(--ww-ink-hairline) 8%, var(--ww-ink-hairline) 92%, transparent);
    margin: 2.2em 0;
  }

  /* Code blocks */
  :deep(pre.ww-code) {
    background: rgb(var(--ww-ink-shade) / .96);
    color: rgb(var(--ww-parchment));
    border: 1px solid rgb(var(--ww-ink) / .8);
    padding: 18px 22px;
    margin: 1.6em 0;
    overflow-x: auto;
    font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, monospace;
    font-size: 14px;
    line-height: 1.55;
    border-radius: 0;
  }
  :deep(pre.ww-code code) {
    font-family: inherit;
    background: transparent;
    color: inherit;
    padding: 0;
    font-size: inherit;
  }
  :deep(code) {
    font-family: 'JetBrains Mono', ui-monospace, monospace;
    background: rgb(var(--ww-ink-hairline));
    padding: 1px 6px;
    font-size: .9em;
    border-radius: 0;
  }

  /* Tables */
  :deep(.ww-table-wrap) {
    overflow-x: auto;
    margin: 1.6em 0;
    border: 1px solid var(--ww-ink-hairline);
  }
  :deep(.ww-table) {
    width: 100%;
    border-collapse: collapse;
    font-family: 'EB Garamond', serif;
    font-size: 15px;
  }
  :deep(.ww-table th), :deep(.ww-table td) {
    border: 1px solid var(--ww-ink-hairline);
    padding: 8px 12px;
    text-align: left;
    vertical-align: top;
    min-width: 80px;
  }
  :deep(.ww-table th) {
    background: rgb(var(--ww-parchment-stain) / .5);
    font-family: 'Cormorant SC', serif;
    font-size: 11px;
    letter-spacing: .2em;
    text-transform: uppercase;
    color: rgb(var(--ww-gold-deep));
    font-weight: 500;
  }
  :deep(.ww-table tr:nth-child(even) td) {
    background: rgb(var(--ww-parchment-stain) / .15);
  }

  /* Callouts */
  :deep(.ww-callout) {
    margin: 1.6em 0;
    padding: 16px 20px 16px 56px;
    border: 1px solid var(--ww-ink-hairline);
    border-left-width: 4px;
    position: relative;
    background: rgb(var(--ww-parchment-stain) / .25);
    font-size: 16px;
  }
  :deep(.ww-callout) :deep(p:last-child),
  :deep(.ww-callout > p:last-child) { margin-bottom: 0; }
  :deep(.ww-callout::before) {
    content: '';
    position: absolute;
    left: 18px;
    top: 18px;
    width: 22px;
    height: 22px;
    background-position: center;
    background-repeat: no-repeat;
    background-size: contain;
    opacity: .85;
  }
  :deep(.ww-callout[data-variant="info"]) {
    border-left-color: rgb(var(--ww-ink-shade));
  }
  :deep(.ww-callout[data-variant="info"]::before) {
    background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23173842' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><circle cx='12' cy='12' r='9'/><path d='M12 8 V13 M12 16.5 V17'/></svg>");
  }
  :deep(.ww-callout[data-variant="warn"]) {
    border-left-color: rgb(var(--ww-vermilion));
    background: rgb(var(--ww-vermilion) / .06);
  }
  :deep(.ww-callout[data-variant="warn"]::before) {
    background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23a8442e' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M12 3 L22 20 L2 20 Z'/><path d='M12 10 V14 M12 17 V17.5'/></svg>");
  }
  :deep(.ww-callout[data-variant="note"]) {
    border-left-color: rgb(var(--ww-gold-deep));
    background: rgb(var(--ww-gold) / .08);
  }
  :deep(.ww-callout[data-variant="note"]::before) {
    background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%237c5e30' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M5 4 H17 L20 7 V20 H5 Z'/><path d='M17 4 V7 H20 M8 11 H16 M8 15 H14'/></svg>");
  }
  :deep(.ww-callout[data-variant="lore"]) {
    border-left-color: rgb(var(--ww-gold));
    background: linear-gradient(135deg, rgb(var(--ww-gold) / .12), rgb(var(--ww-vermilion) / .06));
    font-style: italic;
  }
  :deep(.ww-callout[data-variant="lore"]::before) {
    background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23b8935a' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M12 3 L13.5 10.5 L21 12 L13.5 13.5 L12 21 L10.5 13.5 L3 12 L10.5 10.5 Z'/></svg>");
  }
}
</style>
