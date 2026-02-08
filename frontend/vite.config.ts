import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import { viteMockServe } from 'vite-plugin-mock'
import { VitePWA } from 'vite-plugin-pwa'
import vueDevTools from 'vite-plugin-vue-devtools'
import svgLoader from 'vite-svg-loader'

// 读取环境变量
const ENV = process.env.VITE_ENV || 'dev'
const TARGET_SERVER = process.env.VITE_TARGET_SERVER || ''

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    svgLoader(),
    vueDevTools(),
    tailwindcss(),
    VitePWA({
      registerType: 'autoUpdate',
      devOptions: {
        enabled: true
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,wasm,png,jpg,jpeg,svg}']
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
            type: 'image/png'
          },
          {
            src: 'favicon.png',
            sizes: '512x512',
            type: 'image/png'
          }
        ]
      }
    }),
    // 添加 mock 插件，仅在 mock 模式下启用
    ENV === 'mock' ? viteMockServe({
      mockPath: 'src/mock',  // mock 文件目录
      enable: true,          // 是否启用
      logger: true,          // 是否在控制台显示请求日志
    }) : undefined,
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  build: {
    chunkSizeWarningLimit: 1600,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/vue') || id.includes('node_modules/@vue')) {
            return 'vue-core';
          }
          if (id.includes('node_modules/axios')) {
            return 'axios';
          }
          if (id.includes('node_modules/vue3-lottie') ||
            id.includes('node_modules/lottie-web') ||
            id.includes('node_modules/vue-easy-lightbox')) {
            return 'ui-lib';
          }
          if (id.includes('node_modules/leaflet')) {
            return 'map';
          }
          if (id.includes('node_modules/p5')) {
            return 'p5';
          }
          if (id.includes('node_modules/chinese-days')) {
            return 'utils';
          }
        }
      }
    }
  },
  server: {
    allowedHosts: ['test.lw1314.site', 'free-pine-fa8b.tunnl.gg'],
    host: true, // 监听 0.0.0.0
    port: 5173, // 可选，指定端口
    proxy: ENV === 'dev' ? {
      '/api': {
        target: TARGET_SERVER,
        changeOrigin: true,
        rewrite: path => path, // 保持 /api/xxx
      }
    } : undefined
  },
})