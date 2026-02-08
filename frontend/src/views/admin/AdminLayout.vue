<template>
  <div class="admin-layout flex flex-col h-screen p-2">
    <!-- 顶部导航栏 -->
    <header class="flex-shrink-0 h-18 pb-2">
      <div
        class="generic-card mx-auto p-4 sm:px-6 lg:px-8 h-full flex items-center justify-between"
      >
        <div class="flex items-center">
          <h1 class="text-xl font-(family-name:--font-signature) text-gray-900">后台管理</h1>
        </div>

        <div class="flex items-center space-x-4">
          <!-- 移动端菜单按钮 -->
          <button @click="showMobileMenu = true" class="lg:hidden p-1 rounded-full text-[#FF7500]">
            <BaseIcon :name="showMobileMenu ? 'menu-right' : 'menu-left'" size="w-6 h-6" />
          </button>

          <!-- 用户头像/用户名按钮（所有设备） -->
          <div class="relative">
            <button
              class="flex text-sm rounded-full focus:outline-none"
              @click="showConfirmDialog = true"
            >
              <span
                class="inline-flex items-center justify-center h-8 w-8 rounded-full bg-indigo-100"
              >
                <span class="text-gray-700 font-medium text-sm leading-none">
                  {{ authStore.userInfo?.userName?.charAt(0).toUpperCase() || 'U' }}
                </span>
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- 退出确认对话框 -->
      <GenericDialog
        :open="showConfirmDialog"
        title="退出登录"
        size-class="max-w-md"
        @update:open="showConfirmDialog = $event"
        @cancel="showConfirmDialog = false"
      >
        <template #content>
          <p class="text-gray-700">
            确定要退出登录吗，{{ authStore.userInfo?.userName || '用户' }}？
          </p>
        </template>
        <template #actions>
          <div class="w-full flex">
            <div class="flex-1 text-center cursor-pointer" @click="showConfirmDialog = false">
              取消
            </div>
            <div
              class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
              @click="handleLogout"
            >
              退出
            </div>
          </div>
        </template>
      </GenericDialog>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边栏 -->
      <aside class="w-64 flex-shrink-0 hidden lg:block pr-2">
        <div class="generic-card p-4 h-full w-full">
          <nav class="mt-5 px-2">
            <div class="space-y-1">
              <router-link
                v-for="item in menuItems"
                :key="item.path"
                :to="item.path"
                class="text-[var(--primary-color)] hover:bg-[var(--primary-light)] group flex items-center px-2 py-2 text-md font-medium rounded-md transition-all duration-200"
                :class="route.path === item.path ? 'bg-[var(--primary-color)] text-white' : ''"
              >
                <component :is="item.icon" :size="24" class="mr-3" />
                {{ item.label }}
              </router-link>
            </div>
          </nav>
        </div>
      </aside>

      <!-- 主内容区域，移除标题和按钮 -->
      <main class="flex-1 flex flex-col overflow-hidden">
        <div class="flex-1 h-full w-full generic-card p-4 overflow-hidden">
          <router-view />
        </div>
      </main>
    </div>

    <!-- 移动端抽屉菜单 - 从左侧弹出 -->
    <div v-if="showMobileMenu" class="fixed inset-0 z-50 lg:hidden" @click="showMobileMenu = false">
      <!-- 背景遮罩 -->
      <div class="fixed inset-0 bg-white/50 bg-opacity-80"></div>

      <!-- 抽屉内容 - 从左侧滑入 -->
      <div class="fixed top-0 left-0 h-full w-64 max-w-sm bg-[#EEEDEE]" @click.stop>
        <div class="p-4">
          <div class="flex justify-end mb-4">
            <button @click="showMobileMenu = false" class="text-[var(--primary-color)]">
              <BaseIcon name="close" size="w-6 h-6" color="text-[#FFB61E]" />
            </button>
          </div>

          <nav class="mt-5">
            <div class="space-y-1">
              <router-link
                v-for="item in menuItems"
                :key="item.path"
                :to="item.path"
                class="text-[var(--primary-color)] group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-all duration-200"
                :class="route.path === item.path ? 'bg-[var(--primary-color)] text-white' : ''"
                @click="showMobileMenu = false"
              >
                <BaseIcon :name="item.icon" class="mr-2" size="w-6 h-6" />
                {{ item.label }}
              </router-link>
            </div>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import GenericDialog from '@/components/ui/GenericDialog.vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/utils/toastUtils'

const router = useRouter()
const authStore = useAuthStore()

const showConfirmDialog = ref(false)

const showToast = useToast()

const handleLogout = () => {
  authStore.logout()
  showToast('已成功登出', 'info')
  setTimeout(() => {
    router.push('/')
  }, 300)
}

const route = useRoute()
const showMobileMenu = ref(false)

// 定义菜单项数组
const menuItems = [
  {
    path: '/admin/dashboard',
    label: '仪表盘',
    icon: 'home-heart',
  },
  {
    path: '/admin/users',
    label: '用户管理',
    icon: 'user-heart',
  },
  {
    path: '/admin/content',
    label: '内容管理',
    icon: 'content-heart',
  },
  {
    path: '/admin/settings',
    label: '系统设置',
    icon: 'setting-heart',
  },
]
</script>
