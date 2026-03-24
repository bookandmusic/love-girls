<template>
  <div class="w-full h-screen flex flex-col">
    <ToastNotification
      v-model:show="toast.show"
      :message="toast.message"
      :type="toast.type"
    />
    <RouterView />
  </div>
</template>

<script lang="ts" setup>
import { provide, ref } from "vue";

import ToastNotification from "@/components/ui/ToastNotification.vue";

const toast = ref({
  show: false,
  message: "",
  type: "info" as "success" | "error" | "info",
});

const showToast = (
  message: string,
  type: "success" | "error" | "info" = "info",
) => {
  toast.value = {
    show: true,
    message,
    type,
  };
};

provide("showToast", showToast);
</script>
