import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

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
    // 获取认证token
    getAuthToken() {
      const authStore = useAuthStore()
      const tokenCookie = useCookie('auth-token')
      
      // 优先使用 auth store 中的 token，如果没有则使用 cookie
      const token = authStore.token || tokenCookie.value
      
      console.log('获取到的token:', token ? '存在' : '不存在')
      return token
    },

    // 获取认证头
    getAuthHeaders() {
      const token = this.getAuthToken()
      if (!token) {
        throw new Error('未找到认证token')
      }
      
      return {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    },

    // 执行API请求的通用方法
    async makeRequest(url, options = {}) {
      const config = useRuntimeConfig()
      const fullUrl = `${config.public.apiBase}${url}`
      
      const defaultOptions = {
        headers: this.getAuthHeaders(),
        ...options
      }

      // 如果有body数据，确保正确序列化
      if (defaultOptions.body && typeof defaultOptions.body === 'object') {
        defaultOptions.body = JSON.stringify(defaultOptions.body)
      }

      console.log('发起API请求:', fullUrl, defaultOptions)

      try {
        // 使用全局的$fetch函数
        const response = await $fetch(fullUrl, defaultOptions)
        console.log('API响应:', response)
        return response
      } catch (error) {
        console.error('API请求失败:', error)
        throw error
      }
    },

    // 获取服务器列表
    async fetchServers() {
      this.isLoading = true
      this.error = null
      
      try {
        console.log('正在请求服务器列表...')
        
        const response = await this.makeRequest('/servers')
        
        console.log('服务器列表响应:', response)
        
        this.servers = response.data || []
      } catch (error) {
        this.error = '获取服务器列表失败'
        console.error('获取服务器列表失败:', error)
        console.error('错误详情:', error.data)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // 创建服务器
    async createServer(serverData) {
      try {
        console.log('正在创建服务器:', serverData)
        
        const response = await this.makeRequest('/servers', {
          method: 'POST',
          body: serverData
        })
        
        console.log('创建服务器响应:', response)
        
        // 添加到本地状态
        this.servers.push(response.data)
        return response.data
      } catch (error) {
        this.error = '创建服务器失败'
        console.error('创建服务器失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 更新服务器
    async updateServer(serverId, updateData) {
      try {
        console.log('正在更新服务器:', serverId, updateData)
        
        const response = await this.makeRequest(`/servers/${serverId}`, {
          method: 'PUT',
          body: updateData
        })
        
        console.log('更新服务器响应:', response)
        
        // 更新本地状态
        const index = this.servers.findIndex(server => server.id === serverId)
        if (index !== -1) {
          this.servers[index] = response.data
        }
        
        return response.data
      } catch (error) {
        this.error = '更新服务器失败'
        console.error('更新服务器失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 删除服务器
    async deleteServer(serverId) {
      try {
        console.log('正在删除服务器:', serverId)
        
        await this.makeRequest(`/servers/${serverId}`, {
          method: 'DELETE'
        })
        
        console.log('服务器删除成功')
        
        // 从本地状态移除
        this.servers = this.servers.filter(server => server.id !== serverId)
      } catch (error) {
        this.error = '删除服务器失败'
        console.error('删除服务器失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 获取单个服务器信息
    async getServer(serverId) {
      try {
        console.log('正在获取服务器信息:', serverId)
        
        const response = await this.makeRequest(`/servers/${serverId}`)
        
        console.log('获取服务器信息响应:', response)
        
        return response.data
      } catch (error) {
        this.error = '获取服务器信息失败'
        console.error('获取服务器信息失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 获取服务器RCON信息
    async getServerRCON(serverId) {
      try {
        console.log('正在获取服务器RCON信息:', serverId)
        
        const response = await this.makeRequest(`/servers/${serverId}/rcon`)
        
        console.log('获取服务器RCON信息响应:', response)
        
        return response.data
      } catch (error) {
        this.error = '获取服务器RCON信息失败'
        console.error('获取服务器RCON信息失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 获取镜像状态
    async getImageStatus() {
      try {
        console.log('正在获取镜像状态...')
        
        const response = await this.makeRequest('/servers/images/status')
        
        console.log('获取镜像状态响应:', response)
        
        return response.data
      } catch (error) {
        this.error = '获取镜像状态失败'
        console.error('获取镜像状态失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 启动服务器
    async startServer(serverId) {
      try {
        console.log('正在启动服务器:', serverId)
        
        await this.makeRequest(`/servers/${serverId}/start`, {
          method: 'POST'
        })
        
        console.log('服务器启动请求发送成功')
        
        // 更新本地状态
        this.updateServerStatus(serverId, 'starting')
        
        // 3秒后更新为运行中（模拟启动过程）
        setTimeout(() => {
          this.updateServerStatus(serverId, 'running')
        }, 3000)
      } catch (error) {
        this.error = '启动服务器失败'
        console.error('启动服务器失败:', error)
        console.error('错误详情:', error.data)
        throw error
      }
    },

    // 停止服务器
    async stopServer(serverId) {
      try {
        console.log('正在停止服务器:', serverId)
        
        await this.makeRequest(`/servers/${serverId}/stop`, {
          method: 'POST'
        })
        
        console.log('服务器停止请求发送成功')
        
        // 更新本地状态
        this.updateServerStatus(serverId, 'stopping')
        
        // 2秒后更新为已停止（模拟停止过程）
        setTimeout(() => {
          this.updateServerStatus(serverId, 'stopped')
        }, 2000)
      } catch (error) {
        this.error = '停止服务器失败'
        console.error('停止服务器失败:', error)
        console.error('错误详情:', error.data)
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

    // 下载单个镜像
    async downloadImage(imageName) {
      try {
        console.log('开始下载镜像:', imageName)
        
        const response = await this.makeRequest('/servers/images/pull', {
          method: 'POST',
          body: { image_name: imageName }
        })
        
        console.log('镜像下载请求已发送:', response)
        return response
      } catch (error) {
        this.error = '下载镜像失败'
        console.error('下载镜像失败:', error)
        throw error
      }
    },

    // 更新单个镜像
    async updateImage(imageName) {
      try {
        console.log('开始更新镜像:', imageName)
        
        const response = await this.makeRequest('/servers/images/update', {
          method: 'POST',
          body: { image_name: imageName }
        })
        
        console.log('镜像更新请求已发送:', response)
        return response
      } catch (error) {
        this.error = '更新镜像失败'
        console.error('更新镜像失败:', error)
        throw error
      }
    },

    // 清除错误
    clearError() {
      this.error = null
    }
  }
}) 