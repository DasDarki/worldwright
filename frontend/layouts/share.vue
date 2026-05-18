<template>
  <div class="min-h-screen flex flex-col share-layout">
    <div class="share-bar mx-auto w-full max-w-screen-2xl px-6 md:px-12 py-5 flex items-center justify-between">
      <a href="/" class="brand">
        <span class="dot" aria-hidden="true" />
        worldwright
        <span class="vol">vol. I</span>
      </a>
      <div class="ww-label hint">Public share</div>
    </div>
    <main class="flex-1">
      <slot />
    </main>
    <footer class="share-footer mx-auto w-full max-w-screen-2xl px-6 md:px-12 py-10">
      <div class="ww-rule mb-6" />
      <p class="footer-note">
        This is a public excerpt from a private codex.
        <a href="/" class="ww-link">Pressed in Worldwright.</a>
      </p>
    </footer>
  </div>
</template>

<style scoped lang="scss">
.brand {
  display: flex;
  align-items: baseline;
  gap: 12px;
  font-family: 'Fraunces', serif;
  font-variation-settings: "SOFT" 80, "opsz" 12, "wght" 500;
  font-size: 22px;
  letter-spacing: -0.01em;
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
}
.hint { color: rgb(var(--ww-gold-deep)); }
.footer-note {
  font-family: 'EB Garamond', serif;
  font-style: italic;
  color: var(--ww-ink-faint);
  text-align: center;
  font-size: 14px;
}
</style>
