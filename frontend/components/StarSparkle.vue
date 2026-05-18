<script setup lang="ts">
withDefaults(defineProps<{ size?: number; delay?: number }>(), { size: 18, delay: 0 })
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="currentColor"
    aria-hidden="true"
    class="text-gold"
    :style="{ animation: `twinkle 4s ease-in-out ${delay}s infinite` }"
  >
    <path d="M12 0 L13.5 10.5 L24 12 L13.5 13.5 L12 24 L10.5 13.5 L0 12 L10.5 10.5 Z" />
  </svg>
</template>
