<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-background font-(family-name:--font-body)">
    <RouterView />
  </div>
</template>

<style scoped>
.bg-background {
  background-color: var(--background);
}
</style>
