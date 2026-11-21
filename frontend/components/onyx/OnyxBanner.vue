<script setup lang="ts">
import OnyxCrystal from '~/components/onyx/OnyxCrystal.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const onyx = useOnyx()
const banner = ref<HTMLElement | null>(null)
const busy = ref(false)
const failed = ref(false)
let hideTimer: ReturnType<typeof setTimeout> | undefined
let observer: ResizeObserver | null = null

const onLoginPage = computed(() => route.path === '/login')
const showMigration = computed(() => onyx.needsMigration.value && !onLoginPage.value)
const success = computed(() => onyx.result.value === 'linked')
const message = computed(() => (onyx.result.value ? t(`onyx.result.${onyx.result.value}`) : null))
const visible = computed(() => !onLoginPage.value && (showMigration.value || message.value !== null))
const crystalState = computed(() => (busy.value ? 'loading' : success.value ? 'success' : 'idle'))

watch(
  () => route.query.onyx,
  (value) => {
    if (!isOnyxResult(value) || onLoginPage.value) return
    onyx.result.value = value
    const query = { ...route.query }
    delete query.onyx
    router.replace({ query })
  },
  { immediate: true },
)

watch(
  () => onyx.result.value,
  (value) => {
    clearTimeout(hideTimer)
    if (value === 'linked') hideTimer = setTimeout(() => (onyx.result.value = null), 7000)
  },
  { immediate: true },
)

watch(banner, (element) => {
  observer?.disconnect()
  if (!element) {
    document.documentElement.style.removeProperty('--onyx-banner-height')
    return
  }
  observer = new ResizeObserver(() => {
    document.documentElement.style.setProperty('--onyx-banner-height', `${element.offsetHeight}px`)
  })
  observer.observe(element)
})

onBeforeUnmount(() => {
  observer?.disconnect()
  clearTimeout(hideTimer)
  document.documentElement.style.removeProperty('--onyx-banner-height')
})

async function migrate() {
  busy.value = true
  failed.value = false
  try {
    await onyx.startMigration()
  } catch {
    busy.value = false
    failed.value = true
  }
}

function dismiss() {
  onyx.result.value = null
}
</script>

<template>
  <Transition name="onyx-banner">
    <section
      v-if="visible"
      ref="banner"
      class="onyx-banner"
      :class="{ 'is-success': success, 'is-error': message && !success }"
      role="region"
      aria-label="Onyx"
    >
      <span class="edge-light" aria-hidden="true" />
      <div class="crystal-wrap">
        <OnyxCrystal size="34px" interactive :state="crystalState" />
      </div>

      <div v-if="message" class="copy">
        <p class="title">{{ success ? t('onyx.linkedTitle') : 'Onyx' }}</p>
        <p class="text">{{ message }}</p>
      </div>
      <div v-else class="copy">
        <p class="title">{{ t('onyx.bannerTitle') }}</p>
        <p class="text">
          {{ t('onyx.bannerText') }}
          <span v-if="failed" class="inline-error">{{ t('onyx.migrateFailed') }}</span>
        </p>
      </div>

      <button v-if="message" type="button" class="dismiss" :aria-label="t('onyx.dismiss')" @click="dismiss">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
          <path d="M3.5 3.5 L12.5 12.5 M12.5 3.5 L3.5 12.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
        </svg>
      </button>
      <button v-else type="button" class="migrate" :disabled="busy" @click="migrate">
        <span class="sweep" aria-hidden="true" />
        <span class="label">{{ busy ? t('onyx.redirecting') : t('onyx.migrate') }}</span>
      </button>
    </section>
  </Transition>
</template>

<style scoped lang="scss">
.onyx-banner {
  --stone-0: #030407;
  --stone-1: #0a0d13;
  --seam: #1e2533;
  --ash: #8b95a9;
  --quartz: #e9edf5;
  --edge: #5b9dff;
  --glint: #bcd6ff;
  --ok: #48d597;
  --danger: #ff6178;
  --ease: cubic-bezier(0.2, 0.8, 0.2, 1);

  position: relative;
  z-index: 60;
  font-family: ui-sans-serif, system-ui, -apple-system, 'Segoe UI', sans-serif;
  display: flex;
  align-items: center;
  gap: 1.1rem;
  padding: 0.85rem clamp(1rem, 3vw, 2.25rem);
  background:
    radial-gradient(120% 220% at 0% 50%, rgba(91, 157, 255, 0.13), transparent 55%),
    linear-gradient(180deg, var(--stone-1), var(--stone-0));
  border-bottom: 1px solid var(--seam);
  color: var(--quartz);
  font-size: 0.95rem;
  overflow: hidden;
}

