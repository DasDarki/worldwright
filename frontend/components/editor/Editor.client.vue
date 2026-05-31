<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Table from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableCell from '@tiptap/extension-table-cell'
import TableHeader from '@tiptap/extension-table-header'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import { common, createLowlight } from 'lowlight'
import { Wikilink } from './WikilinkExtension'
import { Spoiler } from './SpoilerExtension'
import { WWImage } from './ImageExtension'
import { Callout } from './CalloutExtension'
import { RelationshipGraphNode } from './RelationshipGraphExtension'

const lowlight = createLowlight(common)

const props = defineProps<{ modelValue: unknown }>()
const emit = defineEmits<{ 'update:modelValue': [value: unknown] }>()

const editor = shallowRef<Editor | null>(null)
const pickerOpen = ref(false)
const imagePickerOpen = ref(false)
const graphPickerOpen = ref(false)

function parseInitial(): any {
  try {
    if (!props.modelValue) return { type: 'doc', content: [{ type: 'paragraph' }] }
    if (typeof props.modelValue === 'string') return JSON.parse(props.modelValue)
    return props.modelValue
  } catch {
    return { type: 'doc', content: [{ type: 'paragraph' }] }
  }
}

onMounted(() => {
  editor.value = new Editor({
    content: parseInitial(),
    extensions: [
      StarterKit.configure({
        heading: { levels: [2, 3] },
        codeBlock: false, // replaced by CodeBlockLowlight below
      }),
      CodeBlockLowlight.configure({ lowlight }),
      Table.configure({ resizable: true }),
      TableRow,
      TableCell,
      TableHeader,
      Callout,
      Wikilink,
      Spoiler,
      WWImage,
      RelationshipGraphNode,
    ],
    onUpdate({ editor }) {
      emit('update:modelValue', editor.getJSON())
    },
    editorProps: {
      attributes: { class: 'tt-editor-content' },
    },
  })
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})

function onPick(payload: { slug: string; label: string }) {
  editor.value?.chain().focus().insertWikilink(payload).run()
  pickerOpen.value = false
}

function onImagePick(payload: { assetId: number; src: string; alt: string }) {
  editor.value?.chain().focus().insertWWImage({ src: payload.src, alt: payload.alt, assetId: payload.assetId }).run()
  imagePickerOpen.value = false
}

function onGraphPick(payload: { entityIds: number[] }) {
  editor.value?.chain().focus().insertRelationshipGraph({ entityIds: payload.entityIds }).run()
  graphPickerOpen.value = false
}
</script>

<template>
  <div class="editor-shell">
    <EditorToolbar
      :editor="editor"
      @wikilink="pickerOpen = true"
      @image="imagePickerOpen = true"
      @graph="graphPickerOpen = true"
    />
    <EditorContent :editor="editor" class="editor-host" />
    <WikilinkPicker :open="pickerOpen" @close="pickerOpen = false" @pick="onPick" />
    <ImagePicker :open="imagePickerOpen" @close="imagePickerOpen = false" @pick="onImagePick" />
    <GraphPicker :open="graphPickerOpen" @close="graphPickerOpen = false" @pick="onGraphPick" />
  </div>
</template>

