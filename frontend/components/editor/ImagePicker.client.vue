<script setup lang="ts">
import type { Asset } from '~/types/api'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; pick: [{ assetId: number; src: string; alt: string }] }>()

const { $api } = useNuxtApp()
const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)

const recent = ref<Asset[]>([])
const pending = ref(false)
const altDraft = ref('')

async function loadRecent() {
  pending.value = true
  try {
    const { assets } = await $api<{ assets: Asset[] }>('/assets?limit=24')
    recent.value = assets
  } catch {
    recent.value = []
  } finally {
    pending.value = false
  }
}

watch(() => props.open, async (isOpen) => {
  if (isOpen) {
    altDraft.value = ''
    await loadRecent()
  }
})

function onUploaded(a: Asset) {
  recent.value = [a, ...recent.value]
  insertAsset(a)
}

function insertAsset(a: Asset) {
  emit('pick', {
    assetId: a.id,
    src: `${apiBase.value}/assets/${a.id}`,
    alt: altDraft.value || a.filename || '',
  })
}
</script>

<template>
  <Transition name="picker">
    <div v-if="open" class="backdrop" @click.self="emit('close')">
      <div class="modal" role="dialog" aria-label="Insert image">
        <div class="ww-eyebrow head">Pin an image</div>

        <AssetUpload @uploaded="onUploaded" />

        <label class="alt">
          <span class="ww-label">Alt text</span>
          <input v-model="altDraft" class="ww-input" placeholder="A description for readers and crawlers" />
        </label>

        <div class="recent-head">Recent images</div>
        <div v-if="recent.length" class="grid">
          <button
            v-for="a in recent"
            :key="a.id"
            type="button"
            class="thumb"
            @click="insertAsset(a)"
            :aria-label="a.filename"
          >
            <img :src="`${apiBase}/assets/${a.id}`" :alt="a.filename" />
          </button>
        </div>
        <p v-else-if="!pending" class="empty">No images uploaded yet.</p>
      </div>
    </div>
  </Transition>
</template>

<style scoped lang="scss">
.backdrop {
  position: fixed; inset: 0;
  background: rgb(0 0 0 / .25);
  backdrop-filter: blur(2px);
  z-index: 200;
  display: flex; align-items: flex-start; justify-content: center;
  padding: 8vh 20px 20px;
}
.modal {
  width: min(680px, 100%);
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  padding: 24px;
  display: grid; gap: 16px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 60px 100px -40px rgb(0 0 0 / .45);
}
.head { margin-bottom: 0; }
.alt { display: grid; gap: 4px; }
.recent-head {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  margin-top: 6px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 8px;
}
.thumb {
  position: relative;
  border: 1px solid var(--ww-ink-hairline);
  background: transparent;
  padding: 0;
  cursor: pointer;
  aspect-ratio: 1;
  overflow: hidden;
  transition: border-color .25s ease, transform .25s cubic-bezier(.22,1,.36,1);
}
.thumb img { width: 100%; height: 100%; object-fit: cover; transition: transform .5s cubic-bezier(.22,1,.36,1); }
.thumb:hover { border-color: rgb(var(--ww-gold)); transform: translateY(-2px); }
.thumb:hover img { transform: scale(1.06); }
.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 20px 0;
}

.picker-enter-active, .picker-leave-active { transition: opacity .25s ease; }
.picker-enter-active .modal, .picker-leave-active .modal {
  transition: transform .35s cubic-bezier(.22,1,.36,1), opacity .25s ease;
}
.picker-enter-from, .picker-leave-to { opacity: 0; }
.picker-enter-from .modal, .picker-leave-to .modal {
  transform: translateY(-8px) scale(.98);
  opacity: 0;
}
</style>
