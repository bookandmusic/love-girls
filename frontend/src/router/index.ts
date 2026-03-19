import { createRouter, createWebHashHistory } from 'vue-router'

import DefaultLayout from '@/layouts/DefaultLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { useSystemStore } from '@/stores/system'
// 后台页面
import AdminLayout from '@/views/admin/AdminLayout.vue'
// 内容管理子页面
import AlbumsManagement from '@/views/admin/content/albums/AlbumsManagementView.vue'
import AnniversariesManagement from '@/views/admin/content/anniversaries/AnniversariesManagementView.vue'
import ContentView from '@/views/admin/content/ContentView.vue'
import MomentsManagement from '@/views/admin/content/moments/MomentsManagementView.vue'
import PlacesManagement from '@/views/admin/content/places/PlacesManagementView.vue'
import DashboardView from '@/views/admin/DashboardView.vue'
import SettingsView from '@/views/admin/SettingsView.vue'
import UsersView from '@/views/admin/users/UsersView.vue'
import AlbumsView from '@/views/frontend/AlbumsView.vue'
import AnniversariesView from '@/views/frontend/AnniversariesView.vue'
// 前台页面
import HomeView from '@/views/frontend/HomeView.vue'
import MomentsView from '@/views/frontend/MomentsView.vue'
import PlacesView from '@/views/frontend/PlacesView.vue'
// 初始化页面
import InitSystemView from '@/views/InitSystemView.vue'
// 登录页面
import LoginView from '@/views/LoginView.vue'

import { getServerUrl, isDesktopMode } from '../utils/platform'
import NotFoundView from '../views/NotFoundView.vue'
import ServerConfigView from '../views/ServerConfigView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/server-config',
      name: 'ServerConfig',
      component: ServerConfigView,
      meta: { requiresAuth: false },
    },
    // 管理页面路由
    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'places',
          name: 'places',
          component: PlacesView,
        },
        {
          path: 'albums',
          component: AlbumsView,
          name: 'albums',
        },
        {
          path: 'moments',
          component: MomentsView,
          name: 'moments',
        },
        {
          path: 'anniversaries',
          component: AnniversariesView,
          name: 'anniversaries',
        },
        {
          path: 'init',
          name: 'init',
          component: InitSystemView,
        },
        // 登录页面路由
        {
          path: 'login',
          name: 'login',
          component: LoginView,
        },
        {
          path: 'admin',
          component: AdminLayout,
          redirect: '/admin/dashboard',
          children: [
            {
              path: 'dashboard',
              name: 'admin-dashboard',
              component: DashboardView,
            },
            {
              path: 'users',
              name: 'admin-users',
              component: UsersView,
            },
            {
              path: 'content',
              name: 'admin-content',
              component: ContentView,
              redirect: '/admin/content/moments',
              children: [
                {
                  path: 'moments',
                  name: 'admin-content-moments',
                  component: MomentsManagement,
                },
                {
                  path: 'anniversaries',
                  name: 'admin-content-anniversaries',
                  component: AnniversariesManagement,
                },
                {
                  path: 'places',
                  name: 'admin-content-places',
                  component: PlacesManagement,
                },
                {
                  path: 'albums',
                  name: 'admin-content-albums',
                  component: AlbumsManagement,
                },
              ],
            },
            {
              path: 'settings',
              name: 'admin-settings',
              component: SettingsView,
            },
          ],
        },
      ],
    },

    // 404页面路由，必须放在最后
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFoundView,
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  // Desktop 模式下，优先检查服务器配置
  if (isDesktopMode()) {
    if (to.path !== '/server-config' && !getServerUrl()) {
      return next('/server-config')
    }
    if (to.path === '/server-config') {
      return next()
    }
  }

  const systemStore = useSystemStore()
  const authStore = useAuthStore()

  const isInitialized = await systemStore.checkInitialization()
  const isAdminRoute = to.path.startsWith('/admin')
  const isAuthenticated = await authStore.checkAuthStatus()

  if (to.path === '/init') {
    if (isInitialized) {
      return next('/')
    }
  }

  if (!isInitialized && to.path !== '/init') {
    return next('/init')
  }

  if (isAdminRoute && !isAuthenticated && to.path !== '/login') {
    return next('/login')
  }
  if (isAuthenticated && to.path === '/login') {
    return next('/admin')
  }

  return next()
})

export default router
