<script setup lang="ts">
import { RouterView } from "vue-router";
import SplashScreen from "@/components/SplashScreen.vue";
import { useUIStore } from "@/stores/ui";
import { onMounted, ref } from "vue";
import { useGlobalRefresh } from "@/composables/useGlobalRefresh";

const uiStore = useUIStore();
const { initGlobalState, startPeriodicRefresh } = useGlobalRefresh();
const showSplash = ref(true);

onMounted(async () => {
  uiStore.setAppReady(false);

  const success = await initGlobalState();
  if (success) {
    startPeriodicRefresh(60000);
    setTimeout(() => {
      showSplash.value = false;
      uiStore.setAppReady(true);
    }, 1500);
  } else {
    setTimeout(() => {
      showSplash.value = false;
    }, 1500);
  }
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
