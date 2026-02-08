<template>
  <div class="generic-card animated-border" :style="borderStyle" @mouseenter="markHovered">
    <!-- Header -->
    <div v-if="$slots.header" class="mb-4">
      <slot name="header" />
    </div>

    <!-- Body -->
    <div>
      <slot />
    </div>

    <!-- Footer -->
    <div v-if="$slots.footer" class="mt-4">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  /**
   * 边框颜色（任意 CSS 颜色）
   * 例：#22c55e / rgb(34,197,94) / hsl(...)
   */
  borderColor: {
    type: String,
    default: 'rgba(59,130,246,0.5)',
  },

  /**
   * 边框宽度
   */
  borderWidth: {
    type: Number,
    default: 3,
  },

  /**
   * 动画时长（秒）
   */
  animationDuration: {
    type: Number,
    default: 0.5,
  },
})

const markHovered = (e: MouseEvent) => {
  const target = e.currentTarget as HTMLElement | null
  if (target) {
    target.classList.add('hovered')
  }
}

/**
 * 注入给 ::after 使用的 CSS 变量
 */
const borderStyle = computed(() => ({
  '--card-border-color': props.borderColor,
  '--card-border-width': `${props.borderWidth}px`,
  '--card-border-duration': `${props.animationDuration}s`,
}))
</script>

<style scoped>
/* =========================
   基础卡片样式
   ========================= */
:global(.dark) .card-base {
  background: rgba(31, 41, 55, 0.5);
  border-color: rgba(55, 65, 81, 0.6);
}

/* =========================
   动画容器
   ========================= */
.animated-border {
  position: relative;
}

/* 伪元素边框 */
.animated-border::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  border: var(--card-border-width) solid var(--card-border-color);
  opacity: 0;
  pointer-events: none;
}

/* =========================
   Hover 进入动画
   ========================= */
.animated-border:hover::after {
  animation: border-spin var(--card-border-duration) ease-out forwards;
}

/* Hover 离开动画（只对 hover 过的生效） */
.animated-border.hovered:not(:hover)::after {
  animation: border-unspin var(--card-border-duration) ease-out forwards;
}

/* =========================
   动画定义
   ========================= */
@keyframes border-spin {
  0% {
    clip-path: inset(0 100% 100% 0);
    opacity: 0;
  }

  50% {
    clip-path: inset(0 0 100% 0);
    opacity: 1;
  }

  100% {
    clip-path: inset(0 0 0 0);
    opacity: 1;
  }
}

@keyframes border-unspin {
  0% {
    clip-path: inset(0 0 0 0);
    opacity: 1;
  }

  50% {
    clip-path: inset(100% 0 0 0);
    opacity: 1;
  }

  100% {
    clip-path: inset(100% 0 0 100%);
    opacity: 0;
  }
}
</style>
