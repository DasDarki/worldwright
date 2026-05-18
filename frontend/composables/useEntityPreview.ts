import type { Entity } from '~/types/api'

type CacheEntry =
  | { status: 'loading'; promise: Promise<Entity | null> }
  | { status: 'ok'; entity: Entity }
  | { status: 'missing' }
  | { status: 'error' }

const cache = new Map<string, CacheEntry>()

export function useEntityPreview() {
  const { $api } = useNuxtApp()

  async function fetchPreview(slug: string): Promise<Entity | null> {
    if (!slug) return null
    const hit = cache.get(slug)
    if (hit) {
      if (hit.status === 'ok') return hit.entity
      if (hit.status === 'missing' || hit.status === 'error') return null
      return hit.promise
    }
    const promise = (async () => {
      try {
        const { entity } = await $api<{ entity: Entity }>(`/entities/by-slug/${encodeURIComponent(slug)}`)
        cache.set(slug, { status: 'ok', entity })
        return entity
      } catch (e: any) {
        const code = e?.statusCode || e?.response?.status
        cache.set(slug, { status: code === 404 ? 'missing' : 'error' })
        return null
      }
    })()
    cache.set(slug, { status: 'loading', promise })
    return promise
  }

  function clear() {
    cache.clear()
  }

  return { fetchPreview, clear }
}
