<template>
  <div class="space-y-4">
    <!-- 编辑模式切换 -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div class="flex flex-wrap gap-2 sm:gap-4">
        <button
          @click="switchToVisualMode"
          :class="[
            'px-3 py-2 rounded-md text-sm font-medium transition-colors',
            editMode === 'visual'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          ]"
        >
          <i class="fas fa-sliders-h mr-2"></i>可视化编辑
        </button>
        <button
          @click="switchToTextMode"
          :class="[
            'px-3 py-2 rounded-md text-sm font-medium transition-colors',
            editMode === 'text'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          ]"
        >
          <i class="fas fa-code mr-2"></i>文本编辑
        </button>
      </div>
      <div class="flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
        <div class="text-sm text-gray-500">
          GameUserSettings.ini - {{ editMode === 'visual' ? '可视化编辑' : '文本编辑' }}
        </div>
        <div v-if="isSyncing" class="flex items-center text-sm text-blue-600">
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
          正在同步...
        </div>
        <div v-else class="flex items-center text-sm text-green-600">
          <i class="fas fa-check-circle mr-1"></i>
          已同步
        </div>
      </div>
    </div>

    <!-- 可视化编辑模式 -->
    <div v-if="editMode === 'visual'" class="space-y-6">
      <div 
        v-for="(section, sectionKey) in gameUserSettingsParams" 
        :key="sectionKey"
        class="bg-gray-50 rounded-lg p-4"
      >
        <h4 class="text-base sm:text-lg font-semibold text-gray-800 mb-4 border-b border-gray-200 pb-2">
          {{ section.title }}
        </h4>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <div 
            v-for="(param, paramKey) in section.params" 
            :key="paramKey"
            class="space-y-2 relative"
          >
            <label class="block text-sm font-medium text-gray-700">
              {{ param.label }}
              <span 
                class="relative inline-block ml-1"
                @mouseenter="showTooltip = paramKey"
                @mouseleave="showTooltip = null"
              >
                <i class="fas fa-info-circle text-gray-400 cursor-help" :title="param.description"></i>
                <div 
                  v-if="showTooltip === paramKey"
                  class="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                  style="bottom: 100%; left: 50%; transform: translateX(-50%); margin-bottom: 8px;"
                >
                  {{ param.description }}
                  <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
                </div>
              </span>
            </label>
            
            <!-- Boolean 类型 -->
            <div v-if="param.type === 'boolean'">
              <ToggleSwitch
                :id="paramKey"
                v-model="visualConfig[paramKey]"
                :label="visualConfig[paramKey] ? '启用' : '禁用'"
              />
            </div>
            
            <!-- Number 类型 -->
            <input
              v-else-if="param.type === 'number'"
              v-model.number="visualConfig[paramKey]"
              type="number"
              :min="param.min"
              :max="param.max"
              :step="param.step || 1"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            
            <!-- Text/Password 类型 -->
            <div v-else-if="param.type === 'password'" class="relative">
              <input
                v-model="visualConfig[paramKey]"
                :type="showPasswords[paramKey] ? 'text' : 'password'"
                class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                @click="togglePasswordVisibility(paramKey)"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-700"
                :title="showPasswords[paramKey] ? '隐藏密码' : '显示密码'"
              >
                <svg v-if="showPasswords[paramKey]" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
              </button>
            </div>
            
            <!-- 其他文本类型 -->
            <input
              v-else
              v-model="visualConfig[paramKey]"
              :type="param.type"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 文本编辑模式 -->
    <div v-else>
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-2 gap-2">
        <label class="block text-sm font-medium text-gray-700">
          GameUserSettings.ini 内容
        </label>
        <div class="flex gap-2">
          <button
            @click="resetToDefault"
            class="text-sm text-blue-600 hover:text-blue-800"
          >
            重置为默认
          </button>
          <button
            @click="formatConfig"
            class="text-sm text-green-600 hover:text-green-800"
          >
            格式化
          </button>
        </div>
      </div>
      <textarea
        v-model="textContent"
        class="w-full h-64 sm:h-96 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
        placeholder="输入 GameUserSettings.ini 配置内容..."
      ></textarea>
      <p class="text-xs text-gray-500 mt-1">
        此文件包含服务器的基本设置，如端口、密码、最大玩家数等
      </p>
    </div>

    <!-- 提示信息 -->
    <div v-if="editMode === 'visual'" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
      <div class="flex flex-col sm:flex-row items-start">
        <i class="fas fa-info-circle text-blue-500 mr-0 sm:mr-3 mb-2 sm:mb-0"></i>
        <div class="text-sm text-blue-700">
          <p class="font-medium">可视化编辑模式</p>
          <p>通过表单控件修改参数，鼠标悬停在参数名称旁的 <i class="fas fa-info-circle text-blue-400"></i> 图标可查看详细说明。修改会自动同步到配置文件，切换到文本模式可查看生成的配置内容。</p>
        </div>
      </div>
    </div>
    
    <!-- 基础参数说明 -->
    <div class="bg-green-50 border border-green-200 rounded-lg p-4">
      <div class="flex flex-col sm:flex-row items-start">
        <i class="fas fa-check-circle text-green-500 mr-0 sm:mr-3 mb-2 sm:mb-0"></i>
        <div class="text-sm text-green-700">
          <p class="font-medium">基础参数说明</p>
          <p>以下基础参数已从服务器基本设置中自动获取，无需在此处重复设置：</p>
          <ul class="list-disc list-inside mt-1 space-y-1">
            <li><strong>管理员密码 (ServerAdminPassword)</strong> - 从服务器基本设置的管理员密码获取</li>
            <li><strong>最大玩家数 (MaxPlayers)</strong> - 固定为70</li>
            <li><strong>RCON启用 (RCONEnabled)</strong> - 固定为True</li>
          </ul>
          <p class="mt-2">这些参数会在启动时自动添加到启动参数字符串中。</p>
        </div>
      </div>
    </div>

    <div v-if="editMode === 'text'" class="bg-green-50 border border-green-200 rounded-lg p-4">
      <div class="flex flex-col sm:flex-row items-start">
        <i class="fas fa-code text-green-500 mr-0 sm:mr-3 mb-2 sm:mb-0"></i>
        <div class="text-sm text-green-700">
          <p class="font-medium">文本编辑模式</p>
          <p>直接编辑 INI 配置文件内容。修改会自动解析并同步到可视化界面，切换到可视化模式可看到解析后的参数设置。</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { gameUserSettingsParams } from '../utils/gameUserSettingsParams.js'
