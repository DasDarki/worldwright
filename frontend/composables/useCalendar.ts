import type { Calendar, CalendarMoon, CalendarMoonPhase, InWorldDate, WorldEvent } from '~/types/api'

let cachedPromise: Promise<Calendar[]> | null = null

export function useCalendars() {
  const { $api } = useNuxtApp()

  async function load(): Promise<Calendar[]> {
    if (!cachedPromise) {
      cachedPromise = $api<{ calendars: Calendar[] }>('/calendars').then((r) => r.calendars)
    }
    return cachedPromise
  }

  function invalidate() {
    cachedPromise = null
  }

  return { load, invalidate }
}

export function ordinal(n: number): string {
  const v = Math.abs(n) % 100
  if (v >= 11 && v <= 13) return `${n}th`
  switch (v % 10) {
    case 1: return `${n}st`
    case 2: return `${n}nd`
    case 3: return `${n}rd`
    default: return `${n}th`
  }
}

export function pickEra(cal: Calendar, year: number) {
  let chosen = cal.eras[0]
  for (const e of cal.eras) {
    if (e.start_year <= year && (!chosen || e.start_year > chosen.start_year)) {
      chosen = e
    }
  }
  return chosen
}

export function formatDate(cal: Calendar | null | undefined, date: InWorldDate | WorldEvent | null | undefined, locale: 'en' | 'de' = 'en'): string {
  if (!cal || !date) return ''
  const month = cal.months.find((m) => m.sort_order === date.month_index)
  const monthName = month?.name ?? `Month ${date.month_index}`
  const era = date.era_id
    ? cal.eras.find((e) => e.id === date.era_id)
    : pickEra(cal, date.year)
  const yearInEra = era ? date.year - era.start_year : date.year
  const dayPart = locale === 'de' ? `${date.day}.` : ordinal(date.day)
  const yearPart = era?.suffix ? `${yearInEra} ${era.suffix}` : `${yearInEra}`
  const joiner = locale === 'de' ? ' ' : ' of '
  return `${dayPart}${joiner}${monthName}, ${yearPart}`
}

export interface MoonState {
  moon: CalendarMoon
  cycle: number
  phase: CalendarMoonPhase | null
}

export function moonStates(cal: Calendar | null | undefined, day: number): MoonState[] {
  if (!cal) return []
  return cal.moons.map((m) => {
    const cycle = ((day - m.offset_days) % m.cycle_days + m.cycle_days) % m.cycle_days / m.cycle_days
    return { moon: m, cycle, phase: resolveMoonPhase(m, cycle) }
  })
}

export function moonStatesForDate(cal: Calendar | null | undefined, date: InWorldDate | null | undefined): MoonState[] {
  if (!cal || !date) return []
  return moonStates(cal, absoluteDay(cal, date))
}

function modularDistance(a: number, b: number): number {
  const d = Math.abs(a - b)
  return Math.min(d, 1 - d)
}

export function resolveMoonPhase(moon: CalendarMoon, cyclePosition: number): CalendarMoonPhase | null {
  const phases = moon.phases ?? []
  if (!phases.length) return null
  let best: { phase: CalendarMoonPhase; distance: number } | null = null
  for (const p of phases) {
    const d = modularDistance(cyclePosition, p.cycle_position)
    if (d <= p.randomness) {
      if (!best || d < best.distance) best = { phase: p, distance: d }
    }
  }
  if (best) return best.phase
  let closest: { phase: CalendarMoonPhase; distance: number } | null = null
  for (const p of phases) {
    const d = modularDistance(cyclePosition, p.cycle_position)
    if (!closest || d < closest.distance) closest = { phase: p, distance: d }
  }
  return closest?.phase ?? null
}

export function moonGlyphPath(cyclePosition: number, R = 48): string {
  if (cyclePosition < 0.005 || cyclePosition > 0.995) return ''
  const psi = cyclePosition * 2 * Math.PI
  const cos = Math.cos(psi)
  const x = Math.abs(R * cos)
  const outer = cyclePosition < 0.5 ? 1 : 0
  const inner = (cyclePosition < 0.25 || (cyclePosition >= 0.5 && cyclePosition < 0.75)) ? 1 : 0
  return `M 0 ${-R} A ${R} ${R} 0 0 ${outer} 0 ${R} A ${x} ${R} 0 0 ${inner} 0 ${-R} Z`
}

export function absoluteDay(cal: Calendar | null | undefined, date: InWorldDate): number {
  if (!cal) return 0
  const monthsInYear = cal.months.length
  const daysPerYear = cal.months.reduce((s, m) => s + m.days, 0)
  let daysIntoYear = 0
  for (const m of cal.months.slice(0, Math.max(0, date.month_index - 1))) {
    daysIntoYear += m.days
  }
  daysIntoYear += date.day - 1
  return date.year * daysPerYear + daysIntoYear
}
