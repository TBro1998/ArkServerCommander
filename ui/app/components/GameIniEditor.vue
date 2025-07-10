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
          <i class="fas fa-sliders-h mr-2"></i>{{ t('servers.editor.visualEdit') }}
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
          <i class="fas fa-code mr-2"></i>{{ t('servers.editor.textEdit') }}
        </button>
      </div>
      <div class="flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:gap-4">
        <div class="text-sm text-gray-500">
          Game.ini - {{ editMode === 'visual' ? t('servers.editor.visualEditMode') : t('servers.editor.textEditMode') }}
        </div>
        <div v-if="isSyncing" class="flex items-center text-sm text-blue-600">
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 mr-2"></div>
          {{ t('servers.editor.syncing') }}
        </div>
        <div v-else class="flex items-center text-sm text-green-600">
          <i class="fas fa-check-circle mr-1"></i>
          {{ t('servers.editor.synced') }}
        </div>
      </div>
    </div>

    <!-- 可视化编辑模式 -->
    <div v-if="editMode === 'visual'" class="space-y-6">
      <!-- 页签导航 -->
      <div class="border-b border-gray-200">
        <nav class="-mb-px flex space-x-8 overflow-x-auto">
          <button
            v-for="category in getAllCategories()"
            :key="category.key"
            @click="activeTab = category.key"
            :class="[
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm',
              activeTab === category.key
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            ]"
          >
            {{ getCategoryName(category.key) }}
          </button>
        </nav>
      </div>

      <!-- 页签内容 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <div
          v-for="category in getAllCategories()"
          :key="category.key"
          v-show="activeTab === category.key"
          class="space-y-4"
        >
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div 
              v-for="{ key: paramKey, param } in category.params" 
              :key="paramKey"
              class="space-y-2 relative"
            >
              <label class="block text-sm font-medium text-gray-700">
                {{ getParamDisplayName(paramKey) }}
                <span 
                  class="relative inline-block ml-1"
                  @mouseenter="showTooltip = paramKey"
                  @mouseleave="showTooltip = null"
                >
                  <i class="fas fa-info-circle text-gray-400 cursor-help"></i>
                  <div 
                    v-if="showTooltip === paramKey"
                    class="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                    style="bottom: 100%; left: 50%; transform: translateX(-50%); margin-bottom: 8px;"
                  >
                    {{ getParamDisplayName(paramKey) }}
                    <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
                  </div>
                </span>
              </label>
              
              <!-- Boolean 类型 -->
              <div v-if="param.type === 'boolean'">
                <ToggleSwitch
                  :id="paramKey"
                  v-model="visualConfig[paramKey]"
                  :label="visualConfig[paramKey] ? t('servers.editor.enabled') : t('servers.editor.disabled')"
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
              
              <!-- Text 类型 -->
              <input
                v-else
                v-model="visualConfig[paramKey]"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 文本编辑模式 -->
    <div v-else>
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-2 gap-2">
        <label class="block text-sm font-medium text-gray-700">
          Game.ini {{ t('servers.editor.content') }}
        </label>
        <div class="flex gap-2">
          <button
            @click="resetToDefault"
            class="text-sm text-blue-600 hover:text-blue-800"
          >
            {{ t('servers.editor.resetToDefault') }}
          </button>
          <button
            @click="formatConfig"
            class="text-sm text-green-600 hover:text-green-800"
          >
            {{ t('servers.editor.format') }}
          </button>
        </div>
      </div>
      <textarea
        v-model="textContent"
        class="w-full h-64 sm:h-96 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
        :placeholder="t('servers.editor.placeholder')"
      ></textarea>
      <p class="text-xs text-gray-500 mt-1">
        {{ t('servers.editor.description') }}
      </p>
    </div>

    <!-- 提示信息 -->
    <div v-if="editMode === 'visual'" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
      <div class="flex flex-col sm:flex-row items-start">
        <i class="fas fa-info-circle text-blue-500 mr-0 sm:mr-3 mb-2 sm:mb-0"></i>
        <div class="text-sm text-blue-700">
          <p class="font-medium">{{ t('servers.editor.visualEditModeDesc') }}</p>
          <p>{{ t('servers.editor.visualEditModeTip') }}</p>
        </div>
      </div>
    </div>

    <div v-if="editMode === 'text'" class="bg-green-50 border border-green-200 rounded-lg p-4">
      <div class="flex flex-col sm:flex-row items-start">
        <i class="fas fa-code text-green-500 mr-0 sm:mr-3 mb-2 sm:mb-0"></i>
        <div class="text-sm text-green-700">
          <p class="font-medium">{{ t('servers.editor.textEditMode') }}</p>
          <p>{{ t('servers.editor.gameIniTextEditDesc') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { gameIniParams, getAllCategories } from '../utils/gameIniParams'
import { extractConfigValues, formatConfigContent, mergeConfigContent } from '../utils/configParser'
import ToggleSwitch from './ToggleSwitch.vue'

// i18n
const { t } = useI18n()

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
const isInitialized = ref(false)
const activeTab = ref('')

// 防抖定时器
let visualSyncTimer = null
let textSyncTimer = null

// i18n 相关函数
const getCategoryName = (categoryKey) => {
  const translated = t(`servers.gameIniCategories.${categoryKey}`)
  return translated !== `servers.gameIniCategories.${categoryKey}` ? translated : categoryKey
}

const getParamDisplayName = (paramKey) => {
  const translated = t(`servers.gameIniParams.${paramKey}`)
  return translated !== `servers.gameIniParams.${paramKey}` ? translated : paramKey
}

// 初始化可视化配置默认值
const initializeVisualConfig = () => {
  try {
    if (!gameIniParams) {
      console.warn('Game.ini 配置参数尚未加载，跳过初始化')
      return
    }

    // 只在未初始化时设置默认值
    if (!isInitialized.value) {
      Object.keys(gameIniParams).forEach(sectionKey => {
        const section = gameIniParams[sectionKey]
        if (section) {
          Object.keys(section).forEach(paramKey => {
            const param = section[paramKey]
            if (param && param.default !== undefined) {
              visualConfig.value[paramKey] = param.default
            }
          })
        }
      })
      isInitialized.value = true
      console.log('Game.ini 可视化配置初始化完成')
    }
  } catch (error) {
    console.error('初始化 Game.ini 可视化配置失败:', error)
  }
}

// 解析文本配置到可视化配置
const parseTextToVisual = () => {
  try {
    if (!gameIniParams) {
      console.warn('Game.ini 配置参数尚未加载')
      return
    }

    if (textContent.value) {
      console.log('开始解析 Game.ini 配置文件到可视化配置')
      const values = extractConfigValues(
        textContent.value, 
        gameIniParams
      )
      if (values && Object.keys(values).length > 0) {
        // 合并解析的值，保留未解析参数的默认值
        Object.keys(values).forEach(key => {
          if (values[key] !== undefined && values[key] !== null) {
            visualConfig.value[key] = values[key]
          }
        })
        console.log('Game.ini 配置文件解析完成，更新了', Object.keys(values).length, '个参数')
      }
    }
  } catch (error) {
    console.error('解析 Game.ini 配置文件失败:', error)
  }
}

// 将可视化配置同步到文本
const syncVisualToText = () => {
  try {
    if (!gameIniParams) {
      console.warn('Game.ini 配置参数尚未加载，无法同步')
      return
    }

    // 使用智能合并，保留用户自定义的配置
    const mergedContent = mergeConfigContent(
      textContent.value,
      visualConfig.value,
      gameIniParams,
      ['/script/shootergame.shootergamemode', '/script/engine.gamesession']
    )

    textContent.value = mergedContent
    emit('update:modelValue', mergedContent)
  } catch (error) {
    console.error('同步 Game.ini 可视化配置到文本失败:', error)
  }
}

// 监听 props.modelValue 变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== textContent.value) {
    textContent.value = newValue || ''
    // 只有在 gameIniParams 已加载且处于可视化模式时才解析
    if (gameIniParams && editMode.value === 'visual') {
      parseTextToVisual()
    }
  }
}, { immediate: true })

