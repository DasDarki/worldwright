export function useReveal() {
  if (import.meta.server) return

  onMounted(() => {
    const observer = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add('in')
          observer.unobserve(entry.target)
        }
      }
    }, { rootMargin: '0px 0px -8% 0px', threshold: 0.05 })

    document.querySelectorAll('.reveal:not(.in)').forEach((el) => observer.observe(el))

    onBeforeUnmount(() => observer.disconnect())
  })
}
