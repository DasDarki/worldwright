<script setup lang="ts">
import type { Entity, EntitySummary, EntityType, FieldDefinition, Visibility } from '~/types/api'
import { slugify } from '~/composables/useSlugify'

const props = defineProps<{
  initial?: Entity | null
  submitting?: boolean
  submitLabel?: string
}>()

interface EntityFormPayload {
  entity_type_id: number
  title: string
  slug: string
  summary: string
  body: unknown
  parent_id: number | null
  visibility: Visibility
  tags: string[]
  field_values: Record<string, string>
}

const emit = defineEmits<{
  submit: [value: EntityFormPayload]
}>()

const { t, locale } = useI18n()
const { $api } = useNuxtApp()

const { data: typesData } = await useAsyncData('form-types', () =>
  $api<{ entity_types: EntityType[] }>('/entity-types'),
)
const { data: entitiesData } = await useAsyncData('form-entities', () =>
  $api<{ entities: EntitySummary[] }>('/entities'),
)

const types = computed(() => typesData.value?.entity_types || [])
const allEntities = computed(() => entitiesData.value?.entities || [])

const entityTypeId = ref<number>(props.initial?.entity_type_id ?? types.value[0]?.id ?? 0)
const title = ref(props.initial?.title ?? '')
const slug = ref(props.initial?.slug ?? '')
const summary = ref(props.initial?.summary ?? '')
const body = ref<unknown>(props.initial?.body ?? { type: 'doc', content: [{ type: 'paragraph' }] })
const parentId = ref<number | null>(props.initial?.parent_id ?? null)
const visibility = ref<Visibility>(props.initial?.visibility ?? 'secret')
const tags = ref<string[]>([...(props.initial?.tags ?? [])])
const fieldValues = ref<Record<string, string>>({ ...(props.initial?.field_values ?? {}) })

const slugTouched = ref(!!props.initial)
watch(title, (t) => {
  if (!slugTouched.value) slug.value = slugify(t)
})

const currentType = computed(() => types.value.find((t) => t.id === entityTypeId.value) || null)

function typeLabel(ty: EntityType) {
  return locale.value === 'de' ? ty.name_de : ty.name_en
}

function fieldLabel(fd: FieldDefinition) {
  return locale.value === 'de' ? fd.label_de : fd.label_en
}

function fieldOptions(fd: FieldDefinition): string[] {
  const cfg = fd.config as any
  if (Array.isArray(cfg?.options)) return cfg.options
  return []
}

function entityRefOptions(fd: FieldDefinition): EntitySummary[] {
  const cfg = fd.config as any
  const targetTypeKey = cfg?.target_type as string | undefined
  if (!targetTypeKey) return allEntities.value
  const targetType = types.value.find((t) => t.key === targetTypeKey)
  if (!targetType) return allEntities.value
  return allEntities.value.filter((e) => e.entity_type_id === targetType.id)
}

const parentOptions = computed(() => {
  return allEntities.value.filter((e) => !props.initial || e.id !== props.initial.id)
})

function onSubmit() {
  if (!title.value || !slug.value || !entityTypeId.value) return
  emit('submit', {
    entity_type_id: entityTypeId.value,
    title: title.value,
    slug: slug.value,
    summary: summary.value,
    body: body.value,
    parent_id: parentId.value,
    visibility: visibility.value,
    tags: tags.value,
    field_values: fieldValues.value,
  })
}
</script>

