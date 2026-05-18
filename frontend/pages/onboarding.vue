<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const { t } = useI18n()
const auth = useAuthStore()
const toasts = useToastsStore()
const router = useRouter()

useHead({ title: () => t('onboarding.title') })

if (!auth.onboarding) await auth.fetchOnboarding()

const pending = ref<'keep' | 'prune' | 'import' | null>(null)
const importSummary = ref<{
  entities_created: number
  maps_created: number
  pins_created: number
  images_imported: number
  images_failed: number
  warnings?: string[]
} | null>(null)

const fileInput = ref<HTMLInputElement | null>(null)
const importFile = ref<File | null>(null)

function pickFile() {
  fileInput.value?.click()
}
function onFileChosen(e: Event) {
  const t = e.target as HTMLInputElement
  importFile.value = t.files && t.files[0] ? t.files[0] : null
}

async function chooseKeepSeed() {
  pending.value = 'keep'
  try {
    await useNuxtApp().$api('/admin/onboarding/keep-seed', { method: 'POST' })
    auth.setOnboardingCompleted(false)
    await router.replace('/')
  } catch (e: any) {
    toasts.error(e?.data?.error || t('onboarding.errors.generic'))
  } finally {
    pending.value = null
  }
}

async function choosePruneSeed() {
  if (!confirm(t('onboarding.confirmPrune'))) return
  pending.value = 'prune'
  try {
    await useNuxtApp().$api('/admin/onboarding/prune-seed', { method: 'POST' })
    auth.setOnboardingCompleted(true)
    await router.replace('/')
  } catch (e: any) {
    toasts.error(e?.data?.error || t('onboarding.errors.generic'))
  } finally {
    pending.value = null
  }
}

async function chooseImportLK() {
  if (!importFile.value) {
    pickFile()
    return
  }
  if (!confirm(t('onboarding.confirmImport'))) return
  pending.value = 'import'
  try {
    const form = new FormData()
    form.append('file', importFile.value)
    const { summary } = await useNuxtApp().$api<{ summary: typeof importSummary.value }>(
      '/admin/onboarding/import/legendkeeper',
      { method: 'POST', body: form },
    )
    importSummary.value = summary
    auth.setOnboardingCompleted(true)
  } catch (e: any) {
    toasts.error(e?.data?.error || t('onboarding.errors.importFailed'))
  } finally {
    pending.value = null
  }
}

function finishImport() {
  router.replace('/')
}
</script>

