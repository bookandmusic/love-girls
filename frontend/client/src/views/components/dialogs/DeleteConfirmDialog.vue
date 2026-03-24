<template>
  <GenericDialog
    :open="open"
    :title="title"
    @cancel="closeDialog"
    :loading="loading"
    size-class="max-w-sm"
  >
    <template #content>
      <p class="text-center text-gray-600">
        {{ message }}
      </p>
    </template>
    <template #actions>
      <button
        class="flex-1 py-3.5 text-center text-gray-500 font-medium hover:bg-gray-50 active:bg-gray-100 transition-colors"
        @click="closeDialog"
      >
        取消
      </button>
      <button
        class="flex-1 py-3.5 text-center text-red-500 font-semibold border-l border-gray-100 hover:bg-red-50 active:bg-red-100 transition-colors"
        @click="handleConfirm"
      >
        删除
      </button>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import GenericDialog from "@/components/ui/GenericDialog.vue";

interface Props {
  open: boolean;
  loading?: boolean;
  title?: string;
  message?: string;
}

withDefaults(defineProps<Props>(), {
  loading: false,
  title: "确认删除",
  message: "确定要删除吗？此操作无法撤销。",
});

const emit = defineEmits<{
  (e: "update:open", open: boolean): void;
  (e: "confirm"): void;
  (e: "cancel"): void;
}>();

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
};

const handleConfirm = () => {
  emit("confirm");
};
</script>
