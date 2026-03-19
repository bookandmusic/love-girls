/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_ENV: string
  readonly VITE_DESKTOP_MODE?: string
  readonly VITE_API_BASE?: string
  readonly VITE_PWA_ENABLED?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare const __DESKTOP_MODE__: boolean
