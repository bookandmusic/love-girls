<template>
  <div
    class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden"
  >
    <div
      class="absolute inset-0 bg-gradient-to-br from-[#fce4ec] via-[#f8bbd9]/30 to-[#f48fb1]/20"
    ></div>
    <div
      class="absolute top-0 right-0 w-96 h-96 bg-[#f0ada0]/40 rounded-full blur-3xl transform translate-x-1/2 -translate-y-1/2"
    ></div>
    <div
      class="absolute bottom-0 left-0 w-[500px] h-[500px] bg-[#d89388]/30 rounded-full blur-3xl transform -translate-x-1/3 translate-y-1/3"
    ></div>

    <div class="w-full max-w-xl relative z-10">
      <div class="text-center mb-6">
        <img
          src="/favicon.png"
          alt="Logo"
          class="w-16 h-16 mx-auto mb-4 rounded-2xl shadow-lg"
        />
        <h1 class="text-2xl font-bold text-gray-800">欢迎使用情侣纪念站点</h1>
        <p class="text-gray-500 mt-2">初始化您的专属纪念空间</p>
      </div>

      <div
        class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden"
      >
        <div class="px-8 py-4 border-b border-gray-100">
          <div class="flex items-center justify-between">
            <div v-for="i in totalSteps" :key="i" class="flex items-center">
              <div
                class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-all"
                :class="[
                  i <= currentStep
                    ? 'bg-gradient-to-r from-[#f0ada0] to-[#d89388] text-white'
                    : 'bg-gray-100 text-gray-400',
                ]"
              >
                {{ i }}
              </div>
              <div
                v-if="i < totalSteps"
                class="w-12 h-1 mx-1"
                :class="
                  i < currentStep
                    ? 'bg-gradient-to-r from-[#f0ada0] to-[#d89388]'
                    : 'bg-gray-100'
                "
              ></div>
            </div>
          </div>
        </div>

        <div class="p-8 max-h-[60vh] overflow-y-auto">
          <Transition name="fade" mode="out-in">
            <div :key="currentStep">
              <div v-if="currentStep === 1" class="space-y-5">
                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    站点名称
                  </label>
                  <input
                    v-model="form.siteName"
                    placeholder="例如：鹿与星的纪念站"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    站点描述（可选）
                  </label>
                  <input
                    v-model="form.siteDescription"
                    placeholder="例如：记录我们的美好时光"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    故事开始的日期
                  </label>
                  <input
                    v-model="form.startDate"
                    type="date"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>
              </div>

              <div v-if="currentStep === 2" class="space-y-5">
                <div class="flex justify-center mb-4">
                  <div
                    class="w-20 h-20 rounded-2xl border-2 border-gray-200 flex items-center justify-center text-2xl font-bold bg-gray-50"
                  >
                    <span class="text-[#f0ada0]">{{
                      form.userAName?.[0] || "A"
                    }}</span>
                  </div>
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    昵称
                  </label>
                  <input
                    v-model="form.userAName"
                    placeholder="昵称"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    角色
                  </label>
                  <select
                    v-model="form.userARole"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all appearance-none"
                  >
                    <option value="boy">男生</option>
                    <option value="girl">女生</option>
                  </select>
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    邮箱（可选）
                  </label>
                  <input
                    v-model="form.userAEmail"
                    type="email"
                    placeholder="邮箱"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>
              </div>

              <div v-if="currentStep === 3" class="space-y-5">
                <div class="flex justify-center mb-4">
                  <div
                    class="w-20 h-20 rounded-2xl border-2 border-gray-200 flex items-center justify-center text-2xl font-bold bg-gray-50"
                  >
                    <span class="text-[#f0ada0]">{{
                      form.userBName?.[0] || "B"
                    }}</span>
                  </div>
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    昵称
                  </label>
                  <input
                    v-model="form.userBName"
                    placeholder="昵称"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    角色
                  </label>
                  <select
                    v-model="form.userBRole"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all appearance-none"
                  >
                    <option value="boy">男生</option>
                    <option value="girl">女生</option>
                  </select>
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    邮箱（可选）
                  </label>
                  <input
                    v-model="form.userBEmail"
                    type="email"
                    placeholder="邮箱"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>
              </div>

              <div v-if="currentStep === 4" class="space-y-5">
                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    站点访问密码
                  </label>
                  <input
                    v-model="form.sitePassword"
                    type="password"
                    placeholder="建议使用强密码"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>

                <div>
                  <label
                    class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                  >
                    确认密码
                  </label>
                  <input
                    v-model="form.sitePasswordConfirm"
                    type="password"
                    placeholder="请再次输入密码"
                    class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
                  />
                </div>
              </div>
            </div>
          </Transition>

          <div class="mt-4 min-h-[24px]">
            <Transition name="fade">
              <div
                v-if="getCurrentError"
                class="text-xs font-medium text-red-500 bg-red-50 p-3 rounded-xl border border-red-100"
              >
                {{ getCurrentError }}
              </div>
            </Transition>
          </div>
        </div>

        <div class="px-8 py-4 border-t border-gray-100 flex space-x-4">
          <button
            v-if="currentStep > 1"
            @click="prevStep"
            class="flex-1 py-3 rounded-xl text-sm font-bold text-gray-500 hover:bg-gray-50 transition-all"
          >
            上一步
          </button>

          <button
            type="button"
            @click.prevent="nextStep"
            class="flex-[2] py-3 rounded-xl text-sm font-bold text-white transition-all disabled:opacity-50 shadow-sm"
            :class="
              uiStore.loading
                ? 'bg-gray-300'
                : 'bg-gradient-to-r from-[#f0ada0] to-[#d89388] hover:from-[#d89388] hover:to-[#c78277]'
            "
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
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";

import { systemApi } from "@/services/system";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

const router = useRouter();
const currentStep = ref(1);
const totalSteps = 4;
const uiStore = useUIStore();
const systemStore = useSystemStore();

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

watch(
  () => form.userARole,
  (newRole: "boy" | "girl") => {
    form.userBRole = newRole === "boy" ? "girl" : "boy";
  },
);

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
        systemStore.setInitialized(true);
        showToast("系统初始化成功！", "success");
        setTimeout(() => router.push("/login"), 300);
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
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
