<script setup lang="ts">
import type { Asset, EntitySummary } from '~/types/api'

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
useHead({ title: () => t('maps.newTitle') })

const router = useRouter()
const { $api } = useNuxtApp()
const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)

const asset = ref<Asset | null>(null)
const name = ref('')
const parent = ref<EntitySummary | null>(null)
const pending = ref(false)
const error = ref<string | null>(null)

function onUploaded(a: Asset) {
  asset.value = a
}

async function submit() {
  if (!asset.value || !name.value) return
  pending.value = true
  error.value = null
  try {
    const res = await $api<{ map: { id: number } }>('/maps', {
      method: 'POST',
      body: {
        name: name.value,
        asset_id: asset.value.id,
        width: asset.value.width || 0,
        height: asset.value.height || 0,
        parent_entity_id: parent.value?.id ?? null,
      },
    })
    await router.push(`/maps/${res.map.id}/edit`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Could not create map'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <NuxtLink to="/maps" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('maps.backList') }}
      </NuxtLink>

      <div class="stagger mt-8 mb-10">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('maps.eyebrowNew') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('maps.newTitle') }}</em>
        </h1>
        <p class="lede">{{ t('maps.newLede') }}</p>
      </div>

      <div class="stagger form">
        <div class="block">
          <h3 class="ww-label">{{ t('maps.form.image') }}</h3>
          <AssetUpload v-if="!asset" @uploaded="onUploaded" />
          <div v-else class="preview">
            <img :src="`${apiBase}/assets/${asset.id}`" :alt="asset.filename" />
            <button type="button" class="replace" @click="asset = null">{{ t('maps.form.replace') }}</button>
          </div>
        </div>

        <div class="block">
          <h3 class="ww-label">{{ t('maps.form.details') }}</h3>
          <label class="field">
            <span class="lbl">{{ t('maps.form.name') }}</span>
            <input v-model="name" class="ww-input" :placeholder="t('maps.form.namePlaceholder')" />
          </label>
          <label class="field">
            <span class="lbl">{{ t('maps.form.parent') }}</span>
            <EntityPicker v-model="parent" :placeholder="t('maps.form.parentHint')" />
          </label>
        </div>

        <Transition name="fade">
          <p v-if="error" class="error">{{ error }}</p>
        </Transition>

        <div class="actions">
          <button type="button" class="ww-btn-primary" :disabled="!asset || !name || pending" @click="submit">
            {{ pending ? t('common.loading') : t('maps.form.bind') }}
            <span class="arrow" aria-hidden="true">
              <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
                <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </span>
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.back { margin-bottom: 14px; }
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 6vw, 80px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 18px; }

.form { display: grid; gap: 40px; max-width: 900px; }
.block { display: grid; gap: 16px; }
.block h3 { margin: 0 0 4px; padding-bottom: 12px; border-bottom: 1px solid var(--ww-ink-hairline); }
.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}

.preview { position: relative; }
.preview img { max-width: 100%; border: 1px solid var(--ww-ink-hairline); }
.replace {
  position: absolute; top: 10px; right: 10px;
  font-family: 'Cormorant SC', serif;
  font-size: 11px; letter-spacing: .22em; text-transform: uppercase;
  padding: 6px 12px;
  background: rgb(var(--ww-parchment) / .92);
  border: 1px solid var(--ww-ink-hairline);
  color: rgb(var(--ww-ink));
  transition: background-color .3s ease, border-color .3s ease;
}
.replace:hover { background: rgb(var(--ww-gold) / .15); border-color: rgb(var(--ww-gold)); }

.actions { display: flex; justify-content: flex-end; }
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }

.error { color: rgb(var(--ww-vermilion)); font-style: italic; margin: 0; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
