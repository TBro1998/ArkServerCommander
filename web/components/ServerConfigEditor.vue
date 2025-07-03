<template>
  <div
    v-if="show"
    class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    @click="$emit('close')"
  >
    <div class="relative top-10 mx-auto p-5 border w-full max-w-7xl shadow-lg rounded-md bg-white" @click.stop>
      <div class="mb-4">
        <h3 class="text-lg font-bold text-gray-900">
          编辑服务器配置文件 - {{ server?.identifier }}
        </h3>
      </div>

      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">加载配置文件中...</p>
      </div>

      <div v-else class="space-y-6">
        <!-- 配置说明 -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="flex items-start">
            <svg class="w-5 h-5 text-blue-600 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div class="text-sm text-blue-800">
              <p class="font-medium mb-1">服务器启动配置</p>
              <p>这些参数控制ARK服务器的启动行为。修改后需要重启服务器才能生效。</p>
            </div>
          </div>
        </div>

        <!-- 基本配置 -->
        <div class="space-y-4">
          <h4 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">基本配置</h4>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">游戏端口</label>
              <input
                v-model.number="config.PORT"
                type="number"
                min="1"
                max="65535"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="7777"
              />
              <p class="text-xs text-gray-500 mt-1">ARK游戏客户端连接端口</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">查询端口</label>
              <input
                v-model.number="config.QUERYPORT"
                type="number"
                min="1"
                max="65535"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="27015"
              />
              <p class="text-xs text-gray-500 mt-1">服务器查询和Steam列表端口</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">地图</label>
              <select
                v-model="config.MAP"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="TheIsland">The Island</option>
                <option value="TheCenter">The Center</option>
                <option value="ScorchedEarth_P">Scorched Earth</option>
                <option value="Aberration_P">Aberration</option>
                <option value="Extinction">Extinction</option>
                <option value="Valguero_P">Valguero</option>
                <option value="Genesis">Genesis</option>
                <option value="CrystalIsles">Crystal Isles</option>
                <option value="Genesis2">Genesis 2</option>
                <option value="LostIsland">Lost Island</option>
                <option value="Fjordur">Fjordur</option>
              </select>
              <p class="text-xs text-gray-500 mt-1">服务器运行的地图</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">最大玩家数</label>
              <input
                v-model.number="config.MAX_PLAYERS"
                type="number"
                min="1"
                max="255"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="70"
              />
              <p class="text-xs text-gray-500 mt-1">服务器最大玩家数量</p>
            </div>
          </div>
        </div>

        <!-- 更新设置 -->
        <div class="space-y-4">
          <h4 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">更新设置</h4>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="flex items-center">
                <input
                  v-model="config.UPDATE_SERVER"
                  type="checkbox"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <span class="ml-2 text-sm font-medium text-gray-700">自动更新服务器</span>
              </label>
              <p class="text-xs text-gray-500 mt-1">启动时自动更新ARK服务器</p>
            </div>
            
            <div>
              <label class="flex items-center">
                <input
                  v-model="config.UPDATE_MODS"
                  type="checkbox"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <span class="ml-2 text-sm font-medium text-gray-700">自动更新模组</span>
              </label>
              <p class="text-xs text-gray-500 mt-1">启动时自动更新已安装的模组</p>
            </div>
          </div>
        </div>

        <!-- 模组设置 -->
        <div class="space-y-4">
          <h4 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">模组设置</h4>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">模组ID列表</label>
            <textarea
              v-model="config.MODIDS"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="895711211,669673294,1136125765,554678442,926028694,676288311,876038468,566885854"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">Steam Workshop模组ID，用逗号分隔</p>
          </div>
        </div>

        <!-- 服务器参数 -->
        <div class="space-y-4">
          <h4 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">服务器启动参数</h4>
          
          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-yellow-600 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
              </svg>
              <div class="text-sm text-yellow-800">
                <p class="font-medium mb-1">高级设置</p>
                <p>这些参数影响服务器的性能和安全性。请谨慎修改，错误的参数可能导致服务器无法启动。</p>
              </div>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">服务器启动参数</label>
            <textarea
              v-model="config.SERVER_ARGS"
              rows="6"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="-NoBattlEye -servergamelog -structurememopts -UseStructureStasisGrid -SecureSendArKPayload -UseItemDupeCheck -UseSecureSpawnRules -nosteamclient -game -server -log -MinimumTimeBetweenInventoryRetrieval=3600 -newsaveformat -usestore"
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">ARK服务器启动命令行参数</p>
          </div>

          <!-- 常用参数说明 -->
          <div class="bg-gray-50 rounded-lg p-4">
            <h5 class="text-sm font-medium text-gray-900 mb-2">常用参数说明：</h5>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs text-gray-600">
              <div>
                <p><strong>-NoBattlEye:</strong> 禁用BattlEye反作弊</p>
                <p><strong>-servergamelog:</strong> 启用服务器游戏日志</p>
                <p><strong>-structurememopts:</strong> 优化结构内存使用</p>
                <p><strong>-UseStructureStasisGrid:</strong> 使用结构静止网格</p>
                <p><strong>-SecureSendArKPayload:</strong> 安全发送ARK数据包</p>
              </div>
              <div>
                <p><strong>-UseItemDupeCheck:</strong> 启用物品复制检查</p>
                <p><strong>-UseSecureSpawnRules:</strong> 使用安全生成规则</p>
                <p><strong>-nosteamclient:</strong> 不使用Steam客户端</p>
                <p><strong>-newsaveformat:</strong> 使用新的保存格式</p>
                <p><strong>-usestore:</strong> 启用商店功能</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 认证设置 -->
        <div class="space-y-4">
          <h4 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">认证设置</h4>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">服务器密码</label>
              <input
                v-model="config.SERVER_PASSWORD"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="留空表示无密码"
              />
              <p class="text-xs text-gray-500 mt-1">玩家连接服务器需要的密码</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">管理员密码</label>
              <div class="relative">
                <input
                  v-model="config.ADMIN_PASSWORD"
                  :type="showAdminPassword ? 'text' : 'password'"
                  class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="管理员密码"
                />
                <button
                  type="button"
                  @click="showAdminPassword = !showAdminPassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-700"
                  :title="showAdminPassword ? '隐藏密码' : '显示密码'"
                >
                  <svg v-if="showAdminPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                </button>
              </div>
              <p class="text-xs text-gray-500 mt-1">管理员权限密码（同时作为RCON密码）</p>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
          <button
            @click="resetToDefault"
            type="button"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            重置为默认值
          </button>
          <button
            @click="generateConfig"
            type="button"
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            生成配置文件
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

