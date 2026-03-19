<template>
  <div class="server-config-page frontend-root">
    <div
      class="absolute inset-0 bg-gradient-to-br from-[#fce4ec] via-[#f8bbd9]/30 to-[#f48fb1]/20"
    ></div>
    <div
      class="absolute top-0 right-0 w-96 h-96 bg-[#f0ada0]/40 rounded-full blur-3xl transform translate-x-1/2 -translate-y-1/2"
    ></div>
    <div
      class="absolute bottom-0 left-0 w-[500px] h-[500px] bg-[#d89388]/30 rounded-full blur-3xl transform -translate-x-1/3 translate-y-1/3"
    ></div>

    <Transition name="toast">
      <div
        v-if="toastMsg"
        class="fixed top-8 left-1/2 -translate-x-1/2 z-50 px-6 py-3 rounded-xl glass-regular border border-white/40 shadow-lg"
        :class="toastType === 'error' ? 'text-red-500' : 'text-green-600'"
      >
        {{ toastMsg }}
      </div>
    </Transition>

    <div class="relative z-10 w-full max-w-md mx-auto px-6">
      <div class="text-center mb-8">
        <img
          src="/favicon.png"
          alt="Love Girl"
          class="w-20 h-20 mx-auto mb-6 rounded-2xl shadow-lg"
        />
        <h1 class="text-2xl font-bold text-[var(--fe-text-primary)] mb-2">欢迎使用 Love Girl</h1>
        <p class="text-sm text-[var(--fe-text-secondary)]">请配置您的服务器地址</p>
      </div>

      <!-- 已保存的服务器列表 - iOS风格下拉展开 -->
      <div v-if="savedServers.length > 0" class="mb-4 relative" ref="serverListRef">
        <button
          @click="showServerList = !showServerList"
          class="w-full glass-regular rounded-xl p-4 border border-white/30 hover:border-white/50 transition-all"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-xl bg-gradient-to-br from-[var(--fe-primary)] to-[var(--fe-primary-dark)] flex items-center justify-center"
              >
                <svg
                  class="w-5 h-5 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
                  />
                </svg>
              </div>
              <div class="text-left">
                <div class="flex items-center gap-2">
                  <span class="text-sm font-medium text-[var(--fe-text-primary)]">{{
                    currentServer?.name || '选择服务器'
                  }}</span>
                  <span
                    v-if="currentServer"
                    class="text-xs px-2 py-0.5 rounded-full bg-[var(--fe-primary)]/10 text-[var(--fe-primary)]"
                    >当前</span
                  >
                </div>
                <p v-if="currentServer" class="text-xs text-[var(--fe-text-secondary)] mt-0.5">
                  {{ currentServer.url }}
                </p>
              </div>
            </div>
            <svg
              class="w-5 h-5 text-[var(--fe-text-secondary)] transition-transform duration-300"
              :class="showServerList ? 'rotate-180' : ''"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
        </button>

        <Transition name="expand">
          <div
            v-if="showServerList"
            class="absolute top-full left-0 right-0 mt-2 z-20 glass-regular rounded-xl border border-white/40 shadow-xl h-48 overflow-y-auto"
          >
            <div class="p-2 space-y-1">
              <div
                v-for="server in savedServers"
                :key="server.url"
                class="rounded-xl p-3 border transition-all cursor-pointer active:scale-[0.98]"
                :class="
                  activeServerUrl === server.url
                    ? 'border-[var(--fe-primary)] bg-[var(--fe-primary)]/5'
                    : 'border-transparent hover:bg-white/30'
                "
                @click="selectServer(server.url)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2">
                      <span class="text-sm font-medium text-[var(--fe-text-primary)]">{{
                        server.name
                      }}</span>
                      <span
                        v-if="activeServerUrl === server.url"
                        class="text-xs px-2 py-0.5 rounded-full bg-[var(--fe-primary)] text-white"
                        >当前</span
                      >
                    </div>
                    <p class="text-xs text-[var(--fe-text-secondary)] truncate mt-0.5">
                      {{ server.url }}
                    </p>
                  </div>
                  <div class="flex items-center gap-1 ml-2">
                    <button
                      @click.stop="testSavedServer(server.url)"
                      :disabled="testingUrl === server.url"
                      class="p-1.5 rounded-lg hover:bg-white/40 active:bg-white/60 transition-colors"
                      title="测试连接"
                    >
                      <svg
                        v-if="testingUrl !== server.url"
                        class="w-4 h-4 text-[var(--fe-text-secondary)]"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.14 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"
                        />
                      </svg>
                      <svg
                        v-else
                        class="w-4 h-4 text-[var(--fe-primary)] animate-spin"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          class="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          stroke-width="4"
                        ></circle>
                        <path
                          class="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                        ></path>
                      </svg>
                    </button>
                    <button
                      @click.stop="deleteServer(server.url)"
                      class="p-1.5 rounded-lg hover:bg-red-50 active:bg-red-100 transition-colors"
                      title="删除"
                    >
                      <svg
                        class="w-4 h-4 text-red-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>

      <!-- 添加新服务器 -->
      <div
        class="glass-regular rounded-[var(--fe-radius-card)] p-6 border border-white/40 shadow-xl"
      >
        <div class="space-y-4">
          <div class="space-y-2">
            <label class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1">
              名称
            </label>
            <input
              v-model="newServerName"
              type="text"
              class="w-full px-4 py-3 rounded-xl border border-white/30 bg-white/50 backdrop-blur-sm text-[var(--fe-text-primary)] placeholder-[var(--fe-text-secondary)] focus:outline-none focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 transition-all"
              placeholder="如：内网、外网"
            />
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1">
              服务器地址
            </label>
            <input
              v-model="newServerUrl"
              type="url"
              class="w-full px-4 py-3 rounded-xl border border-white/30 bg-white/50 backdrop-blur-sm text-[var(--fe-text-primary)] placeholder-[var(--fe-text-secondary)] focus:outline-none focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 transition-all"
              placeholder="http://192.168.1.100:8182"
              @keyup.enter="handleAddServer"
            />
          </div>

          <p v-if="errorMsg" class="text-sm text-red-500 px-1">
            {{ errorMsg }}
          </p>

          <div class="flex gap-3 pt-2">
            <button
              @click="testNewServer"
              :disabled="testing || !newServerUrl.trim()"
              class="flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all glass-ultra-thin border border-white/30 text-[var(--fe-text-primary)] hover:bg-white/50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ testing ? '测试中...' : '测试连接' }}
            </button>
            <button
              @click="handleAddServer"
              :disabled="connecting || !newServerUrl.trim()"
              class="flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all bg-gradient-to-r from-[var(--fe-primary)] to-[var(--fe-primary-dark)] text-white shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ connecting ? '保存中...' : '添加并使用' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 使用当前选中服务器进入 -->
      <button
        v-if="activeServerUrl"
        @click="enterApp"
        class="w-full mt-4 py-3 px-4 rounded-xl font-semibold text-sm transition-all bg-gradient-to-r from-[var(--fe-primary)] to-[var(--fe-primary-dark)] text-white shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98]"
      >
        进入应用
      </button>

      <p class="text-center text-xs text-[var(--fe-text-secondary)] mt-6">
        请输入后端服务地址，如 http://192.168.1.100:8182
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { refreshApiBaseURL } from '@/services/api'
import {
  addServerUrl,
  getActiveServerUrl,
  getServerUrls,
  removeServerUrl,
  type ServerConfig,
  setActiveServerUrl,
  validateServerUrl,
} from '@/utils/platform'

