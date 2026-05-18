<script setup lang="ts">
import type { EntitySummary, MapAsset, MapPin } from '~/types/api'

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const id = computed(() => Number(route.params.id))

const { data, refresh } = await useAsyncData(`map-edit-${id.value}`, () =>
  $api<{ map: MapAsset }>(`/maps/${id.value}`),
)

const map = ref<MapAsset | null>(data.value?.map ?? null)
useHead({ title: () => map.value ? `${t('maps.editingPrefix')} ${map.value.name}` : t('maps.editTitle') })

const selectedPin = ref<MapPin | null>(null)
const dirty = ref<Set<number>>(new Set())
const error = ref<string | null>(null)
const saving = ref(false)

const { data: maps2 } = await useAsyncData('all-maps', () =>
  $api<{ maps: MapAsset[] }>('/maps'),
)
const otherMaps = computed(() => (maps2.value?.maps || []).filter((m) => m.id !== id.value))

async function placePin(pos: { x: number; y: number }) {
  if (!map.value) return
  try {
    const { pin } = await $api<{ pin: MapPin }>(`/maps/${map.value.id}/pins`, {
      method: 'POST',
      body: { x: pos.x, y: pos.y, visibility: 'secret', label: '' },
    })
    map.value.pins = [...map.value.pins, pin]
    selectedPin.value = pin
  } catch (e: any) {
    error.value = e?.data?.error || 'Could not add pin'
  }
}

function movePin(pin: MapPin, pos: { x: number; y: number }) {
  if (!map.value) return
  const found = map.value.pins.find((p) => p.id === pin.id)
  if (!found) return
  found.x = pos.x
  found.y = pos.y
  dirty.value.add(pin.id)
  if (selectedPin.value?.id === pin.id) {
    selectedPin.value = { ...selectedPin.value, x: pos.x, y: pos.y }
  }
}

async function selectPin(pin: MapPin) {
  if (dirty.value.has(pin.id)) {
    await savePin(pin)
  }
  selectedPin.value = pin
}

async function savePin(pin: MapPin) {
  if (!map.value) return
  try {
    await $api(`/maps/${map.value.id}/pins/${pin.id}`, {
      method: 'PATCH',
      body: {
        x: pin.x, y: pin.y,
        label: pin.label || '',
        icon: pin.icon || '',
        target_entity_id: pin.target_entity_id ?? null,
        target_map_id: pin.target_map_id ?? null,
        visibility: pin.visibility,
      },
    })
    dirty.value.delete(pin.id)
  } catch (e: any) {
    error.value = e?.data?.error || 'Could not save pin'
  }
}

async function saveSelected() {
  if (!selectedPin.value) return
  saving.value = true
  try {
    const idx = map.value!.pins.findIndex((p) => p.id === selectedPin.value!.id)
    if (idx >= 0) map.value!.pins[idx] = { ...selectedPin.value }
    await savePin(selectedPin.value)
  } finally {
    saving.value = false
  }
}

async function removeSelected() {
  if (!selectedPin.value || !map.value) return
  if (!confirm(t('maps.confirmDeletePin'))) return
  await $api(`/maps/${map.value.id}/pins/${selectedPin.value.id}`, { method: 'DELETE' })
  map.value.pins = map.value.pins.filter((p) => p.id !== selectedPin.value!.id)
  selectedPin.value = null
}

const linkedEntity = ref<EntitySummary | null>(null)
watch(selectedPin, async (pin) => {
  if (!pin?.target_entity_id) {
    linkedEntity.value = null
    return
  }
  try {
    const all = await $api<{ entities: EntitySummary[] }>('/entities')
    linkedEntity.value = all.entities.find((e) => e.id === pin.target_entity_id) || null
  } catch {
    linkedEntity.value = null
  }
}, { immediate: true })

watch(linkedEntity, (e) => {
  if (!selectedPin.value) return
  const next = e?.id ?? null
  if (selectedPin.value.target_entity_id === next) return
  selectedPin.value = { ...selectedPin.value, target_entity_id: next }
})

