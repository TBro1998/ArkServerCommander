<template>
  <div class="w-full px-2 sm:px-4 py-4 sm:py-8">
    <div class="bg-white rounded-lg shadow-lg">
      <div class="p-4 sm:p-6 border-b border-gray-200">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <h1 class="text-xl sm:text-2xl font-bold text-gray-900">服务器管理</h1>
          <div class="flex flex-col sm:flex-row items-center gap-4">
            <!-- 镜像状态显示 -->
            <div v-if="imageStatus" class="flex items-center gap-2 text-sm">
              <div v-if="imageStatus.any_pulling" class="flex items-center gap-1 text-yellow-600">
                <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                镜像下载中...
              </div>
              <div v-else-if="!imageStatus.can_create_server" class="flex items-center gap-1 text-red-600">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                </svg>
                镜像未就绪
              </div>
              <div v-else class="flex items-center gap-1 text-green-600">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
                镜像就绪
              </div>
              <button
                @click="refreshImageStatus"
                class="text-blue-600 hover:text-blue-800 p-1"
                title="刷新镜像状态"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
              </button>
            </div>
            <button
              @click="showCreateForm = true"
              :disabled="!imageStatus?.can_create_server"
              class="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed text-white px-4 py-2 rounded-lg flex items-center justify-center gap-2 transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
              </svg>
              新增服务器
            </button>
          </div>
        </div>
      </div>

      <div class="p-4 sm:p-6">
        <!-- 消息提示 -->
        <ErrorMessage
          :show="!!errorMessage"
          :message="errorMessage"
          @close="errorMessage = ''"
        />
        <SuccessMessage
          :show="!!successMessage"
          :message="successMessage"
          @close="successMessage = ''"
        />
        
        <!-- 服务器列表 -->
        <div v-if="loading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p class="mt-2 text-gray-600">加载中...</p>
        </div>

        <div v-else-if="servers.length === 0" class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无服务器</h3>
          <p class="mt-1 text-sm text-gray-500">开始创建您的第一个ARK服务器</p>
          <div class="mt-6">
            <button
              @click="showCreateForm = true"
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
            >
              新增服务器
            </button>
          </div>
        </div>

        <div v-else class="flex flex-wrap gap-4 sm:gap-6 justify-center sm:justify-start">
          <div
            v-for="server in servers"
            :key="server.id"
            class="bg-gray-50 rounded-lg p-4 sm:p-6 hover:shadow-md transition-shadow h-fit w-80"
          >
            <div class="flex justify-between items-start mb-4">
              <div class="min-w-0 flex-1">
                <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">{{ server.identifier }}</h3>
              </div>
              <div class="flex gap-1 sm:gap-2 ml-2">
                <!-- 启动/停止按钮 -->
                <button
                  v-if="server.status === 'stopped'"
                  @click="startServer(server)"
                  class="text-green-600 hover:text-green-800 p-1"
                  title="启动服务器"
                >
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z"/>
                  </svg>
                </button>
                <button
                  v-else-if="server.status === 'running'"
                  @click="stopServer(server)"
                  class="text-red-600 hover:text-red-800 p-1"
                  title="停止服务器"
                >
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M6 6h4v12H6zm8-6v18h4V6h-4z"/>
                  </svg>
                </button>
                <button
                  v-else-if="server.status === 'starting'"
                  class="text-yellow-600 p-1 cursor-not-allowed"
                  disabled
                  title="启动中..."
                >
                  <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                  </svg>
                </button>
                <button
                  v-else-if="server.status === 'stopping'"
                  class="text-yellow-600 p-1 cursor-not-allowed"
                  disabled
                  title="停止中..."
                >
                  <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                  </svg>
                </button>
                <button
                  v-else
                  class="text-gray-400 p-1 cursor-not-allowed"
                  disabled
                  title="未知状态"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </button>
                
                <button
                  @click="showRCONInfo(server)"
                  class="text-green-600 hover:text-green-800 p-1"
                  title="RCON信息"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </button>
                <button
                  @click="editServer(server)"
                  class="text-blue-600 hover:text-blue-800 p-1"
                  title="编辑服务器"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                </button>
                <button
                  @click="confirmDelete(server)"
                  class="text-red-600 hover:text-red-800 p-1"
                  title="删除"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- 基本信息 -->
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600">状态:</span>
                <span
                  :class="{
                    'text-green-600': server.status === 'running',
                    'text-red-600': server.status === 'stopped',
                    'text-yellow-600': server.status === 'starting' || server.status === 'stopping'
                  }"
                  class="font-medium"
                >
                  {{ getStatusText(server.status) }}
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">地图:</span>
                <span class="font-medium">{{ server.map || 'TheIsland' }}</span>
              </div>
            </div>

            <!-- 端口信息 -->
            <div class="mt-4 mb-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0"></path>
                </svg>
                端口配置
              </h4>
              <div class="space-y-1 text-sm pl-5">
                <div class="flex justify-between">
                  <span class="text-gray-600">游戏端口:</span>
                  <span class="font-mono text-blue-600">{{ server.port }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600">查询端口:</span>
                  <span class="font-mono text-blue-600">{{ server.query_port }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600">RCON端口:</span>
                  <span class="font-mono text-blue-600">{{ server.rcon_port }}</span>
                </div>
              </div>
            </div>

            <!-- 认证信息 -->
            <div class="mb-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
                </svg>
                认证信息
              </h4>
              <div class="space-y-1 text-sm pl-5">
                <div class="flex justify-between items-center">
                  <span class="text-gray-600">管理员密码:</span>
                  <div class="flex items-center gap-2">
                    <span class="font-mono text-gray-800 text-xs sm:text-sm">{{ showServerPasswords[server.id] ? server.admin_password : '***' }}</span>
                    <button
                      @click="toggleServerPassword(server.id)"
                      class="text-xs text-gray-500 hover:text-gray-700"
                      :title="showServerPasswords[server.id] ? '隐藏密码' : '显示密码'"
                    >
                      {{ showServerPasswords[server.id] ? '隐藏' : '显示' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- 时间信息 -->
            <div class="mb-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                时间信息
              </h4>
              <div class="space-y-1 text-sm pl-5">
                <div class="flex justify-between">
                  <span class="text-gray-600">创建时间:</span>
                  <span class="text-gray-800 text-xs sm:text-sm">{{ formatDate(server.created_at) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600">更新时间:</span>
                  <span class="text-gray-800 text-xs sm:text-sm">{{ formatDate(server.updated_at) }}</span>
                </div>
              </div>
            </div>

            <!-- 服务器ID信息 -->
            <div class="text-xs text-gray-500 border-t pt-2 mt-2">
              <div class="flex justify-between">
                <span>服务器ID:</span>
                <span class="font-mono">#{{ server.id }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 统一的服务器编辑模态框 -->
    <ServerEditModal
      :show="showCreateForm || showEditForm"
      :mode="showCreateForm ? 'create' : 'edit'"
      :server="currentEditServer"
      :loading="loadingServerData"
      :saving="submitting"
      @close="closeForm"
      @save="handleServerSave"
    />


    <!-- 删除确认模态框 -->
    <div
      v-if="showDeleteConfirm"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showDeleteConfirm = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3 text-center">
          <svg class="mx-auto mb-4 w-12 h-12 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
          </svg>
          <h3 class="text-lg font-medium text-gray-900">确认删除</h3>
          <div class="mt-2 px-7 py-3">
            <p class="text-sm text-gray-500">
              您确定要删除服务器 "{{ serverToDelete?.identifier }}" 吗？此操作无法撤销。
            </p>
          </div>
          <div class="flex gap-3 mt-4">
            <button
              @click="deleteServer"
              :disabled="deleting"
              class="flex-1 bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-lg transition-colors disabled:opacity-50"
            >
              {{ deleting ? '删除中...' : '确认删除' }}
            </button>
            <button
              @click="showDeleteConfirm = false"
              class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-700 py-2 px-4 rounded-lg transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- RCON信息模态框 -->
    <div
      v-if="showRCONModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showRCONModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-full max-w-md shadow-lg rounded-md bg-white" @click.stop>
        <div class="mb-4">
          <h3 class="text-lg font-bold text-gray-900">RCON连接信息</h3>
        </div>
        
        <div v-if="rconInfo" class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">服务器标识</label>
            <p class="px-3 py-2 bg-gray-50 rounded border text-sm">{{ rconInfo.server_identifier }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">RCON端口</label>
            <div class="flex">
              <p class="flex-1 px-3 py-2 bg-gray-50 rounded-l border text-sm font-mono">{{ rconInfo.rcon_port }}</p>
              <button
                @click="copyToClipboard(rconInfo.rcon_port)"
                class="px-3 py-2 bg-blue-500 text-white rounded-r hover:bg-blue-600 text-sm"
                title="复制端口"
              >
                复制
              </button>
            </div>
          </div>
                                <div>
             <label class="block text-sm font-medium text-gray-700 mb-1">管理员密码（RCON密码）</label>
             <div class="flex">
               <input
                 :type="showRCONPassword ? 'text' : 'password'"
                 :value="rconInfo.admin_password"
                 readonly
                 class="flex-1 px-3 py-2 bg-gray-50 rounded-l border text-sm font-mono"
               />
               <button
                 @click="showRCONPassword = !showRCONPassword"
                 class="px-3 py-2 bg-gray-500 text-white hover:bg-gray-600 text-sm"
                 :title="showRCONPassword ? '隐藏密码' : '显示密码'"
               >
                 {{ showRCONPassword ? '隐藏' : '显示' }}
               </button>
               <button
                 @click="copyToClipboard(rconInfo.admin_password)"
                 class="px-3 py-2 bg-blue-500 text-white rounded-r hover:bg-blue-600 text-sm"
                 title="复制密码"
               >
                 复制
               </button>
             </div>
           </div>
        </div>
        
        <div class="flex justify-end mt-6">
          <button
            @click="showRCONModal = false"
            class="px-4 py-2 bg-gray-300 hover:bg-gray-400 text-gray-700 rounded-lg transition-colors"
          >
            关闭
          </button>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

// 使用 servers store
const serversStore = useServersStore()

// 响应式数据
const submitting = ref(false)
const deleting = ref(false)
const showCreateForm = ref(false)
const showEditForm = ref(false)
const showDeleteConfirm = ref(false)
const showRCONModal = ref(false)
const serverToDelete = ref(null)
const rconInfo = ref(null)
const showRCONPassword = ref(false)
const showServerPasswords = ref({})
const errorMessage = ref('')
const successMessage = ref('')

// 镜像状态相关
const imageStatus = ref(null)
const imageStatusInterval = ref(null)

// 统一编辑器相关
const currentEditServer = ref(null)
const loadingServerData = ref(false)

// 使用 store 中的数据
const servers = computed(() => serversStore.servers)
const loading = computed(() => serversStore.isLoading)

// 获取镜像状态
const fetchImageStatus = async () => {
  try {
    console.log('开始获取镜像状态...')
    const status = await serversStore.getImageStatus()
    imageStatus.value = status
    console.log('镜像状态获取成功:', status)
  } catch (error) {
    console.error('获取镜像状态失败:', error)
    // 不显示错误消息，因为镜像状态不是关键功能
  }
}

// 刷新镜像状态
const refreshImageStatus = async () => {
  await fetchImageStatus()
}

// 获取服务器列表
const fetchServers = async () => {
  try {
    console.log('开始获取服务器列表...')
    await serversStore.fetchServers()
    console.log('服务器列表获取成功，共', servers.value.length, '个服务器')
  } catch (error) {
    console.error('获取服务器列表失败:', error)
    
    // 检查是否是认证错误
    if (error.status === 401 || error.statusCode === 401) {
      console.log('认证失败，可能需要重新登录')
      errorMessage.value = '认证失败，请重新登录'
      // 清除认证状态并重定向到登录页
      const authStore = useAuthStore()
      authStore.logout()
      await navigateTo('/login')
    } else {
      errorMessage.value = error.data?.error || serversStore.error || '获取服务器列表失败，请稍后重试'
    }
  }
}

// 关闭表单
const closeForm = () => {
  showCreateForm.value = false
  showEditForm.value = false
  currentEditServer.value = null
  loadingServerData.value = false
}

// 编辑服务器（统一入口）
const editServer = async (server) => {
  currentEditServer.value = null
  loadingServerData.value = true
  showEditForm.value = true
  
  try {
    console.log('加载服务器信息:', server.id)
    const serverData = await serversStore.getServer(server.id)
    
    // 设置服务器数据
    currentEditServer.value = serverData
    
    console.log('服务器信息加载成功', serverData)
  } catch (error) {
    console.error('加载服务器信息失败:', error)
    errorMessage.value = '加载服务器信息失败，请稍后重试'
    closeForm()
  } finally {
    loadingServerData.value = false
  }
}

// 统一的服务器保存处理
const handleServerSave = async (formData) => {
  // 检查镜像状态（仅在创建服务器时）
  if (!showEditForm.value && imageStatus.value && !imageStatus.value.can_create_server) {
    errorMessage.value = '镜像未就绪，无法创建服务器，请等待镜像下载完成'
    return
  }
  
  submitting.value = true
  try {
    if (showEditForm.value && currentEditServer.value) {
      // 更新服务器
      await serversStore.updateServer(currentEditServer.value.id, formData)
      successMessage.value = '服务器更新成功'
    } else {
      // 创建服务器
      await serversStore.createServer(formData)
      successMessage.value = '服务器创建成功'
    }
    
    closeForm()
  } catch (error) {
    console.error('操作失败:', error)
    errorMessage.value = error.data?.error || serversStore.error || '操作失败，请稍后重试'
  } finally {
    submitting.value = false
  }
}

// 确认删除
const confirmDelete = (server) => {
  if (server.status === 'running') {
    alert('无法删除正在运行的服务器，请先停止服务器')
    return
  }
  serverToDelete.value = server
  showDeleteConfirm.value = true
}

// 删除服务器
const deleteServer = async () => {
  deleting.value = true
  try {
    await serversStore.deleteServer(serverToDelete.value.id)
    
    showDeleteConfirm.value = false
    serverToDelete.value = null
    successMessage.value = '服务器删除成功'
  } catch (error) {
    console.error('删除失败:', error)
    errorMessage.value = error.data?.error || serversStore.error || '删除失败，请稍后重试'
  } finally {
    deleting.value = false
  }
}

// 启动服务器
const startServer = async (server) => {
  try {
    await serversStore.startServer(server.id)
    successMessage.value = `服务器 "${server.identifier}" 启动中...`
  } catch (error) {
    console.error('启动服务器失败:', error)
    errorMessage.value = serversStore.error || '启动服务器失败，请稍后重试'
  }
}

// 停止服务器
const stopServer = async (server) => {
  try {
    await serversStore.stopServer(server.id)
    successMessage.value = `服务器 "${server.identifier}" 停止中...`
  } catch (error) {
    console.error('停止服务器失败:', error)
    errorMessage.value = serversStore.error || '停止服务器失败，请稍后重试'
  }
}

// 切换服务器密码显示
const toggleServerPassword = (serverId) => {
  showServerPasswords.value[serverId] = !showServerPasswords.value[serverId]
}



// 显示RCON信息
const showRCONInfo = async (server) => {
  try {
    console.log('获取服务器RCON信息:', server.id)
    const rconData = await serversStore.getServerRCON(server.id)
    rconInfo.value = rconData
    showRCONPassword.value = false // 默认隐藏密码
    showRCONModal.value = true
  } catch (error) {
    console.error('获取RCON信息失败:', error)
    errorMessage.value = '获取RCON信息失败，请稍后重试'
  }
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text.toString())
    successMessage.value = '已复制到剪贴板'
  } catch (error) {
    console.error('复制失败:', error)
    errorMessage.value = '复制失败，请手动复制'
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'running': return '运行中'
    case 'stopped': return '已停止'
    case 'starting': return '启动中'
    case 'stopping': return '停止中'
    default: return '未知'
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '暂无'
  try {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return dateString
  }
}





// 页面加载时获取数据
onMounted(() => {
  fetchServers()
  fetchImageStatus()
  
  // 设置定时刷新镜像状态（每5秒检查一次）
  imageStatusInterval.value = setInterval(() => {
    fetchImageStatus()
  }, 5000)
})

// 页面卸载时清理定时器
onUnmounted(() => {
  if (imageStatusInterval.value) {
    clearInterval(imageStatusInterval.value)
  }
})
</script> 