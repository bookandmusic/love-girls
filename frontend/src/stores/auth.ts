import { defineStore } from 'pinia'
import { ref } from 'vue'

import { userApi, type UserInfo } from '@/services/userApi'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('auth_token'))
  const userInfo = ref<UserInfo | null>(null)
  const isAuthenticated = ref(false)

  // 登录
  const login = (newToken: string, userData: UserInfo) => {
    token.value = newToken
    userInfo.value = userData
    isAuthenticated.value = true
    localStorage.setItem('auth_token', newToken)
  }

  // 登出
  const logout = () => {
    token.value = null
    userInfo.value = null
    isAuthenticated.value = false
    localStorage.removeItem('auth_token')
  }

  // 检查登录状态
  const checkAuthStatus = async () => {
    const storedToken = localStorage.getItem('auth_token')
    if (!storedToken) {
      isAuthenticated.value = false
      return false
    }

    // 验证token是否有效
    try {
      const response = await userApi.verifyToken(storedToken)

      if (response && response.code === 0) {
        token.value = storedToken
        userInfo.value = response.data
        isAuthenticated.value = true
        return true
      } else {
        // Token无效，清除本地存储
        logout()
        return false
      }
    } catch (error) {
      console.error('验证token时发生错误:', (error as Error).message)
      logout()
      return false
    }
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    login,
    logout,
    checkAuthStatus,
  }
})
