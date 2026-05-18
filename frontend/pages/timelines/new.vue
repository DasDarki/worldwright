<script setup lang="ts">
definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
useHead({ title: () => t('timelines.newTitle') })

const router = useRouter()
const { $api } = useNuxtApp()

const pending = ref(false)
const error = ref<string | null>(null)

async function onSubmit(payload: any) {
  pending.value = true
  error.value = null
  try {
    const res = await $api<{ timeline: { id: number } }>('/timelines', {
      method: 'POST',
      body: payload,
    })
    await router.push(`/timelines/${res.timeline.id}`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to create timeline'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <NuxtLink to="/timelines" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('timelines.backList') }}
      </NuxtLink>

      <div class="stagger mt-8 mb-10">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('timelines.eyebrowNew') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('timelines.newTitle') }}</em>
        </h1>
        <p class="lede">{{ t('timelines.newLede') }}</p>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>

      <TimelineForm :submitting="pending" :submit-label="t('timelines.form.create')" @submit="onSubmit" />
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
.error { color: rgb(var(--ww-vermilion)); font-style: italic; margin-bottom: 16px; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
