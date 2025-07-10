<template>
  <UContainer class="py-8">
    <!-- 页面头部 -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ $t('servers.title') }}</h1>
          <p class="text-gray-600">{{ $t('servers.serverManagementDesc') }}</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 w-full sm:w-auto">
          <!-- 镜像状态按钮 -->
          <UButton
            @click="showImageStatusModal = true"
            color="gray"
            variant="outline"
            icon="i-lucide-image"
          >
            {{ $t('servers.imageStatus') }}
          </UButton>
          
          <UButton
            @click="showCreateForm = true"
            :disabled="!imageStatus?.can_create_server"
            color="blue"
            icon="i-lucide-plus"
          >
            {{ $t('servers.addServer') }}
          </UButton>
        </div>
      </div>
      
      <!-- 镜像状态指示器 -->
      <div v-if="imageStatus && !imageStatus.can_create_server" class="mt-4">
        <UAlert
          :title="imageStatus.any_pulling ? $t('servers.imageDownloading') : $t('servers.imageNotReady')"
          :description="imageStatus.any_pulling ? $t('servers.imageDownloadingDesc') : $t('servers.imageNotReadyDesc')"
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
      <p class="text-gray-600">{{ $t('common.loading') }}</p>
    </div>

    <div v-else-if="servers.length === 0" class="text-center py-12">
      <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        <UIcon name="i-lucide-server" class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">{{ $t('servers.noServers') }}</h3>
      <p class="text-gray-500 mb-6">{{ $t('servers.noServersDesc') }}</p>
      
      <div v-if="imageStatus?.can_create_server">
        <UButton
          @click="showCreateForm = true"
          color="blue"
          icon="i-lucide-plus"
        >
          {{ $t('servers.addServer') }}
        </UButton>
      </div>
      <div v-else-if="imageStatus && !imageStatus.can_create_server">
        <UAlert
          :title="imageStatus.any_pulling ? $t('servers.imageDownloading') : $t('servers.imageNotReady')"
          :description="imageStatus.any_pulling ? $t('servers.imageDownloadingDesc') : $t('servers.imageNotReadyDesc')"
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
    console.error($t('servers.imageStatusError'), error)
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
    console.error($t('servers.getServerListFailed'), error)
    
    // 检查是否是认证错误
    if (error.status === 401 || error.statusCode === 401) {
      console.log('认证失败，可能需要重新登录')
      errorMessage.value = $t('servers.authenticationFailed')
      // 清除认证状态并重定向到登录页
      const authStore = useAuthStore()
      authStore.logout()
      await navigateTo('/login')
    } else {
      errorMessage.value = error.data?.error || serversStore.error || $t('servers.getServerListFailed')
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
    console.error($t('servers.loadServerInfoFailed'), error)
    errorMessage.value = $t('servers.loadServerInfoFailed')
    closeForm()
  } finally {
    loadingServerData.value = false
  }
}

// 统一的服务器保存处理
const handleServerSave = async (formData) => {
  // 检查镜像状态（仅在创建服务器时）
  if (!showEditForm.value && imageStatus.value && !imageStatus.value.can_create_server) {
    errorMessage.value = $t('servers.imageNotReadyDesc')
    return
  }
  
  submitting.value = true
  try {
    if (showEditForm.value && currentEditServer.value) {
      // 更新服务器
      await serversStore.updateServer(currentEditServer.value.id, formData)
      successMessage.value = $t('servers.serverUpdateSuccess')
    } else {
      // 创建服务器
      await serversStore.createServer(formData)
      successMessage.value = $t('servers.serverCreateSuccess')
    }
    
    closeForm()
  } catch (error) {
    console.error($t('servers.operationFailed'), error)
    errorMessage.value = error.data?.error || serversStore.error || $t('servers.operationFailed')
  } finally {
    submitting.value = false
  }
}

// 确认删除
const confirmDelete = async (server) => {
  if (server.status === 'running') {
    alert($t('servers.cannotDeleteRunning'))
    return
  }
  
  deleting.value = true
  try {
    await serversStore.deleteServer(server.id)
    successMessage.value = $t('servers.serverDeleteSuccess')
  } catch (error) {
    console.error($t('servers.deleteFailed'), error)
    errorMessage.value = error.data?.error || serversStore.error || $t('servers.deleteFailed')
  } finally {
    deleting.value = false
  }
}

// 启动服务器
const startServer = async (server) => {
  // 检查镜像状态
  if (imageStatus.value && !imageStatus.value.can_start_server) {
    errorMessage.value = $t('servers.imageNotReadyDesc')
    return
  }
  
  try {
    await serversStore.startServer(server.id)
    successMessage.value = $t('servers.serverStartInProgress')
  } catch (error) {
    console.error($t('servers.startServerFailed'), error)
    errorMessage.value = serversStore.error || $t('servers.startServerFailed')
  }
}

// 停止服务器
const stopServer = async (server) => {
  try {
    await serversStore.stopServer(server.id)
    successMessage.value = $t('servers.serverStopInProgress')
  } catch (error) {
    console.error($t('servers.stopServerFailed'), error)
    errorMessage.value = serversStore.error || $t('servers.stopServerFailed')
  }
}

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text.toString())
    successMessage.value = $t('servers.copyToClipboard')
  } catch (error) {
    console.error($t('servers.copyFailed'), error)
    errorMessage.value = $t('servers.copyFailed')
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