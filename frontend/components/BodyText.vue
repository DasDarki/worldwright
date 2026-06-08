<script setup lang="ts">
// Recursive renderer for TipTap text nodes with marks. Each call peels off
// one mark, wraps the rest in the matching HTML element, then recurses.
// Unknown mark types are skipped silently so we degrade to plain text
// rather than dropping the run entirely.
//
// Without this, marks (bold / italic / strike / code) stored on text nodes
// by the editor are written into the DB but never rendered on the read view.

interface Mark {
  type: string
  attrs?: Record<string, unknown>
}

const props = defineProps<{ text: string; marks?: Mark[] }>()

const tagFor: Record<string, string> = {
  bold: 'strong',
  strong: 'strong',
  italic: 'em',
  em: 'em',
  strike: 's',
  strikethrough: 's',
  code: 'code',
  underline: 'u',
}

const ordered = computed<Mark[]>(() => (props.marks || []).filter((m) => tagFor[m.type]))
const head = computed<Mark | null>(() => ordered.value[0] ?? null)
const tag = computed(() => (head.value ? tagFor[head.value.type] : ''))
const rest = computed<Mark[]>(() => ordered.value.slice(1))
</script>

<template>
  <component :is="tag" v-if="tag">
    <BodyText :text="text" :marks="rest" />
  </component>
  <template v-else>{{ text }}</template>
</template>