<style lang="scss">
.editor-shell {
  display: flex;
  flex-direction: column;
}
.editor-host {
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment) / .4);
  padding: 28px 30px 32px;
  min-height: 320px;
}
.tt-editor-content {
  outline: none;
  font-family: 'EB Garamond', serif;
  font-size: 18px;
  line-height: 1.65;
  color: rgb(var(--ww-ink-shade));

  p { margin: 0 0 1em; }
  h2, h3 {
    font-family: 'Fraunces', serif;
    font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
    margin: 1.4em 0 .5em;
    letter-spacing: -0.02em;
  }
  blockquote {
    border-left: 3px solid rgb(var(--ww-gold));
    margin: 1.2em 0;
    padding: .2em 0 .2em 1em;
    font-style: italic;
  }
  ul, ol { padding-left: 1.6em; margin: 0 0 1em; }
  ul { list-style: disc outside; }
  ul ul { list-style: circle outside; }
  ul ul ul { list-style: square outside; }
  ol { list-style: decimal outside; }
  ol ol { list-style: lower-alpha outside; }
  ol ol ol { list-style: lower-roman outside; }
  li { margin: .3em 0; }
  li::marker { color: rgb(var(--ww-gold-deep)); font-family: 'Fraunces', serif; }

  hr {
    border: 0; height: 1px;
    background: linear-gradient(to right, transparent, var(--ww-ink-hairline) 8%, var(--ww-ink-hairline) 92%, transparent);
    margin: 1.6em 0;
  }

  pre {
    background: rgb(var(--ww-ink-shade) / .96);
    color: rgb(var(--ww-parchment));
    padding: 16px 20px;
    margin: 1.4em 0;
    overflow-x: auto;
    font-family: 'JetBrains Mono', ui-monospace, monospace;
    font-size: 14px;
    line-height: 1.55;
  }
  pre code { font-family: inherit; background: transparent; color: inherit; padding: 0; }
  code {
    font-family: 'JetBrains Mono', ui-monospace, monospace;
    background: rgb(var(--ww-ink-hairline));
    padding: 1px 6px;
    font-size: .9em;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin: 1.4em 0;
    font-size: 15px;
    table-layout: fixed;
  }
  table th, table td {
    border: 1px solid var(--ww-ink-hairline);
    padding: 8px 10px;
    text-align: left;
    vertical-align: top;
    position: relative;
    min-width: 80px;
  }
  table th {
    background: rgb(var(--ww-parchment-stain) / .5);
    font-family: 'Cormorant SC', serif;
    font-size: 11px;
    letter-spacing: .2em;
    text-transform: uppercase;
    color: rgb(var(--ww-gold-deep));
  }
  table .selectedCell:after {
    content: ''; position: absolute; inset: 0;
    background: rgb(var(--ww-gold) / .18); pointer-events: none;
  }
  .tableWrapper { overflow-x: auto; }

  .ww-callout {
    margin: 1.4em 0;
    padding: 14px 18px 14px 50px;
    border: 1px solid var(--ww-ink-hairline);
    border-left-width: 4px;
    position: relative;
    background: rgb(var(--ww-parchment-stain) / .25);
  }
  .ww-callout p:last-child { margin-bottom: 0; }
  .ww-callout::before {
    content: '';
    position: absolute;
    left: 14px; top: 14px;
    width: 22px; height: 22px;
    background-position: center;
    background-repeat: no-repeat;
    background-size: contain;
    opacity: .85;
  }
  .ww-callout[data-variant="info"]    { border-left-color: rgb(var(--ww-ink-shade)); }
  .ww-callout[data-variant="warn"]    { border-left-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .06); }
  .ww-callout[data-variant="note"]    { border-left-color: rgb(var(--ww-gold-deep)); background: rgb(var(--ww-gold) / .08); }
  .ww-callout[data-variant="lore"]    { border-left-color: rgb(var(--ww-gold)); background: linear-gradient(135deg, rgb(var(--ww-gold) / .12), rgb(var(--ww-vermilion) / .06)); font-style: italic; }
  .ww-callout[data-variant="info"]::before { background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23173842' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><circle cx='12' cy='12' r='9'/><path d='M12 8 V13 M12 16.5 V17'/></svg>"); }
  .ww-callout[data-variant="warn"]::before { background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23a8442e' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M12 3 L22 20 L2 20 Z'/><path d='M12 10 V14 M12 17 V17.5'/></svg>"); }
  .ww-callout[data-variant="note"]::before { background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%237c5e30' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M5 4 H17 L20 7 V20 H5 Z'/><path d='M17 4 V7 H20 M8 11 H16 M8 15 H14'/></svg>"); }
  .ww-callout[data-variant="lore"]::before { background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23b8935a' stroke-width='1.6' stroke-linecap='round' stroke-linejoin='round'><path d='M12 3 L13.5 10.5 L21 12 L13.5 13.5 L12 21 L10.5 13.5 L3 12 L10.5 10.5 Z'/></svg>"); }

  .ww-graph-placeholder {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 1.4em 0;
    padding: 18px 22px;
    border: 1px dashed rgb(var(--ww-gold-deep) / .5);
    background: rgb(var(--ww-gold) / .08);
    color: rgb(var(--ww-gold-deep));
    user-select: none;
  }
  .ww-graph-placeholder__head {
    font-family: 'Cormorant SC', serif;
    font-size: 12px;
    letter-spacing: .28em;
    text-transform: uppercase;
  }
  .ww-graph-placeholder__count {
    font-family: 'EB Garamond', serif;
    font-style: italic;
    color: var(--ww-ink-faint);
  }
  .ww-graph-placeholder.ProseMirror-selectednode {
    outline: 2px solid rgb(var(--ww-vermilion));
    outline-offset: 2px;
  }

  .spoiler-block {
    position: relative;
    padding: 16px 20px 16px 22px;
    margin: 1em 0;
    background: repeating-linear-gradient(
      45deg,
      rgb(var(--ww-ink) / .04),
      rgb(var(--ww-ink) / .04) 8px,
      rgb(var(--ww-ink) / .08) 8px,
      rgb(var(--ww-ink) / .08) 16px
    );
    border-left: 3px solid rgb(var(--ww-vermilion));
    border-radius: 0;
  }
  .spoiler-block::before {
    content: 'SPOILER';
    position: absolute;
    top: 4px; left: 16px;
    font-family: 'Cormorant SC', serif;
    font-size: 8px;
    letter-spacing: .32em;
    color: rgb(var(--ww-vermilion));
  }

  .ww-img-embed {
    display: block;
    max-width: 100%;
    height: auto;
    margin: 1.4em auto;
    border: 1px solid var(--ww-ink-hairline);
    box-shadow: 0 30px 60px -30px rgb(0 0 0 / .35);
    transition: transform .35s cubic-bezier(.22,1,.36,1);
  }
  .ProseMirror-selectednode.ww-img-embed { outline: 2px solid rgb(var(--ww-gold)); outline-offset: 2px; }

  .wikilink-chip {
    display: inline-block;
    padding: 0 6px;
    color: rgb(var(--ww-vermilion));
    background: rgb(var(--ww-vermilion) / .1);
    border-radius: 2px;
    border: 1px solid rgb(var(--ww-vermilion) / .25);
    font-family: 'Cormorant SC', serif;
    font-size: 0.86em;
    letter-spacing: .08em;
    margin: 0 1px;
    cursor: default;
    user-select: all;
    transition: background-color .3s ease, transform .3s ease;
  }
  .wikilink-chip:hover {
    background: rgb(var(--ww-vermilion) / .2);
    transform: translateY(-1px);
  }
  .ProseMirror-focused .wikilink-chip.ProseMirror-selectednode {
    outline: 2px solid rgb(var(--ww-gold));
    outline-offset: 1px;
  }

  p.is-editor-empty:first-child::before {
    content: 'Begin with a coastline, a name, & a god who lied…';
    float: left;
    color: var(--ww-ink-faint);
    pointer-events: none;
    height: 0;
    font-style: italic;
  }
}
</style>
