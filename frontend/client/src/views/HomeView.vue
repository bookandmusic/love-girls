<script setup lang="ts">
import { getLunarDate } from "chinese-days";
import { computed, onMounted, ref } from "vue";

import MainLayout from "@/layouts/MainLayout.vue";
import { useSystemStore } from "@/stores/system";
import { useToast } from "@/utils/toastUtils";

import CoupleAvatar from "./components/CoupleAvatar.vue";
import HomeSkeleton from "./components/HomeSkeleton.vue";

const systemStore = useSystemStore();
const systemInfo = computed(() => systemStore.getSystemInfo);
const date = ref<string>("2025-12-23");
const showToast = useToast();
const isLoading = ref(true);

onMounted(() => {
  systemStore
    .fetchSystemInfo()
    .then(() => {
      if (systemInfo.value) {
        date.value = systemInfo.value.site.startDate;
      }
    })
    .catch(() => {
      showToast("获取系统信息失败", "error");
    })
    .finally(() => {
      isLoading.value = false;
    });
});

const lunarDate = computed(() => {
  return getLunarDate(date.value);
});
</script>

<template>
  <MainLayout
    :title="
      systemInfo?.couple?.boy?.name + ' ❤️ ' + systemInfo?.couple?.girl?.name
    "
    :subtitle="systemInfo?.site?.description"
    :start-date="systemInfo?.site?.startDate"
  >
    <template #main-content>
      <HomeSkeleton v-if="isLoading" />
      <div
        v-else-if="systemInfo"
        class="flex flex-col items-center justify-start min-h-full p-4 space-y-2 md:space-y-6 pt-2 md:pt-10"
      >
        <div class="w-full max-w-4xl flex flex-col items-center">
          <div class="relative w-full py-2 md:py-4">
            <CoupleAvatar
              :boy="systemInfo.couple.boy"
              :girl="systemInfo.couple.girl"
            ></CoupleAvatar>
          </div>

          <div
            class="glass-regular rounded-[var(--fe-radius-card)] p-4 md:p-6 mt-2 md:mt-4 border border-white/40 shadow-lg text-center ios-transition"
          >
            <div
              class="text-[var(--fe-text-secondary)] font-bold mb-1 tracking-widest uppercase text-[10px] opacity-60"
            >
              我们的故事开始于
            </div>
            <div
              class="text-3xl md:text-5xl font-bold text-[var(--fe-primary)] font-(family-name:--font-math) mb-1"
            >
              {{ date }}
            </div>
            <div
              class="text-sm md:text-lg font-medium text-[var(--fe-text-primary)] font-(family-name:--font-decor)"
            >
              {{
                lunarDate.lunarYearCN +
                "年" +
                lunarDate.lunarMonCN +
                lunarDate.lunarDayCN
              }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>
</template>
