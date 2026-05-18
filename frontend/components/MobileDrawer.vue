<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const { t } = useI18n()
const auth = useAuthStore()
const router = useRouter()

async function logout() {
  emit('close')
  await auth.logout()
  router.push('/login')
}

const links = computed(() => [
  { to: '/',           label: t('nav.codex') },
  { to: '/entities',   label: t('nav.entities') },
  { to: '/tags',       label: t('nav.tags') },
  { to: '/events',     label: t('nav.events') },
  { to: '/timelines',  label: t('nav.timelines') },
  { to: '/maps',       label: t('nav.maps') },
  { to: '/calendars',  label: t('nav.calendars') },
  ...(auth.isAdmin ? [{ to: '/assets', label: t('nav.assets') }] : []),
  { to: '/account',    label: t('nav.account') },
])

if (import.meta.client) {
  watch(() => props.open, (open) => {
    document.documentElement.style.overflow = open ? 'hidden' : ''
  })
  onMounted(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape' && props.open) emit('close')
    }
    document.addEventListener('keydown', onKey)
    onBeforeUnmount(() => {
      document.removeEventListener('keydown', onKey)
      document.documentElement.style.overflow = ''
    })
  })
}
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="open" class="drawer-root">
        <div class="backdrop" @click="emit('close')" />
        <aside class="drawer" role="dialog" :aria-label="t('mobile.menu')">
          <div class="head">
            <span class="brand">
              <span class="dot" aria-hidden="true" />
              worldwright
            </span>
            <button type="button" class="close" :aria-label="t('mobile.close')" @click="emit('close')">×</button>
          </div>

          <nav class="links">
            <NuxtLink
              v-for="(l, i) in links"
              :key="l.to"
              :to="l.to"
              class="link"
              :style="{ animationDelay: `${0.04 * i + 0.08}s` }"
              @click="emit('close')"
            >
              <span class="num">{{ String(i + 1).padStart(2, '0') }}</span>
              <span class="label">{{ l.label }}</span>
              <span class="chevron">→</span>
            </NuxtLink>
          </nav>

          <div class="footer">
            <div class="row">
              <LangToggle />
              <ThemeToggle />
            </div>
            <button v-if="auth.isAuthenticated" type="button" class="ww-btn-ghost logout" @click="logout">
              {{ t('nav.signOut') }}
            </button>
          </div>
        </aside>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.drawer-root {
  position: fixed;
  inset: 0;
  z-index: 300;
  display: flex;
}
.backdrop {
  position: absolute;
  inset: 0;
  background: rgb(0 0 0 / .45);
  backdrop-filter: blur(3px);
}
.drawer {
  position: relative;
  width: min(86vw, 360px);
  height: 100%;
  background: rgb(var(--ww-parchment));
  border-right: 1px solid var(--ww-ink-hairline);
  padding: 20px 26px 28px;
  display: flex;
  flex-direction: column;
  box-shadow: 30px 0 60px -20px rgb(0 0 0 / .35);
  overflow-y: auto;
}

.head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 80, "opsz" 12, "wght" 500;
  font-size: 22px;
  letter-spacing: -0.01em;
}
.brand .dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 0 3px rgb(var(--ww-vermilion) / .15);
}
.close {
  width: 34px; height: 34px;
  font-size: 22px;
  background: transparent;
  border: 1px solid var(--ww-ink-hairline);
  color: rgb(var(--ww-ink));
  cursor: pointer;
  transition: background-color .25s ease, border-color .25s ease;
}
.close:hover { background: rgb(var(--ww-vermilion)); color: rgb(var(--ww-parchment)); border-color: rgb(var(--ww-vermilion)); }

.links {
  display: flex;
  flex-direction: column;
  margin-top: 24px;
  gap: 0;
}
.link {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 14px;
  align-items: baseline;
  padding: 16px 4px;
  border-bottom: 1px dashed var(--ww-ink-hairline);
  color: rgb(var(--ww-ink));
  opacity: 0;
  transform: translateX(-12px);
  animation: slide-in .55s cubic-bezier(.22,1,.36,1) forwards;
  transition: padding .35s cubic-bezier(.22,1,.36,1), color .25s ease;
}
.link:hover { padding-left: 12px; color: rgb(var(--ww-vermilion)); }
.link.router-link-active { color: rgb(var(--ww-vermilion)); }
.link.router-link-active::before {
  content: '';
  position: absolute;
  /* not used since position is static; kept simple */
}
.num {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .26em;
  color: var(--ww-ink-faint);
  text-transform: uppercase;
}
.label {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 24, "wght" 500;
  font-size: 22px;
}
.chevron {
  font-family: 'EB Garamond', serif;
  color: var(--ww-ink-faint);
  font-style: italic;
  transition: transform .35s cubic-bezier(.22,1,.36,1), color .25s ease;
}
.link:hover .chevron { transform: translateX(6px); color: rgb(var(--ww-vermilion)); }

.footer {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid var(--ww-ink-hairline);
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.row {
  display: flex;
  align-items: center;
  gap: 14px;
  justify-content: space-between;
}
.logout { align-self: flex-start; }

@keyframes slide-in {
  to { opacity: 1; transform: translateX(0); }
}

.drawer-enter-active, .drawer-leave-active { transition: opacity .35s ease; }
.drawer-enter-active .drawer, .drawer-leave-active .drawer {
  transition: transform .45s cubic-bezier(.22,1,.36,1);
}
.drawer-enter-from, .drawer-leave-to { opacity: 0; }
.drawer-enter-from .drawer, .drawer-leave-to .drawer { transform: translateX(-100%); }
</style>
