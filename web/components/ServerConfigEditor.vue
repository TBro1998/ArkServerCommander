<template>
  <div
    v-if="show"
    class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    @click="$emit('close')"
  >
    <div class="relative top-10 mx-auto p-5 border w-full max-w-6xl shadow-lg rounded-md bg-white" @click.stop>
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
          <!-- GameUserSettings.ini -->
          <div v-if="activeTab === 'game_user_settings'">
            <div class="flex justify-between items-center mb-2">
              <label class="block text-sm font-medium text-gray-700">
                GameUserSettings.ini 内容
              </label>
              <div class="flex gap-2">
                <button
                  @click="resetToDefault('game_user_settings')"
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  重置为默认
                </button>
                <button
                  @click="formatConfig('game_user_settings')"
                  class="text-sm text-green-600 hover:text-green-800"
                >
                  格式化
                </button>
              </div>
            </div>
            <textarea
              v-model="localConfigData.game_user_settings"
              class="w-full h-96 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="输入 GameUserSettings.ini 配置内容..."
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">
              此文件包含服务器的基本设置，如端口、密码、最大玩家数等
            </p>
          </div>

          <!-- Game.ini -->
          <div v-if="activeTab === 'game_ini'">
            <div class="flex justify-between items-center mb-2">
              <label class="block text-sm font-medium text-gray-700">
                Game.ini 内容
              </label>
              <div class="flex gap-2">
                <button
                  @click="resetToDefault('game_ini')"
                  class="text-sm text-blue-600 hover:text-blue-800"
                >
                  重置为默认
                </button>
                <button
                  @click="formatConfig('game_ini')"
                  class="text-sm text-green-600 hover:text-green-800"
                >
                  格式化
                </button>
              </div>
            </div>
            <textarea
              v-model="localConfigData.game_ini"
              class="w-full h-96 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              placeholder="输入 Game.ini 配置内容..."
            ></textarea>
            <p class="text-xs text-gray-500 mt-1">
              此文件包含游戏规则设置，如倍率、难度、PvP设置等
            </p>
          </div>
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

// 重置为默认配置
const resetToDefault = (configType) => {
  if (!props.server) return
  
  const server = props.server
  
  if (configType === 'game_user_settings') {
    localConfigData.value.game_user_settings = `[ServerSettings]
SessionName=${server.name}
ServerPassword=
ServerAdminPassword=${server.admin_password}
Port=${server.port}
QueryPort=${server.query_port}
RCONEnabled=True
RCONPort=${server.rcon_port}
MaxPlayers=${server.max_players}

[SessionSettings]
SessionName=${server.name}

[MessageOfTheDay]
Message=欢迎来到 ${server.name} ARK 服务器！

[/Script/ShooterGame.ShooterGameMode]
bUseSingleplayerSettings=False
bDisableStructurePlacementCollision=False
bAllowFlyerCarryPvE=True
bDisableStructureDecayPvE=False

[RCONSettings]
RCONEnabled=True
RCONPort=${server.rcon_port}`
  } else if (configType === 'game_ini') {
    localConfigData.value.game_ini = `[/script/shootergame.shootergamemode]
bUseSingleplayerSettings=false
bDisableStructurePlacementCollision=false
bAllowFlyerCarryPvE=true
bDisableStructureDecayPvE=false
bAllowUnlimitedRespecs=true
bAllowPlatformSaddleMultiFloors=true
bPassiveDefensesDamageRiderlessDinos=true
MaxNumberOfPlayersInTribe=0

[/script/engine.gamesession]
MaxPlayers=${server.max_players}

[/Script/ShooterGame.ShooterGameMode]
DifficultyOffset=1.0
OverrideOfficialDifficulty=5.0

# 资源重生倍率
ResourcesRespawnPeriodMultiplier=1.0

# 驯服相关设置
TamingSpeedMultiplier=1.0
DinoCharacterFoodDrainMultiplier=1.0
DinoCharacterStaminaDrainMultiplier=1.0
DinoCharacterHealthRecoveryMultiplier=1.0
DinoCountMultiplier=1.0

# 经验值倍率
XPMultiplier=1.0
PlayerCharacterWaterDrainMultiplier=1.0
PlayerCharacterFoodDrainMultiplier=1.0
PlayerCharacterStaminaDrainMultiplier=1.0
PlayerCharacterHealthRecoveryMultiplier=1.0

# 收集倍率
HarvestAmountMultiplier=1.0
HarvestHealthMultiplier=1.0

# 白天/夜晚时间流逝速度
DayCycleSpeedScale=1.0
NightTimeSpeedScale=1.0

# 结构相关
StructureResistanceMultiplier=1.0
StructureDamageMultiplier=1.0
StructureDamageRepairCooldown=180
PvEStructureDecayPeriodMultiplier=1.0

# PvP相关设置
bPvEDisableFriendlyFire=False
bEnablePvPGamma=False
bDisableFriendlyFire=False
bAllowFlyerCarryPvE=True`
  }
}

// 格式化配置文件内容
const formatConfig = (configType) => {
  const configKey = configType
  const lines = localConfigData.value[configKey].split('\n')
  const formatted = []
  let lastWasEmpty = false
  
  for (const line of lines) {
    const trimmed = line.trim()
    
    if (trimmed === '') {
      if (!lastWasEmpty) {
        formatted.push('')
        lastWasEmpty = true
      }
    } else {
      formatted.push(trimmed)
      lastWasEmpty = false
    }
  }
  
  localConfigData.value[configKey] = formatted.join('\n')
}
</script> 