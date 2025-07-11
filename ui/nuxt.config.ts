// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: false },

  modules: [
    '@pinia/nuxt',
    '@nuxt/ui',
    '@nuxt/eslint',
    '@nuxt/fonts',
    '@nuxt/icon',
    '@nuxtjs/i18n'
  ],

  css: ['~/assets/css/main.css'],

  fonts: {
    provider: 'bunny'
  },

  future: {
    compatibilityVersion: 4
  },

  compatibilityDate: '2024-11-27',
  i18n: {
    locales: [
      { code: 'en', language: 'en-US', file: 'en.ts' },
      { code: 'zh', language: 'zh-CN', file: 'zh.ts' }
    ],
    defaultLocale: 'en',
    lazy: true,
    strategy: 'prefix_except_default',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
      alwaysRedirect: true,
      fallbackLocale: 'en'
    },
    bundle: {
      optimizeTranslationDirective: false
    }
  },
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