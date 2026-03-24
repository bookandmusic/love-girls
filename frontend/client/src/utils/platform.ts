const SERVER_URLS_KEY = "serverUrls";
const ACTIVE_SERVER_KEY = "activeServerUrl";
const SERVER_TOKEN_PREFIX = "server_token_";

export interface ServerConfig {
  name: string;
  url: string;
}

const getServerUrlHash = (url: string): string => {
  let hash = 0;
  for (let i = 0; i < url.length; i++) {
    const char = url.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }
  return Math.abs(hash).toString(36);
};

export const getServerToken = (serverUrl: string): string | null => {
  const hash = getServerUrlHash(serverUrl);
  return localStorage.getItem(`${SERVER_TOKEN_PREFIX}${hash}`);
};

export const setServerToken = (serverUrl: string, token: string): void => {
  const hash = getServerUrlHash(serverUrl);
  localStorage.setItem(`${SERVER_TOKEN_PREFIX}${hash}`, token);
};

export const removeServerToken = (serverUrl: string): void => {
  const hash = getServerUrlHash(serverUrl);
  localStorage.removeItem(`${SERVER_TOKEN_PREFIX}${hash}`);
};

export const getActiveServerToken = (): string | null => {
  const activeUrl = getActiveServerUrl();
  if (!activeUrl) return null;
  return getServerToken(activeUrl);
};

export const getServerUrls = (): ServerConfig[] => {
  const data = localStorage.getItem(SERVER_URLS_KEY);
  if (!data) return [];
  try {
    return JSON.parse(data);
  } catch {
    return [];
  }
};

export const setServerUrls = (servers: ServerConfig[]): void => {
  localStorage.setItem(SERVER_URLS_KEY, JSON.stringify(servers));
};

export const addServerUrl = (name: string, url: string): ServerConfig[] => {
  const servers = getServerUrls();
  const exists = servers.some((s) => s.url === url);
  if (!exists) {
    servers.push({ name, url });
    setServerUrls(servers);
  }
  return servers;
};

export const removeServerUrl = (url: string): ServerConfig[] => {
  const servers = getServerUrls().filter((s) => s.url !== url);
  setServerUrls(servers);

  if (getActiveServerUrl() === url) {
    const firstServer = servers[0];
    if (firstServer) {
      setActiveServerUrl(firstServer.url);
    } else {
      clearActiveServerUrl();
    }
  }
  return servers;
};

export const getActiveServerUrl = (): string | null => {
  return localStorage.getItem(ACTIVE_SERVER_KEY);
};

export const setActiveServerUrl = (url: string): void => {
  localStorage.setItem(ACTIVE_SERVER_KEY, url);
};

export const clearActiveServerUrl = (): void => {
  localStorage.removeItem(ACTIVE_SERVER_KEY);
};

export const validateServerUrl = (
  url: string,
): { valid: boolean; error?: string } => {
  if (!url.trim()) {
    return { valid: false, error: "请输入服务器地址" };
  }

  try {
    const parsed = new URL(url);
    if (!["http:", "https:"].includes(parsed.protocol)) {
      return { valid: false, error: "请输入有效的 HTTP/HTTPS 地址" };
    }
    return { valid: true };
  } catch {
    return { valid: false, error: "地址格式不正确" };
  }
};
