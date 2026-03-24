import { createRouter, createWebHashHistory } from "vue-router";

import DefaultLayout from "@/layouts/DefaultLayout.vue";
import { useSystemStore } from "@/stores/system";
import { refreshApiBaseURL } from "@/services/api";

import { getActiveServerUrl } from "@/utils/platform";

import HomeView from "@/views/HomeView.vue";
import PlacesView from "@/views/PlacesView.vue";
import AlbumsView from "@/views/AlbumsView.vue";
import MomentsView from "@/views/MomentsView.vue";
import AnniversariesView from "@/views/AnniversariesView.vue";
import InitSystemView from "@/views/InitSystemView.vue";
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
          path: "init",
          name: "init",
          component: InitSystemView,
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

router.beforeEach(async (to, from, next) => {
  if (to.path !== "/server-config" && !getActiveServerUrl()) {
    return next("/server-config");
  }
  if (to.path === "/server-config") {
    return next();
  }
  refreshApiBaseURL();

  const systemStore = useSystemStore();

  const isInitialized = await systemStore.checkInitialization();

  if (to.path === "/init") {
    if (isInitialized) {
      return next("/");
    }
    return next();
  }

  if (!isInitialized && systemStore.networkError) {
    return next("/server-config");
  }

  if (!isInitialized && to.path !== "/init") {
    return next("/init");
  }

  return next();
});

export default router;
