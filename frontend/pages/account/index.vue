<script setup lang="ts">
import type { OAuthIdentity, User } from '~/types/api'
import { useAuthStore } from '~/stores/auth'
import { useToastsStore } from '~/stores/toasts'

const { t, locale, setLocale } = useI18n()
useHead({ title: () => t('account.title') })

const auth = useAuthStore()
const toasts = useToastsStore()
const { $api } = useNuxtApp()

const displayName = ref(auth.user?.display_name ?? '')
const avatarUrl = ref(auth.user?.avatar_url ?? '')
const userLocale = ref(auth.user?.locale ?? 'en')

const currentPassword = ref('')
const newPassword = ref('')
const newPasswordConfirm = ref('')

const profileSaving = ref(false)
const passwordSaving = ref(false)

const { data: identitiesData } = await useAsyncData('account-identities', () =>
  $api<{ identities: OAuthIdentity[] }>('/auth/me/identities'),
)
const identities = computed(() => identitiesData.value?.identities || [])

async function saveProfile() {
  profileSaving.value = true
  try {
    const { user } = await $api<{ user: User }>('/auth/me', {
      method: 'PATCH',
      body: {
        display_name: displayName.value,
        avatar_url: avatarUrl.value,
        locale: userLocale.value,
      },
    })
    auth.user = user
    if (userLocale.value !== locale.value) {
      await setLocale(userLocale.value as any)
    }
    toasts.success(t('account.profileSaved'))
  } catch (e: any) {
    toasts.error(e?.data?.error || t('account.saveFailed'))
  } finally {
    profileSaving.value = false
  }
}

async function changePassword() {
  if (newPassword.value.length < 6) {
    toasts.error(t('account.passwordTooShort'))
    return
  }
  if (newPassword.value !== newPasswordConfirm.value) {
    toasts.error(t('account.passwordMismatch'))
    return
  }
  passwordSaving.value = true
  try {
    await $api('/auth/me/password', {
      method: 'POST',
      body: {
        current_password: currentPassword.value,
        new_password: newPassword.value,
      },
    })
    currentPassword.value = ''
    newPassword.value = ''
    newPasswordConfirm.value = ''
    toasts.success(t('account.passwordChanged'))
  } catch (e: any) {
    const err = e?.data?.error
    if (err === 'invalid credentials') {
      toasts.error(t('account.currentPasswordWrong'))
    } else {
      toasts.error(err || t('account.saveFailed'))
    }
  } finally {
    passwordSaving.value = false
  }
}

function providerLabel(p: string): string {
  if (p === 'google') return 'Google'
  if (p === 'discord') return 'Discord'
  return p.charAt(0).toUpperCase() + p.slice(1)
}

const pruningSeed = ref(false)
async function pruneSeedData() {
  if (!confirm(t('settings.pruneSeed.confirm'))) return
  pruningSeed.value = true
  try {
    await $api('/admin/seed/prune', { method: 'POST' })
    auth.setOnboardingCompleted(true)
    toasts.success(t('settings.pruneSeed.done'))
  } catch (e: any) {
    toasts.error(e?.data?.error || t('settings.pruneSeed.failed'))
  } finally {
    pruningSeed.value = false
  }
}

useReveal()
</script>

