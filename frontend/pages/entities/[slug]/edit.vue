<script setup lang="ts">
import type { Entity } from '~/types/api'

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { $api } = useNuxtApp()

const slug = computed(() => route.params.slug as string)

const { data, error: fetchError } = await useAsyncData(`edit-entity-${slug.value}`, () =>
  $api<{ entity: Entity }>(`/entities/by-slug/${slug.value}`),
)

if (fetchError.value && (fetchError.value as any).statusCode === 404) {
  throw createError({ statusCode: 404, statusMessage: 'Entry not found' })
}

const entity = computed(() => data.value?.entity || null)

useHead({ title: () => entity.value ? `${t('editor.editingPrefix')} ${entity.value.title}` : t('editor.editTitle') })

const pending = ref(false)
const error = ref<string | null>(null)

async function onSubmit(payload: any) {
  if (!entity.value) return
  pending.value = true
  error.value = null
  try {
    const res = await $api<{ entity: { slug: string } }>(`/entities/${entity.value.id}`, {
      method: 'PATCH',
      body: payload,
    })
    await router.push(`/entities/${res.entity.slug}`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to save changes'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section v-if="entity" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <NuxtLink :to="`/entities/${entity.slug}`" class="ww-btn-ghost back">
        <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
          <path d="M13 5 H1 M5 1 L1 5 L5 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        {{ t('entity.back') }}
      </NuxtLink>

      <div class="stagger mt-8 mb-10">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('editor.eyebrowEdit') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('editor.editingPrefix') }}</em> {{ entity.title }}
        </h1>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>

      <EntityForm
        :initial="entity"
        :submitting="pending"
        :submit-label="t('editor.save')"
        @submit="onSubmit"
      />

      <div class="ww-rule my-16" />

      <RelationshipEditor :entity-id="entity.id" />
    </div>
  </section>
</template>

<style scoped lang="scss">
.back { margin-bottom: 14px; }
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(36px, 5vw, 64px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  margin-bottom: 16px;
}
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