<template>
  <section class="onboarding">
    <div class="mx-auto max-w-screen-xl px-6 md:px-12 py-16 md:py-24">
      <div class="ww-eyebrow flex items-center gap-3 mb-6">
        <span class="eyebrow-rule" aria-hidden="true" />
        {{ t('onboarding.eyebrow') }}
      </div>
      <h1 class="hero-title mb-6">{{ t('onboarding.title') }}</h1>
      <p class="lede mb-12">{{ t('onboarding.subtitle') }}</p>

      <div v-if="importSummary" class="summary-card">
        <h2 class="summary-title mb-4">{{ t('onboarding.summary.title') }}</h2>
        <ul class="summary-list mb-4">
          <li><strong>{{ importSummary.entities_created }}</strong> {{ t('onboarding.summary.entities') }}</li>
          <li><strong>{{ importSummary.maps_created }}</strong> {{ t('onboarding.summary.maps') }}</li>
          <li><strong>{{ importSummary.pins_created }}</strong> {{ t('onboarding.summary.pins') }}</li>
          <li>
            <strong>{{ importSummary.images_imported }}</strong> {{ t('onboarding.summary.images') }}
            <span v-if="importSummary.images_failed > 0" class="failed">
              ({{ importSummary.images_failed }} {{ t('onboarding.summary.imagesFailed') }})
            </span>
          </li>
        </ul>
        <details v-if="importSummary.warnings && importSummary.warnings.length" class="warnings">
          <summary>{{ t('onboarding.summary.warnings') }} ({{ importSummary.warnings.length }})</summary>
          <ul>
            <li v-for="(w, idx) in importSummary.warnings" :key="idx">{{ w }}</li>
          </ul>
        </details>
        <button class="ww-btn-primary mt-6" type="button" @click="finishImport">
          {{ t('onboarding.summary.enter') }}
          <span class="arrow" aria-hidden="true">→</span>
        </button>
      </div>

      <div v-else class="grid md:grid-cols-3 gap-6 md:gap-8">
        <article class="choice">
          <header>
            <span class="numeral">i</span>
            <h2>{{ t('onboarding.options.keep.title') }}</h2>
          </header>
          <p class="lede">{{ t('onboarding.options.keep.body') }}</p>
          <button class="ww-btn-primary" type="button" :disabled="!!pending" @click="chooseKeepSeed">
            {{ pending === 'keep' ? t('common.loading') : t('onboarding.options.keep.cta') }}
          </button>
        </article>

        <article class="choice">
          <header>
            <span class="numeral">ii</span>
            <h2>{{ t('onboarding.options.prune.title') }}</h2>
          </header>
          <p class="lede">{{ t('onboarding.options.prune.body') }}</p>
          <button class="ww-btn-primary" type="button" :disabled="!!pending" @click="choosePruneSeed">
            {{ pending === 'prune' ? t('common.loading') : t('onboarding.options.prune.cta') }}
          </button>
        </article>

        <article class="choice">
          <header>
            <span class="numeral">iii</span>
            <h2>{{ t('onboarding.options.import.title') }}</h2>
          </header>
          <p class="lede">{{ t('onboarding.options.import.body') }}</p>
          <input
            ref="fileInput"
            type="file"
            accept="application/json,.json"
            class="hidden-input"
            @change="onFileChosen"
          />
          <p v-if="importFile" class="file-name">{{ importFile.name }}</p>
          <div class="import-actions">
            <button class="ww-btn-ghost" type="button" :disabled="!!pending" @click="pickFile">
              {{ importFile ? t('onboarding.options.import.changeFile') : t('onboarding.options.import.pickFile') }}
            </button>
            <button class="ww-btn-primary" type="button" :disabled="!!pending || !importFile" @click="chooseImportLK">
              {{ pending === 'import' ? t('onboarding.options.import.running') : t('onboarding.options.import.cta') }}
            </button>
          </div>
          <p class="note">{{ t('onboarding.options.import.note') }}</p>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.eyebrow-rule {
  display: inline-block; width: 36px; height: 1px; background: currentColor;
}
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 5.2vw, 68px);
  line-height: .98;
  letter-spacing: -.02em;
  margin: 0;
}
.lede {
  font-size: 17px;
  line-height: 1.55;
  color: var(--ww-ink-faint);
  max-width: 40em;
}

.choice {
  border: 1px solid var(--ww-ink-hairline);
  background: var(--ww-card-bg);
  padding: 28px 24px 26px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
  isolation: isolate;
  box-shadow: 0 30px 60px -40px rgb(0 0 0 / .25);
}
.choice header {
  display: flex; align-items: baseline; gap: 14px;
}
.choice .numeral {
  font-family: 'Cormorant SC', serif;
  font-size: 18px;
  letter-spacing: .3em;
  color: rgb(var(--ww-gold-deep));
  text-transform: lowercase;
}
.choice h2 {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 460;
  font-size: 24px;
  line-height: 1.15;
  margin: 0;
}
.choice .lede {
  font-size: 15px;
  margin: 0;
  flex: 1;
}
.choice .ww-btn-primary {
  align-self: flex-start;
}

.import-actions {
  display: flex; gap: 12px; flex-wrap: wrap; align-items: center;
}
.hidden-input {
  position: absolute; width: 1px; height: 1px;
  opacity: 0; pointer-events: none;
}
.file-name {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  margin: 0;
}
.note {
  font-style: italic;
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin: 8px 0 0;
}

.summary-card {
  border: 1px solid var(--ww-ink-hairline);
  background: var(--ww-card-bg);
  padding: 32px 36px;
  max-width: 36em;
  box-shadow: 0 30px 60px -40px rgb(0 0 0 / .25);
}
.summary-title {
  font-family: 'Fraunces', serif;
  font-size: 28px;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 480;
  margin: 0;
}
.summary-list { list-style: none; padding: 0; margin: 0; display: grid; gap: 8px; }
.summary-list li { font-size: 16px; }
.summary-list strong {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 50, "opsz" 32, "wght" 600;
  font-size: 22px;
  color: rgb(var(--ww-gold-deep));
  margin-right: 6px;
}
.summary-list .failed {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  font-size: 13px;
}
.warnings { margin-top: 20px; }
.warnings summary {
  cursor: pointer;
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.warnings ul {
  margin-top: 12px;
  padding: 12px 16px;
  font-size: 13px;
  font-family: 'IBM Plex Mono', monospace;
  border-left: 2px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment-stain) / .4);
  list-style: none;
  display: grid;
  gap: 4px;
  max-height: 240px;
  overflow: auto;
}
</style>