const router = useRouter()

const savedServers = ref<ServerConfig[]>([])
const activeServerUrl = ref<string | null>(null)
const showServerList = ref(false)
const serverListRef = ref<HTMLElement | null>(null)
const newServerName = ref('')
const newServerUrl = ref('')
const errorMsg = ref('')
const testing = ref(false)
const connecting = ref(false)
const testingUrl = ref<string | null>(null)
const toastMsg = ref('')
const toastType = ref<'success' | 'error'>('success')

const currentServer = computed(() => savedServers.value.find(s => s.url === activeServerUrl.value))

const handleClickOutside = (event: MouseEvent) => {
  if (serverListRef.value && !serverListRef.value.contains(event.target as Node)) {
    showServerList.value = false
  }
}

onMounted(() => {
  savedServers.value = getServerUrls()
  activeServerUrl.value = getActiveServerUrl()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const showToast = (msg: string, type: 'success' | 'error' = 'success') => {
  toastMsg.value = msg
  toastType.value = type
  setTimeout(() => {
    toastMsg.value = ''
  }, 3000)
}

const testServerUrl = async (url: string): Promise<boolean> => {
  const validation = validateServerUrl(url)
  if (!validation.valid) {
    errorMsg.value = validation.error || ''
    return false
  }

  try {
    const response = await fetch(`${url}/api/v1/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    })

    if (response.ok) {
      return true
    } else {
      errorMsg.value = '服务器响应异常'
      return false
    }
  } catch {
    errorMsg.value = '无法连接到服务器，请检查地址'
    return false
  }
}

const testNewServer = async () => {
  errorMsg.value = ''
  testing.value = true

  const success = await testServerUrl(newServerUrl.value)
  if (success) {
    showToast('连接成功', 'success')
  }

  testing.value = false
}

const testSavedServer = async (url: string) => {
  testingUrl.value = url
  errorMsg.value = ''

  const success = await testServerUrl(url)
  if (success) {
    showToast('连接成功', 'success')
  }

  testingUrl.value = null
}

const handleAddServer = async () => {
  const validation = validateServerUrl(newServerUrl.value)
  if (!validation.valid) {
    errorMsg.value = validation.error || ''
    return
  }

  errorMsg.value = ''
  connecting.value = true

  const success = await testServerUrl(newServerUrl.value)
  if (!success) {
    connecting.value = false
    return
  }

  const name =
    newServerName.value.trim() ||
    (savedServers.value.length === 0 ? '默认' : `服务器${savedServers.value.length + 1}`)
  savedServers.value = addServerUrl(name, newServerUrl.value)
  activeServerUrl.value = newServerUrl.value
  setActiveServerUrl(newServerUrl.value)
  refreshApiBaseURL()

  newServerName.value = ''
  newServerUrl.value = ''

  showToast('添加成功', 'success')
  router.push('/')
  connecting.value = false
}

const selectServer = (url: string) => {
  activeServerUrl.value = url
  setActiveServerUrl(url)
  refreshApiBaseURL()
  showServerList.value = false

  const server = savedServers.value.find(s => s.url === url)
  if (server) {
    newServerName.value = server.name
    newServerUrl.value = server.url
  }
}

const deleteServer = (url: string) => {
  savedServers.value = removeServerUrl(url)
  activeServerUrl.value = getActiveServerUrl()
  if (savedServers.value.length === 0) {
    showServerList.value = false
  }
}

const enterApp = () => {
  if (activeServerUrl.value) {
    refreshApiBaseURL()
    router.push('/')
  }
}
</script>

<style scoped>
.server-config-page {
  min-height: 100vh;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -20px);
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.95);
}
</style>
