<script setup lang="ts">
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="grid grid-cols-2 md:grid-cols-3 gap-4 p-4 md:p-8">
      <div
        v-for="i in 6"
        :key="i"
        class="aspect-square rounded-xl overflow-hidden"
      >
        <SkeletonLoader width="100%" height="100%" rounded="12px" />
      </div>
    </div>
  </div>
</template>
