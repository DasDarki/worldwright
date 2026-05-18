<script setup lang="ts">
import type { WorldEvent } from '~/types/api'

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data } = await useAsyncData(`event-${id.value}`, () =>
  $api<{ event: WorldEvent }>(`/events/${id.value}`),
)
const event = computed(() => data.value?.event || null)

useHead({ title: () => event.value ? `${t('events.editingPrefix')} ${event.value.title}` : t('events.editTitle') })

const pending = ref(false)
const error = ref<string | null>(null)
const showDelete = ref(false)

async function onSubmit(payload: any) {
  if (!event.value) return
  pending.value = true
  error.value = null
  try {
    await $api(`/events/${event.value.id}`, { method: 'PATCH', body: payload })
    await router.push('/events')
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to save event'
  } finally {
    pending.value = false
  }
}

async function onDelete() {
  if (!event.value) return
  pending.value = true
  try {
    await $api(`/events/${event.value.id}`, { method: 'DELETE' })
    await router.push('/events')
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to delete event'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section v-if="event" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <NuxtLink to="/events" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('events.backList') }}
      </NuxtLink>

      <div class="stagger mt-8 mb-10">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('events.eyebrowEdit') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('events.editingPrefix') }}</em> {{ event.title }}
        </h1>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>

      <EventForm
        :initial="event"
        :submitting="pending"
        :submit-label="t('events.form.save')"
        @submit="onSubmit"
      />

      <div class="ww-rule my-12" />

      <div class="danger">
        <Transition name="fade">
          <button v-if="!showDelete" type="button" class="danger-trigger" @click="showDelete = true">
            {{ t('events.delete') }}
          </button>
          <div v-else class="danger-confirm">
            <span class="confirm-text">{{ t('events.deleteConfirm') }}</span>
            <button type="button" class="cancel" @click="showDelete = false">{{ t('events.cancel') }}</button>
            <button type="button" class="ww-btn-primary destroy" @click="onDelete" :disabled="pending">
              {{ t('events.deleteForever') }}
            </button>
          </div>
        </Transition>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.back { margin-bottom: 14px; }
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
.error { color: rgb(var(--ww-vermilion)); font-style: italic; margin-bottom: 16px; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.danger { display: flex; align-items: center; }
.danger-trigger {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion) / .7);
  background: transparent;
  border: 1px dashed rgb(var(--ww-vermilion) / .4);
  padding: 8px 16px;
  transition: color .25s, background-color .25s, border-color .25s;
}
.danger-trigger:hover {
  color: rgb(var(--ww-vermilion));
  background: rgb(var(--ww-vermilion) / .08);
  border-color: rgb(var(--ww-vermilion));
}
.danger-confirm {
  display: flex;
  align-items: center;
  gap: 14px;
}
.confirm-text { font-style: italic; color: rgb(var(--ww-vermilion)); }
.cancel {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  background: transparent;
  border: 0;
}
.destroy { background: rgb(var(--ww-vermilion)); }
.destroy:hover { background: rgb(var(--ww-vermilion-deep)); }
</style>
