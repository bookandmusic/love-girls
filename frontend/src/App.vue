<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-background font-(family-name:--font-body) overflow-x-hidden">
    <RouterView />
  </div>
</template>

<style scoped>
.bg-background {
  background-color: var(--background);
}

.overflow-x-hidden {
  overflow-x: hidden;
  overscroll-behavior: none;
}
</style>
