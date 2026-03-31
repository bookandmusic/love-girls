<script setup lang="ts">
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
</script>

<template>
  <div class="flex flex-col h-full p-4 md:p-8">
    <div
      v-for="i in 3"
      :key="i"
      class="py-6 border-b border-black/5 last:border-0"
    >
      <div class="flex items-start">
        <SkeletonLoader
          width="48px"
          height="48px"
          rounded="8px"
          class-name="mr-4 flex-shrink-0"
        />
        <div class="flex-grow min-w-0">
          <SkeletonLoader width="100px" height="16px" class-name="mb-2" />
          <SkeletonLoader width="100%" height="40px" class-name="mb-3" />
          <div class="flex gap-1 mb-3">
            <SkeletonLoader width="80px" height="80px" rounded="6px" />
            <SkeletonLoader width="80px" height="80px" rounded="6px" />
            <SkeletonLoader width="80px" height="80px" rounded="6px" />
          </div>
          <div class="flex justify-between items-center">
            <SkeletonLoader width="80px" height="12px" />
            <div class="flex gap-2">
              <SkeletonLoader width="50px" height="24px" rounded="6px" />
              <SkeletonLoader width="50px" height="24px" rounded="6px" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
