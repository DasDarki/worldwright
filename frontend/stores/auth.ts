import { defineStore } from 'pinia'
import type { AuthProviders, User } from '~/types/api'

export interface OnboardingStatus {
  completed: boolean
  seed_data_present: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const providers = ref<AuthProviders>({ password: true, oauth: [] })
  const fetched = ref(false)
  const onboarding = ref<OnboardingStatus | null>(null)

  const isAuthenticated = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function fetchMe() {
    fetched.value = true
    try {
      const { user: u } = await useNuxtApp().$api<{ user: User }>('/auth/me')
      user.value = u
    } catch {
      user.value = null
    }
  }

  async function fetchProviders() {
    try {
      providers.value = await useNuxtApp().$api<AuthProviders>('/auth/providers')
    } catch {
      providers.value = { password: true, oauth: [] }
    }
  }

  async function fetchOnboarding() {
    if (!isAdmin.value) {
      onboarding.value = null
      return
    }
    try {
      const { onboarding: status } = await useNuxtApp().$api<{ onboarding: OnboardingStatus }>('/admin/onboarding')
      onboarding.value = status
    } catch {
      onboarding.value = null
    }
  }

  function setOnboardingCompleted(seedRemoved = false) {
    if (onboarding.value) {
      onboarding.value.completed = true
      if (seedRemoved) onboarding.value.seed_data_present = false
    } else {
      onboarding.value = { completed: true, seed_data_present: !seedRemoved }
    }
  }

  async function login(email: string, password: string) {
    const { user: u } = await useNuxtApp().$api<{ user: User }>('/auth/login', {
      method: 'POST',
      body: { email, password },
    })
    user.value = u
  }

  async function logout() {
    try {
      await useNuxtApp().$api('/auth/logout', { method: 'POST' })
    } finally {
      user.value = null
      onboarding.value = null
    }
  }

  function clear() {
    user.value = null
    fetched.value = false
    onboarding.value = null
  }

  return {
    user, providers, fetched, onboarding,
    isAuthenticated, isAdmin,
    fetchMe, fetchProviders, fetchOnboarding,
    setOnboardingCompleted,
    login, logout, clear,
  }
})
