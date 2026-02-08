import { createRouter, createWebHistory } from 'vue-router'

import DefaultLayout from '@/layouts/DefaultLayout.vue'
import { useAuthStore } from '@/stores/auth'
import { useSystemStore } from '@/stores/system'
// 后台页面
import AdminLayout from '@/views/admin/AdminLayout.vue'
import ContentView from '@/views/admin/content/ContentView.vue'
import DashboardView from '@/views/admin/DashboardView.vue'
import SettingsView from '@/views/admin/SettingsView.vue'
import UsersView from '@/views/admin/users/UsersView.vue'
import AlbumsView from '@/views/frontend/AlbumsView.vue'
import AnniversariesView from '@/views/frontend/AnniversariesView.vue'
// 前台页面
import HomeView from '@/views/frontend/HomeView.vue'
import MomentsView from '@/views/frontend/MomentsView.vue'
import PlacesView from '@/views/frontend/PlacesView.vue'
import WishesView from '@/views/frontend/WishesView.vue'
// 初始化页面
import InitSystemView from '@/views/InitSystemView.vue'
// 登录页面
import LoginView from '@/views/LoginView.vue'

import NotFoundView from '../views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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
          path: 'wishes',
          component: WishesView,
          name: 'wishes',
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
  const systemStore = useSystemStore()
  const authStore = useAuthStore()

  const isInitialized = await systemStore.checkInitialization()
  // 检查是否是管理员路由
  const isAdminRoute = to.path.startsWith('/admin')
  const isAuthenticated = await authStore.checkAuthStatus()
  // 如果正在访问初始化页面
  if (to.path === '/init') {
    // 检查系统是否已经初始化
    if (isInitialized) {
      return next('/')
    }
  }

  // 对于其他页面，检查系统是否已初始化
  if (!isInitialized && to.path !== '/init') {
    return next('/init')
  }

  // 对于管理员页面，需要检查登录状态
  if (isAdminRoute && !isAuthenticated && to.path !== '/login') {
    return next('/login')
  }
  if (isAuthenticated && to.path === '/login') {
    return next('/admin')
  }

  // 允许访问目标页面
  return next()
})

export default router
