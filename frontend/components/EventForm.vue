<script setup lang="ts">
import type { Calendar, EntitySummary, EventParticipant, InWorldDate, WorldEvent } from '~/types/api'

const props = defineProps<{
  initial?: WorldEvent | null
  submitting?: boolean
  submitLabel?: string
}>()

interface EventFormPayload {
  title: string
  body: unknown
  calendar_id: number
  era_id: number | null
  year: number
  month_index: number
  day: number
  end_year: number | null
  end_month_index: number | null
  end_day: number | null
  importance: number
  visibility: 'secret' | 'player' | 'public'
  participants: { entity_id: number; role: string }[]
}

const emit = defineEmits<{ submit: [v: EventFormPayload] }>()

const { t } = useI18n()
const { $api } = useNuxtApp()

const { data: calendarsData } = await useAsyncData('event-form-calendars', () =>
  $api<{ calendars: Calendar[] }>('/calendars'),
)
const calendars = computed(() => calendarsData.value?.calendars || [])

const calendarId = ref<number>(props.initial?.calendar_id ?? calendars.value[0]?.id ?? 0)

const { data: calendarData, refresh: refreshCalendar } = await useAsyncData(
  () => `event-form-calendar-${calendarId.value}`,
  async () => {
    if (!calendarId.value) return null
    const res = await $api<{ calendar: Calendar }>(`/calendars/${calendarId.value}`)
    return res.calendar
  },
  { watch: [calendarId] },
)
const calendar = computed<Calendar | null>(() => calendarData.value ?? null)

const title = ref(props.initial?.title ?? '')
const body = ref<unknown>(props.initial?.body ?? { type: 'doc', content: [{ type: 'paragraph' }] })
const importance = ref(props.initial?.importance ?? 0)
const visibility = ref<'secret' | 'player' | 'public'>(props.initial?.visibility ?? 'secret')

const startDate = ref<InWorldDate | null>(props.initial
  ? { calendar_id: props.initial.calendar_id, era_id: props.initial.era_id ?? null, year: props.initial.year, month_index: props.initial.month_index, day: props.initial.day }
  : null)
const hasEnd = ref(!!(props.initial?.end_year != null))
const endDate = ref<InWorldDate | null>(props.initial && props.initial.end_year != null
  ? { calendar_id: props.initial.calendar_id, era_id: props.initial.era_id ?? null, year: props.initial.end_year, month_index: props.initial.end_month_index ?? 1, day: props.initial.end_day ?? 1 }
  : null)

watch(calendar, (c) => {
  if (!c) return
  if (!startDate.value) {
    startDate.value = { calendar_id: c.id, era_id: null, year: c.current_year, month_index: c.current_month_index, day: c.current_day }
  } else {
    startDate.value = { ...startDate.value, calendar_id: c.id }
  }
})

interface DraftParticipant { entity: EntitySummary | null; role: string }
const participants = ref<DraftParticipant[]>(
  (props.initial?.participants ?? []).map((p) => ({
    entity: p.entity ?? null,
    role: p.role || 'participant',
  })),
)

function addParticipant() {
  participants.value.push({ entity: null, role: 'participant' })
}
function removeParticipant(i: number) {
  participants.value.splice(i, 1)
}

const roleOptions = ['participant', 'location', 'cause', 'consequence']

function onSubmit() {
  if (!title.value || !calendarId.value || !startDate.value) return
  emit('submit', {
    title: title.value,
    body: body.value,
    calendar_id: calendarId.value,
    era_id: startDate.value.era_id ?? null,
    year: startDate.value.year,
    month_index: startDate.value.month_index,
    day: startDate.value.day,
    end_year: hasEnd.value && endDate.value ? endDate.value.year : null,
    end_month_index: hasEnd.value && endDate.value ? endDate.value.month_index : null,
    end_day: hasEnd.value && endDate.value ? endDate.value.day : null,
    importance: importance.value,
    visibility: visibility.value,
    participants: participants.value
      .filter((p) => p.entity)
      .map((p) => ({ entity_id: p.entity!.id, role: p.role || 'participant' })),
  })
}
</script>

