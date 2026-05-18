<script setup lang="ts">
import type { EntitySummary, RelationshipEdge, RelationshipType } from '~/types/api'

const props = defineProps<{ entityId: number }>()

const { t, locale } = useI18n()
const { $api } = useNuxtApp()

const { data: edgesData, refresh: refreshEdges } = await useAsyncData(
  () => `editor-rels-${props.entityId}`,
  () => $api<{ relationships: RelationshipEdge[] }>(`/entities/${props.entityId}/relationships`),
  { watch: [() => props.entityId] },
)

const { data: typesData } = await useAsyncData('rel-types', () =>
  $api<{ relationship_types: RelationshipType[] }>('/relationship-types'),
)

const edges = computed(() => edgesData.value?.relationships || [])
const types = computed(() => typesData.value?.relationship_types || [])

const draftTypeId = ref<number>(0)
const draftTarget = ref<EntitySummary | null>(null)
const draftDescription = ref('')
const draftError = ref<string | null>(null)
const pending = ref(false)

function typeLabel(rt: RelationshipType): string {
  const label = locale.value === 'de' ? rt.label_de : rt.label_en
  const inverse = locale.value === 'de' ? rt.inverse_label_de : rt.inverse_label_en
  if (rt.is_symmetric || !inverse || inverse === label) return label
  return `${label} / ${inverse}`
}

function edgeLabel(e: RelationshipEdge): string {
  const rt = e.type
  if (e.direction === 'out' || rt.is_symmetric) return locale.value === 'de' ? rt.label_de : rt.label_en
  return (locale.value === 'de' ? rt.inverse_label_de : rt.inverse_label_en) || rt.label_en
}

async function addRelationship() {
  draftError.value = null
  if (!draftTypeId.value || !draftTarget.value) {
    draftError.value = t('relationships.editor.missing')
    return
  }
  pending.value = true
  try {
    await $api('/relationships', {
      method: 'POST',
      body: {
        from_entity_id: props.entityId,
        to_entity_id: draftTarget.value.id,
        relationship_type_id: draftTypeId.value,
        description: draftDescription.value || undefined,
      },
    })
    draftTypeId.value = 0
    draftTarget.value = null
    draftDescription.value = ''
    await refreshEdges()
  } catch (e: any) {
    draftError.value = e?.data?.error || t('relationships.editor.failed')
  } finally {
    pending.value = false
  }
}

async function removeRelationship(id: number) {
  try {
    await $api(`/relationships/${id}`, { method: 'DELETE' })
    await refreshEdges()
  } catch {
    // ignore
  }
}
</script>

<template>
  <section class="rel-editor">
    <h3 class="ww-label section-head">{{ t('relationships.editor.title') }}</h3>

    <TransitionGroup name="row" tag="ul" class="edges">
      <li v-for="e in edges" :key="e.id" class="edge">
        <span class="role">{{ edgeLabel(e) }}</span>
        <NuxtLink :to="`/entities/${e.other.slug}`" class="link" target="_blank">{{ e.other.title }}</NuxtLink>
        <button
          type="button"
          class="del"
          :aria-label="t('relationships.editor.remove')"
          @click="removeRelationship(e.id)"
        >×</button>
      </li>
      <li v-if="!edges.length" key="empty" class="empty">{{ t('relationships.editor.empty') }}</li>
    </TransitionGroup>

    <div class="add">
      <div class="add-head">{{ t('relationships.editor.add') }}</div>
      <div class="add-grid">
        <label class="field">
          <span class="ww-label lbl">{{ t('relationships.editor.type') }}</span>
          <select v-model="draftTypeId" class="ww-input">
            <option :value="0">—</option>
            <optgroup
              v-for="cat in [...new Set(types.map((t) => t.category))]"
              :key="cat"
              :label="cat"
            >
              <option
                v-for="rt in types.filter((t) => t.category === cat)"
                :key="rt.id"
                :value="rt.id"
              >{{ typeLabel(rt) }}</option>
            </optgroup>
          </select>
        </label>
        <label class="field target">
          <span class="ww-label lbl">{{ t('relationships.editor.target') }}</span>
          <EntityPicker
            v-model="draftTarget"
            :exclude="[entityId]"
            :placeholder="t('relationships.editor.targetHint')"
          />
        </label>
      </div>
      <label class="field">
        <span class="ww-label lbl">{{ t('relationships.editor.description') }}</span>
        <input v-model="draftDescription" class="ww-input" :placeholder="t('relationships.editor.descHint')" />
      </label>

      <Transition name="fade">
        <p v-if="draftError" class="error">{{ draftError }}</p>
      </Transition>

      <div class="actions">
        <button
          type="button"
          class="ww-btn-primary"
          :disabled="pending || !draftTypeId || !draftTarget"
          @click="addRelationship"
        >
          {{ t('relationships.editor.bind') }}
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
.rel-editor {
  display: grid;
  gap: 24px;
  max-width: 900px;
}

.section-head {
  margin: 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}

.edges {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 6px;
}
.edge {
  display: grid;
  grid-template-columns: max-content 1fr max-content;
  gap: 8px 14px;
  align-items: baseline;
  padding: 8px 10px;
  background: rgb(var(--ww-parchment-deep) / .4);
  border: 1px solid var(--ww-ink-hairline);
  transition: background-color .3s ease;
}
.edge:hover { background: rgb(var(--ww-gold) / .08); }
.role {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
.link {
  font-family: 'EB Garamond', serif;
  font-size: 16px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.link:hover { color: rgb(var(--ww-vermilion)); }

.del {
  font-size: 18px;
  line-height: 1;
  color: var(--ww-ink-faint);
  padding: 2px 8px;
  background: transparent;
  border: 0;
  cursor: pointer;
  transition: color .25s ease;
}
.del:hover { color: rgb(var(--ww-vermilion)); }

.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  padding: 14px 0;
  text-align: center;
  font-size: 14px;
}

.add {
  border-top: 1px dashed var(--ww-ink-hairline);
  padding-top: 18px;
  display: grid;
  gap: 14px;
}
.add-head {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.add-grid {
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 18px;
  @media (max-width: 720px) { grid-template-columns: 1fr; }
}
.field { display: grid; gap: 4px; }
.lbl { margin: 0; }
.target { align-self: end; }

.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  margin: 0;
  font-size: 14px;
}

.actions {
  display: flex;
  justify-content: flex-end;
}
.actions .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.actions .ww-btn-primary:hover .arrow { transform: translateX(6px); }
.actions .ww-btn-primary:disabled { opacity: .5; cursor: not-allowed; }

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

.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
