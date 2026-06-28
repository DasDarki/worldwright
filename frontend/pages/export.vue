<script setup lang="ts">
const { t } = useI18n()
const auth = useAuthStore()

useHead({ title: () => t('export.title') })

interface ExportOption {
  key: string
  href: string
  title: string
  body: string
  format: 'md' | 'pdf'
  shape: 'single' | 'archive'
}

const options = computed<ExportOption[]>(() => [
  {
    key: 'md-single',
    href: '/export.md',
    title: t('export.options.mdSingle.title'),
    body: t('export.options.mdSingle.body'),
    format: 'md',
    shape: 'single',
  },
  {
    key: 'md-zip',
    href: '/export.md.zip',
    title: t('export.options.mdZip.title'),
    body: t('export.options.mdZip.body'),
    format: 'md',
    shape: 'archive',
  },
  {
    key: 'pdf-single',
    href: '/export.pdf',
    title: t('export.options.pdfSingle.title'),
    body: t('export.options.pdfSingle.body'),
    format: 'pdf',
    shape: 'single',
  },
  {
    key: 'pdf-zip',
    href: '/export.pdf.zip',
    title: t('export.options.pdfZip.title'),
    body: t('export.options.pdfZip.body'),
    format: 'pdf',
    shape: 'archive',
  },
])

useReveal()
</script>

<template>
  <section class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl px-2">
      <div class="stagger mb-12">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          {{ t('export.eyebrow') }}
        </div>
        <h1 class="hero-title mb-4">{{ t('export.title') }}</h1>
        <p class="lede mb-2">{{ t('export.subtitle') }}</p>
        <p class="lede small">
          {{ auth.isAuthenticated
              ? t('export.scope.authenticated')
              : t('export.scope.anonymous') }}
        </p>
      </div>

      <div class="grid md:grid-cols-2 gap-6 lg:gap-8">
        <a
          v-for="opt in options"
          :key="opt.key"
          :href="opt.href"
          class="option reveal"
          :class="[`opt-${opt.format}`, `opt-${opt.shape}`]"
        >
          <header class="option-head">
            <span class="badge" :data-format="opt.format">{{ opt.format.toUpperCase() }}</span>
            <span class="shape">{{ t(`export.shape.${opt.shape}`) }}</span>
          </header>
          <h2 class="option-title">{{ opt.title }}</h2>
          <p class="option-body">{{ opt.body }}</p>
          <div class="option-foot">
            <code class="url">{{ opt.href }}</code>
            <span class="cta">
              {{ t('export.download') }}
              <svg width="14" height="10" viewBox="0 0 14 10" fill="none" aria-hidden="true">
                <path d="M1 5 H13 M9 1 L13 5 L9 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </span>
          </div>
        </a>
      </div>

      <p class="footnote">{{ t('export.footnote') }}</p>
    </div>
  </section>
</template>

<style scoped lang="scss">
.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 6vw, 76px);
  line-height: .98;
  letter-spacing: -.025em;
  margin: 0;
}
.lede {
  font-size: 18px;
  line-height: 1.55;
  color: var(--ww-ink-faint);
  max-width: 42em;
  margin: 0 0 .4em;
}
.lede.small {
  font-size: 14px;
  font-style: italic;
  color: rgb(var(--ww-gold-deep));
}

.option {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 24px 26px 22px;
  border: 1px solid var(--ww-ink-hairline);
  background: var(--ww-card-bg);
  color: rgb(var(--ww-ink));
  transition: transform .35s cubic-bezier(.22,1,.36,1), border-color .35s ease, box-shadow .35s ease;
}
.option:hover {
  border-color: rgb(var(--ww-gold));
  transform: translateY(-2px);
  box-shadow: 0 30px 60px -35px rgb(0 0 0 / .35);
}

.option-head {
  display: flex;
  align-items: baseline;
  gap: 14px;
}
.badge {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  padding: 3px 10px;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  background: rgb(var(--ww-vermilion) / .08);
}
.badge[data-format="pdf"] {
  color: rgb(var(--ww-gold-deep));
  background: rgb(var(--ww-gold) / .12);
}
.shape {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}

.option-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 460;
  font-size: 26px;
  line-height: 1.15;
  letter-spacing: -.01em;
  margin: 0;
}
.option-body {
  font-size: 15px;
  line-height: 1.55;
  color: rgb(var(--ww-ink-shade));
  margin: 0;
  flex: 1;
}

.option-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-top: 1px dashed var(--ww-ink-hairline);
  padding-top: 12px;
  margin-top: 4px;
  gap: 14px;
  flex-wrap: wrap;
}
.url {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px;
  color: var(--ww-ink-faint);
}
.cta {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .26em;
  text-transform: uppercase;
  color: rgb(var(--ww-gold-deep));
}
.option:hover .cta { color: rgb(var(--ww-vermilion)); }

.footnote {
  margin-top: 32px;
  font-style: italic;
  font-size: 13px;
  color: var(--ww-ink-faint);
  max-width: 44em;
}
</style>
