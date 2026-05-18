<script setup lang="ts">
const colorMode = useColorMode()

function toggle() {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <button
    type="button"
    class="theme-toggle"
    :aria-label="colorMode.value === 'dark' ? 'Switch to light' : 'Switch to dark'"
    @click="toggle"
  >
    <ClientOnly>
      <svg v-if="colorMode.value === 'dark'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" aria-hidden="true">
        <circle cx="12" cy="12" r="4" />
        <path d="M12 2 V 4 M12 20 V 22 M2 12 H 4 M20 12 H 22 M4.9 4.9 L 6.3 6.3 M17.7 17.7 L 19.1 19.1 M4.9 19.1 L 6.3 17.7 M17.7 6.3 L 19.1 4.9" />
      </svg>
      <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
        <path d="M20.5 14.5 A8.5 8.5 0 1 1 9.5 3.5 a7 7 0 0 0 11 11 Z" />
        <circle cx="17" cy="5" r="0.6" fill="currentColor" />
        <circle cx="20" cy="9" r="0.6" fill="currentColor" />
      </svg>
      <template #fallback>
        <svg viewBox="0 0 24 24" width="16" height="16" />
      </template>
    </ClientOnly>
  </button>
</template>

<style scoped lang="scss">
.theme-toggle {
  width: 38px;
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--ww-ink-hairline);
  color: rgb(var(--ww-ink));
  transition: border-color .4s ease, color .4s ease, background-color .4s ease;
  svg { width: 16px; height: 16px; transition: transform .6s cubic-bezier(.65,0,.35,1); }
  &:hover {
    border-color: rgb(var(--ww-gold));
    color: rgb(var(--ww-gold-deep));
    background: rgb(var(--ww-gold) / .12);
    svg { transform: rotate(20deg); }
  }
}
</style>