// 监听 gameIniParams 的变化（处理异步加载）
watch(() => gameIniParams, (newParams) => {
  if (newParams) {
    console.log('Game.ini 参数已加载，重新初始化')
    initializeVisualConfig()
    // 如果已经有文本内容且处于可视化模式，解析配置
    if (textContent.value && editMode.value === 'visual') {
      parseTextToVisual()
    }
  }
}, { immediate: true })

// 监听可视化配置变化，自动同步到文本
watch(() => visualConfig.value, () => {
  if (editMode.value === 'visual' && isInitialized.value) {
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
    // 确保在切换时解析文本配置
    if (textContent.value && gameIniParams) {
      parseTextToVisual()
    }
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
  
  textContent.value = `[/script/shootergame.shootergamemode]
# 基础游戏设置
bUseSingleplayerSettings=false
bDisableStructurePlacementCollision=false
bAllowFlyerCarryPvE=true
bAllowUnlimitedRespecs=true
bPassiveDefensesDamageRiderlessDinos=true
bAllowCustomRecipes=true
bUseCorpseLocator=false
bDisableLootCrates=false
bPvEDisableFriendlyFire=false
bDisableFriendlyFire=false

# 经验值和等级设置
OverrideMaxExperiencePointsPlayer=0
OverrideMaxExperiencePointsDino=0
OverridePlayerLevelEngramPoints=0
KillXPMultiplier=1.0
HarvestXPMultiplier=1.0
CraftXPMultiplier=1.0
GenericXPMultiplier=1.0
SpecialXPMultiplier=1.0

# 繁殖设置
MatingIntervalMultiplier=1.0
MatingSpeedMultiplier=1.0
EggHatchSpeedMultiplier=1.0
BabyMatureSpeedMultiplier=1.0
BabyCuddleIntervalMultiplier=1.0
BabyFoodConsumptionSpeedMultiplier=1.0
BabyImprintingStatScaleMultiplier=1.0
BabyImprintAmountMultiplier=1.0
BabyCuddleGracePeriodMultiplier=1.0
BabyCuddleLoseImprintQualitySpeedMultiplier=1.0

# 物品和资源设置
CropGrowthSpeedMultiplier=1.0
CropDecaySpeedMultiplier=1.0
GlobalSpoilingTimeMultiplier=1.0
GlobalItemDecompositionTimeMultiplier=1.0
GlobalCorpseDecompositionTimeMultiplier=1.0
UseCorpseLifeSpanMultiplier=1.0
SupplyCrateLootQualityMultiplier=1.0
FishingLootQualityMultiplier=1.0
CustomRecipeEffectivenessMultiplier=1.0
CustomRecipeSkillMultiplier=1.0
CraftingSkillBonusMultiplier=1.0
ResourceNoReplenishRadiusPlayers=1.0
ResourceNoReplenishRadiusStructures=1.0
DinoHarvestingDamageMultiplier=1.0
PlayerHarvestingDamageMultiplier=1.0

# 恐龙设置
TamedDinoCharacterFoodDrainMultiplier=1.0
WildDinoCharacterFoodDrainMultiplier=1.0
WildDinoTorporDrainMultiplier=1.0
TamedDinoTorporDrainMultiplier=1.0
PassiveTameIntervalMultiplier=1.0
DinoTurretDamageMultiplier=1.0
DestroyTamesOverLevelClamp=0

# 部落和玩家设置
MaxNumberOfPlayersInTribe=0
MaxAlliancesPerTribe=0
MaxTribesPerAlliance=0
TribeSlotReuseCooldown=0
MaxTribeLogs=100
KickIdlePlayersPeriod=0

# PvP设置
IncreasePvPRespawnIntervalCheckPeriod=300
IncreasePvPRespawnIntervalMultiplier=2.0
IncreasePvPRespawnIntervalBaseAmount=60
AutoPvEStartTimeSeconds=0
AutoPvEStopTimeSeconds=0
PvPZoneStructureDamageMultiplier=1.0

# 建筑和结构设置
StructureDamageRepairCooldown=180
StructureResistanceMultiplier=1.0
StructureDamageMultiplier=1.0
PvEStructureDecayPeriodMultiplier=1.0
FastDecayInterval=0

# 高级设置
HairGrowthSpeedMultiplier=1.0
PoopIntervalMultiplier=1.0
LayEggIntervalMultiplier=1.0
FuelConsumptionIntervalMultiplier=1.0
GlobalPoweredBatteryDurabilityDecreasePerSecond=0
LimitNonPlayerDroppedItemsRange=0
LimitNonPlayerDroppedItemsCount=0
MaxFallSpeedMultiplier=1.0
PreventOfflinePvPConnectionInvincibleInterval=0

[/script/engine.gamesession]
MaxPlayers=70

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

# PvP相关设置
bEnablePvPGamma=false`

  emit('update:modelValue', textContent.value)
}

// 格式化配置文件内容
const formatConfig = () => {
  textContent.value = formatConfigContent(textContent.value)
  emit('update:modelValue', textContent.value)
}

// 组件挂载时初始化
onMounted(() => {
  initializeVisualConfig()
  
  // 初始化第一个页签
  const categories = getAllCategories()
  if (categories.length > 0 && !activeTab.value) {
    activeTab.value = categories[0].key
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