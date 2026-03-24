<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/20 backdrop-blur-sm"
        @click="handleCancel"
      >
        <div
          v-if="open"
          class="w-full admin-dialog p-6 max-h-[80vh transform transition-all flex flex-col overflow-hidden"
          :class="sizeClass"
          @click.stop
        >
          <div class="flex justify-between items-start mb-4 flex-shrink-0">
            <slot name="header">
              <h3 class="font-medium text-xl text-gray-800">
                {{ title }}
              </h3>
            </slot>
            <button
              @click="handleCancel"
              :disabled="loading"
              class="p-1 rounded-full hover:bg-black/5 transition-colors"
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
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import BaseIcon from "./BaseIcon.vue";

interface Props {
  open: boolean;
  title?: string;
  loading?: boolean;
  sizeClass?: string;
}

withDefaults(defineProps<Props>(), {
  title: "对话框",
  loading: false,
  sizeClass: "max-w-2xl",
});

const emit = defineEmits<{
  "update:open": [value: boolean];
  cancel: [];
}>();

const handleCancel = () => {
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
</style>
