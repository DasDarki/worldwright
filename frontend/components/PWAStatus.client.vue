<script setup lang="ts">
import { useToastsStore } from '~/stores/toasts'

const { t } = useI18n()
const toasts = useToastsStore()

const nuxtApp = useNuxtApp()
const pwa = (nuxtApp.$pwa as any) ?? null

const showInstall = ref(false)
const needRefresh = ref(false)
const offlineReady = ref(false)

watch(
  () => pwa?.offlineReady,
  (v) => {
    offlineReady.value = !!v
    if (v) toasts.success(t('pwa.offlineReady'), 4000)
  },
)

watch(
  () => pwa?.needRefresh,
  (v) => {
    needRefresh.value = !!v
  },
)

watch(
  () => pwa?.showInstallPrompt,
  (v) => {
    showInstall.value = !!v
  },
  { immediate: true },
)

async function install() {
  if (!pwa?.install) return
  try {
    await pwa.install()
    showInstall.value = false
  } catch {
    showInstall.value = false
  }
}

function dismissInstall() {
  if (pwa?.cancelInstall) pwa.cancelInstall()
  showInstall.value = false
}

async function applyUpdate() {
  if (pwa?.updateServiceWorker) {
    await pwa.updateServiceWorker(true)
  }
}

function dismissUpdate() {
  if (pwa?.cancelPrompt) pwa.cancelPrompt()
  needRefresh.value = false
}
</script>

<template>
  <Teleport to="body">
    <div class="pwa-stack" aria-live="polite">
      <Transition name="banner">
        <div v-if="needRefresh" class="banner update">
          <div class="message">
            <strong>{{ t('pwa.updateTitle') }}</strong>
            <span>{{ t('pwa.updateBody') }}</span>
          </div>
          <div class="actions">
            <button type="button" class="ghost" @click="dismissUpdate">{{ t('pwa.later') }}</button>
            <button type="button" class="primary" @click="applyUpdate">{{ t('pwa.refresh') }}</button>
          </div>
        </div>
      </Transition>

      <Transition name="banner">
        <div v-if="showInstall && !needRefresh" class="banner install">
          <div class="message">
            <strong>{{ t('pwa.installTitle') }}</strong>
            <span>{{ t('pwa.installBody') }}</span>
          </div>
          <div class="actions">
            <button type="button" class="ghost" @click="dismissInstall">{{ t('pwa.notNow') }}</button>
            <button type="button" class="primary" @click="install">{{ t('pwa.install') }}</button>
          </div>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
.pwa-stack {
  position: fixed;
  left: 24px;
  bottom: 24px;
  z-index: 350;
  display: flex;
  flex-direction: column-reverse;
  gap: 12px;
  pointer-events: none;
  max-width: min(420px, calc(100% - 48px));
}
.banner {
  pointer-events: auto;
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 14px;
  align-items: center;
  padding: 14px 16px;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 24px 40px -22px rgb(0 0 0 / .45);
  position: relative;
  overflow: hidden;
}
.banner::before {
  content: '';
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 3px;
}
.banner.update::before { background: rgb(var(--ww-gold)); }
.banner.install::before { background: rgb(var(--ww-vermilion)); }

.message {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}
.message strong {
  font-family: 'Cormorant SC', serif;
  font-weight: 600;
  font-size: 11px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: rgb(var(--ww-ink));
}
.banner.update .message strong { color: rgb(var(--ww-gold-deep)); }
.banner.install .message strong { color: rgb(var(--ww-vermilion)); }

.message span {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  font-size: 14px;
  color: var(--ww-ink-faint);
  line-height: 1.35;
}

.actions {
  display: flex;
  align-items: center;
  gap: 6px;
}
.actions button {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .22em;
  text-transform: uppercase;
  padding: 7px 12px;
  border: 1px solid var(--ww-ink-hairline);
  background: transparent;
  color: rgb(var(--ww-ink));
  transition: background-color .25s ease, color .25s ease, border-color .25s ease;
}
.actions .ghost:hover {
  background: rgb(var(--ww-ink) / .06);
}
.actions .primary {
  background: rgb(var(--ww-ink));
  color: rgb(var(--ww-parchment));
  border-color: rgb(var(--ww-ink));
}
.actions .primary:hover {
  background: rgb(var(--ww-ink-shade));
}

.banner-enter-active, .banner-leave-active {
  transition: opacity .35s ease, transform .4s cubic-bezier(.22,1,.36,1);
}
.banner-enter-from { opacity: 0; transform: translateY(20px); }
.banner-leave-to { opacity: 0; transform: translateY(20px); }
</style>
