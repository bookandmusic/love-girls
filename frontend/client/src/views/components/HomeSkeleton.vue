<script setup lang="ts">
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
</script>

<template>
  <div
    class="flex flex-col items-center justify-start min-h-full p-4 space-y-2 md:space-y-6 pt-2 md:pt-10"
  >
    <div class="w-full max-w-4xl flex flex-col items-center">
      <div class="relative w-full py-2 md:py-4">
        <div class="flex justify-center items-center gap-8 md:gap-16">
          <div class="flex flex-col items-center">
            <SkeletonLoader width="100px" height="100px" rounded="50%" />
            <SkeletonLoader width="60px" height="20px" class-name="mt-2" />
          </div>
          <div class="flex flex-col items-center">
            <SkeletonLoader width="100px" height="100px" rounded="50%" />
            <SkeletonLoader width="60px" height="20px" class-name="mt-2" />
          </div>
        </div>
      </div>

      <div
        class="w-full max-w-md p-4 md:p-6 mt-2 md:mt-4 rounded-xl border border-white/40 shadow-lg text-center"
      >
        <SkeletonLoader width="120px" height="10px" class-name="mx-auto mb-2" />
        <SkeletonLoader width="200px" height="40px" class-name="mx-auto mb-1" />
        <SkeletonLoader width="160px" height="24px" class-name="mx-auto" />
      </div>
    </div>
  </div>
</template>
