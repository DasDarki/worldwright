export interface OAuthIdentity {
  id: number
  user_id: number
  provider: string
  subject: string
  email?: string
  created_at: string
}

export interface User {
  id: number
  email: string
  has_password: boolean
  role: 'admin' | 'player'
  locale: string
  display_name?: string
  avatar_url?: string
  created_at: string
}

export interface FieldDefinition {
  id: number
  entity_type_id: number
  key: string
  label_en: string
  label_de: string
  data_type: string
  config?: unknown
  sort_order: number
  is_required: boolean
}

export interface EntityType {
  id: number
  key: string
  name_en: string
  name_de: string
  icon?: string
  color?: string
  is_system: boolean
  fields?: FieldDefinition[]
}

export type Visibility = 'secret' | 'player' | 'public'

export interface EntitySummary {
  id: number
  entity_type_id: number
  title: string
  slug: string
  summary?: string
  parent_id?: number
  visibility: Visibility
}

export interface Entity {
  id: number
  entity_type_id: number
  entity_type?: EntityType
  title: string
  slug: string
  summary?: string
  body: unknown
  parent_id?: number
  materialized_path: string
  visibility: Visibility
  created_at: string
  updated_at: string
  tags: string[]
  field_values: Record<string, string>
}

export interface TagWithCount {
  name: string
  count: number
}

export interface Backlink {
  source_entity_id: number
  title: string
  slug: string
  entity_type_id: number
  summary?: string
}

export interface AuthProviders {
  password: boolean
  oauth: string[]
}

export interface Asset {
  id: number
  filename: string
  mime: string
  size: number
  width?: number
  height?: number
  created_at: string
}

export interface MapPin {
  id: number
  map_id: number
  x: number
  y: number
  label?: string
  icon?: string
  target_entity_id?: number | null
  target_entity_slug?: string
  target_map_id?: number | null
  visibility: Visibility
}

export interface MapAsset {
  id: number
  name: string
  asset_id: number
  width: number
  height: number
  parent_entity_id?: number | null
  pins: MapPin[]
}

export interface CalendarMonth {
  id?: number
  calendar_id?: number
  sort_order: number
  name: string
  days: number
}

export interface CalendarWeekday {
  id?: number
  calendar_id?: number
  sort_order: number
  name: string
}

export interface CalendarEra {
  id?: number
  calendar_id?: number
  name: string
  start_year: number
  suffix?: string
}

export interface CalendarMoonPhase {
  id?: number
  moon_id?: number
  sort_order: number
  name: string
  cycle_position: number
  randomness: number
  icon?: string
}

export interface CalendarMoon {
  id?: number
  calendar_id?: number
  name: string
  cycle_days: number
  offset_days: number
  phases: CalendarMoonPhase[]
}

export interface Calendar {
  id: number
  name: string
  epoch_name?: string
  current_year: number
  current_month_index: number
  current_day: number
  leap_rule?: string
  months: CalendarMonth[]
  weekdays: CalendarWeekday[]
  eras: CalendarEra[]
  moons: CalendarMoon[]
}

export interface InWorldDate {
  calendar_id: number
  era_id?: number | null
  year: number
  month_index: number
  day: number
}

export interface EventParticipant {
  event_id?: number
  entity_id: number
  role: string
  entity?: EntitySummary
}

export interface WorldEvent {
  id: number
  title: string
  body: unknown
  calendar_id: number
  era_id?: number | null
  year: number
  month_index: number
  day: number
  end_year?: number | null
  end_month_index?: number | null
  end_day?: number | null
  importance: number
  visibility: Visibility
  participants: EventParticipant[]
}

export interface Timeline {
  id: number
  name: string
  description?: string
  event_ids: number[]
  events?: WorldEvent[]
}

export interface GenealogyNode {
  id: number
  entity_type_id: number
  title: string
  slug: string
  summary?: string
  visibility: Visibility
}

export interface GenealogyEdge {
  from: number
  to: number
  type: 'parent_of' | 'spouse_of' | 'sibling_of' | string
}

export interface Genealogy {
  focal: number
  nodes: GenealogyNode[]
  edges: GenealogyEdge[]
}

export interface RelationshipType {
  id: number
  key: string
  label_en: string
  label_de: string
  inverse_label_en?: string
  inverse_label_de?: string
  is_symmetric: boolean
  category: 'generic' | 'genealogy' | string
}

export interface RelationshipEdge {
  id: number
  type: RelationshipType
  direction: 'in' | 'out'
  other: EntitySummary
  description?: string
}

export interface TipTapNode {
  type: string
  text?: string
  attrs?: Record<string, unknown>
  content?: TipTapNode[]
  marks?: { type: string; attrs?: Record<string, unknown> }[]
}
