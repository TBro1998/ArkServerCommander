<template>
  <div
    v-if="show"
    class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    @click="$emit('close')"
  >
    <div class="relative top-10 mx-auto p-5 border w-full max-w-7xl shadow-lg rounded-md bg-white" @click.stop>
      <div class="mb-4">
        <h3 class="text-lg font-bold text-gray-900">
          编辑服务器配置文件 - {{ server?.name }}
        </h3>
      </div>

      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">加载配置文件中...</p>
      </div>

      <div v-else class="space-y-6">
        <!-- 配置文件选项卡 -->
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-8">
            <button
              @click="activeTab = 'game_user_settings'"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm',
                activeTab === 'game_user_settings'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              GameUserSettings.ini
            </button>
            <button
              @click="activeTab = 'game_ini'"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm',
                activeTab === 'game_ini'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              Game.ini
            </button>
          </nav>
        </div>

        <!-- 配置文件编辑区域 -->
        <div class="space-y-4">
          <!-- GameUserSettings.ini 编辑器 -->
          <GameUserSettingsEditor
            v-if="activeTab === 'game_user_settings'"
            v-model="localConfigData.game_user_settings"
            :server="server"
          />

          <!-- Game.ini 编辑器 -->
          <GameIniEditor
            v-if="activeTab === 'game_ini'"
            v-model="localConfigData.game_ini"
            :server="server"
          />
        </div>

        <!-- 配置文件路径信息 -->
        <div v-if="configPaths" class="bg-gray-50 p-4 rounded-lg">
          <h4 class="text-sm font-medium text-gray-700 mb-2">配置文件路径信息</h4>
          <div class="space-y-1 text-xs text-gray-600">
            <div v-if="activeTab === 'game_user_settings'">
              <span class="font-medium">GameUserSettings.ini:</span>
              <span class="font-mono ml-2">{{ configPaths.game_user_settings_path }}</span>
            </div>
            <div v-if="activeTab === 'game_ini'">
              <span class="font-medium">Game.ini:</span>
              <span class="font-mono ml-2">{{ configPaths.game_ini_path }}</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-3 pt-4">
          <button
            @click="saveConfig"
            :disabled="saving"
            class="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors disabled:opacity-50"
          >
            {{ saving ? '保存中...' : '保存配置' }}
          </button>
          <button
            @click="$emit('close')"
            class="flex-1 bg-gray-300 hover:bg-gray-400 text-gray-700 py-2 px-4 rounded-lg transition-colors"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import GameUserSettingsEditor from './GameUserSettingsEditor.vue'
import GameIniEditor from './GameIniEditor.vue'

// Props
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  server: {
    type: Object,
    default: null
  },
  configData: {
    type: Object,
    default: () => ({
      game_user_settings: '',
      game_ini: ''
    })
  },
  configPaths: {
    type: Object,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  },
  saving: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['close', 'save'])

// 本地状态
const activeTab = ref('game_user_settings')
const localConfigData = ref({
  game_user_settings: '',
  game_ini: ''
})

// 监听 configData 变化，同步到本地状态
watch(() => props.configData, (newData) => {
  localConfigData.value = { ...newData }
}, { deep: true, immediate: true })

// 监听 show 变化，重置选项卡
watch(() => props.show, (newShow) => {
  if (newShow) {
    activeTab.value = 'game_user_settings'
  }
})

// 保存配置
const saveConfig = () => {
  emit('save', localConfigData.value)
}
</script> 