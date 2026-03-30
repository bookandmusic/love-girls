<script setup lang="ts">
import { getSolarDateFromLunar } from "chinese-days";
import { computed, onMounted, ref } from "vue";
import { PullRefresh as VanPullRefresh } from "vant";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import ActionSheet, {
  type ActionSheetAction,
} from "@/components/ui/ActionSheet.vue";
import FloatingAddButton from "@/components/ui/FloatingAddButton.vue";
import { useLongPress } from "@/composables/useLongPress";
import MainLayout from "@/layouts/MainLayout.vue";
import { type Anniversary, anniversaryApi } from "@/services/anniversaryApi";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

import AnniversaryEditDialog from "./components/dialogs/AnniversaryEditDialog.vue";
import DeleteConfirmDialog from "./components/dialogs/DeleteConfirmDialog.vue";

interface NormalizedAnniversary {
  anniversary: Anniversary;
  nextDate: Date;
}

const uiStore = useUIStore();
const systemStore = useSystemStore();

const anniversaries = ref<Anniversary[]>([]);
const systemInfo = computed(() => systemStore.getSystemInfo);
const showToast = useToast();
const isRefreshing = ref(false);
const isAtTop = ref(true);
const scrollContainer = ref<HTMLElement | null>(null);

const fetchAnniversaries = async () => {
  try {
    uiStore.setLoading(true);
    const res = await anniversaryApi.getAnniversaries(1, 100);
    anniversaries.value = res.data.anniversaries;
  } catch {
    showToast("获取纪念日失败", "error");
  } finally {
    uiStore.setLoading(false);
  }
};

const parseMonthDay = (date: string) => {
  const [, month, day] = date.split("-").map(Number);
  if (!month || !day || isNaN(month) || isNaN(day)) return null;
  return { month, day };
};

const toSolarDate = (anniversary: Anniversary, year: number): Date | null => {
  const parsed = parseMonthDay(anniversary.date);
  if (!parsed) return null;

  const { month, day } = parsed;

  try {
    if (anniversary.calendar === "lunar") {
      const lunarStr = `${year.toString().padStart(4, "0")}-${month
        .toString()
        .padStart(2, "0")}-${day.toString().padStart(2, "0")}`;

      const result = getSolarDateFromLunar(lunarStr);

      if (result && result.date) {
        const [y, m, d] = result.date.split("-").map(Number);
        if (
          y != null &&
          m != null &&
          d != null &&
          !isNaN(y) &&
          !isNaN(m) &&
          !isNaN(d)
        ) {
          return new Date(y, m - 1, d);
        }
      }
      return null;
    }

    return new Date(year, month - 1, day);
  } catch (err) {
    console.error("转换农历失败:", err);
    return null;
  }
};

const getNextDate = (anniversary: Anniversary): Date | null => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const year = today.getFullYear();

  let date = toSolarDate(anniversary, year);
  if (!date) return null;

  if (date < today) {
    date = toSolarDate(anniversary, year + 1);
  }

  return date;
};

const normalizedAnniversaries = computed<NormalizedAnniversary[]>(() => {
  const list: NormalizedAnniversary[] = [];

  for (const anniversary of anniversaries.value) {
    const nextDate = getNextDate(anniversary);
    if (nextDate) {
      list.push({ anniversary, nextDate });
    }
  }

  return list.sort((a, b) => a.nextDate.getTime() - b.nextDate.getTime());
});

const sortedAnniversaries = computed(() =>
  normalizedAnniversaries.value.map((item) => item.anniversary),
);

const upcomingAnniversaryInfo = computed(() => {
  if (normalizedAnniversaries.value.length === 0) return null;

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const next = normalizedAnniversaries.value[0];
  if (!next) return null;

  const diffTime = next.nextDate.getTime() - today.getTime();
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

  return {
    title: next.anniversary.title,
    date: next.nextDate,
    diffDays,
    description: next.anniversary.description,
    isToday: diffDays === 0,
    anniversary: next.anniversary,
    days: diffDays,
    message:
      diffDays === 0
        ? `今天是${next.anniversary.title}`
        : `距离${next.anniversary.title}还有${diffDays}天`,
  };
});

const formatLunarDate = (month: number, day: number) => {
  const months = [
    "正月",
    "二月",
    "三月",
    "四月",
    "五月",
    "六月",
    "七月",
    "八月",
    "九月",
    "十月",
    "十一月",
    "十二月",
  ];
  const days = [
    "初一",
    "初二",
    "初三",
    "初四",
    "初五",
    "初六",
    "初七",
    "初八",
    "初九",
    "初十",
    "十一",
    "十二",
    "十三",
    "十四",
    "十五",
    "十六",
    "十七",
    "十八",
    "十九",
    "二十",
    "廿一",
    "廿二",
    "廿三",
    "廿四",
    "廿五",
    "廿六",
    "廿七",
    "廿八",
    "廿九",
    "三十",
    "卅一",
  ];
  return `${months[month - 1]}${days[day - 1]}`;
};

