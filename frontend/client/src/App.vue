<script setup lang="ts">
import { RouterView } from "vue-router";
import SplashScreen from "@/components/SplashScreen.vue";
import { useUIStore } from "@/stores/ui";
import { useNotificationStore } from "@/stores/notification";
import { onMounted, onUnmounted, ref, watch } from "vue";
import { useAuthStore } from "@/stores/auth";

const uiStore = useUIStore();
const notificationStore = useNotificationStore();
const authStore = useAuthStore();
const showSplash = ref(true);

onMounted(() => {
  setTimeout(() => {
    showSplash.value = false;
    uiStore.setAppReady(true);
  }, 1500);
});

watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      notificationStore.startPolling();
    } else {
      notificationStore.stopPolling();
    }
  },
  { immediate: true },
);

onUnmounted(() => {
  notificationStore.stopPolling();
});
</script>

<template>
  <SplashScreen v-if="showSplash" />
  <RouterView v-else />
</template>

<style>
html,
body,
#app {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>
