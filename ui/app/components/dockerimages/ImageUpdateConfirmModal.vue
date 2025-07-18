<template>
  <UModal v-model:open="isOpen" :title="$t('servers.dockerImages.updateConfirm')">
    <template #body>
      <div class="space-y-4">
        <!-- 镜像信息 -->
        <div class="bg-blue-50 p-4 rounded-lg border border-blue-200">
          <h3 class="text-sm font-semibold text-blue-900 mb-2">
            {{ $t('servers.dockerImages.imageInfo') }}
          </h3>
          <p class="text-sm text-blue-800">
            <span class="font-medium">{{ $t('servers.dockerImages.imageName') }}:</span>
            {{ imageName }}
          </p>
        </div>

        <!-- 影响的服务器列表 -->
        <div v-if="affectedServers.length > 0" class="bg-yellow-50 p-4 rounded-lg border border-yellow-200">
          <h3 class="text-sm font-semibold text-yellow-900 mb-2">
            {{ $t('servers.dockerImages.affectedServers') }} ({{ affectedServers.length }})
          </h3>
          <div class="space-y-2 max-h-32 overflow-y-auto">
            <div 
              v-for="server in affectedServers" 
              :key="server.id"
              class="flex items-center justify-between bg-white p-2 rounded border"
            >
              <div>
                <span class="text-sm font-medium text-gray-900">{{ server.identifier }}</span>
                <span class="text-xs text-gray-500 ml-2">({{ server.session_name }})</span>
              </div>
              <span 
                :class="{
                  'text-green-600 bg-green-100': server.status === 'running',
                  'text-gray-600 bg-gray-100': server.status === 'stopped',
                  'text-yellow-600 bg-yellow-100': server.status === 'starting' || server.status === 'stopping'
                }"
                class="text-xs px-2 py-1 rounded-full font-medium"
              >
                {{ $t(`servers.${server.status}`) }}
              </span>
            </div>
          </div>
        </div>

        <!-- 更新风险提示 -->
        <div class="bg-red-50 p-4 rounded-lg border border-red-200">
          <h3 class="text-sm font-semibold text-red-900 mb-2">
            {{ $t('servers.dockerImages.updateWarning') }}
          </h3>
          <ul class="text-sm text-red-800 space-y-1">
            <li>• {{ $t('servers.dockerImages.warningDownloadTime') }}</li>
            <li>• {{ $t('servers.dockerImages.warningContainerRecreate') }}</li>
            <li>• {{ $t('servers.dockerImages.warningDataSafety') }}</li>
          </ul>
        </div>

        <!-- 操作选项 -->
        <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
          <h3 class="text-sm font-semibold text-gray-900 mb-3">
            {{ $t('servers.dockerImages.updateOptions') }}
          </h3>
          <div class="space-y-3">
            <label class="flex items-start space-x-3">
              <input 
                v-model="updateOptions.updateImageOnly" 
                type="radio" 
                :value="true"
                name="updateOption"
                class="mt-1"
              >
              <div>
                <div class="text-sm font-medium text-gray-900">
                  {{ $t('servers.dockerImages.updateImageOnly') }}
                </div>
                <div class="text-xs text-gray-600">
                  {{ $t('servers.dockerImages.updateImageOnlyDesc') }}
                </div>
              </div>
            </label>
            <label class="flex items-start space-x-3">
              <input 
                v-model="updateOptions.updateImageOnly" 
                type="radio" 
                :value="false"
                name="updateOption"
                class="mt-1"
              >
              <div>
                <div class="text-sm font-medium text-gray-900">
                  {{ $t('servers.dockerImages.updateAndRecreate') }}
                </div>
                <div class="text-xs text-gray-600">
                  {{ $t('servers.dockerImages.updateAndRecreateDesc') }}
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>
    </template>

    <template #footer>
      <div class="flex justify-between items-center w-full">
        <UButton 
          color="gray" 
          variant="ghost"
          :label="$t('common.cancel')"
          @click="handleCancel"
        />
        
        <div class="flex gap-2">
          <UButton 
            color="blue" 
            :loading="isUpdating"
            :label="$t('servers.dockerImages.confirmUpdate')"
            @click="handleConfirm"
          />
        </div>
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
  imageName: {
    type: String,
    required: true
  },
  affectedServers: {
    type: Array,
    default: () => []
  }
})

// 定义emits
const emit = defineEmits(['close', 'confirm'])

// 国际化
const { t } = useI18n()

// 本地响应式状态
const isOpen = ref(false)
const isUpdating = ref(false)
const updateOptions = ref({
  updateImageOnly: true
})

// 监听 show prop 的变化
watch(() => props.show, (newValue) => {
  isOpen.value = newValue
  if (newValue) {
    // 重置状态
    isUpdating.value = false
    updateOptions.value.updateImageOnly = true
  }
}, { immediate: true })

// 监听本地状态的变化，并通知父组件
watch(isOpen, (newValue) => {
  if (!newValue && props.show) {
    handleCancel()
  }
})

// 处理取消事件
const handleCancel = () => {
  isOpen.value = false
  emit('close')
}

// 处理确认事件
const handleConfirm = () => {
  isUpdating.value = true
  emit('confirm', {
    imageName: props.imageName,
    updateImageOnly: updateOptions.value.updateImageOnly,
    affectedServers: props.affectedServers
  })
}

// 暴露方法给父组件
defineExpose({
  setUpdating: (value) => {
    isUpdating.value = value
  },
  close: () => {
    isOpen.value = false
  }
})
</script>