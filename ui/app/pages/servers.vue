<template>
  <UContainer class="py-8">
    <!-- 页面头部 -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">服务器管理</h1>
          <p class="text-gray-600">管理您的 ARK 服务器，支持一键启动、停止和配置</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 w-full sm:w-auto">
          <!-- 镜像状态按钮 -->
          <UButton
            @click="showImageStatusModal = true"
            color="gray"
            variant="outline"
            icon="i-lucide-image"
          >
            镜像状态
          </UButton>
          
          <UButton
            @click="showCreateForm = true"
            :disabled="!imageStatus?.can_create_server"
            color="blue"
            icon="i-lucide-plus"
          >
            新增服务器
          </UButton>
        </div>
      </div>
      
      <!-- 镜像状态指示器 -->
      <div v-if="imageStatus && !imageStatus.can_create_server" class="mt-4">
        <UAlert
          :title="imageStatus.any_pulling ? '镜像下载中' : '镜像未就绪'"
          :description="imageStatus.any_pulling ? '正在下载镜像，请稍后创建服务器' : '镜像未就绪，无法创建服务器'"
          :color="imageStatus.any_pulling ? 'yellow' : 'red'"
          variant="soft"
          :icon="imageStatus.any_pulling ? 'i-lucide-loader-2' : 'i-lucide-alert-circle'"
        />
      </div>
    </div>

    <!-- 消息提示 -->
    <UAlert
      v-if="errorMessage"
      :title="errorMessage"
      color="red"
      variant="soft"
      icon="i-lucide-alert-circle"
      @close="errorMessage = ''"
    />
    
    <UAlert
      v-if="successMessage"
      :title="successMessage"
      color="green"
      variant="soft"
      icon="i-lucide-check-circle"
      @close="successMessage = ''"
    />

    <!-- 镜像状态弹窗 -->
    <ImageStatusModal
      :show="showImageStatusModal"
      :image-status="imageStatus"
      @close="showImageStatusModal = false"
      @refresh="refreshImageStatus"
    />

    <!-- 服务器列表 -->
    <div v-if="loading" class="text-center py-12">
      <UIcon name="i-lucide-loader-2" class="w-8 h-8 animate-spin text-blue-600 mx-auto mb-4" />
      <p class="text-gray-600">加载中...</p>
    </div>

    <div v-else-if="servers.length === 0" class="text-center py-12">
      <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        <UIcon name="i-lucide-server" class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">暂无服务器</h3>
      <p class="text-gray-500 mb-6">开始创建您的第一个ARK服务器</p>
      
      <div v-if="imageStatus?.can_create_server">
        <UButton
          @click="showCreateForm = true"
          color="blue"
          icon="i-lucide-plus"
        >
          新增服务器
        </UButton>
      </div>
      <div v-else-if="imageStatus && !imageStatus.can_create_server">
        <UAlert
          :title="imageStatus.any_pulling ? '镜像下载中' : '镜像未就绪'"
          :description="imageStatus.any_pulling ? '请等待镜像下载完成后创建服务器' : '镜像未就绪，无法创建服务器'"
          :color="imageStatus.any_pulling ? 'yellow' : 'red'"
          variant="soft"
          :icon="imageStatus.any_pulling ? 'i-lucide-loader-2' : 'i-lucide-alert-circle'"
        />
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <ServerCard
        v-for="server in servers"
        :key="server.id"
        :server="server"
        :can-start-server="imageStatus?.can_start_server ?? true"
        @start="startServer"
        @stop="stopServer"
        @edit="editServer"
        @delete="confirmDelete"
      />
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
  </UContainer>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useServersStore } from '~/stores/servers'

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
const errorMessage = ref('')
const successMessage = ref('')

// 镜像状态相关
const imageStatus = ref(null)
const imageStatusInterval = ref(null)
const showImageStatusModal = ref(false)
const imageStatusReady = ref(false) // 标记镜像是否已就绪

// 统一编辑器相关
const currentEditServer = ref(null)
const loadingServerData = ref(false)

// 使用 store 中的数据
const servers = computed(() => serversStore.servers)
const loading = computed(() => serversStore.isLoading)

// 计算属性：是否正在轮询镜像状态
const isPollingImageStatus = computed(() => {
  return !imageStatusReady.value && !!imageStatusInterval.value
})

// 获取镜像状态
const fetchImageStatus = async () => {
  try {
    console.log('开始获取镜像状态...')
    const status = await serversStore.getImageStatus()
    imageStatus.value = status
    console.log('镜像状态获取成功:', status)
    
    // 检查镜像是否都已就绪
    if (status && status.can_create_server && status.can_start_server) {
      if (!imageStatusReady.value) {
        console.log('镜像已全部就绪，停止轮询')
        imageStatusReady.value = true
        // 停止轮询
        if (imageStatusInterval.value) {
          clearInterval(imageStatusInterval.value)
          imageStatusInterval.value = null
        }
      }
    } else {
      // 镜像未就绪，重置标记
      imageStatusReady.value = false
    }

  } catch (error) {
    console.error('获取镜像状态失败:', error)
    // 不显示错误消息，因为镜像状态不是关键功能
  }
}

// 刷新镜像状态
const refreshImageStatus = async () => {
  console.log('手动刷新镜像状态')
  
  // 重置就绪状态
  imageStatusReady.value = false
  
  // 如果当前有轮询，先停止
  if (imageStatusInterval.value) {
    clearInterval(imageStatusInterval.value)
    imageStatusInterval.value = null
  }
  
  await fetchImageStatus()
  
  // 如果镜像未就绪，重新启动轮询
  if (!imageStatusReady.value) {
    console.log('镜像未就绪，重新启动轮询')
    imageStatusInterval.value = setInterval(() => {
      fetchImageStatus()
    }, 2000)
  }
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
const confirmDelete = async (server) => {
  if (server.status === 'running') {
    alert('无法删除正在运行的服务器，请先停止服务器')
    return
  }
  
  deleting.value = true
  try {
    await serversStore.deleteServer(server.id)
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
  // 检查镜像状态
  if (imageStatus.value && !imageStatus.value.can_start_server) {
    errorMessage.value = '镜像未就绪，无法启动服务器，请等待镜像下载完成'
    return
  }
  
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

// 页面加载时获取数据
onMounted(async () => {
  fetchServers()
  await fetchImageStatus()
  
  // 只有在镜像未就绪时才启动轮询
  if (!imageStatusReady.value) {
    console.log('镜像未就绪，启动轮询')
    imageStatusInterval.value = setInterval(() => {
      fetchImageStatus()
    }, 2000)
  } else {
    console.log('镜像已就绪，无需轮询')
  }
})

// 页面卸载时清理定时器
onUnmounted(() => {
  if (imageStatusInterval.value) {
    console.log('清理镜像状态轮询定时器')
    clearInterval(imageStatusInterval.value)
    imageStatusInterval.value = null
  }
})
</script> 