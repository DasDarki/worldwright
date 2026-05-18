<script setup lang="ts">
import type { Calendar, CalendarEra, CalendarMonth, CalendarMoon, CalendarWeekday } from '~/types/api'
import { absoluteDay } from '~/composables/useCalendar'
import { useAuthStore } from '~/stores/auth'
import { useToastsStore } from '~/stores/toasts'

const { t } = useI18n()
const auth = useAuthStore()
const toasts = useToastsStore()
const route = useRoute()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data, refresh } = await useAsyncData(`calendar-${id.value}`, () =>
  $api<{ calendar: Calendar }>(`/calendars/${id.value}`),
)

useHead({ title: () => data.value?.calendar.name || t('calendars.title') })

const cal = reactive<Calendar>(JSON.parse(JSON.stringify(data.value!.calendar)))

const pending = ref(false)

function move<T>(arr: T[], i: number, dir: -1 | 1) {
  const j = i + dir
  if (j < 0 || j >= arr.length) return
  ;[arr[i], arr[j]] = [arr[j]!, arr[i]!]
}

function addMonth() { cal.months.push({ sort_order: cal.months.length + 1, name: 'New month', days: 30 } as CalendarMonth) }
function removeMonth(i: number) { cal.months.splice(i, 1) }

function addWeekday() { cal.weekdays.push({ sort_order: cal.weekdays.length + 1, name: 'New day' } as CalendarWeekday) }
function removeWeekday(i: number) { cal.weekdays.splice(i, 1) }

function addEra() { cal.eras.push({ name: 'New era', start_year: 0, suffix: '' } as CalendarEra) }
function removeEra(i: number) { cal.eras.splice(i, 1) }

function addMoon() {
  cal.moons.push({
    name: 'New moon',
    cycle_days: 28,
    offset_days: 0,
    phases: [
      { sort_order: 1, name: 'New',  cycle_position: 0.0, randomness: 0.02, icon: '🌑' },
      { sort_order: 2, name: 'Full', cycle_position: 0.5, randomness: 0.02, icon: '🌕' },
    ],
  } as CalendarMoon)
}
function removeMoon(i: number) { cal.moons.splice(i, 1) }

function addPhase(moon: CalendarMoon) {
  if (!moon.phases) moon.phases = []
  moon.phases.push({
    sort_order: moon.phases.length + 1,
    name: 'New phase',
    cycle_position: 0.25,
    randomness: 0.05,
    icon: '',
  })
}
function removePhase(moon: CalendarMoon, i: number) {
  moon.phases.splice(i, 1)
}

function currentCycleFor(moon: CalendarMoon): number {
  const day = absoluteDay(cal, {
    calendar_id: cal.id,
    year: cal.current_year,
    month_index: cal.current_month_index,
    day: cal.current_day,
  })
  const cycle = moon.cycle_days || 1
  return ((day - moon.offset_days) % cycle + cycle) % cycle / cycle
}

