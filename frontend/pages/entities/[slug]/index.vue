<script setup lang="ts">
import type { Backlink, Entity, Genealogy } from '~/types/api'
import { useAuthStore } from '~/stores/auth'
import { highlightInDom } from '~/composables/useHighlightInDom'

const route = useRoute()
const { t, locale } = useI18n()
const { $api } = useNuxtApp()
const auth = useAuthStore()

const slug = computed(() => route.params.slug as string)

const { data, error } = await useAsyncData(
  () => `entity-${slug.value}`,
  () => $api<{ entity: Entity }>(`/entities/by-slug/${slug.value}`),
  { watch: [slug] },
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

const shareTempOpen = ref(false)
function printPage() {
  if (typeof window !== 'undefined') window.print()
}

if (import.meta.client) {
  // The previous Vue-watch-based approach raced BodyView's nested rendering:
  // by the time the watcher fired with the new entity, the BodyView subtree
  // sometimes hadn't been committed to the DOM yet (or had been committed
  // for the previous entity). Instead of trying to win that race, we just
  // re-poll the DOM up to ~1.5s every 80ms — whenever the markup appears
  // and contains a match, we highlight + scroll and stop polling. Each
  // navigation cancels the previous polling loop so we never highlight
  // against stale content.

  function clearOldHighlights(root: HTMLElement) {
    root.querySelectorAll('mark.search-hl').forEach((m) => {
      const parent = m.parentNode
      if (!parent) return
      parent.replaceChild(document.createTextNode(m.textContent || ''), m)
      parent.normalize?.()
    })
  }

  let pollHandle: number | null = null
  let pollGeneration = 0

  function stopPolling() {
    if (pollHandle != null) {
      window.clearTimeout(pollHandle)
      pollHandle = null
    }
  }

  function pollHighlight(query: string, generation: number, deadline: number) {
    if (generation !== pollGeneration) return
    const root = document.querySelector('.entity-article') as HTMLElement | null
    if (root) {
      clearOldHighlights(root)
      const first = highlightInDom(root, query)
      if (first) {
        first.scrollIntoView({ behavior: 'smooth', block: 'center' })
        return
      }
    }
    if (performance.now() < deadline) {
      pollHandle = window.setTimeout(() => pollHighlight(query, generation, deadline), 80)
    }
  }

  function triggerHighlight() {
    stopPolling()
    const q = route.query.q
    const query = typeof q === 'string' ? q : Array.isArray(q) ? String(q[0] ?? '') : ''
    if (!query) return
    const generation = ++pollGeneration
    const deadline = performance.now() + 1500
    // start immediately AND on the next paint, so a fast render path doesn't
    // wait a needless 80ms
    pollHighlight(query, generation, deadline)
    requestAnimationFrame(() => pollHighlight(query, generation, deadline))
  }

  const router = useRouter()
  let unhook: (() => void) | null = null

  onMounted(() => {
    triggerHighlight()
    unhook = router.afterEach(() => triggerHighlight())
  })

  onBeforeUnmount(() => {
    stopPolling()
    if (unhook) unhook()
  })
}
</script>

<template>
  <article v-if="entity" class="entity-article py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="topbar">
        <NuxtLink to="/entities" class="ww-btn-ghost back">
          <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
            <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
          {{ t('entity.back') }}
        </NuxtLink>
        <div class="actions">
          <ShareButton
            v-if="entity.visibility === 'public'"
            :url="`/share/entity/${entity.slug}`"
            class="share"
          />
          <button
            v-if="auth.isAdmin"
            type="button"
            class="ww-btn-ghost"
            :title="t('entity.shareTemp')"
            @click="shareTempOpen = true"
          >{{ t('entity.shareTemp') }}</button>
          <a
            :href="`/api/entities/by-slug/${entity.slug}/export.md`"
            class="ww-btn-ghost"
            :title="t('entity.exportMd')"
          >.md</a>
          <button
            type="button"
            class="ww-btn-ghost"
            :title="t('entity.printPdf')"
            @click="printPage"
          >PDF</button>
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
      </div>

      <ShareTokenDialog
        v-if="auth.isAdmin && entity"
        :open="shareTempOpen"
        :entity-id="entity.id"
        :entity-slug="entity.slug"
        @close="shareTempOpen = false"
      />

      <div class="grid lg:grid-cols-[1.6fr_1fr] gap-12 lg:gap-16 mt-10">
        <div class="stagger main-col">
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

          <RelationshipGraph
            v-if="hasGenealogy && genealogy"
            :entity-ids="genealogy.nodes.map((n) => n.id)"
            :title="t('graph.lineage')"
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
  flex-wrap: wrap;
}
.actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
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

/* CSS Grid columns default to min-content; long inline text or wide canvases
   in the left column will otherwise blow the layout out and push the side
   panel off-screen. min-width: 0 lets the column actually shrink. */
.main-col { min-width: 0; }
.side { min-width: 0; }

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