async function deleteMap() {
  if (!map.value) return
  if (!confirm(t('maps.confirmDeleteMap'))) return
  await $api(`/maps/${map.value.id}`, { method: 'DELETE' })
  await router.push('/maps')
}
</script>

<template>
  <section v-if="map" class="py-12 md:py-16">
    <div class="mx-auto max-w-screen-2xl">
      <div class="topbar">
        <NuxtLink :to="`/maps/${map.id}`" class="ww-btn-ghost back">
          <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
            <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('maps.backView') }}
        </NuxtLink>
        <button type="button" class="ww-btn-ghost destroy" @click="deleteMap">{{ t('maps.deleteMap') }}</button>
      </div>

      <div class="stagger mt-6 mb-8">
        <div class="ww-eyebrow mb-4">{{ t('maps.eyebrowEdit') }}</div>
        <h1 class="title"><em>{{ t('maps.editingPrefix') }}</em> {{ map.name }}</h1>
        <p class="hint">{{ t('maps.editHint') }}</p>
      </div>

      <div class="layout">
        <MapCanvas
          :asset-id="map.asset_id"
          :pins="map.pins"
          :editable="true"
          :selected-pin-id="selectedPin?.id ?? null"
          @place-pin="placePin"
          @select-pin="selectPin"
          @move-pin="movePin"
        />

        <aside v-if="selectedPin" class="pin-editor ww-panel">
          <div class="ph">
            <h4 class="ww-label">{{ t('maps.pinEditor') }}</h4>
            <button type="button" class="x" @click="selectedPin = null">×</button>
          </div>
          <label class="field">
            <span class="lbl">{{ t('maps.label') }}</span>
            <input v-model="selectedPin.label" class="ww-input" />
          </label>
          <label class="field">
            <span class="lbl">{{ t('maps.visibility') }}</span>
            <select v-model="selectedPin.visibility" class="ww-input">
              <option value="secret">secret</option>
              <option value="player">player</option>
              <option value="public">public</option>
            </select>
          </label>
          <div class="field">
            <span class="lbl">{{ t('maps.linkEntity') }}</span>
            <EntityPicker v-model="linkedEntity" />
          </div>
          <label class="field">
            <span class="lbl">{{ t('maps.linkMap') }}</span>
            <select v-model.number="selectedPin.target_map_id" class="ww-input">
              <option :value="null">—</option>
              <option v-for="m in otherMaps" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </label>

          <div class="pe-actions">
            <button type="button" class="ww-btn-ghost destroy" @click="removeSelected">{{ t('maps.removePin') }}</button>
            <button type="button" class="ww-btn-primary" :disabled="saving" @click="saveSelected">
              {{ saving ? t('common.loading') : t('maps.savePin') }}
            </button>
          </div>
        </aside>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>
    </div>
  </section>
</template>

<style scoped lang="scss">
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(36px, 5vw, 64px);
  line-height: .92;
  letter-spacing: -0.025em;
  margin: 0;
}
.title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.hint {
  font-style: italic;
  font-size: 14px;
  color: var(--ww-ink-faint);
}

.layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
  align-items: start;
  @media (max-width: 1100px) { grid-template-columns: 1fr; }
}

.pin-editor {
  position: sticky;
  top: 96px;
  display: grid;
  gap: 14px;
}
.ph { display: flex; align-items: center; justify-content: space-between; }
.x { background: transparent; border: 0; font-size: 20px; color: var(--ww-ink-faint); cursor: pointer; }
.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.pe-actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 6px;
  padding-top: 14px;
  border-top: 1px dashed var(--ww-ink-hairline);
}
.destroy { color: rgb(var(--ww-vermilion) / .8); border-color: rgb(var(--ww-vermilion) / .4); }
.destroy:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); background: rgb(var(--ww-vermilion) / .08); }

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

.error { color: rgb(var(--ww-vermilion)); margin-top: 16px; font-style: italic; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
