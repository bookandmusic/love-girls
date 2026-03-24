<template>
  <WaveBackground :waveColor="'#f0ada0'" class="flex-grow">
    <div class="w-full h-screen flex flex-col items-center justify-center p-4">
      <div class="max-w-md w-full text-center">
        <BaseIcon name="err404" size="w-24" color="text-[#f0ada0]" />

        <h1 class="text-3xl font-bold mt-6 mb-3 text-[var(--fe-text-primary)]">
          页面未找到
        </h1>

        <p class="mb-8 text-[var(--fe-text-secondary)]">
          抱歉，您访问的页面不存在或已被移除
        </p>

        <div class="flex flex-col sm:flex-row justify-center gap-4">
          <button
            @click="$router.back()"
            class="px-6 py-3 rounded-xl border-2 border-[#f0ada0] text-[#d89388] font-semibold hover:bg-[#f0ada0]/10 transition-colors"
          >
            返回上一页
          </button>
          <button
            @click="$router.push('/')"
            class="px-6 py-3 rounded-xl bg-gradient-to-r from-[#f0ada0] to-[#d89388] text-white font-semibold hover:from-[#d89388] hover:to-[#c78277] transition-all shadow-sm"
          >
            返回首页
          </button>
        </div>
      </div>
    </div>
  </WaveBackground>
</template>

<script setup lang="ts">
import BaseIcon from "@/components/ui/BaseIcon.vue";
import WaveBackground from "@/components/ui/WaveBackground.vue";
</script>
