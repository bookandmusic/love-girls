<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        :class="[
          'fixed inset-0 z-50 flex items-center justify-center p-4',
          variant === 'admin' ? '' : 'bg-black/40 backdrop-blur-sm',
        ]"
        @click="handleCancel"
      >
        <Transition :name="variant === 'admin' ? '' : 'scale'" appear>
          <div
            v-if="open"
            :class="[
              'w-full transform transition-all flex flex-col overflow-hidden',
              variant === 'admin'
                ? 'admin-dialog p-6 max-h-[80vh]'
                : 'bg-white/90 backdrop-blur-xl rounded-2xl shadow-2xl',
              sizeClass,
            ]"
            @click.stop
          >
            <!-- 前台 iOS 风格布局 -->
            <template v-if="variant !== 'admin'">
              <div class="p-5 border-b border-gray-100">
                <h3 class="text-lg font-semibold text-gray-900 text-center">
                  {{ title }}
                </h3>
              </div>
              <div
                class="text-sm text-gray-600 leading-relaxed overflow-y-auto max-h-[60vh] p-5"
              >
                <slot name="content"></slot>
              </div>

              <!-- 底部动作区域 - iOS Alert 风格 -->
              <div class="border-t border-gray-100 flex">
                <slot name="actions"></slot>
              </div>
            </template>

            <!-- 后台原有布局保持不变 -->
            <template v-else>
              <div class="flex justify-between items-start mb-4 flex-shrink-0">
                <slot name="header">
                  <h3 class="font-medium text-xl">
                    {{ title }}
                  </h3>
                </slot>
                <button
                  @click="handleCancel"
                  :disabled="loading"
                  class="p-1 rounded-full hover:bg-black/5"
                >
                  <BaseIcon
                    name="close"
                    size="w-6 h-6"
                    color="var(--admin-accent-color)"
                  />
                </button>
              </div>

              <div class="space-y-4 flex-1 overflow-y-auto min-h-0">
                <slot name="content"></slot>
              </div>

              <div class="mt-6 flex flex-wrap gap-3 justify-end flex-shrink-0">
                <slot name="actions"></slot>
              </div>
            </template>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import BaseIcon from "./BaseIcon.vue";

// 组件的 Props
interface Props {
  open: boolean;
  title?: string;
  loading?: boolean;
  sizeClass?: string;
  variant?: "default" | "admin";
}

withDefaults(defineProps<Props>(), {
  title: "对话框",
  loading: false,
  sizeClass: "max-w-2xl",
  variant: "default",
});

const emit = defineEmits<{
  "update:open": [value: boolean];
  cancel: [];
}>();

const handleCancel = () => {
  // 触发取消事件，通知父组件关闭对话框
  emit("cancel");
  emit("update:open", false);
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.scale-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.scale-leave-active {
  transition: all 0.2s ease-in;
}
.scale-enter-from {
  transform: scale(0.9);
  opacity: 0;
}
.scale-leave-to {
  transform: scale(0.95);
  opacity: 0;
}
</style>
