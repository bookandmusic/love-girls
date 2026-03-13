<template>
  <div
    class="flex items-center justify-center min-h-screen p-4 frontend-root"
    :style="{
      backgroundImage: `url(${bgSrc})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }"
  >
    <!-- 背景遮罩 -->
    <div class="absolute inset-0 bg-white/10 pointer-events-none"></div>

    <div
      class="glass-thick p-8 w-full max-w-md rounded-[var(--fe-radius-card)] border border-white/40 shadow-2xl relative z-10"
    >
      <div class="text-center mb-10">
        <h1 class="text-3xl font-bold text-[var(--fe-text-primary)]">登录</h1>
        <p class="text-[var(--fe-text-secondary)] mt-2 font-medium">请登录以管理您的纪念空间</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label
            for="username"
            class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
            >用户名</label
          >
          <input
            id="username"
            v-model="username"
            type="text"
            class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            placeholder="请输入用户名"
            required
          />
        </div>

        <div>
          <label
            for="password"
            class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
            >密码</label
          >
          <input
            id="password"
            v-model="password"
            type="password"
            class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="pt-4">
          <button
            type="submit"
            class="w-full bg-[var(--fe-primary)] text-white py-3 rounded-xl text-sm font-bold shadow-md shadow-[var(--fe-primary)]/20 ios-transition tap-feedback disabled:opacity-50"
            :disabled="uiStore.loading"
          >
            <span v-if="!uiStore.loading">登录</span>
            <span v-else>登录中...</span>
          </button>
        </div>
      </form>

      <!-- 返回首页 -->
      <div class="mt-8 text-center">
        <RouterLink
          to="/"
          class="text-xs font-bold text-[var(--fe-primary)] uppercase tracking-widest"
        >
          返回首页
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/assets/frontend-theme.css'

import { ref } from 'vue'
import { useRouter } from 'vue-router'

import bgSrc from '@/assets/images/bg.png'
import { userApi } from '@/services/userApi'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

const router = useRouter()
const uiStore = useUIStore()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const showToast = useToast()

const handleLogin = async () => {
  errorMessage.value = ''
  uiStore.setLoading(true)
  try {
    const response = await userApi.login({
      username: username.value,
      password: password.value,
    })
    if (response && response.access_token) {
      // 存储token到localStorage
      const token = response.access_token
      // 验证token并获取用户信息
      const userInfoResponse = await userApi.verifyToken(token)
      if (userInfoResponse && userInfoResponse.code === 0) {
        authStore.login(token, userInfoResponse.data)
        showToast('登录成功！', 'success')
        // 跳转到管理面板首页
        setTimeout(() => {
          router.push('/admin')
        }, 300)
      } else {
        errorMessage.value = '获取用户信息失败'
      }
    } else {
      errorMessage.value = '登录失败，请检查用户名和密码'
    }
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } }
    errorMessage.value = err.response?.data?.message || '登录失败，请稍后重试'
  } finally {
    uiStore.setLoading(false)
  }
  if (!errorMessage.value) {
    return
  }
  showToast(errorMessage.value, 'error')
}
</script>
