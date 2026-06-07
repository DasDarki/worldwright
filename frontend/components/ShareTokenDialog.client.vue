<script setup lang="ts">
interface ShareToken {
  token: string
  entity_id: number
  expires_at: string
  created_at: string
  created_by?: number | null
}

const props = defineProps<{
  open: boolean
  entityId: number
  entitySlug: string
}>()

const emit = defineEmits<{ close: [] }>()

const { t, locale } = useI18n()
const { $api } = useNuxtApp()
const toasts = useToastsStore()
const { absolute } = useSiteUrl()

const tokens = ref<ShareToken[]>([])
const loading = ref(false)
const creating = ref(false)
const ttlHours = ref<number>(24)
const ttlOptions = [1, 6, 24, 72, 168, 720] // 1h, 6h, 24h, 3d, 7d, 30d

async function loadTokens() {
  loading.value = true
  try {
    const res = await $api<{ share_tokens: ShareToken[] }>(`/entities/${props.entityId}/share-tokens`)
    tokens.value = res.share_tokens || []
  } catch (e: any) {
    toasts.error(e?.data?.error || t('entity.shareTokens.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function createToken() {
  creating.value = true
  try {
    const res = await $api<{ share_token: ShareToken }>(
      `/entities/${props.entityId}/share-tokens`,
      { method: 'POST', body: { ttl_seconds: ttlHours.value * 3600 } },
    )
    tokens.value = [res.share_token, ...tokens.value]
    await copyToClipboard(linkFor(res.share_token))
  } catch (e: any) {
    toasts.error(e?.data?.error || t('entity.shareTokens.createFailed'))
  } finally {
    creating.value = false
  }
}

async function revoke(token: string) {
  if (!confirm(t('entity.shareTokens.revokeConfirm'))) return
  try {
    await $api(`/share-tokens/${token}`, { method: 'DELETE' })
    tokens.value = tokens.value.filter((x) => x.token !== token)
    toasts.success(t('entity.shareTokens.revoked'))
  } catch (e: any) {
    toasts.error(e?.data?.error || t('entity.shareTokens.revokeFailed'))
  }
}

function linkFor(tok: ShareToken): string {
  return absolute(`/share/entity/${props.entitySlug}?token=${tok.token}`)
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    toasts.success(t('entity.shareTokens.copied'))
  } catch {
    // fallback: select via prompt
    window.prompt(t('entity.shareTokens.copyFallback'), text)
  }
}

function fmtRelative(iso: string): string {
  const t1 = new Date(iso).getTime()
  if (Number.isNaN(t1)) return iso
  const diff = t1 - Date.now()
  const sign = diff < 0 ? -1 : 1
  const abs = Math.abs(diff)
  const h = Math.floor(abs / 3_600_000)
  const m = Math.floor((abs % 3_600_000) / 60_000)
  if (sign < 0) return locale.value === 'de' ? 'abgelaufen' : 'expired'
  if (h > 24) return locale.value === 'de'
    ? `in ${Math.floor(h / 24)} Tagen`
    : `in ${Math.floor(h / 24)} days`
  if (h >= 1) return locale.value === 'de' ? `in ${h} h ${m} min` : `in ${h}h ${m}m`
  return locale.value === 'de' ? `in ${m} min` : `in ${m}m`
}

watch(() => props.open, (open) => { if (open) loadTokens() })
</script>

<template>
  <Teleport to="body">
    <Transition name="dlg-fade">
      <div v-if="open" class="ww-overlay" @click.self="emit('close')">
        <div class="dialog">
          <header class="head">
            <h3 class="title">{{ t('entity.shareTokens.title') }}</h3>
            <button type="button" class="x" @click="emit('close')" aria-label="Close">&times;</button>
          </header>
          <p class="hint">{{ t('entity.shareTokens.hint') }}</p>

          <div class="create">
            <label class="ttl">
              <span class="ww-label">{{ t('entity.shareTokens.duration') }}</span>
              <select v-model.number="ttlHours" class="ww-input">
                <option v-for="h in ttlOptions" :key="h" :value="h">
                  {{ h < 24 ? `${h} h` : `${h / 24} d` }}
                </option>
              </select>
            </label>
            <button
              type="button"
              class="ww-btn-primary"
              :disabled="creating"
              @click="createToken"
            >{{ creating ? t('common.loading') : t('entity.shareTokens.create') }}</button>
          </div>

          <div class="list-wrap">
            <p v-if="loading" class="state">{{ t('common.loading') }}</p>
            <ul v-else-if="tokens.length" class="list">
              <li v-for="tok in tokens" :key="tok.token" class="token-row">
                <div class="row-top">
                  <code class="link">{{ linkFor(tok) }}</code>
                  <div class="row-actions">
                    <button type="button" class="copy" @click="copyToClipboard(linkFor(tok))">
                      {{ t('entity.shareTokens.copy') }}
                    </button>
                    <button type="button" class="revoke" @click="revoke(tok.token)">
                      {{ t('entity.shareTokens.revoke') }}
                    </button>
                  </div>
                </div>
                <div class="meta">
                  {{ t('entity.shareTokens.expires') }} {{ fmtRelative(tok.expires_at) }}
                  <span class="abs">· {{ new Date(tok.expires_at).toLocaleString(locale === 'de' ? 'de-DE' : 'en-US') }}</span>
                </div>
              </li>
            </ul>
            <p v-else class="state empty">{{ t('entity.shareTokens.none') }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.ww-overlay {
  position: fixed; inset: 0; z-index: 9700;
  background: rgb(0 0 0 / .35);
  display: flex; align-items: center; justify-content: center;
  padding: 24px;
  backdrop-filter: blur(2px);
}
.dialog {
  background: var(--ww-card-bg);
  border: 1px solid var(--ww-ink-hairline);
  width: 640px; max-width: 100%;
  max-height: 80vh;
  display: flex; flex-direction: column;
  box-shadow: 0 50px 80px -40px rgb(0 0 0 / .5);
}
.head {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px 22px 6px;
}
.title {
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 60, "opsz" 32, "wght" 460;
  font-size: 22px;
  margin: 0;
}
.x { background: transparent; border: 0; font-size: 22px; color: var(--ww-ink-faint); cursor: pointer; }
.hint { padding: 0 22px 14px; font-style: italic; color: var(--ww-ink-faint); margin: 0; }

.create {
  display: flex; align-items: flex-end; gap: 14px;
  padding: 0 22px 16px;
  border-bottom: 1px solid var(--ww-ink-hairline);
}
.ttl { display: flex; flex-direction: column; gap: 4px; flex: 1; max-width: 160px; }

.list-wrap { overflow-y: auto; flex: 1; padding: 14px 22px; }
.state { padding: 24px 0; text-align: center; font-style: italic; color: var(--ww-ink-faint); }
.list { list-style: none; margin: 0; padding: 0; display: grid; gap: 14px; }
.token-row {
  border: 1px solid var(--ww-ink-hairline);
  padding: 12px 14px;
  background: rgb(var(--ww-parchment-stain) / .25);
}
.row-top {
  display: flex; align-items: center; gap: 12px;
  flex-wrap: wrap;
}
.link {
  flex: 1; min-width: 0;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  color: rgb(var(--ww-vermilion));
  word-break: break-all;
}
.row-actions { display: flex; gap: 8px; }
.copy, .revoke {
  font-family: 'Cormorant SC', serif;
  font-size: 11px;
  letter-spacing: .22em;
  text-transform: uppercase;
  background: transparent;
  border: 1px solid var(--ww-ink-hairline);
  padding: 4px 12px;
  cursor: pointer;
  transition: color .25s, border-color .25s;
}
.copy:hover { color: rgb(var(--ww-gold-deep)); border-color: rgb(var(--ww-gold)); }
.revoke:hover { color: rgb(var(--ww-vermilion)); border-color: rgb(var(--ww-vermilion)); }
.meta {
  margin-top: 8px;
  font-size: 13px;
  color: var(--ww-ink-faint);
  font-style: italic;
}
.abs { color: rgb(var(--ww-ink) / .35); }

.dlg-fade-enter-active, .dlg-fade-leave-active { transition: opacity .2s ease; }
.dlg-fade-enter-from, .dlg-fade-leave-to { opacity: 0; }
</style>
