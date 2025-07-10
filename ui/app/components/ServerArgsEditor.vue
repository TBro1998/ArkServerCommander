<template>
  <div class="space-y-4">
    <!-- 标题栏 -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
      <div class="flex flex-wrap gap-2 sm:gap-4">
        <div class="text-sm text-gray-500">
          {{ t('servers.argsEditor.title') }}
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="flex flex-wrap gap-2 sm:gap-4">
      <button
        @click="resetToDefaults"
        class="px-3 py-2 rounded-md text-sm font-medium bg-gray-200 text-gray-700 hover:bg-gray-300 transition-colors"
      >
                    <i class="fas fa-undo mr-2"></i>{{ t('servers.argsEditor.resetToDefault') }}
      </button>
    </div>

    <!-- 参数分类标签 -->
    <div class="border-b border-gray-200">
      <nav class="-mb-px flex space-x-4 sm:space-x-8 overflow-x-auto">
        <button
          v-for="category in visibleCategories"
          :key="category.key"
          @click="activeParamCategory = category.key"
          :class="[
            'py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap',
            activeParamCategory === category.key
              ? 'border-blue-500 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          ]"
        >
          {{ category.name }}
        </button>
      </nav>
    </div>

    <!-- 参数编辑区域 -->
    <div class="space-y-6">
      <div 
        v-for="category in visibleCategories" 
        :key="category.key"
        v-show="activeParamCategory === category.key"
        class="bg-gray-50 rounded-lg p-4"
      >
        <h4 class="text-base sm:text-lg font-semibold text-gray-800 mb-4 border-b border-gray-200 pb-2">
          {{ category.name }}
          <span class="text-sm font-normal text-gray-500 ml-2">{{ category.description }}</span>
        </h4>
        
        <!-- 预定义参数 -->
        <div v-if="category.key !== 'custom'" class="space-y-6">
          <!-- 布尔值参数组 -->
          <div v-if="category.params.filter(p => p.param.type === 'boolean').length > 0" class="space-y-4">
            <h5 class="text-sm font-medium text-gray-700 border-b border-gray-200 pb-2">{{ t('servers.argsEditor.switchParams') }}</h5>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              <div 
                v-for="param in category.params.filter(p => p.param.type === 'boolean')" 
                :key="param.key"
                class="space-y-2 relative"
              >
                <label class="block text-sm font-medium text-gray-700">
                  {{ getParamDisplayName(param.key, param.param) }}
                  <span 
                    class="relative inline-block ml-1"
                    @mouseenter="showTooltip = param.key"
                    @mouseleave="showTooltip = null"
                  >
                    <i class="fas fa-info-circle text-gray-400 cursor-help" :title="param.key"></i>
                    <div 
                      v-if="showTooltip === param.key"
                      class="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                      style="bottom: 100%; left: 50%; transform: translateX(-50%); margin-bottom: 8px;"
                    >
                      <div class="font-medium mb-1">{{ param.key }}</div>
                      <div class="text-gray-300">{{ t('servers.argsEditor.switchParams') }}</div>
                      <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
                    </div>
                  </span>
                </label>
                <ToggleSwitch
                  :id="param.key"
                  v-model="paramValues[param.key]"
                  :label="paramValues[param.key] ? t('servers.argsEditor.enabled') : t('servers.argsEditor.disabled')"
                />
              </div>
            </div>
          </div>

          <!-- 数值参数组 -->
          <div v-if="category.params.filter(p => p.param.type === 'number').length > 0" class="space-y-4">
            <h5 class="text-sm font-medium text-gray-700 border-b border-gray-200 pb-2">{{ t('servers.argsEditor.numberParams') }}</h5>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              <div 
                v-for="param in category.params.filter(p => p.param.type === 'number')" 
                :key="param.key"
                class="space-y-2 relative"
              >
                <label class="block text-sm font-medium text-gray-700">
                  {{ getParamDisplayName(param.key, param.param) }}
                  <span 
                    class="relative inline-block ml-1"
                    @mouseenter="showTooltip = param.key"
                    @mouseleave="showTooltip = null"
                  >
                    <i class="fas fa-info-circle text-gray-400 cursor-help" :title="param.key"></i>
                    <div 
                      v-if="showTooltip === param.key"
                      class="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                      style="bottom: 100%; left: 50%; transform: translateX(-50%); margin-bottom: 8px;"
                    >
                      <div class="font-medium mb-1">{{ param.key }}</div>
                      <div class="text-gray-300">{{ t('servers.argsEditor.numberParams') }}</div>
                      <div v-if="param.param.min !== undefined && param.param.max !== undefined" class="text-gray-300 mt-1">
                        {{ t('servers.argsEditor.range') }}: {{ param.param.min }}-{{ param.param.max }}
                      </div>
                      <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
                    </div>
                  </span>
                </label>
                <input
                  v-model.number="paramValues[param.key]"
                  type="number"
                  :min="param.param.min"
                  :max="param.param.max"
                  :step="param.param.step || 1"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  :placeholder="param.param.default?.toString()"
                />
              </div>
            </div>
          </div>

          <!-- 文本参数组 -->
          <div v-if="category.params.filter(p => p.param.type === 'string' || p.param.type === 'select').length > 0" class="space-y-4">
            <h5 class="text-sm font-medium text-gray-700 border-b border-gray-200 pb-2">{{ t('servers.argsEditor.textParams') }}</h5>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              <div 
                v-for="param in category.params.filter(p => p.param.type === 'string' || p.param.type === 'select')" 
                :key="param.key"
                class="space-y-2 relative"
              >
                <label class="block text-sm font-medium text-gray-700">
                  {{ getParamDisplayName(param.key, param.param) }}
                  <span 
                    class="relative inline-block ml-1"
                    @mouseenter="showTooltip = param.key"
                    @mouseleave="showTooltip = null"
                  >
                    <i class="fas fa-info-circle text-gray-400 cursor-help" :title="param.key"></i>
                    <div 
                      v-if="showTooltip === param.key"
                      class="absolute z-30 bg-gray-900 text-white text-xs rounded py-2 px-3 max-w-xs whitespace-normal shadow-lg"
                      style="bottom: 100%; left: 50%; transform: translateX(-50%); margin-bottom: 8px;"
                    >
                      <div class="font-medium mb-1">{{ param.key }}</div>
                      <div class="text-gray-300">{{ param.param.options ? t('servers.argsEditor.selectParams') : t('servers.argsEditor.textParams') }}</div>
                      <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-900"></div>
                    </div>
                  </span>
                </label>
                <template v-if="param.param.options">
                  <select
                    v-model="paramValues[param.key]"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">{{ param.param.default || t('servers.argsEditor.pleaseSelect') }}</option>
                    <option v-for="option in param.param.options" :key="option" :value="option">{{ option }}</option>
                  </select>
                </template>
                <input
                  v-else
                  v-model="paramValues[param.key]"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  :placeholder="param.param.default || ''"
                />
              </div>
            </div>
          </div>
        </div>
        
        <!-- 自定义参数 -->
        <div v-if="category.key === 'custom'" class="space-y-4">
          <div class="text-sm text-gray-600 mb-3">
            {{ t('servers.argsEditor.customArgsDesc') }}
          </div>
          <div v-for="(arg, idx) in customArgs" :key="'custom-'+idx" class="flex gap-2">
            <input 
              v-model="customArgs[idx]" 
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
              :placeholder="t('servers.argsEditor.customArgPlaceholder')" 
            />
            <button 
              @click="deleteCustomArg(idx)" 
              class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-md transition-colors"
              title="删除此参数"
            >
              <i class="fas fa-trash"></i>
            </button>
          </div>
          <button 
            @click="addCustomArg" 
            class="text-blue-600 hover:text-blue-800 text-sm font-medium transition-colors duration-200"
          >
            <i class="fas fa-plus mr-1"></i>{{ t('servers.argsEditor.addCustomArg') }}
          </button>
        </div>
      </div>
    </div>


    

  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import ToggleSwitch from './ToggleSwitch.vue'
