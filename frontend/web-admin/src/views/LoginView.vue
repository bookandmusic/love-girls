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

    <div class="w-full max-w-md relative z-10">
      <div class="text-center mb-8">
        <img
          src="/favicon.png"
          alt="Logo"
          class="w-16 h-16 mx-auto mb-4 rounded-2xl shadow-lg"
        />
        <h1 class="text-2xl font-bold text-gray-800">后台管理</h1>
        <p class="text-gray-500 mt-2">请登录以管理您的纪念空间</p>
      </div>

      <div class="bg-white rounded-2xl shadow-lg border border-gray-100 p-8">
        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label
              for="username"
              class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
            >
              用户名
            </label>
            <input
              id="username"
              v-model="username"
              type="text"
              class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div>
            <label
              for="password"
              class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
            >
              密码
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              class="w-full px-4 py-3 rounded-xl border border-gray-200 bg-gray-50/50 text-sm focus:border-[#f0ada0] focus:ring-2 focus:ring-[#f0ada0]/20 outline-none transition-all"
              placeholder="请输入密码"
              required
            />
          </div>

          <div class="pt-2">
            <button
              type="submit"
              class="w-full py-3 rounded-xl text-sm font-bold text-white transition-all disabled:opacity-50 shadow-sm"
              :class="
                uiStore.loading
                  ? 'bg-gray-300'
                  : 'bg-gradient-to-r from-[#f0ada0] to-[#d89388] hover:from-[#d89388] hover:to-[#c78277]'
              "
              :disabled="uiStore.loading"
            >
              <span v-if="!uiStore.loading">登录</span>
              <span v-else>登录中...</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

import { userApi } from "@/services/userApi";
import { useAuthStore } from "@/stores/auth";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

const router = useRouter();
const uiStore = useUIStore();
const authStore = useAuthStore();
const showToast = useToast();

const username = ref("");
const password = ref("");

const handleLogin = async () => {
  uiStore.setLoading(true);
  try {
    const response = await userApi.login({
      username: username.value,
      password: password.value,
    });
    if (response && response.access_token) {
      const token = response.access_token;
      const userInfoResponse = await userApi.verifyToken(token);
      if (userInfoResponse && userInfoResponse.code === 0) {
        authStore.login(token, userInfoResponse.data);
        showToast("登录成功！", "success");
        setTimeout(() => {
          router.push("/dashboard");
        }, 300);
      } else {
        showToast("获取用户信息失败", "error");
      }
    } else {
      showToast("登录失败，请检查用户名和密码", "error");
    }
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } };
    showToast(err.response?.data?.message || "登录失败，请稍后重试", "error");
  } finally {
    uiStore.setLoading(false);
  }
};
</script>
