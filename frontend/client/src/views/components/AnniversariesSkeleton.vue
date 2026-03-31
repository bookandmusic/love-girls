<script setup lang="ts">
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
</script>

<template>
  <div class="flex flex-col h-full bg-[var(--fe-bg-gray)]/30">
    <div class="p-6 flex justify-center items-center text-center flex-shrink-0">
      <div
        class="w-full max-w-lg rounded-xl p-6 border border-white/40 shadow-lg"
      >
        <SkeletonLoader width="60px" height="10px" class="mx-auto mb-2" />
        <SkeletonLoader width="200px" height="28px" class="mx-auto" />
      </div>
    </div>

    <div class="overflow-y-auto px-4 md:px-8 pb-8">
      <div class="relative max-w-2xl mx-auto py-4">
        <div class="space-y-6">
          <div v-for="i in 3" :key="i" class="relative flex items-center">
            <div
              class="absolute left-6 w-3 h-3 rounded-full bg-gray-300 transform translate-x-[-50%] z-10"
            />
            <div class="ml-12 flex-grow">
              <div class="p-4 rounded-2xl">
                <SkeletonLoader width="120px" height="20px" class="mb-1" />
                <SkeletonLoader width="80px" height="12px" class="mb-2" />
                <SkeletonLoader width="100%" height="32px" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
