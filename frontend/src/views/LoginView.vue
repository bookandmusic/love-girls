<template>
  <div class="flex items-center justify-center min-h-screen p-4">
    <AnimatedBorderCard
      :borderColor="'#f0ada0'"
      :borderWidth="3"
      class="generic-card p-4 w-full max-w-md"
    >
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-800">登录</h1>
        <p class="text-gray-600 mt-2">请登录以继续访问管理面板</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            class="w-full win11-input"
            placeholder="请输入用户名"
            required
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="w-full win11-input"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="flex items-center justify-center">
          <button
            type="submit"
            class="w-full win11-button disabled:opacity-50"
            :disabled="uiStore.loading"
          >
            <span v-if="!uiStore.loading">登录</span>
            <span v-else>登录中...</span>
          </button>
        </div>
      </form>
    </AnimatedBorderCard>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import AnimatedBorderCard from '@/components/ui/AnimatedBorderCard.vue'
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
