<template>
  <div class="fixed inset-0 z-999999 flex items-center justify-center bg-white/20 backdrop-blur-md">
    <div class="flex flex-col items-center gap-8">
      <div class="relative w-24 h-24">
        <div
          class="absolute inset-0 border-4 rounded-full animate-pulse"
          style="border-color: var(--primary-light); border-top-color: var(--primary-color); animation-duration: 0.8s;"
        ></div>
        <div
          class="absolute inset-2 border-4 rounded-full animate-pulse"
          style="border-color: var(--primary-light); border-right-color: var(--primary-color); animation-duration: 1s; animation-delay: -0.2s;"
        ></div>
        <div
          class="absolute inset-4 border-4 rounded-full animate-pulse"
          style="border-color: var(--primary-light); border-bottom-color: var(--primary-color); animation-duration: 1.2s; animation-delay: -0.4s;"
        ></div>
        <div class="absolute inset-0 flex items-center justify-center">
          <span class="text-3xl animate-pulse" style="animation-duration: 1.5s;">💕</span>
        </div>
      </div>
      <div class="text-base tracking-wide font-decor" style="color: var(--primary-color);">正在加载…</div>
    </div>
  </div>
</template>

<style scoped>
@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: rotate(0deg);
  }
  50% {
    opacity: 0.7;
    transform: rotate(180deg);
  }
}

.animate-pulse {
  animation: pulse var(--animation-duration, 1s) ease-in-out infinite;
}
</style>