async function save() {
  if (!auth.isAdmin) return
  pending.value = true
  try {
    cal.months.forEach((m, i) => (m.sort_order = i + 1))
    cal.weekdays.forEach((w, i) => (w.sort_order = i + 1))
    await $api(`/calendars/${cal.id}`, { method: 'PATCH', body: cal })
    toasts.success(t('calendars.saved'))
    await refresh()
  } catch (e: any) {
    toasts.error(e?.data?.error || 'Could not save calendar')
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <NuxtLink to="/calendars" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('calendars.backList') }}
      </NuxtLink>

      <div class="stagger mt-8 mb-16">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('calendars.eyebrow') }}
        </div>
        <h1 class="hero-title mb-2">
          <em>{{ cal.name }}</em>
        </h1>
        <p v-if="cal.epoch_name" class="lede">{{ t('calendars.epochOf') }} {{ cal.epoch_name }}</p>
      </div>

      <div class="cal-preview reveal">
        <div class="ww-eyebrow preview-eyebrow">{{ t('moons.tonight') }}</div>
        <MoonDisplay :calendar="cal" variant="inline" :size="56" :glow="true" />
      </div>

      <fieldset :disabled="!auth.isAdmin" class="form">
        <section class="block">
          <h3 class="ww-label section-head">{{ t('calendars.essentials') }}</h3>
          <div class="grid two">
            <label class="field">
              <span class="lbl">{{ t('calendars.name') }}</span>
              <input v-model="cal.name" class="ww-input" />
            </label>
            <label class="field">
              <span class="lbl">{{ t('calendars.epoch') }}</span>
              <input v-model="cal.epoch_name" class="ww-input" />
            </label>
            <label class="field">
              <span class="lbl">{{ t('calendars.currentYear') }}</span>
              <input v-model.number="cal.current_year" type="number" class="ww-input" />
            </label>
            <label class="field">
              <span class="lbl">{{ t('calendars.currentMonth') }}</span>
              <input v-model.number="cal.current_month_index" type="number" :min="1" :max="cal.months.length" class="ww-input" />
            </label>
            <label class="field">
              <span class="lbl">{{ t('calendars.currentDay') }}</span>
              <input v-model.number="cal.current_day" type="number" :min="1" class="ww-input" />
            </label>
            <label class="field">
              <span class="lbl">{{ t('calendars.leapRule') }}</span>
              <input v-model="cal.leap_rule" class="ww-input" :placeholder="t('calendars.leapHint')" />
            </label>
          </div>
        </section>

        <section class="block">
          <div class="block-head">
            <h3 class="ww-label section-head">{{ t('calendars.months') }}</h3>
            <button v-if="auth.isAdmin" type="button" class="add" @click="addMonth">+ {{ t('calendars.addMonth') }}</button>
          </div>
          <TransitionGroup name="row" tag="ul" class="rows">
            <li v-for="(m, i) in cal.months" :key="`m-${i}`" class="row">
              <span class="idx">{{ i + 1 }}</span>
              <input v-model="m.name" class="ww-input flex-1" :placeholder="t('calendars.monthName')" />
              <label class="days">
                <input v-model.number="m.days" type="number" :min="1" class="ww-input small" />
                <span class="ww-label suffix">d</span>
              </label>
              <div v-if="auth.isAdmin" class="row-actions">
                <button type="button" class="iconbtn" :disabled="i === 0" @click="move(cal.months, i, -1)" aria-label="Move up">↑</button>
                <button type="button" class="iconbtn" :disabled="i === cal.months.length - 1" @click="move(cal.months, i, 1)" aria-label="Move down">↓</button>
                <button type="button" class="iconbtn del" @click="removeMonth(i)" aria-label="Remove">×</button>
              </div>
            </li>
          </TransitionGroup>
        </section>

        <section class="block">
          <div class="block-head">
            <h3 class="ww-label section-head">{{ t('calendars.weekdays') }}</h3>
            <button v-if="auth.isAdmin" type="button" class="add" @click="addWeekday">+ {{ t('calendars.addWeekday') }}</button>
          </div>
          <TransitionGroup name="row" tag="ul" class="rows">
            <li v-for="(w, i) in cal.weekdays" :key="`w-${i}`" class="row">
              <span class="idx">{{ i + 1 }}</span>
              <input v-model="w.name" class="ww-input flex-1" :placeholder="t('calendars.weekdayName')" />
              <div v-if="auth.isAdmin" class="row-actions">
                <button type="button" class="iconbtn" :disabled="i === 0" @click="move(cal.weekdays, i, -1)">↑</button>
                <button type="button" class="iconbtn" :disabled="i === cal.weekdays.length - 1" @click="move(cal.weekdays, i, 1)">↓</button>
                <button type="button" class="iconbtn del" @click="removeWeekday(i)">×</button>
              </div>
            </li>
          </TransitionGroup>
        </section>

        <section class="block">
          <div class="block-head">
            <h3 class="ww-label section-head">{{ t('calendars.eras') }}</h3>
            <button v-if="auth.isAdmin" type="button" class="add" @click="addEra">+ {{ t('calendars.addEra') }}</button>
          </div>
          <TransitionGroup name="row" tag="ul" class="rows">
            <li v-for="(e, i) in cal.eras" :key="`e-${i}`" class="row era-row">
              <input v-model="e.name" class="ww-input flex-1" :placeholder="t('calendars.eraName')" />
              <label class="seg">
                <span class="ww-label tiny">{{ t('calendars.startYear') }}</span>
                <input v-model.number="e.start_year" type="number" class="ww-input small" />
              </label>
              <label class="seg">
                <span class="ww-label tiny">{{ t('calendars.suffix') }}</span>
                <input v-model="e.suffix" class="ww-input small" placeholder="LR" />
              </label>
              <button v-if="auth.isAdmin" type="button" class="iconbtn del" @click="removeEra(i)">×</button>
            </li>
          </TransitionGroup>
        </section>

        <section class="block">
          <div class="block-head">
            <h3 class="ww-label section-head">{{ t('calendars.moons') }}</h3>
            <button v-if="auth.isAdmin" type="button" class="add" @click="addMoon">+ {{ t('calendars.addMoon') }}</button>
          </div>
          <TransitionGroup name="row" tag="div" class="moons">
            <article v-for="(m, i) in cal.moons" :key="`mo-${i}`" class="moon-card">
              <header class="moon-head">
                <MoonGlyph :cycle="0.25" :size="40" />
                <input v-model="m.name" class="ww-input name-input" :placeholder="t('calendars.moonName')" />
                <label class="seg">
                  <span class="ww-label tiny">{{ t('calendars.cycleDays') }}</span>
                  <input v-model.number="m.cycle_days" type="number" step="0.1" class="ww-input small" />
                </label>
                <label class="seg">
                  <span class="ww-label tiny">{{ t('calendars.offsetDays') }}</span>
                  <input v-model.number="m.offset_days" type="number" step="0.1" class="ww-input small" />
                </label>
                <button v-if="auth.isAdmin" type="button" class="iconbtn del" @click="removeMoon(i)">×</button>
              </header>

              <MoonCycleTimeline
                v-if="(m.phases || []).length"
                :moon="m"
                :current-cycle="currentCycleFor(m)"
                class="cycle-preview"
              />

              <div class="phases-block">
                <div class="phases-head">
                  <span class="ww-label">{{ t('calendars.phases') }}</span>
                  <button
                    v-if="auth.isAdmin"
                    type="button"
                    class="add small"
                    @click="addPhase(m)"
                  >+ {{ t('calendars.addPhase') }}</button>
                </div>
                <TransitionGroup name="row" tag="ul" class="phases">
                  <li v-for="(p, j) in (m.phases || [])" :key="`p-${i}-${j}`" class="phase-row">
                    <MoonGlyph :cycle="p.cycle_position" :size="36" />
                    <input v-model="p.icon" class="ww-input icon-input" placeholder="🌑" maxlength="2" />
                    <input v-model="p.name" class="ww-input flex-1" :placeholder="t('calendars.phaseName')" />
                    <div class="slider-seg">
                      <span class="ww-label tiny">{{ t('calendars.position') }}</span>
                      <div class="slider-row">
                        <input v-model.number="p.cycle_position" type="range" min="0" max="1" step="0.005" class="range" />
                        <input v-model.number="p.cycle_position" type="number" min="0" max="1" step="0.005" class="ww-input mini" />
                      </div>
                    </div>
                    <div class="slider-seg">
                      <span class="ww-label tiny">±</span>
                      <div class="slider-row">
                        <input v-model.number="p.randomness" type="range" min="0" max="0.5" step="0.005" class="range" />
                        <input v-model.number="p.randomness" type="number" min="0" max="0.5" step="0.005" class="ww-input mini" />
                      </div>
                    </div>
                    <button v-if="auth.isAdmin" type="button" class="iconbtn del" @click="removePhase(m, j)">×</button>
                  </li>
                  <li v-if="!(m.phases || []).length" :key="`empty-${i}`" class="phase-empty">{{ t('calendars.phasesEmpty') }}</li>
                </TransitionGroup>
              </div>
            </article>
          </TransitionGroup>
        </section>
      </fieldset>

      <div v-if="auth.isAdmin" class="actions">
        <button type="button" class="ww-btn-primary" :disabled="pending" @click="save">
          {{ pending ? t('common.loading') : t('calendars.save') }}
          <span class="arrow" aria-hidden="true">
            <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
              <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
        </button>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.back { margin-bottom: 14px; }
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(48px, 8vw, 96px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 18px; }

