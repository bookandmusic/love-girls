import { createRouter, createWebHistory } from "vue-router";

import DefaultLayout from "@/layouts/DefaultLayout.vue";
import AdminLayout from "@/views/AdminLayout.vue";
import DashboardView from "@/views/DashboardView.vue";
import UsersView from "@/views/users/UsersView.vue";
import ContentView from "@/views/content/ContentView.vue";
import MomentsManagement from "@/views/content/moments/MomentsManagementView.vue";
import AnniversariesManagement from "@/views/content/anniversaries/AnniversariesManagementView.vue";
import PlacesManagement from "@/views/content/places/PlacesManagementView.vue";
import AlbumsManagement from "@/views/content/albums/AlbumsManagementView.vue";
import SettingsView from "@/views/SettingsView.vue";
import LoginView from "@/views/LoginView.vue";
import InitSystemView from "@/views/InitSystemView.vue";
import NotFoundView from "@/views/NotFoundView.vue";

import { useAuthStore } from "@/stores/auth";
import { useSystemStore } from "@/stores/system";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: DefaultLayout,
      children: [
        {
          path: "",
          redirect: "/dashboard",
        },
        {
          path: "init",
          name: "init",
          component: InitSystemView,
        },
        {
          path: "login",
          name: "login",
          component: LoginView,
        },
        {
          path: "/",
          component: AdminLayout,
          children: [
            {
              path: "dashboard",
              name: "dashboard",
              component: DashboardView,
            },
            {
              path: "users",
              name: "users",
              component: UsersView,
            },
            {
              path: "content",
              name: "content",
              component: ContentView,
              redirect: "/content/moments",
              children: [
                {
                  path: "moments",
                  name: "content-moments",
                  component: MomentsManagement,
                },
                {
                  path: "anniversaries",
                  name: "content-anniversaries",
                  component: AnniversariesManagement,
                },
                {
                  path: "places",
                  name: "content-places",
                  component: PlacesManagement,
                },
                {
                  path: "albums",
                  name: "content-albums",
                  component: AlbumsManagement,
                },
              ],
            },
            {
              path: "settings",
              name: "settings",
              component: SettingsView,
            },
          ],
        },
        {
          path: "/:pathMatch(.*)*",
          name: "not-found",
          component: NotFoundView,
        },
      ],
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const systemStore = useSystemStore();
  const authStore = useAuthStore();

  const isInitialized = await systemStore.checkInitialization();

  if (!isInitialized && to.path !== "/init") {
    return next("/init");
  }

  if (isInitialized && to.path === "/init") {
    return next("/login");
  }

  if (isInitialized) {
    const isAuthenticated = await authStore.checkAuthStatus();

    if (!isAuthenticated && to.path !== "/login") {
      return next("/login");
    }

    if (isAuthenticated && to.path === "/login") {
      return next("/dashboard");
    }
  }

  return next();
});

export default router;
