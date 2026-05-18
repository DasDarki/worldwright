import { defineStore } from 'pinia'

export type ToastKind = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  kind: ToastKind
  message: string
  ttl: number
}

export const useToastsStore = defineStore('toasts', () => {
  const items = ref<Toast[]>([])
  let nextId = 1

  function push(kind: ToastKind, message: string, ttl = 3200) {
    if (import.meta.server) return
    const id = nextId++
    items.value = [...items.value, { id, kind, message, ttl }]
    if (ttl > 0) {
      setTimeout(() => dismiss(id), ttl)
    }
    return id
  }

  function dismiss(id: number) {
    items.value = items.value.filter((t) => t.id !== id)
  }

  function success(message: string, ttl?: number) { return push('success', message, ttl) }
  function error(message: string, ttl?: number)   { return push('error',   message, ttl ?? 5000) }
  function info(message: string, ttl?: number)    { return push('info',    message, ttl) }

  return { items, push, dismiss, success, error, info }
})
