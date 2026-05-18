<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

definePageMeta({ layout: 'auth' })

const { t } = useI18n()
useHead({ title: () => t('auth.title') })
const { absolute } = useSiteUrl()
const ogImageUrl = computed(() => absolute('/banner.png'))

useSeoMeta({
  title: () => t('auth.title'),
  description: () => t('app.tagline'),
  ogTitle: 'Worldwright',
  ogDescription: () => t('app.tagline'),
  ogImage: () => ogImageUrl.value,
  ogImageWidth: 1200,
  ogImageHeight: 630,
  twitterCard: 'summary_large_image',
  twitterImage: () => ogImageUrl.value,
})

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

if (!auth.providers.oauth.length) await auth.fetchProviders()
if (auth.isAuthenticated) await router.replace((route.query.redirect as string) || '/')

const email = ref('')
const password = ref('')
const error = ref<string | null>(null)
const pending = ref(false)

async function onSubmit() {
  error.value = null
  pending.value = true
  try {
    await auth.login(email.value, password.value)
    await router.replace((route.query.redirect as string) || '/')
  } catch (e: any) {
    error.value = e?.data?.error === 'invalid credentials'
      ? t('auth.errors.invalid')
      : t('auth.errors.generic')
  } finally {
    pending.value = false
  }
}

const config = useRuntimeConfig()
const apiBase = computed(() => config.public.apiBase as string)
</script>

