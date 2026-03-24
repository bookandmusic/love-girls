import { defineStore } from "pinia";
import { ref } from "vue";

import { userApi, type UserInfo } from "@/services/userApi";
import {
  getActiveServerUrl,
  getActiveServerToken,
  setServerToken,
  removeServerToken,
} from "@/utils/platform";

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(null);
  const userInfo = ref<UserInfo | null>(null);
  const isAuthenticated = ref(false);

  const login = (newToken: string, userData: UserInfo) => {
    const serverUrl = getActiveServerUrl();
    if (!serverUrl) return;

    token.value = newToken;
    userInfo.value = userData;
    isAuthenticated.value = true;
    setServerToken(serverUrl, newToken);
  };

  const logout = () => {
    const serverUrl = getActiveServerUrl();
    if (serverUrl) {
      removeServerToken(serverUrl);
    }
    token.value = null;
    userInfo.value = null;
    isAuthenticated.value = false;
  };

  const checkAuthStatus = async () => {
    const storedToken = getActiveServerToken();
    if (!storedToken) {
      isAuthenticated.value = false;
      token.value = null;
      userInfo.value = null;
      return false;
    }

    try {
      const response = await userApi.verifyToken(storedToken);

      if (response && response.code === 0) {
        token.value = storedToken;
        userInfo.value = response.data;
        isAuthenticated.value = true;
        return true;
      } else {
        logout();
        return false;
      }
    } catch (error) {
      console.error("验证token时发生错误:", (error as Error).message);
      logout();
      return false;
    }
  };

  const loadTokenFromServer = () => {
    const storedToken = getActiveServerToken();
    if (storedToken) {
      token.value = storedToken;
    } else {
      token.value = null;
    }
  };

  return {
    token,
    userInfo,
    isAuthenticated,
    login,
    logout,
    checkAuthStatus,
    loadTokenFromServer,
  };
});
