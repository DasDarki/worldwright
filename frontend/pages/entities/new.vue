<script setup lang="ts">

definePageMeta({ middleware: 'admin-only' })

const { t } = useI18n()
useHead({ title: () => t('editor.newTitle') })

const router = useRouter()
const { $api } = useNuxtApp()

const pending = ref(false)
const error = ref<string | null>(null)

async function onSubmit(payload: any) {
  pending.value = true
  error.value = null
  try {
    const res = await $api<{ entity: { slug: string } }>('/entities', {
      method: 'POST',
      body: payload,
    })
    await router.push(`/entities/${res.entity.slug}`)
  } catch (e: any) {
    error.value = e?.data?.error || 'Failed to create entity'
  } finally {
    pending.value = false
  }
}

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-12">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('editor.eyebrowNew') }}
        </div>
        <h1 class="hero-title mb-4">
          <em>{{ t('editor.newTitle') }}</em>
        </h1>
        <p class="lede">{{ t('editor.newLede') }}</p>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error">{{ error }}</p>
      </Transition>

      <EntityForm
        :submitting="pending"
        :submit-label="t('editor.create')"
        @submit="onSubmit"
      />
    </div>
  </section>
</template>

<style scoped lang="scss">
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 6vw, 80px);
  line-height: .95;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 18px;
  max-width: 36em;
}
.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  margin-bottom: 16px;
}
.fade-enter-active, .fade-leave-active { transition: opacity .3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
