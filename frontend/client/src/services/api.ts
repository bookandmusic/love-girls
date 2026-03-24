import axios from "axios";

import router from "@/router";
import { useAuthStore } from "@/stores/auth";
import { getActiveServerUrl, getActiveServerToken } from "@/utils/platform";

const getBaseURL = (): string => {
  const serverUrl = getActiveServerUrl();
  return serverUrl ? `${serverUrl}/api/v1` : "/api/v1";
};

const api = axios.create({
  baseURL: getBaseURL(),
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  (config) => {
    const token = getActiveServerToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error("API Error:", error);

    if (error.response?.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();

      if (router.currentRoute.value.name !== "ServerConfig") {
        router.push({ name: "ServerConfig" });
      }
    }

    return Promise.reject(error);
  },
);

export const refreshApiBaseURL = (): void => {
  api.defaults.baseURL = getBaseURL();
};

export default api;