.edge-light {
  position: absolute;
  left: 0;
  right: 0;
  bottom: -1px;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--glint), var(--edge), transparent);
  background-size: 40% 100%;
  background-repeat: no-repeat;
  animation: edge-run 5.5s linear infinite;
}

.is-success .edge-light {
  background-image: linear-gradient(90deg, transparent, var(--ok), transparent);
}

.is-error .edge-light {
  background-image: linear-gradient(90deg, transparent, var(--danger), transparent);
}

.crystal-wrap {
  flex: none;
  display: grid;
  place-items: center;
  width: 44px;
  height: 52px;
}

.copy {
  flex: 1;
  min-width: 0;
}

.title {
  margin: 0;
  font-weight: 700;
  color: var(--quartz);
  font-size: 1.05rem;
  letter-spacing: 0.005em;
}

.text {
  margin: 0.15rem 0 0;
  color: var(--ash);
  line-height: 1.45;
  max-width: 70ch;
}

.inline-error {
  display: block;
  margin-top: 0.25rem;
  color: var(--danger);
}

.migrate {
  position: relative;
  flex: none;
  margin-left: auto;
  min-height: 46px;
  padding: 0 1.6rem;
  border: 1px solid rgba(91, 157, 255, 0.65);
  background: linear-gradient(180deg, #12203a, #0b1424);
  color: var(--quartz);
  font: inherit;
  font-weight: 700;
  font-size: 1rem;
  letter-spacing: 0.01em;
  cursor: pointer;
  overflow: hidden;
  clip-path: polygon(0 0, calc(100% - 12px) 0, 100% 12px, 100% 100%, 0 100%);
  transition:
    border-color 250ms var(--ease),
    transform 120ms var(--ease),
    box-shadow 250ms var(--ease);

  &:hover:not(:disabled) {
    border-color: var(--glint);
    box-shadow: inset 0 0 0 1px rgba(188, 214, 255, 0.25);
  }

  &:hover:not(:disabled) .sweep {
    transform: translateX(260%) skewX(-20deg);
    transition: transform 900ms var(--ease);
  }

  &:active:not(:disabled) {
    transform: translateY(1px);
  }

  &:disabled {
    cursor: progress;
    opacity: 0.8;
  }

  &:focus-visible {
    outline: 2px solid var(--glint);
    outline-offset: 2px;
  }
}

.sweep {
  position: absolute;
  inset: 0 auto 0 -40%;
  width: 35%;
  background: linear-gradient(90deg, transparent, rgba(188, 214, 255, 0.28), transparent);
  transform: translateX(0) skewX(-20deg);
  pointer-events: none;
}

.label {
  position: relative;
}

.dismiss {
  flex: none;
  margin-left: auto;
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  color: var(--ash);
  cursor: pointer;

  &:hover {
    color: var(--quartz);
    border-color: var(--seam);
  }

  :deep(svg) {
    width: 18px;
    height: 18px;
  }
}

.onyx-banner-enter-active,
.onyx-banner-leave-active {
  transition:
    opacity 420ms var(--ease),
    transform 420ms var(--ease);
}

.onyx-banner-enter-from,
.onyx-banner-leave-to {
  opacity: 0;
  transform: translateY(-100%);
}

@keyframes edge-run {
  from {
    background-position: -40% 0;
  }
  to {
    background-position: 140% 0;
  }
}

@media (max-width: 720px) {
  .onyx-banner {
    flex-wrap: wrap;
    gap: 0.75rem 0.9rem;
  }

  .copy {
    flex-basis: calc(100% - 60px);
  }

  .migrate {
    width: 100%;
    margin-left: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .edge-light {
    animation: none;
    background-position: 50% 0;
  }

  .sweep {
    display: none;
  }
}
</style>
