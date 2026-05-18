<script setup lang="ts">
import type { Calendar, Timeline, WorldEvent } from '~/types/api'

const props = defineProps<{
  initial?: Timeline | null
  submitting?: boolean
  submitLabel?: string
}>()

interface TimelineFormPayload {
  name: string
  description: string
  event_ids: number[]
}

const emit = defineEmits<{ submit: [v: TimelineFormPayload] }>()

const { t } = useI18n()
const { $api } = useNuxtApp()

const { data: eventsData } = await useAsyncData('tl-form-events', () =>
  $api<{ events: WorldEvent[] }>('/events'),
)
const { data: calsData } = await useAsyncData('tl-form-cals', () =>
  $api<{ calendars: Calendar[] }>('/calendars'),
)

const allEvents = computed(() => eventsData.value?.events || [])
const allCalendars = computed(() => calsData.value?.calendars || [])

const name = ref(props.initial?.name ?? '')
const description = ref(props.initial?.description ?? '')
const eventIds = ref<number[]>([...(props.initial?.event_ids ?? [])])

function onSubmit() {
  if (!name.value) return
  emit('submit', {
    name: name.value,
    description: description.value,
    event_ids: eventIds.value,
  })
}
</script>

<template>
  <form class="tl-form stagger" @submit.prevent="onSubmit">
    <section class="block">
      <h3 class="ww-label section-head">{{ t('timelines.form.essentials') }}</h3>
      <label class="field">
        <span class="lbl">{{ t('timelines.form.name') }}</span>
        <input v-model="name" class="ww-input" required />
      </label>
      <label class="field">
        <span class="lbl">{{ t('timelines.form.description') }}</span>
        <textarea v-model="description" class="ww-input textarea" rows="3" />
      </label>
    </section>

    <section class="block">
      <h3 class="ww-label section-head">{{ t('timelines.form.curate') }}</h3>
      <TimelineCurator v-model="eventIds" :events="allEvents" :calendars="allCalendars" />
    </section>

    <div class="actions">
      <button type="submit" class="ww-btn-primary" :disabled="submitting || !name">
        {{ submitLabel || t('timelines.form.save') }}
        <span class="arrow" aria-hidden="true">
          <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
            <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </span>
      </button>
    </div>
  </form>
</template>

<style scoped lang="scss">
.tl-form { display: grid; gap: 40px; max-width: 900px; }
.block { display: grid; gap: 16px; }
.section-head { margin: 0; padding-bottom: 12px; border-bottom: 1px solid var(--ww-ink-hairline); }
.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.textarea {
  border: 1px solid var(--ww-ink-hairline);
  padding: 10px 12px;
  resize: vertical;
}
.actions {
  display: flex; justify-content: flex-end;
  padding-top: 16px; border-top: 1px solid var(--ww-ink-hairline);
}
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }
</style>
