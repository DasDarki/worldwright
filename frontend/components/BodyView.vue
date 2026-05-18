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

  :deep(p) { margin: 0 0 1em; }
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
  :deep(ul), :deep(ol) { padding-left: 1.4em; }
  :deep(li) { margin: .25em 0; }
}
</style>
