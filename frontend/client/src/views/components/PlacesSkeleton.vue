<script setup lang="ts">
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
</script>

<template>
  <div
    class="flex-grow flex flex-col overflow-hidden bg-[var(--fe-bg-gray)]/30"
  >
    <div class="p-4 md:p-6 flex-shrink-0">
      <SkeletonLoader
        width="100%"
        height="256px"
        rounded="12px"
        class="md:h-96"
      />
    </div>

    <div class="overflow-y-auto px-4 md:px-8 pb-8">
      <SkeletonLoader width="120px" height="12px" class="mb-3 ml-1" />
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="i in 6" :key="i" class="flex items-center p-4 rounded-2xl">
          <SkeletonLoader
            width="48px"
            height="48px"
            rounded="12px"
            class="mr-4"
          />
          <div class="flex-grow">
            <SkeletonLoader width="100px" height="16px" class="mb-1" />
            <SkeletonLoader width="80px" height="12px" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
