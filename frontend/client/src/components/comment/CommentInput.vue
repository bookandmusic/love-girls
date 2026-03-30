<script setup lang="ts">
import { ref, computed } from "vue";

import type { Comment } from "@/services/commentApi";
import { commentApi } from "@/services/commentApi";
import { useToast } from "@/utils/toastUtils";

const props = defineProps<{
  momentId: number;
  replyTo?: Comment | null;
  inline?: boolean;
}>();

const emit = defineEmits<{
  created: [comment: Comment];
  cancel: [];
}>();

const showToast = useToast();

const content = ref("");
const submitting = ref(false);

const placeholder = computed(() => {
  if (props.replyTo) {
    return `回复 ${props.replyTo.author.name}...`;
  }
  return "写下你的评论...";
});

const canSubmit = computed(() => {
  return content.value.trim().length > 0;
});

const handleSubmit = async () => {
  if (!canSubmit.value || submitting.value) return;

  submitting.value = true;
  try {
    const data: {
      content: string;
      parentId?: number;
      replyToId?: number;
    } = {
      content: content.value.trim(),
    };

    if (props.replyTo) {
      data.parentId = props.replyTo.parentId || props.replyTo.id;
      data.replyToId = props.replyTo.id;
    }

    const response = await commentApi.createComment(props.momentId, data);
    if (response.code === 0) {
      showToast("评论成功", "success");
      content.value = "";
      emit("created", response.data);
    } else {
      showToast(response.msg || "评论失败", "error");
    }
  } catch {
    showToast("评论失败", "error");
  } finally {
    submitting.value = false;
  }
};

const handleCancel = () => {
  content.value = "";
  emit("cancel");
};
</script>

<template>
  <div
    :class="
      inline
        ? 'comment-input-inline'
        : 'comment-input bg-gray-50 rounded-xl p-3 mb-3'
    "
  >
    <div v-if="replyTo" class="flex items-center justify-between mb-2">
      <span class="text-xs text-gray-500">
        回复 <span class="text-[#576b95]">{{ replyTo.author.name }}</span>
      </span>
      <button
        @click="handleCancel"
        class="text-xs text-gray-400 hover:text-gray-600 transition-colors"
      >
        取消回复
      </button>
    </div>

    <textarea
      ref="textareaRef"
      v-model="content"
      :placeholder="placeholder"
      :disabled="submitting"
      :rows="inline ? 3 : 3"
      class="w-full bg-transparent border-none resize-none text-sm focus:outline-none placeholder:text-gray-400"
      :class="{ 'bg-gray-50 rounded-lg p-2': inline }"
    ></textarea>

    <div class="flex items-center justify-end gap-2 mt-2">
      <button
        v-if="!inline && replyTo"
        @click="handleCancel"
        class="px-4 py-1.5 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
      >
        取消
      </button>
      <button
        @click="handleSubmit"
        :disabled="!canSubmit || submitting"
        class="px-4 py-1.5 text-sm text-white bg-[var(--fe-primary-dark)] rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed hover:opacity-90"
      >
        {{ submitting ? "发送中..." : "发送" }}
      </button>
    </div>
  </div>
</template>