import { extractConfigValues, formatConfigContent } from '../utils/configParser.js'
import ToggleSwitch from './ToggleSwitch.vue'

// Props
const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  server: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 响应式状态
const editMode = ref('visual')
const showTooltip = ref(null)
const isSyncing = ref(false)
const textContent = ref('')
const visualConfig = ref({})
const showPasswords = ref({})

// 防抖定时器
let visualSyncTimer = null
let textSyncTimer = null

// 初始化可视化配置默认值
const initializeVisualConfig = () => {
  try {
    if (!gameUserSettingsParams) {
      console.warn('GameUserSettings 配置参数尚未加载，跳过初始化')
      return
    }

    // 只在 visualConfig 为空时才设置默认值
    const isVisualConfigEmpty = Object.keys(visualConfig.value).length === 0

    Object.keys(gameUserSettingsParams).forEach(sectionKey => {
      const section = gameUserSettingsParams[sectionKey]
      if (section && section.params) {
        Object.keys(section.params).forEach(paramKey => {
          const param = section.params[paramKey]
          // 只在参数不存在时才设置默认值
          if (param && param.default !== undefined && (isVisualConfigEmpty || visualConfig.value[paramKey] === undefined)) {
            visualConfig.value[paramKey] = param.default
          }
        })
      }
    })
    
    console.log('初始化后的可视化配置:', visualConfig.value)
  } catch (error) {
    console.error('初始化 GameUserSettings 可视化配置失败:', error)
  }
}

// 解析文本配置到可视化配置
const parseTextToVisual = () => {
  try {
    if (!gameUserSettingsParams) {
      console.warn('GameUserSettings 配置参数尚未加载')
      return
    }

    if (textContent.value) {
      console.log('开始解析 GameUserSettings 配置:', textContent.value.substring(0, 200) + '...')
      const values = extractConfigValues(
        textContent.value, 
        gameUserSettingsParams
      )
      console.log('解析得到的配置值:', values)
      if (values) {
        // 合并而不是覆盖，保留未在配置文件中定义的默认值
        Object.assign(visualConfig.value, values)
        console.log('更新后的可视化配置:', visualConfig.value)
      }
    }
  } catch (error) {
    console.error('解析 GameUserSettings 配置文件失败:', error)
  }
}

