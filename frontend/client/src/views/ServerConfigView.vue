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
      <div class="text-center mb-8 h-[140px] flex flex-col justify-end">
        <img
          src="/favicon.png"
          alt="Love Girl"
          class="w-20 h-20 mx-auto mb-6 rounded-2xl shadow-lg"
        />
        <h1 class="text-2xl font-bold text-[var(--fe-text-primary)] mb-2">
          {{ step === "server" ? "欢迎使用 Love Girl" : "登录" }}
        </h1>
        <p class="text-sm text-[var(--fe-text-secondary)]">
          {{
            step === "server" ? "请配置您的服务器地址" : "请输入账号密码登录"
          }}
        </p>
      </div>

      <!-- 步骤指示器 -->
      <div class="flex items-center justify-center gap-2 mb-6 h-6">
        <div
          class="flex items-center gap-1.5"
          :class="
            step === 'server'
              ? 'text-[var(--fe-primary)]'
              : 'text-[var(--fe-text-secondary)]'
          "
        >
          <div
            class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
            :class="
              step === 'server'
                ? 'bg-[var(--fe-primary)] text-white'
                : 'bg-[var(--fe-primary)] text-white'
            "
          >
            1
          </div>
          <span class="text-xs font-medium">服务器</span>
        </div>
        <div class="w-8 h-0.5 bg-[var(--fe-text-secondary)]/30 rounded"></div>
        <div
          class="flex items-center gap-1.5"
          :class="
            step === 'login'
              ? 'text-[var(--fe-primary)]'
              : 'text-[var(--fe-text-secondary)]'
          "
        >
          <div
            class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
            :class="
              step === 'login'
                ? 'bg-[var(--fe-primary)] text-white'
                : 'bg-gray-300 text-white'
            "
          >
            2
          </div>
          <span class="text-xs font-medium">登录</span>
        </div>
      </div>

      <!-- 第一步：服务器配置 -->
      <div class="w-full h-[420px] overflow-y-auto">
        <template v-if="step === 'server'">
          <!-- 已保存的服务器列表 -->
          <div
            v-if="savedServers.length > 0"
            class="mb-4 relative"
            ref="serverListRef"
          >
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
                      <span
                        class="text-sm font-medium text-[var(--fe-text-primary)]"
                        >{{ currentServer?.name || "选择服务器" }}</span
                      >
                      <span
                        v-if="currentServer"
                        class="text-xs px-2 py-0.5 rounded-full bg-[var(--fe-primary)]/10 text-[var(--fe-primary)]"
                        >当前</span
                      >
                    </div>
                    <p
                      v-if="currentServer"
                      class="text-xs text-[var(--fe-text-secondary)] mt-0.5"
                    >
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
                    @click="selectServer(server)"
                  >
                    <div class="flex items-center justify-between">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2">
                          <span
                            class="text-sm font-medium text-[var(--fe-text-primary)]"
                            >{{ server.name }}</span
                          >
                          <span
                            v-if="activeServerUrl === server.url"
                            class="text-xs px-2 py-0.5 rounded-full bg-[var(--fe-primary)] text-white"
                            >当前</span
                          >
                        </div>
                        <p
                          class="text-xs text-[var(--fe-text-secondary)] truncate mt-0.5"
                        >
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
                <label
                  class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1"
                >
                  名称
                </label>
                <input
                  v-model="newServerName"
                  type="text"
                  class="w-full win11-input"
                  placeholder="如：内网、外网"
                />
              </div>

              <div class="space-y-2">
                <label
                  class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1"
                >
                  服务器地址
                </label>
                <input
                  v-model="newServerUrl"
                  type="url"
                  class="w-full win11-input"
                  placeholder="http://192.168.1.100:8182"
                  @keyup.enter="handleConnectServer"
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
                  {{ testing ? "测试中..." : "测试连接" }}
                </button>
                <button
                  @click="handleConnectServer"
                  :disabled="connecting || !newServerUrl.trim()"
                  class="flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all bg-gradient-to-r from-[var(--fe-primary)] to-[var(--fe-primary-dark)] text-white shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ connecting ? "连接中..." : "连接" }}
                </button>
              </div>
            </div>
          </div>

          <p class="text-center text-xs text-[var(--fe-text-secondary)] mt-6">
            请输入后端服务地址，如 http://192.168.1.100:8182
          </p>
        </template>

        <!-- 第二步：登录 -->
        <template v-if="step === 'login'">
          <div
            class="glass-regular rounded-[var(--fe-radius-card)] p-6 border border-white/40 shadow-xl"
          >
            <!-- 当前服务器信息 -->
            <div
              class="mb-6 p-3 rounded-xl bg-white/30 border border-white/40 flex items-center gap-3"
            >
              <div
                class="w-8 h-8 rounded-lg bg-gradient-to-br from-[var(--fe-primary)] to-[var(--fe-primary-dark)] flex items-center justify-center"
              >
                <svg
                  class="w-4 h-4 text-white"
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
              <div class="flex-1 min-w-0">
                <div
                  class="text-sm font-medium text-[var(--fe-text-primary)] truncate"
                >
                  {{ currentServer?.name || "服务器" }}
                </div>
                <div class="text-xs text-[var(--fe-text-secondary)] truncate">
                  {{ currentServer?.url }}
                </div>
              </div>
              <button
                @click="switchServer"
                class="text-xs text-[var(--fe-primary)] font-medium hover:underline"
              >
                切换
              </button>
            </div>

            <!-- 登录表单 -->
            <div class="space-y-4">
              <div class="space-y-2">
                <label
                  class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1"
                >
                  用户名
                </label>
                <input
                  v-model="username"
                  type="text"
                  class="w-full win11-input"
                  placeholder="请输入用户名"
                  @keyup.enter="handleLogin"
                />
              </div>

              <div class="space-y-2">
                <label
                  class="block text-sm font-semibold text-[var(--fe-text-primary)] px-1"
                >
                  密码
                </label>
                <input
                  v-model="password"
                  type="password"
                  class="w-full win11-input"
                  placeholder="请输入密码"
                  @keyup.enter="handleLogin"
                />
              </div>

              <p v-if="errorMsg" class="text-sm text-red-500 px-1">
                {{ errorMsg }}
              </p>

              <div class="flex gap-3 pt-2">
                <button
                  @click="switchServer"
                  class="flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all glass-ultra-thin border border-white/30 text-[var(--fe-text-primary)] hover:bg-white/50"
                >
                  切换服务器
                </button>
                <button
                  @click="handleLogin"
                  :disabled="loggingIn || !username.trim() || !password.trim()"
                  class="flex-1 py-3 px-4 rounded-xl font-semibold text-sm transition-all bg-gradient-to-r from-[var(--fe-primary)] to-[var(--fe-primary-dark)] text-white shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ loggingIn ? "登录中..." : "登录" }}
                </button>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useRouter } from "vue-router";

