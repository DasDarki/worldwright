<script setup lang="ts">
const { t } = useI18n()

const revealed = ref(false)

function toggle() {
  revealed.value = !revealed.value
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    toggle()
  }
}
</script>

<template>
  <div
    :class="['spoiler', { revealed }]"
    role="button"
    tabindex="0"
    :aria-pressed="revealed"
    :aria-label="revealed ? t('spoiler.hide') : t('spoiler.reveal')"
    @click="toggle"
    @keydown="onKey"
  >
    <span class="badge">
      <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round">
        <path d="M2 12 C 5 7, 9 5, 12 5 C 15 5, 19 7, 22 12 C 19 17, 15 19, 12 19 C 9 19, 5 17, 2 12 Z"/>
        <circle cx="12" cy="12" r="3"/>
      </svg>
      <span>{{ revealed ? t('spoiler.shown') : t('spoiler.hidden') }}</span>
    </span>
    <div class="inner">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
.spoiler {
  position: relative;
  padding: 18px 22px;
  margin: 1.2em 0;
  background: repeating-linear-gradient(
    45deg,
    rgb(var(--ww-ink) / .06),
    rgb(var(--ww-ink) / .06) 8px,
    rgb(var(--ww-ink) / .12) 8px,
    rgb(var(--ww-ink) / .12) 16px
  );
  border-left: 3px solid rgb(var(--ww-vermilion));
  cursor: pointer;
  outline: none;
  transition: background-color .45s cubic-bezier(.22,1,.36,1);
}
.spoiler:focus-visible { box-shadow: 0 0 0 2px rgb(var(--ww-gold)); }

.spoiler.revealed {
  background: rgb(var(--ww-parchment-deep) / .35);
}

.badge {
  position: absolute;
  top: -10px; left: 14px;
  display: inline-flex; align-items: center; gap: 6px;
  background: rgb(var(--ww-parchment));
  padding: 2px 10px;
  font-family: 'Cormorant SC', serif;
  font-size: 9px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: rgb(var(--ww-vermilion));
  border: 1px solid rgb(var(--ww-vermilion) / .4);
}

.inner {
  display: grid;
  grid-template-rows: 0fr;
  overflow: hidden;
  transition: grid-template-rows .45s cubic-bezier(.22,1,.36,1), opacity .35s ease;
  opacity: 0;
}
.spoiler.revealed .inner {
  grid-template-rows: 1fr;
  opacity: 1;
}
.inner > :first-child { min-height: 0; }
.inner :deep(> *) { min-height: 0; }
</style>
