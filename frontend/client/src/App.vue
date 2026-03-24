<script setup lang="ts">
import { RouterView } from "vue-router";
import SplashScreen from "@/components/SplashScreen.vue";
import { useUIStore } from "@/stores/ui";
import { onMounted, ref } from "vue";

const uiStore = useUIStore();
const showSplash = ref(true);

onMounted(() => {
  setTimeout(() => {
    showSplash.value = false;
    uiStore.setAppReady(true);
  }, 1500);
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
