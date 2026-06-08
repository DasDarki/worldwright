<script setup lang="ts">
import type { Editor } from '@tiptap/core'

const props = defineProps<{ editor: Editor | null }>()
const emit = defineEmits<{ wikilink: []; image: []; graph: []; vault: [] }>()

function is(name: string, attrs?: Record<string, any>) {
  return props.editor?.isActive(name, attrs) ?? false
}

function insertTable() {
  props.editor?.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
}
function toggleCodeBlock() {
  props.editor?.chain().focus().toggleCodeBlock().run()
}
function toggleCallout() {
  if (is('callout')) {
    props.editor?.chain().focus().cycleCalloutVariant().run()
  } else {
    props.editor?.chain().focus().toggleCallout({ variant: 'info' }).run()
  }
}

const HEADING_LEVELS = [1, 2, 3, 4, 5, 6] as const
const currentHeadingLevel = computed<number | null>(() => {
  const e = props.editor
  if (!e) return null
  for (const lvl of HEADING_LEVELS) {
    if (e.isActive('heading', { level: lvl })) return lvl
  }
  return null
})
function setHeading(level: number | null) {
  const e = props.editor
  if (!e) return
  if (level == null) e.chain().focus().setParagraph().run()
  else e.chain().focus().toggleHeading({ level: level as 1 | 2 | 3 | 4 | 5 | 6 }).run()
}
function onHeadingChange(ev: Event) {
  const target = ev.target as HTMLSelectElement
  const v = target.value
  if (v === 'p') setHeading(null)
  else setHeading(Number(v))
}
</script>

<template>
  <div v-if="editor" class="toolbar">
    <button type="button" class="t-btn" :class="{ on: is('bold') }" @click="editor.chain().focus().toggleBold().run()" aria-label="Bold">
      <b>B</b>
    </button>
    <button type="button" class="t-btn italic" :class="{ on: is('italic') }" @click="editor.chain().focus().toggleItalic().run()" aria-label="Italic">
      <i>I</i>
    </button>
    <span class="sep" />
    <select
      class="t-select"
      :value="currentHeadingLevel != null ? String(currentHeadingLevel) : 'p'"
      aria-label="Heading level"
      @change="onHeadingChange"
    >
      <option value="p">¶</option>
      <option value="1">H1</option>
      <option value="2">H2</option>
      <option value="3">H3</option>
      <option value="4">H4</option>
      <option value="5">H5</option>
      <option value="6">H6</option>
    </select>
    <span class="sep" />
    <button type="button" class="t-btn" :class="{ on: is('blockquote') }" @click="editor.chain().focus().toggleBlockquote().run()" aria-label="Quote">&ldquo;</button>
    <button type="button" class="t-btn" :class="{ on: is('bulletList') }" @click="editor.chain().focus().toggleBulletList().run()" aria-label="Bullet list">•</button>
    <button type="button" class="t-btn" :class="{ on: is('orderedList') }" @click="editor.chain().focus().toggleOrderedList().run()" aria-label="Ordered list">1.</button>
    <span class="sep" />
    <button type="button" class="t-btn" :class="{ on: is('spoiler') }" @click="editor.chain().focus().toggleSpoiler().run()" aria-label="Spoiler block">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
        <path d="M2 12 C 5 7, 9 5, 12 5 C 15 5, 19 7, 22 12 C 19 17, 15 19, 12 19 C 9 19, 5 17, 2 12 Z"/>
        <circle cx="12" cy="12" r="3"/>
        <path d="M5 5 L 19 19"/>
      </svg>
    </button>
    <button type="button" class="t-btn" @click="emit('image')" aria-label="Insert image">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="5" width="18" height="14" />
        <circle cx="9" cy="10" r="1.5" fill="currentColor" stroke="none"/>
        <path d="M3 17 L9 12 L14 16 L18 13 L21 16"/>
      </svg>
    </button>
    <button type="button" class="t-btn wikilink" @click="emit('wikilink')" aria-label="Insert wikilink">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
        <path d="M10 13 a5 5 0 0 0 7 0 l3-3 a5 5 0 0 0-7-7 l-1.5 1.5"/>
        <path d="M14 11 a5 5 0 0 0-7 0 l-3 3 a5 5 0 0 0 7 7 l1.5-1.5"/>
      </svg>
      [[…]]
    </button>

    <span class="sep" />

    <button type="button" class="t-btn" :class="{ on: is('codeBlock') }" @click="toggleCodeBlock" aria-label="Code block">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
        <path d="M9 7 L4 12 L9 17 M15 7 L20 12 L15 17"/>
      </svg>
    </button>
    <button type="button" class="t-btn" @click="insertTable" aria-label="Insert table">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6">
        <rect x="3" y="5" width="18" height="14"/>
        <path d="M3 10 H21 M3 14.5 H21 M8.5 5 V19 M14.5 5 V19"/>
      </svg>
    </button>
    <button type="button" class="t-btn" :class="{ on: is('callout') }" @click="toggleCallout" aria-label="Callout block">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="9"/>
        <path d="M12 8 V13 M12 16.5 V17"/>
      </svg>
    </button>
    <button type="button" class="t-btn graph" @click="emit('graph')" aria-label="Relationship graph widget">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round">
        <circle cx="5" cy="6" r="2"/>
        <circle cx="5" cy="18" r="2"/>
        <circle cx="19" cy="12" r="2"/>
        <path d="M7 6 L17 11 M7 18 L17 13"/>
      </svg>
    </button>
    <button type="button" class="t-btn vault" @click="emit('vault')" aria-label="Sealed section (author-only)">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
        <rect x="4" y="11" width="16" height="10" rx="1"/>
        <path d="M8 11 V7 a4 4 0 0 1 8 0 V11"/>
        <circle cx="12" cy="16" r="1.4" fill="currentColor"/>
      </svg>
    </button>
  </div>
</template>

<style scoped lang="scss">
.toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px;
  border: 1px solid var(--ww-ink-hairline);
  border-bottom: 0;
  background: rgb(var(--ww-parchment) / .6);
  position: sticky;
  top: 76px;
  z-index: 10;
  backdrop-filter: blur(6px);
}
.t-btn {
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  width: 32px; height: 32px;
  display: inline-flex; align-items: center; justify-content: center;
  color: rgb(var(--ww-ink));
  background: transparent;
  transition: background-color .25s ease, color .25s ease;
}
.t-btn:hover { background: rgb(var(--ww-gold) / .12); color: rgb(var(--ww-gold-deep)); }
.t-btn.on { background: rgb(var(--ww-ink)); color: rgb(var(--ww-parchment)); }
.t-btn.italic { font-style: italic; }
.t-select {
  height: 32px;
  background: transparent;
  border: 1px solid transparent;
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .18em;
  color: rgb(var(--ww-ink));
  padding: 0 6px;
  cursor: pointer;
  transition: background-color .25s ease, border-color .25s ease;
}
.t-select:hover { background: rgb(var(--ww-gold) / .12); border-color: var(--ww-ink-hairline); }
.t-select:focus { outline: none; border-color: rgb(var(--ww-gold)); }
.t-btn.wikilink {
  width: auto;
  gap: 6px;
  padding: 0 10px;
  font-family: 'Cormorant SC', serif;
  letter-spacing: .15em;
  font-size: 10px;
  text-transform: uppercase;
}
.sep {
  display: inline-block;
  width: 1px;
  height: 18px;
  background: var(--ww-ink-hairline);
  margin: 0 4px;
}
</style>
