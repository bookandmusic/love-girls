import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'
import vueDevTools from 'vite-plugin-vue-devtools'
import svgLoader from 'vite-svg-loader'

export default defineConfig(() => {
  return {
    plugins: [
      vue(),
      svgLoader(),
      vueDevTools(),
      tailwindcss(),
      VitePWA({
        registerType: 'autoUpdate',
        devOptions: {
          enabled: true,
        },
        workbox: {
          globPatterns: ['**/*.{js,css,html,wasm,png,jpg,jpeg,svg}'],
        },
        manifest: {
          name: 'Love Girl',
          short_name: 'LoveGirl',
          description: 'A romantic app for couples',
          theme_color: '#ff69b4',
          background_color: '#ffffff',
          display: 'standalone',
          icons: [
            {
              src: 'favicon.png',
              sizes: '192x192',
              type: 'image/png',
            },
            {
              src: 'favicon.png',
              sizes: '512x512',
              type: 'image/png',
            },
          ],
        },
      }),
    ],
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
            if (
              id.includes('node_modules/vue3-lottie') ||
              id.includes('node_modules/lottie-web') ||
              id.includes('node_modules/vue-easy-lightbox')
            ) {
              return 'ui-lib'
            }
            if (id.includes('node_modules/leaflet')) {
              return 'map'
            }
            if (id.includes('node_modules/p5')) {
              return 'p5'
            }
            if (id.includes('node_modules/chinese-days')) {
              return 'utils'
            }
          },
        },
      },
    },
  }
})