.cal-preview {
  display: grid;
  gap: 14px;
  padding: 28px 30px;
  margin-bottom: 40px;
  background: linear-gradient(180deg,
    rgb(var(--ww-parchment-deep) / .45),
    rgb(var(--ww-parchment-deep) / .2));
  border: 1px solid var(--ww-ink-hairline);
}
.preview-eyebrow { margin: 0; }

.form { display: grid; gap: 50px; border: 0; padding: 0; margin: 0; min-inline-size: 0; }
.block { display: grid; gap: 16px; }
.section-head {
  margin: 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.block-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
}
.add {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  padding: 4px 12px;
  border: 1px solid rgb(var(--ww-vermilion) / .3);
  transition: background-color .25s ease, color .25s ease;
}
.add:hover { background: rgb(var(--ww-vermilion) / .12); }

.grid { display: grid; gap: 18px 28px; }
.grid.two { grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); }

.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}

.rows { list-style: none; margin: 0; padding: 0; display: grid; gap: 8px; }
.row {
  display: grid;
  grid-template-columns: 28px 1fr auto auto;
  gap: 12px;
  align-items: end;
  padding: 8px 0;
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.era-row, .moon-row { grid-template-columns: 1fr auto auto auto; }

.idx {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  color: var(--ww-ink-faint);
  text-align: center;
  align-self: center;
}
.days { display: inline-flex; align-items: baseline; gap: 4px; }
.suffix { display: inline-block; }
.seg { display: grid; gap: 2px; }
.tiny { font-size: 9px; }
.ww-input.small { width: 96px; }
.ww-input.flex-1 { flex: 1; }

.row-actions { display: inline-flex; gap: 6px; align-self: center; }
.iconbtn {
  width: 26px; height: 26px;
  display: inline-flex; align-items: center; justify-content: center;
  border: 1px solid var(--ww-ink-hairline);
  background: transparent;
  color: var(--ww-ink-faint);
  transition: border-color .25s, color .25s, background-color .25s;
}
.iconbtn:hover:not(:disabled) {
  color: rgb(var(--ww-gold-deep));
  border-color: rgb(var(--ww-gold));
  background: rgb(var(--ww-gold) / .12);
}
.iconbtn:disabled { opacity: .3; cursor: not-allowed; }
.iconbtn.del:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .1); }

