<script setup lang="ts">
const route = useRoute()
const showSidebar = computed(() => route.path.startsWith('/entities'))
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <AppHeader />
    <div class="layout-body mx-auto w-full max-w-screen-2xl px-6 md:px-12 pt-10 pb-24 flex-1">
      <AppSidebar v-if="showSidebar" />
      <main class="content">
        <slot />
      </main>
    </div>
    <footer class="mx-auto w-full max-w-screen-2xl px-6 md:px-12 py-10">
      <div class="ww-rule mb-6" />
      <div class="flex items-center justify-between text-xs ww-label">
        <span>© MMXXVI · Sub Stellis</span>
        <span>EN · DE</span>
      </div>
    </footer>
  </div>
</template>

<style scoped lang="scss">
.layout-body {
  display: flex;
  gap: clamp(20px, 3vw, 44px);
}
.content {
  flex: 1;
  min-width: 0;
}
</style>
