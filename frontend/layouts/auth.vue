<template>
  <div class="min-h-screen flex flex-col">
    <div class="auth-bar mx-auto w-full max-w-screen-2xl px-6 md:px-12 py-5 flex items-center justify-between">
      <NuxtLink to="/login" class="brand">
        <span class="dot" aria-hidden="true" />
        worldwright
        <span class="vol">vol. I</span>
      </NuxtLink>
      <div class="flex items-center gap-3">
        <LangToggle />
        <ThemeToggle />
      </div>
    </div>
    <main class="flex-1">
      <slot />
    </main>
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
</style>
