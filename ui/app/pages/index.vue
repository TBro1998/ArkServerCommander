<template>
  <UContainer class="py-12">
    <div class="text-center space-y-12">
      <!-- 欢迎标题区域 -->
      <div class="space-y-6">
        <div class="mx-auto w-20 h-20 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center mb-6">
          <UIcon name="i-lucide-server" class="w-10 h-10 text-white" />
        </div>
        <h1 class="text-5xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
          {{ $t('home.title') }}
        </h1>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
          {{ $t('home.subtitle') }}
        </p>
      </div>
      
      <!-- 系统信息卡片 -->
      <UCard class="max-w-lg mx-auto">
        <template #header>
          <div class="flex items-center space-x-3">
            <UIcon name="i-lucide-user" class="w-5 h-5 text-blue-600" />
            <h2 class="text-xl font-semibold">{{ $t('home.systemInfo') }}</h2>
          </div>
        </template>
        
        <div class="space-y-4">
          <div class="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0">
            <span class="font-medium text-gray-700">{{ $t('home.username') }}</span>
            <span class="text-gray-900 font-mono">{{ authStore.user?.username }}</span>
          </div>
          <div class="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0">
            <span class="font-medium text-gray-700">{{ $t('home.userID') }}</span>
            <span class="text-gray-900 font-mono">#{{ authStore.user?.id }}</span>
          </div>
        </div>
      </UCard>
      
      <!-- 功能卡片网格 -->
      <div class="space-y-8">
        <h3 class="text-2xl font-semibold text-gray-900">{{ $t('home.features') }}</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
          <!-- 服务器管理卡片 -->
          <NuxtLink to="/servers">
            <UCard class="h-full hover:shadow-xl transition-all duration-300 cursor-pointer group border-2 border-transparent hover:border-blue-200">
              <div class="text-center space-y-6 p-6">
                <div class="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center group-hover:bg-blue-200 transition-colors">
                  <UIcon 
                    name="i-lucide-server" 
                    class="w-8 h-8 text-blue-600 group-hover:text-blue-700"
                  />
                </div>
                <div>
                  <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ $t('home.serverManagement') }}</h3>
                  <p class="text-gray-600 leading-relaxed">{{ $t('home.serverManagementDesc') }}</p>
                </div>
                <UButton 
                  color="blue" 
                  variant="ghost" 
                  size="sm"
                  class="group-hover:bg-blue-50"
                >
                  {{ $t('home.startManage') }}
                  <UIcon name="i-lucide-arrow-right" class="w-4 h-4 ml-1" />
                </UButton>
              </div>
            </UCard>
          </NuxtLink>

          <!-- 占位卡片：玩家管理 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <div class="text-center space-y-6 p-6">
              <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
                <UIcon 
                  name="i-lucide-users" 
                  class="w-8 h-8 text-gray-400"
                />
              </div>
              <div>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.playerManagement') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.playerManagementDesc') }}</p>
              </div>
              <UBadge color="gray" variant="soft" size="sm">
                {{ $t('home.comingSoon') }}
              </UBadge>
            </div>
          </UCard>

          <!-- 占位卡片：日志监控 -->
          <UCard class="h-full opacity-60 border-2 border-gray-100">
            <div class="text-center space-y-6 p-6">
              <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
                <UIcon 
                  name="i-lucide-bar-chart-3" 
                  class="w-8 h-8 text-gray-400"
                />
              </div>
              <div>
                <h3 class="text-xl font-semibold text-gray-500 mb-2">{{ $t('home.logMonitoring') }}</h3>
                <p class="text-gray-400 leading-relaxed">{{ $t('home.logMonitoringDesc') }}</p>
              </div>
              <UBadge color="gray" variant="soft" size="sm">
                {{ $t('home.comingSoon') }}
              </UBadge>
            </div>
          </UCard>
        </div>
      </div>
      
      <!-- 提示信息 -->
      <div class="flex items-center justify-center space-x-2 text-gray-500">
        <UIcon name="i-lucide-lightbulb" class="w-4 h-4" />
        <span class="text-sm">{{ $t('home.tip') }}</span>
      </div>
    </div>
  </UContainer>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

onMounted(() => {
  authStore.initFromStorage()
})
</script> 