import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'
import svgLoader from 'vite-svg-loader'

export default defineConfig(() => {
  const ENV = process.env.VITE_ENV || 'dev'
  const TARGET_SERVER = process.env.VITE_TARGET_SERVER || ''
  const BASE_URL = process.env.VITE_BASE_URL || '/'

  return {
    base: BASE_URL,
    plugins: [vue(), vueDevTools(), tailwindcss(), svgLoader()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    build: {
      chunkSizeWarningLimit: 1600,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('node_modules/vue') || id.includes('node_modules/@vue')) {
              return 'vue-core'
            }
            if (id.includes('node_modules/axios')) {
              return 'axios'
            }
            if (id.includes('node_modules/leaflet')) {
              return 'map'
            }
            if (id.includes('node_modules/chinese-days')) {
              return 'utils'
            }
          },
        },
      },
    },
    server: {
      host: true,
      port: 5174,
      proxy:
        ENV === 'dev'
          ? {
              '/api': {
                target: TARGET_SERVER,
                changeOrigin: true,
                rewrite: path => path,
              },
            }
          : undefined,
    },
  }
})