import { queryParams, commandLineArgs, categories, getParamsByCategory, getDefaultValues } from '~/utils/arkServerParams'

// i18n
const { t } = useI18n()

// Props
const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({ query_params: {}, command_line_args: {}, custom_args: [] })
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 本地状态
const activeParamCategory = ref('basic')
const showTooltip = ref(null)

// 参数值存储
const paramValues = ref({})
const customArgs = ref([])

// 获取参数分类
const paramsByCategory = computed(() => getParamsByCategory())

// 可见的分类（显示所有参数）
const visibleCategories = computed(() => {
  const categories = Object.entries(paramsByCategory.value).map(([key, params]) => ({
    key,
    name: getCategoryName(key),
    description: getCategoryDescription(key),
    params: params
  }))
  
  return categories.filter(cat => cat.params.length > 0 || cat.key === 'custom')
})

// 获取分类名称
function getCategoryName(key) {
  const translation = t(`servers.paramCategories.${key}`, key)
  return translation !== `servers.paramCategories.${key}` ? translation : key
}

// 获取分类描述
function getCategoryDescription(key) {
  const translation = t(`servers.paramCategories.${key}`, '')
  return translation !== `servers.paramCategories.${key}` ? translation : ''
}

// 获取参数显示名称
function getParamDisplayName(key, param) {
  // 首先检查是否是查询参数
  if (queryParams[key]) {
    const translation = t(`servers.queryParams.${key}`, key)
    return translation !== `servers.queryParams.${key}` ? translation : key
  }
  // 然后检查是否是命令行参数
  if (commandLineArgs[key]) {
    const translation = t(`servers.commandLineArgs.${key}`, key)
    return translation !== `servers.commandLineArgs.${key}` ? translation : key
  }
  return key
}

