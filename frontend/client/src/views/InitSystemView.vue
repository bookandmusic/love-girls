<template>
  <div
    class="h-screen w-full flex items-center justify-center p-4 frontend-root overflow-hidden"
    :style="{
      backgroundImage: `url(${bgSrc})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }"
  >
    <!-- 背景遮罩 -->
    <div class="absolute inset-0 bg-white/10 pointer-events-none"></div>

    <div
      class="w-full max-w-xl glass-thick rounded-[var(--fe-radius-card)] border border-white/40 shadow-2xl relative z-10 flex flex-col max-h-[90vh] overflow-hidden"
    >
      <!-- Header -->
      <div class="p-8 pb-4 text-center">
        <h1
          class="text-2xl md:text-3xl font-bold text-[var(--fe-text-primary)] mb-2"
        >
          欢迎使用情侣纪念站点
        </h1>
        <p class="text-[var(--fe-text-secondary)] font-medium">
          初始化您的专属纪念空间
        </p>
      </div>

      <!-- Step Indicator -->
      <div class="px-8 mb-4">
        <div
          class="flex items-center justify-between p-1 glass-ultra-thin rounded-full border border-white/40"
        >
          <div
            v-for="i in totalSteps"
            :key="i"
            class="flex-1 text-[10px] font-bold py-1.5 text-center rounded-full ios-transition"
            :class="[
              i === currentStep
                ? 'bg-white shadow-sm text-[var(--fe-text-primary)]'
                : 'text-[var(--fe-text-secondary)]',
            ]"
          >
            {{ i }} / {{ totalSteps }}
          </div>
        </div>
      </div>

      <!-- Scrollable form container -->
      <div class="flex-grow overflow-y-auto px-8 py-4 custom-scrollbar">
        <!-- Step 1: Site -->
        <div v-if="currentStep === 1" class="space-y-6">
          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >站点名称</label
            >
            <input
              v-model="form.siteName"
              placeholder="例如：鹿与星的纪念站"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >站点描述（可选）</label
            >
            <input
              v-model="form.siteDescription"
              placeholder="例如：记录我们的美好时光"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >故事开始的日期</label
            >
            <input
              v-model="form.startDate"
              type="date"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>
        </div>

        <!-- Step 2: User A -->
        <div v-if="currentStep === 2" class="space-y-6">
          <div class="flex justify-center mb-4">
            <div
              class="w-20 h-20 rounded-2xl glass-ultra-thin border-2 border-white/60 flex items-center justify-center text-2xl text-[var(--fe-primary)] font-bold shadow-sm"
            >
              <img
                v-if="avatarAPreview"
                :src="avatarAPreview"
                class="w-full h-full object-cover rounded-2xl"
                draggable="false"
              />
              <span v-else>{{ form.userAName?.[0] || "A" }}</span>
            </div>
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >昵称</label
            >
            <input
              v-model="form.userAName"
              placeholder="昵称"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >角色</label
            >
            <select
              v-model="form.userARole"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition appearance-none"
            >
              <option value="boy">男生</option>
              <option value="girl">女生</option>
            </select>
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >邮箱（可选）</label
            >
            <input
              v-model="form.userAEmail"
              type="email"
              placeholder="邮箱"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>
        </div>

        <!-- Step 3: User B -->
        <div v-if="currentStep === 3" class="space-y-6">
          <div class="flex justify-center mb-4">
            <div
              class="w-20 h-20 rounded-2xl glass-ultra-thin border-2 border-white/60 flex items-center justify-center text-2xl text-[var(--fe-primary)] font-bold shadow-sm"
            >
              <img
                v-if="avatarBPreview"
                :src="avatarBPreview"
                class="w-full h-full object-cover rounded-2xl"
                draggable="false"
              />
              <span v-else>{{ form.userBName?.[0] || "B" }}</span>
            </div>
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >昵称</label
            >
            <input
              v-model="form.userBName"
              placeholder="昵称"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >角色</label
            >
            <select
              v-model="form.userBRole"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition appearance-none"
            >
              <option value="boy">男生</option>
              <option value="girl">女生</option>
            </select>
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >邮箱（可选）</label
            >
            <input
              v-model="form.userBEmail"
              type="email"
              placeholder="邮箱"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>
        </div>

        <!-- Step 4: Password -->
        <div v-if="currentStep === 4" class="space-y-6">
          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >站点访问密码</label
            >
            <input
              v-model="form.sitePassword"
              type="password"
              placeholder="建议使用强密码"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>

          <div>
            <label
              class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
              >确认密码</label
            >
            <input
              v-model="form.sitePasswordConfirm"
              type="password"
              placeholder="请再次输入密码"
              class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            />
          </div>
        </div>

        <!-- Error Message Display -->
        <div class="mt-4 min-h-[24px]">
          <Transition name="fade">
            <div
              v-if="getCurrentError"
              class="text-xs font-bold text-red-500 bg-red-50/50 p-2 rounded-lg border border-red-200/50"
            >
              {{ getCurrentError }}
            </div>
          </Transition>
        </div>
      </div>

      <!-- Actions -->
      <div class="p-8 pt-4 flex space-x-4">
        <button
          v-if="currentStep > 1"
          @click="prevStep"
          class="flex-1 py-3 rounded-xl text-sm font-bold text-[var(--fe-text-secondary)] ios-transition tap-feedback"
        >
          上一步
        </button>

        <button
          type="button"
          @click.prevent="nextStep"
          class="flex-[2] bg-[var(--fe-primary)] text-white py-3 rounded-xl text-sm font-bold shadow-md shadow-[var(--fe-primary)]/20 ios-transition tap-feedback disabled:opacity-50"
          :disabled="uiStore.loading"
        >
          <span v-if="!uiStore.loading">{{
            currentStep === totalSteps ? "完成初始化" : "下一步"
          }}</span>
          <span v-else>处理中...</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import "@/assets/frontend-theme.css";

import { computed, reactive, ref, watch } from "vue";

import bgSrc from "@/assets/images/bg.png";
import router from "@/router";
import { systemApi } from "@/services/system";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

const currentStep = ref(1);
const totalSteps = 4;
const uiStore = useUIStore();

const avatarAPreview = ref<string | null>(null);
const avatarBPreview = ref<string | null>(null);

const form = reactive({
  siteName: "",
  siteDescription: "",
  startDate: "",
  userAName: "",
  userARole: "boy" as "boy" | "girl",
  userAEmail: "",
  userAPhone: "",
  userBName: "",
  userBRole: "girl" as "boy" | "girl",
  userBEmail: "",
  userBPhone: "",
  sitePassword: "",
  sitePasswordConfirm: "",
});

const showToast = useToast();

const errors = reactive({
  siteName: "",
  startDate: "",
  userAName: "",
  userARole: "",
  userBName: "",
  userBRole: "",
  sitePassword: "",
  sitePasswordConfirm: "",
  passwordMismatch: "",
});

// 监听用户A角色变化，自动切换用户B角色
watch(
  () => form.userARole,
  (newRole: "boy" | "girl") => {
    form.userBRole = newRole === "boy" ? "girl" : "boy";
  },
);

// 计算当前步骤的错误信息
const getCurrentError = computed(() => {
  switch (currentStep.value) {
    case 1:
      return errors.siteName || errors.startDate;
    case 2:
      return errors.userAName || errors.userARole;
    case 3:
      return errors.userBName || errors.userBRole;
    case 4:
      return (
        errors.sitePassword ||
        errors.sitePasswordConfirm ||
        errors.passwordMismatch
      );
    default:
      return "";
  }
});

function clearErrors() {
  const errorKeys = Object.keys(errors) as Array<keyof typeof errors>;
  errorKeys.forEach((key) => {
    errors[key] = "";
  });
}

function validateStep(): boolean {
  clearErrors();

  if (currentStep.value === 1) {
    if (!form.siteName) {
      errors.siteName = "请输入站点名称";
      return false;
    }
    if (!form.startDate) {
      errors.startDate = "请选择在一起日期";
      return false;
    }
  }

  if (currentStep.value === 2) {
    if (!form.userAName) {
      errors.userAName = "请输入用户 A 昵称";
      return false;
    }
    if (!form.userARole) {
      errors.userARole = "请选择用户 A 角色";
      return false;
    }
  }

  if (currentStep.value === 3) {
    if (!form.userBName) {
      errors.userBName = "请输入用户 B 昵称";
      return false;
    }
    if (!form.userBRole) {
      errors.userBRole = "请选择用户 B 角色";
      return false;
    }
  }

  if (currentStep.value === 4) {
    if (!form.sitePassword) {
      errors.sitePassword = "请输入站点访问密码";
      return false;
    }
    if (!form.sitePasswordConfirm) {
      errors.sitePasswordConfirm = "请确认站点访问密码";
      return false;
    }
    if (form.sitePassword !== form.sitePasswordConfirm) {
      errors.passwordMismatch = "两次输入的密码不一致";
      return false;
    }
  }

  return true;
}

function nextStep() {
  if (!validateStep()) return;
  if (currentStep.value < totalSteps) currentStep.value++;
  else submit();
}

function prevStep() {
  if (currentStep.value > 1) currentStep.value--;
}

function submit() {
  uiStore.setLoading(true);

  systemApi
    .initSystem(form)
    .then((res) => {
      if (res.data.code === 0) {
        // 初始化成功后，立即更新系统状态
        const systemStore = useSystemStore();
        systemStore.setInitialized(true);

        showToast("系统初始化成功！", "success");
        setTimeout(() => router.push("/"), 300);
      } else {
        showToast(res.data.msg || "初始化失败", "error");
      }
    })
    .catch((err) => {
      console.error(err);
      showToast("网络错误，请稍后重试", "error");
    })
    .finally(() => {
      uiStore.setLoading(false);
    });
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