import { refreshApiBaseURL } from "@/services/api";
import { userApi } from "@/services/userApi";
import { useAuthStore } from "@/stores/auth";
import { useSystemStore } from "@/stores/system";
import {
  addServerUrl,
  getActiveServerUrl,
  getActiveServerToken,
  getServerUrls,
  removeServerUrl,
  type ServerConfig,
  setActiveServerUrl,
  validateServerUrl,
} from "@/utils/platform";

const router = useRouter();
const systemStore = useSystemStore();
const authStore = useAuthStore();

const step = ref<"server" | "login">("server");
const savedServers = ref<ServerConfig[]>([]);
const activeServerUrl = ref<string | null>(null);
const showServerList = ref(false);
const serverListRef = ref<HTMLElement | null>(null);
const newServerName = ref("");
const newServerUrl = ref("");
const errorMsg = ref("");
const testing = ref(false);
const connecting = ref(false);
const testingUrl = ref<string | null>(null);
const toastMsg = ref("");
const toastType = ref<"success" | "error">("success");

const username = ref("");
const password = ref("");
const loggingIn = ref(false);

const currentServer = computed(() =>
  savedServers.value.find((s) => s.url === activeServerUrl.value),
);

const handleClickOutside = (event: MouseEvent) => {
  if (
    serverListRef.value &&
    !serverListRef.value.contains(event.target as Node)
  ) {
    showServerList.value = false;
  }
};

onMounted(async () => {
  savedServers.value = getServerUrls();
  activeServerUrl.value = getActiveServerUrl();
  document.addEventListener("click", handleClickOutside);

  const activeToken = getActiveServerToken();
  if (activeServerUrl.value && activeToken) {
    refreshApiBaseURL();
    const isValid = await authStore.checkAuthStatus();
    if (isValid) {
      router.push("/");
      return;
    }
  }
});

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
});

