<script setup lang="ts">
import type { RelationshipEdge } from '~/types/api'

const props = defineProps<{ entityId: number }>()

const { t, locale } = useI18n()
const { $api } = useNuxtApp()

const { data } = await useAsyncData(
  () => `rels-${props.entityId}`,
  () => $api<{ relationships: RelationshipEdge[] }>(`/entities/${props.entityId}/relationships`),
  { watch: [() => props.entityId] },
)

const edges = computed(() => data.value?.relationships || [])

function labelFor(edge: RelationshipEdge): string {
  const t = edge.type
  if (edge.direction === 'out' || t.is_symmetric) {
    return locale.value === 'de' ? t.label_de : t.label_en
  }
  const inverse = locale.value === 'de' ? t.inverse_label_de : t.inverse_label_en
  return inverse || (locale.value === 'de' ? t.label_de : t.label_en)
}

const generic = computed(() => edges.value.filter((e) => e.type.category !== 'genealogy'))
const genealogy = computed(() => edges.value.filter((e) => e.type.category === 'genealogy'))

const family = computed(() => {
  const groups = { parents: [] as RelationshipEdge[], children: [] as RelationshipEdge[], spouses: [] as RelationshipEdge[], siblings: [] as RelationshipEdge[] }
  for (const e of genealogy.value) {
    if (e.type.key === 'parent_of') {
      if (e.direction === 'in') groups.parents.push(e)
      else groups.children.push(e)
    } else if (e.type.key === 'spouse_of') {
      groups.spouses.push(e)
    } else if (e.type.key === 'sibling_of') {
      groups.siblings.push(e)
    }
  }
  return groups
})
</script>

<template>
  <div v-if="edges.length" class="ww-panel rel-panel">
    <h4 class="ww-label mb-4">{{ t('relationships.title') }}</h4>

    <div v-if="genealogy.length" class="block genealogy">
      <div class="block-head">{{ t('relationships.genealogy') }}</div>

      <div v-if="family.parents.length" class="group">
        <span class="role">{{ t('relationships.parents') }}</span>
        <ul class="lineage">
          <li v-for="e in family.parents" :key="e.id">
            <NuxtLink :to="`/entities/${e.other.slug}`" class="link">{{ e.other.title }}</NuxtLink>
          </li>
        </ul>
      </div>

      <div v-if="family.spouses.length" class="group">
        <span class="role">{{ t('relationships.spouses') }}</span>
        <ul class="lineage">
          <li v-for="e in family.spouses" :key="e.id">
            <NuxtLink :to="`/entities/${e.other.slug}`" class="link">{{ e.other.title }}</NuxtLink>
          </li>
        </ul>
      </div>

      <div v-if="family.siblings.length" class="group">
        <span class="role">{{ t('relationships.siblings') }}</span>
        <ul class="lineage">
          <li v-for="e in family.siblings" :key="e.id">
            <NuxtLink :to="`/entities/${e.other.slug}`" class="link">{{ e.other.title }}</NuxtLink>
          </li>
        </ul>
      </div>

      <div v-if="family.children.length" class="group">
        <span class="role">{{ t('relationships.children') }}</span>
        <ul class="lineage">
          <li v-for="e in family.children" :key="e.id">
            <NuxtLink :to="`/entities/${e.other.slug}`" class="link">{{ e.other.title }}</NuxtLink>
          </li>
        </ul>
      </div>
    </div>

    <div v-if="generic.length" class="block">
      <div class="block-head">{{ t('relationships.generic') }}</div>
      <ul class="edges">
        <li v-for="e in generic" :key="e.id" class="edge">
          <span class="role">{{ labelFor(e) }}</span>
          <NuxtLink :to="`/entities/${e.other.slug}`" class="link">{{ e.other.title }}</NuxtLink>
          <p v-if="e.description" class="desc">{{ e.description }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped lang="scss">
.rel-panel { position: relative; }
.block { margin-top: 18px; }
.block:first-of-type { margin-top: 0; }
.block-head {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.group {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 8px 16px;
  align-items: baseline;
  margin: 8px 0;
}
.role {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.lineage {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 4px 14px;
}
.edges { list-style: none; margin: 0; padding: 0; display: grid; gap: 10px; }
.edge {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 4px 14px;
  align-items: baseline;
}
.edge .desc {
  grid-column: 2;
  font-size: 12px;
  color: var(--ww-ink-faint);
  font-style: italic;
  margin: 0;
}
.link {
  font-family: 'EB Garamond', serif;
  font-size: 16px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
  transition: border-color .3s ease, color .3s ease;
}
.link:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); }
</style>
