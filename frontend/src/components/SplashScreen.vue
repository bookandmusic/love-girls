<template>
  <Transition name="splash-fade">
    <div v-if="visible" class="splash-screen">
      <div class="splash-bg">
        <div class="gradient-layer"></div>
        <div class="blob blob-1"></div>
        <div class="blob blob-2"></div>
        <div class="blob blob-3"></div>
      </div>

      <div class="splash-content">
        <div class="logo-container">
          <div class="logo-glow"></div>
          <img src="/favicon.png" alt="Love Girl" class="logo" />
        </div>

        <div class="loading-indicator">
          <div class="heart-pulse">
            <svg viewBox="0 0 24 24" class="heart-icon">
              <path
                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                fill="currentColor"
              />
            </svg>
          </div>
          <p class="loading-text">{{ loadingText }}</p>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const visible = ref(true)
const loadingText = ref('正在启动...')
const texts = ['正在启动...', '正在连接服务器...', '正在加载资源...'] as const
let textIndex = 0
let textTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  textTimer = setInterval(() => {
    textIndex = (textIndex + 1) % texts.length
    loadingText.value = texts[textIndex]!
  }, 2000)
})

onUnmounted(() => {
  if (textTimer) {
    clearInterval(textTimer)
  }
})
</script>

<style scoped>
.splash-screen {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.splash-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #fce4ec 0%, #f8bbd9 50%, #f48fb1 100%);
}

.gradient-layer {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.3) 0%, transparent 50%);
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.6;
  animation: float 8s ease-in-out infinite;
}

.blob-1 {
  top: -20%;
  right: -10%;
  width: 60%;
  height: 60%;
  background: linear-gradient(135deg, #f0ada0 0%, #f48fb1 100%);
  animation-delay: 0s;
}

.blob-2 {
  bottom: -30%;
  left: -20%;
  width: 70%;
  height: 70%;
  background: linear-gradient(135deg, #d89388 0%, #f0ada0 100%);
  animation-delay: -2s;
}

.blob-3 {
  top: 40%;
  left: 60%;
  width: 40%;
  height: 40%;
  background: linear-gradient(135deg, #f8bbd9 0%, #f48fb1 100%);
  animation-delay: -4s;
}

@keyframes float {
  0%,
  100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.05);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.95);
  }
}

.splash-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 48px;
}

.logo-container {
  position: relative;
  animation: logoFloat 3s ease-in-out infinite;
}

.logo-glow {
  position: absolute;
  inset: -20px;
  background: radial-gradient(circle, rgba(244, 143, 177, 0.5) 0%, transparent 70%);
  animation: glow 2s ease-in-out infinite alternate;
}

.logo {
  position: relative;
  width: 100px;
  height: 100px;
  border-radius: 24px;
  box-shadow: 0 20px 60px rgba(244, 143, 177, 0.4);
}

@keyframes logoFloat {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes glow {
  0% {
    opacity: 0.5;
    transform: scale(0.9);
  }
  100% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

.loading-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.heart-pulse {
  animation: pulse 1.2s ease-in-out infinite;
}

.heart-icon {
  width: 32px;
  height: 32px;
  color: #d81b60;
  filter: drop-shadow(0 4px 8px rgba(216, 27, 96, 0.3));
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.2);
  }
}

.loading-text {
  font-size: 14px;
  color: #ad1457;
  font-weight: 500;
  letter-spacing: 0.5px;
  animation: textFade 2s ease-in-out infinite;
}

@keyframes textFade {
  0%,
  100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

.splash-fade-enter-active,
.splash-fade-leave-active {
  transition: opacity 0.3s ease-out;
}

.splash-fade-enter-from,
.splash-fade-leave-to {
  opacity: 0;
}
</style>
