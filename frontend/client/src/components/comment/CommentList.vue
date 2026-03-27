<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue";

import type { Comment } from "@/services/commentApi";
import { commentApi } from "@/services/commentApi";
import { useToast } from "@/utils/toastUtils";

import CommentItem from "./CommentItem.vue";
import CommentInput from "./CommentInput.vue";

const props = defineProps<{
  momentId: number;
}>();

const showToast = useToast();

const comments = ref<Comment[]>([]);
const loading = ref(false);
const total = ref(0);
const page = ref(1);
const hasMore = computed(() => comments.value.length < total.value);

const showInput = ref(false);
const replyingTo = ref<Comment | null>(null);

const fetchComments = async (reset = false) => {
  if (loading.value) return;
  if (!reset && !hasMore.value) return;

  loading.value = true;
  try {
    if (reset) {
      page.value = 1;
      comments.value = [];
    }

    const response = await commentApi.getComments(props.momentId, page.value);
    if (response.code === 0) {
      if (reset) {
        comments.value = response.data.comments || [];
      } else {
        comments.value = [...comments.value, ...(response.data.comments || [])];
      }
      total.value = response.data.total;
      page.value = response.data.page + 1;
    } else {
      showToast(response.msg || "获取评论失败", "error");
    }
  } catch (error) {
    console.error("获取评论失败:", error);
    showToast("获取评论失败", "error");
  } finally {
    loading.value = false;
  }
};

const handleReply = (comment: Comment) => {
  replyingTo.value = comment;
  showInput.value = true;
};

const handleCommentCreated = (newComment: Comment) => {
  if (newComment.parentId) {
    const addToParent = (items: Comment[]): boolean => {
      for (const item of items) {
        if (item.id === newComment.parentId) {
          if (!item.children) {
            item.children = [];
          }
          item.children.push(newComment);
          return true;
        }
        if (item.children && addToParent(item.children)) {
          return true;
        }
      }
      return false;
    };
    addToParent(comments.value);
  } else {
    comments.value.push(newComment);
  }
  total.value++;
  showInput.value = false;
  replyingTo.value = null;
};

const handleCommentDeleted = (commentId: number) => {
  const removeFromList = (items: Comment[]): Comment[] => {
    return items
      .filter((item) => item.id !== commentId)
      .map((item) => ({
        ...item,
        children: item.children ? removeFromList(item.children) : [],
      }));
  };
  comments.value = removeFromList(comments.value);
  total.value--;
};

const openNewComment = () => {
  replyingTo.value = null;
  showInput.value = true;
};

const loadComments = () => {
  fetchComments(true);
};

onMounted(() => {
  loadComments();
});

watch(
  () => props.momentId,
  (newVal, oldVal) => {
    if (newVal !== oldVal) {
      fetchComments(true);
    }
  },
);

defineExpose({
  fetchComments,
  total,
});
</script>

<template>
  <div class="comment-list flex flex-col h-full">
    <div class="flex-shrink-0">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-medium text-gray-700">
          评论 ({{ total }})
        </span>
        <button
          @click="openNewComment"
          class="text-sm text-[var(--fe-primary-dark)] font-medium hover:opacity-80 transition-opacity"
        >
          发表评论
        </button>
      </div>

      <CommentInput
        v-if="showInput"
        :moment-id="momentId"
        :reply-to="replyingTo"
        @created="handleCommentCreated"
        @cancel="showInput = false"
      />
    </div>

    <div class="flex-1 overflow-y-auto min-h-0">
      <div v-if="comments.length === 0 && !loading" class="py-6 text-center">
        <p class="text-sm text-gray-500">暂无评论，快来发表第一条评论吧</p>
      </div>

      <div v-else class="space-y-3">
        <CommentItem
          v-for="comment in comments"
          :key="comment.id"
          :comment="comment"
          :moment-id="momentId"
          @reply="handleReply"
          @deleted="handleCommentDeleted"
        />
      </div>

      <div v-if="hasMore" class="py-4 text-center">
        <button
          @click="fetchComments()"
          :disabled="loading"
          class="text-sm text-[var(--fe-primary-dark)] hover:opacity-80 transition-opacity disabled:opacity-50"
        >
          {{ loading ? "加载中..." : "加载更多评论" }}
        </button>
      </div>
    </div>
  </div>
</template>
