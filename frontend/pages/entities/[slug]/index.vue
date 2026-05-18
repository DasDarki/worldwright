<script setup lang="ts">
import type { Backlink, Entity, Genealogy } from '~/types/api'
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const { t, locale } = useI18n()
const { $api } = useNuxtApp()
const auth = useAuthStore()

const slug = computed(() => route.params.slug as string)

const { data, error } = await useAsyncData(`entity-${slug.value}`, () =>
  $api<{ entity: Entity }>(`/entities/by-slug/${slug.value}`)
)

if (error.value && (error.value as any).statusCode === 404) {
  throw createError({ statusCode: 404, statusMessage: 'Entry not found in the codex' })
}

const entity = computed(() => data.value?.entity || null)

const { data: backlinksData } = await useAsyncData(
  () => `backlinks-${entity.value?.id ?? 0}`,
  () => entity.value
    ? $api<{ backlinks: Backlink[] }>(`/entities/${entity.value.id}/backlinks`)
    : Promise.resolve({ backlinks: [] as Backlink[] }),
  { watch: [entity] }
)

const { data: genealogyData } = await useAsyncData(
  () => `genealogy-${entity.value?.id ?? 0}`,
  () => entity.value
    ? $api<{ genealogy: Genealogy }>(`/entities/${entity.value.id}/genealogy?depth=3`)
    : Promise.resolve({ genealogy: { focal: 0, nodes: [], edges: [] } as Genealogy }),
  { watch: [entity] }
)
const genealogy = computed<Genealogy | null>(() => genealogyData.value?.genealogy ?? null)
const hasGenealogy = computed(() => (genealogy.value?.nodes.length ?? 0) > 1)

const typeName = computed(() => {
  if (!entity.value?.entity_type) return ''
  return locale.value === 'de' ? entity.value.entity_type.name_de : entity.value.entity_type.name_en
})

const { absolute } = useSiteUrl()
const ogImageUrl = computed(() => absolute('/banner.png'))
const pageUrl = computed(() => entity.value ? absolute(`/entities/${entity.value.slug}`) : absolute('/'))

useHead({ title: () => entity.value?.title || t('common.loading') })
useSeoMeta({
  title: () => entity.value?.title,
  description: () => entity.value?.summary,
  ogTitle: () => entity.value?.title,
  ogDescription: () => entity.value?.summary,
  ogImage: () => ogImageUrl.value,
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: () => entity.value?.title,
  ogType: 'article',
  ogUrl: () => pageUrl.value,
  twitterCard: 'summary_large_image',
  twitterTitle: () => entity.value?.title,
  twitterDescription: () => entity.value?.summary,
  twitterImage: () => ogImageUrl.value,
})

const fields = computed(() => {
  if (!entity.value?.entity_type?.fields) return []
  return entity.value.entity_type.fields
    .map(fd => ({
      key: fd.key,
      label: locale.value === 'de' ? fd.label_de : fd.label_en,
      value: entity.value!.field_values[fd.key],
    }))
    .filter(f => f.value)
})

useReveal()
</script>