.actions {
  display: flex;
  align-items: center;
  gap: 16px;
  justify-content: flex-end;
  padding-top: 24px;
  margin-top: 20px;
  border-top: 1px solid var(--ww-ink-hairline);
}
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }
.saved {
  font-style: italic;
  color: rgb(var(--ww-gold-deep));
  font-size: 14px;
}

.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  margin: 12px 0 0;
}

.fade-enter-active, .fade-leave-active { transition: opacity .35s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.row-enter-active, .row-leave-active { transition: opacity .35s ease, transform .35s cubic-bezier(.22,1,.36,1); }
.row-enter-from { opacity: 0; transform: translateY(-4px); }
.row-leave-to { opacity: 0; transform: translateX(8px); }
.row-leave-active { position: absolute; }

.moons { display: grid; gap: 18px; }
.moon-card {
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment-deep) / .25);
  padding: 18px 20px 16px;
}
.moon-head {
  display: grid;
  grid-template-columns: 40px 1fr auto auto auto;
  gap: 14px;
  align-items: end;
  padding-bottom: 14px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.name-input { font-family: 'Fraunces', serif; font-size: 18px; }

.cycle-preview {
  margin-top: 16px;
  padding: 0 8px;
  border-top: 1px dashed var(--ww-ink-hairline);
}
.phases-block { padding-top: 14px; }
.phases-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 12px;
}
.add.small {
  padding: 3px 10px;
  font-size: 10px;
}
.phases {
  list-style: none; margin: 0; padding: 0;
  display: grid; gap: 10px;
}
.phase-row {
  display: grid;
  grid-template-columns: 36px 56px 1fr minmax(160px, 1.2fr) minmax(140px, 1fr) auto;
  align-items: end;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px dashed rgb(var(--ww-ink-hairline) / .6);
}
@media (max-width: 880px) {
  .phase-row {
    grid-template-columns: 36px 56px 1fr auto;
  }
  .phase-row .slider-seg { grid-column: 1 / -1; }
}
.phase-empty {
  padding: 14px 0;
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  font-size: 13px;
}

.icon-input {
  width: 100%;
  text-align: center;
  font-size: 18px;
  padding: 6px 0;
}

.slider-seg { display: grid; gap: 4px; min-width: 0; }
.slider-row {
  display: grid;
  grid-template-columns: 1fr 60px;
  gap: 8px;
  align-items: center;
}
.range {
  appearance: none;
  width: 100%;
  height: 22px;
  background: transparent;
  cursor: pointer;
}
.range::-webkit-slider-runnable-track {
  height: 4px;
  background: linear-gradient(to right,
    rgb(var(--ww-ink-shade) / .25) 0%,
    rgb(var(--ww-gold)) 50%,
    rgb(var(--ww-ink-shade) / .25) 100%);
  border: 1px solid var(--ww-ink-hairline);
}
.range::-moz-range-track {
  height: 4px;
  background: linear-gradient(to right,
    rgb(var(--ww-ink-shade) / .25) 0%,
    rgb(var(--ww-gold)) 50%,
    rgb(var(--ww-ink-shade) / .25) 100%);
  border: 1px solid var(--ww-ink-hairline);
}
.range::-webkit-slider-thumb {
  appearance: none;
  width: 16px; height: 16px;
  border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 2px 4px rgb(0 0 0 / .25), 0 0 0 3px rgb(var(--ww-parchment));
  border: 0;
  margin-top: -7px;
  transition: transform .2s ease;
}
.range::-webkit-slider-thumb:hover { transform: scale(1.18); }
.range::-moz-range-thumb {
  width: 16px; height: 16px;
  border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 2px 4px rgb(0 0 0 / .25), 0 0 0 3px rgb(var(--ww-parchment));
  border: 0;
  transition: transform .2s ease;
}
.range::-moz-range-thumb:hover { transform: scale(1.18); }

.ww-input.mini {
  width: 60px;
  padding: 4px 4px;
  font-size: 13px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
</style>
