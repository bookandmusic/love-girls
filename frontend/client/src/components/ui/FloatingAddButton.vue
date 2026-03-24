<template>
  <button
    @click="$emit('click')"
    :disabled="loading"
    class="fixed right-6 bottom-24 md:bottom-8 z-[100] w-14 h-14 rounded-full shadow-lg transition-all active:scale-90 tap-feedback flex items-center justify-center"
    :class="[
      loading
        ? 'bg-gray-300 cursor-not-allowed'
        : 'bg-gradient-to-br from-[var(--fe-primary)] to-[var(--fe-primary-dark)] hover:shadow-xl hover:scale-105',
    ]"
  >
    <svg
      v-if="!loading"
      class="w-7 h-7 text-white"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2.5"
        d="M12 4v16m8-8H4"
      />
    </svg>
    <svg
      v-else
      class="w-6 h-6 text-white animate-spin"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      ></circle>
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
      ></path>
    </svg>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  loading?: boolean;
}>();

defineEmits<{
  (e: "click"): void;
}>();
</script>
