<script setup lang="ts">
import { computed, ref } from "vue";

import ActionSheet, {
  type ActionSheetAction,
} from "@/components/ui/ActionSheet.vue";
import type { Comment } from "@/services/commentApi";
import { useAuthStore } from "@/stores/auth";
import { useToast } from "@/utils/toastUtils";

const props = defineProps<{
  comment: Comment;
  momentId: number;
  depth?: number;
  embedded?: boolean;
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

const showActionSheet = ref(false);

const isOwnComment = computed(() => {
  return props.comment.userId === authStore.userInfo?.userId;
});

const handleCommentClick = () => {
  if (isOwnComment.value) {
    showActionSheet.value = true;
  } else {
    emit("reply", props.comment);
  }
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

const actionSheetActions = computed<ActionSheetAction[]>(() => [
  {
    label: "删除",
    destructive: true,
    handler: handleDelete,
  },
]);
</script>

<template>
  <div class="comment-item py-1">
    <div class="flex items-start">
      <div class="flex-1 min-w-0">
        <p
          @click="handleCommentClick"
          class="text-[13px] leading-relaxed break-words cursor-pointer"
        >
          <span class="text-[#576b95] font-medium">{{
            comment.author.name
          }}</span>
          <template v-if="comment.replyTo">
            <span class="text-gray-500 mx-1">回复</span>
            <span class="text-[#576b95] font-medium">{{
              comment.replyTo.name
            }}</span>
          </template>
          <span class="text-gray-500">：</span>
          <span class="text-gray-800">{{ comment.content }}</span>
        </p>

        <template v-if="comment.children && comment.children.length > 0">
          <div
            v-if="!shouldCollapseChildren"
            class="mt-1 space-y-1 pl-3 border-l-2 border-gray-100"
          >
            <CommentItem
              v-for="child in comment.children"
              :key="child.id"
              :comment="child"
              :moment-id="momentId"
              :depth="currentDepth + 1"
              :embedded="embedded"
              @reply="(c) => emit('reply', c)"
              @deleted="(id) => emit('deleted', id)"
            />
          </div>
          <div v-else class="mt-1">
            <span class="text-xs text-[var(--fe-primary-dark)]">
              {{ comment.children.length }} 条回复
            </span>
          </div>
        </template>
      </div>
    </div>

    <ActionSheet
      v-model="showActionSheet"
      title="评论操作"
      :actions="actionSheetActions"
    />
  </div>
</template>
