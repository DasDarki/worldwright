<script setup lang="ts">
import type { Backlink, Entity } from '~/types/api'

definePageMeta({ layout: 'share' })

const route = useRoute()
const { t, locale } = useI18n()
const { $api } = useNuxtApp()
const config = useRuntimeConfig()

const slug = computed(() => route.params.slug as string)

const { data, error } = await useAsyncData(`share-entity-${slug.value}`, () =>
  $api<{ entity: Entity; backlinks: Backlink[] }>(`/share/entities/${slug.value}`),
)

if (error.value && (error.value as any).statusCode === 404) {
  throw createError({ statusCode: 404, statusMessage: 'Entry not found' })
}

const entity = computed(() => data.value?.entity || null)
const backlinks = computed(() => data.value?.backlinks || [])

const typeName = computed(() => {
  if (!entity.value?.entity_type) return ''
  return locale.value === 'de' ? entity.value.entity_type.name_de : entity.value.entity_type.name_en
})

const siteName = computed(() => (config.public.siteName as string) || 'Worldwright')
const { absolute } = useSiteUrl()
const ogImageUrl = computed(() => absolute('/banner.png'))
const pageUrl = computed(() => entity.value ? absolute(`/share/entity/${entity.value.slug}`) : absolute('/'))

useHead({ title: () => entity.value?.title || t('common.loading') })
useSeoMeta({
  title: () => entity.value?.title,
  description: () => entity.value?.summary,
  ogTitle: () => entity.value ? `${entity.value.title} · ${siteName.value}` : siteName.value,
  ogDescription: () => entity.value?.summary,
  ogImage: () => ogImageUrl.value,
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: () => entity.value?.title || siteName.value,
  ogType: 'article',
  ogUrl: () => pageUrl.value,
  ogSiteName: () => siteName.value,
  ogLocale: () => locale.value === 'de' ? 'de_DE' : 'en_US',
  twitterCard: 'summary_large_image',
  twitterTitle: () => entity.value?.title,
  twitterDescription: () => entity.value?.summary,
  twitterImage: () => ogImageUrl.value,
})

const fields = computed(() => {
  if (!entity.value?.entity_type?.fields) return []
  return entity.value.entity_type.fields
    .map((fd) => ({
      key: fd.key,
      label: locale.value === 'de' ? fd.label_de : fd.label_en,
      value: entity.value!.field_values[fd.key],
    }))
    .filter((f) => f.value)
})

useReveal()
</script>

<template>
  <article v-if="entity" class="share-entity py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl px-6 md:px-12">
      <div class="stagger">
        <div class="entity-meta">
          <span class="type-chip">{{ typeName }}</span>
          <span class="ww-label public-badge">{{ t('share.publicBadge') }}</span>
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
      </div>

      <div v-if="backlinks.length" class="mentioned-in">
        <Ornament>{{ t('entity.mentionedIn') }}</Ornament>
        <ul class="bl-list">
          <li v-for="b in backlinks" :key="b.source_entity_id">
            <a :href="`/share/entity/${b.slug}`" class="ww-link">{{ b.title }}</a>
            <p v-if="b.summary" class="bl-summary">{{ b.summary }}</p>
          </li>
        </ul>
      </div>
    </div>
  </article>
</template>

<style scoped lang="scss">
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
.public-badge {
  color: rgb(var(--ww-gold-deep));
}

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

.mentioned-in { margin-top: 56px; }
.bl-list {
  list-style: none;
  margin: 30px 0 0;
  padding: 0;
  display: grid;
  gap: 16px;
  max-width: 36em;
}
.bl-summary {
  font-size: 13px;
  color: var(--ww-ink-faint);
  margin: 4px 0 0;
  font-style: italic;
}
</style>
