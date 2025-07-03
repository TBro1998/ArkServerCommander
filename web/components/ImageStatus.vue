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
        @click="$emit('refresh')"
        class="text-blue-600 hover:text-blue-800 p-1 transition-colors"
        title="刷新镜像状态"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
      </button>

    </div>
    
    <!-- 总体进度条 -->
    <div v-if="imageStatus.overall_progress !== undefined" class="w-full">
      <div class="flex justify-between text-xs text-gray-600 mb-1">
        <span>总体进度</span>
        <span>{{ imageStatus.overall_progress }}%</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
        <div 
          class="bg-gradient-to-r from-blue-500 to-blue-600 h-2 rounded-full transition-all duration-500 ease-out"
          :style="{ width: imageStatus.overall_progress + '%' }"
        ></div>
      </div>
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
            {{ status.percent }}%
          </span>
        </div>
        <div v-if="status.pulling || !status.ready" class="text-xs text-gray-500 mt-1 break-words">
          {{ status.progress }}
          <span v-if="status.percent !== undefined" class="text-blue-600">({{ status.percent }}%)</span>
        </div>
        <!-- 镜像大小信息 -->
        <div v-if="status.size" class="text-xs text-gray-600 mt-1">
          大小: {{ formatBytes(status.size) }}
        </div>
        <!-- 单个镜像进度条 -->
        <div v-if="status.pulling" class="mt-1">
          <div class="w-full bg-gray-100 rounded-full h-1">
            <div 
              class="bg-yellow-500 h-1 rounded-full transition-all duration-300"
              :style="{ width: status.percent + '%' }"
            ></div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 统计信息 -->
    <div v-if="imageStatus.total_images" class="text-xs text-gray-500 border-t pt-2">
      <div class="flex justify-between">
        <span>镜像总数: {{ imageStatus.total_images }}</span>
        <span v-if="imageStatus.pulling_count > 0">下载中: {{ imageStatus.pulling_count }}</span>
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
const emit = defineEmits(['refresh'])

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
    case 'tbro98/ase-server:latest': return 'ARK服务器'
    case 'alpine:latest': return 'Alpine系统'
    default: return imageName
  }
}
</script> 