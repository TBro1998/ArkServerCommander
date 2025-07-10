<template>
  <div class="relative">
    <UButton
      color="gray"
      variant="ghost"
      size="sm"
      class="flex items-center space-x-2"
      @click="isOpen = !isOpen"
    >
      <UIcon name="i-lucide-globe" class="w-4 h-4" />
      <span class="hidden sm:inline">{{ currentLanguageLabel }}</span>
      <UIcon name="i-lucide-chevron-down" class="w-4 h-4" />
    </UButton>
    
    <!-- 语言选择下拉菜单 -->
    <div v-if="isOpen" class="absolute right-0 mt-2 w-32 bg-white rounded-md shadow-lg border border-gray-200 z-50">
      <div class="py-1">
        <button
          @click="selectLanguage('zh')"
          class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center space-x-2"
          :class="{ 'bg-blue-50 text-blue-700': currentLanguage === 'zh' }"
        >
          <UIcon name="i-lucide-flag" class="w-4 h-4" />
          <span>中文</span>
        </button>
        <button
          @click="selectLanguage('en')"
          class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center space-x-2"
          :class="{ 'bg-blue-50 text-blue-700': currentLanguage === 'en' }"
        >
          <UIcon name="i-lucide-flag" class="w-4 h-4" />
          <span>English</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const { locale, locales, setLocale } = useI18n()
const switchLocalePath = useSwitchLocalePath()

const isOpen = ref(false)
const currentLanguage = computed(() => locale.value)

// 当前语言标签
const currentLanguageLabel = computed(() => {
  const currentLocale = locales.value.find(l => l.code === locale.value)
  return currentLocale?.code === 'zh' ? '中文' : 'English'
})

// 选择语言
const selectLanguage = async (newLocale) => {
  try {
    console.log('切换语言到:', newLocale)
    
    // 使用 setLocale 方法确保正确加载语言文件
    await setLocale(newLocale)
    isOpen.value = false
    
    // 强制重新渲染以确保翻译生效
    await nextTick()
    
    console.log('语言切换完成:', newLocale)
  } catch (error) {
    console.error('语言切换失败:', error)
  }
}

// 点击外部关闭下拉菜单
onMounted(() => {
  const handleClickOutside = (event) => {
    const target = event.target
    if (!target.closest('.relative')) {
      isOpen.value = false
    }
  }
  
  document.addEventListener('click', handleClickOutside)
  
  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
  })
})
</script> 