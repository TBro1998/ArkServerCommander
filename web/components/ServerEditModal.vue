<template>
  <div
    v-if="show"
    class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    @click="$emit('close')"
  >
    <div class="relative top-2 sm:top-10 mx-2 sm:mx-4 lg:mx-8 xl:mx-12 p-4 sm:p-5 border w-auto shadow-lg rounded-md bg-white mb-4" @click.stop>
      <div class="mb-4">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-bold text-gray-900">
            {{ mode === 'create' ? '新增服务器' : '编辑服务器' }} 
            {{ server && mode === 'edit' ? `- ${server.identifier}` : '' }}
          </h3>
          
          <!-- 操作按钮 -->
          <div v-if="!loading" class="flex gap-2">
            <button
              @click="saveServer"
              :disabled="saving"
              class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors disabled:opacity-50"
            >
              {{ saving ? '保存中...' : (mode === 'create' ? '创建服务器' : '保存更改') }}
            </button>
            <button
              @click="$emit('close')"
              class="bg-gray-300 hover:bg-gray-400 text-gray-700 py-2 px-4 rounded-lg transition-colors"
            >
              取消
            </button>
          </div>
        </div>
      </div>

      <div v-if="loading" class="text-center py-8">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">{{ mode === 'create' ? '准备中...' : '加载服务器信息中...' }}</p>
      </div>

      <div v-else class="space-y-6">
        <!-- 选项卡导航 -->
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-4 sm:space-x-8 overflow-x-auto">
            <button
              @click="activeTab = 'basic'"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap',
                activeTab === 'basic'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              基本参数
            </button>
            <button
              @click="activeTab = 'game_user_settings'"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap',
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
                'py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap',
                activeTab === 'game_ini'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              Game.ini
            </button>
          </nav>
        </div>

        <!-- 内容区域 -->
        <div class="space-y-6">
          <!-- 基本参数选项卡 -->
          <div v-if="activeTab === 'basic'" class="space-y-4">
            <form @submit.prevent="saveServer" class="space-y-4">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">服务器标识 *</label>
                  <input
                    v-model="formData.identifier"
                    type="text"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="输入服务器标识"
                  />
                </div>
                <div></div>
              </div>

              <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">游戏端口 *</label>
                  <input
                    v-model.number="formData.port"
                    type="number"
                    required
                    min="1"
                    max="65535"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="7777"
                  />
                </div>
                <div></div>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">查询端口 *</label>
                  <input
                    v-model.number="formData.query_port"
                    type="number"
                    required
                    min="1"
                    max="65535"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="27015"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">RCON端口 *</label>
                  <input
                    v-model.number="formData.rcon_port"
                    type="number"
                    required
                    min="1"
                    max="65535"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="32330"
                  />
                </div>
              </div>

              <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">地图</label>
                  <select
                    v-model="formData.map"
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
                </div>
                <div></div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">管理员密码 *</label>
                <div class="relative">
                  <input
                    v-model="formData.admin_password"
                    :type="showFormPassword ? 'text' : 'password'"
                    required
                    class="w-full px-3 py-2 pr-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="输入管理员密码（同时作为RCON密码）"
                  />
                  <button
                    type="button"
                    @click="toggleFormPassword"
                    class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 hover:text-gray-700"
                    :title="showFormPassword ? '隐藏密码' : '显示密码'"
                  >
                    <svg v-if="showFormPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                    </svg>
                  </button>
                </div>
              </div>
            </form>
          </div>

          <!-- GameUserSettings.ini 选项卡 -->
          <div v-if="activeTab === 'game_user_settings'" class="space-y-4">
            <GameUserSettingsEditor
              v-model="formData.game_user_settings"
              :server="server"
            />
          </div>

          <!-- Game.ini 选项卡 -->
          <div v-if="activeTab === 'game_ini'" class="space-y-4">
            <GameIniEditor
              v-model="formData.game_ini"
              :server="server"
            />
          </div>


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
  mode: {
    type: String, // 'create' 或 'edit'
    default: 'create'
  },
  server: {
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
const activeTab = ref('basic')
const showFormPassword = ref(false)

// 表单数据
const formData = ref({
  identifier: '',
  port: 7777,
  query_port: 27015,
  rcon_port: 27020,
  admin_password: '',
  map: 'TheIsland',
  game_user_settings: '',
  game_ini: ''
})

// 切换密码显示
const toggleFormPassword = () => {
  showFormPassword.value = !showFormPassword.value
}

// 监听 show 变化，重置状态
watch(() => props.show, (newShow) => {
  if (newShow) {
    activeTab.value = 'basic'
    showFormPassword.value = false
    
    if (props.mode === 'create') {
      // 创建模式：重置表单
      formData.value = {
        identifier: '',
        port: 7777,
        query_port: 27015,
        rcon_port: 32330,
        admin_password: '',
        map: 'TheIsland',
        game_user_settings: '',
        game_ini: ''
      }
    } else if (props.mode === 'edit' && props.server) {
      // 编辑模式：填充当前服务器数据
      formData.value = {
        identifier: props.server.identifier || '',
        port: props.server.port || 7777,
        query_port: props.server.query_port || 27015,
        rcon_port: props.server.rcon_port || 32330,
        admin_password: props.server.admin_password || '',
        map: props.server.map || 'TheIsland',
        game_user_settings: props.server.game_user_settings || '',
        game_ini: props.server.game_ini || ''
      }
    }
  }
})

// 监听服务器数据变化（编辑模式下动态更新）
watch(() => props.server, (newServer) => {
  if (props.mode === 'edit' && newServer && props.show) {
    formData.value = {
      identifier: newServer.identifier || '',
      port: newServer.port || 7777,
              query_port: newServer.query_port || 27015,
        rcon_port: newServer.rcon_port || 32330,
        admin_password: newServer.admin_password || '',
      map: newServer.map || 'TheIsland',
      game_user_settings: newServer.game_user_settings || '',
      game_ini: newServer.game_ini || ''
    }
  }
}, { deep: true })

// 保存服务器
const saveServer = () => {
  emit('save', { ...formData.value })
}
</script> 