<template>
  <article v-if="entity" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-2xl">
      <div class="topbar">
        <NuxtLink to="/entities" class="ww-btn-ghost back">
          <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
            <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('entity.back') }}
        </NuxtLink>
        <ShareButton
          v-if="entity.visibility === 'public'"
          :url="`/share/entity/${entity.slug}`"
          class="share"
        />
        <NuxtLink
          v-if="auth.isAdmin"
          :to="`/entities/${entity.slug}/edit`"
          class="ww-btn-ghost edit"
        >
          {{ t('entity.edit') }}
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
            <path d="M1 11 L3 11 L11 3 L9 1 L1 9 Z" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round"/>
          </svg>
        </NuxtLink>
      </div>

      <div class="grid lg:grid-cols-[1.6fr_1fr] gap-12 lg:gap-16 mt-10">
        <div class="stagger">
          <div class="entity-meta">
            <span class="type-chip">{{ typeName }}</span>
            <span class="ww-label visibility" :data-v="entity.visibility">{{ t(`entities.visibility.${entity.visibility}`) }}</span>
          </div>

          <h1 class="entity-title">{{ entity.title }}</h1>
          <p v-if="entity.summary" class="entity-summary">{{ entity.summary }}</p>

          <div v-if="fields.length" class="fields">
            <div v-for="f in fields" :key="f.key" class="field">
              <div class="ww-label mb-1">{{ f.label }}</div>
              <div class="field-val">{{ f.value }}</div>
            </div>
          </div>

          <BodyView :body="entity.body" class="body" />

          <div v-if="entity.tags.length" class="tags">
            <span v-for="tag in entity.tags" :key="tag" class="ww-tag">{{ tag }}</span>
          </div>

          <PedigreeChart
            v-if="hasGenealogy && genealogy"
            :genealogy="genealogy"
            :focal="entity.id"
            class="pedigree"
          />
        </div>

        <aside class="side reveal">
          <RelationshipPanel :entity-id="entity.id" />

          <EventsForEntityPanel :entity-id="entity.id" />

          <div class="ww-panel mentioned">
            <h4 class="ww-label mb-4">{{ t('entity.mentionedIn') }}</h4>
            <ul v-if="(backlinksData?.backlinks || []).length" class="backlinks">
              <li v-for="b in backlinksData!.backlinks" :key="b.source_entity_id">
                <NuxtLink :to="`/entities/${b.slug}`">{{ b.title }}</NuxtLink>
                <p v-if="b.summary" class="bl-summary">{{ b.summary }}</p>
              </li>
            </ul>
            <p v-else class="empty">{{ t('entity.noBacklinks') }}</p>
          </div>
        </aside>
      </div>
    </div>
  </article>
</template>

<style scoped lang="scss">
.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  gap: 14px;
}
.back, .edit { margin: 0; }

.entity-meta {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}
.type-chip {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  color: rgb(var(--ww-vermilion));
  background: rgb(var(--ww-vermilion) / .08);
  padding: 5px 12px;
  text-transform: uppercase;
}
.visibility { font-size: 10px; }
.visibility[data-v="secret"] { color: rgb(var(--ww-vermilion)); }
.visibility[data-v="player"] { color: rgb(var(--ww-gold-deep)); }
.visibility[data-v="public"] { color: rgb(var(--ww-ink) / .55); }

.entity-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: clamp(44px, 7vw, 96px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.entity-summary {
  font-style: italic;
  font-size: 20px;
  line-height: 1.45;
  color: rgb(var(--ww-ink-shade));
  margin: 24px 0 12px;
  max-width: 36em;
}

.fields {
  margin: 36px 0 32px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 22px 28px;
  padding: 28px 0;
  border-top: 1px dashed var(--ww-ink-hairline);
  border-bottom: 1px dashed var(--ww-ink-hairline);
}
.field-val { font-size: 16px; color: rgb(var(--ww-ink)); }

.body { margin-top: 30px; }

.tags { margin-top: 32px; display: flex; gap: 8px; flex-wrap: wrap; }
.pedigree { margin-top: 56px; }

.side {
  position: sticky;
  top: 96px;
  align-self: start;
  display: grid;
  gap: 18px;
}
.mentioned .backlinks {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 14px;
}
.mentioned a {
  font-family: 'EB Garamond', serif;
  font-size: 17px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
  transition: border-color .3s, color .3s;
}
.mentioned a:hover {
  color: rgb(var(--ww-vermilion));
  border-color: rgb(var(--ww-vermilion));
}
.mentioned .bl-summary {
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin: 4px 0 0;
  font-style: italic;
}
.empty {
  font-style: italic;
  color: var(--ww-ink-faint);
  margin: 0;
}
</style>
