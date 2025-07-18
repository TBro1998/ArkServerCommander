<template>
  <UModal v-model:open="isOpen" :title="$t('servers.dockerImages.title')">
    <template #body>
      <!-- 总体状态 -->
      <div class="flex items-center gap-2 mb-6">
        <div v-if="imageStatus.any_pulling" class="flex items-center gap-2 text-yellow-600">
          <svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          <span class="font-medium">{{ imageStatus.overall_status }}</span>
        </div>
        <div v-else-if="!imageStatus.can_start_server" class="flex items-center gap-2 text-red-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
          </svg>
          <span class="font-medium">{{ $t('servers.dockerImages.imageMissingManualDownload') }}</span>
        </div>
        <div v-else class="flex items-center gap-2 text-green-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <span class="font-medium">{{ $t('servers.dockerImages.imageReady') }}</span>
        </div>
        
        <button
          @click="$emit('refresh')"
          class="ml-auto text-blue-600 hover:text-blue-800 p-2 transition-colors rounded-lg hover:bg-blue-50"
          :title="$t('servers.dockerImages.refreshStatus')"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
        </button>
      </div>

      <!-- 镜像横向排列 -->
      <div v-if="imageStatus.images" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div 
          v-for="(status, imageName) in imageStatus.images" 
          :key="imageName"
          class="bg-gray-50 p-4 rounded-lg border-l-4"
          :class="{
            'border-l-green-500': status.ready,
            'border-l-yellow-500': status.pulling,
            'border-l-red-500': !status.ready && !status.pulling
          }"
        >
          <!-- 镜像标题和状态 -->
          <div class="flex justify-between items-center mb-3">
            <h3 class="text-sm font-semibold text-gray-900">{{ getImageDisplayName(imageName) }}</h3>
            <span 
              :class="{
                'text-green-600 bg-green-100': status.ready,
                'text-yellow-600 bg-yellow-100': status.pulling,
                'text-red-600 bg-red-100': !status.ready && !status.pulling
              }"
              class="text-sm font-medium px-2 py-1 rounded-full"
            >
              {{ getStatusText(status) }}
            </span>
          </div>

          <!-- 层级信息 -->
          <div v-if="status.pulling && status.layers" class="space-y-2">
            <div class="text-xs font-medium text-gray-700 mb-2">{{ $t('servers.dockerImages.layerProgress') }}:</div>
            <div 
              v-for="(layer, layerId) in status.layers" 
              :key="layerId"
              class="bg-white p-3 rounded border"
            >
              <div class="flex justify-between items-center mb-2">
                <span class="text-gray-600 font-mono text-xs">
                  {{ layer.id.substring(0, 12) }}...
                </span>
                <span 
                  :class="{
                    'text-blue-600': layer.status === 'downloading',
                    'text-yellow-600': layer.status === 'extracting',
                    'text-purple-600': layer.status === 'verifying',
                    'text-green-600': layer.status === 'complete',
                    'text-gray-500': layer.status === 'pending'
                  }"
                  class="text-xs font-medium"
                >
                  {{ $t(`servers.dockerImages.layerStatus.${layer.status}`, layer.status) }}
                </span>
              </div>
              
              <div v-if="layer.size > 0" class="text-xs text-gray-500 mb-2">
                {{ formatBytes(layer.progress) }} / {{ formatBytes(layer.size) }}
              </div>
              <div v-else class="text-xs text-gray-500 mb-2">
                {{ formatBytes(layer.progress) }} / {{ $t('servers.dockerImages.unknownSize') }}
              </div>
              
              <!-- 层级进度条 -->
              <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div 
                  :class="{
                    'bg-blue-500': layer.status === 'downloading',
                    'bg-yellow-500': layer.status === 'extracting',
                    'bg-purple-500': layer.status === 'verifying',
                    'bg-green-500': layer.status === 'complete',
                    'bg-gray-400': layer.status === 'pending'
                  }"
                  class="h-2 rounded-full transition-all duration-300"
                  :style="{ width: layer.size > 0 ? Math.min((layer.progress / layer.size) * 100, 100) + '%' : (layer.status === 'complete' ? '100%' : '0%') }"
                ></div>
              </div>
            </div>
          </div>

          <!-- 镜像就绪状态 -->
          <div v-else-if="status.ready" class="text-center py-4">
            <svg class="w-8 h-8 text-green-500 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
            <p class="text-sm text-green-600 font-medium">{{ $t('servers.dockerImages.ready') }}</p>
          </div>

          <!-- 镜像未就绪状态 -->
          <div v-else class="text-center py-4">
            <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
            </svg>
            <p class="text-sm text-gray-600">{{ $t('servers.dockerImages.waitingDownload') }}</p>
          </div>
        </div>
      </div>
      
      <!-- 统计信息 -->
      <div v-if="imageStatus.total_images" class="mt-6 pt-4 border-t border-gray-200">
        <div class="flex justify-between items-center text-sm text-gray-600">
          <span>{{ $t('servers.dockerImages.totalImages') }}: {{ imageStatus.total_images }}</span>
          <span v-if="imageStatus.pulling_count > 0" class="text-yellow-600 font-medium">
            {{ $t('servers.dockerImages.downloadingCount') }}: {{ imageStatus.pulling_count }}
          </span>
        </div>
      </div>
    </template>

    <template #footer>
      <div class="flex justify-between items-center w-full">
        <div class="flex gap-2">
          <!-- 手动下载按钮 -->
          <UButton
            v-if="!imageStatus.can_start_server && !imageStatus.any_pulling"
            color="green"
            icon="i-heroicons-arrow-down-tray"
            :label="$t('servers.dockerImages.manualDownload')"
            @click="$emit('manual-download')"
          />
          
          <!-- 检查更新按钮 -->
          <UButton
            v-if="imageStatus.can_start_server"
            color="purple"
            icon="i-heroicons-arrow-path"
            :label="$t('servers.dockerImages.checkUpdates')"
            @click="$emit('check-updates')"
          />
        </div>
        
        <UButton
          color="neutral"
          :label="$t('common.close')"
          @click="handleClose"
        />
      </div>
    </template>
  </UModal>
</template>

<script setup>
// 定义props
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  imageStatus: {
    type: Object,
    required: true
  }
})

// 定义emits
const emit = defineEmits(['close', 'refresh', 'manual-download', 'check-updates'])

// 国际化
const { t } = useI18n()

// 本地响应式状态
const isOpen = ref(false)

// 监听 show prop 的变化
watch(() => props.show, (newValue) => {
  isOpen.value = newValue
}, { immediate: true })

// 监听本地状态的变化，并通知父组件
watch(isOpen, (newValue) => {
  if (!newValue && props.show) {
    handleClose()
  }
})

// 处理关闭事件
const handleClose = () => {
  isOpen.value = false
  emit('close')
}

// 格式化字节大小
const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取镜像显示名称
const getImageDisplayName = (imageName) => {
  switch (imageName) {
    case 'tbro98/ase-server:latest': return t('servers.dockerImages.arkServer')
    case 'alpine:latest': return t('servers.dockerImages.alpineSystem')
    default: return imageName
  }
}

// 获取状态文本
const getStatusText = (status) => {
  if (status.pulling) {
    return t('servers.dockerImages.downloading')
  } else if (status.ready) {
    return t('servers.dockerImages.ready')
  } else {
    return t('servers.dockerImages.notReady')
  }
}
</script> 