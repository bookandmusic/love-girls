<script setup lang="ts">
import { computed } from "vue";

import type { Comment } from "@/services/commentApi";
import { useAuthStore } from "@/stores/auth";
import { useToast } from "@/utils/toastUtils";

const props = defineProps<{
  comment: Comment;
  momentId: number;
  depth?: number;
}>();

const emit = defineEmits<{
  reply: [comment: Comment];
  deleted: [commentId: number];
}>();

const authStore = useAuthStore();
const showToast = useToast();

const currentDepth = computed(() => props.depth ?? props.comment.depth);
const maxDisplayDepth = 3;
const shouldCollapseChildren = computed(
  () => currentDepth.value >= maxDisplayDepth,
);

const handleReply = () => {
  emit("reply", props.comment);
};

const handleDelete = async () => {
  try {
    const { commentApi } = await import("@/services/commentApi");
    const response = await commentApi.deleteComment(props.comment.id);
    if (response.code === 0) {
      showToast("评论已删除", "success");
      emit("deleted", props.comment.id);
    } else {
      showToast(response.msg || "删除失败", "error");
    }
  } catch {
    showToast("删除失败", "error");
  }
};

const canDelete = computed(() => {
  return props.comment.userId === authStore.userInfo?.userId;
});
</script>

<template>
  <div
    class="comment-item"
    :class="{ 'ml-6 mt-2': currentDepth > 0 && currentDepth < maxDisplayDepth }"
  >
    <div class="flex gap-3">
      <div
        class="w-8 h-8 rounded-full overflow-hidden bg-gray-200 flex-shrink-0 flex items-center justify-center"
      >
        <img
          v-if="comment.author.avatar?.thumbnail || comment.author.avatar?.url"
          :src="comment.author.avatar.thumbnail || comment.author.avatar.url"
          :alt="comment.author.name"
          class="w-full h-full object-cover"
        />
        <span v-else class="text-xs font-bold text-gray-500">
          {{ comment.author.name.substring(0, 1) }}
        </span>
      </div>

      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium text-[#576b95]">
            {{ comment.author.name }}
          </span>
          <span class="text-xs text-gray-400">{{ comment.createdAt }}</span>
        </div>

        <p class="text-sm text-gray-800 mt-1 break-words">
          <template v-if="comment.replyTo">
            <span class="text-gray-500">回复</span>
            <span class="text-[#576b95] font-medium mx-1">{{
              comment.replyTo.name
            }}</span>
            <span class="text-gray-500">：</span>
          </template>
          {{ comment.content }}
        </p>

        <div class="flex items-center gap-4 mt-2">
          <button
            @click="handleReply"
            class="text-xs text-gray-500 hover:text-[var(--fe-primary-dark)] transition-colors"
          >
            回复
          </button>
          <button
            v-if="canDelete"
            @click="handleDelete"
            class="text-xs text-gray-500 hover:text-red-500 transition-colors"
          >
            删除
          </button>
        </div>

        <template v-if="comment.children && comment.children.length > 0">
          <div v-if="!shouldCollapseChildren" class="mt-3 space-y-3">
            <CommentItem
              v-for="child in comment.children"
              :key="child.id"
              :comment="child"
              :moment-id="momentId"
              :depth="currentDepth + 1"
              @reply="(c) => emit('reply', c)"
              @deleted="(id) => emit('deleted', id)"
            />
          </div>
          <div v-else class="mt-2">
            <span class="text-xs text-[var(--fe-primary-dark)]">
              {{ comment.children.length }} 条回复
            </span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
