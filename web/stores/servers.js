export const useServersStore = defineStore('servers', {
  state: () => ({
    servers: [],
    isLoading: false,
    error: null
  }),

  getters: {
    serverCount: (state) => state.servers.length,
    runningServers: (state) => state.servers.filter(server => server.status === 'running'),
    stoppedServers: (state) => state.servers.filter(server => server.status === 'stopped'),
    getServerById: (state) => (id) => state.servers.find(server => server.id === id)
  },

  actions: {
    // 获取服务器列表
    async fetchServers() {
      this.isLoading = true
      this.error = null
      
      try {
        const { $fetch } = useNuxtApp()
        const config = useRuntimeConfig()
        const response = await $fetch(`${config.public.apiBase}/servers`, {
          headers: {
            Authorization: `Bearer ${useCookie('auth-token').value}`
          }
        })
        
        this.servers = response.data || []
      } catch (error) {
        this.error = '获取服务器列表失败'
        console.error('获取服务器列表失败:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // 创建服务器
    async createServer(serverData) {
      try {
        const { $fetch } = useNuxtApp()
        const config = useRuntimeConfig()
        const response = await $fetch(`${config.public.apiBase}/servers`, {
          method: 'POST',
          body: serverData,
          headers: {
            Authorization: `Bearer ${useCookie('auth-token').value}`
          }
        })
        
        // 添加到本地状态
        this.servers.push(response.data)
        return response.data
      } catch (error) {
        this.error = '创建服务器失败'
        console.error('创建服务器失败:', error)
        throw error
      }
    },

    // 更新服务器
    async updateServer(serverId, updateData) {
      try {
        const { $fetch } = useNuxtApp()
        const config = useRuntimeConfig()
        const response = await $fetch(`${config.public.apiBase}/servers/${serverId}`, {
          method: 'PUT',
          body: updateData,
          headers: {
            Authorization: `Bearer ${useCookie('auth-token').value}`
          }
        })
        
        // 更新本地状态
        const index = this.servers.findIndex(server => server.id === serverId)
        if (index !== -1) {
          this.servers[index] = response.data
        }
        
        return response.data
      } catch (error) {
        this.error = '更新服务器失败'
        console.error('更新服务器失败:', error)
        throw error
      }
    },

    // 删除服务器
    async deleteServer(serverId) {
      try {
        const { $fetch } = useNuxtApp()
        const config = useRuntimeConfig()
        await $fetch(`${config.public.apiBase}/servers/${serverId}`, {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${useCookie('auth-token').value}`
          }
        })
        
        // 从本地状态移除
        this.servers = this.servers.filter(server => server.id !== serverId)
      } catch (error) {
        this.error = '删除服务器失败'
        console.error('删除服务器失败:', error)
        throw error
      }
    },

    // 获取单个服务器信息
    async getServer(serverId) {
      try {
        const { $fetch } = useNuxtApp()
        const config = useRuntimeConfig()
        const response = await $fetch(`${config.public.apiBase}/servers/${serverId}`, {
          headers: {
            Authorization: `Bearer ${useCookie('auth-token').value}`
          }
        })
        
        return response.data
      } catch (error) {
        this.error = '获取服务器信息失败'
        console.error('获取服务器信息失败:', error)
        throw error
      }
    },

    // 更新服务器状态
    updateServerStatus(serverId, status) {
      const server = this.servers.find(s => s.id === serverId)
      if (server) {
        server.status = status
      }
    },

    // 清除错误
    clearError() {
      this.error = null
    }
  }
}) 