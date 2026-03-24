<script setup lang="ts">
import { useRoute } from "vue-router";

import BaseIcon from "@/components/ui/BaseIcon.vue";

// 定义菜单项类型
interface MenuItem {
  icon: string;
  label: string;
  path: string;
}

defineProps<{
  isMobile?: boolean;
}>();

const route = useRoute();

const menuItems: MenuItem[] = [
  { icon: "home-heart", label: "首页", path: "/" },
  { icon: "moment", label: "动态", path: "/moments" },
  { icon: "photo-heart", label: "相册", path: "/albums" },
  { icon: "place", label: "足迹", path: "/places" },
  { icon: "anniversary", label: "纪念日", path: "/anniversaries" },
];

const isActive = (path: string) => {
  if (path === "/") return route.path === "/";
  return route.path.startsWith(path);
};
</script>

<template>
  <!-- PC端左侧悬浮导航栏 -->
  <div
    v-if="!isMobile"
    class="hidden md:flex flex-col items-center justify-center w-28 py-8 h-full"
  >
    <div
      class="glass-ultra-thin rounded-full py-6 px-3 flex flex-col space-y-6 border border-white/20 shadow-lg"
    >
      <RouterLink
        v-for="(item, index) in menuItems"
        :key="index"
        :to="item.path"
        class="relative group flex flex-col items-center tap-feedback ios-transition"
      >
        <div
          class="p-3 rounded-2xl ios-transition flex items-center justify-center"
          :class="isActive(item.path) ? 'bg-white/40 shadow-sm scale-110' : ''"
        >
          <BaseIcon
            :name="item.icon"
            size="w-7 h-7"
            :color="
              isActive(item.path)
                ? 'text-[var(--fe-primary)]'
                : 'text-[var(--fe-text-primary)] opacity-70'
            "
          />
        </div>
        <!-- Tooltip -->
        <span
          class="absolute left-full ml-4 px-2 py-1 glass-thick rounded-lg text-xs opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none border border-white/20 shadow-sm"
        >
          {{ item.label }}
        </span>
      </RouterLink>
    </div>
  </div>

  <!-- 手机端底部 Tab Bar -->
  <div
    v-else
    class="md:hidden glass-thick border-t border-white/20 px-2 pt-1.5 pb-[calc(0.375rem+var(--fe-safe-area-bottom))] z-50"
  >
    <div class="flex justify-around items-center max-w-lg mx-auto">
      <RouterLink
        v-for="(item, index) in menuItems"
        :key="index"
        :to="item.path"
        class="flex items-center justify-center flex-1 py-1 tap-feedback ios-transition"
      >
        <div
          class="p-1.5 rounded-xl ios-transition flex items-center justify-center"
          :class="isActive(item.path) ? 'bg-white/40' : ''"
        >
          <BaseIcon
            :name="item.icon"
            size="w-7 h-7"
            :color="
              isActive(item.path)
                ? 'text-[var(--fe-primary)]'
                : 'text-[var(--fe-text-secondary)]'
            "
          />
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.tap-feedback:active {
  transform: scale(0.92);
  transition: transform 0.1s ease-out;
}
</style>