<template>
  <form class="event-form stagger" @submit.prevent="onSubmit">
    <section class="block">
      <h3 class="ww-label">{{ t('events.form.essentials') }}</h3>
      <label class="field">
        <span class="lbl">{{ t('events.form.title') }}</span>
        <input v-model="title" class="ww-input" required maxlength="200" />
      </label>
      <div class="grid two">
        <label class="field">
          <span class="lbl">{{ t('events.form.calendar') }}</span>
          <select v-model="calendarId" class="ww-input">
            <option v-for="c in calendars" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </label>
        <label class="field">
          <span class="lbl">{{ t('events.form.visibility') }}</span>
          <select v-model="visibility" class="ww-input">
            <option value="secret">secret</option>
            <option value="player">player</option>
            <option value="public">public</option>
          </select>
        </label>
      </div>
    </section>

    <section class="block">
      <h3 class="ww-label">{{ t('events.form.when') }}</h3>
      <div class="field">
        <span class="lbl">{{ t('events.form.startDate') }}</span>
        <DateInput v-model="startDate" :calendar="calendar" show-moons />
      </div>
      <label class="check">
        <input v-model="hasEnd" type="checkbox" />
        <span>{{ t('events.form.span') }}</span>
      </label>
      <Transition name="slide">
        <div v-if="hasEnd" class="field">
          <span class="lbl">{{ t('events.form.endDate') }}</span>
          <DateInput v-model="endDate" :calendar="calendar" show-moons />
        </div>
      </Transition>
      <div class="grid two">
        <label class="field">
          <span class="lbl">{{ t('events.form.importance') }}</span>
          <input v-model.number="importance" type="number" :min="0" :max="5" class="ww-input" />
        </label>
      </div>
    </section>

    <section class="block">
      <div class="block-head">
        <h3 class="ww-label">{{ t('events.form.participants') }}</h3>
        <button type="button" class="add" @click="addParticipant">+ {{ t('events.form.addParticipant') }}</button>
      </div>
      <TransitionGroup name="row" tag="ul" class="participants">
        <li v-for="(p, i) in participants" :key="i" class="participant">
          <div class="role-wrap">
            <span class="ww-label tiny">{{ t('events.form.role') }}</span>
            <select v-model="p.role" class="ww-input small">
              <option v-for="r in roleOptions" :key="r" :value="r">{{ r }}</option>
            </select>
          </div>
          <div class="picker-wrap">
            <span class="ww-label tiny">{{ t('events.form.entity') }}</span>
            <EntityPicker v-model="p.entity" />
          </div>
          <button type="button" class="iconbtn del" @click="removeParticipant(i)" aria-label="Remove">×</button>
        </li>
        <li v-if="!participants.length" key="empty" class="empty">{{ t('events.form.noParticipants') }}</li>
      </TransitionGroup>
    </section>

    <section class="block body-block">
      <h3 class="ww-label">{{ t('events.form.body') }}</h3>
      <ClientOnly>
        <LazyEditor v-model="body" />
        <template #fallback>
          <div class="editor-fallback">Loading the inkwell…</div>
        </template>
      </ClientOnly>
    </section>

    <div class="actions">
      <button type="submit" class="ww-btn-primary" :disabled="submitting || !title || !startDate">
        {{ submitLabel || t('events.form.save') }}
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
.event-form { display: grid; gap: 50px; max-width: 900px; }
.block { display: grid; gap: 18px; }
.block h3 {
  margin: 0 0 4px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.block-head { display: flex; align-items: baseline; justify-content: space-between; }

.add {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  padding: 4px 12px;
  border: 1px solid rgb(var(--ww-vermilion) / .3);
  transition: background-color .25s ease;
}
.add:hover { background: rgb(var(--ww-vermilion) / .12); }

.grid { display: grid; gap: 18px 28px; }
.grid.two { grid-template-columns: 1fr 1fr; }
@media (max-width: 640px) { .grid.two { grid-template-columns: 1fr; } }

.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.tiny { font-size: 9px; }

.check {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-family: 'EB Garamond', serif;
  font-style: italic;
  color: var(--ww-ink-faint);
}
.check input { width: 14px; height: 14px; accent-color: rgb(var(--ww-gold)); }

.participants {
  list-style: none; margin: 0; padding: 0;
  display: grid; gap: 12px;
}
.participant {
  display: grid;
  grid-template-columns: 1fr 2fr auto;
  gap: 14px;
  align-items: end;
  padding: 12px 0;
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.role-wrap, .picker-wrap { display: grid; gap: 4px; }
.ww-input.small { padding: 8px 0; }

.iconbtn {
  width: 30px; height: 30px;
  display: inline-flex; align-items: center; justify-content: center;
  border: 1px solid var(--ww-ink-hairline);
  background: transparent;
  color: var(--ww-ink-faint);
  transition: border-color .25s, color .25s, background-color .25s;
  align-self: center;
}
.iconbtn.del:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .1); }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  padding: 14px 0;
  text-align: center;
  font-size: 14px;
}

.body-block { grid-column: 1 / -1; }
.editor-fallback {
  border: 1px solid var(--ww-ink-hairline);
  padding: 60px 30px;
  text-align: center;
  font-style: italic;
  color: var(--ww-ink-faint);
}

.actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid var(--ww-ink-hairline);
}
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }

select.ww-input {
  appearance: none;
  background-color: transparent;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23173842'/></svg>");
  background-repeat: no-repeat;
  background-position: right .3em center;
  background-size: 8px 5px;
  padding-right: 1.5em;
}
:root.dark select.ww-input {
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23ecdfc2'/></svg>");
}

.row-enter-active, .row-leave-active { transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1); }
.row-enter-from { opacity: 0; transform: translateY(-4px); }
.row-leave-to { opacity: 0; transform: translateX(8px); }
.row-leave-active { position: absolute; }

.slide-enter-active, .slide-leave-active {
  transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1), max-height .4s cubic-bezier(.22,1,.36,1);
  overflow: hidden;
  max-height: 200px;
}
.slide-enter-from, .slide-leave-to {
  opacity: 0; transform: translateY(-4px); max-height: 0;
}
</style>