<template>
  <form class="entity-form stagger" @submit.prevent="onSubmit">
    <section class="block">
      <h3 class="ww-label">Essentials</h3>
      <label class="field">
        <span class="lbl">Title</span>
        <input v-model="title" class="ww-input" required maxlength="120" />
      </label>
      <label class="field">
        <span class="lbl">Slug</span>
        <input
          v-model="slug"
          class="ww-input mono"
          required
          pattern="[a-z0-9\-]+"
          @input="slugTouched = true"
        />
      </label>
      <label class="field">
        <span class="lbl">Summary</span>
        <input v-model="summary" class="ww-input" maxlength="300" placeholder="A single line a player will remember." />
      </label>
    </section>

    <section class="block">
      <h3 class="ww-label">Type & taxonomy</h3>
      <div class="row two">
        <label class="field">
          <span class="lbl">Type</span>
          <select v-model="entityTypeId" class="ww-input">
            <option v-for="ty in types" :key="ty.id" :value="ty.id">{{ typeLabel(ty) }}</option>
          </select>
        </label>
        <label class="field">
          <span class="lbl">Visibility</span>
          <select v-model="visibility" class="ww-input">
            <option value="secret">secret · author only</option>
            <option value="player">player · all readers</option>
            <option value="public">public · share link</option>
          </select>
        </label>
      </div>
      <label class="field">
        <span class="lbl">Parent</span>
        <select v-model="parentId" class="ww-input">
          <option :value="null">— none —</option>
          <option v-for="e in parentOptions" :key="e.id" :value="e.id">{{ e.title }}</option>
        </select>
      </label>
      <label class="field">
        <span class="lbl">Tags</span>
        <TagInput v-model="tags" />
      </label>
    </section>

    <section v-if="(currentType?.fields || []).length" class="block">
      <h3 class="ww-label">Properties · {{ typeLabel(currentType!) }}</h3>
      <div class="grid">
        <label v-for="fd in currentType!.fields" :key="fd.id" class="field">
          <span class="lbl">{{ fieldLabel(fd) }}<sup v-if="fd.is_required">*</sup></span>

          <textarea
            v-if="fd.data_type === 'longtext'"
            v-model="fieldValues[fd.key]"
            class="ww-input"
            rows="3"
          />
          <input
            v-else-if="fd.data_type === 'number'"
            v-model="fieldValues[fd.key]"
            type="number"
            class="ww-input"
          />
          <select
            v-else-if="fd.data_type === 'select'"
            v-model="fieldValues[fd.key]"
            class="ww-input"
          >
            <option value="">—</option>
            <option v-for="opt in fieldOptions(fd)" :key="opt" :value="opt">{{ opt }}</option>
          </select>
          <select
            v-else-if="fd.data_type === 'entity_ref'"
            v-model="fieldValues[fd.key]"
            class="ww-input"
          >
            <option value="">—</option>
            <option v-for="e in entityRefOptions(fd)" :key="e.id" :value="e.slug">{{ e.title }}</option>
          </select>
          <input
            v-else
            v-model="fieldValues[fd.key]"
            class="ww-input"
            :type="fd.data_type === 'date' ? 'text' : 'text'"
            :placeholder="fd.data_type === 'date' ? '14th of Miravel, 1142 LR' : ''"
          />
        </label>
      </div>
    </section>

    <section class="block body-block">
      <h3 class="ww-label">Body</h3>
      <ClientOnly>
        <LazyEditor v-model="body" />
        <template #fallback>
          <div class="editor-fallback">Loading the inkwell…</div>
        </template>
      </ClientOnly>
    </section>

    <div class="actions">
      <button type="submit" class="ww-btn-primary" :disabled="submitting">
        {{ submitLabel || 'Save' }}
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
.entity-form {
  display: grid;
  gap: 50px;
  max-width: 900px;
}

.block { display: grid; gap: 18px; }
.block h3 {
  margin: 0 0 4px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}

.row { display: grid; gap: 18px; }
.row.two { grid-template-columns: 1fr 1fr; }
@media (max-width: 640px) { .row.two { grid-template-columns: 1fr; } }

.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.field sup { color: rgb(var(--ww-vermilion)); margin-left: 2px; }

.mono { font-family: 'JetBrains Mono', ui-monospace, monospace; font-size: 16px; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px 28px;
}

select.ww-input {
  appearance: none;
  background-color: transparent;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23173842'/></svg>");
  background-repeat: no-repeat;
  background-position: right 0.3em center;
  background-size: 8px 5px;
  padding-right: 1.5em;
}
:root.dark select.ww-input {
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23ecdfc2'/></svg>");
}

textarea.ww-input {
  border: 1px solid var(--ww-ink-hairline);
  padding: 10px 12px;
  resize: vertical;
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
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--ww-ink-hairline);
}
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }
</style>
