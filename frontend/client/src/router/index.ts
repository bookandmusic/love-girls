import { createRouter, createWebHashHistory } from "vue-router";

import DefaultLayout from "@/layouts/DefaultLayout.vue";
import { useAuthStore } from "@/stores/auth";
import { refreshApiBaseURL } from "@/services/api";

import { getActiveServerUrl } from "@/utils/platform";

import HomeView from "@/views/HomeView.vue";
import PlacesView from "@/views/PlacesView.vue";
import AlbumsView from "@/views/AlbumsView.vue";
import MomentsView from "@/views/MomentsView.vue";
import AnniversariesView from "@/views/AnniversariesView.vue";
import NotificationsView from "@/views/NotificationsView.vue";
import ServerConfigView from "@/views/ServerConfigView.vue";
import NotFoundView from "@/views/NotFoundView.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/server-config",
      name: "ServerConfig",
      component: ServerConfigView,
      meta: { requiresAuth: false },
    },
    {
      path: "/",
      component: DefaultLayout,
      children: [
        {
          path: "",
          name: "home",
          component: HomeView,
        },
        {
          path: "places",
          name: "places",
          component: PlacesView,
        },
        {
          path: "albums",
          name: "albums",
          component: AlbumsView,
        },
        {
          path: "moments",
          name: "moments",
          component: MomentsView,
        },
        {
          path: "anniversaries",
          name: "anniversaries",
          component: AnniversariesView,
        },
        {
          path: "notifications",
          name: "notifications",
          component: NotificationsView,
        },
      ],
    },
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: NotFoundView,
    },
  ],
});

router.beforeEach((to, from, next) => {
  if (to.path !== "/server-config" && !getActiveServerUrl()) {
    return next("/server-config");
  }
  if (to.path === "/server-config") {
    return next();
  }
  refreshApiBaseURL();

  const authStore = useAuthStore();

  if (!authStore.isAuthenticated) {
    return next({
      path: "/server-config",
      query: { step: "login", redirect: to.fullPath },
    });
  }

  return next();
});

export default router;
