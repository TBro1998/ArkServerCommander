<template>
  <UContainer class="py-12">
    <main class="text-center space-y-12">
      <!-- 欢迎标题区域 -->
      <header class="space-y-6">
        <UAvatar 
          size="2xl" 
          icon="i-lucide-server" 
          :ui="{ background: 'bg-gradient-to-br from-blue-500 to-indigo-600' }"
          class="mx-auto mb-6"
        />
        <h1 class="text-5xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
          {{ $t('home.title') }}
        </h1>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
          {{ $t('home.subtitle') }}
        </p>
      </header>
      
      <!-- 镜像管理区域 -->
      <section class="space-y-6">
        <h2 class="text-2xl font-semibold text-gray-900">{{ $t('home.imageManagement') }}</h2>
        <div class="max-w-4xl mx-auto">
          <UCard>
            <template #header>
              <div class="flex items-center gap-3">
                <UAvatar
                  size="sm"
                  icon="i-lucide-download"
                  :ui="{ background: 'bg-green-100' }"
                />
                <h3 class="text-lg font-semibold text-gray-900">{{ $t('servers.dockerImages.title') }}</h3>
              </div>
            </template>
            
            <ImageStatus
              :imageStatus="imageStatus"
              @refresh="refreshImageStatus"
              @manual-download="handleManualDownload"
              @check-updates="handleCheckUpdates"
            />
            
            <!-- 镜像状态详细模态框 -->
            <ImageStatusModal
              :show="showImageModal"
              :imageStatus="imageStatus"
              @close="showImageModal = false"
              @refresh="refreshImageStatus"
              @manual-download="handleManualDownload"
              @check-updates="handleCheckUpdates"
            />
            
            <!-- 镜像更新确认模态框 -->
            <ImageUpdateConfirmModal
              :show="showUpdateModal"
              :imageName="selectedImageName"
              :affectedServers="affectedServers"
              @close="showUpdateModal = false"
              @confirm="handleUpdateConfirm"
            />
          </UCard>
        </div>
      </section>

      <!-- 功能卡片网格 -->
      <section class="space-y-8">
        <h2 class="text-2xl font-semibold text-gray-900">{{ $t('home.features') }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <!-- 服务器管理卡片 -->
          <NuxtLink to="/servers">
            <UCard class="h-full hover:shadow-xl transition-all duration-300 cursor-pointer group border-2 border-transparent hover:border-blue-200">
              <article class="text-center space-y-6 p-6">
                <UAvatar
                  size="lg"
                  icon="i-lucide-server"
                  :ui="{ background: 'bg-blue-100 group-hover:bg-blue-200' }"
                  class="mx-auto transition-colors"
                />
                <header>
                  <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ $t('home.serverManagement') }}</h3>
                  <p class="text-gray-600 leading-relaxed">{{ $t('home.serverManagementDesc') }}</p>
                </header>
                <footer>
                  <UButton
                    color="blue"
                    variant="ghost"
                    size="sm"
                    class="group-hover:bg-blue-50"
                  >
                    {{ $t('home.startManage') }}
                    <UIcon name="i-lucide-arrow-right" class="w-4 h-4 ml-1" />
                  </UButton>
                </footer>
              </article>
            </UCard>
          </NuxtLink>

          <!-- 占位卡片：玩家管理 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <article class="text-center space-y-6 p-6">
              <UAvatar
                size="lg"
                icon="i-lucide-users"
                :ui="{ background: 'bg-gray-100' }"
                class="mx-auto"
              />
              <header>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.playerManagement') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.playerManagementDesc') }}</p>
              </header>
              <footer>
                <UBadge color="gray" variant="soft" size="sm">
                  {{ $t('home.comingSoon') }}
                </UBadge>
              </footer>
            </article>
          </UCard>

          <!-- 占位卡片：日志监控 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <article class="text-center space-y-6 p-6">
              <UAvatar
                size="lg"
                icon="i-lucide-bar-chart-3"
                :ui="{ background: 'bg-gray-100' }"
                class="mx-auto"
              />
              <header>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.logMonitoring') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.logMonitoringDesc') }}</p>
              </header>
              <footer>
                <UBadge color="gray" variant="soft" size="sm">
                  {{ $t('home.comingSoon') }}
                </UBadge>
              </footer>
            </article>
          </UCard>
        </div>
      </section>
    </main>
  </UContainer>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import ImageStatus from '~/components/dockerimages/ImageStatus.vue'
