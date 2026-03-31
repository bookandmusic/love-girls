<script setup lang="ts">
defineProps<{
  width?: string;
  height?: string;
  rounded?: string;
  className?: string;
}>();
</script>

<template>
  <div
    :class="[
      'skeleton-loader bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 animate-pulse',
      rounded ? '' : 'rounded-lg',
      className,
    ]"
    :style="{
      width: width || '100%',
      height: height || '1rem',
      borderRadius: rounded || undefined,
    }"
  >
    <slot />
  </div>
</template>

<style scoped>
.skeleton-loader {
  background-size: 200% 100%;
  animation: skeleton-pulse 1.5s ease-in-out infinite;
}

@keyframes skeleton-pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
