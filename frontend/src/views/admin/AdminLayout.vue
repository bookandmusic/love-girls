<template>
  <div class="admin-layout flex flex-col h-screen">
    <!-- 顶部导航栏 -->
    <header class="admin-header flex-shrink-0 p-2 pb-2">
      <div class="admin-card mx-auto p-4 sm:px-6 lg:px-8 h-full flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-3xl">💕</span>
          <h1 class="text-3xl font-bold font-(family-name:--font-signature) admin-text-primary">
            后台管理
          </h1>
        </div>

        <div class="flex items-center space-x-4">
          <!-- 移动端菜单按钮 -->
          <button
            @click="showMobileMenu = true"
            class="lg:hidden p-1 rounded-full hover:bg-white/40 transition"
          >
            <BaseIcon
              :name="showMobileMenu ? 'menu-right' : 'menu-left'"
              size="w-6 h-6"
              color="#f0ada0"
            />
          </button>

          <!-- 用户头像/用户名按钮（所有设备） -->
          <div class="relative">
            <button
              class="flex text-sm rounded-full focus:outline-none"
              @click="showConfirmDialog = true"
            >
              <span
                class="inline-flex items-center justify-center h-10 w-10 rounded-full bg-gradient-to-br from-[#f0ada0] to-[#d89388] text-white font-bold"
              >
                {{ authStore.userInfo?.userName?.charAt(0).toUpperCase() || 'U' }}
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
          <p class="admin-text-primary">
            确定要退出登录吗，{{ authStore.userInfo?.userName || '用户' }}？
          </p>
        </template>
        <template #actions>
          <div class="w-full flex">
            <div
              class="flex-1 text-center cursor-pointer admin-text-secondary hover:admin-text-primary transition"
              @click="showConfirmDialog = false"
            >
              取消
            </div>
            <div
              class="w-1/2 border-l border-white/60 text-center cursor-pointer text-[#E8A8A8] hover:text-[#d89898] transition"
              @click="handleLogout"
            >
              退出
            </div>
          </div>
        </template>
      </GenericDialog>
    </header>

    <div class="flex flex-1 overflow-hidden px-2">
      <!-- 侧边栏 -->
      <aside class="w-64 flex-shrink-0 hidden lg:block pr-2">
        <div class="admin-sidebar admin-card p-4 h-full w-full flex flex-col">
          <nav class="mt-5 px-2 flex-1">
            <div class="space-y-2">
              <router-link
                v-for="item in menuItems"
                :key="item.path"
                :to="item.path"
                class="admin-sidebar-item group flex items-center px-3 py-3 text-md font-medium rounded-xl transition-all duration-200"
                :class="isActiveMenuItem(item.path) ? 'active' : ''"
              >
                <BaseIcon
                  :name="item.icon"
                  size="w-6 h-6"
                  :color="isActiveMenuItem(item.path) ? '#ffffff' : '#f0ada0'"
                  class="mr-3"
                />
                {{ item.label }}
              </router-link>
            </div>
          </nav>
          <div class="px-2 pb-4">
            <div class="border-t border-white/30 pt-4">
              <router-link
                to="/"
                class="flex items-center px-3 py-3 rounded-xl transition-all duration-200 bg-gradient-to-r from-[#f0ada0]/20 to-transparent hover:from-[#f0ada0]/30"
              >
                <div
                  class="w-10 h-10 rounded-full bg-gradient-to-br from-[#f0ada0] to-[#d89388] flex items-center justify-center mr-3"
                >
                  <BaseIcon name="home" size="w-5 h-5" color="#ffffff" />
                </div>
                <div class="flex-1">
                  <div class="text-sm font-semibold admin-text-primary">返回前台</div>
                  <div class="text-xs admin-text-secondary">浏览网站首页</div>
                </div>
              </router-link>
            </div>
          </div>
        </div>
      </aside>

      <!-- 主内容区域 -->
      <main class="flex-1 flex flex-col overflow-hidden">
        <div class="admin-card flex-1 h-full w-full p-4 overflow-hidden">
          <router-view />
        </div>
      </main>
    </div>

    <!-- 移动端抽屉菜单 - 从左侧弹出 -->
    <div v-if="showMobileMenu" class="fixed inset-0 z-50 lg:hidden" @click="showMobileMenu = false">
      <div class="fixed inset-0 bg-black/20 backdrop-blur-sm"></div>

      <div class="fixed top-0 left-0 h-full w-64 max-w-sm admin-layout" @click.stop>
        <div class="p-4 h-full flex flex-col">
          <div class="flex justify-between items-center mb-6">
            <div class="flex items-center gap-2">
              <span class="text-2xl">💕</span>
              <h2 class="text-lg font-(family-name:--font-signature) admin-text-primary">菜单</h2>
            </div>
            <button
              @click="showMobileMenu = false"
              class="p-2 rounded-full hover:bg-white/40 transition"
            >
              <BaseIcon name="close" size="w-5 h-5" color="#f0ada0" />
            </button>
          </div>

          <nav class="flex-1">
            <div class="space-y-2">
              <router-link
                v-for="item in menuItems"
                :key="item.path"
                :to="item.path"
                class="admin-sidebar-item group flex items-center px-3 py-3 text-md font-medium rounded-xl transition-all duration-200"
                :class="isActiveMenuItem(item.path) ? 'active' : ''"
                @click="showMobileMenu = false"
              >
                <BaseIcon
                  :name="item.icon"
                  size="w-6 h-6"
                  :color="isActiveMenuItem(item.path) ? '#ffffff' : '#f0ada0'"
                  class="mr-3"
                />
                {{ item.label }}
              </router-link>
            </div>
          </nav>
          <div class="pb-4">
            <div class="border-t border-white/30 pt-4">
              <router-link
                to="/"
                class="flex items-center px-3 py-3 rounded-xl transition-all duration-200 bg-gradient-to-r from-[#f0ada0]/20 to-transparent hover:from-[#f0ada0]/30"
                @click="showMobileMenu = false"
              >
                <div
                  class="w-10 h-10 rounded-full bg-gradient-to-br from-[#f0ada0] to-[#d89388] flex items-center justify-center mr-3"
                >
                  <BaseIcon name="home" size="w-5 h-5" color="#ffffff" />
                </div>
                <div class="flex-1">
                  <div class="text-sm font-semibold admin-text-primary">返回前台</div>
                  <div class="text-xs admin-text-secondary">浏览网站首页</div>
                </div>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/views/admin/styles/admin-theme.css'

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

const isActiveMenuItem = (itemPath: string) => {
  if (itemPath === '/admin/content') {
    return route.path.startsWith('/admin/content')
  }
  return route.path === itemPath
}

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
