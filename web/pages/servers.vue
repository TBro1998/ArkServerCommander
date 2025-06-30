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
                <h3 class="text-lg font-semibold text-gray-900">{{ server.name }}</h3>
                <p class="text-sm text-gray-600 mt-1">{{ server.description || '无描述' }}</p>
              </div>
              <div class="flex gap-2">
                <button
                  @click="editServer(server)"
                  class="text-blue-600 hover:text-blue-800 p-1"
                  title="编辑"
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
                    'text-yellow-600': server.status === 'starting'
                  }"
                  class="font-medium"
                >
                  {{ getStatusText(server.status) }}
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">地址:</span>
                <span class="font-mono">{{ server.ip }}:{{ server.port }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">地图:</span>
                <span>{{ server.map }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">最大玩家:</span>
                <span>{{ server.max_players }}</span>
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
                v-else
                class="flex-1 bg-gray-400 text-white py-2 rounded-lg text-sm cursor-not-allowed"
                disabled
              >
                启动中...
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
              <label class="block text-sm font-medium text-gray-700 mb-1">服务器名称 *</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="输入服务器名称"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">服务器描述</label>
              <input
                v-model="form.description"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="输入服务器描述"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">服务器IP *</label>
              <input
                v-model="form.ip"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="127.0.0.1"
              />
            </div>
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
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">最大玩家数 *</label>
              <input
                v-model.number="form.max_players"
                type="number"
                required
                min="1"
                max="200"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="70"
              />
            </div>
          </div>

          <div v-if="!showEditForm">
            <label class="block text-sm font-medium text-gray-700 mb-1">RCON密码 *</label>
            <input
              v-model="form.rcon_password"
              type="password"
              :required="!showEditForm"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="输入RCON密码"
            />
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
              您确定要删除服务器 "{{ serverToDelete?.name }}" 吗？此操作无法撤销。
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
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const { $fetch } = useNuxtApp()
const config = useRuntimeConfig()

// 响应式数据
const servers = ref([])
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const showCreateForm = ref(false)
const showEditForm = ref(false)
const showDeleteConfirm = ref(false)
const serverToDelete = ref(null)
const editingServer = ref(null)
const errorMessage = ref('')
const successMessage = ref('')

// 表单数据
const form = ref({
  name: '',
  description: '',
  ip: '127.0.0.1',
  port: 7777,
  query_port: 27015,
  rcon_port: 27020,
  rcon_password: '',
  map: 'TheIsland',
  max_players: 70
})

// 获取服务器列表
const fetchServers = async () => {
  loading.value = true
  try {
    const response = await $fetch(`${config.public.apiBase}/servers`, {
      headers: {
        Authorization: `Bearer ${useCookie('auth-token').value}`
      }
    })
    servers.value = response.data || []
  } catch (error) {
    console.error('获取服务器列表失败:', error)
    errorMessage.value = '获取服务器列表失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  form.value = {
    name: '',
    description: '',
    ip: '127.0.0.1',
    port: 7777,
    query_port: 27015,
    rcon_port: 27020,
    rcon_password: '',
    map: 'TheIsland',
    max_players: 70
  }
}

// 关闭表单
const closeForm = () => {
  showCreateForm.value = false
  showEditForm.value = false
  editingServer.value = null
  resetForm()
}

// 编辑服务器
const editServer = (server) => {
  editingServer.value = server
  form.value = {
    name: server.name,
    description: server.description,
    ip: server.ip,
    port: server.port,
    query_port: server.query_port,
    rcon_port: server.rcon_port,
    rcon_password: '', // 编辑时不显示密码
    map: server.map,
    max_players: server.max_players
  }
  showEditForm.value = true
}

// 提交表单
const submitForm = async () => {
  submitting.value = true
  try {
    if (showEditForm.value) {
      // 更新服务器
      const updateData = { ...form.value }
      delete updateData.rcon_password // 编辑时不更新密码
      
      await $fetch(`${config.public.apiBase}/servers/${editingServer.value.id}`, {
        method: 'PUT',
        body: updateData,
        headers: {
          Authorization: `Bearer ${useCookie('auth-token').value}`
        }
      })
    } else {
      // 创建服务器
      await $fetch(`${config.public.apiBase}/servers`, {
        method: 'POST',
        body: form.value,
        headers: {
          Authorization: `Bearer ${useCookie('auth-token').value}`
        }
      })
    }
    
    closeForm()
    await fetchServers()
    successMessage.value = showEditForm.value ? '服务器更新成功' : '服务器创建成功'
  } catch (error) {
    console.error('操作失败:', error)
    errorMessage.value = error.data?.error || '操作失败，请稍后重试'
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
    await $fetch(`${config.public.apiBase}/servers/${serverToDelete.value.id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${useCookie('auth-token').value}`
      }
    })
    
    showDeleteConfirm.value = false
    serverToDelete.value = null
    await fetchServers()
    successMessage.value = '服务器删除成功'
  } catch (error) {
    console.error('删除失败:', error)
    errorMessage.value = error.data?.error || '删除失败，请稍后重试'
  } finally {
    deleting.value = false
  }
}

// 启动服务器（占位函数）
const startServer = (server) => {
  console.log('启动服务器:', server.name)
  // TODO: 实现启动服务器功能
}

// 停止服务器（占位函数）
const stopServer = (server) => {
  console.log('停止服务器:', server.name)
  // TODO: 实现停止服务器功能
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

// 页面加载时获取数据
onMounted(() => {
  fetchServers()
})
</script> 