<template>
  <section v-if="auth.user" class="py-12 md:py-20">
    <div class="mx-auto max-w-screen-xl">
      <div class="stagger mb-16">
        <div class="ww-eyebrow mb-6 flex items-center gap-3">
          <span class="inline-block w-9 h-px bg-current" aria-hidden="true" />
          Vol. I · {{ t('account.eyebrow') }}
        </div>
        <h1 class="hero-title"><em>{{ t('account.title') }}</em></h1>
        <p class="lede">{{ auth.user.email }}</p>
      </div>

      <div class="grid">
        <section class="block">
          <h3 class="ww-label section-head">{{ t('account.profile') }}</h3>
          <label class="field">
            <span class="lbl">{{ t('account.displayName') }}</span>
            <input v-model="displayName" class="ww-input" :placeholder="t('account.displayNamePlaceholder')" />
          </label>
          <label class="field">
            <span class="lbl">{{ t('account.avatarUrl') }}</span>
            <input v-model="avatarUrl" class="ww-input" placeholder="https://…" />
          </label>
          <label class="field">
            <span class="lbl">{{ t('account.locale') }}</span>
            <select v-model="userLocale" class="ww-input">
              <option value="en">English</option>
              <option value="de">Deutsch</option>
            </select>
          </label>
          <div class="actions">
            <button type="button" class="ww-btn-primary" :disabled="profileSaving" @click="saveProfile">
              {{ profileSaving ? t('common.loading') : t('account.saveProfile') }}
            </button>
          </div>
        </section>

        <section class="block">
          <h3 class="ww-label section-head">{{ t('account.password') }}</h3>
          <label v-if="auth.user.has_password" class="field">
            <span class="lbl">{{ t('account.currentPassword') }}</span>
            <input v-model="currentPassword" type="password" class="ww-input" autocomplete="current-password" />
          </label>
          <p v-else class="hint">{{ t('account.noPasswordYet') }}</p>
          <label class="field">
            <span class="lbl">{{ t('account.newPassword') }}</span>
            <input v-model="newPassword" type="password" class="ww-input" autocomplete="new-password" />
          </label>
          <label class="field">
            <span class="lbl">{{ t('account.newPasswordConfirm') }}</span>
            <input v-model="newPasswordConfirm" type="password" class="ww-input" autocomplete="new-password" />
          </label>
          <div class="actions">
            <button type="button" class="ww-btn-primary" :disabled="passwordSaving || !newPassword" @click="changePassword">
              {{ passwordSaving ? t('common.loading') : t('account.savePassword') }}
            </button>
          </div>
        </section>

        <section class="block">
          <h3 class="ww-label section-head">{{ t('account.identities') }}</h3>
          <ul v-if="identities.length" class="identities">
            <li v-for="i in identities" :key="i.id" class="identity">
              <div class="prov">
                <span class="prov-dot" :data-prov="i.provider" aria-hidden="true" />
                <span class="prov-name">{{ providerLabel(i.provider) }}</span>
              </div>
              <span v-if="i.email" class="prov-email">{{ i.email }}</span>
              <span class="prov-sub">{{ i.subject.slice(0, 8) }}…</span>
            </li>
          </ul>
          <p v-else class="hint">{{ t('account.noIdentities') }}</p>
        </section>

        <section v-if="auth.isAdmin && auth.onboarding?.seed_data_present" class="block">
          <h3 class="ww-label section-head">{{ t('settings.pruneSeed.label') }}</h3>
          <p class="hint">{{ t('settings.pruneSeed.body') }}</p>
          <button
            type="button"
            class="ww-btn-ghost prune-btn"
            :disabled="pruningSeed"
            @click="pruneSeedData"
          >
            {{ pruningSeed ? t('common.loading') : t('settings.pruneSeed.cta') }}
          </button>
        </section>
      </div>
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
.lede { font-style: italic; color: var(--ww-ink-faint); font-size: 17px; }

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
  gap: 30px 40px;
}

.block {
  display: grid;
  gap: 16px;
  background: rgb(var(--ww-parchment-deep) / .25);
  border: 1px solid var(--ww-ink-hairline);
  padding: 24px 26px;
}
.section-head {
  margin: 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.field { display: grid; gap: 4px; }
.lbl {
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  letter-spacing: .28em;
  text-transform: uppercase;
  color: var(--ww-ink-faint);
}
.hint {
  font-style: italic;
  color: var(--ww-ink-faint);
  font-size: 14px;
  margin: 0;
}
.actions {
  margin-top: 6px;
  display: flex;
  justify-content: flex-end;
}

.identities { list-style: none; margin: 0; padding: 0; display: grid; gap: 10px; }
.identity {
  display: grid;
  grid-template-columns: max-content 1fr max-content;
  gap: 14px;
  align-items: baseline;
  padding: 10px 12px;
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-parchment) / .6);
}
.prov { display: inline-flex; align-items: center; gap: 8px; }
.prov-dot {
  width: 10px; height: 10px;
  border-radius: 50%;
  background: rgb(var(--ww-ink-shade) / .4);
}
.prov-dot[data-prov="google"] { background: #4285F4; }
.prov-dot[data-prov="discord"] { background: #5865F2; }
.prov-name {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
}
.prov-email {
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  color: rgb(var(--ww-ink-shade));
}
.prov-sub {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 10px;
  color: var(--ww-ink-faint);
}

.prune-btn {
  color: rgb(var(--ww-vermilion));
  border-bottom-color: rgb(var(--ww-vermilion) / .35);
  margin-top: 8px;
}
.prune-btn:hover {
  color: rgb(var(--ww-vermilion-deep));
  border-bottom-color: rgb(var(--ww-vermilion));
}
.prune-btn:disabled { opacity: .5; cursor: progress; }

select.ww-input {
  appearance: none;
  background-color: transparent;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23173842'/></svg>");
  background-repeat: no-repeat;
  background-position: right .3em center;
  background-size: 8px 5px;
  padding-right: 1.5em;
}
:root.dark select.ww-input {
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 8 5'><path d='M0 0 L4 5 L8 0' fill='%23ecdfc2'/></svg>");
}
</style>
