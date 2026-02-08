<template>
  <div class="fixed inset-0 z-999999 flex items-center justify-center bg-white/70">
    <div class="flex flex-col items-center gap-6">
      <!-- Win11 风格旋转圆点 使用主题色 -->
      <div class="relative w-14 h-14">
        <span
          v-for="i in 8"
          :key="i"
          class="absolute top-1/2 left-1/2 w-2.5 h-2.5 rounded-full origin-[0_0] animate-spin"
          :style="{
            transform: `rotate(${i * 45}deg) translate(0, -24px)`,
            backgroundColor: 'var(--primary-color)',
          }"
        />
      </div>

      <!-- 文本 使用主题色 -->
      <div class="text-sm text-on-background tracking-wide">正在加载…</div>
    </div>
  </div>
</template>

<style scoped>
/* 可以在这里添加额外的样式 */
</style>
