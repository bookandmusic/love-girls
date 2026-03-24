<template>
  <Teleport to="body">
    <Transition name="action-sheet">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[300] flex items-end justify-center"
        @click="handleCancel"
      >
        <div class="absolute inset-0 bg-black/40"></div>

        <Transition name="action-sheet-content">
          <div
            v-if="modelValue"
            class="relative w-full max-w-lg mx-4 mb-8 rounded-2xl overflow-hidden bg-white/95 backdrop-blur-xl shadow-2xl"
            @click.stop
          >
            <div
              v-if="title"
              class="px-4 py-3 text-center border-b border-gray-100"
            >
              <p class="text-sm font-medium text-gray-500">
                {{ title }}
              </p>
            </div>

            <div class="p-2">
              <button
                v-for="(action, index) in visibleActions"
                :key="index"
                @click="handleAction(action)"
                class="w-full py-3.5 px-4 text-center rounded-xl transition-all active:bg-gray-100"
                :class="[
                  action.destructive
                    ? 'text-red-500 font-medium'
                    : 'text-[var(--fe-primary)] font-medium',
                  action.disabled
                    ? 'opacity-50 cursor-not-allowed'
                    : 'hover:bg-gray-50',
                ]"
                :disabled="action.disabled"
              >
                <div class="flex items-center justify-center gap-2">
                  <component
                    v-if="action.icon"
                    :is="action.icon"
                    class="w-5 h-5"
                  />
                  <span>{{ action.label }}</span>
                </div>
              </button>
            </div>

            <div class="p-2 pt-0">
              <button
                @click="handleCancel"
                class="w-full py-3.5 px-4 text-center rounded-xl font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 transition-all active:bg-gray-300"
              >
                取消
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { Component } from "vue";

export interface ActionSheetAction {
  label: string;
  icon?: Component;
  destructive?: boolean;
  disabled?: boolean;
  handler?: () => void;
}

interface Props {
  modelValue: boolean;
  title?: string;
  actions: ActionSheetAction[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", value: boolean): void;
  (e: "cancel"): void;
  (e: "select", action: ActionSheetAction): void;
}>();

const visibleActions = computed(() => props.actions);

const handleAction = (action: ActionSheetAction) => {
  if (action.disabled) return;

  emit("update:modelValue", false);
  emit("select", action);
  action.handler?.();
};

const handleCancel = () => {
  emit("update:modelValue", false);
  emit("cancel");
};
</script>

<style scoped>
.action-sheet-enter-active,
.action-sheet-leave-active {
  transition: opacity 0.25s ease;
}

.action-sheet-enter-from,
.action-sheet-leave-to {
  opacity: 0;
}

.action-sheet-content-enter-active {
  transition: transform 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}

.action-sheet-content-leave-active {
  transition: transform 0.2s ease-in;
}

.action-sheet-content-enter-from,
.action-sheet-content-leave-to {
  transform: translateY(100%);
}
</style>
