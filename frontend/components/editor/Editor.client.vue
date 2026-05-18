<script setup lang="ts">
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { Wikilink } from './WikilinkExtension'
import { Spoiler } from './SpoilerExtension'
import { WWImage } from './ImageExtension'

const props = defineProps<{ modelValue: unknown }>()
const emit = defineEmits<{ 'update:modelValue': [value: unknown] }>()

const editor = shallowRef<Editor | null>(null)
const pickerOpen = ref(false)
const imagePickerOpen = ref(false)

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
      }),
      Wikilink,
      Spoiler,
      WWImage,
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
</script>

<template>
  <div class="editor-shell">
    <EditorToolbar
      :editor="editor"
      @wikilink="pickerOpen = true"
      @image="imagePickerOpen = true"
    />
    <EditorContent :editor="editor" class="editor-host" />
    <WikilinkPicker :open="pickerOpen" @close="pickerOpen = false" @pick="onPick" />
    <ImagePicker :open="imagePickerOpen" @close="imagePickerOpen = false" @pick="onImagePick" />
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
  ul, ol { padding-left: 1.4em; }
  li { margin: .25em 0; }

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
