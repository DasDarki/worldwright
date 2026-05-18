<script setup lang="ts">
const props = defineProps<{ url: string }>()

const { t } = useI18n()
const copied = ref(false)

async function copy() {
  try {
    if (navigator?.clipboard) {
      await navigator.clipboard.writeText(absoluteUrl.value)
    } else {
      const el = document.createElement('textarea')
      el.value = absoluteUrl.value
      document.body.appendChild(el)
      el.select()
      document.execCommand('copy')
      document.body.removeChild(el)
    }
    copied.value = true
    setTimeout(() => { copied.value = false }, 1800)
  } catch {
    copied.value = false
  }
}

const absoluteUrl = computed(() => {
  if (import.meta.server) return props.url
  return new URL(props.url, window.location.origin).toString()
})
</script>

<template>
  <button type="button" class="share-btn ww-btn-ghost" @click="copy">
    <Transition name="swap" mode="out-in">
      <span v-if="!copied" key="copy" class="row">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" aria-hidden="true">
          <circle cx="4"  cy="8" r="2"/>
          <circle cx="12" cy="4" r="2"/>
          <circle cx="12" cy="12" r="2"/>
          <path d="M5.7 7 L 10.3 4.6"/>
          <path d="M5.7 9 L 10.3 11.4"/>
        </svg>
        {{ t('share.copy') }}
      </span>
      <span v-else key="copied" class="row done">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true">
          <path d="M3 8 L7 12 L13 4"/>
        </svg>
        {{ t('share.copied') }}
      </span>
    </Transition>
  </button>
</template>

<style scoped lang="scss">
.share-btn { gap: 8px; }
.row { display: inline-flex; align-items: center; gap: 8px; }
.row.done { color: rgb(var(--ww-gold-deep)); }

.swap-enter-active, .swap-leave-active { transition: opacity .25s ease, transform .25s cubic-bezier(.22,1,.36,1); }
.swap-enter-from { opacity: 0; transform: translateY(-3px); }
.swap-leave-to { opacity: 0; transform: translateY(3px); }
</style>
