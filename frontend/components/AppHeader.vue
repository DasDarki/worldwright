<script setup lang="ts">
import type { Calendar } from '~/types/api'
import { useAuthStore } from '~/stores/auth'
import { useCalendars } from '~/composables/useCalendar'

const auth = useAuthStore()
const { t } = useI18n()
const router = useRouter()

const headerCalendar = ref<Calendar | null>(null)
if (auth.isAuthenticated) {
  try {
    const { load } = useCalendars()
    const list = await load()
    const first = list[0]
    if (first?.id) {
      const { $api } = useNuxtApp()
      const res = await $api<{ calendar: Calendar }>(`/calendars/${first.id}`)
      headerCalendar.value = res.calendar
    }
  } catch {
    headerCalendar.value = null
  }
}

const scrolled = ref(false)
if (import.meta.client) {
  onMounted(() => {
    const onScroll = () => { scrolled.value = window.scrollY > 12 }
    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))
  })
}

async function logout() {
  await auth.logout()
  router.push('/login')
}

const drawerOpen = ref(false)
</script>

<template>
  <header :class="['ww-nav', { scrolled }]">
    <div class="mx-auto max-w-screen-2xl px-6 md:px-12 flex items-center gap-6 py-4">
      <button
        v-if="auth.isAuthenticated"
        type="button"
        class="hamburger"
        :aria-label="t('mobile.menu')"
        :aria-expanded="drawerOpen"
        @click="drawerOpen = true"
      >
        <span class="bar" />
        <span class="bar" />
        <span class="bar" />
      </button>

      <NuxtLink to="/" class="brand">
        <span class="dot" aria-hidden="true" />
        worldwright
        <span class="vol">vol. I</span>
      </NuxtLink>

      <nav class="menu" :aria-label="t('nav.codex')">
        <NuxtLink to="/">{{ t('nav.codex') }}</NuxtLink>
        <NuxtLink to="/entities">{{ t('nav.entities') }}</NuxtLink>
        <NuxtLink to="/tags">{{ t('nav.tags') }}</NuxtLink>
        <NuxtLink to="/events">{{ t('nav.events') }}</NuxtLink>
        <NuxtLink to="/timelines">{{ t('nav.timelines') }}</NuxtLink>
        <NuxtLink to="/maps">{{ t('nav.maps') }}</NuxtLink>
        <NuxtLink to="/calendars">{{ t('nav.calendars') }}</NuxtLink>
        <NuxtLink v-if="auth.isAdmin" to="/assets">{{ t('nav.assets') }}</NuxtLink>
      </nav>

      <div class="spacer" />

      <SearchBox v-if="auth.isAuthenticated" />

      <MoonDisplay
        v-if="auth.isAuthenticated && headerCalendar"
        :calendar="headerCalendar"
        variant="minimal"
        :size="22"
        class="header-moons"
        :title="t('moons.today')"
      />

      <div class="actions">
        <NuxtLink
          v-if="auth.isAuthenticated"
          to="/account"
          class="me"
          :title="auth.user?.email"
        >
          {{ auth.user?.display_name || (auth.user?.email || '').split('@')[0] }}
        </NuxtLink>
        <LangToggle />
        <ThemeToggle />
        <button
          v-if="auth.isAuthenticated"
          type="button"
          class="logout"
          @click="logout"
        >{{ t('nav.signOut') }}</button>
        <NuxtLink
          v-else
          to="/login"
          class="logout"
        >{{ t('nav.signIn') }}</NuxtLink>
      </div>
    </div>
  </header>
  <MobileDrawer :open="drawerOpen" @close="drawerOpen = false" />
</template>

<style scoped lang="scss">
.ww-nav {
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(6px);
  background: linear-gradient(to bottom, rgb(var(--ww-parchment) / .94) 60%, transparent);
  border-bottom: 1px solid transparent;
  transition: background-color .4s ease, border-color .4s ease;
  &.scrolled {
    background: rgb(var(--ww-parchment) / .96);
    border-bottom-color: var(--ww-ink-hairline);
  }
}
.spacer { flex: 1; }
.brand {
  display: flex;
  align-items: baseline;
  gap: 12px;
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 80, "opsz" 12, "wght" 500;
  font-size: 22px;
  letter-spacing: -0.01em;
  flex-shrink: 0;
}
.brand .dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: rgb(var(--ww-vermilion));
  box-shadow: 0 0 0 3px rgb(var(--ww-vermilion) / .15);
  align-self: center;
}
.brand .vol {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .3em;
  color: var(--ww-ink-faint);
  margin-left: 6px;
}
.menu {
  display: flex; gap: 28px;
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  @media (max-width: 1100px) { display: none; }
}

.hamburger {
  display: none;
  width: 36px; height: 36px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 5px;
  border: 1px solid var(--ww-ink-hairline);
  background: transparent;
  cursor: pointer;
  transition: border-color .25s ease, background-color .25s ease;
  @media (max-width: 1100px) { display: inline-flex; }
}
.hamburger:hover {
  border-color: rgb(var(--ww-gold));
  background: rgb(var(--ww-gold) / .08);
}
.hamburger .bar {
  display: block;
  width: 16px;
  height: 1.5px;
  background: rgb(var(--ww-ink));
  transition: transform .35s cubic-bezier(.22,1,.36,1), opacity .25s ease;
}
.hamburger:hover .bar:nth-child(1) { transform: translateX(-2px); }
.hamburger:hover .bar:nth-child(3) { transform: translateX(2px); }
.menu a { position: relative; padding: 6px 0; }
.menu a::after {
  content: '';
  position: absolute; left: 0; right: 0; bottom: -2px;
  height: 1px;
  background: rgb(var(--ww-gold));
  transform: scaleX(0);
  transform-origin: left;
  transition: transform .45s cubic-bezier(.22,1,.36,1);
}
.menu a:hover::after, .menu a.router-link-active::after { transform: scaleX(1); }

.header-moons {
  padding: 0 4px;
  @media (max-width: 880px) { display: none; }
}

.actions {
  display: flex; align-items: center; gap: 14px;
}
.logout {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  padding: 9px 16px;
  border: 1px solid var(--ww-ink-hairline);
  color: rgb(var(--ww-ink));
  transition: background-color .35s ease, color .35s ease, border-color .35s ease;
}
.logout:hover {
  background: rgb(var(--ww-ink));
  color: rgb(var(--ww-parchment));
}
.me {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  font-size: 14px;
  color: rgb(var(--ww-ink));
  border-bottom: 1px solid transparent;
  padding: 4px 0;
  transition: color .25s ease, border-color .25s ease;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  @media (max-width: 880px) { display: none; }
}
.me:hover { color: rgb(var(--ww-vermilion)); border-bottom-color: rgb(var(--ww-vermilion)); }
</style>
