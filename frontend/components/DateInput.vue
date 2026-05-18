<script setup lang="ts">
import type { Calendar, InWorldDate } from '~/types/api'

const props = withDefaults(defineProps<{
  modelValue: InWorldDate | null
  calendar: Calendar | null
  required?: boolean
  showMoons?: boolean
}>(), {
  required: false,
  showMoons: false,
})

const emit = defineEmits<{ 'update:modelValue': [v: InWorldDate | null] }>()

const calendarId = computed(() => props.calendar?.id ?? 0)
const eraId = ref<number | null>(props.modelValue?.era_id ?? null)
const year = ref<number>(props.modelValue?.year ?? props.calendar?.current_year ?? 0)
const monthIndex = ref<number>(props.modelValue?.month_index ?? props.calendar?.current_month_index ?? 1)
const day = ref<number>(props.modelValue?.day ?? props.calendar?.current_day ?? 1)

const months = computed(() => props.calendar?.months ?? [])
const eras = computed(() => props.calendar?.eras ?? [])

const activeMonth = computed(() => months.value.find((m) => m.sort_order === monthIndex.value) ?? null)
const maxDay = computed(() => activeMonth.value?.days ?? 31)

watch([year, monthIndex, day, eraId, calendarId], () => {
  if (!calendarId.value) {
    emit('update:modelValue', null)
    return
  }
  if (day.value > maxDay.value) day.value = maxDay.value
  emit('update:modelValue', {
    calendar_id: calendarId.value,
    era_id: eraId.value ?? null,
    year: year.value,
    month_index: monthIndex.value,
    day: day.value,
  })
})

watch(() => props.modelValue, (v) => {
  if (!v) return
  eraId.value = v.era_id ?? null
  year.value = v.year
  monthIndex.value = v.month_index
  day.value = v.day
})
</script>

<template>
  <div class="date-input-wrap">
    <div class="date-input">
      <label class="seg seg-day">
        <span class="seg-label">d</span>
        <input v-model.number="day" type="number" :min="1" :max="maxDay" />
      </label>
      <label class="seg seg-month">
        <span class="seg-label">m</span>
        <select v-model.number="monthIndex">
          <option v-for="m in months" :key="m.sort_order" :value="m.sort_order">{{ m.name }}</option>
        </select>
      </label>
      <label class="seg seg-year">
        <span class="seg-label">y</span>
        <input v-model.number="year" type="number" />
      </label>
      <label v-if="eras.length" class="seg seg-era">
        <span class="seg-label">era</span>
        <select v-model.number="eraId">
          <option :value="null">auto</option>
          <option v-for="e in eras" :key="e.id" :value="e.id">{{ e.suffix || e.name }}</option>
        </select>
      </label>
    </div>
    <Transition name="moon-preview">
      <MoonDisplay
        v-if="showMoons && calendar && modelValue"
        class="preview"
        :calendar="calendar"
        :date="modelValue"
        :size="26"
        variant="inline"
      />
    </Transition>
  </div>
</template>

<style scoped lang="scss">
.date-input {
  display: inline-flex;
  align-items: stretch;
  gap: 0;
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment) / .5);
  transition: border-color .3s ease;
}
.date-input:focus-within { border-color: rgb(var(--ww-gold)); }
.seg {
  display: inline-flex;
  flex-direction: column;
  justify-content: center;
  padding: 4px 10px;
  border-right: 1px solid var(--ww-ink-hairline);
}
.seg:last-child { border-right: 0; }
.seg-label {
  font-family: 'Cormorant SC', serif;
  font-size: 8px;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.seg input, .seg select {
  border: 0;
  background: transparent;
  outline: none;
  font-family: 'EB Garamond', serif;
  color: rgb(var(--ww-ink));
  font-size: 16px;
  padding: 2px 0;
  min-width: 0;
}
.seg-day input { width: 36px; }
.seg-year input { width: 64px; }
.seg-month select { min-width: 110px; }
.seg-era select { min-width: 64px; }

select { appearance: none; }

.date-input-wrap {
  display: inline-flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-start;
}
.preview {
  padding-left: 2px;
}

.moon-preview-enter-active, .moon-preview-leave-active {
  transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1);
}
.moon-preview-enter-from, .moon-preview-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
