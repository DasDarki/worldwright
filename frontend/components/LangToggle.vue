<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()
</script>

<template>
  <div class="lang-toggle">
    <button
      v-for="l in (locales as { code: string; name: string }[])"
      :key="l.code"
      type="button"
      :class="['lang-btn', { active: locale === l.code }]"
      :aria-pressed="locale === l.code"
      @click="setLocale(l.code as any)"
    >
      {{ l.code.toUpperCase() }}
    </button>
  </div>
</template>

<style scoped lang="scss">
.lang-toggle {
  display: inline-flex;
  align-items: center;
  gap: 0;
  font-family: 'Cormorant SC', serif;
  letter-spacing: 0.28em;
  font-size: 10px;
  border: 1px solid var(--ww-ink-hairline);
  padding: 0 2px;
}
.lang-btn {
  padding: 8px 10px;
  color: rgb(var(--ww-ink) / .55);
  transition: color .3s ease;
  &:hover { color: rgb(var(--ww-gold-deep)); }
  &.active { color: rgb(var(--ww-ink)); }
}
</style>
