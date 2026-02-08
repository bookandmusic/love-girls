<template>
  <li
    class="relative border border-gray-200 rounded-lg mb-3 overflow-hidden shadow-sm hover:shadow-md transition-shadow"
  >
    <!-- 背景层（可选） -->
    <slot name="background"></slot>

    <!-- 遮罩层（可选） -->
    <slot name="overlay"></slot>

    <!-- 内容层 -->
    <div class="relative z-10 block p-4">
      <div>
        <slot name="header"></slot>
      </div>

      <div class="mt-3">
        <slot name="content"></slot>
      </div>

      <div class="mt-6">
        <slot name="footer"></slot>
      </div>
    </div>
  </li>
</template>

<script setup lang="ts">
// 通用卡片组件
// 支持 background / overlay / header / content / footer 分层插槽
</script>
