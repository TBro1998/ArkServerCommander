import { defineStore } from 'pinia'
import type { Ref } from 'vue'

// 定义用户对象类型
interface User {
  id: number;
  username: string;
  is_admin: boolean;
  created_at: string;
}

// 定义认证状态类型
interface AuthState {
  token: string | null;
  user: User | null;
  isLoading: boolean;
}

// 定义登录凭据类型
interface Credentials {
  username?: string;
  password?: string;
}

// 定义API响应类型
interface AuthResponse {
  token: string;
  user: User;
  message: string;
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: null,
    user: null,
    isLoading: false
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.token
  },

  actions: {
    async checkInit(): Promise<boolean> {
      const config = useRuntimeConfig()
      try {
        const response = await $fetch<{ initialized: boolean }>(`${config.public.apiBase}/auth/check-init`)
        return response.initialized
      } catch (error) {
        console.error('检查初始化状态失败:', error)
        return false
      }
    },

    async init(credentials: Credentials): Promise<{ success: boolean; message: string }> {
      const config = useRuntimeConfig()
      this.isLoading = true
      try {
        const response = await $fetch<AuthResponse>(`${config.public.apiBase}/auth/init`, {
          method: 'POST',
          body: credentials
        })
        
        this.token = response.token
        this.user = response.user
        
        // 保存token到cookie
        const tokenCookie: Ref<string | null> = useCookie('auth-token', {
          httpOnly: false,
          secure: false,
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 7 // 7天
        })
        tokenCookie.value = response.token
        
        return { success: true, message: response.message }
      } catch (error: any) {
        console.error('初始化失败:', error)
        return { success: false, message: error.data?.error || '初始化失败' }
      } finally {
        this.isLoading = false
      }
    },

    async login(credentials: Credentials): Promise<{ success: boolean; message: string }> {
      const config = useRuntimeConfig()
      this.isLoading = true
      try {
        const response = await $fetch<AuthResponse>(`${config.public.apiBase}/auth/login`, {
          method: 'POST',
          body: credentials
        })
        
        this.token = response.token
        this.user = response.user
        
        // 保存token到cookie
        const tokenCookie: Ref<string | null> = useCookie('auth-token', {
          httpOnly: false,
          secure: false,
          sameSite: 'lax',
          maxAge: 60 * 60 * 24 * 7 // 7天
        })
        tokenCookie.value = response.token
        
        return { success: true, message: response.message }
      } catch (error: any) {
        console.error('登录失败:', error)
        return { success: false, message: error.data?.error || '登录失败' }
      } finally {
        this.isLoading = false
      }
    },

    async getProfile(): Promise<void> {
      const config = useRuntimeConfig()
      try {
        const response = await $fetch<{ user: User }>(`${config.public.apiBase}/profile`, {
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

    logout(): void {
      this.token = null
      this.user = null
      
      // 删除cookie
      const tokenCookie = useCookie('auth-token')
      tokenCookie.value = null
    },

    initFromStorage(): void {
      // 从cookie初始化token
      const tokenCookie = useCookie<string | null>('auth-token')
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