<template>
  <GenericDialog
    :open="open"
    :title="anniversary?.id ? '编辑纪念日' : '添加纪念日'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <div class="space-y-4 h-full">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            标题 <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            class="w-full win11-input"
            placeholder="请输入纪念日标题"
            :disabled="loading"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            日历类型 <span class="text-red-500">*</span>
          </label>
          <div class="flex items-center bg-gray-100 rounded-lg p-1">
            <button
              type="button"
              @click="form.calendar = 'solar'"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="
                form.calendar === 'solar'
                  ? 'bg-white text-[var(--fe-primary)] shadow-sm'
                  : 'text-gray-500'
              "
              :disabled="loading"
            >
              公历
            </button>
            <button
              type="button"
              @click="form.calendar = 'lunar'"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="
                form.calendar === 'lunar'
                  ? 'bg-white text-[var(--fe-primary)] shadow-sm'
                  : 'text-gray-500'
              "
              :disabled="loading"
            >
              农历
            </button>
          </div>
        </div>

        <div class="relative">
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            日期 <span class="text-red-500">*</span>
          </label>
          <input
            ref="inputRef"
            v-model="form.date"
            type="text"
            class="w-full win11-input cursor-pointer"
            placeholder="请选择日期"
            readonly
            @focus="openCalendar"
            :disabled="loading"
          />
          <div
            v-if="isCalendarVisible"
            ref="calendarRef"
            class="absolute mt-1 z-10 w-80"
            @mousedown.prevent
          >
            <LunarCalendar
              :default-date="calendarDate"
              :show-lunar="true"
              @onSelect="handleDateChange"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            描述
          </label>
          <textarea
            v-model="form.description"
            class="w-full win11-input resize-none"
            rows="4"
            placeholder="请输入纪念日描述"
            :disabled="loading"
          ></textarea>
        </div>
      </div>
    </template>
    <template #actions>
      <button
        class="flex-1 py-3.5 text-center text-gray-700 font-medium hover:bg-gray-100 active:bg-gray-200 transition-colors"
        @click="closeDialog"
      >
        取消
      </button>
      <button
        class="flex-1 py-3.5 text-center text-[var(--fe-primary-dark)] font-semibold border-l border-gray-100 hover:bg-gray-100 active:bg-gray-200 transition-colors"
        @click="handleSave"
      >
        确认
      </button>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { getLunarDate, getSolarDateFromLunar } from "chinese-days";
import { computed, ref, watch } from "vue";

import GenericDialog from "@/components/ui/GenericDialog.vue";
import LunarCalendar from "@/components/ui/LunarCalendar.vue";
import { type Anniversary } from "@/services/anniversaryApi";
import { useToast } from "@/utils/toastUtils";

const showToast = useToast();

const inputRef = ref<HTMLInputElement | null>(null);

interface Props {
  open: boolean;
  loading?: boolean;
  anniversary?: Anniversary | null;
}

const props = withDefaults(defineProps<Props>(), {
  anniversary: null,
  loading: false,
});

interface Emits {
  (e: "update:open", open: boolean): void;
  (e: "confirm", anniversary: Anniversary): void;
  (e: "cancel"): void;
}

const emit = defineEmits<Emits>();

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
};

const DEFAULT_ANNIVERSARY: Anniversary = {
  id: 0,
  title: "",
  date: "",
  description: "",
  calendar: "solar" as "solar" | "lunar",
};

const form = ref<Anniversary>({ ...DEFAULT_ANNIVERSARY });

const isCalendarVisible = ref(false);

const calendarDate = computed(() => {
  switch (form.value.calendar) {
    case "lunar":
      const solarDate = getSolarDateFromLunar(form.value.date);
      return solarDate.date;
    default:
      return form.value.date;
  }
});

const openCalendar = () => {
  if (!props.loading) {
    isCalendarVisible.value = true;
  }
};

const handleDateChange = (selectedDate: string) => {
  let tmpDate = selectedDate;
  switch (form.value.calendar) {
    case "lunar":
      const lunarDate = getLunarDate(selectedDate);
      tmpDate = `${lunarDate.lunarYear}-${lunarDate.lunarMon}-${lunarDate.lunarDay}`;
    default:
      break;
  }
  form.value.date = tmpDate;
  isCalendarVisible.value = false;
  inputRef.value?.blur();
};

const handleSave = async () => {
  if (!form.value.title.trim()) {
    showToast("请输入标题", "error");
    return;
  }
  if (!form.value.date) {
    showToast("请选择日期", "error");
    return;
  }
  emit("confirm", form.value);
};

watch(
  () => props.anniversary,
  (newVal) => {
    if (newVal) {
      form.value = {
        id: newVal.id,
        title: newVal.title,
        date: newVal.date,
        description: newVal.description,
        calendar: newVal.calendar,
      };
    } else {
      const today = new Date();
      const formattedDate = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, "0")}-${String(today.getDate()).padStart(2, "0")}`;
      form.value = {
        ...DEFAULT_ANNIVERSARY,
        date: formattedDate,
      };
    }
  },
  { deep: true, immediate: true },
);
</script>