const showToast = (msg: string, type: "success" | "error" = "success") => {
  toastMsg.value = msg;
  toastType.value = type;
  setTimeout(() => {
    toastMsg.value = "";
  }, 3000);
};

const testServerUrl = async (url: string): Promise<boolean> => {
  const validation = validateServerUrl(url);
  if (!validation.valid) {
    errorMsg.value = validation.error || "";
    return false;
  }

  try {
    const response = await fetch(`${url}/api/v1/health`, {
      method: "GET",
      signal: AbortSignal.timeout(5000),
    });

    if (response.ok) {
      return true;
    } else {
      errorMsg.value = "服务器响应异常";
      return false;
    }
  } catch {
    errorMsg.value = "无法连接到服务器，请检查地址";
    return false;
  }
};

const testNewServer = async () => {
  errorMsg.value = "";
  testing.value = true;

  const success = await testServerUrl(newServerUrl.value);
  if (success) {
    showToast("连接成功", "success");
  }

  testing.value = false;
};

const testSavedServer = async (url: string) => {
  testingUrl.value = url;
  errorMsg.value = "";

  const success = await testServerUrl(url);
  if (success) {
    showToast("连接成功", "success");
  }

  testingUrl.value = null;
};

const handleConnectServer = async () => {
  const validation = validateServerUrl(newServerUrl.value);
  if (!validation.valid) {
    errorMsg.value = validation.error || "";
    return;
  }

  errorMsg.value = "";
  connecting.value = true;

  const success = await testServerUrl(newServerUrl.value);
  if (!success) {
    connecting.value = false;
    return;
  }

  const name =
    newServerName.value.trim() ||
    (savedServers.value.length === 0
      ? "默认"
      : `服务器${savedServers.value.length + 1}`);
  savedServers.value = addServerUrl(name, newServerUrl.value);
  activeServerUrl.value = newServerUrl.value;
  setActiveServerUrl(newServerUrl.value);
  refreshApiBaseURL();
  systemStore.clearCache();

  newServerName.value = "";
  newServerUrl.value = "";

  showToast("连接成功", "success");

  authStore.loadTokenFromServer();
  const isValid = await authStore.checkAuthStatus();
  if (isValid) {
    router.push("/");
  } else {
    step.value = "login";
  }

  connecting.value = false;
};

const selectServer = async (server: ServerConfig) => {
  showServerList.value = false;
  newServerName.value = server.name;
  newServerUrl.value = server.url;

  if (activeServerUrl.value !== server.url) {
    setActiveServerUrl(server.url);
    activeServerUrl.value = server.url;
    refreshApiBaseURL();
    systemStore.clearCache();
    authStore.loadTokenFromServer();

    const isValid = await authStore.checkAuthStatus();
    if (isValid) {
      router.push("/");
    } else {
      step.value = "login";
    }
  } else {
    const isValid = await authStore.checkAuthStatus();
    if (isValid) {
      router.push("/");
    } else {
      step.value = "login";
    }
  }
};

const deleteServer = (url: string) => {
  savedServers.value = removeServerUrl(url);
  activeServerUrl.value = getActiveServerUrl();
  if (savedServers.value.length === 0) {
    showServerList.value = false;
  }
};

const switchServer = () => {
  step.value = "server";
  errorMsg.value = "";
  username.value = "";
  password.value = "";
};

const handleLogin = async () => {
  if (!username.value.trim() || !password.value.trim()) {
    errorMsg.value = "请输入用户名和密码";
    return;
  }

  errorMsg.value = "";
  loggingIn.value = true;

  try {
    const response = await userApi.login({
      username: username.value,
      password: password.value,
    });

    if (response && response.access_token) {
      const token = response.access_token;
      const userInfoResponse = await userApi.verifyToken(token);

      if (userInfoResponse && userInfoResponse.code === 0) {
        authStore.login(token, userInfoResponse.data);
        showToast("登录成功！", "success");
        setTimeout(() => {
          router.push("/");
        }, 300);
      } else {
        errorMsg.value = "获取用户信息失败";
      }
    } else {
      errorMsg.value = "登录失败，请检查用户名和密码";
    }
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } } };
    errorMsg.value = err.response?.data?.message || "登录失败，请稍后重试";
  } finally {
    loggingIn.value = false;
  }

  if (errorMsg.value) {
    showToast(errorMsg.value, "error");
  }
};
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