const formatAnniversaryDate = (anniversary: Anniversary) => {
  const parsed = parseMonthDay(anniversary.date);
  if (!parsed) return "无效日期";

  const { month, day } = parsed;

  return anniversary.calendar === "lunar"
    ? formatLunarDate(month, day)
    : `${month}月${day.toString().padStart(2, "0")}日`;
};

const isPastAnniversary = (anniversary: Anniversary) => {
  if (!normalizedAnniversaries.value.length) return false;

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const currentYear = today.getFullYear();
  const anniversaryInCurrentYear = toSolarDate(anniversary, currentYear);
  if (!anniversaryInCurrentYear) return false;

  const isPastThisYear = anniversaryInCurrentYear < today;
  return isPastThisYear;
};

const showActionSheet = ref(false);
const selectedAnniversary = ref<Anniversary | null>(null);

const {
  onPointerDown,
  onPointerUp,
  onPointerLeave,
  onPointerCancel,
  onPointerMove,
} = useLongPress({
  duration: 500,
  onFinish: () => {
    if (selectedAnniversary.value) {
      showActionSheet.value = true;
    }
  },
});

const handlePointerDown = (anniversary: Anniversary, event: PointerEvent) => {
  selectedAnniversary.value = anniversary;
  onPointerDown(event);
};

const actionSheetActions = computed<ActionSheetAction[]>(() => [
  {
    label: "编辑",
    handler: () => openEditDialog(selectedAnniversary.value),
  },
  {
    label: "删除",
    destructive: true,
    handler: () => openDeleteDialog(selectedAnniversary.value),
  },
]);

const showEditDialog = ref(false);
const editingAnniversary = ref<Anniversary | null>(null);
const savingAnniversary = ref(false);

const openAddDialog = () => {
  editingAnniversary.value = null;
  showEditDialog.value = true;
};

const openEditDialog = (anniversary: Anniversary | null) => {
  editingAnniversary.value = anniversary;
  showEditDialog.value = true;
};

const handleSaveAnniversary = async (anniversary: Anniversary) => {
  savingAnniversary.value = true;
  try {
    if (anniversary.id) {
      await anniversaryApi.updateAnniversary(anniversary.id, anniversary);
      showToast("纪念日更新成功", "success");
    } else {
      await anniversaryApi.createAnniversary(anniversary);
      showToast("纪念日添加成功", "success");
    }
    showEditDialog.value = false;
    await fetchAnniversaries();
  } catch {
    showToast("操作失败", "error");
  } finally {
    savingAnniversary.value = false;
  }
};

const showDeleteDialog = ref(false);
const deletingAnniversary = ref<Anniversary | null>(null);
const deleting = ref(false);

const openDeleteDialog = (anniversary: Anniversary | null) => {
  deletingAnniversary.value = anniversary;
  showDeleteDialog.value = true;
};

const handleDeleteAnniversary = async () => {
  if (!deletingAnniversary.value) return;
  deleting.value = true;
  try {
    await anniversaryApi.deleteAnniversary(deletingAnniversary.value.id);
    showToast("纪念日删除成功", "success");
    showDeleteDialog.value = false;
    await fetchAnniversaries();
  } catch {
    showToast("删除失败", "error");
  } finally {
    deleting.value = false;
  }
};

onMounted(async () => {
  uiStore.setLoading(true);
  await systemStore.fetchSystemInfo();
  await fetchAnniversaries();
  uiStore.setLoading(false);
});

const handleRefresh = async () => {
  await fetchAnniversaries();
  isRefreshing.value = false;
};

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target) return;
  isAtTop.value = target.scrollTop === 0;
};
</script>

