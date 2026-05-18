<script setup lang="ts">
import type { Timeline } from '~/types/api'

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data } = await useAsyncData(`tl-edit-${id.value}`, () =>
  $api<{ timeline: Timeline }>(`/timelines/${id.value}`),
)
const timeline = computed(() => data.value?.timeline || null)

useHead({ title: () => timeline.value ? `${t('timelines.editingPrefix')} ${timeline.value.name}` : t('timelines.editTitle') })

const pending = ref(false)
const error = ref<string | null>(null)

async function onSubmit(payload: any) {
  if (!timeline.value) return
  pending.value = true
  error.value = null
  try {
    await $api(`/timelines/${timeline.value.id}`, { method: 'PATCH', body: payload })
    await router.push(`/timelines/${timeline.value.id}`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to save timeline'
  } finally {
    pending.value = false
  }
}

async function onDelete() {
  if (!timeline.value) return
  if (!confirm(t('timelines.confirmDelete'))) return
  pending.value = true
  try {
    await $api(`/timelines/${timeline.value.id}`, { method: 'DELETE' })
    await router.push('/timelines')
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to delete'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section v-if="timeline" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="topbar">
        <NuxtLink :to="`/timelines/${timeline.id}`" class="ww-btn-ghost back">
          <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
            <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('timelines.backView') }}
        </NuxtLink>
        <button type="button" class="ww-btn-ghost destroy" @click="onDelete" :disabled="pending">{{ t('timelines.delete') }}</button>
      </div>

      <div class="stagger mt-8 mb-10">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('timelines.eyebrowEdit') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('timelines.editingPrefix') }}</em> {{ timeline.name }}
        </h1>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>

      <TimelineForm
        :initial="timeline"
        :submitting="pending"
        :submit-label="t('timelines.form.save')"
        @submit="onSubmit"
      />
    </div>
  </section>
</template>

<style scoped lang="scss">
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(36px, 5vw, 64px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.destroy { color: rgb(var(--ww-vermilion) / .8); border-color: rgb(var(--ww-vermilion) / .4); }
.destroy:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .08); }
.error { color: rgb(var(--ww-vermilion)); font-style: italic; margin-bottom: 16px; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