<template>
  <section class="login-section">
    <div class="mx-auto max-w-screen-xl px-6 md:px-12 py-16 md:py-24 grid md:grid-cols-2 gap-12 md:gap-20 items-center">
      <div class="stagger">
        <div class="ww-eyebrow flex items-center gap-3 mb-6">
          <span class="eyebrow-rule" aria-hidden="true" />
          {{ t('auth.eyebrow') }}
        </div>
        <h1 class="hero-title mb-8">
          {{ t('auth.title') }}<br />
          <em>{{ t('app.tagline') }}</em>
        </h1>
        <p class="lede mb-10">{{ t('auth.subtitle') }}</p>

        <form class="form" @submit.prevent="onSubmit" novalidate>
          <label class="field">
            <span class="ww-label mb-1 block">{{ t('auth.email') }}</span>
            <input
              v-model="email"
              type="email"
              autocomplete="email"
              required
              class="ww-input"
              placeholder="admin@worldwright.local"
            />
          </label>
          <label class="field">
            <span class="ww-label mb-1 block">{{ t('auth.password') }}</span>
            <input
              v-model="password"
              type="password"
              autocomplete="current-password"
              required
              class="ww-input"
              placeholder="••••••••"
            />
          </label>

          <Transition name="fade">
            <p v-if="error" class="error">{{ error }}</p>
          </Transition>

          <button class="ww-btn-primary submit" type="submit" :disabled="pending">
            {{ pending ? t('common.loading') : t('auth.submit') }}
            <span class="arrow" aria-hidden="true">
              <svg width="18" height="10" viewBox="0 0 18 10" fill="none">
                <path d="M1 5 H16 M11 1 L16 5 L11 9" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </span>
          </button>
        </form>

        <div v-if="auth.providers.oauth.length" class="oauth-block">
          <div class="oauth-divider">
            <span>{{ t('auth.or') }}</span>
          </div>
          <div class="oauth-buttons">
            <a
              v-if="auth.providers.oauth.includes('google')"
              :href="`${apiBase}/auth/oauth/google/start`"
              class="oauth-btn"
            >
              <svg viewBox="0 0 18 18" width="18" height="18" aria-hidden="true">
                <path fill="#4285F4" d="M17.64 9.2c0-.64-.06-1.25-.17-1.84H9v3.48h4.84a4.14 4.14 0 0 1-1.79 2.72v2.26h2.9c1.7-1.57 2.69-3.88 2.69-6.62z"/>
                <path fill="#34A853" d="M9 18c2.43 0 4.47-.8 5.96-2.18l-2.9-2.26c-.8.54-1.84.86-3.06.86-2.35 0-4.34-1.59-5.05-3.72H.96v2.33A9 9 0 0 0 9 18z"/>
                <path fill="#FBBC05" d="M3.95 10.7a5.41 5.41 0 0 1 0-3.4V4.96H.96a9 9 0 0 0 0 8.08l3-2.34z"/>
                <path fill="#EA4335" d="M9 3.58c1.32 0 2.5.46 3.44 1.35l2.58-2.58A9 9 0 0 0 9 0a9 9 0 0 0-8.04 4.96l2.99 2.34C4.66 5.17 6.65 3.58 9 3.58z"/>
              </svg>
              {{ t('auth.withGoogle') }}
            </a>
            <a
              v-if="auth.providers.oauth.includes('discord')"
              :href="`${apiBase}/auth/oauth/discord/start`"
              class="oauth-btn discord"
            >
              <svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true" fill="#5865F2">
                <path d="M19.27 5.33A19.94 19.94 0 0 0 14.4 4l-.24.44a18.69 18.69 0 0 1 4.32 1.36c-2-.92-4.18-1.36-6.48-1.36-2.3 0-4.48.44-6.48 1.36a18.69 18.69 0 0 1 4.32-1.36L9.6 4a19.94 19.94 0 0 0-4.87 1.33C1.7 9.83.93 14.23 1.3 18.55a20.6 20.6 0 0 0 6.24 3.13l1.26-1.7a13.07 13.07 0 0 1-2.07-.99l.5-.36a14.27 14.27 0 0 0 12.55 0l.5.36c-.65.38-1.34.7-2.07.99l1.26 1.7a20.55 20.55 0 0 0 6.24-3.13c.45-5.01-.74-9.36-3.43-13.22zM9.45 15.7c-1.18 0-2.16-1.08-2.16-2.41 0-1.33.96-2.42 2.16-2.42 1.21 0 2.19 1.1 2.16 2.42 0 1.33-.95 2.41-2.16 2.41zm5.1 0c-1.18 0-2.16-1.08-2.16-2.41 0-1.33.96-2.42 2.16-2.42 1.21 0 2.18 1.1 2.16 2.42 0 1.33-.95 2.41-2.16 2.41z"/>
              </svg>
              {{ t('auth.withDiscord') }}
            </a>
          </div>
        </div>

        <p class="footer-note">{{ t('auth.footer') }}</p>
      </div>

      <div class="medallion-wrap">
        <div class="medallion" aria-hidden="true">
          <svg class="rings" viewBox="0 0 600 600" fill="none">
            <circle cx="300" cy="300" r="290" stroke="rgb(124 94 48 / .5)" stroke-width="1" stroke-dasharray="2 6"/>
            <circle cx="300" cy="300" r="260" stroke="rgb(23 56 66 / .25)" stroke-width="1"/>
            <g stroke="rgb(23 56 66 / .55)" stroke-width="1">
              <line x1="300" y1="6" x2="300" y2="22"/>
              <line x1="300" y1="578" x2="300" y2="594"/>
              <line x1="6" y1="300" x2="22" y2="300"/>
              <line x1="578" y1="300" x2="594" y2="300"/>
            </g>
            <g font-family="Cormorant SC, serif" font-size="12" fill="rgb(124 94 48)" letter-spacing="3">
              <text x="300" y="40" text-anchor="middle">N</text>
              <text x="300" y="572" text-anchor="middle">S</text>
              <text x="40" y="306" text-anchor="middle">W</text>
              <text x="560" y="306" text-anchor="middle">E</text>
            </g>
            <g fill="rgb(184 147 90)">
              <circle cx="300" cy="14" r="2.5"/>
              <circle cx="586" cy="300" r="2.5"/>
              <circle cx="300" cy="586" r="2.5"/>
              <circle cx="14" cy="300" r="2.5"/>
            </g>
          </svg>
          <div class="plate">
            <img src="/banner.png" alt="Worldwright" />
          </div>
          <StarSparkle class="star s1" :size="22" :delay="0" />
          <StarSparkle class="star s2" :size="16" :delay="1.2" />
          <StarSparkle class="star s3" :size="18" :delay="2.4" />
          <StarSparkle class="star s4" :size="14" :delay="0.6" />
        </div>
        <p class="medallion-caption">{{ t('app.tagline') }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.login-section { position: relative; }

.eyebrow-rule {
  display: inline-block;
  width: 36px;
  height: 1px;
  background: currentColor;
}

.hero-title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 70, "opsz" 144, "wght" 380;
  font-size: clamp(40px, 5.6vw, 76px);
  line-height: 0.95;
  letter-spacing: -0.025em;
  margin: 0;
}
.hero-title em {
  font-style: italic;
  font-variation-settings: "SOFT" 100, "opsz" 144, "wght" 320;
  color: rgb(var(--ww-gold-deep));
}
.lede {
  font-size: 18px;
  line-height: 1.55;
  color: var(--ww-ink-faint);
  max-width: 36em;
}

