import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    user: null,
    isLoading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token
  },

  actions: {
    async checkInit() {
      const config = useRuntimeConfig()
      try {
        const response = await $fetch(`${config.public.apiBase}/auth/check-init`)
        return response.initialized
      } catch (error) {
        console.error('检查初始化状态失败:', error)
        return false
      }
    },

    async init(credentials) {
      const config = useRuntimeConfig()
      this.isLoading = true
      try {
        const response = await $fetch(`${config.public.apiBase}/auth/init`, {
          method: 'POST',
          body: credentials
        })
        
        this.token = response.token
        this.user = response.user
        
        // 保存token到localStorage
        if (process.client) {
          localStorage.setItem('token', response.token)
        }
        
        return { success: true, message: response.message }
      } catch (error) {
        console.error('初始化失败:', error)
        return { success: false, message: error.data?.error || '初始化失败' }
      } finally {
        this.isLoading = false
      }
    },

    async login(credentials) {
      const config = useRuntimeConfig()
      this.isLoading = true
      try {
        const response = await $fetch(`${config.public.apiBase}/auth/login`, {
          method: 'POST',
          body: credentials
        })
        
        this.token = response.token
        this.user = response.user
        
        // 保存token到localStorage
        if (process.client) {
          localStorage.setItem('token', response.token)
        }
        
        return { success: true, message: response.message }
      } catch (error) {
        console.error('登录失败:', error)
        return { success: false, message: error.data?.error || '登录失败' }
      } finally {
        this.isLoading = false
      }
    },

    async getProfile() {
      const config = useRuntimeConfig()
      try {
        const response = await $fetch(`${config.public.apiBase}/profile`, {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })
        this.user = response.user
      } catch (error) {
        console.error('获取用户信息失败:', error)
        this.logout()
      }
    },

    logout() {
      this.token = null
      this.user = null
      if (process.client) {
        localStorage.removeItem('token')
      }
    },

    initFromStorage() {
      if (process.client) {
        const token = localStorage.getItem('token')
        if (token) {
          this.token = token
          this.getProfile()
        }
      }
    }
  }
}) 