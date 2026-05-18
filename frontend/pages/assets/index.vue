<script setup lang="ts">
import type { Asset } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const { t } = useI18n()
useHead({ title: () => t('assets.title') })

const auth = useAuthStore()
const { $api } = useNuxtApp()
const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)

const { data, refresh } = await useAsyncData('asset-library', () =>
  $api<{ assets: Asset[] }>('/assets?limit=200'),
)
const assets = computed<Asset[]>(() => data.value?.assets || [])

const query = ref('')
const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return assets.value
  return assets.value.filter((a) => a.filename.toLowerCase().includes(q))
})

const showUpload = ref(false)
function onUploaded() {
  showUpload.value = false
  refresh()
}

const lightbox = ref<Asset | null>(null)

async function remove(a: Asset) {
  if (!confirm(t('assets.confirmDelete'))) return
  try {
    await $api(`/assets/${a.id}`, { method: 'DELETE' })
    lightbox.value = null
    refresh()
  } catch {
    // ignore
  }
}

function copyEmbed(a: Asset) {
  const url = `${apiBase.value}/assets/${a.id}`
  try {
    navigator?.clipboard?.writeText(url)
  } catch {
    // ignore
  }
}

function formatBytes(n: number): string {
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / (1024 * 1024)).toFixed(2)} MB`
}

function formatDate(s: string): string {
  try { return new Date(s).toLocaleDateString() } catch { return s }
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-2xl">
      <div class="stagger mb-12 flex items-end justify-between gap-6 flex-wrap">
        <div>
          <div class="ww-eyebrow mb-6 flex items-center gap-3">
            <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
            Vol. I · The plate-press
          </div>
          <h1 class="hero-title"><em>{{ t('assets.title') }}</em></h1>
          <p class="lede">{{ t('assets.lede') }}</p>
        </div>
        <button
          v-if="auth.isAdmin"
          type="button"
          class="ww-btn-primary upload-btn"
          @click="showUpload = !showUpload"
        >
          {{ showUpload ? t('assets.closeUpload') : t('assets.uploadNew') }}
        </button>
      </div>

      <Transition name="slide">
        <div v-if="showUpload" class="upload-bay">
          <AssetUpload @uploaded="onUploaded" />
        </div>
      </Transition>

      <div class="filter">
        <input v-model="query" type="search" class="ww-input" :placeholder="t('assets.filter')" />
        <span class="count">{{ filtered.length }} / {{ assets.length }}</span>
      </div>

      <p v-if="!assets.length" class="empty">{{ t('assets.empty') }}</p>

      <TransitionGroup v-else name="asset" tag="ul" class="grid">
        <li
          v-for="a in filtered"
          :key="a.id"
          class="card reveal"
          @click="lightbox = a"
        >
          <div class="thumb">
            <img :src="`${apiBase}/assets/${a.id}`" :alt="a.filename" loading="lazy" />
            <span class="overlay" aria-hidden="true" />
            <div class="hover">
              <span class="ww-label">{{ t('assets.open') }}</span>
            </div>
          </div>
          <div class="meta">
            <div class="filename" :title="a.filename">{{ a.filename }}</div>
            <div class="sub">
              <span v-if="a.width && a.height">{{ a.width }} × {{ a.height }}</span>
              <span class="dot" aria-hidden="true">·</span>
              <span>{{ formatBytes(a.size) }}</span>
            </div>
          </div>
        </li>
      </TransitionGroup>

      <Transition name="lightbox">
        <div v-if="lightbox" class="lightbox" @click.self="lightbox = null">
          <div class="frame">
            <button type="button" class="close" :aria-label="t('common.ok')" @click="lightbox = null">×</button>
            <img :src="`${apiBase}/assets/${lightbox.id}`" :alt="lightbox.filename" />
            <div class="info">
              <div class="line">
                <span class="ww-label">{{ t('assets.filename') }}</span>
                <span class="val">{{ lightbox.filename }}</span>
              </div>
              <div class="line">
                <span class="ww-label">{{ t('assets.dimensions') }}</span>
                <span class="val">{{ lightbox.width && lightbox.height ? `${lightbox.width} × ${lightbox.height}` : '—' }}</span>
              </div>
              <div class="line">
                <span class="ww-label">{{ t('assets.size') }}</span>
                <span class="val">{{ formatBytes(lightbox.size) }}</span>
              </div>
              <div class="line">
                <span class="ww-label">{{ t('assets.uploaded') }}</span>
                <span class="val">{{ formatDate(lightbox.created_at) }}</span>
              </div>
              <div class="line">
                <span class="ww-label">{{ t('assets.url') }}</span>
                <code class="url">{{ `${apiBase}/assets/${lightbox.id}` }}</code>
              </div>

              <div class="actions">
                <button type="button" class="ww-btn-ghost" @click="copyEmbed(lightbox)">
                  {{ t('assets.copyUrl') }}
                </button>
                <button
                  v-if="auth.isAdmin"
                  type="button"
                  class="ww-btn-ghost destroy"
                  @click="remove(lightbox)"
                >{{ t('assets.delete') }}</button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </section>
</template>

<style scoped lang="scss">
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(48px, 8vw, 120px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 18px; margin-top: 6px; }

.upload-btn { white-space: nowrap; }

.upload-bay {
  margin-bottom: 28px;
  padding: 6px;
  background: rgb(var(--ww-parchment-deep) / .35);
  border: 1px dashed var(--ww-ink-hairline);
}

.filter {
  display: flex;
  align-items: baseline;
  gap: 18px;
  margin: 20px 0 28px;
  max-width: 36em;
}
.count {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  color: var(--ww-ink-faint);
}

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  padding: 80px 0;
}

.grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 22px;
}

.card {
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: transform .35s cubic-bezier(.22,1,.36,1);
}
.card:hover { transform: translateY(-3px); }

.thumb {
  position: relative;
  aspect-ratio: 1;
  border: 1px solid var(--ww-ink-hairline);
  overflow: hidden;
  background: rgb(var(--ww-parchment-deep) / .4);
}
.thumb img {
  width: 100%; height: 100%;
  object-fit: cover;
  filter: sepia(.15);
  transition: transform .6s cubic-bezier(.22,1,.36,1), filter .5s ease;
}
.card:hover .thumb img { transform: scale(1.06); filter: sepia(.02); }
.overlay {
  position: absolute;
  inset: 0;
  pointer-events: none;
  box-shadow: inset 0 0 40px rgb(14 36 44 / .25);
}
.hover {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 14px;
  opacity: 0;
  background: linear-gradient(to top, rgb(0 0 0 / .55), transparent 60%);
  transition: opacity .35s ease;
  color: rgb(var(--ww-parchment));
}
.card:hover .hover { opacity: 1; }

.meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.filename {
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: rgb(var(--ww-ink));
}
.sub {
  font-family: 'Cormorant SC', serif;
  font-size: 9px;
  letter-spacing: .22em;
  color: var(--ww-ink-faint);
  display: flex;
  gap: 6px;
  align-items: center;
}
.sub .dot { opacity: .6; }

/* Lightbox */
.lightbox {
  position: fixed;
  inset: 0;
  background: rgb(0 0 0 / .55);
  backdrop-filter: blur(3px);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}
.frame {
  position: relative;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  max-width: 92vw;
  max-height: 90vh;
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 0;
  overflow: hidden;
  box-shadow: 0 60px 100px -40px rgb(0 0 0 / .55);
}
@media (max-width: 880px) {
  .frame { grid-template-columns: 1fr; max-height: 96vh; }
  .frame img { max-height: 60vh; }
}
.frame img {
  display: block;
  width: 100%;
  height: 100%;
  max-height: 90vh;
  object-fit: contain;
  background: rgb(var(--ww-ink-shade) / .85);
}
.close {
  position: absolute;
  top: 12px; right: 12px;
  width: 34px; height: 34px;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  font-size: 22px;
  line-height: 1;
  color: rgb(var(--ww-ink));
  cursor: pointer;
  z-index: 2;
  transition: background-color .25s ease, color .25s ease, border-color .25s ease;
}
.close:hover { background: rgb(var(--ww-vermilion)); color: rgb(var(--ww-parchment)); border-color: rgb(var(--ww-vermilion)); }
.info {
  padding: 24px 26px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  border-left: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment));
  overflow-y: auto;
}
.line { display: flex; flex-direction: column; gap: 4px; }
.val { font-family: 'EB Garamond', serif; font-size: 15px; color: rgb(var(--ww-ink)); }
.url {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px;
  background: rgb(var(--ww-ink-shade) / .08);
  padding: 4px 6px;
  word-break: break-all;
}
.actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  padding-top: 16px;
  border-top: 1px dashed var(--ww-ink-hairline);
}
.destroy { color: rgb(var(--ww-vermilion) / .8); border-color: rgb(var(--ww-vermilion) / .4); }
.destroy:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .08); }

.asset-enter-active, .asset-leave-active {
  transition: opacity .35s ease, transform .4s cubic-bezier(.22,1,.36,1);
}
.asset-enter-from { opacity: 0; transform: translateY(8px) scale(.96); }
.asset-leave-to { opacity: 0; transform: scale(.92); }

.slide-enter-active, .slide-leave-active {
  transition: max-height .45s cubic-bezier(.22,1,.36,1), opacity .35s ease;
  overflow: hidden;
  max-height: 320px;
}
.slide-enter-from, .slide-leave-to { opacity: 0; max-height: 0; }

.lightbox-enter-active, .lightbox-leave-active {
  transition: opacity .3s ease;
}
.lightbox-enter-active .frame, .lightbox-leave-active .frame {
  transition: transform .4s cubic-bezier(.22,1,.36,1), opacity .35s ease;
}
.lightbox-enter-from, .lightbox-leave-to { opacity: 0; }
.lightbox-enter-from .frame, .lightbox-leave-to .frame {
  opacity: 0;
  transform: scale(.95);
}
</style>
