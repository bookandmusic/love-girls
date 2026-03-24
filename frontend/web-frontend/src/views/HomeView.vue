<script setup lang="ts">
import { getLunarDate } from "chinese-days";
import { computed, onMounted, ref } from "vue";

import MainLayout from "@/layouts/MainLayout.vue";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";

import CoupleAvatar from "./components/CoupleAvatar.vue";

const uiStore = useUIStore();
const systemStore = useSystemStore();

const systemInfo = computed(() => systemStore.getSystemInfo);
const loadError = computed(() => systemStore.loadError);
const date = ref<string>("2024-01-01");

onMounted(async () => {
  uiStore.setLoading(true);
  try {
    const info = await systemStore.fetchSystemInfo();
    date.value = info.site.startDate;
  } finally {
    uiStore.setLoading(false);
  }
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
      <div
        class="flex flex-col items-center justify-start min-h-full p-4 space-y-2 md:space-y-6 pt-2 md:pt-10"
      >
        <div class="w-full max-w-4xl flex flex-col items-center">
          <div
            v-if="loadError"
            class="w-full mb-4 p-3 rounded-xl bg-amber-50 border border-amber-200 text-amber-700 text-center text-sm"
          >
            无法连接服务器，请联系管理员初始化项目
          </div>

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
