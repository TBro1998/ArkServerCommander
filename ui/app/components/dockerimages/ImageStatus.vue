<template>
  <div class="flex flex-col gap-2 text-sm w-full sm:w-80 bg-gray-50 p-3 rounded-lg border">
    <!-- 总体状态 -->
    <div class="flex items-center gap-2">
      <div v-if="imageStatus.any_pulling" class="flex items-center gap-1 text-yellow-600">
        <svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        {{ imageStatus.overall_status }}
      </div>
      <div v-else-if="!imageStatus.can_start_server" class="flex items-center gap-1 text-red-600">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
        </svg>
        {{ $t('servers.dockerImages.imageMissingManualDownload') }}
      </div>
      <div v-else class="flex items-center gap-1 text-green-600">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
        </svg>
        {{ $t('servers.dockerImages.imageReady') }}
      </div>
      <UButton
        @click="$emit('refresh')"
        color="blue"
        variant="ghost"
        size="xs"
        icon="i-lucide-refresh-cw"
        :title="$t('servers.dockerImages.refreshStatus')"
      />
      
      <!-- 手动下载按钮 -->
      <UButton
        v-if="!imageStatus.can_start_server && !imageStatus.any_pulling"
        @click="$emit('manual-download')"
        color="green"
        variant="ghost"
        size="xs"
        icon="i-lucide-download"
        :title="$t('servers.dockerImages.manualDownload')"
      />
      
      <!-- 检查更新按钮 -->
      <UButton
        v-if="imageStatus.can_start_server"
        @click="$emit('check-updates')"
        color="purple"
        variant="ghost"
        size="xs"
        icon="i-lucide-refresh-ccw"
        :title="$t('servers.dockerImages.checkUpdates')"
      />

    </div>
    

    
    <!-- 详细镜像状态 -->
    <div v-if="imageStatus.images" class="space-y-2">
      <div 
        v-for="(status, imageName) in imageStatus.images" 
        :key="imageName"
        class="text-xs bg-white p-2 rounded border-l-4"
        :class="{
          'border-l-green-500': status.ready,
          'border-l-yellow-500': status.pulling,
          'border-l-red-500': !status.ready && !status.pulling
        }"
      >
        <div class="flex justify-between items-center">
          <span class="text-gray-700 font-medium truncate max-w-32">{{ getImageDisplayName(imageName) }}</span>
          <span 
            :class="{
              'text-green-600': status.ready,
              'text-yellow-600': status.pulling,
              'text-red-600': !status.ready && !status.pulling
            }"
            class="ml-2 font-semibold"
          >
            {{ getStatusText(status) }}
          </span>
        </div>

        <!-- 层级信息 -->
        <div v-if="status.pulling && status.layers" class="text-xs text-gray-600 mt-2 space-y-1">
          <div class="font-medium text-gray-700">{{ $t('servers.dockerImages.layerProgress') }}:</div>
          <div 
            v-for="(layer, layerId) in status.layers" 
            :key="layerId"
            class="bg-gray-50 p-2 rounded border"
          >
            <div class="flex justify-between items-center mb-1">
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
            
            <div v-if="layer.size > 0" class="text-xs text-gray-500 mb-1">
              {{ formatBytes(layer.progress) }} / {{ formatBytes(layer.size) }}
            </div>
            <div v-else class="text-xs text-gray-500 mb-1">
              {{ formatBytes(layer.progress) }} / {{ $t('servers.dockerImages.unknownSize') }}
            </div>
            
            <!-- 层级进度条 -->
            <div class="w-full bg-gray-200 rounded-full h-1.5 overflow-hidden">
              <div 
                :class="{
                  'bg-blue-500': layer.status === 'downloading',
                  'bg-yellow-500': layer.status === 'extracting',
                  'bg-purple-500': layer.status === 'verifying',
                  'bg-green-500': layer.status === 'complete',
                  'bg-gray-400': layer.status === 'pending'
                }"
                class="h-1.5 rounded-full transition-all duration-300"
                :style="{ width: layer.size > 0 ? Math.min((layer.progress / layer.size) * 100, 100) + '%' : (layer.status === 'complete' ? '100%' : '0%') }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 统计信息 -->
    <div v-if="imageStatus.total_images" class="text-xs text-gray-500 border-t pt-2">
      <div class="flex justify-between">
        <span>{{ $t('servers.dockerImages.totalImages') }}: {{ imageStatus.total_images }}</span>
        <span v-if="imageStatus.pulling_count > 0">{{ $t('servers.dockerImages.downloadingCount') }}: {{ imageStatus.pulling_count }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
// 定义props
const props = defineProps({
  imageStatus: {
    type: Object,
    required: true
  }
})

// 定义emits
const emit = defineEmits(['refresh', 'manual-download', 'check-updates'])

// 国际化
const { t } = useI18n()

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