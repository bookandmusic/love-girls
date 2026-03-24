import { createRouter, createWebHistory } from "vue-router";

import DefaultLayout from "@/layouts/DefaultLayout.vue";

import HomeView from "@/views/HomeView.vue";
import PlacesView from "@/views/PlacesView.vue";
import AlbumsView from "@/views/AlbumsView.vue";
import MomentsView from "@/views/MomentsView.vue";
import AnniversariesView from "@/views/AnniversariesView.vue";
import NotFoundView from "@/views/NotFoundView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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
      ],
    },
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: NotFoundView,
    },
  ],
});

export default router;
