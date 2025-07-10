<template>
  <UCard class="h-full hover:shadow-lg transition-all duration-300 group">
    <!-- 卡片头部 -->
    <template #header>
      <div class="flex justify-between items-start">
        <div class="min-w-0 flex-1">
          <h3 class="text-lg font-semibold text-gray-900 truncate">{{ server.identifier }}</h3>
        </div>
        <div class="flex gap-1 ml-2">
          <!-- 启动/停止按钮 -->
          <UButton
            v-if="server.status === 'stopped'"
            @click="$emit('start', server)"
            :disabled="!canStartServer"
            color="green"
            variant="ghost"
            size="xs"
            :title="canStartServer ? t('servers.card.startServer') : t('servers.card.cannotStartImageNotReady')"
          >
            <UIcon name="i-lucide-play" class="w-4 h-4" />
          </UButton>
          
          <UButton
            v-else-if="server.status === 'running'"
            @click="$emit('stop', server)"
            color="red"
            variant="ghost"
            size="xs"
            :title="t('servers.card.stopServer')"
          >
            <UIcon name="i-lucide-square" class="w-4 h-4" />
          </UButton>
          
          <UButton
            v-else-if="server.status === 'starting'"
            color="yellow"
            variant="ghost"
            size="xs"
            disabled
            :title="t('servers.card.starting')"
          >
            <UIcon name="i-lucide-loader-2" class="w-4 h-4 animate-spin" />
          </UButton>
          
          <UButton
            v-else-if="server.status === 'stopping'"
            color="yellow"
            variant="ghost"
            size="xs"
            disabled
            :title="t('servers.card.stopping')"
          >
            <UIcon name="i-lucide-loader-2" class="w-4 h-4 animate-spin" />
          </UButton>
          
          <UButton
            v-else
            color="gray"
            variant="ghost"
            size="xs"
            disabled
            :title="t('servers.card.unknownStatus')"
          >
            <UIcon name="i-lucide-help-circle" class="w-4 h-4" />
          </UButton>
          
          <!-- RCON信息按钮 -->
          <UPopover>
            <UButton
              color="green"
              variant="ghost"
              size="xs"
              :title="t('servers.card.rconInfo')"
            >
              <UIcon name="i-lucide-info" class="w-4 h-4" />
            </UButton>
            
            <template #content="{ close }">
              <div class="p-4 space-y-4">
                <div class="flex items-center space-x-2 mb-3">
                  <UIcon name="i-lucide-info" class="w-4 h-4 text-green-600" />
                  <h4 class="text-sm font-semibold text-gray-700">{{ t('servers.card.rconConnectionInfo') }}</h4>
                </div>
                
                <div class="space-y-3">
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('servers.card.serverIdentifier') }}</label>
                    <div class="flex gap-2">
                      <span class="text-xs font-mono bg-gray-100 px-2 py-1 rounded flex-1">{{ server.identifier }}</span>
                      <UButton
                        @click="copyToClipboard(server.identifier)"
                        color="blue"
                        variant="ghost"
                        size="xs"
                      >
                        {{ t('servers.card.copy') }}
                      </UButton>
                    </div>
                  </div>
                  
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('servers.card.rconPort') }}</label>
                    <div class="flex gap-2">
                      <span class="text-xs font-mono bg-gray-100 px-2 py-1 rounded flex-1">{{ server.rcon_port }}</span>
                      <UButton
                        @click="copyToClipboard(server.rcon_port)"
                        color="blue"
                        variant="ghost"
                        size="xs"
                      >
                        {{ t('servers.card.copy') }}
                      </UButton>
                    </div>
                  </div>
                  
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">{{ t('servers.card.adminPassword') }}</label>
                    <div class="flex gap-2">
                      <span class="text-xs font-mono bg-gray-100 px-2 py-1 rounded flex-1">
                        {{ showPassword ? server.admin_password : '***' }}
                      </span>
                      <UButton
                        @click="togglePassword"
                        color="gray"
                        variant="ghost"
                        size="xs"
                      >
                        <UIcon :name="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'" class="w-3 h-3" />
                      </UButton>
                      <UButton
                        @click="copyToClipboard(server.admin_password)"
                        color="blue"
                        variant="ghost"
                        size="xs"
                      >
                        {{ t('servers.card.copy') }}
                      </UButton>
                    </div>
                  </div>
                </div>
                
                <div class="flex justify-end pt-2 border-t border-gray-100">
                  <UButton
                    @click="close"
                    color="gray"
                    variant="ghost"
                    size="xs"
                  >
                    {{ t('servers.card.close') }}
                  </UButton>
                </div>
              </div>
            </template>
          </UPopover>
          
          <!-- 编辑按钮 -->
          <UButton
            @click="$emit('edit', server)"
            color="blue"
            variant="ghost"
            size="xs"
            :title="t('servers.card.editServer')"
          >
            <UIcon name="i-lucide-edit" class="w-4 h-4" />
          </UButton>
          
          <!-- 删除按钮 -->
          <UPopover>
            <UButton
              color="red"
              variant="ghost"
              size="xs"
              :title="t('servers.card.deleteServer')"
            >
              <UIcon name="i-lucide-trash-2" class="w-4 h-4" />
            </UButton>
            
            <template #content="{ close }">
              <div class="p-4 space-y-4">
                <div class="flex items-center space-x-2 mb-3">
                  <UIcon name="i-lucide-alert-triangle" class="w-4 h-4 text-red-600" />
                  <h4 class="text-sm font-semibold text-gray-700">{{ t('servers.card.confirmDelete') }}</h4>
                </div>
                
                <p class="text-xs text-gray-600">
                  {{ t('servers.card.confirmDeleteMessage', { identifier: server.identifier }) }}
                </p>
                
                <div class="flex gap-2 pt-2 border-t border-gray-100">
                  <UButton
                    @click="$emit('delete', server); close()"
                    color="red"
                    variant="solid"
                    size="xs"
                  >
                    {{ t('servers.card.deleteServer') }}
                  </UButton>
                  <UButton
                    @click="close"
                    color="gray"
                    variant="ghost"
                    size="xs"
                  >
                    {{ t('common.cancel') }}
                  </UButton>
                </div>
              </div>
            </template>
          </UPopover>
        </div>
      </div>
    </template>

    <!-- 卡片内容 -->
    <div class="space-y-6">
      <!-- 状态徽章 -->
      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-600">{{ t('servers.card.status') }}</span>
        <UBadge
          :color="getStatusColor(server.status)"
          :variant="getStatusVariant(server.status)"
          size="sm"
        >
          {{ getStatusText(server.status) }}
        </UBadge>
      </div>

      <!-- 基本信息 -->
      <div class="space-y-3">
        <div v-if="server.session_name" class="flex justify-between items-center">
          <span class="text-sm text-gray-600">{{ t('servers.card.serverName') }}</span>
          <span class="text-sm font-medium text-gray-900">{{ server.session_name }}</span>
        </div>
        <div v-if="server.cluster_id" class="flex justify-between items-center">
          <span class="text-sm text-gray-600">{{ t('servers.card.clusterId') }}</span>
          <span class="text-sm font-medium text-gray-900">{{ server.cluster_id }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">{{ t('servers.card.map') }}</span>
          <span class="text-sm font-medium text-gray-900">{{ server.map || 'TheIsland' }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">{{ t('servers.card.maxPlayers') }}</span>
          <span class="text-sm font-medium text-gray-900">{{ server.max_players || 70 }}</span>
        </div>
      </div>

      <!-- 端口配置 -->
      <div class="border-t border-gray-100 pt-4">
        <div class="flex items-center space-x-2 mb-3">
          <UIcon name="i-lucide-wifi" class="w-4 h-4 text-blue-600" />
          <h4 class="text-sm font-semibold text-gray-700">{{ t('servers.card.portConfig') }}</h4>
        </div>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between items-center">
            <span class="text-gray-600">{{ t('servers.card.gamePort') }}</span>
            <span class="font-mono text-blue-600 text-xs">{{ server.port }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600">{{ t('servers.card.queryPort') }}</span>
            <span class="font-mono text-blue-600 text-xs">{{ server.query_port }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600">{{ t('servers.card.rconPortLabel') }}</span>
            <span class="font-mono text-blue-600 text-xs">{{ server.rcon_port }}</span>
          </div>
        </div>
      </div>

      <!-- 认证信息 -->
      <div class="border-t border-gray-100 pt-4">
        <div class="flex items-center space-x-2 mb-3">
          <UIcon name="i-lucide-lock" class="w-4 h-4 text-green-600" />
          <h4 class="text-sm font-semibold text-gray-700">{{ t('servers.card.authInfo') }}</h4>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">{{ t('servers.card.adminPassword') }}</span>
          <div class="flex items-center gap-2">
            <span class="font-mono text-gray-800 text-xs">
              {{ showPassword ? server.admin_password : '***' }}
            </span>
            <UButton
              @click="togglePassword"
              color="gray"
              variant="ghost"
              size="xs"
              :title="showPassword ? t('servers.card.hidePassword') : t('servers.card.showPassword')"
            >
              <UIcon :name="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'" class="w-3 h-3" />
            </UButton>
          </div>
        </div>
      </div>

      <!-- 时间信息 -->
      <div class="border-t border-gray-100 pt-4">
        <div class="flex items-center space-x-2 mb-3">
          <UIcon name="i-lucide-clock" class="w-4 h-4 text-purple-600" />
          <h4 class="text-sm font-semibold text-gray-700">{{ t('servers.card.timeInfo') }}</h4>
        </div>
        <div class="space-y-2 text-xs">
          <div class="flex justify-between items-center">
            <span class="text-gray-600">{{ t('servers.card.createdAt') }}</span>
            <span class="text-gray-800">{{ formatDate(server.created_at) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600">{{ t('servers.card.updatedAt') }}</span>
            <span class="text-gray-800">{{ formatDate(server.updated_at) }}</span>
          </div>
        </div>
      </div>

      <!-- 服务器ID -->
      <div class="border-t border-gray-100 pt-4">
        <div class="flex justify-between items-center text-xs text-gray-500">
          <span>{{ t('servers.card.serverId') }}</span>
          <span class="font-mono">#{{ server.id }}</span>
        </div>
      </div>
    </div>
  </UCard>
</template>

<script setup>
// 导入国际化
const { t } = useI18n()

// Props
const props = defineProps({
  server: {
    type: Object,
    required: true
  },
  canStartServer: {
    type: Boolean,
    default: true
  }
})

// Emits
const emit = defineEmits(['start', 'stop', 'edit', 'delete', 'rcon'])

// 响应式数据
const showPassword = ref(false)

// 方法
const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const getStatusText = (status) => {
  switch (status) {
    case 'running': return t('servers.running')
    case 'stopped': return t('servers.stopped')
    case 'starting': return t('servers.starting')
    case 'stopping': return t('servers.stopping')
    default: return t('servers.unknown')
  }
}

const getStatusColor = (status) => {
  switch (status) {
    case 'running': return 'green'
    case 'stopped': return 'red'
    case 'starting':
    case 'stopping': return 'yellow'
    default: return 'gray'
  }
}

const getStatusVariant = (status) => {
  switch (status) {
    case 'running': return 'solid'
    case 'stopped': return 'solid'
    case 'starting':
    case 'stopping': return 'soft'
    default: return 'soft'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return t('common.loading')
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

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text.toString())
    // 使用翻译的提示信息
    alert(t('servers.copyToClipboard'))
  } catch (error) {
    console.error(t('servers.copyFailed'), error)
    alert(t('servers.copyFailed'))
  }
}
</script> 