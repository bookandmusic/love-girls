<template>
  <div class="admin-layout flex flex-col h-screen overflow-hidden">
    <!-- 顶部导航栏 - Ultra-thin 材质 -->
    <header class="admin-header flex-shrink-0 z-50">
      <div class="w-full px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="text-2xl transition-transform cursor-default"> 💕 </span>
          <h1
            class="text-2xl font-bold font-(family-name:--font-signature) admin-text-primary tracking-tight"
          >
            后台管理
          </h1>
        </div>

        <div class="flex items-center space-x-3">
          <!-- 移动端菜单按钮 -->
          <button
            @click="showMobileMenu = true"
            class="lg:hidden p-2 rounded-full hover:bg-black/5 active:bg-black/10 transition-colors"
          >
            <BaseIcon name="menu-right" size="w-6 h-6" color="var(--admin-accent-color)" />
          </button>

          <!-- 用户头像/用户名按钮 -->
          <button
            class="flex items-center gap-2 p-1 pr-3 rounded-full hover:bg-black/5 active:bg-black/10 transition-all border border-transparent hover:border-white/40"
            @click="showConfirmDialog = true"
          >
            <span
              class="inline-flex items-center justify-center h-9 w-9 rounded-full bg-gradient-to-br from-[#f0ada0] to-[#d89388] text-white font-bold shadow-sm"
            >
              {{ authStore.userInfo?.userName?.charAt(0).toUpperCase() || 'U' }}
            </span>
            <span class="hidden sm:inline text-sm font-medium admin-text-primary">
              {{ authStore.userInfo?.userName || '管理员' }}
            </span>
          </button>
        </div>
      </div>

      <!-- 退出确认对话框 - 保持 iOS Alert 风格覆盖 -->
      <GenericDialog
        variant="admin"
        :open="showConfirmDialog"
        title="退出登录"
        size-class="max-w-xs"
        @update:open="showConfirmDialog = $event"
        @cancel="showConfirmDialog = false"
      >
        <template #content>
          <p class="text-center py-2 admin-text-secondary text-sm">确定要退出登录吗？</p>
        </template>
        <template #actions>
          <div class="w-full flex border-t border-black/5">
            <button
              class="flex-1 py-3 text-center admin-text-primary font-medium hover:bg-black/5 transition-colors"
              @click="showConfirmDialog = false"
            >
              取消
            </button>
            <button
              class="flex-1 py-3 border-l border-black/5 text-center text-[#ff3b30] font-semibold hover:bg-black/5 transition-colors"
              @click="handleLogout"
            >
              退出
            </button>
          </div>
        </template>
      </GenericDialog>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边栏 - iPadOS 风格 -->
      <aside
        class="w-64 flex-shrink-0 hidden lg:flex flex-col border-r border-black/5 bg-white/20 backdrop-blur-md"
      >
        <nav class="flex-1 mt-6 px-4 space-y-1">
          <router-link
            v-for="item in menuItems"
            :key="item.path"
            :to="item.path"
            class="admin-sidebar-item group flex items-center px-4 py-3 text-sm font-semibold transition-all duration-200"
            :class="isActiveMenuItem(item.path) ? 'active' : 'hover:bg-black/5'"
          >
            <BaseIcon
              :name="item.icon"
              size="w-5 h-5"
              :color="isActiveMenuItem(item.path) ? '#ffffff' : 'var(--admin-accent-color)'"
              class="mr-3"
            />
            {{ item.label }}
          </router-link>
        </nav>

        <div class="p-4 border-t border-black/5">
          <router-link
            to="/"
            class="flex items-center p-3 rounded-2xl bg-white/40 hover:bg-white/60 transition-all group"
          >
            <div
              class="w-10 h-10 rounded-xl bg-gradient-to-br from-[#f0ada0] to-[#d89388] flex items-center justify-center mr-3 shadow-sm transition-transform"
            >
              <BaseIcon name="home" size="w-5 h-5" color="#ffffff" />
            </div>
            <div class="flex-1">
              <div class="text-sm font-bold admin-text-primary">返回前台</div>
              <div class="text-xs admin-text-secondary">浏览首页</div>
            </div>
          </router-link>
        </div>
      </aside>

      <!-- 主内容区域 -->
      <main class="flex-1 flex flex-col relative bg-gray-50/30 overflow-hidden">
        <div class="flex-1 p-4 sm:p-6 overflow-y-auto ios-touch-scroll">
          <router-view v-slot="{ Component }">
            <transition name="fade-slide" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>

    <!-- 移动端菜单 - iOS Bottom Sheet 风格 -->
    <Transition name="sheet">
      <div v-if="showMobileMenu" class="fixed inset-0 z-[100] lg:hidden">
        <!-- 背景遮罩 -->
        <div
          class="absolute inset-0 bg-black/30 backdrop-blur-sm"
          @click="showMobileMenu = false"
        ></div>

        <!-- Sheet 内容 -->
        <div
          class="absolute bottom-0 inset-x-0 bg-white/98 backdrop-blur-2xl rounded-t-[32px] p-6 shadow-[0_-8px_40px_rgba(0,0,0,0.15)] transition-transform"
        >
          <!-- 顶部把手 -->
          <div class="w-12 h-1.5 bg-black/10 rounded-full mx-auto mb-6"></div>

          <div class="flex justify-between items-center mb-6 px-2">
            <div class="flex items-center gap-2">
              <span class="text-2xl">💕</span>
              <h2 class="text-xl font-bold admin-text-primary">菜单</h2>
            </div>
            <button
              @click="showMobileMenu = false"
              class="p-2 rounded-full bg-black/5 hover:bg-black/10"
            >
              <BaseIcon name="close" size="w-5 h-5" color="var(--admin-accent-secondary)" />
            </button>
          </div>

          <nav class="grid grid-cols-2 gap-4 mb-8">
            <router-link
              v-for="item in menuItems"
              :key="item.path"
              :to="item.path"
              class="flex flex-col items-center justify-center p-4 rounded-2xl transition-all"
              :class="
                isActiveMenuItem(item.path)
                  ? 'bg-[var(--admin-accent-color)] text-white shadow-lg shadow-pink-200'
                  : 'bg-white/50 border border-black/5'
              "
              @click="showMobileMenu = false"
            >
              <BaseIcon
                :name="item.icon"
                size="w-8 h-8"
                :color="isActiveMenuItem(item.path) ? '#ffffff' : 'var(--admin-accent-color)'"
                class="mb-2"
              />
              <span class="text-sm font-bold">{{ item.label }}</span>
            </router-link>
          </nav>

          <div class="pb-safe">
            <router-link
              to="/"
              class="flex items-center p-4 rounded-2xl bg-black/5 active:bg-black/10 transition-all"
              @click="showMobileMenu = false"
            >
              <div
                class="w-10 h-10 rounded-full bg-gradient-to-br from-[#f0ada0] to-[#d89388] flex items-center justify-center mr-3"
              >
                <BaseIcon name="home" size="w-5 h-5" color="#ffffff" />
              </div>
              <div class="flex-1">
                <div class="text-sm font-bold admin-text-primary">返回前台首页</div>
              </div>
              <BaseIcon name="menu-right" size="w-5 h-5" color="var(--admin-text-secondary)" />
            </router-link>
          </div>
        </div>
      </div>
    </Transition>
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

<style scoped>
/* 页面转场动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Bottom Sheet 动画 */
.sheet-enter-active,
.sheet-leave-active {
  transition: all 0.4s cubic-bezier(0.32, 0.72, 0, 1);
}

.sheet-enter-from .bg-white\/90 {
  transform: translateY(100%);
}

.sheet-enter-from .absolute.inset-0 {
  opacity: 0;
}

.sheet-leave-to .bg-white\/90 {
  transform: translateY(100%);
}

.sheet-leave-to .absolute.inset-0 {
  opacity: 0;
}

.pb-safe {
  padding-bottom: env(safe-area-inset-bottom, 20px);
}
</style>