// 初始化参数值
function initializeParamValues() {
  // 初始化查询参数
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      const value = props.modelValue.query_params[key]
      if (param.type === 'boolean') {
        paramValues.value[key] = value === 'True' || value === true
      } else {
        paramValues.value[key] = value || param.default || ''
      }
    }
  })
  
  // 初始化命令行参数
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      const value = props.modelValue.command_line_args[key]
      if (param.type === 'boolean') {
        paramValues.value[key] = value === true || value === 'true'
      } else {
        paramValues.value[key] = value || param.default || ''
      }
    }
  })
  
  // 初始化自定义参数
  customArgs.value = [...(props.modelValue.custom_args || [])]
}

// 重置为默认值
function resetToDefaults() {
  paramValues.value = {}
  
  // 重置查询参数
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      paramValues.value[key] = param.default
    }
  })
  
  // 重置命令行参数
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      paramValues.value[key] = param.default
    }
  })
  
  // 重置自定义参数
  customArgs.value = []
  
  // 触发更新
  updateModelValue()
}

// 添加自定义参数
function addCustomArg() {
  customArgs.value.push('')
  updateModelValue()
}

// 删除自定义参数
function deleteCustomArg(idx) {
  customArgs.value.splice(idx, 1)
  updateModelValue()
}

// 更新模型值
function updateModelValue() {
  const newValue = {
    query_params: {},
    command_line_args: {},
    custom_args: customArgs.value.filter(arg => arg.trim() !== '')
  }
  
  // 更新查询参数
  Object.entries(queryParams).forEach(([key, param]) => {
    if (!param.hidden) {
      const value = paramValues.value[key]
      if (param.type === 'boolean') {
        newValue.query_params[key] = value ? 'True' : 'False'
      } else {
        newValue.query_params[key] = String(value || '')
      }
    }
  })
  
  // 更新命令行参数
  Object.entries(commandLineArgs).forEach(([key, param]) => {
    if (!param.hidden) {
      const value = paramValues.value[key]
      newValue.command_line_args[key] = value
    }
  })
  
  emit('update:modelValue', newValue)
}

// 监听参数值变化
watch(paramValues, updateModelValue, { deep: true })

// 监听自定义参数变化
watch(customArgs, updateModelValue, { deep: true })

// 监听props变化（避免初始化时的循环）
watch(() => props.modelValue, (newValue, oldValue) => {
  if (JSON.stringify(newValue) !== JSON.stringify(oldValue)) {
    initializeParamValues()
  }
}, { deep: true })

// 组件挂载时初始化
onMounted(() => {
  initializeParamValues()
})
</script> 