// 将可视化配置同步到文本
const syncVisualToText = () => {
  try {
    if (!gameUserSettingsParams) {
      console.warn('GameUserSettings 配置参数尚未加载，无法同步')
      return
    }

    let content = '[ServerSettings]\n'
    Object.keys(gameUserSettingsParams).forEach(sectionKey => {
      const section = gameUserSettingsParams[sectionKey]
      if (section && section.params) {
        Object.keys(section.params).forEach(paramKey => {
          const value = visualConfig.value[paramKey]
          if (value !== undefined && value !== null && value !== '') {
            content += `${paramKey}=${value}\n`
          }
        })
      }
    })
    
    content += '\n[SessionSettings]\n'
    content += `SessionName=${visualConfig.value.SessionName || 'ARK Server'}\n`
    content += `Port=${visualConfig.value.Port || 7777}\n`
    content += `QueryPort=${visualConfig.value.QueryPort || 27015}\n`
    
    content += '\n[/Script/Engine.GameSession]\n'
    content += `MaxPlayers=${visualConfig.value.MaxPlayers || 70}\n`
    
    content += '\n[MessageOfTheDay]\n'
    content += 'Message=欢迎来到 ARK 服务器！\n'
    content += 'Duration=30\n'

    textContent.value = content
    emit('update:modelValue', content)
  } catch (error) {
    console.error('同步 GameUserSettings 可视化配置到文本失败:', error)
  }
}

// 监听 props.modelValue 变化
watch(() => props.modelValue, (newValue) => {
  textContent.value = newValue || ''
  // 只有在 gameUserSettingsParams 已加载时才解析
  if (gameUserSettingsParams) {
    parseTextToVisual()
  }
}, { immediate: true })

// 监听 gameUserSettingsParams 的变化（处理异步加载）
watch(() => gameUserSettingsParams, (newParams) => {
  if (newParams) {
    console.log('GameUserSettings 参数已加载，重新初始化')
    initializeVisualConfig()
    if (textContent.value) {
      parseTextToVisual()
    }
  }
}, { immediate: true })

// 监听可视化配置变化，自动同步到文本
watch(() => visualConfig.value, () => {
  if (editMode.value === 'visual') {
    isSyncing.value = true
    if (visualSyncTimer) {
      clearTimeout(visualSyncTimer)
    }
    visualSyncTimer = setTimeout(() => {
      syncVisualToText()
      isSyncing.value = false
    }, 500)
  }
}, { deep: true })

// 监听文本配置变化，自动同步到可视化
watch(() => textContent.value, () => {
  if (editMode.value === 'text') {
    isSyncing.value = true
    if (textSyncTimer) {
      clearTimeout(textSyncTimer)
    }
    textSyncTimer = setTimeout(() => {
      parseTextToVisual()
      emit('update:modelValue', textContent.value)
      isSyncing.value = false
    }, 800)
  }
})

// 切换到可视化模式
const switchToVisualMode = () => {
  if (editMode.value !== 'visual') {
    editMode.value = 'visual'
    parseTextToVisual()
  }
}

// 切换到文本模式
const switchToTextMode = () => {
  if (editMode.value !== 'text') {
    editMode.value = 'text'
    syncVisualToText()
  }
}

// 重置为默认配置
const resetToDefault = () => {
  if (!props.server) return
  
  const server = props.server
  
  textContent.value = `[ServerSettings]
SessionName=${server.identifier}
ServerPassword=
ServerAdminPassword=${server.admin_password}
Port=${server.port}
QueryPort=${server.query_port}
RCONEnabled=True
RCONPort=${server.rcon_port}
MaxPlayers=70

[SessionSettings]
SessionName=${server.identifier}

[MessageOfTheDay]
Message=欢迎来到 ${server.identifier} ARK 服务器！

[/Script/ShooterGame.ShooterGameMode]
bUseSingleplayerSettings=False
bDisableStructurePlacementCollision=False
bAllowFlyerCarryPvE=True
bDisableStructureDecayPvE=False

[RCONSettings]
RCONEnabled=True
RCONPort=${server.rcon_port}`

  emit('update:modelValue', textContent.value)
}

// 格式化配置文件内容
const formatConfig = () => {
  textContent.value = formatConfigContent(textContent.value)
  emit('update:modelValue', textContent.value)
}

// 切换密码显示状态
const togglePasswordVisibility = (paramKey) => {
  showPasswords.value[paramKey] = !showPasswords.value[paramKey]
}

// 组件挂载时初始化
onMounted(() => {
  console.log('GameUserSettingsEditor 组件已挂载')
  console.log('初始 textContent:', textContent.value.substring(0, 100) + '...')
  console.log('gameUserSettingsParams 是否已加载:', !!gameUserSettingsParams)
  
  // 如果参数已经加载，立即初始化
  if (gameUserSettingsParams) {
    initializeVisualConfig()
    if (textContent.value) {
      parseTextToVisual()
    }
  }
})

// 组件卸载时清理定时器
onUnmounted(() => {
  if (visualSyncTimer) {
    clearTimeout(visualSyncTimer)
  }
  if (textSyncTimer) {
    clearTimeout(textSyncTimer)
  }
})
</script>
 