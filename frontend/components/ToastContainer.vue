<script setup lang="ts">
import { useToastsStore } from '~/stores/toasts'

const toasts = useToastsStore()
</script>

<template>
  <Teleport to="body">
    <div class="toast-stack" aria-live="polite">
      <TransitionGroup name="toast" tag="div" class="stack-inner">
        <div
          v-for="t in toasts.items"
          :key="t.id"
          :class="['toast', t.kind]"
          role="status"
        >
          <span class="icon" aria-hidden="true">
            <svg v-if="t.kind === 'success'" viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <path d="M3 8 L7 12 L13 4"/>
            </svg>
            <svg v-else-if="t.kind === 'error'" viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <path d="M4 4 L12 12 M12 4 L4 12"/>
            </svg>
            <svg v-else viewBox="0 0 16 16" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <circle cx="8" cy="8" r="6"/>
              <path d="M8 7 L8 11 M8 5 L8 5.1"/>
            </svg>
          </span>
          <span class="message">{{ t.message }}</span>
          <button type="button" class="x" aria-label="Dismiss" @click="toasts.dismiss(t.id)">×</button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
.toast-stack {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 400;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: min(420px, calc(100% - 48px));
  pointer-events: none;
}
.stack-inner {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.toast {
  pointer-events: auto;
  display: grid;
  grid-template-columns: 20px 1fr auto;
  align-items: center;
  gap: 12px;
  padding: 12px 14px 12px 14px;
  background: rgb(var(--ww-parchment));
  border: 1px solid var(--ww-ink-hairline);
  box-shadow: 0 22px 40px -20px rgb(0 0 0 / .45);
  font-family: 'EB Garamond', serif;
  font-size: 15px;
  color: rgb(var(--ww-ink));
  position: relative;
  overflow: hidden;
}
.toast::before {
  content: '';
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 3px;
}
.toast.success::before { background: rgb(var(--ww-gold)); }
.toast.error::before   { background: rgb(var(--ww-vermilion)); }
.toast.info::before    { background: rgb(var(--ww-ink-shade)); }

.icon { display: inline-flex; align-items: center; justify-content: center; }
.toast.success .icon { color: rgb(var(--ww-gold-deep)); }
.toast.error   .icon { color: rgb(var(--ww-vermilion)); }
.toast.info    .icon { color: rgb(var(--ww-ink-shade)); }

.message {
  font-style: italic;
  line-height: 1.35;
}

.x {
  background: transparent;
  border: 0;
  font-size: 18px;
  line-height: 1;
  color: var(--ww-ink-faint);
  cursor: pointer;
  padding: 0 6px;
  transition: color .25s ease;
}
.x:hover { color: rgb(var(--ww-vermilion)); }

.toast-enter-active, .toast-leave-active {
  transition: opacity .35s ease, transform .4s cubic-bezier(.22,1,.36,1);
}
.toast-enter-from { opacity: 0; transform: translateX(20px); }
.toast-leave-to   { opacity: 0; transform: translateX(20px); }
.toast-move { transition: transform .4s cubic-bezier(.22,1,.36,1); }
</style>
