<template>
  <div class="container mx-auto px-4 py-8">
    <div class="bg-white rounded-lg shadow-lg">
      <div class="p-6 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <h1 class="text-2xl font-bold text-gray-900">服务器管理</h1>
          <button
            @click="showCreateForm = true"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            新增服务器
          </button>
        </div>
      </div>

      <div class="p-6">
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

        <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="server in servers"
            :key="server.id"
            class="bg-gray-50 rounded-lg p-6 hover:shadow-md transition-shadow"
          >
            <div class="flex justify-between items-start mb-4">
              <div>
                <h3 class="text-lg font-semibold text-gray-900">{{ server.identifier }}</h3>
              </div>
              <div class="flex gap-2">
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
                  @click="editServerConfig(server)"
                  class="text-purple-600 hover:text-purple-800 p-1"
                  title="编辑配置文件"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  </svg>
                </button>
                <button
                  @click="editServer(server)"
                  class="text-blue-600 hover:text-blue-800 p-1"
                  title="编辑基本信息"
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
                <span class="text-gray-600">游戏端口:</span>
                <span class="font-mono">{{ server.port }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">地图:</span>
                <span>{{ server.map }}</span>
              </div>

              <div class="flex justify-between">
                <span class="text-gray-600">管理员密码:</span>
                <span class="font-mono">{{ showServerPasswords[server.id] ? server.admin_password : '***' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">创建时间:</span>
                <span>{{ server.created_at }}</span>
              </div>
            </div>

            <div class="mt-4 flex gap-2">
              <button
                v-if="server.status === 'stopped'"
                @click="startServer(server)"
                class="flex-1 bg-green-600 hover:bg-green-700 text-white py-2 rounded-lg text-sm transition-colors"
              >
                启动
              </button>
              <button
                v-else-if="server.status === 'running'"
                @click="stopServer(server)"
                class="flex-1 bg-red-600 hover:bg-red-700 text-white py-2 rounded-lg text-sm transition-colors"
              >
                停止
              </button>
              <button
                v-else-if="server.status === 'starting'"
                class="flex-1 bg-gray-400 text-white py-2 rounded-lg text-sm cursor-not-allowed"
                disabled
              >
                启动中...
              </button>
              <button
                v-else-if="server.status === 'stopping'"
                class="flex-1 bg-gray-400 text-white py-2 rounded-lg text-sm cursor-not-allowed"
                disabled
              >
                停止中...
              </button>
              <button
                v-else
                class="flex-1 bg-gray-400 text-white py-2 rounded-lg text-sm cursor-not-allowed"
                disabled
              >
                未知状态
              </button>
              <button
                @click="toggleServerPassword(server.id)"
                class="px-3 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg text-sm transition-colors"
                :title="showServerPasswords[server.id] ? '隐藏密码' : '显示密码'"
              >
                {{ showServerPasswords[server.id] ? '隐藏' : '显示' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑服务器表单模态框 -->
    <div
      v-if="showCreateForm || showEditForm"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="closeForm"
    >
      <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white" @click.stop>
        <div class="mb-4">
          <h3 class="text-lg font-bold text-gray-900">
            {{ showEditForm ? '编辑服务器' : '新增服务器' }}
          </h3>
        </div>

        <form @submit.prevent="submitForm" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">服务器标识 *</label>
              <input
                v-model="form.identifier"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="输入服务器标识"
              />
            </div>
            <div></div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">游戏端口 *</label>
              <input
                v-model.number="form.port"
                type="number"
                required
                min="1"
                max="65535"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="7777"
              />
            </div>
            <div></div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">查询端口 *</label>
              <input
                v-model.number="form.query_port"
                type="number"
                required
                min="1"
                max="65535"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="27015"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">RCON端口 *</label>
              <input
                v-model.number="form.rcon_port"
                type="number"
                required
                min="1"
                max="65535"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="27020"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">地图</label>
              <select
                v-model="form.map"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="TheIsland">The Island</option>
                <option value="TheCenter">The Center</option>
                <option value="ScorchedEarth_P">Scorched Earth</option>
                <option value="Aberration_P">Aberration</option>
                <option value="Extinction">Extinction</option>
                <option value="Valguero_P">Valguero</option>
                <option value="Genesis">Genesis</option>
                <option value="CrystalIsles">Crystal Isles</option>
                <option value="Genesis2">Genesis 2</option>
                <option value="LostIsland">Lost Island</option>
                <option value="Fjordur">Fjordur</option>
              </select>
            </div>
            <div></div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">管理员密码 *</label>
            <div class="relative">
              <input
                v-model="form.admin_password"
                :type="showFormPassword ? 'text' : 'password'"
                required
                class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="输入管理员密码（同时作为RCON密码）"
              />
              <button
                type="button"
                @click="toggleFormPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-700"
                :title="showFormPassword ? '隐藏密码' : '显示密码'"
              >
                <svg v-if="showFormPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
              </button>
            </div>
          </div>

          <div class="flex gap-3 pt-4">
            <button
              type="submit"
              :disabled="submitting"
              class="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors disabled:opacity-50"
            >
              {{ submitting ? '提交中...' : (showEditForm ? '更新服务器' : '创建服务器') }}
            </button>
            <button
              type="button"
              @click="closeForm"
              class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-700 py-2 px-4 rounded-lg transition-colors"
            >
              取消
            </button>
          </div>
        </form>
      </div>
    </div>

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

    <!-- 配置文件编辑组件 -->
    <ServerConfigEditor
      :show="showConfigModal"
      :server="currentConfigServer"
      :config-data="configData"
      :config-paths="configPaths"
      :loading="loadingConfig"
      :saving="savingConfig"
      @close="closeConfigModal"
      @save="saveConfig"
    />
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
const editingServer = ref(null)
const rconInfo = ref(null)
const showRCONPassword = ref(false)
const showServerPasswords = ref({})
const showFormPassword = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// 配置文件编辑相关
const showConfigModal = ref(false)
const currentConfigServer = ref(null)
const loadingConfig = ref(false)
const savingConfig = ref(false)
const configData = ref({
  game_user_settings: '',
  game_ini: ''
})
const configPaths = ref(null)

// 表单数据
const form = ref({
  identifier: '',
  port: 7777,
  query_port: 27015,
  rcon_port: 27020,
  admin_password: '',
  map: 'TheIsland'
})

// 使用 store 中的数据
const servers = computed(() => serversStore.servers)
const loading = computed(() => serversStore.isLoading)

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

// 重置表单
const resetForm = () => {
  form.value = {
    identifier: '',
    port: 7777,
    query_port: 27015,
    rcon_port: 27020,
    admin_password: '',
    map: 'TheIsland'
  }
}

// 关闭表单
const closeForm = () => {
  showCreateForm.value = false
  showEditForm.value = false
  editingServer.value = null
  showFormPassword.value = false
  resetForm()
}

// 编辑服务器
const editServer = (server) => {
  editingServer.value = server
  form.value = {
    identifier: server.identifier,
    port: server.port,
    query_port: server.query_port,
    rcon_port: server.rcon_port,
    admin_password: server.admin_password, // 显示当前密码
    map: server.map
  }
  showEditForm.value = true
}

// 提交表单
const submitForm = async () => {
  submitting.value = true
  try {
    if (showEditForm.value) {
      // 更新服务器
      await serversStore.updateServer(editingServer.value.id, form.value)
    } else {
      // 创建服务器
      await serversStore.createServer(form.value)
    }
    
    closeForm()
    successMessage.value = showEditForm.value ? '服务器更新成功' : '服务器创建成功'
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

// 切换表单密码显示
const toggleFormPassword = () => {
  showFormPassword.value = !showFormPassword.value
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

// 编辑服务器配置文件
const editServerConfig = async (server) => {
  currentConfigServer.value = server
  showConfigModal.value = true
  loadingConfig.value = true
  
  try {
    console.log('加载服务器配置文件:', server.id)
    const serverData = await serversStore.getServer(server.id)
    
    // 设置配置文件内容
    configData.value = {
      game_user_settings: serverData.game_user_settings || '',
      game_ini: serverData.game_ini || ''
    }
    
    // 设置配置文件路径
    configPaths.value = {
      game_user_settings_path: serverData.game_user_settings_path,
      game_ini_path: serverData.game_ini_path
    }
    
    console.log('配置文件加载成功', configData.value)
  } catch (error) {
    console.error('加载配置文件失败:', error)
    errorMessage.value = '加载配置文件失败，请稍后重试'
    showConfigModal.value = false
  } finally {
    loadingConfig.value = false
  }
}

// 关闭配置文件编辑模态框
const closeConfigModal = () => {
  showConfigModal.value = false
  currentConfigServer.value = null
  activeConfigTab.value = 'game_user_settings'
  configData.value = {
    game_user_settings: '',
    game_ini: ''
  }
  configPaths.value = null
}

// 保存配置文件
const saveConfig = async (updatedConfigData) => {
  if (!currentConfigServer.value) return
  
  savingConfig.value = true
  try {
    console.log('保存配置文件:', currentConfigServer.value.id, updatedConfigData)
    
    const updatedServer = await serversStore.updateServer(currentConfigServer.value.id, updatedConfigData)
    console.log('配置文件保存成功:', updatedServer)
    
    closeConfigModal()
    successMessage.value = '配置文件保存成功'
  } catch (error) {
    console.error('保存配置文件失败:', error)
    errorMessage.value = error.data?.error || '保存配置文件失败，请稍后重试'
  } finally {
    savingConfig.value = false
  }
}



// 页面加载时获取数据
onMounted(() => {
  fetchServers()
})
</script> 