import axios from 'axios'

import { useAuthStore } from '@/stores/auth'

// 创建一个标准的axios实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1', // 基础URL
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 从localStorage获取token并添加到请求头
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    // 可以在这里统一处理响应数据
    return response
  },
  error => {
    // 可以在这里统一处理错误
    console.error('API Error:', error)

    // 如果是401错误，可能是token过期，跳转到登录页
    if (error.response && error.response.status === 401) {
      // 清除本地token
      const authStore = useAuthStore()
      authStore.logout()

      // 如果当前不在登录页面，则跳转到登录页
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  }
)

export default api
