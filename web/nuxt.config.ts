// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  compatibilityDate: '2025-06-30',
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],
  runtimeConfig: {
    public: {
      // 开发环境连接本地后端，生产环境使用相对路径
      apiBase: '/api'
    }
  },
  ssr: false,
  nitro: {
    preset: 'static',
    // 开发环境代理配置
    devProxy: {
      '/api': {
        target: 'http://localhost:8080/api',
        changeOrigin: true
      }
    }
  },
  app: {
    baseURL: '/',
    buildAssetsDir: '/_nuxt/'
  },
  // 开发环境配置
  devServer: {
    port: 3000,
    host: 'localhost'
  }
}) 