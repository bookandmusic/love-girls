<template>
  <div class="relative inline-block" ref="menuRef">
    <button
      @click="toggleMenu"
      class="menu-button p-1.5 md:p-2.5 rounded-lg glass-regular border border-white/30 ios-transition active:scale-95"
      :class="{ 'border-[var(--fe-primary)]/50': showMenu }"
    >
      <BaseIcon name="menu-dots" size="w-3.5 h-3.5 md:w-5 md:h-5" color="var(--fe-text-primary)" />
    </button>

    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="showMenu"
          class="dropdown-menu rounded-lg md:rounded-xl border border-white/20 shadow-xl overflow-hidden"
          :style="{ top: menuPosition.top + 'px', left: menuPosition.left + 'px' }"
          ref="dropdownRef"
        >
          <div class="p-0.5 md:p-1.5">
            <div class="dropdown-info">
              <BaseIcon
                name="server"
                size="w-3 h-3 md:w-4 md:h-4"
                color="var(--fe-text-secondary)"
              />
              <span class="text-[10px] md:text-sm text-[var(--fe-text-primary)] truncate">{{
                currentServerUrl
              }}</span>
            </div>
            <button
              @click="handleSwitchServer"
              class="dropdown-item w-full text-left rounded-md md:rounded-lg"
            >
              <BaseIcon
                name="logout"
                size="w-3 h-3 md:w-4 md:h-4"
                color="var(--fe-text-secondary)"
              />
              <span class="text-[10px] md:text-sm font-medium text-[var(--fe-text-primary)]"
                >切换服务器</span
              >
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import { getActiveServerUrl } from '@/utils/platform'

const router = useRouter()
const showMenu = ref(false)
const menuRef = ref<HTMLElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)

const currentServerUrl = computed(() => getActiveServerUrl() || '未配置')

const menuPosition = ref({ top: 0, left: 0 })

const updatePosition = () => {
  if (menuRef.value) {
    const rect = menuRef.value.getBoundingClientRect()
    const menuWidth = window.innerWidth >= 768 ? 200 : 140
    menuPosition.value = {
      top: rect.bottom + 8,
      left: rect.right - menuWidth,
    }
  }
}

watch(showMenu, val => {
  if (val) {
    updatePosition()
  }
})

const toggleMenu = () => {
  showMenu.value = !showMenu.value
}

const closeMenu = () => {
  showMenu.value = false
}

const handleSwitchServer = () => {
  closeMenu()
  router.push('/server-config')
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as Node
  if (
    menuRef.value &&
    !menuRef.value.contains(target) &&
    dropdownRef.value &&
    !dropdownRef.value.contains(target)
  ) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('resize', updatePosition)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('resize', updatePosition)
})
</script>

<style scoped>
.menu-button {
  display: flex;
  align-items: center;
  justify-content: center;
}

.dropdown-menu {
  position: fixed;
  min-width: 140px;
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  z-index: 9999;
}

@media (min-width: 768px) {
  .dropdown-menu {
    min-width: 200px;
  }
}

.dropdown-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

@media (min-width: 768px) {
  .dropdown-info {
    gap: 10px;
    padding: 10px 12px;
  }
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  transition: all 0.15s ease;
  cursor: pointer;
}

@media (min-width: 768px) {
  .dropdown-item {
    gap: 10px;
    padding: 10px 12px;
  }
}

.dropdown-item:hover {
  background: rgba(0, 0, 0, 0.04);
}

.dropdown-item:active {
  background: rgba(0, 0, 0, 0.08);
  transform: scale(0.98);
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.95);
}
</style>
