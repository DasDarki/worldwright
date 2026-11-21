import { useAuthStore } from '~/stores/auth'
import type { User } from '~/types/api'

export type OnyxResult = 'linked' | 'taken' | 'other' | 'failed' | 'cancelled' | 'expired'

const claimStorageKey = 'worldwright:onyx-claim'
const results: OnyxResult[] = ['linked', 'taken', 'other', 'failed', 'cancelled', 'expired']

export function isOnyxResult(value: unknown): value is OnyxResult {
  return typeof value === 'string' && (results as string[]).includes(value)
}

function readClaim(): string | null {
  if (!import.meta.client) return null
  try {
    return sessionStorage.getItem(claimStorageKey)
  } catch {
    return null
  }
}

function writeClaim(value: string | null) {
  if (!import.meta.client) return
  try {
    if (value) sessionStorage.setItem(claimStorageKey, value)
    else sessionStorage.removeItem(claimStorageKey)
  } catch {
    return
  }
}

function safePath(value: unknown): string {
  return typeof value === 'string' && value.startsWith('/') && !value.startsWith('//') && !value.startsWith('/\\')
    ? value
    : ''
}

export function useOnyx() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()
  const { $api } = useNuxtApp()
  const result = useState<OnyxResult | null>('onyx-result', () => null)

  const enabled = computed(() => !!auth.user?.onyx_enabled)
  const linked = computed(() => !!auth.user?.onyx_linked)
  const needsMigration = computed(() => auth.isAuthenticated && enabled.value && !linked.value)

  function loginUrl(redirect?: unknown) {
    const path = safePath(redirect)
    const query = path ? `?redirect=${encodeURIComponent(path)}` : ''
    return `${config.public.apiBase}/auth/onyx/start${query}`
  }

  function registerUrl(claim: string) {
    return `${config.public.apiBase}/auth/onyx/register?claim=${encodeURIComponent(claim)}`
  }

  async function startMigration() {
    const response = await $api<{ url: string }>('/auth/onyx/link', { method: 'POST' })
    if (!response?.url) throw new Error('onyx link failed')
    window.location.href = response.url
  }

  function rememberClaim(claim: string) {
    writeClaim(claim)
  }

  function forgetClaim() {
    writeClaim(null)
  }

  function hasPendingClaim() {
    return readClaim() !== null
  }

  async function completePendingClaim() {
    const claim = readClaim()
    if (!claim || !auth.isAuthenticated) return
    forgetClaim()
    try {
      const response = await $api<{ result: string, user?: User }>('/auth/onyx/claim', { method: 'POST', body: { claim } })
      result.value = isOnyxResult(response.result) ? response.result : 'failed'
      if (response.user) auth.user = response.user
    } catch {
      result.value = 'failed'
    }
  }

  return {
    result,
    enabled,
    linked,
    needsMigration,
    loginUrl,
    registerUrl,
    startMigration,
    rememberClaim,
    forgetClaim,
    hasPendingClaim,
    completePendingClaim,
  }
}
