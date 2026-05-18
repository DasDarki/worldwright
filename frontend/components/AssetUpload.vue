<script setup lang="ts">
import type { Asset } from '~/types/api'

const props = withDefaults(defineProps<{
  label?: string
  accept?: string
}>(), {
  label: 'Drop a chart or click to choose',
  accept: 'image/png,image/jpeg,image/webp,image/gif',
})

const emit = defineEmits<{ uploaded: [a: Asset] }>()

const { $api } = useNuxtApp()
const inputRef = ref<HTMLInputElement | null>(null)
const dragging = ref(false)
const uploading = ref(false)
const error = ref<string | null>(null)
const progress = ref(0)

async function handleFile(file: File) {
  error.value = null
  uploading.value = true
  progress.value = 0
  try {
    const form = new FormData()
    form.append('file', file)
    const res = await $api<{ asset: Asset }>('/assets', { method: 'POST', body: form })
    progress.value = 1
    emit('uploaded', res.asset)
  } catch (e: any) {
    error.value = e?.data?.error || 'Upload failed'
  } finally {
    setTimeout(() => { uploading.value = false; progress.value = 0 }, 300)
  }
}

function onPick(e: Event) {
  const f = (e.target as HTMLInputElement).files?.[0]
  if (f) handleFile(f)
}

function onDrop(e: DragEvent) {
  e.preventDefault()
  dragging.value = false
  const f = e.dataTransfer?.files?.[0]
  if (f) handleFile(f)
}
</script>

<template>
  <label
    :class="['uploader', { dragging, uploading }]"
    @dragenter.prevent="dragging = true"
    @dragover.prevent="dragging = true"
    @dragleave.prevent="dragging = false"
    @drop="onDrop"
  >
    <input
      ref="inputRef"
      type="file"
      :accept="accept"
      class="sr-only"
      @change="onPick"
    />
    <div v-if="!uploading" class="prompt">
      <svg viewBox="0 0 36 36" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.2" aria-hidden="true">
        <rect x="4" y="6" width="28" height="24" stroke-dasharray="2 3" />
        <circle cx="13" cy="14" r="2" fill="currentColor" stroke="none" />
        <path d="M4 24 L13 16 L20 22 L26 18 L32 24" stroke-linejoin="round" />
      </svg>
      <span class="text">{{ label }}</span>
      <span class="hint">PNG · JPG · WebP · GIF · up to 20 MB</span>
    </div>
    <div v-else class="progress" aria-live="polite">
      <div class="bar" :style="{ width: `${Math.max(progress * 100, 10)}%` }" />
      <span class="status">Pressing the leaf…</span>
    </div>
    <Transition name="fade">
      <p v-if="error" class="error">{{ error }}</p>
    </Transition>
  </label>
</template>

<style scoped lang="scss">
.uploader {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px 30px;
  border: 1px dashed var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment) / .5);
  color: rgb(var(--ww-ink));
  cursor: pointer;
  transition: background-color .35s ease, border-color .35s ease, transform .35s cubic-bezier(.22,1,.36,1);
}
.uploader:hover { background: rgb(var(--ww-gold) / .08); border-color: rgb(var(--ww-gold)); }
.uploader.dragging {
  background: rgb(var(--ww-gold) / .15);
  border-color: rgb(var(--ww-gold));
  transform: scale(1.01);
  border-style: solid;
}
.uploader.uploading { cursor: progress; }

.sr-only {
  position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px;
  overflow: hidden; clip: rect(0,0,0,0); white-space: nowrap; border: 0;
}

.prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: rgb(var(--ww-gold-deep));
}
.prompt .text {
  font-family: 'Cormorant SC', serif;
  letter-spacing: .22em;
  text-transform: uppercase;
  font-size: 12px;
  color: rgb(var(--ww-ink));
}
.prompt .hint {
  font-style: italic;
  font-size: 12px;
  color: var(--ww-ink-faint);
}

.progress {
  width: 100%;
  display: grid;
  gap: 8px;
}
.bar {
  height: 4px;
  background: rgb(var(--ww-gold));
  transition: width .35s cubic-bezier(.22,1,.36,1);
}
.status {
  font-style: italic;
  text-align: center;
  color: var(--ww-ink-faint);
  font-size: 13px;
}

.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  margin: 0;
  font-size: 14px;
}
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