.form { display: grid; gap: 22px; max-width: 30em; }
.field { display: block; }
.submit { margin-top: 6px; }
.submit .arrow { transition: transform .5s cubic-bezier(.22,1,.36,1); }
.submit:hover .arrow { transform: translateX(6px); }
.submit:disabled { opacity: .6; cursor: progress; }

.error {
  color: rgb(var(--ww-vermilion));
  font-style: italic;
  font-size: 14px;
  margin-top: 4px;
}

.fade-enter-active, .fade-leave-active { transition: opacity .35s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.oauth-block { margin-top: 32px; max-width: 30em; }
.oauth-divider {
  display: flex; align-items: center; gap: 16px;
  font-family: 'Cormorant SC', serif;
  font-size: 11px; letter-spacing: .32em; text-transform: uppercase;
  color: var(--ww-ink-faint);
  margin-bottom: 18px;
  &::before, &::after {
    content: ''; height: 1px; flex: 1;
    background: linear-gradient(to right, transparent, var(--ww-ink-hairline) 50%, transparent);
  }
}
.oauth-buttons { display: grid; gap: 10px; }
.oauth-btn {
  display: flex; align-items: center; justify-content: center; gap: 12px;
  padding: 14px 18px;
  border: 1px solid var(--ww-ink-hairline);
  font-family: 'Cormorant SC', serif;
  font-size: 12px;
  letter-spacing: .22em;
  text-transform: uppercase;
  transition: border-color .35s ease, background-color .35s ease, transform .35s ease;
}
.oauth-btn:hover {
  border-color: rgb(var(--ww-gold));
  background: rgb(var(--ww-gold) / .08);
  transform: translateY(-1px);
}

.footer-note {
  font-style: italic;
  color: var(--ww-ink-faint);
  margin-top: 36px;
  font-size: 14px;
  max-width: 30em;
}

/* Medallion */
.medallion-wrap { display: flex; flex-direction: column; align-items: center; }
.medallion {
  position: relative;
  width: clamp(260px, 38vw, 460px);
  aspect-ratio: 1;
  isolation: isolate;
  margin-bottom: 28px;
}
.medallion .rings {
  position: absolute;
  inset: -4%;
  animation: slowSpin 90s linear infinite;
  opacity: 0;
  animation: dawn 1.6s cubic-bezier(.22,1,.36,1) .9s forwards, slowSpin 90s linear infinite;
}
.medallion .plate {
  position: absolute;
  inset: 10%;
  border-radius: 50%;
  background: radial-gradient(circle at 32% 28%, var(--ww-plate-from) 0%, var(--ww-plate-mid) 55%, var(--ww-plate-to) 100%);
  box-shadow:
    inset 0 0 60px rgb(124 94 48 / .25),
    inset 0 0 0 1px rgb(124 94 48 / .35),
    var(--ww-plate-glow);
  overflow: hidden;
  opacity: 0;
  transform: scale(.94);
  animation: dawn 1.6s cubic-bezier(.22,1,.36,1) .4s forwards;
}
.medallion .plate img {
  position: absolute;
  inset: 8% 6% 10%;
  width: auto;
  max-width: 86%;
  max-height: 84%;
  margin: auto;
  filter: contrast(1.04) saturate(.95);
}
.medallion .star { position: absolute; }
.medallion .star.s1 { top: -2%;  left: 18%; }
.medallion .star.s2 { top: 30%;  right: -4%; }
.medallion .star.s3 { bottom: 4%; left: -3%; }
.medallion .star.s4 { bottom: -2%; right: 22%; }

.medallion-caption {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .42em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
  text-align: center;
  max-width: 28em;
}

@media (max-width: 760px) {
  .medallion-wrap { display: none; }
}

@keyframes slowSpin { to { transform: rotate(360deg); } }
@keyframes dawn { to { opacity: 1; transform: scale(1) rotate(0deg); } }
</style>