import ImageStatusModal from '~/components/dockerimages/ImageStatusModal.vue'
import ImageUpdateConfirmModal from '~/components/dockerimages/ImageUpdateConfirmModal.vue'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()
const toast = useToast()

// 镜像状态管理
const imageStatus = ref({
  images: {},
  any_pulling: false,
  any_not_ready: true,
  can_create_server: true,
  can_start_server: false,
  overall_status: '检查中...',
  pulling_count: 0,
  total_images: 2
})

const showImageModal = ref(false)
const showUpdateModal = ref(false)
const selectedImageName = ref('')
const affectedServers = ref([])

// 获取镜像状态
const refreshImageStatus = async () => {
  try {
    const response = await $fetch('/api/servers/images/status', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (response.data) {
      imageStatus.value = response.data
    }
  } catch (error) {
    console.error('获取镜像状态失败:', error)
    toast.add({
      title: '错误',
      description: '获取镜像状态失败',
      color: 'red'
    })
  }
}

// 处理手动下载
const handleManualDownload = async () => {
  // 找到第一个未就绪的镜像
  const notReadyImage = Object.entries(imageStatus.value.images).find(([_, status]) => !status.ready)
  
  if (!notReadyImage) {
    toast.add({
      title: '提示',
      description: '所有镜像都已就绪',
      color: 'blue'
    })
    return
  }

  const imageName = notReadyImage[0]
  
  try {
    await $fetch('/api/servers/images/pull', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: {
        image_name: imageName
      }
    })
    
    toast.add({
      title: '成功',
      description: `镜像 ${imageName} 开始下载`,
      color: 'green'
    })
    
    // 开始轮询状态
    startPolling()
  } catch (error) {
    console.error('下载镜像失败:', error)
    toast.add({
      title: '错误',
      description: '下载镜像失败',
      color: 'red'
    })
  }
}

// 处理检查更新
const handleCheckUpdates = async () => {
  try {
    const response = await $fetch('/api/servers/images/check-updates', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    toast.add({
      title: '检查完成',
      description: '镜像更新检查完成',
      color: 'blue'
    })
  } catch (error) {
    console.error('检查更新失败:', error)
    toast.add({
      title: '错误',
      description: '检查更新失败',
      color: 'red'
    })
  }
}

// 处理更新确认
const handleUpdateConfirm = async (options) => {
  try {
    await $fetch('/api/servers/images/update', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: {
        image_name: options.imageName
      }
    })
    
    showUpdateModal.value = false
    
    toast.add({
      title: '成功',
      description: `镜像 ${options.imageName} 开始更新`,
      color: 'green'
    })
    
    startPolling()
  } catch (error) {
    console.error('更新镜像失败:', error)
    toast.add({
      title: '错误',
      description: '更新镜像失败',
      color: 'red'
    })
  }
}

// 轮询状态
let pollingInterval = null

const startPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
  }
  
  pollingInterval = setInterval(async () => {
    await refreshImageStatus()
    
    // 如果没有镜像在拉取中，停止轮询
    if (!imageStatus.value.any_pulling) {
      clearInterval(pollingInterval)
      pollingInterval = null
    }
  }, 2000)
}

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

onMounted(() => {
  authStore.initFromStorage()
  refreshImageStatus()
})

onUnmounted(() => {
  stopPolling()
})
</script>