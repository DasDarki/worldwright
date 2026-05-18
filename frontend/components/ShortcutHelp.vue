<script setup lang="ts">
const { t } = useI18n()

const open = ref(false)

const shortcuts = computed(() => [
  { keys: ['⌘', 'K'], altKeys: ['Ctrl', 'K'], label: t('shortcuts.search') },
  { keys: ['?'],                              label: t('shortcuts.help') },
  { keys: ['Esc'],                            label: t('shortcuts.dismiss') },
  { keys: ['Enter'],                          label: t('shortcuts.openHit') },
  { keys: ['↑', '↓'],                         label: t('shortcuts.navHits') },
])

if (import.meta.client) {
  onMounted(() => {
    function onKey(e: KeyboardEvent) {
      const tag = (e.target as HTMLElement | null)?.tagName
      const inField = tag === 'INPUT' || tag === 'TEXTAREA' || (e.target as HTMLElement | null)?.isContentEditable
      if (e.key === '?' && !e.metaKey && !e.ctrlKey && !inField) {
        e.preventDefault()
        open.value = !open.value
      } else if (e.key === 'Escape' && open.value) {
        open.value = false
      }
    }
    document.addEventListener('keydown', onKey)
    onBeforeUnmount(() => document.removeEventListener('keydown', onKey))
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="help">
      <div v-if="open" class="help-root" @click.self="open = false">
        <div class="card" role="dialog" :aria-label="t('shortcuts.title')">
          <div class="head">
            <div class="ww-eyebrow">{{ t('shortcuts.eyebrow') }}</div>
            <button type="button" class="x" aria-label="Close" @click="open = false">×</button>
          </div>
          <h3 class="title">{{ t('shortcuts.title') }}</h3>

          <ul class="list">
            <li v-for="(s, i) in shortcuts" :key="i" class="row">
              <span class="keys">
                <template v-for="(k, j) in s.keys" :key="`k-${j}`">
                  <kbd>{{ k }}</kbd>
                  <span v-if="j < s.keys.length - 1" class="plus">+</span>
                </template>
                <span v-if="s.altKeys" class="alt">
                  <span class="or">{{ t('shortcuts.or') }}</span>
                  <template v-for="(k, j) in s.altKeys" :key="`a-${j}`">
                    <kbd>{{ k }}</kbd>
                    <span v-if="j < s.altKeys.length - 1" class="plus">+</span>
                  </template>
                </span>
              </span>
              <span class="label">{{ s.label }}</span>
            </li>
          </ul>

          <p class="footer">{{ t('shortcuts.tip') }}</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.help-root {
  position: fixed;
  inset: 0;
  z-index: 250;
  background: rgb(0 0 0 / .35);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.card {
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  padding: 28px 30px 24px;
  width: min(520px, 100%);
  box-shadow: 0 60px 100px -40px rgb(0 0 0 / .55);
  display: grid;
  gap: 18px;
}
.head { display: flex; align-items: baseline; justify-content: space-between; }
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 144, "wght" 400;
  font-size: 32px;
  line-height: 1;
  letter-spacing: -0.02em;
  margin: 0;
}
.x {
  background: transparent;
  border: 0;
  font-size: 22px;
  line-height: 1;
  cursor: pointer;
  color: var(--ww-ink-faint);
  transition: color .25s ease;
}
.x:hover { color: rgb(var(--ww-vermilion)); }

.list { list-style: none; margin: 0; padding: 0; display: grid; gap: 10px; }
.row {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 16px;
  align-items: center;
}
.keys { display: inline-flex; align-items: center; gap: 4px; flex-wrap: wrap; }
kbd {
  display: inline-block;
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .12em;
  color: rgb(var(--ww-ink));
  background: rgb(var(--ww-parchment-deep) / .6);
  border: 1px solid var(--ww-ink-hairline);
  padding: 3px 7px;
  min-width: 18px;
  text-align: center;
}
.plus {
  color: var(--ww-ink-faint);
  font-size: 10px;
}
.alt {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: 6px;
}
.or {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 11px;
  margin: 0 4px;
}
.label {
  font-family: 'EB Garamond', serif;
  font-size: 15px;
  color: rgb(var(--ww-ink));
}

.footer {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 13px;
  text-align: center;
  margin: 0;
  padding-top: 12px;
  border-top: 1px dashed var(--ww-ink-hairline);
}

.help-enter-active, .help-leave-active { transition: opacity .3s ease; }
.help-enter-active .card, .help-leave-active .card {
  transition: transform .4s cubic-bezier(.22,1,.36,1), opacity .35s ease;
}
.help-enter-from, .help-leave-to { opacity: 0; }
.help-enter-from .card, .help-leave-to .card {
  transform: scale(.96);
  opacity: 0;
}
</style>