<template>
  <MainLayout
    title="专属纪念"
    subtitle="珍藏我们的浪漫约定"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="anniversaries.length === 0"
  >
    <template #empty-state>
      <BaseIcon
        name="anniversary"
        size="w-24"
        style="color: var(--fe-text-secondary)"
      />
      <p class="font-bold text-xl mt-4 text-[var(--fe-text-primary)]">
        暂无纪念日数据
      </p>
      <p class="text-md mt-2 text-[var(--fe-text-secondary)]">
        期待添加第一个纪念日
      </p>
    </template>

    <template #main-content>
      <div class="flex flex-col h-full bg-[var(--fe-bg-gray)]/30">
        <div
          v-if="upcomingAnniversaryInfo"
          class="p-6 flex justify-center items-center text-center flex-shrink-0"
        >
          <div
            class="glass-thick rounded-[var(--fe-radius-card)] p-6 w-full max-w-lg border border-white/40 shadow-lg relative overflow-hidden group"
          >
            <div
              class="absolute -right-4 -top-4 opacity-10 transition-transform duration-700"
            >
              <BaseIcon
                :name="upcomingAnniversaryInfo.isToday ? 'calendar' : 'clock'"
                size="w-32 h-32"
              />
            </div>

            <div class="relative z-10 flex flex-col items-center">
              <span
                class="text-[10px] font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-1 opacity-60"
              >
                {{ upcomingAnniversaryInfo.isToday ? "işte 今天" : "Upcoming" }}
              </span>
              <div
                class="text-xl md:text-2xl font-bold text-[var(--fe-text-primary)]"
              >
                {{ upcomingAnniversaryInfo.message }}
              </div>
            </div>
          </div>
        </div>

        <van-pull-refresh
          v-model="isRefreshing"
          :disabled="!isAtTop"
          @refresh="handleRefresh"
          class="flex-grow"
        >
          <div
            ref="scrollContainer"
            class="overflow-y-auto custom-scrollbar px-4 md:px-8 pb-8"
            @scroll="handleScroll"
          >
            <div class="relative max-w-2xl mx-auto py-4">
              <div
                class="absolute left-6 top-0 bottom-0 w-px bg-black/10"
              ></div>

              <div class="space-y-6">
                <div
                  v-for="anniversary in sortedAnniversaries"
                  :key="anniversary.id"
                  class="relative flex items-center group"
                >
                  <div
                    class="absolute left-6 w-3 h-3 rounded-full border-2 border-white transform translate-x-[-50%] z-10 ios-transition"
                    :class="[
                      upcomingAnniversaryInfo &&
                      anniversary.id === upcomingAnniversaryInfo.anniversary.id
                        ? 'bg-[var(--fe-primary)] scale-125 shadow-[0_0_8px_rgba(240,173,160,0.6)]'
                        : 'bg-[var(--fe-text-secondary)] opacity-30',
                    ]"
                  ></div>

                  <div class="ml-12 flex-grow">
                    <div
                      class="glass-thick p-4 rounded-2xl border border-white/40 shadow-sm ios-transition tap-feedback active:scale-[0.98] cursor-pointer"
                      :class="{
                        'opacity-50 grayscale-[0.2]':
                          isPastAnniversary(anniversary),
                      }"
                      @pointerdown="handlePointerDown(anniversary, $event)"
                      @pointermove="onPointerMove"
                      @pointerup="onPointerUp"
                      @pointerleave="onPointerLeave"
                      @pointercancel="onPointerCancel"
                    >
                      <div class="flex justify-between items-start">
                        <div class="min-w-0">
                          <h3
                            class="font-bold text-[var(--fe-text-primary)] text-lg truncate"
                          >
                            {{ anniversary.title }}
                          </h3>
                          <p
                            class="text-xs font-medium text-[var(--fe-text-secondary)] mt-0.5"
                          >
                            {{ formatAnniversaryDate(anniversary) }}
                            <span
                              v-if="anniversary.calendar === 'lunar'"
                              class="ml-1 text-[var(--fe-primary-dark)]"
                            >
                              (农历)
                            </span>
                          </p>
                        </div>
                        <div
                          v-if="
                            upcomingAnniversaryInfo &&
                            anniversary.id ===
                              upcomingAnniversaryInfo.anniversary.id
                          "
                          class="flex-shrink-0"
                        >
                          <div
                            class="px-2 py-0.5 rounded-full bg-[var(--fe-primary)]/20 text-[var(--fe-primary)] text-[10px] font-bold uppercase tracking-wider"
                          >
                            Next
                          </div>
                        </div>
                      </div>
                      <p
                        class="mt-2 text-sm text-[var(--fe-text-primary)] leading-relaxed"
                      >
                        {{ anniversary.description }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="h-20 md:hidden"></div>
          </div>
        </van-pull-refresh>
      </div>
    </template>
  </MainLayout>

  <FloatingAddButton :loading="savingAnniversary" @click="openAddDialog" />

  <ActionSheet
    v-model="showActionSheet"
    title="纪念日操作"
    :actions="actionSheetActions"
  />

  <AnniversaryEditDialog
    v-model:open="showEditDialog"
    :anniversary="editingAnniversary"
    :loading="savingAnniversary"
    @confirm="handleSaveAnniversary"
  />

  <DeleteConfirmDialog
    v-model:open="showDeleteDialog"
    :loading="deleting"
    title="删除纪念日"
    :message="`确定要删除「${deletingAnniversary?.title || ''}」吗？删除后无法恢复。`"
    @confirm="handleDeleteAnniversary"
  />
</template>
