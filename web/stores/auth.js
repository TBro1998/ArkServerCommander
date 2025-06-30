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
        
        // 保存token到cookie
        const tokenCookie = useCookie('auth-token', {
          httpOnly: false,
          secure: false,
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 7 // 7天
        })
        tokenCookie.value = response.token
        
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
        
        // 保存token到cookie
        const tokenCookie = useCookie('auth-token', {
          httpOnly: false,
          secure: false,
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 7 // 7天
        })
        tokenCookie.value = response.token
        
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
      
      // 删除cookie
      const tokenCookie = useCookie('auth-token')
      tokenCookie.value = null
    },

    initFromStorage() {
      // 从cookie初始化token
      const tokenCookie = useCookie('auth-token')
      console.log('初始化认证状态，cookie中的token:', tokenCookie.value ? '存在' : '不存在')
      
      if (tokenCookie.value) {
        this.token = tokenCookie.value
        console.log('已设置token到auth store')
        
        // 异步获取用户信息，不阻塞页面加载
        this.getProfile().catch(error => {
          console.error('获取用户信息失败，可能token已过期:', error)
          this.logout()
        })
      } else {
        console.log('未找到有效的认证token')
      }
    }
  }
}) 