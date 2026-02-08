<template>
  <WaveBackground :waveColor="'#f0ada0'" class="flex-grow">
    <div class="w-full h-screen flex flex-col items-center justify-center p-4">
      <div class="max-w-md w-full text-center text-[#FF7500]">
        <!-- Win11风格的404图标 -->
        <div class="flex items-center justify-center">
          <BaseIcon name="err404" size="w-24" color="text-[#ff7500]" />
        </div>

        <!-- 标题 -->
        <h1 class="text-3xl font-bold mt-4 mb-2 font-[ZCOOL_KuaiLe]">页面未找到</h1>

        <!-- 描述 -->
        <p class="mb-8 font-[Ma_Shan_Zheng]">抱歉，您访问的页面不存在或已被移除</p>

        <!-- 操作按钮 -->
        <div class="flex flex-col flex-row justify-center gap-4">
          <button @click="$router.back()" class="win11-button outline px-6 py-2 rounded-lg">
            返回上一页
          </button>
          <button @click="$router.push('/')" class="win11-button px-6 py-2 rounded-lg">
            返回首页
          </button>
        </div>
      </div>
    </div>
  </WaveBackground>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import WaveBackground from '@/components/ui/WaveBackground.vue'
</script>

<style scoped></style>