// 响应式数据
const showAdminPassword = ref(false)

// 配置对象
const config = ref({
  PORT: 7777,
  QUERYPORT: 27015,
  MAP: 'TheIsland',
  SERVER_PASSWORD: '',
  ADMIN_PASSWORD: 'Admin',
  MAX_PLAYERS: 70,
  UPDATE_SERVER: false,
  UPDATE_MODS: true,
  MODIDS: '895711211,669673294,1136125765,554678442,926028694,676288311,876038468,566885854',
  SERVER_ARGS: '-NoBattlEye -servergamelog -structurememopts -UseStructureStasisGrid -SecureSendArKPayload -UseItemDupeCheck -UseSecureSpawnRules -nosteamclient -game -server -log -MinimumTimeBetweenInventoryRetrieval=3600 -newsaveformat -usestore'
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

// 监听服务器数据变化
watch(() => props.server, (newServer) => {
  if (newServer) {
    // 使用服务器数据更新配置
    config.value.PORT = newServer.port || 7777
    config.value.QUERYPORT = newServer.query_port || 27015
    config.value.MAP = newServer.map || 'TheIsland'
    config.value.ADMIN_PASSWORD = newServer.admin_password || 'Admin'
  }
}, { immediate: true })

// 解析配置文件内容
const parseConfigContent = (content) => {
  if (!content) return
  
  const lines = content.split('\n')
  for (const line of lines) {
    const trimmedLine = line.trim()
    if (!trimmedLine || trimmedLine.startsWith('#')) continue
    
    const equalIndex = trimmedLine.indexOf('=')
    if (equalIndex === -1) continue
    
    const key = trimmedLine.substring(0, equalIndex).trim()
    const value = trimmedLine.substring(equalIndex + 1).trim()
    
    // 移除引号
    const cleanValue = value.replace(/^["']|["']$/g, '')
    
    switch (key) {
      case 'PORT':
        config.value.PORT = parseInt(cleanValue) || 7777
        break
      case 'QUERYPORT':
        config.value.QUERYPORT = parseInt(cleanValue) || 27015
        break
      case 'MAP':
        config.value.MAP = cleanValue || 'TheIsland'
        break
      case 'SERVER_PASSWORD':
        config.value.SERVER_PASSWORD = cleanValue
        break
      case 'ADMIN_PASSWORD':
        config.value.ADMIN_PASSWORD = cleanValue || 'Admin'
        break
      case 'MAX_PLAYERS':
        config.value.MAX_PLAYERS = parseInt(cleanValue) || 70
        break
      case 'UPDATE_SERVER':
        config.value.UPDATE_SERVER = cleanValue.toLowerCase() === 'true'
        break
      case 'UPDATE_MODS':
        config.value.UPDATE_MODS = cleanValue.toLowerCase() === 'true'
        break
      case 'MODIDS':
        config.value.MODIDS = cleanValue
        break
      case 'SERVER_ARGS':
        config.value.SERVER_ARGS = cleanValue
        break
    }
  }
}

// 生成配置文件内容
const generateConfigContent = () => {
  const lines = [
    `PORT=${config.value.PORT}`,
    `QUERYPORT=${config.value.QUERYPORT}`,
    `MAP=${config.value.MAP}`,
    config.value.SERVER_PASSWORD ? `SERVER_PASSWORD=${config.value.SERVER_PASSWORD}` : 'SERVER_PASSWORD=',
    `ADMIN_PASSWORD=${config.value.ADMIN_PASSWORD}`,
    `MAX_PLAYERS=${config.value.MAX_PLAYERS}`,
    `UPDATE_SERVER=${config.value.UPDATE_SERVER}`,
    `UPDATE_MODS=${config.value.UPDATE_MODS}`,
    `MODIDS="${config.value.MODIDS}"`,
    `SERVER_ARGS="${config.value.SERVER_ARGS}"`
  ]
  
  const content = lines.join('\n')
  emit('save', content)
}

// 重置为默认值
const resetToDefault = () => {
  config.value = {
    PORT: 7777,
    QUERYPORT: 27015,
    MAP: 'TheIsland',
    SERVER_PASSWORD: '',
    ADMIN_PASSWORD: 'Admin',
    MAX_PLAYERS: 70,
    UPDATE_SERVER: false,
    UPDATE_MODS: true,
    MODIDS: '895711211,669673294,1136125765,554678442,926028694,676288311,876038468,566885854',
    SERVER_ARGS: '-NoBattlEye -servergamelog -structurememopts -UseStructureStasisGrid -SecureSendArKPayload -UseItemDupeCheck -UseSecureSpawnRules -nosteamclient -game -server -log -MinimumTimeBetweenInventoryRetrieval=3600 -newsaveformat -usestore'
  }
}

// 生成配置文件
const generateConfig = () => {
  generateConfigContent()
}
</script> 